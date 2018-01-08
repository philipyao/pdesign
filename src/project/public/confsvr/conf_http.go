package main

import(
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"

    "base/log"
)

type AdminLoginRsp struct {
    Userinfo    *SUserinfo           `json:"userinfo"`
}

type SUserinfo struct {
    Username        string          `json:"username"`
    Token           string          `json:"token"`
}

type AdminListRsp struct {
    Entries      []*ConfEntry          `json:"entries"`
}
type ConfEntry struct {
    ID              uint        `json:"id"`
    Namespace       string      `json:"namespace"`
    Key             string      `json:"key"`
    Value           string      `json:"value"`
    Updated         uint32      `json:"updated"`
    Created         uint32      `json:"created"`
    Version         int         `json:"version"`
}

type AdminAddReq struct {
    Namespace       string      `json:"namespace"`
    Key             string      `json:"key"`
    Value           string      `json:"value"`
}
type AdminAddRsp struct {
    Entry           *ConfEntry  `json:"entry"`
}

func handle_admin() {
 
    http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            log.Error("parse form error: %v", err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusMethodNotAllowed)
            return
        }

        userName, passwd, veriCode := r.FormValue("username"), r.FormValue("password"), r.FormValue("code")
        log.Debug("ADMIN LOGIN: [%v] [%v] [%v]", userName, passwd, veriCode)
        w.Header().Set("Content-Type", "application/json")
        var loginRsp AdminLoginRsp
        loginRsp.Userinfo.Username = userName
        loginRsp.Userinfo.Token = "HXS04KSSS"
        doWriteJson(w, loginRsp)
    })

    http.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            log.Info("err handle http request, method %v", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        results := AllConfig()
        var rsp AdminListRsp
        for _, r := range results {
            rsp.Entries = append(rsp.Entries, dumpConfEntry(r))
        }
        doWriteJson(w, rsp)
    })

    http.HandleFunc("/api/add", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            log.Info("err handle http request, method %v", r.Method)
            http.Error(w, "inv method", http.StatusMethodNotAllowed)
            return
        }
        reqdata, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Error("read body error %v", err)
            return
        }
        if len(reqdata) == 0 {
            log.Error("no reqdata for /api/add")
            http.Error(w, "no reqdata for /api/add", http.StatusNoContent)
            return
        }
        var req AdminAddReq
        err = json.Unmarshal(reqdata, &req)
        if err != nil {
            log.Error(err.Error())
            http.Error(w, "error parse json reqdata for /api/add", http.StatusBadRequest)
            return
        }
        c, err := addConfig(req.Namespace, req.Key, req.Value)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        var rsp AdminAddRsp
        rsp.Entry = dumpConfEntry(*c)
        doWriteJson(w, rsp)
    })

    http.HandleFunc("/api/update", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            log.Info("err handle http request, method %v", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }
    })
}

func dumpConfEntry(c Config) *ConfEntry {
    return &ConfEntry{
        ID:             c.ID,
        Namespace:      c.Namespace,
        Key:            c.Key,
        Value:          c.Value,
        Updated:        uint32(c.UpdatedAt.Unix()),
        Created:        uint32(c.CreatedAt.Unix()),
        Version:        c.Version,
    }
}

func doWriteJson(w http.ResponseWriter, pkg interface{}) {
    data, err := json.Marshal(pkg)
    if err != nil {
        log.Error("err marshal: err %v, pkg %+v", err, pkg)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func doWriteError(w http.ResponseWriter, errmsg string) {
    w.Write([]byte(errmsg))
}

//=======================================================
func startHttpServer() *http.Server {
    srv := &http.Server{Addr: ":8999"}

    handle_admin()

    go func() {
        if err := srv.ListenAndServe(); err != nil {
            fmt.Printf("Httpserver: ListenAndServe() error: %s\n", err)
        }
    }()

    // returning reference so caller can call Shutdown()
    return srv
}

func serveHttp(done chan struct{}) {
    srv := startHttpServer()
    go func() {
        <- done
        log.Info("stop http listening on %v.", srv.Addr)
        srv.Shutdown(nil)
    }()

    return
}
