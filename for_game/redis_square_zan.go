package for_game

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/client_square"
	"game_server/pb/share_message"
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
	DYNAMIC_OPERATE    = 1 //点赞或者关注
	DYNAMIC_DELOPERATE = 2 //取消点赞或者取消关注
	UNREAD_ZAN         = 2
)

const (
	//REDIS_SQUARE_DYNAMICZAN      = "redis_square_dynamiczan"      //动态赞信息 map[redis_square_dynamiczan_logid]map[id]string  // 原来的,暂时不删
	REDIS_SQUARE_DYNAMICZAN = "redis_square:zan" //动态赞信息 map[redis_square_dynamiczan_logid]map[id]string
	//REDIS_SQUARE_ATTENTION       = "redis_square_attention"       //玩家被关注信息 map[redis_square_attention_pid]map[id]string
	REDIS_SQUARE_DYNAMICZANINFO = "redis_square:dynamiczaninfo" //被赞的动态信息 map[redis_square_dynamiczaninfo_pid]map[logId][]id
	REDIS_SQUARE_DELZANIDS      = "redis_square:delzanids"      //取消点赞的id

	//存储限制key
	REDIS_SQUARE_SAVEZANID = "savezanid" //上次存库的最后赞id
)

/**
xiong
社交广场登录时
加载自己赞过的动态id
*/
func ReloadPlayerZanIds(playerId int64, ch chan int) {
	defer SquareRecoverAndLog(ch)
	lst, err := GetZanDataListByOperateIdFromDB(playerId)
	easygo.PanicError(err)
	if len(lst) == 0 {
		return
	}
	var ids []int64
	for _, msg := range lst {
		ids = append(ids, msg.GetDynamicId())
	}
	err1 := easygo.RedisMgr.GetC().SAdd(MakeRedisKey(REDIS_SQUARE_MYZANIDS, playerId), ids) //这里包括被删除的动态  但是没关系  用不到
	easygo.PanicError(err1)
}

/**
xiong
登录社交广场时
加载指定动态被赞的信息进redis hset
*/
func ReloadDynamicZan(logIds []int64, ch chan int) { //加载logIds里的动态 被赞的信息
	defer SquareRecoverAndLog(ch)
	lst, err := GetZanDataListByDynamicIdsFromDB(logIds)
	easygo.PanicError(err)
	if len(lst) == 0 {
		return
	}
	info := make(map[int64]map[int64][]int64)
	var maxId int64
	for _, msg := range lst {
		playerId := msg.GetPlayerId()
		id := msg.GetLogId()
		logId := msg.GetDynamicId()
		info1 := info[playerId]
		if info1 == nil {
			info1 = make(map[int64][]int64)
			info1[logId] = []int64{id}
			info[playerId] = info1
		} else {
			info1[logId] = append(info1[logId], id)
		}
		if maxId < id {
			maxId = id
		}
		b, err := json.Marshal(msg)
		if err != nil {
			logs.Error(err)
			continue
		}
		err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, logId), easygo.AnytoA(id), string(b))
		if err1 != nil {
			logs.Error(err1)
		}
	}

	for pid, m := range info {
		err := AddRedisPlayerBeZanIdInfo(pid, m)
		if err != nil {
			logs.Error(err)
		}
	}

	if maxId > GetDynamicSaveZanId() {
		UpdateDynamicSaveZanId(maxId)
	}
}

/**
xiong
从redis 中删除赞信息
*/
func DelDynamicZan(logId int64) {
	redisKey := MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, easygo.AnytoA(logId))
	b, err := easygo.RedisMgr.GetC().Exist(redisKey)
	easygo.PanicError(err)
	if !b {
		return
	}
	b1, err1 := easygo.RedisMgr.GetC().Delete(redisKey) //从redis中删除动态赞的信息
	easygo.PanicError(err1)
	if !b1 {
		logs.Error("redis删除动态赞失败id:", logId)
	}
}

/**
xiong
删除指定用户对某条动态的赞信息
*/
func DelDynamicZanInfo(pid, logId int64) {
	redisKey := MakeRedisKey(REDIS_SQUARE_DYNAMICZANINFO, easygo.AnytoA(pid))
	b, err := easygo.RedisMgr.GetC().HExists(redisKey, easygo.AnytoA(logId))
	easygo.PanicError(err)
	if !b {
		return
	}
	b1, err6 := easygo.RedisMgr.GetC().Hdel(redisKey, easygo.AnytoA(logId)) //从redis中删除动态赞的信息
	easygo.PanicError(err6)
	if !b1 {
		logs.Error("redis删除玩家被赞动态失败id:", logId)
	}

}

