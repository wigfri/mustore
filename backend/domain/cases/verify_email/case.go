package verify_email

import (
	"context"
	"errors"
	"strings"

	"github.com/wigfri/mustore/domain"
	"gorm.io/gorm"
)

type Request struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type Response struct {
	Ok bool `json:"ok"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	const op = "verify_email.Run"
	ctx := context.Background()
	email := strings.TrimSpace(strings.ToLower(r.Email))
	code := strings.TrimSpace(r.Code)
	if email == "" || code == "" {
		return nil, domain.ErrBadRequest("email and code are required")
	}

	user, err := c.Connection().User().GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUnauthorized("invalid email or code")
		}
		c.Services().Logger().Error("verify email lookup failed", "op", op, "error", err.Error())
		return nil, err
	}

	if user.IsVerified {
		return nil, domain.ErrBadRequest("account already verified")
	}

	ok, err := c.Services().OTPStore().VerifyEmailCode(ctx, user.Id, code)
	if err != nil {
		c.Services().Logger().Error("verify email otp store", "op", op, "error", err.Error())
		return nil, err
	}
	if !ok {
		return nil, domain.ErrUnauthorized("invalid email or code")
	}

	user.IsVerified = true
	if _, err := c.Connection().User().Update(user.Id, user); err != nil {
		return nil, err
	}

	c.Services().Logger().Info("email verified", "op", op, "user_id", user.Id.String())
	return &Response{Ok: true}, nil
}
