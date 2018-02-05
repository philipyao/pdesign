package log

import (
    "os"
    "fmt"
    "time"
    "runtime"
    "strconv"
    "path/filepath"

    "base/log/adapter"
)

const (
    _ int              = iota
    LevelDebug
    LevelInfo
    LevelWarn
    LevelError
    LevelFatal

    LevelStringDebug    = "DEBUG"
    LevelStringInfo     = "INFO"
    LevelStringWarn     = "WARN"
    LevelStringError    = "ERROR"
    LevelStringFatal    = "FATAL"
)


type logFlag    int8
const (
    _ logFlag       = (1 << iota)
    LogDate
    LogTime
    LogMicroTime
    LogLongFile
    LogShortFile
    LogStd          = LogDate | LogTime | LogShortFile
)

const (
    LogChanSize           = 1024000
    DefaultCalldepth      = 2
)

const (
    AdapterConsole      = "console"
    AdapterFile         = "file"
    AdapterNet          = "net"
)

var (
    adapters    map[string]adapter.Adapter

    level       string
    lvs         map[string]int

    flag        logFlag

    logChan     chan *logMessage
    doneChan    chan struct{}
)

func init() {
    adapters    = make(map[string]adapter.Adapter)
    lvs         = map[string]int{
        LevelStringDebug:   LevelDebug,
        LevelStringInfo:    LevelInfo,
        LevelStringWarn:    LevelWarn,
        LevelStringError:   LevelError,
        LevelStringFatal:   LevelFatal,
    }

    //默认输出INFO
    level = LevelStringInfo
    flag = LogStd

    logChan = make(chan *logMessage, LogChanSize)
    doneChan = make(chan struct{}, 1)

    go handleWriteLog()
}

func AddAdapter(name string, conf string) error {
    var err error
    if _, ok := adapters[name]; ok {
        return fmt.Errorf("duplicated adapter name %v", name)
    }

    logconf := loadLogConfig(conf)
    if logconf == nil {
        return fmt.Errorf("parse log json config error: %v", conf)
    }
    var adp adapter.Adapter
    if name == AdapterConsole {
        panic(name)
    } else if name == AdapterFile {
        options := &adapter.Options{
            MaxSize: adapter.ByteSize(logconf.MaxSize),
            MaxBackup: logconf.MaxBackup,
        }
        adp, err = adapter.NewAdapterFile(logconf.FileName, options)
    } else if name == AdapterNet {
        panic(name)
    } else {
        err = fmt.Errorf("unknown adapter name %v", name)
    }
    if err != nil {
        return err
    }
    adapters[name] = adp
    return nil
}

func RemoveAdapter(name string) error {
    delete(adapters, name)
    return nil
}

func CheckLevel(lv string) bool {
    _, ok := lvs[lv]
    return ok
}

func SetLevel(lv string) error {
    if _, ok := lvs[lv]; !ok {
        return fmt.Errorf("invalid log level %v", lv)
    }
    level = lv
    return nil
}

func SetFlags(f logFlag) {
    flag = f
}

func Debug(format string, args ...interface{}) {
    if lvs[level] > LevelDebug {
        return
    }
    output(DefaultCalldepth, LevelStringDebug, format, args...)
}

func Info(format string, args ...interface{}) {
    if lvs[level] > LevelInfo {
        return
    }
    output(DefaultCalldepth, LevelStringInfo, format, args...)
}

func Warn(format string, args ...interface{}) {
    if lvs[level] > LevelWarn {
        return
    }
    output(DefaultCalldepth, LevelStringWarn, format, args...)
}
func Error(format string, args ...interface{}) {
    if lvs[level] > LevelError {
        return
    }
    output(DefaultCalldepth, LevelStringError, format, args...)
}
func Fatal(format string, args ...interface{}) {
    if lvs[level] > LevelFatal {
        return
    }
    output(DefaultCalldepth, LevelStringFatal, format, args...)

    //退出前，将所有log flush掉，阻塞等
    Flush()
    os.Exit(1)
}
func Output(calldepth int, format string, args ...interface{}) {
    output(calldepth, "", format, args...)
}

func Flush() {
    //结束log监听写，将剩余log写入后退出for循环
    close(logChan)
    logChan = nil

    //等待所有日志写完（优化TODO：等待一定时间）
    <-doneChan
    for _, adp := range adapters {
        adp.Close()
    }
}


//////////////////////////////////////////////////////////////////////

func output(calldepth int, lvString string, format string, args ...interface{}) {
    if logChan == nil { return }

    var text string
    tmNow := time.Now()
    if flag & LogDate != 0 {
        y, m, d := tmNow.Date()
        text += fmt.Sprintf("%04d-%02d-%02d", y, int(m), d)
        text += " "
    }
    if flag & LogTime != 0 {
        text += tmNow.Format("15:04:05")
        if flag & LogMicroTime != 0 {
            text += fmt.Sprintf(".%06d", tmNow.Nanosecond()/1e3)
        }
        text += " "
    }

    if flag & (LogShortFile | LogLongFile) != 0 {
        _, file, line, ok := runtime.Caller(calldepth)
        if !ok {
            file = "???"
            line = 0
        } else {
            if flag & LogShortFile != 0 {
                file = filepath.Base(file)
            }
        }
        text += file + ":" + strconv.Itoa(line)
        text += " "
    }

    if lvString != "" {
        text += fmt.Sprintf("[%v] ", lvString)
    }
    text += fmt.Sprintf(format, args...)
    text += "\n"

    logMsg := logMessageGet()
    logMsg.Buff = []byte(text)
    if len(logChan) == LogChanSize {
        fmt.Println("FULL!!!!!")
        return
    }
    logChan <- logMsg
}

func handleWriteLog() {
    //要结束该for range循环，可以close(logChan)
    for logMsg := range logChan {
        for _, adp  := range adapters {
            adp.Write(logMsg.Buff)
        }
        logMessagePut(logMsg)
    }

    //至此，所有日志写完毕，通知监听者
    doneChan <- struct{}{}
}
