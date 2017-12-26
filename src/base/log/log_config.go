package log

import (
    "log"
    "encoding/json"
)

type logConfigFile struct {
    fileName    string    `json:"filename"`
    maxSize     int64     `json:"maxsize"`
    maxBackup   int       `json:"maxbackup"`
}
type logConfigNet struct {
    net         string    `json:"net"`
    addr        string    `json:"addr"`
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