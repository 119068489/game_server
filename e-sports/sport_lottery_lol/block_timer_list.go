//LOL定时器管理

package sport_lottery_lol

import (
	"game_server/easygo"
	"game_server/for_game"
)

type Timer struct {
	for_game.Persistence `bson:"-"`
	Timers               map[int64]*easygo.Timer `bson:"-"`
}

func NewTimer() *Timer {
	p := &Timer{}
	p.Init()
	return p
}

func (self *Timer) Init() {
	self.Timers = make(map[int64]*easygo.Timer)
	kwargs := easygo.KWAT{
		"DirtyEventHandler": self.DirtyEventHandler,
	}
	self.Persistence.Init(self, kwargs)
}

func (self *Timer) DirtyEventHandler(isAll ...bool) {

}

func (self *Timer) OnBorn(kwargs ...easygo.KWAT) { // override

}

func (self *Timer) OnLoad() { // override

}

func (self *Timer) GetTimerList() map[int64]*easygo.Timer { // override
	return self.Timers
}

func (self *Timer) DelTimerList(key int64) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	timer, exist := self.Timers[key]
	if exist {
		timer.Stop()
	}

	_, ok := self.Timers[key]
	if ok {
		delete(self.Timers, key)
	}
}

func (self *Timer) AddTimerList(key int64, timer *easygo.Timer) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	self.Timers[key] = timer
}

func (self *Timer) GetTimerById(key int64) *easygo.Timer {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	timer, ok := self.Timers[key]
	if !ok {
		return nil
	}
	return timer
}
