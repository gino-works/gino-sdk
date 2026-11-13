package v1

// EndpointResolved is the effective configuration after groups and location have been applied.
type EndpointResolved struct {
	EndpointTemplate `json:",inline"`

	// Sources records where each field came from, keyed by field name path within the template.
	Sources map[string]ResolvedSource `json:"sources,omitempty"`
	// TagSources records where each tag came from.
	TagSources map[string]ResolvedSource `json:"tagSources,omitempty"`
	// InterfaceTagSources records where each tag came from. Indexes align w/ the interface tags
	// in the resolved template.
	InterfaceTagSources []ResolvedSource `json:"interfaceTagSources,omitempty"`
	// CredentialSources records where each kind's credential came from.
	CredentialSources map[CollectorKind]ResolvedSource `json:"credentialSources,omitempty"`
}

// NewEndpointResolved returns a defaulted EndpointResolved.
func NewEndpointResolved() EndpointResolved {
	return EndpointResolved{
		EndpointTemplate: EndpointTemplate{
			Profile:           "default",
			Tags:              map[string]string{},
			Credentials:       map[CollectorKind]string{},
			DisableCollection: new(false),
		},
		Sources: map[string]ResolvedSource{
			"profile": {
				Kind: ResolvedSourceKindDefault,
			},
			"disableCollection": {
				Kind: ResolvedSourceKindDefault,
			},
		},
		TagSources:        map[string]ResolvedSource{},
		CredentialSources: map[CollectorKind]ResolvedSource{},
	}
}

// MergeLocation takes any set/non-zero value field from other and resolves it into r.
func (r *EndpointResolved) MergeLocation(
	other LocationResolved,
) {
	for k, v := range other.Tags {
		r.Tags[k] = v
		r.TagSources[k] = other.TagSources[k]
	}

	if other.DisableCollection != nil {
		r.DisableCollection = new(*other.DisableCollection)
		r.Sources["disableCollection"] = other.Sources["disableCollection"]
	}

	if other.Cluster != "" {
		r.Cluster = other.Cluster
		r.Sources["cluster"] = other.Sources["cluster"]
	}
}

// Merge takes any set/non-zero value field from other and resolves it onto r.
func (r *EndpointResolved) Merge(
	source ResolvedSource,
	otherTemplate EndpointTemplate,
) {
	if otherTemplate.Vendor != "" {
		r.Vendor = otherTemplate.Vendor
		r.Sources["vendor"] = source
	}

	if otherTemplate.Platform != "" {
		r.Platform = otherTemplate.Platform
		r.Sources["platform"] = source
	}

	if otherTemplate.Profile != "" {
		r.Profile = otherTemplate.Profile
		r.Sources["profile"] = source
	}

	for k, v := range otherTemplate.Tags {
		r.Tags[k] = v
		r.TagSources[k] = source
	}

	for _, interfaceTag := range otherTemplate.InterfaceTags {
		r.InterfaceTags = append(r.InterfaceTags, interfaceTag)
		r.InterfaceTagSources = append(r.InterfaceTagSources, source)
	}

	if otherTemplate.DisableCollection != nil {
		r.DisableCollection = new(*otherTemplate.DisableCollection)
		r.Sources["disableCollection"] = source
	}

	if otherTemplate.Cluster != "" {
		r.Cluster = otherTemplate.Cluster
		r.Sources["cluster"] = source
	}

	for kind, name := range otherTemplate.Credentials {
		r.Credentials[kind] = name
		r.CredentialSources[kind] = source
	}
}
