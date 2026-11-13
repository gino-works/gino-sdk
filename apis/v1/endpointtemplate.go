package v1

// EndpointTemplate is the inheritable subset of an endpoint's settings. It is the single
// definition of "what can come from a group or location": used by EndpointGroup//Location templates
// and mirrored into EndpointResolved.
type EndpointTemplate struct {
	// Vendor is the name of the vendor, this drives association w/ Platform cr for informing
	// workers how to talk to this Endpoint.
	// +optional
	Vendor string `json:"vendor,omitempty"`
	// Platform is the name of the platform, same idea as Vendor. This is optional.
	// +optional
	Platform string `json:"platform,omitempty"`
	// Profile is the profile assigned to the endpoint, this is the final widget for informing the
	// workers how to interact with this Endpoint.
	// +optional
	Profile string `json:"profile,omitempty"`
	// Tags hold additional k:v data for the endpoint. These are not used for selection (thats what
	// the regular labels are for) but instead for expressing additional context about the endpoint.
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
	// InterfaceTags holds a list of InterfaceTag used to associate interfaces w/ tag data. Tags
	// are appended to in resolution, so last match wins when actually using these downstream.
	// +listType=atomic
	// +optional
	InterfaceTags []InterfaceTag `json:"interfaceTags,omitempty"`
	// DisableCollection ensures no workers are spawned for this endpoint.
	// +optional
	DisableCollection *bool `json:"disableCollection,omitempty"`
	// Cluster is the name of the cluster that collection is done for on this endpoint.
	// +optional
	Cluster string `json:"cluster,omitempty"`
	// Credentials names, per collector kind, the Secret in the endpoint's own namespace holding
	// the credentials workers of that kind use. Workers run in the endpoint's namespace and
	// Locations replicate secrets into it, so only a name is needed. Typically set once on a group
	// and inherited. Its the workers responsibility to know which key it wants (i.e. "token").
	// +optional
	Credentials map[CollectorKind]string `json:"credentials,omitempty"`
}

// InterfaceTag holds interface matching -> tag data.
type InterfaceTag struct {
	Match InterfaceTagMatch `json:"match"`
	Tags  map[string]string `json:"tags"`
}

// InterfaceTagMatch selects interfaces; every set field must match (AND). At least one is required.
// +kubebuilder:validation:XValidation:rule="has(self.name) || has(self.description)",message="at least one of name or description must be set"
type InterfaceTagMatch struct {
	// Name is a Go regex matched against the interface name.
	// +kubebuilder:validation:MaxLength=256
	// +optional
	Name *string `json:"name,omitempty"`
	// Description is a Go regex matched against the interface description.
	// +kubebuilder:validation:MaxLength=256
	// +optional
	Description *string `json:"description,omitempty"`
}
