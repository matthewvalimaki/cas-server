package tools

import (
    "fmt"
    "net/http"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// StartServers starts listening per server configuration
func StartServers(config *types.Config) error {
    for _, server := range config.Servers {
        if server.SSL {           
            err := http.ListenAndServeTLS(":" + server.PortToString(), server.CACert, server.CAKey, nil)
            if err != nil {
                LogError(err.Error())
                return err
            }   
            
            Log(fmt.Sprintf("cas-server is now started up and binding to all interfaces with port `%d` using SSL", server.Port))
        } else {           
            err := http.ListenAndServe(":" + server.PortToString(), nil)
            if err != nil {
                LogError(err.Error())
                return err
            }   
            
            Log(fmt.Sprintf("cas-server is now started up and binding to all interfaces with port `%d`", server.Port))
        }
    }
    
    return nil
}