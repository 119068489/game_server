package for_game

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

//=======================================

const WISH_PLAYER_EXIST_TIME = 600 * 1000
const WISH_PLAYER_ADDRESS = "Address"

type WishPlayerEx struct {
	Id         int64
	Account    string
	Channel    int32  //渠道:1001-柠檬im，1002-语音渠道，1003-其他渠道
	NickName   string //昵称
	HeadUrl    string //头像
	PlayerId   int64  //渠道方玩家唯一id,用来获取详细数据的
	Token      string //登录渠道token，用来校验是否有效登录用户
	IsTryOne   bool   //是否试玩过了
	CreateTime int64  //创建时间
	//repeated WishAddress Address = 10;   //地址
	HallSid                 int32  //所在的大厅
	NotOne                  bool   // false-首次玩,true-不是首次玩
	NotOneWish              bool   // false-首次许愿,true-不是首次许愿
	Diamond                 int64  //钻石
	LastExchangeDiamondTime int64  // 最近一次兑换钻石的时间
	Types                   int64  // 用户类型 0-正常用户，1假用户 2-白名单用户
	IsFreeze                bool   //钻石帐户是否冻结
	FreezeTime              int64  //冻结到期时间 秒
	Note                    string //冻结备注
}

type RedisWishPlayerObj struct {
	Id PLAYER_ID
	RedisBase
}

func NewWishPlayerObj(playerId PLAYER_ID, data ...*share_message.WishPlayer) *RedisWishPlayerObj {
	p := &RedisWishPlayerObj{Id: playerId}
	obj := append(data, nil)[0]
	return p.Init(obj)
}
func (self *RedisWishPlayerObj) Init(obj *share_message.WishPlayer) *RedisWishPlayerObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	self.Sid = WishPlayerMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		WishPlayerMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryWishPlayer(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisWishPlayer(obj)
	}
	//	logs.Info("初始化新的PlayerBase管理器:", self.Id)
	return self
}
func (self *RedisWishPlayerObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisWishPlayerObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_WISH_PLAYER, self.Id)
}

