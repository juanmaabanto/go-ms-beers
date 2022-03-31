package query

import (
	"context"
	"errors"
	"testing"

	"github.com/juanmaabanto/go-ms-beers/common/mocks"
	"github.com/juanmaabanto/go-ms-beers/internal/domain/beer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_NewListBeersHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic")
		}
	}()

	NewListBeersHandler(nil)
}

func Test_Handle_ListBeers_Count_Err(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()

	mockRepo.On("Count", ctx, mock.AnythingOfType("primitive.D")).Return(int64(0), errors.New("An error"))

	// Act
	testQuery := NewListBeersHandler(mockRepo)
	_, _, err := testQuery.Handle(ctx, ListBeers{Name: "", Start: 0, PageSize: 50})

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.Equal(t, "An error", err.Error())
}

func Test_Handle_ListBeers_Paginated_Error(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()

	mockRepo.On("Count", ctx, mock.AnythingOfType("primitive.D")).Return(int64(1), nil)
	mockRepo.On("Paginated", ctx, mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D"), int64(50), int64(0), mock.AnythingOfType("*[]beer.Beer")).Return(errors.New("An error"))

	// Act
	testQuery := NewListBeersHandler(mockRepo)
	_, _, err := testQuery.Handle(ctx, ListBeers{Name: "", Start: 0, PageSize: 50})

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.Equal(t, "An error", err.Error())
}

func Test_Handle_ListBeers_Ok(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockRepository)
	ctx := context.Background()

	mockRepo.On("Count", ctx, mock.AnythingOfType("primitive.D")).Return(int64(1), nil)
	mockRepo.On("Paginated", ctx, mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D"), int64(50), int64(0), mock.AnythingOfType("*[]beer.Beer")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(5).(*[]beer.Beer)

		*arg = append(*arg, beer.Beer{
			Name: "test",
		})

	})

	// Act
	testQuery := NewListBeersHandler(mockRepo)
	total, results, _ := testQuery.Handle(ctx, ListBeers{Name: "", Start: 0, PageSize: 50})

	// Assert
	mockRepo.AssertExpectations(t)

	assert.Equal(t, int64(1), total)
	assert.Equal(t, "test", results[0].Name)
}
