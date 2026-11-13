package v1

import (
	apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResolvedSourceKind identifies what kind of object (or pseudo-source) a resolved field came from.
// +kubebuilder:validation:Enum=Endpoint;EndpointGroup;Location;Labels;Default
type ResolvedSourceKind string

// enumerations of resolved sources.
const (
	ResolvedSourceKindEndpoint      ResolvedSourceKind = "Endpoint"
	ResolvedSourceKindEndpointGroup ResolvedSourceKind = "EndpointGroup"
	ResolvedSourceKindLocation      ResolvedSourceKind = "Location"
	ResolvedSourceKindLabels        ResolvedSourceKind = "Labels"
	ResolvedSourceKindDefault       ResolvedSourceKind = "Default"
)

// ResolvedSource identifies the object a resolved field was taken from.
type ResolvedSource struct {
	Kind ResolvedSourceKind `json:"kind"`
	// Name of the source object. Empty for Labels/Default.
	// +optional
	Name string `json:"name,omitempty"`
	// Namespace of the source object, for namespaced kinds.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// TargetSelector picks Endpoints by their labels and/or by the labels of their Location. These are
// AND'd together and neither can be nil -- an empty set means "all" -- so you could do `{}` for
// locations, then simply select based on endpoint labels to do a fleet wide type match. Location
// ancestry is ignored for this -- just for simplicity, could be added in the future, but locations
// should just be created w/ their appropriate labels in the first place.
// +kubebuilder:validation:XValidation:rule="has(self.locations) && has(self.endpoints)",message="both locations and endpoints selectors must be set; use {} to select all"
type TargetSelector struct {
	// Locations selects by Location labels.
	Locations *apimachineryapismetav1.LabelSelector `json:"locations"`
	// Endpoints selects by Endpoint labels.
	Endpoints *apimachineryapismetav1.LabelSelector `json:"endpoints"`
}
