package client_hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type ISquareClient2Hall interface {
	RpcFlushSquareDynamic(reqMsg *FlushInfo) *base.Empty
	RpcFlushSquareDynamic_(reqMsg *FlushInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcNewVersionFlushSquareDynamic(reqMsg *share_message.NewVersionFlushInfo) *base.Empty
	RpcNewVersionFlushSquareDynamic_(reqMsg *share_message.NewVersionFlushInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcFlushSquareDynamicTopic(reqMsg *share_message.FlushSquareDynamicTopicReq) *base.Empty
	RpcFlushSquareDynamicTopic_(reqMsg *share_message.FlushSquareDynamicTopicReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddSquareDynamic(reqMsg *share_message.DynamicData) *client_server.RequestInfo
	RpcAddSquareDynamic_(reqMsg *share_message.DynamicData) (*client_server.RequestInfo, easygo.IRpcInterrupt)
	RpcDelSquareDynamicApi(reqMsg *client_server.RequestInfo) *base.Empty
	RpcDelSquareDynamicApi_(reqMsg *client_server.RequestInfo) (*base.Empty, easygo.IRpcInterrupt)
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
	RpcFirstLoginSquare(reqMsg *FirstLoginSquareReq) *FirstLoginSquareReply
	RpcFirstLoginSquare_(reqMsg *FirstLoginSquareReq) (*FirstLoginSquareReply, easygo.IRpcInterrupt)
	RpcAdvDetail(reqMsg *share_message.AdvSetting) *AdvDetailReply
	RpcAdvDetail_(reqMsg *share_message.AdvSetting) (*AdvDetailReply, easygo.IRpcInterrupt)
	RpcGetAllTopic(reqMsg *base.Empty) *AllTopic
	RpcGetAllTopic_(reqMsg *base.Empty) (*AllTopic, easygo.IRpcInterrupt)
	RpcGetTopicDetailReq(reqMsg *TopicDetailReq) *TopicDetailResp
	RpcGetTopicDetailReq_(reqMsg *TopicDetailReq) (*TopicDetailResp, easygo.IRpcInterrupt)
	RpcGetTopicMainPageList(reqMsg *TopicMainPageListReq) *TopicMainPageListResp
	RpcGetTopicMainPageList_(reqMsg *TopicMainPageListReq) (*TopicMainPageListResp, easygo.IRpcInterrupt)
	RpcGetTopicParticipateList(reqMsg *TopicParticipateListReq) *TopicParticipateListResp
	RpcGetTopicParticipateList_(reqMsg *TopicParticipateListReq) (*TopicParticipateListResp, easygo.IRpcInterrupt)
	RpcAttentionTopic(reqMsg *AttentionTopicReq) *base.Empty
	RpcAttentionTopic_(reqMsg *AttentionTopicReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcMyAttentionTopicList(reqMsg *MyAttentionTopicListReq) *MyAttentionTopicListResp
	RpcMyAttentionTopicList_(reqMsg *MyAttentionTopicListReq) (*MyAttentionTopicListResp, easygo.IRpcInterrupt)
	RpcSearchTopic(reqMsg *SearchTopicReq) *SearchTopicResp
	RpcSearchTopic_(reqMsg *SearchTopicReq) (*SearchTopicResp, easygo.IRpcInterrupt)
	RpcSearchHotTopic(reqMsg *base.Empty) *SearchHotTopicResp
	RpcSearchHotTopic_(reqMsg *base.Empty) (*SearchHotTopicResp, easygo.IRpcInterrupt)
	RpcGetTopicTypeList(reqMsg *base.Empty) *TopicTypeListResp
	RpcGetTopicTypeList_(reqMsg *base.Empty) (*TopicTypeListResp, easygo.IRpcInterrupt)
	RpcGetTopicList(reqMsg *TopicListReq) *TopicListResp
	RpcGetTopicList_(reqMsg *TopicListReq) (*TopicListResp, easygo.IRpcInterrupt)
	RpcFlushTopic(reqMsg *FlushTopicReq) *FlushTopicResp
	RpcFlushTopic_(reqMsg *FlushTopicReq) (*FlushTopicResp, easygo.IRpcInterrupt)
	RpcHotTopicList(reqMsg *HotTopicListReq) *HotTopicListResp
	RpcHotTopicList_(reqMsg *HotTopicListReq) (*HotTopicListResp, easygo.IRpcInterrupt)
	RpcAttentionRecommendPlayer(reqMsg *base.Empty) *AttentionRecommendPlayerResp
	RpcAttentionRecommendPlayer_(reqMsg *base.Empty) (*AttentionRecommendPlayerResp, easygo.IRpcInterrupt)
	RpcSquareAttention(reqMsg *SquareAttentionReq) *SquareAttentionResp
	RpcSquareAttention_(reqMsg *SquareAttentionReq) (*SquareAttentionResp, easygo.IRpcInterrupt)
	RpcTopicHotDynamicParticipatePlayer(reqMsg *TopicParticipateListReq) *TopicParticipateListResp
	RpcTopicHotDynamicParticipatePlayer_(reqMsg *TopicParticipateListReq) (*TopicParticipateListResp, easygo.IRpcInterrupt)
	RpcTopicHead(reqMsg *base.Empty) *TopicHeadResp
	RpcTopicHead_(reqMsg *base.Empty) (*TopicHeadResp, easygo.IRpcInterrupt)
	RpcTopicDevoteList(reqMsg *TopicDevoteListReq) *TopicDevoteListResp
	RpcTopicDevoteList_(reqMsg *TopicDevoteListReq) (*TopicDevoteListResp, easygo.IRpcInterrupt)
	RpcTopicMasterCondition(reqMsg *TopicMasterConditionReq) *TopicMasterConditionResp
	RpcTopicMasterCondition_(reqMsg *TopicMasterConditionReq) (*TopicMasterConditionResp, easygo.IRpcInterrupt)
	RpcApplyTopicMaster(reqMsg *ApplyTopicMasterReq) *ApplyTopicMasterResp
	RpcApplyTopicMaster_(reqMsg *ApplyTopicMasterReq) (*ApplyTopicMasterResp, easygo.IRpcInterrupt)
	RpcTopicMasterEdit(reqMsg *TopicMasterEditReq) *TopicMasterEditResp
	RpcTopicMasterEdit_(reqMsg *TopicMasterEditReq) (*TopicMasterEditResp, easygo.IRpcInterrupt)
	RpcTopicTop(reqMsg *TopicTopReq) *TopicTopResp
	RpcTopicTop_(reqMsg *TopicTopReq) (*TopicTopResp, easygo.IRpcInterrupt)
	RpcTopicTopCancel(reqMsg *TopicTopCancelReq) *TopicTopCancelResp
	RpcTopicTopCancel_(reqMsg *TopicTopCancelReq) (*TopicTopCancelResp, easygo.IRpcInterrupt)
	RpcTopicLeaderBoardDescription(reqMsg *base.Empty) *TopicLeaderBoardDescriptionResp
	RpcTopicLeaderBoardDescription_(reqMsg *base.Empty) (*TopicLeaderBoardDescriptionResp, easygo.IRpcInterrupt)
	RpcQuitTopicMaster(reqMsg *QuitTopicMasterReq) *QuitTopicMasterResp
	RpcQuitTopicMaster_(reqMsg *QuitTopicMasterReq) (*QuitTopicMasterResp, easygo.IRpcInterrupt)
	RpcTopicMasterDelDynamic(reqMsg *TopicMasterDelDynamicReq) *TopicMasterDelDynamicResp
	RpcTopicMasterDelDynamic_(reqMsg *TopicMasterDelDynamicReq) (*TopicMasterDelDynamicResp, easygo.IRpcInterrupt)
}

type SquareClient2Hall struct {
	Sender easygo.IMessageSender
}

func (self *SquareClient2Hall) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *SquareClient2Hall) RpcFlushSquareDynamic(reqMsg *FlushInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcFlushSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *SquareClient2Hall) RpcFlushSquareDynamic_(reqMsg *FlushInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcFlushSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *SquareClient2Hall) RpcNewVersionFlushSquareDynamic(reqMsg *share_message.NewVersionFlushInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcNewVersionFlushSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *SquareClient2Hall) RpcNewVersionFlushSquareDynamic_(reqMsg *share_message.NewVersionFlushInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNewVersionFlushSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *SquareClient2Hall) RpcFlushSquareDynamicTopic(reqMsg *share_message.FlushSquareDynamicTopicReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcFlushSquareDynamicTopic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *SquareClient2Hall) RpcFlushSquareDynamicTopic_(reqMsg *share_message.FlushSquareDynamicTopicReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcFlushSquareDynamicTopic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *SquareClient2Hall) RpcAddSquareDynamic(reqMsg *share_message.DynamicData) *client_server.RequestInfo {
	msg, e := self.Sender.CallRpcMethod("RpcAddSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.RequestInfo)
}

func (self *SquareClient2Hall) RpcAddSquareDynamic_(reqMsg *share_message.DynamicData) (*client_server.RequestInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.RequestInfo), e
}
func (self *SquareClient2Hall) RpcDelSquareDynamicApi(reqMsg *client_server.RequestInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelSquareDynamicApi", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *SquareClient2Hall) RpcDelSquareDynamicApi_(reqMsg *client_server.RequestInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelSquareDynamicApi", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *SquareClient2Hall) RpcGetDynamicInfo(reqMsg *client_server.IdInfo) *share_message.DynamicData {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DynamicData)
}

func (self *SquareClient2Hall) RpcGetDynamicInfo_(reqMsg *client_server.IdInfo) (*share_message.DynamicData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DynamicData), e
}
func (self *SquareClient2Hall) RpcGetDynamicMainComment(reqMsg *client_server.IdInfo) *share_message.CommentList {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicMainComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CommentList)
}

