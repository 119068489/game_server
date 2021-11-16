package for_game

import (
	"game_server/easygo"
)

type ConnectorManager struct {
	Mutex      easygo.RLock
	Connectors map[SERVER_ID]interface{} //服务器id为key
}

func NewConnectorManager() *ConnectorManager {
	p := &ConnectorManager{}
	p.Init()
	return p
}

func (self *ConnectorManager) Init() {
	self.Connectors = make(map[SERVER_ID]interface{})
}

func (self *ConnectorManager) LoadConnector(serverId SERVER_ID) interface{} {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value, ok := self.Connectors[serverId]
	if ok {
		return value
	}
	return nil
}

func (self *ConnectorManager) Delete(serverId SERVER_ID) { // override 为了让使用者清楚参数的类型
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	delete(self.Connectors, serverId)
}

func (self *ConnectorManager) Length() int {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	return len(self.Connectors)
}
func (self *ConnectorManager) Store(serverId SERVER_ID, con interface{}) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	self.Connectors[serverId] = con
}

func (self *ConnectorManager) GetConnectors() map[SERVER_ID]interface{} {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	cons := make(map[SERVER_ID]interface{})
	for k, v := range self.Connectors {
		cons[k] = v
	}
	return cons
}

func (self *ConnectorManager) GetRandServerId() int32 {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for serverId, _ := range self.Connectors {
		return serverId
	}
	return 0
}

//----------------------------------------------

// Login connector 管理器
//var LoginConMgr = NewLoginConnectorManager() // 以 serverId 关联
