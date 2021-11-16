package backstage

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"strconv"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//====================================================================================================统一生成报表
//每小时要生成的报表任务
func MakeReports() {
	MakePlayerActiveReport()
	MakePlayerWeekActiveReport()
	MakePlayerMonthActiveReport()
	// MakeChannelReport()
	// logs.Info("更新渠道报表")
	MakeSquareReport()
	MakeAdvReport()
	MakeNearbyAdvReport()
	for_game.SaveRedisButtonClickReportToMongo() //保存redis中的按钮点击行为报表数据
}

//每天要生成的报表任务
func MakeReportsForDay() {
	MakePageRegLogReport()          //注册登录页面埋点报表
	MakeCoinProductReport()         //虚拟商城报表
	for_game.UpdateInOutCashCount() //更新出入款汇总报表出入款人数

	MakeBetSlipReport() //注单统计报表
}

//====================================================================================================报表任务进度处理
//写报表任务进度
func MakeReportJob(report string, time int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORTJOB)
	defer closeFun()

	v := &share_message.ReportJob{
		Report: easygo.NewString(report),
		Time:   easygo.NewInt64(time),
	}
	_, err := col.Upsert(bson.M{"_id": report}, bson.M{"$set": v})
	easygo.PanicError(err)
}

//读报表任务进度
func GetReportJob(report string) *share_message.ReportJob {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORTJOB)
	defer closeFun()

	job := &share_message.ReportJob{}
	err := col.Find(bson.M{"_id": report}).One(job)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return job
}

//===============================================================================================================查询用户留存报表
//生成用户留存报表
func MakePlayerKeepReport(startTime, endTime TIME_64) {
	if endTime == 0 {
		endTime = time.Now().Unix()
	}

	if startTime == 0 {
		InitPlayerKeepReport() //初始化用户留存报表
	} else {
		DeletePlayerKeepReport(startTime, endTime)
	}

	logList := GetOnlineTimeLogList(startTime, endTime)

	var ids []int64
	for _, p := range logList {
		if for_game.IsContains(p.GetPlayerId(), ids) == -1 {
			ids = append(ids, p.GetPlayerId())
		}
	}

	mapPlayer := make(map[PLAYER_ID]*share_message.PlayerBase)
	players := QueryplayerlistByIds(ids)
	for _, player := range players {
		mapPlayer[player.GetPlayerId()] = player
	}

	for _, item := range logList {
		report := for_game.QueryPlayerKeepReport(item.GetCreateTime())
		if report == nil {
			report = &share_message.PlayerKeepReport{
				CreateTime: easygo.NewInt64(item.GetCreateTime()),
			}
		}

		createTime := mapPlayer[item.GetPlayerId()].GetCreateTime() / 1000
		createtime0timestamp := easygo.Get0ClockTimestamp(createTime) //注册日0点
		//更新今日注册用户数
		if createtime0timestamp == item.GetCreateTime() {
			report.TodayRegister = easygo.NewInt32(report.GetTodayRegister() + 1)
		} else {
			days := easygo.GetDifferenceDay(createTime, item.GetCreateTime())
			reportold := for_game.QueryPlayerKeepReport(createtime0timestamp) //查询注册日的报表
			if reportold != nil {
				report = reportold
			}
			//更新留存人数
			switch days {
			case 1:
				report.NextKeep = easygo.NewInt32(report.GetNextKeep() + 1)
			case 2:
				report.ThreeKeep = easygo.NewInt32(report.GetThreeKeep() + 1)
			case 3:
				report.FourKeep = easygo.NewInt32(report.GetFourKeep() + 1)
			case 4:
				report.FiveKeep = easygo.NewInt32(report.GetFiveKeep() + 1)
			case 5:
				report.SixKeep = easygo.NewInt32(report.GetSixKeep() + 1)
			case 6:
				report.SevenKeep = easygo.NewInt32(report.GetSevenKeep() + 1)
			case 7:
				report.EightKeep = easygo.NewInt32(report.GetEightKeep() + 1)
			case 8:
				report.NineKeep = easygo.NewInt32(report.GetNineKeep() + 1)
			case 9:
				report.TenKeep = easygo.NewInt32(report.GetTenKeep() + 1)
			case 10:
				report.ElevenKeep = easygo.NewInt32(report.GetElevenKeep() + 1)
			case 11:
				report.TwelveKeep = easygo.NewInt32(report.GetTwelveKeep() + 1)
			case 12:
				report.ThirteenKeep = easygo.NewInt32(report.GetTwelveKeep() + 1)
			case 13:
				report.FourteenKeep = easygo.NewInt32(report.GetFourteenKeep() + 1)
			case 14:
				report.FifteenKeep = easygo.NewInt32(report.GetFifteenKeep() + 1)
			case 15:
				report.SixteenKeep = easygo.NewInt32(report.GetSixteenKeep() + 1)
			case 16:
				report.SeventeenKeep = easygo.NewInt32(report.GetSeventeenKeep() + 1)
			case 17:
				report.EighteenKeep = easygo.NewInt32(report.GetEighteenKeep() + 1)
			case 18:
				report.NineteenKeep = easygo.NewInt32(report.GetNineteenKeep() + 1)
			case 19:
				report.TwentyKeep = easygo.NewInt32(report.GetTwentyKeep() + 1)
			case 20:
				report.TwentyOneKeep = easygo.NewInt32(report.GetTwentyOneKeep() + 1)
			case 21:
				report.TwentyTwoKeep = easygo.NewInt32(report.GetTwentyTwoKeep() + 1)
			case 22:
				report.TwentyThreeKeep = easygo.NewInt32(report.GetTwentyThreeKeep() + 1)
			case 23:
				report.TwentyFourKeep = easygo.NewInt32(report.GetTwentyFourKeep() + 1)
			case 24:
				report.TwentyFiveKeep = easygo.NewInt32(report.GetTwentyFiveKeep() + 1)
			case 25:
				report.TwentySixKeep = easygo.NewInt32(report.GetTwentySixKeep() + 1)
			case 26:
				report.TwentySevenKeep = easygo.NewInt32(report.GetTwentySevenKeep() + 1)
			case 27:
				report.TwentyEightKeep = easygo.NewInt32(report.GetTwentyEightKeep() + 1)
			case 28:
				report.TwentyNineKeep = easygo.NewInt32(report.GetTwentyNineKeep() + 1)
			case 29:
				report.Thirtykeep = easygo.NewInt32(report.GetThirtykeep())
			}
		}

		UpdatePlayerKeepReport(report)
	}
}

//初始化用户留存报表
func InitPlayerKeepReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYERKEEPREPORT)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//删除指定时间范围用户留存报表数据
func DeletePlayerKeepReport(startTime, endTime TIME_64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYERKEEPREPORT)
	defer closeFun()
	queryBson := bson.M{"_id": bson.M{"$gte": startTime, "$lte": endTime}}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//查询指定时间留存报表
func GetPlayerKeepReportBytime(querytime int64) *share_message.PlayerKeepReport {
	querytime = easygo.Get0ClockTimestamp(querytime)
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYERKEEPREPORT)
	defer closeFun()
	var obj *share_message.PlayerKeepReport
	err := col.Find(bson.M{"_id": querytime}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//更新用户留存报表
func UpdatePlayerKeepReport(report *share_message.PlayerKeepReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYERKEEPREPORT)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report})
	easygo.PanicError(err)
}

//查询用户留存报表
func GetPlayerKeepReport(reqMsg *brower_backstage.ListRequest, user *share_message.Manager) ([]*share_message.PlayerKeepReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYERKEEPREPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		sort = reqMsg.GetKeyword()
	}

	var list []*share_message.PlayerKeepReport
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//===============================================================================================================用户活跃报表
//生成用户日活跃报表
func MakePlayerActiveReport() {
	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_PLAYER_ACTIVE_REPORT)
	if job != nil {
		if job.GetTime() < easygo.GetToday0ClockTimestamp() {
			startTime = job.GetTime()
		} else {
			startTime = easygo.GetToday0ClockTimestamp()
			DeleteTodayPlayerActiveReport(startTime)
		}
	} else {
		InitPlayerActiveReport() //初始化用户日留存报表
	}

	logList := GetOnlineTimeLogList(startTime, endTime)
	report := &share_message.PlayerActiveReport{}
	createTime := int64(0)

	for _, item := range logList {
		if createTime == 0 {
			report = &share_message.PlayerActiveReport{
				CreateTime:    easygo.NewInt64(easygo.Get0ClockTimestamp(item.GetCreateTime())),
				LoginCount:    easygo.NewInt64(1),
				RegisterCount: easygo.NewInt64(0),
				SumOnlineTime: item.OnlineTime,
				AveOnlineTime: item.OnlineTime,
			}
			createTime = item.GetCreateTime()
		} else {
			if item.GetCreateTime() > createTime {
				UpdatePlayerActiveReport(report)
				report = &share_message.PlayerActiveReport{
					CreateTime:    easygo.NewInt64(easygo.Get0ClockTimestamp(item.GetCreateTime())),
					LoginCount:    easygo.NewInt64(1),
					RegisterCount: easygo.NewInt64(0),
					SumOnlineTime: item.OnlineTime,
					AveOnlineTime: item.OnlineTime,
				}
				createTime = item.GetCreateTime()

			}
			report.LoginCount = easygo.NewInt64(report.GetLoginCount() + 1)
			report.SumOnlineTime = easygo.NewInt64(report.GetSumOnlineTime() + item.GetOnlineTime())
			report.AveOnlineTime = easygo.NewInt64(report.GetSumOnlineTime() / report.GetLoginCount())
		}
		player := for_game.GetRedisPlayerBase(item.GetPlayerId())
		if easygo.GetDifferenceDay(player.GetCreateTime(), item.GetCreateTime()) == 0 {
			if len(player.GetLabelList()) > 0 {
				report.RegisterCount = easygo.NewInt64(report.GetRegisterCount() + 1)
			}
		}
	}

	if report.CreateTime != nil {
		UpdatePlayerActiveReport(report)
	}
	//写报表生成进度
	MakeReportJob(for_game.TABLE_PLAYER_ACTIVE_REPORT, endTime)
}

//初始化用户活跃报表
func InitPlayerActiveReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_ACTIVE_REPORT)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//删除当天数据
func DeleteTodayPlayerActiveReport(todaytime int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_ACTIVE_REPORT)
	defer closeFun()

	err := col.RemoveId(todaytime) //数据库删除
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
}

//更新用户活跃报表
func UpdatePlayerActiveReport(report *share_message.PlayerActiveReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_ACTIVE_REPORT)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report})
	easygo.PanicError(err)
	logs.Info("更新日活跃报表")
}

