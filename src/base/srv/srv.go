package srv

import (
    "fmt"
    "os"
    "os/signal"
    "flag"
    "errors"
    "sync"
    "path/filepath"
    "syscall"

    "github.com/philipyao/toolbox/util"
    "github.com/philipyao/phttp"
    "base/prpc"
)

type server struct {
    pName       string

    bInited     bool

    done        chan struct{}
    wg          sync.WaitGroup

    argIndex    *int
    argCluster  *int
    argIP       *string
    argPort     *int

    initFunc        func(chan struct{}) error
    shutdownFunc    func()
    logFunc         func(format string, args ...interface{})

    rpc         *prpc.Worker
    http        *phttp.HTTPWorker
}

type HTTPWorker struct {
    *phttp.HTTPWorker
}

var (
    ptrWanIP          *string
)

var defaultSrv  = &server{done: make(chan struct{})}
func init() {
    defaultSrv.prepare()
}

func (sv *server) addr() string {
    return fmt.Sprintf("%v:%v", *sv.argIP, *sv.argPort)
}

func (sv *server) setRpc(r *prpc.Worker) {
    sv.rpc = r
}

func (sv *server) setHttp(h *phttp.HTTPWorker) {
    sv.http = h
}

func (sv *server) prepare() {
    sv.readArgs()
}

func (sv *server) init() error {
    sv.logFunc("server start...")

    if sv.bInited {
        panic("already inited.")
    }

    err := sv.initFunc(sv.done)
    if err != nil {
        return err
    }
    sv.bInited = true
    sv.logFunc("server init ok.")
    return nil
}

func (sv *server) serve() {
    if !sv.bInited {
        panic("not inited")
    }
    if sv.rpc != nil {
        sv.wg.Add(1)
        go sv.rpc.Serve(sv.done, &sv.wg)
    }
    if sv.http != nil {
        sv.wg.Add(1)
        go sv.http.Serve(sv.done, &sv.wg)
    }
    sv.writePid()

    sv.wg.Add(1)
    go sv.listenInterupt()

    sv.wg.Wait()

    sv.shutdownFunc()
    sv.removePid()
}

//====================================

func (sv *server) readArgs() {
    sv.argCluster = flag.Int("c", 0, "server clusterid")
    sv.argIndex = flag.Int("i", 0, "server instance index")
    sv.argIP = flag.String("l", "0.0.0.0", "server local ip")
    sv.argPort = flag.Int("p", 0, "server rpc port")
    ptrWanIP = flag.String("w", "0.0.0.0", "server wan ip")

    flag.Parse()
    if *sv.argPort <= 0 {
        panic("no server port specified!")
    }
    if *sv.argCluster <= 0 {
        panic("no server cluster id specified")
    }
}

func (sv *server) listenInterupt() {
    defer sv.wg.Done()

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    <-sigs
    sv.shutdown()
}

func (sv *server) shutdown() {
    sv.logFunc("graceful shutdown...")
    close(sv.done)
}

func (sv *server) writePid() {
    pName := sv.processName()
    pidFile := util.GenPidFilePath(pName)
    util.WritePidToFile(pidFile, os.Getpid())
}

func (sv *server) removePid() {
    pName := sv.processName()
    pidFile := util.GenPidFilePath(pName)
    util.DeletePidFile(pidFile)
}

func (sv *server) processName() string {
    if sv.pName != "" {
        return sv.pName
    }
    sv.pName = filepath.Base(os.Args[0])
    if *sv.argIndex > 0 {
        sv.pName = fmt.Sprintf("%v%v", sv.pName, *sv.argIndex)
    }
    return sv.pName
}


//=====================================================

//必须实现，server基础接口
func HandleBase(onInit func(chan struct{}) error,
                onShutdown func()) error {
    if onInit == nil {
        return errors.New("nil onInit.")
    }
    if onShutdown == nil {
        return errors.New("nil onShutdown.")
    }
    defaultSrv.initFunc = onInit
    defaultSrv.shutdownFunc = onShutdown
    if defaultSrv.logFunc == nil {
        defaultSrv.logFunc = defaultLogFunc()
    }
    return defaultSrv.init()
}

// 可选，注册rpc服务
func HandleRpc(rpcName string, rpcWorker interface{}) error {
    rpcW := prpc.New(defaultSrv.addr(), rpcName, rpcWorker)
    if rpcW == nil {
        return errors.New(prpc.ErrMsg())
    }
    rpcW.SetLog(defaultSrv.logFunc)
    defaultSrv.setRpc(rpcW)

    return nil
}

// 可选，注册http服务
func HandleHttp(addr string) (*HTTPWorker, error) {
    httpW := phttp.New(addr)
    if httpW == nil {
        return nil, errors.New("init http error")
    }
    httpW.SetLog(defaultSrv.logFunc)
    defaultSrv.setHttp(httpW)

    return &HTTPWorker{httpW}, nil
}

//可选，自定义log输出
func SetLogger(l func(int, string, ...interface{})) {
    defaultSrv.logFunc = customLogFunc(l)
}

// server运行入口函数
func Run() {
    defaultSrv.serve()
}

// 获取server进程名字
func ProcessName() string {
    return defaultSrv.processName()
}

