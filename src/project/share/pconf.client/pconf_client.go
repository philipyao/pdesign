package pconfclient

import (
    "fmt"
    "reflect"

    "project/share/pconf.client/core"
)

const (
    PConfTag            = "pconf"
)

func RegisterConfDef(confDef interface{}) error {
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
    var err error
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

func Load() chan bool {
    done := make(chan bool, 1)
    //开始从远程服务器加载需要的配置

    return done
}
