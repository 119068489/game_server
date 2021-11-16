package brower_backstage

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type IBrower2Backstage interface {
	RpcUploadFile(reqMsg *UploadRequest) *UploadResponse
	RpcUploadFile_(reqMsg *UploadRequest) (*UploadResponse, easygo.IRpcInterrupt)
	RpcDelUploadFile(reqMsg *UploadRequest) *base.Empty
	RpcDelUploadFile_(reqMsg *UploadRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcUploadFileList(reqMsg *ListRequest) *UploadListResponse
	RpcUploadFileList_(reqMsg *ListRequest) (*UploadListResponse, easygo.IRpcInterrupt)
	RpcGoogleCode(reqMsg *GoogleCodeRequest) *GoogleCodeResponse
	RpcGoogleCode_(reqMsg *GoogleCodeRequest) (*GoogleCodeResponse, easygo.IRpcInterrupt)
	RpcGetCode(reqMsg *SigninRequest) *CodeResponse
	RpcGetCode_(reqMsg *SigninRequest) (*CodeResponse, easygo.IRpcInterrupt)
	RpcSignin(reqMsg *SigninRequest) *SigninResponse
	RpcSignin_(reqMsg *SigninRequest) (*SigninResponse, easygo.IRpcInterrupt)
	RpcGetBsCode(reqMsg *SigninRequest) *CodeResponse
	RpcGetBsCode_(reqMsg *SigninRequest) (*CodeResponse, easygo.IRpcInterrupt)
	RpcVerCode(reqMsg *SigninRequest) *base.Empty
	RpcVerCode_(reqMsg *SigninRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcLogin(reqMsg *LoginRequest) *LoginResponse
	RpcLogin_(reqMsg *LoginRequest) (*LoginResponse, easygo.IRpcInterrupt)
	RpcLogout(reqMsg *base.Empty) *base.Empty
	RpcLogout_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryVersion(reqMsg *base.Empty) *VersionData
	RpcQueryVersion_(reqMsg *base.Empty) (*VersionData, easygo.IRpcInterrupt)
	RpcUpdateVersion(reqMsg *VersionData) *base.Empty
	RpcUpdateVersion_(reqMsg *VersionData) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryTfserver(reqMsg *base.Empty) *Tfserver
	RpcQueryTfserver_(reqMsg *base.Empty) (*Tfserver, easygo.IRpcInterrupt)
	RpcUpdateTfserver(reqMsg *Tfserver) *base.Empty
	RpcUpdateTfserver_(reqMsg *Tfserver) (*base.Empty, easygo.IRpcInterrupt)
	RpcManagerList(reqMsg *GetPlayerListRequest) *GetManagerListResponse
	RpcManagerList_(reqMsg *GetPlayerListRequest) (*GetManagerListResponse, easygo.IRpcInterrupt)
	RpcEditManager(reqMsg *share_message.Manager) *base.Empty
	RpcEditManager_(reqMsg *share_message.Manager) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddManager(reqMsg *share_message.Manager) *base.Empty
	RpcAddManager_(reqMsg *share_message.Manager) (*base.Empty, easygo.IRpcInterrupt)
	RpcAdminFreeze(reqMsg *QueryDataByIds) *base.Empty
	RpcAdminFreeze_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcAdminUnFreeze(reqMsg *QueryDataByIds) *base.Empty
	RpcAdminUnFreeze_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryManagerLog(reqMsg *ListRequest) *ManagerLogResponse
	RpcQueryManagerLog_(reqMsg *ListRequest) (*ManagerLogResponse, easygo.IRpcInterrupt)
	RpcManagerLogTypesKeyValue(reqMsg *base.Empty) *KeyValueResponse
	RpcManagerLogTypesKeyValue_(reqMsg *base.Empty) (*KeyValueResponse, easygo.IRpcInterrupt)
	RpcQueryManagerTypes(reqMsg *ListRequest) *ManagerTypesResponse
	RpcQueryManagerTypes_(reqMsg *ListRequest) (*ManagerTypesResponse, easygo.IRpcInterrupt)
	RpcManagerTypesKeyValue(reqMsg *base.Empty) *KeyValueResponseTag
	RpcManagerTypesKeyValue_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcEditManagerTypes(reqMsg *share_message.ManagerTypes) *base.Empty
	RpcEditManagerTypes_(reqMsg *share_message.ManagerTypes) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryDataOverview(reqMsg *base.Empty) *DataOverview
	RpcQueryDataOverview_(reqMsg *base.Empty) (*DataOverview, easygo.IRpcInterrupt)
	RpcRegisterLoginReportLine(reqMsg *ListRequest) *LineChartResponse
	RpcRegisterLoginReportLine_(reqMsg *ListRequest) (*LineChartResponse, easygo.IRpcInterrupt)
	RpcInterestTagLine(reqMsg *ListRequest) *LineChartResponse
	RpcInterestTagLine_(reqMsg *ListRequest) (*LineChartResponse, easygo.IRpcInterrupt)
	RpcPhoneBrandLine(reqMsg *base.Empty) *NameValueResponseTag
	RpcPhoneBrandLine_(reqMsg *base.Empty) (*NameValueResponseTag, easygo.IRpcInterrupt)
	RpcPlayerOnlineLine(reqMsg *ListRequest) *LineChartResponse
	RpcPlayerOnlineLine_(reqMsg *ListRequest) (*LineChartResponse, easygo.IRpcInterrupt)
	RpcPlayerLogLocation(reqMsg *ListRequest) *NameValueResponseTag
	RpcPlayerLogLocation_(reqMsg *ListRequest) (*NameValueResponseTag, easygo.IRpcInterrupt)
	RpcPlayerPortrait(reqMsg *base.Empty) *PlayerPortraitResponse
	RpcPlayerPortrait_(reqMsg *base.Empty) (*PlayerPortraitResponse, easygo.IRpcInterrupt)
	RpcQueryRolePower(reqMsg *ListRequest) *QueryRolePowerList
	RpcQueryRolePower_(reqMsg *ListRequest) (*QueryRolePowerList, easygo.IRpcInterrupt)
	RpcUpdateRolePower(reqMsg *share_message.RolePower) *base.Empty
	RpcUpdateRolePower_(reqMsg *share_message.RolePower) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteRolePower(reqMsg *QueryDataByIds) *base.Empty
	RpcDeleteRolePower_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetPowerRouter(reqMsg *QueryDataById) *share_message.RolePower
	RpcGetPowerRouter_(reqMsg *QueryDataById) (*share_message.RolePower, easygo.IRpcInterrupt)
	RpcGetRolePowerList(reqMsg *base.Empty) *KeyValueResponseTag
	RpcGetRolePowerList_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcPlayerList(reqMsg *GetPlayerListRequest) *GetPlayerListResponse
	RpcPlayerList_(reqMsg *GetPlayerListRequest) (*GetPlayerListResponse, easygo.IRpcInterrupt)
	RpcGetPlayerById(reqMsg *QueryDataById) *share_message.PlayerBase
	RpcGetPlayerById_(reqMsg *QueryDataById) (*share_message.PlayerBase, easygo.IRpcInterrupt)
	RpcGetPlayerByAccount(reqMsg *QueryDataById) *share_message.PlayerBase
	RpcGetPlayerByAccount_(reqMsg *QueryDataById) (*share_message.PlayerBase, easygo.IRpcInterrupt)
	RpcEditPlayer(reqMsg *share_message.PlayerBase) *base.Empty
	RpcEditPlayer_(reqMsg *share_message.PlayerBase) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddPlayer(reqMsg *SigninRequest) *base.Empty
	RpcAddPlayer_(reqMsg *SigninRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddWaiter(reqMsg *AddWaiterRequest) *AddWaiterResponse
	RpcAddWaiter_(reqMsg *AddWaiterRequest) (*AddWaiterResponse, easygo.IRpcInterrupt)
	RpcPlayerFreeze(reqMsg *QueryDataByIds) *base.Empty
	RpcPlayerFreeze_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcPlayerUnFreeze(reqMsg *QueryDataByIds) *base.Empty
	RpcPlayerUnFreeze_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryPlayerComplaint(reqMsg *ListRequest) *PlayerComplaintResponse
	RpcQueryPlayerComplaint_(reqMsg *ListRequest) (*PlayerComplaintResponse, easygo.IRpcInterrupt)
	RpcQueryPlayerComplaintOther(reqMsg *ListRequest) *PlayerComplaintResponse
	RpcQueryPlayerComplaintOther_(reqMsg *ListRequest) (*PlayerComplaintResponse, easygo.IRpcInterrupt)
	RpcReplyPlayerComplaint(reqMsg *share_message.PlayerComplaint) *base.Empty
	RpcReplyPlayerComplaint_(reqMsg *share_message.PlayerComplaint) (*base.Empty, easygo.IRpcInterrupt)
	RpcEditPlayerCustomTag(reqMsg *QueryDataByIds) *base.Empty
	RpcEditPlayerCustomTag_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcEditPlayerLable(reqMsg *QueryDataByIds) *base.Empty
	RpcEditPlayerLable_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcEditPersonalityTags(reqMsg *QueryDataByIds) *base.Empty
	RpcEditPersonalityTags_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetPersonalityTags(reqMsg *base.Empty) *KeyValueResponseTag
	RpcGetPersonalityTags_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcInterestTypeList(reqMsg *ListRequest) *InterestTypeResponse
	RpcInterestTypeList_(reqMsg *ListRequest) (*InterestTypeResponse, easygo.IRpcInterrupt)
	RpcEditInterestType(reqMsg *share_message.InterestType) *base.Empty
	RpcEditInterestType_(reqMsg *share_message.InterestType) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetInterestTypeList(reqMsg *base.Empty) *KeyValueResponseTag
	RpcGetInterestTypeList_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcInterestTagList(reqMsg *ListRequest) *InterestTagResponse
	RpcInterestTagList_(reqMsg *ListRequest) (*InterestTagResponse, easygo.IRpcInterrupt)
	RpcEditInterestTag(reqMsg *share_message.InterestTag) *base.Empty
	RpcEditInterestTag_(reqMsg *share_message.InterestTag) (*base.Empty, easygo.IRpcInterrupt)
	RpcInterestGroupList(reqMsg *ListRequest) *InterestGroupResponse
	RpcInterestGroupList_(reqMsg *ListRequest) (*InterestGroupResponse, easygo.IRpcInterrupt)
	RpcEditInterestGroup(reqMsg *share_message.InterestGroup) *base.Empty
	RpcEditInterestGroup_(reqMsg *share_message.InterestGroup) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelInterestGroups(reqMsg *QueryDataByIds) *base.Empty
	RpcDelInterestGroups_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetInterestTagList(reqMsg *ListRequest) *KeyValueResponseTag
	RpcGetInterestTagList_(reqMsg *ListRequest) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcCustomTagList(reqMsg *ListRequest) *CustomTagResponse
	RpcCustomTagList_(reqMsg *ListRequest) (*CustomTagResponse, easygo.IRpcInterrupt)
	RpcEditCustomTag(reqMsg *share_message.CustomTag) *base.Empty
	RpcEditCustomTag_(reqMsg *share_message.CustomTag) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetCustomTagList(reqMsg *QueryDataById) *KeyValueResponseTag
	RpcGetCustomTagList_(reqMsg *QueryDataById) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcToPlayerCustomTag(reqMsg *QueryDataByIds) *base.Empty
	RpcToPlayerCustomTag_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcGrabTagList(reqMsg *ListRequest) *GrabTagResponse
	RpcGrabTagList_(reqMsg *ListRequest) (*GrabTagResponse, easygo.IRpcInterrupt)
	RpcEditGrabTag(reqMsg *share_message.GrabTag) *base.Empty
	RpcEditGrabTag_(reqMsg *share_message.GrabTag) (*base.Empty, easygo.IRpcInterrupt)
	RpcCrawlWordsList(reqMsg *ListRequest) *CrawlWordsResponse
	RpcCrawlWordsList_(reqMsg *ListRequest) (*CrawlWordsResponse, easygo.IRpcInterrupt)
	RpcEditCrawlWords(reqMsg *share_message.CrawlWords) *base.Empty
	RpcEditCrawlWords_(reqMsg *share_message.CrawlWords) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetGrabTagList(reqMsg *base.Empty) *KeyValueResponseTag
	RpcGetGrabTagList_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcDelCrawlWords(reqMsg *QueryDataByIds) *base.Empty
	RpcDelCrawlWords_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryPlayerWordsList(reqMsg *QueryDataById) *PlayerCrawlWordsResponse
	RpcQueryPlayerWordsList_(reqMsg *QueryDataById) (*PlayerCrawlWordsResponse, easygo.IRpcInterrupt)
	RpcQueryFriendPlayerList(reqMsg *GetPlayerFriendListRequest) *GetPlayerListResponse
	RpcQueryFriendPlayerList_(reqMsg *GetPlayerFriendListRequest) (*GetPlayerListResponse, easygo.IRpcInterrupt)
	RpcQueryPlayerInfo(reqMsg *QueryDataById) *PlayerFriendInfo
	RpcQueryPlayerInfo_(reqMsg *QueryDataById) (*PlayerFriendInfo, easygo.IRpcInterrupt)
	RpcAddFriend(reqMsg *AddPlayerFriendInfo) *base.Empty
	RpcAddFriend_(reqMsg *AddPlayerFriendInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcPlayerCancleAccountList(reqMsg *ListRequest) *PlayerCancleAccountListResponse
	RpcPlayerCancleAccountList_(reqMsg *ListRequest) (*PlayerCancleAccountListResponse, easygo.IRpcInterrupt)
	RpcEditPlayerCancleAccount(reqMsg *share_message.PlayerCancleAccount) *base.Empty
	RpcEditPlayerCancleAccount_(reqMsg *share_message.PlayerCancleAccount) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryPersonalChatLog(reqMsg *ChatLogRequest) *PersonalChatLogResponse
	RpcQueryPersonalChatLog_(reqMsg *ChatLogRequest) (*PersonalChatLogResponse, easygo.IRpcInterrupt)
	RpcQueryPersonalChatLogByObj(reqMsg *ChatLogRequest) *PersonalChatLogResponse
	RpcQueryPersonalChatLogByObj_(reqMsg *ChatLogRequest) (*PersonalChatLogResponse, easygo.IRpcInterrupt)
	RpcQueryTeamChatLog(reqMsg *ChatLogRequest) *TeamChatLogResponse
	RpcQueryTeamChatLog_(reqMsg *ChatLogRequest) (*TeamChatLogResponse, easygo.IRpcInterrupt)
	RpcCheckChatLogWhitelist(reqMsg *ChatLogRequest) *CommonResponse
	RpcCheckChatLogWhitelist_(reqMsg *ChatLogRequest) (*CommonResponse, easygo.IRpcInterrupt)
	RpcQueryTeamList(reqMsg *GetTeamListRequest) *GetTeamListResponse
	RpcQueryTeamList_(reqMsg *GetTeamListRequest) (*GetTeamListResponse, easygo.IRpcInterrupt)
	RpcGetTeamById(reqMsg *QueryDataById) *share_message.TeamData
	RpcGetTeamById_(reqMsg *QueryDataById) (*share_message.TeamData, easygo.IRpcInterrupt)
	RpcEditTeam(reqMsg *share_message.TeamData) *base.Empty
	RpcEditTeam_(reqMsg *share_message.TeamData) (*base.Empty, easygo.IRpcInterrupt)
	RpcDefunctTeam(reqMsg *QueryDataById) *base.Empty
	RpcDefunctTeam_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcTeamMemberOpt(reqMsg *MemberOptRequest) *base.Empty
	RpcTeamMemberOpt_(reqMsg *MemberOptRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryTeamMember(reqMsg *TeamMemberRequest) *TeamMemberResponse
	RpcQueryTeamMember_(reqMsg *TeamMemberRequest) (*TeamMemberResponse, easygo.IRpcInterrupt)
	RpcExportChatRecord(reqMsg *ListRequest) *ExportChatRecordResponse
	RpcExportChatRecord_(reqMsg *ListRequest) (*ExportChatRecordResponse, easygo.IRpcInterrupt)
	RpcQueryTeamMessage(reqMsg *ListRequest) *ExportChatRecordResponse
	RpcQueryTeamMessage_(reqMsg *ListRequest) (*ExportChatRecordResponse, easygo.IRpcInterrupt)
	RpcCreateTeamMessage(reqMsg *CreateTeamInfo) *share_message.TeamData
	RpcCreateTeamMessage_(reqMsg *CreateTeamInfo) (*share_message.TeamData, easygo.IRpcInterrupt)
	RpcQueryTeamPlayerList(reqMsg *GetTeamPlayerListRequest) *GetTeamPlayerListResponse
	RpcQueryTeamPlayerList_(reqMsg *GetTeamPlayerListRequest) (*GetTeamPlayerListResponse, easygo.IRpcInterrupt)
	RpcTeamBan(reqMsg *QueryDataByIds) *base.Empty
	RpcTeamBan_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcTeamUnBan(reqMsg *QueryDataByIds) *base.Empty
	RpcTeamUnBan_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcTeamCloseAndOpen(reqMsg *TeamManager) *ErrMessage
	RpcTeamCloseAndOpen_(reqMsg *TeamManager) (*ErrMessage, easygo.IRpcInterrupt)
	RpcTeamMemCloseAndOpen(reqMsg *TeamManager) *ErrMessage
	RpcTeamMemCloseAndOpen_(reqMsg *TeamManager) (*ErrMessage, easygo.IRpcInterrupt)
	RpcWarnLord(reqMsg *QueryDataByIds) *base.Empty
	RpcWarnLord_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQuerySouceType(reqMsg *SourceTypeRequest) *SourceTypeResponse
	RpcQuerySouceType_(reqMsg *SourceTypeRequest) (*SourceTypeResponse, easygo.IRpcInterrupt)
	RpcQueryGeneralQuota(reqMsg *base.Empty) *share_message.GeneralQuota
	RpcQueryGeneralQuota_(reqMsg *base.Empty) (*share_message.GeneralQuota, easygo.IRpcInterrupt)
	RpcEditGeneralQuota(reqMsg *share_message.GeneralQuota) *base.Empty
	RpcEditGeneralQuota_(reqMsg *share_message.GeneralQuota) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryPaymentSetting(reqMsg *ListRequest) *PaymentSettingResponse
	RpcQueryPaymentSetting_(reqMsg *ListRequest) (*PaymentSettingResponse, easygo.IRpcInterrupt)
	RpcEditPaymentSetting(reqMsg *share_message.PaymentSetting) *base.Empty
	RpcEditPaymentSetting_(reqMsg *share_message.PaymentSetting) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelPaymentSetting(reqMsg *QueryDataByIds) *base.Empty
	RpcDelPaymentSetting_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryPayType(reqMsg *QueryDataById) *PayTypeResponse
	RpcQueryPayType_(reqMsg *QueryDataById) (*PayTypeResponse, easygo.IRpcInterrupt)
	RpcEditPayType(reqMsg *share_message.PayType) *base.Empty
	RpcEditPayType_(reqMsg *share_message.PayType) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelPayType(reqMsg *QueryDataByIds) *base.Empty
	RpcDelPayType_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryPayScene(reqMsg *base.Empty) *PaySceneResponse
	RpcQueryPayScene_(reqMsg *base.Empty) (*PaySceneResponse, easygo.IRpcInterrupt)
	RpcEditPayScene(reqMsg *share_message.PayScene) *base.Empty
	RpcEditPayScene_(reqMsg *share_message.PayScene) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelPayScene(reqMsg *QueryDataByIds) *base.Empty
	RpcDelPayScene_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryPaymentPlatform(reqMsg *PlatformChannelRequest) *PaymentPlatformResponse
	RpcQueryPaymentPlatform_(reqMsg *PlatformChannelRequest) (*PaymentPlatformResponse, easygo.IRpcInterrupt)
	RpcEditPaymentPlatform(reqMsg *share_message.PaymentPlatform) *base.Empty
	RpcEditPaymentPlatform_(reqMsg *share_message.PaymentPlatform) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelPaymentPlatform(reqMsg *QueryDataByIds) *base.Empty
	RpcDelPaymentPlatform_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryPlatformChannel(reqMsg *PlatformChannelRequest) *PlatformChannelResponse
	RpcQueryPlatformChannel_(reqMsg *PlatformChannelRequest) (*PlatformChannelResponse, easygo.IRpcInterrupt)
	RpcEditPlatformChannel(reqMsg *share_message.PlatformChannel) *base.Empty
	RpcEditPlatformChannel_(reqMsg *share_message.PlatformChannel) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelPlatformChannel(reqMsg *QueryDataByIds) *base.Empty
	RpcDelPlatformChannel_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcBatchClosePlatformChannel(reqMsg *QueryDataByIds) *base.Empty
	RpcBatchClosePlatformChannel_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddGold(reqMsg *AddGoldResult) *base.Empty
	RpcAddGold_(reqMsg *AddGoldResult) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryGoldLog(reqMsg *QueryGoldLogRequest) *QueryGoldLogResponse
	RpcQueryGoldLog_(reqMsg *QueryGoldLogRequest) (*QueryGoldLogResponse, easygo.IRpcInterrupt)
	RpcQueryOrderList(reqMsg *QueryOrderRequest) *QueryOrderResponse
	RpcQueryOrderList_(reqMsg *QueryOrderRequest) (*QueryOrderResponse, easygo.IRpcInterrupt)
	RpcOptOrder(reqMsg *OptOrderRequest) *base.Empty
	RpcOptOrder_(reqMsg *OptOrderRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcUpdateOrderList(reqMsg *base.Empty) *base.Empty
	RpcUpdateOrderList_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcCheckOrder(reqMsg *OptOrderRequest) *base.Empty
	RpcCheckOrder_(reqMsg *OptOrderRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcMakePlayerKeepReport(reqMsg *ListRequest) *PlayerKeepReportResponse
	RpcMakePlayerKeepReport_(reqMsg *ListRequest) (*PlayerKeepReportResponse, easygo.IRpcInterrupt)
	RpcPlayerKeepReport(reqMsg *ListRequest) *PlayerKeepReportResponse
	RpcPlayerKeepReport_(reqMsg *ListRequest) (*PlayerKeepReportResponse, easygo.IRpcInterrupt)
	RpcPlayerActiveReport(reqMsg *ListRequest) *PlayerActiveReportResponse
	RpcPlayerActiveReport_(reqMsg *ListRequest) (*PlayerActiveReportResponse, easygo.IRpcInterrupt)
	RpcPlayerWeekActiveReport(reqMsg *ListRequest) *PlayerActiveReportResponse
	RpcPlayerWeekActiveReport_(reqMsg *ListRequest) (*PlayerActiveReportResponse, easygo.IRpcInterrupt)
	RpcPlayerMonthActiveReport(reqMsg *ListRequest) *PlayerActiveReportResponse
	RpcPlayerMonthActiveReport_(reqMsg *ListRequest) (*PlayerActiveReportResponse, easygo.IRpcInterrupt)
	RpcPlayerBehaviorReport(reqMsg *ListRequest) *PlayerBehaviorReportResponse
	RpcPlayerBehaviorReport_(reqMsg *ListRequest) (*PlayerBehaviorReportResponse, easygo.IRpcInterrupt)
	RpcInOutCashSumReport(reqMsg *ListRequest) *InOutCashSumReportResponse
	RpcInOutCashSumReport_(reqMsg *ListRequest) (*InOutCashSumReportResponse, easygo.IRpcInterrupt)
	RpcRegisterLoginReport(reqMsg *ListRequest) *RegisterLoginReportResponse
	RpcRegisterLoginReport_(reqMsg *ListRequest) (*RegisterLoginReportResponse, easygo.IRpcInterrupt)
	RpcOperationChannelReport(reqMsg *ListRequest) *OperationChannelReportResponse
	RpcOperationChannelReport_(reqMsg *ListRequest) (*OperationChannelReportResponse, easygo.IRpcInterrupt)
	RpcChannelReport(reqMsg *ListRequest) *ChannelReportResponse
	RpcChannelReport_(reqMsg *ListRequest) (*ChannelReportResponse, easygo.IRpcInterrupt)
	RpcOperationChannelLine(reqMsg *ListRequest) *OperationChannelReportLineResponse
	RpcOperationChannelLine_(reqMsg *ListRequest) (*OperationChannelReportLineResponse, easygo.IRpcInterrupt)
	RpcArticleReport(reqMsg *ListRequest) *ArticleReportResponse
	RpcArticleReport_(reqMsg *ListRequest) (*ArticleReportResponse, easygo.IRpcInterrupt)
	RpcNoticeReport(reqMsg *ListRequest) *ArticleReportResponse
	RpcNoticeReport_(reqMsg *ListRequest) (*ArticleReportResponse, easygo.IRpcInterrupt)
	RpcSquareReport(reqMsg *ListRequest) *SquareReportResponse
	RpcSquareReport_(reqMsg *ListRequest) (*SquareReportResponse, easygo.IRpcInterrupt)
	RpcQueryActivityReport(reqMsg *ListRequest) *ActivityReportResponse
	RpcQueryActivityReport_(reqMsg *ListRequest) (*ActivityReportResponse, easygo.IRpcInterrupt)
	RpcAdvReport(reqMsg *ListRequest) *AdvReportResponse
	RpcAdvReport_(reqMsg *ListRequest) (*AdvReportResponse, easygo.IRpcInterrupt)
	RpcEditRegisterLoginReport(reqMsg *share_message.RegisterLoginReport) *base.Empty
	RpcEditRegisterLoginReport_(reqMsg *share_message.RegisterLoginReport) (*base.Empty, easygo.IRpcInterrupt)
	RpcEditPlayerKeepReport(reqMsg *share_message.PlayerKeepReport) *base.Empty
	RpcEditPlayerKeepReport_(reqMsg *share_message.PlayerKeepReport) (*base.Empty, easygo.IRpcInterrupt)
	RpcEditOperationChannelReport(reqMsg *share_message.OperationChannelReport) *base.Empty
	RpcEditOperationChannelReport_(reqMsg *share_message.OperationChannelReport) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryRecallReport(reqMsg *ListRequest) *RecallReportResponse
	RpcQueryRecallReport_(reqMsg *ListRequest) (*RecallReportResponse, easygo.IRpcInterrupt)
	RpcQueryRecallPlayerLog(reqMsg *ListRequest) *RecallplayerLogResponse
	RpcQueryRecallPlayerLog_(reqMsg *ListRequest) (*RecallplayerLogResponse, easygo.IRpcInterrupt)
	RpcCoinProductReport(reqMsg *ListRequest) *CoinProductReportResponse
	RpcCoinProductReport_(reqMsg *ListRequest) (*CoinProductReportResponse, easygo.IRpcInterrupt)
	RpcCoinProductDetailReport(reqMsg *ListRequest) *CoinProductReportResponse
	RpcCoinProductDetailReport_(reqMsg *ListRequest) (*CoinProductReportResponse, easygo.IRpcInterrupt)
	RpcNearbyAdvReport(reqMsg *ListRequest) *NearReportResponse
	RpcNearbyAdvReport_(reqMsg *ListRequest) (*NearReportResponse, easygo.IRpcInterrupt)
	RpcButtonClickReport(reqMsg *ListRequest) *ButtonClickReportResponse
	RpcButtonClickReport_(reqMsg *ListRequest) (*ButtonClickReportResponse, easygo.IRpcInterrupt)
	RpcPageRegLogReport(reqMsg *ListRequest) *PageRegLogReportResponse
	RpcPageRegLogReport_(reqMsg *ListRequest) (*PageRegLogReportResponse, easygo.IRpcInterrupt)
	RpcQueryAppPushMessage(reqMsg *QueryFeaturesRequest) *QueryFeaturesResponse
	RpcQueryAppPushMessage_(reqMsg *QueryFeaturesRequest) (*QueryFeaturesResponse, easygo.IRpcInterrupt)
	RpcEditAppPushMessage(reqMsg *share_message.AppPushMessage) *base.Empty
	RpcEditAppPushMessage_(reqMsg *share_message.AppPushMessage) (*base.Empty, easygo.IRpcInterrupt)
	RpcQuerySystemNoticeMessage(reqMsg *QueryFeaturesRequest) *QuerySystemNoticeResponse
	RpcQuerySystemNoticeMessage_(reqMsg *QueryFeaturesRequest) (*QuerySystemNoticeResponse, easygo.IRpcInterrupt)
	RpcEditSystemNoticeMessage(reqMsg *share_message.SystemNotice) *base.Empty
	RpcEditSystemNoticeMessage_(reqMsg *share_message.SystemNotice) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelSystemNoticeMessage(reqMsg *QueryDataByIds) *base.Empty
	RpcDelSystemNoticeMessage_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryTweets(reqMsg *QueryArticleOrTweetsRequest) *QueryTweetsResponse
	RpcQueryTweets_(reqMsg *QueryArticleOrTweetsRequest) (*QueryTweetsResponse, easygo.IRpcInterrupt)
	RpcAddTweets(reqMsg *share_message.Tweets) *base.Empty
	RpcAddTweets_(reqMsg *share_message.Tweets) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelTweets(reqMsg *QueryDataByIds) *base.Empty
	RpcDelTweets_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcSendTweets(reqMsg *share_message.Tweets) *base.Empty
	RpcSendTweets_(reqMsg *share_message.Tweets) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryArticle(reqMsg *QueryArticleOrTweetsRequest) *QueryArticleResponse
	RpcQueryArticle_(reqMsg *QueryArticleOrTweetsRequest) (*QueryArticleResponse, easygo.IRpcInterrupt)
	RpcEditArticle(reqMsg *share_message.Article) *base.Empty
	RpcEditArticle_(reqMsg *share_message.Article) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelArticle(reqMsg *QueryDataByIds) *base.Empty
	RpcDelArticle_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQuerySysParameterById(reqMsg *QueryDataById) *share_message.SysParameter
	RpcQuerySysParameterById_(reqMsg *QueryDataById) (*share_message.SysParameter, easygo.IRpcInterrupt)
	RpcEditSysParameter(reqMsg *share_message.SysParameter) *base.Empty
	RpcEditSysParameter_(reqMsg *share_message.SysParameter) (*base.Empty, easygo.IRpcInterrupt)
	RpcCheckShieldScore(reqMsg *CheckScoreRequest) *CheckScoreResponse
	RpcCheckShieldScore_(reqMsg *CheckScoreRequest) (*CheckScoreResponse, easygo.IRpcInterrupt)
	RpcAddRegisterPush(reqMsg *share_message.RegisterPush) *base.Empty
	RpcAddRegisterPush_(reqMsg *share_message.RegisterPush) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryRegisterPush(reqMsg *QueryArticleOrTweetsRequest) *QueryRegisterPushResponse
	RpcQueryRegisterPush_(reqMsg *QueryArticleOrTweetsRequest) (*QueryRegisterPushResponse, easygo.IRpcInterrupt)
	RpcQueryArticleComment(reqMsg *ListRequest) *ArticleCommentResponse
	RpcQueryArticleComment_(reqMsg *ListRequest) (*ArticleCommentResponse, easygo.IRpcInterrupt)
	RpcDelArticleComment(reqMsg *QueryDataByIds) *base.Empty
	RpcDelArticleComment_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryNearLead(reqMsg *ListRequest) *QueryNearSetResponse
	RpcQueryNearLead_(reqMsg *ListRequest) (*QueryNearSetResponse, easygo.IRpcInterrupt)
	RpcSaveNearLead(reqMsg *share_message.NearSet) *base.Empty
	RpcSaveNearLead_(reqMsg *share_message.NearSet) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelNearLead(reqMsg *QueryDataByIds) *base.Empty
	RpcDelNearLead_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryNearFastTerm(reqMsg *ListRequest) *QueryNearSetResponse
	RpcQueryNearFastTerm_(reqMsg *ListRequest) (*QueryNearSetResponse, easygo.IRpcInterrupt)
	RpcSaveNearFastTerm(reqMsg *share_message.NearSet) *base.Empty
	RpcSaveNearFastTerm_(reqMsg *share_message.NearSet) (*base.Empty, easygo.IRpcInterrupt)
	RpcOperationChannelList(reqMsg *ListRequest) *OperationListResponse
	RpcOperationChannelList_(reqMsg *ListRequest) (*OperationListResponse, easygo.IRpcInterrupt)
	RpcEditOperationChannel(reqMsg *share_message.OperationChannel) *base.Empty
	RpcEditOperationChannel_(reqMsg *share_message.OperationChannel) (*base.Empty, easygo.IRpcInterrupt)
	RpGetChannelList(reqMsg *base.Empty) *KeyValueResponse
	RpGetChannelList_(reqMsg *base.Empty) (*KeyValueResponse, easygo.IRpcInterrupt)
	RpcQueryDirtyWords(reqMsg *ListRequest) *DirtyWordsResponse
	RpcQueryDirtyWords_(reqMsg *ListRequest) (*DirtyWordsResponse, easygo.IRpcInterrupt)
	RpcDelDirtyWords(reqMsg *QueryDataByIds) *base.Empty
	RpcDelDirtyWords_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddDirtyWords(reqMsg *QueryDataByIds) *base.Empty
	RpcAddDirtyWords_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQuerySignature(reqMsg *ListRequest) *SignatureResponse
	RpcQuerySignature_(reqMsg *ListRequest) (*SignatureResponse, easygo.IRpcInterrupt)
	RpcDelSignature(reqMsg *QueryDataByIds) *base.Empty
	RpcDelSignature_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddSignature(reqMsg *QueryDataByIds) *base.Empty
	RpcAddSignature_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryShopItem(reqMsg *QueryShopItemRequest) *QueryShopItemResponse
	RpcQueryShopItem_(reqMsg *QueryShopItemRequest) (*QueryShopItemResponse, easygo.IRpcInterrupt)
	RpcQueryShopItemDetailById(reqMsg *QueryDataById) *QueryShopItemDetailResponse
	RpcQueryShopItemDetailById_(reqMsg *QueryDataById) (*QueryShopItemDetailResponse, easygo.IRpcInterrupt)
	RpcShopSoldOut(reqMsg *QueryDataById) *base.Empty
	RpcShopSoldOut_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetShopItemTypeDropDown(reqMsg *base.Empty) *GetShopItemTypeDropDownResponse
	RpcGetShopItemTypeDropDown_(reqMsg *base.Empty) (*GetShopItemTypeDropDownResponse, easygo.IRpcInterrupt)
	RpcReleaseShopItem(reqMsg *ReleaseEditShopItemObject) *base.Empty
	RpcReleaseShopItem_(reqMsg *ReleaseEditShopItemObject) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetEditShopItemDetailById(reqMsg *QueryDataById) *ReleaseEditShopItemObject
	RpcGetEditShopItemDetailById_(reqMsg *QueryDataById) (*ReleaseEditShopItemObject, easygo.IRpcInterrupt)
	RpcEditShopItem(reqMsg *ReleaseEditShopItemObject) *base.Empty
	RpcEditShopItem_(reqMsg *ReleaseEditShopItemObject) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryShopComment(reqMsg *QueryShopCommentRequest) *QueryShopCommentResponse
	RpcQueryShopComment_(reqMsg *QueryShopCommentRequest) (*QueryShopCommentResponse, easygo.IRpcInterrupt)
	RpcEditShopComment(reqMsg *EditShopCommentRequest) *base.Empty
	RpcEditShopComment_(reqMsg *EditShopCommentRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteShopComment(reqMsg *QueryDataById) *base.Empty
	RpcDeleteShopComment_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryShopOrder(reqMsg *QueryShopOrderRequest) *QueryShopOrderResponse
	RpcQueryShopOrder_(reqMsg *QueryShopOrderRequest) (*QueryShopOrderResponse, easygo.IRpcInterrupt)
	RpcGetExpressComDropDown(reqMsg *base.Empty) *GetExpressComDropDownResponse
	RpcGetExpressComDropDown_(reqMsg *base.Empty) (*GetExpressComDropDownResponse, easygo.IRpcInterrupt)
	RpcSendShopOrder(reqMsg *SendShopOrderRequest) *base.Empty
	RpcSendShopOrder_(reqMsg *SendShopOrderRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryShopOrderExpress(reqMsg *QueryDataById) *base.Empty
	RpcQueryShopOrderExpress_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryShopReceiveAddress(reqMsg *QueryDataById) *QueryShopReceiveAddressResponse
	RpcQueryShopReceiveAddress_(reqMsg *QueryDataById) (*QueryShopReceiveAddressResponse, easygo.IRpcInterrupt)
	RpcQueryShopDeliverAddress(reqMsg *QueryDataById) *QueryShopDeliverAddressResponse
	RpcQueryShopDeliverAddress_(reqMsg *QueryDataById) (*QueryShopDeliverAddressResponse, easygo.IRpcInterrupt)
	RpcImportShopPointCard(reqMsg *ImportShopPointCardRequest) *ImportShopPointCardResponse
	RpcImportShopPointCard_(reqMsg *ImportShopPointCardRequest) (*ImportShopPointCardResponse, easygo.IRpcInterrupt)
	RpcQueryShopPointCard(reqMsg *QueryShopPointCardRequest) *QueryShopPointCardResponse
	RpcQueryShopPointCard_(reqMsg *QueryShopPointCardRequest) (*QueryShopPointCardResponse, easygo.IRpcInterrupt)
	RpcGetShopPointCardDropDown(reqMsg *GetShopPointCardDropDownRequest) *GetShopPointCardDropDownResponse
	RpcGetShopPointCardDropDown_(reqMsg *GetShopPointCardDropDownRequest) (*GetShopPointCardDropDownResponse, easygo.IRpcInterrupt)
	RpcCancelShopOrder(reqMsg *CancelShopOrderRequest) *base.Empty
	RpcCancelShopOrder_(reqMsg *CancelShopOrderRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcCancelShopOrderForWaitSend(reqMsg *CancelShopOrderRequest) *base.Empty
	RpcCancelShopOrderForWaitSend_(reqMsg *CancelShopOrderRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcCancelShopOrderForWaitReceive(reqMsg *CancelShopOrderRequest) *base.Empty
	RpcCancelShopOrderForWaitReceive_(reqMsg *CancelShopOrderRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetIMmessageCount(reqMsg *base.Empty) *IMmessageResponse
	RpcGetIMmessageCount_(reqMsg *base.Empty) (*IMmessageResponse, easygo.IRpcInterrupt)
	RpcGetWaiterMsg(reqMsg *base.Empty) *IMmessageNopageResponse
	RpcGetWaiterMsg_(reqMsg *base.Empty) (*IMmessageNopageResponse, easygo.IRpcInterrupt)
	RpcGetWaiterMsgByMid(reqMsg *ListRequest) *share_message.IMmessage
	RpcGetWaiterMsgByMid_(reqMsg *ListRequest) (*share_message.IMmessage, easygo.IRpcInterrupt)
	RpcWaiterSendMsgToPlayer(reqMsg *share_message.IMmessage) *base.Empty
	RpcWaiterSendMsgToPlayer_(reqMsg *share_message.IMmessage) (*base.Empty, easygo.IRpcInterrupt)
	RpcWaiterOverMsgToPlayer(reqMsg *share_message.IMmessage) *share_message.IMmessage
	RpcWaiterOverMsgToPlayer_(reqMsg *share_message.IMmessage) (*share_message.IMmessage, easygo.IRpcInterrupt)
	RpcWaiterPerformanceList(reqMsg *ListRequest) *WaiterPerformanceResponse
	RpcWaiterPerformanceList_(reqMsg *ListRequest) (*WaiterPerformanceResponse, easygo.IRpcInterrupt)
	RpcWaiterPerformance(reqMsg *QueryDataById) *share_message.WaiterPerformance
	RpcWaiterPerformance_(reqMsg *QueryDataById) (*share_message.WaiterPerformance, easygo.IRpcInterrupt)
	RpcWaiterReception(reqMsg *QueryDataById) *base.Empty
	RpcWaiterReception_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcWaiterRest(reqMsg *QueryDataById) *base.Empty
	RpcWaiterRest_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcWaiterChatLogList(reqMsg *ListRequest) *IMmessageResponse
	RpcWaiterChatLogList_(reqMsg *ListRequest) (*IMmessageResponse, easygo.IRpcInterrupt)
	RpcWaiterFAQList(reqMsg *ListRequest) *WaiterFAQResponse
	RpcWaiterFAQList_(reqMsg *ListRequest) (*WaiterFAQResponse, easygo.IRpcInterrupt)
	RpcEditWaiterFAQ(reqMsg *share_message.WaiterFAQ) *base.Empty
	RpcEditWaiterFAQ_(reqMsg *share_message.WaiterFAQ) (*base.Empty, easygo.IRpcInterrupt)
	RpcWaiterFastReply(reqMsg *ListRequest) *WaiterFastReplyResponse
	RpcWaiterFastReply_(reqMsg *ListRequest) (*WaiterFastReplyResponse, easygo.IRpcInterrupt)
	RpcWaiterFastReplyNopage(reqMsg *base.Empty) *WaiterFastReplyResponse
	RpcWaiterFastReplyNopage_(reqMsg *base.Empty) (*WaiterFastReplyResponse, easygo.IRpcInterrupt)
	RpcEditWaiterFastReply(reqMsg *share_message.WaiterFastReply) *base.Empty
	RpcEditWaiterFastReply_(reqMsg *share_message.WaiterFastReply) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelWaiterFastReply(reqMsg *QueryDataByIds) *base.Empty
	RpcDelWaiterFastReply_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryDynamic(reqMsg *DynamicListRequest) *DynamicListResponse
	RpcQueryDynamic_(reqMsg *DynamicListRequest) (*DynamicListResponse, easygo.IRpcInterrupt)
	RpcQueryDynamicDetails(reqMsg *QueryDataById) *share_message.DynamicData
	RpcQueryDynamicDetails_(reqMsg *QueryDataById) (*share_message.DynamicData, easygo.IRpcInterrupt)
	RpcQueryCommentDetails(reqMsg *DynamicListRequest) *CommentList
	RpcQueryCommentDetails_(reqMsg *DynamicListRequest) (*CommentList, easygo.IRpcInterrupt)
	RpcUpdateDynamic(reqMsg *share_message.DynamicData) *base.Empty
	RpcUpdateDynamic_(reqMsg *share_message.DynamicData) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteDynamic(reqMsg *DelDynamicRequest) *base.Empty
	RpcDeleteDynamic_(reqMsg *DelDynamicRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteUnDynamic(reqMsg *DelDynamicRequest) *base.Empty
	RpcDeleteUnDynamic_(reqMsg *DelDynamicRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcShieldDynamic(reqMsg *QueryDataByIds) *base.Empty
	RpcShieldDynamic_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteCommentDatas(reqMsg *QueryDataByIds) *base.Empty
	RpcDeleteCommentDatas_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcReviewDynamic(reqMsg *QueryDataById) *base.Empty
	RpcReviewDynamic_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetTopicByIds(reqMsg *QueryDataByIds) *TopicResponse
	RpcGetTopicByIds_(reqMsg *QueryDataByIds) (*TopicResponse, easygo.IRpcInterrupt)
	RpcQueryTopicTypeList(reqMsg *ListRequest) *TopicTypeResponse
	RpcQueryTopicTypeList_(reqMsg *ListRequest) (*TopicTypeResponse, easygo.IRpcInterrupt)
	RpcUpdateTopicType(reqMsg *share_message.TopicType) *base.Empty
	RpcUpdateTopicType_(reqMsg *share_message.TopicType) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryTopicList(reqMsg *ListRequest) *TopicResponse
	RpcQueryTopicList_(reqMsg *ListRequest) (*TopicResponse, easygo.IRpcInterrupt)
	RpcUpdateTopic(reqMsg *share_message.Topic) *base.Empty
	RpcUpdateTopic_(reqMsg *share_message.Topic) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryTopicApplyList(reqMsg *ListRequest) *TopicApplyListRes
	RpcQueryTopicApplyList_(reqMsg *ListRequest) (*TopicApplyListRes, easygo.IRpcInterrupt)
	RpcQueryTopicApply(reqMsg *QueryDataById) *QueryTopicApplyRes
	RpcQueryTopicApply_(reqMsg *QueryDataById) (*QueryTopicApplyRes, easygo.IRpcInterrupt)
	RpcAuditTopicApply(reqMsg *AuditTopicApplyReq) *base.Empty
	RpcAuditTopicApply_(reqMsg *AuditTopicApplyReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcApplyTopicMasterList(reqMsg *ListRequest) *ApplyTopicMasterRes
	RpcApplyTopicMasterList_(reqMsg *ListRequest) (*ApplyTopicMasterRes, easygo.IRpcInterrupt)
	RpcApplyTopicMaster(reqMsg *QueryDataById) *base.Empty
	RpcApplyTopicMaster_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryAdvList(reqMsg *ListRequest) *AdvListResponse
	RpcQueryAdvList_(reqMsg *ListRequest) (*AdvListResponse, easygo.IRpcInterrupt)
	RpcEditAdvData(reqMsg *share_message.AdvSetting) *base.Empty
	RpcEditAdvData_(reqMsg *share_message.AdvSetting) (*base.Empty, easygo.IRpcInterrupt)
	RpcAdvOnShelf(reqMsg *QueryDataByIds) *base.Empty
	RpcAdvOnShelf_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcAdvOffShelf(reqMsg *QueryDataByIds) *base.Empty
	RpcAdvOffShelf_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcUpdateAdvSort(reqMsg *QueryDataByIds) *base.Empty
	RpcUpdateAdvSort_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelAdvData(reqMsg *QueryDataByIds) *base.Empty
	RpcDelAdvData_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcIndexTipsList(reqMsg *ListRequest) *IndexTipsResponse
	RpcIndexTipsList_(reqMsg *ListRequest) (*IndexTipsResponse, easygo.IRpcInterrupt)
	RpcSaveIndexTips(reqMsg *share_message.IndexTips) *base.Empty
	RpcSaveIndexTips_(reqMsg *share_message.IndexTips) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelIndexTips(reqMsg *QueryDataByIds) *base.Empty
	RpcDelIndexTips_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcPopSuspendList(reqMsg *ListRequest) *PopSuspendResponse
	RpcPopSuspendList_(reqMsg *ListRequest) (*PopSuspendResponse, easygo.IRpcInterrupt)
	RpcSavePopSuspendList(reqMsg *PopSuspendResponse) *base.Empty
	RpcSavePopSuspendList_(reqMsg *PopSuspendResponse) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryAdvDownList(reqMsg *ListRequest) *KeyValueResponseTag
	RpcQueryAdvDownList_(reqMsg *ListRequest) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcCoinItemList(reqMsg *ListRequest) *CoinItemListResponse
	RpcCoinItemList_(reqMsg *ListRequest) (*CoinItemListResponse, easygo.IRpcInterrupt)
	RpcSaveCoinItem(reqMsg *share_message.CoinRecharge) *base.Empty
	RpcSaveCoinItem_(reqMsg *share_message.CoinRecharge) (*base.Empty, easygo.IRpcInterrupt)
	RpcGiveCoin(reqMsg *QueryDataById) *base.Empty
	RpcGiveCoin_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryCoinChangeLog(reqMsg *ListRequest) *CoinChangeLogResponse
	RpcQueryCoinChangeLog_(reqMsg *ListRequest) (*CoinChangeLogResponse, easygo.IRpcInterrupt)
	RpcQueryPropsItemList(reqMsg *ListRequest) *PropsItemResponse
	RpcQueryPropsItemList_(reqMsg *ListRequest) (*PropsItemResponse, easygo.IRpcInterrupt)
	RpcSavePropsItem(reqMsg *share_message.PropsItem) *base.Empty
	RpcSavePropsItem_(reqMsg *share_message.PropsItem) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryPropsItemByIds(reqMsg *QueryDataByIds) *PropsItemResponse
	RpcQueryPropsItemByIds_(reqMsg *QueryDataByIds) (*PropsItemResponse, easygo.IRpcInterrupt)
	RpcCoinProductList(reqMsg *ListRequest) *CoinProductResponse
	RpcCoinProductList_(reqMsg *ListRequest) (*CoinProductResponse, easygo.IRpcInterrupt)
	RpcSaveCoinProduct(reqMsg *share_message.CoinProduct) *base.Empty
	RpcSaveCoinProduct_(reqMsg *share_message.CoinProduct) (*base.Empty, easygo.IRpcInterrupt)
	RpcPlayerBagItem(reqMsg *ListRequest) *PlayerBagItemResponse
	RpcPlayerBagItem_(reqMsg *ListRequest) (*PlayerBagItemResponse, easygo.IRpcInterrupt)
	RpcPlayerGetPropsLogList(reqMsg *ListRequest) *PlayerGetPropsLogResponse
	RpcPlayerGetPropsLogList_(reqMsg *ListRequest) (*PlayerGetPropsLogResponse, easygo.IRpcInterrupt)
	RpcRecycleProps(reqMsg *QueryDataByIds) *base.Empty
	RpcRecycleProps_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcSysGiveProps(reqMsg *QueryDataById) *base.Empty
	RpcSysGiveProps_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryCoinProductArry(reqMsg *ListRequest) *KeyValueResponse
	RpcQueryCoinProductArry_(reqMsg *ListRequest) (*KeyValueResponse, easygo.IRpcInterrupt)
	RpcBagRecycleProps(reqMsg *QueryDataById) *base.Empty
	RpcBagRecycleProps_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcToolWishBoxItemList(reqMsg *QueryDataById) *ToolWishBoxItemListRes
	RpcToolWishBoxItemList_(reqMsg *QueryDataById) (*ToolWishBoxItemListRes, easygo.IRpcInterrupt)
	RpcToolSaveWishBoxItem(reqMsg *ToolSaveWishBoxItemReq) *base.Empty
	RpcToolSaveWishBoxItem_(reqMsg *ToolSaveWishBoxItemReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcToolDelWishBoxItem(reqMsg *QueryDataByIds) *base.Empty
	RpcToolDelWishBoxItem_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcToolRateList(reqMsg *ToolRateReq) *ToolRateRes
	RpcToolRateList_(reqMsg *ToolRateReq) (*ToolRateRes, easygo.IRpcInterrupt)
	RpcToolLucky(reqMsg *ToolLuckyReq) *ToolLuckyRes
	RpcToolLucky_(reqMsg *ToolLuckyReq) (*ToolLuckyRes, easygo.IRpcInterrupt)
	RpcToolOutputData(reqMsg *QueryDataById) *ToolOutputDataRes
	RpcToolOutputData_(reqMsg *QueryDataById) (*ToolOutputDataRes, easygo.IRpcInterrupt)
	RpcToolOutputitemList(reqMsg *QueryDataById) *ToolOutputitemRes
	RpcToolOutputitemList_(reqMsg *QueryDataById) (*ToolOutputitemRes, easygo.IRpcInterrupt)
	RpcToolToolPumping(reqMsg *QueryDataById) *ToolPumping
	RpcToolToolPumping_(reqMsg *QueryDataById) (*ToolPumping, easygo.IRpcInterrupt)
	RpcToolResetWishPool(reqMsg *QueryDataById) *base.Empty
	RpcToolResetWishPool_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcToolGetWishPool(reqMsg *QueryDataById) *WishPool
	RpcToolGetWishPool_(reqMsg *QueryDataById) (*WishPool, easygo.IRpcInterrupt)
	RpcToolWishBoxList(reqMsg *ListRequest) *WishBoxList
	RpcToolWishBoxList_(reqMsg *ListRequest) (*WishBoxList, easygo.IRpcInterrupt)
	RpcToolWishBoxSave(reqMsg *share_message.WishBox) *base.Empty
	RpcToolWishBoxSave_(reqMsg *share_message.WishBox) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishBoxList(reqMsg *WishBoxListRequest) *WishBoxList
	RpcQueryWishBoxList_(reqMsg *WishBoxListRequest) (*WishBoxList, easygo.IRpcInterrupt)
	RpcUpdateWishBox(reqMsg *WishBox) *base.Empty
	RpcUpdateWishBox_(reqMsg *WishBox) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishBoxDetail(reqMsg *QueryDataById) *WishBox
	RpcGetWishBoxDetail_(reqMsg *QueryDataById) (*WishBox, easygo.IRpcInterrupt)
	RpcQueryWishBoxGoodsItemList(reqMsg *ListRequest) *WishBoxGoodsItemList
	RpcQueryWishBoxGoodsItemList_(reqMsg *ListRequest) (*WishBoxGoodsItemList, easygo.IRpcInterrupt)
	RpcQueryWishBoxWinCfgList(reqMsg *ListRequest) *WishBoxWinCfgList
	RpcQueryWishBoxWinCfgList_(reqMsg *ListRequest) (*WishBoxWinCfgList, easygo.IRpcInterrupt)
	RpcGetWishBoxKvs(reqMsg *base.Empty) *KeyValueResponseTag
	RpcGetWishBoxKvs_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcWishBoxLottery(reqMsg *WishBoxLotteryReq) *WishBoxLotteryResp
	RpcWishBoxLottery_(reqMsg *WishBoxLotteryReq) (*WishBoxLotteryResp, easygo.IRpcInterrupt)
	RpcGetGoodsListByBoxId(reqMsg *QueryDataById) *WishBoxGoodsSelectedList
	RpcGetGoodsListByBoxId_(reqMsg *QueryDataById) (*WishBoxGoodsSelectedList, easygo.IRpcInterrupt)
	RpcQueryWishGoodsList(reqMsg *WishBoxGoodsListRequest) *WishBoxGoodsList
	RpcQueryWishGoodsList_(reqMsg *WishBoxGoodsListRequest) (*WishBoxGoodsList, easygo.IRpcInterrupt)
	RpcUpdateWishGoods(reqMsg *WishBoxGoods) *base.Empty
	RpcUpdateWishGoods_(reqMsg *WishBoxGoods) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishGoodsDetail(reqMsg *QueryDataById) *WishBoxGoods
	RpcGetWishGoodsDetail_(reqMsg *QueryDataById) (*WishBoxGoods, easygo.IRpcInterrupt)
	RpcQueryWishGoodsBrandList(reqMsg *ListRequest) *WishGoodsBrandList
	RpcQueryWishGoodsBrandList_(reqMsg *ListRequest) (*WishGoodsBrandList, easygo.IRpcInterrupt)
	RpcUpdateWishGoodsBrand(reqMsg *WishGoodsBrand) *base.Empty
	RpcUpdateWishGoodsBrand_(reqMsg *WishGoodsBrand) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishGoodsBrandKvs(reqMsg *ListRequest) *KeyValueResponseTag
	RpcGetWishGoodsBrandKvs_(reqMsg *ListRequest) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcQueryWishGoodsTypeList(reqMsg *ListRequest) *WishGoodsTypeList
	RpcQueryWishGoodsTypeList_(reqMsg *ListRequest) (*WishGoodsTypeList, easygo.IRpcInterrupt)
	RpcUpdateWishGoodsType(reqMsg *WishGoodsType) *base.Empty
	RpcUpdateWishGoodsType_(reqMsg *WishGoodsType) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishGoodsTypeKvs(reqMsg *ListRequest) *KeyValueResponseTag
	RpcGetWishGoodsTypeKvs_(reqMsg *ListRequest) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcQueryWishDeliveryOrderList(reqMsg *ListRequest) *WishDeliveryOrderList
	RpcQueryWishDeliveryOrderList_(reqMsg *ListRequest) (*WishDeliveryOrderList, easygo.IRpcInterrupt)
	RpcUpdateDeliveryOrderCourierInfo(reqMsg *UpdateDeliveryOrderCourierInfo) *base.Empty
	RpcUpdateDeliveryOrderCourierInfo_(reqMsg *UpdateDeliveryOrderCourierInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcUpdateDeliveryOrderStatus(reqMsg *UpdateStatusRequest) *base.Empty
	RpcUpdateDeliveryOrderStatus_(reqMsg *UpdateStatusRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishRecycleOrderList(reqMsg *ListRequest) *WishRecycleOrderList
	RpcQueryWishRecycleOrderList_(reqMsg *ListRequest) (*WishRecycleOrderList, easygo.IRpcInterrupt)
	RpcGetWishRecycleOrderDetail(reqMsg *QueryDataById) *WishRecycleOrderDetailList
	RpcGetWishRecycleOrderDetail_(reqMsg *QueryDataById) (*WishRecycleOrderDetailList, easygo.IRpcInterrupt)
	RpcUpdateWishRecycleOrderStatus(reqMsg *UpdateStatusRequest) *base.Empty
	RpcUpdateWishRecycleOrderStatus_(reqMsg *UpdateStatusRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishRecycleOrderUserInfo(reqMsg *QueryDataById) *WishRecycleOrderUserInfo
	RpcGetWishRecycleOrderUserInfo_(reqMsg *QueryDataById) (*WishRecycleOrderUserInfo, easygo.IRpcInterrupt)
	RpcQueryWishOrderList(reqMsg *QueryWishOrderRequest) *QueryOrderResponse
	RpcQueryWishOrderList_(reqMsg *QueryWishOrderRequest) (*QueryOrderResponse, easygo.IRpcInterrupt)
	RpcOptWishOrder(reqMsg *OptOrderRequest) *base.Empty
	RpcOptWishOrder_(reqMsg *OptOrderRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcWishPlayerList(reqMsg *ListRequest) *WishPlayerListResponse
	RpcWishPlayerList_(reqMsg *ListRequest) (*WishPlayerListResponse, easygo.IRpcInterrupt)
	RpcWishPlayerFreezeDiamond(reqMsg *QueryDataByIds) *base.Empty
	RpcWishPlayerFreezeDiamond_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcWishPlayerUnFreezeDiamond(reqMsg *QueryDataByIds) *base.Empty
	RpcWishPlayerUnFreezeDiamond_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishPoolReportList(reqMsg *ListRequest) *WishPoolReportList
	RpcQueryWishPoolReportList_(reqMsg *ListRequest) (*WishPoolReportList, easygo.IRpcInterrupt)
	RpcQueryWishBoxReportList(reqMsg *ListRequest) *WishBoxReportList
	RpcQueryWishBoxReportList_(reqMsg *ListRequest) (*WishBoxReportList, easygo.IRpcInterrupt)
	RpcQueryWishBoxDetailReportList(reqMsg *ListRequest) *WishBoxDetailReportList
	RpcQueryWishBoxDetailReportList_(reqMsg *ListRequest) (*WishBoxDetailReportList, easygo.IRpcInterrupt)
	RpcQueryWishItemReportList(reqMsg *ListRequest) *WishItemReportList
	RpcQueryWishItemReportList_(reqMsg *ListRequest) (*WishItemReportList, easygo.IRpcInterrupt)
	RpcQueryTestPlayerWishItemList(reqMsg *ListRequest) *TestPlayerWishItemList
	RpcQueryTestPlayerWishItemList_(reqMsg *ListRequest) (*TestPlayerWishItemList, easygo.IRpcInterrupt)
	RpcQueryTestWishPoolLogList(reqMsg *ListRequest) *TestWishPoolLogList
	RpcQueryTestWishPoolLogList_(reqMsg *ListRequest) (*TestWishPoolLogList, easygo.IRpcInterrupt)
	RpcQueryTestWishPoolPumpLogList(reqMsg *ListRequest) *TestWishPoolPumpLogList
	RpcQueryTestWishPoolPumpLogList_(reqMsg *ListRequest) (*TestWishPoolPumpLogList, easygo.IRpcInterrupt)
	RpcQueryTestWishPoolBoxPoolInfoList(reqMsg *ListRequest) *TestWishPoolBoxPoolInfoList
	RpcQueryTestWishPoolBoxPoolInfoList_(reqMsg *ListRequest) (*TestWishPoolBoxPoolInfoList, easygo.IRpcInterrupt)
	RpcQueryDrawRecordList(reqMsg *ListRequest) *DrawRecordList
	RpcQueryDrawRecordList_(reqMsg *ListRequest) (*DrawRecordList, easygo.IRpcInterrupt)
	RpcQueryAddBoxRecordList(reqMsg *ListRequest) *AddBoxRecordList
	RpcQueryAddBoxRecordList_(reqMsg *ListRequest) (*AddBoxRecordList, easygo.IRpcInterrupt)
	RpcQueryWishGoodsRecordList(reqMsg *ListRequest) *WishGoodsRecordList
	RpcQueryWishGoodsRecordList_(reqMsg *ListRequest) (*WishGoodsRecordList, easygo.IRpcInterrupt)
	RpcQueryDrawBoxRecordList(reqMsg *ListRequest) *DrawBoxRecordList
	RpcQueryDrawBoxRecordList_(reqMsg *ListRequest) (*DrawBoxRecordList, easygo.IRpcInterrupt)
	RpcQueryHaveItemList(reqMsg *ListRequest) *HaveItemList
	RpcQueryHaveItemList_(reqMsg *ListRequest) (*HaveItemList, easygo.IRpcInterrupt)
	RpcDeleteHaveItem(reqMsg *QueryDataByIds) *HaveItemList
	RpcDeleteHaveItem_(reqMsg *QueryDataByIds) (*HaveItemList, easygo.IRpcInterrupt)
	RpcQueryWinRecordList(reqMsg *ListRequest) *WinRecordList
	RpcQueryWinRecordList_(reqMsg *ListRequest) (*WinRecordList, easygo.IRpcInterrupt)
	RpcQueryWishDelItemList(reqMsg *ListRequest) *HaveItemList
	RpcQueryWishDelItemList_(reqMsg *ListRequest) (*HaveItemList, easygo.IRpcInterrupt)
	RpcQueryWishPoolList(reqMsg *ListRequest) *WishPoolList
	RpcQueryWishPoolList_(reqMsg *ListRequest) (*WishPoolList, easygo.IRpcInterrupt)
	RpcUpdateWishPool(reqMsg *WishPool) *base.Empty
	RpcUpdateWishPool_(reqMsg *WishPool) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteWishPool(reqMsg *QueryDataById) *base.Empty
	RpcDeleteWishPool_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishPoolKvs(reqMsg *base.Empty) *KeyValueResponseTag
	RpcGetWishPoolKvs_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcUpdateDefaultWish(reqMsg *QueryDataById) *base.Empty
	RpcUpdateDefaultWish_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishPool(reqMsg *QueryDataById) *WishPool
	RpcGetWishPool_(reqMsg *QueryDataById) (*WishPool, easygo.IRpcInterrupt)
	RpcResetWishPool(reqMsg *QueryDataById) *base.Empty
	RpcResetWishPool_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcDiamondItemList(reqMsg *ListRequest) *DiamondItemListResponse
	RpcDiamondItemList_(reqMsg *ListRequest) (*DiamondItemListResponse, easygo.IRpcInterrupt)
	RpcSaveDiamondItem(reqMsg *share_message.DiamondRecharge) *base.Empty
	RpcSaveDiamondItem_(reqMsg *share_message.DiamondRecharge) (*base.Empty, easygo.IRpcInterrupt)
	RpcGiveDiamond(reqMsg *QueryDataById) *base.Empty
	RpcGiveDiamond_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryDiamondChangeLog(reqMsg *ListRequest) *DiamondChangeLogResponse
	RpcQueryDiamondChangeLog_(reqMsg *ListRequest) (*DiamondChangeLogResponse, easygo.IRpcInterrupt)
	RpcGetPriceSection(reqMsg *base.Empty) *PriceSection
	RpcGetPriceSection_(reqMsg *base.Empty) (*PriceSection, easygo.IRpcInterrupt)
	RpcUpdatePriceSection(reqMsg *PriceSection) *base.Empty
	RpcUpdatePriceSection_(reqMsg *PriceSection) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetMailSection(reqMsg *base.Empty) *WishMailSection
	RpcGetMailSection_(reqMsg *base.Empty) (*WishMailSection, easygo.IRpcInterrupt)
	RpcUpdateMailSection(reqMsg *WishMailSection) *base.Empty
	RpcUpdateMailSection_(reqMsg *WishMailSection) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishRecycleSection(reqMsg *base.Empty) *WishRecycleSection
	RpcGetWishRecycleSection_(reqMsg *base.Empty) (*WishRecycleSection, easygo.IRpcInterrupt)
	RpcUpdateWishRecycleSection(reqMsg *WishRecycleSection) *base.Empty
	RpcUpdateWishRecycleSection_(reqMsg *WishRecycleSection) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishPayWarnCfg(reqMsg *base.Empty) *WishPayWarnCfg
	RpcGetWishPayWarnCfg_(reqMsg *base.Empty) (*WishPayWarnCfg, easygo.IRpcInterrupt)
	RpcUpdateWishPayWarnCfg(reqMsg *WishPayWarnCfg) *base.Empty
	RpcUpdateWishPayWarnCfg_(reqMsg *WishPayWarnCfg) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishCoolDownConfig(reqMsg *base.Empty) *WishCoolDownConfig
	RpcGetWishCoolDownConfig_(reqMsg *base.Empty) (*WishCoolDownConfig, easygo.IRpcInterrupt)
	RpcUpdateWishCoolDownConfig(reqMsg *WishCoolDownConfig) *base.Empty
	RpcUpdateWishCoolDownConfig_(reqMsg *WishCoolDownConfig) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishCurrencyConversionCfg(reqMsg *base.Empty) *WishCurrencyConversionCfg
	RpcGetWishCurrencyConversionCfg_(reqMsg *base.Empty) (*WishCurrencyConversionCfg, easygo.IRpcInterrupt)
	RpcUpdateWishCurrencyConversionCfg(reqMsg *WishCurrencyConversionCfg) *base.Empty
	RpcUpdateWishCurrencyConversionCfg_(reqMsg *WishCurrencyConversionCfg) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetWishGuardianCfg(reqMsg *base.Empty) *WishGuardianCfg
	RpcGetWishGuardianCfg_(reqMsg *base.Empty) (*WishGuardianCfg, easygo.IRpcInterrupt)
	RpcUpdateWishGuardianCfg(reqMsg *WishGuardianCfg) *base.Empty
	RpcUpdateWishGuardianCfg_(reqMsg *WishGuardianCfg) (*base.Empty, easygo.IRpcInterrupt)
	RpcSaveRecycleNoteCfg(reqMsg *RecycleNoteCfg) *base.Empty
	RpcSaveRecycleNoteCfg_(reqMsg *RecycleNoteCfg) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetRecycleNoteCfg(reqMsg *base.Empty) *RecycleNoteCfg
	RpcGetRecycleNoteCfg_(reqMsg *base.Empty) (*RecycleNoteCfg, easygo.IRpcInterrupt)
	RpcPayPlayerLocation(reqMsg *ListRequest) *NameValueResponseTag
	RpcPayPlayerLocation_(reqMsg *ListRequest) (*NameValueResponseTag, easygo.IRpcInterrupt)
	RpcWishCoinRechargeActivityCfgList(reqMsg *ListRequest) *WishCoinRechargeActivityCfgRes
	RpcWishCoinRechargeActivityCfgList_(reqMsg *ListRequest) (*WishCoinRechargeActivityCfgRes, easygo.IRpcInterrupt)
	RpcWishCoinRechargeActivityCfgUpdate(reqMsg *share_message.WishCoinRechargeActivityCfg) *base.Empty
	RpcWishCoinRechargeActivityCfgUpdate_(reqMsg *share_message.WishCoinRechargeActivityCfg) (*base.Empty, easygo.IRpcInterrupt)
	RpcWishCoinRechargeActivityCfgDel(reqMsg *QueryDataByIds) *base.Empty
	RpcWishCoinRechargeActivityCfgDel_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishActPool(reqMsg *ListRequest) *WishActPoolList
	RpcQueryWishActPool_(reqMsg *ListRequest) (*WishActPoolList, easygo.IRpcInterrupt)
	RpcUpdateWishActPool(reqMsg *share_message.WishActPool) *base.Empty
	RpcUpdateWishActPool_(reqMsg *share_message.WishActPool) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteWishActPool(reqMsg *QueryDataByIds) *base.Empty
	RpcDeleteWishActPool_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcAddWishAllowList(reqMsg *AddWishAllowListReq) *base.Empty
	RpcAddWishAllowList_(reqMsg *AddWishAllowListReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteWishAllowList(reqMsg *QueryDataByIds) *base.Empty
	RpcDeleteWishAllowList_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcWishAllowList(reqMsg *base.Empty) *WishAllowListResp
	RpcWishAllowList_(reqMsg *base.Empty) (*WishAllowListResp, easygo.IRpcInterrupt)
	RpcQueryWishActPoolDetail(reqMsg *ListRequest) *WishActPoolDetail
	RpcQueryWishActPoolDetail_(reqMsg *ListRequest) (*WishActPoolDetail, easygo.IRpcInterrupt)
	RpcGetWishActPoolTypeKvs(reqMsg *base.Empty) *KeyValueResponseTag
	RpcGetWishActPoolTypeKvs_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt)
	RpcGetWishActCfg(reqMsg *ListRequest) *share_message.Activity
	RpcGetWishActCfg_(reqMsg *ListRequest) (*share_message.Activity, easygo.IRpcInterrupt)
	RpcUpdateWishActCfg(reqMsg *share_message.Activity) *base.Empty
	RpcUpdateWishActCfg_(reqMsg *share_message.Activity) (*base.Empty, easygo.IRpcInterrupt)
	RpcWishPoolActivityReportList(reqMsg *ListRequest) *WishPoolActivityReportResp
	RpcWishPoolActivityReportList_(reqMsg *ListRequest) (*WishPoolActivityReportResp, easygo.IRpcInterrupt)
	RpcQueryWishActPoolRuleDay(reqMsg *ListRequest) *WishActPoolRuleList
	RpcQueryWishActPoolRuleDay_(reqMsg *ListRequest) (*WishActPoolRuleList, easygo.IRpcInterrupt)
	RpcAddWishActPoolRuleDay(reqMsg *AddWishActPoolRuleRequest) *base.Empty
	RpcAddWishActPoolRuleDay_(reqMsg *AddWishActPoolRuleRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcUpdateWishActPoolRuleDay(reqMsg *WishActPoolRule) *base.Empty
	RpcUpdateWishActPoolRuleDay_(reqMsg *WishActPoolRule) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteWishActPoolRuleDay(reqMsg *QueryDataByIds) *base.Empty
	RpcDeleteWishActPoolRuleDay_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishActPoolRuleCount(reqMsg *ListRequest) *WishActPoolRuleList
	RpcQueryWishActPoolRuleCount_(reqMsg *ListRequest) (*WishActPoolRuleList, easygo.IRpcInterrupt)
	RpcAddWishActPoolRuleCount(reqMsg *AddWishActPoolRuleRequest) *base.Empty
	RpcAddWishActPoolRuleCount_(reqMsg *AddWishActPoolRuleRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcUpdateWishActPoolRuleCount(reqMsg *WishActPoolRule) *base.Empty
	RpcUpdateWishActPoolRuleCount_(reqMsg *WishActPoolRule) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteWishActPoolRuleCount(reqMsg *QueryDataByIds) *base.Empty
	RpcDeleteWishActPoolRuleCount_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishActPoolRuleWeekMonth(reqMsg *ListRequest) *WishActPoolRuleList
	RpcQueryWishActPoolRuleWeekMonth_(reqMsg *ListRequest) (*WishActPoolRuleList, easygo.IRpcInterrupt)
	RpcAddWishActPoolRuleWeekMonth(reqMsg *AddWishActPoolRuleRequest) *base.Empty
	RpcAddWishActPoolRuleWeekMonth_(reqMsg *AddWishActPoolRuleRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcUpdateWishActPoolRuleItemWeekMonth(reqMsg *WishActPoolAwardItem) *base.Empty
	RpcUpdateWishActPoolRuleItemWeekMonth_(reqMsg *WishActPoolAwardItem) (*base.Empty, easygo.IRpcInterrupt)
	RpcDeleteWishActPoolRuleItemWeekMonth(reqMsg *QueryDataByIds) *base.Empty
	RpcDeleteWishActPoolRuleItemWeekMonth_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishActPlayerRecordList(reqMsg *ListRequest) *WishActPlayerRecordList
	RpcQueryWishActPlayerRecordList_(reqMsg *ListRequest) (*WishActPlayerRecordList, easygo.IRpcInterrupt)
	RpcQueryWishActPlayerWinRecordList(reqMsg *ListRequest) *WishActPlayerWinRecordList
	RpcQueryWishActPlayerWinRecordList_(reqMsg *ListRequest) (*WishActPlayerWinRecordList, easygo.IRpcInterrupt)
	RpcQueryWishActPlayerDrawRecordList(reqMsg *ListRequest) *WishActPlayerDrawRecordList
	RpcQueryWishActPlayerDrawRecordList_(reqMsg *ListRequest) (*WishActPlayerDrawRecordList, easygo.IRpcInterrupt)
	RpcQueryWishLogReport(reqMsg *ListRequest) *QueryWishLogReportRes
	RpcQueryWishLogReport_(reqMsg *ListRequest) (*QueryWishLogReportRes, easygo.IRpcInterrupt)
	RpcQueryWishActivityPrizeLog(reqMsg *ListRequest) *QueryWishActivityPrizeLogRes
	RpcQueryWishActivityPrizeLog_(reqMsg *ListRequest) (*QueryWishActivityPrizeLogRes, easygo.IRpcInterrupt)
	RpcPlayerCardList(reqMsg *ListRequest) *InterestTagResponse
	RpcPlayerCardList_(reqMsg *ListRequest) (*InterestTagResponse, easygo.IRpcInterrupt)
	RpcUpdatePlayerCard(reqMsg *QueryDataById) *base.Empty
	RpcUpdatePlayerCard_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcCharacterTagList(reqMsg *ListRequest) *InterestTagResponse
	RpcCharacterTagList_(reqMsg *ListRequest) (*InterestTagResponse, easygo.IRpcInterrupt)
	RpcSaveCharacterTag(reqMsg *share_message.InterestTag) *base.Empty
	RpcSaveCharacterTag_(reqMsg *share_message.InterestTag) (*base.Empty, easygo.IRpcInterrupt)
	RpcPlayerVoiceWorkList(reqMsg *ListRequest) *VoiceWorkListResponse
	RpcPlayerVoiceWorkList_(reqMsg *ListRequest) (*VoiceWorkListResponse, easygo.IRpcInterrupt)
	RpcReviewePlayerVoiceWork(reqMsg *QueryDataById) *base.Empty
	RpcReviewePlayerVoiceWork_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelPlayerVoiceWork(reqMsg *QueryDataByIds) *base.Empty
	RpcDelPlayerVoiceWork_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcUploadPlayerVoiceWork(reqMsg *share_message.PlayerMixVoiceVideo) *base.Empty
	RpcUploadPlayerVoiceWork_(reqMsg *share_message.PlayerMixVoiceVideo) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetPlayerVoiceWorkUse(reqMsg *QueryDataById) *share_message.PlayerMixVoiceVideo
	RpcGetPlayerVoiceWorkUse_(reqMsg *QueryDataById) (*share_message.PlayerMixVoiceVideo, easygo.IRpcInterrupt)
	RpcBgTagList(reqMsg *ListRequest) *InterestTagResponse
	RpcBgTagList_(reqMsg *ListRequest) (*InterestTagResponse, easygo.IRpcInterrupt)
	RpcUpdateBgTag(reqMsg *share_message.InterestTag) *base.Empty
	RpcUpdateBgTag_(reqMsg *share_message.InterestTag) (*base.Empty, easygo.IRpcInterrupt)
	RpcBgVoiceVideoList(reqMsg *ListRequest) *BgVoiceVideoListResponse
	RpcBgVoiceVideoList_(reqMsg *ListRequest) (*BgVoiceVideoListResponse, easygo.IRpcInterrupt)
	RpcUpdateBgVoiceVideo(reqMsg *share_message.BgVoiceVideo) *base.Empty
	RpcUpdateBgVoiceVideo_(reqMsg *share_message.BgVoiceVideo) (*base.Empty, easygo.IRpcInterrupt)
	RpcRevieweBgVoiceVideo(reqMsg *QueryDataById) *base.Empty
	RpcRevieweBgVoiceVideo_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelBgVoiceVideo(reqMsg *QueryDataByIds) *base.Empty
	RpcDelBgVoiceVideo_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcMatchGuideList(reqMsg *ListRequest) *MatchGuideListResponse
	RpcMatchGuideList_(reqMsg *ListRequest) (*MatchGuideListResponse, easygo.IRpcInterrupt)
	RpcUpdateMatchGuide(reqMsg *QueryDataByIds) *base.Empty
	RpcUpdateMatchGuide_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelMatchGuide(reqMsg *QueryDataByIds) *base.Empty
	RpcDelMatchGuide_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcSayHiList(reqMsg *ListRequest) *MatchGuideListResponse
	RpcSayHiList_(reqMsg *ListRequest) (*MatchGuideListResponse, easygo.IRpcInterrupt)
	RpcUpdateSayHi(reqMsg *QueryDataByIds) *base.Empty
	RpcUpdateSayHi_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelSayHi(reqMsg *QueryDataByIds) *base.Empty
	RpcDelSayHi_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcSystemBgImageList(reqMsg *ListRequest) *SystemBgImageListResponse
	RpcSystemBgImageList_(reqMsg *ListRequest) (*SystemBgImageListResponse, easygo.IRpcInterrupt)
	RpcSaveSystemBgImage(reqMsg *QueryDataByIds) *base.Empty
	RpcSaveSystemBgImage_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelSystemBgImage(reqMsg *QueryDataByIds) *base.Empty
	RpcDelSystemBgImage_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryIntimacyConfig(reqMsg *base.Empty) *IntimacyConfigRes
	RpcQueryIntimacyConfig_(reqMsg *base.Empty) (*IntimacyConfigRes, easygo.IRpcInterrupt)
	RpcUpdateIntimacyConfig(reqMsg *IntimacyConfigRes) *base.Empty
	RpcUpdateIntimacyConfig_(reqMsg *IntimacyConfigRes) (*base.Empty, easygo.IRpcInterrupt)
	RpcVCBuryingPointReport(reqMsg *ListRequest) *VCBuryingPointReportRes
	RpcVCBuryingPointReport_(reqMsg *ListRequest) (*VCBuryingPointReportRes, easygo.IRpcInterrupt)
	RpcCrawlPull(reqMsg *QueryDataById) *base.Empty
	RpcCrawlPull_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcCrawlJobList(reqMsg *ListRequest) *CrawlJobResponse
	RpcCrawlJobList_(reqMsg *ListRequest) (*CrawlJobResponse, easygo.IRpcInterrupt)
	RpcNewsSource(reqMsg *ListRequest) *NewsSourceResponse
	RpcNewsSource_(reqMsg *ListRequest) (*NewsSourceResponse, easygo.IRpcInterrupt)
	RpcVideoSource(reqMsg *ListRequest) *VideoSourceResponse
	RpcVideoSource_(reqMsg *ListRequest) (*VideoSourceResponse, easygo.IRpcInterrupt)
	RpcNewsList(reqMsg *ListRequest) *NewsSourceResponse
	RpcNewsList_(reqMsg *ListRequest) (*NewsSourceResponse, easygo.IRpcInterrupt)
	RpcSaveNews(reqMsg *share_message.TableESPortsRealTimeInfo) *base.Empty
	RpcSaveNews_(reqMsg *share_message.TableESPortsRealTimeInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelNewsSource(reqMsg *QueryDataById) *base.Empty
	RpcDelNewsSource_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelNews(reqMsg *QueryDataByIds) *base.Empty
	RpcDelNews_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcVideoList(reqMsg *ListRequest) *VideoSourceResponse
	RpcVideoList_(reqMsg *ListRequest) (*VideoSourceResponse, easygo.IRpcInterrupt)
	RpcSaveVideo(reqMsg *share_message.TableESPortsVideoInfo) *base.Empty
	RpcSaveVideo_(reqMsg *share_message.TableESPortsVideoInfo) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelVideoSource(reqMsg *QueryDataById) *base.Empty
	RpcDelVideoSource_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelVideo(reqMsg *QueryDataByIds) *base.Empty
	RpcDelVideo_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcChekVideo(reqMsg *QueryDataById) *base.Empty
	RpcChekVideo_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt)
	RpcBanVideo(reqMsg *QueryDataByIds) *base.Empty
	RpcBanVideo_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetGameList(reqMsg *ListRequest) *GameListResponse
	RpcGetGameList_(reqMsg *ListRequest) (*GameListResponse, easygo.IRpcInterrupt)
	RpcGetGameGuess(reqMsg *GameGuessRequest) *GameGuessResponse
	RpcGetGameGuess_(reqMsg *GameGuessRequest) (*GameGuessResponse, easygo.IRpcInterrupt)
	RpcEditGameGuess(reqMsg *EditGameGuessRequest) *base.Empty
	RpcEditGameGuess_(reqMsg *EditGameGuessRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetGameTeamInfo(reqMsg *GameGuessRequest) *GameTeamInfoResponse
	RpcGetGameTeamInfo_(reqMsg *GameGuessRequest) (*GameTeamInfoResponse, easygo.IRpcInterrupt)
	RpcGetGameRealTimeData(reqMsg *GameGuessRequest) *GameRealTimeResponse
	RpcGetGameRealTimeData_(reqMsg *GameGuessRequest) (*GameRealTimeResponse, easygo.IRpcInterrupt)
	RpcEditGameRealTimeData(reqMsg *EditGameRealTimeRequest) *base.Empty
	RpcEditGameRealTimeData_(reqMsg *EditGameRealTimeRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcCommentList(reqMsg *ListRequest) *CommentListResponse
	RpcCommentList_(reqMsg *ListRequest) (*CommentListResponse, easygo.IRpcInterrupt)
	RpcDelComment(reqMsg *CommentDelRequest) *base.Empty
	RpcDelComment_(reqMsg *CommentDelRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcUploadComment(reqMsg *CommentUploadRequest) *base.Empty
	RpcUploadComment_(reqMsg *CommentUploadRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetAppLabel(reqMsg *ListRequest) *SysLabelResponse
	RpcGetAppLabel_(reqMsg *ListRequest) (*SysLabelResponse, easygo.IRpcInterrupt)
	RpcSaveAppLabel(reqMsg *share_message.TableESPortsLabel) *base.Empty
	RpcSaveAppLabel_(reqMsg *share_message.TableESPortsLabel) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetSysLabel(reqMsg *ListRequest) *SysLabelResponse
	RpcGetSysLabel_(reqMsg *ListRequest) (*SysLabelResponse, easygo.IRpcInterrupt)
	RpcSaveSysLabel(reqMsg *share_message.TableESPortsLabel) *base.Empty
	RpcSaveSysLabel_(reqMsg *share_message.TableESPortsLabel) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetCarouselList(reqMsg *ListRequest) *CarouselResponse
	RpcGetCarouselList_(reqMsg *ListRequest) (*CarouselResponse, easygo.IRpcInterrupt)
	RpcSaveCarousel(reqMsg *share_message.TableESPortsCarousel) *base.Empty
	RpcSaveCarousel_(reqMsg *share_message.TableESPortsCarousel) (*base.Empty, easygo.IRpcInterrupt)
	RpcSportSysNotice(reqMsg *ListRequest) *SportSysNoticeResponse
	RpcSportSysNotice_(reqMsg *ListRequest) (*SportSysNoticeResponse, easygo.IRpcInterrupt)
	RpcSendSportSysNotice(reqMsg *share_message.TableESPortsSysMsg) *base.Empty
	RpcSendSportSysNotice_(reqMsg *share_message.TableESPortsSysMsg) (*base.Empty, easygo.IRpcInterrupt)
	RpcBetSlipList(reqMsg *ListRequest) *BetSlipListResponse
	RpcBetSlipList_(reqMsg *ListRequest) (*BetSlipListResponse, easygo.IRpcInterrupt)
	RpcBetWinLosStatistics(reqMsg *ListRequest) *BetSlipStatisticsResponse
	RpcBetWinLosStatistics_(reqMsg *ListRequest) (*BetSlipStatisticsResponse, easygo.IRpcInterrupt)
	RpcBetGameStatistics(reqMsg *ListRequest) *BetSlipStatisticsResponse
	RpcBetGameStatistics_(reqMsg *ListRequest) (*BetSlipStatisticsResponse, easygo.IRpcInterrupt)
	RpcBetSlipReportLine(reqMsg *ListRequest) *LineChartsResponse
	RpcBetSlipReportLine_(reqMsg *ListRequest) (*LineChartsResponse, easygo.IRpcInterrupt)
	RpcBetSlipReportBar(reqMsg *ListRequest) *LineChartsResponse
	RpcBetSlipReportBar_(reqMsg *ListRequest) (*LineChartsResponse, easygo.IRpcInterrupt)
	RpcBetSlipOperate(reqMsg *RpcBetSlipOperateRequest) *base.Empty
	RpcBetSlipOperate_(reqMsg *RpcBetSlipOperateRequest) (*base.Empty, easygo.IRpcInterrupt)
	RpcGiveWhiteList(reqMsg *ListRequest) *GiveWhiteListRes
	RpcGiveWhiteList_(reqMsg *ListRequest) (*GiveWhiteListRes, easygo.IRpcInterrupt)
	RpcAddGiveWhiteList(reqMsg *QueryDataByIds) *base.Empty
	RpcAddGiveWhiteList_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelGiveWhiteList(reqMsg *QueryDataByIds) *base.Empty
	RpcDelGiveWhiteList_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryRechargeEsAct(reqMsg *base.Empty) *share_message.Activity
	RpcQueryRechargeEsAct_(reqMsg *base.Empty) (*share_message.Activity, easygo.IRpcInterrupt)
	RpcUpdateRechargeEsAct(reqMsg *share_message.Activity) *base.Empty
	RpcUpdateRechargeEsAct_(reqMsg *share_message.Activity) (*base.Empty, easygo.IRpcInterrupt)
	RpcRechargeEsCfg(reqMsg *base.Empty) *RechargeEsCfgRes
	RpcRechargeEsCfg_(reqMsg *base.Empty) (*RechargeEsCfgRes, easygo.IRpcInterrupt)
	RpcSaveRechargeEsCfg(reqMsg *share_message.TableESportsExchangeCfg) *base.Empty
	RpcSaveRechargeEsCfg_(reqMsg *share_message.TableESportsExchangeCfg) (*base.Empty, easygo.IRpcInterrupt)
	RpcDelRechargeEsCfg(reqMsg *QueryDataByIds) *base.Empty
	RpcDelRechargeEsCfg_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt)
	RpcPointsReportList(reqMsg *ListRequest) *PointsReportRes
	RpcPointsReportList_(reqMsg *ListRequest) (*PointsReportRes, easygo.IRpcInterrupt)
	RpcQueryeSportCoinLog(reqMsg *ListRequest) *SportCoinLogResponse
	RpcQueryeSportCoinLog_(reqMsg *ListRequest) (*SportCoinLogResponse, easygo.IRpcInterrupt)
}

type Brower2Backstage struct {
	Sender easygo.IMessageSender
}

func (self *Brower2Backstage) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Brower2Backstage) RpcUploadFile(reqMsg *UploadRequest) *UploadResponse {
	msg, e := self.Sender.CallRpcMethod("RpcUploadFile", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*UploadResponse)
}

func (self *Brower2Backstage) RpcUploadFile_(reqMsg *UploadRequest) (*UploadResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUploadFile", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*UploadResponse), e
}
func (self *Brower2Backstage) RpcDelUploadFile(reqMsg *UploadRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelUploadFile", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelUploadFile_(reqMsg *UploadRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelUploadFile", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcUploadFileList(reqMsg *ListRequest) *UploadListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcUploadFileList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*UploadListResponse)
}

func (self *Brower2Backstage) RpcUploadFileList_(reqMsg *ListRequest) (*UploadListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUploadFileList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*UploadListResponse), e
}
func (self *Brower2Backstage) RpcGoogleCode(reqMsg *GoogleCodeRequest) *GoogleCodeResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGoogleCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GoogleCodeResponse)
}

func (self *Brower2Backstage) RpcGoogleCode_(reqMsg *GoogleCodeRequest) (*GoogleCodeResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGoogleCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GoogleCodeResponse), e
}
func (self *Brower2Backstage) RpcGetCode(reqMsg *SigninRequest) *CodeResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CodeResponse)
}

func (self *Brower2Backstage) RpcGetCode_(reqMsg *SigninRequest) (*CodeResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CodeResponse), e
}
func (self *Brower2Backstage) RpcSignin(reqMsg *SigninRequest) *SigninResponse {
	msg, e := self.Sender.CallRpcMethod("RpcSignin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SigninResponse)
}

func (self *Brower2Backstage) RpcSignin_(reqMsg *SigninRequest) (*SigninResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSignin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SigninResponse), e
}
func (self *Brower2Backstage) RpcGetBsCode(reqMsg *SigninRequest) *CodeResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetBsCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CodeResponse)
}

func (self *Brower2Backstage) RpcGetBsCode_(reqMsg *SigninRequest) (*CodeResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetBsCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CodeResponse), e
}
func (self *Brower2Backstage) RpcVerCode(reqMsg *SigninRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcVerCode", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcVerCode_(reqMsg *SigninRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcVerCode", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcLogin(reqMsg *LoginRequest) *LoginResponse {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LoginResponse)
}

func (self *Brower2Backstage) RpcLogin_(reqMsg *LoginRequest) (*LoginResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LoginResponse), e
}
func (self *Brower2Backstage) RpcLogout(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcLogout", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcLogout_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLogout", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryVersion(reqMsg *base.Empty) *VersionData {
	msg, e := self.Sender.CallRpcMethod("RpcQueryVersion", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*VersionData)
}

func (self *Brower2Backstage) RpcQueryVersion_(reqMsg *base.Empty) (*VersionData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryVersion", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*VersionData), e
}
func (self *Brower2Backstage) RpcUpdateVersion(reqMsg *VersionData) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateVersion", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateVersion_(reqMsg *VersionData) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateVersion", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryTfserver(reqMsg *base.Empty) *Tfserver {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTfserver", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*Tfserver)
}

func (self *Brower2Backstage) RpcQueryTfserver_(reqMsg *base.Empty) (*Tfserver, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTfserver", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*Tfserver), e
}
func (self *Brower2Backstage) RpcUpdateTfserver(reqMsg *Tfserver) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateTfserver", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateTfserver_(reqMsg *Tfserver) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateTfserver", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcManagerList(reqMsg *GetPlayerListRequest) *GetManagerListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcManagerList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetManagerListResponse)
}

func (self *Brower2Backstage) RpcManagerList_(reqMsg *GetPlayerListRequest) (*GetManagerListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcManagerList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetManagerListResponse), e
}
func (self *Brower2Backstage) RpcEditManager(reqMsg *share_message.Manager) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditManager", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditManager_(reqMsg *share_message.Manager) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditManager", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAddManager(reqMsg *share_message.Manager) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddManager", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddManager_(reqMsg *share_message.Manager) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddManager", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAdminFreeze(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAdminFreeze", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAdminFreeze_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAdminFreeze", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAdminUnFreeze(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAdminUnFreeze", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAdminUnFreeze_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAdminUnFreeze", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryManagerLog(reqMsg *ListRequest) *ManagerLogResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryManagerLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ManagerLogResponse)
}

func (self *Brower2Backstage) RpcQueryManagerLog_(reqMsg *ListRequest) (*ManagerLogResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryManagerLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ManagerLogResponse), e
}
func (self *Brower2Backstage) RpcManagerLogTypesKeyValue(reqMsg *base.Empty) *KeyValueResponse {
	msg, e := self.Sender.CallRpcMethod("RpcManagerLogTypesKeyValue", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponse)
}

func (self *Brower2Backstage) RpcManagerLogTypesKeyValue_(reqMsg *base.Empty) (*KeyValueResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcManagerLogTypesKeyValue", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponse), e
}
func (self *Brower2Backstage) RpcQueryManagerTypes(reqMsg *ListRequest) *ManagerTypesResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryManagerTypes", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ManagerTypesResponse)
}

func (self *Brower2Backstage) RpcQueryManagerTypes_(reqMsg *ListRequest) (*ManagerTypesResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryManagerTypes", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ManagerTypesResponse), e
}
func (self *Brower2Backstage) RpcManagerTypesKeyValue(reqMsg *base.Empty) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcManagerTypesKeyValue", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcManagerTypesKeyValue_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcManagerTypesKeyValue", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcEditManagerTypes(reqMsg *share_message.ManagerTypes) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditManagerTypes", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditManagerTypes_(reqMsg *share_message.ManagerTypes) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditManagerTypes", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryDataOverview(reqMsg *base.Empty) *DataOverview {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDataOverview", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DataOverview)
}

func (self *Brower2Backstage) RpcQueryDataOverview_(reqMsg *base.Empty) (*DataOverview, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDataOverview", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DataOverview), e
}
func (self *Brower2Backstage) RpcRegisterLoginReportLine(reqMsg *ListRequest) *LineChartResponse {
	msg, e := self.Sender.CallRpcMethod("RpcRegisterLoginReportLine", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LineChartResponse)
}

func (self *Brower2Backstage) RpcRegisterLoginReportLine_(reqMsg *ListRequest) (*LineChartResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRegisterLoginReportLine", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LineChartResponse), e
}
func (self *Brower2Backstage) RpcInterestTagLine(reqMsg *ListRequest) *LineChartResponse {
	msg, e := self.Sender.CallRpcMethod("RpcInterestTagLine", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LineChartResponse)
}

func (self *Brower2Backstage) RpcInterestTagLine_(reqMsg *ListRequest) (*LineChartResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcInterestTagLine", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LineChartResponse), e
}
func (self *Brower2Backstage) RpcPhoneBrandLine(reqMsg *base.Empty) *NameValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcPhoneBrandLine", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NameValueResponseTag)
}

func (self *Brower2Backstage) RpcPhoneBrandLine_(reqMsg *base.Empty) (*NameValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPhoneBrandLine", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NameValueResponseTag), e
}
func (self *Brower2Backstage) RpcPlayerOnlineLine(reqMsg *ListRequest) *LineChartResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerOnlineLine", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LineChartResponse)
}

func (self *Brower2Backstage) RpcPlayerOnlineLine_(reqMsg *ListRequest) (*LineChartResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerOnlineLine", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LineChartResponse), e
}
func (self *Brower2Backstage) RpcPlayerLogLocation(reqMsg *ListRequest) *NameValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerLogLocation", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NameValueResponseTag)
}

func (self *Brower2Backstage) RpcPlayerLogLocation_(reqMsg *ListRequest) (*NameValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerLogLocation", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NameValueResponseTag), e
}
func (self *Brower2Backstage) RpcPlayerPortrait(reqMsg *base.Empty) *PlayerPortraitResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerPortrait", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerPortraitResponse)
}

func (self *Brower2Backstage) RpcPlayerPortrait_(reqMsg *base.Empty) (*PlayerPortraitResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerPortrait", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerPortraitResponse), e
}
func (self *Brower2Backstage) RpcQueryRolePower(reqMsg *ListRequest) *QueryRolePowerList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryRolePower", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryRolePowerList)
}

func (self *Brower2Backstage) RpcQueryRolePower_(reqMsg *ListRequest) (*QueryRolePowerList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryRolePower", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryRolePowerList), e
}
func (self *Brower2Backstage) RpcUpdateRolePower(reqMsg *share_message.RolePower) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateRolePower", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateRolePower_(reqMsg *share_message.RolePower) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateRolePower", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteRolePower(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteRolePower", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteRolePower_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteRolePower", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetPowerRouter(reqMsg *QueryDataById) *share_message.RolePower {
	msg, e := self.Sender.CallRpcMethod("RpcGetPowerRouter", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.RolePower)
}

func (self *Brower2Backstage) RpcGetPowerRouter_(reqMsg *QueryDataById) (*share_message.RolePower, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPowerRouter", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.RolePower), e
}
func (self *Brower2Backstage) RpcGetRolePowerList(reqMsg *base.Empty) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetRolePowerList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetRolePowerList_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetRolePowerList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcPlayerList(reqMsg *GetPlayerListRequest) *GetPlayerListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetPlayerListResponse)
}

func (self *Brower2Backstage) RpcPlayerList_(reqMsg *GetPlayerListRequest) (*GetPlayerListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetPlayerListResponse), e
}
func (self *Brower2Backstage) RpcGetPlayerById(reqMsg *QueryDataById) *share_message.PlayerBase {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerById", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.PlayerBase)
}

func (self *Brower2Backstage) RpcGetPlayerById_(reqMsg *QueryDataById) (*share_message.PlayerBase, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerById", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.PlayerBase), e
}
func (self *Brower2Backstage) RpcGetPlayerByAccount(reqMsg *QueryDataById) *share_message.PlayerBase {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerByAccount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.PlayerBase)
}

func (self *Brower2Backstage) RpcGetPlayerByAccount_(reqMsg *QueryDataById) (*share_message.PlayerBase, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerByAccount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.PlayerBase), e
}
func (self *Brower2Backstage) RpcEditPlayer(reqMsg *share_message.PlayerBase) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPlayer_(reqMsg *share_message.PlayerBase) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAddPlayer(reqMsg *SigninRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddPlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddPlayer_(reqMsg *SigninRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddPlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAddWaiter(reqMsg *AddWaiterRequest) *AddWaiterResponse {
	msg, e := self.Sender.CallRpcMethod("RpcAddWaiter", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AddWaiterResponse)
}

func (self *Brower2Backstage) RpcAddWaiter_(reqMsg *AddWaiterRequest) (*AddWaiterResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddWaiter", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AddWaiterResponse), e
}
func (self *Brower2Backstage) RpcPlayerFreeze(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerFreeze", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcPlayerFreeze_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerFreeze", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcPlayerUnFreeze(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerUnFreeze", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcPlayerUnFreeze_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerUnFreeze", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryPlayerComplaint(reqMsg *ListRequest) *PlayerComplaintResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPlayerComplaint", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerComplaintResponse)
}

func (self *Brower2Backstage) RpcQueryPlayerComplaint_(reqMsg *ListRequest) (*PlayerComplaintResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPlayerComplaint", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerComplaintResponse), e
}
func (self *Brower2Backstage) RpcQueryPlayerComplaintOther(reqMsg *ListRequest) *PlayerComplaintResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPlayerComplaintOther", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerComplaintResponse)
}

func (self *Brower2Backstage) RpcQueryPlayerComplaintOther_(reqMsg *ListRequest) (*PlayerComplaintResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPlayerComplaintOther", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerComplaintResponse), e
}
func (self *Brower2Backstage) RpcReplyPlayerComplaint(reqMsg *share_message.PlayerComplaint) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReplyPlayerComplaint", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcReplyPlayerComplaint_(reqMsg *share_message.PlayerComplaint) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReplyPlayerComplaint", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcEditPlayerCustomTag(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlayerCustomTag", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPlayerCustomTag_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlayerCustomTag", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcEditPlayerLable(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlayerLable", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPlayerLable_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlayerLable", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcEditPersonalityTags(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPersonalityTags", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPersonalityTags_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPersonalityTags", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetPersonalityTags(reqMsg *base.Empty) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetPersonalityTags", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetPersonalityTags_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPersonalityTags", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcInterestTypeList(reqMsg *ListRequest) *InterestTypeResponse {
	msg, e := self.Sender.CallRpcMethod("RpcInterestTypeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*InterestTypeResponse)
}

func (self *Brower2Backstage) RpcInterestTypeList_(reqMsg *ListRequest) (*InterestTypeResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcInterestTypeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*InterestTypeResponse), e
}
func (self *Brower2Backstage) RpcEditInterestType(reqMsg *share_message.InterestType) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditInterestType", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditInterestType_(reqMsg *share_message.InterestType) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditInterestType", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetInterestTypeList(reqMsg *base.Empty) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetInterestTypeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetInterestTypeList_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetInterestTypeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcInterestTagList(reqMsg *ListRequest) *InterestTagResponse {
	msg, e := self.Sender.CallRpcMethod("RpcInterestTagList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*InterestTagResponse)
}

func (self *Brower2Backstage) RpcInterestTagList_(reqMsg *ListRequest) (*InterestTagResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcInterestTagList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*InterestTagResponse), e
}
func (self *Brower2Backstage) RpcEditInterestTag(reqMsg *share_message.InterestTag) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditInterestTag", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditInterestTag_(reqMsg *share_message.InterestTag) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditInterestTag", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcInterestGroupList(reqMsg *ListRequest) *InterestGroupResponse {
	msg, e := self.Sender.CallRpcMethod("RpcInterestGroupList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*InterestGroupResponse)
}

func (self *Brower2Backstage) RpcInterestGroupList_(reqMsg *ListRequest) (*InterestGroupResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcInterestGroupList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*InterestGroupResponse), e
}
func (self *Brower2Backstage) RpcEditInterestGroup(reqMsg *share_message.InterestGroup) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditInterestGroup", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditInterestGroup_(reqMsg *share_message.InterestGroup) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditInterestGroup", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelInterestGroups(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelInterestGroups", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelInterestGroups_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelInterestGroups", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetInterestTagList(reqMsg *ListRequest) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetInterestTagList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetInterestTagList_(reqMsg *ListRequest) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetInterestTagList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcCustomTagList(reqMsg *ListRequest) *CustomTagResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCustomTagList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CustomTagResponse)
}

func (self *Brower2Backstage) RpcCustomTagList_(reqMsg *ListRequest) (*CustomTagResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCustomTagList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CustomTagResponse), e
}
func (self *Brower2Backstage) RpcEditCustomTag(reqMsg *share_message.CustomTag) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditCustomTag", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditCustomTag_(reqMsg *share_message.CustomTag) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditCustomTag", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetCustomTagList(reqMsg *QueryDataById) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetCustomTagList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetCustomTagList_(reqMsg *QueryDataById) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCustomTagList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcToPlayerCustomTag(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcToPlayerCustomTag", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcToPlayerCustomTag_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToPlayerCustomTag", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGrabTagList(reqMsg *ListRequest) *GrabTagResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGrabTagList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GrabTagResponse)
}

func (self *Brower2Backstage) RpcGrabTagList_(reqMsg *ListRequest) (*GrabTagResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGrabTagList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GrabTagResponse), e
}
func (self *Brower2Backstage) RpcEditGrabTag(reqMsg *share_message.GrabTag) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditGrabTag", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditGrabTag_(reqMsg *share_message.GrabTag) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditGrabTag", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcCrawlWordsList(reqMsg *ListRequest) *CrawlWordsResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCrawlWordsList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CrawlWordsResponse)
}

