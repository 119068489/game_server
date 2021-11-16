package for_game

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"sync"
	"time"
)

//管理redis对象
type RegisObjManager struct {
	sync.Map
	//Mutex easygo.RLock
	Mutex easygo.Mutex
	Sid   int32 //所在服务器id
}

func NewRedisObjManager(sid int32) *RegisObjManager {
	p := &RegisObjManager{
		Sid: sid,
	}
	return p
}

var OrderMgr *RegisObjManager
var AccountMgr *RegisObjManager
var PlayerBaseMgr *RegisObjManager
var RedPacketMgr *RegisObjManager
var RedPacketTotalMgr *RegisObjManager
var TransferMoneyMgr *RegisObjManager
var GoldChangeLogMgr *RegisObjManager
var CoinChangeLogMgr *RegisObjManager
var ESportCoinChangeLogMgr *RegisObjManager
var PersonalChatLogMgr *RegisObjManager
var TeamMgr *RegisObjManager
var TeamPersonalMgr *RegisObjManager
var TeamChatLogMgr *RegisObjManager
var RegisterLoginReportMgr *RegisObjManager
var ArticleReportMgr *RegisObjManager
var InOutCashSumReportMgr *RegisObjManager
var NoticeReportMgr *RegisObjManager
var PlayerKeepReportMgr *RegisObjManager
var PlayerBehaviorReportMgr *RegisObjManager
var OperationChannelReportMgr *RegisObjManager
var RecallReportMgr *RegisObjManager
var PlayerBagItemMgr *RegisObjManager
var PlayerEquipmentMgr *RegisObjManager

var ESportPlayerMgr *RegisObjManager
var ESportLiveRoomPlayerMgr *RegisObjManager
var ESportFollowMgr *RegisObjManager
var ESportThumbsUpMgr *RegisObjManager
var ESportRoomChatMgr *RegisObjManager
var ESportBpsDurationMgr *RegisObjManager
var ESportBpsClickMgr *RegisObjManager
var ESportTableFedUpdateMgr *RegisObjManager

var ChatSessionMgr *RegisObjManager
var ButtonClickReportMgr *RegisObjManager
var WishPlayerMgr *RegisObjManager
var DiamondChangeLogMgr *RegisObjManager
var PlayerIntimacyMgr *RegisObjManager
var VCBuryingPointMgr *RegisObjManager

var WishLogReportMgr *RegisObjManager

//初始化redis模块
func InitRedisObjManager(sid int32) {
	OrderMgr = NewRedisObjManager(sid)
	AccountMgr = NewRedisObjManager(sid)
	PlayerBaseMgr = NewRedisObjManager(sid)
	RedPacketMgr = NewRedisObjManager(sid)
	RedPacketTotalMgr = NewRedisObjManager(sid)
	TransferMoneyMgr = NewRedisObjManager(sid)
	GoldChangeLogMgr = NewRedisObjManager(sid)
	CoinChangeLogMgr = NewRedisObjManager(sid)
	ESportCoinChangeLogMgr = NewRedisObjManager(sid)
	PersonalChatLogMgr = NewRedisObjManager(sid)
	TeamMgr = NewRedisObjManager(sid)
	TeamPersonalMgr = NewRedisObjManager(sid)
	TeamChatLogMgr = NewRedisObjManager(sid)

	RegisterLoginReportMgr = NewRedisObjManager(sid)
	ArticleReportMgr = NewRedisObjManager(sid)
	InOutCashSumReportMgr = NewRedisObjManager(sid)
	NoticeReportMgr = NewRedisObjManager(sid)
	PlayerKeepReportMgr = NewRedisObjManager(sid)
	PlayerBehaviorReportMgr = NewRedisObjManager(sid)
	OperationChannelReportMgr = NewRedisObjManager(sid)
	RecallReportMgr = NewRedisObjManager(sid)
	PlayerBagItemMgr = NewRedisObjManager(sid)
	PlayerEquipmentMgr = NewRedisObjManager(sid)
	ChatSessionMgr = NewRedisObjManager(sid)
	ButtonClickReportMgr = NewRedisObjManager(sid)
	WishPlayerMgr = NewRedisObjManager(sid)
	DiamondChangeLogMgr = NewRedisObjManager(sid)
	PlayerIntimacyMgr = NewRedisObjManager(sid)
	VCBuryingPointMgr = NewRedisObjManager(sid)
	//电竞reids管理器初始化
	ESportPlayerMgr = NewRedisObjManager(sid)
	ESportLiveRoomPlayerMgr = NewRedisObjManager(sid)
	ESportFollowMgr = NewRedisObjManager(sid)
	ESportThumbsUpMgr = NewRedisObjManager(sid)
	//放映厅聊天管理
	ESportRoomChatMgr = NewRedisObjManager(sid)
	//埋点redis管理
	ESportBpsDurationMgr = NewRedisObjManager(sid)    //停留时长埋点
	ESportBpsClickMgr = NewRedisObjManager(sid)       //点击埋点
	ESportTableFedUpdateMgr = NewRedisObjManager(sid) //数据表某个字段修改
	//许愿池埋点redis管理
	WishLogReportMgr = NewRedisObjManager(sid)
}