//查询用户活跃报表
func GetPlayerActiveReportList(reqMsg *brower_backstage.ListRequest) ([]*share_message.PlayerActiveReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_ACTIVE_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		sort = reqMsg.GetKeyword()
	}

	var list []*share_message.PlayerActiveReport
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//===============================================================================================================周活跃报表
//生成周活跃报表
func MakePlayerWeekActiveReport() {
	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_PLAYER_WEEK_ACTIVE_REPORT)
	if job != nil {
		if job.GetTime() < easygo.GetWeek0ClockTimestamp() {
			startTime = job.GetTime()
		} else {
			startTime = easygo.GetWeek0ClockTimestamp()
			DeleteWeekPlayerActiveReport(startTime)
		}
	} else {
		InitPlayerWeekActiveReport() //初始化周活跃报表
	}

	logList := GetOnlineTimeLogList(startTime, endTime)
	report := &share_message.PlayerActiveReport{}
	createTime := int64(0)

	for _, item := range logList {
		id := easygo.GetWeek0ClockOfTimestamp(item.GetCreateTime())
		if createTime == 0 {
			report = &share_message.PlayerActiveReport{
				CreateTime:    easygo.NewInt64(id),
				RegisterCount: easygo.NewInt64(0),
				LoginCount:    easygo.NewInt64(1),
				SumOnlineTime: item.OnlineTime,
				AveOnlineTime: item.OnlineTime,
			}
			createTime = id
		} else {
			if id > createTime {
				UpdatePlayerWeekActiveReport(report)
				report = &share_message.PlayerActiveReport{
					CreateTime:    easygo.NewInt64(id),
					RegisterCount: easygo.NewInt64(0),
					LoginCount:    easygo.NewInt64(1),
					SumOnlineTime: item.OnlineTime,
					AveOnlineTime: item.OnlineTime,
				}
				createTime = id
			}
			report.LoginCount = easygo.NewInt64(report.GetLoginCount() + 1)
			report.SumOnlineTime = easygo.NewInt64(report.GetSumOnlineTime() + item.GetOnlineTime())
			report.AveOnlineTime = easygo.NewInt64(report.GetSumOnlineTime() / report.GetLoginCount())
		}
		player := for_game.GetRedisPlayerBase(item.GetPlayerId())
		if player != nil {
			if easygo.GetDifferenceDay(player.GetCreateTime(), item.GetCreateTime()) == 0 {
				report.RegisterCount = easygo.NewInt64(report.GetRegisterCount() + 1)
			}
		} else {
			logs.Error("不存在的玩家id", item.GetPlayerId())
		}

	}
	if report.CreateTime != nil {
		UpdatePlayerWeekActiveReport(report)
	}

	//写报表生成进度
	MakeReportJob(for_game.TABLE_PLAYER_WEEK_ACTIVE_REPORT, endTime)
}

//初始化用户活跃报表
func InitPlayerWeekActiveReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WEEK_ACTIVE_REPORT)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//删除当周数据
func DeleteWeekPlayerActiveReport(todaytime int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_MONTH_ACTIVE_REPORT)
	defer closeFun()

	err := col.RemoveId(todaytime) //数据库删除
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
}

//更新周活跃报表
func UpdatePlayerWeekActiveReport(report *share_message.PlayerActiveReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WEEK_ACTIVE_REPORT)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report})
	easygo.PanicError(err)
	logs.Info("更新周活跃报表")
}

//查询用户周活跃报表
func GetPlayerWeekActiveReportList(reqMsg *brower_backstage.ListRequest) ([]*share_message.PlayerActiveReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WEEK_ACTIVE_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		sort = reqMsg.GetKeyword()
	}

	var list []*share_message.PlayerActiveReport
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//===============================================================================================================周活跃报表
//生成月活跃报表
func MakePlayerMonthActiveReport() {
	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_PLAYER_MONTH_ACTIVE_REPORT)
	if job != nil {
		if job.GetTime() < easygo.GetMonth0ClockTimestamp() {
			startTime = job.GetTime()
		} else {
			startTime = easygo.GetMonth0ClockTimestamp()
			DeleteMonthPlayerActiveReport(startTime)
		}
	} else {
		InitPlayerMonthActiveReport() //初始化周活跃报表
	}

	logList := GetOnlineTimeLogList(startTime, endTime)
	report := &share_message.PlayerActiveReport{}
	createTime := int64(0)

	for _, item := range logList {
		id := easygo.GetMonth0ClockOfTimestamp(item.GetCreateTime())
		if createTime == 0 {
			report = &share_message.PlayerActiveReport{
				CreateTime:    easygo.NewInt64(id),
				RegisterCount: easygo.NewInt64(0),
				LoginCount:    easygo.NewInt64(1),
				SumOnlineTime: item.OnlineTime,
				AveOnlineTime: item.OnlineTime,
			}
			createTime = id
		} else {
			if id > createTime {
				UpdatePlayerMonthActiveReport(report)
				report = &share_message.PlayerActiveReport{
					CreateTime:    easygo.NewInt64(id),
					RegisterCount: easygo.NewInt64(0),
					LoginCount:    easygo.NewInt64(1),
					SumOnlineTime: item.OnlineTime,
					AveOnlineTime: item.OnlineTime,
				}
				createTime = id
			}
			report.LoginCount = easygo.NewInt64(report.GetLoginCount() + 1)
			report.SumOnlineTime = easygo.NewInt64(report.GetSumOnlineTime() + item.GetOnlineTime())
			report.AveOnlineTime = easygo.NewInt64(report.GetSumOnlineTime() / report.GetLoginCount())
		}
		player := for_game.GetRedisPlayerBase(item.GetPlayerId())
		if player != nil && easygo.GetDifferenceDay(player.GetCreateTime(), item.GetCreateTime()) == 0 {
			report.RegisterCount = easygo.NewInt64(report.GetRegisterCount() + 1)
		}
	}
	if report.CreateTime != nil {
		UpdatePlayerMonthActiveReport(report)
	}

	//写报表生成进度
	MakeReportJob(for_game.TABLE_PLAYER_MONTH_ACTIVE_REPORT, endTime)
}

//初始化用户月活跃报表
func InitPlayerMonthActiveReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_MONTH_ACTIVE_REPORT)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//删除当月数据
func DeleteMonthPlayerActiveReport(todaytime int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_MONTH_ACTIVE_REPORT)
	defer closeFun()

	err := col.RemoveId(todaytime) //数据库删除
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
}

//更新用户月活跃报表
func UpdatePlayerMonthActiveReport(report *share_message.PlayerActiveReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_MONTH_ACTIVE_REPORT)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report})
	easygo.PanicError(err)
	logs.Info("更新月活跃报表")
}

//查询用户月活跃报表
func GetPlayerMonthActiveReportList(reqMsg *brower_backstage.ListRequest) ([]*share_message.PlayerActiveReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_MONTH_ACTIVE_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		sort = reqMsg.GetKeyword()
	}

	var list []*share_message.PlayerActiveReport
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//===============================================================================================================查询用户行为报表
func GetPlayerBehaviorReport(reqMsg *brower_backstage.ListRequest) ([]*share_message.PlayerBehaviorReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BEHAVIOR_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		sort = reqMsg.GetKeyword()
	}

	var list []*share_message.PlayerBehaviorReport
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//补充用户行为报表附近的人打招呼数
func MakePlayerBehaviorReportNeerbyData() {
	startTime := easygo.GetYesterday0ClockTimestamp() * 1000
	endTime := easygo.GetYesterday24ClockTimestamp() * 1000

	m := []bson.M{
		{"$match": bson.M{"IsFirst": true, "CreateTime": bson.M{"$gte": startTime, "$lt": endTime}}},
		{"$group": bson.M{"_id": "$SendPlayerId"}},
	}

	ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_NEARBY_MESSAGE_NEW_LOG, m, 0, 0)
	start := len(ls)

	m = []bson.M{
		{"$match": bson.M{"IsFirst": bson.M{"$ne": true}, "CreateTime": bson.M{"$gte": startTime, "$lt": endTime}}},
		{"$group": bson.M{"_id": "$SendPlayerId"}},
	}

	ls = for_game.FindPipeAll(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_NEARBY_MESSAGE_NEW_LOG, m, 0, 0)
	reply := len(ls)

	reportId := startTime / 1000
	if start > 0 {
		for_game.SetRedisPlayerBehaviorReportFildVal(reportId, int64(start), "Start")

	}
	if reply > 0 {
		for_game.SetRedisPlayerBehaviorReportFildVal(reportId, int64(reply), "Reply")
	}
}

//===============================================================================================================查询出入款汇总报表
func GetInOutCashSumReport(reqMsg *brower_backstage.ListRequest) ([]*share_message.InOutCashSumReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_INOUTCASHSUM_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		sort = reqMsg.GetKeyword()
	}

	var list []*share_message.InOutCashSumReport
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//==============================================================================================================埋点注册登录报表
// 查询埋点注册登录报表
func GetRegisterLoginReport(reqMsg *brower_backstage.ListRequest, user *share_message.Manager) ([]*share_message.RegisterLoginReport, int) {
	// col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, "bak_report_register_login")
	// if user.GetRole() == 0 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_LOGIN_REGISTER_REPORT)
	// }
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"
	if reqMsg.Sort != nil && reqMsg.GetSort() != "" {
		sort = reqMsg.GetSort()
	}

	var list []*share_message.RegisterLoginReport
	var errc error
	if pageSize > 0 {
		errc = query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	} else {
		errc = query.Sort(sort).All(&list)
	}

	easygo.PanicError(errc)

	return list, count
}

//查询指定时间埋点注册登录报表
func GetRegisterLoginReportByTime(querytime int64) *share_message.RegisterLoginReport {
	querytime = easygo.Get0ClockTimestamp(querytime)
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_LOGIN_REGISTER_REPORT)
	defer closeFun()
	var obj *share_message.RegisterLoginReport
	err := col.Find(bson.M{"_id": querytime}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//=========================================================================================================运营渠道数据汇总报表
//查询运营渠道数据汇总报表
func GetOperationChannelReportList(reqMsg *brower_backstage.ListRequest, user *share_message.Manager) ([]*share_message.OperationChannelReport, int) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_OPERATION_CHANNEL_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["ChannelName"] = reqMsg.GetKeyword()
		case 2:
			types := 0
			switch reqMsg.GetKeyword() {
			case "cpa":
				types = 1
			case "cps":
				types = 2
			case "cpc":
				types = 3
			case "cpd":
				types = 4
			}
			queryBson["Cooperation"] = types
		}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.OperationChannelReport
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//时间查询运营渠道数据汇总报表
func GetOperationChannelReportListByTime(startTime int64, endTime int64) []*share_message.OperationChannelReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_OPERATION_CHANNEL_REPORT)
	defer closeFun()

	queryBson := bson.M{}
	queryBson["CreateTime"] = bson.M{"$gte": startTime, "$lte": endTime}

	query := col.Find(queryBson)
	sort := "-CreateTime"

	var list []*share_message.OperationChannelReport
	errc := query.Sort(sort).All(&list)
	easygo.PanicError(errc)

	return list
}

