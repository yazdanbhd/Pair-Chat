package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pairchat/internal/model"
)

func TestCreateSession(t *testing.T) {
	hub := &model.Hub{
		Sessions: make(map[string]*model.Session),
	}

	req := httptest.NewRequest("POST", "/session/create", nil)
	rr := httptest.NewRecorder()

	handler := CreateSession(hub)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedContentType := "application/json"
	if ctype := rr.Header().Get("Content-Type"); ctype != expectedContentType {
		t.Errorf("handler returned wrong content type: got %s want %s",
			ctype, expectedContentType)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not decode response body: %v", err)
	}

	sessionID, ok := response["sessionID"]
	if !ok {
		t.Fatal("Response did not contain sessionID")
	}
	if _, exists := hub.Sessions[sessionID]; !exists {
		t.Errorf("handler did not create a session in the hub")
	}
}
