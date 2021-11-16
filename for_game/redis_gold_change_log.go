package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"sort"
	"time"
)

/*
redis金币变化日志
*/

const (
	GOLD_CHANGE_EXIST_LOG  = "redis_gold_change_log" //新金币变化记录
	GOLD_CHANGE_EXIST_TIME = 1000 * 600              //redis的key删除时间:毫秒
)

type RedisGoldLogObj struct {
	Id string
	RedisBase
}

//写入redis
func NewRedisGoldLog() *RedisGoldLogObj {
	p := &RedisGoldLogObj{
		Id: TABLE_GOLDCHANGELOG, //固定id
	}
	return p.Init()
}

func (self *RedisGoldLogObj) Init() *RedisGoldLogObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_GOLDCHANGELOG)
	self.Sid = GoldChangeLogMgr.GetSid()
	//增加到管理器
	GoldChangeLogMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//if self.IsExistKey() {
	//	self.SetSaveStatus(true)
	//}
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	//self.SetSaveStatus(true)
	return self
}
func (self *RedisGoldLogObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisGoldLogObj) GetKeyId() string { //override
	return TABLE_GOLDCHANGELOG
}

//重写保存方法
func (self *RedisGoldLogObj) SaveToMongo() { //override
	logsData := self.GetRedisGoldLog()
	logs.Info("redis 存储 ：", self.GetKeyId(), logsData)
	var keys []string
	if len(logsData) > 0 {
		var newlst []interface{}
		for _, log := range logsData {
			keys = append(keys, easygo.AnytoA(log.GetLogId()))
			newlst = append(newlst, log)
		}
		col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_GOLDCHANGELOG)
		defer closeFun()
		err1 := col.Insert(newlst...)
		easygo.PanicError(err1)
		_, err := easygo.RedisMgr.GetC().Hdel(self.GetKeyId(), keys...)
		easygo.PanicError(err)
	}
	self.SetSaveStatus(false)
}

//定时更新数据
func (self *RedisGoldLogObj) UpdateData() { //override
	if !self.IsExistKey() {
		GoldChangeLogMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	//if self.GetSaveStatus() { //需要保存的数据进行存储
	//	self.SaveToMongo()
	//}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > GOLD_CHANGE_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		GoldChangeLogMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisGoldLogObj) InitRedis() { //override
	//增加到管理器
	GoldChangeLogMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//重新激活定时器
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisGoldLogObj) GetRedisSaveData() interface{} { //override
	return nil
}
func (self *RedisGoldLogObj) SaveOtherData() { //override
}
func (self *RedisGoldLogObj) QueryGoldLog() *CommonGold {
	data := self.QueryMongoData(self.Id)
	if data != nil {
		var log CommonGold
		StructToOtherStruct(data, &log)
		return &log
	}
	return nil
}

func (self *RedisGoldLogObj) AddRedisGoldLog(log *CommonGold) {
	chatInfo := make(map[int64]string)
	s, _ := json.Marshal(log)
	chatInfo[log.GetLogId()] = string(s)
	err1 := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), chatInfo)
	easygo.PanicError(err1)
	self.SetSaveSid()
	self.SetSaveStatus(true)
}
func (self *RedisGoldLogObj) GetRedisGoldLog() []*CommonGold {
	var lst []*CommonGold
	if !self.IsExistKey() {
		return lst
	}
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	for _, m := range values {
		log := &CommonGold{}
		_ = json.Unmarshal([]byte(m), &log)
		lst = append(lst, log)
	}
	return lst
}

//对外方法
func GetRedisGoldLogObj() *RedisGoldLogObj {
	return GoldChangeLogMgr.GetRedisGoldLogObj()
}
func AddGoldChangeLog(log *CommonGold) {
	GoldChangeLogMgr.AddGoldChangeLog(log)
}

