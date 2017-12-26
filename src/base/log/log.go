package log

import (
    "fmt"
    "time"
    //"log"
    "runtime"
    "strings"
    "path/filepath"

    "base/log/adapter"
)

const (
    LevelDebug          = iota
    LevelInfo
    LevelWarn
    LevelError
    LevelCrit

    LevelStringDebug    = "DEBUG"
    LevelStringInfo     = "INFO"
    LevelStringWarn     = "WARN"
    LevelStringError    = "ERROR"
    LevelStringCrit     = "CRIT"
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
    LogChanSize           = 10200
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

    lvs[LevelStringDebug]   = LevelDebug
    lvs[LevelStringInfo]    = LevelInfo
    lvs[LevelStringWarn]    = LevelWarn
    lvs[LevelStringError]   = LevelError
    lvs[LevelStringCrit]    = LevelCrit
    //默认输出INFO
    level = LevelStringInfo

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

    } else if name == AdapterFile {
        options := &adapter.Options{
            MaxSize: adapter.ByteSize(logconf.maxSize),
            MaxBackup: logconf.maxBackup,
        }
        adp, err = adapter.NewAdapterFile(logconf.fileName, options)
    } else if name == AdapterNet {

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
    output(LevelStringDebug, format, args)
}

func Info(format string, args ...interface{}) {
    if lvs[level] > LevelInfo {
        return
    }
    output(LevelStringInfo, format, args)
}

func Warn(format string, args ...interface{}) {
    if lvs[level] > LevelWarn {
        return
    }
    output(LevelStringWarn, format, args)
}
func Error(format string, args ...interface{}) {
    if lvs[level] > LevelError {
        return
    }
    output(LevelStringError, format, args)
}
func Crit(format string, args ...interface{}) {
    if lvs[level] > LevelCrit {
        return
    }
    output(LevelStringCrit, format, args)
}

func Flush() {
    close(logChan)
    <-doneChan
    for _, adp := range adapters {
        adp.Close()
    }
}

func output(lvString string, format string, args ...interface{}) {
    tmNow := time.Now()
    text := tmNow.Format(time.RFC3339)
    _, file, line, ok := runtime.Caller(2)
    if !ok {
        file = "???"
        line = 0
    } else {
        if flag & LogShortFile != 0 {
            file = filepath.Base(file)
        }
    }

    fileName := file + ":" + fmt.Sprintf("%v", line)
    lvStr := fmt.Sprintf("[%v]", lvString)
    msg := fmt.Sprintf(format, args)
    text = strings.Join([]string{text, fileName, lvStr, msg}, " ")
    text += "\n"

    logMsg := logMessageGet()
    logMsg.Buff = []byte(text)
    logChan <- logMsg
}

func handleWriteLog() {
    for logMsg := range logChan {
        for _, adp  := range adapters {
            adp.Write(logMsg.Buff)
        }
        logMessagePut(logMsg)
    }
    doneChan <- struct{}{}
}
