package main

import (
    "fmt"
    "net/http"
    
    "github.com/matthewvalimaki/cas-server/types"
    "github.com/matthewvalimaki/cas-server/tools"
    "github.com/matthewvalimaki/cas-server/spec"
    "github.com/matthewvalimaki/cas-server/admin"
)

var (
    servers []*types.Server
)

// StartServers starts listening per server configuration
func startServers(config *types.Config) error {
    mux := getCorsMux()
    
    for _, server := range config.Servers {
        if server.SSL {                  
            go startServer(mux, config.Cors, server.PortToString(), server.CACert, server.CAKey)
            
            tools.Log(fmt.Sprintf("cas-server is now started up and binding to all interfaces with port `%d` using SSL", server.Port))
        } else {           
            go startServer(mux, config.Cors, server.PortToString(), "", "")
            
            tools.Log(fmt.Sprintf("cas-server is now started up and binding to all interfaces with port `%d`", server.Port))
        }
    }
    
    return nil
}

func startServer(mux *http.ServeMux, cors *types.Cors, port string, cacert string, cakey string) {
    var err error
    
    if cacert != "" {
        err = http.ListenAndServeTLS(":" + port, cacert, cakey, corsHandler(mux, cors))
    } else {
        err = http.ListenAndServe(":" + port, corsHandler(mux, cors))
    }
    
    if err != nil {
        tools.LogError(err.Error())
        return
    }
}

func getCorsMux() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/admin/services", admin.HandleServices)
    
    // v1
    mux.HandleFunc("/login", spec.HandleLogin)
    mux.HandleFunc("/validate", spec.HandleValidate)
    
    // v2
    mux.HandleFunc("/serviceValidate", spec.HandleValidateV2)
    mux.HandleFunc("/proxyValidate", spec.HandleValidateV2)
    mux.HandleFunc("/proxy", spec.HandleProxyV2)
    
    // v3
    mux.HandleFunc("/p3/serviceValidate", spec.HandleValidateV2)
    mux.HandleFunc("/p3/proxyValidate", spec.HandleValidateV2)
    
    return mux
}

func corsHandler(next http.Handler, cors *types.Cors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if cors != nil {
            if cors.Origin != nil {
                w.Header().Set("Access-Control-Allow-Origin", cors.OriginToString())
            }
            
            if cors.Methods != nil {
                w.Header().Set("Access-Control-Allow-Methods", cors.MethodsToString())
            }
            
            if cors.Credentials != false {
                w.Header().Set("Access-Control-Allow-Credentials", "true")               
            } else {
                w.Header().Set("Access-Control-Allow-Credentials", "false")
            }
        }
        
		// If this was preflight options request let's write empty ok response and return
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			w.Write(nil)
			return
		}
        
        next.ServeHTTP(w, r)
    })
}
