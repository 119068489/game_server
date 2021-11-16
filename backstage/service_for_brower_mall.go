// 管理后台为[浏览器]提供的服务

package backstage

//虚拟市场管理
import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"

	"github.com/astaxie/beego/logs"

	"github.com/akqp2019/mgo/bson"
)

//硬币管理列表
func (self *cls4) RpcCoinItemList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-Sort", "EndTime"}
	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["Platform"] = reqMsg.GetListType()
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_RECHARGE, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.CoinRecharge
	for _, li := range lis {
		one := &share_message.CoinRecharge{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.CoinItemListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//硬币保存
func (self *cls4) RpcSaveCoinItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.CoinRecharge) easygo.IMessage {
	if reqMsg.Coin == nil && reqMsg.GetCoin() < 1 {
		return easygo.NewFailMsg("硬币数量错误")
	}

	if for_game.IS_FORMAL_SERVER {
		if reqMsg.GetPrice() < 100 {
			return easygo.NewFailMsg("价格不能小于1元")
		}
	}

	if reqMsg.GetPlatform() == 0 {
		return easygo.NewFailMsg("请先选择平台")
	}

	if reqMsg.GetMonthFirst() > 0 && reqMsg.GetRebate() < 100 {
		return easygo.NewFailMsg("月首次购买赠送和折扣不能同时设置")
	}

	msg := fmt.Sprintf("修改硬币商品:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_COIN_RECHARGE))
		msg = fmt.Sprintf("添加硬币商品:%d", reqMsg.GetId())
	}

	nowTime := easygo.NowTimestamp()
	if reqMsg.GetEndTime() > nowTime && reqMsg.GetDisPrice() == 0 {
		return easygo.NewFailMsg("活动期间，折扣价不能为空")
	}

	if reqMsg.GetEndTime() == 0 && reqMsg.GetDisPrice() > 0 && reqMsg.GetRebate() < 100 {
		if reqMsg.GetStatus() != for_game.COIN_PRODUCT_STATUS_DOWN {
			return easygo.NewFailMsg("请先设置活动时间")
		}
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_RECHARGE, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.COINS_MANAGE, msg)

	return easygo.EmptyMsg
}

//系统赠送硬币
func (self *cls4) RpcGiveCoin(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.IdStr == nil && reqMsg.GetIdStr() == "" {
		return easygo.NewFailMsg("用户柠檬号不能为空")
	}

	if reqMsg.GetId32() == 0 {
		return easygo.NewFailMsg("赠送数量不能为0")
	}

	playerId := reqMsg.GetId64()
	if playerId == 0 {
		player := QueryPlayerbyAccount(reqMsg.GetIdStr())
		if player != nil {
			playerId = player.GetPlayerId()
		}
	}

	var st int32
	var msg string
	coins := reqMsg.GetId32()
	switch reqMsg.GetNote() {
	case "give":
		if coins < 0 {
			coins = -coins
		}
		st = for_game.COIN_TYPE_SYSTEM_IN
		msg = fmt.Sprintf("系统赠送[%s]硬币[%d]个", reqMsg.GetIdStr(), reqMsg.GetId32())
	case "back":
		if coins > 0 {
			coins = -coins
		}
		switch reqMsg.GetObjId() {
		case 1: //回收个人兑换
			st = for_game.COIN_TYPE_CONFISCATE_OUT
			msg = fmt.Sprintf("系统罚没[%s]硬币[%d]个", reqMsg.GetIdStr(), reqMsg.GetId32())
		case 2: //回收平台赠送
			st = for_game.COIN_TYPE_SYSTEM_OUT
			msg = fmt.Sprintf("系统回收[%s]硬币[%d]个", reqMsg.GetIdStr(), reqMsg.GetId32())
		default:
			return easygo.NewFailMsg("ObjId参数错误")
		}

	default:
		return easygo.NewFailMsg("IdStr参数错误")
	}

	req := &server_server.Recharge{
		PlayerId:     easygo.NewInt64(playerId),
		RechargeGold: easygo.NewInt64(coins),
		SourceType:   easygo.NewInt32(st),
		Note:         easygo.NewString(msg),
	}
	result := SendToPlayer(playerId, "RpcChangeCoinsToHall", req) //通知大厅
	err := for_game.ParseReturnDataErr(result)
	if err != nil {
		return err
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.COINS_MANAGE, msg)

	return easygo.EmptyMsg
}

//硬币流水日志
func (self *cls4) RpcQueryCoinChangeLog(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	if reqMsg.GetCurPage() == 1 {
		for_game.SaveCoinChangeLogToMongoDB()
	}

	findBson := bson.M{}
	sort := []string{"-CreateTime"}
	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["PayType"] = reqMsg.GetListType()
	}
	if reqMsg.DownType != nil && reqMsg.GetDownType() != 0 {
		findBson["SourceType"] = reqMsg.GetDownType()
	}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			findBson["Extend.OrderId"] = reqMsg.GetKeyword()
		case 2:
			base := QueryPlayerbyAccount(reqMsg.GetKeyword())
			if base != nil {
				findBson["PlayerId"] = base.GetPlayerId()
			}
		case 3: //注单号查询
			findBson["Extend.MerchantId"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
		}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_COINCHANGELOG, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*for_game.CoinLog
	for _, li := range lis {
		lione := &for_game.CoinLog{}
		for_game.StructToOtherStruct(li, lione)
		list = append(list, lione)
	}

	var msg []*brower_backstage.CoinLogList
	for _, l := range list {
		base := QueryPlayerbyId(l.GetPlayerId())
		one := &brower_backstage.CoinLogList{
			InLine: &share_message.CoinChangeLog{
				LogId:      easygo.NewInt64(l.GetLogId()),
				PlayerId:   easygo.NewInt64(l.GetPlayerId()),
				Account:    easygo.NewString(base.GetAccount()),
				ChangeCoin: easygo.NewInt64(l.GetChangeCoin()),
				PayType:    easygo.NewInt32(l.GetPayType()),
				SourceType: easygo.NewInt32(l.GetSourceType()),
				CurCoin:    easygo.NewInt64(l.GetCurCoin()),
				Coin:       easygo.NewInt64(l.GetCoin()),
				CurBCoin:   easygo.NewInt64(l.GetCurBCoin()),
				BCoin:      easygo.NewInt64(l.GetBCoin()),
				Note:       easygo.NewString(l.GetNote()),
				CreateTime: easygo.NewInt64(l.GetCreateTime()),
			},
			Extend: l.Extend,
		}
		msg = append(msg, one)
	}

	return &brower_backstage.CoinChangeLogResponse{
		List:      msg,
		PageCount: easygo.NewInt32(count),
	}
}

