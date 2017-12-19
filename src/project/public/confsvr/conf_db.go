package main

import (
    "errors"
    "fmt"
    "time"

    //"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"

    "github.com/go-xorm/xorm"
    "github.com/go-xorm/core"
)

var (
    //db *gorm.DB
    engine *xorm.Engine

    confs []DBConfig
)

const (
    TableNameConf       = "tbl_conf"

    ConfNamespaceOms    = "oms"
    ConfNamespaceCommon = "common"
)

//配置表
type DBConfig struct {
    ID              uint        `xorm:"pk autoincr 'id'"`
    Namespace       string      `xorm:"varchar(32) notnull"`
    Key             string      `xorm:"varchar(64) notnull"`
    Value           string      `xorm:"varchar(128) notnull"`

    UpdatedAt       time.Time   `xorm:"updated"`
    CreatedAt       time.Time   `xorm:"created"`
    Version         int         `xorm:"version"`    //自动更新版本号
}
func (c DBConfig) TableName() string {
    return TableNameConf
}

type UserSimu struct {
    Charid          uint64      `xorm:"pk notnull 'charid'"`
    Accid           uint        `xorm:"notnull"`
    Name            string      `xorm:"varchar(128) notnull"`
    BaseData        []uint8     `xorm:"mediumblob"`

    UpdatedAt       time.Time   `xorm:"updated"`
    CreatedAt       time.Time   `xorm:"created"`
}

func DBInit() error {
    var err error
    //db, err = gorm.Open("mysql", "hgame:Hgame188@tcp(10.1.164.20:3306)/db_new_oms?charset=utf8&parseTime=True&loc=Local")
    //if err != nil {
    //    Log.Printf("gorm.Open error %v\n", err)
    //    return err
    //}
    //db = db.Debug()
    //row := db.Select("VERSION()").Row()
    //Log.Printf("version: %+v\n", row)

    engine, err = xorm.NewEngine("mysql", "hgame:Hgame188@tcp(10.1.164.20:3306)/db_new_oms?charset=utf8")
    if err != nil {
        Log.Printf("xorm.NewEngine error %v\n", err)
        return err
    }
    engine.ShowSQL(true)
    engine.Logger().SetLevel(core.LOG_DEBUG)
    engine.SetMapper(core.GonicMapper{})
    err = engine.Ping()
    if err != nil {
        Log.Printf("engine.Ping error %v\n", err)
        return err
    }
    return nil
}

func DBLoadConfig() error {
    var err error

    //if db == nil {
    //    return errors.New("null db")
    //}

    //err = db.Table(TableNameConf).AutoMigrate(&Config{}).Error
    //if err != nil {
    //    Log.Println("AutoMigrate %v error: %v\n", TableNameConf, err)
    //    return err
    //}
    //
    //var confs []Config
    //err = db.Table(TableNameConf).Find(&confs).Error
    //if err != nil {
    //    Log.Println("Find %v error: %v\n", TableNameConf, err)
    //    return err
    //}
    //Log.Printf("Load conf, total %v records\n", len(confs))

    if engine == nil {
        return errors.New("null engine")
    }
    err = engine.Sync2(new(DBConfig))
    if err != nil {
        Log.Println(err)
    }
    confs = make([]DBConfig, 0)
    err = engine.Find(&confs)
    if err != nil {
        Log.Println(err)
    }
    Log.Printf("confs: %+v\n", confs)
    return nil
}

func SimuCreateMulti() error {
    const TBL_USER_NUM  = 2
    var err error
    //for i := 1; i <= TBL_USER_NUM; i++{
    //    tblName := fmt.Sprintf("t_simu_user_%02d", i)
    //    err = db.DropTableIfExists(tblName).Error
    //    if err != nil {
    //        Log.Println("DropTableIfExists %v error: %v\n", tblName, err)
    //        return err
    //    }
    //    err = db.Table(tblName).Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&UserSimu{}).Error
    //    if err != nil {
    //        Log.Println("CreateTable %v error: %v\n", tblName, err)
    //        return err
    //    }
    //}
    for i := 1; i <= TBL_USER_NUM; i++{
        tblName := fmt.Sprintf("t_simu_user_%02d", i)
        err = engine.Table(tblName).Sync2(&UserSimu{})
        if err != nil {
            Log.Println("Sync2 %v error: %v\n", tblName, err)
            return err
        }
    }
    return nil
}

func ConfigWithNamespace(nameSpace string) map[string]string {
    rets := make(map[string]string)
    //common的固定返回
    for _, c := range confs {
        if c.Namespace == ConfNamespaceCommon {
            rets[c.Key] = c.Value
        }
    }
    //特定namespace的同key的可以覆盖common中的
    for _, c := range confs {
        if c.Namespace == nameSpace {
            rets[c.Key] = c.Value
        }
    }
    return rets
}

func DBFini() {
    if engine != nil {
        engine.Close()
    }
}
