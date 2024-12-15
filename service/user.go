package service

import (
	"auth_service/domain"
	"auth_service/repository"
	"auth_service/usecases"
	"github.com/google/uuid"
)

type User struct {
	authService usecases.Auth
	repo        repository.User
}

func NewUser(authService usecases.Auth, userRepo repository.User) usecases.User {
	return &User{
		authService: authService,
		repo:        userRepo,
	}
}

func (s *User) GenTokens(uuid uuid.UUID, IP string) (domain.AccessToken, domain.RefreshToken, error) {
	user, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return "", "", err
	}
	user.IP = IP

	aToken, rToken, err := s.authService.CreateTokenPair(&user)
	if err != nil {
		return "", "", err
	}

	err = s.repo.Put(user)
	if err != nil {
		return "", "", err
	}

	return aToken, rToken, nil
}

func (s *User) RefreshTokens(token domain.RefreshToken, IP string) (domain.AccessToken, domain.RefreshToken, error) {
	aToken, rToken, err := s.authService.RefreshTokenPair(token, IP)
	if err != nil {
		return "", "", err
	}
	userData, _, err := s.authService.ParseRefreshToken(token)
	if err != nil {
		return "", "", err
	}
	if userData.IP != IP {
		err = s.NotifyUser(userData.UserID)
	}
	return aToken, rToken, err
}

func (s *User) NotifyUser(uuid uuid.UUID) error {
	return nil
}
