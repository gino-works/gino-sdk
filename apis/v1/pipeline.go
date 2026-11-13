package v1

import (
	apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CollectorKind is the kind of collector/protocol used to talk to an endpoint.
// +kubebuilder:validation:Enum=tcp;cli;http
type CollectorKind string

// enumerations of CollectorKind.
const (
	CollectorKindTCP  CollectorKind = "tcp"
	CollectorKindCLI  CollectorKind = "cli"
	CollectorKindHTTP CollectorKind = "http"
)

// Pipeline is an object that defines how endpoints get mapped to a pipeline -- that is the process
// of collecting, normalizing, and loading data. The result of a Pipeline being reconciled is the
// WorkerAssignment crs that hold the resolved information about the collection deployment(s), the
// consumers(queues) for the messaging service, and the deployments that handle the transform/enrich
// and load components.
// Short name gpipe -> "gino pipeline".
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=pipelines,scope=Cluster,shortName=gpipe
// +kubebuilder:validation:XValidation:rule="self.metadata.name.matches('^[a-z0-9]([-a-z0-9]*[a-z0-9])?$') && size(self.metadata.name) <= 24",message="metadata.name must be a valid DNS label (no dots, max 24 characters)"
// +kubebuilder:printcolumn:name="Resolved",type=string,JSONPath=".status.conditions[?(@.type=='Resolved')].status"
type Pipeline struct {
	apimachineryapismetav1.TypeMeta   `json:",inline"`
	apimachineryapismetav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec"`
	Status PipelineStatus `json:"status,omitempty"`
}

// PipelineList is the list type for Pipeline.
// +kubebuilder:object:root=true
type PipelineList struct {
	apimachineryapismetav1.TypeMeta `json:",inline"`
	apimachineryapismetav1.ListMeta `json:"metadata,omitempty"`

	Items []Pipeline `json:"items"`
}

// PipelineSpec is the spec for a Pipeline resource.
type PipelineSpec struct {
	// Selector is how this Pipeline selects sites or endpoints it will manifest the workers for.
	Selector TargetSelector `json:"selector"`

	// Models is the listing of models this pipeline is setup to collect/emit.
	Models []PipelineModel `json:"models"`

	// Collection defines how to collect things for this Pipeline.
	Collection CollectionSpec `json:"collection"`

	// Processing defines how processing should be setup for this Pipeline.
	Processing ProcessingSpec `json:"processing"`
}

// PipelineModel is a model the Pipeline is working to build.
type PipelineModel struct {
	// Name of the model.
	Name string `json:"name"`
	// Version the version of the model -- defaults to the models target version if unset.
	// +optional
	Version string `json:"version,omitempty"`
}

// PipelineStatus is the spec for a Pipeline resource.
type PipelineStatus struct {
	// ObservedGeneration is the metadata.generation this status was computed from. The apiserver
	// bumps metadata.generation on every spec change (not on status or metadata changes), so a
	// reader can tell whether status reflects the current spec: if observedGeneration <
	// metadata.generation the controller has not caught up yet and the rest of status is stale.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions summarize the pipeline's state. Known types:
	//   - Resolved: TODO
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []apimachineryapismetav1.Condition `json:"conditions,omitempty" patchMergeKey:"type" patchStrategy:"merge"`

	// Resolved holds the fully resolved pipeline data.
	Resolved *PipelineResolved `json:"resolved,omitempty"`
}

// MessagingSecretName returns MessagingSecretNameFor this pipeline.
func (p *Pipeline) MessagingSecretName() string {
	return MessagingSecretNameFor(p.Name)
}
