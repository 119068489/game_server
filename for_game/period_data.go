// 周期变量。到期自动重置
package for_game

import (
	"game_server/easygo"
	"log"
	"strconv"
)

var _ = log.Println

type IPeriodData interface {
	IPersistence

	GetPeriodNo() int
	CreatePeriodSeqFunctionSet() IPeriodSeqFunctionSet
}

type PeriodData struct {
	Persistence `bson:"-"`
	Me          IPeriodData                       `bson:"-"`
	KeepPeriod  int                               `bson:"-"`
	Data        map[string]map[string]interface{} `bson:"Data,omitempty"`
	FunctionSet IPeriodSeqFunctionSet             `bson:"-"`
}

/*
抽象类，不能实例化，只能被继承
func NewPeriodData(keepPeriod int, handlers ...func()) *PeriodData {
}*/

func (self *PeriodData) Init(me IPeriodData, keepPeriod int, kwargs ...easygo.KWAT) {
	if keepPeriod < 1 {
		panic("保存周期至少为 1")
	}
	self.KeepPeriod = keepPeriod
	self.Me = me
	self.Persistence.Init(me, kwargs...)

	self.Data = make(map[string]map[string]interface{})
	self.FunctionSet = self.Me.CreatePeriodSeqFunctionSet()
}

func (self *PeriodData) GetPeriodNo() int {
	panic("抽象方法，请在子类实现")
}

func (self *PeriodData) CreatePeriodSeqFunctionSet() IPeriodSeqFunctionSet {
	return PeriodSeq
}

func (self *PeriodData) OnLoad() { // override
	curPeriodNo := self.Me.GetPeriodNo()

	var dirty bool
	for periodNo, v := range self.Data {
		no := easygo.Atoi(periodNo)
		if curPeriodNo >= no+self.KeepPeriod {
			dirty = true
		} else {
			n := strconv.Itoa(no)
			self.Data[n] = v
		}
	}
	if dirty {
		self.Me.MarkDirty()
	}
}

func (self *PeriodData) Set(key string, value interface{}) { //
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	curNo := self.Me.GetPeriodNo()
	curPeriodNo := strconv.Itoa(curNo)
	data, ok := self.Data[curPeriodNo]
	if !ok {
		data = make(map[string]interface{})
		self.Data[curPeriodNo] = data
	}

	switch v := value.(type) { // 如果是数值，确保只存成 float64 和 int64
	case float32:
		data[key] = float64(v)
	case float64:
		data[key] = v

	case int:
		data[key] = int64(v)
	case int32:
		data[key] = int64(v)
	case int64:
		data[key] = v

	case uint:
		data[key] = int64(v)
	case uint32:
		data[key] = int64(v)
	case uint64:
		data[key] = int64(v)
	default:
		// 字符串，bool型， 复合结构之类的原样存下去
		data[key] = value
	}
	self.Me.MarkDirty()
}

func (self *PeriodData) Delete(key string) { //
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	curNo := self.Me.GetPeriodNo()
	curPeriodNo := strconv.Itoa(curNo)
	data, ok := self.Data[curPeriodNo]
	if !ok {
		return
	}
	delete(data, key)
	self.Me.MarkDirty()
}

func (self *PeriodData) AddInteger(key string, add interface{}, whichPeriod ...int) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	which := append(whichPeriod, 0)[0]
	curNo := self.Me.GetPeriodNo() + which
	curPeriodNo := strconv.Itoa(curNo)
	data, ok := self.Data[curPeriodNo]
	var old int64
	if ok {
		switch value := data[key].(type) {
		case nil:
			old = 0
		case int64:
			old = value
		case float64:
			old = int64(value)
		default:
			panic("为什么会存了其他类型在里面，哪里没有拦住")
		}
	} else {
		data = make(map[string]interface{})
		self.Data[curPeriodNo] = data
	}

	a := easygo.NewInt64(add)
	new := old + *a
	data[key] = new
	self.Me.MarkDirty()
}

func (self *PeriodData) AddFloat(key string, add interface{}, whichPeriod ...int) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	which := append(whichPeriod, 0)[0]
	curNo := self.Me.GetPeriodNo() + which
	curPeriodNo := strconv.Itoa(curNo)
	data, ok := self.Data[curPeriodNo]
	var old float64
	if ok {
		switch value := data[key].(type) {
		case nil:
			old = float64(0)
		case int64:
			old = float64(value)
		case float64:
			old = value
		default:
			panic("为什么会存了其他类型在里面，哪里没有拦住")
		}
	} else {
		data = make(map[string]interface{})
		self.Data[curPeriodNo] = data
	}

	a := easygo.NewFloat64(add)
	new := old + *a
	data[key] = new
	self.Me.MarkDirty()
}

func (self *PeriodData) Fetch(key string, whichPeriod ...int) interface{} { // 转不实现 "找不到时返回指定值"
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	return self.FetchNoLock(key, whichPeriod...)
}

