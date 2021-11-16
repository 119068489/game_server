package sport_apply

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
)

const SERVER_NAME = "sport_apply"

var CronJob *easygo.CronJob //定时器任务
var MongoMgr easygo.IMongoDBManager
var RedisMgr easygo.IRedisManager //redis连接
var MongoLogMgr easygo.IMongoDBManager

var SportApplyInstance *SportApply

// 服务器配置信息
var PServerInfo *share_message.ServerInfo
var PServerInfoMgr *for_game.ServerInfoManager

var PClient3KVMgr *for_game.Client3KVManager
var PSysParameterMgr *for_game.SysParameterManager

var PWebApiForClient *WebHttpForClient

var PWebApiForServer *WebHttpForServer

func Initialize() {
	MongoMgr = easygo.MongoMgr
	MongoLogMgr = easygo.MongoLogMgr
	if MongoMgr == nil {
		panic("初始化Master DB失败")
	}
	if MongoLogMgr == nil {
		panic("初始化Slave DB失败")
	}
	RedisMgr = easygo.RedisMgr
	if RedisMgr == nil {
		panic("Redis初始化失败")
	}

	PServerInfoMgr = for_game.NewServerInfoManager()
	PServerInfo = for_game.ReadServerInfoByYaml()
	for_game.InitRedisObjManager(PServerInfo.GetSid()) //InitRedisObjManager(PServerInfo.GetSid())

	SportApplyInstance = NewSportApply()

	PClient3KVMgr = for_game.NewClient3KVManager(SERVER_NAME, PServerInfo)
	CronJob = easygo.NewCronJob() //定时器任务
	CronJob.Serve()
	PackageInitialized.Trigger()
}
