// JSON 对象持久化
package for_game

import (
	"game_server/easygo"
	"log"
)

var _ = log.Println

type IEasyPersist interface {
	IPersistence
}

type EasyPersist struct {
	Persistence `bson:"-"`
	Me          IEasyPersist `bson:"-"`

	Data map[string]interface{} `bson:"Data,omitempty"`
}

func NewEasyPersist(kwargs ...easygo.KWAT) *EasyPersist {
	p := &EasyPersist{}
	p.Init(p, kwargs...)
	return p
}

func (self *EasyPersist) Init(me IEasyPersist, kwargs ...easygo.KWAT) {
	self.Me = me
	self.Persistence.Init(me, kwargs...)
	self.Data = make(map[string]interface{})
}

func (self *EasyPersist) Set(key string, value interface{}) { //
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	switch v := value.(type) { // 如果是数值，确保只存成 float64 和 int64
	case float32:
		self.Data[key] = float64(v)
	case float64:
		self.Data[key] = v

	case int:
		self.Data[key] = int64(v)
	case int32:
		self.Data[key] = int64(v)
	case int64:
		self.Data[key] = v

	case uint:
		self.Data[key] = int64(v)
	case uint32:
		self.Data[key] = int64(v)
	case uint64:
		self.Data[key] = int64(v)
	default:
		// 字符串，bool型， 复合结构之类的原样存下去
		self.Data[key] = value
	}
	self.Me.MarkDirty()
}

func (self *EasyPersist) Delete(key string) { //
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	delete(self.Data, key)
	self.Me.MarkDirty()
}

func (self *EasyPersist) AddInteger(key string, add interface{}, defaultVal ...int) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	var old int64
	v, ok := self.Data[key]
	if ok {
		switch value := v.(type) {
		case int64:
			old = value
		case float64:
			old = int64(value)
		default:
			panic("为什么会存了其他类型在里面，哪里没有拦住")
		}
	}

	var a int64
	switch value := add.(type) {

	case int:
		a = int64(value)
	case int32:
		a = int64(value)
	case int64:
		a = value
	case uint:
		a = int64(value)
	case uint32:
		a = int64(value)
	case uint64:
		a = int64(value)
	default:
		panic("不支持的类型,自行补充")
	}

	new := old + a
	self.Data[key] = new
	self.Me.MarkDirty()
}

func (self *EasyPersist) AddFloat(key string, add interface{}, defaultVal ...float64) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	var old float64
	v, ok := self.Data[key]
	if ok {
		switch value := v.(type) {
		case float64:
			old = value
		case int64:
			old = float64(value)
		default:
			panic("为什么会存了其他类型在里面，哪里没有拦住")
		}
	}

	var a float64
	switch value := add.(type) {
	case float32:
		old = float64(value)
	case float64:
		old = value
	case int64:
		a = float64(value)
	case int:
		a = float64(value)
	case int32:
		a = float64(value)
	case uint32:
		a = float64(value)
	case uint:
		a = float64(value)
	default:
		panic("不支持的类型,自行补充")
	}

	new := old + a
	self.Data[key] = new
	self.Me.MarkDirty()
}

func (self *EasyPersist) Fetch(key string) interface{} {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	return self.Data[key]
}

func (self *EasyPersist) FetchString(key string) string {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value, ok := self.Data[key]
	if !ok {
		return ""
	}
	v := value.(string) // 如果实际存储的不是 string ，那是调用者的责任
	return v
}

func (self *EasyPersist) FetchInt(key string, defVal ...int) int { // 找不到时，且你没有拱供默认值则返回 0
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value, ok := self.Data[key]
	if !ok {
		return append(defVal, 0)[0]
	}
	switch v := value.(type) { // AddInteger AddFloat 保证了实际存储的一定是 int64 或 float64
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		panic("为什么会存了其他类型在里面，哪里没有拦住")
	}
}

func (self *EasyPersist) FetchInt32(key string, defVal ...int32) int32 { // 找不到时，且你没有拱供默认值则返回 0
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value, ok := self.Data[key]
	if !ok {
		return append(defVal, 0)[0]
	}
	switch v := value.(type) { // AddInteger AddFloat 保证了实际存储的一定是 int64 或 float64
	case int64:
		return int32(v)
	case float64:
		return int32(v)
	default:
		panic("为什么会存了其他类型在里面，哪里没有拦住")
	}
}

func (self *EasyPersist) FetchInt64(key string, defVal ...int64) int64 { // 找不到时，且你没有拱供默认值则返回 0
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value, ok := self.Data[key]
	if !ok {
		return append(defVal, 0)[0]
	}
	switch v := value.(type) { // AddInteger AddFloat 保证了实际存储的一定是 int64 或 float64
	case int64:
		return v
	case float64:
		return int64(v)
	default:
		panic("为什么会存了其他类型在里面，哪里没有拦住")
	}
}

func (self *EasyPersist) FetchFloat32(key string, defVal ...float32) float32 { // 找不到时，且你没有拱供默认值则返回 0
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value, ok := self.Data[key]
	if !ok {
		return append(defVal, 0)[0]
	}
	switch v := value.(type) { // AddInteger AddFloat 保证了实际存储的一定是 int64 或 float64
	case int64:
		return float32(v)
	case float64:
		return float32(v)
	default:
		panic("为什么会存了其他类型在里面，哪里没有拦住")
	}
}

func (self *EasyPersist) FetchFloat64(key string, defVal ...float64) float64 { // 找不到时，且你没有拱供默认值则返回 0
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value, ok := self.Data[key]
	if !ok {
		return append(defVal, 0)[0]
	}
	switch v := value.(type) { // AddInteger AddFloat 保证了实际存储的一定是 int64 或 float64
	case int64:
		return float64(v)
	case float64:
		return v
	default:
		panic("为什么会存了其他类型在里面，哪里没有拦住")
	}
}
