package adapter

type Adapter interface {
    Write(b []byte) error
    Close()
}

//日志大小
type ByteSize int64
const (
    _ ByteSize      = 1 << (iota * 10)
    KB
    MB
    GB
    TB
    PB
)

//日志选项
type Options struct {
    MaxSize         ByteSize
    MaxBackup       int
}
