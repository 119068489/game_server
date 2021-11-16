package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

const (
	PLAYER_ACCOUNT_EXIST_LIST = "player_account_exist_list" //玩家账号修改存储列表
	ACCOUNT_EXIST_TIME        = 1000 * 600                  //redis的key删除时间:毫秒
)

type RedisAccountObj struct {
	PlayerId PLAYER_ID
	RedisBase
}

type PlayerAccountEx struct {
	PlayerId    int64 `json:"_id"`
	Account     string
	Password    string
	Email       string
	Token       string
	PayPassword string
	OpenId      string
	CreateTime  int64
	IsBind      bool
	UnionId     string
}

func NewRedisAccount(playerId PLAYER_ID, data ...*share_message.PlayerAccount) *RedisAccountObj {
	p := &RedisAccountObj{
		PlayerId: playerId,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}
func (self *RedisAccountObj) Init(obj *share_message.PlayerAccount) *RedisAccountObj {
	self.RedisBase.Init(self, self.PlayerId, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_ACCOUNT)
	self.Sid = AccountMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		AccountMgr.Store(self.PlayerId, self)
		self.AddToExistList(self.PlayerId)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryPlayerAccount(self.PlayerId)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisAccount(obj)
	}
	logs.Info("初始化新的account管理器", self.PlayerId)
	return self
}
func (self *RedisAccountObj) GetId() interface{} { //override
	return self.PlayerId
}
func (self *RedisAccountObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_PLAYER_ACCOUNT, self.PlayerId)
}

//定时更新数据
func (self *RedisAccountObj) UpdateData() { //override
	if !self.IsExistKey() {
		AccountMgr.Delete(self.PlayerId) // 释放对象
		self.DelToExistList(self.PlayerId)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > ACCOUNT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.PlayerId)
			self.DelRedisKey() //redis删除
		}
		AccountMgr.Delete(self.PlayerId) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}

func (self *RedisAccountObj) InitRedis() { //override
	obj := self.QueryPlayerAccount(self.PlayerId)
	if obj == nil {
		return
	}
	self.SetRedisAccount(obj)
	//重新激活定时器
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisAccountObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisAccount()
	return data
}
func (self *RedisAccountObj) SaveOtherData() { //override
}

//通过playerId从mongo中读取登录玩家数据
func (self *RedisAccountObj) QueryPlayerAccount(id PLAYER_ID) *share_message.PlayerAccount {
	data := self.QueryMongoData(id)
	if data != nil {
		var account share_message.PlayerAccount
		StructToOtherStruct(data, &account)
		return &account
	}
	return nil
}

//设置玩家账号信息
func (self *RedisAccountObj) SetRedisAccount(obj *share_message.PlayerAccount) {
	AccountMgr.Store(obj.GetPlayerId(), self)
	self.AddToExistList(obj.GetPlayerId())
	///重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	account := &PlayerAccountEx{}
	StructToOtherStruct(obj, account)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), account)
	easygo.PanicError(err)

}
func (self *RedisAccountObj) GetRedisAccount() *share_message.PlayerAccount {
	obj := &PlayerAccountEx{}
	key := self.GetKeyId()
	value, err := easygo.RedisMgr.GetC().HGetAll(key)
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.PlayerAccount{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//更新支付密码
func (self *RedisAccountObj) SetPayPassword(pass string) {
	pwd := Md5(pass)
	self.SetOneValue("PayPassword", pwd)
}
func (self *RedisAccountObj) GetPayPassword() string {
	var val string
	self.GetOneValue("PayPassword", &val)
	return val
}
func (self *RedisAccountObj) GetPlayerId() int64 {
	var val int64
	self.GetOneValue("PlayerId", &val)
	return val
}
func (self *RedisAccountObj) GetOpenId() string {
	var val string
	self.GetOneValue("OpenId", &val)
	return val
}
func (self *RedisAccountObj) GetUnionId() string {
	var val string
	self.GetOneValue("UnionId", &val)
	return val
}
func (self *RedisAccountObj) GetIsBind() bool {
	var val bool
	self.GetOneValue("IsBind", &val)
	return val
}
func (self *RedisAccountObj) SetIsBind(isBind bool) {
	self.SetOneValue("IsBind", isBind)
}
func (self *RedisAccountObj) GetCreateTime() int64 {
	var val int64
	self.GetOneValue("CreateTime", &val)
	return val
}

//更新登录密码
func (self *RedisAccountObj) SetPassword(pass string) {
	pwd := Md5(pass)
	self.SetOneValue("Password", pwd)
}
func (self *RedisAccountObj) GetPassword() string {
	var val string
	self.GetOneValue("Password", &val)
	return val
}

//更新token码
func (self *RedisAccountObj) SetToken(token string) {
	self.SetOneValue("Token", token)
}

func (self *RedisAccountObj) GetToken() string {
	var val string
	self.GetOneValue("Token", &val)
	return val
}

//更新openId
func (self *RedisAccountObj) SetOpenId(openId string) {
	self.SetOneValue("OpenId", openId)
}

//更新 unionId
func (self *RedisAccountObj) SetUnionId(unionId string) {
	self.SetOneValue("UnionId", unionId)
}

func (self *RedisAccountObj) GetAccount() string {
	var val string
	self.GetOneValue("Account", &val)
	return val
}

func (self *RedisAccountObj) SetAccount(phone string) {
	self.SetOneValue("Account", phone)
}

//对外方法，获取玩家对象，如果为nil表示redis内存不存在，数据库也不存在
func GetRedisAccountObj(id PLAYER_ID) *RedisAccountObj {
	return AccountMgr.GetRedisAccountObj(id)
}

//创建新的账号
func CreateRedisAccount(account *share_message.PlayerAccount) *RedisAccountObj {
	return AccountMgr.CreateRedisAccount(account)
}

//通过account从mongo中读取登录玩家数据
func GetRedisAccountByPhone(phone string) *RedisAccountObj {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_ACCOUNT)
	defer closeFun()

	player := &share_message.PlayerAccount{}
	err := col.Find(bson.M{"Account": phone}).One(&player)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return CreateRedisAccount(player)
}

//停服保存处理，保存需要存储的数据
func SaveRedisAccountToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_PLAYER_ACCOUNT, &ids)
	//修改
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisAccountObj(id)
		if obj != nil {
			account := obj.GetRedisAccount()
			saveData = append(saveData, bson.M{"_id": account.GetPlayerId()}, account)
			obj.SetSaveStatus(false)
		}
	}
	//更新订单数据
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_ACCOUNT, saveData)
	}
}
