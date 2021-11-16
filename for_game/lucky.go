/**
author:ax
十一幸运抽奖模块
*/
package for_game

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
	"sort"
	"sync"
	"time"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

const (
	LUCKY_PLAYER_IS_SIGN_IN     = "playerIsSignIn"      // 今天是否签到
	LUCKY_PLAYER_IS_SHARE       = "playerIsShare"       // 今天是否已分享
	LUCKY_PLAYER_IS_SENDDYNAMIC = "playerIsSendDynamic" // 今天已发布动态
	LUCKY_PLAYER_IS_SENDREDPACK = "playerIsSendRedPack" // 今天已发红包
	LUCKY_PLAYER_IS_LOGIN       = "playerIsLogin"       // 今天已登录
	LUCKY_DAY_FIRST_LUCKY       = "dayFirstLucky"       // 今日是否首次抽卡
)
const (
	REDIS_LUCKY_PLAYER      = "lucky_player"
	REDIS_FULL_LUCKY_PLAYER = "full_lucky_player"   // map[fullTime]pid
	OPEN_BEFORE_TENTIME     = "2020-10-08 19:50:00" // 开奖前10分钟
	LUCKY_START_TIME        = "2020-09-28 00:00:00"

	FULL_START_TIME  = "2020-09-27 00:00:00"
	FULL_MIDDLE_TIME = "2020-09-30 00:00:00"
	FULL_END_TIME    = "2020-10-08 19:50:00"

	FULL_COUNT = "full_count"

	OPEN_MIN_MONEY = 400  //分
	OPEN_MAX_MONEY = 1000 //分

)

type LuckyPlayer struct{}

var RedisLuckyPlayer LuckyPlayer

// 递增抽奖次数
func (lp *LuckyPlayer) IncrLuckyCountToRedis(pid, count PLAYER_ID) {
	// 先判断是否存在这个key
	lp.GetLuckyCountFromRedis(pid)
	// 该用户添加抽奖次数进数据库
	err := UpsetLuckyPlayerToDB(pid, count, 0)
	// 数据同步进redis
	easygo.RedisMgr.GetC().HIncrBy(MakeRedisKey(REDIS_LUCKY_PLAYER, pid), "LuckyCount", count)
	easygo.PanicError(err)
}

func (lp *LuckyPlayer) GetLuckyCountFromRedis(pid PLAYER_ID) int32 {
	b, err := easygo.RedisMgr.GetC().HExists(MakeRedisKey(REDIS_LUCKY_PLAYER, pid), "LuckyCount")
	easygo.PanicError(err)
	if !b { // 从数据库获取,然后设置进redis
		lp := GetLuckyPlayerFromDB(pid)
		if lp == nil {
			return 0
		}
		err := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_LUCKY_PLAYER, pid), "LuckyCount", lp.GetLuckyCount())
		easygo.PanicError(err)
		err1 := easygo.RedisMgr.GetC().Expire(MakeRedisKey(REDIS_LUCKY_PLAYER, pid), 30*86400) // 30天过期
		easygo.PanicError(err1)
		return lp.GetLuckyCount()
	}
	luckyCount, err1 := easygo.RedisMgr.GetC().HGet(MakeRedisKey(REDIS_LUCKY_PLAYER, pid), "LuckyCount")
	easygo.PanicError(err1)
	return easygo.AtoInt32(string(luckyCount))
}

func (lp *LuckyPlayer) SetLuckyPlayerInfoToRedis(pid PLAYER_ID, info map[string]interface{}) {
	err := easygo.RedisMgr.GetC().HMSet(MakeRedisKey(REDIS_LUCKY_PLAYER, pid), info)
	easygo.PanicError(err)
}

func (lp *LuckyPlayer) SetFullLuckyPlayerInfoToRedis(info map[int64]int64) {
	err := easygo.RedisMgr.GetC().HMSet(REDIS_FULL_LUCKY_PLAYER, info)
	easygo.PanicError(err)
	err1 := easygo.RedisMgr.GetC().Expire(REDIS_FULL_LUCKY_PLAYER, 30*86400) // 30天过期
	easygo.PanicError(err1)
}

// 设置集满人数进redis
func (lp *LuckyPlayer) SetFullCountToRedis(count int64) {
	easygo.RedisMgr.GetC().HIncrBy(FULL_COUNT, FULL_COUNT, count)
	err1 := easygo.RedisMgr.GetC().Expire(FULL_COUNT, 30*86400)
	easygo.PanicError(err1)
}

