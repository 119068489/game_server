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
运营渠道数据汇总报表内存数据管理
*/
const (
	OPERATIONCHANNELREPORT_EXIST_LIST = "operationchannelreport_exist_list" //redis内存中存在的key
	OPERATIONCHANNELREPORT_EXIST_TIME = 1000 * 600                          //redis的key删除时间:毫秒
)

type OperationChannelReportObj struct {
	Id string
	RedisBase
}
type RedisOperationChannelReport struct {
	Id               string `json:"_id"`
	ChannelNo        string
	ChannelName      string
	Cooperation      int32
	CreateTime       int64
	RegCount         int64
	ValidRegCount    int64
	LoginCount       int64
	ActDevCount      int64
	ValidActDevCount int64
	DownLoadCount    int64
	UvCount          int64
	OnlineSum        int64
	NextKeep         int64
	// ShopOrderSumCount int64
	// ShopOrderNewCount int64
	// ShopOrderOldCount int64
	// ShopDealSumAmount int64
	// ShopDealNewAmount int64
	// ShopDealOldAmount int64
	// RechargeSumAmount int64
	// RechargeNewAmount int64
	// RechargeOldAmount int64
	// WithdrawSumAmount int64
	// WithdrawNewAmount int64
	// WithdrawOldAmount int64
}

func NewRedisOperationChannelReport(id string, data ...*share_message.OperationChannelReport) *OperationChannelReportObj {
	p := &OperationChannelReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}

func (self *OperationChannelReportObj) Init(obj *share_message.OperationChannelReport) *OperationChannelReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_OPERATION_CHANNEL_REPORT)
	self.Sid = OperationChannelReportMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		OperationChannelReportMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryOperationChannelReport(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisOperationChannelReport(obj)
	}

	logs.Info("初始化新的OperationChannelReport管理器:", self.Id)
	return self
}
func (self *OperationChannelReportObj) GetId() interface{} { //override
	return self.Id
}
func (self *OperationChannelReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_OPERATION_CHANNEL_REPORT, self.Id)
}
func (self *OperationChannelReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		OperationChannelReportMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > OPERATIONCHANNELREPORT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		OperationChannelReportMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *OperationChannelReportObj) InitRedis() { //override
	obj := QueryOperationChannelReport(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisOperationChannelReport(obj)
}
func (self *OperationChannelReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisOperationChannelReport()
	return data
}
func (self *OperationChannelReportObj) SaveOtherData() { //override
}
func (self *OperationChannelReportObj) SetRedisOperationChannelReport(obj *share_message.OperationChannelReport) {
	//增加到管理器
	OperationChannelReportMgr.Store(obj.GetId(), self)
	self.AddToExistList(obj.GetId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	redisObj := &RedisOperationChannelReport{}
	StructToOtherStruct(obj, redisObj)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redisObj)
	easygo.PanicError(err)

}
func (self *OperationChannelReportObj) GetRedisOperationChannelReport() *share_message.OperationChannelReport {

	obj := &RedisOperationChannelReport{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.OperationChannelReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置数据
func SetRedisOperationChannelReport(report *share_message.OperationChannelReport) {
	channle, createTime := report.GetChannelNo(), report.GetCreateTime()
	if report.GetCreateTime() == 0 || report.GetChannelNo() == "" {
		return
	}
	id := MakeNewString(channle, createTime)
	obj := OperationChannelReportMgr.GetRedisOperationChannelReport(id, channle, createTime)
	obj.SetRedisOperationChannelReport(report)
}

//设置RedisOperationChannelReport字段值
func SetRedisOperationChannelReportFildVal(createTime, val int64, channle string, fild string) {
	if createTime == 0 || channle == "" || fild == "" {
		return
	}
	createTime = easygo.Get0ClockTimestamp(createTime)
	id := MakeNewString(channle, createTime)
	obj := OperationChannelReportMgr.GetRedisOperationChannelReport(id, channle, createTime)
	obj.IncrOneValue(fild, val)
}

//获取数据
func GetRedisOperationChannelReport(id string) *OperationChannelReportObj {
	obj := OperationChannelReportMgr.GetRedisOperationChannelReport(id, "", 0)
	return obj
}

//从mongo中读取转账
func QueryOperationChannelReport(id string) *share_message.OperationChannelReport {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_OPERATION_CHANNEL_REPORT)
	defer closeFun()
	var obj *share_message.OperationChannelReport
	err := col.Find(bson.M{"_id": id}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//批量保存需要存储的数据
func SaveRedisOperationChannelReportToMongo() {
	ids := []string{}
	GetAllRedisSaveList(TABLE_OPERATION_CHANNEL_REPORT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisOperationChannelReport(id)
		if obj != nil {
			data := obj.GetRedisOperationChannelReport()
			saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_OPERATION_CHANNEL_REPORT, saveData)
	}
}
