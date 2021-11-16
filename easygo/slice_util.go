// 我想实现 slice 的通用算法
package easygo

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
)

// 网上抄来的
func Insert(slice interface{}, pos int, value interface{}) interface{} {
	v := reflect.ValueOf(slice)
	v = reflect.Append(v, reflect.ValueOf(value))
	reflect.Copy(v.Slice(pos+1, v.Len()), v.Slice(pos, v.Len()))
	v.Index(pos).Set(reflect.ValueOf(value))
	return v.Interface()
}

// 判断 item 是否在 container 中，container 支持的类型 arrary,slice,map。
// 用了反射，估计性能会差一些(网上抄的)
func Contain(container interface{}, item interface{}) bool {
	containerValue := reflect.ValueOf(container)
	switch reflect.TypeOf(container).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerValue.Len(); i++ {
			if containerValue.Index(i).Interface() == item {
				return true
			}
		}
	case reflect.Map:
		if containerValue.MapIndex(reflect.ValueOf(item)).IsValid() {
			return true
		}
	}
	return false
}

func In(slice interface{}, item interface{}) bool {
	if list, ok := slice.([]int); ok {
		for _, v := range list {
			if v == item {
				return true
			}
		}
		return false
	} else if list, ok := slice.([]int32); ok {
		for _, v := range list {
			if v == item {
				return true
			}
		}
		return false
	} else if list, ok := slice.([]int64); ok {
		for _, v := range list {
			if v == item {
				return true
			}
		}
		return false
	} else {
		panic("不支持的切片类型，你要不要完善一下代码")
	}
}

func Del(slice interface{}, item interface{}) interface{} {
	index := Index(slice, item)
	if index == -1 {
		return slice
	}
	if list, ok := slice.([]int); ok {
		return append(list[:index], list[index+1:]...)
	} else if list, ok := slice.([]int32); ok {
		return append(list[:index], list[index+1:]...)
	} else if list, ok := slice.([]int64); ok {
		return append(list[:index], list[index+1:]...)
	} else if list, ok := slice.([]float32); ok {
		return append(list[:index], list[index+1:]...)
	} else if list, ok := slice.([]string); ok {
		return append(list[:index], list[index+1:]...)
	}
	panic("不支持的数据类型:")
}

// 查找指定元素所在的索引。-1 表示没有找到
func Index(container interface{}, item interface{}) int {
	k := reflect.TypeOf(container).Kind()
	if k != reflect.Slice && k != reflect.Array {
		panic("仅支持 slice 和 array")
	}

	containerValue := reflect.ValueOf(container)
	for i := 0; i < containerValue.Len(); i++ {
		tem := containerValue.Index(i).Interface()
		if tem == item {
			return i
		}
	}
	return -1
}

/*
func Index(slice interface{}, item interface{}) int {
	if list, ok := slice.([]int); ok {
		for i, v := range list {
			if v == item {
				return i
			}
		}
		return -1
	} else if list, ok := slice.([]int32); ok {
		for i, v := range list {
			if v == item {
				return i
			}
		}
		return -1
	} else if list, ok := slice.([]int64); ok {
		for i, v := range list {
			if v == item {
				return i
			}
		}
		return -1
	} else {
		panic("不支持的切片类型，你要不要完善一下代码")
	}
}*/

func ToInterfaceArr(arr interface{}) []interface{} {
	if reflect.TypeOf(arr).Kind() != reflect.Slice {
		return nil
	}
	arrValue := reflect.ValueOf(arr)
	retArr := make([]interface{}, arrValue.Len())
	for k := 0; k < arrValue.Len(); k++ {
		retArr[k] = arrValue.Index(k).Interface()
	}
	return retArr
}
func RandGetNItemFromSlice(silce interface{}, n int) []interface{} {
	interfaceSlice := ToInterfaceArr(silce)

	var result []interface{}
	for i := 0; i < n; i++ {
		interfaceSliceLen := len(interfaceSlice)
		intn := rand.Intn(interfaceSliceLen)
		result = append(result, interfaceSlice[intn])

		var left []interface{}
		for lIndex, lInterface := range interfaceSlice {
			if lIndex != intn {
				left = append(left, lInterface)
			}
		}
		interfaceSlice = left
	}

	return result
}

type bodyWrapper struct {
	Bodys []interface{}
	by    func(p, q *interface{}) bool
}

func (self bodyWrapper) Len() int {
	return len(self.Bodys)
}
func (self bodyWrapper) Swap(i, j int) {
	self.Bodys[i], self.Bodys[j] = self.Bodys[j], self.Bodys[i]
}
func (self bodyWrapper) Less(i, j int) bool {
	return self.by(&self.Bodys[i], &self.Bodys[j])
}
func SortStructSlice(structSlice interface{}, sortByField string, isAscOrDesc bool) {
	sort.Sort(bodyWrapper{ToInterfaceArr(structSlice), func(p, q *interface{}) bool {
		v := reflect.ValueOf(*p)
		i := v.FieldByName(sortByField)
		v = reflect.ValueOf(*q)
		j := v.FieldByName(sortByField)
		if isAscOrDesc {
			return i.String() < j.String()
		} else {
			return i.String() > j.String()
		}
	}})
}

func SortSliceInt32(Slice []int32, isAscOrDesc bool) {
	sort.Slice(Slice, func(i, j int) bool {
		if isAscOrDesc {
			return Slice[i] < Slice[j] // 生序
		} else {
			return Slice[i] > Slice[j] // 降序
		}

	})
}

func SortSliceInt64(Slice []int64, isAscOrDesc bool) {
	sort.Slice(Slice, func(i, j int) bool {
		if isAscOrDesc {
			return Slice[i] < Slice[j] // 生序
		} else {
			return Slice[i] > Slice[j] // 降序
		}

	})
}

// 倒序排序数组,如[23,44,12,45]-->[45,12,44,23]
func DescInt64Slice(s []int64) []int64 {
	ns := make([]int64, 0)
	if len(s) == 0 {
		return ns
	}
	for i := len(s) - 1; i >= 0; i-- {
		ns = append(ns, s[i])
	}
	return ns
}

//通过指定开始和结束值，返回数字数组[start,end]
func GetInt64List(start, end int64) []int64 {
	val := make([]int64, 0)
	for i := start; i <= end; i++ {
		val = append(val, i)
	}
	return val
}

/**
排列组合 源{1,2,3,4,5}
返回{[1 2], [2 1], [1 3] ,[3 1], [1 4] ,[4 1], [1 5] ,[5 1] ,[2 3], [3 2] ,[2 4] ,[4 2] ,[2 5] ,[5 2] ,[3 4] ,[4 3], [3 5], [5 3], [4 5] ,[5 4]}
*/
func SlicePermutations(pa []int32) [][]int32 {
	result := [][]int32{}
	if len(pa) == 0 {
		fmt.Println(result)
		return result
	}
	if len(pa) < 2 {
		result = append(result, pa)
		fmt.Println(result)
		return result
	}
	for i := 0; i < len(pa)-1; i++ {
		for j := i + 1; j < len(pa); j++ {
			p := []int32{pa[i], pa[j]}
			p2 := []int32{pa[j], pa[i]}
			result = append(result, p, p2)
		}
	}
	return result
}
