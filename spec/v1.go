package spec

import (
    "fmt"
    "net/http"
    "text/template"
    
    "github.com/matthewvalimaki/cas-server/validators"
    "github.com/matthewvalimaki/cas-server/types"
    "github.com/matthewvalimaki/cas-server/storage"
    "github.com/matthewvalimaki/cas-server/security"
)

type loginResponseFnType func(*types.CasError, *types.Ticket, http.ResponseWriter, *http.Request)
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
    validateResponseFn = validateResponse
}

func login() {
    http.HandleFunc("/v1/login", handleLogin)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    service := r.URL.Query().Get("service")
    
    err := validators.ValidateService(service, config)
    
    if err != nil {
        loginResponseFn(err, nil, w, r)
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
    http.HandleFunc("/v1/validate", handleValidate)
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
    service := r.URL.Query().Get("service")
    ticket := r.URL.Query().Get("ticket")
    
    err := validators.ValidateTicket(ticket)
    
    if err != nil {
        validateResponseFn(false, err, nil, w, r)
        return
    }
    
    if len(service) == 0 {
        http.Error(w, "Query Parameter `service` must be defined.", http.StatusInternalServerError)
        return
    }

    serviceTicket := &types.Ticket{Service: service, Ticket: ticket}
    err = security.ValidateServiceTicket(strg, serviceTicket)

    if err != nil {
        validateResponseFn(false, err, serviceTicket, w, r)
        return
    }
    validateResponseFn(true, nil, serviceTicket, w, r)
}

func validateResponse(valid bool, casError *types.CasError, ticket *types.Ticket, w http.ResponseWriter, r *http.Request) {
    if valid {
        fmt.Fprintf(w, "yes")
    } else {
        fmt.Fprintf(w, "no")
    }
}
