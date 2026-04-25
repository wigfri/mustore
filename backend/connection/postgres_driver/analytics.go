package postgres_driver

import (
	"time"

	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain/models"
	"github.com/wigfri/mustore/domain/repositories"
	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

type loginEvent struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId    uuid.UUID `gorm:"type:uuid;index"`
	CreatedAt time.Time `gorm:"index"`
}

func (e loginEvent) TableName() string {
	return "login_events"
}

func makeLoginEvent(model *models.LoginEvent) loginEvent {
	return loginEvent{
		Id:        model.Id,
		UserId:    model.UserId,
		CreatedAt: model.CreatedAt,
	}
}

func (r *analyticsRepository) RecordLogin(userID string) error {
	parsed, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	event := makeLoginEvent(&models.LoginEvent{
		Id:        uuid.New(),
		UserId:    parsed,
		CreatedAt: time.Now().UTC(),
	})
	return r.db.Create(&event).Error
}

func (r *analyticsRepository) GetUserAnalytics(periodStart, periodEnd time.Time) (*repositories.UserAnalytics, error) {
	stats := &repositories.UserAnalytics{
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
	}

	if err := r.db.Model(&user{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&user{}).Where("is_verified = ?", true).Count(&stats.VerifiedUsers).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&user{}).Where("role = ?", string(models.Admin)).Count(&stats.AdminsCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&user{}).
		Where("created_at >= ? AND created_at < ?", periodStart, periodEnd).
		Count(&stats.RegistrationsDay).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&loginEvent{}).
		Where("created_at >= ? AND created_at < ?", periodStart, periodEnd).
		Count(&stats.LoginsDay).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&loginEvent{}).
		Where("created_at >= ? AND created_at < ?", periodStart, periodEnd).
		Distinct("user_id").
		Count(&stats.ActiveUsersDay).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

