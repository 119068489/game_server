package client_hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
)

type _ = base.NoReturn

type IClient2ESportsGame interface {
	RpcESportGetGameList(reqMsg *ESportGameListRequest) *ESportGameListResult
	RpcESportGetGameList_(reqMsg *ESportGameListRequest) (*ESportGameListResult, easygo.IRpcInterrupt)
	RpcESportGetGameDetail(reqMsg *GameDetailRequest) *ESportGameDetailResult
	RpcESportGetGameDetail_(reqMsg *GameDetailRequest) (*ESportGameDetailResult, easygo.IRpcInterrupt)
	RpcESportGetGameGuessCartPoll(reqMsg *GameGuessCartRequest) *GameGuessCartResult
	RpcESportGetGameGuessCartPoll_(reqMsg *GameGuessCartRequest) (*GameGuessCartResult, easygo.IRpcInterrupt)
	RpcESportGameGuessBet(reqMsg *GameGuessBetRequest) *GameGuessBetResult
	RpcESportGameGuessBet_(reqMsg *GameGuessBetRequest) (*GameGuessBetResult, easygo.IRpcInterrupt)
	RpcESportGameHistoryData(reqMsg *GameHistoryRequest) *GameHistoryResult
	RpcESportGameHistoryData_(reqMsg *GameHistoryRequest) (*GameHistoryResult, easygo.IRpcInterrupt)
	RpcESportGameRealTimeData(reqMsg *GameRealTimeRequest) *GameRealTimeResult
	RpcESportGameRealTimeData_(reqMsg *GameRealTimeRequest) (*GameRealTimeResult, easygo.IRpcInterrupt)
}

type Client2ESportsGame struct {
	Sender easygo.IMessageSender
}

func (self *Client2ESportsGame) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Client2ESportsGame) RpcESportGetGameList(reqMsg *ESportGameListRequest) *ESportGameListResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetGameList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportGameListResult)
}

func (self *Client2ESportsGame) RpcESportGetGameList_(reqMsg *ESportGameListRequest) (*ESportGameListResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetGameList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportGameListResult), e
}
func (self *Client2ESportsGame) RpcESportGetGameDetail(reqMsg *GameDetailRequest) *ESportGameDetailResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetGameDetail", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportGameDetailResult)
}

func (self *Client2ESportsGame) RpcESportGetGameDetail_(reqMsg *GameDetailRequest) (*ESportGameDetailResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetGameDetail", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportGameDetailResult), e
}
func (self *Client2ESportsGame) RpcESportGetGameGuessCartPoll(reqMsg *GameGuessCartRequest) *GameGuessCartResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetGameGuessCartPoll", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GameGuessCartResult)
}

func (self *Client2ESportsGame) RpcESportGetGameGuessCartPoll_(reqMsg *GameGuessCartRequest) (*GameGuessCartResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetGameGuessCartPoll", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GameGuessCartResult), e
}
func (self *Client2ESportsGame) RpcESportGameGuessBet(reqMsg *GameGuessBetRequest) *GameGuessBetResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGameGuessBet", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GameGuessBetResult)
}

func (self *Client2ESportsGame) RpcESportGameGuessBet_(reqMsg *GameGuessBetRequest) (*GameGuessBetResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGameGuessBet", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GameGuessBetResult), e
}
func (self *Client2ESportsGame) RpcESportGameHistoryData(reqMsg *GameHistoryRequest) *GameHistoryResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGameHistoryData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GameHistoryResult)
}

func (self *Client2ESportsGame) RpcESportGameHistoryData_(reqMsg *GameHistoryRequest) (*GameHistoryResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGameHistoryData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GameHistoryResult), e
}
func (self *Client2ESportsGame) RpcESportGameRealTimeData(reqMsg *GameRealTimeRequest) *GameRealTimeResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGameRealTimeData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GameRealTimeResult)
}

func (self *Client2ESportsGame) RpcESportGameRealTimeData_(reqMsg *GameRealTimeRequest) (*GameRealTimeResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGameRealTimeData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GameRealTimeResult), e
}

// ==========================================================
type IESportsGame2Client interface {
}

type ESportsGame2Client struct {
	Sender easygo.IMessageSender
}

func (self *ESportsGame2Client) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------
