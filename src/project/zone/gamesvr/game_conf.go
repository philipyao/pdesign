package main

import (
    "fmt"
    "strconv"

    "base/log"

    pconfclient "project/share/pconf.client"
)

var (
    ConfMgr *GameConf
)
const (
    Namespace       string      = "gamesvr"
)

type GameConf struct {
    logLevel        string          `pconf:"log_level"`
    enableVip       int             `pconf:"enable_vip"`
}

func InitConf() error {
    gconf := new(GameConf)
    err := pconfclient.RegisterConfDef(Namespace, gconf)
    if err != nil {
        log.Error("Init conf error: %v", err)
        return err
    }
    return nil
}

func LoadConf(done chan struct{}) error {
    return pconfclient.Load(done)
}

func (gc *GameConf) LogLevel() string {
    return gc.logLevel
}
func (gc *GameConf) SetLogLevel(val string) error {
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
    //log.SetLevel(val)
}

func (gc *GameConf) EnableVip() int {
    return gc.enableVip
}
func (gc *GameConf) SetEnableVip(val string) error {
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





