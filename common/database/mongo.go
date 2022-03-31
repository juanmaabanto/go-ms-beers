package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConnection struct {
	Client   mongo.Client
	Database mongo.Database
}

func NewMongoConnection(ctx context.Context, databaseName string, uri string) MongoConnection {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	return MongoConnection{
		Client:   *client,
		Database: *client.Database(databaseName),
	}
}
