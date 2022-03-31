package ports

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/juanmaabanto/go-ms-beers/common/errors"
	"github.com/juanmaabanto/go-ms-beers/common/responses"
	"github.com/juanmaabanto/go-ms-beers/internal/app"
	"github.com/juanmaabanto/go-ms-beers/internal/app/command"
	"github.com/juanmaabanto/go-ms-beers/internal/app/query"
	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}

// CreateBeer godoc
// @Summary Create a new beer.
// @Tags Beers
// @Accept json
// @Produce json
// @Param command body command.CreateBeer true "Object to be created."
// @Success 201 {int64} string "Id of the created object"
// @Failure 400 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/beers [post]
func (h HttpServer) AddBeer(c echo.Context) error {
	item := command.CreateBeer{}

	if err := c.Bind(&item); err != nil {
		panic(err)
	}

	if err := c.Validate(item); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		panic(errors.NewValidationError(Simple(validationErrors)))
	}

	id, err := h.app.Commands.CreateBeer.Handle(c.Request().Context(), item)

	if err != nil {
		panic(err)
	}

	c.Response().Header().Set("location", c.Request().URL.String()+"/"+fmt.Sprint(id))

	return c.JSON(http.StatusCreated, id)
}

// GetBeer godoc
// @Summary Get a beer by Id.
// @Tags Beers
// @Accept json
// @Produce json
// @Param id path int64  true  "Beer Id"
// @Success 200 {object} response.BeerResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/beers/{beerId} [get]
func (h HttpServer) GetBeer(c echo.Context) error {
	var beerId int64
	if n, err := strconv.Atoi(c.Param("beerId")); err == nil {
		beerId = int64(n)
	} else {
		panic(err)
	}

	item := query.GetBeerById{Id: beerId}

	result, err := h.app.Queries.GetBeerById.Handle(c.Request().Context(), item)

	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)
}

// ListBeer godoc
// @Summary Return a Beers List.
// @Tags Beers
// @Accept json
// @Produce json
// @Param name query string  false  "word to search"
// @Param pageSize query int  false  "Number of results per page"
// @Param start query string  false  "Page number"
// @Success 200 {object} responses.PaginatedResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/beers [get]
func (h HttpServer) ListBeer(c echo.Context) error {
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))

	if err != nil {
		pageSize = 50
	}

	start, err := strconv.Atoi(c.QueryParam("start"))

	if err != nil {
		start = 0
	}

	total, items, err := h.app.Queries.ListBeers.Handle(c.Request().Context(), query.ListBeers{
		Name:     c.QueryParam("name"),
		Start:    int64(start),
		PageSize: int64(pageSize),
	})

	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, responses.PaginatedResponse{
		Start:    int64(start),
		PageSize: int64(pageSize),
		Total:    total,
		Data:     items,
	})
}

// GetBoxPrice godoc
// @Summary Return total price.
// @Tags Beers
// @Accept json
// @Produce json
// @Param id path int64  true    "search one beer for Id"
// @Param currency query int  false  "money to pay"
// @Param quantity query int  false  "quantity"
// @Success 200 {object} responses.PaginatedResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/beers/{beerId}/boxprice [get]
func (h HttpServer) GetBoxPrice(c echo.Context) error {
	var beerId int64
	if n, err := strconv.Atoi(c.Param("beerId")); err == nil {
		beerId = int64(n)
	} else {
		panic(err)
	}

	quantity, err := strconv.Atoi(c.QueryParam("quantity"))

	if err != nil {
		quantity = 6
	}

	result, err := h.app.Queries.GetBoxPrice.Handle(c.Request().Context(), query.GetBoxPrice{
		Currency: c.QueryParam("currency"),
		Id:       beerId,
		Quantity: int64(quantity),
	})

	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)
}

func Simple(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}

	return errs
}
