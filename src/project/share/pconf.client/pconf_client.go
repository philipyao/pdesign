package pconfclient

import (
    "fmt"
    "reflect"
    "errors"

    "project/share"
    "project/share/pconf.client/core"
)

const (
    PConfTag            = "pconf"
    NameKeyZKAddr       = "zkaddr"
)

var (
    currNamespace   string
)

func RegisterConfDef(namespace string, confDef interface{}) error {
    var err error

    confAddr := "10.1.164.45:12001"
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
        fmt.Printf("register field %v %v ok\n", tag, goName)
    }
    if tagFound == false {
        return fmt.Errorf("no 'pconf' tag found in provided confdef")
    }
    return nil
}

func Load(done chan struct{}) error {
    //开始从远程服务器加载需要的配置
    keys := core.EntryKeys()
    keys = append(keys, NameKeyZKAddr)
    fmt.Printf("load: keys %+v\n", keys)

    confs, err := core.FetchConfFromServer(currNamespace, keys)
    if err != nil {
        return err
    }
    //TODO get zkaddr
    var zkaddr string
    for _, c := range confs {
        if c.Key == share.ConfigKeyZKAddr {
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
    for i, c := range confs {
        fmt.Printf("fetched confs[%v]: %+v\n", i, c)
        if c.Key == share.ConfigKeyZKAddr {
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
        fmt.Printf("watch entry <%v %v>ok\n", c.Namespace, c.Key)
    }

    // listen updates
    go handleWatch(notify, done)

    return nil
}

func handleWatch(notify chan string, done chan struct{}) {
    for {
        select {
        case <- done:
            return
        case updateKey := <- notify:
            fmt.Printf("update key %v\n", updateKey)
            handleUpdate(updateKey)
        }
    }
}

func handleUpdate(key string) {
    confs, err := core.FetchConfFromServer(currNamespace, []string{key})
    if err != nil {
        fmt.Println(err)
        return
    }
    if len(confs) != 1 {
        fmt.Println("inv conf counts")
        return
    }
    if confs[0].Key != key {
        fmt.Printf("mismatch key %v %+v\n", key, confs[0])
        return
    }
    fmt.Printf("update: %+v\n", confs[0])
    //TODO
    core.UpdateEntry(key, "", confs[0].Value)
}