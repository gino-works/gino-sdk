package v1

import apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// EndpointGroup is an object that holds information about a group of gino endpoints.
// Short name gepg -> "gino endpoint group".
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=endpointgroups,scope=Cluster,shortName=gepg
type EndpointGroup struct {
	apimachineryapismetav1.TypeMeta   `json:",inline"`
	apimachineryapismetav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EndpointGroupSpec   `json:"spec"`
	Status EndpointGroupStatus `json:"status,omitempty"`
}

// EndpointGroupList is the list type for EndpointGroups.
// +kubebuilder:object:root=true
type EndpointGroupList struct {
	apimachineryapismetav1.TypeMeta `json:",inline"`
	apimachineryapismetav1.ListMeta `json:"metadata,omitempty"`

	Items []EndpointGroup `json:"items"`
}

// EndpointGroupSpec is the spec for an EndpointGroup resource.
type EndpointGroupSpec struct {
	EndpointTemplate `json:",inline"`

	// Selector picks the Endpoints this group's template is applied to. Empty selects nothing.
	// +optional
	Selector TargetSelector `json:"selector,omitempty"`

	// Weight is used to decide order of resolution for group data onto endpoints. Higher wins.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=255
	// +kubebuilder:default=128
	// +optional
	Weight int32 `json:"weight"`
}

// EndpointGroupStatus is the status for an EndpointGroup resource.
type EndpointGroupStatus struct{}
