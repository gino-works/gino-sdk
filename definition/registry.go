package definition

import (
	"errors"
	"fmt"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
)

// Errors a registry can report.
var (
	ErrUnknownPlatform = errors.New("unknown platform")
	ErrUnknownInput    = errors.New("unknown input")
	ErrUnknownModel    = errors.New("unknown model")
	ErrDuplicate       = errors.New("duplicate definition")
	ErrInvalid         = errors.New("invalid definition")
)

// Registry holds every definition a plugin binary knows, and is both the manifest source and the
// dispatch table for transform/enrich/load functionality.
type Registry struct {
	// the actual bits we need in the registry to plumb things together.
	platforms map[string]Platform
	models    map[string][]Model // all versions

	// optional other stuff we can store/dump to yaml so its all in the same place as the other
	// things.
	clusters  map[string]ginosdkapisv1.Cluster
	locations map[string]ginosdkapisv1.Location
	endpoints map[string]ginosdkapisv1.Endpoint
	groups    map[string]ginosdkapisv1.EndpointGroup

	// TODO - pipelines but not sure if i want them down here or more "required".
}

// New builds a Registry.
func New(
	platforms []Platform,
	models []Model,
	options ...Option,
) (*Registry, error) {
	r := &Registry{
		platforms: make(map[string]Platform, len(platforms)),
		models:    map[string][]Model{},
		clusters:  map[string]ginosdkapisv1.Cluster{},
		locations: map[string]ginosdkapisv1.Location{},
		endpoints: map[string]ginosdkapisv1.Endpoint{},
		groups:    map[string]ginosdkapisv1.EndpointGroup{},
	}

	for _, opt := range options {
		opt(r)
	}

	for _, p := range platforms {
		if _, dup := r.platforms[p.Name]; dup {
			return nil, fmt.Errorf("%w: platform %q", ErrDuplicate, p.Name)
		}

		r.platforms[p.Name] = p
	}

	for _, m := range models {
		for _, existing := range r.models[m.Name] {
			if existing.Version == m.Version {
				return nil, fmt.Errorf("%w: model %s/%s", ErrDuplicate, m.Name, m.Version)
			}
		}

		for _, in := range m.Inputs {
			if (in.Freshness > 0) == (in.At != "") {
				return nil, fmt.Errorf(
					"%w: model %s/%s input %q must set exactly one of freshness or at",
					ErrInvalid,
					m.Name,
					m.Version,
					in.Name,
				)
			}
		}

		r.models[m.Name] = append(r.models[m.Name], m)
	}

	for name, versions := range r.models {
		targets := 0

		for _, m := range versions {
			if m.IsTarget {
				targets++
			}
		}

		if targets != 1 {
			return nil, fmt.Errorf(
				"model %q must have exactly one target version, has %d",
				name,
				targets,
			)
		}
	}

	return r, nil
}
