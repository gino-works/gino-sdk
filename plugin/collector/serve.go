package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
	ginosdkconstants "github.com/gino-works/gino-sdk/constants"
	ginosdkpluginv1 "github.com/gino-works/gino-sdk/plugin/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func decodeEndpoint(b []byte) (*ginosdkapisv1.WorkerAssignmentEndpoint, string, error) {
	endpoint := &ginosdkapisv1.WorkerAssignmentEndpoint{}

	err := json.Unmarshal(b, endpoint)
	if err != nil {
		return nil, "", status.Errorf(codes.InvalidArgument, "decoding endpoint: %v", err)
	}

	if endpoint.Name == "" || endpoint.Namespace == "" {
		return nil, "", status.Error(codes.InvalidArgument, "endpoint has no namespace/name")
	}

	return endpoint, endpoint.Namespace + "/" + endpoint.Name, nil
}

// NewWorker returns a Worker that the server can manage.
type NewWorker func(
	processCtx context.Context,
	l *slog.Logger,
	endpoint *ginosdkapisv1.WorkerAssignmentEndpoint,
	creds map[string][]byte,
) (Worker, error)

// Worker worker is the worker implementation the plugin server can manage.
type Worker interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	Collect(
		ctx context.Context,
		task *ginosdkpluginv1.CollectTask,
	) ([]*ginosdkpluginv1.Collected, error)
}

// Run runs the collector server.
func Run(
	ctx context.Context,
	l *slog.Logger,
	version,
	name string,
	newWorker NewWorker,
) error {
	s := &server{
		l:           l,
		ctx:         ctx,
		version:     version,
		name:        name,
		newWorker:   newWorker,
		workers:     map[string]Worker{},
		workersLock: sync.Mutex{},
	}

	lc := net.ListenConfig{}

	lis, err := lc.Listen(ctx, "unix", ginosdkconstants.GinoPluginSocketPath)
	if err != nil {
		return fmt.Errorf("listening on %q: %w", ginosdkconstants.GinoPluginSocketPath, err)
	}

	_ = os.Chmod(ginosdkconstants.GinoPluginSocketPath, ginosdkconstants.GinoPluginSocketPerms)

	srv := grpc.NewServer()
	ginosdkpluginv1.RegisterPluginServiceServer(srv, s)
	ginosdkpluginv1.RegisterCollectorServiceServer(srv, s)

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	return srv.Serve(lis)
}

type server struct {
	ginosdkpluginv1.UnimplementedPluginServiceServer
	ginosdkpluginv1.UnimplementedCollectorServiceServer

	// the root ctx, for workers so it lives past an rpc.
	ctx context.Context //nolint:containedctx

	l *slog.Logger

	version   string
	name      string
	newWorker NewWorker

	workers     map[string]Worker
	workersLock sync.Mutex
}

// Ping satisfies the plugin service server initialization.
func (s *server) Ping(
	_ context.Context,
	_ *ginosdkpluginv1.PingRequest,
) (*ginosdkpluginv1.PingResponse, error) {
	return &ginosdkpluginv1.PingResponse{
		Version:       ginosdkconstants.GinoPluginContractVersion,
		PluginName:    s.name,
		PluginVersion: s.version,
		Services: []ginosdkpluginv1.Service{
			ginosdkpluginv1.Service_SERVICE_COLLECTOR,
		},
	}, nil
}

func (s *server) Shutdown(
	ctx context.Context,
	_ *ginosdkpluginv1.ShutdownRequest,
) (*ginosdkpluginv1.ShutdownResponse, error) {
	s.workersLock.Lock()
	defer s.workersLock.Unlock()

	wg := sync.WaitGroup{}

	for endpoint := range s.workers {
		wg.Go(
			func() {
				_ = s.workers[endpoint].Close(ctx)
			},
		)
	}

	wg.Wait()

	for endpoint := range s.workers {
		delete(s.workers, endpoint)
	}

	return &ginosdkpluginv1.ShutdownResponse{}, nil
}

// Open opens a collector worker for the given request data.
func (s *server) Open(
	ctx context.Context,
	req *ginosdkpluginv1.OpenRequest,
) (*ginosdkpluginv1.OpenResponse, error) {
	endpoint, k, err := decodeEndpoint(req.GetEndpoint())
	if err != nil {
		return nil, err
	}

	s.workersLock.Lock()

	_, ok := s.workers[k]
	if ok {
		_ = s.workers[k].Close(ctx)
	}

	delete(s.workers, k)

	w, err := s.newWorker( //nolint: contextcheck
		// rpc ctx ends with the rpc; the worker lives with the server
		s.ctx,
		s.l.WithGroup(k),
		endpoint,
		req.GetCredentials(),
	)
	if err == nil {
		s.workers[k] = w
	}

	s.workersLock.Unlock()

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "creating worker: %v", err)
	}

	err = w.Open(ctx)
	if err != nil {
		s.workersLock.Lock()

		delete(s.workers, k)

		s.workersLock.Unlock()

		return nil, status.Errorf(codes.Unavailable, "opening session: %v", err)
	}

	return &ginosdkpluginv1.OpenResponse{}, nil
}

func (s *server) Collect(
	ctx context.Context,
	req *ginosdkpluginv1.CollectRequest,
) (*ginosdkpluginv1.CollectResponse, error) {
	_, k, err := decodeEndpoint(req.GetEndpoint())
	if err != nil {
		return nil, err
	}

	s.workersLock.Lock()

	w, ok := s.workers[k]

	s.workersLock.Unlock()

	if !ok {
		return nil, status.Errorf(codes.Unavailable, "worker for endpoint does not exist")
	}

	task := req.GetTask()
	if task == nil || len(task.GetImpl()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "collect request has no task")
	}

	collected, err := w.Collect(ctx, task)
	if err != nil {
		return nil, err
	}

	return &ginosdkpluginv1.CollectResponse{
		Collected: collected,
	}, nil
}

func (s *server) Close(
	ctx context.Context,
	req *ginosdkpluginv1.CloseRequest,
) (*ginosdkpluginv1.CloseResponse, error) {
	_, k, err := decodeEndpoint(req.GetEndpoint())
	if err != nil {
		return nil, err
	}

	s.workersLock.Lock()

	w, ok := s.workers[k]
	if ok {
		delete(s.workers, k)
	}

	s.workersLock.Unlock()

	if !ok {
		return &ginosdkpluginv1.CloseResponse{}, nil
	}

	_ = w.Close(ctx)

	return &ginosdkpluginv1.CloseResponse{}, nil
}
