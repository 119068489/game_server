// 管理后台为[浏览器]提供的服务
//用户管理

package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

//普通用户列表
func (self *cls4) RpcPlayerList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.GetPlayerListRequest) easygo.IMessage {
	if reqMsg.GetCurPage() == 1 {
		for_game.SaveRedisPlayerBaseToMongo()
	}
	list, count := GetPlayerList(user, reqMsg)

	msg := &brower_backstage.GetPlayerListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//修改用户
func (self *cls4) RpcEditPlayer(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PlayerBase) easygo.IMessage {
	// reqMsg.Password = easygo.NewString(reqMsg.GetPassword())
	if reqMsg.NickName == nil || reqMsg.GetNickName() == "" {
		return easygo.NewFailMsg("昵称不能为空")
	}
	if reqMsg.Phone == nil || reqMsg.GetPhone() == "" {
		return easygo.NewFailMsg("手机号不能为空")
	}
	playerAccount := for_game.GetRedisAccountObj(reqMsg.GetPlayerId())
	if playerAccount.GetAccount() != reqMsg.GetPhone() {
		oldAccount := for_game.GetPlayerByPhone(reqMsg.GetPhone())
		if oldAccount != nil {
			return easygo.NewFailMsg("手机号已存在")
		}
	}
	EditPlayer(reqMsg, false)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, "修改用户:"+reqMsg.GetAccount()+"详细资料")

	return easygo.EmptyMsg
}

//批量修改玩家兴趣标签
func (self *cls4) RpcEditPlayerLable(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	if reqMsg.Ids64 == nil || len(reqMsg.GetIds64()) == 0 {
		return easygo.NewFailMsg("玩家Id不能为空")
	}
	if reqMsg.Ids32 == nil || len(reqMsg.GetIds32()) == 0 {
		return easygo.NewFailMsg("标签不能为空")
	}

	var ids string
	players := GetPlayerBaseByIds(reqMsg.GetIds64())
	count := len(players)
	for i, p := range players {
		if p.GetTypes() == 1 || p.GetTypes() == 3 {
			return easygo.NewFailMsg("只允许修改运营账号的兴趣标签")
		}

		if i < count {
			ids += p.GetAccount() + ","
		} else {
			ids += p.GetAccount()
		}
	}

	EditPlayerLable(reqMsg)

	msg := fmt.Sprintf("批量修改用户标签: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)

	return easygo.EmptyMsg
}

//创建用户
func (self *cls4) RpcAddPlayer(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.SigninRequest) easygo.IMessage {
	if reqMsg.Phone == nil {
		return easygo.NewFailMsg("手机号不能为空")
	}
	// if reqMsg.Code == nil {
	// 	return easygo.NewFailMsg("验证码不能为空")
	// }
	if reqMsg.Password == nil {
		return easygo.NewFailMsg("密码不能为空")
	}

	// data := for_game.MessageMarkInfo.GetMessageMarkInfo(2, reqMsg.GetPhone())
	// if data == nil {
	// 	res := "验证码不存在"
	// 	return easygo.NewFailMsg(res, "1001")
	// }
	// if data.Mark != reqMsg.GetCode() {
	// 	res := "验证码不正确"
	// 	return easygo.NewFailMsg(res)
	// }
	acc := for_game.GetRedisAccountByPhone(reqMsg.GetPhone())
	if acc != nil {
		return easygo.NewFailMsg("账号重复！请修改账号重新确认添加")
	}
	ip := ep.GetConnection().RemoteAddr().String()
	data := &share_message.CreateAccountData{
		Phone:    easygo.NewString(reqMsg.GetPhone()),
		PassWord: easygo.NewString(reqMsg.GetPassword()),
		Ip:       easygo.NewString(ip),
		Types:    easygo.NewInt32(reqMsg.GetTypes()),
	}
	b, playerId := for_game.CreateAccount(data)
	if b {
		player := for_game.GetRedisPlayerBase(playerId)
		player.SetNickName(player.GetAccount())
		player.SetDeviceType(3)
		player.SetCreateTime()
		player.SetHeadIcon(for_game.GetDefaultHeadicon(2))
		player.SaveToMongo()

		AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, "创建用户:"+player.GetAccount())
	}
	return easygo.EmptyMsg
}

