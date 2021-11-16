package client_hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type IClient2ESports interface {
	RpcESportGetHomeInfo(reqMsg *ESportInfoRequest) *ESportMenuHomeInfo
	RpcESportGetHomeInfo_(reqMsg *ESportInfoRequest) (*ESportMenuHomeInfo, easygo.IRpcInterrupt)
	RpcESportGetRealtimeList(reqMsg *ESportPageRequest) *ESportRealtimeListResult
	RpcESportGetRealtimeList_(reqMsg *ESportPageRequest) (*ESportRealtimeListResult, easygo.IRpcInterrupt)
	RpcESportGetRealtimeInfo(reqMsg *ESportInfoRequest) *ESportRealTimeResult
	RpcESportGetRealtimeInfo_(reqMsg *ESportInfoRequest) (*ESportRealTimeResult, easygo.IRpcInterrupt)
	RpcESportThumbsUp(reqMsg *ESportInfoRequest) *ESportThumbsUpResult
	RpcESportThumbsUp_(reqMsg *ESportInfoRequest) (*ESportThumbsUpResult, easygo.IRpcInterrupt)
	RpcESportSendComment(reqMsg *ESportCommentInfo) *ESportCommonResult
	RpcESportSendComment_(reqMsg *ESportCommentInfo) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESportDeleteComment(reqMsg *ESportDeleteCommentInfo) *ESportCommonResult
	RpcESportDeleteComment_(reqMsg *ESportDeleteCommentInfo) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESportGetComment(reqMsg *ESportCommentRequest) *ESportCommentReplyListResult
	RpcESportGetComment_(reqMsg *ESportCommentRequest) (*ESportCommentReplyListResult, easygo.IRpcInterrupt)
	RpcESportGetCommentReply(reqMsg *ESportCommentRequest) *ESportCommentReplyListResult
	RpcESportGetCommentReply_(reqMsg *ESportCommentRequest) (*ESportCommentReplyListResult, easygo.IRpcInterrupt)
	RpcESportGetAllLabelList(reqMsg *ESportInfoRequest) *ESportLabelList
	RpcESportGetAllLabelList_(reqMsg *ESportInfoRequest) (*ESportLabelList, easygo.IRpcInterrupt)
	RpcESportGetCarouselList(reqMsg *ESportInfoRequest) *ESportCarouselList
	RpcESportGetCarouselList_(reqMsg *ESportInfoRequest) (*ESportCarouselList, easygo.IRpcInterrupt)
	RpcESportSaveGameLabelConfig(reqMsg *ESportLabelList) *ESportCommonResult
	RpcESportSaveGameLabelConfig_(reqMsg *ESportLabelList) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESportGetSysMsgList(reqMsg *ESportInfoRequest) *ESPortsSysMsgList
	RpcESportGetSysMsgList_(reqMsg *ESportInfoRequest) (*ESPortsSysMsgList, easygo.IRpcInterrupt)
	RpcESportGetVideoList(reqMsg *ESportVideoPageRequest) *ESportVideoListResult
	RpcESportGetVideoList_(reqMsg *ESportVideoPageRequest) (*ESportVideoListResult, easygo.IRpcInterrupt)
	RpcESportGetVideoInfo(reqMsg *ESportVideoRequest) *ESportVideoResult
	RpcESportGetVideoInfo_(reqMsg *ESportVideoRequest) (*ESportVideoResult, easygo.IRpcInterrupt)
	RpcESportGetMyHistoryVideoList(reqMsg *ESportVideoPageRequest) *ESportVideoListResult
	RpcESportGetMyHistoryVideoList_(reqMsg *ESportVideoPageRequest) (*ESportVideoListResult, easygo.IRpcInterrupt)
	RpcESportGetLiveHomeInfo(reqMsg *ESportInfoRequest) *ESportLiveHomeInfo
	RpcESportGetLiveHomeInfo_(reqMsg *ESportInfoRequest) (*ESportLiveHomeInfo, easygo.IRpcInterrupt)
	RpcESportAddFollowLive(reqMsg *ESportInfoRequest) *ESportCommonResult
	RpcESportAddFollowLive_(reqMsg *ESportInfoRequest) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESportGetMyFollowLiveList(reqMsg *ESportPageRequest) *ESportVideoListResult
	RpcESportGetMyFollowLiveList_(reqMsg *ESportPageRequest) (*ESportVideoListResult, easygo.IRpcInterrupt)
	RpcESportApplyOpenLive(reqMsg *ESportMyLiveRoomInfo) *ESportCommonResult
	RpcESportApplyOpenLive_(reqMsg *ESportMyLiveRoomInfo) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESportSendLiveRoomMsg(reqMsg *ESportCommentInfo) *ESportCommonResult
	RpcESportSendLiveRoomMsg_(reqMsg *ESportCommentInfo) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESportLeaveLive(reqMsg *ESportCommonResult) *ESportCommonResult
	RpcESportLeaveLive_(reqMsg *ESportCommonResult) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESportGetESPortsGameViewList(reqMsg *ESportPageRequest) *ESPortsGameItemViewResult
	RpcESportGetESPortsGameViewList_(reqMsg *ESportPageRequest) (*ESPortsGameItemViewResult, easygo.IRpcInterrupt)
	RpcESportGetGameVideoList(reqMsg *ESportGameViewPageRequest) *ESportVideoListResult
	RpcESportGetGameVideoList_(reqMsg *ESportGameViewPageRequest) (*ESportVideoListResult, easygo.IRpcInterrupt)
	RpcESPortsBpsClick(reqMsg *ESPortsBpsClickRequest) *ESportCommonResult
	RpcESPortsBpsClick_(reqMsg *ESPortsBpsClickRequest) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESPortsBpsClickList(reqMsg *ESPortsBpsClickListRequest) *ESportCommonResult
	RpcESPortsBpsClickList_(reqMsg *ESPortsBpsClickListRequest) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESPortsBpsDuration(reqMsg *ESPortsBpsDurationRequest) *ESportCommonResult
	RpcESPortsBpsDuration_(reqMsg *ESPortsBpsDurationRequest) (*ESportCommonResult, easygo.IRpcInterrupt)
	RpcESPortsCoinView(reqMsg *base.Empty) *ESPortsCoinViewResult
	RpcESPortsCoinView_(reqMsg *base.Empty) (*ESPortsCoinViewResult, easygo.IRpcInterrupt)
	RpcESPortsCoinExChange(reqMsg *ESPortsCoinExChangeRequest) *ESPortsCoinExChangeResult
	RpcESPortsCoinExChange_(reqMsg *ESPortsCoinExChangeRequest) (*ESPortsCoinExChangeResult, easygo.IRpcInterrupt)
	RpcESPortsCoinExChangeRecord(reqMsg *ESPortsCoinExChangeRecordRequest) *ESPortsCoinExChangeRecordResult
	RpcESPortsCoinExChangeRecord_(reqMsg *ESPortsCoinExChangeRecordRequest) (*ESPortsCoinExChangeRecordResult, easygo.IRpcInterrupt)
	RpcESPortsApiOrigin(reqMsg *base.Empty) *RpcESPortsApiOriginResult
	RpcESPortsApiOrigin_(reqMsg *base.Empty) (*RpcESPortsApiOriginResult, easygo.IRpcInterrupt)
}

