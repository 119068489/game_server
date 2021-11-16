package hall

import (
	"encoding/base64"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
	"time"

	"github.com/astaxie/beego/logs"
)

func CreateTeam(owner int64, plst []int64, info *client_hall.CreateTeam) *share_message.TeamData {
	teamId := for_game.NextId(for_game.TABLE_TEAM_DATA, for_game.INIT_TEAM_ID)
	teamchat := for_game.GetRandAccount("nmq", teamId)
	ti := for_game.GetMillSecond()
	headURL := info.GetHeadUrl()
	adminId := info.GetAdminId()
	if headURL == "" {
		headURL = for_game.GetRandTeamHeadIcon()
	}
	team := &share_message.TeamData{
		Id:             easygo.NewInt64(teamId),
		Name:           easygo.NewString(info.GetTeamName()),
		HeadUrl:        easygo.NewString(headURL),
		GongGao:        easygo.NewString(""),
		Owner:          easygo.NewInt64(owner),
		QRCode:         easygo.NewString(""),
		CreateTime:     easygo.NewInt64(ti),
		LastTalkTime:   easygo.NewInt64(ti),
		MaxMember:      easygo.NewInt32(for_game.NORMAL_TEAMMERBER),
		TeamChat:       easygo.NewString(teamchat),
		MessageSetting: GetDefaultMessageSetting(),
		Status:         easygo.NewInt32(0),
		IsRecommend:    easygo.NewBool(false),
		DissolveTime:   easygo.NewInt64(0),
		RefreshTime:    easygo.NewInt64(0),
		Level:          easygo.NewInt32(1), //默认等级1
		WelcomeWord:    easygo.NewString("欢迎进入本群。文明交流，大家一起热闹，禁止发广告~"),
		Topic:          easygo.NewString(info.GetTopic()),
		TopicDesc:      easygo.NewString(info.GetTopicDesc()),
	}
	if adminId != 0 {
		team.AdminID = easygo.NewInt64(adminId)
	}
	team.MemberList = plst
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	defer closeFun()
	err := col.Insert(team)
	easygo.PanicError(err)
	return team
}

func GetDefaultMessageSetting() *share_message.MessageSetting {
	msg := &share_message.MessageSetting{
		IsTimeClean:         easygo.NewBool(false),
		IsReadClean:         easygo.NewBool(false),
		IsScreenShotNotify:  easygo.NewBool(false),
		IsStopTalk:          easygo.NewBool(false),
		TeamHelp:            easygo.NewInt32(0),
		IsAddFriend:         easygo.NewBool(false),
		IsInvite:            easygo.NewBool(false),
		IsStopAddTeam:       easygo.NewBool(false),
		IsOpenTeamMoneyCode: easygo.NewBool(false),
		IsOpenWelcomeWord:   easygo.NewBool(true),
		IsManagerEdit:       easygo.NewBool(false),
	}
	return msg
}

func AddTeamMember(teamId, maxId int64, plst []int64, b bool, t int32, reason string, serverId int32, msg *share_message.TeamChannel) bool {
	success := for_game.AddTeamPersonData(teamId, maxId, plst, for_game.TEAM_MASSES, reason, msg)
	if success {
		teamObj := for_game.GetRedisTeamObj(teamId)
		teamObj.AddTeamMember(plst)
		msg := for_game.GetTeamMsgForHall(teamId)
		msg.IsShow = easygo.NewBool(b)
		if t != 0 {
			msg.Type = easygo.NewInt32(t)
		}
		TeamSendMessage(plst, 0, serverId, "RpcAddTeamResult", msg)
	}

	return success
}

func DeleteTeamMember(teamId int64, plst []int64, serverId int32) bool {
	success := RemoveTeamMember(teamId, plst)
	if success {
		msg := &client_server.TeamInfo{
			TeamId: easygo.NewInt64(teamId),
		}
		TeamSendMessage(plst, 0, serverId, "RpcTeamOutPlayer", msg)
	}
	return success
}

