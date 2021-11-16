package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

/*
redis数据管理基类模板
*/

const REDIS_SAVE_TIME = 600 * time.Second //保存时间

const REDIS_SAVE_KEY_LIST = "save_keys"   //需要保存的keys集合
const REDIS_EXIST_KEY_LIST = "exist_keys" //需要保存的keys集合

type IRedisBase interface {
	GetId() interface{}
	GetKeyId() string
	UpdateData()
	GetRedisSaveData() interface{}
	InitRedis()
	SaveOtherData() //保存其他数据
}
type RedisBase struct {
	Me         IRedisBase
	DB         easygo.IMongoDBManager
	DBName     string //mongo数据库名
	TBName     string //表名 redis key名并且mongo表名
	ExistKey   string //reidis现存key值d
	SaveKeys   string //存储列表
	CreateTime int64  //对象创建时间
	//SaveStatus bool         //存储状态:true，需要存储，false不需要存储
	Mutex   easygo.RLock //局部数据锁,只对同一服务器内有用
	IncrId  int64
	Sid     int32
	IsCheck bool
}

func (self *RedisBase) Init(me IRedisBase, id interface{}, db easygo.IMongoDBManager, dbName, tbName string) {
	self.DB = db
	self.DBName = dbName
	self.TBName = tbName
	self.CreateTime = GetMillSecond()
	//self.SaveStatus = false //默认不需要存储
	self.IsCheck = true
	self.Me = me
	self.SaveKeys = MakeRedisKey(REDIS_SAVE_KEY_LIST, tbName)
	self.ExistKey = MakeRedisKey(REDIS_EXIST_KEY_LIST, tbName)
}

//增加现存keys
func (self *RedisBase) AddToExistList(id interface{}) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	err := easygo.RedisMgr.GetC().HSet(self.ExistKey, easygo.AnytoA(id), self.Sid)
	easygo.PanicError(err)
}

//删除key
func (self *RedisBase) DelToExistList(id interface{}) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	_, err := easygo.RedisMgr.GetC().Hdel(self.ExistKey, easygo.AnytoA(id))
	easygo.PanicError(err)
}

//增加到需要保存列表
func (self *RedisBase) AddToSaveList() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	err := easygo.RedisMgr.GetC().SAdd(self.SaveKeys, self.Me.GetId())
	easygo.PanicError(err)
}

//从需保存列表删除key
func (self *RedisBase) DelToSaveList() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	err := easygo.RedisMgr.GetC().SRem(self.SaveKeys, self.Me.GetId())
	easygo.PanicError(err)
}

//获取所有要保存的Keys
func GetAllRedisSaveList(key string, data interface{}) {
	val, err := easygo.RedisMgr.GetC().Smembers(MakeRedisKey(REDIS_SAVE_KEY_LIST, key))
	easygo.PanicError(err)
	switch value := data.(type) {
	case *[]string:
		InterfersToStrings(val, value)
		data = value
	case *[]int64:
		InterfersToInt64s(val, value)
		data = value
	case *[]int32:
		InterfersToInt32s(val, value)
		data = value
	default:
		panic("找不到类型，请自行定义")
	}
}
func GetAllRedisExistList(key string, data interface{}) {
	val, err := easygo.RedisMgr.GetC().HKeys(MakeRedisKey(REDIS_EXIST_KEY_LIST, key))
	easygo.PanicError(err)
	switch value := data.(type) {
	case *[]string:
		InterfersToStrings(val, value)
		data = value
	case *[]int64:
		InterfersToInt64s(val, value)
		data = value
	case *[]int32:
		InterfersToInt32s(val, value)
		data = value
	default:
		panic("找不到类型，请自行定义")
	}
}

func (self *RedisBase) SetSaveSid() {
	self.AddToExistList(self.Me.GetId())
}
func (self *RedisBase) CheckIsDelRedisKey() bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if !self.IsExistKey() {
		return false
	}
	sid, err := redis.Int(easygo.RedisMgr.GetC().HGet(self.ExistKey, easygo.AnytoA(self.Me.GetId())))
	if err != nil {
		logs.Error("CheckIsDelRedisKey 报错:", err.Error())
		return false
	}

	//logs.Info("key:", sid, self.Sid)
	return int32(sid) == self.Sid
}

//设置数据是否需要存储
func (self *RedisBase) SetSaveStatus(b bool) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if b {
		exist := easygo.RedisMgr.GetC().SIsMember(self.SaveKeys, self.Me.GetId())
		if !exist {
			self.AddToSaveList()
		}
	} else {
		self.DelToSaveList()
	}
}

//存储
func (self *RedisBase) GetSaveStatus() bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	b := easygo.RedisMgr.GetC().SIsMember(self.SaveKeys, self.Me.GetId())
	return b
}

//指定值增加
func (self *RedisBase) IncrOneValue(key string, val int64) int64 {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if !self.IsExistKey() {
		//如果key值不存在，先获取
		self.Me.InitRedis()
	}
	res := easygo.RedisMgr.GetC().HIncrBy(self.Me.GetKeyId(), key, val)
	self.SetSaveStatus(true)
	self.CreateTime = GetMillSecond()
	self.SetSaveSid()
	return res

}

