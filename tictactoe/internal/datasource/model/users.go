package model

import (
	"time"
	
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type PublicUserInfo struct {
	ID           string `db:"id"`
	Login        string    `db:"login"`
}