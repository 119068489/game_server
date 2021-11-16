package for_game

import (
	b64 "encoding/base64"
	"encoding/json"
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

/*
转账内存数据管理
*/

const (
	//列表key
	TRANSFERMONEY_EXIST_LIST = "transfer_exist_lists" //修改的转账列表
	TRANSFERMONEY_EXIST_TIME = 1000 * 600             //redis的key删除时间:毫秒
)

type RedisTransferMoneyObj struct {
	Id int64
	RedisBase
}
type RedisTransferMoney struct {
	Id         int64  `json:"_id"`
	Sender     int64  //转账人的id
	TargetId   int64  //接收人的id
	Way        int32  //转账方式:99 零钱 1微信 2支付宝 3银行卡
	Card       string //银行卡号，零钱不需要
	Gold       int64  //转账金额
	Content    string //留言
	CreateTime int64  //转账时间
	State      int32  //状态：1未领取，2已领取，3转账主动退款,4已过期(24小时未领取退回)
	OpenTime   int64  //领取或者退还时间
	OrderId    string //交易id
}

//写入redis
func NewRedisTransferMoney(id int64, data ...*share_message.TransferMoney) *RedisTransferMoneyObj {
	p := &RedisTransferMoneyObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}
func (self *RedisTransferMoneyObj) Init(obj *share_message.TransferMoney) *RedisTransferMoneyObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_TRANSFER_MONEY)
	self.Sid = TransferMoneyMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		TransferMoneyMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryTransferMoney()
			if obj == nil {
				return nil
			}
		}
		self.SetRedisTransferMoney(obj)
	}

	logs.Info("初始化新的TransferMoney管理器:", self.Id)
	return self
}
func (self *RedisTransferMoneyObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisTransferMoneyObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_TRANSFER_MONEY, self.Id)
}

//定时更新数据
func (self *RedisTransferMoneyObj) UpdateData() { //override
	if !self.IsExistKey() {
		TransferMoneyMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > TRANSFERMONEY_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		TransferMoneyMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisTransferMoneyObj) InitRedis() { //override
	obj := self.QueryTransferMoney()
	if obj == nil {
		return
	}
	self.SetRedisTransferMoney(obj)
}
func (self *RedisTransferMoneyObj) GetRedisSaveData() interface{} { //override
	data := self.GetTransferMoney()
	return data
}
func (self *RedisTransferMoneyObj) SaveOtherData() { //override
}
func (self *RedisTransferMoneyObj) QueryTransferMoney() *share_message.TransferMoney {
	data := self.QueryMongoData(self.Id)
	if data != nil {
		var transfer share_message.TransferMoney
		StructToOtherStruct(data, &transfer)
		return &transfer
	}
	return nil
}
func (self *RedisTransferMoneyObj) SetRedisTransferMoney(obj *share_message.TransferMoney) {
	//增加到管理器
	TransferMoneyMgr.Store(obj.GetId(), self)
	self.AddToExistList(obj.GetId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	transfer := &RedisTransferMoney{}
	StructToOtherStruct(obj, transfer)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), transfer)
	easygo.PanicError(err)

}
func (self *RedisTransferMoneyObj) GetTransferMoney() *share_message.TransferMoney {
	obj := &RedisTransferMoney{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.TransferMoney{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//获取传输给前端的字符串结构
func (self *RedisTransferMoneyObj) GetTransferMoneyBase64Str() string {
	data := self.GetTransferMoney()
	content, err := json.Marshal(data)
	easygo.PanicError(err)
	sEnc := b64.StdEncoding.EncodeToString(content)
	return sEnc
}
func (self *RedisTransferMoneyObj) GetSender() int64 {
	var val int64
	self.GetOneValue("Sender", &val)
	return val
}
func (self *RedisTransferMoneyObj) GetTargetId() int64 {
	var val int64
	self.GetOneValue("TargetId", &val)
	return val
}
func (self *RedisTransferMoneyObj) GetGold() int64 {
	var val int64
	self.GetOneValue("Gold", &val)
	return val
}
func (self *RedisTransferMoneyObj) GetContent() string {
	var val string
	self.GetOneValue("Content", &val)
	return val
}

func (self *RedisTransferMoneyObj) GetWay() int32 {
	var val int32
	self.GetOneValue("Way", &val)
	return val
}
func (self *RedisTransferMoneyObj) GetState() int32 {
	var val int32
	self.GetOneValue("State", &val)
	return val
}
func (self *RedisTransferMoneyObj) SetState(val int32) {
	self.SetOneValue("State", val)
}
func (self *RedisTransferMoneyObj) SetOpenTime(val int64) {
	self.SetOneValue("OpenTime", val)
}

//对外方法
func GetRedisTransferMoneyObj(id int64) *RedisTransferMoneyObj {
	return TransferMoneyMgr.GetRedisTransferMoneyObj(id)
}
func GetOpenTransferBase64Str(msg *client_hall.OpenTransfer) string {
	content, err := json.Marshal(msg)
	easygo.PanicError(err)
	sEnc := b64.StdEncoding.EncodeToString(content)
	return sEnc
}

//发转账
func CreateTransferMoney(msg *share_message.TransferMoney, orderId string) *RedisTransferMoneyObj {
	obj := TransferMoneyMgr.CreateTransferMoney(msg, orderId)
	obj.SetSaveStatus(true)
	return obj
}

func GetAllExistTransferList() []int64 {
	ids := []int64{}
	b, err := easygo.RedisMgr.GetC().Exist(TRANSFERMONEY_EXIST_LIST)
	if !b {
		return ids
	}
	val, err := easygo.RedisMgr.GetC().Smembers(TRANSFERMONEY_EXIST_LIST)
	if len(val) > 0 && err != nil {
		easygo.PanicError(err)
	}
	InterfersToInt64s(val, &ids)
	return ids
}

//========================================================================<

//检测玩家是否有发送未领取完的红包
func CheckPlayerTransferMoney(playerId PLAYER_ID) bool {
	//redis检测
	ids := GetAllExistTransferList()
	for _, id := range ids {
		transfer := GetRedisTransferMoneyObj(id)
		if transfer.GetSender() == playerId && transfer.GetState() == TRANSFER_MONEY_OPEN {
			return false
		}
	}
	//数据库检测
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TRANSFER_MONEY)
	defer closeFun()
	queryBson := bson.M{"Sender": playerId, "State": TRANSFER_MONEY_OPEN}
	lst := []*share_message.TransferMoney{}
	err := col.Find(queryBson).All(&lst)
	if err != nil && err != mgo.ErrNotFound {
		return false
	}
	if len(lst) > 0 {
		return false
	}
	return true
}

//批量保存需要存储的数据
func SaveRedisTransferMoneyToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_TRANSFER_MONEY, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisTransferMoneyObj(id)
		if obj != nil {
			data := obj.GetTransferMoney()
			saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_TRANSFER_MONEY, saveData)
	}
}
