package types

import (
    "strconv"
)

// Config contains all of configuration
type Config struct {
    Port int
    Services map[string]Service
    FlatServiceIDList []string
}

// PortToString converts integer port to string
func (c Config) PortToString() string {
    return strconv.Itoa(c.Port)
}

// FlattenServiceIDs takes all service ids and flattens them
func (c *Config) FlattenServiceIDs() {
    for key := range c.Services {
        c.FlatServiceIDList = append(c.FlatServiceIDList, c.Services[key].ID...)
    }
}