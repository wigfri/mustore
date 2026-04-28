package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/cases/create_user"
	"github.com/wigfri/mustore/domain/cases/get_all_users"
	"github.com/wigfri/mustore/domain/cases/get_user"
)

// CreateUser godoc
// @Summary Create a new User
// @Description Create a new User with the provided details
// @Tags Users
// @Accept  json
// @Produce  json
// @Param  create_user.Request  body  create_user.Request  true  "User request"
// @Success 200 {object} RawResponse "Returns the ID of the created User"
// @Failure 400 {object} RawResponse "Bad request"
// @Failure 500 {object} RawResponse "Internal server error"
// @Router /api/v1/user/ [post]
func CreateUser(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	var req create_user.Request

	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(err)
	}

	id, err := create_user.Run(c, req)
	if err != nil {
		return DomainError(err)
	}

	return OK(id)
}

// GetAllUsers godoc
// @Summary Get all Users
// @Description Retrieve a list of all Users
// @Tags Users
// @Produce  json
// @Success 200 {object} RawResponse "Returns a list of Users"
// @Failure 500 {object} RawResponse "Internal server error"
// @Router /api/v1/user/ [get]
func GetAllUser(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	Users, err := get_all_users.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(Users)
}

// GetUser godoc
// @Summary Get an User by ID
// @Description Retrieve a single User by its ID
// @Tags Users
// @Produce  json
// @Param  id  path  string  true  "User ID"
// @Success 200 {object} RawResponse "Returns the requested User"
// @Failure 500 {object} RawResponse "Internal server error"
// @Router /api/v1/user/{id} [get]
func GetUser(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	id := ctx.Params("id")

	var req get_user.Request
	req.Id = id

	User, err := get_user.Run(c, req)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(User)
}
