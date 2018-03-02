package def

import (
    "time"
)

const (
    ConfNamespaceOms    = "oms"
    ConfNamespaceCommon = "common"

    AdminUsername       = "admin"
    AdminPasswd         = "hellopconf"
    DefaultSaltLen      = 32
    DefaultCliPasswdLen = 40    //sha1输出为40位
    ClientSaltPart      = "^rR@8=YlsU"

    TableNameUser       = "tbl_user"
    TableNameConf       = "tbl_conf"
    TableNameNamespace  = "tbl_namespace"
    TableNameConfOplog  = "tbl_conf_oplog"
)

//用户表
type User struct {
    Username        string      `xorm:"pk varchar(32) notnull"`
    Passwd          string      `xorm:"varchar(128) notnull"`
    Enabled         uint        `xorm:"'enabled'"`
    IsSuper         uint        `xorm:"'is_super'"` //是否超级用户
    Salt            string      `xorm:"varchar(32) notnull"`    //随机密码盐
    UpdatedAt       time.Time   `xorm:"updated"`
    CreatedAt       time.Time   `xorm:"created"`
}
func (u User) TableName() string {
    return TableNameUser
}

//namespace表
type Namespace struct {
    Name            string      `xorm:"pk varchar(32) notnull"`
    Desc            string      `xorm:"varchar(128)"`
    Creator         string      `xorm:"varchar(32)"`
    CreatedAt       time.Time   `xorm:"created"`
}
func (n Namespace) TableName() string {
    return TableNameNamespace
}

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

