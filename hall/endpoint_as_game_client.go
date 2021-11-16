// endpoint 代表一个玩家的连接

package hall

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"

	"github.com/astaxie/beego/logs"
)

var SpecialRpcName = []string{"RpcLogin", "RpcHeartbeat", "RpcMessageRequest"}

// H5 和 U3D 客户端连上来都是用这个表示
type IGameClientEndpoint interface {
	for_game.IClientEndpoint

	client_hall.IAccountHall2Client
	client_hall.IChatHall2Client
	client_hall.IHall2ShopClient
	client_hall.IHall2SquareClient
	client_hall.IHall2CoinShopClient
	client_hall.ILoveHall2Client
	client_hall.IESports2Client
	SetAssociativePlayer(player *Player)
	GetPlayer() *Player
}

//-----------------------------------

// 假装自己是一个 EndpointWithSocket,其实是被 mixin 到具体的 endpoint 对象
type MixIn struct {
	client_server.Server2Client
	client_hall.AccountHall2Client
	client_hall.ChatHall2Client
	client_hall.Hall2ShopClient
	client_hall.Hall2SquareClient
	client_hall.Hall2CoinShopClient
	client_hall.LoveHall2Client
	client_hall.ESports2Client
	Player       *Player
	GameClientEp IGameClientEndpoint
}

func (self *MixIn) Init(gameClientEp IGameClientEndpoint) {
	self.GameClientEp = gameClientEp
	self.AccountHall2Client.Init(gameClientEp)
	self.ChatHall2Client.Init(gameClientEp)
	self.Server2Client.Init(gameClientEp)
	self.Hall2ShopClient.Init(gameClientEp)
	self.Hall2SquareClient.Init(gameClientEp)
	self.Hall2CoinShopClient.Init(gameClientEp)
	self.LoveHall2Client.Init(gameClientEp)
}

func (self *MixIn) GetContexForDealRequest(methodName string, requestId uint64) (interface{}, interface{}) { // [override]
	epId := self.GameClientEp.GetEndpointId()
	ep := ClientEpMgr.LoadEndpoint(epId)
	if IsStopServer {
		if methodName == "RpcLogin" {
			return ep, easygo.NewFailMsg("服务器已关闭 ,放弃处理消息", for_game.FAIL_MSG_CODE_1011)
		} else {
			return ep, fmt.Sprintf("服务器已关闭 ,放弃处理消息 %s", methodName)
		}
	}
	for _, name := range SpecialRpcName {
		if name == methodName {
			return ep, self.GameClientEp
		}
	}
	if self.Player == nil && !easygo.Contain(for_game.NOT_NEED_LOGIN_PB, methodName) {
		return ep, fmt.Sprintf("收到 %s,非法操作，请先登录", methodName)
	}
	return ep, self.Player
}

func (self *MixIn) SetAssociativePlayer(player *Player) {
	self.Player = player
}
func (self *MixIn) GetPlayer() *Player {
	return self.Player
}

func (self *MixIn) CreateGoroutineForRequest() easygo.IGoroutine { // [override]
	return NewRequestGoroutine(self.GameClientEp.GetEndpointId())
}

// func (self *MixIn) CreateRateLimiter() *ratelimit.Bucket { // [override]
// 	return ratelimit.NewBucketWithQuantum(1000*time.Millisecond, 20, 2)
// }

//---------------------------------

// U3D 客户端
type GameClientTcpEndpoint struct {
	MixIn // 据测试，函数有二义性时，是广度搜索。就是想优先查找 MixIn 的函数，而不是 easygo.EndpointWithTcpSocket
	easygo.EndpointWithTcpSocket
}

func NewGameClientTcpEndpoint(serviceMap map[string]easygo.Pair, stubMap map[string]easygo.Pair, endpointId ENDPOINT_ID) *GameClientTcpEndpoint {
	p := &GameClientTcpEndpoint{}
	p.Init(serviceMap, stubMap, endpointId)
	return p
}

func (self *GameClientTcpEndpoint) Init(serviceMap map[string]easygo.Pair, stubMap map[string]easygo.Pair, endpointId ENDPOINT_ID) {
	self.MixIn.Init(self)
	self.EndpointWithTcpSocket.Init(self, _ServicesForGameClient, serviceMap, stubMap, endpointId)
}

