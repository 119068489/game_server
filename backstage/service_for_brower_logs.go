// 管理后台为[浏览器]提供的服务
//用户管理

package backstage

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/mongo_init"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/share_message"
)

// 获取管理员日志
func (self *cls4) RpcQueryManagerLog(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	// if reqMsg.GetListType() == 2 {
	// 	_, err := strconv.Atoi(reqMsg.GetKeyword())
	// 	easygo.PanicError(err, "被操作人查询内容只能为数值")
	// }
	list, count := QueryBackstageLog(reqMsg)
	msg := &brower_backstage.ManagerLogResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//查询管理员日志分类下拉列表
func (self *cls4) RpcManagerLogTypesKeyValue(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	list := mongo_init.InitManageLogTypes()
	return &brower_backstage.KeyValueResponse{
		List: list,
	}
}

//查询现金变化类型
func (self *cls4) RpcQueryGoldLog(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryGoldLogRequest) easygo.IMessage {
	list, count := GetGoldLogList(reqMsg)

	var playerids []int64
	for _, item := range list {
		if for_game.IsContains(item.PlayerId, playerids) == -1 {
			playerids = append(playerids, item.PlayerId)
		}
	}

	players := for_game.GetAllPlayerBase(playerids)

	reqType := &brower_backstage.SourceTypeRequest{}
	soucetype := QuerySouceTypeList(reqType)
	sts := make(map[int32]*share_message.SourceType)
	for _, s := range soucetype {
		sts[s.GetKey()] = s
	}

	var msg []*brower_backstage.GoldLogList
	for _, l := range list {
		one := &brower_backstage.GoldLogList{
			InLine: &share_message.GoldChangeLog{
				LogId:          easygo.NewInt64(l.LogId),
				PlayerId:       easygo.NewInt64(l.PlayerId),
				Account:        easygo.NewString(players[l.PlayerId].GetAccount()),
				ChangeGold:     easygo.NewInt64(l.ChangeGold),
				PayType:        easygo.NewInt32(l.PayType),
				SourceType:     easygo.NewInt32(l.SourceType),
				SourceTypeName: easygo.NewString(sts[l.SourceType].GetValue()),
				CurGold:        easygo.NewInt64(l.CurGold),
				Gold:           easygo.NewInt64(l.Gold),
				Note:           easygo.NewString(l.Note),
				CreateTime:     easygo.NewInt64(l.CreateTime),
			},
			Extend: l.Extend,
		}
		msg = append(msg, one)
	}

	return &brower_backstage.QueryGoldLogResponse{
		List:      msg,
		PageCount: easygo.NewInt32(count),
	}
}
