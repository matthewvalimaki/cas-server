package spec

import (
    "fmt"
    "errors"
    "net/http"
    
    "github.com/matthewvalimaki/cas-server/tools"
    "github.com/matthewvalimaki/cas-server/validators"
    "github.com/matthewvalimaki/cas-server/types"
    "github.com/matthewvalimaki/cas-server/storage"
    "github.com/matthewvalimaki/cas-server/security"
)

var (
    specTemplatePath = "spec/tmpl/"
    
    strg storage.IStorage
    config *types.Config
)

// SupportV1 enables spec v1 support
func SupportV1(strgObject storage.IStorage, cfg *types.Config) {
    strg = strgObject
    config = cfg
          
    http.HandleFunc("/login", handleLogin)
    http.HandleFunc("/validate", setupValidate)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    if config == nil {
        err := &types.CasError{Error: errors.New("`config` has not been set"), CasErrorCode: types.CAS_ERROR_CODE_INTERNAL_ERROR}
        loginResponse(err, nil, w, r)
        tools.LogRequest(r, err.Error.Error())
        return
    }
    
    if strg == nil {
        err := &types.CasError{Error: errors.New("`strg` has not been set"), CasErrorCode: types.CAS_ERROR_CODE_INTERNAL_ERROR}
        loginResponse(err, nil, w, r)
        tools.LogRequest(r, err.Error.Error())
        return
    }
    
    err := validators.ValidateRequest(r)
    if err != nil {
        loginResponse(err, nil, w, r)
        tools.LogRequest(r, err.Error.Error())
        return
    }    
    
    service := r.URL.Query().Get("service")  
    err = validators.ValidateService(service, config)
    if err != nil {
        loginResponse(err, nil, w, r)
        tools.LogService(service, err.Error.Error())
        security.ProcessFailedLogin(r.RemoteAddr)
        return
    }
    
    var serviceTicket, _ = security.CreateNewServiceTicket(strg, service)
    
    loginResponse(nil, serviceTicket, w, r)
}

func loginResponse(casError *types.CasError, ticket *types.Ticket, w http.ResponseWriter, r *http.Request) {
    if casError != nil {
        if casError.CasErrorCode == types.CAS_ERROR_CODE_INTERNAL_ERROR {
            w.WriteHeader(http.StatusInternalServerError)
        }
        
        fmt.Fprintf(w, casError.Error.Error())
        return
    }
    
    http.Redirect(w, r, ticket.Service + "?ticket=" + ticket.Ticket, http.StatusFound)
}

func setupValidate(w http.ResponseWriter, r *http.Request) {
    serviceTicket, err := runValidators(w, r)
    
    if err != nil {
        validateResponse(false, err, nil, w, r)
        return
    }
    
    validateResponse(true, nil, serviceTicket, w, r)
}

func runValidators(w http.ResponseWriter, r *http.Request) (*types.Ticket, *types.CasError) {
    ticket := r.URL.Query().Get("ticket")
    
    err := validators.ValidateTicket(ticket)
    if err != nil {
        return nil, err
    }

    service := r.URL.Query().Get("service")    
    serviceError := validators.ValidateService(service, config)
    if serviceError != nil {
        return nil, serviceError
    }
    
    serviceTicket := &types.Ticket{Service: service, Ticket: ticket}
    err = security.ValidateServiceTicket(strg, serviceTicket)
    if err != nil {
        return nil, err
    }
    
    return serviceTicket, nil
}

func validateResponse(valid bool, casError *types.CasError, ticket *types.Ticket, w http.ResponseWriter, r *http.Request) {
    if valid {
        fmt.Fprintf(w, "yes")
    } else {
        fmt.Fprintf(w, "no")
    }
}