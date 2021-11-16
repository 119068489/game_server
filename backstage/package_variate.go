// 包变量全在这里定义吧
package backstage

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/share_message"
)

var (
	Server4Brower  *ServerForBrower
	IS_GOOGLE_AUTH bool
	GOOGLE_SECRET  string
)

const SERVER_NAME = "backstage"

// 浏览器 endpoint 管理器
var InspectSvr *InspectServer
var BrowerEpMgr = NewBrowerEndpointManager() // 用 endpointId 作 key 的
var BrowerEpMp = NewBrowerEndpointMapping()  // 用 userId 作 key 的
var BackstageEpMgr = for_game.NewEndpointManager()

var MongoMgr easygo.IMongoDBManager
var MongoLogMgr easygo.IMongoDBManager

var TimerMgr = NewTimer()                 //通知推送定时管理器
var ArticleTimeMgr = NewTimer()           //文章推送定时管理器
var UserSweetsTimeMgr = NewTimer()        //小助手推送定时管理器
var SendDynamicTimeMgr = NewTimer()       //社交广场推送定时管理器
var SendEsportsNewsTimeMgr = NewTimer()   //电竞新闻定时发布管理器
var SendEsportsVideoTimeMgr = NewTimer()  //电竞视频定时发布管理器
var SendEsportsSysMsgTimeMgr = NewTimer() //电竞系统消息定时发布管理器
var EsportsActCloseTimeMgr = NewTimer()   //电竞活动关闭任务管理器

var WebApiServreMgr *WebHttpServer

var PSysParameterMgr *for_game.SysParameterManager

var PWebApiForClient *WebHttpForClient

var PWebApiForServer *WebHttpForServer

// // 大厅连接器
// var _HallConnector *HallConnector

// 定时任务
var CronJob *easygo.CronJob

// 服务器配置信息
var PServerInfo *share_message.ServerInfo
var PServerInfoMgr *for_game.ServerInfoManager

//大厅连接管理
var HallConMgr = for_game.NewConnectorManager()

//后台连接管理
var BackstageConMgr = for_game.NewConnectorManager()

//后台用户在线管理器
var UserOnlineMgr = for_game.NewPlayerOnlineManager()

//玩家在线管理器
var PlayerOnlineMgr = for_game.NewPlayerOnlineManager()

//玩家IP查询管理器
var IpSearchMgr = for_game.NewIpSearchMgr()

var PClient3KVMgr *for_game.Client3KVManager

var QQbucket = util.NewQQbucket() //腾讯云存储桶

var PCsvObjManager *for_game.CsvObjManager

func Initialize() {

	MongoMgr = easygo.MongoMgr
	MongoLogMgr = easygo.MongoLogMgr
	if MongoMgr == nil {
		panic("初始化Master DB失败")
	}
	if MongoLogMgr == nil {
		panic("初始化Slave DB失败")
	}
	PServerInfo = for_game.ReadServerInfoByYaml()
	PServerInfo.ConNum = easygo.NewInt32(1000)
	for_game.InitRedisObjManager(PServerInfo.GetSid())
	PClient3KVMgr = for_game.NewClient3KVManager(SERVER_NAME, PServerInfo)
	PServerInfoMgr = for_game.NewServerInfoManager()
	InitPoolManager() // 许愿池redis
	//初始化csv配置
	PCsvObjManager = for_game.NewCsvObjManager()

	PSysParameterMgr = for_game.NewSysParameterManager()
	//address := for_game.MakeAddress(PServerInfo.GetIp(), PServerInfo.GetClientPort())
	address := for_game.MakeAddress("0.0.0.0", PServerInfo.GetClientWSPort())
	Server4Brower = NewServerForBrower(address)

	//设置服务器为存储服务器
	if for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_BACKSTAGE_SID) == 0 {
		for_game.SetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_BACKSTAGE_SID)
	}

	for_game.PDirtyWordsMgr = for_game.NewDirtyWordMgr() //屏蔽词库

	IS_GOOGLE_AUTH = easygo.YamlCfg.GetValueAsBool("IS_GOOGLE_AUTH")
	GOOGLE_SECRET = easygo.YamlCfg.GetValueAsString("GOOGLE_SECRET")

	//启动http webapi服务
	WebApiServreMgr = NewWebHttpServer()
	//easygo.Spawn(WebApiServreMgr.Serve)

	CronJob = easygo.Cronjob
}
