package backstage

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/h5_wish"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"strings"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/astaxie/beego/logs"

	"github.com/akqp2019/mgo/bson"
)

// slice 分页
func PaginateForInt64(x []int64, skip int, size int) []int64 {
	limit := func() int {
		if skip+size > len(x) {
			return len(x)
		} else {
			return skip + size
		}

	}

	start := func() int {
		if skip > len(x) {
			return len(x)
		} else {
			return skip
		}

	}
	return x[start():limit()]
}

// 获取活动状态
func GetActCfgIsOnline(t int32) bool {
	cfg := GetWishActCfg(t)
	curTime := time.Now().Unix()
	if curTime >= cfg.GetStartTime() && curTime <= cfg.GetEndTime() && cfg.GetStatus() == 0 {
		return true
	}

	return false
}

// 获取全部活动状态
func GetAllActCfgIsOnline() bool {
	cfg1 := GetActCfgIsOnline(for_game.WISH_ACT_COUNT)
	cfg2 := GetActCfgIsOnline(for_game.WISH_ACT_DAY)
	cfg3 := GetActCfgIsOnline(for_game.WISH_ACT_WEEK_MONTH)
	if cfg1 || cfg2 || cfg3 {
		return true
	}

	return false
}

// 获取活动奖池管理
func QueryWishActPool(req *brower_backstage.ListRequest) ([]*share_message.WishActPool, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-CreateTime"}

	queryBson := bson.M{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishActPool
	var errc error

	errc = query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 获取活动奖池管理
func GetWishActPoolAll() []*share_message.WishActPool {
	sort := []string{"-CreateTime"}

	queryBson := bson.M{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL)
	defer closeFun()

	query := col.Find(queryBson)

	var list []*share_message.WishActPool
	var errc error

	errc = query.Sort(sort...).All(&list)
	easygo.PanicError(errc)

	return list
}

// 新增/更新活动奖池
func UpdateWishActPool(upData *share_message.WishActPool) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL)
	defer closeFun()

	curTime := time.Now().Unix()
	// 新增
	if upData.GetId() == 0 {
		upData.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_ACT_POOL))
		upData.CreateTime = easygo.NewInt64(curTime)
	}

	//upData.UpdateTime = easygo.NewInt64(curTime)

	_, err1 := col.Upsert(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err1)
}

// 删除活动奖池
func DeleteWishActPool(ids []int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL)
	defer closeFun()

	_, err1 := col.RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	easygo.PanicError(err1)
}

// 删除活动奖池
func DeleteWishActPoolRule(ids []int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()

	_, err1 := col.RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	easygo.PanicError(err1)
}

// 删除活动奖池
func DeleteWishActPoolRuleByPools(ids []int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()

	_, err1 := col.RemoveAll(bson.M{"WishActPoolId": bson.M{"$in": ids}})
	easygo.PanicError(err1)
}

//获取活动奖池
func GetWishActPool(id int64) *share_message.WishActPool {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL)
	defer closeFun()
	player := &share_message.WishActPool{}
	err := col.Find(bson.M{"_id": id}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

//获取活动奖池
func GetWishActPoolByBoxId(id int64) *share_message.WishActPool {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL)
	defer closeFun()
	player := &share_message.WishActPool{}
	err := col.Find(bson.M{"BoxIds": id}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

// 获取活动设置
func GetWishActCfg(t int32) *share_message.Activity {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_LUCKY_ACTIVITY)
	defer closeFun()

	list := &share_message.Activity{}
	queryBson := bson.M{"Types": t}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return list
	}

	return list
}

func GetWishActCfgs(t []int32) []*share_message.Activity {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_LUCKY_ACTIVITY)
	defer closeFun()

	var list []*share_message.Activity
	queryBson := bson.M{"Types": bson.M{"$in": t}}
	err := col.Find(queryBson).All(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return list
	}

	return list
}

// 更新活动设置
func UpdateWishActCfg(data *share_message.Activity) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_LUCKY_ACTIVITY)
	defer closeFun()

	if data.GetId() == 0 {
		data.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_LUCKY_ACTIVITY))
	}

	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

