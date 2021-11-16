package for_game

import (
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"

	"github.com/akqp2019/mgo/bson"
)

//玩家之间亲密度
const REDIS_PLAYER_INTIMACY_EXIST_TIME = 1000 * 600
const PLAYER_INTIMACY_MAX_LV = 5

const PLAYER_INTIMACY_REDUCE_DAY = 3 * 86400 //秒

type RedisPlayerIntimacyObj struct {
	Id string
	RedisBase
}
type PlayerIntimacyEx struct {
	Id          string `json:"_id"` //亲密度id：玩家ID_玩家ID 小到大排序组装成
	IntimacyVal int64  //亲密度值
	IntimacyLv  int32  //亲密度等级
	LastTime    int64  //0表示没加过，上一次加亲密度时间
	IsSayHi     bool   //是否sayHi过
}

//写入redis
func NewRedisPlayerIntimacy(id string, data ...*share_message.PlayerIntimacy) *RedisPlayerIntimacyObj {
	p := &RedisPlayerIntimacyObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}

func (self *RedisPlayerIntimacyObj) Init(obj *share_message.PlayerIntimacy) *RedisPlayerIntimacyObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_INTIMACY)
	self.Sid = PlayerIntimacyMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		PlayerIntimacyMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryPlayerIntimacy(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisPlayerIntimacy(obj)
	}
	return self
}
func (self *RedisPlayerIntimacyObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisPlayerIntimacyObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_PLAYER_INTIMACY, self.Id)
}

