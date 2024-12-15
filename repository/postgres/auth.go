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
		INSERT INTO users_tokens (user_id, refresh_token_hashed, issued_at, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET refresh_token_hashed = $2, 
			issued_at = $3, 
			expires_at = $4`
	_, err := repo.DB.Exec(query, auth.UserID, auth.RefreshHashed, auth.Iat, auth.Exp)
	return err
}

func (repo *AuthRepository) GetByUUID(uuid uuid.UUID) (domain.Auth, error) {
	auth := domain.Auth{UserID: uuid}
	query := `
		SELECT refresh_token_hashed, issued_at, expires_at
		FROM users_tokens
		WHERE user_id = $1`
	err := repo.DB.QueryRow(query, uuid).Scan(&auth.RefreshHashed, &auth.Iat, &auth.Exp)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth, domain.ErrNotFound
		}
		return auth, err
	}
	return auth, nil
}
