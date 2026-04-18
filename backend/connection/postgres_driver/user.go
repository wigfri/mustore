package postgres_driver

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (e *userRepository) Update(id uuid.UUID, u *models.User) (*models.User, error) {
	if u == nil {
		return nil, errors.New("user is nil")
	}

	var n int64
	if err := e.db.Model(&user{}).Where("id = ?", id).Count(&n).Error; err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	updates := map[string]interface{}{
		"username":   u.Username,
		"updated_at": time.Now(),
	}

	if err := e.db.Model(&user{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	return e.GetUser(id)
}

func (e *userRepository) DeleteFromDb(id uuid.UUID) error {
	return e.db.Unscoped().Delete(&models.User{Id: id}).Error
}

type user struct {
	Id           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	IsVerified   bool
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func (e user) model() *models.User {
	return &models.User{
		Id:           e.Id,
		Username:     e.Username,
		Email:        e.Email,
		PasswordHash: e.PasswordHash,
		IsVerified:   e.IsVerified,
		Role:         models.Role(e.Role),
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		DeletedAt:    e.DeletedAt,
	}
}

func makeUser(e *models.User) user {
	return user{
		Id:           e.Id,
		Username:     e.Username,
		Email:        e.Email,
		PasswordHash: e.PasswordHash,
		IsVerified:   e.IsVerified,
		Role:         string(e.Role),
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		DeletedAt:    e.DeletedAt,
	}
}

func (e *userRepository) Insert(user *models.User) (string, error) {
	userModel := makeUser(user)

	if err := e.db.Create(userModel).Error; err != nil {
		return "", err
	}

	return userModel.Id.String(), nil
}

func (e *userRepository) GetUser(id uuid.UUID) (*models.User, error) {
	var result user

	if err := e.db.Where(models.User{Id: id}).First(&result).Error; err != nil {
		return nil, err
	}

	return result.model(), nil
}

func (e *userRepository) GetUserByEmail(email string) (*models.User, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return nil, gorm.ErrRecordNotFound
	}

	var result user
	if err := e.db.Where("LOWER(email) = ?", normalized).First(&result).Error; err != nil {
		return nil, err
	}

	return result.model(), nil
}

func (e *userRepository) All() ([]*models.User, error) {
	var result []user

	if err := e.db.Find(&result).Error; err != nil {
		return nil, err
	}

	out := make([]*models.User, len(result))

	for i, em := range result {
		out[i] = em.model()
	}

	return out, nil
}
