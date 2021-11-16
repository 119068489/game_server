package client_login

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/client_server"
)

type _ = base.NoReturn

type IClient2Login interface {
	RpcLoginHall(reqMsg *LoginMsg) *LoginResult
	RpcLoginHall_(reqMsg *LoginMsg) (*LoginResult, easygo.IRpcInterrupt)
	RpcClientGetCode(reqMsg *client_server.GetCodeRequest) *base.Empty
	RpcClientGetCode_(reqMsg *client_server.GetCodeRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcCheckMessageCode(reqMsg *client_server.CodeResponse) *base.Empty
	RpcCheckMessageCode_(reqMsg *client_server.CodeResponse) (*base.Empty, easygo.IRpcInterrupt)
	RpcForgetLoginPassword(reqMsg *LoginMsg) *base.Empty
	RpcForgetLoginPassword_(reqMsg *LoginMsg) (*base.Empty, easygo.IRpcInterrupt)
	RpcCheckAccountVaild(reqMsg *client_server.CheckInfo) *client_server.CheckInfo
	RpcCheckAccountVaild_(reqMsg *client_server.CheckInfo) (*client_server.CheckInfo, easygo.IRpcInterrupt)
	RpcAccountCancel(reqMsg *AccountCancel) *AccountCancel
	RpcAccountCancel_(reqMsg *AccountCancel) (*AccountCancel, easygo.IRpcInterrupt)
}

type Client2Login struct {
	Sender easygo.IMessageSender
}

func (self *Client2Login) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Client2Login) RpcLoginHall(reqMsg *LoginMsg) *LoginResult {
	msg, e := self.Sender.CallRpcMethod("RpcLoginHall", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LoginResult)
}

func (self *Client2Login) RpcLoginHall_(reqMsg *LoginMsg) (*LoginResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLoginHall", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LoginResult), e
}
func (self *Client2Login) RpcClientGetCode(reqMsg *client_server.GetCodeRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcClientGetCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Login) RpcClientGetCode_(reqMsg *client_server.GetCodeRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcClientGetCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Login) RpcCheckMessageCode(reqMsg *client_server.CodeResponse) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCheckMessageCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Login) RpcCheckMessageCode_(reqMsg *client_server.CodeResponse) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckMessageCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Login) RpcForgetLoginPassword(reqMsg *LoginMsg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcForgetLoginPassword", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Login) RpcForgetLoginPassword_(reqMsg *LoginMsg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcForgetLoginPassword", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Login) RpcCheckAccountVaild(reqMsg *client_server.CheckInfo) *client_server.CheckInfo {
	msg, e := self.Sender.CallRpcMethod("RpcCheckAccountVaild", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.CheckInfo)
}

func (self *Client2Login) RpcCheckAccountVaild_(reqMsg *client_server.CheckInfo) (*client_server.CheckInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckAccountVaild", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.CheckInfo), e
}
func (self *Client2Login) RpcAccountCancel(reqMsg *AccountCancel) *AccountCancel {
	msg, e := self.Sender.CallRpcMethod("RpcAccountCancel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AccountCancel)
}

func (self *Client2Login) RpcAccountCancel_(reqMsg *AccountCancel) (*AccountCancel, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAccountCancel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AccountCancel), e
}

// ==========================================================
type ILogin2Client interface {
}

type Login2Client struct {
	Sender easygo.IMessageSender
}

func (self *Login2Client) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------
