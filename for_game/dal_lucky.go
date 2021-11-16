package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"

	"github.com/astaxie/beego/logs"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

const (
	ACTIVITY_OPEN  = 0 //活动开启
	ACTIVITY_CLOSE = 1 //活动关闭

	ACTIVITY_TYPE_WISH = 3 // 许愿池活动类型

	TABLE_LUCKY_ACTIVITY             = "lucky_activity"             // 活动表
	TABLE_LUCKY_PROPS                = "lucky_props"                // 道具表
	TABLE_LUCKY_DAY_PROPS            = "lucky_day_props"            // 每天掉落
	TABLE_LUCKY_PLAYER_PROPS         = "lucky_player_props"         // 用户道具背包
	TABLE_LUCKY_PLAYER_USE_PROPS_LOG = "lucky_player_use_props_log" // 用户道具使用日志
	TABLE_LUCKY_PLAYER               = "lucky_player"               // 参加抽奖的玩家
	TABLE_LUCKY_PLAYER_RELATED       = "lucky_player_related"       // 玩家关联关系
	TABLE_LUCKY_PROPS_RATE           = "lucky_props_rate"           //抽卡概率
	TABLE_LUCKY_SYS_FULL_COUNT       = "lucky_sys_full_count"       //当前集满人数(系统自身生成的,不包含真实用户,如需真实用户,需要在LuckyPlayer表中累加).
	TABLE_ACTIVITY_REPORT            = "report_activity"            //活动统计表
)

// 卡片id常量
const (
	ID_HE   = 1 // 和
	ID_NING = 2 // 拧
	ID_MENG = 3 // 檬
	ID_QU   = 4 // 趣
	ID_LV   = 5 // 旅
	ID_XING = 6 // 行
)

// 卡片日志类型
const (
	LOG_GET     = 1 // 获得
	LOG_DEL     = 2 // 删除
	LOG_USE     = 3 // 使用
	LOG_GIVE    = 4 // 赠送
	LOG_RECEIVE = 5 // 获赠
)

//===================活动===========================
// 从数据库中获取活动详情.
func GetActivityFromDB(id PLAYER_ID) *share_message.Activity {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_ACTIVITY)
	defer closeFun()
	var activity *share_message.Activity
	err := col.Find(bson.M{"_id": id}).One(&activity)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		return nil
	}
	return activity
}

// 从数据库中获取活动详情.
func GetActivityByType(t int32) *share_message.Activity {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_ACTIVITY)
	defer closeFun()
	var activity *share_message.Activity
	err := col.Find(bson.M{"Types": t}).One(&activity)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		logs.Error("GetActivityByType err: %s", err.Error())
		return nil
	}
	return activity
}

// 从数据库中获取活动详情.
func GetActivityByTypes(t []int32) []*share_message.Activity {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_ACTIVITY)
	defer closeFun()
	activityList := make([]*share_message.Activity, 0)
	err := col.Find(bson.M{"Types": bson.M{"$in": t}}).All(&activityList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		logs.Error("GetActivityByTypes err: %s", err.Error())
		return nil
	}
	return activityList
}

//===================玩家道具背包===========================
// 从数据库中获取玩家道具列表.
func GetPlayerPropsListByPidFromDB(pid PLAYER_ID) []*share_message.PlayerProps {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER_PROPS)
	defer closeFun()
	playerPropsList := make([]*share_message.PlayerProps, 0)
	err := col.Find(bson.M{"PlayerId": pid}).All(&playerPropsList)
	easygo.PanicError(err)
	return playerPropsList
}
func GetPlayerPropsByPidAndProIdFromDB(pid, propsId PLAYER_ID) *share_message.PlayerProps {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER_PROPS)
	defer closeFun()
	playerProps := new(share_message.PlayerProps)
	err := col.Find(bson.M{"PlayerId": pid, "PropsId": propsId}).One(&playerProps)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return playerProps
}

// 玩家新增道具
func UpsetPlayerPropsToDB(playerId, propsId int64, count int) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER_PROPS)
	defer closeFun()

	_, err := col.Upsert(bson.M{"PlayerId": playerId, "PropsId": propsId}, bson.M{"$inc": bson.M{"Count": count}})
	easygo.PanicError(err)
	return err
}

//===================玩家道具背包===========================

//==================参加抽奖的玩家=======================
func UpsetLuckyPlayerToDB(pid, luckyCount, fullTime PLAYER_ID, isFull ...bool) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER)
	defer closeFun()
	upbson := bson.M{}
	isF := append(isFull, false)[0]
	if isF {
		count := GetLocalFullCount() + 1
		upbson = bson.M{"$inc": bson.M{"LuckyCount": luckyCount, "FullPlaces": count}, "$set": bson.M{"IsFull": isF, "FullTime": fullTime}}
	} else {
		upbson = bson.M{"$inc": bson.M{"LuckyCount": luckyCount}}
	}

	_, err := col.Upsert(bson.M{"_id": pid}, upbson)
	easygo.PanicError(err)
	return err
}

