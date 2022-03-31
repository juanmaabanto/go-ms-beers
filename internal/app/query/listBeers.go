package query

import (
	"context"

	"github.com/juanmaabanto/go-ms-beers/internal/domain/beer"
	"github.com/juanmaabanto/go-ms-beers/internal/ports/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListBeers struct {
	Name     string
	Start    int64
	PageSize int64
}

type ListBeersHandler struct {
	repo beer.Repository
}

func NewListBeersHandler(repo beer.Repository) ListBeersHandler {
	if repo == nil {
		panic("nil repo")
	}

	return ListBeersHandler{repo}
}

func (h ListBeersHandler) Handle(ctx context.Context, query ListBeers) (int64, []response.BeerResponse, error) {
	var items []beer.Beer
	results := []response.BeerResponse{}

	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"name", primitive.Regex{
					Pattern: query.Name,
					Options: "i",
				}}},
			}},
	}

	total, err := h.repo.Count(ctx, filter)

	if err != nil {
		return 0, results, err
	}

	err = h.repo.Paginated(ctx, filter, bson.D{}, query.PageSize, query.Start, &items)

	if err != nil {
		return 0, results, err
	}

	for _, element := range items {
		results = append(results, response.BeerResponse{
			Id:       element.Id,
			Name:     element.Name,
			Brewery:  element.Brewery,
			Country:  element.Country,
			Price:    element.Price,
			Currency: element.Currency,
		})
	}

	return total, results, err
}
