package definition

import (
	"fmt"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
)

// Option is a functional option for the Registry -- so we can set non required additional things
// on it if we want to use this to generate more manifests.
type Option func(r *Registry)

// WithClusters registers clusters to the registry.
func WithClusters(clusters []ginosdkapisv1.Cluster) Option {
	return func(r *Registry) {
		for idx := range clusters {
			r.clusters[clusters[idx].Name] = clusters[idx]
		}
	}
}

// WithLocations registers locations to the registry.
func WithLocations(locations []ginosdkapisv1.Location) Option {
	return func(r *Registry) {
		for idx := range locations {
			r.locations[locations[idx].Name] = locations[idx]
		}
	}
}

// WithEndpoints registers endpoints to the registry.
func WithEndpoints(endpoints []ginosdkapisv1.Endpoint) Option {
	return func(r *Registry) {
		for idx := range endpoints {
			r.endpoints[fmt.Sprintf(
				"%s/%s", endpoints[idx].Namespace, endpoints[idx].Name,
			)] = endpoints[idx]
		}
	}
}

// WithEndpointGroups registers endpoint groups to the registry.
func WithEndpointGroups(groups []ginosdkapisv1.EndpointGroup) Option {
	return func(r *Registry) {
		for idx := range groups {
			r.groups[groups[idx].Name] = groups[idx]
		}
	}
}
