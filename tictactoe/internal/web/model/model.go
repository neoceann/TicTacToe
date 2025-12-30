package model

type MoveRequest struct {
	Field [][]int `json:"field"`
}

type MoveResponse struct {
	GameID string     `json:"game_id"`
	Field  [][]int    `json:"field"`
	Status string     `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateGameResponse struct {
	GameID string     `json:"game_id"`
	Field  [][]int    `json:"field"`
	Status string     `json:"status"`
}