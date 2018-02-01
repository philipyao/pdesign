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

//getter
func (gc *GameConf) LogLevel() string {
    return gc.logLevel
}
//setter
func (gc *GameConf) SetLogLevel(val string) error {
    if !log.CheckLevel(val) {
        return fmt.Errorf("invalid log level: %v", val)
    }

    gc.logLevel = val
    return nil
}
func (gc *GameConf) OnUpdateLogLevel(oldVal, val string) {
    //log.SetLevel(val)
}

//getter
func (gc *GameConf) EnableVip() int {
    return gc.enableVip
}
//setter
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

//============================================================

func LoadConf(done chan struct{}) error {
    gconf := new(GameConf)
    err := pconfclient.RegisterConfDef(Namespace, gconf)
    if err != nil {
        return fmt.Errorf("RegisterConfDef: %v", err)
    }
    err = pconfclient.Load(done)
    if err != nil {
        return fmt.Errorf("Load: %v", err)
    }
    log.Info("LoadConf ok.")
    return nil
}




