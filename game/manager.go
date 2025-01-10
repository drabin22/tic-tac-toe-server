package game

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

// GameManager manages active Tic-Tac-Toe games
type GameManager struct {
	games map[string]*GameSession // Map of game IDs to GameSession objects
	mu    sync.Mutex              // Mutex for thread-safe access
}

// GameSession represents a single game session with multiple players
type GameSession struct {
	Game        *Game                    // The Tic-Tac-Toe game
	Connections map[*websocket.Conn]bool // WebSocket connections for players
	PlayerX     *websocket.Conn          // WebSocket for player "X"
	PlayerO     *websocket.Conn          // WebSocket for player "O"
	mu          sync.Mutex               // Mutex for thread-safe access
}

// NewGameManager initializes a new GameManager
func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*GameSession),
	}
}

// CreateGame creates a new game session and returns its ID
func (gm *GameManager) CreateGame(gameID string) (*GameSession, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Check if gameID already exists
	if _, exists := gm.games[gameID]; exists {
		return nil, errors.New("game ID already exists")
	}

	// Create a new game session
	session := &GameSession{
		Game:        NewGame(),
		Connections: make(map[*websocket.Conn]bool),
	}
	gm.games[gameID] = session
	return session, nil
}

// GetGame retrieves a game session by its ID
func (gm *GameManager) GetGame(gameID string) (*GameSession, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Retrieve the game session
	session, exists := gm.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}
	return session, nil
}

// DeleteGame removes a game from the manager
func (gm *GameManager) DeleteGame(gameID string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	delete(gm.games, gameID)
}

// AddConnection adds a WebSocket connection to a game session
func (gm *GameManager) AddConnection(gameID string, conn *websocket.Conn) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Retrieve the game session
	session, exists := gm.games[gameID]
	if !exists {
		return errors.New("game not found")
	}

	// Add the connection
	session.Connections[conn] = true
	return nil
}

// RemoveConnection removes a WebSocket connection from a game session
func (gm *GameManager) RemoveConnection(gameID string, conn *websocket.Conn) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Retrieve the game session
	if session, exists := gm.games[gameID]; exists {
		delete(session.Connections, conn)
	}
}

// AssignPlayer assigns a player to "X" or "O"
func (session *GameSession) AssignPlayer(conn *websocket.Conn) (string, error) {
	session.mu.Lock()
	defer session.mu.Unlock()

	if session.PlayerX == nil {
		session.PlayerX = conn
		return "X", nil
	}
	if session.PlayerO == nil {
		session.PlayerO = conn
		return "O", nil
	}
	return "", errors.New("game is already full")
}

// RemovePlayer removes a player from the game
func (session *GameSession) RemovePlayer(conn *websocket.Conn) {
	session.mu.Lock()
	defer session.mu.Unlock()

	if session.PlayerX == conn {
		session.PlayerX = nil
	}
	if session.PlayerO == conn {
		session.PlayerO = nil
	}
	delete(session.Connections, conn)
}

// GetPlayerSymbol determines if the connection is "X" or "O"
func (session *GameSession) GetPlayerSymbol(conn *websocket.Conn) (string, error) {
	session.mu.Lock()
	defer session.mu.Unlock()

	if session.PlayerX == conn {
		return "X", nil
	}
	if session.PlayerO == conn {
		return "O", nil
	}
	return "", errors.New("connection is not a player in this game")
}