func (lp *LuckyPlayer) GetFullCountFromRedis() int64 {
	b, err := easygo.RedisMgr.GetC().HExists(FULL_COUNT, FULL_COUNT)
	easygo.PanicError(err)
	if !b { // 从数据库获取,然后设置进redis
		cc := GetFullCountFromDB()
		if cc.GetFullCount() > 0 { // 设计 redis
			lp.SetFullCountToRedis(cc.GetFullCount())
		}
		return cc.GetFullCount()
	}
	count, err1 := easygo.RedisMgr.GetC().HGet(FULL_COUNT, FULL_COUNT)
	easygo.PanicError(err1)
	return easygo.AtoInt64(string(count))

}

// 内部方法,抽卡,返回卡片id
func luckyCard() int64 {
	// 获取权重
	var v []*share_message.Props
	SysPropsData.RateMap.Range(func(key, value interface{}) bool {
		if easygo.GetToday0ClockTimestamp() == key.(int64) {
			v = value.([]*share_message.Props)
		}
		return true
	})
	if len(v) == 0 {
		logs.Error("luckyCard 获取权重出错,可能时间有问题,当天的整点时间为:", easygo.GetToday0ClockTimestamp())
		return -1
	}
	rate := make([]float32, 0)
	for _, value := range v {
		rate = append(rate, float32(value.GetRate())/100)
	}
	logs.Info("luckyCard 权重列表--->", rate)
	// 一万次
	var index int
	for i := 0; i < 10000; i++ {
		index = WeightedRandomIndex(rate)
	}

	id := SysPropsData.PropsSlice[index] // 抽到的卡的id
	if id == ID_QU {                     // 只有趣字才有控制
		if count := IncrDayProps(id, -1); count < 0 {
			IncrDayProps(id, 1) // 减到负数,重置为0
			//return luckyCard()
			return 2
		}
	}
	logs.Info("luckyCard 卡片索引为----->", index)
	logs.Info("luckyCard 卡片id为----->", id)
	return id
}

func ReloadSysPropsToDayProps() {
	sysList := GetSysPropsListFromDB()
	dayPropsList := make([]*share_message.DayProps, 0)
	for _, v := range sysList {
		dayPropsList = append(dayPropsList, &share_message.DayProps{
			Id:         easygo.NewInt64(v.GetId()),
			Count:      easygo.NewInt64(v.GetCount()),
			UpdateTime: easygo.NewInt64(time.Now().Unix()),
		})
	}

	var data []interface{}
	for _, v := range dayPropsList {
		b1 := bson.M{"_id": v.GetId()}
		b2 := v
		data = append(data, b1, b2)
	}
	UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_LUCKY_DAY_PROPS, data)
}

// 逐日递增集满人数
func AddFullLuckyPlayer() {
	//today := easygo.GetToday0ClockTimestamp()
	today := time.Now().Unix()
	time28 := easygo.Str2Time(FULL_START_TIME) // 09-28号,这里填写27号是下面获取28号的时间戳用的.
	unix28 := easygo.Get24ClockTimestamp(time28.Unix())
	time01 := easygo.Str2Time(FULL_MIDDLE_TIME) //10-1号
	unix01 := easygo.Get24ClockTimestamp(time01.Unix())
	time08 := easygo.Str2Time(FULL_END_TIME).Unix() //10-8号
	logs.Info("-----------今天的时间为: %d,unix28:%d,unix01:%d,unix08:%d", today, unix28, unix01, time08)
	if today >= unix28 && today <= unix01 { // 增加1-3人
		// 随机增加人数
		c := RangeRand(1, 3)
		logs.Info("增加1-3人,实际增加的人数为: ", c)
		UpsetSysFullCount(c)
		RedisLuckyPlayer.SetFullCountToRedis(c) // 设置进redis
	}
	if today > unix01 && today <= time08 { // 增加130-150人
		// 随机增加人数
		c := RangeRand(130, 150)
		logs.Info("增加130-150人,实际增加的人数为: ", c)
		UpsetSysFullCount(c)
		RedisLuckyPlayer.SetFullCountToRedis(c) // 设置进redis
	}
}

