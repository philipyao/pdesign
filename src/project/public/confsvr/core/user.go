package core

import (
    "crypto/sha256"
    "crypto/sha1"
    "encoding/hex"

    "base/util"
    "base/log"

    "project/public/confsvr/def"
    "project/public/confsvr/db"
)

//预先生成管理员账号
func createAdmin() error {
    name := def.AdminUsername
    exist, err := db.ExistUser(name)
    if err != nil {
        return err
    }
    if exist {
        return nil
    }
    var adminUser def.User
    adminUser.Username = name
    randStr, err := util.GenerateRandomString(def.DefaultSaltLen)
    if err != nil {
        return err
    }
    adminUser.Salt = randStr
    adminUser.Passwd = encodeRawPasswd(adminUser.Username, adminUser.Salt, def.AdminPasswd)
    log.Info("passwd: %v", adminUser.Passwd)

    adminUser.Enabled = 1
    adminUser.IsSuper = 1

    err = db.InsertUser(&adminUser)
    if err != nil {
        return err
    }
    return nil
}

//校验用户登录密码
func VerifyUser(userName, cliPasswd string) (bool, error) {
    //cliPasswd为客户端初次加密后的密码
    user, err := db.QueryUser(userName)
    if err != nil {
        return false, err
    }
    if user.Enabled == 0 {
        //返回错误码
        return false, def.CodeUserDisabled
    }
    //服务器二次加密
    encPwd := encodePasswd(user.Salt, cliPasswd)
    log.Debug("encPwd %v, passwd %v", encPwd, user.Passwd)
    return encPwd == user.Passwd, nil
}

//创建普通账号
func CreateUser(userName, cliPasswd string) (*def.User, int) {
    exist, err := db.ExistUser(userName)
    if err != nil {
        log.Error("error existUser: %v", err)
        return nil, def.ErrSystem
    }
    if exist {
        return nil, def.ErrAccountExist
    }

    var user def.User
    user.Username = userName
    //为用户创建随机盐
    randStr, err := util.GenerateRandomString(def.DefaultSaltLen)
    if err != nil {
        log.Error("err CreateUser: %v", err)
        return nil, def.ErrSystem
    }
    user.Salt = randStr
    user.Passwd = encodePasswd(user.Salt, cliPasswd)
    log.Info("createUser user: <name=%v passwd=%v salt=%v>",
                user.Username, user.Passwd, user.Salt)

    user.Enabled = 1
    user.IsSuper = 0    //都是普通用户

    err = db.InsertUser(&user)
    if err != nil {
        log.Error("error insertUser: %v", err)
        return nil, def.ErrSystem
    }
    return &user, def.ErrOK
}

func QueryUser(userName string) (*def.User, error) {
    return db.QueryUser(userName)
}

//禁用某一普通账号
func DisableUser(userName string) error {
    user, err := db.QueryUser(userName)
    if err != nil {
        return err
    }
    if user.Enabled == 0 {
        return nil
    }
    if user.IsSuper > 0 {
        //忽略超级用户
        return nil
    }
    user.Enabled = 0
    return db.UpdateUser(user)
}

//启用某一普通账号
func EnableUser(userName string) error {
    user, err := db.QueryUser(userName)
    if err != nil {
        return err
    }
    if user.Enabled > 0 {
        return nil
    }
    if user.IsSuper > 0 {
        //忽略超级用户
        return nil
    }
    user.Enabled = 1
    return db.UpdateUser(user)
}

func CheckUserPrivilege(userName string) (bool, error) {
    user, err := db.QueryUser(userName)
    if err != nil {
        return false, err
    }
    return user.IsSuper > 0, nil
}

func ListUser() ([]*def.User, error) {
    dbUsers, err := db.ListUser()
    if err != nil {
        return nil, err
    }
    users := make([]*def.User, 0, len(dbUsers))
    for _, du := range dbUsers {
        if du.IsSuper > 0 {
            continue
        }
        users = append(users, du)
    }
    return users, nil
}

//=====================================================================

//客户端初步加密后的密码加密
func encodePasswd(salt, encPasswd string) string {
    text := salt + encPasswd
    log.Debug("encodePasswd %v", text)
    hash := sha256.New()
    hash.Write([]byte(text))
    return hex.EncodeToString(hash.Sum(nil))
}

//原始明文密码加密
func encodeRawPasswd(userName, salt, passwd string) string {
    //先模拟客户端进行初次加密
    text := def.ClientSaltPart + passwd + userName
    hash := sha1.New()
    hash.Write([]byte(text))
    encPasswd := hex.EncodeToString(hash.Sum(nil))
    log.Debug("simulate client encode password: raw %v enc %v", text, encPasswd)
    //服务器再次加密
    return encodePasswd(salt, encPasswd)
}
