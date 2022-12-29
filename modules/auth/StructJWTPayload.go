package auth

import "github.com/golang-jwt/jwt/v4"

type JWTPayload struct {
	Email string
}

func (payload *JWTPayload) ToJWTClaims() *jwt.MapClaims {
	return &jwt.MapClaims{
		"Email": payload.Email,
	}
}
