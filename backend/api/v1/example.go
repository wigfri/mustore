package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/cases/create_example"
	"github.com/wigfri/mustore/domain/cases/get_all_example"
	"github.com/wigfri/mustore/domain/cases/get_example"
)

// CreateExample godoc
// @Summary Create a new example
// @Description Create a new example with the provided details
// @Tags examples
// @Accept  json
// @Produce  json
// @Param  create_example.Request  body  create_example.Request  true  "Example request"
// @Success 200 {object} RawResponse "Returns the ID of the created example"
// @Failure 400 {object} RawResponse "Bad request"
// @Failure 500 {object} RawResponse "Internal server error"
// @Router /api/v1/example/ [post]
func CreateExample(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	var req create_example.Request

	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(err)
	}

	id, err := create_example.Run(c, req)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(id)
}

// GetAllExamples godoc
// @Summary Get all examples
// @Description Retrieve a list of all examples
// @Tags examples
// @Produce  json
// @Success 200 {object} RawResponse "Returns a list of examples"
// @Failure 500 {object} RawResponse "Internal server error"
// @Router /api/v1/example/ [get]
func GetAllExamples(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	examples, err := get_all_example.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(examples)
}

// GetExample godoc
// @Summary Get an example by ID
// @Description Retrieve a single example by its ID
// @Tags examples
// @Produce  json
// @Param  id  path  string  true  "Example ID"
// @Success 200 {object} RawResponse "Returns the requested example"
// @Failure 500 {object} RawResponse "Internal server error"
// @Router /api/v1/example/{id} [get]
func GetExample(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	id := ctx.Params("id")

	var req get_example.Request
	req.Id = id

	example, err := get_example.Run(c, req)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(example)
}
