package spec

import (
    "net/http"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// SupportV3 enables spec v3 support
func SupportV3() {
    loginV3()
}

func loginV3() {
    http.HandleFunc("/v3/login", handleLoginV3)
}

func handleLoginV3(w http.ResponseWriter, r *http.Request) {    
    method := r.URL.Query().Get("method")
    
    if len(method) != 0 {
        if method != "POST" {
            http.Error(w, "Query Parameter `method` must have value `POST` or not present at all.", http.StatusInternalServerError)
            return
        }
        
        if method == "POST" {
            // loginResponseFn = loginResponseV3
        }
    }
    
    handleLogin(w, r)
}

func loginResponseV3(casError *types.CasError, ticket *types.Ticket, w http.ResponseWriter, r *http.Request) {
    r.Header.Set("User-Agent", "cas-server")
    
    http.Post(ticket.Service + "?ticket=" + ticket.Ticket , "text/plain", nil)
}