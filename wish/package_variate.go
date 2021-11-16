package wish

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

const SERVER_NAME = "wish"

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

var WishBoxMgr *WishBoxManager
var PoolMgr *PoolManager

func Initialize() {
	MongoMgr = easygo.MongoMgr
	if MongoMgr == nil {
		panic("初始化失败")
	}
	PServerInfoMgr = for_game.NewServerInfoManager()
	WishBoxMgr = NewWishBoxManager()
	PoolMgr = NewPoolManager()

	PServerInfo = for_game.ReadServerInfoByYaml()
	for_game.InitRedisObjManager(PServerInfo.GetSid())
	PClient3KVMgr = for_game.NewClient3KVManager(SERVER_NAME, PServerInfo)
	//初始化csv配置
	PCsvObjManager = for_game.NewCsvObjManager()
	CronJob = easygo.NewCronJob() //定时器任务
	CronJob.Serve()

	PackageInitialized.Trigger()
	TaskPoolJob() // 定时存储数据.
	//easygo.Spawn(PlatformBackJob) // 定时任务回收
	easygo.Spawn(TaskUpBox) // 检查定时上线盲盒的定时任务.
	TimedTask()
	CleanRedisLock()
}
func InitializeDependDB() {
	//初始化许愿池的数据
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_MENU, "初始化许愿池首页菜单", true, bson.M{}, for_game.InitWishMenu())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BRAND, "初始化许愿池物品品牌", true, bson.M{}, for_game.InitWishBrand())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM_TYPE, "初始化许愿池物品类型", true, bson.M{}, for_game.InitWishItemType())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_STYLE, "初始化许愿池物品款式", true, bson.M{}, for_game.InitWishStyle())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM, "初始化许愿池物品配置", true, bson.M{}, for_game.InitWishItem())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM, "初始化许愿池盲盒商品", true, bson.M{}, for_game.InitWishBoxItem())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX, "初始化许愿池盲盒", true, bson.M{}, for_game.InitWishBox())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_COLLECTION, "初始化已收藏的愿望盒数据", true, bson.M{}, for_game.InitPlayerWishCollection())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA, "初始化已许愿的愿望盒数据", true, bson.M{}, for_game.InitPlayerWishData())
	for_game.InitIdGenerator(for_game.TABLE_WISH_PLAYER, INIT_WISH_PLAYER_ID)
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER, "初始化许愿用户表数据", false, bson.M{}, for_game.InitWishPlayer())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_REASON, "初始化回收理由数据", false, bson.M{}, for_game.InitWishRecycleReason())
}

// 每半小时存储赞数据.
func TaskPoolJob() {
	logs.Info("定时存储水池数据")
	SavePoolDataToMongo()
	easygo.AfterFunc(600*time.Second, TaskPoolJob) // 600秒更新,十分钟
}

// 处理回收硬币
//func PlatformBackJob() {
//	logs.Info("定时回收许愿池初始化")
//	// 查询有效物品数据
//	list := for_game.GetPlayerWishItemByStatus()
//	if len(list) == 0 {
//		return
//	}
//	now := time.Now().Unix()
//	for _, v := range list {
//		t := v.GetExpireTime() - now
//		if t < 0 {
//			continue
//		}
//		a := v
//		fun := func() {
//			PlantBack([]int64{a.GetId()}, a.GetPlayerId())
//		}
//		easygo.AfterFunc(time.Duration(t)*time.Second, fun)
//	}
//}

func TaskUpBox() {
	logs.Info("检查需要恢复上架的盲盒定时任务")
	boxList := for_game.GetTaskBox()
	if len(boxList) == 0 {
		return
	}
	for _, v := range boxList {
		a := v
		CheckBoxStatus(a)
	}
}

func TimedTask() {
	CronJob = easygo.Cronjob
	//整点定时任务
	CronJob.HourEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
		SumGuardianCoinNum()
	})
	CronJob.WeekEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
		SumWeekMonth(for_game.WISH_ACT_H5_WEEK)
	})
	//每月定时任务
	CronJob.MonthEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
		SumWeekMonth(for_game.WISH_ACT_H5_MONTH)
	})
}

//清理redis锁
func CleanRedisLock() {
	easygo.RedisMgr.GetC().Delete(WISH_LOGIN_MUTEX_LOCK)
}
