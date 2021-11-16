// 管理后台为[浏览器]提供的服务
//用户管理

package backstage

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo/bson"
)

//重置指定时间范围的用户留存报表
func (self *cls4) RpcMakePlayerKeepReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	startTime := reqMsg.GetBeginTimestamp()
	endTime := reqMsg.GetEndTimestamp()
	MakePlayerKeepReport(startTime, endTime)

	list, count := GetPlayerKeepReport(reqMsg, user)
	return &brower_backstage.PlayerKeepReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询用户留存报表
func (self *cls4) RpcPlayerKeepReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetPlayerKeepReport(reqMsg, user)

	return &brower_backstage.PlayerKeepReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询日活跃报表
func (self *cls4) RpcPlayerActiveReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetPlayerActiveReportList(reqMsg)

	return &brower_backstage.PlayerActiveReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询周活跃报表
func (self *cls4) RpcPlayerWeekActiveReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetPlayerWeekActiveReportList(reqMsg)

	return &brower_backstage.PlayerActiveReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询月活跃报表
func (self *cls4) RpcPlayerMonthActiveReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetPlayerMonthActiveReportList(reqMsg)

	return &brower_backstage.PlayerActiveReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询用户行为报表
func (self *cls4) RpcPlayerBehaviorReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetPlayerBehaviorReport(reqMsg)

	return &brower_backstage.PlayerBehaviorReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询出入款汇总报表
func (self *cls4) RpcInOutCashSumReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetInOutCashSumReport(reqMsg)

	return &brower_backstage.InOutCashSumReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询埋点注册登录报表
func (self *cls4) RpcRegisterLoginReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	if reqMsg.GetCurPage() == 1 {
		for_game.SaveRedisRegisterLoginReportToMongo() //保存redis中的数据
	}
	list, count := GetRegisterLoginReport(reqMsg, user)

	return &brower_backstage.RegisterLoginReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询运营渠道汇总报表
func (self *cls4) RpcOperationChannelReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetOperationChannelReportList(reqMsg, user)

	return &brower_backstage.OperationChannelReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询渠道报表
func (self *cls4) RpcChannelReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetChannelReportList(reqMsg)

	return &brower_backstage.ChannelReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询运营渠道汇总报表曲线图
func (self *cls4) RpcOperationChannelLine(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	if reqMsg.Keyword == nil || reqMsg.GetKeyword() == "" {
		return easygo.NewFailMsg("请先选择渠道")
	}

	list, _ := GetOperationChannelReportList(reqMsg, user)
	line := &brower_backstage.OperationChannelReportLine{}

	for _, k := range list {
		// 折线图数据
		line.TaxDate = append(line.TaxDate, k.GetCreateTime())
		line.RegCount = append(line.RegCount, k.GetRegCount())
		line.LoginCount = append(line.LoginCount, k.GetLoginCount())
		// line.ShopOrderSumCount = append(line.ShopOrderSumCount, k.GetShopOrderSumCount())
		// line.ShopDealSumAmount = append(line.ShopDealSumAmount, k.GetShopDealSumAmount())
		// line.RechargeSumAmount = append(line.RechargeSumAmount, k.GetRechargeSumAmount())
		// line.WithdrawSumAmount = append(line.WithdrawSumAmount, k.GetWithdrawSumAmount())
	}

	return &brower_backstage.OperationChannelReportLineResponse{
		Line: line,
	}
}

//查询文章报表
func (self *cls4) RpcArticleReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetArticleReportList(reqMsg)
	return &brower_backstage.ArticleReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询推送通知报表
func (self *cls4) RpcNoticeReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetNoticeReportList(reqMsg)

	return &brower_backstage.ArticleReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询社交广场报表
func (self *cls4) RpcSquareReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetSquareReportList(reqMsg)

	return &brower_backstage.SquareReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//活动报表查询
func (self *cls4) RpcQueryActivityReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetActivityReportList(reqMsg)

	return &brower_backstage.ActivityReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//广告报表查询
