package topic_hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type ITopic2Hall interface {
	RpcTopic2HallClient(reqMsg *share_message.MsgToClient)
}

type Topic2Hall struct {
	Sender easygo.IMessageSender
}

func (self *Topic2Hall) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Topic2Hall) RpcTopic2HallClient(reqMsg *share_message.MsgToClient) {
	self.Sender.CallRpcMethod("RpcTopic2HallClient", reqMsg)
}

// ==========================================================
type IHall2Topic interface {
	RpcHall2Topic(reqMsg *share_message.MsgToClient)
}

type Hall2Topic struct {
	Sender easygo.IMessageSender
}

func (self *Hall2Topic) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Hall2Topic) RpcHall2Topic(reqMsg *share_message.MsgToClient) {
	self.Sender.CallRpcMethod("RpcHall2Topic", reqMsg)
}