//=========================================================================================================渠道报表
//生成渠道报表
/*
func MakeChannelReport() {
	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_CHANNEL_REPORT)
	if job != nil {
		startTime = easygo.GetToday0ClockTimestamp()
	} else {
		InitChannelReport()
	}

	logList := GetOperationChannelReportListByTime(startTime, endTime)
	report := &share_message.ChannelReport{}

	for _, item := range logList {
		report = GetChannelReportByNo(item.GetCreateTime(), item.GetChannelNo())
		if report == nil {
			report = &share_message.ChannelReport{
				Id:          easygo.NewInt64(for_game.NextId(for_game.TABLE_CHANNEL_REPORT)),
				ChannelNo:   easygo.NewString(item.GetChannelNo()),   //渠道编号
				CreateTime:  easygo.NewInt64(item.GetCreateTime()),   //报表时间
				ChannelName: easygo.NewString(item.GetChannelName()), //渠道名称
				Cooperation: easygo.NewInt32(item.GetCooperation()),  //渠道类型： -1 暂无，0 全部，1 cpa,2 cps,3 cpc,4 cpd
				ChannelCost: easygo.NewInt64(0),                      //渠道成本
				ActCost:     easygo.NewInt64(0),                      //激活成本
				LoginCost:   easygo.NewInt64(0),                      //登录成本
				RegCost:     easygo.NewInt64(0),                      //注册成本
				ROI:         easygo.NewFloat64(0),                    //ROI
				RegRoiRate:  easygo.NewFloat64(0),                    //注册转化率
				KeepRate:    easygo.NewFloat64(0),                    //留存率
			}
		} else {
			report.ChannelCost = easygo.NewInt64(0)  //渠道成本
			report.ActCost = easygo.NewInt64(0)      //激活成本
			report.LoginCost = easygo.NewInt64(0)    //登录成本
			report.RegCost = easygo.NewInt64(0)      //注册成本
			report.ROI = easygo.NewFloat64(0)        //ROI
			report.RegRoiRate = easygo.NewFloat64(0) //注册转化率
			report.KeepRate = easygo.NewFloat64(0)   //留存率
		}

		channel := for_game.QueryOperationByNo(item.GetChannelNo()) //查询渠道信息
		switch item.GetCooperation() {
		case 1: //cpa的渠道成本=单价*注册人数
			report.ChannelCost = easygo.NewInt64(channel.GetPrice() * item.GetRegCount()) //渠道成本
		case 2: //cps的渠道成本=成交金额*佣金比例
			report.ChannelCost = easygo.NewInt64(channel.GetRate() / 1000 * item.GetShopDealSumAmount()) //渠道成本
		case 3: //cpc的渠道成本=单价*UV数
			report.ChannelCost = easygo.NewInt64(channel.GetPrice() * item.GetUvCount()) //渠道成本
		case 4: //cpd的渠道成本=单价*下载数
			report.ChannelCost = easygo.NewInt64(channel.GetPrice() * item.GetDownLoadCount()) //渠道成本
		}

		cost := report.GetChannelCost()
		if cost > 0 {
			if item.GetActDevCount() != 0 {
				report.ActCost = easygo.NewInt64(cost / item.GetActDevCount())
			}
			if item.GetLoginCount() != 0 {
				report.LoginCost = easygo.NewInt64(cost / item.GetLoginCount())
			}
			if item.GetRegCount() != 0 {
				report.RegCost = easygo.NewInt64(cost / item.GetRegCount())
			}
			report.ROI = easygo.NewFloat64(item.GetShopDealSumAmount() / cost)
			if item.GetUvCount() != 0 {
				report.RegRoiRate = easygo.NewFloat64(item.GetRegCount() / item.GetUvCount() * 100)
			}

			yesTime := item.GetCreateTime() - 86400
			yesreport := for_game.GetOperationChannelReport(item.GetChannelNo(), yesTime)
			if yesreport.GetRegCount() > 0 {
				report.KeepRate = easygo.NewFloat64(item.GetLoginCount() / yesreport.GetRegCount() * 100)
			}
		}

		UpdateChannelReport(report)
	}

	//写报表生成进度
	MakeReportJob(for_game.TABLE_CHANNEL_REPORT, endTime)
}*/

//初始化渠道报表
func InitChannelReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CHANNEL_REPORT)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//更新渠道报表
func UpdateChannelReport(report *share_message.ChannelReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CHANNEL_REPORT)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetId()}, bson.M{"$set": report})
	easygo.PanicError(err)
}

//查询时间渠道号渠道报表
func GetChannelReportByNo(querytime int64, channelNo string) *share_message.ChannelReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CHANNEL_REPORT)
	defer closeFun()
	var obj *share_message.ChannelReport
	err := col.Find(bson.M{"CreateTime": querytime, "ChannelNo": channelNo}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//查询渠道报表
func GetChannelReportList(reqMsg *brower_backstage.ListRequest) ([]*share_message.ChannelReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CHANNEL_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetListType() {
		case 1:
			queryBson["ChannelName"] = reqMsg.GetKeyword()
		case 2:
			types := 0
			switch reqMsg.GetKeyword() {
			case "cpa":
				types = 1
			case "cps":
				types = 2
			case "cpc":
				types = 3
			case "cpd":
				types = 4
			}
			queryBson["Cooperation"] = types
		}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-CreateTime"
	var list []*share_message.ChannelReport
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//=======================================================================================================文章报表
//查询文章报表
func GetArticleReportList(reqMsg *brower_backstage.ListRequest) ([]*share_message.ArticleReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ARTICLE_REPORT)
	defer closeFun()
	queryBson := bson.M{}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		queryBson["Article.ArticleType"] = reqMsg.GetListType()
	}
	if reqMsg.DownType != nil && reqMsg.GetDownType() != 0 {
		queryBson["Article.IsMain"] = reqMsg.GetDownType()
	}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Article.Title"] = reqMsg.GetKeyword()
		case 2:
			i, _ := strconv.ParseInt(reqMsg.GetKeyword(), 10, 64)
			queryBson["Article._id"] = i
		}
	}

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	//联表查询
	m := []bson.M{
		{"$lookup": bson.M{
			"from":         for_game.TABLE_ARTICLE,
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "Article",
		}},
		{"$match": queryBson},
		{"$project": bson.M{"_id": 1, "CreateTime": 1, "UpdateTime": 1, "PushPlayer": 1, "Clicks": 1, "Jumps": 1, "Title": "$Article.Title", "Types": "$Article.ArticleType", "IsMain": "$Article.IsMain"}},
		{"$unwind": "$Title"},
		{"$unwind": "$Types"},
		{"$unwind": "$IsMain"},
		{"$sort": bson.M{"CreateTime": -1}},
		// {"$skip": curPage * pageSize},
		// {"$limit": pageSize},
	}
	var list []*share_message.ArticleReport
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	count := len(list)
	m = append(m, bson.M{"$skip": curPage * pageSize}, bson.M{"$limit": pageSize})
	querypage := col.Pipe(m)

	errc := querypage.All(&list)
	easygo.PanicError(errc)
	return list, count
}

//=======================================================================================================推送通知报表
//查询推送通知报表
func GetNoticeReportList(reqMsg *brower_backstage.ListRequest) ([]*share_message.ArticleReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_NOTICE_REPORT)
	defer closeFun()
	queryBson := bson.M{}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Article.Title"] = reqMsg.GetKeyword()
		case 2:
			i, _ := strconv.ParseInt(reqMsg.GetKeyword(), 10, 64)
			queryBson["Article._id"] = i
		}
	}

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	//联表查询
	m := []bson.M{
		{"$lookup": bson.M{
			"from":         for_game.TABLE_FEATURES_APPPUSHMSG,
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "Article",
		}},
		{"$match": queryBson},
		{"$project": bson.M{"_id": 1, "CreateTime": 1, "PushPlayer": 1, "Clicks": 1, "Title": "$Article.Title"}},
		{"$unwind": "$Title"},
		{"$sort": bson.M{"CreateTime": -1}},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	query := col.Pipe(m)

	var list []*share_message.ArticleReport
	errc := query.All(&list)
	easygo.PanicError(errc)
	count, _ := col.Count()

	return list, count
}

//=======================================================================================================社交广场报表
//生成社交广场报表
func MakeSquareReport() {
	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_SQUARE_REPORT)
	if job != nil {
		startTime = job.GetTime()
	} else {
		InitSquareReport()
	}

	var report *share_message.SquareReport
	dList := QueryDynamicByTime(startTime, endTime)
	if len(dList) == 0 {
		return
	}

	for _, item := range dList {
		report = GetSquareReportByTime(item.GetSendTime())
		if report == nil {
			report = &share_message.SquareReport{
				CreateTime: easygo.NewInt64(easygo.Get0ClockTimestamp(item.GetCreateTime())),
			}
		}

		pmg := for_game.GetRedisPlayerBase(item.GetPlayerId())
		if pmg != nil {
			switch pmg.GetTypes() {
			case for_game.PLAYER_MARKET, for_game.PLAYER_SHOP, for_game.PLAYER_MANAGE, for_game.PLAYER_OFFICIAL:
				report.OperatSubCount = easygo.NewInt64(report.GetOperatSubCount() + 1)
			default:
				report.PlayerSubCount = easygo.NewInt64(report.GetPlayerSubCount() + 1)
			}
		}

		//删除的动态
		if item.GetStatue() > 0 {
			//动态状态 1后台删除，2前端删除
			switch item.GetStatue() {
			case 1:
				report.BackstageDelCount = easygo.NewInt64(report.GetBackstageDelCount() + 1)
			case 2:
				report.PlayerDelCount = easygo.NewInt64(report.GetPlayerDelCount() + 1)
			}
		}

		report.BsZanCount = easygo.NewInt64(report.GetBsZanCount() + int64(item.GetTrueZan()))

		UpdateSquareReport(report)
	}

	clist := QueryDynamicCommentByTime(startTime, endTime)
	if len(clist) > 0 {
		for _, c := range clist {
			report = GetSquareReportByTime(c.GetCreateTime())
			if report == nil {
				report = &share_message.SquareReport{
					CreateTime: easygo.NewInt64(easygo.Get0ClockTimestamp(c.GetCreateTime())),
				}
			}
			pmg := for_game.GetRedisPlayerBase(c.GetPlayerId())
			if pmg != nil {
				switch pmg.GetTypes() {
				case for_game.PLAYER_NORMAL:
					report.PlayerComm = easygo.NewInt64(report.GetPlayerComm() + 1)
				case for_game.PLAYER_MARKET, for_game.PLAYER_SHOP, for_game.PLAYER_MANAGE, for_game.PLAYER_OFFICIAL:
					report.OperatComm = easygo.NewInt64(report.GetOperatComm() + 1)
				}
				UpdateSquareReport(report)
			}
		}
	}

	zlist := QueryZanDataByTime(startTime, endTime)
	if len(zlist) > 0 {
		for _, z := range zlist {
			report = GetSquareReportByTime(z.GetCreateTime())
			if report == nil {
				report = &share_message.SquareReport{
					CreateTime: easygo.NewInt64(easygo.Get0ClockTimestamp(z.GetCreateTime())),
				}
			}
			report.ZanCount = easygo.NewInt64(report.GetZanCount() + 1)
			UpdateSquareReport(report)
		}
	}

	//写报表生成进度
	MakeReportJob(for_game.TABLE_SQUARE_REPORT, endTime)
}

//初始化社交广场报表
func InitSquareReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_REPORT)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//更新社交广场报表
func UpdateSquareReport(report *share_message.SquareReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_REPORT)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report})
	easygo.PanicError(err)
	logs.Info("更新社交广场报表")
}

//查询指定时间的社交广场报表
func GetSquareReportByTime(querytime TIME_64) *share_message.SquareReport {
	todaytime := easygo.Get0ClockTimestamp(querytime)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_REPORT)
	defer closeFun()

	var obj *share_message.SquareReport
	err := col.Find(bson.M{"_id": todaytime}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//查询社交广场报表
func GetSquareReportList(reqMsg *brower_backstage.ListRequest) ([]*share_message.SquareReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		startTime := reqMsg.GetBeginTimestamp()
		endTime := reqMsg.GetEndTimestamp()
		if reqMsg.GetBeginTimestamp() > 1000000000000 {
			startTime /= 1000
			endTime /= 1000
		}

		queryBson["_id"] = bson.M{"$gte": startTime, "$lte": endTime}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		sort = reqMsg.GetKeyword()
	}

	var list []*share_message.SquareReport
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//=======================================================================================================活动报表
//查询活动报表
func GetActivityReportList(reqMsg *brower_backstage.ListRequest) ([]*share_message.ActivityReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ACTIVITY_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"

	var list []*share_message.ActivityReport
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//======================================================================================================广告报表
//生成广告报表
func MakeAdvReport() {
	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_ADV_REPORT)
	if job != nil {
		startTime = job.GetTime()
	} else {
		InitAdvReport()
	}

	var report *share_message.AdvReport
	dList := for_game.QueryAdvLogByTime(startTime, endTime)
	if len(dList) == 0 {
		return
	}

	for _, item := range dList {
		id := easygo.AnytoA(easygo.Get0ClockTimestamp(item.GetOpTime())) + easygo.AnytoA(item.GetAdvId())
		report = GetAdvReportById(id)
		if report == nil {
			report = &share_message.AdvReport{
				Id:         easygo.NewString(id),
				AdvId:      easygo.NewInt64(item.GetAdvId()),
				CreateTime: easygo.NewInt64(easygo.Get0ClockTimestamp(item.GetOpTime())),
			}
		}

		switch item.GetOpType() {
		case for_game.ADV_LOG_OP_TYPE_1:
			report.PvCount = easygo.NewInt64(report.GetPvCount() + 1)
		case for_game.ADV_LOG_OP_TYPE_2:
			report.UvCount = easygo.NewInt64(report.GetUvCount() + 1)
		case for_game.ADV_LOG_OP_TYPE_3:
			report.Clicks = easygo.NewInt64(report.GetClicks() + 1)
		case for_game.ADV_LOG_OP_TYPE_4:
			report.ClickPlayers = easygo.NewInt64(report.GetClickPlayers() + 1)
		}

		UpdateAdvReport(report)
	}

	//写报表生成进度
	MakeReportJob(for_game.TABLE_ADV_REPORT, endTime)
}

//初始化广告报表
func InitAdvReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ADV_REPORT)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//更新广告报表
func UpdateAdvReport(report *share_message.AdvReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ADV_REPORT)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetId()}, bson.M{"$set": report})
	easygo.PanicError(err)
	// logs.Info("更新广告报表")
}

