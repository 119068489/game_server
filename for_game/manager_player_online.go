package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"

	"github.com/astaxie/beego/logs"
)

//玩家在线管理器，主要存放在线玩家当前在那个大厅上
type PlayerOnlineManager struct {
	PlayerList    map[PLAYER_ID]int32       //在线玩家列表[PLAYER_ID]SERVER_ID
	CutrList      map[PLAYER_ID]bool        //切后台玩家列表[PLAYER_ID]SERVER_ID
	OutReturnList map[PLAYER_ID]interface{} //玩家切后台时间信息列表
	Mutex         easygo.Mutex
}

func NewPlayerOnlineManager() *PlayerOnlineManager {
	p := &PlayerOnlineManager{}
	p.PlayerList = make(map[PLAYER_ID]int32)
	p.CutrList = make(map[PLAYER_ID]bool)
	p.OutReturnList = make(map[PLAYER_ID]interface{})
	return p
}

//玩家上线
func (self *PlayerOnlineManager) PlayerOnline(playerId PLAYER_ID, serverId int32) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.PlayerList[playerId] = serverId
	logs.Info("玩家上线了:", serverId, playerId)
}

//玩家下线
func (self *PlayerOnlineManager) PlayerOffline(playerId PLAYER_ID) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	delete(self.PlayerList, playerId)
	delete(self.CutrList, playerId)
}

func (self *PlayerOnlineManager) GetPlayerServerId(playerId PLAYER_ID) int32 {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.PlayerList[playerId]
}
func (self *PlayerOnlineManager) GetAllOnlinePlayers() *share_message.PlayerOnlineInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	onLines := &share_message.PlayerOnlineInfo{}
	for pid, sid := range self.PlayerList {
		p := &share_message.PlayerState{
			PlayerId: easygo.NewInt64(pid),
			ServerId: easygo.NewInt32(sid),
		}
		onLines.OnLines = append(onLines.OnLines, p)
	}
	return onLines
}

//检测玩家是否在线
func (self *PlayerOnlineManager) CheckPlayerIsOnLine(playerId PLAYER_ID) bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if self.PlayerList[playerId] != 0 {
		return true
	}
	return false
}

//设置玩家切后台状态
func (self *PlayerOnlineManager) SetPlayerIsCutBackstage(playerId PLAYER_ID, b bool) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.CutrList[playerId] = b
}

//检测玩家是否切后台
func (self *PlayerOnlineManager) CheckPlayerIsCutBackstage(playerId PLAYER_ID) bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.CutrList[playerId]
}

func (self *PlayerOnlineManager) SetPlayerOutReturnTime(playerId PLAYER_ID, state string, time int64) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	switch result := self.OutReturnList[playerId].(type) {
	case nil:
		self.OutReturnList[playerId] = make(map[string]int64)
		self.OutReturnList[playerId].(map[string]int64)[state] = time
	case map[string]int64:
		result[state] = time
	}
}

func (self *PlayerOnlineManager) GetPlayerOutReturnTime(playerId PLAYER_ID, state string) int64 {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	switch result := self.OutReturnList[playerId].(type) {
	case nil:
		return 0
	case map[string]int64:
		return result[state]
	}
	return 0
}
