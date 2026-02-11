package repository

import (
	"context"
	"tictactoe/internal/domain/model"
	"github.com/google/uuid"
)

type GameRepository interface {
	Save(ctx context.Context, game *model.Game) error
	
	Get(ctx context.Context, gameID uuid.UUID, userID uuid.UUID) (*model.Game, error)

	GetUserIdByGameId(ctx context.Context, gameID uuid.UUID) (uuid.UUID, error)

	UpdateAfterJoin(ctx context.Context, game *model.Game) error

	GetWaitingGames(ctx context.Context) ([]*model.WaitingGames, error)

	GetPublicUserInfoById(ctx context.Context, userID uuid.UUID) (*model.PublicUserInfo, error)
}