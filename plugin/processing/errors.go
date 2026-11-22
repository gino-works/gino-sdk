package server

import (
	"errors"

	ginosdkdefinition "github.com/gino-works/gino-sdk/definition"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StatusFor dispatches the correct status code for the input error.
func StatusFor(err error) error {
	switch {
	case errors.Is(err, ginosdkdefinition.ErrUnknownPlatform),
		errors.Is(err, ginosdkdefinition.ErrUnknownInput),
		errors.Is(err, ginosdkdefinition.ErrUnknownModel):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
