package pconfclient

import (
    "fmt"
    "reflect"
    "errors"

    "project/share/commdef"
    "project/share/pconf.client/core"
)

const (
    PConfTag            = "pconf"
    NameKeyZKAddr       = "zkaddr"
)

var (
    currNamespace   string
    logFunc         func(string, ...interface{})
)

func init() {
    defaultLogger := func(format string, args ...interface{}) {
        fmt.Printf(format, args)
    }
    SetLogger(defaultLogger)
}

func RegisterConfDef(namespace string, confDef interface{}) error {
    var err error

    confAddr := "10.1.164.99:12001"
    err = core.InitFetcher(confAddr)
    if err != nil {
        return err
    }

    if namespace == "" {
        return errors.New("empty conf namespace.")
    }
    currNamespace = namespace

    t := reflect.TypeOf(confDef)
    v := reflect.ValueOf(confDef)
    if t.Kind() != reflect.Ptr {
        return fmt.Errorf("confdef should be pointer.")
    }
    t = t.Elem()
    if t.Kind() != reflect.Struct {
        return fmt.Errorf("confdef should be struct pointer. %v", reflect.TypeOf(t.Elem()).Kind())
    }
    if t.NumField() == 0 {
        return fmt.Errorf("confdef with no fields.")
    }
    //找出'pconf' tag的字段
    tagFound := false
    for i := 0; i < t.NumField(); i++ {
        sf := t.Field(i)
        tag, ok := sf.Tag.Lookup(PConfTag)
        if !ok { continue }
        if tag == "" {
            return fmt.Errorf("empty value of 'pconf' tag for field <%v> is not allowed.", sf.Name)
        }
        tagFound = true
        goName := tag2GoName(tag)
        err = core.RegisterEntry(tag, goName, v)
        if err != nil {
            return err
        }
    }
    if tagFound == false {
        return fmt.Errorf("no 'pconf' tag found in provided confdef")
    }
    return nil
}

func SetLogger(l func(string, ...interface{})) {
    logFunc = func(format string, args ...interface{}) {
        l("[pconfclient] " + format, args)
    }
    core.SetLogger(l)
}

func Load(done chan struct{}) error {
    //开始从远程服务器加载需要的配置
    keys := core.EntryKeys()
    logFunc("start loading confs: count %v", len(keys))
    keys = append(keys, NameKeyZKAddr)
    confs, err := core.FetchConfFromServer(currNamespace, keys)
    if err != nil {
        return err
    }
    logFunc("fetch confs from confsvr ok.")

    // get zkaddr
    var zkaddr string
    for _, c := range confs {
        if c.Key == commdef.ConfigKeyZKAddr {
            zkaddr = c.Value
        }
    }
    if zkaddr == "" {
        return errors.New("zkaddr config not found")
    }
    err = core.InitWatcher(zkaddr)
    if err != nil {
        return err
    }

    notify := make(chan string)
    for _, c := range confs {
        if c.Key == commdef.ConfigKeyZKAddr {
            continue
        }
        err = core.InitEntry(c.Key, c.Value)
        if err != nil {
            return err
        }
        err = core.WatchEntryUpdate(c.Namespace, c.Key, notify, done)
        if err != nil {
            return err
        }
        logFunc("watch entry <%v %v>.", c.Namespace, c.Key)
    }

    // listen updates
    go handleWatch(notify, done)

    logFunc("finished loading confs.")
    return nil
}

func handleWatch(notify chan string, done chan struct{}) {
    for {
        select {
        case <- done:
            return
        case updateKey := <- notify:
            handleUpdate(updateKey)
        }
    }
}

func handleUpdate(key string) {
    confs, err := core.FetchConfFromServer(currNamespace, []string{key})
    if err != nil {
        logFunc("FetchConfFromServer: %v", err)
        return
    }
    if len(confs) != 1 {
        logFunc("inv conf counts")
        return
    }
    if confs[0].Key != key {
        logFunc("mismatch key %v %+v", key, confs[0])
        return
    }
    logFunc("handleUpdate: key %v", key)
    //TODO
    core.UpdateEntry(key, "", confs[0].Value)
}