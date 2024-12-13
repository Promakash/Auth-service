package usecases

import (
	"auth_service/domain"
)

type Auth interface {
	CreateTokenPair(user *domain.User, auth *domain.Auth) (domain.AccessToken, domain.RefreshToken, error)
	CreateAccessToken(user *domain.User, auth *domain.Auth) (domain.AccessToken, error)
	CreateRefreshToken(user *domain.User, auth *domain.Auth) (domain.RefreshToken, error)
	ParseAccessToken(token domain.AccessToken) (*domain.AccessTokenClaims, error)
	ParseRefreshToken(token domain.RefreshToken) (domain.User, error)
}
