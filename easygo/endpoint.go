package easygo

import (
	"reflect"
	"sync"
	"time"

	"fmt"
	"game_server/easygo/base"
	"net"
	"sync/atomic"

	"github.com/astaxie/beego/logs"
	"github.com/juju/ratelimit"
)

type ENDPOINT_ID = int32

//服务器类型
const (
	SERVER_TYPE_LOGIN      int32 = 1 //登录服
	SERVER_TYPE_HALL       int32 = 2 //大厅服
	SERVER_TYPE_BACKSTAGE  int32 = 3 //后台服
	SERVER_TYPE_SHOP       int32 = 4 //商场服
	SERVER_TYPE_STATISTICS int32 = 5 //统计服
	SERVER_TYPE_SQUARE     int32 = 6 //社交广场服
)

var _NoReturnMsg = &base.NoReturn{}
var FailMsg = &base.Fail{}
var EmptyMsg = &base.Empty{}

//========================================================================================
type IEndpointBase interface {
	IFinalable
	IMessageSender
	GetEndpointId() ENDPOINT_ID
	ProcessResponse([]byte)
	SendFailMsg(requestId uint64, failMsg IMessage)
	SendPacket(packet []byte)

	CreateFunctionSet() IFunctionSet
	CreateAsyncResult() IAsyncResult
	GetContexForDealRequest(methodName string, requestId uint64) (interface{}, interface{})
	NextRequestId() uint64
	GetMethodTimeout(methodName string) time.Duration

	RecvPacket(bytes []byte) (bool, IMessage, int64)
	DealRequest(recvTimeStamp int64, request *base.Request)

	// BeforeDispatchRpc
	Start()
	Shutdown(timeout ...time.Duration)
	String() string

	CreateFailMsgWhenPanic(o interface{}, methodName string, requestId uint64) *base.Fail

	GetDisconnectedEvent() *Event
	SetAddr(net.Addr)
	GetAddr() net.Addr

	OnDisconnected()

	CreateGoroutineForRequest() IGoroutine
	FetchMetexForDealRequest(methodName string, reqId uint64) sync.Locker
	GetLastRecvStamp() int64              // 最后收包时间戳
	CreateRateLimiter() *ratelimit.Bucket // 令牌桶,限速器
	SetUid(uid string)
	GetUid() string
	SetFlag(b bool)
	GetFlag() bool
	RoutingToOtherServer(data []byte, request *base.Request)
	SetIsSend(b bool)
	GetIsSend() bool
}

type EndpointBase struct {
	Finalable
	Me IEndpointBase

	Mutex Mutex

	FuncSet           IFunctionSet
	endpointId        ENDPOINT_ID
	Services          map[string]reflect.Value
	lastRequestId     uint64
	PeerLastRequestId uint64 // 对端最后请求 id
	pendingId         map[uint64]IAsyncResult
	ServiceMap        map[string]Pair
	stubMap           map[string]Pair

	Disconnected *Event

	Addr          net.Addr
	GenMutexByKey IGenMutexByKey
	LastRecvStamp int64             //最后收包时间戳
	Bucket        *ratelimit.Bucket // 令牌桶限速
	Uid           string
	Flag          bool //是否顶号
	IsSend        bool //是否可以发送
}

type c1 = EndpointBase

/*抽象类，不提供此方法
func NewEndpoint(maxSize int) *EndpointBase {}
*/

// services 的 value 是结构体或结构体指针都可以

func (self *c1) Init(me IEndpointBase, services map[string]interface{}, serviceMap map[string]Pair, stubMap map[string]Pair, endpointId ...ENDPOINT_ID) {
	self.Finalable.Init(me)

	self.Services = make(map[string]reflect.Value)
	for modName, service := range services {
		v := reflect.ValueOf(service)
		if v.Kind() == reflect.Ptr {
			if v.Elem().Kind() != reflect.Struct {
				panic("参数 service 必须是个 struct 对象或 struct 指针")
			}
		} else if v.Kind() != reflect.Struct {
			panic("参数 service 必须是个 struct 对象或 struct 指针")
		}

		self.Services[modName] = v
	}

	self.Me = me
	self.pendingId = make(map[uint64]IAsyncResult)
	self.FuncSet = me.CreateFunctionSet()
	self.endpointId = append(endpointId, 0)[0]
	self.lastRequestId = 0
	self.ServiceMap = serviceMap
	self.stubMap = stubMap
	self.Disconnected = NewEvent()
	self.GenMutexByKey = NewGenMutexByKey()
	self.LastRecvStamp = self.FuncSet.GetTimestamp()
	self.Bucket = self.Me.CreateRateLimiter()
	self.Flag = false
	self.IsSend = true
}

