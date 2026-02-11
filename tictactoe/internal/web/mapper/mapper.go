package mapper

import (
	"encoding/json"
	"fmt"
	"net/http"
	domainModel "tictactoe/internal/domain/model"
	webModel "tictactoe/internal/web/model"
)

func ToMoveResponse(game *domainModel.Game) *webModel.MoveResponse {
	if game == nil {
		return nil
	}

	playerTurnID := game.UserID.String()

	if game.PlayerTurn == domainModel.OPlayerIcon {
		playerTurnID = game.User2ID.String()
	}

	return &webModel.MoveResponse{
		GameID: game.ID.String(),
		Field:  game.Field,
		Status: string(game.State),
		Turn: game.PlayerTurn,
		TurnID:	playerTurnID,
	}
}

func FinishedGameResponse(game *domainModel.Game) *webModel.FinishResponse {
	if game == nil {
		return nil
	}

	playerWinnerID := ""
	switch game.State {
	case domainModel.StatePlayerXWon:
		playerWinnerID = game.UserID.String()
	case domainModel.StatePlayerOWon:
		playerWinnerID = game.User2ID.String()
	}

	return &webModel.FinishResponse{
		GameID: game.ID.String(),
		Field: game.Field,
		Status: string(game.State),
		WinnerID: playerWinnerID,
	}
}

func ToErrorResponse(err error) *webModel.ErrorResponse {
	return &webModel.ErrorResponse{
		Error: err.Error(),
	}
}

func FieldFromRequest(field [][]int) domainModel.GameField {
	return domainModel.GameField(field)
}

func ValidateField(field [][]int, expectedSize int) error {
	if field == nil {
		return fmt.Errorf("field is required")
	}

	if len(field) != expectedSize {
		return fmt.Errorf("field must be %dx%d, got %dx%d", 
			expectedSize, expectedSize, len(field), len(field))
	}

	for i, row := range field {
		if len(row) != expectedSize {
			return fmt.Errorf("row %d has invalid length: expected %d, got %d", 
				i, expectedSize, len(row))
		}
		
		for j, cell := range row {
			if cell < 0 || cell > 2 {
				return fmt.Errorf("invalid cell value at [%d][%d]: %d (must be 0, 1, or 2)", 
					i, j, cell)
			}
		}
	}

	return nil
}

func ParseMoveRequest(body []byte) (*webModel.MoveRequest, error) {
	var req webModel.MoveRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}
	return &req, nil
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}