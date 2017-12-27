package zkcli

import (
    "fmt"
    "strings"
    "time"

    "github.com/samuel/go-zookeeper/zk"
)

//回调：监视节点变化
type FuncWatchCallback func(path string, data []byte, err error)
//回调：监视子节点改变
type FuncWatchChildrenCallback func(path string, children []string, err error)

const (
    //zk session的超时时间，在此时间内，session可以自动保活
    //如果在此时间内，与特定server的连接断开，自动尝试重连其他服务器
    //重连成功之后，临时节点和watch依然有效
    DefaultConnectTimeout       = 5
)

type Conn struct {
    conn *zk.Conn
}
func (c *Conn) Conn() *zk.Conn {
    return c.conn
}
func (c *Conn) SetConn(conn *zk.Conn) {
    c.conn = conn
}

func (c *Conn) Write(path string, data []byte) error {
    exist, stat, err := c.conn.Exists(path)
    if err != nil {
        return err
    }
    if exist {
        fmt.Printf("set %v\n", path)
        _, err = c.conn.Set(path, data, stat.Version)
    } else {
        //不存在则创建
        // 永久节点
        fmt.Printf("create %v\n", path)
        flags := int32(0)
        _, err = doCreate(c.conn, path, data, flags)
    }
    return err
}

func Connect(zkAddr string) (*Conn, error) {
    if len(zkAddr) == 0 {
        return nil, fmt.Errorf("empty zkAddr")
    }
    zks := strings.Split(zkAddr, ",")
    conn, _, err := zk.Connect(zks, time.Second * DefaultConnectTimeout)
    if err != nil {
        return nil, fmt.Errorf("err connect to zk<%v>: %v", zkAddr, err)
    }

    c := new(Conn)
    c.SetConn(conn)
    return c, nil
}

func CreateEphemeral(zkConn *Conn, path string, data []byte) error {
    exist, _, err := zkConn.Conn().Exists(path)
    if err != nil {
        return err
    }
    if exist {
        return nil
    }
    // 临时节点
    _, err = doCreate(zkConn.Conn(), path, data, int32(zk.FlagEphemeral))
    return err
}

func CreateSequence(zkConn *Conn, path string, data []byte) (string, error) {
    flags := int32(zk.FlagSequence | zk.FlagEphemeral)
    return doCreate(zkConn.Conn(), path, data, flags)
}

func Watch(zkConn *Conn, path string, cb FuncWatchCallback, stopCh chan struct{}) error {
    _, ch, err := getW(zkConn.Conn(), path)
    if err != nil {
        return err
    }
    go func() {
        var data []byte
        for {
            select {
            case <-stopCh:
                return
            case ev := <-ch:
                if ev.Err != nil {
                    //错误回调
                    cb(path, nil, ev.Err)
                    return
                }
                if ev.Path != path {
                    cb(path, nil, fmt.Errorf("mismatched path %v %v", ev.Path, path))
                    return
                }
            }
            // 获取变化后的节点数据
            // 并更新watcher（zookeeper的watcher是一次性的）
            data, ch, err = getW(zkConn.Conn(), path)
            if err != nil {
                //错误回调
                cb(path, nil, err)
                return
            }
            //数据回调
            cb(path, data, nil)
        }
    }()

    return nil
}

func WatchChildren(zkConn *Conn, path string, cb FuncWatchChildrenCallback, stopCh chan struct{}) error {
    _, ch, err := childrenW(zkConn.Conn(), path)
    if err != nil {
        return err
    }
    go func() {
        for {
            select {
            case <-stopCh:
                return
            case ev := <-ch:
                if ev.Err != nil {
                    //错误回调
                    cb(path, nil, ev.Err)
                    return
                }
                if ev.Path != path {
                    cb(path, nil, fmt.Errorf("mismatched path %v %v", ev.Path, path))
                    return
                }
            }
            // 获取变化后的节点数据
            // 并更新watcher（zookeeper的watcher是一次性的）
            var children []string
            children, ch, err = childrenW(zkConn.Conn(), path)
            if err != nil {
                //错误回调
                cb(path, nil, err)
                return
            }
            //数据回调
            cb(path, children, nil)
        }
    }()

    return nil
}

////////////////////////////////////////////////////////////////////////////

func getW(zkConn *zk.Conn, path string) ([]byte, <-chan zk.Event, error) {
    data, _, ch, err := zkConn.GetW(path)
    return data, ch, err
}
func childrenW(zkConn *zk.Conn, path string) ([]string, <-chan zk.Event, error) {
    children, _, ch, err := zkConn.ChildrenW(path)
    return children, ch, err
}

func doCreate(zkConn *zk.Conn, path string, data []byte, flags int32) (string, error) {
    //TODO 权限控制
    acl := zk.WorldACL(zk.PermAll)
    path, err := zkConn.Create(path, data, flags, acl)
    if err != nil {
        return "", err
    }

    return path, nil
}




