package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	pb "github.com/sumit-tembe/grpc-svc/pkg/grpc/user"
	logger "github.com/sumit-tembe/grpc-svc/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

var (
	customFunc = func(code codes.Code) logrus.Level {
		switch code {
		case codes.OK:
			return logrus.InfoLevel
		default:
			return logrus.ErrorLevel
		}
	}
)

type server struct {
	pb.UnimplementedUsersServer
}

func (s *server) GetUsers(ctx context.Context, in *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	return &pb.GetUsersResponse{
		Users: []*pb.User{{Name: "Sumit"}},
	}, nil
}

func main() {

	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)
	// bind OS events to the signal channel
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Logrus entry is used, allowing pre-definition of certain fields by the user.
	logrusEntry := logrus.NewEntry(logger.Logger)
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(customFunc),
	}
	// Make sure that log statements internal to gRPC library are logged using the logrus Logger as well.
	// grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_logrus.StreamServerInterceptor(logrusEntry, opts...),
		),
	)
	pb.RegisterUsersServer(s, &server{})
	reflection.Register(s)

	go func() {
		if err := s.Serve(listener); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		logrusEntry.Errorf("Fatal error: %+v", err)
	case <-stopChan:
		logger.Logger.Info("Received graceful shutdown signal")
	}
}
