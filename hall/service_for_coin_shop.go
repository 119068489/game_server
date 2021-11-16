// 大厅服务器为[游戏客户端]提供的服务

package hall

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"
	"sort"

	"github.com/astaxie/beego/logs"
)

const (
	//虚拟商场提示
	COIN_RECHARGE_SUCCESS int32 = 1 //成功
	COIN_SHOP_FAIL_2      int32 = 2 //钱不足
	COIN_SHOP_FAIL_3      int32 = 3 //商品不存在
	COIN_SHOP_FAIL_4      int32 = 4 //
	COIN_SHOP_FAIL_5      int32 = 5 //
	COIN_SHOP_FAIL_6      int32 = 6 //
	COIN_SHOP_FAIL_7      int32 = 7 //道具使用重复
	COIN_SHOP_FAIL_8      int32 = 8 //商品已下架
	COIN_SHOP_FAIL_9      int32 = 9 //
	COIN_SHOP_FAIL_10     int32 = 10
)

// 硬币商场配置道具请求
func (self *cls1) RpcGetPropsItems(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcPropsItemsReq", who.GetPlayerId(), common)
	items := for_game.GetPropsItemsCfg()
	//logs.Info("返回商场道具:", items)
	return items
}

// 硬币充值商场请求
func (self *cls1) RpcGetCoinRechargeList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CoinRechargeList, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetCoinRechargeList", who.GetPlayerId(), common)
	items := for_game.GetCoinRechargeCfg(reqMsg.GetWay(), who.GetPlayerId())
	return items
}

//获取硬币商场指定类型商品
func (self *cls1) RpcGetCoinShopList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CoinShopList, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetCoinShopList", reqMsg, who.GetPlayerId(), common)
	items := for_game.GetCoinShopList(reqMsg.GetType())
	for _, li := range items {
		playermgr := for_game.GetRedisPlayerBase(who.GetPlayerId())
		if playermgr != nil && playermgr.GetDeviceType() == 1 && li.GetIosCoin() > 0 {
			li.Coin = easygo.NewInt64(li.GetIosCoin()) //如果设置了IOS价格，ios用户返回IOS价格给前端
		}
	}
	reqMsg.Items = items
	return reqMsg
}

//兑换硬币
func (self *cls1) RpcCoinRecharge(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CoinRechargeReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcCoinRecharge", reqMsg, who.GetPlayerId(), common)
	item := for_game.GetCoinRecharge(reqMsg.GetId())
	resp := &client_hall.CoinRechargeResp{
		Id:     easygo.NewInt64(reqMsg.GetId()),
		Result: easygo.NewInt32(COIN_RECHARGE_SUCCESS),
	}
	//青少年保护模式
	if who.GetYoungPassWord() != "" {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_8)
		return resp
	}
	//道具不存在
	if item == nil {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_3)
		return resp
	}
	player := for_game.GetRedisPlayerBase(who.GetPlayerId())
	//活动打折
	price := item.GetPrice()
	t := for_game.GetMillSecond()
	if t >= item.GetStartTime() && t < item.GetEndTime() {
		//活动期间
		price = item.GetDisPrice()
	}
	//零钱不足
	if player.GetGold() < price {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_2)
		return resp
	}
	if player.GetPayPassword() == "" {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_7)
		logs.Error("未设置零钱支付密码")
		return resp
	}
	if reqMsg.GetIsCheck() {
		resp.Result = easygo.NewInt32(COIN_RECHARGE_SUCCESS)
		return resp
	}
	//检查支付密码
	if player.GetPayPassword() != reqMsg.GetPassWord() {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_6)
		logs.Error("支付密码错误")
		return resp
	}
	//扣除零钱，并增加硬币
	orderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, for_game.GOLD_TYPE_EXCHANGE_COIN)
	reason := "兑换硬币消费"
	msg := &share_message.GoldExtendLog{
		OrderId: easygo.NewString(orderId),
		PayType: easygo.NewInt32(for_game.PAY_TYPE_GOLD),
		Title:   easygo.NewString(reason),
		Gold:    easygo.NewInt64(item.GetCoin()),
	}

	if NotifyAddGold(player.GetPlayerId(), -price, reason, for_game.GOLD_TYPE_EXCHANGE_COIN, msg) != "" {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_4) //服务器数据异常
		logs.Error("兑换扣除零钱失败")
		return resp
	}
	//增加硬币
	//每月首充额外赠送
	addCoin := int64(0)
	logObj := for_game.GetRedisCoinLogObj()
	if !logObj.CheckMonthRecharge(player.GetPlayerId(), item.GetId()) && for_game.CheckPlayerMonthRecharge(player.GetPlayerId(), item.GetId()) {
		addCoin = item.GetMonthFirst()
	}
	reason = "兑换获得"
	// orderId = for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, for_game.COIN_TYPE_EXCHANGE_IN)
	msg = &share_message.GoldExtendLog{
		RedPacketId: easygo.NewInt64(item.GetId()),
		Title:       easygo.NewString(reason),
		OrderId:     easygo.NewString(orderId),
	}
	if NotifyAddCoin(player.GetPlayerId(), item.GetCoin()+addCoin, reason, for_game.COIN_TYPE_EXCHANGE_IN, msg) != "" {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_5) //服务器数据异常
		logs.Error("兑换增加硬币失败")
		return resp
	}
	return resp
}

