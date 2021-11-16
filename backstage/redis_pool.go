package backstage

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"strings"
	"sync"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//=======================================

var PoolMgr *PoolManager

// 水池管理器.
type PoolManager struct {
	Mutex easygo.RLock
	sync.Map
}

func InitPoolManager() {
	PoolMgr = &PoolManager{}
}

func Delete(keys string) error {
	_, e := easygo.RedisMgr.GetC().Delete(keys)
	return e
}

func (self *PoolManager) LoadPool(poolId int64) *PoolObj {
	value, ok := self.Load(poolId)
	if ok {
		return value.(*PoolObj)
	}
	return nil
}

// 水池对象
type PoolObj struct {
	PoolId int64
	Mutex  easygo.RLock
	Mutex1 easygo.RLock
}

func NewPool(poolId int64) *PoolObj {
	p := &PoolObj{
		PoolId: poolId,
	}
	PoolMgr.Store(poolId, p)
	return p
}

//=====================对外接口====================

// 对外接口，获取盲盒数据.
func GetPoolObj(poolId int64) *PoolObj {
	//PoolMgr.Mutex.Lock()
	//defer PoolMgr.Mutex.Unlock()
	pool := PoolMgr.LoadPool(poolId)
	if pool != nil {
		return pool
	}
	obj := NewPool(poolId)
	return obj
}

// 修改水池金额
func (self *PoolObj) UpdatePollPriceToRedis(poolId, value int64, key string) (int64, error) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return for_game.UpdatePollPriceToRedis(poolId, value, key)
}

// 获取水池的信息
func (self *PoolObj) GetPollInfoFromRedis() *share_message.WishPool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return for_game.GetPollInfoFromRedis(self.PoolId)
}

// 更新水池的信息
func (self *PoolObj) RefreshRedisPollInfo() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for_game.RefreshRedisPollInfo(self.PoolId)
}

// 删除水池的信息
func (self *PoolObj) DeleteRedisPollInfo() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	// 判断是否存在,如果不存在
	conn := easygo.RedisMgr.GetC()
	redisKey := for_game.MakeRedisKey(for_game.REDIS_POOL, self.PoolId)
	_, err := conn.Exist(redisKey)
	easygo.PanicError(err)
	Delete(redisKey)
}

// 是否可以开大奖
func (self *PoolObj) GetReward(poolId, price int64, redisKey string) (int64, bool) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	pool := for_game.GetPollInfoFromRedis(poolId)
	incomeValue := pool.GetIncomeValue() // 水池当前的收入
	b := incomeValue-pool.GetInitialValue() > price
	var poolPrice int64
	if b { // 中大奖,更新水池金额
		pp, err := for_game.UpdatePollPriceToRedis(poolId, 0-price, redisKey)
		easygo.PanicError(err)
		poolPrice = pp
	}
	return poolPrice, b

}

func (self *PoolObj) SetIsOpenAward(poolId int64, b bool) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for_game.SetIsOpenAward(poolId, b)
}

// 获取所有的水池的key
func GetAllPoolKeyList() []int64 {
	ids := make([]int64, 0)
	keys, err1 := easygo.RedisMgr.GetC().Scan(for_game.REDIS_POOL)

	easygo.PanicError(err1)
	for _, key := range keys {
		id := strings.SplitN(key, ":", 2)[1] // pool:1
		ids = append(ids, easygo.AtoInt64(id))
	}
	logs.Info("---->", ids)
	return ids
}

//停服保存处理，保存需要存储的数据
func SavePoolDataToMongo() {
	ids := GetAllPoolKeyList()
	if len(ids) == 0 {
		return
	}
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetPoolObj(id)
		data := obj.GetPollInfoFromRedis()
		saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
	}
	if len(saveData) > 0 {
		for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL, saveData)
	}
}
