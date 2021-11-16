package easygo

import (
	"sync"
	"time"
	//"reflect"
	//"fmt"
)

type IQueue interface {
	Put(interface{})
	Take() interface{}
	Peek(block bool, timeout time.Duration) interface{}
	Empty() bool
	Full() bool
	Size() int
	Capacity() int
	Close()
}

type Queue struct {
	Me       IQueue
	mutex    Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	maxSize  int
	queue    []interface{}
}

func NewQueue(maxSize ...int) *Queue {
	//默认为 0, 表示无上限,撑爆内存为止
	size := append(maxSize, 0)[0]
	p := &Queue{}
	p.Init(p, size)
	return p
}
func (self *Queue) Init(me IQueue, maxSize ...int) {
	self.Me = me

	self.maxSize = append(maxSize, 0)[0]
	self.notFull = sync.NewCond(&self.mutex)
	self.notEmpty = sync.NewCond(&self.mutex)
}
func (self *Queue) Put(item interface{}) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.maxSize != 0 {
		for len(self.queue) >= self.maxSize {
			self.notFull.Wait()
		}
	}
	self.queue = append(self.queue, item)
	self.notEmpty.Signal()
}

func (self *Queue) Take() interface{} {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	for len(self.queue) <= 0 {
		self.notEmpty.Wait()
	}
	// assert(len(self.queue)>0);
	front := self.queue[0]
	self.queue = append(self.queue[:0], self.queue[1:]...)

	self.notFull.Signal()
	return front
}

func (self *Queue) Close() {
	self.Me.Put(nil)
}
func (self *Queue) Peek(block bool, timeout time.Duration) interface{} {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	for len(self.queue) <= 0 {
		self.notEmpty.Wait()
	}
	front := self.queue[0]
	return front
}

func (self *Queue) Empty() bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return len(self.queue) == 0
}

func (self *Queue) Full() bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return len(self.queue) >= self.maxSize
}

func (self *Queue) Size() int {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return len(self.queue)
}

func (self *Queue) Capacity() int {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return cap(self.queue)
}

/*
func main() {

	var qq []interface{}
	// qq = append(qq,4)
	fmt.Println(len(qq))

	queue := Queue{}
	queue.Init(4)

	// wg := sync.WaitGroup

	// producer
	go func() {
		for i := 0; i < 20; i++ {
			if i > 10 {
				time.Sleep(2 * time.Second)
			}
			queue.Put(i)
			fmt.Println("put item = ", i)
		}
	}()

	// consumer
	go func() {
		time.Sleep(2 * time.Second)
		for {
			item := queue.Take()
			fmt.Println("take item = ", item)

			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(60 * time.Second)

}
*/
