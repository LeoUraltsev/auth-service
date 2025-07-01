package grpc

import (
	"context"
	"github.com/LeoUraltsev/auth-service/internal/domain/users"
	auth1 "github.com/LeoUraltsev/proto/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

// todo: покрыть тестами
// todo: мапперы

type userGRPCApi struct {
	auth1.UnimplementedUserServiceServer
	service users.UserServiceHandler
	log     *slog.Logger
}

func Register(gRPC *grpc.Server, service users.UserServiceHandler, log *slog.Logger) {
	auth1.RegisterUserServiceServer(gRPC, &userGRPCApi{
		service: service,
		log:     log,
	})
}

//todo: интерсептор для получения id запроса
//todo: разные типы ошибок (невалидные данные, пользователь уже есть, внутренняя ошибка)

func (a *userGRPCApi) CreateUser(ctx context.Context, request *auth1.CreateUserRequest) (*auth1.CreateUserResponse, error) {
	log := a.log
	log.Info("creating new user")

	id, err := a.service.CreateUser(ctx, request.Name, request.Email, request.Password)
	if err != nil {
		log.Error("failed to create user", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &auth1.CreateUserResponse{Id: id.String()}, nil
}

//todo: интерсептор для проверки может ли пользователь дергать данную ручку

func (a *userGRPCApi) GetUser(ctx context.Context, request *auth1.GetUserRequest) (*auth1.GetUserResponse, error) {
	log := a.log
	log.Info("getting user")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		log.Error("failed to parse uuid", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "incorrect id")
	}
	user, err := a.service.GetUser(ctx, id)
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	//todo: пароль не должен возвращаться
	return &auth1.GetUserResponse{
		User: &auth1.User{
			Id:        user.ID().String(),
			Name:      user.Name().String(),
			Email:     user.Email().String(),
			Password:  "",
			CreatedAt: timestamppb.New(user.CreatedAt()),
			UpdatedAt: timestamppb.New(user.UpdatedAt()),
		},
	}, nil
}

//todo: интерсептор для проверки может ли пользователь дергать данную ручку
//todo: логи

func (a *userGRPCApi) GetListUsers(ctx context.Context, request *auth1.GetListUserRequest) (*auth1.GetListUserResponse, error) {
	usrs, err := a.service.GetListUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get users")
	}
	res := make([]*auth1.User, 0, len(usrs))
	for _, usr := range usrs {
		a.log.Debug("got user", slog.Any("user", usr.ID()))
		res = append(res, &auth1.User{
			Id:        usr.ID().String(),
			Name:      usr.Name().String(),
			Email:     usr.Email().String(),
			Password:  "",
			CreatedAt: timestamppb.New(usr.CreatedAt()),
			UpdatedAt: timestamppb.New(usr.UpdatedAt()),
		})
	}

	return &auth1.GetListUserResponse{Users: res}, nil
}

//todo: интерсептор для проверки может ли пользователь дергать данную ручку

func (a *userGRPCApi) UpdateUser(ctx context.Context, request *auth1.UpdateUserRequest) (*auth1.UpdateUserResponse, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "incorrect id")
	}
	err = a.service.UpdateUser(ctx, id, request.Name, request.Email, request.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
	}
	return &auth1.UpdateUserResponse{}, nil
}

//todo: интерсептор для проверки может ли пользователь дергать данную ручку

func (a *userGRPCApi) DeleteUser(ctx context.Context, request *auth1.DeleteUserRequest) (*auth1.DeleteUserResponse, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "incorrect id")
	}
	err = a.service.DeleteUser(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete user")
	}
	return &auth1.DeleteUserResponse{Success: true}, nil
}

func (a *userGRPCApi) Login(ctx context.Context, request *auth1.LoginRequest) (*auth1.LoginResponse, error) {
	log := a.log.With("email", request.Email)
	log.Info("logging in")
	token, err := a.service.Login(ctx, request.Email, request.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}
	return &auth1.LoginResponse{Token: token}, nil
}

func (a *userGRPCApi) mustEmbedUnimplementedUserServiceServer() {}
