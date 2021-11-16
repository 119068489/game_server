// 大厅发给子游戏服的消息处理

package hall

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"reflect"

	"github.com/astaxie/beego/logs"
)

//===================================================================

type ServiceForShop struct {
	Service reflect.Value
}

//分发消息
func (self *ServiceForShop) RpcServerReport(common *base.Common, reqMsg *share_message.ServerInfo) easygo.IMessage {
	logs.Info("RpcServerReport:", reqMsg)
	//ep.SetServerId(reqMsg.GetSid())
	//ShopEpMgr.Store(reqMsg.GetSid(), ep)
	PServerInfoMgr.AddServerInfo(reqMsg)
	//TODO:对在线的玩家重新分配连接
	//endpoints := ClientEpMgr.GetEndpoints()
	//for _, p := range endpoints {
	//	player := p.GetPlayer()
	//	if player != nil {
	//		srv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_SHOP)
	//		if srv != nil {
	//			player.SetShopServerId(srv.GetSid())
	//		}
	//	}
	//}
	return nil
}

/*func (self *ServiceForShop) RpcShopItemMessageNotify(common *base.Common, reqMsg *share_message.ShopItemMessageInfoWithWho) easygo.IMessage {
	logs.Info("RpcShopItemMessageNotify")
	epx := ClientEpMp.LoadEndpoint(reqMsg.GetPlayerId())
	if epx != nil {
		//发送到前端
		msg := &share_message.ShopItemMessageInfo{
			Type:        reqMsg.Type,
			ShopMessage: reqMsg.GetShopMessage()}
		epx.CallRpcMethod("RpcShopItemMessageNotify", msg)
	}

	return nil
}*/

/*func (self *ServiceForShop) RpcShopOrderNotify(common *base.Common, reqMsg *share_message.ShopOrderNotifyInfoWithWho) easygo.IMessage {
	logs.Info("RpcShopOrderNotify")
	epx := ClientEpMp.LoadEndpoint(reqMsg.GetPlayerId())
	if epx != nil {
		//发送到前端
		msg := &share_message.ShopOrderNotifyInfo{
			OrderId: easygo.NewInt64(reqMsg.GetOrderId())}
		epx.CallRpcMethod("RpcShopOrderNotify", msg)
	}
	return nil
}
*/
func (self *ServiceForShop) RpcSendShopOrderExpress(common *base.Common, reqMsg *share_message.QueryExpressInfosResult) easygo.IMessage {
	logs.Info("RpcSendShopOrderExpress")
	//发送到前端
	var expressInfos []*server_server.ShopOrderExpressBody = []*server_server.ShopOrderExpressBody{}
	if reqMsg.GetExpressInfos() != nil {
		for _, value := range reqMsg.GetExpressInfos() {
			expressInfos = append(expressInfos, &server_server.ShopOrderExpressBody{
				DateTime: easygo.NewString(value.GetDateTime()),
				Remark:   easygo.NewString(value.GetRemark()),
			})
		}
	}
	msg := &server_server.ShopOrderExpressInfos{
		ExpressInfos: expressInfos,
		ExpressPhone: easygo.NewString(reqMsg.GetExpressPhone()),
		ExpressName:  easygo.NewString(reqMsg.GetExpressName()),
		UserId:       easygo.NewInt64(reqMsg.GetUserId()),
	}

	// ChooseOneBackStage("RpcSendShopOrderExpress", msg)
	admin := for_game.GetRedisAdmin(reqMsg.GetUserId())
	if admin != nil {
		sid := admin.ServerId
		SendMsgToServerNew(sid, "RpcSendShopOrderExpress", msg)
	}

	return nil
}

func (self *ServiceForShop) RpcShopPaySeller(common *base.Common, reqMsg *share_message.PaySellerInfo) easygo.IMessage {
	logs.Info("======RpcShopPaySeller====== reqMsg=%v", reqMsg)
	HandlePaySeller(reqMsg)
	return easygo.EmptyMsg
}

