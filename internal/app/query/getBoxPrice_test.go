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

func Test_NewGetBoxPriceHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic")
		}
	}()

	NewGetBoxPriceHandler(nil)
}

func Test_Handle_GetBoxPrice_Found(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()
	expected := response.PriceResponse{PriceTotal: 10}

	mockRepo.On("FindById", ctx, int64(1), mock.AnythingOfType("*beer.Beer")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*beer.Beer)
		arg.Id = 1
		arg.Name = "test"
		arg.Price = 10
	})

	// Act
	testQuery := NewGetBoxPriceHandler(mockRepo)
	result, err := testQuery.Handle(ctx, GetBoxPrice{Id: 1, Quantity: 1})

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Equal(t, expected.PriceTotal, result.PriceTotal)
	assert.Nil(t, err)
}

func Test_Handler_GetBoxPrice_Not_Found(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()

	mockRepo.On("FindById", ctx, mock.AnythingOfType("int64"), mock.AnythingOfType("*beer.Beer")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*beer.Beer)
		arg.Id = 0
	})

	// Act
	testQuery := NewGetBoxPriceHandler(mockRepo)
	_, err := testQuery.Handle(ctx, GetBoxPrice{Id: 1})

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.IsType(t, errors.ApplicationError{}, err)
	assert.Equal(t, errors.ErrorTypeNotFound, err.(errors.ApplicationError).ErrorType())
}

func Test_Handler_GetBoxPrice_Error(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()

	mockRepo.On("FindById", ctx, mock.AnythingOfType("int64"), mock.AnythingOfType("*beer.Beer")).Return(errorsN.New("An error has occurred"))

	// Act
	testQuery := NewGetBoxPriceHandler(mockRepo)
	_, err := testQuery.Handle(ctx, GetBoxPrice{Id: 1})

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "An error has occurred")
}

func Test_convertCurrent(t *testing.T) {
	type args struct {
		source string
		target string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "EUR to PEN",
			args:    args{source: "EUR", target: "PEN"},
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertCurrent(tt.args.source, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertCurrent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got <= 0 {
				t.Errorf("convertCurrent() = %v, want %v", got, tt.want)
			}
		})
	}
}
