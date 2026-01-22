package queue

import (
	m "github.com/maxbrt/game-backend/match-service/internal/models"
	"sync"
)

type InMemoryQueue struct {
	players    []m.Player
	ticketsIDs map[string]int
	mu         sync.Mutex
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		players:    make([]m.Player, 0),
		ticketsIDs: make(map[string]int),
		mu:         sync.Mutex{},
	}
}

func (q *InMemoryQueue) Enqueue(player m.Player) (position int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.players = append(q.players, player)
	q.ticketsIDs[player.TicketID] = len(q.players) - 1
	return len(q.players) - 1
}

func (q *InMemoryQueue) Dequeue() (player m.Player, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.players) == 0 {
		return player, false
	}

	player = q.players[0]
	q.players = q.players[1:]
	return player, true
}

func (q *InMemoryQueue) IsEmpty() bool {
	return len(q.players) == 0
}

func (q *InMemoryQueue) Len() int {
	return len(q.players)
}
func (q *InMemoryQueue) Peek(n int) (players []m.Player) {
	return q.players[:n]
}

func (q *InMemoryQueue) Remove(TicketID string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if ok := q.Contains(TicketID); !ok {
		return false
	}

	i := q.ticketsIDs[TicketID]
	q.players = append(q.players[:i], q.players[i+1:]...)
	delete(q.ticketsIDs, TicketID)

	return true
}

func (q *InMemoryQueue) Contains(TicketID string) bool {
	return q.ticketsIDs[TicketID] != -1
}

func (q *InMemoryQueue) GetPosition(TicketID string) (position int, ok bool) {
	return 0, false
}
