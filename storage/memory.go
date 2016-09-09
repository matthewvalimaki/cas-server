package storage

import (
	"time"

	"github.com/matthewvalimaki/cas-server/types"
)

// MemoryStorage is memory based Storage
type MemoryStorage struct {
	tickets map[string]*types.Ticket
}

// NewMemoryStorage returns new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tickets: make(map[string]*types.Ticket),
	}
}

// DoesTicketExist checks if given ticket exists
func (s *MemoryStorage) DoesTicketExist(ticket string) *types.Ticket {
	if t, ok := s.tickets[ticket]; ok {
		// check if ticket should be deleted
		if t.Old() {
			s.DeleteTicket(ticket)
			return nil
		}
		return t
	}

	return nil
}

// SaveTicket stores the Ticket
func (s *MemoryStorage) SaveTicket(ticket *types.Ticket) {
	ticket.Created = time.Now()
	s.tickets[ticket.Ticket] = ticket
}

// DeleteTicket deletes given ticket
func (s *MemoryStorage) DeleteTicket(ticket string) {
	delete(s.tickets, ticket)
}
