package repositories

import "github.com/wigfri/mustore/domain/models"

type Example interface {
	Insert(example *models.Example) (string, error)
	GetExample(id string) (*models.Example, error)
	All() ([]*models.Example, error)
}
