package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/cases/get_admin_analytics"
)

// AdminAnalytics godoc
// @Summary Admin analytics dashboard
// @Description Daily user activity metrics (registrations, logins, active users) and totals.
// @Tags Admin
// @Produce  json
// @Param Authorization header string true "Bearer access_token (admin)"
// @Success 200 {object} RawResponse
// @Failure 401 {object} RawResponse
// @Failure 403 {object} RawResponse
// @Failure 500 {object} RawResponse
// @Router /api/v1/admin/analytics [get]
func AdminAnalytics(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	if err := requireAdmin(c, ctx); err != nil {
		return DomainError(err)
	}

	res, err := get_admin_analytics.Run(c)
	if err != nil {
		return DomainError(err)
	}

	return OK(res)
}

