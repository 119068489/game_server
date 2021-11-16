// 大厅发给子游戏服的消息处理

package hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_server"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"reflect"

	"github.com/astaxie/beego/logs"
)

//===================================================================

type ServiceForHall struct {
	Service reflect.Value
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcServerReport(common *base.Common, reqMsg *share_message.ServerInfo) easygo.IMessage {
	logs.Info("RpcServerReport ", reqMsg)
	//ep.SetServerId(reqMsg.GetSid())
	//HallEpMgr.Store(reqMsg.GetSid(), ep)
	PServerInfoMgr.AddServerInfo(reqMsg)
	msg := PlayerOnlineMgr.GetAllOnlinePlayers()
	return msg
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

func (self *ServiceForHall) RpcAgreeAddFriend(common *base.Common, reqMsg *server_server.AddFriend) easygo.IMessage {
	pid := reqMsg.GetPlayerId()
	friendId := reqMsg.GetFriendId()
	//player := PlayerMgr.LoadPlayer(pid)
	//if player != nil {
	ep1 := ClientEpMp.LoadEndpoint(pid)
	if ep1 != nil {
		msg := GetFriendInfo(friendId)
		if reqMsg.GetOpenWindows() != 0 {
			msg.OpenWindows = easygo.NewInt32(reqMsg.GetOpenWindows()) //不打开聊天窗口
		}
		ep1.RpcAgreeFriendResponse(msg)
	}
	//}

	return nil
}

func (self *ServiceForHall) RpcAddFriendRequest(common *base.Common, reqMsg *server_server.AddFriend) easygo.IMessage {
	pid := reqMsg.GetPlayerId()
	ep1 := ClientEpMp.LoadEndpoint(pid)
	if ep1 != nil {
		fq := for_game.GetFriendBase(pid)
		msg := fq.GetNewVersionAllFriendRequestForOne() //发送所有好友申请信息
		ep1.RpcNoticeAddFriend(msg)
	}
	return nil
}

//消息广播转发给指定玩家
//func (self *ServiceForHall) RpcBroadCastMsgToClient(common *base.Common, reqMsg *hall_hall.MsgToClient) easygo.IMessage {
//	logs.Info("RpcBroadCastMsgToClient:", reqMsg)
//	for _, pid := range reqMsg.GetPlayerIds() {
//		//player := PlayerMgr.LoadPlayer(pid)
//		//if player != nil {
//		//直接推送给客户端
//		if reqMsg.GetIsSend() {
//			ep := ClientEpMp.LoadEndpoint(pid)
//			if ep != nil {
//				msg := easygo.NewMessage(reqMsg.GetMsgName())
//				err := msg.Unmarshal(reqMsg.GetMsg())
//				easygo.PanicError(err)
//				_, err1 := ep.CallRpcMethod(reqMsg.GetRpcName(), msg)
//				easygo.PanicError(err1)
//			} else {
//				//指定用户处理逻辑
//				methodName := reqMsg.GetRpcName()
//				self.Service = reflect.ValueOf(self)
//				method := self.Service.MethodByName(methodName)
//				msg := easygo.NewMessage(reqMsg.GetMsgName())
//				if !method.IsValid() || method.Kind() != reflect.Func {
//					logs.Info("无效的rpc请求，找不到methodName:", methodName)
//					return nil
//				}
//				args := make([]reflect.Value, 0, 3)
//				//args = append(args, reflect.ValueOf(ep))
//				//args = append(args, reflect.ValueOf(ctx))
//				args = append(args, reflect.ValueOf(msg))
//				backMsg := method.Call(args) // 分发到指定的rpc
//				if backMsg != nil {
//					return backMsg[0].Interface().(easygo.IMessage)
//				}
//			}
//
//		}
//	}
//	return nil
//}

//消息广播转指定玩家处理逻辑
func (self *ServiceForHall) RpcBroadCastMsgToHall(common *base.Common, reqMsg *server_server.MsgToHall) easygo.IMessage {
	logs.Info("RpcBroadCastMsgToHall:", reqMsg)
	player := PlayerMgr.LoadPlayer(reqMsg.GetPlayerId())
	if player != nil {
		epx := ClientEpMp.LoadEndpoint(player.GetPlayerId())
		if epx != nil {
			msg := easygo.NewMessage(reqMsg.GetMsgName())
			err := msg.Unmarshal(reqMsg.GetMsg())
			easygo.PanicError(err)
			methodName := reqMsg.GetRpcName()
			self.Service = reflect.ValueOf(self)
			method := self.Service.MethodByName(methodName)

			if !method.IsValid() || method.Kind() != reflect.Func {
				logs.Info("无效的rpc请求，找不到methodName:", methodName)
				return nil
			}
			args := make([]reflect.Value, 0, 2)
			args = append(args, reflect.ValueOf(common))
			args = append(args, reflect.ValueOf(msg))
			backMsg := method.Call(args) // 分发到指定的rpc
			if backMsg != nil {
				return backMsg[0].Interface().(easygo.IMessage)
			}
		}
	}
	return nil
}

//自定义RPC
func (self *ServiceForHall) RpcRechargeToHall(common *base.Common, reqMsg *server_server.Recharge) easygo.IMessage {
	logs.Info("其他大厅转发过来的 RpcRechargeToHall", reqMsg)
	HandleAfterRecharge(reqMsg.GetOrderId())
	return nil
}

func (self *ServiceForHall) RpcAssistantNotify(common *base.Common, reqMsg *client_server.AssistantMsg) easygo.IMessage {
	ep2 := ClientEpMp.LoadEndpoint(reqMsg.GetPlayerId())
	if ep2 != nil {
		ep2.RpcAssistantNotify(reqMsg)
	}
	return nil
}

//心跳
func (self *ServiceForHall) RpcHeartBeat(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	//logs.Info("RpcHeartBeat ")
	return nil
}

//其他大厅顶号登录通知
func (self *ServiceForHall) RpcOtherHallLogin(common *base.Common, reqMsg *server_server.PlayerIdInfo) easygo.IMessage {
	pid := reqMsg.GetPlayerId()
	reqMsg.Success = easygo.NewBool(false)
	if PlayerOnlineMgr.CheckPlayerIsOnLine(pid) {
		ep := ClientEpMp.LoadEndpoint(pid)
		if ep != nil {
			ep.RpcReLogin(nil) //通知前端被顶号了，返回登录界面
			ep.Shutdown()
			reqMsg.Success = easygo.NewBool(true)
		}
	}
	return reqMsg
}

//其他地方红包退款通知

//func (self *ServiceForHall) RpcHallToClient(common *base.Common, reqMsg *share_message.MsgData) easygo.IMessage {
//	for _, pid := range reqMsg.GetPlayerids() {
//		player := PlayerMgr.LoadPlayer(pid)
//		if player != nil {
//			ep := ClientEpMp.LoadEndpoint(player.GetPlayerId())
//			if ep != nil {
//				if reqMsg.GetMsgName() != "" {
//					msg := easygo.NewMessage(reqMsg.GetMsgName())
//					err := msg.Unmarshal(reqMsg.GetMsg())
//					easygo.PanicError(err)
//					_, err1 := ep.CallRpcMethod(reqMsg.GetRpcName(), msg)
//					easygo.PanicError(err1)
//				} else {
//					_, err1 := ep.CallRpcMethod(reqMsg.GetRpcName(), nil)
//					easygo.PanicError(err1)
//				}
//
//			}
//		}
//	}
//	return nil
//}

// 其他大厅发送聊天记录
func (self *ServiceForHall) RpcChatToOther(common *base.Common, reqMsg *server_server.ChatToOtherReq) easygo.IMessage {
	logs.Info("==========其他大厅发送聊天记录 RpcChatToOther =========", reqMsg)
	pid := reqMsg.GetPlayerId()
	player := GetPlayerObj(pid)
	if player == nil {
		logs.Error("对象为空")
		return nil
	}
	//ep1 := ClientEpMp.LoadEndpoint(pid)
	//if ep1 != nil {
	cls1 := new(ServiceForGameClient)
	cls1.RpcChatNew(nil, player, reqMsg.GetChat())
	//}
	return nil
}