//定时更新数据
func (self *RedisWishPlayerObj) UpdateData() { //override
	if !self.IsExistKey() {
		WishPlayerMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > WISH_PLAYER_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		WishPlayerMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisWishPlayerObj) InitRedis() { //override
	obj := self.QueryWishPlayer(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisWishPlayer(obj)
}
func (self *RedisWishPlayerObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisWishPlayer()
	return data
}
func (self *RedisWishPlayerObj) SaveOtherData() { //override
}

//通过playerId从mongo中读取登录玩家数据
func (self *RedisWishPlayerObj) QueryWishPlayer(id PLAYER_ID) *share_message.WishPlayer {
	data := self.QueryMongoData(id)
	if data != nil {
		var player share_message.WishPlayer
		StructToOtherStruct(data, &player)
		return &player
	}
	//redis中获取数据初始化
	if self.IsExistKey() {
		return self.GetRedisWishPlayer()
	}
	return nil
}
func (self *RedisWishPlayerObj) SetRedisWishPlayer(obj *share_message.WishPlayer) {
	WishPlayerMgr.Store(obj.GetPlayerId(), self)
	self.AddToExistList(obj.GetPlayerId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	//玩家基础信息
	player := &WishPlayerEx{}
	StructToOtherStruct(obj, player)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), player)
	easygo.PanicError(err)
	//其他玩家数据设置到redis
	self.SetPlayerOtherData(obj, false)
}

//玩家其他数据写入redis
func (self *RedisWishPlayerObj) SetPlayerOtherData(obj *share_message.WishPlayer, save ...bool) {
	//	logs.Info("设置玩家其他数据")
	self.SetStringValueToRedis(WISH_PLAYER_ADDRESS, obj.GetAddress(), save...)
}

//获取玩家其他redis数据
func (self *RedisWishPlayerObj) GetPlayerOtherData(obj *share_message.WishPlayer) {
	self.GetStringValueToRedis(WISH_PLAYER_ADDRESS, &obj.Address)
}

//获取玩家信息
func (self *RedisWishPlayerObj) GetRedisWishPlayer() *share_message.WishPlayer {
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	if len(value) == 0 {
		return nil
	}
	easygo.PanicError(err)
	var base WishPlayerEx
	err = redis.ScanStruct(value, &base)
	easygo.PanicError(err)
	var newBase *share_message.WishPlayer
	StructToOtherStruct(base, &newBase)
	//其他数据
	self.GetPlayerOtherData(newBase)
	return newBase
}

//获取用户PlayerId
func (self *RedisWishPlayerObj) GetPlayerId() int64 {
	var val int64
	self.GetOneValue("PlayerId", &val)
	return val
}

//获取金币
func (self *RedisWishPlayerObj) GetDiamond() int64 {
	var val int64
	self.GetOneValue("Diamond", &val)
	return val
}
func (self *RedisWishPlayerObj) IncrDiamond(val int64) int64 {
	return self.IncrOneValue("Diamond", val)
}

func (self *RedisWishPlayerObj) GetNickName() string {
	var val string
	self.GetOneValue("NickName", &val)
	return val
}

func (self *RedisWishPlayerObj) GetTypes() int32 {
	var val int32
	self.GetOneValue("Types", &val)
	return val
}

//获取用户地址
func (self *RedisWishPlayerObj) GetAddressList() []*share_message.WishAddress {
	var val []*share_message.WishAddress
	self.GetStringValueToRedis(WISH_PLAYER_ADDRESS, &val)
	return val
}

//获取用户地址
func (self *RedisWishPlayerObj) GetOneAddress(addressId int64) *share_message.WishAddress {
	val := self.GetAddressList()
	for _, v := range val {
		if v.GetAddressId() == addressId {
			return v
		}
	}
	return nil
}

func (self *RedisWishPlayerObj) SetAddressList(val []*share_message.WishAddress) {
	self.SetStringValueToRedis(WISH_PLAYER_ADDRESS, val)
}

//删除地址
func (self *RedisWishPlayerObj) DelAddress(addressId int64) {
	address := self.GetAddressList()
	for i, b := range address {
		if b.GetAddressId() == addressId {
			address = append(address[:i], address[i+1:]...) //easygo.Del(bankInfos, b).([]*share_message.BankInfo)
			break
		}
	}
	self.SetStringValueToRedis(WISH_PLAYER_ADDRESS, address)
}

// 获取最新一次兑换钻石的时间
func (self *RedisWishPlayerObj) GetLastExchangeDiamondTime() int64 {
	var val int64
	self.GetOneValue("LastExchangeDiamondTime", &val)
	return val
}
func (self *RedisWishPlayerObj) SetLastExchangeDiamondTime(val int64) {
	self.SetOneValue("LastExchangeDiamondTime", val)
}
func (self *RedisWishPlayerObj) GetHeadIcon() string {
	var val string
	self.GetOneValue("HeadUrl", &val)
	return val
}
func (self *RedisWishPlayerObj) GetAccount() string {
	var val string
	self.GetOneValue("Account", &val)
	return val
}

func (self *RedisWishPlayerObj) GetIsFreeze() bool {
	var val bool
	self.GetOneValue("IsFreeze", &val)
	return val
}
func (self *RedisWishPlayerObj) GetFreezeTime() int64 {
	var val int64
	self.GetOneValue("FreezeTime", &val)
	return val
}

//修改许愿池用户Account
func (self *RedisWishPlayerObj) SetAccount(account string) {
	self.SetOneValue("Account", account)
	self.SaveOneRedisDataToMongo("Account", account)
}

//设置已试玩
func (self *RedisWishPlayerObj) SetTryOne() {
	self.SetOneValue("IsTryOne", true)
	self.SaveOneRedisDataToMongo("IsTryOne", true)
}

//设置token
func (self *RedisWishPlayerObj) SetToken(token string) {
	self.SetOneValue("Token", token)
	self.SaveOneRedisDataToMongo("Token", token)
}

//设置NoOne
func (self *RedisWishPlayerObj) SetNoOne(b bool) {
	self.SetOneValue("NotOne", b)
	self.SaveOneRedisDataToMongo("NotOne", b)
}

// false-首次许愿,true-不是首次许愿
func (self *RedisWishPlayerObj) SetNotOneWish(b bool) {
	self.SetOneValue("NotOneWish", b)
	self.SaveOneRedisDataToMongo("NotOneWish", b)
}

//冻结钻石帐户
func (self *RedisWishPlayerObj) SetIsFreeze(b bool) {
	self.SetOneValue("IsFreeze", b)
	self.SaveOneRedisDataToMongo("IsFreeze", b)
}

//
func (self *RedisWishPlayerObj) SetFreezeTime(val int64) {
	self.SetOneValue("FreezeTime", val)
	self.SaveOneRedisDataToMongo("FreezeTime", val)
}
func (self *RedisWishPlayerObj) SetNote(val string) {
	self.SetOneValue("Note", val)
	self.SaveOneRedisDataToMongo("Note", val)
}
func (self *RedisWishPlayerObj) SetOperator(val string) {
	self.SetOneValue("Operator", val)
	self.SaveOneRedisDataToMongo("Operator", val)
}

//返回值:fail==nil,成功，否则失败
func (self *RedisWishPlayerObj) AddDiamond(value int64, reason string, sourceType int32, extendLog interface{}) (*base.Fail, int64) {
	newDiamond := self.IncrDiamond(value)
	if newDiamond < 0 {
		//减少钱加回来
		newDiamond = self.IncrDiamond(-value)
		return easygo.NewFailMsg("钻石不足"), newDiamond
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": self.GetId()}, bson.M{"$inc": bson.M{"Diamond": value}})
	easygo.PanicError(err)
	self.AddDiamondLog(value, reason, sourceType, newDiamond-value, newDiamond, extendLog)
	return nil, newDiamond
}
func (self *RedisWishPlayerObj) AddDiamondLog(value int64, reason string, sourceType int32, oldDiamond int64, newDiamond int64, extendLog interface{}) {
	if value == 0 {
		return
	}
	st := GettSourceTypeById(sourceType)
	logId := NextId(TABLE_DIAMOND_CHANGELOG)
	log := share_message.DiamondChangeLog{
		LogId:         &logId,
		PlayerId:      easygo.NewInt64(self.Id),
		ChangeDiamond: &value,
		SourceType:    &sourceType,
		PayType:       st.Type,
		Note:          &reason,

		CurDiamond: &oldDiamond,
		Diamond:    &newDiamond,
		CreateTime: easygo.NewInt64(GetMillSecond()),
	}
	m := &CommonDiamond{
		DiamondChangeLog: log,
		Extend:           extendLog,
	}
	AddDiamondChangeLog(m)
}

//处理回收物品 1-已兑换，2-已回收
func (self *RedisWishPlayerObj) DealPlayerWishItem(dealType int, reqMsg *h5_wish.WishBoxReq) ([]*share_message.PlayerWishItem, int, []int64) {
	whiteList := GetWishWhiteList() // 白名单列表
	whiteIds := make([]int64, 0)
	for _, v := range whiteList {
		whiteIds = append(whiteIds, v.GetId())
	}
	//  是白名单或者是运营号
	if self.GetTypes() >= 2 || util.Int64InSlice(easygo.AtoInt64(easygo.AnytoA(self.GetId())), whiteIds) {
		return nil, DEAL_PLAYER_LIMIT, nil
	}
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	playerId := self.GetId().(int64)
	ids := reqMsg.GetIdList() //PlayerWishItem的id
	//获取用户兑换的物品中间表ID
	playerWishItems, err1 := GetPlayerWishItemByIds(playerId, ids)
	if err1 != nil {
		logs.Error("兑换/回收时查询玩家兑换物品数据失败: %v", err1)
		return nil, DEAL_FAULT, nil
	}
	editStatusIds := make([]int64, 0)
	if len(playerWishItems) < 1 {
		logs.Error("兑换/回收物品为非待兑换状态:playerId：%v, playerWishItemIds: %v", playerId, ids)
		return nil, DEAL_FAULT, nil
	}
	for _, v := range playerWishItems {
		editStatusIds = append(editStatusIds, v.GetId())
	}
	//更改玩家兑换的物品状态
	ExchangeToPlayerItem(editStatusIds, dealType)

	return playerWishItems, DEAL_SUCCESS, editStatusIds
}

// 获取token
func (self *RedisWishPlayerObj) GetToke() string {
	var val string
	self.GetOneValue("Token", &val)
	return val
}

//对外接口
// 封装外部方法，获取玩家的信息
func GetRedisWishPlayer(playerId PLAYER_ID, player ...*share_message.WishPlayer) *RedisWishPlayerObj {
	return WishPlayerMgr.GetRedisWishPlayerObj(playerId, player...)
}
func GetRedisWishPlayerByPid(playerId PLAYER_ID, player ...*share_message.WishPlayer) *RedisWishPlayerObj {
	p := GetWishPlayerInfo(playerId)
	return WishPlayerMgr.GetRedisWishPlayerObj(p.GetId(), player...)
}

func GetRedisWishPlayerByImPid(playerId PLAYER_ID, player ...*share_message.WishPlayer) *RedisWishPlayerObj {
	p := GetWishPlayerInfoByImId(playerId)
	return WishPlayerMgr.GetRedisWishPlayerObj(p.GetId(), player...)
}
