package easygo

import (
	"sync"
)

type IGenMutexByKey interface {
	GenMutex(key interface{}) sync.Locker
}

type GenMutexByKey struct {
	Me IGenMutexByKey

	Mutex     RLock
	MutexMaps map[interface{}]*RLock
}

func NewGenMutexByKey() *GenMutexByKey {
	p := &GenMutexByKey{}
	p.Init(p)
	return p
}

func (self *GenMutexByKey) Init(Me IGenMutexByKey) {
	self.Me = Me

	self.MutexMaps = make(map[interface{}]*RLock)
}

func (self *GenMutexByKey) GenMutex(key interface{}) sync.Locker {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	mutex, ok := self.MutexMaps[key]
	if !ok {
		mutex = new(RLock)
		self.MutexMaps[key] = mutex // think, 什么时候移除呢？
	}
	return mutex
}