func (self *c2) CreateGoroutineForRequest() IGoroutine {
	return NewGoroutine()
}

func (self *c1) GetDisconnectedEvent() *Event {
	return self.Disconnected
}

func (self *c1) SetAddr(addr net.Addr) {
	self.Addr = addr
}

func (self *c1) GetAddr() net.Addr {
	return self.Addr
}
func (self *c1) SetUid(uid string) {
	self.Uid = uid
}

func (self *c1) GetUid() string {
	return self.Uid
}

func (self *c1) Start() {
	panic("请在子类实现")
}

// 生成request id
func (self *c1) NextRequestId() uint64 {
	for {
		v := atomic.AddUint64(&self.lastRequestId, 1) // 溢出后自动回转
		if v != 0 {                                   // 不能是0,0用于表示无需回复的请求
			return v
		}
	}
}

func (self *c1) SendPacket(packet []byte) { // 允许多个生产者往send_queue塞数据
	panic("请在子类实现")
}

func (self *c1) GetEndpointId() ENDPOINT_ID {
	return self.endpointId
}

// 主动关闭,关闭写,半关闭(应该理解成是一个 rpc 调用)
func (self *c1) Shutdown(timeout ...time.Duration) { // default=1
	panic("请在子类实现")
}

func (self *c1) CreatePool() {
	// return gevent.pool.Pool(100)  // 协程池,限制最大并发数
}

func (self *c1) CreateRateLimiter() *ratelimit.Bucket {
	return nil
}

func (self *c1) DealRequest(recvTimeStamp int64, request *base.Request) {
	// f = util.Functor(self.ProcessRequest)
	// job = self.pool.spawn(f, recvTimeStamp, request)
	// job.link(util.Functor(self.__remove_worker_job))
	// self.worker_job_group.add(job)

	// 临时用下面的替代. todo 要控制总量
	g := self.Me.CreateGoroutineForRequest()
	g.Start(self.ProcessRequest, recvTimeStamp, request)
}

// func (self *c1) __remove_worker_job(job){
// 	self.worker_job_group.discard(job)
// }

func (self *c1) String() string {
	return fmt.Sprintf(`endpoint id=%v,addr=%v`, self.endpointId, self.Addr)
}

// todo 这个接口不好理解。。要重构
func (self *c1) RecvPacket(bytes []byte) (bool, IMessage, int64) {
	if self.Bucket != nil {
		self.Bucket.Wait(1) // 取令牌,取不到就卡一会吧
	}
	packet := &base.Packet{}
	err := packet.Unmarshal(bytes)
	if err != nil {
		logs.Debug("收到的逻辑包有误,反序列化出错,收到字节流是%v", string(bytes))
		return false, nil, 0
	}
	if packet.GetType() == base.PacketType_TYPE_REQUEST {
		if len(self.ServiceMap) == 0 {
			logs.Debug("收到一个请求,但是 endpoint 不接受任何请求.")
			return false, nil, 0
		}
		request := &base.Request{}
		err := request.Unmarshal(packet.Serialized)
		if err != nil {
			logs.Warn("请求的逻辑包有误,反序列化出错 %v", packet.Serialized)
			return false, nil, 0
		}
		now := self.FuncSet.GetTimestamp()
		self.LastRecvStamp = now // 更新最后收包时间戳
		return true, request, now

	} else if packet.GetType() == base.PacketType_TYPE_RESPONSE {
		self.Me.ProcessResponse(packet.GetSerialized())
		return false, nil, 0
	} else {
		logs.Debug("对端恶意:收到未知类型的网络包:%v", packet)
		return false, nil, 0
	}
}

func (self *c1) GetLastRecvStamp() int64 {
	return self.LastRecvStamp
}
func (self *c1) SendFailMsg(requestId uint64, failMsg IMessage) {
	// logs.Debug("Send:ReqId:%v, Method:%v, Args:{%v}", requestId, "failMsg", failMsg)

	bytes := Marshal(failMsg)
	st := base.ResponseType_TYPE_FAIL
	response := &base.Response{ResponseId: &requestId, SubType: &st, Serialized: bytes}
	bytes = Marshal(response)

	t := base.PacketType_TYPE_RESPONSE
	packet := &base.Packet{Type: &t, Serialized: bytes}
	bytes = Marshal(packet)
	self.Me.SendPacket(bytes)
}

func (self *c1) GetContexForDealRequest(methodName string, requestId uint64) (interface{}, interface{}) { // 获取 context,在处理对端的请求时
	return self.Me, self.Me
}

func (self *c1) CreateFunctionSet() IFunctionSet {
	return FuncSet
}

