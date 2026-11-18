package v1

import apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// CliOptions options/settings for cli input(s). When set on the "Input" is the default
// for every child task, tasks override this if set, otherwise inherit it.
type CliOptions struct {
	Timeout *apimachineryapismetav1.Duration `json:"timeout,omitempty"`
}

// CliPlatformInput holds input data for a cli input.
type CliPlatformInput struct {
	CliOptions `json:",inline"`

	// Commands are all run, in order, on the endpoint's session. Each command is one part (task)
	// of the input: transformed on its own, cached by the enricher under its name, and the input
	// is complete once every part is present.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	Commands []CliCommand `json:"commands"`
}

// CliCommand is a single command/input to a cli device.
type CliCommand struct {
	CliOptions `json:",inline"`

	// Name is the part (task) name: unique within the input, what transforms switch on and what
	// the enricher caches the part under.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Name string `json:"name"`
	// Input -- the actual input sent to the device.
	Input string `json:"input"`
}
