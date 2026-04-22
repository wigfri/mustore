package repositories

import (
	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain/models"
)

type InstrumentListFilter struct {
	Search      string
	Category    string
	Brand       string
	MinPrice    *int64
	MaxPrice    *int64
	OnlyActive  *bool
	InStockOnly bool
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   string
}

type Instrument interface {
	Insert(instrument *models.Instrument) (string, error)
	GetByID(id uuid.UUID) (*models.Instrument, error)
	GetBySlug(slug string) (*models.Instrument, error)
	All(filter InstrumentListFilter) ([]*models.Instrument, error)
	Update(id uuid.UUID, instrument *models.Instrument) (*models.Instrument, error)
	DeleteFromDb(id uuid.UUID) error
}
