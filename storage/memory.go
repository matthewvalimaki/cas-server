package storage

import (
    "github.com/matthewvalimaki/cas-server/types"
)

// MemoryStorage is memory based Storage
type MemoryStorage struct {
    tickets map[string]*types.Ticket
}

// NewMemoryStorage returns new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{tickets: make(map[string]*types.Ticket)}
}

// DoesTicketExist checks if given ticket exists
func (s MemoryStorage) DoesTicketExist(ticket string) bool {   
    if _, ok := s.tickets[ticket]; ok {
        return true
    }
    
    return false
}

// SaveTicket stores the Ticket
func (s MemoryStorage) SaveTicket(ticket *types.Ticket) {
    s.tickets[ticket.Ticket] = ticket
}

// DeleteTicket deletes given ticket
func (s MemoryStorage) DeleteTicket(ticket string) {
    delete(s.tickets, ticket)
}
