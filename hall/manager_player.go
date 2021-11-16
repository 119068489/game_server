package hall

import (
	"sync"
)

type PlayerManager struct {
	sync.Map
}

func NewPlayerManager() *PlayerManager {
	p := &PlayerManager{}
	return p
}

func (self *PlayerManager) LoadPlayer(playerId PLAYER_ID) *Player {
	value, ok := self.Load(playerId)
	if ok {
		return value.(*Player)
	}
	return nil
}

var PlayerMgr = NewPlayerManager()
