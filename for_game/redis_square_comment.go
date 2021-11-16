package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/pb/client_server"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

/*
社交广场 redis数据
author : 狗哥
*/
//todo  删除redis数据机制

const (
	DYNAMIC_COMMENT_STATUE_COMMON        = 0 //正常的评论
	DYNAMIC_COMMENT_STATUE_DELETE        = 1 //被删除的评论
	DYNAMIC_COMMENT_STATUE_DELETE_CLIENT = 2 //被删除的评论,app端删除.
	UNREAD_COMMENT                       = 1
)

const (
	//REDIS_SQUARE_COMMENT = "redis_square_comment" //社交广场动态评论 map[redis_square_comment_logid]map[id]string // 原来的,暂时不删
	REDIS_SQUARE_COMMENT               = "redis_square:comment"          //社交广场动态评论 map[redis_square_comment_logid]map[id]string
	REDIS_SQUARE_MAIN_COMMENT_MYZANIDS = "redis_square:comment_myzanids" //我点过赞的主评论id
	//存储限制key
	REDIS_SQUARE_SAVECOMMENTID     = "savecommentid"                  //上次存库的最后评论id lst []int64
	REDIS_SQUARE_COMMENT_MAX_LOGID = "redis_square:comment_max_logid" //评论获取id
)

/**
xiong
加载动态评论进redis,并记录这次存储最大的游标位置.
*/
func ReloadDynamicComment(logIds []int64, ch chan int) {
	defer SquareRecoverAndLog(ch)
	var newlst []int64
	for _, logId := range logIds { //排除redis中已经存在的评论
		b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(REDIS_SQUARE_COMMENT, logId))
		if err != nil {
			logs.Error(err)
			continue
		}
		if !b {
			newlst = append(newlst, logId)
		}
	}
	if len(newlst) == 0 {
		return
	}
	commentInfo, err2 := GetCommentListByLogIdsAndStatusFromDB(newlst, DYNAMIC_COMMENT_STATUE_COMMON) // 根据动态id列表和状态获取评论列表
	easygo.PanicError(err2)
	info1 := make(map[int64][]*share_message.CommentData)
	for _, m := range commentInfo {
		logId := m.GetLogId()
		info1[logId] = append(info1[logId], m)
	}
	var maxId int64
	for logId, lst := range info1 { //加载所有评论到redis
		info2 := make(map[int64]string)
		for _, m := range lst {
			id := m.GetId()
			b, _ := json.Marshal(m)
			info2[id] = string(b)
			if id > maxId {
				maxId = id
			}
		}
		err3 := easygo.RedisMgr.GetC().HMSet(MakeRedisKey(REDIS_SQUARE_COMMENT, logId), info2)
		easygo.PanicError(err3)
	}
	if maxId > GetDynamicSaveCommentId() { // 记录上一次库存最大的id
		UpdateDynamicSaveCommentId(maxId)
	}

}

/**
xiong
登陆时加载新的评论进redis
*/
func ReloadMyCommentInfo(logIds []int64, ch chan int, playerId int64) { //加载自己发布的动态评论数据
	defer SquareRecoverAndLog(ch)
	info := make(map[int64][]int64)
	newlst := []int64{}
	for _, logId := range logIds { //如果某条动态的评论已经在redis中就不查数据库  直接从redis中找
		b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(REDIS_SQUARE_COMMENT, logId))
		if err != nil {
			logs.Error(err)
			continue
		}
		if !b {
			newlst = append(newlst, logId)
		} else {
			values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(REDIS_SQUARE_COMMENT, logId)))
			easygo.PanicError(err)
			for id, _ := range values {
				info[logId] = append(info[logId], id)
			}
		}
	}
	if len(newlst) > 0 { //如果不在redis中 就去mongo中找
		commentInfo, err2 := GetCommentListByLogIdsAndStatusFromDB(newlst, DYNAMIC_COMMENT_STATUE_COMMON)
		easygo.PanicError(err2)
		info1 := make(map[int64][]*share_message.CommentData)
		for _, m := range commentInfo {
			logId := m.GetLogId()
			info1[logId] = append(info1[logId], m)
		}
		for logId, lst := range info1 { //加载所有评论到redis
			info2 := make(map[int64]string)
			for _, m := range lst {
				id := m.GetId()
				b, _ := json.Marshal(m)
				info2[id] = string(b)
				info[logId] = append(info[logId], m.GetId())
			}
			err3 := easygo.RedisMgr.GetC().HMSet(MakeRedisKey(REDIS_SQUARE_COMMENT, logId), info2)
			easygo.PanicError(err3)
		}
	}

	if len(info) != 0 { //加载个人动态 回复自己的消息id列表
		AddRedisPlayerMessageInfo(playerId, info)
	}
}