func (self *RegisObjManager) GetSid() int32 {
	return self.Sid
}

//对外方法，获取订单管理对象，如果为nil表示redis内存不存在，数据库也不存在
//========================Order start==============================================
func (self *RegisObjManager) LoadOrderObj(id string) *RedisOrderObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisOrderObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisOrderObj(id string) *RedisOrderObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadOrderObj(id)
	if obj == nil {
		obj = NewRedisOrder(id)
	}
	return obj
}

////创建新订单
func (self *RegisObjManager) CreateRedisOrderObj(order *share_message.Order) *RedisOrderObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var orderNo string
	if order.OrderId == nil || order.GetOrderId() == "" {
		orderNo = RedisCreateOrderNo(order.GetChangeType(), order.GetSourceType())
		order.OrderId = easygo.NewString(orderNo)
	} else {
		orderNo = order.GetOrderId()
	}
	obj := NewRedisOrder(orderNo, order)
	obj.SetSaveStatus(true) //标志数据要存储
	return obj
}

//========================Order end==============================================
//========================Account start==============================================
func (self *RegisObjManager) LoadAccountObj(id PLAYER_ID) *RedisAccountObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisAccountObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisAccountObj(id PLAYER_ID) *RedisAccountObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := AccountMgr.LoadAccountObj(id)
	if obj == nil {
		obj = NewRedisAccount(id)
	}
	return obj
}

////创建新订单
func (self *RegisObjManager) CreateRedisAccount(account *share_message.PlayerAccount) *RedisAccountObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := NewRedisAccount(account.GetPlayerId(), account)
	obj.SetSaveStatus(true)
	return obj
}

//========================Account end==============================================
//========================PlayerBase start==============================================
func (self *RegisObjManager) LoadPlayerBaseObj(id PLAYER_ID) *RedisPlayerBaseObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisPlayerBaseObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisPlayerBaseObj(id PLAYER_ID, player ...*share_message.PlayerBase) *RedisPlayerBaseObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadPlayerBaseObj(id)
	if obj == nil {
		obj = NewRedisPlayerBase(id, player...)
	}
	return obj
}

//========================PlayerBase end==============================================
//========================RedPacket start==============================================
func (self *RegisObjManager) LoadRedPacketObj(id REDPACKET_ID) *RedisRedPacketObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisRedPacketObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisRedPacketObj(id REDPACKET_ID) *RedisRedPacketObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := RedPacketMgr.LoadRedPacketObj(id)
	if obj == nil {
		obj = NewRedisRedPacket(id)
	}
	return obj
}

