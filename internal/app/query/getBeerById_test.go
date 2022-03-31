package query

import (
	"context"
	errorsN "errors"
	"testing"

	"github.com/juanmaabanto/go-ms-beers/common/errors"
	"github.com/juanmaabanto/go-ms-beers/common/mocks"
	"github.com/juanmaabanto/go-ms-beers/internal/domain/beer"
	"github.com/juanmaabanto/go-ms-beers/internal/ports/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_NewGetBeerByIdHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic")
		}
	}()

	NewGetBeerByIdHandler(nil)
}

func Test_Handle_GetBeerById_Found(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()
	expected := response.BeerResponse{Id: 1, Name: "test"}

	mockRepo.On("FindById", ctx, expected.Id, mock.AnythingOfType("*beer.Beer")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*beer.Beer)
		arg.Id = expected.Id
		arg.Name = expected.Name
	})

	// Act
	testQuery := NewGetBeerByIdHandler(mockRepo)
	result, err := testQuery.Handle(ctx, GetBeerById{Id: expected.Id})

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Equal(t, expected.Id, result.Id)
	assert.Equal(t, expected.Name, result.Name)
	assert.Nil(t, err)
}

func Test_Handler_GetBeerById_Not_Found(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()

	mockRepo.On("FindById", ctx, mock.AnythingOfType("int64"), mock.AnythingOfType("*beer.Beer")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*beer.Beer)
		arg.Id = 0
	})

	// Act
	testQuery := NewGetBeerByIdHandler(mockRepo)
	_, err := testQuery.Handle(ctx, GetBeerById{Id: 1})

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.IsType(t, errors.ApplicationError{}, err)
	assert.Equal(t, errors.ErrorTypeNotFound, err.(errors.ApplicationError).ErrorType())
}

func Test_Handler_GetBeerById_Error(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()

	mockRepo.On("FindById", ctx, mock.AnythingOfType("int64"), mock.AnythingOfType("*beer.Beer")).Return(errorsN.New("An error has occurred"))

	// Act
	testQuery := NewGetBeerByIdHandler(mockRepo)
	_, err := testQuery.Handle(ctx, GetBeerById{Id: 1})

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "An error has occurred")
}
