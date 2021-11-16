package sport_common_dal

import (
	"game_server/easygo"
	"game_server/for_game"
	//"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

// 递增 id 表
type RerurnVal struct {
	Key   *int64 `bson:"_id"`
	Value *int32 `bson:"ThumbsUpCount,omitempty"`
}

func DataOptThumbsUp(table string, preId, exId int64, playerId int64, meunid int32) (int32, string, int32, int32) {

	dataId := exId
	parentId := int64(0)
	if preId > 0 {
		parentId = preId
	}

	count := int32(0)
	isthup := int32(0)
	redisObj, b := IsThumbsUp(meunid, parentId, dataId, playerId)
	if b {
		count = -1
	} else {
		count = 1
		isthup = 1
	}
	col, closeFun := GetC(table)
	defer closeFun()

	//data.ThumbsUpCount = easygo.NewInt32(data.GetThumbsUpCount() + count)
	//err := col.Update(bson.M{"_id": dataId}, bson.M{"$inc": bson.M{"ThumbsUpCount": count}})
	rerurnVal := &RerurnVal{
		Key:   &dataId,
		Value: easygo.NewInt32(count),
	}
	_, err := col.Find(bson.M{"_id": dataId}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"ThumbsUpCount": count}},
		Upsert:    true,
		ReturnNew: true,
	}, &rerurnVal)
	allcount := int32(*rerurnVal.Value)
	if err == nil {
		if count > 0 {
			redisObj.AddThumbsUp(dataId)
		} else {
			redisObj.CancelThumbsUp(dataId)
		}
		return for_game.C_OPT_SUCCESS, "操作成功", isthup, allcount
	} else {
		return for_game.C_INFO_NOT_EXISTS, "系统错误,数据丢失", 0, 0
	}
}

func IsThumbsUp(menuId int32, dataId, exId int64, playerId int64) (*for_game.RedisESportThumbsUpObj, bool) {
	obj := for_game.GetRedisESportThumbsUpObj(menuId, dataId, playerId)
	return obj, obj.IsThumbsUp(exId)
}

///redis某个int32的字段 +1 -1
func UpdateFedAddition_xv(table string, fed string, dataId int64, count int64) (int32, string) {
	//obj := for_game.GetRedisESportTableFedUpdateObj(table)
	//obj.RvpValue(fed, dataId, count)
	return UpdateFedAddition(table, fed, dataId, count)
}

///某个int32的字段 +1 -1
func UpdateFedAddition(table string, fed string, dataId int64, count int64) (int32, string) {
	col, closeFun := GetC(table)
	defer closeFun()
	//data := bson.M{}
	/*err := col.Find(bson.M{"_id": dataId}).Select(bson.M{fed: 1}).One(&data)
	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	if err == mgo.ErrNotFound {
		return for_game.C_INFO_NOT_EXISTS, "数据不存在"
	}*/
	//err := col.Update(bson.M{"_id": dataId}, bson.M{"$inc": bson.M{fed: count}})
	_, err := col.Find(bson.M{"_id": dataId}).Apply(mgo.Change{
		Update: bson.M{"$inc": bson.M{fed: count}},
		Upsert: true,
	}, nil)

	if err == nil {
		return for_game.C_OPT_SUCCESS, "操作成功"
	} else {
		return for_game.C_INFO_NOT_EXISTS, "系统错误,数据丢失"
	}
}
func UpdateFedAdditionEx(col *mgo.Collection, fed string, dataId int64, count int64) (int32, string) {

	/*
		data := bson.M{}
		err := col.Find(bson.M{"_id": dataId}).Select(bson.M{fed: 1}).One(&data)
		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
			return for_game.C_SYS_ERROR, "系统错误"
		}
		if err == mgo.ErrNotFound {
			return for_game.C_INFO_NOT_EXISTS, "数据不存在"
		}
	*/
	logs.Info("UpdateFedAdditionEx +", dataId, "+", count)
	//err = col.Update(bson.M{"_id": dataId}, bson.M{"$inc": bson.M{fed: count}})
	_, err := col.Find(bson.M{"_id": dataId}).Apply(mgo.Change{
		Update: bson.M{"$inc": bson.M{fed: count}},
		Upsert: true,
	}, nil)

	//logs.Info("UpdateFedAdditionEx rerurnVal", val)
	if err == nil {
		return for_game.C_OPT_SUCCESS, "操作成功"
	} else {
		logs.Error(err)
		return for_game.C_INFO_NOT_EXISTS, "系统错误,数据丢失"
	}

}
