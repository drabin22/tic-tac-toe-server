package main

import (
	"log"
	"net/http"

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

// Handler for WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")
	// Echo messages back for now (test connection)
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received: %s", message)
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
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