func (self *PeriodData) FetchNoLock(key string, whichPeriod ...int) interface{} { // 转不实现 "找不到时返回指定值"
	which := append(whichPeriod, 0)[0]
	no := self.Me.GetPeriodNo() + which
	n := strconv.Itoa(no)
	data, ok := self.Data[n]
	if !ok {
		return nil
	}
	return data[key]
}

func (self *PeriodData) FetchString(key string, whichPeriod ...int) string {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value := self.FetchNoLock(key, whichPeriod...)
	if value == nil {
		return ""
	}
	return value.(string) // 如果实际存储的不是 string ，那是调用者的责任
}
func (self *PeriodData) FetchBool(key string, whichPeriod ...int) bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value := self.FetchNoLock(key, whichPeriod...)
	if value == nil {
		return false
	}
	return value.(bool) // 如果实际存储的不是 string ，那是调用者的责任
}

func (self *PeriodData) FetchInt(key string, whichPeriod ...int) int { // 找不到返回 0
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value := self.FetchNoLock(key, whichPeriod...)
	if value == nil {
		return 0
	}
	a := easygo.NewInt(value)
	return *a
}

func (self *PeriodData) FetchInt32(key string, whichPeriod ...int) int32 { // 找不到返回 0
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value := self.FetchNoLock(key, whichPeriod...)
	if value == nil {
		return 0
	}
	a := easygo.NewInt32(value)
	return *a
}

func (self *PeriodData) FetchInt64(key string, whichPeriod ...int) int64 { // 找不到返回 0
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value := self.FetchNoLock(key, whichPeriod...)
	if value == nil {
		return 0
	}
	a := easygo.NewInt64(value)
	return *a
}

func (self *PeriodData) FetchFloat64(key string, whichPeriod ...int) float64 { // 找不到返回 0
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value := self.FetchNoLock(key, whichPeriod...)
	if value == nil {
		return 0
	}
	a := easygo.NewFloat64(value)
	return *a
}

func (self *PeriodData) FetchByte(key string, whichPeriod ...int) []byte { // 找不到返回 []byte
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	value := self.FetchNoLock(key, whichPeriod...)
	if value == nil {
		return nil
	}
	a := value.([]byte)
	return a
}

//-------------------------------
// 天变量
type DayPeriodData struct {
	PeriodData `bson:",inline,omitempty"`
}

func NewDayPeriodData(keepPeriod int, kwargs ...easygo.KWAT) *DayPeriodData {
	p := &DayPeriodData{}
	p.Init(keepPeriod, kwargs...)
	return p
}

func (self *DayPeriodData) Init(keepPeriod int, kwargs ...easygo.KWAT) {
	self.PeriodData.Init(self, keepPeriod, kwargs...)
}

func (self *DayPeriodData) GetPeriodNo() int { // override
	return self.FunctionSet.GetDayNo()
}

//-------------------------------

// 周变量
type WeekPeriodData struct {
	PeriodData `bson:",inline,omitempty"`
}

func NewWeekPeriodData(keepPeriod int, kwargs ...easygo.KWAT) *WeekPeriodData {
	p := &WeekPeriodData{}
	p.Init(keepPeriod, kwargs...)
	return p
}

func (self *WeekPeriodData) Init(keepPeriod int, kwargs ...easygo.KWAT) {
	self.PeriodData.Init(self, keepPeriod, kwargs...)
}

func (self *WeekPeriodData) GetPeriodNo() int { // override
	return self.FunctionSet.GetWeekNo()
}

//-------------------------------

// 月变量
type MonthPeriodData struct {
	PeriodData `bson:",inline,omitempty"`
}

func NewMonthPeriodData(keepPeriod int, kwargs ...easygo.KWAT) *MonthPeriodData {
	p := &MonthPeriodData{}
	p.Init(keepPeriod, kwargs...)
	return p
}

func (self *MonthPeriodData) Init(keepPeriod int, kwargs ...easygo.KWAT) {
	self.PeriodData.Init(self, keepPeriod, kwargs...)
}

func (self *MonthPeriodData) GetPeriodNo() int { // override
	return self.FunctionSet.GetMonthNo()
}

//-------------------------------

// 半年变量
type HalfYearPeriodData struct {
	PeriodData `bson:",inline,omitempty"`
}

func NewHaltYearPeriodData(keepPeriod int, kwargs ...easygo.KWAT) *HalfYearPeriodData {
	p := &HalfYearPeriodData{}
	p.Init(keepPeriod, kwargs...)
	return p
}

func (self *HalfYearPeriodData) Init(keepPeriod int, kwargs ...easygo.KWAT) {
	self.PeriodData.Init(self, keepPeriod, kwargs...)
}

func (self *HalfYearPeriodData) GetPeriodNo() int { // override
	return self.FunctionSet.GetHalfYearNo()
}
