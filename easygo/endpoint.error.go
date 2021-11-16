package easygo

import (
	"errors"
	"fmt"
	"game_server/easygo/base"
	"strings"
)

//---------------------------------------------------------------------------------------

type IRpcInterrupt interface {
	error
	__IRpcInterrupt__()
	Reason() string
	Code() string
	AddPrefix(prefix interface{})
	AddPostfix(postfix interface{})
}

type RpcInterrupt struct {
	error
	Me      IRpcInterrupt
	failMsg *base.Fail
}

/* 抽象类，不提供实例化方法
func NewRpcInterrupt(text string)*RpcInterrupt{}
*/

func (self *RpcInterrupt) Init(me IRpcInterrupt, methodName string, failMsg *base.Fail) {
	self.Me = me

	var list []string
	if methodName != "" {
		list = append(list, "method: "+methodName)
	}

	code := failMsg.GetCode()
	if code != "" {
		list = append(list, "code: "+code)
	}

	reason := failMsg.GetReason()
	if reason != "" {
		list = append(list, "reason: "+reason)
	}

	text := strings.Join(list, ";")
	self.error = errors.New(text)
	self.failMsg = failMsg
}

func (self *RpcInterrupt) Reason() string {
	return self.failMsg.GetReason() // 没有把 methodName 拼在这里
}

func (self *RpcInterrupt) Code() string {
	return self.failMsg.GetCode()
}

func (self *RpcInterrupt) AddPrefix(prefix interface{}) {
	s := fmt.Sprintf("%v;%v", prefix, self.error)
	self.error = errors.New(s)
}

func (self *RpcInterrupt) AddPostfix(postfix interface{}) {
	s := fmt.Sprintf("%v;%v", self.error, postfix)
	self.error = errors.New(s)
}

func (self *RpcInterrupt) __IRpcInterrupt__() {}

//----------------------------------------------------------------------------------------

type IRpcFail interface {
	IRpcInterrupt
	__IRpcFail__()
}

type RpcFail struct {
	RpcInterrupt
}

func NewRpcFail(methodName string, failMsg *base.Fail) *RpcFail {
	p := &RpcFail{}
	p.Init(p, methodName, failMsg)
	return p
}

func (self *RpcFail) __IRpcFail__() {}

//----------------------------------------------------------------------------------------

type IRpcTimeout interface {
	IRpcInterrupt
	__IRpcTimeout__()
}
type RpcTimeout struct {
	RpcInterrupt
}

func NewRpcTimeout(methodName string, failMsg *base.Fail) *RpcTimeout {
	p := &RpcTimeout{}
	p.Init(p, methodName, failMsg)
	return p
}

func (self *RpcTimeout) __IRpcTimeout__() {}
