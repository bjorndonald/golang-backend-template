package helpers

import "github.com/golang-jwt/jwt"

type OtpVerify struct {
	Token string `json:"token" validate:"required,len=5"`
	Email string `json:"email" validate:"required,email"`
}

type IError struct {
	Field string
	Tag   string
	Value string
}

type AuthTokenJwtClaim struct {
	Email  string
	Name   string
	UserId string
	jwt.StandardClaims
}

type AccountStatus int

type EmailInput struct {
	Email string `json:"email" validate:"required,email"`
}
