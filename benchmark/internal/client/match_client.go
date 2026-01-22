package client

import (
	cfg "github.com/maxbrt/game-backend/benchmark/internal/config"
	"net/http"
)

type MatchClient struct {
	client http.Client
	config *cfg.Config
}

func NewMatchClient(config *cfg.Config) *MatchClient {
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
	}

	return &MatchClient{
		client: http.Client{
			Timeout:   config.HTTPTimeout,
			Transport: transport,
		},
		config: config,
	}
}

type JoinRequest struct {
	Player PlayerInfo `json:"player"`
}

type JoinResponse struct {
	Success  bool   `json:"success"`
	QueuePos int    `json:"queuePos"`
	TicketID string `json:"ticketID"`
	Role     string `json:"role"`
}

type StatusResponse struct {
	Status        string       `json:"status"`
	QueuePosition int          `json:"queuePosition,omitempty"`
	MatchID       string       `json:"matchID,omitempty"`
	Survivors     []PlayerInfo `json:"survivors,omitempty"`
	Killer        PlayerInfo   `json:"killer"`
}

type PlayerInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}
