package security

import (
	"fmt"
	"math/rand"

	"github.com/matthewvalimaki/cas-server/storage"
	"github.com/matthewvalimaki/cas-server/tools"
	"github.com/matthewvalimaki/cas-server/types"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// createNewTicket creates a new ticket
func createNewTicket(ticketType string) types.Ticket {
	c := 100
	b := make([]byte, c)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return types.Ticket{
		Ticket: ticketType + "-" + string(b),
		Type:   ticketType,
	}
}

// CreateNewProxyGrantingTicket creates new PGT
func CreateNewProxyGrantingTicket() (*types.Ticket, *types.CasError) {
	ticket := createNewTicket("PGT")

	return &ticket, nil
}

// CreateNewProxyGrantingTicketIOU creates new proxy granting ticket for a service
// see: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html#proxy-granting-ticket-iou
func CreateNewProxyGrantingTicketIOU() (*types.Ticket, *types.CasError) {
	ticket := createNewTicket("PGTIOU")

	return &ticket, nil
}

// CreateNewProxyTicket creates new Proxy Ticket (PT)
func CreateNewProxyTicket() (*types.Ticket, *types.CasError) {
	ticket := createNewTicket("PT")

	return &ticket, nil
}

// CreateNewServiceTicket creates a new Service Ticket
// see: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html#service-ticket
func CreateNewServiceTicket(strg storage.IStorage, serviceID string) (*types.Ticket, error) {
	ticket := createNewTicket("ST")
	ticket.Service = serviceID

	strg.SaveTicket(&ticket)

	tools.LogST(&ticket, "ticket created")

	return &ticket, nil
}

// ValidateServiceTicket checks if given Service Ticket exists and is valid
func ValidateServiceTicket(strg storage.IStorage, t *types.Ticket) *types.CasError {
	tools.LogST(t, t.Type+" validation requested")

	ticket := strg.DoesTicketExist(t.Ticket)
	if ticket != nil {
		tools.LogST(ticket, ticket.Type+" validation succeeded (ticket was found)")

		strg.DeleteTicket(ticket.Ticket)

		tools.LogST(ticket, ticket.Type+" deleted")

		return nil
	}

	tools.LogST(t, ticket.Type+" validation failed (ticket was not found)")

	return &types.CasError{Error: fmt.Errorf("The ticket `%s` is invalid.", t.Ticket), CasErrorCode: types.CAS_ERROR_CODE_INVALID_TICKET}
}

// ValidateProxyGrantingTicket checks if given PGT exists
func ValidateProxyGrantingTicket(strg storage.IStorage, t *types.Ticket) *types.CasError {
	tools.LogPGT(t, "PGT validation requested")

	ticket := strg.DoesTicketExist(t.Ticket)
	if ticket != nil {
		tools.LogPGT(ticket, "PGT validation succeeded (ticket was found)")

		// only deleted PGT ticket if its old
		if ticket.Old() {
			strg.DeleteTicket(ticket.Ticket)
			tools.LogPGT(ticket, "PGT deleted")
		}

		return nil
	}

	tools.LogPGT(t, "PGT validation failed (ticket was not found)")

	return &types.CasError{Error: fmt.Errorf("The ticket `%s` is invalid.", t.Ticket), CasErrorCode: types.CAS_ERROR_CODE_INVALID_TICKET}
}
