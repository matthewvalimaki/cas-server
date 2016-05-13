package xml

import (
    "fmt"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// V2ValidationFailure produces XML string for failure
func V2ValidationFailure(casError *types.CasError, format string) string {
    if format == "XML" {
        return fmt.Sprintf(`<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas"><cas:authenticationFailure code="%s">%s</cas:authenticationFailure></cas:serviceResponse>`, 
            casError.CasErrorCode.String(), casError.Error.Error())
    }
    
    return fmt.Sprintf(`{
        "serviceResponse": {
            "authenticationFailure": {
                "code": "%s",
                "description": "%s"
            }
        }
    }`, 
        casError.CasErrorCode.String(), casError.Error.Error())
}

// V2ValidationSuccess produces XML string for success
func V2ValidationSuccess(username string, proxyGrantingTicket *types.Ticket, format string) string {
    ticket := ""
    
    if format == "XML" {   
        if proxyGrantingTicket != nil {
            ticket = fmt.Sprintf("<cas:proxyGrantingTicket>%s</cas:proxyGrantingTicket>", proxyGrantingTicket.Ticket) 
        }
        
        return fmt.Sprintf(`<cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>
        <cas:authenticationSuccess>%s
            <cas:user>%s</cas:user>
            <cas:attributes>
                <cas:id>2</cas:id>
            </cas:attributes>
        </cas:authenticationSuccess>
    </cas:serviceResponse>`, 
        ticket, username)
    }
    
    if proxyGrantingTicket != nil {
        ticket = fmt.Sprintf(`"proxyGrantingTicket": "%s",`, proxyGrantingTicket.Ticket) 
    }
    
    return fmt.Sprintf(`{
        "serviceResponse": {
            "authenticationSuccess": {%s
                "user": "%s"
            }
        }
    }`, ticket, username)
}