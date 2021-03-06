package main

import (
    "fmt"
    "log"
    "net"
    "net/rpc"
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

func serveRPC(port int, clusterID, index int) {
    rpc.RegisterName(RpcName, new(RpcWorker))

    l, e := net.Listen("tcp", fmt.Sprintf(":%v", port))
    if e != nil {
        log.Fatalf("Error: listen %d error: %v", port, e)
    }

    go func() {
        for {
            conn, err := l.Accept()
            if err != nil {
                log.Print("Error: accept rpc connection", err.Error())
                continue
            }
            go rpc.ServeConn(conn)
        }
    }()

    //注册rpc地址到zk TODO
    serverID := fmt.Sprintf("%v.%v.%v", clusterID, serverType, index)
    addr := fmt.Sprintf("%v:%v", "xxxx", port)
    _ = serverID
    _= addr

    Log.Printf("server %v rpc serve %v\n", serverID, addr)
}

