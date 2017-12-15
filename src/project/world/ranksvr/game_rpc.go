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

    Log.Println("1")
    l, e := net.Listen("tcp", fmt.Sprintf(":%v", port))
    if e != nil {
        log.Fatal("Error: listen %d error:", port, e)
    }

    wg.Add(1)
    Log.Println("2")
    go func() {
        Log.Println("0")
        for {
            conn, err := l.Accept()
            if err != nil {
                log.Print("Error: accept rpc connection", err.Error())
                continue
            }
            go rpc.ServeConn(conn)
        }
        wg.Done()
    }()

    //注册rpc地址到zk TODO
    serverID := fmt.Sprintf("%v.%v.%v", clusterID, serverType, index)
    addr := fmt.Sprintf("%v:%v", "xxxx", port)
    _ = serverID
    _= addr

    Log.Println("3")
    Log.Printf("server %v rpc serve %v\n", serverID, addr)
}

