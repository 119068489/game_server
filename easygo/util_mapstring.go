//Package easygo ...
package easygo

type MyStrMap struct {
	Data map[string]interface{}
}

func (self MyStrMap) Add(key string, val interface{}) {
	self.Data[key] = val
}
func (self MyStrMap) Del(key string) {
	delete(self.Data, key)
}
func (self MyStrMap) GetInt64(key string) int64 {
	val, ok := self.Data[key]
	if !ok {
		return 0
	}
	return AtoInt64(AnytoA(val))
}

func (self MyStrMap) GetFloat64(key string) float64 {
	val, ok := self.Data[key]
	if !ok {
		return 0
	}
	return AtoFloat64(AnytoA(val))
}
func (self MyStrMap) GetString(key string) string {
	val, ok := self.Data[key]
	if !ok {
		return ""
	}
	return AnytoA(val)
}
