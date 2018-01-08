package main

import (
    "fmt"
    "os"
    "path/filepath"

    "base/log"
    "base/srv"

    "project/share"

)

var (
    serverType  int
)

func onInit(done chan struct{}) error {
    serverType  = share.ServerTypeGamesvr
    setLog()
    var err error
    err = InitConf()
    if err != nil {
        log.Error("InitConf: %v", err)
        return err
    }
    err = LoadConf(done)
    if err != nil {
        log.Error("LoadConf: %v", err)
        return err
    }

    return nil
}

func setLog() {
    config := `{"filename": "%v", "maxsize": 102400, "maxbackup": 10}`
    wd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    logName := filepath.Join(wd, "log", srv.ProcessName())
    config = fmt.Sprintf(config, logName)
    fmt.Printf("log config: %+v\n", config)
    err = log.AddAdapter(log.AdapterFile, config)
    if err != nil {
        panic(err)
    }
    log.SetLevel(log.LevelStringDebug)
    log.SetFlags(log.LogDate | log.LogTime | log.LogMicroTime | log.LogLongFile)
}

func onShutdown() {
    log.Info("onShutdown gamesvr")
    log.Flush()
}

func main() {
    var err error
    err = srv.Handlebase(onInit, onShutdown, log.Info)
    if err != nil {
        log.Error("srv.Handlebase() err: %v", err.Error())
        os.Exit(1)
    }

    name, worker := NewRpc()
    err = srv.HandleRpc(name, worker)
    if err != nil {
        log.Error("srv.HandleRpc() err: %v", err.Error())
        os.Exit(1)
    }
    srv.Run()
}



