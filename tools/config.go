package tools

import (
    "github.com/matthewvalimaki/cas-server/types"
    
    "github.com/BurntSushi/toml"
)

// NewConfig loads configuratio
func NewConfig(location string) (*types.Config, error) {
    var config types.Config
    _, err := toml.DecodeFile(location, &config)
    
    if err != nil {
        return nil, err
    }
    
    config.FlattenServiceIDs()
    
    return &config, nil
}