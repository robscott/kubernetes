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
	"fmt"
	"testing"

	discovery "k8s.io/api/discovery/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilpointer "k8s.io/utils/pointer"
)

func TestGetEndpointsMap(t *testing.T) {
	testCases := map[string]struct {
		endpointSlices []*discovery.EndpointSlice
		hostname       string
		namespacedName types.NamespacedName
		expectedMap    map[ServicePortName][]*BaseEndpointInfo
	}{
		// A simple use case with 3 endpoints, all ready, and 2 ports
		"1 slice, 2 hosts, ports 80,443": {
			namespacedName: types.NamespacedName{Name: "svc1", Namespace: "ns1"},
			hostname:       "host1",
			endpointSlices: []*discovery.EndpointSlice{
				generateEndpointSlice("svc1", "ns1", 1, 3, 999, []string{"host1", "host2"}, []int32{80, 443}),
			},
			expectedMap: map[ServicePortName][]*BaseEndpointInfo{
				makeServicePortName("ns1", "svc1", "port-0"): {
					&BaseEndpointInfo{Endpoint: "10.0.1.1:80", IsLocal: false},
					&BaseEndpointInfo{Endpoint: "10.0.1.2:80", IsLocal: true},
					&BaseEndpointInfo{Endpoint: "10.0.1.3:80", IsLocal: false},
				},
				makeServicePortName("ns1", "svc1", "port-1"): {
					&BaseEndpointInfo{Endpoint: "10.0.1.1:443", IsLocal: false},
					&BaseEndpointInfo{Endpoint: "10.0.1.2:443", IsLocal: true},
					&BaseEndpointInfo{Endpoint: "10.0.1.3:443", IsLocal: false},
				},
			},
		},
		// 2 slices, both referencing the same port
		"2 slices, same port": {
			namespacedName: types.NamespacedName{Name: "svc1", Namespace: "ns1"},
			endpointSlices: []*discovery.EndpointSlice{
				generateEndpointSlice("svc1", "ns1", 1, 3, 999, []string{}, []int32{80}),
				generateEndpointSlice("svc1", "ns1", 2, 3, 999, []string{}, []int32{80}),
			},
			expectedMap: map[ServicePortName][]*BaseEndpointInfo{
				makeServicePortName("ns1", "svc1", "port-0"): {
					&BaseEndpointInfo{Endpoint: "10.0.1.1:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.2:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.3:80"},
					&BaseEndpointInfo{Endpoint: "10.0.2.1:80"},
					&BaseEndpointInfo{Endpoint: "10.0.2.2:80"},
					&BaseEndpointInfo{Endpoint: "10.0.2.3:80"},
				},
			},
		},
		// 2 slices, with some overlapping endpoints, result should be a union of the 2
		"2 overlapping slices, same port": {
			namespacedName: types.NamespacedName{Name: "svc1", Namespace: "ns1"},
			endpointSlices: []*discovery.EndpointSlice{
				generateEndpointSlice("svc1", "ns1", 1, 3, 999, []string{}, []int32{80}),
				generateEndpointSlice("svc1", "ns1", 1, 4, 999, []string{}, []int32{80}),
			},
			expectedMap: map[ServicePortName][]*BaseEndpointInfo{
				makeServicePortName("ns1", "svc1", "port-0"): {
					&BaseEndpointInfo{Endpoint: "10.0.1.1:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.2:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.3:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.4:80"},
				},
			},
		},
		// 2 slices with all endpoints overlapping, more unready in first than second.
		// If an endpoint is marked ready, we add it to the EndpointsMap, even if
		//   conditions.Ready isn't true for another matching endpoint
		"2 slices, overlapping endpoints, some endpoints unready in 1 or both": {
			namespacedName: types.NamespacedName{Name: "svc1", Namespace: "ns1"},
			endpointSlices: []*discovery.EndpointSlice{
				generateEndpointSlice("svc1", "ns1", 1, 10, 3, []string{}, []int32{80}),
				generateEndpointSlice("svc1", "ns1", 1, 10, 6, []string{}, []int32{80}),
			},
			expectedMap: map[ServicePortName][]*BaseEndpointInfo{
				makeServicePortName("ns1", "svc1", "port-0"): {
					&BaseEndpointInfo{Endpoint: "10.0.1.1:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.10:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.2:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.3:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.4:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.5:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.7:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.8:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.9:80"},
				},
			},
		},
		// 2 slices with all endpoints unready, result should be empty endpointsMap
		"2 slices, overlapping endpoints, all unready": {
			namespacedName: types.NamespacedName{Name: "svc1", Namespace: "ns1"},
			endpointSlices: []*discovery.EndpointSlice{
				generateEndpointSlice("svc1", "ns1", 1, 10, 1, []string{}, []int32{80}),
				generateEndpointSlice("svc1", "ns1", 1, 10, 1, []string{}, []int32{80}),
			},
			expectedMap: map[ServicePortName][]*BaseEndpointInfo{},
		},
		// Verifying only the appropriate caches are updated in EndpointsMap
		"3 slices with different services and namespaces": {
			namespacedName: types.NamespacedName{Name: "svc1", Namespace: "ns1"},
			endpointSlices: []*discovery.EndpointSlice{
				generateEndpointSlice("svc1", "ns1", 1, 3, 999, []string{}, []int32{80}),
				generateEndpointSlice("svc2", "ns1", 2, 3, 999, []string{}, []int32{80}),
				generateEndpointSlice("svc1", "ns2", 3, 3, 999, []string{}, []int32{80}),
			},
			expectedMap: map[ServicePortName][]*BaseEndpointInfo{
				makeServicePortName("ns1", "svc1", "port-0"): {
					&BaseEndpointInfo{Endpoint: "10.0.1.1:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.2:80"},
					&BaseEndpointInfo{Endpoint: "10.0.1.3:80"},
				},
			},
		},
	}

	for name, tc := range testCases {
		esCache := NewEndpointSliceCache(tc.hostname, nil, nil, nil)

		for _, endpointSlice := range tc.endpointSlices {
			esCache.Update(endpointSlice)
		}

		compareEndpointsMapsStr(t, name, esCache.GetEndpointsMap(tc.namespacedName), tc.expectedMap)
	}
}

