package api

import (
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

// GroupName is the group name use in this package
const GroupName = "nuage.kubernetes.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = unversioned.GroupVersion{Group: GroupName, Version: "v1"}

const NuageNetworkPolicyResourcePlural = "nuagenetworkpolicies"

type NuageNetworkPolicy struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata"`
	Spec                 NuageNetworkPolicySpec `json:"spec"`
}

type NuageNetworkPolicySpec struct {
	PodSelector unversioned.LabelSelector       `json:"podSelector"`
	Ingress     []NuageNetworkPolicyIngressRule `json:"ingress"`
}

type NuageNetworkPolicyIngressRule struct {
	Ports []NuageNetworkPolicyPort `json:"ports,omitempty"`
	From  []NuageNetworkPolicyPeer `json:"from,omitempty"`
}

type NuageNetworkPolicyPort struct {
	Protocol *api.Protocol `json:"protocol,omitempty"`

	Port intstr.IntOrString `json:"port,omitempty"`
}

type NuageNetworkPolicyPeer struct {
	PodSelector   *unversioned.LabelSelector `json:"podSelector,omitempty"`
	FieldSelector *fields.Selector           `json:"fieldSelector,omitempty"`
}

type NuageNetworkPolicyList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []NuageNetworkPolicy `json:"items"`
}
