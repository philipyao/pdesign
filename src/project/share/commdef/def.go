package commdef

//服务器类型定义
const (
    //public
    ServerTypePlatsvr               = 1
    ServerTypePaysvr                = 2
    ServerTypeConfsvr               = 3

    //world级别
    ServerTypeRanksvr               = 21
    ServerTypeReplaysvr             = 22

    //zone级别
    ServerTypeGatesvr               = 51
    ServerTypeGamesvr               = 52
    ServerTypeSessionsvr            = 53
    ServerTypeDBsvr                 = 54
)

const (
    ZKPrefixConfig      = "/config"
)

const (
    ConfigKeyZKAddr     = "zkaddr"
)
