package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/client_hall"
	"game_server/pb/client_square"
	"game_server/pb/share_message"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

/*
社交广场 redis数据
author : 狗哥
*/
//todo  删除redis数据机制

const (
	DYNAMIC_REQUEST_NUM        = 20
	DYNAMIC_STATUE_COMMON      = 0 //正常的动态
	DYNAMIC_STATUE_DELETE      = 1 //后台删除
	DYNAMIC_STATUE_APP_DELETE  = 2 //前端删除
	DYNAMIC_STATUE_UNPUBLISHED = 3 //未发布
	DYNAMIC_STATUE_EXPIRED     = 4 //已过期

	DYNAMIC_TOP_NUM      = 3 // 置顶限制数量
	DYNAMIC_CHECK_OK     = 1 //审核状态 0未处理,1已审核,2已拒绝，3自动审核
	DYNAMIC_CHECK_REFUSE = 2
	DYNAMIC_CHECK_AUTOOK = 3
)

const (
	REDIS_SQUARE_OF_TOP_DYNAMIC = "redis_square:of_top_dynamic" // 所有官方置顶的动态.
	REDIS_SQUARE_BS_TOP_DYNAMIC = "redis_square:bs_top_dynamic" // 所有后台置顶的动态.
	REDIS_SQUARE_TOP_DYNAMIC    = "redis_square:top_dynamic"    // 所有app置顶的动态.

	//REDIS_SQUARE_DYNAMIC = "redis_square_dynamic" //所有的社交广场发布动态 map[redis_square_dynamic]map[logid]string  // 原来的,暂时不删
	REDIS_SQUARE_DYNAMIC1 = "redis_square:dynamic" //所有的社交广场发布动态 map[redis_square_dynamic]map[logid]string
	//REDIS_SQUARE_LOGIDS          = "redis_square_logids"          //当前redis中存在的动态id
	REDIS_SQUARE_MESSAGEIDS  = "redis_square:messageids"  //收到的回复消息列表idmap[redis_square_messageids_pid]map[logId][]id
	REDIS_SQUARE_UNREADINFO  = "redis_square:unreadinfo"  //未读消息信息
	REDIS_SQUARE_UNIQUE_CODE = "redis_square:unique_code" // 前端发布动态的时候传递的唯一标识.断网的时候使用校验.
	REDIS_SQUARE_MYZANIDS    = "redis_square:myzanids"    //我点过赞的动态id

	//存储限制key
	REDIS_SQUARE_SAVEINFO  = "redis_square:saveinfo" //上次存库的最后id记录
	REDIS_SQUARE_SAVELIST  = "redis_square:savelist"
	REDIS_SQUARE_MAX_LOGID = "redis_square:max_logid"

	//话题贡献榜
	TOPIC_DEVOTE = "topic_devote" //话题贡献度
)

const (
	OF_TOP_NUM  = 1 // 官方置顶的数量
	BS_TOP_NUM  = 2 // 后台置顶的数量
	APP_TOP_NUM = 2 // app置顶的数量

	PLAYER_TOP_NUM = 3 // 每个用户最多置顶的条数
)

const (
	DATA_TYPE_DYNAMIC = 0 // 代表是动态
	DATA_TYPE_ADV     = 1 // 代表是广告
	DATA_TYPE_TOPIC   = 2 // 代表是话题
)

const (
	SQUARE_ATTENTION = 2 // 关注
	SQUARE_DYNAMIC   = 1 // 广场
)

// 动态广场广告插入的位置,5-10条
const (
	ADV_MIN_INDEX = 5
	ADV_MAX_INDEX = 10
)

var TimerMgr *SquareTimerMgr

const (
	HOT_TYPE_0 = 0 // 没有热门
	HOT_TYPE_1 = 1 // 普通热门
)

// SquareTimerMgr 社交广场置顶定时任务管理器
type SquareTimerMgr struct {
	TimerMap sync.Map
}

func NewSquareTimerMgr() *SquareTimerMgr {
	TimerMgr = new(SquareTimerMgr)
	return TimerMgr
}

func SquareRecoverAndLog(ch chan int) { //保证有panic时  通道能被退出而不会阻塞
	recoverVal := recover()
	if recoverVal != nil {
		easygo.LogPanicAndStack(recoverVal)
	}
	if ch != nil {
		ch <- 0
	}

}

/**
xiong
程序初始化加载100条到redis中
*/
func ReloadSquareInfo() {

	keys, err := easygo.RedisMgr.GetC().Scan(REDIS_SQUARE_DYNAMIC1)
	easygo.PanicError(err)
	if len(keys) > 0 {
		return
	}

	lst := GetDynamicListByLimitFromDB(DataLimit) // 100 条
	if len(lst) > 0 {
		var logIds []int64 //所有动态id
		info := make(map[int64]string)
		for _, m := range lst {
			logId := m.GetLogId()
			s, _ := json.Marshal(m)
			info[logId] = string(s)
			logIds = append(logIds, logId)
		}

		easygo.Spawn(func() {
			for _, value := range lst {
				SetSquareDynamicToRedis(value)
			}
		})

		easygo.Spawn(func() { ReloadDynamicComment(logIds, nil) }) //加载logIds的所有评论
	}
}

/**
xiong
社交广场登录时
加载自己的动态
*/
func ReloadMyDynamicInfo(playerId int64) {
	player := GetRedisPlayerBase(playerId)
	if player == nil {
		return
	}
	if player.GetIsRelodSquare() {
		return
	}
	player.SetIsRelodSquare(true)
	lst, err := GetDynamicListByPlayerIdAndStatusFromDB(playerId, DYNAMIC_STATUE_COMMON)
	easygo.PanicError(err)
	if len(lst) == 0 { //如果他没有发布过动态  就加载他的点赞信息和被关注信息
		easygo.Spawn(func() { ReloadPlayerAttention(playerId, nil) }) //加载玩家被关注信息
		easygo.Spawn(func() { ReloadPlayerZanIds(playerId, nil) })    //加载自己赞过的动态id
		return
	}
	logIds := make([]int64, 0)
	// 记录不存在的动态
	ds := make([]*share_message.DynamicData, 0)
	for _, log := range lst {
		logId := log.GetLogId()
		logIds = append(logIds, logId)
		b, err := CheckHasInRedis(MakeRedisKey(REDIS_SQUARE_DYNAMIC1, logId)) //如果已经在redis中 就不加载
		if err != nil {
			logs.Error(err)
			continue
		}
		if !b {
			ds = append(ds, log)
		}
	}
	if len(ds) != 0 {
		easygo.Spawn(func() {
			for _, v := range ds {
				SetSquareDynamicToRedis(v)
			}
		})
	}

	ch := make(chan int, 4)
	easygo.Spawn(func() { ReloadMyCommentInfo(logIds, ch, playerId) }) //加载newIds的所有评论
	easygo.Spawn(func() { ReloadDynamicZan(logIds, ch) })              //加载这些动态被赞的信息
	easygo.Spawn(func() { ReloadPlayerAttention(playerId, ch) })       //加载玩家被关注信息
	easygo.Spawn(func() { ReloadPlayerZanIds(playerId, ch) })          //加载自己赞过的动态id
	for i := 0; i < 4; i++ {
		<-ch
	}
	close(ch)
}

/**
xiong
加载logIds中的动态
*/
func ReloadSomeSquareInfo(logIds []int64) {
	lst, err := GetDynamicListByLogIdsAndStatusFromDB(logIds, DYNAMIC_STATUE_COMMON)
	easygo.PanicError(err)
	if len(lst) == 0 {
		return
	}
	var newLogIds []int64 //所有动态id
	ds := make([]*share_message.DynamicData, 0)
	for _, m := range lst {
		logId := m.GetLogId()
		ds = append(ds, m)
		newLogIds = append(newLogIds, logId)
	}
	if len(ds) > 0 {
		easygo.Spawn(func() {
			for _, v := range ds {
				SetSquareDynamicToRedis(v)
			}
		})
	}

	ch := make(chan int)
	easygo.Spawn(func() { ReloadDynamicComment(newLogIds, ch) }) //加载logIds的所有评论
	<-ch
	close(ch)
}

/**
获取redis中动态id列表(不包含置顶的动态)
*/
func GetRedisSquareLogIds() []int64 {
	ids := make([]int64, 0)
	keys, err1 := easygo.RedisMgr.GetC().Scan(REDIS_SQUARE_DYNAMIC1)
	easygo.PanicError(err1)
	for _, key := range keys {
		id := strings.SplitN(key, ":", 3)[2] // redis_square:dynamic:1
		ids = append(ids, easygo.AtoInt64(id))
	}

	// 排除置顶的动态
	topIds := GetAllTopDynamicIds()
	if len(topIds) == 0 {
		easygo.SortSliceInt64(ids, false)
		return ids
	}
	noTopIds := make([]int64, 0)
	noTopIds = util.Slice1DelSlice2(ids, topIds)
	easygo.SortSliceInt64(noTopIds, false)
	return noTopIds
}

// 获取所有的置顶动态id列表
func GetAllTopDynamicIds() []int64 {
	bsTopIds := make([]int64, 0) // 后台置顶
	redisCon := easygo.RedisMgr.GetC()
	value, err1 := redisCon.HKeys(REDIS_SQUARE_BS_TOP_DYNAMIC)
	easygo.PanicError(err1)
	InterfersToInt64s(value, &bsTopIds)

	appTopIds := make([]int64, 0) // app 置顶
	value, err1 = redisCon.HKeys(REDIS_SQUARE_TOP_DYNAMIC)
	easygo.PanicError(err1)
	InterfersToInt64s(value, &appTopIds)

	OFTopIds := make([]int64, 0) // 官方置顶
	value, err1 = redisCon.HKeys(REDIS_SQUARE_OF_TOP_DYNAMIC)
	easygo.PanicError(err1)
	InterfersToInt64s(value, &OFTopIds)

	allTop := make([]int64, 0)
	allTop = append(allTop, bsTopIds...)
	allTop = append(allTop, appTopIds...)
	allTop = append(allTop, OFTopIds...)
	return allTop
}

