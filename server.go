package main

import (
    "flag"
    "os"
    "os/signal"
    "syscall"
    "time"
    "runtime"
    
    "github.com/matthewvalimaki/cas-server/admin"
    "github.com/matthewvalimaki/cas-server/tools"
    "github.com/matthewvalimaki/cas-server/spec"
    "github.com/matthewvalimaki/cas-server/test"
    "github.com/matthewvalimaki/cas-server/storage"
)

func main() {
    var configPath string
    // get configuration location
    flag.StringVar(&configPath, "config", "", "Path to config file")
    flag.Parse()
    
    if configPath == "" {
        tools.Log("Command line argument `-config` must be set")
        return
    }
    
    config, err := tools.NewConfig(configPath)
    
    if err != nil {
        tools.LogError(err.Error())
        return
    }
    
    admin.SupportServices(config)
    
    storage := storage.NewMemoryStorage()
    
    spec.SupportV1(storage, config)
    spec.SupportV2()
    spec.SupportV3()
    
    test.SupportTest()

    tools.StartServers(config)
    
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    signal.Notify(c, syscall.SIGTERM)
    go func() {
        <-c
        os.Exit(1)
    }()
    
    for {
        runtime.Gosched()
        time.Sleep(5 * time.Second)
    }
}