package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	IP string `json:"ip"`
}

type AccessToken = string

type RefreshToken = string

type Auth struct {
	UserID        uuid.UUID
	RefreshHashed []byte
	Iat           int64
	Exp           int64
}
