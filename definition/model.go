package definition

import "time"

// Model is what a model is made of. One Model value per version; versions of the same name share
// one Model CR, with TargetVersion pointing at whichever is current.
type Model struct {
	Name     string
	Version  string
	IsTarget bool
	Inputs   []ModelInput
	Enrich   Enrich
}

// ModelInput is one input a model version wants and how fresh: exactly one of Freshness or At,
// checked by New.
type ModelInput struct {
	Name      string
	Freshness time.Duration
	At        string
	Optional  bool
}
