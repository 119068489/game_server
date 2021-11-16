// 管理后台为[浏览器]提供的服务
//用户管理

package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"math"
	"net"

	"github.com/astaxie/beego/logs"
)

//查询通用支付额度设置
func (self *cls4) RpcQueryGeneralQuota(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	ret := for_game.GetGeneralQuota()

	return ret
}

//修改通用支付额度设置
func (self *cls4) RpcEditGeneralQuota(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.GeneralQuota) easygo.IMessage {
	msg := "修改通用支付额度设置"
	EditGeneralQuota(reqMsg)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询支付方式
func (self *cls4) RpcQueryPayType(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	list := QueryPayType(reqMsg)
	msg := &brower_backstage.PayTypeResponse{
		List: list,
	}

	return msg
}

//修改支付方式
func (self *cls4) RpcEditPayType(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PayType) easygo.IMessage {
	msg := "修改支付方式:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_PAYTYPE))
		msg = "添加支付方式:"
	}

	EditPayType(reqMsg)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除支付方式
func (self *cls4) RpcDelPayType(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	for _, id := range idList {
		pfc := QueryPlatformChannelByPid(int32(id), 3)
		if pfc != nil {
			return easygo.NewFailMsg("请先删除使用此方式的支付通道")
		}
	}
	err := DelDataById(for_game.TABLE_PAYTYPE, idList)
	easygo.PanicError(err)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除支付方式: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询支付场景
func (self *cls4) RpcQueryPayScene(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	list := QueryPayScene()

	return &brower_backstage.PaySceneResponse{
		List: list,
	}
}

//修改支付场景
func (self *cls4) RpcEditPayScene(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PayScene) easygo.IMessage {
	msg := "修改支付场景:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_PAYSCENE))
		msg = "添加支付场景:"
	}

	EditPayScene(reqMsg)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除支付场景
func (self *cls4) RpcDelPayScene(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	for _, id := range idList {
		pfc := QueryPlatformChannelByPid(int32(id), 2)
		if pfc != nil {
			return easygo.NewFailMsg("请先删除使用此场景的支付通道")
		}
	}
	err := DelDataById(for_game.TABLE_PAYSCENE, idList)
	easygo.PanicError(err)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除支付场景: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询支付设定
func (self *cls4) RpcQueryPaymentSetting(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := QueryPaymentSettingtList(reqMsg)

	return &brower_backstage.PaymentSettingResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//修改支付设定
func (self *cls4) RpcEditPaymentSetting(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PaymentSetting) easygo.IMessage {
	msg := "修改支付设定:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_PAYMENTSETTING))
		msg = "添加支付设定:"
	}

	EditPaymentSetting(reqMsg)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除支付设定
func (self *cls4) RpcDelPaymentSetting(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	for _, v := range idList {
		if v == 1 || v == 2 {
			easygo.NewFailMsg("不能删除ID为1和2的默认设置")
		}
	}
	err := DelDataById(for_game.TABLE_PAYMENTSETTING, idList)
	easygo.PanicError(err)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除支付设定: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询支付平台
func (self *cls4) RpcQueryPaymentPlatform(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.PlatformChannelRequest) easygo.IMessage {
	list, count := QueryPaymentPlatform(reqMsg)

	return &brower_backstage.PaymentPlatformResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//修改支付平台
func (self *cls4) RpcEditPaymentPlatform(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PaymentPlatform) easygo.IMessage {
	msg := "修改支付平台:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("ID不能为空")
		//reqMsg.Id = easygo.NewInt32(for_game.NextId(MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PAYMENTPLATFORM))
		//msg = "添加支付平台"
	} else {
		pf := QuerPaymentPlatformById(reqMsg.GetId())
		if pf == nil {
			msg = "添加支付平台:"
		}
	}

	EditPaymentPlatform(reqMsg)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除支付平台
func (self *cls4) RpcDelPaymentPlatform(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	for _, id := range idList {
		pfc := QueryPlatformChannelByPid(int32(id), 1)
		if pfc != nil {
			return easygo.NewFailMsg("请先删除此平台的支付通道")
		}
	}

	err := DelDataById(for_game.TABLE_PAYMENTPLATFORM, idList)
	easygo.PanicError(err)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除支付平台: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询支付平台通道列表
func (self *cls4) RpcQueryPlatformChannel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.PlatformChannelRequest) easygo.IMessage {
	list, count := QueryPlatformChannelList(reqMsg)
	reList := []*share_message.PlatformChannel{}
	for _, item := range list {
		item.PlatformName = QuerPaymentPlatformById(item.GetPlatformId()).Name
		item.PayTypeName = QueryPayTypeById(item.GetPayTypeId()).Name
		item.PaySceneName = QueryPaySceneById(item.GetPaySceneId()).Name
		item.PaymentSettingName = QueryPaymentSettingById(item.GetPaymentSettingId()).Name
		reList = append(reList, item)
	}

	return &brower_backstage.PlatformChannelResponse{
		List:      reList,
		PageCount: easygo.NewInt32(count),
	}
}

//修改支付平台通道
func (self *cls4) RpcEditPlatformChannel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PlatformChannel) easygo.IMessage {
	msg := "修改支付平台通道:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("ID不能为空")
		// reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_PLATFORM_CHANNEL))
		// msg = "添加支付平台通道:"
	}

	if reqMsg.PaymentSettingId == nil {
		easygo.NewFailMsg("支付设定必须选择")
	}

	if reqMsg.PayTypeId == nil {
		easygo.NewFailMsg("支付方式必须选择")
	}

	if reqMsg.PaySceneId == nil {
		easygo.NewFailMsg("支付场景必须选择")
	}

	if reqMsg.Types == nil {
		easygo.NewFailMsg("类型必须选择")
	}

	if reqMsg.PlatformId == nil {
		easygo.NewFailMsg("所属平台必须选择")
	}

	EditPlatformChannel(reqMsg)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除支付平台通道
func (self *cls4) RpcDelPlatformChannel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	err := DelDataById(for_game.TABLE_PLATFORM_CHANNEL, idList)
	easygo.PanicError(err)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除支付平台通道: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//批量修改通道状态 1开启，2关闭
func (self *cls4) RpcBatchClosePlatformChannel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	BatchClosePlatformChannel(reqMsg.GetIds32(), 2)
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcPaySetChangeToHall", nil) //通知大厅重载支付配置

	var ids string
	idsarr := reqMsg.GetIds32()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量关闭支付通道: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询现金变化类型
func (self *cls4) RpcQuerySouceType(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.SourceTypeRequest) easygo.IMessage {
	list := QuerySouceTypeList(reqMsg)
	msg := &brower_backstage.SourceTypeResponse{
		List: list,
	}
	return msg
}

//人工出入款
func (self *cls4) RpcAddGold(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.AddGoldResult) easygo.IMessage {
	if reqMsg.GetPlayerId() == 0 || reqMsg.PlayerId == nil {
		easygo.NewFailMsg("玩家ID不正确")
	}
	if reqMsg.GetGold() <= 0 {
		easygo.NewFailMsg("金额错误")
	}
	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	remoteAddr, _, err := net.SplitHostPort(ep.GetAddr().String()) //获取客户端IP
	easygo.PanicError(err)
	order := &share_message.Order{}
	status := easygo.NewInt32(0)
	PayStatus := easygo.NewInt32(0)
	tax := int64(0) //税收
	//操作上分
	amount := reqMsg.GetGold()
	changegold := amount
	//人工出入款修改支付状态为已完成
	if reqMsg.GetChanneltype() == 1 {
		PayStatus = easygo.NewInt32(1)
	}

	//出款操作
	if reqMsg.GetChangeType() == 2 {
		if (player.GetGold() - changegold) < 0 {
			return easygo.NewFailMsg("金额不足")
		}
		if reqMsg.GetSourceType() == for_game.GOLD_TYPE_CASH_OUT {
			paymentSetting := QueryPaymentSettingById(2)
			tax = int64(math.Ceil(float64(changegold) * (float64(paymentSetting.GetFeeRate()) / 1000.0)))

		}
		changegold = -(changegold - tax)
	}

	order = &share_message.Order{
		PlayerId:    reqMsg.PlayerId,
		Account:     easygo.NewString(player.GetAccount()),
		NickName:    easygo.NewString(player.GetNickName()),
		RealName:    easygo.NewString(player.GetRealName()),
		SourceType:  reqMsg.SourceType,
		ChangeType:  reqMsg.ChangeType,
		Channeltype: reqMsg.Channeltype,
		PayChannel:  easygo.NewInt32(0),
		PayType:     easygo.NewInt32(0),
		CurGold:     easygo.NewInt64(player.GetGold()),
		ChangeGold:  easygo.NewInt64(changegold),
		Gold:        easygo.NewInt64(player.GetGold() + changegold - tax),
		Amount:      easygo.NewInt64(amount),
		CreateTime:  easygo.NewInt64(for_game.GetMillSecond()),
		CreateIP:    easygo.NewString(remoteAddr),
		Status:      status,
		PayStatus:   PayStatus,
		Note:        reqMsg.Note,
		Tax:         easygo.NewInt64(-tax),
		Operator:    user.Account,
		PayWay:      easygo.NewInt32(for_game.PAY_TYPE_MONEY),
	}
	//下订单
	obj := for_game.CreateRedisOrder(order)

	//人工出入款直接完成订单
	if reqMsg.GetChanneltype() == 1 {
		//完成订单
		err := FinishOrder(obj.GetOrderId(), user.GetAccount(), "")
		if err != nil {
			return easygo.NewFailMsg(err.GetReason())
		}
	}

	s := fmt.Sprintf("为用户%s人工入款:%.2f", player.GetAccount(), easygo.Decimal(float64(changegold/int64(100)), 2))
	//提现操作
	if reqMsg.GetChangeType() == 2 {
		s = fmt.Sprintf("为用户%s人工出款:%.2f", player.GetAccount(), easygo.Decimal(float64(changegold/int64(100)), 2))
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, s)

	return easygo.EmptyMsg
}

//查询线上充值订单
func (self *cls4) RpcQueryOrderList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryOrderRequest) easygo.IMessage {
	logs.Info("RpcQueryOrderList:", reqMsg)
	list, count := QueryOrderList(reqMsg)

	return &brower_backstage.QueryOrderResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//更新订单列表
func (self *cls4) RpcUpdateOrderList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	for_game.SaveRedisOrderToMongo()
	return easygo.EmptyMsg
}

//订单操作
func (self *cls4) RpcOptOrder(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.OptOrderRequest) easygo.IMessage {
	if user.GetRole() > 0 {
		role := GetPowerRouter(user.GetRoleType())
		if for_game.IsContainsStr("orderManage", role.GetMenuIds()) == -1 {
			easygo.NewFailMsg("权限不足")
		}
	}
	oid := reqMsg.GetOid()
	if oid == "" {
		easygo.NewFailMsg("订单号错误")
	}

	order := for_game.GetRedisOrderObj(oid)
	if order == nil {
		easygo.NewFailMsg("订单不存在")
	}

	switch order.GetStatus() {
	case 1:
		easygo.NewFailMsg("订单已完成")
	case 2:
		easygo.NewFailMsg("订单已审核")
	case 3:
		easygo.NewFailMsg("订单已取消")
	case 4:
		easygo.NewFailMsg("订单已拒绝")
	}

	s := fmt.Sprintf("订单操作:%s", order.GetOrderId())
	var payNotice string
	switch reqMsg.GetOpt() {
	case 1: // 完成订单
		err := FinishOrder(order.GetOrderId(), user.GetAccount(), reqMsg.GetNote())
		if err != nil {
			return err
		}

		s = fmt.Sprintf("人工完成订单:%s", order.GetOrderId())
	case 2: //审核出款订单
		req := &server_server.AuditOrder{
			OrderId: easygo.NewString(oid),
		}
		ChooseOneHall(0, "RpcBsAuditOrder", req)
		s = fmt.Sprintf("审核订单:%s", order.GetOrderId())
	case 3, 4: // 取消/拒绝订单
		err := OptOrder(order.GetOrderId(), reqMsg.GetOpt(), user.GetAccount(), reqMsg.GetNote())
		if err != nil {
			return err
		}
		if reqMsg.GetOpt() == 3 {
			s = fmt.Sprintf("取消订单:%s", order.GetOrderId())
			payNotice = fmt.Sprintf("您在%s发起的提现订单已被取消，如有疑问请咨询官方客服", easygo.Stamp2Str(order.GetCreateTime()))
		} else {
			s = fmt.Sprintf("拒绝订单:%s", order.GetOrderId())
			payNotice = fmt.Sprintf("您在%s发起的提现订单已被拒绝，如有疑问请咨询官方客服", easygo.Stamp2Str(order.GetCreateTime()))
		}
		//提现操作通知
		if order.GetChangeType() == 2 {
			SendSystemNotice(order.GetPlayerId(), "支付通知", payNotice)
		}
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, s)

	return easygo.EmptyMsg
}

//人工补单(查询完成订单)
func (self *cls4) RpcCheckOrder(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.OptOrderRequest) easygo.IMessage {
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

	if for_game.CheckOrderDoing(oid) {
		return easygo.NewFailMsg("订单正在处理，稍后再试")
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

	result := SendToPlayer(order.GetPlayerId(), "RpcCheckOrderToHall", reqMsg) //通知大厅
	err := for_game.ParseReturnDataErr(result)
	if err != nil {
		return err
	}

	order.SetNote("人工补单")
	order.SetOperator(user.GetAccount())
	order.SaveToMongo() //保存Redis中的订单数据

	s := fmt.Sprintf("人工补单:%s", order.GetOrderId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, s)
	return easygo.EmptyMsg
}