//查询广告报表
func GetAdvReportById(id string) *share_message.AdvReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ADV_REPORT)
	defer closeFun()

	var obj *share_message.AdvReport
	err := col.Find(bson.M{"_id": id}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//查询广告报表
func GetAdvReportList(reqMsg *brower_backstage.ListRequest) ([]*share_message.AdvReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ADV_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["AdvId"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		queryBson["Adv.Location"] = reqMsg.GetListType()
	}

	//联表查询
	m := []bson.M{
		{"$lookup": bson.M{
			"from":         for_game.TABLE_ADV_DATA,
			"localField":   "AdvId",
			"foreignField": "_id",
			"as":           "Adv",
		}},
		{"$match": queryBson},
		{"$unwind": "$Adv"},
		{"$sort": bson.M{"CreateTime": -1}},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	query := col.Pipe(m)
	count, err := col.Count()
	easygo.PanicError(err)

	var list []*share_message.AdvReport
	errc := query.All(&list)
	easygo.PanicError(errc)

	return list, count
}

//======================================================================================================附近的人引导项报表
//生成附近的人引导项报表
func MakeNearbyAdvReport() {
	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_NEARBY_ADV_REPORT)
	if job != nil {
		startTime = job.GetTime()
	} else {
		InitNearbyAdvReport()
	}

	var report *share_message.AdvReport
	dList := for_game.QueryNearbyAdvLogByTime(startTime, endTime)
	if len(dList) == 0 {
		return
	}

	for _, item := range dList {
		id := easygo.AnytoA(easygo.Get0ClockTimestamp(item.GetOpTime())) + easygo.AnytoA(item.GetAdvId())
		report = GetNearbyAdvReportById(id)
		if report == nil {
			report = &share_message.AdvReport{
				Id:         easygo.NewString(id),
				AdvId:      easygo.NewInt64(item.GetAdvId()),
				CreateTime: easygo.NewInt64(easygo.Get0ClockTimestamp(item.GetOpTime())),
			}
		}

		switch item.GetOpType() {
		case for_game.ADV_LOG_OP_TYPE_1:
			report.PvCount = easygo.NewInt64(report.GetPvCount() + 1)
		case for_game.ADV_LOG_OP_TYPE_2:
			report.UvCount = easygo.NewInt64(report.GetUvCount() + 1)
		case for_game.ADV_LOG_OP_TYPE_3:
			report.Clicks = easygo.NewInt64(report.GetClicks() + 1)
		case for_game.ADV_LOG_OP_TYPE_4:
			report.ClickPlayers = easygo.NewInt64(report.GetClickPlayers() + 1)
		}

		UpdateNearbyAdvReport(report)
	}

	//写报表生成进度
	MakeReportJob(for_game.TABLE_NEARBY_ADV_REPORT, endTime)
}

//初始化附近的人引导项报表
func InitNearbyAdvReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_NEARBY_ADV_REPORT)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//更新附近的人引导项报表
func UpdateNearbyAdvReport(report *share_message.AdvReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_NEARBY_ADV_REPORT)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetId()}, bson.M{"$set": report})
	easygo.PanicError(err)
	logs.Info("更新广告报表")
}

//查询附近的人引导项报表
func GetNearbyAdvReportById(id string) *share_message.AdvReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_NEARBY_ADV_REPORT)
	defer closeFun()

	var obj *share_message.AdvReport
	err := col.Find(bson.M{"_id": id}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//查询附近的人引导项报表
func GetNearbyAdvReportList(reqMsg *brower_backstage.ListRequest) ([]*share_message.NearReport, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_NEARBY_ADV_REPORT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		var ids []int64
		lis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_NEAR_LEAD, bson.M{"Name": bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}}, 0, 0)
		for _, i := range lis {
			ids = append(ids, i.(bson.M)["_id"].(int64))
		}
		queryBson["AdvId"] = bson.M{"$in": ids}

	}

	//联表查询
	m := []bson.M{
		{"$lookup": bson.M{
			"from":         for_game.TABLE_NEAR_LEAD,
			"localField":   "AdvId",
			"foreignField": "_id",
			"as":           "Near",
		}},
		{"$match": queryBson},
		{"$unwind": "$Near"},
		{"$sort": bson.M{"CreateTime": -1}},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	query := col.Pipe(m)
	count, err := col.Count()
	easygo.PanicError(err)

	var list []*share_message.NearReport
	errc := query.All(&list)
	easygo.PanicError(errc)

	return list, count
}

//=====================================================================================================start
//报表
func BakReport() {
	startTime := easygo.Get0ClockMillTimestamp(easygo.NowTimestamp() - 86400)
	endTime := easygo.Get24ClockMillTimestamp(easygo.NowTimestamp() - 86400)
	regCount := GetCountForPlayer(bson.M{"CreateTime": bson.M{"$gte": startTime, "$lt": endTime}})                //注册人数
	regValCount := GetCountForPlayer(bson.M{"CreateTime": bson.M{"$gte": startTime, "$lt": endTime}, "Types": 6}) //新增有效注册人数

	one := GetRegisterLoginReportByTime(startTime) //埋点报表
	if one != nil {
		one.RegSumCount = easygo.NewInt64(regCount)
		if regValCount > 0 {
			one.PhoneRegCount = easygo.NewInt64(one.GetPhoneRegCount() + regValCount)
			one.ValidRegSumCount = easygo.NewInt64(one.GetValidRegSumCount() + regValCount)
			one.ValidPhoneRegCount = easygo.NewInt64(one.GetValidPhoneRegCount() + regValCount)
			one.LoginSumCount = easygo.NewInt64(one.GetLoginSumCount() + regValCount)
			one.LoginTimesCount = easygo.NewInt64(one.GetLoginTimesCount() + regValCount)
			// one.PhoneLoginCount = easygo.NewInt64(one.GetPhoneLoginCount() + regValCount)
			// one.ActDevCount = easygo.NewInt64(one.GetActDevCount() + (regValCount + int64(util.RandIntn(int(regValCount)))))
			// one.ValidActDevCount = easygo.NewInt64(one.GetValidActDevCount() + regValCount)
		}

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, "bak_report_register_login")
		count, _ := col.Count()
		if count == 0 {
			CopyRegLogTable()
		}

		_, err := col.Upsert(bson.M{"_id": one.GetCreateTime()}, bson.M{"$set": one})
		easygo.PanicError(err)
		closeFun()
	}
}

//查询指定玩家数量
func GetCountForPlayer(queryBson bson.M) int64 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	return int64(count)
}

//注册登录日志表
func CopyRegLogTable() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, "report_register_login")
	defer closeFun()
	log := &share_message.RegisterLoginReport{}
	iter := col.Find(bson.M{}).Iter() //mongodb迭代器
	col1, closeFun1 := MongoMgr.GetC(for_game.MONGODB_NINGMENG, "bak_report_register_login")
	for iter.Next(&log) {
		col1.Upsert(bson.M{"_id": log.GetCreateTime()}, bson.M{"$set": log})
	}
	closeFun1()
	if err := iter.Close(); err != nil {
		easygo.PanicError(err)
	}
}

//=====================================================================================================虚拟商城报表
func MakeCoinProductReport() {
	startTime := easygo.GetToday0ClockTimestamp()
	endTime := easygo.GetToday24ClockTimestamp()
	if for_game.IS_FORMAL_SERVER {
		startTime -= easygo.A_DAY_SECOND
		endTime -= easygo.A_DAY_SECOND
	}

	queryBson := bson.M{"CreateTime": bson.M{"$gte": startTime, "$lt": endTime}, "GetType": bson.M{"$ne": for_game.COIN_ITEM_GETTYPE_BACK}, "ProductId": bson.M{"$gte": 0}}
	getPropsLog, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_GETPROPS_LOG, queryBson, 0, 0)
	if count == 0 {
		return
	}

	dataMap := make(map[string]*share_message.CoinProductReport)
	for _, li := range getPropsLog {
		log := &share_message.PlayerGetPropsLog{}
		for_game.StructToOtherStruct(li, log)
		productId := log.GetProductId()
		reportId := easygo.AnytoA(startTime) + easygo.AnytoA(productId)
		oldreport := dataMap[reportId]
		report := &share_message.CoinProductReport{}
		if oldreport == nil {
			report = &share_message.CoinProductReport{
				Id:         easygo.NewString(reportId),
				CreateTime: easygo.NewInt64(startTime),
				ProductId:  easygo.NewInt64(productId),
			}
			switch log.GetGetType() {
			case for_game.COIN_ITEM_GETTYPE_BUY:
				report.BuyNum = easygo.NewInt64(1)
				report.BuyCount = easygo.NewInt64(1)
				switch log.GetBuyWay() {
				case for_game.COIN_PROPS_BUYWAY_COIN:
					coinLog := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_COINCHANGELOG, bson.M{"Extend.OrderId": log.GetOrderId()})
					if coinLog != nil {
						report.CoinSum = easygo.NewInt64(-coinLog.(bson.M)["ChangeCoin"].(int64))
					}
				case for_game.COIN_PROPS_BUYWAY_MONEY:
					goldLog := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_GOLDCHANGELOG, bson.M{"Extend.OrderId": log.GetOrderId()})
					if goldLog != nil {
						report.GoldSum = easygo.NewInt64(-goldLog.(bson.M)["ChangeGold"].(int64))
					}
				}
			case for_game.COIN_ITEM_GETTYPE_SEND:
				report.GiveNum = easygo.NewInt64(1)
				report.GiveCount = easygo.NewInt64(1)
			case for_game.COIN_ITEM_GETTYPE_PLAYER_SEND:
				report.UserGiveNum = easygo.NewInt64(1)
				report.UserGiveCount = easygo.NewInt64(1)
			case for_game.COIN_ITEM_GETTYPE_ACTIVITY:
				report.ActGiveNum = easygo.NewInt64(1)
				report.ActGiveCount = easygo.NewInt64(1)
			}
		} else {
			for_game.StructToOtherStruct(oldreport, report)
			switch log.GetGetType() {
			case for_game.COIN_ITEM_GETTYPE_BUY:
				report.BuyNum = easygo.NewInt64(report.GetBuyNum() + 1)
				report.BuyCount = easygo.NewInt64(report.GetBuyCount() + 1)
				switch log.GetBuyWay() {
				case for_game.COIN_PROPS_BUYWAY_COIN:
					coinLog := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_COINCHANGELOG, bson.M{"Extend.OrderId": log.GetOrderId()})
					if coinLog != nil {
						report.CoinSum = easygo.NewInt64(report.GetCoinSum() - coinLog.(bson.M)["ChangeCoin"].(int64))
					}
				case for_game.COIN_PROPS_BUYWAY_MONEY:
					goldLog := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_GOLDCHANGELOG, bson.M{"Extend.OrderId": log.GetOrderId()})
					if goldLog != nil {
						report.GoldSum = easygo.NewInt64(report.GetGoldSum() - goldLog.(bson.M)["ChangeGold"].(int64))
					}
				}
			case for_game.COIN_ITEM_GETTYPE_SEND:
				report.GiveNum = easygo.NewInt64(report.GetGiveNum() + 1)
				report.GiveCount = easygo.NewInt64(report.GetGiveCount() + 1)
			case for_game.COIN_ITEM_GETTYPE_PLAYER_SEND:
				report.UserGiveNum = easygo.NewInt64(report.GetUserGiveNum() + 1)
				report.UserGiveCount = easygo.NewInt64(report.GetUserGiveCount() + 1)
			case for_game.COIN_ITEM_GETTYPE_ACTIVITY:
				report.ActGiveNum = easygo.NewInt64(report.GetActGiveNum() + 1)
				report.ActGiveCount = easygo.NewInt64(report.GetActGiveCount() + 1)
			}
		}
		dataMap[reportId] = report
	}
	//校正购买人数和赠送人数
	for _, dli := range dataMap {
		//人数字段要重新获取真实人数
		if dli.GetBuyNum() > 1 {
			findBson := bson.M{"CreateTime": bson.M{"$gte": startTime, "$lt": endTime}, "GetType": for_game.COIN_ITEM_GETTYPE_BUY, "ProductId": dli.GetProductId()}
			m := []bson.M{
				{"$match": findBson},
				{"$group": bson.M{"_id": "$PlayerId"}},
			}
			ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_GETPROPS_LOG, m, 0, 0)
			dli.BuyNum = easygo.NewInt64(len(ls))
		}
		if dli.GetGiveNum() > 1 {
			findBson := bson.M{"CreateTime": bson.M{"$gte": startTime, "$lt": endTime}, "GetType": for_game.COIN_ITEM_GETTYPE_SEND, "ProductId": dli.GetProductId()}
			m := []bson.M{
				{"$match": findBson},
				{"$group": bson.M{"_id": "$PlayerId"}},
			}
			ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_GETPROPS_LOG, m, 0, 0)
			dli.GiveNum = easygo.NewInt64(len(ls))
		}
		if dli.GetUserGiveNum() > 1 {
			findBson := bson.M{"CreateTime": bson.M{"$gte": startTime, "$lt": endTime}, "GetType": for_game.COIN_ITEM_GETTYPE_PLAYER_SEND, "ProductId": dli.GetProductId()}
			m := []bson.M{
				{"$match": findBson},
				{"$group": bson.M{"_id": "$PlayerId"}},
			}
			ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_GETPROPS_LOG, m, 0, 0)
			dli.UserGiveNum = easygo.NewInt64(len(ls))
		}
		if dli.GetActGiveNum() > 1 {
			findBson := bson.M{"CreateTime": bson.M{"$gte": startTime, "$lt": endTime}, "GetType": for_game.COIN_ITEM_GETTYPE_ACTIVITY, "ProductId": dli.GetProductId()}
			m := []bson.M{
				{"$match": findBson},
				{"$group": bson.M{"_id": "$PlayerId"}},
			}
			ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_GETPROPS_LOG, m, 0, 0)
			dli.ActGiveNum = easygo.NewInt64(len(ls))
		}

		//礼包类型要除掉礼包的多个道具重复次数
		product := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT, bson.M{"_id": dli.GetProductId()})
		if product != nil {
			propsId := product.(bson.M)["PropsId"].(int64)
			iprops := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_PROPS_ITEM, bson.M{"_id": propsId})
			if iprops != nil {
				props := &share_message.PropsItem{}
				for_game.StructToOtherStruct(iprops, props)
				if props.GetPropsType() == for_game.COIN_PROPS_TYPE_LB {
					checkItems := []int64{}
					err := json.Unmarshal([]byte(props.GetUseValue()), &checkItems)
					easygo.PanicError(err)
					nmb := len(checkItems)
					dli.BuyCount = easygo.NewInt64(dli.GetBuyCount() / int64(nmb))
					dli.GiveCount = easygo.NewInt64(dli.GetGiveCount() / int64(nmb))
					dli.UserGiveCount = easygo.NewInt64(dli.GetUserGiveCount() / int64(nmb))
					dli.ActGiveCount = easygo.NewInt64(dli.GetActGiveCount() / int64(nmb))
					dli.GoldSum = easygo.NewInt64(dli.GetGoldSum() / int64(nmb))
					dli.CoinSum = easygo.NewInt64(dli.GetCoinSum() / int64(nmb))
				}
			}

		}
	}

	UpsertCoinProductReport(dataMap)
}

