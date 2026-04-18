package repositories

import (
	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain/models"
)

type User interface {
	Insert(user *models.User) (string, error)
	GetUser(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	All() ([]*models.User, error)
	Update(id uuid.UUID, user *models.User) (*models.User, error)
	DeleteFromDb(id uuid.UUID) error
}
