// 管理后台给浏览器连接的 server

package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"

	"github.com/astaxie/beego/logs"
)

type ServerForBrower struct {
	easygo.WebSocketServer
}

func NewServerForBrower(address string) *ServerForBrower {
	p := &ServerForBrower{}
	p.Init(p, address, "管理后台 for Web前端")
	return p
}

func (self *ServerForBrower) CreateEndpoint(endpointId ENDPOINT_ID) easygo.IEndpointWithSocket { // override
	ep := NewBrowerEndpoint(endpointId)
	BrowerEpMgr.Store(endpointId, ep)

	ep.Disconnected.AddHandler(self.OnDisConnect)
	return ep
}

func (self *ServerForBrower) OnDisConnect(ep IBrowerEndpoint) {
	epId := ep.GetEndpointId()
	BrowerEpMgr.Delete(epId)
	user := ep.GetUser()
	var userId USER_ID
	if user != nil {
		userId = user.GetId()
		BrowerEpMp.Delete(userId)
		for_game.DelRedisWaiter(userId)
		for_game.DelRedisAdmin(userId)
	}
	if user != nil {
		user.IsOnlie = easygo.NewBool(false)
		//更新登录数据
		EditManage(user.GetSite(), user, "logout")
	}
	s := fmt.Sprintf("浏览器断开与服务器的连接,ep=%d,user=%d\n", epId, userId)
	logs.Info(s)
}
