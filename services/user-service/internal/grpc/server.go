package grpc

import (
	"fmt"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/repository"
	userpb "github.com/squ1ky/avigo-c2c-marketplace/services/user-service/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(cfg config.GRPCConfig, userRepo repository.UserRepository) (*Server, error) {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", cfg.Address, err)
	}

	kaParams := keepalive.ServerParameters{
		MaxConnectionIdle:     cfg.MaxConnectionIdle,
		MaxConnectionAge:      cfg.MaxConnectionAge,
		MaxConnectionAgeGrace: cfg.MaxConnectionAgeGrace,
		Time:                  cfg.Time,
		Timeout:               cfg.Timeout,
	}

	kaEnforcementPolicy := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(kaParams),
		grpc.KeepaliveEnforcementPolicy(kaEnforcementPolicy),
		grpc.MaxRecvMsgSize(1024*1024*10), // 10MB
		grpc.MaxSendMsgSize(1024*1024*10), // 10MB
	)

	userServiceServer := NewUserServiceServer(userRepo)
	userpb.RegisterUserServiceServer(grpcServer, userServiceServer)

	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		listener:   listener,
	}, nil
}

func (s *Server) Start() error {
	log.Printf("gRPC server listening on %s", s.listener.Addr().String())
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) GracefulStop() {
	log.Println("Gracefully stopping gRPC server...")
	s.grpcServer.GracefulStop()
	log.Println("gRPC server stopped")
}

func (s *Server) Stop() {
	log.Println("Force stopping gRPC server...")
	s.grpcServer.Stop()
	log.Println("gRPC server force stopped")
}
