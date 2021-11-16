package for_game

import (
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo/bson"
	"github.com/garyburd/redigo/redis"
)

//玩家装备详情
const REDIS_PLAYER_EQUIPMENT_EXIST_LIST = "player_equipment_exist_list"
const REDIS_PLAYER_EQUIPMENT_EXIST_TIME = 1000 * 600

const (
	EQUIPMENT_UP   = 1
	EQUIPMENT_DOWN = 2
)

type RedisPlayerEquipmentObj struct {
	Id int64 //玩家id
	RedisBase
}
type PlayerEquipmentEx struct {
	PlayerId int64
	GJ       int64
	QP       int64
	MP       int64
	QTX      int64
	MZBS     int64
}

//写入redis
func NewRedisPlayerEquipment(id PLAYER_ID, equipment ...*share_message.PlayerEquipment) *RedisPlayerEquipmentObj {
	p := &RedisPlayerEquipmentObj{
		Id: id,
	}
	obj := append(equipment, nil)[0]
	return p.Init(obj)
}

func (self *RedisPlayerEquipmentObj) Init(obj *share_message.PlayerEquipment) *RedisPlayerEquipmentObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_EQUIPMENT)
	self.Sid = PlayerEquipmentMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		PlayerEquipmentMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//self.SetSaveStatus(true)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryPlayerEquipment(self.Id)
			//if obj == nil {
			//	return nil
			//}
		}
		self.SetRedisPlayerEquipment(obj)
	}
	return self
}
func (self *RedisPlayerEquipmentObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisPlayerEquipmentObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_PLAYER_EQUIPMENT, self.Id)
}

//定时更新数据
func (self *RedisPlayerEquipmentObj) UpdateData() { //override
	if !self.IsExistKey() {
		PlayerEquipmentMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDIS_PLAYER_EQUIPMENT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.SaveToMongo() //删除前保存一次数据
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		PlayerEquipmentMgr.Delete(self.Id) // 释放对象
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisPlayerEquipmentObj) InitRedis() { //override
	obj := self.QueryPlayerEquipment(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisPlayerEquipment(obj)
}
func (self *RedisPlayerEquipmentObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisPlayerEquipment()
	return data
}

func (self *RedisPlayerEquipmentObj) SaveOtherData() { //override
}
func (self *RedisPlayerEquipmentObj) QueryPlayerEquipment(id int64) *share_message.PlayerEquipment {
	data := self.QueryMongoData(id)
	if data != nil {
		var equipment share_message.PlayerEquipment
		StructToOtherStruct(data, &equipment)
		return &equipment
	}
	//redis中获取数据初始化
	if self.IsExistKey() {
		return self.GetRedisPlayerEquipment()
	}
	return nil
}
func (self *RedisPlayerEquipmentObj) SetRedisPlayerEquipment(equipment *share_message.PlayerEquipment) {
	PlayerEquipmentMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	eq := &PlayerEquipmentEx{
		PlayerId: self.Id,
		GJ:       0,
		QP:       0,
		MP:       0,
		QTX:      0,
		MZBS:     0,
	}
	if equipment != nil {
		StructToOtherStruct(equipment, eq)
	}
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), eq)
	easygo.PanicError(err)
}
func (self *RedisPlayerEquipmentObj) GetRedisPlayerEquipment() *share_message.PlayerEquipment {
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	var ex PlayerEquipmentEx
	err = redis.ScanStruct(value, &ex)
	easygo.PanicError(err)
	equipment := &share_message.PlayerEquipment{}
	StructToOtherStruct(ex, equipment)
	return equipment
}

//获取挂件
func (self *RedisPlayerEquipmentObj) GetGJ() int64 {
	var val int64
	self.GetOneValue("GJ", &val)
	return val
}

func (self *RedisPlayerEquipmentObj) SetGJ(id int64) {
	self.SetOneValue("GJ", id)
}

//获取气泡
func (self *RedisPlayerEquipmentObj) GetQP() int64 {
	var val int64
	self.GetOneValue("QP", &val)
	return val
}

func (self *RedisPlayerEquipmentObj) SetQP(id int64) {
	self.SetOneValue("QP", id)
}

//获取铭牌
func (self *RedisPlayerEquipmentObj) GetMP() int64 {
	var val int64
	self.GetOneValue("MP", &val)
	return val
}

func (self *RedisPlayerEquipmentObj) SetMP(id int64) {
	self.SetOneValue("MP", id)
}

//获取群特效
func (self *RedisPlayerEquipmentObj) GetQTX() int64 {
	var val int64
	self.GetOneValue("QTX", &val)
	return val
}

func (self *RedisPlayerEquipmentObj) SetQTX(id int64) {
	self.SetOneValue("QTX", id)
}

