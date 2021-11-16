package easygo

import (
	"game_server/easygo/base"
	"io"
	"time"

	"github.com/astaxie/beego/logs"
)

var _STOP_SENTRY interface{} = nil

type IEndpointWithSocket interface {
	IEndpointBase

	InterceptAndDeal(packet []byte) (bool, []byte)

	Join()

	CreateSendQueue() IQueue

	SendProc()
	RecvProc()

	ReadMsg() ([][]byte, error)
	SendMsg(packet []byte) error
	ShutdownWrite() error
	CloseConnection()
	SetConnection(conn EasygoConn) IEndpointWithSocket
	GetConnection() EasygoConn
}

type EndpointWithSocket struct {
	EndpointBase
	Me IEndpointWithSocket

	SendQueue IQueue

	RecvJob       IGoroutine
	SendJob       IGoroutine
	HasStopSentry bool
	Closed        bool
}

type c2 = EndpointWithSocket

func NewEndpointWithSocket(services map[string]interface{}, serviceMap map[string]Pair, stubMap map[string]Pair, endpointId ...ENDPOINT_ID) *EndpointWithSocket {
	p := &EndpointWithSocket{}
	p.Init(p, services, serviceMap, stubMap, endpointId...)
	return p
}

// const SEND_QUEUE_SIZE = 0 // 消息包的个数,若是无上限,则赋值为 0

func (self *c2) Init(me IEndpointWithSocket, services map[string]interface{}, serviceMap map[string]Pair, stubMap map[string]Pair, endpointId ...ENDPOINT_ID) { //  services map[string]interface{},
	self.Me = me
	self.EndpointBase.Init(me, services, serviceMap, stubMap, endpointId...)

	self.RecvJob, self.SendJob = nil, nil
	self.SendQueue = self.Me.CreateSendQueue()
	self.HasStopSentry = false
}

// protected
func (self *c2) ReadMsg() ([][]byte, error) {
	panic("抽象方法，子类必须实现")
}

func (self *c2) SendMsg(packet []byte) error {
	panic("抽象方法，子类必须实现")
}

// protected
func (self *c2) ShutdownWrite() error {
	panic("抽象方法，子类必须实现")
}

// protected
func (self *c2) CloseConnection() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.Closed = true
}

func (self *c2) SetConnection(conn EasygoConn) IEndpointWithSocket {
	panic("抽象方法，子类必须实现")
}

func (self *c2) GetConnection() EasygoConn {
	panic("抽象方法，子类必须实现")
}

func (self *c2) CreateSendQueue() IQueue {
	return NewQueue()
}

func (self *c2) Start() { // implement
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	var f func()

	f = self.Me.RecvProc
	self.RecvJob = Spawn(f)

	f = self.Me.SendProc
	self.SendJob = Spawn(f)
}

func (self *c2) Join() {
	JoinAllJobs(self.RecvJob, self.SendJob)
}

func (self *c2) SendProc() {
	defer func() {
		if r := recover(); r != nil {
			self.Me.CloseConnection() // 确保 ReadProc 能退出来
			panic(r)
		}
	}()
	for {
		packet := self.SendQueue.Take()
		if packet == _STOP_SENTRY {
			break
		}
		data := packet.([]byte)
		socketError := self.Me.SendMsg(data)
		if socketError != nil {
			return
		}
		// self.Me.WriteComplete()
	}
	_ = self.Me.ShutdownWrite() // 半关闭，关闭写
}

func (self *c2) RecvProc() {

	defer func() {
		if r := recover(); r != nil {
			self.SendQueue.Put(_STOP_SENTRY) // 确保 SendProc 能退出来
			self.Me.OnDisconnected()
			panic(r) // 抛出异常，子类自行确保不要上抛网络相关的异常上来
		}
	}()

	var msgs [][]byte
	var socketError error
	for {
		msgs, socketError = self.Me.ReadMsg()
		if socketError != nil {
			break
		}
		for _, packet := range msgs {
			intercept, newPacket := self.Me.InterceptAndDeal(packet)
			if intercept {
				continue
			}
			isRequest, request, now := self.Me.RecvPacket(newPacket)
			if isRequest {
				req := request.(*base.Request)
				var serverType int32
				if com := req.GetCommon(); com != nil {
					serverType = com.GetServerType()
				}
				switch serverType {
				case SERVER_TYPE_SHOP, SERVER_TYPE_SQUARE:
					logs.Info("消息转发给其他服务器:", serverType, req.GetMethodName())
					self.Me.RoutingToOtherServer(newPacket, req)
				default:
					self.Me.DealRequest(now, req)
				}
			}
		}
	}
	if socketError == io.EOF { // 不是 io.EOF, 则不是 GRACEFULLY 关闭
		self.Mutex.Lock()
		if !self.HasStopSentry { // 我方没有调用过 Shutdown.是对方主动 shutdown
			self.HasStopSentry = true
			self.SendQueue.Put(_STOP_SENTRY) // 让发送队列中的发送完毕

			self.Mutex.Unlock()
			self.SendJob.Join(1 * time.Second) // 如果我方已经放了 _STOP_SENTRY,也可能正发送中就收到 fin
			self.Mutex.Lock()
		}
		self.Mutex.Unlock()
		self.Me.CloseConnection() // 真正地全关闭
		self.Me.OnDisconnected()
	} else { //客户端主动关闭情况下 也需要主动调用CloseConnection方法  否则tcp连接会处于close_wait状态
		self.SendQueue.Put(_STOP_SENTRY) // 确保 SendProc 能退出来
		self.Me.CloseConnection()        // 真正地全关闭
		self.Me.OnDisconnected()
	}
}

func (self *c2) InterceptAndDeal(packet []byte) (bool, []byte) { // 是否拦截并处理
	return false, packet
}

func (self *c2) SendPacket(packet []byte) { // override
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	if self.Closed || self.HasStopSentry {
		return
	}
	/* 此特性先不开启
	if SEND_QUEUE_SIZE != 0 && self.SendQueue.Size() >= SEND_QUEUE_SIZE { // 对待接收缓慢者,严格的手段,断开连接(服务器之间的连接不能这么做)
		text := fmt.Sprintf("发送队列满了,包数量 %v,即将关闭socket ", self.SendQueue.Size())
		logs.Info(text)

		self.Me.CloseConnection()
		self.RecvJob.DoCancel()
		return
	}
	*/
	self.SendQueue.Put(packet)
}

func (self *c2) Shutdown(timeouts ...time.Duration) { // override
	timeout := append(timeouts, time.Second*1)[0]

	self.Mutex.Lock()
	if self.Closed || self.HasStopSentry {
		self.Mutex.Unlock()
		return
	}
	self.HasStopSentry = true
	self.Mutex.Unlock()

	self.SendQueue.Put(_STOP_SENTRY) // 加入哨兵，让已经进了队列的数据发完.

	value := self.RecvJob.Get(true, timeout) // n 秒内要收到客户端的 fin,因为客户端可能收 fin 后但是就不愿意 shutdown (恶意的客户端)
	if _, ok := value.(ITimeoutError); ok {
		self.Me.CloseConnection() // 强行结束
	}
}
func (self *c2) RoutingToOtherServer(data []byte, request *base.Request) {
	logs.Info("c2==============")
}
