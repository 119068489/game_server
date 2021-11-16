package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"sort"
	"time"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

/*
redis金币变化日志
*/

const (
	DIAMOND_CHANGE_EXIST_TIME = 1000 * 600 //redis的key删除时间:毫秒
)

type RedisDiamondLogObj struct {
	Id string
	RedisBase
}

//写入redis
func NewRedisDiamondLog() *RedisDiamondLogObj {
	p := &RedisDiamondLogObj{
		Id: TABLE_DIAMOND_CHANGELOG, //固定id
	}
	return p.Init()
}

func (self *RedisDiamondLogObj) Init() *RedisDiamondLogObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_DIAMOND_CHANGELOG)
	self.Sid = DiamondChangeLogMgr.GetSid()
	//增加到管理器
	DiamondChangeLogMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	return self
}
func (self *RedisDiamondLogObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisDiamondLogObj) GetKeyId() string { //override
	return TABLE_DIAMOND_CHANGELOG
}

//重写保存方法
func (self *RedisDiamondLogObj) SaveToMongo() { //override
	logsData := self.GetRedisDiamondLog()
	logs.Info("redis 存储 ：", self.GetKeyId(), logsData)
	var keys []string
	if len(logsData) > 0 {
		var newlst []interface{}
		for _, log := range logsData {
			keys = append(keys, easygo.AnytoA(log.GetLogId()))
			newlst = append(newlst, bson.M{"_id": log.GetLogId()}, log)
		}
		//col, closeFun := easygo.MongoMgr.GetC(self.DBName, self.TBName)
		//defer closeFun()
		//err1 := col.Insert(newlst...)
		//easygo.PanicError(err1)
		UpsertAll(easygo.MongoMgr, self.DBName, self.TBName, newlst)
		_, err := easygo.RedisMgr.GetC().Hdel(self.GetKeyId(), keys...)
		easygo.PanicError(err)
	}
	logs.Info("RedisDiamondLogObj  redis 存储完成")
	self.SetSaveStatus(false)
}

//定时更新数据
func (self *RedisDiamondLogObj) UpdateData() { //override
	if !self.IsExistKey() {
		DiamondChangeLogMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > DIAMOND_CHANGE_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		DiamondChangeLogMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisDiamondLogObj) InitRedis() { //override
	//增加到管理器
	DiamondChangeLogMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//重新激活定时器
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisDiamondLogObj) GetRedisSaveData() interface{} { //override
	return nil
}
func (self *RedisDiamondLogObj) SaveOtherData() { //override
}
func (self *RedisDiamondLogObj) QueryDiamondLog() *CommonDiamond {
	data := self.QueryMongoData(self.Id)
	if data != nil {
		var log CommonDiamond
		StructToOtherStruct(data, &log)
		return &log
	}
	return nil
}

func (self *RedisDiamondLogObj) AddRedisDiamondLog(log *CommonDiamond) {
	chatInfo := make(map[int64]string)
	s, _ := json.Marshal(log)
	chatInfo[log.GetLogId()] = string(s)

	err1 := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), chatInfo)
	easygo.PanicError(err1)
	self.SetSaveSid()
	self.SetSaveStatus(true)
}
func (self *RedisDiamondLogObj) GetRedisDiamondLog() []*CommonDiamond {
	var lst []*CommonDiamond
	if !self.IsExistKey() {
		return lst
	}
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	for _, m := range values {
		log := &CommonDiamond{}
		_ = json.Unmarshal([]byte(m), &log)
		lst = append(lst, log)
	}
	return lst
}

//对外方法
func GetRedisDiamondLogObj() *RedisDiamondLogObj {
	return DiamondChangeLogMgr.GetRedisDiamondLogObj()
}
func AddDiamondChangeLog(log *CommonDiamond) {
	DiamondChangeLogMgr.AddDiamondChangeLog(log)
}

func GetPlayerDiamondChangeLogForDiamondExtendLog(pid int64, t []int32) []*DiamondLog {
	obj := GetRedisDiamondLogObj()
	lst := obj.GetRedisDiamondLog()
	var newlst []*DiamondLog
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
		msg := &DiamondLog{}
		StructToOtherStruct(log, msg)
		newlst = append(newlst, msg)
	}
	return newlst
}

func GetPageDiamondLogs(pid int64, gt []int32, page, num, year, month int) []*DiamondLog {
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
	obj := GetRedisDiamondLogObj()
	lst := obj.GetRedisDiamondLog()
	var newlst []*DiamondLog
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
		msg := &DiamondLog{}
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
		return []*DiamondLog{}
	}
}

func GetDiamondLogInfoForOrderIds(pid int64, gt int32, orderIds []string) []*DiamondLog {
	obj := GetRedisDiamondLogObj()
	lst := obj.GetRedisDiamondLog()
	var newlst []*DiamondLog
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
		msg := &DiamondLog{}
		StructToOtherStruct(log, msg)
		newlst = append(newlst, msg)
	}
	return newlst
}

func SaveDiamondChangeLogToMongoDB() {
	obj := GetRedisDiamondLogObj()
	obj.SaveToMongo()
}

//是否存在购买记录
func (self *RedisDiamondLogObj) CheckMonthRecharge(playerId, id int64) bool {
	logsData := self.GetRedisDiamondLog()
	for _, log := range logsData {
		if log.GetPlayerId() == playerId && log.GetCreateTime() > easygo.GetMonth0ClockOfTimestamp(easygo.NowTimestamp())*1000 {
			extend := &share_message.GoldExtendLog{}
			b, err := json.Marshal(log.Extend)
			easygo.PanicError(err)
			err = json.Unmarshal(b, &extend)
			easygo.PanicError(err)
			if extend.GetRedPacketId() == id {
				return true
			}
		}
	}
	return false
}

// 测试代码,后续删除
//func (self *RedisDiamondLogObj) CheckMonthRecharge1(playerId, id int64) bool {
//	logsData := self.GetRedisDiamondLog()
//	for _, log := range logsData {
//		//if log.GetPlayerId() == playerId && log.GetCreateTime() > easygo.GetMonth0ClockOfTimestamp(easygo.NowTimestamp())*1000 {
//		if log.GetPlayerId() == playerId && log.GetCreateTime()+(10*60*1000) > GetMillSecond() { // 10分钟内
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

// 获取每月充值记录
func CheckPlayerDiamondMonthRecharge(playerId, id int64) bool {

	n := FindAllCount(MONGODB_NINGMENG, TABLE_DIAMOND_CHANGELOG, bson.M{"PlayerId": playerId, "SourceType": DIAMOND_TYPE_EXCHANGE_IN, "CreateTime": bson.M{"$gte": easygo.GetMonth0ClockOfTimestamp(easygo.NowTimestamp()) * 1000}, "Extend.RedPacketId": id})
	return n == 0
}

//func CheckPlayerDiamondMonthRecharge1(playerId, id int64) bool {
//
//	n := FindAllCount(MONGODB_NINGMENG, TABLE_DIAMOND_CHANGELOG, bson.M{"PlayerId": playerId, "SourceType": DIAMOND_TYPE_EXCHANGE_IN, "CreateTime": bson.M{"$gte": GetMillSecond() - (10 * 60 * 1000)}, "Extend.RedPacketId": id})
//	return n == 0
//}
