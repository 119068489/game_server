// 大厅服务器为[游戏客户端]提供的服务

package wish

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/h5_wish"

	"github.com/astaxie/beego/logs"
)

//我的愿望盒-已收藏的
func (self *WebHttpForClient) RpcGetCollectionBox(common *base.Common, reqMsg *h5_wish.DataPageReq) easygo.IMessage {
	logs.Info("===我的愿望盒-已收藏的 RpcGetCollectionBox pid:%v, reqMsg: %v", common.GetUserId(), reqMsg)
	playerId := common.GetUserId()
	page := int(reqMsg.GetPage())
	pageSize := int(reqMsg.GetPageSize())
	reqType := int(reqMsg.GetReqType())

	var toExchangeCount, exchangedCount, recycleCount int

	if reqType != WISH_SALEOUT && reqType != WISH_ONSALE {
		logs.Error("我的愿望盒-已收藏 类型有误, opType: %d", reqType)
		return easygo.NewFailMsg("参数有误")
	}

	skip := (page - 1) * pageSize
	var resp *h5_wish.MyCollectedBoxResp
	wishBoxList, onSaleCount, SaleOutCount := GetCollectedWishBoxList(playerId, skip, pageSize, reqType)

	if reqType == WISH_ONSALE && page == 1 {
		toExchangeCount, exchangedCount, recycleCount = for_game.CountPlayerWish(playerId)
	}
	//当type为WISH_SALEOUT时OnSaleCount与ExchangeCount无效
	resp = &h5_wish.MyCollectedBoxResp{
		Boxes:           wishBoxList,
		OnSaleCount:     easygo.NewInt32(onSaleCount),
		SaleOutCount:    easygo.NewInt32(SaleOutCount),
		ToExchangeCount: easygo.NewInt32(toExchangeCount),
		ExchangedCount:  easygo.NewInt32(exchangedCount),
		RecycleCount:    easygo.NewInt32(recycleCount),
	}

	return resp
}

//我的愿望盒-所有已收藏
func (self *WebHttpForClient) RpcGetAllCollectionBox(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===我的愿望盒-所有已收藏 RpcGetAllCollectionBox pid:%v", common.GetUserId())
	playerId := common.GetUserId()

	wishBoxList, err := GetAllCollectedWishBoxList(playerId)
	if err != nil {
		return easygo.NewFailMsg(err.Error())
	}
	toExchangeCount, exchangedCount, recycleCount := for_game.CountPlayerWish(playerId)
	resp := &h5_wish.MyCollectedBoxResp{
		Boxes:           wishBoxList,
		ToExchangeCount: easygo.NewInt32(toExchangeCount),
		ExchangedCount:  easygo.NewInt32(exchangedCount),
		RecycleCount:    easygo.NewInt32(recycleCount),
	}
	return resp
}

//我的愿望盒-已许愿
func (self *WebHttpForClient) RpcGetWishBoxList(common *base.Common, reqMsg *h5_wish.DataPageReq) easygo.IMessage {
	logs.Info("===我的愿望盒-已许愿 RpcGetWishBoxList pid:%v, reqMsg: %v", common.GetUserId(), reqMsg)
	playerId := common.GetUserId()
	page := int(reqMsg.GetPage())
	pageSize := int(reqMsg.GetPageSize())
	skip := (page - 1) * pageSize
	reqType := int(reqMsg.GetReqType())

	if reqType != WISH_SALEOUT && reqType != WISH_ONSALE && reqType != WISH_ALL {
		logs.Error("我的愿望盒-已许愿 类型有误, opType: %d", reqType)
		return easygo.NewFailMsg("参数有误")
	}

	wishBoxList, onSaleTotal, saleOutTotal := GetWishDataList(playerId, skip, pageSize, reqType)
	//当type为WISH_SALEOUT时OnSaleCount无效
	resp := &h5_wish.MyCollectedBoxResp{
		Boxes:        wishBoxList,
		OnSaleCount:  easygo.NewInt32(onSaleTotal),
		SaleOutCount: easygo.NewInt32(saleOutTotal),
	}
	return resp
}

//我的愿望盒-已许愿
func (self *WebHttpForClient) RpcGetAllWishBoxList(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===我的愿望盒-已许愿 RpcGetAllWishBoxList pid:%v", common.GetUserId())
	playerId := common.GetUserId()
	wishBoxList, err := GetAllWishDataList(playerId)
	if err != nil {
		return easygo.NewFailMsg(err.Error())
	}
	resp := &h5_wish.MyCollectedBoxResp{
		Boxes: wishBoxList,
	}
	return resp
}