//批量更新虚拟商城报表
func UpsertCoinProductReport(lis map[string]*share_message.CoinProductReport) {
	var data []interface{}
	for _, v := range lis {
		b1 := bson.M{"_id": v.GetId()}
		b2 := v
		data = append(data, b1, b2)
	}
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT_REPORT, data)
}

//==========================================注册登录页面埋点报表
func MakePageRegLogReport() {
	list := for_game.GetChannelListNopage()
	for _, li := range list {
		SavePageRegLogReport(1, li.GetKey())
		SavePageRegLogReport(2, li.GetKey())
	}
	SavePageRegLogReport(1, "")
	SavePageRegLogReport(2, "")

}

func SavePageRegLogReport(dicType int32, Channel string) {
	queryTime := easygo.GetYesterday0ClockTimestamp()
	endTime := easygo.GetYesterday24ClockTimestamp()
	var loginTimes int64 = 0
	var loginCount int64 = 0
	var oneLoginTimes int64 = 0
	var oneLoginCount int64 = 0
	var wxLoginTimes int64 = 0
	var wxLoginCount int64 = 0
	var phoneRegTimes int64 = 0
	var phoneRegCount int64 = 0
	var regCodeTimes int64 = 0
	var regCodeCount int64 = 0
	var useInfoTimes int64 = 0
	var useInfoCount int64 = 0
	var intWallTimes int64 = 0
	var intWallCount int64 = 0
	var recPageTimes int64 = 0
	var recPageCount int64 = 0
	var actDevCount int64 = 0
	var validActDevCount int64 = 0

	timeBson := bson.M{"$gte": queryTime * 1000, "$lt": endTime * 1000}
	queryBson := bson.M{"DicType": dicType, "Channel": Channel, "CreateTime": timeBson}
	tM := []bson.M{
		{"$match": queryBson},
		{"$group": bson.M{"_id": "$Type", "Count": bson.M{"$sum": 1}}},
	}
	list := for_game.FindPipeAll(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PAGE_REGLOG, tM, 0, 0)
	for _, li := range list {
		one := &share_message.PipeIntCount{}
		for_game.StructToOtherStruct(li, one)
		switch one.GetId() {
		/*Type类型按钮值
		1-登录页面浏览次数
		2-一键登录页面浏览次数
		3-微信绑定页浏览次数
		4-手机号注册页浏览次数
		5-验证码填写页浏览次数
		6-个人信息页浏览次数
		7-兴趣墙页的浏览次数
		8-推荐页的浏览次数
		*/
		case 1:
			loginTimes = one.GetCount()
		case 2:
			oneLoginTimes = one.GetCount()
		case 3:
			wxLoginTimes = one.GetCount()
		case 4:
			phoneRegTimes = one.GetCount()
		case 5:
			regCodeTimes = one.GetCount()
		case 6:
			useInfoTimes = one.GetCount()
		case 7:
			intWallTimes = one.GetCount()
		case 8:
			recPageTimes = one.GetCount()
		}
	}

	for i := 1; i < 9; i++ {
		cM := []bson.M{
			{"$match": bson.M{"DicType": dicType, "Channel": Channel, "Type": i, "CreateTime": timeBson}},
			{"$group": bson.M{"_id": "$Code", "Count": bson.M{"$sum": 1}}},
			{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
		}
		count := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PAGE_REGLOG, cM)
		switch i {
		case 1:
			loginCount = count
		case 2:
			oneLoginCount = count
		case 3:
			wxLoginCount = count
		case 4:
			phoneRegCount = count
		case 5:
			regCodeCount = count
		case 6:
			useInfoCount = count
		case 7:
			intWallCount = count
		case 8:
			recPageCount = count
		}
	}

	report := &share_message.PageRegLogReport{
		Id:            easygo.NewString(easygo.AnytoA(queryTime) + Channel + easygo.AnytoA(dicType)),
		CreateTime:    easygo.NewInt64(queryTime),
		Channel:       easygo.NewString(Channel),
		DicType:       easygo.NewInt32(dicType),
		LoginTimes:    easygo.NewInt64(loginTimes),
		LoginCount:    easygo.NewInt64(loginCount),
		OneLoginTimes: easygo.NewInt64(oneLoginTimes),
		OneLoginCount: easygo.NewInt64(oneLoginCount),
		WxLoginTimes:  easygo.NewInt64(wxLoginTimes),
		WxLoginCount:  easygo.NewInt64(wxLoginCount),
		PhoneRegTimes: easygo.NewInt64(phoneRegTimes),
		PhoneRegCount: easygo.NewInt64(phoneRegCount),
		RegCodeTimes:  easygo.NewInt64(regCodeTimes),
		RegCodeCount:  easygo.NewInt64(regCodeCount),
		UseInfoTimes:  easygo.NewInt64(useInfoTimes),
		UseInfoCount:  easygo.NewInt64(useInfoCount),
		IntWallTimes:  easygo.NewInt64(intWallTimes),
		IntWallCount:  easygo.NewInt64(intWallCount),
		RecPageTimes:  easygo.NewInt64(recPageTimes),
		RecPageCount:  easygo.NewInt64(recPageCount),
	}

	if dicType == 1 {
		lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_POS_DEVICECODE, bson.M{"Channle": Channel, "_id": timeBson}, 0, 0)
		actDevCount = int64(count)
		for _, l := range lis {
			players := for_game.QueryPlayersByDeviceCode(l.(bson.M)["DeviceCode"].(string))
			if len(players) > 0 {
				validActDevCount += 1
			}
		}
		report.ActDevCount = easygo.NewInt64(actDevCount)
		report.ValidActDevCount = easygo.NewInt64(validActDevCount)
	}

	if loginTimes+loginCount+oneLoginTimes+oneLoginCount+wxLoginTimes+wxLoginCount+phoneRegTimes+phoneRegCount+regCodeTimes+regCodeCount+useInfoTimes+useInfoCount+intWallTimes+intWallCount+recPageTimes+recPageCount+actDevCount+validActDevCount > 0 {
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_PAGE_REGLOG_REPORT, bson.M{"_id": report.GetId()}, report, true)
	}
}

//=====================================================================================================注单统计报表
//生成注单统计报表
func MakeBetSlipReport() {
	startTime := easygo.GetToday0ClockTimestamp()
	endTime := easygo.GetToday24ClockTimestamp()
	if for_game.IS_FORMAL_SERVER {
		startTime -= easygo.A_DAY_SECOND
		endTime -= easygo.A_DAY_SECOND
	} else {
		for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_BET_SLIP_REPORT, bson.M{"CreateTime": startTime})
	}

	queryBson := bson.M{"CreateTime": bson.M{"$gte": startTime, "$lt": endTime}}
	betSlips, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, queryBson, 0, 0)
	if count == 0 {
		return
	}

	dataMap := make(map[int64]*share_message.BetSlipReport)
	for _, li := range betSlips {
		log := &share_message.TableESPortsGuessBetRecord{}
		for_game.StructToOtherStruct(li, log)
		appLabelId := log.GetAppLabelId()
		reportId := easygo.StringToInt64noErr(easygo.AnytoA(startTime) + easygo.AnytoA(appLabelId))
		oldreport := dataMap[reportId]
		report := &share_message.BetSlipReport{}
		if oldreport == nil {
			report = &share_message.BetSlipReport{
				Id:         easygo.NewInt64(reportId),
				CreateTime: easygo.NewInt64(startTime),
				AppLabelID: easygo.NewInt64(appLabelId),
			}
		} else {
			for_game.StructToOtherStruct(oldreport, report)
		}
		report.BetAmount = easygo.NewInt64(report.GetBetAmount() + log.GetBetAmount())
		report.SuccessAmount = easygo.NewInt64(report.GetSuccessAmount() + log.GetSuccessAmount())
		report.FailAmount = easygo.NewInt64(report.GetFailAmount() + log.GetFailAmount())
		report.DisableAmount = easygo.NewInt64(report.GetDisableAmount() + log.GetDisableAmount())
		report.IllegalAmount = easygo.NewInt64(report.GetIllegalAmount() + log.GetIllegalAmount())
		report.SumAmount = easygo.NewInt64(report.GetSumAmount() + log.GetBetAmount() - log.GetSuccessAmount() - log.GetDisableAmount())
		dataMap[reportId] = report
	}

	UpsertBetSlipReport(dataMap)
}

