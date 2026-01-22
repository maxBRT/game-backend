package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	cfg "github.com/maxbrt/game-backend/benchmark/internal/config"
)

type StoreClient struct {
	client http.Client
	config *cfg.Config
}

func NewStoreClient(config *cfg.Config) *StoreClient {
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
	}

	return &StoreClient{
		client: http.Client{
			Timeout:   config.HTTPTimeout,
			Transport: transport,
		},
		config: config,
	}
}

func (c *StoreClient) Purchase(ctx context.Context, request PurchaseRequest) (*PurchaseResult, error) {
	start := time.Now()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	url := c.config.PlayerServiceURL + "/store/purchase"
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {

		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Request timed out")
			return nil, err
		}

		if errors.Is(err, context.Canceled) {
			log.Println("Request canceled")
			return nil, err
		}

		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		switch resp.StatusCode {
		case http.StatusUnprocessableEntity:
			return nil, fmt.Errorf("%w, message: %s", ErrInsufficientFunds, string(body))
		case http.StatusNotFound:
			return nil, fmt.Errorf("%w, message: %s", ErrPlayerNotFound, string(body))
		case http.StatusBadRequest:
			return nil, fmt.Errorf("%w, message: %s", ErrItemNotFound, string(body))
		default:
			return nil, fmt.Errorf("purchase failed (status %d): %s", resp.StatusCode, string(body))
		}
	}

	var response PurchaseResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &PurchaseResult{
		Result:  response,
		Latency: time.Since(start),
	}, nil
}

type PurchaseRequest struct {
	PlayerId int `json:"playerId"`
	ItemId   int `json:"itemId"`
	Amount   int `json:"amount"`
}

type PurchaseResponse struct {
	Success    bool   `json:"success"`
	NewBalance int    `json:"newBalance"`
	ItemName   string `json:"itemName"`
}

type PurchaseResult struct {
	Result  PurchaseResponse `json:"result"`
	Latency time.Duration    `json:"latency"`
}