// 获取活动奖池管理
func QueryWishActPoolRule(req *brower_backstage.ListRequest, t int32) ([]*share_message.WishActPoolRule, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-CreateTime"}

	queryBson := bson.M{}
	if t != 0 {
		queryBson["Type"] = t
	}

	// 状态
	if req.ListType != nil && req.GetListType() != 0 {
		queryBson["WishActPoolId"] = req.GetListType()
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishActPoolRule
	var errc error

	errc = query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 获取活动奖池管理
func QueryWishActPoolRuleWeekMonth(req *brower_backstage.ListRequest, t int32) ([]*share_message.WishActPoolRule, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"Key"}

	queryBson := bson.M{}
	if t != 0 {
		queryBson["Type"] = t
	} else {
		queryBson["Type"] = bson.M{"$in": []int32{3, 4}}
	}

	// 状态
	if req.ListType != nil && req.GetListType() != 0 {
		queryBson["WishActPoolId"] = req.GetListType()
	}

	// 奖励类型
	if req.Type != nil && req.GetType() != 0 {
		queryBson["AwardType"] = req.GetType()
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishActPoolRule
	var errc error
	if pageSize > 0 {
		errc = query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	} else {
		errc = query.Sort(sort...).All(&list)
	}

	easygo.PanicError(errc)

	return list, int32(count)
}

// 新增累计规则
func UpdateWishActPoolRule(data *share_message.WishActPoolRule) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()

	if data.GetId() == 0 {
		data.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_ACT_POOL_RULE))
	}

	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

// 获取累计规则
func GetWishActPoolRule(pid int64) *share_message.WishActPoolRule {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()

	list := &share_message.WishActPoolRule{}
	queryBson := bson.M{"_id": pid}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	return list
}

// 获取累计规则
func GetWishActPoolRuleByPoolIdAndKey(pid int64, key, t int32) *share_message.WishActPoolRule {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()

	list := &share_message.WishActPoolRule{}
	queryBson := bson.M{"WishActPoolId": pid, "Key": key, "Type": t}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	return list
}

// 获取活动期间抽奖记录列表
func QueryActDrawRecordList(req *brower_backstage.ListRequest) *brower_backstage.WishActPlayerRecordList {

	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG)
	defer closeFun()

	match := bson.M{}

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 柠檬号
			user := for_game.GetPlayerByAccount(req.GetKeyword())
			wishP := GetWishPlayerInfoByPid(user.GetPlayerId())
			match["DareId"] = wishP.GetId()
		case 2: // 用户昵称
			user := for_game.GetPlayerByNickName(req.GetKeyword())
			wishP := GetWishPlayerInfoByPid(user.GetPlayerId())
			match["DareId"] = wishP.GetId()
		case 3: // 手机号码
			user := for_game.GetPlayerByPhone(req.GetKeyword())
			wishP := GetWishPlayerInfoByPid(user.GetPlayerId())
			match["DareId"] = wishP.GetId()
		default:

		}
	}

	if req.GetBeginTimestamp() != 0 {
		match["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	} else {
		// 默认查询活动时间期间的
		cfg := GetWishActCfg(for_game.WISH_ACT_COUNT)
		match["CreateTime"] = bson.M{"$gte": cfg.GetStartTime(), "$lte": cfg.GetEndTime()}
	}

	group := bson.M{
		"_id":              "$DareId",
		"DrawDiamondTotal": bson.M{"$sum": "$DarePrice"},
		"LastDrawTime":     bson.M{"$last": "$CreateTime"},
		"DrawTotal":        bson.M{"$sum": 1},
		"UserId":           bson.M{"$last": "$DareId"},
	}

	project := bson.M{
		"UserId":           1,
		"DrawDiamondTotal": 1,
		"LastDrawTime":     1,
		"DrawTotal":        1,
	}

	countType := "PageCount"

	countCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$count": countType},
	}

	retList := &brower_backstage.WishActPlayerRecordList{}
	err := col.Pipe(countCond).One(&retList)

	sort := bson.M{"LastDrawTime": -1}
	if req.GetSort() != "" {
		if strings.Contains(req.GetSort(), "-") {
			s := strings.TrimLeft(req.GetSort(), "-")
			sort = bson.M{s: -1}
		} else {
			sort = bson.M{req.GetSort(): 1}
		}
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
		{"$sort": sort},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	var list []*brower_backstage.WishActPlayerRecord
	err = col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)

	retList.List = list
	return retList
}

