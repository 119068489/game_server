package hall

import (
	"encoding/base64"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
	"game_server/pb/h5_wish"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/akqp2019/protobuf/proto"
	"github.com/astaxie/beego/logs"
)

const (
	ADD_TEAM_TYPE_QRCODE    int32 = iota + 1 // 1 二维码
	ADD_TEAM_TYPE_CARD                       // 2 群名片
	ADD_TEAM_TYPE_PASSWORD                   // 3 群口令
	ADD_TEAM_TYPE_BS                         // 4 后台添加
	ADD_TEAM_TYPE_ADV                        // 5广告
	ADD_TEAM_TYPE_BACKSTAGE                  // 6 后台推送
	ADD_TEAM_TYPE_TOPIC                      // 7 通过话题进群
)

const (
	VC_MATCH_LEVEL_TOP    = 0 //特级匹配
	VC_MATCH_LEVEL_FIRST  = 1 //一级匹配
	VC_MATCH_LEVEL_SECOND = 2 //二级匹配
	VC_MATCH_LEVEL_THIRD  = 3 //三级匹配
	VC_MATCH_LEVEL_FOURTH = 4 //四级匹配
)

type PLAYER_ID = int64
type ENDPOINT_ID = easygo.ENDPOINT_ID
type DB_NAME = string
type SERVER_ID = int32
type INSTANCE_ID = int32
type GAME_TYPE = int32
type LEVEL = int32

type Int32List []int32

func (s Int32List) Len() int           { return len(s) }
func (s Int32List) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Int32List) Less(i, j int) bool { return s[i] < s[j] }

var _ClientEndpointId ENDPOINT_ID = 0
var GetPlayerPeriod = for_game.GetPlayerPeriod

const PENGJU_PAY_MIN_MONEY = 1000    //最小10元
const PENGJU_PAY_MAX_MONEY = 5000000 //最大50000元

// 生成游戏客户端的 endpoint id (U3D 客户端 + H5 客户端)
func GenClientEndpointId() ENDPOINT_ID {
	v := atomic.AddInt32(&_ClientEndpointId, 1) // 溢出后自动回转
	return v
}

/**
 支付预警功能服务器需要根据设定的预警规则，给预警人员发送短信
pid: 用户id
changeType: 交易类型. 1-充值;2-提现
amount: 金额
*/
func CheckWarningSMS(changeType int32, amount int64) {
	// 从配置表中读取配置.
	msg := PSysParameterMgr.GetSysParameter(for_game.WARNING_PARAMETER)
	if msg == nil {
		logs.Error("读取配置表中的警告参数失败,请求id为: ", for_game.WARNING_PARAMETER)
		return
	}
	if len(msg.GetPhoneList()) == 0 {
		logs.Error("读取配置表中的警告参数,需要发短信的手机号为空,请求id为: ", for_game.WARNING_PARAMETER)
		return
	}
	// 兼容手机号
	phones := make([]string, 0)
	for _, p := range msg.GetPhoneList() {
		if !strings.HasPrefix(p, "+86") {
			p = fmt.Sprintf("%s%s", "+86", p)
		}
		phones = append(phones, p)
	}
	msg.PhoneList = phones
	switch changeType {
	case for_game.GOLD_CHANGE_TYPE_IN: // 充值
		rechargeTime := msg.GetRechargeTime()         // 充值预警频次时间 取出来的单位是秒
		rechargeTimes := msg.GetRechargeTimes()       // 充值预警频次次数
		rechargeGoldRate := msg.GetRechargeGoldRate() //  充值预警金额时间 取出来的单位是秒
		rechargeGold := msg.GetRechargeGold()         // 充值预警金额
		// 设置充值预警次数
		count := for_game.SetRechargeCountToRedis(1, rechargeTime)
		if count >= rechargeTimes {
			logs.Error("用户充值次数超出预警,%d 秒限制 %d 次,现在是 %d 次了", rechargeTime, rechargeTimes, count)
			//  发送腾讯云警告短信
			for_game.NewSMSInst(for_game.SMS_BUSINESS_TC).SendWarningSMS(msg.GetPhoneList(), []string{easygo.Stamp2Str(time.Now().Unix())})
			return
		}
		// 设置充值金额
		gold := for_game.SetRechargeAmountToRedis(amount, rechargeGoldRate) // amount 单位是秒
		if gold >= rechargeGold {
			logs.Error("用户 充值金额超出预警,%d 秒限制 %v 元,现在是 %v 元了", rechargeGoldRate, rechargeGold/100, gold/100)
			//  发送腾讯云警告短信
			for_game.NewSMSInst(for_game.SMS_BUSINESS_TC).SendWarningSMS(msg.GetPhoneList(), []string{easygo.Stamp2Str(time.Now().Unix())})
			return
		}
		// 设置充值预警金额
	case for_game.GOLD_CHANGE_TYPE_OUT: // 提现
		withdrawalTime := msg.GetWithdrawalTime()         // 提现预警频次时间 取出来的单位是秒
		withdrawalTimes := msg.GetWithdrawalTimes()       // 提现预警频次次数
		withdrawalGoldRate := msg.GetWithdrawalGoldRate() // 提现预警金额时间 取出来的单位是秒
		withdrawalGold := msg.GetWithdrawalGold()         // 提现预警总金额
		// 设置提现预警次数
		count := for_game.SetWithdrawCountToRedis(1, withdrawalTime)
		if count >= withdrawalTimes {
			logs.Error("用户 提现次数超出预警,%d 秒限制 %d 次,现在是 %d 次了", withdrawalTime, withdrawalTimes, count)
			//  发送腾讯云警告短信
			for_game.NewSMSInst(for_game.SMS_BUSINESS_TC).SendWarningSMS(msg.GetPhoneList(), []string{easygo.Stamp2Str(time.Now().Unix())})
			return
		}
		gold := for_game.SetWithdrawAmountToRedis(amount, withdrawalGoldRate)
		if gold >= withdrawalGold {
			logs.Error("用户 提现金额超出预警,%d 秒限制 %v 元,现在是 %v 元了", withdrawalGoldRate, withdrawalGold/100, gold/100)
			//  发送腾讯云警告短信
			for_game.NewSMSInst(for_game.SMS_BUSINESS_TC).SendWarningSMS(msg.GetPhoneList(), []string{easygo.Stamp2Str(time.Now().Unix())})
			return

		}
	default:
		logs.Error("操作类型有误,changeType: ", changeType)
		return
	}
}

