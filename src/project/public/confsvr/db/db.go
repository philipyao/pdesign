package db

import (
    "errors"
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"

    "github.com/go-xorm/xorm"
    "github.com/go-xorm/core"

    "base/log"

    "project/public/confsvr/def"
)

var (
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

func Init(objs ...interface{}) error {
    var err error
    //TODO 数据库地址
    engine, err = xorm.NewEngine("mysql",
        "hgame:Hgame188@tcp(10.1.164.20:3306)/db_new_oms?charset=utf8")
    if err != nil {
        return err
    }

    log.Info("open connection to db engine.")

    engine.ShowSQL(true)
    engine.Logger().SetLevel(core.LOG_DEBUG)
    engine.SetMapper(core.GonicMapper{})
    err = engine.Ping()
    if err != nil {
        return err
    }
    for _, obj := range objs {
        err = engine.Sync2(obj)
        if err != nil {
            return err
        }
    }

    log.Info("sync db tables ok.")
    return nil
}

func ListUser() (users []*def.User, err error) {
    if engine == nil {
        return nil, errors.New("null engine")
    }

    users = make([]*def.User, 0)
    err = engine.Find(&users)
    if err != nil {
        log.Error("engine.Find() error %v", err)
        return nil, err
    }
    return
}

func QueryUser(userName string) (*def.User, error) {
    if engine == nil {
        return nil, errors.New("null engine")
    }
    var user def.User
    has, err := engine.Id(userName).Get(&user)
    if err != nil {
        return nil, fmt.Errorf("queryUser() error %v, userName %v", err, userName)
    }
    if !has {
        return nil, nil
    }
    return &user, nil
}

func ExistUser(userName string) (bool, error) {
    if engine == nil {
        return false, errors.New("null engine")
    }
    user := &def.User{Username: userName}
    total, err := engine.Count(user)
    if err != nil {
        return false, err
    }
    return total > 0, nil
}

func InsertUser(user *def.User) error  {
    if engine == nil {
        return errors.New("null engine")
    }
    affected, err := engine.Insert(user)
    if err != nil {
        return fmt.Errorf("engine.Insert() error %v, user %+v", err, user)
    }
    if affected != 1 {
        return fmt.Errorf("inv affected %v", affected)
    }
    return nil
}

func UpdateUser(user *def.User) error {
    if engine == nil {
        return errors.New("null engine")
    }
    affected, err := engine.Id(user.Username).Cols("enabled").Update(user)
    if err != nil {
        log.Error("engine.Update() error %v, user %+v", err, user)
        return err
    }
    if affected != 1 {
        return fmt.Errorf("inv affected %v", affected)
    }
    return nil
}

func DeleteUser(userName string) error {
    user := new(def.User)
    _, err := engine.Id(userName).Delete(user)
    return err
}

func ExistNamespace(name string) (bool, error) {
    if engine == nil {
        return false, errors.New("null engine")
    }
    ns := &def.Namespace{Name: name}
    total, err := engine.Count(ns)
    if err != nil {
        log.Error("engine.Count() error %v, name %v", err, name)
        return false, err
    }
    return total > 0, nil
}

func InsertNamespace(ns *def.Namespace) error  {
    if engine == nil {
        return errors.New("null engine")
    }
    affected, err := engine.Insert(ns)
    if err != nil {
        log.Error("engine.Insert() error %v, ns %+v", err, ns)
        return err
    }
    if affected != 1 {
        return fmt.Errorf("inv affected %v", affected)
    }
    return nil
}

func LoadConfigAll() (confs []*def.Config, namespaces []string, err error) {
    if engine == nil {
        return nil, nil, errors.New("null engine")
    }

    confs = make([]*def.Config, 0)
    err = engine.Find(&confs)
    if err != nil {
        log.Error("engine.Find() error %v", err)
        return nil, nil, err
    }
    tmpMap := make(map[string]bool)
    for _, c := range confs {
        tmpMap[c.Namespace] = true
    }
    namespaces = make([]string, 0, len(tmpMap))
    for k := range tmpMap {
        namespaces = append(namespaces, k)
    }
    return
}

func UpdateConfig(opConf *def.Config) error {
    if engine == nil {
        return errors.New("null engine")
    }
    affected, err := engine.Id(opConf.ID).Cols("value").Update(opConf)
    if err != nil {
        log.Error("engine.Update() error %v, opConf %+v", err, opConf)
        return err
    }
    if affected != 1 {
        return fmt.Errorf("inv affected %v", affected)
    }
    return nil
}

func InsertConfig(opConf *def.Config) error {
    if engine == nil {
        return errors.New("null engine")
    }
    affected, err := engine.Insert(opConf)
    if err != nil {
        log.Error("engine.Insert() error %v, opConf %+v", err, opConf)
        return err
    }
    if affected != 1 {
        return fmt.Errorf("inv affected %v", affected)
    }
    return nil
}

func InsertOplog(oplog *def.ConfigOplog) error {
    if engine == nil {
        return errors.New("null engine")
    }
    affected, err := engine.Insert(oplog)
    if err != nil {
        log.Error("engine.Insert() error %v, log %+v", err, oplog)
        return err
    }
    if affected != 1 {
        return fmt.Errorf("inv affected %v", affected)
    }
    return nil
}


func SimuCreateMulti() error {
    const TBL_USER_NUM  = 2
    var err error
    for i := 1; i <= TBL_USER_NUM; i++{
        tblName := fmt.Sprintf("t_simu_user_%02d", i)
        err = engine.Table(tblName).Sync2(&UserSimu{})
        if err != nil {
            log.Error("Sync2 %v error: %v", tblName, err)
            return err
        }
    }
    return nil
}

func Fini() {
    if engine != nil {
        log.Info("close connection to db engine.")
        engine.Close()
    }
}
