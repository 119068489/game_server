package sport_common_dal

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

func GetESPortsGameItemViewList(cPage, pSize int32, lableId int64) ([]*share_message.TableESPortsGame, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESPortsGame
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME)
	defer closeFun()
	queryBson := bson.M{}

	queryBson["app_label_id"] = lableId
	// //比赛状态 0 未开始，1 进行中，2 已结束(api字段)(0和1的时候结合begin_time判断)
	queryBson["game_status"] = bson.M{"$ne": for_game.GAME_STATUS_2}
	queryBson["release_flag"] = for_game.GAME_RELEASE_FLAG_2 //2==已发布

	query := col.Find(queryBson)
	count, err := query.Count()
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	err = query.Select(bson.M{"match_name": 1, "match_stage": 1, "bo": 1, "team_a": 1, "team_b": 1}).Sort("begin_time_int").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if err != nil {
		logs.Error(err)
		return nil, 0
	}

	return list, count
}
func GetESPortsGameItem(upgameid int64) *share_message.TableESPortsGame {

	rd := &share_message.TableESPortsGame{}
	col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["_id"] = upgameid
	err := col.Find(queryBson).One(&rd)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return rd
}
