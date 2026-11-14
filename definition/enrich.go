package definition

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
)

// Errors for enrichment.
var (
	ErrMissingRequiredInput = errors.New("missing required input")
)

// EnrichInputs are the transformed payloads an enrich receives, by input name.
type EnrichInputs map[string]json.RawMessage

// Enrich accepts inputs and returns the unified output.
type Enrich func(
	ctx context.Context,
	l *slog.Logger,
	endpoint ginosdkapisv1.WorkerAssignmentEndpoint,
	inputs EnrichInputs,
) (any, error)

// DecodeEnrichInput returns the message decoded into the type -- if its optional and not here ok
// is false. Only returns an error on unmarshal issues.
func DecodeEnrichInput[T any]( //nolint: ireturn
	in EnrichInputs,
	name string,
) (value T, ok bool, err error) { //nolint: nonamedreturns
	raw, ok := in[name]
	if !ok {
		return value, false, nil
	}

	err = json.Unmarshal(raw, &value)
	if err != nil {
		return value, true, fmt.Errorf("input %q: %w", name, err)
	}

	return value, true, nil
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
