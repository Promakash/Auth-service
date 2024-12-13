package domain

import "github.com/google/uuid"

type User struct {
	UserID uuid.UUID `json:"uuid"`
	IP     string    `json:"ip"`
	Email  string    `json:"email"`
}