func RemoveTeamMember(teamId int64, plst []int64) bool {
	var b bool
	teamObj := for_game.GetRedisTeamObj(teamId)
	memberObj := for_game.GetRedisTeamPersonalObj(teamId)
	memberList := teamObj.GetTeamMemberList()
	managerList := teamObj.GetTeamManageList()
	delManagerList, delMemberList := []int64{}, []int64{}
	for _, playerId := range plst {
		if !util.Int64InSlice(playerId, memberList) { //
			continue
		}
		memberList = easygo.Del(memberList, playerId).([]int64)
		delMemberList = append(delMemberList, playerId)
		pos := for_game.GetTeamPlayerPos(teamId, playerId)
		if pos == for_game.TEAM_OWNER { //如果群主退群
			b = true
		} else if pos == for_game.TEAM_MANAGER { //如果管理员退群或者被踢出
			delManagerList = append(delManagerList, playerId)
		}
	}
	if b {
		var id int64
		if len(managerList) != 0 {
			id = managerList[0]
		} else {
			if len(memberList) != 0 {
				id = memberList[0]
			}
		}
		if id != 0 { // 说明群主是最后一个成员
			oldOwner := teamObj.GetTeamOwner()
			teamObj.SetTeamOwner(id)

			for_game.SetTeamPlayerPos(teamId, id, for_game.TEAM_OWNER)
			NoticeTeamMessage(teamId, oldOwner, for_game.CHANGE_OWNER, []int64{id})
			if util.Int64InSlice(id, managerList) { //如果随机群主是管理员 把它从管理员列表中删除
				delManagerList = append(delManagerList, id)
			}
		}
	}

	teamObj.DelTeamMemberList(delMemberList)
	memberObj.DelTeamPersonData(delMemberList)
	if len(delManagerList) != 0 {
		teamObj.DelTeamManageList(delManagerList)
	}
	if len(teamObj.GetTeamMemberList()) == 0 {
		teamObj.SetTeamStatus(1)
		teamObj.SetTeamDissolveTime(time.Now().Unix())
	}
	return true
}

//设置群职位
func SetTeamMemberPosition(teamId int64, playerList []int64, pos int32) {
	if len(playerList) == 0 {
		panic("玩家列表为空")
	}
	teamObj := for_game.GetRedisTeamObj(teamId)
	managerList := teamObj.GetTeamManageList()
	if pos == for_game.TEAM_MANAGER { //设置为管理员
		addManagerList := []int64{}
		for _, pid := range playerList {
			if util.Int64InSlice(pid, managerList) {
				continue
			}
			addManagerList = append(addManagerList, pid)
			for_game.SetTeamPlayerPos(teamId, pid, for_game.TEAM_MANAGER)
		}
		if len(addManagerList) > 0 {
			teamObj.AddTeamManage(addManagerList)
		}
	} else if pos == for_game.TEAM_OWNER { //设置为群主
		oldOwner := teamObj.GetTeamOwner()
		pid := playerList[0]
		teamObj.SetTeamOwner(pid)
		for_game.SetTeamPlayerPos(teamId, pid, for_game.TEAM_OWNER)
		for_game.SetTeamPlayerPos(teamId, oldOwner, for_game.TEAM_MASSES)
		if util.Int64InSlice(pid, managerList) {
			teamObj.DelTeamManageList([]int64{pid})
		}
	} else if pos == for_game.TEAM_MASSES {
		delManageList := []int64{}
		for _, pid := range playerList {
			if util.Int64InSlice(pid, managerList) {
				delManageList = append(delManageList, pid)
			}
			for_game.SetTeamPlayerPos(teamId, pid, for_game.TEAM_MASSES)
		}
		if len(delManageList) != 0 {
			teamObj.DelTeamManageList(delManageList)
		}
	}
}

