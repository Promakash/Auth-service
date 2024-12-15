package service

import (
	"auth_service/domain"
	pkgtime "auth_service/pkg/time"
	"auth_service/repository"
	"auth_service/usecases"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type Auth struct {
	repo            repository.Auth
	secret          string
	refreshTokenExp pkgtime.UnixTime
	accessTokenExp  pkgtime.UnixTime
}

func NewAuth(repo repository.Auth, secret string, refreshExp, accessExp time.Duration) usecases.Auth {
	return &Auth{
		repo:            repo,
		secret:          secret,
		refreshTokenExp: pkgtime.DaysToUnix(refreshExp),
		accessTokenExp:  pkgtime.DaysToUnix(accessExp),
	}
}

func (s *Auth) CreateAccessToken(user *domain.User, auth *domain.Auth) (domain.AccessToken, error) {
	claims := &domain.AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Unix(auth.Iat, 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(auth.Iat+s.accessTokenExp, 0)),
		},
		IP: user.IP,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	t, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func (s *Auth) ParseAccessToken(token domain.AccessToken) (*domain.AccessTokenClaims, error) {
	claims := &domain.AccessTokenClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
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
	rawToken := fmt.Sprintf("%s/%s/%d", user.UserID, user.IP, auth.Iat)
	b := []byte(rawToken)
	token := base64.URLEncoding.EncodeToString(b)
	hashedToken, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	auth.RefreshHashed = hashedToken

	return token, nil
}

func (s *Auth) ParseRefreshToken(token domain.RefreshToken) (domain.User, pkgtime.UnixTime, error) {
	const rTokenParts = 3
	var user domain.User
	var err error

	decodedToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return user, 0, err
	}

	parts := strings.Split(domain.RefreshToken(decodedToken), "/")
	if len(parts) != rTokenParts {
		return user, 0, domain.ErrInvalidToken
	}

	user.UserID, err = uuid.Parse(parts[0])
	if err != nil {
		return user, 0, err
	}

	user.IP = parts[1]

	lastUpdate, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return user, 0, err
	}
	return user, lastUpdate, nil
}

func (s *Auth) CreateTokenPair(user *domain.User) (domain.AccessToken, domain.RefreshToken, error) {
	auth := domain.Auth{
		UserID: user.UserID,
	}
	s.setRefreshTokenTimestamps(&auth)

	aToken, err := s.CreateAccessToken(user, &auth)
	if err != nil {
		return "", "", err
	}

	rToken, err := s.CreateRefreshToken(user, &auth)
	if err != nil {
		return "", "", err
	}

	err = s.repo.Put(auth)
	if err != nil {
		return "", "", err
	}

	return aToken, rToken, nil
}

func (s *Auth) RefreshTokenPair(refreshToken domain.RefreshToken, IP string) (domain.AccessToken, domain.RefreshToken, error) {
	userData, iat, err := s.ParseRefreshToken(refreshToken)
	if err != nil {
		slog.Warn("Error occured while parsing refresh token", "token", refreshToken, "error", err)
		return "", "", domain.ErrInvalidToken
	}
	userData.IP = IP

	auth, err := s.repo.GetByUUID(userData.UserID)
	if err != nil {
		return "", "", err
	}

	verified, err := s.verifyToken(refreshToken, iat, &auth)
	if err != nil {
		return "", "", domain.ErrUnauthorized
	}
	if verified != true {
		return "", "", domain.ErrInvalidToken
	}

	s.setRefreshTokenTimestamps(&auth)

	aToken, err := s.CreateAccessToken(&userData, &auth)
	if err != nil {
		return "", "", err
	}

	rToken, err := s.CreateRefreshToken(&userData, &auth)
	if err != nil {
		return "", "", err
	}

	err = s.repo.Put(auth)
	if err != nil {
		return "", "", err
	}

	return aToken, rToken, nil
}

func (s *Auth) verifyToken(refreshToken domain.RefreshToken, iat pkgtime.UnixTime, auth *domain.Auth) (bool, error) {
	if iat != auth.Iat {
		return false, domain.ErrUnauthorized
	}
	if iat+s.refreshTokenExp < time.Now().Unix() {
		return false, domain.ErrUnauthorized
	}

	decodedToken, err := base64.URLEncoding.DecodeString(refreshToken)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(auth.RefreshHashed, decodedToken)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Auth) setRefreshTokenTimestamps(auth *domain.Auth) {
	timestamp := time.Now().Unix()
	auth.Iat = timestamp
	auth.Exp = timestamp + s.refreshTokenExp
}