//消息转发给指定服务器
func (self *GameClientTcpEndpoint) RoutingToOtherServer(data []byte, request *base.Request) {
	logs.Info("GameClientTcpEndpoint========")
	methodName, requestId := request.GetMethodName(), request.GetRequestId()
	pair, ok := self.ServiceMap[methodName]
	if !ok {
		logs.Warn(`收到 "%s" 请求,但是在 ServiceMap 中找不到这个 key.`, methodName)
		return
	}
	_, responseCls := pair[0], pair[1]
	var failInfo *base.Fail = nil
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
	player, ok1 := context.(*Player)
	if !ok1 { // 表示玩家不存在
		failInfo = easygo.NewFailMsg("无效的用户请求")
		logs.Debug("无效的用户请求")
		return
	}
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
				if e, ok := o.(easygo.IRpcInterrupt); ok { // 有可能在处理请求的过程中调用了另一个 rpc 方法,而这个 rpc 方法抛出了 “失败”
					reason, code := e.Reason(), e.Code()
					failInfo = &base.Fail{Reason: &reason, Code: &code} // 抛异常了,回个包给对端,免得对端死等回复
				} else {
					failInfo = self.Me.CreateFailMsgWhenPanic(o, methodName, requestId)
				}
				self.Me.SendFailMsg(requestId, failInfo)
			}
			if _, ok := o.(easygo.IRpcInterrupt); ok { // 重抛出非 IRpcInterrupt 异常
				panic(o) // 不能丢失 o 原来的类型
			} else {
				s := fmt.Sprintf("在分发 %s 时;%v", methodName, o) // 补上 rpc 方法名
				panic(s)                                       // o 失去了原来的类型，现在是 string 类型
			}
		}
	}()
	var serverType int32
	if com := request.GetCommon(); com != nil {
		serverType = com.GetServerType()
		request.Common.UserId = easygo.NewInt64(player.GetPlayerId())
		request.Common.Flag = easygo.NewInt32(for_game.MSG_FLAG)
	}
	server := PServerInfoMgr.GetIdelServer(serverType)
	if server == nil {
		logs.Error("转发时找不到服务器信息", serverType)
		if responseCls != "base.NoReturn" {
			failInfo = easygo.NewFailMsg("转发时找不到服务器信息:", easygo.AnytoA(serverType))
		}
		return
	}
	port := server.GetServerApiPort()
	if port == 0 {
		logs.Error("服务器api端口怎么会为0呢")
		if responseCls != "base.NoReturn" {
			failInfo = easygo.NewFailMsg("转发时找不到服务器信息", easygo.AnytoA(serverType))
		}
		return
	}
	u := "http://" + server.GetInternalIP() + ":" + easygo.AnytoA(port) + "/api"
	//分发到各 服务器 业务处理函数
	bs, err := for_game.DoBytesPost(u, data, request.GetCommon())
	if err != nil {
		if responseCls != "base.NoReturn" {
			failInfo = easygo.NewFailMsg("转发时找不到服务器信息", easygo.AnytoA(serverType))
		}
		return
	}
	if responseCls != "base.NoReturn" {
		self.Me.SendPacket(bs)
	}
	sended = true
}

//--------------------------------------

// H5 客户端
type GameClientWebEndpoint struct {
	MixIn // 据测试，函数有二义性时，是广度搜索。就是想优先查找 MixIn 的函数，而不是 easygo.EndpointWithWebSocket
	easygo.EndpointWithWebSocket
}

func NewGameClientWebEndpoint(serviceMap map[string]easygo.Pair, stubMap map[string]easygo.Pair, endpointId ENDPOINT_ID) *GameClientWebEndpoint {
	p := &GameClientWebEndpoint{}
	p.Init(serviceMap, stubMap, endpointId)
	return p
}

func (self *GameClientWebEndpoint) Init(serviceMap map[string]easygo.Pair, stubMap map[string]easygo.Pair, endpointId ENDPOINT_ID) {
	self.MixIn.Init(self)
	self.EndpointWithWebSocket.Init(self, _ServicesForGameClient, serviceMap, stubMap, endpointId)
}

