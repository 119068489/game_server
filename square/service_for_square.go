// 大厅发给子游戏服的消息处理

package square

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"reflect"
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

// 玩家广场上线通知
func (self *ServiceForSquare) RpcPlayerOnLine(common *base.Common, reqMsg *share_message.PlayerState) easygo.IMessage {
	PlayerOnlineMgr.PlayerOnline(reqMsg.GetPlayerId(), reqMsg.GetServerId())
	return nil
}

// 处理其他广场发送的定时任务逻辑
func (self *ServiceForSquare) RpcSquareNotifyTop(common *base.Common, reqMsg *share_message.BackstageNotifyTopReq) easygo.IMessage {
	if err := for_game.ProcessTopTimer(reqMsg); err != nil {
		logs.Error("处理其他广场发送的定时任务逻辑 RpcSquareNotifyTop, err: ", err)
		return err
	}
	return nil
}
