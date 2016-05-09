package spec

import (
    "net/http"
    "testing"
    "strings"
    "net/http/httptest"
    "github.com/matthewvalimaki/cas-server/types"
    "github.com/matthewvalimaki/cas-server/storage"
)

func TestHandleLoginStorageNotSet(t *testing.T) {
    config = &types.Config{}
    specTemplatePath = "tmpl/"
    
    req, _ := http.NewRequest("GET", "?service=http://test.com", nil)
    req.RemoteAddr = "127.0.0.1:1234"
    w := httptest.NewRecorder()
    
    handleLogin(w, req)
    
    if w.Code != http.StatusInternalServerError {
        t.Errorf("Login did not return `%d`", http.StatusInternalServerError)
    }
    
    if !strings.Contains(w.Body.String(), "`strg` has not been set") { 
        t.Errorf("Validation should have returned ``strg` has not been set`.")
    }
}

func TestHandleLoginUnknownService(t *testing.T) {
    config = &types.Config{}
    strg = &storage.MemoryStorage{}
    specTemplatePath = "tmpl/"
    
    req, _ := http.NewRequest("GET", "?service=http://test.com", nil)
    req.RemoteAddr = "127.0.0.1:1234"
    w := httptest.NewRecorder()
    
    handleLogin(w, req)
    
    if !strings.Contains(w.Body.String(), "Unknown service") { 
        t.Errorf("Validation should have returned `Unknown service`.")
    }
}

func TestHandleLoginSuccess(t *testing.T) {
    strg = storage.NewMemoryStorage()
    services := make(map[string]*types.Service)
    services["test"] = &types.Service{ID: []string{"http://test.com"}}
    config = &types.Config{Services: services}
    config.FlattenServiceIDs()
    
    specTemplatePath = "tmpl/"
    
    req, _ := http.NewRequest("GET", "?service=http://test.com", nil)
    req.RemoteAddr = "127.0.0.1:1234"
    w := httptest.NewRecorder()
    
    handleLogin(w, req)
    
    if w.Code != http.StatusFound {
        t.Errorf("Login did not return `%d`", http.StatusFound)
    }

    if !strings.Contains(w.HeaderMap.Get("location"), "ticket=") {
        t.Errorf("Redirect did not contain query parameter `ticket`")
    }
}

func TestValidate(t *testing.T) {
    strg = storage.NewMemoryStorage()
    services := make(map[string]*types.Service)
    services["test"] = &types.Service{ID: []string{"http://test.com"}}
    config = &types.Config{Services: services}
    config.FlattenServiceIDs()
    
    specTemplatePath = "tmpl/"
    
    // without ticket
    req, _ := http.NewRequest("GET", "?service=http://test.com", nil)
    req.RemoteAddr = "127.0.0.1:1234"
    w := httptest.NewRecorder()
    
    setupValidate(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Validate returned `%d` was expecting `%d`", w.Code, http.StatusOK)
    }

    if w.Body.String() != "no" {
        t.Errorf("Response body contained `%s` was expecting `no`", w.Body.String())
    }
    
    // with ticket
    req, _ = http.NewRequest("GET", "?service=http://test.com", nil)
    req.RemoteAddr = "127.0.0.1:1234"
    w = httptest.NewRecorder()
    
    handleLogin(w, req)
    ticket := w.HeaderMap.Get("location")[strings.LastIndex(w.HeaderMap.Get("location"), "ticket=") + 7 : len(w.HeaderMap.Get("location"))]
    
    req, _ = http.NewRequest("GET", "?service=http://test.com&ticket=" + ticket, nil)
    req.RemoteAddr = "127.0.0.1:1234"
    w = httptest.NewRecorder()
    setupValidate(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Validate returned `%d` was expecting `%d`", w.Code, http.StatusOK)
    }
    
    if w.Body.String() != "yes" {
        t.Errorf("Response body contained `%s` was expecting `yes` for ticket `%s`", w.Body.String(), ticket)
    }
}