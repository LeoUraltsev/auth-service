package jwt

import (
	"fmt"
	"github.com/LeoUraltsev/auth-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type AuthClaims struct {
	*jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

type Token struct {
	log *slog.Logger
	cfg *config.Config
}

func NewToken(log *slog.Logger, cfg *config.Config) *Token {
	return &Token{
		log: log,
		cfg: cfg,
	}
}

func (t *Token) GenerateToken(userID uuid.UUID) (string, error) {
	log := t.log
	log.Info("Generating token")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(t.cfg.JWT.Expiration)),
		},
		UserID: userID,
	})
	signedString, err := token.SignedString([]byte(t.cfg.JWT.Secret))
	if err != nil {
		log.Warn("Failed to sign token")
		return "", err
	}
	log.Info("Signed token generated", slog.String("token", signedString))
	return signedString, err
}

func (t *Token) ValidateToken(token string) (*AuthClaims, error) {
	tkn, err := jwt.ParseWithClaims(
		token,
		&AuthClaims{},
		func(j *jwt.Token) (interface{}, error) {
			return []byte(t.cfg.JWT.Secret), nil
		},
	)
	if err != nil {
		t.log.Warn("Failed to parse token", slog.String("err", err.Error()))
		return nil, err
	}

	claims, ok := tkn.Claims.(*AuthClaims)
	if !ok {
		t.log.Warn("Failed to get claims", slog.String("token", token))
		return nil, fmt.Errorf("unknown claims type")
	}
	exp := claims.RegisteredClaims.ExpiresAt.Time
	if time.Now().UTC().After(exp) {
		t.log.Warn("Token expired", slog.String("token", token))
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