////创建新订单
func (self *RegisObjManager) CreateRedisRedPacket(msg *share_message.RedPacket) *RedisRedPacketObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	id := NextId(TABLE_RED_PACKET) //获取红包Redis自增Id
	sendId := msg.GetSender()
	player := GetRedisPlayerBase(sendId)

	name := GetTeamReName(msg.GetType(), msg.GetTargetId(), sendId, player.GetNickName())
	redPacket := &share_message.RedPacket{
		Id:         easygo.NewInt64(id),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
		Type:       easygo.NewInt32(msg.GetType()),
		SenderName: easygo.NewString(name),
		SenderHead: easygo.NewString(player.GetHeadIcon()),
		Sender:     easygo.NewInt64(msg.GetSender()),
		TargetId:   easygo.NewInt64(msg.GetTargetId()),
		TotalMoney: easygo.NewInt64(msg.GetTotalMoney()),
		TotalCount: easygo.NewInt32(msg.GetTotalCount()),
		PerMoney:   easygo.NewInt64(msg.GetPerMoney()),
		Content:    easygo.NewString(msg.GetContent()),
		CurMoney:   easygo.NewInt64(msg.GetTotalMoney()),
		CurCount:   easygo.NewInt32(msg.GetTotalCount()),
		State:      easygo.NewInt32(PACKET_MONEY_OPEN),
		Packets:    msg.GetPackets(),
		PlayerList: msg.GetPlayerList(),
		OrderId:    easygo.NewString(msg.GetOrderId()),
		PayWay:     easygo.NewInt32(msg.GetPayWay()),
		Sex:        easygo.NewInt32(player.GetSex()),
	}
	obj := NewRedisRedPacket(id, redPacket)
	obj.SetSaveStatus(true)
	return obj
}

//========================RedPacket end==============================================
//========================RedPacketTotal start==============================================
func (self *RegisObjManager) LoadRedPacketTotalObj(id string) *RedisRedPacketTotalObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisRedPacketTotalObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisRedPackeTotalObj(id string) *RedisRedPacketTotalObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadRedPacketTotalObj(id)
	if obj == nil {
		obj = NewRedPacketTotal(id)
	}
	return obj
}

////创建新订单
func (self *RegisObjManager) UpdateRedisRedPacketTotal(pid, ti, value int64, t int32) {
	id := MakeRedisRedPacketTotalId(pid, ti)
	obj := GetRedisRedPacketTotal(id)
	if obj == nil {
		obj = CreateRedPacketTotalInfo(pid, ti)
	}
	switch t {
	case REDPACKET_STATISTICS_RECV:
		obj.IncrRecTotalMoney(value)
		obj.IncrRecCount(1)
	case REDPACKET_STATISTICS_SEND:
		obj.IncrSendTotalMoney(value)
		obj.IncrSendCount(1)
	case REDPACKET_STATISTICS_LUCK:
		obj.IncrLuckCnt(1)
	}
}

//========================RedPacketTotal end==============================================
//========================TransferMoney start==============================================
func (self *RegisObjManager) LoadTransferMoneyObj(id int64) *RedisTransferMoneyObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisTransferMoneyObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisTransferMoneyObj(id int64) *RedisTransferMoneyObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadTransferMoneyObj(id)
	if obj == nil {
		obj = NewRedisTransferMoney(id)
	}
	return obj
}

////创建新订单
func (self *RegisObjManager) CreateTransferMoney(msg *share_message.TransferMoney, orderId string) *RedisTransferMoneyObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	id := NextId(TABLE_TRANSFER_MONEY) //获取Redis自增Id
	transferAccounts := &share_message.TransferMoney{
		Id:         easygo.NewInt64(id),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
		Sender:     easygo.NewInt64(msg.GetSender()),
		TargetId:   easygo.NewInt64(msg.GetTargetId()),
		Way:        easygo.NewInt32(msg.GetWay()),
		Card:       easygo.NewString(msg.GetCard()),
		Gold:       easygo.NewInt64(msg.GetGold()),
		Content:    easygo.NewString(msg.GetContent()),
		State:      easygo.NewInt32(TRANSFER_MONEY_OPEN),
		OrderId:    easygo.NewString(orderId),
	}
	obj := NewRedisTransferMoney(id, transferAccounts)
	return obj
}

//========================TransferMoney end==============================================
//========================GoldChangeLog start==============================================
func (self *RegisObjManager) LoadGoldLogObj(id string) *RedisGoldLogObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisGoldLogObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisGoldLogObj() *RedisGoldLogObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadGoldLogObj(TABLE_GOLDCHANGELOG)
	if obj == nil {
		obj = NewRedisGoldLog()
	}
	return obj
}

