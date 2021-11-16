package statistics

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"time"
)

const SERVER_NAME = "statistics"

var MongoMgr easygo.IMongoDBManager
var MongoLogMgr easygo.IMongoDBManager

// 服务器配置信息
var PServerInfo *share_message.ServerInfo
var PServerInfoMgr *for_game.ServerInfoManager
var PClient3KVMgr *for_game.Client3KVManager
var PWebApiForClient *WebHttpForClient

var PWebApiForServer *WebHttpForServer

//关键词抓取管理
var PWordManager *WordManager

//商城阿里验证变量
var ali_item_scan_by_id *AliItemScanMap = &AliItemScanMap{}

var PCsvObjManager *for_game.CsvObjManager

// 定时任务
var CronJob *easygo.CronJob

func Initialize() {
	MongoMgr = easygo.MongoMgr
	if MongoMgr == nil {
		panic("MongoMgr初始化失败")
	}
	MongoLogMgr = easygo.MongoLogMgr
	if MongoLogMgr == nil {
		panic("MongoLogMgr初始化失败")
	}
	InitializeDependDB()
	PServerInfoMgr = for_game.NewServerInfoManager()
	PServerInfo = for_game.ReadServerInfoByYaml()
	//初始化csv配置
	PCsvObjManager = for_game.NewCsvObjManager()
	for_game.InitRedisObjManager(PServerInfo.GetSid())
	PClient3KVMgr = for_game.NewClient3KVManager(SERVER_NAME, PServerInfo)
	TimedTask()
}

// 依赖数据库的初始化
func InitializeDependDB() {
	//定时清理一些无效的redis keys
	easygo.AfterFunc(time.Second, DelRedisTimeOutkeys)
	//定时处理注销订单:10分钟处理一次
	easygo.AfterFunc(time.Second*600, DealAccountCancel)
	//easygo.AfterFunc(time.Second*200, DealAccountCancel)
	// //每小时衰减热门分值
	// easygo.AfterFunc(time.Second*3600, CoolingSquare)
	//商城的定时任务
	//商城超时取消订单(第一次启动立即执行)
	easygo.Spawn(UpdateOrderList)
	//商城阿里云验证发送(第一次启动等待5s执行)
	easygo.AfterFunc(time.Second*5, DoAliItemScanSend)
	//商城阿里云验证图片和视频异步结果取得(第一次启动等待10s执行)
	easygo.AfterFunc(time.Second*10, GetImageVideoScanRst)
	//腾讯云定期删除违禁图片
	easygo.AfterFunc(time.Second*time.Duration(COS_IMAGE_DELETE_TIME), UpdateDelImageFromTX)
}

// 定时任务
func TimedTask() {
	UpdateTopicTopExp()
	CronJob = easygo.Cronjob
	//整点定时任务
	CronJob.HourEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
		if for_game.IS_FORMAL_SERVER {
			CoolingSquare()
			UpdateHotTopic()
		}
	})

	//每整10分钟定时任务
	CronJob.TenMinEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
		//测试服每10分钟更新
		if !for_game.IS_FORMAL_SERVER {
			CoolingSquare()
			UpdateHotTopic()
		}
		UpdateMallItemSaleStatus() //下架商城限时售卖的商品
		//TODO 暂时处理亲密度
		//easygo.Spawn(UpdatePlayerIntimacy) //定时处理玩家亲密度，超过3天的要处理
	})

	//每天定时任务
	CronJob.DayEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
		easygo.Spawn(UpdateBCoinExpiration) //定时任务,硬币商城硬币回收
		easygo.Spawn(UpdatePlayerIntimacy)  //定时处理玩家亲密度，超过3天的要处理
	})

	// 每天八点执行通知过期
	CronJob.DayEightClockEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
		DealPreBCoinExpiration() //处理一天内将要过期的硬币
		NoticeProductExp()       // 通知物品过期
	})

}
