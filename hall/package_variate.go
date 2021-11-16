// web接口逻辑模块
package hall

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/mongo_init"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo/bson"

	"github.com/astaxie/beego/logs"

	"github.com/juju/ratelimit"
)

const SERVER_NAME = "hall"

// var DBMgr = for_game.NewDBManager(CreateTableDDL, AlterTableDDL) // AlterTableDDL
const TIME_TO_SAVE_REDIS = 600 * time.Second         // 600
const TIME_TO_ADD_FULL_PLAYER = 1800 * time.Second   // 30分钟
const TIME_TO_SEND_TOPIC_DYNAMIC = 600 * time.Second // 10分钟

var MongoMgr easygo.IMongoDBManager
var MongoLogMgr easygo.IMongoDBManager

var PackageInitialized = easygo.NewEvent()

var IsStopServer bool
var HallServerId int32

// 游客登录速率控制
var GuestBucket *ratelimit.Bucket
var WebServreMgr *WebHttpServer

// 服务器配置信息
var PServerInfo *share_message.ServerInfo

//服务器信息管理
var PServerInfoMgr *for_game.ServerInfoManager

//玩家在线管理器
var PlayerOnlineMgr = for_game.NewPlayerOnlineManager()

//玩家IP查询管理器
var IpSearchMgr = for_game.NewIpSearchMgr()

//群在线管理器
var TeamOnlineMgr = for_game.NewTeamOnlineManager()
var UpdateMgr *UpdateManager

var PClient3KVMgr *for_game.Client3KVManager

var PWebMiaoDaoPay *WebMiaoDaoPay
var PWebHuiChaoPay *WebHuiChaoPay
var PWebTongTongPay *WebTongTongPay
var PWebYunTongShangPay *WebYunTongShangPay

var PPayChannelMgr *for_game.PayChannelManager

var PSysParameterMgr *for_game.SysParameterManager

var RedisMgr easygo.IRedisManager //redis连接
// 定时任务
var CronJob *easygo.CronJob

var PWebApiForClient *WebHttpForClient

var PWebApiForServer *WebHttpForServer
var PCsvObjManager *for_game.CsvObjManager

// var QQbucket = util.NewQQbucket() //腾讯云存储桶

func Initialize(b ...bool) {
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

	UpdateMgr = NewUpdateManager()
	GuestBucket = ratelimit.NewBucketWithQuantum(1*time.Minute, 300, 30)
	PServerInfo = for_game.ReadServerInfoByYaml()
	for_game.InitRedisObjManager(PServerInfo.GetSid())
	//设置服务器为存储服务器
	if for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_HALL_SID) == 0 {
		for_game.SetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_HALL_SID)
	}
	if for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_HALL_SID) == PServerInfo.GetSid() {
		for_game.CleanOrderDoing() //清理订单锁
	}

	PClient3KVMgr = for_game.NewClient3KVManager(SERVER_NAME, PServerInfo)
	PServerInfoMgr = for_game.NewServerInfoManager()

	PWebMiaoDaoPay = NewWebMiaoDaoPay()
	PWebHuiChaoPay = NewWebHuiChaoPay()
	PWebTongTongPay = NewWebTongTongPay()
	PWebYunTongShangPay = NewWebYunTongShangPay()
	//SubGameInfoMgr = NewSubGameInfoManager()
	//初始化csv配置
	PCsvObjManager = for_game.NewCsvObjManager()
	PackageInitialized.Trigger()
	for_game.PDirtyWordsMgr = for_game.NewDirtyWordMgr()

	CronJob = easygo.Cronjob
}

