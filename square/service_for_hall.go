package square

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"reflect"

	"github.com/astaxie/beego/logs"
)

type ServiceForHall struct {
	Service reflect.Value
}

func (self *ServiceForHall) RpcHallNotifyParamChange(comm *base.Common, reqMsg *server_server.SysteamModId) easygo.IMessage {
	logs.Info("后台修改配置:", reqMsg)
	PSysParameterMgr.UpLoad(reqMsg.GetId())
	return nil
}

// 处理大厅发给广场置顶消息
func (self *ServiceForHall) RpcHallNotifyTop(common *base.Common, reqMsg *share_message.BackstageNotifyTopReq) easygo.IMessage {
	logs.Info("============ 处理大厅发给广场置顶消息 RpcHallNotifyTop==============", reqMsg)
	// 因为是广播给全部广场,所以不是当前的服务器直接返回,其他服务器肯定有一个是可以收到本基处理的.
	saveSid := for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_SQUARE_SID)
	if saveSid != PServerInfo.GetSid() {
		logs.Error("不在同一个广场,直接返回nil,其他广场会处理,当前广场为: %d,存储的广场为: %d", PServerInfo.GetSid(), saveSid)
		return nil
	}
	// 后台置顶
	reqMsg.IsBsTop = easygo.NewBool(true)
	// 判断是否在当前服务器,如果是,直接添加定时任务,添加定时任务管理器里面.否则,通知管理定时任务的服务添加定时任务,存放数据库
	if err := for_game.ProcessTopTimer(reqMsg); err != nil {
		return err
	}
	return nil
}

// 处理大厅发给广场加载广场数据
func (self *ServiceForHall) RpcHallNotifyReloadSquare(common *base.Common, reqMsg *server_server.ReloadDynamicReq) easygo.IMessage {
	logs.Info("============ 处理大厅发给广场加载广场数据 RpcHallNotifyReloadSquare==============", reqMsg)
	for_game.ReloadMyDynamicInfo(reqMsg.GetPlayerId())
	return easygo.EmptyMsg
}
