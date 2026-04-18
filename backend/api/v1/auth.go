package v1

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/cases/login_user"
	"github.com/wigfri/mustore/domain/cases/logout_user"
	"github.com/wigfri/mustore/domain/cases/resend_verification"
	"github.com/wigfri/mustore/domain/cases/verify_email"
)

const bearerPrefix = "Bearer "

// Login godoc
// @Summary Log in
// @Description Step 1: email+password → sends login code by email (requires_code). Step 2: same credentials + code → Bearer token.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param  login_user.Request  body  login_user.Request  true  "Credentials (code optional on second step)"
// @Success 200 {object} RawResponse "Access token or requires_code"
// @Failure 400 {object} RawResponse "Bad request"
// @Failure 401 {object} RawResponse "Unauthorized"
// @Failure 403 {object} RawResponse "Forbidden"
// @Failure 500 {object} RawResponse "Internal server error"
// @Router /api/v1/auth/login [post]
func Login(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	var req login_user.Request
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(err)
	}

	res, err := login_user.Run(c, req)
	if err != nil {
		return DomainError(err)
	}
	return OK(res)
}

// Logout godoc
// @Summary Log out
// @Description Validates the Bearer token if present; clients should discard the token after success.
// @Tags Auth
// @Produce  json
// @Param Authorization header string false "Bearer access_token"
// @Success 200 {object} RawResponse "Logged out"
// @Failure 401 {object} RawResponse "Unauthorized"
// @Failure 500 {object} RawResponse "Internal server error"
// @Router /api/v1/auth/logout [post]
func Logout(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	auth := ctx.Get(fiber.HeaderAuthorization)
	token := ""
	if strings.HasPrefix(auth, bearerPrefix) {
		token = strings.TrimSpace(strings.TrimPrefix(auth, bearerPrefix))
	}

	if err := logout_user.Run(c, logout_user.Request{AccessToken: token}); err != nil {
		return DomainError(err)
	}
	return OK(nil)
}

// VerifyEmail godoc
// @Summary Confirm email after registration
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param  verify_email.Request  body  verify_email.Request  true  "Email and code from letter"
// @Success 200 {object} RawResponse "Verified"
// @Failure 400 {object} RawResponse "Bad request"
// @Failure 401 {object} RawResponse "Unauthorized"
// @Failure 500 {object} RawResponse "Internal server error"
// @Router /api/v1/auth/verify-email [post]
func VerifyEmail(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	var req verify_email.Request
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(err)
	}
	res, err := verify_email.Run(c, req)
	if err != nil {
		return DomainError(err)
	}
	return OK(res)
}

// ResendVerification godoc
// @Summary Resend registration confirmation code
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param  resend_verification.Request  body  resend_verification.Request  true  "Email"
// @Success 200 {object} RawResponse "Sent (or unknown email — same response)"
// @Failure 400 {object} RawResponse "Bad request"
// @Failure 500 {object} RawResponse "Internal server error"
// @Router /api/v1/auth/resend-verification [post]
func ResendVerification(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	var req resend_verification.Request
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(err)
	}
	res, err := resend_verification.Run(c, req)
	if err != nil {
		return DomainError(err)
	}
	return OK(res)
}
