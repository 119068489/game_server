package client_hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
)

type _ = base.NoReturn

type ILoveClient2Hall interface {
	RpcGetVoiceCards(reqMsg *base.Empty) *LoveMatchResp
	RpcGetVoiceCards_(reqMsg *base.Empty) (*LoveMatchResp, easygo.IRpcInterrupt)
	RpcGetVoiceCardsNew(reqMsg *IsFirstLogin) *LoveMatchResp
	RpcGetVoiceCardsNew_(reqMsg *IsFirstLogin) (*LoveMatchResp, easygo.IRpcInterrupt)
	RpcZanVoiceCard(reqMsg *PlayerInfoReq) *VCZanResult
	RpcZanVoiceCard_(reqMsg *PlayerInfoReq) (*VCZanResult, easygo.IRpcInterrupt)
	RpcGetLoveMeList(reqMsg *LoveMeReq) *LoveMeResp
	RpcGetLoveMeList_(reqMsg *LoveMeReq) (*LoveMeResp, easygo.IRpcInterrupt)
	RpcGetMyLoveList(reqMsg *MyLoveReq) *MyLoveResp
	RpcGetMyLoveList_(reqMsg *MyLoveReq) (*MyLoveResp, easygo.IRpcInterrupt)
	RpcGetVoiceCard(reqMsg *PlayerInfoReq) *VoiceCard
	RpcGetVoiceCard_(reqMsg *PlayerInfoReq) (*VoiceCard, easygo.IRpcInterrupt)
	RpcChangeSystemBgImage(reqMsg *SystemBgReq) *SystemBgResp
	RpcChangeSystemBgImage_(reqMsg *SystemBgReq) (*SystemBgResp, easygo.IRpcInterrupt)
	RpcSysPersonalityTags(reqMsg *PersonalityTagReq) *PersonalityTagReq
	RpcSysPersonalityTags_(reqMsg *PersonalityTagReq) (*PersonalityTagReq, easygo.IRpcInterrupt)
	RpcMakeVoiceVideo(reqMsg *VoiceVideo) *VoiceVideo
	RpcMakeVoiceVideo_(reqMsg *VoiceVideo) (*VoiceVideo, easygo.IRpcInterrupt)
	RpcMixVoiceVideo(reqMsg *MixVoiceVideo) *MixVoiceVideo
	RpcMixVoiceVideo_(reqMsg *MixVoiceVideo) (*MixVoiceVideo, easygo.IRpcInterrupt)
	RpcGetVoiceVideo(reqMsg *GetVoiceVideoReq) *GetVoiceVideoResp
	RpcGetVoiceVideo_(reqMsg *GetVoiceVideoReq) (*GetVoiceVideoResp, easygo.IRpcInterrupt)
	RpcSearchVoiceVideo(reqMsg *SearchVoiceVideoReq) *SearchVoiceVideoResp
	RpcSearchVoiceVideo_(reqMsg *SearchVoiceVideoReq) (*SearchVoiceVideoResp, easygo.IRpcInterrupt)
	RpcGetVoiceCardList(reqMsg *PlayerInfoReq) *VoiceCardListResp
	RpcGetVoiceCardList_(reqMsg *PlayerInfoReq) (*VoiceCardListResp, easygo.IRpcInterrupt)
	RpcGetPersonalityTags(reqMsg *PlayerInfoReq) *PersonalityTagsResp
	RpcGetPersonalityTags_(reqMsg *PlayerInfoReq) (*PersonalityTagsResp, easygo.IRpcInterrupt)
	RpcModifyVoiceCard(reqMsg *SetVoiceCard) *SetVoiceCard
	RpcModifyVoiceCard_(reqMsg *SetVoiceCard) (*SetVoiceCard, easygo.IRpcInterrupt)
	RpcSayHiToPlayer(reqMsg *PlayerInfoReq) *base.Empty
	RpcSayHiToPlayer_(reqMsg *PlayerInfoReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetIntimacyInfo(reqMsg *PlayerInfoReq) *IntimacyInfoResp
	RpcGetIntimacyInfo_(reqMsg *PlayerInfoReq) (*IntimacyInfoResp, easygo.IRpcInterrupt)
	RpcGetAllConstellation(reqMsg *base.Empty) *ConstellationResp
	RpcGetAllConstellation_(reqMsg *base.Empty) (*ConstellationResp, easygo.IRpcInterrupt)
	RpcSetMyConstellation(reqMsg *SetConstellation) *SetConstellation
	RpcSetMyConstellation_(reqMsg *SetConstellation) (*SetConstellation, easygo.IRpcInterrupt)
	RpcGetLoveMeNewNum(reqMsg *base.Empty) *LoveMeData
	RpcGetLoveMeNewNum_(reqMsg *base.Empty) (*LoveMeData, easygo.IRpcInterrupt)
	RpcReadLoveMeLog(reqMsg *base.Empty) *base.Empty
	RpcReadLoveMeLog_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetVoiceTags(reqMsg *SearchVoiceVideoReq) *PersonalityTagsResp
	RpcGetVoiceTags_(reqMsg *SearchVoiceVideoReq) (*PersonalityTagsResp, easygo.IRpcInterrupt)
	RpcGetHotEpisode(reqMsg *SearchVoiceVideoReq) *HotEpisodeResp
	RpcGetHotEpisode_(reqMsg *SearchVoiceVideoReq) (*HotEpisodeResp, easygo.IRpcInterrupt)
	RpcGetMayLikeEpisode(reqMsg *SearchVoiceVideoReq) *SearchVoiceVideoResp
	RpcGetMayLikeEpisode_(reqMsg *SearchVoiceVideoReq) (*SearchVoiceVideoResp, easygo.IRpcInterrupt)
	RpcGetVoiceProduct(reqMsg *VoiceProduct) *SearchVoiceVideoResp
	RpcGetVoiceProduct_(reqMsg *VoiceProduct) (*SearchVoiceVideoResp, easygo.IRpcInterrupt)
	RpcDelVoiceCard(reqMsg *DelVoiceCard) *DelVoiceCard
	RpcDelVoiceCard_(reqMsg *DelVoiceCard) (*DelVoiceCard, easygo.IRpcInterrupt)
	RpcGetSessionSayHiLog(reqMsg *SayHiLog) *SayHiLog
	RpcGetSessionSayHiLog_(reqMsg *SayHiLog) (*SayHiLog, easygo.IRpcInterrupt)
}

type LoveClient2Hall struct {
	Sender easygo.IMessageSender
}

func (self *LoveClient2Hall) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *LoveClient2Hall) RpcGetVoiceCards(reqMsg *base.Empty) *LoveMatchResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceCards", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LoveMatchResp)
}

