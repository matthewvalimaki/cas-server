package types

// User holds user related data
type User struct {
    FailedLoginCount int
    IP []string
}

// NewUser creates new user
func NewUser(ip string) User {
    var user = User{FailedLoginCount: 0}
    user.IP = append(user.IP, ip)
    
    return user
}