////创建新订单
func (self *RegisObjManager) AddGoldChangeLog(log *CommonGold) {
	obj := self.GetRedisGoldLogObj()
	obj.AddRedisGoldLog(log)
}

//========================GoldChangeLog end==============================================
//========================CoinChangeLog start==============================================
func (self *RegisObjManager) LoadCoinLogObj(id string) *RedisCoinLogObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisCoinLogObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisCoinLogObj() *RedisCoinLogObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadCoinLogObj(TABLE_COINCHANGELOG)
	if obj == nil {
		obj = NewRedisCoinLog()
	}
	return obj
}

////创建新订单
func (self *RegisObjManager) AddCoinChangeLog(log *CommonCoin) {
	obj := self.GetRedisCoinLogObj()
	obj.AddRedisCoinLog(log)
}

//========================CoinChangeLog end==============================================

//========================ESportCoinChangeLog start==============================================
func (self *RegisObjManager) LoadESportCoinLogObj(id string) *RedisESportCoinLogObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisESportCoinLogObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisESportCoinLogObj() *RedisESportCoinLogObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadESportCoinLogObj(TABLE_ESPORTCHANGELOG)
	if obj == nil {
		obj = NewRedisESportCoinLog()
	}
	return obj
}

////创建新订单
func (self *RegisObjManager) AddESportCoinChangeLog(log *CommonESportCoin) {
	obj := self.GetRedisESportCoinLogObj()
	obj.AddRedisESportCoinLog(log)
}

//========================ESportCoinChangeLog end==============================================

//========================PersonalChatLog start==============================================
func (self *RegisObjManager) LoadPersonalChatLogObj(id string) *RedisPersonalChatLogObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisPersonalChatLogObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisPersonalChatLogObj(sessionId string) *RedisPersonalChatLogObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadPersonalChatLogObj(sessionId)
	if obj == nil {
		obj = NewRedisPersonalChatLog(sessionId)
	} else {
		obj.SetSaveStatus(true) //恢复已经存在redis数据库的存储
	}
	return obj
}

////创建新聊天记录
func (self *RegisObjManager) AddPersonalChatLog(msg *share_message.Chat) int64 {
	obj := self.GetRedisPersonalChatLogObj(msg.GetSessionId())
	talkId := msg.GetSourceId()
	targetId := msg.GetTargetId()
	logId := NextId(TABLE_PERSONAL_CHAT_LOG)
	sessionObj := GetRedisChatSessionObj(msg.GetSessionId())
	talkLogId := sessionObj.GetNextMaxLogId()
	log := &share_message.PersonalChatLog{
		LogId:     easygo.NewInt64(logId),
		Talker:    easygo.NewInt64(talkId),
		TargetId:  easygo.NewInt64(targetId),
		Time:      easygo.NewInt64(msg.GetTime()),
		Content:   easygo.NewString(msg.GetContent()),
		Type:      easygo.NewInt32(msg.GetContentType()),
		Cite:      easygo.NewString(msg.GetCite()),
		QPId:      easygo.NewInt64(msg.GetQPId()),
		SessionId: easygo.NewString(msg.GetSessionId()),
		Status:    easygo.NewInt32(TALK_STATUS_NORMAL),
		TalkLogId: easygo.NewInt64(talkLogId),
		Mark:      easygo.NewString(msg.GetMark()),
	}
	obj.AddRedisPersonalChatLog(log)
	return talkLogId
}

//========================PersonalChatLog end==============================================
//========================Team start==============================================
func (self *RegisObjManager) LoadTeamObj(id PLAYER_ID) *RedisTeamObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisTeamObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisTeamObj(id int64, team ...*share_message.TeamData) *RedisTeamObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadTeamObj(id)
	if obj == nil {
		obj = NewRedisTeamObj(id, team...)
	}
	return obj
}

//========================Team end==============================================
//========================TeamPersonal start==============================================
func (self *RegisObjManager) LoadTeamPersonalObj(id int64) *RedisTeamPersonalObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisTeamPersonalObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisTeamPersonalObj(id int64, data ...[]*share_message.PersonalTeamData) *RedisTeamPersonalObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadTeamPersonalObj(id)
	if obj == nil {
		obj = NewRedisTeamPersonalObj(id, data...)
	}
	return obj
}