//批量更新注单统计
func UpsertBetSlipReport(lis map[int64]*share_message.BetSlipReport) {
	var data []interface{}
	for _, v := range lis {
		b1 := bson.M{"_id": v.GetId()}
		b2 := v
		data = append(data, b1, b2)
	}
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_BET_SLIP_REPORT, data)
}

//====================================================================================================电竞埋点报表
type PointReportComm struct {
	StartTime int64
	EndTime   int64
	// 时间查询条件
	TimeBson bson.M
	// APP登录人数
	LoginCount int64
	// 底部按钮点击总次数
	FootBtnCLick int64
	// 底部电竞按钮点击人数
	FootEsBtnCount int64
	// 底部电竞按钮点击次数
	FootEsBtnCLick int64
	// 单日使用电竞总时长
	EsUseTime int64
}

// 生成电竞埋点报表入口方法
func MakeEsportPointReport(types string) {
	var startTime, endTime TIME_64
	switch types {
	case "day":
		logs.Info("生成电竞埋点日报表")
		startTime = easygo.GetYesterday0ClockTimestamp()
		endTime = easygo.GetYesterday24ClockTimestamp()
	case "week":
		logs.Info("生成电竞埋点周报表")
		startTime = easygo.GetWeek0ClockOfTimestamp(easygo.GetYesterday0ClockTimestamp())
		endTime = easygo.GetToday0ClockTimestamp()
		if easygo.GetDifferenceDay(startTime, endTime) < 7 {
			return
		}
	case "month":
		logs.Info("生成电竞埋点月报表")
		startTime = easygo.GetMonth0ClockOfTimestamp(easygo.GetYesterday0ClockTimestamp())
		endTime = easygo.GetToday0ClockTimestamp()
		if easygo.GetDifferenceDay(startTime, endTime) < 28 {
			return
		}
	default:
		logs.Error("请求生成电竞报表类型错误")
		return
	}

	// 时间查询条件
	timeBson := bson.M{"$gte": startTime, "$lt": endTime}
	// APP登录人数
	lcM := []bson.M{
		{"$match": bson.M{"Type": bson.M{"$lt": for_game.LOGINREGISTER_PHONEREGISTER}, "Time": bson.M{"$gte": startTime * 1000, "$lt": endTime * 1000}}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	loginCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_LOGIN_REGISTER_LOG, lcM)
	// 底部按钮点击总次数
	fbcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_1, "ButtonId": 0, "CreateTime": timeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
	}
	footBtnCLick := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, fbcM)
	// 底部电竞按钮点击人数
	febcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_1, "NavigationId": for_game.ESPORT_MODLE_4, "ButtonId": 0, "CreateTime": timeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	footEsBtnCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, febcM)
	// 底部电竞按钮点击次数
	febccM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_1, "NavigationId": for_game.ESPORT_MODLE_4, "ButtonId": 0, "CreateTime": timeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
	}
	footEsBtnCLick := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, febccM)
	// 单日使用电竞总时长
	eutM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_1, "CreateTime": timeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$Duration"}}},
	}
	esUseTime := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_DURATION_LOG, eutM)

	prc := &PointReportComm{
		StartTime:      startTime,
		EndTime:        endTime,
		TimeBson:       timeBson,
		LoginCount:     loginCount,
		FootBtnCLick:   footBtnCLick,
		FootEsBtnCount: footEsBtnCount,
		FootEsBtnCLick: footEsBtnCLick,
		EsUseTime:      esUseTime,
	}
	prc.MakeBasisPointsReport(types)
	prc.MakeMenuPointsReport(types)
	prc.MakeLabelPointsReport(types)
	prc.MakeNewsAmusePointsReport(types)
	prc.MakeVdoHallPointsReport(types)
	prc.MakeApplyVdoHallPointsReport(types)
	prc.MakeMatchLsPointsReport(types)
	prc.MakeMatchDilPointsReport(types)
	prc.MakeGuessPointsReport(types)
	prc.MakeMsgPointsReport(types)
	prc.MakeEsportCoinPointsReport(types)
}

//生成基础埋点报表
func (p *PointReportComm) MakeBasisPointsReport(types string) {
	// 浏览页面总数
	psM := []bson.M{
		{"$match": bson.M{"NavigationId": for_game.ESPORT_MODLE_4, "ButtonId": bson.M{"$lt": 1}, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
	}
	pvSum := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, psM)

	var nextDayCount, threeDayCount, sevenDayCount, nextNewCount, threeNewCount, sevenNewCount int = 0, 0, 0, 0, 0, 0
	var loyalUserCount int64 = 0
	var reportTable string
	switch types {
	case "day":
		reportTable = for_game.TABLE_ESPORTS_BASIS_POINTS_REPORT_DAY
		// 次留人数
		loginTimeBson := bson.M{"$gte": p.StartTime, "$lt": p.EndTime}
		nextDayCreateTimeBson := bson.M{"$gte": p.StartTime - for_game.DAY_SECOND, "$lt": p.EndTime - for_game.DAY_SECOND}
		nextDayCount = for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_PLAYER, bson.M{"LastLoginTime": loginTimeBson, "CreateTime": nextDayCreateTimeBson})
		// 3留人数
		threeDayCreateTimeBson := bson.M{"$gte": p.StartTime - 3*for_game.DAY_SECOND, "$lt": p.EndTime - 3*for_game.DAY_SECOND}
		threeDayCount = for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_PLAYER, bson.M{"LastLoginTime": loginTimeBson, "CreateTime": threeDayCreateTimeBson})
		// 7留人数
		sevenDayCreateTimeBson := bson.M{"$gte": p.StartTime - 7*for_game.DAY_SECOND, "$lt": p.EndTime - 7*for_game.DAY_SECOND}
		sevenDayCount = for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_PLAYER, bson.M{"LastLoginTime": loginTimeBson, "CreateTime": sevenDayCreateTimeBson})

		nextNewCount = for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_PLAYER, bson.M{"CreateTime": nextDayCreateTimeBson})
		threeNewCount = for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_PLAYER, bson.M{"CreateTime": threeDayCreateTimeBson})
		sevenNewCount = for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_PLAYER, bson.M{"CreateTime": sevenDayCreateTimeBson})
	case "week":
		reportTable = for_game.TABLE_ESPORTS_BASIS_POINTS_REPORT_WEEK
	case "month":
		reportTable = for_game.TABLE_ESPORTS_BASIS_POINTS_REPORT_MONTH
		// 底部电竞按钮点击人数
		lucM := []bson.M{
			{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_1, "NavigationId": for_game.ESPORT_MODLE_4, "ButtonId": 0, "CreateTime": p.TimeBson}},
			{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
			{"$match": bson.M{"Count": bson.M{"$gte": 14}}}, //查询登录次数大于2周的用户
			{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
		}
		loyalUserCount = for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, lucM)

	}

	perAvgTime := int64(0)
	if p.FootEsBtnCount > 0 {
		perAvgTime = p.EsUseTime / p.FootEsBtnCount
	}
	sigAvgTime := int64(0)
	if p.FootEsBtnCLick > 0 {
		sigAvgTime = p.EsUseTime / p.FootEsBtnCLick
	}

	report := &share_message.BasisPointsReport{
		CreateTime:     easygo.NewInt64(p.StartTime),
		SumClick:       easygo.NewInt64(p.FootBtnCLick),
		EsClickCount:   easygo.NewInt64(p.FootEsBtnCount),
		EsClick:        easygo.NewInt64(p.FootEsBtnCLick),
		PerAvgTime:     easygo.NewInt64(perAvgTime),
		SigAvgTime:     easygo.NewInt64(sigAvgTime),
		LoyalUserCount: easygo.NewInt64(loyalUserCount),
		PvSum:          easygo.NewInt64(pvSum),
		NextDayCount:   easygo.NewInt64(nextDayCount),
		ThreeDayCount:  easygo.NewInt64(threeDayCount),
		SevenDayCount:  easygo.NewInt64(sevenDayCount),
		LoginCount:     easygo.NewInt64(p.LoginCount),
		NextNewCount:   easygo.NewInt64(nextNewCount),
		ThreeNewCount:  easygo.NewInt64(threeNewCount),
		SevenNewCount:  easygo.NewInt64(sevenNewCount),
	}
	if perAvgTime+sigAvgTime+loyalUserCount+pvSum+int64(nextDayCount)+int64(threeDayCount)+int64(sevenDayCount) > 0 {
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report}, true)
	}
}

//Tab菜单埋点报表
func (p *PointReportComm) MakeMenuPointsReport(types string) {
	// 资讯点击人数
	var newsCount int64 = 0
	// 娱乐点击人数
	var amuseCount int64 = 0
	// 放映厅点击人数
	var vdoHallCount int64 = 0
	// 赛事点击人数
	var matchCount int64 = 0

	nM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "NavigationId": for_game.ESPORT_MODLE_4, "MenuId": for_game.ESPORTMENU_REALTIME, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId"}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	newsCount = for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, nM)

	aM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "NavigationId": for_game.ESPORT_MODLE_4, "MenuId": for_game.ESPORTMENU_RECREATION, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId"}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	amuseCount = for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, aM)

	vM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "NavigationId": for_game.ESPORT_MODLE_4, "MenuId": for_game.ESPORTMENU_LIVE, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId"}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	vdoHallCount = for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, vM)

	mM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "NavigationId": for_game.ESPORT_MODLE_4, "MenuId": for_game.ESPORTMENU_GAME, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId"}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	matchCount = for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, mM)

	report := &share_message.MenuPointsReport{
		CreateTime:   easygo.NewInt64(p.StartTime),
		NewsCount:    easygo.NewInt64(newsCount),
		AmuseCount:   easygo.NewInt64(amuseCount),
		VdoHallCount: easygo.NewInt64(vdoHallCount),
		MatchCount:   easygo.NewInt64(matchCount),
		EsClickCount: easygo.NewInt64(p.FootEsBtnCount),
	}
	var reportTable string
	switch types {
	case "day":
		reportTable = for_game.TABLE_ESPORTS_MENU_POINTS_REPORT_DAY
	case "week":
		reportTable = for_game.TABLE_ESPORTS_MENU_POINTS_REPORT_WEEK
	case "month":
		reportTable = for_game.TABLE_ESPORTS_MENU_POINTS_REPORT_MONTH
	}
	if newsCount+amuseCount+vdoHallCount+matchCount > 0 {
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report}, true)
	}
}