// 充值成功后的处理
func HandleAfterRecharge(orderId string) string {
	//更新充值信息
	logs.Info("处理订单：", orderId)
	order := for_game.GetRedisOrderObj(orderId)
	logs.Info("订单:%+v", order)
	if order == nil {
		logs.Info("无效的支付订单")
		return "无效的支付订单:" + orderId
	}

	//如果是充值，并且订单取消，直接返回
	if order.GetChangeType() == for_game.GOLD_CHANGE_TYPE_IN && order.GetStatus() == for_game.ORDER_ST_CANCEL {
		logs.Info("充值订单已取消")
		return "充值订单已取消"
	}
	playerId := order.GetPlayerId()
	base := GetPlayerObj(playerId)
	if base == nil {
		s := fmt.Sprintf("玩家Id: %d 不存在", playerId)
		return s
	}

	//许愿池订单状态修改
	if order.GetOrderType() == for_game.ORDER_ST_WISH {
		srv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)
		if srv == nil {
			logs.Error("许愿池回收回调::找不到许愿池服务器？商城许愿池断线？")
			return "许愿池回收回调::找不到许愿池服务器？许愿池网络断线？"
		} else {
			SendMsgToServerNew(srv.GetSid(), "RpcRecycleCallBack", &h5_wish.OrderMsgResp{OrderId: easygo.NewString(order.GetId())})
		}
		//订单处理完成
		order.SetOverTime(for_game.GetMillSecond())
		order.SetStatus(for_game.ORDER_ST_FINISH)
		//redis数据存储到mongo
		order.SaveToMongo()
		return ""
	}

	var reason string
	switch order.GetSourceType() {
	case for_game.GOLD_TYPE_CASH_AFIN:
		reason = "人工入款"
	case for_game.GOLD_TYPE_CASH_IN:
		reason = "充值成功"
	case for_game.GOLD_TYPE_GET_REDPACKET:
		reason = "收红包"
	case for_game.GOLD_TYPE_GET_TRANSFER_MONEY:
		reason = "转入成功"
	case for_game.GOLD_TYPE_GET_MONEY:
		reason = "收款成功"
	case for_game.GOLD_TYPE_REDPACKET_OVERTIME:
		reason = "红包退款"
	case for_game.GOLD_TYPE_TRANSFER_MONEY_OVER:
		reason = "转账退款"
	case for_game.GOLD_TYPE_BACK_MONEY:
		reason = "商家退款"
	case for_game.GOLD_TYPE_CASH_AFOUT:
		reason = "人工出款"
	case for_game.GOLD_TYPE_CASH_OUT:
		reason = "提现成功"
	case for_game.GOLD_TYPE_SEND_REDPACKET:
		reason = "发红包"
	case for_game.GOLD_TYPE_SEND_TRANSFER_MONEY:
		reason = "转出成功"
	case for_game.GOLD_TYPE_PAY_MONEY:
		reason = "付款成功"
	case for_game.GOLD_TYPE_FINE_MONEY:
		reason = "罚没"
	case for_game.GOLD_TYPE_EXTRA_MONEY:
		reason = "手续费"
	case for_game.GOLD_TYPE_CASH_OUT_BACK:
		reason = "取消提款"

	}
	//订单LOG扩展数据
	var bankInfo, bankName, bankId string
	payType := order.GetPayType()
	if payType == for_game.PAY_TYPE_BANKCARD { //银行卡支付
		Id := order.GetBankInfo()
		if Id == "" {
			Id = order.GetAccountNo()
		}
		bankId := Id[len(Id)-4:]
		code := for_game.GetBankCodeForBankId(Id)
		bankName = for_game.BankName[code]
		bankInfo = bankName + bankId
	}
	var extendLog interface{}
	if order.GetSourceType() == for_game.GOLD_TYPE_CASH_IN { //充值订单
		info := map[string]interface{}{
			"Type": payType,
		}

		if payType == for_game.PAY_TYPE_BANKCARD { //银行卡支付
			info["BankId"] = bankId
			info["BankName"] = bankName
		}
		reason1 := for_game.GetGoldChangeNote(order.GetSourceType(), 0, info)
		extendLog = &share_message.GoldExtendLog{
			OrderId:  easygo.NewString(orderId),
			Title:    easygo.NewString(reason1),
			PayType:  easygo.NewInt32(payType),
			Gold:     easygo.NewInt64(base.GetGold() + order.GetChangeGold()),
			BankName: easygo.NewString(bankInfo),
		}
	} else if order.GetSourceType() == for_game.GOLD_TYPE_CASH_AFIN {
		var tt int32
		if order.GetChangeType() == for_game.GOLD_CHANGE_TYPE_IN {
			tt = for_game.PAY_TYPE_BACKSTAGE_IN
		} else if order.GetChangeType() == for_game.GOLD_CHANGE_TYPE_OUT {
			tt = for_game.PAY_TYPE_BACKSTAGE_OUT
		}
		info := map[string]interface{}{
			"Type": tt,
		}
		reason1 := for_game.GetGoldChangeNote(order.GetSourceType(), 0, info)
		extendLog = &share_message.GoldExtendLog{
			OrderId: easygo.NewString(orderId),
			Title:   easygo.NewString(reason1),
			PayType: easygo.NewInt32(payType),
			Gold:    easygo.NewInt64(base.GetGold() + order.GetChangeGold()),
		}
	} else {
		extendLog = GetExtendLog(order.GetRedisOrder())
	}

	logs.Info("type:", order.GetChangeType())
	var result string
	switch order.GetChangeType() {
	case for_game.GOLD_CHANGE_TYPE_IN:
		changeGold := order.GetChangeGold() - (-(order.GetTax())) //不含税的订单金额
		//充值零钱到账
		result = NotifyAddGold(base.GetPlayerId(), changeGold, reason, order.GetSourceType(), extendLog)
		logs.Info("result:", result)
		if result != "" {
			logs.Info("type:", len(result))
			return result
		}
		//如果有税收进行扣税操作
		if order.GetTax() < 0 {
			extendLog := &share_message.RechargeExtend{
				PayChannel:  easygo.NewInt32(order.GetPayChannel()),
				PayType:     easygo.NewInt32(order.GetPayType()),
				Channeltype: easygo.NewInt32(order.GetChanneltype()),
				CreateIP:    easygo.NewString(order.GetCreateIP()),
				OrderId:     easygo.NewString(order.GetOrderId()),
				Amount:      easygo.NewInt64(order.GetTax()),
				Operator:    easygo.NewString(order.GetOperator()),
			}
			NotifyAddGold(base.GetPlayerId(), order.GetTax(), "手续费", for_game.GOLD_TYPE_EXTRA_MONEY, extendLog)
		}
	case for_game.GOLD_CHANGE_TYPE_OUT:
		changeGold := order.GetChangeGold()
		if order.GetStatus() == 3 { //取消提现退款
			changeGold = -(order.GetChangeGold()) //含税的订单金额

			reason1 := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_CASH_OUT_BACK, 0, nil)
			extendLog = &share_message.GoldExtendLog{
				OrderId: easygo.NewString(orderId),
				Title:   easygo.NewString(reason1),
				Gold:    easygo.NewInt64(base.GetGold() + order.GetChangeGold()),
			}

			reason = reason1
			//充值零钱到账
			result = NotifyAddGold(base.GetPlayerId(), changeGold, reason, for_game.GOLD_TYPE_CASH_OUT_BACK, extendLog)
			return result

		} else {
			//人工出款扣钱,//下订单的时候已经扣款了
			if order.GetSourceType() != for_game.GOLD_TYPE_CASH_OUT {
				result = NotifyAddGold(base.GetPlayerId(), changeGold, reason, order.GetSourceType(), extendLog)
				if result != "" {
					logs.Info("type:", len(result))
					return result
				}
			} else {
				//记录玩家提现次数和额度
				period := for_game.GetPlayerPeriod(order.GetPlayerId())
				period.DayPeriod.AddInteger("OutTimes", 1)
				period.DayPeriod.AddInteger("OutSum", order.GetAmount())
			}
		}

	}
	//如果充值，通知前端充值完成
	ep := ClientEpMgr.LoadEndpointByPid(base.GetPlayerId())

	//如果是充值发红包

	payWay := order.GetPayWay()
	//logs.Info("pay:", payWay)
	if payWay == for_game.PAY_TYPE_REPACKET_PERSIONAL || order.GetPayWay() == for_game.PAY_TYPE_REPACKET_TEAM {
		//个人发红包
		player := GetPlayerObj(playerId)
		cl := &cls1{}
		//设置发红包参数
		tp := 1
		if order.GetPayWay() == for_game.PAY_TYPE_REPACKET_TEAM {
			tp = 2
		}
		msg := &share_message.RedPacket{
			Type:       easygo.NewInt32(tp),
			Sender:     easygo.NewInt64(playerId),
			TargetId:   easygo.NewInt64(order.GetPayTargetId()),
			TotalMoney: easygo.NewInt64(order.GetChangeGold()),
			TotalCount: easygo.NewInt32(order.GetTotalCount()),
			PerMoney:   easygo.NewInt64(0),
			Content:    easygo.NewString(order.GetContent()),
			PayWay:     easygo.NewInt32(payType),
		}
		if payType == for_game.PAY_TYPE_BANKCARD {
			msg.BankCard = easygo.NewString(bankInfo)
		}
		//res := cl.RpcSendRedPacket(ep, player, msg)
		//if res != nil {
		//	logs.Info("res:", res)
		//	panic("充值发红包异常")
		//}
		//logs.Info("发红包结束")
		chatLog := CreateRedPacket(msg, orderId)
		cl.RpcChatNew(nil, player, chatLog)
		logs.Info("发红包结束")
	} else if order.GetPayWay() == for_game.PAY_TYPE_TRANSFER {
		//个人转账
		player := GetPlayerObj(playerId)
		cl := &cls1{}
		//设置转账参数
		msg := &share_message.TransferMoney{
			Sender:      easygo.NewInt64(playerId),
			TargetId:    easygo.NewInt64(order.GetPayTargetId()),
			Way:         easygo.NewInt32(payType),
			Gold:        easygo.NewInt64(order.GetChangeGold()),
			Content:     easygo.NewString(order.GetContent()),
			PayPassWord: easygo.NewString(player.GetPassword()),
		}
		if payType == for_game.PAY_TYPE_BANKCARD {
			msg.BankInfo = easygo.NewString(bankInfo)
		}
		chatLog := CreateTransferMoney(msg, orderId)
		cl.RpcChatNew(nil, player, chatLog)
	} else if order.GetPayWay() == for_game.PAY_TYPE_SHOP {
		//TODO 购买充值:
		srv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_SHOP)
		if srv == nil {
			//恢复库存
			rst1 := for_game.ShopRecoverStock(order.GetPayTargetId())
			if rst1 != "" {
				logs.Error(rst1)
				return rst1
			}
			logs.Error("商城购买支付回调::找不到商城服务器？商城网络断线？")
			return "商城购买支付回调::找不到商城服务器？商城网络断线？"
		} else {
			player := GetPlayerObj(playerId)
			result := ShopOrderPay(order.GetPayTargetId(), player, payType, bankInfo, orderId)
			if result != "" {
				logs.Error(result)
				//恢复库存
				rst1 := for_game.ShopRecoverStock(order.GetPayTargetId())
				if rst1 != "" {
					logs.Error(rst1)
					return rst1
				}
			} else {
				SendMsgToServerNew(srv.GetSid(), "RpcPayCallBack", &share_message.OrderID{OrderId: easygo.NewInt64(order.GetPayTargetId())})
			}
		}
	} else if order.GetPayWay() == for_game.PAY_TYPE_CODE { //扫码付款充值
		player := GetPlayerObj(playerId)
		cl := &cls1{}
		msg := &client_hall.PayInfo{
			Gold:     easygo.NewInt64(order.GetChangeGold()),
			Content:  easygo.NewString(order.GetContent()),
			Type:     easygo.NewInt32(payType),
			PlayerId: easygo.NewInt64(order.GetPayTargetId()),
			IsWay:    easygo.NewBool(true),
			OrderId:  easygo.NewString(orderId),
		}
		if payType == for_game.PAY_TYPE_BANKCARD {
			msg.BankInfo = easygo.NewString(bankInfo)
		}
		cl.RpcPayForCode(ep, player, msg)
	} else if order.GetPayWay() == for_game.PAY_TYPE_TEAMCODE { //扫群二维码付款充值
		player := GetPlayerObj(playerId)
		cl := &cls1{}
		msg := &client_hall.PayInfo{
			Gold:     easygo.NewInt64(order.GetChangeGold()),
			Content:  easygo.NewString(order.GetContent()),
			Type:     easygo.NewInt32(payType),
			PlayerId: easygo.NewInt64(order.GetPayTargetId()),
			IsWay:    easygo.NewBool(true),
			OrderId:  easygo.NewString(orderId),
		}
		if payType == for_game.PAY_TYPE_BANKCARD {
			msg.BankInfo = easygo.NewString(bankInfo)
		}
		cl.RpcPayForCode(ep, player, msg)
	} else if order.GetPayWay() == for_game.PAY_TYPE_COIN {
		//充值兑换硬币
		player := GetPlayerObj(playerId)
		cl := &cls1{}
		if order.GetExtendValue() == "Activity" {
			msg := &client_hall.RechargeActReq{
				ActCfgId:    easygo.NewInt64(order.GetPayTargetId()),
				PassWord:    easygo.NewString(player.GetPayPassword()),
				GiveType:    easygo.NewInt32(order.GetTotalCount()),
				OrderId:     easygo.NewString(orderId),
				OrderAmount: easygo.NewInt64(order.GetAmount()),
			}
			if payType == for_game.PAY_TYPE_BANKCARD {
				msg.BankCard = easygo.NewString(bankInfo)
			}
			CoinRechargeActService(msg, player, nil)
			//cl.RpcCoinRechargeAct(ep, player, msg)
		} else {
			msg := &client_hall.CoinRechargeReq{
				Id:       easygo.NewInt64(order.GetPayTargetId()),
				PassWord: easygo.NewString(player.GetPayPassword()),
			}
			cl.RpcCoinRecharge(ep, player, msg)
		}
	} else if order.GetPayWay() == for_game.PAY_TYPE_COIN_ITEM {
		//充值购买硬币商品
		player := GetPlayerObj(playerId)
		cl := &cls1{}
		msg := &client_hall.BuyCoinItem{
			Id:       easygo.NewInt64(order.GetPayTargetId()),
			Num:      easygo.NewInt32(order.GetTotalCount()),
			Way:      easygo.NewInt32(for_game.COIN_PROPS_BUYWAY_MONEY),
			IsBuy:    easygo.NewBool(true),
			PassWord: easygo.NewString(player.GetPayPassword()),
		}
		cl.RpcBuyCoinItem(ep, player, msg)
	}
	//订单处理完成
	order.SetOverTime(for_game.GetMillSecond())
	order.SetStatus(for_game.ORDER_ST_FINISH)
	// if order.GetSourceType() == for_game.GOLD_TYPE_CASH_IN || order.GetSourceType() == for_game.GOLD_TYPE_CASH_OUT {
	// 	easygo.Spawn(func() {
	// 		for_game.MakeOperationChannelReport(5, order.GetPlayerId(), "", nil, order.GetRedisOrder())
	// 	}) //生成运营渠道数据汇总报表 已优化到Redis
	// }
	//redis数据存储到mongo
	order.SaveToMongo()
	return result
}