// 收藏或者取消收藏
func (self *WebHttpForClient) RpcRemoveCollectionBox(common *base.Common, reqMsg *h5_wish.CollectionBoxReq) easygo.IMessage {
	logs.Info("=====RpcRemoveCollectionBox=====,common=%v,reqMsg=%v", common, reqMsg)
	ids := reqMsg.GetIdList() //盲盒Id
	if len(ids) < 1 {
		return easygo.NewFailMsg("参数有误")
	}
	playerId := common.GetUserId()
	res := 0
	err := CollectionBoxService(playerId, reqMsg)
	if err != nil {
		res = 1
	}
	return &h5_wish.DefaultResp{
		Result: easygo.NewInt32(res),
	}
}

//我的愿望盒子请求
func (self *WebHttpForClient) RpcMyAllWish(common *base.Common, reqMsg *h5_wish.MyWishReq) easygo.IMessage {
	logs.Info("===我的愿望盒子请求 RpcMyAllWish pid:%v, reqMsg: %v", common.GetUserId(), reqMsg)
	playerId := common.GetUserId()
	playerWishItem, err := MyWishService(playerId, reqMsg)
	if err != nil {
		return easygo.NewFailMsg(err.Error())
	}
	resp := &h5_wish.ProductResp{
		ItemList: playerWishItem,
	}

	if reqMsg.GetPage() == 1 {
		toExchangeCount, exchangedCount, recycleCount := for_game.CountPlayerWish(playerId)
		resp.ToExchangeCount = easygo.NewInt32(toExchangeCount)
		resp.ExchangedCount = easygo.NewInt32(exchangedCount)
		resp.RecycleCount = easygo.NewInt32(recycleCount)
	}

	return resp
}

//试玩一次
func (self *WebHttpForClient) RpcTryOnce(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===试玩一次 RpcTryOnce pid:%v", common.GetUserId())
	err := TryOnceService(common.GetUserId())
	return RpcReturnCommon("RpcTryOnce", nil, err)
}

//收藏盲盒或删除收藏盲盒
func (self *WebHttpForClient) RpcCollectionBox(common *base.Common, reqMsg *h5_wish.CollectionBoxReq) easygo.IMessage {
	logs.Info("=====RpcCollectionBox=====,common=%v,reqMsg=%v", common, reqMsg)
	ids := reqMsg.GetIdList() //盲盒Id
	if len(ids) < 1 {
		return easygo.NewFailMsg("参数有误")
	}
	playerId := common.GetUserId()
	res := 0
	err := CollectionBoxService(playerId, reqMsg)
	if err != nil {
		res = 1
	}
	return &h5_wish.DefaultResp{
		Result: easygo.NewInt32(res),
	}
}

