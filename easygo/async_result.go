package easygo

import (
	"time"
)

//-----------------------------
type IAsyncResult interface {
	Value() interface{}
	Successful() bool

	Set(value interface{})

	Get() (interface{}, ITimeoutError)
	GetUntilTimeout(timeout time.Duration) (interface{}, ITimeoutError)
	GetNonBlock() (interface{}, ITimeoutError)
	GetImplement(block bool, timeout time.Duration) (interface{}, ITimeoutError)
}

//-----------------------------

type AsyncResult struct {
	Me       IAsyncResult
	mutex    Mutex
	setEvent chan bool
	value    interface{}
}

func NewAsyncResult() *AsyncResult {
	p := &AsyncResult{}
	p.Init(p)
	return p
}
func (self *AsyncResult) Init(me IAsyncResult) {
	self.Me = me
	self.setEvent = make(chan bool)
	self.value = None
}
func (self *AsyncResult) Set(value interface{}) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.value == None {
		self.value = value
		close(self.setEvent) // 重复 close 会 panic
	}
}
func (self *AsyncResult) Value() interface{} {
	if self.value == None {
		return nil
	} else {
		return self.value
	}
}
func (self *AsyncResult) Successful() bool {
	return self.value != None
}

/*block forever*/
func (self *AsyncResult) Get() (interface{}, ITimeoutError) {
	return self.Me.GetImplement(true, 0)
}

/*block to timeout*/
func (self *AsyncResult) GetUntilTimeout(timeout time.Duration) (interface{}, ITimeoutError) {
	return self.Me.GetImplement(true, timeout)
}

/*non-block,如果没有 set 返回 TimeoutError */
func (self *AsyncResult) GetNonBlock() (interface{}, ITimeoutError) {
	return self.Me.GetImplement(false, 0)
}

/*
具体实现，timeout 为0表示不超时
*/
func (self *AsyncResult) GetImplement(block bool, timeout time.Duration) (interface{}, ITimeoutError) {
	var tc <-chan time.Time = nil
	var timer *time.Timer = nil

	self.mutex.Lock()
	if self.value != None {
		value := self.value
		self.mutex.Unlock()
		return value, nil
	}
	if !block {
		self.mutex.Unlock()
		return nil, NewTimeoutError("timeout") // todo 把具体时间放进 NewTimeoutError 对象里面
	}
	if timeout != 0 {
		timer = time.NewTimer(timeout)
		tc = timer.C
	}
	self.mutex.Unlock()
	select {
	case <-self.setEvent:
		if timer != nil {
			timer.Stop()
		}
		self.mutex.Lock() // 取 self.value 要上锁
		defer self.mutex.Unlock()
		return self.value, nil
	case <-tc: // 如果是 nil chan 会永久阻塞
		return nil, NewTimeoutError("timeout") // todo 把具体时间放进 error 对象里面
	}
}
