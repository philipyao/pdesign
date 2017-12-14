package server_conf

type OmsHost struct{
	Host            string

	IP              string
	WanIP           string
	SSHPort         uint32
	CpuNum          uint32

	DBAddr          string
}

//部署cluster
//1、部署cluster到特定机器的特定目录
//2、生成服务器进程配置表数据
type OmsCluster struct {
	Cluster         string      //cluster，比如world2000@2(同一个cluster可以部署在多台机器，所以区分了@1，@2...)
	ClusterID       uint32      //比如1(public), 2000(world), 7007(zone)等
	ClusterLayer    string      //public | world | zone
	Host            string      //部署在哪台机器上
	Loc             string      //地区版本
}

//server类型定义
type OmsServerType struct {
    TypeName        string      //如gamesvr，gatesvr
    TypeID          uint32      //进程类型编号
}

//server
type OmsServer struct {
	ID              uint32      //自增ID
	TypeName        string      //如gamesvr，gatesvr
    StartIdx        uint32      //部署开始索引
    EndIdx          uint32      //部署结束索引
	Cluster         string      //cluster，比如world2000@2
	ClusterID       uint32      //比如2000
	PortBase        uint32      //起始端口号，开多个实例则端口递增
}

//部署DB
//1、根据OMS的DB配置去对应的机器上建库表(db_biz游戏逻辑数据库，db_conf_server服务器配置库, db_route_db游戏逻辑表路由库)
//2、所有表库建好后生成路由表tbl_route_table数据

type OmsBizDBTable struct {
	ID              uint32      //自增ID
	TableName       string      //表名
	DBName          string      //所在的库名
	StartIdx        uint32      //起始索引
	EndIdx          uint32      //结束索引
	Host            string      //数据表建在哪台机器上
}

const (
	OmsPkgDownloadUrl                    = "PkgDownloadUrl"
	OmsInstallDir                        = "InstallDir"
	OmsUser                              = "OmsUser"    //进行部署操作的user,需要预先设置sudo权限
	OmsParamCommonZKAddr                 = "ZKAddr"
	OmsParamCommonZKUser                 = "ZKUser"
	OmsParamCommonDBUser                 = "DBUser"
	OmsParamCommonDBPasswd               = "DBPasswd"
	OmsParamCommonLogMaxSize             = "LogMaxSize"
	OmsParamCommonLogLevel               = "LogLevel"
	OmsParamCommonRedisAddr              = "RedisAddr"
	OmsParamCommonRedisPasswd            = "RedisPasswd"
)

//部署配置
//1、从OMS的DB中拷贝服务器参数表数据(主要是通用参数和特定类型服务器进程参数)
type OmsParam struct {
	Params          map[string]string
}
