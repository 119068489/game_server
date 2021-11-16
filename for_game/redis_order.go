package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo/bson"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

/*
订单Redis数据管理
*/

const (
	//现存订单列表
	ORDER_EXIST_TIME = 1000 * 600               //redis的key删除时间:毫秒
	ORDER_DOING_LOCK = "order:order_doing_list" //正在处理订单列表
)

type RedisOrderObj struct {
	Id string
	RedisBase
}
type RedisOrder struct {
	OrderId     string `json:"_id"`
	PlayerId    int64
	Account     string
	NickName    string
	RealName    string
	SourceType  int32
	ChangeType  int32
	Channeltype int32
	CurGold     int64
	ChangeGold  int64
	Gold        int64
	ExternalNo  string
	PayChannel  int32
	PayType     int32
	Amount      int64
	CreateTime  int64
	CreateIP    string
	Status      int32
	PayStatus   int32
	Note        string
	Tax         int64
	Operator    string
	OverTime    int64
	BankInfo    string
	PayWay      int32
	PayTargetId int64
	PayOpenId   string
	TotalCount  int32
	Content     string
	ExtendValue string
	//提现信息
	BankCode    string
	AccountType string
	AccountNo   string
	AccountName string
	AccountProp string
	OrderDate   string
	IsCheck     string
	PlatformTax int64
	RealTax     int64
	OrderType   int32
}

//写入redis
func NewRedisOrder(id string, order ...*share_message.Order) *RedisOrderObj {
	p := &RedisOrderObj{
		Id: id,
	}
	obj := append(order, nil)[0]
	return p.Init(obj)
}

func (self *RedisOrderObj) Init(obj *share_message.Order) *RedisOrderObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_ORDER)
	self.Sid = OrderMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		OrderMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryOrder(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisOrder(obj)
	}
	logs.Info("初始化新的order订单管理器:", self.Id)
	return self
}
func (self *RedisOrderObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisOrderObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_ORDER, self.Id)
}

