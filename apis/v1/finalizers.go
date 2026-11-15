package v1

// Finalizer is a gino finalizer.
type Finalizer string

// enumerations of Finalizer.
const (
	FinalizerMessaging   Finalizer = "gino.works/messaging-cleanup"
	FinalizerAssignments Finalizer = "gino.works/assignment-cleanup"
	FinalizerCluster     Finalizer = "gino.works/cluster-cleanup"
)
