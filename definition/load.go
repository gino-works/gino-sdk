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

// LoadInput is a single document for the loader server to load.
type LoadInput struct {
	Endpoint   ginosdkapisv1.WorkerAssignmentEndpoint
	EnrichedAt time.Time
	// Inputs is when each input that went into this document was collected, by input name. For
	// a multi-part input that is its oldest part. For a time-series row this is the measurement
	// time; EnrichedAt is only when the document was assembled.
	Inputs  map[string]time.Time
	Payload json.RawMessage
}

// TypedLoadInput is a LoadInput with Payload decoded.
type TypedLoadInput[T any] struct {
	Endpoint   ginosdkapisv1.WorkerAssignmentEndpoint
	EnrichedAt time.Time
	Inputs     map[string]time.Time
	Payload    T
}

// CollectedAt returns the collection time of the named input, or false if that input did not go
// into this document.
func (t TypedLoadInput[T]) CollectedAt(input string) (time.Time, bool) {
	at, ok := t.Inputs[input]

	return at, ok
}

// OldestCollectedAt returns the earliest collection time across the inputs that went into this
// document -- the conservative "how stale is this" answer. Zero if there are none.
func (t TypedLoadInput[T]) OldestCollectedAt() time.Time {
	var oldest time.Time

	for _, at := range t.Inputs {
		if oldest.IsZero() || at.Before(oldest) {
			oldest = at
		}
	}

	return oldest
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
			Inputs:     in[idx].Inputs,
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
