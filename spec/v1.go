package spec

import (
    "fmt"
    "net/http"
    "text/template"
    
    "github.com/matthewvalimaki/cas-server/tools"
    "github.com/matthewvalimaki/cas-server/validators"
    "github.com/matthewvalimaki/cas-server/types"
    "github.com/matthewvalimaki/cas-server/storage"
    "github.com/matthewvalimaki/cas-server/security"
)

type loginResponseFnType func(*types.CasError, *types.Ticket, http.ResponseWriter, *http.Request)
type validateValidatorFnType func(ID string, ticket string)
type validateResponseFnType func(bool, *types.CasError, *types.Ticket, http.ResponseWriter, *http.Request)

var (
    loginResponseFn loginResponseFnType
    validateResponseFn validateResponseFnType
    
    strg storage.IStorage
    config *types.Config
)

// SupportV1 enables spec v1 support
func SupportV1(strgObject storage.IStorage, cfg *types.Config) {
    strg = strgObject
    config = cfg
          
    login()
    validate()
    
    loginResponseFn = loginResponse
}

func login() {
    http.HandleFunc("/v1/login", handleLogin)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    err := validators.ValidateRequest(r)
    if err != nil {
        loginResponseFn(err, nil, w, r)
        tools.LogRequest(r, err.Error.Error())
        return
    }    
    
    service := r.URL.Query().Get("service")
    
    err = validators.ValidateService(service, config)
    
    if err != nil {
        loginResponseFn(err, nil, w, r)
        tools.LogService(service, err.Error.Error())
        security.ProcessFailedLogin(r.RemoteAddr)
        return
    }
    
    var serviceTicket, _ = security.CreateNewServiceTicket(strg, service)
    
    loginResponseFn(nil, serviceTicket, w, r)
}

func loginResponse(casError *types.CasError, ticket *types.Ticket, w http.ResponseWriter, r *http.Request) {
    if casError != nil {
        w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
        t, _ := template.ParseFiles("spec/tmpl/v2ValidationFailure.tmpl")
        
        t.Execute(w, map[string] string {"Error": casError.Error.Error(), "CasErrorCode": casError.CasErrorCode.String()})
        return
    }
    
    http.Redirect(w, r, ticket.Service + "?ticket=" + ticket.Ticket, http.StatusFound)
}

func validate() {
    http.HandleFunc("/v1/validate", setupValidate)
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
