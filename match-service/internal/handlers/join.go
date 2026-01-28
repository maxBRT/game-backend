package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	m "github.com/maxbrt/game-backend/match-service/internal/models"
	q "github.com/maxbrt/game-backend/match-service/internal/queue"
)

type JoinRequest struct {
	Player m.PlayerInfo `json:"player"`
}

type JoinResponse struct {
	Success  bool   `json:"success"`
	QueuePos int    `json:"queuePos"`
	TicketID string `json:"ticketID"`
	Role     string `json:"role"`
}

type JoinHandler struct {
	QueueManager *q.QueueManager
}

func NewJoinHandler(queueManager *q.QueueManager) *JoinHandler {
	return &JoinHandler{
		QueueManager: queueManager,
	}
}

func (h *JoinHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req JoinRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playerId := req.Player.ID
	playerName := req.Player.Name
	playerRole := req.Player.Role

	if playerId == "" || playerName == "" || playerRole == "" {
		fmt.Println("player id, name and role are required")
		http.Error(w, fmt.Sprintf("player id, name and role are required"), http.StatusBadRequest)
		return
	}

	if playerRole != "survivor" && playerRole != "killer" {
		fmt.Println("player role must be either survivor or killer")
		http.Error(w, fmt.Sprintf("player role must be either survivor or killer"), http.StatusBadRequest)
		return
	}

	ticketID := uuid.New().String()

	p := m.Player{
		ID:       playerId,
		Name:     playerName,
		Role:     playerRole,
		Position: -1,
		TicketID: ticketID,
	}

	position := h.QueueManager.AddPlayer(p)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := JoinResponse{
		Success:  true,
		QueuePos: position,
		TicketID: ticketID,
		Role:     playerRole,
	}

	log.Printf("player %s joined queue at position %d as %s", playerName, position, playerRole)

	json.NewEncoder(w).Encode(resp)
}
