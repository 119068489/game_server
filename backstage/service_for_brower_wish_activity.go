package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

// 获取活动奖池管理
func (self *cls4) RpcQueryWishActPool(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishActPool, %+v：", reqMsg)
	list, count := QueryWishActPool(reqMsg)
	ret := &brower_backstage.WishActPoolList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 新增/更新活动奖池
func (self *cls4) RpcUpdateWishActPool(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.WishActPool) easygo.IMessage {
	logs.Info("RpcUpdateWishActPool, %+v：", reqMsg)
	ids := reqMsg.GetBoxIds()
	msg := fmt.Sprintf("更新活动奖池,%s", reqMsg.GetName())
	if reqMsg.GetId() == 0 {
		msg = fmt.Sprintf("新增活动奖池,%s", reqMsg.GetName())
		if GetAllActCfgIsOnline() {
			return easygo.NewFailMsg("活动已经开始无法进行删除")
		}
		if for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL, bson.M{}) >= 2 {
			return easygo.NewFailMsg("奖池不能超过2个")
		}
	} else {

		pool := GetWishActPool(reqMsg.GetId())
		havaIds := pool.GetBoxIds()
		havaIdM := make(map[int64]int64)
		for _, id := range havaIds {
			havaIdM[id] = id
		}

		for _, id := range ids {
			if havaIdM[id] > 0 {
				havaIdM[id] = 0
			}
		}

		for _, v := range havaIdM {
			if v > 0 {
				if GetAllActCfgIsOnline() {
					return easygo.NewFailMsg("活动已经开始无法进行删除")
				}
			}
		}
	}

	for _, id := range ids {
		pool := GetWishActPoolByBoxId(id)
		if pool.GetId() > 0 && pool.GetId() != reqMsg.GetId() {
			return easygo.NewFailMsg(fmt.Sprintf("该id为%d的盲盒已经存在于%s奖池", id, pool.GetName()))
		}
	}

	reqMsg.BoxNum = easygo.NewInt32(len(ids))

	UpdateWishActPool(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 删除活动奖池
func (self *cls4) RpcDeleteWishActPool(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	logs.Info("RpcDeleteWishActPool, %+v：", reqMsg)
	if GetAllActCfgIsOnline() {
		return easygo.NewFailMsg("活动已经开始无法进行删除")
	}

	ids := reqMsg.GetIds64()
	poolcount := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL, bson.M{})
	if poolcount-len(ids) < 2 {
		return easygo.NewFailMsg("奖池不能少于2个")
	}
	DeleteWishActPool(ids)
	DeleteWishActPoolRuleByPools(ids)
	msg := fmt.Sprintf("删除活动奖池,%v", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 获取活动奖池详情
func (self *cls4) RpcQueryWishActPoolDetail(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishActPooDetail, %+v：", reqMsg)
	var ids []int64

	//已有的盲盒列表
	var haveWishBox []*share_message.WishBox
	var activeWishBox []*share_message.WishBox

	//可用的盲盒列表
	activeBoxes := GetActiveWishBoxList()
	var list []*brower_backstage.WishActPoolItem
	if reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1: // 盲盒id
			id := easygo.StringToInt64noErr(reqMsg.GetKeyword())
			searchBoxes := GetWishBoxList([]int64{id})
			searchBox := searchBoxes[0]
			searchOne := &brower_backstage.WishActPoolItem{
				BoxId:  easygo.NewInt64(searchBox.GetId()),
				Name:   easygo.NewString(searchBox.GetName()),
				Price:  easygo.NewInt64(searchBox.GetPrice()),
				IsHave: easygo.NewBool(true),
			}
			for _, v := range activeBoxes {
				if v.GetId() == searchBoxes[0].GetId() {
					searchOne.IsHave = easygo.NewBool(false)
					break
				}
			}
			list = append(list, searchOne)
		default:

		}
	} else {
		//当前可用的盲盒列表
		activeWishBox = activeBoxes
		if reqMsg.GetId() > 0 {
			logs.Info("reqMsg.GetId():", reqMsg.GetId())
			act := GetWishActPool(reqMsg.GetId())
			ids = act.GetBoxIds()
			haveWishBox = GetWishBoxList(ids)
		}

	}

	//if len(ids) > 10 && reqMsg.GetPageSize() > 0 {
	//	pageSize, curPage := SetMgoPage(reqMsg.GetPageSize(), reqMsg.GetCurPage())
	//	ids = PaginateForInt64(ids, curPage*pageSize, pageSize)
	//}

	havaIdM := make(map[int64]int64)

	for _, id := range ids {
		havaIdM[id] = id
	}

	/*for _, v := range boxs {
		one := &brower_backstage.WishActPoolItem{
			BoxId: easygo.NewInt64(v.GetId()),
			Name:  easygo.NewString(v.GetName()),
			Price: easygo.NewInt64(v.GetPrice()),
		}
		if havaIdM[v.GetId()] > 0 {
			one.IsHave = easygo.NewBool(true)
		} else {
			one.IsHave = easygo.NewBool(false)
		}
		list = append(list, one)
	}*/
	if len(activeWishBox) > 0 {
		for _, v := range activeWishBox {
			one := &brower_backstage.WishActPoolItem{
				BoxId:  easygo.NewInt64(v.GetId()),
				Name:   easygo.NewString(v.GetName()),
				Price:  easygo.NewInt64(v.GetPrice()),
				IsHave: easygo.NewBool(false),
			}
			list = append(list, one)
		}
	}

	if len(haveWishBox) > 0 {
		for _, v := range haveWishBox {
			haveOne := &brower_backstage.WishActPoolItem{
				BoxId:  easygo.NewInt64(v.GetId()),
				Name:   easygo.NewString(v.GetName()),
				Price:  easygo.NewInt64(v.GetPrice()),
				IsHave: easygo.NewBool(true),
			}
			list = append(list, haveOne)
		}
	}

	//logs.Info("list---", list)
	ret := &brower_backstage.WishActPoolDetail{
		List:      list,
		PageCount: easygo.NewInt32(len(list)),
	}

	return ret
}

// 获取活动奖池列表键值对
func (self *cls4) RpcGetWishActPoolTypeKvs(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {
	pools := GetWishActPoolAll()
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

// 活动设置
// 获取活动设置
func (self *cls4) RpcGetWishActCfg(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcGetWishActCfg, %+v：", reqMsg)
	t := reqMsg.GetListType()
	var data *share_message.Activity
	curTime := time.Now().Unix()
	if t == 0 {
		luckyActivities := GetWishActCfgs([]int32{3, 4, 5})
		status := 1 //0开启，1关闭
		for _, v := range luckyActivities {
			if v.GetStatus() == 0 && v.GetEndTime() > curTime {
				status = 0
				break
			}
		}
		data = &share_message.Activity{
			Status: easygo.NewInt32(status),
		}
	} else {
		data = GetWishActCfg(t)
		if data.GetEndTime() <= curTime {
			data.Status = easygo.NewInt32(1)
			UpdateWishActCfg(data)
			data = GetWishActCfg(t)
		}
	}
	return data
}

//添加白名单
func (self *cls4) RpcAddWishAllowList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.AddWishAllowListReq) easygo.IMessage {
	logs.Info("RpcAddWishAllowList, %+v：", reqMsg)
	remark := reqMsg.GetRemark()
	accountsMap := reqMsg.GetAccounts()
	playerBaseList := QueryplayerlistByAccounts(accountsMap)
	if len(playerBaseList) != len(accountsMap) {
		return easygo.NewFailMsg("柠檬号有误，请检查后，重新添加")
	}

	pacts := ""
	for o, p := range playerBaseList {
		if p.GetTypes() < 1 {
			if len(playerBaseList) == o+1 {
				pacts += p.GetAccount()
			} else {
				pacts += p.GetAccount() + ","
			}
		}
	}
	if pacts != "" {
		return easygo.NewFailMsg("非运营号不可加入白名单:" + pacts)
	}

	count := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WHITE, bson.M{"Account": bson.M{"$in": accountsMap}})
	if count > 0 {
		return easygo.NewFailMsg("白名单中已包含该柠檬号！")
	}
	AddWishAllowList(playerBaseList, remark)
	return easygo.EmptyMsg
}

//删除白名单
func (self *cls4) RpcDeleteWishAllowList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	logs.Info("RpcDeleteWishAllowList, %+v：", reqMsg)
	Ids := reqMsg.GetIds64()
	var toDeleteIds []int64
	for _, wishPlayerId := range Ids {
		wishPlayer := GetWishPlayerInfo(wishPlayerId)
		wishBoxlist := GetWishBoxListByGuardianIds(wishPlayerId)
		//如果删除的白名单中有盲盒守护者
		if len(wishBoxlist) > 0 {
			req := &h5_wish.BackstageSetGuardianReq{
				Account:  easygo.NewString(wishPlayer.GetAccount()),
				Channel:  easygo.NewInt32(1001),
				NickName: easygo.NewString(wishPlayer.GetNickName()),
				HeadUrl:  easygo.NewString(wishPlayer.GetHeadUrl()),
				PlayerId: easygo.NewInt64(wishPlayer.GetPlayerId()),
				Token:    easygo.NewString(""),
				OpType:   easygo.NewInt32(2),
			}
			var toDeleteBool bool
			for _, wishBox := range wishBoxlist {
				if wishBox.GetGuardianId() == wishPlayer.GetId() {
					req.BoxId = easygo.NewInt64(wishBox.GetId())
					err := CallRpcBackstageSetGuardian(req)
					if err == nil {
						toDeleteBool = true
						continue
					} else {
						toDeleteBool = false
						logs.Warn("通知许愿池服务器更新盲盒守护者失败，rpc=RpcBackstageSetGuardian，req=%+v, err=%v", req, err.GetReason())
						break
					}

				}
			}
			if toDeleteBool {
				toDeleteIds = append(toDeleteIds, wishPlayerId)
			}
		} else {
			toDeleteIds = append(toDeleteIds, wishPlayerId)
		}
	}

	//logs.Info("toDeleteIds----", toDeleteIds)
	if len(toDeleteIds) > 0 {
		findBson := bson.M{"_id": bson.M{"$in": toDeleteIds}}
		//删除指定的白名单
		for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WHITE, findBson)
	}
	return easygo.EmptyMsg
}