func (self *Brower2Backstage) RpcCrawlWordsList_(reqMsg *ListRequest) (*CrawlWordsResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCrawlWordsList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CrawlWordsResponse), e
}
func (self *Brower2Backstage) RpcEditCrawlWords(reqMsg *share_message.CrawlWords) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditCrawlWords", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditCrawlWords_(reqMsg *share_message.CrawlWords) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditCrawlWords", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetGrabTagList(reqMsg *base.Empty) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetGrabTagList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetGrabTagList_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetGrabTagList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcDelCrawlWords(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelCrawlWords", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelCrawlWords_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelCrawlWords", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryPlayerWordsList(reqMsg *QueryDataById) *PlayerCrawlWordsResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPlayerWordsList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerCrawlWordsResponse)
}

func (self *Brower2Backstage) RpcQueryPlayerWordsList_(reqMsg *QueryDataById) (*PlayerCrawlWordsResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPlayerWordsList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerCrawlWordsResponse), e
}
func (self *Brower2Backstage) RpcQueryFriendPlayerList(reqMsg *GetPlayerFriendListRequest) *GetPlayerListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryFriendPlayerList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetPlayerListResponse)
}

func (self *Brower2Backstage) RpcQueryFriendPlayerList_(reqMsg *GetPlayerFriendListRequest) (*GetPlayerListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryFriendPlayerList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetPlayerListResponse), e
}
func (self *Brower2Backstage) RpcQueryPlayerInfo(reqMsg *QueryDataById) *PlayerFriendInfo {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPlayerInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerFriendInfo)
}

