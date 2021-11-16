package for_game

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

const REDIS_PERSONALCHAT_EXIST_LIST = "personal_chat_exist_list"
const REDIS_PERSONALCHAT_EXIST_TIME = 1000 * 600

type RedisPersonalChatLogObj struct {
	Id string
	RedisBase
}

//写入redis
func NewRedisPersonalChatLog(id string, data ...[]*share_message.PersonalChatLog) *RedisPersonalChatLogObj {
	p := &RedisPersonalChatLogObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	p.Init(obj)
	return p
}

func (self *RedisPersonalChatLogObj) Init(obj []*share_message.PersonalChatLog) *RedisPersonalChatLogObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG)
	self.Sid = PersonalChatLogMgr.GetSid()
	//增加到管理器
	PersonalChatLogMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//if self.IsExistKey() {
	//	self.SetSaveStatus(true)
	//}
	//self.IsCheck = false
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	return self
}
func (self *RedisPersonalChatLogObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisPersonalChatLogObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_PERSONAL_CHAT_LOG, self.Id)
}

//重写保存方法
func (self *RedisPersonalChatLogObj) SaveToMongo() { //override
	logsData := self.GetRedisPersonalChatLog()
	//logs.Info("redis 存储 ：", logsData)
	var keys []string
	if len(logsData) > 0 {
		var saveData []interface{}
		for _, log := range logsData {
			keys = append(keys, easygo.AnytoA(log.GetTalkLogId()))
			saveData = append(saveData, bson.M{"_id": log.GetLogId()}, log)
		}
		if len(saveData) > 0 {
			UpsertAll(easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG, saveData)
		}
	}
	if len(keys) > 0 {
		b, err := easygo.RedisMgr.GetC().Hdel(self.GetKeyId(), keys...)
		if !b {
			logs.Error("删除keys失败,redis不存在：", keys, err)
		}
	}
	self.SetSaveStatus(false)
}

//定时更新数据
func (self *RedisPersonalChatLogObj) UpdateData() { //override
	if !self.IsExistKey() {
		PersonalChatLogMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDIS_PERSONALCHAT_EXIST_TIME { //单位：毫秒
		logs.Info("释放对象")
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		PersonalChatLogMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisPersonalChatLogObj) InitRedis() { //override
	//if self.IsExistKey() {
	//	return
	//}
	//PersonalChatLogMgr.Store(self.Id, self)
	//self.AddToExistList(self.Id)
	////重置过期时间
	//self.CreateTime = GetMillSecond()
	////重新激活定时器
	//easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisPersonalChatLogObj) GetRedisSaveData() interface{} { //override
	return nil
}

func (self *RedisPersonalChatLogObj) SaveOtherData() { //override
}

//获取查询一条记录
func (self *RedisPersonalChatLogObj) QueryPersonalOneChatLog(id int64) *share_message.PersonalChatLog {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()
	var lst *share_message.PersonalChatLog
	err1 := col.Find(bson.M{"SessionId": self.Id, "TalkLogId": id}).One(&lst)
	if err1 != nil && err1.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err1)
	}
	return lst
}

//查询全部聊天记录
func (self *RedisPersonalChatLogObj) QueryPersonalAllChatLog() []*share_message.PersonalChatLog {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()
	var lst []*share_message.PersonalChatLog
	err1 := col.Find(bson.M{"SessionId": self.Id}).Sort("-_id").Limit(100).All(&lst)
	easygo.PanicError(err1)
	return lst
}

//获取当前redis记录
func (self *RedisPersonalChatLogObj) GetRedisPersonalChatLog() []*share_message.PersonalChatLog {
	var chatList []*share_message.PersonalChatLog
	if !self.IsExistKey() {
		return chatList
	}
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	for _, m := range values {
		log := &share_message.PersonalChatLog{}
		_ = json.Unmarshal([]byte(m), &log)
		chatList = append(chatList, log)
	}
	return chatList
}

//设置sayHi为删除状态
func (self *RedisPersonalChatLogObj) DelSayHiLog() *share_message.PersonalChatLog {
	logsX := self.GetRedisPersonalChatLog()
	for _, log := range logsX {
		if log.GetSessionId() == self.Id && log.GetType() == TALK_CONTENT_SAY_HI && log.GetStatus() != TALK_STATUS_DELETE {
			log.Status = easygo.NewInt32(TALK_STATUS_DELETE)
			self.UpdatePersonalChatLog(log)
			return log
		}
	}
	//不在redis，去mongo处理
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()
	var log *share_message.PersonalChatLog
	err := col.Find(bson.M{"SessionId": self.Id, "Type": TALK_CONTENT_SAY_HI, "Status": TALK_STATUS_NORMAL}).Sort("-CreateTime").One(&log)
	if err != nil {
		logs.Error("DelSayHiLog err:", err)
	}
	log.Status = easygo.NewInt32(TALK_STATUS_DELETE)
	//更新数据
	self.UpdatePersonalChatLog(log)
	return log
}

