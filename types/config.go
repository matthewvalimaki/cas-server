package types

import (
    "strconv"
)

// Config contains all of configuration
type Config struct {
    Port int
    Services map[string]Service
    
    FlatServiceIDList map[string][]string
}

// PortToString converts integer port to string
func (c Config) PortToString() string {
    return strconv.Itoa(c.Port)
}

// FlattenServiceIDs takes all service ids and flattens them
func (c *Config) FlattenServiceIDs() {   
    c.FlatServiceIDList = make(map[string][]string)
    
    for key := range c.Services {        
        for _, serviceID := range c.Services[key].ID {
            
            c.FlatServiceIDList[serviceID] = append(c.FlatServiceIDList[serviceID], key)
        }
    }
}