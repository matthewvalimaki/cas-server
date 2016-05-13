package types

import (
    "sort"
)

// Config contains all of configuration
type Config struct {
    Servers map[string]*Server
    Services map[string]*Service
    Cors *Cors
    
    FlatServiceIDList map[string][]string
}

// FlattenServiceIDs takes all service ids and flattens them
func (c *Config) FlattenServiceIDs() {   
    c.FlatServiceIDList = make(map[string][]string)
    
    for key := range c.Services {        
        for _, serviceID := range c.Services[key].ID {
            c.FlatServiceIDList[serviceID] = append(c.FlatServiceIDList[serviceID], key)
        }
        
        // sorting improves sort.SearchStrings performance
        // so do it once here
        sort.Strings(c.Services[key].ProxyServices)
    }
}