package repository

import (
	"sync"

	"fmt"

	"tictactoe/internal/datasource/model"
)

type GameStorage struct {
	storage sync.Map
}

func NewGameStorage() *GameStorage {
	return &GameStorage{storage: sync.Map{}}
}

func (storage *GameStorage) Save(game *model.GameModel) {
	storage.storage.Store(game.ID, game)
}

func (storage *GameStorage) Get(gameID string) (*model.GameModel, error) {
	gameAsInterface, exists := storage.storage.Load(gameID)

	if !exists {
		return  nil, fmt.Errorf("cant find game by this ID")
	}

	game, ok := gameAsInterface.(*model.GameModel)

	if !ok {
		return nil, fmt.Errorf("gameModel type assertion failed")
	}

	return game, nil
}