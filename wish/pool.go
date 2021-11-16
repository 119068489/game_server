package wish

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/share_message"
	"strings"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//=======================================

// 水池对象
type PoolObj struct {
	PoolId int64
	Mutex  easygo.RLock
	Mutex1 easygo.RLock
	Mutex2 easygo.RLock // 水池安全锁.
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
	PoolMgr.Mutex.Lock()
	defer PoolMgr.Mutex.Unlock()
	pool := PoolMgr.LoadPool(poolId)
	if pool != nil {
		return pool
	}
	obj := NewPool(poolId)
	return obj
}

// 修改水池金额
func (self *PoolObj) UpdatePollPriceToRedis(poolId, value int64, key string) (int64, error) {
	//self.Mutex.Lock()
	//defer self.Mutex.Unlock()
	i, e := for_game.UpdatePollPriceToRedis(poolId, value, key)
	return i, e
}

// 获取水池的信息
func (self *PoolObj) GetPollInfoFromRedis() *share_message.WishPool {
	//self.Mutex.Lock()
	//defer self.Mutex.Unlock()
	p := for_game.GetPollInfoFromRedis(self.PoolId)
	return p
}

// 更新水池的信息
func (self *PoolObj) RefreshRedisPollInfo() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for_game.RefreshRedisPollInfo(self.PoolId)
}

// 是否可以开大奖
func (self *PoolObj) GetReward(pBase *for_game.RedisWishPlayerObj, pid, poolId, price int64, redisKey string, goroutineID uint64) (int64, bool) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	pool := for_game.GetPollInfoFromRedis(poolId)
	incomeValue := pool.GetIncomeValue() // 水池当前的收入
	b := incomeValue-pool.GetInitialValue() > price
	var poolPrice int64
	if b { // 中大奖,更新水池金额

		whiteList := for_game.GetWishWhiteList() // 白名单列表
		whiteIds := make([]int64, 0)
		for _, v := range whiteList {
			whiteIds = append(whiteIds, v.GetId())
		}
		// 不是白名单并且不是运营号
		if !util.Int64InSlice(pid, whiteIds) {
			//if pBase.GetTypes() != for_game.WISH_PLAYER_BASE_TYPES_2 {
			// todo 白名单 判断是否是白名单用户,白名单用户不扣除
			pp, err := for_game.UpdatePollPriceToRedis(poolId, 0-price, redisKey)
			easygo.PanicError(err)
			poolPrice = pp
		} else {
			logs.Info("goroutineID: %v,此用户为白名单用户,大奖不参与水池水量变化,id为: %d", goroutineID, pBase.GetPlayerId())
		}

	}
	return poolPrice, b

}

// 后台工具是否可以开大奖
func (self *PoolObj) GetToolReward(poolId, price int64, redisKey string) (int64, bool) {
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
func (self *PoolObj) SetLocalStatus(poolId, status int64) {
	logs.Info("我是修改水池的状态的")
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for_game.SetLocalStatus(poolId, status)
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

// 获取水池当前的状态
func GetPoolStatus(poolId int64, goroutineID ...uint64) int32 {
	poolObj := GetPoolObj(poolId)
	pool := poolObj.GetPollInfoFromRedis()

	f := easygo.Decimal(easygo.AtoFloat64(easygo.AnytoA(pool.GetIncomeValue()))/easygo.AtoFloat64(easygo.AnytoA(pool.GetPoolLimit())), 2)
	poolStatusLimitFloat64 := f * 100
	//poolStatusLimit := easygo.AtoInt64(easygo.AnytoA(poolStatusLimitFloat64))
	poolStatusLimit := int64(easygo.Decimal(poolStatusLimitFloat64, 0))
	logs.Info("goroutineID: %v,poolStatusLimit---------->%v", goroutineID, poolStatusLimit)
	var poolStatus int32 // 1-大亏;2-小亏;3-普通;4-大盈;5-小盈
	bigLoss := pool.GetBigLoss()
	if bigLoss != nil && (poolStatusLimit <= bigLoss.GetShowMaxValue() && poolStatusLimit >= bigLoss.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_BIGLOSS
	}
	smallLoss := pool.GetSmallLoss()
	if smallLoss != nil && (poolStatusLimit <= smallLoss.GetShowMaxValue() && poolStatusLimit >= smallLoss.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_SMALLLOSS
	}
	common := pool.GetCommon()
	if common != nil && (poolStatusLimit <= common.GetShowMaxValue() && poolStatusLimit >= common.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_COMMON
	}
	bigWin := pool.GetBigWin()
	if bigWin != nil && (poolStatusLimit <= bigWin.GetShowMaxValue() && poolStatusLimit >= bigWin.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_BIGWIN
	}
	smallWin := pool.GetSmallWin()
	if smallWin != nil && (poolStatusLimit <= smallWin.GetShowMaxValue() && poolStatusLimit >= smallWin.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_SMALLWIN
	}
	if bigWin != nil && poolStatusLimit >= bigWin.GetShowMaxValue() { // 大于100%
		poolStatus = for_game.POOL_STATUS_BIGWIN
	}
	if bigLoss != nil && (poolStatusLimit <= bigLoss.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_BIGLOSS
	}
	return poolStatus
}
