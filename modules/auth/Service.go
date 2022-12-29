package auth

import (
	"context"
	"errors"
	"fmt"
	database "newswav/http-server-sample/modules/mongodb"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Database  *database.MongodbService
	JWTSecret string
}

func (service *AuthService) ValidateLogin(email, password string) *JWTPayload {
	var user SchemaUser
	collection := service.Database.WithCollection(user.SchemaName())
	err := collection.FindOne(
		context.Background(),
		bson.M{
			"contactInfo.email": email,
		},
	).Decode(&user)

	if err != nil {
		return nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedP), []byte(password))
	if err != nil {
		return nil
	}

	payload := JWTPayload{}
	payload.Email = user.ContactInfo.Email

	return &payload
}

func (service *AuthService) ValidateToken(jwtToken string) (bool, error) {
	if len(service.JWTSecret) == 0 {
		return false, errors.New("empty token secret")
	}
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unsupported signing method: %v", token.Header["alg"])
		}

		return []byte(service.JWTSecret), nil
	})

	if err != nil {
		// all error from jwt.Parse are jwt string related error
		return false, nil
	}

	return token.Valid, nil
}

func (service *AuthService) SignJWT(payload JWTPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload.ToJWTClaims())

	tokenString, err := token.SignedString([]byte(service.JWTSecret))

	return tokenString, err
}