func (self *Brower2Backstage) RpcQueryPlayerInfo_(reqMsg *QueryDataById) (*PlayerFriendInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPlayerInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerFriendInfo), e
}
func (self *Brower2Backstage) RpcAddFriend(reqMsg *AddPlayerFriendInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddFriend", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddFriend_(reqMsg *AddPlayerFriendInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddFriend", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcPlayerCancleAccountList(reqMsg *ListRequest) *PlayerCancleAccountListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerCancleAccountList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerCancleAccountListResponse)
}

func (self *Brower2Backstage) RpcPlayerCancleAccountList_(reqMsg *ListRequest) (*PlayerCancleAccountListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerCancleAccountList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerCancleAccountListResponse), e
}
func (self *Brower2Backstage) RpcEditPlayerCancleAccount(reqMsg *share_message.PlayerCancleAccount) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlayerCancleAccount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPlayerCancleAccount_(reqMsg *share_message.PlayerCancleAccount) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlayerCancleAccount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryPersonalChatLog(reqMsg *ChatLogRequest) *PersonalChatLogResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPersonalChatLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PersonalChatLogResponse)
}

func (self *Brower2Backstage) RpcQueryPersonalChatLog_(reqMsg *ChatLogRequest) (*PersonalChatLogResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPersonalChatLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PersonalChatLogResponse), e
}
func (self *Brower2Backstage) RpcQueryPersonalChatLogByObj(reqMsg *ChatLogRequest) *PersonalChatLogResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPersonalChatLogByObj", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PersonalChatLogResponse)
}

