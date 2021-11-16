package easygo

type IMapping interface {
}

// 线程安全，且是不实际存储对象
type Mapping struct {
	Me    IMapping
	Data  map[interface{}]interface{}
	Mutex RLock
}

func (self *Mapping) Init(me IMapping) {
	self.Me = me
	self.Data = make(map[interface{}]interface{})
}
func (self *Mapping) Store(key interface{}, fetch func() interface{}) bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value := fetch()
	if value == nil {
		// panic("要 Store 的对象找不到")
		return false
	}
	finalable, ok := value.(IFinalable)
	if !ok {
		panic("要 Store 的对象必须是 IFinalable")
	}
	self.Data[key] = fetch

	ft := func(final IFinalable) {
		self.Mutex.Lock()
		defer self.Mutex.Unlock()

		v, ok := self.Data[key]
		if ok {
			fetch := v.(func() interface{})
			if fetch() == final {
				delete(self.Data, key)
			} //else {
			// 	log.Println("幸亏检测到了，不然又是一次误删")
			// }
		}
	}
	_ = finalable
	_ = ft
	// finalable.AddFinalizer(ft) // finalable 析构时自动移除 fetch 函数
	return true
}

// 到实际存储的地方取对象
func (self *Mapping) Load(key interface{}) interface{} {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value, ok := self.Data[key]
	if !ok {
		return nil
	}
	fetch := value.(func() interface{})
	return fetch()
}

func (self *Mapping) GetCopy() map[interface{}]interface{} {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	length := len(self.Data)
	m := make(map[interface{}]interface{}, length)
	for k, v := range self.Data {
		m[k] = v
	}
	return m
}

func (self *Mapping) Delete(key interface{}) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	delete(self.Data, key)
}

func (self *Mapping) Length() int {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	return len(self.Data)
}
