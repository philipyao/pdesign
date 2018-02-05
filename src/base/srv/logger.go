package srv

import (
    "fmt"
    "log"
)

const (
    LogPrefix           = "[server]"
    LogCalldepth        = 3
)

func defaultLogFunc() func(format string, args ...interface{}) {
    return func(format string, args ...interface{}) {
        //默认输出到stdout, 借助golang官方log包
        s := LogPrefix + fmt.Sprintf(format, args...)
        log.Output(LogCalldepth, s)
    }
}

//自定义log输出
func customLogFunc(l func(int, string, ...interface{})) func(format string, args ...interface{}) {
    return func(format string, args ...interface{}) {
        fmt.Printf("custom log: %v\n", format)
        l(LogCalldepth, LogPrefix + format, args...)
    }
}

