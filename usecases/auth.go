package usecases

import (
	"auth_service/domain"
	pkgtime "auth_service/pkg/time"
)

type Auth interface {
	CreateTokenPair(user *domain.User) (domain.AccessToken, domain.RefreshToken, error)
	RefreshTokenPair(refreshToken domain.RefreshToken, IP string) (domain.AccessToken, domain.RefreshToken, error)
	CreateAccessToken(user *domain.User, auth *domain.Auth) (domain.AccessToken, error)
	CreateRefreshToken(user *domain.User, auth *domain.Auth) (domain.RefreshToken, error)
	ParseAccessToken(token domain.AccessToken) (*domain.AccessTokenClaims, error)
	ParseRefreshToken(token domain.RefreshToken) (domain.User, pkgtime.UnixTime, error)
}
