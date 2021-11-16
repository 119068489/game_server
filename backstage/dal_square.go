package backstage

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"math"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//启服加载定时发送任务
func TimedSendDynamic() {
	logs.Info("加载定时发送任务")
	//list := GetDynamicByStatus(for_game.DYNAMIC_STATUE_UNPUBLISHED) //查询未发布的动态
	list := GetDynamicByStatusFromSnap(for_game.DYNAMIC_STATUE_UNPUBLISHED) //查询临时表中未发布的动态
	for _, item := range list {
		AddSendDynamic(item)
	}
}

//添加发送任务定时器
func AddSendDynamic(item *share_message.DynamicData) {
	if SendDynamicTimeMgr.GetTimerById(item.GetLogId()) != nil {
		SendDynamicTimeMgr.DelTimerList(item.GetLogId())
	}
	triggerTime := time.Duration(item.GetSendTime()-time.Now().Unix()) * time.Second
	if triggerTime >= 0 {
		//  插入临时库
		InsertDynamicToSnapDB(item)
		timer := easygo.AfterFunc(triggerTime, func() {
			SendDynamic(item)
		})
		SendDynamicTimeMgr.AddTimerList(item.GetLogId(), timer)
	}
	//} else {
	//	UpdateDynamicStatus(item.GetLogId(), for_game.DYNAMIC_STATUE_EXPIRED) // 过期动态
	//}
}

//删除动态发小助手通知
func SendSystemNotice(pid PLAYER_ID, title, content string) {
	noticMsg := &share_message.SystemNotice{
		Id:      easygo.NewInt64(pid),
		Title:   easygo.NewString(title),
		Content: easygo.NewString(content),
	}
	SendToPlayer(pid, "RpcSendSystemNoticeToPlayer", noticMsg)
}

//定时发布动态
func SendDynamic(item *share_message.DynamicData) {
	if item.GetStatue() != for_game.DYNAMIC_STATUE_UNPUBLISHED {
		return
	}
	// 从原来的库中删除记录
	DelDynamicFromSnapDB(item.GetLogId())
	SendDynamicTimeMgr.DelTimerList(item.GetLogId()) // 删除定时任务,临时表的id
	item.Statue = easygo.NewInt32(for_game.DYNAMIC_STATUE_COMMON)
	item.LogId = easygo.NewInt64(for_game.NextId(for_game.TABLE_SQUARE_DYNAMIC))
	// 把动态存在正式动态表中
	insertErr := for_game.InsertDynamicToDB(item)
	easygo.PanicError(insertErr)
	for_game.UpdateRedisSquareDynamic(item)
	//UpdateDynamicStatus(item.GetLogId(), for_game.DYNAMIC_STATUE_COMMON)
	who := for_game.GetRedisPlayerBase(item.GetPlayerId())
	who.AddRedisPlayerDynamicList(item.GetLogId())
	easygo.Spawn(SetHotDynamic, item.GetLogId(), item.GetIsHot())

	if item.GetIsBsTop() && item.GetStatue() == for_game.DYNAMIC_STATUE_COMMON {
		msg := &server_server.TopRequest{
			LogId: item.LogId,
		}
		// BroadCastToAllHall("RpcBackstageTop", msg)
		ChooseOneHall(0, "RpcBackstageTop", msg) //通知大厅置顶社交动态
	}

}

