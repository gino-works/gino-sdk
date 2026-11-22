package server

import (
	"context"
	"fmt"
	"net"
	"os"

	ginosdkconstants "github.com/gino-works/gino-sdk/constants"
	"google.golang.org/grpc"
)

// Serve is a convenience func to setup the plugin listener and register services and run it all.
func Serve(
	ctx context.Context,
	registerF func(srv *grpc.Server),
) error {
	lc := net.ListenConfig{}

	lis, err := lc.Listen(ctx, "unix", ginosdkconstants.GinoPluginSocketPath)
	if err != nil {
		return fmt.Errorf("listening on %q: %w", ginosdkconstants.GinoPluginSocketPath, err)
	}

	_ = os.Chmod(ginosdkconstants.GinoPluginSocketPath, ginosdkconstants.GinoPluginSocketPerms)

	srv := grpc.NewServer()
	registerF(srv)

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	return srv.Serve(lis)
}
