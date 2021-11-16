package util

// MergeMap 合并两个map[int]int
func MergeMap(args ...map[int]int) map[int]int {
	mergedMap := map[int]int{}
	for _, m := range args {
		for k, v := range m {
			if _, exists := mergedMap[k]; exists {
				mergedMap[k] += v
			} else {
				mergedMap[k] = v
			}
		}
	}
	return mergedMap
}

// GetMapMinValue 获取map[int]int中最小值, 并返回所有key的集合
// value约定是大于0的
func GetMapMinValue(m map[int]int) (int, []int) {
	// value到keys的对应关系
	keys := []int{}
	// 最小值
	minValue := -1
	for k, v := range m {
		if minValue == -1 || v < minValue {
			minValue = v
			keys = []int{k}
		} else if v == minValue {
			keys = append(keys, k)
		}
	}
	return minValue, keys
}

// GetMapMaxValue 获取map[int]int中最大值, 并返回所有key的集合
// value约定是大于0的
func GetMapMaxValue(m map[int]int) (int, []int) {
	// value到keys的对应关系
	keys := []int{}
	// 最大值
	maxValue := -1
	for k, v := range m {
		if v > maxValue {
			maxValue = v
			keys = []int{k}
		} else if v == maxValue {
			keys = append(keys, k)
		}
	}
	return maxValue, keys
}

// GetMapValues 获取map[int]int结构的所有value
// 返回结果去重
func GetMapValues(m map[int]int) []int {
	values := make([]int, 0)
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
