package for_game

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

type RedisESportTableFedUpdateObj struct {
	Table string
	RedisBase
}
type RedisESportTableFedUpdate struct {
	Table string
	Fed   string
	Value int64
	Id    int64
}

const REDIS_ESPORT_TABLE_FED_UPDATE_TIME = 1000 * 600 //毫秒，key值存在时间
const REDIS_ESPORT_FED_SAVE_TIME = 1000 * 300         //5分钟存一次

var REDIS_ESPORT_TABLE_FED_UPDAT_KEY = ESportExN("fed_update") //redis内存中存在的key

func (self *RedisESportTableFedUpdateObj) Init(table string) *RedisESportTableFedUpdateObj {
	self.Table = table
	self.RedisBase.Init(self, self.GetId(), easygo.MongoMgr, MONGODB_NINGMENG_LOG, REDIS_ESPORT_TABLE_FED_UPDAT_KEY)
	self.Sid = ESportTableFedUpdateMgr.GetSid()
	self.AddToExistList(self.GetId())
	ESportTableFedUpdateMgr.Store(self.GetId(), self)
	easygo.AfterFunc(REDIS_ESPORT_FED_SAVE_TIME, self.UpdateData)
	return self
}
func GetTableFedUpdateId(fed string, id int64) string {
	return fmt.Sprintf("%s_%d", fed, id)
}
func (self *RedisESportTableFedUpdateObj) GetId() interface{} { //override
	return self.Table
}
func (self *RedisESportTableFedUpdateObj) GetKeyId() string { //override
	return MakeRedisKey(REDIS_ESPORT_TABLE_FED_UPDAT_KEY, self.GetId())
}
func (self *RedisESportTableFedUpdateObj) GetItemKeyId(fed string, id int64) string { //override
	return GetTableFedUpdateId(fed, id)
}
func (self *RedisESportTableFedUpdateObj) UpdateData() { //override

	self.SaveToMongo()

	if self.CreateTime > REDIS_ESPORT_TABLE_FED_UPDATE_TIME { //用户退出的时候控制
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.GetId())
			self.DelRedisKey() //redis删除
		}
		ESportTableFedUpdateMgr.Delete(self.GetId()) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_ESPORT_FED_SAVE_TIME, self.UpdateData)
}

///某个int32的字段 +1 -1
func (self *RedisESportTableFedUpdateObj) UpdateFedAddition(table string, fed string, dataId int64, count int64) bool {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, table)
	defer closeFun()
	_, err := col.Find(bson.M{"_id": dataId}).Apply(mgo.Change{
		Update: bson.M{"$inc": bson.M{fed: count}},
		Upsert: true,
	}, nil)
	if err != nil {
		logs.Error(err)
	}
	return err == nil
}

//重写保存方法
func (self *RedisESportTableFedUpdateObj) SaveToMongo() {

	lst := self.GetTableItem()
	for _, it := range lst {
		if it.Id > 0 && it.Value > 0 {
			b := self.UpdateFedAddition(it.Table, it.Fed, it.Id, it.Value)
			if b {
				logs.Info("累加:", it)
				self.DeleteFed(it.Fed, it.Id)
			}
		}
	}
	self.SetSaveStatus(false)
}

//刪除键
func (self *RedisESportTableFedUpdateObj) DeleteFed(fed string, id int64) {
	_, err := easygo.RedisMgr.GetC().Hdel(self.GetKeyId(), self.GetItemKeyId(fed, id))
	easygo.PanicError(err)
}

//当检测到redis key不存在时
func (self *RedisESportTableFedUpdateObj) InitRedis() { //override

}
func (self *RedisESportTableFedUpdateObj) GetRedisSaveData() interface{} { //override
	return nil
}
func (self *RedisESportTableFedUpdateObj) SaveOtherData() { //override

}

func (self *RedisESportTableFedUpdateObj) GetTableItem() map[string]*RedisESportTableFedUpdate {

	mps := self.GetTableItemStr()
	savelist := make(map[string]*RedisESportTableFedUpdate)
	for k, v := range mps {
		it := &RedisESportTableFedUpdate{}
		_ = json.Unmarshal([]byte(v), &it)
		savelist[k] = it
	}
	return savelist
}

func (self *RedisESportTableFedUpdateObj) GetTableItemStr() map[string]string {
	values, err := StrkeyStringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	return values
}

//累加数据
func (self *RedisESportTableFedUpdateObj) RvpValue(fed string, id, value int64) {
	ky := self.GetItemKeyId(fed, id)
	var val string
	self.GetOneValue(ky, &val)
	if val != "" {
		it := &RedisESportTableFedUpdate{}
		_ = json.Unmarshal([]byte(val), &it)
		if it.Id > 0 {
			it.Value += value
			s, _ := json.Marshal(it)
			self.SetOneValue(ky, string(s))
		}
	} else {
		it := &RedisESportTableFedUpdate{
			Table: self.Table,
			Fed:   fed,
			Value: value,
			Id:    id,
		}
		s, _ := json.Marshal(it)
		self.SetOneValue(ky, string(s))
	}

}

func NewRedisESportTableFedUpdateObj(table string) *RedisESportTableFedUpdateObj {
	obj := RedisESportTableFedUpdateObj{}
	return obj.Init(table)
}

//对外方法，获取对象，如果为nil表示redis内存不存在，数据库也不存在
func GetRedisESportTableFedUpdateObj(table string) *RedisESportTableFedUpdateObj {
	obj, ok := ESportTableFedUpdateMgr.Load(table)
	if ok && obj != nil {
		return obj.(*RedisESportTableFedUpdateObj)
	} else {
		return NewRedisESportTableFedUpdateObj(table)
	}
}
