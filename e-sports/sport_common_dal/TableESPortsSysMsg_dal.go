package sport_common_dal

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	//   "time"
)

func AddTableESPortsSysMsg(indx int32) {
	b := true
	info := share_message.TableESPortsSysMsg{
		RecipientType: easygo.NewInt64(2),
		Title:         easygo.NewString(fmt.Sprintf("发个毛系统消息%d", indx)),
		Content:       easygo.NewString("烦人的信息"),
		Status:        easygo.NewInt32(1),
		CreateTime:    easygo.NewInt64(easygo.NowTimestamp()),
		JumpInfo: &share_message.ESPortsJumpInfo{
			//跳转类型 1外部跳转，2内部跳转,3跳转其他APP
			JumpType: easygo.NewInt32(2),
			//跳转对象Id
			JumpObjId: easygo.NewInt64(1),
			//跳转位置 1 主界面，2 柠檬团队，3 柠檬助手，4附近的人，5社交广场-主界面，6社交广场-新增关注，7社交广场-指定动态：通过填写动态ID指定，8好物-主界面，9好物-指定商品：通过填写商品ID指定,10群-指定群id,11社交广场发布页,12零钱,13话题-指定话题,14-指定的动态评论,15-话题主界面,16-硬币商城主页,17-电竞币充值页,18-指定资讯详情,19-指定视频详情,20-电竞主页
			JumpObject: easygo.NewInt32(2),
			//跳转URL
			JumpUrl: easygo.NewString(""),
			//跳转对象样式 0默认，1隐藏头部
			JumpStyle: easygo.NewInt32(1),
		},
		EffectiveTime:   easygo.NewInt64(easygo.NowTimestamp()),
		EffectiveType:   easygo.NewInt64(1),
		IsPush:          &b,
		IsMessageCenter: &b,
		Icon:            easygo.NewString("boy_5"),
	}
	CreateTableESPortsSysMsg(&info)
}

func CreateTableESPortsSysMsg(info *share_message.TableESPortsSysMsg) (int32, string) {
	id := for_game.NextId(for_game.TABLE_ESPORTS_SYS_MSG)
	col, closeFun := GetC(for_game.TABLE_ESPORTS_SYS_MSG)
	defer closeFun()

	info.Id = easygo.NewInt64(id)
	_, err := col.Upsert(bson.M{"_id": id}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}
func GetTableESportAllSysList(lasttime int64, recipientType int32) ([]*share_message.TableESPortsSysMsg, int32) {
	var list []*share_message.TableESPortsSysMsg
	col, closeFun := GetC(for_game.TABLE_ESPORTS_SYS_MSG)
	defer closeFun()
	queryBson := bson.M{}
	nowUnit := easygo.NowTimestamp() //time.Now().Unix()
	////过滤掉时间未到的
	queryBson["Status"] = int32(1)
	queryBson["IsMessageCenter"] = true
	//queryBson["$or"] = []bson.M{bson.M{"RecipientType": 0}, bson.M{"RecipientType": recipientType}}
	queryBson["$and"] = []bson.M{
		bson.M{"EffectiveTime": bson.M{"$lte": nowUnit}},
		bson.M{"EffectiveTime": bson.M{"$gt": lasttime}},
		bson.M{"$and": []bson.M{
			bson.M{"$or": []bson.M{
				bson.M{"FailureTime": 0},
				bson.M{"FailureTime": bson.M{"$gt": nowUnit}}},
			},
			bson.M{"$or": []bson.M{
				bson.M{"RecipientType": 0},
				bson.M{"RecipientType": recipientType}},
			}},
		},
	}
	//queryBson["$or"] = []bson.M{bson.M{"FailureTime": 0}, bson.M{"FailureTime": bson.M{"$gt": nowUnit}}}
	query := col.Find(queryBson)
	err := query.Sort("-EffectiveTime").All(&list)
	if err != nil {
		logs.Error(err)
		return nil, for_game.C_SYS_ERROR
	}

	return list, for_game.C_OPT_SUCCESS
}

func GetTableESportAllSysMsgCount(lasttime int64, recipientType int32) int32 {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_SYS_MSG)
	defer closeFun()
	queryBson := bson.M{}
	nowUnit := easygo.NowTimestamp() //time.Now().Unix()
	////过滤掉时间未到的
	queryBson["Status"] = int32(1)
	queryBson["IsMessageCenter"] = true
	/*
		queryBson["$or"] = []bson.M{bson.M{"RecipientType": 0}, bson.M{"RecipientType": recipientType}}

		queryBson["$and"] = []bson.M{bson.M{"EffectiveTime": bson.M{"$lte": nowUnit}}, bson.M{"EffectiveTime": bson.M{"$gt": lasttime}}} //bson.M{"$or": []bson.M{bson.M{"FailureTime": 0}, bson.M{"FailureTime": bson.M{"$lte": nowUnit}}}},

		queryBson["$or"] = []bson.M{bson.M{"FailureTime": 0}, bson.M{"FailureTime": bson.M{"$gt": nowUnit}}}
	*/
	queryBson["$and"] = []bson.M{
		bson.M{"EffectiveTime": bson.M{"$lte": nowUnit}},
		bson.M{"EffectiveTime": bson.M{"$gt": lasttime}},
		bson.M{"$and": []bson.M{
			bson.M{"$or": []bson.M{
				bson.M{"FailureTime": 0},
				bson.M{"FailureTime": bson.M{"$gt": nowUnit}}},
			},
			bson.M{"$or": []bson.M{
				bson.M{"RecipientType": 0},
				bson.M{"RecipientType": recipientType}},
			}},
		},
	}

	query := col.Find(queryBson)
	n, err := query.Count()
	logs.Info("GetTableESportAllSysMsgCount", n)
	if err != nil {
		logs.Error(err)
	}
	return int32(n)
}

