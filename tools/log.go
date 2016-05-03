package tools

import (
    "fmt"
    "time"
    
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

// LogST prints ST related log message
func LogST(ticket *types.Ticket, message string) {
    log("ST", fmt.Sprintf("[%s] [%s] %s", ticket.Service, ticket.Ticket, message))
}

// LogAdmin logs admin message
func LogAdmin(message string) {
    log("admin", message)
}