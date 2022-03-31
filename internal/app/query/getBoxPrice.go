package query

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/juanmaabanto/go-ms-beers/common/errors"
	"github.com/juanmaabanto/go-ms-beers/internal/domain/beer"
	"github.com/juanmaabanto/go-ms-beers/internal/ports/response"
)

type GetBoxPrice struct {
	Id       int64
	Currency string
	Quantity int64
}

type GetBoxPriceHandler struct {
	repo beer.Repository
}

func NewGetBoxPriceHandler(repo beer.Repository) GetBoxPriceHandler {
	if repo == nil {
		panic("nil repo")
	}

	return GetBoxPriceHandler{repo}
}

func (h GetBoxPriceHandler) Handle(ctx context.Context, query GetBoxPrice) (*response.PriceResponse, error) {
	receiver := beer.Beer{}

	err := h.repo.FindById(ctx, query.Id, &receiver)

	if err != nil {
		return nil, err
	}

	if receiver.Id == 0 {
		return nil, errors.NewNotFoundError("beer")
	}

	endPrice := 0.00

	if len(query.Currency) < 1 || query.Currency == receiver.Currency {
		endPrice = receiver.Price
	} else {
		ratio, err := convertCurrent(receiver.Currency, query.Currency)

		if err != nil {
			return nil, err
		}
		endPrice = receiver.Price * ratio
	}

	response := response.PriceResponse{
		PriceTotal: endPrice * float64(query.Quantity),
	}

	return &response, nil

}

func convertCurrent(source string, target string) (float64, error) {
	client := &http.Client{}
	uri := "http://api.currencylayer.com/live?access_key=b17cd86be31831425bf635f9e8be3f26&currencies=" + source + "," + target

	req, err := http.NewRequest(http.MethodGet, uri, nil)

	if err != nil {
		return 0, err
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	c := make(map[string]json.RawMessage)

	err = json.Unmarshal(body, &c)

	if err != nil {
		return 0, nil
	}

	results := make(map[string]json.RawMessage)

	err = json.Unmarshal(c["quotes"], &results)

	if err != nil {
		return 0, nil
	}

	m1, err := strconv.ParseFloat(string(results["USD"+source]), 64)

	if err != nil {
		return 0, nil
	}

	m2, err := strconv.ParseFloat(string(results["USD"+target]), 64)

	if err != nil {
		return 0, nil
	}

	return (1.00 / m1) * m2, nil
}

