package main

import (
    "os"

    "base/log"
    "base/srv"

    "project/share"
    "project/share/commdef"
)

var (
    serverType  int
)

func onInit(done chan struct{}) error {
    serverType  = commdef.ServerTypeConfsvr

    //通用添加log支持
    logSize := 102400
    share.SetServerLog(logSize)

    err := initCore()
    if err != nil {
        log.Error(err.Error())
        return err
    }

    return nil
}

func onShutdown() {
    log.Info("onShutdown confsvr.")
    finiCore()
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

    err = srv.HandleHttp(":8999", httpHandler)
    if err != nil {
        log.Error("srv.HandleHttp() err: %v", err.Error())
        os.Exit(1)
    }

    srv.Run()
}
