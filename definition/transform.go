package definition

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
)

// TransformRequest is the payload of one task of one input to transform. An input can be made of
// several tasks (cli commands, http requests, tcp probes); each is transformed on its own, so a
// transform for such an input switches on Task. The conventional output is a partial instance of
// the model(s) the input feeds -- enrich then merges the tasks' outputs.
type TransformRequest struct {
	// Endpoint is the endpoint the payload came from: tags, interface tags, location -- context
	// for *how to read* the payload (or what to drop early). Decorating the output with endpoint
	// data is enrich's job.
	Endpoint ginosdkapisv1.WorkerAssignmentEndpoint
	// Input is the input (ModelInput.name) the task belongs to.
	Input string
	// Task is the task's name: CliCommand/HTTPRequest/TCPProbe.name.
	Task        string
	CollectedAt time.Time
	Payload     []byte
}

// Transform accepts one task's payload and returns the transformed version of it.
type Transform func(
	ctx context.Context,
	l *slog.Logger,
	req TransformRequest,
) (any, error)

// Transform runs the registered transform for platform/input. This is what a transformer plugin's
// Transform RPC calls; the fn name on the wire is FnName(platform, input).
func (r *Registry) Transform(
	ctx context.Context,
	l *slog.Logger,
	fn string,
	req TransformRequest,
) (any, error) {
	platform, input, ok := strings.Cut(fn, "/")
	if !ok {
		return nil, fmt.Errorf("%w: fn %q is not <platform>/<input>", ErrUnknownInput, fn)
	}

	p, ok := r.platforms[platform]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownPlatform, platform)
	}

	in, ok := p.Inputs[input]
	if !ok {
		return nil, fmt.Errorf("%w: %q on platform %q", ErrUnknownInput, input, platform)
	}

	if req.Input == "" {
		req.Input = input
	}

	value, err := in.Transform(ctx, l, req)
	if err != nil {
		return nil, fmt.Errorf("transforming %s/%s (%s): %w", platform, input, req.Task, err)
	}

	return value, nil
}
