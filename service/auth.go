package service

import (
	"auth_service/domain"
	"auth_service/repository"
	"auth_service/usecases"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type Auth struct {
	repo       repository.Auth
	secret     string
	refreshExp time.Duration
	accessExp  time.Duration
}

func NewAuth(repo repository.Auth, secret string, refreshExp, accessExp time.Duration) usecases.Auth {
	return &Auth{
		repo:       repo,
		secret:     secret,
		refreshExp: refreshExp,
		accessExp:  accessExp,
	}
}

func (s *Auth) CreateAccessToken(user *domain.User, auth *domain.Auth) (domain.AccessToken, error) {
	claims := &domain.AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(auth.Iat),
			ExpiresAt: jwt.NewNumericDate(auth.Exp),
		},
		IP: user.IP,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	t, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (s *Auth) ParseAccessToken(token domain.AccessToken) (*domain.AccessTokenClaims, error) {
	claims := &domain.AccessTokenClaims{}

	jwtToken, err := jwt.ParseWithClaims(string(token), claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if !jwtToken.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

func (s *Auth) CreateRefreshToken(user *domain.User, auth *domain.Auth) (domain.RefreshToken, error) {
	rawToken := fmt.Sprintf("%s-%s", user.UserID, user.IP)
	b := []byte(rawToken)
	token := domain.RefreshToken(base64.URLEncoding.EncodeToString(b))
	hashedToken, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	auth.RefreshHashed = hashedToken

	return token, nil
}

func (s *Auth) ParseRefreshToken(token domain.RefreshToken) (domain.User, error) {
	var user domain.User

	b, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return user, domain.ErrInvalidToken
	}

	tokenStr := domain.RefreshToken(b)
	idx := strings.LastIndex(tokenStr, "-")
	user.IP = tokenStr[idx+1:]
	user.UserID = uuid.UUID(b[:idx])
	return user, nil
}

func (s *Auth) CreateTokenPair(user *domain.User, auth *domain.Auth) (domain.AccessToken, domain.RefreshToken, error) {
	aToken, err := s.CreateAccessToken(user, auth)
	if err != nil {
		return "", "", err
	}

	rToken, err := s.CreateRefreshToken(user, auth)
	if err != nil {
		return "", "", err
	}

	auth.LastUpdate = auth.Iat
	err = s.repo.Put(*auth)
	if err != nil {
		return "", "", err
	}

	return aToken, rToken, nil
}
