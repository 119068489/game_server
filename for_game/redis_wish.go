package for_game

import (
	"encoding/json"
	"errors"
	"fmt"
	"game_server/easygo"
	"game_server/pb/share_message"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/garyburd/redigo/redis"
)

const (
	REDIS_POOL = "pool"       // 水池
	COOL_DOWN  = "cool_down"  // 用户冷却数据 cool_down:1
	DARE_COUNT = "dare_count" // 需要冷却的抽奖次数 dare_count:1
)
const WISH_REDIS_LOCK_KEY = "wish:redis_lock:"

const (
	WARNING_RECYCLE_COUNT   = "wish:warning:recycle_count"   // hset warning_recycle_count pid : count
	WARNING_RECYCLE_AMOUNT  = "wish:warning:recycle_amount"  // hset warning_recycle_amount pid : count
	WARNING_RECYCLE_DIAMOND = "wish:warning:recycle_diamond" // hset warning_recycle_DIAMOND pid : count

	WISH_PLAYER_ACCESS_LOG = "wish:player_access_log"
)

type WishPoolEX struct {
	PoolId       int64 `json:"_id"` //玩家id
	PoolLimit    int64 //水池上限
	InitialValue int64 // 水池初始值(存库)
	IncomeValue  int64 //水池当前收入值
	Recycle      int64 //回收阀值(存库)
	Commission   int64 //官方抽水(回收比例)(存库)
	StartAward   int64 //开启放奖金额(存库)
	CloseAward   int64 //关闭放奖金额(存库)

	Name       string // 名称
	CreateTime int64  //水池创建时间
	UpdateTime int64  //水池更新时间

	ShowInitialValue int64 // 水池初始值(百分比%，后台展示)
	ShowRecycle      int64 // 回收阀值(百分比%，后台展示)
	ShowCommission   int64 // 官方抽水(回收比例)(百分比%，后台展示)
	ShowStartAward   int64 // 开启放奖金额(百分比%，后台展示)
	ShowCloseAward   int64 // 关闭放奖金额(百分比%，后台展示)

	IsOpenAward bool // 是否开启放奖 true 放奖 false 关闭

	IsDefault   bool  // 是否默认水池 true 默认 false 非默认
	LocalStatus int64 // 当前水池的状态;1-大亏;2-小亏;3-普通;4-大赢;5-小盈

	PoolConfigId int64 // 水池模板配置id
}

type WishPoolStatusEX struct {
	MaxValue     int64 // 上限值大于等于(存库) 80
	MinValue     int64 // 下限值大于等于(存库)
	ShowMaxValue int64 // 上限值大于等于(%)(百分比%，后台展示)
	ShowMinValue int64 // 下限值大于等于(%)(百分比%，后台展示)
}

type WishPlayerAccessLogEx struct {
	WishTime      int64 //上次访问许愿池日期（天）
	ExchangeTime  int64 //上次访问兑换钻石页日期（天）
	DareTime      int64 //上次访问挑战赛日期（天）
	ChallengeTime int64 //上次点击挑战赛盲盒日期（天）
	VExchangeTime int64 //上次兑换钻石成功日期（天）
	RetainedDay   int32 //最新留存天数
}

