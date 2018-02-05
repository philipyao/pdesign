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


    var err error
    //加载服务器配置
    err = LoadConf(done)
    if err != nil {
        return err
    }

    return nil
}

func onShutdown() {
    log.Info("gamesvr onShutdown ok.")
    log.Flush()
}

func init () {
    //初始化业务日志
    logSize := 102400
    share.InitBizLog(logSize)

    //自定义srv框架log
    share.SetSrvLogger()
}

func main() {
    var err error

    //服务器基础：启动，关闭
    err = srv.HandleBase(onInit, onShutdown)
    if err != nil {
        log.Fatal("srv.HandleBase() err: %v", err)
    }

    //进程间通信：RPC服务
    name, worker := NewRpc()
    err = srv.HandleRpc(name, worker)
    if err != nil {
        log.Fatal("srv.HandleRpc() err: %v", err)
    }

    srv.Run()
}



