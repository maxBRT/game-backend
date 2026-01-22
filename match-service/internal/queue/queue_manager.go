package queue

import (
	"sync"

	m "github.com/maxbrt/game-backend/match-service/internal/models"
)

type QueueManager struct {
	survivorsQueue Queue
	killerQueue    Queue
	store          MatchStore
	mu             sync.Mutex
}

func NewQueueManager(survivorsQ Queue, killerQ Queue, store MatchStore) *QueueManager {
	return &QueueManager{
		survivorsQueue: survivorsQ,
		killerQueue:    killerQ,
		store:          store,
	}
}

func (m *QueueManager) AddPlayer(player m.Player) {
	if player.IsSurvivor() {
		m.survivorsQueue.Enqueue(player)
	} else {
		m.killerQueue.Enqueue(player)
	}
}

func (m *QueueManager) GetPlayerStatus(TicketID string) (matched bool, position int, match *m.Match) {
	match, ok := m.store.GetMatch(TicketID)
	if ok {
		return true, -1, match
	}

	position, ok = m.survivorsQueue.GetPosition(TicketID)
	if ok {
		return true, position, nil
	}

	position, ok = m.killerQueue.GetPosition(TicketID)
	if ok {
		return true, position, nil
	}

	return false, -1, nil

}
