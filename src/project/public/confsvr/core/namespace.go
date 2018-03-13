package core

import (
    "errors"
    "base/log"

    "project/public/confsvr/def"
    "project/public/confsvr/db"
)

var (
    ns Namespace = Namespace{}
)

type Namespace []string

func (n Namespace) Load(ns []string) {
    copy(n, ns)
}

func (n Namespace) Exist(val string) bool {
    for _, entry := range n {
        if entry == val {
            return true
        }
    }
    return false
}

//预先生成公共空间
func (n Namespace) CreateCommon() error {
    name := def.ConfNamespaceCommon
    exist, err := db.ExistNamespace(name)
    if err != nil {
        return err
    }
    if exist {
        return nil
    }
    namespace := &def.Namespace {
        Name: name,
        Desc: "公共配置区间，配置项可以被私有同名配置项覆盖",
        Creator: def.AdminUsername,
    }
    log.Info("create namespace [common] ok")
    err = db.InsertNamespace(namespace)
    if err != nil {
        return err
    }
    n = append(n, name)
    return nil
}

//创建普通私有空间
func (n Namespace)Create(creator, name, desc string) error {
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
    namespace := &def.Namespace {
        Name: name,
        Desc: desc,
        Creator: creator,
    }
    err = db.InsertNamespace(namespace)
    if err != nil {
        return err
    }
    log.Info("create namespace [%v] ok", name)
    n = append(n, name)
    return nil
}
