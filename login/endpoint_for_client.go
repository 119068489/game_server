// endpoint 代表一个玩家的连接

package login

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/client_login"
	"game_server/pb/client_server"
)

var SpecialRpcName = []string{"RpcLoginHall", "RpcHeartbeat", "RpcRegister", "RpcClientGetCode", "RpcCheckMessageCode", "RpcCheckAccountVaild", "RpcForgetLoginPassword", "RpcBtnClick", "RpcPageRegLogLoad"}

// H5 和 U3D 客户端连上来都是用这个表示
type IGameClientEndpoint interface {
	for_game.IClientEndpoint
	client_login.ILogin2Client
	SetAssociativePlayer(player *Player)
	GetPlayer() *Player
	GetConType() int
	SetConType(t int)
	AddAddrs(addr string)
	GetClientAddr() string
}

//-----------------------------------

// 假装自己是一个 EndpointWithSocket,其实是被 mixin 到具体的 endpoint 对象
type MixIn struct {
	client_server.Server2Client
	client_login.Login2Client

	Player       *Player
	GameClientEp IGameClientEndpoint
	ConnType     int
	Addrs        []string //连接的地址:包括转发
}

func (self *MixIn) Init(gameClientEp IGameClientEndpoint) {
	self.GameClientEp = gameClientEp
	self.Server2Client.Init(gameClientEp)
	self.Login2Client.Init(gameClientEp)
}

func (self *MixIn) GetContexForDealRequest(methodName string, requestId uint64) (interface{}, interface{}) { // [override]
	epId := self.GameClientEp.GetEndpointId()
	ep := ClientEpMgr.LoadEndpoint(epId)

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
func (self *MixIn) GetConType() int {
	return self.ConnType
}
func (self *MixIn) SetConType(t int) {
	self.ConnType = t
}

//存储连接地址，包括转发
func (self *MixIn) AddAddrs(addr string) {
	self.Addrs = append(self.Addrs, addr)
}

//获取client真是连接地址
func (self *MixIn) GetClientAddr() string {
	if len(self.Addrs) > 0 {
		return self.Addrs[len(self.Addrs)-1]
	}
	return ""
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
