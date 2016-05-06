package security

import (
    "fmt"
    "math/rand"
    
    "github.com/matthewvalimaki/cas-server/tools"
    "github.com/matthewvalimaki/cas-server/types"
    "github.com/matthewvalimaki/cas-server/storage"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// createNewTicket creates a new ticket
func createNewTicket(ticketType string) types.Ticket {
	c := 100
    b := make([]byte, c)
    for i := range b {
        b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
    }
    
    return types.Ticket{Ticket: ticketType + "-" + string(b)} 
}

// CreateNewProxyGrantingTicket creates new proxy granting ticket for a service
func CreateNewProxyGrantingTicket(strg storage.IStorage, service *types.Service, proxyService string) (*types.Ticket, *types.CasError) {
    ticket := createNewTicket("PGTIOU")
    
    return &ticket, nil
}

// CreateNewServiceTicket creates a new Service Ticket
// see: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html#service-ticket
func CreateNewServiceTicket(strg storage.IStorage, serviceID string) (*types.Ticket, error) {
	ticket := createNewTicket("ST")
    ticket.Service = serviceID
    
    strg.SaveNewServiceTicket(&ticket)
    
    tools.LogST(&ticket, "ticket created")
    
    return &ticket, nil
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
