package queue

import (
	"time"

	"github.com/google/uuid"
	m "github.com/maxbrt/game-backend/match-service/internal/models"
)

type Matcher struct {
	QueueManager  *QueueManager
	stopChan      chan struct{}
	matchInterval time.Duration
}

func NewMatcher(queueManager *QueueManager, matchInterval time.Duration) *Matcher {
	return &Matcher{
		QueueManager:  queueManager,
		stopChan:      make(chan struct{}),
		matchInterval: matchInterval,
	}
}

func (m *Matcher) Start() {
	go func() {
		for {
			select {
			case <-m.stopChan:
				return
			case <-time.After(m.matchInterval):
				m.TryMatch()
			}
		}
	}()
}

func (m *Matcher) Stop() {
	close(m.stopChan)
}

func (m *Matcher) TryMatch() {
	m.QueueManager.mu.Lock()
	defer m.QueueManager.mu.Unlock()

	if m.QueueManager.survivorsQueue.Len() >= 4 && m.QueueManager.killerQueue.Len() >= 1 {
		survivors := m.QueueManager.survivorsQueue.Peek(4)
		killer := m.QueueManager.killerQueue.Peek(1)

		match := newMatch(survivors, killer[0])

		m.QueueManager.store.StoreMatch(match)

		for _, survivor := range survivors {
			m.QueueManager.survivorsQueue.Remove(survivor.TicketID)
		}

		m.QueueManager.killerQueue.Remove(killer[0].TicketID)

		m.QueueManager.survivorsQueue.Recalculate()
		m.QueueManager.killerQueue.Recalculate()
	}

}

func newMatch(survivors []m.Player, killer m.Player) *m.Match {
	return &m.Match{
		ID:        uuid.New().String(),
		Survivors: survivors,
		Killer:    killer,
		CreatedAt: time.Now().String(),
	}
}
