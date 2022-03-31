package command

import (
	"context"
	errorsN "errors"
	"testing"

	"github.com/juanmaabanto/go-ms-beers/common/errors"
	"github.com/juanmaabanto/go-ms-beers/common/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_NewCreateBeerHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic")
		}
	}()

	NewCreateBeerHandler(nil)
}

func Test_Handle_CreateBeer_Count_Error(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()
	item := CreateBeer{Name: "test"}

	mockRepo.On("Count", ctx, mock.AnythingOfType("primitive.M")).Return(int64(0), errorsN.New("An error has occurred"))

	// Act
	testCommand := NewCreateBeerHandler(mockRepo)
	_, err := testCommand.Handle(ctx, item)

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "An error has occurred")
}

func Test_Handle_CreateBeer_Count_Is_Greater_Than_Zero(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()
	item := CreateBeer{Name: "test"}

	mockRepo.On("Count", ctx, mock.AnythingOfType("primitive.M")).Return(int64(1), nil)

	// Act
	testCommand := NewCreateBeerHandler(mockRepo)
	_, err := testCommand.Handle(ctx, item)

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.IsType(t, errors.ApplicationError{}, err)
	assert.Equal(t, errors.ErrorTypeConflict, err.(errors.ApplicationError).ErrorType())
}

func Test_Handle_CreateBeer_Insert_Completed(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()
	item := CreateBeer{Name: "test"}
	expected := int64(1)

	mockRepo.On("Count", ctx, mock.AnythingOfType("primitive.M")).Return(int64(0), nil)
	mockRepo.On("InsertOne", ctx, mock.AnythingOfType("beer.Beer")).Return(expected, nil)

	// Act
	testCommand := NewCreateBeerHandler(mockRepo)
	result, err := testCommand.Handle(ctx, item)

	// Assert

	mockRepo.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_Handle_CreateBeer_Insert_Error(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()
	item := CreateBeer{Name: "test"}

	mockRepo.On("Count", ctx, mock.AnythingOfType("primitive.M")).Return(int64(0), nil)
	mockRepo.On("InsertOne", ctx, mock.AnythingOfType("beer.Beer")).Return(int64(0), errorsN.New("An error has occurred"))

	// Act
	testCommand := NewCreateBeerHandler(mockRepo)
	_, err := testCommand.Handle(ctx, item)

	// Assert

	mockRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "An error has occurred")
}
