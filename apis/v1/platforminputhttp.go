package v1

import apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// HTTPOptions options/settings for http input(s). When set on the "Input" is the default
// for every child task, tasks override this if set, otherwise inherit it.
type HTTPOptions struct {
	Timeout *apimachineryapismetav1.Duration `json:"timeout,omitempty"`
}

// HTTPPlatformInput holds input data for a http input. HTTP is too varied to describe (auth,
// paging, token refresh, rate limits...), so an input names implementation functions the
// plugin provides, and the plugin does whatever that API needs.
type HTTPPlatformInput struct {
	HTTPOptions `json:",inline"`

	// Requests are all run, in order, against the endpoint's session. Each request is one task of the input, producing one (or, for fan-out APIs, many) payloads that are
	// transformed on their own; the input is complete once every task is present.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	Requests []HTTPRequest `json:"requests"`
}

// HTTPRequest is one named call the plugin knows how to make.
type HTTPRequest struct {
	HTTPOptions `json:",inline"`

	// Name is the task name: unique within the input, what transforms switch on and what the
	// enricher caches the task under.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Name string `json:"name"`
	// Fn is the plugin function that makes the request.
	Fn string `json:"fn"`
	// Args are passed to Fn as-is; whatever that function documents (a path, an account id, a
	// page size). Keeps one Fn serving several inputs without a code change.
	// +optional
	Args map[string]string `json:"args,omitempty"`
}
