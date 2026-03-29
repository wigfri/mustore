package get_user

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
	User *models.User `json:"user"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	user, err := c.Connection().User().GetUser(uuid.MustParse(r.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to get example due [%s]", err)
	}

	return &Response{User: user}, nil
}
