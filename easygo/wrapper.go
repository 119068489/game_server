package easygo

import (
	"encoding/json"
	"math"
	"strconv"
	// "time"
	//"reflect"
	//"fmt"
)

//-----------------------------

//字符串转int
func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

//字符串转int32
func AtoInt32(s string) int32 {
	i := Atoi(s)
	return int32(i)
}

//字符串转int64
func AtoInt64(s string) int64 {
	if s == "" {
		return 0
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

//字符串转float32
func AtoFloat32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}
	return float32(f)
}

//字符串转float64
func AtoFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

//int,int32,int64,float32,float64转字符串
func AnytoA(val interface{}) string {
	var s string
	switch value := val.(type) {
	case uint8:
		s = strconv.Itoa(int(value))
	case int:
		s = strconv.Itoa(value)
	case int32:
		s = strconv.FormatInt(int64(value), 10)
	case int64:
		s = strconv.FormatInt(int64(value), 10)
	case uint64:
		s = strconv.FormatInt(int64(value), 10)
	case float32:
		s = strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		s = strconv.FormatFloat(value, 'f', -1, 64)
	case string:
		s = value
	case map[string]interface{}:
		js, err := json.Marshal(value)
		if err != nil {
			panic(err)
		}
		s = string(js)
	default:

		panic("AnytoA 未定义的字符类型转成string，请自行添加")
	}
	return s
}

//将float64转成精确的int64
func FtoInt64(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

//将int64恢复成正常的float64
func ItoFloat64(num int64, retain int) float64 {
	return float64(num) / math.Pow10(retain)
}

//精准float64
func FToFloat64(num float64, retain int) float64 {
	return num * math.Pow10(retain)
}

//精准int64
func IToInt64(num int64, retain int) int64 {
	return int64(ItoFloat64(num, retain))
}

// 无需再包装
// func Itoa(i int){
// 	strconv.Itoa(i)
// }