func (self *LoveClient2Hall) RpcGetVoiceCards_(reqMsg *base.Empty) (*LoveMatchResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceCards", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LoveMatchResp), e
}
func (self *LoveClient2Hall) RpcGetVoiceCardsNew(reqMsg *IsFirstLogin) *LoveMatchResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceCardsNew", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LoveMatchResp)
}

func (self *LoveClient2Hall) RpcGetVoiceCardsNew_(reqMsg *IsFirstLogin) (*LoveMatchResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceCardsNew", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LoveMatchResp), e
}
func (self *LoveClient2Hall) RpcZanVoiceCard(reqMsg *PlayerInfoReq) *VCZanResult {
	msg, e := self.Sender.CallRpcMethod("RpcZanVoiceCard", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*VCZanResult)
}

func (self *LoveClient2Hall) RpcZanVoiceCard_(reqMsg *PlayerInfoReq) (*VCZanResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcZanVoiceCard", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*VCZanResult), e
}
func (self *LoveClient2Hall) RpcGetLoveMeList(reqMsg *LoveMeReq) *LoveMeResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetLoveMeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LoveMeResp)
}

func (self *LoveClient2Hall) RpcGetLoveMeList_(reqMsg *LoveMeReq) (*LoveMeResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetLoveMeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LoveMeResp), e
}
func (self *LoveClient2Hall) RpcGetMyLoveList(reqMsg *MyLoveReq) *MyLoveResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetMyLoveList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MyLoveResp)
}

func (self *LoveClient2Hall) RpcGetMyLoveList_(reqMsg *MyLoveReq) (*MyLoveResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetMyLoveList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MyLoveResp), e
}
func (self *LoveClient2Hall) RpcGetVoiceCard(reqMsg *PlayerInfoReq) *VoiceCard {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceCard", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*VoiceCard)
}

