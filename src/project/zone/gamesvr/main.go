package main

import (
    "base/log"
    "base/srv"

    "project/share"
    "project/share/commdef"
)

var (
    serverType  int
)

func onInit(done chan struct{}) error {
    serverType  = commdef.ServerTypeGamesvr

    //通用添加log支持
    logSize := 102400
    share.SetServerLog(logSize)

    var err error
    //加载服务器配置
    err = LoadConf(done)
    if err != nil {
        return err
    }

    return nil
}

func onShutdown() {
    log.Info("onShutdown gamesvr")
    log.Flush()
}

func main() {
    var err error

    //服务器基础：启动，关闭
    err = srv.Handlebase(onInit, onShutdown, log.Info)
    if err != nil {
        log.Fatal("srv.Handlebase() err: %v", err)
    }

    //进程间通信：RPC服务
    name, worker := NewRpc()
    err = srv.HandleRpc(name, worker)
    if err != nil {
        log.Fatal("srv.HandleRpc() err: %v", err)
    }

    srv.Run()
}



