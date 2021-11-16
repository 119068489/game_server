package for_game

import (
	"fmt"
	"game_server/easygo"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

type RedisESportLiveFollowObj struct {
	PlayerId int64
	RedisBase
}

func NewRedisESportLiveFollowObj(playerId int64) *RedisESportLiveFollowObj {
	obj := RedisESportLiveFollowObj{}
	obj.Init(playerId)

	return &obj
}

var EPORTS_FOLLOW = ESportExN("follow") //redis内存中存在的key
func (self *RedisESportLiveFollowObj) Init(playerId int64) *RedisESportLiveFollowObj {
	self.PlayerId = playerId
	self.RedisBase.Init(self, playerId, easygo.MongoMgr, MONGODB_NINGMENG, EPORTS_FOLLOW)
	if !self.IsExistKey() {
		lst := GetPlayerAllLiveFollowList(playerId)
		if len(lst) > 0 {
			mps := make(map[string]int)
			for _, v := range lst {
				did := v.GetDataId()
				mps[self.MakeItemKey(did)] = 1
			}
			self.Mutex.Lock()
			defer self.Mutex.Unlock()
			ric := easygo.RedisMgr.GetC()
			err := ric.HMSet(self.Me.GetKeyId(), mps)
			if err != nil {
				logs.Error(err)
			}
		}
	}

	return self
}
func (self *RedisESportLiveFollowObj) GetId() interface{} { //override
	return self.PlayerId
}
func (self *RedisESportLiveFollowObj) GetKeyId() string { //override
	return self.MakeKey()
}
func (self *RedisESportLiveFollowObj) UpdateData() { //override

}

//当检测到redis key不存在时
func (self *RedisESportLiveFollowObj) InitRedis() { //override

}
func (self *RedisESportLiveFollowObj) GetRedisSaveData() interface{} { //override
	return nil
}
func (self *RedisESportLiveFollowObj) SaveOtherData() { //override

}
func (self *RedisESportLiveFollowObj) MakeItemKey(itemId int64) string {
	return fmt.Sprintf("%d", itemId)
}
func (self *RedisESportLiveFollowObj) MakeKey() string {
	return fmt.Sprintf("%s:%d", EPORTS_FOLLOW, self.PlayerId) //视频或者资讯关注
}

//关注
func (self *RedisESportLiveFollowObj) AddFollow(itemId int64) {
	self.SetOneValue(self.MakeItemKey(itemId), 1)

}

//这是否关注
func (self *RedisESportLiveFollowObj) IsFollow(itemId int64) bool {
	vals, err := easygo.RedisMgr.GetC().HGetAll(self.MakeKey())
	b, _ := ObjListExistStrKey(vals, err, self.MakeItemKey(itemId))
	return b
}

//取消
func (self *RedisESportLiveFollowObj) CancelFollow(itemId int64) {
	_, err := easygo.RedisMgr.GetC().Hdel(self.MakeKey(), self.MakeItemKey(itemId))
	easygo.PanicError(err)
}
func (self *RedisESportLiveFollowObj) GetItemsIds() map[int64]string {
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.MakeKey()))
	easygo.PanicError(err)
	return values
}

func GetPlayerAllLiveFollowList(plyId int64) []*share_message.TableESPortsFlowInfo {
	var list []*share_message.TableESPortsFlowInfo
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_FLOW_LIVE_FOLLOW_HISTORY)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["PlayerId"] = plyId
	query := col.Find(queryBson)
	err := query.Select(bson.M{"DataId": 1}).All(&list)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return list
}

//对外方法，获取玩家对象，如果为nil表示redis内存不存在，数据库也不存在
func GetRedisESportLiveFollowObj(id PLAYER_ID) *RedisESportLiveFollowObj {
	obj, ok := ESportFollowMgr.Load(id)
	if ok && obj != nil {
		return obj.(*RedisESportLiveFollowObj)
	} else {
		return NewRedisESportLiveFollowObj(id)
	}
}
