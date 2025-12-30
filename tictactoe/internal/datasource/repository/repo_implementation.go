package repository

import (
	"context"
	"fmt"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/datasource/mapper"
	"github.com/google/uuid"
)

type GameRepositoryImpl struct {
	storage *GameStorage
}

func NewGameRepo(storage *GameStorage) GameRepository {
	return &GameRepositoryImpl{storage: storage}
}

func (r *GameRepositoryImpl) Save(ctx context.Context, game *model.Game) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	
	gameModel, err := mapper.FromDomainToDs(game)
	if err != nil {
		return fmt.Errorf("failed to convert game to model: %w", err)
	}
	
	r.storage.Save(gameModel)
	
	return nil
}

func (r *GameRepositoryImpl) Get(ctx context.Context, gameID uuid.UUID) (*model.Game, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	
	gameModel, err := r.storage.Get(gameID.String())
	if err != nil {
		return nil, err
	}
	
	game, err := mapper.FromDsToDomain(gameModel)
	if err != nil {
		return nil, fmt.Errorf("failed to convert model to game: %w", err)
	}
	
	return game, nil
}