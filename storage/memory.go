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

// DoesServiceTicketExist checks if given Service Ticket exists
func (s MemoryStorage) DoesServiceTicketExist(st string) bool {   
    if _, ok := s.tickets[st]; ok {
        return true
    }
    
    return false
}

// SaveTicket stores the Ticket
func (s MemoryStorage) SaveTicket(ticket *types.Ticket) {
    s.tickets[ticket.Ticket] = ticket
}

// DeleteServiceTicket deletes given ticket
func (s MemoryStorage) DeleteServiceTicket(st string) {
    delete(s.tickets, st)
}