func (self *LoveClient2Hall) RpcGetVoiceCard_(reqMsg *PlayerInfoReq) (*VoiceCard, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceCard", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*VoiceCard), e
}
func (self *LoveClient2Hall) RpcChangeSystemBgImage(reqMsg *SystemBgReq) *SystemBgResp {
	msg, e := self.Sender.CallRpcMethod("RpcChangeSystemBgImage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SystemBgResp)
}

func (self *LoveClient2Hall) RpcChangeSystemBgImage_(reqMsg *SystemBgReq) (*SystemBgResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcChangeSystemBgImage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SystemBgResp), e
}
func (self *LoveClient2Hall) RpcSysPersonalityTags(reqMsg *PersonalityTagReq) *PersonalityTagReq {
	msg, e := self.Sender.CallRpcMethod("RpcSysPersonalityTags", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PersonalityTagReq)
}

func (self *LoveClient2Hall) RpcSysPersonalityTags_(reqMsg *PersonalityTagReq) (*PersonalityTagReq, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSysPersonalityTags", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PersonalityTagReq), e
}
func (self *LoveClient2Hall) RpcMakeVoiceVideo(reqMsg *VoiceVideo) *VoiceVideo {
	msg, e := self.Sender.CallRpcMethod("RpcMakeVoiceVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*VoiceVideo)
}

func (self *LoveClient2Hall) RpcMakeVoiceVideo_(reqMsg *VoiceVideo) (*VoiceVideo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMakeVoiceVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*VoiceVideo), e
}
func (self *LoveClient2Hall) RpcMixVoiceVideo(reqMsg *MixVoiceVideo) *MixVoiceVideo {
	msg, e := self.Sender.CallRpcMethod("RpcMixVoiceVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MixVoiceVideo)
}

func (self *LoveClient2Hall) RpcMixVoiceVideo_(reqMsg *MixVoiceVideo) (*MixVoiceVideo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMixVoiceVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MixVoiceVideo), e
}
func (self *LoveClient2Hall) RpcGetVoiceVideo(reqMsg *GetVoiceVideoReq) *GetVoiceVideoResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetVoiceVideoResp)
}

func (self *LoveClient2Hall) RpcGetVoiceVideo_(reqMsg *GetVoiceVideoReq) (*GetVoiceVideoResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetVoiceVideoResp), e
}
func (self *LoveClient2Hall) RpcSearchVoiceVideo(reqMsg *SearchVoiceVideoReq) *SearchVoiceVideoResp {
	msg, e := self.Sender.CallRpcMethod("RpcSearchVoiceVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SearchVoiceVideoResp)
}

func (self *LoveClient2Hall) RpcSearchVoiceVideo_(reqMsg *SearchVoiceVideoReq) (*SearchVoiceVideoResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSearchVoiceVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SearchVoiceVideoResp), e
}
func (self *LoveClient2Hall) RpcGetVoiceCardList(reqMsg *PlayerInfoReq) *VoiceCardListResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceCardList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*VoiceCardListResp)
}

func (self *LoveClient2Hall) RpcGetVoiceCardList_(reqMsg *PlayerInfoReq) (*VoiceCardListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceCardList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*VoiceCardListResp), e
}
func (self *LoveClient2Hall) RpcGetPersonalityTags(reqMsg *PlayerInfoReq) *PersonalityTagsResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetPersonalityTags", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PersonalityTagsResp)
}

func (self *LoveClient2Hall) RpcGetPersonalityTags_(reqMsg *PlayerInfoReq) (*PersonalityTagsResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPersonalityTags", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PersonalityTagsResp), e
}
func (self *LoveClient2Hall) RpcModifyVoiceCard(reqMsg *SetVoiceCard) *SetVoiceCard {
	msg, e := self.Sender.CallRpcMethod("RpcModifyVoiceCard", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SetVoiceCard)
}

func (self *LoveClient2Hall) RpcModifyVoiceCard_(reqMsg *SetVoiceCard) (*SetVoiceCard, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcModifyVoiceCard", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SetVoiceCard), e
}
func (self *LoveClient2Hall) RpcSayHiToPlayer(reqMsg *PlayerInfoReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSayHiToPlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *LoveClient2Hall) RpcSayHiToPlayer_(reqMsg *PlayerInfoReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSayHiToPlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *LoveClient2Hall) RpcGetIntimacyInfo(reqMsg *PlayerInfoReq) *IntimacyInfoResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetIntimacyInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*IntimacyInfoResp)
}

func (self *LoveClient2Hall) RpcGetIntimacyInfo_(reqMsg *PlayerInfoReq) (*IntimacyInfoResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetIntimacyInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*IntimacyInfoResp), e
}
func (self *LoveClient2Hall) RpcGetAllConstellation(reqMsg *base.Empty) *ConstellationResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllConstellation", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ConstellationResp)
}

