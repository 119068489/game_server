package util

import (
	"math/rand"
	"strconv"
	"time"
)

// GenRangeInt 创建一个长度为length的slice
// 值为连续的整形，第一个数为from
func GenRangeInt(length int, from int) []int {
	_range := make([]int, length, length)
	for i := 0; i < length; i++ {
		_range[i] = i + from
	}
	return _range
}

// IntInSlice 判断某个int值是否在切片中
func IntInSlice(finder int, slice []int) bool {
	exists := false
	for _, v := range slice {
		if v == finder {
			exists = true
			break
		}
	}
	return exists
}

// Int32InSlice 判断某个int值是否在切片中
func Int32InSlice(finder int32, slice []int32) bool {
	exists := false
	for _, v := range slice {
		if v == finder {
			exists = true
			break
		}
	}
	return exists
}

// Int64InSlice 判断某个int64值是否在切片中
func Int64InSlice(finder int64, slice []int64) bool {
	exists := false
	for _, v := range slice {
		if v == finder {
			exists = true
			break
		}
	}
	return exists
}

// ShuffleSliceInt 打乱一个切片
func ShuffleSliceInt(src []int) []int {
	dest := make([]int, len(src))

	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		dest[v] = src[i]
	}

	return dest
}

// IsSameSlice 判断两个slice是否相等
func IsSameSlice(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

// SliceDel 删除slice中的某些元素
func SliceDel(slice []int, values ...int) []int {
	if slice == nil || len(values) == 0 {
		return slice
	}
	for _, value := range values {
		slice = SliceDelOne(slice, value)
	}
	return slice
}

func SliceDelOne(slice []int, value int) []int {
	if slice == nil {
		return slice
	}
	for i, j := range slice {
		if j == value {
			return append(append([]int{}, slice[:i]...), slice[i+1:]...)
		}
	}
	return slice
}

// SliceDel 删除slice中的某些元素
func SliceDelInt32(slice []int32, values ...int32) []int32 {
	if slice == nil || len(values) == 0 {
		return slice
	}
	for _, value := range values {
		slice = SliceDelOneInt32(slice, value)
	}
	return slice
}

func SliceDelOneInt32(slice []int32, value int32) []int32 {
	if slice == nil {
		return slice
	}
	for i, j := range slice {
		if j == value {
			return append(append([]int32{}, slice[:i]...), slice[i+1:]...)
		}
	}
	return slice
}

// SliceDel 删除slice中的某些元素
func SliceDelInt64(slice []int64, values ...int64) []int64 {
	if slice == nil || len(values) == 0 {
		return slice
	}
	for _, value := range values {
		slice = SliceDelOneInt64(slice, value)
	}
	return slice
}

func SliceDelOneInt64(slice []int64, value int64) []int64 {
	if slice == nil {
		return slice
	}
	for i, j := range slice {
		if j == value {
			return append(append([]int64{}, slice[:i]...), slice[i+1:]...)
		}
	}
	return slice
}

// SliceCopy 拷贝一个切片
func SliceCopy(s []int) []int {
	var slice = make([]int, len(s))
	copy(slice, s)
	return slice
}

func SliceCopyInt64(s []int64) []int64 {
	var slice = make([]int64, len(s))
	copy(slice, s)
	return slice
}

func SliceCopyString(s []string) []string {
	var slice = make([]string, len(s))
	copy(slice, s)
	return slice
}

// SliceJoin 将一个slice拼接成一个字符串
func SliceJoin(s []int, joinString string) string {
	var str = ""
	if length := len(s); length > 0 {
		str = strconv.Itoa(s[0])
		for i := 1; i < length; i++ {
			str += joinString
			str += strconv.Itoa(s[i])
		}
	}
	return str
}

// InStringSlice 判断某个string值是否在切片中
func InStringSlice(finder string, slice []string) bool {
	exists := false
	for _, v := range slice {
		if v == finder {
			exists = true
			break
		}
	}
	return exists
}

// SliceDelString 删除slice中的某些元素
func SliceDelString(slice []string, values ...string) []string {
	if slice == nil || len(values) == 0 {
		return slice
	}
	for _, value := range values {
		slice = sliceDelString(slice, value)
	}
	return slice
}

func sliceDelString(slice []string, value string) []string {
	if slice == nil {
		return slice
	}
	for i, j := range slice {
		if j == value {
			// return append(slice[:i], slice[i+1:]...)
			return append(append([]string{}, slice[:i]...), slice[i+1:]...)
		}
	}
	return slice
}

// SliceMaxInt 取int类型的最大值
func SliceMaxInt(s []int) int {
	var max = 0
	for _, v := range s {
		if v > max {
			max = v
		}
	}
	return max
}

// SliceMaxInt 取int类型的最大值
func Slice64MaxInt(s []int64) int64 {
	var max int64
	for _, v := range s {
		if v > max {
			max = v
		}
	}
	return max
}

// SliceUniqueInt 去重
func SliceUniqueInt(s []int) []int {
	uniquedSlice := []int{}
	m := make(map[int]bool)
	for _, v := range s {
		if _, exists := m[v]; !exists {
			m[v] = true
			uniquedSlice = append(uniquedSlice, v)
		}
	}
	return uniquedSlice
}

// SliceToMap 将[]int 转化成map[int]count
func SliceToMap(slice []int) map[int]int {
	var m = map[int]int{}
	for _, j := range slice {
		var _, ok = m[j]
		if ok {
			m[j]++
		} else {
			m[j] = 1
		}
	}
	return m
}

// MapToSlice 将map[int]count 转成 []int
func MapToSlice(m map[int]int) []int {
	tiles := []int{}
	for tile, cnt := range m {
		for i := 0; i < cnt; i++ {
			tiles = append(tiles, tile)
		}
	}
	return tiles
}

// 从数组1中去重数组2的内容
func Slice1DelSlice2(slice1 []int64, slice2 []int64) []int64 {
	for _, v1 := range slice2 {
		for k2, v2 := range slice1 {
			if v1 == v2 {
				// 删除对应的k2
				slice1 = append(slice1[:k2], slice1[k2+1:]...)
			}
		}
	}
	return slice1
}

// 从数组中随机抽取数组
func RandSliceFromSlice(ops []string, num int) []string {
	result := make([]string, 0)
	if len(ops) == 0 {
		return result
	}

	if num >= len(ops) {
		return ops
	}
	for {
		index := rand.Intn(len(ops) - 1)
		if InStringSlice(ops[index], result) {
			continue
		}
		result = append(result, ops[index])
		if len(result) == num {
			break
		}
	}
	return result
}
