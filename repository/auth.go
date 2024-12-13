package repository

import (
	"auth_service/domain"
	"github.com/google/uuid"
)

type Auth interface {
	GetByUUID(uuid uuid.UUID) (domain.Auth, error)
	GetByRefreshHashed(token []byte) (domain.Auth, error)
	Put(auth domain.Auth) error
}
