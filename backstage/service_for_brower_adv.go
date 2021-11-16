// 管理后台为[浏览器]提供的服务
//广告管理

package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/share_message"

	"github.com/astaxie/beego/logs"

	"github.com/akqp2019/mgo/bson"
)

//查询广告列表
func (self *cls4) RpcQueryAdvList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	UpdateAdvListToDown()

	list, count := QueryAdvListToDB(reqMsg)

	msg := &brower_backstage.AdvListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return msg
}

//修改广告
func (self *cls4) RpcEditAdvData(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.AdvSetting) easygo.IMessage {
	// if reqMsg.GetId() == 0 && for_game.QueryAdvByLacSort(reqMsg.GetLocation(), reqMsg.GetWeights()) != nil {
	// 	return easygo.NewFailMsg("排序权重已存在")
	// }
	reqMsg.CreateTime = easygo.NewInt64(util.GetMilliTime())
	msg := fmt.Sprintf("修改广告:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ADV_DATA))
		msg = fmt.Sprintf("添加广告:%d", reqMsg.GetId())
	}
	nowTime := util.GetMilliTime()
	if reqMsg.GetStatus() == for_game.ADV_ON_SHELF && reqMsg.GetEndTime() < nowTime {
		return easygo.NewFailMsg("投放时间已过期")
	}

	if reqMsg.GetIsShield() && reqMsg.GetLocation() != 3 {
		return easygo.NewFailMsg("参数设置错误")
	}

	if reqMsg.GetLocation() == for_game.ADV_LOCATION_BANNER_PERS && reqMsg.GetStatus() == for_game.ADV_ON_SHELF {
		onCount := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_ADV_DATA, bson.M{"Location": for_game.ADV_LOCATION_BANNER_PERS, "Status": for_game.ADV_ON_SHELF})
		if onCount >= 5 {
			return easygo.NewFailMsg("个人信息页banner位,只允许上架5个广告")
		}
	}

	if reqMsg.GetIsTop() {
		old := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ADV_DATA, bson.M{"IsTop": true})
		if old != nil {
			one := &share_message.AdvSetting{}
			for_game.StructToOtherStruct(old, one)
			if one.GetId() != reqMsg.GetId() {
				one.IsTop = easygo.NewBool(false)
				for_game.UpdateAdvListToDB(one)
			}
		}
	}

	for_game.UpdateAdvListToDB(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ADV_MANAGE, msg)

	return easygo.EmptyMsg
}

//批量上架广告
func (self *cls4) RpcAdvOnShelf(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	var ids string
	idsarr := reqMsg.GetIds64()
	for_game.UpdateAdvShelfToDB(idsarr, for_game.ADV_ON_SHELF)
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}

	msg := fmt.Sprintf("批量上架广告: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ADV_MANAGE, msg)

	return easygo.EmptyMsg
}

//批量下架广告
func (self *cls4) RpcAdvOffShelf(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	var ids string
	idsarr := reqMsg.GetIds64()
	for_game.UpdateAdvShelfToDB(idsarr, for_game.ADV_OFF_SHELF)
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}

	msg := fmt.Sprintf("批量下架广告: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ADV_MANAGE, msg)

	return easygo.EmptyMsg
}

//修改广告排序
func (self *cls4) RpcUpdateAdvSort(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	logs.Debug("RpcUpdateAdvSort:", reqMsg)
	ids := reqMsg.GetIds64()
	if len(ids) == 0 {
		return easygo.NewFailMsg("ids不能为空")
	}

	for_game.UpdateAdvSort(ids)
	msg := fmt.Sprintf("修改广告排序")
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ADV_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除广告
func (self *cls4) RpcDelAdvData(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	logs.Debug("RpcDelAdvData:", reqMsg)
	aids := reqMsg.GetIds64()
	if len(aids) == 0 {
		return easygo.NewFailMsg("ids不能为空")
	}

	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_ADV_DATA, bson.M{"_id": bson.M{"$in": aids}, "Status": for_game.ADV_OFF_SHELF})

	count := len(aids)
	var ids string
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(aids[i])) + ","

		} else {
			ids += easygo.AnytoA(aids[i])
		}
	}

	msg := fmt.Sprintf("批量删除广告: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ADV_MANAGE, msg)

	return easygo.EmptyMsg
}

//从数据库查询广告列表
func QueryAdvListToDB(reqMsg *brower_backstage.ListRequest) ([]*share_message.AdvSetting, int32) {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ADV_DATA)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1: //id查询
			id := easygo.StringToInt64noErr(reqMsg.GetKeyword())
			queryBson["_id"] = id
		case 2: //标题
			queryBson["Title"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		queryBson["Location"] = reqMsg.GetListType()
	}

	var list []*share_message.AdvSetting
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-Weights").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