func GetPlayerGoldChangeLogForGoldExtendLog(pid int64, t []int32) []*GoldLog {
	obj := GetRedisGoldLogObj()
	lst := obj.GetRedisGoldLog()
	var newlst []*GoldLog
	for _, log := range lst {
		if log.GetPlayerId() != pid {
			continue
		}
		sourceType := log.GetSourceType()
		if sourceType == GOLD_TYPE_CASH_AFIN || sourceType == GOLD_TYPE_CASH_AFOUT { //等于人工出入款就跳过
			continue
		}
		if len(t) != 0 { //如果有要求类型就判断 t为空代表查询所有类型
			if !util.Int32InSlice(sourceType, t) {
				continue
			}
		}

		extend := &share_message.GoldExtendLog{}
		b, err := json.Marshal(log.Extend)
		easygo.PanicError(err)
		err = json.Unmarshal(b, &extend)
		easygo.PanicError(err)
		log.Extend = extend
		msg := &GoldLog{}
		StructToOtherStruct(log, msg)
		newlst = append(newlst, msg)
	}
	return newlst
}

func GetPageGoldLogs(pid int64, gt []int32, page, num, year, month int) []*GoldLog {
	var year2, month2 int
	if month == 12 {
		year2 = year + 1
		month2 = 1
	} else {
		year2 = year
		month2 = month + 1
	}
	t1 := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).UnixNano() / 1e6
	t2 := time.Date(year2, time.Month(month2), 1, 0, 0, 0, 0, time.Local).UnixNano() / 1e6
	obj := GetRedisGoldLogObj()
	lst := obj.GetRedisGoldLog()
	var newlst []*GoldLog
	for _, log := range lst {
		if log.GetPlayerId() != pid {
			continue
		}
		if len(gt) != 0 { //如果有要求类型就判断 t为空代表查询所有类型
			if !util.Int32InSlice(log.GetSourceType(), gt) {
				continue
			}
		}

		if log.GetCreateTime() < t1 || log.GetCreateTime() > t2 {
			continue
		}

		extend := &share_message.GoldExtendLog{}
		b, err := json.Marshal(log.Extend)
		easygo.PanicError(err)
		err = json.Unmarshal(b, &extend)
		easygo.PanicError(err)
		log.Extend = extend
		msg := &GoldLog{}
		StructToOtherStruct(log, msg)
		newlst = append(newlst, msg)
	}
	sort.Slice(newlst, func(i, j int) bool {
		return newlst[i].GetLogId() > newlst[j].GetLogId() // 降序
		//return newlst[i].GetLogId() < newlst[j].GetLogId() // 升序
	})
	start := (page - 1) * num
	end := page * num
	leng := len(newlst)
	if leng > start { //如果数量够
		if leng >= end {
			return newlst[start:end]
		} else {
			return newlst[start:]
		}
	} else {
		return []*GoldLog{}
	}
}

func GetGoldLogInfoForOrderIds(pid int64, gt int32, orderIds []string) []*GoldLog {
	obj := GetRedisGoldLogObj()
	lst := obj.GetRedisGoldLog()
	var newlst []*GoldLog
	for _, log := range lst {
		if log.GetPlayerId() != pid {
			continue
		}
		if log.GetSourceType() != gt {
			continue
		}

		m := log.Extend.(map[string]interface{})
		orderId := m["OrderId"].(string)
		if !util.InStringSlice(orderId, orderIds) {
			continue
		}

		extend := &share_message.GoldExtendLog{}
		b, err := json.Marshal(log.Extend)
		easygo.PanicError(err)
		err = json.Unmarshal(b, &extend)
		easygo.PanicError(err)
		log.Extend = extend
		msg := &GoldLog{}
		StructToOtherStruct(log, msg)
		newlst = append(newlst, msg)
	}
	return newlst
}

func SaveGoldChangeLogToMongoDB() {
	obj := GetRedisGoldLogObj()
	obj.SaveToMongo()
}