func QueryDynamic(reqMsg *brower_backstage.DynamicListRequest) (dataList []*share_message.DynamicData, countPage int32) {
	// 判断是否是未发布的动态
	var col *mgo.Collection
	var closeFun func()
	if reqMsg.GetStatus() == for_game.DYNAMIC_STATUE_UNPUBLISHED {
		col, closeFun = easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC_SNAP)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC)
	}
	defer closeFun()
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)
	queryBson := bson.M{} //"Statue": 0

	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.BeginTimestamp != nil && reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		queryBson["SendTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp() / 1000, "$lte": reqMsg.GetEndTimestamp() / 1000}
	}

	if reqMsg.GetType() != 0 && reqMsg.Type != nil && reqMsg.GetKeyword() != "" && reqMsg.Keyword != nil {
		switch reqMsg.GetType() {
		case 1:
			pl := QueryPlayerbyAccount(reqMsg.GetKeyword())
			if pl != nil {
				queryBson["PlayerId"] = pl.GetPlayerId()
			}
		case 8:
			queryId := easygo.StringToIntnoErr(reqMsg.GetKeyword())
			queryBson["_id"] = queryId
		case 9:
			topBson := bson.M{"Name": bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}}
			topicidList, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, topBson, 0, 0)
			var topicIds []int64
			for _, tidli := range topicidList {
				topicIds = append(topicIds, tidli.(bson.M)["_id"].(int64))
			}
			if len(topicIds) > 0 {
				queryBson["TopicId"] = bson.M{"$in": topicIds}
			} else {
				return nil, 0
			}
		}
	}

	if reqMsg.Check != nil && reqMsg.GetCheck() != 1000 {
		if reqMsg.GetCheck() == 0 {
			queryBson["Check"] = nil
		} else {
			queryBson["Check"] = reqMsg.GetCheck()
		}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 1000 {
		queryBson["Statue"] = reqMsg.GetStatus()
	}

	if reqMsg.IsTop != nil && reqMsg.GetIsTop() != 0 {
		queryBson["IsTop"] = easygo.If(reqMsg.GetIsTop() == 1, true, bson.M{"$ne": true})
	}

	if reqMsg.IsBsTop != nil && reqMsg.GetIsBsTop() != 0 {
		queryBson["IsBsTop"] = easygo.If(reqMsg.GetIsBsTop() == 1, true, bson.M{"$ne": true})
	}

	if reqMsg.IsShield != nil && reqMsg.GetIsShield() != 0 {
		queryBson["IsShield"] = easygo.If(reqMsg.GetIsShield() == 1, true, bson.M{"$ne": true})
	}

	if reqMsg.ListType != nil {
		switch reqMsg.GetListType() {
		case 1: //普通用户查询
			queryBson["SenderType"] = bson.M{"$lte": 1} //小于等于1普通用户
		case 2: //运营用户查询
			queryBson["SenderType"] = bson.M{"$gt": 1} //大于1运营用户
		}
	}

	if reqMsg.ReportCountMin != nil && reqMsg.ReportCountMax != nil {
		queryBson["ReportCount"] = bson.M{"$gte": reqMsg.GetReportCountMin(), "$lte": reqMsg.GetReportCountMax()}
	}

	if reqMsg.TopicType != nil && reqMsg.GetTopicType() > 0 && len(reqMsg.GetTopicId()) == 0 {
		idList, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, bson.M{"TopicTypeId": reqMsg.GetTopicType()}, 0, 0)
		var ids []int64
		for _, idli := range idList {
			ids = append(ids, idli.(bson.M)["_id"].(int64))
		}
		if len(ids) > 0 {
			queryBson["TopicId"] = bson.M{"$in": ids}
		}
	}

	if len(reqMsg.GetTopicId()) > 0 {
		queryBson["TopicId"] = bson.M{"$in": reqMsg.GetTopicId()}
	}

	var list []*share_message.DynamicData
	query := col.Find(queryBson)
	count, err1 := query.Count()
	easygo.PanicError(err1)

	err2 := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(err2)

	for _, value := range list {
		player := for_game.GetRedisPlayerBase(value.GetPlayerId())
		value.Account = easygo.NewString(player.GetAccount())
		value.NickName = easygo.NewString(player.GetNickName())
		if reqMsg.GetStatus() == for_game.DYNAMIC_STATUE_COMMON {
			value.CommentNum = easygo.NewInt64(GetDynamicCommentNum(value.GetLogId()))
			value.Zan = easygo.NewInt32(for_game.GetRedisDynamicZanNum(value.GetLogId()))
		}
	}
	return list, int32(count)
}

//更新动态状态 // 0：正常 1:后台删除 2：前端删除,3未发布，4已过期
func UpdateDynamicStatus(id int64, status int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	err := col.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"Statue": status}})
	easygo.PanicError(err)
}

//更新动态状态 //审核状态 0未处理,1已审核,2已拒绝，3自动审核
func UpdateDynamicCheck(id int64, check int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	err := col.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"Check": check}})
	easygo.PanicError(err)
}