//========================TeamPersonal end==============================================
//========================TeamChatLog start==============================================
func (self *RegisObjManager) LoadTeamChatLogObj(id int64) *RedisTeamChatLogObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisTeamChatLogObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisTeamChatLogObj(id int64, data ...[]*share_message.TeamChatLog) *RedisTeamChatLogObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadTeamChatLogObj(id)
	if obj == nil {
		obj = NewRedisTeamChatLog(id, data...)
	}
	return obj
}

//========================TeamChatLog end==============================================

//==================================================================================================报表相关
//========================RegisterLoginReport start==============================================
func (self *RegisObjManager) LoadRegisterLoginObj(id int64) *RegisterLoginReportObj {
	id = easygo.Get0ClockTimestamp(id)
	value, ok := self.Load(id)
	if ok {
		return value.(*RegisterLoginReportObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisRegisterLoginObj(id int64) *RegisterLoginReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadRegisterLoginObj(id)
	if obj == nil {
		obj = NewRedisRegisterLoginReport(id)
		if obj == nil {
			obj = self.CreateRedisRegisterLoginObj(id)
		}
	}
	return obj
}

//创建一个新的key
func (self *RegisObjManager) CreateRedisRegisterLoginObj(id int64) *RegisterLoginReportObj {
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadRegisterLoginObj(id)
	if obj == nil {
		report := &share_message.RegisterLoginReport{
			CreateTime: easygo.NewInt64(id),
		}
		obj = NewRedisRegisterLoginReport(id, report)
	}
	return obj
}

//========================RegisterLoginReport end==============================================

//========================ArticleReport start==============================================
func (self *RegisObjManager) LoadArticleReportObj(id int64) *ArticleReportObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*ArticleReportObj)
	}
	return nil
}

func (self *RegisObjManager) GetRedisArticleReport(id int64) *ArticleReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadArticleReportObj(id)
	if obj == nil {
		obj = NewRedisArticleReport(id)
		if obj == nil {
			obj = self.CreateRedisArticleReport(id)
		}
	}
	return obj
}

func (self *RegisObjManager) CreateRedisArticleReport(id int64) *ArticleReportObj {
	obj := self.LoadArticleReportObj(id)
	if obj == nil {
		report := &share_message.ArticleReport{
			Id:         easygo.NewInt64(id),
			CreateTime: easygo.NewInt64(util.GetMilliTime()),
		}
		obj = NewRedisArticleReport(id, report)
	}
	return obj
}

//========================ArticleReport end==============================================

//========================InOutCashSumReport start==============================================
func (self *RegisObjManager) LoadInOutCashSumReportObj(id int64) *InOutCashSumReportObj {
	id = easygo.Get0ClockTimestamp(id)
	value, ok := self.Load(id)
	if ok {
		return value.(*InOutCashSumReportObj)
	}
	return nil
}

func (self *RegisObjManager) GetRedisInOutCashSumReport(id int64) *InOutCashSumReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadInOutCashSumReportObj(id)
	if obj == nil {
		obj = NewRedisInOutCashSumReport(id)
		if obj == nil {
			obj = self.CreateRedisInOutCashSumReport(id)
		}
	}
	return obj
}

func (self *RegisObjManager) CreateRedisInOutCashSumReport(id int64) *InOutCashSumReportObj {
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadInOutCashSumReportObj(id)
	if obj == nil {
		report := &share_message.InOutCashSumReport{
			CreateTime: easygo.NewInt64(id),
		}
		obj = NewRedisInOutCashSumReport(id, report)
	}
	return obj
}

//========================InOutCashSumReport end==============================================

//========================NoticeReport start==============================================
func (self *RegisObjManager) LoadNoticeReportObj(id int32) *NoticeReportObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*NoticeReportObj)
	}
	return nil
}

