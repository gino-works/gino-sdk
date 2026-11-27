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
	// Emit is when the enricher emits a document for this version, see ModelEnrichEmitSpec. The
	// zero value leaves it to the api server's default (allRefreshed: one document per full
	// cycle); to emit on particular inputs set Type to EnrichEmitInputs and list them.
	Emit   ginosdkapisv1.ModelEnrichEmitSpec
	Enrich Enrich
	Load   Load
}

// ModelInput is one input a model version wants and how fresh: exactly one of Freshness or At,
// checked by New. Splay (only with At) spreads the collections across endpoints, see
// ginosdkapisv1.ModelRequirement.
type ModelInput struct {
	Name      string
	Freshness time.Duration
	At        string
	Splay     time.Duration
}