// 获取活动期间用户信息列表
func QueryActPlayerInfoList(req *brower_backstage.ListRequest) ([]*share_message.WishPlayerActivity, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	match := bson.M{} //"_id": req.GetId()
	if req.GetBeginTimestamp() != 0 {
		match["UpdateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	} else {
		// 默认查询活动时间期间的
		cfg := GetWishActCfg(for_game.WISH_ACT_COUNT)
		startTime := cfg.GetStartTime()
		endTime := cfg.GetEndTime()
		if startTime < 9999999999 {
			startTime = startTime * 1000
			endTime = endTime * 1000
		}
		match["UpdateTime"] = bson.M{"$gte": startTime, "$lte": endTime}
	}

	if req.GetKeyword() != "" {
		playerBase := QueryPlayerByAccountOrPhone(req.GetKeyword())
		wishPlayer := GetWishPlayerInfoByPid(playerBase.GetPlayerId())
		match["_id"] = wishPlayer.GetId()
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER_ACTIVITY)
	defer closeFun()

	query := col.Find(match)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishPlayerActivity
	errc := query.Sort("-UpdateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 根据用户id获取活动领奖次数
func GetActReceiveAwardCountByUserIds(id []int64) []*brower_backstage.WishActPlayerRecord {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACTIVITY_PRIZE_LOG)
	defer closeFun()

	match := bson.M{"PlayerId": bson.M{"$in": id}}
	match["Status"] = 1
	// 默认查询活动时间期间的
	cfg := GetWishActCfg(for_game.WISH_ACT_COUNT)
	startTime := cfg.GetStartTime()
	endTime := cfg.GetEndTime()
	if startTime < 9999999999 {
		startTime = startTime * 1000
		endTime = endTime * 1000
	}
	match["FinishTime"] = bson.M{"$gte": startTime, "$lte": endTime}
	group := bson.M{
		"_id":        "$PlayerId",
		"AwardTotal": bson.M{"$sum": 1},
		"UserId":     bson.M{"$last": "$PlayerId"},
	}

	project := bson.M{
		"UserId":     1,
		"AwardTotal": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var counts []*brower_backstage.WishActPlayerRecord
	err := col.Pipe(pipeCond).All(&counts)
	easygo.PanicError(err)

	return counts
}

// 根据用户id获取活动领奖次数
func GetActDayCountByUserIds(id []int64) []*brower_backstage.WishActPlayerRecord {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER_ACTIVITY)
	defer closeFun()

	match := bson.M{"PlayerId": bson.M{"$in": id}}
	match["Status"] = 1
	match["Data.Type"] = 1
	// 默认查询活动时间期间的
	//cfg := GetWishActCfg()
	//match["FinishTime"] = bson.M{"$gte": cfg.GetStartTime(), "$lte": cfg.GetEndTime()}
	group := bson.M{
		"_id":        "$PlayerId",
		"AwardTotal": bson.M{"$sum": 1},
		"UserId":     bson.M{"$last": "$PlayerId"},
	}

	project := bson.M{
		"UserId":     1,
		"AwardTotal": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var counts []*brower_backstage.WishActPlayerRecord
	err := col.Pipe(pipeCond).All(&counts)
	easygo.PanicError(err)

	return counts
}

// 获取活动期间获奖记录列表
func QueryActWinRecordList(req *brower_backstage.ListRequest) ([]*share_message.WishActivityPrizeLog, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	match := bson.M{"PlayerId": req.GetId()}
	if req.GetBeginTimestamp() != 0 {
		match["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	} else {
		// 默认查询活动时间期间的
		cfg := GetWishActCfg(for_game.WISH_ACT_COUNT)
		startTime := cfg.GetStartTime()
		endTime := cfg.GetEndTime()
		if startTime < 9999999999 {
			startTime = startTime * 1000
			endTime = endTime * 1000
		}
		match["CreateTime"] = bson.M{"$gte": startTime, "$lte": endTime}
	}

	if req.ListType != nil && req.GetListType() != 1000 {
		match["Status"] = req.GetListType()
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACTIVITY_PRIZE_LOG)
	defer closeFun()

	query := col.Find(match)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishActivityPrizeLog
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)

}

// 获取活动期间获奖记录列表
func QueryActDrawRecordListByPid(req *brower_backstage.ListRequest) ([]*share_message.WishLog, int32) {

	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	match := bson.M{"DareId": req.GetId()}
	if req.GetBeginTimestamp() != 0 {
		match["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	} else {
		// 默认查询活动时间期间的
		cfg := GetWishActCfg(for_game.WISH_ACT_COUNT) // todo
		match["CreateTime"] = bson.M{"$gte": cfg.GetStartTime(), "$lte": cfg.GetEndTime()}
	}
	//id查询
	if req.GetKeyword() != "" {
		id := easygo.StringToInt64noErr(req.GetKeyword())
		match["WishBoxId"] = id
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG)
	defer closeFun()

	query := col.Find(match)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishLog
	var errc error

	errc = query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)

}

// 生成许愿池活动相关报表
func MakeWishActReports() {
	MakeWishActReport()
}

// 生成许愿池活动报表
func MakeWishActReport() {
	logs.Info("开始更新许愿池活动报表")
	endTime := easygo.GetToday0ClockTimestamp()
	startTime := endTime - 24*3600

	wishPoolReport := GetWishPoolReport(startTime)
	//次数
	countList := GroupWishActivityPrizeLog(1, startTime, endTime)
	//天数
	dayList := GroupWishActivityPrizeLog(2, startTime, endTime)
	report := &share_message.WishActivityReport{
		CreateTime:        easygo.NewInt64(startTime),
		InPoolCount:       easygo.NewInt64(wishPoolReport.GetInPoolCount()),
		InPoolPlayerCount: easygo.NewInt32(wishPoolReport.GetInPoolPlayerCount()),
		CounterData:       countList,
		DayCountData:      dayList,
	}
	UpdateWishActReport(report)
	logs.Info("完成更新许愿池活动报表")
	//写报表生成进度
	//MakeReportJob(for_game.TABLE_REPORT_WISH_ACTIVITY, endTime)
}

//初始化许愿池活动报表
func InitWishActReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_ACTIVITY)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//更新许愿池活动报表
func UpdateWishActReport(report *share_message.WishActivityReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_ACTIVITY)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report})
	easygo.PanicError(err)

}

