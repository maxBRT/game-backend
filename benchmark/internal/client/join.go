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
	"strings"
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
		msg := string(body)
		switch {
		case strings.Contains(msg, "player id , name or role field missing"):
			return nil, fmt.Errorf("%w, message: %s", ErrPlayerFieldMissing, string(body))
		case strings.Contains(msg, "role must be either 'survivor' or 'killer'"):
			return nil, fmt.Errorf("%w, message: %s", ErrInvalidRole, string(body))
		default:
			return nil, &APIError{StatusCode: resp.StatusCode, Message: string(body)}
		}
	}

	var response JoinResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
