package easygo

import (
	"runtime"
	"sync/atomic"

	sync "github.com/sasha-s/go-deadlock"
)

//-------------------------------

type Mutex struct {
	sync.Mutex
}

type RWMutex struct {
	sync.RWMutex
}

//-------------------------------

// 可重入锁
type RLock struct {
	Mutex   Mutex
	OwnerId int32 // 全程用 atomic 相关函数操作
	Count   int32 // 始终受 Mutex 保护
}

func (self *RLock) Lock() {
	currentId := GetGoroutineId()
	if atomic.LoadInt32(&self.OwnerId) == currentId {
		self.Count++
		return
	}
	self.Mutex.Lock()
	atomic.StoreInt32(&self.OwnerId, currentId)
	self.Count = 1 // 这一行已经受 self.Mutex.Lock() 保护了，可以直接赋值

}
func (self *RLock) Unlock() {
	currentId := GetGoroutineId()
	if atomic.LoadInt32(&self.OwnerId) != currentId {
		panic("Lock 和 Unlock 必须在同一个协程")
	}
	// 来到这里，肯定已经受到 Lock 保护了
	self.Count--
	if self.Count == 0 {
		atomic.StoreInt32(&self.OwnerId, 0)
		self.Mutex.Unlock()
	} else if self.Count < 0 {
		panic("怎么可能小于 0")
	}
}

func GetGoroutineId() int32 {
	var buf [64]byte
	runtime.Stack(buf[:], false)

	var id int32
	for i := 10; ; i++ { // 字符串是以 "goroutine " 开头的，长度为 10 个 byte
		c := buf[i]
		if c == ' ' {
			break
		}
		v := int32(c - '0')
		id = id*10 + v
	}
	return id
}

//------------------------------------------------------------------
/*
func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			result := sum(10000)
			log.Printf("第 %d 个任务，结果 = %d ", i, result)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

var rlock = RLock{}

// 递归累加
func sum(i int) int {
	// rlock.Lock()
	// defer rlock.Unlock()

	if i == 1 {
		return 1
	}
	return i + sum(i-1)
}
*/
