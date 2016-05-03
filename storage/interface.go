package storage

import (
    "github.com/matthewvalimaki/cas-server/types"
)

// IStorage interface for all Storages
type IStorage interface {
    SaveNewServiceTicket(*types.Ticket)
    DoesServiceTicketExist(st string) bool
    DeleteServiceTicket(st string)
}
