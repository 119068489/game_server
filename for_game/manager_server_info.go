package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"
	"sort"
)

//服务器管理器，管理连接服务器信息
type ServerInfoManager struct {
	HallServerInfo map[SERVER_ID]*share_message.ServerInfo
	Mutex          easygo.RLock
}

func NewServerInfoManager() *ServerInfoManager { // services map[string]interface{},
	p := &ServerInfoManager{}
	p.Init()
	return p
}

func (self *ServerInfoManager) Init() {
	self.HallServerInfo = make(map[SERVER_ID]*share_message.ServerInfo)
}
func (self *ServerInfoManager) AddServerInfo(srv *share_message.ServerInfo) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.HallServerInfo[srv.GetSid()] = srv
}
func (self *ServerInfoManager) DelServerInfo(id SERVER_ID) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	delete(self.HallServerInfo, id)
}
func (self *ServerInfoManager) GetServerInfo(serverId SERVER_ID) *share_message.ServerInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	srvInfo, ok := self.HallServerInfo[serverId]
	if !ok {
		return nil
	}
	return srvInfo
}
func (self *ServerInfoManager) ChangeServerState(serverId SERVER_ID, st int32) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	srv := self.GetServerInfo(serverId)
	srv.State = easygo.NewInt32(st)
}

//负载均衡，分配一台大厅服务器
func (self *ServerInfoManager) GetIdelServer(t int32) *share_message.ServerInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var temMap map[int32]int32
	temMap = make(map[int32]int32, len(self.HallServerInfo))
	for k, v := range self.HallServerInfo {
		if v.GetType() == t {
			temMap[k] = v.GetConNum()
		}
	}
	if len(temMap) > 0 {
		sid := self.SortMapByValue(temMap)
		return self.HallServerInfo[sid]
	}
	return self.HallServerInfo[0]
}

func (self *ServerInfoManager) GetAllServers(t int32) []*share_message.ServerInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var temMap []*share_message.ServerInfo
	temMap = make([]*share_message.ServerInfo, 0, len(self.HallServerInfo))
	for _, v := range self.HallServerInfo {
		if v.GetType() == t {
			temMap = append(temMap, v)
		}
	}
	return temMap
}

//连接数增加
func (self *ServerInfoManager) AddConNum(sid SERVER_ID) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, v := range self.HallServerInfo {
		if v.GetSid() == sid {
			v.ConNum = easygo.NewInt32(v.GetConNum() + 1)
			break
		}
	}
}

//连接数减少
func (self *ServerInfoManager) DelConNum(sid SERVER_ID) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, v := range self.HallServerInfo {
		if v.GetSid() == sid {
			v.ConNum = easygo.NewInt32(v.GetConNum() - 1)
			break
		}
	}
}

//------------------------------
// A data structure to hold a key/value pair.
type EGOPair struct {
	Key   int32
	Value int32
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type EGOPairList []EGOPair

func (p EGOPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p EGOPairList) Len() int           { return len(p) }
func (p EGOPairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// 返回人数最少的服务器ID
func (self *ServerInfoManager) SortMapByValue(m map[int32]int32) int32 {
	p := make(EGOPairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = EGOPair{k, v}
		i++
	}
	sort.Sort(p)
	return p[0].Key
}
