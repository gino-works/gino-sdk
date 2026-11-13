package v1

import (
	apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IndexEndpointLocation is the index field for looking up endpoints for a given location quickly.
const IndexEndpointLocation = ".spec.location"

// Endpoint is an object that holds information about a gino endpoint -- that is, a router, switch,
// firewall, or *something* that gino is interested in.
// Short name gep -> "gino endpoint".
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=endpoints,shortName=gep
// +kubebuilder:printcolumn:name="Resolved",type=string,JSONPath=".status.conditions[?(@.type=='Resolved')].status"
type Endpoint struct {
	apimachineryapismetav1.TypeMeta   `json:",inline"`
	apimachineryapismetav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EndpointSpec   `json:"spec"`
	Status EndpointStatus `json:"status,omitempty"`
}

// EndpointList is the list type for Endpoints.
// +kubebuilder:object:root=true
type EndpointList struct {
	apimachineryapismetav1.TypeMeta `json:",inline"`
	apimachineryapismetav1.ListMeta `json:"metadata,omitempty"`

	Items []Endpoint `json:"items"`
}

// EndpointSpec is the spec for an Endpoint resource.
type EndpointSpec struct {
	EndpointTemplate `json:",inline"`

	// Location is the location of this endpoint.
	Location string `json:"location"`
	// Host is the name/v4/v6 hostname to use to connect to the device.
	Host string `json:"host"`
}

// EndpointStatus is the status for an Endpoint resource.
type EndpointStatus struct {
	// ObservedGeneration is the metadata.generation this status was computed from. The apiserver
	// bumps metadata.generation on every spec change (not on status or metadata changes), so a
	// reader can tell whether status reflects the current spec: if observedGeneration <
	// metadata.generation the controller has not caught up yet and the rest of status is stale.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions summarize the endpoint's state. Known types:
	//   - Resolved: groups/location have been applied; False with reason NoVendor etc. when the
	//     effective configuration is unusable.
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []apimachineryapismetav1.Condition `json:"conditions,omitempty" patchMergeKey:"type" patchStrategy:"merge"`

	// Resolved holds the fully resolved endpoint data -- meaning any related groups have that
	// refer to this endpoint have been fully processed and we are left w/ the full set of
	// endpoint data *with* info about which fields came from where.
	Resolved *EndpointResolved `json:"resolved,omitempty"`
}
