package service

import (
	"github.com/google/uuid"
	"context"
	"tictactoe/internal/domain/model"
)

type GameService interface {
	GetNextMove(ctx context.Context, gameID uuid.UUID) (*model.Game, error)
	ValidateField(ctx context.Context, gameID uuid.UUID, field model.GameField) (bool, error)
	GetGameState(ctx context.Context, gameID uuid.UUID) (string, error)
	MakePlayerMove(ctx context.Context, gameID uuid.UUID, row, col, player int) (*model.Game, error)
    CreateGame(ctx context.Context, size int) (*model.Game, error) 
    GetGame(ctx context.Context, gameID uuid.UUID) (*model.Game, error)
}

type MinimaxAlgorithm interface {
    FindBestMove(game *model.Game) (row, col int)
}