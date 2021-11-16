package for_game

import (
	"game_server/easygo"

	"github.com/astaxie/beego/logs"
)

const (
	INIT_PLAYER_ID int64 = 1887436000
	INIT_TEAM_ID   int64 = 18800000
)

const (
	REDIS_SAVE_SERVER           = "save_server" //作为存储数据到mongo的大厅服务器编号
	REDIS_SAVE_HALL_SID         = "hall"        //作为存储数据到mongo的大厅服务器编号
	REDIS_SAVE_SQUARE_SID       = "square"      // 作为存储定时任务管理器的社交广场服务器编号
	REDIS_SAVE_BACKSTAGE_SID    = "backstage"   //作为储定时任务的后台服务器编号
	REDIS_SAVE_ESPORTS_CSGO_SID = "sport_csgo"  //作为储定时任务的后台服务器编号
	REDIS_SAVE_ESPORTS_DOTA_SID = "sport_dota"  //作为储定时任务的后台服务器编号
	REDIS_SAVE_ESPORTS_LOL_SID  = "sport_lol"   //作为储定时任务的后台服务器编号
	REDIS_SAVE_ESPORTS_WZRY_SID = "sport_wzry"  //作为储定时任务的后台服务器编号

	REDIS_SAVE_ESPORTS_API_SID = "sport_api" //作为储定时任务的后台服务器编号
)

// 获取当前存储任务的服务器id
func GetCurrentSaveServerSid(sid int32, redisKey string) int32 {
	key := MakeRedisKey(REDIS_SAVE_SERVER, redisKey)
	b, err := easygo.RedisMgr.GetC().Exist(key)
	easygo.PanicError(err)
	if !b {
		SetCurrentSaveServerSid(sid, redisKey)
		return sid
	}
	val := int32(0)
	err = easygo.RedisMgr.GetC().StringGet(key, &val)
	easygo.PanicError(err)
	return val
}

//设置存储数据的大厅
func SetCurrentSaveServerSid(sid int32, redisKey string) {
	err := easygo.RedisMgr.GetC().StringSet(MakeRedisKey(REDIS_SAVE_SERVER, redisKey), sid)
	easygo.PanicError(err)
	logs.Info("当前存储 key 为: %s,服务器为: %d", redisKey, sid)
}

// 获取当前存储任务的服务器id(该方法只有对api数据服务启动只接受一次有效)
func GetCurrentSaveServerSidForSportApi(sid int32, redisKey string) int32 {
	key := MakeRedisKey(REDIS_SAVE_SERVER, redisKey)
	b, err := easygo.RedisMgr.GetC().Exist(key)
	easygo.PanicError(err)
	if !b {
		return int32(0)
	}
	val := int32(0)
	err = easygo.RedisMgr.GetC().StringGet(key, &val)
	easygo.PanicError(err)
	return val
}
