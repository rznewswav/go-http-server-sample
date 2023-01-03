package mongodb

import (
	"container/list"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbServiceMock struct {
	collectionMock *MongodbCollectionMock
	nextToPop      map[string]*list.List
	nextErrorToPop map[string]*list.List
}

type MongodbCollectionMock struct {
	service *MongodbServiceMock
}

func (service *MongodbServiceMock) popNext(context string) interface{} {
	valuesToPop := service.nextToPop[context]
	if valuesToPop == nil {
		return nil
	}

	valueToReturn := valuesToPop.Front().Value
	fmt.Printf("popNext: %s\n", valueToReturn)
	return valueToReturn
}

func (service *MongodbServiceMock) popNextError(context string) error {
	valuesToPop := service.nextErrorToPop[context]
	if valuesToPop == nil {
		return nil
	}

	valueToReturn := valuesToPop.Front().Value.(error)
	fmt.Printf("popNextError: %s\n", valueToReturn)
	return valueToReturn
}

func (collection *MongodbCollectionMock) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	value := mongo.NewSingleResultFromDocument(collection.service.popNext("FindOne"), nil, nil)
	return value
}

func (collection *MongodbCollectionMock) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return collection.service.popNext("UpdateOne").(*mongo.UpdateResult), collection.service.popNextError("UpdateOne")
}

func (service *MongodbServiceMock) Init(databaseUri string, databaseName string) {}

func (service *MongodbServiceMock) WithCollection(collection string) IMongodbCollection {
	return service.collectionMock
}

func (service *MongodbServiceMock) GetClient() *mongo.Client {
	return nil
}

type popNextInstructionObject struct {
	functionName string
	value        interface{}
	err          error
}

func PopNextInstruction(functionName string, value interface{}, err error) popNextInstructionObject {
	return popNextInstructionObject{
		functionName: functionName,
		value:        value,
		err:          err,
	}
}

func PrepareMockService(popNexts ...popNextInstructionObject) IMongodbService {
	service := MongodbServiceMock{}
	collection := MongodbCollectionMock{
		service: &service,
	}

	service.collectionMock = &collection
	for _, pnio := range popNexts {
		if service.nextToPop == nil {
			service.nextToPop = map[string]*list.List{}
		}

		if service.nextToPop[pnio.functionName] == nil {
			service.nextToPop[pnio.functionName] = list.New()
		}

		service.nextToPop[pnio.functionName].PushBack(pnio.value)

		if service.nextErrorToPop == nil {
			service.nextErrorToPop = map[string]*list.List{}
		}

		if service.nextErrorToPop[pnio.functionName] == nil {
			service.nextErrorToPop[pnio.functionName] = list.New()
		}

		service.nextErrorToPop[pnio.functionName].PushBack(pnio.err)
	}
	return &service
}

func (service *MongodbServiceMock) DecodeSingleResult(singleResult *mongo.SingleResult, object interface{}) (any, error) {
	return service.popNext("Decode"), nil
}
