package game

import (
	"errors"
	"strings"
)

// Game represents a Tic-Tac-Toe game
type Game struct {
	Board  [3][3]string // 3x3 game board
	Turn   string       // Current player's turn ("X" or "O")
	Winner string       // Winner ("X", "O", or "tie")
}

// NewGame initializes a new Tic-Tac-Toe game
func NewGame() *Game {
	return &Game{
		Board:  [3][3]string{},
		Turn:   "X", // "X" always starts
		Winner: "",
	}
}

// MakeMove processes a player's move
func (g *Game) MakeMove(row, col int) error {
	// Check if move if within bounds
	if row < 0 || row > 2 || col < 0 || col > 2 {
		return errors.New("move out of boudns")
	}

	// Check if the cell is already occupied
	if g.Board[row][col] != "" {
		return errors.New("cell already occupied")
	}

	// Make the move
	g.Board[row][col] = g.Turn

	// Check if this move wins the game
	if g.checkWin() {
		g.Winner = g.Turn
		return nil
	}

	// Check if the game is tied
	if g.checkTie() {
		g.Winner = "tie"
		return nil
	}

	// Switch turns
	if g.Turn == "X" {
		g.Turn = "O"
	} else {
		g.Turn = "X"
	}

	return nil
}

// checkWin checks if the current player has won
func (g *Game) checkWin() bool {
	// Check rows, columns, and diagonals
	for i := 0; i < 3; i++ {
		// Check rows
		if g.Board[i][0] == g.Turn && g.Board[i][1] == g.Turn && g.Board[i][2] == g.Turn {
			return true
		}

		// Check Cols
		if g.Board[0][i] == g.Turn && g.Board[1][i] == g.Turn && g.Board[2][i] == g.Turn {
			return true
		}
	}

	// Check diagonal 1 (\)
	if g.Board[0][0] == g.Turn && g.Board[1][1] == g.Turn && g.Board[2][2] == g.Turn {
		return true
	}
	// Check diagonal 2 (/)
	if g.Board[0][2] == g.Turn && g.Board[1][1] == g.Turn && g.Board[2][0] == g.Turn {
		return true
	}

	return false
}

// checkTie checks if the game is tied
func (g *Game) checkTie() bool {
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == "" {
				// After checking for wins, checks every tile
				// Game is only tied when all tiles have been taken up
				return false
			}
		}
	}
	return true
}

// String returns a visual representation of the board
func (g *Game) String() string {
	var sb strings.Builder
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == "" {
				sb.WriteString("-")
			} else {
				sb.WriteString(cell)
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