func (self *RedisPersonalChatLogObj) AddRedisPersonalChatLog(obj *share_message.PersonalChatLog) {
	s, err := json.Marshal(obj)
	easygo.PanicError(err)
	err = easygo.RedisMgr.GetC().HSet(self.GetKeyId(), easygo.AnytoA(obj.GetTalkLogId()), string(s))
	easygo.PanicError(err)
	self.SetSaveSid()
	self.SetSaveStatus(true)
}

func (self *RedisPersonalChatLogObj) GetPersonalChatLog(id int64) *share_message.PersonalChatLog {
	var log *share_message.PersonalChatLog

	by, err := easygo.RedisMgr.GetC().HGet(self.GetKeyId(), easygo.AnytoA(id))
	if err != nil {
		log = self.QueryPersonalOneChatLog(id)
		return log
	}
	err = json.Unmarshal(by, &log)
	easygo.PanicError(err)
	return log
}
func (self *RedisPersonalChatLogObj) UpdatePersonalChatLog(log *share_message.PersonalChatLog) {
	self.AddRedisPersonalChatLog(log)
}
func AddPersonalChatLog(msg *share_message.Chat) int64 {
	return PersonalChatLogMgr.AddPersonalChatLog(msg)
}

//func ReadPersonalMessage(logIds []int64) { //发送者默认已读 所以不计字段
//	var slogs []string
//	for _, id := range logIds {
//		if util.InStringSlice(easygo.AnytoA(id), slogs) {
//			continue
//		}
//		slogs = append(slogs, easygo.AnytoA(id))
//	}
//	chatInfo := make(map[int64]string)
//	values, err := easygo.RedisMgr.GetC().HMGet(TABLE_PERSONAL_CHAT_LOG, slogs...)
//	easygo.PanicError(err)
//	var newlst []int
//	for index, m := range values {
//		if m == nil {
//			newlst = append(newlst, index)
//			continue
//		}
//		var log *share_message.PersonalChatLog
//		err1 := json.Unmarshal(m.([]byte), &log)
//		easygo.PanicError(err1)
//		log.IsRead = easygo.NewBool(true)
//		s, _ := json.Marshal(log)
//		chatInfo[log.GetLogId()] = string(s)
//	}
//	if len(chatInfo) != 0 { //redis 内存
//		err2 := easygo.RedisMgr.GetC().HMSet(TABLE_PERSONAL_CHAT_LOG, chatInfo)
//		easygo.PanicError(err2)
//	}
//
//	if len(newlst) != 0 { //已经存到数据库的
//		var ids []int64
//		for _, index := range newlst {
//			ids = append(ids, easygo.AtoInt64(slogs[index]))
//		}
//		col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG)
//		defer closeFun()
//		_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"IsRead": true}}) //未读消息
//		easygo.PanicError(err)
//	}
//}

func WithdrawMessage(pid, logId int64, sessionId string) (bool, int64) { //撤回消息
	obj := GetRedisPersonalChatLogObj(sessionId)
	log := obj.GetPersonalChatLog(logId)
	var content string
	t := GetMillSecond()
	if t-log.GetTime() > 5*60*1000 {
		return false, 0
	}
	player := GetRedisPlayerBase(pid)
	content = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`"%s"撤回了一条消息`, player.GetNickName())))
	log.Type = easygo.NewInt32(TALK_CONTENT_WITHDRAW)
	log.Content = easygo.NewString(content)
	obj.UpdatePersonalChatLog(log)
	return true, log.GetTime()
}

func GetPersonalSessionChatLog(start, end int64, sessionId string) []*share_message.PersonalChatLog {
	//先保存下数据
	obj := GetRedisPersonalChatLogObj(sessionId)
	obj.SaveToMongo()
	//SavePersonalChatToMongoDB()
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()
	var chatInfo []*share_message.PersonalChatLog
	err := col.Find(bson.M{"SessionId": sessionId, "TalkLogId": bson.M{"$gte": start, "$lte": end}}).Sort("-TalkLogId").All(&chatInfo)
	easygo.PanicError(err)
	return chatInfo
}

