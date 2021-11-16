package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/h5_wish"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"sort"
	"time"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

// 工具盲盒基础数据
func (self *cls4) RpcToolWishBoxList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_WISH_BOX, bson.M{"GuardianId": user.GetId()}, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()))

	var list []*brower_backstage.WishBox
	for _, li := range lis {
		one := &brower_backstage.WishBox{}
		for_game.StructToOtherStruct(li, one)
		one.Id = easygo.NewInt64(li.(bson.M)["_id"])
		list = append(list, one)
	}

	return &brower_backstage.WishBoxList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//工具添加盲盒
func (self *cls4) RpcToolWishBoxSave(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.WishBox) easygo.IMessage {
	if reqMsg.GetName() == "" {
		return easygo.NewFailMsg("盲盒名字不能为空")
	}

	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		id := for_game.NextId(for_game.TABLE_TOOL_WISH_BOX)
		reqMsg.Id = easygo.NewInt64(id)
	}

	reqMsg.GuardianId = easygo.NewInt64(user.GetId())

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_WISH_BOX, queryBson, updateBson, true)

	return easygo.EmptyMsg
}

// 商品基础数据
func (self *cls4) RpcToolWishBoxItemList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {

	lis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_WISH_BOX_ITEM, bson.M{"WishBoxId": reqMsg.GetId64()}, 0, 0)

	var list []*share_message.WishBoxItem
	for _, li := range lis {
		one := &share_message.WishBoxItem{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}

	return &brower_backstage.ToolWishBoxItemListRes{
		List: list,
	}
}

//保存商品基础数据
func (self *cls4) RpcToolSaveWishBoxItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ToolSaveWishBoxItemReq) easygo.IMessage {
	items := reqMsg.GetList()
	if len(reqMsg.GetList()) == 0 {
		return easygo.NewFailMsg("请先填写商品信息")
	}

	for _, it := range items {
		if it.GetWishItemName() == "" { //|| it.GetPrice() == 0 || it.GetRewardLv() == 0 || it.GetCommon() == 0 || it.GetSmallWin() == 0 || it.GetBigWin() == 0 || it.GetCommonAddWeight() == 0 || it.GetBigWinAddWeight() == 0 || it.GetSmallWinAddWeight() == 0 {
			continue
		}
		if it.GetWishBoxId() == 0 {
			return easygo.NewFailMsg("盲盒ID不能为空")
		}

		if it.Id == nil || it.GetId() == 0 {
			id := for_game.NextId(for_game.TABLE_TOOL_WISH_BOX_ITEM)
			it.Id = easygo.NewInt64(id)
			it.WishItemId = easygo.NewInt64(id)
		}

		queryBson := bson.M{"_id": it.GetId()}
		updateBson := bson.M{"$set": it}
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_WISH_BOX_ITEM, queryBson, updateBson, true)
	}

	return easygo.EmptyMsg
}

//删除商品基础数据
func (s *cls4) RpcToolDelWishBoxItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	ids := reqMsg.GetIds64()
	if len(ids) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}
	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_WISH_BOX_ITEM, bson.M{"_id": bson.M{"$in": ids}})

	return easygo.EmptyMsg
}

//概率参数计算表
func (self *cls4) RpcToolRateList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ToolRateReq) easygo.IMessage {
	price := reqMsg.GetChallengeRmb()    //挑战价格
	poolPat := reqMsg.GetPoolPat()       //水池浮动参数
	weightsPat := reqMsg.GetWeightsPat() //权重参数
	if reqMsg.GetChallengeDiamond() == 0 || price == 0 || poolPat == 0 || weightsPat == 0 {
		return easygo.NewFailMsg("缺少必填参数")
	}
	items, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_WISH_BOX_ITEM, bson.M{"WishBoxId": reqMsg.GetWishBoxId()}, 0, 0)
	if count == 0 {
		return nil
	}

	var ltCount int64 = 0    //小于挑战价格的物品数量
	var ltPriceSum int64 = 0 //小于挑战价格的物品价格总和

	shops := make(map[int64]*share_message.WishBoxItem)
	for _, it := range items {
		one := &share_message.WishBoxItem{}
		for_game.StructToOtherStruct(it, one)
		shops[one.GetId()] = one
		if one.GetPrice() < price {
			ltCount += 1
			ltPriceSum += one.GetPrice()
		}
	}

	poolSum := 0.00
	var list []*brower_backstage.ToolRate
	for _, s := range shops {
		valueSum := float64(s.GetPrice()) / ((float64(price*ltCount) - float64(ltPriceSum)) / float64(ltCount))
		pool := valueSum * float64(price)
		li := &brower_backstage.ToolRate{
			Id:       easygo.NewInt64(s.GetId()),
			Name:     easygo.NewString(s.GetWishItemName()),
			Price:    easygo.NewInt64(s.GetPrice()),
			ValueSum: easygo.NewFloat64(valueSum),
			Pool:     easygo.NewFloat64(pool),
		}
		list = append(list, li)
		poolSum += pool
	}

	poolSumFloat := poolSum * poolPat
	ptWeightsSum := 0.00
	ptAppendSum := 0.00
	xyWeightsSum := 0.00
	xyAppendSum := 0.00
	dyWeightsSum := 0.00
	dyAppendSum := 0.00

	for _, l := range list {
		weights := poolSumFloat / float64(l.GetPrice()) * float64(weightsPat)
		ptWeights := weights * (float64(shops[l.GetId()].GetCommon()) / 1000)
		ptAppend := weights * (float64(shops[l.GetId()].GetCommonAddWeight()) / 1000)
		xyWeights := weights * (float64(shops[l.GetId()].GetSmallWin()) / 1000)
		xyAppend := weights * (float64(shops[l.GetId()].GetSmallWinAddWeight()) / 1000)
		dyWeights := weights * (float64(shops[l.GetId()].GetBigWin()) / 1000)
		dyAppend := weights * (float64(shops[l.GetId()].GetBigWinAddWeight()) / 1000)

		l.Weights = easygo.NewFloat64(weights)
		l.PtWeights = easygo.NewFloat64(ptWeights)
		l.PtAppend = easygo.NewFloat64(ptAppend)
		l.XyWeights = easygo.NewFloat64(xyWeights)
		l.XyAppend = easygo.NewFloat64(xyAppend)
		l.DyWeights = easygo.NewFloat64(dyWeights)
		l.DyAppend = easygo.NewFloat64(dyAppend)

		ptWeightsSum += ptWeights
		ptAppendSum += ptAppend
		xyWeightsSum += xyWeights
		xyAppendSum += xyAppend
		dyWeightsSum += dyWeights
		dyAppendSum += dyAppend
	}

	minPtRate := 0.00
	maxPtRate := 0.00
	minPtAppendRate := 0.00
	maxPtAppendRate := 0.00
	minXyRate := 0.00
	maxXyRate := 0.00
	minXyAppendRate := 0.00
	maxXyAppendRate := 0.00
	minDyRate := 0.00
	maxDyRate := 0.00
	minDyAppendRate := 0.00
	maxDyAppendRate := 0.00

	for _, lis := range list {
		ptRate := lis.GetPtWeights() / ptWeightsSum
		ptAppendRate := lis.GetPtAppend() / ptAppendSum
		xyRate := lis.GetXyWeights() / xyWeightsSum
		xyAppendRate := lis.GetXyAppend() / xyAppendSum
		dyRate := lis.GetDyWeights() / dyWeightsSum
		dyAppendRate := lis.GetDyAppend() / dyAppendSum

		lis.PtRate = easygo.NewFloat64(ptRate)
		lis.PtAppendRate = easygo.NewFloat64(ptAppendRate)
		lis.XyRate = easygo.NewFloat64(xyRate)
		lis.XyAppendRate = easygo.NewFloat64(xyAppendRate)
		lis.DyRate = easygo.NewFloat64(dyRate)
		lis.DyAppendRate = easygo.NewFloat64(dyAppendRate)
		if lis.GetPrice() > price {
			minPtRate += ptRate
			minPtAppendRate += ptAppendRate
			minXyRate += xyRate
			minXyAppendRate += xyAppendRate
			minDyRate += dyRate
			minDyAppendRate += dyAppendRate
		} else {
			maxPtRate += ptRate
			maxPtAppendRate += ptAppendRate
			maxXyRate += xyRate
			maxXyAppendRate += xyAppendRate
			maxDyRate += dyRate
			maxDyAppendRate += dyAppendRate
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].GetId() < list[j].GetId()
	})

	minli := &brower_backstage.ToolRate{
		Name:         easygo.NewString("(商品价格>抽取价格)合计"),
		PtRate:       easygo.NewFloat64(minPtRate),
		PtAppendRate: easygo.NewFloat64(minPtAppendRate),
		XyRate:       easygo.NewFloat64(minXyRate),
		XyAppendRate: easygo.NewFloat64(minXyAppendRate),
		DyRate:       easygo.NewFloat64(minDyRate),
		DyAppendRate: easygo.NewFloat64(minDyAppendRate),
	}
	maxli := &brower_backstage.ToolRate{
		Name:         easygo.NewString("(商品价格<=抽取价格)合计"),
		PtRate:       easygo.NewFloat64(maxPtRate),
		PtAppendRate: easygo.NewFloat64(maxPtAppendRate),
		XyRate:       easygo.NewFloat64(maxXyRate),
		XyAppendRate: easygo.NewFloat64(maxXyAppendRate),
		DyRate:       easygo.NewFloat64(maxDyRate),
		DyAppendRate: easygo.NewFloat64(maxDyAppendRate),
	}
	list = append(list, minli, maxli)

	return &brower_backstage.ToolRateRes{
		List: list,
	}
}

//抽奖
func (self *cls4) RpcToolLucky(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ToolLuckyReq) easygo.IMessage {
	logs.Info("RpcToolLucky:", reqMsg)
	if reqMsg.GetPoolId() == 0 || reqMsg.GetRunTimes() == 0 || reqMsg.GetChallengeDiamond() == 0 {
		return easygo.NewFailMsg("缺少必填参数")
	}
	if reqMsg.GetWishBoxId() == 0 {
		return easygo.NewFailMsg("盲盒ID不能为空")
	}
	wishbox := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_WISH_BOX, bson.M{"_id": reqMsg.GetWishBoxId()})
	if wishbox == nil {
		return easygo.NewFailMsg("盲盒不存在")
	}

	poolcfgId := reqMsg.GetPoolId()
	if poolcfgId > 0 {
		poolcfg := GetWishPoolCfg(poolcfgId)
		if poolcfg == nil {
			return easygo.NewFailMsg("水池配置不存在")
		}
		pool := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL, bson.M{"PoolConfigId": poolcfgId})
		one := &brower_backstage.WishPool{}
		if pool != nil {
			for_game.StructToOtherStruct(pool, one)
		} else {
			boxPoolId := for_game.NextId(for_game.TABLE_WISH_POOL)
			one = &brower_backstage.WishPool{
				Id:               easygo.NewInt64(boxPoolId),
				PoolLimit:        easygo.NewInt32(poolcfg.GetPoolLimit()),
				Name:             easygo.NewString(poolcfg.GetName()),
				ShowInitialValue: easygo.NewInt32(poolcfg.GetShowInitialValue()),
				ShowRecycle:      easygo.NewInt32(poolcfg.GetShowRecycle()),
				ShowCommission:   easygo.NewInt32(poolcfg.GetShowCommission()),
				ShowStartAward:   easygo.NewInt32(poolcfg.GetShowStartAward()),
				ShowCloseAward:   easygo.NewInt32(poolcfg.GetShowCloseAward()),
				IsDefault:        easygo.NewBool(poolcfg.GetIsDefault()),
				LocalStatus:      easygo.NewInt64(0),
				IsOpenAward:      easygo.NewBool(false),
			}
		}

		one.PoolCfgId = easygo.NewInt64(poolcfgId)

		one.SmallLoss = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(poolcfg.SmallLoss.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(poolcfg.SmallLoss.GetShowMinValue()),
		}

		one.SmallWin = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(poolcfg.SmallWin.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(poolcfg.SmallWin.GetShowMinValue()),
		}

		one.BigLoss = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(poolcfg.BigLoss.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(poolcfg.BigLoss.GetShowMinValue()),
		}

		one.BigWin = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(poolcfg.BigWin.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(poolcfg.BigWin.GetShowMinValue()),
		}
		one.Common = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(poolcfg.Common.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(poolcfg.Common.GetShowMinValue()),
		}

		UpdateWishPool2(one)
		// 前端传过来的是，水池配置id,存储的是盲盒水池id
		reqMsg.PoolId = easygo.NewInt64(one.GetId())
	}

	req := &h5_wish.BackstageDareToolReq{
		PoolId:    easygo.NewInt64(reqMsg.GetPoolId()),
		Count:     easygo.NewInt64(reqMsg.GetRunTimes()),
		Diamond:   easygo.NewInt64(reqMsg.GetChallengeDiamond()),
		UserId:    easygo.NewInt64(user.GetId()),
		WishBoxId: easygo.NewInt64(reqMsg.GetWishBoxId()),
	}
	serInfo := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)

	easygo.Spawn(func() {
		httpErr := ChooseOneWish(serInfo.GetSid(), "RpcBackstageDareTool", req)
		err := for_game.ParseReturnDataErr(httpErr)
		if err != nil {
			logs.Warn("通知许愿池服务器更新钻石失败，rpc=RpcBackstageAddDiamond, req=%+v", req)
		}
	})

	return &brower_backstage.ToolLuckyRes{
		Result: easygo.NewString("抽奖正在进行"),
	}
}

