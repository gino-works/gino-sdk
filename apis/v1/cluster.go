package v1

import (
	apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ClusterNameCore is the name of the Cluster representing the core cluster itself. It is created
	// by the core init and is never relayed to.
	ClusterNameCore = "core"
	// ClusterMessagingSecretPrefix is the prefix for the messaging secret for a cluster -- all
	// workers share a single messaging secret for each remote cluster (or local), this is the
	// prefix of that name (the remaining bit is the cluster name).
	ClusterMessagingSecretPrefix = "gino-worker-messaging-"
)

// Cluster is a cluster that realizes work: the core cluster itself (named ClusterNameCore, created
// at init) or a remote one. Reconciling it ensures the cluster's messaging identity (a cert-manager
// Certificate whose Secret workers in that cluster mount) and, for remotes, the consumers that
// carry configuration from core.
// Short name gc -> "gino cluster".
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=clusters,scope=Cluster,shortName=gc
type Cluster struct {
	apimachineryapismetav1.TypeMeta   `json:",inline"`
	apimachineryapismetav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec"`
	Status ClusterStatus `json:"status,omitempty"`
}

// ClusterList is the list type for Clusters.
// +kubebuilder:object:root=true
type ClusterList struct {
	apimachineryapismetav1.TypeMeta `json:",inline"`
	apimachineryapismetav1.ListMeta `json:"metadata,omitempty"`

	Items []Cluster `json:"items"`
}

// ClusterSpec is the spec for a Cluster resource.
type ClusterSpec struct{}

// ClusterStatus is the status for a Cluster resource.
type ClusterStatus struct{}

// MessagingSecretNameFor returns the name of the Secret holding the named cluster's worker
// messaging identity: issued (by the Cluster controller, via cert-manager) in the manager
// namespace, copied by the realizer into every namespace that runs workers for that cluster, and
// mounted by those workers. A WorkerAssignment does not carry it; it is derived from spec.cluster.
func MessagingSecretNameFor(cluster string) string {
	return ClusterMessagingSecretPrefix + cluster
}

// MessagingSecretName returns MessagingSecretNameFor this cluster.
func (c *Cluster) MessagingSecretName() string {
	return MessagingSecretNameFor(c.Name)
}

// MessagingSecretName returns the name of the Secret the workers for this assignment mount for
// their messaging identity.
func (a *WorkerAssignment) MessagingSecretName() string {
	return MessagingSecretNameFor(a.Spec.Cluster)
}
