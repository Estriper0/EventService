package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Estriper0/EventService/internal/config"
	event_handler "github.com/Estriper0/EventService/internal/handlers/event"
	"github.com/Estriper0/EventService/internal/service"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	logger     *slog.Logger
	config     *config.Config
	grpcServer *grpc.Server
}

func New(
	logger *slog.Logger,
	config *config.Config,
	eventService service.IEventService,
) *GRPCServer {
	grpcServer := grpc.NewServer()

	event_handler.Register(grpcServer, eventService)

	return &GRPCServer{
		logger:     logger,
		config:     config,
		grpcServer: grpcServer,
	}
}

func (s *GRPCServer) Run() {
	s.logger.Info(
		"Starting gRPC server",
		slog.Int("port", s.config.Port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		panic(err)
	}

	s.logger.Info(
		"GRPC server is running",
		slog.String("addr", l.Addr().String()),
	)
	if err := s.grpcServer.Serve(l); err != nil {
		panic(err)
	}
}

func (s *GRPCServer) Stop() {
	s.logger.Info(
		"Stopping grpc server",
		slog.Int("port", s.config.Port),
	)

	s.grpcServer.GracefulStop()
}