//创建运营用户
func (self *cls4) RpcAddWaiter(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.AddWaiterRequest) easygo.IMessage {
	if reqMsg.Type == nil {
		return easygo.NewFailMsg("用户类型不能为空")
	}
	if reqMsg.Count == nil || reqMsg.GetCount() == 0 {
		return easygo.NewFailMsg("创建数量错误")
	}

	if reqMsg.GetCount() < 1 || reqMsg.GetCount() > 99 {
		return easygo.NewFailMsg("创建数量只能是1~99")
	}

	if reqMsg.Password == nil {
		return easygo.NewFailMsg("密码不能为空")
	}

	if reqMsg.Payword != nil && len(reqMsg.GetPayword()) != 6 {
		return easygo.NewFailMsg("支付密码长度错误")
	}

	// if reqMsg.Label == nil || len(reqMsg.GetLabel()) == 0 {
	// 	return easygo.NewFailMsg("兴趣标签不能为空")
	// }

	accode := 0
	account := 0
	accounts := ""
	reaccount := ""
	count := int(reqMsg.GetCount())
	ip := ep.GetConnection().RemoteAddr().String()
	if reqMsg.GetType() == 6 {
		brands := []string{"HUAWEI", "iPhone", "Xiaomi", "Meizu", "OPPO", "vivo", "HONOR", "samsung", "Redmi"}
		if reqMsg.ChannelNo == nil && reqMsg.GetChannelNo() == "" {
			return easygo.NewFailMsg("渠道号不能为空")
		}
		li := for_game.QueryOperationByNo(reqMsg.GetChannelNo())
		if li == nil {
			return easygo.NewFailMsg("渠道不存在")
		}
		jsonFile := easygo.YamlCfg.GetValueAsString("CLIENT_VERSION_DATA")
		versiondata := for_game.FindVersion(jsonFile)
		pass := easygo.If(reqMsg.GetPassword() == "", "1111", reqMsg.GetPassword())
		RandAddPlayer(count, li, versiondata, brands, easygo.AnytoA(pass))
	} else {
		for i := 1; i <= count; i++ {
			switch reqMsg.GetType() {
			case 2:
				accode = int(for_game.NextId("waiter_count"))
				account = 10010000000
			case 3:
				accode = int(for_game.NextId("user_shop_count"))
				account = 12010000000
			case 4:
				accode = int(for_game.NextId("manage_count"))
				account = 20010000000
			case 5:
				accode = int(for_game.NextId("user_official_count"))
				account = 11010000000
			case 6:
			default:
				return easygo.NewFailMsg("创建用户类型错误")
			}

			account += accode
			sAccount := easygo.IntToString(account)
			data := &share_message.CreateAccountData{
				Phone:    easygo.NewString(sAccount),
				PassWord: easygo.NewString(reqMsg.GetPassword()),
				Ip:       easygo.NewString(ip),
				Types:    easygo.NewInt32(reqMsg.GetType()),
			}
			b, a := for_game.CreateAccount(data)
			if b {
				if i < count {
					accounts += sAccount + ","
					if i == 1 {
						reaccount += sAccount + "~"
					}
				} else {
					accounts += sAccount
					reaccount += sAccount
				}

				if reqMsg.Payword != nil {
					// playeracc := QueryPlayerAccountbyId(a)
					playerAcc := for_game.GetRedisPlayerBase(a)
					if playerAcc == nil {
						continue
					}
					playerAcc.SetPayPassword(reqMsg.GetPayword())
					playerAcc.SaveToMongo()
				}

				player := for_game.GetRedisPlayerBase(a)
				player.SetNickName(for_game.GetRandNickName())
				if reqMsg.GetIsSlogan() {
					player.SetSignature(for_game.GetRandSignature())
				}
				sexType := reqMsg.GetSex()
				if sexType == -1 {
					sexType = int32(util.RandIntn(2) + 1)
				}
				if reqMsg.GetIsCity() {
					region := for_game.GetRandProvice()
					player.SetProvice(region)
					player.SetCity(for_game.GetRandCity(region))
				}
				player.SetDeviceType(3)
				player.SetCreateTime()
				if reqMsg.GetApprove() {
					player.SetPeopleId(sAccount)
					player.SetRealName("客服")
				}
				player.SetRedisLabelList(reqMsg.GetLabel())
				player.SetSex(int32(sexType))
				player.SetHeadIcon(for_game.GetRandRealHeadIcon(int(sexType)))
				player.SaveToMongo()
			}
		}
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, "批量创建运营用户"+accounts)

	msg := &brower_backstage.AddWaiterResponse{
		Account:  easygo.NewString(reaccount), //lg:10010000001~1001000011
		Password: reqMsg.Password,             //密码
		Payword:  reqMsg.Payword,              //支付密码 6位数字
		Approve:  reqMsg.Approve,              //是否自动实名认证
	}
	return msg
}