//查询许愿池商品报表报表
func GetWishActReport(startTime int64) *share_message.WishActivityReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_ACTIVITY)
	defer closeFun()
	report := &share_message.WishActivityReport{}
	err := col.Find(bson.M{"_id": startTime}).One(&report)

	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return report
}

// 查询活动奖项日志记录统计
func GroupWishActivityPrizeLog(t int32, startTime, endTime int64) []*share_message.WishActivityUnit {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACTIVITY_PRIZE_LOG)
	defer closeFun()

	match := bson.M{"Type": t}
	match["CreateTime"] = bson.M{"$gte": startTime, "$lt": endTime}
	group := bson.M{
		"_id":   "$WishActPoolRuleId",
		"Value": bson.M{"$sum": 1},
		//"Key":   bson.M{"$last": "$ActType"},
	}

	project := bson.M{
		"Value":             1,
		"WishActPoolRuleId": "$_id",
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var list []*share_message.WishActivityUnit
	err := col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)

	return list

}

//添加白名单
func AddWishAllowList(playerBaseList []*share_message.PlayerBase, remark string) {
	//playerBaseList := QueryplayerlistByAccounts(accountsMap)
	playerIds := make([]int64, 0, 0) //有效的playerId

	for k, v := range playerBaseList {
		playerIds = append(playerIds, v.GetPlayerId())
		//playerBase的types = 1对应wishplayer的types= 0
		if v.GetTypes() == 1 {
			playerBaseList[k].Types = easygo.NewInt32(0)
		}

	}

	if len(playerIds) > 0 {
		//从WishPlayer中找到PlayerIds中的用户
		wishPlayerList := GetWishPlayerListByPlayers(playerIds)

		//已有的白名单列表
		wishWhiteList := GetWishWhiteListByPlayers(playerIds)
		var insertWishWhite []interface{}
		//var updateData []interface{}
		var insertData []interface{}
		for _, playerBase := range playerBaseList {
			playerId := playerBase.GetPlayerId()
			var isExistWishPlayer bool
			for _, wishPlayer := range wishPlayerList {
				if wishPlayer.GetPlayerId() == playerId {
					//updateData = append(updateData, bson.M{"PlayerId": playerId}, bson.M{"$set": bson.M{"Types": 2}})
					isExistWishPlayer = true

					//wishPlayer中存在记录的情况下，wishWhite可能没有对应记录
					var isExistWishWhite bool
					for _, wishWhite := range wishWhiteList {
						if wishWhite.GetPlayerId() == playerId {
							isExistWishWhite = true
						}
					}
					if !isExistWishWhite {
						wishWhite := &share_message.WishWhite{
							Id:       easygo.NewInt64(wishPlayer.GetId()),
							NickName: easygo.NewString(playerBase.GetNickName()),
							Account:  easygo.NewString(playerBase.GetAccount()),
							PlayerId: easygo.NewInt64(playerId),
							Note:     easygo.NewString(remark),
						}
						insertWishWhite = append(insertWishWhite, wishWhite)
					}
					break
				}
			}
			if !isExistWishPlayer {
				wishPlayerId := for_game.NextId(for_game.TABLE_WISH_PLAYER)
				uData := &share_message.WishPlayer{
					Id:         easygo.NewInt64(wishPlayerId),
					NickName:   easygo.NewString(playerBase.GetNickName()),
					HeadUrl:    easygo.NewString(playerBase.GetHeadIcon()),
					PlayerId:   easygo.NewInt64(playerId),
					Types:      easygo.NewInt32(playerBase.GetTypes()),
					CreateTime: easygo.NewInt64(time.Now().Unix()),
				}
				insertData = append(insertData, uData)
				var ifExist bool
				for _, wishWhite := range wishWhiteList {
					if wishWhite.GetPlayerId() == playerId {
						ifExist = true
					}
				}
				if !ifExist {
					wishWhite := &share_message.WishWhite{
						Id:       easygo.NewInt64(wishPlayerId),
						NickName: easygo.NewString(playerBase.GetNickName()),
						Account:  easygo.NewString(playerBase.GetAccount()),
						PlayerId: easygo.NewInt64(playerId),
						Note:     easygo.NewString(remark),
					}
					insertWishWhite = append(insertWishWhite, wishWhite)
				}
			}
		}
		//新增白名单
		if len(insertWishWhite) > 0 {
			logs.Info("insertWishWhite:", insertWishWhite)
			for_game.InsertAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WHITE, insertWishWhite...)
		}

		//新增wishPlayer
		if len(insertData) > 0 {
			logs.Info("insertData:", insertData)
			for_game.InsertAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER, insertData...)
		}

		/*if len(updateData) > 0 {
			for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER, updateData)
		}*/
	}
}

