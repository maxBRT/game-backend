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

// TestPurchase_Success verifies that a successful purchase request:
// - Uses POST method with application/json Content-Type
// - Correctly parses the response into PurchaseResult
// - Returns the expected Success, NewBalance, and ItemName fields
func TestPurchase_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and headers
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Verify request body can be decoded
		var req PurchaseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		// Return successful purchase response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PurchaseResponse{
			Success:    true,
			NewBalance: 900,
			ItemName:   "Health Potion",
		})
	}))
	defer server.Close()

	config := &cfg.Config{
		PlayerServiceURL:    server.URL,
		HTTPTimeout:         5 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
	}

	client := NewStoreClient(config)

	req := PurchaseRequest{
		PlayerId: 1,
		ItemId:   100,
		Amount:   1,
	}

	result, err := client.Purchase(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if !result.Result.Success {
		t.Errorf("expected Success to be true, got false")
	}
	if result.Result.NewBalance != 900 {
		t.Errorf("expected NewBalance 900, got %d", result.Result.NewBalance)
	}
	if result.Result.ItemName != "Health Potion" {
		t.Errorf("expected ItemName 'Health Potion', got %s", result.Result.ItemName)
	}
}

// TestPurchase_InsufficientFunds verifies that a 422 Unprocessable Entity response
// is correctly mapped to ErrInsufficientFunds error.
func TestPurchase_InsufficientFunds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("insufficient funds for purchase"))
	}))
	defer server.Close()

	config := &cfg.Config{
		PlayerServiceURL:    server.URL,
		HTTPTimeout:         5 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
	}

	client := NewStoreClient(config)

	req := PurchaseRequest{
		PlayerId: 1,
		ItemId:   100,
		Amount:   999,
	}

	result, err := client.Purchase(context.Background(), req)

	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}
}

// TestPurchase_ItemNotFound verifies that a 400 Bad Request response
// is correctly mapped to ErrItemNotFound error.
func TestPurchase_ItemNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("item not found"))
	}))
	defer server.Close()

	config := &cfg.Config{
		PlayerServiceURL:    server.URL,
		HTTPTimeout:         5 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
	}

	client := NewStoreClient(config)

	req := PurchaseRequest{
		PlayerId: 1,
		ItemId:   99999, // Non-existent item
		Amount:   1,
	}

	result, err := client.Purchase(context.Background(), req)

	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrItemNotFound) {
		t.Errorf("expected ErrItemNotFound, got %v", err)
	}
}

// TestPurchase_PlayerNotFound verifies that a 404 Not Found response
// is correctly mapped to ErrPlayerNotFound error.
func TestPurchase_PlayerNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("player not found"))
	}))
	defer server.Close()

	config := &cfg.Config{
		PlayerServiceURL:    server.URL,
		HTTPTimeout:         5 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
	}

	client := NewStoreClient(config)

	req := PurchaseRequest{
		PlayerId: 99999, // Non-existent player
		ItemId:   100,
		Amount:   1,
	}

	result, err := client.Purchase(context.Background(), req)

	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrPlayerNotFound) {
		t.Errorf("expected ErrPlayerNotFound, got %v", err)
	}
}

// TestPurchase_ServerError verifies that a 500 Internal Server Error response
// returns an error containing the status code for debugging purposes.
func TestPurchase_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	config := &cfg.Config{
		PlayerServiceURL:    server.URL,
		HTTPTimeout:         5 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
	}

	client := NewStoreClient(config)

	req := PurchaseRequest{
		PlayerId: 1,
		ItemId:   100,
		Amount:   1,
	}

	result, err := client.Purchase(context.Background(), req)

	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	errMsg := err.Error()
	if !contains(errMsg, "500") {
		t.Errorf("expected error message to contain '500', got %s", errMsg)
	}
}

// TestPurchase_LatencyTracking verifies that the client accurately measures
// request latency. Uses an artificial 50ms delay and checks that the measured
// latency falls within expected bounds (>= delay, < delay + 100ms overhead).
func TestPurchase_LatencyTracking(t *testing.T) {
	delay := 50 * time.Millisecond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate network/processing delay
		time.Sleep(delay)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PurchaseResponse{
			Success:    true,
			NewBalance: 900,
			ItemName:   "Health Potion",
		})
	}))
	defer server.Close()

	config := &cfg.Config{
		PlayerServiceURL:    server.URL,
		HTTPTimeout:         5 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
	}

	client := NewStoreClient(config)

	req := PurchaseRequest{
		PlayerId: 1,
		ItemId:   100,
		Amount:   1,
	}

	result, err := client.Purchase(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Latency < delay {
		t.Errorf("expected latency >= %v, got %v", delay, result.Latency)
	}
	maxExpectedLatency := delay + 100*time.Millisecond
	if result.Latency > maxExpectedLatency {
		t.Errorf("expected latency < %v, got %v", maxExpectedLatency, result.Latency)
	}
}

// TestPurchase_ContextCancellation verifies that the client respects context
// timeouts. Creates a context with 50ms timeout against a server that takes
// 500ms to respond, expecting context.DeadlineExceeded error.
func TestPurchase_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay longer than the context timeout
		time.Sleep(500 * time.Millisecond)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PurchaseResponse{
			Success:    true,
			NewBalance: 900,
			ItemName:   "Health Potion",
		})
	}))
	defer server.Close()

	config := &cfg.Config{
		PlayerServiceURL:    server.URL,
		HTTPTimeout:         5 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
	}

	client := NewStoreClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	req := PurchaseRequest{
		PlayerId: 1,
		ItemId:   100,
		Amount:   1,
	}

	result, err := client.Purchase(ctx, req)

	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

// TestPurchase_ContextAlreadyCancelled verifies that the client returns
// immediately with context.Canceled if the context is already cancelled
// before the request is made. The server handler should never be invoked.
func TestPurchase_ContextAlreadyCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This should never be reached
		t.Error("Handler should not be called when context is already cancelled")
	}))
	defer server.Close()

	config := &cfg.Config{
		PlayerServiceURL:    server.URL,
		HTTPTimeout:         5 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
	}

	client := NewStoreClient(config)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := PurchaseRequest{
		PlayerId: 1,
		ItemId:   100,
		Amount:   1,
	}

	result, err := client.Purchase(ctx, req)

	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

// TestPurchase_RequestBodyFormat verifies that the client correctly serializes
// PurchaseRequest to JSON with the expected field names (playerId, itemId, amount).
func TestPurchase_RequestBodyFormat(t *testing.T) {
	var receivedRequest PurchaseRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the request body for verification
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PurchaseResponse{
			Success:    true,
			NewBalance: 900,
			ItemName:   "Test Item",
		})
	}))
	defer server.Close()

	config := &cfg.Config{
		PlayerServiceURL:    server.URL,
		HTTPTimeout:         5 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
	}

	client := NewStoreClient(config)

	req := PurchaseRequest{
		PlayerId: 42,
		ItemId:   123,
		Amount:   5,
	}

	_, _ = client.Purchase(context.Background(), req)

	if receivedRequest.PlayerId != 42 {
		t.Errorf("expected PlayerId 42, got %d", receivedRequest.PlayerId)
	}
	if receivedRequest.ItemId != 123 {
		t.Errorf("expected ItemId 123, got %d", receivedRequest.ItemId)
	}
	if receivedRequest.Amount != 5 {
		t.Errorf("expected Amount 5, got %d", receivedRequest.Amount)
	}
}

// contains checks if substr exists within s.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
