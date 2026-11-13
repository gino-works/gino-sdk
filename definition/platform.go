package definition

import ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"

// Platform is one vendor/platform/profile's worth of knowledge: how to ask for each input and how
// to read the answer. Projected to a gino Platform CR. Name is the name of the actual CR resource
// that will get yeeted out, and also is the prefix for function names associated with the platform.
type Platform struct {
	Name   string
	Match  ginosdkapisv1.PlatformMatch
	Inputs map[string]PlatformInput
}

// PlatformInput holds the implementation that actually collects/generates the input and the
// function that does the transforming for it.
type PlatformInput struct {
	Impl      ginosdkapisv1.PlatformInput
	Transform Transform
}
