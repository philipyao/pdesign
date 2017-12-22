package pconfclient

import (
    "fmt"
    "sync"
    "strconv"
    "reflect"
    "strings"
)


func RegisterConfDef(confDef interface{}) error {
    t := reflect.TypeOf(confDef)
    v := reflect.ValueOf(confDef)
    if t.Kind() != reflect.Ptr {
        return fmt.Errorf("confdef should be pointer.")
    }
    if t.NumMethod() == 0 {
        return fmt.Errorf("confdef with XXX|SetXXX|OnUpdateXXX should be defined.")
    }
    for i := 0; i < t.NumMethod(); i++ {
        name := t.Method(i).Name
        if strings.HasPrefix(name, PrefixSet) {
            name = strings.TrimPrefix(name, PrefixSet)
        } else if strings.HasPrefix(name, PrefixOnUpdate) {
            name = strings.TrimPrefix(name, PrefixOnUpdate)
        }
        nameKey := gonicCasedName(name)
        _, ok := processors[nameKey]
        if !ok {
            processors[nameKey] = &confProcessor{
                NameKey:    nameKey,
                Name:       name,
            }
        }
    }
    return initProcessorFunc(v)
}

func Load() chan bool {
    done := make(chan bool, 1)
    //开始从远程服务器加载需要的配置

    return done
}



func gonicCasedName(name string) string {
    newstr := make([]rune, 0, len(name)+3)
    for idx, chr := range name {
        if isASCIIUpper(chr) && idx > 0 {
            if !isASCIIUpper(newstr[len(newstr)-1]) {
                newstr = append(newstr, '_')
            }
        }

        if !isASCIIUpper(chr) && idx > 1 {
            l := len(newstr)
            if isASCIIUpper(newstr[l-1]) && isASCIIUpper(newstr[l-2]) {
                newstr = append(newstr, newstr[l-1])
                newstr[l-1] = '_'
            }
        }

        newstr = append(newstr, chr)
    }
    return strings.ToLower(string(newstr))
}

func isASCIIUpper(r rune) bool {
    return 'A' <= r && r <= 'Z'
}