//删除许愿盲盒物品
func (self *WebHttpForClient) RpcDelWishBox(common *base.Common, reqMsg *h5_wish.WishBoxReq) easygo.IMessage {
	logs.Info("===删除许愿盲盒物品 RpcDelWishBox pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	//TODO 逻辑处理
	ids := reqMsg.GetIdList() //盲盒Id
	if len(ids) < 1 {
		return easygo.NewFailMsg("参数有误")
	}

	res := 0
	playerId := common.GetUserId()
	err := for_game.RemovePlayerWishDataByIds(playerId, ids)
	if err != nil {
		res = 1
	}
	return &h5_wish.DefaultResp{
		Result: easygo.NewInt32(res),
	}
}

//兑换盲盒物品 1-失败，2-成功，3-邮费不足
func (self *WebHttpForClient) RpcExchangeBox(common *base.Common, reqMsg *h5_wish.WishBoxReq) easygo.IMessage {
	logs.Info("===兑换盲盒物品 RpcExchangeBox pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	playerId := common.GetUserId()
	if CheckIsFreeze(playerId) {
		return &h5_wish.DefaultResp{
			Result: easygo.NewInt32(for_game.RECYCLE_STATUS_FREEZE),
			Msg:    easygo.NewString("您的许愿池账号已被冻结，请联系官方客服"),
		}
	}

	resp := &h5_wish.DefaultResp{
		Result: easygo.NewInt32(for_game.RECYCLE_STATUS_ERROR),
	}

	freeNum, err := for_game.GetWishMailSection()
	if err != nil {
		logs.Error("获取邮费信息失败：%v", err)
		resp.Msg = easygo.NewString(EXCHANGE_ERR_STR)
		return resp
	}

	player := for_game.GetRedisWishPlayer(playerId)
	if player == nil {
		resp.Msg = easygo.NewString(EXCHANGE_ERR_STR)
		return resp
	}
	//查询地址
	address := player.GetOneAddress(reqMsg.GetAddressId())
	if address == nil {
		logs.Info("兑换时查询用户添加的地址数据为空")
		resp.Msg = easygo.NewString(EXCHANGE_ERR_STR)
		return resp
	}

	playerItems, errNum, editStatusIds := player.DealPlayerWishItem(for_game.WISH_EXCHANGED, reqMsg)
	if errNum == for_game.DEAL_FAULT {
		resp.Msg = easygo.NewString(EXCHANGE_ERR_STR)
		return resp
	} else if errNum == for_game.DEAL_PLAYER_LIMIT {
		resp.Msg = easygo.NewString("此账号不允许兑换")
		return resp
	}

	result := ExchangeGoods(player, playerItems, freeNum.GetFreeNumber(), address)
	resp.Result = easygo.NewInt32(result)
	if result != for_game.EXCHANGE_STATUS_SUCCESS {
		resp.Msg = easygo.NewString(EXCHANGE_ERR_STR)
		//失败重制玩家物品状态，若为全为预售商品则无需处理则
		if result != for_game.EXCHANGE_STATUS_ISPRESALE {
			for_game.ExchangeToPlayerItem(editStatusIds, for_game.WISH_TO_EXCHANGE)
		}
	}
	return resp
}

//回收盲盒物品
func (self *WebHttpForClient) RpcRecycleGoods(common *base.Common, reqMsg *h5_wish.WishBoxReq) easygo.IMessage {
	logs.Info("===回收盲盒物品 RpcRecycleGoods pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	playerId := common.GetUserId()
	if CheckIsFreeze(playerId) {
		return &h5_wish.DefaultResp{
			Result: easygo.NewInt32(for_game.EXCHANGE_STATUS_FREEZE),
			Msg:    easygo.NewString("您的许愿池账号已被冻结，请联系官方客服"),
		}
	}

	resp := &h5_wish.DefaultResp{
		Result: easygo.NewInt32(for_game.EXCHANGE_STATUS_ERROR),
	}
	recycleRatio, err := for_game.GetRecycleRatio(for_game.RECYCLE_BY_PLAYER)
	if err != nil || recycleRatio == 0 {
		logs.Error("回收比例为0")
		resp.Msg = easygo.NewString(RECYCLE_ERR_STR)
		return resp
	}

	player := for_game.GetRedisWishPlayer(playerId)
	if player == nil {
		logs.Error("回收盲盒物品获取玩家(%v)信息失败", common.GetUserId())
		resp.Msg = easygo.NewString(RECYCLE_ERR_STR)
		return resp
	}
	playerItems, errNum, editStatusIds := player.DealPlayerWishItem(for_game.WISH_RECYCLED, reqMsg)
	if errNum == for_game.DEAL_FAULT {
		resp.Msg = easygo.NewString(RECYCLE_ERR_STR)
		return resp
	} else if errNum == for_game.DEAL_PLAYER_LIMIT {
		resp.Msg = easygo.NewString("此账号不允许回收")
		return resp
	}

	result := for_game.RECYCLE_STATUS_ERROR
	if reqMsg.GetBankCardId() != "" {
		result, err = RecycleGoods(player, playerItems, for_game.RECYCLE_BY_PLAYER, recycleRatio, reqMsg)
	} else {
		result, err = RecycleGoodsToDiamond(player, playerItems, for_game.RECYCLE_BY_PLAYER, recycleRatio, reqMsg)
	}
	resp.Result = easygo.NewInt32(result)
	if err != nil {
		resp.Msg = easygo.NewString(err.Error())
		//失败重制玩家物品状态
		for_game.ExchangeToPlayerItem(editStatusIds, for_game.WISH_TO_EXCHANGE)
	}
	return resp
}

//请求收货地址
func (self *WebHttpForClient) RpcGetAddressList(common *base.Common, reqMsg *h5_wish.DataPageReq) easygo.IMessage {
	logs.Info("===请求收货地址 RpcGetAddressList pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	//TODO 校验参数
	playerId := common.GetUserId()
	page, pageSize := for_game.MakePageAndPageSize(reqMsg.GetPage(), reqMsg.GetPageSize())
	list := GetAddressListByUid(playerId, page, pageSize)
	return &h5_wish.AddressListResp{
		List:  list,
		Count: easygo.NewInt32(len(list)),
	}
}

//添加收货地址
func (self *WebHttpForClient) RpcAddAddress(common *base.Common, reqMsg *h5_wish.WishAddress) easygo.IMessage {
	logs.Info("===添加收货地址 RpcAddAddress pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	playerId := common.GetUserId()
	done := AddAddress(playerId, reqMsg)
	if !done {
		return &h5_wish.DefaultResp{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString("添加收货地址失败，请联系客服"),
		}
	}
	return &h5_wish.DefaultResp{
		Result: easygo.NewInt32(0),
	}
}

//修改收货地址
func (self *WebHttpForClient) RpcEditAddress(common *base.Common, reqMsg *h5_wish.WishAddress) easygo.IMessage {
	logs.Info("===修改收货地址 RpcEditAddress pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	playerId := common.GetUserId()
	done := EditAddress(playerId, reqMsg)
	if !done {
		return &h5_wish.DefaultResp{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString("修改收货地址失败，请联系客服"),
		}
	}
	return &h5_wish.DefaultResp{
		Result: easygo.NewInt32(0),
	}
}

//移除收货地址
func (self *WebHttpForClient) RpcRemoveAddress(common *base.Common, reqMsg *h5_wish.RemoveAddressReq) easygo.IMessage {
	logs.Info("===移除收货地址 RpcRemoveAddress pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	playerId := common.GetUserId()
	addressId := reqMsg.GetAddressId()
	done := RemoveAddress(playerId, addressId)
	if !done {
		return &h5_wish.DefaultResp{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString("移除收货地址失败，请联系客服"),
		}
	}
	return &h5_wish.DefaultResp{
		Result: easygo.NewInt32(0),
	}
}

//获取未读待兑换的盲盒个数
func (self *WebHttpForClient) RpcUnReadWishNum(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===获取未读待兑换的盲盒个数 RpcUnReadWishNum pid:%v", common.GetUserId())
	count := for_game.CountUnReadToExchangePlayerWish(common.GetUserId())
	resp := &h5_wish.JustNumberResp{Result: easygo.NewInt32(count)}
	return resp
}

//获取待兑换的盲盒个数
func (self *WebHttpForClient) RpcToExchangeWishNum(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===获取待兑换的盲盒个数 RpcToExchangeWishNum pid:%v", common.GetUserId())
	count := for_game.CountToExchangePlayerWish(common.GetUserId())
	resp := &h5_wish.JustNumberResp{Result: easygo.NewInt32(count)}
	return resp
}

//获取回收比例
func (self *WebHttpForClient) RpcRecycleRatio(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===获取回收比例 RpcRecycleRatio pid:%v", common.GetUserId())
	//TODO 等待表格设计
	//count := for_game.CountToExchangePlayerWish(common.GetUserId())
	ratio, err := for_game.GetRecycleRatio(for_game.RECYCLE_BY_PLAYER)
	if err != nil {
		return easygo.NewFailMsg("回收比例失败")
	}
	resp := &h5_wish.JustNumberResp{Result: easygo.NewInt32(ratio)}
	return resp
}

//获取不同地区邮费
func (self *WebHttpForClient) RpcAreaPostage(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===获取不同地区邮费 RpcAreaPostage pid:%v", common.GetUserId())
	postage, err := for_game.GetWishMailSection()

	if err != nil {
		return easygo.NewFailMsg("获取邮费失败")
	}
	return postage
}

//盲盒中所有商品详情
func (self *WebHttpForClient) RpcBoxProduct(common *base.Common, reqMsg *h5_wish.DareReq) easygo.IMessage {
	logs.Info("===盲盒中所有商品详情 RpcBoxProduct pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	boxId := reqMsg.GetBoxId()
	if boxId == 0 {
		logs.Error("---商品详情 productId 为 0")
		return easygo.NewFailMsg("参数有误")
	}
	//TODO 逻辑处理
	BoxDetailArr, err := BoxProductDetailService(boxId)
	if err != nil {
		return easygo.NewFailMsg(err.Error())
	}
	return &h5_wish.BoxProductResp{
		ProductList: BoxDetailArr,
	}
}

//获取用户银行卡信息
func (self *WebHttpForClient) RpcGetUserIdBankCards(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=====获取用户银行卡信息 RpcGetPlayerBankCards pid:%v", common.GetUserId())
	player := for_game.GetRedisWishPlayer(common.GetUserId())
	if player == nil {
		logs.Error("获取用户信息失败")
		return easygo.NewFailMsg("获取用户信息失败")
	}

	user := &h5_wish.UserInfoReq{
		UserId: easygo.NewInt64(player.GetPlayerId()),
	}
	resp, err := SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcGetUserBankCardsInfo", user, player.GetPlayerId())
	if err != nil {
		logs.Error("获取用户银行卡信息失败：%v", err)
		return easygo.NewFailMsg("获取用户银行卡信息失败")
	}
	return resp
}

//获取默认配置
func (self *WebHttpForClient) RpcGetConfig(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	userId := common.GetUserId()
	logs.Info("=====获取默认配置 RpcGetConfig pid:%v", userId)
	errKey := make([]string, 0)
	postage, err := for_game.GetWishMailSection()
	if err != nil {
		errKey = append(errKey, "postage")
	}
	recycleReason := make([]*h5_wish.Menu, 0)
	reason, err2 := for_game.GetRecycleReason()
	if err2 != nil || len(reason) == 0 {
		errKey = append(errKey, "RecycleRatio")
	} else {
		for _, v := range reason {
			recycleReason = append(recycleReason,
				&h5_wish.Menu{
					Id:   easygo.NewInt64(v.GetId()),
					Name: easygo.NewString(v.GetReason())})
		}
	}

	currency := for_game.GetCurrencyCfg()
	destCurrency := &h5_wish.WishCurrencyConversionCfg{}
	if currency == nil {
		errKey = append(errKey, "Conversion")
	} else {
		for_game.StructToOtherStruct(currency, destCurrency)
	}
	resp := &h5_wish.ConfigResp{
		Postage:       postage,
		RecycleReason: recycleReason,
		Conversion:    destCurrency,
	}

	//获取用户回收比例
	payConfig := for_game.GetConfigWishPayment()
	if payConfig != nil {
		resp.RecycleRatio = easygo.NewInt32(payConfig.GetPlayer())
	} else {
		errKey = append(errKey, "payConfig")
	}

	player := for_game.GetRedisWishPlayer(userId)
	if player == nil {
		errKey = append(errKey, "player")
		resp.ErrKey = errKey
		return resp
	}
	//用户单日回收参数
	dayPeriod := for_game.GetPlayerPeriod(player.GetPlayerId()).DayPeriod
	if dayPeriod != nil {
		resp.PlayerRecycleMoneyTime = easygo.NewInt32(dayPeriod.FetchInt64(for_game.WISH_OUT_MONEY_TIME))
		resp.PlayerRecycleDiamondSum = easygo.NewInt64(dayPeriod.FetchInt64(for_game.WISH_OUT_DIAMOND_SUM))
		resp.PlayerRecycleDiamondTime = easygo.NewInt32(dayPeriod.FetchInt64(for_game.WISH_OUT_DIAMOND_TIME))
		resp.PlayerRecycleMoneySum = easygo.NewInt64(dayPeriod.FetchInt64(for_game.WISH_OUT_MONEY_SUM))
	} else {
		errKey = append(errKey, "dayPeriod")
	}
	payment := for_game.GetConfigWishPayment()
	if payment != nil {
		resp.DayRecycleMoneyTime = easygo.NewInt32(payment.GetDayMoneyTopCount())
		resp.DayRecycleMoneySum = easygo.NewInt64(payment.GetDayMoneyTop())
		resp.DayRecycleDiamondTime = easygo.NewInt32(payment.GetDayDiamondTopCount())
		resp.DayRecycleDiamondSum = easygo.NewInt64(payment.GetDayDiamondTop())
	} else {
		errKey = append(errKey, "payment")
	}
	resp.ErrKey = errKey
	return resp
}

//添加许愿池埋点记录  请求类型：5-点击发起挑战
func (self *WebHttpForClient) RpcReportWishLog(common *base.Common, reqMsg *h5_wish.TypeReq) easygo.IMessage {
	logs.Info("=====添加许愿池埋点记录 RpcAddReportWishLog pid:%v, type", common.GetUserId(), reqMsg)
	reqType := reqMsg.GetType()
	if reqType <= for_game.WISH_REPORT_CHALLENGE {
		AddReportWishLogService(common.GetUserId(), int(reqType))
	} else {
		return easygo.NewFailMsg("参数有误")
	}
	return &base.Empty{}
}

//回收责任说明
func (self *WebHttpForClient) RpcRecycleDesc(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=====回收责任说明 RpcRecycleDesc pid:%v", common.GetUserId())
	note := for_game.GetOneRecycleNote()
	return &h5_wish.DefaultResp{
		Msg: easygo.NewString(note),
	}
}