//购买虚拟商品
func (self *cls1) RpcBuyCoinItem(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.BuyCoinItem, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcBuyCoinItem", reqMsg, who.GetPlayerId(), common)
	playerId := who.GetPlayerId()
	item := for_game.GetCoinShopItem(reqMsg.GetId())
	//道具不存在
	if item == nil {
		reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_3)
		logs.Error("RpcBuyCoinItem 道具不存在!", reqMsg.GetId())
		reqMsg.Reason = easygo.NewString("道具不存在!")
		return reqMsg
	}
	//青少年保护模式
	if who.GetYoungPassWord() != "" {
		reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_9)
		reqMsg.Reason = easygo.NewString("青少年模式下无法购买")
		return reqMsg
	}
	if item.GetStatus() != for_game.COIN_PRODUCT_STATUS_UP {
		reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_10)
		logs.Error("RpcBuyCoinItem 商品已下架!", reqMsg.GetId())
		reqMsg.Reason = easygo.NewString("商品已下架!")
		return reqMsg
	}

	if item.GetSaleStartTime() > 0 && item.GetSaleEndTime() > 0 {
		nowTime := util.GetMilliTime()
		if item.GetSaleStartTime() > nowTime || item.GetSaleEndTime() < nowTime {
			reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_10)
			logs.Error("RpcBuyCoinItem 商品不在销售期!", reqMsg.GetId())
			reqMsg.Reason = easygo.NewString("商品不在销售期!")
			return reqMsg
		}
	}

	bagObj := for_game.GetRedisPlayerBagItemObj(playerId)
	//检测是否时永久道具，永久道具只能购买一个
	if !reqMsg.GetIsBuy() {
		if item.GetEffectiveTime() == for_game.COIN_PROPS_FOREVER {
			cfgItem := for_game.GetPropsItemInfo(item.GetPropsId())
			checkItems := []int64{}
			if item.GetPropsType() == for_game.COIN_PROPS_TYPE_LB {
				//礼包
				val := cfgItem.GetUseValue()
				err := json.Unmarshal([]byte(val), &checkItems)
				easygo.PanicError(err)
			} else {
				checkItems = append(checkItems, cfgItem.GetId())
			}
			names := bagObj.CheckHadForeverItem(checkItems)
			if len(names) > 0 {
				if len(checkItems) > 1 {
					reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_7)
					str := ""
					for i := 0; i < len(names); i++ {
						if i == len(names)-1 {
							str += names[i]
						} else {
							str += names[i] + "、"
						}
					}
					s := fmt.Sprintf("礼包中已包含永久商品%s，购买后相同属性物品会叠加", str)
					reqMsg.Reason = easygo.NewString(s)
				} else {
					reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_4)
					reqMsg.Reason = easygo.NewString("已经拥有该永久道具，无法重复购买!")
				}
				logs.Error("RpcBuyCoinItem 永久道具重复购买", reqMsg.GetId())
				return reqMsg
			}
		}
	}
	player := for_game.GetRedisPlayerBase(playerId)
	buyNum := int64(reqMsg.GetNum())
	if buyNum <= 0 {
		reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_6)
		logs.Error("RpcBuyCoinItem 购买数量为0!", reqMsg.GetId())
		return reqMsg
	}
	logs.Info("item:", item)
	propsItem := for_game.GetPropsItemInfo(item.GetPropsId())
	if propsItem == nil {
		reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_3)
		logs.Error("RpcBuyCoinItem 道具不存在!", reqMsg.GetId())
		return reqMsg
	}
	if reqMsg.GetIsCheck() {
		reqMsg.Result = easygo.NewInt32(COIN_RECHARGE_SUCCESS)
		return reqMsg
	}
	info := map[string]interface{}{
		"Name": item.GetName(),
		"Num":  easygo.AnytoA(buyNum),
	}
	//上边数量变成购买数量
	item.ProductNum = easygo.NewInt64(buyNum)
	reason := for_game.GetGoldChangeNote(for_game.COIN_TYPE_SHOP_OUT, playerId, info)
	var orderId string
	switch reqMsg.GetWay() {
	case for_game.COIN_PROPS_BUYWAY_COIN:
		//计算商品价值
		coin := item.GetCoin()
		playermgr := for_game.GetRedisPlayerBase(playerId)
		if playermgr != nil && playermgr.GetDeviceType() == for_game.TYPE_IOS && item.GetIosCoin() > 0 {
			coin = item.GetIosCoin()
		}

		t := for_game.GetMillSecond()
		if t >= item.GetStartTime() && t < item.GetEndTime() && item.GetCoinRebate() > 0 {
			coin = item.GetDisCoin()
		}
		coin = coin * buyNum

		if player.GetAllCoin() < coin {
			reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_2)
			return reqMsg
		}
		//扣掉对应金币
		orderId = for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, for_game.COIN_TYPE_SHOP_OUT)
		msg := &share_message.GoldExtendLog{
			RedPacketId: easygo.NewInt64(item.GetId()),
			OrderId:     easygo.NewString(orderId),
			Title:       easygo.NewString(reason),
		}
		if NotifyAddCoin(playerId, -coin, reason, for_game.COIN_TYPE_SHOP_OUT, msg) != "" {
			reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_5) //服务器数据异常
			logs.Error("消费减少硬币失败")
			return reqMsg
		}
	case for_game.COIN_PROPS_BUYWAY_MONEY:
		if reqMsg.GetWay() == for_game.COIN_PROPS_BUYWAY_MONEY && (reqMsg.GetPassWord() == "" || player.GetPayPassword() != reqMsg.GetPassWord()) {
			reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_8)
			logs.Error("RpcBuyCoinItem 支付密码错误!", reqMsg.GetId())
			return reqMsg
		}
		gold := item.GetPrice() * buyNum
		if player.GetGold() < gold {
			reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_2)
			return reqMsg
		}
		//扣掉对应零钱
		orderId = for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, for_game.GOLD_TYPE_XN_SHOP_MONEY)
		msg := &share_message.GoldExtendLog{
			RedPacketId: easygo.NewInt64(item.GetId()),
			OrderId:     easygo.NewString(orderId),
			Title:       easygo.NewString(reason),
			PayType:     easygo.NewInt32(for_game.PAY_TYPE_GOLD),
		}
		if NotifyAddGold(playerId, -gold, reason, for_game.GOLD_TYPE_XN_SHOP_MONEY, msg) != "" {
			reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_5) //服务器数据异常
			logs.Error("消费减少零钱失败")
			return reqMsg
		}
	}
	//增加物品到背包
	NotifyAddBagItems(playerId, []*share_message.CoinProduct{item}, for_game.COIN_ITEM_GETTYPE_BUY, reqMsg.GetWay(), orderId, "")
	reqMsg.Result = easygo.NewInt32(COIN_RECHARGE_SUCCESS)
	return reqMsg
}

