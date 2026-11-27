package v1

import (
	apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnrichEmit is when a model version's document is emitted, see ModelEnrichEmitSpec.
// +kubebuilder:validation:Enum=allRefreshed;anyRefreshed;inputs;passthrough
type EnrichEmit string

// enumerations of EnrichEmit.
const (
	EnrichEmitAllRefreshed EnrichEmit = "allRefreshed"
	EnrichEmitAnyRefreshed EnrichEmit = "anyRefreshed"
	EnrichEmitInputs       EnrichEmit = "inputs"
	EnrichEmitPassthrough  EnrichEmit = "passthrough"
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
// +kubebuilder:validation:XValidation:rule="!has(self.enrichSpec) || !has(self.enrichSpec.emit) || !has(self.enrichSpec.emit.inputs) || self.enrichSpec.emit.inputs.all(i, self.inputs.exists(x, x.name == i))",message="enrichSpec.emit.inputs may only name inputs of this version"
// +kubebuilder:validation:XValidation:rule="!has(self.enrichSpec) || !has(self.enrichSpec.emit) || !has(self.enrichSpec.emit.type) || self.enrichSpec.emit.type != 'passthrough' || size(self.inputs) == 1",message="passthrough models take exactly one input"
type ModelVersion struct {
	Version string `json:"version"`
	// Inputs are every input this version needs; all are required.
	// +kubebuilder:validation:MaxItems=64
	Inputs []ModelInput `json:"inputs"`
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
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name"`
	// Require defines when or how fresh the model requires this input to be.
	Require ModelRequirement `json:"require"`
}

// ModelRequirement holds info about when/at what time the model wants this input.
// +kubebuilder:validation:XValidation:rule="(has(self.freshness) ? 1 : 0) + (has(self.at) ? 1 : 0) == 1",message="Exactly one requirement type must be specified"
// +kubebuilder:validation:XValidation:rule="!has(self.splay) || has(self.at)",message="splay is only allowed with at"
type ModelRequirement struct {
	Freshness *apimachineryapismetav1.Duration `json:"freshness,omitempty"`
	At        *string                          `json:"at,omitempty"`
	// Splay helps spread out the "at" requirements for the pipeline. Namespace/Name of the endpoint
	// is hashed and used as offset from that at time to a max of the splay value.
	// +kubebuilder:validation:XValidation:rule="duration(self) > duration('0s') && duration(self) <= duration('24h')",message="splay must be greater than 0 and at most 24h"
	// +optional
	Splay *apimachineryapismetav1.Duration `json:"splay,omitempty"`
}

// ModelEnrichSpec defines when the enricher produces a document for this model version.
type ModelEnrichSpec struct {
	// Emit is when a document is enriched and published.
	// +kubebuilder:default={}
	// +optional
	Emit ModelEnrichEmitSpec `json:"emit,omitzero"`
}

// ModelEnrichEmitSpec is when a document is enriched (except passthrough) and published. Whatever
// the type, nothing is emitted until every input of the version has arrived at least once.
// +kubebuilder:validation:XValidation:rule="(has(self.type) && self.type == 'inputs') == (has(self.inputs) && size(self.inputs) > 0)",message="inputs must be set (non empty) when type is inputs, and only then"
type ModelEnrichEmitSpec struct {
	// Type is when to emit:
	//   allRefreshed -- once every input has been collected again since the last document: one
	//                   document per full collection cycle.
	//   anyRefreshed -- on every arrival: freshest values, one document per input arrival.
	//   inputs       -- on every arrival of one of ModelInputs; arrivals of other inputs only
	//                   update what the next document is built from.
	//   passthrough  -- basically no cache write/input lookups involved, gino side will simply emit
	//                   the message straight back out for a loader to pick it up.
	// +kubebuilder:default=allRefreshed
	// +optional
	Type EnrichEmit `json:"type,omitempty"`
	// ModelInputs are the inputs whose arrival emits a document; set when (and only when) Type is
	// inputs. Names must be inputs of this model version.
	// +listType=set
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=63
	// +optional
	ModelInputs []string `json:"inputs,omitempty"`
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
