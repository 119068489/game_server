package client_square

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type IClient2Square interface {
	RpcLogin(reqMsg *LoginMsg) *base.Empty
	RpcLogin_(reqMsg *LoginMsg) (*base.Empty, easygo.IRpcInterrupt)
	RpcFlushSquareDynamic(reqMsg *FlushInfo) *base.Empty
	RpcFlushSquareDynamic_(reqMsg *FlushInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcNewVersionFlushSquareDynamic(reqMsg *NewVersionFlushInfo) *base.Empty
	RpcNewVersionFlushSquareDynamic_(reqMsg *NewVersionFlushInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddSquareDynamic(reqMsg *share_message.DynamicData) *client_server.RequestInfo
	RpcAddSquareDynamic_(reqMsg *share_message.DynamicData) (*client_server.RequestInfo, easygo.IRpcInterrupt)
	RpcDelSquareDynamic(reqMsg *client_server.RequestInfo) *base.Empty
	RpcDelSquareDynamic_(reqMsg *client_server.RequestInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcZanOperateSquareDynamic(reqMsg *client_server.ZanInfo) *base.Empty
	RpcZanOperateSquareDynamic_(reqMsg *client_server.ZanInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddCommentSquareDynamic(reqMsg *share_message.CommentData) *share_message.CommentData
	RpcAddCommentSquareDynamic_(reqMsg *share_message.CommentData) (*share_message.CommentData, easygo.IRpcInterrupt)
	RpcDelCommentSquareDynamic(reqMsg *client_server.IdInfo) *base.Empty
	RpcDelCommentSquareDynamic_(reqMsg *client_server.IdInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcAttentioPlayer(reqMsg *client_server.AttenInfo) *base.Empty
	RpcAttentioPlayer_(reqMsg *client_server.AttenInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetDynamicInfo(reqMsg *client_server.IdInfo) *share_message.DynamicData
	RpcGetDynamicInfo_(reqMsg *client_server.IdInfo) (*share_message.DynamicData, easygo.IRpcInterrupt)
	RpcGetDynamicMainComment(reqMsg *client_server.IdInfo) *share_message.CommentList
	RpcGetDynamicMainComment_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt)
	RpcGetDynamicSecondaryComment(reqMsg *client_server.IdInfo) *share_message.CommentList
	RpcGetDynamicSecondaryComment_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt)
	RpcGetDynamicInfoNew(reqMsg *client_server.IdInfo) *share_message.DynamicData
	RpcGetDynamicInfoNew_(reqMsg *client_server.IdInfo) (*share_message.DynamicData, easygo.IRpcInterrupt)
	RpcGetDynamicMainCommentNew(reqMsg *client_server.IdInfo) *share_message.CommentList
	RpcGetDynamicMainCommentNew_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt)
	RpcGetDynamicSecondaryCommentNew(reqMsg *client_server.IdInfo) *share_message.CommentList
	RpcGetDynamicSecondaryCommentNew_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt)
	RpcGetSquareMessage(reqMsg *client_server.IdInfo) *MessageMainInfo
	RpcGetSquareMessage_(reqMsg *client_server.IdInfo) (*MessageMainInfo, easygo.IRpcInterrupt)
	RpcGetPlayerZanInfo(reqMsg *client_server.RequestInfo) *ZanList
	RpcGetPlayerZanInfo_(reqMsg *client_server.RequestInfo) (*ZanList, easygo.IRpcInterrupt)
	RpcGetPlayerAttentionInfo(reqMsg *client_server.RequestInfo) *AttentionList
	RpcGetPlayerAttentionInfo_(reqMsg *client_server.RequestInfo) (*AttentionList, easygo.IRpcInterrupt)
	RpcDynamicTop(reqMsg *DynamicTopReq) *base.Empty
	RpcDynamicTop_(reqMsg *DynamicTopReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcReadPlayerInfo(reqMsg *UnReadInfo) *base.Empty
	RpcReadPlayerInfo_(reqMsg *UnReadInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcLogOut(reqMsg *base.Empty) *base.Empty
	RpcLogOut_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcFirstLoginSquare(reqMsg *FirstLoginSquareReq) *FirstLoginSquareReply
	RpcFirstLoginSquare_(reqMsg *FirstLoginSquareReq) (*FirstLoginSquareReply, easygo.IRpcInterrupt)
	RpcAdvDetail(reqMsg *share_message.AdvSetting) *AdvDetailReply
	RpcAdvDetail_(reqMsg *share_message.AdvSetting) (*AdvDetailReply, easygo.IRpcInterrupt)
	RpcAddAdvLog(reqMsg *share_message.AdvLogReq) *base.Empty
	RpcAddAdvLog_(reqMsg *share_message.AdvLogReq) (*base.Empty, easygo.IRpcInterrupt)
}

type Client2Square struct {
	Sender easygo.IMessageSender
}

func (self *Client2Square) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Client2Square) RpcLogin(reqMsg *LoginMsg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcLogin_(reqMsg *LoginMsg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Square) RpcFlushSquareDynamic(reqMsg *FlushInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcFlushSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcFlushSquareDynamic_(reqMsg *FlushInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcFlushSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Square) RpcNewVersionFlushSquareDynamic(reqMsg *NewVersionFlushInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcNewVersionFlushSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcNewVersionFlushSquareDynamic_(reqMsg *NewVersionFlushInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNewVersionFlushSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Square) RpcAddSquareDynamic(reqMsg *share_message.DynamicData) *client_server.RequestInfo {
	msg, e := self.Sender.CallRpcMethod("RpcAddSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.RequestInfo)
}

func (self *Client2Square) RpcAddSquareDynamic_(reqMsg *share_message.DynamicData) (*client_server.RequestInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.RequestInfo), e
}
func (self *Client2Square) RpcDelSquareDynamic(reqMsg *client_server.RequestInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcDelSquareDynamic_(reqMsg *client_server.RequestInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Square) RpcZanOperateSquareDynamic(reqMsg *client_server.ZanInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcZanOperateSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcZanOperateSquareDynamic_(reqMsg *client_server.ZanInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcZanOperateSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Square) RpcAddCommentSquareDynamic(reqMsg *share_message.CommentData) *share_message.CommentData {
	msg, e := self.Sender.CallRpcMethod("RpcAddCommentSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CommentData)
}

func (self *Client2Square) RpcAddCommentSquareDynamic_(reqMsg *share_message.CommentData) (*share_message.CommentData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddCommentSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CommentData), e
}
func (self *Client2Square) RpcDelCommentSquareDynamic(reqMsg *client_server.IdInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelCommentSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcDelCommentSquareDynamic_(reqMsg *client_server.IdInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelCommentSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Square) RpcAttentioPlayer(reqMsg *client_server.AttenInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAttentioPlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcAttentioPlayer_(reqMsg *client_server.AttenInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAttentioPlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Square) RpcGetDynamicInfo(reqMsg *client_server.IdInfo) *share_message.DynamicData {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DynamicData)
}

func (self *Client2Square) RpcGetDynamicInfo_(reqMsg *client_server.IdInfo) (*share_message.DynamicData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DynamicData), e
}
func (self *Client2Square) RpcGetDynamicMainComment(reqMsg *client_server.IdInfo) *share_message.CommentList {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicMainComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CommentList)
}

func (self *Client2Square) RpcGetDynamicMainComment_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicMainComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CommentList), e
}
func (self *Client2Square) RpcGetDynamicSecondaryComment(reqMsg *client_server.IdInfo) *share_message.CommentList {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicSecondaryComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CommentList)
}

func (self *Client2Square) RpcGetDynamicSecondaryComment_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicSecondaryComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CommentList), e
}
func (self *Client2Square) RpcGetDynamicInfoNew(reqMsg *client_server.IdInfo) *share_message.DynamicData {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicInfoNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DynamicData)
}

func (self *Client2Square) RpcGetDynamicInfoNew_(reqMsg *client_server.IdInfo) (*share_message.DynamicData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicInfoNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DynamicData), e
}
func (self *Client2Square) RpcGetDynamicMainCommentNew(reqMsg *client_server.IdInfo) *share_message.CommentList {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicMainCommentNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CommentList)
}

