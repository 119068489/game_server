package for_game

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"sort"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

/*
群聊天内存数据管理
*/

const REDIS_TEAMCHAT_EXIST_LIST = "team_chat_exist_list"
const REDIS_TEAMCHAT_EXIST_TIME = 1000 * 600 //毫秒，key值存在时间

type RedisTeamChatLogObj struct {
	Id int64 //群id
	RedisBase
}

//写入redis
func NewRedisTeamChatLog(id int64, data ...[]*share_message.TeamChatLog) *RedisTeamChatLogObj {
	p := &RedisTeamChatLogObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	p.Init(obj)
	return p
}
func (self *RedisTeamChatLogObj) Init(obj []*share_message.TeamChatLog) *RedisTeamChatLogObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG)
	self.Sid = TeamChatLogMgr.GetSid()
	//redis已经存在key了，但管理器已销毁
	TeamChatLogMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//if self.IsExistKey() {
	//	self.SetSaveStatus(true)
	//}
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	//logs.Info("初始化新的TeamChatLog管理器:", self.Id)
	return self
}
func (self *RedisTeamChatLogObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisTeamChatLogObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_TEAM_CHAT_LOG, self.Id)
}

//重写保存方法
func (self *RedisTeamChatLogObj) SaveToMongo() {
	logsData := self.GetRedisTeamChatLog()
	//logs.Info("data:", logsData)
	var keys []string
	var saveData []interface{}
	for _, log := range logsData {
		keys = append(keys, easygo.AnytoA(log.GetTeamLogId()))
		saveData = append(saveData, bson.M{"_id": log.GetLogId()}, log)
	}
	if len(saveData) > 0 {
		//logs.Info("RedisTeamChatLogObj redis 存储 ：", saveData)
		UpsertAll(easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG, saveData)
	}
	if len(keys) > 0 {
		b, err := easygo.RedisMgr.GetC().Hdel(self.GetKeyId(), keys...)
		if !b {
			logs.Error("删除redis key:", err)
		}
	}
	self.SetSaveStatus(false)
}

//定时更新数据
func (self *RedisTeamChatLogObj) UpdateData() { //override
	if !self.IsExistKey() {
		TeamChatLogMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDIS_TEAMCHAT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		TeamChatLogMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}

func (self *RedisTeamChatLogObj) InitRedis() { //override
	//if self.IsExistKey() {
	//	return
	//}
	//obj := self.QueryTeamChatLog()
	//if obj == nil {
	//	return
	//}
	////重置过期时间
	//self.CreateTime = GetMillSecond()
	//self.SetRedisTeamChatLog(obj)
	////重新激活定时器
	//easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisTeamChatLogObj) GetRedisSaveData() interface{} { //override

	return nil
}
func (self *RedisTeamChatLogObj) SaveOtherData() { //override
	//保存表情数据
}
func (self *RedisTeamChatLogObj) QueryTeamChatLog() []*share_message.TeamChatLog {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG)
	defer closeFun()
	var lst []*share_message.TeamChatLog
	err1 := col.Find(bson.M{"TeamId": self.Id}).Sort("-_id").Limit(100).All(&lst)
	easygo.PanicError(err1)
	return lst
}
func (self *RedisTeamChatLogObj) QueryTeamOneChatLog(logId int64) *share_message.TeamChatLog {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG)
	defer closeFun()
	var log *share_message.TeamChatLog
	err1 := col.Find(bson.M{"TeamId": self.Id, "TeamLogId": logId}).Sort("-_id").One(&log)
	if err1 != nil && err1.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err1)
	}
	return log
}

//获取当前redis记录
func (self *RedisTeamChatLogObj) GetRedisTeamChatLog() []*share_message.TeamChatLog {
	var chatList []*share_message.TeamChatLog
	if !self.IsExistKey() {
		return chatList
	}
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	for _, m := range values {
		log := &share_message.TeamChatLog{}
		_ = json.Unmarshal([]byte(m), &log)
		chatList = append(chatList, log)
	}
	return chatList
}
func (self *RedisTeamChatLogObj) AddTeamChatLog(obj *share_message.TeamChatLog) {
	s, err := json.Marshal(obj)
	easygo.PanicError(err)
	err = easygo.RedisMgr.GetC().HSet(self.GetKeyId(), easygo.AnytoA(obj.GetTeamLogId()), string(s))
	easygo.PanicError(err)
	self.SetSaveSid()
	self.SetSaveStatus(true)
}

//获取存在的logs
func (self *RedisTeamChatLogObj) GetChatLogKeys() []int64 {
	val, err := easygo.RedisMgr.GetC().HKeys(self.GetKeyId())
	easygo.PanicError(err)
	var logs []int64
	InterfersToInt64s(val, &logs)
	return logs
}

