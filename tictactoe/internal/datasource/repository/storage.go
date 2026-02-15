package repository

import (
	"context"
	"fmt"
	"time"

	"tictactoe/internal/datasource/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GameStorage struct {
	DBConnection *pgx.Conn
}

func NewGameStorage() (*GameStorage, error) {
	connectionData := "postgresql://postgres:root@localhost:5432/postgres_tictactoe"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, connectionData)

	if err != nil {
		return nil, fmt.Errorf("failed connection to db: %w", err)
	}

	storage := &GameStorage{DBConnection: conn}

	if err := storage.createTable(ctx); err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to init tables: %w", err)
	}

	return storage, nil
}

func (storage *GameStorage) createTable(ctx context.Context) error {

	usersQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			login VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	
	_, err := storage.DBConnection.Exec(ctx, usersQuery)

	gameQuery := `
		CREATE TABLE IF NOT EXISTS games (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			user2_id UUID NOT NULL,
			field TEXT NOT NULL,
			state VARCHAR(20) NOT NULL,
			player_turn INTEGER NOT NULL,
			size INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			opponent VARCHAR(20) NOT NULL
		)
	`
	_, err = storage.DBConnection.Exec(ctx, gameQuery)

	return err
}

func (storage *GameStorage) Save(game *model.GameModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO games (id, user_id, user2_id, field, state, player_turn, size, created_at, updated_at, opponent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			field = EXCLUDED.field,
			state = EXCLUDED.state,
			player_turn = EXCLUDED.player_turn,
			updated_at = EXCLUDED.updated_at
	`
	now := time.Now()

	if game.CreatedAt.IsZero() {
		game.CreatedAt = now
	}
	game.UpdatedAt = now

	_, err := storage.DBConnection.Exec(ctx, query,
		game.ID,
		game.UserID,
		game.User2ID,
		game.Field,
		game.State,
		game.PlayerTurn,
		game.Size,
		game.CreatedAt,
		game.UpdatedAt,
		game.Opponent,
	)

	if err != nil {
		return err
	}

	return nil
}

func (storage *GameStorage) Get(gameID string, userID string) (*model.GameModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, user2_id, field, state, player_turn, size, created_at, updated_at, opponent
		FROM games
		WHERE id = $1 AND (user_id = $2 OR user2_id = $2)
	`

	game := &model.GameModel{}
	err := storage.DBConnection.QueryRow(ctx, query, gameID, userID).Scan(
		&game.ID,
		&game.UserID,
		&game.User2ID,
		&game.Field,
		&game.State,
		&game.PlayerTurn,
		&game.Size,
		&game.CreatedAt,
		&game.UpdatedAt,
		&game.Opponent,
	)

	if err != nil {
		if err == pgx.ErrNoRows || userID != game.UserID {
			return nil, fmt.Errorf("game with ID %s not found or access denied", gameID)
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return game, nil
}

func (storage *GameStorage) UpdateAfterJoinToGame(game *model.GameModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE games
		SET user2_id = $1,
		state = $2,
		updated_at = $3
		WHERE id = $4
	`
	now := time.Now()

	game.UpdatedAt = now

	_, err := storage.DBConnection.Exec(ctx, query,
		game.User2ID,
		game.State,
		game.UpdatedAt,
		game.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (storage *GameStorage) GetUserIdByGameId(gameID string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT user_id
		FROM games
		WHERE id = $1
	`

	userID := ""
	err := storage.DBConnection.QueryRow(ctx, query, gameID).Scan(
		&userID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, fmt.Errorf("game with ID %s not found or access denied", gameID)
		}
		return uuid.Nil, fmt.Errorf("failed to get game: %w", err)
	}

	userUUID, _ := uuid.Parse(userID)

	return userUUID, nil
}

func (storage *GameStorage) GetWaitingGames() ([]*model.WaitingGamesModel, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id
		FROM games
		WHERE state = 'Waiting for player O'
	`

	rows, err := storage.DBConnection.Query(ctx, query)

    if err != nil {
        return nil, fmt.Errorf("get waiting games: %w", err)
    }
    defer rows.Close()
    
    var games []*model.WaitingGamesModel
    for rows.Next() {
        game := &model.WaitingGamesModel{}
        err := rows.Scan(
            &game.ID,
            &game.UserID,
        )
        if err != nil {
            return nil, fmt.Errorf("scan game: %w", err)
        }
        games = append(games, game)
    }
    
    return games, nil
}

func (storage *GameStorage) SaveUser(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (id, login, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := storage.DBConnection.Exec(ctx, query,
		user.ID,
		user.Login,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (storage *GameStorage) FindUserByLogin(login string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, login, password_hash, created_at, updated_at
		FROM users
		WHERE login = $1
	`

	user := &model.User{}
	err := storage.DBConnection.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (storage *GameStorage) GetPublicUserInfoById(id string) (*model.PublicUserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, login
		FROM users
		WHERE id = $1
	`

	user := &model.PublicUserInfo{}
	err := storage.DBConnection.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Login,

	)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (storage *GameStorage) GetFinishedGames(id string) ([]*model.FinishedGamesInfo, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
			SELECT id, state, opponent, created_at
			FROM games
			WHERE 
				(state IN ('Player X won', 'Draw') AND user_id = $1) OR
				(state IN ('Player O won', 'Draw') AND user2_id = $1)
	`

	rows, err := storage.DBConnection.Query(ctx, query, id)

    if err != nil {
        return nil, fmt.Errorf("get finished games: %w", err)
    }
    defer rows.Close()
    
    var games []*model.FinishedGamesInfo
    for rows.Next() {
        game := &model.FinishedGamesInfo{}
        err := rows.Scan(
            &game.GameID,
            &game.State,
			&game.Opponent,
			&game.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("scan game: %w", err)
        }
        games = append(games, game)
    }
    
    return games, nil
}

func (storage *GameStorage) GetLeaderboard(limit int) ([]*model.PlayerWinrateInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			u.id AS user_id,
			ROUND (COUNT(CASE 
				WHEN (g.state = 'Player X won' AND g.user_id = u.id) OR
					(g.state = 'Player O won' AND g.user2_id = u.id) 
				THEN 1 END) * 1.0 / NULLIF(COUNT(*), 0), 1) AS win_ratio
		FROM users u
		LEFT JOIN games g ON (u.id = g.user_id OR u.id = g.user2_id)
		WHERE g.state IS NOT NULL
		GROUP BY u.id
		HAVING COUNT(*) > 0
		ORDER BY win_ratio DESC
		LIMIT $1;
	`

	rows, err := storage.DBConnection.Query(ctx, query, limit)

    if err != nil {
        return nil, fmt.Errorf("get leaderboard: %w", err)
    }
    defer rows.Close()
    
    var players []*model.PlayerWinrateInfo
    for rows.Next() {
        player := &model.PlayerWinrateInfo{}
        err := rows.Scan(
            &player.ID,
            &player.Rating,
        )
        if err != nil {
            return nil, fmt.Errorf("scan game: %w", err)
        }
        players = append(players, player)
    }
    
    return players, nil
}