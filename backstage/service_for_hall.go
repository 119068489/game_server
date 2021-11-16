// 管理后台为[浏览器]提供的服务

package backstage

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"reflect"

	"github.com/astaxie/beego/logs"
)

type ServiceForHall struct {
	Service reflect.Value
}

func (self *ServiceForHall) RpcReplaceLogin(common *base.Common, reqMsg *server_server.PlayerSI) easygo.IMessage {
	logs.Info("RpcReplaceLogin ", reqMsg)
	ReplaceLogin(reqMsg.GetPlayerId())
	return nil
}

// 玩家上线通知
func (self *ServiceForHall) RpcPlayerOnLine(common *base.Common, reqMsg *share_message.PlayerState) easygo.IMessage {
	logs.Info("RpcPlayerOnLine ", reqMsg)
	PlayerOnlineMgr.PlayerOnline(reqMsg.GetPlayerId(), reqMsg.GetServerId())
	return nil
}

// 玩家离线通知
func (self *ServiceForHall) RpcPlayerOffLine(common *base.Common, reqMsg *share_message.PlayerState) easygo.IMessage {
	logs.Info("RpcPlayerOffLine ", reqMsg)
	PlayerOnlineMgr.PlayerOffline(reqMsg.GetPlayerId())
	return nil
}

//只用于通知广播，无返回
func (self *ServiceForHall) RpcHallToBackStage(common *base.Common, reqMsg *share_message.MsgData) easygo.IMessage {
	//收到大厅通知
	logs.Info("大厅发送过来的消息->RpcHallToBackStage:", reqMsg)
	//msg := easygo.NewMessage(reqMsg.GetMsgName())
	//err := msg.Unmarshal(reqMsg.GetMsg())
	//easygo.PanicError(err)
	//methodName := reqMsg.GetMsgName()
	//self.Service = reflect.ValueOf(self)
	//method := self.Service.MethodByName(methodName)
	//if !method.IsValid() || method.Kind() != reflect.Func {
	//	logs.Info("无效的rpc请求，找不到methodName:", methodName)
	//	return nil
	//}
	//args := make([]reflect.Value, 0, 3)
	//args = append(args, reflect.ValueOf(ep))
	//args = append(args, reflect.ValueOf(ctx))
	//args = append(args, reflect.ValueOf(msg))
	//backMsg := method.Call(args) // 分发
	//if backMsg != nil {
	//	return backMsg[0].Interface().(easygo.IMessage)
	//}
	return nil
}

//玩家上下线通知到后台
func (self *ServiceForHall) RpcNewLoginPush(common *base.Common, reqMsg *server_server.PlayerSI) easygo.IMessage {
	//推送玩家上下线通知消息到后台前端
	AllPush()
	return nil
}

//红包拆包处理
func (self *ServiceForHall) RpcDealReaPacket(common *base.Common, reqMsg *share_message.RedPacket) easygo.IMessage {
	//推送玩家上下线通知消息到后台前端
	logs.Info("RpcDealReaPacket", reqMsg)
	DealRedPacket(reqMsg)
	return reqMsg
}

//收到大厅发送的玩家消息
func (self *ServiceForHall) RpcHallSendIMmessage(common *base.Common, reqMsg *share_message.IMmessage) easygo.IMessage {
	SendIMmessage(reqMsg)
	return easygo.EmptyMsg
}

//收到大厅发送的物流信息
func (self *ServiceForHall) RpcSendShopOrderExpress(common *base.Common, reqMsg *server_server.ShopOrderExpressInfos) easygo.IMessage {
	logs.Info("RpcSendShopOrderExpress")
	epx := BrowerEpMp.LoadEndpoint(reqMsg.GetUserId())
	if epx != nil {
		var expressInfos []*brower_backstage.QueryShopOrderExpressBody = []*brower_backstage.QueryShopOrderExpressBody{}
		if reqMsg.GetExpressInfos() != nil {
			for _, value := range reqMsg.GetExpressInfos() {
				expressInfos = append(expressInfos, &brower_backstage.QueryShopOrderExpressBody{
					DateTime: easygo.NewString(value.GetDateTime()),
					Remark:   easygo.NewString(value.GetRemark()),
				})
			}
		}
		//发送到前端
		msg := &brower_backstage.QueryShopOrderExpressResponse{
			ExpressInfos: expressInfos,
			ExpressPhone: easygo.NewString(reqMsg.GetExpressPhone()),
			ExpressName:  easygo.NewString(reqMsg.GetExpressName()),
		}
		epx.CallRpcMethod("RpcSendShopOrderExpress", msg)
	} else {
		logs.Info("===============推送物流信息给后台客户端,ep断开")
	}
	return nil
}
