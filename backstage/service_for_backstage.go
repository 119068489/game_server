// 管理后台为[浏览器]提供的服务

package backstage

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
)

//分发消息
func (self *WebHttpForServer) RpcBServerReport(common *base.Common, reqMsg *share_message.ServerInfo) easygo.IMessage {
	logs.Info("RpcBServerReport:", reqMsg)
	//ep.SetServerId(reqMsg.GetSid())
	//BackstageEpMgr.Store(reqMsg.GetSid(), ep)
	PServerInfoMgr.AddServerInfo(reqMsg)
	return nil
}

//心跳
func (self *WebHttpForServer) RpcBHeartBeat(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	//logs.Info("=========== 其他后台发过来的心跳RpcHeartBeat===========", ep.GetServerId())
	return nil
}

//后台用户上线通知
func (self *WebHttpForServer) RpcUserOnLine(common *base.Common, reqMsg *share_message.PlayerState) easygo.IMessage {
	UserOnlineMgr.PlayerOnline(reqMsg.GetPlayerId(), reqMsg.GetServerId())
	return nil
}
