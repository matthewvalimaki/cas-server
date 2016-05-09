package tools

import (
    "fmt"
    "time"
    "net"
    "net/http"
    
    "github.com/matthewvalimaki/cas-server/types"
)

func timestamp() string {
    var timeNow = time.Now()
    
    return timeNow.Format(time.RFC3339)
}

func log(messageType string, message string) {
    fmt.Println(fmt.Sprintf("[%s] [%s] %s", timestamp(), messageType, message))
}

// Log prints generic log message
func Log(message string) {   
    log("generic", message)
}

func LogError(message string) {
    log("error", message)
}

// LogST prints ST related log message
func LogST(ticket *types.Ticket, message string) {
    log("ST", fmt.Sprintf("[%s] [%s] %s", ticket.Service, ticket.Ticket, message))
}

// LogPGT prints PGT related log message
func LogPGT(ticket *types.Ticket, message string) {
    log("LogPGT", fmt.Sprintf("[%s] [%s] %s", ticket.Service, ticket.Ticket, message))
}

// LogService logs service related message
func LogService(ID string, message string) {
    log("service", fmt.Sprintf("[%s] %s", ID, message))
}

// LogAdmin logs admin message
func LogAdmin(message string) {
    log("admin", message)
}

// LogRequest logs plain request
func LogRequest(r *http.Request, message string) {
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    
    log("request", fmt.Sprintf("[%s] %s", ip, message))
}