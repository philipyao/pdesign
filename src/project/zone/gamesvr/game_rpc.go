package main

import (
    "fmt"
    "net"
    "time"
    "net/rpc"

    "base/log"
)

const (
    RpcName         = "Game"
)

type GameHelloArg struct {
    A, B int
}
type GameHelloRep struct {
    C int
}

type RpcWorker int

func (r *RpcWorker) GamesvrHello(args *GameHelloArg, reply *GameHelloRep) error {
    reply.C = args.A * args.B
    return nil
}


///=====================================================================

func serveRPC(done chan struct{}, port int, clusterID, index int) error {
    rpc.RegisterName(RpcName, new(RpcWorker))

    addr := fmt.Sprintf("%v:%v", *ptrIP, port)
    laddr, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil {
        log.Error(err.Error())
        return err
    }

    l, e := net.ListenTCP("tcp", laddr)
    if e != nil {
        log.Error("Error: listen on ", laddr, e)
        return e
    }

    wg.Add(1)
    go doServe(l)

    //注册rpc地址到zk TODO
    serverID := fmt.Sprintf("%v.%v.%v", clusterID, serverType, index)
    log.Info("server %v rpc serve %v\n", serverID, addr)
    return nil
}

func doServe(listener *net.TCPListener) {
    defer wg.Done()
    defer listener.Close()

    for {
        select {
        case <-done:
            log.Info("stopping rpc listening on %v...", listener.Addr())
            return
        default:
        }
        listener.SetDeadline(time.Now().Add(1e9))
        conn, err := listener.AcceptTCP()
        if err != nil {
            if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                continue
            }
            log.Error("Error: accept rpc connection, %v", err.Error())
        }
        //TODO wg.Add(1)
        go rpc.ServeConn(conn)
    }
}

