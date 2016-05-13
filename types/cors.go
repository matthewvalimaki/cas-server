package types

// Cors contains CORS configuration
type Cors struct {
    Origin []string
    Methods []string
    Credentials bool
}

// OriginToString converts origin configuration to string
func (c Cors) OriginToString() string {
    var values string
    for _, value := range c.Origin {
        if values != "" {
            values += ","
        }
        values += value
    }
    
    return values
}

// MethodsToString converts methods configuration to string
func (c Cors) MethodsToString() string {
    var values string
    for _, value := range c.Methods {
        if values != "" {
            values += ","
        }
        values += value
    }
    
    return values
}