func ShopOrderPay(order_id int64, player *Player, payType int32, bankInfo string, orId ...string) string {
	var pay_gold int32 = 0
	var bill *share_message.TableBill = &share_message.TableBill{}
	var order_list []int64 = []int64{}

	colBill, closeFunBill := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_BILLS)
	defer closeFunBill()

	colOrder, closeFunOrder := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFunOrder()

	eBill := colBill.Find(bson.M{"_id": order_id}).One(bill)
	if eBill != nil && eBill != mgo.ErrNotFound {
		logs.Error(eBill)
		return "操作失败"
	}

	if eBill == mgo.ErrNotFound {
		order := share_message.TableShopOrder{}
		eOrder := colOrder.Find(bson.M{"_id": order_id}).One(&order)

		if eOrder != nil && eOrder != mgo.ErrNotFound {
			logs.Error(eOrder)
			return "操作失败"
		}

		if eOrder == mgo.ErrNotFound {
			logs.Error(eOrder)
			return "订单不存在"
		}

		pay_gold = order.GetItems().GetPrice() * order.GetItems().GetCount()
		order_list = []int64{order.GetOrderId()}
	} else {
		pay_gold = bill.GetPrice()
		order_list = bill.OrderList
	}
	base := for_game.GetPlayerById(player.GetPlayerId())
	if base == nil {
		return "找不到该用户"
	}

	if player.GetGold() < int64(pay_gold) {
		return "余额不足"
	}

	for _, id := range order_list {
		order := share_message.TableShopOrder{}
		eOrder := colOrder.Find(bson.M{"_id": id}).One(&order)

		if eOrder != nil && eOrder != mgo.ErrNotFound {
			logs.Error(eOrder)
			return "操作失败"
		}
		if eOrder == mgo.ErrNotFound {
			logs.Error("找不到订单%v", id)
			return "订单不存在"
		}

		sponsor := GetPlayerObj(order.GetSponsorId())
		if sponsor == nil {
			logs.Error("玩家对象为空")
			return "操作失败"
		}

		sponsor_name := sponsor.GetNickName()
		gold := int64(order.GetItems().GetPrice() * order.GetItems().GetCount())
		//orderId, errMsg := for_game.PlaceOrder(player.GetPlayerId(), -gold, for_game.GOLD_TYPE_SHOP_MONEY)
		orderId := append(orId, "")[0]
		if orderId == "" {
			orderId = for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, for_game.GOLD_TYPE_SHOP_MONEY)
		}
		//if errMsg.GetCode() != for_game.FAIL_MSG_CODE_SUCCESS {
		//		//	return errMsg.GetReason()
		//		//}
		info := map[string]interface{}{
			"SponsorName": sponsor_name,
		}
		MerchantId := easygo.IntToString(int(order_id))
		reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_SHOP_MONEY, player.GetPlayerId(), info)
		msg := &share_message.GoldExtendLog{
			OrderId:    easygo.NewString(orderId),
			MerchantId: &MerchantId,
			PayType:    easygo.NewInt32(payType),
			Title:      &reason,
			Gold:       easygo.NewInt64(base.GetGold() - gold),
			Account:    easygo.NewString(player.GetAccount()),
		}
		if bankInfo != "" {
			msg.BankName = easygo.NewString(bankInfo)
		}
		if payType == for_game.PAY_TYPE_GOLD {
			msg.Gold = easygo.NewInt64(player.GetGold() - gold)
		}
		return NotifyAddGold(player.GetPlayerId(), -gold, reason, for_game.GOLD_TYPE_SHOP_MONEY, msg) //买家付款
		//errMsgFinish := FinishOrder(orderId, msg)
		//if errMsgFinish.GetCode() != for_game.FAIL_MSG_CODE_SUCCESS {
		//	return errMsgFinish.GetReason()
		//}
	}
	return ""
}

