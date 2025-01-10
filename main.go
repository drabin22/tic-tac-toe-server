package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/drabin22/tic-tac-toe-server/game"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP requests to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (update for production)
	},
}

func main() {
	// Create a new router
	router := mux.NewRouter()

	// RESTful API endpoints
	router.HandleFunc("/stats", GetStatsHandler).Methods("GET")
	router.HandleFunc("/stats", PostStatsHandler).Methods("POST")

	// WebSocket endpoint
	router.HandleFunc("/ws", WebSocketHandler)

	// Start the HTTP server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// // Handler for WebSocket connections
// func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("Failed to upgrade connection: %v", err)
// 		return
// 	}
// 	defer conn.Close()

// 	log.Println("Client connected")
// 	// Echo messages back for now (test connection)
// 	for {
// 		messageType, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Printf("Error reading message: %v", err)
// 			break
// 		}

// 		log.Printf("Received: %s", message)
// 		if err := conn.WriteMessage(messageType, message); err != nil {
// 			log.Printf("Error writing message: %v", err)
// 			break
// 		}
// 	}
// }

var gameManager = game.NewGameManager()

// Handler for WebSocket Connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Game session variables
	var currentSession *game.GameSession
	var gameID string
	var playerSymbol string

	for {
		// Read JSON message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			if currentSession != nil && gameID != "" {
				currentSession.RemovePlayer(conn)
				broadcastToGame(currentSession, "Player disconnected", playerSymbol+" has left the game")

				// End the game if no players remain
				if len(currentSession.Connections) == 0 {
					gameManager.DeleteGame(gameID)
					log.Printf("Game %s ended due to no players remaining", gameID)
				}
			}
			break
		}

		// Parse JSON
		var msg game.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			sendError(conn, "Invalid JSON format")
			continue
		}

		switch msg.Type {
		case "create":
			// Create a new game
			gameID = msg.GameID
			currentSession, err = gameManager.CreateGame(gameID)
			if err != nil {
				sendError(conn, err.Error())
				continue
			}
			playerSymbol, _ = currentSession.AssignPlayer(conn)
			gameManager.AddConnection(gameID, conn)
			sendResponse(conn, "Game created. You are "+playerSymbol, gameID)
		case "join":
			// Join an existing game
			gameID = msg.GameID
			currentSession, err = gameManager.GetGame(gameID)
			if err != nil {
				sendError(conn, err.Error())
				continue
			}
			playerSymbol, err = currentSession.AssignPlayer(conn)
			if err != nil {
				sendError(conn, err.Error())
				continue
			}
			gameManager.AddConnection(gameID, conn)
			sendResponse(conn, "Game joined. You are "+playerSymbol, gameID)
		case "move":
			// Make a move
			if currentSession == nil {
				sendError(conn, "No game joined")
				continue
			}

			// Validate that it's the player's turn
			if playerSymbol != currentSession.Game.Turn {
				sendError(conn, "Not your turn")
				continue
			}

			err := currentSession.Game.MakeMove(msg.Row, msg.Col)
			if err != nil {
				sendError(conn, err.Error())
				continue
			}
			// Broadcast updated board to all connections in the game
			broadcastToGame(currentSession, "Board updated", currentSession.Game.String())
			if currentSession.Game.Winner != "" {
				broadcastToGame(currentSession, "Winner", currentSession.Game.Winner)
			}
		default:
			sendError(conn, "Unknown message type")
		}

	}
}

// Helper to broadcast a message to all players in a game session
func broadcastToGame(session *game.GameSession, messageType, payload string) {
	for conn := range session.Connections {
		sendResponse(conn, messageType, payload)
	}
}

// Helper to send error messages
func sendError(conn *websocket.Conn, errorMsg string) {
	msg := game.Message{
		Type:    "error",
		Payload: errorMsg,
	}
	conn.WriteJSON(msg)
}

// Helper to send response messages
func sendResponse(conn *websocket.Conn, messageType, payload string) {
	msg := game.Message{
		Type:    messageType,
		Payload: payload,
	}
	conn.WriteJSON(msg)
}

// Placeholder for retrieving stats
func GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Retrieve stats"))
}

// Placeholder for posting stats
func PostStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post stats"))
}
