package client_hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type IAccountClient2Hall interface {
	RpcLogin(reqMsg *LoginMsg) *LoginMsg
	RpcLogin_(reqMsg *LoginMsg) (*LoginMsg, easygo.IRpcInterrupt)
	RpcModifyPlayerMsg(reqMsg *client_server.ChangePlayerInfo) *base.Empty
	RpcModifyPlayerMsg_(reqMsg *client_server.ChangePlayerInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddFriend(reqMsg *AddPlayerInfo) *base.Empty
	RpcAddFriend_(reqMsg *AddPlayerInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcAgreeFriend(reqMsg *AccountInfo) *base.Empty
	RpcAgreeFriend_(reqMsg *AccountInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcFindFriend(reqMsg *AccountInfo) *base.Empty
	RpcFindFriend_(reqMsg *AccountInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcCreateTeam(reqMsg *CreateTeam) *base.Empty
	RpcCreateTeam_(reqMsg *CreateTeam) (*base.Empty, easygo.IRpcInterrupt)
	RpcCreateTopicTeam(reqMsg *CreateTeam) *base.Empty
	RpcCreateTopicTeam_(reqMsg *CreateTeam) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetFriendRequest(reqMsg *base.Empty) *share_message.AllAddPlayerMsg
	RpcGetFriendRequest_(reqMsg *base.Empty) (*share_message.AllAddPlayerMsg, easygo.IRpcInterrupt)
	RpcNewVersionGetFriendRequest(reqMsg *base.Empty) *share_message.AllAddPlayerMsg
	RpcNewVersionGetFriendRequest_(reqMsg *base.Empty) (*share_message.AllAddPlayerMsg, easygo.IRpcInterrupt)
	RpcNewVersionGetFriendNumRequest(reqMsg *base.Empty) *FriendNum
	RpcNewVersionGetFriendNumRequest_(reqMsg *base.Empty) (*FriendNum, easygo.IRpcInterrupt)
	RpcReadFriendRequest(reqMsg *client_server.ReadInfo) *base.Empty
	RpcReadFriendRequest_(reqMsg *client_server.ReadInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcSearchAddressBook(reqMsg *AddBookInfo) *client_server.AllPlayerInfo
	RpcSearchAddressBook_(reqMsg *AddBookInfo) (*client_server.AllPlayerInfo, easygo.IRpcInterrupt)
	RpcSetPassword(reqMsg *client_server.PasswordInfo) *base.Empty
	RpcSetPassword_(reqMsg *client_server.PasswordInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcChangePassword(reqMsg *client_server.PasswordInfo) *base.Empty
	RpcChangePassword_(reqMsg *client_server.PasswordInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcForgetPayPassword(reqMsg *client_server.PasswordInfo) *base.Empty
	RpcForgetPayPassword_(reqMsg *client_server.PasswordInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcClientGetCode(reqMsg *client_server.GetCodeRequest) *base.Empty
	RpcClientGetCode_(reqMsg *client_server.GetCodeRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcCheckMessageCode(reqMsg *client_server.CodeResponse) *base.Empty
	RpcCheckMessageCode_(reqMsg *client_server.CodeResponse) (*base.Empty, easygo.IRpcInterrupt)
	RpcClosePassword(reqMsg *base.Empty) *base.Empty
	RpcClosePassword_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcCheckPlayerPeopleId(reqMsg *PeopleIdInfo) *base.Empty
	RpcCheckPlayerPeopleId_(reqMsg *PeopleIdInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcSetPeopleAuth(reqMsg *PeopleIdInfo) *base.Empty
	RpcSetPeopleAuth_(reqMsg *PeopleIdInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcActiveAddTeamMember(reqMsg *client_server.TeamReq) *base.Empty
	RpcActiveAddTeamMember_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcLogOut(reqMsg *base.Empty) *base.Empty
	RpcLogOut_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcCheckPassword(reqMsg *client_server.PasswordInfo) *base.Empty
	RpcCheckPassword_(reqMsg *client_server.PasswordInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcFrezzeAccount(reqMsg *AccountInfo) *base.Empty
	RpcFrezzeAccount_(reqMsg *AccountInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcUnFrezzeAccount(reqMsg *AccountInfo) *base.Empty
	RpcUnFrezzeAccount_(reqMsg *AccountInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetVersion(reqMsg *base.Empty) *VersionInfo
	RpcGetVersion_(reqMsg *base.Empty) (*VersionInfo, easygo.IRpcInterrupt)
	RpcGetBlackList(reqMsg *base.Empty) *client_server.AllPlayerInfo
	RpcGetBlackList_(reqMsg *base.Empty) (*client_server.AllPlayerInfo, easygo.IRpcInterrupt)
	RpcBlackOperate(reqMsg *BlackInfo) *BlackInfo
	RpcBlackOperate_(reqMsg *BlackInfo) (*BlackInfo, easygo.IRpcInterrupt)
	RpcDelFriend(reqMsg *DelFriendInfo) *DelFriendInfo
	RpcDelFriend_(reqMsg *DelFriendInfo) (*DelFriendInfo, easygo.IRpcInterrupt)
	RpcGetCLAssistant(reqMsg *AInfo) *client_server.AssistantInfo
	RpcGetCLAssistant_(reqMsg *AInfo) (*client_server.AssistantInfo, easygo.IRpcInterrupt)
	RpcGetMoneyAssistant(reqMsg *PageInfo) *AllMoneyAssistantInfo
	RpcGetMoneyAssistant_(reqMsg *PageInfo) (*AllMoneyAssistantInfo, easygo.IRpcInterrupt)
	RpcGetLocationInfo(reqMsg *LocationInfo) *AllLocationPlayerInfo
	RpcGetLocationInfo_(reqMsg *LocationInfo) (*AllLocationPlayerInfo, easygo.IRpcInterrupt)
	RpcGetAllNearByInfo(reqMsg *base.Empty) *AllNearByMessage
	RpcGetAllNearByInfo_(reqMsg *base.Empty) (*AllNearByMessage, easygo.IRpcInterrupt)
	RpcGetNearByInfo(reqMsg *base.Empty) *AllNearByMessage
	RpcGetNearByInfo_(reqMsg *base.Empty) (*AllNearByMessage, easygo.IRpcInterrupt)
	RpcAddNearByMessage(reqMsg *NearByMessage) *base.Empty
	RpcAddNearByMessage_(reqMsg *NearByMessage) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteNearByInfo(reqMsg *base.Empty) *base.Empty
	RpcDeleteNearByInfo_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteMyNearByInfo(reqMsg *base.Empty) *base.Empty
	RpcDeleteMyNearByInfo_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcBindBankCode(reqMsg *BankMessage) *base.Empty
	RpcBindBankCode_(reqMsg *BankMessage) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddBank(reqMsg *BankMessage) *base.Empty
	RpcAddBank_(reqMsg *BankMessage) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelBank(reqMsg *BankMessage) *base.Empty
	RpcDelBank_(reqMsg *BankMessage) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetBankCode(reqMsg *BankMessage) *BankMessage
	RpcGetBankCode_(reqMsg *BankMessage) (*BankMessage, easygo.IRpcInterrupt)
	RpcSyncTime(reqMsg *base.Empty) *client_server.NTP
	RpcSyncTime_(reqMsg *base.Empty) (*client_server.NTP, easygo.IRpcInterrupt)
	RpcGetMoneyOrderInfo(reqMsg *MoneyType) *AllOrderInfo
	RpcGetMoneyOrderInfo_(reqMsg *MoneyType) (*AllOrderInfo, easygo.IRpcInterrupt)
	RpcGetCashInfo(reqMsg *PageInfo) *AllCashInfo
	RpcGetCashInfo_(reqMsg *PageInfo) (*AllCashInfo, easygo.IRpcInterrupt)
	RpcGetRedPacketInfo(reqMsg *RedPacketInfo) *AllRedPacketInfo
	RpcGetRedPacketInfo_(reqMsg *RedPacketInfo) (*AllRedPacketInfo, easygo.IRpcInterrupt)
	RpcClearLocalTime(reqMsg *base.Empty)
	RpcRechargeMoney(reqMsg *share_message.PayOrderInfo) *share_message.PayOrderResult
	RpcRechargeMoney_(reqMsg *share_message.PayOrderInfo) (*share_message.PayOrderResult, easygo.IRpcInterrupt)
	RpcAddComplaintInfo(reqMsg *share_message.PlayerComplaint) *base.Empty
	RpcAddComplaintInfo_(reqMsg *share_message.PlayerComplaint) (*base.Empty, easygo.IRpcInterrupt)
	RpcWithdrawRequest(reqMsg *WithdrawInfo) *base.Empty
	RpcWithdrawRequest_(reqMsg *WithdrawInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcCheckPeopleIdValid(reqMsg *PeopleIdInfo) *CheckPeopleInfo
	RpcCheckPeopleIdValid_(reqMsg *PeopleIdInfo) (*CheckPeopleInfo, easygo.IRpcInterrupt)
	RpcFirstLoginSetInfo(reqMsg *FirstInfo) *FirstReturnInfo
	RpcFirstLoginSetInfo_(reqMsg *FirstInfo) (*FirstReturnInfo, easygo.IRpcInterrupt)
	RpcAliPay(reqMsg *share_message.AliPayData) *base.Empty
	RpcAliPay_(reqMsg *share_message.AliPayData) (*base.Empty, easygo.IRpcInterrupt)
	RpcPayForCode(reqMsg *PayInfo) *PayForCodeInfo
	RpcPayForCode_(reqMsg *PayInfo) (*PayForCodeInfo, easygo.IRpcInterrupt)
	RpcGetCollectIndexList(reqMsg *base.Empty) *CollectIndex
	RpcGetCollectIndexList_(reqMsg *base.Empty) (*CollectIndex, easygo.IRpcInterrupt)
	RpcGetCollectIndexInfo(reqMsg *CollectIndex) *AllCollectInfo
	RpcGetCollectIndexInfo_(reqMsg *CollectIndex) (*AllCollectInfo, easygo.IRpcInterrupt)
	RpcAddCollectInfo(reqMsg *share_message.CollectInfo) *share_message.CollectInfo
	RpcAddCollectInfo_(reqMsg *share_message.CollectInfo) (*share_message.CollectInfo, easygo.IRpcInterrupt)
	RpcGetCollectInfo(reqMsg *GetCollectInfo) *AllCollectInfo
	RpcGetCollectInfo_(reqMsg *GetCollectInfo) (*AllCollectInfo, easygo.IRpcInterrupt)
	RpcDelCollectInfo(reqMsg *DelCollectInfo) *base.Empty
	RpcDelCollectInfo_(reqMsg *DelCollectInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcSearchCollectInfo(reqMsg *SearchCollectInfo) *AllCollectInfo
	RpcSearchCollectInfo_(reqMsg *SearchCollectInfo) (*AllCollectInfo, easygo.IRpcInterrupt)
	RpcCheckout(reqMsg *share_message.OrderID) *base.Empty
	RpcCheckout_(reqMsg *share_message.OrderID) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetCodeInfo(reqMsg *CodeInfo) *CodeInfo
	RpcGetCodeInfo_(reqMsg *CodeInfo) (*CodeInfo, easygo.IRpcInterrupt)
	RpcCheckUnGetMoney(reqMsg *UnGetMoneyInfo) *UnGetMoneyInfo
	RpcCheckUnGetMoney_(reqMsg *UnGetMoneyInfo) (*UnGetMoneyInfo, easygo.IRpcInterrupt)
	RpcCutBackstage(reqMsg *base.Empty) *base.Empty
	RpcCutBackstage_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcReturnApp(reqMsg *base.Empty) *base.Empty
	RpcReturnApp_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcReqPaySMS(reqMsg *share_message.PayOrderInfo) *base.Empty
	RpcReqPaySMS_(reqMsg *share_message.PayOrderInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcReqSMSPay(reqMsg *BankPaySMS) *base.Empty
	RpcReqSMSPay_(reqMsg *BankPaySMS) (*base.Empty, easygo.IRpcInterrupt)
	RpcRecommendRequest(reqMsg *RecommendMsg) *base.Empty
	RpcRecommendRequest_(reqMsg *RecommendMsg) (*base.Empty, easygo.IRpcInterrupt)
	RpcFreshRecommendInfo(reqMsg *RecommendRefreshInfo) *client_server.RecommendInfo
	RpcFreshRecommendInfo_(reqMsg *RecommendRefreshInfo) (*client_server.RecommendInfo, easygo.IRpcInterrupt)
	RpcRequestCallInfo(reqMsg *base.Empty) *share_message.CallInfo
	RpcRequestCallInfo_(reqMsg *base.Empty) (*share_message.CallInfo, easygo.IRpcInterrupt)
	RpcCheckAccountVaild(reqMsg *client_server.CheckInfo) *client_server.CheckInfo
	RpcCheckAccountVaild_(reqMsg *client_server.CheckInfo) (*client_server.CheckInfo, easygo.IRpcInterrupt)
	RpcUserReturnApp(reqMsg *base.Empty) *client_server.TweetsListResponse
	RpcUserReturnApp_(reqMsg *base.Empty) (*client_server.TweetsListResponse, easygo.IRpcInterrupt)
	RpcOperateBindWechat(reqMsg *WechatInfo) *base.Empty
	RpcOperateBindWechat_(reqMsg *WechatInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcOperateBindPhone(reqMsg *WechatInfo) *base.Empty
	RpcOperateBindPhone_(reqMsg *WechatInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelTweets(reqMsg *client_server.TweetsIdsRequest) *base.Empty
	RpcDelTweets_(reqMsg *client_server.TweetsIdsRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcReturnTweets(reqMsg *base.Empty) *client_server.TweetsListResponse
	RpcReturnTweets_(reqMsg *base.Empty) (*client_server.TweetsListResponse, easygo.IRpcInterrupt)
	RpcModifyEmoticon(reqMsg *share_message.PlayerEmoticon) *base.Empty
	RpcModifyEmoticon_(reqMsg *share_message.PlayerEmoticon) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelEmoticon(reqMsg *share_message.PlayerEmoticon) *base.Empty
	RpcDelEmoticon_(reqMsg *share_message.PlayerEmoticon) (*base.Empty, easygo.IRpcInterrupt)
	RpcMarkPlayer(reqMsg *MarkName) *base.Empty
	RpcMarkPlayer_(reqMsg *MarkName) (*base.Empty, easygo.IRpcInterrupt)
	RpcSumClicksJumps(reqMsg *ArticleOptRequest) *base.Empty
	RpcSumClicksJumps_(reqMsg *ArticleOptRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcNoticeClicks(reqMsg *NoticeOptRequest) *base.Empty
	RpcNoticeClicks_(reqMsg *NoticeOptRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcOpenMyMainPage(reqMsg *OpenMyMainPageByPageReq) *client_server.MyMainPageInfo
	RpcOpenMyMainPage_(reqMsg *OpenMyMainPageByPageReq) (*client_server.MyMainPageInfo, easygo.IRpcInterrupt)
	RpcOpenMyMainPageByPage(reqMsg *OpenMyMainPageByPageReq) *client_server.MyMainPageInfo
	RpcOpenMyMainPageByPage_(reqMsg *OpenMyMainPageByPageReq) (*client_server.MyMainPageInfo, easygo.IRpcInterrupt)
	RpcGetMyFansAttentionInfo(reqMsg *MainInfo) *AllFansInfo
	RpcGetMyFansAttentionInfo_(reqMsg *MainInfo) (*AllFansInfo, easygo.IRpcInterrupt)
	RpcGetSomeDynamic(reqMsg *DynamicIdInfo) *DynamicInfo
	RpcGetSomeDynamic_(reqMsg *DynamicIdInfo) (*DynamicInfo, easygo.IRpcInterrupt)
	RpcGetPlayerOtherData(reqMsg *base.Empty) *client_server.AllPlayerMsg
	RpcGetPlayerOtherData_(reqMsg *base.Empty) (*client_server.AllPlayerMsg, easygo.IRpcInterrupt)
	RpcGetTeamMembers(reqMsg *TeamMembers) *TeamMembers
	RpcGetTeamMembers_(reqMsg *TeamMembers) (*TeamMembers, easygo.IRpcInterrupt)
	RpcCheckDirtyWord(reqMsg *CheckDirtyWord) *CheckDirtyWord
	RpcCheckDirtyWord_(reqMsg *CheckDirtyWord) (*CheckDirtyWord, easygo.IRpcInterrupt)
	RpcCheckCancelAccount(reqMsg *base.Empty) *CheckCancelAccount
	RpcCheckCancelAccount_(reqMsg *base.Empty) (*CheckCancelAccount, easygo.IRpcInterrupt)
	RpcSubmitCancelAccount(reqMsg *CancelAccountData) *base.Empty
	RpcSubmitCancelAccount_(reqMsg *CancelAccountData) (*base.Empty, easygo.IRpcInterrupt)
	RpcAttentioPlayer(reqMsg *client_server.AttenInfo) *base.Empty
	RpcAttentioPlayer_(reqMsg *client_server.AttenInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcESportAttentioPlayer(reqMsg *client_server.AttenInfo) *base.Empty
	RpcESportAttentioPlayer_(reqMsg *client_server.AttenInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcZanOperateSquareDynamic(reqMsg *client_server.ZanInfo) *base.Empty
	RpcZanOperateSquareDynamic_(reqMsg *client_server.ZanInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddCommentSquareDynamic(reqMsg *share_message.CommentData) *base.Empty
	RpcAddCommentSquareDynamic_(reqMsg *share_message.CommentData) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelCommentSquareDynamic(reqMsg *client_server.IdInfo) *base.Empty
	RpcDelCommentSquareDynamic_(reqMsg *client_server.IdInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddSquareCommentZan(reqMsg *share_message.CommentDataZan) *base.Empty
	RpcAddSquareCommentZan_(reqMsg *share_message.CommentDataZan) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelSquareDynamic(reqMsg *client_server.RequestInfo) *base.Empty
	RpcDelSquareDynamic_(reqMsg *client_server.RequestInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcShopOrderComplaintDetail(reqMsg *ComplaintID) *ShopOrderComplaintDetailRsp
	RpcShopOrderComplaintDetail_(reqMsg *ComplaintID) (*ShopOrderComplaintDetailRsp, easygo.IRpcInterrupt)
	RpcGetAdvData(reqMsg *AdvSettingRequest) *AdvSettingResponse
	RpcGetAdvData_(reqMsg *AdvSettingRequest) (*AdvSettingResponse, easygo.IRpcInterrupt)
	RpcGetNearByInfo2(reqMsg *LocationInfo) *NearByInfoReply
	RpcGetNearByInfo2_(reqMsg *LocationInfo) (*NearByInfoReply, easygo.IRpcInterrupt)
	RpcGetNearInfoByPage(reqMsg *LocationInfoByPage) *NearByInfoReplyByPage
	RpcGetNearInfoByPage_(reqMsg *LocationInfoByPage) (*NearByInfoReplyByPage, easygo.IRpcInterrupt)
	RpcAddAdvLog(reqMsg *share_message.AdvLogReq) *base.Empty
	RpcAddAdvLog_(reqMsg *share_message.AdvLogReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelNewFriendList(reqMsg *DelNewFriendListReq) *base.Empty
	RpcDelNewFriendList_(reqMsg *DelNewFriendListReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcNewTeamSetting(reqMsg *NewTeamSettingReq) *TeamSettingNotify
	RpcNewTeamSetting_(reqMsg *NewTeamSettingReq) (*TeamSettingNotify, easygo.IRpcInterrupt)
	RpcLocationInfoNew(reqMsg *LocationInfoNewReq) *LocationInfoNewResp
	RpcLocationInfoNew_(reqMsg *LocationInfoNewReq) (*LocationInfoNewResp, easygo.IRpcInterrupt)
	RpcNearRecommend(reqMsg *NearRecommendReq) *NearRecommendResp
	RpcNearRecommend_(reqMsg *NearRecommendReq) (*NearRecommendResp, easygo.IRpcInterrupt)
	RpcNearSayMessage(reqMsg *share_message.NearSessionList) *base.Empty
	RpcNearSayMessage_(reqMsg *share_message.NearSessionList) (*base.Empty, easygo.IRpcInterrupt)
	RpcNearSessionList(reqMsg *NearSessionListReq) *NearSessionListResp
	RpcNearSessionList_(reqMsg *NearSessionListReq) (*NearSessionListResp, easygo.IRpcInterrupt)
	RpcSendPlayerMessageList(reqMsg *SendPlayerMessageListReq) *NearSessionListResp
	RpcSendPlayerMessageList_(reqMsg *SendPlayerMessageListReq) (*NearSessionListResp, easygo.IRpcInterrupt)
	RpcGetNearChatList(reqMsg *GetNearChatListReq) *NearSessionListResp
	RpcGetNearChatList_(reqMsg *GetNearChatListReq) (*NearSessionListResp, easygo.IRpcInterrupt)
	RpcDelNearMessage(reqMsg *DelNearMessageReq) *base.Empty
	RpcDelNearMessage_(reqMsg *DelNearMessageReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcUpdateIsReadReq(reqMsg *UpdateIsReadReq) *base.Empty
	RpcUpdateIsReadReq_(reqMsg *UpdateIsReadReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcSetYoungPassWord(reqMsg *YoungPassWord) *base.Empty
	RpcSetYoungPassWord_(reqMsg *YoungPassWord) (*base.Empty, easygo.IRpcInterrupt)
	RpcNearAddAdvLog(reqMsg *share_message.AdvLogReq) *base.Empty
	RpcNearAddAdvLog_(reqMsg *share_message.AdvLogReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetPlayerOtherDataNew(reqMsg *base.Empty) *client_server.AllPlayerMsg
	RpcGetPlayerOtherDataNew_(reqMsg *base.Empty) (*client_server.AllPlayerMsg, easygo.IRpcInterrupt)
	RpcGetPlayerFriends(reqMsg *base.Empty) *client_server.AllPlayerMsg
	RpcGetPlayerFriends_(reqMsg *base.Empty) (*client_server.AllPlayerMsg, easygo.IRpcInterrupt)
	RpcGetPlayerNewFriends(reqMsg *GetNewFriends) *client_server.NewFriends
	RpcGetPlayerNewFriends_(reqMsg *GetNewFriends) (*client_server.NewFriends, easygo.IRpcInterrupt)
	RpcGetTeamMemberChange(reqMsg *TeamMembersChange) *TeamMembersChange
	RpcGetTeamMemberChange_(reqMsg *TeamMembersChange) (*TeamMembersChange, easygo.IRpcInterrupt)
	RpcGetSupportBankList(reqMsg *SupportBankList) *SupportBankList
	RpcGetSupportBankList_(reqMsg *SupportBankList) (*SupportBankList, easygo.IRpcInterrupt)
	RpcGetOnLineNum(reqMsg *base.Empty) *OnLineNum
	RpcGetOnLineNum_(reqMsg *base.Empty) (*OnLineNum, easygo.IRpcInterrupt)
	RpcGetMsgPageAdvList(reqMsg *base.Empty) *MsgAdv
	RpcGetMsgPageAdvList_(reqMsg *base.Empty) (*MsgAdv, easygo.IRpcInterrupt)
	RpcGetDiamond(reqMsg *base.Empty) *PlayerDiamond
	RpcGetDiamond_(reqMsg *base.Empty) (*PlayerDiamond, easygo.IRpcInterrupt)
	RpcGetSupportPayChannel(reqMsg *PayChannels) *PayChannels
	RpcGetSupportPayChannel_(reqMsg *PayChannels) (*PayChannels, easygo.IRpcInterrupt)
	RpcGetAllLabelMsg(reqMsg *base.Empty) *client_server.LabelMsg
	RpcGetAllLabelMsg_(reqMsg *base.Empty) (*client_server.LabelMsg, easygo.IRpcInterrupt)
	RpcSetPlayerLabel(reqMsg *FirstInfo) *base.Empty
	RpcSetPlayerLabel_(reqMsg *FirstInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcTopicTeamDynamicList(reqMsg *TopicTeamDynamicReq) *TopicTeamDynamicResp
	RpcTopicTeamDynamicList_(reqMsg *TopicTeamDynamicReq) (*TopicTeamDynamicResp, easygo.IRpcInterrupt)
	RpcGetTopicTeams(reqMsg *TopicTeams) *TopicTeams
	RpcGetTopicTeams_(reqMsg *TopicTeams) (*TopicTeams, easygo.IRpcInterrupt)
	RpcDefunctTeam(reqMsg *DefunctTeam) *DefunctTeam
	RpcDefunctTeam_(reqMsg *DefunctTeam) (*DefunctTeam, easygo.IRpcInterrupt)
	RpcGetAllMainMenu(reqMsg *base.Empty) *AllMainMenu
	RpcGetAllMainMenu_(reqMsg *base.Empty) (*AllMainMenu, easygo.IRpcInterrupt)
	RpcGetAllTipAdvs(reqMsg *base.Empty) *AllTipAdv
	RpcGetAllTipAdvs_(reqMsg *base.Empty) (*AllTipAdv, easygo.IRpcInterrupt)
	RpcGetStartPageAdvList(reqMsg *base.Empty) *MsgAdv
	RpcGetStartPageAdvList_(reqMsg *base.Empty) (*MsgAdv, easygo.IRpcInterrupt)
}

type AccountClient2Hall struct {
	Sender easygo.IMessageSender
}

func (self *AccountClient2Hall) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *AccountClient2Hall) RpcLogin(reqMsg *LoginMsg) *LoginMsg {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LoginMsg)
}

func (self *AccountClient2Hall) RpcLogin_(reqMsg *LoginMsg) (*LoginMsg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LoginMsg), e
}
func (self *AccountClient2Hall) RpcModifyPlayerMsg(reqMsg *client_server.ChangePlayerInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcModifyPlayerMsg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcModifyPlayerMsg_(reqMsg *client_server.ChangePlayerInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcModifyPlayerMsg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcAddFriend(reqMsg *AddPlayerInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddFriend", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcAddFriend_(reqMsg *AddPlayerInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddFriend", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcAgreeFriend(reqMsg *AccountInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAgreeFriend", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcAgreeFriend_(reqMsg *AccountInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAgreeFriend", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcFindFriend(reqMsg *AccountInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcFindFriend", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcFindFriend_(reqMsg *AccountInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcFindFriend", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcCreateTeam(reqMsg *CreateTeam) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCreateTeam", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcCreateTeam_(reqMsg *CreateTeam) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCreateTeam", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcCreateTopicTeam(reqMsg *CreateTeam) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCreateTopicTeam", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcCreateTopicTeam_(reqMsg *CreateTeam) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCreateTopicTeam", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcGetFriendRequest(reqMsg *base.Empty) *share_message.AllAddPlayerMsg {
	msg, e := self.Sender.CallRpcMethod("RpcGetFriendRequest", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.AllAddPlayerMsg)
}

func (self *AccountClient2Hall) RpcGetFriendRequest_(reqMsg *base.Empty) (*share_message.AllAddPlayerMsg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetFriendRequest", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.AllAddPlayerMsg), e
}
func (self *AccountClient2Hall) RpcNewVersionGetFriendRequest(reqMsg *base.Empty) *share_message.AllAddPlayerMsg {
	msg, e := self.Sender.CallRpcMethod("RpcNewVersionGetFriendRequest", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.AllAddPlayerMsg)
}

func (self *AccountClient2Hall) RpcNewVersionGetFriendRequest_(reqMsg *base.Empty) (*share_message.AllAddPlayerMsg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNewVersionGetFriendRequest", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.AllAddPlayerMsg), e
}
func (self *AccountClient2Hall) RpcNewVersionGetFriendNumRequest(reqMsg *base.Empty) *FriendNum {
	msg, e := self.Sender.CallRpcMethod("RpcNewVersionGetFriendNumRequest", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*FriendNum)
}

func (self *AccountClient2Hall) RpcNewVersionGetFriendNumRequest_(reqMsg *base.Empty) (*FriendNum, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNewVersionGetFriendNumRequest", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*FriendNum), e
}
func (self *AccountClient2Hall) RpcReadFriendRequest(reqMsg *client_server.ReadInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReadFriendRequest", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcReadFriendRequest_(reqMsg *client_server.ReadInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReadFriendRequest", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcSearchAddressBook(reqMsg *AddBookInfo) *client_server.AllPlayerInfo {
	msg, e := self.Sender.CallRpcMethod("RpcSearchAddressBook", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.AllPlayerInfo)
}

func (self *AccountClient2Hall) RpcSearchAddressBook_(reqMsg *AddBookInfo) (*client_server.AllPlayerInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSearchAddressBook", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.AllPlayerInfo), e
}
func (self *AccountClient2Hall) RpcSetPassword(reqMsg *client_server.PasswordInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSetPassword", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcSetPassword_(reqMsg *client_server.PasswordInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSetPassword", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcChangePassword(reqMsg *client_server.PasswordInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcChangePassword", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcChangePassword_(reqMsg *client_server.PasswordInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcChangePassword", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcForgetPayPassword(reqMsg *client_server.PasswordInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcForgetPayPassword", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcForgetPayPassword_(reqMsg *client_server.PasswordInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcForgetPayPassword", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcClientGetCode(reqMsg *client_server.GetCodeRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcClientGetCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcClientGetCode_(reqMsg *client_server.GetCodeRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcClientGetCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcCheckMessageCode(reqMsg *client_server.CodeResponse) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCheckMessageCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcCheckMessageCode_(reqMsg *client_server.CodeResponse) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckMessageCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcClosePassword(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcClosePassword", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcClosePassword_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcClosePassword", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcCheckPlayerPeopleId(reqMsg *PeopleIdInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCheckPlayerPeopleId", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcCheckPlayerPeopleId_(reqMsg *PeopleIdInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckPlayerPeopleId", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcSetPeopleAuth(reqMsg *PeopleIdInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSetPeopleAuth", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcSetPeopleAuth_(reqMsg *PeopleIdInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSetPeopleAuth", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcActiveAddTeamMember(reqMsg *client_server.TeamReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcActiveAddTeamMember", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcActiveAddTeamMember_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcActiveAddTeamMember", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcLogOut(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcLogOut", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcLogOut_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLogOut", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcCheckPassword(reqMsg *client_server.PasswordInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCheckPassword", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcCheckPassword_(reqMsg *client_server.PasswordInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckPassword", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcFrezzeAccount(reqMsg *AccountInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcFrezzeAccount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcFrezzeAccount_(reqMsg *AccountInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcFrezzeAccount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcUnFrezzeAccount(reqMsg *AccountInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUnFrezzeAccount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcUnFrezzeAccount_(reqMsg *AccountInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUnFrezzeAccount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcGetVersion(reqMsg *base.Empty) *VersionInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetVersion", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*VersionInfo)
}

func (self *AccountClient2Hall) RpcGetVersion_(reqMsg *base.Empty) (*VersionInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetVersion", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*VersionInfo), e
}
func (self *AccountClient2Hall) RpcGetBlackList(reqMsg *base.Empty) *client_server.AllPlayerInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetBlackList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.AllPlayerInfo)
}

func (self *AccountClient2Hall) RpcGetBlackList_(reqMsg *base.Empty) (*client_server.AllPlayerInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetBlackList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.AllPlayerInfo), e
}
func (self *AccountClient2Hall) RpcBlackOperate(reqMsg *BlackInfo) *BlackInfo {
	msg, e := self.Sender.CallRpcMethod("RpcBlackOperate", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BlackInfo)
}

func (self *AccountClient2Hall) RpcBlackOperate_(reqMsg *BlackInfo) (*BlackInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBlackOperate", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BlackInfo), e
}
func (self *AccountClient2Hall) RpcDelFriend(reqMsg *DelFriendInfo) *DelFriendInfo {
	msg, e := self.Sender.CallRpcMethod("RpcDelFriend", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DelFriendInfo)
}

func (self *AccountClient2Hall) RpcDelFriend_(reqMsg *DelFriendInfo) (*DelFriendInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelFriend", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DelFriendInfo), e
}
func (self *AccountClient2Hall) RpcGetCLAssistant(reqMsg *AInfo) *client_server.AssistantInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetCLAssistant", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.AssistantInfo)
}

func (self *AccountClient2Hall) RpcGetCLAssistant_(reqMsg *AInfo) (*client_server.AssistantInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCLAssistant", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.AssistantInfo), e
}
func (self *AccountClient2Hall) RpcGetMoneyAssistant(reqMsg *PageInfo) *AllMoneyAssistantInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetMoneyAssistant", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllMoneyAssistantInfo)
}

func (self *AccountClient2Hall) RpcGetMoneyAssistant_(reqMsg *PageInfo) (*AllMoneyAssistantInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetMoneyAssistant", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllMoneyAssistantInfo), e
}
func (self *AccountClient2Hall) RpcGetLocationInfo(reqMsg *LocationInfo) *AllLocationPlayerInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetLocationInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllLocationPlayerInfo)
}

func (self *AccountClient2Hall) RpcGetLocationInfo_(reqMsg *LocationInfo) (*AllLocationPlayerInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetLocationInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllLocationPlayerInfo), e
}
func (self *AccountClient2Hall) RpcGetAllNearByInfo(reqMsg *base.Empty) *AllNearByMessage {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllNearByInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllNearByMessage)
}

func (self *AccountClient2Hall) RpcGetAllNearByInfo_(reqMsg *base.Empty) (*AllNearByMessage, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllNearByInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllNearByMessage), e
}
func (self *AccountClient2Hall) RpcGetNearByInfo(reqMsg *base.Empty) *AllNearByMessage {
	msg, e := self.Sender.CallRpcMethod("RpcGetNearByInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllNearByMessage)
}

func (self *AccountClient2Hall) RpcGetNearByInfo_(reqMsg *base.Empty) (*AllNearByMessage, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetNearByInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllNearByMessage), e
}
func (self *AccountClient2Hall) RpcAddNearByMessage(reqMsg *NearByMessage) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddNearByMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcAddNearByMessage_(reqMsg *NearByMessage) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddNearByMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcDeleteNearByInfo(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteNearByInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcDeleteNearByInfo_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteNearByInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcDeleteMyNearByInfo(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteMyNearByInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcDeleteMyNearByInfo_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteMyNearByInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcBindBankCode(reqMsg *BankMessage) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcBindBankCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcBindBankCode_(reqMsg *BankMessage) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBindBankCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcAddBank(reqMsg *BankMessage) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddBank", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcAddBank_(reqMsg *BankMessage) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddBank", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcDelBank(reqMsg *BankMessage) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelBank", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcDelBank_(reqMsg *BankMessage) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelBank", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcGetBankCode(reqMsg *BankMessage) *BankMessage {
	msg, e := self.Sender.CallRpcMethod("RpcGetBankCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BankMessage)
}

func (self *AccountClient2Hall) RpcGetBankCode_(reqMsg *BankMessage) (*BankMessage, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetBankCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BankMessage), e
}
func (self *AccountClient2Hall) RpcSyncTime(reqMsg *base.Empty) *client_server.NTP {
	msg, e := self.Sender.CallRpcMethod("RpcSyncTime", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.NTP)
}

func (self *AccountClient2Hall) RpcSyncTime_(reqMsg *base.Empty) (*client_server.NTP, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSyncTime", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.NTP), e
}
func (self *AccountClient2Hall) RpcGetMoneyOrderInfo(reqMsg *MoneyType) *AllOrderInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetMoneyOrderInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllOrderInfo)
}

func (self *AccountClient2Hall) RpcGetMoneyOrderInfo_(reqMsg *MoneyType) (*AllOrderInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetMoneyOrderInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllOrderInfo), e
}
func (self *AccountClient2Hall) RpcGetCashInfo(reqMsg *PageInfo) *AllCashInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetCashInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllCashInfo)
}

func (self *AccountClient2Hall) RpcGetCashInfo_(reqMsg *PageInfo) (*AllCashInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCashInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllCashInfo), e
}
func (self *AccountClient2Hall) RpcGetRedPacketInfo(reqMsg *RedPacketInfo) *AllRedPacketInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetRedPacketInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllRedPacketInfo)
}

func (self *AccountClient2Hall) RpcGetRedPacketInfo_(reqMsg *RedPacketInfo) (*AllRedPacketInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetRedPacketInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllRedPacketInfo), e
}
func (self *AccountClient2Hall) RpcClearLocalTime(reqMsg *base.Empty) {
	self.Sender.CallRpcMethod("RpcClearLocalTime", reqMsg)
}
func (self *AccountClient2Hall) RpcRechargeMoney(reqMsg *share_message.PayOrderInfo) *share_message.PayOrderResult {
	msg, e := self.Sender.CallRpcMethod("RpcRechargeMoney", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.PayOrderResult)
}

func (self *AccountClient2Hall) RpcRechargeMoney_(reqMsg *share_message.PayOrderInfo) (*share_message.PayOrderResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRechargeMoney", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.PayOrderResult), e
}
func (self *AccountClient2Hall) RpcAddComplaintInfo(reqMsg *share_message.PlayerComplaint) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddComplaintInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcAddComplaintInfo_(reqMsg *share_message.PlayerComplaint) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddComplaintInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcWithdrawRequest(reqMsg *WithdrawInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcWithdrawRequest", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcWithdrawRequest_(reqMsg *WithdrawInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWithdrawRequest", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcCheckPeopleIdValid(reqMsg *PeopleIdInfo) *CheckPeopleInfo {
	msg, e := self.Sender.CallRpcMethod("RpcCheckPeopleIdValid", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CheckPeopleInfo)
}

func (self *AccountClient2Hall) RpcCheckPeopleIdValid_(reqMsg *PeopleIdInfo) (*CheckPeopleInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckPeopleIdValid", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CheckPeopleInfo), e
}
func (self *AccountClient2Hall) RpcFirstLoginSetInfo(reqMsg *FirstInfo) *FirstReturnInfo {
	msg, e := self.Sender.CallRpcMethod("RpcFirstLoginSetInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*FirstReturnInfo)
}

func (self *AccountClient2Hall) RpcFirstLoginSetInfo_(reqMsg *FirstInfo) (*FirstReturnInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcFirstLoginSetInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*FirstReturnInfo), e
}
func (self *AccountClient2Hall) RpcAliPay(reqMsg *share_message.AliPayData) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAliPay", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcAliPay_(reqMsg *share_message.AliPayData) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAliPay", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcPayForCode(reqMsg *PayInfo) *PayForCodeInfo {
	msg, e := self.Sender.CallRpcMethod("RpcPayForCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PayForCodeInfo)
}

func (self *AccountClient2Hall) RpcPayForCode_(reqMsg *PayInfo) (*PayForCodeInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPayForCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PayForCodeInfo), e
}
func (self *AccountClient2Hall) RpcGetCollectIndexList(reqMsg *base.Empty) *CollectIndex {
	msg, e := self.Sender.CallRpcMethod("RpcGetCollectIndexList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CollectIndex)
}

func (self *AccountClient2Hall) RpcGetCollectIndexList_(reqMsg *base.Empty) (*CollectIndex, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCollectIndexList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CollectIndex), e
}
func (self *AccountClient2Hall) RpcGetCollectIndexInfo(reqMsg *CollectIndex) *AllCollectInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetCollectIndexInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllCollectInfo)
}

func (self *AccountClient2Hall) RpcGetCollectIndexInfo_(reqMsg *CollectIndex) (*AllCollectInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCollectIndexInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllCollectInfo), e
}
func (self *AccountClient2Hall) RpcAddCollectInfo(reqMsg *share_message.CollectInfo) *share_message.CollectInfo {
	msg, e := self.Sender.CallRpcMethod("RpcAddCollectInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CollectInfo)
}

func (self *AccountClient2Hall) RpcAddCollectInfo_(reqMsg *share_message.CollectInfo) (*share_message.CollectInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddCollectInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CollectInfo), e
}
func (self *AccountClient2Hall) RpcGetCollectInfo(reqMsg *GetCollectInfo) *AllCollectInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetCollectInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllCollectInfo)
}

func (self *AccountClient2Hall) RpcGetCollectInfo_(reqMsg *GetCollectInfo) (*AllCollectInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCollectInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllCollectInfo), e
}
func (self *AccountClient2Hall) RpcDelCollectInfo(reqMsg *DelCollectInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelCollectInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcDelCollectInfo_(reqMsg *DelCollectInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelCollectInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcSearchCollectInfo(reqMsg *SearchCollectInfo) *AllCollectInfo {
	msg, e := self.Sender.CallRpcMethod("RpcSearchCollectInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllCollectInfo)
}

func (self *AccountClient2Hall) RpcSearchCollectInfo_(reqMsg *SearchCollectInfo) (*AllCollectInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSearchCollectInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllCollectInfo), e
}
func (self *AccountClient2Hall) RpcCheckout(reqMsg *share_message.OrderID) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCheckout", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcCheckout_(reqMsg *share_message.OrderID) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckout", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcGetCodeInfo(reqMsg *CodeInfo) *CodeInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetCodeInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CodeInfo)
}

func (self *AccountClient2Hall) RpcGetCodeInfo_(reqMsg *CodeInfo) (*CodeInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCodeInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CodeInfo), e
}
func (self *AccountClient2Hall) RpcCheckUnGetMoney(reqMsg *UnGetMoneyInfo) *UnGetMoneyInfo {
	msg, e := self.Sender.CallRpcMethod("RpcCheckUnGetMoney", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*UnGetMoneyInfo)
}

func (self *AccountClient2Hall) RpcCheckUnGetMoney_(reqMsg *UnGetMoneyInfo) (*UnGetMoneyInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckUnGetMoney", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*UnGetMoneyInfo), e
}
func (self *AccountClient2Hall) RpcCutBackstage(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCutBackstage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcCutBackstage_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCutBackstage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcReturnApp(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReturnApp", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcReturnApp_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReturnApp", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcReqPaySMS(reqMsg *share_message.PayOrderInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReqPaySMS", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcReqPaySMS_(reqMsg *share_message.PayOrderInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReqPaySMS", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcReqSMSPay(reqMsg *BankPaySMS) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReqSMSPay", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcReqSMSPay_(reqMsg *BankPaySMS) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReqSMSPay", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcRecommendRequest(reqMsg *RecommendMsg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcRecommendRequest", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcRecommendRequest_(reqMsg *RecommendMsg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRecommendRequest", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcFreshRecommendInfo(reqMsg *RecommendRefreshInfo) *client_server.RecommendInfo {
	msg, e := self.Sender.CallRpcMethod("RpcFreshRecommendInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.RecommendInfo)
}

func (self *AccountClient2Hall) RpcFreshRecommendInfo_(reqMsg *RecommendRefreshInfo) (*client_server.RecommendInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcFreshRecommendInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.RecommendInfo), e
}
func (self *AccountClient2Hall) RpcRequestCallInfo(reqMsg *base.Empty) *share_message.CallInfo {
	msg, e := self.Sender.CallRpcMethod("RpcRequestCallInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CallInfo)
}

func (self *AccountClient2Hall) RpcRequestCallInfo_(reqMsg *base.Empty) (*share_message.CallInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRequestCallInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CallInfo), e
}
func (self *AccountClient2Hall) RpcCheckAccountVaild(reqMsg *client_server.CheckInfo) *client_server.CheckInfo {
	msg, e := self.Sender.CallRpcMethod("RpcCheckAccountVaild", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.CheckInfo)
}

func (self *AccountClient2Hall) RpcCheckAccountVaild_(reqMsg *client_server.CheckInfo) (*client_server.CheckInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckAccountVaild", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.CheckInfo), e
}
func (self *AccountClient2Hall) RpcUserReturnApp(reqMsg *base.Empty) *client_server.TweetsListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcUserReturnApp", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.TweetsListResponse)
}

func (self *AccountClient2Hall) RpcUserReturnApp_(reqMsg *base.Empty) (*client_server.TweetsListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUserReturnApp", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.TweetsListResponse), e
}
func (self *AccountClient2Hall) RpcOperateBindWechat(reqMsg *WechatInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcOperateBindWechat", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcOperateBindWechat_(reqMsg *WechatInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOperateBindWechat", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcOperateBindPhone(reqMsg *WechatInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcOperateBindPhone", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcOperateBindPhone_(reqMsg *WechatInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOperateBindPhone", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcDelTweets(reqMsg *client_server.TweetsIdsRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelTweets", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcDelTweets_(reqMsg *client_server.TweetsIdsRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelTweets", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcReturnTweets(reqMsg *base.Empty) *client_server.TweetsListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcReturnTweets", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.TweetsListResponse)
}

func (self *AccountClient2Hall) RpcReturnTweets_(reqMsg *base.Empty) (*client_server.TweetsListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReturnTweets", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.TweetsListResponse), e
}
func (self *AccountClient2Hall) RpcModifyEmoticon(reqMsg *share_message.PlayerEmoticon) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcModifyEmoticon", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcModifyEmoticon_(reqMsg *share_message.PlayerEmoticon) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcModifyEmoticon", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcDelEmoticon(reqMsg *share_message.PlayerEmoticon) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelEmoticon", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcDelEmoticon_(reqMsg *share_message.PlayerEmoticon) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelEmoticon", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcMarkPlayer(reqMsg *MarkName) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcMarkPlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcMarkPlayer_(reqMsg *MarkName) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMarkPlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcSumClicksJumps(reqMsg *ArticleOptRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSumClicksJumps", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcSumClicksJumps_(reqMsg *ArticleOptRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSumClicksJumps", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcNoticeClicks(reqMsg *NoticeOptRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcNoticeClicks", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcNoticeClicks_(reqMsg *NoticeOptRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNoticeClicks", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcOpenMyMainPage(reqMsg *OpenMyMainPageByPageReq) *client_server.MyMainPageInfo {
	msg, e := self.Sender.CallRpcMethod("RpcOpenMyMainPage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.MyMainPageInfo)
}

func (self *AccountClient2Hall) RpcOpenMyMainPage_(reqMsg *OpenMyMainPageByPageReq) (*client_server.MyMainPageInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOpenMyMainPage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.MyMainPageInfo), e
}
func (self *AccountClient2Hall) RpcOpenMyMainPageByPage(reqMsg *OpenMyMainPageByPageReq) *client_server.MyMainPageInfo {
	msg, e := self.Sender.CallRpcMethod("RpcOpenMyMainPageByPage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.MyMainPageInfo)
}

func (self *AccountClient2Hall) RpcOpenMyMainPageByPage_(reqMsg *OpenMyMainPageByPageReq) (*client_server.MyMainPageInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOpenMyMainPageByPage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.MyMainPageInfo), e
}
func (self *AccountClient2Hall) RpcGetMyFansAttentionInfo(reqMsg *MainInfo) *AllFansInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetMyFansAttentionInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllFansInfo)
}

func (self *AccountClient2Hall) RpcGetMyFansAttentionInfo_(reqMsg *MainInfo) (*AllFansInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetMyFansAttentionInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllFansInfo), e
}
func (self *AccountClient2Hall) RpcGetSomeDynamic(reqMsg *DynamicIdInfo) *DynamicInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetSomeDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DynamicInfo)
}

func (self *AccountClient2Hall) RpcGetSomeDynamic_(reqMsg *DynamicIdInfo) (*DynamicInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSomeDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DynamicInfo), e
}
func (self *AccountClient2Hall) RpcGetPlayerOtherData(reqMsg *base.Empty) *client_server.AllPlayerMsg {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerOtherData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.AllPlayerMsg)
}

func (self *AccountClient2Hall) RpcGetPlayerOtherData_(reqMsg *base.Empty) (*client_server.AllPlayerMsg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerOtherData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.AllPlayerMsg), e
}
func (self *AccountClient2Hall) RpcGetTeamMembers(reqMsg *TeamMembers) *TeamMembers {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamMembers", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TeamMembers)
}

func (self *AccountClient2Hall) RpcGetTeamMembers_(reqMsg *TeamMembers) (*TeamMembers, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamMembers", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TeamMembers), e
}
func (self *AccountClient2Hall) RpcCheckDirtyWord(reqMsg *CheckDirtyWord) *CheckDirtyWord {
	msg, e := self.Sender.CallRpcMethod("RpcCheckDirtyWord", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CheckDirtyWord)
}

func (self *AccountClient2Hall) RpcCheckDirtyWord_(reqMsg *CheckDirtyWord) (*CheckDirtyWord, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckDirtyWord", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CheckDirtyWord), e
}
func (self *AccountClient2Hall) RpcCheckCancelAccount(reqMsg *base.Empty) *CheckCancelAccount {
	msg, e := self.Sender.CallRpcMethod("RpcCheckCancelAccount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CheckCancelAccount)
}

func (self *AccountClient2Hall) RpcCheckCancelAccount_(reqMsg *base.Empty) (*CheckCancelAccount, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckCancelAccount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CheckCancelAccount), e
}
func (self *AccountClient2Hall) RpcSubmitCancelAccount(reqMsg *CancelAccountData) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSubmitCancelAccount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcSubmitCancelAccount_(reqMsg *CancelAccountData) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSubmitCancelAccount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcAttentioPlayer(reqMsg *client_server.AttenInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAttentioPlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcAttentioPlayer_(reqMsg *client_server.AttenInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAttentioPlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcESportAttentioPlayer(reqMsg *client_server.AttenInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcESportAttentioPlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcESportAttentioPlayer_(reqMsg *client_server.AttenInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcESportAttentioPlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcZanOperateSquareDynamic(reqMsg *client_server.ZanInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcZanOperateSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcZanOperateSquareDynamic_(reqMsg *client_server.ZanInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcZanOperateSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcAddCommentSquareDynamic(reqMsg *share_message.CommentData) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddCommentSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcAddCommentSquareDynamic_(reqMsg *share_message.CommentData) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddCommentSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcDelCommentSquareDynamic(reqMsg *client_server.IdInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelCommentSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcDelCommentSquareDynamic_(reqMsg *client_server.IdInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelCommentSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcAddSquareCommentZan(reqMsg *share_message.CommentDataZan) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddSquareCommentZan", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcAddSquareCommentZan_(reqMsg *share_message.CommentDataZan) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddSquareCommentZan", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcDelSquareDynamic(reqMsg *client_server.RequestInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelSquareDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcDelSquareDynamic_(reqMsg *client_server.RequestInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelSquareDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcShopOrderComplaintDetail(reqMsg *ComplaintID) *ShopOrderComplaintDetailRsp {
	msg, e := self.Sender.CallRpcMethod("RpcShopOrderComplaintDetail", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ShopOrderComplaintDetailRsp)
}

func (self *AccountClient2Hall) RpcShopOrderComplaintDetail_(reqMsg *ComplaintID) (*ShopOrderComplaintDetailRsp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopOrderComplaintDetail", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ShopOrderComplaintDetailRsp), e
}
func (self *AccountClient2Hall) RpcGetAdvData(reqMsg *AdvSettingRequest) *AdvSettingResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetAdvData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AdvSettingResponse)
}

func (self *AccountClient2Hall) RpcGetAdvData_(reqMsg *AdvSettingRequest) (*AdvSettingResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAdvData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AdvSettingResponse), e
}
func (self *AccountClient2Hall) RpcGetNearByInfo2(reqMsg *LocationInfo) *NearByInfoReply {
	msg, e := self.Sender.CallRpcMethod("RpcGetNearByInfo2", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NearByInfoReply)
}

func (self *AccountClient2Hall) RpcGetNearByInfo2_(reqMsg *LocationInfo) (*NearByInfoReply, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetNearByInfo2", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NearByInfoReply), e
}
func (self *AccountClient2Hall) RpcGetNearInfoByPage(reqMsg *LocationInfoByPage) *NearByInfoReplyByPage {
	msg, e := self.Sender.CallRpcMethod("RpcGetNearInfoByPage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NearByInfoReplyByPage)
}

func (self *AccountClient2Hall) RpcGetNearInfoByPage_(reqMsg *LocationInfoByPage) (*NearByInfoReplyByPage, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetNearInfoByPage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NearByInfoReplyByPage), e
}
func (self *AccountClient2Hall) RpcAddAdvLog(reqMsg *share_message.AdvLogReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddAdvLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcAddAdvLog_(reqMsg *share_message.AdvLogReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddAdvLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcDelNewFriendList(reqMsg *DelNewFriendListReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelNewFriendList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcDelNewFriendList_(reqMsg *DelNewFriendListReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelNewFriendList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcNewTeamSetting(reqMsg *NewTeamSettingReq) *TeamSettingNotify {
	msg, e := self.Sender.CallRpcMethod("RpcNewTeamSetting", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TeamSettingNotify)
}

func (self *AccountClient2Hall) RpcNewTeamSetting_(reqMsg *NewTeamSettingReq) (*TeamSettingNotify, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNewTeamSetting", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TeamSettingNotify), e
}
func (self *AccountClient2Hall) RpcLocationInfoNew(reqMsg *LocationInfoNewReq) *LocationInfoNewResp {
	msg, e := self.Sender.CallRpcMethod("RpcLocationInfoNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LocationInfoNewResp)
}

func (self *AccountClient2Hall) RpcLocationInfoNew_(reqMsg *LocationInfoNewReq) (*LocationInfoNewResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLocationInfoNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LocationInfoNewResp), e
}
func (self *AccountClient2Hall) RpcNearRecommend(reqMsg *NearRecommendReq) *NearRecommendResp {
	msg, e := self.Sender.CallRpcMethod("RpcNearRecommend", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NearRecommendResp)
}

func (self *AccountClient2Hall) RpcNearRecommend_(reqMsg *NearRecommendReq) (*NearRecommendResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNearRecommend", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NearRecommendResp), e
}
func (self *AccountClient2Hall) RpcNearSayMessage(reqMsg *share_message.NearSessionList) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcNearSayMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcNearSayMessage_(reqMsg *share_message.NearSessionList) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNearSayMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcNearSessionList(reqMsg *NearSessionListReq) *NearSessionListResp {
	msg, e := self.Sender.CallRpcMethod("RpcNearSessionList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NearSessionListResp)
}

func (self *AccountClient2Hall) RpcNearSessionList_(reqMsg *NearSessionListReq) (*NearSessionListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNearSessionList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NearSessionListResp), e
}
func (self *AccountClient2Hall) RpcSendPlayerMessageList(reqMsg *SendPlayerMessageListReq) *NearSessionListResp {
	msg, e := self.Sender.CallRpcMethod("RpcSendPlayerMessageList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NearSessionListResp)
}

func (self *AccountClient2Hall) RpcSendPlayerMessageList_(reqMsg *SendPlayerMessageListReq) (*NearSessionListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSendPlayerMessageList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NearSessionListResp), e
}
func (self *AccountClient2Hall) RpcGetNearChatList(reqMsg *GetNearChatListReq) *NearSessionListResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetNearChatList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NearSessionListResp)
}

func (self *AccountClient2Hall) RpcGetNearChatList_(reqMsg *GetNearChatListReq) (*NearSessionListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetNearChatList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NearSessionListResp), e
}
func (self *AccountClient2Hall) RpcDelNearMessage(reqMsg *DelNearMessageReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelNearMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcDelNearMessage_(reqMsg *DelNearMessageReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelNearMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcUpdateIsReadReq(reqMsg *UpdateIsReadReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateIsReadReq", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcUpdateIsReadReq_(reqMsg *UpdateIsReadReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateIsReadReq", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcSetYoungPassWord(reqMsg *YoungPassWord) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSetYoungPassWord", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcSetYoungPassWord_(reqMsg *YoungPassWord) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSetYoungPassWord", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcNearAddAdvLog(reqMsg *share_message.AdvLogReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcNearAddAdvLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcNearAddAdvLog_(reqMsg *share_message.AdvLogReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNearAddAdvLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcGetPlayerOtherDataNew(reqMsg *base.Empty) *client_server.AllPlayerMsg {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerOtherDataNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.AllPlayerMsg)
}

func (self *AccountClient2Hall) RpcGetPlayerOtherDataNew_(reqMsg *base.Empty) (*client_server.AllPlayerMsg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerOtherDataNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.AllPlayerMsg), e
}
func (self *AccountClient2Hall) RpcGetPlayerFriends(reqMsg *base.Empty) *client_server.AllPlayerMsg {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerFriends", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.AllPlayerMsg)
}

func (self *AccountClient2Hall) RpcGetPlayerFriends_(reqMsg *base.Empty) (*client_server.AllPlayerMsg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerFriends", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.AllPlayerMsg), e
}
func (self *AccountClient2Hall) RpcGetPlayerNewFriends(reqMsg *GetNewFriends) *client_server.NewFriends {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerNewFriends", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.NewFriends)
}

func (self *AccountClient2Hall) RpcGetPlayerNewFriends_(reqMsg *GetNewFriends) (*client_server.NewFriends, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerNewFriends", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.NewFriends), e
}
func (self *AccountClient2Hall) RpcGetTeamMemberChange(reqMsg *TeamMembersChange) *TeamMembersChange {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamMemberChange", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TeamMembersChange)
}

func (self *AccountClient2Hall) RpcGetTeamMemberChange_(reqMsg *TeamMembersChange) (*TeamMembersChange, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamMemberChange", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TeamMembersChange), e
}
func (self *AccountClient2Hall) RpcGetSupportBankList(reqMsg *SupportBankList) *SupportBankList {
	msg, e := self.Sender.CallRpcMethod("RpcGetSupportBankList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SupportBankList)
}

func (self *AccountClient2Hall) RpcGetSupportBankList_(reqMsg *SupportBankList) (*SupportBankList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSupportBankList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SupportBankList), e
}
func (self *AccountClient2Hall) RpcGetOnLineNum(reqMsg *base.Empty) *OnLineNum {
	msg, e := self.Sender.CallRpcMethod("RpcGetOnLineNum", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*OnLineNum)
}

func (self *AccountClient2Hall) RpcGetOnLineNum_(reqMsg *base.Empty) (*OnLineNum, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetOnLineNum", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*OnLineNum), e
}
func (self *AccountClient2Hall) RpcGetMsgPageAdvList(reqMsg *base.Empty) *MsgAdv {
	msg, e := self.Sender.CallRpcMethod("RpcGetMsgPageAdvList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MsgAdv)
}

func (self *AccountClient2Hall) RpcGetMsgPageAdvList_(reqMsg *base.Empty) (*MsgAdv, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetMsgPageAdvList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MsgAdv), e
}
func (self *AccountClient2Hall) RpcGetDiamond(reqMsg *base.Empty) *PlayerDiamond {
	msg, e := self.Sender.CallRpcMethod("RpcGetDiamond", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerDiamond)
}

func (self *AccountClient2Hall) RpcGetDiamond_(reqMsg *base.Empty) (*PlayerDiamond, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDiamond", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerDiamond), e
}
func (self *AccountClient2Hall) RpcGetSupportPayChannel(reqMsg *PayChannels) *PayChannels {
	msg, e := self.Sender.CallRpcMethod("RpcGetSupportPayChannel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PayChannels)
}

func (self *AccountClient2Hall) RpcGetSupportPayChannel_(reqMsg *PayChannels) (*PayChannels, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSupportPayChannel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PayChannels), e
}
func (self *AccountClient2Hall) RpcGetAllLabelMsg(reqMsg *base.Empty) *client_server.LabelMsg {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllLabelMsg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.LabelMsg)
}

func (self *AccountClient2Hall) RpcGetAllLabelMsg_(reqMsg *base.Empty) (*client_server.LabelMsg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllLabelMsg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.LabelMsg), e
}
func (self *AccountClient2Hall) RpcSetPlayerLabel(reqMsg *FirstInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSetPlayerLabel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountClient2Hall) RpcSetPlayerLabel_(reqMsg *FirstInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSetPlayerLabel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountClient2Hall) RpcTopicTeamDynamicList(reqMsg *TopicTeamDynamicReq) *TopicTeamDynamicResp {
	msg, e := self.Sender.CallRpcMethod("RpcTopicTeamDynamicList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicTeamDynamicResp)
}

func (self *AccountClient2Hall) RpcTopicTeamDynamicList_(reqMsg *TopicTeamDynamicReq) (*TopicTeamDynamicResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTopicTeamDynamicList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicTeamDynamicResp), e
}
func (self *AccountClient2Hall) RpcGetTopicTeams(reqMsg *TopicTeams) *TopicTeams {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicTeams", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicTeams)
}

func (self *AccountClient2Hall) RpcGetTopicTeams_(reqMsg *TopicTeams) (*TopicTeams, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicTeams", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicTeams), e
}
func (self *AccountClient2Hall) RpcDefunctTeam(reqMsg *DefunctTeam) *DefunctTeam {
	msg, e := self.Sender.CallRpcMethod("RpcDefunctTeam", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DefunctTeam)
}

func (self *AccountClient2Hall) RpcDefunctTeam_(reqMsg *DefunctTeam) (*DefunctTeam, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDefunctTeam", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DefunctTeam), e
}
func (self *AccountClient2Hall) RpcGetAllMainMenu(reqMsg *base.Empty) *AllMainMenu {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllMainMenu", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllMainMenu)
}

func (self *AccountClient2Hall) RpcGetAllMainMenu_(reqMsg *base.Empty) (*AllMainMenu, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllMainMenu", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllMainMenu), e
}
func (self *AccountClient2Hall) RpcGetAllTipAdvs(reqMsg *base.Empty) *AllTipAdv {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllTipAdvs", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllTipAdv)
}

func (self *AccountClient2Hall) RpcGetAllTipAdvs_(reqMsg *base.Empty) (*AllTipAdv, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllTipAdvs", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllTipAdv), e
}
func (self *AccountClient2Hall) RpcGetStartPageAdvList(reqMsg *base.Empty) *MsgAdv {
	msg, e := self.Sender.CallRpcMethod("RpcGetStartPageAdvList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MsgAdv)
}

func (self *AccountClient2Hall) RpcGetStartPageAdvList_(reqMsg *base.Empty) (*MsgAdv, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetStartPageAdvList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MsgAdv), e
}

// ==========================================================
type IAccountHall2Client interface {
	RpcPlayerLoginResponse(reqMsg *client_server.AllPlayerMsg)
	RpcUpdateGold(reqMsg *UpdateGold)
	RpcReLogin(reqMsg *base.Empty)
	RpcKickOut(reqMsg *base.Empty)
	RpcNoticeAddFriend(reqMsg *share_message.AllAddPlayerMsg)
	RpcNoticeAddFriendNew(reqMsg *share_message.AddPlayerRequest)
	RpcNoticeAgreeFriend(reqMsg *client_server.PlayerMsg)
	RpcCreateTeamResult(reqMsg *client_server.TeamMsg)
	RpcAddTeamResult(reqMsg *client_server.TeamMsg)
	RpcFindFriendResponse(reqMsg *client_server.PlayerMsg)
	RpcAgreeFriendResponse(reqMsg *client_server.PlayerMsg)
	RpcNoticeNearByInfo(reqMsg *base.Empty)
	RpcRechargeMoneyResult(reqMsg *share_message.RechargeOrderResult)
	RpcRechargeMoneyFinish(reqMsg *share_message.RechargeFinish)
	RpcAssistantNotify(reqMsg *client_server.AssistantMsg)
	RpcAliPayResult(reqMsg *share_message.AliPayData) *base.Empty
	RpcAliPayResult_(reqMsg *share_message.AliPayData) (*base.Empty, easygo.IRpcInterrupt)
	RpcShopItemMessageNotify(reqMsg *share_message.ShopItemMessageInfo)
	RpcDaiFuResult(reqMsg *WithdrawInfo)
	RpcFreezePlayerLogOut(reqMsg *base.Empty)
	RpcBindBankCodeResult(reqMsg *BankMessage)
	RpcReqPaySMSResult(reqMsg *BankPaySMS)
	RpcReturnRecommendInfo(reqMsg *client_server.RecommendInfo)
	RpcAssistantNotifyArticle(reqMsg *client_server.ArticleListResponse)
	RpcRegisterPush(reqMsg *client_server.TweetsListResponse)
	RpcNoticeAddGold(reqMsg *AddGoldInfo)
	RpcSysParameterChange(reqMsg *share_message.SysParameter)
	RpcPayChannelChange(reqMsg *client_server.AllPlayerMsg)
	RpcNewMessage(reqMsg *NewUnReadMessageResp)
	RpcNoNewMessage(reqMsg *base.Empty)
	RpcShopOrderNotify(reqMsg *share_message.ShopOrderNotifyInfo)
	RpcTeamSettingNotify(reqMsg *TeamSettingNotify)
	RpcNotifyNewSayMessage(reqMsg *NotifyNewSayMessageReq)
	RpcHasUnReadNear(reqMsg *HasUnReadNearResp)
}

type AccountHall2Client struct {
	Sender easygo.IMessageSender
}

func (self *AccountHall2Client) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *AccountHall2Client) RpcPlayerLoginResponse(reqMsg *client_server.AllPlayerMsg) {
	self.Sender.CallRpcMethod("RpcPlayerLoginResponse", reqMsg)
}
func (self *AccountHall2Client) RpcUpdateGold(reqMsg *UpdateGold) {
	self.Sender.CallRpcMethod("RpcUpdateGold", reqMsg)
}
func (self *AccountHall2Client) RpcReLogin(reqMsg *base.Empty) {
	self.Sender.CallRpcMethod("RpcReLogin", reqMsg)
}
func (self *AccountHall2Client) RpcKickOut(reqMsg *base.Empty) {
	self.Sender.CallRpcMethod("RpcKickOut", reqMsg)
}
func (self *AccountHall2Client) RpcNoticeAddFriend(reqMsg *share_message.AllAddPlayerMsg) {
	self.Sender.CallRpcMethod("RpcNoticeAddFriend", reqMsg)
}
func (self *AccountHall2Client) RpcNoticeAddFriendNew(reqMsg *share_message.AddPlayerRequest) {
	self.Sender.CallRpcMethod("RpcNoticeAddFriendNew", reqMsg)
}
func (self *AccountHall2Client) RpcNoticeAgreeFriend(reqMsg *client_server.PlayerMsg) {
	self.Sender.CallRpcMethod("RpcNoticeAgreeFriend", reqMsg)
}
func (self *AccountHall2Client) RpcCreateTeamResult(reqMsg *client_server.TeamMsg) {
	self.Sender.CallRpcMethod("RpcCreateTeamResult", reqMsg)
}
func (self *AccountHall2Client) RpcAddTeamResult(reqMsg *client_server.TeamMsg) {
	self.Sender.CallRpcMethod("RpcAddTeamResult", reqMsg)
}
func (self *AccountHall2Client) RpcFindFriendResponse(reqMsg *client_server.PlayerMsg) {
	self.Sender.CallRpcMethod("RpcFindFriendResponse", reqMsg)
}
func (self *AccountHall2Client) RpcAgreeFriendResponse(reqMsg *client_server.PlayerMsg) {
	self.Sender.CallRpcMethod("RpcAgreeFriendResponse", reqMsg)
}
func (self *AccountHall2Client) RpcNoticeNearByInfo(reqMsg *base.Empty) {
	self.Sender.CallRpcMethod("RpcNoticeNearByInfo", reqMsg)
}
func (self *AccountHall2Client) RpcRechargeMoneyResult(reqMsg *share_message.RechargeOrderResult) {
	self.Sender.CallRpcMethod("RpcRechargeMoneyResult", reqMsg)
}
func (self *AccountHall2Client) RpcRechargeMoneyFinish(reqMsg *share_message.RechargeFinish) {
	self.Sender.CallRpcMethod("RpcRechargeMoneyFinish", reqMsg)
}
func (self *AccountHall2Client) RpcAssistantNotify(reqMsg *client_server.AssistantMsg) {
	self.Sender.CallRpcMethod("RpcAssistantNotify", reqMsg)
}
func (self *AccountHall2Client) RpcAliPayResult(reqMsg *share_message.AliPayData) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAliPayResult", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *AccountHall2Client) RpcAliPayResult_(reqMsg *share_message.AliPayData) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAliPayResult", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *AccountHall2Client) RpcShopItemMessageNotify(reqMsg *share_message.ShopItemMessageInfo) {
	self.Sender.CallRpcMethod("RpcShopItemMessageNotify", reqMsg)
}
func (self *AccountHall2Client) RpcDaiFuResult(reqMsg *WithdrawInfo) {
	self.Sender.CallRpcMethod("RpcDaiFuResult", reqMsg)
}
func (self *AccountHall2Client) RpcFreezePlayerLogOut(reqMsg *base.Empty) {
	self.Sender.CallRpcMethod("RpcFreezePlayerLogOut", reqMsg)
}
func (self *AccountHall2Client) RpcBindBankCodeResult(reqMsg *BankMessage) {
	self.Sender.CallRpcMethod("RpcBindBankCodeResult", reqMsg)
}
func (self *AccountHall2Client) RpcReqPaySMSResult(reqMsg *BankPaySMS) {
	self.Sender.CallRpcMethod("RpcReqPaySMSResult", reqMsg)
}
func (self *AccountHall2Client) RpcReturnRecommendInfo(reqMsg *client_server.RecommendInfo) {
	self.Sender.CallRpcMethod("RpcReturnRecommendInfo", reqMsg)
}
func (self *AccountHall2Client) RpcAssistantNotifyArticle(reqMsg *client_server.ArticleListResponse) {
	self.Sender.CallRpcMethod("RpcAssistantNotifyArticle", reqMsg)
}
func (self *AccountHall2Client) RpcRegisterPush(reqMsg *client_server.TweetsListResponse) {
	self.Sender.CallRpcMethod("RpcRegisterPush", reqMsg)
}
func (self *AccountHall2Client) RpcNoticeAddGold(reqMsg *AddGoldInfo) {
	self.Sender.CallRpcMethod("RpcNoticeAddGold", reqMsg)
}
func (self *AccountHall2Client) RpcSysParameterChange(reqMsg *share_message.SysParameter) {
	self.Sender.CallRpcMethod("RpcSysParameterChange", reqMsg)
}
func (self *AccountHall2Client) RpcPayChannelChange(reqMsg *client_server.AllPlayerMsg) {
	self.Sender.CallRpcMethod("RpcPayChannelChange", reqMsg)
}
func (self *AccountHall2Client) RpcNewMessage(reqMsg *NewUnReadMessageResp) {
	self.Sender.CallRpcMethod("RpcNewMessage", reqMsg)
}
func (self *AccountHall2Client) RpcNoNewMessage(reqMsg *base.Empty) {
	self.Sender.CallRpcMethod("RpcNoNewMessage", reqMsg)
}
func (self *AccountHall2Client) RpcShopOrderNotify(reqMsg *share_message.ShopOrderNotifyInfo) {
	self.Sender.CallRpcMethod("RpcShopOrderNotify", reqMsg)
}
func (self *AccountHall2Client) RpcTeamSettingNotify(reqMsg *TeamSettingNotify) {
	self.Sender.CallRpcMethod("RpcTeamSettingNotify", reqMsg)
}
func (self *AccountHall2Client) RpcNotifyNewSayMessage(reqMsg *NotifyNewSayMessageReq) {
	self.Sender.CallRpcMethod("RpcNotifyNewSayMessage", reqMsg)
}
func (self *AccountHall2Client) RpcHasUnReadNear(reqMsg *HasUnReadNearResp) {
	self.Sender.CallRpcMethod("RpcHasUnReadNear", reqMsg)
}