//通知其他服务器，玩家上线
func NotifyPlayerOnLine(playerId PLAYER_ID) {
	msg := &share_message.PlayerState{
		PlayerId: easygo.NewInt64(playerId),
		ServerId: easygo.NewInt32(PServerInfo.GetSid()),
	}
	logs.Info("通知其他服务器上线，", msg)
	//通知其他大厅
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPlayerOnLine", msg)
	// 通知广场.
	req := &server_server.ReloadDynamicReq{
		PlayerId: easygo.NewInt64(playerId),
	}
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_SQUARE, "RpcHallNotifyReloadSquare", req)
	//用户上线通知电竞接口服务器
	SendMsgToESportsApply("RpcESPortsPlayerOnLine", msg)

}

//通知其他服务器，玩家离线
func NotifyPlayerOffLine(playerId PLAYER_ID) {
	msg := &share_message.PlayerState{
		PlayerId: easygo.NewInt64(playerId),
	}
	//通知其他大厅
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPlayerOffLine", msg)
	//通知接口电竞服务器
	SendMsgToESportsApply("RpcESPortsPlayerOffLine", msg)
}

//广播消息给所有在线客户端
func BroadCastMsgToOnlineClient(methodName string, msg easygo.IMessage) {
	endPoints := ClientEpMgr.GetEndpoints()
	for _, ep := range endPoints {
		if ep != nil {
			epx := ep.(IGameClientEndpoint)
			epx.CallRpcMethod(methodName, msg)
		}
	}
}

//负载均衡一台后台处理事务
func ChooseOneBackStage(methodName string, msg easygo.IMessage) easygo.IMessage {
	srv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_BACKSTAGE)
	if srv == nil {
		return easygo.NewFailMsg("找不到服务器")
	}
	resp, err := SendMsgToServerNew(srv.GetSid(), methodName, msg)
	if err != nil {
		return err
	}
	return resp
}

//随机给指定类型服务器发送
//func SendToIdelOtherServer(t int32, methodName string, msg easygo.IMessage) (easygo.IMessage, *base.Fail) {
//	srv := PServerInfoMgr.GetIdelServer(t)
//	if srv != nil {
//		return PWebApiForServer.SendToServer(srv, methodName, msg)
//	}
//	return nil, easygo.NewFailMsg("找不到服务器")
//}

//玩家充值发货
func RechargeGoldToPlayer(playerId PLAYER_ID, orderId string) {
	//先检测玩家是否在本服
	epPlayer := ClientEpMp.LoadEndpoint(playerId)
	if epPlayer != nil {
		HandleAfterRecharge(orderId)
	} else {
		//玩家不在本服，先检测玩家是否在线
		if PlayerOnlineMgr.CheckPlayerIsOnLine(playerId) {
			serverId := PlayerOnlineMgr.GetPlayerServerId(playerId)
			reqMsg := &server_server.Recharge{
				PlayerId:     easygo.NewInt64(playerId),
				OrderId:      easygo.NewString(orderId),
				RechargeGold: easygo.NewInt64(0), //字段没用到，但proto定义为require，随便赋值
			}
			SendMsgToServerNew(serverId, "RpcRechargeToHall", reqMsg)
		} else { //不在线直接发货
			HandleAfterRecharge(orderId)
		}
	}
}

func SendLoginMessage(pid PLAYER_ID, name string) {
	title := "欢迎来到柠檬畅聊"
	text := `
%s，您可以在这里获得以下服务：
1、结识附近的朋友，点击“广场”——点击“附近的人”，即可与您附近的朋友畅聊。
2、点“好物”有过万件精选商品，供您任意挑选心仪好物。
`
	//将您的闲置物品分享给其他的六十万用户，点击”好物“——点击右下角“+号”可一键发布您的闲置物品。
	content := fmt.Sprintf(text, name)
	NoticeAssistant(pid, 1, title, content)
}

func NoticeAssistant(pid PLAYER_ID, t int32, title, content string) {
	now := util.GetMilliTime()
	assistant_id := for_game.NextId("Assistant")
	msg := &share_message.Assistant{
		CreateTime: easygo.NewInt64(now),
		Type:       &t,
		Id:         easygo.NewInt64(assistant_id),
		Title:      easygo.NewString(title),
		Content:    easygo.NewString(content),
		PlayerId:   easygo.NewInt64(pid),
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ASSISTANT)
	defer closeFun()
	e := col.Insert(msg)
	if e != nil {
		logs.Error(e)
	}
	notice := &client_server.AssistantMsg{}
	if t == 1 {
		notice = &client_server.AssistantMsg{
			LogId:      easygo.NewInt64(assistant_id),
			MsgType:    &t,
			DateTime:   &now,
			Title:      &title,
			SysContent: &content,
			PlayerId:   &pid,
		}
	} else {
		notice = &client_server.AssistantMsg{
			LogId:            easygo.NewInt64(assistant_id),
			MsgType:          &t,
			DateTime:         &now,
			Title:            &title,
			SysNoticeContent: &content,
			PlayerId:         &pid,
		}
	}
	if PlayerOnlineMgr.CheckPlayerIsOnLine(pid) { //如果在线
		ep2 := ClientEpMp.LoadEndpoint(pid)
		if ep2 != nil {
			ep2.RpcAssistantNotify(notice)
		} else {
			serverId := PlayerOnlineMgr.GetPlayerServerId(pid)
			SendMsgToServerNew(serverId, "RpcAssistantNotify", notice)
		}
	} else {
		logs.Error("不在线,pid------------------------------------->", pid)
	}
}