//===============以下对外方法=======================
/**
签到逻辑:
判断该用户是否存在
用户表身上加一个是否签到和抽奖次数,是否集齐了这3个字段.
判断用户是否已签到
修改签到状态
签到后的玩家id设置redis
*/
func LuckySignIn(aid, pid PLAYER_ID) *base.Fail {
	act := GetActivityFromDB(aid)
	//0开启活动，1关闭活动
	if act.GetStatus() > 0 {
		return easygo.NewFailMsg("活动已关闭")
	}
	t := time.Now().Unix()
	//   判断当前时间是否还没开始
	if act.GetStartTime() > t {
		logs.Error("活动还没开始,当前时间为: ", t)
		return easygo.NewFailMsg("活动还没开始,不能签到")
	}
	//   判断当前时间是否与开奖时间差10分钟
	if act.GetEndTime() < t {
		logs.Error("当前时间离开奖时间不足10分钟,当前时间为: ", t)
		return easygo.NewFailMsg("离开奖还有10分钟,不能签到")
	}

	if player := GetRedisPlayerBase(pid); player == nil {
		logs.Error("LuckySignIn 签到逻辑,用户不存在,playerId为: ", pid)
		return easygo.NewFailMsg("用户不存在")
	}
	period := GetPlayerPeriod(pid)
	if isSignIn := period.DayPeriod.FetchBool(LUCKY_PLAYER_IS_SIGN_IN); isSignIn {
		logs.Error("LuckySignIn 签到操作,玩家今天已经签到过了,playerId为: ", pid)
		return easygo.NewFailMsg("今天已经签到过了")
	}
	period.DayPeriod.Set(LUCKY_PLAYER_IS_SIGN_IN, true) // 修改状态已签到状态
	// 数据同步进redis和数据库
	RedisLuckyPlayer.IncrLuckyCountToRedis(pid, 1)
	return nil
}

// 获取我的背包卡片列表
func GetPlayerPropsList(pid PLAYER_ID) ([]*share_message.PlayerProps, *base.Fail) {
	if player := GetRedisPlayerBase(pid); player == nil {
		logs.Error("GetPlayerPropsList 获取我的背包卡片列表,用户不存在,playerId为: ", pid)
		return nil, easygo.NewFailMsg("用户不存在")
	}
	list := GetPlayerPropsListByPidFromDB(pid)
	return list, nil
}

/**
邀请新人
friendId:邀请人
phone: 被邀请人
有昵称并且有头像的才算是老用户,不做任何处理,直接返回
不满足上面条件的,说明是新用户,进行绑定关系.
新人进行修改个人资料,昵称,头像时,查找看有没有绑定关系,有的话,在对应的邀请人身上添加抽卡次数.
*/
func InviteNewFriend(aid, friendId PLAYER_ID, phone string) *base.Fail {
	act := GetActivityFromDB(aid)
	//0开启活动，1关闭活动
	if act.GetStatus() > 0 {
		return easygo.NewFailMsg("活动已关闭")
	}
	t := time.Now().Unix()
	//   判断当前时间是否还没开始
	if act.GetStartTime() > t {
		logs.Error("活动还没开始,当前时间为: ", t)
		return easygo.NewFailMsg("活动还没开始,不能邀请新人")
	}
	//   判断当前时间是否与开奖时间差10分钟
	if act.GetEndTime() < t {
		logs.Error("当前时间离开奖时间不足10分钟,当前时间为: ", t)
		return easygo.NewFailMsg("离开奖还有10分钟,不能邀请新人")
	}

	// 判断当前用户是否是有效用户
	if p := GetRedisPlayerBase(friendId); p == nil {
		logs.Error("InviteNewFriend 邀请新人,当前用户不存在,pId: ", friendId)
		return easygo.NewFailMsg("用户不存在")
	}
	if newFriend := GetPlayerByPhone(phone); newFriend != nil &&
		newFriend.GetHeadIcon() != "" && newFriend.GetNickName() != "" {
		logs.Error("InviteNewFriend 被邀请的人已经是柠檬的老用户,不处理逻辑,phone: ", phone)
		return easygo.NewFailMsg("邀请的用户是柠檬的老用户")
	}
	// 绑定关联关系
	pr := &share_message.LuckyPlayerRelated{
		Id:          easygo.NewString(fmt.Sprintf("%s_%s", easygo.AnytoA(friendId), phone)),
		PlayerId:    easygo.NewInt64(friendId),
		FriendPhone: easygo.NewString(phone),
		RelatedTime: easygo.NewInt64(time.Now().Unix()),
	}
	if err := UpsetLuckyPlayerRelatedToDB(pr); err != nil {
		easygo.PanicError(err)
		return easygo.NewFailMsg("邀请失败")
	}

	// 异步埋点:新增绑定关系人数.
	easygo.Spawn(func() {
		UpdateActivityReport("NewBindCount", 1)
	})
	return nil
}

