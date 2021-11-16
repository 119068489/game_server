package for_game

import (
	"game_server/easygo"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

var _IdentityMutex easygo.Mutex

// 递增 id 表
type Identity struct {
	Key   *string `bson:"_id"`
	Value *uint64 `bson:"Value,omitempty"`
}

// 生成递增的 id ,每一次都会读写 Database
//只做一个表，统一写再easygo.MongoMgr, for_game.MONGODB_NINGMENG
func NextId(key string, val ...int64) int64 {
	v := append(val, 0)[0]
	_IdentityMutex.Lock()
	defer _IdentityMutex.Unlock()
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ID_GENERATOR)
	defer closeFun()
	identity := &Identity{
		Key:   &key,
		Value: easygo.NewUint64(v),
	}
	_, err := col.Find(bson.M{"_id": key}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"Value": int64(1)}},
		Upsert:    true,
		ReturnNew: true,
	}, &identity)
	easygo.PanicError(err)

	return int64(*identity.Value)
}

// 获取当前id值 ,每一次都会读写 Database
func CurrentId(key string) uint64 {
	_IdentityMutex.Lock()
	defer _IdentityMutex.Unlock()
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ID_GENERATOR)
	defer closeFun()
	identity := &Identity{}
	err := col.Find(bson.M{"_id": key}).One(&identity)
	easygo.PanicError(err)

	return *identity.Value
}

//设置新的id
func SetCurrentId(key string, v int64) {
	_IdentityMutex.Lock()
	defer _IdentityMutex.Unlock()
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ID_GENERATOR)
	defer closeFun()
	identity := &Identity{
		Key:   &key,
		Value: easygo.NewUint64(v),
	}
	_, err := col.Upsert(bson.M{"_id": key}, bson.M{"$set": identity})
	easygo.PanicError(err)
}
