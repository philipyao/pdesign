package main

import (
    "fmt"
    "os"
    "os/signal"
    "log"
    "time"
    "flag"
    "sync"
    "path/filepath"
    "syscall"

    "base/def"
    "base/util"
)

var (
    ptrIndex       *int
    ptrPort        *int
    ptrClusterID   *int
    ptrIP             *string
    ptrWanIP          *string

    serverType  int

    Log         *log.Logger

    done        chan struct{}
    wg          *sync.WaitGroup
)

func init() {
    serverType  = def.ServerTypeConfsvr

    ptrIndex = flag.Int("i", 0, "server instance index")
    ptrPort = flag.Int("p", 0, "server rpc port")
    ptrClusterID = flag.Int("c", 0, "server clusterid")
    ptrIP = flag.String("l", "0.0.0.0", "server local ip")
    ptrWanIP = flag.String("w", "0.0.0.0", "server wan ip")

    done = make(chan struct{})
    wg = &sync.WaitGroup{}
}

func shutdown() {
    Log.Println("shutdown server")
    close(done)
}

func main() {
    log.Println()

    readFlags()
    setLog()
    Log.Println("hello server!")

    handleSignal()

    serveRPC(done, *ptrPort, *ptrClusterID, *ptrIndex)

    writePid()

    wg.Wait()

    removePid()

    Log.Println("server exit.")
}

func readFlags() {
    flag.Parse()
    if *ptrPort <= 0 {
        log.Fatalf("invalid port: %v", *ptrPort)
    }
    if *ptrClusterID <= 0 {
        log.Fatalf("invalid clusterid: %v", *ptrClusterID)
    }
    log.Printf("port %v, clusterid %v\n", *ptrPort, *ptrClusterID)
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
    svrname := filepath.Base(os.Args[0])
    log.Printf("svrname %v\n", svrname)
    _, month, day := time.Now().Date()
    logname := processName() + fmt.Sprintf(".%v%v", int(month), day) + ".log"
    logname = filepath.Join("log", logname)
    log.Printf("svrname: %v, logname: %v\n", svrname, logname)
    f, err := os.OpenFile(logname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("open logfile error : %v", err)
    }

    Log = log.New(f, "", log.Ldate|log.Lmicroseconds|log.Llongfile)
}

func writePid() {
    pName := processName()
    pidFile := util.GenPidFilePath(pName)
    util.WritePidToFile(pidFile, os.Getpid())
    Log.Printf("pidfile %v, pid %v", pidFile, os.Getpid())
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
        Log.Printf("receive sig %v\n", <-sigs)
        shutdown()
    }()
    return
}
