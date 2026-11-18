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

// EnrichPart is one transformed part (task) of an input.
type EnrichPart struct {
	// Task is the part: the CliCommand/HTTPRequest/TCPProbe name.
	Task        string
	CollectedAt time.Time
	Payload     json.RawMessage
}

// EnrichInput is one input of the model version and all of its parts, in the platform's order.
type EnrichInput struct {
	Name  string
	Parts []EnrichPart
}

// EnrichInputs is everything an enrich receives: every input of the model version, in the model
// version's input order, each with all of its parts in the platform's order. Every input is
// required, so an enrich only runs once all of them -- and all of their parts -- are present.
//
// Convention: when parts are partial models, merge them in this order (inputs in model order,
// parts in platform order), later wins. DecodeEnrichParts hands them over in exactly that order.
type EnrichInputs []EnrichInput

// Input returns the named input.
func (in EnrichInputs) Input(name string) (EnrichInput, bool) {
	for idx := range in {
		if in[idx].Name == name {
			return in[idx], true
		}
	}

	return EnrichInput{}, false
}

// Enrich accepts inputs and returns the unified output.
type Enrich func(
	ctx context.Context,
	l *slog.Logger,
	endpoint ginosdkapisv1.WorkerAssignmentEndpoint,
	inputs EnrichInputs,
) (any, error)

// EnrichedPart is one decoded part.
type EnrichedPart[T any] struct {
	Input       string
	Task        string
	CollectedAt time.Time
	Value       T
}

// DecodeEnrichPart decodes one part of one input into T. A missing input or part is an error --
// gino only enriches once everything is present, so it means the definition and the platform
// disagree.
func DecodeEnrichPart[T any](
	in EnrichInputs,
	input,
	task string,
) (EnrichedPart[T], error) {
	ei, ok := in.Input(input)
	if !ok {
		return EnrichedPart[T]{}, fmt.Errorf("%w: input %q", ErrMissingRequiredInput, input)
	}

	for idx := range ei.Parts {
		if ei.Parts[idx].Task != task {
			continue
		}

		return decodeEnrichPart[T](input, &ei.Parts[idx])
	}

	return EnrichedPart[T]{}, fmt.Errorf(
		"%w: input %q has no part %q",
		ErrMissingRequiredInput,
		input,
		task,
	)
}

// DecodeEnrichParts decodes every part of every input into T, in merge order (inputs in model
// order, parts in platform order).
func DecodeEnrichParts[T any](
	in EnrichInputs,
) ([]EnrichedPart[T], error) {
	count := 0

	for idx := range in {
		count += len(in[idx].Parts)
	}

	out := make([]EnrichedPart[T], 0, count)

	for idx := range in {
		for partIdx := range in[idx].Parts {
			part, err := decodeEnrichPart[T](in[idx].Name, &in[idx].Parts[partIdx])
			if err != nil {
				return nil, err
			}

			out = append(out, part)
		}
	}

	return out, nil
}

func decodeEnrichPart[T any](
	input string,
	part *EnrichPart,
) (EnrichedPart[T], error) {
	out := EnrichedPart[T]{
		Input:       input,
		Task:        part.Task,
		CollectedAt: part.CollectedAt,
	}

	err := json.Unmarshal(part.Payload, &out.Value)
	if err != nil {
		return EnrichedPart[T]{}, fmt.Errorf("input %q part %q: %w", input, part.Task, err)
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
