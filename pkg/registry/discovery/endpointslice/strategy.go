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

package endpointslice

import (
	"context"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/discovery"
	"k8s.io/kubernetes/pkg/apis/discovery/validation"
)

// endpointSliceStrategy implements verification logic for Replication.
type endpointSliceStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Replication EndpointSlice objects.
var Strategy = endpointSliceStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// NamespaceScoped returns true because all EndpointSlices need to be within a namespace.
func (endpointSliceStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears the status of an EndpointSlice before creation.
func (endpointSliceStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	endpointSlice := obj.(*discovery.EndpointSlice)
	endpointSlice.Generation = 1
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (endpointSliceStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newEndpointSlice := obj.(*discovery.EndpointSlice)
	oldEndpointSlice := old.(*discovery.EndpointSlice)

	// Any changes to the endpoints increment the generation number, any changes to the
	// status should reflect the generation number of the corresponding object.
	// See metav1.ObjectMeta description for more information on Generation.
	if !apiequality.Semantic.DeepEqual(oldEndpointSlice.Endpoints, newEndpointSlice.Endpoints) {
		newEndpointSlice.Generation = oldEndpointSlice.Generation + 1
	}

}

// Validate validates a new EndpointSlice.
func (endpointSliceStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	endpointSlice := obj.(*discovery.EndpointSlice)
	err := validation.ValidateEndpointSlice(endpointSlice)
	return err
}

// Canonicalize normalizes the object after validation.
func (endpointSliceStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for EndpointSlice; this means POST is needed to create one.
func (endpointSliceStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (endpointSliceStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateEndpointSlice(obj.(*discovery.EndpointSlice))
}

// AllowUnconditionalUpdate is the default update policy for EndpointSlice objects.
func (endpointSliceStrategy) AllowUnconditionalUpdate() bool {
	return true
}

type endpointSliceStatusStrategy struct {
	endpointSliceStrategy
}
