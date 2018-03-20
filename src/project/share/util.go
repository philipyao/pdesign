package share

import (
    "fmt"
    "os"
    "path/filepath"

    log "github.com/philipyao/toolbox/logging"
    "base/srv"
)

//初始化业务日志
func InitBizLog(maxSize int) {
    //@1 业务log初始化
    config := `{"filename": "%v", "maxsize": %v, "maxbackup": 10}`
    wd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    logName := filepath.Join(wd, "log", srv.ProcessName())
    config = fmt.Sprintf(config, logName, maxSize)
    fmt.Printf("log config: %+v\n", config)
    err = log.AddAdapter(log.AdapterFile, config)
    if err != nil {
        panic(err)
    }
    log.SetLevel(log.LevelStringDebug)
    log.SetFlags(log.LogDate | log.LogTime | log.LogMicroTime | log.LogLongFile)
}

//设置srv框架自定义日志
func SetSrvLogger() {
    srv.SetLogger(log.Output)
}
