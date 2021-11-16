package login

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/server_server"
	"github.com/astaxie/beego/logs"
)

//TODO 接收分发方法
func (self *WebHttpForServer) RpcHallNotifyParamChange(comm *base.Common, reqMsg *server_server.SysteamModId) easygo.IMessage {
	logs.Info("后台修改配置:", reqMsg)
	PSysParameterMgr.UpLoad(reqMsg.GetId())
	return nil
}
