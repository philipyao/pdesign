package main

import (
    "fmt"
    "time"
    "errors"

    "base/log"
    "project/share"
)

const (
    ConfNamespaceOms    = "oms"
    ConfNamespaceCommon = "common"
    TableNameConf       = "tbl_conf"
)
//配置表
type Config struct {
    ID              uint        `xorm:"pk autoincr 'id'"`
    Namespace       string      `xorm:"varchar(32) notnull"`
    Key             string      `xorm:"varchar(64) notnull"`
    Value           string      `xorm:"varchar(128) notnull"`

    UpdatedAt       time.Time   `xorm:"updated"`
    CreatedAt       time.Time   `xorm:"created"`
    Version         int         `xorm:"version"`    //自动更新版本号
}
func (c Config) TableName() string {
    return TableNameConf
}

var (
    namespaces []string
    confs []*Config
)

func initCore() error {
    var err error
    err = initDB(new(Config))
    if err != nil {
        log.Error(err.Error())
        return err
    }

    log.Debug("loadConfigFromDB")
    confs, namespaces, err = loadConfigFromDB()
    if err != nil {
        log.Error(err.Error())
        return err
    }
    log.Debug("confs: %+v", confs)

    var zkaddr string
    for _, c := range confs {
        if c.Namespace == ConfNamespaceCommon && c.Key == share.ConfigKeyZKAddr {
            zkaddr = c.Value
            break
        }
    }
    if zkaddr == "" {
        return errors.New("no zkaddr config specified!")
    }

    err = initZK(zkaddr)
    if err != nil {
        return err
    }

    for _, n := range namespaces {
        err = attachNamespaceWithZK(n)
        if err != nil {
            log.Error(err.Error())
        }
    }
    for _, c := range confs {
        log.Debug("attach %v %v", c.Namespace, c.Key)
        err = attachWithZK(c.Namespace, c.Key)
        if err != nil {
            log.Error(err.Error())
        }
    }

    return nil
}

func finiCore()  {
    finiDB()
    finiZK()
}

func updateConfig(namespace, key, value string) error {
    var opConf *Config
    for _, conf := range confs {
        if conf.Namespace == namespace && conf.Key == key {
            opConf = conf
            break
        }
    }
    if opConf == nil {
        return fmt.Errorf("error update: config<%v %v> not found.", namespace, key)
    }
    if opConf.Value == value {
        return fmt.Errorf("error update: config<%v %v> unchanged", namespace, key)
    }
    return updateByConfig(opConf, value)
}

func addConfig(namespace, key, value string) error {
    for _, conf := range confs {
        if conf.Namespace == namespace && conf.Key == key {
            return fmt.Errorf("duplicated entry: %v %v", namespace, key)
        }
    }
    var err error
    var addConf Config
    addConf.Namespace = namespace
    addConf.Key = key
    addConf.Value = value
    err = addDB(&addConf)
    if err != nil {
        return err
    }
    return nil
}

func ConfigWithNamespaceKey(nameSpace string, keys []string) (map[string][]string, error) {
    rets := make(map[string][]string)
    //common的固定返回
    for _, key := range keys {
        //先取common的值
        for _, c := range confs {
            if c.Key == key && c.Namespace == ConfNamespaceCommon {
                rets[key] = []string{c.Namespace, c.Value}
                break
            }
        }
        //再取特定namespace的值，同key的覆盖
        for _, c := range confs {
            if c.Key == key && c.Namespace == nameSpace {
                rets[key] = []string{c.Namespace, c.Value}
                break
            }
        }
        if _, ok := rets[key]; !ok {
            return nil, fmt.Errorf("config for key <%v> not specified!", key)
        }
    }

    return rets, nil
}

func AllConfig() []Config {
    var results []Config
    for _, c := range confs {
        results = append(results, *c)
    }
    return results
}

func updateByConfig(opConf *Config, value string) error {
    opConf.Value = value
    err := updateDB(opConf)
    if err != nil {
        return err
    }
    return notifyWithZK(opConf.Namespace, opConf.Key)
}