package v1

// Condition types. All are positive-polarity: True is the good state.
const (
	// ConditionResolved: groups and location have been applied and the effective
	// configuration is usable.
	ConditionResolved = "Resolved"
	// ConditionUnschedulable: used for pipeline targets that cannot be scheduled.
	ConditionUnschedulable = "Unschedulable"
	// ConditionSecretsReplicated: every secret a managed location namespace asks for exists in
	// a source and has been copied in.
	ConditionSecretsReplicated = "SecretsReplicated"
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

// UnschedulableReason identifies why an endpoint selected by a pipeline produced no work. Listed
// in the order they are checked; the first that applies is reported.
// +kubebuilder:validation:Enum=Unresolved;Disabled;NoCluster;NoPlatformMatch;NoInputs;NoCredentials
type UnschedulableReason string

// Reasons for ConditionUnschedulable.
const (
	UnschedulableReasonUnresolved      UnschedulableReason = "Unresolved"
	UnschedulableReasonDisabled        UnschedulableReason = "Disabled"
	UnschedulableReasonNoCluster       UnschedulableReason = "NoCluster"
	UnschedulableReasonNoPlatformMatch UnschedulableReason = "NoPlatformMatch"
	UnschedulableReasonNoInputs        UnschedulableReason = "NoInputs"
	UnschedulableReasonNoCredentials   UnschedulableReason = "NoCredentials"
)
