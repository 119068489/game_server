package client_esport

import (
	"game_server/easygo"
	"game_server/easygo/base"
)

type _ = base.NoReturn

type IClient2ESports interface {
	RpcLogin(reqMsg *LoginMsg) *base.Empty
	RpcLogin_(reqMsg *LoginMsg) (*base.Empty, easygo.IRpcInterrupt)
	RpcLogOut(reqMsg *base.Empty) *base.Empty
	RpcLogOut_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
}

type Client2ESports struct {
	Sender easygo.IMessageSender
}

func (self *Client2ESports) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Client2ESports) RpcLogin(reqMsg *LoginMsg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2ESports) RpcLogin_(reqMsg *LoginMsg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2ESports) RpcLogOut(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcLogOut", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2ESports) RpcLogOut_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLogOut", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}

// ==========================================================
type IESports2Client interface {
	RpcPushBroadcast(reqMsg *BroadcastMsg)
}

type ESports2Client struct {
	Sender easygo.IMessageSender
}

func (self *ESports2Client) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *ESports2Client) RpcPushBroadcast(reqMsg *BroadcastMsg) {
	self.Sender.CallRpcMethod("RpcPushBroadcast", reqMsg)
}
