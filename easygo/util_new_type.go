package easygo

import (
	"fmt"
	"reflect"
)

func NewBool(val interface{}) *bool {

	if val == nil {
		return nil
	}

	v := new(bool)

	switch value := val.(type) {
	case *bool:
		*v = *value
	case bool:
		*v = value
	default:
		s := fmt.Sprintf("NewBool 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewString(val interface{}) *string {

	if val == nil {
		return nil
	}

	v := new(string)

	switch value := val.(type) {
	case *string:
		*v = *value
	case string:
		*v = value
	case int64:
		*v = AnytoA(value)
	case int:
		*v = AnytoA(value)
	default:
		s := fmt.Sprintf("NewString 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewInt(val interface{}) *int {

	if val == nil {
		return nil
	}

	v := new(int)

	switch value := val.(type) {
	case *int:
		*v = *value
	case int:
		*v = value
	case int8:
		*v = int(value)
	case int16:
		*v = int(value)
	case int32:
		*v = int(value)
	case int64:
		*v = int(value)
	case float32:
		*v = int(value)
	case float64:
		*v = int(value)
	case uint:
		*v = int(value)
	case uint8:
		*v = int(value)
	case uint16:
		*v = int(value)
	case uint32:
		*v = int(value)
	case uint64:
		*v = int(value)
	default:
		s := fmt.Sprintf("NewInt 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewInt32(val interface{}) *int32 {

	if val == nil {
		return nil
	}

	v := new(int32)

	switch value := val.(type) {
	case *int32:
		*v = *value
	case int32:
		*v = value
	case int:
		*v = int32(value)
	case int8:
		*v = int32(value)
	case int16:
		*v = int32(value)
	case int64:
		*v = int32(value)
	case float32:
		*v = int32(value)
	case float64:
		*v = int32(value)
	case uint:
		*v = int32(value)
	case uint8:
		*v = int32(value)
	case uint16:
		*v = int32(value)
	case uint32:
		*v = int32(value)
	case uint64:
		*v = int32(value)
	case string:
		*v = AtoInt32(value)
	default:
		s := fmt.Sprintf("NewInt32 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewFloat32(val interface{}) *float32 {

	if val == nil {
		return nil
	}

	v := new(float32)

	switch value := val.(type) {
	case *float32:
		*v = *value
	case float32:
		*v = value
	case int:
		*v = float32(value)
	case int8:
		*v = float32(value)
	case int16:
		*v = float32(value)
	case int32:
		*v = float32(value)
	case int64:
		*v = float32(value)
	case float64:
		*v = float32(value)
	case uint:
		*v = float32(value)
	case uint8:
		*v = float32(value)
	case uint16:
		*v = float32(value)
	case uint32:
		*v = float32(value)
	case uint64:
		*v = float32(value)
	default:
		s := fmt.Sprintf("NewFloat32 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewInt64(val interface{}) *int64 {

	if val == nil {
		return nil
	}

	v := new(int64)

	switch value := val.(type) {
	case *int64:
		*v = *value
	case int64:
		*v = value
	case int:
		*v = int64(value)
	case int8:
		*v = int64(value)
	case int16:
		*v = int64(value)
	case int32:
		*v = int64(value)
	case float32:
		*v = int64(value)
	case float64:
		*v = int64(value)
	case uint:
		*v = int64(value)
	case uint8:
		*v = int64(value)
	case uint16:
		*v = int64(value)
	case uint32:
		*v = int64(value)
	case uint64:
		*v = int64(value)
	case string:
		*v = AtoInt64(value)
	default:
		s := fmt.Sprintf("NewInt64 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewFloat64(val interface{}) *float64 {

	if val == nil {
		return nil
	}

	v := new(float64)

	switch value := val.(type) {
	case *float64:
		*v = *value
	case float64:
		*v = value
	case int:
		*v = float64(value)
	case int8:
		*v = float64(value)
	case int16:
		*v = float64(value)
	case int32:
		*v = float64(value)
	case int64:
		*v = float64(value)
	case float32:
		*v = float64(value)
	case uint:
		*v = float64(value)
	case uint8:
		*v = float64(value)
	case uint16:
		*v = float64(value)
	case uint32:
		*v = float64(value)
	case uint64:
		*v = float64(value)
	case string:
		*v = AtoFloat64(value)
	default:
		s := fmt.Sprintf("NewFloat64 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewUint(val interface{}) *uint {
	if val == nil {
		return nil
	}
	v := new(uint)

	switch value := val.(type) {
	case *uint:
		*v = *value
	case uint:
		*v = value
	case int:
		*v = uint(value)
	case int8:
		*v = uint(value)
	case int16:
		*v = uint(value)
	case int32:
		*v = uint(value)
	case int64:
		*v = uint(value)
	case uint8:
		*v = uint(value)
	case uint16:
		*v = uint(value)
	case uint32:
		*v = uint(value)
	case uint64:
		*v = uint(value)
	default:
		s := fmt.Sprintf("NewUint 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewUint32(val interface{}) *uint32 {

	if val == nil {
		return nil
	}

	v := new(uint32)

	switch value := val.(type) {
	case *uint32:
		*v = *value
	case uint32:
		*v = value
	case int:
		*v = uint32(value)
	case int8:
		*v = uint32(value)
	case int16:
		*v = uint32(value)
	case int32:
		*v = uint32(value)
	case int64:
		*v = uint32(value)
	case uint:
		*v = uint32(value)
	case uint8:
		*v = uint32(value)
	case uint16:
		*v = uint32(value)
	case uint64:
		*v = uint32(value)
	default:
		s := fmt.Sprintf("NewUint32 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewUint64(val interface{}) *uint64 {
	if val == nil {
		return nil
	}
	v := new(uint64)

	switch value := val.(type) {
	case *uint64:
		*v = *value
	case uint64:
		*v = value
	case int:
		*v = uint64(value)
	case int8:
		*v = uint64(value)
	case int16:
		*v = uint64(value)
	case int32:
		*v = uint64(value)
	case int64:
		*v = uint64(value)
	case uint32:
		*v = uint64(value)
	case uint:
		*v = uint64(value)
	case uint8:
		*v = uint64(value)
	case uint16:
		*v = uint64(value)
	default:
		s := fmt.Sprintf("NewUint64 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewUint8(val interface{}) *uint8 {

	if val == nil {
		return nil
	}

	v := new(uint8)

	switch value := val.(type) {
	case *uint8:
		*v = *value
	case uint8:
		*v = value
	case int:
		*v = uint8(value)
	case int8:
		*v = uint8(value)
	case int16:
		*v = uint8(value)
	case int32:
		*v = uint8(value)
	case int64:
		*v = uint8(value)
	case uint:
		*v = uint8(value)
	case uint16:
		*v = uint8(value)
	case uint32:
		*v = uint8(value)
	case uint64:
		*v = uint8(value)
	default:
		s := fmt.Sprintf("NewUint8 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewUint16(val interface{}) *uint16 {

	if val == nil {
		return nil
	}

	v := new(uint16)

	switch value := val.(type) {
	case *uint16:
		*v = *value
	case uint16:
		*v = value
	case int:
		*v = uint16(value)
	case int8:
		*v = uint16(value)
	case int16:
		*v = uint16(value)
	case int32:
		*v = uint16(value)
	case int64:
		*v = uint16(value)
	case uint8:
		*v = uint16(value)
	case uint:
		*v = uint16(value)
	case uint32:
		*v = uint16(value)
	case uint64:
		*v = uint16(value)
	default:
		s := fmt.Sprintf("NewUint16 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewInt16(val interface{}) *int16 {

	if val == nil {
		return nil
	}

	v := new(int16)

	switch value := val.(type) {
	case *int16:
		*v = *value
	case int16:
		*v = value
	case int:
		*v = int16(value)
	case int8:
		*v = int16(value)
	case int32:
		*v = int16(value)
	case int64:
		*v = int16(value)
	case uint8:
		*v = int16(value)
	case uint16:
		*v = int16(value)
	case uint:
		*v = int16(value)
	case uint32:
		*v = int16(value)
	case uint64:
		*v = int16(value)
	default:
		s := fmt.Sprintf("NewInt16 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}

func NewInt8(val interface{}) *int8 {

	if val == nil {
		return nil
	}

	v := new(int8)

	switch value := val.(type) {
	case *int8:
		*v = *value
	case int8:
		*v = value
	case int:
		*v = int8(value)
	case int16:
		*v = int8(value)
	case int32:
		*v = int8(value)
	case int64:
		*v = int8(value)
	case uint8:
		*v = int8(value)
	case uint16:
		*v = int8(value)
	case uint:
		*v = int8(value)
	case uint32:
		*v = int8(value)
	case uint64:
		*v = int8(value)
	default:
		s := fmt.Sprintf("NewInt8 不支持 %v 类型,自行补充", reflect.TypeOf(val))
		panic(s)
	}
	return v
}
