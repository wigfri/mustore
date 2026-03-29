package get_all_users

import (
	"fmt"

	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/models"
)

type Response struct {
	Users []*models.User `json:"users"`
}

func Run(c domain.Context) (*Response, error) {
	users, err := c.Connection().User().All()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve roles due [%s]", err)
	}

	return &Response{Users: users}, nil
}
