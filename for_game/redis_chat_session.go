package for_game

import (
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

//会话数据
const REDIS_CHAT_SESSION_EXIST_LIST = "chat_session_exist_list"
const REDIS_CHAT_SESSION_EXIST_TIME = 1000 * 600
const REDIS_CHAT_SESSION_PLAYERS = "PlayerIds"
const REDIS_CHAT_SESSION_READINFO = "ReadInfo"

type ChatSessionEx struct {
	Id             string `json:"_id"`
	Type           int32  //1、私人会话，2群会话
	SessionName    string //会话名称
	TeamName       string //群名
	SessionHeadUrl string //会话头像
	MaxLogId       int64  // 消息总条数
	Topic          string //话题
}
type RedisChatSessionObj struct {
	Id string //会话id
	RedisBase
}

//写入redis
func NewRedisChatSessionObj(id string, session ...*share_message.ChatSession) *RedisChatSessionObj {
	p := &RedisChatSessionObj{
		Id: id,
	}
	obj := append(session, nil)[0]
	return p.Init(obj)
}

func (self *RedisChatSessionObj) Init(obj *share_message.ChatSession) *RedisChatSessionObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_MONGODB_CHAT_SESSION)
	self.Sid = ChatSessionMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		ChatSessionMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryChatSession(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisChatSession(obj)
	}
	return self
}
func (self *RedisChatSessionObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisChatSessionObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_MONGODB_CHAT_SESSION, self.Id)
}

//定时更新数据
func (self *RedisChatSessionObj) UpdateData() { //override
	if !self.IsExistKey() {
		ChatSessionMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDIS_CHAT_SESSION_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		ChatSessionMgr.Delete(self.Id) // 释放对象
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisChatSessionObj) InitRedis() { //override
	obj := self.QueryChatSession(self.Id)
	if obj == nil {
		ChatSessionMgr.Delete(self.Id) // 释放对象
		return
	}
	self.SetRedisChatSession(obj)
}
func (self *RedisChatSessionObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisChatSession()
	return data
}

func (self *RedisChatSessionObj) SaveOtherData() { //override
}
func (self *RedisChatSessionObj) QueryChatSession(id string) *share_message.ChatSession {
	var data *share_message.ChatSession
	col, closeFun := self.DB.GetC(self.DBName, self.TBName)
	defer closeFun()
	//过滤过期的道具
	err := col.Find(bson.M{"_id": id}).One(&data)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return data
}
func (self *RedisChatSessionObj) SetRedisChatSession(session *share_message.ChatSession) {
	ChatSessionMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	//会话基础信息
	ex := &ChatSessionEx{}
	StructToOtherStruct(session, ex)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), ex)
	easygo.PanicError(err)
	//其他玩家数据设置到redis
	self.SetSessionOtherData(session, false)
}
func (self *RedisChatSessionObj) SetSessionOtherData(session *share_message.ChatSession, save ...bool) {
	self.SetStringValueToRedis(REDIS_CHAT_SESSION_PLAYERS, session.GetPlayerIds(), save...)
	self.SetStringValueToRedis(REDIS_CHAT_SESSION_READINFO, session.GetReadInfo(), save...)
}
func (self *RedisChatSessionObj) GetSessionOtherData(session *share_message.ChatSession) {
	self.GetStringValueToRedis(REDIS_CHAT_SESSION_PLAYERS, &session.PlayerIds)
	self.GetStringValueToRedis(REDIS_CHAT_SESSION_READINFO, &session.ReadInfo)
}
func (self *RedisChatSessionObj) GetRedisChatSession() *share_message.ChatSession {
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	var ex ChatSessionEx
	err = redis.ScanStruct(value, &ex)
	easygo.PanicError(err)
	var session *share_message.ChatSession
	StructToOtherStruct(ex, &session)
	//其他数据
	self.GetSessionOtherData(session)
	return session
}

//获取会话成员
func (self *RedisChatSessionObj) GetRedisChatSessionPlayers() []int64 {
	playerIds := make([]int64, 0)
	self.GetStringValueToRedis(REDIS_CHAT_SESSION_PLAYERS, &playerIds)
	return playerIds
}

