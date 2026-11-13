package v1

import (
	"strings"

	k8scorev1 "k8s.io/api/core/v1"
)

// Image is a container image reference; unset fields inherit from the enclosing level.
type Image struct { //nolint: recvcheck
	// Repository is the image repository to use.
	// +optional
	Repository string `json:"repository,omitempty"`
	// Tag is the image tag to use.
	// +optional
	Tag string `json:"tag,omitempty"`
	// Digest pins the image (sha256:...); when set it is used instead of Tag.
	// +kubebuilder:validation:Pattern=`^sha256:[a-f0-9]{64}$`
	// +optional
	Digest string `json:"digest,omitempty"`
	// PullPolicy is the image pull policy for the worker.
	// +kubebuilder:validation:Enum=IfNotPresent;Always;Never
	// +optional
	PullPolicy string `json:"pullPolicy,omitempty"`
}

// WorkerSpec holds common worker spec data -- image and resources for example.
type WorkerSpec struct {
	// Image is the override image to set on this worker, defaulted to whatever the controller was
	// deployed w/ for this otherwise.
	// +optional
	Image Image `json:"image,omitempty"`
	// Resources define standard k8s resources for the worker.
	// +optional
	Resources k8scorev1.ResourceRequirements `json:"resources,omitempty"`
}

// Merge returns i with every set field of over laid on top.
func (i Image) Merge(over Image) Image {
	if over.Repository != "" {
		i.Repository = over.Repository
	}

	if over.Tag != "" {
		i.Tag = over.Tag
	}

	if over.Digest != "" {
		i.Digest = over.Digest
	}

	if over.PullPolicy != "" {
		i.PullPolicy = over.PullPolicy
	}

	return i
}

// Ref renders the image reference to run: repository@digest when pinned, else repository:tag,
// else just the repository.
func (i Image) Ref() string {
	switch {
	case i.Digest != "":
		return i.Repository + "@" + i.Digest
	case i.Tag != "":
		return i.Repository + ":" + i.Tag
	default:
		return i.Repository
	}
}

// ParseImageRef splits a reference as Ref renders it -- "repo", "repo:tag" or "repo@digest" --
// into an Image. A ':' before the last '/' is a registry port, not a tag.
func ParseImageRef(ref string) Image {
	if repository, digest, ok := strings.Cut(ref, "@"); ok {
		return Image{Repository: repository, Digest: digest}
	}

	slash := strings.LastIndex(ref, "/")

	if colon := strings.LastIndex(ref, ":"); colon > slash {
		return Image{Repository: ref[:colon], Tag: ref[colon+1:]}
	}

	return Image{Repository: ref}
}

// WorkerLoggingSpec holds common worker logging configs.
type WorkerLoggingSpec struct {
	// LogLevel sets the logger level.
	// +kubebuilder:validation:Enum=error;warn;info;debug
	// +kubebuilder:default=warn
	// +optional
	LogLevel string `json:"logLevel,omitempty"`
	// LogDevMode enables the dev mode with pretty printing/colors.
	// +optional
	LogDevMode bool `json:"logDevMode,omitempty"`
	// LogRecorderMode enables the logger to record messages and emit them when an error is
	// encountered.
	// +optional
	LogRecorderMode bool `json:"logRecorderMode,omitempty"`
}