func GetTeamPlayerInfo(teamId, pid, playerId int64, currentPage, pageSize int32) *share_message.TeamPlayerInfo {
	rename := for_game.GetFriendsReName(playerId, pid)
	memberObj := for_game.GetRedisTeamPersonalObj(teamId)
	teamName := memberObj.GetTeamMemberReName(pid)
	channel := memberObj.GetTeamMemberChannel(pid)
	base := for_game.GetRedisPlayerBase(pid)
	if base == nil {
		return nil
	}
	//pageMap := base.GetRedisPlayerDynamicListByPage(playerId, currentPage, pageSize) // 分页
	//logIds := pageMap["arr"].([]int64)
	ds, count := base.GetRedisPlayerDynamicListByPage(playerId, currentPage, pageSize) // 分页
	dynamicData := &share_message.DynamicDataListPage{
		//DynamicData: for_game.GetRedisDynamicForSomeLogId(isTop, pid, logIds, playerId),
		DynamicData: for_game.GetRedisDynamicForSomeLogId2(pid, ds, playerId),
		TotalCount:  easygo.NewInt32(count),
	}
	// 获取后台加赞
	trueZan := for_game.GetAllTrueZan(base.GetPlayerId())
	var addFriendType int32
	if friendBase := for_game.GetFriendBase(playerId); friendBase != nil {
		for _, v := range friendBase.GetFriends() {
			if v.GetPlayerId() == pid {
				addFriendType = v.GetType()
			}
		}
	}
	info := &share_message.TeamPlayerInfo{
		PlayerId:           easygo.NewInt64(pid),
		Account:            easygo.NewString(base.GetAccount()),
		NickName:           easygo.NewString(base.GetNickName()),
		HeadIcon:           easygo.NewString(base.GetHeadIcon()),
		Sex:                easygo.NewInt32(base.GetSex()),
		Photo:              base.GetPhoto(),
		Phone:              easygo.NewString(base.GetPhone()),
		Signature:          easygo.NewString(base.GetSignature()),
		Provice:            easygo.NewString(base.GetProvice()),
		City:               easygo.NewString(base.GetCity()),
		Channel:            easygo.NewString(channel),
		ReName:             easygo.NewString(rename),   //好友备注
		TeamName:           easygo.NewString(teamName), //
		TeamId:             easygo.NewInt64(teamId),
		Zans:               easygo.NewInt32(base.GetZan() + trueZan),
		Fans:               easygo.NewInt32(len(base.GetFans())),
		Attentions:         easygo.NewInt32(len(base.GetAttention())),
		Icon:               easygo.NewInt32(0),
		DynamicData:        dynamicData,
		AccountState:       easygo.NewInt32(base.GetStatus()),
		BackgroundImageURL: easygo.NewString(base.GetBackgroundImageURL()),
		Types:              easygo.NewInt32(base.GetTypes()),
		AddFriendType:      easygo.NewInt32(addFriendType),
		TeamPosition:       easygo.NewInt32(memberObj.GetTeamPosition(pid)),
	}
	return info
}

