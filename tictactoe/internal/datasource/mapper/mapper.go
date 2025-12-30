package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	dsModel "tictactoe/internal/datasource/model"
	domainModel "tictactoe/internal/domain/model"
)

func FromDsToDomain(model *dsModel.GameModel) (*domainModel.Game, error) {
	if model == nil {
		return nil, fmt.Errorf("model is nil")
	}

	id, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, fmt.Errorf("bad UUID")
	}

	var field domainModel.GameField

	if err := json.Unmarshal([]byte(model.Field), &field); err != nil {
		return nil, fmt.Errorf("cant parse json field")
	}

	return &domainModel.Game{
		ID: id,
		Field: field,
		State: model.State,
		PlayerTurn: model.PlayerTurn,
		Size: model.Size,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func FromDomainToDs(game *domainModel.Game) (*dsModel.GameModel, error) {
	if game == nil {
		return nil, fmt.Errorf("game is nil")
	}
	
	fieldJSON, err := json.Marshal(game.Field)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field: %w", err)
	}
	
	return &dsModel.GameModel{
		ID:        game.ID.String(), 
		Field:     string(fieldJSON),
		State:     game.State,
		PlayerTurn: game.PlayerTurn,
		Size:      game.Size,
		CreatedAt: game.CreatedAt,
		UpdatedAt: game.UpdatedAt,
	}, nil
}