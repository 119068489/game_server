package for_game

//广告数据查询
import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"

	"github.com/astaxie/beego/logs"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

//广告数据类型
type (
	ADV_ID       = int64   //id类型
	ADV_IDS      = []int64 //id数组类型
	ADV_LOCATION = int32   //投放类型
	ADV_TYPES    = int32   //素材类型
	ADV_STATUS   = int32   //状态类型
)

const (
	//Location广告投放的位置类型
	ADV_LOCATION_SQUARE           = 1 //社交广场动态流广告
	ADV_LOCATION_START            = 2 //启动页广告
	ADV_LOCATION_BANNER_SQUARE    = 3 //banner位横幅广场页广告
	ADV_LOCATION_BANNER_PERS      = 4 //banner位横幅个人页广告
	ADV_LOCATION_BANNER_COIN      = 5 //banner硬币页
	ADV_LOCATION_BANNER_MSG_LEFT  = 6 //banner消息页广告左侧
	ADV_LOCATION_VOICE_LOVE       = 7 //恋爱匹配信息流广告
	ADV_LOCATION_BANNER_MSG_RIGHT = 8 //banner消息页广告右侧
	//广告状态类型
	ADV_ON_SHELF  = 1 //广告状态上架
	ADV_OFF_SHELF = 2 //广告状态下架

	//跳转对象
	ADV_JUMP_OBJ_LINK   = 1 //外链
	ADV_JUMP_OBJ_SQUARE = 2 //社交广场
	ADV_JUMP_OBJ_TEAM   = 3 //群

)

// 查询广告
func QueryAdvToDB(id ADV_ID) *share_message.AdvSetting {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ADV_DATA)
	defer closeFun()
	siteOne := &share_message.AdvSetting{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

// 查询广告
func QueryAdvByLacSort(location ADV_LOCATION, sort int32) *share_message.AdvSetting {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ADV_DATA)
	defer closeFun()
	siteOne := &share_message.AdvSetting{}
	err := col.Find(bson.M{"Location": location, "Weights": sort}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//插入更新广告数据
func UpdateAdvListToDB(reqMsg *share_message.AdvSetting) {
	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(NextId(TABLE_ADV_DATA))
		reqMsg.CreateTime = easygo.NewInt64(util.GetMilliTime())
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ADV_DATA)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//批量广告上下架
func UpdateAdvShelfToDB(ids ADV_IDS, opt ADV_STATUS) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ADV_DATA)
	defer closeFun()

	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"Status": opt}})
	easygo.PanicError(err)
}

//按权重排序查询需投放的广告列表  Location广告投放的位置类型：1 社交动态列表广告,2启动页广告,3 banner位横幅广场页广告,4 banner位横幅个人页广告 5 banner硬币页，6banner消息页，7恋爱匹配广告
func QueryAdvListToDB(location ADV_LOCATION, ip ...string) []*share_message.AdvSetting {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ADV_DATA)
	defer closeFun()
	nowTime := util.GetMilliTime()
	queryBson := bson.M{"Status": ADV_ON_SHELF, "StartTime": bson.M{"$lt": nowTime}, "EndTime": bson.M{"$gt": nowTime}, "Location": location}
	if len(ip) > 0 && IsShieldUser(ip[0], false) {
		queryBson["IsShield"] = bson.M{"$ne": true}
	}
	var list []*share_message.AdvSetting
	query := col.Find(queryBson)
	err := query.Sort("-Weights").All(&list)
	easygo.PanicError(err)

	return list
}

func QueryMsgPageAdvToDB(ip ...string) ([]*share_message.AdvSetting, []*share_message.AdvSetting) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ADV_DATA)
	defer closeFun()
	nowTime := util.GetMilliTime()
	leftBson := bson.M{"Status": ADV_ON_SHELF, "StartTime": bson.M{"$lt": nowTime}, "EndTime": bson.M{"$gt": nowTime}, "Location": ADV_LOCATION_BANNER_MSG_LEFT}
	if len(ip) > 0 && IsShieldUser(ip[0], false) {
		leftBson["IsShield"] = bson.M{"$ne": true}
	}
	leftM := []bson.M{
		{"$match": leftBson},
		{"$sample": bson.M{"size": 1}},
	}
	var leftList []*share_message.AdvSetting
	err := col.Pipe(leftM).All(&leftList)
	easygo.PanicError(err)

	var rightList []*share_message.AdvSetting
	rightBson := bson.M{"Status": ADV_ON_SHELF, "StartTime": bson.M{"$lt": nowTime}, "EndTime": bson.M{"$gt": nowTime}, "Location": ADV_LOCATION_BANNER_MSG_RIGHT}
	if len(ip) > 0 && IsShieldUser(ip[0], false) {
		rightBson["IsShield"] = bson.M{"$ne": true}
	}
	rightM := []bson.M{
		{"$match": rightBson},
		{"$sample": bson.M{"size": 1}},
	}
	err = col.Pipe(rightM).All(&rightList)
	easygo.PanicError(err)

	return leftList, rightList
}

