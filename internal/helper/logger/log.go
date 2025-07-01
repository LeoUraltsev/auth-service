package logger

import (
	"context"
	"github.com/LeoUraltsev/auth-service/internal/infrastructure/interceptors"
	"github.com/google/uuid"
	"log/slog"
)

func LogWithContext(ctx context.Context, log *slog.Logger) *slog.Logger {
	l := log
	reqID, ok := ctx.Value(interceptors.KeyCtxRequestID).(uuid.UUID)
	if ok {
		l = log.With("request_id", reqID.String())
	}

	userID, ok := ctx.Value(interceptors.KeyCtxUserID).(uuid.UUID)
	if ok {
		l = log.With("user_id", userID.String())
	}

	return l
}