func HandlePaySeller(reqMsg *share_message.PaySellerInfo) { //卖家收款
	//确认收货
	if reqMsg.GetPayType() == 0 {

		receiver := for_game.GetRedisPlayerBase(reqMsg.GetReceiverId())
		sponsor := GetPlayerObj(reqMsg.GetSponsor_Id())
		gold := int64(reqMsg.GetMoney())

		s := fmt.Sprintf("确认收货卖家收款：商城订单号(%v)，卖家收取钱总额(%v)，卖家ID(%v)，卖家昵称(%s)，买家ID(%v)，买家昵称(%s)",
			reqMsg.GetOrderId(),
			gold,
			reqMsg.GetSponsor_Id(),
			sponsor.GetNickName(),
			reqMsg.GetReceiverId(),
			receiver.GetNickName(),
		)

		logs.Info(s)

		//orderId, _ := for_game.PlaceOrder(reqMsg.GetSponsor_Id(), gold, for_game.GOLD_TYPE_SHOP_ITEM_MONEY)
		orderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, for_game.GOLD_TYPE_SHOP_ITEM_MONEY)
		info := map[string]interface{}{
			"ReceiverName": receiver.GetNickName(),
		}

		MerchantId := easygo.IntToString(int(reqMsg.GetOrderId()))
		reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_SHOP_ITEM_MONEY, reqMsg.GetSponsor_Id(), info)
		msg := &share_message.GoldExtendLog{
			OrderId:    easygo.NewString(orderId),
			MerchantId: &MerchantId,
			PayType:    easygo.NewInt32(for_game.PAY_TYPE_GOLD),
			Title:      &reason,
			Gold:       easygo.NewInt64(sponsor.GetGold() + gold),
			Account:    easygo.NewString(receiver.GetAccount()),
		}
		NotifyAddGold(sponsor.GetPlayerId(), gold, reason, for_game.GOLD_TYPE_SHOP_ITEM_MONEY, msg)

		//后台取消
	} else if reqMsg.GetPayType() == 1 {

		receiver := GetPlayerObj(reqMsg.GetReceiverId())
		sponsor := for_game.GetRedisPlayerBase(reqMsg.GetSponsor_Id())

		gold := int64(reqMsg.GetMoney())

		s := fmt.Sprintf("后台取消退款给买家：商城订单号(%v)，退款总额(%v)，卖家ID(%v)，卖家昵称(%s)，买家ID(%v)，买家昵称(%s)",
			reqMsg.GetOrderId(),
			gold,
			reqMsg.GetSponsor_Id(),
			sponsor.GetNickName(),
			reqMsg.GetReceiverId(),
			receiver.GetNickName(),
		)

		logs.Info(s)
		//orderId, _ := for_game.PlaceOrder(reqMsg.GetSponsor_Id(), gold, for_game.GOLD_TYPE_SHOP_ITEM_MONEY)
		orderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, for_game.GOLD_TYPE_BACK_MONEY)
		info := map[string]interface{}{
			"ShopOrderID": reqMsg.GetOrderId(),
		}

		MerchantId := easygo.IntToString(int(reqMsg.GetOrderId()))
		reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_BACK_MONEY, reqMsg.GetReceiverId(), info)
		msg := &share_message.GoldExtendLog{
			OrderId:    easygo.NewString(orderId),
			MerchantId: &MerchantId,
			PayType:    easygo.NewInt32(for_game.PAY_TYPE_GOLD),
			Title:      &reason,
			Gold:       easygo.NewInt64(receiver.GetGold() + gold),
		}
		NotifyAddGold(receiver.GetPlayerId(), gold, reason, for_game.GOLD_TYPE_BACK_MONEY, msg)
	}

}

//心跳
func (self *ServiceForShop) RpcHeartBeat(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	//logs.Info("RpcHeartBeat ")
	return nil
}

//直接转发给客户端
/*func (self *ServiceForShop) RpcShop2HallClient(common *base.Common, reqMsg *share_message.MsgToClient) easygo.IMessage {
	logs.Info("======RpcShop2HallClient=======reqMsg=%v", reqMsg)
	for _, playerId := range reqMsg.GetPlayerIds() {
		ep := ClientEpMp.LoadEndpoint(playerId)
		if ep != nil {
			msg := easygo.NewMessage(reqMsg.GetMsgName())
			err := msg.Unmarshal(reqMsg.GetMsg())
			easygo.PanicError(err)
			_, err1 := ep.CallRpcMethod(reqMsg.GetRpcName(), msg)
			easygo.PanicError(err1)
		} else {
			logs.Info("--------------->玩家不在线了,pid: %d,reqMsg: %+v", playerId, reqMsg)
		}
	}
	return nil
}
*/
func (self *ServiceForShop) RpcShop2HallClient(common *base.Common, reqMsg *share_message.MsgToClient) easygo.IMessage {
	logs.Info("=====大厅服接收商城的下发 RpcShop2HallClient====,", reqMsg)
	for _, playerId := range reqMsg.GetPlayerIds() {
		//直接推送给客户端
		if reqMsg.GetIsSend() {
			ep := ClientEpMp.LoadEndpoint(playerId)
			if ep != nil {
				msg := easygo.NewMessage(reqMsg.GetMsgName())
				err := msg.Unmarshal(reqMsg.GetMsg())
				easygo.PanicError(err)
				_, err1 := ep.CallRpcMethod(reqMsg.GetRpcName(), msg)
				easygo.PanicError(err1)
			}
		} else {
			logs.Info("--------不用发给前端,自己内部转发-------")
			//指定用户处理逻辑
			methodName := reqMsg.GetRpcName()
			self.Service = reflect.ValueOf(self)
			method := self.Service.MethodByName(methodName)
			msg := easygo.NewMessage(reqMsg.GetMsgName())
			if !method.IsValid() || method.Kind() != reflect.Func {
				logs.Info("无效的rpc请求，找不到methodName:", methodName)
				return nil
			}
			args := make([]reflect.Value, 0, 2)
			args = append(args, reflect.ValueOf(common))
			//args = append(args, reflect.ValueOf(ctx))
			args = append(args, reflect.ValueOf(msg))
			backMsg := method.Call(args) // 分发到指定的rpc
			if backMsg != nil {
				return backMsg[0].Interface().(easygo.IMessage)
			}
		}
	}
	return nil
}