func (self *Client2Square) RpcGetDynamicMainCommentNew_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicMainCommentNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CommentList), e
}
func (self *Client2Square) RpcGetDynamicSecondaryCommentNew(reqMsg *client_server.IdInfo) *share_message.CommentList {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicSecondaryCommentNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CommentList)
}

func (self *Client2Square) RpcGetDynamicSecondaryCommentNew_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicSecondaryCommentNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CommentList), e
}
func (self *Client2Square) RpcGetSquareMessage(reqMsg *client_server.IdInfo) *MessageMainInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetSquareMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MessageMainInfo)
}

func (self *Client2Square) RpcGetSquareMessage_(reqMsg *client_server.IdInfo) (*MessageMainInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSquareMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MessageMainInfo), e
}
func (self *Client2Square) RpcGetPlayerZanInfo(reqMsg *client_server.RequestInfo) *ZanList {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerZanInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ZanList)
}

func (self *Client2Square) RpcGetPlayerZanInfo_(reqMsg *client_server.RequestInfo) (*ZanList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerZanInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ZanList), e
}
func (self *Client2Square) RpcGetPlayerAttentionInfo(reqMsg *client_server.RequestInfo) *AttentionList {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerAttentionInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AttentionList)
}

func (self *Client2Square) RpcGetPlayerAttentionInfo_(reqMsg *client_server.RequestInfo) (*AttentionList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerAttentionInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AttentionList), e
}
func (self *Client2Square) RpcDynamicTop(reqMsg *DynamicTopReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDynamicTop", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcDynamicTop_(reqMsg *DynamicTopReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDynamicTop", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Square) RpcReadPlayerInfo(reqMsg *UnReadInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReadPlayerInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcReadPlayerInfo_(reqMsg *UnReadInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReadPlayerInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Square) RpcLogOut(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcLogOut", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcLogOut_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLogOut", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Client2Square) RpcFirstLoginSquare(reqMsg *FirstLoginSquareReq) *FirstLoginSquareReply {
	msg, e := self.Sender.CallRpcMethod("RpcFirstLoginSquare", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*FirstLoginSquareReply)
}

func (self *Client2Square) RpcFirstLoginSquare_(reqMsg *FirstLoginSquareReq) (*FirstLoginSquareReply, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcFirstLoginSquare", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*FirstLoginSquareReply), e
}
func (self *Client2Square) RpcAdvDetail(reqMsg *share_message.AdvSetting) *AdvDetailReply {
	msg, e := self.Sender.CallRpcMethod("RpcAdvDetail", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AdvDetailReply)
}

func (self *Client2Square) RpcAdvDetail_(reqMsg *share_message.AdvSetting) (*AdvDetailReply, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAdvDetail", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AdvDetailReply), e
}
func (self *Client2Square) RpcAddAdvLog(reqMsg *share_message.AdvLogReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddAdvLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Client2Square) RpcAddAdvLog_(reqMsg *share_message.AdvLogReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddAdvLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}

// ==========================================================
type ISquare2Client interface {
	RpcSquareAllDynamic(reqMsg *AllInfo)
	RpcNewMessage(reqMsg *NewUnReadMessageResp)
	RpcNoNewMessage(reqMsg *base.Empty)
	RpcNewVersionSquareAllDynamic(reqMsg *NewVersionAllInfo)
}

type Square2Client struct {
	Sender easygo.IMessageSender
}

func (self *Square2Client) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Square2Client) RpcSquareAllDynamic(reqMsg *AllInfo) {
	self.Sender.CallRpcMethod("RpcSquareAllDynamic", reqMsg)
}
func (self *Square2Client) RpcNewMessage(reqMsg *NewUnReadMessageResp) {
	self.Sender.CallRpcMethod("RpcNewMessage", reqMsg)
}
func (self *Square2Client) RpcNoNewMessage(reqMsg *base.Empty) {
	self.Sender.CallRpcMethod("RpcNoNewMessage", reqMsg)
}
func (self *Square2Client) RpcNewVersionSquareAllDynamic(reqMsg *NewVersionAllInfo) {
	self.Sender.CallRpcMethod("RpcNewVersionSquareAllDynamic", reqMsg)
}
