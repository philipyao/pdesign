package core

import (
    "fmt"
    "log"
)

const (
    LogPrefix           = "[pconfclient]"
    LogCalldepth        = 4
)

var (
    logger  func(string, ...interface{})
)

func init() {
    //默认输出到stdout, 借助golang官方log包
    logger = func(format string, args ...interface{}) {
        s := LogPrefix + fmt.Sprintf(format, args...)
        log.Output(LogCalldepth, s)
    }
}

//自定义log输出
func SetLogger(l func(int, string, ...interface{})) {
    logger = func(format string, args ...interface{}) {
        l(LogCalldepth, LogPrefix + format, args...)
    }
}

func Log(format string, args ...interface{}) {
    logger(format, args...)
}
