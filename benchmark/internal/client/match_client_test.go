package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cfg "github.com/maxbrt/game-backend/benchmark/internal/config"
)

// testConfig creates a Config pointing to the given test server URL
func testConfig(serverURL string) *cfg.Config {
	return &cfg.Config{
		MatchMakingServiceURL: serverURL,
		HTTPTimeout:           5 * time.Second,
		MaxIdleConns:          10,
		MaxIdleConnsPerHost:   10,
		MaxConnsPerHost:       10,
	}
}

// Verifies Join sends correct request and parses response.
func TestJoin_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/match/join" {
			t.Errorf("expected /match/join, got %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Validate the request body was serialized correctly
		var req JoinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if req.Player.ID != "player-1" {
			t.Errorf("expected player ID 'player-1', got %s", req.Player.ID)
		}
		if req.Player.Name != "TestPlayer" {
			t.Errorf("expected player name 'TestPlayer', got %s", req.Player.Name)
		}
		if req.Player.Role != "survivor" {
			t.Errorf("expected player role 'survivor', got %s", req.Player.Role)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "queuePos": 5, "ticketID": "ticket-abc123", "role": "survivor"}`))
	}))
	defer server.Close()

	client := NewMatchClient(testConfig(server.URL))
	resp, err := client.Join(context.Background(), PlayerInfo{
		ID:   "player-1",
		Name: "TestPlayer",
		Role: "survivor",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !resp.Success {
		t.Error("expected Success to be true")
	}
	if resp.TicketID != "ticket-abc123" {
		t.Errorf("expected ticketID 'ticket-abc123', got %s", resp.TicketID)
	}
	if resp.QueuePos != 5 {
		t.Errorf("expected queuePos 5, got %d", resp.QueuePos)
	}
	if resp.Role != "survivor" {
		t.Errorf("expected role 'survivor', got %s", resp.Role)
	}
}

// Verifies Join returns error when player ID is missing.
func TestJoin_MissingPlayerFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/match/join" {
			t.Errorf("expected /match/join, got %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		var req JoinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if req.Player.ID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "queuePos": 5, "ticketID": "ticket-abc123", "role": "survivor"}`))
	}))
	defer server.Close()

	client := NewMatchClient(testConfig(server.URL))
	_, err := client.Join(context.Background(), PlayerInfo{
		Name: "TestPlayer",
		Role: "survivor",
	})
	if err == nil {
		t.Fatal("expected error for missing player ID")
	}
}

// Verifies Join returns error when role is not "survivor" or "killer".
func TestJoin_InvalidRole(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/match/join" {
			t.Errorf("expected /match/join, got %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		var req JoinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if req.Player.Role != "survivor" && req.Player.Role != "killer" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "queuePos": 5, "ticketID": "ticket-abc123", "role": "survivor"}`))
	}))
	defer server.Close()

	client := NewMatchClient(testConfig(server.URL))
	_, err := client.Join(context.Background(), PlayerInfo{
		ID:   "player-1",
		Name: "TestPlayer",
		Role: "invalid",
	})
	if err == nil {
		t.Fatal("expected error for invalid role")
	}
}

// Verifies Status parses "waiting" response with queue position.
func TestStatus_WaitingStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/match/status/ticket-abc123" {
			t.Errorf("expected /match/status/ticket-abc123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "waiting", "queuePosition": 5}`))
	}))
	defer server.Close()

	client := NewMatchClient(testConfig(server.URL))
	resp, err := client.Status(context.Background(), "ticket-abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "waiting" {
		t.Errorf("expected status 'waiting', got %s", resp.Status)
	}
	if resp.QueuePosition != 5 {
		t.Errorf("expected queue position 5, got %d", resp.QueuePosition)
	}
}

// Verifies Status parses "matched" response with match details.
func TestStatus_MatchedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/match/status/ticket-abc123" {
			t.Errorf("expected /match/status/ticket-abc123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "matched", "matchID": "match-abc123", "survivors": [{"id": "survivor-1", "name": "TestSurvivor"}], "killer": {"id": "killer-1", "name": "TestKiller"}}`))
	}))
	defer server.Close()

	client := NewMatchClient(testConfig(server.URL))
	resp, err := client.Status(context.Background(), "ticket-abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "matched" {
		t.Errorf("expected status 'matched', got %s", resp.Status)
	}
	if resp.MatchID != "match-abc123" {
		t.Errorf("expected match ID 'match-abc123', got %s", resp.MatchID)
	}
	if len(resp.Survivors) != 1 {
		t.Errorf("expected 1 survivor, got %d", len(resp.Survivors))
	}
	if resp.Survivors[0].ID != "survivor-1" {
		t.Errorf("expected survivor ID 'survivor-1', got %s", resp.Survivors[0].ID)
	}
	if resp.Survivors[0].Name != "TestSurvivor" {
		t.Errorf("expected survivor name 'TestSurvivor', got %s", resp.Survivors[0].Name)
	}
	if resp.Killer.ID != "killer-1" {
		t.Errorf("expected killer ID 'killer-1', got %s", resp.Killer.ID)
	}
	if resp.Killer.Name != "TestKiller" {
		t.Errorf("expected killer name 'TestKiller', got %s", resp.Killer.Name)
	}
}

// Verifies Status returns error for invalid ticket ID.
func TestStatus_InvalidTicketID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"statusCode": 400, "message": "invalid ticket id"}`))
	}))
	defer server.Close()

	client := NewMatchClient(testConfig(server.URL))
	_, err := client.Status(context.Background(), "invalid-ticket")
	if err == nil {
		t.Fatal("expected error for invalid ticket ID")
	}
}

// Verifies Join returns context.Canceled when context is cancelled mid-request.
func TestJoin_ContextCancellation(t *testing.T) {
	serverHit := make(chan struct{})
	unblock := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHit <- struct{}{} // signal that request arrived
		<-unblock               // wait until test unblocks us
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "ticketID": "ticket-abc123"}`))
	}))
	defer server.Close()
	defer close(unblock) // unblock server on test exit

	client := NewMatchClient(testConfig(server.URL))
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		_, err := client.Join(ctx, PlayerInfo{
			ID:   "player-1",
			Name: "TestPlayer",
			Role: "survivor",
		})
		errCh <- err
	}()

	<-serverHit // wait for request to reach server
	cancel()    // cancel context while request is in flight

	err := <-errCh
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

// Verifies Status returns context.Canceled when context is cancelled mid-request.
func TestStatus_ContextCancellation(t *testing.T) {
	serverHit := make(chan struct{})
	unblock := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHit <- struct{}{} // signal that request arrived
		<-unblock               // wait until test unblocks us
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "waiting", "queuePosition": 5}`))
	}))
	defer server.Close()
	defer close(unblock) // unblock server on test exit

	client := NewMatchClient(testConfig(server.URL))
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		_, err := client.Status(ctx, "ticket-abc123")
		errCh <- err
	}()

	<-serverHit // wait for request to reach server
	cancel()    // cancel context while request is in flight

	err := <-errCh
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