//白名单列表
func (self *cls4) RpcWishAllowList(ep IBrowerEndpoint, ctx interface{}, reqMsg interface{}) easygo.IMessage {
	logs.Info("RpcWishAllowList, %+v：", reqMsg)
	wishPlayers := GetAllWishWhite()
	return &brower_backstage.WishAllowListResp{
		List: wishPlayers,
	}
}

// 更新活动设置
func (self *cls4) RpcUpdateWishActCfg(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.Activity) easygo.IMessage {
	logs.Info("RpcUpdateWishActCfg, %+v：", reqMsg)
	UpdateWishActCfg(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, "更新活动设置")
	return easygo.EmptyMsg
}

// 许愿池活动规则管理
// 累计天数规则管理
// 获取累计天数规则
func (self *cls4) RpcQueryWishActPoolRuleDay(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishActPoolRuleDay, %+v：", reqMsg)
	list, count := QueryWishActPoolRule(reqMsg, 2)

	var retList []*brower_backstage.WishActPoolRule

	for _, v := range list {
		one := &brower_backstage.WishActPoolRule{
			Id:            easygo.NewInt64(v.GetId()),
			WishActPoolId: easygo.NewInt64(v.GetWishActPoolId()),
			Name:          easygo.NewString(""),
			Key:           easygo.NewInt32(v.GetKey()),
			Diamond:       easygo.NewInt64(v.GetDiamond()),
		}
		pool := GetWishActPool(v.GetWishActPoolId())
		one.Name = easygo.NewString(pool.GetName())
		retList = append(retList, one)
	}

	ret := &brower_backstage.WishActPoolRuleList{
		List:      retList,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 新增累计天数规则
func (self *cls4) RpcAddWishActPoolRuleDay(ep IBrowerEndpoint, user *share_message.Manager, req *brower_backstage.AddWishActPoolRuleRequest) easygo.IMessage {
	logs.Info("RpcAddWishActPoolRuleDay, %+v：", req)

	list := req.GetList()

	for _, v := range list {
		if v.GetKey() == 0 && v.GetDiamond() == 0 {
			continue
		}

		if (v.GetKey() == 0 && v.GetDiamond() != 0) || (v.GetKey() != 0 && v.GetDiamond() == 0) {
			return easygo.NewFailMsg("参数错误")
		}

		id := v.GetWishActPoolId()
		rule := GetWishActPoolRuleByPoolIdAndKey(id, v.GetKey(), for_game.WISH_ACTIVITY_DATA_TYPE_2)
		if rule.GetId() > 0 {
			if GetActCfgIsOnline(for_game.WISH_ACT_DAY) {
				return easygo.NewFailMsg("活动已经开始无法进行操作")
			}
			rule.Diamond = easygo.NewInt64(v.GetDiamond())
			UpdateWishActPoolRule(rule)
		} else {
			one := &share_message.WishActPoolRule{
				WishActPoolId: easygo.NewInt64(id),
				Key:           easygo.NewInt32(v.GetKey()),
				Diamond:       easygo.NewInt64(v.GetDiamond()),
				AwardType:     easygo.NewInt32(1),
				Type:          easygo.NewInt32(2),
			}

			if for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL_RULE, bson.M{"WishActPoolId": id, "Type": 2}) >= 4 {
				return easygo.NewFailMsg("累计天数规则不能大于4个")
			}

			UpdateWishActPoolRule(one)
		}

	}

	msg := fmt.Sprintf("新增累计天数规则,奖池id")
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 更新累计天数规则
func (self *cls4) RpcUpdateWishActPoolRuleDay(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishActPoolRule) easygo.IMessage {
	logs.Info("RpcUpdateWishActPoolRuleDay, %+v：", reqMsg)
	if GetActCfgIsOnline(for_game.WISH_ACT_DAY) {
		return easygo.NewFailMsg("活动已经开始无法进行操作")
	}
	id := reqMsg.GetWishActPoolId()
	rule := GetWishActPoolRuleByPoolIdAndKey(id, reqMsg.GetKey(), for_game.WISH_ACTIVITY_DATA_TYPE_2)
	if rule.GetId() > 0 {
		rule.Diamond = easygo.NewInt64(reqMsg.GetDiamond())
		UpdateWishActPoolRule(rule)
	} else {
		r := GetWishActPoolRule(reqMsg.GetId())
		r.Key = easygo.NewInt32(reqMsg.GetKey())
		r.Diamond = easygo.NewInt64(reqMsg.GetDiamond())
		UpdateWishActPoolRule(r)
	}
	msg := fmt.Sprintf("更新累计天数规则,奖池id%d", reqMsg.GetWishActPoolId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 删除累计天数规则
func (self *cls4) RpcDeleteWishActPoolRuleDay(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	logs.Info("RpcDeleteWishActPoolRuleDay, %+v：", reqMsg)
	if GetActCfgIsOnline(for_game.WISH_ACT_DAY) {
		return easygo.NewFailMsg("活动已经开始无法进行删除")
	}

	ids := reqMsg.GetIds64()
	DeleteWishActPoolRule(ids)
	msg := fmt.Sprintf("删除累计天数规则,%v", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 许愿池活动规则管理
// 累计次数规则管理
// 获取累计次数规则
func (self *cls4) RpcQueryWishActPoolRuleCount(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishActPoolRuleCount, %+v：", reqMsg)
	list, count := QueryWishActPoolRule(reqMsg, 1)

	var retList []*brower_backstage.WishActPoolRule

	for _, v := range list {
		one := &brower_backstage.WishActPoolRule{
			Id:            easygo.NewInt64(v.GetId()),
			WishActPoolId: easygo.NewInt64(v.GetWishActPoolId()),
			Name:          easygo.NewString(""),
			Key:           easygo.NewInt32(v.GetKey()),
			Diamond:       easygo.NewInt64(v.GetDiamond()),
		}
		pool := GetWishActPool(v.GetWishActPoolId())
		one.Name = easygo.NewString(pool.GetName())
		retList = append(retList, one)
	}

	ret := &brower_backstage.WishActPoolRuleList{
		List:      retList,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 新增累计次数规则
func (self *cls4) RpcAddWishActPoolRuleCount(ep IBrowerEndpoint, user *share_message.Manager, req *brower_backstage.AddWishActPoolRuleRequest) easygo.IMessage {
	logs.Info("RpcAddWishActPoolRuleCount, %+v：", req)

	list := req.GetList()

	for _, v := range list {
		if v.GetKey() == 0 && v.GetDiamond() == 0 {
			continue
		}

		if (v.GetKey() == 0 && v.GetDiamond() != 0) || (v.GetKey() != 0 && v.GetDiamond() == 0) {
			return easygo.NewFailMsg("参数错误")
		}

		id := v.GetWishActPoolId()
		rule := GetWishActPoolRuleByPoolIdAndKey(id, v.GetKey(), for_game.WISH_ACTIVITY_DATA_TYPE_1)
		if rule.GetId() > 0 {
			if GetActCfgIsOnline(for_game.WISH_ACT_COUNT) {
				return easygo.NewFailMsg("活动已经开始无法进行操作")
			}
			rule.Diamond = easygo.NewInt64(v.GetDiamond())
			UpdateWishActPoolRule(rule)
		} else {

			one := &share_message.WishActPoolRule{
				WishActPoolId: easygo.NewInt64(id),
				Key:           easygo.NewInt32(v.GetKey()),
				Diamond:       easygo.NewInt64(v.GetDiamond()),
				AwardType:     easygo.NewInt32(1),
				Type:          easygo.NewInt32(1),
			}
			if for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL_RULE, bson.M{"WishActPoolId": id, "Type": 1}) >= 5 {
				return easygo.NewFailMsg("累计次数规则不能大于5个")
			}
			UpdateWishActPoolRule(one)
		}

	}

	msg := fmt.Sprintf("新增累计次数规则")
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 更新累计次数规则
func (self *cls4) RpcUpdateWishActPoolRuleCount(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishActPoolRule) easygo.IMessage {
	logs.Info("RpcUpdateWishActPoolRuleCount, %+v：", reqMsg)
	if GetActCfgIsOnline(for_game.WISH_ACT_COUNT) {
		return easygo.NewFailMsg("活动已经开始无法进行操作")
	}

	id := reqMsg.GetId()
	rule := GetWishActPoolRuleByPoolIdAndKey(id, reqMsg.GetKey(), for_game.WISH_ACTIVITY_DATA_TYPE_1)
	if rule.GetId() > 0 {
		rule.Diamond = easygo.NewInt64(reqMsg.GetDiamond())
		UpdateWishActPoolRule(rule)
	} else {
		r := GetWishActPoolRule(reqMsg.GetId())
		r.Key = easygo.NewInt32(reqMsg.GetKey())
		r.Diamond = easygo.NewInt64(reqMsg.GetDiamond())
		UpdateWishActPoolRule(r)
	}
	msg := fmt.Sprintf("更新累计次数规则,奖池id%d", reqMsg.GetWishActPoolId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 删除累计次数规则
func (self *cls4) RpcDeleteWishActPoolRuleCount(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	logs.Info("RpcDeleteWishActPoolRuleCount, %+v：", reqMsg)
	if GetActCfgIsOnline(for_game.WISH_ACT_COUNT) {
		return easygo.NewFailMsg("活动已经开始无法进行删除")
	}

	ids := reqMsg.GetIds64()
	DeleteWishActPoolRule(ids)
	msg := fmt.Sprintf("删除累计次数规则,%v", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 获取累计金额规则
func (self *cls4) RpcQueryWishActPoolRuleWeekMonth(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishActPoolRuleWeekMonth, %+v：", reqMsg)
	t := 0
	switch reqMsg.GetListType() {
	case 1:
		t = 3
	case 2:
		t = 4
	default:
	}
	reqMsg.ListType = easygo.NewInt32(0)
	list, count := QueryWishActPoolRuleWeekMonth(reqMsg, int32(t))

	//var retList []*brower_backstage.WishActPoolRule
	var retList2 []*brower_backstage.WishActPoolAwardItem
	for _, v := range list {
		switch v.GetAwardType() {
		case 1:
			one := &brower_backstage.WishActPoolAwardItem{
				Id:            easygo.NewInt64(v.GetId()),
				WishActPoolId: easygo.NewInt64(v.GetWishActPoolId()),
				Name:          easygo.NewString(""),
				Key:           easygo.NewInt32(v.GetKey()),
				AwardDiamond:  easygo.NewInt64(v.GetDiamond()),
				AwardType:     easygo.NewInt32(1),
			}
			switch v.GetType() {
			case 3:
				one.RuleType = easygo.NewInt32(1)
			case 4:
				one.RuleType = easygo.NewInt32(2)
			default:

			}
			pool := GetWishActPool(v.GetWishActPoolId())
			one.Name = easygo.NewString(pool.GetName())
			retList2 = append(retList2, one)
		case 2:
			one := &brower_backstage.WishActPoolAwardItem{
				Id:            easygo.NewInt64(v.GetId()),
				WishActPoolId: easygo.NewInt64(v.GetWishActPoolId()),
				Key:           easygo.NewInt32(v.GetKey()),
				WishItemId:    easygo.NewInt64(v.GetWishItemId()),
				AwardType:     easygo.NewInt32(2),
			}
			switch v.GetType() {
			case 3:
				one.RuleType = easygo.NewInt32(1)
			case 4:
				one.RuleType = easygo.NewInt32(2)
			default:

			}

			item := GetWishItemById(one.GetWishItemId())
			one.Name = easygo.NewString(item.GetName())
			one.Icon = easygo.NewString(item.GetIcon())
			one.Money = easygo.NewInt64(item.GetPrice())
			one.Diamond = easygo.NewInt64(item.GetDiamond())
			retList2 = append(retList2, one)

		default:

		}

	}

	ret := &brower_backstage.WishActPoolRuleList{
		ItemList:  retList2,
		PageCount: easygo.NewInt32(count),
	}
	return ret
}

// 新增累计金额规则
func (self *cls4) RpcAddWishActPoolRuleWeekMonth(ep IBrowerEndpoint, user *share_message.Manager, req *brower_backstage.AddWishActPoolRuleRequest) easygo.IMessage {
	logs.Info("RpcAddWishActPoolRuleWeekMonth, %+v：", req)

	at := req.GetAwardType() // 1、钻石奖励 2、实物奖励
	rt := req.GetRuleType()
	switch rt {
	case 1:
		rt = for_game.WISH_ACTIVITY_DATA_TYPE_3 // 周榜
	case 2:
		rt = for_game.WISH_ACTIVITY_DATA_TYPE_4 // 月榜
	default:

	}
	switch at {
	case 1:
		list := req.GetWeekMonthList()
		for _, v := range list {

			if v.GetKey() == 0 && v.GetDiamond() == 0 && v.GetId() != 0 {
				DeleteWishActPoolRule([]int64{v.GetId()})
				continue
			}

			if (v.GetKey() == 0 && v.GetAwardDiamond() != 0) || (v.GetKey() != 0 && v.GetAwardDiamond() == 0) {
				return easygo.NewFailMsg("参数错误")
			}

			id := v.GetWishActPoolId()
			rule := GetWishActPoolRuleByPoolIdAndKey(id, v.GetKey(), rt)
			if rule.GetId() > 0 {
				if GetActCfgIsOnline(for_game.WISH_ACT_WEEK_MONTH) && rule.GetDiamond() != v.GetAwardDiamond() {
					return easygo.NewFailMsg("活动已经开始无法进行操作")
				}
				// 相同排名不能同时存在两种奖励
				if rule.GetAwardType() != 1 {
					return easygo.NewFailMsg(fmt.Sprintf("已存在排名%d的实物奖励", v.GetKey()))
				}
				rule.Diamond = easygo.NewInt64(v.GetAwardDiamond())
				UpdateWishActPoolRule(rule)
			} else {

				one := &share_message.WishActPoolRule{
					WishActPoolId: easygo.NewInt64(id),
					Key:           easygo.NewInt32(v.GetKey()),
					Diamond:       easygo.NewInt64(v.GetAwardDiamond()),
					AwardType:     easygo.NewInt32(1),
					Type:          easygo.NewInt32(rt),
				}
				UpdateWishActPoolRule(one)
			}
		}
	case 2:
		list := req.GetWeekMonthList()
		for _, v := range list {
			if v.GetKey() == 0 && v.GetWishItemId() == 0 && v.GetId() != 0 {
				DeleteWishActPoolRule([]int64{v.GetId()})
				continue
			}

			if (v.GetKey() == 0 && v.GetName() != "") || (v.GetKey() != 0 && v.GetName() == "") {
				return easygo.NewFailMsg("参数错误")
			}

			id := v.GetWishActPoolId()
			rule := GetWishActPoolRuleByPoolIdAndKey(id, v.GetKey(), rt)
			if rule.GetId() > 0 {
				item := GetWishItemById(rule.GetWishItemId())
				// 判断实物信息是否有更改，活动开启期间，不允许更新
				if GetActCfgIsOnline(for_game.WISH_ACT_WEEK_MONTH) && (item.GetName() != v.GetName() || item.GetDiamond() != v.GetDiamond() || item.GetPrice() != v.GetMoney()) {
					return easygo.NewFailMsg("活动已经开始无法进行操作")
				}
				// 相同排名不能同时存在两种奖励
				if rule.GetAwardType() != 2 {
					return easygo.NewFailMsg(fmt.Sprintf("已存在排名%d的钻石奖励", v.GetKey()))
				}

				item.Name = easygo.NewString(v.GetName())
				item.Icon = easygo.NewString(v.GetIcon())
				item.Price = easygo.NewInt64(v.GetMoney())
				item.Diamond = easygo.NewInt64(v.GetDiamond())
				UpdateWishItem(item)
			} else {
				item := &share_message.WishItem{
					Name:    easygo.NewString(v.GetName()),
					Icon:    easygo.NewString(v.GetIcon()),
					Price:   easygo.NewInt64(v.GetMoney()),
					Diamond: easygo.NewInt64(v.GetDiamond()),
					Status:  easygo.NewInt32(0),
					UseType: easygo.NewInt32(1),
				}
				itemId := AddWishItem(item)
				one := &share_message.WishActPoolRule{
					WishActPoolId: easygo.NewInt64(id),
					Key:           easygo.NewInt32(v.GetKey()),
					AwardType:     easygo.NewInt32(2),
					Type:          easygo.NewInt32(rt),
				}
				one.WishItemId = easygo.NewInt64(itemId)
				UpdateWishActPoolRule(one)
			}
		}
	default:

	}

	msg := fmt.Sprintf("新增累计金额规则")
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 更新累计金额规则-钻石奖励
func (self *cls4) RpcUpdateWishActPoolRuleWeekMonth(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishActPoolRule) easygo.IMessage {
	logs.Info("RpcUpdateWishActPoolRuleWeekMonth, %+v：", reqMsg)
	if GetActCfgIsOnline(for_game.WISH_ACT_WEEK_MONTH) {
		return easygo.NewFailMsg("活动已经开始无法进行操作")
	}

	id := reqMsg.GetId()
	//one := &share_message.WishActPoolRule{
	//	Id:      easygo.NewInt64(id),
	//	Key:     easygo.NewInt32(reqMsg.GetKey()),
	//	Diamond: easygo.NewInt64(reqMsg.GetDiamond()),
	//}
	//UpdateWishActPoolRule(one)

	rule := GetWishActPoolRuleByPoolIdAndKey(id, reqMsg.GetKey(), for_game.WISH_ACTIVITY_DATA_TYPE_2)
	if rule.GetId() > 0 {
		rule.Diamond = easygo.NewInt64(reqMsg.GetDiamond())
		UpdateWishActPoolRule(rule)
	} else {
		return easygo.NewFailMsg("数据有误")
	}

	msg := fmt.Sprintf("更新累计次数规则,奖池id%d", reqMsg.GetWishActPoolId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 更新累计金额规则-钻石奖励/实物奖励
func (self *cls4) RpcUpdateWishActPoolRuleItemWeekMonth(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.WishActPoolAwardItem) easygo.IMessage {
	logs.Info("RpcUpdateWishActPoolRuleItemWeekMonth, %+v：", reqMsg)
	if GetActCfgIsOnline(for_game.WISH_ACT_WEEK_MONTH) {
		return easygo.NewFailMsg("活动已经开始无法进行操作")
	}

	id := reqMsg.GetWishActPoolId()
	rt := reqMsg.GetRuleType()
	switch rt {
	case 1:
		rt = for_game.WISH_ACTIVITY_DATA_TYPE_3 // 周榜
	case 2:
		rt = for_game.WISH_ACTIVITY_DATA_TYPE_4 // 月榜
	default:

	}
	rule := GetWishActPoolRuleByPoolIdAndKey(id, reqMsg.GetKey(), rt)
	if rule.GetId() > 0 {
		switch reqMsg.GetAwardType() {
		case 1: //  钻石奖励
			rule.Diamond = easygo.NewInt64(reqMsg.GetDiamond())
			UpdateWishActPoolRule(rule)
		case 2: //  实物奖励
			item := GetWishItemById(rule.GetWishItemId())
			item.Name = easygo.NewString(reqMsg.GetName())
			item.Icon = easygo.NewString(reqMsg.GetIcon())
			item.Price = easygo.NewInt64(reqMsg.GetMoney())
			item.Diamond = easygo.NewInt64(reqMsg.GetDiamond())
			UpdateWishItem(item)
		default:

		}

	} else {
		return easygo.NewFailMsg("数据有误")
	}
	msg := fmt.Sprintf("更新累计次数规则,奖池id%d", reqMsg.GetWishActPoolId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

// 删除累计金额规则
func (self *cls4) RpcDeleteWishActPoolRuleItemWeekMonth(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	logs.Info("RpcDeleteWishActPoolRuleItemWeekMonth, %+v：", reqMsg)
	if GetActCfgIsOnline(for_game.WISH_ACT_WEEK_MONTH) {
		return easygo.NewFailMsg("活动已经开始无法进行删除")
	}

	ids := reqMsg.GetIds64()
	DeleteWishActPoolRule(ids)
	msg := fmt.Sprintf("删除累计金额规则,%v", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

//充值活动配置
func (self *cls4) RpcWishCoinRechargeActivityCfgList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcWishCoinRechargeActivityCfgList, %+v：", reqMsg)
	sort := []string{"_id"}
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_COIN_RECHARGE_ACT_CFG, bson.M{}, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.WishCoinRechargeActivityCfg
	for _, li := range lis {
		one := &share_message.WishCoinRechargeActivityCfg{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}

	return &brower_backstage.WishCoinRechargeActivityCfgRes{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//充值活动配置修改新增
func (s *cls4) RpcWishCoinRechargeActivityCfgUpdate(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.WishCoinRechargeActivityCfg) easygo.IMessage {
	logs.Info("RpcWishCoinRechargeActivityCfgUpdate, %+v：", reqMsg)
	if reqMsg.Id == nil || reqMsg.GetId() < 1 {
		return easygo.NewFailMsg("硬币数量不能小于1")
	}

	if reqMsg.Amount == nil || reqMsg.GetAmount() < 1 {
		return easygo.NewFailMsg("人民币价格不能小于1")
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_COIN_RECHARGE_ACT_CFG, queryBson, updateBson, true)
	msg := "修改充值活动配置"
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)
	return easygo.EmptyMsg
}

//充值活动配置删除
func (s *cls4) RpcWishCoinRechargeActivityCfgDel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}

	findBson := bson.M{"_id": bson.M{"$in": idList}}
	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_COIN_RECHARGE_ACT_CFG, findBson)

	var ids string
	count := len(idList)
	for i, t := range idList {
		ids += easygo.AnytoA(t)
		if i < count-1 {
			ids += ","
		}
	}

	msg := fmt.Sprintf("批量删除充值活动配置: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WISH_MANAGE, msg)

	return easygo.EmptyMsg
}

// 许愿池活动用户记录-活动用户记录列表
func (self *cls4) RpcQueryWishActPlayerRecordList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishActPlayerRecordList, %+v：", reqMsg)

	list, count := QueryActPlayerInfoList(reqMsg)

	var ids []int64
	var retList []*brower_backstage.WishActPlayerRecord
	userM := make(map[int64]*brower_backstage.WishActPlayerRecord)
	//获取所有奖池
	wishPools, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL, bson.M{}, 0, 0)
	zeroTime := easygo.Get0ClockTimestamp(time.Now().Unix()) * 1000 //wish_day_activity_log存的是毫秒
	for i := range list {
		actInfo := list[i]
		playerId := actInfo.GetPlayerId()
		ids = append(ids, playerId)

		one := &brower_backstage.WishActPlayerRecord{}

		userM[list[i].GetPlayerId()] = one

		u := for_game.GetWishPlayerByPid(actInfo.GetPlayerId())
		u2 := GetPlayerInfoByPid(u.GetPlayerId())
		one.UserId = easygo.NewInt64(list[i].GetPlayerId())
		one.UserAccount = easygo.NewString(u2.GetAccount())
		//activityDatas := actInfo.GetData()
		for _, li := range wishPools {
			actPool := &share_message.WishActPool{}
			for_game.StructToOtherStruct(li, actPool)

			drawTotal, dayTotal := for_game.SumContinuedCount(playerId, actPool.GetId())
			countBson := bson.M{"PlayerId": playerId, "WishActPoolId": actPool.GetId(), "CreateTime": bson.M{"$in": []int64{zeroTime, zeroTime - 24*3600*1000}}}
			count := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_DAY_ACTIVITY_LOG, countBson)
			//断抽
			if count == 0 {
				dayTotal = 0
			}
			//logs.Info("playerId: %d, actPool: %d, drawTotal: %d, dayTotal: %d", u.GetPlayerId(), actPool.GetId(),drawTotal, dayTotal)
			unit := &brower_backstage.PoolData{
				Id:        easygo.NewInt64(actPool.GetId()),
				Name:      easygo.NewString(actPool.GetName()),
				DrawTotal: easygo.NewInt32(drawTotal),
				DayTotal:  easygo.NewInt32(dayTotal),
			}
			one.Data = append(one.Data, unit)
		}
		// 活动拓展信息

		//actList := actInfo.GetData()
		/*for _, act := range actList {
			switch act.GetType() {
			case 1: //  次数
				one.DrawTotal = easygo.NewInt32(act.GetValue())
			case 2: //  天数
				one.DayTotal = easygo.NewInt32(act.GetValue())
			default:

			}
		}*/
		one.DrawDiamondTotal = easygo.NewInt64(actInfo.GetDiamond())
		retList = append(retList, one)
	}

	// 领奖次数
	awardTotal := GetActReceiveAwardCountByUserIds(ids)
	for _, v := range awardTotal {
		if mv, ok := userM[v.GetUserId()]; ok {
			mv.AwardTotal = easygo.NewInt32(v.GetAwardTotal())
		}
	}

	ret := &brower_backstage.WishActPlayerRecordList{
		List:      retList,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 许愿池活动用户记录-获奖记录列表
func (self *cls4) RpcQueryWishActPlayerWinRecordList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishActPlayerWinRecordList, %+v：", reqMsg)
	list, count := QueryActWinRecordList(reqMsg)

	var retList []*brower_backstage.WishActPlayerWinRecord

	for _, v := range list {
		one := &brower_backstage.WishActPlayerWinRecord{
			Id:         easygo.NewInt64(v.GetId()),
			Status:     easygo.NewInt32(v.GetStatus()),
			CreateTime: easygo.NewInt64(v.GetCreateTime()),
			WinTime:    easygo.NewInt64(v.GetFinishTime()),
		}

		// 获奖类型处理
		typeN := ""
		//rule := GetWishActPoolRule(v.GetActType())
		switch v.GetType() {
		case 1:
			typeN = fmt.Sprintf("累计次数达到%d次", v.GetActType())
		case 2:
			typeN = fmt.Sprintf("累计天数达到%d天", v.GetActType())
		case 3:
			typeN = fmt.Sprintf("周榜单第%d天", v.GetActType())
		case 4:
			typeN = fmt.Sprintf("月榜单第%d天", v.GetActType())
		default:

		}

		// 备注处理
		note := ""
		switch v.GetPrizeType() {
		case 1: // 钻石奖励
			note = fmt.Sprintf("钻石*%d", v.GetPrizeValue())
		case 2: // 实物奖励
			item := GetWishItemById(v.GetPrizeValue())
			note = item.GetName()
		default:

		}
		one.TypeName = easygo.NewString(typeN)
		one.Note = easygo.NewString(note)
		retList = append(retList, one)
	}

	ret := &brower_backstage.WishActPlayerWinRecordList{
		List:      retList,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

// 许愿池活动用户记录-抽取记录列表
func (self *cls4) RpcQueryWishActPlayerDrawRecordList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishActPlayerDrawRecordList, %+v：", reqMsg)
	list, count := QueryActDrawRecordListByPid(reqMsg)

	var retList []*brower_backstage.WishActPlayerDrawRecord

	for _, v := range list {
		one := &brower_backstage.WishActPlayerDrawRecord{
			Price:      easygo.NewInt64(v.GetDarePrice()),
			CreateTime: easygo.NewInt64(v.GetCreateTime()),
			BoxId:      easygo.NewInt64(v.GetWishBoxId()),
		}
		box := GetWishBoxById(v.GetWishBoxId())
		one.ActBoxName = easygo.NewString(box.GetName())
		pool := GetWishActPoolByBoxId(v.GetWishBoxId())
		one.ActPoolName = easygo.NewString(pool.GetName())
		retList = append(retList, one)
	}

	ret := &brower_backstage.WishActPlayerDrawRecordList{
		List:      retList,
		PageCount: easygo.NewInt32(count),
	}

	return ret
}

//许愿池活动日志查询
func (self *cls4) RpcQueryWishActivityPrizeLog(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcQueryWishActivityPrizeLog, %+v：", reqMsg)
	findBson := bson.M{}
	sort := []string{"-_id"}
	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACTIVITY_PRIZE_LOG, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.WishActivityPrizeLog
	for _, l := range lis {
		one := &share_message.WishActivityPrizeLog{}
		for_game.StructToOtherStruct(l, one)

		u := for_game.GetWishPlayerByPid(one.GetPlayerId())
		playerbase := GetPlayerInfoByPid(u.GetPlayerId())
		//one.ActTypeTitle = easygo.NewString(actTypeNameMap[one.GetActType()])
		var name string
		//1、次数 2、天数 3、周排名 4、月排名
		switch one.GetType() {
		case 1:
			name = "累计次数达到%d次"
		case 2:
			name = "累计天数达到%d天"
		case 3:
			name = "周榜单第%d名"
		case 4:
			name = "月榜单第%d名"
		}
		title := fmt.Sprintf(name, one.GetActType())
		one.ActTypeTitle = easygo.NewString(title)
		one.PlayerAccount = easygo.NewString(playerbase.GetAccount())
		list = append(list, one)
	}
	return &brower_backstage.QueryWishActivityPrizeLogRes{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

// 许愿池活动报表
func (self *cls4) RpcWishPoolActivityReportList(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcWishPoolActivityReportList, %+v：", reqMsg)
	list, count := QueryWishPoolActivityReportList(reqMsg)
	ret := &brower_backstage.WishPoolActivityReportResp{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return ret
}
