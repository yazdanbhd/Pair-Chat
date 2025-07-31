package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"

	"pairchat/internal/model"
)

// CreateSession is the handler for starting a new chat session.
func CreateSession(hub *model.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := uuid.New().String()
		hub.Sessions[sessionID] = &model.Session{}

		response := map[string]string{
			"message":   "New session created.",
			"sessionID": sessionID,
			"link_1":    "/chat/" + sessionID + "/1",
			"link_2":    "/chat/" + sessionID + "/2",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding json response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// WebSocket is the handler for upgrading the connection and managing messages.
func WebSocket(hub *model.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		sessionID := parts[2]
		userID := parts[3]

		wsHandler := websocket.Handler(func(ws *websocket.Conn) {
			defer func() {
				if err := ws.Close(); err != nil {
					log.Printf("Error closing websocket for user %s: %v", userID, err)
				}
			}()

			log.Printf("User %s connected to session %s", userID, sessionID)

			session, exists := hub.Sessions[sessionID]
			if !exists {
				log.Printf("Session %s not found", sessionID)
				return
			}

			if userID == "1" {
				session.User1 = ws
			} else {
				session.User2 = ws
			}

			for {
				var msg string
				if err := websocket.Message.Receive(ws, &msg); err != nil {
					log.Printf("Connection closed for user %s in session %s: %v", userID, sessionID, err)
					break
				}

				log.Printf("Received message from User %s in Session %s: %s", userID, sessionID, msg)

				var recipient *websocket.Conn
				if userID == "1" {
					recipient = session.User1
				} else {
					recipient = session.User2
				}

				if recipient != nil {
					if err := websocket.Message.Send(recipient, msg); err != nil {
						log.Printf("Failed to send message to other user: %v", err)
						break
					}
				}
			}
		})

		wsHandler.ServeHTTP(w, r)
	}
}