func NoticeAssistant2(player *Player, friendId int64, t int32, friendType int32, assistantType int32) {
	if player == nil {
		return
	}
	now := util.GetMilliTime()
	assistant_id := for_game.NextId("Assistant")
	msg := &share_message.Assistant{
		PlayerId:      easygo.NewInt64(friendId),
		CreateTime:    easygo.NewInt64(now),
		Phone:         easygo.NewString(player.GetPhone()),
		Type:          &assistantType,
		FriendType:    &friendType,
		Id:            easygo.NewInt64(assistant_id),
		Nickname:      easygo.NewString(player.GetNickName()),
		Avatar:        easygo.NewString(player.GetHeadIcon()),
		Account:       easygo.NewString(player.GetAccount()),
		Photo:         player.GetPhoto(),
		AddFriendType: easygo.NewInt32(int32(t)),
		Signature:     easygo.NewString(player.GetSignature()),
		Pid:           easygo.NewInt64(player.GetPlayerId()),
		Sex:           easygo.NewInt32(player.GetSex()),
	}

	notice := &client_server.AssistantMsg{
		LogId:         easygo.NewInt64(assistant_id),
		MsgType:       easygo.NewInt32(assistantType),
		DateTime:      easygo.NewInt64(now),
		AddPalType:    easygo.NewInt32(friendType),
		NickName:      easygo.NewString(player.GetNickName()),
		Account:       easygo.NewString(player.GetAccount()),
		HeadIcon:      easygo.NewString(player.GetHeadIcon()),
		Phone:         easygo.NewString(player.GetPhone()),
		AddFriendType: easygo.NewInt32(t),
		Photo:         player.GetPhoto(),
		Signature:     easygo.NewString(player.GetSignature()),
		PlayerId:      easygo.NewInt64(player.GetPlayerId()),
		Sex:           easygo.NewInt32(player.GetSex()),
		Types:         easygo.NewInt32(player.GetTypes()),
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ASSISTANT)
	defer closeFun()
	e := col.Insert(msg)
	if e != nil {
		logs.Error(e)
	}

	if !PlayerOnlineMgr.CheckPlayerIsOnLine(friendId) { //如果不在线
		return
	}
	// 如果在线
	ep2 := ClientEpMp.LoadEndpoint(friendId)
	if ep2 != nil {
		ep2.RpcAssistantNotify(notice)
		return
	}
	serverId := PlayerOnlineMgr.GetPlayerServerId(friendId)
	SendMsgToServerNew(serverId, "RpcAssistantNotify", notice)
}

//只推送给在线用户
func NoticeAssistant3(pid PLAYER_ID, id int64, title string, content string) {
	if !PlayerOnlineMgr.CheckPlayerIsOnLine(pid) { //如果不在线
		return
	}
	ep2 := ClientEpMp.LoadEndpoint(pid)
	if ep2 == nil {
		return
	}
	now := util.GetMilliTime()
	notice := &client_server.AssistantMsg{
		NoticeId:         easygo.NewInt64(id),
		MsgType:          easygo.NewInt32(3),
		DateTime:         &now,
		Title:            &title,
		SysNoticeContent: &content,
	}
	ep2.RpcAssistantNotify(notice)
}

//给电竞用户发系统消息
func NoticeAssistantEsportSysMsg(msg *share_message.TableESPortsSysMsg) {
	reqMsg := &client_hall.ESPortsSysMsgList{
		SysMsgList: []*share_message.TableESPortsSysMsg{msg},
	}
	methodName := "RpcESportNewSysMessage"
	recipientType := msg.GetRecipientType()                          //接收者类型 0 全体，1 IOS,2 Android
	var SendToPlyer = func(pinfo *Player, epx IGameClientEndpoint) { //发送给用户
		plyid := pinfo.GetPlayerId()
		obj := for_game.GetRedisESportPlayerObj(plyid)
		if obj != nil {
			reqMsg.PlayerId = easygo.NewInt64(plyid)
			_, err := epx.CallRpcMethod(methodName, reqMsg)
			if err == nil {
				obj.SetLastPullTime(easygo.NowTimestamp())
			}
		}
	}
	endPoints := ClientEpMgr.GetEndpoints() //所有在线用户
	for _, ep := range endPoints {
		if ep != nil {
			epx := ep.(IGameClientEndpoint)
			pinfo := epx.GetPlayer()
			if pinfo != nil {
				pDeviceType := pinfo.GetDeviceType() ////用户设备类型 1 IOS，2 Android，3 PC
				if recipientType == 0 {              //发送全部
					SendToPlyer(pinfo, epx)
				} else if recipientType == 1 && pDeviceType == 1 { //1 IOS, recipientType ==  pDeviceType
					SendToPlyer(pinfo, epx)
				} else if recipientType == 2 && pDeviceType == 2 { //2 Android, recipientType ==  pDeviceType
					SendToPlyer(pinfo, epx)
				}
			}

		}
	}
}
func BuildUnreadAssistantList(player *Player, login_type int32) []*client_server.AssistantMsg {
	if player == nil {
		return []*client_server.AssistantMsg{}
	}
	list := make([]*client_server.AssistantMsg, 0)
	t := player.GetLastAssistantTime()
	if t == 0 {
		t = player.GetCreateTime()
	}

	assistants := make([]*share_message.Assistant, 0)
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ASSISTANT)
	defer closeFun()
	e := col.Find(bson.M{"player_id": player.GetPlayerId(), "create_time": bson.M{"$gt": t}}).All(&assistants)
	if e != mgo.ErrNotFound && e != nil {
		logs.Error(e)
	}

	for _, value := range assistants {
		createTime := value.GetCreateTime()
		switch value.GetType() { //消息类型 0：好友 1：系统消息 2：客服反馈
		case 0:
			var types, sex int32
			if pp := for_game.GetRedisPlayerBase(value.GetPid()); pp != nil {
				types = pp.GetTypes()
				sex = pp.GetSex()
			}
			list = append(list, &client_server.AssistantMsg{
				LogId:         value.Id,
				MsgType:       value.Type,
				DateTime:      &createTime,
				AddPalType:    value.FriendType,
				NickName:      value.Nickname,
				Account:       value.Account,
				HeadIcon:      value.Avatar,
				Phone:         value.Phone,
				AddFriendType: easygo.NewInt32(value.GetAddFriendType()),
				Photo:         value.GetPhoto(),
				Signature:     easygo.NewString(value.GetSignature()),
				PlayerId:      easygo.NewInt64(value.GetPid()),
				Types:         easygo.NewInt32(types),
				Sex:           easygo.NewInt32(sex),
			})
		case 1:
			list = append(list, &client_server.AssistantMsg{LogId: value.Id, MsgType: value.Type, DateTime: &createTime,
				Title: value.Title, SysContent: value.Content})
		case 2:
			list = append(list, &client_server.AssistantMsg{LogId: value.Id, MsgType: value.Type, DateTime: &createTime,
				Title: value.Title, SysNoticeContent: value.Content})

		}
	}
	notices := make([]*share_message.SystemNotice, 0)
	col2, closeFun2 := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_NOTICE)
	defer closeFun2()
	e = col2.Find(bson.M{}).All(&notices)
	if e != mgo.ErrNotFound && e != nil {
		logs.Error(e)
	}
	var tt int32 = 3
	for _, value := range notices {
		if value.GetState() != 2 {
			continue
		}
		if value.GetUserType() != 0 && player.GetDeviceType() != value.GetUserType() {
			continue
		}
		//只发一次的，检测获取时间t是否大于发布时间
		if t > value.GetEditTime() && value.GetCountType() == 0 {
			continue
		}
		//检测每次重新发的
		if value.GetCountType() == 1 && login_type == 0 {
			continue
		}
		//if !(t > value.GetEditTime() && value.GetCountType() == 0) || !(value.GetCountType() == 1 && login_type == 1) {
		//	continue
		//}
		if value.GetType() == 0 {
			createTime := value.GetCreateTime()
			list = append(list, &client_server.AssistantMsg{NoticeId: value.Id, MsgType: &tt, DateTime: &createTime,
				Title: value.Title, SysNoticeContent: value.Content})
		} else {
			editTime := value.GetEditTime()
			list = append(list, &client_server.AssistantMsg{NoticeId: value.Id, MsgType: &tt, DateTime: &editTime,
				Title: value.Title, SysNoticeContent: value.Content})

		}
	}
	return list
}

