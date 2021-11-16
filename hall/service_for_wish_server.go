package hall

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"
	"math"

	"github.com/astaxie/beego/logs"
)

const (
	WISH_ORDER_TYPE = 1 //许愿池订单类型
)

//TODO 接收分发方法，函数方法名自己定，对端发送时传什么，这里就接收什么,需要返回值，则直接返回，但在接收的地方要强转
func (self *WebHttpForServer) RpcCheckWishToken(common *base.Common, reqMsg *h5_wish.LoginReq) easygo.IMessage {
	logs.Info("RpcCheckWishToken:", reqMsg, common.GetUserId())
	resp := &h5_wish.LoginResp{
		Result: easygo.NewInt32(2),
	}
	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	if player != nil {
		if player.GetToken() == reqMsg.GetToken() {
			resp.Result = easygo.NewInt32(1)
			resp.Token = easygo.NewString(player.GetToken())
			resp.HallSid = easygo.NewInt32(player.GetSid())
			return resp
		}
	}
	return resp
}

func (self *WebHttpForServer) RpcAddCoin(common *base.Common, reqMsg *h5_wish.AddCoinReq) easygo.IMessage {
	logs.Info("RpcAddCoin:", reqMsg)
	resp := &h5_wish.AddCoinResp{
		Result: easygo.NewInt32(1),
	}
	var reason string
	switch reqMsg.GetSourceType() {
	case for_game.COIN_TYPE_WISH_ADD:
		reason = "许愿池回收"
	case for_game.COIN_TYPE_WISH_PRO_ADD:
		reason = "许愿池守护者收益"
	case for_game.COIN_TYPE_WISH_PAY:
		reason = "许愿池兑换钻石"
	case for_game.COIN_TYPE_WISH_PLATFORM_ADD:
		reason = "许愿池平台回收"
	case for_game.COIN_TYPE_WISH_DARE_BACK:
		reason = "许愿池抽奖返利"
	}

	err := NotifyAddCoinEx(reqMsg.GetUserId(), reqMsg.GetCoin(), reason, reqMsg.GetSourceType(), nil, reqMsg.GetDiamond())
	if err != "" {
		resp.Result = easygo.NewInt32(2)
	}
	if base := for_game.GetRedisPlayerBase(reqMsg.GetUserId()); base != nil {
		resp.Coin = easygo.NewInt64(base.GetCoin())
	}
	return resp
}

// 零钱变化  NotifyAddGold
func (self *WebHttpForServer) RpcAddGold(common *base.Common, reqMsg *h5_wish.AddGoldReq) easygo.IMessage {
	logs.Info("RpcAddGold:", reqMsg)
	resp := &h5_wish.AddCoinResp{
		Result: easygo.NewInt32(1),
	}
	var reason string
	switch reqMsg.GetSourceType() {
	case for_game.COIN_TYPE_WISH_ADD:
		reason = "许愿池回收"
	case for_game.COIN_TYPE_WISH_PLATFORM_ADD:
		reason = "许愿池平台回收"
	default:
		logs.Error("许愿池回收类型有误,sourceType: %d", for_game.COIN_TYPE_WISH_ADD)
		resp.Result = easygo.NewInt32(2)
		return resp
	}

	err := NotifyAddGold(reqMsg.GetUserId(), reqMsg.GetCoin(), reason, reqMsg.GetSourceType(), nil)
	if err != "" {
		resp.Result = easygo.NewInt32(2)
	}
	return resp
}
func (self *WebHttpForServer) RpcGetUseInfo(common *base.Common, reqMsg *h5_wish.UserInfoReq) easygo.IMessage {
	logs.Info("RpcGetUseInfo: %v", reqMsg)
	resp := &h5_wish.UserInfoResp{}
	player := for_game.GetRedisPlayerBase(reqMsg.GetUserId())
	var types int32
	if player != nil {
		resp.UserId = easygo.NewInt64(player.GetPlayerId())
		resp.Name = easygo.NewString(player.GetNickName())
		resp.Sex = easygo.NewInt32(player.GetSex())
		resp.HeadUrl = easygo.NewString(player.GetHeadIcon())
		resp.Coin = easygo.NewInt64(player.GetAllCoin())
		resp.Account = easygo.NewString(player.GetPhone())
		if player.GetTypes() == 1 {
			types = 0
		} else if player.GetTypes() > 1 {
			types = player.GetTypes()
		}
		resp.Types = easygo.NewInt32(types)
	}
	return resp
}