func (self *Brower2Backstage) RpcQueryPersonalChatLogByObj_(reqMsg *ChatLogRequest) (*PersonalChatLogResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPersonalChatLogByObj", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PersonalChatLogResponse), e
}
func (self *Brower2Backstage) RpcQueryTeamChatLog(reqMsg *ChatLogRequest) *TeamChatLogResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTeamChatLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TeamChatLogResponse)
}

func (self *Brower2Backstage) RpcQueryTeamChatLog_(reqMsg *ChatLogRequest) (*TeamChatLogResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTeamChatLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TeamChatLogResponse), e
}
func (self *Brower2Backstage) RpcCheckChatLogWhitelist(reqMsg *ChatLogRequest) *CommonResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCheckChatLogWhitelist", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CommonResponse)
}

func (self *Brower2Backstage) RpcCheckChatLogWhitelist_(reqMsg *ChatLogRequest) (*CommonResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckChatLogWhitelist", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CommonResponse), e
}
func (self *Brower2Backstage) RpcQueryTeamList(reqMsg *GetTeamListRequest) *GetTeamListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTeamList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetTeamListResponse)
}

func (self *Brower2Backstage) RpcQueryTeamList_(reqMsg *GetTeamListRequest) (*GetTeamListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTeamList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetTeamListResponse), e
}
func (self *Brower2Backstage) RpcGetTeamById(reqMsg *QueryDataById) *share_message.TeamData {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamById", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.TeamData)
}