func (self *LoveClient2Hall) RpcGetAllConstellation_(reqMsg *base.Empty) (*ConstellationResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllConstellation", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ConstellationResp), e
}
func (self *LoveClient2Hall) RpcSetMyConstellation(reqMsg *SetConstellation) *SetConstellation {
	msg, e := self.Sender.CallRpcMethod("RpcSetMyConstellation", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SetConstellation)
}

func (self *LoveClient2Hall) RpcSetMyConstellation_(reqMsg *SetConstellation) (*SetConstellation, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSetMyConstellation", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SetConstellation), e
}
func (self *LoveClient2Hall) RpcGetLoveMeNewNum(reqMsg *base.Empty) *LoveMeData {
	msg, e := self.Sender.CallRpcMethod("RpcGetLoveMeNewNum", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LoveMeData)
}

func (self *LoveClient2Hall) RpcGetLoveMeNewNum_(reqMsg *base.Empty) (*LoveMeData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetLoveMeNewNum", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LoveMeData), e
}
func (self *LoveClient2Hall) RpcReadLoveMeLog(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReadLoveMeLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *LoveClient2Hall) RpcReadLoveMeLog_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReadLoveMeLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *LoveClient2Hall) RpcGetVoiceTags(reqMsg *SearchVoiceVideoReq) *PersonalityTagsResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceTags", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PersonalityTagsResp)
}

func (self *LoveClient2Hall) RpcGetVoiceTags_(reqMsg *SearchVoiceVideoReq) (*PersonalityTagsResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceTags", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PersonalityTagsResp), e
}
func (self *LoveClient2Hall) RpcGetHotEpisode(reqMsg *SearchVoiceVideoReq) *HotEpisodeResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetHotEpisode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*HotEpisodeResp)
}

func (self *LoveClient2Hall) RpcGetHotEpisode_(reqMsg *SearchVoiceVideoReq) (*HotEpisodeResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetHotEpisode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*HotEpisodeResp), e
}
func (self *LoveClient2Hall) RpcGetMayLikeEpisode(reqMsg *SearchVoiceVideoReq) *SearchVoiceVideoResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetMayLikeEpisode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SearchVoiceVideoResp)
}

func (self *LoveClient2Hall) RpcGetMayLikeEpisode_(reqMsg *SearchVoiceVideoReq) (*SearchVoiceVideoResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetMayLikeEpisode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SearchVoiceVideoResp), e
}
func (self *LoveClient2Hall) RpcGetVoiceProduct(reqMsg *VoiceProduct) *SearchVoiceVideoResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceProduct", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SearchVoiceVideoResp)
}

func (self *LoveClient2Hall) RpcGetVoiceProduct_(reqMsg *VoiceProduct) (*SearchVoiceVideoResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetVoiceProduct", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SearchVoiceVideoResp), e
}
func (self *LoveClient2Hall) RpcDelVoiceCard(reqMsg *DelVoiceCard) *DelVoiceCard {
	msg, e := self.Sender.CallRpcMethod("RpcDelVoiceCard", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DelVoiceCard)
}

func (self *LoveClient2Hall) RpcDelVoiceCard_(reqMsg *DelVoiceCard) (*DelVoiceCard, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelVoiceCard", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DelVoiceCard), e
}
func (self *LoveClient2Hall) RpcGetSessionSayHiLog(reqMsg *SayHiLog) *SayHiLog {
	msg, e := self.Sender.CallRpcMethod("RpcGetSessionSayHiLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SayHiLog)
}

func (self *LoveClient2Hall) RpcGetSessionSayHiLog_(reqMsg *SayHiLog) (*SayHiLog, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSessionSayHiLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SayHiLog), e
}

// ==========================================================
type ILoveHall2Client interface {
	RpcChangeIntimacy(reqMsg *IntimacyInfoResp)
	RpcDelSayHiLog(reqMsg *SayHiLog)
}

type LoveHall2Client struct {
	Sender easygo.IMessageSender
}

func (self *LoveHall2Client) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *LoveHall2Client) RpcChangeIntimacy(reqMsg *IntimacyInfoResp) {
	self.Sender.CallRpcMethod("RpcChangeIntimacy", reqMsg)
}
func (self *LoveHall2Client) RpcDelSayHiLog(reqMsg *SayHiLog) {
	self.Sender.CallRpcMethod("RpcDelSayHiLog", reqMsg)
}
