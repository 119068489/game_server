package for_game

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/client_square"
	"game_server/pb/share_message"
	"sort"
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
	UNREAD_ATTENTION = 3
)

const (

	//REDIS_SQUARE_ATTENTION       = "redis_square_attention"       //玩家被关注信息 map[redis_square_attention_pid]map[id]string
	REDIS_SQUARE_ATTENTION       = "redis_square:attention"       //玩家被关注信息 map[redis_square_attention_pid]map[id]string // 原来的,暂时不删 redis_square_attention
	REDIS_SQUARE_DELATTENTIONIDS = "redis_square:delattentionids" //取消关注的id
	//存储限制key
	REDIS_SQUARE_SAVEATTID = "saveattid" //上次存库的最后关注id
)

/**
xiong
社交广场登录时
加载谁关注了我的信息
*/
func ReloadPlayerAttention(playerId int64, ch chan int) {
	defer SquareRecoverAndLog(ch)
	lst, err := GetAttentionListByPlayerIdFromDB(playerId)
	easygo.PanicError(err)
	if len(lst) == 0 {
		return
	}
	var maxId int64
	for _, msg := range lst {
		id := msg.GetLogId()
		if maxId < id {
			maxId = id
		}
		b, err := json.Marshal(msg)
		if err != nil {
			logs.Error(err)
			continue
		}
		err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_ATTENTION, playerId), easygo.AnytoA(id), string(b))
		if err1 != nil {
			logs.Error(err1)
		}
	}
	if maxId > GetDynamicSaveAttId() {
		UpdateDynamicSaveAttId(maxId)
	}

}

// =========================动态相关===========================

// ==========================Operate===============================

func GetRedisDynamicDelAttId() []int64 { //获取取消关注的日志id
	ids := make([]int64, 0)
	value, err := easygo.RedisMgr.GetC().Smembers(REDIS_SQUARE_DELATTENTIONIDS)
	easygo.PanicError(err)
	InterfersToInt64s(value, &ids)
	return ids
}

func AddRedisDynamicDelAttId(id int64) { //增加取消关注的日志id
	err := easygo.RedisMgr.GetC().SAdd(REDIS_SQUARE_DELATTENTIONIDS, id)
	easygo.PanicError(err)
}

func ClearRedisDynamicDelAttId() { //清除取消关注的日志id
	_, err := easygo.RedisMgr.GetC().Delete(REDIS_SQUARE_DELATTENTIONIDS)
	easygo.PanicError(err)
}

// 操作关注操作
func OperateRedisDynamicAttention(operate int32, operateId, playerId int64) bool {
	return OperateRedisDynamicAttentionEx(operate, operateId, playerId, 0)
}

// 操作关注操作 值说明： source=1 电竞来源 ，0或者其他im来源
func OperateRedisDynamicAttentionEx(operate int32, operateId, playerId int64, source int32) bool {
	if operate == DYNAMIC_OPERATE { //operate 是操作  1是关注  2是取消关注
		//Id := NextId(REDIS_SQUARE_ATTENTION)
		Id := NextId(TABLE_SQUARE_ATTENTION)
		msg := &share_message.AttentionData{
			LogId:      easygo.NewInt64(Id),
			PlayerId:   easygo.NewInt64(playerId),
			OperateId:  easygo.NewInt64(operateId),
			CreateTime: easygo.NewInt64(time.Now().Unix()),
			Source:     easygo.NewInt32(source),
		}
		b, err := json.Marshal(msg)
		easygo.PanicError(err)
		// 判断当前关注列表中是否已存在了
		for _, attentionData := range GetRedisSquareAllAttentionInfoNoLimit(playerId) {
			if attentionData.GetOperateId() == operateId {
				logs.Error("redis 中存在相同的操作人员,redis key 为: redis_square_attention_%d,操作者为: %d", playerId, operateId)
				return false
			}
		}
		err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_ATTENTION, playerId), easygo.AnytoA(Id), string(b))
		easygo.PanicError(err1)
		AddPlayerUnreadInfo(playerId, UNREAD_ATTENTION, 1)
		return true
	}

	// 2 是取消关注
	operater := GetRedisPlayerBase(operateId)
	operateAttention := operater.GetAttention()
	if !util.Int64InSlice(playerId, operateAttention) {
		return true
	}

	value, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(REDIS_SQUARE_ATTENTION, playerId)))
	easygo.PanicError(err)
	var id int64
	for _, s := range value {
		var msg *share_message.AttentionData
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
	if id == 0 { // redis 中找不到,从数据库中找
		col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_ATTENTION)
		defer closeFunc()
		var attentionData *share_message.AttentionData
		err := col.Find(bson.M{"PlayerId": playerId, "OperateId": operateId}).One(&attentionData)
		if err != nil && err != mgo.ErrNotFound {
			easygo.PanicError(err)
			//return false
		}
		if attentionData.GetPlayerId() <= 0 {
			logs.Error("redis 中找不到这个关注信息,数据库也找不到,operateId: %s,playerId: %s", operateId, playerId)
			//return false
		}
		// 设置进redis
		b, _ := json.Marshal(attentionData)
		err1 := easygo.RedisMgr.GetC().HSet(MakeRedisKey(REDIS_SQUARE_ATTENTION, playerId), easygo.AnytoA(attentionData.GetPlayerId()), string(b))
		easygo.PanicError(err1)
		id = attentionData.GetPlayerId()
	}
	b, err1 := easygo.RedisMgr.GetC().Hdel(MakeRedisKey(REDIS_SQUARE_ATTENTION, playerId), easygo.AnytoA(id))
	easygo.PanicError(err1)
	if !b {
		logs.Error(fmt.Sprintf("取消关注失败:%d,%d", playerId, operateId))
		//return false
	}
	AddRedisDynamicDelAttId(id)
	AddPlayerUnreadInfo(playerId, UNREAD_ATTENTION, -1)
	return true
}

