package minimax

import (
    "tictactoe/internal/domain/model"
)

const (
    MaxScore = 1000
    MinScore = -1000
)

type Minimax struct {
    computerPlayer int
    humanPlayer    int
}

func NewMinimax(computerPlayer int) *Minimax {
    humanPlayer := 1
    if computerPlayer == 1 {
        humanPlayer = 2
    }
    
    return &Minimax{
        computerPlayer: computerPlayer,
        humanPlayer:    humanPlayer,
    }
}

func (m *Minimax) FindBestMove(game *model.Game) (int, int) {
    gameCopy := game.DeepCopy()
    
    bestScore := MinScore
    bestRow, bestCol := -1, -1
    alpha := MinScore
    beta := MaxScore
    
    priorityMoves := m.getPriorityMoves(gameCopy)
    
    for _, move := range priorityMoves {
        row, col := move[0], move[1]
        if gameCopy.Field.IsEmpty(row, col) {
            currentGame := gameCopy.DeepCopy()
            currentGame.MakeMove(row, col, m.computerPlayer)
            
            score := m.minimax(currentGame, 0, false, alpha, beta)
            
            if score > bestScore {
                bestScore = score
                bestRow, bestCol = row, col
            }
            
            alpha = max(alpha, bestScore)
            if beta <= alpha {
                break
            }
        }
    }
    
    if bestRow == -1 {
        return m.findFirstEmpty(gameCopy)
    }
    
    return bestRow, bestCol
}

func (m *Minimax) minimax(game *model.Game, depth int, isMaximizing bool, alpha, beta int) int {
    winner := game.CheckWinner()
    if winner == m.computerPlayer {
        return MaxScore - depth
    }
    if winner == m.humanPlayer {
        return depth - MaxScore
    }
    if game.IsFull() {
        return 0
    }
    
    if isMaximizing {
        maxScore := MinScore
        for i := 0; i < game.Size; i++ {
            for j := 0; j < game.Size; j++ {
                if game.Field.IsEmpty(i, j) {
                    gameCopy := game.DeepCopy()
                    gameCopy.MakeMove(i, j, m.computerPlayer)
                    
                    score := m.minimax(gameCopy, depth+1, false, alpha, beta)
                    maxScore = max(maxScore, score)
                    alpha = max(alpha, score)
                    
                    if beta <= alpha {
                        break
                    }
                }
            }
        }
        return maxScore
    } else {
        minScore := MaxScore
        for i := 0; i < game.Size; i++ {
            for j := 0; j < game.Size; j++ {
                if game.Field.IsEmpty(i, j) {
                    gameCopy := game.DeepCopy()
                    gameCopy.MakeMove(i, j, m.humanPlayer)
                    
                    score := m.minimax(gameCopy, depth+1, true, alpha, beta)
                    minScore = min(minScore, score)
                    beta = min(beta, score)
                    
                    if beta <= alpha {
                        break
                    }
                }
            }
        }
        return minScore
    }
}

func (m *Minimax) getPriorityMoves(game *model.Game) [][2]int {
    size := game.Size
    var moves [][2]int
    
    if size%2 == 1 {
        center := size / 2
        moves = append(moves, [2]int{center, center})
    }
    
    corners := [][2]int{
        {0, 0}, {0, size - 1},
        {size - 1, 0}, {size - 1, size - 1},
    }
    moves = append(moves, corners...)
    
    for i := 0; i < size; i++ {
        for j := 0; j < size; j++ {
            isCenter := size%2 == 1 && i == size/2 && j == size/2
            isCorner := (i == 0 || i == size-1) && (j == 0 || j == size-1)
            
            if !isCenter && !isCorner {
                moves = append(moves, [2]int{i, j})
            }
        }
    }
    
    return moves
}

func (m *Minimax) findFirstEmpty(game *model.Game) (int, int) {
    for i := 0; i < game.Size; i++ {
        for j := 0; j < game.Size; j++ {
            if game.Field.IsEmpty(i, j) {
                return i, j
            }
        }
    }
    return -1, -1
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}