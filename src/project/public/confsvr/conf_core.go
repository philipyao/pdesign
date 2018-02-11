package main

import (
    "fmt"
    "time"
    "errors"

    "base/log"
    "project/share/commdef"
)

var (
    namespaces []string
    confs []*Config
)

func initCore() error {
    var (
        err error
        dbUser User
        dbNamespace Namespace
        dbConfig Config
        dbOpLog ConfigOplog
    )
    err = initDB(&dbUser, &dbNamespace, &dbConfig, &dbOpLog)
    if err != nil {
        return err
    }

    err = prepareDBData()

    log.Debug("loadConfigFromDB...")
    confs, namespaces, err = loadConfigFromDB()
    if err != nil {
        return err
    }
    log.Debug("loadConfigFromDB ok, count: %v", len(confs))

    var zkaddr string
    for _, c := range confs {
        if c.Namespace == ConfNamespaceCommon && c.Key == commdef.ConfigKeyZKAddr {
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

    for _, c := range confs {
        err = attachWithZK(c.Namespace, c.Key)
        if err != nil {
            return err
        }
    }

    return nil
}

func finiCore()  {
    finiDB()
    finiZK()
}

func prepareDBData() error {
    var err error

    //create admin user
    err = createAdmin()
    if err != nil {
        return err
    }
    // create common namespace
    err = createNamespaceCommon()
    if err != nil {
        return err
    }
    return nil
}

func updateConfig(id uint, value string, version int) error {
    var opConf *Config
    for _, conf := range confs {
        if conf.ID == id {
            opConf = conf
            break
        }
    }
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

func addOplog(name, comment, author string, changes []*OpChange) {
    oplog := &ConfigOplog{
        Name: name,
        Comment: comment,
        Changes: changes,
        Author: author,
        OpTime: time.Now(),
    }
    err := dbAddOplog(oplog)
    if err != nil {
        log.Error("add oplog<%+v> error: %v", oplog, err)
        return
    }
}

func addConfig(namespace, key, value string) (*Config, error) {
    for _, conf := range confs {
        if conf.Namespace == namespace && conf.Key == key {
            return nil, fmt.Errorf("duplicated entry: %v %v", namespace, key)
        }
    }
    var err error
    var addConf Config
    addConf.Namespace = namespace
    addConf.Key = key
    addConf.Value = value
    err = addDB(&addConf)
    if err != nil {
        return nil, err
    }
    confs = append(confs, &addConf)
    addNamespace := true
    for _, n := range namespaces {
        if n == namespace {
            addNamespace = false
            break
        }
    }
    if addNamespace {
        namespaces = append(namespaces, namespace)
    }
    return &addConf, nil
}

func configByID(id uint) *Config {
	for _, c := range confs {
		if c.ID == id {
			return c
		}
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
