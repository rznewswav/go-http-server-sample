package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongodbService interface {
	Init(databaseUri string, databaseName string)
	WithCollection(collection string) IMongodbCollection
	GetClient() *mongo.Client
	DecodeSingleResult(singleResult *mongo.SingleResult, object interface{}) (interface{}, error)
}

type MongodbService struct {
	Client       *mongo.Client
	databaseName string
}

type IMongodbCollection interface {
	FindOne(ctx context.Context, filter interface{},
		opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
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

func (service *MongodbService) WithCollection(collection string) IMongodbCollection {
	collectionHandler := service.Client.Database(service.databaseName).Collection(collection)
	return collectionHandler
}

func (service *MongodbService) GetClient() *mongo.Client {
	return service.Client
}

func (service *MongodbService) DecodeSingleResult(singleResult *mongo.SingleResult, object interface{}) (interface{}, error) {
	err := singleResult.Decode(object)
	return object, err
}
