package square

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"time"

	"github.com/astaxie/beego/logs"
)

const SERVER_NAME = "square"

var CronJob *easygo.CronJob //定时器任务
var MongoMgr easygo.IMongoDBManager

// 服务器配置信息
var PServerInfo *share_message.ServerInfo
var PServerInfoMgr *for_game.ServerInfoManager

//玩家在线管理器
var PlayerOnlineMgr = for_game.NewPlayerOnlineManager()

var PClient3KVMgr *for_game.Client3KVManager

var PPayChannelMgr *for_game.PayChannelManager
var PSysParameterMgr *for_game.SysParameterManager

var PWebApiForClient *WebHttpForClient

var PWebApiForServer *WebHttpForServer
var PCsvObjManager *for_game.CsvObjManager

func Initialize() {
	MongoMgr = easygo.MongoMgr
	if MongoMgr == nil {
		panic("初始化失败")
	}
	PServerInfoMgr = for_game.NewServerInfoManager()

	PServerInfo = for_game.ReadServerInfoByYaml()
	for_game.InitRedisObjManager(PServerInfo.GetSid())
	PClient3KVMgr = for_game.NewClient3KVManager(SERVER_NAME, PServerInfo)
	//初始化csv配置
	PCsvObjManager = for_game.NewCsvObjManager()
	CronJob = easygo.NewCronJob() //定时器任务
	CronJob.Serve()

	PackageInitialized.Trigger()
	for_game.ReloadSquareInfo()                                   //从数据库加载社区动态
	easygo.AfterFunc(1*time.Hour, for_game.DealSquareDynamicKeys) //一小时检测一次社交广场 把没有访问的清空
	TaskZanJob()                                                  // 每隔半小时存一次赞信息进数据库
}

// 初始化定时任务
func ReloadTopTimerJob() {
	time.Sleep(3 * time.Second)
	servers := PServerInfoMgr.GetAllServers(for_game.SERVER_TYPE_SQUARE)
	if len(servers) == 0 {
		for_game.SetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_SQUARE_SID)
		sid := for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_SQUARE_SID)
		logs.Info("存储社交广场定时任务的sid:----->%d,当前服务器id:---> %d", sid, PServerInfo.GetSid())
		if PServerInfo.GetSid() == sid {
			// 加载数据库中所有的定时任务管理器
			ls := for_game.GetAllSquareTopTimeMgrFromDB()
			if len(ls) == 0 {
				logs.Error("社交广场重启,没有定时任务需要启动")
			}
			// 重新启动定时任务
			for _, t := range ls {
				for_game.ProcessTopTimer(t)
			}
		}
	}
}

// 每半小时存储赞数据.
func TaskZanJob() {
	logs.Info("定时存储赞数据")
	if for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_SQUARE_SID) == PServerInfo.GetSid() {
		for_game.TaskSaveZanData()
	}
	easygo.AfterFunc(1800*time.Second, TaskZanJob)
}
