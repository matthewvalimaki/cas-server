package spec

import (
    "net/http"
    "text/template"
    "encoding/json"

    "github.com/matthewvalimaki/cas-server/types"
    "github.com/matthewvalimaki/cas-server/validators"
    "github.com/matthewvalimaki/cas-server/security"
)

// SupportV2 enables spec v2 support
func SupportV2() {
    validateV2()
    proxyV2()
}

func validateV2() {
    http.HandleFunc("/v2/serviceValidate", setupValidateV2)
}

func proxyV2() {
    http.HandleFunc("/v2/proxy", setupProxyV2)
}

func setupValidateV2(w http.ResponseWriter, r *http.Request) {
    // do format check here as it will affect response format
    format := r.URL.Query().Get("format")
    if len(format) > 0 {
        err := validators.ValidateFormat(format)
        validateResponseV2("XML", err, nil, w, r)
        return
    }
    if len(format) == 0 {
        format = "XML"
    }
    
    _, err := runValidators(w, r)
    if err != nil {
        validateResponseV2(format, err, nil, w, r)
        return
    }
    
    _, pgtURL, proxyGrantingTicket, proxyGrantingTicketIOU, err := runValidatorsV2(w, r)
    if err != nil {
        validateResponseV2(format, err, nil, w, r)
        return
    }
    
    // see: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html#servicevalidate-cas-20
    if pgtURL != "" {
        strg.SaveTicket(proxyGrantingTicket)
        
        validateResponseV2(format, nil, proxyGrantingTicketIOU, w, r)
    }

    validateResponseV2(format, nil, nil, w, r)
}

func runValidatorsV2(w http.ResponseWriter, r *http.Request) (service *types.Service, pgtURL string, proxyGrantingTicket *types.Ticket, proxyGrantingTicketIOU *types.Ticket, err *types.CasError) {   
    pgtURL = r.URL.Query().Get("pgtUrl")
    if len(pgtURL) > 0 {
        serviceParameter := r.URL.Query().Get("service")
        
        // make sure that pgtURL can be used with service 
        service, err := validators.ValidateProxyGrantingURL(config, serviceParameter, pgtURL)
        if err != nil {
            return nil, "", nil, nil, err
        }
        
        // Make sure endpoint can be reached and uses SSL as dictated by CAS spec
        // see: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html#head2.5.4
        err = validators.ValidateProxyURLEndpoint(pgtURL)
        if err != nil {
            return nil, "", nil, nil, err
        }
        
        // Generate PGT (ProxyGrantingTicket) and PGTIOU (ProxyGgrantingTicketIOU) 
        proxyGrantingTicket, err := security.CreateNewProxyGrantingTicket()
        if err != nil {
            return nil, "", nil, nil, err
        }
        
        proxyGrantingTicketIOU, err := security.CreateNewProxyGrantingTicketIOU()
        if err != nil {
            return nil, "", nil, nil, err
        }
        
        // reach out to proxy and then validate behavior
        err = validators.SendAndValidateProxyIDAndIOU(pgtURL, proxyGrantingTicket, proxyGrantingTicketIOU)
        if err != nil {
            return nil, "", nil, nil, err
        }        
        
        return service, pgtURL, proxyGrantingTicket, proxyGrantingTicketIOU, nil
    }
    
    return nil, "", nil, nil, nil
}

func validateResponseV2(format string, casError *types.CasError, proxyGrantingTicketIOU *types.Ticket, w http.ResponseWriter, r *http.Request) {
    if format == "XML" {
        w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
        
        if casError != nil {
            t, _ := template.ParseFiles("spec/tmpl/v2ValidationFailure.tmpl")
            
            t.Execute(w, map[string] string {"Error": casError.Error.Error(), "CasErrorCode": casError.CasErrorCode.String()})
        } else {
            t, _ := template.ParseFiles("spec/tmpl/v2ValidationSuccess.tmpl")
            
            t.Execute(w, map[string] string {"Username": "test", "proxyGrantingTicketIOU": proxyGrantingTicketIOU.Ticket})
        }
    } else {
        // response := new({})
        
        // js, err := json.Marshal({"serviceResponse": {}})
        // if err != nil {
        //     http.Error(w, err.Error(), http.StatusInternalServerError)
        //     return
        // }
  
        w.Header().Set("Content-Type", "application/json;charset=UTF-8")
        w.Write()
    }
}

func setupProxyV2(w http.ResponseWriter, r *http.Request) {
    err := validators.ValidateProxyTicket()
}