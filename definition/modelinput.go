package definition

import ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"

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

// CliPlatformInput is a convenience func to return a cli platform input.
func CliPlatformInput(
	commands []ginosdkapisv1.CliCommand,
	transform Transform,
) PlatformInput {
	return PlatformInput{
		Impl: ginosdkapisv1.PlatformInput{
			Cli: &ginosdkapisv1.CliPlatformInput{
				Commands: commands,
			},
		},
		Transform: transform,
	}
}

// HTTPPlatformInput is a convenience func to return an http platform input.
func HTTPPlatformInput(
	transform Transform,
) PlatformInput {
	return PlatformInput{
		Impl: ginosdkapisv1.PlatformInput{
			HTTP: &ginosdkapisv1.HTTPPlatformInput{},
		},
		Transform: transform,
	}
}