//冻结用户
func (self *cls4) RpcPlayerFreeze(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	overTime := reqMsg.GetObjIds()
	// if reqMsg.ObjIds == nil || len(overTime) == 0 {
	// 	return easygo.NewFailMsg("冻结结束时间不能为空")
	// }
	// if overTime[0] < easygo.NowTimestamp() {
	// 	return easygo.NewFailMsg("冻结结束时间不能小于当前时间")
	// }

	overTime = append(overTime, easygo.NowTimestamp()+60)

	req := &server_server.PlayerIds{
		PlayerIds:   reqMsg.Ids64,
		Node:        reqMsg.Note,
		Operator:    easygo.NewString(user.GetAccount()),
		BanOverTime: easygo.NewInt64(overTime[0]),
	}
	ChooseOneHall(0, "RpcFreezePlayer", req)

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i+1 < count {
			ids += QueryPlayerbyId(idsarr[i]).GetAccount() + ","

		} else {
			ids += QueryPlayerbyId(idsarr[i]).GetAccount()
		}
	}

	msg := fmt.Sprintf("批量冻结用户: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)

	AddPlayerFreezeLogs(reqMsg.GetIds64(), reqMsg.GetNote())
	return easygo.EmptyMsg
}

//解冻用户
func (self *cls4) RpcPlayerUnFreeze(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	UpPlayerStatus(reqMsg.GetIds64(), 0)
	for _, id := range reqMsg.GetIds64() {
		player := for_game.GetRedisPlayerBase(id)
		if player == nil {
			continue
		}
		player.SetStatus(0)
		player.SetNote("")
		player.SetOperator("")
		player.SetBanOverTime(0)
		for_game.DelFreezeAccount(player.GetAccount())
		player.SaveToMongo()
	}

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i+1 < count {
			ids += QueryPlayerbyId(idsarr[i]).GetAccount() + ","

		} else {
			ids += QueryPlayerbyId(idsarr[i]).GetAccount()
		}
	}

	msg := fmt.Sprintf("批量解冻用户: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)
	return easygo.EmptyMsg
}

//ID查询用户资料
func (self *cls4) RpcGetPlayerById(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.Id64 == nil || reqMsg.GetId64() == 0 {
		return easygo.NewFailMsg("ID不能为空")
	}
	pmrg := for_game.GetRedisPlayerBase(reqMsg.GetId64())
	player := pmrg.GetRedisPlayerBase()
	// player := QueryPlayerbyId(reqMsg.GetId64())
	if player == nil {
		return easygo.NewFailMsg("ID不存在")
	} else {
		regInfo := for_game.GetRegisterLoginLogById(player.GetPlayerId())
		player.RegType = easygo.NewInt32(regInfo.GetType())
	}

	//处理默认头像不显示问题
	if find := strings.Contains(player.GetHeadIcon(), "http"); !find {
		url := "https://im-resource-1253887233.cos.accelerate.myqcloud.com/defaulticon/"
		if player.GetHeadIcon() == "" {
			switch player.GetSex() {
			case 1:
				player.HeadIcon = easygo.NewString("boy_1")
			default:
				player.HeadIcon = easygo.NewString("girl_1")
			}

		}
		player.HeadIcon = easygo.NewString(url + player.GetHeadIcon() + ".png")
	}

	player.Diamond = easygo.NewInt64(0)
	wishP := for_game.GetRedisWishPlayerByImPid(player.GetPlayerId())
	if wishP != nil {
		player.Diamond = easygo.NewInt64(wishP.GetDiamond())
	}

	return player
}

//柠檬号或手机号查询用户资料
func (self *cls4) RpcGetPlayerByAccount(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.IdStr == nil || reqMsg.GetIdStr() == "" {
		return easygo.NewFailMsg("账号不能为空")
	}
	player := QueryPlayerbyAccount(reqMsg.GetIdStr())
	if player == nil {
		player = QueryPlayerbyPhone(reqMsg.GetIdStr())
		if player == nil {
			return easygo.NewFailMsg("账号不存在")
		}
	}

	switch player.GetStatus() {
	case for_game.ACCOUNT_USER_FROZEN, for_game.ACCOUNT_ADMIN_FROZEN: // 1 //用户冻结
		return easygo.NewFailMsg("玩家柠檬号已冻结!")
	case for_game.ACCOUNT_CANCELING: // 3 //注销中
		return easygo.NewFailMsg("玩家柠檬号注销中!")
	case for_game.ACCOUNT_CANCELED: // 4 //已注销
		return easygo.NewFailMsg("玩家柠檬号已注销!")
	}

	player.Diamond = easygo.NewInt64(0)
	wishP := for_game.GetRedisWishPlayerByImPid(player.GetPlayerId())
	if wishP != nil {
		player.Diamond = easygo.NewInt64(wishP.GetDiamond())
	}

	return player
}

