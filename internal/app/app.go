package app

import (
	"github.com/juanmaabanto/go-ms-beers/internal/app/command"
	"github.com/juanmaabanto/go-ms-beers/internal/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateBeer command.CreateBeerHandler
}

type Queries struct {
	GetBeerById query.GetBeerByIdHandler
	ListBeers   query.ListBeersHandler
	GetBoxPrice query.GetBoxPriceHandler
}
