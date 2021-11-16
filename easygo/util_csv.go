//Package easygo ...
package easygo

type CsvData struct {
	data  map[string]interface{}
	mutex RLock
} //csv数据格式
func NewCsvData() *CsvData {
	p := &CsvData{}
	p.Init()
	return p
}
func (self *CsvData) Init() {
	self.data = make(map[string]interface{}, 0)
}
func (self *CsvData) Add(key string, val interface{}) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.data[key] = val
}
func (self *CsvData) Del(key string) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	delete(self.data, key)
}

func (self *CsvData) GetString(key string) string {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	val, ok := self.data[key]
	if ok {
		return AnytoA(val)
	}
	return ""

}
func (self *CsvData) GetInt(key string) int {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	val, ok := self.data[key].(int)
	if ok {
		return val
	}
	sVal := self.GetString(key)
	if sVal != "" {
		return Atoi(sVal)
	}
	return 0
}
func (self *CsvData) GetInt32(key string) int32 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	val, ok := self.data[key].(int32)
	if ok {
		return val
	}
	sVal := self.GetString(key)
	if sVal != "" {
		return AtoInt32(sVal)
	}
	return 0
}
func (self *CsvData) GetInt64(key string) int64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	val, ok := self.data[key].(int64)
	if ok {
		return val
	}
	sVal := self.GetString(key)
	if sVal != "" {
		return AtoInt64(sVal)
	}
	return 0
}
func (self *CsvData) GetFloat32(key string) float32 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	val, ok := self.data[key].(float32)
	if ok {
		return val
	}
	sVal := self.GetString(key)
	if sVal != "" {
		return AtoFloat32(sVal)
	}
	return 0
}
func (self *CsvData) GetFloat64(key string) float64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	val, ok := self.data[key].(float64)
	if ok {
		return val
	}
	sVal := self.GetString(key)
	if sVal != "" {
		return AtoFloat64(sVal)
	}
	return 0
}
