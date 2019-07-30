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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	apivalidation "k8s.io/kubernetes/pkg/apis/core/validation"
	"k8s.io/kubernetes/pkg/apis/discovery"
)

func ValidateEndpointSlice(endpointSlice *discovery.EndpointSlice) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&endpointSlice.ObjectMeta, true, apivalidation.ValidateReplicationControllerName, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateEndpointSliceEndpoints(&endpointSlice.Endpoints, field.NewPath("endpoints"))...)
	allErrs = append(allErrs, ValidateEndpointSlicePorts(&endpointSlice.Ports, field.NewPath("ports"))...)
	return allErrs
}

func ValidateEndpointSliceEndpoints(endpoints *[]discovery.Endpoint, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, endpoint := range *endpoints {
		idxPath := fldPath.Index(i)

		if len(endpoint.Targets) < 1 {
			allErrs = append(allErrs, field.Required(idxPath.Child("targets"), "At least 1 target required"))
		}
	}

	return allErrs
}

func ValidateEndpointSlicePorts(endpointPorts *[]discovery.EndpointPort, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	supportedPortProtocols := sets.NewString(string(v1.ProtocolTCP), string(v1.ProtocolUDP), string(v1.ProtocolSCTP))
	portNames := sets.String{}

	for i, endpointPort := range *endpointPorts {
		idxPath := fldPath.Index(i)

		if len(*endpointPort.Name) > 0 {
			for _, msg := range validation.IsValidPortName(*endpointPort.Name) {
				allErrs = append(allErrs, field.Invalid(idxPath.Child("name"), *endpointPort.Name, msg))
			}
		}

		if portNames.Has(*endpointPort.Name) {
			allErrs = append(allErrs, field.Duplicate(idxPath.Child("name"), *endpointPort.Name))
		} else {
			portNames.Insert(*endpointPort.Name)
		}

		if len(*endpointPort.Protocol) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("protocol"), ""))
		} else if !supportedPortProtocols.Has(string(*endpointPort.Protocol)) {
			allErrs = append(allErrs, field.NotSupported(idxPath.Child("protocol"), *endpointPort.Protocol, supportedPortProtocols.List()))
		}
	}

	return allErrs
}