func (self *RegisObjManager) GetRedisNoticeReport(id int32) *NoticeReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadNoticeReportObj(id)
	if obj == nil {
		obj = NewRedisNoticeReport(id)
		if obj == nil {
			obj = self.CreateRedisNoticeReport(id)
		}
	}
	return obj
}

func (self *RegisObjManager) CreateRedisNoticeReport(id int32) *NoticeReportObj {
	obj := self.LoadNoticeReportObj(id)
	if obj == nil {
		report := &share_message.ArticleReport{
			Id:         easygo.NewInt64(id),
			CreateTime: easygo.NewInt64(util.GetMilliTime()),
		}
		obj = NewRedisNoticeReport(id, report)
	}
	return obj
}

//========================NoticeReport end==============================================

//========================ButtonClickReport start==============================================
func (self *RegisObjManager) LoadButtonClickReportObj(id int64) *ButtonClickReportObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*ButtonClickReportObj)
	}
	return nil
}

func (self *RegisObjManager) GetRedisButtonClickReport(id int64) *ButtonClickReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadButtonClickReportObj(id)
	if obj == nil {
		obj = NewRedisButtonClickReport(id)
		if obj == nil {
			obj = self.CreateRedisButtonClickReport(id)
		}
	}
	return obj
}

func (self *RegisObjManager) CreateRedisButtonClickReport(id int64) *ButtonClickReportObj {
	obj := self.LoadButtonClickReportObj(id)
	if obj == nil {
		report := &share_message.ButtonClickReport{
			CreateTime: easygo.NewInt64(easygo.GetToday0ClockTimestamp()),
		}
		obj = NewRedisButtonClickReport(id, report)
	}
	return obj
}

//========================ButtonClickReport end==============================================

//========================PlayerKeepReport start==============================================
func (self *RegisObjManager) LoadPlayerKeepReportObj(id int64) *PlayerKeepReportObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*PlayerKeepReportObj)
	}
	return nil
}

func (self *RegisObjManager) GetRedisPlayerKeepReport(id int64) *PlayerKeepReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadPlayerKeepReportObj(id)
	if obj == nil {
		obj = NewRedisPlayerKeepReport(id)
		if obj == nil {
			obj = self.CreateRedisPlayerKeepReport(id)
		}
	}
	return obj
}

func (self *RegisObjManager) CreateRedisPlayerKeepReport(id int64) *PlayerKeepReportObj {
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadPlayerKeepReportObj(id)
	if obj == nil {
		report := &share_message.PlayerKeepReport{
			CreateTime: easygo.NewInt64(id),
		}
		obj = NewRedisPlayerKeepReport(id, report)
	}
	return obj
}

//========================PlayerKeepReport end==============================================

//========================PlayerBehaviorReport start==============================================
func (self *RegisObjManager) LoadPlayerBehaviorReportObj(id int64) *PlayerBehaviorReportObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*PlayerBehaviorReportObj)
	}
	return nil
}

func (self *RegisObjManager) GetRedisPlayerBehaviorReport(id int64) *PlayerBehaviorReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadPlayerBehaviorReportObj(id)
	if obj == nil {
		obj = NewRedisPlayerBehaviorReport(id)
		if obj == nil {
			obj = self.CreateRedisPlayerBehaviorReport(id)
		}
	}
	return obj
}

func (self *RegisObjManager) CreateRedisPlayerBehaviorReport(id int64) *PlayerBehaviorReportObj {
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadPlayerBehaviorReportObj(id)
	if obj == nil {
		report := &share_message.PlayerBehaviorReport{
			CreateTime: easygo.NewInt64(id),
		}
		obj = NewRedisPlayerBehaviorReport(id, report)
	}
	return obj
}

//========================PlayerBehaviorReport end==============================================

//========================OperationChannelReport start==============================================
func (self *RegisObjManager) LoadOperationChannelReportObj(id string) *OperationChannelReportObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*OperationChannelReportObj)
	}
	return nil
}

func (self *RegisObjManager) GetRedisOperationChannelReport(id, channle string, createTime int64) *OperationChannelReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadOperationChannelReportObj(id)
	if obj == nil {
		obj = NewRedisOperationChannelReport(id)
		if obj == nil {
			obj = self.CreateRedisOperationChannelReport(id, channle, createTime)
		}
	}
	return obj
}