func GetRedisDynamicZanNum(logId int64) int32 { //获取该动态的赞的总数
	redisCon := easygo.RedisMgr.GetC()
	value, err := redisCon.HLen(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, logId))
	easygo.PanicError(err)
	count := int32(value)
	if int32(value) <= 0 { // 去数据库查找
		data, err := GetZanDataListByDynamicIdsFromDB([]int64{logId})
		easygo.PanicError(err)
		if len(data) > 0 {
			count = int32(len(data))
			m := make(map[int64]string)
			for _, v := range data {
				bytes, _ := json.Marshal(v)
				m[v.GetLogId()] = string(bytes)
			}
			err = redisCon.HMSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, logId), m)
			easygo.PanicError(err)
		}
	}
	return count
}

func GetRedisDynamicZanNumEx(logIds []int64) map[int64]int32 { //获取该动态的赞的总数
	m := make(map[int64]int32)
	unFindList := []int64{}
	redisCon := easygo.RedisMgr.GetC()
	for _, id := range logIds {
		value, err := redisCon.HLen(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, id))
		easygo.PanicError(err)
		if int32(value) <= 0 {
			unFindList = append(unFindList, id)
		} else {
			m[id] = int32(value)
		}
	}
	if len(unFindList) > 0 {
		data, err := GetZanDataListByDynamicIdsFromDB(unFindList)
		easygo.PanicError(err)
		if len(data) > 0 {
			for _, d := range data {
				count := m[d.GetDynamicId()]
				count += 1
				m[d.GetDynamicId()] = count
				/*
					m1 := make(map[int64]string)
					for _, v := range data {
						bytes, _ := json.Marshal(v)
						m1[v.GetLogId()] = string(bytes)
					}
					err = redisCon.HMSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, d.GetDynamicId()), m1)
					easygo.PanicError(err)*/

			}
			fun := func(data []*share_message.ZanData) {
				for _, d := range data {
					m1 := make(map[int64]string)
					for _, v := range data {
						bytes, _ := json.Marshal(v)
						m1[v.GetLogId()] = string(bytes)
					}
					err = redisCon.HMSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, d.GetDynamicId()), m1)
					easygo.PanicError(err)
				}

			}
			easygo.Spawn(fun, data)
		}
	}
	return m
}

func GetRedisDynamicZanNumEx1(logIds []int64) map[int64]int32 { //获取该动态的赞的总数
	m := make(map[int64]int32)
	unFindList := make([]int64, 0)
	redisCon := easygo.RedisMgr.GetC()
	for _, id := range logIds {
		value, err := redisCon.HLen(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, id))
		easygo.PanicError(err)
		if int32(value) <= 0 {
			unFindList = append(unFindList, id)
		} else {
			m[id] = int32(value)
		}
	}
	if len(unFindList) > 0 {
		data, err := GetZanDataListByDynamicIdsFromDB(unFindList)
		easygo.PanicError(err)
		if len(data) > 0 {
			for _, d := range data {
				count := m[d.GetDynamicId()]
				count += 1
				m[d.GetDynamicId()] = count
			}
			// 异步存redis
			fun := func(data []*share_message.ZanData) {
				m1 := make(map[int64]map[int64]string)
				// 封装同一个动态的赞信息
				for _, v := range data {
					did := v.GetDynamicId()
					vv, ok := m1[did]
					if !ok {
						vv = make(map[int64]string)
					}
					bytes, _ := json.Marshal(v)
					vv[v.GetLogId()] = string(bytes)
					m1[did] = vv
				}
				for key, value := range m1 {
					err = redisCon.HMSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, key), value)
					easygo.PanicError(err)
				}
			}
			easygo.Spawn(fun, data)
		}
	}
	return m
}

// ===========================评论相关====================================

// ==========================Operate===============================

