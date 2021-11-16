// 连接器。用于发起 TCP 连接到 Server 端
// todo 此类的线程安全需要再次检查一次
package easygo

import (
	"fmt"
	"github.com/gorilla/websocket"
	"math"
	"net/url"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

//var MIN_DELAY, MAX_DELAY time.Duration = 50 * time.Millisecond, 3 * time.Second

type IWebConnector interface {
	CreateEndpoint() IEndpointWithWebSocket

	OnConnectOk(ep IEndpointWithWebSocket)
	OnDisconnected(ep IEndpointWithWebSocket)
	OnReConnectOk(ep IEndpointWithWebSocket)

	DoReConnect()

	ConnectLoop()
	String() string
}

// 线程安全的连接器
type WebConnector struct {
	Me IWebConnector

	Mutex     RLock
	ConnectOK *sync.Cond

	Endpoint       IEndpointWithWebSocket
	Address        string
	TargetName     string
	AutoReConnect  bool
	LoopConnecting bool
	IsStop         bool //是否停止
}

/* 抽象类，不能实例化
func NewConnector(address string, targetName string, autoReConnect ...bool) *Connector {
	p := &Connector{}
	p.Init(p, address, targetName, autoReConnect...)
	return p
}*/

func (self *WebConnector) Init(me IWebConnector, address string, targetName string, autoReConnect ...bool) {
	self.Me = me
	self.Address = address
	self.TargetName = targetName
	self.AutoReConnect = append(autoReConnect, true)[0]
	self.ConnectOK = sync.NewCond(&self.Mutex)
	self.LoopConnecting = false
	self.IsStop = false
}

func (self *WebConnector) ConnectUntilOk() IEndpointWithWebSocket { // 不断尝试连接，直到成功
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if self.Endpoint == nil && !self.LoopConnecting {
		Spawn(self.Me.ConnectLoop)
	}

	for self.Endpoint == nil {
		self.ConnectOK.Wait()
	}
	return self.Endpoint
}

func (self *WebConnector) ConnectLoop() {
	// 此函数不要长时间上锁，否则导致 GetEndpoint() 不能快速返回
	self.Mutex.Lock()
	if self.LoopConnecting {
		self.Mutex.Unlock()
		return
	}
	self.LoopConnecting = true
	self.Mutex.Unlock()

	defer func() {
		self.Mutex.Lock()
		defer self.Mutex.Unlock()
		self.LoopConnecting = false
	}()
	delay := MIN_DELAY
	for {
		if self.IsStop {
			self.ConnectOK.Broadcast()
			return
		}
		ep := self.FetchEndpoint() // 获取时会自动连
		if ep != nil {             // 连接成功
			self.ConnectOK.Broadcast()
			return
		}
		a, b := float64(MAX_DELAY), float64(delay*2)
		delay = time.Duration(math.Min(a, b))
		time.Sleep(delay)
	}
}
func (self *WebConnector) CreateEndpoint() IEndpointWithWebSocket {
	panic("抽象方法，请在子类实现")
}

// 尝试一次连接
func (self *WebConnector) ConnectOnce() IEndpointWithWebSocket {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	if self.Endpoint != nil { // 已建立了连接了 // && !reflect.ValueOf(self.Endpoint).IsNil()
		return self.Endpoint
	}
	logs.Info("尝试连接到(%s)%v", self.TargetName, self.Address)
	u := url.URL{Scheme: "ws", Host: self.Address, Path: ""}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logs.Info("连接(%s)失败.%v", self.TargetName, err)
		return nil
	}

	//addr, err := net.ResolveTCPAddr("tcp4", self.Address)
	//PanicError(err)
	//conn, err := net.DialTCP("tcp4", nil, addr)

	logs.Info("连接(%s)成功,%v", self.TargetName, self.Address)

	ep := self.Me.CreateEndpoint()
	event := ep.GetDisconnectedEvent()
	event.AddHandler(self.Me.OnDisconnected)
	address := conn.RemoteAddr()
	ep.SetAddr(address)
	ep.SetConnection(conn)
	ep.Start()
	self.Me.OnConnectOk(ep)
	self.Endpoint = ep
	return ep
}

func (self *WebConnector) OnConnectOk(ep IEndpointWithWebSocket) {
}

// 断线处理函数
func (self *WebConnector) OnDisconnected(ep IEndpointWithWebSocket) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	self.Endpoint = nil // 断线置 nil
	if self.AutoReConnect {
		Spawn(self.Me.DoReConnect) // 在新协程中不断重连
	}
}

func (self *WebConnector) OnReConnectOk(ep IEndpointWithWebSocket) {

}

func (self *WebConnector) DoReConnect() {
	ep := self.ConnectUntilOk()
	self.Me.OnReConnectOk(ep)
}

// 如果连接断开，会自动连接一次
func (self *WebConnector) FetchEndpoint() IEndpointWithWebSocket {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if self.GetIsStop() {
		return self.Endpoint
	}
	if self.Endpoint == nil {
		self.ConnectOnce() // 连接1次,有可能是不成功的,特别是服务器宕掉后
	}
	return self.Endpoint
}

func (self *WebConnector) GetEndpoint() IEndpointWithWebSocket {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.Endpoint
}

func (self *WebConnector) String() string {
	return fmt.Sprintf(`Connector=%v`, self.Endpoint)
}
func (self *WebConnector) SetIsStop(b bool) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.IsStop = b
}
func (self *WebConnector) GetIsStop() bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.IsStop
}
