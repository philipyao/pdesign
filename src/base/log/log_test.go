package log

import (
    "testing"
)

/*
func TestLogFile(t *testing.T) {
    config := `{"filename": "flagtest", "maxsize": 102400, "maxbackup": 10}`
    err := AddAdapter(AdapterFile, config)
    if err != nil {
        t.Error("增加adapter失败", err)
    }
    Warn("warn log: %v", 111)
    Info("info log: %v", 222)
    Error("error log: %v", 333)

    SetLevel(LevelStringDebug)
    Debug("debug log: %v", 1)
    Warn("warn log: %v", 111)
    Info("info log: %v", 222)
    Error("error log: %v", 333)
    SetLevel(LevelStringError)
    Debug("debug log: %v", 1)
    Warn("warn log: %v", 111)
    Info("info log: %v", 222)
    Error("error log: %v", 333)

    SetFlags(LogTime | LogLongFile)
    Warn("warn log: %v", 111)
    Info("info log: %v", 222)
    Error("error log: %v", 333)

    Flush()
}
*/

func BenchmarkLogFile(b *testing.B) {
    config := `{"filename": "benchmark", "maxsize": 102400000, "maxbackup": 10}`
    err := AddAdapter(AdapterFile, config)
    if err != nil {
        b.Error("增加adapter失败", err)
    }
    SetFlags(LogDate | LogTime | LogLongFile)
    SetLevel(LevelStringDebug)

    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        for j := 0; j < 100; j++ {
            Debug("debug log: %v", 1)
        }
        for j := 0; j < 100; j++ {
            Info("info log: %v", 222)
        }
        for j := 0; j < 100; j++ {
            Warn("warn log: %v", 111)
        }
        for j := 0; j < 100; j++ {
            Error("error log: %v", 333)
        }
    }
}