// 优化前的代码,稳定后删除
//func BuildUnreadAssistantList(player *Player, login_type int32) []*client_server.AssistantMsg {
//	if player == nil {
//		return []*client_server.AssistantMsg{}
//	}
//	list := make([]*client_server.AssistantMsg, 0)
//	t := player.GetLastAssistantTime()
//
//	if t == 0 {
//		t = player.GetCreateTime()
//	}
//
//	assistants := make([]*share_message.Assistant, 0)
//	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ASSISTANT)
//	defer closeFun()
//	e := col.Find(bson.M{"player_id": player.GetPlayerId(), "create_time": bson.M{"$gt": t}}).All(&assistants)
//	if e != mgo.ErrNotFound && e != nil {
//		logs.Error(e)
//	}
//
//	for _, value := range assistants {
//		createTime := value.GetCreateTime()
//		if value.GetType() == 0 {
//			var types int32
//			if pp := for_game.GetRedisPlayerBase(value.GetPid()); pp != nil {
//				types = pp.GetTypes()
//			}
//			list = append(list, &client_server.AssistantMsg{
//				LogId:         value.Id,
//				MsgType:       value.Type,
//				DateTime:      &createTime,
//				AddPalType:    value.FriendType,
//				NickName:      value.Nickname,
//				Account:       value.Account,
//				HeadIcon:      value.Avatar,
//				Phone:         value.Phone,
//				AddFriendType: easygo.NewInt32(value.GetAddFriendType()),
//				Photo:         value.GetPhoto(),
//				Signature:     easygo.NewString(value.GetSignature()),
//				PlayerId:      easygo.NewInt64(value.GetPid()),
//				Types:         easygo.NewInt32(types),
//			})
//		}
//
//		if value.GetType() == 1 {
//			list = append(list, &client_server.AssistantMsg{LogId: value.Id, MsgType: value.Type, DateTime: &createTime,
//				Title: value.Title, SysContent: value.Content})
//		}
//
//		if value.GetType() == 2 {
//			list = append(list, &client_server.AssistantMsg{LogId: value.Id, MsgType: value.Type, DateTime: &createTime,
//				Title: value.Title, SysNoticeContent: value.Content})
//		}
//
//	}
//	notices := make([]*share_message.SystemNotice, 0)
//	col2, closeFun2 := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_NOTICE)
//	defer closeFun2()
//	e = col2.Find(bson.M{}).All(&notices)
//	if e != mgo.ErrNotFound && e != nil {
//		logs.Error(e)
//	}
//	var tt int32 = 3
//	for _, value := range notices {
//		if value.GetState() != 2 {
//			continue
//		}
//		if value.GetUserType() == 0 || player.GetDeviceType() == value.GetUserType() {
//			if (t > value.GetEditTime() && value.GetCountType() == 0) || (value.GetCountType() == 1 && login_type == 1) {
//				if value.GetType() == 0 {
//					createTime := value.GetCreateTime()
//					list = append(list, &client_server.AssistantMsg{NoticeId: value.Id, MsgType: &tt, DateTime: &createTime,
//						Title: value.Title, SysNoticeContent: value.Content})
//				} else {
//					editTime := value.GetEditTime()
//					list = append(list, &client_server.AssistantMsg{NoticeId: value.Id, MsgType: &tt, DateTime: &editTime,
//						Title: value.Title, SysNoticeContent: value.Content})
//
//				}
//			}
//		}
//	}
//	return list
//}

//生成金币变化扩展日志
func GetExtendLog(order *share_message.Order) *share_message.RechargeExtend {
	extendLog := &share_message.RechargeExtend{
		PayChannel:  easygo.NewInt32(order.GetPayChannel()),
		PayType:     easygo.NewInt32(order.GetPayType()),
		Channeltype: easygo.NewInt32(order.GetChanneltype()),
		CreateIP:    easygo.NewString(order.GetCreateIP()),
		OrderId:     easygo.NewString(order.GetOrderId()),
		Amount:      easygo.NewInt64(order.GetAmount()),
		Operator:    easygo.NewString(order.GetOperator()),
	}
	return extendLog
}

func GetCallTime(ti int64) string {
	min := ti / 60
	sec := ti % 60
	var smin, ssec string
	if min < 10 {
		smin = fmt.Sprintf("0%d", min)
	} else {
		smin = strconv.Itoa(int(min))
	}
	if sec < 10 {
		ssec = fmt.Sprintf("0%d", sec)
	} else {
		ssec = strconv.Itoa(int(sec))
	}
	sti := fmt.Sprintf("%s:%s", smin, ssec)
	return sti
}

//文本检测，content为base64编码字符串
func TextModeration(content string, talker, targetId int64, isSave ...bool) int32 {
	save := append(isSave, false)[0]
	res := for_game.TextModeration(content)
	if res == nil {
		return 0
	}
	logs.Info("文本检测结果:", res)
	if res.EvilFlag != 0 || res.EvilType != 100 {
		score := PSysParameterMgr.GetTextModeration(res.EvilType)

		curScore := int32(0)
		for _, v := range res.DetailResult {
			if v.EvilType == res.EvilType {
				curScore = v.Score
				break
			}
		}
		if curScore > score {
			if save { //需要记录备案的
				saveContent, err := base64.StdEncoding.DecodeString(content)
				easygo.PanicError(err)
				log := share_message.PlayerTalkLog{
					PlayerId: easygo.NewInt64(talker),
					TargetId: easygo.NewInt64(targetId),
					Connect:  easygo.NewString(string(saveContent)),
					Words:    res.Keywords,
					EvilType: easygo.NewInt32(res.EvilType),
				}
				col, closeFun := easygo.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PLAYER_TALK_LOG)
				defer closeFun()
				err = col.Insert(log)
				easygo.PanicError(err)
			}
			return res.EvilType
		} else {
			return 100
		}
	}
	return 100
}

