package auth_test

import (
	"context"
	"fmt"
	"newswav/http-server-sample/modules/auth"
	"newswav/http-server-sample/modules/mongodb"
	"os"
	"runtime/debug"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type MockReturnValue struct {
	hello string
}

func TestService(t *testing.T) {
	database := mongodb.PrepareMockService(
		mongodb.PopNextInstruction("Decode", MockReturnValue{
			hello: "world",
		}, nil),
	)

	service := auth.AuthService{
		Database:  database,
		JWTSecret: "abcdef",
	}

	collection := service.Database.WithCollection("nil")

	singleResult := collection.FindOne(context.TODO(), nil, nil)
	expected, err := service.Database.DecodeSingleResult(singleResult, nil)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if expected.(MockReturnValue).hello != "world" {
		t.Fail()
	} else {
		return
	}
}

func TestLogin(t *testing.T) {
	var hashedPassword string
	password := "abcdef"
	if hashedP, err := bcrypt.GenerateFromPassword([]byte(password), 10); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n%s\n", err.Error(), debug.Stack())
		t.Fail()
		return
	} else {
		hashedPassword = string(hashedP)
	}

	database := mongodb.PrepareMockService(
		mongodb.PopNextInstruction("Decode", &auth.SchemaUser{
			Name:    "Name",
			HashedP: hashedPassword,
			ContactInfo: auth.SchemaUserContactInfo{
				PhoneNumber: "PhoneNumber",
				Email:       "Email",
			},
		}, nil),
	)

	service := auth.AuthService{
		Database:  database,
		JWTSecret: "abcdef",
	}

	payload := service.ValidateLogin("Email", password)

	if payload == nil {
		t.Fail()
	}
}