func GetRedisDynamicZan(logId int64) map[int64]string { //获取所有该动态的点赞信息
	// 先从数据库查询出来,然后在set进去,避免redis中没有,数据库中有.
	zanDataLis := GetZanDataListByDynamicId(logId)
	m := make(map[int64]string)
	for _, value := range zanDataLis {
		zanBytes, _ := json.Marshal(value)
		m[value.GetLogId()] = string(zanBytes)
	}
	if len(m) > 0 {
		err1 := easygo.RedisMgr.GetC().HMSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, logId), m)
		easygo.PanicError(err1)
	}
	value, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, logId)))
	easygo.PanicError(err)
	return value
}

func GetRedisPlayerBeZanIdInfo(playerId int64) map[int64][]int64 { //获取玩家被赞动态的id
	info := make(map[int64][]int64)
	value, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(REDIS_SQUARE_DYNAMICZANINFO, playerId)))
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

func GetRedisPlayerBeZanIdListForLogId(playerId, logId int64) []int64 { //获取玩家被赞动态的id
	var lst []int64
	b, err := easygo.RedisMgr.GetC().HExists(MakeRedisKey(REDIS_SQUARE_DYNAMICZANINFO, playerId), easygo.AnytoA(logId))
	easygo.PanicError(err)
	if !b {
		return lst
	}
	value, err := easygo.RedisMgr.GetC().HGet(MakeRedisKey(REDIS_SQUARE_DYNAMICZANINFO, playerId), easygo.AnytoA(logId))
	easygo.PanicError(err)

	err1 := json.Unmarshal(value, &lst)
	easygo.PanicError(err1)
	return lst
}

func AddRedisPlayerBeZanId(playerId, logId, Id int64) {
	lst := GetRedisPlayerBeZanIdListForLogId(playerId, logId)
	lst = append(lst, Id)
	b, _ := json.Marshal(lst)
	err := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZANINFO, playerId), easygo.AnytoA(logId), string(b))
	easygo.PanicError(err)
}

func DelRedisPlayerBeZanId(playerId, logId, Id int64) {
	lst := GetRedisPlayerBeZanIdListForLogId(playerId, logId)
	if len(lst) == 0 {
		return
	}
	var ind int
	for index, id := range lst {
		if id == Id {
			ind = index
			break
		}
	}

	newlst := append(lst[:ind], lst[ind+1:]...)
	b, _ := json.Marshal(newlst)
	err := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZANINFO, playerId), easygo.AnytoA(logId), string(b))
	easygo.PanicError(err)
}

func AddRedisPlayerBeZanIdInfo(playerId int64, m map[int64][]int64) error { //加载所有被赞的
	info := make(map[int64]string) //因为redis hashmap中value必须是string类型  所以把[]int64 json成字符串存进去
	for key, lst := range m {
		b, err := json.Marshal(lst)
		if err != nil {
			logs.Error(err)
			continue
		}
		info[key] = string(b)
	}
	err := easygo.RedisMgr.GetC().HMSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZANINFO, playerId), info)
	return err
}

func GetRedisDynamicDelZanId() []int64 { //获取取消赞的日志id
	ids := []int64{}
	value, err := easygo.RedisMgr.GetC().Smembers(REDIS_SQUARE_DELZANIDS)
	easygo.PanicError(err)
	InterfersToInt64s(value, &ids)
	return ids
}

func AddRedisDynamicDelZanId(id int64) {
	err := easygo.RedisMgr.GetC().SAdd(REDIS_SQUARE_DELZANIDS, id)
	easygo.PanicError(err)
}

func ClearRedisDynamicDelZanId() {
	_, err := easygo.RedisMgr.GetC().Delete(REDIS_SQUARE_DELZANIDS)
	easygo.PanicError(err)
}

func GetRedisDynamicIsZan(logId, playerId int64) bool { //判断用户是否赞过该动态
	ids := GetRedisPlayerZanIds(playerId)
	if util.Int64InSlice(logId, ids) {
		return true
	}
	return false
}

func GetRedisPlayerZanIds(playerId int64) []int64 { //获取人物赞过的动态id列表  倒序
	ids := []int64{}
	value, err := easygo.RedisMgr.GetC().Smembers(MakeRedisKey(REDIS_SQUARE_MYZANIDS, playerId))
	easygo.PanicError(err)
	InterfersToInt64s(value, &ids)
	easygo.SortSliceInt64(ids, false)
	return ids
}

