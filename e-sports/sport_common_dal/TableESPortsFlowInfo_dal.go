package sport_common_dal

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"time"
	//   "time"
)

func AddTableESPortsFlowInfo(dataTpye int64, playerId, dataId int64) (int32, string) {

	table := ""
	if dataTpye == for_game.ESPORT_FLOW_LIVE_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_LIVE_HISTORY
	} else if dataTpye == for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_LIVE_FOLLOW_HISTORY
	} else if dataTpye == for_game.ESPORT_FLOW_VIDEO_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_VIDEO_HISTORY
	}
	if table == "" {
		return for_game.C_INFO_NOT_EXISTS, "数据类型不存在"
	}
	info := &share_message.TableESPortsFlowInfo{
		PlayerId:   easygo.NewInt64(playerId),
		DataId:     easygo.NewInt64(dataId),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
	}
	col, closeFun := GetC(table)
	defer closeFun()

	//DeletePalyerTableESPortsFlowInfo(col, dataId, playerId)
	return CreateTableESPortsFlowInfoEx(col, table, info)
}

func FlowInfoIdKey(PlayerId, dataId int64) string {
	return fmt.Sprintf("%d_%d", PlayerId, dataId)
}

func CreateTableESPortsFlowInfoEx(col *mgo.Collection, table string, info *share_message.TableESPortsFlowInfo) (int32, string) {
	//id := NextId(table)
	info.CreateTime = easygo.NewInt64(time.Now().Unix())
	info.Id = easygo.NewString(FlowInfoIdKey(info.GetPlayerId(), info.GetDataId()))
	_, err := col.Upsert(bson.M{"_id": info.GetId()}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}

func DeletePalyerTableESPortsFlowInfo(col *mgo.Collection, dataid, playerId int64) bool {

	_, err := col.RemoveAll(bson.M{"DataId": dataid, "PlayerId": playerId})
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}

//添加關注
func AddTableESPortsFlowfollowInfo(dataTpye int64, playerId, dataId int64) (int32, string) {

	table := ""
	if dataTpye == for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_LIVE_FOLLOW_HISTORY
	}
	if table == "" {
		return for_game.C_INFO_NOT_EXISTS, "数据类型不存在"
	}
	info := &share_message.TableESPortsFlowInfo{
		PlayerId:   easygo.NewInt64(playerId),
		DataId:     easygo.NewInt64(dataId),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
	}
	col, closeFun := GetC(table)
	defer closeFun()

	return CreateTableESPortsFlowInfoEx(col, table, info)

}

func DeleteTableESPortsFlowInfo(dataTpye int32, id int64) bool {

	table := ""
	if dataTpye == for_game.ESPORT_FLOW_LIVE_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_LIVE_HISTORY
	} else if dataTpye == for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_LIVE_FOLLOW_HISTORY
	} else if dataTpye == for_game.ESPORT_FLOW_VIDEO_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_VIDEO_HISTORY
	}

	col, closeFun := GetC(table)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": id})
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}

func DeleteTableESPortsFlowInfoEx(dataTpye int64, plyid, id int64) bool {

	table := ""
	if dataTpye == for_game.ESPORT_FLOW_LIVE_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_LIVE_HISTORY
	} else if dataTpye == for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_LIVE_FOLLOW_HISTORY
	} else if dataTpye == for_game.ESPORT_FLOW_VIDEO_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_VIDEO_HISTORY
	}

	col, closeFun := GetC(table)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"DataId": id, "PlayerId": plyid})
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}

///获取历史或者关注视频或放映厅
func GetFlowVideoOrLiveListByPlayerId(dataTpye int64, cPage, pSize int32, sort string, playerId int64) ([]*share_message.TableESPortsVideoInfo, int) {
	table := ""
	/**/
	if dataTpye == for_game.ESPORT_FLOW_LIVE_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_LIVE_HISTORY
	} else if dataTpye == for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_LIVE_FOLLOW_HISTORY
	} else if dataTpye == for_game.ESPORT_FLOW_VIDEO_HISTORY {
		table = for_game.TABLE_ESPORTS_FLOW_VIDEO_HISTORY
	}
	if dataTpye == for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY { //关注放服务器 ，其他app本地存储
		table = for_game.TABLE_ESPORTS_FLOW_LIVE_FOLLOW_HISTORY
	}
	if table == "" {
		return nil, 0
	}

	return GetTableESPortsFlowLiveFollow_t(table, for_game.TABLE_ESPORTS_VIDEO, cPage, pSize, sort, playerId)
}