//设置会话成员
func (self *RedisChatSessionObj) SetRedisChatSessionPlayers(ids []int64) {
	self.SetStringValueToRedis(REDIS_CHAT_SESSION_PLAYERS, ids)
}

//增加会话成员
func (self *RedisChatSessionObj) AddRedisChatSessionPlayers(ids []int64) {
	newIds := self.GetRedisChatSessionPlayers()
	newIds = append(newIds, ids...)
	self.SetStringValueToRedis(REDIS_CHAT_SESSION_PLAYERS, newIds)
}

//删除会话成员
func (self *RedisChatSessionObj) DelRedisChatSessionPlayers(ids []int64) {
	newIds := self.GetRedisChatSessionPlayers()
	for _, id := range ids {
		newIds = easygo.Del(newIds, id).([]int64)
	}
	self.SetStringValueToRedis(REDIS_CHAT_SESSION_PLAYERS, newIds)
}

//获取下一个日志id
func (self *RedisChatSessionObj) GetNextMaxLogId() int64 {
	logId := self.IncrOneValue("MaxLogId", 1)
	return logId
}

//获取当前日志id
func (self *RedisChatSessionObj) GetMaxLogId() int64 {
	var val int64
	self.GetOneValue("MaxLogId", &val)
	return val
}

//设置maxLogId
func (self *RedisChatSessionObj) SetMaxLogId(id int64) {
	self.SetOneValue("MaxLogId", id)
}

//会话类型
func (self *RedisChatSessionObj) GetType() int32 {
	var val int32
	self.GetOneValue("Type", &val)
	return val
}

//头像
func (self *RedisChatSessionObj) GetSessionHeadUrl() string {
	var val string
	self.GetOneValue("SessionHeadUrl", &val)
	return val
}

//话题
func (self *RedisChatSessionObj) GetTopic() string {
	var val string
	self.GetOneValue("Topic", &val)
	return val
}
func (self *RedisChatSessionObj) SetSessionHeadUrl(headUrl string) {
	self.SetOneValue("SessionHeadUrl", headUrl)
}

//名称
func (self *RedisChatSessionObj) GetSessionName() string {
	var val string
	self.GetOneValue("SessionName", &val)
	return val
}
func (self *RedisChatSessionObj) SetSessionName(name string) {
	self.SetOneValue("SessionName", name)
}

//群名
func (self *RedisChatSessionObj) GetTeamName() string {
	var val string
	self.GetOneValue("TeamName", &val)
	return val
}
func (self *RedisChatSessionObj) SetTeamName(name string) {
	self.GetOneValue("TeamName", name)
	self.SetSessionName(name)
}

//已读记录信息
func (self *RedisChatSessionObj) GetReadInfo() []*share_message.ReadLogInfo {
	val := make([]*share_message.ReadLogInfo, 0)
	self.GetStringValueToRedis(REDIS_CHAT_SESSION_READINFO, &val)
	return val
}

//获取已读id
func (self *RedisChatSessionObj) GetPersonalReadId(pid int64) int64 {
	info := self.GetReadInfo()
	for _, p := range info {
		if p.GetPlayerId() == pid {
			return p.GetReadId()
		}
	}
	return 0
}

//检查是否能跟陌生人说话
func (self *RedisChatSessionObj) CheckCanStrangerTalk(pid int64, cn int32) int32 {
	t := time.Now().Unix()
	info := self.GetReadInfo()
	num := int32(0)
	for _, p := range info {
		if p.GetPlayerId() == pid {
			num = p.GetTalkNum()
			if easygo.CheckTheSameDay(t, p.GetLastTalkTime()) {
				if p.GetTalkNum() >= cn {
					return cn
				} else {
					p.TalkNum = easygo.NewInt32(p.GetTalkNum() + 1)
					p.LastTalkTime = easygo.NewInt64(t)
				}
			} else {
				//不同天重置时间
				logs.Error("不同天")
				if p.LastTalkTime == nil || p.TalkNum == nil || p.GetTalkNum() > 0 {
					logs.Error("重置不同天")
					p.TalkNum = easygo.NewInt32(1)
					p.LastTalkTime = easygo.NewInt64(t)
				}
			}
			break
		}
	}
	//存储新的
	self.SetStringValueToRedis(REDIS_CHAT_SESSION_READINFO, info)
	return num
}

