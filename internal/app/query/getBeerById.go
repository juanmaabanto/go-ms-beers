package query

import (
	"context"

	"github.com/juanmaabanto/go-ms-beers/common/errors"
	"github.com/juanmaabanto/go-ms-beers/internal/domain/beer"
	"github.com/juanmaabanto/go-ms-beers/internal/ports/response"
)

type GetBeerById struct {
	Id int64
}

type GetBeerByIdHandler struct {
	repo beer.Repository
}

func NewGetBeerByIdHandler(repo beer.Repository) GetBeerByIdHandler {
	if repo == nil {
		panic("nil repo")
	}

	return GetBeerByIdHandler{repo}
}

func (h GetBeerByIdHandler) Handle(ctx context.Context, query GetBeerById) (*response.BeerResponse, error) {
	receiver := beer.Beer{}

	err := h.repo.FindById(ctx, query.Id, &receiver)

	if err != nil {
		return nil, err
	}

	if receiver.Id == 0 {
		return nil, errors.NewNotFoundError("beer")
	}

	response := response.BeerResponse{
		Id:       receiver.Id,
		Name:     receiver.Name,
		Brewery:  receiver.Brewery,
		Country:  receiver.Country,
		Price:    receiver.Price,
		Currency: receiver.Currency,
	}

	return &response, nil
}
