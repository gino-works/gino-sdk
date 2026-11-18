package definition

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
)

// TransformRequest is one part (task) of one input to transform. Inputs can be made of several
// parts (cli commands, http requests, tcp probes); each is transformed on its own, so a transform
// for a multi-part input switches on Task. The conventional output is a partial instance of the
// model(s) the input feeds -- enrich then merges the parts.
type TransformRequest struct {
	// Endpoint is the endpoint the payload came from: tags, interface tags, location -- context
	// for *how to read* the payload (or what to drop early). Decorating the output with endpoint
	// data is enrich's job.
	Endpoint ginosdkapisv1.WorkerAssignmentEndpoint
	// Input is the logical input (ModelInput.name) this part belongs to.
	Input string
	// Task is the part of the input: the CliCommand/HTTPRequest/TCPProbe name.
	Task        string
	CollectedAt time.Time
	Payload     []byte
}

// Transform accepts one part of an input and returns the transformed version of it.
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
