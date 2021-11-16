package for_game

import (
	"game_server/easygo"
)

//映射已经实例化的群信息到指定服务器上
type TeamOnlineManager struct {
	TeamList map[int64]int32 //在线群:在游戏服
	Mutex    easygo.Mutex
}

func NewTeamOnlineManager() *TeamOnlineManager {
	p := &TeamOnlineManager{}
	p.TeamList = make(map[int64]int32)
	return p
}

//玩家上线
func (self *TeamOnlineManager) TeamOnline(teamId int64, serverId int32) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.TeamList[teamId] = serverId
}

//玩家下线
func (self *TeamOnlineManager) TeamOffline(teamId int64) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	delete(self.TeamList, teamId)
}

func (self *TeamOnlineManager) GetTeamServerId(teamId int64) int32 {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.TeamList[teamId]
}

//检测玩家是否在线
func (self *TeamOnlineManager) CheckTeamIsOnLine(teamId int64) bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if self.TeamList[teamId] != 0 {
		return true
	}
	return false
}
