package main

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


///=====================================================================
func NewRpc() (string, interface{}) {
    return RpcName, new(RpcWorker)
}