//func GetRedisSquareMaxId() int64 { //获取当前最大的动态id
//	ids := GetRedisSquareLogIds()
//	if len(ids) == 0 {
//		return 0
//	}
//	return util.Slice64MaxInt(ids)
//}

/**
xiong
内部逻辑使用,判断该动态是否存在
如果不存在redis就去数据库中找  如果数据库中有就加载到redis
*/
func RedisDynamicIsExits(logId int64) bool {
	obj := GetRedisSquareDynamic(logId)
	if obj == nil { //redis内存中没有  去数据库找一下有没有
		obj1, err := GetDynamicByLogIdFromDB(logId)
		if err != nil && err != mgo.ErrNotFound {
			easygo.PanicError(err)
		}
		if err == mgo.ErrNotFound {
			return false
		}
		if obj1 == nil {
			return false
		}
		if obj1.GetStatue() == DYNAMIC_STATUE_DELETE {
			return false
		}

		UpdateRedisSquareDynamic(obj1)
		ch := make(chan int)
		easygo.Spawn(func() { ReloadDynamicComment([]int64{logId}, ch) }) //加载logIds的所有评论
		<-ch
		close(ch)
	}
	return true
}

/**
xiong
间接调用
从redis中获取一条动态
*/
func GetRedisSquareDynamic(logId int64) *share_message.DynamicData { //获取其中一条动态数据
	b, err := CheckHasInRedis(MakeRedisKey(REDIS_SQUARE_DYNAMIC1, logId))
	easygo.PanicError(err)
	if !b {
		return nil
	}
	value, err1 := easygo.RedisMgr.GetC().HGet(MakeRedisKey(REDIS_SQUARE_DYNAMIC1, logId), easygo.AnytoA(logId))
	easygo.PanicError(err1)
	var obj *share_message.DynamicData
	err2 := json.Unmarshal(value, &obj)
	easygo.PanicError(err2)
	if obj.GetStatue() != DYNAMIC_STATUE_COMMON {
		return nil
	}
	return obj
}

/**
xiong
对外接口
从redis中获取一条动态
*/
func GetRedisDynamic(logId int64) *share_message.DynamicData { //检查有没有在redis中 没有就去数据库找 找不到就不存在
	b := RedisDynamicIsExits(logId)
	if !b {
		return nil
	}
	dynamic := GetRedisSquareDynamic(logId)
	if player := GetRedisPlayerBase(dynamic.GetPlayerId()); player != nil {
		dynamic.Types = easygo.NewInt32(player.GetTypes())
	}
	return dynamic
}

/**
xiong
增加一条动态  立马存库并且加载到redis中
*/
func AddRedisSquareDynamic(msg *share_message.DynamicData) {
	playerMgr := GetRedisPlayerBase(msg.GetPlayerId())
	if playerMgr != nil {
		if playerMgr.GetCheckNum() >= 50 {
			msg.Check = easygo.NewInt32(DYNAMIC_CHECK_AUTOOK) //白名单自动审核动态
		}
		msg.SenderType = easygo.NewInt32(playerMgr.GetTypes())
	}

	err := InsertDynamicToDB(msg)
	easygo.PanicError(err)
	UpdateRedisSquareDynamic(msg)
}

/**
xiong
删除动态
内部逻辑调用
*/
func DelSquareDynamicById(logId int64) bool {
	ok, err := easygo.RedisMgr.GetC().Delete(MakeRedisKey(REDIS_SQUARE_DYNAMIC1, easygo.AnytoA(logId)))
	easygo.PanicError(err)
	return ok
}

/**
xiong
对外接口
删除动态
*/
func DelRedisSquareDynamic(pid, logId int64, note string, status ...int32) bool { //删除一条动态
	log := GetRedisDynamic(logId)
	if log == nil {
		logs.Error("不存在动态，id:" + easygo.AnytoA(logId))
		return false
	}
	if log.GetPlayerId() != pid {
		logs.Error("不是你发布的id，不能删除，id:" + easygo.AnytoA(logId))
		return false
	}
	// 修改数据库状态
	var delStatus int32 = DYNAMIC_STATUE_APP_DELETE
	if len(status) > 0 {
		delStatus = status[0]
	}
	if err := UpSertDynamicStatusByLogIdFromDB(logId, delStatus, note); err != nil {
		logs.Error("数据库删除动态失败id:", logId, err)
		return false
	}
	DelSquareDynamicById(logId)
	DelDynamicComment(logId)
	DelDynamicMessage(pid, logId)
	DelDynamicZan(logId)
	DelDynamicZanInfo(pid, logId)
	return true
}

/**
xiong
删除指定用户对某条动态的消息id
*/
func DelDynamicMessage(pid, logId int64) {
	redisKey := MakeRedisKey(REDIS_SQUARE_MESSAGEIDS, easygo.AnytoA(pid))
	b, err := easygo.RedisMgr.GetC().HExists(redisKey, easygo.AnytoA(logId))
	easygo.PanicError(err)
	if !b {
		return
	}
	b1, err4 := easygo.RedisMgr.GetC().Hdel(redisKey, easygo.AnytoA(logId)) //从redis中删除消息id
	easygo.PanicError(err4)
	if !b1 {
		logs.Error("redis删除消息失败id:", logId)
	}
}

/**
xiong
更新动态信息进redis
*/
func UpdateRedisSquareDynamic(obj *share_message.DynamicData) {
	SetSquareDynamicToRedis(obj)
}

// isTop 是否是置顶的数据.
func GetRedisDynamicForSomeLogId(isTop bool, pid int64, logIds []int64, opId ...int64) []*share_message.DynamicData {
	lst := make([]*share_message.DynamicData, 0)
	if len(logIds) == 0 {
		return lst
	}
	player := GetRedisPlayerBase(pid)
	if player == nil {
		logs.Error("GetRedisDynamicForSomeLogId 玩家对象怎么会为空")
		return lst
	}

	name := player.GetNickName()
	headIcon := player.GetHeadIcon()
	sex := player.GetSex()
	if !isTop {
		easygo.SortSliceInt64(logIds, false) //降序排序
	}

	var num int
	values := make(map[int64][]byte)
	redisCon := easygo.RedisMgr.GetC()
	for _, id := range logIds {
		result, err := redisCon.HGet(MakeRedisKey(REDIS_SQUARE_DYNAMIC1, id), easygo.AnytoA(id))
		if err != nil {
			//easygo.PanicError(err)
			continue
		}
		values[id] = result
	}

	var zanPid int64 // 不同场景进来,操作者不一样
	zanPid = append(opId, pid)[0]
	newLst := make([]int64, 0)
	for index, m := range values {
		if m == nil {
			newLst = append(newLst, index)
			continue
		}
		if num > DYNAMIC_REQUEST_NUM {
			break
		}
		var log *share_message.DynamicData
		_ = json.Unmarshal(m, &log)
		logId := log.GetLogId()
		log.CommentNum = easygo.NewInt64(GetRedisDynamicCommentNum(logId))
		log.Zan = easygo.NewInt32(GetRedisDynamicZanNum(logId) + log.GetTrueZan())
		log.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, zanPid))
		log.NickName = easygo.NewString(name)
		log.HeadIcon = easygo.NewString(headIcon)
		log.Sex = easygo.NewInt32(sex)
		lst = append(lst, log)
	}
	if len(newLst) > 0 { // 去数据库查询
		ds := make([]*share_message.DynamicData, 0)
		col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
		defer closeFun()
		err := col.Find(bson.M{"_id": bson.M{"$in": newLst}, "Statue": 0}).All(&ds) //未读消息
		easygo.PanicError(err)
		for _, v := range ds {
			logId := v.GetLogId()
			v.CommentNum = easygo.NewInt64(GetRedisDynamicCommentNum(logId))
			v.Zan = easygo.NewInt32(GetRedisDynamicZanNum(logId) + v.GetTrueZan())
			v.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, zanPid))
			v.NickName = easygo.NewString(name)
			v.HeadIcon = easygo.NewString(headIcon)
			v.Sex = easygo.NewInt32(sex)
			lst = append(lst, v)
		}
	}

	sort.Slice(lst, func(i, j int) bool {
		return lst[i].GetLogId() > lst[j].GetLogId()
	})

	return lst
}

// isTop 是否是置顶的数据.
func GetRedisDynamicForSomeLogId2(pid int64, ds []*share_message.DynamicData, opId ...int64) []*share_message.DynamicData {
	player := GetRedisPlayerBase(pid)
	if player == nil {
		panic("玩家对象怎么会为空")
	}
	lst := make([]*share_message.DynamicData, 0)
	if len(ds) == 0 {
		return lst
	}
	name := player.GetNickName()
	headIcon := player.GetHeadIcon()
	sex := player.GetSex()

	var zanPid int64 // 不同场景进来,操作者不一样
	zanPid = append(opId, pid)[0]
	for _, log := range ds {
		logId := log.GetLogId()
		log.CommentNum = easygo.NewInt64(GetRedisDynamicCommentNum(logId))
		log.Zan = easygo.NewInt32(GetRedisDynamicZanNum(logId) + log.GetTrueZan())
		log.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, zanPid))
		log.NickName = easygo.NewString(name)
		log.HeadIcon = easygo.NewString(headIcon)
		log.Sex = easygo.NewInt32(sex)
		lst = append(lst, log)
	}
	return lst
}

