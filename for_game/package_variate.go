// 包变量
package for_game

var (
	// easygo.MongoMgr         easygo.IMongoDBManager
	// YamlCfg          IYamlConfig
	EDITION          string // 发行版
	IS_FORMAL_SERVER bool   //是否正式服：true正式服, false测试服

	SERVER_CENTER_ADDR string //etcd中心地址

	IS_TFSERVER bool //是否走转发
	//RedisPoolObj IRedisPool
)

var MessageMarkInfo *MessageMarkInfoMgr
var PDirtyWordsMgr *DirtyWordMgr