// 重置水池Id64-id： 水池id
func (self *cls4) RpcToolResetWishPool(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("RpcToolResetWishPool, %+v：", reqMsg)

	pool := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL, bson.M{"PoolConfigId": reqMsg.GetId64()})
	if pool == nil {
		return easygo.NewFailMsg("水池不存在")
	}
	redis := GetPoolObj(pool.(bson.M)["_id"].(int64))
	v := redis.GetPollInfoFromRedis()
	one := &brower_backstage.WishPool{
		Id:               easygo.NewInt64(v.GetId()),
		PoolLimit:        easygo.NewInt32(v.GetPoolLimit()),
		Name:             easygo.NewString(v.GetName()),
		ShowInitialValue: easygo.NewInt32(v.GetShowInitialValue()),
		ShowRecycle:      easygo.NewInt32(v.GetShowRecycle()),
		ShowCommission:   easygo.NewInt32(v.GetShowCommission()),
		ShowStartAward:   easygo.NewInt32(v.GetShowStartAward()),
		ShowCloseAward:   easygo.NewInt32(v.GetShowCloseAward()),
		IsDefault:        easygo.NewBool(v.GetIsDefault()),
		LocalStatus:      easygo.NewInt64(0),
		IsOpenAward:      easygo.NewBool(false),
	}

	one.PoolCfgId = easygo.NewInt64(v.GetPoolConfigId())

	one.SmallLoss = &brower_backstage.WishPoolStatus{
		ShowMaxValue: easygo.NewInt32(v.SmallLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.SmallLoss.GetShowMinValue()),
	}

	one.SmallWin = &brower_backstage.WishPoolStatus{
		ShowMaxValue: easygo.NewInt32(v.SmallWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.SmallWin.GetShowMinValue()),
	}

	one.BigLoss = &brower_backstage.WishPoolStatus{
		ShowMaxValue: easygo.NewInt32(v.BigLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.BigLoss.GetShowMinValue()),
	}

	one.BigWin = &brower_backstage.WishPoolStatus{
		ShowMaxValue: easygo.NewInt32(v.BigWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.BigWin.GetShowMinValue()),
	}
	one.Common = &brower_backstage.WishPoolStatus{
		ShowMaxValue: easygo.NewInt32(v.Common.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.Common.GetShowMinValue()),
	}

	UpdateWishPool2(one)
	redis.DeleteRedisPollInfo()
	DeleteToolWishPoolLog(one.GetId())
	DeleteToolWishPoolPumpLog(one.GetId())
	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_PLAYER_WISH_ITEM, bson.M{"WishBoxId": reqMsg.GetObjId()})

	redis.RefreshRedisPollInfo()
	msg := fmt.Sprintf("重置水池,水池id:%d", reqMsg.GetId64())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 获取水池Id64-id： 水池id
func (self *cls4) RpcToolGetWishPool(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("RpcGetWishPool, %+v：", reqMsg)
	pool := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL, bson.M{"PoolConfigId": reqMsg.GetId64()})
	if pool == nil {
		return &brower_backstage.WishPool{}
	}
	redis := GetPoolObj(pool.(bson.M)["_id"].(int64))
	v := redis.GetPollInfoFromRedis()
	if v == nil && reqMsg.GetIdStr() != "tool" {
		return easygo.NewFailMsg("水池不存在")
	}
	one := &brower_backstage.WishPool{
		Id:               easygo.NewInt64(v.GetId()),
		PoolLimit:        easygo.NewInt32(v.GetPoolLimit()),
		InitialValue:     easygo.NewInt32(v.GetInitialValue()),
		IncomeValue:      easygo.NewInt32(v.GetIncomeValue()),
		Name:             easygo.NewString(v.GetName()),
		CreateTime:       easygo.NewInt64(v.GetCreateTime()),
		Recycle:          easygo.NewInt32(v.GetRecycle()),
		Commission:       easygo.NewInt32(v.GetCommission()),
		StartAward:       easygo.NewInt32(v.GetStartAward()),
		CloseAward:       easygo.NewInt32(v.GetCloseAward()),
		ShowInitialValue: easygo.NewInt32(v.GetShowInitialValue()),
		ShowRecycle:      easygo.NewInt32(v.GetShowRecycle()),
		ShowCommission:   easygo.NewInt32(v.GetShowCommission()),
		ShowStartAward:   easygo.NewInt32(v.GetShowStartAward()),
		ShowCloseAward:   easygo.NewInt32(v.GetShowCloseAward()),
		IsDefault:        easygo.NewBool(v.GetIsDefault()),
		IsOpenAward:      easygo.NewBool(v.GetIsOpenAward()),
		LocalStatus:      easygo.NewInt64(v.GetLocalStatus()),
	}

	one.SmallLoss = &brower_backstage.WishPoolStatus{
		MaxValue:     easygo.NewInt32(v.SmallLoss.GetMaxValue()),
		MinValue:     easygo.NewInt32(v.SmallLoss.GetMinValue()),
		ShowMaxValue: easygo.NewInt32(v.SmallLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.SmallLoss.GetShowMinValue()),
	}

	one.SmallWin = &brower_backstage.WishPoolStatus{
		MaxValue:     easygo.NewInt32(v.SmallWin.GetMaxValue()),
		MinValue:     easygo.NewInt32(v.SmallWin.GetMinValue()),
		ShowMaxValue: easygo.NewInt32(v.SmallWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.SmallWin.GetShowMinValue()),
	}

	one.BigLoss = &brower_backstage.WishPoolStatus{
		MaxValue:     easygo.NewInt32(v.BigLoss.GetMaxValue()),
		MinValue:     easygo.NewInt32(v.BigLoss.GetMinValue()),
		ShowMaxValue: easygo.NewInt32(v.BigLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.BigLoss.GetShowMinValue()),
	}

	one.BigWin = &brower_backstage.WishPoolStatus{
		MaxValue:     easygo.NewInt32(v.BigWin.GetMaxValue()),
		MinValue:     easygo.NewInt32(v.BigWin.GetMinValue()),
		ShowMaxValue: easygo.NewInt32(v.BigWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.BigWin.GetShowMinValue()),
	}
	one.Common = &brower_backstage.WishPoolStatus{
		MaxValue:     easygo.NewInt32(v.Common.GetMaxValue()),
		MinValue:     easygo.NewInt32(v.Common.GetMinValue()),
		ShowMaxValue: easygo.NewInt32(v.Common.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.Common.GetShowMinValue()),
	}

	return one
}

//输出数据
func (self *cls4) RpcToolOutputData(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("RpcToolOutputData, %+v：", reqMsg)
	groupBson := bson.M{"_id": nil}
	groupBson["DareDiamond"] = bson.M{"$sum": "$DareDiamond"}
	groupBson["WishItemDiamond"] = bson.M{"$sum": "$WishItemDiamond"}
	m := []bson.M{
		{"$match": bson.M{"WishBoxId": reqMsg.GetId64()}},
		{"$group": groupBson},
	}
	one := for_game.FindPipeOne(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_PLAYER_WISH_ITEM, m)
	result := &brower_backstage.ToolOutputDataRes{}
	if one != nil {
		for_game.StructToOtherStruct(one, result)
	}
	return result
}

//产出物品
func (self *cls4) RpcToolOutputitemList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("RpcToolOutputitemList, %+v：", reqMsg)
	lis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_PLAYER_WISH_ITEM, bson.M{"WishBoxId": reqMsg.GetId64()}, 0, 0)

	listMap := make(map[int64]*brower_backstage.ToolOutputitem)
	for _, li := range lis {
		one := &share_message.PlayerWishItem{}
		for_game.StructToOtherStruct(li, one)
		shopId := one.GetChallengeItemId()
		if listMap[shopId] == nil {
			listMap[shopId] = &brower_backstage.ToolOutputitem{}
		}
		listMap[shopId].Id = easygo.NewInt64(shopId)
		listMap[shopId].Name = easygo.NewString(one.GetProductName())
		if one.GetAfterIsOpenAward() {
			switch one.GetAfterLocalStatus() {
			//1-大亏;2-小亏;3-普通;4-大赢;5-小盈
			case 3:
				listMap[shopId].PtO = easygo.NewInt64(listMap[shopId].GetPtO() + int64(1))
			case 4:
				listMap[shopId].DyO = easygo.NewInt64(listMap[shopId].GetDyO() + int64(1))
			case 5:
				listMap[shopId].XyO = easygo.NewInt64(listMap[shopId].GetXyO() + int64(1))
			}
		} else {
			switch one.GetAfterLocalStatus() {
			//1-大亏;2-小亏;3-普通;4-大赢;5-小盈
			case 3:
				listMap[shopId].PtC = easygo.NewInt64(listMap[shopId].GetPtC() + int64(1))
			case 4:
				listMap[shopId].DyC = easygo.NewInt64(listMap[shopId].GetDyC() + int64(1))
			case 5:
				listMap[shopId].XyC = easygo.NewInt64(listMap[shopId].GetXyC() + int64(1))
			}
		}
	}
	var list []*brower_backstage.ToolOutputitem
	for _, lm := range listMap {
		list = append(list, lm)
	}
	return &brower_backstage.ToolOutputitemRes{
		List: list,
	}
}

//抽水表
func (self *cls4) RpcToolToolPumping(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("RpcToolToolPumping, %+v：", reqMsg)
	pool := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL, bson.M{"PoolConfigId": reqMsg.GetObjId()})
	if pool == nil {
		return easygo.NewFailMsg("水池不存在")
	}
	groupBson := bson.M{"_id": nil}
	groupBson["PumpingTimes"] = bson.M{"$sum": 1}
	groupBson["PumpingSumDiamond"] = bson.M{"$sum": "$Price"}
	groupBson["PumpingAvgDiamond"] = bson.M{"$avg": "$Price"}
	m := []bson.M{
		{"$match": bson.M{"BoxId": reqMsg.GetId64(), "PoolId": pool.(bson.M)["_id"].(int64)}},
		{"$group": groupBson},
	}
	one := for_game.FindPipeOne(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_WISH_POOL_PUMP_LOG, m)
	result := &brower_backstage.ToolPumping{}
	if one != nil {
		result.PumpingTimes = easygo.NewInt64(one.(bson.M)["PumpingTimes"])
		result.PumpingSumDiamond = easygo.NewInt64(one.(bson.M)["PumpingSumDiamond"])
		result.PumpingAvgDiamond = easygo.NewInt64(one.(bson.M)["PumpingAvgDiamond"])
	}
	return result
}

// 获取盲盒列表
func (self *cls4) RpcQueryWishBoxList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.WishBoxListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishBoxList, %+v：", reqMsg)
	list, count := QueryWishBoxList(reqMsg)

	for i := range list {
		goods, _ := GetWishBoxGoodsItemList(&brower_backstage.ListRequest{Id: easygo.NewInt64(list[i].GetId())})
		list[i].ItemList = goods
		redis := GetPoolObj(list[i].GetWishPoolId())
		v := redis.GetPollInfoFromRedis()
		list[i].LocalStatus = easygo.NewInt64(v.GetLocalStatus())
		// 返回给前端传的是水池配置id,从数据库拿的是盲盒水池id
		list[i].WishPoolId = easygo.NewInt64(v.GetPoolConfigId())
		// 盲盒水池id
		list[i].BoxPoolId = easygo.NewInt64(v.GetId())
		list[i].IsOpenAward = easygo.NewBool(v.GetIsOpenAward())
	}

	ret := &brower_backstage.WishBoxList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 新增/更新盲盒
