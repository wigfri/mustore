package logout_user

import (
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/helpers"
)

type Request struct {
	AccessToken string
}

func Run(c domain.Context, r Request) error {
	const op = "logout_user.Run"
	logger := c.Services().Logger()

	if r.AccessToken == "" {
		return nil
	}

	cfg := c.Services().Config()
	claims, err := helpers.ParseAccessToken(r.AccessToken, []byte(cfg.JwtSecret()))
	if err != nil {
		return domain.ErrUnauthorized("invalid or expired token")
	}

	userID, err := helpers.UserIDFromClaims(claims)
	if err != nil {
		return domain.ErrUnauthorized("invalid or expired token")
	}

	logger.Info("user logged out", "op", op, "user_id", userID.String())
	return nil
}
