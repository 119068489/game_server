package for_game

import (
	"fmt"
	"game_server/easygo"
)

type RedisESportThumbsUpObj struct {
	PlayerId    int64
	MenuId      int32
	DataId      int64
	LastActTime int64 //上次访问时间
	Items       map[int64]string
	RedisBase
}

func newRedisESportThumbsUp(menuId int32, dataId, playerId int64) *RedisESportThumbsUpObj {
	obj := RedisESportThumbsUpObj{}
	obj.Init(menuId, dataId, playerId)
	return &obj
}

var EPORTS_THUMBSUP = ESportExN("thumbsup") //redis内存中存在的key

func (self *RedisESportThumbsUpObj) Init(menuId int32, dataId, playerId int64) *RedisESportThumbsUpObj {
	self.MenuId = menuId
	self.DataId = dataId
	self.PlayerId = playerId
	self.RedisBase.Init(self, self.DataId, easygo.MongoMgr, MONGODB_NINGMENG, EPORTS_THUMBSUP)
	return self
}
func (self *RedisESportThumbsUpObj) GetId() interface{} { //override
	return self.PlayerId
}
func (self *RedisESportThumbsUpObj) GetKeyId() string { //override
	return self.MakeKey()
}
func (self *RedisESportThumbsUpObj) UpdateData() { //override
	t := GetMillSecond()
	if t-self.CreateTime > ESPORT_REDIS_PLAYER_EXIST_TIME {
		kstr := GetRedisESportThumbsUpObjKey(self.MenuId, self.DataId, self.PlayerId)
		ESportThumbsUpMgr.Delete(kstr) // 释放对象
		return
	}

	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
}

//当检测到redis key不存在时
func (self *RedisESportThumbsUpObj) InitRedis() { //override

}
func (self *RedisESportThumbsUpObj) GetRedisSaveData() interface{} { //override
	return nil
}
func (self *RedisESportThumbsUpObj) SaveOtherData() { //override

}
func (self *RedisESportThumbsUpObj) MakeItemKey(itemId int64) string {
	return fmt.Sprintf("%d", itemId)
}
func (self *RedisESportThumbsUpObj) MakeKey() string {
	if self.DataId < 1 {
		return fmt.Sprintf("%s:%d:%d", EPORTS_THUMBSUP, self.PlayerId, self.MenuId) //视频或者资讯点赞
	} else {
		return fmt.Sprintf("%s:%d:%d:nv%d", EPORTS_THUMBSUP, self.PlayerId, self.MenuId, self.DataId) //评论点赞 //self.DataId=视频或资讯的ID
	}
}

//点赞
func (self *RedisESportThumbsUpObj) AddThumbsUp(itemId int64) {
	self.SetOneValue(self.MakeItemKey(itemId), 1)
}

//这文章是否点赞
func (self *RedisESportThumbsUpObj) IsThumbsUp(itemId int64) bool {
	vals, err := easygo.RedisMgr.GetC().HGetAll(self.MakeKey())
	b, _ := ObjListExistStrKey(vals, err, self.MakeItemKey(itemId))
	return b
	//self.Items = self.GetItemsIds()
	//it := self.Items[itemId]
	//return it != "" && it == "1"
}

//文章点赞列表
func (self *RedisESportThumbsUpObj) GetThumbsList() []string {
	vals, err := easygo.RedisMgr.GetC().HGetAll(self.MakeKey())
	b, _ := ObjListToStrKeyList(vals, err)
	return b
}

//这文章是否点赞
func (self *RedisESportThumbsUpObj) IsInThumbsList(list []string, itemId int64) bool {

	itkey := self.MakeItemKey(itemId)
	for _, it := range list {
		if itkey == it {
			return true
		}
	}
	return false
}

//取消
func (self *RedisESportThumbsUpObj) CancelThumbsUp(itemId int64) {
	_, err := easygo.RedisMgr.GetC().Hdel(self.MakeKey(), self.MakeItemKey(itemId))
	easygo.PanicError(err)
}

func GetRedisESportThumbsUpObjKey(menuId int32, dataId, playerId int64) string {

	keystr := fmt.Sprintf("%d_%d_%d", menuId, dataId, playerId)

	return keystr
}

func GetRedisESportThumbsUpObj(menuId int32, dataId, playerId int64) *RedisESportThumbsUpObj {

	kstr := GetRedisESportThumbsUpObjKey(menuId, dataId, playerId)
	obj, ok := ESportThumbsUpMgr.Load(kstr)
	if ok && obj != nil {
		tobj := obj.(*RedisESportThumbsUpObj)
		return tobj
	} else {
		obj := newRedisESportThumbsUp(menuId, dataId, playerId)
		if obj != nil {
			ESportThumbsUpMgr.Store(kstr, obj)
			easygo.AfterFunc(REDIS_SAVE_TIME, obj.UpdateData)
		}
		return obj
	}
}
