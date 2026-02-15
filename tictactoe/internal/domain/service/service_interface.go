package service

import (
	"github.com/google/uuid"
	"context"
	"tictactoe/internal/domain/model"
)

type GameService interface {
	GetNextMove(ctx context.Context, gameID uuid.UUID, userID uuid.UUID) (*model.Game, error)
	ValidateField(ctx context.Context, gameID uuid.UUID, userID uuid.UUID, field model.GameField) (bool, error)
	GetGameState(ctx context.Context, gameID uuid.UUID, userID uuid.UUID) (string, error)
	MakePlayerMove(ctx context.Context, gameID uuid.UUID, userID uuid.UUID, row, col, player int) (*model.Game, error)
    CreateGame(ctx context.Context, userID uuid.UUID, size int, opponent string) (*model.Game, error) 
    GetGame(ctx context.Context, gameID uuid.UUID, userID uuid.UUID) (*model.Game, error)

	GetUserIdByGameId(ctx context.Context, gameID uuid.UUID) (uuid.UUID, error)
	UpdateAfterJoin(ctx context.Context, game *model.Game) error
	GetWaitingGames(ctx context.Context) ([]*model.WaitingGames, error)
	GetPublicUserInfoById(ctx context.Context, userID uuid.UUID) (*model.PublicUserInfo, error)

	GetFinishedGamesByID(ctx context.Context, userID uuid.UUID) ([]*model.FinishedGamesInfo, error)
	GetLeaderboard(ctx context.Context, limit int) ([]*model.PlayerWinrateInfo, error)
}

type MinimaxAlgorithm interface {
    FindBestMove(game *model.Game) (row, col int)
}