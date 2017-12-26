package log

import (
    "sync"
)
type logMessage struct {
    Buff       []byte
}

var (
    logPool    sync.Pool
)

func init() {
    logPool.New = func() interface{} {
        return new(logMessage)
    }
}

func logMessageGet() *logMessage {
    return logPool.Get().(*logMessage)
}

func logMessagePut(l *logMessage) {
    logPool.Put(l)
}

