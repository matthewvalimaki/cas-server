package spec

import (
    "net/http"
    "text/template"

    "github.com/matthewvalimaki/cas-server/types"
    "github.com/matthewvalimaki/cas-server/validators"
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
        validateResponseV2(false, err, nil, w, r)
        return
    }
    
    err = runValidatorsV2(w, r)
    if err != nil {
        validateResponseV2(false, err, nil, w, r)
        return
    }

    validateResponseV2(true, nil, serviceTicket, w, r)
}

func runValidatorsV2(w http.ResponseWriter, r *http.Request) *types.CasError {
    pgtURL := r.URL.Query().Get("pgtUrl")

    if len(pgtURL) > 0 {
        service := r.URL.Query().Get("service")
        
        // make sure that pgtURL can be used with service 
        err := validators.ValidateProxyGrantingURL(config, service, pgtURL)
        if err != nil {
            return err
        }
        
        err = validators.ValidateProxyURLEndpoint(pgtURL)
        if err != nil {
            return err
        }
    }
    
    return nil
}

func validateResponseV2(valid bool, casError *types.CasError, ticket *types.Ticket, w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
    
    if valid {
        t, _ := template.ParseFiles("spec/tmpl/v2ValidationSuccess.tmpl")
        
        t.Execute(w, map[string] string {"Username": "test"})
    } else {
        t, _ := template.ParseFiles("spec/tmpl/v2ValidationFailure.tmpl")
        
        t.Execute(w, map[string] string {"Error": casError.Error.Error(), "CasErrorCode": casError.CasErrorCode.String()})
    }
}