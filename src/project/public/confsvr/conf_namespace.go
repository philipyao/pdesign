package main

import (
    "errors"
    "base/log"
)

//预先生成公共空间
func createNamespaceCommon() error {
    name := ConfNamespaceCommon
    exist, err := existNamespace(name)
    if err != nil {
        return err
    }
    if exist {
        return nil
    }
    ns := &Namespace {
        Name: name,
        Desc: "公共配置区间，可以被私有同名配置覆盖",
        Creator: AdminUsername,
    }
    log.Info("create namespace [common] ok")
    return insertNamespace(ns)
}

//创建普通私有空间
func createNamespace(sess *SessionStore, name, desc string) error {
    if name == "" {
        return errors.New("error create namespace: empty name")
    }
    exist, err := existNamespace(name)
    if err != nil {
        return err
    }
    if exist {
        return errors.New("error create namespace: already exist")
    }
    ns := &Namespace {
        Name: name,
        Desc: desc,
        Creator: sess.Get(KeyUserName).(string),
    }
    log.Info("create namespace [%v] ok", name)
    return insertNamespace(ns)
}
