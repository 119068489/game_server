package client_topic

import (
	"game_server/easygo"
	"game_server/easygo/base"
)

type _ = base.NoReturn

type IClient2Topic interface {
	RpcLogin(reqMsg *LoginMsg) *base.Empty
	RpcLogin_(reqMsg *LoginMsg) (*base.Empty, easygo.IRpcInterrupt)
}

type Client2Topic struct {
	Sender easygo.IMessageSender
}

func (self *Client2Topic) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Client2Topic) RpcLogin(reqMsg *LoginMsg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Topic) RpcLogin_(reqMsg *LoginMsg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}

// ==========================================================
type ITopic2Client interface {
	RpcLoginResp(reqMsg *LoginMsg) *base.Empty
	RpcLoginResp_(reqMsg *LoginMsg) (*base.Empty, easygo.IRpcInterrupt)
}

type Topic2Client struct {
	Sender easygo.IMessageSender
}

func (self *Topic2Client) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Topic2Client) RpcLoginResp(reqMsg *LoginMsg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcLoginResp", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Topic2Client) RpcLoginResp_(reqMsg *LoginMsg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLoginResp", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
