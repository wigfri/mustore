package login_user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/helpers"
	"github.com/wigfri/mustore/domain/models"
	"github.com/wigfri/mustore/domain/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code,omitempty"`
}

type Response struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RequiresCode bool   `json:"requires_code,omitempty"`
	OtpExpiresIn int64  `json:"otp_expires_in,omitempty"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	const op = "login_user.Run"
	ctx := context.Background()
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
		return nil, domain.ErrForbidden("confirm your email before signing in")
	}

	code := strings.TrimSpace(r.Code)
	if code == "" {
		return issueLoginOTP(ctx, c, user, op, logger)
	}

	ok, err := c.Services().OTPStore().VerifyLoginOTP(ctx, user.Id, code)
	if err != nil {
		logger.Error("login otp verify", "op", op, "error", err.Error())
		return nil, err
	}
	if !ok {
		return nil, domain.ErrUnauthorized("login code was not requested, expired, or invalid")
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

const loginOTPTTL = 10 * time.Minute

func issueLoginOTP(ctx context.Context, c domain.Context, user *models.User, op string, logger services.Logger) (*Response, error) {
	plain, err := helpers.GenerateOTPCode(6)
	if err != nil {
		return nil, err
	}
	if err := c.Services().OTPStore().SetLoginOTP(ctx, user.Id, plain, loginOTPTTL); err != nil {
		return nil, err
	}
	if err := c.Services().MailQueue().PublishLoginCodeEmail(user.Email, plain); err != nil {
		logger.Error("failed to send login code email", "op", op, "error", err.Error())
		return nil, err
	}

	logger.Info("login otp issued", "op", op, "user_id", user.Id.String())
	return &Response{
		RequiresCode: true,
		OtpExpiresIn: int64(loginOTPTTL.Seconds()),
	}, nil
}