func (self *Brower2Backstage) RpcGetTeamById_(reqMsg *QueryDataById) (*share_message.TeamData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTeamById", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.TeamData), e
}
func (self *Brower2Backstage) RpcEditTeam(reqMsg *share_message.TeamData) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditTeam", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditTeam_(reqMsg *share_message.TeamData) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditTeam", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDefunctTeam(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDefunctTeam", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDefunctTeam_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDefunctTeam", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcTeamMemberOpt(reqMsg *MemberOptRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcTeamMemberOpt", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcTeamMemberOpt_(reqMsg *MemberOptRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTeamMemberOpt", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryTeamMember(reqMsg *TeamMemberRequest) *TeamMemberResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTeamMember", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TeamMemberResponse)
}

func (self *Brower2Backstage) RpcQueryTeamMember_(reqMsg *TeamMemberRequest) (*TeamMemberResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTeamMember", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TeamMemberResponse), e
}
func (self *Brower2Backstage) RpcExportChatRecord(reqMsg *ListRequest) *ExportChatRecordResponse {
	msg, e := self.Sender.CallRpcMethod("RpcExportChatRecord", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ExportChatRecordResponse)
}

func (self *Brower2Backstage) RpcExportChatRecord_(reqMsg *ListRequest) (*ExportChatRecordResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcExportChatRecord", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ExportChatRecordResponse), e
}
func (self *Brower2Backstage) RpcQueryTeamMessage(reqMsg *ListRequest) *ExportChatRecordResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTeamMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ExportChatRecordResponse)
}

func (self *Brower2Backstage) RpcQueryTeamMessage_(reqMsg *ListRequest) (*ExportChatRecordResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTeamMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ExportChatRecordResponse), e
}
func (self *Brower2Backstage) RpcCreateTeamMessage(reqMsg *CreateTeamInfo) *share_message.TeamData {
	msg, e := self.Sender.CallRpcMethod("RpcCreateTeamMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.TeamData)
}

func (self *Brower2Backstage) RpcCreateTeamMessage_(reqMsg *CreateTeamInfo) (*share_message.TeamData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCreateTeamMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.TeamData), e
}
func (self *Brower2Backstage) RpcQueryTeamPlayerList(reqMsg *GetTeamPlayerListRequest) *GetTeamPlayerListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTeamPlayerList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetTeamPlayerListResponse)
}

func (self *Brower2Backstage) RpcQueryTeamPlayerList_(reqMsg *GetTeamPlayerListRequest) (*GetTeamPlayerListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTeamPlayerList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetTeamPlayerListResponse), e
}
func (self *Brower2Backstage) RpcTeamBan(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcTeamBan", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcTeamBan_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTeamBan", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcTeamUnBan(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcTeamUnBan", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcTeamUnBan_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTeamUnBan", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcTeamCloseAndOpen(reqMsg *TeamManager) *ErrMessage {
	msg, e := self.Sender.CallRpcMethod("RpcTeamCloseAndOpen", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ErrMessage)
}

func (self *Brower2Backstage) RpcTeamCloseAndOpen_(reqMsg *TeamManager) (*ErrMessage, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTeamCloseAndOpen", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ErrMessage), e
}
func (self *Brower2Backstage) RpcTeamMemCloseAndOpen(reqMsg *TeamManager) *ErrMessage {
	msg, e := self.Sender.CallRpcMethod("RpcTeamMemCloseAndOpen", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ErrMessage)
}

func (self *Brower2Backstage) RpcTeamMemCloseAndOpen_(reqMsg *TeamManager) (*ErrMessage, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTeamMemCloseAndOpen", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ErrMessage), e
}
func (self *Brower2Backstage) RpcWarnLord(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcWarnLord", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcWarnLord_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWarnLord", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQuerySouceType(reqMsg *SourceTypeRequest) *SourceTypeResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQuerySouceType", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SourceTypeResponse)
}

func (self *Brower2Backstage) RpcQuerySouceType_(reqMsg *SourceTypeRequest) (*SourceTypeResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQuerySouceType", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SourceTypeResponse), e
}
func (self *Brower2Backstage) RpcQueryGeneralQuota(reqMsg *base.Empty) *share_message.GeneralQuota {
	msg, e := self.Sender.CallRpcMethod("RpcQueryGeneralQuota", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.GeneralQuota)
}

func (self *Brower2Backstage) RpcQueryGeneralQuota_(reqMsg *base.Empty) (*share_message.GeneralQuota, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryGeneralQuota", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.GeneralQuota), e
}
func (self *Brower2Backstage) RpcEditGeneralQuota(reqMsg *share_message.GeneralQuota) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditGeneralQuota", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditGeneralQuota_(reqMsg *share_message.GeneralQuota) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditGeneralQuota", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryPaymentSetting(reqMsg *ListRequest) *PaymentSettingResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPaymentSetting", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PaymentSettingResponse)
}

func (self *Brower2Backstage) RpcQueryPaymentSetting_(reqMsg *ListRequest) (*PaymentSettingResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPaymentSetting", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PaymentSettingResponse), e
}
func (self *Brower2Backstage) RpcEditPaymentSetting(reqMsg *share_message.PaymentSetting) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPaymentSetting", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPaymentSetting_(reqMsg *share_message.PaymentSetting) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPaymentSetting", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelPaymentSetting(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelPaymentSetting", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelPaymentSetting_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelPaymentSetting", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryPayType(reqMsg *QueryDataById) *PayTypeResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPayType", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PayTypeResponse)
}

func (self *Brower2Backstage) RpcQueryPayType_(reqMsg *QueryDataById) (*PayTypeResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPayType", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PayTypeResponse), e
}
func (self *Brower2Backstage) RpcEditPayType(reqMsg *share_message.PayType) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPayType", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPayType_(reqMsg *share_message.PayType) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPayType", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelPayType(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelPayType", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelPayType_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelPayType", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryPayScene(reqMsg *base.Empty) *PaySceneResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPayScene", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PaySceneResponse)
}

func (self *Brower2Backstage) RpcQueryPayScene_(reqMsg *base.Empty) (*PaySceneResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPayScene", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PaySceneResponse), e
}
func (self *Brower2Backstage) RpcEditPayScene(reqMsg *share_message.PayScene) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPayScene", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPayScene_(reqMsg *share_message.PayScene) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPayScene", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelPayScene(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelPayScene", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelPayScene_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelPayScene", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryPaymentPlatform(reqMsg *PlatformChannelRequest) *PaymentPlatformResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPaymentPlatform", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PaymentPlatformResponse)
}

func (self *Brower2Backstage) RpcQueryPaymentPlatform_(reqMsg *PlatformChannelRequest) (*PaymentPlatformResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPaymentPlatform", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PaymentPlatformResponse), e
}
func (self *Brower2Backstage) RpcEditPaymentPlatform(reqMsg *share_message.PaymentPlatform) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPaymentPlatform", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPaymentPlatform_(reqMsg *share_message.PaymentPlatform) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPaymentPlatform", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelPaymentPlatform(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelPaymentPlatform", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelPaymentPlatform_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelPaymentPlatform", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryPlatformChannel(reqMsg *PlatformChannelRequest) *PlatformChannelResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPlatformChannel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlatformChannelResponse)
}

func (self *Brower2Backstage) RpcQueryPlatformChannel_(reqMsg *PlatformChannelRequest) (*PlatformChannelResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPlatformChannel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlatformChannelResponse), e
}
func (self *Brower2Backstage) RpcEditPlatformChannel(reqMsg *share_message.PlatformChannel) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlatformChannel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPlatformChannel_(reqMsg *share_message.PlatformChannel) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlatformChannel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelPlatformChannel(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelPlatformChannel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelPlatformChannel_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelPlatformChannel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcBatchClosePlatformChannel(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcBatchClosePlatformChannel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcBatchClosePlatformChannel_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBatchClosePlatformChannel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAddGold(reqMsg *AddGoldResult) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddGold", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddGold_(reqMsg *AddGoldResult) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddGold", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryGoldLog(reqMsg *QueryGoldLogRequest) *QueryGoldLogResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryGoldLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryGoldLogResponse)
}

func (self *Brower2Backstage) RpcQueryGoldLog_(reqMsg *QueryGoldLogRequest) (*QueryGoldLogResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryGoldLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryGoldLogResponse), e
}
func (self *Brower2Backstage) RpcQueryOrderList(reqMsg *QueryOrderRequest) *QueryOrderResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryOrderList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryOrderResponse)
}

func (self *Brower2Backstage) RpcQueryOrderList_(reqMsg *QueryOrderRequest) (*QueryOrderResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryOrderList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryOrderResponse), e
}
func (self *Brower2Backstage) RpcOptOrder(reqMsg *OptOrderRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcOptOrder", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcOptOrder_(reqMsg *OptOrderRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOptOrder", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcUpdateOrderList(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateOrderList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateOrderList_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateOrderList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcCheckOrder(reqMsg *OptOrderRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCheckOrder", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcCheckOrder_(reqMsg *OptOrderRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckOrder", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcMakePlayerKeepReport(reqMsg *ListRequest) *PlayerKeepReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcMakePlayerKeepReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerKeepReportResponse)
}

func (self *Brower2Backstage) RpcMakePlayerKeepReport_(reqMsg *ListRequest) (*PlayerKeepReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMakePlayerKeepReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerKeepReportResponse), e
}
func (self *Brower2Backstage) RpcPlayerKeepReport(reqMsg *ListRequest) *PlayerKeepReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerKeepReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerKeepReportResponse)
}

func (self *Brower2Backstage) RpcPlayerKeepReport_(reqMsg *ListRequest) (*PlayerKeepReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerKeepReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerKeepReportResponse), e
}
func (self *Brower2Backstage) RpcPlayerActiveReport(reqMsg *ListRequest) *PlayerActiveReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerActiveReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerActiveReportResponse)
}

func (self *Brower2Backstage) RpcPlayerActiveReport_(reqMsg *ListRequest) (*PlayerActiveReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerActiveReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerActiveReportResponse), e
}
func (self *Brower2Backstage) RpcPlayerWeekActiveReport(reqMsg *ListRequest) *PlayerActiveReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerWeekActiveReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerActiveReportResponse)
}

func (self *Brower2Backstage) RpcPlayerWeekActiveReport_(reqMsg *ListRequest) (*PlayerActiveReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerWeekActiveReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerActiveReportResponse), e
}
func (self *Brower2Backstage) RpcPlayerMonthActiveReport(reqMsg *ListRequest) *PlayerActiveReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerMonthActiveReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerActiveReportResponse)
}

func (self *Brower2Backstage) RpcPlayerMonthActiveReport_(reqMsg *ListRequest) (*PlayerActiveReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerMonthActiveReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerActiveReportResponse), e
}
func (self *Brower2Backstage) RpcPlayerBehaviorReport(reqMsg *ListRequest) *PlayerBehaviorReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerBehaviorReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerBehaviorReportResponse)
}

func (self *Brower2Backstage) RpcPlayerBehaviorReport_(reqMsg *ListRequest) (*PlayerBehaviorReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerBehaviorReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerBehaviorReportResponse), e
}
func (self *Brower2Backstage) RpcInOutCashSumReport(reqMsg *ListRequest) *InOutCashSumReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcInOutCashSumReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*InOutCashSumReportResponse)
}

func (self *Brower2Backstage) RpcInOutCashSumReport_(reqMsg *ListRequest) (*InOutCashSumReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcInOutCashSumReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*InOutCashSumReportResponse), e
}
func (self *Brower2Backstage) RpcRegisterLoginReport(reqMsg *ListRequest) *RegisterLoginReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcRegisterLoginReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*RegisterLoginReportResponse)
}

func (self *Brower2Backstage) RpcRegisterLoginReport_(reqMsg *ListRequest) (*RegisterLoginReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRegisterLoginReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*RegisterLoginReportResponse), e
}
func (self *Brower2Backstage) RpcOperationChannelReport(reqMsg *ListRequest) *OperationChannelReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcOperationChannelReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*OperationChannelReportResponse)
}

func (self *Brower2Backstage) RpcOperationChannelReport_(reqMsg *ListRequest) (*OperationChannelReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOperationChannelReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*OperationChannelReportResponse), e
}
func (self *Brower2Backstage) RpcChannelReport(reqMsg *ListRequest) *ChannelReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcChannelReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ChannelReportResponse)
}

func (self *Brower2Backstage) RpcChannelReport_(reqMsg *ListRequest) (*ChannelReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcChannelReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ChannelReportResponse), e
}
func (self *Brower2Backstage) RpcOperationChannelLine(reqMsg *ListRequest) *OperationChannelReportLineResponse {
	msg, e := self.Sender.CallRpcMethod("RpcOperationChannelLine", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*OperationChannelReportLineResponse)
}

func (self *Brower2Backstage) RpcOperationChannelLine_(reqMsg *ListRequest) (*OperationChannelReportLineResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOperationChannelLine", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*OperationChannelReportLineResponse), e
}
func (self *Brower2Backstage) RpcArticleReport(reqMsg *ListRequest) *ArticleReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcArticleReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ArticleReportResponse)
}

func (self *Brower2Backstage) RpcArticleReport_(reqMsg *ListRequest) (*ArticleReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcArticleReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ArticleReportResponse), e
}
func (self *Brower2Backstage) RpcNoticeReport(reqMsg *ListRequest) *ArticleReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcNoticeReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ArticleReportResponse)
}

func (self *Brower2Backstage) RpcNoticeReport_(reqMsg *ListRequest) (*ArticleReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNoticeReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ArticleReportResponse), e
}
func (self *Brower2Backstage) RpcSquareReport(reqMsg *ListRequest) *SquareReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcSquareReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SquareReportResponse)
}

func (self *Brower2Backstage) RpcSquareReport_(reqMsg *ListRequest) (*SquareReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSquareReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SquareReportResponse), e
}
func (self *Brower2Backstage) RpcQueryActivityReport(reqMsg *ListRequest) *ActivityReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryActivityReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ActivityReportResponse)
}

func (self *Brower2Backstage) RpcQueryActivityReport_(reqMsg *ListRequest) (*ActivityReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryActivityReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ActivityReportResponse), e
}
func (self *Brower2Backstage) RpcAdvReport(reqMsg *ListRequest) *AdvReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcAdvReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AdvReportResponse)
}

func (self *Brower2Backstage) RpcAdvReport_(reqMsg *ListRequest) (*AdvReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAdvReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AdvReportResponse), e
}
func (self *Brower2Backstage) RpcEditRegisterLoginReport(reqMsg *share_message.RegisterLoginReport) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditRegisterLoginReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditRegisterLoginReport_(reqMsg *share_message.RegisterLoginReport) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditRegisterLoginReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcEditPlayerKeepReport(reqMsg *share_message.PlayerKeepReport) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlayerKeepReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditPlayerKeepReport_(reqMsg *share_message.PlayerKeepReport) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditPlayerKeepReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcEditOperationChannelReport(reqMsg *share_message.OperationChannelReport) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditOperationChannelReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditOperationChannelReport_(reqMsg *share_message.OperationChannelReport) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditOperationChannelReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryRecallReport(reqMsg *ListRequest) *RecallReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryRecallReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*RecallReportResponse)
}

func (self *Brower2Backstage) RpcQueryRecallReport_(reqMsg *ListRequest) (*RecallReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryRecallReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*RecallReportResponse), e
}
func (self *Brower2Backstage) RpcQueryRecallPlayerLog(reqMsg *ListRequest) *RecallplayerLogResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryRecallPlayerLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*RecallplayerLogResponse)
}

func (self *Brower2Backstage) RpcQueryRecallPlayerLog_(reqMsg *ListRequest) (*RecallplayerLogResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryRecallPlayerLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*RecallplayerLogResponse), e
}
func (self *Brower2Backstage) RpcCoinProductReport(reqMsg *ListRequest) *CoinProductReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCoinProductReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CoinProductReportResponse)
}

func (self *Brower2Backstage) RpcCoinProductReport_(reqMsg *ListRequest) (*CoinProductReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCoinProductReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CoinProductReportResponse), e
}
func (self *Brower2Backstage) RpcCoinProductDetailReport(reqMsg *ListRequest) *CoinProductReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCoinProductDetailReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CoinProductReportResponse)
}

func (self *Brower2Backstage) RpcCoinProductDetailReport_(reqMsg *ListRequest) (*CoinProductReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCoinProductDetailReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CoinProductReportResponse), e
}
func (self *Brower2Backstage) RpcNearbyAdvReport(reqMsg *ListRequest) *NearReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcNearbyAdvReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NearReportResponse)
}

func (self *Brower2Backstage) RpcNearbyAdvReport_(reqMsg *ListRequest) (*NearReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNearbyAdvReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NearReportResponse), e
}
func (self *Brower2Backstage) RpcButtonClickReport(reqMsg *ListRequest) *ButtonClickReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcButtonClickReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ButtonClickReportResponse)
}

func (self *Brower2Backstage) RpcButtonClickReport_(reqMsg *ListRequest) (*ButtonClickReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcButtonClickReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ButtonClickReportResponse), e
}
func (self *Brower2Backstage) RpcPageRegLogReport(reqMsg *ListRequest) *PageRegLogReportResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPageRegLogReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PageRegLogReportResponse)
}

func (self *Brower2Backstage) RpcPageRegLogReport_(reqMsg *ListRequest) (*PageRegLogReportResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPageRegLogReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PageRegLogReportResponse), e
}
func (self *Brower2Backstage) RpcQueryAppPushMessage(reqMsg *QueryFeaturesRequest) *QueryFeaturesResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryAppPushMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryFeaturesResponse)
}

func (self *Brower2Backstage) RpcQueryAppPushMessage_(reqMsg *QueryFeaturesRequest) (*QueryFeaturesResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryAppPushMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryFeaturesResponse), e
}
func (self *Brower2Backstage) RpcEditAppPushMessage(reqMsg *share_message.AppPushMessage) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditAppPushMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditAppPushMessage_(reqMsg *share_message.AppPushMessage) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditAppPushMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQuerySystemNoticeMessage(reqMsg *QueryFeaturesRequest) *QuerySystemNoticeResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQuerySystemNoticeMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QuerySystemNoticeResponse)
}

func (self *Brower2Backstage) RpcQuerySystemNoticeMessage_(reqMsg *QueryFeaturesRequest) (*QuerySystemNoticeResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQuerySystemNoticeMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QuerySystemNoticeResponse), e
}
func (self *Brower2Backstage) RpcEditSystemNoticeMessage(reqMsg *share_message.SystemNotice) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditSystemNoticeMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditSystemNoticeMessage_(reqMsg *share_message.SystemNotice) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditSystemNoticeMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelSystemNoticeMessage(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelSystemNoticeMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelSystemNoticeMessage_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelSystemNoticeMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryTweets(reqMsg *QueryArticleOrTweetsRequest) *QueryTweetsResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTweets", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryTweetsResponse)
}

func (self *Brower2Backstage) RpcQueryTweets_(reqMsg *QueryArticleOrTweetsRequest) (*QueryTweetsResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTweets", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryTweetsResponse), e
}
func (self *Brower2Backstage) RpcAddTweets(reqMsg *share_message.Tweets) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddTweets", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddTweets_(reqMsg *share_message.Tweets) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddTweets", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelTweets(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelTweets", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelTweets_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelTweets", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcSendTweets(reqMsg *share_message.Tweets) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSendTweets", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSendTweets_(reqMsg *share_message.Tweets) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSendTweets", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryArticle(reqMsg *QueryArticleOrTweetsRequest) *QueryArticleResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryArticle", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryArticleResponse)
}

func (self *Brower2Backstage) RpcQueryArticle_(reqMsg *QueryArticleOrTweetsRequest) (*QueryArticleResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryArticle", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryArticleResponse), e
}
func (self *Brower2Backstage) RpcEditArticle(reqMsg *share_message.Article) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditArticle", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditArticle_(reqMsg *share_message.Article) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditArticle", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelArticle(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelArticle", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelArticle_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelArticle", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQuerySysParameterById(reqMsg *QueryDataById) *share_message.SysParameter {
	msg, e := self.Sender.CallRpcMethod("RpcQuerySysParameterById", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.SysParameter)
}

func (self *Brower2Backstage) RpcQuerySysParameterById_(reqMsg *QueryDataById) (*share_message.SysParameter, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQuerySysParameterById", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.SysParameter), e
}
func (self *Brower2Backstage) RpcEditSysParameter(reqMsg *share_message.SysParameter) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditSysParameter", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditSysParameter_(reqMsg *share_message.SysParameter) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditSysParameter", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcCheckShieldScore(reqMsg *CheckScoreRequest) *CheckScoreResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCheckShieldScore", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CheckScoreResponse)
}

func (self *Brower2Backstage) RpcCheckShieldScore_(reqMsg *CheckScoreRequest) (*CheckScoreResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCheckShieldScore", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CheckScoreResponse), e
}
func (self *Brower2Backstage) RpcAddRegisterPush(reqMsg *share_message.RegisterPush) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddRegisterPush", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddRegisterPush_(reqMsg *share_message.RegisterPush) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddRegisterPush", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryRegisterPush(reqMsg *QueryArticleOrTweetsRequest) *QueryRegisterPushResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryRegisterPush", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryRegisterPushResponse)
}

func (self *Brower2Backstage) RpcQueryRegisterPush_(reqMsg *QueryArticleOrTweetsRequest) (*QueryRegisterPushResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryRegisterPush", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryRegisterPushResponse), e
}
func (self *Brower2Backstage) RpcQueryArticleComment(reqMsg *ListRequest) *ArticleCommentResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryArticleComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ArticleCommentResponse)
}

func (self *Brower2Backstage) RpcQueryArticleComment_(reqMsg *ListRequest) (*ArticleCommentResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryArticleComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ArticleCommentResponse), e
}
func (self *Brower2Backstage) RpcDelArticleComment(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelArticleComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelArticleComment_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelArticleComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryNearLead(reqMsg *ListRequest) *QueryNearSetResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryNearLead", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryNearSetResponse)
}

func (self *Brower2Backstage) RpcQueryNearLead_(reqMsg *ListRequest) (*QueryNearSetResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryNearLead", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryNearSetResponse), e
}
func (self *Brower2Backstage) RpcSaveNearLead(reqMsg *share_message.NearSet) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveNearLead", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveNearLead_(reqMsg *share_message.NearSet) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveNearLead", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelNearLead(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelNearLead", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelNearLead_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelNearLead", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryNearFastTerm(reqMsg *ListRequest) *QueryNearSetResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryNearFastTerm", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryNearSetResponse)
}

func (self *Brower2Backstage) RpcQueryNearFastTerm_(reqMsg *ListRequest) (*QueryNearSetResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryNearFastTerm", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryNearSetResponse), e
}
func (self *Brower2Backstage) RpcSaveNearFastTerm(reqMsg *share_message.NearSet) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveNearFastTerm", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveNearFastTerm_(reqMsg *share_message.NearSet) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveNearFastTerm", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcOperationChannelList(reqMsg *ListRequest) *OperationListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcOperationChannelList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*OperationListResponse)
}