//图片检测
func ImageModeration(url string, talker, targetId int64, isSave ...bool) int32 {
	save := append(isSave, false)[0]
	res := for_game.ImageModeration(url)
	if res == nil {
		return 100
	}
	logs.Info(" 图片检测结果:", res)
	if res.EvilFlag != 0 || res.EvilType != 100 {
		isGet := false
		score := PSysParameterMgr.GetImageModeration(res.EvilType)
		switch res.EvilType {
		case 20001:
			if res.PolityDetect.Score > score {
				isGet = true
			}
		case 20002:
			if res.PornDetect.Score > score {
				isGet = true
			}
		case 20006:
			if res.IllegalDetect.Score > score {
				isGet = true
			}
		case 20103:
			if res.HotDetect.Score > score {
				isGet = true
			}
		case 24001:
			if res.TerrorDetect.Score > score {
				isGet = true
			}
		}
		if isGet {
			if save {
				log := share_message.PlayerTalkLog{
					PlayerId:   easygo.NewInt64(talker),
					TargetId:   easygo.NewInt64(targetId),
					Connect:    easygo.NewString(string(url)),
					EvilType:   easygo.NewInt32(res.EvilType),
					CreateTime: easygo.NewInt64(time.Now().Unix()),
				}
				col, closeFun := easygo.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PLAYER_TALK_LOG)
				defer closeFun()
				err := col.Insert(log)
				easygo.PanicError(err)
			}
			return res.EvilType
		}
	}
	return 100
}

//零钱变化
func NotifyAddGold(playerId int64, value int64, reason string, sourceType int32, extendLog interface{}) string {
	player := for_game.GetRedisPlayerBase(playerId)
	result := player.TryAddGold(value, reason, sourceType, extendLog)
	msg := player.GetMyInfo()
	msg.Diamond = easygo.NewInt64(GetDiamondFromWishServer(playerId))
	b := easygo.PCall(SendMsgToHallClientNew, playerId, "RpcPlayerAttrChange", msg)
	if !b {
		logs.Error("通知前端失败")
	}
	return result
}

//硬币变化
func NotifyAddCoin(playerId int64, value int64, reason string, sourceType int32, extendLog interface{}, isUseCoin ...bool) string {
	player := for_game.GetRedisPlayerBase(playerId)
	if player == nil {
		logs.Error("硬币变化,获取用户信息失败,用户id为: %d,reason: %s", playerId, reason)
		return "用户不存在"
	}
	result := player.AddCoin(value, reason, sourceType, extendLog, isUseCoin...)
	msg := player.GetMyInfo()
	msg.Diamond = easygo.NewInt64(GetDiamondFromWishServer(playerId))
	b := easygo.PCall(SendMsgToHallClientNew, playerId, "RpcPlayerAttrChange", msg)
	if !b {
		logs.Error("通知前端失败")
	}
	return result
}

//硬币变化扩展
func NotifyAddCoinEx(playerId int64, value int64, reason string, sourceType int32, extendLog interface{}, diamond int64, isUseCoin ...bool) string {
	player := for_game.GetRedisPlayerBase(playerId)
	if player == nil {
		logs.Error("硬币变化,获取用户信息失败,用户id为: %d,reason: %s", playerId, reason)
		return "用户不存在"
	}
	result := player.AddCoin(value, reason, sourceType, extendLog, isUseCoin...)
	msg := player.GetMyInfo()
	msg.Diamond = easygo.NewInt64(diamond)
	b := easygo.PCall(SendMsgToHallClientNew, playerId, "RpcPlayerAttrChange", msg)
	if !b {
		logs.Error("通知前端失败")
	}
	return result
}

//电竞币变化
func NotifyAddESportCoin(playerId int64, value int64, reason string, sourceType int32, extendLog interface{}) string {
	player := for_game.GetRedisPlayerBase(playerId)
	result := player.TryAddESportCoin(value, reason, sourceType, extendLog)
	msg := player.GetMyInfo()
	msg.Diamond = easygo.NewInt64(GetDiamondFromWishServer(playerId))
	b := easygo.PCall(SendMsgToHallClientNew, playerId, "RpcPlayerAttrChange", msg)
	if !b {
		logs.Error("通知前端失败")
	}
	return result
}

//玩家增加道具,并通知变化 bugWay:1-硬币购买,2-现金购买 orderId:产生的订单id
func NotifyAddBagItems(playerId int64, items []*share_message.CoinProduct, way, bugWay int32, orderId string, operator string, givePlayerId ...int64) {
	bagObj := for_game.GetRedisPlayerBagItemObj(playerId)
	newItems := bagObj.AddItems(items, way, bugWay, orderId, operator, givePlayerId...)
	//检测新装备是否自动安装
	types := []int32{}
	for _, it := range newItems {
		cfg := for_game.GetPropsItemInfo(it.GetPropsId())
		if cfg == nil {
			logs.Info("找不到道具配置:", it.GetId())
			continue
		}
		types = append(types, cfg.GetPropsType())
		if cfg.GetUseType() == for_game.COIN_PROPS_USETYPE_EQUIPMENT {
			equipmentObj := for_game.GetRedisPlayerEquipmentObj(playerId)
			curId := equipmentObj.GetCurEquipment(cfg.GetPropsType())
			if curId == 0 {
				//该部位为空，自动安装
				it.Status = easygo.NewInt32(for_game.COIN_BAG_ITEM_USED)
				bagObj.SetItemStatus(it.GetId(), for_game.COIN_BAG_ITEM_USED)
				equipmentObj.Equipment(cfg.GetPropsType(), it.GetId())
				//通知前端装备修改
				newEquipment := equipmentObj.GetEquipmentForClient()
				SendMsgToClient([]int64{playerId}, "RpcModifyEquipment", newEquipment)

			}
		}
	}
	//通知前端道具修改
	bMsg := &client_hall.BagItems{
		Items: newItems,
	}
	SendMsgToClient([]int64{playerId}, "RpcModifyBagItem", bMsg)
	//通知前端有新道具
	notice := &client_hall.NewBagItemsTip{Types: types}
	SendMsgToClient([]int64{playerId}, "RpcNewBagItemsTip", notice)
}

//道具回收:t回收时间，-1表示永久
func NoticeReduceBagItem(playerId, netId, t int64) {
	bagObj := for_game.GetRedisPlayerBagItemObj(playerId)
	item := bagObj.ReduceUseTime(netId, t)
	bMsg := &client_hall.BagItems{
		Items: []*share_message.PlayerBagItem{item},
	}
	SendMsgToClient([]int64{playerId}, "RpcModifyBagItem", bMsg)
	if item.GetStatus() == for_game.COIN_BAG_ITEM_EXPIRED {
		//通知用户装备改变
		equipmentObj := for_game.GetRedisPlayerEquipmentObj(playerId)
		equipment := equipmentObj.GetEquipmentForClient()
		SendMsgToClient([]int64{playerId}, "RpcModifyEquipment", equipment)
	}
}

//发送消息给客户端:包括本服客户端，和其他服
func SendMsgToClient(playerIds []PLAYER_ID, methodName string, msg easygo.IMessage) {
	serverId := PServerInfo.GetSid()
	info := make(map[int32][]PLAYER_ID)
	for _, pid := range playerIds { //群发 每个人都发
		if PlayerOnlineMgr.CheckPlayerIsOnLine(pid) { //如果在线
			// 统计成字典 然后统一发
			serverId := PlayerOnlineMgr.GetPlayerServerId(pid)
			info[serverId] = append(info[serverId], pid)
		}
	}
	for sId, lst := range info {
		if sId == serverId {
			SendToCurrentHallClient(lst, methodName, msg)
		} else {
			BroadCastMsgToHallClientNew(lst, methodName, msg)
		}
	}
}

//发送给本服客户端
func SendToCurrentHallClient(playerIds []PLAYER_ID, methodName string, msg easygo.IMessage) {
	for _, pid := range playerIds {
		ep := ClientEpMp.LoadEndpoint(pid)
		if ep != nil {
			_, err2 := ep.CallRpcMethod(methodName, msg)
			easygo.PanicError(err2)
		}
	}
}

