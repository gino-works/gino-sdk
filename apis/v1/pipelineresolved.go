package v1

// PipelineResolved is what a Pipeline wants collected, before deciding who collects it:
// every (endpoint, kind) pair the pipeline owns, grouped by the cluster that will realize it.
type PipelineResolved struct {
	// Unscheduled explains every selected endpoint that produced no work.
	// +optional
	Unscheduled []UnscheduledEndpoint `json:"unscheduled"`
	// Models is the list of resolved model name/versions for this pipeline.
	Models []ResolvedModel `json:"models"`
	// Targets[cluster][kind] to the endpoints and their inputs for that container type.
	Targets map[string]map[CollectorKind][]AssignedEndpoint `json:"targets"`
	// MissingModelCoverage is a mapping showing the platform -> model -> inputs that *cannot* be
	// provided/resolved.
	MissingModelCoverage map[string]map[string][]string `json:"missingModelCoverage"`
}

// NewPipelineResolved returns a fresh PipelineResolved.
func NewPipelineResolved() PipelineResolved {
	return PipelineResolved{
		MissingModelCoverage: map[string]map[string][]string{},
		Targets:              map[string]map[CollectorKind][]AssignedEndpoint{},
	}
}

// ResolvedModel is a fully resolved model.
type ResolvedModel struct {
	// Name is the name of the model.
	Name string `json:"name"`
	// Version is the version of the model.
	Version string `json:"version"`
}

// AssignedInput is one input to collect from one endpoint, fully resolved: the platform's
// implementation, the strictest requirement across every model that wants it, and who wants it.
type AssignedInput struct {
	Name      string           `json:"name"`
	Kind      CollectorKind    `json:"kind"`
	Impl      PlatformInput    `json:"impl"`
	Require   ModelRequirement `json:"require"`
	Optional  bool             `json:"optional,omitempty"` // optional for *every* consumer
	Consumers []ResolvedModel  `json:"consumers"`          // model names
}

// AssignedEndpoint is one endpoint's worth of work for one collector kind.
type AssignedEndpoint struct {
	WorkerAssignmentEndpoint `json:",inline"`

	Inputs []AssignedInput `json:"inputs"`
}

// UnscheduledEndpoint holds reason why an endpoint (and which endpoint) cant be scheduled in this
// Pipeline.
type UnscheduledEndpoint struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	// Kind is set when only one collector kind is affected (NoCredentials); empty when the whole
	// endpoint is.
	// +optional
	Kind   CollectorKind       `json:"kind,omitempty"`
	Reason UnschedulableReason `json:"reason"`
}
