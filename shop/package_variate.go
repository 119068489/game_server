package shop

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
)

const SERVER_NAME = "shop"

var CronJob *easygo.CronJob //定时器任务
var MongoMgr easygo.IMongoDBManager
var ShopInstance *Shop
var WebServreMgr *WebHttpServer

// 服务器配置信息
var PServerInfo *share_message.ServerInfo
var PServerInfoMgr *for_game.ServerInfoManager

//玩家在线管理器
var PlayerOnlineMgr = for_game.NewPlayerOnlineManager()

var PClient3KVMgr *for_game.Client3KVManager

var PPayChannelMgr *for_game.PayChannelManager
var PWebApiForClient *WebHttpForClient

var PWebApiForServer *WebHttpForServer
var PCsvObjManager *for_game.CsvObjManager

func Initialize() {
	MongoMgr = easygo.MongoMgr
	if MongoMgr == nil {
		panic("初始化失败")
	}
	PServerInfoMgr = for_game.NewServerInfoManager()
	InitializeDependDB()
	//ShopInstance.Init()
	ShopInstance = NewShop()
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

}
