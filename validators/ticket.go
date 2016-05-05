package validators

import (
    "fmt"
    "errors"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// ValidateTicket validates ticket
func ValidateTicket(ticket string) (*types.CasError) {
    err := validateTicketLength(ticket)
    if (err != nil) {
        return err
    }
    
    err = validateTicketFormat(ticket)
    if (err != nil) {
        return err
    }
    
    err = validateTicketTimestamp(ticket)
    if (err != nil) {
        return err
    }

    return nil
}

func validateTicketLength(ticket string) *types.CasError {
    if len(ticket) == 0 {
        return &types.CasError{Error: errors.New("Required query parameter `ticket` was not defined."), CasErrorCode: types.CAS_ERROR_CODE_INVALID_REQUEST}
    }
    
    if len(ticket) < 32 {
        return &types.CasError{Error: fmt.Errorf("Ticket is not long enough. Minimum length is `%d` but length was `%d`.", 32, len(ticket)), CasErrorCode: types.CAS_ERROR_CODE_INVALID_TICKET_SPEC}
    }
    
    if len(ticket) > 256 {
        return &types.CasError{Error: fmt.Errorf("Ticket is too long. Maximum length is `%d` but length was `%d`.", 256, len(ticket)), CasErrorCode: types.CAS_ERROR_CODE_INVALID_TICKET_SPEC}
    }
    
    return nil
}

func validateTicketFormat(ticket string) *types.CasError {
    if ticket[0:3] != "ST-" {
        return &types.CasError{Error: fmt.Errorf("Required prefix `%s` is missing.", "ST-"), CasErrorCode: types.CAS_ERROR_CODE_INVALID_TICKET_SPEC}
    }
    
    return nil
}

func validateTicketTimestamp(ticket string) *types.CasError {
    return nil
}
