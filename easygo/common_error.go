package easygo

import (
	"errors"
	"fmt"
)

type ICommonError interface {
	error
	Code() int
}

type CommonError struct {
	error
	code int
	Me   ICommonError
}

func NewCommonError(text string, code ...int) *CommonError {
	p := &CommonError{}
	p.Init(p, text, code...)
	return p
}

func (self *CommonError) Init(me ICommonError, text string, code ...int) {
	self.Me = me
	self.error = errors.New(text)
	self.code = append(code, 0)[0]
}

func (self *CommonError) Code() int {
	return self.code
}

// 为什么不工作，暂时没有调通
func (self *CommonError) String() string {
	code, err := self.Me.Code(), self.Me.Error()
	if code == 0 {
		return fmt.Sprintf("error=%s", err)
	} else {
		return fmt.Sprintf("code=%d,error=%s", code, err)
	}
}
