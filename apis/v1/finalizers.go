package v1

// Finalizer is a gino finalizer.
type Finalizer string

// enumerations of Finalizer.
const (
	FinalizerMessaging Finalizer = "gino.io/messaging-cleanup"
)
