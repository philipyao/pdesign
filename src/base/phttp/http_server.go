package phttp

import (
    "fmt"
    "sync"
    "net/http"
)

type Worker struct {
    srv         *http.Server
    logFunc     func(format string, args ...interface{})
}

func New(addr string) *Worker {
    srv := &http.Server{Addr: addr}

    //handle_admin()

    return &Worker{srv: srv}
}

func (w *Worker) Serve(done chan struct{}, wg *sync.WaitGroup) {
    defer wg.Done()

    if w.logFunc != nil {
        w.logFunc("[http] start listening on %v.", w.srv.Addr)
    }

    go func() {
        if err := w.srv.ListenAndServe(); err != nil {
            if err != http.ErrServerClosed {
                if w.logFunc != nil {
                    w.logFunc("[http] ListenAndServe() error %v.", err)
                } else {
                    fmt.Printf("Httpserver: ListenAndServe() error: %s\n", err.Error())
                }
            }
        }
    }()

    <- done
    if w.logFunc != nil {
        w.logFunc("[http] stop listening on %v.", w.srv.Addr)
    }
    w.srv.Shutdown(nil)
}

func (w *Worker) SetLog(l func(format string, args ...interface{})) {
    w.logFunc = l
}
