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

// TargetCount returns total count of target endpoints for this pipeline.
func (r PipelineResolved) TargetCount() int {
	var count int

	for _, kinds := range r.Targets {
		for _, endpoints := range kinds {
			count += len(endpoints)
		}
	}

	return count
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
	ModelInput string           `json:"modelInput"`
	Kind       CollectorKind    `json:"kind"`
	Tasks      []AssignedTask   `json:"tasks"`
	Require    ModelRequirement `json:"require"`
	Consumers  []ResolvedModel  `json:"consumers"` // model versions wanting this input
}

// AssignedTask is a singular resolved task for an assignment.
// +kubebuilder:validation:XValidation:rule="(has(self.cli) ? 1 : 0) + (has(self.http) ? 1 : 0) + (has(self.tcp) ? 1 : 0) == 1",message="Exactly one input type must be specified"
type AssignedTask struct {
	// Name is redundant here since the tasks include the name but this is used for list sorting in
	// the AssignedInput Tasks list.
	Name string `json:"name"`

	Cli  *CliCommand  `json:"cli,omitempty"`
	HTTP *HTTPRequest `json:"http,omitempty"`
	TCP  *TCPProbe    `json:"tcp,omitempty"`
}

// Impl returns the set task object for the AssignedTask.
func (t *AssignedTask) Impl() any {
	switch {
	case t.Cli != nil:
		return t.Cli
	case t.HTTP != nil:
		return t.HTTP
	case t.TCP != nil:
		return t.TCP
	default:
		// shouldnt happen ever since cel rules will ensure one is set
		return nil
	}
}

// AssignedEndpoint is one endpoint's worth of work for one collector kind.
type AssignedEndpoint struct {
	WorkerAssignmentEndpoint `json:",inline"`

	Inputs []AssignedInput `json:"inputs"`
}

// UnscheduledEndpoint holds reason why an endpoint (and which endpoint) cant be scheduled in this
// Pipeline.
type UnscheduledEndpoint struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Reason    UnschedulableReason `json:"reason"`
}