// 获取用戶银行卡信息
func (self *WebHttpForServer) RpcGetUserBankCardsInfo(common *base.Common, reqMsg *h5_wish.UserInfoReq) easygo.IMessage {
	logs.Info("RpcGetPlayerBankCardsInfo: %v", reqMsg)
	player := for_game.GetRedisPlayerBase(reqMsg.GetUserId())
	cards := make([]*h5_wish.BankCardInfo, 0)
	channel := PPayChannelMgr.GetCurPayChannel()
	for i := len(player.GetBankInfos()) - 1; i >= 0; i-- {
		bankInfo := player.GetBankInfos()[i]
		cards = append(cards, &h5_wish.BankCardInfo{
			BankId:    easygo.NewString(bankInfo.GetBankId()),
			BankCode:  easygo.NewString(bankInfo.GetBankCode()),
			BankName:  easygo.NewString(bankInfo.GetBankName()),
			IsSupport: easygo.NewBool(CheckSupportBank(channel.GetId(), bankInfo.GetBankCode())),
		})
	}
	return &h5_wish.BankCardResp{Cards: cards}
}

// 检查用户银行卡是否有效
func (self *WebHttpForServer) RpcCheckBankCardInfo(common *base.Common, reqMsg *h5_wish.RecycleToHall) easygo.IMessage {
	logs.Info("RpcGetPlayerBankCardsInfo: %v", reqMsg)
	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	if player == nil {
		logs.Error("检查用户银行卡是否有效时用户信息为空")
		return easygo.NewFailMsg("用户信息有误")
	}
	bankInfo := player.GetBankInfo(reqMsg.GetBankCardId())
	if nil == bankInfo {
		return easygo.NewFailMsg("玩家并未绑定该银行卡")
	}
	channel := PPayChannelMgr.GetCurPayChannel()
	if !CheckSupportBank(channel.GetId(), bankInfo.GetBankCode()) {
		return easygo.NewFailMsg("该银行卡并不支持回收")
	}
	return nil
}

