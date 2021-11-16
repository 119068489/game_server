package sport_crawl

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/share_message"
)

const SERVER_NAME = "sport_crawl"

var CronJob *easygo.CronJob //定时器任务
var MongoMgr easygo.IMongoDBManager
var RedisMgr easygo.IRedisManager //redis连接
var MongoLogMgr easygo.IMongoDBManager

var SportCrawlInstance *SportCrawl

// 服务器配置信息
var PServerInfo *share_message.ServerInfo
var PServerInfoMgr *for_game.ServerInfoManager

var PClient3KVMgr *for_game.Client3KVManager
var PWebApiForClient *WebHttpForClient

var PWebApiForServer *WebHttpForServer

var QQbucket = util.NewQQbucket() //腾讯云存储桶

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
	for_game.InitRedisObjManager(PServerInfo.GetSid())

	SportCrawlInstance = NewSportCrawl()

	PClient3KVMgr = for_game.NewClient3KVManager(SERVER_NAME, PServerInfo)
	CronJob = easygo.NewCronJob() //定时器任务
	CronJob.Serve()
	PackageInitialized.Trigger()
}

//定时任务
func TimedTask() {
	//整点定时任务
	CronJob.HourEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
		CrawlDataRun()
	})
	//每天定时任务
	CronJob.DayEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
		CrawlScoreHistoryDataRun()
	})
	//每整10分钟定时任务
	CronJob.TenMinEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {

	})

}
