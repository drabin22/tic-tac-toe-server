package game

import "testing"

func TestGame(t *testing.T) {
	g := NewGame()

	// Test initial state
	if g.Turn != "X" {
		t.Errorf("Expected turn to be 'X', got %s", g.Turn)
	}

	// Test making valid moves
	if err := g.MakeMove(0, 0); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if g.Board[0][0] != "X" {
		t.Errorf("Expected 'X' at (0, 0), got %s", g.Board[0][0])
	}

	// Print the board after the first move
	t.Log("\nBoard after first move:\n" + g.String())

	// Test switching turns
	if g.Turn != "O" {
		t.Errorf("Expected turn to be 'O', got %s", g.Turn)
	}

	// Test invalid move (cell already occupied)
	if err := g.MakeMove(0, 0); err == nil {
		t.Error("Expected error for occupied cell, got nil")
	}

	// Print the board after invalid move attempt
	t.Log("\nBoard after invalid move attempt:\n" + g.String())

	// Test winning condition
	g.MakeMove(1, 0) // O
	g.MakeMove(0, 1) // X
	g.MakeMove(1, 1) // O
	g.MakeMove(0, 2) // X (wins)

	// Print the board after the winning move
	t.Log("\nBoard after winning move:\n" + g.String())

	if g.Winner != "X" {
		t.Errorf("Expected winner to be 'X', got %s", g.Winner)
	}
}
