package statistics

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
)

//TODO 消息接收分发
func (self *WebHttpForServer) RpcReportServer(common *base.Common, reqMsg *share_message.ServerInfo) easygo.IMessage {
	logs.Info("serverInfo 报道:", reqMsg)
	return nil
}
