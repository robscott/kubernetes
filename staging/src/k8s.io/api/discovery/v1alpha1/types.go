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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EndpointSlice represents a set of endpoints that implement a service
type EndpointSlice struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// This field specifies the list of ports associated with each endpoint in the EndpointSlice
	// Each EndpointPort must have a unique port name.
	// +optional
	// +listType=atomic
	Ports []EndpointPort `json:"ports" protobuf:"bytes,2,rep,name=ports"`
	// This field specifies the type of targets associated with each endpoint in the EndpointSlice
	// Default is IP.
	// +optional
	TargetType *TargetType `json:"targetType" protobuf:"bytes,3,rep,name=targetType"`
	// The set of endpoints backing this slice
	// +listType=atomic
	Endpoints []Endpoint `json:"endpoints" protobuf:"bytes,4,rep,name=endpoints"`
}

// Endpoint represents a set of backends implementing a service
type Endpoint struct {
	// Required: Targets of the endpoint. Must contain at least one target.
	// Different consumers (e.g. kube-proxy) handle different types of
	// targets in the context of its own capabilities.
	// +listType=set
	Targets []string `json:"targets" protobuf:"bytes,1,rep,name=targets"`
	// Required: the conditions of the endpoint.
	Conditions EndpointConditions `json:"conditions,omitempty" protobuf:"bytes,2,opt,name=conditions"`
	// The stated hostname of this endpoint. This field may be used by consumers
	// of endpoint to distinguish endpoints from each other (e.g. in DNS names).
	// Multiple endpoints which use the same hostname should be considered fungible (e.g. multiple A values in DNS).
	// +optional
	Hostname *string `json:"hostname,omitempty" protobuf:"bytes,3,opt,name=hostname"`
	// Reference to object providing the endpoint.
	// +optional
	TargetRef *v1.ObjectReference `json:"targetRef,omitempty" protobuf:"bytes,4,opt,name=targetRef"`
	// Topology can contain arbitrary topology information associated with the
	// endpoint.
	// Key/value pairs contained in topology must conform with the label format.
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels
	// Topology may include, but is not limited to the following well known keys:
	// kubernetes.io/hostname: the value indicates the hostname of the node
	// where the endpoint is located. This should match the corresponding
	// node label.
	// topology.kubernetes.io/zone: the value indicates the zone where the
	// endpoint is located. This should match the corresponding node label.
	// topology.kubernetes.io/region: the value indicates the region where the
	// endpoint is located. This should match the corresponding node label.
	// +optional
	Topology map[string]string `json:"topology,omitempty" protobuf:"bytes,5,opt,name=topology"`
}

// EndpointConditions represents the current condition of an endpoint
type EndpointConditions struct {
	// Ready indicates if the endpoint is ready to serve traffic
	Ready bool `json:"ready,omitempty" protobuf:"bytes,1,name=ready"`
}

// EndpointPort represents a Port used by an EndpointSlice
type EndpointPort struct {
	// The name of this port.
	// If the EndpointSlice is dervied from K8s service, this corresponds to ServicePort.Name.
	// Name must be a IANA_SVC_NAME or an empty string.
	// Default is empty string.
	// +optional
	Name *string `json:"name,omitempty" protobuf:"bytes,1,name=name"`
	// The IP protocol for this port.
	// Must be UDP, TCP, or SCTP.
	// Default is TCP.
	// +optional
	Protocol *v1.Protocol `json:"protocol,omitempty" protobuf:"bytes,2,name=protocol"`
	// The port number of the endpoint.
	// If this is not specified, ports are not restricted and must be
	// interpreted in the context of the specific consumer.
	// +optional
	Port *int32 `json:"port,omitempty" protobuf:"bytes,3,opt,name=port"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EndpointSliceList represents a list of endpoint slices
type EndpointSliceList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of endpoint slices
	// +listType=set
	Items []EndpointSlice `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// TargetType represents the type of Target referred to by an Endpoint
type TargetType string

const (
	// IPTargetType represents an IP Target
	IPTargetType = TargetType("IP")
)