func (self *Brower2Backstage) RpcOperationChannelList_(reqMsg *ListRequest) (*OperationListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOperationChannelList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*OperationListResponse), e
}
func (self *Brower2Backstage) RpcEditOperationChannel(reqMsg *share_message.OperationChannel) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditOperationChannel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditOperationChannel_(reqMsg *share_message.OperationChannel) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditOperationChannel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpGetChannelList(reqMsg *base.Empty) *KeyValueResponse {
	msg, e := self.Sender.CallRpcMethod("RpGetChannelList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponse)
}

func (self *Brower2Backstage) RpGetChannelList_(reqMsg *base.Empty) (*KeyValueResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpGetChannelList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponse), e
}
func (self *Brower2Backstage) RpcQueryDirtyWords(reqMsg *ListRequest) *DirtyWordsResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDirtyWords", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DirtyWordsResponse)
}

func (self *Brower2Backstage) RpcQueryDirtyWords_(reqMsg *ListRequest) (*DirtyWordsResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDirtyWords", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DirtyWordsResponse), e
}
func (self *Brower2Backstage) RpcDelDirtyWords(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelDirtyWords", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelDirtyWords_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelDirtyWords", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAddDirtyWords(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddDirtyWords", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddDirtyWords_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddDirtyWords", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQuerySignature(reqMsg *ListRequest) *SignatureResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQuerySignature", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SignatureResponse)
}

func (self *Brower2Backstage) RpcQuerySignature_(reqMsg *ListRequest) (*SignatureResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQuerySignature", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SignatureResponse), e
}
func (self *Brower2Backstage) RpcDelSignature(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelSignature", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelSignature_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelSignature", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAddSignature(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddSignature", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddSignature_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddSignature", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryShopItem(reqMsg *QueryShopItemRequest) *QueryShopItemResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryShopItemResponse)
}

func (self *Brower2Backstage) RpcQueryShopItem_(reqMsg *QueryShopItemRequest) (*QueryShopItemResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryShopItemResponse), e
}
func (self *Brower2Backstage) RpcQueryShopItemDetailById(reqMsg *QueryDataById) *QueryShopItemDetailResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopItemDetailById", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryShopItemDetailResponse)
}

func (self *Brower2Backstage) RpcQueryShopItemDetailById_(reqMsg *QueryDataById) (*QueryShopItemDetailResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopItemDetailById", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryShopItemDetailResponse), e
}
func (self *Brower2Backstage) RpcShopSoldOut(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcShopSoldOut", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcShopSoldOut_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopSoldOut", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetShopItemTypeDropDown(reqMsg *base.Empty) *GetShopItemTypeDropDownResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetShopItemTypeDropDown", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetShopItemTypeDropDownResponse)
}

func (self *Brower2Backstage) RpcGetShopItemTypeDropDown_(reqMsg *base.Empty) (*GetShopItemTypeDropDownResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetShopItemTypeDropDown", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetShopItemTypeDropDownResponse), e
}
func (self *Brower2Backstage) RpcReleaseShopItem(reqMsg *ReleaseEditShopItemObject) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReleaseShopItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcReleaseShopItem_(reqMsg *ReleaseEditShopItemObject) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReleaseShopItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetEditShopItemDetailById(reqMsg *QueryDataById) *ReleaseEditShopItemObject {
	msg, e := self.Sender.CallRpcMethod("RpcGetEditShopItemDetailById", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ReleaseEditShopItemObject)
}

func (self *Brower2Backstage) RpcGetEditShopItemDetailById_(reqMsg *QueryDataById) (*ReleaseEditShopItemObject, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetEditShopItemDetailById", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ReleaseEditShopItemObject), e
}
func (self *Brower2Backstage) RpcEditShopItem(reqMsg *ReleaseEditShopItemObject) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditShopItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditShopItem_(reqMsg *ReleaseEditShopItemObject) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditShopItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryShopComment(reqMsg *QueryShopCommentRequest) *QueryShopCommentResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryShopCommentResponse)
}

func (self *Brower2Backstage) RpcQueryShopComment_(reqMsg *QueryShopCommentRequest) (*QueryShopCommentResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryShopCommentResponse), e
}
func (self *Brower2Backstage) RpcEditShopComment(reqMsg *EditShopCommentRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditShopComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditShopComment_(reqMsg *EditShopCommentRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditShopComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteShopComment(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteShopComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteShopComment_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteShopComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryShopOrder(reqMsg *QueryShopOrderRequest) *QueryShopOrderResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopOrder", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryShopOrderResponse)
}

func (self *Brower2Backstage) RpcQueryShopOrder_(reqMsg *QueryShopOrderRequest) (*QueryShopOrderResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopOrder", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryShopOrderResponse), e
}
func (self *Brower2Backstage) RpcGetExpressComDropDown(reqMsg *base.Empty) *GetExpressComDropDownResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetExpressComDropDown", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetExpressComDropDownResponse)
}

func (self *Brower2Backstage) RpcGetExpressComDropDown_(reqMsg *base.Empty) (*GetExpressComDropDownResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetExpressComDropDown", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetExpressComDropDownResponse), e
}
func (self *Brower2Backstage) RpcSendShopOrder(reqMsg *SendShopOrderRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSendShopOrder", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSendShopOrder_(reqMsg *SendShopOrderRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSendShopOrder", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryShopOrderExpress(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopOrderExpress", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcQueryShopOrderExpress_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopOrderExpress", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryShopReceiveAddress(reqMsg *QueryDataById) *QueryShopReceiveAddressResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopReceiveAddress", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryShopReceiveAddressResponse)
}

func (self *Brower2Backstage) RpcQueryShopReceiveAddress_(reqMsg *QueryDataById) (*QueryShopReceiveAddressResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopReceiveAddress", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryShopReceiveAddressResponse), e
}
func (self *Brower2Backstage) RpcQueryShopDeliverAddress(reqMsg *QueryDataById) *QueryShopDeliverAddressResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopDeliverAddress", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryShopDeliverAddressResponse)
}

func (self *Brower2Backstage) RpcQueryShopDeliverAddress_(reqMsg *QueryDataById) (*QueryShopDeliverAddressResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopDeliverAddress", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryShopDeliverAddressResponse), e
}
func (self *Brower2Backstage) RpcImportShopPointCard(reqMsg *ImportShopPointCardRequest) *ImportShopPointCardResponse {
	msg, e := self.Sender.CallRpcMethod("RpcImportShopPointCard", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ImportShopPointCardResponse)
}

func (self *Brower2Backstage) RpcImportShopPointCard_(reqMsg *ImportShopPointCardRequest) (*ImportShopPointCardResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcImportShopPointCard", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ImportShopPointCardResponse), e
}
func (self *Brower2Backstage) RpcQueryShopPointCard(reqMsg *QueryShopPointCardRequest) *QueryShopPointCardResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopPointCard", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryShopPointCardResponse)
}

func (self *Brower2Backstage) RpcQueryShopPointCard_(reqMsg *QueryShopPointCardRequest) (*QueryShopPointCardResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryShopPointCard", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryShopPointCardResponse), e
}
func (self *Brower2Backstage) RpcGetShopPointCardDropDown(reqMsg *GetShopPointCardDropDownRequest) *GetShopPointCardDropDownResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetShopPointCardDropDown", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetShopPointCardDropDownResponse)
}

func (self *Brower2Backstage) RpcGetShopPointCardDropDown_(reqMsg *GetShopPointCardDropDownRequest) (*GetShopPointCardDropDownResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetShopPointCardDropDown", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetShopPointCardDropDownResponse), e
}
func (self *Brower2Backstage) RpcCancelShopOrder(reqMsg *CancelShopOrderRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCancelShopOrder", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcCancelShopOrder_(reqMsg *CancelShopOrderRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCancelShopOrder", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcCancelShopOrderForWaitSend(reqMsg *CancelShopOrderRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCancelShopOrderForWaitSend", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcCancelShopOrderForWaitSend_(reqMsg *CancelShopOrderRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCancelShopOrderForWaitSend", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcCancelShopOrderForWaitReceive(reqMsg *CancelShopOrderRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCancelShopOrderForWaitReceive", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcCancelShopOrderForWaitReceive_(reqMsg *CancelShopOrderRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCancelShopOrderForWaitReceive", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetIMmessageCount(reqMsg *base.Empty) *IMmessageResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetIMmessageCount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*IMmessageResponse)
}

func (self *Brower2Backstage) RpcGetIMmessageCount_(reqMsg *base.Empty) (*IMmessageResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetIMmessageCount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*IMmessageResponse), e
}
func (self *Brower2Backstage) RpcGetWaiterMsg(reqMsg *base.Empty) *IMmessageNopageResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetWaiterMsg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*IMmessageNopageResponse)
}

func (self *Brower2Backstage) RpcGetWaiterMsg_(reqMsg *base.Empty) (*IMmessageNopageResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWaiterMsg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*IMmessageNopageResponse), e
}
func (self *Brower2Backstage) RpcGetWaiterMsgByMid(reqMsg *ListRequest) *share_message.IMmessage {
	msg, e := self.Sender.CallRpcMethod("RpcGetWaiterMsgByMid", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.IMmessage)
}

func (self *Brower2Backstage) RpcGetWaiterMsgByMid_(reqMsg *ListRequest) (*share_message.IMmessage, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWaiterMsgByMid", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.IMmessage), e
}
func (self *Brower2Backstage) RpcWaiterSendMsgToPlayer(reqMsg *share_message.IMmessage) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterSendMsgToPlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcWaiterSendMsgToPlayer_(reqMsg *share_message.IMmessage) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterSendMsgToPlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcWaiterOverMsgToPlayer(reqMsg *share_message.IMmessage) *share_message.IMmessage {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterOverMsgToPlayer", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.IMmessage)
}

func (self *Brower2Backstage) RpcWaiterOverMsgToPlayer_(reqMsg *share_message.IMmessage) (*share_message.IMmessage, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterOverMsgToPlayer", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.IMmessage), e
}
func (self *Brower2Backstage) RpcWaiterPerformanceList(reqMsg *ListRequest) *WaiterPerformanceResponse {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterPerformanceList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WaiterPerformanceResponse)
}

func (self *Brower2Backstage) RpcWaiterPerformanceList_(reqMsg *ListRequest) (*WaiterPerformanceResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterPerformanceList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WaiterPerformanceResponse), e
}
func (self *Brower2Backstage) RpcWaiterPerformance(reqMsg *QueryDataById) *share_message.WaiterPerformance {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterPerformance", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.WaiterPerformance)
}

func (self *Brower2Backstage) RpcWaiterPerformance_(reqMsg *QueryDataById) (*share_message.WaiterPerformance, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterPerformance", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.WaiterPerformance), e
}
func (self *Brower2Backstage) RpcWaiterReception(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterReception", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcWaiterReception_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterReception", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcWaiterRest(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterRest", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcWaiterRest_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterRest", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcWaiterChatLogList(reqMsg *ListRequest) *IMmessageResponse {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterChatLogList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*IMmessageResponse)
}

func (self *Brower2Backstage) RpcWaiterChatLogList_(reqMsg *ListRequest) (*IMmessageResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterChatLogList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*IMmessageResponse), e
}
func (self *Brower2Backstage) RpcWaiterFAQList(reqMsg *ListRequest) *WaiterFAQResponse {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterFAQList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WaiterFAQResponse)
}

func (self *Brower2Backstage) RpcWaiterFAQList_(reqMsg *ListRequest) (*WaiterFAQResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterFAQList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WaiterFAQResponse), e
}
func (self *Brower2Backstage) RpcEditWaiterFAQ(reqMsg *share_message.WaiterFAQ) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditWaiterFAQ", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditWaiterFAQ_(reqMsg *share_message.WaiterFAQ) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditWaiterFAQ", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcWaiterFastReply(reqMsg *ListRequest) *WaiterFastReplyResponse {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterFastReply", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WaiterFastReplyResponse)
}

func (self *Brower2Backstage) RpcWaiterFastReply_(reqMsg *ListRequest) (*WaiterFastReplyResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterFastReply", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WaiterFastReplyResponse), e
}
func (self *Brower2Backstage) RpcWaiterFastReplyNopage(reqMsg *base.Empty) *WaiterFastReplyResponse {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterFastReplyNopage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WaiterFastReplyResponse)
}

func (self *Brower2Backstage) RpcWaiterFastReplyNopage_(reqMsg *base.Empty) (*WaiterFastReplyResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWaiterFastReplyNopage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WaiterFastReplyResponse), e
}
func (self *Brower2Backstage) RpcEditWaiterFastReply(reqMsg *share_message.WaiterFastReply) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditWaiterFastReply", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditWaiterFastReply_(reqMsg *share_message.WaiterFastReply) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditWaiterFastReply", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelWaiterFastReply(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelWaiterFastReply", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelWaiterFastReply_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelWaiterFastReply", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryDynamic(reqMsg *DynamicListRequest) *DynamicListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DynamicListResponse)
}

func (self *Brower2Backstage) RpcQueryDynamic_(reqMsg *DynamicListRequest) (*DynamicListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DynamicListResponse), e
}
func (self *Brower2Backstage) RpcQueryDynamicDetails(reqMsg *QueryDataById) *share_message.DynamicData {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDynamicDetails", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DynamicData)
}

func (self *Brower2Backstage) RpcQueryDynamicDetails_(reqMsg *QueryDataById) (*share_message.DynamicData, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDynamicDetails", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DynamicData), e
}
func (self *Brower2Backstage) RpcQueryCommentDetails(reqMsg *DynamicListRequest) *CommentList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryCommentDetails", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CommentList)
}

func (self *Brower2Backstage) RpcQueryCommentDetails_(reqMsg *DynamicListRequest) (*CommentList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryCommentDetails", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CommentList), e
}
func (self *Brower2Backstage) RpcUpdateDynamic(reqMsg *share_message.DynamicData) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateDynamic_(reqMsg *share_message.DynamicData) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteDynamic(reqMsg *DelDynamicRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteDynamic_(reqMsg *DelDynamicRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteUnDynamic(reqMsg *DelDynamicRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteUnDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteUnDynamic_(reqMsg *DelDynamicRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteUnDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcShieldDynamic(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcShieldDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcShieldDynamic_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShieldDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteCommentDatas(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteCommentDatas", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteCommentDatas_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteCommentDatas", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcReviewDynamic(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReviewDynamic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcReviewDynamic_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReviewDynamic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetTopicByIds(reqMsg *QueryDataByIds) *TopicResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicByIds", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicResponse)
}

func (self *Brower2Backstage) RpcGetTopicByIds_(reqMsg *QueryDataByIds) (*TopicResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetTopicByIds", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicResponse), e
}
func (self *Brower2Backstage) RpcQueryTopicTypeList(reqMsg *ListRequest) *TopicTypeResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTopicTypeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicTypeResponse)
}

func (self *Brower2Backstage) RpcQueryTopicTypeList_(reqMsg *ListRequest) (*TopicTypeResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTopicTypeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicTypeResponse), e
}
func (self *Brower2Backstage) RpcUpdateTopicType(reqMsg *share_message.TopicType) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateTopicType", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateTopicType_(reqMsg *share_message.TopicType) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateTopicType", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryTopicList(reqMsg *ListRequest) *TopicResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTopicList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicResponse)
}

func (self *Brower2Backstage) RpcQueryTopicList_(reqMsg *ListRequest) (*TopicResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTopicList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicResponse), e
}
func (self *Brower2Backstage) RpcUpdateTopic(reqMsg *share_message.Topic) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateTopic", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateTopic_(reqMsg *share_message.Topic) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateTopic", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryTopicApplyList(reqMsg *ListRequest) *TopicApplyListRes {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTopicApplyList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TopicApplyListRes)
}

func (self *Brower2Backstage) RpcQueryTopicApplyList_(reqMsg *ListRequest) (*TopicApplyListRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTopicApplyList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TopicApplyListRes), e
}
func (self *Brower2Backstage) RpcQueryTopicApply(reqMsg *QueryDataById) *QueryTopicApplyRes {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTopicApply", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryTopicApplyRes)
}

func (self *Brower2Backstage) RpcQueryTopicApply_(reqMsg *QueryDataById) (*QueryTopicApplyRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTopicApply", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryTopicApplyRes), e
}
func (self *Brower2Backstage) RpcAuditTopicApply(reqMsg *AuditTopicApplyReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAuditTopicApply", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAuditTopicApply_(reqMsg *AuditTopicApplyReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAuditTopicApply", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcApplyTopicMasterList(reqMsg *ListRequest) *ApplyTopicMasterRes {
	msg, e := self.Sender.CallRpcMethod("RpcApplyTopicMasterList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ApplyTopicMasterRes)
}

func (self *Brower2Backstage) RpcApplyTopicMasterList_(reqMsg *ListRequest) (*ApplyTopicMasterRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcApplyTopicMasterList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ApplyTopicMasterRes), e
}
func (self *Brower2Backstage) RpcApplyTopicMaster(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcApplyTopicMaster", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcApplyTopicMaster_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcApplyTopicMaster", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryAdvList(reqMsg *ListRequest) *AdvListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryAdvList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AdvListResponse)
}

func (self *Brower2Backstage) RpcQueryAdvList_(reqMsg *ListRequest) (*AdvListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryAdvList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AdvListResponse), e
}
func (self *Brower2Backstage) RpcEditAdvData(reqMsg *share_message.AdvSetting) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditAdvData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditAdvData_(reqMsg *share_message.AdvSetting) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditAdvData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAdvOnShelf(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAdvOnShelf", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAdvOnShelf_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAdvOnShelf", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAdvOffShelf(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAdvOffShelf", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAdvOffShelf_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAdvOffShelf", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcUpdateAdvSort(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateAdvSort", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateAdvSort_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateAdvSort", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelAdvData(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelAdvData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelAdvData_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelAdvData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcIndexTipsList(reqMsg *ListRequest) *IndexTipsResponse {
	msg, e := self.Sender.CallRpcMethod("RpcIndexTipsList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*IndexTipsResponse)
}

func (self *Brower2Backstage) RpcIndexTipsList_(reqMsg *ListRequest) (*IndexTipsResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcIndexTipsList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*IndexTipsResponse), e
}
func (self *Brower2Backstage) RpcSaveIndexTips(reqMsg *share_message.IndexTips) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveIndexTips", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveIndexTips_(reqMsg *share_message.IndexTips) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveIndexTips", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelIndexTips(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelIndexTips", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelIndexTips_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelIndexTips", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcPopSuspendList(reqMsg *ListRequest) *PopSuspendResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPopSuspendList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PopSuspendResponse)
}

func (self *Brower2Backstage) RpcPopSuspendList_(reqMsg *ListRequest) (*PopSuspendResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPopSuspendList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PopSuspendResponse), e
}
func (self *Brower2Backstage) RpcSavePopSuspendList(reqMsg *PopSuspendResponse) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSavePopSuspendList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSavePopSuspendList_(reqMsg *PopSuspendResponse) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSavePopSuspendList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryAdvDownList(reqMsg *ListRequest) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcQueryAdvDownList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcQueryAdvDownList_(reqMsg *ListRequest) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryAdvDownList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcCoinItemList(reqMsg *ListRequest) *CoinItemListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCoinItemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CoinItemListResponse)
}

func (self *Brower2Backstage) RpcCoinItemList_(reqMsg *ListRequest) (*CoinItemListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCoinItemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CoinItemListResponse), e
}
func (self *Brower2Backstage) RpcSaveCoinItem(reqMsg *share_message.CoinRecharge) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveCoinItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveCoinItem_(reqMsg *share_message.CoinRecharge) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveCoinItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGiveCoin(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcGiveCoin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcGiveCoin_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGiveCoin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryCoinChangeLog(reqMsg *ListRequest) *CoinChangeLogResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryCoinChangeLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CoinChangeLogResponse)
}

func (self *Brower2Backstage) RpcQueryCoinChangeLog_(reqMsg *ListRequest) (*CoinChangeLogResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryCoinChangeLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CoinChangeLogResponse), e
}
func (self *Brower2Backstage) RpcQueryPropsItemList(reqMsg *ListRequest) *PropsItemResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPropsItemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PropsItemResponse)
}

func (self *Brower2Backstage) RpcQueryPropsItemList_(reqMsg *ListRequest) (*PropsItemResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPropsItemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PropsItemResponse), e
}
func (self *Brower2Backstage) RpcSavePropsItem(reqMsg *share_message.PropsItem) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSavePropsItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSavePropsItem_(reqMsg *share_message.PropsItem) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSavePropsItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryPropsItemByIds(reqMsg *QueryDataByIds) *PropsItemResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPropsItemByIds", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PropsItemResponse)
}

func (self *Brower2Backstage) RpcQueryPropsItemByIds_(reqMsg *QueryDataByIds) (*PropsItemResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryPropsItemByIds", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PropsItemResponse), e
}
func (self *Brower2Backstage) RpcCoinProductList(reqMsg *ListRequest) *CoinProductResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCoinProductList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CoinProductResponse)
}

func (self *Brower2Backstage) RpcCoinProductList_(reqMsg *ListRequest) (*CoinProductResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCoinProductList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CoinProductResponse), e
}
func (self *Brower2Backstage) RpcSaveCoinProduct(reqMsg *share_message.CoinProduct) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveCoinProduct", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveCoinProduct_(reqMsg *share_message.CoinProduct) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveCoinProduct", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcPlayerBagItem(reqMsg *ListRequest) *PlayerBagItemResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerBagItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerBagItemResponse)
}

func (self *Brower2Backstage) RpcPlayerBagItem_(reqMsg *ListRequest) (*PlayerBagItemResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerBagItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerBagItemResponse), e
}
func (self *Brower2Backstage) RpcPlayerGetPropsLogList(reqMsg *ListRequest) *PlayerGetPropsLogResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerGetPropsLogList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayerGetPropsLogResponse)
}

func (self *Brower2Backstage) RpcPlayerGetPropsLogList_(reqMsg *ListRequest) (*PlayerGetPropsLogResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerGetPropsLogList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayerGetPropsLogResponse), e
}
func (self *Brower2Backstage) RpcRecycleProps(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcRecycleProps", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcRecycleProps_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRecycleProps", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcSysGiveProps(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSysGiveProps", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSysGiveProps_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSysGiveProps", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryCoinProductArry(reqMsg *ListRequest) *KeyValueResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryCoinProductArry", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponse)
}

func (self *Brower2Backstage) RpcQueryCoinProductArry_(reqMsg *ListRequest) (*KeyValueResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryCoinProductArry", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponse), e
}
func (self *Brower2Backstage) RpcBagRecycleProps(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcBagRecycleProps", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcBagRecycleProps_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBagRecycleProps", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcToolWishBoxItemList(reqMsg *QueryDataById) *ToolWishBoxItemListRes {
	msg, e := self.Sender.CallRpcMethod("RpcToolWishBoxItemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ToolWishBoxItemListRes)
}

func (self *Brower2Backstage) RpcToolWishBoxItemList_(reqMsg *QueryDataById) (*ToolWishBoxItemListRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolWishBoxItemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ToolWishBoxItemListRes), e
}
func (self *Brower2Backstage) RpcToolSaveWishBoxItem(reqMsg *ToolSaveWishBoxItemReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcToolSaveWishBoxItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcToolSaveWishBoxItem_(reqMsg *ToolSaveWishBoxItemReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolSaveWishBoxItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcToolDelWishBoxItem(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcToolDelWishBoxItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcToolDelWishBoxItem_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolDelWishBoxItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcToolRateList(reqMsg *ToolRateReq) *ToolRateRes {
	msg, e := self.Sender.CallRpcMethod("RpcToolRateList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ToolRateRes)
}

func (self *Brower2Backstage) RpcToolRateList_(reqMsg *ToolRateReq) (*ToolRateRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolRateList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ToolRateRes), e
}
func (self *Brower2Backstage) RpcToolLucky(reqMsg *ToolLuckyReq) *ToolLuckyRes {
	msg, e := self.Sender.CallRpcMethod("RpcToolLucky", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ToolLuckyRes)
}

func (self *Brower2Backstage) RpcToolLucky_(reqMsg *ToolLuckyReq) (*ToolLuckyRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolLucky", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ToolLuckyRes), e
}
func (self *Brower2Backstage) RpcToolOutputData(reqMsg *QueryDataById) *ToolOutputDataRes {
	msg, e := self.Sender.CallRpcMethod("RpcToolOutputData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ToolOutputDataRes)
}

func (self *Brower2Backstage) RpcToolOutputData_(reqMsg *QueryDataById) (*ToolOutputDataRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolOutputData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ToolOutputDataRes), e
}
func (self *Brower2Backstage) RpcToolOutputitemList(reqMsg *QueryDataById) *ToolOutputitemRes {
	msg, e := self.Sender.CallRpcMethod("RpcToolOutputitemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ToolOutputitemRes)
}

func (self *Brower2Backstage) RpcToolOutputitemList_(reqMsg *QueryDataById) (*ToolOutputitemRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolOutputitemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ToolOutputitemRes), e
}
func (self *Brower2Backstage) RpcToolToolPumping(reqMsg *QueryDataById) *ToolPumping {
	msg, e := self.Sender.CallRpcMethod("RpcToolToolPumping", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ToolPumping)
}

func (self *Brower2Backstage) RpcToolToolPumping_(reqMsg *QueryDataById) (*ToolPumping, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolToolPumping", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ToolPumping), e
}
func (self *Brower2Backstage) RpcToolResetWishPool(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcToolResetWishPool", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcToolResetWishPool_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolResetWishPool", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcToolGetWishPool(reqMsg *QueryDataById) *WishPool {
	msg, e := self.Sender.CallRpcMethod("RpcToolGetWishPool", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishPool)
}

func (self *Brower2Backstage) RpcToolGetWishPool_(reqMsg *QueryDataById) (*WishPool, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolGetWishPool", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishPool), e
}
func (self *Brower2Backstage) RpcToolWishBoxList(reqMsg *ListRequest) *WishBoxList {
	msg, e := self.Sender.CallRpcMethod("RpcToolWishBoxList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxList)
}

func (self *Brower2Backstage) RpcToolWishBoxList_(reqMsg *ListRequest) (*WishBoxList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolWishBoxList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxList), e
}
func (self *Brower2Backstage) RpcToolWishBoxSave(reqMsg *share_message.WishBox) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcToolWishBoxSave", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcToolWishBoxSave_(reqMsg *share_message.WishBox) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToolWishBoxSave", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryWishBoxList(reqMsg *WishBoxListRequest) *WishBoxList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxList)
}

func (self *Brower2Backstage) RpcQueryWishBoxList_(reqMsg *WishBoxListRequest) (*WishBoxList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxList), e
}
func (self *Brower2Backstage) RpcUpdateWishBox(reqMsg *WishBox) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishBox_(reqMsg *WishBox) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishBoxDetail(reqMsg *QueryDataById) *WishBox {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishBoxDetail", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBox)
}

func (self *Brower2Backstage) RpcGetWishBoxDetail_(reqMsg *QueryDataById) (*WishBox, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishBoxDetail", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBox), e
}
func (self *Brower2Backstage) RpcQueryWishBoxGoodsItemList(reqMsg *ListRequest) *WishBoxGoodsItemList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxGoodsItemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxGoodsItemList)
}

func (self *Brower2Backstage) RpcQueryWishBoxGoodsItemList_(reqMsg *ListRequest) (*WishBoxGoodsItemList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxGoodsItemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxGoodsItemList), e
}
func (self *Brower2Backstage) RpcQueryWishBoxWinCfgList(reqMsg *ListRequest) *WishBoxWinCfgList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxWinCfgList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxWinCfgList)
}

