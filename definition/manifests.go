package definition

import (
	"cmp"
	"fmt"
	"maps"
	"os"
	"slices"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
	apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// Manifests holds all the objects a registry is made of for dumping to disk/a cluster.
type Manifests struct {
	Platforms []ginosdkapisv1.Platform
	Models    []ginosdkapisv1.Model

	Clusters       []ginosdkapisv1.Cluster
	Locations      []ginosdkapisv1.Location
	Endpoints      []ginosdkapisv1.Endpoint
	EndpointGroups []ginosdkapisv1.EndpointGroup
}

// Dump writes all the manifests out to path.
func (m Manifests) Dump(path string) error { //nolint: gocyclo
	for _, platform := range m.Platforms { //nolint: gocritic
		b, err := yaml.Marshal(platform)
		if err != nil {
			return fmt.Errorf("marshaling platform %q: %w", platform.Name, err)
		}

		err = os.WriteFile(
			fmt.Sprintf("%s/platform_%s.yaml", path, platform.Name),
			b,
			0o600, //nolint: mnd
		)
		if err != nil {
			return fmt.Errorf("writing platform %q: %w", platform.Name, err)
		}
	}

	for _, model := range m.Models { //nolint: gocritic
		b, err := yaml.Marshal(model)
		if err != nil {
			return fmt.Errorf("marshaling model %q: %w", model.Name, err)
		}

		err = os.WriteFile(
			fmt.Sprintf("%s/model_%s.yaml", path, model.Name),
			b,
			0o600, //nolint: mnd
		)
		if err != nil {
			return fmt.Errorf("writing model %q: %w", model.Name, err)
		}
	}

	for _, cluster := range m.Clusters { //nolint: gocritic
		b, err := yaml.Marshal(cluster)
		if err != nil {
			return fmt.Errorf("marshaling cluster %q: %w", cluster.Name, err)
		}

		err = os.WriteFile(
			fmt.Sprintf("%s/cluster_%s.yaml", path, cluster.Name),
			b,
			0o600, //nolint: mnd
		)
		if err != nil {
			return fmt.Errorf("writing cluster %q: %w", cluster.Name, err)
		}
	}

	for _, location := range m.Locations { //nolint: gocritic
		b, err := yaml.Marshal(location)
		if err != nil {
			return fmt.Errorf("marshaling location %q: %w", location.Name, err)
		}

		err = os.WriteFile(
			fmt.Sprintf("%s/location_%s.yaml", path, location.Name),
			b,
			0o600, //nolint: mnd
		)
		if err != nil {
			return fmt.Errorf("writing location %q: %w", location.Name, err)
		}
	}

	for _, endpoint := range m.Endpoints { //nolint: gocritic
		b, err := yaml.Marshal(endpoint)
		if err != nil {
			return fmt.Errorf("marshaling endpoint %q: %w", endpoint.Name, err)
		}

		err = os.WriteFile(
			fmt.Sprintf("%s/endpoint_%s_%s.yaml", path, endpoint.Namespace, endpoint.Name),
			b,
			0o600, //nolint: mnd
		)
		if err != nil {
			return fmt.Errorf(
				"writing endpoint '%s/%s': %w",
				endpoint.Namespace,
				endpoint.Name,
				err,
			)
		}
	}

	for _, endpointGroup := range m.EndpointGroups { //nolint: gocritic
		b, err := yaml.Marshal(endpointGroup)
		if err != nil {
			return fmt.Errorf("marshaling endpoint group %q: %w", endpointGroup.Name, err)
		}

		err = os.WriteFile(
			fmt.Sprintf("%s/endpointgroup_%s.yaml", path, endpointGroup.Name),
			b,
			0o600, //nolint: mnd
		)
		if err != nil {
			return fmt.Errorf("writing endpoint group %q: %w", endpointGroup.Name, err)
		}
	}

	return nil
}

// FnName is the transform fn name gino's Platform CR carries for a platform's input, and what
// the transformer plugin dispatches on. Created here so we dont have to type it in multiple places
// and have mismatches.
func FnName(platform, input string) string {
	return platform + "/" + input
}

// Manifests projects every definition to the gino CRs that describe it, in a stable order (by
// name, inputs and versions sorted) so the output diffs cleanly. This is what the manifests cli
// writes; nobody hand-edits these.
func (r *Registry) Manifests() Manifests {
	platforms := make([]ginosdkapisv1.Platform, 0, len(r.platforms))

	for _, name := range slices.Sorted(maps.Keys(r.platforms)) {
		p := r.platforms[name]

		cr := ginosdkapisv1.Platform{
			TypeMeta: apimachineryapismetav1.TypeMeta{
				APIVersion: ginosdkapisv1.SchemeGroupVersion.String(),
				Kind:       "Platform",
			},
			ObjectMeta: apimachineryapismetav1.ObjectMeta{Name: p.Name},
			Spec: ginosdkapisv1.PlatformSpec{
				Match:     p.Match,
				Inputs:    make(map[string]ginosdkapisv1.PlatformInput, len(p.Inputs)),
				Transform: make(map[string][]ginosdkapisv1.PlatformTransform, len(p.Inputs)),
			},
		}

		for _, name := range slices.Sorted(maps.Keys(p.Inputs)) {
			in := p.Inputs[name]
			fn := FnName(p.Name, name)

			cr.Spec.Inputs[name] = in.Impl
			cr.Spec.Transform[name] = []ginosdkapisv1.PlatformTransform{{Fn: &fn}}
		}

		platforms = append(platforms, cr)
	}

	models := make([]ginosdkapisv1.Model, 0, len(r.models))

	for _, name := range slices.Sorted(maps.Keys(r.models)) {
		versions := slices.Clone(r.models[name])

		slices.SortFunc(versions, func(a, b Model) int { return cmp.Compare(a.Version, b.Version) })

		cr := ginosdkapisv1.Model{
			TypeMeta: apimachineryapismetav1.TypeMeta{
				APIVersion: ginosdkapisv1.SchemeGroupVersion.String(),
				Kind:       "Model",
			},
			ObjectMeta: apimachineryapismetav1.ObjectMeta{Name: name},
		}

		for _, m := range versions {
			if m.IsTarget {
				cr.Spec.TargetVersion = m.Version
			}

			mv := ginosdkapisv1.ModelVersion{Version: m.Version}

			for _, in := range m.Inputs {
				mi := ginosdkapisv1.ModelInput{
					Name:     in.Name,
					Optional: in.Optional,
				}

				if in.Freshness > 0 {
					mi.Require.Freshness = &apimachineryapismetav1.Duration{Duration: in.Freshness}
				} else {
					at := in.At
					mi.Require.At = &at
				}

				mv.Inputs = append(mv.Inputs, mi)
			}

			cr.Spec.Versions = append(cr.Spec.Versions, mv)
		}

		models = append(models, cr)
	}

	clusters := make([]ginosdkapisv1.Cluster, 0, len(r.clusters))

	for _, v := range r.clusters { //nolint: gocritic
		clusters = append(clusters, v)
	}

	locations := make([]ginosdkapisv1.Location, 0, len(r.locations))

	for _, v := range r.locations { //nolint: gocritic
		locations = append(locations, v)
	}

	endpoints := make([]ginosdkapisv1.Endpoint, 0, len(r.endpoints))

	for _, v := range r.endpoints { //nolint: gocritic
		endpoints = append(endpoints, v)
	}

	groups := make([]ginosdkapisv1.EndpointGroup, 0, len(r.groups))

	for _, v := range r.groups { //nolint: gocritic
		groups = append(groups, v)
	}

	return Manifests{
		Platforms:      platforms,
		Models:         models,
		Clusters:       clusters,
		Locations:      locations,
		Endpoints:      endpoints,
		EndpointGroups: groups,
	}
}
