package zkcli

import (
    "time"
    "fmt"
    "testing"

    "github.com/samuel/go-zookeeper/zk"
)

func TestWatch(t *testing.T) {
    zkConn, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
    if err != nil {
        t.Fatalf("Connect returned error: %+v", err)
    }
    defer zkConn.Close()

    if err := zkConn.Delete("/gozk-test-w1", -1); err != nil && err != zk.ErrNoNode {
        t.Fatalf("Delete returned error: %+v", err)
    }

    testPath, err := zkConn.Create("/gozk-test-w1", []byte{}, 0, zk.WorldACL(zk.PermAll))
    if err != nil {
        t.Fatalf("Create returned: %+v", err)
    }
    stop := make(chan struct{}, 1)
    err = Watch(zkConn, testPath, func(p string, d []byte, e error){
        fmt.Println("w1: ", p, d, e)
    }, stop)
    if err != nil {
        t.Fatal(err)
    }
    stop2 := make(chan struct{}, 1)
    err = Watch(zkConn, "/notexist", func(p string, d []byte, e error){
        fmt.Println("w2: ", p, d, e)
    }, stop2)
    if err != nil {
        t.Fatal(err)
    }

    time.Sleep(time.Second * 60)
    stop <- struct{}{}
    stop2 <- struct{}{}
    time.Sleep(time.Second)
}

func TestWatchChildren(t *testing.T) {
    zkConn, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
    if err != nil {
        t.Fatalf("Connect returned error: %+v", err)
    }
    defer zkConn.Close()

    if err := zkConn.Delete("/gozk-test-wc", -1); err != nil && err != zk.ErrNoNode {
        t.Fatalf("Delete returned error: %+v", err)
    }
    testPath, err := zkConn.Create("/gozk-test-wc", []byte{}, 0, zk.WorldACL(zk.PermAll))
    if err != nil {
        t.Fatalf("Create returned: %+v", err)
    }
    stop := make(chan struct{}, 1)
    err = WatchChildren(zkConn, testPath, func(p string, c []string, e error){
        fmt.Println("wc: ", p, c, e)
    }, stop)
    if err != nil {
        t.Fatal(err)
    }
    time.Sleep(time.Second * 120)
    stop <- struct{}{}
    time.Sleep(time.Second)
}

func TestCreateEphemeral(t *testing.T) {
    zkConn, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
    if err != nil {
        t.Fatalf("Connect returned error: %+v", err)
    }
    defer zkConn.Close()

    err = CreateEphemeral(zkConn, "/go-zktest-ephemeral1", []byte{})
    if err != nil {
        t.Fatal(err)
    }
    //模拟session断开
    zkConn.TmpCloseConn()
    time.Sleep(60 * time.Second)
}