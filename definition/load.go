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

// Errors for loading.
var (
	ErrSinkUnavailable = errors.New("sink target is unavailable")
	ErrInvalidInput    = errors.New("input cant be decoded")
)

// LoadInput is a single message for the loader server to load.
type LoadInput struct {
	Endpoint   ginosdkapisv1.WorkerAssignmentEndpoint
	EnrichedAt time.Time
	Payload    json.RawMessage
}

// TypedLoadInput is a LoadInput with Data decoded.
type TypedLoadInput[T any] struct {
	Endpoint   ginosdkapisv1.WorkerAssignmentEndpoint
	EnrichedAt time.Time
	Payload    T
}

// LoadInputs is a list of LoadInput -- the loader client (gino) will send load inputs in batches
// like this (as configured by the model).
type LoadInputs []LoadInput

// Load accepts finalized models for the loader server to load.
type Load func(
	ctx context.Context,
	l *slog.Logger,
	inputs LoadInputs,
) error

// DecodeLoadInputs returns the load inputs but typed so the loader server can do whatever it wants
// with that information.
func DecodeLoadInputs[T any](
	in LoadInputs,
) ([]TypedLoadInput[T], error) {
	out := make([]TypedLoadInput[T], len(in))

	for idx := range in {
		var value T

		err := json.Unmarshal(in[idx].Payload, &value)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: input %d (%s/%s): %w",
				ErrInvalidInput,
				idx,
				in[idx].Endpoint.Namespace,
				in[idx].Endpoint.Name,
				err,
			)
		}

		out[idx] = TypedLoadInput[T]{
			Endpoint:   in[idx].Endpoint,
			EnrichedAt: in[idx].EnrichedAt,
			Payload:    value,
		}
	}

	return out, nil
}

// Load runs the registered load op(s) for model/version.
func (r *Registry) Load(
	ctx context.Context,
	l *slog.Logger,
	model, version string,
	inputs LoadInputs,
) error {
	for _, m := range r.models[model] {
		if m.Version != version {
			continue
		}

		err := m.Load(ctx, l, inputs)
		if err != nil {
			return fmt.Errorf("loading %s/%s: %w", model, version, err)
		}

		return nil
	}

	return fmt.Errorf("%w: %s/%s", ErrUnknownModel, model, version)
}
