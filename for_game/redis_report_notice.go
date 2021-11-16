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
系统通知报表
*/
const (
	NOTICEREPORT_EXIST_LIST = "noticereport_exist_list" //redis内存中存在的key
	NOTICEREPORT_EXIST_TIME = 1000 * 600                //redis的key删除时间:毫秒
)

type NoticeReportObj struct {
	Id int32
	RedisBase
}

type RedisNoticeReport struct {
	Id         int32 `json:"_id"`
	CreateTime int64
	PushPlayer int64
	Clicks     int64
}

func NewRedisNoticeReport(id int32, data ...*share_message.ArticleReport) *NoticeReportObj {
	p := &NoticeReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}
func (self *NoticeReportObj) Init(obj *share_message.ArticleReport) *NoticeReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_NOTICE_REPORT)
	self.Sid = NoticeReportMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		NoticeReportMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryNoticeReport(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisNoticeReport(obj)
	}

	logs.Info("初始化新的NoticeReport管理器:", self.Id)
	return self
}
func (self *NoticeReportObj) GetId() interface{} { //override
	return self.Id
}
func (self *NoticeReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_NOTICE_REPORT, self.Id)
}
func (self *NoticeReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		NoticeReportMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > NOTICEREPORT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		NoticeReportMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *NoticeReportObj) InitRedis() { //override
	if self.IsExistKey() {
		return
	}
	obj := QueryNoticeReport(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisNoticeReport(obj)
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//重新激活定时器
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *NoticeReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisNoticeReport()
	return data
}
func (self *NoticeReportObj) SaveOtherData() { //override
}
func (self *NoticeReportObj) SetRedisNoticeReport(obj *share_message.ArticleReport) {
	//增加到管理器
	NoticeReportMgr.Store(obj.GetId(), self)
	self.AddToExistList(obj.GetId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	redisObj := &RedisNoticeReport{}
	StructToOtherStruct(obj, redisObj)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redisObj)
	easygo.PanicError(err)
}
func (self *NoticeReportObj) GetRedisNoticeReport() *share_message.ArticleReport {

	obj := &RedisNoticeReport{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.ArticleReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置数据
func SetRedisNoticeReport(report *share_message.ArticleReport) {
	obj := NoticeReportMgr.GetRedisNoticeReport(int32(report.GetId()))
	obj.SetRedisNoticeReport(report)
}

//设置RedisNoticeReport字段值
func SetRedisNoticeReportFildVal(id int32, val int64, fild string) {
	obj := NoticeReportMgr.GetRedisNoticeReport(id)
	obj.IncrOneValue(fild, val)
}

//获取数据
func GetRedisNoticeReport(id int32) *NoticeReportObj {
	obj := NoticeReportMgr.GetRedisNoticeReport(id)
	return obj
}

//id查询通知报表
func QueryNoticeReport(id int32) *share_message.ArticleReport {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_NOTICE_REPORT)
	defer closeFun()
	var obj *share_message.ArticleReport
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
func SaveRedisNoticeReportToMongo() {
	ids := []int32{}
	GetAllRedisSaveList(TABLE_NOTICE_REPORT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisNoticeReport(id)
		if obj != nil {
			data := obj.GetRedisNoticeReport()
			saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_NOTICE_REPORT, saveData)
	}
}
