package main

import (
    "base/srv"
)

const (
    rpcName     = "Rank"
)

func main() {
    name, worker := NewRpc()
    srv.HandleRpc(name, worker)
    srv.Run()
}

