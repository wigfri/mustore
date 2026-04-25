package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InstrumentCategory string

const (
	CategoryGuitar     InstrumentCategory = "guitar"
	CategoryPiano      InstrumentCategory = "piano"
	CategoryDrums      InstrumentCategory = "drums"
	CategoryWind       InstrumentCategory = "wind"
	CategoryString     InstrumentCategory = "string"
	CategoryAccessory  InstrumentCategory = "accessory"
	CategoryElectronic InstrumentCategory = "electronic"
)

type Instrument struct {
	Id          uuid.UUID      `json:"id"`
	Name        string         `json:"name" gorm:"index"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;not null"`
	Brand       string         `json:"brand" gorm:"index"`
	Category    string         `json:"category" gorm:"index"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Price       int64          `json:"price" gorm:"index"`
	Currency    string         `json:"currency"`
	Stock       int            `json:"stock"`
	SKU         string         `json:"sku" gorm:"uniqueIndex;not null"`
	ImageURL    string         `json:"image_url"`
	IsActive    bool           `json:"is_active" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (i *Instrument) Validate() error {
	if i == nil {
		return errors.New("instrument is nil")
	}
	if strings.TrimSpace(i.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(i.Brand) == "" {
		return errors.New("brand is required")
	}
	if strings.TrimSpace(i.Category) == "" {
		return errors.New("category is required")
	}
	if _, err := ParseInstrumentCategory(i.Category); err != nil {
		return err
	}
	if i.Price < 0 {
		return errors.New("price must be greater or equal to 0")
	}
	if i.Stock < 0 {
		return errors.New("stock must be greater or equal to 0")
	}
	switch strings.ToUpper(strings.TrimSpace(i.Currency)) {
	case "RUB", "USD", "EUR":
	default:
		return errors.New("currency must be one of RUB, USD, EUR")
	}
	return nil
}

func ParseInstrumentCategory(input string) (InstrumentCategory, error) {
	switch InstrumentCategory(strings.ToLower(strings.TrimSpace(input))) {
	case CategoryGuitar, CategoryPiano, CategoryDrums, CategoryWind, CategoryString, CategoryAccessory, CategoryElectronic:
		return InstrumentCategory(strings.ToLower(strings.TrimSpace(input))), nil
	default:
		return "", fmt.Errorf("invalid category: %s", input)
	}
}