func AddRedisPlayerZanIds(playerId, id int64) { //增加人物赞过的动态id
	err := easygo.RedisMgr.GetC().SAdd(MakeRedisKey(REDIS_SQUARE_MYZANIDS, playerId), id)
	easygo.PanicError(err)
}

func DelRedisPlayerZanIds(playerId, id int64) { //删除人物赞过的动态id列表
	err := easygo.RedisMgr.GetC().SRem(MakeRedisKey(REDIS_SQUARE_MYZANIDS, playerId), id)
	easygo.PanicError(err)
}

//判断用户是否赞过该主评论
func GetRedisMainCommentIsZan(mainCommentId, playerId int64) bool {
	ids := GetRedisPlayerMainCommentZanIds(playerId)
	if util.Int64InSlice(mainCommentId, ids) {
		return true
	}
	return false
}

//获取人物赞过的主评论id列表  倒序
func GetRedisPlayerMainCommentZanIds(playerId int64) []int64 { //获取人物赞过的主评论id列表  倒序
	ids := []int64{}
	value, err := easygo.RedisMgr.GetC().Smembers(MakeRedisKey(REDIS_SQUARE_MAIN_COMMENT_MYZANIDS, playerId))
	easygo.PanicError(err)
	InterfersToInt64s(value, &ids)
	easygo.SortSliceInt64(ids, false)
	return ids
}

//增加人物赞过的主评论id
func AddRedisPlayerMainCommentZanIds(playerId, mainCommentId int64) {
	err := easygo.RedisMgr.GetC().SAdd(MakeRedisKey(REDIS_SQUARE_MAIN_COMMENT_MYZANIDS, playerId), mainCommentId)
	easygo.PanicError(err)
}

//删除人物赞过的主评论id列表
func DelRedisPlayerMainCommentZanIds(playerId, mainCommentId int64) {
	err := easygo.RedisMgr.GetC().SRem(MakeRedisKey(REDIS_SQUARE_MAIN_COMMENT_MYZANIDS, playerId), mainCommentId)
	easygo.PanicError(err)
}

func GetRedisDynamicAllZanInfo(logId, maxId int64) []*share_message.ZanData {
	info := GetRedisDynamicZan(logId)
	lst := make([]*share_message.ZanData, 0)
	for _, s := range info {
		var msg *share_message.ZanData
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			logs.Error(err)
			continue
		}
		if maxId > msg.GetLogId() {
			continue
		}
		lst = append(lst, msg)
	}
	return lst

}

func OperateRedisDynamicZan(operate int32, operateId, logId, dynamicPid int64, sysp *share_message.SysParameter) bool { //点赞操作  operateId 代表操作的人的id
	if operate == DYNAMIC_OPERATE { //operate 是操作  1是点赞  2是取消点赞
		//Id := NextId(REDIS_SQUARE_DYNAMICZAN)
		Id := NextId(TABLE_SQUARE_ZAN)
		msg := &share_message.ZanData{
			LogId:      easygo.NewInt64(Id),
			OperateId:  easygo.NewInt64(operateId),
			CreateTime: easygo.NewInt64(time.Now().Unix()),
		}
		msg.PlayerId = easygo.NewInt64(dynamicPid)
		msg.DynamicId = easygo.NewInt64(logId)
		b, err := json.Marshal(msg)
		easygo.PanicError(err)
		err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, logId), easygo.AnytoA(Id), string(b))
		easygo.PanicError(err1)
		AddRedisPlayerBeZanId(dynamicPid, logId, Id) // 记录的是 我发的这条动态  谁赞过我
		AddRedisPlayerZanIds(operateId, logId)       // 记录的是  我赞过的动态
		target := GetRedisPlayerBase(dynamicPid)
		target.AddZan(1)
		// 自己点赞的不添加在未读消息里面
		if operateId != dynamicPid {
			AddPlayerUnreadInfo(dynamicPid, UNREAD_ZAN, 1)
		}
		// 点赞推送
		easygo.Spawn(func() {
			// 如果是自己,不推送
			if operateId == dynamicPid {
				return
			}
			if dynamicPlayer := GetRedisPlayerBase(dynamicPid); dynamicPlayer != nil {
				// 如果开关关了,不推送
				if dynamicPlayer.GetIsOpenSquare() {
					return
				}
			}
			oper := GetRedisPlayerBase(operateId)
			ids := GetJGIds([]int64{dynamicPid})
			m := PushMessage{
				//Title:       "温馨提示",
				ContentType: JG_TYPE_SQUARE,
				Location:    7,             // 指定动态
				ItemId:      PUSH_ITEM_201, // 点赞我的动态
				ObjectId:    logId,
			}

			if sysp == nil || len(sysp.GetPushSet()) == 0 {
				return
			}

			for _, ps := range sysp.GetPushSet() {
				if ps.GetObjId() == m.ItemId && ps.GetIsPush() {
					m.Content = fmt.Sprintf(ps.GetObjContent(), oper.GetNickName())
				}
			}
			JGSendMessage(ids, m, sysp)

		})

		return true
	}

	var id int64
	info := GetRedisDynamicZan(logId)
	for _, s := range info {
		var msg *share_message.ZanData
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			logs.Error(err)
			continue
		}
		if msg.GetOperateId() == operateId {
			id = msg.GetLogId()
			break
		}
	}

	if id == 0 {
		return false
	}

	// 删除赞数据
	DelZanDataFromRedisAndDB(logId, id)

	target := GetRedisPlayerBase(dynamicPid)
	target.AddZan(-1)
	DelRedisPlayerZanIds(operateId, logId)
	DelRedisPlayerBeZanId(dynamicPid, logId, id)
	AddRedisDynamicDelZanId(id)

	// 自己点赞的不添加在未读消息里面
	if operateId != dynamicPid {
		AddPlayerUnreadInfo(dynamicPid, UNREAD_ZAN, -1)
	}
	return true
}

