package service

import (
	"context"
	"os"

	"github.com/juanmaabanto/go-ms-beers/common/database"
	"github.com/juanmaabanto/go-ms-beers/internal/app"
	"github.com/juanmaabanto/go-ms-beers/internal/app/command"
	"github.com/juanmaabanto/go-ms-beers/internal/app/query"
	"github.com/juanmaabanto/go-ms-beers/internal/domain/beer"
	"github.com/juanmaabanto/go-ms-beers/internal/infrastructure"
)

func NewApplication(ctx context.Context) app.Application {
	conn := database.NewMongoConnection(ctx, os.Getenv("MONGODB_NAME"), os.Getenv("MONGODB_URI"))
	document := new(beer.Beer)

	beerRepository := infrastructure.NewBeerRepository(conn, *document)

	return app.Application{
		Commands: app.Commands{
			CreateBeer: command.NewCreateBeerHandler(beerRepository),
		},
		Queries: app.Queries{
			GetBeerById: query.NewGetBeerByIdHandler(beerRepository),
			ListBeers:   query.NewListBeersHandler(beerRepository),
			GetBoxPrice: query.NewGetBoxPriceHandler(beerRepository),
		},
	}
}