//使用虚拟物品:目前只有装备功能
func (self *cls1) RpcUseCoinItem(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.UseCoinItem, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcUseCoinItem", reqMsg, who.GetPlayerId(), common)
	playerId := who.GetPlayerId()
	bagObj := for_game.GetRedisPlayerBagItemObj(playerId)
	item := bagObj.GetItemNetId(reqMsg.GetId())
	if item == nil {
		reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_2)
		logs.Error("道具不存在或者已过期id=", reqMsg.GetId())
		return reqMsg
	}
	if item.GetStatus() == for_game.COIN_BAG_ITEM_EXPIRED {
		reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_3)
		logs.Error("道具已过期id=", reqMsg.GetId())
		return reqMsg
	}
	propsItem := for_game.GetPropsItemInfo(item.GetPropsId())
	if propsItem == nil {
		reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_2)
		logs.Error("道具不存在id=", reqMsg.GetId())
		return reqMsg
	}
	if propsItem.GetUseType() == for_game.COIN_PROPS_USETYPE_EQUIPMENT {
		equipmentObj := for_game.GetRedisPlayerEquipmentObj(playerId)
		newItems := make([]*share_message.PlayerBagItem, 0)
		//卸下旧的装备
		oldId := equipmentObj.GetCurEquipment(propsItem.GetPropsType())
		isChange := false
		if reqMsg.GetWay() == for_game.EQUIPMENT_UP && item.GetStatus() == for_game.COIN_BAG_ITEM_UNUSE {
			//装新的装备
			if oldId == item.GetId() {
				reqMsg.Result = easygo.NewInt32(COIN_SHOP_FAIL_7)
				logs.Error("重复使用道具", reqMsg.GetId())
				return reqMsg
			}
			if oldId != 0 {
				bagObj.SetItemStatus(oldId, for_game.COIN_BAG_ITEM_UNUSE)
				oldItem := bagObj.GetItemNetId(oldId)
				newItems = append(newItems, oldItem)
			}
			equipmentObj.Equipment(propsItem.GetPropsType(), item.GetId())
			bagObj.SetItemStatus(item.GetId(), for_game.COIN_BAG_ITEM_USED)
			newItem := bagObj.GetItemNetId(item.GetId())
			newItems = append(newItems, newItem)
			isChange = true
		} else if reqMsg.GetWay() == for_game.EQUIPMENT_DOWN && item.GetStatus() == for_game.COIN_BAG_ITEM_USED {
			bagObj.SetItemStatus(item.GetId(), for_game.COIN_BAG_ITEM_UNUSE)
			newItem := bagObj.GetItemNetId(item.GetId())
			newItems = append(newItems, newItem)
			//卸下装备
			if oldId == reqMsg.GetId() {
				equipmentObj.Equipment(propsItem.GetPropsType(), 0)
				isChange = true
			}
		} else {
			logs.Error("前端传参异常")
		}
		if len(newItems) > 0 {
			//通知玩家道具修改
			ep.RpcModifyBagItem(&client_hall.BagItems{Items: newItems})
		}
		//通知前端装备修改
		if isChange {
			newEquipment := equipmentObj.GetEquipmentForClient()
			ep.RpcModifyEquipment(newEquipment)
		}
	} else {
		logs.Info("使用消耗型道具")
	}
	reqMsg.Result = easygo.NewInt32(COIN_RECHARGE_SUCCESS)
	return reqMsg
}