/**
xiong
从redis中删除动态的评论
*/
func DelDynamicComment(logId int64) {
	redisKey := MakeRedisKey(REDIS_SQUARE_COMMENT, easygo.AnytoA(logId))
	b, err := easygo.RedisMgr.GetC().Exist(redisKey)
	easygo.PanicError(err)
	if !b {
		return
	}
	b1, err3 := easygo.RedisMgr.GetC().Delete(redisKey) //从redis中删除动态评论
	easygo.PanicError(err3)
	if !b1 {
		logs.Error("redis删除动态评论失败id:", logId)
	}
}

// ===========================评论相关====================================
func GetRedisDynamicCommentNum(logId int64) int64 { //获取评论总数
	redisCon := easygo.RedisMgr.GetC()
	value, err := redisCon.HLen(MakeRedisKey(REDIS_SQUARE_COMMENT, logId))
	easygo.PanicError(err)
	count := int64(value)
	if int64(value) <= 0 { // 从数据库查询.如果是有,再放进redis
		data, err := GetCommentListByLogIdsAndStatusFromDB([]int64{logId}, DYNAMIC_STATUE_COMMON)
		easygo.PanicError(err)
		if len(data) > 0 {
			count = int64(len(data))
			// 放进redis
			m := make(map[int64]string)
			for _, v := range data {
				bytes, _ := json.Marshal(v)
				m[v.GetId()] = string(bytes)
			}
			err = redisCon.HMSet(MakeRedisKey(REDIS_SQUARE_COMMENT, logId), m)
			easygo.PanicError(err)
		}
	}
	return count
}

// todo 备份优化前.
//func GetRedisDynamicCommentNumEx(logIds []int64) map[int64]int32 { //获取评论总数
//	m := make(map[int64]int32)
//	unFindList := make([]int64, 0)
//	redisCon := easygo.RedisMgr.GetC()
//	for _, id := range logIds {
//		value, err := redisCon.HLen(MakeRedisKey(REDIS_SQUARE_COMMENT, id))
//		easygo.PanicError(err)
//		if int32(value) <= 0 {
//			unFindList = append(unFindList, id)
//		} else {
//			m[id] = int32(value)
//		}
//	}
//	if len(unFindList) > 0 {
//		data, err := GetCommentListByLogIdsAndStatusFromDB(unFindList, DYNAMIC_STATUE_COMMON)
//		easygo.PanicError(err)
//		if len(data) > 0 {
//			for _, d := range data {
//				count := m[d.GetLogId()]
//				count += 1
//				m[d.GetLogId()] = count
//
//				m1 := make(map[int64]map[int64]string)
//				// 封装同一个动态的赞信息
//				for _, v := range data {
//					did := v.GetLogId()
//					vv, ok := m1[did]
//					if !ok {
//						vv = make(map[int64]string)
//					}
//					bytes, _ := json.Marshal(v)
//					vv[v.GetId()] = string(bytes)
//					m1[did] = vv
//				}
//				for key, value := range m1 {
//					err = redisCon.HMSet(MakeRedisKey(REDIS_SQUARE_COMMENT, key), value)
//					easygo.PanicError(err)
//				}
//
//			}
//		}
//	}
//	return m
//
//}
func GetRedisDynamicCommentNumEx(logIds []int64) map[int64]int32 { //获取评论总数
	m := make(map[int64]int32)
	unFindList := make([]int64, 0)
	redisCon := easygo.RedisMgr.GetC()
	for _, id := range logIds {
		value, err := redisCon.HLen(MakeRedisKey(REDIS_SQUARE_COMMENT, id))
		easygo.PanicError(err)
		if int32(value) <= 0 {
			unFindList = append(unFindList, id)
		} else {
			m[id] = int32(value)
		}
	}
	if len(unFindList) > 0 {
		data, err := GetCommentListByLogIdsAndStatusFromDB(unFindList, DYNAMIC_STATUE_COMMON)
		easygo.PanicError(err)
		if len(data) > 0 {
			for _, d := range data {
				count := m[d.GetLogId()]
				count += 1
				m[d.GetLogId()] = count
			}

			fun := func(data []*share_message.CommentData) {
				m1 := make(map[int64]map[int64]string)
				// 封装同一个动态的赞信息
				for _, v := range data {
					did := v.GetLogId()
					vv, ok := m1[did]
					if !ok {
						vv = make(map[int64]string)
					}
					bytes, _ := json.Marshal(v)
					vv[v.GetId()] = string(bytes)
					m1[did] = vv
				}
				for key, value := range m1 {
					err = redisCon.HMSet(MakeRedisKey(REDIS_SQUARE_COMMENT, key), value)
					easygo.PanicError(err)
				}
			}

			easygo.Spawn(fun, data)

		}
	}
	return m

}