func GetRedisTeamChatLog(id int64) *RedisTeamChatLogObj {
	return TeamChatLogMgr.GetRedisTeamChatLogObj(id)
}

func AddTeamChatLog(teamId, playerId, logId int64, chat *share_message.Chat, notice *share_message.NoticeMsg, isRead bool, msg ...*share_message.TeamMessage) int64 {
	Id := NextId(TABLE_TEAM_CHAT_LOG)
	chatObj := TeamChatLogMgr.GetRedisTeamChatLogObj(teamId)
	//logId := chatObj.GetTeamMaxChatLogId() + 1

	log := &share_message.TeamChatLog{
		LogId:         easygo.NewInt64(Id),
		TeamId:        easygo.NewInt64(teamId),
		Talker:        easygo.NewInt64(playerId),
		Content:       easygo.NewString(chat.GetContent()),
		Time:          easygo.NewInt64(chat.GetTime()),
		IsSave:        easygo.NewBool(true),
		Type:          easygo.NewInt32(chat.GetContentType()),
		Cite:          easygo.NewString(chat.GetCite()),
		TeamLogId:     easygo.NewInt64(logId),
		TalkerName:    easygo.NewString(chat.GetSourceName()),
		TalkerHeadUrl: easygo.NewString(chat.GetSourceHeadIcon()),
		QPId:          easygo.NewInt64(chat.GetQPId()),
		SessionId:     easygo.NewString(chat.GetSessionId()),
		Status:        easygo.NewInt32(TALK_STATUS_NORMAL),
		Mark:          easygo.NewString(chat.GetMark()),
		IsWelComeWord: easygo.NewBool(chat.GetIsWelcome()),
	}
	if chat.GetContentType() == TALK_CONTENT_REDPACKET_LOG {
		//领取红包记录
		log.PlayerIds = chat.PlayIds
	}
	if notice != nil {
		log.NoticeInfo = notice
	}
	//SetTeamLastTalkTime(teamId, chat.GetTime())
	if chat.GetContentType() == TALK_CONTENT_SYSTEM || chat.GetContentType() == TALK_CONTENT_GROUPNOTICE {
		m := append(msg, &share_message.TeamMessage{})[0]
		log.TeamMessage = m
	}
	chatObj.AddTeamChatLog(log)
	if chat.GetContentType() != 0 && isRead {
		memberObj := GetRedisTeamPersonalObj(teamId)
		memberObj.ReadTeamChatLog(playerId, logId)
	}
	return logId
}

//获取指定聊天日志
func (self *RedisTeamChatLogObj) GetTeamChatLog(logId int64) *share_message.TeamChatLog {
	var log *share_message.TeamChatLog
	b, err := easygo.RedisMgr.GetC().HExists(self.GetKeyId(), easygo.AnytoA(logId))
	if !b {
		log = self.QueryTeamOneChatLog(logId)
		return log
	}
	val, err := easygo.RedisMgr.GetC().HGet(self.GetKeyId(), easygo.AnytoA(logId))
	easygo.PanicError(err)
	err = json.Unmarshal(val, &log)
	easygo.PanicError(err)
	return log
}

//获取自己能显示的最近一条消息
func (self *RedisTeamChatLogObj) GetOneShowLog(pid, logId, maxLogId int64) *share_message.TeamChatLog {
	self.SaveToMongo()
	pos := GetTeamPlayerPos(self.Id, pid)
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG)
	defer closeFun()
	//查找我能显示的最新一条群聊天:
	var log *share_message.TeamChatLog
	q1 := bson.M{"Type": TALK_CONTENT_SYSTEM, "TeamMessage.ShowPos": bson.M{"$gte": pos}}
	q2 := bson.M{"Type": TALK_CONTENT_REDPACKET_LOG, "PlayerIds": pid}
	q3 := bson.M{"Type": bson.M{"$nin": []int32{TALK_CONTENT_REDPACKET_LOG, TALK_CONTENT_SYSTEM}}}
	q := bson.M{"TeamId": self.Id, "$or": []bson.M{q1, q2, q3}}
	err := col.Find(q).Sort("-_id").One(&log)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	logs.Info("返回日志:", log)
	return log

}

//获取系统日志和红包领取日志
func (self *RedisTeamChatLogObj) GetTeamOptLog(pid, reaId, maxId int64) []*share_message.TeamChatLog {
	logs := self.GetSomeTeamChatLog(reaId, maxId)
	newLogs := make([]*share_message.TeamChatLog, 0)
	for _, log := range logs {
		if log.GetTeamMessage() != nil {
			newLogs = append(newLogs, log)
		}
		if log.GetType() == TALK_CONTENT_REDPACKET_LOG {
			if !easygo.Contain(log.GetPlayerIds(), pid) {
				newLogs = append(newLogs, log)
			}
		}
		if pid == log.GetTalker() && !log.GetIsWelComeWord() {
			newLogs = append(newLogs, log)
		}
	}
	return newLogs
}

