package backstage

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"strconv"

	"github.com/akqp2019/mgo/bson"
)

//=========================================================================================================推送消息===>
//推送消息给后台客户端客服
func PushMsgToWeb(pushMsg *share_message.IMmessage) {
	ep := BrowerEpMp.LoadEndpoint(pushMsg.GetWaiterId())
	if ep != nil {
		result := ep.RpcPushIMmessage(pushMsg)
		if result.Id64 != nil && result.GetId64() != 0 {
			//更新消息的已读状态
			UpdateIMmessageStatus(result.GetId64(), 0)
			// log.Println("============更新消息为已读状态")
		}
	}
}

//消息发送器
func SendIMmessage(reqMsg *share_message.IMmessage) {
	// if MessageFilter(user, reqMsg) == false {
	// 	return
	// }

	uptime := util.GetMilliTime()
	if reqMsg.GetUpdateTime()+600000 < uptime {
		for_game.UpdateRedisWaiterCount(reqMsg.GetWaiterId(), 1) //激活10分钟前失活的消息连接，修改客服接待数量
	}

	if reqMsg.WaiterName == nil || reqMsg.GetWaiterName() == "" {
		waiter := GetUser(reqMsg.GetWaiterId())
		reqMsg.WaiterName = easygo.NewString(waiter.GetRealName())
	}

	content := reqMsg.Content[0]
	content.Sendtime = easygo.NewInt64(uptime)
	reqMsg.Content[0] = content
	reqMsg.UpdateTime = easygo.NewInt64(util.GetMilliTime())
	for_game.SaveMessageToDB(reqMsg) //保存消息到数据库

	switch reqMsg.Content[0].GetMtype() {
	case 1:
		//推送消息给后台客服
		PushMsgToWeb(reqMsg)
	case 2:
		//推送消息给大厅玩家
		SendToPlayer(reqMsg.GetPlayerId(), "RpcSendMessageToPlayer", reqMsg)
	default:
		panic("消息发送类型有误")
	}
}

// //留言消息发送器
// func SendLeaveMessage(user *share_message.Manager, reqMsg *share_message.IMmessage) {
// 	userId := reqMsg.GetWaiterId()
// 	if userId == 0 {

// 		//补充接收消息的客服到消息中
// 		reqMsg.WaiterId = easygo.NewInt64(userId)

// 		col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_WAITER_MESSAGE)
// 		defer closeFun()
// 		content := reqMsg.Content
// 		err := col.Update(
// 			bson.M{"_id": reqMsg.GetId()},
// 			bson.M{"$set": bson.M{"Snew": 1, "WaiterId": userId, "Content": content}}, //向数据文档中追加数据，并改变消息阅读状态为未读
// 		)
// 		easygo.PanicError(err)
// 	}

// 	//推送消息给后台客户端
// 	PushMsgToWeb(reqMsg)

// }

// //查询客服的离线留言消息
// func GetOfflineMessage(user *share_message.Manager) []*share_message.IMmessage {
// 	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_WAITER_MESSAGE)
// 	defer closeFun()

// 	var list []*share_message.IMmessage
// 	queryBson := bson.M{"WaiterId": user.GetId()}
// 	// queryBson = bson.M{"$or": []bson.M{bson.M{"WaiterId": bson.M{"$eq": nil}}, bson.M{"WaiterId": user.GetId()}}}
// 	queryBson["Status"] = bson.M{"$lt": 2}
// 	queryBson["Snew"] = 1
// 	query := col.Find(queryBson)
// 	err := query.Sort("UpdateTime").All(&list)
// 	easygo.PanicError(err)

// 	return list
// }

//查询指定客服的活跃消息数量
func GetActiveIMmessageCount(user *share_message.Manager) int {
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_WAITER_MESSAGE)
	defer closeFun()

	var list []*share_message.IMmessage
	// queryBson := bson.M{"$or": []bson.M{bson.M{"UpdateTime": bson.M{"$gte": util.GetMilliTime() - 60*10000}}, bson.M{"Snew": bson.M{"$gt": 0}}}}
	queryBson := bson.M{"UpdateTime": bson.M{"$gte": util.GetMilliTime() - 60*10000}}
	queryBson["Status"] = bson.M{"$lt": 2}
	queryBson["WaiterId"] = user.GetId()
	query := col.Find(queryBson)
	count, cerr := query.Count()
	easygo.PanicError(cerr)
	err := query.Sort("UpdateTime").All(&list)
	easygo.PanicError(err)

	return count
}