func (self *SquareClient2Hall) RpcGetDynamicMainComment_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicMainComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CommentList), e
}
func (self *SquareClient2Hall) RpcGetDynamicSecondaryComment(reqMsg *client_server.IdInfo) *share_message.CommentList {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicSecondaryComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CommentList)
}

func (self *SquareClient2Hall) RpcGetDynamicSecondaryComment_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicSecondaryComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CommentList), e
}
func (self *SquareClient2Hall) RpcGetDynamicInfoNew(reqMsg *client_server.IdInfo) *share_message.DynamicData {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicInfoNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DynamicData)
}

func (self *SquareClient2Hall) RpcGetDynamicInfoNew_(reqMsg *client_server.IdInfo) (*share_message.DynamicData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicInfoNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DynamicData), e
}
func (self *SquareClient2Hall) RpcGetDynamicMainCommentNew(reqMsg *client_server.IdInfo) *share_message.CommentList {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicMainCommentNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CommentList)
}

func (self *SquareClient2Hall) RpcGetDynamicMainCommentNew_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicMainCommentNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CommentList), e
}
func (self *SquareClient2Hall) RpcGetDynamicSecondaryCommentNew(reqMsg *client_server.IdInfo) *share_message.CommentList {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicSecondaryCommentNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CommentList)
}