//保存动态
func SaveDynamic(data *share_message.DynamicData) {
	col, closeFunc := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC)
	defer closeFunc()

	_, err := col.Upsert(bson.M{"_id": data.GetLogId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

//根据状态查询动态 // 0：正常 1:后台删除 2：前端删除,3未发布，4已过期
func GetDynamicByStatus(status int) []*share_message.DynamicData {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	queryBson := bson.M{"Statue": status}
	query := col.Find(queryBson)
	var list []*share_message.DynamicData
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

func QueryDynamicComment(reqMsg *brower_backstage.DynamicListRequest) ([]*share_message.CommentData, int32) {
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_COMMENT)
	col2, closeFun2 := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	defer closeFun2()

	queryBson := bson.M{}
	queryBson2 := bson.M{}

	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}
	playerId := make([]int64, 0)
	if reqMsg.GetType() != 0 && reqMsg.Type != nil && reqMsg.GetKeyword() != "" && reqMsg.Keyword != nil {
		var players []*share_message.PlayerBase

		switch reqMsg.GetType() {
		case 3, 5:
			queryBson2["NickName"] = reqMsg.GetKeyword()
		case 4, 6:
			queryBson2["Account"] = reqMsg.GetKeyword()
		}
		col2.Find(queryBson2).All(&players)
		if len(players) != 0 {
			for _, player := range players {
				playerId = append(playerId, player.GetPlayerId())
			}
		}
	}

	queryBson["LogId"] = reqMsg.GetLogId()

	switch reqMsg.GetType() {
	case 3, 4:
		queryBson["PlayerId"] = bson.M{"$in": playerId}
	case 5, 6:
		queryBson["TargetId"] = bson.M{"$in": playerId}
	case 7:
		queryBson["Content"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "i"}}
	}

	var list []*share_message.CommentData
	query := col.Find(queryBson)
	count, err1 := query.Count()
	if err1 != nil && err1 != mgo.ErrNotFound {
		easygo.PanicError(err1)
	}
	err2 := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if err2 != nil && err2 != mgo.ErrNotFound {
		easygo.PanicError(err2)
	}

	for _, value := range list {
		value.Name = easygo.NewString(for_game.GetRedisPlayerBase(value.GetPlayerId()).GetNickName())
		value.OtherName = easygo.NewString(for_game.GetRedisPlayerBase(value.GetTargetId()).GetNickName())
	}
	return list, int32(count)
}

func QueryDynamicById(id int64) *share_message.DynamicData {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	siteOne := &share_message.DynamicData{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	player := QueryPlayerbyId(siteOne.GetPlayerId())
	siteOne.HeadIcon = easygo.NewString(player.GetHeadIcon())
	siteOne.NickName = easygo.NewString(player.GetNickName())
	siteOne.Account = easygo.NewString(player.GetAccount())
	siteOne.PlayerTypes = easygo.NewInt32(player.GetTypes())
	siteOne.Sex = easygo.NewInt32(player.GetSex())
	siteOne.Zan = easygo.NewInt32(for_game.GetRedisDynamicZanNum(id) + siteOne.GetTrueZan())
	siteOne.CommentNum = easygo.NewInt64(for_game.GetRedisDynamicCommentNum(id))

	return siteOne
}

func DelCommentDatas(ids []int64, note string) {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_COMMENT)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"Statue": for_game.DYNAMIC_COMMENT_STATUE_DELETE, "Note": note}})
	easygo.PanicError(err)
}

//查询指定时间发布的动态列表
func QueryDynamicByTime(startTime, endTime TIME_64) []*share_message.DynamicData {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	queryBson := bson.M{"SendTime": bson.M{"$gte": startTime, "$lte": endTime}}

	var list []*share_message.DynamicData
	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)

	return list
}

//查询指定时间的动态评论
func QueryDynamicCommentByTime(startTime, endTime TIME_64) []*share_message.CommentData {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_COMMENT)
	defer closeFun()

	queryBson := bson.M{"CreateTime": bson.M{"$gte": startTime, "$lte": endTime}}

	var list []*share_message.CommentData
	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)

	return list
}

//查询指定时间的点赞
func QueryZanDataByTime(startTime, endTime TIME_64) []*share_message.ZanData {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_ZAN)
	defer closeFun()

	queryBson := bson.M{"CreateTime": bson.M{"$gte": startTime, "$lte": endTime}}

	var list []*share_message.ZanData
	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)

	return list
}

//查询评论说
func GetDynamicCommentNum(logId int64) int64 {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_COMMENT)
	defer closeFun()

	query := col.Find(bson.M{"LogId": logId})
	count, err := query.Count()
	easygo.PanicError(err)
	return int64(count)
}

