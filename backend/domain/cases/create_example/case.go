package create_example

import (
	"fmt"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/models"
	"math/rand"
	"strconv"
)

type Request struct {
	Example *models.Example `json:"example"`
}

type Response struct {
	Id string `json:"id"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	id := rand.Intn(1000)

	r.Example.Id = strconv.Itoa(id)

	exampleId, err := c.Connection().Example().Insert(r.Example)
	if err != nil {
		return nil, fmt.Errorf("unable to create example due [%s]", err)
	}

	return &Response{Id: exampleId}, nil
}
