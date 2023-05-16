package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserId string `json:"id"`
	jwt.StandardClaims
}
