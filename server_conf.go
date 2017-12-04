package server_conf

///////////////////////////////////////////////////////////////////////////////////////////
///// 本package定义了服务器配置格式，对应了这些配置参数在数据库中的存储格式
///// 配置分为三个层级：
//          1）通用配置，所有服务器进程通用的配置
//          2）类型配置，区分不同类型进程的配置，同一类型的进程共享同一套配置，但是不同类型进程配置会不一样。
//              例如所有gatesvr进程监听外网连接时，可以配置统一的AcceptTimeout。
//          3）进程配置，例如各个服务器进程自己监听的rpc通信port，都不相同
///// 进程配置可以依次向上覆盖同一名字的配置，例如LogLevel配置可以被覆盖。
///////////////////////////////////////////////////////////////////////////////////////////


//=================================================
//参数表key定义
const (
	//通用
	ParamCommonZKAddr                 = "zkaddr"
	ParamCommonZKUser                 = "zkuser"
    ParamCommonZKPasswd               = "zkpasswd"
	ParamCommonDBUser                 = "dbuser"
	ParamCommonDBPasswd               = "dbpasswd"
	ParamCommonLogMaxSize             = "log_maxsize"
	ParamCommonLogLevel               = "log_level"

	//gatesvr类型
	ParamGatesvrMaxConn               = "gate_max_conn"
    ParamGatesvrMaxOnline             = "gate_max_online"
	ParamGatesvrTCPAccpetTimeout      = "gate_tcp_accept_timeout"
    ParamGatesvrTCPMaxChanpkgs        = "gate_tcp_max_chanpkgs"
    ParamGatesvrKCPResend             = "gate_kcp_resend"
    ParamGatesvrKCPMtu                = "gate_kcp_mtu"

	//gamesvr类型
	ParamGamesvrLogMaxSize            = "log_maxsize"       //覆盖common
	ParamGamesvrLogLevel              = "log_level"         //覆盖common
	ParamGamesvrLuaNum                = "game_lua_num"
    ParamGamesvrEnableVip             = "game_enable_vip"
    ParamGamesvrVerno                 = "game_verno"
    ParamGamesvrEnableZip             = "game_enable_zip"

    //namesvr类型
    ParamNamesvrNameMaxlen            = "name_maxlen"
    ParamNamesvrMaxtry                = "name_maxtry"
    ParamNamesvrMaxfindn              = "name_maxfindn"
    ParamNamesvrMinlen                = "name_minlen"

    //ranksvr类型
    ParamRanksvrExtURL                = "rank_ext_url"
    ParamRanksvrURL                   = "rank_url"
    ParamRanksvrRedisAddr             = "rank_redis_addr"
    ParamRanksvrRedisPasswd           = "rank_redis_passwd"

)

//==================================================
//所有server通用配置
type ParamDefCommon struct {
	//key/value键值对
	Params      map[string]string
}

//==================================================
//各server类型配置
type ParamGatesvrDef struct {
	//key/value键值对
	Params      map[string]string
}
type ParamGamesvrDef struct {
	//key/value键值对
	Params      map[string]string
}
type ParamNamesvrDef struct {
	//key/value键值对
	Params      map[string]string
}
type ParamRanksvrDef struct {
    //key/value键值对
    Params      map[string]string
}


//===================================================
//各server最终生成的配置
type ParamServerDef struct {
	Typename            string
    Typeid              uint32              //进程编号

	//（同名字段向上覆盖类型配置和通用配置中的字段）
	Params              map[string]string   //进程参数
}

//数据表路由表
type DBTableRoute struct {
	Table               string          //数据表名，TableName+Seq, 如tbl_user_1
	DBName              string          //库名, 如db_main
	DBAddr              string          //数据库连接地址
}




