package def

import (
    "errors"
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

var (
    CodeUserNotExist     = errors.New("user not exist")
    CodeUserDisabled     = errors.New("user disabled")
)
