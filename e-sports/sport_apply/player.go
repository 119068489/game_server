package sport_apply

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/client_server"
)

//=======================================

type Player struct {
	*for_game.RedisPlayerBaseObj
	PlayerId PLAYER_ID
	Mutex    easygo.RLock
}

func NewPlayer(playerId PLAYER_ID) *Player {
	p := &Player{}
	p.Init(playerId)
	return p
}

func (self *Player) Init(playerId PLAYER_ID) {
	self.PlayerId = playerId
}

func (self *Player) GetPlayerId() PLAYER_ID {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.PlayerId
}

func (self *Player) OnLoadFromDB() {
	self.RedisPlayerBaseObj = for_game.GetRedisPlayerBase(self.PlayerId)
}

func (self *Player) GetAllPlayerInfo() *client_server.AllPlayerMsg {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return &client_server.AllPlayerMsg{}
}
