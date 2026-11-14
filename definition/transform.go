package definition

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// Transform accepts the payload and returns the transformed version of it.
type Transform func(
	ctx context.Context,
	l *slog.Logger,
	payload []byte,
) (any, error)

// Transform runs the registered transform for platform/input. This is what a transformer plugin's
// Transform RPC calls; the fn name on the wire is FnName(platform, input).
func (r *Registry) Transform(
	ctx context.Context,
	l *slog.Logger,
	fn string,
	payload []byte,
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

	value, err := in.Transform(ctx, l, payload)
	if err != nil {
		return nil, fmt.Errorf("transforming %s/%s: %w", platform, input, err)
	}

	return value, nil
}