//获取玩家装备信息
func (self *cls1) RpcGetPlayerEquipment(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.EquipmentReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetPlayerEquipment", reqMsg, who.GetPlayerId(), common)
	if reqMsg.GetId() <= 0 {
		return easygo.EmptyMsg
	}
	equipmentObj := for_game.GetRedisPlayerEquipmentObj(reqMsg.GetId())
	equipment := equipmentObj.GetEquipmentForClient()
	return equipment
}

//获取玩家背包信息
func (self *cls1) RpcGetPlayerBagItems(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.BagItems, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetPlayerBagItems", reqMsg, who.GetPlayerId(), common)
	bagObj := for_game.GetRedisPlayerBagItemObj(who.GetPlayerId())
	items := bagObj.GetRedisPlayerBagItem(reqMsg.GetType())
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetCreateTime() > items[j].GetCreateTime() // 降序
	})
	bMsg := &client_hall.BagItems{
		Items: items,
	}
	// 修改背包isNew
	for _, v := range items {
		v1 := new(share_message.PlayerBagItem)
		for_game.StructToOtherStruct(v, v1)
		v1.IsNew = easygo.NewBool(false)
		bagObj.UpdateItem(v1)
	}
	return bMsg
}

//兑换硬币活动
func (self *cls1) RpcCoinRechargeAct(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.RechargeActReq, common ...*base.Common) easygo.IMessage {
	playerId := who.GetPlayerId()
	logs.Info("RpcCoinRechargeAct", reqMsg, playerId, common)
	item := for_game.GetCoinRechargeActCfg(reqMsg.GetActCfgId())

	resp := &client_hall.CoinRechargeResp{
		Id:     easygo.NewInt64(reqMsg.GetActCfgId()),
		Result: easygo.NewInt32(COIN_RECHARGE_SUCCESS),
	}
	//青少年保护模式
	if who.GetYoungPassWord() != "" {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_8)
		return resp
	}
	//道具不存在
	if item == nil {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_3)
		return resp
	}
	giveType := reqMsg.GetGiveType()
	if giveType != for_game.RECHARGE_GIVE_DIAMOND && giveType != for_game.RECHARGE_GIVE_ESCOIN {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_3)
		return resp
	}

	price := item.GetAmount()
	if reqMsg.GetPayWay() == for_game.PAY_TYPE_GOLD {
		//零钱不足
		if who.GetGold() < price {
			resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_2)
			return resp
		}
		if who.GetPayPassword() == "" {
			resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_7)
			logs.Error("未设置零钱支付密码")
			return resp
		}
		if who.GetPayPassword() != reqMsg.GetPassWord() {
			resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_6)
			logs.Error("支付密码错误")
			return resp
		}
	} else {
		payInfo := reqMsg.GetPayInfo()
		if payInfo == nil {
			resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_10)
			logs.Error("PayInfo为空")
			return resp
		}
		payInfo.ExtendValue = easygo.NewString("Activity")
		ep.RpcTunedUpPayInfo(reqMsg.GetPayInfo()) //通知客户端调起支付
		return resp
	}
	return CoinRechargeActService(reqMsg, who, item)
}

