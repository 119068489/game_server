package util

import (
	"sync"
)

// SMapLen 获取sync.Map的元素个数
func SMapLen(smap *sync.Map) int {
	length := 0
	smap.Range(func(k, v interface{}) bool {
		length++
		return true
	})
	return length
}

// SMapIsEmpty 判断sync.Map是否为空
func SMapIsEmpty(smap *sync.Map) bool {
	empty := true
	smap.Range(func(k, v interface{}) bool {
		empty = false
		return false
	})
	return empty
}
