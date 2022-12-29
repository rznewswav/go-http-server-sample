package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbService struct {
	Client       *mongo.Client
	databaseName string
}

func (service *MongodbService) Init(databaseUri string, databaseName string) {
	service.databaseName = databaseName
	ctx, cancel := context.WithTimeout(context.Background(), 30000)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseUri))
	if err != nil {
		panic(err)
	}
	service.Client = client
}

func (service *MongodbService) WithCollection(collection string) *mongo.Collection {
	collectionHandler := service.Client.Database(service.databaseName).Collection(collection)
	return collectionHandler
}
