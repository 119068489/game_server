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
短信召回统计报表内存数据管理
*/
const (
	RECALLREPORT_EXIST_LIST = "recallReport_exist_list" //redis内存中存在的key
	RECALLREPORT_EXIST_TIME = 1000 * 600                //redis的key删除时间:毫秒
)

type RecallReportObj struct {
	Id int64
	RedisBase
}

type RedisRecallReport struct {
	CreateTime  int64 `json:"_id"`
	Pv          int64
	Uv          int64
	DownCount   int64
	RecallCount int64
}

func NewRedisRecallReport(id int64, data ...*share_message.RecallReport) *RecallReportObj {
	p := &RecallReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}

func (self *RecallReportObj) Init(obj *share_message.RecallReport) *RecallReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_RECALL_REPORT)
	self.Sid = RecallReportMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		RecallReportMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryRecallReport(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisRecallReport(obj)
	}

	logs.Info("初始化新的RecallReport管理器:", self.Id)
	return self
}
func (self *RecallReportObj) GetId() interface{} { //override
	return self.Id
}
func (self *RecallReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_RECALL_REPORT, self.Id)
}
func (self *RecallReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		RecallReportMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > RECALLREPORT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		RecallReportMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RecallReportObj) InitRedis() { //override
	obj := QueryRecallReport(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisRecallReport(obj)
}
func (self *RecallReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisRecallReport()
	return data
}
func (self *RecallReportObj) SaveOtherData() { //override
}
func (self *RecallReportObj) SetRedisRecallReport(obj *share_message.RecallReport) {

	//增加到管理器
	RecallReportMgr.Store(obj.GetCreateTime(), self)
	self.AddToExistList(obj.GetCreateTime())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	redisObj := &RedisRecallReport{}
	StructToOtherStruct(obj, redisObj)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redisObj)
	easygo.PanicError(err)

}
func (self *RecallReportObj) GetRedisRecallReport() *share_message.RecallReport {
	obj := &RedisRecallReport{}
	key := self.GetKeyId()
	value, err := easygo.RedisMgr.GetC().HGetAll(key)
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.RecallReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置RedisRecallReport字段值
func SetRedisRecallReportFildVal(id, val int64, fild string) {
	obj := RecallReportMgr.GetRedisRecallReport(id)
	obj.IncrOneValue(fild, val)
}

//获取数据
func GetRedisRecallReport(id int64) *share_message.RecallReport {
	obj := RecallReportMgr.GetRedisRecallReport(id)
	return obj.GetRedisRecallReport()
}

//id查询召回报表
func QueryRecallReport(id int64) *share_message.RecallReport {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_RECALL_REPORT)
	defer closeFun()
	var obj *share_message.RecallReport
	err := col.Find(bson.M{"_id": id}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//管理器
func GetRedisRecallReportMgr(id int64) *RecallReportObj {
	return RecallReportMgr.GetRedisRecallReport(id)
}

//批量保存需要存储的数据
func SaveRedisRecallReportToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_RECALL_REPORT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisRecallReportMgr(id)
		if obj != nil {
			data := obj.GetRedisRecallReport()
			if data.GetPv()+data.GetUv()+data.GetDownCount()+data.GetRecallCount() > 0 {
				saveData = append(saveData, bson.M{"_id": data.GetCreateTime()}, data)
				obj.SetSaveStatus(false)
			}
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_RECALL_REPORT, saveData)
	}
}
