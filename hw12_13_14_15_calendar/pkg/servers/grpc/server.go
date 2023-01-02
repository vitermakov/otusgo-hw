package grpc

import (
	"errors"
	"net"
	"strconv"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers"
	"google.golang.org/grpc"
)

type RegisterHandlerFunc func(s *grpc.Server)

type Server struct {
	*grpc.Server
	config      servers.Config
	Logger      logger.Logger
	AuthService servers.AuthService
}

func NewServer(config servers.Config, authSrv servers.AuthService, logger logger.Logger) *Server {
	return &Server{Server: grpc.NewServer(), config: config, Logger: logger, AuthService: authSrv}
}

func (s *Server) RegisterHandler(handlerFunc RegisterHandlerFunc) {
	if handlerFunc == nil {
		return
	}
	handlerFunc(s.Server)
}

func (s *Server) Start() error {
	s.Logger.Info("GRPC server starting")
	address := net.JoinHostPort(s.config.GetHost(), strconv.Itoa(s.config.GetPort()))
	socket, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	err = s.Serve(socket)
	if err == nil || errors.Is(err, grpc.ErrServerStopped) {
		return nil
	}
	s.Logger.Error("Failed to start GRPC server: %w", err)
	return err
}

func (s *Server) Stop() {
	s.Server.GracefulStop()
	s.Logger.Info("GRPC server stopped")
}