type Client2ESports struct {
	Sender easygo.IMessageSender
}

func (self *Client2ESports) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Client2ESports) RpcESportGetHomeInfo(reqMsg *ESportInfoRequest) *ESportMenuHomeInfo {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetHomeInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportMenuHomeInfo)
}

func (self *Client2ESports) RpcESportGetHomeInfo_(reqMsg *ESportInfoRequest) (*ESportMenuHomeInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetHomeInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportMenuHomeInfo), e
}
func (self *Client2ESports) RpcESportGetRealtimeList(reqMsg *ESportPageRequest) *ESportRealtimeListResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetRealtimeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportRealtimeListResult)
}

func (self *Client2ESports) RpcESportGetRealtimeList_(reqMsg *ESportPageRequest) (*ESportRealtimeListResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetRealtimeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportRealtimeListResult), e
}
func (self *Client2ESports) RpcESportGetRealtimeInfo(reqMsg *ESportInfoRequest) *ESportRealTimeResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetRealtimeInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportRealTimeResult)
}

func (self *Client2ESports) RpcESportGetRealtimeInfo_(reqMsg *ESportInfoRequest) (*ESportRealTimeResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetRealtimeInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportRealTimeResult), e
}
func (self *Client2ESports) RpcESportThumbsUp(reqMsg *ESportInfoRequest) *ESportThumbsUpResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportThumbsUp", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportThumbsUpResult)
}

func (self *Client2ESports) RpcESportThumbsUp_(reqMsg *ESportInfoRequest) (*ESportThumbsUpResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportThumbsUp", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportThumbsUpResult), e
}
func (self *Client2ESports) RpcESportSendComment(reqMsg *ESportCommentInfo) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportSendComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Client2ESports) RpcESportSendComment_(reqMsg *ESportCommentInfo) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportSendComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Client2ESports) RpcESportDeleteComment(reqMsg *ESportDeleteCommentInfo) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportDeleteComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Client2ESports) RpcESportDeleteComment_(reqMsg *ESportDeleteCommentInfo) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportDeleteComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Client2ESports) RpcESportGetComment(reqMsg *ESportCommentRequest) *ESportCommentReplyListResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommentReplyListResult)
}

