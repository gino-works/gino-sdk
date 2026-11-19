package definition

import (
	"time"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
)

// Model is what a model is made of. One Model value per version; versions of the same name share
// one Model CR, with TargetVersion pointing at whichever is current.
type Model struct {
	Name     string
	Version  string
	IsTarget bool
	Inputs   []ModelInput
	// Emit is when the enricher emits a document for this version, see ModelEnrichSpec. Empty
	// leaves it to the api server's default (allRefreshed: one document per full cycle).
	Emit   ginosdkapisv1.EnrichEmit
	Enrich Enrich
	Load   Load
}

// ModelInput is one input a model version wants and how fresh: exactly one of Freshness or At,
// checked by New.
type ModelInput struct {
	Name      string
	Freshness time.Duration
	At        string
}
