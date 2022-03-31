package common

import (
	"context"

	"github.com/juanmaabanto/go-ms-beers/common/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IBaseRepository interface {
	Count(ctx context.Context, filter interface{}) (int64, error)
	DeleteById(ctx context.Context, id int64) (int64, error)
	FilterBy(ctx context.Context, filter interface{}, receiver []interface{}) error
	FindById(ctx context.Context, id int64, receiver interface{}) error
	FindOne(ctx context.Context, filter interface{}, receiver interface{}) error
	InsertMany(ctx context.Context, documents []interface{}) ([]int64, error)
	InsertOne(ctx context.Context, document interface{}) (int64, error)
	Paginated(ctx context.Context, filter interface{}, sort interface{}, pageSize int64, start int64, receiver interface{}) error
	UpdateOne(ctx context.Context, id int64, document interface{}) error
}

type BaseRepository struct {
	collection mongo.Collection
}

func NewBaseRepository(connection database.MongoConnection, document IDocument) BaseRepository {
	repository := BaseRepository{
		collection: *connection.Database.Collection(document.GetCollectionName()),
	}

	return repository
}

func (repo BaseRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
	result, err := repo.collection.CountDocuments(ctx, filter)

	return result, err
}

func (repo BaseRepository) DeleteById(ctx context.Context, id int64) (int64, error) {
	result, err := repo.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})

	return result.DeletedCount, err
}

func (repo BaseRepository) FilterBy(ctx context.Context, filter interface{}, receiver []interface{}) error {
	cursor, err := repo.collection.Find(ctx, filter)

	if err != nil {
		return err
	}

	cursor.Decode(receiver)

	return nil
}

func (repo BaseRepository) FindById(ctx context.Context, id int64, receiver interface{}) error {
	coll := repo.collection
	result := coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}})

	if result.Err() != nil && result.Err() != mongo.ErrNoDocuments {
		return result.Err()
	}

	if result.Err() == mongo.ErrNoDocuments {
		return nil
	}

	return result.Decode(receiver)
}

func (repo BaseRepository) FindOne(ctx context.Context, filter interface{}, receiver interface{}) error {
	result := repo.collection.FindOne(ctx, filter)

	if result.Err() != nil && result.Err() != mongo.ErrNoDocuments {
		return result.Err()
	}

	return result.Decode(receiver)
}

func (repo BaseRepository) InsertMany(ctx context.Context, documents []interface{}) ([]int64, error) {
	result, err := repo.collection.InsertMany(ctx, documents)

	if err != nil {
		panic(err)
	}

	array := []int64{}

	for i := range result.InsertedIDs {
		array = append(array, result.InsertedIDs[i].(int64))
	}

	return array, err
}

func (repo BaseRepository) InsertOne(ctx context.Context, document interface{}) (int64, error) {
	result, err := repo.collection.InsertOne(ctx, document)

	return result.InsertedID.(int64), err
}

func (repo BaseRepository) Paginated(ctx context.Context, filter interface{}, sort interface{}, pageSize int64, start int64, receiver interface{}) error {
	options := options.Find()

	options.SetSort(sort)
	options.SetSkip(start)
	options.SetLimit(pageSize)

	cursor, err := repo.collection.Find(ctx, filter, options)

	if err != nil {
		return err
	}

	return cursor.All(ctx, receiver)
}

func (repo BaseRepository) UpdateOne(ctx context.Context, id int64, document interface{}) error {
	_, err := repo.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.M{"$set": document})

	return err
}
