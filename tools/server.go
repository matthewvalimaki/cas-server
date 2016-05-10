package tools

import (
    "fmt"
    "net/http"
    
    "github.com/matthewvalimaki/cas-server/types"
)

var (
    servers []*types.Server
)

// StartServers starts listening per server configuration
func StartServers(config *types.Config) error {
    for _, server := range config.Servers {
        if server.SSL {                  
            go startServer(server.PortToString(), server.CACert, server.CAKey)
            
            Log(fmt.Sprintf("cas-server is now started up and binding to all interfaces with port `%d` using SSL", server.Port))
        } else {           
            go startServer(server.PortToString(), "", "")
            
            Log(fmt.Sprintf("cas-server is now started up and binding to all interfaces with port `%d`", server.Port))
        }
    }
    
    return nil
}

func startServer(port string, cacert string, cakey string) {
    var err error
    
    if cacert != "" {
        err = http.ListenAndServeTLS(":" + port, cacert, cakey, nil)
    } else {
        err = http.ListenAndServe(":" + port, nil)
    }
    
    if err != nil {
        LogError(err.Error())
        return
    }
}