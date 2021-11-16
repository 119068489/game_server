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
按钮点击行为报表内存数据管理
*/
const (
	BTNCREPORT_EXIST_LIST = "buttonclickreport_exist_list" //redis内存中存在的key
	BTNCREPORT_EXIST_TIME = 1000 * 600                     //redis的key删除时间:毫秒
)

type ButtonClickReportObj struct {
	Id int64
	RedisBase
}

type RedisButtonClickReport struct {
	CreateTime         int64 `json:"_id"`
	AgreementYes       int64
	AgreementNo        int64
	Phone              int64
	WeChat             int64
	OneClick           int64
	OneClickCount      int64
	OtherClick         int64
	LoginBack          int64
	SendCode           int64
	ReSendCode         int64
	InterestOk         int64
	InterestOkCount    int64
	InterestBack       int64
	InterestBackCount  int64
	RecommendSkip      int64
	RecommendSkipCount int64
	RecommendNext      int64
	RecommendNextCount int64
	InNmBtnCLick       int64
	InNmBtnCount       int64
	InPhoneBack        int64
	InCodeBack         int64
	InfoBack           int64
}

func NewRedisButtonClickReport(id int64, data ...*share_message.ButtonClickReport) *ButtonClickReportObj {
	p := &ButtonClickReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}

func (self *ButtonClickReportObj) Init(obj *share_message.ButtonClickReport) *ButtonClickReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_BUTTON_CLICK_REPORT)
	self.Sid = ButtonClickReportMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		ButtonClickReportMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryButtonClickReport(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisButtonClickReport(obj)
	}

	logs.Info("初始化新的ButtonClickReport管理器:", self.Id)
	return self
}
func (self *ButtonClickReportObj) GetId() interface{} { //override
	return self.Id
}
func (self *ButtonClickReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_BUTTON_CLICK_REPORT, self.Id)
}
func (self *ButtonClickReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		ButtonClickReportMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > BTNCREPORT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		ButtonClickReportMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *ButtonClickReportObj) InitRedis() { //override
	obj := QueryButtonClickReport(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisButtonClickReport(obj)
}
func (self *ButtonClickReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisButtonClickReport()
	return data
}
func (self *ButtonClickReportObj) SaveOtherData() { //override
}
func (self *ButtonClickReportObj) SetRedisButtonClickReport(obj *share_message.ButtonClickReport) {

	//增加到管理器
	ButtonClickReportMgr.Store(obj.GetCreateTime(), self)
	self.AddToExistList(obj.GetCreateTime())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	redisObj := &RedisButtonClickReport{}
	StructToOtherStruct(obj, redisObj)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redisObj)
	easygo.PanicError(err)

}
func (self *ButtonClickReportObj) GetRedisButtonClickReport() *share_message.ButtonClickReport {
	obj := &RedisButtonClickReport{}
	key := self.GetKeyId()
	value, err := easygo.RedisMgr.GetC().HGetAll(key)
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.ButtonClickReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置RedisButtonClickReport字段值
func SetRedisButtonClickReportFildVal(id, val int64, fild string) {
	obj := ButtonClickReportMgr.GetRedisButtonClickReport(id)
	obj.IncrOneValue(fild, val)
}

//获取数据
func GetRedisButtonClickReport(id int64) *share_message.ButtonClickReport {
	obj := ButtonClickReportMgr.GetRedisButtonClickReport(id)
	return obj.GetRedisButtonClickReport()
}

//id查询报表
func QueryButtonClickReport(id int64) *share_message.ButtonClickReport {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BUTTON_CLICK_REPORT)
	defer closeFun()
	var obj *share_message.ButtonClickReport
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
func GetRedisButtonClickReportMgr(id int64) *ButtonClickReportObj {
	return ButtonClickReportMgr.GetRedisButtonClickReport(id)
}

//批量保存需要存储的数据
func SaveRedisButtonClickReportToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_BUTTON_CLICK_REPORT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisButtonClickReportMgr(id)
		if obj != nil {
			data := obj.GetRedisButtonClickReport()
			saveData = append(saveData, bson.M{"_id": data.GetCreateTime()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_BUTTON_CLICK_REPORT, saveData)
	}
}
