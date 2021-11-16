//初始化数据和定时任务入口
//
package backstage

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/mongo_init"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//初始化预设数据到数据库
func InitData() {
	SynchronizeOperationChannel()             //同步渠道表修改
	for_game.InitUsers(mongo_init.GetUsers()) //初始化后台管理员帐号
	//活动初始化
	for_game.InitActivity(mongo_init.InitActivity())   //初始化活动
	for_game.InitProps(mongo_init.InitProps())         //初始化道具
	for_game.InitPropsRate(mongo_init.InitPropsRate()) //初始化抽卡概率
	for_game.InitPropsToMap()                          // 初始化道具进内存,一定要现有道具和概率,否则初始化内存会失败.
}

//定时任务
func TimedTask() {
	sid := PServerInfo.GetSid()
	saveServerSid := for_game.GetCurrentSaveServerSid(sid, for_game.REDIS_SAVE_BACKSTAGE_SID)
	if saveServerSid == sid {
		//整点定时任务
		CronJob.HourEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
			if for_game.IS_FORMAL_SERVER {
				MakeReports() //整点生成报表
			}
			TimingCreatePlayer()
			for_game.UpdatePlayerOnlineReport(QueryPlayerOnline()) //更新在线玩家报表

			// 许愿池
			//easygo.Spawn(MakeWishReports)
		})
		//每天定时任务
		CronJob.DayEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
			// BakReport()
			SynchronizeOperationChannel() //同步渠道表修改
			MakePlayerBehaviorReportNeerbyData()
			OrderExpiredDay()
			UpdateVCBuryingPointReport(easygo.GetYesterday0ClockTimestamp())
			MakeReportsForDay()
			MakeWishReports() // 许愿池
		})
		//每整10分钟定时任务
		CronJob.TenMinEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
			if !for_game.IS_FORMAL_SERVER {
				MakeReports()       //测试服每10分钟更新报表数据
				MakeReportsForDay() //测试服每10分钟更新报表数据
				MakeWishReports()   // 许愿池
			}
			CheckWaiterCount() //每10分钟检查客服活跃消息数量并更新
		})

		/*easygo.Spawn(MakeWishBoxReportsOfWeek)
		easygo.Spawn(MakeWishBoxReportsOfMonth)*/

		//每周定时任务
		CronJob.WeekEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
			MakeEsportPointReport("week")
			//每周周一生成一次
			easygo.Spawn(MakeWishBoxReportsOfWeek)
		})
		//每月定时任务
		CronJob.MonthEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
			MakeEsportPointReport("month")
			easygo.Spawn(MakeWishBoxReportsOfMonth)
		})
		//每天5点定时任务
		CronJob.DayFiveClockEvent.AddHandler(func(year int, month int, day int, hour int, Minute int, week int) {
			MakeEsportPointReport("day")
		})
		//检查并更新定时任务管理器
		TimedAppPushMessage()
		//小助手文章定时管理器
		TimedAppPushTweets()
		//启服加载社交动态定时任务
		TimedSendDynamic()
		//启服加载电竞定时任务
		TimeSendEsports()
	}
}

//=====测试和预留函数========================================================================================================>

//临时使用,批量修改IsFollow值
func CheckFollow(site SITE) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"$or": []bson.M{{"IsFollow": true}, {"IsFollow": false}, {"IsFollow": bson.M{"$eq": nil}}}}, bson.M{"$set": bson.M{"IsFollow": 0}})
	easygo.PanicError(err)
}

//数据库自定义函数，存储过程调用
func GetPlayerName(id int64) string {
	col, closeFun := MongoMgr.GetDB(for_game.MONGODB_NINGMENG)
	defer closeFun()
	result := bson.M{}
	av := fmt.Sprintf("getPlayerName(%d)", id) //获取自增id的自定义函数
	logs.Info(av)
	err := col.Run(bson.M{"eval": av}, &result)
	easygo.PanicError(err)
	return easygo.AnytoA(result["retval"])
}

//mongodb迭代器 游标
func GetPlayerNickName() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	player := &share_message.PlayerBase{}
	iter := col.Find(bson.M{}).Iter() //mongodb迭代器
	for iter.Next(&player) {
		logs.Info("======", player.GetNickName())
	}

	// pli := []*share_message.PlayerBase{}
	// err := iter.All(&pli)
	// easygo.PanicError(err)
	// logs.Info("======", pli)
	if err := iter.Close(); err != nil {
		easygo.PanicError(err)
	}
}

//查询数据库并更新
func DbUpdate() {
	// var a *share_message.DataArea
	find := bson.M{"_id": "华东"}
	update := bson.M{"$set": bson.M{"Country": "上海"}}
	one := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_DATA_AREA, find, update, true)
	logs.Info(one)
}

//获取B站评论
func GetBstationContent(av, page int) {
	msg := for_game.GetBstationContent("https://api.bilibili.com/x/v2/reply?jsonp=jsonp&pn=" + easygo.AnytoA(page) + "&type=1&oid=" + easygo.AnytoA(av) + "&sort=1&nohot=1")
	if len(msg.Replies) == 0 {
		av -= 1
	} else {
		page += 1
	}
	logs.Info(msg.Replies, len(msg.Replies))
	//easygo.AfterFunc(5*time.Second, func() {
	//	GetBstationContent(av, page)
	//})
}

//获取坐标附近的人 需要2d索引
func GetNear() {
	col, closeFun := MongoMgr.GetDB(for_game.MONGODB_NINGMENG)
	defer closeFun()
	var result bson.M

	// result := make([]*share_message.PlayerBase, 0)
	av := "getNear()" //获取自增id的自定义函数
	logs.Info(av)
	err := col.Run(bson.M{"eval": av}, &result)
	easygo.PanicError(err)

	result1 := result["retval"].(bson.M)["_batch"]
	var li []*share_message.PlayerBase
	// logs.Debug("----------", reflect.TypeOf(result1))
	b, _ := json.Marshal(result1)
	// logs.Debug("----------", string(b))
	json.Unmarshal(b, &li)
	logs.Debug("===========", li[0].GetAccount())
}

//复杂聚合查询
func GetMapReduce() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT_REPORT)
	defer closeFun()

	job := &mgo.MapReduce{
		Map:     "function() { emit(this.CreateTime, this.CoinSum) }",
		Reduce:  "function(key, values) { return Array.sum(values) }",
		Verbose: true,
	}
	var result []struct {
		Id    int `json:"_id"`
		Value int
	}
	info, err := col.Find(bson.M{}).MapReduce(job, &result)

	easygo.PanicError(err)
	logs.Debug("result", result, "info", info)

}
