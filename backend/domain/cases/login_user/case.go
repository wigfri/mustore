package login_user

import (
	"errors"

	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/helpers"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	const op = "login_user.Run"
	logger := c.Services().Logger()

	if r.Email == "" || r.Password == "" {
		return nil, domain.ErrUnauthorized("invalid email or password")
	}

	user, err := c.Connection().User().GetUserByEmail(r.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUnauthorized("invalid email or password")
		}
		logger.Error("login lookup failed", "op", op, "error", err.Error())
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(r.Password)); err != nil {
		return nil, domain.ErrUnauthorized("invalid email or password")
	}

	if !user.IsVerified {
		return nil, domain.ErrForbidden("account is not verified")
	}

	cfg := c.Services().Config()
	ttl := cfg.JwtTTL()
	token, err := helpers.SignAccessToken(user, []byte(cfg.JwtSecret()), ttl)
	if err != nil {
		logger.Error("failed to sign access token", "op", op, "error", err.Error())
		return nil, err
	}

	logger.Info("user logged in", "op", op, "user_id", user.Id.String())
	return &Response{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(ttl.Seconds()),
	}, nil
}
