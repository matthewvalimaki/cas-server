package types

import (
    "time"
)

// Ticket is used for ticket information generalization
type Ticket struct {
    Ticket string
    Service string
    CreatedTimestamp time.Time
}