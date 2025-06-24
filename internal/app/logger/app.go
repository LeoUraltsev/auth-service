package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Env string

const development Env = "development"
const prod Env = "prod"

type Logger struct {
	Log *slog.Logger
}

func NewLogger(env string) (*Logger, error) {
	var log *slog.Logger

	switch env {
	case development.String():
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,
			Level:       slog.LevelDebug,
			ReplaceAttr: nil,
		}))

	case prod.String():
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelInfo,
		}))
	default:
		return nil, fmt.Errorf("unknown environment: %s", env)
	}
	return &Logger{Log: log}, nil
}

func (e Env) String() string {
	return string(e)
}
