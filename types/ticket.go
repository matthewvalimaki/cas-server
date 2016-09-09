package types

import "time"

// Ticket is used for ticket information generalization
type Ticket struct {
	Ticket  string
	Type    string
	Service string
	Created time.Time
}

// Old calculates time from creation and compares to lifetime
// configuration
func (t *Ticket) Old() bool {
	// lifetime calculation
	return int(time.Now().Sub(t.Created).Minutes()) > 5
}
