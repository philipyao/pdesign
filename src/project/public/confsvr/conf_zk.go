package main

import (
    "strings"
    "base/zkcli"

    "project/share"
)

var (
    conn *zkcli.Conn
)

func initZK(zkaddr string) error {
    Log.Printf("initZK: %v", zkaddr)
    c, err := zkcli.Connect(zkaddr)
    if err != nil {
        return err
    }
    conn = c
    return conn.Write(share.ZKPrefixConfig, []byte{})
}

func attachNamespaceWithZK(namespace string) error {
    configPath := strings.Join([]string{share.ZKPrefixConfig, namespace}, "/")
    return conn.Write(configPath, []byte{})
}

func attachWithZK(namespace, key string) error {
    configPath := strings.Join([]string{share.ZKPrefixConfig, namespace, key}, "/")
    return conn.Write(configPath, []byte{})
}

func notifyWithZK(namespace, key string) error {
    configPath := strings.Join([]string{share.ZKPrefixConfig, namespace, key}, "/")
    return conn.Write(configPath, []byte{})
}
