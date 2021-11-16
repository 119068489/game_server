package client_hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type IChatClient2Hall interface {
	RpcChat(reqMsg *share_message.Chat) *base.Empty
	RpcChat_(reqMsg *share_message.Chat) (*base.Empty, easygo.IRpcInterrupt)
	RpcReadMessage(reqMsg *client_server.ReadInfo) *base.Empty
	RpcReadMessage_(reqMsg *client_server.ReadInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcSendRedPacket(reqMsg *share_message.RedPacket) *base.Empty
	RpcSendRedPacket_(reqMsg *share_message.RedPacket) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddTeamMember(reqMsg *client_server.TeamReq) *base.Empty
	RpcAddTeamMember_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcRemoveTeamMember(reqMsg *client_server.TeamReq) *base.Empty
	RpcRemoveTeamMember_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetTeamSetting(reqMsg *client_server.TeamInfo) *base.Empty
	RpcGetTeamSetting_(reqMsg *client_server.TeamInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetTeamManageSetting(reqMsg *client_server.TeamInfo) *base.Empty
	RpcGetTeamManageSetting_(reqMsg *client_server.TeamInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcChangeTeamPersonalSetting(reqMsg *CgTeamPerSetting) *CgTeamPerSetting
	RpcChangeTeamPersonalSetting_(reqMsg *CgTeamPerSetting) (*CgTeamPerSetting, easygo.IRpcInterrupt)
	RpcChangeTeamManageSetting(reqMsg *CgTeamManageSetting) *CgTeamManageSetting
	RpcChangeTeamManageSetting_(reqMsg *CgTeamManageSetting) (*CgTeamManageSetting, easygo.IRpcInterrupt)
	RpcAddTeamManager(reqMsg *client_server.TeamReq) *base.Empty
	RpcAddTeamManager_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelTeamManager(reqMsg *client_server.TeamReq) *base.Empty
	RpcDelTeamManager_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcChangeTeamOwner(reqMsg *client_server.TeamReq) *base.Empty
	RpcChangeTeamOwner_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcOpenRedPacket(reqMsg *OpenRedPacket) *base.Empty
	RpcOpenRedPacket_(reqMsg *OpenRedPacket) (*base.Empty, easygo.IRpcInterrupt)
	RpcCheckRedPacket(reqMsg *OpenRedPacket) *base.Empty
	RpcCheckRedPacket_(reqMsg *OpenRedPacket) (*base.Empty, easygo.IRpcInterrupt)
	RpcTransferMoney(reqMsg *share_message.TransferMoney) *base.Empty
	RpcTransferMoney_(reqMsg *share_message.TransferMoney) (*base.Empty, easygo.IRpcInterrupt)
	RpcOpenTransferMoney(reqMsg *share_message.TransferMoney) *base.Empty
	RpcOpenTransferMoney_(reqMsg *share_message.TransferMoney) (*base.Empty, easygo.IRpcInterrupt)
	RpcAcceptAddTeam(reqMsg *client_server.TeamInfo) *base.Empty
	RpcAcceptAddTeam_(reqMsg *client_server.TeamInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcRefuseAddTeam(reqMsg *client_server.TeamInfo) *base.Empty
	RpcRefuseAddTeam_(reqMsg *client_server.TeamInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcChangeFriendSetting(reqMsg *FriendSetting) *FriendSetting
	RpcChangeFriendSetting_(reqMsg *FriendSetting) (*FriendSetting, easygo.IRpcInterrupt)
	RpcExitTeam(reqMsg *client_server.TeamInfo) *base.Empty
	RpcExitTeam_(reqMsg *client_server.TeamInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetTeamPlayerInfo(reqMsg *client_server.TeamInfo) *share_message.TeamPlayerInfo
	RpcGetTeamPlayerInfo_(reqMsg *client_server.TeamInfo) (*share_message.TeamPlayerInfo, easygo.IRpcInterrupt)
	RpcGetBaseTeamInfo(reqMsg *client_server.TeamReq) *TeamDataInfo
	RpcGetBaseTeamInfo_(reqMsg *client_server.TeamReq) (*TeamDataInfo, easygo.IRpcInterrupt)
	RpcGetPlayerCardInfo(reqMsg *client_server.PlayerReq) *share_message.TeamPlayerInfo
	RpcGetPlayerCardInfo_(reqMsg *client_server.PlayerReq) (*share_message.TeamPlayerInfo, easygo.IRpcInterrupt)
	RpcWithdrawMessage(reqMsg *client_server.ReadInfo) *LogInfo
	RpcWithdrawMessage_(reqMsg *client_server.ReadInfo) (*LogInfo, easygo.IRpcInterrupt)
	RpcRequestSpecialChat(reqMsg *SpecialChatInfo) *base.Empty
	RpcRequestSpecialChat_(reqMsg *SpecialChatInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcOperateSpecialChat(reqMsg *SpecialChatInfo) *base.Empty
	RpcOperateSpecialChat_(reqMsg *SpecialChatInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcRequestChatInfo(reqMsg *ChatInfo) *base.Empty
	RpcRequestChatInfo_(reqMsg *ChatInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWaiterMsg(reqMsg *WaiterMsgRequest) *share_message.IMmessage
	RpcGetWaiterMsg_(reqMsg *WaiterMsgRequest) (*share_message.IMmessage, easygo.IRpcInterrupt)
	RpcSendWaiterMsg(reqMsg *share_message.IMmessage) *share_message.IMmessage
	RpcSendWaiterMsg_(reqMsg *share_message.IMmessage) (*share_message.IMmessage, easygo.IRpcInterrupt)
	RpcRequestWaiterTypes(reqMsg *base.Empty) *WaiterTypesResponse
	RpcRequestWaiterTypes_(reqMsg *base.Empty) (*WaiterTypesResponse, easygo.IRpcInterrupt)
	RpcRequestWaiterService(reqMsg *WaiterMsgRequest) *WaiterMsgResponse
	RpcRequestWaiterService_(reqMsg *WaiterMsgRequest) (*WaiterMsgResponse, easygo.IRpcInterrupt)
	RpcWaiterGrade(reqMsg *share_message.IMmessage) *base.Empty
	RpcWaiterGrade_(reqMsg *share_message.IMmessage) (*base.Empty, easygo.IRpcInterrupt)
	RpcSearchForKey(reqMsg *SearchFaqRequest) *SearchFaqResponse
	RpcSearchForKey_(reqMsg *SearchFaqRequest) (*SearchFaqResponse, easygo.IRpcInterrupt)
	RpcOpenFaqById(reqMsg *OpenFaqRequest) *share_message.WaiterFAQ
	RpcOpenFaqById_(reqMsg *OpenFaqRequest) (*share_message.WaiterFAQ, easygo.IRpcInterrupt)
	RpcBroadCastQTX(reqMsg *BroadCastQTX) *base.Empty
	RpcBroadCastQTX_(reqMsg *BroadCastQTX) (*base.Empty, easygo.IRpcInterrupt)
	RpcChatNew(reqMsg *share_message.Chat) *share_message.Chat
	RpcChatNew_(reqMsg *share_message.Chat) (*share_message.Chat, easygo.IRpcInterrupt)
	RpcGetSessionData(reqMsg *AllSessionData) *AllSessionData
	RpcGetSessionData_(reqMsg *AllSessionData) (*AllSessionData, easygo.IRpcInterrupt)
	RpcGetSessionDetail(reqMsg *SessionData) *SessionData
	RpcGetSessionDetail_(reqMsg *SessionData) (*SessionData, easygo.IRpcInterrupt)
	RpcGetSessionChat(reqMsg *SessionChatData) *SessionChatData
	RpcGetSessionChat_(reqMsg *SessionChatData) (*SessionChatData, easygo.IRpcInterrupt)
	RpcGetTeamMemberData(reqMsg *TeamMemberData) *TeamMemberData
	RpcGetTeamMemberData_(reqMsg *TeamMemberData) (*TeamMemberData, easygo.IRpcInterrupt)
	RpcGetTeamDetailData(reqMsg *TeamDetailData) *TeamDetailData
	RpcGetTeamDetailData_(reqMsg *TeamDetailData) (*TeamDetailData, easygo.IRpcInterrupt)
	RpcGetTeamAtData(reqMsg *TeamAtData) *TeamAtData
	RpcGetTeamAtData_(reqMsg *TeamAtData) (*TeamAtData, easygo.IRpcInterrupt)
	RpcDeleteMessage(reqMsg *client_server.ReadInfo) *base.Empty
	RpcDeleteMessage_(reqMsg *client_server.ReadInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetTeamSettingNew(reqMsg *client_server.TeamInfo) *share_message.TeamSetting
	RpcGetTeamSettingNew_(reqMsg *client_server.TeamInfo) (*share_message.TeamSetting, easygo.IRpcInterrupt)
	RpcGetTeamManageSettingNew(reqMsg *client_server.TeamInfo) *client_server.TeamManagerSetting
	RpcGetTeamManageSettingNew_(reqMsg *client_server.TeamInfo) (*client_server.TeamManagerSetting, easygo.IRpcInterrupt)
	RpcGetOneSessionData(reqMsg *SessionData) *SessionData
	RpcGetOneSessionData_(reqMsg *SessionData) (*SessionData, easygo.IRpcInterrupt)
	RpcCheckIsTeamMember(reqMsg *CheckTeamMember) *CheckTeamMember
	RpcCheckIsTeamMember_(reqMsg *CheckTeamMember) (*CheckTeamMember, easygo.IRpcInterrupt)
	RpcGetSaveTeamSessions(reqMsg *base.Empty) *AllSessionData
	RpcGetSaveTeamSessions_(reqMsg *base.Empty) (*AllSessionData, easygo.IRpcInterrupt)
	RpcCheckIsMyTeamSession(reqMsg *CheckIsMySession) *CheckIsMySession
	RpcCheckIsMyTeamSession_(reqMsg *CheckIsMySession) (*CheckIsMySession, easygo.IRpcInterrupt)
	RpcGetOneShowLog(reqMsg *GetOneShowLog) *GetOneShowLog
	RpcGetOneShowLog_(reqMsg *GetOneShowLog) (*GetOneShowLog, easygo.IRpcInterrupt)
}

type ChatClient2Hall struct {
	Sender easygo.IMessageSender
}

func (self *ChatClient2Hall) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *ChatClient2Hall) RpcChat(reqMsg *share_message.Chat) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcChat", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcChat_(reqMsg *share_message.Chat) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcChat", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcReadMessage(reqMsg *client_server.ReadInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReadMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcReadMessage_(reqMsg *client_server.ReadInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReadMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcSendRedPacket(reqMsg *share_message.RedPacket) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSendRedPacket", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcSendRedPacket_(reqMsg *share_message.RedPacket) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSendRedPacket", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcAddTeamMember(reqMsg *client_server.TeamReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddTeamMember", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcAddTeamMember_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddTeamMember", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcRemoveTeamMember(reqMsg *client_server.TeamReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcRemoveTeamMember", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcRemoveTeamMember_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRemoveTeamMember", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcGetTeamSetting(reqMsg *client_server.TeamInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamSetting", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcGetTeamSetting_(reqMsg *client_server.TeamInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamSetting", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcGetTeamManageSetting(reqMsg *client_server.TeamInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamManageSetting", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcGetTeamManageSetting_(reqMsg *client_server.TeamInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamManageSetting", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcChangeTeamPersonalSetting(reqMsg *CgTeamPerSetting) *CgTeamPerSetting {
	msg, e := self.Sender.CallRpcMethod("RpcChangeTeamPersonalSetting", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CgTeamPerSetting)
}

func (self *ChatClient2Hall) RpcChangeTeamPersonalSetting_(reqMsg *CgTeamPerSetting) (*CgTeamPerSetting, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcChangeTeamPersonalSetting", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CgTeamPerSetting), e
}
func (self *ChatClient2Hall) RpcChangeTeamManageSetting(reqMsg *CgTeamManageSetting) *CgTeamManageSetting {
	msg, e := self.Sender.CallRpcMethod("RpcChangeTeamManageSetting", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CgTeamManageSetting)
}

func (self *ChatClient2Hall) RpcChangeTeamManageSetting_(reqMsg *CgTeamManageSetting) (*CgTeamManageSetting, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcChangeTeamManageSetting", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CgTeamManageSetting), e
}
func (self *ChatClient2Hall) RpcAddTeamManager(reqMsg *client_server.TeamReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddTeamManager", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcAddTeamManager_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddTeamManager", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcDelTeamManager(reqMsg *client_server.TeamReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelTeamManager", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcDelTeamManager_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelTeamManager", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcChangeTeamOwner(reqMsg *client_server.TeamReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcChangeTeamOwner", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcChangeTeamOwner_(reqMsg *client_server.TeamReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcChangeTeamOwner", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcOpenRedPacket(reqMsg *OpenRedPacket) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcOpenRedPacket", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcOpenRedPacket_(reqMsg *OpenRedPacket) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOpenRedPacket", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcCheckRedPacket(reqMsg *OpenRedPacket) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCheckRedPacket", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcCheckRedPacket_(reqMsg *OpenRedPacket) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckRedPacket", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcTransferMoney(reqMsg *share_message.TransferMoney) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcTransferMoney", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcTransferMoney_(reqMsg *share_message.TransferMoney) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTransferMoney", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcOpenTransferMoney(reqMsg *share_message.TransferMoney) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcOpenTransferMoney", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcOpenTransferMoney_(reqMsg *share_message.TransferMoney) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOpenTransferMoney", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcAcceptAddTeam(reqMsg *client_server.TeamInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAcceptAddTeam", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcAcceptAddTeam_(reqMsg *client_server.TeamInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAcceptAddTeam", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcRefuseAddTeam(reqMsg *client_server.TeamInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcRefuseAddTeam", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcRefuseAddTeam_(reqMsg *client_server.TeamInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRefuseAddTeam", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcChangeFriendSetting(reqMsg *FriendSetting) *FriendSetting {
	msg, e := self.Sender.CallRpcMethod("RpcChangeFriendSetting", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*FriendSetting)
}

func (self *ChatClient2Hall) RpcChangeFriendSetting_(reqMsg *FriendSetting) (*FriendSetting, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcChangeFriendSetting", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*FriendSetting), e
}
func (self *ChatClient2Hall) RpcExitTeam(reqMsg *client_server.TeamInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcExitTeam", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcExitTeam_(reqMsg *client_server.TeamInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcExitTeam", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcGetTeamPlayerInfo(reqMsg *client_server.TeamInfo) *share_message.TeamPlayerInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamPlayerInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.TeamPlayerInfo)
}

func (self *ChatClient2Hall) RpcGetTeamPlayerInfo_(reqMsg *client_server.TeamInfo) (*share_message.TeamPlayerInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamPlayerInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.TeamPlayerInfo), e
}
func (self *ChatClient2Hall) RpcGetBaseTeamInfo(reqMsg *client_server.TeamReq) *TeamDataInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetBaseTeamInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TeamDataInfo)
}

func (self *ChatClient2Hall) RpcGetBaseTeamInfo_(reqMsg *client_server.TeamReq) (*TeamDataInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetBaseTeamInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TeamDataInfo), e
}
func (self *ChatClient2Hall) RpcGetPlayerCardInfo(reqMsg *client_server.PlayerReq) *share_message.TeamPlayerInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerCardInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.TeamPlayerInfo)
}

func (self *ChatClient2Hall) RpcGetPlayerCardInfo_(reqMsg *client_server.PlayerReq) (*share_message.TeamPlayerInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerCardInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.TeamPlayerInfo), e
}
func (self *ChatClient2Hall) RpcWithdrawMessage(reqMsg *client_server.ReadInfo) *LogInfo {
	msg, e := self.Sender.CallRpcMethod("RpcWithdrawMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LogInfo)
}

func (self *ChatClient2Hall) RpcWithdrawMessage_(reqMsg *client_server.ReadInfo) (*LogInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWithdrawMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LogInfo), e
}
func (self *ChatClient2Hall) RpcRequestSpecialChat(reqMsg *SpecialChatInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcRequestSpecialChat", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcRequestSpecialChat_(reqMsg *SpecialChatInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRequestSpecialChat", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcOperateSpecialChat(reqMsg *SpecialChatInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcOperateSpecialChat", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcOperateSpecialChat_(reqMsg *SpecialChatInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOperateSpecialChat", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcRequestChatInfo(reqMsg *ChatInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcRequestChatInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcRequestChatInfo_(reqMsg *ChatInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRequestChatInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcGetWaiterMsg(reqMsg *WaiterMsgRequest) *share_message.IMmessage {
	msg, e := self.Sender.CallRpcMethod("RpcGetWaiterMsg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.IMmessage)
}

func (self *ChatClient2Hall) RpcGetWaiterMsg_(reqMsg *WaiterMsgRequest) (*share_message.IMmessage, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWaiterMsg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.IMmessage), e
}
func (self *ChatClient2Hall) RpcSendWaiterMsg(reqMsg *share_message.IMmessage) *share_message.IMmessage {
	msg, e := self.Sender.CallRpcMethod("RpcSendWaiterMsg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.IMmessage)
}

func (self *ChatClient2Hall) RpcSendWaiterMsg_(reqMsg *share_message.IMmessage) (*share_message.IMmessage, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSendWaiterMsg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.IMmessage), e
}
func (self *ChatClient2Hall) RpcRequestWaiterTypes(reqMsg *base.Empty) *WaiterTypesResponse {
	msg, e := self.Sender.CallRpcMethod("RpcRequestWaiterTypes", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WaiterTypesResponse)
}

func (self *ChatClient2Hall) RpcRequestWaiterTypes_(reqMsg *base.Empty) (*WaiterTypesResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRequestWaiterTypes", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WaiterTypesResponse), e
}
func (self *ChatClient2Hall) RpcRequestWaiterService(reqMsg *WaiterMsgRequest) *WaiterMsgResponse {
	msg, e := self.Sender.CallRpcMethod("RpcRequestWaiterService", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WaiterMsgResponse)
}

func (self *ChatClient2Hall) RpcRequestWaiterService_(reqMsg *WaiterMsgRequest) (*WaiterMsgResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRequestWaiterService", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WaiterMsgResponse), e
}
func (self *ChatClient2Hall) RpcWaiterGrade(reqMsg *share_message.IMmessage) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterGrade", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcWaiterGrade_(reqMsg *share_message.IMmessage) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterGrade", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcSearchForKey(reqMsg *SearchFaqRequest) *SearchFaqResponse {
	msg, e := self.Sender.CallRpcMethod("RpcSearchForKey", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SearchFaqResponse)
}

func (self *ChatClient2Hall) RpcSearchForKey_(reqMsg *SearchFaqRequest) (*SearchFaqResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSearchForKey", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SearchFaqResponse), e
}
func (self *ChatClient2Hall) RpcOpenFaqById(reqMsg *OpenFaqRequest) *share_message.WaiterFAQ {
	msg, e := self.Sender.CallRpcMethod("RpcOpenFaqById", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.WaiterFAQ)
}

func (self *ChatClient2Hall) RpcOpenFaqById_(reqMsg *OpenFaqRequest) (*share_message.WaiterFAQ, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOpenFaqById", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.WaiterFAQ), e
}
func (self *ChatClient2Hall) RpcBroadCastQTX(reqMsg *BroadCastQTX) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcBroadCastQTX", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcBroadCastQTX_(reqMsg *BroadCastQTX) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBroadCastQTX", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcChatNew(reqMsg *share_message.Chat) *share_message.Chat {
	msg, e := self.Sender.CallRpcMethod("RpcChatNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.Chat)
}

func (self *ChatClient2Hall) RpcChatNew_(reqMsg *share_message.Chat) (*share_message.Chat, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcChatNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.Chat), e
}
func (self *ChatClient2Hall) RpcGetSessionData(reqMsg *AllSessionData) *AllSessionData {
	msg, e := self.Sender.CallRpcMethod("RpcGetSessionData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllSessionData)
}

func (self *ChatClient2Hall) RpcGetSessionData_(reqMsg *AllSessionData) (*AllSessionData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSessionData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllSessionData), e
}
func (self *ChatClient2Hall) RpcGetSessionDetail(reqMsg *SessionData) *SessionData {
	msg, e := self.Sender.CallRpcMethod("RpcGetSessionDetail", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SessionData)
}

func (self *ChatClient2Hall) RpcGetSessionDetail_(reqMsg *SessionData) (*SessionData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSessionDetail", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SessionData), e
}
func (self *ChatClient2Hall) RpcGetSessionChat(reqMsg *SessionChatData) *SessionChatData {
	msg, e := self.Sender.CallRpcMethod("RpcGetSessionChat", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SessionChatData)
}

func (self *ChatClient2Hall) RpcGetSessionChat_(reqMsg *SessionChatData) (*SessionChatData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSessionChat", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SessionChatData), e
}
func (self *ChatClient2Hall) RpcGetTeamMemberData(reqMsg *TeamMemberData) *TeamMemberData {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamMemberData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TeamMemberData)
}

func (self *ChatClient2Hall) RpcGetTeamMemberData_(reqMsg *TeamMemberData) (*TeamMemberData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamMemberData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TeamMemberData), e
}
func (self *ChatClient2Hall) RpcGetTeamDetailData(reqMsg *TeamDetailData) *TeamDetailData {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamDetailData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TeamDetailData)
}

func (self *ChatClient2Hall) RpcGetTeamDetailData_(reqMsg *TeamDetailData) (*TeamDetailData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamDetailData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TeamDetailData), e
}
func (self *ChatClient2Hall) RpcGetTeamAtData(reqMsg *TeamAtData) *TeamAtData {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamAtData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TeamAtData)
}

func (self *ChatClient2Hall) RpcGetTeamAtData_(reqMsg *TeamAtData) (*TeamAtData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamAtData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TeamAtData), e
}
func (self *ChatClient2Hall) RpcDeleteMessage(reqMsg *client_server.ReadInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ChatClient2Hall) RpcDeleteMessage_(reqMsg *client_server.ReadInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *ChatClient2Hall) RpcGetTeamSettingNew(reqMsg *client_server.TeamInfo) *share_message.TeamSetting {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamSettingNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.TeamSetting)
}

func (self *ChatClient2Hall) RpcGetTeamSettingNew_(reqMsg *client_server.TeamInfo) (*share_message.TeamSetting, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamSettingNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.TeamSetting), e
}
func (self *ChatClient2Hall) RpcGetTeamManageSettingNew(reqMsg *client_server.TeamInfo) *client_server.TeamManagerSetting {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamManageSettingNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*client_server.TeamManagerSetting)
}

func (self *ChatClient2Hall) RpcGetTeamManageSettingNew_(reqMsg *client_server.TeamInfo) (*client_server.TeamManagerSetting, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamManageSettingNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*client_server.TeamManagerSetting), e
}
func (self *ChatClient2Hall) RpcGetOneSessionData(reqMsg *SessionData) *SessionData {
	msg, e := self.Sender.CallRpcMethod("RpcGetOneSessionData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SessionData)
}

func (self *ChatClient2Hall) RpcGetOneSessionData_(reqMsg *SessionData) (*SessionData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetOneSessionData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SessionData), e
}
func (self *ChatClient2Hall) RpcCheckIsTeamMember(reqMsg *CheckTeamMember) *CheckTeamMember {
	msg, e := self.Sender.CallRpcMethod("RpcCheckIsTeamMember", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CheckTeamMember)
}

func (self *ChatClient2Hall) RpcCheckIsTeamMember_(reqMsg *CheckTeamMember) (*CheckTeamMember, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckIsTeamMember", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CheckTeamMember), e
}
func (self *ChatClient2Hall) RpcGetSaveTeamSessions(reqMsg *base.Empty) *AllSessionData {
	msg, e := self.Sender.CallRpcMethod("RpcGetSaveTeamSessions", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AllSessionData)
}

func (self *ChatClient2Hall) RpcGetSaveTeamSessions_(reqMsg *base.Empty) (*AllSessionData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSaveTeamSessions", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AllSessionData), e
}
func (self *ChatClient2Hall) RpcCheckIsMyTeamSession(reqMsg *CheckIsMySession) *CheckIsMySession {
	msg, e := self.Sender.CallRpcMethod("RpcCheckIsMyTeamSession", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CheckIsMySession)
}

func (self *ChatClient2Hall) RpcCheckIsMyTeamSession_(reqMsg *CheckIsMySession) (*CheckIsMySession, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckIsMyTeamSession", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CheckIsMySession), e
}
func (self *ChatClient2Hall) RpcGetOneShowLog(reqMsg *GetOneShowLog) *GetOneShowLog {
	msg, e := self.Sender.CallRpcMethod("RpcGetOneShowLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetOneShowLog)
}

func (self *ChatClient2Hall) RpcGetOneShowLog_(reqMsg *GetOneShowLog) (*GetOneShowLog, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetOneShowLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetOneShowLog), e
}

// ==========================================================
type IChatHall2Client interface {
	RpcChatToClient(reqMsg *share_message.Chat)
	RpcChatNewSession(reqMsg *SessionData)
	RpcTeamSettingResponse(reqMsg *client_server.TeamMsg)
	RpcTeamManageSettingResponse(reqMsg *client_server.TeamManagerSetting)
	RpcTeamOutPlayer(reqMsg *client_server.TeamInfo)
	RpcOpenRedPacketResult(reqMsg *OpenRedPacket)
	RpcCheckRedPacketResult(reqMsg *share_message.RedPacket)
	RpcOpenTransferMoneyResult(reqMsg *share_message.TransferMoney)
	RpcTeamNoticeMessage(reqMsg *share_message.TeamMessage)
	RpcRequestADDTeamInviteSuccess(reqMsg *base.Empty)
	RpcRefreshTeamInviteInfo(reqMsg *AllInviteInfo)
	RpcDealAddTeamRequest(reqMsg *DealInviteInfo)
	RpcRefreshTeamPersonalName(reqMsg *ChangeNameInfo)
	RpcTunedUpPayInfo(reqMsg *share_message.PayOrderInfo)
	RpcWithdrawMessageResponse(reqMsg *LogInfo)
	RpcRequestSpecialChatResponse(reqMsg *SpecialChatInfo)
	RpcOperateSpecialChatResponse(reqMsg *SpecialChatInfo)
	RpcReturnChatInfo(reqMsg *ReturnChatInfo)
	RpcTeamChangeInfo(reqMsg *OperatorMessage)
	RpcTeamMemChangeInfo(reqMsg *OperatorMessage)
	RpcNewWaiterMsg(reqMsg *share_message.IMmessage)
	RpcEndWaiterMsg(reqMsg *share_message.IMmessage)
	RpcBroadCastQTXResp(reqMsg *BroadCastQTX)
}

type ChatHall2Client struct {
	Sender easygo.IMessageSender
}

func (self *ChatHall2Client) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *ChatHall2Client) RpcChatToClient(reqMsg *share_message.Chat) {
	self.Sender.CallRpcMethod("RpcChatToClient", reqMsg)
}
func (self *ChatHall2Client) RpcChatNewSession(reqMsg *SessionData) {
	self.Sender.CallRpcMethod("RpcChatNewSession", reqMsg)
}
func (self *ChatHall2Client) RpcTeamSettingResponse(reqMsg *client_server.TeamMsg) {
	self.Sender.CallRpcMethod("RpcTeamSettingResponse", reqMsg)
}
func (self *ChatHall2Client) RpcTeamManageSettingResponse(reqMsg *client_server.TeamManagerSetting) {
	self.Sender.CallRpcMethod("RpcTeamManageSettingResponse", reqMsg)
}
func (self *ChatHall2Client) RpcTeamOutPlayer(reqMsg *client_server.TeamInfo) {
	self.Sender.CallRpcMethod("RpcTeamOutPlayer", reqMsg)
}
func (self *ChatHall2Client) RpcOpenRedPacketResult(reqMsg *OpenRedPacket) {
	self.Sender.CallRpcMethod("RpcOpenRedPacketResult", reqMsg)
}
func (self *ChatHall2Client) RpcCheckRedPacketResult(reqMsg *share_message.RedPacket) {
	self.Sender.CallRpcMethod("RpcCheckRedPacketResult", reqMsg)
}
func (self *ChatHall2Client) RpcOpenTransferMoneyResult(reqMsg *share_message.TransferMoney) {
	self.Sender.CallRpcMethod("RpcOpenTransferMoneyResult", reqMsg)
}
func (self *ChatHall2Client) RpcTeamNoticeMessage(reqMsg *share_message.TeamMessage) {
	self.Sender.CallRpcMethod("RpcTeamNoticeMessage", reqMsg)
}
func (self *ChatHall2Client) RpcRequestADDTeamInviteSuccess(reqMsg *base.Empty) {
	self.Sender.CallRpcMethod("RpcRequestADDTeamInviteSuccess", reqMsg)
}
func (self *ChatHall2Client) RpcRefreshTeamInviteInfo(reqMsg *AllInviteInfo) {
	self.Sender.CallRpcMethod("RpcRefreshTeamInviteInfo", reqMsg)
}
func (self *ChatHall2Client) RpcDealAddTeamRequest(reqMsg *DealInviteInfo) {
	self.Sender.CallRpcMethod("RpcDealAddTeamRequest", reqMsg)
}
func (self *ChatHall2Client) RpcRefreshTeamPersonalName(reqMsg *ChangeNameInfo) {
	self.Sender.CallRpcMethod("RpcRefreshTeamPersonalName", reqMsg)
}
func (self *ChatHall2Client) RpcTunedUpPayInfo(reqMsg *share_message.PayOrderInfo) {
	self.Sender.CallRpcMethod("RpcTunedUpPayInfo", reqMsg)
}
func (self *ChatHall2Client) RpcWithdrawMessageResponse(reqMsg *LogInfo) {
	self.Sender.CallRpcMethod("RpcWithdrawMessageResponse", reqMsg)
}
func (self *ChatHall2Client) RpcRequestSpecialChatResponse(reqMsg *SpecialChatInfo) {
	self.Sender.CallRpcMethod("RpcRequestSpecialChatResponse", reqMsg)
}
func (self *ChatHall2Client) RpcOperateSpecialChatResponse(reqMsg *SpecialChatInfo) {
	self.Sender.CallRpcMethod("RpcOperateSpecialChatResponse", reqMsg)
}
func (self *ChatHall2Client) RpcReturnChatInfo(reqMsg *ReturnChatInfo) {
	self.Sender.CallRpcMethod("RpcReturnChatInfo", reqMsg)
}
func (self *ChatHall2Client) RpcTeamChangeInfo(reqMsg *OperatorMessage) {
	self.Sender.CallRpcMethod("RpcTeamChangeInfo", reqMsg)
}
func (self *ChatHall2Client) RpcTeamMemChangeInfo(reqMsg *OperatorMessage) {
	self.Sender.CallRpcMethod("RpcTeamMemChangeInfo", reqMsg)
}
func (self *ChatHall2Client) RpcNewWaiterMsg(reqMsg *share_message.IMmessage) {
	self.Sender.CallRpcMethod("RpcNewWaiterMsg", reqMsg)
}
func (self *ChatHall2Client) RpcEndWaiterMsg(reqMsg *share_message.IMmessage) {
	self.Sender.CallRpcMethod("RpcEndWaiterMsg", reqMsg)
}
func (self *ChatHall2Client) RpcBroadCastQTXResp(reqMsg *BroadCastQTX) {
	self.Sender.CallRpcMethod("RpcBroadCastQTXResp", reqMsg)
}