//道具库列表
func (self *cls4) RpcQueryPropsItemList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-_id"}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["PropsType"] = reqMsg.GetListType()
	}

	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		findBson["UpdateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			findBson["_id"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		case 2:
			findBson["Name"] = reqMsg.GetKeyword()
		default:
			return easygo.NewFailMsg("搜索类型错误")
		}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PROPS_ITEM, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.PropsItem
	for _, li := range lis {
		one := &share_message.PropsItem{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.PropsItemResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//道具保存
func (self *cls4) RpcSavePropsItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PropsItem) easygo.IMessage {
	if reqMsg.Name == nil && reqMsg.GetName() == "" {
		return easygo.NewFailMsg("请先填写道具名称")
	}

	if reqMsg.PropsType == nil && reqMsg.GetPropsType() == 0 {
		return easygo.NewFailMsg("请先选择道具类型")
	}

	if reqMsg.Id == nil {
		return easygo.NewFailMsg("道具ID不能为空")
	}

	if reqMsg.GetPropsType() == 1 && reqMsg.GetUseType() > 1 {
		return easygo.NewFailMsg("道具使用类型错误")
	}

	one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_PROPS_ITEM, bson.M{"_id": reqMsg.GetId()})
	if one != nil && reqMsg.UpdateTime == nil {
		return easygo.NewFailMsg("道具ID重复")
	}

	msg := fmt.Sprintf("修改硬币商品:%d", reqMsg.GetId())
	if reqMsg.UpdateTime == nil && reqMsg.GetUpdateTime() == 0 {
		msg = fmt.Sprintf("添加硬币商品:%d", reqMsg.GetId())
	}
	reqMsg.UpdateTime = easygo.NewInt64(util.GetMilliTime())

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_PROPS_ITEM, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PROPS_MANAGE, msg)

	return easygo.EmptyMsg
}

