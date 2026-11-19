package v1

import (
	apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnrichEmit is when a model version's document is emitted.
type EnrichEmit string

// enumerations of EnrichEmit.
const (
	EnrichEmitAllRefreshed EnrichEmit = "allRefreshed"
	EnrichEmitAnyRefreshed EnrichEmit = "anyRefreshed"
)

// Model is an object that holds information about a gino model -- the actual thing that we are
// trying to produce -- the schema that gets written into caches/databases.
// Short name gm -> "gino model".
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=models,scope=Cluster,shortName=gm
type Model struct {
	apimachineryapismetav1.TypeMeta   `json:",inline"`
	apimachineryapismetav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelSpec   `json:"spec"`
	Status ModelStatus `json:"status,omitempty"`
}

// ModelList is the list type for Models.
// +kubebuilder:object:root=true
type ModelList struct {
	apimachineryapismetav1.TypeMeta `json:",inline"`
	apimachineryapismetav1.ListMeta `json:"metadata,omitempty"`

	Items []Model `json:"items"`
}

// ModelSpec is the spec for a Model resource.
// +kubebuilder:validation:XValidation:rule="self.versions.exists(v, v.version == self.targetVersion)",message="Target Version must be in list of Versions"
type ModelSpec struct {
	// TargetVersion is the preferred/target version of this model, required to ensure previous
	// model versions are still supported for migration period.
	TargetVersion string `json:"targetVersion"`
	// Versions hold the supported versions -- TargetVersion *must* be present.
	// +listType=map
	// +listMapKey=version
	// +kubebuilder:validation:MaxItems=16
	Versions []ModelVersion `json:"versions"`
}

// ModelVersion holds a version of the Model.
type ModelVersion struct {
	Version string       `json:"version"`
	Inputs  []ModelInput `json:"inputs"`
	// +kubebuilder:default={}
	// +optional
	EnrichSpec ModelEnrichSpec `json:"enrichSpec,omitzero"`
	// +kubebuilder:default={}
	// +optional
	LoadSpec ModelLoadSpec `json:"loadSpec,omitzero"`
}

// ModelInput defines an input this model wants -- the input name and when/how often it wants the
// input.
type ModelInput struct {
	// Name is the name of the input this model wants.
	Name string `json:"name"`
	// Require defines when or how fresh the model requires this input to be.
	Require ModelRequirement `json:"require"`
}

// ModelRequirement holds info about when/at what time the model wants this input.
// +kubebuilder:validation:XValidation:rule="(has(self.freshness) ? 1 : 0) + (has(self.at) ? 1 : 0) == 1",message="Exactly one requirement type must be specified"
type ModelRequirement struct {
	Freshness *apimachineryapismetav1.Duration `json:"freshness,omitempty"`
	At        *string                          `json:"at,omitempty"`
}

// ModelEnrichSpec defines when the enricher produces a document for this model version.
type ModelEnrichSpec struct {
	// Emit is when a document is enriched and published once every input has arrived at least
	// once: allRefreshed waits until every input has been collected again since the last
	// document (one document per full cycle); anyRefreshed emits on every arrival (freshest
	// values, a document per input arrival).
	// +kubebuilder:validation:Enum=allRefreshed;anyRefreshed
	// +kubebuilder:default=allRefreshed
	// +optional
	Emit EnrichEmit `json:"emit,omitempty"`
}

// ModelLoadSpec defines the criteria for when a loader passes the received enriched messages to
// the loader server.
type ModelLoadSpec struct {
	// MaxWait is how long the loader holds a batch open waiting for more documents. It is capped so
	// that the loader's consumer ack wait can be a constant that always exceeds it.
	// +kubebuilder:default="3m"
	// +kubebuilder:validation:XValidation:rule="duration(self) > duration('0s') && duration(self) <= duration('5m')",message="maxWait must be greater than 0 and at most 5m"
	// +optional
	MaxWait apimachineryapismetav1.Duration `json:"maxWait,omitzero"`
	// MaxCount is how many documents the loader batches before loading; whichever of MaxWait or MaxCount is hit first wins.
	// +kubebuilder:default=100
	// +kubebuilder:validation:Minimum=1
	// +optional
	MaxCount int `json:"maxCount,omitzero"`
}

// ModelStatus is the spec for a Model resource.
type ModelStatus struct{}
