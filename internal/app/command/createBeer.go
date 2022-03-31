package command

import (
	"context"
	"time"

	"github.com/juanmaabanto/go-ms-beers/common/errors"
	"github.com/juanmaabanto/go-ms-beers/internal/domain/beer"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateBeer struct {
	Id       int64   `json:"id" validate:"required"`
	Name     string  `json:"name" validate:"required,max=30"`
	Brewery  string  `json:"brewery" validate:"required,max=30"`
	Country  string  `json:"country" validate:"required,max=20"`
	Price    float64 `json:"price" validate:"required"`
	Currency string  `json:"currency" validate:"required,max=5"`
}

type CreateBeerHandler struct {
	repo beer.Repository
}

func NewCreateBeerHandler(repo beer.Repository) CreateBeerHandler {
	if repo == nil {
		panic("nil repo beer")
	}

	return CreateBeerHandler{repo: repo}
}

func (h CreateBeerHandler) Handle(ctx context.Context, command CreateBeer) (int64, error) {
	count, err := h.repo.Count(ctx, bson.M{"_id": command.Id})

	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, errors.NewConflictError("An element with the same id already exists.")
	}

	item := beer.Beer{}
	item.Id = command.Id
	item.Name = command.Name
	item.Brewery = command.Brewery
	item.Country = command.Country
	item.Price = command.Price
	item.Currency = command.Currency
	item.CreatedAt = time.Now()
	item.CreatedBy = "admin"

	id, err := h.repo.InsertOne(ctx, item)

	if err != nil {
		return id, err
	}

	return id, nil
}
