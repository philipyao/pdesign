package core

import (
    "net"
	"net/rpc"
	"fmt"
    "time"
	"errors"
	"project/share/proto"
)

var (
	rpcClient   *rpc.Client
)

//TODO confsvr的rpc地址
func InitFetcher(svrAddr string) error {
    conn, err := net.DialTimeout("tcp", svrAddr, time.Second)
    if err != nil {
        return err
    }
    rpcClient = rpc.NewClient(conn)
	return nil
}

func FetchConfFromServer(namespace string, keys []string) ([]*proto.ConfigEntry, error){
	if rpcClient == nil {
		panic("fetcher not init.")
	}
	if namespace == "" || len(keys) == 0 {
		panic("inv input for fetch.")
	}
	args := &proto.FetchConfigArg{
		Namespace: namespace,
		Keys: keys,
	}
	var response proto.FetchConfigRes
	err := rpcClient.Call("Conf.FetchConfig", args, &response)
	if err != nil {
		return nil, fmt.Errorf("FetchConfig call error %v\n", err)
	}
	if response.Errmsg != "" {
		return nil, errors.New(response.Errmsg)
	}
	return response.Confs, nil
}