//道具Ids查询道具列表
func (self *cls4) RpcQueryPropsItemByIds(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	ids := reqMsg.GetIds64()
	findBson := bson.M{"_id": bson.M{"$in": ids}}
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PROPS_ITEM, findBson, 0, 0)
	var list []*share_message.PropsItem
	for _, li := range lis {
		one := &share_message.PropsItem{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.PropsItemResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//商城商品管理列表
func (self *cls4) RpcCoinProductList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-Sort", "-CreateTime"}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["PropsType"] = reqMsg.GetListType()
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.Keyword != nil {
		switch reqMsg.GetType() {
		case 1:
			findBson["Name"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
		case 2:
			id := easygo.StringToInt64noErr(reqMsg.GetKeyword())
			findBson["_id"] = id
		}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.CoinProduct
	for _, li := range lis {
		one := &share_message.CoinProduct{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.CoinProductResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//商城商品数组
func (self *cls4) RpcQueryCoinProductArry(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-Sort", "-CreateTime"}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		findBson["Name"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["PropsType"] = reqMsg.GetListType()
	}

	lis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT, findBson, 0, 0, sort...)
	var list []*brower_backstage.KeyValueStr
	for _, li := range lis {
		one := &share_message.CoinProduct{}
		for_game.StructToOtherStruct(li, one)

		name := one.GetName() + easygo.AnytoA(easygo.If(one.GetEffectiveTime() > 0, easygo.AnytoA(one.GetEffectiveTime())+"天", "永久"))
		oneLi := &brower_backstage.KeyValueStr{
			Value: easygo.NewString(name),
			Key:   easygo.NewString(one.GetId()),
		}
		list = append(list, oneLi)
	}

	return &brower_backstage.KeyValueResponse{
		ListStr: list,
	}
}

//商品保存
func (self *cls4) RpcSaveCoinProduct(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.CoinProduct) easygo.IMessage {
	if reqMsg.Name == nil && reqMsg.GetName() == "" {
		return easygo.NewFailMsg("请先填写商品名")
	}

	if reqMsg.PropsType == nil && reqMsg.GetPropsType() == 0 {
		return easygo.NewFailMsg("请先选择商品类型")
	}

	if reqMsg.EffectiveTime == nil {
		return easygo.NewFailMsg("请先选择商品时效")
	}

	if reqMsg.Sort == nil {
		return easygo.NewFailMsg("请先填写排序权重")
	}

	if reqMsg.PropsId == nil && reqMsg.GetPropsId() == 0 {
		return easygo.NewFailMsg("请先选择道具")
	}

	if reqMsg.Coin == nil && reqMsg.Price == nil {
		return easygo.NewFailMsg("价格不能为空")
	}

	if reqMsg.Coin == nil {
		reqMsg.Coin = easygo.NewInt64(0)
	}

	if reqMsg.IosCoin == nil {
		reqMsg.IosCoin = easygo.NewInt64(0)
	}

	if reqMsg.Price == nil {
		reqMsg.Price = easygo.NewInt64(0)
	}

	nowTime := easygo.NowTimestamp()
	if reqMsg.GetEndTime() > nowTime && reqMsg.GetDisCoin() == 0 {
		if reqMsg.GetStatus() != for_game.COIN_PRODUCT_STATUS_DOWN {
			return easygo.NewFailMsg("活动期间，折扣价不能为空")
		}
	}

	msg := fmt.Sprintf("修改商品:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_COIN_PRODUCT))
		reqMsg.CreateTime = easygo.NewInt64(util.GetMilliTime())
		reqMsg.ProductNum = easygo.NewInt64(-1)
		msg = fmt.Sprintf("添加商品:%d", reqMsg.GetId())
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.MALL_MANAGE, msg)

	return easygo.EmptyMsg
}

//玩家背包列表
/*
func (self *cls4) RpcPlayerBagItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	if reqMsg.Id != nil {
		bagObj := for_game.GetRedisPlayerBagItemObj(reqMsg.GetId())
		bagObj.SaveToMongo()
	}

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	findBson := bson.M{}
	queryBson := bson.M{"PlayerId": reqMsg.GetId(), "Status": bson.M{"$ne": for_game.COIN_BAG_ITEM_EXPIRED}}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["PropsType"] = reqMsg.GetListType()
	}

	if reqMsg.Keyword != nil {
		switch reqMsg.GetType() {
		case 1:
			id := easygo.StringToInt64noErr(reqMsg.GetKeyword())
			queryBson["PropsId"] = id
		case 2:
			findBson["PropsName"] = reqMsg.GetKeyword()
		}
	}

	m := []bson.M{
		{"$match": queryBson},
		{"$lookup": bson.M{
			"from":         for_game.TABLE_PROPS_ITEM,
			"localField":   "PropsId",
			"foreignField": "_id",
			"as":           "Props",
		}},
		{"$project": bson.M{"PropsId": 1, "GetType": 1, "OverTime": 1, "PropsName": "$Props.Name", "PropsType": "$Props.PropsType"}},
		{"$unwind": "$PropsName"},
		{"$unwind": "$PropsType"},
		{"$match": findBson},
		{"$sort": bson.M{"OverTime": -1}},
		// {"$skip": curPage * pageSize},
		// {"$limit": pageSize},
	}

	ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BAGITEM, m)
	count := len(ls)

	m = append(m, bson.M{"$skip": curPage * pageSize}, bson.M{"$limit": pageSize})
	lis := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BAGITEM, m)
	var list []*share_message.PlayerBagItem
	for _, li := range lis {
		one := &share_message.PlayerBagItem{}
		for_game.StructToOtherStruct(li, one)
		// one.PropsName = easygo.NewString(li.(bson.M)["PropsName"])
		// one.PropsType = easygo.NewInt32(li.(bson.M)["PropsType"])
		list = append(list, one)
	}
	msg := &brower_backstage.PlayerBagItemResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}*/

func (self *cls4) RpcPlayerBagItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	if reqMsg.Id != nil && reqMsg.GetId() > 0 {
		bagObj := for_game.GetRedisPlayerBagItemObj(reqMsg.GetId())
		bagObj.SaveToMongo()
	}

	pageSize := int(reqMsg.GetPageSize())
	curPage := int(reqMsg.GetCurPage())

	queryBson := bson.M{"PlayerId": reqMsg.GetId(), "Status": bson.M{"$ne": for_game.COIN_BAG_ITEM_EXPIRED}}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		queryBson["PropsType"] = reqMsg.GetListType()
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			id := easygo.StringToInt64noErr(reqMsg.GetKeyword())
			queryBson["PropsId"] = id
		case 2:
			queryBson["PropsName"] = reqMsg.GetKeyword()
		}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BAGITEM, queryBson, pageSize, curPage, "-CreateTime")
	var list []*share_message.PlayerBagItem
	for _, li := range lis {
		one := &share_message.PlayerBagItem{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}

	msg := &brower_backstage.PlayerBagItemResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//用户道具获得日志列表
func (self *cls4) RpcPlayerGetPropsLogList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	queryBson := bson.M{}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["PropsType"] = reqMsg.GetListType()
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			id := easygo.StringToInt64noErr(reqMsg.GetKeyword())
			queryBson["PropsId"] = id
		}
	}

	if reqMsg.SrtType != nil && reqMsg.GetSrtType() != "" {
		switch reqMsg.GetDownType() {
		case 1:
			id := QueryPlayerbyPhone(reqMsg.GetSrtType()).GetPlayerId()
			queryBson["PlayerId"] = id
		case 2:
			id := QueryPlayerbyAccount(reqMsg.GetSrtType()).GetPlayerId()
			queryBson["PlayerId"] = id
		}
	}

	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		start := reqMsg.GetBeginTimestamp()
		end := reqMsg.GetEndTimestamp()
		if start > 9999999999 {
			start /= 1000
		}
		if end > 9999999999 {
			end /= 1000
		}
		findBson["CreateTime"] = bson.M{"$gte": start, "$lte": end}
	}

	m := []bson.M{
		{"$match": queryBson},
		{"$lookup": bson.M{
			"from":         for_game.TABLE_PROPS_ITEM,
			"localField":   "PropsId",
			"foreignField": "_id",
			"as":           "Props",
		}},
		{"$project": bson.M{"_id": 1, "CreateTime": 1, "PlayerId": 1, "GivePlayerId": 1, "PropsId": 1, "PropsNum": 1, "GetType": 1, "EffectiveTime": 1, "BagId": 1, "RecycleTime": 1, "Operator": 1, "Note": 1, "BuyWay": 1, "OrderId": 1, "PropsName": "$Props.Name", "PropsType": "$Props.PropsType"}},
		{"$unwind": "$PropsName"},
		{"$unwind": "$PropsType"},
		{"$match": findBson},
		{"$sort": bson.M{"CreateTime": -1}},
		// {"$skip": curPage * pageSize},
		// {"$limit": pageSize},
	}

	ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_GETPROPS_LOG, m, 0, 0)
	count := len(ls)

	lis := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_GETPROPS_LOG, m, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()))
	var list []*share_message.PlayerGetPropsLog
	for _, li := range lis {
		one := &share_message.PlayerGetPropsLog{}
		for_game.StructToOtherStruct(li, one)
		one.PropsName = easygo.NewString(li.(bson.M)["PropsName"])
		one.PropsType = easygo.NewInt32(li.(bson.M)["PropsType"])
		if one.GetPlayerId() > 0 {
			p1 := QueryPlayerbyId(one.GetPlayerId())
			if p1 != nil {
				one.Account = easygo.NewString(p1.GetAccount())
			}
		}
		if one.GetGivePlayerId() > 0 {
			p2 := QueryPlayerbyId(one.GetGivePlayerId())
			if p2 != nil {
				one.Account = easygo.NewString(p2.GetAccount())
			}
		}

		list = append(list, one)
	}
	msg := &brower_backstage.PlayerGetPropsLogResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//日志回收用户道具
