package validators

import (
    "net"
    "net/http"
    "errors"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// ValidateRequest executes validation against plain request
func ValidateRequest(r *http.Request) *types.CasError {
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return &types.CasError{Error: errors.New("Could not parse remote IP:Port."), CasErrorCode: types.CAS_ERROR_CODE_INTERNAL_ERROR}
    }
    
    casError := isRemoteAddrAllowed(ip)
    if casError != nil {
        return casError
    }
    
    return nil
}

func isRemoteAddrAllowed(ip string) *types.CasError {
    return nil
    //return &types.CasError{Error: errors.New("The IP is currently not allowed."), CasErrorCode: types.CAS_ERROR_CODE_INTERNAL_ERROR}
}