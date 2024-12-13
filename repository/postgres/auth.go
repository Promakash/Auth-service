package postgres

import (
	"auth_service/domain"
	"auth_service/repository"
	"database/sql"
	"github.com/google/uuid"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) repository.Auth {
	return &AuthRepository{
		DB: db,
	}
}

func (repo *AuthRepository) Put(auth domain.Auth) error {
	query := `
		INSERT INTO users_tokens (user_id, refresh_token_hashed, issued_at, expires_at, last_time_updated)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id)
		DO UPDATE SET refresh_token_hashed = $2, 
			issued_at = $3, 
			expires_at = $4, 
			last_time_updated = $5`
	_, err := repo.DB.Exec(query, auth.UserID, auth.RefreshHashed, auth.Iat, auth.Exp, auth.LastUpdate)
	return err
}

func (repo *AuthRepository) GetByRefreshHashed(token []byte) (domain.Auth, error) {
	auth := domain.Auth{RefreshHashed: token}
	query := `
		SELECT user_id, issued_at, expires_at, last_time_updated
		FROM users_tokens
		WHERE refresh_token_hashed = $1`
	err := repo.DB.QueryRow(query, token).Scan(&auth.UserID, &auth.Iat, &auth.Exp, &auth.LastUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth, domain.ErrNotFound
		}
		return auth, err
	}
	return auth, nil
}

func (repo *AuthRepository) GetByUUID(uuid uuid.UUID) (domain.Auth, error) {
	auth := domain.Auth{UserID: uuid}
	query := `
		SELECT refresh_token_hashed, issued_at, expires_at, last_time_updated
		FROM users_tokens
		WHERE refresh_token_hashed = $1`
	err := repo.DB.QueryRow(query, uuid).Scan(&auth.RefreshHashed, &auth.Iat, &auth.Exp, &auth.LastUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth, domain.ErrNotFound
		}
		return auth, err
	}
	return auth, nil
}