//回收请求
func (self *WebHttpForServer) RpcRecycleToBandCard(common *base.Common, reqMsg *h5_wish.RecycleToHall) easygo.IMessage {
	logs.Info("RpcRecycleToBandCard", reqMsg)
	//if !for_game.IS_FORMAL_SERVER {
	//	return easygo.NewFailMsg("回收失败,测试服不能回收")
	//}
	//青少年保护模式
	dealReqMsg := &client_hall.WithdrawInfo{
		AccountNo: easygo.NewString(reqMsg.GetBankCardId()),
		Amount:    easygo.NewInt64(reqMsg.GetPrice()),
	}
	who := GetPlayerObj(reqMsg.GetPlayerId())
	if who == nil {
		return easygo.NewFailMsg("用户信息有误")
	}
	if who.GetYoungPassWord() != "" {
		return easygo.NewFailMsg("青少年模式下无法回收")
	}

	config := for_game.GetConfigWishPayment()
	if !config.GetStatus() {
		return easygo.NewFailMsg("回收失败，回收功能暂未开启")
	}

	//8：00-22：00才可以回收
	_, _, _, h, _, _ := easygo.GetTimeData()
	if h < 8 || h >= 22 {
		return easygo.NewFailMsg("回收失败,每天8:00-22:00才可以回收")
	}

	//许愿池取消限制单笔最大最小限额
	/*
		sysConfig := PSysParameterMgr.GetSysParameter(for_game.LIMIT_PARAMETER)
		if !sysConfig.GetIsWithdrawal() {
			return easygo.NewFailMsg("回收入口已关闭，回收失败")
		}
		//每笔最小限额
		if gold < sysConfig.GetWithdrawalMin() { //最小10元
			return easygo.NewFailMsg("回收失败,金额不能小于" + easygo.AnytoA(float64(sysConfig.GetWithdrawalMin())/100.0) + "元")
		}
		//每笔最大限额
		if gold > sysConfig.GetWithdrawalMax() {
			return easygo.NewFailMsg("回收失败,金额不能大于" + easygo.AnytoA(float64(sysConfig.GetWithdrawalMax())/100.0) + "元")
		}
	*/

	if for_game.IS_FORMAL_SERVER {
		if who.GetPeopleId() == "" {
			res := "请先通过实名认证"
			return easygo.NewFailMsg(res)
		}
		if !who.GetIsPayPassword() {
			return easygo.NewFailMsg("回收失败,请先设置支付密码")
		}
	}
	bankId := dealReqMsg.GetAccountNo()
	bankInfo := who.GetBankMsg(bankId)
	if bankInfo == nil {
		return easygo.NewFailMsg("卡号未绑定，请绑定后在操作")
	}
	//当前激活的代付通道
	channel := PPayChannelMgr.GetCurPayChannel()
	logs.Info("代付channel----->%+v", channel)
	// 汇潮代付 需要判断是否有省份,城市
	if channel.GetId() == for_game.PAY_CHANNEL_HUICHAO_DF {
		if bankInfo.GetProvice() == "" || bankInfo.GetCity() == "" {
			return easygo.NewFailMsg("银行卡信息不完善", for_game.FAIL_MSG_CODE_1012)
		}
	}
	code := bankInfo.GetBankCode()
	name := for_game.BankName[code]
	if for_game.IS_FORMAL_SERVER && name == "" {
		return easygo.NewFailMsg("不存在的银行卡号")
	}
	if !CheckSupportBank(channel.GetId(), bankInfo.GetBankCode()) {
		return easygo.NewFailMsg("该银行卡并不支持回收")
	}

	gold := dealReqMsg.GetAmount() //回收金额
	taxGold := int64(0)            //手续费

	period := for_game.GetPlayerPeriod(who.GetPlayerId())
	//每日回收总额限制
	curSum := period.DayPeriod.FetchInt64(for_game.WISH_OUT_MONEY_SUM)
	//回收后台配置信息
	setting := &share_message.PaymentSetting{
		FeeRate:     easygo.NewInt32(config.GetFeeRate()),
		PlatformTax: easygo.NewInt64(config.GetPlatformTax()),
		RealTax:     easygo.NewInt64(config.GetRealTax()),
	}
	//setting := PPayChannelMgr.GetCurPaymentSetting()

	if gold+curSum > config.GetDayMoneyTop() {
		return easygo.NewFailMsg("回收失败，今天回收总额超" + easygo.AnytoA(float64(config.GetDayMoneyTop())/100.0) + "元")
	}
	//每日回收次数限制
	curTimes := period.DayPeriod.FetchInt32(for_game.WISH_OUT_MONEY_TIME)
	if curTimes >= config.GetDayMoneyTopCount() {
		return easygo.NewFailMsg("回收失败，今天回收次数已达" + easygo.AnytoA(config.GetDayMoneyTopCount()) + "次")
	}
	//需手续费:向上取整
	if setting.GetFeeRate() > 1 {
		taxGold = int64(math.Ceil(float64(gold) * float64(setting.GetFeeRate()) / 1000.0)) //手续费: X>2000/995 && X> 最小起提:15
	}

	// 如果是汇聚去到,需要校验回收的银行卡
	if channel.GetId() == for_game.PAY_CHANNEL_HUIJU_DF {
		payNo := for_game.BankPayNo[code]
		if payNo == "" {
			return easygo.NewFailMsg("暂不支持" + for_game.BankName[code] + "回收")
		}
		dealReqMsg.BankCode = easygo.NewString(payNo)
	}

	//发起代付订单
	dealReqMsg.AccountProp = easygo.NewString("0")
	dealReqMsg.AccountName = easygo.NewString(who.GetRealName())
	dealReqMsg.AccountType = easygo.NewString("00")
	dealReqMsg.Tax = easygo.NewInt64(taxGold)
	var orderId string
	logs.Info("channel.GetId() ------------->", channel.GetId())
	//生成内部回收订单
	if channel.GetId() == for_game.PAY_CHANNEL_PENGJU {
		orderId = PWebPengJuPay.CreateOrder(dealReqMsg, who, setting)
	} else if channel.GetId() == for_game.PAY_CHANNEL_HUIJU_DF {
		logs.Info("汇聚提现")
		orderId = PWebHuiJuPay.CreateDFOrder(dealReqMsg, who, setting)
	} else if channel.GetId() == for_game.PAY_CHANNEL_HUICHAO_DF {
		logs.Info("汇潮代付")
		orderId = PWebHuiChaoPay.CreateDFOrder(dealReqMsg, who, setting)
	}

	if orderId == "" {
		//订单号未生成，回收失败
		logs.Error("代付订单id为空")
		return easygo.NewFailMsg("系统异常，回收下单失败")
	}

	// 体现预警通知
	easygo.Spawn(for_game.CheckWishWarningSMS, dealReqMsg.GetAmount(), int64(0))

	dealReqMsg.OrderId = easygo.NewString(orderId)
	//该订单中的金币并无任何意义
	order := for_game.GetRedisOrderObj(orderId)
	order.SetSourceType(for_game.DIAMOND_TYPE_WISH_BACK)
	order.SetOrderType(WISH_ORDER_TYPE)
	order.SetNote("回收")

	errStr := ""
	restStatus := for_game.RECYCLE_STATUS_ERROR
	if gold < config.GetOrderThresholdMoney() {
		//第三方发起回收请求
		errMsg := &base.Fail{}
		if channel.GetId() == for_game.PAY_CHANNEL_PENGJU {
			errMsg = PWebPengJuPay.RechargeOrder(order.GetRedisOrder())
		} else if channel.GetId() == for_game.PAY_CHANNEL_HUIJU_DF {
			errMsg = PWebHuiJuPay.ReqSinglePay(order.GetRedisOrder())
		} else if channel.GetId() == for_game.PAY_CHANNEL_HUICHAO_DF {
			errMsg = PWebHuiChaoPay.TransferFixed(order.GetRedisOrder())
		}
		if errMsg.GetCode() == for_game.FAIL_MSG_CODE_SUCCESS {
			//第三方下单成功
			restStatus = for_game.RECYCLE_STATUS_SUCCESS
			errStr = errMsg.GetReason()
			logs.Info("回收下单成功:")
			order.SetPayStatus(for_game.PAY_ST_DOING)
		} else {
			//渠道下单失败，如果已经生成订单，则转为人工处理
			order.SetChanneltype(for_game.CHANNEL_MAN_MAKE)
			errStr = errMsg.GetReason()
			restStatus = for_game.RECYCLE_STATUS_CHECK
			logs.Info("回收第三方下单失败,转为人工处理订单:", orderId)
		}
	} else {
		//回收值超过风控，进入人工审核
		logs.Info("回收超过风控值:", config.GetOrderThresholdMoney())
		//order.Channeltype = easygo.NewInt32(for_game.CHANNEL_MAN_MAKE)
		errStr = "回收额度超过风控值，进行人工审核"
		restStatus = for_game.RECYCLE_STATUS_CHECK
	}
	order.SetExtendValue(errStr) // 成功/失败原因
	order.SaveToMongo()
	//增加当天回收次数与记录
	period.DayPeriod.AddInteger(for_game.WISH_OUT_MONEY_SUM, gold)
	period.DayPeriod.AddInteger(for_game.WISH_OUT_MONEY_TIME, 1)
	return &h5_wish.OrderMsgResp{
		OrderId: easygo.NewString(orderId),
		Status:  easygo.NewInt32(restStatus),
	}
}

// 小助手通知
func (self *WebHttpForServer) RpcWishNoticeAssistant(common *base.Common, reqMsg *h5_wish.WishNoticeAssistantReq) easygo.IMessage {
	logs.Info("==========大厅收到了许愿池小助手通知请求RpcWishNoticeAssistant: %v", reqMsg)
	pid := reqMsg.GetPlayerId()
	if pid == 0 {
		logs.Error("许愿池抽奖推送的柠檬id为空,请求参数为reqMsg: %v ", reqMsg)
		return easygo.NewFailMsg("推送失败")

	}
	productName := reqMsg.GetProductName()
	if productName == "" {
		logs.Error("许愿池抽奖推送的物品名字为空,请求参数为reqMsg: %v ", reqMsg)
		return easygo.NewFailMsg("推送失败")
	}
	// 过期小助手通知
	content := fmt.Sprintf("您在许愿池中抽中物品%s，快去我的愿望盒看看吧", productName)
	NoticeAssistant(pid, 1, "温馨提示", content)
	return nil
}