//func GetRedisDynamicForSomeLogId1(pid int64, logIds []int64, playerId ...int64) []*share_message.DynamicData {
//	player := GetRedisPlayerBase(pid)
//	if player == nil {
//		panic("玩家对象怎么会为空")
//	}
//	lst := make([]*share_message.DynamicData, 0)
//	if len(logIds) == 0 {
//		return lst
//	}
//	var name, headIcon string
//	if len(playerId) > 0 {
//		p := GetRedisPlayerBase(playerId[0])
//		name = p.GetNickName()
//		headIcon = p.GetHeadIcon()
//	} else {
//		name = player.GetNickName()
//		headIcon = player.GetHeadIcon()
//	}
//
//	sex := player.GetSex()
//	//easygo.SortSliceInt64(logIds, false) //降序排序
//	var ids []string
//	for _, id := range logIds {
//		ids = append(ids, easygo.AnytoA(id))
//	}
//
//	var num int
//	values, err := easygo.RedisMgr.GetC().HMGet(REDIS_SQUARE_DYNAMIC, ids...)
//	easygo.PanicError(err)
//	for _, m := range values {
//		if m == nil {
//			continue
//		}
//		if num > DYNAMIC_REQUEST_NUM {
//			break
//		}
//		var log *share_message.DynamicData
//		_ = json.Unmarshal(m.([]byte), &log)
//		logId := log.GetLogId()
//		log.CommentNum = easygo.NewInt64(GetRedisDynamicCommentNum(logId))
//		log.Zan = easygo.NewInt32(GetRedisDynamicZanNum(logId) + log.GetTrueZan())
//		log.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, pid))
//		log.NickName = easygo.NewString(name)
//		log.HeadIcon = easygo.NewString(headIcon)
//		log.Sex = easygo.NewInt32(sex)
//		lst = append(lst, log)
//	}
//	return lst
//}
func GetRedisDynamicForSomeLogId1(pid int64, ds []*share_message.DynamicData, playerId ...int64) []*share_message.DynamicData {
	player := GetRedisPlayerBase(pid)
	if player == nil {
		panic("玩家对象怎么会为空")
	}
	lst := make([]*share_message.DynamicData, 0)
	if len(ds) == 0 {
		return lst
	}
	var name, headIcon string
	var sex int32
	if len(playerId) > 0 {
		p := GetRedisPlayerBase(playerId[0])
		name = p.GetNickName()
		headIcon = p.GetHeadIcon()
		sex = p.GetSex()
	} else {
		name = player.GetNickName()
		headIcon = player.GetHeadIcon()
		sex = player.GetSex()
	}
	for _, log := range ds {
		logId := log.GetLogId()
		log.CommentNum = easygo.NewInt64(GetRedisDynamicCommentNum(logId))
		log.Zan = easygo.NewInt32(GetRedisDynamicZanNum(logId) + log.GetTrueZan())
		log.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, pid))
		log.NickName = easygo.NewString(name)
		log.HeadIcon = easygo.NewString(headIcon)
		log.Sex = easygo.NewInt32(sex)
		lst = append(lst, log)
	}
	return lst
}

/**
xiong
时间约新的动态越靠前
*/
func SortDynamicSliceByTime(s []int64) []int64 {
	m := make(map[int64]int64)
	timeSlice := make([]int64, 0)
	for _, logId := range s {
		if dynamic := GetRedisDynamic(logId); dynamic != nil {
			m[dynamic.GetCreateTime()] = dynamic.GetLogId()
			timeSlice = append(timeSlice, dynamic.GetCreateTime())
		}
	}
	//时间排序
	easygo.SortSliceInt64(timeSlice, false)
	//重新得到排序后的logId切片
	result := make([]int64, 0)
	for _, v := range timeSlice {
		result = append(result, m[v])
	}
	return result
}

/**
xiong
时间约新的动态越靠前
*/
func SortDynamicSliceByTime1(s []*share_message.DynamicData) []*share_message.DynamicData {
	m := make(map[int64]*share_message.DynamicData)
	timeSlice := make([]int64, 0)
	for _, dynamic := range s {
		m[dynamic.GetSendTime()] = dynamic
		timeSlice = append(timeSlice, dynamic.GetSendTime())
	}
	//时间排序
	easygo.SortSliceInt64(timeSlice, false)
	//重新得到排序后的logId切片
	result := make([]*share_message.DynamicData, 0)
	for _, v := range timeSlice {
		result = append(result, m[v])
	}
	return result
}

// todo 备份修改前的
//func GetRedisSomeDynamic1(pid int64, dsList []*share_message.DynamicData, attentionList []int64) []*share_message.DynamicData { //获取redis中的数据
//	var playerIds []int64
//	lst := make([]*share_message.DynamicData, 0)
//	for _, log := range dsList {
//		if log.GetIsShield() && pid != log.GetPlayerId() { //如果是屏蔽状态并且不是自己请求动态数据  就跳过
//			continue
//		}
//		if util.Int64InSlice(log.GetPlayerId(), attentionList) {
//			log.IsAtten = easygo.NewBool(true)
//		}
//		logId := log.GetLogId()
//		playerIds = append(playerIds, log.GetPlayerId())
//
//		log.CommentNum = easygo.NewInt64(GetRedisDynamicCommentNum(logId))
//		logs.Info("444444444444444")
//		log.Zan = easygo.NewInt32(GetRedisDynamicZanNum(logId) + log.GetTrueZan())
//		logs.Info("5555555555555555")
//		log.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, pid))
//		lst = append(lst, log)
//		// todo 弥补senderType字段,此字段是后续加的
//		if log.GetSenderType() == 0 {
//			if p := GetRedisPlayerBase(log.GetPlayerId()); p != nil {
//				log.SenderType = easygo.NewInt32(p.GetTypes())
//				easygo.Spawn(UpdateDynamicSenderTypeToDB, log.GetLogId(), log.GetSenderType())
//			}
//		}
//	}
//
//	playerInfo := GetAllPlayerBase(playerIds)
//	for _, log := range lst {
//		pid := log.GetPlayerId()
//		player := playerInfo[pid]
//		log.NickName = easygo.NewString(player.GetNickName())
//		log.HeadIcon = easygo.NewString(player.GetHeadIcon())
//		log.Sex = easygo.NewInt32(player.GetSex())
//	}
//	return lst
//}

func GetRedisSomeDynamic1(pid int64, dsList []*share_message.DynamicData, attentionList []int64) []*share_message.DynamicData { //获取redis中的数据
	var playerIds []int64
	lst := make([]*share_message.DynamicData, 0)

	logIds := make([]int64, 0) // 动态id,不重复
	for _, log := range dsList {
		logId := log.GetLogId()
		logIds = append(logIds, logId)
	}
	mm := GetRedisDynamicZanNumEx1(logIds)
	comMap := GetRedisDynamicCommentNumEx(logIds)
	for _, log := range dsList {
		if log.GetIsShield() && pid != log.GetPlayerId() { //如果是屏蔽状态并且不是自己请求动态数据  就跳过
			continue
		}
		if util.Int64InSlice(log.GetPlayerId(), attentionList) {
			log.IsAtten = easygo.NewBool(true)
		}
		logId := log.GetLogId()
		playerIds = append(playerIds, log.GetPlayerId())
		//log.CommentNum = easygo.NewInt64(GetRedisDynamicCommentNum(logId))
		log.CommentNum = easygo.NewInt64(comMap[logId])
		//log.Zan = easygo.NewInt32(GetRedisDynamicZanNum(logId) + log.GetTrueZan())
		log.Zan = easygo.NewInt32(mm[logId] + log.GetTrueZan())
		log.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, pid))
		lst = append(lst, log)
		// todo 弥补senderType字段,此字段是后续加的
		if log.GetSenderType() == 0 {
			if p := GetRedisPlayerBase(log.GetPlayerId()); p != nil {
				log.SenderType = easygo.NewInt32(p.GetTypes())
				easygo.Spawn(UpdateDynamicSenderTypeToDB, log.GetLogId(), log.GetSenderType())
			}
		}
	}
	playerInfo := GetAllPlayerBase(playerIds)
	for _, log := range lst {
		pid := log.GetPlayerId()
		player := playerInfo[pid]
		log.NickName = easygo.NewString(player.GetNickName())
		log.HeadIcon = easygo.NewString(player.GetHeadIcon())
		log.Sex = easygo.NewInt32(player.GetSex())
		log.Types = easygo.NewInt32(player.GetTypes())
	}
	return lst
}

// 附近的人使用
func GetRedisSomeDynamic1Ex(pid int64, player *share_message.PlayerBase, dsList []*share_message.DynamicData, attentionList []int64) { //获取redis中的数据
	for _, log := range dsList {
		if log.GetIsShield() && pid != log.GetPlayerId() { //如果是屏蔽状态并且不是自己请求动态数据  就跳过
			continue
		}
		if util.Int64InSlice(log.GetPlayerId(), attentionList) {
			log.IsAtten = easygo.NewBool(true)
		}
		logId := log.GetLogId()
		log.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, pid))
		log.NickName = easygo.NewString(player.GetNickName())
		log.HeadIcon = easygo.NewString(player.GetHeadIcon())
		log.Sex = easygo.NewInt32(player.GetSex())
	}
}

