package share

import (
    "fmt"
    "os"
    "path/filepath"

    "base/log"
    "base/srv"
)

//服务器通用的设置log接口
func SetServerLog(maxSize int) {
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
