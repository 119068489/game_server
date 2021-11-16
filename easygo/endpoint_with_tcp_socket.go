package easygo

import (
	"errors"
	"game_server/easygo/base"
	"github.com/astaxie/beego/logs"
	"io"
	"net"
)

type IEndpointWithTcpSocket interface {
	IEndpointWithSocket

	CreateBuffer() IBuffer
	CreateDecoder() IDecoder
	CreateEncoder() IEncoder
}

type EndpointWithTcpSocket struct {
	EndpointWithSocket
	Me IEndpointWithTcpSocket

	Decoder IDecoder
	Encoder IEncoder

	Conn   net.Conn
	Buffer IBuffer
}

type c5 = EndpointWithTcpSocket

func NewEndpointWithTcpSocket(services map[string]interface{}, serviceMap map[string]Pair, stubMap map[string]Pair, endpointId ...ENDPOINT_ID) *EndpointWithTcpSocket {
	p := &EndpointWithTcpSocket{}
	p.Init(p, services, serviceMap, stubMap, endpointId...)
	return p
}

func (self *c5) Init(me IEndpointWithTcpSocket, services map[string]interface{}, serviceMap map[string]Pair, stubMap map[string]Pair, endpointId ...ENDPOINT_ID) {
	self.Me = me
	self.EndpointWithSocket.Init(me, services, serviceMap, stubMap, endpointId...)
	self.Conn = nil
	self.Decoder, self.Encoder = self.Me.CreateDecoder(), self.Me.CreateEncoder()
	self.Buffer = self.Me.CreateBuffer()
}

func (self *c5) ReadMsg() ([][]byte, error) { // override
	bytes := self.Buffer.PeekWrite()
	size, err := self.Conn.Read(bytes)
	// 对端强行关闭 err 为 An existing connection was forcibly closed by the remote host
	// 本端执行 close 后 err 为 use of closed network connection

	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		if _, ok := err.(*net.OpError); !ok { // 按理 Read 只能返回 *net.OpError 错误，不太确定，我还是断言一下
			panic(err)
		}
		return nil, err
	}
	self.Buffer.HasWritten(size)
	msgs, reason := self.Decoder.Decode(self.Buffer)
	if reason == "" {
		return msgs, nil
	} else {
		return msgs, errors.New(reason)
	}
}

func (self *c5) ShutdownWrite() error { // override
	conn := self.Conn.(*net.TCPConn)
	err := conn.CloseWrite()
	return err
}

func (self *c5) CloseConnection() { // override
	self.EndpointWithSocket.CloseConnection()
	err := self.Conn.Close()
	_ = err
	//PanicError(err)
}

func (self *c5) SendMsg(packet []byte) error { // override
	b := self.Encoder.Encode(packet)
	return self.SendAll(self.Conn, b)
}

// todo 搬到公共的地方，变成模块函数
func (self *c5) SendAll(conn net.Conn, packet []byte) error {
	for i := 0; i < len(packet); {
		// 本端强行 close 后： use of closed network connection
		// 对方close 后 An established connection was aborted by the software in your host machine
		n, err := conn.Write(packet[i:])
		if err != nil {
			if _, ok := err.(*net.OpError); !ok { // 按理 Write 只能返回 *net.OpError 错误，不太确定，我还是断言一下
				panic(err)
			}
			return err
		}
		i += n
	}
	return nil
}

func (self *c5) CreateBuffer() IBuffer {
	return NewBuffer(1024)
}

func (self *c5) SetConnection(conn EasygoConn) IEndpointWithSocket { // implement
	self.Conn = conn.(net.Conn)
	// sock.setsockopt(socket.SOL_TCP, socket.TCP_NODELAY, 1) // todo
	return self.Me
}

func (self *c5) GetConnection() EasygoConn {
	return self.Conn
}

func (self *c5) CreateDecoder() IDecoder {
	return NewDecoder(1 * 1024 * 1024) // 最大能接收的逻辑包大小,1兆
}

func (self *c5) CreateEncoder() IEncoder {
	return NewEncoder()
}

//消息转发给指定服务器
func (self *c5) RoutingToOtherServer(data []byte, request *base.Request) {
	logs.Info("c5========")
}