// 完善动态信息,包含是否热门
func PerfectDynamicContainHot(pid int64, dsList []*share_message.DynamicData, attentionList []int64, hotScore int32) []*share_message.DynamicData { //获取redis中的数据
	var playerIds []int64
	lst := make([]*share_message.DynamicData, 0)
	for _, log := range dsList {
		if log.GetIsShield() && pid != log.GetPlayerId() { //如果是屏蔽状态并且不是自己请求动态数据  就跳过
			continue
		}
		playerIds = append(playerIds, log.GetPlayerId())
	}
	playerInfo := GetAllPlayerBase(playerIds, false)
	for _, log := range dsList {
		if log.GetIsShield() && pid != log.GetPlayerId() { //如果是屏蔽状态并且不是自己请求动态数据  就跳过
			continue
		}
		if util.Int64InSlice(log.GetPlayerId(), attentionList) {
			log.IsAtten = easygo.NewBool(true)
		}
		if log.GetHostScore() >= hotScore {
			log.HotType = easygo.NewInt32(1)
		}
		logId := log.GetLogId()
		p := playerInfo[log.GetPlayerId()]
		log.CommentNum = easygo.NewInt64(GetRedisDynamicCommentNum(logId))
		log.Zan = easygo.NewInt32(GetRedisDynamicZanNum(logId) + log.GetTrueZan())
		log.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, pid))
		if p != nil {
			log.NickName = easygo.NewString(p.GetNickName())
			log.HeadIcon = easygo.NewString(p.GetHeadIcon())
			log.Sex = easygo.NewInt32(p.GetSex())
			log.Types = easygo.NewInt32(p.GetTypes())
		}
		lst = append(lst, log)
		// 弥补senderType字段,此字段是后续加的
		if log.GetSenderType() == 0 {
			if p != nil {
				log.SenderType = easygo.NewInt32(p.GetTypes())
				easygo.Spawn(UpdateDynamicSenderTypeToDB, log.GetLogId(), log.GetSenderType())
			}
		}
	}

	return lst
}

/**
广场分页
in 查询出置顶的动态列表(in)
分页查询出动态数据.
遍历动态,判断是否已关注.
如果是屏蔽状态并且不是自己请求动态数据  就跳过

关注分页
找到我关注的人的id列表,找出这些人的动态列表 in
分页查询出动态数据.
*/
func GetRedisNewDynamic1(t, page, pageSize int32, pid int64, hotScore int32, advId ...int64) *share_message.DynamicDataListPage {
	player := GetRedisPlayerBase(pid)
	attentionList := player.GetAttention()
	dsList := make([]*share_message.DynamicData, 0)
	if t == SQUARE_DYNAMIC {
		attentionList = append(attentionList, pid) // 包含了自己的
	}
	// 第一页的话,需要置顶消息
	if page == 1 {
		// 是否需要关注人的 id
		atList := make([]int64, 0)
		if t == SQUARE_ATTENTION { // 需要关注人的id
			atList = append(atList, attentionList...)
		}
		//获取后台置顶的动态列表(in)
		bsTopDynamicList := GetBSTopDynamicListByIDsFromDB(pid, atList)
		if len(bsTopDynamicList) > 0 {
			slice := GetDynamicSliceByRandFromSlice(bsTopDynamicList, BS_TOP_NUM)
			//  时间最新的最靠前
			sortSlice := SortDynamicSliceByTime1(slice)
			dsList = append(dsList, sortSlice...)
		}

		//获取app置顶的动态列表(in)
		appTopDynamicList := GetAppTopDynamicListByIDsFromDB(pid, atList)
		if len(appTopDynamicList) > 0 {
			slice := GetDynamicSliceByRandFromSlice(appTopDynamicList, APP_TOP_NUM)
			//// 时间最新的最靠前
			sortSlice := SortDynamicSliceByTime1(slice)
			dsList = append(dsList, sortSlice...)
		}
	}

	var ds []*share_message.DynamicData
	var count int
	switch t {
	case SQUARE_DYNAMIC: // 广场
		// 分页查询出动态数据
		//ds, count = GetNoTopDynamicListByPageFromDB(int(page), int(pageSize))
		ds, count = getNoTopDynamicList(pid, int(page), int(pageSize), easygo.AnytoA(pid))
		dsList = append(dsList, ds...)
		dsList = ParseHotDynamic(dsList, hotScore)
	case SQUARE_ATTENTION: // 关注
		//ds, count = GetNoTopDynamicByPIDsFromDB(int(page), int(pageSize), attentionList)
		maxLogIdKey := MakeNewString(pid, "attention") // 关注.
		ds, count = GetNoTopDynamicByPIDs(pid, int(page), int(pageSize), attentionList, maxLogIdKey)

		dsList = append(dsList, ds...)
		dsList = ParseHotDynamic(dsList, hotScore)
	}
	if t == SQUARE_DYNAMIC { // 广场动态,封装广告数据
		ds = GetRedisSomeDynamic1(pid, dsList, attentionList) // 遍历动态,判断是否已关注.
		adv := append(advId, 0)[0]
		ds = GetDynamicADV(ds, adv, player.GetLastLoginIP())
		return &share_message.DynamicDataListPage{
			DynamicData: ds,
			TotalCount:  easygo.NewInt32(count),
		}
	}
	ds = GetRedisSomeDynamic1(pid, dsList, attentionList) // 遍历动态,判断是否已关注.
	return &share_message.DynamicDataListPage{
		DynamicData: ds,
		TotalCount:  easygo.NewInt32(count),
	}
}

// 封装包含话题的动态内容
func GetRedisNewDynamicTopic(page, pageSize int32, pid int64, hotScore int32, advId ...int64) *share_message.DynamicDataListPage {
	player := GetRedisPlayerBase(pid)
	attentionList := player.GetAttention()
	dsList := make([]*share_message.DynamicData, 0)

	attentionList = append(attentionList, pid) // 包含了自己的

	// 第一页的话,需要置顶消息
	if page == 1 {
		//获取后台置顶的动态列表(in)
		bsTopDynamicList := GetBSTopDynamicListByIDsFromDB(pid, []int64{})
		if len(bsTopDynamicList) > 0 {
			slice := GetDynamicSliceByRandFromSlice(bsTopDynamicList, BS_TOP_NUM)
			//  时间最新的最靠前
			sortSlice := SortDynamicSliceByTime1(slice)
			dsList = append(dsList, sortSlice...)
		}

		//获取app置顶的动态列表(in)
		appTopDynamicList := GetAppTopDynamicListByIDsFromDB(pid, []int64{})
		if len(appTopDynamicList) > 0 {
			slice := GetDynamicSliceByRandFromSlice(appTopDynamicList, APP_TOP_NUM)
			//// 时间最新的最靠前
			sortSlice := SortDynamicSliceByTime1(slice)
			dsList = append(dsList, sortSlice...)
		}
	}

	var ds []*share_message.DynamicData
	var count int
	// 分页查询出动态数据
	ds, count = getNoTopDynamicList(pid, int(page), int(pageSize), easygo.AnytoA(pid))
	dsList = append(dsList, ds...)
	dsList = ParseHotDynamic(dsList, hotScore)
	// 广场动态,封装广告数据

	ds = GetRedisSomeDynamic1(pid, dsList, attentionList) // 遍历动态,判断是否已关注.

	adv := append(advId, 0)[0]
	ds = GetDynamicADV(ds, adv, player.GetLastLoginIP())
	// 添加话题
	dynamicTopic := &share_message.DynamicData{
		DataType:  easygo.NewInt32(DATA_TYPE_TOPIC),
		TopicList: GetDynamicTopicList(),
	}
	ds = append(ds, dynamicTopic)
	return &share_message.DynamicDataListPage{
		DynamicData: ds,
		TotalCount:  easygo.NewInt32(count),
	}

}

// 设置首页的最大的动态id
func SetFirstPageMaxLogIdToRedis(key, value string) {
	err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_MAX_LOGID, key), key, value)
	if err1 != nil {
		easygo.PanicError(err1)
	}
}

// 获取首页的最大的动态id
func GetFirstPageMaxLogIdFromRedis(key string) int64 {
	b, err1 := easygo.RedisMgr.GetC().HGet(MakeRedisKey(REDIS_SQUARE_MAX_LOGID, key), key)
	if err1 != nil {
		easygo.PanicError(err1)
	}
	i, _ := strconv.Atoi(string(b))
	return int64(i)
}

func getNoTopDynamicList(opId int64, page, pageSize int, pid string) ([]*share_message.DynamicData, int) {
	var dynamicId int64
	if page > 1 { // 从redis中获取第一页的最大动态id
		dynamicId = GetFirstPageMaxLogIdFromRedis(pid)
	}
	ds, count := GetNoTopDynamicListByPageFromDB(opId, int(page), int(pageSize), dynamicId)

	if page == 1 { // 找到最大的id
		var maxId int64
		for _, v := range ds {
			if v.GetLogId() > maxId {
				maxId = v.GetLogId()
			}
		}
		// 设置id最大的值进redis
		SetFirstPageMaxLogIdToRedis(pid, easygo.AnytoA(maxId))
	}

	return ds, count
}

// 处理热门的动态
func ParseHotDynamic(ds []*share_message.DynamicData, hotScore int32) []*share_message.DynamicData {
	if len(ds) == 0 || hotScore == 0 {
		return ds
	}
	hotResult := make([]*share_message.DynamicData, 0)
	for _, dynamic := range ds {
		dynamic.HotType = easygo.NewInt32(0)
		// 判断时间是否是7天内.
		t := easygo.NowTimestamp()
		// 只有不是置顶的才给热门分
		if !dynamic.GetIsTop() && !dynamic.GetIsBsTop() {
			if dynamic.GetHostScore() >= hotScore && t-dynamic.GetSendTime() < 7*86400 { // 大于7天的不显示热门. *86400
				dynamic.HotType = easygo.NewInt32(HOT_TYPE_1)
			}
		} else {
			dynamic.HotType = easygo.NewInt32(HOT_TYPE_0) // 如果有置顶,关闭热门
		}
		hotResult = append(hotResult, dynamic)
	}
	return hotResult
}

