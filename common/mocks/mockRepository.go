package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mock MockRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
	args := mock.Called(ctx, filter)
	result := args.Get(0)

	return result.(int64), args.Error(1)
}

func (mock MockRepository) DeleteById(ctx context.Context, id int64) (int64, error) {
	args := mock.Called(ctx, id)
	result := args.Get(0)

	return result.(int64), args.Error(1)
}

func (mock MockRepository) FilterBy(ctx context.Context, filter interface{}, receiver []interface{}) error {
	args := mock.Called()

	return args.Error(0)
}

func (mock MockRepository) FindById(ctx context.Context, id int64, receiver interface{}) error {
	args := mock.Called(ctx, id, receiver)

	return args.Error(0)
}

func (mock MockRepository) FindOne(ctx context.Context, filter interface{}, receiver interface{}) error {
	args := mock.Called()

	return args.Error(0)
}

func (mock MockRepository) InsertMany(ctx context.Context, documents []interface{}) ([]int64, error) {
	args := mock.Called()
	result := args.Get(0)

	return result.([]int64), args.Error(1)
}

func (mock MockRepository) InsertOne(ctx context.Context, document interface{}) (int64, error) {
	args := mock.Called(ctx, document)
	result := args.Get(0)

	return result.(int64), args.Error(1)
}

func (mock MockRepository) Paginated(ctx context.Context, filter interface{}, sort interface{}, pageSize int64, start int64, receiver interface{}) error {
	args := mock.Called(ctx, filter, sort, pageSize, start, receiver)

	return args.Error(0)
}

func (mock MockRepository) UpdateOne(ctx context.Context, id int64, document interface{}) error {
	args := mock.Called(ctx, id, document)

	return args.Error(0)
}
