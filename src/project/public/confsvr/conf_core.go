package main

import (
    "fmt"
    "time"
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
    confs []Config
)

func initCore() error {
    var err error
    err = initDB(new(Config))
    if err != nil {
        return err
    }
    confs, err = loadConfigFromDB()
    if err != nil {
        return err
    }
    Log.Printf("confs: %+v\n", confs)

    return nil
}

func ConfigWithNamespaceKey(nameSpace string, keys []string) (map[string]string, error) {
    rets := make(map[string]string)
    //common的固定返回
    for _, key := range keys {
        //先取common的值
        for _, c := range confs {
            if c.Key == key && c.Namespace == ConfNamespaceCommon {
                rets[key] = c.Value
                break
            }
        }
        //再取特定namespace的值，同key的覆盖
        for _, c := range confs {
            if c.Key == key && c.Namespace == nameSpace {
                rets[key] = c.Value
                break
            }
        }
        if _, ok := rets[key]; !ok {
            return nil, fmt.Errorf("config for key <%v> not specified!", key)
        }
    }

    return rets, nil
}
