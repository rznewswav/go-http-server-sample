package auth

import (
	"context"
	"newswav/http-server-sample/modules/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

type IUserRepository interface {
	GetUserByEmail(email string) (SchemaUser, error)
}

type UserRepository struct {
	Database *mongodb.MongodbService
}

func (repo *UserRepository) GetUserByEmail(email string) (SchemaUser, error) {
	var userTmp SchemaUser
	collection := repo.Database.WithCollection(userTmp.SchemaName())
	err := collection.FindOne(
		context.Background(),
		bson.M{
			"contactInfo.email": email,
		},
	).Decode(&userTmp)

	return userTmp, err
}
