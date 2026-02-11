package module

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/domain/service"
	"tictactoe/internal/web/mapper"
	webModel "tictactoe/internal/web/model"

	"github.com/google/uuid"
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

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		mapper.WriteJSON(w, http.StatusUnauthorized,
			mapper.ToErrorResponse(fmt.Errorf("authentication required")))
		return
	}

	currentGame, err := h.gameService.GetGame(r.Context(), gameID, userID)
	if err != nil {
		mapper.WriteJSON(w, http.StatusNotFound,
			mapper.ToErrorResponse(fmt.Errorf("game not found: %v", err)))
		return
	}

	if currentGame.PlayerTurn == model.XPlayerIcon && userID != currentGame.UserID ||
		currentGame.PlayerTurn == model.OPlayerIcon && userID != currentGame.User2ID {
			mapper.WriteJSON(w, http.StatusBadRequest,
				mapper.ToErrorResponse(fmt.Errorf("not your turn! Current Player: %d", currentGame.PlayerTurn)))
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
		userID,
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

	userRow, userCol := h.findUserMove(currentGame.Field, req.Field, currentGame.PlayerTurn)
	if userRow == -1 {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("invalid move: no valid user move found or multiple moves detected")))
		return
	}

	gameAfterUserMove, err := h.gameService.MakePlayerMove(r.Context(), gameID, userID, userRow, userCol, currentGame.PlayerTurn)
	if err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(err))
		return
	}

	status, err := h.gameService.GetGameState(r.Context(), gameID, userID)
	if err != nil {
		mapper.WriteJSON(w, http.StatusInternalServerError,
			mapper.ToErrorResponse(err))
		return
	}
	
	if status != model.StateInProgress {
		mapper.WriteJSON(w, http.StatusOK, mapper.FinishedGameResponse(gameAfterUserMove))
		return
	}
	
	if status == model.StateInProgress && currentGame.Opponent == model.HumanOpponent {
		mapper.WriteJSON(w, http.StatusOK, mapper.ToMoveResponse(gameAfterUserMove))
		return
	}

	if currentGame.Opponent != model.HumanOpponent {
		gameAfterAIMove, err := h.gameService.GetNextMove(r.Context(), gameID, userID)
		if err != nil {
			mapper.WriteJSON(w, http.StatusInternalServerError,
				mapper.ToErrorResponse(err))
			return
		}

		status, err := h.gameService.GetGameState(r.Context(), gameID, userID)
		if err != nil {
			mapper.WriteJSON(w, http.StatusInternalServerError,
				mapper.ToErrorResponse(err))
			return
		}
		
		if status != model.StateInProgress {
			mapper.WriteJSON(w, http.StatusOK, mapper.FinishedGameResponse(gameAfterAIMove))
			return
		} else {
			mapper.WriteJSON(w, http.StatusOK, mapper.ToMoveResponse(gameAfterAIMove))
		}
	}
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		mapper.WriteJSON(w, http.StatusMethodNotAllowed,
			mapper.ToErrorResponse(fmt.Errorf("method not allowed")))
		return
	}
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		mapper.WriteJSON(w, http.StatusUnauthorized,
			mapper.ToErrorResponse(fmt.Errorf("authentication required")))
		return
	}

	var opponent webModel.OpponentInfo
    json.NewDecoder(r.Body).Decode(&opponent)

	if opponent.Opponent != model.HumanOpponent && opponent.Opponent != model.AIOpponent {
		log.Printf("opp: %s", opponent.Opponent)
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("bad opponent option")))
		return
	}

	game, err := h.gameService.CreateGame(r.Context(), userID, model.FieldSize, opponent.Opponent)
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

		
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		mapper.WriteJSON(w, http.StatusUnauthorized,
			mapper.ToErrorResponse(fmt.Errorf("authentication required")))
		return
	}

	game, err := h.gameService.GetGame(r.Context(), gameID, userID)
	if err != nil {
		mapper.WriteJSON(w, http.StatusNotFound,
			mapper.ToErrorResponse(err))
		return
	}

	if game.State == model.StateInProgress {
		mapper.WriteJSON(w, http.StatusOK, mapper.ToMoveResponse(game))
	} else {
		mapper.WriteJSON(w, http.StatusOK, mapper.FinishedGameResponse(game))
	}
}

func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request){
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
	
	gameIDStr := pathParts[3]
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("invalid game UUID: %v", err)))
		return
	}

	ownerUserId, err := h.gameService.GetUserIdByGameId(r.Context(), gameID)

	game, err := h.gameService.GetGame(r.Context(), gameID, ownerUserId)

    if game.State != model.StateWaitingForPlayer {
        mapper.WriteJSON(w, http.StatusBadRequest,
            mapper.ToErrorResponse(fmt.Errorf("game is not waiting for players")))
        return
    }
    
    if game.Opponent != model.HumanOpponent {
        mapper.WriteJSON(w, http.StatusBadRequest,
            mapper.ToErrorResponse(fmt.Errorf("this is a computer game")))
        return
    }
    
	currentUserId := r.Context().Value("user_id").(uuid.UUID)
    if game.UserID == currentUserId {
        mapper.WriteJSON(w, http.StatusBadRequest,
            mapper.ToErrorResponse(fmt.Errorf("cannot join your own game")))
        return
    }

	game.State = model.StateInProgress
	game.User2ID = currentUserId

	h.gameService.UpdateAfterJoin(r.Context(), game)

	mapper.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "player 2 joined to game",
		"gameID":	game.ID.String(),
		"ownerID":	game.UserID.String(),
		"joinedID":	game.User2ID.String(),
	})
}

func (h *GameHandler) GetWaitingGames(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
		mapper.WriteJSON(w, http.StatusMethodNotAllowed,
			mapper.ToErrorResponse(fmt.Errorf("method not allowed")))
		return
	}

	games, err := h.gameService.GetWaitingGames(r.Context())
	
	if err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest, mapper.ToErrorResponse(fmt.Errorf("cant get waitings games")))
		return
	}

    response := make([]map[string]interface{}, 0, len(games))
    for _, game := range games {
        
        response = append(response, map[string]interface{}{
            "game_id":   game.ID.String(),
			"user_id":	 game.UserID.String(),
        })
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    encoder := json.NewEncoder(w)
    encoder.SetEscapeHTML(false) 
    encoder.SetIndent("", "  ")
	encoder.Encode(response)
}


func (h *GameHandler) GetPublicUserInfoById(w http.ResponseWriter, r *http.Request) {
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
	
	userIDStr := pathParts[2]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("invalid user UUID: %v", err)))
		return
	}

	info, err := h.gameService.GetPublicUserInfoById(r.Context(), userID)

	if err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest,
			mapper.ToErrorResponse(fmt.Errorf("error info by user id: %w", err)))
		return
	}

	mapper.WriteJSON(w, http.StatusOK, map[string]string{
		"UserID:":info.ID,
		"UserLogin": info.Login,
	})

}

func (h *GameHandler) findUserMove(originalField, newField [][]int, player_turn int) (int, int) {
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
				if newField[i][j] != player_turn {
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