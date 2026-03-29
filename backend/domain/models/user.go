package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	Admin    Role = "admin"
	BaseUser Role = "base_user"
)

type User struct {
	Id           uuid.UUID      `json:"id"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	PasswordHash string         `json:"-"`
	IsVerified   bool           `json:"is_verified" example:"true"`
	Role         Role           `json:"role" example:"admin"`
	CreatedAt    time.Time      `json:"created_at" example:"2024-08-15T05:46:10.457952Z"`
	UpdatedAt    time.Time      `json:"updated_at" example:"2024-08-15T05:46:10.457952Z"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

func IsValidRole(input string) (Role, error) {
	role := Role(input)
	switch role {
	case Admin, BaseUser:
		return role, nil
	default:
		return "", errors.New("invalid role")
	}
}

func IsAdmin(role Role) bool {
	return role == Admin
}

func IsBaseUser(role Role) bool {
	return role == BaseUser
}

func NewBaseUser(id uuid.UUID, email, username, password string) *User {
	return &User{
		Id:           id,
		Email:        email,
		Username:     username,
		PasswordHash: password,
		Role:         BaseUser,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
