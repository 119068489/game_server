package wish

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
)

//TODO 接收分发方法，函数方法名自己定，对端发送时传什么，这里就接收什么,需要返回值，则直接返回，但在接收的地方要强转
//比如对端服务器调用:SendMsgToServerNew(sid,"RpcTestXXXXMSG",msg),msg属于share_message.ServerInfo
//接收示例如下:
func (self *WebHttpForServer) RpcTestXXXXMSG(common *base.Common, reqMsg *share_message.ServerInfo) easygo.IMessage {
	return nil
}

//回收回调  支付订单完成后修改物品状态
func (self *WebHttpForServer) RpcRecycleCallBack(common *base.Common, reqMsg *h5_wish.OrderMsgResp) easygo.IMessage {
	logs.Info("===回收回调 RpcRecycleCallBack reqMsg: %v", reqMsg)
	for_game.UpdateWishOrderByPayOrder(reqMsg.GetOrderId())
	return nil
}