//定时更新数据
func (self *RedisPlayerIntimacyObj) UpdateData() { //override
	if !self.IsExistKey() {
		PlayerIntimacyMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDIS_PLAYER_INTIMACY_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		PlayerIntimacyMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisPlayerIntimacyObj) InitRedis() { //override
	obj := self.QueryPlayerIntimacy(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisPlayerIntimacy(obj)
}
func (self *RedisPlayerIntimacyObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisPlayerIntimacy()
	return data
}

func (self *RedisPlayerIntimacyObj) SaveOtherData() { //override
}

//获取查询一条记录
func (self *RedisPlayerIntimacyObj) QueryPlayerIntimacy(id string) *share_message.PlayerIntimacy {
	data := self.QueryMongoData(id)
	if data != nil {
		var newData share_message.PlayerIntimacy
		StructToOtherStruct(data, &newData)
		return &newData
	}
	//redis中获取数据初始化
	if self.IsExistKey() {
		return self.GetRedisPlayerIntimacy()
	}
	return nil
}
func (self *RedisPlayerIntimacyObj) SetRedisPlayerIntimacy(obj *share_message.PlayerIntimacy) {
	PlayerIntimacyMgr.Store(obj.GetId(), self)
	self.AddToExistList(obj.GetId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	ex := &PlayerIntimacyEx{}
	StructToOtherStruct(obj, ex)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), ex)
	easygo.PanicError(err)

}

//获取当前redis记录
func (self *RedisPlayerIntimacyObj) GetRedisPlayerIntimacy() *share_message.PlayerIntimacy {
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	if len(value) == 0 {
		return nil
	}
	easygo.PanicError(err)
	var base PlayerIntimacyEx
	err = redis.ScanStruct(value, &base)
	easygo.PanicError(err)
	var newBase *share_message.PlayerIntimacy
	StructToOtherStruct(base, &newBase)
	return newBase
}

//获取亲密度数据给客户端
func (self *RedisPlayerIntimacyObj) GetPlayerIntimacyToClient() *client_hall.IntimacyInfoResp {
	config := GetConfigIntimacy(self.GetIntimacyLv())
	resp := &client_hall.IntimacyInfoResp{
		Id:             easygo.NewString(self.Id),
		IntimacyLv:     easygo.NewInt32(self.GetIntimacyLv()),
		IntimacyVal:    easygo.NewInt64(self.GetIntimacyVal()),
		IntimacyMaxVal: easygo.NewInt64(config.GetMaxVal()),
	}
	return resp
}

//获取当前亲密度值
func (self *RedisPlayerIntimacyObj) GetIntimacyVal() int64 {
	var val int64
	self.GetOneValue("IntimacyVal", &val)
	return val
}

//设置亲密增值
func (self *RedisPlayerIntimacyObj) SetIntimacyVal(val int64) {
	self.SetOneValue("IntimacyVal", val)
}

//增加亲密增值
func (self *RedisPlayerIntimacyObj) AddIntimacyVal(val int64) int64 {
	return self.IncrOneValue("IntimacyVal", val)
}

//获取当前亲密度等级
func (self *RedisPlayerIntimacyObj) GetIntimacyLv() int32 {
	var val int32
	self.GetOneValue("IntimacyLv", &val)
	return val
}

//设置亲密增等级
func (self *RedisPlayerIntimacyObj) SetIntimacyLv(lv int32) {
	self.SetOneValue("IntimacyLv", lv)
}

//获取是否sayHi
func (self *RedisPlayerIntimacyObj) GetIsSayHi() bool {
	var val bool
	self.GetOneValue("IsSayHi", &val)
	return val
}

//设置是否sayHi
func (self *RedisPlayerIntimacyObj) SetIIsSayHi(b bool) {
	self.SetOneValue("IsSayHi", b)
	self.SaveToMongo()
}

//增加亲密增等级
func (self *RedisPlayerIntimacyObj) AddIntimacyLv(lv int32) int32 {
	return int32(self.IncrOneValue("IntimacyLv", int64(lv)))
}

//获取上一次亲密度增加时间
func (self *RedisPlayerIntimacyObj) GetLastTime() int64 {
	var val int64
	self.GetOneValue("LastTime", &val)
	return val
}

//设置亲密增加时间
func (self *RedisPlayerIntimacyObj) SetLastTime(t int64) {
	self.SetOneValue("LastTime", t)
}

//处理亲密度变化逻辑
func (self *RedisPlayerIntimacyObj) DealAddIntimacyVal(v int64) {
	var IsLastLv bool //是否达到了最高等级
	lv := self.GetIntimacyLv()
	config := GetConfigIntimacy(lv)
	if config == nil {
		logs.Error("数据异常，找不到亲密度配置:", lv)
		return
	}
	if v > 0 {
		//表示亲密度增加
		self.SetLastTime(time.Now().Unix())
		if lv == PLAYER_INTIMACY_MAX_LV {
			IsLastLv = true
			if config.GetMaxVal() == self.GetIntimacyVal() {
				logs.Error("玩家亲密度已满:", self.Id, self.GetIntimacyVal())
				return
			}
		}
	}
	val := self.AddIntimacyVal(v)
	if val > 0 {
		needVal := val - config.GetMaxVal()
		if needVal >= 0 {
			if !IsLastLv {
				//升级了
				self.AddIntimacyLv(1)
				self.SetIntimacyVal(needVal)
			} else {
				//经验值满了，但等级已达最高级，无法再升级
				self.SetIntimacyVal(config.GetMaxVal())
			}
			self.SaveToMongo()
		}
	} else {
		//降级了
		newLv := self.AddIntimacyLv(-1)
		if newLv < 0 {
			self.SetIntimacyLv(0)
			self.SetIntimacyVal(0)
		} else {
			newConfig := GetConfigIntimacy(int32(newLv))
			if newConfig == nil {
				logs.Error("数据异常，找不到亲密度配置:", newLv)
				return
			}
			self.SetIntimacyVal(newConfig.GetMaxVal())
		}
		self.SaveToMongo()
	}
}

//超过3天联系，每天亲密度值减少
func (self *RedisPlayerIntimacyObj) PerDayReduce() {
	lv := self.GetIntimacyLv()
	config := GetConfigIntimacy(lv)
	if config == nil {
		logs.Error("亲密度配置数据异常", lv)
		return
	}
	reduceVal := config.GetMaxVal() * int64(config.GetPerDayVal()) / 100
	newVal := self.GetIntimacyVal() - reduceVal
	if newVal < 0 {
		lv -= 1
		if lv < 0 {
			lv = 0
			newVal = 0
		} else { //降级后，val值等于当前级最大值
			newConfig := GetConfigIntimacy(lv)
			newVal = newConfig.GetMaxVal()
		}
	}
	self.SetIntimacyLv(lv)
	self.SetIntimacyVal(newVal)
}

//删除亲密度记录
func (self *RedisPlayerIntimacyObj) CleanPlayerIntimacy() {
	//1 redis删除key
	self.SetSaveStatus(false)
	self.DelToExistList(self.Id)
	self.DelRedisKey()                //redis删除
	PlayerIntimacyMgr.Delete(self.Id) // 释放对象
	//2 数据库删除
	CleanPlayerIntimacy(self.Id)
}

//对外接口 ==============================
func SavePlayerIntimacyoMongoDB() {
	ids := []string{}
	GetAllRedisSaveList(TABLE_PLAYER_INTIMACY, &ids)
	logs.Info(" ids:", ids)
	list := make([]*share_message.PlayerIntimacy, 0)
	for _, id := range ids {
		obj := GetRedisPlayerIntimacyObj(id)
		if obj != nil {
			intimacy := obj.GetRedisPlayerIntimacy()
			list = append(list, intimacy)
			obj.SetSaveStatus(false)
		}
	}
	logs.Info("要存储的数据:", list)
	if len(list) > 0 {
		saveData := make([]interface{}, 0)
		for _, intimacy := range list {
			saveData = append(saveData, bson.M{"_id": intimacy.GetId()}, intimacy)
		}
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_INTIMACY, saveData)
	}
}

func GetRedisPlayerIntimacyObj(id string) *RedisPlayerIntimacyObj {
	obj := PlayerIntimacyMgr.GetRedisPlayerIntimacyObj(id)
	return obj
}

//亲密度增加
func AddPlayerIntimacy(id string, val int64) {
	obj := PlayerIntimacyMgr.GetRedisPlayerIntimacyObj(id)
	if obj == nil {
		return
	}
	obj.DealAddIntimacyVal(val)
}
