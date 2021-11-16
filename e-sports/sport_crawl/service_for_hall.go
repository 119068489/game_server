package sport_crawl

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"reflect"
)

//===================================================================

type ServiceForHall struct {
	Service reflect.Value
}

func (self *ServiceForHall) RpcHall2ShopMsg(common *base.Common, reqMsg *share_message.MsgToServer) easygo.IMessage {
	logs.Debug("RpcHall2ShopMsg", reqMsg)

	return nil
}
