package easygo

import (
	// "sync"
	// "time"
	//"reflect"
	//"fmt"

	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"runtime"
	"strings"
)

func PanicError(e error, prefix ...string) {
	if e != nil {
		if len(prefix) == 0 {
			panic(e)
		} else {
			// 这里会丢失原来的 e  对象类型，变成了 string 类型的 e
			s := fmt.Sprintf("%s;%v", prefix[0], e)
			panic(s)
		}
	}
}

func PanicLog() interface{} {
	if err := recover(); err != nil {
		kb := 4

		s := []byte("/src/runtime/panic.go")
		e := []byte("\ngoroutine ")
		line := []byte("\n")
		stack := make([]byte, kb<<10) //4KB
		length := runtime.Stack(stack, true)
		start := bytes.Index(stack, s)
		stack = stack[start:length]
		start = bytes.Index(stack, line) + 1
		stack = stack[start:]
		end := bytes.LastIndex(stack, line)
		if end != -1 {
			stack = stack[:end]
		}
		end = bytes.Index(stack, e)
		if end != -1 {
			stack = stack[:end]
		}
		stack = bytes.TrimRight(stack, "\n")

		logs.Error(fmt.Sprintf("%s", string(stack)))

		return err
	}

	return nil
}

//-----------------------------

// 给 async_result 和 goroutine 使用 //好像没有必要用这个来表示超时
type ITimeoutError interface {
	ICommonError
	__ITimeoutError__()
}
type TimeoutError struct {
	CommonError
}

func NewTimeoutError(text string) *TimeoutError {
	p := &TimeoutError{}
	p.Init(p, text)
	return p
}
func (self *TimeoutError) __ITimeoutError__() {}

//-----------------------------

var PanicWriter = NewFileWriter("logs/panic.log")
var RpcFailWriter = NewFileWriter("logs/rpc_fail.log")

//-----------------------------

// 取得调用栈
func CallStack(skips ...int) string {
	//debug.PrintStack()
	skip := append(skips, 0)[0]
	var slice []string
	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.Contains(file, `/Go/src/`) { // 跳过系统源码
			continue
		}
		stack := fmt.Sprintf("%s:%d\n", file, line)
		slice = append(slice, stack)
	}
	stack := strings.Join(slice, "")
	return stack
}

func LogPanicAndStack(recoverVal interface{}, skips ...int) string {
	stack := CallStack(skips...)
	text := LogException(recoverVal, stack)
	return text
}

func LogException(recoverVal interface{}, traceBack string) string {
	text := fmt.Sprintf("================================================\n@@ %v\n%s", recoverVal, traceBack)
	if _, ok := recoverVal.(IRpcInterrupt); ok {
		toStdOut := fmt.Sprintf("================================================\n[只是个警告;rpc 失败] %v\n%s", recoverVal, traceBack)
		os.Stdout.Write([]byte(toStdOut))
		RpcFailWriter.Write(text)
	} else {
		os.Stdout.Write([]byte(text))
		PanicWriter.Write(text)
	}
	return text
}

func RecoverAndLog(skips ...int) { // 能直接被 defer 使用。
	recoverVal := recover()
	if recoverVal != nil {
		LogPanicAndStack(recoverVal, skips...)
	}
}

func RecoverAndRePanic(text interface{}, isPrefix bool) { // 能直接被 defer 使用。
	recoverVal := recover()
	if recoverVal == nil {
		return
	}
	if e, ok := recoverVal.(IRpcInterrupt); ok {
		if isPrefix {
			e.AddPrefix(text)
		} else {
			e.AddPostfix(text)
		}
		panic(recoverVal) // 不能丢了类型,只能原样抛出
	} else {
		var s string
		if isPrefix {
			s = fmt.Sprintf("%v;%v", text, recoverVal)
		} else {
			s = fmt.Sprintf("%v;%v", recoverVal, text)
		}
		panic(s) // 这里会丢失原来的 e  对象类型，变成了 string 类型的 e
	}
}