//撤销聊天记录
func (self *RedisTeamChatLogObj) TeamWithdrawChatLog(playerId, logId int64) bool {
	log := self.GetTeamChatLog(logId)
	if log.GetLogId() == 0 {
		panic("log对象为空")
	}
	t := GetMillSecond()
	if t-log.GetTime() > 5*60*1000 {
		return false
	}
	memberObj := GetRedisTeamPersonalObj(self.Id)
	name := memberObj.GetTeamMemberReName(playerId)
	if name == "" {
		player := GetRedisPlayerBase(playerId)
		name = player.GetNickName()
	}

	content := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s撤回了一条消息", name)))
	log.Type = easygo.NewInt32(TALK_CONTENT_WITHDRAW)
	log.Content = easygo.NewString(content)
	log.Time = easygo.NewInt64(t)
	self.AddTeamChatLog(log)
	return true
}

func (self *RedisTeamChatLogObj) GetSomeTeamChatLog(startId, overId int64) []*share_message.TeamChatLog {
	var unlogIds []int64 //找出不在redis内存中的聊天日志 去mongo里面找
	logs := make([]*share_message.TeamChatLog, 0)
	idList := easygo.GetInt64List(startId, overId)
	logsRedis := self.GetRedisTeamChatLog()
	ids := self.GetChatLogKeys()
	//redis内存中没有日志
	if len(logsRedis) == 0 {
		//去mongo查询
		logs = self.GetTeamChatLogForMongoDB(idList)
		return logs
	}
	//redis中存在，
	for _, log := range logsRedis {
		if easygo.Contain(idList, log.GetTeamLogId()) {
			logs = append(logs, log)
		}
	}
	//找出redis中不存在
	for _, id := range idList {
		if !easygo.Contain(ids, id) {
			unlogIds = append(unlogIds, id)
		}
	}
	//从数据库查询
	if len(unlogIds) > 0 {
		unlogs := self.GetTeamChatLogForMongoDB(unlogIds)
		logs = append(logs, unlogs...)
	}
	return logs
}

func (self *RedisTeamChatLogObj) GetTeamChatLogForMongoDB(logs []int64) []*share_message.TeamChatLog {
	var lst []*share_message.TeamChatLog
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG)
	defer closeFun()
	err := col.Find(bson.M{"TeamId": self.Id, "TeamLogId": bson.M{"$in": logs}}).All(&lst)
	easygo.PanicError(err)
	return lst
}

func (self *RedisTeamChatLogObj) GetTeamChatLogListForLogId(startId, overId int64, count int32) *client_hall.ReturnChatInfo {
	lst := self.GetSomeTeamChatLog(startId, overId)
	for _, v := range lst {
		if player := GetRedisPlayerBase(v.GetTalker()); player != nil {
			v.Types = easygo.NewInt32(player.GetTypes())
		}

	}
	msg := &client_hall.ReturnChatInfo{
		ChatList: lst,
		Count:    easygo.NewInt32(count),
	}
	return msg
}

func (self *RedisTeamChatLogObj) GetTeamUnReadChatLogs(maxId int64, playerObj *RedisPlayerBaseObj) *share_message.TeamChatInfo { //获取未读群聊天记录
	msg := &share_message.TeamChatInfo{}
	//maxId := self.GetTeamMaxChatLogId()
	if maxId != 0 {
		memberObj := GetRedisTeamPersonalObj(self.Id)
		readId := memberObj.GetTeamReadChatLogId(playerObj.Id)
		num := maxId - readId
		var firstList []*share_message.TeamChatLog
		if num > 100 {
			firstList = self.GetSomeTeamChatLog(int64(readId+1), int64(readId+51))
			overList := self.GetSomeTeamChatLog(int64(maxId-50), int64(maxId))
			msg.LastChat = overList
			msg.Count = easygo.NewInt32(num - 100)
		} else {
			if num != 0 {
				firstList = self.GetSomeTeamChatLog(int64(readId+1), int64(maxId))
			}
		}
		msg.FirstChat = firstList
		msg.ReLogIds = self.GetTeamWithDrawLogIds(playerObj, readId)
		msg.IsNoticeMessage = easygo.NewBool(self.GetTeamNoticeMessage(playerObj.Id, readId))
	}
	return msg
}
func (self *RedisTeamChatLogObj) GetTeamSessionChatLogs(start, end int64, sessionId string) []*share_message.TeamChatLog {
	//先保存数据
	self.SaveToMongo()
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG)
	defer closeFun()
	var chatInfo []*share_message.TeamChatLog
	err := col.Find(bson.M{"SessionId": sessionId, "TeamLogId": bson.M{"$gte": start, "$lte": end}}).Sort("-TeamLogId").All(&chatInfo)
	easygo.PanicError(err)
	return chatInfo
}

