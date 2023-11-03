package internalgrpc

import (
	"log/slog"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/pb"
	internalserver "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/server"
)

type server struct {
	pb.UnimplementedEventServiceServer
	grpcSrv *grpc.Server
	logger  *slog.Logger
	cfg     *config.GRPCServerConf
	app     app.App
}

func NewServer(
	logger *slog.Logger,
	app app.App,
	cfg *config.GRPCServerConf,
	baseSrv *grpc.Server,
) internalserver.Server {
	if baseSrv == nil {
		baseSrv = grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				newLoggerInterceptor(logger).UnaryLoggingInterceptor,
				UnaryAuthInterceptor,
			),
		)
	}

	srv := &server{
		grpcSrv: baseSrv,
		logger:  logger,
		cfg:     cfg,
		app:     app,
	}
	pb.RegisterEventServiceServer(srv.grpcSrv, srv)

	return srv
}

func (s *server) GetType() string {
	return "gRPC"
}

func (s *server) GetAddress() string {
	return s.cfg.Address
}

func (s *server) Start(context.Context) error {
	lsn, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		return err
	}

	return s.grpcSrv.Serve(lsn)
}

func (s *server) Stop(context.Context) error {
	s.grpcSrv.GracefulStop()
	return nil
}
