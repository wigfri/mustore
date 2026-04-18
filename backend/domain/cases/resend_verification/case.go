package resend_verification

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/helpers"
	"gorm.io/gorm"
)

type Request struct {
	Email string `json:"email"`
}

type Response struct {
	Ok bool `json:"ok"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	const op = "resend_verification.Run"
	ctx := context.Background()
	email := strings.TrimSpace(strings.ToLower(r.Email))
	if email == "" {
		return nil, domain.ErrBadRequest("email is required")
	}

	user, err := c.Connection().User().GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &Response{Ok: true}, nil
		}
		c.Services().Logger().Error("resend verification lookup failed", "op", op, "error", err.Error())
		return nil, err
	}

	if user.IsVerified {
		return nil, domain.ErrBadRequest("account already verified")
	}

	plain, err := helpers.GenerateOTPCode(6)
	if err != nil {
		return nil, err
	}
	if err := c.Services().OTPStore().SetEmailVerification(ctx, user.Id, plain, 24*time.Hour); err != nil {
		return nil, err
	}

	if err := c.Services().MailQueue().PublishVerificationEmail(user.Email, plain); err != nil {
		c.Services().Logger().Error("failed to send verification email", "op", op, "error", err.Error())
		return nil, err
	}

	c.Services().Logger().Info("verification email resent", "op", op, "user_id", user.Id.String())
	return &Response{Ok: true}, nil
}
