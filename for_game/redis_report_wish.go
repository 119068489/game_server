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
当天许愿池埋点表
*/
const (
	REPORT_WISH_LOG = "wish:report" //当天许愿池埋点日志
)

type WishLogReportObj struct {
	Id int64
	RedisBase
}

type RedisWishLogReport struct {
	CreateTime     int64 `json:"_id"`
	WishTime       int64
	NewPlayer      int32
	OldPlayer      int32
	ExchangeMen    int64
	VExchangeTime  int64
	VExchangeMen   int32
	DareMen        int32
	ChallengeMen   int32
	TwoDayKeep     int32
	ThreeDayKeep   int32
	SevenDayKeep   int32
	FifteenDayKeep int32
	ThirtyDayKeep  int32
}

func NewRedisWishLogReport(id int64, data ...*share_message.WishLogReport) *WishLogReportObj {
	p := &WishLogReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}

func (self *WishLogReportObj) Init(obj *share_message.WishLogReport) *WishLogReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_REPORT_WISH_LOG)
	self.Sid = WishLogReportMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		WishLogReportMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryWishLogReport(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisWishLogReport(obj)
	}

	logs.Info("初始化新的WishLogReport管理器:", self.Id)
	return self
}

func (self *WishLogReportObj) GetId() interface{} { //override
	return self.Id
}

func (self *WishLogReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_REPORT_WISH_LOG, self.Id)
}

func (self *WishLogReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		WishLogReportMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > PLAYERKEEPREPORT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		WishLogReportMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}

func (self *WishLogReportObj) InitRedis() { //override
	obj := QueryWishLogReport(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisWishLogReport(obj)
}

func (self *WishLogReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisWishLogReport()
	return data
}

func (self *WishLogReportObj) SaveOtherData() { //override
}

func (self *WishLogReportObj) SetRedisWishLogReport(obj *share_message.WishLogReport) {
	//增加到管理器
	WishLogReportMgr.Store(obj.GetCreateTime(), self)
	self.AddToExistList(obj.GetCreateTime())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	redisObj := &RedisWishLogReport{}
	StructToOtherStruct(obj, redisObj)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redisObj)
	easygo.PanicError(err)

}

func (self *WishLogReportObj) GetRedisWishLogReport() *share_message.WishLogReport {
	obj := &RedisWishLogReport{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.WishLogReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置数据
func SetRedisWishLogReport(report *share_message.WishLogReport) {
	obj := WishLogReportMgr.GetRedisWishLogReport(report.GetCreateTime())
	obj.SetRedisWishLogReport(report)
}

//设置RedisWishLogReport字段值
func SetRedisWishLogReportFieldVal(createTime, val int64, field string) {
	obj := WishLogReportMgr.GetRedisWishLogReport(createTime)
	obj.IncrOneValue(field, val)
}

//获取数据
func GetRedisWishLogReport(id int64) *WishLogReportObj {
	obj := WishLogReportMgr.GetRedisWishLogReport(id)
	return obj
}

//从mongo中读取数据
func QueryWishLogReport(logintime int64) *share_message.WishLogReport {
	time0 := easygo.Get0ClockTimestamp(logintime)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_REPORT_WISH_LOG)
	defer closeFun()
	var obj *share_message.WishLogReport
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
func SaveRedisWishLogReportToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_REPORT_WISH_LOG, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisWishLogReport(id)
		if obj != nil {
			data := obj.GetRedisWishLogReport()
			saveData = append(saveData, bson.M{"_id": data.GetCreateTime()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_REPORT_WISH_LOG, saveData)
	}
}