/*
//获取放映厅关注列表
func GetTableESPortsFlowLiveFollow(maintable, othetable string, cPage, pSize int32, sort string, playerId int64) ([]*share_message.TableESPortsVideoInfo, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESPortsVideoInfo
	col, closeFun := GetC(maintable)
	defer closeFun()
	pipe := col.Pipe([]bson.M{
		{"$match": bson.M{"PlayerId": playerId}},
		{
			"$lookup": bson.M{
				"from":         othetable,
				"localField":   "DataId",
				"foreignField": "_id",
				"as":           "union",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": bson.M{
					"$mergeObjects": []interface{}{
						bson.M{"$arrayElemAt": []interface{}{"$union", 0}},
						"$$ROOT",
					},
				},
			},
		},
		{"$project": bson.M{
			"union": 0,
			"_id":   0,
		}},
		bson.M{"$sort": bson.M{sort: -1}}, // 1 升序 -1降序
		bson.M{"$skip": curPage * int(pageSize)},
		bson.M{"$limit": pageSize},
	})
	//var data []interface{}
	err := pipe.All(&list)
	if list != nil {
		for _, v := range list {
			logs.Info(v)
		}
	}

	if err != nil {
		logs.Error(err)
		return nil, 0
	}

	return list, 0
}
*/
//获取放映厅关注列表
func GetTableESPortsFlowLiveFollow_t(maintable, othetable string, cPage, pSize int32, sort string, playerId int64) ([]*share_message.TableESPortsVideoInfo, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESPortsVideoInfo
	col, closeFun := GetC(maintable)
	defer closeFun()

	var getfed = func(name string) string {
		return "$" + othetable + "." + name
	}
	pipe := col.Pipe([]bson.M{
		{"$match": bson.M{"PlayerId": playerId}},
		{
			"$lookup": bson.M{
				"from":         othetable,
				"localField":   "DataId",
				"foreignField": "_id",
				"as":           othetable,
			},
		},
		{
			"$unwind": bson.M{
				"path": "$" + othetable,
			},
		},
		{"$project": bson.M{
			"_id":              getfed("_id"),
			"DataId":           1,
			"CreateTime":       1,
			"VideoUrl":         getfed("VideoUrl"),      //"$" + othetable + ".VideoUrl",
			"CoverImageUrl":    getfed("CoverImageUrl"), //""$" + othetable + ".CoverImageUrl",
			"Status":           getfed("Status"),        //""$" + othetable + ".Status",
			"MatchName":        getfed("MatchName"),     //" "$" + othetable + ".MatchName",
			"ThumbsUpCount":    getfed("ThumbsUpCount"),
			"ThumbsUpCountSys": getfed("ThumbsUpCountSys"),
			"AppLabelName":     getfed("AppLabelName"),
			"AppLabelID":       getfed("AppLabelID"),
			"FlowCount":        getfed("FlowCount"),
			"LookCount":        getfed("LookCount"),
			"LookCountSys":     getfed("LookCountSys"),
			"Title":            getfed("Title"),
			"AuthorPlayerId":   getfed("AuthorPlayerId"),
			"UniqueGameId":     getfed("UniqueGameId"),
			"UniqueGameName":   getfed("UniqueGameName"),
		}},
		{"$match": bson.M{"Status": 1}},
		bson.M{"$sort": bson.M{sort: -1}}, // 1 升序 -1降序
		bson.M{"$skip": curPage * int(pageSize)},
		bson.M{"$limit": pageSize},
	})
	count, err := col.Find(bson.M{"PlayerId": playerId}).Count()
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	//var data []interface{}
	err = pipe.All(&list)

	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	/*if list != nil {
		logs.Info(len(list))
		for _, v := range list {
			//logs.Info("%d,%d,%s,%d", v.GetId(), v.GetDataId(), v.GetCoverImageUrl(), v.GetCreateTime())
			logs.Info(v)
		}
	}*/
	return list, count
}
