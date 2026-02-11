package model

import "time"

type GameModel struct {
	ID        string `db:"id"`
	UserID	  string `db:"user_id"`
	User2ID	  string `db:"user2_id"`
	Field     string `db:"field"`
	State     string `db:"state"`
	PlayerTurn int	`db:"player_turn"`
	Size      int	`db:"size"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Opponent string `db:"opponent"`
}

type WaitingGamesModel struct {
	ID string `db:"id"`
	UserID string `db:"user_id"`
}