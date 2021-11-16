package client_hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type IHall2ESports interface {
	RpcESPortsPlayerOnLine(reqMsg *share_message.PlayerState) *ESportCommonResult
	RpcESPortsPlayerOnLine_(reqMsg *share_message.PlayerState) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESPortsPlayerOffLine(reqMsg *share_message.PlayerState) *ESportCommonResult
	RpcESPortsPlayerOffLine_(reqMsg *share_message.PlayerState) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESPortsPushGameOrderSysMsg(reqMsg *share_message.TableESPortsGameOrderSysMsg) *ESportCommonResult
	RpcESPortsPushGameOrderSysMsg_(reqMsg *share_message.TableESPortsGameOrderSysMsg) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESportDataStatusInfo(reqMsg *ESportDataStatusInfo) *ESportCommonResult
	RpcESportDataStatusInfo_(reqMsg *ESportDataStatusInfo) (*ESportCommonResult, easygo.IRpcInterrupt)
}

type Hall2ESports struct {
	Sender easygo.IMessageSender
}

func (self *Hall2ESports) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Hall2ESports) RpcESPortsPlayerOnLine(reqMsg *share_message.PlayerState) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsPlayerOnLine", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Hall2ESports) RpcESPortsPlayerOnLine_(reqMsg *share_message.PlayerState) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsPlayerOnLine", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Hall2ESports) RpcESPortsPlayerOffLine(reqMsg *share_message.PlayerState) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsPlayerOffLine", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Hall2ESports) RpcESPortsPlayerOffLine_(reqMsg *share_message.PlayerState) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsPlayerOffLine", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Hall2ESports) RpcESPortsPushGameOrderSysMsg(reqMsg *share_message.TableESPortsGameOrderSysMsg) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsPushGameOrderSysMsg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Hall2ESports) RpcESPortsPushGameOrderSysMsg_(reqMsg *share_message.TableESPortsGameOrderSysMsg) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsPushGameOrderSysMsg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Hall2ESports) RpcESportDataStatusInfo(reqMsg *ESportDataStatusInfo) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportDataStatusInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Hall2ESports) RpcESportDataStatusInfo_(reqMsg *ESportDataStatusInfo) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportDataStatusInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