//获取所有撤回消息
func GetAllWithDrawLogIds(session string, pid, readId int64) []int64 {
	//获取数据库的
	player := GetRedisPlayerBase(pid)
	ids := GetWithDrawLogIdsForMongo(session, player, readId)
	ids1 := GetWithDrawLogIdsForRedis(session, player, readId)
	ids = append(ids, ids1...)
	return ids
}

//获取打招呼:
func GetAllSayHiLog(sessions []string, t int32) []*share_message.PersonalChatLog {
	logs1, newSessions := GetSayHiLogForRedis(sessions, t)
	logs2 := GetSayHiLogForMongo(newSessions, t)
	return append(logs1, logs2...)
}

//获取单条打招呼信息
func GetOneSayHiLog(session string, t int32) *share_message.PersonalChatLog {
	logs := GetAllSayHiLog([]string{session}, t)
	if len(logs) > 0 {
		return logs[0]
	}
	return nil
}

//redis里查找打招呼
func GetSayHiLogForRedis(sessions []string, t int32) ([]*share_message.PersonalChatLog, []string) {
	data := make([]*share_message.PersonalChatLog, 0)
	newSession := make([]string, 0)
	found := false
	for _, session := range sessions {
		obj := GetRedisPersonalChatLogObj(session)
		m := obj.GetRedisPersonalChatLog()
		for _, log := range m {
			if log.GetSessionId() == session && log.GetType() == t {
				data = append(data, log)
				found = true
				break
			}
		}
		if !found {
			newSession = append(newSession, session)
		} else {
			found = false
		}
	}
	return data, newSession
}

//数据库里查找打招呼
func GetSayHiLogForMongo(sessions []string, t int32) []*share_message.PersonalChatLog {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()
	var chatInfo []*share_message.PersonalChatLog
	err := col.Find(bson.M{"Type": t, "SessionId": bson.M{"$in": sessions}}).All(&chatInfo)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("数据库里查找打招呼出错：%v", err)
	}
	return chatInfo
}

//获取玩家撤回id列表:
func GetWithDrawLogIdsForMongo(session string, player *RedisPlayerBaseObj, readId int64) []int64 {
	data := make([]int64, 0)
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()
	var chatInfo []*share_message.PersonalChatLog
	err := col.Find(bson.M{"SessionId": session, "Type": TALK_CONTENT_WITHDRAW, "TalkLogId": bson.M{"$lte": readId}, "Time": bson.M{"$gte": player.GetLastLogOutTime()}}).Sort("-TalkLogId").All(&chatInfo)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	for _, p := range chatInfo {
		data = append(data, p.GetTalkLogId())
	}
	return data
}

//获取玩家撤回id列表:
func GetWithDrawLogIdsForRedis(session string, player *RedisPlayerBaseObj, readId int64) []int64 {
	obj := GetRedisPersonalChatLogObj(session)
	data := make([]int64, 0)
	m := obj.GetRedisPersonalChatLog()
	for _, log := range m {
		if log.GetSessionId() == session && log.GetType() == TALK_CONTENT_WITHDRAW {
			if log.GetTalkLogId() <= readId && log.GetTime() >= player.GetLastLogOutTime()-5*60*1000 {
				//已经读过的消息撤回，5分钟有效
				data = append(data, log.GetTalkLogId())
			} else if log.GetTalkLogId() > readId && log.GetTime() >= player.GetLastLogOutTime() {
				data = append(data, log.GetTalkLogId())
			}
		}
	}
	return data
}

func SavePersonalChatToMongoDB() {
	ids := []string{}
	GetAllRedisExistList(TABLE_PERSONAL_CHAT_LOG, &ids)
	logList := make([]*share_message.PersonalChatLog, 0)
	for _, id := range ids {
		obj := GetRedisPersonalChatLogObj(id)
		if obj != nil {
			logs := obj.GetRedisPersonalChatLog()
			logList = append(logList, logs...)
			obj.SetSaveStatus(false)
		}
	}
	if len(logList) > 0 {
		saveData := make([]interface{}, 0)
		for _, log := range logList {
			saveData = append(saveData, bson.M{"_id": log.GetLogId()}, log)
		}
		UpsertAll(easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_PERSONAL_CHAT_LOG, saveData)
	}
}
func GetRedisPersonalChatLogObj(sessionId string) *RedisPersonalChatLogObj {
	obj := PersonalChatLogMgr.GetRedisPersonalChatLogObj(sessionId)
	return obj
}
