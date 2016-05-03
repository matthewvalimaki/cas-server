package tools

import (
    "github.com/matthewvalimaki/cas-server/types"
    
    "github.com/BurntSushi/toml"
)

// NewConfig loads configuratio
func NewConfig() (*types.Config, error) {
    var config types.Config
    _, err := toml.DecodeFile("config/default.toml", &config)
    
    if err != nil {
        return &config, err
    }
    
    config.FlattenServiceIDs()
    
    return &config, nil
}


