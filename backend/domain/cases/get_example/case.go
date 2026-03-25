package get_example

import (
	"fmt"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/models"
)

type Request struct {
	Id string
}

type Response struct {
	Example *models.Example `json:"example"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	example, err := c.Connection().Example().GetExample(r.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get example due [%s]", err)
	}

	return &Response{Example: example}, nil
}