func (self *SquareClient2Hall) RpcGetDynamicSecondaryCommentNew_(reqMsg *client_server.IdInfo) (*share_message.CommentList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDynamicSecondaryCommentNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CommentList), e
}
func (self *SquareClient2Hall) RpcGetSquareMessage(reqMsg *client_server.IdInfo) *MessageMainInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetSquareMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MessageMainInfo)
}

func (self *SquareClient2Hall) RpcGetSquareMessage_(reqMsg *client_server.IdInfo) (*MessageMainInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSquareMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MessageMainInfo), e
}
func (self *SquareClient2Hall) RpcGetPlayerZanInfo(reqMsg *client_server.RequestInfo) *ZanList {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerZanInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ZanList)
}

func (self *SquareClient2Hall) RpcGetPlayerZanInfo_(reqMsg *client_server.RequestInfo) (*ZanList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerZanInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ZanList), e
}
func (self *SquareClient2Hall) RpcGetPlayerAttentionInfo(reqMsg *client_server.RequestInfo) *AttentionList {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerAttentionInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AttentionList)
}

func (self *SquareClient2Hall) RpcGetPlayerAttentionInfo_(reqMsg *client_server.RequestInfo) (*AttentionList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerAttentionInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AttentionList), e
}
func (self *SquareClient2Hall) RpcDynamicTop(reqMsg *DynamicTopReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDynamicTop", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *SquareClient2Hall) RpcDynamicTop_(reqMsg *DynamicTopReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDynamicTop", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *SquareClient2Hall) RpcReadPlayerInfo(reqMsg *UnReadInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReadPlayerInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *SquareClient2Hall) RpcReadPlayerInfo_(reqMsg *UnReadInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReadPlayerInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *SquareClient2Hall) RpcFirstLoginSquare(reqMsg *FirstLoginSquareReq) *FirstLoginSquareReply {
	msg, e := self.Sender.CallRpcMethod("RpcFirstLoginSquare", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*FirstLoginSquareReply)
}

