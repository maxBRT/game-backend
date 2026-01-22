package queue

import (
	m "github.com/maxbrt/game-backend/match-service/internal/models"
)

type MatchStore interface {
	StoreMatch(match *m.Match, ticketID string)
	GetMatch(ticketID string) (match *m.Match, ok bool)
	RemoveMatch(ticketID string) bool
}
