package easygo

import (
	"reflect"
)

type IEvent interface {
	AddHandler(function interface{})
	Trigger(args ...interface{})
}

type Event struct {
	handler []reflect.Value
}

func NewEvent() *Event {
	p := &Event{}
	p.Init()
	return p
}
func (self *Event) Init() { //

}

func (self *Event) AddHandler(function interface{}) { //
	f := reflect.ValueOf(function)
	if f.Kind() != reflect.Func {
		panic("参数必须是个 function")
	}
	self.handler = append(self.handler, f)
}

func (self *Event) Trigger(args ...interface{}) { //
	for _, f := range self.handler {
		in := make([]reflect.Value, len(args))
		i := 0
		for _, arg := range args {
			in[i] = reflect.ValueOf(arg)
			i++
		}
		// 捕获每个 handler 的异常，不能因为一个 handler 异常导致全部调用中断
		func() {
			defer RecoverAndLog()
			_ = f.Call(in) // f.CallSlice(in)
		}()
	}
}
