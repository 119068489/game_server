package server_server

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type IServer2Server interface {
	RpcMsgToHallClient(reqMsg *share_message.MsgToClient)
	RpcMsgToOtherServer(reqMsg *share_message.MsgToServer) *share_message.MsgToServer
	RpcMsgToOtherServer_(reqMsg *share_message.MsgToServer) (*share_message.MsgToServer, easygo.IRpcInterrupt)
}

type Server2Server struct {
	Sender easygo.IMessageSender
}

func (self *Server2Server) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Server2Server) RpcMsgToHallClient(reqMsg *share_message.MsgToClient) {
	self.Sender.CallRpcMethod("RpcMsgToHallClient", reqMsg)
}
func (self *Server2Server) RpcMsgToOtherServer(reqMsg *share_message.MsgToServer) *share_message.MsgToServer {
	msg, e := self.Sender.CallRpcMethod("RpcMsgToOtherServer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.MsgToServer)
}

func (self *Server2Server) RpcMsgToOtherServer_(reqMsg *share_message.MsgToServer) (*share_message.MsgToServer, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMsgToOtherServer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.MsgToServer), e
}
