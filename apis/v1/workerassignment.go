package v1

import (
	apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IndexWorkerAssignmentSecret is the index field for looking up what secrets are mapped to what
// workeraissgnment objects.
const IndexWorkerAssignmentSecret = ".spec.endpoints[*].credentials, distinct" //nolint: gosec

// WorkerAssignment is maybe the most important gino crd -- it defines how and where workers get
// placed -- the workers being the *thing* that does the collection of data, ingesting it into the
// gino pipeline.
// Short name gwa -> "gino worker assignment".
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path="workerassignments"
// +kubebuilder:resource:shortName=gwa
type WorkerAssignment struct {
	apimachineryapismetav1.TypeMeta   `json:",inline"`
	apimachineryapismetav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkerAssignmentSpec   `json:"spec"`
	Status WorkerAssignmentStatus `json:"status,omitempty"`
}

// WorkerAssignmentList is the list type for WorkerAssignments.
// +kubebuilder:object:root=true
type WorkerAssignmentList struct {
	apimachineryapismetav1.TypeMeta `json:",inline"`
	apimachineryapismetav1.ListMeta `json:"metadata,omitempty"`

	Items []WorkerAssignment `json:"items"`
}

// WorkerAssignmentSpec holds assignment data.
type WorkerAssignmentSpec struct {
	// Pipeline is the name of the pipeline this assignment is a child of.
	Pipeline string `json:"pipeline"`
	// Cluster is the name of the cluster this assignment runs in.
	Cluster string `json:"cluster"`
	// Kind is the kind of the collector this worker will realize.
	Kind CollectorKind `json:"kind"`
	// Shard is this assignment's index among the assignments for the same pipeline/cluster/kind in
	// this namespace. It is how placement stays sticky across reconciles: endpoints already on a
	// shard stay there. Always 0 when shardSize is 1 (one assignment per endpoint).
	// +optional
	Shard int32 `json:"shard,omitempty"`
	// WorkerSpec holds image/resource spec information for this worker.
	WorkerSpec WorkerSpec `json:"workerSpec"`
	// Logging holds the logging configuration for this worker.
	// +optional
	Logging WorkerLoggingSpec `json:"logging,omitempty"`
	// InputSets are the resolved inputs for this kind, keyed by Platform name. Every endpoint on
	// a platform collects the same inputs, so they are stored once and referenced.
	InputSets map[string][]AssignedInput `json:"inputSets"`
	// Endpoints hold the endpoint info required for the worker to work.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	Endpoints []WorkerAssignmentEndpoint `json:"endpoints"`
}

// WorkerAssignmentEndpoint is a single endpoint this worker works for, including the info about the
// endpoint needed to do the work. This is the final/resolved content that also rides along with all
// messages related to this endpoint (such that downstream workers have "full" information about
// the endpoint).
type WorkerAssignmentEndpoint struct {
	// Name is the endpoint name.
	Name string `json:"name"`
	// Namespace is the namespace the endpoint lives in.
	Namespace string `json:"namespace"`
	// Location is the resolved (final) location of the endpoint (if set).
	// +optional
	Location string `json:"location,omitempty"`
	// LocationType is the resolved location type (if set).
	// +optional
	LocationType string `json:"locationType,omitempty"`
	// Host is the host to connect to for this endpoint.
	Host string `json:"host"`
	// Vendor, Platform and Profile are the endpoint's resolved identity. They ride on every
	// collected message as metadata; the Platform CR that applies is InputSet.
	Vendor   string `json:"vendor,omitempty"`
	Platform string `json:"platform,omitempty"`
	Profile  string `json:"profile,omitempty"`
	// PlatformName is the key into spec.inputSets -- the matched Platform's name.
	PlatformName string `json:"platformName"`
	// Credentials names the Secret in the endpoint's (and so the worker's) namespace holding this
	// endpoint's credentials for this kind.
	Credentials string            `json:"credentials,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	// InterfaceTags are the resolved interface classification rules, evaluated by the enricher
	// against the interfaces it sees; they ride on collected messages like Tags.
	// +listType=atomic
	// +optional
	InterfaceTags []InterfaceTag `json:"interfaceTags,omitempty"`
}

// WorkerAssignmentStatus is the status for a WorkerAssignment resource.
type WorkerAssignmentStatus struct{}
