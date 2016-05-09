package validators

import (
    "fmt"
    "errors"
    "sort"
    "net/http"
    "strings"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// ValidateProxyGrantingURL validates proxy pgtURL
func ValidateProxyGrantingURL(config *types.Config, serviceID string, pgtURL string) (*types.Service, *types.CasError) {
    // make sure we have the pgtURL specified as service in the first place
    proxyServices, ok := config.FlatServiceIDList[pgtURL]
    if !ok {
        return nil, &types.CasError{Error: errors.New("`pgtUrl` was not found from configuration"), CasErrorCode: types.CAS_ERROR_CODE_UNAUTHORIZED_SERVICE_PROXY}
    }
    
    // make sure the service exists
    services, ok := config.FlatServiceIDList[serviceID]
    if !ok {
        return nil, &types.CasError{Error: errors.New("`service` was not found from configuration"), CasErrorCode: types.CAS_ERROR_CODE_INVALID_SERVICE}
    }
    
    for _, service := range services {
        for _, proxyService := range proxyServices {    
            // make sure the pgtURL is found with the service
            sort.Strings(config.Services[service].ProxyServices)
            
            i := sort.SearchStrings(config.Services[service].ProxyServices, proxyService)
            if i < len(config.Services[service].ProxyServices) && config.Services[service].ProxyServices[i] == proxyService {
                return config.Services[service], nil
            }
        }
    }

    return nil, &types.CasError{Error: fmt.Errorf("Proxy service `%s` is not allowed to get proxy ticket with service `%s`", pgtURL, serviceID), CasErrorCode: types.CAS_ERROR_CODE_UNAUTHORIZED_SERVICE_PROXY}
}

// ValidateProxyURLEndpoint reaches out to the proxy URL
// see: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html#head2.5.4
func ValidateProxyURLEndpoint(pgtURL string) *types.CasError {
    _, err := http.Get(pgtURL)
    if err != nil {
        return &types.CasError{Error: fmt.Errorf("Proxy service `%s` validation failed with error: `%s`", pgtURL, err.Error()), CasErrorCode: types.CAS_ERROR_CODE_INVALID_PROXY_CALLBACK}
    }
    
    return nil
}

// SendAndValidateProxyIDAndIOU reaches out to the proxy URL with query parameters
func SendAndValidateProxyIDAndIOU(pgtURL string, proxyGrantingTicket *types.Ticket, proxyGrantingTicketIOU *types.Ticket) *types.CasError {
    pgtURLWithParameters := pgtURL
    
    if strings.Contains(pgtURL, "?") {
        pgtURLWithParameters += "&"
    }
    
    pgtURLWithParameters += "pgtId=" + proxyGrantingTicket.Ticket
    pgtURLWithParameters += "pgtIou=" + proxyGrantingTicketIOU.Ticket
    
    response, err := http.Get(pgtURLWithParameters)
    if err != nil {
        return &types.CasError{Error: fmt.Errorf("Proxy service `%s` validation failed with error: `%s`", pgtURL, err.Error()), CasErrorCode: types.CAS_ERROR_CODE_INVALID_PROXY_CALLBACK}
    }
    
    // enforce required status code check
    if response.StatusCode != http.StatusOK {
        return &types.CasError{Error: fmt.Errorf("Proxy service with CAS query parameters `%s` returned status code `%d` while `%d` is required", pgtURLWithParameters, response.StatusCode, http.StatusOK), CasErrorCode: types.CAS_ERROR_CODE_INVALID_PROXY_CALLBACK}
    }
    
    return nil
}
