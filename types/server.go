package types

import (
    "strconv"
)

// Server contains server specific configuration
type Server struct {
    Port int
    SSL bool
    CACert string
    CAKey string
}

// PortToString converts integer port to string
func (s Server) PortToString() string {
    return strconv.Itoa(s.Port)
}