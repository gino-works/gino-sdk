package v1

// LocationTemplate is the inheritable subset of a location's settings.
type LocationTemplate struct {
	// Namespace defines if/how a namespace corresponding to this object will be managed.
	// +optional
	Namespace LocationNamespaceSpec `json:"namespace,omitempty"`
	// Cluster is the name of the cluster that endpoints for this site are collected from. This
	// is the default home of this information but individual Endpoints can override it.
	// Unset in spec means we will try to resolve it from our parent(s), but unset will mean no
	// collections can be deployed for this location (as in, if it is not resolvable).
	// +optional
	Cluster string `json:"cluster,omitempty"`
	// DisableCollection ensures no workers are spawned for this endpoint.
	// +optional
	DisableCollection *bool `json:"disableCollection,omitempty"`
	// TYpe is the type of the location.
	// +optional
	Type string `json:"type,omitempty"`
	// Tags hold additional k:v data for the endpoint. These are not used for selection (thats what
	// the regular labels are for) but instead for expressing additional context about the endpoint.
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

// LocationNamespaceSpec holds info about if/how to manage a namespace for this location. Prefix and
// secrets inherit down the location tree; manage does not -- owning a namespace is a statement
// about one location.
type LocationNamespaceSpec struct {
	// Manage indicates if a namespace corresponding to this location should be created and owned
	// by this object. Not inherited.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="namespace.manage is immutable"
	// +optional
	Manage bool `json:"manage,omitempty"`
	// Prefix is prepended to the location name to form the namespace name, as in
	// `fmt.Sprintf("%s-%s", prefix, name)`. Inherited.
	// +kubebuilder:validation:MaxLength=32
	// +optional
	Prefix string `json:"prefix,omitempty"`
	// Secrets are names of replicated secrets (from the manager's mount or labelled Secrets) to
	// place in the namespace. Inherited down the location tree as a union. Ignored when manage is
	// false.
	// +listType=set
	// +kubebuilder:validation:MaxItems=64
	// +optional
	Secrets []string `json:"secrets,omitempty"`
}
