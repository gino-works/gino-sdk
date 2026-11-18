package v1

// TCPOptions options/settings for tcp input(s).
type TCPOptions struct{}

// TCPPlatformInput holds input data for a tcp input. When set on the "Input" is the default
// for every child task, tasks override this if set, otherwise inherit it.
type TCPPlatformInput struct {
	TCPOptions `json:",inline"`

	// Probes are all run, in order, against the endpoint. Each probe is one task of the
	// input: transformed on its own, cached by the enricher under its name, and the input is
	// complete once every task is present.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	Probes []TCPProbe `json:"probes"`
}

// TCPProbe holds data about a tcp probe.
type TCPProbe struct {
	TCPOptions `json:",inline"`

	// Name is the task name: unique within the input, what transforms switch on and what the
	// enricher caches the task under.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Name string `json:"name"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int `json:"port"`
}