func (self *cls4) RpcAdvReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetAdvReportList(reqMsg)

	return &brower_backstage.AdvReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//附近的人引导项报表查询
func (self *cls4) RpcNearbyAdvReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetNearbyAdvReportList(reqMsg)

	return &brower_backstage.NearReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

func (self *cls4) RpcEditRegisterLoginReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.RegisterLoginReport) easygo.IMessage {
	find := bson.M{"_id": reqMsg.GetCreateTime()}
	update := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, "bak_report_register_login", find, update, true)
	return easygo.EmptyMsg
}

func (self *cls4) RpcEditPlayerKeepReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PlayerKeepReport) easygo.IMessage {
	find := bson.M{"_id": reqMsg.GetCreateTime()}
	update := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, "bak_report_player_keep", find, update, true)
	return easygo.EmptyMsg
}

func (self *cls4) RpcEditOperationChannelReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.OperationChannelReport) easygo.IMessage {
	find := bson.M{"_id": reqMsg.GetId()}
	update := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, "bak_report_operation_channel", find, update, true)
	return easygo.EmptyMsg
}

//查询用户回归报表
func (self *cls4) RpcQueryRecallReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	if reqMsg.GetCurPage() == 1 {
		base := for_game.GetRedisRecallReportMgr(easygo.GetToday0ClockTimestamp())
		base.SaveToMongo()
	}

	findBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		findBson["_id"] = bson.M{"$gte": easygo.Get0ClockTimestamp(reqMsg.GetBeginTimestamp()), "$lte": easygo.Get0ClockTimestamp(reqMsg.GetEndTimestamp())}
	}
	ls, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_RECALL_REPORT, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), "-_id")

	var list []*share_message.RecallReport
	for _, li := range ls {
		one := &share_message.RecallReport{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}

	return &brower_backstage.RecallReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询用户回归日志
func (self *cls4) RpcQueryRecallPlayerLog(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	// logs.Info("RpcQueryRecallPlayerLog:", reqMsg)
	if reqMsg.GetCurPage() == 1 {
		base := for_game.GetRedisRecallReportMgr(reqMsg.GetBeginTimestamp())
		base.SaveToMongo()
	}

	findBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		findBson["RecallTime"] = bson.M{"$gte": easygo.Get0ClockTimestamp(reqMsg.GetBeginTimestamp()), "$lt": easygo.Get24ClockTimestamp(reqMsg.GetEndTimestamp())}
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			playerid := QueryPlayerbyAccount(reqMsg.GetKeyword()).GetPlayerId()
			findBson["PlayerId"] = playerid
		case 2:
			var playerids PLAYER_IDS
			plis := GetPlayerLikeNickname(for_game.MONGODB_NINGMENG, reqMsg.GetKeyword())
			for _, pli := range plis {
				playerids = append(playerids, pli.GetPlayerId())
			}
			findBson["PlayerId"] = bson.M{"$in": playerids}
		}

	}

	ls, count := for_game.FindAll(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_RECALLPLAYER_LOG, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), "-RecallTime")

	var list []*share_message.RecallPlayerLog
	for _, li := range ls {
		one := &share_message.RecallPlayerLog{}
		for_game.StructToOtherStruct(li, one)
		player := QueryPlayerbyId(one.GetPlayerId())
		one.NickName = easygo.NewString(player.GetNickName())
		one.Account = easygo.NewString(player.GetAccount())
		list = append(list, one)
	}

	return &brower_backstage.RecallplayerLogResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询虚拟商城日统计报表
