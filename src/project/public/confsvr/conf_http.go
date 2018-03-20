package main

import(
    "fmt"

    log "github.com/philipyao/toolbox/logging"
    "base/srv"

    "github.com/philipyao/phttp"

    "project/public/confsvr/def"
    "project/public/confsvr/core"
)

const (
    KeyUserName     = "username"
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
    EncPasswd       string              `json:"enc_passwd"`    //客户端初次加密后的密码
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

func serveHttp(worker *srv.HTTPWorker) error {
    var err error

    //static file serving
    //err = worker.Static("/", "./dist")
    //if err != nil {
    //    return err
    //}

    //global log middleware
    worker.Use(func(ctx *phttp.Context, next phttp.Next) {
        log.Debug("start log middleware...")
        next()
        log.Debug("end log middleware...")
    })

    //global session middleware
    worker.Use(phttp.AttachSession)

    //登录
    worker.Post("/api/login", func(ctx *phttp.Context) error {
        userName, passwd:= ctx.Request().FormValue("username"), ctx.Request().FormValue("password")
        log.Debug("admin user LOGIN: %v@%v", userName, passwd)
        if userName == "" || passwd == "" {
            doWriteError(ctx, def.ErrParamInvalid, "")
            return nil
        }

        pass, err := core.VerifyUser(userName, passwd)
        if err != nil {
            if err == def.CodeUserDisabled {
                doWriteError(ctx, def.ErrAccountDisabled, "")
                return nil
            } else if err == def.CodeUserNotExist{
                doWriteError(ctx, def.ErrAccountPasswd, "")
                return nil
            }
            doWriteError(ctx, def.ErrSystem, err.Error())
            return nil
        }
        if !pass {
            doWriteError(ctx, def.ErrAccountPasswd, "")
            return nil
        }

        //拉取user信息
        user, err := core.QueryUser(userName)
        if err != nil {
            doWriteError(ctx, def.ErrSystem, err.Error())
            return nil
        }
        if user == nil {
            doWriteError(ctx, def.ErrSystem, "user not exist")
            return nil
        }

        //关联session和user
        sess := ctx.Session()
        sess.Set(KeyUserName, userName)

        var loginRsp AdminLoginRsp
        loginRsp.Userinfo = &SUserinfo{
            Username: userName,
            IsSuper: user.IsSuper,
        }
        doWriteJson(ctx, loginRsp)

        return nil
    })

    logicGroup := worker.NewGroup()
    //logicGroup拦截器：检查登录状态
    logicGroup.Use(func(ctx *phttp.Context, next phttp.Next){
        log.Debug("check login...")
        sess := ctx.Session()
        if sess == nil {
            panic("session not attached!")
        }
        username := sess.Get(KeyUserName)
        if username == nil {
            //session过期，需要重新登录
            log.Info("session expired, need relogin.")
            doWriteError(ctx, def.ErrSessionExpired, "")
            return
        }
        log.Debug("check login ok.")
        next()
    })

    //逻辑接口：列出所有配置
    logicGroup.Post("/api/list", func(ctx *phttp.Context) error {
        log.Debug("list all config.")
        results := core.AllConfig()
        var rsp AdminListRsp
        for _, r := range results {
            rsp.Entries = append(rsp.Entries, dumpConfEntry(r))
        }
        doWriteJson(ctx, rsp)
        return nil
    })

    //逻辑接口：修改配置
    logicGroup.Post("/api/change", func(ctx *phttp.Context) error {
        var req AdminUpdateReq
        err = ctx.Request().JsonBody(&req)
        if err != nil {
            doWriteError(ctx, def.ErrParamParseBody, err.Error())
            return nil
        }
        log.Debug("update config req: %+v", req)

        //开始参数校验
        if len(req.Updates) == 0 && len(req.Adds) == 0 {
            doWriteError(ctx, def.ErrParamInvalid, "no adds or updates provided")
            return nil
        }
        if req.Name == "" || req.Author == "" {
            doWriteError(ctx, def.ErrParamInvalid, "no name or author provided, pls check")
            return nil
        }

        var rsp AdminUpdateRsp
        var failed []string
        var changes []*def.OpChange
        for _, upd := range req.Updates {
            c := core.ConfigByID(upd.ID)

            var change def.OpChange
            change.Namespace = c.Namespace
            change.Key = c.Key
            change.OldValue = c.Value
            change.Value = upd.Value

            err = core.UpdateConfig(upd.ID, upd.Value, upd.Version)
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
            var change def.OpChange
            change.Namespace = add.Namespace
            change.Key = add.Key
            change.OldValue = ""
            change.Value = add.Value

            c, err := core.AddConfig(add.Namespace, add.Key, add.Value)
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
            core.AddOplog(req.Name, req.Comment, req.Author, changes)
        }
        rsp.Failed = failed
        doWriteJson(ctx, rsp)
        return nil
    })

    //逻辑接口：列出所有user
    logicGroup.Get("/api/user/list", func(ctx *phttp.Context) error {
        log.Debug("list all user")

        //校验用户权限
        sess := ctx.Session()
        username := sess.Get(KeyUserName)
        cando, err := core.CheckUserPrivilege(username.(string))
        if err != nil {
            doWriteError(ctx, def.ErrSystem, err.Error())
            return nil
        }
        if !cando {
            doWriteError(ctx, def.ErrUnauthorized, "")
            return nil
        }
        var rsp AdminListUserRsp
        users, err := core.ListUser()
        if err != nil {
            rsp.Errmsg = err.Error()
        } else {
            for _, u := range users {
                rsp.Entries = append(rsp.Entries, dumpUserEntry(*u))
            }
        }
        doWriteJson(ctx, rsp)
        return nil
    })

    //逻辑接口：新建user
    logicGroup.Post("/api/user/create", func(ctx *phttp.Context) error {
        sess := ctx.Session()
        username := sess.Get(KeyUserName)
        cando, err := core.CheckUserPrivilege(username.(string))
        if err != nil {
            doWriteError(ctx, def.ErrSystem, err.Error())
            return nil
        }
        if !cando {
            doWriteError(ctx, def.ErrUnauthorized, "")
            return nil
        }
        var req AdminCreateUserReq
        err = ctx.Request().JsonBody(&req)
        if err != nil {
            doWriteError(ctx, def.ErrParamParseBody, err.Error())
            return nil
        }
        log.Debug("user create req: %+v, passwd len %v", req, len(req.EncPasswd))
        if req.Username == "" || len(req.EncPasswd) != def.DefaultCliPasswdLen {
            doWriteError(ctx, def.ErrParamInvalid, "empty username or mismatch encpasswd length")
            return nil
        }

        user, retcode := core.CreateUser(req.Username, req.EncPasswd)
        if retcode != def.ErrOK {
            doWriteError(ctx, retcode, "")
            return nil
        }
        var rsp AdminCreateUserRsp
        rsp.Entry = dumpUserEntry(*user)
        doWriteJson(ctx, rsp)
        return nil
    })

    //逻辑接口：修改用户
    logicGroup.Post("/api/user/change", func(ctx *phttp.Context) error {
        sess := ctx.Session()
        username := sess.Get(KeyUserName)
        cando, err := core.CheckUserPrivilege(username.(string))
        if err != nil {
            doWriteError(ctx, def.ErrSystem, err.Error())
            return nil
        }
        if !cando {
            doWriteError(ctx, def.ErrUnauthorized, "")
            return nil
        }

        var req AdminChangeUserReq
        err = ctx.Request().JsonBody(&req)
        if err != nil {
            doWriteError(ctx, def.ErrParamParseBody, err.Error())
            return nil
        }
        log.Debug("user change req: %+v", req)
        if req.Username == "" {
            doWriteError(ctx, def.ErrParamInvalid, "empty username")
            return nil
        }

        var rsp AdminChangeUserRsp
        if req.Enable == 0 {
            err = core.DisableUser(req.Username)
        } else {
            err = core.EnableUser(req.Username)
        }
        if err != nil {
            rsp.Errmsg = err.Error()
        }
        doWriteJson(ctx, rsp)
        return nil
    })

    return nil
}

//===================================================================

func dumpConfEntry(c def.Config) *ConfEntry {
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

func dumpUserEntry (u def.User) *UserEntry {
    return &UserEntry{
        Username: u.Username,
        Enabled: u.Enabled,
        CreatedAt: int(u.CreatedAt.Unix()),
    }
}

func doWriteJson(ctx *phttp.Context, pkg interface{}) {
    err := ctx.Response().Json(pkg)
    if err != nil {
        log.Error("WriteJson error: %v", err)
    }
}

func doWriteError(ctx *phttp.Context, errcode int, errmsg string) {
    var pkg AdminError
    pkg.Errcode = errcode
    pkg.Errmsg = errmsg
    log.Debug("doWriteError: %v %v", errcode, errmsg)
    doWriteJson(ctx, pkg)
}

