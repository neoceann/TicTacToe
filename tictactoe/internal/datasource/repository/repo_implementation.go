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
	
	if err := r.storage.Save(gameModel); err != nil {
		return fmt.Errorf("%w", err)
	}
	
	return nil
}

func (r *GameRepositoryImpl) Get(ctx context.Context, gameID uuid.UUID, userID uuid.UUID) (*model.Game, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	
	gameModel, err := r.storage.Get(gameID.String(), userID.String())
	if err != nil {
		return nil, err
	}
	
	game, err := mapper.FromDsToDomain(gameModel)
	if err != nil {
		return nil, fmt.Errorf("failed to convert model to game: %w", err)
	}
	
	return game, nil
}

func (r *GameRepositoryImpl) GetWaitingGames(ctx context.Context) ([]*model.WaitingGames, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var games []*model.WaitingGames

	gameModels, err := r.storage.GetWaitingGames()

	if err != nil {
		return nil, fmt.Errorf("failed to get waitings games: %w", err)
	}

	for _, m := range gameModels {
		gID, err := uuid.Parse(m.ID)
		uID, err := uuid.Parse(m.UserID)
		if err != nil {
			return nil ,fmt.Errorf("cant prase uuid %w", err)
		}
		games = append(games, &model.WaitingGames{ID: gID, UserID: uID})
	}

	return games, nil
	
}

func (r *GameRepositoryImpl) GetUserIdByGameId(ctx context.Context, gameID uuid.UUID) (uuid.UUID, error) {
	if err := ctx.Err(); err != nil {
		return uuid.Nil, err
	}
	
	userID, err := r.storage.GetUserIdByGameId(gameID.String())
	if err != nil {
		return uuid.Nil, err
	}
	
	return userID, nil
}

func (r *GameRepositoryImpl) UpdateAfterJoin(ctx context.Context, game *model.Game) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	
	gameModel, err := mapper.FromDomainToDs(game)
	if err != nil {
		return fmt.Errorf("failed to convert game to model: %w", err)
	}

	if r.storage.UpdateAfterJoinToGame(gameModel) != nil {
		return err
	}
	
	return nil
}

func (r *GameRepositoryImpl) GetPublicUserInfoById(ctx context.Context, userID uuid.UUID) (*model.PublicUserInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	info, err := r.storage.GetPublicUserInfoById(userID.String())

	
	if err != nil {
		return nil, fmt.Errorf("failed to get used by ID: %w", err)
	}

	return &model.PublicUserInfo{ID: info.ID, Login: info.Login}, nil	

}
