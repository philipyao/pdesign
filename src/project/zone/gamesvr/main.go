package main

import (
    "fmt"
    "os"
    "os/signal"
    //"time"
    "flag"
    "sync"
    "path/filepath"
    "syscall"

    "base/util"
    "base/log"
    "project/share"

)

var (
    ptrIndex       *int
    ptrPort        *int
    ptrClusterID   *int
    ptrIP             *string
    ptrWanIP          *string

    serverType  int

    done        chan struct{}
    wg          *sync.WaitGroup
)

func init() {
    serverType  = share.ServerTypeGamesvr

    ptrIndex = flag.Int("i", 0, "server instance index")
    ptrPort = flag.Int("p", 0, "server rpc port")
    ptrClusterID = flag.Int("c", 0, "server clusterid")
    ptrIP = flag.String("l", "0.0.0.0", "server local ip")
    ptrWanIP = flag.String("w", "0.0.0.0", "server wan ip")

    done = make(chan struct{})
    wg = &sync.WaitGroup{}
}

func shutdown() {
    log.Info("shutdown gamesvr")
    close(done)

    log.Flush()
}

func main() {
    fmt.Println("\n")

    readFlags()
    setLog()

    handleSignal()
    TryGetGamesvrConfig()
    RegisterConfDef()
    HandleConfChange("log_level", log.LevelStringInfo, log.LevelStringDebug)
    serveRPC(done, *ptrPort, *ptrClusterID, *ptrIndex)

    writePid()

    wg.Wait()

    removePid()
}

func readFlags() {
    flag.Parse()
    if *ptrPort <= 0 {
        fmt.Printf("invalid port: %v", *ptrPort)
        panic("port")
    }
    if *ptrClusterID <= 0 {
        fmt.Printf("invalid clusterid: %v", *ptrClusterID)
        panic("clusterid")
    }
    fmt.Printf("port %v, clusterid %v\n", *ptrPort, *ptrClusterID)
}

func processName() string {
    svrname := filepath.Base(os.Args[0])
    pName := svrname
    if *ptrIndex > 0 {
        pName = fmt.Sprintf("%v%v", pName, *ptrIndex)
    }
    return pName
}

func setLog() {
    config := `{"filename": "%v", "maxsize": 102400, "maxbackup": 10}`
    wd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    logName := filepath.Join(wd, "log", processName())
    config = fmt.Sprintf(config, logName)
    fmt.Printf("log config: %+v\n", config)
    err = log.AddAdapter(log.AdapterFile, config)
    if err != nil {
        panic(err)
    }
    log.SetLevel(log.LevelStringDebug)
    log.SetFlags(log.LogDate | log.LogTime | log.LogMicroTime | log.LogLongFile)
}

func writePid() {
    pName := processName()
    pidFile := util.GenPidFilePath(pName)
    util.WritePidToFile(pidFile, os.Getpid())
    log.Info("pidfile %v, pid %v", pidFile, os.Getpid())
}

func removePid() {
    pName := processName()
    pidFile := util.GenPidFilePath(pName)
    util.DeletePidFile(pidFile)
}

func handleSignal() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGTERM)
    go func() {
        log.Info("receive sig %v\n", <-sigs)
        shutdown()
    }()
    return
}