//重置与陌生人对话的次数
func (self *RedisChatSessionObj) ResetStrangerTalkNum(pid int64) {
	t := time.Now().Unix()
	info := self.GetReadInfo()
	for _, p := range info {
		if p.GetPlayerId() == pid {
			p.TalkNum = easygo.NewInt32(0)
			p.LastTalkTime = easygo.NewInt64(t)
			//存储新的
			self.SetStringValueToRedis(REDIS_CHAT_SESSION_READINFO, info)
			break
		}
	}
}

//设置已读信息
func (self *RedisChatSessionObj) SetReadInfo(pid int64, logId int64) {
	readInfo := self.GetReadInfo()
	for _, info := range readInfo {
		if info.GetPlayerId() == pid {
			if info.GetReadId() > logId {
				return
			}
			info.ReadId = easygo.NewInt64(logId)
			break
		}
	}
	self.SetStringValueToRedis(REDIS_CHAT_SESSION_READINFO, readInfo)
}

//获取对客户端数据:isAll=是否获取全部会话详细信息
func (self *RedisChatSessionObj) GetChatSessionDataForClient(pid int64, all ...bool) *client_hall.SessionData {
	isAll := append(all, false)[0]
	//if pid == 0 {
	//	logs.Error("玩家id怎么会为0呢？")
	//	return nil
	//}
	if self.GetType() == CHAT_TYPE_PRIVATE {
		return self.GetChatSessionDataForPersonal(pid, isAll)
	} else {
		return self.GetChatSessionDataForTeam(pid, isAll)
	}
}

//检测用户在此会话是否有数据变更
func (self *RedisChatSessionObj) CheckSessionForClient(pid int64) bool {
	if pid == 0 {
		logs.Error("玩家id怎么会为0呢？")
		return false
	}
	if self.GetType() == CHAT_TYPE_PRIVATE {
		return self.CheckPersonalSession(pid)
	} else {
		return self.CheckTeamSession(pid)
	}
}
func (self *RedisChatSessionObj) CheckPersonalSession(pid int64) bool {
	var readId int64
	for _, r := range self.GetReadInfo() {
		if r.GetPlayerId() == pid {
			readId = r.GetReadId()
			break
		}
	}
	return readId != self.GetMaxLogId()
}
func (self *RedisChatSessionObj) CheckTeamSession(pid int64) bool {
	teamId := easygo.AtoInt64(self.Id)
	memObj := GetRedisTeamPersonalObj(teamId)
	member := memObj.GetTeamMember(pid)
	return member.GetReadId() != self.GetMaxLogId()
}

//获取私聊会话数据:过滤打招呼卡片
func (self *RedisChatSessionObj) GetChatSessionDataForPersonal(pid int64, isAll bool) *client_hall.SessionData {
	chatObj := GetRedisPersonalChatLogObj(self.Id)
	log := chatObj.GetPersonalChatLog(self.GetMaxLogId())
	var readId, otherId int64
	sayHiLog := GetOneSayHiLog(self.Id, TALK_CONTENT_SAY_HI)
	for _, r := range self.GetReadInfo() {
		if r.GetPlayerId() == pid {
			readId = r.GetReadId()
		} else {
			otherId = r.GetPlayerId()
		}
	}
	resp := &client_hall.SessionData{
		Id:       easygo.NewString(self.Id),
		Type:     easygo.NewInt32(self.GetType()),
		MaxLogId: easygo.NewInt64(self.GetMaxLogId()),
		TargetId: easygo.NewInt64(otherId),
	}
	fq := GetFriendBase(pid)
	f := fq.GetFriend(otherId)
	//会话昵称设置
	base := GetRedisPlayerBase(otherId)
	name := f.GetReName()
	if base != nil {
		resp.SessionHeadUrl = easygo.NewString(base.GetHeadIcon())
		if name == "" {
			name = base.GetNickName()
		}
	}
	resp.SessionName = easygo.NewString(name)
	if log != nil && log.GetTalker() == otherId {
		log.TalkerNickName = easygo.NewString(name)
	}
	resp.PersonalChat = log
	if isAll {
		resp.ReadId = easygo.NewInt64(readId)
		num := self.GetMaxLogId() - readId
		if num < 0 {
			num = 0
		}
		if sayHiLog != nil && sayHiLog.GetTalkLogId() >= readId {
			num -= 1
		}
		resp.NewNum = easygo.NewInt64(num)
		if f != nil {

			setting := f.GetSetting()
			//resp.SessionName = f.ReName
			resp.IsTopChat = setting.IsTopChat
			resp.IsNoDisturb = setting.IsNoDisturb
			resp.IsAfterReadClear = setting.IsAfterReadClear
			resp.IsScreenShotNotify = setting.IsScreenShotNotify
			//resp.WithdrawList = GetAllWithDrawLogIds(resp.GetId(), pid, readId)
		}
	}
	return resp
}

