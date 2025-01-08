package game

// MEssage represents a JSON message sent between the client and server
type Message struct {
	Type    string `json:"type"`    // Message type: "create", "join", "move"
	GameID  string `json:"gameID"`  // Game ID (if applicable)
	Row     int    `json:"row"`     // Row for a move (if applicable)
	Col     int    `json:"col"`     // Column for a move (if applicable)
	Payload string `json:"payload"` // Additional data (e.g., error messages, board state, etc.)
}
