//

package main

import (
	"game_server/easygo"
	"github.com/astaxie/beego/logs"
)

// 子游戏连接到大厅的连接器
type LoginConnector struct {
	easygo.WebConnector
	IsReport bool
}

func NewLoginConnector(address string) *LoginConnector {
	p := &LoginConnector{
		IsReport: false,
	}
	p.Init(p, address, "Login服务器")
	return p
}

func (self *LoginConnector) FetchEndpoint() ILoginEndpoint { // overwrite 强制断言成合适的类型
	endpoint := self.WebConnector.FetchEndpoint()
	ep, ok := endpoint.(ILoginEndpoint)
	if !ok { // 如果是 nil 接口强制断言会 panic 的，所以得试探
		return nil
	}
	return ep
}

func (self *LoginConnector) ConnectUntilOk() ILoginEndpoint { // overwrite 强制断言成合适的类型
	endpoint := self.WebConnector.ConnectUntilOk()
	ep, ok := endpoint.(ILoginEndpoint)
	if !ok { // 如果是 nil 接口强制断言会 panic 的，所以得试探
		return nil
	}
	return ep
}
func (self *LoginConnector) OnConnectOk(ep easygo.IEndpointWithWebSocket) { // override
	//连接上了，向登录服务器报道
	_, ok := ep.(ILoginEndpoint)
	if !ok { // 如果是 nil 接口强制断言会 panic 的，所以得试探
		return
	}
	logs.Info("连接Login服务器  succse")
}

func (self *LoginConnector) CreateEndpoint() easygo.IEndpointWithWebSocket { // override
	return NewLoginEndpoint()
}

func (self *LoginConnector) OnReConnectOk(ep easygo.IEndpointWithWebSocket) { // override
	//重连逻辑处理
}
