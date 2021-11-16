package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"
	"sync"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

var ESPORT_REDIS_PLAYER_KEY = ESportExN("player") //电竞玩家

const (
	ESPORT_REDIS_PLAYER_EXIST_TIME = 1000 * 600 //redis的key删除时间:毫秒
)

type RedisESportPlayerObj struct {
	PlayerId PLAYER_ID
	RedisBase
	BpsDuration sync.Map
}

type RedisESportPlayer struct {
	Id int64 `json:"_id"`
	//状态
	Status int32
	//上次拉取消息的时间
	LastPullTime int64
	//当前所在放映厅的ID
	CurrentRoomLiveId int64
	LastLoginTime     int64
	CreateTime        int64
}

func NewRedisESportPlayerObj(playerId PLAYER_ID, data ...*share_message.TableESPortsPlayer) *RedisESportPlayerObj {
	p := &RedisESportPlayerObj{
		PlayerId: playerId,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}
func (self *RedisESportPlayerObj) Init(obj *share_message.TableESPortsPlayer) *RedisESportPlayerObj {
	self.RedisBase.Init(self, self.PlayerId, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_ESPORTS_PLAYER)
	self.Sid = ESportPlayerMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		ESportPlayerMgr.Store(self.PlayerId, self)
		self.AddToExistList(self.PlayerId)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryESPortPlayer(self.PlayerId)

			if obj == nil {
				return nil
			}

		}
		self.SetRedisESportPlayer(obj)
		self.SetLastLoginTime(easygo.NowTimestamp())
	}
	logs.Info("初始化新的电竞用户管理器", self.PlayerId)
	return self
}
func (self *RedisESportPlayerObj) GetId() interface{} { //override
	return self.PlayerId
}
func (self *RedisESportPlayerObj) GetKeyId() string { //override
	return MakeRedisKey(ESPORT_REDIS_PLAYER_KEY, self.PlayerId)
}

//定时更新数据
func (self *RedisESportPlayerObj) UpdateData() { //override
	if !self.IsExistKey() {
		ESportPlayerMgr.Delete(self.PlayerId) // 释放对象
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > ESPORT_REDIS_PLAYER_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.PlayerId)
			self.DelRedisKey() //redis删除
		}
		ESportPlayerMgr.Delete(self.PlayerId) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}

//重写保存方法
func (self *RedisESportPlayerObj) SaveToMongoEx() {

	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	data := self.Me.GetRedisSaveData()
	if data == nil {
		logs.Error("SaveToMongoEx 存储时获取到数据为空:", self.Me.GetKeyId())
		return
	}
	col, closeFun := self.DB.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_PLAYER)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": self.Me.GetId()}, data)
	//logs.Info("存储redis模块数据:", self.Me.GetKeyId(), self.Me.GetId(), data)
	easygo.PanicError(err)
	self.DelRedisKey()
}

func (self *RedisESportPlayerObj) InitRedis() { //override
	obj := self.QueryESPortPlayer(self.PlayerId)
	if obj == nil {
		return
	}
	self.SetRedisESportPlayer(obj)
	//重新激活定时器
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisESportPlayerObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisESPortsPlayer()
	return data
}
func (self *RedisESportPlayerObj) SaveOtherData() { //override
}

//通过playerId从mongo中读取登录玩家数据
func (self *RedisESportPlayerObj) QueryESPortPlayer(id PLAYER_ID) *share_message.TableESPortsPlayer {
	data := self.QueryMongoData(id)
	if data != nil {
		var account share_message.TableESPortsPlayer
		StructToOtherStruct(data, &account)
		return &account
	}
	return nil
}

//设置玩家账号信息
func (self *RedisESportPlayerObj) SetRedisESportPlayer(obj *share_message.TableESPortsPlayer) {
	ESportPlayerMgr.Store(obj.GetId(), self)
	self.AddToExistList(obj.GetId())
	///重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	account := &RedisESportPlayer{}
	StructToOtherStruct(obj, account)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), account)
	easygo.PanicError(err)

}
func (self *RedisESportPlayerObj) GetRedisESPortsPlayer() *share_message.TableESPortsPlayer {
	obj := &RedisESportPlayer{}
	key := self.GetKeyId()
	value, err := easygo.RedisMgr.GetC().HGetAll(key)
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.TableESPortsPlayer{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

func (self *RedisESportPlayerObj) SetLastLoginTime(t int64) {

	self.SetOneValue("LastLoginTime", t)
}

func (self *RedisESportPlayerObj) SetLastPullTime(t int64) {

	self.SetOneValue("LastPullTime", t)
}
func (self *RedisESportPlayerObj) GetLastPullTime() int64 {
	var val int64
	self.GetOneValue("LastPullTime", &val)
	return val
}
func (self *RedisESportPlayerObj) SetCurrentRoomLiveId(id int64) {
	logs.Info("设置用户：%d进入直播间:%d", self.PlayerId, id)
	self.SetOneValue("CurrentRoomLiveId", id)
}
func (self *RedisESportPlayerObj) GetCurrentRoomLiveId() int64 {
	var val int64
	self.GetOneValue("CurrentRoomLiveId", &val)
	return val
}
func (self *RedisESportPlayerObj) SetCreateTime(t int64) {
	self.SetOneValue("CreateTime", t)
}
func (self *RedisESportPlayerObj) GetCreateTime(t int64) int64 {
	var val int64
	self.GetOneValue("CreateTime", &val)
	return val
}

//对外方法，获取玩家对象，如果为nil表示redis内存不存在，数据库也不存在
func GetRedisESportPlayerObj(id PLAYER_ID) *RedisESportPlayerObj {
	obj, ok := ESportPlayerMgr.Load(id)
	if ok && obj != nil {
		return obj.(*RedisESportPlayerObj)
	} else {
		return NewRedisESportPlayerObj(id)
	}
}

//停服保存处理，保存需要存储的数据
func SaveRedisESportPlayerToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(ESPORT_REDIS_PLAYER_KEY, &ids)
	//修改
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisESportPlayerObj(id)
		if obj != nil {
			account := obj.GetRedisESPortsPlayer()
			saveData = append(saveData, bson.M{"_id": account.GetId()}, account)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_ESPORTS_PLAYER, saveData)
	}
}
