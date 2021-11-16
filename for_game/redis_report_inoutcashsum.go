package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

/*
出入款汇总报表
*/
const (
	INOUTCASHSUMREPORT_EXIST_LIST = "inoutcashsumreport_exist_list" //redis内存中存在的key
	INOUTCASHSUMREPORT_EXIST_TIME = 1000 * 600                      //redis的key删除时间:毫秒
)

type InOutCashSumReportObj struct {
	Id int64
	RedisBase
}

type RedisInOutCashSumReport struct {
	CreateTime    int64 `json:"_id"`
	Recharge      int64
	Withdraw      int64
	Redundant     int64
	RechargeTimes int64
	RechargeCount int64
	WithdrawTimes int64
	WithdrawCount int64
}

func NewRedisInOutCashSumReport(id int64, data ...*share_message.InOutCashSumReport) *InOutCashSumReportObj {
	p := &InOutCashSumReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}
func (self *InOutCashSumReportObj) Init(obj *share_message.InOutCashSumReport) *InOutCashSumReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_INOUTCASHSUM_REPORT)
	self.Sid = InOutCashSumReportMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		InOutCashSumReportMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryInOutCashSumReport(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisInOutCashSumReport(obj)
	}

	logs.Info("初始化新的InOutCashSumReport管理器:", self.Id)
	return self
}
func (self *InOutCashSumReportObj) GetId() interface{} { //override
	return self.Id
}
func (self *InOutCashSumReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_INOUTCASHSUM_REPORT, self.Id)
}
func (self *InOutCashSumReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		InOutCashSumReportMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > INOUTCASHSUMREPORT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		InOutCashSumReportMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *InOutCashSumReportObj) InitRedis() { //override
	obj := QueryInOutCashSumReport(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisInOutCashSumReport(obj)
}
func (self *InOutCashSumReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisInOutCashSumReport()
	return data
}
func (self *InOutCashSumReportObj) SaveOtherData() { //override
}
func (self *InOutCashSumReportObj) SetRedisInOutCashSumReport(obj *share_message.InOutCashSumReport) {
	//增加到管理器
	InOutCashSumReportMgr.Store(obj.GetCreateTime(), self)
	self.AddToExistList(obj.GetCreateTime())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	redisObj := &RedisInOutCashSumReport{}
	StructToOtherStruct(obj, redisObj)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redisObj)
	easygo.PanicError(err)
}
func (self *InOutCashSumReportObj) GetRedisInOutCashSumReport() *share_message.InOutCashSumReport {

	obj := &RedisInOutCashSumReport{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.InOutCashSumReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置数据
func SetRedisInOutCashSumReport(report *share_message.InOutCashSumReport) {
	obj := InOutCashSumReportMgr.GetRedisInOutCashSumReport(report.GetCreateTime())
	obj.SetRedisInOutCashSumReport(report)
}

//设置字段值自增
func SetRedisInOutCashSumReportFildVal(id int64, val int64, fild string) {
	obj := InOutCashSumReportMgr.GetRedisInOutCashSumReport(id)
	obj.IncrOneValue(fild, val)
}

//设置字段值更新
func UpdateRedisInOutCashSumReportFildVal(id int64, val int64, fild string) {
	obj := InOutCashSumReportMgr.GetRedisInOutCashSumReport(id)
	obj.SetOneValue(fild, val)
}

//获取数据
func GetRedisInOutCashSumReport(id int64) *InOutCashSumReportObj {
	obj := InOutCashSumReportMgr.GetRedisInOutCashSumReport(id)
	return obj

}

//从mongo中读取数据
func QueryInOutCashSumReport(querytime int64) *share_message.InOutCashSumReport {
	time0 := easygo.Get0ClockTimestamp(querytime)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_INOUTCASHSUM_REPORT)
	defer closeFun()
	var obj *share_message.InOutCashSumReport
	err := col.Find(bson.M{"_id": time0}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//批量保存需要存储的数据
func SaveRedisInOutCashSumReportToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_INOUTCASHSUM_REPORT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisInOutCashSumReport(id)
		if obj != nil {
			data := obj.GetRedisInOutCashSumReport()
			saveData = append(saveData, bson.M{"_id": data.GetCreateTime()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_INOUTCASHSUM_REPORT, saveData)
	}
}
