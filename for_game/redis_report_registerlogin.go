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
注册登录报表内存数据管理
*/
const (
	REGISTERLOGIN_EXIST_LIST = "registerlogin_exist_list" //redis内存中存在的key
	REGISTERLOGIN_EXIST_TIME = 1000 * 600                 //redis的key删除时间:毫秒
)

type RegisterLoginReportObj struct {
	Id int64
	RedisBase
}
type RedisRegisterLoginReport struct {
	CreateTime         int64 `json:"_id"`
	RegSumCount        int64
	WxRegCount         int64
	PhoneRegCount      int64
	ValidRegSumCount   int64
	ValidWxRegCount    int64
	ValidPhoneRegCount int64
	LoginSumCount      int64
	LoginTimesCount    int64
	RealNameCount      int64
	BankCardCount      int64
	PvCount            int64
	UvCount            int64
	LabelCount         int64
}

func NewRedisRegisterLoginReport(id int64, data ...*share_message.RegisterLoginReport) *RegisterLoginReportObj {
	p := &RegisterLoginReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}
func (self *RegisterLoginReportObj) Init(obj *share_message.RegisterLoginReport) *RegisterLoginReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_LOGIN_REGISTER_REPORT)
	self.Sid = RegisterLoginReportMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		RegisterLoginReportMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryRegisterLoginReport(self.Id, 0)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisRegisterLoginReport(obj)
	}

	logs.Info("初始化新的RegisterLoginReport管理器:", self.Id)
	return self
}
func (self *RegisterLoginReportObj) GetId() interface{} { //override
	return self.Id
}
func (self *RegisterLoginReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_LOGIN_REGISTER_REPORT, self.Id)
}
func (self *RegisterLoginReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		RegisterLoginReportMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REGISTERLOGIN_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		RegisterLoginReportMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RegisterLoginReportObj) InitRedis() { //override
	obj := QueryRegisterLoginReport(self.Id, 0)
	if obj == nil {
		return
	}
	self.SetRedisRegisterLoginReport(obj)
}
func (self *RegisterLoginReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisRegisterLoginReport()
	return data
}
func (self *RegisterLoginReportObj) SaveOtherData() { //override
}
func (self *RegisterLoginReportObj) SetRedisRegisterLoginReport(obj *share_message.RegisterLoginReport) {
	//增加到管理器
	RegisterLoginReportMgr.Store(obj.GetCreateTime(), self)
	self.AddToExistList(obj.GetCreateTime())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	redisObj := &RedisRegisterLoginReport{}
	StructToOtherStruct(obj, redisObj)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redisObj)
	easygo.PanicError(err)

}
func (self *RegisterLoginReportObj) GetRedisRegisterLoginReport() *share_message.RegisterLoginReport {

	obj := &RedisRegisterLoginReport{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.RegisterLoginReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

func SetRedisRegisterLoginReportFildVal(createTime, val int64, fild string) {
	createTime = easygo.Get0ClockTimestamp(createTime)
	obj := RegisterLoginReportMgr.GetRedisRegisterLoginObj(createTime)
	obj.IncrOneValue(fild, val)
}

//管理器
func GetRedisRegisterLoginReportMgr(id int64) *RegisterLoginReportObj {
	return RegisterLoginReportMgr.GetRedisRegisterLoginObj(id)
}

//查询指定时间的埋点注册登录报表
func QueryRegisterLoginReport(querytime int64, role ...int32) *share_message.RegisterLoginReport {
	querytime = easygo.Get0ClockTimestamp(querytime)
	// col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, "bak_report_register_login")
	// if len(role) > 0 {
	// 	if role[0] == 0 {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_LOGIN_REGISTER_REPORT)
	// 	}
	// }

	defer closeFun()
	var obj *share_message.RegisterLoginReport
	err := col.Find(bson.M{"_id": querytime}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//批量保存需要存储的数据
func SaveRedisRegisterLoginReportToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_LOGIN_REGISTER_REPORT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisRegisterLoginReportMgr(id)
		if obj != nil {
			data := obj.GetRedisRegisterLoginReport()
			saveData = append(saveData, bson.M{"_id": data.GetCreateTime()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_LOGIN_REGISTER_REPORT, saveData)
	}
}
