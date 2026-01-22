package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *MatchClient) Status(ctx context.Context, ticketID string) (*StatusResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	url := c.config.MatchMakingServiceURL + "/match/status/" + ticketID
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

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
		case http.StatusBadRequest:
			return nil, fmt.Errorf("%w, message: %s", ErrInvalidTicketID, string(body))
		default:
			return nil, &APIError{StatusCode: resp.StatusCode, Message: string(body)}
		}
	}

	var response StatusResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