//一键设置热门
func SetHotDynamic(logId int64, isHot bool) *base.Fail {
	if !isHot {
		return nil
	}
	dynamic := QueryDynamicById(logId)
	if dynamic == nil {
		return easygo.NewFailMsg("动态不存在")
	}

	if dynamic.GetIsTop() || dynamic.GetIsBsTop() {
		return easygo.NewFailMsg("置顶状态的动态不能设置热门")
	}

	if easygo.GetDifferenceDay(dynamic.GetSendTime(), easygo.NowTimestamp()) > 7 {
		return easygo.NewFailMsg("动态发布已超过7天,不能设置热门")
	}

	config := for_game.QuerySysParameterById(for_game.SQUAREHOT_PARAMETER)
	if config == nil {
		return easygo.NewFailMsg("系统热门参数未设置")
	}

	if config.GetZanScore() == 0 {
		return nil
	}

	score := config.GetHotScore() - dynamic.GetHostScore()
	if score <= 0 {
		return nil
	}

	pnub := math.Ceil(float64(score) / float64(config.GetZanScore()))

	players := for_game.GetRandPlayerByTypes([]int32{2, 3, 4, 5}, int32(pnub))
	sysp := PSysParameterMgr.GetSysParameter(for_game.PUSH_PARAMETER)
	for _, i := range players {
		for_game.OperateRedisDynamicZan(for_game.DYNAMIC_OPERATE, i.GetPlayerId(), logId, dynamic.GetPlayerId(), sysp)
		// 对该动态的分数加累加
		for_game.UpsetDynamicScore(logId, config.GetZanScore())
	}

	return nil
}

// 插入动态进临时表.
func InsertDynamicToSnapDB(dynamic *share_message.DynamicData) {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC_SNAP)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": dynamic.GetLogId()}, dynamic)
	easygo.PanicError(err)
}

func DelDynamicFromSnapDB(id int64) {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC_SNAP)
	defer closeFun()
	err := col.RemoveId(id)
	easygo.PanicError(err)
}

// 从临时表中读取所有定时的列表.
func GetDynamicByStatusFromSnap(status int) []*share_message.DynamicData {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC_SNAP)
	defer closeFun()
	queryBson := bson.M{"Statue": status}
	query := col.Find(queryBson)
	var list []*share_message.DynamicData
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}
func GetDynamicByIdFromSnap(logId int64, status int) *share_message.DynamicData {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC_SNAP)
	defer closeFun()
	queryBson := bson.M{"Statue": status, "_id": logId}
	query := col.Find(queryBson)
	var dynamic *share_message.DynamicData
	err := query.One(&dynamic)
	easygo.PanicError(err)

	return dynamic
}

//======================话题=====》
//查询话题分类列表查询
func QueryTopicTypeList(reqMsg *brower_backstage.ListRequest) ([]*share_message.TopicType, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_TYPE)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Name"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
		}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() > 0 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	var list []*share_message.TopicType
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-Sort").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	for i, li := range list {
		list[i].TopicCount = easygo.NewInt64(for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, bson.M{"TopicTypeId": li.GetId()}))
	}

	return list, count
}

//查询话题列表查询
func QueryTopicList(reqMsg *brower_backstage.ListRequest) ([]*share_message.Topic, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Name"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
		case 2:
			queryBson["Owner"] = reqMsg.GetKeyword()
		case 3:
			queryBson["Admin"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() > 0 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() >= 0 {
		queryBson["TopicTypeId"] = reqMsg.GetListType()
	}

	switch reqMsg.GetDownType() {
	case 1:
		queryBson["IsHot"] = true

	case 2:
		queryBson["IsRecommend"] = true
	case 3:
		queryBson["IsRecommend"] = bson.M{"$ne": true}
		queryBson["IsHot"] = bson.M{"$ne": true}
	}

	sort := "HotScore"
	sortType := -1
	if reqMsg.Sort != nil && reqMsg.GetSort() != "" {
		sort = reqMsg.GetSort()
	} else {
		sort = "HotScore"
	}

	if reqMsg.GetSrtType() == "asc" {
		sortType = 1
	} else {
		sortType = -1
	}

	logs.Debug(queryBson)
	m := []bson.M{
		{"$match": queryBson},
		{"$addFields": bson.M{"HotScore": bson.M{"$add": []string{"$ViewingNum", "$FansNum", "$ParticipationNum", "$AddViewingNum", "$AddParticipationNum", "$AddFansNum"}}}},
		{"$sort": bson.M{sort: sortType}},
		{"$project": bson.M{"HotScore": 0}},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	query := col.Pipe(m)
	var list []*share_message.Topic
	err := query.All(&list)
	easygo.PanicError(err)

	querycount := col.Find(queryBson)
	count, err := querycount.Count()
	easygo.PanicError(err)

	return list, count
}