/**
每日分享
*/
func LuckyShare(aid, pid PLAYER_ID) *base.Fail {
	act := GetActivityFromDB(aid)
	//0开启活动，1关闭活动
	if act.GetStatus() > 0 {
		return easygo.NewFailMsg("活动已关闭")
	}
	t := time.Now().Unix()
	//   判断当前时间是否还没开始
	if act.GetStartTime() > t {
		logs.Error("活动还没开始,当前时间为: ", t)
		return easygo.NewFailMsg("活动还没开始,不能分享")
	}
	//   判断当前时间是否与开奖时间差10分钟
	if act.GetEndTime() < t {
		logs.Error("当前时间离开奖时间不足10分钟,当前时间为: ", t)
		return easygo.NewFailMsg("离开奖还有10分钟,不能分享")
	}

	// 判断当前用户是否是有效用户
	if p := GetRedisPlayerBase(pid); p == nil {
		logs.Error("LuckyShare 每日分享,当前用户不存在,pId: ", pid)
		return easygo.NewFailMsg("用户不存在")
	}
	// 判断今天是否分享过人了
	period := GetPlayerPeriod(pid)
	if period.DayPeriod.FetchBool(LUCKY_PLAYER_IS_SHARE) {
		logs.Error("LuckyShare 今天已经分享过好友了,pId: ", pid)
		return nil
	}
	// 添加抽奖次数进数据库和redis
	RedisLuckyPlayer.IncrLuckyCountToRedis(pid, 1)
	// 设置今天已分享
	period.DayPeriod.Set(LUCKY_PLAYER_IS_SHARE, true)
	return nil
}

/**
获得当日集卡秘籍情况(玩家任务列表)
*/
func GetLuckyPlayerTask(pid PLAYER_ID) (*LuckyPlayerTask, *base.Fail) {
	// 判断当前用户是否是有效用户
	if p := GetRedisPlayerBase(pid); p == nil {
		logs.Error("GetLuckyPlayerTask 获得当日集卡秘籍情况,当前用户不存在,pId: ", pid)
		return nil, easygo.NewFailMsg("用户不存在")
	}
	period := GetPlayerPeriod(pid)
	return &LuckyPlayerTask{
		PlayerId:      pid,
		IsSignIn:      period.DayPeriod.FetchBool(LUCKY_PLAYER_IS_SIGN_IN),
		IsShare:       period.DayPeriod.FetchBool(LUCKY_PLAYER_IS_SHARE),
		IsSendDynamic: period.DayPeriod.FetchBool(LUCKY_PLAYER_IS_SENDDYNAMIC),
		IsRedpacket:   period.DayPeriod.FetchBool(LUCKY_PLAYER_IS_SENDREDPACK),
	}, nil
}

