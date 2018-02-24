package main

import(
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"

    "base/log"
)

const (
    CookieName      = "sessid"
    KeyUserName     = "username"
)

type AdminError struct {
    Errmsg      string      `json:"errmsg"`
}

type AdminLoginRsp struct {
    AdminError
    Userinfo    *SUserinfo           `json:"userinfo"`
}

type SUserinfo struct {
    Username        string          `json:"username"`
    Token           string          `json:"token"`
}

type AdminListRsp struct {
    AdminError
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
    AdminError
    Entry           *ConfEntry  `json:"entry"`
}

type AdminUpdateReq struct {
    Adds            []*AddEntry     `json:"adds"`
	Updates			[]*UpdateEntry	`json:"updates"`
    Name            string          `json:"name"`
    Comment         string          `json:"comment"`
    Author          string          `json:"author"`
}
type AddEntry struct {
    Namespace       string      `json:"namespace"`
    Key             string      `json:"key"`
    Value           string      `json:"value"`
}
type UpdateEntry struct {
	ID              uint        `json:"id"`
    Value           string      `json:"value"`
    Version         int         `json:"version"`    //保证客户端和服务器当前的version一致
}

type AdminUpdateRsp struct {
    AdminError
    Entries      []*ConfEntry          `json:"entries"`
	Failed	     []string			   `json:"errmsgs"`
}

type AdminListUserRsp struct {
    AdminError
    Entries      []*UserEntry          `json:"entries"`
}
type UserEntry struct {
    Username        string              `json:"username"`
    Enabled         uint                `json:"enabled"`
    CreatedAt       int                 `json:"created_at"`
}

type AdminCreateUserReq struct {
    Username        string              `json:"username"`
    EncPasswd       string              `json:"enc_passwd""`    //客户端初次加密后的密码
}
type AdminCreateUserRsp struct {
    AdminError
    Entry        *UserEntry          `json:"entry"`
}

type AdminChangeUserReq struct {
    Username        string              `json:"username"`
    Enable          uint32              `json:"enable"` //启用或者禁用
}
type AdminChangeUserRsp struct {
    AdminError
}

var (
    smgr *SessionMgr = NewManager(3600)
)

