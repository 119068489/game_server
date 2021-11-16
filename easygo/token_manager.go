// 别看了，作废， token 机制太容易死锁
package easygo

import (
	"sync"
)

type ITokenManager interface {
	Put(key interface{}, item interface{})
	Size() int
}

type TokenManager struct {
	Me ITokenManager

	globalMutex Mutex

	objMaps      map[interface{}]interface{} // 不存在这一项时表示对象不存在，存在这一项但 value 为 nil 时表示被借走了
	refCountMaps map[interface{}]int         // 引用计数

	mutexMaps map[interface{}]*Mutex
	condMaps  map[interface{}]*sync.Cond
}

func NewTokenManager() *TokenManager {
	p := &TokenManager{}
	p.Init(p)
	return p
}

func (self *TokenManager) Init(me ITokenManager) {
	self.Me = me

	self.objMaps = make(map[interface{}]interface{})
	self.refCountMaps = make(map[interface{}]int)
	self.mutexMaps = make(map[interface{}]*Mutex)
	self.condMaps = make(map[interface{}]*sync.Cond)

}

// 第一次放入对象，同一对象只允许放一次
func (self *TokenManager) Put(key interface{}, item interface{}) {
	self.globalMutex.Lock()
	defer self.globalMutex.Unlock()

	self.objMaps[key] = item
	self.refCountMaps[key] = 1

	mutex := new(Mutex)
	self.mutexMaps[key] = mutex
	self.condMaps[key] = sync.NewCond(mutex)
}

// 要区分没有此对象 还是 被别人借走了，没有此对象返回 nil,被借走了则阻塞等待
func (self *TokenManager) Borrow(key interface{}) (interface{}, func()) {
	mutex, cond := self.getCondMutexByKey(key, false)
	if mutex == nil || cond == nil { // 没有此对象
		return nil, nil
	}
	var ok bool
	var item interface{}
	mutex.Lock() //似乎不用 Lock 也可以
	defer mutex.Unlock()
	for {
		self.globalMutex.Lock()
		item, ok = self.objMaps[key]
		if ok && item != nil {
			self.objMaps[key] = nil // 表示借走了，不能删除,只能置为 nil
		}
		self.globalMutex.Unlock()

		if !ok { //根本没有这对象
			return nil, nil
		}
		if item != nil {
			// log.Println("find yes")
			break
		}
		// log.Println("wait begin")

		cond.Wait() // todo 加入超时机制
		// log.Println("wait end")

	}
	giveBackFunc := func() {
		self.GiveBackToken(key, item)
	}
	return item, giveBackFunc
}

//被 defer 使用
func (self *TokenManager) GiveBackToken(key interface{}, item interface{}) {
	mutex, cond := self.getCondMutexByKey(key, false)
	if mutex == nil || cond == nil {
		panic("wtf")
	}

	self.globalMutex.Lock()
	defer self.globalMutex.Unlock()

	self.objMaps[key] = item

	mutex.Lock()
	defer mutex.Unlock()
	cond.Signal() //调用本方法时，建议（但并非必须）保持c.L的锁定。

}
func (self *TokenManager) Size() int {
	self.globalMutex.Lock()
	defer self.globalMutex.Unlock()
	return len(self.objMaps) // 包含了被借走的
}

func (self *TokenManager) CopyKeys() []interface{} {
	self.globalMutex.Lock()
	defer self.globalMutex.Unlock()

	len := len(self.objMaps)
	keys := make([]interface{}, len)
	for k, _ := range self.objMaps {
		keys = append(keys, k)
	}
	return keys

}

// 增加引用计数
func (self *TokenManager) IncreaseRefCount(referrer IFinalable, key interface{}) {
	self.globalMutex.Lock()
	defer self.globalMutex.Unlock()

	self.refCountMaps[key] += 1

	f := Functor(self.DecreaseRefCount, key)
	referrer.AddFinalizer(f)

}

// 回调函数
// runtime.SetFinalizer(ptr, tokenTokenManagerObj.Finalizer) 还得想办法闭包一个 id 过来
// 真实参数是 obj *Foo
func (self *TokenManager) DecreaseRefCount(referrer IFinalable, key interface{}) {
	self.globalMutex.Lock()
	defer self.globalMutex.Unlock()
	// _ = referrer.Key()

	//根据 id 减少引用计数
	self.refCountMaps[key] -= 1
	if self.refCountMaps[key] == 0 {
		delete(self.objMaps, key)
		delete(self.refCountMaps, key)
		delete(self.mutexMaps, key)
		delete(self.condMaps, key)
	}
}

func (self *TokenManager) getCondMutexByKey(key interface{}, gen bool) (*Mutex, *sync.Cond) {
	self.globalMutex.Lock()
	defer self.globalMutex.Unlock()

	mutex, ok := self.mutexMaps[key]
	if !ok && gen {
		mutex = new(Mutex)
		self.mutexMaps[key] = mutex
	}

	cond, ok := self.condMaps[key]
	if !ok && gen {
		cond = sync.NewCond(mutex)
		self.condMaps[key] = cond
	}
	return mutex, cond

}
