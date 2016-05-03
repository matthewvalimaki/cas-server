package test

import (
    "net/http"
)

// SupportTest adds support for internal managed tests 
func SupportTest() {
    validate()
}

func validate() {
    http.HandleFunc("/test/login-redirect", handleValidate)
}

func handleValidate(w http.ResponseWriter, r *http.Request) {    
    ticket := r.URL.Query().Get("ticket")
    if len(ticket) == 0 {
        http.Error(w, "Query Parameter `ticket` must be defined.", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "https://cas.matthewvalimaki.com/cas/serviceValidate?service=http://127.0.0.1:10000/test/login-redirect&ticket=" + ticket, http.StatusFound)
}