func (self *Client2ESports) RpcESportGetComment_(reqMsg *ESportCommentRequest) (*ESportCommentReplyListResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommentReplyListResult), e
}
func (self *Client2ESports) RpcESportGetCommentReply(reqMsg *ESportCommentRequest) *ESportCommentReplyListResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetCommentReply", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommentReplyListResult)
}

func (self *Client2ESports) RpcESportGetCommentReply_(reqMsg *ESportCommentRequest) (*ESportCommentReplyListResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetCommentReply", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommentReplyListResult), e
}
func (self *Client2ESports) RpcESportGetAllLabelList(reqMsg *ESportInfoRequest) *ESportLabelList {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetAllLabelList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportLabelList)
}

func (self *Client2ESports) RpcESportGetAllLabelList_(reqMsg *ESportInfoRequest) (*ESportLabelList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetAllLabelList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportLabelList), e
}
func (self *Client2ESports) RpcESportGetCarouselList(reqMsg *ESportInfoRequest) *ESportCarouselList {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetCarouselList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCarouselList)
}

func (self *Client2ESports) RpcESportGetCarouselList_(reqMsg *ESportInfoRequest) (*ESportCarouselList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetCarouselList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCarouselList), e
}
func (self *Client2ESports) RpcESportSaveGameLabelConfig(reqMsg *ESportLabelList) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportSaveGameLabelConfig", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Client2ESports) RpcESportSaveGameLabelConfig_(reqMsg *ESportLabelList) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportSaveGameLabelConfig", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Client2ESports) RpcESportGetSysMsgList(reqMsg *ESportInfoRequest) *ESPortsSysMsgList {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetSysMsgList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESPortsSysMsgList)
}

func (self *Client2ESports) RpcESportGetSysMsgList_(reqMsg *ESportInfoRequest) (*ESPortsSysMsgList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetSysMsgList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESPortsSysMsgList), e
}
func (self *Client2ESports) RpcESportGetVideoList(reqMsg *ESportVideoPageRequest) *ESportVideoListResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetVideoList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportVideoListResult)
}

func (self *Client2ESports) RpcESportGetVideoList_(reqMsg *ESportVideoPageRequest) (*ESportVideoListResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetVideoList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportVideoListResult), e
}
func (self *Client2ESports) RpcESportGetVideoInfo(reqMsg *ESportVideoRequest) *ESportVideoResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetVideoInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportVideoResult)
}

func (self *Client2ESports) RpcESportGetVideoInfo_(reqMsg *ESportVideoRequest) (*ESportVideoResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetVideoInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportVideoResult), e
}
func (self *Client2ESports) RpcESportGetMyHistoryVideoList(reqMsg *ESportVideoPageRequest) *ESportVideoListResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetMyHistoryVideoList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportVideoListResult)
}

func (self *Client2ESports) RpcESportGetMyHistoryVideoList_(reqMsg *ESportVideoPageRequest) (*ESportVideoListResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetMyHistoryVideoList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportVideoListResult), e
}
func (self *Client2ESports) RpcESportGetLiveHomeInfo(reqMsg *ESportInfoRequest) *ESportLiveHomeInfo {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetLiveHomeInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportLiveHomeInfo)
}

func (self *Client2ESports) RpcESportGetLiveHomeInfo_(reqMsg *ESportInfoRequest) (*ESportLiveHomeInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetLiveHomeInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportLiveHomeInfo), e
}
func (self *Client2ESports) RpcESportAddFollowLive(reqMsg *ESportInfoRequest) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportAddFollowLive", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Client2ESports) RpcESportAddFollowLive_(reqMsg *ESportInfoRequest) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportAddFollowLive", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Client2ESports) RpcESportGetMyFollowLiveList(reqMsg *ESportPageRequest) *ESportVideoListResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetMyFollowLiveList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportVideoListResult)
}

func (self *Client2ESports) RpcESportGetMyFollowLiveList_(reqMsg *ESportPageRequest) (*ESportVideoListResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetMyFollowLiveList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportVideoListResult), e
}
func (self *Client2ESports) RpcESportApplyOpenLive(reqMsg *ESportMyLiveRoomInfo) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportApplyOpenLive", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Client2ESports) RpcESportApplyOpenLive_(reqMsg *ESportMyLiveRoomInfo) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportApplyOpenLive", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Client2ESports) RpcESportSendLiveRoomMsg(reqMsg *ESportCommentInfo) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportSendLiveRoomMsg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Client2ESports) RpcESportSendLiveRoomMsg_(reqMsg *ESportCommentInfo) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportSendLiveRoomMsg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Client2ESports) RpcESportLeaveLive(reqMsg *ESportCommonResult) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportLeaveLive", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Client2ESports) RpcESportLeaveLive_(reqMsg *ESportCommonResult) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportLeaveLive", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Client2ESports) RpcESportGetESPortsGameViewList(reqMsg *ESportPageRequest) *ESPortsGameItemViewResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetESPortsGameViewList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESPortsGameItemViewResult)
}

