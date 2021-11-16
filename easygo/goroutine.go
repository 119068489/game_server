package easygo

import (
	"reflect"

	"time"
)

type IGoroutine interface {
	Start(function interface{}, args ...interface{})

	// StartLater(seconds time.Duration, function interface{}, args ...interface{})
	Join(timeout ...time.Duration)
	Get(block bool, timeout ...time.Duration) interface{}
	// Successful() bool
	// Ready()
	// Exception()

	GetCompleteEvent() chan bool
	GetValue() interface{}
	OnPanic(callStack string, recoverVal interface{})
}
type Goroutine struct {
	Me       IGoroutine
	function reflect.Value
	args     []interface{}

	mutex         Mutex
	CompleteEvent chan bool
	Value         interface{} // 要确保 value 线程安全
	started       bool
}

func NewGoroutine() *Goroutine {
	p := &Goroutine{}
	p.Init(p)
	return p
}

func (self *Goroutine) Init(me IGoroutine) {
	self.Me = me
	self.CompleteEvent = make(chan bool)
	self.Value = None
}

func (self *Goroutine) Start(function interface{}, args ...interface{}) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.started {
		return
	}
	self.started = true

	if function == nil {
		panic("function 参数是个 nil")
	}
	f := reflect.ValueOf(function)
	if f.Kind() != reflect.Func {
		panic("function 参数必须是个 function")
	}

	go func() {
		var value interface{}
		defer func() {
			r := recover() // recover 必须在 defer 函数中直接调用才能拦截异常，不能间接调用(试验出来的结果)
			if r != nil {
				callStack := LogPanicAndStack(r)
				value = r
				self.Me.OnPanic(callStack, r)
			}
			self.mutex.Lock()
			defer self.mutex.Unlock()
			self.Value = value
			close(self.CompleteEvent)
		}()

		in := make([]reflect.Value, len(args))
		i := int32(0)
		for _, arg := range args {
			in[i] = reflect.ValueOf(arg)
			i++
		}

		values := f.Call(in)
		if len(values) == 0 {
			value = nil
		} else if len(values) == 1 {
			value = values[0].Interface()
		} else {
			list := make([]interface{}, 0, len(values))
			for _, v := range values {
				list = append(list, v.Interface())
			}
			value = list
		}
	}()
}

func (self *Goroutine) OnPanic(callStack string, recoverVal interface{}) {

}

func (self *Goroutine) Join(timeout ...time.Duration) {
	self.Me.Get(true, timeout...)
}

// 如果超时 ，则返回 TimeoutError。有异常则向上波及。否则返回正常执行的结果
func (self *Goroutine) Get(block bool, timeout ...time.Duration) interface{} {
	to := append(timeout, 0)[0]

	self.mutex.Lock()
	value := self.Value
	self.mutex.Unlock()

	if value != None {
		if _, ok := value.(error); ok { // 各种下标越界，除 0 异常都是符合 error 接口
			panic(value) // 向上波及
		} else { // 正常的执行调用得到的结果
			return value
		}
	}
	if !block {
		return NewTimeoutError("timeout") // todo 把具体时间放进 error 对象里面
	}

	var tc <-chan time.Time = nil
	var timer *time.Timer = nil
	if to != 0 {
		timer = time.NewTimer(to)
		tc = timer.C
	}

	select {
	case <-self.CompleteEvent:
		if timer != nil {
			timer.Stop()
		}
		self.mutex.Lock() // 取 self.Value 要上锁
		defer self.mutex.Unlock()
		return self.Value

	case <-tc: // 如果 tc 是 nil chan 会永久阻塞,也就没有机会进这个分支
		return NewTimeoutError("timeout") // todo 把具体时间放进 error 对象里面
	}
}

// func (self *Goroutine) Successful() bool {
// 	if self.Value == None {
// 		return false
// 	}else if _, ok := self.Value.(error); ok { //
// 		return false
// 	}
// 	return true
// }

func (self *Goroutine) GetCompleteEvent() chan bool {
	return self.CompleteEvent
}

func (self *Goroutine) GetValue() interface{} {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.Value
}

//-------------------------------------------------

func JoinAllJobs(jobs ...IGoroutine) { // 不支持 timeout
	if len(jobs) <= 0 {
		return
	}
	JoinJobsImplement(jobs, true, len(jobs))
}

func JoinSomeJobs(count int, jobs ...IGoroutine) { // 不支持 timeout
	if len(jobs) <= 0 {
		return
	}
	JoinJobsImplement(jobs, true, count)
}

func JoinAllNotPanic(jobs ...IGoroutine) { // 不支持 timeout
	if len(jobs) <= 0 {
		return
	}
	JoinJobsImplement(jobs, false, len(jobs))
}

func JoinSomeNotPanic(count int, jobs ...IGoroutine) { // 不支持 timeout
	if len(jobs) <= 0 {
		return
	}
	JoinJobsImplement(jobs, false, count)
}

func JoinJobsImplement(jobs []IGoroutine, panicError bool, count int) { // 不支持 timeout
	cases := make([]reflect.SelectCase, len(jobs))
	for i, job := range jobs {
		event := job.GetCompleteEvent()
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(event)}
	}
	for i := 0; i < count; i++ {
		chosen, recv, recvOK := reflect.Select(cases)
		// ok will be true if the channel has not been closed.
		_, _ = recv, recvOK
		job := jobs[chosen]

		cases = append(cases[:chosen], cases[chosen+1:]...)
		jobs = append(jobs[:chosen], jobs[chosen+1:]...)
		if panicError {
			value := job.GetValue()
			if _, ok := value.(error); ok { // to think
				panic(value)
			}
		}
	}
}

func Spawn(function interface{}, args ...interface{}) IGoroutine {
	job := NewGoroutine()
	job.Start(function, args...)
	return job
}
