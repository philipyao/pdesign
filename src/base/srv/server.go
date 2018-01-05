package srv

import (
    "fmt"
    "os"
    "os/signal"
    "flag"
    "log"
    "errors"
    "sync"
    "path/filepath"
    "syscall"

    "base/util"
    "base/prpc"
)

type Func func()
type server struct {
    done        chan struct{}
    wg          sync.WaitGroup

    argIndex    *int
    argCluster  *int
    argIP       *string
    argPort     *int

    initFunc    Func

    rpc         *prpc.Worker
}

var (
    ptrWanIP          *string
)

func init() {
    log.SetFlags(log.LstdFlags)
}

var defaultSrv  = &server{done: make(chan struct{})}

func (sv *server) Addr() string {
    return fmt.Sprintf("%v:%v", *sv.argIP, *sv.argPort)
}

func (sv *server) SetRpc(r *prpc.Worker) {
    sv.rpc = r
}

func (sv *server) Serve() {
    log.Println("server start...")

    sv.initArgs()
    sv.readArgs()

    if sv.initFunc != nil {
        sv.initFunc()
    }

    if sv.rpc != nil {
        sv.wg.Add(1)
        sv.rpc.Serve(sv.done, &sv.wg)
    }
    //serveHttp(done)

    sv.writePid()

    sv.wg.Add(1)
    go sv.listenInterupt()

    sv.wg.Wait()

    //finiCore()
    sv.removePid()

    log.Println("server stop.")
}

func (sv *server) initArgs() {
    sv.argCluster = flag.Int("c", 0, "server clusterid")
    sv.argIndex = flag.Int("i", 0, "server instance index")
    sv.argIP = flag.String("l", "0.0.0.0", "server local ip")
    sv.argPort = flag.Int("p", 0, "server rpc port")

    ptrWanIP = flag.String("w", "0.0.0.0", "server wan ip")
}

func (sv *server) readArgs() {
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
    log.Printf("graceful shutdown...\n")
    close(sv.done)
}

func (sv *server) writePid() {
    pName := sv.processName()
    pidFile := util.GenPidFilePath(pName)
    util.WritePidToFile(pidFile, os.Getpid())
    log.Printf("writePid: pidfile %v, pid %v\n", pidFile, os.Getpid())
}

func (sv *server) removePid() {
    pName := sv.processName()
    pidFile := util.GenPidFilePath(pName)
    util.DeletePidFile(pidFile)
    log.Printf("removePid: pidfile %v\n", pidFile)
}

func (sv *server) processName() string {
    svrname := filepath.Base(os.Args[0])
    pName := svrname
    if *sv.argIndex > 0 {
        pName = fmt.Sprintf("%v%v", pName, *sv.argIndex)
    }
    return pName
}

//=====================================================

// 可选，注册rpc服务
func HandleRpc(rpcName string, rpcWorker interface{}) error {
    rpcW := prpc.New(defaultSrv.Addr(), rpcName, rpcWorker)
    if rpcW == nil {
        return errors.New(prpc.ErrMsg())
    }
    defaultSrv.SetRpc(rpcW)
    return nil
}

// 可选，注册http服务
func HandleHttp() error {
    return nil
}

// server运行入口函数
func Run() {
    defaultSrv.Serve()
}


