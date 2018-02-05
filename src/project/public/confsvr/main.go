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

    //服务器基础：启动，关闭
    err = srv.Handlebase(onInit, onShutdown)
    if err != nil {
        log.Fatal("srv.Handlebase() err: %v", err)
    }

    //进程间通信：RPC服务
    name, worker := NewRpc()
    err = srv.HandleRpc(name, worker)
    if err != nil {
        log.Fatal("srv.HandleRpc() err: %v", err)
    }

    //对外HTTP服务
    err = srv.HandleHttp(":8999", httpHandler)
    if err != nil {
        log.Error("srv.HandleHttp() err: %v", err.Error())
        os.Exit(1)
    }

    srv.Run()
}
