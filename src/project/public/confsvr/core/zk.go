package core

import (
    "strings"
    "base/zkcli"
    "base/log"

    "project/share/commdef"
)

var (
    conn *zkcli.Conn
)

func initZK(zkaddr string) error {
    if conn != nil {
        panic("duplicated initZK")
    }
    c, err := zkcli.Connect(zkaddr)
    if err != nil {
        log.Error("initZK err: %v", err.Error())
        return err
    }
    conn = c
    log.Info("open connection to ZK: %v", zkaddr)
    return conn.Write(commdef.ZKPrefixConfig, []byte{})
}

func finiZK() {
    if conn != nil {
        log.Info("close connection to ZK.")
        conn.Close()
        conn = nil
    }
}
func attachNamespaceWithZK(namespace string) error {
    configPath := strings.Join([]string{commdef.ZKPrefixConfig, namespace}, "/")
    return conn.Write(configPath, []byte{})
}

func attachWithZK(namespace, key string) error {
    configPath := strings.Join([]string{commdef.ZKPrefixConfig, namespace}, "/")
    exist, err := conn.Exists(configPath)
    if err != nil {
        return err
    }
    if !exist {
        err = conn.Write(configPath, []byte{})
        if err != nil {
            return err
        }
    }
    log.Debug("attach: %v %v", namespace, key)
    configPath = strings.Join([]string{commdef.ZKPrefixConfig, namespace, key}, "/")
    exist, err = conn.Exists(configPath)
    if err != nil {
        return err
    }
    if !exist {
        log.Info("init config path: %v", configPath)
        return conn.Write(configPath, []byte{})
    }
    return nil
}

//pub变更消息给ZK
func notifyWithZK(namespace, key string) error {
    log.Info("notifyWithZK: %v %v", namespace, key)
    configPath := strings.Join([]string{commdef.ZKPrefixConfig, namespace, key}, "/")
    return conn.Write(configPath, []byte{})
}