func (self *cls4) RpcRecycleProps(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	// logs.Debug("RpcRecycleProps:", reqMsg)
	logIds := reqMsg.GetIds64()
	logsCount := len(logIds)
	if logsCount == 0 {
		return easygo.NewFailMsg("请先选择操作项")
	}

	req := &server_server.PlayerIds{
		PlayerIds: logIds,
	}
	result := ChooseOneHall(0, "RpcRecyclePropsToHall", req)
	err := for_game.ParseReturnDataErr(result)
	if err != nil {
		return err
	} else {
		res := result.(*server_server.PlayerIds)
		okids := res.GetPlayerIds()
		for_game.UpdateAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_GETPROPS_LOG, bson.M{"_id": bson.M{"$in": okids}}, bson.M{"$set": bson.M{"RecycleTime": easygo.NowTimestamp(), "Operator": user.GetAccount(), "Note": reqMsg.GetNote()}})

		count := len(okids)
		var ids string
		for i := 0; i < count; i++ {
			if i < count {
				ids += easygo.IntToString(int(okids[i])) + ","
			} else {
				ids += easygo.IntToString(int(okids[i]))
			}
		}

		msg := fmt.Sprintf("批量回收道具日志Id: %s", ids)
		AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PROPS_MANAGE, msg)

		if logsCount != len(okids) {
			var errIds string
			for j, i := range logIds {
				if for_game.IsContains(i, okids) == -1 {
					if j < logsCount-1 {
						errIds += easygo.IntToString(int(i)) + ","
					} else {
						errIds += easygo.IntToString(int(i))
					}
				}
			}
			errmsg := fmt.Sprintf("道具%s回收失败,已过期", errIds)
			return easygo.NewFailMsg(errmsg)
		}
	}

	return easygo.EmptyMsg
}