/*
func (self *c1) is_timeout(request, now, recvTimeStamp, methodName){
	if request.timeout <= 0:  // 对端对执行此请求无时间要求
		return false
	if request.time_stamp:
		send_stamp = request.time_stamp
	else:
		send_stamp = recvTimeStamp * 1000  // 对端没有填此域则以收包时间算

	if now * 1000 - send_stamp < request.timeout:
		return false

	s = '收到"{}"请求,还没有调用就发现已经超时.发送时间戳{},当前时间戳{},要求{}毫秒内要处理.'
	s = s.format(methodName, request.time_stamp, now * 1000, request.timeout)
	logs.Debug(s)

	return true
}
*/
func (self *c1) BeforeDispatchRpc(methodName string) IMessage {
	return nil
}

func (self *c1) CreateFailMsgWhenPanic(o interface{}, methodName string, requestId uint64) *base.Fail {
	s := fmt.Sprintf("%v", o) // 如果 o 是 error 也能正确地使用 %v 格式化
	code := "upstream_except" // 上游发生异常
	return &base.Fail{Reason: &s, Code: &code}
}

func (self *c1) ProcessRequest(recvTimeStamp int64, request *base.Request) {
	methodName, requestId := request.GetMethodName(), request.GetRequestId()
	//logs.Info("ServiceMap----------->", methodName, self.ServiceMap)
	pair, ok := self.ServiceMap[methodName]
	if !ok {
		logs.Warn(`收到 "%s" 请求,但是在 ServiceMap 中找不到这个 key.`, methodName)
		return
	}
	requestCls, responseCls := pair[0], pair[1]

	var failInfo *base.Fail = nil
	sended := false
	defer func() {
		o := recover()
		if o == nil {
			if responseCls != "base.NoReturn" && !sended {
				if failInfo == nil {
					s := fmt.Sprintf("在分发 %s 时;必须回复对端", methodName)
					panic(s)
				}
				self.Me.SendFailMsg(requestId, failInfo)
			}
		} else {
			if responseCls != "base.NoReturn" {
				if e, ok := o.(IRpcInterrupt); ok { // 有可能在处理请求的过程中调用了另一个 rpc 方法,而这个 rpc 方法抛出了 “失败”
					reason, code := e.Reason(), e.Code()
					failInfo = &base.Fail{Reason: &reason, Code: &code} // 抛异常了,回个包给对端,免得对端死等回复
				} else {
					failInfo = self.Me.CreateFailMsgWhenPanic(o, methodName, requestId)
				}
				self.Me.SendFailMsg(requestId, failInfo)
			}
			if _, ok := o.(IRpcInterrupt); ok { // 重抛出非 IRpcInterrupt 异常
				panic(o) // 不能丢失 o 原来的类型
			} else {
				s := fmt.Sprintf("在分发 %s 时;%v", methodName, o) // 补上 rpc 方法名
				panic(s)                                       // o 失去了原来的类型，现在是 string 类型
			}
		}
	}()

	// now := self.FuncSet.GetTimestamp()

	// 尚未调用方法,已经超时
	// if self.is_timeout(request, now, recvTimeStamp, methodName):
	// 	if responseCls != "base.NoReturn":
	// 		failInfo = f'无法在要求的时间内完成处理{methodName}.', ''
	// 	return

	reqMsg := NewMessage(requestCls)
	if request.Serialized != nil { // 当消息体没有域时不走这里，比如 base.Empty
		err := reqMsg.Unmarshal(request.Serialized)
		if err != nil {
			s := fmt.Sprintf("rpc调用%v时,反序列化 msg 失败,很可能消息定义不一致.err=%v", methodName, err)
			logs.Debug(s)
			if responseCls != "base.NoReturn" {
				failInfo = &base.Fail{Reason: &s}
			}
			return
		}
	}
	ep, context := self.Me.GetContexForDealRequest(methodName, requestId)
	if ep == nil {
		s := "获取不到 第 1 个 Contex"
		if requestId != 0 {
			failInfo = &base.Fail{Reason: &s}
		}
		logs.Debug(s)
		return
	}
	if s, ok := context.(string); ok { // 表示没有通过授权或者上下文不存在
		if responseCls != "base.NoReturn" {
			failInfo = &base.Fail{Reason: &s}
		}
		logs.Debug(s)
		return
	}
	if fail, ok := context.(*base.Fail); ok { // 表示没有通过授权或者上下文不存在
		if responseCls != "base.NoReturn" {
			failInfo = fail
		}
		logs.Debug(fail)
		return
	}

	var method reflect.Value
	for _, service := range self.Services {
		method = service.MethodByName(methodName)
		if method.IsValid() {
			break
		}
	}
	if !method.IsValid() {
		s := fmt.Sprintf("找不到 %v 方法", methodName)
		logs.Debug(s)
		if responseCls != "base.NoReturn" {
			failInfo = &base.Fail{Reason: &s}
		}
		return
	}
	if method.Kind() != reflect.Func {
		s := fmt.Sprintf("%v 不是方法", methodName)
		logs.Debug(s)
		if responseCls != "base.NoReturn" {
			failInfo = &base.Fail{Reason: &s}
		}
		return
	}
	args := []reflect.Value{
		reflect.ValueOf(ep),
		reflect.ValueOf(context),
		reflect.ValueOf(reqMsg),
	}
	if request.GetCommon() != nil {
		args = append(args, reflect.ValueOf(request.GetCommon()))
	}
	// logs.Debug("Receive:ReqId:%v, Method:%v, Args:{%v}", requestId, methodName, reqMsg)

	mutex := self.Me.FetchMetexForDealRequest(methodName, requestId)
	if mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}
	trace := fmt.Sprintf("addr:%v,endpoint id:%d,message 是 %v", self.Addr, self.endpointId, reqMsg)
	defer RecoverAndRePanic(trace, false)
	results := method.Call(args) // 分发到各 rpc 业务处理函数

	// TODO 记录时间消耗
	if len(results) == 0 {
		s := fmt.Sprintf("%v 方法函数原型错了,返回值类型必须是 IMessage", methodName)
		panic(s)
	}

	result := results[0].Interface()

	var respMsg IMessage
	if result == nil {
		if responseCls == "base.NoReturn" {
			return
		}
		if responseCls == "base.Empty" {
			respMsg = EmptyMsg
		} else {
			s := fmt.Sprintf("%v 的 response 不是 base.NoReturn 或 base.Empty，所以返回值不能是 nil", methodName)
			panic(s)
		}
	} else {
		if responseCls == "base.NoReturn" {
			s := fmt.Sprintf("%v 的 response 是 base.NoReturn，所以返回值必须是 nil", methodName)
			panic(s)
		}
		respMsg = result.(IMessage)
		if msg, ok := respMsg.(*base.Fail); ok { // 是失败消息
			failInfo = msg
			return
		}
	}

	serialized, err := respMsg.Marshal()
	// logs.Debug("Send:ReqId:%v, Method:%v, Args:{%v}", request.GetRequestId(), methodName, respMsg)
	if err != nil {
		s := fmt.Sprintf("序列化回复 %v 调用的 response 消息出错;%v", methodName, err.Error())
		panic(s)
	}

	// 调用完成后发现超时了,就不发送给对端,省流量,反正对端也提前结束了
	// if (request.time_stamp > 0 and request.timeout > 0) and (
	// 		self.FuncSet.GetTimestamp() * 1000 - request.time_stamp >= request.timeout):
	// 	s = '调用完{methodName}方法后回复响应时,已经超时.'

	// 	failInfo = s, ''
	// 	return
	st := base.ResponseType_TYPE_SUCCESS
	response := base.Response{ResponseId: &requestId, SubType: &st, Serialized: serialized}
	bytes := Marshal(&response)

	t := base.PacketType_TYPE_RESPONSE
	packet := base.Packet{Type: &t, Serialized: bytes}
	bytes = Marshal(&packet)
	self.Me.SendPacket(bytes)
	sended = true

}

