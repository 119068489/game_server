package easygo

import (
	"runtime"
)

// 与 python 不同，golang 不能对同一个对象设置多次 runtime.SetFinalizer
// 所以设计出下面这个类

type IFinalable interface {
	AddFinalizer(function interface{})
}

type Finalable struct {
	Me              IFinalable
	hasSetFinalizer bool
	event           IEvent // *Event
}

func NewFinalable() *Finalable {
	p := &Finalable{}
	p.Init(p)
	return p
}

func (self *Finalable) Init(me IFinalable) {
	self.Me = me
	self.event = NewEvent()

}
func (self *Finalable) AddFinalizer(function interface{}) {
	self.event.AddHandler(function)
	if self.hasSetFinalizer {
		return
	}
	self.hasSetFinalizer = true

	// 1.如果 self.event 是对象,则 event := self.event 会拷贝一次，导致多次调用 AddFinalizer 只会调用 1 次 callback,因为只闭包了只有一个回调的 event upvalue
	// 2.如果 self.event 是对象,且 event := &self.event ,仍然会引导起循环引用(摸不清闭包的原则)。
	// 3.所以 self.event 类型必须是指针或接口类型
	event := self.event

	finalizer := func(obj IFinalable) {
		event.Trigger(obj) // 绑定的 event 是 upvalue,千万别绑定 self.event,会引起循环引用
	}
	runtime.SetFinalizer(self.Me, finalizer) // finalizer 不可以是成员函数，猜测:如果是成员函数会循环引用，不会析构
}
