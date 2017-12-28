package main

import (
    "strings"
    "base/zkcli"
    "base/log"

    "project/share"
)

var (
    conn *zkcli.Conn
)

func initZK(zkaddr string) error {
    log.Info("initZK: %v", zkaddr)
    c, err := zkcli.Connect(zkaddr)
    if err != nil {
        log.Error("initZK err: %v", err.Error())
        return err
    }
    conn = c
    return conn.Write(share.ZKPrefixConfig, []byte{})
}

func finiZK() {
    if conn != nil {
        log.Info("close connection to ZK.")
        conn.Close()
    }
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
