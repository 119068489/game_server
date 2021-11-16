package easygo

import (
	"game_server/easygo/base"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

type IEndpointWithWebSocket interface {
	IEndpointWithSocket
}

type EndpointWithWebSocket struct {
	EndpointWithSocket
	Me IEndpointWithWebSocket

	Conn *websocket.Conn
}

type c6 = EndpointWithWebSocket

func NewEndpointWithWebSocket(services map[string]interface{}, serviceMap map[string]Pair, stubMap map[string]Pair, endpointId ...ENDPOINT_ID) *EndpointWithWebSocket {
	p := &EndpointWithWebSocket{}
	p.Init(p, services, serviceMap, stubMap, endpointId...)
	return p
}

func (self *c6) Init(me IEndpointWithWebSocket, services map[string]interface{}, serviceMap map[string]Pair, stubMap map[string]Pair, endpointId ...ENDPOINT_ID) {
	self.Me = me
	self.EndpointWithSocket.Init(me, services, serviceMap, stubMap, endpointId...)
	self.Conn = nil

}

func (self *c6) ReadMsg() ([][]byte, error) { // override
	t, p, e := self.Conn.ReadMessage()
	_ = t
	if e != nil {
		return nil, e
	}
	return [][]byte{p}, nil
}

func (self *c6) ShutdownWrite() error { // override
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "") // to check
	err := self.Conn.WriteMessage(websocket.CloseMessage, msg)
	return err
}

func (self *c6) CloseConnection() { // override
	self.EndpointWithSocket.CloseConnection()
	err := self.Conn.Close()
	_ = err
	//PanicError(err)
}

func (self *c6) SendMsg(packet []byte) error { // override
	err := self.Conn.WriteMessage(websocket.BinaryMessage, packet)
	return err
}

func (self *c6) SetConnection(conn EasygoConn) IEndpointWithSocket { // override
	self.Conn = conn.(*websocket.Conn)
	return self.Me
}

func (self *c6) GetConnection() EasygoConn {
	return self.Conn
}

//消息转发给指定服务器
func (self *c6) RoutingToOtherServer(data []byte, request *base.Request) {
	logs.Info("c6========")
}