//获取评论
func GetRedisDynamicComment(logId, id int64) *share_message.CommentData { //获取单条评论的数据
	b, err := easygo.RedisMgr.GetC().HExists(MakeRedisKey(REDIS_SQUARE_COMMENT, logId), easygo.AnytoA(id))
	easygo.PanicError(err)
	if !b {
		return nil
	}
	value, err1 := easygo.RedisMgr.GetC().HGet(MakeRedisKey(REDIS_SQUARE_COMMENT, logId), easygo.AnytoA(id))
	easygo.PanicError(err1)
	var comment *share_message.CommentData
	err2 := json.Unmarshal(value, &comment)
	easygo.PanicError(err2)
	return comment
}

func GetRedisDynamicAllComment(logId, maxId int64) []*share_message.CommentData { //获取大于maxId 的所有评论数据
	lst := make([]*share_message.CommentData, 0)
	value, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(REDIS_SQUARE_COMMENT, logId)))
	easygo.PanicError(err)
	for _, s := range value {
		var comment *share_message.CommentData
		err1 := json.Unmarshal([]byte(s), &comment)
		if err1 != nil {
			logs.Error(err1)
			continue
		}
		if comment.GetId() < maxId {
			continue
		}
		lst = append(lst, comment)
	}
	return lst
}

/**
xiong
添加评论进redis 和 mongodb
*/
func AddRedisDynamicComment(msg *share_message.CommentData) bool {
	err2 := InsertCommentToDB(msg)
	easygo.PanicError(err2)

	b, err := json.Marshal(msg)
	if err != nil {
		logs.Error(err)
		return false
	}

	err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_COMMENT, msg.GetLogId()), easygo.AnytoA(msg.GetId()), string(b))
	if err1 != nil {
		logs.Error(err1)
		return false
	}
	return true
}

//修改评论赞数
func ModifyRedisDynamicCommentScore(commentId, val int64, isComment ...bool) {
	//修改库
	comment := IncrCommentZanNumToDB(commentId, val, isComment...)
	//修改redis
	b, err := json.Marshal(comment)
	if err != nil {
		logs.Error(err)
		return
	}
	err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_COMMENT, comment.GetLogId()), easygo.AnytoA(comment.GetId()), string(b))
	if err1 != nil {
		logs.Error(err1)
	}
}