//修改广告排序
func UpdateAdvSortToDB(id ADV_ID, weights int32) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ADV_DATA)
	defer closeFun()
	upBson := bson.M{"$set": bson.M{"Weights": weights}}
	err := col.Update(bson.M{"_id": id}, upBson)
	easygo.PanicError(err)
}

//拖拽修改广告排序
func UpdateAdvSort(ids ADV_IDS) {
	w := int32(len(ids))
	for _, i := range ids {
		UpdateAdvSortToDB(i, w)
		w--
	}
}

//==================广告埋点数据============================
func AddAdvLogToDB(log *share_message.AdvLogReq) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_ADV_LOG)
	defer closeFun()
	err := col.Insert(log)
	easygo.PanicError(err)

}

// 根据操作判断是否有埋点数据
func GetAdvLogByPidAndOpFromDB(pid, advId PLAYER_ID, opType int32) []*share_message.AdvLogReq {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_ADV_LOG)
	defer closeFun()
	advLog := make([]*share_message.AdvLogReq, 0)
	err := col.Find(bson.M{"PlayerId": pid, "OpType": opType, "AdvId": advId}).Sort("-OpTime").All(&advLog)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	if len(advLog) == 0 {
		return nil
	}
	// 判断操作时间是否已是今天.
	logs.Info("opType : %d,数据库0点时间为: %d, 当前0点时间 %d", opType, easygo.Get0ClockTimestamp(advLog[0].GetOpTime()), easygo.GetToday0ClockTimestamp())
	if easygo.Get0ClockTimestamp(advLog[0].GetOpTime()) == easygo.GetToday0ClockTimestamp() {
		return advLog
	}
	return nil
}

//==================广告埋点数据============================

//==================附近的人广告埋点数据============================
func AddNearAdvLogToDB(log *share_message.AdvLogReq) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_ADV_LOG)
	defer closeFun()
	err := col.Insert(log)
	easygo.PanicError(err)

}

// 根据操作判断是否有埋点数据
func GetNearAdvLogByPidAndOpFromDB(pid, advId PLAYER_ID, opType int32) []*share_message.AdvLogReq {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_ADV_LOG)
	defer closeFun()
	advLog := make([]*share_message.AdvLogReq, 0)
	err := col.Find(bson.M{"PlayerId": pid, "OpType": opType, "AdvId": advId}).Sort("-OpTime").All(&advLog)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	if len(advLog) == 0 {
		return nil
	}
	// 判断操作时间是否已是今天.
	logs.Info("opType : %d,数据库0点时间为: %d, 当前0点时间 %d", opType, easygo.Get0ClockTimestamp(advLog[0].GetOpTime()), easygo.GetToday0ClockTimestamp())
	if easygo.Get0ClockTimestamp(advLog[0].GetOpTime()) == easygo.GetToday0ClockTimestamp() {
		return advLog
	}
	return nil
}

//==================附近的人广告埋点数据============================
//通过广告id获取所有广告信息
func GetAllAdvsByIds(ids []int64, ip ...string) map[int64]*share_message.AdvSetting {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ADV_DATA)
	defer closeFun()
	nowTime := util.GetMilliTime()
	queryBson := bson.M{"Status": ADV_ON_SHELF, "StartTime": bson.M{"$lt": nowTime}, "EndTime": bson.M{"$gt": nowTime}, "_id": bson.M{"$in": ids}}
	if len(ip) > 0 && IsShieldUser(ip[0], false) {
		queryBson["IsShield"] = bson.M{"$ne": true}
	}
	var list []*share_message.AdvSetting
	query := col.Find(queryBson)
	err := query.Sort("-Weights").All(&list)
	easygo.PanicError(err)
	mAdvs := make(map[int64]*share_message.AdvSetting)
	for _, l := range list {
		mAdvs[l.GetId()] = l
	}
	return mAdvs
}

//随机获取指定条数广告列表  Location广告投放的位置类型：1 社交动态列表广告,2启动页广告,3 banner位横幅广场页广告,4 banner位横幅个人页广告 5 banner硬币页，6banner消息页，7恋爱匹配广告
func QueryRandAdvListToDB(location ADV_LOCATION, num int32, ip ...string) []*share_message.AdvSetting {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ADV_DATA)
	defer closeFun()
	nowTime := util.GetMilliTime()
	queryBson := bson.M{"Status": ADV_ON_SHELF, "StartTime": bson.M{"$lt": nowTime}, "EndTime": bson.M{"$gt": nowTime}, "Location": location}
	if len(ip) > 0 && IsShieldUser(ip[0], false) {
		queryBson["IsShield"] = bson.M{"$ne": true}
	}
	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}

	var list []*share_message.AdvSetting
	err := col.Pipe(m).All(&list)
	easygo.PanicError(err)
	return list
}
