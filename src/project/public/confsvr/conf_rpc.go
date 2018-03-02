package main

import (
    "base/log"
    "project/share/proto"

    "project/public/confsvr/core"
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
    confMap, err := core.ConfigWithNamespaceKey(args.Namespace, args.Keys)
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
func NewRpc() (string, interface{}) {
    return RpcName, new(RpcWorker)
}