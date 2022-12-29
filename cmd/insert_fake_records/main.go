package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	brcypt "golang.org/x/crypto/bcrypt"

	auth "newswav/http-server-sample/modules/auth"
	"newswav/http-server-sample/services/utils"

	database "newswav/http-server-sample/modules/mongodb"
)

func main() {
	fake := faker.New()

	mongodbService := database.MongodbService{}
	mongodbService.Init("mongodb://127.0.0.1:27017/", "golang-poc")

	task := func(cancel *func()) {
		var user auth.SchemaUser
		collection := mongodbService.WithCollection(user.SchemaName())
		fmt.Println("Initializing user...")

		user.Name = fake.Person().FirstName()
		if hashedP, err := brcypt.GenerateFromPassword([]byte("abcdef"), 10); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n%s\n", err.Error(), debug.Stack())
			(*cancel)()
			return
		} else {
			user.HashedP = string(hashedP)
		}

		fakeContact := fake.Person().Contact()
		user.ContactInfo.Email = fakeContact.Email
		user.ContactInfo.PhoneNumber = fakeContact.Phone

		fmt.Println("Upserting user...")

		upsertResult, err := collection.UpdateOne(
			context.Background(),
			bson.M{
				"contactInfo.email": user.ContactInfo.Email,
			},
			bson.M{"$set": user},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n%s\n", err.Error(), debug.Stack())
			(*cancel)()
			return
		}

		document := database.CompileDocument(upsertResult.UpsertedID, &user)

		if b, err := json.MarshalIndent(document, "", "  "); err != nil {
			println("Unable to parse config into JSON:", err.Error())
		} else {
			println("Upserted one record:", string(b))
		}
	}
	utils.Repeat(100, &task)
}
