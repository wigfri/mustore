package v1

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/helpers"
	"github.com/wigfri/mustore/domain/models"
)

func requireAdmin(c domain.Context, ctx *fiber.Ctx) error {
	token := extractToken(ctx)
	if token == "" {
		return domain.ErrUnauthorized("missing access token")
	}
	claims, err := helpers.ParseAccessToken(token, []byte(c.Services().Config().JwtSecret()))
	if err != nil {
		return domain.ErrUnauthorized("invalid access token")
	}
	if !models.IsAdmin(models.Role(claims.Role)) {
		return domain.ErrForbidden("admin role required")
	}
	return nil
}

func extractToken(ctx *fiber.Ctx) string {
	auth := ctx.Get(fiber.HeaderAuthorization)
	if strings.HasPrefix(auth, bearerPrefix) {
		return strings.TrimSpace(strings.TrimPrefix(auth, bearerPrefix))
	}
	if cookieToken := strings.TrimSpace(ctx.Cookies(accessTokenCookieName)); cookieToken != "" {
		return cookieToken
	}
	return ""
}
