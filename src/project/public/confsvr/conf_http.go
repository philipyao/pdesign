package main

import(
    "fmt"
    "errors"
    "io/ioutil"
    "net/http"
    "encoding/json"

    "base/log"
)

const (
    CookieName      = "sessid"
    KeyUserName     = "username"
)

//错误码
const (
    ErrOK                   = 0

    ErrSystem               = -1                //系统错误
    ErrParamParseForm       = -20001            //非法的form值
    ErrParamParseBody       = -20002            //读取body失败
    ErrParamInvalid         = -20003            //参数非法
    ErrMethod               = -20101            //非法的method

    ErrUnauthorized         = 40001             //未授权
    ErrAccountDisabled      = 40002             //账号被禁用
    ErrAccountPasswd        = 40003             //用户名或密码错误
    ErrSessionExpired       = 40004             //服务器session过期，需要重新登录
    ErrAccountExist         = 40005             //用户已存在
)

type AdminError struct {
    Errcode     int         `json:"errcode"`
    Errmsg      string      `json:"errmsg"`
}

type AdminLoginRsp struct {
    AdminError
    Userinfo    *SUserinfo           `json:"userinfo"`
}

type SUserinfo struct {
    Username        string          `json:"username"`
    IsSuper         uint            `json:"is_super"` //是否超级用户
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

    //登录
    "/api/login": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, ErrSystem, err.Error())
            return
        }

        err = r.ParseForm()
        if err != nil {
            log.Error("parse form error: %v", err)
            doWriteError(w, ErrParamParseForm, err.Error())
            return
        }
        if r.Method != "POST" {
            log.Error("handle http request, inv method %v", r.Method)
            doWriteError(w, ErrMethod, "")
            return
        }

        userName, passwd:= r.FormValue("username"), r.FormValue("password")
        log.Debug("admin user LOGIN: %v@%v", userName, passwd)
        if userName == "" || passwd == "" {
            doWriteError(w, ErrParamInvalid, "")
            return
        }

        pass, err := verifyUser(userName, passwd)
        if err != nil {
            if err == ErrUserDisabled {
                doWriteError(w, ErrAccountDisabled, "")
                return
            }
            doWriteError(w, ErrSystem, err.Error())
            return
        }
        if !pass {
            doWriteError(w, ErrAccountPasswd, "")
            return
        }

        //拉取user信息
        user, err := QueryUser(userName)
        if err != nil {
            doWriteError(w, ErrSystem, err.Error())
            return
        }
        if user == nil {
            doWriteError(w, ErrSystem, "user not exist")
            return
        }

        //关联session和user
        sess.Set(KeyUserName, userName)

        w.Header().Set("Content-Type", "application/json")
        var loginRsp AdminLoginRsp
		loginRsp.Userinfo = &SUserinfo{
        	Username: userName,
            IsSuper: user.IsSuper,
		}
        doWriteJson(w, loginRsp)
    },

    //列出所有配置
    "/api/list": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, ErrSystem, err.Error())
            return
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //session过期，需要重新登录
            doWriteError(w, ErrSessionExpired, "")
            return
        }

        if r.Method != "POST" {
            log.Error("handle http request, inv method %v", r.Method)
            doWriteError(w, ErrMethod, "")
        }

        results := AllConfig()
        var rsp AdminListRsp
        for _, r := range results {
            rsp.Entries = append(rsp.Entries, dumpConfEntry(r))
        }
        doWriteJson(w, rsp)
    },

    //修改配置
    "/api/change": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, ErrSystem, err.Error())
            return
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //session过期，需要重新登录
            doWriteError(w, ErrSessionExpired, "")
            return
        }

        if r.Method != "POST" {
            log.Error("handle http request, inv method %v", r.Method)
            doWriteError(w, ErrMethod, "")
            return
        }
        var req AdminUpdateReq
        err = readBodyData(r, &req)
        if err != nil {
            doWriteError(w, ErrParamParseBody, err.Error())
            return
        }
		log.Debug("update req: %+v", req)


        //开始参数校验
		if len(req.Updates) == 0 && len(req.Adds) == 0 {
            doWriteError(w, ErrParamInvalid, "no adds or updates provided")
			return
		}
        if req.Name == "" || req.Author == "" {
            doWriteError(w, ErrParamInvalid, "no name or author provided, pls check")
            return
        }

        var rsp AdminUpdateRsp
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
        doWriteJson(w, rsp)
    },

    "/api/user/list": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, ErrSystem, err.Error())
            return
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //session过期，需要重新登录
            doWriteError(w, ErrSessionExpired, "")
            return
        }
        if r.Method != "GET" {
            log.Error("handle http request, inv method %v", r.Method)
            doWriteError(w, ErrMethod, "")
            return
        }

        //校验用户权限
        cando, err := CheckUserPrivilege(username.(string))
        if err != nil {
            doWriteError(w, ErrSystem, err.Error())
            return
        }
        if !cando {
            doWriteError(w, ErrUnauthorized, "")
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
            doWriteError(w, ErrSystem, err.Error())
            return
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //session过期，需要重新登录
            doWriteError(w, ErrSessionExpired, "")
            return
        }
        if r.Method != "POST" {
            log.Error("handle http request, inv method %v", r.Method)
            doWriteError(w, ErrMethod, "")
            return
        }

        cando, err := CheckUserPrivilege(username.(string))
        if err != nil {
            doWriteError(w, ErrSystem, err.Error())
            return
        }
        if !cando {
            doWriteError(w, ErrUnauthorized, "")
            return
        }
        var req AdminCreateUserReq
        err = readBodyData(r, &req)
        if err != nil {
            doWriteError(w, ErrParamParseBody, err.Error())
            return
        }
        log.Debug("user create req: %+v, passwd len %v", req, len(req.EncPasswd))
        if req.Username == "" || len(req.EncPasswd) != DefaultCliPasswdLen {
            doWriteError(w, ErrParamInvalid, "empty username or mismatch encpasswd length")
            return
        }

        user, retcode := CreateUser(req.Username, req.EncPasswd)
        if retcode != ErrOK {
            doWriteError(w, retcode, "")
            return
        }
        var rsp AdminCreateUserRsp
        rsp.Entry = dumpUserEntry(*user)
        doWriteJson(w, rsp)
    },

    "/api/user/change": func(w http.ResponseWriter, r *http.Request) {
        sess, err := smgr.SessionAttach(w, r)
        if err != nil {
            doWriteError(w, ErrSystem, err.Error())
            return
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //session过期，需要重新登录
            doWriteError(w, ErrSessionExpired, "")
            return
        }
        if r.Method != "POST" {
            log.Error("handle http request, inv method %v", r.Method)
            doWriteError(w, ErrMethod, "")
            return
        }

        cando, err := CheckUserPrivilege(username.(string))
        if err != nil {
            doWriteError(w, ErrSystem, err.Error())
            return
        }
        if !cando {
            doWriteError(w, ErrUnauthorized, "")
            return
        }

        var req AdminChangeUserReq
        err = readBodyData(r, &req)
        if err != nil {
            doWriteError(w, ErrParamParseBody, err.Error())
            return
        }
        log.Debug("user change req: %+v", req)
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

func readBodyData(r *http.Request, object interface{}) error {
    reqdata, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return fmt.Errorf("read http body error %v", err)
    }
    if len(reqdata) == 0 {
        return errors.New("no body data found")
    }
    return json.Unmarshal(reqdata, object)
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
    log.Debug("doWriteJson %+v", pkg)
    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func doWriteError(w http.ResponseWriter, errcode int, errmsg string) {
    var rsp AdminError
    rsp.Errcode = errcode
    rsp.Errmsg = errmsg
    log.Debug("doWriteError: %v %v", errcode, errmsg)
    doWriteJson(w, &rsp)
}