func (self *cls4) RpcUpdateWishBox(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishBox) easygo.IMessage {
	logs.Info("RpcUpdateWishBox, %+v：", reqMsg)
	msg := fmt.Sprintf("更新盲盒" + reqMsg.GetName())
	if reqMsg.GetId() == 0 {
		msg = fmt.Sprintf("新增盲盒" + reqMsg.GetName())
	}

	// 检查守护者账号
	imPlayer := &share_message.PlayerBase{}
	if reqMsg.GetUserAccount() != "" {
		imPlayer = for_game.GetPlayerByAccount(reqMsg.GetUserAccount())
		if imPlayer == nil {
			return easygo.NewFailMsg("守护者账号不存在")
		}

		wishUser := GetWishPlayerInfoByPid(imPlayer.GetPlayerId())
		if wishUser.GetPlayerId() == 0 {
			return easygo.NewFailMsg("守护者账号不存在")
		}
		reqMsg.UserId = easygo.NewInt64(wishUser.GetId())
	}

	var backGuardReq *h5_wish.BackstageSetGuardianReq
	curTime := time.Now().Unix()
	var guardId int64 = 0              //守护者Id
	var ifChangedGuardian bool = false //是否变更了守护者
	if reqMsg.GetId() != 0 {
		//编辑
		box := GetWishBoxById(reqMsg.GetId())
		if box != nil {

			// 需要上架时间
			if box.GetStatus() != reqMsg.GetStatus() {
				if reqMsg.GetStatus() == 1 {
					reqMsg.UploadTime = easygo.NewInt64(curTime)
				} else {
					upBox := &share_message.WishBox{
						Id:     easygo.NewInt64(box.GetId()),
						Status: easygo.NewInt32(0),
					}
					upBox.UpdateTime = easygo.NewInt64(curTime)
					upBox.PutOnTime = easygo.NewInt64(0)
					err2 := for_game.UpWishBox(upBox.GetId(), upBox)
					easygo.PanicError(err2)
					return easygo.EmptyMsg
				}

			}

			if box.GetName() != reqMsg.GetName() {
				// 盲盒名称是唯一的
				box := GetWishBoxByName(reqMsg.GetName())
				if box != nil {
					msg := fmt.Sprintf("已存在“%s”盲盒", reqMsg.GetName())
					return easygo.NewFailMsg(msg)
				}
			}

			// 更新守护者信息
			//if box.GetGuardianId() != imPlayer.GetUserId() {
			if box.GetGuardianId() != imPlayer.GetPlayerId() {
				ifChangedGuardian = true
				//serInfo := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)
				req := &h5_wish.BackstageSetGuardianReq{
					Account:  easygo.NewString(imPlayer.GetPhone()),
					Channel:  easygo.NewInt32(1001),
					NickName: easygo.NewString(imPlayer.GetNickName()),
					HeadUrl:  easygo.NewString(imPlayer.GetHeadIcon()),
					PlayerId: easygo.NewInt64(imPlayer.GetPlayerId()),
					Token:    easygo.NewString(""),
					BoxId:    easygo.NewInt64(box.GetId()),
				}

				if reqMsg.GetUserAccount() == "" {
					req.OpType = easygo.NewInt32(2)
				} else {
					req.OpType = easygo.NewInt32(1)
				}
				backGuardReq = req
			}

			// 更新盲盒水池
			if reqMsg.GetWishPoolId() > 0 {
				boxPool := GetWishPool(box.GetWishPoolId())
				pool := GetWishPoolCfg(reqMsg.GetWishPoolId())
				if boxPool.GetPoolConfigId() != reqMsg.GetWishPoolId() {

					boxPoolId := for_game.NextId(for_game.TABLE_WISH_POOL)
					one := &brower_backstage.WishPool{
						Id:               easygo.NewInt64(boxPoolId),
						PoolLimit:        easygo.NewInt32(pool.GetPoolLimit()),
						Name:             easygo.NewString(pool.GetName()),
						ShowInitialValue: easygo.NewInt32(pool.GetShowInitialValue()),
						ShowRecycle:      easygo.NewInt32(pool.GetShowRecycle()),
						ShowCommission:   easygo.NewInt32(pool.GetShowCommission()),
						ShowStartAward:   easygo.NewInt32(pool.GetShowStartAward()),
						ShowCloseAward:   easygo.NewInt32(pool.GetShowCloseAward()),
						IsDefault:        easygo.NewBool(pool.GetIsDefault()),
						LocalStatus:      easygo.NewInt64(0),
						IsOpenAward:      easygo.NewBool(false),
					}

					one.PoolCfgId = easygo.NewInt64(reqMsg.GetWishPoolId())

					one.SmallLoss = &brower_backstage.WishPoolStatus{
						ShowMaxValue: easygo.NewInt32(pool.SmallLoss.GetShowMaxValue()),
						ShowMinValue: easygo.NewInt32(pool.SmallLoss.GetShowMinValue()),
					}

					one.SmallWin = &brower_backstage.WishPoolStatus{
						ShowMaxValue: easygo.NewInt32(pool.SmallWin.GetShowMaxValue()),
						ShowMinValue: easygo.NewInt32(pool.SmallWin.GetShowMinValue()),
					}

					one.BigLoss = &brower_backstage.WishPoolStatus{
						ShowMaxValue: easygo.NewInt32(pool.BigLoss.GetShowMaxValue()),
						ShowMinValue: easygo.NewInt32(pool.BigLoss.GetShowMinValue()),
					}

					one.BigWin = &brower_backstage.WishPoolStatus{
						ShowMaxValue: easygo.NewInt32(pool.BigWin.GetShowMaxValue()),
						ShowMinValue: easygo.NewInt32(pool.BigWin.GetShowMinValue()),
					}
					one.Common = &brower_backstage.WishPoolStatus{
						ShowMaxValue: easygo.NewInt32(pool.Common.GetShowMaxValue()),
						ShowMinValue: easygo.NewInt32(pool.Common.GetShowMinValue()),
					}

					UpdateWishPool2(one)
					logs.Info("盲盒更换新的盲盒水池，info=%+v", one)
					// 前端传过来的是，水池配置id,存储的是盲盒水池id
					reqMsg.WishPoolId = easygo.NewInt64(boxPoolId)
				} else {
					reqMsg.WishPoolId = nil
				}

			}
			//编辑时不改变盲盒守护者由CallRpcBackstageSetGuardian变更
			guardId = box.GetGuardianId()
		}
	} else {
		//新增盲盒
		// 盲盒名称是唯一的
		box := GetWishBoxByName(reqMsg.GetName())
		if box != nil {
			msg := fmt.Sprintf("已存在“%s”盲盒", reqMsg.GetName())
			return easygo.NewFailMsg(msg)
		}

		// 上架
		if reqMsg.GetStatus() == 1 {
			reqMsg.UploadTime = easygo.NewInt64(curTime)
		}

		//新增时添加了守护者
		if reqMsg.GetUserAccount() != "" {
			ifChangedGuardian = true
			req := &h5_wish.BackstageSetGuardianReq{
				Account:  easygo.NewString(imPlayer.GetPhone()),
				Channel:  easygo.NewInt32(1001),
				NickName: easygo.NewString(imPlayer.GetNickName()),
				HeadUrl:  easygo.NewString(imPlayer.GetHeadIcon()),
				PlayerId: easygo.NewInt64(imPlayer.GetPlayerId()),
				Token:    easygo.NewString(""),
				BoxId:    easygo.NewInt64(box.GetId()),
				OpType:   easygo.NewInt32(1),
			}
			backGuardReq = req
		}

		// 更新盲盒水池
		if reqMsg.GetWishPoolId() > 0 {
			pool := GetWishPoolCfg(reqMsg.GetWishPoolId())
			boxPoolId := for_game.NextId(for_game.TABLE_WISH_POOL)
			one := &brower_backstage.WishPool{
				Id:               easygo.NewInt64(boxPoolId),
				PoolLimit:        easygo.NewInt32(pool.GetPoolLimit()),
				Name:             easygo.NewString(pool.GetName()),
				ShowInitialValue: easygo.NewInt32(pool.GetShowInitialValue()),
				ShowRecycle:      easygo.NewInt32(pool.GetShowRecycle()),
				ShowCommission:   easygo.NewInt32(pool.GetShowCommission()),
				ShowStartAward:   easygo.NewInt32(pool.GetShowStartAward()),
				ShowCloseAward:   easygo.NewInt32(pool.GetShowCloseAward()),
				IsDefault:        easygo.NewBool(pool.GetIsDefault()),
				LocalStatus:      easygo.NewInt64(0),
				IsOpenAward:      easygo.NewBool(false),
			}

			one.PoolCfgId = easygo.NewInt64(reqMsg.GetWishPoolId())

			one.SmallLoss = &brower_backstage.WishPoolStatus{
				ShowMaxValue: easygo.NewInt32(pool.SmallLoss.GetShowMaxValue()),
				ShowMinValue: easygo.NewInt32(pool.SmallLoss.GetShowMinValue()),
			}

			one.SmallWin = &brower_backstage.WishPoolStatus{
				ShowMaxValue: easygo.NewInt32(pool.SmallWin.GetShowMaxValue()),
				ShowMinValue: easygo.NewInt32(pool.SmallWin.GetShowMinValue()),
			}

			one.BigLoss = &brower_backstage.WishPoolStatus{
				ShowMaxValue: easygo.NewInt32(pool.BigLoss.GetShowMaxValue()),
				ShowMinValue: easygo.NewInt32(pool.BigLoss.GetShowMinValue()),
			}

			one.BigWin = &brower_backstage.WishPoolStatus{
				ShowMaxValue: easygo.NewInt32(pool.BigWin.GetShowMaxValue()),
				ShowMinValue: easygo.NewInt32(pool.BigWin.GetShowMinValue()),
			}
			one.Common = &brower_backstage.WishPoolStatus{
				ShowMaxValue: easygo.NewInt32(pool.Common.GetShowMaxValue()),
				ShowMinValue: easygo.NewInt32(pool.Common.GetShowMinValue()),
			}

			UpdateWishPool2(one)
			// 前端传过来的是，水池配置id,存储的是盲盒水池id
			reqMsg.WishPoolId = easygo.NewInt64(boxPoolId)
		}
	}

	itemList := reqMsg.GetItemList()

	// 盲盒商品数量>0
	if reqMsg.GetStatus() == 1 && len(itemList) < 1 {
		return easygo.NewFailMsg("需要盲盒商品数量>0")
	}

	var (
		delItemIds  []int64 // 被删除的盲盒商品id
		itemIds     []int64 // 盲盒商品id
		wishItemIds []int64 // 商品id
		brandIds    []int64 // 品牌id
		typeIds     []int64 // 类型id
		stylesIds   []int32 // 款式id
		rareNum     int32

		isNew bool
		isWin bool // 盲盒是否包含必中商品
	)

	brandM := make(map[int64]struct{})  // 品牌id map
	typeM := make(map[int64]struct{})   // 类型id map
	stylesM := make(map[int32]struct{}) // 款式id map
	// 旧商品价格map
	itemMp := make(map[int64]int64)
	itemIdMp := make(map[int]int64)
	delItemM := make(map[int64]int64) // 被删除的盲盒商品id

	if reqMsg.GetId() != 0 {
		// 被删除的盲盒商品id
		itemReq := &brower_backstage.ListRequest{}
		itemReq.Id = easygo.NewInt64(reqMsg.GetId())
		list, _ := GetWishBoxGoodsItemList(itemReq)
		for _, v := range list {
			delItemM[v.GetId()] = v.GetId()
		}

	}

	goodsIds := []int64{}
	for i, v := range itemList {
		item := for_game.GetWishItemByIdFromDB(v.GetGoodsId())
		if item.GetStatus() == 0 {
			goodsIds = append(goodsIds, v.GetGoodsId())
		}
		itemMp[v.GetGoodsId()] = item.GetPrice()
		brandM[int64(item.GetBrand())] = struct{}{}
		typeM[int64(item.GetType())] = struct{}{}
		v.ArrivalTime = easygo.NewInt64(item.GetPreHaveTime())

		if item.GetIsPreSale() {
			reqMsg.ProductStatus = easygo.NewInt32(2)
		}

		itemIdMp[i] = v.GetBoxItemId()
	}

	if len(goodsIds) > 0 {
		idstr := ""
		for i, j := range goodsIds {
			idstr += easygo.AnytoA(j)
			if i < len(goodsIds)-1 {
				idstr += ","
			}
		}
		return easygo.NewFailMsg(fmt.Sprintf("盲盒商品有变化，商品已下架。商品ID:[%s]", idstr))
	}

	if reqMsg.GetId() == 0 {
		isNew = true
		bid := for_game.NextId(for_game.TABLE_WISH_BOX)
		reqMsg.Id = easygo.NewInt64(bid)
		reqMsg.CreateTime = easygo.NewInt64(curTime)
	} else {
		DeleteWishBoxItem(reqMsg.GetId())
		reqMsg.UpdateTime = easygo.NewInt64(curTime)
	}

	// 设置盲盒商品列表
	itemCount := len(itemList)
	for i, v := range itemList {

		if len(itemMp) != 0 {
			if itemMp[v.GetGoodsId()] != v.GetPrice() {
				return easygo.NewFailMsg(fmt.Sprintf("盲盒商品价格发生变化，请重新配置中奖配置。商品ID:%d", v.GetGoodsId()))
			}
		}

		if v.GetPrice() < 0 || v.Price == nil {
			return easygo.NewFailMsg("数值错误")
		}
		if v.GetBigLoss() < 0 || v.Price == nil {
			return easygo.NewFailMsg("中奖配置未配置")
		}
		if v.GetSmallLoss() < 0 || v.Price == nil {
			return easygo.NewFailMsg("中奖配置未配置")
		}
		if v.GetCommon() < 0 || v.Price == nil {
			return easygo.NewFailMsg("中奖配置未配置")
		}
		if v.GetBigWin() < 0 || v.Price == nil {
			return easygo.NewFailMsg("中奖配置未配置")
		}
		if v.GetSmallWin() < 0 || v.Price == nil {
			return easygo.NewFailMsg("中奖配置未配置")
		}
		if v.GetReplenishAmount() < 0 {
			return easygo.NewFailMsg("单次补货量下限为1")
		}
		if v.ReplenishAmount == nil {
			v.ReplenishAmount = easygo.NewInt32(0)
		}

		// 保存被删除的物品重新添加后，id不变
		v.Id = easygo.NewInt64(itemIdMp[i])
		if v.GetId() == 0 {
			v.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_BOX_ITEM))
		}

		data := &brower_backstage.WishBoxGoodsWin{
			Id:                    easygo.NewInt64(v.GetId()),
			GoodsId:               easygo.NewInt64(v.GetGoodsId()),
			Name:                  easygo.NewString(v.GetName()),
			Price:                 easygo.NewInt64(v.GetPrice()),
			IsInfallible:          easygo.NewBool(v.GetIsInfallible()),
			ReplenishAmount:       easygo.NewInt32(v.GetReplenishAmount()),
			ReplenishIntervalTime: easygo.NewInt32(v.GetReplenishIntervalTime()),
			PerRate:               easygo.NewInt32(v.GetPerRate()),
			GoodsType:             easygo.NewInt32(v.GetGoodsType()),
			WishBoxId:             easygo.NewInt64(reqMsg.GetId()),
			BigLoss:               easygo.NewInt32(v.GetBigLoss()),
			SmallLoss:             easygo.NewInt32(v.GetSmallLoss()),
			Common:                easygo.NewInt32(v.GetCommon()),
			BigWin:                easygo.NewInt32(v.GetBigWin()),
			SmallWin:              easygo.NewInt32(v.GetSmallWin()),
			RewardLv:              easygo.NewInt32(v.GetRewardLv()),
			CommonAddWeight:       easygo.NewInt32(v.GetCommonAddWeight()),
			BigWinAddWeight:       easygo.NewInt32(v.GetBigWinAddWeight()),
			SmallWinAddWeight:     easygo.NewInt32(v.GetSmallWinAddWeight()),
			Diamond:               easygo.NewInt64(v.GetDiamond()),
			ArrivalTime:           easygo.NewInt64(v.GetArrivalTime()),
		}

		if v.GetIsInfallible() {
			isWin = true
		}

		data.CreateTime = easygo.NewInt64(curTime + int64(itemCount-i))
		err := UpdateWishBoxGoodsItem(data)
		if err != nil {
			return easygo.NewFailMsg(err.Error())
		}
		itemIds = append(itemIds, v.GetId())
		wishItemIds = append(wishItemIds, v.GetGoodsId())
		stylesM[v.GetGoodsType()] = struct{}{}

		if v.GetGoodsType() != 1 {
			rareNum++
		}

		delItemM[v.GetId()] = 0
	}

	for k, _ := range brandM {
		brandIds = append(brandIds, k)
	}

	for k, _ := range typeM {
		typeIds = append(typeIds, k)
	}

	for k, _ := range stylesM {
		stylesIds = append(stylesIds, k)
	}

	for k, v := range delItemM {
		if v != 0 {
			delItemIds = append(delItemIds, k)
		}
	}

	reqMsg.RareNum = easygo.NewInt32(rareNum)
	reqMsg.Brands = brandIds
	reqMsg.Types = typeIds
	reqMsg.Styles = stylesIds
	reqMsg.Items = itemIds
	reqMsg.WishItems = wishItemIds
	reqMsg.GoodsAmount = easygo.NewInt32(len(itemIds))
	reqMsg.HaveIsWin = easygo.NewBool(isWin)

	reqMsg.UserId = easygo.NewInt64(guardId)
	err := UpdateWishBox(reqMsg)
	if err != nil {
		logs.Error("UpdateWishBox err: ", err)
		return easygo.NewFailMsg("操作盲盒失败！")
	}
	if ifChangedGuardian {
		//变更了守护者
		err := CallRpcBackstageSetGuardian(backGuardReq)
		if err != nil {
			logs.Warn("通知许愿池服务器更新盲盒守护者失败，rpc=RpcBackstageSetGuardian，req=%+v, err=%v", backGuardReq, err.GetReason())
			return easygo.NewFailMsg("更新盲盒守护者失败")
		}
	}

	if !isNew {
		// 删除盲盒商品相关用户许愿数据
		DeleteWishPlayerData(reqMsg.GetId(), delItemIds)
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 获取盲盒中奖配置列表 盲盒id:int64
func (self *cls4) RpcQueryWishBoxWinCfgList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := QueryWishBoxWinCfgList(reqMsg)
	ret := &brower_backstage.WishBoxWinCfgList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 获取盲盒列表键值对
func (self *cls4) RpcGetWishBoxKvs(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {
	pools := GetAllWishBoxOnline()
	var list []*brower_backstage.KeyValueTag
	for _, v := range pools {
		one := &brower_backstage.KeyValueTag{
			Key:   easygo.NewInt32(v.GetId()),
			Value: easygo.NewString(v.GetName()),
		}
		list = append(list, one)
	}
	ret := &brower_backstage.KeyValueResponseTag{
		List: list,
	}

	return ret
}

// 查看盲盒详情 Id64-id： 盲盒id
func (self *cls4) RpcGetWishBoxDetail(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	box := GetWishBoxById(reqMsg.GetId64())
	one := &brower_backstage.WishBox{
		Id:          easygo.NewInt64(box.GetId()),
		Name:        easygo.NewString(box.GetName()),
		Icon:        easygo.NewString(box.GetIcon()),
		GoodsAmount: easygo.NewInt32(box.GetTotalNum()),
		Attribute:   box.GetMenu(),
		UserId:      easygo.NewInt64(box.GetGuardianId()),
		Price:       easygo.NewInt64(box.GetPrice()),
		Status:      easygo.NewInt32(box.GetStatus()),
		CreateTime:  easygo.NewInt64(box.GetCreateTime()),
		UpdateTime:  easygo.NewInt64(box.GetUpdateTime()),
		SortWeight:  easygo.NewInt64(box.GetSortWeight()),
		WishPoolId:  easygo.NewInt64(box.GetWishPoolId()),
		IsRecommend: easygo.NewBool(box.GetIsRecommend()),
	}

	// 是否挑战赛
	if box.GetMatch() == 1 {
		one.IsChallenge = easygo.NewBool(true)
	} else {
		one.IsChallenge = easygo.NewBool(false)
	}

	u := for_game.GetWishPlayerByPid(one.GetUserId())
	one.UserId = easygo.NewInt64(u.GetPlayerId())
	user2 := GetPlayerInfoByPid(u.GetPlayerId())
	one.UserAccount = easygo.NewString(user2.GetAccount())

	req := &brower_backstage.ListRequest{Id: easygo.NewInt64(reqMsg.GetId64())}
	goods, _ := GetWishBoxGoodsItemList(req)

	one.ItemList = goods

	redis := GetPoolObj(one.GetWishPoolId())
	v := redis.GetPollInfoFromRedis()
	// 返回给前端传的是水池配置id,从数据库拿的是盲盒水池id
	one.WishPoolId = easygo.NewInt64(v.GetPoolConfigId())

	return one
}

//盲盒下的商品列表（下拉框）
func (self *cls4) RpcGetGoodsListByBoxId(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("RpcGetGoodsListByBoxId, %+v：", reqMsg)
	goodsList := GetGoodsItemByBoxId(reqMsg.GetId64())
	return &brower_backstage.WishBoxGoodsSelectedList{
		List: goodsList,
	}
}

// 盲盒抽奖
func (self *cls4) RpcWishBoxLottery(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishBoxLotteryReq) easygo.IMessage {
	logs.Info("RpcWishBoxLottery, %+v", reqMsg)
	count := reqMsg.GetCount()
	if count < 1 {
		return easygo.NewFailMsg("抽奖失败,抽奖次数不得小于1次！")
	}
	account := reqMsg.GetAccount()
	result := &brower_backstage.WishBoxLotteryResp{
		Result: easygo.NewInt32(2),
	}
	wishWhiteInfo, err := GetWishWhiteInfoByAccount(account)
	if err != nil {
		//return easygo.NewFailMsg("抽奖失败,请填写有效的白名单账号！")
		result.Msg = easygo.NewString("抽奖失败,请填写有效的白名单账号！")
		return result
	}
	playerId := wishWhiteInfo.GetId() //wishPlayer表的id字段
	wishPlayerInfo := GetWishPlayerInfo(playerId)
	userDiamond := wishPlayerInfo.GetDiamond()

	boxId := reqMsg.GetBoxId()
	wishBoxInfo := GetWishBoxById(boxId)
	var dareType int32              // 1-挑战赛;2-非挑战赛
	match := wishBoxInfo.GetMatch() //0普通赛，1挑战赛
	if match == 0 {
		match = 2
	}
	dareType = match

	if userDiamond < wishBoxInfo.GetPrice()*int64(count) {
		//return easygo.NewFailMsg("抽奖失败,账号钻石余额不足！")
		result.Msg = easygo.NewString("抽奖失败,账号钻石余额不足！")
		return result
	}

	productId := reqMsg.GetProductId()
	if productId == 0 {
		return easygo.NewFailMsg("商品ID不能为空")
	}

	wishProduct := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM, bson.M{"WishItemId": productId})
	if wishProduct == nil {
		return easygo.NewFailMsg("许愿商品不存在")
	}

	var successCount int32 = 0
	serInfo := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)
	logs.Info("serInfo:", serInfo)
	if serInfo == nil {
		return easygo.NewFailMsg("不存在的服务器连接")
	}
	for i := 0; i < int(count); i++ {
		//发起许愿
		wishReq := &h5_wish.WishReq{
			BoxId:     easygo.NewInt64(boxId),
			ProductId: easygo.NewInt64(wishProduct.(bson.M)["_id"]),
			OpType:    easygo.NewInt32(1),
		}

		wishResp, wishErr := SendMsgToServerNew(serInfo.GetSid(), "RpcWish", wishReq)
		if wishErr != nil {
			if wishErr.GetReason() != for_game.WISH_ERR_1 {
				continue
			}
		} else {
			resp, ok := wishResp.(*h5_wish.WishResp)
			if !(ok && resp.GetResult() == 1) {
				continue
			}
		}

		//发起挑战
		doDareReq := &h5_wish.DoDareReq{
			DareType:  easygo.NewInt32(dareType),
			WishBoxId: easygo.NewInt64(boxId),
		}
		_, dareErr := SendMsgToServerNew(serInfo.GetSid(), "RpcDoDare", doDareReq, playerId)
		if dareErr == nil {
			//resp, ok := doDareResp.(*h5_wish.DoDareResp)
			//if !(ok && resp.GetResult() == 1) {
			//	continue
			//}
			successCount++
		} else {
			logs.Info("RpcDoDare err:", dareErr.GetReason())
		}
	}
	result.Result = easygo.NewInt32(1)
	result.Msg = easygo.NewString(fmt.Sprintf("成功完成%d次抽奖！是否继续抽奖？", successCount))
	return result
}

