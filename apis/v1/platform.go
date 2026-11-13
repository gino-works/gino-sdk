package v1

import apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Platform is the way to configure available inputs for each vendor/profile/platform. This data
// is given to workers to know how to collect things from their target Endpoints.
// Short name gplat -> "gino platform".
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=platforms,scope=Cluster,shortName=gplat
type Platform struct {
	apimachineryapismetav1.TypeMeta   `json:",inline"`
	apimachineryapismetav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlatformSpec   `json:"spec"`
	Status PlatformStatus `json:"status,omitempty"`
}

// PlatformList is the list type for Platforms.
// +kubebuilder:object:root=true
type PlatformList struct {
	apimachineryapismetav1.TypeMeta `json:",inline"`
	apimachineryapismetav1.ListMeta `json:"metadata,omitempty"`

	Items []Platform `json:"items"`
}

// PlatformSpec is the spec for a Platform resource.
type PlatformSpec struct {
	// Match defines what vendor/platform/profile endpoints this Platform targets.
	Match PlatformMatch `json:"match"`
	// Inputs define the inputs that this platform supports.
	Inputs map[string]PlatformInput `json:"inputs"`
	// Transform defines how to transform inputs for this platform. Transforms are executed in order
	// of the list of PlatformTransform for the given input.
	Transform map[string][]PlatformTransform `json:"transform"`
}

// PlatformStatus is the status for a Platform resource.
type PlatformStatus struct{}

// PlatformMatch defines matching information to align a Platform to an Endpoint.
type PlatformMatch struct {
	// Vendor is the name of the vendor to match on.
	Vendor string `json:"vendor"`
	// Platform is the name of the platform to match on -- wildcard is allowed ("*").
	Platform string `json:"platform"`
	// Profile is the name of the profile to match on -- wildcard is allowed ("*").
	Profile string `json:"profile"`
}

// PlatformInput defines an input for this platform.
// +kubebuilder:validation:XValidation:rule="(has(self.cli) ? 1 : 0) + (has(self.http) ? 1 : 0) == 1",message="Exactly one input type must be specified"
type PlatformInput struct {
	// TCP is the tcp flavor of input.
	// +optional
	TCP *TCPPlatformInput `json:"tcp,omitempty"`
	// Cli is the cli flavor of input.
	// +optional
	Cli *CliPlatformInput `json:"cli,omitempty"`
	// HTTP is the http flavor of input.
	// +optional
	HTTP *HTTPPlatformInput `json:"http,omitempty"`
}

// TCPPlatformInput holds input data for a tcp input.
type TCPPlatformInput struct {
	// Probes are all run, in order, against the endpoint.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	Probes []TCPProbe `json:"probes"`
}

// TCPProbe holds data about a tcp probe.
type TCPProbe struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int `json:"port"`
}

// CliPlatformInput holds input data for a cli input.
type CliPlatformInput struct {
	// Commands are all run, in order, on the endpoint's session as one unit. Their outputs are
	// handed to the transform together, keyed by name.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	Commands []CliCommand `json:"commands"`
}

// CliCommand is a single command/input to a cli device.
type CliCommand struct {
	// Name -- name that transformers read to map to a transform.
	Name string `json:"name"`
	// Input -- the actual input sent to the device.
	Input string `json:"input"`
}

// HTTPPlatformInput holds input data for a http input.
type HTTPPlatformInput struct{}

// PlatformTransform holds information about how to transform a given input.
// +kubebuilder:validation:XValidation:rule="(has(self.fn) ? 1 : 0) == 1",message="Exactly one transform type must be specified"
type PlatformTransform struct {
	Fn *string `json:"fn,omitempty"`
}

// Kind returns the collector kind of whichever implementation is set. The CRD's CEL rule
// guarantees exactly one is; an empty PlatformInput (only possible from Go) returns "".
func (p *PlatformInput) Kind() CollectorKind {
	switch {
	case p.Cli != nil:
		return CollectorKindCLI
	case p.HTTP != nil:
		return CollectorKindHTTP
	case p.TCP != nil:
		return CollectorKindTCP
	default:
		return ""
	}
}
