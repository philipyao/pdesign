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


)

type UserSimu struct {
    Charid          uint64      `xorm:"pk notnull 'charid'"`
    Accid           uint        `xorm:"notnull"`
    Name            string      `xorm:"varchar(128) notnull"`
    BaseData        []uint8     `xorm:"mediumblob"`

    UpdatedAt       time.Time   `xorm:"updated"`
    CreatedAt       time.Time   `xorm:"created"`
}

func initDB(obj *Config) error {
    var err error
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
    err = engine.Sync2(obj)
    if err != nil {
        Log.Println(err)
    }
    return nil
}

func loadConfigFromDB() (confs []Config, err error) {
    if engine == nil {
        return nil, errors.New("null engine")
    }

    confs = make([]Config, 0)
    err = engine.Find(&confs)
    if err != nil {
        return nil, err
    }

    return
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

func DBFini() {
    if engine != nil {
        engine.Close()
    }
}