func (self *cls4) RpcCoinProductReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		findBson["CreateTime"] = bson.M{"$gte": easygo.Get0ClockTimestamp(reqMsg.GetBeginTimestamp()), "$lte": easygo.Get0ClockTimestamp(reqMsg.GetEndTimestamp())}
	}

	groupBson := bson.M{"_id": "$CreateTime"}
	// groupBson["BuyNum"] = bson.M{"$sum": "$BuyNum"}
	groupBson["BuyCount"] = bson.M{"$sum": "$BuyCount"}
	// groupBson["GiveNum"] = bson.M{"$sum": "$GiveNum"}
	groupBson["GiveCount"] = bson.M{"$sum": "$GiveCount"}
	// groupBson["UserGiveNum"] = bson.M{"$sum": "$UserGiveNum"}
	groupBson["UserGiveCount"] = bson.M{"$sum": "$UserGiveCount"}
	// groupBson["ActGiveNum"] = bson.M{"$sum": "$ActGiveNum"}
	groupBson["ActGiveCount"] = bson.M{"$sum": "$ActGiveCount"}
	groupBson["GoldSum"] = bson.M{"$sum": "$GoldSum"}
	groupBson["CoinSum"] = bson.M{"$sum": "$CoinSum"}

	m := []bson.M{
		{"$match": findBson},
		{"$group": groupBson},
		{"$sort": bson.M{"_id": -1}},
	}

	ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT_REPORT, m, 0, 0)
	count := len(ls)

	lis := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT_REPORT, m, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()))
	var list []*share_message.CoinProductReport
	for _, li := range lis {
		one := &share_message.CoinProductReport{
			CreateTime: easygo.NewInt64(li.(bson.M)["_id"].(int64)),
			// BuyNum:        easygo.NewInt64(easygo.StringToInt64noErr(easygo.AnytoA(li.(bson.M)["BuyNum"]))),
			BuyCount: easygo.NewInt64(easygo.StringToInt64noErr(easygo.AnytoA(li.(bson.M)["BuyCount"]))),
			// GiveNum:       easygo.NewInt64(easygo.StringToInt64noErr(easygo.AnytoA(li.(bson.M)["GiveNum"]))),
			GiveCount: easygo.NewInt64(easygo.StringToInt64noErr(easygo.AnytoA(li.(bson.M)["GiveCount"]))),
			// UserGiveNum:   easygo.NewInt64(easygo.StringToInt64noErr(easygo.AnytoA(li.(bson.M)["UserGiveNum"]))),
			UserGiveCount: easygo.NewInt64(easygo.StringToInt64noErr(easygo.AnytoA(li.(bson.M)["UserGiveCount"]))),
			// ActGiveNum:    easygo.NewInt64(easygo.StringToInt64noErr(easygo.AnytoA(li.(bson.M)["ActGiveNum"]))),
			ActGiveCount: easygo.NewInt64(easygo.StringToInt64noErr(easygo.AnytoA(li.(bson.M)["ActGiveCount"]))),
			GoldSum:      easygo.NewInt64(easygo.StringToInt64noErr(easygo.AnytoA(li.(bson.M)["GoldSum"]))),
			CoinSum:      easygo.NewInt64(easygo.StringToInt64noErr(easygo.AnytoA(li.(bson.M)["CoinSum"]))),
		}
		start := easygo.Get0ClockTimestamp(one.GetCreateTime())
		end := easygo.Get24ClockTimestamp(one.GetCreateTime())
		findBson := bson.M{}
		m := []bson.M{}
		var ls []interface{}

		types := []int{for_game.COIN_ITEM_GETTYPE_BUY, for_game.COIN_ITEM_GETTYPE_SEND, for_game.COIN_ITEM_GETTYPE_PLAYER_SEND, for_game.COIN_ITEM_GETTYPE_ACTIVITY}
		for _, i := range types { //重新计算当天各项人数
			findBson = bson.M{"CreateTime": bson.M{"$gte": start, "$lt": end}, "GetType": i}
			m = []bson.M{
				{"$match": findBson},
				{"$group": bson.M{"_id": "$PlayerId"}},
			}
			ls = for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_GETPROPS_LOG, m, 0, 0)

			switch i {
			case for_game.COIN_ITEM_GETTYPE_BUY:
				one.BuyNum = easygo.NewInt64(len(ls))
			case for_game.COIN_ITEM_GETTYPE_SEND:
				one.GiveNum = easygo.NewInt64(len(ls))
			case for_game.COIN_ITEM_GETTYPE_PLAYER_SEND:
				one.UserGiveNum = easygo.NewInt64(len(ls))
			case for_game.COIN_ITEM_GETTYPE_ACTIVITY:
				one.ActGiveNum = easygo.NewInt64(len(ls))
			}
		}

		list = append(list, one)
	}
	msg := &brower_backstage.CoinProductReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//查询虚拟商城日明细报表
func (self *cls4) RpcCoinProductDetailReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	queryTime := easygo.GetToday0ClockTimestamp()
	if reqMsg.Id != nil && reqMsg.GetId() != 0 {
		queryTime = reqMsg.GetId()
		if queryTime > 9999999999 {
			queryTime /= 1000
		}
		findBson["CreateTime"] = queryTime
	}

	pageSize := int(reqMsg.GetPageSize())
	curPage := int(reqMsg.GetCurPage())

	projectBson := bson.M{"ProductId": 1}
	projectBson["BuyNum"] = 1
	projectBson["BuyCount"] = 1
	projectBson["GiveNum"] = 1
	projectBson["GiveCount"] = 1
	projectBson["UserGiveNum"] = 1
	projectBson["UserGiveCount"] = 1
	projectBson["ActGiveNum"] = 1
	projectBson["ActGiveCount"] = 1
	projectBson["GoldSum"] = 1
	projectBson["CoinSum"] = 1
	projectBson["ProductType"] = "$Product.PropsType"
	projectBson["ProductName"] = "$Product.Name"
	projectBson["EffectiveTime"] = "$Product.EffectiveTime"

	m := []bson.M{
		{"$match": findBson},
		{"$lookup": bson.M{
			"from":         for_game.TABLE_COIN_PRODUCT,
			"localField":   "ProductId",
			"foreignField": "_id",
			"as":           "Product",
		}},
		{"$project": projectBson},
		{"$unwind": "$ProductName"},
		{"$unwind": "$ProductType"},
		{"$unwind": "$EffectiveTime"},
		{"$sort": bson.M{"BuyCount": -1}},
	}

	ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT_REPORT, m, 0, 0)
	count := len(ls)
	var list []*share_message.CoinProductReport
	if count > 0 {
		lis := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT_REPORT, m, pageSize, curPage)

		for _, li := range lis {
			one := &share_message.CoinProductReport{}
			for_game.StructToOtherStruct(li, one)
			one.ProductName = easygo.NewString(li.(bson.M)["ProductName"])
			one.ProductType = easygo.NewInt32(li.(bson.M)["ProductType"])
			one.EffectiveTime = easygo.NewInt32(li.(bson.M)["EffectiveTime"])
			list = append(list, one)
		}
	}

	msg := &brower_backstage.CoinProductReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return msg
}

