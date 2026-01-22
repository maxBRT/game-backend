package queue

import (
	"sync"

	m "github.com/maxbrt/game-backend/match-service/internal/models"
)

type InMemoryMatchStore struct {
	matches       map[string]*m.Match
	ticketToMatch map[string]string
	mu            sync.Mutex
}

func NewMatchStore() *InMemoryMatchStore {
	return &InMemoryMatchStore{
		matches:       make(map[string]*m.Match),
		ticketToMatch: make(map[string]string),
		mu:            sync.Mutex{},
	}
}

func (s *InMemoryMatchStore) StoreMatch(match *m.Match) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.matches[match.ID] = match

	for _, survivor := range match.Survivors {
		s.ticketToMatch[survivor.TicketID] = match.ID
	}

	s.ticketToMatch[match.Killer.TicketID] = match.ID
}

func (s *InMemoryMatchStore) GetMatch(ticketID string) (match *m.Match, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ok := s.ticketToMatch[ticketID] != ""; !ok {
		return match, false
	}

	matchID := s.ticketToMatch[ticketID]
	match = s.matches[matchID]
	return match, true
}

func (s *InMemoryMatchStore) RemoveMatch(ticketID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ok := s.ticketToMatch[ticketID] != ""; !ok {
		return false
	}

	matchID := s.ticketToMatch[ticketID]
	delete(s.matches, matchID)
	delete(s.ticketToMatch, ticketID)

	return true
}
