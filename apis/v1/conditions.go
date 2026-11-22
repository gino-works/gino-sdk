package v1

// Condition types. All are positive-polarity: True is the good state.
const (
	// ConditionResolved: groups and location have been applied and the effective
	// configuration is usable.
	ConditionResolved = "Resolved"
	// ConditionUnschedulable: used for pipeline targets that cannot be scheduled.
	ConditionUnschedulable = "Unschedulable"
	// ConditionSecretsReplicated: every secret a managed location namespace asks for exists in
	// a source and has been copied in; on a pipeline, every secret its tiers mount is in the
	// manager namespace (already there, or replicated from the manager's replicate dir).
	ConditionSecretsReplicated = "SecretsReplicated"
	// ConditionReady: the assignment's worker Deployment exists and is available -- set by
	// whichever controller runs it (core for core assignments, the remote via the relay for
	// remote ones).
	ConditionReady = "Ready"
)

// Reasons for ConditionSecretsReplicated.
const (
	ReasonSecretsReplicated = "Replicated"
	ReasonSecretNotFound    = "SecretNotFound"
	ReasonNamespaceNotOwned = "NamespaceNotOwned"
	ReasonSecretNotOwned    = "SecretNotOwned"
)

// Reasons for ConditionResolved.
const (
	ReasonResolved         = "Resolved"
	ReasonLocationNotFound = "LocationNotFound"
	ReasonAncestry         = "Ancestry"
	ReasonBadSelector      = "BadSelector"
)

// Reasons for ConditionReady.
const (
	ReasonReady                 = "Ready"
	ReasonDeploymentProgressing = "DeploymentProgressing"
	ReasonSecretNotReady        = "SecretNotReady" // messaging identity or a credential not there yet
	ReasonRelayPending          = "RelayPending"   // published, remote has not reported yet
	ReasonRelayStale            = "RelayStale"     // remote stopped reporting
)

// UnschedulableReason identifies why an endpoint selected by a pipeline produced no work. Listed
// in the order they are checked; the first that applies is reported.
// +kubebuilder:validation:Enum=Unresolved;Disabled;NoCluster;NoPlatformMatch;NoModelCoverage;NoInputs
type UnschedulableReason string

// Reasons for ConditionUnschedulable.
const (
	UnschedulableReasonUnresolved      UnschedulableReason = "Unresolved"
	UnschedulableReasonDisabled        UnschedulableReason = "Disabled"
	UnschedulableReasonNoCluster       UnschedulableReason = "NoCluster"
	UnschedulableReasonNoPlatformMatch UnschedulableReason = "NoPlatformMatch"
	// UnschedulableReasonNoModelCoverage: the platform matched but provides every input of none of
	//  the pipeline's model versions, so no model can be produced for the endpoint. The gaps are in
	// PipelineResolved.MissingModelCoverage for the platform.
	UnschedulableReasonNoModelCoverage UnschedulableReason = "NoModelCoverage"
	// UnschedulableReasonNoInputs: the platform covers a model version but its inputs resolve to no
	// work (should not happen -- a covered model has at least one input the platform provides).
	UnschedulableReasonNoInputs UnschedulableReason = "NoInputs"
)