//查询修改过期的广告为下架状态
func UpdateAdvListToDown() {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ADV_DATA)
	defer closeFun()

	nowTime := util.GetMilliTime()
	queryBson := bson.M{"Status": for_game.ADV_ON_SHELF, "EndTime": bson.M{"$lt": nowTime}}

	upBson := bson.M{"$set": bson.M{"Status": for_game.ADV_OFF_SHELF}}
	_, err := col.UpdateAll(queryBson, upBson)
	easygo.PanicError(err)
}

//匹配首页Tips配置列表
func (self *cls4) RpcIndexTipsList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	sort := []string{}
	if len(sort) == 0 {
		sort = append(sort, "_id")
	}
	findBson := bson.M{}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_INDEX_TIPS, findBson, 0, 0, sort...)
	var list []*share_message.IndexTips
	for _, li := range lis {
		one := &share_message.IndexTips{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}

	return &brower_backstage.IndexTipsResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//保存首页Tips配置
func (l *cls4) RpcSaveIndexTips(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.IndexTips) easygo.IMessage {
	if for_game.GetAllAdvsByIds([]int64{reqMsg.GetAdvId()}) == nil {
		s := fmt.Sprintf("广告配置错误:可能广告已过期")
		return easygo.NewFailMsg(s)
	}
	msg := fmt.Sprintf("修改首页Tips配置:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_INDEX_TIPS))
		msg = fmt.Sprintf("添加首页Tips配置:%d", reqMsg.GetId())
		reqMsg.Types = easygo.NewInt32(2)
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_INDEX_TIPS, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ADV_MANAGE, msg)

	return easygo.EmptyMsg
}

//批量删除首页Tips配置
func (l *cls4) RpcDelIndexTips(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds32()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}
	for _, id := range idList {
		switch id {
		case 1, 2, 3, 4:
			return easygo.NewFailMsg("1,2,3,4不能删除")
		}
	}

	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_INDEX_TIPS, bson.M{"_id": bson.M{"$in": idList}})

	count := len(idList)
	var ids string
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idList[i])) + ","
		} else {
			ids += easygo.IntToString(int(idList[i]))
		}
	}
	msg := fmt.Sprintf("批量删除首页Tips配置: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ADV_MANAGE, msg)
	return easygo.EmptyMsg
}

//弹窗悬浮球配置列表
func (self *cls4) RpcPopSuspendList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	sort := []string{}
	if len(sort) == 0 {
		sort = append(sort, "_id")
	}
	findBson := bson.M{}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_POP_SUSPEND, findBson, 0, 0, sort...)
	var list []*share_message.PopSuspend
	for _, li := range lis {
		one := &share_message.PopSuspend{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}

	return &brower_backstage.PopSuspendResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//批量保存弹窗悬浮球配置
func (l *cls4) RpcSavePopSuspendList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.PopSuspendResponse) easygo.IMessage {
	msg := "修改弹窗悬浮球配置"
	if len(reqMsg.GetList()) == 0 {
		return easygo.NewFailMsg("请提交要修改的项")
	}

	var advids []int32
	var saveData []interface{}
	for _, v := range reqMsg.GetList() {
		if for_game.GetAllAdvsByIds([]int64{v.GetAdvId()}) == nil {
			advids = append(advids, v.GetId())
		}
		saveData = append(saveData, bson.M{"_id": v.GetId()}, v)
	}

	count := len(advids)
	if count > 0 {
		ids := ""
		for i, j := range advids {
			ids += easygo.AnytoA(j)
			if i < count-1 {
				ids += ","
			}
		}
		s := fmt.Sprintf("以下位置广告配置错误:%s", ids)
		return easygo.NewFailMsg(s)
	}
	// 批量修改时间
	if len(saveData) > 0 {
		for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_POP_SUSPEND, saveData)
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ADV_MANAGE, msg)

	return easygo.EmptyMsg
}

//广告列表下拉
func (l *cls4) RpcQueryAdvDownList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	sort := []string{}
	if len(sort) == 0 {
		sort = append(sort, "_id")
	}
	findBson := bson.M{}
	if reqMsg.Type != nil {
		findBson["Location"] = reqMsg.GetType()
	}
	if reqMsg.Status != nil {
		findBson["Status"] = reqMsg.GetStatus()
	}

	lis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ADV_DATA, findBson, 0, 0, sort...)
	var list []*brower_backstage.KeyValueTag
	for _, li := range lis {
		one := &share_message.AdvSetting{}
		for_game.StructToOtherStruct(li, one)
		li := &brower_backstage.KeyValueTag{
			Key:   easygo.NewInt32(one.GetId()),
			Value: easygo.NewString((easygo.AnytoA(one.GetId())) + one.GetTitle()),
		}
		list = append(list, li)
	}

	return &brower_backstage.KeyValueResponseTag{
		List: list,
	}
}
