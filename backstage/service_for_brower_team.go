// 管理后台为[浏览器]提供的服务
//用户管理

package backstage

import (
	"encoding/base64"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/client_hall"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"strings"

	"github.com/astaxie/beego/logs"
)

//查询群列表
func (self *cls4) RpcQueryTeamList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.GetTeamListRequest) easygo.IMessage {
	for_game.SaveRedisTeamToMongo()
	list, count := GetTeamList(user, reqMsg)

	if reqMsg.GetType() == 2 && len(list) == 0 {
		reqMsg.Type = easygo.NewInt32(1)
		list, count = GetTeamList(user, reqMsg)
	}
	ids := make([]int64, 0)
	returnList := make([]*share_message.TeamData, 0)
	for _, item := range list {
		owner := QueryPlayerbyId(item.GetOwner())
		item.OwnerNickName = easygo.NewString(owner.GetNickName())
		item.OwnerAccount = easygo.NewString(owner.GetAccount())
		if item.AdminID != nil && item.GetAdminID() != 0 {
			admin := QueryManageByID(item.GetAdminID())
			item.CreateName = easygo.NewString(admin.GetRealName())
		} else {
			item.CreateName = easygo.NewString(owner.GetNickName())
		}

		if item.GetCreateTime() > 9999999999 {
			item.CreateTime = easygo.NewInt64(item.GetCreateTime() / 1000)
		}

		if item.MessageSetting.GetIsBan() && item.MessageSetting.GetUnBanTime() != -1 && item.MessageSetting.GetUnBanTime() < util.GetMilliTime() {
			ids = append(ids, item.GetId())
		}

		returnList = append(returnList, item)
	}

	if len(ids) > 0 {
		DelTeanInfoByIds(ids, for_game.SYSTEM)
	}

	msg := &brower_backstage.GetTeamListResponse{
		List:      returnList,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//修改群资料
func (self *cls4) RpcEditTeam(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.TeamData) easygo.IMessage {
	EditTeamForHall(reqMsg) //修改群资料通知到大厅
	SaveTeamEdit(reqMsg)    //修改数据库
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.TEAM_MANAGE, "修改群资料:"+reqMsg.GetTeamChat())

	return easygo.EmptyMsg
}

//ID查询群资料
func (self *cls4) RpcGetTeamById(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	var teamData *share_message.TeamData

	if reqMsg.GetIdStr() != "" && reqMsg.IdStr != nil {
		teamData = QueryTeambyNick(reqMsg.GetIdStr())
	}

	if reqMsg.Id64 != nil {
		teamObj := for_game.GetRedisTeamObj(reqMsg.GetId64())
		teamData = teamObj.GetRedisTeam()
	}
	if teamData == nil {
		logs.Error("群数据对象为nil")
		return teamData
	}
	owner := for_game.GetRedisPlayerBase(teamData.GetOwner())
	if owner != nil {
		teamData.OwnerNickName = easygo.NewString(owner.GetNickName())
		teamData.OwnerAccount = easygo.NewString(owner.GetAccount())
		teamData.CreateName = easygo.NewString(owner.GetNickName())
	}
	if teamData.GetCreateTime() > 9999999999 {
		teamData.CreateTime = easygo.NewInt64(teamData.GetCreateTime() / 1000)
	}

	if teamData.GetHeadUrl() == "" {
		//TODO处理群默认头像
	}

	return teamData
}

//解散群
func (self *cls4) RpcDefunctTeam(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	reqMsg.IdStr = easygo.NewString(user.GetAccount())
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetId64())
	if teamObj == nil {
		return easygo.NewFailMsg("群不存在")
	}
	msg := fmt.Sprintf("解散群[%s]", teamObj.GetTeamChat())
	DefunctTeamForHall(reqMsg) //修改群资料通知到大厅
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.TEAM_MANAGE, msg)

	return easygo.EmptyMsg
}

