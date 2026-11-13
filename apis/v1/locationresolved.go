package v1

import "slices"

// LocationResolved is the effective configuration any parents have been applied.
type LocationResolved struct {
	LocationTemplate `json:",inline"`

	// Sources records where each field came from, keyed by field path within the template.
	Sources map[string]ResolvedSource `json:"sources,omitempty"`
	// TagSources records where each tag came from.
	TagSources map[string]ResolvedSource `json:"tagSources,omitempty"`
	// SecretSources records where each namespace secret came from (the nearest location naming it).
	SecretSources map[string]ResolvedSource `json:"secretSources,omitempty"`
}

// NewLocationResolved returns a defaulted LocationResolved.
func NewLocationResolved() LocationResolved {
	return LocationResolved{
		LocationTemplate: LocationTemplate{
			Tags: map[string]string{},
		},
		Sources:       map[string]ResolvedSource{},
		TagSources:    map[string]ResolvedSource{},
		SecretSources: map[string]ResolvedSource{},
	}
}

// NamespaceName returns the namespace the named location manages under this resolved config.
func (r *LocationResolved) NamespaceName(locationName string) string {
	if r.Namespace.Prefix == "" {
		return locationName
	}

	return r.Namespace.Prefix + "-" + locationName
}

// Merge takes any set/non-zero value field from other and puts it into r.
func (r *LocationResolved) Merge(
	source ResolvedSource,
	otherTemplate LocationTemplate,
) {
	for k, v := range otherTemplate.Tags {
		r.Tags[k] = v
		r.TagSources[k] = source
	}

	if otherTemplate.DisableCollection != nil {
		r.DisableCollection = new(*otherTemplate.DisableCollection)
		r.Sources["disableCollection"] = source
	}

	if otherTemplate.Cluster != "" {
		r.Cluster = otherTemplate.Cluster
		r.Sources["cluster"] = source
	}

	if otherTemplate.Type != "" {
		r.Type = otherTemplate.Type
		r.Sources["type"] = source
	}

	if otherTemplate.Namespace.Prefix != "" {
		r.Namespace.Prefix = otherTemplate.Namespace.Prefix
		r.Sources["namespace.prefix"] = source
	}

	for _, name := range otherTemplate.Namespace.Secrets {
		if !slices.Contains(r.Namespace.Secrets, name) {
			r.Namespace.Secrets = append(r.Namespace.Secrets, name)
		}

		r.SecretSources[name] = source
	}
}
