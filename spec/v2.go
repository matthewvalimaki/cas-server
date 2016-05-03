package spec

import (
    "net/http"
    "text/template"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// SupportV2 enables spec v2 support
func SupportV2() {
    validateV2()
}

func validateV2() {
    http.HandleFunc("/v2/validateService", handleValidate)
    
    validateResponseFn = validateResponseV2
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