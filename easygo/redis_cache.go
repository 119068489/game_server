package easygo

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/garyburd/redigo/redis"
)

var (
	REDIS_DEFAULT     = time.Duration(0)  // 过期时间 不设置
	REDIS_FOREVER     = time.Duration(-1) // 过期时间不设置
	REDIS_EXPIRE_TIME = int64(1800)       //过期时间
)

type RedisCache struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

// 返回cache 对象, 在多个工具之间建立一个 中间初始化的时候使用
func NewRedisCache(db int, host string, defaultExpiration time.Duration, pass ...string) RedisCache {
	ps := append(pass, "")[0]
	pool := &redis.Pool{
		MaxActive:   50000,                            //  最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）
		MaxIdle:     200,                              // 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
		IdleTimeout: time.Duration(100) * time.Second, // 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用
		Wait:        true,                             // 当超过最大连接数 是报错还是等待， true 等待 false 报错
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host, redis.DialDatabase(db), redis.DialPassword(ps))
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return conn, nil
		},
	}
	return RedisCache{pool: pool, defaultExpiration: defaultExpiration}
}

func Serialization(value interface{}) ([]byte, error) {
	if bytes, ok := value.([]byte); ok {
		return bytes, nil
	}
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
		//case reflect.String:
		//	return []byte(value.(string)), nil
	}
	k, err := json.Marshal(value)
	return k, err
}

func Deserialization(byt []byte, ptr interface{}) (err error) {
	if bytes, ok := ptr.(*[]byte); ok {
		*bytes = byt
		return
	}
	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr {
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var i int64
			i, err = strconv.ParseInt(string(byt), 10, 64)
			if err != nil {
				fmt.Printf("Deserialization: failed to parse int '%s': %s", string(byt), err)
			} else {
				p.SetInt(i)
			}
			return

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var i uint64
			i, err = strconv.ParseUint(string(byt), 10, 64)
			if err != nil {
				fmt.Printf("Deserialization: failed to parse uint '%s': %s", string(byt), err)
			} else {
				p.SetUint(i)
			}
		}
	}
	if len(byt) == 0 {
		return
	}
	err = json.Unmarshal(byt, &ptr)
	return
}

// string 类型 添加, v 可以是任意类型
func (c RedisCache) StringSet(name string, v interface{}) error {
	conn := c.pool.Get()
	s, _ := Serialization(v) // 序列化
	defer conn.Close()
	_, err := conn.Do("SET", name, s)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return err
}

// 获取 字符串类型的值
func (c RedisCache) StringGet(name string, v interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	temp, _ := redis.Bytes(conn.Do("Get", name))
	err := Deserialization(temp, &v) // 反序列化
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return err
}

func (c RedisCache) stringGetTest() {
	//var need []string
	//var need int32
	var need int64
	//Deserialization(aa, &need)
	c.StringGet("yang", &need)
}

// 判断所在的 key 是否存在
func (c RedisCache) Exist(name string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("EXISTS", name))
	return v, err
}

//域是否存在
func (c RedisCache) HExistsEx(name string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("HEXISTS", name))
	return v, err
}

// 自增
func (c RedisCache) StringIncr(name string) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Int(conn.Do("INCR", name))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return v, err
}

func (c RedisCache) Get(name string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	temp, err := redis.Bytes(conn.Do("Get", name))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return string(temp), err
}

func (c RedisCache) Set(name string, v interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", name, v)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return err
}

// 设置过期时间 （单位 秒）
func (c RedisCache) Expire(name string, newSecondsLifeTime int64) error {
	b, er := c.Exist(name)
	if !b {
		return er
	}

	// 设置key 的过期时间
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", name, newSecondsLifeTime)
	return err
}

// 删除指定的键
func (c RedisCache) Delete(keys ...interface{}) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("DEL", keys...))
	return v, err
}

// 查看指定的长度
func (c RedisCache) StrLen(name string) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Int(conn.Do("STRLEN", name))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return v, err
}