//系统赠送商品内的道具给用户
func (self *cls4) RpcSysGiveProps(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Debug("RpcSysGiveProps", reqMsg)
	playerId := reqMsg.GetId64()
	if playerId == 0 {
		return easygo.NewFailMsg("赠送人不能为空")
	}
	shopId := reqMsg.GetObjId()
	if shopId == 0 {
		return easygo.NewFailMsg("赠送商品不能为空")
	}
	count := reqMsg.GetId32()
	if count == 0 {
		return easygo.NewFailMsg("赠送道具数量不能小于0")
	}

	req := &server_server.SysGivePropsRequest{
		PlayerId:  easygo.NewInt64(playerId),
		ProductId: easygo.NewInt64(shopId),
		Num:       easygo.NewInt64(count),
		Operator:  easygo.NewString(user.GetAccount()),
	}

	result := SendToPlayer(playerId, "RpcSysGivePropsToHall", req)
	err := for_game.ParseReturnDataErr(result)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("后台赠送[%d][%d]个道具商品[%d]", playerId, count, shopId)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PROPS_MANAGE, msg)

	return easygo.EmptyMsg
}

//用户背包回收用户道具
func (self *cls4) RpcBagRecycleProps(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	// logs.Debug("RpcRecycleProps:", reqMsg)
	id := reqMsg.GetId64()
	day := reqMsg.GetId32()
	playerId := reqMsg.GetObjId()
	if id == 0 {
		return easygo.NewFailMsg("id参数错误")
	}
	if day == 0 {
		return easygo.NewFailMsg("回收天数错误")
	}
	if playerId == 0 {
		return easygo.NewFailMsg("用户id不能为空")
	}

	reqMsg.IdStr = easygo.NewString(user.GetAccount())
	result := SendToPlayer(playerId, "RpcBagRecyclePropsToHall", reqMsg)
	err := for_game.ParseReturnDataErr(result)
	if err != nil {
		return err
	} else {
		one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BAGITEM, bson.M{"_id": id})
		player := QueryPlayerbyId(one.(bson.M)["PlayerId"].(int64))
		msg := fmt.Sprintf("系统回收[%s]的道具[%d]%d天", player.GetNickName(), one.(bson.M)["PropsId"], day)
		AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PROPS_MANAGE, msg)
	}

	return easygo.EmptyMsg
}
