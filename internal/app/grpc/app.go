package grpc

import (
	"github.com/LeoUraltsev/auth-service/internal/application"
	userGrpc "github.com/LeoUraltsev/auth-service/internal/infrastructure/grpc"
	"github.com/LeoUraltsev/auth-service/internal/infrastructure/interceptors"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log           *slog.Logger
	gRPC          *grpc.Server
	tokenVerifier interceptors.TokenVerifier
	address       string
}

func NewApp(service application.UserService, log *slog.Logger, tokenVerifier interceptors.TokenVerifier, address string) *App {

	i := interceptors.New(log, tokenVerifier)

	gRPC := grpc.NewServer(
		grpc.ChainUnaryInterceptor(i.RequestID, i.Auth),
	)

	userGrpc.Register(gRPC, service, log)
	return &App{
		log:           log,
		gRPC:          gRPC,
		tokenVerifier: tokenVerifier,
		address:       address,
	}
}

func (a *App) Start() error {
	lis, err := net.Listen("tcp", a.address)
	if err != nil {
		return err
	}
	a.log.Info("grpc server listening on ", slog.String("addr", lis.Addr().String()))
	err = a.gRPC.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("shutting down grpc server")
	a.gRPC.GracefulStop()
	a.log.Info("grpc server stopped")
}