func (self *c1) FetchMetexForDealRequest(methodName string, reqId uint64) sync.Locker {
	mutex := self.GenMutexByKey.GenMutex(methodName)
	return mutex
}

func (self *c1) ProcessResponse(serialized []byte) {
	response := base.Response{}
	err := response.Unmarshal(serialized)
	if err != nil {
		logs.Debug("收到的响应包有误,反序列化出错,字节流是 %v", serialized)
		return
	}
	respId := response.GetResponseId()
	if respId == 0 {
		logs.Debug("收到 rpc 回复，但是 ResponseId 为 0")
		return
	}
	self.Mutex.Lock()
	result, ok := self.pendingId[respId]
	self.Mutex.Unlock()
	if !ok { // 超时被弹掉或不存在的 id
		logs.Debug("收到 rpc 回复，但是超时了，所以丢弃")
		return
	}
	if response.GetSubType() == base.ResponseType_TYPE_SUCCESS {
		ri := &ResponseInfo{true, response.Serialized}
		result.Set(ri)
	} else {
		ri := &ResponseInfo{false, response.Serialized}
		result.Set(ri)
	}
}

// func (self *c1) cancel_pending_by_rpc_name(rpc_name){  // 不等回复了,某一类的rpc全部不等
// 	d = self.pending_rpc_name.get(rpc_name)
// 	if not d:
// 		return
// 	for requestId, oAsyncResult in d.items():
// 		e = gevent.GreenletExit()
// 		oAsyncResult.set_exception(e)
// 		logs.Debug(f'中止rpc:{rpc_name},{self}')
// }