func (self *RegisObjManager) CreateRedisOperationChannelReport(id, channle string, createTime int64) *OperationChannelReportObj {
	obj := self.LoadOperationChannelReportObj(id)
	if obj == nil {
		channel := QueryOperationByNo(channle) //查询渠道信息
		report := &share_message.OperationChannelReport{
			Id:          easygo.NewString(id),
			ChannelNo:   easygo.NewString(channel.GetChannelNo()),
			ChannelName: easygo.NewString(channel.GetName()),
			Cooperation: easygo.NewInt32(channel.GetCooperation()),
			CreateTime:  easygo.NewInt64(easygo.Get0ClockTimestamp(createTime)),
		}

		obj = NewRedisOperationChannelReport(id, report)
	}
	return obj
}

//========================OperationChannelReport end==============================================

//========================RecallReport start==============================================
func (self *RegisObjManager) LoadRecallReportObj(id int64) *RecallReportObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RecallReportObj)
	}
	return nil
}

func (self *RegisObjManager) GetRedisRecallReport(id int64) *RecallReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadRecallReportObj(id)
	if obj == nil {
		obj = NewRedisRecallReport(id)
		if obj == nil {
			obj = self.CreateRedisRecallReport(id)
		}
	}
	return obj
}

func (self *RegisObjManager) CreateRedisRecallReport(id int64) *RecallReportObj {
	obj := self.LoadRecallReportObj(id)
	if obj == nil {
		report := &share_message.RecallReport{
			CreateTime: easygo.NewInt64(id),
		}
		obj = NewRedisRecallReport(id, report)
	}
	return obj
}

//========================RecallReport end==============================================
//========================PlayerBagItem start==============================================
func (self *RegisObjManager) LoadPlayerBagItemObj(id int64) *RedisPlayerBagItemObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisPlayerBagItemObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisPlayerBagItemObj(id int64, data ...[]*share_message.PlayerBagItem) *RedisPlayerBagItemObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadPlayerBagItemObj(id)
	if obj == nil {
		obj = NewRedisPlayerBagItem(id, data...)
	}
	return obj
}

//========================PlayerBagItem end
//========================PlayerEquipment start==============================================
func (self *RegisObjManager) LoadPlayerEquipmentObj(id PLAYER_ID) *RedisPlayerEquipmentObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisPlayerEquipmentObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisPlayerEquipmentObj(id PLAYER_ID, equipment ...*share_message.PlayerEquipment) *RedisPlayerEquipmentObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadPlayerEquipmentObj(id)
	if obj == nil {
		obj = NewRedisPlayerEquipment(id, equipment...)
	}
	return obj
}

//========================PlayerEquipment end==============================================

//========================ChatSession start==============================================
func (self *RegisObjManager) LoadChatSessionObj(id string) *RedisChatSessionObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisChatSessionObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisChatSessionObj(id string, session ...*share_message.ChatSession) *RedisChatSessionObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadChatSessionObj(id)
	if obj == nil {
		obj = NewRedisChatSessionObj(id, session...)
	}
	return obj
}

//========================ChatSession end==============================================
//========================WishPlayer start==============================================
func (self *RegisObjManager) LoadWishPlayerObj(id PLAYER_ID) *RedisWishPlayerObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisWishPlayerObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisWishPlayerObj(id PLAYER_ID, player ...*share_message.WishPlayer) *RedisWishPlayerObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadWishPlayerObj(id)
	if obj == nil {
		obj = NewWishPlayerObj(id, player...)
	}
	return obj
}

//========================WishPlayer end==============================================
//========================DiamondChangeLog start==============================================
func (self *RegisObjManager) LoadDiamondLogObj(id string) *RedisDiamondLogObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisDiamondLogObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisDiamondLogObj() *RedisDiamondLogObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadDiamondLogObj(TABLE_DIAMOND_CHANGELOG)
	if obj == nil {
		obj = NewRedisDiamondLog()
	}
	return obj
}
func (self *RegisObjManager) GetRedisDiamondLogObjUnLock() *RedisDiamondLogObj {
	obj := self.LoadDiamondLogObj(TABLE_DIAMOND_CHANGELOG)
	if obj == nil {
		obj = NewRedisDiamondLog()
	}
	return obj
}