/**
xiong
社交广场获取最新动态列表,默认50条
含了置顶操作
*/
func GetRedisNewDynamic(t int32, logId, pid int64, advId ...int64) []*share_message.DynamicData {
	ids := GetRedisSquareLogIds()
	if len(ids) == 0 {
		return nil
	}
	player := GetRedisPlayerBase(pid)
	attenList := player.GetAttention()
	newIds := make([]int64, 0)
	// todo 获取官方置顶动态
	// 获取后台置顶的动态id列表
	keys := GetTopKeysFromRedis(attenList, REDIS_SQUARE_BS_TOP_DYNAMIC, pid) // 要包含自己的置顶动态
	if len(keys) > 0 {
		slice := GetSliceByRandFromSlice(keys, BS_TOP_NUM)
		// 时间最新的最靠前
		sortSlice := SortDynamicSliceByTime(slice)
		logs.Info("返回后台置顶动态,排序后的----->", sortSlice)
		newIds = append(newIds, sortSlice...)
	}
	// 获取app置顶的keys
	appTopKeys := GetTopKeysFromRedis(attenList, REDIS_SQUARE_TOP_DYNAMIC)
	if len(appTopKeys) > 0 {
		slice := GetSliceByRandFromSlice(appTopKeys, APP_TOP_NUM)
		// 时间最新的最靠前
		sortSlice := SortDynamicSliceByTime(slice)
		newIds = append(newIds, sortSlice...)
	}
	if t == 1 { //广场
		for _, id := range ids {
			if id <= logId {
				break
			}
			if len(newIds) >= DYNAMIC_REQUEST_NUM {
				break
			}
			newIds = append(newIds, id)
		}
	} else if t == 2 { //关注
		//attenList = append(attenList, pid) 自己的
		allIds := make([]int64, 0)
		noIds := make([]int64, 0)
		Info := GetAllPlayerBase(attenList)
		for _, m := range Info {
			allIds = append(allIds, m.GetDynamicList()...)
		}
		if len(allIds) == 0 {
			return nil
		}
		easygo.SortSliceInt64(allIds, false) //降序排序
		for _, id := range allIds {
			if id <= logId {
				break
			}
			if len(newIds) >= DYNAMIC_REQUEST_NUM {
				break
			}
			if !util.Int64InSlice(id, ids) { //如果动态id没有在redis内存中
				noIds = append(noIds, id)
			}
			newIds = append(newIds, id)
		}
		if len(noIds) != 0 {
			ReloadSomeSquareInfo(noIds) //去数据库中加载
		}
	}
	if len(newIds) == 0 {
		return nil
	}
	if t == 1 { // 广场动态,封装广告数据
		lst := GetRedisSomeDynamic(pid, newIds, attenList)
		adv := append(advId, 0)[0]
		lst = GetDynamicADV(lst, adv, player.GetLastLoginIP())

		return lst
	}
	lst := GetRedisSomeDynamic(pid, newIds, attenList)
	return lst
}

// todo 优化备份使用.
//func GetRedisDynamicZanNumEx1(logIds []int64) map[int64]int32 { //获取该动态的赞的总数
//	m := make(map[int64]int32)
//	unFindList := []int64{}
//	redisCon := easygo.RedisMgr.GetC()
//	for _, id := range logIds {
//		value, err := redisCon.HLen(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, id))
//		easygo.PanicError(err)
//		if int32(value) <= 0 {
//			unFindList = append(unFindList, id)
//		} else {
//			m[id] = int32(value)
//		}
//	}
//	if len(unFindList) > 0 {
//		data, err := GetZanDataListByDynamicIdsFromDB(unFindList)
//		easygo.PanicError(err)
//		if len(data) > 0 {
//			for _, d := range data {
//				count := m[d.GetDynamicId()]
//				count += 1
//				m[d.GetDynamicId()] = count
//
//				m1 := make(map[int64]map[int64]string)
//				// 封装同一个动态的赞信息
//				for _, v := range data {
//					did := v.GetDynamicId()
//					vv, ok := m1[did]
//					if !ok {
//						vv = make(map[int64]string)
//					}
//					bytes, _ := json.Marshal(v)
//					vv[v.GetLogId()] = string(bytes)
//					m1[did] = vv
//				}
//
//				for key, value := range m1 {
//					err = redisCon.HMSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, key), value)
//					easygo.PanicError(err)
//				}
//
//			}
//
//		}
//	}
//	return m
//}

// ===========================评论相关====================================

func GetRedisPlayerMessageList(playerId int64) map[int64][]int64 { //获取所有回复我的消息logid
	info := make(map[int64][]int64)
	value, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(REDIS_SQUARE_MESSAGEIDS, playerId)))
	easygo.PanicError(err)
	for key, value := range value {
		var lst []int64
		err := json.Unmarshal([]byte(value), &lst)
		if err != nil {
			logs.Error(err)
			continue
		}
		info[key] = lst
	}
	return info
}

func AddRedisPlayerMessageId(playerId, logId, id int64) {
	b, err := easygo.RedisMgr.GetC().HExists(MakeRedisKey(REDIS_SQUARE_MESSAGEIDS, playerId), easygo.AnytoA(logId))
	easygo.PanicError(err)
	var lst []int64
	if b {
		value, err := easygo.RedisMgr.GetC().HGet(MakeRedisKey(REDIS_SQUARE_MESSAGEIDS, playerId), easygo.AnytoA(logId))
		easygo.PanicError(err)
		err1 := json.Unmarshal(value, &lst)
		easygo.PanicError(err1)
	}
	lst = append(lst, id)
	by, err2 := json.Marshal(lst)
	easygo.PanicError(err2)
	err3 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_MESSAGEIDS, playerId), easygo.AnytoA(logId), string(by))
	easygo.PanicError(err3)
}

func AddRedisPlayerMessageInfo(playerId int64, m map[int64][]int64) {
	info := make(map[int64]string) //因为redis hashmap中value必须是string类型  所以把[]int64 json成字符串存进去
	for key, lst := range m {
		b, err := json.Marshal(lst)
		if err != nil {
			logs.Error(err)
			continue
		}
		info[key] = string(b)
	}
	err := easygo.RedisMgr.GetC().HMSet(MakeRedisKey(REDIS_SQUARE_MESSAGEIDS, playerId), info)
	easygo.PanicError(err)
}

// ==========================Operate===============================

// =======================================================================
/**
xiong
根据用户id,动态id列表,获取动态列表
*/
func GetRedisSomeDynamic(pid int64, logIds []int64, attenList []int64) []*share_message.DynamicData { //获取redis中的数据
	var ids []string
	for _, id := range logIds {
		ids = append(ids, easygo.AnytoA(id))
	}
	var playerIds []int64
	lst := make([]*share_message.DynamicData, 0)

	values := make(map[int64][]byte)
	redisCon := easygo.RedisMgr.GetC()
	for _, id := range logIds {
		result, err := redisCon.HGet(MakeRedisKey(REDIS_SQUARE_DYNAMIC1, id), easygo.AnytoA(id))
		if err != nil {
			easygo.PanicError(err)
			continue
		}
		values[id] = result
	}

	//values, err := easygo.RedisMgr.GetC().HMGet(REDIS_SQUARE_DYNAMIC, ids...)

	//easygo.PanicError(err)
	for _, m := range values {
		if m == nil {
			continue
		}
		var log *share_message.DynamicData
		_ = json.Unmarshal(m, &log)
		if log.GetIsShield() && pid != log.GetPlayerId() { //如果是屏蔽状态并且不是自己请求动态数据  就跳过
			continue
		}
		if util.Int64InSlice(log.GetPlayerId(), attenList) {
			log.IsAtten = easygo.NewBool(true)
		}
		logId := log.GetLogId()
		playerIds = append(playerIds, log.GetPlayerId())
		log.CommentNum = easygo.NewInt64(GetRedisDynamicCommentNum(logId))
		log.Zan = easygo.NewInt32(GetRedisDynamicZanNum(logId) + log.GetTrueZan())
		log.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, pid))
		lst = append(lst, log)
	}

	playerInfo := GetAllPlayerBase(playerIds)
	for _, log := range lst {
		pid := log.GetPlayerId()
		player := playerInfo[pid]
		log.NickName = easygo.NewString(player.GetNickName())
		log.HeadIcon = easygo.NewString(player.GetHeadIcon())
		log.Sex = easygo.NewInt32(player.GetSex())
	}
	return lst
}

