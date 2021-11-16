// 可持久化对象
package for_game

import (
	"game_server/easygo"
	"log"
	"sync"
)

var _ = log.Println

type IPersistence interface {
	OnBorn(kwargs ...easygo.KWAT)

	OnLoad()

	MarkDirty(isSet ...bool)
	OnDirty()
	IsDirty() bool
	ClearDirtyFlags()

	CreateLocker() sync.Locker
	GetLocker() sync.Locker
	SetLocker(lock sync.Locker)
}

type Persistence struct {
	Me         IPersistence
	DirtyEvent easygo.IEvent
	DirtyFlag  bool
	Mutex      sync.Locker
	IsAll      bool //是否全部保存
}

func (self *Persistence) Init(me IPersistence, kwargs ...easygo.KWAT) {
	self.Me = me
	self.DirtyEvent = easygo.NewEvent()
	dict := append(kwargs, easygo.KWAO)[0]
	f, ok := dict["DirtyEventHandler"]
	if ok {
		handler := f.(func(...bool))
		self.DirtyEvent.AddHandler(handler)
	}

	l, ok := dict["Locker"]
	if ok {
		self.Mutex = l.(sync.Locker)
	} else {
		self.Mutex = self.Me.CreateLocker()
	}
}

func (self *Persistence) OnBorn(kwargs ...easygo.KWAT) {
}

func (self *Persistence) OnLoad() {
}

func (self *Persistence) MarkDirty(isSet ...bool) {
	self.DirtyFlag = true
	self.IsAll = append(isSet, false)[0]
	self.Me.OnDirty()
}

func (self *Persistence) OnDirty() {
	self.DirtyEvent.Trigger(self.IsAll)
}

func (self *Persistence) IsDirty() bool { // impletement
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	return self.DirtyFlag
}

func (self *Persistence) ClearDirtyFlags() { // impletement
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	self.DirtyFlag = false
}

func (self *Persistence) CreateLocker() sync.Locker {
	return &easygo.RLock{}
}

func (self *Persistence) GetLocker() sync.Locker {
	return self.Mutex
}

func (self *Persistence) SetLocker(lock sync.Locker) {
	self.Mutex = lock
}
