package security

import (
    "fmt"
    "math/rand"
    
    "github.com/matthewvalimaki/cas-server/tools"
    "github.com/matthewvalimaki/cas-server/types"
    "github.com/matthewvalimaki/cas-server/storage"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// CreateNewServiceTicket creates a new Service Ticket
// see: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html#service-ticket
func CreateNewServiceTicket(strg storage.IStorage, srvc string) (ticket *types.Ticket, err error) {
	c := 100
    b := make([]byte, c)
    for i := range b {
        b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
    }

    var randomST = "ST-" + string(b)
    
    newTicket := &types.Ticket{Ticket: randomST, Service: srvc}
    
    strg.SaveNewServiceTicket(newTicket)
    
    tools.LogST(newTicket, "ticket created")
    
    return newTicket, nil
}

// ValidateServiceTicket validates checks if given Service Ticket exists and is valid
func ValidateServiceTicket(strg storage.IStorage, ticket *types.Ticket) *types.CasError {   
    tools.LogST(ticket, "ticket validation requested")
    
    if strg.DoesServiceTicketExist(ticket.Ticket) {
        tools.LogST(ticket, "ticket validation succeeded (ticket was found)")

        strg.DeleteServiceTicket(ticket.Ticket)

        tools.LogST(ticket, "deleting ticket")
        
        return nil
    }
    
    tools.LogST(ticket, "ticket validation failed (ticket was not found)")
    
    return &types.CasError{Error: fmt.Errorf("The ticket `%s` is invalid.", ticket.Ticket), CasErrorCode: types.CAS_ERROR_CODE_INVALID_TICKET}
}