// 依赖数据库的初始化
func InitializeDependDB() {
	for_game.InitIdGenerator(for_game.TABLE_PLAYER_ACCOUNT, for_game.INIT_PLAYER_ID)
	for_game.InitIdGenerator(for_game.TABLE_TEAM_DATA, for_game.INIT_TEAM_ID)
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_SOURCETYPE, "初始化金币源类型", false, bson.M{}, mongo_init.InitSourcetypeCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_SYS_PARAMETER, "初始化系统参数", false, bson.M{}, mongo_init.InitSysParameterCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_SYS_PARAMETER, "初始化系统参数推送配置", false, bson.M{"_id": for_game.PUSH_PARAMETER}, mongo_init.InitPushSetCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_RECHARGE, "初始化硬币配置", true, bson.M{}, mongo_init.InitCoinRechargeCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_PROPS_ITEM, "初始化道具配置", false, bson.M{}, mongo_init.InitPropsItemCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT, "初始化道具商品配置", true, bson.M{}, mongo_init.InitCoinProductCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTTYPE, "初始化标签类型", true, bson.M{}, mongo_init.InitInterestTagTypeCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTTAG, "初始化标签", true, bson.M{}, mongo_init.InitInterestTagCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_TYPE, "初始化话题类别", true, bson.M{}, mongo_init.InitTopicCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER_TYPES, "初始化客服类型", true, bson.M{}, mongo_init.InitManageTypes())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_MATCH_GUIDE, "初始化匹配引导语", false, bson.M{}, mongo_init.InitMatchGuideCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_SAY_HI, "初始化SayHi", false, bson.M{}, mongo_init.InitSayHiCfg())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_INTIMACY_COINFIG, "初始化亲密度", false, bson.M{}, mongo_init.InitIntimacyConfig())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_CHARACTER_TAG, "初始化个性化标签", false, bson.M{}, mongo_init.InitSysPersonalityTags())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_TAG, "初始系统背景标签", false, bson.M{}, mongo_init.InitSysBgVoiceTags())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_BG_IMAGE, "初始系统背景图", false, bson.M{}, mongo_init.InitSysBgImageTags())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_STARSIGNS_TAG, "初始化12星座标签", false, bson.M{}, mongo_init.InitStarSignsTag())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_STARSIGNS_TAG, "初始化12星座标签", false, bson.M{}, mongo_init.InitStarSignsTag())
	//for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_OPERATE, "初始化运营号", false, bson.M{}, mongo_init.InitOperatePlayer())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_INDEX_TIPS, "初始化首页Tips配置", true, bson.M{}, mongo_init.InitIndexTips())
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_POP_SUSPEND, "初始化弹窗悬浮球配置", true, bson.M{}, mongo_init.InitPopSuspend())

	//for_game.ChangePersonnalChat(MongoLogMgr)
	//for_game.ChangeTeamChat(MongoLogMgr)
	//=============初始化许愿池活动==========
	for_game.InitToMongo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_COIN_RECHARGE_ACT_CFG, "初始化充值活动配置", false, bson.M{}, mongo_init.InitWishCoinRechargeActivityCfg())

	for_game.InitGeneralQuota(mongo_init.InitPayGeneralQuota())    //初始化通用额度配置
	for_game.InitPayType(mongo_init.InitPayType())                 //初始化支付类型
	for_game.InitPayScene(mongo_init.InitPayScene())               //初始化支付场景
	for_game.InitPaySetting(mongo_init.InitPaymentSetting())       //初始化支付限定
	for_game.InitPaymentPlatform(mongo_init.InitPaymentPlatform()) //初始化支付平台
	for_game.InitPlatformChannel(mongo_init.InitPlatformChannel()) //初始化支付平台通道

	// ==========十一活动===========
	for_game.InitActivity(mongo_init.InitActivity())   //初始化活动
	for_game.InitProps(mongo_init.InitProps())         //初始化道具
	for_game.InitPropsRate(mongo_init.InitPropsRate()) //初始化抽卡概率
	for_game.InitPropsToMap()                          //初始化道具进内存,一定要现有道具和概率,否则初始化内存会失败.
	for_game.ReloadSysPropsToDayProps()                //加载掉落到数据库
	InitStarMap()                                      // 初始化星座匹配
	InitStarNameMap()                                  // 初始化星座id和名字
	TaskJobDayProps()                                  // 定时任务更新每天掉落
	TaskFullLuckyPlayer()                              // 定时加人数.
	TopicTeamUpdateDynamic()                           //话题群组定时发送动态
}

func InitStarMap() {
	//s := `[
	//		{
	//			"_id": 1,
	//			"Name": "自信满满白羊座"
	//		},
	//		{
	//			"_id": 2,
	//			"Name": "慢热达人金牛座"
	//		},
	//		{
	//			"_id": 3,
	//			"Name": "好奇宝宝双子座"
	//		},
	//		{
	//			"_id": 4,
	//			"Name": "温柔体贴巨蟹座"
	//		},
	//		{
	//			"_id": 5,
	//			"Name": "慷慨大方狮子座"
	//		},
	//		{
	//			"_id": 6,
	//			"Name": "完美主义处女座"
	//		},
	//		{
	//			"_id": 7,
	//			"Name": "世界和平天枰座"
	//		},
	//		{
	//			"_id": 8,
	//			"Name": "爱憎分明天蝎座"
	//		},
	//		{
	//			"_id": 9,
	//			"Name": "崇尚自由射手座"
	//		},
	//		{
	//			"_id": 10,
	//			"Name": "沉着冷静魔羯座"
	//		},
	//		{
	//			"_id": 11,
	//			"Name": "聪明过人水瓶座"
	//		},
	//		{
	//			"_id": 12,
	//			"Name": "梦幻达人双鱼座"
	//		}
	//	]`
	for_game.ConstellationMap = make(map[int32]int32)
	for_game.ConstellationMap[1] = 12
	for_game.ConstellationMap[2] = 11
	for_game.ConstellationMap[3] = 10
	for_game.ConstellationMap[4] = 9
	for_game.ConstellationMap[5] = 8
	for_game.ConstellationMap[6] = 7

	for_game.ConstellationMap[12] = 1
	for_game.ConstellationMap[11] = 2
	for_game.ConstellationMap[10] = 3
	for_game.ConstellationMap[9] = 4
	for_game.ConstellationMap[8] = 5
	for_game.ConstellationMap[7] = 6
}