// 修改玩家是否是首次抽卡.
func UpsetIsNewLuckyToDB(pid PLAYER_ID, isNewLucky bool) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": pid}, bson.M{"$set": bson.M{"IsNewLucky": isNewLucky}})
	easygo.PanicError(err)
	return err
}

func GetLuckyPlayerFromDB(pid PLAYER_ID) *share_message.LuckyPlayer {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER)
	defer closeFun()
	lp := new(share_message.LuckyPlayer)
	err := col.Find(bson.M{"_id": pid}).One(&lp)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		panic(err)
	}
	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		return nil
	}
	return lp
}

// 获取集满的人物
func GetFullLuckyPlayerListFromDB(pageSize int) []*share_message.LuckyPlayer {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER)
	defer closeFun()
	lp := make([]*share_message.LuckyPlayer, 0)
	if err := col.Find(bson.M{"IsFull": true}).Limit(pageSize).Sort("-_id").All(&lp); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return lp
}

func UpsetLuckyMoneyToDB(lp *share_message.LuckyPlayer) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": lp.GetPlayerId()}, bson.M{"$set": lp})
	easygo.PanicError(err)
	return err
}

//==================参加抽奖的玩家=======================

//===========玩家关联关系======================
func UpsetLuckyPlayerRelatedToDB(pr *share_message.LuckyPlayerRelated) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER_RELATED)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": pr.GetId()}, bson.M{"$set": pr})
	return err
}

func GetRelatedByPhoneFromDB(friendPhone string) *share_message.LuckyPlayerRelated {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER_RELATED)
	defer closeFun()
	r := new(share_message.LuckyPlayerRelated)
	if err := col.Find(bson.M{"FriendPhone": friendPhone}).One(&r); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return r
}

//===========玩家关联关系======================

//=================系统道具================
func GetSysPropsListFromDB() []*share_message.Props {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PROPS)
	defer closeFun()
	ps := make([]*share_message.Props, 0)
	if err := col.Find(bson.M{}).All(&ps); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return ps
}

func GetSysPropsByIdFromDB(id int64) *share_message.Props {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PROPS)
	defer closeFun()
	ps := new(share_message.Props)
	err := col.Find(bson.M{"_id": id}).One(&ps)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return ps
}

//=================系统道具================

//=================概率表================
func GetSysPropsRateFromDB() []*share_message.PropsRate {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PROPS_RATE)
	defer closeFun()
	rate := make([]*share_message.PropsRate, 0)
	if err := col.Find(bson.M{}).All(&rate); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return rate
}

//=================概率表================

//==========掉落数据===================
func UpsetDayProps(dp *share_message.DayProps) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_DAY_PROPS)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": dp.GetId()}, bson.M{"$set": dp})
	easygo.PanicError(err)
}

func GetDayPropsById(id int64) *share_message.DayProps {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_DAY_PROPS)
	defer closeFun()
	props := new(share_message.DayProps)
	err := col.Find(bson.M{"_id": id}).One(&props)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return props
}

// 返回递增后的数量
func IncrDayProps(id, count int64) int64 {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_DAY_PROPS)
	defer closeFun()
	var dp *share_message.DayProps
	_, err := col.Find(bson.M{"_id": id}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"Count": count}},
		Upsert:    true,
		ReturnNew: true,
	}, &dp)
	easygo.PanicError(err)

	return dp.GetCount()

}

//==========掉落数据===================

//=======道具使用日志==========
func InsertUsePropsLogToDB(propsLogs []interface{}) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER_USE_PROPS_LOG)
	defer closeFun()
	err := col.Insert(propsLogs...)
	easygo.PanicError(err)
}

// 获取日志列表
func GetLogListByPid(pid PLAYER_ID, types int32) []*share_message.PlayerUsePropsLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_PLAYER_USE_PROPS_LOG)
	defer closeFun()
	logs := make([]*share_message.PlayerUsePropsLog, 0)
	if err := col.Find(bson.M{"PlayerId": pid, "Types": types}).Sort("-CreateTime").All(&logs); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return logs
}

//=======道具使用日志==========

//=======系统自己生成的当前集满人数==========
func UpsetSysFullCount(fullCount int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_SYS_FULL_COUNT)
	defer closeFun()
	_, err := col.Upsert(bson.M{}, bson.M{"$inc": bson.M{"FullCount": fullCount}})
	easygo.PanicError(err)
}

func GetFullCountFromDB() *share_message.LuckySysFullCount {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LUCKY_SYS_FULL_COUNT)
	defer closeFun()
	c := new(share_message.LuckySysFullCount)
	if err := col.Find(bson.M{}).One(&c); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return c
}

//=======系统自己生成的当前集满人数==========
//活动报表 fild字段名,val变化数量, createTime当天的报表数据更改可以不传时间
func UpdateActivityReport(fild string, val int, createTime ...int64) {
	cTime := easygo.GetToday0ClockTimestamp()
	if len(createTime) > 0 {
		cTime = easygo.Get0ClockTimestamp(createTime[0])
	}

	FindAndModify(MONGODB_NINGMENG, TABLE_ACTIVITY_REPORT, bson.M{"_id": cTime}, bson.M{"$inc": bson.M{fild: val}}, true)
}
