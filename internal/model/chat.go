package model

import "golang.org/x/net/websocket"

// Hub is a temporary, in-memory store for all active chat sessions.
// The key of the map will be the session ID.
type Hub struct {
	Sessions map[string]*Session
}

// Session represents a single two-person chat room.
// It holds the WebSocket connections for the two participants.
type Session struct {
	User1 *websocket.Conn
	User2 *websocket.Conn
}