func InitStarNameMap() {
	constellations := for_game.GetConfigConstellationFormDB()
	for_game.ConstellationNameMap = make(map[int32]string)
	for _, v := range constellations {
		for_game.ConstellationNameMap[v.GetId()] = v.GetName()
	}
}

//把需要初始化到reids的配置初始化
func InitCoinfigToRedis() {
	for_game.InitAllConfig()
}

func SaveRedisData() {
	logs.Info("停服保存数据redis数据存储到mongo----------->>>>>>>>>>>")
	//for_game.SaveCoinChangeLogToMongoDB()
	//for_game.SaveGoldChangeLogToMongoDB()
	//for_game.SaveRedisOrderToMongo()
	//for_game.SaveRedisAccountToMongo()
	//for_game.SaveRedisPlayerBagItemToMongo()
	//for_game.SaveRedisPlayerBaseToMongo()
	//for_game.SavePersonalChatToMongoDB()
	//for_game.SaveRedisPlayerEquipmentToMongo()
	//for_game.SaveRedisRedPacketToMongo()
	//for_game.SaveRedPacketTotalToMongoDB()
	//for_game.SaveRedisArticleReportToMongo()
	//for_game.SaveRedisInOutCashSumReportToMongo()
	//for_game.SaveRedisNoticeReportToMongo()
	//for_game.SaveRedisOperationChannelReportToMongo()
	//for_game.SaveRedisPlayerBehaviorReportToMongo()
	//for_game.SaveRedisPlayerKeepReportToMongo()
	//for_game.SaveRedisRecallReportToMongo()
	//for_game.SaveRedisRegisterLoginReportToMongo()
	//for_game.SaveRedisTeamToMongo()
	//for_game.SaveRedisTeamChatLogToMongo()
	//for_game.SaveRedisTeamPersonalToMongo()
	//for_game.SaveRedisTransferMoneyToMongo()
	//定时存储社区广场信息
	for_game.UpdateSquareDynamic()
	// 点赞的数据
	for_game.TaskSaveZanData()
	logs.Info("停服保存数据完成----------->>>>>>>>>>>")

}

func TaskJobDayProps() {
	CronJob.DayEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
		logs.Info("=================>执行定时任务充值背包记录")
		for_game.ReloadSysPropsToDayProps()
	})
}

// 每半小时增加集满人数.
func TaskFullLuckyPlayer() {
	t := time.Now().Unix()
	if for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_HALL_SID) == PServerInfo.GetSid() {
		if easygo.Str2Time(for_game.LUCKY_START_TIME).Unix() <= t && easygo.Str2Time(for_game.OPEN_BEFORE_TENTIME).Unix() >= t {
			for_game.AddFullLuckyPlayer()
		}
	}
	if easygo.Str2Time(for_game.OPEN_BEFORE_TENTIME).Unix() < t {
		logs.Info("================超出定时任务时间,直接return")
		return
	}
	easygo.AfterFunc(TIME_TO_ADD_FULL_PLAYER, TaskFullLuckyPlayer)
}

//每过一定时间群组发送动态()
func TopicTeamUpdateDynamic() {
	logs.Info("TopicTeamUpdateDynamic 话题群组定时发送动态-------->>>>>")
	teams := for_game.GetAllUnSendDynamicTopicTeam()
	for _, t := range teams {
		teamObj := for_game.GetRedisTeamObj(t.GetId(), t)
		if teamObj != nil {
			chat := teamObj.GetTopicTeamDynamic()
			if chat == nil {
				continue
			}
			player := GetPlayerObj(teamObj.GetTeamOwner())
			//随机延迟发
			rand := time.Duration(for_game.RandInt(60, 600))
			fun := func() {
				cl := &cls1{}
				cl.RpcChatNew(nil, player, chat)
				teamObj.SetLastDynamicTime(time.Now().Unix())
			}
			logs.Info("群id=%d将于%d秒后发送动态消息:%v", teamObj.GetId(), rand, chat)
			easygo.AfterFunc(rand*time.Second, fun)
		} else {
			logs.Error("teamobj is nil")
		}
	}
	easygo.AfterFunc(TIME_TO_SEND_TOPIC_DYNAMIC, TopicTeamUpdateDynamic)
}
