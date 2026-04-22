package v1

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/cases/create_instrument"
	"github.com/wigfri/mustore/domain/cases/delete_instrument"
	"github.com/wigfri/mustore/domain/cases/get_all_instruments"
	"github.com/wigfri/mustore/domain/cases/get_instrument"
	"github.com/wigfri/mustore/domain/cases/update_instrument"
	"github.com/wigfri/mustore/domain/repositories"
)

// CreateInstrument godoc
// @Summary Create a new instrument
// @Description Create instrument card in catalog (admin only)
// @Tags Instruments
// @Accept  json
// @Produce  json
// @Param request body create_instrument.Request true "Instrument payload"
// @Success 200 {object} RawResponse
// @Failure 400 {object} RawResponse
// @Failure 401 {object} RawResponse
// @Failure 403 {object} RawResponse
// @Failure 500 {object} RawResponse
// @Router /api/v1/instruments/ [post]
func CreateInstrument(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	if err := requireAdmin(c, ctx); err != nil {
		return DomainError(err)
	}
	var req create_instrument.Request
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(err)
	}
	res, err := create_instrument.Run(c, req)
	if err != nil {
		return DomainError(err)
	}
	return OK(res)
}

// GetInstrument godoc
// @Summary Get instrument by ID
// @Tags Instruments
// @Produce  json
// @Param id path string true "Instrument ID"
// @Success 200 {object} RawResponse
// @Failure 400 {object} RawResponse
// @Failure 500 {object} RawResponse
// @Router /api/v1/instruments/{id} [get]
func GetInstrument(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	res, err := get_instrument.Run(c, get_instrument.Request{Id: ctx.Params("id")})
	if err != nil {
		return DomainError(err)
	}
	return OK(res)
}

// GetAllInstruments godoc
// @Summary List instruments
// @Tags Instruments
// @Produce  json
// @Param search query string false "Search by name or brand"
// @Param category query string false "Category"
// @Param brand query string false "Brand"
// @Param min_price query int false "Minimal price"
// @Param max_price query int false "Maximum price"
// @Param only_active query bool false "Filter by active flag"
// @Param in_stock query bool false "Only in stock"
// @Param limit query int false "Page size"
// @Param offset query int false "Page offset"
// @Param sort_by query string false "created_at|price|name"
// @Param sort_order query string false "asc|desc"
// @Success 200 {object} RawResponse
// @Failure 500 {object} RawResponse
// @Router /api/v1/instruments/ [get]
func GetAllInstruments(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	filter := repositories.InstrumentListFilter{
		Search:      ctx.Query("search"),
		Category:    ctx.Query("category"),
		Brand:       ctx.Query("brand"),
		InStockOnly: ctx.QueryBool("in_stock", false),
		Limit:       ctx.QueryInt("limit", 20),
		Offset:      ctx.QueryInt("offset", 0),
		SortBy:      ctx.Query("sort_by", "created_at"),
		SortOrder:   ctx.Query("sort_order", "desc"),
	}
	if v := ctx.Query("min_price"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.MinPrice = &n
		}
	}
	if v := ctx.Query("max_price"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.MaxPrice = &n
		}
	}
	if v := ctx.Query("only_active"); v != "" {
		b := ctx.QueryBool("only_active", true)
		filter.OnlyActive = &b
	}

	res, err := get_all_instruments.Run(c, get_all_instruments.Request{Filter: filter})
	if err != nil {
		return DomainError(err)
	}
	return OK(res)
}

// UpdateInstrument godoc
// @Summary Update instrument
// @Description Update instrument card (admin only)
// @Tags Instruments
// @Accept  json
// @Produce  json
// @Param id path string true "Instrument ID"
// @Param request body update_instrument.Request true "Instrument update payload"
// @Success 200 {object} RawResponse
// @Failure 400 {object} RawResponse
// @Failure 401 {object} RawResponse
// @Failure 403 {object} RawResponse
// @Failure 500 {object} RawResponse
// @Router /api/v1/instruments/{id} [put]
func UpdateInstrument(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	if err := requireAdmin(c, ctx); err != nil {
		return DomainError(err)
	}
	var req update_instrument.Request
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(err)
	}
	req.Id = ctx.Params("id")
	res, err := update_instrument.Run(c, req)
	if err != nil {
		return DomainError(err)
	}
	return OK(res)
}

// DeleteInstrument godoc
// @Summary Delete instrument
// @Description Soft delete instrument (admin only)
// @Tags Instruments
// @Produce  json
// @Param id path string true "Instrument ID"
// @Success 200 {object} RawResponse
// @Failure 400 {object} RawResponse
// @Failure 401 {object} RawResponse
// @Failure 403 {object} RawResponse
// @Failure 500 {object} RawResponse
// @Router /api/v1/instruments/{id} [delete]
func DeleteInstrument(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	if err := requireAdmin(c, ctx); err != nil {
		return DomainError(err)
	}
	if err := delete_instrument.Run(c, delete_instrument.Request{Id: ctx.Params("id")}); err != nil {
		return DomainError(err)
	}
	return OK(nil)
}