// //////////////////  hash ///////////
// 删除指定的 hash 键
func (c RedisCache) Hdel(name string, keys ...string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	var err error
	args := []interface{}{name}
	for _, field := range keys {
		args = append(args, field)
	}
	v, err := redis.Bool(conn.Do("HDEL", args...))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return v, err
}

// 查看hash 中指定是否存在
func (c RedisCache) HExists(name, field string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	var err error
	v, err := redis.Bool(conn.Do("HEXISTS", name, field))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return v, err
}

// 获取hash 的键的个数
func (c RedisCache) HLen(name string) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Int(conn.Do("HLEN", name))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return v, err
}

func (c RedisCache) HKeys(name string) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	value, err := redis.Values(conn.Do("HKEYS", name))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return value, err
}

// 传入的 字段列表获得对应的值
func (c RedisCache) HMGet(name string, fields ...string) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	args := []interface{}{name}
	for _, field := range fields {
		args = append(args, field)
	}
	value, err := redis.Values(conn.Do("HMGET", args...))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return value, err
}

// 传入的 字段列表获得对应的值
func (c RedisCache) HGetAll(name string) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	value, err := redis.Values(conn.Do("HGETALL", name))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return value, err
}

// 设置多个值 , obj 可以是指针 slice map struct
func (c RedisCache) HMSet(name string, obj interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()
	_, err = conn.Do("HMSET", redis.Args{}.Add(name).AddFlat(obj)...)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return
}

//// 设置单个值, value 还可以是一个 map slice 等
//func (c RedisCache) HSet(name string, key string, value interface{}) (err error) {
//	conn := c.pool.Get()
//	defer conn.Close()
//	v, _ := Serialization(value)
//	_, err = conn.Do("HSET", name, key, v)
//	return
//}
//
////获取单个hash 中的值
//func (c RedisCache) HGet(name, field string, v interface{}) error {
//	conn := c.pool.Get()
//	defer conn.Close()
//	temp, err := redis.Bytes(conn.Do("HGet", name, field))
//	err = Deserialization(temp, &v)
//	return err
//}

// 设置单个值, value 还可以是一个 map slice 等
func (c RedisCache) HSet(name string, key string, value interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()
	//v, _ := Serialization(value)
	_, err = conn.Do("HSET", name, key, value)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return
}

//获取单个hash 中的值
func (c RedisCache) HGet(name, field string) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()
	temp, err := redis.Bytes(conn.Do("HGet", name, field))
	//err = Deserialization(temp, &v)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return temp, err
}
func (c RedisCache) HGetEx(name, field string) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	temp, err := redis.Values(conn.Do("HGet", name, field))
	//err = Deserialization(temp, &v)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return temp, err
}

// H数值增加
func (c RedisCache) HIncrBy(name string, key string, value int64) int64 {
	conn := c.pool.Get()
	defer conn.Close()
	res, err := conn.Do("HIncrBy", name, key, value)
	PanicError(err)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return res.(int64)
}

// set 集合中添加元素
func (c RedisCache) SAdd(name string, data interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()
	m := redis.Args{}.Add(name).AddFlat(data)
	_, err = conn.Do("SAdd", m...)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return err
}

// set 集合中删除元素
func (c RedisCache) SRem(name string, data interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SRem", redis.Args{}.Add(name).AddFlat(data)...)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return err
}

// set 集合是否包含元素
func (c RedisCache) SIsMember(name string, val interface{}) bool {
	conn := c.pool.Get()
	defer conn.Close()
	b, _ := redis.Bool(conn.Do("SIsMember", name, val))
	return b
}

// 获取 set 集合中所有的元素, 想要什么类型的自己指定
func (c RedisCache) Smembers(name string) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	temp, err := redis.Values(conn.Do("smembers", name))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return temp, err
}

// 获取集合中元素的个数
func (c RedisCache) ScardInt64s(name string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Int64(conn.Do("SCARD", name))
	return v, err
}

