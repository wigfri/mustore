package get_admin_analytics

import (
	"time"

	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/repositories"
)

type Response struct {
	User *repositories.UserAnalytics `json:"user"`
}

func Run(c domain.Context) (*Response, error) {
	now := time.Now().UTC()
	periodStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.Add(24 * time.Hour)

	analytics, err := c.Connection().Analytics().GetUserAnalytics(periodStart, periodEnd)
	if err != nil {
		return nil, err
	}

	return &Response{
		User: analytics,
	}, nil
}

