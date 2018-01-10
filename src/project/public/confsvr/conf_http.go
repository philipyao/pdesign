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

type AdminUpdateReq struct {
	Updates			[]*UpdateEntry	`json:"updates"`
}
type UpdateEntry struct {
    ID              uint        `json:"id"`
    Value           string      `json:"value"`
	Comment			string		`json:"comment"`
	Author			string		`json:"author"`
}

type AdminUpdateRsp struct {
    Entries      []*ConfEntry          `json:"entries"`
	Errmsgs	     []string			   `json:"errmsgs"`
}

var httpHandler = map[string]func(w http.ResponseWriter, r *http.Request){

    "/api/login": func(w http.ResponseWriter, r *http.Request) {
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
		loginRsp.Userinfo = &SUserinfo{
        	Username: userName,
        	Token: "HXS04KSSS",
		}
        doWriteJson(w, loginRsp)
    },

    "/api/list": func(w http.ResponseWriter, r *http.Request) {
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
    },

    "/api/add": func(w http.ResponseWriter, r *http.Request) {
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
    },

    "/api/update": func(w http.ResponseWriter, r *http.Request) {
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
            log.Error("no reqdata for /api/update")
            http.Error(w, "no reqdata for /api/update", http.StatusNoContent)
            return
        }
        var req AdminUpdateReq
        err = json.Unmarshal(reqdata, &req)
        if err != nil {
            log.Error(err.Error())
            http.Error(w, "error parse json reqdata for /api/update", http.StatusBadRequest)
            return
        }
		if len(req.Updates) == 0 {
            http.Error(w, "no updates provided", http.StatusBadRequest)
			return
		}
		log.Debug("updates: %+v", req.Updates)
        var rsp AdminUpdateRsp
		var errMsgs []string
		for _, upd := range req.Updates {
			err = updateConfig(upd.ID, upd.Value)
			log.Debug("try update: %v %v, err %v", upd.ID, upd.Value, err)
			if err != nil {
				errMsg := fmt.Sprintf("config<id:%v> update error: %v; ", upd.ID, err.Error())	
				errMsgs = append(errMsgs, errMsg)
				continue
			}
			c := configByID(upd.ID)
			log.Debug("updated: %+v", c)
			rsp.Entries = append(rsp.Entries, dumpConfEntry(*c))
		}
		rsp.Errmsgs = errMsgs
        doWriteJson(w, rsp)
    },
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
