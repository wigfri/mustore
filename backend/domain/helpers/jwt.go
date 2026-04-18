package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain/models"
)

type AccessClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func SignAccessToken(user *models.User, secret []byte, ttl time.Duration) (string, error) {
	if len(secret) == 0 {
		return "", fmt.Errorf("jwt secret is empty")
	}
	now := time.Now()
	claims := AccessClaims{
		Email: user.Email,
		Role:  string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &claims).SignedString(secret)
}

func ParseAccessToken(tokenString string, secret []byte) (*AccessClaims, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("jwt secret is empty")
	}
	tok, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*AccessClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func UserIDFromClaims(c *AccessClaims) (uuid.UUID, error) {
	if c.Subject == "" {
		return uuid.Nil, fmt.Errorf("missing subject")
	}
	return uuid.Parse(c.Subject)
}