//按钮点击行为报表
func (self *cls4) RpcButtonClickReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	if reqMsg.GetCurPage() == 1 {
		for_game.SaveRedisButtonClickReportToMongo() //保存redis中的按钮点击行为报表数据
	}

	findBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		findBson["_id"] = bson.M{"$gte": easygo.Get0ClockTimestamp(reqMsg.GetBeginTimestamp()), "$lte": easygo.Get0ClockTimestamp(reqMsg.GetEndTimestamp())}
	}
	ls, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_BUTTON_CLICK_REPORT, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), "-_id")

	var list []*share_message.ButtonClickReport
	for _, li := range ls {
		one := &share_message.ButtonClickReport{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}

	return &brower_backstage.ButtonClickReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//注册登录页面埋点报表
func (self *cls4) RpcPageRegLogReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	queryBson := bson.M{}
	if reqMsg.SrtType != nil {
		queryBson["Channel"] = reqMsg.GetSrtType()
	}
	if reqMsg.Type != nil {
		queryBson["DicType"] = reqMsg.GetType()
	}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": easygo.Get0ClockTimestamp(reqMsg.GetBeginTimestamp()), "$lte": easygo.Get0ClockTimestamp(reqMsg.GetEndTimestamp())}
	}

	groupBson := bson.M{}
	groupBson["_id"] = "$CreateTime"
	groupBson["CreateTime"] = bson.M{"$addToSet": "$CreateTime"}
	groupBson["Channel"] = bson.M{"$addToSet": "$Channel"}
	groupBson["DicType"] = bson.M{"$addToSet": "$DicType"}
	groupBson["LoginTimes"] = bson.M{"$sum": "$LoginTimes"}
	groupBson["LoginCount"] = bson.M{"$sum": "$LoginCount"}
	groupBson["OneLoginTimes"] = bson.M{"$sum": "$OneLoginTimes"}
	groupBson["OneLoginCount"] = bson.M{"$sum": "$OneLoginCount"}
	groupBson["WxLoginTimes"] = bson.M{"$sum": "$WxLoginTimes"}
	groupBson["WxLoginCount"] = bson.M{"$sum": "$WxLoginCount"}
	groupBson["PhoneRegTimes"] = bson.M{"$sum": "$PhoneRegTimes"}
	groupBson["PhoneRegCount"] = bson.M{"$sum": "$PhoneRegCount"}
	groupBson["RegCodeTimes"] = bson.M{"$sum": "$RegCodeTimes"}
	groupBson["RegCodeCount"] = bson.M{"$sum": "$RegCodeCount"}
	groupBson["UseInfoTimes"] = bson.M{"$sum": "$UseInfoTimes"}
	groupBson["UseInfoCount"] = bson.M{"$sum": "$UseInfoCount"}
	groupBson["ActDevCount"] = bson.M{"$sum": "$ActDevCount"}
	groupBson["ValidActDevCount"] = bson.M{"$sum": "$ValidActDevCount"}

	M := []bson.M{
		{"$match": queryBson},
		{"$group": groupBson},
		{"$sort": bson.M{"_id": -1}},
	}
	ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PAGE_REGLOG_REPORT, M, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()))
	count := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PAGE_REGLOG_REPORT, M, 0, 0)
	var list []*share_message.PageRegLogReport
	for _, li := range ls {
		one := &share_message.PageRegLogReport{
			CreateTime:       easygo.NewInt64(li.(bson.M)["_id"]),
			LoginTimes:       easygo.NewInt64(li.(bson.M)["LoginTimes"]),
			LoginCount:       easygo.NewInt64(li.(bson.M)["LoginCount"]),
			OneLoginTimes:    easygo.NewInt64(li.(bson.M)["OneLoginTimes"]),
			OneLoginCount:    easygo.NewInt64(li.(bson.M)["OneLoginCount"]),
			WxLoginTimes:     easygo.NewInt64(li.(bson.M)["WxLoginTimes"]),
			WxLoginCount:     easygo.NewInt64(li.(bson.M)["WxLoginCount"]),
			PhoneRegTimes:    easygo.NewInt64(li.(bson.M)["PhoneRegTimes"]),
			PhoneRegCount:    easygo.NewInt64(li.(bson.M)["PhoneRegCount"]),
			RegCodeTimes:     easygo.NewInt64(li.(bson.M)["RegCodeTimes"]),
			RegCodeCount:     easygo.NewInt64(li.(bson.M)["RegCodeCount"]),
			UseInfoTimes:     easygo.NewInt64(li.(bson.M)["UseInfoTimes"]),
			UseInfoCount:     easygo.NewInt64(li.(bson.M)["UseInfoCount"]),
			ActDevCount:      easygo.NewInt64(li.(bson.M)["ActDevCount"]),
			ValidActDevCount: easygo.NewInt64(li.(bson.M)["ValidActDevCount"]),
		}
		list = append(list, one)
	}

	return &brower_backstage.PageRegLogReportResponse{
		List:      list,
		PageCount: easygo.NewInt32(len(count)),
	}
}