func NoticeTeamMessage(teamId, pid int64, t int32, value interface{}, arg ...int64) {
	adminId := append(arg, 0)[0]
	var content string
	memberObj := for_game.GetRedisTeamPersonalObj(teamId)
	teamObj := for_game.GetRedisTeamObj(teamId)
	name := memberObj.GetTeamMemberReName(pid)
	if name == "" {
		if base := for_game.GetRedisPlayerBase(pid); base != nil {
			name = base.GetNickName()
		}
	}

	ti := for_game.GetMillSecond()
	msg := &share_message.TeamMessage{
		PlayerId:    easygo.NewInt64(pid),
		Type:        easygo.NewInt32(t),
		IsAllNotice: easygo.NewBool(true),
		Time:        easygo.NewInt64(ti),
		TeamId:      easygo.NewInt64(teamId),
		Name:        easygo.NewString(name),
		ShowPos:     easygo.NewInt32(for_game.TEAM_MASSES), //默认所有群成员可以看见
	}
	if t != for_game.EXIT_PLAYER {
		msg.Pos = easygo.NewInt32(for_game.GetTeamPlayerPos(teamId, pid))
	}

	if adminId != 0 {
		msg.AdminID = easygo.NewInt64(adminId)
	}
	switch v := value.(type) {
	case []int64:
		msg.PlayerList = []int64(v)
	case bool:
		msg.Value = easygo.NewBool(bool(v))
	case string:
		msg.Value1 = easygo.NewString(string(v))
	}

	if t == for_game.CHANGE_OWNER {
		targetId := msg.PlayerList[0]
		name1 := memberObj.GetTeamMemberReName(targetId)
		if name1 == "" {
			base1 := GetPlayerObj(targetId)
			name1 = base1.GetNickName()
		}
		content = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`"%s"转让群主给"%s"`, name, name1)))
		msg.Value1 = easygo.NewString(content)
		teamObj.SetTeamIsOpenTeamMoneyCode(false)
		teamObj.SetTeamQRCode("")
	} else if t == for_game.ADD_MANAGER || t == for_game.DEL_MANAGER {
		msg.AllPlayerList = teamObj.GetTeamManageList()

		namelst := ""
		for _, pid := range msg.PlayerList {
			name := memberObj.GetTeamMemberReName(pid)
			if name == "" {
				player := GetPlayerObj(pid)
				name = player.GetNickName()
			}
			namelst += fmt.Sprintf(`"%s"、`, name)
		}
		if len(namelst) > 0 {
			namelst = namelst[:len(namelst)-3] //为什么减3  因为一个顿号在字符串中占3个索引
		}
		content := base64.StdEncoding.EncodeToString([]byte(namelst))
		msg.Value1 = easygo.NewString(content)
	} else if t == for_game.INVITE_PLAYER || t == for_game.ACTIVE_ADDTEAM || t == for_game.ADV_TEAM_MEM {
		content := ""
		if len(msg.GetPlayerList()) == 1 {
			pid := msg.GetPlayerList()[0]
			name1 := memberObj.GetTeamMemberReName(pid)
			player := for_game.GetRedisPlayerBase(pid)
			if name1 == "" {
				name1 = player.GetNickName()
			}
			content += fmt.Sprintf(`"%s"`, name1)
			msg.Members = append(msg.Members, memberObj.GetTeamMember(pid, player.GetRedisPlayerBase()))
		} else {
			for n, pid := range msg.GetPlayerList() {
				name1 := memberObj.GetTeamMemberReName(pid)
				player := for_game.GetRedisPlayerBase(pid)
				if name1 == "" {
					name1 = player.GetNickName()
				}
				if n < 3 {
					content += fmt.Sprintf(`"%s"、`, name1)
				}
				msg.Members = append(msg.Members, memberObj.GetTeamMember(pid, player.GetRedisPlayerBase()))
			}
			if len(msg.GetPlayerList()) > 3 {
				s := fmt.Sprintf("等%d人", len(msg.GetPlayerList()))
				content += s
			}
		}
		logs.Info("contemt:", content)
		msg.Value1 = easygo.NewString(base64.StdEncoding.EncodeToString([]byte(content)))
	} else if t == for_game.REQUEST_ADDTEAM {
		id := msg.PlayerList[0]
		base := GetPlayerObj(id)
		msg.Value1 = easygo.NewString(base64.StdEncoding.EncodeToString([]byte(base.GetNickName())))
		msg.Account = easygo.NewString(base.GetAccount())
		msg.ShowPos = easygo.NewInt32(for_game.TEAM_MANAGER)
	} else if t == for_game.WITHDRAW_MESSAGE {
		msg.Value1 = easygo.NewString(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s撤回了一条消息", name))))
	} else if t == for_game.EXIT_PLAYER {
		content = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`"%s"退出了群聊`, name)))
		msg.Value1 = easygo.NewString(content)
		//退群信息只有管理员权限以上用户才能看见
		msg.ShowPos = easygo.NewInt32(for_game.TEAM_MANAGER)
	} else if t == for_game.DEL_PLAYER {
		namelst := ""
		for _, pid := range msg.PlayerList {
			name := memberObj.GetTeamMemberReName(pid)
			if name == "" {
				base := GetPlayerObj(pid)
				name = base.GetNickName()
			}
			namelst += fmt.Sprintf(`"%s"、`, name)
		}
		if len(namelst) > 0 {
			namelst = namelst[:len(namelst)-3] //为什么减3  因为一个顿号在字符串中占3个索引
		}
		msg.Value1 = easygo.NewString(base64.StdEncoding.EncodeToString([]byte(namelst)))
		//退群信息只有管理员权限以上用户才能看见
		msg.ShowPos = easygo.NewInt32(for_game.TEAM_MANAGER)
	}
	/*else if t == for_game.WELCOME_WORD {
		if msg.GetValue() {
			content = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`"%s开启了群欢迎语"`, name)))
		} else {
			content = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`"%s关闭了群欢迎语"`, name)))
		}
		msg.Value1 = easygo.NewString(content)
	} else if t == for_game.WELCOME_WORD_MANAGER {
		if msg.GetValue() {
			content = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`"群主开启了管理员可编辑欢迎语功能"`)))
		} else {
			content = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`"群主关闭了管理员可编辑欢迎语功能"`)))
		}
		msg.Value1 = easygo.NewString(content)
	} else if t == for_game.EDIT_WELCOME_WORD {
		content = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`"%s把群欢迎语修改为:%s"`, name, msg.GetValue())))
		msg.Value1 = easygo.NewString(content)
	}
	*/

	if t != for_game.WITHDRAW_MESSAGE {

		chat := &share_message.Chat{
			Content:     easygo.NewString(content),
			ContentType: easygo.NewInt32(for_game.TALK_CONTENT_SYSTEM), //默认系统通知
			Time:        easygo.NewInt64(ti),
			SessionId:   easygo.NewString(easygo.AnytoA(teamId)),
		}
		if t == for_game.TEAM_GONGGAO {
			chat.Content = easygo.NewString(msg.GetValue1())
			chat.ContentType = easygo.NewInt32(for_game.TALK_CONTENT_GROUPNOTICE) //群公告
			chat.SessionId = easygo.NewString(easygo.AnytoA(teamId))
			chat.TargetId = easygo.NewInt64(teamId)
			chat.ChatType = easygo.NewInt32(for_game.CHAT_TYPE_TEAM)
			player := for_game.GetRedisPlayerBase(pid)
			if player != nil {
				chat.SourceId = easygo.NewInt64(memberObj.GetId())
				chat.SourceName = easygo.NewString(player.GetNickName())
				chat.SourceHeadIcon = easygo.NewString(player.GetHeadIcon())
				chat.NoticeInfo = &share_message.NoticeInfo{
					IsAll: easygo.NewBool(true),
				}
				equiptment := for_game.GetRedisPlayerEquipmentObj(pid)
				if equiptment != nil {
					eq := equiptment.GetEquipmentForClient()
					chat.QPId = easygo.NewInt64(eq.GetQP().GetPropsId())
				}
				cl := &cls1{}
				p := GetPlayerObj(pid)
				cl.RpcChatNew(nil, p, chat)
			}

		} else {
			session := for_game.GetRedisChatSessionObj(easygo.AnytoA(teamId))
			logId := session.GetNextMaxLogId()
			for_game.AddTeamChatLog(teamId, pid, logId, chat, nil, true, msg)
			msg.LogId = easygo.NewInt64(logId)
		}
	} else {
		msg.LogId = easygo.NewInt64(value.(int64))
	}
	serverId := PlayerOnlineMgr.GetPlayerServerId(pid)
	members := teamObj.GetTeamMemberList()
	if t == for_game.WITHDRAW_MESSAGE {
		members = easygo.Del(members, pid).([]int64)
	}
	TeamSendMessage(members, 0, serverId, "RpcTeamNoticeMessage", msg)
	//进群发送欢迎语
	if t == for_game.INVITE_PLAYER || t == for_game.ACTIVE_ADDTEAM || t == for_game.ADV_TEAM_MEM {
		setting := teamObj.GetTeamMessageSetting()
		if setting.GetIsOpenWelcomeWord() {
			logs.Info("欢迎语新成员:", teamId, msg.GetPlayerList())
			cl := &cls1{}
			player := GetPlayerObj(teamObj.GetTeamOwner())
			chat := teamObj.GetSendWelComeWord(msg.GetPlayerList())
			chat.IsWelcome = easygo.NewBool(true) // 设置是入欢迎语
			//ep := ClientEpMgr.LoadEndpointByPid(player.GetPlayerId())
			//cl.RpcChat(ep, player, chat)
			cl.RpcChatNew(nil, player, chat)
		}
	}
}

func TeamSendMessage(plst []int64, ownerId int64, serverId int32, methodName string, msg easygo.IMessage) {
	info := GetServerPlayerMap(plst, ownerId)
	for sId, lst := range info {
		if sId == PServerInfo.GetSid() { //如果在本服务器的玩家，直接通知就好
			SendToCurrentHallClient(lst, methodName, msg)
		} else {
			BroadCastMsgToHallClientNew(lst, methodName, msg)
		}
	}
}
