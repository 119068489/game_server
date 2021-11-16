package sport_lottery_csgo

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
)

const SERVER_NAME = "sport_lottery_csgo"

var CronJob *easygo.CronJob //定时器任务
var MongoMgr easygo.IMongoDBManager
var RedisMgr easygo.IRedisManager //redis连接
var MongoLogMgr easygo.IMongoDBManager

// 服务器配置信息
var PServerInfo *share_message.ServerInfo
var PServerInfoMgr *for_game.ServerInfoManager

var PClient3KVMgr *for_game.Client3KVManager
var PWebApiForClient *WebHttpForClient

var PWebApiForServer *WebHttpForServer

var delayLotteryTime int64
var blockedBeforeTime int64

var CSGOSysMsgTimeMgr = NewTimer() //CSGO定时发布管理器

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

	if for_game.IS_FORMAL_SERVER {
		delayLotteryTime = int64(easygo.YamlCfg.GetValueAsInt("DELAY_LOTTERY_TIME"))
		blockedBeforeTime = int64(easygo.YamlCfg.GetValueAsInt("BLOCKED_BEFORE_TIME"))
	} else {
		delayLotteryTime = 10800
		blockedBeforeTime = 300
	}

	InitializeDependDB()

	PServerInfoMgr = for_game.NewServerInfoManager()
	PServerInfo = for_game.ReadServerInfoByYaml()
	for_game.InitRedisObjManager(PServerInfo.GetSid())

	//设置服务器为存储服务器
	if for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_ESPORTS_CSGO_SID) == 0 {
		for_game.SetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_ESPORTS_CSGO_SID)
	}

	PClient3KVMgr = for_game.NewClient3KVManager(SERVER_NAME, PServerInfo)
	CronJob = easygo.NewCronJob() //定时器任务
	CronJob.Serve()
	PackageInitialized.Trigger()
}

// 依赖数据库的初始化
func InitializeDependDB() {
}
