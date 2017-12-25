package core

import (
	"net/rpc"
	"fmt"
	"errors"
	"project/share/proto"
)

var (
	rpcClient   *rpc.Client
)

//TODO confsvr的rpc地址
func InitFetcher(svrAddr string) error {
	client, err := rpc.Dial("tcp", svrAddr)
	if err != nil {
		return fmt.Errorf("err dialing: confsvr<%v>, err %v\n", svrAddr, err)
	}
	rpcClient = client
	return nil
}

func FetchConfFromServer(namespace string, keys []string) ([]*proto.ConfigEntry, error){
	if rpcClient == nil {
		panic("fetcher not init.")
	}
	if namespace == "" || len(keys) == 0 {
		panic("inv input for fetch.")
	}
	fmt.Printf("FetchConfFromServer...\n")
	args := &proto.FetchConfigArg{
		Namespace: namespace,
		Keys: keys,
	}
	var response proto.FetchConfigRes
	err := rpcClient.Call("Conf.FetchConfig", args, &response)
	if err != nil {
		return fmt.Errorf("FetchConfig call error %v\n", err)
	}
	if response.Errmsg != "" {
		return nil, errors.New(response.Errmsg)
	}
	return response.Confs, nil
}

