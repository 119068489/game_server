package for_game

import (
	"game_server/easygo"
)

type EndpointManager struct {
	Mutex     easygo.RLock
	Endpoints map[SERVER_ID]interface{} //服务器endpoint Id为key
}

func NewEndpointManager() *EndpointManager {
	p := &EndpointManager{}
	p.Init()
	return p
}

func (self *EndpointManager) Init() {
	self.Endpoints = make(map[SERVER_ID]interface{})
}

func (self *EndpointManager) LoadEndpoint(sid SERVER_ID) interface{} {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value, ok := self.Endpoints[sid]
	if ok {
		return value
	}
	return nil
}

func (self *EndpointManager) Delete(sid SERVER_ID) { // override 为了让使用者清楚参数的类型
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	delete(self.Endpoints, sid)
}

func (self *EndpointManager) Length() int {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	return len(self.Endpoints)
}
func (self *EndpointManager) Store(sid SERVER_ID, ep interface{}) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	self.Endpoints[sid] = ep
}

func (self *EndpointManager) GetEndpoints() map[SERVER_ID]interface{} {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	endPoints := make(map[SERVER_ID]interface{})
	for k, v := range self.Endpoints {
		endPoints[k] = v
	}
	return endPoints
}

//----------------------------------------------

// 后台 endpoint 管理器
//var BackStageEpMgr = NewEndpointManager() // 以 serverId 关联