func GetRedisDynamicAllInfo(jumpMainCommentId, logId, pid int64, attenList []int64, isNew ...bool) *share_message.DynamicData {
	obj := GetRedisDynamic(logId)
	if obj == nil {
		return &share_message.DynamicData{}
	}
	playerId := obj.GetPlayerId()
	player := GetRedisPlayerBase(playerId)
	if util.Int64InSlice(playerId, attenList) {
		obj.IsAtten = easygo.NewBool(true)
	}
	obj.CommentNum = easygo.NewInt64(GetRedisDynamicCommentNum(logId))
	if isNew == nil {
		obj.CommentList = GetRedisDynamicCommentInfo(logId, 0)
	}
	// 判断是否有跳转的评论id,如果有,插入进第一条.
	if jumpMainCommentId > 0 {
		cd := GetCommentDataFromDB(jumpMainCommentId)
		if cd != nil { // 重新复制
			cl := new(share_message.CommentList)
			cl.HotList = obj.GetCommentList().GetHotList()
			cl.CommentInfo = append(cl.CommentInfo, cd)
			for _, v := range obj.GetCommentList().GetCommentInfo() {
				cl.CommentInfo = append(cl.CommentInfo, v)
			}
			obj.CommentList = cl
		}
	}
	obj.Zan = easygo.NewInt32(GetRedisDynamicZanNum(logId) + obj.GetTrueZan())
	obj.IsZan = easygo.NewBool(GetRedisDynamicIsZan(logId, pid))
	obj.NickName = easygo.NewString(player.GetNickName())
	obj.HeadIcon = easygo.NewString(player.GetHeadIcon())
	obj.Sex = easygo.NewInt32(player.GetSex())
	obj.TrueZan = easygo.NewInt32(obj.GetTrueZan())
	obj.HostScore = easygo.NewInt32(obj.GetHostScore())
	return obj
}

func GetRedisMessageAllInfo(playerId, Id int64) *client_square.MessageMainInfo { //获取消息界面的信息
	msg := &client_square.MessageMainInfo{
		UnreadComment:   easygo.NewInt32(GetPlayerUnreadInfo(playerId, UNREAD_COMMENT)),
		UnreadZan:       easygo.NewInt32(GetPlayerUnreadInfo(playerId, UNREAD_ZAN)),
		UnreadAttention: easygo.NewInt32(GetPlayerUnreadInfo(playerId, UNREAD_ATTENTION)),
	}
	info := GetRedisPlayerMessageList(playerId) //获取人物身上所有评论id  按评论id排序 倒序排序
	if len(info) == 0 {
		return msg
	}

	messageInfo := make(map[int64]int64)
	for logId, ids := range info { // key--> 动态的id,value--> 评论的内容的id数组
		for _, id := range ids {
			messageInfo[id] = logId
		}
	}
	var messageIds []int64
	for id := range messageInfo {
		if id >= Id && Id != 0 {
			continue
		}
		messageIds = append(messageIds, id)
	}
	easygo.SortSliceInt64(messageIds, false)

	if len(messageIds) >= DYNAMIC_REQUEST_NUM {
		messageIds = messageIds[:DYNAMIC_REQUEST_NUM]
	}

	var commentList []*share_message.CommentData
	var dynamicList []*share_message.DynamicData
	var dynamicIds []int64
	for _, id := range messageIds { //获取回复信息
		logId := messageInfo[id]
		if logId == 0 {
			logs.Error("动态id怎么会是0,id:", id)
			continue
		}
		b1, err := easygo.RedisMgr.GetC().HExists(MakeRedisKey(REDIS_SQUARE_COMMENT, logId), easygo.AnytoA(id))
		if err != nil {
			logs.Error(err)
			continue
		}
		if b1 {
			if !util.Int64InSlice(logId, dynamicIds) {
				dynamicIds = append(dynamicIds, logId)
				dynamic := GetRedisDynamic(logId)
				if dynamic == nil {
					continue
				}
				player := GetRedisPlayerBase(dynamic.GetPlayerId())
				name := player.GetNickName()
				headIcon := player.GetHeadIcon()
				sex := player.GetSex()
				dynamic.NickName = easygo.NewString(name)
				dynamic.HeadIcon = easygo.NewString(headIcon)
				dynamic.Sex = easygo.NewInt32(sex)
				dynamicList = append(dynamicList, dynamic)
			}
			commentInfo := GetOneCommentInfo(logId, id)
			if playerId != commentInfo.GetPlayerId() { // 如果是自己评论的,就不添加进去
				commentList = append(commentList, commentInfo)
			}
		} else {
			logs.Error("动态id怎么会没有消息id:", id)
		}
	}

	msg.DynamicData = dynamicList
	msg.CommentData = commentList
	return msg
}

func AddPlayerUnreadInfo(pid int64, t int32, value int64) {
	b, err := easygo.RedisMgr.GetC().HExists(MakeRedisKey(REDIS_SQUARE_UNREADINFO, pid), easygo.AnytoA(t))
	easygo.PanicError(err)
	if !b { //如果不存在
		if value > 0 {
			easygo.RedisMgr.GetC().HIncrBy(MakeRedisKey(REDIS_SQUARE_UNREADINFO, pid), easygo.AnytoA(t), value)
		} else {
			_ = easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_UNREADINFO, pid), easygo.AnytoA(t), 0)
		}
	} else {
		by, err := easygo.RedisMgr.GetC().HGet(MakeRedisKey(REDIS_SQUARE_UNREADINFO, pid), easygo.AnytoA(t))
		easygo.PanicError(err)
		num := easygo.AtoInt64(string(by))
		if num > 0 {
			easygo.RedisMgr.GetC().HIncrBy(MakeRedisKey(REDIS_SQUARE_UNREADINFO, pid), easygo.AnytoA(t), value)
		} else if num <= 0 && value > 0 {
			easygo.RedisMgr.GetC().HIncrBy(MakeRedisKey(REDIS_SQUARE_UNREADINFO, pid), easygo.AnytoA(t), value-num)
		}
	}

}

func GetPlayerUnreadInfo(pid int64, t int32) int32 {
	b, err := easygo.RedisMgr.GetC().HExists(MakeRedisKey(REDIS_SQUARE_UNREADINFO, pid), easygo.AnytoA(t))
	easygo.PanicError(err)
	if !b {
		return 0
	}

	value, err1 := easygo.RedisMgr.GetC().HGet(MakeRedisKey(REDIS_SQUARE_UNREADINFO, pid), easygo.AnytoA(t))
	easygo.PanicError(err1)
	return easygo.AtoInt32(string(value))
}

func DelPlayerUnreadInfo(pid int64, t int32) {
	_, err := easygo.RedisMgr.GetC().Hdel(MakeRedisKey(REDIS_SQUARE_UNREADINFO, pid), easygo.AnytoA(t))
	easygo.PanicError(err)
}

func GetIsNewMessage(playerId int64) bool {
	comment := GetPlayerUnreadInfo(playerId, UNREAD_COMMENT)
	zan := GetPlayerUnreadInfo(playerId, UNREAD_ZAN)
	att := GetPlayerUnreadInfo(playerId, UNREAD_ATTENTION)
	if comment == 0 && zan == 0 && att == 0 {
		return false
	}
	return true
}

func GetNewUnReadMessageFromRedis(playerId int64) *client_square.NewUnReadMessageResp {
	return &client_square.NewUnReadMessageResp{
		UnreadComment:   easygo.NewInt32(GetPlayerUnreadInfo(playerId, UNREAD_COMMENT)),
		UnreadZan:       easygo.NewInt32(GetPlayerUnreadInfo(playerId, UNREAD_ZAN)),
		UnreadAttention: easygo.NewInt32(GetPlayerUnreadInfo(playerId, UNREAD_ATTENTION)),
	}
}

//==============================savedatabase==============================

func UpdateSquareDynamic() {

	logIds := GetRedisSquareLogIds()
	if len(logIds) > 0 {
		for _, logId := range logIds {
			maxCommentId := GetDynamicSaveCommentId()
			commentList := GetRedisDynamicAllComment(logId, maxCommentId) //存储评论数据
			if len(commentList) > 0 {
				var data []interface{}
				var maxId int64
				for _, v := range commentList {
					id := v.GetId()
					if id > maxId {
						maxId = id
					}
					b1 := bson.M{"_id": id}
					b2 := v
					data = append(data, b1, b2)
				}
				UpdateDynamicSaveCommentId(maxId)
				easygo.Spawn(func() { UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_SQUARE_COMMENT, data) })
			}

			maxZanId := GetDynamicSaveZanId()
			zanlst := GetRedisDynamicAllZanInfo(logId, maxZanId) //存储赞的数据
			if len(zanlst) > 0 {
				var data1 []interface{}
				var maxId int64
				for _, v := range zanlst {
					id := v.GetLogId()
					if id > maxId {
						maxId = id
					}
					b1 := bson.M{"_id": id}
					b2 := v
					data1 = append(data1, b1, b2)
				}
				UpdateDynamicSaveZanId(maxId)
				easygo.Spawn(func() { UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_SQUARE_ZAN, data1) })
			}
			delzanlst := GetRedisDynamicDelZanId()
			if len(delzanlst) > 0 {
				col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_ZAN)
				defer closeFunc()
				_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": delzanlst}})
				if err != nil {
					logs.Error("delzanlst err", err)
					continue
				}
				ClearRedisDynamicDelZanId()
			}
		}
	}

	keys, err := easygo.RedisMgr.GetC().Scan(REDIS_SQUARE_ATTENTION) //存储所有关注的信息
	easygo.PanicError(err)
	var lst []*share_message.AttentionData
	maxAttId := GetDynamicSaveAttId()
	for _, key := range keys {
		value, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(key))
		easygo.PanicError(err)
		for _, s := range value {
			var msg *share_message.AttentionData
			err := json.Unmarshal([]byte(s), &msg)
			if err != nil {
				logs.Error(err)
				continue
			}
			id := msg.GetLogId()
			if id > maxAttId {
				lst = append(lst, msg)
			}
		}
	}
	if len(lst) > 0 {
		var data []interface{}
		var maxId int64
		for _, v := range lst {
			id := v.GetLogId()
			if id > maxId {
				maxId = id
			}
			b1 := bson.M{"_id": id}
			b2 := v
			data = append(data, b1, b2)
		}
		UpdateDynamicSaveAttId(maxId)
		easygo.Spawn(func() { UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_SQUARE_ATTENTION, data) })
	}
	delattlst := GetRedisDynamicDelAttId()
	if len(delattlst) > 0 {
		col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_ATTENTION)
		defer closeFunc()
		_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": delattlst}})
		easygo.PanicError(err)
		ClearRedisDynamicDelAttId()
	}
}