// 对外方法: 抽卡   pid: 玩家id
func LuckyCard(pid, aid PLAYER_ID) (*share_message.PlayerProps, *base.Fail) {
	act := GetActivityFromDB(aid)
	//0开启活动，1关闭活动
	if act.GetStatus() > 0 {
		return nil, easygo.NewFailMsg("活动已关闭")
	}
	t := time.Now().Unix()
	//   判断当前时间是否还没开始
	if act.GetStartTime() > t {
		logs.Error("活动还没开始,当前时间为: ", t)
		return nil, easygo.NewFailMsg("活动还没开始,不能抽卡")
	}
	//   判断当前时间是否与开奖时间差10分钟
	if act.GetEndTime() < t {
		logs.Error("当前时间离开奖时间不足10分钟,当前时间为: ", t)
		return nil, easygo.NewFailMsg("离开奖还有10分钟,不能抽卡")
	}
	// 判断当前用户是否是有效用户
	if p := GetRedisPlayerBase(pid); p == nil {
		logs.Error("LuckyCard 抽卡逻辑,当前用户不存在,pId: ", pid)
		return nil, easygo.NewFailMsg("用户不存在")
	}
	if lc := RedisLuckyPlayer.GetLuckyCountFromRedis(pid); lc <= 0 {
		logs.Error("LuckyCard 抽卡逻辑,抽卡次数小于等于0,,pId: %d,luckyCount: %d", pid, lc)
		return nil, easygo.NewFailMsg("做任务增加抽卡次数")
	}
	// 判断玩家表中是否有首次抽奖
	lp := GetLuckyPlayerFromDB(pid)
	//  异步埋点
	easygo.Spawn(func() {
		// 判断今日是否首次抽卡.
		period := GetPlayerPeriod(pid)
		if !period.DayPeriod.FetchBool(LUCKY_DAY_FIRST_LUCKY) { // 今日没有抽卡
			UpdateActivityReport("PlayerCount", 1)            //设置今日是否首次抽卡.
			period.DayPeriod.Set(LUCKY_DAY_FIRST_LUCKY, true) // 修改已抽卡
		}

		if lp == nil { // 数据库中没有,是新增的玩家,直接设置埋点数据
			_ = UpsetIsNewLuckyToDB(pid, true)        // 修改数据库是否首次
			UpdateActivityReport("NewPlayerCount", 1) // 埋点
			return
		}
		if !lp.GetIsNewLucky() {
			_ = UpsetIsNewLuckyToDB(pid, true)        // 修改数据库是否首次
			UpdateActivityReport("NewPlayerCount", 1) // 埋点
			return
		}
	})

	//  抽卡逻辑.
	cardId := luckyCard()
	if cardId < 0 {
		return nil, easygo.NewFailMsg("抽卡时间未到")
	}
	// 抽卡次数减1
	RedisLuckyPlayer.IncrLuckyCountToRedis(pid, -1)
	// 把卡片添加到用户背包.
	_ = UpsetPlayerPropsToDB(pid, cardId, 1)
	// 根据玩家id和卡片id查询玩家背包信息
	pp := GetPlayerPropsByPidAndProIdFromDB(pid, cardId)

	//  异步判断卡片是否集满了,是,存数据库,存redis
	easygo.Spawn(func() {
		if list := GetPlayerPropsListByPidFromDB(pid); len(list) >= 6 { // 集满了6张卡了.
			// 设置已集满,设置时间
			if lp.GetIsFull() {
				_ = UpsetLuckyPlayerToDB(pid, 0, time.Now().Unix())
			} else {
				_ = UpsetLuckyPlayerToDB(pid, 0, time.Now().Unix(), true)

				// 持久化redis 30 天
				m := make(map[int64]int64)
				m[time.Now().Unix()] = pid
				RedisLuckyPlayer.SetFullLuckyPlayerInfoToRedis(m)
				// 集满人数+1
				UpsetSysFullCount(1)
				RedisLuckyPlayer.SetFullCountToRedis(1) // 设置进redis
				// 添加集齐玩家埋点数
				UpdateActivityReport("FullPlayerCount", 1)
			}

		}
	})
	// 封装用户卡片数据并返回
	return &share_message.PlayerProps{
		PropsId:  easygo.NewInt64(pp.GetPropsId()),
		PlayerId: easygo.NewInt64(pp.GetPlayerId()),
		Count:    easygo.NewInt64(pp.GetCount()),
	}, nil
}

// 获取用户剩余抽卡次数
func GetLuckyCount(pid PLAYER_ID) int32 {
	return RedisLuckyPlayer.GetLuckyCountFromRedis(pid)
}