func (self *Client2ESports) RpcESportGetESPortsGameViewList_(reqMsg *ESportPageRequest) (*ESPortsGameItemViewResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetESPortsGameViewList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESPortsGameItemViewResult), e
}
func (self *Client2ESports) RpcESportGetGameVideoList(reqMsg *ESportGameViewPageRequest) *ESportVideoListResult {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetGameVideoList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportVideoListResult)
}

func (self *Client2ESports) RpcESportGetGameVideoList_(reqMsg *ESportGameViewPageRequest) (*ESportVideoListResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportGetGameVideoList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportVideoListResult), e
}
func (self *Client2ESports) RpcESPortsBpsClick(reqMsg *ESPortsBpsClickRequest) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsBpsClick", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Client2ESports) RpcESPortsBpsClick_(reqMsg *ESPortsBpsClickRequest) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsBpsClick", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Client2ESports) RpcESPortsBpsClickList(reqMsg *ESPortsBpsClickListRequest) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsBpsClickList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Client2ESports) RpcESPortsBpsClickList_(reqMsg *ESPortsBpsClickListRequest) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsBpsClickList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Client2ESports) RpcESPortsBpsDuration(reqMsg *ESPortsBpsDurationRequest) *ESportCommonResult {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsBpsDuration", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESportCommonResult)
}

func (self *Client2ESports) RpcESPortsBpsDuration_(reqMsg *ESPortsBpsDurationRequest) (*ESportCommonResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsBpsDuration", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESportCommonResult), e
}
func (self *Client2ESports) RpcESPortsCoinView(reqMsg *base.Empty) *ESPortsCoinViewResult {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsCoinView", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESPortsCoinViewResult)
}

func (self *Client2ESports) RpcESPortsCoinView_(reqMsg *base.Empty) (*ESPortsCoinViewResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsCoinView", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESPortsCoinViewResult), e
}
func (self *Client2ESports) RpcESPortsCoinExChange(reqMsg *ESPortsCoinExChangeRequest) *ESPortsCoinExChangeResult {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsCoinExChange", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESPortsCoinExChangeResult)
}

func (self *Client2ESports) RpcESPortsCoinExChange_(reqMsg *ESPortsCoinExChangeRequest) (*ESPortsCoinExChangeResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsCoinExChange", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESPortsCoinExChangeResult), e
}
func (self *Client2ESports) RpcESPortsCoinExChangeRecord(reqMsg *ESPortsCoinExChangeRecordRequest) *ESPortsCoinExChangeRecordResult {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsCoinExChangeRecord", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ESPortsCoinExChangeRecordResult)
}

func (self *Client2ESports) RpcESPortsCoinExChangeRecord_(reqMsg *ESPortsCoinExChangeRecordRequest) (*ESPortsCoinExChangeRecordResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsCoinExChangeRecord", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ESPortsCoinExChangeRecordResult), e
}
func (self *Client2ESports) RpcESPortsApiOrigin(reqMsg *base.Empty) *RpcESPortsApiOriginResult {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsApiOrigin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*RpcESPortsApiOriginResult)
}

func (self *Client2ESports) RpcESPortsApiOrigin_(reqMsg *base.Empty) (*RpcESPortsApiOriginResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESPortsApiOrigin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*RpcESPortsApiOriginResult), e
}

// ==========================================================
type IESports2Client interface {
	RpcESportNewSysMessage(reqMsg *ESPortsSysMsgList)
	RpcESportNewRoomMsg(reqMsg *share_message.TableESPortsLiveRoomMsgLog)
	RpcESportDataStatusInfo(reqMsg *ESportDataStatusInfo)
}

type ESports2Client struct {
	Sender easygo.IMessageSender
}

func (self *ESports2Client) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *ESports2Client) RpcESportNewSysMessage(reqMsg *ESPortsSysMsgList) {
	self.Sender.CallRpcMethod("RpcESportNewSysMessage", reqMsg)
}
func (self *ESports2Client) RpcESportNewRoomMsg(reqMsg *share_message.TableESPortsLiveRoomMsgLog) {
	self.Sender.CallRpcMethod("RpcESportNewRoomMsg", reqMsg)
}
func (self *ESports2Client) RpcESportDataStatusInfo(reqMsg *ESportDataStatusInfo) {
	self.Sender.CallRpcMethod("RpcESportDataStatusInfo", reqMsg)
}