//创建待开奖信息
func CreateTableESPortsGameOrderSysMsg(info *share_message.TableESPortsGameOrderSysMsg) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG)
	defer closeFun()
	info.CreateTime = easygo.NewInt64(easygo.NowTimestamp())
	info.UpdateTime = easygo.NewInt64(easygo.NowTimestamp())
	_, err := col.Upsert(bson.M{"_id": info.GetOrderId()}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}

func CreateTableESPortsGameOrderSysMsgEx(infos []*share_message.TableESPortsGameOrderSysMsg) (int32, string) {

	var saveData []interface{}
	for _, info := range infos {
		saveData = append(saveData, bson.M{"_id": info.GetOrderId()}, info)
	}
	if len(saveData) > 0 {
		for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG, saveData)
	}

	return for_game.C_OPT_SUCCESS, "創建成功"
}

func GetTableESportAllGameOrderSysMsgList(playerId int64) ([]*share_message.TableESPortsGameOrderSysMsg, int32) {
	var list []*share_message.TableESPortsGameOrderSysMsg
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG)
	defer closeFun()
	queryBson := bson.M{}

	////过滤掉时间未到的
	queryBson["PlayerId"] = playerId
	query := col.Find(queryBson)
	err := query.Sort("-CreateTime").All(&list)
	if err != nil {
		logs.Error(err)
		return nil, for_game.C_SYS_ERROR
	}
	return list, for_game.C_OPT_SUCCESS
}
func GetTableESportAllGameOrderSysMsgCount(playerId int64) int32 {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG)
	defer closeFun()
	queryBson := bson.M{}

	queryBson["PlayerId"] = playerId
	query := col.Find(queryBson)
	n, err := query.Count()
	if err != nil {
		logs.Error(err)
	}
	return int32(n)
}
func DeleteTableESPortsGameOrderSysMsg(PlayerId int64) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"PlayerId": PlayerId})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "删除成功"
}

