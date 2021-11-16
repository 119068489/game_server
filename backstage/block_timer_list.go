//后台定时器管理

package backstage

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

func (t *Timer) Init() {
	t.Timers = make(map[int64]*easygo.Timer)
	kwargs := easygo.KWAT{
		"DirtyEventHandler": t.DirtyEventHandler,
	}
	t.Persistence.Init(t, kwargs)
}

func (t *Timer) DirtyEventHandler(isAll ...bool) {

}

func (t *Timer) OnBorn(kwargs ...easygo.KWAT) { // override

}

func (t *Timer) OnLoad() { // override

}

func (t *Timer) GetTimerList() map[int64]*easygo.Timer { // override
	return t.Timers
}

func (t *Timer) DelTimerList(key int64) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	timer, exist := t.Timers[key]
	if exist {
		timer.Stop()
	}

	_, ok := t.Timers[key]
	if ok {
		delete(t.Timers, key)
	}
}

func (t *Timer) AddTimerList(key int64, timer *easygo.Timer) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	t.Timers[key] = timer
}

func (t *Timer) GetTimerById(key int64) *easygo.Timer {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	timer, ok := t.Timers[key]
	if !ok {
		return nil
	}
	return timer
}