// =======================================================================

func GetRedisSquareAllZanInfo(playerId, minId int64) *client_square.ZanList { //获取赞的信息
	info := GetRedisPlayerBeZanIdInfo(playerId) // [动态id][]{赞id1,赞id2}
	var lst []*share_message.ZanData
	var dynamicIds []int64
	var dynamicList []*share_message.DynamicData
	zanInfo := make(map[int64]int64) // [赞id]动态id
	for logId, ids := range info {
		for _, id := range ids {
			zanInfo[id] = logId
		}
	}
	var zanIds []int64
	for id := range zanInfo {
		if id >= minId && minId != 0 {
			continue
		}
		zanIds = append(zanIds, id)
	}
	easygo.SortSliceInt64(zanIds, false)

	var num int32
	player := GetRedisPlayerBase(playerId)
	name := player.GetNickName()
	headIcon := player.GetHeadIcon()
	types := player.GetTypes()
	sex := player.GetSex()
	var operateIds []int64
	for _, id := range zanIds {
		if num >= DYNAMIC_REQUEST_NUM {
			break
		}
		logId := zanInfo[id]
		if logId == 0 {
			logs.Error("动态id怎么会是0,id:", logId, id)
			continue
		}
		b, err := easygo.RedisMgr.GetC().HExists(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, logId), easygo.AnytoA(id))
		if err != nil {
			logs.Error(err)
			continue
		}
		if !b { // 从数据库中查找.
			logs.Error("动态id: %d,怎么会没有赞id:%d 在redis中,从数据库中查询.", logId, id)
			zd := GetZanDataByIdAndDynamicId(id, logId)
			if zd == nil {
				continue
			}
			// 放进redis
			b, _ := json.Marshal(zd)
			if err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, logId), easygo.AnytoA(zd.GetLogId()), string(b)); err1 != nil {
				easygo.PanicError(err1)
				continue
			}
		}
		value, err := easygo.RedisMgr.GetC().HGet(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, logId), easygo.AnytoA(id))
		if err != nil {
			logs.Error(err)
			continue
		}
		var msg *share_message.ZanData
		err1 := json.Unmarshal(value, &msg)
		if err1 != nil {
			logs.Error(err1)
			continue
		}

		if msg.GetOperateId() == playerId { //如果是自己的点赞就不显示了
			continue
		}
		if !util.Int64InSlice(msg.GetOperateId(), operateIds) { //一次性获取所有的点赞人物id
			operateIds = append(operateIds, msg.GetOperateId())
		}
		if !util.Int64InSlice(logId, dynamicIds) {
			dynamicIds = append(dynamicIds, logId)
			dynamic := GetRedisDynamic(logId)
			if dynamic == nil {
				continue
			}
			dynamic.NickName = easygo.NewString(name)
			dynamic.HeadIcon = easygo.NewString(headIcon)
			dynamic.Sex = easygo.NewInt32(sex)
			dynamic.Types = easygo.NewInt32(types)
			dynamicList = append(dynamicList, dynamic)
		}
		num += 1
		lst = append(lst, msg)

	}

	playerInfo := GetAllPlayerBase(operateIds)
	for _, m := range lst {
		info := playerInfo[m.GetOperateId()]
		m.Sex = easygo.NewInt32(info.GetSex())
		m.Name = easygo.NewString(info.GetNickName())
		m.HeadIcon = easygo.NewString(info.GetHeadIcon())
		m.Types = easygo.NewInt32(info.GetTypes())
	}

	msg := &client_square.ZanList{
		ZanData:     lst,
		DynamicData: dynamicList,
	}

	return msg
}