//定时更新数据
func (self *RedisOrderObj) UpdateData() { //override
	if !self.IsExistKey() {
		OrderMgr.Delete(self.Id) // 释放对象
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > ORDER_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		OrderMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}

//当检测到redis key不存在时
func (self *RedisOrderObj) InitRedis() { //override
	obj := self.QueryOrder(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisOrder(obj)
}
func (self *RedisOrderObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisOrder()
	return data
}
func (self *RedisOrderObj) SaveOtherData() { //override
}
func (self *RedisOrderObj) QueryOrder(id ORDER_ID) *share_message.Order {
	data := self.QueryMongoData(id)
	if data != nil {
		var order share_message.Order
		StructToOtherStruct(data, &order)
		return &order
	}
	return nil
}

//设置订单数据
func (self *RedisOrderObj) SetRedisOrder(obj *share_message.Order) {
	//增加到管理器
	OrderMgr.Store(obj.GetOrderId(), self)
	self.AddToExistList(obj.GetOrderId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	order := &RedisOrder{}
	StructToOtherStruct(obj, order)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), order)
	easygo.PanicError(err)
}

//获取订单数据
func (self *RedisOrderObj) GetRedisOrder() *share_message.Order {
	obj := &RedisOrder{}
	key := self.GetKeyId()
	value, err := easygo.RedisMgr.GetC().HGetAll(key)
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.Order{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置成员值模块
func (self *RedisOrderObj) SetOverTime(val int64) {
	self.SetOneValue("OverTime", val)
}
func (self *RedisOrderObj) SetStatus(val int32) {
	self.SetOneValue("Status", val)
}
func (self *RedisOrderObj) SetPayStatus(val int32) {
	self.SetOneValue("PayStatus", val)
}
func (self *RedisOrderObj) SetChanneltype(val int32) {
	self.SetOneValue("Channeltype", val)
}
func (self *RedisOrderObj) SetExtendValue(val string) {
	self.SetOneValue("ExtendValue", val)
}
func (self *RedisOrderObj) SetExternalNo(val string) {
	self.SetOneValue("ExternalNo", val)
}
func (self *RedisOrderObj) SetNote(val string) {
	self.SetOneValue("Note", val)
}
func (self *RedisOrderObj) SetOperator(val string) {
	self.SetOneValue("Operator", val)
}
func (self *RedisOrderObj) SetSourceType(val int32) {
	self.SetOneValue("SourceType", val)
}
func (self *RedisOrderObj) SetOrderType(val int32) {
	self.SetOneValue("OrderType", val)
}

//获取成员值模块
func (self *RedisOrderObj) GetChangeType() int32 {
	var val int32
	self.GetOneValue("ChangeType", &val)
	return val
}
func (self *RedisOrderObj) GetPayStatus() int32 {
	var val int32
	self.GetOneValue("PayStatus", &val)
	return val
}
func (self *RedisOrderObj) GetStatus() int32 {
	var val int32
	self.GetOneValue("Status", &val)
	return val
}
func (self *RedisOrderObj) GetOverTime() int64 {
	var val int64
	self.GetOneValue("OverTime", &val)
	return val
}
func (self *RedisOrderObj) GetPlayerId() int64 {
	var val int64
	self.GetOneValue("PlayerId", &val)
	return val
}
func (self *RedisOrderObj) GetChangeGold() int64 {
	var val int64
	self.GetOneValue("ChangeGold", &val)
	return val
}
func (self *RedisOrderObj) GetAmount() int64 {
	var val int64
	self.GetOneValue("Amount", &val)
	return val
}
func (self *RedisOrderObj) GetTax() int64 {
	var val int64
	self.GetOneValue("Tax", &val)
	return val
}
func (self *RedisOrderObj) GetPayTargetId() int64 {
	var val int64
	self.GetOneValue("PayTargetId", &val)
	return val
}
func (self *RedisOrderObj) GetPayChannel() int32 {
	var val int32
	self.GetOneValue("PayChannel", &val)
	return val
}
func (self *RedisOrderObj) GetTotalCount() int32 {
	var val int32
	self.GetOneValue("TotalCount", &val)
	return val
}
func (self *RedisOrderObj) GetPayWay() int32 {
	var val int32
	self.GetOneValue("PayWay", &val)
	return val
}
func (self *RedisOrderObj) GetChanneltype() int32 {
	var val int32
	self.GetOneValue("Channeltype", &val)
	return val
}
func (self *RedisOrderObj) GetContent() string {
	var val string
	self.GetOneValue("Content", &val)
	return val
}
func (self *RedisOrderObj) GetCreateIP() string {
	var val string
	self.GetOneValue("CreateIP", &val)
	return val
}
func (self *RedisOrderObj) GetExternalNo() string {
	var val string
	self.GetOneValue("ExternalNo", &val)
	return val
}
func (self *RedisOrderObj) GetOrderId() string {
	return self.Id
}
func (self *RedisOrderObj) GetOperator() string {
	var val string
	self.GetOneValue("Operator", &val)
	return val
}
func (self *RedisOrderObj) GetSourceType() int32 {
	var val int32
	self.GetOneValue("SourceType", &val)
	return val
}
func (self *RedisOrderObj) GetPayType() int32 {
	var val int32
	self.GetOneValue("PayType", &val)
	return val
}
func (self *RedisOrderObj) GetBankInfo() string {
	var val string
	self.GetOneValue("BankInfo", &val)
	return val
}
func (self *RedisOrderObj) GetAccountNo() string {
	var val string
	self.GetOneValue("AccountNo", &val)
	return val
}
func (self *RedisOrderObj) GetCreateTime() int64 {
	var val int64
	self.GetOneValue("CreateTime", &val)
	return val
}

func (self *RedisOrderObj) GetOrderType() int32 {
	var val int32
	self.GetOneValue("OrderType", &val)
	return val
}

func (self *RedisOrderObj) GetExtendValue() string {
	var val string
	self.GetOneValue("ExtendValue", &val)
	return val
}

//对外方法，获取订单管理对象，如果为nil表示redis内存不存在，数据库也不存在
func GetRedisOrderObj(id string) *RedisOrderObj {
	return OrderMgr.GetRedisOrderObj(id)
}

//创建新的订单号
func RedisCreateOrderNo(changeType, sourceType int32) string {
	id := NextId(TABLE_ORDER)
	r := int64(RandInt(10000, 99999))
	orderNo := easygo.AnytoA(changeType) + easygo.AnytoA(sourceType) + easygo.AnytoA(GetMillSecond()+id+1000000000) + easygo.AnytoA(r)
	return orderNo
}

//创建新订单
func CreateRedisOrder(order *share_message.Order) *RedisOrderObj {
	return OrderMgr.CreateRedisOrderObj(order)
}

//批量保存需要存储的数据
func SaveRedisOrderToMongo() {
	ids := []string{}
	GetAllRedisSaveList(TABLE_ORDER, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisOrderObj(id)
		if obj != nil {
			data := obj.GetRedisOrder()
			saveData = append(saveData, bson.M{"_id": data.GetOrderId()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_ORDER, saveData)
	}
}

//============充值发货订单
func CheckOrderDoing(id string) bool {
	b := easygo.RedisMgr.GetC().SIsMember(ORDER_DOING_LOCK, id)
	return b
}
func LockOrderDoing(id string) {
	err := easygo.RedisMgr.GetC().SAdd(ORDER_DOING_LOCK, id)
	easygo.PanicError(err)
}
func UnLockOrderDoing(id string) {
	err := easygo.RedisMgr.GetC().SRem(ORDER_DOING_LOCK, id)
	easygo.PanicError(err)
}
func CleanOrderDoing() {
	_, err := easygo.RedisMgr.GetC().Delete(ORDER_DOING_LOCK)
	easygo.PanicError(err)
}
