package login

import (
	"time"

	"github.com/astaxie/beego/logs"

	"game_server/easygo"
)

var _CHECK_ENDPOINT_INTERVAL = 1 * time.Second //5 * time.Minute // 定时器间隔
var _KICK_ENDPOINT_DELAY = int64(15 * 1000)    //int64(9 * 60 * 1000)      // 踢除判断条件。毫秒
var _KICK_SEND_DELAY = int64(2 * 1000)         //int64(9 * 60 * 1000)      // 发送判断条件。毫秒

type GameClientEndpointManager struct {
	Mutex     easygo.RLock
	Endpoints map[ENDPOINT_ID]IGameClientEndpoint
}

func NewGameClientEndpointManager() *GameClientEndpointManager {
	p := &GameClientEndpointManager{}
	p.Init()
	return p
}

func (self *GameClientEndpointManager) Init() {
	self.Endpoints = make(map[ENDPOINT_ID]IGameClientEndpoint)
	easygo.AfterFunc(_CHECK_ENDPOINT_INTERVAL, self.CheckDisconnected)
}

// 定时检查太久没有发包的连接，踢掉
func (self *GameClientEndpointManager) CheckDisconnected() {
	defer easygo.AfterFunc(_CHECK_ENDPOINT_INTERVAL, self.CheckDisconnected)

	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, ep := range self.Endpoints {

		player := ep.GetPlayer()
		if player != nil {
			continue
		}
		stamp := ep.GetLastRecvStamp()
		now := easygo.FuncSet.GetTimestamp()
		if now > stamp+_KICK_ENDPOINT_DELAY {
			if player != nil {
				logs.Info("登录玩家心跳超时---", player.GetPlayerId())
			} else {
				logs.Info("登录连接心跳超时")
			}
			ep.Shutdown()
		} else if now > stamp+_KICK_SEND_DELAY {
			ep.SetIsSend(false)
		}
	}
}

func (self *GameClientEndpointManager) LoadEndpoint(endpointId ENDPOINT_ID) IGameClientEndpoint {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value, ok := self.Endpoints[endpointId]
	if ok {
		return value
	}
	return nil
}
func (self *GameClientEndpointManager) LoadEndpointByPid(pid int64) IGameClientEndpoint {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, ep := range self.Endpoints {
		player := ep.GetPlayer()
		if player.GetPlayerId() == pid {
			return ep
		}
	}
	return nil
}

func (self *GameClientEndpointManager) Delete(endpointId ENDPOINT_ID) { // override 为了让使用者清楚参数的类型
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	delete(self.Endpoints, endpointId)
}

func (self *GameClientEndpointManager) Length() int {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	return len(self.Endpoints)
}
func (self *GameClientEndpointManager) Store(endpointId ENDPOINT_ID, ep IGameClientEndpoint) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	self.Endpoints[endpointId] = ep
}

func (self *GameClientEndpointManager) GetEndpoints() map[ENDPOINT_ID]IGameClientEndpoint {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	endPoints := make(map[ENDPOINT_ID]IGameClientEndpoint)
	for k, v := range self.Endpoints {
		endPoints[k] = v
	}
	return endPoints
}

//-----------------------------------------------------------

type GameClientEndpointMapping struct {
	easygo.Mapping
}

func NewGameClientEndpointMapping() *GameClientEndpointMapping {
	p := &GameClientEndpointMapping{}
	p.Init(p)
	return p
}

func (self *GameClientEndpointMapping) LoadEndpoint(playerId PLAYER_ID) IGameClientEndpoint {
	value := self.Mapping.Load(playerId)
	if value != nil {
		return value.(IGameClientEndpoint)
	}
	return nil
}

func (self *GameClientEndpointMapping) StoreEndpoint(playerId PLAYER_ID, endpointId ENDPOINT_ID) bool {
	if endpointId == 0 {
		panic("endpoint id 不可能是 0")
	}
	fetch := func() interface{} {
		v := ClientEpMgr.LoadEndpoint(endpointId)
		if v == nil {
			return nil
		}
		return v
	}
	return self.Store(playerId, fetch)
}

func (self *GameClientEndpointMapping) Delete(playerId PLAYER_ID) { // overwrite 为了让使用者清楚参数的类型
	self.Mapping.Delete(playerId)
}

//----------------------------------------------

// 玩家 endpoint 管理器
var ClientEpMgr = NewGameClientEndpointManager() // 以 endpoint id 关联
var ClientEpMp = NewGameClientEndpointMapping()  // 以 玩家 id 关联
