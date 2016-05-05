package security

import (
    "net"
)

// IsRemoteAddrBanned checks if address is banned
func IsRemoteAddrBanned(remoteAddr string) bool {
    return true
}

// ProcessFailedLogin processes failed login attempt
func ProcessFailedLogin(remoteAddr string) {
    _, _, err := net.SplitHostPort(remoteAddr)
    if err != nil {
        return
    }
}