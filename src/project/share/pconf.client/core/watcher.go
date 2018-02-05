package core

import (
    "fmt"
    "strings"

    "base/zkcli"
    "project/share/commdef"
)
var (
    conn    *zkcli.Conn
)

func InitWatcher(zkaddr string) error {
    c, err := zkcli.Connect(zkaddr)
    if err != nil {
        return err
    }
    conn = c
    return nil
}

func WatchEntryUpdate(namespace, key string, notify chan string, done chan struct{}) error {
    if conn == nil {
        return fmt.Errorf("nil zk conn")
    }
    entryPath := strings.Join([]string{commdef.ZKPrefixConfig, namespace, key}, "/")
    return conn.Watch(entryPath, func(p string, d []byte, e error){
        Log("fire watch for local<%v %v>, code %v\n", namespace, key, e)
        if e != nil {
            //todo watch出错
        }
        notify <- key
    }, done)
}
