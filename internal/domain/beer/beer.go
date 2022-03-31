package beer

import "github.com/juanmaabanto/go-ms-beers/common"

type Beer struct {
	Id              int64   `bson:"_id"`
	Name            string  `bson:"name"`
	Brewery         string  `bson:"brewery"`
	Country         string  `bson:"country"`
	Price           float64 `bson:"price"`
	Currency        string  `bson:"currency"`
	common.Document `bson:"inline"`
}

func (_ Beer) GetCollectionName() string {
	return "beer"
}
