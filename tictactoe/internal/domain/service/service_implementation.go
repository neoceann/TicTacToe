package service

import (
	"context"
	"fmt"
	"tictactoe/internal/datasource/repository"
	"tictactoe/internal/domain/model"
	"time"

	"github.com/google/uuid"
)

type GameServiceImpl struct {
	repo repository.GameRepository
	algo MinimaxAlgorithm
}

func NewGameService(repo repository.GameRepository, algo MinimaxAlgorithm) GameService {
	return &GameServiceImpl{
		repo: repo,
		algo: algo,
	}
}

func (s *GameServiceImpl) GetNextMove(ctx context.Context, gameID uuid.UUID) (*model.Game, error) {
	game, err := s.repo.Get(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	row, col := s.algo.FindBestMove(game)
	
	if err := game.MakeMove(row, col, 2); err != nil {
		return nil, fmt.Errorf("move AI failed: %w", err)
	}

	if winner := game.CheckWinner(); winner != 0 {
        if winner == 1 {
            game.State = model.StatePlayerWon
        } else {
            game.State = model.StateAIWon
        }
    } else if game.IsFull() {
        game.State = model.StateDraw
    }
	
	if err := s.repo.Save(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to save game: %w", err)
	}
	
	return game, nil
}

func (s *GameServiceImpl) ValidateField(ctx context.Context, gameID uuid.UUID, newField model.GameField) (bool, error) {
    originalGame, err := s.repo.Get(ctx, gameID)
    if err != nil {
        return false, fmt.Errorf("failed to get game: %w", err)
    }
    
    return s.isValidContinuation(originalGame.Field, newField), nil
}

func (s *GameServiceImpl) isValidContinuation(oldField, newField model.GameField) bool {
    if len(oldField) != len(newField) {
        return false
    }
    for i := range oldField {
        if len(oldField[i]) != len(newField[i]) {
            return false
        }
    }
    
    size := len(oldField)
    changes := 0
    
    for i := 0; i < size; i++ {
        for j := 0; j < size; j++ {
            oldVal := oldField[i][j]
            newVal := newField[i][j]
            
            if oldVal != newVal {
                if oldVal != 0 {
                    return false
                }
                
                if newVal != 1 && newVal != 2 {
                    return false
                }
                
                changes++
            }
        }
    }
    
    if changes != 1 {
        return false
    }
        
    return true
}

func (s *GameServiceImpl) GetGameState(ctx context.Context, gameID uuid.UUID) (string, error) {
	game, err := s.repo.Get(ctx, gameID)
	if err != nil {
		return "", fmt.Errorf("failed to get game: %w", err)
	}
	
	return game.State, nil
}

func (s *GameServiceImpl) MakePlayerMove(ctx context.Context, gameID uuid.UUID, row, col, player int) (*model.Game, error) {
	game, err := s.repo.Get(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	
	if game.State != model.StateInProgress {
		return nil, fmt.Errorf("game is already finished")
	}
	
	if err := game.MakeMove(row, col, player); err != nil {
		return nil, fmt.Errorf("move failed: %w", err)
	}

	if winner := game.CheckWinner(); winner != 0 {
        if winner == 1 {
            game.State = model.StatePlayerWon
        } else {
            game.State = model.StateAIWon
        }
    } else if game.IsFull() {
        game.State = model.StateDraw
    }
	
	if err := s.repo.Save(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to save game: %w", err)
	}
	
	return game, nil
}

func (s *GameServiceImpl) CreateGame(ctx context.Context, size int) (*model.Game, error) {
	if size < 3 || size > 10 {
		return nil, fmt.Errorf("invalid size: must be between 3 and 10")
	}
	
	field := make(model.GameField, size)
	for i := range field {
		field[i] = make([]int, size)
	}
	
	game := &model.Game{
		ID:         uuid.New(),
		Field:      field,
		State:      model.StateInProgress,
		PlayerTurn: 1,
		Size:       size,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	if err := s.repo.Save(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to save game: %w", err)
	}
	
	return game, nil
}

func (s *GameServiceImpl) GetGame(ctx context.Context, gameID uuid.UUID) (*model.Game, error) {
	return s.repo.Get(ctx, gameID)
}

func (s *GameServiceImpl) FindBestMove(game *model.Game) (row, col int) {
	return s.algo.FindBestMove(game)
}