func (self *SquareClient2Hall) RpcFirstLoginSquare_(reqMsg *FirstLoginSquareReq) (*FirstLoginSquareReply, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcFirstLoginSquare", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*FirstLoginSquareReply), e
}
func (self *SquareClient2Hall) RpcAdvDetail(reqMsg *share_message.AdvSetting) *AdvDetailReply {
	msg, e := self.Sender.CallRpcMethod("RpcAdvDetail", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AdvDetailReply)
}

func (self *SquareClient2Hall) RpcAdvDetail_(reqMsg *share_message.AdvSetting) (*AdvDetailReply, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAdvDetail", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AdvDetailReply), e
}
func (self *SquareClient2Hall) RpcGetAllTopic(reqMsg *base.Empty) *AllTopic {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllTopic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllTopic)
}

func (self *SquareClient2Hall) RpcGetAllTopic_(reqMsg *base.Empty) (*AllTopic, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllTopic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllTopic), e
}
func (self *SquareClient2Hall) RpcGetTopicDetailReq(reqMsg *TopicDetailReq) *TopicDetailResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicDetailReq", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicDetailResp)
}

func (self *SquareClient2Hall) RpcGetTopicDetailReq_(reqMsg *TopicDetailReq) (*TopicDetailResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicDetailReq", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicDetailResp), e
}
func (self *SquareClient2Hall) RpcGetTopicMainPageList(reqMsg *TopicMainPageListReq) *TopicMainPageListResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicMainPageList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicMainPageListResp)
}

func (self *SquareClient2Hall) RpcGetTopicMainPageList_(reqMsg *TopicMainPageListReq) (*TopicMainPageListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicMainPageList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicMainPageListResp), e
}
func (self *SquareClient2Hall) RpcGetTopicParticipateList(reqMsg *TopicParticipateListReq) *TopicParticipateListResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicParticipateList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicParticipateListResp)
}

func (self *SquareClient2Hall) RpcGetTopicParticipateList_(reqMsg *TopicParticipateListReq) (*TopicParticipateListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicParticipateList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicParticipateListResp), e
}
func (self *SquareClient2Hall) RpcAttentionTopic(reqMsg *AttentionTopicReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAttentionTopic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *SquareClient2Hall) RpcAttentionTopic_(reqMsg *AttentionTopicReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAttentionTopic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *SquareClient2Hall) RpcMyAttentionTopicList(reqMsg *MyAttentionTopicListReq) *MyAttentionTopicListResp {
	msg, e := self.Sender.CallRpcMethod("RpcMyAttentionTopicList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MyAttentionTopicListResp)
}

func (self *SquareClient2Hall) RpcMyAttentionTopicList_(reqMsg *MyAttentionTopicListReq) (*MyAttentionTopicListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMyAttentionTopicList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MyAttentionTopicListResp), e
}
func (self *SquareClient2Hall) RpcSearchTopic(reqMsg *SearchTopicReq) *SearchTopicResp {
	msg, e := self.Sender.CallRpcMethod("RpcSearchTopic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SearchTopicResp)
}

func (self *SquareClient2Hall) RpcSearchTopic_(reqMsg *SearchTopicReq) (*SearchTopicResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSearchTopic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SearchTopicResp), e
}
func (self *SquareClient2Hall) RpcSearchHotTopic(reqMsg *base.Empty) *SearchHotTopicResp {
	msg, e := self.Sender.CallRpcMethod("RpcSearchHotTopic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SearchHotTopicResp)
}

func (self *SquareClient2Hall) RpcSearchHotTopic_(reqMsg *base.Empty) (*SearchHotTopicResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSearchHotTopic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SearchHotTopicResp), e
}
func (self *SquareClient2Hall) RpcGetTopicTypeList(reqMsg *base.Empty) *TopicTypeListResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicTypeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicTypeListResp)
}

