package definition

import (
	"time"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
	apimachineryapismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CliOption sets an input level option on a cli platform input: the default for every command
// of the input that does not set its own.
type CliOption func(o *ginosdkapisv1.CliOptions)

// WithCliTimeout sets the input's default command timeout.
func WithCliTimeout(d time.Duration) CliOption {
	return func(o *ginosdkapisv1.CliOptions) {
		o.Timeout = &apimachineryapismetav1.Duration{Duration: d}
	}
}

// HTTPOption sets an input level option on an http platform input: the default for every
// request of the input that does not set its own.
type HTTPOption func(o *ginosdkapisv1.HTTPOptions)

// WithHTTPTimeout sets the input's default request timeout.
func WithHTTPTimeout(d time.Duration) HTTPOption {
	return func(o *ginosdkapisv1.HTTPOptions) {
		o.Timeout = &apimachineryapismetav1.Duration{Duration: d}
	}
}

// TCPPlatformInput is a convenience func to return a tcp platform input.
func TCPPlatformInput(
	probes []ginosdkapisv1.TCPProbe,
	transform Transform,
) PlatformInput {
	return PlatformInput{
		Impl: ginosdkapisv1.PlatformInput{
			TCP: &ginosdkapisv1.TCPPlatformInput{
				Probes: probes,
			},
		},
		Transform: transform,
	}
}

// CliPlatformInput is a convenience func to return a cli platform input. Options set the
// input level defaults; a command's own options win over them.
func CliPlatformInput(
	commands []ginosdkapisv1.CliCommand,
	transform Transform,
	opts ...CliOption,
) PlatformInput {
	in := &ginosdkapisv1.CliPlatformInput{
		Commands: commands,
	}

	for _, opt := range opts {
		opt(&in.CliOptions)
	}

	return PlatformInput{
		Impl: ginosdkapisv1.PlatformInput{
			Cli: in,
		},
		Transform: transform,
	}
}

// HTTPPlatformInput is a convenience func to return an http platform input. Options set the
// input level defaults; a request's own options win over them.
func HTTPPlatformInput(
	requests []ginosdkapisv1.HTTPRequest,
	transform Transform,
	opts ...HTTPOption,
) PlatformInput {
	in := &ginosdkapisv1.HTTPPlatformInput{
		Requests: requests,
	}

	for _, opt := range opts {
		opt(&in.HTTPOptions)
	}

	return PlatformInput{
		Impl: ginosdkapisv1.PlatformInput{
			HTTP: in,
		},
		Transform: transform,
	}
}