func AddSaveDynamicIdList(logId int64) {
	err := easygo.RedisMgr.GetC().SAdd(REDIS_SQUARE_SAVELIST, logId)
	easygo.PanicError(err)
}

func GetSaveDynamicIdList() []int64 {
	ids := []int64{}
	value, err := easygo.RedisMgr.GetC().Smembers(REDIS_SQUARE_SAVELIST)
	easygo.PanicError(err)
	InterfersToInt64s(value, &ids)
	return ids
}

func DelSaveDynamicIdList() {
	_, err := easygo.RedisMgr.GetC().Delete(REDIS_SQUARE_SAVELIST)
	easygo.PanicError(err)
}

const SaveDynamicCnt = 300

func DealSquareDynamicKeys() { //清空社交广场长时间不用的动态
	defer easygo.AfterFunc(1*time.Hour, DealSquareDynamicKeys)
	logIds := GetRedisSquareLogIds()
	if len(logIds) > SaveDynamicCnt {
		easygo.SortSliceInt64(logIds, false) //倒序
		lst := logIds[SaveDynamicCnt-1:]     //前300条永远不清
		saveList := GetSaveDynamicIdList()
		easygo.SortSliceInt64(saveList, false) //倒序
		var newlst []int64                     //需要清除的ids
		for _, logId := range lst {
			if util.Int64InSlice(logId, saveList) {
				continue
			}
			newlst = append(newlst, logId)
		}
		if len(newlst) > 0 {
			for _, logId := range newlst {
				DelSquareDynamicById(logId)
				DelDynamicComment(logId)
			}
		}
		DelSaveDynamicIdList()
	}

}

func GetRedisPlayerSomeDynamic(playerId, logId int64, t int32) *client_hall.DynamicInfo {
	player := GetRedisPlayerBase(playerId)
	if player == nil {
		panic("玩家对象怎么会为空")
	}
	logList := player.GetDynamicList()
	easygo.SortSliceInt64(logList, false)
	var ids []int64
	var num int
	if t == 1 { //最新
		for _, id := range logList {
			if id > logId {
				ids = append(ids, id)
				num += 1
			}
			if num >= DYNAMIC_REQUEST_NUM || id <= logId {
				break
			}
		}
	} else { //下一页
		for _, id := range logList {
			if id < logId {
				ids = append(ids, id)
				num += 1
			}
			if num >= DYNAMIC_REQUEST_NUM {
				break
			}
		}
	}
	lst := GetRedisDynamicForSomeLogId(false, playerId, ids)
	msg := &client_hall.DynamicInfo{DynamicData: lst}
	return msg
}

// 判断前端是否有提交过动态,并返回动态内容
func CheckDuplicateDynamicData(uniqueCode string) (*share_message.DynamicData, error) {
	var data *share_message.DynamicData
	b, err := CheckHasInRedis(MakeRedisKey(REDIS_SQUARE_UNIQUE_CODE, uniqueCode))
	if err != nil {
		return data, err
	}

	if !b { // 不存在,从数据库中读取,然后放进redis
		col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
		defer closeFunc()
		err := col.Find(bson.M{"ClientUniqueCode": uniqueCode, "Statue": DYNAMIC_STATUE_COMMON}).One(&data)
		if err != nil {
			return data, err
		}
		bytes, _ := json.Marshal(data)
		// 存放进redis
		err = easygo.RedisMgr.GetC().Set(MakeRedisKey(REDIS_SQUARE_UNIQUE_CODE, uniqueCode), string(bytes))
		if err != nil {
			return data, err
		}
		err = easygo.RedisMgr.GetC().Expire(MakeRedisKey(REDIS_SQUARE_UNIQUE_CODE, uniqueCode), 10) // 10秒过期
		if err != nil {
			return data, err
		}
		return data, nil
	}

	s, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(REDIS_SQUARE_UNIQUE_CODE, uniqueCode))
	if err != nil {
		return data, err
	}
	_ = json.Unmarshal([]byte(s), &data)
	return data, nil
}

// 存重复动态标记
func SaveDuplicateDynamicDataToRedis(uniqueCode string, data *share_message.DynamicData) {
	bytes, _ := json.Marshal(data)
	err := easygo.RedisMgr.GetC().Set(MakeRedisKey(REDIS_SQUARE_UNIQUE_CODE, uniqueCode), string(bytes))
	easygo.PanicError(err)
	err = easygo.RedisMgr.GetC().Expire(MakeRedisKey(REDIS_SQUARE_UNIQUE_CODE, uniqueCode), 600) // 10十分钟过期
	easygo.PanicError(err)
}

// 修改置顶状态(),置顶/取消置顶
func UpdateDynamicInfo(dynamic *share_message.DynamicData) {
	// 修改 hset中的动态信息.
	UpdateRedisSquareDynamic(dynamic)
	// 修改数据库的动态信息
	SaveSquareDynamic(dynamic.GetLogId())
}

// t秒后执行取消置顶.
func AfterCancelTop(t time.Duration, dynamicId int64) *easygo.Timer {
	return easygo.AfterFunc(t, func() {
		dynamic := GetDynamicByStatusSFromDB(dynamicId, []int{DYNAMIC_STATUE_COMMON, DYNAMIC_STATUE_UNPUBLISHED})

		dynamic.IsBsTop = easygo.NewBool(false)

		dynamic.IsTop = easygo.NewBool(false)
		dynamic.TopOverTime = easygo.NewInt64(0)

		// 修改动态置顶状态
		UpdateDynamicInfo(dynamic)
		// 从map删除定时任务管理器
		TimerMgr.TimerMap.Delete(dynamicId)
		// 从数据库中删除管理器记录
		DelSquareTopTimeMgrByIdFromDB(dynamicId)
	})
}

// 设置置顶进redis 根据传进来的key确定是后台置顶还是app置顶
func SetTopDynamicToRedis(redisName string, data *share_message.DynamicData) {
	b, err := json.Marshal(data)
	easygo.PanicError(err)
	err1 := easygo.RedisMgr.GetC().HSet(redisName, easygo.AnytoA(data.GetLogId()), string(b))
	easygo.PanicError(err1)
}

// 获取置顶的logId列表
func GetTopKeysFromRedis(attentionList []int64, name string, pid ...int64) []int64 {
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(name))
	easygo.PanicError(err)
	ids := make([]int64, 0)
	// 获取redis中置顶的动态,然后判断是否是自己关注的人的动态,取出动态id列表
	if len(pid) > 0 {
		attentionList = append(attentionList, pid...)
	}
	for _, v := range attentionList {
		for _, m := range values {
			dynamic := &share_message.DynamicData{}
			_ = json.Unmarshal([]byte(m), &dynamic)
			if dynamic.GetPlayerId() == 0 {
				continue
			}
			if dynamic.GetPlayerId() == v {
				ids = append(ids, dynamic.GetLogId())
			}
		}
	}
	return ids
}

// 我的主页获取置顶的logId列表
func GetMyMainTopKeysFromRedis(pid int64, name string) []int64 {
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(name))
	easygo.PanicError(err)
	ids := make([]int64, 0)
	// 获取redis中置顶的动态,然后判断是否是自己关注的人的动态,取出动态id列表
	for _, m := range values {
		dynamic := &share_message.DynamicData{}
		json.Unmarshal([]byte(m), &dynamic)
		if dynamic.GetPlayerId() == 0 {
			continue
		}
		if dynamic.GetPlayerId() == pid {
			ids = append(ids, dynamic.GetLogId())
		}
	}
	return ids
}

// 删除置顶列表中的数据
func DelTopDynamicListFromRedis(redisName string, keys ...string) {
	_, err := easygo.RedisMgr.GetC().Hdel(redisName, keys...)
	easygo.PanicError(err)
}

func GetADVSetMap(location int32, ips ...string) (map[int64]*share_message.AdvSetting, map[int64]int, []*share_message.AdvSetting) {
	ip := ""
	if len(ips) > 0 {
		ip = ips[0]
	}
	// 获取广告,设置进map
	advList := QueryAdvListToDB(location, ip)
	m := make(map[int64]*share_message.AdvSetting)
	idMap := make(map[int64]int) // index 和广告id对应的map   map[广告列表第一条广告的id]第几条广告
	if len(advList) > 0 {
		for k, v := range advList {
			m[v.GetId()] = v
			idMap[v.GetId()] = k
		}
		return m, idMap, advList
	}
	return m, idMap, advList
}

