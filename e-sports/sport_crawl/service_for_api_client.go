// 大厅服务器为[游戏客户端]提供的服务

package sport_crawl

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/client_hall"
	"reflect"

	"github.com/astaxie/beego/logs"
)

type ServiceForClient struct {
	Service reflect.Value
}
type sfc = ServiceForClient

func (self *sfc) RpcESportEnter(common *base.Common, reqMsg *client_hall.ESportCommonResult) easygo.IMessage {
	logs.Info("===api RpcESportEnter===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	rd := client_hall.ESportCommonResult{
		Code: easygo.NewInt32(1),
		Msg:  easygo.NewString("跑错地方了吧，小老弟"),
	}
	return &rd
}

//RpcESportEnter(client_hall.EnterMsg)returns(client_hall.CommonResult);
