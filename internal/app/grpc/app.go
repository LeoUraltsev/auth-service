package grpc

import (
	"github.com/LeoUraltsev/auth-service/internal/domain/users"
	userGrpc "github.com/LeoUraltsev/auth-service/internal/infrastructure/grpc"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log     *slog.Logger
	gRPC    *grpc.Server
	address string
}

func NewApp(service users.UserServiceHandler, log *slog.Logger, address string) *App {
	gRPC := grpc.NewServer()

	userGrpc.Register(gRPC, service, log)
	return &App{
		log:     log,
		gRPC:    gRPC,
		address: address,
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
