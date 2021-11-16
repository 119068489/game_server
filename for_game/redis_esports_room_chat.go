package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

var REDIS_ESPORT_ROOM_CHAT_LOG_KEY = ESportExN("room_chat_log")

const REDIS_ESPORT_ROOM_CHAT_EXIST_TIME = 1000 * 600

type RedisESportRoomChatLogObj struct {
	Id string
	RedisBase
}

//写入redis
func NewRedisESportRoomChatLogObj(id int64) *RedisESportRoomChatLogObj {
	p := &RedisESportRoomChatLogObj{
		Id: easygo.AnytoA(id),
	}
	p.Init()
	return p
}

func (self *RedisESportRoomChatLogObj) Init() *RedisESportRoomChatLogObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_ESPORTS_ROOM_CHAT_MSG_LOG)
	self.Sid = ESportRoomChatMgr.GetSid()
	//增加到管理器
	ESportRoomChatMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	return self
}
func (self *RedisESportRoomChatLogObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisESportRoomChatLogObj) GetKeyId() string { //override
	return MakeRedisKey(REDIS_ESPORT_ROOM_CHAT_LOG_KEY, self.Id)
}

//重写保存方法
func (self *RedisESportRoomChatLogObj) SaveToMongo() { //override
	logsData := self.GetRedisESportRoomChatLog()
	//logs.Info("redis 存储 ：", logsData)
	var keys []string
	if len(logsData) > 0 {
		var saveData []interface{}
		for _, log := range logsData {
			keys = append(keys, easygo.AnytoA(log.GetId()))
			saveData = append(saveData, bson.M{"_id": log.GetId()}, log)
		}
		if len(saveData) > 0 {
			UpsertAll(easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_ESPORTS_ROOM_CHAT_MSG_LOG, saveData)
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
func (self *RedisESportRoomChatLogObj) UpdateData() { //override
	if !self.IsExistKey() {
		ESportRoomChatMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDIS_ESPORT_ROOM_CHAT_EXIST_TIME { //单位：毫秒
		logs.Info("释放对象")
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		ESportRoomChatMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisESportRoomChatLogObj) InitRedis() { //override
	//if self.IsExistKey() {
	//	return
	//}
	//ESportRoomChatMgr.Store(self.Id, self)
	//self.AddToExistList(self.Id)
	////重置过期时间
	//self.CreateTime = GetMillSecond()
	////重新激活定时器
	//easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisESportRoomChatLogObj) GetRedisSaveData() interface{} { //override
	return nil
}

func (self *RedisESportRoomChatLogObj) SaveOtherData() { //override
}

//获取查询一条记录
func (self *RedisESportRoomChatLogObj) QueryESPortsRoomOneChatLog(id int64) *share_message.TableESPortsLiveRoomMsgLog {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_ESPORTS_ROOM_CHAT_MSG_LOG)
	defer closeFun()
	var lst *share_message.TableESPortsLiveRoomMsgLog
	err1 := col.Find(bson.M{"SessionId": self.Id, "TalkLogId": id}).One(&lst)
	if err1 != nil && err1.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err1)
	}
	return lst
}

//查询全部聊天记录
func (self *RedisESportRoomChatLogObj) QueryESPortsRoomAllChatLog() []*share_message.TableESPortsLiveRoomMsgLog {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_ESPORTS_ROOM_CHAT_MSG_LOG)
	defer closeFun()
	var lst []*share_message.TableESPortsLiveRoomMsgLog
	err1 := col.Find(bson.M{"SessionId": self.Id}).Sort("-_id").Limit(100).All(&lst)
	easygo.PanicError(err1)
	return lst
}

//获取当前redis记录
func (self *RedisESportRoomChatLogObj) GetRedisESportRoomChatLog() []*share_message.TableESPortsLiveRoomMsgLog {
	var chatList []*share_message.TableESPortsLiveRoomMsgLog
	if !self.IsExistKey() {
		return chatList
	}
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	for _, m := range values {
		log := &share_message.TableESPortsLiveRoomMsgLog{}
		_ = json.Unmarshal([]byte(m), &log)
		chatList = append(chatList, log)
	}
	return chatList
}

func (self *RedisESportRoomChatLogObj) AddRedisESportRoomChatLog(log *share_message.TableESPortsLiveRoomMsgLog) {
	s, err := json.Marshal(log)
	easygo.PanicError(err)
	err = easygo.RedisMgr.GetC().HSet(self.GetKeyId(), easygo.AnytoA(log.GetId()), string(s))
	easygo.PanicError(err)
	self.SetSaveSid()
	self.SetSaveStatus(true)
}

func (self *RedisESportRoomChatLogObj) GetTableESPortsLiveRoomMsgLog(id int64) *share_message.TableESPortsLiveRoomMsgLog {
	var log *share_message.TableESPortsLiveRoomMsgLog

	by, err := easygo.RedisMgr.GetC().HGet(self.GetKeyId(), easygo.AnytoA(id))
	if err != nil {
		log = self.QueryESPortsRoomOneChatLog(id)
		return log
	}
	err = json.Unmarshal(by, &log)
	easygo.PanicError(err)
	return log
}
func (self *RedisESportRoomChatLogObj) UpdateTableESPortsLiveRoomMsgLog(log *share_message.TableESPortsLiveRoomMsgLog) {
	self.AddRedisESportRoomChatLog(log)
}

func GetRedisESportRoomChatLogObj(id int64) *RedisESportRoomChatLogObj {
	kstr := easygo.AnytoA(id)
	obj, ok := ESportRoomChatMgr.Load(kstr)
	if ok && obj != nil {
		tobj := obj.(*RedisESportRoomChatLogObj)
		return tobj
	} else {
		obj := NewRedisESportRoomChatLogObj(id)
		if obj != nil {
			ESportRoomChatMgr.Store(kstr, obj)
		}
		return obj
	}
}