//增减群成员
func (self *cls4) RpcTeamMemberOpt(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.MemberOptRequest) easygo.IMessage {
	logs.Info("RpcTeamMemberOpt 增减群成员")
	switch reqMsg.GetTypes() {
	case 1: //新增
		team := QueryTeambyId(reqMsg.GetTeamId())
		max := team.GetMaxMember()
		memberCount := len(team.GetMemberList())
		addMemberCount := len(reqMsg.GetAccount())
		if addMemberCount == 0 {
			return easygo.NewFailMsg("请先选择要加入的成员")
		}
		if (memberCount + addMemberCount) > int(max) {
			return easygo.NewFailMsg("增加成员失败：超过群成员上限")
		}

		TeamMemberOptForHall(user, reqMsg) // 新增群成员通知到notify
		AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.TEAM_MANAGE, "添加成员:"+easygo.StringArrayToString(reqMsg.GetAccount()))
	case 2: //删除
		TeamMemberOptForHall(user, reqMsg) //删除群成员通知到notify
		AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.TEAM_MANAGE, "删除成员:"+easygo.StringArrayToString(reqMsg.GetAccount()))
	default:
		return easygo.NewFailMsg("操作类型有误")
	}
	return easygo.EmptyMsg
}

//查询群成员列表
func (self *cls4) RpcQueryTeamMember(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.TeamMemberRequest) easygo.IMessage {
	memberObj := for_game.GetRedisTeamPersonalObj(reqMsg.GetTeamId())
	if reqMsg.GetCurPage() == 1 {
		memberObj.SaveToMongo()
	}

	list, count := QueryTeamMember(reqMsg)
	var ids PLAYER_IDS
	for _, obj := range list {
		obj.OperatorInfoPer = memberObj.GetRedisOperatorInfoPer(obj.GetPlayerId(), nil)
		if infos := obj.GetOperatorInfoPer(); infos != nil {
			for _, info := range infos {
				if info.GetCloseTime() > util.GetMilliTime() || obj.GetStatus() == 1 { //已过禁言时间
					continue
				}
				ids = append(ids, obj.GetPlayerId())
			}
		}
	}

	if len(ids) > 0 {
		DelTeanMemInfoByIds(reqMsg.GetTeamId(), ids, 1)
	}

	msg := &brower_backstage.TeamMemberResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//查询群聊天记录列表
func (self *cls4) RpcExportChatRecord(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	if reqMsg.GetId() == 0 || reqMsg.Id == nil {
		return easygo.NewFailMsg("群ID错误")
	}

	list, count := QueryChatRecord(reqMsg)
	logs.Info("team", list)

	msg := &brower_backstage.ExportChatRecordResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

func (self *cls4) RpcQueryTeamMessage(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {

	if reqMsg.GetKeyword() == "" || reqMsg.Keyword == nil {
		return easygo.NewFailMsg("群ID错误")
	}

	list, count := QueryTeamMessage(reqMsg)

	msg := &brower_backstage.ExportChatRecordResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//创建群
func (self *cls4) RpcCreateTeamMessage(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.CreateTeamInfo) easygo.IMessage {
	logs.Info("RpcCreateTeamMessage 后台创建群")
	if reqMsg.PlayerID == nil {
		return easygo.NewFailMsg("无效的玩家ID")
	}
	result := SendToPlayer(reqMsg.GetPlayerID(), "RpcCreateTeam", &server_server.CreateTeamInfo{PlayerID: reqMsg.PlayerID, TeamName: reqMsg.TeamName, AdminID: user.Id, AdminName: user.RealName})
	err := for_game.ParseReturnDataErr(result)
	if err != nil {
		return err
	}
	msg := result.(*server_server.CreateTeamResult)
	team := QueryTeambyId(msg.GetTeamID())
	owner := QueryPlayerbyId(team.GetOwner())
	if owner != nil {
		team.OwnerNickName = easygo.NewString(owner.GetNickName())
		team.OwnerAccount = easygo.NewString(owner.GetAccount())
	}
	if team.AdminID != nil {
		manager := QueryManageByID(team.GetAdminID())
		team.CreateName = easygo.NewString(manager.GetRealName())
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.TEAM_MANAGE, "创建群:"+team.GetTeamChat())
	return team
}

//获取拉人进群用户列表
func (self *cls4) RpcQueryTeamPlayerList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.GetTeamPlayerListRequest) easygo.IMessage {
	logs.Info("RpcQueryTeamPlayerList 获取拉人进群用户列表")
	if reqMsg.TeamId == nil {
		return easygo.NewFailMsg("无效群ID")
	}
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetTeamId())
	team := teamObj.GetRedisTeam()
	if team == nil {
		return easygo.NewFailMsg("找不到这个群")
	}

	list, count := GetTeamPlayerList(team, reqMsg)

	msg := &brower_backstage.GetTeamPlayerListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
		Team:      team,
	}
	return msg
}

//群封禁
func (self *cls4) RpcTeamBan(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	req := &server_server.TeamBan{
		Ids:     reqMsg.Ids64,
		BanTime: easygo.NewInt64(reqMsg.GetIds32()[0]),
		Status:  easygo.NewInt32(1),
		Node:    reqMsg.Note,
	}
	// BroadCastToAllHall("RpcTeamBanHall", req)
	ChooseOneHall(0, "RpcTeamBanHall", req)

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}

	msg := fmt.Sprintf("批量封禁群: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)
	return easygo.EmptyMsg
}

//群解封
func (self *cls4) RpcTeamUnBan(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	req := &server_server.TeamBan{
		Ids: reqMsg.Ids64,
		// BanTime: easygo.NewInt64(reqMsg.GetIds32()[0]),
		Status: easygo.NewInt32(2),
		Node:   reqMsg.Note,
	}
	// BroadCastToAllHall("RpcTeamBanHall", req)
	ChooseOneHall(0, "RpcTeamBanHall", req)

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}

	msg := fmt.Sprintf("批量解封群: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)
	return easygo.EmptyMsg
}

//解封群
func (self *cls4) RpcTeamCloseAndOpen(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.TeamManager) easygo.IMessage {
	message := &brower_backstage.ErrMessage{}
	message.Err = nil
	if reqMsg.GetTeamIds() == nil {
		return easygo.NewFailMsg("teamids为空")
	}
	now := util.GetMilliTime()

	//天转毫秒
	dayToMill := int64(int(reqMsg.GetDay()) * 24 * 3600000)
	//小时转毫秒
	hourToMill := int64(int(reqMsg.GetHour()) * 3600000)
	//分钟转毫秒
	minToMill := int64(int(reqMsg.GetMinutes()) * 60000)
	closeTime := now + hourToMill + minToMill + dayToMill
	if reqMsg.GetCloseTime() < 0 {
		closeTime = reqMsg.GetCloseTime()
	}

	info := &share_message.OperatorInfo{
		Operator:  easygo.NewString(user.Account),
		Time:      easygo.NewInt64(now),
		Flag:      easygo.NewInt32(1),
		CloseTime: easygo.NewInt64(closeTime),
		Reason:    easygo.NewString(reqMsg.GetReason()),
	}

	var flag1 string
	var flat string
	isBan := false
	switch reqMsg.GetFlag() {
	case 1:
		flag1 = "批量解封群"
		flat = "关闭"
		DelTeanInfoByIds(reqMsg.GetTeamIds(), for_game.SYSTEM) //移除封禁信息
		message.Err = easygo.NewString("成功解除了对群" + strings.Join(reqMsg.GetNickName(), "、") + "的封禁")
	case 2:
		flag1 = "批量封禁群"
		flat = "开启"
		isBan = true
		UpdateTeanInfoByIds(reqMsg.GetTeamIds(), info) //更新群状态信息
	}

	for _, teamId := range reqMsg.GetTeamIds() {
		msg2 := &share_message.TeamMessage{
			Type:        easygo.NewInt32(for_game.BACKSTAGE_BAN_TEAM),
			PlayerId:    easygo.NewInt64(0),
			IsAllNotice: easygo.NewBool(true),
			Time:        easygo.NewInt64(closeTime),
			TeamId:      easygo.NewInt64(teamId),
			Name:        easygo.NewString(user.RealName),
			Pos:         easygo.NewInt32(0),
			SendTime:    easygo.NewInt64(now),
			Value:       easygo.NewBool(isBan),
		}
		req := &server_server.TeamManager{
			TeamIds:   reqMsg.GetTeamIds(),
			SendTime:  easygo.NewInt64(now),
			Flag:      easygo.NewInt32(reqMsg.GetFlag()),
			Reason:    easygo.NewString(reqMsg.GetReason()),
			CloseTime: easygo.NewInt64(closeTime),
			Name:      easygo.NewString(user.RealName),
		}

		content := fmt.Sprintf("系统已%s全员禁言", flat)
		ti := for_game.GetMillSecond()
		chat := &share_message.Chat{
			Content:     easygo.NewString(content),
			ContentType: easygo.NewInt32(0),
			Time:        easygo.NewInt64(ti),
		}
		teamObj := for_game.GetRedisTeamObj(teamId)
		session := for_game.GetRedisChatSessionObj(easygo.AnytoA(teamId))
		logId := session.GetNextMaxLogId()
		for_game.AddTeamChatLog(teamId, 0, logId, chat, nil, true, msg2)
		req.LogId = easygo.NewInt64(logId)

		msgteammer := &client_hall.OperatorMessage{
			Name:      easygo.NewString(req.GetName()),
			TeamId:    easygo.NewInt64(teamId),
			SendTime:  easygo.NewInt64(req.GetSendTime()),
			Flag:      easygo.NewInt64(req.GetFlag()),
			CloseTime: easygo.NewInt64(req.GetCloseTime()),
			LogId:     easygo.NewInt64(req.GetLogId()),
		}
		playerIds := teamObj.GetTeamMemberList() //获取playerid
		SendMsgToHallClientNew(playerIds, "RpcTeamChangeInfo", msgteammer)
	}
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetTeamIds()[0])
	msg := fmt.Sprintf("%s: %v", flag1, teamObj.GetRedisTeam().GetTeamChat())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)
	return message
}

