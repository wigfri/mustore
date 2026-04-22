package delete_instrument

import (
	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain"
)

type Request struct {
	Id string
}

func Run(c domain.Context, r Request) error {
	id, err := uuid.Parse(r.Id)
	if err != nil {
		return domain.ErrBadRequest("invalid instrument id")
	}
	return c.Connection().Instrument().DeleteFromDb(id)
}
