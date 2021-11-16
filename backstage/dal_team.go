package backstage

import (
	"encoding/base64"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//查询群列表
func GetTeamList(user *share_message.Manager, reqMsg *brower_backstage.GetTeamListRequest) ([]*share_message.TeamData, int) {
	var list []*share_message.TeamData
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)
	permission := true
	if user.GetRole() != 0 {
		permission = !QueryPermissionById(user.GetSite(), user.GetRoleType(), "groupManage-oneself")
	}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	defer closeFun()

	// 关键词查询账号
	queryBson := bson.M{}

	if !permission {
		queryBson["AdminID"] = user.GetId()
	}
	if reqMsg.GetKeyword() != "" && reqMsg.Keyword != nil {
		switch reqMsg.GetType() {
		case 1:
			queryBson["TeamChat"] = easygo.If(reqMsg.GetKeyword() != "", reqMsg.GetKeyword(), bson.M{"$ne": nil})
		case 2:
			queryBson["Name"] = easygo.If(reqMsg.GetKeyword() != "", reqMsg.GetKeyword(), bson.M{"$ne": nil})
		case 3:
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			queryBson["Owner"] = player.GetPlayerId()
		}
	}

	if reqMsg.State != nil {
		queryBson["Status"] = reqMsg.GetState()
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		switch reqMsg.GetListType() {
		case 1:
			queryBson["IsRecommend"] = true
		case 2:
			queryBson["IsRecommend"] = false
		}
	}

	// // 判断有日期才按日期查询
	// if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
	// 	queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	// }

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//Id查询群
func QueryTeambyId(id TEAM_ID) *share_message.TeamData {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	defer closeFun()

	siteOne := &share_message.TeamData{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//Ids查询群
func QueryTeambyIds(ids []int64) []*share_message.TeamData {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	defer closeFun()
	queryBson := bson.M{"_id": bson.M{"$in": ids}}
	var list []*share_message.TeamData
	query := col.Find(queryBson)
	errc := query.Sort("-_id").All(&list)
	easygo.PanicError(errc)
	return list
}

//Id查询群
func QueryTeambyNick(nickName string) *share_message.TeamData {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	defer closeFun()

	siteOne := &share_message.TeamData{}
	err := col.Find(bson.M{"TeamChat": nickName}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//Id查询群成员详情
func QueryTeamMember(reqMsg *brower_backstage.TeamMemberRequest) ([]*share_message.PersonalTeamData, int) {
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAMMEMBER)
	defer closeFun()

	list := []*share_message.PersonalTeamData{}
	queryBson := bson.M{"TeamId": reqMsg.GetTeamId()}
	//
	if reqMsg.GetState() != 0 && reqMsg.State != nil {
		queryBson["Status"] = reqMsg.GetState()
	}

	// 关键词查询
	if reqMsg.Type != nil && reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		flag := reqMsg.GetType()
		switch flag {
		case 1: //柠檬号
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			if player != nil {
				queryBson["PlayerId"] = player.GetPlayerId()
			}
		case 2: //群昵称
			var playerids PLAYER_IDS
			plis := GetPlayerLikeNickname(for_game.MONGODB_NINGMENG, reqMsg.GetKeyword())
			for _, pli := range plis {
				playerids = append(playerids, pli.GetPlayerId())
			}
			queryBson["PlayerId"] = bson.M{"$in": playerids}
			// queryBson["NickName"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "i"}}
		}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("—_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	var playerids []int64
	for _, item := range list {
		playerids = append(playerids, item.GetPlayerId())
	}

	players := for_game.GetAllPlayerBase(playerids)

	for _, item := range list {
		info := make([]*share_message.OperatorInfoPer, 0)
		for _, value := range item.GetOperatorInfoPer() {
			if value.GetCloseTime() <= util.GetMilliTime() {
				// DelTeanMemInfoByIds(item.GetTeamId(), []int64{item.GetPlayerId()}, 1)
			} else {
				info = append(info, value)
			}
		}
		item.OperatorInfoPer = info
		p := players[item.GetPlayerId()]
		item.Account = easygo.NewString(p.GetAccount())
		item.PerNickName = easygo.NewString(p.GetNickName())
		if item.TeamChannel == nil {
			teamChannel := &share_message.TeamChannel{
				Name: easygo.NewString("客服"),
				Type: easygo.NewInt32(5),
			}
			item.TeamChannel = teamChannel
		}
	}

	return list, count
}

//查询群聊天记录列表
func QueryChatRecord(reqMsg *brower_backstage.ListRequest) ([]*share_message.TeamChatLog, int) {
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)
	teamId := reqMsg.GetId()
	// tableName := for_game.GetMongoTableName(teamId, for_game.TABLE_TEAM_CHAT_LOG) //每个群聊有自己的日志id
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_TEAM_CHAT_LOG)
	defer closeFun()

	queryBson := bson.M{"TeamId": teamId}
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	var list []*share_message.TeamChatLog
	errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	//解析聊天内容
	if reqMsg.GetListType() == 1 {
		for i, item := range list {
			if item.GetContent() != "" {
				content, _ := base64.StdEncoding.DecodeString(item.GetContent())
				list[i].Content = easygo.NewString(string(content))
				team := QueryTeambyId(item.GetTeamId())
				list[i].TeamAccount = easygo.NewString(team.GetTeamChat())
				list[i].TeamName = easygo.NewString(team.Name)
			}
			if item.GetTalker() != 0 {
				player := QueryPlayerbyId(item.GetTalker())
				list[i].TalkerAccount = easygo.NewString(player.GetAccount())
				list[i].TalkerName = easygo.NewString(player.GetNickName())
			}

		}
	}

	return list, count
}

//查询群历史列表
func QueryTeamMessage(reqMsg *brower_backstage.ListRequest) ([]*share_message.TeamChatLog, int) {
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)
	teamId := easygo.StringToIntnoErr(reqMsg.GetKeyword())
	// tableName := for_game.GetMongoTableName(teamId, for_game.TABLE_TEAM_CHAT_LOG) //每个群聊有自己的日志id
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_TEAM_CHAT_LOG)
	defer closeFun()

	queryBson := bson.M{"$or": []bson.M{
		{"TeamMessage.Type": for_game.INVITE_PLAYER}, {"TeamMessage.Type": for_game.DEL_PLAYER},
		{"TeamMessage.Type": for_game.ADD_MANAGER}, {"TeamMessage.Type": for_game.DEL_MANAGER},
		{"TeamMessage.Type": for_game.CHANGE_OWNER}, {"TeamMessage.Type": for_game.TEAM_GONGGAO},
		{"TeamMessage.Type": for_game.BACKSTAGE_BAN_TEAM}, {"TeamMessage.Type": for_game.STOP_TALK},
	}}

	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["Time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	queryBson["TeamId"] = teamId
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	var list []*share_message.TeamChatLog = []*share_message.TeamChatLog{}

	team := QueryTeambyId(int64(teamId)) //拿到群的信息
	if team == nil {
		logs.Error("群[%d]找不到", teamId)
		return list, 0
	}

	if reqMsg.BeginTimestamp == nil || reqMsg.GetBeginTimestamp() == 0 || (reqMsg.GetBeginTimestamp() >= team.GetCreateTime() && team.GetCreateTime() <= reqMsg.GetEndTimestamp()) {
		if (curPage+1)*pageSize > count {
			errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
			easygo.PanicError(errc)
			player := QueryPlayerbyId(team.GetOwner()) //找群主

			if player == nil {
				return list, 0
			}
			AdminName := ""
			if team.GetAdminID() != 0 {
				admin := QueryManageByID(team.GetAdminID()) //找后台管理员
				if admin != nil {
					AdminName = admin.GetRealName()
				}
			}

			list = append(list, &share_message.TeamChatLog{Time: easygo.NewInt64(team.GetCreateTime()), Content: easygo.NewString(fmt.Sprintf("%s创建了群", player.GetNickName())), TeamMessage: &share_message.TeamMessage{
				PlayerId:  team.Owner,
				Type:      easygo.NewInt32(99),
				Pos:       easygo.NewInt32(1),
				Name:      player.NickName,
				Account:   player.Account,
				TeamId:    team.Id,
				AdminID:   team.AdminID,
				AdminName: &AdminName,
			}})
		} else {
			errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
			easygo.PanicError(errc)
		}
		count = count + 1
	} else {
		errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
		easygo.PanicError(errc)
	}

	for _, item := range list {
		AdminName := ""
		if item.TeamMessage.GetAdminID() != 0 {
			admin := QueryManageByID(item.TeamMessage.GetAdminID())
			if admin != nil {
				AdminName = admin.GetRealName()
			}
		}

		item.TeamMessage.AdminName = easygo.NewString(AdminName)
		optName := item.TeamMessage.GetName()
		val, _ := base64.StdEncoding.DecodeString(item.TeamMessage.GetValue1())
		switch item.TeamMessage.GetType() {
		case for_game.INVITE_PLAYER:
			invite_name := ""
			for _, member := range item.TeamMessage.Members {
				invite_name += fmt.Sprintf("%s、", member.GetNickName())
			}
			if len(invite_name) > 0 {
				invite_name = invite_name[:len(invite_name)-3] //为什么减3  因为一个顿号在字符串中占3个索引
			}
			item.Content = easygo.NewString(fmt.Sprintf("%s邀请了%s入群", optName, invite_name))
		case for_game.DEL_PLAYER:
			item.Content = easygo.NewString(fmt.Sprintf("%s将%s移出了群", optName, val))
		case for_game.ADD_MANAGER:
			item.Content = easygo.NewString(fmt.Sprintf("%s将%s设定为了管理员", optName, val))
		case for_game.DEL_MANAGER:
			item.Content = easygo.NewString(fmt.Sprintf("%s移除了%s的管理员身份", optName, val))
		case for_game.CHANGE_OWNER:
			content, _ := base64.StdEncoding.DecodeString(item.GetContent())
			item.Content = easygo.NewString(string(content))
		case for_game.TEAM_GONGGAO:
			item.Content = easygo.NewString(fmt.Sprintf("%s修改了群公告", optName))
		case for_game.STOP_TALK:
			item.Content = easygo.NewString("前端全员禁言")
		}

	}

	return list, count
}

//写后台修改的群设置
func SaveTeamEdit(reqMsg *share_message.TeamData) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	defer closeFun()
	data := bson.M{"MaxMember": reqMsg.GetMaxMember(), "IsRecommend": reqMsg.GetIsRecommend(), "GongGao": reqMsg.GetGongGao(), "Name": reqMsg.GetName(), "Level": reqMsg.GetLevel()}
	err1 := col.Update(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err1)
}

//更新群封禁人信息
func UpdateTeanInfoByIds(ids []int64, reqMsg *share_message.OperatorInfo) {
	for _, id := range ids {
		teamObj := for_game.GetRedisTeamObj(id)
		teamObj.SetTeamIsBan2(id, true)
		teamObj.SetTeamStatus(for_game.BANNED)
		teamObj.SetTeamUnBanTime(id, reqMsg.GetCloseTime())
		teamObj.SetRedisOperatorInfo(reqMsg)
		teamObj.SaveToMongo()
	}
}

func DelTeanInfoByIds(ids []int64, flag int32) {

	for _, id := range ids {
		teamObj := for_game.GetRedisTeamObj(id)
		message := teamObj.GetTeamMessageSetting()
		teamObj.SetTeamIsBan2(id, false)
		if !message.GetIsStopTalk() { //设置状态前判断前端是否封禁
			teamObj.SetTeamStatus(for_game.NORMAL)
		}
		teamObj.DelRedisOperatorInfo(for_game.SYSTEM)
		teamObj.SaveToMongo()
	}

	//col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	//defer closeFun()
	//_,err := col.UpdateAll(bson.M{"_id": bson.M{"$in":ids}}, bson.M{"$pull": bson.M{"OperatorInfo": bson.M{"Flag":flag}}})
	//if err != nil && err == mgo.ErrNotFound {
	//	easygo.PanicError(err)
	//}
}

func UpdateTeanMemInfoByIds(teamId int64, ids []int64, reqMsg *share_message.OperatorInfoPer) {
	memberObj := for_game.GetRedisTeamPersonalObj(teamId)
	for _, id := range ids {
		info := memberObj.GetTeamMember(id)
		if info == nil {
			continue
		}
		info.Status = easygo.NewInt32(2)
		info.IsSave = easygo.NewBool(true)
		memberObj.UpdateTeamMemberInfo(info)
		memberObj.SetRedisOperatorInfoPer(id, reqMsg)
		memberObj.SaveToMongo()
	}
}

//移除成员封禁信息
func DelTeanMemInfoByIds(teamId int64, ids []int64, flag int32) {
	memberObj := for_game.GetRedisTeamPersonalObj(teamId)
	for _, id := range ids {
		info := memberObj.GetTeamMember(id)
		info.Status = easygo.NewInt32(1)
		info.IsSave = easygo.NewBool(true)
		memberObj.UpdateTeamMemberInfo(info)
		memberObj.DelRedisOperatorInfoPer(flag, id)
	}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAMMEMBER)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"TeamId": teamId, "PlayerId": bson.M{"$in": ids}}, bson.M{"$pull": bson.M{"OperatorInfoPer": bson.M{"Flag": flag}}, "$set": bson.M{"Status": 1}})
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
}
