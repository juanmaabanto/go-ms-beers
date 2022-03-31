package infrastructure

import (
	"github.com/juanmaabanto/go-ms-beers/common"
	"github.com/juanmaabanto/go-ms-beers/common/database"
	"github.com/juanmaabanto/go-ms-beers/internal/domain/beer"
)

type BeerRepository struct {
	common.BaseRepository
}

func NewBeerRepository(connection database.MongoConnection, document beer.Beer) BeerRepository {
	repository := BeerRepository{
		BaseRepository: common.NewBaseRepository(connection, document),
	}

	return repository
}
