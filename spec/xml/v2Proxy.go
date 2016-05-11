package xml

import (
    "fmt"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// V2ProxyFailure produces XML string for failure
func V2ProxyFailure(casError *types.CasError) string {
    return fmt.Sprintf(`<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas"><cas:proxyFailure code="%s">%s</cas:proxyFailure></cas:serviceResponse>`, 
        casError.CasErrorCode.String(), casError.Error.Error())
}

// V2ProxySuccess produces XML string for success
func V2ProxySuccess(proxyTicket *types.Ticket) string {
    return fmt.Sprintf(`<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
    <cas:proxySuccess>
        <cas:proxyTicket>%s</cas:proxyTicket>
    </cas:proxySuccess>
</cas:serviceResponse>`, 
        proxyTicket.Ticket)
}