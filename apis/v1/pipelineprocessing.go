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
	WorkerAugments    `json:",inline"`
}

// EnrichSpec holds configuration data for the enrichment tier.
type EnrichSpec struct {
	WorkerSpec        `json:",inline"`
	WorkerLoggingSpec `json:",inline"`
	WorkerAugments    `json:",inline"`
}

// LoadSpec holds configuration data for the loading tier.
type LoadSpec struct {
	WorkerSpec        `json:",inline"`
	WorkerLoggingSpec `json:",inline"`
	WorkerAugments    `json:",inline"`
}

// EnvVar is a simplified env var to use with gino workers.
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