func (self *Brower2Backstage) RpcQueryWishBoxWinCfgList_(reqMsg *ListRequest) (*WishBoxWinCfgList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxWinCfgList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxWinCfgList), e
}
func (self *Brower2Backstage) RpcGetWishBoxKvs(reqMsg *base.Empty) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishBoxKvs", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetWishBoxKvs_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishBoxKvs", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcWishBoxLottery(reqMsg *WishBoxLotteryReq) *WishBoxLotteryResp {
	msg, e := self.Sender.CallRpcMethod("RpcWishBoxLottery", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxLotteryResp)
}

func (self *Brower2Backstage) RpcWishBoxLottery_(reqMsg *WishBoxLotteryReq) (*WishBoxLotteryResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWishBoxLottery", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxLotteryResp), e
}
func (self *Brower2Backstage) RpcGetGoodsListByBoxId(reqMsg *QueryDataById) *WishBoxGoodsSelectedList {
	msg, e := self.Sender.CallRpcMethod("RpcGetGoodsListByBoxId", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxGoodsSelectedList)
}

func (self *Brower2Backstage) RpcGetGoodsListByBoxId_(reqMsg *QueryDataById) (*WishBoxGoodsSelectedList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetGoodsListByBoxId", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxGoodsSelectedList), e
}
func (self *Brower2Backstage) RpcQueryWishGoodsList(reqMsg *WishBoxGoodsListRequest) *WishBoxGoodsList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxGoodsList)
}

func (self *Brower2Backstage) RpcQueryWishGoodsList_(reqMsg *WishBoxGoodsListRequest) (*WishBoxGoodsList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxGoodsList), e
}
func (self *Brower2Backstage) RpcUpdateWishGoods(reqMsg *WishBoxGoods) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoods", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishGoods_(reqMsg *WishBoxGoods) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoods", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishGoodsDetail(reqMsg *QueryDataById) *WishBoxGoods {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishGoodsDetail", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxGoods)
}

func (self *Brower2Backstage) RpcGetWishGoodsDetail_(reqMsg *QueryDataById) (*WishBoxGoods, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishGoodsDetail", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxGoods), e
}
func (self *Brower2Backstage) RpcQueryWishGoodsBrandList(reqMsg *ListRequest) *WishGoodsBrandList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsBrandList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishGoodsBrandList)
}

func (self *Brower2Backstage) RpcQueryWishGoodsBrandList_(reqMsg *ListRequest) (*WishGoodsBrandList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsBrandList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishGoodsBrandList), e
}
func (self *Brower2Backstage) RpcUpdateWishGoodsBrand(reqMsg *WishGoodsBrand) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoodsBrand", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishGoodsBrand_(reqMsg *WishGoodsBrand) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoodsBrand", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishGoodsBrandKvs(reqMsg *ListRequest) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishGoodsBrandKvs", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetWishGoodsBrandKvs_(reqMsg *ListRequest) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishGoodsBrandKvs", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcQueryWishGoodsTypeList(reqMsg *ListRequest) *WishGoodsTypeList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsTypeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishGoodsTypeList)
}

func (self *Brower2Backstage) RpcQueryWishGoodsTypeList_(reqMsg *ListRequest) (*WishGoodsTypeList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsTypeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishGoodsTypeList), e
}
func (self *Brower2Backstage) RpcUpdateWishGoodsType(reqMsg *WishGoodsType) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoodsType", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishGoodsType_(reqMsg *WishGoodsType) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoodsType", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishGoodsTypeKvs(reqMsg *ListRequest) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishGoodsTypeKvs", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetWishGoodsTypeKvs_(reqMsg *ListRequest) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishGoodsTypeKvs", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcQueryWishDeliveryOrderList(reqMsg *ListRequest) *WishDeliveryOrderList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishDeliveryOrderList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishDeliveryOrderList)
}

func (self *Brower2Backstage) RpcQueryWishDeliveryOrderList_(reqMsg *ListRequest) (*WishDeliveryOrderList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishDeliveryOrderList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishDeliveryOrderList), e
}
func (self *Brower2Backstage) RpcUpdateDeliveryOrderCourierInfo(reqMsg *UpdateDeliveryOrderCourierInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateDeliveryOrderCourierInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateDeliveryOrderCourierInfo_(reqMsg *UpdateDeliveryOrderCourierInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateDeliveryOrderCourierInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcUpdateDeliveryOrderStatus(reqMsg *UpdateStatusRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateDeliveryOrderStatus", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateDeliveryOrderStatus_(reqMsg *UpdateStatusRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateDeliveryOrderStatus", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryWishRecycleOrderList(reqMsg *ListRequest) *WishRecycleOrderList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishRecycleOrderList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishRecycleOrderList)
}

func (self *Brower2Backstage) RpcQueryWishRecycleOrderList_(reqMsg *ListRequest) (*WishRecycleOrderList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishRecycleOrderList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishRecycleOrderList), e
}
func (self *Brower2Backstage) RpcGetWishRecycleOrderDetail(reqMsg *QueryDataById) *WishRecycleOrderDetailList {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishRecycleOrderDetail", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishRecycleOrderDetailList)
}

func (self *Brower2Backstage) RpcGetWishRecycleOrderDetail_(reqMsg *QueryDataById) (*WishRecycleOrderDetailList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishRecycleOrderDetail", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishRecycleOrderDetailList), e
}
func (self *Brower2Backstage) RpcUpdateWishRecycleOrderStatus(reqMsg *UpdateStatusRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishRecycleOrderStatus", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishRecycleOrderStatus_(reqMsg *UpdateStatusRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishRecycleOrderStatus", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishRecycleOrderUserInfo(reqMsg *QueryDataById) *WishRecycleOrderUserInfo {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishRecycleOrderUserInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishRecycleOrderUserInfo)
}

func (self *Brower2Backstage) RpcGetWishRecycleOrderUserInfo_(reqMsg *QueryDataById) (*WishRecycleOrderUserInfo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishRecycleOrderUserInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishRecycleOrderUserInfo), e
}
func (self *Brower2Backstage) RpcQueryWishOrderList(reqMsg *QueryWishOrderRequest) *QueryOrderResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishOrderList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryOrderResponse)
}

func (self *Brower2Backstage) RpcQueryWishOrderList_(reqMsg *QueryWishOrderRequest) (*QueryOrderResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishOrderList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryOrderResponse), e
}
func (self *Brower2Backstage) RpcOptWishOrder(reqMsg *OptOrderRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcOptWishOrder", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcOptWishOrder_(reqMsg *OptOrderRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOptWishOrder", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcWishPlayerList(reqMsg *ListRequest) *WishPlayerListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcWishPlayerList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishPlayerListResponse)
}

func (self *Brower2Backstage) RpcWishPlayerList_(reqMsg *ListRequest) (*WishPlayerListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWishPlayerList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishPlayerListResponse), e
}
func (self *Brower2Backstage) RpcWishPlayerFreezeDiamond(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcWishPlayerFreezeDiamond", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcWishPlayerFreezeDiamond_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWishPlayerFreezeDiamond", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcWishPlayerUnFreezeDiamond(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcWishPlayerUnFreezeDiamond", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcWishPlayerUnFreezeDiamond_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWishPlayerUnFreezeDiamond", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryWishPoolReportList(reqMsg *ListRequest) *WishPoolReportList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishPoolReportList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishPoolReportList)
}

func (self *Brower2Backstage) RpcQueryWishPoolReportList_(reqMsg *ListRequest) (*WishPoolReportList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishPoolReportList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishPoolReportList), e
}
func (self *Brower2Backstage) RpcQueryWishBoxReportList(reqMsg *ListRequest) *WishBoxReportList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxReportList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxReportList)
}

func (self *Brower2Backstage) RpcQueryWishBoxReportList_(reqMsg *ListRequest) (*WishBoxReportList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxReportList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxReportList), e
}
func (self *Brower2Backstage) RpcQueryWishBoxDetailReportList(reqMsg *ListRequest) *WishBoxDetailReportList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxDetailReportList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxDetailReportList)
}

func (self *Brower2Backstage) RpcQueryWishBoxDetailReportList_(reqMsg *ListRequest) (*WishBoxDetailReportList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxDetailReportList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxDetailReportList), e
}
func (self *Brower2Backstage) RpcQueryWishItemReportList(reqMsg *ListRequest) *WishItemReportList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishItemReportList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishItemReportList)
}

func (self *Brower2Backstage) RpcQueryWishItemReportList_(reqMsg *ListRequest) (*WishItemReportList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishItemReportList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishItemReportList), e
}
func (self *Brower2Backstage) RpcQueryTestPlayerWishItemList(reqMsg *ListRequest) *TestPlayerWishItemList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTestPlayerWishItemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TestPlayerWishItemList)
}

func (self *Brower2Backstage) RpcQueryTestPlayerWishItemList_(reqMsg *ListRequest) (*TestPlayerWishItemList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTestPlayerWishItemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TestPlayerWishItemList), e
}
func (self *Brower2Backstage) RpcQueryTestWishPoolLogList(reqMsg *ListRequest) *TestWishPoolLogList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTestWishPoolLogList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TestWishPoolLogList)
}

func (self *Brower2Backstage) RpcQueryTestWishPoolLogList_(reqMsg *ListRequest) (*TestWishPoolLogList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTestWishPoolLogList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TestWishPoolLogList), e
}
func (self *Brower2Backstage) RpcQueryTestWishPoolPumpLogList(reqMsg *ListRequest) *TestWishPoolPumpLogList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTestWishPoolPumpLogList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TestWishPoolPumpLogList)
}

func (self *Brower2Backstage) RpcQueryTestWishPoolPumpLogList_(reqMsg *ListRequest) (*TestWishPoolPumpLogList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTestWishPoolPumpLogList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TestWishPoolPumpLogList), e
}
func (self *Brower2Backstage) RpcQueryTestWishPoolBoxPoolInfoList(reqMsg *ListRequest) *TestWishPoolBoxPoolInfoList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTestWishPoolBoxPoolInfoList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TestWishPoolBoxPoolInfoList)
}

func (self *Brower2Backstage) RpcQueryTestWishPoolBoxPoolInfoList_(reqMsg *ListRequest) (*TestWishPoolBoxPoolInfoList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryTestWishPoolBoxPoolInfoList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TestWishPoolBoxPoolInfoList), e
}
func (self *Brower2Backstage) RpcQueryDrawRecordList(reqMsg *ListRequest) *DrawRecordList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDrawRecordList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DrawRecordList)
}

func (self *Brower2Backstage) RpcQueryDrawRecordList_(reqMsg *ListRequest) (*DrawRecordList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDrawRecordList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DrawRecordList), e
}
func (self *Brower2Backstage) RpcQueryAddBoxRecordList(reqMsg *ListRequest) *AddBoxRecordList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryAddBoxRecordList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AddBoxRecordList)
}

func (self *Brower2Backstage) RpcQueryAddBoxRecordList_(reqMsg *ListRequest) (*AddBoxRecordList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryAddBoxRecordList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AddBoxRecordList), e
}
func (self *Brower2Backstage) RpcQueryWishGoodsRecordList(reqMsg *ListRequest) *WishGoodsRecordList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsRecordList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishGoodsRecordList)
}

func (self *Brower2Backstage) RpcQueryWishGoodsRecordList_(reqMsg *ListRequest) (*WishGoodsRecordList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsRecordList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishGoodsRecordList), e
}
func (self *Brower2Backstage) RpcQueryDrawBoxRecordList(reqMsg *ListRequest) *DrawBoxRecordList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDrawBoxRecordList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DrawBoxRecordList)
}

func (self *Brower2Backstage) RpcQueryDrawBoxRecordList_(reqMsg *ListRequest) (*DrawBoxRecordList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDrawBoxRecordList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DrawBoxRecordList), e
}
func (self *Brower2Backstage) RpcQueryHaveItemList(reqMsg *ListRequest) *HaveItemList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryHaveItemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*HaveItemList)
}

func (self *Brower2Backstage) RpcQueryHaveItemList_(reqMsg *ListRequest) (*HaveItemList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryHaveItemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*HaveItemList), e
}
func (self *Brower2Backstage) RpcDeleteHaveItem(reqMsg *QueryDataByIds) *HaveItemList {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteHaveItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*HaveItemList)
}

func (self *Brower2Backstage) RpcDeleteHaveItem_(reqMsg *QueryDataByIds) (*HaveItemList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteHaveItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*HaveItemList), e
}
func (self *Brower2Backstage) RpcQueryWinRecordList(reqMsg *ListRequest) *WinRecordList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWinRecordList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WinRecordList)
}

func (self *Brower2Backstage) RpcQueryWinRecordList_(reqMsg *ListRequest) (*WinRecordList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWinRecordList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WinRecordList), e
}
func (self *Brower2Backstage) RpcQueryWishDelItemList(reqMsg *ListRequest) *HaveItemList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishDelItemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*HaveItemList)
}

func (self *Brower2Backstage) RpcQueryWishDelItemList_(reqMsg *ListRequest) (*HaveItemList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishDelItemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*HaveItemList), e
}
func (self *Brower2Backstage) RpcQueryWishPoolList(reqMsg *ListRequest) *WishPoolList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishPoolList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishPoolList)
}

func (self *Brower2Backstage) RpcQueryWishPoolList_(reqMsg *ListRequest) (*WishPoolList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishPoolList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishPoolList), e
}
func (self *Brower2Backstage) RpcUpdateWishPool(reqMsg *WishPool) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishPool", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishPool_(reqMsg *WishPool) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishPool", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteWishPool(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishPool", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteWishPool_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishPool", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishPoolKvs(reqMsg *base.Empty) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishPoolKvs", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetWishPoolKvs_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishPoolKvs", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcUpdateDefaultWish(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateDefaultWish", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateDefaultWish_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateDefaultWish", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishPool(reqMsg *QueryDataById) *WishPool {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishPool", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishPool)
}

func (self *Brower2Backstage) RpcGetWishPool_(reqMsg *QueryDataById) (*WishPool, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishPool", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishPool), e
}
func (self *Brower2Backstage) RpcResetWishPool(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcResetWishPool", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcResetWishPool_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcResetWishPool", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDiamondItemList(reqMsg *ListRequest) *DiamondItemListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcDiamondItemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DiamondItemListResponse)
}

func (self *Brower2Backstage) RpcDiamondItemList_(reqMsg *ListRequest) (*DiamondItemListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDiamondItemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DiamondItemListResponse), e
}
func (self *Brower2Backstage) RpcSaveDiamondItem(reqMsg *share_message.DiamondRecharge) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveDiamondItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveDiamondItem_(reqMsg *share_message.DiamondRecharge) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveDiamondItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGiveDiamond(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcGiveDiamond", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcGiveDiamond_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGiveDiamond", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryDiamondChangeLog(reqMsg *ListRequest) *DiamondChangeLogResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDiamondChangeLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DiamondChangeLogResponse)
}

func (self *Brower2Backstage) RpcQueryDiamondChangeLog_(reqMsg *ListRequest) (*DiamondChangeLogResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryDiamondChangeLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DiamondChangeLogResponse), e
}
func (self *Brower2Backstage) RpcGetPriceSection(reqMsg *base.Empty) *PriceSection {
	msg, e := self.Sender.CallRpcMethod("RpcGetPriceSection", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PriceSection)
}

func (self *Brower2Backstage) RpcGetPriceSection_(reqMsg *base.Empty) (*PriceSection, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPriceSection", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PriceSection), e
}
func (self *Brower2Backstage) RpcUpdatePriceSection(reqMsg *PriceSection) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdatePriceSection", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdatePriceSection_(reqMsg *PriceSection) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdatePriceSection", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetMailSection(reqMsg *base.Empty) *WishMailSection {
	msg, e := self.Sender.CallRpcMethod("RpcGetMailSection", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishMailSection)
}

func (self *Brower2Backstage) RpcGetMailSection_(reqMsg *base.Empty) (*WishMailSection, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetMailSection", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishMailSection), e
}
func (self *Brower2Backstage) RpcUpdateMailSection(reqMsg *WishMailSection) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateMailSection", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateMailSection_(reqMsg *WishMailSection) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateMailSection", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishRecycleSection(reqMsg *base.Empty) *WishRecycleSection {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishRecycleSection", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishRecycleSection)
}

func (self *Brower2Backstage) RpcGetWishRecycleSection_(reqMsg *base.Empty) (*WishRecycleSection, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishRecycleSection", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishRecycleSection), e
}
func (self *Brower2Backstage) RpcUpdateWishRecycleSection(reqMsg *WishRecycleSection) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishRecycleSection", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishRecycleSection_(reqMsg *WishRecycleSection) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishRecycleSection", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishPayWarnCfg(reqMsg *base.Empty) *WishPayWarnCfg {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishPayWarnCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishPayWarnCfg)
}

func (self *Brower2Backstage) RpcGetWishPayWarnCfg_(reqMsg *base.Empty) (*WishPayWarnCfg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishPayWarnCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishPayWarnCfg), e
}
func (self *Brower2Backstage) RpcUpdateWishPayWarnCfg(reqMsg *WishPayWarnCfg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishPayWarnCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishPayWarnCfg_(reqMsg *WishPayWarnCfg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishPayWarnCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishCoolDownConfig(reqMsg *base.Empty) *WishCoolDownConfig {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishCoolDownConfig", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishCoolDownConfig)
}

func (self *Brower2Backstage) RpcGetWishCoolDownConfig_(reqMsg *base.Empty) (*WishCoolDownConfig, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishCoolDownConfig", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishCoolDownConfig), e
}
func (self *Brower2Backstage) RpcUpdateWishCoolDownConfig(reqMsg *WishCoolDownConfig) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishCoolDownConfig", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishCoolDownConfig_(reqMsg *WishCoolDownConfig) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishCoolDownConfig", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishCurrencyConversionCfg(reqMsg *base.Empty) *WishCurrencyConversionCfg {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishCurrencyConversionCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishCurrencyConversionCfg)
}

func (self *Brower2Backstage) RpcGetWishCurrencyConversionCfg_(reqMsg *base.Empty) (*WishCurrencyConversionCfg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishCurrencyConversionCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishCurrencyConversionCfg), e
}
func (self *Brower2Backstage) RpcUpdateWishCurrencyConversionCfg(reqMsg *WishCurrencyConversionCfg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishCurrencyConversionCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishCurrencyConversionCfg_(reqMsg *WishCurrencyConversionCfg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishCurrencyConversionCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetWishGuardianCfg(reqMsg *base.Empty) *WishGuardianCfg {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishGuardianCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishGuardianCfg)
}

func (self *Brower2Backstage) RpcGetWishGuardianCfg_(reqMsg *base.Empty) (*WishGuardianCfg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishGuardianCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishGuardianCfg), e
}
func (self *Brower2Backstage) RpcUpdateWishGuardianCfg(reqMsg *WishGuardianCfg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGuardianCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishGuardianCfg_(reqMsg *WishGuardianCfg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGuardianCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcSaveRecycleNoteCfg(reqMsg *RecycleNoteCfg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveRecycleNoteCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveRecycleNoteCfg_(reqMsg *RecycleNoteCfg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveRecycleNoteCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetRecycleNoteCfg(reqMsg *base.Empty) *RecycleNoteCfg {
	msg, e := self.Sender.CallRpcMethod("RpcGetRecycleNoteCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*RecycleNoteCfg)
}

func (self *Brower2Backstage) RpcGetRecycleNoteCfg_(reqMsg *base.Empty) (*RecycleNoteCfg, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetRecycleNoteCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*RecycleNoteCfg), e
}
func (self *Brower2Backstage) RpcPayPlayerLocation(reqMsg *ListRequest) *NameValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcPayPlayerLocation", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NameValueResponseTag)
}

func (self *Brower2Backstage) RpcPayPlayerLocation_(reqMsg *ListRequest) (*NameValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPayPlayerLocation", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NameValueResponseTag), e
}
func (self *Brower2Backstage) RpcWishCoinRechargeActivityCfgList(reqMsg *ListRequest) *WishCoinRechargeActivityCfgRes {
	msg, e := self.Sender.CallRpcMethod("RpcWishCoinRechargeActivityCfgList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishCoinRechargeActivityCfgRes)
}

func (self *Brower2Backstage) RpcWishCoinRechargeActivityCfgList_(reqMsg *ListRequest) (*WishCoinRechargeActivityCfgRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWishCoinRechargeActivityCfgList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishCoinRechargeActivityCfgRes), e
}
func (self *Brower2Backstage) RpcWishCoinRechargeActivityCfgUpdate(reqMsg *share_message.WishCoinRechargeActivityCfg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcWishCoinRechargeActivityCfgUpdate", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcWishCoinRechargeActivityCfgUpdate_(reqMsg *share_message.WishCoinRechargeActivityCfg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWishCoinRechargeActivityCfgUpdate", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcWishCoinRechargeActivityCfgDel(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcWishCoinRechargeActivityCfgDel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcWishCoinRechargeActivityCfgDel_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWishCoinRechargeActivityCfgDel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryWishActPool(reqMsg *ListRequest) *WishActPoolList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPool", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishActPoolList)
}

func (self *Brower2Backstage) RpcQueryWishActPool_(reqMsg *ListRequest) (*WishActPoolList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPool", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishActPoolList), e
}
func (self *Brower2Backstage) RpcUpdateWishActPool(reqMsg *share_message.WishActPool) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishActPool", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishActPool_(reqMsg *share_message.WishActPool) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishActPool", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteWishActPool(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishActPool", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteWishActPool_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishActPool", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcAddWishAllowList(reqMsg *AddWishAllowListReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddWishAllowList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddWishAllowList_(reqMsg *AddWishAllowListReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddWishAllowList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteWishAllowList(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishAllowList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteWishAllowList_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishAllowList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcWishAllowList(reqMsg *base.Empty) *WishAllowListResp {
	msg, e := self.Sender.CallRpcMethod("RpcWishAllowList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishAllowListResp)
}

func (self *Brower2Backstage) RpcWishAllowList_(reqMsg *base.Empty) (*WishAllowListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWishAllowList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishAllowListResp), e
}
func (self *Brower2Backstage) RpcQueryWishActPoolDetail(reqMsg *ListRequest) *WishActPoolDetail {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPoolDetail", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishActPoolDetail)
}

func (self *Brower2Backstage) RpcQueryWishActPoolDetail_(reqMsg *ListRequest) (*WishActPoolDetail, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPoolDetail", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishActPoolDetail), e
}
func (self *Brower2Backstage) RpcGetWishActPoolTypeKvs(reqMsg *base.Empty) *KeyValueResponseTag {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishActPoolTypeKvs", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*KeyValueResponseTag)
}

func (self *Brower2Backstage) RpcGetWishActPoolTypeKvs_(reqMsg *base.Empty) (*KeyValueResponseTag, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishActPoolTypeKvs", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*KeyValueResponseTag), e
}
func (self *Brower2Backstage) RpcGetWishActCfg(reqMsg *ListRequest) *share_message.Activity {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishActCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.Activity)
}

func (self *Brower2Backstage) RpcGetWishActCfg_(reqMsg *ListRequest) (*share_message.Activity, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishActCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.Activity), e
}
func (self *Brower2Backstage) RpcUpdateWishActCfg(reqMsg *share_message.Activity) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishActCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishActCfg_(reqMsg *share_message.Activity) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishActCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcWishPoolActivityReportList(reqMsg *ListRequest) *WishPoolActivityReportResp {
	msg, e := self.Sender.CallRpcMethod("RpcWishPoolActivityReportList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishPoolActivityReportResp)
}

func (self *Brower2Backstage) RpcWishPoolActivityReportList_(reqMsg *ListRequest) (*WishPoolActivityReportResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWishPoolActivityReportList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishPoolActivityReportResp), e
}
func (self *Brower2Backstage) RpcQueryWishActPoolRuleDay(reqMsg *ListRequest) *WishActPoolRuleList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPoolRuleDay", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishActPoolRuleList)
}

func (self *Brower2Backstage) RpcQueryWishActPoolRuleDay_(reqMsg *ListRequest) (*WishActPoolRuleList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPoolRuleDay", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishActPoolRuleList), e
}
func (self *Brower2Backstage) RpcAddWishActPoolRuleDay(reqMsg *AddWishActPoolRuleRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddWishActPoolRuleDay", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddWishActPoolRuleDay_(reqMsg *AddWishActPoolRuleRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddWishActPoolRuleDay", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcUpdateWishActPoolRuleDay(reqMsg *WishActPoolRule) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishActPoolRuleDay", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishActPoolRuleDay_(reqMsg *WishActPoolRule) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishActPoolRuleDay", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteWishActPoolRuleDay(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishActPoolRuleDay", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteWishActPoolRuleDay_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishActPoolRuleDay", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryWishActPoolRuleCount(reqMsg *ListRequest) *WishActPoolRuleList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPoolRuleCount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishActPoolRuleList)
}

func (self *Brower2Backstage) RpcQueryWishActPoolRuleCount_(reqMsg *ListRequest) (*WishActPoolRuleList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPoolRuleCount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishActPoolRuleList), e
}
func (self *Brower2Backstage) RpcAddWishActPoolRuleCount(reqMsg *AddWishActPoolRuleRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddWishActPoolRuleCount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddWishActPoolRuleCount_(reqMsg *AddWishActPoolRuleRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddWishActPoolRuleCount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcUpdateWishActPoolRuleCount(reqMsg *WishActPoolRule) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishActPoolRuleCount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishActPoolRuleCount_(reqMsg *WishActPoolRule) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishActPoolRuleCount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteWishActPoolRuleCount(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishActPoolRuleCount", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteWishActPoolRuleCount_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishActPoolRuleCount", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryWishActPoolRuleWeekMonth(reqMsg *ListRequest) *WishActPoolRuleList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPoolRuleWeekMonth", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishActPoolRuleList)
}

func (self *Brower2Backstage) RpcQueryWishActPoolRuleWeekMonth_(reqMsg *ListRequest) (*WishActPoolRuleList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPoolRuleWeekMonth", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishActPoolRuleList), e
}
func (self *Brower2Backstage) RpcAddWishActPoolRuleWeekMonth(reqMsg *AddWishActPoolRuleRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddWishActPoolRuleWeekMonth", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddWishActPoolRuleWeekMonth_(reqMsg *AddWishActPoolRuleRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddWishActPoolRuleWeekMonth", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcUpdateWishActPoolRuleItemWeekMonth(reqMsg *WishActPoolAwardItem) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishActPoolRuleItemWeekMonth", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateWishActPoolRuleItemWeekMonth_(reqMsg *WishActPoolAwardItem) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishActPoolRuleItemWeekMonth", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDeleteWishActPoolRuleItemWeekMonth(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishActPoolRuleItemWeekMonth", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDeleteWishActPoolRuleItemWeekMonth_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteWishActPoolRuleItemWeekMonth", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryWishActPlayerRecordList(reqMsg *ListRequest) *WishActPlayerRecordList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPlayerRecordList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishActPlayerRecordList)
}

func (self *Brower2Backstage) RpcQueryWishActPlayerRecordList_(reqMsg *ListRequest) (*WishActPlayerRecordList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPlayerRecordList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishActPlayerRecordList), e
}
func (self *Brower2Backstage) RpcQueryWishActPlayerWinRecordList(reqMsg *ListRequest) *WishActPlayerWinRecordList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPlayerWinRecordList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishActPlayerWinRecordList)
}

func (self *Brower2Backstage) RpcQueryWishActPlayerWinRecordList_(reqMsg *ListRequest) (*WishActPlayerWinRecordList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPlayerWinRecordList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishActPlayerWinRecordList), e
}
func (self *Brower2Backstage) RpcQueryWishActPlayerDrawRecordList(reqMsg *ListRequest) *WishActPlayerDrawRecordList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPlayerDrawRecordList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishActPlayerDrawRecordList)
}

