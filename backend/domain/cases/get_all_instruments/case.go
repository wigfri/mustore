package get_all_instruments

import (
	"fmt"

	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/models"
	"github.com/wigfri/mustore/domain/repositories"
)

type Request struct {
	Filter repositories.InstrumentListFilter
}

type Response struct {
	Instruments []*models.Instrument `json:"instruments"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	items, err := c.Connection().Instrument().All(r.Filter)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve instruments due [%s]", err)
	}
	return &Response{Instruments: items}, nil
}
