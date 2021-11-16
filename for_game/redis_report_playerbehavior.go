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
用户行为报表内存数据管理
*/
const (
	PLAYERBEHAVIORREPORT_EXIST_LIST = "playerbehaviorreport_exist_list" //redis内存中存在的key
	PLAYERBEHAVIORREPORT_EXIST_TIME = 1000 * 600                        //redis的key删除时间:毫秒
)

type PlayerBehaviorReportObj struct {
	Id int64
	RedisBase
}
type RedisPlayerBehaviorReport struct {
	CreateTime               int64 `json:"_id"`
	SendMsgCount             int64
	SendRedpacketPlayerCount int64
	SendRedpacketCount       int64
	SendRedpacketMoney       int64
	RobRedpacketPlayerCount  int64
	RobRedpacketCount        int64
	RobRedpacketMoney        int64
	TransferPlayerCount      int64
	TransferCount            int64
	TransferMoney            int64
	ShopOrderCount           int64
	ShopOrderMoney           int64
	OneDialogue              int64
	BindCard                 int64
	Start                    int64
	Reply                    int64
}

func NewRedisPlayerBehaviorReport(id int64, data ...*share_message.PlayerBehaviorReport) *PlayerBehaviorReportObj {
	p := &PlayerBehaviorReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}
func (self *PlayerBehaviorReportObj) Init(obj *share_message.PlayerBehaviorReport) *PlayerBehaviorReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_BEHAVIOR_REPORT)
	self.Sid = PlayerBehaviorReportMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		PlayerBehaviorReportMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryPlayerBehaviorReport(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisPlayerBehaviorReport(obj)
	}
	logs.Info("初始化新的PlayerBehaviorReport管理器:", self.Id)
	return self
}
func (self *PlayerBehaviorReportObj) GetId() interface{} { //override
	return self.Id
}
func (self *PlayerBehaviorReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_PLAYER_BEHAVIOR_REPORT, self.Id)
}
func (self *PlayerBehaviorReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		PlayerBehaviorReportMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > PLAYERBEHAVIORREPORT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		PlayerBehaviorReportMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *PlayerBehaviorReportObj) InitRedis() { //override
	if self.IsExistKey() {
		return
	}
	obj := QueryPlayerBehaviorReport(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisPlayerBehaviorReport(obj)
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//重新激活定时器
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *PlayerBehaviorReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisPlayerBehaviorReport()
	return data
}
func (self *PlayerBehaviorReportObj) SaveOtherData() { //override
}
func (self *PlayerBehaviorReportObj) SetRedisPlayerBehaviorReport(obj *share_message.PlayerBehaviorReport) {
	//增加到管理器
	PlayerBehaviorReportMgr.Store(obj.GetCreateTime(), self)
	self.AddToExistList(obj.GetCreateTime())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	redisObj := &RedisPlayerBehaviorReport{}
	StructToOtherStruct(obj, redisObj)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redisObj)
	easygo.PanicError(err)
}
func (self *PlayerBehaviorReportObj) GetRedisPlayerBehaviorReport() *share_message.PlayerBehaviorReport {

	obj := &RedisPlayerBehaviorReport{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.PlayerBehaviorReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置数据
//func SetRedisPlayerBehaviorReport(report *share_message.PlayerBehaviorReport) {
//	obj := PlayerBehaviorReportMgr.GetRedisPlayerBehaviorReport(report.GetCreateTime())
//	obj.SetRedisPlayerBehaviorReport(report)
//}

//设置RedisPlayerBehaviorReport字段值
func SetRedisPlayerBehaviorReportFildVal(createTime, val int64, fild string) {
	obj := PlayerBehaviorReportMgr.GetRedisPlayerBehaviorReport(createTime)
	obj.IncrOneValue(fild, val)
}

//获取数据
func GetRedisPlayerBehaviorReport(id int64) *PlayerBehaviorReportObj {
	obj := PlayerBehaviorReportMgr.GetRedisPlayerBehaviorReport(id)
	return obj
}

//从mongo中读取转账
func QueryPlayerBehaviorReport(querytime int64) *share_message.PlayerBehaviorReport {
	time0 := easygo.Get0ClockTimestamp(querytime)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BEHAVIOR_REPORT)
	defer closeFun()
	var obj *share_message.PlayerBehaviorReport
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
func SaveRedisPlayerBehaviorReportToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_PLAYER_BEHAVIOR_REPORT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisPlayerBehaviorReport(id)
		if obj != nil {
			data := obj.GetRedisPlayerBehaviorReport()
			saveData = append(saveData, bson.M{"_id": data.GetCreateTime()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_BEHAVIOR_REPORT, saveData)
	}
}
