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
)

func (c *MatchClient) Join(ctx context.Context, player PlayerInfo) (*JoinResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	url := c.config.MatchMakingServiceURL + "/match/join"
	jsonBody, err := json.Marshal(JoinRequest{
		Player: player,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
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
		return nil, fmt.Errorf("join failed (status %d): %s", resp.StatusCode, body)
	}

	var response JoinResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