func (self *Brower2Backstage) RpcQueryWishActPlayerDrawRecordList_(reqMsg *ListRequest) (*WishActPlayerDrawRecordList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActPlayerDrawRecordList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishActPlayerDrawRecordList), e
}
func (self *Brower2Backstage) RpcQueryWishLogReport(reqMsg *ListRequest) *QueryWishLogReportRes {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishLogReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryWishLogReportRes)
}

func (self *Brower2Backstage) RpcQueryWishLogReport_(reqMsg *ListRequest) (*QueryWishLogReportRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishLogReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryWishLogReportRes), e
}
func (self *Brower2Backstage) RpcQueryWishActivityPrizeLog(reqMsg *ListRequest) *QueryWishActivityPrizeLogRes {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActivityPrizeLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryWishActivityPrizeLogRes)
}

func (self *Brower2Backstage) RpcQueryWishActivityPrizeLog_(reqMsg *ListRequest) (*QueryWishActivityPrizeLogRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishActivityPrizeLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryWishActivityPrizeLogRes), e
}
func (self *Brower2Backstage) RpcPlayerCardList(reqMsg *ListRequest) *InterestTagResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerCardList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*InterestTagResponse)
}

func (self *Brower2Backstage) RpcPlayerCardList_(reqMsg *ListRequest) (*InterestTagResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerCardList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*InterestTagResponse), e
}
func (self *Brower2Backstage) RpcUpdatePlayerCard(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdatePlayerCard", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdatePlayerCard_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdatePlayerCard", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcCharacterTagList(reqMsg *ListRequest) *InterestTagResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCharacterTagList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*InterestTagResponse)
}

func (self *Brower2Backstage) RpcCharacterTagList_(reqMsg *ListRequest) (*InterestTagResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCharacterTagList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*InterestTagResponse), e
}
func (self *Brower2Backstage) RpcSaveCharacterTag(reqMsg *share_message.InterestTag) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveCharacterTag", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveCharacterTag_(reqMsg *share_message.InterestTag) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveCharacterTag", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcPlayerVoiceWorkList(reqMsg *ListRequest) *VoiceWorkListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerVoiceWorkList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*VoiceWorkListResponse)
}

func (self *Brower2Backstage) RpcPlayerVoiceWorkList_(reqMsg *ListRequest) (*VoiceWorkListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayerVoiceWorkList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*VoiceWorkListResponse), e
}
func (self *Brower2Backstage) RpcReviewePlayerVoiceWork(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReviewePlayerVoiceWork", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcReviewePlayerVoiceWork_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReviewePlayerVoiceWork", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelPlayerVoiceWork(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelPlayerVoiceWork", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelPlayerVoiceWork_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelPlayerVoiceWork", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcUploadPlayerVoiceWork(reqMsg *share_message.PlayerMixVoiceVideo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUploadPlayerVoiceWork", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUploadPlayerVoiceWork_(reqMsg *share_message.PlayerMixVoiceVideo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUploadPlayerVoiceWork", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetPlayerVoiceWorkUse(reqMsg *QueryDataById) *share_message.PlayerMixVoiceVideo {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerVoiceWorkUse", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.PlayerMixVoiceVideo)
}

func (self *Brower2Backstage) RpcGetPlayerVoiceWorkUse_(reqMsg *QueryDataById) (*share_message.PlayerMixVoiceVideo, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerVoiceWorkUse", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.PlayerMixVoiceVideo), e
}
func (self *Brower2Backstage) RpcBgTagList(reqMsg *ListRequest) *InterestTagResponse {
	msg, e := self.Sender.CallRpcMethod("RpcBgTagList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*InterestTagResponse)
}

func (self *Brower2Backstage) RpcBgTagList_(reqMsg *ListRequest) (*InterestTagResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBgTagList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*InterestTagResponse), e
}
func (self *Brower2Backstage) RpcUpdateBgTag(reqMsg *share_message.InterestTag) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateBgTag", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateBgTag_(reqMsg *share_message.InterestTag) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateBgTag", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcBgVoiceVideoList(reqMsg *ListRequest) *BgVoiceVideoListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcBgVoiceVideoList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BgVoiceVideoListResponse)
}

func (self *Brower2Backstage) RpcBgVoiceVideoList_(reqMsg *ListRequest) (*BgVoiceVideoListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBgVoiceVideoList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BgVoiceVideoListResponse), e
}
func (self *Brower2Backstage) RpcUpdateBgVoiceVideo(reqMsg *share_message.BgVoiceVideo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateBgVoiceVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateBgVoiceVideo_(reqMsg *share_message.BgVoiceVideo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateBgVoiceVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcRevieweBgVoiceVideo(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcRevieweBgVoiceVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcRevieweBgVoiceVideo_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRevieweBgVoiceVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelBgVoiceVideo(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelBgVoiceVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelBgVoiceVideo_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelBgVoiceVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcMatchGuideList(reqMsg *ListRequest) *MatchGuideListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcMatchGuideList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MatchGuideListResponse)
}

func (self *Brower2Backstage) RpcMatchGuideList_(reqMsg *ListRequest) (*MatchGuideListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMatchGuideList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MatchGuideListResponse), e
}
func (self *Brower2Backstage) RpcUpdateMatchGuide(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateMatchGuide", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateMatchGuide_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateMatchGuide", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelMatchGuide(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelMatchGuide", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelMatchGuide_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelMatchGuide", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcSayHiList(reqMsg *ListRequest) *MatchGuideListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcSayHiList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MatchGuideListResponse)
}

func (self *Brower2Backstage) RpcSayHiList_(reqMsg *ListRequest) (*MatchGuideListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSayHiList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MatchGuideListResponse), e
}
func (self *Brower2Backstage) RpcUpdateSayHi(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateSayHi", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateSayHi_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateSayHi", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelSayHi(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelSayHi", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelSayHi_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelSayHi", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcSystemBgImageList(reqMsg *ListRequest) *SystemBgImageListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcSystemBgImageList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SystemBgImageListResponse)
}

func (self *Brower2Backstage) RpcSystemBgImageList_(reqMsg *ListRequest) (*SystemBgImageListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSystemBgImageList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SystemBgImageListResponse), e
}
func (self *Brower2Backstage) RpcSaveSystemBgImage(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveSystemBgImage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveSystemBgImage_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveSystemBgImage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelSystemBgImage(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelSystemBgImage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelSystemBgImage_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelSystemBgImage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryIntimacyConfig(reqMsg *base.Empty) *IntimacyConfigRes {
	msg, e := self.Sender.CallRpcMethod("RpcQueryIntimacyConfig", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*IntimacyConfigRes)
}

func (self *Brower2Backstage) RpcQueryIntimacyConfig_(reqMsg *base.Empty) (*IntimacyConfigRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryIntimacyConfig", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*IntimacyConfigRes), e
}
func (self *Brower2Backstage) RpcUpdateIntimacyConfig(reqMsg *IntimacyConfigRes) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateIntimacyConfig", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateIntimacyConfig_(reqMsg *IntimacyConfigRes) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateIntimacyConfig", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcVCBuryingPointReport(reqMsg *ListRequest) *VCBuryingPointReportRes {
	msg, e := self.Sender.CallRpcMethod("RpcVCBuryingPointReport", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*VCBuryingPointReportRes)
}

func (self *Brower2Backstage) RpcVCBuryingPointReport_(reqMsg *ListRequest) (*VCBuryingPointReportRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcVCBuryingPointReport", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*VCBuryingPointReportRes), e
}
func (self *Brower2Backstage) RpcCrawlPull(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcCrawlPull", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcCrawlPull_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCrawlPull", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcCrawlJobList(reqMsg *ListRequest) *CrawlJobResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCrawlJobList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CrawlJobResponse)
}

func (self *Brower2Backstage) RpcCrawlJobList_(reqMsg *ListRequest) (*CrawlJobResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCrawlJobList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CrawlJobResponse), e
}
func (self *Brower2Backstage) RpcNewsSource(reqMsg *ListRequest) *NewsSourceResponse {
	msg, e := self.Sender.CallRpcMethod("RpcNewsSource", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NewsSourceResponse)
}

func (self *Brower2Backstage) RpcNewsSource_(reqMsg *ListRequest) (*NewsSourceResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNewsSource", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NewsSourceResponse), e
}
func (self *Brower2Backstage) RpcVideoSource(reqMsg *ListRequest) *VideoSourceResponse {
	msg, e := self.Sender.CallRpcMethod("RpcVideoSource", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*VideoSourceResponse)
}

func (self *Brower2Backstage) RpcVideoSource_(reqMsg *ListRequest) (*VideoSourceResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcVideoSource", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*VideoSourceResponse), e
}
func (self *Brower2Backstage) RpcNewsList(reqMsg *ListRequest) *NewsSourceResponse {
	msg, e := self.Sender.CallRpcMethod("RpcNewsList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*NewsSourceResponse)
}

func (self *Brower2Backstage) RpcNewsList_(reqMsg *ListRequest) (*NewsSourceResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNewsList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*NewsSourceResponse), e
}
func (self *Brower2Backstage) RpcSaveNews(reqMsg *share_message.TableESPortsRealTimeInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveNews", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveNews_(reqMsg *share_message.TableESPortsRealTimeInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveNews", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelNewsSource(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelNewsSource", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelNewsSource_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelNewsSource", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelNews(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelNews", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelNews_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelNews", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcVideoList(reqMsg *ListRequest) *VideoSourceResponse {
	msg, e := self.Sender.CallRpcMethod("RpcVideoList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*VideoSourceResponse)
}

func (self *Brower2Backstage) RpcVideoList_(reqMsg *ListRequest) (*VideoSourceResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcVideoList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*VideoSourceResponse), e
}
func (self *Brower2Backstage) RpcSaveVideo(reqMsg *share_message.TableESPortsVideoInfo) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveVideo_(reqMsg *share_message.TableESPortsVideoInfo) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelVideoSource(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelVideoSource", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelVideoSource_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelVideoSource", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelVideo(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelVideo_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcChekVideo(reqMsg *QueryDataById) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcChekVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcChekVideo_(reqMsg *QueryDataById) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcChekVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcBanVideo(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcBanVideo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcBanVideo_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBanVideo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetGameList(reqMsg *ListRequest) *GameListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetGameList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GameListResponse)
}

func (self *Brower2Backstage) RpcGetGameList_(reqMsg *ListRequest) (*GameListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetGameList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GameListResponse), e
}
func (self *Brower2Backstage) RpcGetGameGuess(reqMsg *GameGuessRequest) *GameGuessResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetGameGuess", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GameGuessResponse)
}

func (self *Brower2Backstage) RpcGetGameGuess_(reqMsg *GameGuessRequest) (*GameGuessResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetGameGuess", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GameGuessResponse), e
}
func (self *Brower2Backstage) RpcEditGameGuess(reqMsg *EditGameGuessRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditGameGuess", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditGameGuess_(reqMsg *EditGameGuessRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditGameGuess", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetGameTeamInfo(reqMsg *GameGuessRequest) *GameTeamInfoResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetGameTeamInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GameTeamInfoResponse)
}

func (self *Brower2Backstage) RpcGetGameTeamInfo_(reqMsg *GameGuessRequest) (*GameTeamInfoResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetGameTeamInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GameTeamInfoResponse), e
}
func (self *Brower2Backstage) RpcGetGameRealTimeData(reqMsg *GameGuessRequest) *GameRealTimeResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetGameRealTimeData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GameRealTimeResponse)
}

func (self *Brower2Backstage) RpcGetGameRealTimeData_(reqMsg *GameGuessRequest) (*GameRealTimeResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetGameRealTimeData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GameRealTimeResponse), e
}
func (self *Brower2Backstage) RpcEditGameRealTimeData(reqMsg *EditGameRealTimeRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcEditGameRealTimeData", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcEditGameRealTimeData_(reqMsg *EditGameRealTimeRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditGameRealTimeData", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcCommentList(reqMsg *ListRequest) *CommentListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcCommentList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CommentListResponse)
}

func (self *Brower2Backstage) RpcCommentList_(reqMsg *ListRequest) (*CommentListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCommentList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CommentListResponse), e
}
func (self *Brower2Backstage) RpcDelComment(reqMsg *CommentDelRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelComment_(reqMsg *CommentDelRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcUploadComment(reqMsg *CommentUploadRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUploadComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUploadComment_(reqMsg *CommentUploadRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUploadComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetAppLabel(reqMsg *ListRequest) *SysLabelResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetAppLabel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SysLabelResponse)
}

func (self *Brower2Backstage) RpcGetAppLabel_(reqMsg *ListRequest) (*SysLabelResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAppLabel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SysLabelResponse), e
}
func (self *Brower2Backstage) RpcSaveAppLabel(reqMsg *share_message.TableESPortsLabel) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveAppLabel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveAppLabel_(reqMsg *share_message.TableESPortsLabel) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveAppLabel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetSysLabel(reqMsg *ListRequest) *SysLabelResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetSysLabel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SysLabelResponse)
}

func (self *Brower2Backstage) RpcGetSysLabel_(reqMsg *ListRequest) (*SysLabelResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetSysLabel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SysLabelResponse), e
}
func (self *Brower2Backstage) RpcSaveSysLabel(reqMsg *share_message.TableESPortsLabel) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveSysLabel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveSysLabel_(reqMsg *share_message.TableESPortsLabel) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveSysLabel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGetCarouselList(reqMsg *ListRequest) *CarouselResponse {
	msg, e := self.Sender.CallRpcMethod("RpcGetCarouselList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CarouselResponse)
}

func (self *Brower2Backstage) RpcGetCarouselList_(reqMsg *ListRequest) (*CarouselResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCarouselList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CarouselResponse), e
}
func (self *Brower2Backstage) RpcSaveCarousel(reqMsg *share_message.TableESPortsCarousel) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveCarousel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveCarousel_(reqMsg *share_message.TableESPortsCarousel) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveCarousel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcSportSysNotice(reqMsg *ListRequest) *SportSysNoticeResponse {
	msg, e := self.Sender.CallRpcMethod("RpcSportSysNotice", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SportSysNoticeResponse)
}

func (self *Brower2Backstage) RpcSportSysNotice_(reqMsg *ListRequest) (*SportSysNoticeResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSportSysNotice", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SportSysNoticeResponse), e
}
func (self *Brower2Backstage) RpcSendSportSysNotice(reqMsg *share_message.TableESPortsSysMsg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSendSportSysNotice", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSendSportSysNotice_(reqMsg *share_message.TableESPortsSysMsg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSendSportSysNotice", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcBetSlipList(reqMsg *ListRequest) *BetSlipListResponse {
	msg, e := self.Sender.CallRpcMethod("RpcBetSlipList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BetSlipListResponse)
}

func (self *Brower2Backstage) RpcBetSlipList_(reqMsg *ListRequest) (*BetSlipListResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBetSlipList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BetSlipListResponse), e
}
func (self *Brower2Backstage) RpcBetWinLosStatistics(reqMsg *ListRequest) *BetSlipStatisticsResponse {
	msg, e := self.Sender.CallRpcMethod("RpcBetWinLosStatistics", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BetSlipStatisticsResponse)
}

func (self *Brower2Backstage) RpcBetWinLosStatistics_(reqMsg *ListRequest) (*BetSlipStatisticsResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBetWinLosStatistics", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BetSlipStatisticsResponse), e
}
func (self *Brower2Backstage) RpcBetGameStatistics(reqMsg *ListRequest) *BetSlipStatisticsResponse {
	msg, e := self.Sender.CallRpcMethod("RpcBetGameStatistics", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BetSlipStatisticsResponse)
}

func (self *Brower2Backstage) RpcBetGameStatistics_(reqMsg *ListRequest) (*BetSlipStatisticsResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBetGameStatistics", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BetSlipStatisticsResponse), e
}
func (self *Brower2Backstage) RpcBetSlipReportLine(reqMsg *ListRequest) *LineChartsResponse {
	msg, e := self.Sender.CallRpcMethod("RpcBetSlipReportLine", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LineChartsResponse)
}

func (self *Brower2Backstage) RpcBetSlipReportLine_(reqMsg *ListRequest) (*LineChartsResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBetSlipReportLine", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LineChartsResponse), e
}
func (self *Brower2Backstage) RpcBetSlipReportBar(reqMsg *ListRequest) *LineChartsResponse {
	msg, e := self.Sender.CallRpcMethod("RpcBetSlipReportBar", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LineChartsResponse)
}

func (self *Brower2Backstage) RpcBetSlipReportBar_(reqMsg *ListRequest) (*LineChartsResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBetSlipReportBar", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LineChartsResponse), e
}
func (self *Brower2Backstage) RpcBetSlipOperate(reqMsg *RpcBetSlipOperateRequest) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcBetSlipOperate", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcBetSlipOperate_(reqMsg *RpcBetSlipOperateRequest) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBetSlipOperate", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcGiveWhiteList(reqMsg *ListRequest) *GiveWhiteListRes {
	msg, e := self.Sender.CallRpcMethod("RpcGiveWhiteList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GiveWhiteListRes)
}

func (self *Brower2Backstage) RpcGiveWhiteList_(reqMsg *ListRequest) (*GiveWhiteListRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGiveWhiteList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GiveWhiteListRes), e
}
func (self *Brower2Backstage) RpcAddGiveWhiteList(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcAddGiveWhiteList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcAddGiveWhiteList_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddGiveWhiteList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelGiveWhiteList(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelGiveWhiteList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelGiveWhiteList_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelGiveWhiteList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcQueryRechargeEsAct(reqMsg *base.Empty) *share_message.Activity {
	msg, e := self.Sender.CallRpcMethod("RpcQueryRechargeEsAct", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.Activity)
}

func (self *Brower2Backstage) RpcQueryRechargeEsAct_(reqMsg *base.Empty) (*share_message.Activity, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryRechargeEsAct", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.Activity), e
}
func (self *Brower2Backstage) RpcUpdateRechargeEsAct(reqMsg *share_message.Activity) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateRechargeEsAct", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcUpdateRechargeEsAct_(reqMsg *share_message.Activity) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateRechargeEsAct", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcRechargeEsCfg(reqMsg *base.Empty) *RechargeEsCfgRes {
	msg, e := self.Sender.CallRpcMethod("RpcRechargeEsCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*RechargeEsCfgRes)
}

func (self *Brower2Backstage) RpcRechargeEsCfg_(reqMsg *base.Empty) (*RechargeEsCfgRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRechargeEsCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*RechargeEsCfgRes), e
}
func (self *Brower2Backstage) RpcSaveRechargeEsCfg(reqMsg *share_message.TableESportsExchangeCfg) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSaveRechargeEsCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcSaveRechargeEsCfg_(reqMsg *share_message.TableESportsExchangeCfg) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSaveRechargeEsCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcDelRechargeEsCfg(reqMsg *QueryDataByIds) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcDelRechargeEsCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Backstage) RpcDelRechargeEsCfg_(reqMsg *QueryDataByIds) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelRechargeEsCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Backstage) RpcPointsReportList(reqMsg *ListRequest) *PointsReportRes {
	msg, e := self.Sender.CallRpcMethod("RpcPointsReportList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PointsReportRes)
}

func (self *Brower2Backstage) RpcPointsReportList_(reqMsg *ListRequest) (*PointsReportRes, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPointsReportList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PointsReportRes), e
}
func (self *Brower2Backstage) RpcQueryeSportCoinLog(reqMsg *ListRequest) *SportCoinLogResponse {
	msg, e := self.Sender.CallRpcMethod("RpcQueryeSportCoinLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SportCoinLogResponse)
}

func (self *Brower2Backstage) RpcQueryeSportCoinLog_(reqMsg *ListRequest) (*SportCoinLogResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryeSportCoinLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SportCoinLogResponse), e
}

// ==========================================================
type IBackstage2Brower interface {
	RpcTestPush(reqMsg *PushRequest) *PushResponse
	RpcTestPush_(reqMsg *PushRequest) (*PushResponse, easygo.IRpcInterrupt)
	RpcNewPush(reqMsg *base.Empty)
	RpcReplacePush(reqMsg *ErrMessage)
	RpcPushIMmessage(reqMsg *share_message.IMmessage) *QueryDataById
	RpcPushIMmessage_(reqMsg *share_message.IMmessage) (*QueryDataById, easygo.IRpcInterrupt)
	RpcSendShopOrderExpress(reqMsg *QueryShopOrderExpressResponse)
	RpcCrawlPush(reqMsg *share_message.TableCrawlJob)
	RpcToolLuckyPush(reqMsg *base.Empty)
}

type Backstage2Brower struct {
	Sender easygo.IMessageSender
}

func (self *Backstage2Brower) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Backstage2Brower) RpcTestPush(reqMsg *PushRequest) *PushResponse {
	msg, e := self.Sender.CallRpcMethod("RpcTestPush", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PushResponse)
}

func (self *Backstage2Brower) RpcTestPush_(reqMsg *PushRequest) (*PushResponse, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTestPush", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PushResponse), e
}
func (self *Backstage2Brower) RpcNewPush(reqMsg *base.Empty) {
	self.Sender.CallRpcMethod("RpcNewPush", reqMsg)
}
func (self *Backstage2Brower) RpcReplacePush(reqMsg *ErrMessage) {
	self.Sender.CallRpcMethod("RpcReplacePush", reqMsg)
}
func (self *Backstage2Brower) RpcPushIMmessage(reqMsg *share_message.IMmessage) *QueryDataById {
	msg, e := self.Sender.CallRpcMethod("RpcPushIMmessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryDataById)
}

func (self *Backstage2Brower) RpcPushIMmessage_(reqMsg *share_message.IMmessage) (*QueryDataById, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPushIMmessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryDataById), e
}
func (self *Backstage2Brower) RpcSendShopOrderExpress(reqMsg *QueryShopOrderExpressResponse) {
	self.Sender.CallRpcMethod("RpcSendShopOrderExpress", reqMsg)
}
func (self *Backstage2Brower) RpcCrawlPush(reqMsg *share_message.TableCrawlJob) {
	self.Sender.CallRpcMethod("RpcCrawlPush", reqMsg)
}
func (self *Backstage2Brower) RpcToolLuckyPush(reqMsg *base.Empty) {
	self.Sender.CallRpcMethod("RpcToolLuckyPush", reqMsg)
}
