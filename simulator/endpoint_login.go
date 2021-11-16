//

package main

import (
	"game_server/easygo"
	"game_server/pb/client_login"
	"game_server/pb/client_server"
)

// 代表连接到大厅的一个连接
type ILoginEndpoint interface {
	easygo.IEndpointWithWebSocket
	client_login.IClient2Login
	client_server.IClient2server
	GetServerId() int32
	SetServerId(id int32)
}

type LoginEndpoint struct {
	easygo.EndpointWithWebSocket
	client_login.Client2Login
	client_server.Client2server

	Me        ILoginEndpoint
	ServerId  int32 // 服务器编号
	HasReport bool
}

func NewLoginEndpoint() *LoginEndpoint {
	services := map[string]interface{}{
		"ServiceForLogin": &ServiceForLogin{},
	}

	p := &LoginEndpoint{}
	p.Init(p, services)
	return p
}

func (self *LoginEndpoint) Init(me ILoginEndpoint, services map[string]interface{}, endpointId ...ENDPOINT_ID) { // 覆写，解决二义性
	self.Me = me
	downRpc := easygo.CombineRpcMap(client_login.DownRpc, client_server.DownRpc)
	upRpc := easygo.CombineRpcMap(client_login.UpRpc, client_server.UpRpc)
	self.EndpointWithWebSocket.Init(me, services, downRpc, upRpc, endpointId...)
	self.Client2Login.Init(self)
	self.Client2server.Init(self)
}
func (self *LoginEndpoint) SetServerId(id int32) {
	self.ServerId = id
}
func (self *LoginEndpoint) GetServerId() int32 {
	return self.ServerId
}