func Set(key string, v interface{}, expire int64) error {
	if expire < 0 {
		return errors.New("不允许设置永久有效的key")
	}
	err := easygo.RedisMgr.GetC().Set(key, v)
	if err != nil {
		return err
	}
	err = Expire(key, expire)
	if err != nil {
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	d, err := easygo.RedisMgr.GetC().Get(key)
	return d, err
}

func Expire(name string, newSecondsLifeTime int64) error {
	err := easygo.RedisMgr.GetC().Expire(name, newSecondsLifeTime)
	return err
}

func Delete(keys string) error {
	_, e := easygo.RedisMgr.GetC().Delete(keys)
	return e
}

func GetDareRedisKey(userId int64) string {
	dateStr := time.Now().Format("2006-01-02")
	key := MakeRedisKey(fmt.Sprintf("%s%d", wish, userId), dateStr)
	return key
}

//设置用户挑战次数
func SetDareFrequency(userId int64) error {
	oneDay := easygo.GetToday24ClockTimestamp() - time.Now().Unix() // 0 点清零
	key := GetDareRedisKey(userId)
	frequency := GetDareFrequency(userId)
	if frequency < 1 {
		frequency = 1
	} else {
		frequency += 1
	}
	err := Set(key, frequency, oneDay)
	return err
}

//设置用户挑战次数
func SetDareFrequency1(count, userId, boxId, boxItemId int64) int64 {
	// 先判断是否存在,如果不存,就要设置过期时间,如果存在,不用设置过期时间
	key := GetDareRedisKey1(userId, boxId, boxItemId)
	b, err := easygo.RedisMgr.GetC().Exist(key)
	easygo.PanicError(err)
	value := easygo.RedisMgr.GetC().IncrBy(key, count)
	if !b { // 不存在的时候,设置过期时间.
		oneDay := easygo.GetToday24ClockTimestamp() - time.Now().Unix() // 0 点清零
		//oneDay := int64(180)                             // 0 点清零
		err = easygo.RedisMgr.GetC().Expire(key, oneDay) // 设置过期时间.
	}
	easygo.PanicError(err)
	return value
}

// 修改愿望要先干掉以前的
func DelDareFrequency1(userId, boxId, boxItemId int64) error {
	key := GetDareRedisKey1(userId, boxId, boxItemId)
	return Delete(key)
}

var wish = "wish:dare:"

func GetDareRedisKey1(userId, boxId, boxItemId int64) string {
	//dateStr := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s%d_%d_%d", wish, userId, boxId, boxItemId)
	return key
}
func GetPlayerMakeWishKey(userId, boxId int64) string {
	//dateStr := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s%d_%d", wish, userId, boxId)
	return key
}

//获取用户挑战次数
func GetDareFrequency(userId int64) int {
	key := GetDareRedisKey(userId)
	frequency, _ := Get(key)
	return easygo.StringToIntnoErr(frequency)
}

//获取用户挑战次数
func GetDareFrequency1(userId, boxId, boxItemId int64) int {
	key := GetDareRedisKey1(userId, boxId, boxItemId)
	frequency, _ := Get(key)
	return easygo.StringToIntnoErr(frequency)
}

// 设置物品的类型
func SetWishItemType(value string) {
	key := MakeRedisKey("wish", "itemType")
	oneDay := 24 * 3600 //1天的时间
	err := Set(key, value, int64(oneDay))
	easygo.PanicError(err)
}

// 获取物品的类型
func GetWishItemType() string {
	key := MakeRedisKey("wish", "itemType")
	value, err := Get(key)
	if err != nil && err.Error() != redis.ErrNil.Error() {
		easygo.PanicError(err)
	}
	return value
}

// 设置物品的品牌
func SetWishBrand(value string) {
	key := MakeRedisKey("wish", "brand")
	oneDay := 24 * 3600 //1天的时间
	err := Set(key, value, int64(oneDay))
	easygo.PanicError(err)
}

// 获取物品的品牌
func GetWishBrand() string {
	key := MakeRedisKey("wish", "brand")
	value, err := Get(key)
	if err != nil && err.Error() != redis.ErrNil.Error() {
		easygo.PanicError(err)
	}
	return value
}

// 水池价格变化
func UpdatePollPriceToRedis(poolId, value int64, key string) (int64, error) {
	// 判断是否存在,如果不存在
	conn := easygo.RedisMgr.GetC()
	redisKey := MakeRedisKey(REDIS_POOL, poolId)
	b, err := conn.Exist(redisKey)
	easygo.PanicError(err)
	if !b {
		// 查询数据库设置进redis
		pool := GetPoolInfoFromDB(poolId)
		if pool == nil {
			return 0, errors.New("水池数据为空")
		}
		wishPool := &WishPoolEX{}
		StructToOtherStruct(pool, wishPool)
		err = conn.HMSet(redisKey, wishPool)
		easygo.PanicError(err)
		SetWishPoolOtherData(pool)
	}

	result := conn.HIncrBy(redisKey, key, value)
	return result, nil
}

// 获取水池的信息
func GetPollInfoFromRedis(poolId int64) *share_message.WishPool {
	// 判断是否存在,如果不存在
	conn := easygo.RedisMgr.GetC()
	redisKey := MakeRedisKey(REDIS_POOL, poolId)
	b, err := conn.Exist(redisKey)
	easygo.PanicError(err)
	var data *share_message.WishPool
	if !b {
		// 查询数据库设置进redis
		pool := GetPoolInfoFromDB(poolId)
		if pool == nil {
			return data
		}
		wishPool := &WishPoolEX{}
		StructToOtherStruct(pool, wishPool)
		err = conn.HMSet(redisKey, wishPool)
		easygo.PanicError(err)
		SetWishPoolOtherData(pool)
	}
	value, err := conn.HGetAll(redisKey)
	if err != nil {
		easygo.PanicError(err)
		return data
	}
	if len(value) == 0 {
		return data
	}
	var base WishPoolEX
	err = redis.ScanStruct(value, &base)
	easygo.PanicError(err)
	StructToOtherStruct(base, &data)
	GetWishPoolOtherData(data)
	return data
}

// 获取水池的信息
func RefreshRedisPollInfo(poolId int64) {
	// 判断是否存在,如果不存在
	conn := easygo.RedisMgr.GetC()
	redisKey := MakeRedisKey(REDIS_POOL, poolId)
	_, err := conn.Exist(redisKey)
	easygo.PanicError(err)
	// 查询数据库设置进redis
	pool := GetPoolInfoFromDB(poolId)
	wishPool := &WishPoolEX{}
	StructToOtherStruct(pool, wishPool)
	err = conn.HMSet(redisKey, wishPool)
	easygo.PanicError(err)
	SetWishPoolOtherData(pool)
}

//单key取得redis分布式锁没有有重试机制   timeout 单位: 秒
func DoRedisLockNoRetry(key string, timeout int32) error {
	err := easygo.RedisMgr.GetC().DoRedisLockNoRetry(key, timeout)
	return err
}

// 单key取得redis分布式锁有重试机制   timeout 单位: 秒
func DoRedisLockWithRetry(key string, timeout int32) error {
	return easygo.RedisMgr.GetC().DoRedisLockWithRetry(key, timeout)
}

// 释放 redis 分布式锁
func DoRedisUnlock(key string) {
	easygo.RedisMgr.GetC().DoRedisUnlock(key)
}

// 设置是否开启放奖
func SetIsOpenAward(poolId int64, b bool) {
	redisKey := MakeRedisKey(REDIS_POOL, poolId)
	_ = easygo.RedisMgr.GetC().HSet(redisKey, "IsOpenAward", b)
}

func SetLocalStatus(poolId, status int64) {
	redisKey := MakeRedisKey(REDIS_POOL, poolId)
	_ = easygo.RedisMgr.GetC().HSet(redisKey, "LocalStatus", status)
}

// 设置水池的其他结构体进redis
func SetWishPoolOtherData(obj *share_message.WishPool) {
	redisKey := MakeRedisKey(REDIS_POOL, obj.GetId())
	SetStringValueToRedis(redisKey, "BigLoss", obj.GetBigLoss())
	SetStringValueToRedis(redisKey, "SmallLoss", obj.GetSmallLoss())
	SetStringValueToRedis(redisKey, "Common", obj.GetCommon())
	SetStringValueToRedis(redisKey, "BigWin", obj.GetBigWin())
	SetStringValueToRedis(redisKey, "SmallWin", obj.GetSmallWin())
}

// 结构体中的就结构体
func SetStringValueToRedis(key, name string, data interface{}) {
	val, err := json.Marshal(data)
	err = easygo.RedisMgr.GetC().HSet(key, name, string(val))
	easygo.PanicError(err)
}

func GetWishPoolOtherData(obj *share_message.WishPool) {
	redisKey := MakeRedisKey(REDIS_POOL, obj.GetId())
	GetStringValueToRedis(redisKey, "BigLoss", &obj.BigLoss)
	GetStringValueToRedis(redisKey, "SmallLoss", &obj.SmallLoss)
	GetStringValueToRedis(redisKey, "Common", &obj.Common)
	GetStringValueToRedis(redisKey, "BigWin", &obj.BigWin)
	GetStringValueToRedis(redisKey, "SmallWin", &obj.SmallWin)
}

// 从redis中获取字符串重新赋值给结构体
func GetStringValueToRedis(key, name string, data interface{}) {
	val, err := easygo.RedisMgr.GetC().HGet(key, name)
	if len(val) == 0 {
		return
	}
	easygo.PanicError(err)
	err = json.Unmarshal(val, &data)
	easygo.PanicError(err)
}

// 用户挑战次数
func AddPlayerDareCountFromRedis(pid, exp int64) int64 {
	//DARE_COUNT
	redisKey := MakeRedisKey(DARE_COUNT, pid)
	conn := easygo.RedisMgr.GetC()
	b, e := conn.Exist(redisKey)
	easygo.PanicError(e)
	count := conn.HIncrBy(redisKey, "count", 1)
	if !b {
		err := conn.Expire(redisKey, exp)
		easygo.PanicError(err)
	}
	return count
}

// 设置用户的冷却时间间隔
func SetPlayerCoolDownTimeFromRedis(pid int64, exp int64) {
	redisKey := MakeRedisKey(COOL_DOWN, pid)
	conn := easygo.RedisMgr.GetC()
	m := make(map[string]interface{})
	m["playerId"] = pid
	m["createTime"] = time.Now().Unix()
	err := conn.HMSet(redisKey, m)
	easygo.PanicError(err)
	err = conn.Expire(redisKey, exp)
	easygo.PanicError(err)
}

// 是否在冷却期
func IsExistPlayerCoolDownTimeFromRedis(pid int64) (int64, bool) {
	redisKey := MakeRedisKey(COOL_DOWN, pid)
	b, err := easygo.RedisMgr.GetC().Exist(redisKey)
	easygo.PanicError(err)
	var createTime int64
	if b {
		by, err := easygo.RedisMgr.GetC().HGet(redisKey, "createTime")
		easygo.PanicError(err)
		createTime = easygo.AtoInt64(string(by))
	}
	return createTime, b
}

// 设置回收次数
func SetRecycleCountToRedis(count, expire int64) int64 {
	// 先判断是否存在,如果不存,就要设置过期时间,如果存在,不用设置过期时间
	b, err := easygo.RedisMgr.GetC().Exist(WARNING_RECYCLE_COUNT)
	easygo.PanicError(err)
	value := easygo.RedisMgr.GetC().IncrBy(WARNING_RECYCLE_COUNT, count)
	if !b { // 不存在的时候,设置过期时间.
		err = easygo.RedisMgr.GetC().Expire(WARNING_RECYCLE_COUNT, expire) // 设置过期时间.
	}
	easygo.PanicError(err)
	return value
}

// 设置回收总金额 单位:分
func SetRecycleAmountToRedis(amount, expire int64) int64 {
	// 先判断是否存在,如果不存,就要设置过期时间,如果存在,不用设置过期时间
	b, err := easygo.RedisMgr.GetC().Exist(WARNING_RECYCLE_AMOUNT)
	easygo.PanicError(err)
	value := easygo.RedisMgr.GetC().IncrBy(WARNING_RECYCLE_AMOUNT, amount)
	if !b { // 不存在的时候,设置过期时间.
		err = easygo.RedisMgr.GetC().Expire(WARNING_RECYCLE_AMOUNT, expire) // 设置过期时间.
	}
	easygo.PanicError(err)
	return value
}

// 设置回收总钻石
func SetRecycleDiamondToRedis(count, expire int64) int64 {
	// 先判断是否存在,如果不存,就要设置过期时间,如果存在,不用设置过期时间
	b, err := easygo.RedisMgr.GetC().Exist(WARNING_RECYCLE_DIAMOND)
	easygo.PanicError(err)
	value := easygo.RedisMgr.GetC().IncrBy(WARNING_RECYCLE_DIAMOND, count)
	if !b { // 不存在的时候,设置过期时间.
		err = easygo.RedisMgr.GetC().Expire(WARNING_RECYCLE_DIAMOND, expire) // 设置过期时间.
	}
	easygo.PanicError(err)
	return value
}

/**
回收预警功能服务器需要根据设定的预警规则，给预警人员发送短信
*/
func CheckWishWarningSMS(amount, diamond int64) {
	// 从配置表中读取配置.
	config := GetConfigWishPayWarn()
	if config == nil {
		return
	}
	// 兼容手机号
	phones := make([]string, 0)
	for _, p := range config.GetPhoneList() {
		if !strings.HasPrefix(p, "+86") {
			p = fmt.Sprintf("%s%s", "+86", p)
		}
		phones = append(phones, p)
	}

	recycleTime := config.GetWithdrawalTime()               // 回收预警频次时间 取出来的单位是秒
	recycleTimes := config.GetWithdrawalTimes()             // 回收预警频次次数
	recycleGoldRate := config.GetWithdrawalGoldRate()       // 回收预警金额时间 取出来的单位是秒
	recycleGold := config.GetWithdrawalGold()               // 回收预警总金额
	recycleDiamondRate := config.GetWithdrawalDiamondRate() // 回收预警钻石时间 取出来的单位是秒
	recycleDiamond := config.GetWithdrawalDiamond()         // 回收预警总钻石

	// 设置回收预警次数
	count := SetRecycleCountToRedis(1, recycleTime)
	if count >= recycleTimes {
		logs.Error("用户 回收次数超出预警,%d 秒限制 %d 次,现在是 %d 次了", recycleTime, recycleTimes, count)
		//  发送腾讯云警告短信
		NewSMSInst(SMS_BUSINESS_TC).SendWarningSMS(phones, []string{easygo.Stamp2Str(time.Now().Unix())})
		return
	}
	if amount > 0 {
		gold := SetRecycleAmountToRedis(amount, recycleGoldRate)
		if gold >= recycleGold {
			logs.Error("用户 回收金额超出预警,%d 秒限制 %v 元,现在是 %v 元了", recycleGoldRate, recycleGold/100, gold/100)
			//  发送腾讯云警告短信
			NewSMSInst(SMS_BUSINESS_TC).SendWarningSMS(phones, []string{easygo.Stamp2Str(time.Now().Unix())})
			return
		}
	}
	if diamond > 0 {
		warnDiamond := SetRecycleDiamondToRedis(diamond, recycleGoldRate)
		if warnDiamond >= recycleDiamond {
			logs.Error("用户 回收钻石超出预警,%d 秒限制 %v 元,现在是 %v 元了", recycleDiamondRate, recycleDiamond, warnDiamond)
			//  发送腾讯云警告短信
			NewSMSInst(SMS_BUSINESS_TC).SendWarningSMS(phones, []string{easygo.Stamp2Str(time.Now().Unix())})
			return
		}
	}
}

func createWishPlayerAccessLogToRedis(pid int64) bool {
	conn := easygo.RedisMgr.GetC()
	redisKey := MakeRedisKey(WISH_PLAYER_ACCESS_LOG, pid)
	b, err := conn.Exist(redisKey)
	easygo.PanicError(err)

	if !b {
		// 查询数据库设置进redis
		data := GetWishPlayerAccessFormDB(pid)
		if data == nil {
			return false
		}
		accessLogs := &WishPlayerAccessLogEx{}
		StructToOtherStruct(data, accessLogs)
		err = conn.HMSet(redisKey, accessLogs)
		easygo.PanicError(err)
		err = conn.Expire(redisKey, 600)
		easygo.PanicError(err)
	}
	return true
}

// 设置许愿池埋点日志
func GetWishPlayerAccessLogFromRedis(pid int64) *share_message.WishPlayerAccessLog {
	// 先判断是否存在,如果不存,就要设置过期时间,如果存在,不用设置过期时间
	conn := easygo.RedisMgr.GetC()
	redisKey := MakeRedisKey(WISH_PLAYER_ACCESS_LOG, pid)
	b, err := conn.Exist(redisKey)
	easygo.PanicError(err)

	if !b {
		// 查询数据库设置进redis
		data := GetWishPlayerAccessFormDB(pid)
		if data == nil {
			return nil
		}
		accessLogs := &WishPlayerAccessLogEx{}
		StructToOtherStruct(data, accessLogs)
		err = conn.HMSet(redisKey, accessLogs)
		easygo.PanicError(err)
		err = conn.Expire(redisKey, 600)
		easygo.PanicError(err)
	}
	value, err := conn.HGetAll(redisKey)
	if err != nil {
		easygo.PanicError(err)
		return nil
	}
	if len(value) == 0 {
		return nil
	}
	var data *share_message.WishPlayerAccessLog
	var base WishPlayerAccessLogEx
	err = redis.ScanStruct(value, &base)
	easygo.PanicError(err)
	StructToOtherStruct(base, &data)
	return data
}

//设置RedisWishLogReport字段值
func SetWishPlayerAccessLogToRedis(pid, val int64, field string) {
	conn := easygo.RedisMgr.GetC()
	redisKey := MakeRedisKey(WISH_PLAYER_ACCESS_LOG, pid)
	b, err := conn.Exist(redisKey)
	easygo.PanicError(err)
	if b {
		_ = easygo.RedisMgr.GetC().HSet(MakeRedisKey(WISH_PLAYER_ACCESS_LOG, pid), field, val)
	}
	UpsertWishPlayerAccessFormDB(pid, val, field)
}

//设置RedisWishLogReport字段值
func IncrWishPlayerAccessLogToRedis(pid, val int64, field string) int64 {
	conn := easygo.RedisMgr.GetC()
	redisKey := MakeRedisKey(WISH_PLAYER_ACCESS_LOG, pid)
	b, err := conn.Exist(redisKey)
	easygo.PanicError(err)
	if !b {
		// 查询数据库设置进redis
		data := GetWishPlayerAccessFormDB(pid)
		if data == nil {
			return 0
		}
		accessLogs := &WishPlayerAccessLogEx{}
		StructToOtherStruct(data, accessLogs)
		err = conn.HMSet(redisKey, accessLogs)
		easygo.PanicError(err)
		err = conn.Expire(redisKey, 600)
		easygo.PanicError(err)
	}
	val = conn.HIncrBy(MakeRedisKey(WISH_PLAYER_ACCESS_LOG, pid), field, val)
	UpsertWishPlayerAccessFormDB(pid, val, field)
	return val
}

// 设置许愿必中记录
func SetPlayerMakeWish(count, userId, boxId int64) int64 {
	key := GetPlayerMakeWishKey(userId, boxId)
	value := easygo.RedisMgr.GetC().IncrBy(key, count)
	return value
}

func GetPlayerMakeWish(userId, boxId int64) int {
	key := GetPlayerMakeWishKey(userId, boxId)
	frequency, _ := Get(key)
	return easygo.StringToIntnoErr(frequency)
}
