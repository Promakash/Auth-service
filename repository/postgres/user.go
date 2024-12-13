package postgres

import (
	"auth_service/domain"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Put(user domain.User) error {
	query := `UPDATE users 
			  SET IP = $2, email = $3 
			  WHERE user_id = $1`
	_, err := r.DB.Exec(query, user.UserID, user.IP, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByUUID(uuid uuid.UUID) (domain.User, error) {
	user := domain.User{UserID: uuid}
	query := "SELECT IP, email FROM users WHERE user_id = $1"
	row := r.DB.QueryRow(query, uuid)

	if err := row.Scan(&user.IP, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return user, err
	}

	return user, nil
}
