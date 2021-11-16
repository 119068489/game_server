// tcpsocket 服务器
package easygo

import (
	"fmt"
	"net"

	"github.com/astaxie/beego/logs"
)

type ITcpSocketServer interface {
	ISocketServer
}

type TcpSocketServer struct {
	SocketServer
	Me   ITcpSocketServer
	Name string
}

// 地址长成这样 0.0.0.0:1211    localhost:1211
func NewTcpSocketServer(address string, name string) *TcpSocketServer {
	p := &TcpSocketServer{}
	p.Init(p, address, name)
	return p
}

func (self *TcpSocketServer) Init(me ITcpSocketServer, address string, name string) {
	self.Me = me
	self.Name = name
	self.SocketServer.Init(me, address)
}
func (self *TcpSocketServer) Serve() { // impletement
	addr, err := net.ResolveTCPAddr("tcp4", self.Address)
	PanicError(err)

	listener, err := net.ListenTCP("tcp4", addr) //监听
	PanicError(err)
	defer listener.Close()

	s := fmt.Sprintf("(%s) 正在监听: %v", self.Name, addr)
	logs.Info(s)
	for {
		conn, err := listener.AcceptTCP()
		PanicError(err)
		s := fmt.Sprintf("在 (%s) 的 %v 接收到连接: %s", self.Name, self.Address, conn.RemoteAddr())
		logs.Debug(s)
		Spawn(self.Me.HandleConnection, conn) // 起一个goroutine处理
	}
}
