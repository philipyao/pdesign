package main

import (
    "fmt"
    "sync"
    "strconv"
    "reflect"
    "strings"

    "base/log"

    _ "project/share/pconf.client/core"
)

type GameConf struct {
    logLevel        string          `pconf:"log_level"`
    vipLevel        uint32          `pconf:"vip_level"`
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


func HandleConfChange(key, oldVal, val string) {

}