var httpHandler = map[string]func(w http.ResponseWriter, r *http.Request){

    "/api/login": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, err.Error())
            return
        }

        err = r.ParseForm()
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
        pass, err := verifyUser(userName, passwd)
        if err != nil {
            errcode := http.StatusInternalServerError
            http.Error(w, http.StatusText(errcode), errcode)
            return
        }
        if !pass {
            errcode := http.StatusUnauthorized
            http.Error(w, http.StatusText(errcode), errcode)
            return
        }
        sess.Set(KeyUserName, userName)

        w.Header().Set("Content-Type", "application/json")
        var loginRsp AdminLoginRsp
		loginRsp.Userinfo = &SUserinfo{
        	Username: userName,
		}
        doWriteJson(w, loginRsp)
    },

    "/api/list": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, err.Error())
            return
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //需要重新登录
            doWriteError(w, "need login")
            return
        }

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

    "/api/change": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, err.Error())
            return
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //需要重新登录
            doWriteError(w, "need login")
            return
        }

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
		log.Debug("update req: %+v", req)
        var rsp AdminUpdateRsp
        defer func() {
            log.Debug("rsp %+v", rsp)
            doWriteJson(w, rsp)
        }()

        //开始参数校验
		if len(req.Updates) == 0 && len(req.Adds) == 0 {
            rsp.Errmsg = "no adds or updates provided"
			return
		}
        if req.Name == "" || req.Author == "" {
            rsp.Errmsg = "no name or author provided, pls check"
            return
        }

		var failed []string
        var changes []*OpChange
		for _, upd := range req.Updates {
            c := configByID(upd.ID)

            var change OpChange
            change.Namespace = c.Namespace
            change.Key = c.Key
            change.OldValue = c.Value
            change.Value = upd.Value

			err = updateConfig(upd.ID, upd.Value, upd.Version)
			log.Debug("try update: %v %v, err %v", upd.ID, upd.Value, err)
			if err != nil {
				errMsg := fmt.Sprintf("config<id:%v> update error: %v; ", upd.ID, err.Error())
                failed = append(failed, errMsg)
				continue
			}

			log.Debug("updated ok: %+v", c)
			rsp.Entries = append(rsp.Entries, dumpConfEntry(*c))
            changes = append(changes, &change)
		}
        for _, add := range req.Adds {
            var change OpChange
            change.Namespace = add.Namespace
            change.Key = add.Key
            change.OldValue = ""
            change.Value = add.Value

            c, err := addConfig(add.Namespace, add.Key, add.Value)
            if err != nil {
                errMsg := fmt.Sprintf("config<%+v> add error: %v; ", add, err.Error())
                failed = append(failed, errMsg)
                continue
            }
            log.Debug("added ok: %+v", c)
            rsp.Entries = append(rsp.Entries, dumpConfEntry(*c))
            changes = append(changes, &change)
        }
        if len(changes) > 0 {
            //TODO 失败的要不要记录？
            addOplog(req.Name, req.Comment, req.Author, changes)
        }
		rsp.Failed = failed
    },

    "/api/user/list": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, err.Error())
            return
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //需要重新登录
            doWriteError(w, "need login")
            return
        }

        if r.Method != "GET" {
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        cando, err := CheckUserPrivilege(username.(string))
        if err != nil {
            errcode := http.StatusInternalServerError
            http.Error(w, http.StatusText(errcode), errcode)
            return
        }
        if !cando {
            errcode := http.StatusUnauthorized
            http.Error(w, http.StatusText(errcode), errcode)
            return
        }
        var rsp AdminListUserRsp
        users, err := ListUser()
        if err != nil {
            rsp.Errmsg = err.Error()
        } else {
            for _, u := range users {
                rsp.Entries = append(rsp.Entries, dumpUserEntry(*u))
            }
        }
        doWriteJson(w, rsp)
    },

    "/api/user/create": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, err.Error())
            return
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //需要重新登录
            doWriteError(w, "need login")
            return
        }

        if r.Method != "POST" {
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        cando, err := CheckUserPrivilege(username.(string))
        if err != nil {
            errcode := http.StatusInternalServerError
            http.Error(w, http.StatusText(errcode), errcode)
            return
        }
        if !cando {
            errcode := http.StatusUnauthorized
            http.Error(w, http.StatusText(errcode), errcode)
            return
        }

        reqdata, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Error("read body error %v", err)
            return
        }
        if len(reqdata) == 0 {
            log.Error("no reqdata for /api/user/create")
            http.Error(w, "no reqdata", http.StatusNoContent)
            return
        }
        var req AdminCreateUserReq
        err = json.Unmarshal(reqdata, &req)
        if err != nil {
            log.Error(err.Error())
            http.Error(w, "error parse json reqdata", http.StatusBadRequest)
            return
        }
        log.Debug("user create req: %+v", req)
        if req.Username == "" || req.EncPasswd == "" {
            http.Error(w, "invalid req data", http.StatusBadRequest)
            return
        }

        var rsp AdminCreateUserRsp
        user, err := CreateUser(req.Username, req.EncPasswd)
        if err != nil {
            rsp.Errmsg = err.Error()
        } else {
            rsp.Entry = dumpUserEntry(*user)
        }
        doWriteJson(w, rsp)
    },

    "/api/user/change": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, err.Error())
            return
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //需要重新登录
            doWriteError(w, "need login")
            return
        }

        if r.Method != "POST" {
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        cando, err := CheckUserPrivilege(username.(string))
        if err != nil {
            errcode := http.StatusInternalServerError
            http.Error(w, http.StatusText(errcode), errcode)
            return
        }
        if !cando {
            errcode := http.StatusUnauthorized
            http.Error(w, http.StatusText(errcode), errcode)
            return
        }

        reqdata, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Error("read body error %v", err)
            return
        }
        if len(reqdata) == 0 {
            log.Error("no reqdata for /api/user/change")
            http.Error(w, "no reqdata", http.StatusNoContent)
            return
        }
        var req AdminChangeUserReq
        err = json.Unmarshal(reqdata, &req)
        if err != nil {
            log.Error(err.Error())
            http.Error(w, "error parse json reqdata", http.StatusBadRequest)
            return
        }
        log.Debug("user create req: %+v", req)
        if req.Username == "" {
            http.Error(w, "invalid req data", http.StatusBadRequest)
            return
        }

        var rsp AdminChangeUserRsp
        if req.Enable == 0 {
             err = disableUser(req.Username)
        } else {
            err = enableUser(req.Username)
        }
        if err != nil {
            rsp.Errmsg = err.Error()
        }
        doWriteJson(w, rsp)
    },
}

//===================================================================
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

func dumpUserEntry (u User) *UserEntry {
    return &UserEntry{
        Username: u.Username,
        Enabled: u.Enabled,
        CreatedAt: int(u.CreatedAt.Unix()),
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
    var rsp AdminError
    rsp.Errmsg = errmsg
    doWriteJson(w, &rsp)
}
