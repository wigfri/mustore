package get_instrument

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/models"
)

type Request struct {
	Id string
}

type Response struct {
	Instrument *models.Instrument `json:"instrument"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	id, err := uuid.Parse(r.Id)
	if err != nil {
		return nil, domain.ErrBadRequest("invalid instrument id")
	}
	instrument, err := c.Connection().Instrument().GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get instrument due [%s]", err)
	}
	return &Response{Instrument: instrument}, nil
}