func (self *SquareClient2Hall) RpcGetTopicTypeList_(reqMsg *base.Empty) (*TopicTypeListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicTypeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicTypeListResp), e
}
func (self *SquareClient2Hall) RpcGetTopicList(reqMsg *TopicListReq) *TopicListResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicListResp)
}

func (self *SquareClient2Hall) RpcGetTopicList_(reqMsg *TopicListReq) (*TopicListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicListResp), e
}
func (self *SquareClient2Hall) RpcFlushTopic(reqMsg *FlushTopicReq) *FlushTopicResp {
	msg, e := self.Sender.CallRpcMethod("RpcFlushTopic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*FlushTopicResp)
}

func (self *SquareClient2Hall) RpcFlushTopic_(reqMsg *FlushTopicReq) (*FlushTopicResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcFlushTopic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*FlushTopicResp), e
}
func (self *SquareClient2Hall) RpcHotTopicList(reqMsg *HotTopicListReq) *HotTopicListResp {
	msg, e := self.Sender.CallRpcMethod("RpcHotTopicList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*HotTopicListResp)
}

func (self *SquareClient2Hall) RpcHotTopicList_(reqMsg *HotTopicListReq) (*HotTopicListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcHotTopicList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*HotTopicListResp), e
}
func (self *SquareClient2Hall) RpcAttentionRecommendPlayer(reqMsg *base.Empty) *AttentionRecommendPlayerResp {
	msg, e := self.Sender.CallRpcMethod("RpcAttentionRecommendPlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AttentionRecommendPlayerResp)
}

func (self *SquareClient2Hall) RpcAttentionRecommendPlayer_(reqMsg *base.Empty) (*AttentionRecommendPlayerResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAttentionRecommendPlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AttentionRecommendPlayerResp), e
}
func (self *SquareClient2Hall) RpcSquareAttention(reqMsg *SquareAttentionReq) *SquareAttentionResp {
	msg, e := self.Sender.CallRpcMethod("RpcSquareAttention", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SquareAttentionResp)
}

func (self *SquareClient2Hall) RpcSquareAttention_(reqMsg *SquareAttentionReq) (*SquareAttentionResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSquareAttention", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SquareAttentionResp), e
}
func (self *SquareClient2Hall) RpcTopicHotDynamicParticipatePlayer(reqMsg *TopicParticipateListReq) *TopicParticipateListResp {
	msg, e := self.Sender.CallRpcMethod("RpcTopicHotDynamicParticipatePlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicParticipateListResp)
}

func (self *SquareClient2Hall) RpcTopicHotDynamicParticipatePlayer_(reqMsg *TopicParticipateListReq) (*TopicParticipateListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTopicHotDynamicParticipatePlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicParticipateListResp), e
}
func (self *SquareClient2Hall) RpcTopicHead(reqMsg *base.Empty) *TopicHeadResp {
	msg, e := self.Sender.CallRpcMethod("RpcTopicHead", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicHeadResp)
}

func (self *SquareClient2Hall) RpcTopicHead_(reqMsg *base.Empty) (*TopicHeadResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTopicHead", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicHeadResp), e
}
func (self *SquareClient2Hall) RpcTopicDevoteList(reqMsg *TopicDevoteListReq) *TopicDevoteListResp {
	msg, e := self.Sender.CallRpcMethod("RpcTopicDevoteList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicDevoteListResp)
}

func (self *SquareClient2Hall) RpcTopicDevoteList_(reqMsg *TopicDevoteListReq) (*TopicDevoteListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTopicDevoteList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicDevoteListResp), e
}
func (self *SquareClient2Hall) RpcTopicMasterCondition(reqMsg *TopicMasterConditionReq) *TopicMasterConditionResp {
	msg, e := self.Sender.CallRpcMethod("RpcTopicMasterCondition", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicMasterConditionResp)
}

func (self *SquareClient2Hall) RpcTopicMasterCondition_(reqMsg *TopicMasterConditionReq) (*TopicMasterConditionResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTopicMasterCondition", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicMasterConditionResp), e
}
func (self *SquareClient2Hall) RpcApplyTopicMaster(reqMsg *ApplyTopicMasterReq) *ApplyTopicMasterResp {
	msg, e := self.Sender.CallRpcMethod("RpcApplyTopicMaster", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ApplyTopicMasterResp)
}

