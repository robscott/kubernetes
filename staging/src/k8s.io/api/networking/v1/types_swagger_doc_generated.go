/*
Copyright The Kubernetes Authors.

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

package v1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-generated-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE. DO NOT EDIT.
var map_HTTPIngressPath = map[string]string{
	"":         "HTTPIngressPath associates a path with a backend. Incoming urls matching the path are forwarded to the backend.",
	"path":     "Path is matched against the path of an incoming request. Currently it can contain characters disallowed from the conventional \"path\" part of a URL as defined by RFC 3986. Paths must begin with a '/'. When unspecified, all paths from incoming requests are matched.",
	"pathType": "PathType determines the interpretation of the Path matching. PathType can be one of Exact, Prefix, or ImplementationSpecific. Implementations are required to support all path types. Defaults to ImplementationSpecific.",
	"backend":  "Backend defines the referenced service endpoint to which the traffic will be forwarded to.",
}

func (HTTPIngressPath) SwaggerDoc() map[string]string {
	return map_HTTPIngressPath
}

var map_HTTPIngressRuleValue = map[string]string{
	"":      "HTTPIngressRuleValue is a list of http selectors pointing to backends. In the example: http://<host>/<path>?<searchpart> -> backend where where parts of the url correspond to RFC 3986, this resource will be used to match against everything after the last '/' and before the first '?' or '#'.",
	"paths": "A collection of paths that map requests to backends.",
}

func (HTTPIngressRuleValue) SwaggerDoc() map[string]string {
	return map_HTTPIngressRuleValue
}

var map_IPBlock = map[string]string{
	"":       "IPBlock describes a particular CIDR (Ex. \"192.168.1.1/24\",\"2001:db9::/64\") that is allowed to the pods matched by a NetworkPolicySpec's podSelector. The except entry describes CIDRs that should not be included within this rule.",
	"cidr":   "CIDR is a string representing the IP Block Valid examples are \"192.168.1.1/24\" or \"2001:db9::/64\"",
	"except": "Except is a slice of CIDRs that should not be included within an IP Block Valid examples are \"192.168.1.1/24\" or \"2001:db9::/64\" Except values will be rejected if they are outside the CIDR range",
}

func (IPBlock) SwaggerDoc() map[string]string {
	return map_IPBlock
}

var map_Ingress = map[string]string{
	"":         "Ingress is a collection of rules that allow inbound connections to reach the endpoints defined by a backend. An Ingress can be configured to give services externally-reachable urls, load balance traffic, terminate SSL, offer name based virtual hosting etc.",
	"metadata": "Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata",
	"spec":     "Spec is the desired state of the Ingress. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status",
	"status":   "Status is the current state of the Ingress. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status",
}

func (Ingress) SwaggerDoc() map[string]string {
	return map_Ingress
}

var map_IngressBackend = map[string]string{
	"":            "IngressBackend describes all endpoints for a given service and port.",
	"serviceName": "Specifies the name of the referenced service.",
	"servicePort": "Specifies the port of the referenced service.",
}

func (IngressBackend) SwaggerDoc() map[string]string {
	return map_IngressBackend
}

var map_IngressClass = map[string]string{
	"":         "IngressClass represents the class of the Ingress, referenced by the Ingress Spec. The `ingressclass.kubernetes.io/is-default-class` annotation can be used to indicate that an IngressClass should be considered default. When a single IngressClass resource has this annotation set to true, new Ingress resources without a class specified will be assigned this default class.",
	"metadata": "Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"spec":     "Spec is the desired state of the IngressClass. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status",
}

func (IngressClass) SwaggerDoc() map[string]string {
	return map_IngressClass
}

var map_IngressClassList = map[string]string{
	"":         "IngressClassList is a collection of IngressClasses.",
	"metadata": "Standard list metadata.",
	"items":    "Items is the list of IngressClasses.",
}

func (IngressClassList) SwaggerDoc() map[string]string {
	return map_IngressClassList
}

var map_IngressClassSpec = map[string]string{
	"":           "IngressClassSpec provides information about the class of an Ingress.",
	"controller": "Controller refers to the name of the controller that should handle this class. This allows for different \"flavors\" that are controlled by the same controller. For example, you may have different Parameters for the same implementing controller. This should be specified as a domain-prefixed path no more than 250 characters in length, e.g. \"acme.io/ingress-controller\". This field is immutable.",
	"parameters": "Parameters is a link to a resource containing additional configuration for the controller. This is optional if the controller does not require extra parameters. Example configuration resources include `core.ConfigMap` or a controller specific Custom Resource.",
}

func (IngressClassSpec) SwaggerDoc() map[string]string {
	return map_IngressClassSpec
}

var map_IngressList = map[string]string{
	"":         "IngressList is a collection of Ingress.",
	"metadata": "Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata",
	"items":    "Items is the list of Ingress.",
}

func (IngressList) SwaggerDoc() map[string]string {
	return map_IngressList
}

var map_IngressRule = map[string]string{
	"":     "IngressRule represents the rules mapping the paths under a specified host to the related backend services. Incoming requests are first evaluated for a host match, then routed to the backend associated with the matching IngressRuleValue.",
	"host": "Host is the fully qualified domain name of a network host, as defined by RFC 3986. Note the following deviations from the \"host\" part of the URI as defined in the RFC: 1. IPs are not allowed. Currently an IngressRuleValue can only apply to\n   the IP in the Spec of the parent Ingress.\n2. The `:` delimiter is not respected because ports are not allowed.\n\t  Currently the port of an Ingress is implicitly :80 for http and\n\t  :443 for https.\nBoth these may change in the future. Incoming requests are matched against the host before the IngressRuleValue. If the host is unspecified, the Ingress routes all traffic based on the specified IngressRuleValue.",
}

func (IngressRule) SwaggerDoc() map[string]string {
	return map_IngressRule
}

var map_IngressRuleValue = map[string]string{
	"": "IngressRuleValue represents a rule to apply against incoming requests. If the rule is satisfied, the request is routed to the specified backend. Currently mixing different types of rules in a single Ingress is disallowed, so exactly one of the following must be set.",
}

func (IngressRuleValue) SwaggerDoc() map[string]string {
	return map_IngressRuleValue
}

var map_IngressSpec = map[string]string{
	"":               "IngressSpec describes the Ingress the user wishes to exist.",
	"class":          "Class is the name of the IngressClass cluster resource. The associated IngressClass defines which controller will implement the resource. This replaces the deprecated `kubernetes.io/ingress.class` annotation. For backwards compatibility, when that annotation is set, it must be given precedence over this field. The controller may emit a warning if the field and annotation have different values. Implementations of this API should ignore Ingresses without a class specified. An IngressClass resource may be marked as default, which can be used to set a default value for this field. For more information, refer to the IngressClass documentation.",
	"defaultBackend": "The default backend capable of servicing requests that don't match any rule. If rules is empty, DefaultBackend must be specified. If DefaultBackend is not set, the handling of requests that do not match any of the rules will be up to the Ingress controller.",
	"tls":            "TLS configuration. Currently the Ingress only supports a single TLS port, 443. If multiple members of this list specify different hosts, they will be multiplexed on the same port according to the hostname specified through the SNI TLS extension, if the ingress controller fulfilling the ingress supports SNI.",
	"rules":          "A list of host rules used to configure the Ingress. If unspecified, or no rule matches, all traffic is sent to the default backend.",
}

func (IngressSpec) SwaggerDoc() map[string]string {
	return map_IngressSpec
}

var map_IngressStatus = map[string]string{
	"":             "IngressStatus describe the current state of the Ingress.",
	"loadBalancer": "LoadBalancer contains the current status of the load-balancer.",
}

func (IngressStatus) SwaggerDoc() map[string]string {
	return map_IngressStatus
}

var map_IngressTLS = map[string]string{
	"":           "IngressTLS describes the transport layer security associated with an Ingress.",
	"hosts":      "Hosts are a list of hosts included in the TLS certificate. The values in this list must match the name/s used in the tlsSecret. Defaults to the wildcard host setting for the loadbalancer controller fulfilling this Ingress, if left unspecified.",
	"secretName": "SecretName is the name of the secret used to terminate SSL traffic on 443. Field is left optional to allow SSL routing based on SNI hostname alone. If the SNI host in a listener conflicts with the \"Host\" header field used by an IngressRule, the SNI host is used for termination and value of the Host header is used for routing.",
}

func (IngressTLS) SwaggerDoc() map[string]string {
	return map_IngressTLS
}

var map_NetworkPolicy = map[string]string{
	"":         "NetworkPolicy describes what network traffic is allowed for a set of Pods",
	"metadata": "Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"spec":     "Specification of the desired behavior for this NetworkPolicy.",
}

func (NetworkPolicy) SwaggerDoc() map[string]string {
	return map_NetworkPolicy
}

var map_NetworkPolicyEgressRule = map[string]string{
	"":      "NetworkPolicyEgressRule describes a particular set of traffic that is allowed out of pods matched by a NetworkPolicySpec's podSelector. The traffic must match both ports and to. This type is beta-level in 1.8",
	"ports": "List of destination ports for outgoing traffic. Each item in this list is combined using a logical OR. If this field is empty or missing, this rule matches all ports (traffic not restricted by port). If this field is present and contains at least one item, then this rule allows traffic only if the traffic matches at least one port in the list.",
	"to":    "List of destinations for outgoing traffic of pods selected for this rule. Items in this list are combined using a logical OR operation. If this field is empty or missing, this rule matches all destinations (traffic not restricted by destination). If this field is present and contains at least one item, this rule allows traffic only if the traffic matches at least one item in the to list.",
}

func (NetworkPolicyEgressRule) SwaggerDoc() map[string]string {
	return map_NetworkPolicyEgressRule
}

var map_NetworkPolicyIngressRule = map[string]string{
	"":      "NetworkPolicyIngressRule describes a particular set of traffic that is allowed to the pods matched by a NetworkPolicySpec's podSelector. The traffic must match both ports and from.",
	"ports": "List of ports which should be made accessible on the pods selected for this rule. Each item in this list is combined using a logical OR. If this field is empty or missing, this rule matches all ports (traffic not restricted by port). If this field is present and contains at least one item, then this rule allows traffic only if the traffic matches at least one port in the list.",
	"from":  "List of sources which should be able to access the pods selected for this rule. Items in this list are combined using a logical OR operation. If this field is empty or missing, this rule matches all sources (traffic not restricted by source). If this field is present and contains at least one item, this rule allows traffic only if the traffic matches at least one item in the from list.",
}

func (NetworkPolicyIngressRule) SwaggerDoc() map[string]string {
	return map_NetworkPolicyIngressRule
}

var map_NetworkPolicyList = map[string]string{
	"":         "NetworkPolicyList is a list of NetworkPolicy objects.",
	"metadata": "Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"items":    "Items is a list of schema objects.",
}

func (NetworkPolicyList) SwaggerDoc() map[string]string {
	return map_NetworkPolicyList
}

var map_NetworkPolicyPeer = map[string]string{
	"":                  "NetworkPolicyPeer describes a peer to allow traffic from. Only certain combinations of fields are allowed",
	"podSelector":       "This is a label selector which selects Pods. This field follows standard label selector semantics; if present but empty, it selects all pods.\n\nIf NamespaceSelector is also set, then the NetworkPolicyPeer as a whole selects the Pods matching PodSelector in the Namespaces selected by NamespaceSelector. Otherwise it selects the Pods matching PodSelector in the policy's own Namespace.",
	"namespaceSelector": "Selects Namespaces using cluster-scoped labels. This field follows standard label selector semantics; if present but empty, it selects all namespaces.\n\nIf PodSelector is also set, then the NetworkPolicyPeer as a whole selects the Pods matching PodSelector in the Namespaces selected by NamespaceSelector. Otherwise it selects all Pods in the Namespaces selected by NamespaceSelector.",
	"ipBlock":           "IPBlock defines policy on a particular IPBlock. If this field is set then neither of the other fields can be.",
}

func (NetworkPolicyPeer) SwaggerDoc() map[string]string {
	return map_NetworkPolicyPeer
}

var map_NetworkPolicyPort = map[string]string{
	"":         "NetworkPolicyPort describes a port to allow traffic on",
	"protocol": "The protocol (TCP, UDP, or SCTP) which traffic must match. If not specified, this field defaults to TCP.",
	"port":     "The port on the given protocol. This can either be a numerical or named port on a pod. If this field is not provided, this matches all port names and numbers.",
}

func (NetworkPolicyPort) SwaggerDoc() map[string]string {
	return map_NetworkPolicyPort
}

var map_NetworkPolicySpec = map[string]string{
	"":            "NetworkPolicySpec provides the specification of a NetworkPolicy",
	"podSelector": "Selects the pods to which this NetworkPolicy object applies. The array of ingress rules is applied to any pods selected by this field. Multiple network policies can select the same set of pods. In this case, the ingress rules for each are combined additively. This field is NOT optional and follows standard label selector semantics. An empty podSelector matches all pods in this namespace.",
	"ingress":     "List of ingress rules to be applied to the selected pods. Traffic is allowed to a pod if there are no NetworkPolicies selecting the pod (and cluster policy otherwise allows the traffic), OR if the traffic source is the pod's local node, OR if the traffic matches at least one ingress rule across all of the NetworkPolicy objects whose podSelector matches the pod. If this field is empty then this NetworkPolicy does not allow any traffic (and serves solely to ensure that the pods it selects are isolated by default)",
	"egress":      "List of egress rules to be applied to the selected pods. Outgoing traffic is allowed if there are no NetworkPolicies selecting the pod (and cluster policy otherwise allows the traffic), OR if the traffic matches at least one egress rule across all of the NetworkPolicy objects whose podSelector matches the pod. If this field is empty then this NetworkPolicy limits all outgoing traffic (and serves solely to ensure that the pods it selects are isolated by default). This field is beta-level in 1.8",
	"policyTypes": "List of rule types that the NetworkPolicy relates to. Valid options are \"Ingress\", \"Egress\", or \"Ingress,Egress\". If this field is not specified, it will default based on the existence of Ingress or Egress rules; policies that contain an Egress section are assumed to affect Egress, and all policies (whether or not they contain an Ingress section) are assumed to affect Ingress. If you want to write an egress-only policy, you must explicitly specify policyTypes [ \"Egress\" ]. Likewise, if you want to write a policy that specifies that no egress is allowed, you must specify a policyTypes value that include \"Egress\" (since such a policy would not include an Egress section and would otherwise default to just [ \"Ingress\" ]). This field is beta-level in 1.8",
}

func (NetworkPolicySpec) SwaggerDoc() map[string]string {
	return map_NetworkPolicySpec
}

// AUTO-GENERATED FUNCTIONS END HERE
