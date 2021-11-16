package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo/bson"
)

//写后台操作日志的对外方法
func AddBackstageLog(account USER_ACCOUNT, ip string, optType string, remark string) {
	defer func() {
		site := for_game.MONGODB_NINGMENG
		htLog := map[string]interface{}{
			"site":    site,
			"account": account,
			"ip":      ip,
			"optType": optType,
			"remark":  remark,
		}

		AddBackstageOptLog(htLog)
	}()

}

//写后台的操作日志
func AddBackstageOptLog(logMap map[string]interface{}) {
	site := for_game.MONGODB_NINGMENG
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_HTOPTLOG_NAME)
	defer closeFun()
	htLog := &share_message.BackstageOptLog{
		Site:       easygo.NewString(site),
		Account:    easygo.NewString(logMap["account"]),
		Ip:         easygo.NewString(logMap["ip"]),
		OptType:    easygo.NewString(logMap["optType"]),
		Remarks:    easygo.NewString(logMap["remark"]),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
	}
	targetUserId, ok := logMap["targetUserId"]
	if ok && targetUserId != nil {
		tUserId := []int64{}
		switch tt := targetUserId.(type) {
		case int:
			tUserId = append(tUserId, targetUserId.(int64))
		case int64, int32:
			tUserId = append(tUserId, targetUserId.(int64))
		case []int64:
			tUserId = targetUserId.([]int64)
		default:
			errStr := fmt.Sprintf("添加操作记录错误，不支持被操作人的类型: %T , 内容: %v", tt, targetUserId)
			fmt.Println(errStr)
			htLog.Remarks = easygo.NewString(errStr)
		}

		htLog.TargetUserId = tUserId
	}

	err := col.Insert(htLog)
	easygo.PanicError(err)
}

// 查询管理员日志
func QueryBackstageLog(reqMsg *brower_backstage.ListRequest) ([]*share_message.BackstageOptLog, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_HTOPTLOG_NAME)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.BeginTimestamp != nil {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	switch reqMsg.GetType() {
	case 1:
		if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
			queryBson["Account"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.SrtType != nil && reqMsg.GetSrtType() != "" {
		queryBson["OptType"] = reqMsg.GetSrtType()
	}

	var list []*share_message.BackstageOptLog
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//金币日志的通用查询
func GetGoldLogList(reqMsg *brower_backstage.QueryGoldLogRequest) ([]*for_game.GoldChangeLog, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_GOLDCHANGELOG)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	//不查关键词
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetKeyType() {
		case 1: //订单号
			queryBson["Extend.OrderId"] = reqMsg.GetKeyword()
		case 2: //柠檬号
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			queryBson["PlayerId"] = player.GetPlayerId()
		default:
			easygo.NewFailMsg("查询条件有误")
		}
	}

	if reqMsg.SourceType != nil && len(reqMsg.GetSourceType()) > 0 {
		types := reqMsg.GetSourceType()
		queryBson["SourceType"] = bson.M{"$in": types}

	}

	if reqMsg.PayType != nil && reqMsg.GetPayType() != 0 {
		queryBson["PayType"] = reqMsg.GetPayType()
	}

	// keys := []string{"PayType"}

	// for_game.EnsureIndexKey(for_game.MONGODB_NINGMENG, for_game.TABLE_GOLDCHANGELOG, keys)

	query := col.Find(queryBson).Hint("_id")
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*for_game.GoldChangeLog
	errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//查询指定时间的用户在线时长日志
func GetOnlineTimeLogList(startTime, endTime int64) []*share_message.OnlineTimeLog {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ONLINETIMELOG)
	defer closeFun()

	queryBson := bson.M{"CreateTime": bson.M{"$gte": startTime, "$lte": endTime}}
	query := col.Find(queryBson)
	var list []*share_message.OnlineTimeLog
	errc := query.All(&list)
	easygo.PanicError(errc)

	return list
}
