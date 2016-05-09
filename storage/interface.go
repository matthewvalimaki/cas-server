package storage

import (
    "github.com/matthewvalimaki/cas-server/types"
)

// IStorage interface for all Storages
type IStorage interface {
    SaveTicket(*types.Ticket)
    DoesTicketExist(ticket string) bool
    DeleteTicket(ticket string)
}