/**
赠送逻辑:
 pid: 送卡人
friendId:被送卡的人
propsId:卡片id
*/
func GiveCard(aid, pid, friendId, propsId PLAYER_ID) *base.Fail {
	act := GetActivityFromDB(aid)
	//0开启活动，1关闭活动
	if act.GetStatus() > 0 {
		return easygo.NewFailMsg("活动已关闭")
	}
	t := time.Now().Unix()
	//   判断当前时间是否还没开始
	if act.GetStartTime() > t {
		logs.Error("活动还没开始,当前时间为: ", t)
		return easygo.NewFailMsg("活动还没开始,不能送卡")
	}
	//   判断当前时间是否与开奖时间差10分钟
	if act.GetEndTime() < t {
		logs.Error("当前时间离开奖时间不足10分钟,当前时间为: ", t)
		return easygo.NewFailMsg("离开奖还有10分钟,不能送卡")
	}
	// 判断当前用户是否是有效用户
	if p := GetRedisPlayerBase(pid); p == nil {
		logs.Error("GiveCard 抽卡逻辑,当前用户不存在,pId: ", pid)
		return easygo.NewFailMsg("用户不存在")
	}
	if p := GetRedisPlayerBase(friendId); p == nil {
		logs.Error("GiveCard 抽卡逻辑,好友不存在,friendId: ", friendId)
		return easygo.NewFailMsg("好友不存在")
	}
	// 判断是否是只有一张卡
	playerProps := GetPlayerPropsByPidAndProIdFromDB(pid, propsId)
	if playerProps.GetCount() <= 1 {
		logs.Error("玩家背包卡片不足,pid: %d,卡片id: %d,当前卡片数量: %d", pid, propsId, playerProps.GetCount())
		return easygo.NewFailMsg("只有一张卡，不能赠送")
	}
	// 送卡人的对应的卡片-1
	_ = UpsetPlayerPropsToDB(pid, propsId, -1)
	// 好友的卡片背包+1
	_ = UpsetPlayerPropsToDB(friendId, propsId, 1)
	// 异步记录送卡记录.
	easygo.Spawn(func() {
		// 赠送人卡片日志
		log1 := &share_message.PlayerUsePropsLog{
			Id:          easygo.NewInt64(NextId(TABLE_LUCKY_PLAYER_USE_PROPS_LOG)),
			PropsId:     easygo.NewInt64(propsId),
			PlayerId:    easygo.NewInt64(pid),
			CreateTime:  easygo.NewInt64(time.Now().Unix()),
			Types:       easygo.NewInt32(LOG_GIVE), // 赠送
			RevPlayerId: easygo.NewInt64(friendId),
			Count:       easygo.NewInt64(1),
		}
		logList := make([]interface{}, 0)
		logList = append(logList, log1)
		InsertUsePropsLogToDB(logList)
	})

	// 判断收卡人的片是否集满了,是,存数据库,存redis
	easygo.Spawn(func() {
		if list := GetPlayerPropsListByPidFromDB(friendId); len(list) >= 6 { // 集满了6张卡了.
			// 设置已集满,设置时间
			lp := GetLuckyPlayerFromDB(friendId)
			if lp.GetIsFull() {
				_ = UpsetLuckyPlayerToDB(friendId, 0, time.Now().Unix())
			} else {
				_ = UpsetLuckyPlayerToDB(friendId, 0, time.Now().Unix(), true)

				// 持久化redis 30 天
				m := make(map[int64]int64)
				m[time.Now().Unix()] = friendId
				RedisLuckyPlayer.SetFullLuckyPlayerInfoToRedis(m)
				// 集满人数+1
				UpsetSysFullCount(1)
				RedisLuckyPlayer.SetFullCountToRedis(1) // 设置进redis
				// 添加集齐玩家埋点数
				UpdateActivityReport("FullPlayerCount", 1)
			}

		}
	})
	return nil
}

// 获取赠送列表
func GetGiveList(pid PLAYER_ID) []*share_message.PlayerUsePropsLog {
	lis := GetLogListByPid(pid, LOG_GIVE)
	ids := make([]PLAYER_ID, 0)
	for _, p := range lis {
		ids = append(ids, p.GetRevPlayerId())
	}
	pls := GetPlayerListByIds(ids)
	pmap := make(map[PLAYER_ID]*share_message.PlayerBase)
	for _, pb := range pls {
		pmap[pb.GetPlayerId()] = pb
	}
	for _, item := range lis {
		item.HeadIcon = easygo.NewString(pmap[item.GetRevPlayerId()].GetHeadIcon())
		item.NickName = easygo.NewString(pmap[item.GetRevPlayerId()].GetNickName())
	}
	return lis
}

// 获取当前集满的人数
func GetLocalFullCount() int64 {
	return RedisLuckyPlayer.GetFullCountFromRedis()
}

