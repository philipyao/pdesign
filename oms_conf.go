package server_conf

type OmsHost struct{
	Host            string

	IP              string
	WanIP           string
	SSHPort         uint32
	CpuNum          uint32

	DBAddr          string
}

type OmsCluster struct {
	Cluster         string      //cluster，比如world2000@2(同一个cluster可以部署在多台机器，所以区分了@1，@2...)
	ClusterID       uint32      //比如1(public), 2000(world), 7007(zone)等
	ClusterLayer    string      //public | world | zone
	Host            string      //部署在哪台机器上
	Loc             string      //地区版本
}

type OmsServer struct {
	ID              uint32      //自增ID
	TypeID          uint32      //进程类型编号
	TypeName        string     //如gamesvr，gatesvr
	Cluster         string      //cluster，比如world2000@2
	ClusterID       uint32      //比如2000
	StartIdx        uint32      //部署开始索引
	EndIdx          uint32      //部署结束索引
	PortBase        uint32      //起始端口号，开多个实例则端口递增
}

//部署DB
//1、根据OMS的DB配置去对应的机器上建库表
//2、生成所有表的访问路由,供业务访问(和rpc路由一样写入zk)
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
	OmsParamCommonDBAddr                 = "DBAddr"
	OmsParamCommonDBUser                 = "DBUser"
	OmsParamCommonDBPasswd               = "DBPasswd"
	OmsParamCommonLogMaxSize             = "LogMaxSize"
	OmsParamCommonLogLevel               = "LogLevel"
	OmsParamCommonRedisAddr              = "RedisAddr"
	OmsParamCommonRedisPasswd            = "RedisPasswd"
)

type OmsParam struct {
	Params          map[string]string
}
