package core

import (
    "errors"
    "base/log"

    "project/public/confsvr/def"
    "project/public/confsvr/db"
)

//预先生成公共空间
func createNamespaceCommon() error {
    name := def.ConfNamespaceCommon
    exist, err := db.ExistNamespace(name)
    if err != nil {
        return err
    }
    if exist {
        return nil
    }
    ns := &def.Namespace {
        Name: name,
        Desc: "公共配置区间，配置项可以被私有同名配置项覆盖",
        Creator: def.AdminUsername,
    }
    log.Info("create namespace [common] ok")
    return db.InsertNamespace(ns)
}

//创建普通私有空间
func createNamespace(creator, name, desc string) error {
    if name == "" {
        return errors.New("error create namespace: empty name")
    }
    exist, err := db.ExistNamespace(name)
    if err != nil {
        return err
    }
    if exist {
        return errors.New("error create namespace: already exist")
    }
    ns := &def.Namespace {
        Name: name,
        Desc: desc,
        Creator: creator,
    }
    log.Info("create namespace [%v] ok", name)
    return db.InsertNamespace(ns)
}