//修改指定值
func (self *RedisBase) SetOneValue(key string, val interface{}) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if !self.IsExistKey() {
		//如果key值不存在，先获取
		self.Me.InitRedis()
	}
	err := easygo.RedisMgr.GetC().HSet(self.Me.GetKeyId(), key, val)
	easygo.PanicError(err)
	self.SetSaveStatus(true)
	self.CreateTime = GetMillSecond()
	self.SetSaveSid()
}

//获取指定值
func (self *RedisBase) GetOneValue(key string, val interface{}) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if !self.IsExistKey() {
		//如果key值不存在，先获取
		self.Me.InitRedis()
	}
	res, err := easygo.RedisMgr.GetC().HMGet(self.Me.GetKeyId(), key)
	easygo.PanicError(err)
	_, err = redis.Scan(res, val)
	self.CreateTime = GetMillSecond()
	self.SetSaveSid()

}

//ridis存储指定hash值
func (self *RedisBase) SetStringValueToRedis(name string, data interface{}, save ...bool) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if !self.IsExistKey() {
		//如果key值不存在，先获取
		self.Me.InitRedis()
	}
	//logs.Info("设置玩家数据:", name, data)
	val, err := json.Marshal(data)
	err = easygo.RedisMgr.GetC().HSet(self.Me.GetKeyId(), name, string(val))
	easygo.PanicError(err)
	isSave := append(save, true)[0]
	if isSave {
		self.SetSaveStatus(true)
	}
	self.CreateTime = GetMillSecond()
	self.SetSaveSid()
}

//redis获取指定hash值
func (self *RedisBase) GetStringValueToRedis(name string, data interface{}) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if !self.IsExistKey() {
		//如果key值不存在，先获取
		self.Me.InitRedis()
	}
	val, err := easygo.RedisMgr.GetC().HGet(self.Me.GetKeyId(), name)
	if len(val) == 0 {
		return
	}
	easygo.PanicError(err)
	err = json.Unmarshal(val, &data)
	easygo.PanicError(err)
	self.CreateTime = GetMillSecond()
	self.SetSaveSid()
}

//从mongo中读取数据
func (self *RedisBase) QueryMongoData(id interface{}) interface{} {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	col, closeFun := self.DB.GetC(self.DBName, self.TBName)
	defer closeFun()

	var data interface{}
	err := col.Find(bson.M{"_id": id}).One(&data)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return data
}

//检测是否存在reids数据
func (self *RedisBase) IsExistKey() bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if self.IsCheck {
		res, err := easygo.RedisMgr.GetC().Exist(self.Me.GetKeyId())
		easygo.PanicError(err)
		return res
	}
	return true
}

//全部数据写到mongo
func (self *RedisBase) SaveToMongo() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	if self.IsExistKey() { //只有redis 存在key才进行存储
		data := self.Me.GetRedisSaveData()
		if data == nil {
			logs.Error("SaveToMongo 存储时获取到数据为空:", self.Me.GetKeyId())
			return
		}
		self.Me.SaveOtherData()
		//logs.Info("redis 存储 ：", self.Me.GetKeyId(), data)

		fun := func() {
			col, closeFun := self.DB.GetC(self.DBName, self.TBName)
			defer closeFun()
			_, err := col.Upsert(bson.M{"_id": self.Me.GetId()}, data)
			//logs.Info("存储redis模块数据:", self.Me.GetKeyId(), self.Me.GetId(), data)
			easygo.PanicError(err)
		}
		easygo.Spawn(fun)
		self.SetSaveStatus(false)
	}

}

//保存单个字段
func (self *RedisBase) SaveOneRedisDataToMongo(file string, value interface{}) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	col, closeFun := self.DB.GetC(self.DBName, self.TBName)
	defer closeFun()
	err := col.Update(bson.M{"_id": self.Me.GetId()}, bson.M{"$set": bson.M{file: value}})
	if err != nil {
		logs.Error("err:", err, file, value)
	}
}

//把key从redis删除
func (self *RedisBase) DelRedisKey() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	_, err := easygo.RedisMgr.GetC().Delete(self.Me.GetKeyId())
	easygo.PanicError(err)
}

//对外通用方法
func DelAllKeyFromRedis(dbName, existName string, val interface{}) {
	var delKeys []interface{}
	switch ids := val.(type) {
	case []int64:
		for _, id := range ids {
			delKeys = append(delKeys, MakeRedisKey(dbName, id))
		}
	case []string:
		for _, id := range ids {
			delKeys = append(delKeys, MakeRedisKey(dbName, id))
		}
	}
	delKeys = append(delKeys, existName)
	if len(delKeys) > 0 {
		_, err := easygo.RedisMgr.GetC().Delete(delKeys...)
		easygo.PanicError(err)
	}
}

//:TODO 后续优化分布式锁
func (self *RedisBase) Lock() {
	_, err := easygo.RedisMgr.GetC().Delete(self.Me.GetKeyId())
	easygo.PanicError(err)
	logs.Info("删除redis key:", self.Me.GetKeyId())
}
func (self *RedisBase) UnLock() {
	_, err := easygo.RedisMgr.GetC().Delete(self.Me.GetKeyId())
	easygo.PanicError(err)
	logs.Info("删除redis key:", self.Me.GetKeyId())
}
func (self *RedisBase) CheckLock() {
	_, err := easygo.RedisMgr.GetC().Delete(self.Me.GetKeyId())
	easygo.PanicError(err)
	logs.Info("删除redis key:", self.Me.GetKeyId())
}
