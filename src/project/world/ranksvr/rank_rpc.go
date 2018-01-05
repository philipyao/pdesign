package main

import (
    "base/log"
    "project/share/proto"
)

const (
    RpcName         = "Rank"
)

func NewRpc() (string, interface{}) {
    return RpcName, new(RpcWorker)
}

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
    return nil
}
