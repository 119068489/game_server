package sport_apply

import (
	dal "game_server/e-sports/sport_common_dal"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"reflect"
)

//===================================================================

type ServiceForBackStage struct {
	Service reflect.Value
}

// 后台注单操作后通知电竞消息
func (self *ServiceForBackStage) RpcBetSlipOperateNotify(common *base.Common, reqMsg *share_message.TableESPortsGameOrderSysMsg) easygo.IMessage {
	logs.Info("=======RpcBetSlipOperateNotify========= ", reqMsg)

	notifyRst := dal.PushGameOrderSysMsg(PServerInfoMgr, reqMsg)

	if notifyRst.GetCode() != for_game.C_OPT_SUCCESS {
		return easygo.NewFailMsg(notifyRst.GetMsg())
	}

	return nil
}