//标签埋点报表
func (p *PointReportComm) MakeLabelPointsReport(types string) {
	var menuIds []int32
	menuIds = append(menuIds, for_game.ESPORTMENU_REALTIME, for_game.ESPORTMENU_RECREATION, for_game.ESPORTMENU_GAME, for_game.ESPORTMENU_LIVE)

	var flisLab []int64                                           //固定标签
	var labMap = make(map[int64]*share_message.TableESPortsLabel) //标签map
	labelList, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_LABEL, bson.M{}, 0, 0)
	for _, lab := range labelList {
		one := &share_message.TableESPortsLabel{}
		for_game.StructToOtherStruct(lab, one)
		//不等于2是系统固定标签类型
		if one.GetLabelType() != 2 {
			flisLab = append(flisLab, one.GetId())
		}
		labMap[one.GetId()] = one
	}

	var sum int64 = 0
	for _, mid := range menuIds {
		lM := []bson.M{
			{"$match": bson.M{"CreateTime": p.TimeBson, "MenuId": mid, "PageType": for_game.ESPORT_BPS_PAGE_TYPE_3, "ButtonId": 0, "LabelId": -1000}},
			{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		}
		labMagBtnCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, lM)

		cM := []bson.M{
			{"$match": bson.M{"CreateTime": p.TimeBson, "MenuId": mid, "PageType": for_game.ESPORT_BPS_PAGE_TYPE_3, "LabelId": -1000, "ButtonId": 1}},
			{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		}
		labChaBtnCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, cM)

		var clisLab []int64 //当前菜单的自定义标签
		for _, clab := range labMap {
			if clab.GetLabelType() == 2 && clab.GetMenuId() == mid {
				clisLab = append(clisLab, clab.GetId())
			}
		}

		var clis []*share_message.LabelPointsstruct

		for _, c := range clisLab {
			M := []bson.M{
				{"$match": bson.M{"CreateTime": p.TimeBson, "MenuId": mid, "PageType": for_game.ESPORT_BPS_PAGE_TYPE_3, "ButtonId": 0, "LabelId": c}},
				{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
				{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
			}
			count := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, M)
			cone := &share_message.LabelPointsstruct{
				Id:    easygo.NewInt64(c),
				Title: easygo.NewString(labMap[c].GetTitle()),
				Count: easygo.NewInt64(count),
			}
			clis = append(clis, cone)
			sum += count
		}

		var reportTable string
		switch types {
		case "day":
			reportTable = for_game.TABLE_ESPORTS_LABEL_POINTS_REPORT_DAY
		case "week":
			reportTable = for_game.TABLE_ESPORTS_LABEL_POINTS_REPORT_WEEK
		case "month":
			reportTable = for_game.TABLE_ESPORTS_LABEL_POINTS_REPORT_MONTH
		}

		var flis []*share_message.LabelPointsstruct
		for _, f := range flisLab {
			M := []bson.M{
				{"$match": bson.M{"CreateTime": p.TimeBson, "MenuId": mid, "PageType": for_game.ESPORT_BPS_PAGE_TYPE_3, "ButtonId": 0, "LabelId": f}},
				{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
				{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
			}
			count := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, M)
			fone := &share_message.LabelPointsstruct{
				Id:    easygo.NewInt64(f),
				Title: easygo.NewString(labMap[f].GetTitle()),
				Count: easygo.NewInt64(count),
			}
			flis = append(flis, fone)
			sum += count
		}

		repoertId := easygo.AnytoA(p.StartTime) + easygo.AnytoA(mid)
		report := &share_message.LabelPointsReport{
			Id:             easygo.NewString(repoertId),
			CreateTime:     easygo.NewInt64(p.StartTime),
			MenuId:         easygo.NewInt32(mid),
			EsClickCount:   easygo.NewInt64(p.FootEsBtnCount),
			LabMagBtnCount: easygo.NewInt64(labMagBtnCount),
			LabChaBtnCount: easygo.NewInt64(labChaBtnCount),
			Custom:         clis,
			Fixed:          flis,
		}
		if sum+labMagBtnCount+labChaBtnCount > 0 {
			for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": repoertId}, bson.M{"$set": report}, true)
		}
	}
}

//资讯娱乐埋点报表
func (p *PointReportComm) MakeNewsAmusePointsReport(types string) {
	var reportTable string
	switch types {
	case "day":
		reportTable = for_game.TABLE_ESPORTS_NEWS_AMUSE_POINTS_REPORT_DAY
	case "week":
		reportTable = for_game.TABLE_ESPORTS_NEWS_AMUSE_POINTS_REPORT_WEEK
	case "month":
		reportTable = for_game.TABLE_ESPORTS_NEWS_AMUSE_POINTS_REPORT_MONTH
	}

	var menuIds []int32
	menuIds = append(menuIds, for_game.ESPORTMENU_REALTIME, for_game.ESPORTMENU_RECREATION)
	for _, mid := range menuIds {
		eutM := []bson.M{
			{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "MenuId": mid, "CreateTime": p.TimeBson}},
			{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$Duration"}}},
		}
		esUseTime := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_DURATION_LOG, eutM)

		mbcM := []bson.M{
			{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "MenuId": mid, "ButtonId": 0, "CreateTime": p.TimeBson}},
			{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
		}
		menuBtnCLick := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, mbcM)

		rsM := []bson.M{
			{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_5, "MenuId": mid, "ButtonId": 0, "CreateTime": p.TimeBson}}, //ButtonId=0表示非按钮
			{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
		}
		readSum := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, rsM)

		avgTime := int64(0)
		if p.FootEsBtnCount > 0 {
			avgTime = esUseTime / p.FootEsBtnCount
		}
		sigAvgTime := int64(0)
		if menuBtnCLick > 0 {
			sigAvgTime = esUseTime / menuBtnCLick
		}

		repoertId := easygo.AnytoA(p.StartTime) + easygo.AnytoA(mid)
		report := &share_message.NewsAmusePointsReport{
			Id:         easygo.NewString(repoertId),
			CreateTime: easygo.NewInt64(p.StartTime),
			MenuId:     easygo.NewInt32(mid),
			AvgTime:    easygo.NewInt64(avgTime),
			SigAvgTime: easygo.NewInt64(sigAvgTime),
			ReadSum:    easygo.NewInt64(readSum),
		}
		if avgTime+sigAvgTime+readSum > 0 {
			for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": repoertId}, bson.M{"$set": report}, true)
		}
	}
}

//放映厅埋点报表
func (p *PointReportComm) MakeVdoHallPointsReport(types string) {
	eutM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "MenuId": for_game.ESPORTMENU_LIVE, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$Duration"}}},
	}
	esUseTime := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_DURATION_LOG, eutM)

	rsM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_5, "MenuId": for_game.ESPORTMENU_LIVE, "ButtonId": 0, "CreateTime": p.TimeBson}}, //ButtonId=0表示非按钮
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
	}
	readSum := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, rsM)

	csM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "MenuId": for_game.ESPORTMENU_LIVE, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
	}
	clickSum := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, csM)

	eucM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "MenuId": for_game.ESPORTMENU_LIVE, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	esUserCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, eucM)

	oucM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_5, "MenuId": for_game.ESPORTMENU_LIVE, "DataType": 3, "ButtonId": 1, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	openUserCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, oucM)

	vlcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_4, "MenuId": for_game.ESPORTMENU_LIVE, "ExTabId": 1, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	vdoLsCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, vlcM)

	mfcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_4, "MenuId": for_game.ESPORTMENU_LIVE, "ExTabId": 2, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	myFollowCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, mfcM)

	plcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_4, "MenuId": for_game.ESPORTMENU_LIVE, "ExTabId": 3, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	playLogCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, plcM)

	fhcM := []bson.M{
		{"$match": bson.M{"CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	followHallCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_FLOW_LIVE_FOLLOW_HISTORY, fhcM)

	fucM := []bson.M{
		{"$match": bson.M{"CreateTime": p.TimeBson, "Source": 1}}, //电竞页的主播关注按钮为1
		{"$group": bson.M{"_id": "$OperateId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	followUserCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_ATTENTION, fucM)

	avgTime := int64(0)
	if p.FootEsBtnCount > 0 {
		avgTime = esUseTime / p.FootEsBtnCount
	}
	sigAvgTime := int64(0)
	if readSum > 0 {
		sigAvgTime = esUseTime / readSum
	}

	report := &share_message.VdoHallPointsReport{
		CreateTime:      easygo.NewInt64(p.StartTime),
		ReadSum:         easygo.NewInt64(readSum),
		ClickSum:        easygo.NewInt64(clickSum),
		EsUserCount:     easygo.NewInt64(esUserCount),
		FollowHallCount: easygo.NewInt64(followHallCount),
		FollowUserCount: easygo.NewInt64(followUserCount),
		AvgTime:         easygo.NewInt64(avgTime),
		SigAvgTime:      easygo.NewInt64(sigAvgTime),
		OpenUserCount:   easygo.NewInt64(openUserCount),
		VdoLsCount:      easygo.NewInt64(vdoLsCount),
		MyFollowCount:   easygo.NewInt64(myFollowCount),
		PlayLogCount:    easygo.NewInt64(playLogCount),
	}

	var reportTable string
	switch types {
	case "day":
		reportTable = for_game.TABLE_ESPORTS_VDOHALL_POINTS_REPORT_DAY
	case "week":
		reportTable = for_game.TABLE_ESPORTS_VDOHALL_POINTS_REPORT_WEEK
	case "month":
		reportTable = for_game.TABLE_ESPORTS_VDOHALL_POINTS_REPORT_MONTH
	}
	if readSum+clickSum+esUserCount+followHallCount+followUserCount+avgTime+sigAvgTime+openUserCount+vdoLsCount+myFollowCount+playLogCount > 0 {
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report}, true)
	}
}

//申请放映厅埋点报表
func (p *PointReportComm) MakeApplyVdoHallPointsReport(types string) {
	ocM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_5, "MenuId": for_game.ESPORTMENU_LIVE, "DataType": 3, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	openCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, ocM)

	acM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_5, "MenuId": for_game.ESPORTMENU_LIVE, "DataType": 3, "ButtonId": 1, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	applyCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, acM)

	abcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_5, "MenuId": for_game.ESPORTMENU_LIVE, "DataType": 3, "ButtonId": 2, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	applyBackCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, abcM)

	report := &share_message.ApplyVdoHallPointsReport{
		CreateTime:     easygo.NewInt64(p.StartTime),
		OpenCount:      easygo.NewInt64(openCount),
		ApplyCount:     easygo.NewInt64(applyCount),
		ApplyBackCount: easygo.NewInt64(applyBackCount),
	}
	var reportTable string
	switch types {
	case "day":
		reportTable = for_game.TABLE_ESPORTS_APPLYVDOHALL_POINTS_REPORT_DAY
	case "week":
		reportTable = for_game.TABLE_ESPORTS_APPLYVDOHALL_POINTS_REPORT_WEEK
	case "month":
		reportTable = for_game.TABLE_ESPORTS_APPLYVDOHALL_POINTS_REPORT_MONTH
	}
	if openCount+applyCount+applyBackCount > 0 {
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report}, true)
	}
}

//赛事列表埋点报表
func (p *PointReportComm) MakeMatchLsPointsReport(types string) {
	eutM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "MenuId": for_game.ESPORTMENU_GAME, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$Duration"}}},
	}
	esUseTime := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_DURATION_LOG, eutM)

	mlcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "MenuId": for_game.ESPORTMENU_GAME, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	matchLsCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, mlcM)

	mlM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_5, "MenuId": for_game.ESPORTMENU_GAME, "ButtonId": 0, "CreateTime": p.TimeBson}}, //ButtonId=0表示非按钮
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
	}
	matchLsClick := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, mlM)

	tcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_4, "MenuId": for_game.ESPORTMENU_GAME, "ExTabId": 1, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	todayCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, tcM)

	mcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_4, "MenuId": for_game.ESPORTMENU_GAME, "ExTabId": 2, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	matchCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, mcM)

	rcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_4, "MenuId": for_game.ESPORTMENU_GAME, "ExTabId": 3, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	rollCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, rcM)

	ocM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_4, "MenuId": for_game.ESPORTMENU_GAME, "ExTabId": 3, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	overCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, ocM)

	avgTime := int64(0)
	if p.FootEsBtnCount > 0 {
		avgTime = esUseTime / p.FootEsBtnCount
	}
	sigAvgTime := int64(0)
	if matchLsClick > 0 {
		sigAvgTime = esUseTime / matchLsClick
	}
	report := &share_message.MatchLsPointsReport{
		CreateTime:   easygo.NewInt64(p.StartTime),
		MatchLsClick: easygo.NewInt64(matchLsClick),
		AvgTime:      easygo.NewInt64(avgTime),
		SigAvgTime:   easygo.NewInt64(sigAvgTime),
		TodayCount:   easygo.NewInt64(todayCount),
		MatchCount:   easygo.NewInt64(matchCount),
		RollCount:    easygo.NewInt64(rollCount),
		OverCount:    easygo.NewInt64(overCount),
		MatchLsCount: easygo.NewInt64(matchLsCount),
	}
	var reportTable string
	switch types {
	case "day":
		reportTable = for_game.TABLE_ESPORTS_MATCHLS_POINTS_REPORT_DAY
	case "week":
		reportTable = for_game.TABLE_ESPORTS_MATCHLS_POINTS_REPORT_WEEK
	case "month":
		reportTable = for_game.TABLE_ESPORTS_MATCHLS_POINTS_REPORT_MONTH
	}
	if matchLsClick+avgTime+sigAvgTime+todayCount+matchCount+rollCount+overCount+matchLsCount > 0 {
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report}, true)
	}
}

