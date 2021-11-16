package wish

import (
	"game_server/easygo"
	"sync"
)

// 水池管理器.
type PoolManager struct {
	Mutex easygo.RLock
	sync.Map
}

func NewPoolManager() *PoolManager {
	p := &PoolManager{}
	return p
}

func (self *PoolManager) LoadPool(poolId int64) *PoolObj {
	value, ok := self.Load(poolId)
	if ok {
		return value.(*PoolObj)
	}
	return nil
}
