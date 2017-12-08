package main

import (
    "fmt"
    "os"
    "log"
    "flag"
    "path/filepath"
)

var (
    ptrIndex       *int
    ptrPort        *int
    ptrClusterID   *int

    serverType  int

    Log      *log.Logger
)

func init() {
    serverType  = ServerTypeGamesvr

    ptrIndex = flag.Int("i", 0, "server instance index")
    ptrPort = flag.Int("p", 0, "server rpc port")
    ptrClusterID = flag.Int("c", 0, "server clusterid")
}

func main() {
    readFlags()
    setLog()
    Log.Println("hello world!")
    serveRPC(*ptrPort, *ptrClusterID, *ptrIndex)
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

func setLog() {
    svrname := filepath.Base(os.Args[0])
    log.Printf("svrname %v\n", svrname)
    logname := svrname
    if *ptrIndex > 0 {
        logname = fmt.Sprintf("%v%v.log", svrname, *ptrIndex)
    } else {
        logname = fmt.Sprintf("%v.log", svrname)
    }
    log.Printf("svrname: %v, logname: %v\n", svrname, logname)
    f, err := os.OpenFile(logname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("open logfile error : %v", err)
    }

    Log = log.New(f, "", log.Ldate|log.Lmicroseconds|log.Llongfile)
}
