package prpc

import (
    "fmt"
    "time"
    "net"
    "sync"
    "net/rpc"
)

type Worker struct {
    listener    *net.TCPListener
    logFunc     func(format string, args ...interface{})
}

var (
    errMsg  string
)

func ErrMsg() string {
    return errMsg
}

func New(addr, rpcName string, rpcWorker interface{}) *Worker {
    var err error
    err = rpc.RegisterName(rpcName, rpcWorker)
    if err != nil {
        errMsg = fmt.Sprintf("[rpc] RegisterName() error: %v", err)
        return nil
    }
    laddr, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil {
        errMsg = fmt.Sprintf("[rpc] ResolveTCPAddr() error: addr %v, err %v", addr, err)
        return nil
    }

    l, err := net.ListenTCP("tcp", laddr)
    if err != nil {
        errMsg = fmt.Sprintf("[rpc ] rpc listen on %v, %v", laddr, err)
        return nil
    }

    return &Worker{ listener: l }
}

func (w *Worker) Serve(done chan struct{}, wg *sync.WaitGroup) {
    if w.logFunc != nil {
        w.logFunc("[rpc] start listening on %v...", w.listener.Addr())
    }
    w.doServe(w.listener, done, wg)
}

func (w *Worker) SetLog(l func(format string, args ...interface{})) {
    w.logFunc = l
}

func (w *Worker) doServe(listener *net.TCPListener, done chan struct{}, wg *sync.WaitGroup) {
    defer wg.Done()
    defer listener.Close()

    for {
        select {
        case <-done:
            if w.logFunc != nil {
                w.logFunc("[rpc] stop listening on %v...", listener.Addr())
            }
            return
        default:
        }
        listener.SetDeadline(time.Now().Add(1e9))
        conn, err := listener.AcceptTCP()
        if err != nil {
            if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                continue
            }
            if w.logFunc != nil {
                w.logFunc("[rpc] Error: accept connection, %v", err.Error())
            }
        }
        go rpc.ServeConn(conn)
    }
}
