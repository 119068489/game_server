// 大厅发给子游戏服的消息处理

package hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"reflect"

	"github.com/astaxie/beego/logs"
)

//===================================================================

type ServiceForSquare struct {
	Service reflect.Value
}

//分发消息
func (self *ServiceForSquare) RpcServerReport(common *base.Common, reqMsg *share_message.ServerInfo) easygo.IMessage {
	logs.Info("RpcServerReport:", reqMsg)
	//ep.SetServerId(reqMsg.GetSid())
	//SquareEpMgr.Store(reqMsg.GetSid(), ep)
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

//心跳
func (self *ServiceForSquare) RpcHeartBeat(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	//logs.Info("RpcHeartBeat ")
	return nil
}

// 接受广场下发ep给前端  RpcSquare2HallClient
func (self *ServiceForSquare) RpcSquare2HallClient(common *base.Common, reqMsg *share_message.MsgToClient) easygo.IMessage {
	logs.Info("=====广场接受大厅的下发 RpcSquare2HallClient====,", reqMsg)
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

// 推送话题动态置顶/取消置顶信息
func (self *WebHttpForServer) RpcTopicDynamicTopStatus(common *base.Common, reqMsg *server_server.TopicDynamicTopLittleHelper) easygo.IMessage {
	logs.Info("==========大厅收到了话题小助手通知请求RpcTopicDynamicTopStatus: %v", reqMsg)
	pid := reqMsg.GetPlayerId()
	if pid == 0 {
		logs.Error("推送失败,请求参数为reqMsg: %v ", reqMsg)
		return easygo.NewFailMsg("推送失败")
	}
	// 小助手通知
	NoticeAssistant(pid, 1, reqMsg.GetTitle(), reqMsg.GetContent())
	return nil
}