func CreateTableESPortsGameEndOrderSysMsg(info *share_message.TableESPortsGameOrderSysMsg) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG_E)
	defer closeFun()
	info.CreateTime = easygo.NewInt64(easygo.NowTimestamp())
	info.UpdateTime = easygo.NewInt64(easygo.NowTimestamp())
	_, err := col.Upsert(bson.M{"_id": info.GetOrderId()}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}

func CreateTableESPortsGameEndOrderSysMsgEx(infos []*share_message.TableESPortsGameOrderSysMsg) (int32, string) {

	var saveData []interface{}
	for _, info := range infos {
		saveData = append(saveData, bson.M{"_id": info.GetOrderId()}, info)
	}
	if len(saveData) > 0 {
		for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG_E, saveData)
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}

func GetTableESportAllGameOrderEndSysMsgList(playerId int64) ([]*share_message.TableESPortsGameOrderSysMsg, int32) {
	var list []*share_message.TableESPortsGameOrderSysMsg
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG_E)
	defer closeFun()
	queryBson := bson.M{}

	////过滤掉时间未开奖的
	//queryBson["BetResult"] = bson.M{"$ne": for_game.GAME_GUESS_BET_RESULT_1}
	queryBson["PlayerId"] = playerId
	query := col.Find(queryBson)
	err := query.Sort("-CreateTime").All(&list)
	if err != nil {
		logs.Error(err)
		return nil, for_game.C_SYS_ERROR
	}
	return list, for_game.C_OPT_SUCCESS
}

func GetTableESportAllGameOrderEndSysMsgCount(playerId int64) int32 {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG_E)
	defer closeFun()
	queryBson := bson.M{}

	////过滤掉时间未开奖的
	queryBson["BetResult"] = bson.M{"$ne": for_game.GAME_GUESS_BET_RESULT_1}
	queryBson["PlayerId"] = playerId
	query := col.Find(queryBson)
	n, err := query.Count()
	if err != nil {
		logs.Error(err)
	}
	return int32(n)
}

func ClaerTableESPortsGameOrderEndSysMsg(PlayerId int64) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG_E)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"PlayerId": PlayerId, "BetResult": bson.M{"$ne": for_game.GAME_GUESS_BET_RESULT_1}})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "清空成功"
}
func DeleteTableESPortsGameOrderEndSysMsgByOrderId(playerId, orderId int64) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG_E)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"PlayerId": playerId, "OrderId": orderId})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "删除成功"
}

//修改订单
func UpdateTableESPortsGameOrderEndSysMsg(orderId int64, betResult string, resultAmount int64) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG_E)
	defer closeFun()
	updatedata := bson.M{}

	updatedata["BetResult"] = betResult
	updatedata["ResultAmount"] = resultAmount
	updatedata["UpdateTime"] = easygo.NowTimestamp()

	_, err := col.Upsert(bson.M{"_id": orderId}, bson.M{"$set": updatedata})
	if err != nil {
		logs.Error(err)
	}
	if err == nil {
		return for_game.C_OPT_SUCCESS, "修改成功"
	}
	return for_game.C_SYS_ERROR, "系统错误"
}

//获取是否有系统消息
func GetSysMsgCount(lasttime, playerId int64, recipientType int32) int32 {
	n := GetTableESportAllGameOrderEndSysMsgCount(playerId)
	if n > 0 {
		return n
	}
	n = GetTableESportAllGameOrderSysMsgCount(playerId)
	if n > 0 {
		return n
	}
	n = GetTableESportAllSysMsgCount(lasttime, recipientType)
	if n > 0 {
		return n
	}
	return 0
}

func AddTableESPortsRoomChatMsg(info *share_message.TableESPortsLiveRoomMsgLog, id int64) (int32, string) {

	col, closeFun := easygo.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ESPORTS_ROOM_CHAT_MSG_LOG)
	defer closeFun()

	info.Id = easygo.NewInt64(id)
	_, err := col.Upsert(bson.M{"_id": id}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}
