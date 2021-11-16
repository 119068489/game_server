// socket 服务器。目前有 2 个子类，分别是 TcpSocketServer 和 WebSocketServer
package easygo

import (
	"net"
	"strconv"
	"sync"
	"sync/atomic"
)

type ISocketServer interface {
	GenEndpointId() ENDPOINT_ID
	HandleConnection(conn EasygoConn, uid ...string)
	CreateEndpoint(endpointId ENDPOINT_ID) IEndpointWithSocket
	Serve()
}

type SocketServer struct {
	Me             ISocketServer
	Address        string
	LastEndpointId int32
	ConnMap        sync.Map
}

/* 抽象类，不能实例化
func NewSocketServer(address string) *SocketServer {
}*/

func (self *SocketServer) Init(me ISocketServer, address string) {
	self.Me = me
	self.Address = address
}

func (self *SocketServer) CreateEndpoint(endpointId ENDPOINT_ID) IEndpointWithSocket {
	panic("抽象方法，请在子类实现")
}
func (self *SocketServer) Serve() {
	panic("抽象方法，请在子类实现")
}

func (self *SocketServer) GenEndpointId() ENDPOINT_ID { // 线程安全
	v := atomic.AddInt32(&self.LastEndpointId, 1) // 溢出后自动回转
	return v
}

func (self *SocketServer) HandleConnection(conn EasygoConn, conId ...string) {
	uid := append(conId, "")[0]
	id := self.Me.GenEndpointId()
	ep := self.Me.CreateEndpoint(id)
	ep.SetConnection(conn)
	address := conn.RemoteAddr()
	ep.SetAddr(address)
	if uid != "" {
		ep.SetUid(uid)
		self.ConnMap.Store(uid, ep)
	}
	ep.Start()
}

// 好像没有什么用了，先留着吧
func (self *SocketServer) GetPort() int32 {
	_, sport, err := net.SplitHostPort(self.Address) //获取客户端IP
	PanicError(err)
	p, err := strconv.Atoi(sport)
	PanicError(err)
	return int32(p)
}