//远程调用RpcBackstageSetGuardian
func CallRpcBackstageSetGuardian(msg *h5_wish.BackstageSetGuardianReq) *base.Fail {
	serInfo := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)
	sid := serInfo.GetSid()
	response, err := SendMsgToServerNew(sid, "RpcBackstageSetGuardian", msg)
	if err != nil {
		logs.Error("CallRpcBackstageSetGuardian err:", err)
		return err
	}
	resp, ok := response.(*h5_wish.BackstageSetGuardianResp)
	if !ok {
		logs.Error("CallRpcBackstageSetGuardian err:not h5_wish.BackstageSetGuardianResp")
		return &base.Fail{
			Reason: easygo.NewString("操作失败"),
			Code:   easygo.NewString("调用失败"),
		}
	}
	if resp.GetResult() != 1 {
		logs.Error("CallRpcBackstageSetGuardian err: resp=", resp)
		return &base.Fail{
			Reason: easygo.NewString("操作失败"),
			Code:   easygo.NewString("调用失败"),
		}
	}
	return nil
}

//远程调用RpcSetWishAccount
func CallRpcSetWishAccount(msg *server_server.PlayerSI) *base.Fail {
	serInfo := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)
	sid := serInfo.GetSid()
	response, err := SendMsgToServerNew(sid, "RpcSetWishAccount", msg)
	if err == nil {
		resp, ok := response.(*h5_wish.BackstageSetGuardianResp)
		if ok && resp.GetResult() == 1 {
			return nil
		}
		return &base.Fail{
			Reason: easygo.NewString("修改用户数据失败"),
		}
	}
	return err
}

//根据时间统计进入许愿池的次数
func CountWishTime(startTime, endTime int64) int64 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_LOG)
	defer closeFun()
	inPoolPipeBson := []bson.M{
		{"$match": bson.M{"CreateTime": bson.M{"$gte": startTime, "$lt": endTime}}},
		{"$group": bson.M{"_id": nil, "WishTime": bson.M{"$sum": "$WishTime"}}},
		{"$project": bson.M{"WishTime": 1, "CreateTime": 1}},
	}
	var one *share_message.WishLogReport
	err := col.Pipe(inPoolPipeBson).One(&one)
	if err != nil {
		return 0
	}
	return one.GetWishTime()
}

//统计兑换的钻石总额
func CountExchangeDiamond(startTime, endTime int64) int64 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_DIAMOND_CHANGELOG)
	defer closeFun()
	inPoolPipeBson := []bson.M{
		{"$match": bson.M{"CreateTime": bson.M{"$gte": startTime, "$lt": endTime}}},
		{"$group": bson.M{"_id": nil, "ChangeDiamond": bson.M{"$sum": "$ChangeDiamond"}}},
		{"$project": bson.M{"ChangeDiamond": 1, "CreateTime": 1}},
	}
	var one *share_message.DiamondChangeLog
	err := col.Pipe(inPoolPipeBson).One(&one)
	if err != nil {
		return 0
	}
	return one.GetChangeDiamond()
}
