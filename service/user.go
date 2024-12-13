package service

import (
	"auth_service/domain"
	"auth_service/repository"
	"auth_service/usecases"
	"github.com/google/uuid"
)

type User struct {
	repo repository.User
}

func NewUser(repo repository.User) usecases.User {
	return &User{
		repo: repo,
	}
}

func (s *User) GetByUUID(uuid uuid.UUID) (domain.User, error) {
	return s.repo.GetByUUID(uuid)
}
func (s *User) NotifyUser(user domain.User) error {
	return nil
}
