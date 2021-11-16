package for_game

import (
	"game_server/easygo"

	"github.com/garyburd/redigo/redis"
)

/*
redis后台管理员管理器
*/
const (
	ADMIN_ONLINE_LIST = "admin_online_list"
)

type RedisAdmin struct {
	UserId    int64  //ID
	Role      int32  //管理员类型 0超管，1管理员
	ServerId  int32  //服务器Id
	Timestamp int64  //登录时间戳
	Token     string //登录token
}

//添加管理员到Redis
func SetRedisAdmin(obj *RedisAdmin) {
	//js, err := json.Marshal(obj)
	//easygo.PanicError(err)
	err1 := easygo.RedisMgr.GetC().HMSet(MakeRedisKey(ADMIN_ONLINE_LIST, obj.UserId), obj)
	easygo.PanicError(err1)
}

//查询管理员列表
func GetRedisAdminList() []*RedisAdmin {
	var lst []*RedisAdmin
	keys, err := easygo.RedisMgr.GetC().Scan(ADMIN_ONLINE_LIST)
	if err != nil {
		return lst
	}
	for _, key := range keys {
		value, err := easygo.RedisMgr.GetC().HGetAll(key)
		if err != nil {
			continue
		}
		var admin RedisAdmin
		err = redis.ScanStruct(value, &admin)
		if err != nil {
			continue
		}
		lst = append(lst, &admin)
	}
	return lst
}

//根据id查询
func GetRedisAdmin(UserId int64) *RedisAdmin {
	val, err := easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(ADMIN_ONLINE_LIST, UserId))
	if err != nil {
		return nil
	}

	var admin RedisAdmin
	err = redis.ScanStruct(val, &admin)
	if err != nil {
		return nil
	}
	return &admin
}

//删除Redis中的数据
func DelRedisAdmin(id int64) {
	_, err := easygo.RedisMgr.GetC().Delete(MakeRedisKey(ADMIN_ONLINE_LIST, id))
	easygo.PanicError(err)
}

//删除指定服务器客服
func DelRedisAdminForSid(id int32) {
	res, err := easygo.RedisMgr.GetC().Exist(ADMIN_ONLINE_LIST)
	easygo.PanicError(err)
	if !res {
		return
	}

	lis := GetRedisAdminList()
	keys := make([]interface{}, 0)
	for _, v := range lis {
		if v.ServerId == id {
			keys = append(keys, MakeRedisKey(ADMIN_ONLINE_LIST, v.UserId))
		}
	}
	if len(keys) > 0 {
		_, err = easygo.RedisMgr.GetC().Delete(keys...)
		easygo.PanicError(err)
	}
}