func (self *c1) CreateAsyncResult() IAsyncResult {
	return NewAsyncResult()
}

func (self *c1) CallRpcMethod(methodName string, reqMsg IMessage, common ...*base.Common) (IMessage, IRpcInterrupt) { // override
	// log.Println("CallRpcMethod")
	var requestId uint64 = 0
	var timeout time.Duration = 0
	var result IAsyncResult

	pair, ok := self.stubMap[methodName]
	if !ok {
		s := fmt.Sprintf(`"%s"在 rpc info map 中没有定义`, methodName)
		panic(s)
	}
	requestCls, responseCls := pair[0], pair[1]

	if requestCls == "base.Empty" { //
		if IsInterfaceNil(reqMsg) {
			reqMsg = EmptyMsg
		} else {
			panic("request 消息是 base.Empty 时传 nil 过来就行了")
		}
	}

	requestId = self.Me.NextRequestId()
	if responseCls != "base.NoReturn" { // 需要回复的
		result = self.Me.CreateAsyncResult()
		self.Mutex.Lock()
		self.pendingId[requestId] = result
		self.Mutex.Unlock()

		defer func() {
			self.Mutex.Lock()
			delete(self.pendingId, requestId)
			self.Mutex.Unlock()

			// async_results = self.pending_rpc_name.get(methodName, nil) // todo 此特性后期再加
			// if async_results:
			// 	async_results.pop(requestId, nil)
		}()
		// self.pending_rpc_name.setdefault(methodName, {})[requestId] = result // todo 此特性后期再加
		// timeout = getattr(reqMsg, 'Timeout', nil)  // 如果请求消息里本身有这个timeout字段(default值也可以拿到的)
		// if timeout is nil:
		// 	timeout = self.GetMethodTimeout(methodName)
		// todo 改为尝试从 reqMsg 字段中拿
		timeout = self.Me.GetMethodTimeout(methodName)
	}
	packet := self.FuncSet.MakeRequestPacket(methodName, reqMsg, requestId, timeout, common...)
	self.Me.SendPacket(packet)
	// logs.Debug("Send:ReqId:%v, Method:%v, Args:{%v}", requestId, methodName, reqMsg)

	// 调用虚函数，让使用者可以获得发了什么字节流(战斗录相用得到这个特性)
	if responseCls == "base.NoReturn" { // 无需回复的
		// 肯定会成功。返回一个 NoReturn 实例 ,简化后面的强制转换等操作
		return _NoReturnMsg, nil
	}

	// TODO: getcurrent()当前协程放入哪里进行监控
	value, err := result.GetUntilTimeout(timeout)
	if err != nil {
		reason, code := " 对端超时未回复", "rpc_timeout"
		failMsg := &base.Fail{Reason: &reason, Code: &code} // 伪造一个消息,其实不是对端回复的
		return nil, NewRpcTimeout(methodName, failMsg)
	}
	info := value.(*ResponseInfo)
	if info.RpcOk {
		respMsg := NewMessage(responseCls)
		e := respMsg.Unmarshal(info.Bytes)
		if e != nil { // 反序列化出错
			logs.Debug("收到的 Response 包有误,反序列化出错,收到字节流是%v", info.Bytes)
		}
		return respMsg, nil
	} else {
		failMsg := &base.Fail{}
		e := failMsg.Unmarshal(info.Bytes)
		if e != nil { // 反序列化出错
			logs.Debug("收到的 Fail 包有误,反序列化出错,收到字节流是%v", info.Bytes)
		}
		return nil, NewRpcFail(methodName, failMsg)
	}
}

func (self *c1) GetMethodTimeout(methodName string) time.Duration { // 返回 0 表示不超时
	return 20 * time.Second // 20秒,永不超时不好啊!若是对端永远不回,某个协程就永远跳不回去了
}

// func (self *c1) follow_up(timeout=3){
// }

func (self *c1) OnDisconnected() {
	self.Disconnected.Trigger(self.Me)
}
func (self *c1) SetFlag(b bool) {
	self.Flag = b
}
func (self *c1) GetFlag() bool {
	return self.Flag
}
func (self *c1) SetIsSend(b bool) {
	self.IsSend = b
}
func (self *c1) GetIsSend() bool {
	return self.IsSend
}
