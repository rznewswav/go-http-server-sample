package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	ValidateLogin(email, password string) *JWTPayload
	ValidateToken(jwtToken string) (bool, error)
	SignJWT(payload JWTPayload) (string, error)
}

type AuthService struct {
	UserRepo  IUserRepository
	JWTSecret string
}

func (service *AuthService) ValidateLogin(email, password string) *JWTPayload {
	user, err := service.UserRepo.GetUserByEmail(email)

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
