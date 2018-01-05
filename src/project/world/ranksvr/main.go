package main

import (
    "base/srv"
)

const (
    rpcName     = "Rank"
)

func main() {
    name, worker := NewRpc()
    err := srv.HandleRpc(name, worker)
    if err != nil {
        panic(err)
    }
    srv.Run()
}

