package repository

import (
	"auth_service/domain"
	"github.com/google/uuid"
)

type User interface {
	Put(user domain.User) error
	GetByUUID(uuid uuid.UUID) (domain.User, error)
}