//退出会话
func (self *RedisChatSessionObj) ExitChatSession(plist []int64) {
	ids := self.GetRedisChatSessionPlayers()
	for _, pid := range plist {
		ids = easygo.Del(ids, pid).([]int64)
	}
	self.SetRedisChatSessionPlayers(ids)
}

// 获取自己能显示的日志
func (self *RedisChatSessionObj) GetOneShowLog(pid, logId, readId int64) *share_message.TeamChatLog {
	chatObj := GetRedisTeamChatLog(easygo.AtoInt64(self.Id))
	log := chatObj.GetOneShowLog(pid, readId, self.GetMaxLogId())
	return log
}

//获取群聊会话数据
func (self *RedisChatSessionObj) GetChatSessionDataForTeam(pid int64, isAll bool) *client_hall.SessionData {
	teamId := easygo.AtoInt64(self.Id)
	teamObj := GetRedisTeamObj(teamId)
	//确保消息事最新的昵称
	memObj := GetRedisTeamPersonalObj(teamId)
	member := &share_message.PersonalTeamData{}
	if pid != 0 {
		member = memObj.GetTeamMember(pid)
	}
	chatObj := GetRedisTeamChatLog(teamId)
	log := chatObj.GetTeamChatLog(self.GetMaxLogId())
	if member.GetId() == 0 {
		//不是群成员
		resp := &client_hall.SessionData{
			Id:             easygo.NewString(self.Id),
			Type:           easygo.NewInt32(self.GetType()),
			TeamChat:       log,
			MaxLogId:       easygo.NewInt64(self.GetMaxLogId()),
			SessionHeadUrl: easygo.NewString(self.GetSessionHeadUrl()),
			SessionName:    easygo.NewString(self.GetSessionName()),
			TargetId:       easygo.NewInt64(self.Id),
			TeamName:       easygo.NewString(self.GetTeamName()),
			PlayerNum:      easygo.NewInt32(len(self.GetRedisChatSessionPlayers())),
			NewNum:         easygo.NewInt64(0),
			Topic:          easygo.NewString(self.GetTopic()),
		}
		if pid != 0 && isAll {
			self.ExitChatSession([]int64{pid})
			return nil
		}
		return resp
	}
	if log != nil {
		if member.GetNickName() == "" {
			talker := GetRedisPlayerBase(log.GetTalker())
			if talker != nil {
				log.TalkerName = easygo.NewString(talker.GetNickName())
			}
		} else {
			log.TalkerName = easygo.NewString(member.GetNickName())
		}
	}
	num := self.GetMaxLogId() - member.GetReadId()
	//noNum := chatObj.GetNoShowNum(pid)
	if num < 0 {
		num = 0
	}
	resp := &client_hall.SessionData{
		Id:             easygo.NewString(self.Id),
		Type:           easygo.NewInt32(self.GetType()),
		TeamChat:       log,
		MaxLogId:       easygo.NewInt64(self.GetMaxLogId()),
		SessionHeadUrl: easygo.NewString(self.GetSessionHeadUrl()),
		SessionName:    easygo.NewString(self.GetSessionName()),
		TargetId:       easygo.NewInt64(self.Id),
		TeamName:       easygo.NewString(teamObj.GetTeamName()),
		PlayerNum:      easygo.NewInt32(len(self.GetRedisChatSessionPlayers())),
		Topic:          easygo.NewString(self.GetTopic()),
		//NewNum:         easygo.NewInt64(num),
	}
	if pid != 0 && isAll {
		if member.GetPlayerId() == 0 {
			//玩家不在本会话中
			self.ExitChatSession([]int64{pid})
			return nil
		}
		optLog := chatObj.GetTeamOptLog(pid, member.GetReadId()+1, self.GetMaxLogId())
		unShowNum := int64(0)
		pos := memObj.GetTeamPosition(pid)
		for _, lo := range optLog {
			message := lo.GetTeamMessage()
			//过滤不显示的系统消息
			if message != nil && pos > message.GetShowPos() {
				unShowNum += 1
			}
			//过滤不显示的领取红包记录
			if lo.GetType() == TALK_CONTENT_REDPACKET_LOG && !easygo.Contain(lo.GetPlayerIds(), pid) {
				unShowNum += 1
			}
			if lo.GetType() != TALK_CONTENT_SYSTEM && lo.GetTalker() == pid {
				unShowNum += 1
			}
		}
		resp.NewNum = easygo.NewInt64(num - unShowNum)
		isNotice := chatObj.GetTeamNoticeMessage(pid, member.GetReadId())
		resp.ReadId = easygo.NewInt64(member.GetReadId())
		resp.IsNotice = easygo.NewBool(isNotice)
		//resp.TeamOptLog = optLog
		setting := member.GetSetting()
		if setting != nil {
			resp.IsTopChat = setting.IsTopChat
			resp.IsNoDisturb = setting.IsNoDisturb
		}
		//resp.WithdrawList = chatObj.GetAllTeamWithDrawLogIds(pid, member.GetReadId())
	}
	return resp
}

