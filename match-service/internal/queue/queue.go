package queue

import (
	m "github.com/maxbrt/game-backend/match-service/internal/models"
)

type Queue interface {
	Enqueue(m.Player) (position int)
	Dequeue() (player m.Player, ok bool)
	IsEmpty() bool
	Len() int
	Peek(n int) (players []m.Player)
	Remove(TicketID string) bool
	Contains(TicketID string) bool
	GetPosition(TicketID string) (position int, ok bool)
}