/**
开奖
开奖金额在7 正负3元,精确到分.
aid: 活动表的id
返回: 幸运玩家的对象,错误信息
*/
func OpenLucky(pid, aid PLAYER_ID) (*share_message.LuckyPlayer, *base.Fail) {
	// 判断玩家是否存在
	if playerBase := GetRedisPlayerBase(pid); playerBase == nil {
		logs.Error("OpenLucky 开奖逻辑,当前用户不存在,pId: ", pid)
		return nil, easygo.NewFailMsg("用户不存在")
	}

	// 判断开奖时间是否在10-8 :20:00:00 ~10-10 23:59
	act := GetActivityFromDB(aid)
	//0开启活动，1关闭活动
	if act.GetStatus() != 0 {
		return nil, easygo.NewFailMsg("活动已关闭")
	}
	t := time.Now().Unix()
	//   判断开奖是否还没开始
	logs.Info("数据库的开奖时间为----->%d,当前时间为:----->%d", act.GetOpenTime(), t)
	if act.GetOpenTime() >= t {
		logs.Error("开奖还没开始,当前时间为: ", t)
		return nil, easygo.NewFailMsg("开奖还没开始,不能开奖")
	}
	//   判断开奖是否结束了
	if act.GetCloseTime() < t {
		logs.Error("开奖结束了,当前时间为: ", t)
		return nil, easygo.NewFailMsg("开奖结束了,不能开奖")
	}
	// 判断该用户是否集卡满了
	lp := GetLuckyPlayerFromDB(pid)
	if lp.GetIsOpen() { // 判断该玩家是否已抽过奖了
		logs.Error("该用户已经抽过奖了,pid: %d,抽奖时间为: %d", pid, lp.GetOpenTime())
		return nil, easygo.NewFailMsg("您已经抽过奖了")
	}
	if !lp.GetIsFull() {
		logs.Error("该用户未集满卡片,pid: ", pid)
		return nil, easygo.NewFailMsg("未集齐用户无缘抽奖")
	}
	money := RangeRand(OPEN_MIN_MONEY, OPEN_MAX_MONEY)
	logs.Info("玩家: %d 抽到的金额为: %d", pid, money)
	lp.LuckyMoney = easygo.NewInt64(money)
	lp.IsOpen = easygo.NewBool(true)
	lp.OpenTime = easygo.NewInt64(t) // 秒
	// 修改数据库数据
	_ = UpsetLuckyMoneyToDB(lp)
	return lp, nil
}

func GetFullLuckyPlayerList() []string {
	lis := GetFullLuckyPlayerListFromDB(10)
	pids := make([]int64, 0)
	for _, p := range lis {
		pids = append(pids, p.GetPlayerId())
	}
	plist := GetPlayerListByIds(pids)
	names := make([]string, 0)
	for _, n := range plist {
		if n.GetNickName() != "" {
			names = append(names, n.GetNickName())
		}

	}

	if len(names) < 10 {
		nlis := GetRandNickNames(10 - len(names))
		names = append(names, nlis...)
	}

	sort.Strings(names)

	return names
}

// ===========================================================
type SysProps struct {
	PropsMap   sync.Map
	PropsSlice []int64
	RateMap    sync.Map
}

var SysPropsData *SysProps

func newSysProps() {
	SysPropsData = &SysProps{PropsSlice: make([]int64, 0)}
}

// 初始化系统道具包进内存
func InitPropsToMap() {
	newSysProps()
	sysPropsList := GetSysPropsListFromDB()
	for _, v := range sysPropsList {
		SysPropsData.PropsMap.Store(v.GetName(), v.GetId())
		SysPropsData.PropsSlice = append(SysPropsData.PropsSlice, v.GetId())
	}
	// 初始化概率进内存
	rateList := GetSysPropsRateFromDB()

	for _, v := range rateList {
		SysPropsData.RateMap.Store(v.GetCreateTime(), v.GetRate())
	}
}

/*//=======系统自己生成的当前集满人数==========
//活动报表 fild字段名,val变化数量, createTime当天的报表数据更改可以不传时间
func UpdateActivityReport(fild string, val int64, createTime ...int64) {
	cTime := easygo.GetToday0ClockTimestamp()
	if len(createTime) > 0 {
		cTime = easygo.Get0ClockTimestamp(createTime[0])
	}

	FindAndModify(MONGODB_NINGMENG, TABLE_ACTIVITY_REPORT, bson.M{"_id": cTime}, bson.M{"$inc": bson.M{fild: val}}, true)
}
*/
