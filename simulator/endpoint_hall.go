package main

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
)

type IHallEndpoint interface {
	easygo.IEndpointWithWebSocket
	client_hall.IAccountClient2Hall
	client_server.IClient2server
	SetServerId(id SERVER_ID)
	GetServerId() SERVER_ID
}

type HallEndpoint struct {
	easygo.EndpointWithWebSocket
	client_hall.AccountClient2Hall
	client_server.Client2server
	ServerId SERVER_ID
	Me       IHallEndpoint
}

func NewHallEndpoint() *HallEndpoint { // services map[string]interface{},
	services := map[string]interface{}{
		"ServiceForHall": &ServiceForHall{},
	}

	p := &HallEndpoint{}
	p.Init(p, services)
	return p
}

func (self *HallEndpoint) Init(me IHallEndpoint, services map[string]interface{}, endpointId ...easygo.ENDPOINT_ID) {
	downRpc := easygo.CombineRpcMap(client_hall.DownRpc, client_server.DownRpc)
	upRpc := easygo.CombineRpcMap(client_hall.UpRpc, client_server.UpRpc)
	self.EndpointWithWebSocket.Init(me, services, downRpc, upRpc, endpointId...)
	self.AccountClient2Hall.Init(me)
	self.Client2server.Init(me)

}
func (self *HallEndpoint) SetServerId(id SERVER_ID) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.ServerId = id

}
func (self *HallEndpoint) GetServerId() SERVER_ID {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.ServerId
}
func (self *HallEndpoint) SendToOtherServer(data []byte, request *base.Request) {

}