//查询用户投诉
func (self *cls4) RpcQueryPlayerComplaint(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetPlayerComplaint(reqMsg)
	return &brower_backstage.PlayerComplaintResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询其他投诉
func (self *cls4) RpcQueryPlayerComplaintOther(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetPlayerComplaintOther(reqMsg)
	return &brower_backstage.PlayerComplaintResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//回复用户投诉
func (self *cls4) RpcReplyPlayerComplaint(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PlayerComplaint) easygo.IMessage {

	if reqMsg.ReContent == nil || reqMsg.GetReContent() == "" {
		return easygo.NewFailMsg("回复内容不能为空")
	}
	reqMsg.ReTime = easygo.NewInt64(time.Now().Unix())
	reqMsg.Operator = user.Account
	reqMsg.Status = easygo.NewInt32(2)

	ReplyPlayerComplaint(reqMsg)
	ReplyPlayerComplaintForHall(reqMsg) //回复用户投诉到大厅

	msg := fmt.Sprintf("处理投诉单号: %d", reqMsg.GetId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询添加用户列表
func (self *cls4) RpcQueryFriendPlayerList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.GetPlayerFriendListRequest) easygo.IMessage {

	list, count := GetFriendPlayerList(reqMsg)
	msg := &brower_backstage.GetPlayerListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//批量添加好友
func (self *cls4) RpcAddFriend(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.AddPlayerFriendInfo) easygo.IMessage {
	friend_base := for_game.GetFriendBase(reqMsg.GetPlayerID())

	if friend_base == nil {
		return easygo.NewFailMsg("找不到该用户")
	}
	friend_count := 0
	if friend_base != nil {
		friend_count = len(friend_base.GetFriendIds())
	}

	if friend_count+len(reqMsg.GetList()) > for_game.MAX_FRIEND {
		return easygo.NewFailMsg("好友数量达到上限")
	}

	SendToPlayer(reqMsg.GetPlayerID(), "RpcAddFriend", &server_server.AddPlayerFriendInfo{PlayerID: reqMsg.PlayerID, List: reqMsg.List})
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, "批量添加好友")
	return nil
}

//查询好友数量
func (self *cls4) RpcQueryPlayerInfo(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.GetIdStr() == "" {
		return easygo.NewFailMsg("柠檬号不能为空")
	}
	player := QueryPlayerbyAccount(reqMsg.GetIdStr())

	if player == nil {
		return easygo.NewFailMsg("该用户不存在")
	}
	friend_base := for_game.GetFriendBase(player.GetPlayerId())
	friend_count := 0
	if friend_base != nil {
		friend_count = len(friend_base.GetFriendIds())
	}

	return &brower_backstage.PlayerFriendInfo{FriendCount: easygo.NewInt32(friend_count), MaxFriendCount: easygo.NewInt32(for_game.MAX_FRIEND), Info: player}
}

//修改用户自定义标签
func (self *cls4) RpcEditPlayerCustomTag(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	EditPlayerCustomTag(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, "修改用户:"+easygo.AnytoA(reqMsg.GetIds64()[0])+"自定义标签")

	return easygo.EmptyMsg
}

//修改用户个性标签
func (self *cls4) RpcEditPersonalityTags(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	players := reqMsg.GetIds64()
	if len(players) < 1 {
		return easygo.NewFailMsg("用户id不能为空")
	}

	pMgr := for_game.GetRedisPlayerBase(players[0])
	v := reqMsg.GetIds32()
	if len(v) > 0 {
		easygo.SortSliceInt32(v, true)
		pMgr.SetPersonalityTags(v)
		pMgr.SaveOneRedisDataToMongo("PersonalityTags", v)
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, "修改用户:"+easygo.AnytoA(reqMsg.GetIds64()[0])+"个性标签")

	return easygo.EmptyMsg
}

//注销账号记录列表
func (self *cls4) RpcPlayerCancleAccountList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetPlayerCancleAccountList(reqMsg)
	return &brower_backstage.PlayerCancleAccountListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//审核注销账号
func (self *cls4) RpcEditPlayerCancleAccount(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PlayerCancleAccount) easygo.IMessage {
	logs.Info("RpcEditPlayerCancleAccount:", reqMsg)
	if reqMsg.Note == nil || reqMsg.GetNote() == "" {
		return easygo.NewFailMsg("原因不能为空")
	}
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("Id不能为空")
	}
	EditPlayerCancleAccount(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, "审核:"+string(reqMsg.GetAccount())+"账号注销")

	return easygo.EmptyMsg
}

//查询个人聊天记录
func (self *cls4) RpcQueryPersonalChatLog(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ChatLogRequest) easygo.IMessage {
	logs.Debug("=========RpcQueryPersonalChatLog=========", reqMsg)
	if (reqMsg.Keyword1 == nil && reqMsg.Keyword2 == nil) || (reqMsg.GetKeyword1() == "" && reqMsg.GetKeyword2() == "") {
		return &brower_backstage.PersonalChatLogResponse{}
	}
	for_game.SavePersonalChatToMongoDB()
	list, count := GetPersonalChatLogList(reqMsg)

	return &brower_backstage.PersonalChatLogResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询指定个人聊天记录
func (self *cls4) RpcQueryPersonalChatLogByObj(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ChatLogRequest) easygo.IMessage {
	if reqMsg.Keyword1 == nil || reqMsg.Keyword2 == nil {
		return easygo.NewFailMsg("指定的聊天对象不能为空")
	}
	for_game.SavePersonalChatToMongoDB()
	list, count, ids := QueryPersonalChatLogByObj(reqMsg)

	players := QueryplayerlistByIds(ids)
	playerMap := make(map[int64]*share_message.PlayerBase)
	for _, p := range players {
		playerMap[p.GetPlayerId()] = p
	}

	for i, item := range list {
		list[i].TalkerAccount = easygo.NewString(playerMap[item.GetTalker()].GetAccount())
		list[i].TalkerNickName = easygo.NewString(playerMap[item.GetTalker()].GetNickName())

		list[i].TargetAccount = easygo.NewString(playerMap[item.GetTargetId()].GetAccount())
		list[i].TargetNickName = easygo.NewString(playerMap[item.GetTargetId()].GetNickName())
	}
	return &brower_backstage.PersonalChatLogResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//查询群聊天记录
func (self *cls4) RpcQueryTeamChatLog(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ChatLogRequest) easygo.IMessage {
	list, count := GetTeamChatLog(reqMsg)
	var ids, tids []int64
	for _, li := range list {
		ids = append(ids, li.GetTalker())
		tids = append(tids, li.GetTeamId())
	}
	players := QueryplayerlistByIds(ids)
	playerMap := make(map[int64]*share_message.PlayerBase)
	for _, p := range players {
		playerMap[p.GetPlayerId()] = p
	}

	teams := QueryTeambyIds(tids)
	teamMap := make(map[int64]*share_message.TeamData)
	for _, t := range teams {
		teamMap[t.GetId()] = t
	}

	for i, item := range list {
		list[i].TalkerAccount = easygo.NewString(playerMap[item.GetTalker()].GetAccount())
		list[i].TeamAccount = easygo.NewString(teamMap[item.GetTeamId()].GetTeamChat())
		list[i].TeamName = easygo.NewString(teamMap[item.GetTeamId()].GetName())
	}
	return &brower_backstage.TeamChatLogResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//检查用户报名单
func (self *cls4) RpcCheckChatLogWhitelist(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ChatLogRequest) easygo.IMessage {
	var ids PLAYER_IDS
	var checkReturn bool
	switch reqMsg.GetTypes() {
	case 1:
		ids = append(ids, easygo.StringToInt64noErr(reqMsg.GetKeyword1()), easygo.StringToInt64noErr(reqMsg.GetKeyword2()))
		players := CheckUserWhitelist(ids)
		if len(players) > 0 {
			checkReturn = true
		}
	case 2:
		id := easygo.StringToInt64noErr(reqMsg.GetKeyword1())
		mMgr := for_game.GetRedisTeamPersonalObj(id)
		if mMgr == nil {
			return easygo.NewFailMsg("群ID错误")
		}
		members := mMgr.GetRedisTeamPersonal()
		for _, item := range members {
			ids = append(ids, item.GetPlayerId())
		}

		players := CheckUserWhitelist(ids)
		if len(players) > 0 {
			checkReturn = true
		}

	}

	return &brower_backstage.CommonResponse{
		BoolField: easygo.NewBool(checkReturn),
	}
}
