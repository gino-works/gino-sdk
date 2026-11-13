package v1

// CollectionSpec is the spec for a Collection resource.
type CollectionSpec struct {
	// Workers holds configurations for the workers in this CollectionSpec.
	Workers CollectionWorkers `json:"workers"`
}

// CollectionWorkers holds configuration data (including if the worker should be enabled) about
// the workers for this CollectionSpec.
type CollectionWorkers struct {
	TCP  TCPWorkerSpec  `json:"tcp"`
	Cli  CliWorkerSpec  `json:"cli"`
	HTTP HTTPWorkerSpec `json:"http"`
}

// CollectionSharding holds sharding information for the collection. Its own type to be used
// across all worker types with `inline`.
type CollectionSharding struct {
	// ShardSize is the number of target endpoints per replica -- setting this to 1 means 1:1
	// endpoint to collector.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10000
	// +kubebuilder:default=10
	// +optional
	ShardSize int32 `json:"shardSize,omitempty"`
}

// TCPWorkerSpec defines configuration for tcp workers.
type TCPWorkerSpec struct {
	WorkerSpec         `json:",inline"`
	WorkerLoggingSpec  `json:",inline"`
	CollectionSharding `json:",inline"`
}

// CliWorkerSpec defines configuration for cli workers.
type CliWorkerSpec struct {
	WorkerSpec         `json:",inline"`
	WorkerLoggingSpec  `json:",inline"`
	CollectionSharding `json:",inline"`
}

// HTTPWorkerSpec defines configuration for http workers.
type HTTPWorkerSpec struct {
	WorkerSpec         `json:",inline"`
	WorkerLoggingSpec  `json:",inline"`
	CollectionSharding `json:",inline"`
}
