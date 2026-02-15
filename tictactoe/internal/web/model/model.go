package model

type OpponentInfo struct {
	Opponent string `json:"opponent"`
}

type MoveRequest struct {
	Field [][]int `json:"field"`
}

type MoveResponse struct {
	GameID string     `json:"game_id"`
	Field  [][]int    `json:"field"`
	Status string     `json:"status"`
	Turn int          `json:"player_turn"`
	TurnID string	  `json:"player_turn_id"`
}

type FinishResponse struct {
	GameID string     `json:"game_id"`
	Field  [][]int    `json:"field"`
	Status string     `json:"status"`
	WinnerID string   `json:"winner_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateGameResponse struct {
	GameID string     `json:"game_id"`
	Field  [][]int    `json:"field"`
	Status string     `json:"status"`
}