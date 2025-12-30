package model

import (
	"time"
	"fmt"
	"github.com/google/uuid"
)

const FieldSize int = 3

type GameField [][]int

const (
	StateInProgress = "Game in progress"
	StatePlayerWon = "Player won"
	StateAIWon = "AI won"
	StateDraw = "Draw"
)

type Game struct {
	ID        uuid.UUID
	Field     GameField
	State     string
	PlayerTurn int
	Size      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewField(size int) GameField {
    field := make(GameField, size)
    for i := range field {
        field[i] = make([]int, size)
    }
    return field
}

func (g *Game) DeepCopy() *Game {
    return &Game{
        ID:         g.ID,
        Field:      g.Field.DeepCopy(),
        State:      g.State,
        PlayerTurn: g.PlayerTurn,
        Size:       g.Size,
        CreatedAt:  g.CreatedAt,
        UpdatedAt:  g.UpdatedAt,
    }
}

func (f GameField) DeepCopy() GameField {
    size := len(f)
    copy := make(GameField, size)
    for i := range f {
        copy[i] = make([]int, size)
        for j := range f[i] {
            copy[i][j] = f[i][j]
        }
    }
    return copy
}

func (g *Game) MakeMove(row, col, player int) error {
    if row < 0 || row >= g.Size || col < 0 || col >= g.Size {
        return fmt.Errorf("invalid coordinates")
    }
    if !g.Field.IsEmpty(row, col) {
        return fmt.Errorf("cell already occupied")
    }
    
    g.Field[row][col] = player
    g.PlayerTurn = 3 - player
    g.UpdatedAt = time.Now()
    return nil
}

func (f GameField) IsEmpty(row, col int) bool {
    return f[row][col] == 0
}

func (g *Game) CheckWinner() int {
    size := g.Size
    
    for i := 0; i < size; i++ {
        first := g.Field[i][0]
        if first == 0 {
            continue
        }
        win := true
        for j := 1; j < size; j++ {
            if g.Field[i][j] != first {
                win = false
                break
            }
        }
        if win {
            return first
        }
    }
    
    for j := 0; j < size; j++ {
        first := g.Field[0][j]
        if first == 0 {
            continue
        }
        win := true
        for i := 1; i < size; i++ {
            if g.Field[i][j] != first {
                win = false
                break
            }
        }
        if win {
            return first
        }
    }
    
    first := g.Field[0][0]
    if first != 0 {
        win := true
        for i := 1; i < size; i++ {
            if g.Field[i][i] != first {
                win = false
                break
            }
        }
        if win {
            return first
        }
    }
    
    first = g.Field[0][size-1]
    if first != 0 {
        win := true
        for i := 1; i < size; i++ {
            if g.Field[i][size-1-i] != first {
                win = false
                break
            }
        }
        if win {
            return first
        }
    }
    
    return 0
}

func (g *Game) IsFull() bool {
    for i := 0; i < g.Size; i++ {
        for j := 0; j < g.Size; j++ {
            if g.Field[i][j] == 0 {
                return false
            }
        }
    }
    return true
}