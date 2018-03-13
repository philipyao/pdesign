package core

import (
    "fmt"
    "time"
    "errors"

    "base/log"
    "project/share/commdef"

    "project/public/confsvr/def"
    "project/public/confsvr/db"
)

var (
    c Configure = Configure{}
)

func Init() error {
    var (
        err error
        dbUser def.User
        dbNamespace def.Namespace
        dbConfig def.Config
        dbOpLog def.ConfigOplog
    )
    err = db.Init(&dbUser, &dbNamespace, &dbConfig, &dbOpLog)
    if err != nil {
        return err
    }

    err = prepareDBData()

    log.Debug("loadConfigFromDB...")
    confs, namespaces, err := db.LoadConfigAll()
    if err != nil {
        return err
    }
    log.Debug("loadConfigFromDB ok, count: %v", len(confs))
    c.Load(confs)
    ns.Load(namespaces)

    zkConf := c.GetBy(def.ConfNamespaceCommon, commdef.ConfigKeyZKAddr)
    if zkConf == nil {
        return errors.New("no zkaddr config specified!")
    }
    err = initZK(zkConf.Value)
    if err != nil {
        return err
    }

    err = c.Foreach(func(entry *def.Config) error {
        return attachWithZK(entry.Namespace, entry.Key)
    })
    if err != nil {
        return err
    }

    return nil
}

func Fini()  {
    db.Fini()
    finiZK()
}

func UpdateConfig(id uint, value string, version int) error {
    var opConf *def.Config = c.Get(id)
    if opConf == nil {
        return fmt.Errorf("error update: config<%v> not found", id)
    }
	if opConf.Version != version {
        return fmt.Errorf("error update: config<%v> version mismatch<%v %v>", id, opConf.Version, version)
	}
    if opConf.Value == value {
        return fmt.Errorf("error update: config<%v> unchanged", id)
    }
    return updateByConfig(opConf, value)
}

func AddOplog(name, comment, author string, changes []*def.OpChange) {
    oplog := &def.ConfigOplog{
        Name: name,
        Comment: comment,
        Changes: changes,
        Author: author,
        OpTime: time.Now(),
    }
    err := db.InsertOplog(oplog)
    if err != nil {
        log.Error("add oplog<%+v> error: %v", oplog, err)
        return
    }
}

func AddConfig(namespace, key, value string) (*def.Config, error) {
    if !ns.Exist(namespace) {
        return nil, errors.New("namespace not exist")
    }
    if c.GetBy(namespace, key) != nil {
        return nil, fmt.Errorf("duplicated entry: %v %v", namespace, key)
    }

    conf, err := c.Set(namespace, key, value)
    if err != nil {
        return nil, err
    }
    return conf, nil
}

func ConfigByID(id uint) *def.Config {
    return c.Get(id)
}

func ConfigWithNamespaceKey(nameSpace string, keys []string) (map[string][]string, error) {
    rets := make(map[string][]string)
    //common的固定返回
    var conf *def.Config
    for _, key := range keys {
        //先取common的值
        conf = c.GetBy(def.ConfNamespaceCommon, key)
        if conf != nil {
            rets[key] = []string{conf.Namespace, conf.Value}
        }
        //再取特定namespace的值，同key的覆盖
        conf = c.GetBy(nameSpace, key)
        if conf != nil {
            rets[key] = []string{conf.Namespace, conf.Value}
        }
        if _, ok := rets[key]; !ok {
            return nil, fmt.Errorf("config for key <%v> not specified!", key)
        }
    }

    return rets, nil
}

func AllConfig() []def.Config {
    var results []def.Config
    c.Foreach(func(entry *def.Config) error {
        results = append(results, *entry)
        return nil
    })
    return results
}

//==========================================================

func prepareDBData() error {
    var err error

    //create admin user
    err = createAdmin()
    if err != nil {
        return err
    }
    // create common namespace
    err = ns.CreateCommon()
    if err != nil {
        return err
    }
    return nil
}

func updateByConfig(opConf *def.Config, value string) error {
    opConf.Value = value
    err := db.UpdateConfig(opConf)
    if err != nil {
        return err
    }
    return notifyWithZK(opConf.Namespace, opConf.Key)
}
