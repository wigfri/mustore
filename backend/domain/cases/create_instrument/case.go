package create_instrument

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/models"
	"gorm.io/gorm"
)

type Request struct {
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Currency    string `json:"currency"`
	Stock       int    `json:"stock"`
	SKU         string `json:"sku"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`
}

type Response struct {
	Id string `json:"id"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	id := uuid.New()
	slug := makeSlug(r.Name)
	if slug == "" {
		return nil, domain.ErrBadRequest("name is required")
	}

	if existing, err := c.Connection().Instrument().GetBySlug(slug); err == nil && existing != nil {
		slug = slug + "-" + id.String()[:8]
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	inst := &models.Instrument{
		Id:          id,
		Name:        strings.TrimSpace(r.Name),
		Slug:        slug,
		Brand:       strings.TrimSpace(r.Brand),
		Category:    strings.ToLower(strings.TrimSpace(r.Category)),
		Type:        strings.TrimSpace(r.Type),
		Description: strings.TrimSpace(r.Description),
		Price:       r.Price,
		Currency:    strings.ToUpper(strings.TrimSpace(r.Currency)),
		Stock:       r.Stock,
		SKU:         strings.TrimSpace(r.SKU),
		ImageURL:    strings.TrimSpace(r.ImageURL),
		IsActive:    r.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := inst.Validate(); err != nil {
		return nil, domain.ErrBadRequest(err.Error())
	}
	createdID, err := c.Connection().Instrument().Insert(inst)
	if err != nil {
		return nil, err
	}
	return &Response{Id: createdID}, nil
}

func makeSlug(input string) string {
	s := strings.ToLower(strings.TrimSpace(input))
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}