// 传入的 字段列表获得对应的值
func (c RedisCache) LPush(name string, data interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("LPush", redis.Args{}.Add(name).AddFlat(data)...)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return err
}

// 传入的 字段列表获得对应的值
func (c RedisCache) LPop(name string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	value, err := conn.Do("LPop", name)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	if value == nil {
		return "", err
	}
	return string(value.([]byte)), err
}

func (c RedisCache) RPush(name string, data interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("RPush", redis.Args{}.Add(name).AddFlat(data)...)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return err
}

func (c RedisCache) RPop(name string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	value, err := conn.Do("RPop", name)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return string(value.([]byte)), err
}

func (c RedisCache) LLen(name string) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	value, err := redis.Int(conn.Do("LLen", name))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return value, err
}

func (c RedisCache) LRange(name string, start, end interface{}) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	value, err := redis.Values(conn.Do("LRange", name, start, end))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return value, err
}

//返回值为0代表删除失败，返回几代表删除几个元素
func (c RedisCache) LRem(name string, count, val interface{}) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	n, err := redis.Int(conn.Do("lrem", name, count, val))
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return n, err
}

//模糊查询keys
//func (c RedisCache) Keys(str string) ([]string, error) {
//	conn := c.pool.Get()
//	defer conn.Close()
//	val, err := redis.Strings(conn.Do("keys", str))
//	return val, err
//}
//模糊查询keys
func (c RedisCache) Scan(str string) ([]string, error) {
	conn := c.pool.Get()
	//*扫描所有key，每次50条
	var cursor int64
	keys := make([]string, 0)
	for {
		data, err := redis.Values(conn.Do("Scan", cursor, "match", str+"*", "count", 50))
		if err != nil {
			return keys, err
		}
		cursor = AtoInt64(string(data[0].([]uint8)))
		ids := data[1].([]interface{})

		for i := range ids {
			v := string(ids[i].([]uint8))
			keys = append(keys, v)
		}
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

//单key取得redis分布式锁没有有重试机制
func (c RedisCache) DoRedisLockNoRetry(key string, timeout int32) (err error) {
	//从连接池中娶一个con链接，pool可以自己定义。
	conn := c.pool.Get()
	defer conn.Close()
	//这里需要redis.String包一下，才能返回redis.ErrNil
	_, err = redis.String(conn.Do("set", key, 1, "ex", timeout, "nx"))
	if err != nil {

		s := fmt.Sprintf("DoRedisLockNoRetry 单key取得redis分布式无重试机制锁失败,redis key is %v", key)
		logs.Error(s)
		return err
	}
	return nil
}

//单key取得redis分布式锁有重试机制
func (c RedisCache) DoRedisLockWithRetry(key string, timeout int32) (err error) {
	//重试退出开始时间
	nowStartTime := time.Now().Unix()

	err1 := c.DoRedisLockNoRetry(key, timeout)

	if err1 != nil {

		// 需要阻塞重试取得锁
		for {
			errLoop := c.DoRedisLockNoRetry(key, timeout)

			if errLoop == nil {
				return nil
			}

			//取得锁或者重试超过8秒左右退出
			//最好要小于过期时间几秒左右,目前这个重试锁过期时间设置的可以根据业务来设置:操作设置10秒,定时设置20秒
			//过期时间保证业务在此时间内执行完毕,避免有些定时业务未执行完毕而时间过期
			if time.Now().Unix()-nowStartTime > 8 {
				s := fmt.Sprintf("DoRedisLockWithRetry 单key取得redis分布式有重试机制锁失败,redis key is %v,重试时间%v", key, time.Now().Unix()-nowStartTime)
				logs.Error(s)
				return redis.ErrNil
			}

		}
	}

	return nil
}

//单key分布式解锁
//如果是err说明业务太长 已过超时时间 key失效,所以上层业务不用关心 继续可以获得锁(根据业务调整过期时间)
func (c RedisCache) DoRedisUnlock(key string) {

	_, err := c.Delete(key)
	//如果是err说明业务太长 已过超时时间 所以上层业务不用关心 继续可以获得锁
	if err != nil {
		s := fmt.Sprintf("DoRedisUnlock 单key分布式解锁失败,可以无视错误,redis key is %v", key)
		logs.Error(s)
		logs.Error(err)
	}
}

//多key取得redis分布式锁没有有重试机制
func (c RedisCache) DoBatchRedisLockNoRetry(keys []string, timeout int32) (err error) {
	//从连接池中娶一个con链接，pool可以自己定义。
	conn := c.pool.Get()
	defer conn.Close()

	//因为传入得是字符串数组不是字符串指针数组所以用以下的循环方式
	for i := 0; i < len(keys); i++ {

		//这里需要redis.String包一下，才能返回redis.ErrNil
		_, err = redis.String(conn.Do("set", keys[i], 1, "ex", timeout, "nx"))
		if err != nil {
			return err
		}
	}
	return nil
}

//多key取得redis分布式锁有重试机制
func (c RedisCache) DoBatchRedisLockWithRetry(keys []string, timeout int32) (err error) {

	//重试退出开始时间
	nowStartTime := time.Now().Unix()
	err1 := c.DoBatchRedisLockNoRetry(keys, timeout)

	if err1 != nil {
		// 需要阻塞重试取得锁
		for {

			errLoop := c.DoBatchRedisLockNoRetry(keys, timeout)

			if errLoop == nil {
				//logs.Info(lockCnt)
				return nil
			}

			//取得锁或者重试超过8秒左右退出
			//最好要小于过期时间几秒左右,目前这个重试锁过期时间设置的可以根据业务来设置:操作设置10秒,定时设置20秒
			//过期时间保证业务在此时间内执行完毕,避免有些定时业务未执行完毕而时间过期
			if time.Now().Unix()-nowStartTime > 8 {
				s := fmt.Sprintf("DoBatchRedisLockWithRetry 多key取得redis分布式有重试机制锁失败,redis keys is %v,重试时间%v", keys, time.Now().Unix()-nowStartTime)
				logs.Error(s)
				return redis.ErrNil
			}
		}
	}

	return nil
}

//多key分布式解锁
//如果是err说明业务太长 已过超时时间 所以上层业务不用关心 继续可以获得锁(根据业务调整过期时间)
func (c RedisCache) DoBatchRedisUnlock(keys []string) {

	_, err := c.Delete(redis.Args{}.Add().AddFlat(keys)...)
	if err != nil {
		s := fmt.Sprintf("DoBatchRedisUnlock 多key分布式解锁失败,可以无视错误,redis keys is %v", keys)
		logs.Error(s)
		logs.Error(err)
	}
}

// 自增返回int64,生成分布式的时候生成订单id用不能设置过期时间
func (c RedisCache) StringIncrForInt64(name string) int64 {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Int64(conn.Do("INCR", name))
	PanicError(err)
	return v
}

// H数值增加
func (c RedisCache) IncrBy(key string, value int64) int64 {
	conn := c.pool.Get()
	defer conn.Close()
	res, err := conn.Do("IncrBy", key, value)
	PanicError(err)
	//_ = c.Expire(name, REDIS_EXPIRE_TIME)
	return res.(int64)
}

// ZADD 实现分页
func (c RedisCache) ZAdd(key string, obj interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", redis.Args{}.Add(key).AddFlat(obj)...)
	PanicError(err)
	return err

}

// ZADD 实现分页;page 当前页,pageSize 页面大小
func (c RedisCache) ZRange(key string, page, pageSize int) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	res, err := redis.Values(conn.Do("ZRANGE", key, page, pageSize))
	if err != nil {
		PanicError(err)
	}
	return res, err
}

//电竞设置key
func (c RedisCache) SetWithTime(name string, v interface{}, timeout int64) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", name, v)
	_ = c.Expire(name, timeout)
	return err
}
