package game

import (
	"errors"
	"sync"
)

// GameManager manages active Tic-Tac-Toe games
type GameManager struct {
	games map[string]*Game // Map of game IDs to Game objects
	mu    sync.Mutex       // Mutex for thread-safe access
}

// NewGameManager initializes a new GameManager
func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*Game),
	}
}

// CreateGame creates a new game and returns its ID
func (gm *GameManager) CreateGame(gameID string) (*Game, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Check if gameID already exists
	if _, exists := gm.games[gameID]; exists {
		return nil, errors.New("game ID already exists")
	}

	// Create a new game
	game := NewGame()
	gm.games[gameID] = game
	return game, nil
}

// GetGame retrieves a game by its ID
func (gm *GameManager) GetGame(gameID string) (*Game, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Retrieve the game
	game, exists := gm.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}
	return game, nil
}

// DeleteGame removes a game from the manager
func (gm *GameManager) DeleteGame(gameID string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	delete(gm.games, gameID)
}