//前端删除聊天记录
func (self *RedisChatSessionObj) DeleteMessage(logId int64) {
	if self.GetType() == CHAT_TYPE_PRIVATE {
		chatObj := GetRedisPersonalChatLogObj(self.Id)
		log := chatObj.GetPersonalChatLog(logId)
		if log != nil {
			log.Status = easygo.NewInt32(1) // 状态改1，表示删除的消息
			chatObj.UpdatePersonalChatLog(log)
		}
	} else {
		teamId := easygo.AtoInt64(self.Id)
		chatObj := GetRedisTeamChatLog(teamId)
		log := chatObj.GetTeamChatLog(logId)
		if log != nil {
			log.Status = easygo.NewInt32(1) // 状态改1，表示删除的消息
			chatObj.AddTeamChatLog(log)
		}
	}
}

//对外接口
func GetRedisChatSessionObj(id string, data ...*share_message.ChatSession) *RedisChatSessionObj {
	return ChatSessionMgr.GetRedisChatSessionObj(id, data...)
}

//创建一个新的会话
func CreateRedisChatSessionObj(id int64, data *share_message.Chat) *RedisChatSessionObj {
	session := &share_message.ChatSession{
		Id:       easygo.NewString(data.GetSessionId()),
		Type:     easygo.NewInt32(data.GetChatType()),
		MaxLogId: easygo.NewInt64(0),
	}
	if data.GetChatType() == CHAT_TYPE_PRIVATE {
		//私聊会话
		playerIds := []int64{id, data.GetTargetId()}
		session.PlayerIds = playerIds
		readInfo := make([]*share_message.ReadLogInfo, 0)
		for _, pid := range playerIds {
			r := &share_message.ReadLogInfo{
				PlayerId: easygo.NewInt64(pid),
				ReadId:   easygo.NewInt64(0),
				TalkNum:  easygo.NewInt32(0),
			}
			readInfo = append(readInfo, r)
		}
		session.ReadInfo = readInfo

	} else if data.GetChatType() == CHAT_TYPE_TEAM {
		//群聊会话
		teamId := data.GetTargetId()
		teamObj := GetRedisTeamObj(teamId)
		members := teamObj.GetTeamMemberList()
		session.SessionName = easygo.NewString(teamObj.GetTeamName())
		session.SessionHeadUrl = easygo.NewString(teamObj.GetHeadUrl())
		session.PlayerIds = members
	}

	return ChatSessionMgr.GetRedisChatSessionObj(session.GetId(), session)
}