//获取名字变色
func (self *RedisPlayerEquipmentObj) GetMZBS() int64 {
	var val int64
	self.GetOneValue("MZBS", &val)
	return val
}

func (self *RedisPlayerEquipmentObj) SetMZBS(id int64) {
	self.SetOneValue("MZBS", id)
}

//获取当前装备
func (self *RedisPlayerEquipmentObj) GetCurEquipment(pos int32) int64 {
	oldId := int64(0)
	switch pos {
	case COIN_PROPS_TYPE_GJ:
		oldId = self.GetGJ()
	case COIN_PROPS_TYPE_QP:
		oldId = self.GetQP()
	case COIN_PROPS_TYPE_MP:
		oldId = self.GetMP()
	case COIN_PROPS_TYPE_QTX:
		oldId = self.GetQTX()
	case COIN_PROPS_TYPE_MZBS:
		oldId = self.GetMZBS()
	}
	return oldId
}

//卸下装备
func (self *RedisPlayerEquipmentObj) EquipmentDown(pos int32) {
	switch pos {
	case COIN_PROPS_TYPE_GJ:
		self.SetGJ(0)
	case COIN_PROPS_TYPE_QP:
		self.SetQP(0)
	case COIN_PROPS_TYPE_MP:
		self.SetMP(0)
	case COIN_PROPS_TYPE_QTX:
		self.SetQTX(0)
	case COIN_PROPS_TYPE_MZBS:
		self.SetMZBS(0)
	}
}

//使用装备 pos 部位
func (self *RedisPlayerEquipmentObj) Equipment(pos int32, id int64) {
	switch pos {
	case COIN_PROPS_TYPE_GJ:
		self.SetGJ(id)
	case COIN_PROPS_TYPE_QP:
		self.SetQP(id)
	case COIN_PROPS_TYPE_MP:
		self.SetMP(id)
	case COIN_PROPS_TYPE_QTX:
		self.SetQTX(id)
	case COIN_PROPS_TYPE_MZBS:
		self.SetMZBS(id)
	}
}

//获取传给前端的装备配置
func (self *RedisPlayerEquipmentObj) GetEquipmentForClient() *client_hall.EquipmentReq {
	bagObj := GetRedisPlayerBagItemObj(self.Id)
	equipment := &client_hall.EquipmentReq{
		Id: easygo.NewInt64(self.Id),
	}
	gj := bagObj.GetItemNetId(self.GetGJ())
	if gj != nil && gj.GetStatus() == COIN_BAG_ITEM_USED {
		equipment.GJ = &client_hall.Equipment{
			BagId:   easygo.NewInt64(gj.GetId()),
			PropsId: easygo.NewInt64(gj.GetPropsId()),
		}
	}
	qp := bagObj.GetItemNetId(self.GetQP())
	if qp != nil && qp.GetStatus() == COIN_BAG_ITEM_USED {
		equipment.QP = &client_hall.Equipment{
			BagId:   easygo.NewInt64(qp.GetId()),
			PropsId: easygo.NewInt64(qp.GetPropsId()),
		}
	}
	mp := bagObj.GetItemNetId(self.GetMP())
	if mp != nil && mp.GetStatus() == COIN_BAG_ITEM_USED {
		equipment.MP = &client_hall.Equipment{
			BagId:   easygo.NewInt64(mp.GetId()),
			PropsId: easygo.NewInt64(mp.GetPropsId()),
		}
	}
	qtx := bagObj.GetItemNetId(self.GetQTX())
	if qtx != nil && qtx.GetStatus() == COIN_BAG_ITEM_USED {
		equipment.QTX = &client_hall.Equipment{
			BagId:   easygo.NewInt64(qtx.GetId()),
			PropsId: easygo.NewInt64(qtx.GetPropsId()),
		}
	}
	mzbs := bagObj.GetItemNetId(self.GetMZBS())
	if mzbs != nil && mzbs.GetStatus() == COIN_BAG_ITEM_USED {
		equipment.MZBS = &client_hall.Equipment{
			BagId:   easygo.NewInt64(mzbs.GetId()),
			PropsId: easygo.NewInt64(mzbs.GetPropsId()),
		}
	}
	return equipment
}

//对外接口
func GetRedisPlayerEquipmentObj(id int64, data ...*share_message.PlayerEquipment) *RedisPlayerEquipmentObj {
	return PlayerEquipmentMgr.GetRedisPlayerEquipmentObj(id, data...)
}

//停服保存处理，保存需要存储的数据
func SaveRedisPlayerEquipmentToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_PLAYER_EQUIPMENT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisPlayerEquipmentObj(id)
		if obj != nil {
			data := obj.GetRedisPlayerEquipment()
			saveData = append(saveData, bson.M{"_id": data.GetPlayerId()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_EQUIPMENT, saveData)
	}
}
