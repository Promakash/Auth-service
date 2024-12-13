package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
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
	Iat           time.Time
	Exp           time.Time
	LastUpdate    time.Time
}