//充值活动服务
func CoinRechargeActService(reqMsg *client_hall.RechargeActReq, player *Player, item *share_message.WishCoinRechargeActivityCfg) easygo.IMessage {
	logs.Info("------- 充值活动服务 CoinRechargeActService -------")
	//扣除零钱，并增加硬币
	orderId := reqMsg.GetOrderId()
	var money int64

	if orderId == "" {
		orderId = for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, for_game.GOLD_TYPE_EXCHANGE_COIN)
	}
	resp := &client_hall.CoinRechargeResp{
		Id:     easygo.NewInt64(reqMsg.GetActCfgId()),
		Result: easygo.NewInt32(COIN_RECHARGE_SUCCESS),
	}

	if item == nil { //使用第三方支付
		item = for_game.GetCoinRechargeActCfg(reqMsg.GetActCfgId())
		if item == nil { //道具不存在
			logs.Error("用户(%v)使用第三方支付充值后兑换硬币活动使用活动id(%v)无效", player.GetId(), reqMsg.GetActCfgId())
			resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_3)
			return resp
		}
		money = reqMsg.GetOrderAmount()
	} else { //使用零钱支付
		money = item.GetAmount()
	}
	playerId := player.GetPlayerId()
	actCfgId := item.GetId()
	giveType := reqMsg.GetGiveType()

	//获取赠送币数量以及类型
	isFirst := true
	var giveCoin int64
	playerRecharge := for_game.GetPlayerRechargeActFirst(playerId)
	for _, v := range playerRecharge.GetLevels() {
		if v == reqMsg.GetActCfgId() {
			isFirst = false
			break
		}
	}
	if isFirst {
		if giveType == for_game.RECHARGE_GIVE_DIAMOND {
			giveCoin = item.GetFirstDiamond()
		} else {
			giveCoin = item.GetFirstEsCoin()
		}
	} else {
		if giveType == for_game.RECHARGE_GIVE_DIAMOND {
			giveCoin = item.GetDailyDiamond()
		} else {
			giveCoin = item.GetDailyEsCoin()
		}
	}

	reason := "充值活动消费"
	msg := &share_message.GoldExtendLog{
		OrderId:  easygo.NewString(orderId),
		PayType:  easygo.NewInt32(for_game.PAY_TYPE_GOLD),
		Title:    easygo.NewString(reason),
		Gold:     easygo.NewInt64(actCfgId), //增加硬币数量
		BankName: easygo.NewString(reqMsg.GetBankCard()),
	}
	if NotifyAddGold(playerId, -money, reason, for_game.GOLD_TYPE_EXCHANGE_COIN, msg) != "" {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_4) //服务器数据异常
		logs.Error("充值扣除零钱失败")
		return resp
	}

	reason = "充值活动充值获得"
	msg = &share_message.GoldExtendLog{
		RedPacketId: easygo.NewInt64(actCfgId),
		Title:       easygo.NewString(reason),
		OrderId:     easygo.NewString(orderId),
	}
	if NotifyAddCoin(playerId, actCfgId, reason, for_game.COIN_TYPE_EXCHANGE_IN, msg) != "" {
		resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_5) //服务器数据异常
		logs.Error("用户(%v)兑换增加硬币失败", playerId)
		return resp
	}

	//赠送币种
	if giveType == for_game.RECHARGE_GIVE_DIAMOND {
		wishSrv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)
		msg := &h5_wish.BackStageAddDiamondReq{
			PlayerId: easygo.NewInt64(playerId),
			Account:  easygo.NewString(player.GetPhone()),
			Channel:  easygo.NewInt32(1001),
			NickName: easygo.NewString(player.GetNickName()),
			HeadUrl:  easygo.NewString(player.GetHeadIcon()),
			Token:    easygo.NewString(""), //暂不需要
			Diamond:  easygo.NewInt64(giveCoin),
		}
		_, err := SendMsgToServerNew(wishSrv.GetSid(), "RpcBackstageAddDiamond", msg)
		if err != nil {
			logs.Error("用户(%v)充值活动赠送获得许愿池钻石失败: %v", playerId, err)
		}
	} else {
		if err := NotifyAddESportCoin(playerId, giveCoin, reason, for_game.ESPORTCOIN_TYPE_EXCHANGE_GIVE_IN, msg); err != "" {
			resp.Result = easygo.NewInt32(COIN_SHOP_FAIL_4) //服务器数据异常
			logs.Error("用户(%v)充值活动赠送获得电竞币失败: %v", playerId, err)
		}
	}
	if isFirst {
		for_game.AddPlayerRechargeActFirst(playerId, reqMsg.GetActCfgId())
	}
	for_game.AddPlayerRechargeActLog(playerId, money, actCfgId, giveCoin, giveType)
	return resp
}
