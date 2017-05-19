package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api"
)

const NuageNetworkPolicyResourcePlural = "nuagenetworkpolicies"

type NuageNetworkPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              NuageNetworkPolicySpec `json:"spec"`
}

type NuageNetworkPolicySpec struct {
	PodSelector metav1.LabelSelector            `json:"podSelector"`
	Ingress     []NuageNetworkPolicyIngressRule `json:"ingress"`
}

type NuageNetworkPolicyIngressRule struct {
	Ports []NuageNetworkPolicyPort `json:"ports,omitempty"`
	From  []NuageNetworkPolicyPeer `json:"from,omitempty"`
}

type NuageNetworkPolicyPort struct {
	Protocol *api.Protocol      `json:"protocol,omitempty"`
	Port     intstr.IntOrString `json:"port,omitempty"`
}

type NuageNetworkPolicyPeer struct {
	PodSelector   *metav1.LabelSelector `json:"podSelector,omitempty"`
	FieldSelector *metav1.LabelSelector `json:"fieldSelector,omitempty"`
}

type NuageNetworkPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NuageNetworkPolicy `json:"items"`
}
