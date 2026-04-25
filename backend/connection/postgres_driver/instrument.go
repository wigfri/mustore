package postgres_driver

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain/models"
	"github.com/wigfri/mustore/domain/repositories"
	"gorm.io/gorm"
)

type instrumentRepository struct {
	db *gorm.DB
}

type instrument struct {
	Id          uuid.UUID
	Name        string
	Slug        string
	Brand       string
	Category    string
	Type        string
	Description string
	Price       int64
	Currency    string
	Stock       int
	SKU         string
	ImageURL    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func (m instrument) model() *models.Instrument {
	return &models.Instrument{
		Id:          m.Id,
		Name:        m.Name,
		Slug:        m.Slug,
		Brand:       m.Brand,
		Category:    m.Category,
		Type:        m.Type,
		Description: m.Description,
		Price:       m.Price,
		Currency:    m.Currency,
		Stock:       m.Stock,
		SKU:         m.SKU,
		ImageURL:    m.ImageURL,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

func makeInstrument(m *models.Instrument) instrument {
	return instrument{
		Id:          m.Id,
		Name:        m.Name,
		Slug:        m.Slug,
		Brand:       m.Brand,
		Category:    m.Category,
		Type:        m.Type,
		Description: m.Description,
		Price:       m.Price,
		Currency:    m.Currency,
		Stock:       m.Stock,
		SKU:         m.SKU,
		ImageURL:    m.ImageURL,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

func (e *instrumentRepository) Insert(item *models.Instrument) (string, error) {
	model := makeInstrument(item)
	if err := e.db.Create(&model).Error; err != nil {
		return "", err
	}
	return model.Id.String(), nil
}

func (e *instrumentRepository) GetByID(id uuid.UUID) (*models.Instrument, error) {
	var result instrument
	if err := e.db.Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result.model(), nil
}

func (e *instrumentRepository) GetBySlug(slug string) (*models.Instrument, error) {
	var result instrument
	if err := e.db.Where("slug = ?", strings.ToLower(strings.TrimSpace(slug))).First(&result).Error; err != nil {
		return nil, err
	}
	return result.model(), nil
}

func (e *instrumentRepository) All(filter repositories.InstrumentListFilter) ([]*models.Instrument, error) {
	q := e.db.Model(&instrument{})
	if filter.Search != "" {
		p := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		q = q.Where("LOWER(name) LIKE ? OR LOWER(brand) LIKE ?", p, p)
	}
	if filter.Category != "" {
		q = q.Where("LOWER(category) = ?", strings.ToLower(strings.TrimSpace(filter.Category)))
	}
	if filter.Brand != "" {
		q = q.Where("LOWER(brand) = ?", strings.ToLower(strings.TrimSpace(filter.Brand)))
	}
	if filter.MinPrice != nil {
		q = q.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		q = q.Where("price <= ?", *filter.MaxPrice)
	}
	if filter.OnlyActive != nil {
		q = q.Where("is_active = ?", *filter.OnlyActive)
	}
	if filter.InStockOnly {
		q = q.Where("stock > 0")
	}

	sortBy := "created_at"
	switch filter.SortBy {
	case "price", "name", "created_at":
		sortBy = filter.SortBy
	}
	sortOrder := "desc"
	if strings.EqualFold(filter.SortOrder, "asc") {
		sortOrder = "asc"
	}
	q = q.Order(sortBy + " " + sortOrder)

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	q = q.Limit(limit).Offset(offset)

	var rows []instrument
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*models.Instrument, len(rows))
	for i, row := range rows {
		out[i] = row.model()
	}
	return out, nil
}

func (e *instrumentRepository) Update(id uuid.UUID, item *models.Instrument) (*models.Instrument, error) {
	if item == nil {
		return nil, errors.New("instrument is nil")
	}
	updates := map[string]interface{}{
		"name":        item.Name,
		"slug":        item.Slug,
		"brand":       item.Brand,
		"category":    item.Category,
		"type":        item.Type,
		"description": item.Description,
		"price":       item.Price,
		"currency":    item.Currency,
		"stock":       item.Stock,
		"sku":         item.SKU,
		"image_url":   item.ImageURL,
		"is_active":   item.IsActive,
		"updated_at":  time.Now(),
	}
	if err := e.db.Model(&instrument{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	return e.GetByID(id)
}

func (e *instrumentRepository) DeleteFromDb(id uuid.UUID) error {
	return e.db.Delete(&instrument{Id: id}).Error
}
