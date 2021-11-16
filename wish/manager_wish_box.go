package wish

import (
	"game_server/easygo"
	"sync"
)

type WishBoxManager struct {
	Mutex easygo.RLock
	sync.Map
}

func NewWishBoxManager() *WishBoxManager {
	p := &WishBoxManager{}
	return p
}

func (self *WishBoxManager) LoadWishBox(boxId int64) *WishBoxObj {
	value, ok := self.Load(boxId)
	if ok {
		return value.(*WishBoxObj)
	}
	return nil
}
