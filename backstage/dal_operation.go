package backstage

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo/bson"
)

//查询运营渠道列表
func GetOperationList(reqMsg *brower_backstage.ListRequest) ([]*share_message.OperationChannel, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_OPERATION_CHANNEL)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.BeginTimestamp != nil {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Name"] = reqMsg.GetKeyword()
		case 2:
			queryBson["ChannelNo"] = reqMsg.GetKeyword()
		case 3:
			queryBson["CompanyName"] = reqMsg.GetKeyword()
		case 4:
			queryBson["Type"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		queryBson["Cooperation"] = reqMsg.GetListType()
	}

	if reqMsg.DownType != nil && reqMsg.GetDownType() == 0 {
		queryBson["Status"] = 1
	} else {
		queryBson["Status"] = reqMsg.GetDownType()
	}

	var list []*share_message.OperationChannel
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改运营渠道
func EditOperationChannel(reqMsg *share_message.OperationChannel) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_OPERATION_CHANNEL)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//修改运营渠道执行表
func EditOperationChannelUse(reqMsg *share_message.OperationChannel) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_OPERATION_CHANNEL_USE)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//查询所有渠道数据
func GetOperationListAll() []*share_message.OperationChannel {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_OPERATION_CHANNEL)
	defer closeFun()

	queryBson := bson.M{}
	var list []*share_message.OperationChannel
	query := col.Find(queryBson)
	errc := query.All(&list)
	easygo.PanicError(errc)

	return list
}

//同步渠道表修改
func SynchronizeOperationChannel() {
	lis := GetOperationListAll()
	for _, i := range lis {
		EditOperationChannelUse(i)
	}
}
