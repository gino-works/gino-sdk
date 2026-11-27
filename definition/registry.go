package definition

import (
	"errors"
	"fmt"
	"slices"
	"time"

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
	pipelines map[string]ginosdkapisv1.Pipeline
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
		pipelines: map[string]ginosdkapisv1.Pipeline{},
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

			err := validateSplay(in)
			if err != nil {
				return nil, fmt.Errorf(
					"%w: model %s/%s input %q: %w",
					ErrInvalid,
					m.Name,
					m.Version,
					in.Name,
					err,
				)
			}
		}

		err := validateEmit(&m)
		if err != nil {
			return nil, fmt.Errorf("%w: model %s/%s: %w", ErrInvalid, m.Name, m.Version, err)
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

// maxSplay is the most an input's collections are spread; the crd's limit too.
const maxSplay = 24 * time.Hour

// validateSplay checks in's splay: only with at, and within (0, maxSplay]. The cron's own interval
// is the collector's to cap at -- the sdk does not parse crons.
func validateSplay(in ModelInput) error {
	switch {
	case in.Splay == 0:
		return nil
	case in.At == "":
		return errors.New("splay is only allowed with at")
	case in.Splay < 0 || in.Splay > maxSplay:
		return fmt.Errorf("splay must be greater than 0 and at most %s, is %s", maxSplay, in.Splay)
	default:
		return nil
	}
}

func validateEmit(m *Model) error {
	switch m.Emit.Type {
	case "",
		ginosdkapisv1.EnrichEmitAllRefreshed,
		ginosdkapisv1.EnrichEmitAnyRefreshed:
		if len(m.Emit.ModelInputs) > 0 {
			return fmt.Errorf(
				"emit inputs are only allowed with emit type %q",
				ginosdkapisv1.EnrichEmitInputs,
			)
		}

		return nil
	case ginosdkapisv1.EnrichEmitPassthrough:
		if len(m.Emit.ModelInputs) > 0 {
			return fmt.Errorf(
				"emit inputs are only allowed with emit type %q",
				ginosdkapisv1.EnrichEmitInputs,
			)
		}

		// each arrival goes straight out as the document: there is nothing to assemble it with
		if len(m.Inputs) != 1 {
			return fmt.Errorf(
				"emit type %q takes exactly one input, has %d",
				ginosdkapisv1.EnrichEmitPassthrough,
				len(m.Inputs),
			)
		}

		return nil
	case ginosdkapisv1.EnrichEmitInputs:
		if len(m.Emit.ModelInputs) == 0 {
			return fmt.Errorf(
				"emit type %q needs at least one input",
				ginosdkapisv1.EnrichEmitInputs,
			)
		}
	default:
		return fmt.Errorf(
			"emit type %q is not one of %q, %q, %q",
			m.Emit.Type,
			ginosdkapisv1.EnrichEmitAllRefreshed,
			ginosdkapisv1.EnrichEmitAnyRefreshed,
			ginosdkapisv1.EnrichEmitInputs,
		)
	}

	seen := make(map[string]bool, len(m.Emit.ModelInputs))

	for _, name := range m.Emit.ModelInputs {
		if seen[name] {
			return fmt.Errorf("emit input %q is listed more than once", name)
		}

		seen[name] = true

		if !slices.ContainsFunc(m.Inputs, func(in ModelInput) bool { return in.Name == name }) {
			return fmt.Errorf("emit input %q is not an input of the model", name)
		}
	}

	return nil
}
