package main

import (
    "fmt"
    "net"
    "time"
    "net/rpc"

    "log"
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

func serveRPC(done chan struct{}, port int, clusterID, index int) {
    rpc.RegisterName(RpcName, new(RpcWorker))

    addr := fmt.Sprintf("%v:%v", *ptrIP, port)
    laddr, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil {
        log.Fatalln(err)
    }

    l, e := net.ListenTCP("tcp", laddr)
    if e != nil {
        log.Fatal("Error: listen on ", laddr, e)
    }

    wg.Add(1)
    go doServe(l)

    //注册rpc地址到zk TODO
    serverID := fmt.Sprintf("%v.%v.%v", clusterID, serverType, index)
    fmt.Printf("server %v rpc serve %v\n", serverID, addr)
}

func doServe(listener *net.TCPListener) {
    defer wg.Done()
    defer listener.Close()

    for {
        select {
        case <-done:
            fmt.Printf("stopping rpc listening on %v...\n", listener.Addr())
            return
        default:
        }
        listener.SetDeadline(time.Now().Add(1e9))
        conn, err := listener.AcceptTCP()
        if err != nil {
            if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                continue
            }
            log.Printf("Error: accept rpc connection, %v\n", err.Error())
        }
        //TODO wg.Add(1)
        go rpc.ServeConn(conn)
    }
}

