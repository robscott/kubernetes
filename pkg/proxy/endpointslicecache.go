/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxy

import (
	"sort"

	discovery "k8s.io/api/discovery/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	utilproxy "k8s.io/kubernetes/pkg/proxy/util"
	utilnet "k8s.io/utils/net"
)

// EndpointSliceCache is used as a cache of EndpointSlice information
type EndpointSliceCache struct {
	slicesByService  map[types.NamespacedName]map[string]EndpointSliceInfo
	makeEndpointInfo makeEndpointFunc
	hostname         string
	isIPv6Mode       *bool
	recorder         record.EventRecorder
}

// EndpointSliceInfo contains just the attributes kube-proxy cares about.
// Used for caching. Intentionally small to limit memory util.
type EndpointSliceInfo struct {
	Ports     []discovery.EndpointPort
	Endpoints []EndpointInfo
}

// EndpointInfo contains just the attributes kube-proxy cares about.
// Used for caching. Intentionally small to limit memory util.
type EndpointInfo struct {
	Targets  []string
	Topology map[string]string
}

// NewEndpointSliceCache initializes an EndpointSliceCache
func NewEndpointSliceCache(hostname string, isIPv6Mode *bool, recorder record.EventRecorder, makeEndpointInfo makeEndpointFunc) *EndpointSliceCache {
	return &EndpointSliceCache{
		slicesByService:  map[types.NamespacedName]map[string]EndpointSliceInfo{},
		hostname:         hostname,
		isIPv6Mode:       isIPv6Mode,
		makeEndpointInfo: makeEndpointInfo,
		recorder:         recorder,
	}
}

// Update a slice in the cache
func (e *EndpointSliceCache) Update(endpointSlice *discovery.EndpointSlice) {
	serviceKey, sliceKey := EndpointSliceCacheKeys(endpointSlice)
	esInfo := EndpointSliceInfo{
		Ports:     endpointSlice.Ports,
		Endpoints: []EndpointInfo{},
	}
	for _, endpoint := range endpointSlice.Endpoints {
		if endpoint.Conditions.Ready {
			esInfo.Endpoints = append(esInfo.Endpoints, EndpointInfo{
				Targets:  endpoint.Targets,
				Topology: endpoint.Topology,
			})
		}
	}
	if _, exists := e.slicesByService[serviceKey]; !exists {
		e.slicesByService[serviceKey] = map[string]EndpointSliceInfo{}
	}
	e.slicesByService[serviceKey][sliceKey] = esInfo
}

// Delete a slice from the cache
func (e *EndpointSliceCache) Delete(endpointSlice *discovery.EndpointSlice) {
	serviceKey, sliceKey := EndpointSliceCacheKeys(endpointSlice)
	delete(e.slicesByService[serviceKey], sliceKey)
}

// GetEndpointsMap computes an EndpointsMap for a given service
func (e *EndpointSliceCache) GetEndpointsMap(serviceNN types.NamespacedName) EndpointsMap {
	endpointsMap := EndpointsMap{}
	endpointInfoByServicePort := map[ServicePortName]map[string]Endpoint{}
	endpointSlices := e.slicesByService[serviceNN]

	// We need to build a map of portname -> all ip:ports for that portname.
	for _, endpointSlice := range endpointSlices {
		for i := range endpointSlice.Ports {
			port := &endpointSlice.Ports[i]
			if *port.Port == int32(0) {
				klog.Warningf("ignoring invalid endpoint port %s", *port.Name)
				continue
			}

			svcPortName := ServicePortName{NamespacedName: serviceNN}
			if port.Name != nil {
				svcPortName.Port = *port.Name
			}
			if _, exists := endpointInfoByServicePort[svcPortName]; !exists {
				endpointInfoByServicePort[svcPortName] = map[string]Endpoint{}
			}

			// iterate through endpoints to add them to endpointInfoByServicePort
			for i := range endpointSlice.Endpoints {
				endpoint := &endpointSlice.Endpoints[i]
				if len(endpoint.Targets) < 1 {
					klog.Warningf("ignoring invalid endpoint port %s with empty host", *port.Name)
					continue
				}

				// Filter out the incorrect IP version case.
				// Any endpoint port that contains incorrect IP version will be ignored.
				if e.isIPv6Mode != nil && utilnet.IsIPv6String(endpoint.Targets[0]) != *e.isIPv6Mode {
					// Emit event on the corresponding service which had a different
					// IP version than the endpoint.
					utilproxy.LogAndEmitIncorrectIPVersionEvent(e.recorder, "endpointslice", endpoint.Targets[0], serviceNN.Name, serviceNN.Namespace, "")
					continue
				}

				isLocal := len(e.hostname) > 0 && endpoint.Topology["kubernetes.io/hostname"] == e.hostname
				endpointInfo := newBaseEndpointInfo(endpoint.Targets[0], int(*port.Port), isLocal)

				// This logic ensures we're deduping potential overlapping endpoints
				// isLocal should not vary between matching IPs, but if it does, we favor a true value here if it exists
				if _, exists := endpointInfoByServicePort[svcPortName][endpointInfo.IP()]; !exists || isLocal {
					if e.makeEndpointInfo != nil {
						endpointInfoByServicePort[svcPortName][endpointInfo.IP()] = e.makeEndpointInfo(endpointInfo)
					} else {
						endpointInfoByServicePort[svcPortName][endpointInfo.IP()] = endpointInfo
					}
				}
			}
		}
	}

	// transform endpointInfoByServicePort into an endpointsMap with sorted IPs
	for svcPortName, endpointInfoByIP := range endpointInfoByServicePort {
		if len(endpointInfoByIP) > 0 {
			endpointsMap[svcPortName] = []Endpoint{}
			for _, endpointInfo := range endpointInfoByIP {
				endpointsMap[svcPortName] = append(endpointsMap[svcPortName], endpointInfo)

			}
			// Ensure IPs are always returned in the same order to simplify diffing
			sort.Sort(byIP(endpointsMap[svcPortName]))

			if klog.V(3) {
				newEPList := []string{}
				for _, ep := range endpointsMap[svcPortName] {
					newEPList = append(newEPList, ep.String())
				}
				klog.Infof("Setting endpoints for %q to %+v", svcPortName, newEPList)
			}
		}
	}

	return endpointsMap
}

// EndpointSliceCacheKeys returns cache keys used for a given EndpointSlice
func EndpointSliceCacheKeys(endpointSlice *discovery.EndpointSlice) (types.NamespacedName, string) {
	if len(endpointSlice.OwnerReferences) > 0 {
		ownerRef := endpointSlice.OwnerReferences[0]
		return types.NamespacedName{Namespace: endpointSlice.Namespace, Name: ownerRef.Name}, endpointSlice.Name
	}
	return types.NamespacedName{}, endpointSlice.Name
}

type byIP []Endpoint

func (e byIP) Len() int {
	return len(e)
}
func (e byIP) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
func (e byIP) Less(i, j int) bool {
	return e[i].IP() < e[j].IP()
}
