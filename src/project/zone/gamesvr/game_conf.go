package main

import (
    "fmt"
    "sync"
    "strconv"
    "reflect"
    "strings"

    "base/log"
)

const (
    PrefixSet       = "Set"
    PrefixOnUpdate  = "OnUpdate"
)

var (
    ConfMgr *GameConf

    processors  map[string]*ConfProcessor
)

type FuncProcessSet func(string) error
type FuncProcessOnUpdate func(string, string)

type ConfProcessor struct {
    NameKey         string          //小写下划线配置名字
    Name            string
    FuncGet         reflect.Value   //获取函数
    FuncSet         reflect.Value   //设置函数
    FuncOnUpdate    reflect.Value   //更新回调
}

// gamesvr 用到的conf定义
type GameConf struct {
    lock *sync.RWMutex

    logLevel        string
    enableVip       int
}

func init() {
    ConfMgr = new(GameConf)
    ConfMgr.lock = new(sync.RWMutex)
    processors = make(map[string]*ConfProcessor)
}

func (gc *GameConf) LogLevel() string {
    gc.lock.RLock()
    defer gc.lock.RUnlock()

    return gc.logLevel
}
func (gc *GameConf) SetLogLevel(val string) error {
    gc.lock.Lock()
    defer gc.lock.Unlock()

    if val != log.LevelStringDebug &&
        val != log.LevelStringInfo &&
        val != log.LevelStringWarn &&
        val != log.LevelStringError &&
        val != log.LevelStringCrit {
        return fmt.Errorf("invalid log level: %v", val)
    }
    gc.logLevel = val
    return nil
}
func (gc *GameConf) OnUpdateLogLevel(oldVal, val string) {
    log.SetLevel(val)
}

func (gc *GameConf) EnableVip() int {
    gc.lock.RLock()
    defer gc.lock.RUnlock()

    return gc.enableVip
}
func (gc *GameConf) SetEnableVip(val string) error {
    gc.lock.Lock()
    defer gc.lock.Unlock()

    if len(val) == 0 {
        return fmt.Errorf("empty set string")
    }
    i, err := strconv.Atoi(val)
    if err != nil {
        return err
    }
    if i != 0 && i != 1 {
        return fmt.Errorf("invalid enable_vip value: %v", i)
    }
    gc.enableVip = i
    return nil
}
func (gc *GameConf) OnUpdateEnableVip(oldVal, val string) {
    //do nothing
}

//============================================================
func RegisterConfDef() interface{} {
    mgr := ConfMgr
    t := reflect.TypeOf(mgr)
    v := reflect.ValueOf(mgr)
    name := reflect.Indirect(v).Type().Name()
    log.Debug("confdef: %v, gonic: %v", name, gonicCasedName(name))
    if t.Kind() != reflect.Ptr {
        log.Error("confdef should be pointer.")
        panic(t)
    }
    for i := 0; i < t.NumMethod(); i++ {
        log.Debug("method: %+v", t.Method(i))
        name = t.Method(i).Name
        if strings.HasPrefix(name, PrefixSet) {
            name = strings.TrimPrefix(name, PrefixSet)
        } else if strings.HasPrefix(name, PrefixOnUpdate) {
            name = strings.TrimPrefix(name, PrefixOnUpdate)
        }
        nameKey := gonicCasedName(name)
        _, ok := processors[nameKey]
        if !ok {
            processors[nameKey] = &ConfProcessor{
                NameKey:    nameKey,
                Name:       name,
            }
        }
    }

    errorInterface := reflect.TypeOf((*error)(nil)).Elem()
    for k, p := range processors {
        p.FuncGet = v.MethodByName(p.Name)
        if p.FuncGet.IsValid() == false {
            log.Error("error conf<key %v, name %v>: get method not found!", k, p.Name)
        } else {
            if p.FuncGet.Type().NumIn() > 0 {
                log.Error("error conf<key %v, name %v>: get method must NOT have parameters", k, p.Name)
            }
            if p.FuncGet.Type().NumOut() == 0 {
                log.Error("error conf<key %v, name %v>: get method without output parameters!", k, p.Name)
            }
        }

        p.FuncSet = v.MethodByName(PrefixSet + p.Name)
        if p.FuncSet.IsValid() == false {
            log.Error("error conf<key %v, name %v>: set method not found!", k, p.Name)
        } else {
            if p.FuncSet.Type().NumIn() != 1 {
                log.Error("error conf<key %v, name %v>: Set method with input parameter count not equal to 1", k, p.Name)
            } else {
                param := p.FuncSet.Type().In(0)
                if param.Kind() != reflect.String {
                    log.Error("error conf<key %v, name %v>: Set method parameter must be string", k, p.Name)
                }
            }
            if p.FuncSet.Type().NumOut() != 1 {
                log.Error("error conf<key %v, name %v>: set method must has 1 erro return value!", k, p.Name)
            } else {
                param := p.FuncSet.Type().Out(0)
                if param.Kind() != reflect.Interface || param.Implements(errorInterface) == false {
                    log.Error("error conf<key %v, name %v>: Set method return must be type error, kind %v", k, p.Name, param.Kind())
                }
            }
        }

        p.FuncOnUpdate = v.MethodByName(PrefixOnUpdate + p.Name)
        if p.FuncOnUpdate.IsValid() == false {
            log.Error("error conf<key %v, name %v>: on-update method not found!", k, p.Name)
        } else {
            //TODO update的参数校验
        }
    }
    return mgr
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

