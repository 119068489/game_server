package backstage

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"strconv"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

//查询现金源类型
func QuerySouceTypeList(reqMsg *brower_backstage.SourceTypeRequest) []*share_message.SourceType {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SOURCETYPE)
	defer closeFun()

	var list []*share_message.SourceType
	queryBson := bson.M{}
	if reqMsg.Types != nil && reqMsg.GetTypes() != 0 {
		queryBson["Type"] = reqMsg.GetTypes()
	}
	if reqMsg.Channel != nil && reqMsg.GetChannel() != 0 {
		queryBson["Channel"] = reqMsg.GetChannel()
	}

	q := col.Find(queryBson)
	err := q.All(&list)
	easygo.PanicError(err)

	return list
}

//修改通用支付额度设置
func EditGeneralQuota(reqMsg *share_message.GeneralQuota) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_GENERAL_QUOTA)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//查询支付方式
func QueryPayType(reqMsg *brower_backstage.QueryDataById) []*share_message.PayType {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYTYPE)
	defer closeFun()

	var list []*share_message.PayType
	queryBson := bson.M{}
	if reqMsg.Id32 != nil && reqMsg.GetId32() != 0 {
		queryBson["Types"] = reqMsg.GetId32()
	}

	q := col.Find(queryBson)
	err := q.Sort("Sort").All(&list)
	easygo.PanicError(err)

	return list
}

//修改支付方式
func EditPayType(reqMsg *share_message.PayType) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYTYPE)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//id查询支付方式
func QueryPayTypeById(id int32) *share_message.PayType {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYTYPE)
	defer closeFun()

	siteOne := &share_message.PayType{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//查询支付场景
func QueryPayScene() []*share_message.PayScene {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYSCENE)
	defer closeFun()

	var list []*share_message.PayScene
	queryBson := bson.M{}

	q := col.Find(queryBson)
	err := q.All(&list)
	easygo.PanicError(err)

	return list
}

//修改支付场景
func EditPayScene(reqMsg *share_message.PayScene) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYSCENE)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//id查询支付场景
func QueryPaySceneById(id int32) *share_message.PayScene {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYSCENE)
	defer closeFun()

	siteOne := &share_message.PayScene{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//查询支付设定列表
func QueryPaymentSettingtList(reqMsg *brower_backstage.ListRequest) ([]*share_message.PaymentSetting, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYMENTSETTING)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		queryBson["Types"] = reqMsg.GetListType()
	}
	switch reqMsg.GetType() {
	case 1:
		queryBson["Name"] = reqMsg.GetKeyword()
	}

	var list []*share_message.PaymentSetting
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//id查询支付设定
func QueryPaymentSettingById(id int32) *share_message.PaymentSetting {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYMENTSETTING)
	defer closeFun()

	siteOne := &share_message.PaymentSetting{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//修改支付设定
func EditPaymentSetting(reqMsg *share_message.PaymentSetting) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYMENTSETTING)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//查询支付平台
func QueryPaymentPlatform(reqMsg *brower_backstage.PlatformChannelRequest) ([]*share_message.PaymentPlatform, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYMENTPLATFORM)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	switch reqMsg.GetTypes() {
	case 1:
		queryBson["Name"] = reqMsg.GetKeyword()
	}

	var list []*share_message.PaymentPlatform
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改支付平台
func EditPaymentPlatform(reqMsg *share_message.PaymentPlatform) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYMENTPLATFORM)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//名称查询支付平台
func QuerPaymentPlatformByName(name string) *share_message.PaymentPlatform {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYMENTPLATFORM)
	defer closeFun()

	siteOne := &share_message.PaymentPlatform{}
	err := col.Find(bson.M{"Name": name}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//id查询支付平台
func QuerPaymentPlatformById(id int32) *share_message.PaymentPlatform {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PAYMENTPLATFORM)
	defer closeFun()

	siteOne := &share_message.PaymentPlatform{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//查询支付平台通道列表
func QueryPlatformChannelList(reqMsg *brower_backstage.PlatformChannelRequest) ([]*share_message.PlatformChannel, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLATFORM_CHANNEL)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		queryBson["Types"] = reqMsg.GetListType()
	}
	if reqMsg.PayType != nil && reqMsg.GetPayType() != 0 {
		queryBson["PayTypeId"] = reqMsg.GetPayType()
	}
	if reqMsg.Status != nil && reqMsg.GetStatus() != 1000 {
		queryBson["Status"] = reqMsg.GetStatus()
	}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetTypes() {
		case 1:
			queryBson["Name"] = reqMsg.GetKeyword()
		case 2:
			i, _ := strconv.Atoi(reqMsg.GetKeyword()) //搜索查询不需要返回错误
			queryBson["Weights"] = i
		case 3:
			platform := QuerPaymentPlatformByName(reqMsg.GetKeyword())
			if platform != nil {
				queryBson["PlatformId"] = platform.GetId()
			} else {
				queryBson["PlatformId"] = 0
			}
		}
	}

	var list []*share_message.PlatformChannel
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改支付平台通道
func EditPlatformChannel(reqMsg *share_message.PlatformChannel) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLATFORM_CHANNEL)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//批量修改通道状态 1开启，2关闭
func BatchClosePlatformChannel(ids []int32, status int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLATFORM_CHANNEL)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"Status": status}})
	easygo.PanicError(err)
}