// deleteStatus 1-后台管理系统删除; 2-app端删除.
func DelRedisDynamicComment(logId, Id, deleteStatus int64, note ...string) bool { //删除评论
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT)
	defer closeFunc()
	noteStr := append(note, "")[0]
	err := col.Update(bson.M{"_id": Id}, bson.M{"$set": bson.M{"Statue": deleteStatus, "Note": noteStr}})
	if err != nil {
		logs.Error("删除评论数据库失败id:", Id, err)
		return false
	}
	_, err1 := easygo.RedisMgr.GetC().Hdel(MakeRedisKey(REDIS_SQUARE_COMMENT, logId), easygo.AnytoA(Id))
	if err1 != nil {
		logs.Error("删除评论redis失败id:", Id, err1)
		return false
	}

	log := &share_message.CommentData{}
	iter := col.Find(bson.M{"BelongId": Id}).Iter()
	for iter.Next(&log) {
		_, err1 = easygo.RedisMgr.GetC().Hdel(MakeRedisKey(REDIS_SQUARE_COMMENT, logId), easygo.AnytoA(log.GetId()))
		if err1 != nil {
			logs.Error("删除评论redis失败id:", Id, err1)
			return false
		}
	}

	if err := iter.Close(); err != nil {
		easygo.PanicError(err)
	}

	_, err = col.UpdateAll(bson.M{"BelongId": Id}, bson.M{"$set": bson.M{"Statue": deleteStatus, "Note": noteStr}})
	if err != nil {
		logs.Error("删除评论数据库失败id:", Id, err)
		return false
	}
	return true
}

//通过页码获取动态评论
func GetRedisDynamicCommentInfoByPage(pid, logId int64, reqMsg *client_server.IdInfo, params *share_message.SysParameter) *share_message.CommentList {
	page := reqMsg.GetPage()
	size := reqMsg.GetPageSize()
	hotList := reqMsg.GetHotList()
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if size == 0 {
		size = DYNAMIC_REQUEST_NUM
	}
	result := GetDynamicCommentInfoByPage(logId, int(page), int(size), hotList, params)
	// 判断是否有跳转的评论id,如果有,插入进第一条.
	if reqMsg.GetJumpMainCommentId() > 0 {
		cd := GetCommentDataFromDB(reqMsg.GetJumpMainCommentId()) // 评论
		// 如果是子评论,就要查询出主评论
		if cd.GetBelongId() > 0 {
			cd = GetCommentDataFromDB(cd.GetBelongId())
		}
		if cd != nil { // 重新复制
			cl := new(share_message.CommentList)
			cl.HotList = result.GetHotList()
			cl.CommentInfo = append(cl.CommentInfo, cd)
			for _, v := range result.GetCommentInfo() {
				cl.CommentInfo = append(cl.CommentInfo, v)
			}
			result = cl
		}
	}
	//  判断自己是否已关注
	for _, c := range result.CommentInfo {
		c.IsZan = easygo.NewBool(GetRedisMainCommentIsZan(c.GetId(), pid))
		// 封装评论信息.
		player := GetRedisPlayerBase(c.GetPlayerId())
		c.Sex = easygo.NewInt32(player.GetSex())
		c.Name = easygo.NewString(player.GetNickName())
		c.HeadIcon = easygo.NewString(player.GetHeadIcon())
		c.Types = easygo.NewInt32(player.GetTypes())
		// 查询子评论数量.
		if num := GetSecondaryCommentNum(c.GetId()); num > 0 {
			c.TotalNum = easygo.NewInt64(num)
		}
	}

	return result
}

