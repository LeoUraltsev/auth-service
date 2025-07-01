package interceptors

import (
	"context"
	"github.com/LeoUraltsev/auth-service/internal/infrastructure/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"strings"
)

var KeyCtxRequestID = "request_id"

type TokenVerifier interface {
	ValidateToken(token string) (*jwt.AuthClaims, error)
}

type Interceptors struct {
	log           *slog.Logger
	tokenVerifier TokenVerifier
}

func New(log *slog.Logger, verifier TokenVerifier) *Interceptors {
	return &Interceptors{
		log:           log,
		tokenVerifier: verifier,
	}
}

func (i *Interceptors) RequestID(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	requestID := uuid.New()
	log := i.log.With("request_id", requestID)
	log.Info("new call", slog.String("method", info.FullMethod))
	ctx = context.WithValue(ctx, KeyCtxRequestID, requestID)

	return handler(ctx, req)

}

func (i *Interceptors) Auth(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if info.FullMethod == "/auth.UserService/Login" || info.FullMethod == "/auth.UserService/CreateUser" {
		return handler(ctx, req)
	}
	log := i.log.With("method", info.FullMethod)
	id, ok := ctx.Value(KeyCtxRequestID).(uuid.UUID)
	if !ok {
		log.Warn("context value for key 'request_id' not found")
	}

	log = log.With("requestID", id)

	log.Info("access verification")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Warn("context metadata not found")
		return nil, status.Error(codes.Unauthenticated, "no metadata found")
	}

	a := md.Get("authorization")
	if len(a) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no authorization metadata")
	}

	if a[0] == "" {
		log.Warn("no token found")
		return nil, status.Error(codes.Unauthenticated, "no token found")
	}

	token := strings.TrimPrefix(a[0], "Bearer ")

	_, err := i.tokenVerifier.ValidateToken(token)
	if err != nil {
		log.Warn("invalid token", slog.String("token", token), slog.String("error", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "no token found")
	}

	return handler(ctx, req)
}
