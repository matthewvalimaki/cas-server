package storage

import (
    "github.com/matthewvalimaki/cas-server/types"
)

// MemoryStorage is memory based Storage
type MemoryStorage struct {
    serviceTickets map[string] *types.Ticket
}

// NewMemoryStorage returns new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{serviceTickets: make(map[string] *types.Ticket)}
}

// DoesServiceTicketExist checks if given Service Ticket exists
func (s MemoryStorage) DoesServiceTicketExist(st string) bool {   
    if _, ok := s.serviceTickets[st]; ok {
        return true
    }
    
    return false
}

// SaveNewServiceTicket stores the Ticket
func (s MemoryStorage) SaveNewServiceTicket(ticket *types.Ticket) {
    var st = ticket.Ticket
    s.serviceTickets[st] = ticket
}

// DeleteServiceTicket deletes given ticket
func (s MemoryStorage) DeleteServiceTicket(st string) {
    delete(s.serviceTickets, st)
}