// 新增/更新盲盒中奖配置列表
func (self *cls4) RpcUpdateWishBoxWinCfgList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishBoxWinCfgList) easygo.IMessage {
	logs.Info("RpcUpdateWishBoxWinCfgList, %+v：", reqMsg)

	list := reqMsg.GetList()

	for _, v := range list {
		if v.GetPrice() < 0 {
			return easygo.NewFailMsg("数值错误")
		}
		if v.GetBigLoss() < 0 {
			return easygo.NewFailMsg("数值错误")
		}
		if v.GetSmallLoss() < 0 {
			return easygo.NewFailMsg("数值错误")
		}
		if v.GetCommon() < 0 {
			return easygo.NewFailMsg("数值错误")
		}
		if v.GetBigWin() < 0 {
			return easygo.NewFailMsg("数值错误")
		}
		if v.GetSmallWin() < 0 {
			return easygo.NewFailMsg("数值错误")
		}
		UpdateWishBoxWinCfgList(v)
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, "新增/更新盲盒中奖配置列表")
	return easygo.EmptyMsg
}

// 商品管理
// 获取商品列表
func (self *cls4) RpcQueryWishGoodsList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.WishBoxGoodsListRequest) easygo.IMessage {
	list, count := QueryWishGoodsList(reqMsg)

	ret := &brower_backstage.WishBoxGoodsList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 获取商品品牌列表键值对
func (self *cls4) RpcGetWishGoodsBrandKvs(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, _ := QueryWishGoodsBrandList(reqMsg)

	var kvs []*brower_backstage.KeyValueTag
	for _, v := range list {
		one := &brower_backstage.KeyValueTag{
			Key:   easygo.NewInt32(v.GetId()),
			Value: easygo.NewString(v.GetName()),
		}
		kvs = append(kvs, one)
	}

	ret := &brower_backstage.KeyValueResponseTag{
		List: kvs,
	}

	return ret
}

// 新增/更新商品
func (self *cls4) RpcUpdateWishGoods(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishBoxGoods) easygo.IMessage {
	logs.Info("RpcUpdateWishGoods, %+v：", reqMsg)

	curTime := time.Now().Unix()
	if reqMsg.GetId() != 0 {
		goods := GetGoodsById(reqMsg.GetId())

		if reqMsg.GetStatus() != goods.GetStatus() {

			if reqMsg.GetStatus() == 0 {
				// 检查其他盲盒是否已经包含了该商品
				boxs := GetBoxByItemId(reqMsg.GetId())
				if len(boxs) != 0 {
					var boxids []int64
					for _, box := range boxs {
						boxids = append(boxids, box.GetId())
					}
					return easygo.NewFailMsg(fmt.Sprintf("当前上架盲盒ID%v含有该商品，只有下架盲盒才能下架商品", boxids))
				}
				//if CheckBoxHasGoods(reqMsg.GetId()) {
				//	return easygo.NewFailMsg("有盲盒有该商品，只有下架盲盒才能下架商品")
				//}
				reqMsg.SoldOutTime = easygo.NewInt64(curTime)
				reqMsg.UploadTime = easygo.NewInt64(0)
				reqMsg.UpdateTime = easygo.NewInt64(curTime)
			} else {
				reqMsg.UploadTime = easygo.NewInt64(curTime)
				reqMsg.UpdateTime = easygo.NewInt64(curTime)
				reqMsg.SoldOutTime = easygo.NewInt64(0)
			}
		}
	}

	msg := fmt.Sprintf("更新商品,%s", reqMsg.GetName())
	if reqMsg.GetId() == 0 {
		msg = fmt.Sprintf("新增商品,%s", reqMsg.GetName())
	}
	UpdateWishGoods(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 获取商品品牌列表
func (self *cls4) RpcQueryWishGoodsBrandList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := QueryWishGoodsBrandList(reqMsg)
	ret := &brower_backstage.WishGoodsBrandList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 新增/更新商品品牌
func (self *cls4) RpcUpdateWishGoodsBrand(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishGoodsBrand) easygo.IMessage {
	msg := fmt.Sprintf("更新商品品牌,%s", reqMsg.GetName())
	if reqMsg.GetId() == 0 {
		msg = fmt.Sprintf("新增商品品牌,%s", reqMsg.GetName())
	}
	UpdateWishGoodsBrand(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 获取商品类型列表
func (self *cls4) RpcQueryWishGoodsTypeList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := QueryWishGoodsTypeList(reqMsg)

	ret := &brower_backstage.WishGoodsTypeList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 获取商品类型列表键值对
func (self *cls4) RpcGetWishGoodsTypeKvs(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	list, _ := QueryWishGoodsTypeList(reqMsg)

	var kvs []*brower_backstage.KeyValueTag
	for _, v := range list {
		one := &brower_backstage.KeyValueTag{
			Key:   easygo.NewInt32(v.GetId()),
			Value: easygo.NewString(v.GetName()),
		}
		kvs = append(kvs, one)
	}

	ret := &brower_backstage.KeyValueResponseTag{
		List: kvs,
	}

	return ret
}

// 查看盲盒详情 Id64-id： 盲盒id
func (self *cls4) RpcGetWishGoodsDetail(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {

	box := GetWishItemById(reqMsg.GetId64())
	one := &brower_backstage.WishBoxGoods{
		Id:             easygo.NewInt64(box.GetId()),
		Name:           easygo.NewString(box.GetName()),
		Icon:           easygo.NewString(box.GetIcon()),
		Price:          easygo.NewInt64(box.GetPrice()),
		Status:         easygo.NewInt32(box.GetStatus()),
		WishBrandId:    easygo.NewInt32(box.GetBrand()),
		WishItemTypeId: easygo.NewInt32(box.GetType()),
		IsPreSale:      easygo.NewBool(box.GetIsPreSale()),
		ArrivalTime:    easygo.NewInt64(box.GetPreHaveTime()),
		StockAmount:    easygo.NewInt32(box.GetStockAmount()),
		UploadTime:     easygo.NewInt64(box.GetUploadTime()),
		SoldOutTime:    easygo.NewInt64(box.GetSoldOutTime()),
	}

	return one
}

// 新增/更新商品类型
func (self *cls4) RpcUpdateWishGoodsType(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishGoodsType) easygo.IMessage {
	logs.Info("RpcUpdateWishGoodsBrand, %+v：", reqMsg)
	msg := fmt.Sprintf("更新商品类型,%s", reqMsg.GetName())
	if reqMsg.GetId() == 0 {
		msg = fmt.Sprintf("新增商品类型,%s", reqMsg.GetName())
	}
	UpdateWishGoodsType(reqMsg)

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 订单管理
// 获取发货订单列表
func (self *cls4) RpcQueryWishDeliveryOrderList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishDeliveryOrderList, %+v：", reqMsg)
	list, count := QueryWishDeliveryOrderList(reqMsg)

	for i := range list {
		u := for_game.GetWishPlayerByPid(list[i].GetUserId())
		u2 := GetPlayerInfoByPid(u.GetPlayerId())
		list[i].UserId = easygo.NewInt64(u2.GetPlayerId())
		list[i].UserAccount = easygo.NewString(u2.GetAccount())
	}

	ret := &brower_backstage.WishDeliveryOrderList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 填写发货信息
func (self *cls4) RpcUpdateDeliveryOrderCourierInfo(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.UpdateDeliveryOrderCourierInfo) easygo.IMessage {
	logs.Info("RpcUpdateDeliveryOrderCourierInfo, %+v：", reqMsg)
	if reqMsg.GetOdd() == "" {
		return easygo.NewFailMsg("缺少快递单号")
	}
	order := GetWishPlayerExchangeLog(reqMsg.GetOrderId())
	if order == nil {
		return easygo.NewFailMsg("订单不存在")
	}
	if order.GetStatus() > 0 {
		return easygo.NewFailMsg("订单已不是待发货状态")
	}
	UpdateDeliveryOrderCourierInfo(reqMsg, user.GetAccount())
	UpdatePlayerWishItemStatusByPlayerExchangeLog([]int64{order.GetPlayerItemId()}, 1)

	msgText := fmt.Sprintf("您在我的愿望盒兑换的商品[%s]已发货，快递公司：%s；快递单号：%s", order.GetProductName(), reqMsg.GetCompany(), reqMsg.GetOdd())
	SendSystemNotice(order.GetUserId(), "发货通知", msgText)

	msg := fmt.Sprintf("填写发货信息,订单id:%d", reqMsg.GetOrderId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 取消发货
func (self *cls4) RpcUpdateDeliveryOrderStatus(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.UpdateStatusRequest) easygo.IMessage {
	logs.Info("RpcUpdateDeliveryOrderStatus, %+v：", reqMsg)
	UpdateDeliveryOrderStatus(reqMsg, user.GetAccount())
	order := GetWishPlayerExchangeLog(reqMsg.GetId())
	var status int32
	if reqMsg.GetStatus() == 2 {
		status = 0
	}
	UpdatePlayerWishItemStatusByPlayerExchangeLog([]int64{order.GetPlayerItemId()}, status)

	// 返回邮费-处理批量兑换订单
	orders := GetWishPlayerExchangeLogListByOrderId(order.GetOrderId())
	count := 0
	diamond := int64(0)
	for _, v := range orders {
		if v.GetStatus() == 2 {
			count++
		}
		diamond = v.GetPostage()
	}

	if count == len(orders) {
		logs.Info("返回运费")
		wishP := for_game.GetWishPlayerByPid(order.GetPlayerId())
		player := for_game.GetPlayerById(wishP.GetPlayerId())
		serInfo := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)

		req := &h5_wish.BackStageUpdateDiamondReq{
			Account:    easygo.NewString(player.GetPhone()),
			Channel:    easygo.NewInt32(1001),
			NickName:   easygo.NewString(player.GetNickName()),
			HeadUrl:    easygo.NewString(player.GetHeadIcon()),
			PlayerId:   easygo.NewInt64(player.GetPlayerId()),
			Token:      easygo.NewString(""),
			Diamond:    easygo.NewInt64(diamond),
			Reason:     easygo.NewString("运费返回"),
			SourceType: easygo.NewInt32(for_game.DIAMOND_TYPE_POSTAGE_FAILD),
		}
		httpErr := ChooseOneWish(serInfo.GetSid(), "RpcBackstageUpdateDiamond", req)
		err := for_game.ParseReturnDataErr(httpErr)
		if err != nil {
			logs.Warn("通知许愿池服务器更新钻石失败，rpc=RpcBackstageAddDiamond, req=%+v", req)
			return easygo.NewFailMsg("更新钻石失败")
		}
	}
	msg := fmt.Sprintf("取消发货,订单id:%d", reqMsg.GetId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

//  更新订单状态 Id : 订单id   Note: 原因（status: 2-已取消 3-已拒绝）
func (self *cls4) RpcUpdateWishRecycleOrderStatus(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.UpdateStatusRequest) easygo.IMessage {
	logs.Info("RpcUpdateWishRecycleOrderStatus, %+v：", reqMsg)
	order := GetWishRecycleOrder(reqMsg.GetId())
	if order.GetStatus() > 0 {
		return easygo.NewFailMsg("该订单已不是待审核状态")
	}
	UpdateWishRecycleOrderStatus(reqMsg, user.GetAccount())

	msg := ""
	switch reqMsg.GetStatus() {
	case 1:
		order := GetWishRecycleOrder(reqMsg.GetId())
		if order.GetStatus() == 1 {
			if order.GetBankCardId() != "" {
				player := GetWishPlayerInfoByPid(order.GetUserId())
				if player.GetPlayerId() == 0 {
					return easygo.NewFailMsg(fmt.Sprintf("用户不存在，id:%d", order.GetUserId()))
				}
				// 提现用户银行卡
				recycleMsg := &h5_wish.RecycleToHall{
					PlayerId:   easygo.NewInt64(player.GetPlayerId()),
					BankCardId: easygo.NewString(order.GetBankCardId()),
					Price:      easygo.NewInt64(order.GetRecyclePriceTotal()),
				}

				serInfo := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_HALL)

				_, serverErr := SendMsgToServerNew(serInfo.GetSid(), "RpcRecycleToBandCard", recycleMsg, player.GetPlayerId())
				if serverErr != nil {
					reqMsg.Status = easygo.NewInt32(0)
					UpdateWishRecycleOrderStatus(reqMsg, user.GetAccount())
					logs.Error("提现回用户银行卡：resp: %v, err: %v", serverErr, serverErr.GetReason())
					return easygo.NewFailMsg("提现到银行卡失败")
				}
			} else {
				coins := order.GetRecycleDiamond()
				player := GetWishPlayerInfoByPid(order.GetUserId())
				if player.GetPlayerId() == 0 {
					return easygo.NewFailMsg(fmt.Sprintf("用户不存在，id:%d", order.GetUserId()))
				}
				wishPlayer := for_game.GetRedisWishPlayer(player.GetId())
				if wishPlayer == nil {
					logs.Error("获取许愿池用户失败")
					return easygo.NewFailMsg("参数有误")
				}
				err2, _ := wishPlayer.AddDiamond(int64(coins), "回收钻石", for_game.DIAMOND_TYPE_WISH_BACK, nil)
				if err2 != nil {
					reqMsg.Status = easygo.NewInt32(0)
					UpdateWishRecycleOrderStatus(reqMsg, user.GetAccount())
					logs.Error("更新钻石失败，err:%v", err2)
					return easygo.NewFailMsg("更新钻石失败")
				}
			}
		}

		var ids []int64
		for _, v := range order.GetRecycleItemList() {
			ids = append(ids, v.GetPlayerItemId())
		}
		// 更新所回收物品的状态
		UpdatePlayerWishItemStatusByOrder(ids, 2)
		msg = fmt.Sprintf("通过回收,订单id:%d", reqMsg.GetId())
	case 2:
		order := GetWishRecycleOrder(reqMsg.GetId())
		var ids []int64
		for _, v := range order.GetRecycleItemList() {
			ids = append(ids, v.GetPlayerItemId())
		}
		UpdatePlayerWishItemStatusByOrder(ids, 0)
		msg = fmt.Sprintf("取消回收,订单id:%d", reqMsg.GetId())
		note := fmt.Sprintf("您在%s 发起的回收申请审核不通过，您可再次发起回收申请，如有疑问请咨询官方客服", easygo.Stamp2Str(order.GetInitTime()))
		SendSystemNotice(order.GetUserId(), "许愿池通知", note)
	case 3:
		order := GetWishRecycleOrder(reqMsg.GetId())
		var ids []int64
		for _, v := range order.GetRecycleItemList() {
			ids = append(ids, v.GetPlayerItemId())
		}
		UpdatePlayerWishItemStatusByOrder(ids, 4, user.GetAccount())
		msg = fmt.Sprintf("拒绝回收,订单id:%d", reqMsg.GetId())
		note := fmt.Sprintf("您在%s 发起的回收申请审核不通过，已扣除该笔订单所有物品，如有疑问请咨询官方客服", easygo.Stamp2Str(order.GetInitTime()))
		SendSystemNotice(order.GetUserId(), "许愿池通知", note)
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 回收订单用户审核详情 Id64-id： 订单id
func (self *cls4) RpcGetWishRecycleOrderUserInfo(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	order := GetWishRecycleOrder(reqMsg.GetId64())
	wishP := GetWishPlayerInfo(order.GetPlayerId())
	player := for_game.GetPlayerById(wishP.GetPlayerId())

	data := &brower_backstage.WishRecycleOrderUserInfo{
		PlayerId:     easygo.NewInt64(player.GetPlayerId()),
		Nickname:     easygo.NewString(player.GetNickName()),
		RealName:     easygo.NewString(player.GetRealName()),
		Account:      easygo.NewString(player.GetAccount()),
		PlayerType:   easygo.NewInt64(player.GetTypes()),
		Diamond:      easygo.NewInt32(order.GetRecycleDiamond()),
		CurDiamond:   easygo.NewInt64(wishP.GetDiamond()),
		RegisterTime: easygo.NewInt64(player.GetCreateTime()),
		LoginAddr:    easygo.NewString(player.GetLastLoginIP()),
	}

	return data
}

// 获取回收订单列表
func (self *cls4) RpcQueryWishRecycleOrderList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	list, count := QueryWishRecycleOrderList(reqMsg)

	for i := range list {
		u := for_game.GetPlayerById(list[i].GetUserId())
		list[i].UserAccount = easygo.NewString(u.GetAccount())
		list[i].UserId = easygo.NewInt64(u.GetPlayerId())
	}
	ret := &brower_backstage.WishRecycleOrderList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 回收订单详情 Id64-id： 订单id
func (self *cls4) RpcGetWishRecycleOrderDetail(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {

	list := GetWishRecycleItemList(reqMsg.GetId64())
	if reqMsg.GetIdStr() != "" {
		list = GetWishRecycleItemListByPaymentOrderId(reqMsg.GetIdStr())
	}

	ret := &brower_backstage.WishRecycleOrderDetailList{
		List:      list,
		PageCount: easygo.NewInt32(len(list)),
	}

	return ret
}

//查询线上充值订单
func (self *cls4) RpcQueryWishOrderList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryWishOrderRequest) easygo.IMessage {
	reqMsg.BeginTimestamp = easygo.NewInt64(reqMsg.GetBeginTimestamp() * 1000)
	reqMsg.EndTimestamp = easygo.NewInt64(reqMsg.GetEndTimestamp() * 1000)
	list, count := QueryWishOrderList(reqMsg)

	for i := range list {
		o := GetWishRecycleOrderByPaymentOrderId(list[i].GetOrderId())
		reason := GetWishRecycleReason(o.GetRecycleNote())
		list[i].PlayerReason = easygo.NewString(reason.GetReason())
	}

	return &brower_backstage.QueryOrderResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//订单操作
func (self *cls4) RpcOptWishOrder(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.OptOrderRequest) easygo.IMessage {
	logs.Info("RpcOptWishOrder:", reqMsg)
	if user.GetRole() > 0 {
		role := GetPowerRouter(user.GetRoleType())
		if for_game.IsContainsStr("orderManage", role.GetMenuIds()) == -1 {
			return easygo.NewFailMsg("权限不足")
		}
	}
	oid := reqMsg.GetOid()
	if oid == "" {
		return easygo.NewFailMsg("订单号错误")
	}

	order := for_game.GetRedisOrderObj(oid)
	if order == nil {
		return easygo.NewFailMsg("订单不存在")
	}

	switch order.GetStatus() {
	case 1:
		return easygo.NewFailMsg("订单已完成")
	case 2:
		return easygo.NewFailMsg("订单已审核")
	case 3:
		return easygo.NewFailMsg("订单已取消")
	case 4:
		return easygo.NewFailMsg("订单已拒绝")
	}

	s := fmt.Sprintf("订单操作:%s", order.GetOrderId())
	o := GetWishRecycleOrderByPaymentOrderId(order.GetOrderId())
	var ids []int64
	for _, v := range o.GetRecycleItemList() {
		ids = append(ids, v.GetPlayerItemId())
	}
	switch reqMsg.GetOpt() {
	case 1: // 完成订单
		err := FinishOrder(order.GetOrderId(), user.GetAccount(), reqMsg.GetNote())
		if err != nil {
			return err
		}
		UpdateWishOrder(order.GetOrderId(), 1)
		s = fmt.Sprintf("人工完成订单:%s", order.GetOrderId())

		// 更新所回收物品的状态
		UpdatePlayerWishItemStatusByOrder(ids, 2)
	case 2: //审核出款订单
		req := &server_server.AuditOrder{
			OrderId: easygo.NewString(oid),
		}
		result := ChooseOneHall(0, "RpcBsAuditOrder", req)
		err := for_game.ParseReturnDataErr(result)
		if err != nil {
			return easygo.NewFailMsg("请求第三方服务异常")
		} else {
			o := for_game.GetRedisOrderObj(oid)
			if o == nil {
				return easygo.NewFailMsg("订单不存在")
			}
			if o.GetStatus() == 1 {
				UpdateWishOrder(order.GetOrderId(), 1)
				// 更新所回收物品的状态
				UpdatePlayerWishItemStatusByOrder(ids, 2)
			} else if o.GetStatus() == 3 || o.GetStatus() == 4 {
				// 更新所回收物品的状态
				UpdatePlayerWishItemStatusByOrder(ids, 0)
			}
			s = fmt.Sprintf("审核订单:%s", order.GetOrderId())
		}

	case 3: // 取消/拒绝订单
		err := OptOrder(order.GetOrderId(), reqMsg.GetOpt(), user.GetAccount(), reqMsg.GetNote())
		if err != nil {
			return err
		}

		s = fmt.Sprintf("取消订单:%s", order.GetOrderId())

		UpdateWishOrder(order.GetOrderId(), 3)
		// 更新所回收物品的状态
		UpdatePlayerWishItemStatusByOrder(ids, 0)
		note := fmt.Sprintf("您在%s 发起的回收申请审核不通过，您可再次发起回收申请，如有疑问请咨询官方客服", easygo.Stamp2Str(order.GetCreateTime()))
		SendSystemNotice(order.GetPlayerId(), "许愿池通知", note)
	case 4:
		err := OptOrder(order.GetOrderId(), reqMsg.GetOpt(), user.GetAccount(), reqMsg.GetNote())
		if err != nil {
			return err
		}

		s = fmt.Sprintf("拒绝订单:%s", order.GetOrderId())

		UpdateWishOrder(order.GetOrderId(), 2)
		// 更新所回收物品的状态
		UpdatePlayerWishItemStatusByOrder(ids, 4, user.GetAccount())
		note := fmt.Sprintf("您在%s 发起的回收申请审核不通过，已扣除该笔订单所有物品，如有疑问请咨询官方客服", easygo.Stamp2Str(order.GetCreateTime()))
		SendSystemNotice(order.GetPlayerId(), "许愿池通知", note)
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, s)

	return easygo.EmptyMsg
}

// 许愿池报表列表
func (self *cls4) RpcQueryWishPoolReportList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishPoolReportList, %+v：", reqMsg)
	list, count := QueryWishPoolReportList(reqMsg)
	ret := &brower_backstage.WishPoolReportList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 盲盒报表列表
func (self *cls4) RpcQueryWishBoxReportList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishBoxReportList, %+v：", reqMsg)
	var (
		list  []*share_message.WishBoxReport
		count int32
	)

	list, count = QueryWishBoxReportList(reqMsg)
	ret := &brower_backstage.WishBoxReportList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 盲盒详情报表列表
func (self *cls4) RpcQueryWishBoxDetailReportList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishBoxDetailReportList, %+v：", reqMsg)
	var (
		list  []*share_message.WishBoxDetailReport
		count int32
	)
	if reqMsg.GetEndTimestamp() > 0 {
		MakeWishBoxDetailReportByTime(int64(reqMsg.GetType()), reqMsg.GetBeginTimestamp(), reqMsg.GetEndTimestamp())
		list, count = QueryWishBoxDetailReportTempList(reqMsg)
		DeleteWishBoxDetailReportTempList(reqMsg)
	} else {
		list, count = QueryWishBoxDetailReportList(reqMsg)
	}

	for i := range list {
		item := GetWishBoxItemById(list[i].GetItemId())
		goods := GetWishItemById(item.GetWishItemId())
		list[i].ItemName = easygo.NewString(goods.GetName())
	}

	ret := &brower_backstage.WishBoxDetailReportList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 商品报表列表
func (self *cls4) RpcQueryWishItemReportList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	list, count := QueryWishItemReportList(reqMsg)
	ret := &brower_backstage.WishItemReportList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// +++++++++++++++++++++++++++++++++++++++++++++++++盲盒详情导出测试数据
// 导出测试数据-玩家物品列表
func (self *cls4) RpcQueryTestPlayerWishItemList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	list, count := QueryPlayerWishItemList(reqMsg)
	ret := &brower_backstage.TestPlayerWishItemList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	logs.Info("导出测试数据-玩家物品列表，数据量：", len(list))
	return ret
}

// 导出测试数据-水池流水日志
func (self *cls4) RpcQueryTestWishPoolLogList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	list, count := QueryWishPoolLogList(reqMsg)
	ret := &brower_backstage.TestWishPoolLogList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	logs.Info("导出测试数据-水池流水日志，数据量：", len(list))
	return ret
}

// 导出测试数据-水池抽水日志
func (self *cls4) RpcQueryTestWishPoolPumpLogList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	list, count := QueryWishPoolPumpLogList(reqMsg)

	for i := range list {
		box := GetWishBoxById(list[i].GetBoxId())
		list[i].BoxName = easygo.NewString(box.GetName())

		pool1 := GetWishPool(list[i].GetPoolId())
		pool2 := GetWishPoolCfg(pool1.GetPoolConfigId())
		list[i].PoolName = easygo.NewString(pool2.GetName())

		// todo
		if list[i].GetBoxName() == "" {
			list[i].PoolName = easygo.NewString("测试数据")
		}
	}

	ret := &brower_backstage.TestWishPoolPumpLogList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	logs.Info("导出测试数据-水池抽水日志，数据量：", len(list))
	return ret
}

// 导出测试数据-盲盒水池信息
func (self *cls4) RpcQueryTestWishPoolBoxPoolInfoList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	box := GetWishBoxById(reqMsg.GetId())
	redis := GetPoolObj(box.GetWishPoolId())
	v := redis.GetPollInfoFromRedis()

	var list []*share_message.WishPool

	list = append(list, v)

	ret := &brower_backstage.TestWishPoolBoxPoolInfoList{
		List:      list,
		PageCount: easygo.NewInt32(1),
	}
	logs.Info("导出测试数据-水池信息，%+v：", ret)
	return ret
}

// 抽奖记录列表
func (self *cls4) RpcQueryDrawRecordList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryDrawRecordList，%+v：", reqMsg)
	ret := QueryDrawRecordList(reqMsg)
	list := ret.GetList()

	var ids []int64
	userM := make(map[int64]*brower_backstage.DrawRecord)
	for i := range list {
		ids = append(ids, list[i].GetUserId())
		userM[list[i].GetUserId()] = list[i]

		u := for_game.GetWishPlayerByPid(list[i].GetUserId())
		u2 := GetPlayerInfoByPid(u.GetPlayerId())
		list[i].UserId = easygo.NewInt64(u2.GetPlayerId())
		list[i].UserAccount = easygo.NewString(u2.GetAccount())
		list[i].Phone = easygo.NewString(u2.GetPhone())
		list[i].UserNickname = easygo.NewString(u2.GetNickName())

		var userType int32 //1普通用户 2运营号 3白名单
		if u.GetTypes() <= 2 {
			userType = 1
		} else {
			userType = 2
		}
		wishWhiteCount := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WHITE, bson.M{"_id": u.GetId()})
		if wishWhiteCount > 0 {
			userType = 3
		}
		list[i].UserType = easygo.NewInt32(userType)

	}

	addBoxCounts := GetAddBoxCountByUserIds(ids)
	wishItemCount := GetWishDataCountByUserIds(ids)
	haveItemCount := GetPlayerWishItems(ids)
	delItemCount := GetPlayerDelWishItemsCount(ids)
	for _, v := range addBoxCounts {
		if mv, ok := userM[v.GetUserId()]; ok {
			mv.AddBoxCount = easygo.NewInt32(v.GetAddBoxCount())
		}
	}

	for _, v := range wishItemCount {
		if mv, ok := userM[v.GetUserId()]; ok {
			mv.WishItemCount = easygo.NewInt32(v.GetWishItemCount())
		}
	}

	for _, v := range haveItemCount {
		if mv, ok := userM[v.GetUserId()]; ok {
			mv.HaveItemCount = easygo.NewInt32(v.GetHaveItemCount())
		}
	}

	for _, v := range delItemCount {
		if mv, ok := userM[v.GetUserId()]; ok {
			mv.DelItemCount = easygo.NewInt32(v.GetDelItemCount())
		}
	}

	logs.Info("ret: %v", ret)
	return ret
}

// 收藏盲盒记录列表
func (self *cls4) RpcQueryAddBoxRecordList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	player := GetWishPlayerInfoByPid(reqMsg.GetId())
	reqMsg.Id = easygo.NewInt64(player.GetId())
	ret := QueryAddBoxRecordList(reqMsg)

	return ret
}

// 许愿物品记录列表
func (self *cls4) RpcQueryWishGoodsRecordList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	player := GetWishPlayerInfoByPid(reqMsg.GetId())
	reqMsg.Id = easygo.NewInt64(player.GetId())
	ret := QueryWishGoodsRecordList(reqMsg)

	return ret
}

// 抽奖盲盒记录列表
func (self *cls4) RpcQueryDrawBoxRecordList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	player := GetWishPlayerInfoByPid(reqMsg.GetId())
	reqMsg.Id = easygo.NewInt64(player.GetId())
	ret := DrawBoxRecordList(reqMsg)
	var ids []int64
	userM := make(map[int64]*brower_backstage.DrawBoxRecord)
	for _, v := range ret.GetList() {
		ids = append(ids, v.GetBoxId())
		userM[v.GetBoxId()] = v
	}

	boxs := GetWishBoxList(ids)

	for _, v := range boxs {
		if mv, ok := userM[v.GetId()]; ok {
			mv.BoxName = easygo.NewString(v.GetName())
		}
	}

	return ret
}

// 现有物品列表
func (self *cls4) RpcQueryHaveItemList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	player := GetWishPlayerInfoByPid(reqMsg.GetId())
	reqMsg.Id = easygo.NewInt64(player.GetId())
	ret := PlayerWishItemList(reqMsg)

	return ret
}

// 扣除用户现有物品
func (self *cls4) RpcDeleteHaveItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	logs.Info("RpcDeleteHaveItem, %+v：", reqMsg)
	pids := reqMsg.GetObjIds()
	if len(pids) == 0 {
		return easygo.NewFailMsg("要扣除的用户ID参数错误")
	}
	ids := reqMsg.GetIds64()
	msg := ""
	if len(ids) > 0 {
		DeletePlayerWishItem(ids, user.GetAccount())
		idstr := ""
		for i, j := range ids {
			idstr += easygo.AnytoA(j)
			if i < len(ids)-1 {
				idstr += ","
			}
		}
		msgText := "系统检测到您有违规行为，已扣除您部分物品，如有疑问请咨询官方客服"
		SendSystemNotice(pids[0], "扣除物品通知", msgText)
		msg = fmt.Sprintf("扣除用户待兑换物品%s", idstr)
	} else {
		wish := for_game.GetWishPlayerInfoByImId(pids[0])
		if wish == nil {
			return easygo.NewFailMsg("许愿池帐户不存在")
		}
		DeletePlayerWishItemByPid(wish.GetPlayerId(), user.GetAccount())
		msgAllText := "系统检测到您有违规行为，已扣除您全部物品，如有疑问请咨询官方客服"
		SendSystemNotice(int64(pids[0]), "扣除物品通知", msgAllText)
		msg = "扣除用户所有待兑换物品"
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 中奖记录记录列表
func (self *cls4) RpcQueryWinRecordList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWinRecordList, %+v：", reqMsg)
	player := GetWishPlayerInfoByPid(reqMsg.GetId())
	uid := player.GetId()
	reqMsg.Id = easygo.NewInt64(uid)
	ret := QueryWinRecordList(reqMsg)
	list := ret.List
	// 解决数据量太大，一次性查询
	var itemIds []int64
	var goodsIds []int64
	itemM := make(map[int64]*share_message.PlayerWishData)
	goodsM := make(map[int64]*share_message.WishItem)
	for _, v := range list {
		itemIds = append(itemIds, v.GetBoxItemId())
		goodsIds = append(goodsIds, v.GetGoodsId())
	}

	items := GetPlayerWishDataByWishBoxItemIds(uid, itemIds)
	for i := range items {
		itemM[items[i].GetWishBoxItemId()] = items[i]
	}

	for _, v := range list {

		data := itemM[v.GetBoxItemId()]
		if data.GetId() != 0 {
			v.HasWish = easygo.NewBool(true)
		} else {
			v.HasWish = easygo.NewBool(false)
		}
	}

	goods := GetWishItemByIds(goodsIds)
	for _, v := range goods {
		goodsM[v.GetId()] = v
	}

	for _, v := range list {

		data := goodsM[v.GetGoodsId()]
		v.GoodsName = easygo.NewString(data.GetName())
	}

	return ret
}

// 扣除物品列表
func (self *cls4) RpcQueryWishDelItemList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	player := GetWishPlayerInfoByPid(reqMsg.GetId())
	reqMsg.Id = easygo.NewInt64(player.GetId())
	ret := PlayerWishDelItemList(reqMsg)

	return ret
}

// 水池管理
// ++++++++++++++++++++++++++++++++++++++++++++++++++许愿池+++++++++++++++++++
// 获取水池列表
func (self *cls4) RpcQueryWishPoolList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	list, count := QueryWishPoolList(reqMsg)

	ret := &brower_backstage.WishPoolList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return ret
}

