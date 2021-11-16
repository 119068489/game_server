package for_game

import (
	"fmt"
	"game_server/easygo"
)

//在某个某个比赛房间的用户
type RedisLiveRoomPlayerObj struct {
	Id int64
	RedisBase
}

var EPORTS_LIVE_ROOMP_LAYER = ESportExN("live_room_player") //redis内存中存在的key

func (self *RedisLiveRoomPlayerObj) Init(LiveId int64) *RedisLiveRoomPlayerObj {
	self.Id = LiveId
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, EPORTS_LIVE_ROOMP_LAYER)
	self.Sid = ESportLiveRoomPlayerMgr.GetSid()
	ESportLiveRoomPlayerMgr.Store(self.Id, self)
	return self
}
func (self *RedisLiveRoomPlayerObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisLiveRoomPlayerObj) GetKeyId() string { //override
	return MakeRedisKey(EPORTS_LIVE_ROOMP_LAYER, self.Id)
}
func (self *RedisLiveRoomPlayerObj) UpdateData() { //override

}

//当检测到redis key不存在时
func (self *RedisLiveRoomPlayerObj) InitRedis() { //override

}
func (self *RedisLiveRoomPlayerObj) GetRedisSaveData() interface{} { //override
	return nil
}
func (self *RedisLiveRoomPlayerObj) SaveOtherData() { //override

}
func (self *RedisLiveRoomPlayerObj) MakePlyerKey(playerId int64) string {
	return fmt.Sprintf("%d", playerId)
}

//进入房间
func (self *RedisLiveRoomPlayerObj) EnterRoom(playerId int64) {
	self.SetOneValue(self.MakePlyerKey(playerId), playerId)
}

//离开房间
func (self *RedisLiveRoomPlayerObj) LeaveRoom(playerId int64) {
	_, err := easygo.RedisMgr.GetC().Hdel(self.GetKeyId(), self.MakePlyerKey(playerId))
	easygo.PanicError(err)
}
func (self *RedisLiveRoomPlayerObj) GetPlayerIds() []int64 {

	values, err := ObjInt64List(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	return values
}

func NewRedisLiveRoomPlayerObj(id int64) *RedisLiveRoomPlayerObj {
	obj := RedisLiveRoomPlayerObj{}
	return obj.Init(id)
}

//对外方法，获取对象，如果为nil表示redis内存不存在，数据库也不存在
func GetRedisLiveRoomPlayerObj(id int64) *RedisLiveRoomPlayerObj {
	obj, ok := ESportLiveRoomPlayerMgr.Load(id)
	if ok && obj != nil {
		return obj.(*RedisLiveRoomPlayerObj)
	} else {
		return NewRedisLiveRoomPlayerObj(id)
	}
}
