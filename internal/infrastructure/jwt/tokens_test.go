package jwt

import (
	"github.com/LeoUraltsev/auth-service/internal/app/logger"
	"github.com/LeoUraltsev/auth-service/internal/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToken_GenerateToken(t *testing.T) {
	log, _ := logger.NewLogger("development")
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "qwerty",
			Expiration: 1 * time.Hour,
		},
	}
	tkn := NewToken(log.Log, cfg)
	id := uuid.New()
	token, err := tkn.GenerateToken(id)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	err = tkn.ValidateToken(token)
	assert.NoError(t, err)

}
