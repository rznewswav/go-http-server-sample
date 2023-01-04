package tests

import (
	"fmt"
	"newswav/http-server-sample/modules/auth"
	"newswav/http-server-sample/modules/auth/mocks"
	"os"
	"runtime/debug"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

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

	mockRepo := mocks.GetMockUserRepository(&auth.SchemaUser{
		Name:    "Name",
		HashedP: hashedPassword,
		ContactInfo: auth.SchemaUserContactInfo{
			PhoneNumber: "PhoneNumber",
			Email:       "Email",
		},
	})

	service := auth.AuthService{
		UserRepo:  mockRepo,
		JWTSecret: "abcdef",
	}

	payload := service.ValidateLogin("Email", password)

	if payload == nil {
		t.Fail()
	}
}
