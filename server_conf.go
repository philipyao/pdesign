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
	ParamCommonZKAddr                 = "ZKAddr"
	ParamCommonZKUser                 = "ZKUser"
	ParamCommonDBAddr                 = "DBAddr"
	ParamCommonDBUser                 = "DBUser"
	ParamCommonDBPasswd               = "DBPasswd"
	ParamCommonLogMaxSize             = "LogMaxSize"
	ParamCommonLogLevel               = "LogLevel"
	ParamCommonRedisAddr              = "RedisAddr"
	ParamCommonRedisPasswd            = "RedisPasswd"

	//gatesvr类型
	ParamGatesvrMaxSession            = "MaxSession"
	ParamGatesvrTcpAccpetTimeout      = "TcpAcceptTimeout"
	ParamGatesvrLogLevel              = "LogLevel"

	//gamesvr类型
	ParamGamesvrLogMaxSize            = "LogMaxSize"
	ParamGamesvrLogLevel              = "LogLevel"
	ParamGamesvrLuaNum                = "LuaNum"

	//namesvr类型
	ParamNamesvrNameMaxlen            = "NameMaxlen"

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

//===================================================
//各server独有配置
type ParamServerDef struct {
	//以下为通用字段
	Server              string              /*  clusterid.typeid.index 三段式服务器进程标志 */
	Cluster             string              /*  部署的cluster，如world2000@2 */
	Typestr             string              /*  进程类型，如ranksvr，注意不是typeid */
	Addr                string              /*  ip地址 */
	Port                uint16              /*  端口 */
	Profurl             string              /*  进程分析url */

	//自定义字段（可以覆盖类型配置和通用配置中同一名字的字段）
	Custom              map[string]string   /*  各进程独有的字段，key/value格式 */
}

//数据表路由表
type DBTableRoute struct {
	Table               string          //数据表名，TableName+Seq, 如tbl_user_1
	DBName              string          //库名, 如db_main
	DBAddr              string          //数据库连接地址
}