//获取撤回消息id
func (self *RedisTeamChatLogObj) GetTeamWithDrawLogIds(player *RedisPlayerBaseObj, readId int64) []*share_message.WithDrawInfo {
	var lst []*share_message.WithDrawInfo
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	for _, s := range values {
		var obj *share_message.TeamChatLog
		json.Unmarshal([]byte(s), &obj)
		if obj.GetTeamLogId() < readId {
			continue
		}
		if obj.GetType() == TALK_CONTENT_WITHDRAW && obj.GetTime() > player.GetLastLogOutTime() {
			msg := &share_message.WithDrawInfo{
				LogId:    easygo.NewInt64(obj.GetTeamLogId()),
				PlayerId: easygo.NewInt64(obj.GetTalker()),
			}
			lst = append(lst, msg)
		}
	}
	sort.Slice(lst, func(i, j int) bool {
		return lst[i].GetLogId() > lst[j].GetLogId() // 降序
		//return newlst[i].GetLogId() < newlst[j].GetLogId() // 升序
	})

	return lst
}

//获取所有撤回消息
func (self *RedisTeamChatLogObj) GetAllTeamWithDrawLogIds(pid, readId int64) []int64 {
	player := GetRedisPlayerBase(pid)
	ids := self.GetTeamWithDrawLogIdsRedis(player, readId)
	ids1 := self.GetTeamWithDrawLogIdsMongo(player, readId)
	ids = append(ids, ids1...)
	return ids
}

//获取撤回消息的id
func (self *RedisTeamChatLogObj) GetTeamWithDrawLogIdsRedis(player *RedisPlayerBaseObj, readId int64) []int64 {
	data := make([]int64, 0)
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(TABLE_TEAM_CHAT_LOG, self.Id)))
	easygo.PanicError(err)
	for _, s := range values {
		var obj *share_message.TeamChatLog
		json.Unmarshal([]byte(s), &obj)
		if obj.GetTeamLogId() < readId {
			//5分钟以前的
			if obj.GetType() == TALK_CONTENT_WITHDRAW && obj.GetTime() >= player.GetLastLogOutTime()-5*60*1000 {
				data = append(data, obj.GetTeamLogId())
			}
		} else {
			if obj.GetType() == TALK_CONTENT_WITHDRAW && obj.GetTime() > player.GetLastLogOutTime() {
				data = append(data, obj.GetTeamLogId())
			}
		}
	}
	return data
}

//
func (self *RedisTeamChatLogObj) GetTeamWithDrawLogIdsMongo(player *RedisPlayerBaseObj, readId int64) []int64 {
	data := make([]int64, 0)
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()
	var chatInfo []*share_message.TeamChatLog
	err := col.Find(bson.M{"TeamId": self.Id, "Type": TALK_CONTENT_WITHDRAW, "TeamLogId": bson.M{"$lte": readId}, "Time": bson.M{"$gte": player.GetLastLogOutTime()}}).Sort("-TeamLogId").All(&chatInfo)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	for _, p := range chatInfo {
		data = append(data, p.GetTeamLogId())
	}
	return data
}

func (self *RedisTeamChatLogObj) GetTeamNoticeMessage(playerId, readId int64) bool {
	var b bool
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	for _, s := range values {
		var obj *share_message.TeamChatLog
		json.Unmarshal([]byte(s), &obj)
		if obj.GetTeamLogId() < readId {
			continue
		}
		notice := obj.GetNoticeInfo()
		if notice == nil {
			continue
		}
		if notice.GetIsAll() {
			return true
		}
		playerIds := notice.GetPlayerId()
		if util.Int64InSlice(playerId, playerIds) {
			return true
		}
	}
	return b
}

//批量保存需要存储的数据
func SaveRedisTeamChatLogToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_TEAM_CHAT_LOG, &ids)
	logList := make([]*share_message.TeamChatLog, 0)
	for _, id := range ids {
		obj := GetRedisTeamChatLog(id)
		if obj != nil {
			logs := obj.GetRedisTeamChatLog()
			logList = append(logList, logs...)
			obj.SetSaveStatus(false)
		}
	}
	if len(logList) > 0 {
		saveData := make([]interface{}, 0)
		for _, log := range logList {
			saveData = append(saveData, bson.M{"_id": log.GetLogId()}, log)
		}
		UpsertAll(easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG, saveData)
	}
}