//消息转发给指定服务器
func (self *GameClientWebEndpoint) RoutingToOtherServer(data []byte, request *base.Request) {
	methodName, requestId := request.GetMethodName(), request.GetRequestId()
	pair, ok := self.ServiceMap[methodName]
	if !ok {
		logs.Warn(`收到 "%s" 请求,但是在 ServiceMap 中找不到这个 key.`, methodName)
		return
	}
	_, responseCls := pair[0], pair[1]
	var failInfo *base.Fail = nil
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
	player, ok1 := context.(*Player)
	if !ok1 { // 表示玩家不存在
		failInfo = easygo.NewFailMsg("无效的用户请求")
		logs.Debug("无效的用户请求")
		return
	}
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
				if e, ok := o.(easygo.IRpcInterrupt); ok { // 有可能在处理请求的过程中调用了另一个 rpc 方法,而这个 rpc 方法抛出了 “失败”
					reason, code := e.Reason(), e.Code()
					failInfo = &base.Fail{Reason: &reason, Code: &code} // 抛异常了,回个包给对端,免得对端死等回复
				} else {
					failInfo = self.Me.CreateFailMsgWhenPanic(o, methodName, requestId)
				}
				self.Me.SendFailMsg(requestId, failInfo)
			}
			if _, ok := o.(easygo.IRpcInterrupt); ok { // 重抛出非 IRpcInterrupt 异常
				panic(o) // 不能丢失 o 原来的类型
			} else {
				s := fmt.Sprintf("在分发 %s 时;%v", methodName, o) // 补上 rpc 方法名
				panic(s)                                       // o 失去了原来的类型，现在是 string 类型
			}
		}
	}()
	var serverType int32
	if com := request.GetCommon(); com != nil {
		serverType = com.GetServerType()
		request.Common.UserId = easygo.NewInt64(player.GetPlayerId())
		request.Common.Flag = easygo.NewInt32(for_game.MSG_FLAG)
	}
	server := PServerInfoMgr.GetIdelServer(serverType)
	if server == nil {
		logs.Error("转发时找不到服务器信息", serverType)
		if responseCls != "base.NoReturn" {
			failInfo = easygo.NewFailMsg("转发时找不到服务器信息:", easygo.AnytoA(serverType))
		}
		return
	}
	port := server.GetServerApiPort()
	if port == 0 {
		logs.Error("服务器api端口怎么会为0呢")
		if responseCls != "base.NoReturn" {
			failInfo = easygo.NewFailMsg("转发时找不到服务器信息", easygo.AnytoA(serverType))
		}
		return
	}
	u := "http://" + server.GetInternalIP() + ":" + easygo.AnytoA(port) + "/api"
	//分发到各 服务器 业务处理函数
	bs, err := for_game.DoBytesPost(u, data, request.GetCommon())
	if err != nil {
		if responseCls != "base.NoReturn" {
			failInfo = easygo.NewFailMsg("转发时找不到服务器信息", easygo.AnytoA(serverType))
		}
		return
	}
	if responseCls != "base.NoReturn" {
		self.Me.SendPacket(bs)
	}
	sended = true
}

//--------------------------------------

type RequestGoroutine struct {
	easygo.Goroutine
	EndpointId ENDPOINT_ID
}

func NewRequestGoroutine(endpointId ENDPOINT_ID) *RequestGoroutine {
	p := &RequestGoroutine{}
	p.Init(endpointId)
	return p
}

func (self *RequestGoroutine) Init(endpointId ENDPOINT_ID) {
	self.EndpointId = endpointId
	self.Goroutine.Init(self)
}

// 发生异常时向客户端发送调用栈并提示到界面上
func (self *RequestGoroutine) OnPanic(callStack string, recoverVal interface{}) { // override
	if for_game.IS_FORMAL_SERVER {
		return
	}

	ep := ClientEpMgr.LoadEndpoint(self.EndpointId)
	if ep != nil {
		for_game.RpcToast(ep, "\n【服务端出错了，请截图到群里】\n"+callStack)
	}
}

var _ServicesForGameClient = map[string]interface{}{}

func RegisterServiceForGameClient(modName string, service interface{}) {
	_ServicesForGameClient[modName] = service
}