//发送给本服客户端
func SendToCurrentHallClientEx(playerId PLAYER_ID, methodName string, msg easygo.IMessage) bool {
	pid := playerId
	ep := ClientEpMp.LoadEndpoint(pid)
	if ep != nil {
		_, err2 := ep.CallRpcMethod(methodName, msg)
		if err2 != nil {
			logs.Error(err2)
		} else {
			return true
		}
	}
	return false
}

//服务器间通讯通用
func SendMsgToServerNew(sid SERVER_ID, methodName string, msg easygo.IMessage, pid ...int64) (easygo.IMessage, *base.Fail) {
	srv := PServerInfoMgr.GetServerInfo(sid)
	if srv == nil {
		return nil, easygo.NewFailMsg("无效的服务器id =" + easygo.AnytoA(sid))
	}
	var msgByte []byte
	if msg != nil {
		b, err := msg.Marshal()
		easygo.PanicError(err)
		msgByte = b
	} else {
		msgByte = []byte{}
	}
	playerId := append(pid, 0)[0]
	req := &share_message.MsgToServer{
		PlayerId: easygo.NewInt64(playerId),
		RpcName:  easygo.NewString(methodName),
		MsgName:  easygo.NewString(proto.MessageName(msg)),
		Msg:      msgByte,
	}
	return PWebApiForServer.SendToServer(srv, "RpcMsgToOtherServer", req)
}

//广播给指定类型服务器
func BroadCastMsgToServerNew(t int32, methodName string, msg easygo.IMessage, pid ...int64) {
	servers := PServerInfoMgr.GetAllServers(t)
	for _, srv := range servers {
		if srv == nil {
			continue
		}
		if srv.GetSid() == PServerInfo.GetSid() {
			continue
		}
		var msgByte []byte
		if msg != nil {
			b, err := msg.Marshal()
			easygo.PanicError(err)
			msgByte = b
		} else {
			msgByte = []byte{}
		}

		playerId := append(pid, 0)[0]
		req := &share_message.MsgToServer{
			PlayerId: easygo.NewInt64(playerId),
			RpcName:  easygo.NewString(methodName),
			MsgName:  easygo.NewString(proto.MessageName(msg)),
			Msg:      msgByte,
		}
		PWebApiForServer.SendToServer(srv, "RpcMsgToOtherServer", req)
	}
}

//广播给其他大厅客户端
func BroadCastMsgToHallClientNew(playerIds []int64, methodName string, msg easygo.IMessage, pid ...int64) {
	servers := PServerInfoMgr.GetAllServers(for_game.SERVER_TYPE_HALL)
	//if len(servers) == 0 && PServerInfo.GetType() == for_game.SERVER_TYPE_HALL {
	//	SendToCurrentHallClient(playerIds, methodName, msg)
	//	return
	//}
	//组装发送消息体
	var msgByte []byte
	if msg != nil {
		b, err := msg.Marshal()
		easygo.PanicError(err)
		msgByte = b
	} else {
		msgByte = []byte{}
	}

	req := &share_message.MsgToClient{
		PlayerIds: playerIds,
		RpcName:   easygo.NewString(methodName),
		MsgName:   easygo.NewString(proto.MessageName(msg)),
		Msg:       msgByte,
	}
	//把本服也增加上
	servers = append(servers, PServerInfo)
	//广播所有大厅服务器
	for _, srv := range servers {
		if srv == nil {
			continue
		}
		PWebApiForServer.SendToServer(srv, "RpcMsgToHallClient", req)
	}
}

//发送给指定玩家发送,在本服和不在本服
func SendMsgToHallClientNew(pid int64, methodName string, msg easygo.IMessage) {
	player := for_game.GetRedisPlayerBase(pid)
	if player == nil || !player.GetIsOnLine() {
		logs.Error("玩家不在线或者不存在")
		return
	}
	sid := player.GetSid()
	if sid == PServerInfo.GetSid() {
		SendToCurrentHallClient([]int64{pid}, methodName, msg)
		return
	}
	srv := PServerInfoMgr.GetServerInfo(sid)
	if srv == nil {
		logs.Error("找不到服务器:sid=", sid)
		return
	}
	var msgByte []byte
	if msg != nil {
		b, err := msg.Marshal()
		easygo.PanicError(err)
		msgByte = b
	} else {
		msgByte = []byte{}
	}

	req := &share_message.MsgToClient{
		PlayerIds: []int64{pid},
		RpcName:   easygo.NewString(methodName),
		MsgName:   easygo.NewString(proto.MessageName(msg)),
		Msg:       msgByte,
	}
	PWebApiForServer.SendToServer(srv, "RpcMsgToHallClient", req)
}

//与当前版本比较大小
func CompareVersion(ver string) bool {
	src := strings.Split(ver, ".")
	target := strings.Split(PServerInfo.GetVersion(), ".")
	for i := 0; i < len(src); i++ {
		a := easygo.AtoInt32(src[i])
		b := easygo.AtoInt32(target[i])
		if a != b {
			return a < b
		}
	}
	return true
}

// 取两个切片的交集
func Intersect(slice1, slice2 []int32) []int32 {
	m := make(map[int32]int)
	n := make([]int32, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
			m[v]++ //取消重复
		}
	}
	return n
}

// 两个玩家之间的匹配度，返回匹配值，共同标签
func GetPlayerMatchDegree(player1 *for_game.RedisPlayerBaseObj, player2 *share_message.PlayerBase) (int32, []int32) {
	matchDegree := player1.GetMatchingDegree(player2.GetPlayerId())
	commonTags := Intersect(player1.GetPersonalityTags(), player2.GetPersonalityTags())
	return matchDegree, commonTags
}

//发送给随机电竞接口服务器
func SendMsgToESportsApply(methodName string, msg easygo.IMessage) {
	srv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_SPORT_APPLY)
	if srv == nil {
		logs.Error("电竞服务器尚未开启")
		return
	}
	PWebApiForServer.SendToServer(srv, methodName, msg)
}

//发送给指定玩家发送(单发送)
func SendMsgToHallClientNewEx(playerId int64, methodName string, msg easygo.IMessage) bool {

	if playerId < 1 {
		return false
	}
	player := for_game.GetRedisPlayerBase(playerId)
	if player == nil || !player.GetIsOnLine() {
		return false
	}
	sid := player.GetSid()
	if sid == PServerInfo.GetSid() {
		//在本服，直接发送
		SendToCurrentHallClientEx(playerId, methodName, msg)
	} else {
		srv := PServerInfoMgr.GetServerInfo(sid)
		if srv == nil {
			return false
		}
		var msgByte []byte
		if msg != nil {
			b, err := msg.Marshal()
			easygo.PanicError(err)
			msgByte = b
		} else {
			msgByte = []byte{}
		}

		req := &share_message.MsgToClient{
			PlayerIds: []int64{playerId},
			RpcName:   easygo.NewString(methodName),
			MsgName:   easygo.NewString(proto.MessageName(msg)),
			Msg:       msgByte,
		}
		PWebApiForServer.SendToServer(srv, "RpcMsgToHallClient", req)

	}
	return true
}

//检查是否支持的银行卡
func CheckSupportBank(channel int32, cardCode string) bool {
	clt := &cls1{}
	supportReq := &client_hall.SupportBankList{
		PayId: easygo.NewInt32(channel),
	}
	clt.RpcGetSupportBankList(nil, nil, supportReq)
	for _, v := range supportReq.GetBanks() {
		if v.GetCode() == cardCode {
			return true
		}
	}
	return false
}