// 获取动态内容(包含广告)
func GetDynamicADV(dns []*share_message.DynamicData, advId int64, ips ...string) []*share_message.DynamicData {
	if len(dns) == 0 {
		return dns
	}
	ip := ""
	if len(ips) > 0 {
		ip = ips[0]
	}
	// 查询广告数据进map
	m, advIdMap, advList := GetADVSetMap(ADV_LOCATION_SQUARE, ip)
	if len(advList) == 0 || len(m) == 0 {
		return dns
	}
	var nextAdvId int64
	// 如果advId为0,使用第一条广告
	if advId == 0 {
		nextAdvId = advList[0].GetId()
	} else { // 如果advId不为0 ,使用advId+1,判断是否超了
		// 取出广告id对应的位置
		index, b := advIdMap[advId]
		if !b {
			return dns
		}
		// 列表中取出广告
		if len(advList)-1 < index+1 { // 没有广告可以取了
			return dns
		}
		adv := advList[index+1]
		nextAdvId = adv.GetId()
	}
	if nextAdvId == 0 { // 没有广告
		return dns
	}
	var insertIndex int64
	for in := 0; in < len(dns); in++ {
		rand := RangeRand(ADV_MIN_INDEX, ADV_MAX_INDEX) // 得出广告的插入位置
		insertIndex += int64(in) + rand
		if int(insertIndex) > len(dns) {
			logs.Error("============> 超过最大索引了")
			break
		}
		// 获得广告
		adv, b := m[nextAdvId]
		if !b {
			return dns
		}

		dn := &share_message.DynamicData{}
		adv.CreateTime = easygo.NewInt64(dns[insertIndex-1].GetSendTime()) // 取出上一条动态的值,获取时间,封装进广告的createTime
		dn.AdvSetting = adv
		dn.DataType = easygo.NewInt32(DATA_TYPE_ADV)
		dns = easygo.Insert(dns, int(insertIndex), dn).([]*share_message.DynamicData)
		//  下一条广告.
		index, b := advIdMap[adv.GetId()]
		if !b {
			return dns
		}
		// 列表中取出广告
		if len(advList)-1 < index+1 { // 没有广告可以取了
			break
		}
		nextAdv := advList[index+1]
		nextAdvId = nextAdv.GetId()
	}
	return dns
}

func GetDynamicSliceByRandFromSlice(ds []*share_message.DynamicData, count int) []*share_message.DynamicData {
	logs.Info("所有后台置顶条数为=====>", len(ds))
	result := make([]*share_message.DynamicData, 0)
	if len(ds) == 0 {
		return result
	}
	if len(ds) <= count {
		return ds
	}
	m := make(map[int64]*share_message.DynamicData)
	ids := make([]int64, 0)
	tempIds := make([]int64, 0)
	for i := 0; i < len(ds); i++ {
		tempIds = append(tempIds, ds[i].GetLogId())
		m[ds[i].GetLogId()] = ds[i]
	}
	for i := 0; i < count; i++ {
		id := RandInt(0, len(tempIds))
		ids = append(ids, tempIds[id])
		tempIds = easygo.Del(tempIds, tempIds[id]).([]int64)
	}
	for _, i := range ids {
		result = append(result, m[i])
	}
	return result
}

//==================================
// 设置到redis
func SetSquareDynamicToRedis(obj *share_message.DynamicData) {
	b, err := json.Marshal(obj)
	easygo.PanicError(err)
	err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_DYNAMIC1, obj.GetLogId()), easygo.AnytoA(obj.GetLogId()), string(b))
	easygo.PanicError(err1)
}

// 检验是否存在redis中
func CheckHasInRedis(name string) (bool, error) {
	b, err := easygo.RedisMgr.GetC().Exist(name)
	easygo.PanicError(err)
	return b, err
}

//设置用户每日发布获取的贡献度
func SetDevoteDynamic(playerId, topicId int64, topicName string) int64 {
	data := time.Now().Format("2006-01-02")
	key := MakeRedisKey(TOPIC_DEVOTE, data, topicId, playerId)
	hv, _ := easygo.RedisMgr.GetC().HGet(key, "dynamic")
	iv := easygo.AtoInt64(string(hv))
	if iv >= DynamicMax {
		return iv
	}
	v := easygo.RedisMgr.GetC().HIncrBy(key, "dynamic", DynamicInc)
	easygo.RedisMgr.GetC().Expire(key, SurplusTime())
	if v <= DynamicMax {
		SetTopicPlayerDevoteDay(playerId, topicId, DynamicInc, topicName)
		SetTopicPlayerDevoteMonth(playerId, topicId, DynamicInc, topicName)
		SetTopicPlayerDevoteTotal(playerId, topicId, DynamicInc, topicName)
	}
	return v
}

//设置用户每日评论获取的贡献度
func SetDevoteComment(playerId, topicId int64, topicName string) int64 {
	data := time.Now().Format("2006-01-02")
	key := MakeRedisKey(TOPIC_DEVOTE, data, topicId, playerId)
	hv, _ := easygo.RedisMgr.GetC().HGet(key, "comment")
	iv := easygo.AtoInt64(string(hv))
	if iv >= CommentMax {
		return iv
	}
	v := easygo.RedisMgr.GetC().HIncrBy(key, "comment", CommentInc)
	easygo.RedisMgr.GetC().Expire(key, SurplusTime())
	if v <= CommentMax {
		SetTopicPlayerDevoteDay(playerId, topicId, CommentInc, topicName)
		SetTopicPlayerDevoteMonth(playerId, topicId, CommentInc, topicName)
		SetTopicPlayerDevoteTotal(playerId, topicId, CommentInc, topicName)
	}
	return v
}

//设置用户每日点赞获取的贡献度
func SetDevoteLike(playerId, topicId int64, topicName string) int64 {
	data := time.Now().Format("2006-01-02")
	key := MakeRedisKey(TOPIC_DEVOTE, data, topicId, playerId)
	hv, _ := easygo.RedisMgr.GetC().HGet(key, "like")
	iv := easygo.AtoInt64(string(hv))
	if iv >= LikeMax {
		return iv
	}
	v := easygo.RedisMgr.GetC().HIncrBy(key, "like", LikeInc)
	easygo.RedisMgr.GetC().Expire(key, SurplusTime())
	if v <= LikeMax {
		SetTopicPlayerDevoteDay(playerId, topicId, LikeInc, topicName)
		SetTopicPlayerDevoteMonth(playerId, topicId, LikeInc, topicName)
		SetTopicPlayerDevoteTotal(playerId, topicId, LikeInc, topicName)
	}
	return v
}

//减去用户每日发布动态获取的贡献度
func LessDevoteDynamic(playerId, topicId, logId int64) int64 {
	data := time.Now().Format("2006-01-02")
	dynamicIncV := int64(-DynamicInc)

	//删除旧动态，不需要redis操作，直接更新数据库数据
	dynamicInfo, _ := GetDynamicByLogIdFromDB(logId)
	dynamicCreateData := time.Unix(dynamicInfo.GetCreateTime(), 0).Format("2006-01-02")
	if dynamicCreateData != data {
		year := easygo.AtoInt32(time.Unix(dynamicInfo.GetCreateTime(), 0).Format("2006"))
		month := easygo.AtoInt32(time.Unix(dynamicInfo.GetCreateTime(), 0).Format("01"))
		day := easygo.AtoInt32(time.Unix(dynamicInfo.GetCreateTime(), 0).Format("02"))
		SetSpecifyTopicPlayerDevoteDay(playerId, topicId, dynamicIncV, "", year, month, day)
		SetSpecifyTopicPlayerDevoteMonth(playerId, topicId, dynamicIncV, "", year, month)
		SetTopicPlayerDevoteTotal(playerId, topicId, dynamicIncV, "")
		return 0
	}

	key := MakeRedisKey(TOPIC_DEVOTE, data, topicId, playerId)
	hv, _ := easygo.RedisMgr.GetC().HGet(key, "dynamic")
	iv := easygo.AtoInt64(string(hv))
	if iv <= 0 {
		return iv
	}
	v := easygo.RedisMgr.GetC().HIncrBy(key, "dynamic", dynamicIncV)
	easygo.RedisMgr.GetC().Expire(key, SurplusTime())
	if v <= DynamicMax {
		SetTopicPlayerDevoteDay(playerId, topicId, dynamicIncV, "")
		SetTopicPlayerDevoteMonth(playerId, topicId, dynamicIncV, "")
		SetTopicPlayerDevoteTotal(playerId, topicId, dynamicIncV, "")
	}
	return v
}

//减去用户每日评论获取的贡献度
func LessDevoteComment(playerId, topicId int64) int64 {
	data := time.Now().Format("2006-01-02")
	key := MakeRedisKey(TOPIC_DEVOTE, data, topicId, playerId)

	hv, _ := easygo.RedisMgr.GetC().HGet(key, "comment")
	iv := easygo.AtoInt64(string(hv))
	if iv <= 0 {
		return iv
	}
	commentIncV := int64(-CommentInc)
	v := easygo.RedisMgr.GetC().HIncrBy(key, "comment", commentIncV)
	easygo.RedisMgr.GetC().Expire(key, SurplusTime())
	if v <= CommentMax {
		SetTopicPlayerDevoteDay(playerId, topicId, commentIncV, "")
		SetTopicPlayerDevoteMonth(playerId, topicId, commentIncV, "")
		SetTopicPlayerDevoteTotal(playerId, topicId, commentIncV, "")
	}
	return v
}

//减去用户每日点赞获取的贡献度
func LessDevoteLike(playerId, topicId int64) int64 {
	data := time.Now().Format("2006-01-02")
	key := MakeRedisKey(TOPIC_DEVOTE, data, topicId, playerId)

	hv, _ := easygo.RedisMgr.GetC().HGet(key, "like")
	iv := easygo.AtoInt64(string(hv))
	if iv <= 0 {
		return iv
	}
	likeIncV := int64(-LikeInc)
	v := easygo.RedisMgr.GetC().HIncrBy(key, "like", likeIncV)
	easygo.RedisMgr.GetC().Expire(key, SurplusTime())
	if v <= LikeMax {
		SetTopicPlayerDevoteDay(playerId, topicId, likeIncV, "")
		SetTopicPlayerDevoteMonth(playerId, topicId, likeIncV, "")
		SetTopicPlayerDevoteTotal(playerId, topicId, likeIncV, "")
	}
	return v
}
