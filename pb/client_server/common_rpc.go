package client_server

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type IClient2server interface {
	RpcHeartbeat(reqMsg *NTP) *NTP
	RpcHeartbeat_(reqMsg *NTP) (*NTP, easygo.IRpcInterrupt)
	RpcTFToServer(reqMsg *ClientInfo)
	RpcBtnClick(reqMsg *BtnClickInfo) *base.Empty
	RpcBtnClick_(reqMsg *BtnClickInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcPageRegLogLoad(reqMsg *PageRegLogLoad) *base.Empty
	RpcPageRegLogLoad_(reqMsg *PageRegLogLoad) (*base.Empty, easygo.IRpcInterrupt)
}

type Client2server struct {
	Sender easygo.IMessageSender
}

func (self *Client2server) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Client2server) RpcHeartbeat(reqMsg *NTP) *NTP {
	msg, e := self.Sender.CallRpcMethod("RpcHeartbeat", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NTP)
}

func (self *Client2server) RpcHeartbeat_(reqMsg *NTP) (*NTP, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcHeartbeat", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NTP), e
}
func (self *Client2server) RpcTFToServer(reqMsg *ClientInfo) {
	self.Sender.CallRpcMethod("RpcTFToServer", reqMsg)
}
func (self *Client2server) RpcBtnClick(reqMsg *BtnClickInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcBtnClick", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2server) RpcBtnClick_(reqMsg *BtnClickInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBtnClick", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2server) RpcPageRegLogLoad(reqMsg *PageRegLogLoad) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcPageRegLogLoad", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2server) RpcPageRegLogLoad_(reqMsg *PageRegLogLoad) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPageRegLogLoad", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}

// ==========================================================
type IServer2Client interface {
	RpcToast(reqMsg *ToastMsg)
	RpcBroadCastMsg(reqMsg *share_message.BroadCastMsg)
	RpcPlayerAttrChange(reqMsg *PlayerMsg)
	RpcStopBroad(reqMsg *BroadIdReq)
	RpcPlayerTimeoutBeKick(reqMsg *PlayerTimeoutBeKick)
}

type Server2Client struct {
	Sender easygo.IMessageSender
}

func (self *Server2Client) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Server2Client) RpcToast(reqMsg *ToastMsg) {
	self.Sender.CallRpcMethod("RpcToast", reqMsg)
}
func (self *Server2Client) RpcBroadCastMsg(reqMsg *share_message.BroadCastMsg) {
	self.Sender.CallRpcMethod("RpcBroadCastMsg", reqMsg)
}
func (self *Server2Client) RpcPlayerAttrChange(reqMsg *PlayerMsg) {
	self.Sender.CallRpcMethod("RpcPlayerAttrChange", reqMsg)
}
func (self *Server2Client) RpcStopBroad(reqMsg *BroadIdReq) {
	self.Sender.CallRpcMethod("RpcStopBroad", reqMsg)
}
func (self *Server2Client) RpcPlayerTimeoutBeKick(reqMsg *PlayerTimeoutBeKick) {
	self.Sender.CallRpcMethod("RpcPlayerTimeoutBeKick", reqMsg)
}
