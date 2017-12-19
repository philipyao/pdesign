package main

import (
    "fmt"
    "log"
    "time"
    "net"
    "net/rpc"

    "project/share/proto"
)

const (
    RpcName         = "Conf"
)



type RpcWorker int

func (r *RpcWorker) GamesvrHello(args *proto.GameHelloArg,
                                 reply *proto.GameHelloRep) error {
    reply.C = args.A * args.B
    return nil
}

//根据特定namespace获取配置键值对
func (r *RpcWorker) ConfigWithNamespace(args *proto.ConfigWithNamespaceArg,
                                        reply *proto.ConfigWithNamespaceRep) error {
    Log.Printf("ConfigWithNamespace: args %+v\n", args)
    confMap := ConfigWithNamespace(args.Namespace)
    for k, v := range confMap {
        reply.Confs = append(reply.Confs, &proto.Config{
            Key:    k,
            Value:  v,
        })
    }
    return nil
}


//==============================================================================

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
    Log.Printf("server %v rpc serve %v\n", serverID, addr)
}

func doServe(listener *net.TCPListener) {
    defer wg.Done()
    defer listener.Close()

    for {
        select {
        case <-done:
            Log.Printf("stopping rpc listening on %v...", listener.Addr())
            return
        default:
        }
        listener.SetDeadline(time.Now().Add(1e9))
        conn, err := listener.AcceptTCP()
        if err != nil {
            if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                continue
            }
            Log.Printf("Error: accept rpc connection, %v", err.Error())
        }
        //TODO wg.Add(1)
        go rpc.ServeConn(conn)
    }
}