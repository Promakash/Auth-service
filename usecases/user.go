package usecases

import (
	"auth_service/domain"
	"github.com/google/uuid"
)

type User interface {
	GenTokens(uuid uuid.UUID, IP string) (domain.AccessToken, domain.RefreshToken, error)
	RefreshTokens(token domain.RefreshToken, IP string) (domain.AccessToken, domain.RefreshToken, error)
	NotifyUser(uuid uuid.UUID) error
}
