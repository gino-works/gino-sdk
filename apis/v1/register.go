package v1

import (
	ginosdkapis "github.com/gino-works/gino-sdk/apis"
	apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachineryruntime "k8s.io/apimachinery/pkg/runtime"
	apimachineryruntimeschema "k8s.io/apimachinery/pkg/runtime/schema"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Version is the API version.
	Version = "v1"
)

var (
	schemeBuilder      = apimachineryruntime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &schemeBuilder
)

// SchemeGroupVersion is group version used to register these objects.
var SchemeGroupVersion = apimachineryruntimeschema.GroupVersion{
	Group:   ginosdkapis.Group,
	Version: Version,
}

// AddToScheme adds the types in this group-version to the given scheme.
var AddToScheme = localSchemeBuilder.AddToScheme

func addKnownTypes(scheme *apimachineryruntime.Scheme) error {
	_, objects := GetAPIs()

	for _, object := range objects {
		scheme.AddKnownTypes(SchemeGroupVersion, object)
	}

	apimachineryapismetav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}

// GetAPIs returns the information necessary to register this package's types to a scheme.
func GetAPIs() (apimachineryruntimeschema.GroupVersion, []apimachineryruntime.Object) {
	return SchemeGroupVersion, []apimachineryruntime.Object{
		&Cluster{},
		&ClusterList{},
		&Endpoint{},
		&EndpointList{},
		&EndpointGroup{},
		&EndpointGroupList{},
		&Location{},
		&LocationList{},
		&Model{},
		&ModelList{},
		&Pipeline{},
		&PipelineList{},
		&Platform{},
		&PlatformList{},
		&WorkerAssignment{},
		&WorkerAssignmentList{},
	}
}

// GetObjects returns the objects (non list types) for gino.
func GetObjects() []ctrlruntimeclient.Object {
	return []ctrlruntimeclient.Object{
		&Cluster{},
		&Endpoint{},
		&EndpointGroup{},
		&Location{},
		&Model{},
		&Pipeline{},
		&Platform{},
		&WorkerAssignment{},
	}
}