////创建新日志
func (self *RegisObjManager) AddDiamondChangeLog(log *CommonDiamond) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.GetRedisDiamondLogObjUnLock()
	obj.AddRedisDiamondLog(log)
}

//========================GoldChangeLog end==============================================
//========================PlayerIntimacy start==============================================
func (self *RegisObjManager) LoadPlayerIntimacyObj(id string) *RedisPlayerIntimacyObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*RedisPlayerIntimacyObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisPlayerIntimacyObj(id string, session ...*share_message.PlayerIntimacy) *RedisPlayerIntimacyObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadPlayerIntimacyObj(id)
	if obj == nil {
		obj = NewRedisPlayerIntimacy(id, session...)
		if obj == nil {
			if len(id) == 21 { //id组装规则，一定是这个长度
				return self.CreateNewPlayerIntimacy(id)
			}
		}
	}
	return obj
}

//创建新的亲密度
func (self *RegisObjManager) CreateNewPlayerIntimacy(id string) *RedisPlayerIntimacyObj {
	data := &share_message.PlayerIntimacy{
		Id:          easygo.NewString(id),
		IntimacyVal: easygo.NewInt64(0),
		IntimacyLv:  easygo.NewInt32(0),
		LastTime:    easygo.NewInt64(0),
		IsSayHi:     easygo.NewBool(false),
	}

	obj := NewRedisPlayerIntimacy(id, data)
	//创建后，先存下数据库
	obj.SaveToMongo()
	return obj
}

//========================PlayerIntimacy end==============================================

//========================VCBuryingPoint start==============================================
func (self *RegisObjManager) LoadVCBuryingPointReportObj(id int64) *VCBuryingPointReportObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*VCBuryingPointReportObj)
	}
	return nil
}
func (self *RegisObjManager) GetRedisVCBuryingPointReportObj(id int64, data ...*share_message.VCBuryingPointReport) *VCBuryingPointReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadVCBuryingPointReportObj(id)
	if obj == nil {
		obj = NewRedisVCBuryingPointReportObj(id, data...)
		if obj == nil {
			//创建新报表
			obj = self.CreateNewVCBuryingPointReport(id)
		}
	}
	return obj
}
func (self *RegisObjManager) CreateNewVCBuryingPointReport(id int64) *VCBuryingPointReportObj {
	data := &share_message.VCBuryingPointReport{
		Id: easygo.NewInt64(id),
	}

	obj := NewRedisVCBuryingPointReportObj(id, data)
	//创建后，先不存下数据库
	// obj.SaveToMongo()
	return obj
}

//========================VCBuryingPoint end==============================================

//========================WishLogReport start==============================================
func (self *RegisObjManager) LoadWishLogReportObj(id int64) *WishLogReportObj {
	value, ok := self.Load(id)
	if ok {
		return value.(*WishLogReportObj)
	}
	return nil
}

func (self *RegisObjManager) GetRedisWishLogReport(id int64) *WishLogReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadWishLogReportObj(id)
	if obj == nil {
		obj = NewRedisWishLogReport(id)
		if obj == nil {
			obj = self.CreateRedisWishLogReport(id)
		}
	}
	return obj
}

func (self *RegisObjManager) SetRedisWishLogReport(id int64) *WishLogReportObj {
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadWishLogReportObj(id)
	if obj == nil {
		report := &share_message.WishLogReport{
			CreateTime: easygo.NewInt64(id),
		}
		obj = NewRedisWishLogReport(id, report)
	}
	return obj
}

func (self *RegisObjManager) CreateRedisWishLogReport(id int64) *WishLogReportObj {
	id = easygo.Get0ClockTimestamp(id)
	obj := self.LoadWishLogReportObj(id)
	if obj == nil {
		report := &share_message.WishLogReport{
			CreateTime: easygo.NewInt64(id),
		}
		obj = NewRedisWishLogReport(id, report)
	}
	return obj
}

//========================WishLogReport end==============================================
