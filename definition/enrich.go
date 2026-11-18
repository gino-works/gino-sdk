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

// EnrichTask is the transformed payload of one task of an input.
type EnrichTask struct {
	// Name is the task's name (CliCommand/HTTPRequest/TCPProbe.name).
	Name        string
	CollectedAt time.Time
	Payload     json.RawMessage
}

// EnrichInput is one input of the model version and all of its tasks, in the platform's order.
type EnrichInput struct {
	Name  string
	Tasks []EnrichTask
}

// EnrichInputs is everything an enrich receives: every input of the model version, in the model
// version's input order, each with all of its tasks in the platform's order. Every input is
// required, so an enrich only runs once all of them -- and all of their tasks -- are present.
//
// Convention: when transforms emit partial models, merge them in this order (inputs in model
// order, tasks in platform order), later wins. DecodeEnrichTasks hands them over in exactly that
// order.
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

// EnrichedTask is one decoded task payload.
type EnrichedTask[T any] struct {
	ModelInput  string
	Task        string
	CollectedAt time.Time
	Value       T
}

// DecodeEnrichTask decodes one task of one input into T. A missing input or task is an error --
// gino only enriches once everything is present, so it means the definition and the platform
// disagree.
func DecodeEnrichTask[T any](
	in EnrichInputs,
	input, task string,
) (EnrichedTask[T], error) {
	ei, ok := in.Input(input)
	if !ok {
		return EnrichedTask[T]{}, fmt.Errorf("%w: input %q", ErrMissingRequiredInput, input)
	}

	for idx := range ei.Tasks {
		if ei.Tasks[idx].Name != task {
			continue
		}

		return decodeEnrichTask[T](input, &ei.Tasks[idx])
	}

	return EnrichedTask[T]{}, fmt.Errorf(
		"%w: input %q has no task %q",
		ErrMissingRequiredInput,
		input,
		task,
	)
}

// DecodeEnrichTasks decodes every task of every input into T, in merge order (inputs in model
// order, tasks in platform order). This is the "every transform emits a partial model" case:
// decode them all as the model type and fold them together in the order returned.
func DecodeEnrichTasks[T any](in EnrichInputs) ([]EnrichedTask[T], error) {
	count := 0

	for idx := range in {
		count += len(in[idx].Tasks)
	}

	out := make([]EnrichedTask[T], 0, count)

	for idx := range in {
		for taskIdx := range in[idx].Tasks {
			task, err := decodeEnrichTask[T](in[idx].Name, &in[idx].Tasks[taskIdx])
			if err != nil {
				return nil, err
			}

			out = append(out, task)
		}
	}

	return out, nil
}

func decodeEnrichTask[T any](modelInput string, task *EnrichTask) (EnrichedTask[T], error) {
	out := EnrichedTask[T]{
		ModelInput:  modelInput,
		Task:        task.Name,
		CollectedAt: task.CollectedAt,
	}

	err := json.Unmarshal(task.Payload, &out.Value)
	if err != nil {
		return EnrichedTask[T]{}, fmt.Errorf(
			"modelInput %q task %q: %w",
			modelInput,
			task.Name,
			err,
		)
	}

	return out, nil
}

// Enrich runs the registered enrich for model/version.
func (r *Registry) Enrich(
	ctx context.Context,
	l *slog.Logger,
	model,
	version string,
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