//多种id查询支付平台通道  types=1平台Id查询，2场景id查询，3支付类型id查询
func QueryPlatformChannelByPid(pid int32, types int32) *share_message.PlatformChannel {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLATFORM_CHANNEL)
	defer closeFun()

	siteOne := &share_message.PlatformChannel{}
	querybson := bson.M{}
	switch types {
	case 1:
		querybson["PlatformId"] = pid
	case 2:
		querybson["PaySceneId"] = pid
	case 3:
		querybson["PayTypeId"] = pid
	}
	err := col.Find(querybson).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//订单列表查询
func QueryOrderList(reqMsg *brower_backstage.QueryOrderRequest) ([]*share_message.Order, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ORDER)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{"OrderType": bson.M{"$ne": 1}}
	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	//不查关键词
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetTypes() {
		case 1: //订单号
			queryBson["_id"] = reqMsg.GetKeyword()
		case 2: //柠檬号
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			queryBson["PlayerId"] = player.GetPlayerId()
		case 3:
			queryBson["Operator"] = reqMsg.GetKeyword()
		default:
			easygo.NewFailMsg("查询条件有误")
		}
	}

	if reqMsg.PayStatus != nil && reqMsg.GetPayStatus() != 1000 {
		queryBson["PayStatus"] = reqMsg.GetPayStatus()

	}

	if reqMsg.SourceType != nil && reqMsg.GetSourceType() > 0 {
		// types := reqMsg.GetSourceType()
		// queryBson["SourceType"] = bson.M{"$in": types}
		queryBson["SourceType"] = reqMsg.GetSourceType()

	}

	if reqMsg.PayType != nil && reqMsg.GetPayType() != 0 {
		queryBson["PayChannel"] = reqMsg.GetPayType()
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 1000 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.ChangeType != nil && reqMsg.GetChangeType() != 0 {
		queryBson["ChangeType"] = reqMsg.GetChangeType()
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.Order
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//查询玩家指定时间指定状态的订单  status=-1 查询全部订单
func GetOderByPlayerIdAndTime(id PLAYER_ID, startTime int64, endTime int64, status int32) []*share_message.Order {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ORDER)
	defer closeFun()

	var list []*share_message.Order
	queryBson := bson.M{"PlayerId": id, "CreateTime": bson.M{"$gte": startTime, "$lt": endTime}}

	if status > -1 {
		queryBson["Status"] = status
	}

	q := col.Find(queryBson)
	err := q.All(&list)
	easygo.PanicError(err)

	return list
}

//完成订单
func FinishOrder(orderid string, operator string, note string) *base.Fail {
	order := for_game.GetRedisOrderObj(orderid)
	if order == nil {
		return easygo.NewFailMsg("订单不存在")
	}

	order.SetNote(note)
	// order.Status = easygo.NewInt32(1)
	// order.OverTime = easygo.NewInt64(for_game.GetMillSecond())
	req := &server_server.Recharge{
		PlayerId:     easygo.NewInt64(order.GetPlayerId()),
		RechargeGold: easygo.NewInt64(order.GetChangeGold()), //此处的订单金额，是扣税后的金额
		OrderId:      easygo.NewString(order.GetOrderId()),
		SourceType:   easygo.NewInt32(order.GetSourceType()),
	}
	//通知大厅完成订单
	SendToPlayer(order.GetPlayerId(), "RpcRechargeToHall", req)

	// if order.GetTax() < 0 {
	// 	req := &backstage_hall.Recharge{
	// 		PlayerId:     order.PlayerId,
	// 		RechargeGold: easygo.NewInt64(order.GetTax()), //此处是扣税
	// 		OrderId:      order.OrderId,
	// 		SourceType:   easygo.NewInt32(for_game.GOLD_TYPE_EXTRA_MONEY),
	// 	}
	// 	//通知大厅完成扣税
	// 	SendToRandHall("RpcRechargeToHall", req)
	// }
	//修改订单
	// for_game.SetOrder(order)

	return nil
}

//订单操作
func OptOrder(orderid string, opt int32, operator string, note string) *base.Fail {
	order := for_game.GetRedisOrderObj(orderid)
	if order == nil {
		return easygo.NewFailMsg("订单不存在")
	}

	if order.GetStatus() == opt {
		return easygo.NewFailMsg("操作已完成")
	}

	order.SetStatus(opt)
	order.SetNote(note)
	order.SetOperator(operator)
	order.SetOverTime(for_game.GetMillSecond())
	order.SaveToMongo() //保存Redis中的订单数据

	if opt == 3 {
		req := &server_server.Recharge{
			PlayerId:     easygo.NewInt64(order.GetPlayerId()),
			RechargeGold: easygo.NewInt64(-order.GetChangeGold()), //此处的订单金额，是扣税后的金额
			OrderId:      easygo.NewString(order.GetOrderId()),
			SourceType:   easygo.NewInt32(order.GetSourceType()),
		}
		//通知大厅完成订单
		SendToPlayer(order.GetPlayerId(), "RpcRechargeToHall", req)
	}

	return nil
}
