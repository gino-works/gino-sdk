package v1

// ProcessingSpec holds configuration data for the processing tiers.
type ProcessingSpec struct {
	// TargetsPerReplica determines the replica count for transform, enrich, and load
	// workers when explicit replica counts are not configured. The replica count is
	// calculated by dividing the total number of pipeline targets by this value,
	// rounding up to the nearest integer.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	// +kubebuilder:default=25
	TargetReplicaMultiplier int `json:"targetReplicaMultiplier,omitempty"`

	Transform TransformSpec `json:"transform"`
	Enrich    EnrichSpec    `json:"enrich"`
	Load      LoadSpec      `json:"load"`
}

// TransformSpec holds configuration data for the transformation tier.
type TransformSpec struct {
	WorkerSpec        `json:",inline"`
	WorkerLoggingSpec `json:",inline"`

	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=1
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
}

// EnrichSpec holds configuration data for the enrichment tier.
type EnrichSpec struct {
	WorkerSpec        `json:",inline"`
	WorkerLoggingSpec `json:",inline"`

	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=1
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
}

// LoadSpec holds configuration data for the loading tier.
type LoadSpec struct {
	WorkerSpec        `json:",inline"`
	WorkerLoggingSpec `json:",inline"`

	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=1
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Secrets hold any secrets required to be mounted on the loader -- assumed to be in the main
	// gino namespace where the loader will be running.
	Secrets []string `json:"secrets,omitempty"`

	// ExtraEnv are extra env vars that will be mounted on the loader.
	ExtraEnv []EnvVar `json:"extraEnv,omitempty"`
}

// EnvVar is a simplified env var to use with gino workers.
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
