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

package validation

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/discovery"
)

func TestValidateEndpointSlice(t *testing.T) {
	standardMeta := metav1.ObjectMeta{
		Name:      "hello",
		Namespace: "world",
	}

	nameHTTP := "http"
	nameEmpty := ""
	nameCaps := "ABCDEF"
	nameChars := "almost_valid"
	protoTCP := api.ProtocolTCP
	targetTypeIP := discovery.IPTargetType
	targetTypeOther := discovery.TargetType("Other")

	testCases := map[string]struct {
		isExpectedFailure bool
		endpointSlice     *discovery.EndpointSlice
	}{
		"good-slice": {
			isExpectedFailure: false,
			endpointSlice: &discovery.EndpointSlice{
				ObjectMeta: standardMeta,
				TargetType: &targetTypeIP,
				Ports: []discovery.EndpointPort{{
					Name:     &nameHTTP,
					Protocol: &protoTCP,
				}},
				Endpoints: []discovery.Endpoint{{
					Targets: []string{"1.2.3.4"},
				}},
			},
		},
		"empty-port-name": {
			isExpectedFailure: false,
			endpointSlice: &discovery.EndpointSlice{
				ObjectMeta: standardMeta,
				TargetType: &targetTypeIP,
				Ports: []discovery.EndpointPort{{
					Name:     &nameEmpty,
					Protocol: &protoTCP,
				}, {
					Name:     &nameHTTP,
					Protocol: &protoTCP,
				}},
				Endpoints: []discovery.Endpoint{{
					Targets: []string{"1.2.3.4"},
				}},
			},
		},
		"empty-ports-and-endpoints": {
			isExpectedFailure: false,
			endpointSlice: &discovery.EndpointSlice{
				ObjectMeta: standardMeta,
				TargetType: &targetTypeIP,
				Ports:      []discovery.EndpointPort{},
				Endpoints:  []discovery.Endpoint{},
			},
		},

		// expected failures
		"duplicate-port-name": {
			isExpectedFailure: true,
			endpointSlice: &discovery.EndpointSlice{
				ObjectMeta: standardMeta,
				TargetType: &targetTypeIP,
				Ports: []discovery.EndpointPort{{
					Name:     &nameEmpty,
					Protocol: &protoTCP,
				}, {
					Name:     &nameEmpty,
					Protocol: &protoTCP,
				}},
				Endpoints: []discovery.Endpoint{},
			},
		},
		"bad-port-name-caps": {
			isExpectedFailure: true,
			endpointSlice: &discovery.EndpointSlice{
				ObjectMeta: standardMeta,
				TargetType: &targetTypeIP,
				Ports: []discovery.EndpointPort{{
					Name:     &nameCaps,
					Protocol: &protoTCP,
				}},
				Endpoints: []discovery.Endpoint{},
			},
		},
		"bad-port-name-chars": {
			isExpectedFailure: true,
			endpointSlice: &discovery.EndpointSlice{
				ObjectMeta: standardMeta,
				TargetType: &targetTypeIP,
				Ports: []discovery.EndpointPort{{
					Name:     &nameChars,
					Protocol: &protoTCP,
				}},
				Endpoints: []discovery.Endpoint{},
			},
		},
		"no-endpoint-targets": {
			isExpectedFailure: true,
			endpointSlice: &discovery.EndpointSlice{
				ObjectMeta: standardMeta,
				TargetType: &targetTypeIP,
				Ports: []discovery.EndpointPort{{
					Name:     &nameHTTP,
					Protocol: &protoTCP,
				}},
				Endpoints: []discovery.Endpoint{{
					Targets: []string{},
				}},
			},
		},
		"bad-target-type": {
			isExpectedFailure: true,
			endpointSlice: &discovery.EndpointSlice{
				ObjectMeta: standardMeta,
				TargetType: &targetTypeOther,
				Ports: []discovery.EndpointPort{{
					Name:     &nameHTTP,
					Protocol: &protoTCP,
				}},
				Endpoints: []discovery.Endpoint{{
					Targets: []string{},
				}},
			},
		},
		"empty-everything": {
			isExpectedFailure: true,
			endpointSlice:     &discovery.EndpointSlice{},
		},
	}

	for name, testCase := range testCases {
		errs := ValidateEndpointSlice(testCase.endpointSlice)
		if len(errs) == 0 && testCase.isExpectedFailure {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !testCase.isExpectedFailure {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}
	}
}
