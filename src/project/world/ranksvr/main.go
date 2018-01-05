package main

import (
    "base/srv"
)

func main() {

    //handle rpc
    name, worker := NewRpc()
    err := srv.HandleRpc(name, worker)
    if err != nil {
        panic(err)
    }

    //start run
    srv.Run()
}

