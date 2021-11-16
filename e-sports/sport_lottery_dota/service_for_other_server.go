package sport_lottery_dota

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
)

//TODO 接收分发方法，函数方法名自己定，对端发送时传什么，这里就接收什么,需要返回值，则直接返回，但在接收的地方要强转
//比如对端服务器调用:SendMsgToServerNew(sid,"RpcTestXXXXMSG",msg),msg属于share_message.ServerInfo
//接收示例如下:
func (self *WebHttpForServer) RpcTestXXXXMSG(common *base.Common, reqMsg *share_message.ServerInfo) easygo.IMessage {
	return nil
}
