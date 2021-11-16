package for_game

import (
	"encoding/json"
	"game_server/easygo"
)

/*
redis电镜比变化日志
*/

const (
	ESPORT_COIN_CHANGE_EXIST_LOG  = "redis_esport_coin_change_log" //新电竞币变化记录
	ESPORT_COIN_CHANGE_EXIST_TIME = 1000 * 600                     //redis的key删除时间:毫秒
)

type RedisESportCoinLogObj struct {
	Id string
	RedisBase
}

//写入redis
func NewRedisESportCoinLog() *RedisESportCoinLogObj {
	p := &RedisESportCoinLogObj{
		Id: TABLE_ESPORTCHANGELOG, //固定id
	}
	return p.Init()
}

func (self *RedisESportCoinLogObj) Init() *RedisESportCoinLogObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_ESPORTCHANGELOG)
	//增加到管理器
	self.Sid = ESportCoinChangeLogMgr.GetSid()
	ESportCoinChangeLogMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)

	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	return self
}
func (self *RedisESportCoinLogObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisESportCoinLogObj) GetKeyId() string { //override
	return TABLE_ESPORTCHANGELOG
}

//重写保存方法
func (self *RedisESportCoinLogObj) SaveToMongo() { //override
	logsData := self.GetRedisESportCoinLog()
	//logs.Info("redis 存储 ：", logsData)
	var keys []string
	if len(logsData) > 0 {
		var newlst []interface{}
		for _, log := range logsData {
			keys = append(keys, easygo.AnytoA(log.GetLogId()))
			newlst = append(newlst, log)
		}
		col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTCHANGELOG)
		defer closeFun()
		err1 := col.Insert(newlst...)
		easygo.PanicError(err1)
		_, err := easygo.RedisMgr.GetC().Hdel(self.TBName, keys...)
		easygo.PanicError(err)
	}
	self.SetSaveStatus(false)
}

//定时更新数据
func (self *RedisESportCoinLogObj) UpdateData() { //override
	if !self.IsExistKey() {
		ESportCoinChangeLogMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	//if self.GetSaveStatus() { //需要保存的数据进行存储
	self.SaveToMongo()
	//}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > ESPORT_COIN_CHANGE_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		ESportCoinChangeLogMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisESportCoinLogObj) InitRedis() { //override
	//增加到管理器
	ESportCoinChangeLogMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//重新激活定时器
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisESportCoinLogObj) GetRedisSaveData() interface{} { //override
	return nil
}
func (self *RedisESportCoinLogObj) SaveOtherData() { //override
}
func (self *RedisESportCoinLogObj) QueryESportCoinLog() *CommonESportCoin {
	data := self.QueryMongoData(self.Id)
	if data != nil {
		var log CommonESportCoin
		StructToOtherStruct(data, &log)
		return &log
	}
	return nil
}

func (self *RedisESportCoinLogObj) AddRedisESportCoinLog(log *CommonESportCoin) {
	chatInfo := make(map[int64]string)
	s, _ := json.Marshal(log)
	chatInfo[log.GetLogId()] = string(s)
	err1 := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), chatInfo)
	easygo.PanicError(err1)
	self.SetSaveSid()
	self.SetSaveStatus(true)
}
func (self *RedisESportCoinLogObj) GetRedisESportCoinLog() []*CommonESportCoin {
	var lst []*CommonESportCoin
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	for _, m := range values {
		log := &CommonESportCoin{}
		_ = json.Unmarshal([]byte(m), &log)
		lst = append(lst, log)
	}
	return lst
}

////是否存在购买记录
//func (self *RedisESportCoinLogObj) CheckMonthRecharge(playerId, id int64) bool {
//	logsData := self.GetRedisESportCoinLog()
//	for _, log := range logsData {
//		if log.GetPlayerId() == playerId && log.GetCreateTime() > easygo.GetMonth0ClockOfTimestamp(easygo.NowTimestamp())*1000 {
//			extend := &share_message.GoldExtendLog{}
//			b, err := json.Marshal(log.Extend)
//			easygo.PanicError(err)
//			err = json.Unmarshal(b, &extend)
//			easygo.PanicError(err)
//			if extend.GetRedPacketId() == id {
//				return true
//			}
//		}
//	}
//	return false
//}

//对外方法
func GetRedisESportCoinLogObj() *RedisESportCoinLogObj {
	return ESportCoinChangeLogMgr.GetRedisESportCoinLogObj()
}
func AddESportCoinChangeLog(log *CommonESportCoin) {
	ESportCoinChangeLogMgr.AddESportCoinChangeLog(log)
}

//保存日志到数据库
func SaveESportCoinChangeLogToMongoDB() {
	obj := GetRedisESportCoinLogObj()
	obj.SaveToMongo()
}
