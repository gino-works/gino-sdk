package v1

import apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Location is an object that holds information about a gino "location" -- that is, some place that
// gino endpoints live at (or really not, its obviously a logical grouping!).
// Short name gl -> "gino location".
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=locations,scope=Cluster,shortName=gl
type Location struct {
	apimachineryapismetav1.TypeMeta   `json:",inline"`
	apimachineryapismetav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LocationSpec   `json:"spec"`
	Status LocationStatus `json:"status,omitempty"`
}

// LocationList is the list type for Locations.
// +kubebuilder:object:root=true
type LocationList struct {
	apimachineryapismetav1.TypeMeta `json:",inline"`
	apimachineryapismetav1.ListMeta `json:"metadata,omitempty"`

	Items []Location `json:"items"`
}

// LocationSpec is the spec for a Location resource.
type LocationSpec struct {
	LocationTemplate `json:",inline"`

	// Parent is the optional parent of this location -- for things like site/location in netbox.
	// +optional
	Parent *string `json:"parent,omitempty"`
}

// LocationStatus is the status for a Location resource.
type LocationStatus struct {
	// ObservedGeneration is the metadata.generation this status was computed from. The apiserver
	// bumps metadata.generation on every spec change (not on status or metadata changes), so a
	// reader can tell whether status reflects the current spec: if observedGeneration <
	// metadata.generation the controller has not caught up yet and the rest of status is stale.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions summarize the locations's state. Known types:
	//   - Resolved: template things have been applied; False with reason NoSecret etc. when the
	//     effective configuration is unusable.
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []apimachineryapismetav1.Condition `json:"conditions,omitempty" patchMergeKey:"type" patchStrategy:"merge"`

	// Resolved holds the fully resolved location data -- meaning we have fully processed any
	// parent locations and resolved our data compeletely.
	Resolved *LocationResolved `json:"resolved,omitempty"`
}
