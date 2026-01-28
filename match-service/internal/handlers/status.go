package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	m "github.com/maxbrt/game-backend/match-service/internal/models"
	q "github.com/maxbrt/game-backend/match-service/internal/queue"
	"github.com/maxbrt/game-backend/match-service/internal/utils"
)

type WaitingResponse struct {
	Status        string `json:"status"`
	QueuePosition int    `json:"queuePosition"`
}

type MatchResponse struct {
	Status    string         `json:"status"`
	MatchID   string         `json:"matchID"`
	Survivors []m.PlayerInfo `json:"survivors"`
	Killer    m.PlayerInfo   `json:"killer"`
}

type StatusHandler struct {
	QueueManager *q.QueueManager
	pollTimeout  time.Duration
	pollInterval time.Duration
}

func NewStatusHandler(queueManager *q.QueueManager, pollTimeout, pollInterval time.Duration) *StatusHandler {
	return &StatusHandler{
		QueueManager: queueManager,
		pollTimeout:  pollTimeout,
		pollInterval: pollInterval,
	}
}

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ticketID := r.PathValue("ticketID")
	if ticketID == "" {
		fmt.Println("TicketID is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Checking ticket: %s\n", ticketID)

	if ok := h.QueueManager.Contains(ticketID); !ok {
		fmt.Printf("ticketID not found: %s\n", ticketID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), h.pollTimeout)
	defer cancel()

	ticker := time.NewTicker(h.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				_, position, _ := h.QueueManager.GetPlayerStatus(ticketID)
				resp := WaitingResponse{
					Status:        "waiting",
					QueuePosition: position,
				}
				json.NewEncoder(w).Encode(resp)
			}
			// Return on any context cancellation (timeout or client disconnect)
			return
		case <-ticker.C:
			// Check if request context was cancelled (client disconnect)
			if r.Context().Err() != nil {
				return
			}
			matched, _, match := h.QueueManager.GetPlayerStatus(ticketID)
			if matched {
				resp := MatchResponse{
					Status:    "matched",
					MatchID:   match.ID,
					Survivors: utils.PlayersToInfo(match.Survivors),
					Killer:    utils.PlayerToInfo(match.Killer),
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	}
}
