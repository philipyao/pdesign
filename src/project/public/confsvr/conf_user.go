package main

import (
    "crypto/sha256"
    "crypto/sha1"
    "encoding/hex"

    "base/util"
    "base/log"
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
    user, err := queryUser(userName)
    if err != nil {
        return false, err
    }
    if user.Enabled == 0 {
        //返回错误码
    }
    //服务器二次加密
    encPwd := encodePasswd(user.Salt, cliPasswd)
    log.Debug("encPwd %v, passwd %v", encPwd, user.Passwd)
    return encPwd == user.Passwd, nil
}

//创建普通账号
func createUser(userName, cliPasswd string) error {
    var user User
    user.Username = userName
    randStr, err := util.GenerateRandomString(DefaultSaltLen)
    if err != nil {
        return err
    }
    user.Salt = randStr
    user.Passwd = encodePasswd(user.Salt, cliPasswd)
    log.Info("createUser user: <name=%v passwd=%v salt=%v>",
                user.Username, user.Passwd, user.Salt)

    user.Enabled = 1
    user.IsSuper = 0

    err = insertUser(&user)
    if err != nil {
        return err
    }
    return nil
}

//禁用某一普通账号
func disableUser(userName string) error {
    user, err := queryUser(userName)
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
    user, err := queryUser(userName)
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