// 新增/更新水池
func (self *cls4) RpcUpdateWishPool(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishPool) easygo.IMessage {
	logs.Info("RpcUpdateWishPool, %+v：", reqMsg)

	if reqMsg.GetPoolLimit() < 0 || reqMsg.GetPoolLimit() > 987654321 {
		return easygo.NewFailMsg("参数错误")
	}

	if reqMsg.GetShowInitialValue() < 0 || reqMsg.GetShowInitialValue() > 100 {
		return easygo.NewFailMsg("参数错误")
	}

	if reqMsg.GetShowCloseAward() < 0 || reqMsg.GetShowCloseAward() > 100 {
		return easygo.NewFailMsg("参数错误")
	}

	if reqMsg.GetShowStartAward() < 0 || reqMsg.GetShowStartAward() > 100 {
		return easygo.NewFailMsg("参数错误")
	}

	if reqMsg.GetShowCommission() < 0 || reqMsg.GetShowCommission() > 100 {
		return easygo.NewFailMsg("参数错误")
	}
	if reqMsg.GetShowRecycle() < 0 || reqMsg.GetShowCommission() > 100 {
		return easygo.NewFailMsg("参数错误")
	}

	if reqMsg.SmallWin.GetShowMaxValue() < 0 || reqMsg.SmallWin.GetShowMinValue() < 0 {
		return easygo.NewFailMsg("参数错误")
	}

	if reqMsg.SmallLoss.GetShowMaxValue() < 0 || reqMsg.SmallLoss.GetShowMinValue() < 0 {
		return easygo.NewFailMsg("参数错误")
	}

	if reqMsg.Common.GetShowMaxValue() < 0 || reqMsg.Common.GetShowMinValue() < 0 {
		return easygo.NewFailMsg("参数错误")
	}

	if reqMsg.BigWin.GetShowMaxValue() < 0 || reqMsg.BigWin.GetShowMinValue() < 0 {
		return easygo.NewFailMsg("参数错误")
	}
	if reqMsg.BigLoss.GetShowMaxValue() < 0 || reqMsg.BigLoss.GetShowMinValue() < 0 {
		return easygo.NewFailMsg("参数错误")
	}

	msg := fmt.Sprintf("更新水池,水池id:%d", reqMsg.GetId())
	if reqMsg.GetId() == 0 {
		msg = fmt.Sprintf("新增水池")
	}

	UpdateWishPool(reqMsg)
	// 同步redis信息
	if reqMsg.GetId() > 0 {
		poolId := reqMsg.GetId()
		ids := GetSubWishPoolIds(reqMsg.GetId())

		for _, id := range ids {
			one := reqMsg
			one.Id = easygo.NewInt64(id)
			one.LocalStatus = easygo.NewInt64(0)
			one.IsOpenAward = easygo.NewBool(false)
			one.PoolCfgId = easygo.NewInt64(poolId)
			redis := GetPoolObj(id)
			UpdateWishPool2(one)

			redis.DeleteRedisPollInfo()
			DeleteWishPoolLog(id)
			DeleteWishPoolPumpLog(id)
			redis.RefreshRedisPollInfo()
		}
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 删除水池Id64-id： 水池id
func (self *cls4) RpcDeleteWishPool(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	logs.Info("RpcDeleteWishPool, %+v：", reqMsg)

	qBson := bson.M{"WishPoolId": bson.M{"$in": reqMsg.GetIds64()}}
	data := GetWishBoxByBson(qBson)
	if data != nil {
		return easygo.NewFailMsg("正被调用的水池，不可被删除")
	}
	DeleteWishPool(reqMsg.GetIds64())
	msg := fmt.Sprintf("删除水池,水池id:%d", reqMsg.GetIds64())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 重置水池Id64-id： 水池id
func (self *cls4) RpcResetWishPool(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("RpcResetWishPool, %+v：", reqMsg)
	redis := GetPoolObj(reqMsg.GetId64())
	v := redis.GetPollInfoFromRedis()
	one := &brower_backstage.WishPool{
		Id:               easygo.NewInt64(v.GetId()),
		PoolLimit:        easygo.NewInt32(v.GetPoolLimit()),
		Name:             easygo.NewString(v.GetName()),
		ShowInitialValue: easygo.NewInt32(v.GetShowInitialValue()),
		ShowRecycle:      easygo.NewInt32(v.GetShowRecycle()),
		ShowCommission:   easygo.NewInt32(v.GetShowCommission()),
		ShowStartAward:   easygo.NewInt32(v.GetShowStartAward()),
		ShowCloseAward:   easygo.NewInt32(v.GetShowCloseAward()),
		IsDefault:        easygo.NewBool(v.GetIsDefault()),
		LocalStatus:      easygo.NewInt64(0),
		IsOpenAward:      easygo.NewBool(false),
	}

	one.PoolCfgId = easygo.NewInt64(v.GetPoolConfigId())

	one.SmallLoss = &brower_backstage.WishPoolStatus{
		ShowMaxValue: easygo.NewInt32(v.SmallLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.SmallLoss.GetShowMinValue()),
	}

	one.SmallWin = &brower_backstage.WishPoolStatus{
		ShowMaxValue: easygo.NewInt32(v.SmallWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.SmallWin.GetShowMinValue()),
	}

	one.BigLoss = &brower_backstage.WishPoolStatus{
		ShowMaxValue: easygo.NewInt32(v.BigLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.BigLoss.GetShowMinValue()),
	}

	one.BigWin = &brower_backstage.WishPoolStatus{
		ShowMaxValue: easygo.NewInt32(v.BigWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.BigWin.GetShowMinValue()),
	}
	one.Common = &brower_backstage.WishPoolStatus{
		ShowMaxValue: easygo.NewInt32(v.Common.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.Common.GetShowMinValue()),
	}

	UpdateWishPool2(one)
	redis.DeleteRedisPollInfo()
	DeleteWishPoolLog(one.GetId())
	DeleteWishPoolPumpLog(one.GetId())
	redis.RefreshRedisPollInfo()
	msg := fmt.Sprintf("重置水池,水池id:%d", reqMsg.GetId64())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 获取水池Id64-id： 水池id
func (self *cls4) RpcGetWishPool(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("RpcGetWishPool, %+v：", reqMsg)
	redis := GetPoolObj(reqMsg.GetId64())
	v := redis.GetPollInfoFromRedis()
	if v == nil && reqMsg.GetIdStr() != "tool" {
		return easygo.NewFailMsg("水池不存在")
	}
	one := &brower_backstage.WishPool{
		Id:               easygo.NewInt64(v.GetId()),
		PoolLimit:        easygo.NewInt32(v.GetPoolLimit()),
		InitialValue:     easygo.NewInt32(v.GetInitialValue()),
		IncomeValue:      easygo.NewInt32(v.GetIncomeValue()),
		Name:             easygo.NewString(v.GetName()),
		CreateTime:       easygo.NewInt64(v.GetCreateTime()),
		Recycle:          easygo.NewInt32(v.GetRecycle()),
		Commission:       easygo.NewInt32(v.GetCommission()),
		StartAward:       easygo.NewInt32(v.GetStartAward()),
		CloseAward:       easygo.NewInt32(v.GetCloseAward()),
		ShowInitialValue: easygo.NewInt32(v.GetShowInitialValue()),
		ShowRecycle:      easygo.NewInt32(v.GetShowRecycle()),
		ShowCommission:   easygo.NewInt32(v.GetShowCommission()),
		ShowStartAward:   easygo.NewInt32(v.GetShowStartAward()),
		ShowCloseAward:   easygo.NewInt32(v.GetShowCloseAward()),
		IsDefault:        easygo.NewBool(v.GetIsDefault()),
		IsOpenAward:      easygo.NewBool(v.GetIsOpenAward()),
		LocalStatus:      easygo.NewInt64(v.GetLocalStatus()),
	}

	one.SmallLoss = &brower_backstage.WishPoolStatus{
		MaxValue:     easygo.NewInt32(v.SmallLoss.GetMaxValue()),
		MinValue:     easygo.NewInt32(v.SmallLoss.GetMinValue()),
		ShowMaxValue: easygo.NewInt32(v.SmallLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.SmallLoss.GetShowMinValue()),
	}

	one.SmallWin = &brower_backstage.WishPoolStatus{
		MaxValue:     easygo.NewInt32(v.SmallWin.GetMaxValue()),
		MinValue:     easygo.NewInt32(v.SmallWin.GetMinValue()),
		ShowMaxValue: easygo.NewInt32(v.SmallWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.SmallWin.GetShowMinValue()),
	}

	one.BigLoss = &brower_backstage.WishPoolStatus{
		MaxValue:     easygo.NewInt32(v.BigLoss.GetMaxValue()),
		MinValue:     easygo.NewInt32(v.BigLoss.GetMinValue()),
		ShowMaxValue: easygo.NewInt32(v.BigLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.BigLoss.GetShowMinValue()),
	}

	one.BigWin = &brower_backstage.WishPoolStatus{
		MaxValue:     easygo.NewInt32(v.BigWin.GetMaxValue()),
		MinValue:     easygo.NewInt32(v.BigWin.GetMinValue()),
		ShowMaxValue: easygo.NewInt32(v.BigWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.BigWin.GetShowMinValue()),
	}
	one.Common = &brower_backstage.WishPoolStatus{
		MaxValue:     easygo.NewInt32(v.Common.GetMaxValue()),
		MinValue:     easygo.NewInt32(v.Common.GetMinValue()),
		ShowMaxValue: easygo.NewInt32(v.Common.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt32(v.Common.GetShowMinValue()),
	}
	logs.Debug(one)
	return one
}

// 更新默认水池水池 int64:新的 int32 旧的
func (self *cls4) RpcUpdateDefaultWish(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("RpcUpdateDefaultWish, %+v：", reqMsg)

	if reqMsg.GetId64() != 0 {
		new := &share_message.WishPool{
			Id:        easygo.NewInt64(reqMsg.GetId64()),
			IsDefault: easygo.NewBool(true),
		}
		UpdateWishPoolDB(new)
	}

	if reqMsg.GetId32() != 0 {
		old := &share_message.WishPool{
			Id:        easygo.NewInt64(reqMsg.GetId32()),
			IsDefault: easygo.NewBool(false),
		}
		UpdateWishPoolDB(old)
	}
	msg := fmt.Sprintf("更新默认水池,水池id:%d", reqMsg.GetId32())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)

	return easygo.EmptyMsg
}

// 水池键值对
func (self *cls4) RpcGetWishPoolKvs(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {
	list, _ := QueryWishPoolList(&brower_backstage.ListRequest{})
	var kvs []*brower_backstage.KeyValueTag
	for _, v := range list {
		one := &brower_backstage.KeyValueTag{
			Key:   easygo.NewInt32(v.GetId()),
			Value: easygo.NewString(v.GetName()),
		}
		if v.GetIsDefault() {
			one.Value = easygo.NewString(easygo.IntToString(int(v.GetId())) + "（默认）")
		} else {
			one.Value = easygo.NewString(easygo.IntToString(int(v.GetId())))
		}

		kvs = append(kvs, one)
	}

	ret := &brower_backstage.KeyValueResponseTag{
		List: kvs,
	}

	return ret
}

// 参数设置
// 获取价格参数设置
func (self *cls4) RpcGetPriceSection(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {

	data := GetPriceSection()
	return data
}

// 更新价格参数设置
func (self *cls4) RpcUpdatePriceSection(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.PriceSection) easygo.IMessage {
	logs.Info("RpcUpdatePriceSection, %+v：", reqMsg)
	UpdatePriceSection(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, "更新价格参数设置")
	return easygo.EmptyMsg
}

// 邮寄参数设置
// 获取邮寄参数设置
func (self *cls4) RpcGetMailSection(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {

	data := GetMailSection()
	return data
}

// 更新邮寄参数设置
func (self *cls4) RpcUpdateMailSection(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishMailSection) easygo.IMessage {
	logs.Info("RpcUpdateMailSection, %+v：", reqMsg)
	UpdateMailSection(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, "更新邮寄参数设置")
	return easygo.EmptyMsg
}

//保存回收说明
func (self *cls4) RpcSaveRecycleNoteCfg(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.RecycleNoteCfg) easygo.IMessage {
	logs.Info("RpcSaveRecycleNoteCfg, %+v：", reqMsg)
	data := &share_message.RecycleNoteCfg{
		Id:   easygo.NewString(for_game.TABLE_WISH_RECYCLE_NOTE_CFG),
		Text: reqMsg.GetText(),
	}
	SaveRecycleNoteCfg(data)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, "保存回收说明")
	return easygo.EmptyMsg
}

//查看回收说明
func (self *cls4) RpcGetRecycleNoteCfg(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("RpcGetRecycleNoteCfg, %+v：", reqMsg)
	data := GetRecycleNoteCfg()
	return &brower_backstage.RecycleNoteCfg{
		Text: data.GetText(),
	}
}

// 物品回收参数设置
// 获取物品回收参数设置
func (self *cls4) RpcGetWishRecycleSection(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {

	data := GetWishRecycleSection()

	return data
}

// 更新物品回收参数设置
func (self *cls4) RpcUpdateWishRecycleSection(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishRecycleSection) easygo.IMessage {
	logs.Info("RpcUpdateWishRecycleSection, %+v：", reqMsg)
	UpdateWishRecycleSection(reqMsg)
	for_game.InitConfigWishPayment()
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, "更新物品回收参数设置")
	return easygo.EmptyMsg
}

// 支付预警设置
// 获取支付预警
func (self *cls4) RpcGetWishPayWarnCfg(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {

	data := GetWishPayWarnCfg()

	return data
}

// 更新支付预警
func (self *cls4) RpcUpdateWishPayWarnCfg(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishPayWarnCfg) easygo.IMessage {
	logs.Info("RpcUpdateWishPayWarnCfg, %+v：", reqMsg)
	UpdateWishPayWarnCfg(reqMsg)
	for_game.InitConfigWishPayWarn()
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, "更新支付预警")
	return easygo.EmptyMsg
}

// 获取冷却期参数设置
func (self *cls4) RpcGetWishCoolDownConfig(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {

	data := GetWishCoolDownConfigFromDB()

	return data
}

// 更新冷却期参数设置
func (self *cls4) RpcUpdateWishCoolDownConfig(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishCoolDownConfig) easygo.IMessage {
	logs.Info("RpcUpdateWishCoolDownConfig, %+v：", reqMsg)
	UpdateWishCoolDownConfigFromDB(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, "更新冷却期参数设置")
	return easygo.EmptyMsg
}

// 获取货币换算参数设置
func (self *cls4) RpcGetWishCurrencyConversionCfg(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {

	data := GetWishCurrencyConversionCfg()

	return data
}

// 更新货币换算参数设置
func (self *cls4) RpcUpdateWishCurrencyConversionCfg(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishCurrencyConversionCfg) easygo.IMessage {
	logs.Info("RpcUpdateWishCurrencyConversionCfg, %+v：", reqMsg)
	UpdateWishCurrencyConversionCfg(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, "更新货币换算参数设置")
	return easygo.EmptyMsg
}

// 获取守护者收益设置
func (self *cls4) RpcGetWishGuardianCfg(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {

	data := GetWishGuardianCfg()

	return data
}

// 更新守护者收益设置
func (self *cls4) RpcUpdateWishGuardianCfg(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishGuardianCfg) easygo.IMessage {
	logs.Info("RpcUpdateWishGuardianCfg, %+v：", reqMsg)
	UpdateWishGuardianCfg(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, "更新守护者收益设置")
	return easygo.EmptyMsg
}

//钻石管理列表
func (self *cls4) RpcDiamondItemList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-Sort", "EndTime"}
	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["Platform"] = reqMsg.GetListType()
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_DIAMOND_RECHARGE, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.DiamondRecharge
	for _, li := range lis {
		one := &share_message.DiamondRecharge{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.DiamondItemListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

// 钻石保存
func (self *cls4) RpcSaveDiamondItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.DiamondRecharge) easygo.IMessage {
	if reqMsg.Diamond == nil && reqMsg.GetDiamond() < 1 {
		return easygo.NewFailMsg("钻石数量错误")
	}

	if reqMsg.GetCoinPrice() < 1 {
		return easygo.NewFailMsg("价格不能小于1硬币")
	}

	if reqMsg.GetMonthFirst() > 0 && reqMsg.GetGiveDiamond() > 0 {
		return easygo.NewFailMsg("月首次赠送和其他赠送优惠不可同时存在，两者选其一")
	}

	msg := fmt.Sprintf("修改钻石商品:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_DIAMOND_RECHARGE))
		msg = fmt.Sprintf("添加钻石商品:%d", reqMsg.GetId())
	}

	if reqMsg.GetEndTime() == 0 && reqMsg.GetGiveDiamond() > 0 {
		if reqMsg.GetStatus() != 2 {
			return easygo.NewFailMsg("请先设置活动时间")
		}
	}

	if reqMsg.MonthFirst == nil {
		reqMsg.MonthFirst = easygo.NewInt64(0)
	}

	if reqMsg.GiveDiamond == nil {
		reqMsg.GiveDiamond = easygo.NewInt64(0)
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_DIAMOND_RECHARGE, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)

	return easygo.EmptyMsg
}

//系统赠送硬币
func (self *cls4) RpcGiveDiamond(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("RpcGiveDiamond reqMsg: %v", reqMsg)
	if reqMsg.IdStr == nil && reqMsg.GetIdStr() == "" {
		return easygo.NewFailMsg("用户柠檬号不能为空")
	}

	if reqMsg.GetId32() == 0 {
		return easygo.NewFailMsg("赠送数量不能为0")
	}

	playerId := reqMsg.GetId64()
	player := &share_message.PlayerBase{}
	if playerId == 0 {
		player = QueryPlayerbyAccount(reqMsg.GetIdStr())
	} else {
		player = QueryPlayerbyId(playerId)
	}

	if player == nil {
		return easygo.NewFailMsg("账号信息有误")
	}

	//var st int32
	var msg string
	var msgText string
	coins := reqMsg.GetId32()
	//note := ""
	switch reqMsg.GetNote() {
	case "give":
		if coins < 0 {
			coins = -coins
		}
		//note = "后台赠送"
		//st = for_game.DIAMOND_TYPE_BACK_GIVE
		msg = fmt.Sprintf("系统赠送[%s]钻石[%d]个", reqMsg.GetIdStr(), reqMsg.GetId32())
	case "back":
		if coins > 0 {
			coins = -coins
		}
		//note = "后台扣除"
		//st = for_game.DIAMOND_TYPE_BACK_RECYCLE
		msgText = fmt.Sprintf("系统检测到您有违规行为，已扣除您%d钻石，如有疑问请咨询官方客服", -coins)
		msg = fmt.Sprintf("系统回收[%s]钻石[%d]个", reqMsg.GetIdStr(), reqMsg.GetId32())
	case "freeze": // TODO

	default:
		return easygo.NewFailMsg("IdStr参数错误")
	}

	serInfo := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)

	req := &h5_wish.BackStageAddDiamondReq{
		Account:  easygo.NewString(player.GetPhone()),
		Channel:  easygo.NewInt32(1001),
		NickName: easygo.NewString(player.GetNickName()),
		HeadUrl:  easygo.NewString(player.GetHeadIcon()),
		PlayerId: easygo.NewInt64(player.GetPlayerId()),
		Token:    easygo.NewString(""),
		Diamond:  easygo.NewInt64(coins),
	}
	httpErr := ChooseOneWish(serInfo.GetSid(), "RpcBackstageAddDiamond", req)
	err := for_game.ParseReturnDataErr(httpErr)
	if err != nil {
		logs.Warn("通知许愿池服务器更新钻石失败，rpc=RpcBackstageAddDiamond, req=%+v", req)
		return easygo.NewFailMsg("更新钻石失败")
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	if len(msgText) > 0 {
		SendSystemNotice(player.GetPlayerId(), "扣除钻石通知", msgText)
	}

	return easygo.EmptyMsg
}

// 钻石流水日志
func (self *cls4) RpcQueryDiamondChangeLog(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	if reqMsg.GetCurPage() == 1 {
		for_game.SaveDiamondChangeLogToMongoDB()
	}
	logs.Info("RpcUpdateWishGuardianCfg, %+v：", reqMsg)
	list, count := QueryDiamondChangeLog(reqMsg)

	for i := range list {
		u := for_game.GetWishPlayerByPid(list[i].GetPlayerId())
		base := QueryPlayerbyId(u.GetPlayerId())
		list[i].PlayerId = easygo.NewInt64(base.GetPlayerId())
		list[i].Account = easygo.NewString(base.GetAccount())
	}

	return &brower_backstage.DiamondChangeLogResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//付费用户地理位置分布图
func (self *cls4) RpcPayPlayerLocation(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	queryBson := bson.M{}
	if reqMsg.GetListType() > 0 {
		queryBson["DeviceType"] = reqMsg.GetListType()
	}

	switch reqMsg.GetType() {
	case 1:
		queryBson["Piece"] = "CN"
	case 2:
		queryBson["Piece"] = bson.M{"$ne": "CN"}
	default:
		queryBson["Piece"] = bson.M{"$ne": nil}
	}

	m := []bson.M{
		{"$match": queryBson},
		{"$group": bson.M{"_id": "$Position", "Count": bson.M{"$sum": 1}}},
	}

	list := for_game.FindPipeAll(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_WISH_PAYPLAYER_LOCATION_LOG, m, 0, 0)
	line := []*brower_backstage.NameValueTag{}
	var total int64 = 0
	for _, k := range list {
		one := &share_message.PipeStringCount{}
		for_game.StructToOtherStruct(k, one)
		lineOne := &brower_backstage.NameValueTag{
			Name:  easygo.NewString(one.GetId()),
			Value: easygo.NewInt64(one.GetCount()),
		}
		total += int64(one.GetCount())
		line = append(line, lineOne)
	}

	sort.Slice(line, func(i int, j int) bool {
		return line[i].GetValue() > line[j].GetValue()
	})

	return &brower_backstage.NameValueResponseTag{
		List:  line,
		Total: easygo.NewInt64(total),
	}
}

//许愿池埋点报表查询
func (self *cls4) RpcQueryWishLogReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-_id"}
	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_LOG, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.WishLogReport
	for _, l := range lis {
		one := &share_message.WishLogReport{}
		for_game.StructToOtherStruct(l, one)
		list = append(list, one)
	}
	return &brower_backstage.QueryWishLogReportRes{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//许愿池用户管理 Type下拉查询项,ListType渠道类型:1001-柠檬im 1002-语音渠道 1003-其他渠道,Status状态:1-冻结 2解冻
func (self *cls4) RpcWishPlayerList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcWishPlayerList, %+v：", reqMsg)
	findBson := bson.M{}
	sort := []string{"-_id"}
	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}
	if reqMsg.Status != nil && reqMsg.GetStatus() > 0 {
		findBson["IsFreeze"] = easygo.If(reqMsg.GetStatus() == 1, true, false)
	}
	if reqMsg.ListType != nil && reqMsg.GetListType() > 0 {
		findBson["Channel"] = reqMsg.GetListType()
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1: //许愿池帐号查询
			findBson["Account"] = reqMsg.GetKeyword()
		case 2: //昵称查询
			findBson["NickName"] = reqMsg.GetKeyword()
		case 3: //渠道帐号查询
			findBson["PlayerId"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		}
	}
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.WishPlayer
	for _, li := range lis {
		one := &share_message.WishPlayer{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	return &brower_backstage.WishPlayerListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//批量冻结钻石帐户
func (s *cls4) RpcWishPlayerFreezeDiamond(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要冻结的项")
	}

	times := reqMsg.GetObjIds()
	if len(times) == 0 {
		return easygo.NewFailMsg("冻结到期时间不能为空")
	}

	for _, id := range idList {
		mgr := for_game.GetRedisWishPlayer(id)
		mgr.SetIsFreeze(true)
		mgr.SetFreezeTime(times[0])
		mgr.SetNote(reqMsg.GetNote())
		mgr.SetOperator(user.GetAccount())
		SendSystemNotice(mgr.GetPlayerId(), "冻结许愿池账号通知", "系统检测到您有违规行为，已冻结您的许愿池账号，如有疑问请咨询官方客服") //柠檬助手通知
	}

	var ids string
	count := len(idList)
	for i, t := range idList {
		ids += easygo.AnytoA(t)
		if i < count-1 {
			ids += ","
		}
	}

	msg := fmt.Sprintf("批量冻结钻石帐户: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)

	return easygo.EmptyMsg
}

//批量解冻钻石帐户
func (s *cls4) RpcWishPlayerUnFreezeDiamond(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要冻结的项")
	}

	for _, id := range idList {
		mgr := for_game.GetRedisWishPlayer(id)
		mgr.SetIsFreeze(false)
		mgr.SetNote("")
		mgr.SetOperator(user.GetAccount())
	}

	var ids string
	count := len(idList)
	for i, t := range idList {
		ids += easygo.AnytoA(t)
		if i < count-1 {
			ids += ","
		}
	}

	msg := fmt.Sprintf("批量解冻钻石帐户: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)

	return easygo.EmptyMsg
}