//更新消息阅读状态
func UpdateIMmessageStatus(id int64, status int32) {
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_WAITER_MESSAGE)
	defer closeFun()
	upbson := bson.M{"$set": bson.M{"Snew": status}}
	err := col.Update(bson.M{"_id": id}, upbson)
	easygo.PanicError(err)
}

//========================================================================================================其他业务====>
//绩效列表查询
func QueryWaiterPerformanceList(reqMsg *brower_backstage.ListRequest) ([]*share_message.WaiterPerformance, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WAITER_PERFORMANCE)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		i, _ := strconv.Atoi(reqMsg.GetKeyword()) //搜索查询不需要返回错误
		queryBson["WaiterId"] = i
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.BeginTimestamp != nil {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	var list []*share_message.WaiterPerformance
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//客服查询正在沟通的消息列表
func GetWaiterMsg(uid USER_ID) []*share_message.IMmessage {
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_WAITER_MESSAGE)
	defer closeFun()

	queryBson := bson.M{}
	queryBson["WaiterId"] = uid
	queryBson["Status"] = for_game.WAITER_MESSAGE_ING

	var list []*share_message.IMmessage
	query := col.Find(queryBson)
	err := query.Sort("-CreateTime").All(&list)
	easygo.PanicError(err)

	return list
}

//聊天记录查询
func QueryWaiterChatLogList(reqMsg *brower_backstage.ListRequest) ([]*share_message.IMmessage, int) {
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_WAITER_MESSAGE)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Account"] = reqMsg.GetKeyword()
		case 2:
			queryBson["NickName"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.BeginTimestamp != nil {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}
	queryBson["WaiterId"] = reqMsg.GetId()

	var list []*share_message.IMmessage
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//客服常见问题列表
func QueryWaiterFAQList(reqMsg *brower_backstage.ListRequest) ([]*share_message.WaiterFAQ, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WAITER_FAQ)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Title"] = bson.M{"$regex": reqMsg.GetKeyword()}
		case 2:
			var keys []string
			keys = append(keys, reqMsg.GetKeyword())
			queryBson["KeyWord"] = bson.M{"$in": keys}
		}
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.BeginTimestamp != nil {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	queryBson["Status"] = easygo.If(reqMsg.GetListType() != 0, reqMsg.GetListType(), bson.M{"$ne": nil})

	sort := reqMsg.GetSort()
	if sort == "" {
		sort = "-Clicks"
	}

	var list []*share_message.WaiterFAQ
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//客服常见问题修改
func EditWaiterFAQ(reqMsg *share_message.WaiterFAQ) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WAITER_FAQ)
	defer closeFun()
	var err error
	if len(reqMsg.GetKeyWord()) == 0 {
		_, err = col.Upsert(bson.M{"_id": reqMsg.GetId()}, reqMsg)
	} else {
		_, err = col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	}
	easygo.PanicError(err)
}

//删除客服常见问题
func DelWaiterFAQ(ids []int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WAITER_FAQ)
	defer closeFun()

	_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	easygo.PanicError(err)
}

//客服常用语列表
func QueryWaiterFastReply(reqMsg *brower_backstage.ListRequest) ([]*share_message.WaiterFastReply, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WAITER_FASTREPLY)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		queryBson["Content"] = bson.M{"$regex": reqMsg.GetKeyword()}
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.BeginTimestamp != nil {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	queryBson["Status"] = easygo.If(reqMsg.GetListType() != 0, reqMsg.GetListType(), bson.M{"$ne": nil})

	var list []*share_message.WaiterFastReply
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//客服常用语列表不分页
func QueryWaiterFastReplyNopage() ([]*share_message.WaiterFastReply, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WAITER_FASTREPLY)
	defer closeFun()

	queryBson := bson.M{}
	queryBson["Status"] = 1

	var list []*share_message.WaiterFastReply
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("_id").All(&list)
	easygo.PanicError(errc)

	return list, count
}

//客服常用语修改
func EditWaiterFastReply(reqMsg *share_message.WaiterFastReply) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WAITER_FASTREPLY)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//删除客服常用语
func DelWaiterFastReply(ids []int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WAITER_FASTREPLY)
	defer closeFun()

	_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	easygo.PanicError(err)
}
