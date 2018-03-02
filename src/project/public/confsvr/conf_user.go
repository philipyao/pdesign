package main

import (
    "errors"
    "crypto/sha256"
    "crypto/sha1"
    "encoding/hex"

    "base/util"
    "base/log"
)

var (
    ErrUserDisabled     = errors.New("user disabled")
)

//预先生成管理员账号
func createAdmin() error {
    name := AdminUsername
    exist, err := existUser(name)
    if err != nil {
        return err
    }
    if exist {
        deleteUser(name)
        //return nil
    }
    var adminUser User
    adminUser.Username = name
    randStr, err := util.GenerateRandomString(DefaultSaltLen)
    if err != nil {
        return err
    }
    adminUser.Salt = randStr
    adminUser.Passwd = encodeRawPasswd(adminUser.Username, adminUser.Salt, AdminPasswd)
    log.Info("passwd: %v", adminUser.Passwd)

    adminUser.Enabled = 1
    adminUser.IsSuper = 1

    err = insertUser(&adminUser)
    if err != nil {
        return err
    }
    return nil
}

//校验用户登录密码
func verifyUser(userName, cliPasswd string) (bool, error) {
    //cliPasswd为客户端初次加密后的密码
    user, err := dbQueryUser(userName)
    if err != nil {
        return false, err
    }
    if user.Enabled == 0 {
        //返回错误码
        return false, ErrUserDisabled
    }
    //服务器二次加密
    encPwd := encodePasswd(user.Salt, cliPasswd)
    log.Debug("encPwd %v, passwd %v", encPwd, user.Passwd)
    return encPwd == user.Passwd, nil
}

//创建普通账号
func CreateUser(userName, cliPasswd string) (*User, int) {
    exist, err := existUser(userName)
    if err != nil {
        log.Error("error existUser: %v", err)
        return nil, ErrSystem
    }
    if exist {
        return nil, ErrAccountExist
    }

    var user User
    user.Username = userName
    //为用户创建随机盐
    randStr, err := util.GenerateRandomString(DefaultSaltLen)
    if err != nil {
        log.Error("err CreateUser: %v", err)
        return nil, ErrSystem
    }
    user.Salt = randStr
    user.Passwd = encodePasswd(user.Salt, cliPasswd)
    log.Info("createUser user: <name=%v passwd=%v salt=%v>",
                user.Username, user.Passwd, user.Salt)

    user.Enabled = 1
    user.IsSuper = 0    //都是普通用户

    err = insertUser(&user)
    if err != nil {
        log.Error("error insertUser: %v", err)
        return nil, ErrSystem
    }
    return &user, ErrOK
}

func QueryUser(userName string) (*User, error) {
    return dbQueryUser(userName)
}

//禁用某一普通账号
func disableUser(userName string) error {
    user, err := dbQueryUser(userName)
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
    return updateUser(user)
}

//启用某一普通账号
func enableUser(userName string) error {
    user, err := dbQueryUser(userName)
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
    return updateUser(user)
}

func CheckUserPrivilege(userName string) (bool, error) {
    user, err := dbQueryUser(userName)
    if err != nil {
        return false, err
    }
    return user.IsSuper > 0, nil
}

func ListUser() ([]*User, error) {
    dbUsers, err := listUser()
    if err != nil {
        return nil, err
    }
    users := make([]*User, 0, len(dbUsers))
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
    text := ClientSaltPart + passwd + userName
    hash := sha1.New()
    hash.Write([]byte(text))
    encPasswd := hex.EncodeToString(hash.Sum(nil))
    log.Debug("simulate client encode password: raw %v enc %v", text, encPasswd)
    //服务器再次加密
    return encodePasswd(salt, encPasswd)
}