func generateEndpointSliceWithOffset(serviceName, namespace string, sliceNum, offset, numEndpoints, unreadyMod int, hosts []string, portNums []int32) *discovery.EndpointSlice {
	ipTargetType := discovery.IPTargetType
	endpointSlice := &discovery.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf("%s-%d", serviceName, sliceNum),
			Namespace:       namespace,
			OwnerReferences: []metav1.OwnerReference{{Kind: "Service", Name: serviceName}},
		},
		Ports:      []discovery.EndpointPort{},
		TargetType: &ipTargetType,
		Endpoints:  []discovery.Endpoint{},
	}

	for i, portNum := range portNums {
		endpointSlice.Ports = append(endpointSlice.Ports, discovery.EndpointPort{
			Name: utilpointer.StringPtr(fmt.Sprintf("port-%d", i)),
			Port: utilpointer.Int32Ptr(portNum),
		})
	}

	for i := 1; i <= numEndpoints; i++ {
		endpoint := discovery.Endpoint{
			Targets:    []string{fmt.Sprintf("10.0.%d.%d", offset, i)},
			Conditions: discovery.EndpointConditions{Ready: i%unreadyMod != 0},
		}

		if len(hosts) > 0 {
			endpoint.Topology = map[string]string{
				"kubernetes.io/hostname": hosts[i%len(hosts)],
			}
		}

		endpointSlice.Endpoints = append(endpointSlice.Endpoints, endpoint)
	}

	return endpointSlice
}

func generateEndpointSlice(serviceName, namespace string, sliceNum, numEndpoints, unreadyMod int, hosts []string, portNums []int32) *discovery.EndpointSlice {
	return generateEndpointSliceWithOffset(serviceName, namespace, sliceNum, sliceNum, numEndpoints, unreadyMod, hosts, portNums)
}