// 获取大于Id的50条主评论
func GetRedisDynamicCommentInfo(logId, mainId int64) *share_message.CommentList {
	var lst []int64
	info := make(map[int64]int64)
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(REDIS_SQUARE_COMMENT, logId)))
	easygo.PanicError(err)
	for id, s := range values { // key-->评论的id,value-->评论的内容
		var comment *share_message.CommentData
		_ = json.Unmarshal([]byte(s), &comment)
		if comment.GetStatue() != DYNAMIC_COMMENT_STATUE_COMMON {
			continue
		}
		if mainId == 0 {
			if comment.GetBelongId() == 0 {
				lst = append(lst, id)
			}
		} else {
			if comment.GetBelongId() == 0 && id < mainId {
				lst = append(lst, id)
			}
		}
		if comment.GetBelongId() != 0 {
			info[comment.GetBelongId()] += 1
		}

	}
	if len(lst) == 0 {
		return &share_message.CommentList{}
	}

	easygo.SortSliceInt64(lst, false) //顺序
	var commentList []*share_message.CommentData
	var count int
	for _, id := range lst {
		if count >= DYNAMIC_REQUEST_NUM {
			break
		}
		comment := GetOneCommentInfo(logId, id)
		comment.TotalNum = easygo.NewInt64(info[id])
		commentList = append(commentList, comment)
		count += 1
	}
	msg := &share_message.CommentList{
		CommentInfo: commentList,
	}
	return msg
}

// 获取mainId主评论下的大于secondId的子评论
func GetRedisDynamicSecondComment(logId, mainId, secondId int64) *share_message.CommentList {
	lst := make([]int64, 0)
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(REDIS_SQUARE_COMMENT, logId)))
	easygo.PanicError(err)
	for id, s := range values { // 这里以后如果数据量大的话  可以考虑用chan开多协程并发去计算
		var comment *share_message.CommentData
		_ = json.Unmarshal([]byte(s), &comment)
		if comment.GetBelongId() == mainId && comment.GetId() > secondId {
			lst = append(lst, id)
		}
	}
	if len(lst) == 0 {
		return &share_message.CommentList{}
	}
	easygo.SortSliceInt64(lst, false) //顺序
	commentList := make([]*share_message.CommentData, 0)
	var count int
	for _, id := range lst {
		if count >= DYNAMIC_REQUEST_NUM {
			break
		}
		comment := GetOneCommentInfo(logId, id)
		commentList = append(commentList, comment)
		count += 1
	}
	msg := &share_message.CommentList{
		CommentInfo: commentList,
	}
	return msg
}

// =======================================================================

func PackageCommentDataInfo(msg *share_message.CommentData) { //封装评论信息
	pid := msg.GetPlayerId()
	player := GetRedisPlayerBase(pid)
	msg.Sex = easygo.NewInt32(player.GetSex())
	msg.Name = easygo.NewString(player.GetNickName())
	msg.HeadIcon = easygo.NewString(player.GetHeadIcon())
	msg.Types = easygo.NewInt32(player.GetTypes())
	if msg.GetBelongId() != 0 { //代表是子评论
		targetId := msg.GetTargetId()
		target := GetRedisPlayerBase(targetId)
		if target == nil {
			return
		}
		msg.OtherSex = easygo.NewInt32(target.GetSex())
		msg.OtherName = easygo.NewString(target.GetNickName())
	}
}

func GetOneCommentInfo(logId, Id int64) *share_message.CommentData {
	comment := GetRedisDynamicComment(logId, Id)
	if comment == nil {
		return comment
	}
	PackageCommentDataInfo(comment)
	return comment
}

//==============================savedatabase==============================
func UpdateDynamicSaveCommentId(maxId int64) {
	err := easygo.RedisMgr.GetC().HSet(MakeSquareInfoKey(REDIS_SQUARE_SAVECOMMENTID), REDIS_SQUARE_SAVECOMMENTID, maxId)
	easygo.PanicError(err)
}

func GetDynamicSaveCommentId() int64 {
	b, err := easygo.RedisMgr.GetC().HExists(MakeSquareInfoKey(REDIS_SQUARE_SAVECOMMENTID), REDIS_SQUARE_SAVECOMMENTID)
	easygo.PanicError(err)
	if !b {
		return 0
	}
	Id, err1 := easygo.RedisMgr.GetC().HGet(MakeSquareInfoKey(REDIS_SQUARE_SAVECOMMENTID), REDIS_SQUARE_SAVECOMMENTID)
	easygo.PanicError(err1)
	return easygo.AtoInt64(string(Id))
}
