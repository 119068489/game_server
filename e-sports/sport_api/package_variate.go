package sport_api

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"net/url"
	"time"
)

const SERVER_NAME = "sport_api"

var yeZiUrl string
var yeZiAppKey string
var yeZiAppSecret string

var CronJob *easygo.CronJob //定时器任务
var MongoMgr easygo.IMongoDBManager
var RedisMgr easygo.IRedisManager //redis连接
var MongoLogMgr easygo.IMongoDBManager

var SportApiInstance *SportApi
var WebServreMgr *WebHttpServer

// 服务器配置信息
var PServerInfo *share_message.ServerInfo
var PServerInfoMgr *for_game.ServerInfoManager

var PClient3KVMgr *for_game.Client3KVManager
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

	if for_game.IS_FORMAL_SERVER {
		yeZiUrl = easygo.YamlCfg.GetValueAsString("YEZI_API_FORMAL_URL")
	} else {
		yeZiUrl = easygo.YamlCfg.GetValueAsString("YEZI_API_TEST_URL")
	}

	_, err := url.Parse(yeZiUrl)
	if err != nil {
		panic("野子科技url配置错误!")
	}

	yeZiAppKey = easygo.YamlCfg.GetValueAsString("YEZI_APP_KEY")
	yeZiAppSecret = easygo.YamlCfg.GetValueAsString("YEZI_APP_SECRET")

	//游戏英雄、装备数据
	//TODO:目前没有实时数据
	//for_game.SetRedisGameRealTimeBase()

	PServerInfoMgr = for_game.NewServerInfoManager()
	PServerInfo = for_game.ReadServerInfoByYaml()
	for_game.InitRedisObjManager(PServerInfo.GetSid())

	SportApiInstance = NewSportApi()

	//设置服务器为存储服务器、并且在用的时候不管怎么样、只要初始化一次成功以后启动不会再做初始化
	//如果要做初始化就要删除该redis的key
	if for_game.GetCurrentSaveServerSidForSportApi(PServerInfo.GetSid(), for_game.REDIS_SAVE_ESPORTS_API_SID) == 0 {

		for_game.SetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_ESPORTS_API_SID)
		//初始化
		InitializeDependDB()
	} else {
		sid := PServerInfo.GetSid()
		saveServerSid := for_game.GetCurrentSaveServerSid(sid, for_game.REDIS_SAVE_ESPORTS_API_SID)
		if saveServerSid == sid {
			//初始化未开始、进行中的比赛的详情、比赛动态信息表(早盘、滚盘)信息
			//两队历史交锋、两队胜败统计、两队天敌克制统计(统计一次即可)
			DoGetYeZiAllGameDetailsFunc(for_game.ESPORTS_INCREMENT_FLAG_2)
		}
	}

	PClient3KVMgr = for_game.NewClient3KVManager(SERVER_NAME, PServerInfo)
	CronJob = easygo.NewCronJob() //定时器任务
	CronJob.Serve()
	PackageInitialized.Trigger()

	//TODO 测试环境用模拟回调
	if for_game.IS_FORMAL_SERVER {
	} else {
		TimedTask()
	}
}

// 依赖数据库的初始化
func InitializeDependDB() {

	//初始化未开始比赛列表:王者荣耀、DOTA2、英雄联盟、CSGO
	DoGetYeZiAllGamesFunc()

	//初始化未开始、进行中的比赛的详情、比赛动态信息表(早盘、滚盘)信息
	//两队历史交锋、两队胜败统计、两队天敌克制统计(统计一次即可)
	DoGetYeZiAllGameDetailsFunc(for_game.ESPORTS_INCREMENT_FLAG_1)

}

// 定时任务
func TimedTask() {

	easygo.AfterFunc(time.Duration(3600)*time.Second, func() {
		GameListCallBack()
	})

	easygo.AfterFunc(time.Duration(300)*time.Second, func() {
		GameDetailCallBack()
	})

	easygo.AfterFunc(time.Duration(600)*time.Second, func() {
		MornBetCallBack()
	})

	easygo.AfterFunc(time.Duration(60)*time.Second, func() {
		RollBetCallBack()
	})

	//easygo.AfterFunc(time.Duration(20)*time.Second, func() {
	//	RealTimeCallBack()
	//})
}
