package model

import "time"

type GameModel struct {
	ID        string
	Field     string
	State     string
	PlayerTurn int
	Size      int
	CreatedAt time.Time
	UpdatedAt time.Time
}