func (self *SquareClient2Hall) RpcApplyTopicMaster_(reqMsg *ApplyTopicMasterReq) (*ApplyTopicMasterResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcApplyTopicMaster", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ApplyTopicMasterResp), e
}
func (self *SquareClient2Hall) RpcTopicMasterEdit(reqMsg *TopicMasterEditReq) *TopicMasterEditResp {
	msg, e := self.Sender.CallRpcMethod("RpcTopicMasterEdit", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicMasterEditResp)
}

func (self *SquareClient2Hall) RpcTopicMasterEdit_(reqMsg *TopicMasterEditReq) (*TopicMasterEditResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTopicMasterEdit", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicMasterEditResp), e
}
func (self *SquareClient2Hall) RpcTopicTop(reqMsg *TopicTopReq) *TopicTopResp {
	msg, e := self.Sender.CallRpcMethod("RpcTopicTop", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicTopResp)
}

func (self *SquareClient2Hall) RpcTopicTop_(reqMsg *TopicTopReq) (*TopicTopResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTopicTop", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicTopResp), e
}
func (self *SquareClient2Hall) RpcTopicTopCancel(reqMsg *TopicTopCancelReq) *TopicTopCancelResp {
	msg, e := self.Sender.CallRpcMethod("RpcTopicTopCancel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicTopCancelResp)
}

func (self *SquareClient2Hall) RpcTopicTopCancel_(reqMsg *TopicTopCancelReq) (*TopicTopCancelResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTopicTopCancel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicTopCancelResp), e
}
func (self *SquareClient2Hall) RpcTopicLeaderBoardDescription(reqMsg *base.Empty) *TopicLeaderBoardDescriptionResp {
	msg, e := self.Sender.CallRpcMethod("RpcTopicLeaderBoardDescription", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicLeaderBoardDescriptionResp)
}

func (self *SquareClient2Hall) RpcTopicLeaderBoardDescription_(reqMsg *base.Empty) (*TopicLeaderBoardDescriptionResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTopicLeaderBoardDescription", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicLeaderBoardDescriptionResp), e
}
func (self *SquareClient2Hall) RpcQuitTopicMaster(reqMsg *QuitTopicMasterReq) *QuitTopicMasterResp {
	msg, e := self.Sender.CallRpcMethod("RpcQuitTopicMaster", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QuitTopicMasterResp)
}

func (self *SquareClient2Hall) RpcQuitTopicMaster_(reqMsg *QuitTopicMasterReq) (*QuitTopicMasterResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQuitTopicMaster", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QuitTopicMasterResp), e
}
func (self *SquareClient2Hall) RpcTopicMasterDelDynamic(reqMsg *TopicMasterDelDynamicReq) *TopicMasterDelDynamicResp {
	msg, e := self.Sender.CallRpcMethod("RpcTopicMasterDelDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicMasterDelDynamicResp)
}

func (self *SquareClient2Hall) RpcTopicMasterDelDynamic_(reqMsg *TopicMasterDelDynamicReq) (*TopicMasterDelDynamicResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTopicMasterDelDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicMasterDelDynamicResp), e
}

// ==========================================================
type IHall2SquareClient interface {
	RpcSquareAllDynamic(reqMsg *AllInfo)
	RpcNewMessageForApi(reqMsg *NewUnReadMessageRespForApi)
	RpcNewVersionSquareAllDynamic(reqMsg *NewVersionAllInfo)
}

type Hall2SquareClient struct {
	Sender easygo.IMessageSender
}

func (self *Hall2SquareClient) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Hall2SquareClient) RpcSquareAllDynamic(reqMsg *AllInfo) {
	self.Sender.CallRpcMethod("RpcSquareAllDynamic", reqMsg)
}
func (self *Hall2SquareClient) RpcNewMessageForApi(reqMsg *NewUnReadMessageRespForApi) {
	self.Sender.CallRpcMethod("RpcNewMessageForApi", reqMsg)
}
func (self *Hall2SquareClient) RpcNewVersionSquareAllDynamic(reqMsg *NewVersionAllInfo) {
	self.Sender.CallRpcMethod("RpcNewVersionSquareAllDynamic", reqMsg)
}