// =======================================================================

func GetRedisSquareAttentionInfo(playerId, minId int64) *client_square.AttentionList {
	lst := GetRedisSquareAllAttentionInfo(playerId, minId)
	for _, info := range lst {
		pid := info.GetOperateId()
		player := GetRedisPlayerBase(pid)
		info.HeadIcon = easygo.NewString(player.GetHeadIcon())
		info.Name = easygo.NewString(player.GetNickName())
		info.Sex = easygo.NewInt32(player.GetSex())
		info.Types = easygo.NewInt32(player.GetTypes())
	}
	msg := &client_square.AttentionList{
		AttentionData: lst,
	}
	return msg
}

func GetRedisSquareAllAttentionInfo(playerId, minId int64) []*share_message.AttentionData {
	value, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(REDIS_SQUARE_ATTENTION, playerId)))
	easygo.PanicError(err)
	var lst []*share_message.AttentionData
	for _, s := range value {
		var msg *share_message.AttentionData
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			logs.Error(err)
			continue
		}
		if msg.GetLogId() >= minId && minId != 0 {
			continue
		}

		lst = append(lst, msg)
	}

	sort.Slice(lst, func(i, j int) bool {
		return lst[i].GetLogId() > lst[j].GetLogId() // 降序
	})
	if len(lst) > DYNAMIC_REQUEST_NUM {
		lst = lst[:DYNAMIC_REQUEST_NUM]
	}
	return lst
}

// 没有限制条数的关注列表
func GetRedisSquareAllAttentionInfoNoLimit(playerId int64) []*share_message.AttentionData {
	value, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(MakeRedisKey(REDIS_SQUARE_ATTENTION, playerId)))
	easygo.PanicError(err)
	var lst []*share_message.AttentionData
	for _, s := range value {
		var msg *share_message.AttentionData
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			logs.Error(err)
			continue
		}
		lst = append(lst, msg)
	}
	return lst
}

//==============================savedatabase==============================

func UpdateDynamicSaveAttId(maxId int64) {
	err := easygo.RedisMgr.GetC().HSet(MakeSquareInfoKey(REDIS_SQUARE_SAVEATTID), REDIS_SQUARE_SAVEATTID, maxId)
	easygo.PanicError(err)
}

func GetDynamicSaveAttId() int64 {
	b, err := easygo.RedisMgr.GetC().HExists(MakeSquareInfoKey(REDIS_SQUARE_SAVEATTID), REDIS_SQUARE_SAVEATTID)
	easygo.PanicError(err)
	if !b {
		return 0
	}
	Id, err1 := easygo.RedisMgr.GetC().HGet(MakeSquareInfoKey(REDIS_SQUARE_SAVEATTID), REDIS_SQUARE_SAVEATTID)
	easygo.PanicError(err1)
	return easygo.AtoInt64(string(Id))
}
