package main

import (
    "fmt"
    "net/rpc"

    "base/log"
    "project/share/proto"
)

func TryGetGamesvrConfig() {
    client, err := rpc.Dial("tcp", "10.1.164.45:12001")
    if err != nil {
        fmt.Printf("dialing:", err)
        return
    }
    fmt.Printf("TryGetGamesvrConfig...\n")
    args := &proto.ConfigWithNamespaceArg{
        Namespace: "gamesvr",
    }
    var reply proto.ConfigWithNamespaceRep
    err = client.Call("Conf.ConfigWithNamespace", args, &reply)
    if err != nil {
        fmt.Printf("TryGetGamesvrConfig call error %v\n", err)
        return
    }
    for i, c := range reply.Confs {
        fmt.Printf("TryGetGamesvrConfig reply, conf[%v]: %+v\n", i, c)
        log.Debug("TryGetGamesvrConfig reply, conf[%v]: %+v", i, c)
    }

}
