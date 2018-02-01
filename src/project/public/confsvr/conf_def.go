package main

import (
    "time"
)

const (
    ConfNamespaceOms    = "oms"
    ConfNamespaceCommon = "common"
    TableNameConf       = "tbl_conf"
    TableNameConfOplog  = "tbl_conf_oplog"
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

//操作日志
type ConfigOplog struct {
    OpID            uint        `xorm:"pk autoincr 'opid'"`
    Name            string      `xorm:"varchar(64) notnull"`
    Comment         string      `xorm:"varchar(128) notnull"`
    Changes         []*OpChange `xorm:"text notnull"`
    Author          string      `xorm:"varchar(32) notnull"`
    OpTime          time.Time   `xorm:"DateTime notnull"`
}
type OpChange struct {
    Namespace       string
    Key             string
    OldValue        string
    Value           string
}
func (co ConfigOplog) TableName() string {
    return TableNameConfOplog
}