//赛事详情埋点报表
func (p *PointReportComm) MakeMatchDilPointsReport(types string) {
	eutM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "MenuId": for_game.ESPORTMENU_GAME, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$Duration"}}},
	}
	esUseTime := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_DURATION_LOG, eutM)

	mlM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_5, "MenuId": for_game.ESPORTMENU_GAME, "ButtonId": 0, "CreateTime": p.TimeBson}}, //ButtonId=0表示非按钮
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
	}
	matchLsClick := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, mlM)

	lcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_6, "MenuId": for_game.ESPORTMENU_GAME, "ExId": 1, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	lineupCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, lcM)

	dcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_6, "MenuId": for_game.ESPORTMENU_GAME, "ExId": 2, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	dataCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, dcM)

	dtM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_6, "MenuId": for_game.ESPORTMENU_GAME, "ExId": 2, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$Duration"}}},
	}
	dataTime := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_DURATION_LOG, dtM)

	datacM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_6, "MenuId": for_game.ESPORTMENU_GAME, "ExId": 2, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
	}
	dataClick := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, datacM)

	gcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_6, "MenuId": for_game.ESPORTMENU_GAME, "ExId": 3, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	guessCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, gcM)

	gtM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_6, "MenuId": for_game.ESPORTMENU_GAME, "ExId": 3, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$Duration"}}},
	}
	guessTime := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_DURATION_LOG, gtM)

	guessM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_6, "MenuId": for_game.ESPORTMENU_GAME, "ExId": 3, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$ActCount"}}},
	}
	guessClick := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, guessM)

	mdcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_5, "MenuId": for_game.ESPORTMENU_GAME, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	matchDilCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, mdcM)

	avgTime := int64(0)
	if p.FootEsBtnCount > 0 {
		avgTime = esUseTime / p.FootEsBtnCount
	}
	sigAvgTime := int64(0)
	if matchLsClick > 0 {
		sigAvgTime = esUseTime / matchLsClick
	}

	dataAvgTime := int64(0)
	if p.FootEsBtnCount > 0 {
		dataAvgTime = dataTime / p.FootEsBtnCount
	}
	dataSigAvgTime := int64(0)
	if dataClick > 0 {
		dataSigAvgTime = dataTime / dataClick
	}

	guAvgTime := int64(0)
	if p.FootEsBtnCount > 0 {
		dataAvgTime = guessTime / p.FootEsBtnCount
	}
	guSigAvgTime := int64(0)
	if guessClick > 0 {
		dataSigAvgTime = guessTime / guessClick
	}

	report := &share_message.MatchDilPointsReport{
		CreateTime:     easygo.NewInt64(p.StartTime),
		MatchClick:     easygo.NewInt64(matchLsClick),
		EsClickCount:   easygo.NewInt64(p.FootEsBtnCount),
		AvgTime:        easygo.NewInt64(avgTime),
		SigAvgTime:     easygo.NewInt64(sigAvgTime),
		LineupCount:    easygo.NewInt64(lineupCount),
		DataCount:      easygo.NewInt64(dataCount),
		DataAvgTime:    easygo.NewInt64(dataAvgTime),
		DataSigAvgTime: easygo.NewInt64(dataSigAvgTime),
		GuessCount:     easygo.NewInt64(guessCount),
		GuessAvgTime:   easygo.NewInt64(guAvgTime),
		GuessSigTime:   easygo.NewInt64(guSigAvgTime),
		MatchDilCount:  easygo.NewInt64(matchDilCount),
	}
	var reportTable string
	switch types {
	case "day":
		reportTable = for_game.TABLE_ESPORTS_MATCHDIL_POINTS_REPORT_DAY
	case "week":
		reportTable = for_game.TABLE_ESPORTS_MATCHDIL_POINTS_REPORT_WEEK
	case "month":
		reportTable = for_game.TABLE_ESPORTS_MATCHDIL_POINTS_REPORT_MONTH
	}
	if matchLsClick+avgTime+sigAvgTime+lineupCount+dataCount+dataAvgTime+dataSigAvgTime+guessCount+guAvgTime+guSigAvgTime+matchDilCount > 0 {
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report}, true)
	}
}

//竞猜页埋点报表
func (p *PointReportComm) MakeGuessPointsReport(types string) {
	bcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_6, "MenuId": for_game.ESPORTMENU_GAME, "ExId": 3, "ButtonId": 3, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	betCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, bcM)

	bocM := []bson.M{
		{"$match": bson.M{"CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayInfo.PlayId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	betOkCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, bocM)

	report := &share_message.GuessPointsReport{
		CreateTime:   easygo.NewInt64(p.StartTime),
		EsClickCount: easygo.NewInt64(p.FootEsBtnCount),
		BetCount:     easygo.NewInt64(betCount),
		BetOkCount:   easygo.NewInt64(betOkCount),
	}
	var reportTable string
	switch types {
	case "day":
		reportTable = for_game.TABLE_ESPORTS_GUESS_POINTS_REPORT_DAY
	case "week":
		reportTable = for_game.TABLE_ESPORTS_GUESS_POINTS_REPORT_WEEK
	case "month":
		reportTable = for_game.TABLE_ESPORTS_GUESS_POINTS_REPORT_MONTH
	}
	if betCount+betOkCount > 0 {
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report}, true)
	}
}

//消息页埋点报表
func (p *PointReportComm) MakeMsgPointsReport(types string) {
	mcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "MenuId": for_game.ESPORTMENU_SYS_MSG, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	msgCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, mcM)

	smcM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_3, "MenuId": for_game.ESPORTMENU_SYS_MSG, "LabelId": 1, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	sysMsgCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, smcM)

	usM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_3, "MenuId": for_game.ESPORTMENU_SYS_MSG, "LabelId": 2, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	unSettle := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, usM)

	sM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_3, "MenuId": for_game.ESPORTMENU_SYS_MSG, "LabelId": 3, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	settle := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, sM)

	report := &share_message.MsgPointsReport{
		CreateTime:   easygo.NewInt64(p.StartTime),
		EsClickCount: easygo.NewInt64(p.FootEsBtnCount),
		MsgCount:     easygo.NewInt64(msgCount),
		SysMsgCount:  easygo.NewInt64(sysMsgCount),
		UnSettle:     easygo.NewInt64(unSettle),
		Settle:       easygo.NewInt64(settle),
	}
	var reportTable string
	switch types {
	case "day":
		reportTable = for_game.TABLE_ESPORTS_MSG_POINTS_REPORT_DAY
	case "week":
		reportTable = for_game.TABLE_ESPORTS_MSG_POINTS_REPORT_WEEK
	case "month":
		reportTable = for_game.TABLE_ESPORTS_MSG_POINTS_REPORT_MONTH
	}
	if msgCount+sysMsgCount+unSettle+settle > 0 {
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report}, true)
	}
}

//电竞币埋点报表
func (p *PointReportComm) MakeEsportCoinPointsReport(types string) {
	ecM := []bson.M{
		{"$match": bson.M{"PageType": for_game.ESPORT_BPS_PAGE_TYPE_2, "MenuId": for_game.ESPORTMENU_SHOP, "ButtonId": 0, "CreateTime": p.TimeBson}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	esCoinCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_BPS_CLICK_LOG, ecM)

	eoM := []bson.M{
		{"$match": bson.M{"SourceType": for_game.ESPORTCOIN_TYPE_EXCHANGE_IN, "CreateTime": bson.M{"$gte": p.StartTime * 1000, "$lt": p.EndTime * 1000}}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	exchangeOkCount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTCHANGELOG, eoM)

	report := &share_message.EsportCoinPointsReport{
		CreateTime:      easygo.NewInt64(p.StartTime),
		EsClickCount:    easygo.NewInt64(p.FootEsBtnCount),
		EsCoinCount:     easygo.NewInt64(esCoinCount),
		ExchangeOkCount: easygo.NewInt64(exchangeOkCount),
	}
	var reportTable string
	switch types {
	case "day":
		reportTable = for_game.TABLE_ESPORTS_COIN_POINTS_REPORT_DAY
	case "week":
		reportTable = for_game.TABLE_ESPORTS_COIN_POINTS_REPORT_WEEK
	case "month":
		reportTable = for_game.TABLE_ESPORTS_COIN_POINTS_REPORT_MONTH
	}
	if esCoinCount+exchangeOkCount > 0 {
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, reportTable, bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report}, true)
	}
}

//==========================================匹配埋点报表更新前一天的数据
func UpdateVCBuryingPointReport(id int64) {
	rmg := for_game.GetRedisVCBuryingPointReportObj(id)
	report := rmg.GetRedisVCBuryingPointReport()
	report.MainEnterPeopleNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_MAIN_BG_ENTER, id))
	report.MainSayHiNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_MAIN_DJ_SAYHI, id))
	report.MainZanNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_MAIN_DJ_XIHUAN, id))
	report.MainHeadNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_MAIN_DJ_GRTX, id))
	report.MainRecordNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_MAIN_DJ_LZMP, id))
	report.MainFinishNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_MAIN_DJ_FINISH, id))
	report.MainHeadCardNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_MAIN_DJ_MPTX, id))
	report.MainHandShake = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_MAIN_DJ_YYY, id))
	report.LZMPEnterPeopleNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_LZMP_BG_ENTER, id))
	report.LZMPBackNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_LZMP_DJ_EXIT, id))
	report.LZMPDuBai = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_LZMP_DJ_DUBAI, id))
	report.LZMPdypy = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_LZMP_DJ_DYPY, id))
	report.LZMPcyc = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_LZMP_DJ_CYC, id))
	report.LZMPly = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_LZMP_DJ_LY, id))
	report.LZMPssgd = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_LZMP_DJ_SSGD, id))
	report.SXHWEnterPeopleNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_SXHW_BG_ENTER, id))
	report.SXHWxhw = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_SXHW_DJ_XHW, id))
	report.SXHWwxh = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_SXHW_DJ_WXH, id))
	report.SXHWxhwHuiFu = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_SXHW_XHW_DJ_HF, id))
	report.SXHWxhwSayHi = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_SXHW_XHW_DJ_SAYHI, id))
	report.SXHWxhwZan = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_SXHW_XHW_DJ_XIHUAN, id))
	report.SSGDEnterPeopleNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_SSGD_BG_ENTER, id))
	report.SSGDsc = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_SSGD_DJ_SC, id))
	report.SSGDscBackNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_SSGD_SC_DJ_FH, id))
	report.SSGDscTJNum = easygo.NewInt64(FindBuryingPointLogEventCount(for_game.BP_VC_SSGD_SC_DJ_TJ, id))
	rmg.UpdateRedisVCBuryingPointReport(report)
	rmg.SaveToMongo()
}

//计算匹配埋点人数
func FindBuryingPointLogEventCount(eventId int32, ctime int64) int64 {
	stime := easygo.Get0ClockMillTimestamp(ctime)
	etime := easygo.Get24ClockMillTimestamp(ctime)
	findBson := bson.M{"Time": bson.M{"$gte": stime, "$lt": etime}, "EventType": eventId}
	m := []bson.M{
		{"$match": findBson},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	count := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_BURYING_POINT_LOG, m)
	return count
}
