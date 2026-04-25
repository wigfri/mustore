package update_instrument

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/models"
)

type Request struct {
	Id          string `json:"id"`
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
	Instrument *models.Instrument `json:"instrument"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	id, err := uuid.Parse(r.Id)
	if err != nil {
		return nil, domain.ErrBadRequest("invalid instrument id")
	}
	existing, err := c.Connection().Instrument().GetByID(id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(r.Name) != "" && existing.Name != r.Name {
		existing.Name = strings.TrimSpace(r.Name)
		existing.Slug = makeSlug(r.Name)
	}
	if strings.TrimSpace(r.Brand) != "" {
		existing.Brand = strings.TrimSpace(r.Brand)
	}
	if strings.TrimSpace(r.Category) != "" {
		existing.Category = strings.ToLower(strings.TrimSpace(r.Category))
	}
	if strings.TrimSpace(r.Type) != "" {
		existing.Type = strings.TrimSpace(r.Type)
	}
	if strings.TrimSpace(r.Description) != "" {
		existing.Description = strings.TrimSpace(r.Description)
	}
	if r.Price >= 0 {
		existing.Price = r.Price
	}
	if strings.TrimSpace(r.Currency) != "" {
		existing.Currency = strings.ToUpper(strings.TrimSpace(r.Currency))
	}
	if r.Stock >= 0 {
		existing.Stock = r.Stock
	}
	if strings.TrimSpace(r.SKU) != "" {
		existing.SKU = strings.TrimSpace(r.SKU)
	}
	if strings.TrimSpace(r.ImageURL) != "" {
		existing.ImageURL = strings.TrimSpace(r.ImageURL)
	}
	existing.IsActive = r.IsActive
	existing.UpdatedAt = time.Now()

	if err := existing.Validate(); err != nil {
		return nil, domain.ErrBadRequest(err.Error())
	}
	updated, err := c.Connection().Instrument().Update(id, existing)
	if err != nil {
		return nil, err
	}
	return &Response{Instrument: updated}, nil
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