//==============================savedatabase==============================

func MakeSquareInfoKey(k string) string {
	return MakeRedisKey(REDIS_SQUARE_SAVEINFO, k)
}

func UpdateDynamicSaveZanId(maxId int64) {
	err := easygo.RedisMgr.GetC().HSet(MakeSquareInfoKey(REDIS_SQUARE_SAVEZANID), REDIS_SQUARE_SAVEZANID, maxId)
	easygo.PanicError(err)
}

func GetDynamicSaveZanId() int64 {
	b, err := easygo.RedisMgr.GetC().HExists(MakeSquareInfoKey(REDIS_SQUARE_SAVEZANID), REDIS_SQUARE_SAVEZANID)
	easygo.PanicError(err)
	if !b {
		return 0
	}
	Id, err1 := easygo.RedisMgr.GetC().HGet(MakeSquareInfoKey(REDIS_SQUARE_SAVEZANID), REDIS_SQUARE_SAVEZANID)
	easygo.PanicError(err1)
	return easygo.AtoInt64(string(Id))
}

//判断玩家是否对评论进行点赞过
func CheckIsDynamicCommentZan(playerId, dynamicId, commentId int64) int64 {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT_ZAN)
	defer closeFunc()
	var comment *share_message.CommentDataZan
	err := col.Find(bson.M{"PlayerId": playerId, "DynamicId": dynamicId, "CommentId": commentId}).One(&comment)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return comment.GetId()
}

func GetSomeZanInfoFromRedis() ([]string, []*share_message.ZanData) {
	keys, err := easygo.RedisMgr.GetC().Scan(REDIS_SQUARE_DYNAMICZAN) //存储所有赞的信息
	easygo.PanicError(err)
	var lst []*share_message.ZanData
	for _, key := range keys {
		value, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(key))
		easygo.PanicError(err)
		for _, s := range value {
			var msg *share_message.ZanData
			err := json.Unmarshal([]byte(s), &msg)
			if err != nil {
				logs.Error(err)
				continue
			}
			lst = append(lst, msg)
		}
	}

	return keys, lst
}

// 删除redis中的赞数据.
func DelZanDatasFromRedis(keys []string) {
	for _, v := range keys {
		b, err := easygo.RedisMgr.GetC().Exist(v)
		if err != nil || !b {
			logs.Error(err)
			continue
		}
		b, err1 := easygo.RedisMgr.GetC().Delete(v)
		easygo.PanicError(err1)
	}
}

// 定时储存赞的数据
func TaskSaveZanData() {
	// 获得redis中的赞数据
	keys, zanS := GetSomeZanInfoFromRedis()
	// 存储进数据库
	UpsetAllZanDataToDB(zanS)
	// 删除redis中的数据
	DelZanDatasFromRedis(keys)
}

// 从redis中删除赞数据,并且删除数据库中的数据.
func DelZanDataFromRedisAndDB(dynamicId, zanId int64) bool {
	_, err1 := easygo.RedisMgr.GetC().Hdel(MakeRedisKey(REDIS_SQUARE_DYNAMICZAN, dynamicId), easygo.AnytoA(zanId))
	easygo.PanicError(err1)
	// 删除数据库
	err := DelZanDataToDB(zanId)
	if err != nil {
		return false
	}
	return true
}
