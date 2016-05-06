package spec

import (
    "net/http"
    "text/template"

    "github.com/matthewvalimaki/cas-server/types"
    "github.com/matthewvalimaki/cas-server/validators"
    "github.com/matthewvalimaki/cas-server/security"
)

// SupportV2 enables spec v2 support
func SupportV2() {
    validateV2()
}

func validateV2() {
    http.HandleFunc("/v2/serviceValidate", setupValidateV2)
}

func setupValidateV2(w http.ResponseWriter, r *http.Request) {
    serviceTicket, err := runValidators(w, r)
    if err != nil {
        validateResponseV2(err, nil, nil, w, r)
        return
    }
    
    service, pgtURL, err := runValidatorsV2(w, r)
    if err != nil {
        validateResponseV2(err, nil, nil, w, r)
        return
    }
    
    if pgtURL != "" {
        proxyGrantingTicket, err := security.CreateNewProxyGrantingTicket(strg, service, pgtURL)
        if err != nil {
            validateResponseV2(err, nil, nil, w, r)
            return
        }
        
        validateResponseV2(nil, serviceTicket, proxyGrantingTicket, w, r)
    }

    validateResponseV2(nil, serviceTicket, nil, w, r)
}

func runValidatorsV2(w http.ResponseWriter, r *http.Request) (*types.Service, string, *types.CasError) {
    pgtURL := r.URL.Query().Get("pgtUrl")

    if len(pgtURL) > 0 {
        serviceParameter := r.URL.Query().Get("service")
        
        // make sure that pgtURL can be used with service 
        service, err := validators.ValidateProxyGrantingURL(config, serviceParameter, pgtURL)
        if err != nil {
            return nil, "", err
        }
        
        err = validators.ValidateProxyURLEndpoint(pgtURL)
        if err != nil {
            return nil, "", err
        }
        
        return service, pgtURL, nil
    }
    
    return nil, "", nil
}

func validateResponseV2(casError *types.CasError, ticket *types.Ticket, proxyGrantingTicket *types.Ticket, w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
    
    if casError != nil {
        t, _ := template.ParseFiles("spec/tmpl/v2ValidationFailure.tmpl")
        
        t.Execute(w, map[string] string {"Error": casError.Error.Error(), "CasErrorCode": casError.CasErrorCode.String()})
    } else {
        t, _ := template.ParseFiles("spec/tmpl/v2ValidationSuccess.tmpl")
        
        t.Execute(w, map[string] string {"Username": "test", "ProxyGrantingTicket": proxyGrantingTicket.Ticket})
    }
}