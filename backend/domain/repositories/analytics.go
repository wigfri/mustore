package repositories

import "time"

type UserAnalytics struct {
	TotalUsers       int64     `json:"total_users"`
	VerifiedUsers    int64     `json:"verified_users"`
	AdminsCount      int64     `json:"admins_count"`
	RegistrationsDay int64     `json:"registrations_day"`
	LoginsDay        int64     `json:"logins_day"`
	ActiveUsersDay   int64     `json:"active_users_day"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
}

type Analytics interface {
	GetUserAnalytics(periodStart, periodEnd time.Time) (*UserAnalytics, error)
	RecordLogin(userID string) error
}

