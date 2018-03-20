package main

import (
    log "github.com/philipyao/toolbox/logging"
    "base/srv"

    "project/share"
    "project/share/commdef"

    "project/public/confsvr/core"
)

var (
    serverType  int
)

func onInit(done chan struct{}) error {
    serverType  = commdef.ServerTypeConfsvr

    err := core.Init()
    if err != nil {
        return err
    }

    return nil
}

func onShutdown() {
    core.Fini()

    log.Info("confsvr onShutdown ok.")
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

    //对外HTTP服务
    httpWorker, err := srv.HandleHttp(":8999")
    if err != nil {
        log.Fatal("srv.HandleHttp() err: %v", err.Error())
    }
    err = serveHttp(httpWorker)
    if err != nil {
        log.Fatal("serveHttp err: %v", err.Error())
    }

    srv.Run()
}
