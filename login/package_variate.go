// web接口逻辑模块
package login

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"time"

	"github.com/juju/ratelimit"
)

const SERVER_NAME = "login"

var MongoMgr easygo.IMongoDBManager

//var PServerForHall *ServerForHall

var CronJob *easygo.CronJob //定时器任务

var PackageInitialized = easygo.NewEvent()
var IsStopServer bool //服务器已关闭

// 游客登录速率控制
var GuestBucket *ratelimit.Bucket
var PServerInfoMgr *for_game.ServerInfoManager

// 服务器配置信息
var PServerInfo *share_message.ServerInfo

//大厅 ep管理器
var HallEpMgr = for_game.NewEndpointManager()

var PClient3KVMgr *for_game.Client3KVManager

var RedisMgr easygo.IRedisManager //redis连接

var PWebApiForClient *WebHttpForClient

var PWebApiForServer *WebHttpForServer

var PSysParameterMgr *for_game.SysParameterManager
var PCsvObjManager *for_game.CsvObjManager

func Initialize() {
	MongoMgr = easygo.MongoMgr
	if MongoMgr == nil {
		panic("MongoDB初始化失败")
	}
	RedisMgr = easygo.RedisMgr
	if RedisMgr == nil {
		panic("Redis初始化失败")
	}
	//连接服务器管理器
	PServerInfoMgr = for_game.NewServerInfoManager()
	GuestBucket = ratelimit.NewBucketWithQuantum(1*time.Minute, 300, 30)
	PServerInfo = for_game.ReadServerInfoByYaml()
	for_game.InitRedisObjManager(PServerInfo.GetSid())
	PClient3KVMgr = for_game.NewClient3KVManager(SERVER_NAME, PServerInfo)

	//初始化csv配置
	PCsvObjManager = for_game.NewCsvObjManager()

	CronJob = easygo.NewCronJob() //定时器任务
	CronJob.Serve()
	PackageInitialized.Trigger()

}

// 依赖数据库的初始化
func InitializeDependDB() {
	for_game.InitIdGenerator(for_game.TABLE_PLAYER_ACCOUNT, for_game.INIT_PLAYER_ID)
	for_game.InitIdGenerator(for_game.TABLE_TEAM_DATA, for_game.INIT_TEAM_ID)
}
