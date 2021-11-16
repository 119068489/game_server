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
用户留存报表
*/
const (
	PLAYERKEEPREPORT_EXIST_LIST = "playerkeepreport_exist_list" //redis内存中存在的key
	PLAYERKEEPREPORT_EXIST_TIME = 1000 * 600                    //redis的key删除时间:毫秒
)

type PlayerKeepReportObj struct {
	Id int64
	RedisBase
}
type RedisPlayerKeepReport struct {
	CreateTime      int64 `json:"_id"`
	TodayRegister   int32
	NextKeep        int32
	ThreeKeep       int32
	FourKeep        int32
	FiveKeep        int32
	SixKeep         int32
	SevenKeep       int32
	EightKeep       int32
	NineKeep        int32
	TenKeep         int32
	ElevenKeep      int32
	TwelveKeep      int32
	ThirteenKeep    int32
	FourteenKeep    int32
	FifteenKeep     int32
	SixteenKeep     int32
	SeventeenKeep   int32
	EighteenKeep    int32
	NineteenKeep    int32
	TwentyKeep      int32
	TwentyOneKeep   int32
	TwentyTwoKeep   int32
	TwentyThreeKeep int32
	TwentyFourKeep  int32
	TwentyFiveKeep  int32
	TwentySixKeep   int32
	TwentySevenKeep int32
	TwentyEightKeep int32
	TwentyNineKeep  int32
	Thirtykeep      int32
}

func NewRedisPlayerKeepReport(id int64, data ...*share_message.PlayerKeepReport) *PlayerKeepReportObj {
	p := &PlayerKeepReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}
func (self *PlayerKeepReportObj) Init(obj *share_message.PlayerKeepReport) *PlayerKeepReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYERKEEPREPORT)
	self.Sid = PlayerKeepReportMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		PlayerKeepReportMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryPlayerKeepReport(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisPlayerKeepReport(obj)
	}

	logs.Info("初始化新的PlayerKeepReport管理器:", self.Id)
	return self
}
func (self *PlayerKeepReportObj) GetId() interface{} { //override
	return self.Id
}
func (self *PlayerKeepReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_PLAYERKEEPREPORT, self.Id)
}
func (self *PlayerKeepReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		PlayerKeepReportMgr.Delete(self.Id) // 释放对象
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
		PlayerKeepReportMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *PlayerKeepReportObj) InitRedis() { //override
	obj := QueryPlayerKeepReport(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisPlayerKeepReport(obj)
}
func (self *PlayerKeepReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisPlayerKeepReport()
	return data
}
func (self *PlayerKeepReportObj) SaveOtherData() { //override
}
func (self *PlayerKeepReportObj) SetRedisPlayerKeepReport(obj *share_message.PlayerKeepReport) {
	//增加到管理器
	PlayerKeepReportMgr.Store(obj.GetCreateTime(), self)
	self.AddToExistList(obj.GetCreateTime())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	redisObj := &RedisPlayerKeepReport{}
	StructToOtherStruct(obj, redisObj)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redisObj)
	easygo.PanicError(err)

}
func (self *PlayerKeepReportObj) GetRedisPlayerKeepReport() *share_message.PlayerKeepReport {

	obj := &RedisPlayerKeepReport{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.PlayerKeepReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置数据
func SetRedisPlayerKeepReport(report *share_message.PlayerKeepReport) {
	obj := PlayerKeepReportMgr.GetRedisPlayerKeepReport(report.GetCreateTime())
	obj.SetRedisPlayerKeepReport(report)
}

//设置RedisPlayerKeepReport字段值
func SetRedisPlayerKeepReportFildVal(createTime, val int64, fild string) {
	obj := PlayerKeepReportMgr.GetRedisPlayerKeepReport(createTime)
	obj.IncrOneValue(fild, val)
}

//获取数据
func GetRedisPlayerKeepReport(id int64) *PlayerKeepReportObj {
	obj := PlayerKeepReportMgr.GetRedisPlayerKeepReport(id)
	return obj
}

//从mongo中读取数据
func QueryPlayerKeepReport(logintime int64) *share_message.PlayerKeepReport {
	time0 := easygo.Get0ClockTimestamp(logintime)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYERKEEPREPORT)
	defer closeFun()
	var obj *share_message.PlayerKeepReport
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
func SaveRedisPlayerKeepReportToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_PLAYERKEEPREPORT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisPlayerKeepReport(id)
		if obj != nil {
			data := obj.GetRedisPlayerKeepReport()
			saveData = append(saveData, bson.M{"_id": data.GetCreateTime()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYERKEEPREPORT, saveData)
	}
}
