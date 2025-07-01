package jwt

import (
	"github.com/LeoUraltsev/auth-service/internal/app/logger"
	"github.com/LeoUraltsev/auth-service/internal/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//todo: отдельный логгер заглушка для тестов

func TestToken_GenerateToken(t *testing.T) {
	log, _ := logger.NewLogger("development")
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "14f982080eacd7e38bd7a74fc0519946",
			Expiration: 1 * time.Hour,
		},
	}
	tkn := NewToken(log.Log, cfg)
	id := uuid.New()
	token, err := tkn.GenerateToken(id)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	claims, err := tkn.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
	tokenExp := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTEzNzkxMTQsInVzZXJfaWQiOiJmMTkzMzY1Mi0zOTcyLTQ4ZDQtYjhlNy1hODM2YjZjNmE4NDkifQ.m00yEpITDpyPt8d4ksfcImo6Vm3oH_BbPlA0A2idZyA"
	_, err = tkn.ValidateToken(tokenExp)
	assert.Error(t, err)

}
