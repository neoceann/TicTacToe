package module

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/google/uuid"
	"tictactoe/internal/domain/service"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/web/mapper"
	webModel "tictactoe/internal/web/model"

)

type GameHandler struct {
	gameService service.GameService
}

func NewGameHandler(gameService service.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		mapper.WriteJSON(w, http.StatusMethodNotAllowed, 
			mapper.ToErrorResponse(fmt.Errorf("method not allowed")))
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("invalid URL format")))
		return
	}
	
	gameIDStr := pathParts[2]
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("invalid game UUID: %v", err)))
		return
	}

	var req webModel.MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("invalid JSON: %v", err)))
		return
}

	currentGame, err := h.gameService.GetGame(r.Context(), gameID)
	if err != nil {
		mapper.WriteJSON(w, http.StatusNotFound,
			mapper.ToErrorResponse(fmt.Errorf("game not found: %v", err)))
		return
	}

	if err := mapper.ValidateField(req.Field, currentGame.Size); err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(err))
		return
	}

	isValid, err := h.gameService.ValidateField(
		r.Context(),
		gameID,
		mapper.FieldFromRequest(req.Field),
	)
	if err != nil {
		mapper.WriteJSON(w, http.StatusInternalServerError,
			mapper.ToErrorResponse(err))
		return
	}
	
	if !isValid {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("invalid game field: previous moves have been changed")))
		return
	}

	userRow, userCol := h.findUserMove(currentGame.Field, req.Field)
	if userRow == -1 {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("invalid move: no valid user move found or multiple moves detected")))
		return
	}

	gameAfterUserMove, err := h.gameService.MakePlayerMove(r.Context(), gameID, userRow, userCol, 1)
	if err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(err))
		return
	}

	status, err := h.gameService.GetGameState(r.Context(), gameID)
	if err != nil {
		mapper.WriteJSON(w, http.StatusInternalServerError,
			mapper.ToErrorResponse(err))
		return
	}
	
	if status != model.StateInProgress {
		mapper.WriteJSON(w, http.StatusOK, mapper.ToMoveResponse(gameAfterUserMove))
		return
	}

	gameAfterAIMove, err := h.gameService.GetNextMove(r.Context(), gameID)
	if err != nil {
		mapper.WriteJSON(w, http.StatusInternalServerError,
			mapper.ToErrorResponse(err))
		return
	}

	mapper.WriteJSON(w, http.StatusOK, mapper.ToMoveResponse(gameAfterAIMove))
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		mapper.WriteJSON(w, http.StatusMethodNotAllowed,
			mapper.ToErrorResponse(fmt.Errorf("method not allowed")))
		return
	}

	game, err := h.gameService.CreateGame(r.Context(), model.FieldSize)
	if err != nil {
		mapper.WriteJSON(w, http.StatusInternalServerError,
			mapper.ToErrorResponse(err))
		return
	}

	mapper.WriteJSON(w, http.StatusCreated, &webModel.CreateGameResponse{
		GameID: game.ID.String(),
		Field:  game.Field,
		Status: string(game.State),
	})
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		mapper.WriteJSON(w, http.StatusMethodNotAllowed,
			mapper.ToErrorResponse(fmt.Errorf("method not allowed")))
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("invalid URL format")))
		return
	}
	
	gameIDStr := pathParts[2]
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("invalid game UUID: %v", err)))
		return
	}

	game, err := h.gameService.GetGame(r.Context(), gameID)
	if err != nil {
		mapper.WriteJSON(w, http.StatusNotFound,
			mapper.ToErrorResponse(err))
		return
	}

	mapper.WriteJSON(w, http.StatusOK, mapper.ToMoveResponse(game))
}

func (h *GameHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		mapper.WriteJSON(w, http.StatusMethodNotAllowed,
			mapper.ToErrorResponse(fmt.Errorf("method not allowed")))
		return
	}

	mapper.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "tictactoe",
	})
}

func (h *GameHandler) findUserMove(originalField, newField [][]int) (int, int) {
	if len(originalField) != len(newField) {
		return -1, -1
	}

	row, col := -1, -1
	changes := 0

	for i := 0; i < len(originalField); i++ {
		if len(originalField[i]) != len(newField[i]) {
			return -1, -1
		}
		
		for j := 0; j < len(originalField[i]); j++ {
			if originalField[i][j] != newField[i][j] {
				if originalField[i][j] != 0 {
					return -1, -1
				}
				if newField[i][j] != 1 {
					return -1, -1
				}
				row, col = i, j
				changes++
			}
		}
	}

	if changes == 1 {
		return row, col
	}
	return -1, -1
}