//解封群成员
func (self *cls4) RpcTeamMemCloseAndOpen(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.TeamManager) easygo.IMessage {
	logs.Info("解封群成员RpcTeamMemCloseAndOpen:", reqMsg)
	message := &brower_backstage.ErrMessage{}
	message.Err = nil
	if reqMsg.GetTeamIds() == nil {
		return easygo.NewFailMsg("teamids为空")
	}
	now := util.GetMilliTime()

	//天转毫秒
	dayToMill := int64(int(reqMsg.GetDay()) * 24 * 3600000)
	//小时转毫秒
	hourToMill := int64(int(reqMsg.GetHour()) * 3600000)
	//分钟转毫秒
	minToMill := int64(int(reqMsg.GetMinutes()) * 60000)

	closeTime := now + hourToMill + minToMill + dayToMill

	info := &share_message.OperatorInfoPer{
		Operator:  user.Account,
		Time:      easygo.NewInt64(now),
		Flag:      easygo.NewInt32(1),
		CloseTime: easygo.NewInt64(closeTime),
		Reason:    easygo.NewString(reqMsg.GetReason()),
	}

	req := &server_server.TeamManager{
		TeamId:    easygo.NewInt64(reqMsg.GetTeamId()),
		NickName:  reqMsg.GetNickName(),
		Name:      easygo.NewString(user.RealName),
		TeamIds:   reqMsg.GetTeamIds(),
		Flag:      easygo.NewInt32(reqMsg.GetFlag()),
		CloseTime: easygo.NewInt64(closeTime),
		SendTime:  easygo.NewInt64(now),
		Day:       easygo.NewInt32(reqMsg.GetDay()),
		Minutes:   easygo.NewInt32(reqMsg.GetMinutes()),
		Hour:      easygo.NewInt32(reqMsg.GetHour()),
	}

	msg2 := &share_message.TeamMessage{
		Type:        easygo.NewInt32(for_game.BACKSTAGE_BAN_TEAM_MEM),
		PlayerId:    easygo.NewInt64(0),
		IsAllNotice: easygo.NewBool(true),
		Time:        easygo.NewInt64(closeTime),
		TeamId:      easygo.NewInt64(reqMsg.GetTeamId()),
		Name:        easygo.NewString(user.RealName),
		Pos:         easygo.NewInt32(0),
		SendTime:    easygo.NewInt64(now),
		PlayerList:  reqMsg.GetTeamIds(),
	}

	var flag string

	switch reqMsg.GetFlag() {
	case 1:
		flag = "批量解封群成员"
		DelTeanMemInfoByIds(reqMsg.GetTeamId(), reqMsg.GetTeamIds(), 1)
		message.Err = easygo.NewString("成功解除了成员" + strings.Join(reqMsg.GetNickName(), "、") + "的封禁")
		msg2.Value = easygo.NewBool(false)
		msg2.Value1 = easygo.NewString(base64.StdEncoding.EncodeToString([]byte(strings.Join(reqMsg.GetNickName(), "、"))))
	case 2:
		flag = "批量封禁群成员"
		UpdateTeanMemInfoByIds(reqMsg.GetTeamId(), reqMsg.GetTeamIds(), info) //更新群状态信息
		msg2.Value = easygo.NewBool(true)
		namelst := ""
		memberObj := for_game.GetRedisTeamPersonalObj(reqMsg.GetTeamId())
		for _, pid := range msg2.PlayerList {
			name := memberObj.GetTeamMemberReName(pid)
			if name == "" {
				base := for_game.GetRedisPlayerBase(pid)
				name = base.GetNickName()
			}
			namelst += fmt.Sprintf(`"%s"、`, name)
		}
		if len(namelst) > 0 {
			namelst = namelst[:len(namelst)-3] //为什么减3  因为一个顿号在字符串中占3个索引
		}
		msg2.Value1 = easygo.NewString(base64.StdEncoding.EncodeToString([]byte(namelst)))
	}

	account := []string{}

	for _, value := range reqMsg.GetTeamIds() {
		name := for_game.GetRedisPlayerBase(value).GetAccount()
		account = append(account, name)
	}
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetTeamId())
	msg := fmt.Sprintf("%s: %v:%v", flag, teamObj.GetRedisTeam().GetTeamChat(), account)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)
	chat := &share_message.Chat{
		Content:     easygo.NewString(""),
		ContentType: easygo.NewInt32(0),
		Time:        easygo.NewInt64(now),
	}
	session := for_game.GetRedisChatSessionObj(easygo.AnytoA(reqMsg.GetTeamId()))
	logId := session.GetNextMaxLogId()
	for_game.AddTeamChatLog(reqMsg.GetTeamId(), 0, logId, chat, nil, true, msg2)
	req.LogId = easygo.NewInt64(logId)
	playerIds := teamObj.GetTeamMemberList() //获取playerid
	for _, playerId := range playerIds {
		p := for_game.GetRedisPlayerBase(playerId)
		if p != nil && p.GetIsOnLine() {
			req.PlayerId = easygo.NewInt64(playerId)
			SendToPlayer(playerId, "RpcTeamMemCloseAndOpen", req)
		}
	}

	return message
}

//警告群主
func (self *cls4) RpcWarnLord(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	if len(reqMsg.GetIds64()) == 0 {
		return easygo.NewFailMsg("至少选择一个群")
	}
	s := "您创建的群:%s,经用户投诉，并经审核确认，您的群在使用柠檬畅聊过程中存在违规行为，请注意规范使用柠檬畅聊账号、文明沟通。多次出现违规行为，将导致账号被进行相应处理，感谢您的理解与支持。"
	reqMsg.Note = easygo.NewString(s)
	WarnLordForHall(reqMsg)

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}

	msg := fmt.Sprintf("警告群: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.TEAM_MANAGE, msg)
	return easygo.EmptyMsg
}
