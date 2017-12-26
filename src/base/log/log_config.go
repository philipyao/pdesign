package log

import (
    "log"
    "encoding/json"
)

type logConfigFile struct {
    FileName    string    `json:"filename"`
    MaxSize     int64     `json:"maxsize"`
    MaxBackup   int       `json:"maxbackup"`
}
type logConfigNet struct {
    Net         string    `json:"net"`
    Addr        string    `json:"addr"`
}

type logConfig struct {
    logConfigFile
    logConfigNet
}

func loadLogConfig(conf string) *logConfig {
    var lc logConfig
    err := json.Unmarshal([]byte(conf), &lc)
    if err != nil {
        log.Printf("load log config error: %v, conf %v", err, conf)
    }
    return &lc
}