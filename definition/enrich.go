package definition

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
)

// Errors for enrichment.
var (
	ErrMissingRequiredInput = errors.New("missing required input")
)

// EnrichInput holds the input and the collected at time of an input (by name).
type EnrichInput struct {
	CollectedAt time.Time
	Payload     json.RawMessage
}

// EnrichInputs are the transformed payloads an enrich receives, by input name.
type EnrichInputs map[string]EnrichInput

// Enrich accepts inputs and returns the unified output.
type Enrich func(
	ctx context.Context,
	l *slog.Logger,
	endpoint ginosdkapisv1.WorkerAssignmentEndpoint,
	inputs EnrichInputs,
) (any, error)

// EnrichedInput is one decoded input. Found is false only for an optional input that has not
// been collected yet -- the enricher does not run until every required input is present.
type EnrichedInput[T any] struct {
	Value       T
	CollectedAt time.Time
	Found       bool
}

// DecodeEnrichInput returns the message decoded into the type -- if its optional and not here ok
// is false. Only returns an error on unmarshal issues.
func DecodeEnrichInput[T any](
	in EnrichInputs,
	name string,
) (EnrichedInput[T], error) {
	raw, ok := in[name]
	if !ok {
		return EnrichedInput[T]{}, nil
	}

	out := EnrichedInput[T]{
		CollectedAt: raw.CollectedAt,
		Found:       true,
	}

	err := json.Unmarshal(raw.Payload, &out.Value)
	if err != nil {
		return EnrichedInput[T]{}, fmt.Errorf("input %q: %w", name, err)
	}

	return out, nil
}

// Enrich runs the registered enrich for model/version.
func (r *Registry) Enrich(
	ctx context.Context,
	l *slog.Logger,
	model, version string,
	endpoint ginosdkapisv1.WorkerAssignmentEndpoint,
	inputs EnrichInputs,
) (any, error) {
	for _, m := range r.models[model] {
		if m.Version != version {
			continue
		}

		value, err := m.Enrich(ctx, l, endpoint, inputs)
		if err != nil {
			return nil, fmt.Errorf("enriching %s/%s: %w", model, version, err)
		}

		return value, nil
	}

	return nil, fmt.Errorf("%w: %s/%s", ErrUnknownModel, model, version)
}
