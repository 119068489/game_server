package for_game

import (
	"encoding/json"
	"game_server/easygo"

	"github.com/astaxie/beego/logs"
)

/*
redis后台客服管理器
*/
const (
	WAITER_ONLINE_LIST = "waiter_online_list"
)

type RedisWaiter struct {
	UserId    int64 //客服ID
	Role      int   //客服类型 1普通客服，2主管客服
	ConnCount int   //接待数
	ServerId  int32 //服务器Id
	Status    int32 //状态 0接待，1休息
	Types     int32 //客服分类
}

//添加在线客服到Redis
func SetRedisWaiter(obj *RedisWaiter) {
	js, err := json.Marshal(obj)
	easygo.PanicError(err)
	err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(WAITER_ONLINE_LIST), easygo.AnytoA(obj.UserId), string(js))
	easygo.PanicError(err1)
}

//查询客服列表
func GetRedisWaiterList() []*RedisWaiter {
	var lst []*RedisWaiter
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(WAITER_ONLINE_LIST))
	easygo.PanicError(err)
	for _, m := range values {
		waiter := &RedisWaiter{}
		json.Unmarshal([]byte(m), &waiter)
		lst = append(lst, waiter)
	}
	return lst
}

//根据id查询
func GetRedisWaiter(UserId int64) *RedisWaiter {
	waiter := &RedisWaiter{}
	val, err := easygo.RedisMgr.GetC().HGet(WAITER_ONLINE_LIST, easygo.AnytoA(UserId))
	if err == nil {
		json.Unmarshal(val, &waiter)
		return waiter
	} else {
		return nil
	}
}

//设置客服工作状态 0接待，1休息
func SetRedisWaiterStatus(UserId, val int64) {
	waiter := GetRedisWaiter(UserId)
	if waiter == nil {
		return
	}

	waiter.Status = int32(val)
	SetRedisWaiter(waiter)
}

//修改客服链接数 count 变化的数量
func UpdateRedisWaiterCount(userId int64, count int) {
	waiter := GetRedisWaiter(userId)
	if waiter != nil {
		waiter.ConnCount += count
		if waiter.ConnCount < 0 {
			waiter.ConnCount = 0
		}
		SetRedisWaiter(waiter)
	}
}

//重置客服连接数 count 实际的数量
func ReloadRedisWaiterCount(userId int64, count int) {
	waiter := GetRedisWaiter(userId)
	if waiter != nil {
		waiter.ConnCount = count
		SetRedisWaiter(waiter)
	}
}

//删除Redis中的数据
func DelRedisWaiter(id int64) {
	_, err := easygo.RedisMgr.GetC().Hdel(WAITER_ONLINE_LIST, easygo.AnytoA(id))
	easygo.PanicError(err)
}

//删除指定服务器客服
func DelRedisWaiterForSid(id int32) {
	res, err := easygo.RedisMgr.GetC().Exist(WAITER_ONLINE_LIST)
	easygo.PanicError(err)
	if !res {
		return
	}

	lis := GetRedisWaiterList()
	var keys []string
	for _, v := range lis {
		if v.ServerId == id {
			keys = append(keys, easygo.AnytoA(v.UserId))
		}
	}

	b, _ := easygo.RedisMgr.GetC().Hdel(WAITER_ONLINE_LIST, keys...)
	if !b {
		logs.Info("删除错误:", keys)
	}

}

//分配接收消息的客服Id
func GetWaiterOnline(types int32) int64 {
	Waiters := GetRedisWaiterList()
	var wlist []*RedisWaiter
	for _, li := range Waiters {
		if li.Status == 0 && li.Types == types {
			wlist = append(wlist, li)
		}
	}

	var userId int64
	count := 0
	list := make(map[int64]*RedisWaiter)
	for _, v := range wlist {
		list[v.UserId] = v
	}

	for _, v := range list {
		if v.ConnCount == 0 {
			userId = v.UserId

			break
		} else if v.ConnCount >= 10 {
			continue
		}

		if v.ConnCount < count || count == 0 {
			userId = v.UserId
			count = v.ConnCount

		}
	}

	return userId
}
