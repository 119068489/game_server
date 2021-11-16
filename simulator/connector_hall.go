package main

import (
	"game_server/easygo"
	"github.com/astaxie/beego/logs"
)

// 商店连接到通知的连接器
type HallConnector struct {
	easygo.WebConnector
	IsReport bool
}

func NewHallConnector(address string) *HallConnector {
	p := &HallConnector{
		IsReport: false,
	}
	p.Init(p, address, "大厅服务器")
	return p
}

func (self *HallConnector) FetchEndpoint() IHallEndpoint { // overwrite 强制断言成合适的类型
	endpoint := self.WebConnector.FetchEndpoint()
	ep, ok := endpoint.(IHallEndpoint)
	if !ok { // 如果是 nil 接口强制断言会 panic 的，所以得试探
		return nil
	}
	return ep
}

func (self *HallConnector) ConnectUntilOk() IHallEndpoint { // overwrite 强制断言成合适的类型
	endpoint := self.WebConnector.ConnectUntilOk()
	ep, ok := endpoint.(IHallEndpoint)
	if !ok { // 如果是 nil 接口强制断言会 panic 的，所以得试探
		return nil
	}
	return ep
}
func (self *HallConnector) OnConnectOk(ep easygo.IEndpointWithWebSocket) { // override
	//连接上了，向登录服务器报道
	logs.Info("连接Hall服务器  succse")
}

func (self *HallConnector) CreateEndpoint() easygo.IEndpointWithWebSocket { // override
	return NewHallEndpoint()
}

func (self *HallConnector) OnReConnectOk(ep easygo.IEndpointWithWebSocket) { // override
	//重连逻辑处理
}
