package main

import (
    "fmt"
    "time"
    "net"
    "net/rpc"

    "base/log"
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
func (r *RpcWorker) FetchConfig(args *proto.FetchConfigArg,
                                response *proto.FetchConfigRes) error {
    log.Debug("[rpc]FetchConfig: args %+v", args)
    confMap, err := ConfigWithNamespaceKey(args.Namespace, args.Keys)
    if err != nil {
        response.Errmsg = err.Error()
        return nil
    }
    for k, v := range confMap {
        response.Confs = append(response.Confs, &proto.ConfigEntry{
            Namespace: v[0],
            Key:    k,
            Value:  v[1],
        })
    }
    return nil
}


//==============================================================================

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
        log.Error("Error: listen on %v, %v", laddr, e)
        return e
    }

    wg.Add(1)
    go doServe(l)

    //注册rpc地址到zk TODO
    serverID := fmt.Sprintf("%v.%v.%v", clusterID, serverType, index)
    log.Info("server %v rpc serve %v", serverID, addr)

    return nil
}

func doServe(listener *net.TCPListener) {
    defer wg.Done()
    defer listener.Close()

    for {
        select {
        case <-done:
            log.Info("stop rpc listening on %v...", listener.Addr())
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
        go rpc.ServeConn(conn)
    }
}