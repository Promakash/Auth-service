package usecases

import (
	"auth_service/domain"
	"github.com/google/uuid"
)

type User interface {
	GetByUUID(uuid uuid.UUID) (domain.User, error)
	NotifyUser(user domain.User) error
}
