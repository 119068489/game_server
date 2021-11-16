package backstage

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

//查询埋点报表注册人数和登录人数汇总
func GetRegAndLogSumCount(user *share_message.Manager) (int64, int64) {
	// col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, "bak_report_register_login")
	// if user.GetRole() == 0 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_LOGIN_REGISTER_REPORT)
	// }
	defer closeFun()

	queryBson := bson.M{"_id": bson.M{"$lte": easygo.GetToday0ClockTimestamp()}} //查询今天以前的所有数据
	//聚合查询
	m := []bson.M{
		{"$match": queryBson},
		{"$group": bson.M{"_id": nil, "RegSumCount": bson.M{"$sum": "$RegSumCount"}, "LoginSumCount": bson.M{"$sum": "$LoginSumCount"}}},
	}

	query := col.Pipe(m)
	list := &share_message.RegisterLoginReport{}
	err := query.One(list)

	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return 0, 0
	}

	return list.GetRegSumCount(), list.GetLoginSumCount()
}

//兴趣爱好分类查询玩家数量
func GetInterestTagSumCount(tag []int32) []*share_message.PipeIntCount {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	queryBson := bson.M{"Label": bson.M{"$in": tag}} //查询指定的标签用户
	//聚合查询
	m := []bson.M{
		{"$unwind": "$Label"},
		{"$match": queryBson},
		{"$group": bson.M{"_id": "$Label", "Count": bson.M{"$sum": 1}}},
	}

	query := col.Pipe(m)
	var all []*share_message.PipeIntCount
	err := query.All(&all)
	easygo.PanicError(err)

	return all
}

//手机品牌玩家数量
func GetPhoneBrandSumCount(brand string) int64 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	queryBson := bson.M{"Brand": bson.M{"$regex": bson.RegEx{Pattern: brand, Options: "i"}}} //忽略大小写模糊匹配手机品牌
	//聚合查询
	m := []bson.M{
		{"$match": queryBson},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}

	query := col.Pipe(m)
	var one *for_game.PipeIntCount
	err := query.One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return 0
	}

	return *one.Count
}

//上网热度柱状图
func GetPlayerOnlineLineReport(reqMsg *brower_backstage.ListRequest) []*share_message.PlayerOnlineReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYERONLINE_REPORT)
	defer closeFun()

	queryBson := bson.M{}
	if reqMsg.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}
	query := col.Find(queryBson)
	var list []*share_message.PlayerOnlineReport
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

//用户登录地区分布图
func GetPlayerLogLocation(reqMsg *brower_backstage.ListRequest) []*for_game.PipeStringCount {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYERLOG_LOCATION_REPORT)
	defer closeFun()

	dayTime := easygo.Get0ClockTimestamp(easygo.NowTimestamp() - 86400) //查询前一天的报表数据
	queryBson := bson.M{"DayTime": dayTime}
	if reqMsg.GetListType() > 0 {
		queryBson["DeviceType"] = reqMsg.GetListType()
	}

	switch reqMsg.GetType() {
	case 1:
		queryBson["Piece"] = "CN"
	case 2:
		queryBson["Piece"] = bson.M{"$ne": "CN"}
	default:
		queryBson["Piece"] = bson.M{"$ne": nil}
	}

	//聚合查询
	m := []bson.M{
		{"$match": queryBson},
		{"$group": bson.M{"_id": "$Position", "Count": bson.M{"$sum": "$Count"}}},
	}

	query := col.Pipe(m)
	var list []*for_game.PipeStringCount
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

//注册用户区域分布数据
func GetPlayerRegLocation() []*for_game.PipeStringCount {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	queryBson := bson.M{"Provice": bson.M{"$ne": ""}}
	//聚合查询
	m := []bson.M{
		{"$match": queryBson},
		{"$group": bson.M{"_id": "$Provice", "Count": bson.M{"$sum": 1}}},
	}

	query := col.Pipe(m)
	var list []*for_game.PipeStringCount
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

//玩家男女分布数据
func GetPlayerGenderCount() []*for_game.PipeIntCount {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	queryBson := bson.M{"Sex": bson.M{"$gt": 0}}
	//聚合查询
	m := []bson.M{
		{"$match": queryBson},
		{"$group": bson.M{"_id": "$Sex", "Count": bson.M{"$sum": 1}}},
	}

	query := col.Pipe(m)
	var list []*for_game.PipeIntCount
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}
