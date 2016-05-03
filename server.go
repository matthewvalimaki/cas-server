package main

import (
    "fmt"
    "net/http"
    
    "github.com/matthewvalimaki/cas-server/admin"
    "github.com/matthewvalimaki/cas-server/tools"
    "github.com/matthewvalimaki/cas-server/spec"
    "github.com/matthewvalimaki/cas-server/test"
    "github.com/matthewvalimaki/cas-server/storage"
)

func main() {
    config, err := tools.NewConfig()
    
    if err != nil {
        tools.Log(err.Error())
    }
    
    admin.SupportServices(config)
    
    storage := storage.NewMemoryStorage()
    
    spec.SupportV1(storage, config)
    spec.SupportV2()
    spec.SupportV3()
    
    test.SupportTest()
    
    tools.Log(fmt.Sprintf("cas-server is starting up and binding to all interfaces with port `%d`", config.Port))
    
    http.ListenAndServe(":" + config.PortToString(), nil)
}