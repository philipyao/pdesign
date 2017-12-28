package main

import(
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"

    "base/log"
)

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
            rsp.Entries = append(rsp.Entries, &ConfEntry{
                ID:             r.ID,
                Namespace:      r.Namespace,
                Key:            r.Key,
                Value:          r.Value,
                Updated:        uint32(r.UpdatedAt.Unix()),
                Created:        uint32(r.CreatedAt.Unix()),
                Version:        r.Version,
            })
        }
        doWriteJson(w, rsp)
    })

    http.HandleFunc("/api/add", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            log.Info("err handle http request, method %v", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }
        reqdata, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Error("read body error %v", err)
            return
        }
        if len(reqdata) == 0 {
            log.Error("no reqdata for /api/add")
            http.Error(w, "no reqdata for /api/add", http.StatusBadRequest)
            return
        }
        var req AdminAddReq
        err = json.Unmarshal(reqdata, &req)
        if err != nil {
            log.Error(err.Error())
            http.Error(w, "error parse json reqdata for /api/add", http.StatusBadRequest)
            return
        }
    })

    http.HandleFunc("/api/update", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            log.Info("err handle http request, method %v", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }
    })
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
    srv := &http.Server{Addr: ":8080"}

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
        log.Info("stop http listening.")
        srv.Shutdown(nil)
    }()

    return
}
