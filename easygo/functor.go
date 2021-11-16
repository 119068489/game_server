package easygo

import (
	"reflect"
)

func Functor(function interface{}, args1 ...interface{}) func(...interface{}) []reflect.Value { //
	f := reflect.ValueOf(function)
	if f.Kind() != reflect.Func {
		panic("第1个参数必须是个 function")
	}
	return func(args2 ...interface{}) []reflect.Value {
		in := make([]reflect.Value, len(args1)+len(args2))
		i := 0
		for _, arg := range args2 {
			in[i] = reflect.ValueOf(arg)
			i++
		}
		for _, arg := range args1 {
			in[i] = reflect.ValueOf(arg)
			i++
		}
		// result := f.CallSlice(in)
		result := f.Call(in)
		return result
	}
}

/*
func bar(s string,a int ,b int ,c int)(string,int ,int ,int){
	fmt.Println("s=",s , reflect.TypeOf(s))
	fmt.Println("a=",a , reflect.TypeOf(a))
	fmt.Println("b=",b , reflect.TypeOf(b))
	return s,a,b,c
}

func main() {
	// fmt.Println("s=",s , reflect.TypeOf(s))
	// f1 := Functor(foo)
	// f1()

	fmt.Println("----------------------")

	// f2 := Functor(bar, "ssss", 1, 2, 3)

	// f2() //

	// f3 := Functor(bar, 1, 2, 3)
	// f3("ssss")

	// f4 := Functor(bar, 1, 2)
	// f4("ssss", 3)

	f5 := Functor(bar, 1)
	v := f5("ssss", 2, 3)
	fmt.Println("v=",v,reflect.TypeOf(v))
	// f6 := Functor(bar)

	// s, a, b, c :=f6("ssss", 1, 2, 3)

	// slice :=f6("ssss", 1, 2, 3)
	// s, a, b, c := slice[0], slice[1], slice[2], slice[3]

	// fmt.Println(s, a, b, c)

	//var callback func() //
	// var callback func(...interface{}) //

	// callback = func(args2 ...interface{})([]reflect.Value){
	// 	return make([]reflect.Value, 0)
	// }

	// callback()

	// f7 = Functor(bar, 1)
	// f7("ssss", 2, 3)
}
*/
