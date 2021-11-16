// 大厅服务器为[游戏客户端]提供的服务

package hall

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"

	"github.com/astaxie/beego/logs"
)

type ServiceForGameClient struct {
}
type cls1 = ServiceForGameClient

func init() {
	RegisterServiceForGameClient("HallService", &ServiceForGameClient{})
}

//=============================mainlogic.proto=============================

// 登录大厅都用这一个
func (self *cls1) RpcLogin(ep IGameClientEndpoint, ctx interface{}, reqMsg *client_hall.LoginMsg, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcLogin,msg=%v", reqMsg) // 别删，永久留存
	playerId := reqMsg.GetPlayerId()
	Info := for_game.GetRedisAccountObj(playerId)
	if Info == nil {
		res := "玩家对象没有创建成功"
		logs.Error("玩家对象没有创建成功:")
		return easygo.NewFailMsg(res)
	}
	base := for_game.GetRedisPlayerBase(playerId)
	oldEp := ClientEpMp.LoadEndpoint(playerId)
	if for_game.GetMillSecond()-base.GetLastOnLineTime() < 500 && oldEp != nil {
		res := "玩家登录过于频繁,请稍后再试"
		logs.Error("玩家登录过于频繁:", for_game.GetMillSecond(), base.GetLastOnLineTime())
		return easygo.NewFailMsg(res, for_game.FAIL_MSG_CODE_1014)
	}
	token := reqMsg.GetToken()
	if token != base.GetToken() {
		res := "token验证码错误"
		logs.Error("顶号:", res, token, base.GetToken(), playerId)
		return easygo.NewFailMsg(res, for_game.FAIL_MSG_CODE_1001)
	}
	if base.GetStatus() == for_game.ACCOUNT_ADMIN_FROZEN || base.GetStatus() == for_game.ACCOUNT_USER_FROZEN {
		return easygo.NewFailMsg("账号已被冻结，登录失败")
	}
	if base.GetStatus() == for_game.ACCOUNT_CANCELING {
		return easygo.NewFailMsg("账号注销中，继续登录将取消注销", for_game.FAIL_MSG_CODE_1013)
	}
	if PlayerOnlineMgr.CheckPlayerIsOnLine(playerId) { //如果玩家在线
		serverId := PlayerOnlineMgr.GetPlayerServerId(playerId)
		if serverId != PServerInfo.GetSid() { //如果不是在当前大厅 要通知那个大厅 把他踢下线
			otherHall := PServerInfoMgr.GetServerInfo(serverId)
			if otherHall != nil {
				msg := &server_server.PlayerIdInfo{
					PlayerId: easygo.NewInt64(playerId),
				}
				m, err := SendMsgToServerNew(otherHall.GetSid(), "RpcOtherHallLogin", msg)
				if err != nil {
					logs.Error("登录顶号失败:", err)
				}
				if m1, ok := m.(*server_server.PlayerIdInfo); ok {
					if m1.GetSuccess() == false {
						logs.Error("登陆顶号失败", playerId, serverId)
					}
				}
			} else {
				logs.Error("获取大厅失败，大厅id:", serverId)
			}
		} else {
			if oldEp != nil && oldEp != ep {
				logs.Error("玩家顶号通知:", playerId)
				oldEp.RpcReLogin(nil) //通知前端被顶号了，返回登录界面
				oldEp.SetFlag(true)
				oldEp.Shutdown()
			}
		}
		Info = for_game.GetRedisAccountObj(playerId) //shutdown 会把account的redis数据清空  所有要重新生成
	}
	//shopSrv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_SHOP)
	//if shopSrv == nil {
	//	ep.Shutdown()
	//	logs.Info("登录大厅时，找不到有效的商场服务器")
	//	return nil
	//}
	player := PlayerMgr.LoadPlayer(playerId)
	if player == nil {
		player = NewPlayer(playerId)
		player.OnLoadFromDB()
		PlayerMgr.Store(playerId, player)
	} else {
		player.OnLoadFromDB()
	}
	easygo.Spawn(func() {
		for_game.MakePlayerKeepReport(playerId, 1)
	}) //生成玩家留存报表 已优化到Redis

	//设置不是后台
	if player.GetLoginTimes() == 1 || (player.GetLoginTimes() != 1 && for_game.GetMillSecond()-player.GetLastLogOutTime() > 15*86400*1000) { //如果不是第一次登陆并且离上次登陆超过15天
		SendLoginMessage(player.GetPlayerId(), player.GetNickName())
	}
	player.CheckRegistrationIdOrChannel(reqMsg.GetRegistrationId(), reqMsg.GetChannel())
	//关联玩家ep前先检测有没有通话中
	callDetail := player.GetDetailCallInfo()
	if callDetail != nil && callDetail.GetOperate() != 0 {
		//帮玩家挂断通话
		callDetail.Operate = easygo.NewInt32(3)
		callDetail.OperateId = easygo.NewInt64(player.GetPlayerId())
		self.RpcOperateSpecialChat(ep, player, callDetail)
		logs.Info("登录时清理通话状态------------->>>>>>", player.GetPlayerId())
	}
	ep.SetAssociativePlayer(player) //关联玩家ep
	player.SetDeviceType(reqMsg.GetDeviceType())
	player.SetVersion(reqMsg.GetVersionNumber())
	player.SetBrand(reqMsg.GetBrand())
	msg := player.GetAllPlayerInfo(reqMsg.GetLoginType())
	callInfo := player.GetCallInfo()
	if callInfo.GetPlayerId() != 0 { //如果登陆的时候 有人在给你打电话  下发通话协议
		smsg := callInfo.GetSMsg()
		var m *client_hall.SpecialChatInfo
		_ = json.Unmarshal([]byte(smsg), &m)
		ep.RpcRequestSpecialChatResponse(m)
	}
	ep.RpcPlayerLoginResponse(msg)
	//存储玩家节点
	ClientEpMp.StoreEndpoint(playerId, ep.GetEndpointId())
	//设置玩家为上线状态
	PlayerOnlineMgr.PlayerOnline(playerId, PServerInfo.GetSid())
	//记录玩家当前所在大厅
	player.SetSid(PServerInfo.GetSid())
	//登录成功，通知其他服务器玩家进入
	NotifyPlayerOnLine(playerId)

	easygo.Spawn(func() {
		if easygo.Get0ClockTimestamp(player.GetLastLogOutTime()) < easygo.GetToday0ClockTimestamp() {
			ip := ep.GetAddr().String()
			for_game.MakePlayerLogLocationReport(reqMsg.GetDeviceType(), ip)
		}
		for_game.AddStatisticsInfo(reqMsg.GetType(), playerId, reqMsg.GetLoginType())
	})
	waiterMsg := for_game.QueryIMmessageByPid(playerId, for_game.WAITER_MESSAGE_ING)
	if waiterMsg.GetCnew() > 0 {
		im := &share_message.IMmessage{}
		im.Id = easygo.NewInt64(waiterMsg.GetId())
		im.Cnew = easygo.NewInt32(waiterMsg.GetCnew())
		ep.RpcNewWaiterMsg(im)
	}
	logs.Info("玩家登录完成", player.GetPlayerId(), ep)
	return reqMsg
}

//获取玩家其他数据信息
//func (self *cls1) RpcGetPlayerOtherData(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
//	logs.Info("请求玩家其他数据,pid----->", who.GetPlayerId())
//	player := GetPlayerObj(who.GetPlayerId())
//	msg := &client_server.AllPlayerMsg{}
//	if player != nil {
//		myChatLogs := for_game.GetUnReadeChatMessage(player.GetPlayerId())
//		//logs.Info("我的聊天记录:", myChatLogs, who.GetPlayerId())
//		msg = &client_server.AllPlayerMsg{
//			Friends:     player.GetFriendsInfo(),
//			Teams:       player.GetTeamsInfo(),
//			ChatMsg:     myChatLogs,
//			Pay:         PPayChannelMgr.GetPlatformChannelList(),
//			PayConfig:   PPayChannelMgr.GetPaymentSettingList(),                     //暂时默认用鹏聚的提现配置
//			LimitConfig: PSysParameterMgr.GetSysParameter(for_game.LIMIT_PARAMETER), //写死默认的
//			Types:       easygo.NewInt32(player.GetTypes()),
//		}
//	}
//	logs.Info("请求玩家其他数据完成,pid----->", who.GetPlayerId())
//	//检测玩家有没有新道具
//	bagItemObj := for_game.GetRedisPlayerBagItemObj(player.GetPlayerId())
//	if bagItemObj != nil {
//		types := bagItemObj.CheckNewItemNotice()
//		if len(types) > 0 {
//			//通知前端有新道具
//			notice := &client_hall.NewBagItemsTip{Types: types}
//			ep.RpcNewBagItemsTip(notice)
//		}
//	}
//	return msg
//}

//获取玩家其他数据信息
func (self *cls1) RpcGetPlayerOtherDataNew(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetPlayerOtherDataNew,pid----->", who.GetPlayerId())
	//player := GetPlayerObj(who.GetPlayerId())
	player := for_game.GetRedisPlayerBase(who.GetPlayerId())
	if player == nil {
		return easygo.NewFailMsg("无效的玩家数据")
	}
	msg := &client_server.AllPlayerMsg{
		IsNearBy:    easygo.NewBool(player.GetIsNearBy()),
		Pay:         PPayChannelMgr.GetPlatformChannelList(),
		PayConfig:   PPayChannelMgr.GetPaymentSettingList(),                     //暂时默认用鹏聚的提现配置
		LimitConfig: PSysParameterMgr.GetSysParameter(for_game.LIMIT_PARAMETER), //写死默认的
	}

	if len(player.GetLabelList()) == 0 { //代表没有完善标签信息  第一次登陆
		msgList := for_game.GetInterestTagAllList()
		sysInfo := for_game.QuerySysParameterById(for_game.INTEREST_PARAMETER)
		m := &client_server.LabelMsg{
			LabelInfo:    msgList,
			Max:          easygo.NewInt32(sysInfo.GetInterestMax()),
			Min:          easygo.NewInt32(sysInfo.GetInterestMin()),
			InterestType: for_game.GetRedisLabelInfo(),
		}
		msg.LabelMsg = m
		msg.RandName = easygo.NewString(for_game.GetRandNickName()) //  随机昵称库
	}

	if player.GetIsRecommendOver() {
		msg.RecommendInfo = for_game.GetRecommendInfo(player.GetPlayerId(), 0)
	}

	//赞相关信息
	msg.FanNum = easygo.NewInt32(len(player.GetFans()))
	msg.AttentionNum = easygo.NewInt32(len(player.GetAttention()))
	msg.ZanNum = easygo.NewInt32(player.GetZan() + for_game.GetAllTrueZan(player.GetPlayerId()))

	bagItemObj := for_game.GetRedisPlayerBagItemObj(player.GetPlayerId())
	if bagItemObj != nil {
		types := bagItemObj.CheckNewItemNotice()
		if len(types) > 0 {
			//通知前端有新道具
			notice := &client_hall.NewBagItemsTip{Types: types}
			ep.RpcNewBagItemsTip(notice)
		}
	}
	return msg
}

//获取玩家好友信息
func (self *cls1) RpcGetPlayerFriends(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetPlayerFriends,pid----->", who.GetPlayerId())
	player := GetPlayerObj(who.GetPlayerId())
	msg := &client_server.AllPlayerMsg{}
	if player != nil {
		msg = &client_server.AllPlayerMsg{
			Friends: player.GetFriendsInfo(),
		}
	}
	return msg
}

//获取离线期间新加的好友
func (self *cls1) RpcGetPlayerNewFriends(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.GetNewFriends, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetPlayerNewFriends,pid----->", who.GetPlayerId())
	msg := &client_server.NewFriends{
		Friends: who.GetFriendsInfo(reqMsg.GetTime()),
	}
	logs.Info("返回 RpcGetPlayerNewFriends 列表:", msg)
	return msg
}

//请求新好友数目
func (self *cls1) RpcNewVersionGetFriendNumRequest(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetPlayerFriends,pid----->", who.GetPlayerId())
	playerId := who.GetPlayerId()
	player := GetPlayerObj(playerId)
	num := 0
	if player != nil {
		fq := for_game.GetFriendBase(playerId)
		req := fq.GetNewVersionAllFriendRequestForOne()
		num = len(req.GetAddPlayerRequest())
	}
	msg := &client_hall.FriendNum{
		Num: easygo.NewInt32(num),
	}
	logs.Info("RpcGetPlayerFriends,pid----->", who.GetPlayerId(), num)
	return msg
}

//退出登录
func (self *cls1) RpcLogOut(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	playerId := who.GetPlayerId()
	//ClientEpMgr.Delete(ep.GetEndpointId())
	//ClientEpMp.Delete(playerId)
	//who.UpdateSaveChatLog(true) //离线前把聊天记录保存
	//HandleAfterLogout(ep, who)
	PlayerMgr.Delete(playerId)
	return nil
}

//添加好友
func (self *cls1) RpcAddFriend(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.AddPlayerInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("添加好友信息：", reqMsg, who.GetPlayerId())

	t := reqMsg.GetType()
	isMark := reqMsg.GetMark()
	var friendId int64
	var base *for_game.RedisPlayerBaseObj

	switch t {
	case share_message.AddFriend_Type_PHONE: //搜索账号添加好友
		//BankIds := who.GetBankId()
		//RealName := who.GetRealName()
		//if len(BankIds) == 0 || RealName == "" {
		//	res := "绑定银行卡或者真实姓名没有填写完全，不能添加好友"
		//	return easygo.NewFailMsg(res)
		//}
		account := reqMsg.GetAccount()
		if account == who.GetAccount() {
			return easygo.NewFailMsg("不可以添加自己为好友")
		}
		Info := for_game.GetRedisAccountByPhone(account)
		if Info == nil {
			return easygo.NewFailMsg("该玩家账号不存在")
		}

		friendId = Info.GetPlayerId()
		base = for_game.GetRedisPlayerBase(friendId)
		if base == nil || base.GetStatus() == for_game.ACCOUNT_CANCELED {
			res := "该账号异常"
			return easygo.NewFailMsg(res)
		}
		if !base.GetIsPhone() {
			return easygo.NewFailMsg("无法找到该用户")
		}
	case share_message.AddFriend_Type_TEAM:
		friendId = reqMsg.GetPlayerId()
		base = for_game.GetRedisPlayerBase(friendId)
		if base == nil || base.GetStatus() == for_game.ACCOUNT_CANCELED {
			res := "该账号异常"
			return easygo.NewFailMsg(res)
		}
		if !base.GetIsTeamChat() {
			return easygo.NewFailMsg("由于对方的隐私设置，你无法通过群聊添加对方")
		}
	case share_message.AddFriend_Type_CODE:
		friendId = reqMsg.GetPlayerId()
		base = for_game.GetRedisPlayerBase(friendId)
		if base == nil || base.GetStatus() == for_game.ACCOUNT_CANCELED {
			res := "该账号异常"
			return easygo.NewFailMsg(res)
		}
		if !base.GetIsCode() {
			return easygo.NewFailMsg("由于对方的隐私设置，你无法通过二维码添加对方")
		}
	case share_message.AddFriend_Type_ACCOUNT:
		account := reqMsg.GetAccount()
		friendId = for_game.GetPlayerIdForAccount(account)
		if friendId == 0 {
			return easygo.NewFailMsg("该玩家账号不存在")
		}
		base = for_game.GetRedisPlayerBase(friendId)
		if base == nil || base.GetStatus() == for_game.ACCOUNT_CANCELED {
			res := "该账号异常"
			return easygo.NewFailMsg(res)
		}
		if !base.GetIsAccount() {
			return easygo.NewFailMsg("无法找到该用户")
		}
	case share_message.AddFriend_Type_PLAYERID, share_message.AddFriend_Type_SQUARE, share_message.AddFriend_Type_NEARBY:
		friendId = reqMsg.GetPlayerId()
		base = for_game.GetRedisPlayerBase(friendId)
		if base == nil || base.GetStatus() == for_game.ACCOUNT_CANCELED {
			res := "该账号异常"
			return easygo.NewFailMsg(res)
		}
	case share_message.AddFriend_Type_CARD:
		friendId = reqMsg.GetPlayerId()
		base = for_game.GetRedisPlayerBase(friendId)
		if base == nil || base.GetStatus() == for_game.ACCOUNT_CANCELED {
			res := "该账号异常"
			return easygo.NewFailMsg(res)
		}
		if !base.GetIsCard() {
			return easygo.NewFailMsg("由于对方的隐私设置，你无法通过名片添加对方")
		}
	case share_message.AddFriend_Type_STRANGER:
		friendId = reqMsg.GetPlayerId()
		base = for_game.GetRedisPlayerBase(friendId)
		if base == nil || base.GetStatus() == for_game.ACCOUNT_CANCELED {
			res := "该账号异常"
			return easygo.NewFailMsg(res)
		}
	case share_message.AddFriend_Type_REGISTER:
		friendId = reqMsg.GetPlayerId()
		base = for_game.GetRedisPlayerBase(friendId)
		if base == nil || base.GetStatus() == for_game.ACCOUNT_CANCELED {
			res := "该账号异常"
			return easygo.NewFailMsg(res)
		}

	default:
		logs.Error("添加好友传过来的类型有误,type: ", t)
		return easygo.NewFailMsg("添加好友类型错误.")
	}

	if friendId == 0 {
		panic("玩家id 怎么会是0")
	}

	if friendId == who.GetPlayerId() {
		return easygo.NewFailMsg("不可以添加自己为好友")
	}

	if util.Int64InSlice(friendId, who.GetFriends()) {
		return easygo.NewFailMsg("该玩家已经是你的好友")
	}

	//if base == nil {
	//	base = GetPlayerObj(friendId)
	//}
	if util.Int64InSlice(who.GetPlayerId(), base.GetBlackList()) {
		return easygo.NewFailMsg("玩家拒绝接受你的消息")
	}
	wq := for_game.GetFriendBase(who.GetPlayerId())
	if !wq.CheckMaxNum() {
		return easygo.NewFailMsg("你的好友数量已经达到上限")
	}
	photo := base.GetHeadIcon()
	//photoList := base.GetPhoto()
	//if len(photoList) != 0 {
	//	photo = photoList[0]
	//}

	var isAss bool
	fq := for_game.GetFriendBase(friendId)
	friend := for_game.GetRedisPlayerBase(friendId)
	newReqMsg := &share_message.AddPlayerRequest{}
	if !util.Int64InSlice(who.GetPlayerId(), friend.GetFriends()) {
		nextId := for_game.NextId("Add_Friend")
		newReqMsg = &share_message.AddPlayerRequest{
			PlayerId:  easygo.NewInt64(who.GetPlayerId()),
			Time:      easygo.NewInt64(time.Now().Unix()),
			Text:      easygo.NewString(reqMsg.GetText()),
			Type:      &t,
			Result:    easygo.NewInt32(for_game.ADDFRIEND_NODEAL), //未处理
			Id:        easygo.NewInt64(nextId),
			IsRead:    easygo.NewBool(false),
			Photo:     easygo.NewString(photo),
			Signature: easygo.NewString(base.GetSignature()),
			Sex:       easygo.NewInt32(who.GetSex()),
		}
		fq.AddFriendRequest(newReqMsg)
		isAss = true
	}
	//如果对方已经把你添加好友了
	if util.Int64InSlice(who.GetPlayerId(), friend.GetFriends()) {
		logs.Info("已经是对方添加好友:", who.GetPlayerId(), friendId)
		fq.AgreeAddFriend(who.GetPlayerId()) // 同一个大厅并且在线
		wq.AddFriend(friendId, int32(t))
		msg1 := GetFriendInfo(friendId)
		msg1.OpenWindows = easygo.NewInt32(1)
		ep.RpcNoticeAgreeFriend(msg1)
	} else {
		//对方没把我加为好友
		if friend.GetIsAddFriend() { //开启好友验证
			wq.AddFriend(friendId, int32(t))
			who.AddAttention(friendId)
			who.AddFans(friendId)
			for_game.OperateRedisDynamicAttention(1, who.GetPlayerId(), friendId)
			//msg := fq.GetNewVersionAllFriendRequestForOne() //发送所有好友申请信息
			//SendMsgToHallClientNew(friendId, "RpcNoticeAddFriend", msg)
			newReqMsg.NickName = easygo.NewString(who.GetNickName())
			newReqMsg.HeadIcon = easygo.NewString(who.GetHeadIcon())
			SendMsgToHallClientNew(friendId, "RpcNoticeAddFriendNew", newReqMsg)
			msg1 := GetFriendInfo(friendId)
			ep.RpcNoticeAgreeFriend(msg1)
		} else {
			if !fq.CheckMaxNum() {
				return nil
			}
			//fq.ReadFriendRequest([]int64{who.GetPlayerId()})
			fq.AgreeAddFriend(who.GetPlayerId()) // 同一个大厅并且在线
			wq.AddFriend(friendId, int32(t))
			fq.AddFriend(who.GetPlayerId(), int32(t))

			msg := GetFriendInfo(who.GetPlayerId())
			msg.IsMark = easygo.NewBool(isMark)
			msg.AddType = &t
			SendMsgToHallClientNew(friendId, "RpcAgreeFriendResponse", msg)

			msg1 := GetFriendInfo(friendId)
			msg1.OpenWindows = easygo.NewInt32(1)
			ep.RpcNoticeAgreeFriend(msg1)
			who.AddAttention(friendId)
			friend.AddFans(who.GetPlayerId())
			friend.AddAttention(who.GetPlayerId())
			who.AddFans(friendId)
			for_game.OperateRedisDynamicAttention(1, who.GetPlayerId(), friendId)
			for_game.OperateRedisDynamicAttention(1, friendId, who.GetPlayerId())
		}
	}
	//if PlayerOnlineMgr.CheckPlayerIsOnLine(friendId) { //如果在线
	//	serverId := PlayerOnlineMgr.GetPlayerServerId(friendId)
	//	if serverId == PServerInfo.GetSid() { //在同一个大厅
	//		if ep1 := ClientEpMp.LoadEndpoint(friendId); ep1 != nil {
	//			if util.Int64InSlice(who.GetPlayerId(), friend.GetFriends()) {
	//				fq.AgreeAddFriend(who.GetPlayerId()) // 同一个大厅并且在线
	//				wq.AddFriend(friendId, int32(t))
	//				msg1 := GetFriendInfo(friendId)
	//				msg1.OpenWindows = easygo.NewInt32(1)
	//				ep.RpcNoticeAgreeFriend(msg1)
	//			} else {
	//				if friend.GetIsAddFriend() { //开启好友验证
	//					msg := fq.GetNewVersionAllFriendRequestForOne() //发送所有好友申请信息
	//					ep1.RpcNoticeAddFriend(msg)
	//				} else {
	//					if !fq.CheckMaxNum() {
	//						return nil
	//					}
	//					//fq.ReadFriendRequest([]int64{who.GetPlayerId()})
	//					fq.AgreeAddFriend(who.GetPlayerId()) // 同一个大厅并且在线
	//					wq.AddFriend(friendId, int32(t))
	//					fq.AddFriend(who.GetPlayerId(), int32(t))
	//
	//					msg := GetFriendInfo(who.GetPlayerId())
	//					msg.IsMark = easygo.NewBool(isMark)
	//					msg.AddType = &t
	//					ep1.RpcAgreeFriendResponse(msg)
	//
	//					msg1 := GetFriendInfo(friendId)
	//					msg1.OpenWindows = easygo.NewInt32(1)
	//					ep.RpcNoticeAgreeFriend(msg1)
	//					who.AddAttention(friendId)
	//					friend.AddFans(who.GetPlayerId())
	//					friend.AddAttention(who.GetPlayerId())
	//					who.AddFans(friendId)
	//					for_game.OperateRedisDynamicAttention(1, who.GetPlayerId(), friendId)
	//					for_game.OperateRedisDynamicAttention(1, friendId, who.GetPlayerId())
	//				}
	//			}
	//		}
	//	} else { //在线但是不在同一个大厅
	//		otherHall := PServerInfoMgr.GetServerInfo(serverId)
	//		if otherHall != nil {
	//			if util.Int64InSlice(who.GetPlayerId(), fq.GetFriendIds()) {
	//				fq.AgreeAddFriend(who.GetPlayerId())
	//				wq.AddFriend(friendId, int32(t))
	//				msg1 := GetFriendInfo(friendId)
	//				msg1.OpenWindows = easygo.NewInt32(1)
	//				ep.RpcNoticeAgreeFriend(msg1)
	//			} else {
	//				if base.GetIsAddFriend() { //开启好友验证
	//					msg := fq.GetNewVersionAllFriendRequestForOne() //发送所有好友申请信息
	//					SendMsgToHallClientNew(friendId, "RpcNoticeAddFriend", msg)
	//				} else {
	//					if !fq.CheckMaxNum() {
	//						return nil
	//					}
	//
	//					fq.AgreeAddFriend(who.GetPlayerId()) // 同一个大厅并且在线
	//					wq.AddFriend(friendId, int32(t))
	//					fq.AddFriend(who.GetPlayerId(), int32(t))
	//					msg := GetFriendInfo(who.GetPlayerId())
	//					msg.IsMark = easygo.NewBool(isMark)
	//					msg.AddType = &t
	//					SendMsgToHallClientNew(friendId, "RpcAgreeFriendResponse", msg)
	//
	//					msg1 := GetFriendInfo(friendId)
	//					msg1.OpenWindows = easygo.NewInt32(1) //打开聊天窗口
	//					ep.RpcNoticeAgreeFriend(msg1)
	//					//社交广场互粉+互关注
	//					who.AddAttention(friendId)
	//					who.AddFans(friendId)
	//					friend.AddAttention(who.GetPlayerId())
	//					friend.AddFans(who.GetPlayerId())
	//					for_game.OperateRedisDynamicAttention(1, who.GetPlayerId(), friendId)
	//					for_game.OperateRedisDynamicAttention(1, friendId, who.GetPlayerId())
	//				}
	//			}
	//		} else {
	//			logs.Info("多大厅获取ep错误，", friendId, serverId)
	//		}
	//	}
	//} else {
	//	//logs.Info("玩家不在线---------", fq)
	//	if !base.GetIsAddFriend() || util.Int64InSlice(who.GetPlayerId(), fq.GetFriendIds()) {
	//		fq.AddFriend(who.GetPlayerId(), int32(t))
	//		wq.AddFriend(friendId, int32(t))
	//		fq.ReadFriendRequest([]int64{who.GetPlayerId()})
	//		msg := GetFriendInfo(friendId)
	//		msg.AddType = &t
	//		msg.IsMark = easygo.NewBool(isMark)
	//		ep.RpcAgreeFriendResponse(msg)
	//		//社交广场互粉+互关注
	//		who.AddAttention(friendId)
	//		who.AddFans(friendId)
	//		friend.AddAttention(who.GetPlayerId())
	//		friend.AddFans(who.GetPlayerId())
	//		for_game.OperateRedisDynamicAttention(1, who.GetPlayerId(), friendId)
	//		for_game.OperateRedisDynamicAttention(1, friendId, who.GetPlayerId())
	//	}
	//}
	if isAss {
		NoticeAssistant2(who, friendId, int32(t), 0, 0)
	}
	//发送打招呼
	content := base64.StdEncoding.EncodeToString([]byte("Hi,你好吗?"))
	sessionId := for_game.MakeSessionKey(friendId, who.GetPlayerId())
	chatMsg := &share_message.Chat{
		SessionId:   easygo.NewString(sessionId),
		SourceId:    easygo.NewInt64(who.GetPlayerId()),
		TargetId:    easygo.NewInt64(friendId),
		Content:     easygo.NewString(content),
		ChatType:    easygo.NewInt32(for_game.CHAT_TYPE_PRIVATE),
		ContentType: easygo.NewInt32(for_game.TALK_CONTENT_WORD),
		SayType:     easygo.NewInt32(int32(t)),
	}
	self.RpcChatNew(nil, who, chatMsg)

	return nil
}

//阅读好友请求消息
func (self *cls1) RpcReadFriendRequest(ep IGameClientEndpoint, who *Player, reqMsg *client_server.ReadInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("=====RpcReadFriendRequest====, reqMst=%v", reqMsg)
	ids := reqMsg.GetFriendId()
	fq := for_game.GetFriendBase(who.GetPlayerId())
	fq.ReadFriendRequest(ids)
	return nil
}

// 查找好友 优化后
func (self *cls1) RpcFindFriend(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.AccountInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("搜索好友信息：", reqMsg)
	account := reqMsg.GetAccount()
	t := reqMsg.GetType()
	var pid PLAYER_ID

	var player *Player
	switch t {
	case share_message.AddFriend_Type_PHONE: // 手机号码搜索
		info := for_game.GetRedisAccountByPhone(account)
		if info == nil {
			return easygo.NewFailMsg("该用户不存在")
		}
		pid = info.GetPlayerId()
		player = GetPlayerObj(pid)
		if player != nil && !player.GetIsPhone() {
			res := "无法找到该用户"
			return easygo.NewFailMsg(res)
		}
	case share_message.AddFriend_Type_ACCOUNT: // 柠檬号搜索
		pid = for_game.GetPlayerIdForAccount(account)
		if pid == 0 {
			res := fmt.Sprintf("该用户不存在")
			return easygo.NewFailMsg(res)
		}
		player = GetPlayerObj(pid)
		if player != nil && !player.GetIsAccount() {
			res := "无法找到该用户"
			return easygo.NewFailMsg(res)
		}
	default:
		logs.Error("查找好友,查找类型有误,type: ", t)
		return easygo.NewFailMsg("查找类型有误")
	}

	// 判断是否注销账号
	if player == nil || player.GetStatus() == for_game.ACCOUNT_CANCELED {
		return easygo.NewFailMsg("该用户不存在")
	}

	if util.Int64InSlice(pid, who.GetFriends()) {
		res := "该玩家已经是你的好友"
		return easygo.NewFailMsg(res)
	}
	msg := GetFriendInfo(pid)
	msg.AddType = &t
	ep.RpcFindFriendResponse(msg)
	return nil
}

//同意添加好友
func (self *cls1) RpcAgreeFriend(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.AccountInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("同意添加好友信息：", reqMsg, who.GetPlayerId())
	t := reqMsg.GetType()
	friendId := reqMsg.GetPlayerId() // 666
	if t == share_message.AddFriend_Type_NEARBY {
		for_game.AgreeNearByInfo(who.GetPlayerId(), friendId)
	}

	if friendId == 0 {
		panic("玩家id怎么会是0")
	}

	wq := for_game.GetFriendBase(who.GetPlayerId())
	if util.Int64InSlice(friendId, who.GetFriends()) {
		res := "该玩家已经是你的好友"
		return easygo.NewFailMsg(res)
	}

	if !wq.CheckMaxNum() {
		res := "你的好友数量已经到达上限，无法加入"
		return easygo.NewFailMsg(res)
	}

	fq := for_game.GetFriendBase(friendId)
	if !fq.CheckMaxNum() {
		res := "对方好友数量已经到达上限，无法加入"
		return easygo.NewFailMsg(res)
	}

	NoticeAssistant2(who, friendId, int32(t), 1, 0)
	wq.AgreeAddFriend(friendId)
	wq.AddFriend(friendId, int32(t))
	AgreeAddFriend(who.GetPlayerId(), friendId, int32(t))
	msg := GetFriendInfo(friendId)
	msg.OpenWindows = easygo.NewInt32(1)
	msg.AddType = &t
	//社交广场相互关注+粉
	who.AddAttention(friendId)
	who.AddFans(friendId)
	friend := for_game.GetRedisPlayerBase(friendId)
	friend.AddAttention(who.GetPlayerId())
	friend.AddFans(who.GetPlayerId())
	for_game.OperateRedisDynamicAttention(1, who.GetPlayerId(), friendId)
	for_game.OperateRedisDynamicAttention(1, friendId, who.GetPlayerId())
	ep.RpcAgreeFriendResponse(msg)
	return nil
}

func AgreeAddFriend(pid, friendId PLAYER_ID, t int32) {
	fq := for_game.GetFriendBase(friendId)
	fq.AgreeAddFriend(pid)
	fq.AddFriend(pid, t)
	msg1 := GetFriendInfo(pid)
	SendMsgToHallClientNew(friendId, "RpcNoticeAgreeFriend", msg1)
	if PlayerOnlineMgr.CheckPlayerIsOnLine(friendId) { //如果在线
		//friend := PlayerMgr.LoadPlayer(friendId)
		msg1 := GetFriendInfo(pid)
		//if friend != nil { //如果同大厅
		ep1 := ClientEpMp.LoadEndpoint(friendId)
		if ep1 != nil {
			ep1.RpcNoticeAgreeFriend(msg1)
		} else { //如果不同大厅
			SendMsgToHallClientNew(friendId, "RpcNoticeAgreeFriend", msg1)
		}
	}
}

//获取所有好友申请信息
func (self *cls1) RpcGetFriendRequest(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("====RpcGetFriendRequest=== reqMst=%v", reqMsg)
	fq := for_game.GetFriendBase(who.GetPlayerId())
	msg := fq.GetNewVersionAllFriendRequestForOne()
	return msg
}

// 新版本获取所有好友申请信息
func (self *cls1) RpcNewVersionGetFriendRequest(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("====新版本获取新的好友RpcNewVersionGetFriendRequest=== reqMst=%v", reqMsg)
	fq := for_game.GetFriendBase(who.GetPlayerId())
	if fq == nil {
		return easygo.EmptyMsg
	}
	msg := fq.GetNewVersionAllFriendRequestForOne()
	return msg
}

// 新的朋友列表中,删除记录.
func (self *cls1) RpcDelNewFriendList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.DelNewFriendListReq, common ...*base.Common) easygo.IMessage {
	logs.Info("====RpcDelNewFriendList=== reqMst=%v", reqMsg)
	if len(reqMsg.GetPlayerIds()) == 0 {
		logs.Error("RpcDelNewFriendList 需要删除的id列表为空")
		return easygo.EmptyMsg
	}
	fq := for_game.GetFriendBase(who.GetPlayerId())
	fq.DelAddFriend(reqMsg.GetPlayerIds())
	return easygo.EmptyMsg
}

//心跳协议
func (self *cls1) RpcHeartbeat(ep IGameClientEndpoint, ctx interface{}, reqMsg *client_server.NTP, common ...*base.Common) easygo.IMessage {
	//logs.Info("===================心跳协议 RpcHeartbeat", reqMsg, ctx) // 别删，永久留存
	reqMsg.T2 = easygo.NewInt64(time.Now().Unix())
	return reqMsg
}

//修改个人信息:头像，性别，昵称
func (self *cls1) RpcModifyPlayerMsg(ep IGameClientEndpoint, who *Player, reqMsg *client_server.ChangePlayerInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("修改个人信息====RpcModifyPlayerMsg=====", who, reqMsg)
	reason := who.ChangePlayerInfo(reqMsg)
	if reason != "" {
		return easygo.NewFailMsg(reason)
	}
	ep.RpcPlayerAttrChange(who.GetPlayerInfo())
	if reqMsg.GetType() == 1 || reqMsg.GetType() == 2 || reqMsg.GetType() == 3 {
		who.SendMsgToShop("RpcModifyPlayerMsg", reqMsg)
	}
	return nil
}

//查找通讯录
func (self *cls1) RpcSearchAddressBook(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.AddBookInfo, common ...*base.Common) easygo.IMessage {
	plist := reqMsg.GetPhoneList()
	var lst []*client_server.PlayerMsg
	for _, phone := range plist {
		player := for_game.GetRedisAccountByPhone(phone)
		if player == nil {
			continue
		}
		msg := GetFriendInfo(player.GetPlayerId())
		lst = append(lst, msg)
	}
	msg := &client_server.AllPlayerInfo{
		PlayerMsg: lst,
	}
	return msg
}

func CreateTeamOperate(ownerId int64, lst []int64, info *client_hall.CreateTeam) *client_server.TeamMsg {
	//创建会话
	headURL := info.GetHeadUrl()
	teamName := info.GetTeamName()
	adminId := info.GetAdminId()
	memberList := []int64{ownerId}
	memberList = append(memberList, lst...)
	if headURL == "" {
		headURL = for_game.GetRandTeamHeadIcon()
	}
	team := CreateTeam(ownerId, memberList, info)
	if teamName == "" {
		//没设置群名称，默认前三个成员
		for n, id := range memberList {
			p := for_game.GetRedisPlayerBase(id)
			if p != nil {
				teamName += p.GetNickName()
				if n != 2 {
					teamName += " 、"
				}
			}
			if n == 2 {
				break
			}
		}
	}
	teamId := team.GetId()
	teamObj := for_game.GetRedisTeamObj(teamId, team)
	//for_game.AddTeamDataToRedis(team)

	mlst := team.GetMemberList()
	reason := "群主邀请入群"
	m := &share_message.TeamChannel{
		Name: easygo.NewString("群主"),
	}
	if adminId != 0 {
		m.Type = easygo.NewInt32(5) //后台邀请
	} else {
		m.Type = easygo.NewInt32(4)
	}

	for_game.AddTeamPersonData(teamId, teamObj.GetLogMaxId(), mlst, for_game.TEAM_MASSES, reason, m, ownerId)

	for _, pid := range mlst {
		player := for_game.GetRedisPlayerBase(pid)
		if player != nil {
			player.AddTeamId(teamId)
		}
	}

	msg := for_game.GetTeamMsgForHall(teamId)
	msg.IsShow = easygo.NewBool(false)
	serverId := PlayerOnlineMgr.GetPlayerServerId(ownerId)
	TeamSendMessage(mlst, ownerId, serverId, "RpcAddTeamResult", msg)

	session := &share_message.ChatSession{
		Id:             easygo.NewString(teamId),
		Type:           easygo.NewInt32(for_game.CHAT_TYPE_TEAM),
		MaxLogId:       easygo.NewInt64(0),
		SessionName:    easygo.NewString(teamName),
		SessionHeadUrl: easygo.NewString(team.GetHeadUrl()),
		PlayerIds:      team.GetMemberList(),
		TeamName:       easygo.NewString(teamName),
		Topic:          easygo.NewString(team.GetTopic()),
	}
	sessionObj := for_game.GetRedisChatSessionObj(session.GetId(), session)
	sessionObj.SaveToMongo()
	// NoticeTeamMessage(teamId, inviteId, messageType, plst) //

	return msg
}

//创建群聊
func (self *cls1) RpcCreateTeam(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CreateTeam, common ...*base.Common) easygo.IMessage {
	mark := reqMsg.GetMark()
	plst := reqMsg.GetPlayerList()
	if for_game.IS_FORMAL_SERVER {
		if len(plst) < 2 {
			return easygo.NewFailMsg("不可以创建两人群聊")
		}
	}

	var lst []PLAYER_ID
	if mark == 1 { //全选好友
		lst = who.GetFriends()
	} else if mark == 2 { //选中的玩家 加入群
		lst = plst
	} else if mark == 3 { //除了选中的玩家 都加入群
		friendlst := who.GetFriends()
		var newlst []PLAYER_ID
		for _, pid := range friendlst {
			if util.Int64InSlice(pid, plst) {
				continue
			}
			newlst = append(newlst, pid)
		}
		lst = newlst
	} else {
		panic(fmt.Sprintf("mark值错误,%d", mark))
	}
	msg := CreateTeamOperate(who.GetPlayerId(), lst, reqMsg)
	ep.RpcCreateTeamResult(msg)
	return nil
}

//创建话题群组
func (self *cls1) RpcCreateTopicTeam(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CreateTeam, common ...*base.Common) easygo.IMessage {
	//TODO 判断玩家是不是话题主,最多只能创建6个该话题的群
	logs.Info("RpcCreateTopicTeam:", reqMsg)
	if reqMsg.GetTopic() == "" {
		return easygo.NewFailMsg("话题为空")
	}
	if reqMsg.GetTopicDesc() == "" {
		return easygo.NewFailMsg("话题简介为空")
	}
	topic := for_game.GetTopicByNameFromDB(reqMsg.GetTopic())
	if topic == nil {
		return easygo.NewFailMsg("不存在的话题")
	}
	if topic.GetTopicMaster() != who.GetPlayerId() {
		return easygo.NewFailMsg("不是话题管理员，不能创建群组")
	}
	if for_game.GetTopicTeamNumByTopic(reqMsg.GetTopic()) >= 6 {
		return easygo.NewFailMsg("话题群已达上限哦")
	}
	msg := CreateTeamOperate(who.GetPlayerId(), []int64{}, reqMsg)
	ep.RpcCreateTeamResult(msg)
	return nil
}

//主动申请进群
func (self *cls1) RpcActiveAddTeamMember(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamReq, common ...*base.Common) easygo.IMessage {
	logs.Info("===========RpcActive AddTeamMember=================", reqMsg)
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		//panic(fmt.Sprintf("群对象为空：%d", teamId))
		return easygo.NewFailMsg("该群已经解散")
	}
	playerId := who.GetPlayerId()

	if teamObj.GetTeamStatus() == for_game.DISSOLVE {
		return easygo.NewFailMsg("该群已经解散")
	}

	memberlist := teamObj.GetTeamMemberList()
	if util.Int64InSlice(playerId, memberlist) {
		return easygo.NewFailMsg("你已经在该群中")
	}

	if !teamObj.CheckTeamMaxMember(1) {
		//return easygo.NewFailMsg("添加成员数量超过该群聊数量上限")
		return easygo.NewFailMsg("群成员超出上限")
	}

	setting := teamObj.GetTeamMessageSetting()
	if setting.GetIsStopAddTeam() {
		return easygo.NewFailMsg("该群已禁止加群行为")
	}

	inviteId := reqMsg.GetInviteId()
	t := reqMsg.GetType()
	var name string
	if base := for_game.GetRedisPlayerBase(inviteId); base != nil {
		name = base.GetNickName()
	}

	var reason string
	messageType := int32(for_game.REQUEST_ADDTEAM)
	switch t {
	case ADD_TEAM_TYPE_QRCODE: // 1
		reason = fmt.Sprintf("通过\"%s\"发送的二维码进群", name)
	case ADD_TEAM_TYPE_CARD: // 2
		reason = fmt.Sprintf("通过\"%s\"发送的群名片进群", name)
	case ADD_TEAM_TYPE_PASSWORD: // 3
		reason = fmt.Sprintf("通过\"%s\"发送的群链接进群", name)
		messageType = int32(for_game.ACTIVE_ADDTEAM)
	case ADD_TEAM_TYPE_ADV, ADD_TEAM_TYPE_BACKSTAGE: // 5
		reason = fmt.Sprintf("\"%s\"成功加入群聊", name)
		messageType = int32(for_game.ADV_TEAM_MEM)
	case ADD_TEAM_TYPE_TOPIC:
		reason = fmt.Sprintf("通过话题群组页加入群聊")
		messageType = int32(for_game.ACTIVE_ADDTEAM)
	}
	m := &share_message.TeamChannel{
		Name: easygo.NewString(name),
		Type: easygo.NewInt32(t),
	}
	plst := []int64{playerId}
	pos := for_game.GetTeamPlayerPos(teamId, inviteId)
	if setting.GetIsInvite() && pos >= for_game.TEAM_MASSES {
		if !teamObj.GetTeamVaildInviteState(playerId) {
			teamObj.AddTeamInvite(inviteId, name, plst, reason, m)
			NoticeTeamMessage(teamId, inviteId, for_game.REQUEST_ADDTEAM, plst)
		}
		if PlayerOnlineMgr.CheckPlayerIsOnLine(playerId) {
			ep1 := ClientEpMp.LoadEndpoint(playerId) //同一个大厅
			if ep1 != nil {
				ep1.RpcRequestADDTeamInviteSuccess(nil)
			}
		}
	} else {
		serverId := PlayerOnlineMgr.GetPlayerServerId(playerId)
		AddTeamMember(teamId, teamObj.GetSessionLogMaxId(), plst, true, t, reason, serverId, m)
		//NoticeTeamMessage(teamId, inviteId, for_game.ACTIVE_ADDTEAM, plst)
		NoticeTeamMessage(teamId, inviteId, messageType, plst)
		who.AddTeamId(teamId)
	}
	//默认关注话题群
	//如果是话题群，默认关注话题
	if teamObj.GetTopic() != "" {
		topic := for_game.GetTopicByNameFromDB(teamObj.GetTopic())
		if topic != nil {
			for _, pid := range plst {
				for_game.OperateTopic(for_game.OPERATE_TOPIC_ATTENTION, pid, []int64{topic.GetId()})
			}
		}
	}
	return nil
}

//客户端请求验证码
func (self *cls1) RpcClientGetCode(ep IGameClientEndpoint, ctx interface{}, reqMsg *client_server.GetCodeRequest, common ...*base.Common) easygo.IMessage {
	logs.Info("=========请求验证码==========", reqMsg)
	phone := reqMsg.GetPhone()
	t := reqMsg.GetType()
	return self.sendSMS(phone, t, reqMsg.GetAreaCode())
}

/**
sendSMS 内部发送短信公用方法.添加银行卡的时候避免客户端调用两次rpc请求.
phone: 手机号
t:短信类型
*/
func (self *cls1) sendSMS(phone string, t int32, areaCode string) easygo.IMessage {
	if !for_game.IS_FORMAL_SERVER && t != for_game.CLIENT_CODE_BINDBANK {
		//return easygo.NewFailMsg("测试服不需要发送验证码")
		return nil
	}

	if !for_game.MessageMarkInfo.CheckPhoneVaild(phone) {
		return easygo.NewFailMsg("你操作频繁过快，请稍后再试")
	}

	if t == for_game.CLIENT_CODE_PLAYERMESSAGE {
		bHad := for_game.GetRedisAccountByPhone(phone) //检查这个手机号码是否已经注册了账号
		if bHad != nil {
			return easygo.NewFailMsg("该手机号已注册过帐号", for_game.FAIL_MSG_CODE_1008)
		}
	}

	data := for_game.MessageMarkInfo.GetMessageMarkInfo(t, phone)
	if data != nil {
		leaveTime := time.Now().Unix() - data.Timestamp
		if leaveTime <= 120 {
			return easygo.NewFailMsg("验证码已发送!")
		}
	}
	codes := for_game.SendCodeToClientUser(t, phone, areaCode)
	if codes != "" {
		return nil
	}
	return easygo.NewFailMsg("验证码发送失败！")
}

//检查验证码
func (self *cls1) RpcCheckMessageCode(ep IGameClientEndpoint, ctx interface{}, reqMsg *client_server.CodeResponse, common ...*base.Common) easygo.IMessage {
	code := reqMsg.GetCode()
	t := reqMsg.GetType()
	phone := reqMsg.GetPhone()
	return for_game.CheckMessageCode(phone, code, t)
}

//关闭安全密码
func (self *cls1) RpcClosePassword(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	who.DelSafePassword()
	ep.RpcPlayerAttrChange(who.GetPlayerInfo())
	return nil
}

//设置支付 安全密码
func (self *cls1) RpcSetPassword(ep IGameClientEndpoint, who *Player, reqMsg *client_server.PasswordInfo, common ...*base.Common) easygo.IMessage {
	password := reqMsg.GetPassword()
	t := reqMsg.GetType()
	if t == 1 {
		if who.GetPayPassword() != "" {
			res := "支付密码不为空"
			return easygo.NewFailMsg(res)
		}
		who.SetPayPassword(password)
		// who.SetIsPayPassword(true)
		ep.RpcPlayerAttrChange(who.GetPlayerInfo())
	} else if t == 2 {
		if who.GetSafePassword() != "" {
			res := "安全密码不为空"
			return easygo.NewFailMsg(res)
		}
		who.SetSafePassword(password)
		ep.RpcPlayerAttrChange(who.GetPlayerInfo())
	} else if t == 3 {
		old := reqMsg.GetOldPassword()
		accountInfo := for_game.GetRedisAccountObj(who.GetPlayerId())
		if accountInfo == nil {
			res := "账号数据异常"
			return easygo.NewFailMsg(res)
		}
		if accountInfo.GetPassword() != for_game.Md5(old) {
			res := "登录密码错误"
			return easygo.NewFailMsg(res)
		}
		accountInfo.SetPassword(password)
	}
	return nil
}

//修改支付 安全 登录密码
func (self *cls1) RpcChangePassword(ep IGameClientEndpoint, who *Player, reqMsg *client_server.PasswordInfo, common ...*base.Common) easygo.IMessage {
	password := reqMsg.GetPassword()
	t := reqMsg.GetType()
	accountInfo := for_game.GetRedisAccountObj(who.GetPlayerId())
	if accountInfo == nil {
		res := "账号数据异常"
		return easygo.NewFailMsg(res)
	}
	if t == 1 {
		if who.GetPayPassword() == "" {
			res := "支付密码为空，请先设置支付密码"
			return easygo.NewFailMsg(res)
		}
		who.SetPayPassword(password)
	} else if t == 2 {
		if who.GetSafePassword() == "" {
			res := "安全密码为空，请先设置支付密码"
			return easygo.NewFailMsg(res)
		}
		who.SetSafePassword(password)
	} else if t == 3 {
		if accountInfo.GetPassword() == "" {
			res := "登录密码为空"
			return easygo.NewFailMsg(res)
		}
		accountInfo.SetPassword(password)
	}
	return nil
}

//忘记支付密码重新设置
func (self *cls1) RpcForgetPayPassword(ep IGameClientEndpoint, who *Player, reqMsg *client_server.PasswordInfo, common ...*base.Common) easygo.IMessage {
	pd := reqMsg.GetPassword()
	who.SetPayPassword(pd)
	who.ClearBankId()
	return nil
}

//身份验证
func (self *cls1) RpcCheckPlayerPeopleId(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PeopleIdInfo, common ...*base.Common) easygo.IMessage {
	name := reqMsg.GetName()
	peopleId := reqMsg.GetPeopleId()
	if who.GetPeopleId() == peopleId && who.GetRealName() == name {
		return nil
	}
	res := "身份验证错误，请重新输入"
	return easygo.NewFailMsg(res)
}

//身份实名认证
func (self *cls1) RpcSetPeopleAuth(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PeopleIdInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("=== 身份实名认证RpcSetPeopleAuth,reqMsg:", reqMsg)
	id := reqMsg.GetPeopleId()
	if len(id) != 18 {
		return easygo.NewFailMsg("身份证号码不必须是18位")
	}

	period := for_game.GetPlayerPeriod(who.GetPlayerId())
	var num int64
	m := period.DayPeriod.Fetch("check_peopleId")
	if m == nil {
		num = 0
	} else {
		num = m.(int64)
	}

	if num >= 3 {
		return easygo.NewFailMsg("今日身份验证次数已达上限")
	}
	period.DayPeriod.AddInteger("check_peopleId", 1)
	name := reqMsg.GetName()

	if !for_game.CheckPeopleIdIsValid(id) {
		return easygo.NewFailMsg("身份证号码已被使用")
	}

	if for_game.IS_FORMAL_SERVER {
		if !for_game.AuthPeopleIdName(id, name) {
			return easygo.NewFailMsg("姓名与身份证不符")
		}
	}
	who.SetPeopleAuth(id, name)
	ep.RpcPlayerAttrChange(who.GetPlayerInfo())

	//msg := &shop_hall.AuthInfo{
	//	Name:     &name,
	//	PeopleId: &id,
	//	PlayerId: easygo.NewInt64(who.GetPlayerId()),
	//}
	//who.SendMsgToShop("RpcPeopleAuthCallBack", msg)
	for_game.SetRedisRegisterLoginReportFildVal(easygo.Get0ClockTimestamp(util.GetMilliTime()), 1, "RealNameCount") //更新埋点报表实名认证人数
	return nil
}

//检查身份证id的有效性
func (self *cls1) RpcCheckPeopleIdValid(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PeopleIdInfo, common ...*base.Common) easygo.IMessage {
	msg := &client_hall.CheckPeopleInfo{}
	id := reqMsg.GetPeopleId()
	if len(id) != 18 {
		msg.Type = easygo.NewInt32(1)
	} else if !for_game.CheckPeopleIdIsValid(id) { //如果身份证id 已经被人用过
		msg.Type = easygo.NewInt32(2)
	} else {
		msg.Type = easygo.NewInt32(3)
	}
	return msg
}

//通过银行卡号获取银行名字
func (self *cls1) RpcGetBankCode(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.BankMessage, common ...*base.Common) easygo.IMessage {
	bankId := reqMsg.GetBankCardNo()
	//if !for_game.IS_FORMAL_SERVER {
	//	reqMsg.BankCode = easygo.NewString("HKB")
	//	return reqMsg
	//}
	code := for_game.GetBankCodeForBankId(bankId)
	reqMsg.BankCode = easygo.NewString(code)
	return reqMsg
}

//获取绑卡短信验证码
// RpcBindBankCode 每分钟1次，一天6次
func (self *cls1) RpcBindBankCode(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.BankMessage, common ...*base.Common) easygo.IMessage {
	if !who.CheckPeopleAuth() {
		res := "请先进行实名认证"
		return easygo.NewFailMsg(res)
	}
	logs.Info("RpcBindBankCode:", reqMsg)

	bankId := reqMsg.GetBankCardNo()
	if reqMsg.GetIsModify() {
		//如果是修改卡区域信息
		smsRespErr := self.sendSMS(who.GetPhone(), for_game.SmsTypeBindBankCard, who.GetArea())
		ep.RpcBindBankCodeResult(reqMsg)
		return smsRespErr
	}
	if util.InStringSlice(bankId, who.GetBankId()) {
		res := "已绑定该银行卡"
		return easygo.NewFailMsg(res)
	}
	//暂时只支持汇潮的银行卡绑定，后续绑卡去掉限制
	//if !easygo.Contain(PWebHuiChaoPay.BankList, reqMsg.GetBankCode()) {
	//	res := "暂不支持此类银行卡"
	//	return easygo.NewFailMsg(res)
	//}
	if time.Now().Unix() < who.LastOSMTime+60 {
		return easygo.NewFailMsg("验证码请求过于频繁")
	}
	period := for_game.GetPlayerPeriod(who.GetPlayerId())
	//每日短信请求次数
	osmSum := period.DayPeriod.FetchInt64("check_bank")
	if osmSum >= HJ_OSM_DATE_OSM_TIME {
		return easygo.NewFailMsg("今日银行卡次数已达上限")
	}
	// 验证绑卡人银行卡信息是否有误
	if b, msg := for_game.AuthBankIdName(reqMsg.GetBankCardNo(), who.GetRealName(), reqMsg.GetIdNo(), reqMsg.GetMobileNo()); !b {
		logs.Error("认证银行卡错误:", msg)
		return easygo.NewFailMsg("认证信息不匹配")
	}
	// 调用发送内部的短信验证码接口 t=5
	//默认只能绑定国内卡，所以不能用国际短信
	smsRespErr := self.sendSMS(reqMsg.GetMobileNo(), for_game.SmsTypeBindBankCard, who.GetArea())
	ep.RpcBindBankCodeResult(reqMsg)
	return smsRespErr
}

// RpcAddBank 增加银行卡
func (self *cls1) RpcAddBank(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.BankMessage, common ...*base.Common) easygo.IMessage {
	code := reqMsg.GetMsgCode()
	if code == "" || len(code) != 6 {
		return easygo.NewFailMsg("请输入有效的验证码")
	}
	// 验证短信验证码
	if err := for_game.CheckMessageCode(reqMsg.GetMobileNo(), reqMsg.GetMsgCode(), for_game.CLIENT_CODE_BINDBANK); err != nil {
		logs.Error("RpcAddBank 增加银行卡,验证平台发送的短信验证码验证失败,手机号为: %s,验证码为: %s,短信类型为: %d",
			who.GetPhone(), reqMsg.GetMsgCode(), for_game.CLIENT_CODE_BINDBANK)
		return err
	}
	if reqMsg.GetIsModify() {
		//单纯的修改卡区域
		who.ModifyBankInfo(reqMsg)
		ep.RpcPlayerAttrChange(who.GetPlayerInfo())
		return nil
	}
	bankName := for_game.BankName[reqMsg.GetBankCode()]
	msg := &share_message.BankInfo{
		BankId:    easygo.NewString(reqMsg.GetBankCardNo()),
		BankCode:  easygo.NewString(reqMsg.GetBankCode()),
		Time:      easygo.NewInt64(time.Now().Unix()),
		SignNo:    easygo.NewString(reqMsg.GetSignNo()),
		BankName:  easygo.NewString(bankName),
		BankPhone: easygo.NewString(reqMsg.GetMobileNo()),
		City:      easygo.NewString(reqMsg.GetCity()),
		Provice:   easygo.NewString(reqMsg.GetProvice()),
		Area:      easygo.NewString(reqMsg.GetArea()),
	}
	who.AddBankInfo(msg)
	ep.RpcPlayerAttrChange(who.GetPlayerInfo())
	easygo.Spawn(func() {
		for_game.MakePlayerBehaviorReport(3, who.GetPlayerId(), nil, nil, nil, nil) //生成用户行为报表绑卡字段 已优化到Redis
		for_game.SetRegisterLoginReportBankCardCount(who.GetPlayerId())             //更新埋点报表绑卡人数
	})
	//绑定今日次数累加
	period := for_game.GetPlayerPeriod(who.GetPlayerId())
	if period != nil {
		period.DayPeriod.AddInteger("check_bank", 1)
	}
	return nil
}

//删除银行卡
func (self *cls1) RpcDelBank(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.BankMessage, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcDelBank:", who.GetPlayerId(), reqMsg)
	bankId := reqMsg.GetBankCardNo()
	if !util.InStringSlice(bankId, who.GetBankId()) {
		res := "没绑定该银行卡"
		return easygo.NewFailMsg(res)
	}
	who.DelBankInfo(bankId)
	ep.RpcPlayerAttrChange(who.GetPlayerInfo())
	logs.Info("删除银行卡完成")
	return nil
}

//检查密码
func (self *cls1) RpcCheckPassword(ep IGameClientEndpoint, who *Player, reqMsg *client_server.PasswordInfo, common ...*base.Common) easygo.IMessage {
	pd := reqMsg.GetPassword()
	t := reqMsg.GetType()
	if t == 1 {
		if !who.GetIsPayPassword() {
			res := "未设置支付密码"
			return easygo.NewFailMsg(res)
		}
		period := for_game.GetPlayerPeriod(who.GetPlayerId())
		var num int64
		m := period.DayPeriod.Fetch(for_game.CHECK_PAYPASSWORD)
		if m == nil {
			num = 0
		} else {
			num = m.(int64)
		}
		if num >= 3 {
			res := fmt.Sprintf("今日密码输入错误次数已达上限，请明天再试")
			return easygo.NewFailMsg(res)
		}
		if who.GetPayPassword() != for_game.Md5(pd) {
			period.DayPeriod.AddInteger(for_game.CHECK_PAYPASSWORD, 1)
			count := num + 1
			if count == 3 {
				res := fmt.Sprintf("今日密码输入错误次数已达上限，请明天再试")
				return easygo.NewFailMsg(res)
			}
			res := fmt.Sprintf("支付密码错误,你还可以输入%d次", 3-count)
			return easygo.NewFailMsg(res)
		} else {
			period.DayPeriod.Set(for_game.CHECK_PAYPASSWORD, 0)
		}
	} else if t == 3 {
		account := for_game.GetRedisAccountObj(who.GetPlayerId())
		if account == nil {
			res := "账号为空"
			return easygo.NewFailMsg(res)
		}

		if for_game.Md5(pd) != account.GetPassword() {
			res := "登录密码错误"
			return easygo.NewFailMsg(res)
		}
	}
	return nil
}

//冻结账号
func (self *cls1) RpcFrezzeAccount(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.AccountInfo, common ...*base.Common) easygo.IMessage {
	account := reqMsg.GetAccount()
	if account != who.GetPhone() {
		res := "该账号不是该用户的账号"
		return easygo.NewFailMsg(res)
	}
	if !for_game.LoginAccountAuth(account) {
		res := "该账号已经被冻结"
		return easygo.NewFailMsg(res)
	}
	who.SetStatus(1)
	for_game.AddFreezeAccount(account)
	return nil
}

//解冻账号
func (self *cls1) RpcUnFrezzeAccount(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.AccountInfo, common ...*base.Common) easygo.IMessage {
	account := reqMsg.GetAccount()
	if account != who.GetPhone() {
		res := "该账号不是该用户的账号"
		return easygo.NewFailMsg(res)
	}
	if for_game.LoginAccountAuth(account) {
		res := "该账号未被冻结"
		return easygo.NewFailMsg(res)
	}
	who.SetStatus(0)
	for_game.DelFreezeAccount(account)
	return nil
}

//获取版本号
func (self *cls1) RpcGetVersion(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	version := easygo.YamlCfg.GetValueAsString("VERSION_NUMBER")
	msg := &client_hall.VersionInfo{
		Version: easygo.NewString(version),
	}
	return msg
}

//获取黑名单
func (self *cls1) RpcGetBlackList(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	msg := who.GetBlackInfo()
	return msg
}

//删除好友
func (self *cls1) RpcDelFriend(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.DelFriendInfo, common ...*base.Common) easygo.IMessage {
	pid := reqMsg.GetPlayerId()
	isBlack := reqMsg.GetIsBlack()

	if !util.Int64InSlice(pid, who.GetFriends()) {
		res := "玩家不是你的好友"
		return easygo.NewFailMsg(res)
	}

	if isBlack {
		blacklst := who.GetBlackList()
		if !util.Int64InSlice(pid, blacklst) {
			blacklst = append(blacklst, pid)
			who.SetRedisBlackList(blacklst)
		}
	}
	wq := for_game.GetFriendBase(who.GetPlayerId())
	wq.DelFriend(pid)

	//亲密度删除
	intimacyObj := for_game.GetRedisPlayerIntimacyObj(for_game.MakeSessionKey(who.GetPlayerId(), pid))
	if intimacyObj != nil {
		intimacyObj.CleanPlayerIntimacy()
	}
	//sayHi信息删除
	for_game.DelVoiceCardSayHiLog(who.GetPlayerId(), pid)
	//重置陌生人聊天次数
	id := for_game.MakeSessionKey(who.GetPlayerId(), pid)
	session := for_game.GetRedisChatSessionObj(id)
	if session != nil {
		session.ResetStrangerTalkNum(who.GetPlayerId())
	}
	//聊天记录匹配信息变删除状态
	//chatObj := for_game.GetRedisPersonalChatLogObj(id)
	//if chatObj != nil {
	//	log := chatObj.DelSayHiLog()
	//	//通知双方客户
	//	delMsg := &client_hall.SayHiLog{
	//		SessionId: easygo.NewString(id),
	//		Log:       log,
	//	}
	//	SendMsgToClient([]int64{who.GetPlayerId(), pid}, "RpcDelSayHiLog", delMsg)
	//}
	return reqMsg
}

//黑名单操作
func (self *cls1) RpcBlackOperate(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.BlackInfo, common ...*base.Common) easygo.IMessage {
	logs.Debug("RpcBlackOperate", reqMsg)
	t := reqMsg.GetType()
	pid := reqMsg.GetPlayerId()
	blacklst := who.GetBlackList()
	if for_game.GetPlayerById(pid) == nil { //如果没有这个账号
		res := fmt.Sprintf("无效玩家id:%d", pid)
		return easygo.NewFailMsg(res)
	}
	if t == 1 {
		if util.Int64InSlice(pid, blacklst) {
			res := "玩家已经在你的黑名单中"
			return easygo.NewFailMsg(res)
		}
		blacklst = append(blacklst, pid)
		who.SetRedisBlackList(blacklst)
	} else if t == 2 {
		if !util.Int64InSlice(pid, blacklst) {
			res := "玩家不在你的黑名单中"
			return easygo.NewFailMsg(res)
		}
		blacklst = util.SliceDelOneInt64(blacklst, pid)
		who.SetRedisBlackList(blacklst)
	} else {
		res := "未知类型"
		return easygo.NewFailMsg(res)
	}
	//msg := &shop_hall.BlackList{List: blacklst, PlayerId: easygo.NewInt64(who.GetPlayerId())}
	//who.SendMsgToShop("RpcBlackListCallBack", msg)
	return reqMsg
}

//畅聊助手
func (self *cls1) RpcGetCLAssistant(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.AInfo, common ...*base.Common) easygo.IMessage {
	list := BuildUnreadAssistantList(who, 0)
	if reqMsg.GetType() == 1 {
		who.SetLastAssistantTime(util.GetMilliTime())
	}

	return &client_server.AssistantInfo{AssistantInfoList: list}
}

//零钱助手
func (self *cls1) RpcGetMoneyAssistant(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PageInfo, common ...*base.Common) easygo.IMessage {
	page := int(reqMsg.GetPage())
	num := int(reqMsg.GetNum())
	lst := for_game.GetGoldLogInfo(who.GetPlayerId(), page, num)
	var allmsg []*client_hall.MoneyAssistantInfo

	for _, info := range lst {
		extend := info.Extend
		msg := &client_hall.MoneyAssistantInfo{
			LogId:      easygo.NewInt64(info.GetLogId()),
			PlayerId:   easygo.NewInt64(info.GetPlayerId()),
			ChangeGold: easygo.NewInt64(info.GetChangeGold()),
			SourceType: easygo.NewInt32(info.GetSourceType()),
			Note:       easygo.NewString(info.GetNote()),
			CreateTime: easygo.NewInt64(info.GetCreateTime()),
			Gold:       easygo.NewInt64(extend.GetGold()),
			OrderId:    easygo.NewString(extend.GetOrderId()),
		}
		allmsg = append(allmsg, msg)
	}

	msg := &client_hall.AllMoneyAssistantInfo{
		Info: allmsg,
	}
	return msg
}

//附近的人
func (self *cls1) RpcGetLocationInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.LocationInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("==========大厅附近的人 RpcGetLocationInfo========", reqMsg)
	x := reqMsg.GetX()
	y := reqMsg.GetY()
	pid := who.GetPlayerId()
	base := GetPlayerObj(pid)
	base.SetStation(x, y)
	//base.SetProvice(reqMsg.GetProvince())
	//base.SetCity(reqMsg.GetCity())
	msg := for_game.GetLocationInfo(pid, reqMsg)
	return msg
}

//所有在附近的人给我打招呼的信息
func (self *cls1) RpcGetAllNearByInfo(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	lst := for_game.GetAllNearByInfo(who.GetPlayerId())
	msg := &client_hall.AllNearByMessage{
		NearByMessage: lst,
	}
	return msg
}

//附近的人打招呼信息
func (self *cls1) RpcGetNearByInfo(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	lst := for_game.GetAllNewNearByInfo(who.GetPlayerId())
	msg := &client_hall.AllNearByMessage{
		NearByMessage: lst,
	}
	return msg
}

// 整合附近的人的信息
func (self *cls1) RpcGetNearByInfo2(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.LocationInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("==========大厅附近的人 RpcGetNearByInfo2========", reqMsg)

	x := reqMsg.GetX()
	y := reqMsg.GetY()
	pid := who.GetPlayerId()
	base := GetPlayerObj(pid)
	base.SetStation(x, y)
	//base.SetProvice(reqMsg.GetProvince())
	//base.SetCity(reqMsg.GetCity())
	locationInfo := for_game.GetLocationInfo(pid, reqMsg)

	//附近的人打招呼信息
	lst2 := for_game.GetAllNewNearByInfo(who.GetPlayerId())
	msg := &client_hall.NearByInfoReply{
		AllNearByInfo:    locationInfo,
		AllNewNearByInfo: lst2,
	}

	return msg
}

/*func (self *cls1) RpcGetNearInfoByPage(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.LocationInfoByPage, common... *base.Common) easygo.IMessage {
	logs.Info("==========大厅分页获取附近的人 RpcGetNearInfoByPage========", reqMsg)
	x := reqMsg.GetX()
	y := reqMsg.GetY()
	pid := who.GetPlayerId()
	base := GetPlayerObj(pid)
	base.SetStation(x, y)
	//base.SetProvice(reqMsg.GetProvince())
	//base.SetCity(reqMsg.GetCity())
	locationInfo := for_game.GetLocationInfo(pid, reqMsg)

	//附近的人打招呼信息
	lst2 := for_game.GetAllNewNearByInfo(who.GetPlayerId())
	msg := &client_hall.NearByInfoReply{
		AllNearByInfo:    locationInfo,
		AllNewNearByInfo: lst2,
	}
	return msg
}*/

//给附近的人打招呼
func (self *cls1) RpcAddNearByMessage(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.NearByMessage, common ...*base.Common) easygo.IMessage {
	pid := reqMsg.GetPlayerId()
	if pid >= for_game.Min_Robot_PlayerId && pid <= for_game.Max_Robot_PlayerId {
		return nil
	}

	base := GetPlayerObj(pid)
	if util.Int64InSlice(who.GetPlayerId(), base.GetBlackList()) {
		res := "玩家拒绝接受你的消息"
		return easygo.NewFailMsg(res)
	}

	content := reqMsg.GetContent()
	for_game.AddNearByInfo(pid, who.GetPlayerId(), content)
	if PlayerOnlineMgr.CheckPlayerIsOnLine(pid) {
		//player := PlayerMgr.LoadPlayer(pid)
		//if player != nil {
		ep := ClientEpMp.LoadEndpoint(pid)
		if ep != nil {
			ep.RpcNoticeNearByInfo(nil)
		} else {
			SendMsgToHallClientNew(pid, "RpcNoticeNearByInfo", nil)
		}
	} else {
		base := GetPlayerObj(pid)
		base.SetIsNearBy(true)
	}
	return nil
}

//清除自己的附近信息
func (self *cls1) RpcDeleteNearByInfo(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	base := GetPlayerObj(who.GetPlayerId())
	base.SetStation(0, 0)
	return nil
}

//清除给自己的打招呼信息
func (self *cls1) RpcDeleteMyNearByInfo(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	for_game.DelNearByInfo(who.GetPlayerId())
	return nil
}

//同步时间
func (self *cls1) RpcSyncTime(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	msg := &client_server.NTP{
		T1: easygo.NewInt64(time.Now().Unix()),
	}
	return msg
}

//获取账单消息
func (self *cls1) RpcGetMoneyOrderInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.MoneyType, common ...*base.Common) easygo.IMessage {
	logs.Info("----------获取账单消息----------", reqMsg)
	t := reqMsg.GetType()
	typeList := make([]int32, 0)
	switch t {
	case 0: //全部
	case 1: //充值提现
		typeList = append(typeList, for_game.GOLD_TYPE_CASH_IN, for_game.GOLD_TYPE_CASH_OUT)
	case 2: //零钱红包
		typeList = append(typeList, for_game.GOLD_TYPE_GET_REDPACKET, for_game.GOLD_TYPE_SEND_REDPACKET)
	case 3: //商户消费
		typeList = append(typeList, for_game.GOLD_TYPE_SHOP_MONEY, for_game.GOLD_TYPE_SHOP_ITEM_MONEY)
	case 4: //退款
		typeList = append(typeList, for_game.GOLD_TYPE_REDPACKET_OVERTIME, for_game.GOLD_TYPE_TRANSFER_MONEY_OVER, for_game.GOLD_TYPE_BACK_MONEY, for_game.GOLD_TYPE_CASH_OUT_BACK)
	case 5: //二维码收付款
		typeList = append(typeList, for_game.GOLD_TYPE_GET_MONEY, for_game.GOLD_TYPE_PAY_MONEY)
	case 6: //转账
		typeList = append(typeList, for_game.GOLD_TYPE_GET_TRANSFER_MONEY, for_game.GOLD_TYPE_SEND_TRANSFER_MONEY)
	case 7: //其他
		typeList = append(typeList, for_game.GOLD_TYPE_FINE_MONEY, for_game.GOLD_TYPE_EXTRA_MONEY)
	default:
		panic(fmt.Sprintf("错误的类型t:%d", t))
	}

	page := int(reqMsg.GetPage())
	num := int(reqMsg.GetNum())
	year := int(reqMsg.GetYear())
	month := int(reqMsg.GetMonth())
	pageList := for_game.GetPageGoldLogs(who.GetPlayerId(), typeList, page, num, year, month)
	if len(pageList) < num {
		count := num - len(pageList)
		newList := for_game.GetMoneyOrderInfo(who.GetPlayerId(), typeList, count, year, month)
		pageList = append(pageList, newList...)
	}
	var lst []*client_hall.OrderInfo
	for _, info := range pageList {
		extend := info.Extend
		orderId := extend.GetOrderId()
		orderObj := for_game.GetRedisOrderObj(orderId)
		var order *share_message.Order
		if orderObj != nil {
			order = orderObj.GetRedisOrder()
		}

		if order != nil && order.GetPayWay() > for_game.PAY_TYPE_MONEY {
			continue
		}

		t := info.GetSourceType()
		m := &client_hall.OrderInfo{
			Name:     easygo.NewString(extend.GetTitle()),
			Time:     easygo.NewInt64(info.GetCreateTime()),
			Money:    easygo.NewInt64(info.GetChangeGold()),
			HeadIcon: easygo.NewString(extend.GetHeadIcon()),
			Type:     easygo.NewInt32(t),
			OrderId:  easygo.NewString(orderId),
			Text:     easygo.NewString(extend.GetTransferText()),
		}
		redId := extend.GetRedPacketId()
		if redId != 0 && (t == for_game.GOLD_TYPE_SEND_REDPACKET || t == for_game.GOLD_TYPE_GET_REDPACKET) {
			redPacker := for_game.GetRedisRedPacket(redId)
			if redPacker != nil {
				m.RedPacket = redPacker.GetRedisRedPacket()
			} else {
				logs.Error("无效的红包id:", redId)
			}

		}
		payType := extend.GetPayType()
		if payType == for_game.PAY_TYPE_GOLD {
			m.PayName = easygo.NewString("零钱")
		} else if payType == for_game.PAY_TYPE_WEIXIN {
			m.PayName = easygo.NewString("微信支付")
		} else if payType == for_game.PAY_TYPE_ZHIFUBAO {
			m.PayName = easygo.NewString("支付宝支付")
		} else if payType == for_game.PAY_TYPE_BANKCARD {
			m.PayName = easygo.NewString(extend.GetBankName())
		}

		if t == for_game.GOLD_TYPE_SEND_TRANSFER_MONEY || t == for_game.GOLD_TYPE_GET_TRANSFER_MONEY ||
			t == for_game.GOLD_TYPE_PAY_MONEY || t == for_game.GOLD_TYPE_GET_MONEY {
			m.TransferText = easygo.NewString(extend.GetTransferText())
		}
		if order == nil {
			m.Statue = easygo.NewInt32(for_game.ORDER_ST_FINISH)
		} else {
			m.Statue = easygo.NewInt32(order.GetStatus())
		}
		if t == for_game.GOLD_TYPE_CASH_OUT {
			if order != nil {
				m.ReceiveTime = easygo.NewInt64(order.GetOverTime())
			}
			m.BankName = easygo.NewString(extend.GetBankName())
		}
		lst = append(lst, m)
	}
	msg := &client_hall.AllOrderInfo{
		OrderInfo: lst,
	}

	return msg
}

//获取零钱明细
func (self *cls1) RpcGetCashInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PageInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("======获取零钱明细 RpcGetCashInfo======,reqMsg=%+v", reqMsg)
	page := reqMsg.GetPage()
	num := reqMsg.GetNum()

	orderList := for_game.GetCashOrderInfo(who.GetPlayerId(), int(page), int(num))
	newlst := for_game.GetPlayerGoldChangeLogForGoldExtendLog(who.GetPlayerId(), []int32{})
	orderList = append(orderList, newlst...)
	var lst []*client_hall.CashInfo
	for _, info := range orderList {
		extend := info.Extend
		if extend.GetChannel() > for_game.PAY_TYPE_MONEY {
			continue
		}
		m := &client_hall.CashInfo{
			OrderId:    easygo.NewString(extend.GetOrderId()),
			Time:       easygo.NewInt64(info.GetCreateTime()),
			ChangeGold: easygo.NewInt64(info.GetChangeGold()),
			//Gold:       easygo.NewInt64(extend.GetGold()),
			Gold: easygo.NewInt64(info.GetGold()),
			Type: easygo.NewInt32(info.GetSourceType()),
		}
		lst = append(lst, m)
	}
	msg := &client_hall.AllCashInfo{
		CashInfo: lst,
	}

	return msg
}

//红包明细
func (self *cls1) RpcGetRedPacketInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.RedPacketInfo, common ...*base.Common) easygo.IMessage {
	var gt int32
	t := reqMsg.GetType()
	if t == 1 {
		gt = for_game.GOLD_TYPE_GET_REDPACKET
	} else if t == 2 {
		gt = for_game.GOLD_TYPE_SEND_REDPACKET
	}
	year := int(reqMsg.GetYear())
	month := int(reqMsg.GetMonth())
	page := int(reqMsg.GetPage())
	num := int(reqMsg.GetNum())

	redisTotal := for_game.GetRedisRedPacketTotalEx(who.GetPlayerId(), year, month)

	pageList := for_game.GetPageGoldLogs(who.GetPlayerId(), []int32{gt}, page, num, year, month) //取redis当前页数数量的红包
	if len(pageList) < num {                                                                     //如果当前页数的数量不够 就去查库
		count := num - len(pageList)
		lst := for_game.GetRedPacketForPageInfo(who.GetPlayerId(), gt, count, year, month)
		pageList = append(pageList, lst...)
	}

	var lst []*share_message.RedPacket
	redPacketInfo := for_game.GetAllRedPacketInfo(who.GetPlayerId(), t, year, month) //拿数据库的当前月的所有红包
	newlst := for_game.GetAllRedPackForRedis(who.GetPlayerId(), t, year, month)      //redis缓存的当月红包
	for id, log := range newlst {
		redPacketInfo[id] = log
	}

	for _, info := range pageList {
		extend := info.Extend
		redId := extend.GetRedPacketId()
		redPacket := redPacketInfo[redId]
		if redPacket == nil {
			continue
		}
		lst = append(lst, redPacket)
	}
	msg := &client_hall.AllRedPacketInfo{
		AllInfo: lst,
		Type:    easygo.NewInt32(t),
	}
	if t == 1 { //收红包
		msg.LuckCnt = easygo.NewInt32(redisTotal.GetLuckCnt())
		msg.TotalCnt = easygo.NewInt32(redisTotal.GetRecCount())
		msg.TotalGold = easygo.NewInt64(redisTotal.GetRecTotalMoney())
	} else {
		msg.TotalCnt = easygo.NewInt32(redisTotal.GetSendCount())
		msg.TotalGold = easygo.NewInt64(redisTotal.GetSendTotalMoney())
	}
	return msg

}

//清除本地聊天记录
func (self *cls1) RpcClearLocalTime(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	who.SetClearLocalLogTime()
	return nil
}

//意见  投诉
func (self *cls1) RpcAddComplaintInfo(ep IGameClientEndpoint, who *Player, reqMsg *share_message.PlayerComplaint, common ...*base.Common) easygo.IMessage {
	logs.Info("========意见投诉 RpcAddComplaintInfo=======,reqMsg=%v", reqMsg)
	var title string
	var content string
	logs.Debug("RpcAddComplaintInfo", reqMsg)
	t := reqMsg.GetType()
	now := util.GetMilliTime()
	reqMsg.CreateTime = easygo.NewInt64(now)
	diff := int64(int(30) * 60000)
	var flag string
	var obj int64
	if t == 1 { //意见反馈
		who.SetComplaintTime()
		flag = "您已提交反馈"
		title = "意见反馈"
		content = "尊敬的用户您好，您的投诉建议我们已收到。感谢大家的支持，相信我们会做的更好。"
	} else if t == 2 { //投诉
		title = "个人投诉反馈"
		flag = "您已投诉过此账号"
		content = fmt.Sprintf("您对%v的投诉内容我们已收到，客服会尽快查证核实，会第一时间给予您投诉反馈。", reqMsg.GetName())
		obj = reqMsg.GetRespondentId()

	} else if t == 3 { //商城订单投诉
		flag = "您已投诉过此订单"
		title = "订单投诉反馈"
		content = fmt.Sprintf("您对订单%v的投诉我们已经收到，客服会在第一时间反馈查证跟进，请耐心等待。", reqMsg.GetOrderId())
		obj = reqMsg.GetOrderId()
	} else if t == 4 { //商城物品投诉
		flag = "您已投诉过此商品"
		base := GetPlayerObj(reqMsg.GetRespondentId())
		if base == nil {
			logs.Error("玩家ID:%v不存在", reqMsg.GetRespondentId())
			return nil
		}
		title = "商品投诉反馈"
		content = fmt.Sprintf("您对卖家%v（ID:%v）的商品(%v)投诉我们已经收到。客服会尽快查实，会第一时间给予您投诉反馈。", base.GetNickName(), base.GetAccount(), reqMsg.GetName())
		logs.Info("卖家:", base.GetAccount())
		obj = reqMsg.GetGoodsId()
	} else if t == 5 { //群投诉
		flag = "您已投诉过此群"
		title = "群投诉反馈"
		content = fmt.Sprintf("您对%v群的投诉我们已经收到。客服会尽快查证核实，会第一时间给予您投诉反馈。", reqMsg.GetName())
		obj = reqMsg.GetRespondentId()
	} else if t == 6 { //硬币反馈
		flag = "您已提交反馈"
		title = "意见反馈"
		content = "尊敬的用户您好，您的投诉建议我们已收到。感谢大家的支持，相信我们会做的更好。"
	} else if t == 7 { // 动态投诉
		title = "动态投诉反馈"
		flag = "动态投诉"
		//通过动态id查询发布者昵称
		dy := for_game.GetRedisDynamic(reqMsg.GetDynamicId())
		p := for_game.GetRedisPlayerBase(dy.GetPlayerId())
		content = fmt.Sprintf(`您对"%s"所发动态的投诉内容我们已收到，客服会尽快查证核实，会第一时间给予您投诉反馈。`, p.GetNickName())
		obj = reqMsg.GetDynamicId()
	} else if t == 8 { // 话题举报
		title = "话题举报反馈"
		flag = "话题举报"
		//通过动态id查询发布者昵称
		content = fmt.Sprintf(`您的举报已经受理，我们回尽快处理，谢谢您让我们变得更好。`)
		obj = reqMsg.GetTopicId()
	} else {
		panic(fmt.Sprintf("未知类型 %d", t))
	}

	playerComp := GetPlayerComplaint(who.GetPlayerId(), obj, t)
	if playerComp != nil {
		if playerComp.GetCreateTime() != 0 && now-playerComp.GetCreateTime() < diff {
			logs.Info("playerComp:", playerComp)
			return easygo.NewFailMsg(fmt.Sprintf("%v半小时内请勿重复投诉", flag))
		}
	}

	if title != "" && content != "" {
		NoticeAssistant(who.GetPlayerId(), 2, title, content)
	}

	reqMsg.Status = easygo.NewInt32(1)
	if reqMsg != nil {
		if err := for_game.AddPlayerComplaint(reqMsg); err != nil {
			return easygo.NewFailMsg("投诉失败")
		}
		if t == 7 {
			for_game.UpdateComplaintNum(reqMsg.GetDynamicId(), 1) // 动态投诉次数+1
		}
	}

	return nil
}

//提现请求
func (self *cls1) RpcWithdrawRequest(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.WithdrawInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcWithdrawRequest", reqMsg)
	wAmount := reqMsg.GetAmount()
	//if !for_game.IS_FORMAL_SERVER {
	//	return easygo.NewFailMsg("提现失败,测试服不能提现")
	//}
	//青少年保护模式
	if who.GetYoungPassWord() != "" {
		return easygo.NewFailMsg("青少年模式下无法提现")
	}
	sysConfig := PSysParameterMgr.GetSysParameter(for_game.LIMIT_PARAMETER)
	if !sysConfig.GetIsWithdrawal() {
		return easygo.NewFailMsg("提现入口已关闭，提现失败")
	}

	player := GetPlayerObj(who.GetPlayerId())
	if player == nil {
		return easygo.NewFailMsg("提现失败,不存在的玩家")
	}
	//8：00-22：00才可以提现
	_, _, _, h, _, _ := easygo.GetTimeData()
	if h < 8 || h >= 22 {
		return easygo.NewFailMsg("提现失败,每天8:00-22:00才可以提现")
	}
	if for_game.IS_FORMAL_SERVER {
		if who.GetPeopleId() == "" {
			res := "请先通过实名认证"
			return easygo.NewFailMsg(res)
		}
		if !who.GetIsPayPassword() {
			return easygo.NewFailMsg("提现失败,请先设置支付密码")
		}
	}
	bankId := reqMsg.GetAccountNo()

	if !util.InStringSlice(bankId, who.GetBankId()) {
		return easygo.NewFailMsg("提现失败,不存在这张银行卡")
	}
	bankInfo := who.GetBankMsg(bankId)
	if bankInfo == nil {
		return easygo.NewFailMsg("卡号未绑定，请绑定后在操作")
	}
	//当前激活的代付通道
	channel := PPayChannelMgr.GetCurPayChannel()
	logs.Info("代付channel----->%+v", channel)
	// 汇潮代付 需要判断是否有省份,城市
	if channel.GetId() == for_game.PAY_CHANNEL_HUICHAO_DF {
		if bankInfo.GetProvice() == "" || bankInfo.GetCity() == "" {
			return easygo.NewFailMsg("银行卡信息不完善", for_game.FAIL_MSG_CODE_1012)
		}
	}

	gold := reqMsg.GetAmount() //提现金额
	taxGold := int64(0)        //手续费
	//提现后台配置信息
	setting := PPayChannelMgr.GetCurPaymentSetting()

	if setting == nil {
		return easygo.NewFailMsg("系统错误无法提现,请联系客服")
	}
	period := for_game.GetPlayerPeriod(player.GetPlayerId())
	//每日提现总额限制
	curSum := period.DayPeriod.FetchInt64("OutSum")
	if gold+curSum > sysConfig.GetOutSum() {
		return easygo.NewFailMsg("提现失败，今天提现总额超" + easygo.AnytoA(float64(sysConfig.GetOutSum())/100.0) + "元")
	}
	//每日提现次数限制
	curTimes := period.DayPeriod.FetchInt32("OutTimes")
	if curTimes >= sysConfig.GetOutTimes() {
		return easygo.NewFailMsg("提现失败，今天提现次数已达" + easygo.AnytoA(sysConfig.GetOutTimes()) + "次")
	}
	//每笔最小限额
	if gold < sysConfig.GetWithdrawalMin() { //最小10元
		return easygo.NewFailMsg("提现失败,金额不能小于" + easygo.AnytoA(float64(sysConfig.GetWithdrawalMin())/100.0) + "元")
	}
	//每笔最大限额
	if gold > sysConfig.GetWithdrawalMax() {
		return easygo.NewFailMsg("提现失败,金额不能大于" + easygo.AnytoA(float64(sysConfig.GetWithdrawalMax())/100.0) + "元")
	}

	//需手续费:向上取整
	taxGold = int64(math.Ceil(float64(gold) * float64(setting.GetFeeRate()) / 1000.0)) //手续费: X>2000/995 && X> 最小起提:15

	if gold > who.GetGold() {
		return easygo.NewFailMsg("提现失败,账号金额不足" + easygo.AnytoA(float64(gold)/100.0) + "元")
	}

	code := bankInfo.GetBankCode()
	name := for_game.BankName[code]
	if for_game.IS_FORMAL_SERVER && name == "" {
		return easygo.NewFailMsg("不存在的银行卡号")
	}
	// 如果是汇聚去到,需要校验提现的银行卡
	if channel.GetId() == for_game.PAY_CHANNEL_HUIJU_DF {
		payNo := for_game.BankPayNo[code]
		if payNo == "" {
			return easygo.NewFailMsg("暂不支持" + for_game.BankName[code] + "提现")
		}
		reqMsg.BankCode = easygo.NewString(payNo)
	}

	info := map[string]interface{}{
		"BankId":   bankId[len(bankId)-4:],
		"BankName": name,
	}
	bankName := name + fmt.Sprintf("(%s)", bankId[len(bankId)-4:])
	//发起代付订单
	reqMsg.AccountProp = easygo.NewString("0")
	reqMsg.AccountName = easygo.NewString(player.GetRealName())
	reqMsg.AccountType = easygo.NewString("00")
	reqMsg.Tax = easygo.NewInt64(taxGold)
	var orderId string
	logs.Info("channel.GetId() ------------->", channel.GetId())
	//生成内部提现订单
	if channel.GetId() == for_game.PAY_CHANNEL_PENGJU {
		orderId = PWebPengJuPay.CreateOrder(reqMsg, player, setting)
	} else if channel.GetId() == for_game.PAY_CHANNEL_HUIJU_DF {
		logs.Info("汇聚提现")
		orderId = PWebHuiJuPay.CreateDFOrder(reqMsg, player, setting)
	} else if channel.GetId() == for_game.PAY_CHANNEL_HUICHAO_DF {
		logs.Info("汇潮代付")
		orderId = PWebHuiChaoPay.CreateDFOrder(reqMsg, player, setting)
	}

	// 体现预警通知
	easygo.Spawn(CheckWarningSMS, for_game.GOLD_CHANGE_TYPE_OUT, wAmount)

	if orderId == "" {
		//订单号未生成，提现失败
		logs.Error("代付订单id为空")
		return easygo.NewFailMsg("系统异常，提现下单失败")
	}
	reqMsg.OrderId = easygo.NewString(orderId)
	//扣除提现金额,分三部分：实际到账金额=提现金额-手续费(提现金额的0.008)-银行手续费(不同银行不同)
	reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_CASH_OUT, 0, info)
	msg := &share_message.GoldExtendLog{
		OrderId:  easygo.NewString(orderId),
		Title:    easygo.NewString(reason),
		PayType:  easygo.NewInt32(for_game.PAY_TYPE_GOLD),
		BankName: easygo.NewString(bankName),
		Gold:     easygo.NewInt64(who.GetGold() - gold),
	}
	NotifyAddGold(player.GetPlayerId(), -gold, "提现金额", for_game.GOLD_TYPE_CASH_OUT, msg)
	errStr := ""
	order := for_game.GetRedisOrderObj(orderId)
	//风控值检测
	if gold < sysConfig.GetRiskControl() {
		//第三方发起提现请求
		errMsg := &base.Fail{}
		if channel.GetId() == for_game.PAY_CHANNEL_PENGJU {
			errMsg = PWebPengJuPay.RechargeOrder(order.GetRedisOrder())
		} else if channel.GetId() == for_game.PAY_CHANNEL_HUIJU_DF {
			errMsg = PWebHuiJuPay.ReqSinglePay(order.GetRedisOrder())
		} else if channel.GetId() == for_game.PAY_CHANNEL_HUICHAO_DF {
			errMsg = PWebHuiChaoPay.TransferFixed(order.GetRedisOrder())
		}
		if errMsg.GetCode() == for_game.FAIL_MSG_CODE_SUCCESS {
			//第三方下单成功
			errStr = errMsg.GetReason()
			logs.Info("提现下单成功:")
			order.SetPayStatus(for_game.PAY_ST_DOING)
		} else {
			//渠道下单失败，如果已经生成订单，则转为人工处理
			order.SetChanneltype(for_game.CHANNEL_MAN_MAKE)
			errStr = errMsg.GetReason()
			logs.Info("提现第三方下单失败,转为人工处理订单:", orderId)
		}
	} else {
		//提现值超过风控，进入人工审核
		logs.Info("超过风控值:", sysConfig.GetRiskControl())
		//order.SetChanneltype(for_game.CHANNEL_MAN_MAKE)
		errStr = "提现额度超过风控值，进行人工审核"
	}
	order.SetExtendValue(errStr) // 成功/失败原因
	return nil
}

//第一次登录设置密码和手机号码
func (self *cls1) RpcFirstLoginSetInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.FirstInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("===========RpcFirstLoginSetInfo===========", reqMsg)
	msg := &client_hall.FirstReturnInfo{}
	pwd := reqMsg.GetPassword()
	sex := reqMsg.GetSex()
	name := reqMsg.GetNickName()
	signature := reqMsg.GetSignature()
	headIcon := reqMsg.GetHeadIcon()
	account := for_game.GetRedisAccountObj(who.GetPlayerId())
	if account == nil {
		panic("账号数据异常" + who.GetPhone())
	}
	if pwd != "" {
		account.SetPassword(pwd)
	}
	if sex != 0 {
		who.SetSex(sex)
	}
	if name != "" {
		who.SetNickName(name)
	}
	if signature != "" {
		who.SetSignature(signature)
	}
	if headIcon != "" {
		who.SetHeadIcon(headIcon)
	}
	label := reqMsg.GetLabel()
	//注册屏蔽标签页面
	/*if len(label) <= 0 {
		panic("标签数量怎么会小于等于0")
	}*/
	if len(label) > 0 {
		who.SetRedisLabelList(label)
		msg.Type = easygo.NewInt32(0)

		easygo.SortSliceInt32(label, true)
		var photo string
		var playTime int32
		// 判断 label 是否是大于1个
		if len(label) == 1 {
			msg := for_game.GetInterestTag(label[0])
			photo = msg.GetPopIcon()
			playTime = msg.GetPlayTime()
		} else {
			m := for_game.QueryInterestGroupByGroup(label)
			if m == nil {
				rand.Seed(time.Now().Unix())
				index := rand.Intn(len(label))
				id := label[index]
				msg := for_game.GetInterestTag(id)
				photo = msg.GetPopIcon()
				playTime = msg.GetPlayTime()

			} else {
				photo = m.GetPopIcon()
				playTime = m.GetPlayTime()
			}
		}
		newMsg := for_game.GetRecommendInfo(who.GetPlayerId(), 0)
		newMsg.PlayTime = easygo.NewInt32(playTime)
		newMsg.Photo = easygo.NewString(photo)
		ep.RpcReturnRecommendInfo(newMsg) //发送推荐的好友群信息
	}
	who.SetCreateTime()
	ep.RpcPlayerAttrChange(who.GetPlayerInfo())
	//把玩家补全信息存储到mongo
	who.SaveToMongo()

	//补全资料 算今天的注册人数
	easygo.Spawn(func() {
		for_game.MakePlayerKeepReport(who.GetPlayerId(), 2) //生成玩家留存报表 已优化到Redis
		for_game.SetRedisOperationChannelReportFildVal(easygo.Get0ClockTimestamp(util.GetMilliTime()), 1, who.GetChannel(), "ValidRegCount")
		for_game.SetRedisRegisterLoginReportFildVal(easygo.Get0ClockTimestamp(util.GetMilliTime()), 1, "ValidRegSumCount") //更新有效注册数
		amgr := for_game.GetRedisAccountObj(who.GetPlayerId())
		if amgr != nil && amgr.GetUnionId() != "" {
			for_game.SetRedisRegisterLoginReportFildVal(easygo.Get0ClockTimestamp(util.GetMilliTime()), 1, "ValidWxRegCount") //更新有效微信注册数
		} else {
			for_game.SetRedisRegisterLoginReportFildVal(easygo.Get0ClockTimestamp(util.GetMilliTime()), 1, "ValidPhoneRegCount") //更新有效手机号注册数
		}
	})

	// 通过手机号,获取当前的用户是否在抽奖邀请关联表中
	related := for_game.GetRelatedByPhoneFromDB(who.GetPhone())
	if related.GetId() != "" && related.GetPlayerId() != 0 { // 有关联关系
		easygo.Spawn(func() {
			// 添加抽奖次数进数据库和redis
			for_game.RedisLuckyPlayer.IncrLuckyCountToRedis(related.GetPlayerId(), 2)
			for_game.UpdateActivityReport("NewInviteCount", 1) //  邀请新增人数埋点.
		})
	}
	registerPush := GetRegisterPush(who, 2)
	ep.RpcRegisterPush(registerPush)
	return msg
}

// 推荐好友操作
func (self *cls1) RpcRecommendRequest(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.RecommendMsg, common ...*base.Common) easygo.IMessage {
	logs.Info("======RpcRecommendRequest===== reqMst=%v", reqMsg)
	playerIds := reqMsg.GetPlayerIds()
	who.SetIsRecommendOver(false)
	if len(playerIds) <= 0 {
		return nil
	}
	//fq := for_game.GetFriendBase(who.GetPlayerId())
	//t := share_message.AddFriend_Type_REGISTER
	//for _, pid := range playerIds {
	//	player := GetPlayerObj(pid)
	//	photos := player.GetPhoto()
	//	photo := ""
	//	if len(photos) > 0 {
	//		photo = photos[0]
	//	}
	//	nextId := for_game.NextId("Add_Friend")
	//	msg := &share_message.AddPlayerRequest{
	//		PlayerId:  easygo.NewInt64(pid),
	//		Time:      easygo.NewInt64(time.Now().Unix()),
	//		Text:      easygo.NewString(""),
	//		Type:      &t,
	//		Result:    easygo.NewInt32(for_game.ADDFRIEND_NODEAL), //未处理
	//		Id:        easygo.NewInt64(nextId),
	//		IsRead:    easygo.NewBool(false),
	//		Photo:     easygo.NewString(photo),
	//		Signature: easygo.NewString(player.GetSignature()),
	//		Sex:       easygo.NewInt32(player.GetSex()),
	//	}
	//	fq.AddFriendRequest(msg)
	//}
	//msg := fq.GetNewVersionAllFriendRequestForOne()
	//ep.RpcNoticeAddFriend(msg)
	//who.SetIsRecommendOver(false)
	//content := base64.StdEncoding.EncodeToString([]byte("你好，很高兴认识您！"))
	t := share_message.AddFriend_Type_REGISTER
	for _, pid := range playerIds {
		addMsg := &client_hall.AddPlayerInfo{
			PlayerId: easygo.NewInt64(pid),
			Type:     &t,
		}
		self.RpcAddFriend(ep, who, addMsg)
	}
	return nil
}

func (self *cls1) RpcFreshRecommendInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.RecommendRefreshInfo, common ...*base.Common) easygo.IMessage {
	t := reqMsg.GetType()
	msg := for_game.GetRecommendInfo(who.GetPlayerId(), t)
	return msg
}

//扫描二维码获取信息
func (self *cls1) RpcGetCodeInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CodeInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("======RpcGetCodeInfo========", reqMsg)
	teamId := reqMsg.GetTeamId()
	pid := reqMsg.GetPlayerId()
	playerId := who.GetPlayerId()
	if teamId != 0 { //群收款二维码
		teamObj := for_game.GetRedisTeamObj(teamId)
		if teamObj == nil {
			reqMsg.FailReason = easygo.NewString("该群不存在")
		}
		if for_game.GetTeamPlayerPos(teamId, playerId) == for_game.TEAM_UNUSE {
			reqMsg.FailReason = easygo.NewString("你不是该群成员，无法进行付款")
		}
		if pid != teamObj.GetTeamOwner() {
			reqMsg.FailReason = easygo.NewString("该群群主已变化,该二维码已不可用")
		}
		if reqMsg.GetCode() != teamObj.GetTeamQRCode() {
			reqMsg.FailReason = easygo.NewString("该群二维码已发生变化，当前二维码不可用")
		}
		if time.Now().Unix()-teamObj.GetTeamRefreshTime() > 3*86400 {
			reqMsg.FailReason = easygo.NewString("群收款二维码过期")
		}
		if reqMsg.GetFailReason() != "" {
			return reqMsg
		}
	}

	player := GetPlayerObj(pid)
	if player == nil {
		panic("玩家对象不存在")
	}
	reqMsg.Name = easygo.NewString(player.GetNickName())
	reqMsg.HeadIcon = easygo.NewString(player.GetHeadIcon())

	return reqMsg
}

//扫码付款
func (self *cls1) RpcPayForCode(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PayInfo, common ...*base.Common) easygo.IMessage {
	//检测备注是否包含敏感词
	sysParams := PSysParameterMgr.GetSysParameter(for_game.LIMIT_PARAMETER)
	if !sysParams.GetIsQRcode() {
		s := "二维码支付功能已关闭，暂时无法使用该功能"
		logs.Error(s)
		return easygo.NewFailMsg(s)
	}
	content := reqMsg.GetContent()
	if len(content) > 0 {
		evilType, _ := for_game.PDirtyWordsMgr.CheckWord(content)
		logs.Info(evilType)
		if evilType {
			return easygo.NewFailMsg("备注包含敏感字，请重新编辑")
		}
	}
	gold := reqMsg.GetGold()
	t := reqMsg.GetType()
	pid := reqMsg.GetPlayerId()
	receiver := GetPlayerObj(pid)
	//注销的账号不允许收钱
	if receiver == nil || receiver.GetStatus() == for_game.ACCOUNT_CANCELED {
		res := "玩家对象不存在"
		return easygo.NewFailMsg(res)
	}
	newMsg := &client_hall.PayForCodeInfo{}
	if t == for_game.PAY_TYPE_GOLD {
		if !who.GetIsPayPassword() {
			return easygo.NewFailMsg("请先设置支付密码")
		}
		if !who.CheckPayPassWord(for_game.Md5(reqMsg.GetPassword())) {
			return easygo.NewFailMsg("支付密码错误")
		}
		if gold > who.GetGold() {
			return easygo.NewFailMsg("零钱不足")
		}
	} else {
		if !reqMsg.GetIsWay() {
			ep.RpcTunedUpPayInfo(reqMsg.GetPayOrderInfo())
			newMsg.Type = easygo.NewInt32(2)
			return newMsg
		}
	}

	// 扣除付款金额
	//orderId, _ := for_game.PlaceOrder(who.GetPlayerId(), -gold, for_game.GOLD_TYPE_PAY_MONEY)
	orderId := reqMsg.GetOrderId()
	if orderId == "" {
		orderId = for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, for_game.GOLD_TYPE_PAY_MONEY)
	}

	reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_PAY_MONEY, pid, nil)
	msg := &share_message.GoldExtendLog{
		OrderId:      easygo.NewString(orderId),
		PayType:      easygo.NewInt32(t),
		TransferText: easygo.NewString(content),
		HeadIcon:     easygo.NewString(receiver.GetHeadIcon()),
		Title:        easygo.NewString(reason),
		Account:      easygo.NewString(receiver.GetAccount()),
	}
	if reqMsg.GetBankInfo() != "" {
		msg.BankName = easygo.NewString(reqMsg.GetBankInfo())
	}
	//if t == for_game.PAY_TYPE_GOLD {
	msg.Gold = easygo.NewInt64(who.GetGold() - gold)
	//}
	NotifyAddGold(who.GetPlayerId(), -gold, reason, for_game.GOLD_TYPE_PAY_MONEY, msg)

	//增加收款金额
	//orderId1, _ := for_game.PlaceOrder(pid, gold, for_game.GOLD_TYPE_GET_MONEY)
	orderId1 := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, for_game.GOLD_TYPE_GET_MONEY)
	reason1 := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_GET_MONEY, who.GetPlayerId(), nil)
	msg1 := &share_message.GoldExtendLog{
		OrderId:      easygo.NewString(orderId1),
		PayType:      easygo.NewInt32(t),
		TransferText: easygo.NewString(content),
		HeadIcon:     easygo.NewString(who.GetHeadIcon()),
		Title:        easygo.NewString(reason1),
		Gold:         easygo.NewInt64(receiver.GetGold() + gold),
		Account:      easygo.NewString(who.GetAccount()),
	}
	base := GetPlayerObj(pid)
	NotifyAddGold(base.GetPlayerId(), gold, reason1, for_game.GOLD_TYPE_GET_MONEY, msg1)
	ep1 := ClientEpMp.LoadEndpoint(pid)
	if ep1 != nil {
		msg := &client_hall.AddGoldInfo{
			Gold:     easygo.NewInt64(gold),
			Name:     easygo.NewString(receiver.GetNickName()),
			PlayerId: easygo.NewInt64(pid),
		}
		ep1.RpcNoticeAddGold(msg)
	}
	title := "收款到账通知"
	text := fmt.Sprintf("收到一笔二维码收款，收款金额￥%.2f。通过“我的”→“零钱”→“账单”可查看详情", float64(gold)/100)
	NoticeAssistant(pid, 1, title, text)

	newMsg.Type = easygo.NewInt32(1)
	return newMsg
}

//获取所有收藏信息索引
func (self *cls1) RpcGetCollectIndexList(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	lst := who.GetAllCollectIndexList()
	msg := &client_hall.CollectIndex{
		IndexList: lst,
	}
	logs.Info("RpcGetCollectIndexList", msg)
	return msg
}

//获取索引的收藏信息
func (self *cls1) RpcGetCollectIndexInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CollectIndex, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetCollectIndexInfo", reqMsg)
	indexList := reqMsg.GetIndexList()
	msg := who.GetCollectInfoForIndex(indexList)
	return msg
}

//分页获取收藏信息
func (self *cls1) RpcGetCollectInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.GetCollectInfo, common ...*base.Common) easygo.IMessage {
	page := reqMsg.GetPage()
	num := reqMsg.GetNum()
	t := reqMsg.GetType()
	info := who.GetPageCollectInfo(page, num, t)
	msg := &client_hall.AllCollectInfo{
		CollectInfo: info,
	}
	return msg
}

//增加收藏信息
func (self *cls1) RpcAddCollectInfo(ep IGameClientEndpoint, who *Player, reqMsg *share_message.CollectInfo, common ...*base.Common) easygo.IMessage {
	if len(who.GetCollectInfo()) >= 200 {
		return easygo.NewFailMsg("收藏已达上限，收藏失败")
	}
	msg := who.AddCollectInfo(reqMsg)
	return msg
}

//删除收藏信息
func (self *cls1) RpcDelCollectInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.DelCollectInfo, common ...*base.Common) easygo.IMessage {
	index := reqMsg.GetIndex()
	if index > who.GetMaxCollectIndex() {
		return easygo.NewFailMsg("删除失败，索引超出界限")
	}
	who.DelCollectInfo(index)
	return nil
}

//搜索收藏信息
func (self *cls1) RpcSearchCollectInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SearchCollectInfo, common ...*base.Common) easygo.IMessage {
	content := reqMsg.GetContent()
	info := who.GetSearchCollectInfo(content)
	msg := &client_hall.AllCollectInfo{
		CollectInfo: info,
	}
	return msg
}

//检查是否有未领取的红包和转账
func (self *cls1) RpcCheckUnGetMoney(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.UnGetMoneyInfo, common ...*base.Common) easygo.IMessage {
	targetId := reqMsg.GetTargetId()
	reqMsg.IsHave = easygo.NewBool(for_game.GetUnGetMoneyInfo(targetId, reqMsg.GetType()))
	return reqMsg
}

//软件切出后台
func (self *cls1) RpcCutBackstage(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcCutBackstage", who.GetPlayerId(), who.GetNickName())
	PlayerOnlineMgr.SetPlayerIsCutBackstage(who.GetPlayerId(), true)
	PlayerOnlineMgr.SetPlayerOutReturnTime(who.GetPlayerId(), "out", util.GetMilliTime())
	return nil
}

//软件切回app
func (self *cls1) RpcReturnApp(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcReturnApp", who.GetPlayerId(), who.GetNickName())
	PlayerOnlineMgr.SetPlayerIsCutBackstage(who.GetPlayerId(), false)
	return nil
}

func (self *cls1) RpcRequestCallInfo(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	callInfo := who.GetCallInfo()
	if callInfo == nil {
		return &share_message.CallInfo{}
	}
	return callInfo
}

//绑定微信
func (self *cls1) RpcOperateBindWechat(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.WechatInfo, common ...*base.Common) easygo.IMessage {
	_, _, unionId := for_game.GetWeChatInfo(reqMsg.GetCode(), reqMsg.GetApkCode())
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_ACCOUNT)
	defer closeFun()
	account := for_game.GetRedisAccountObj(who.GetPlayerId())
	t := reqMsg.GetType()
	if t == 1 { //绑定微信
		var msg *share_message.PlayerAccount
		err := col.Find(bson.M{"UnionId": unionId}).One(msg)
		if err != nil && err != mgo.ErrNotFound {
			panic(err)
		}
		if err == mgo.ErrNotFound { //数据库中找不到
			account.SetUnionId(unionId)
			err := col.Update(bson.M{"_id": who.GetPlayerId()}, bson.M{"$set": bson.M{"UnionId": unionId}}) //需要立马存库
			if err != nil || err != mgo.ErrNotFound {
				panic(err)
			}
			ep.RpcPlayerAttrChange(who.GetPlayerInfo())
			return nil
		}
		return easygo.NewFailMsg("该微信已经被绑定")
	} else { //解绑微信
		if account.GetUnionId() != "" {
			account.SetUnionId("")
			err := col.Update(bson.M{"_id": who.GetPlayerId()}, bson.M{"$set": bson.M{"UnionId": ""}}) //需要立马存库
			if err != nil || err != mgo.ErrNotFound {
				panic(err)
			}
			ep.RpcPlayerAttrChange(who.GetPlayerInfo())
			return nil
		}
		return easygo.NewFailMsg("该账号未绑定微信绑定")
	}
}

//绑定手机号
func (self *cls1) RpcOperateBindPhone(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.WechatInfo, common ...*base.Common) easygo.IMessage {
	phone := reqMsg.GetCode()
	account := for_game.GetRedisAccountByPhone(phone)
	if account != nil {
		return easygo.NewFailMsg("该手机号码已经被绑定")
	}

	if who.GetPhone() != "" {
		return easygo.NewFailMsg("该账号已经绑定过手机")
	}

	if for_game.IS_FORMAL_SERVER {
		phoneCode := reqMsg.GetPhoneCode()
		data := for_game.MessageMarkInfo.GetMessageMarkInfo(for_game.CLIENT_CODE_BINDPHONE, phone)
		if data == nil {
			return easygo.NewFailMsg("该手机号码获取不到验证码")
		}
		if data.Mark != phoneCode {
			return easygo.NewFailMsg("验证码错误")
		}
	}
	account = for_game.GetRedisAccountObj(who.GetPlayerId())
	account.SetAccount(phone)
	who.SetPhone(phone)
	ep.RpcPlayerAttrChange(who.GetPlayerInfo())

	return nil
}

//=============================chat.proto=============================
//撤回信息
func (self *cls1) RpcWithdrawMessage(ep IGameClientEndpoint, who *Player, reqMsg *client_server.ReadInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("===========RpcWithdrawMessage=================", reqMsg)
	chatType := reqMsg.GetType()
	logId := reqMsg.GetLogId()

	switch chatType {
	case for_game.CHAT_TYPE_PRIVATE: //私聊
		sessionId := for_game.MakeSessionKey(who.GetPlayerId(), reqMsg.GetPlayerId())
		b, ti := for_game.WithdrawMessage(who.GetPlayerId(), logId, sessionId)
		if !b {
			return easygo.NewFailMsg("该消息已超过5分钟，不能撤回")
		}
		friendId := reqMsg.GetPlayerId()
		newMsg := &client_hall.LogInfo{
			LogId:    easygo.NewInt64(logId),
			PlayerId: easygo.NewInt64(who.GetPlayerId()),
			Time:     easygo.NewInt64(ti),
			TargetId: easygo.NewInt64(friendId),
		}
		if PlayerOnlineMgr.CheckPlayerIsOnLine(friendId) {
			serverId := PlayerOnlineMgr.GetPlayerServerId(friendId)
			if serverId != PServerInfo.GetSid() { //不同大厅
				SendMsgToHallClientNew(friendId, "RpcWithdrawMessageResponse", newMsg)
			} else {
				ep1 := ClientEpMp.LoadEndpoint(friendId)
				if ep1 != nil {
					ep1.RpcWithdrawMessageResponse(newMsg)
				}
			}
		}
		return newMsg
		//ep.RpcWithdrawMessageResponse(newMsg)
	case for_game.CHAT_TYPE_TEAM: //群聊
		teamId := reqMsg.GetTeamId()
		teamObj := for_game.GetRedisTeamObj(teamId)
		if teamObj == nil {
			panic(fmt.Sprintf("群对象为空：%d", teamId))
		}
		if !util.Int64InSlice(who.GetPlayerId(), teamObj.GetTeamMemberList()) {
			return easygo.NewFailMsg("你不在这个群里不能撤回消息")
		}
		chatObj := for_game.GetRedisTeamChatLog(teamId)
		if !chatObj.TeamWithdrawChatLog(who.GetPlayerId(), logId) {
			return easygo.NewFailMsg("该消息已超过5分钟，不能撤回")
		}
		//通知其他成员
		NoticeTeamMessage(teamId, who.GetPlayerId(), for_game.WITHDRAW_MESSAGE, logId)
		//响应给自己
		ti := for_game.GetMillSecond()
		newMsg := &client_hall.LogInfo{
			LogId:    easygo.NewInt64(logId),
			PlayerId: easygo.NewInt64(who.GetPlayerId()),
			Time:     easygo.NewInt64(ti),
			TargetId: easygo.NewInt64(teamId),
		}
		return newMsg
	}
	return nil
}

//请求语音视屏聊天
func (self *cls1) RpcRequestSpecialChat(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SpecialChatInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcRequestSpecialChat", reqMsg)
	Type := reqMsg.GetType()
	chatType := reqMsg.GetChatType()
	targetId := reqMsg.GetTargetId()
	switch Type {
	case for_game.CHAT_TYPE_PRIVATE: //私聊
		target := GetPlayerObj(targetId)
		//   判断用户是否已注销
		if target == nil || target.GetStatus() == for_game.ACCOUNT_CANCELED {
			return easygo.NewFailMsg("该账号异常", for_game.FAIL_MSG_CODE_1010)
		}
		if !util.Int64InSlice(targetId, who.GetFriends()) {
			return easygo.NewFailMsg("对方不是你的好友", for_game.FAIL_MSG_CODE_1010)
		}
		if !util.Int64InSlice(who.GetPlayerId(), target.GetFriends()) {
			return easygo.NewFailMsg("你不是对方好友，操作失败", for_game.FAIL_MSG_CODE_1010)
		}
		if util.Int64InSlice(who.GetPlayerId(), target.GetBlackList()) {
			return easygo.NewFailMsg("通话请求已发送，但被对方拒收", for_game.FAIL_MSG_CODE_1010)
		}

		if who.GetCallInfoPlayerId() != 0 { //如果我正在打电话
			return easygo.NewFailMsg("你正在通话中，无法呼叫")
		}

		if target.GetCallInfoPlayerId() != 0 { //如果对方正在拨打语音视屏时
			chatLog := &share_message.Chat{
				SessionId:   easygo.NewString(for_game.MakeSessionKey(reqMsg.GetSendId(), targetId)),
				SourceId:    easygo.NewInt64(reqMsg.GetSendId()),
				TargetId:    easygo.NewInt64(targetId),
				ChatType:    easygo.NewInt32(Type),
				ContentType: easygo.NewInt32(chatType),
				Content:     easygo.NewString(base64.StdEncoding.EncodeToString([]byte("忙线中"))),
			}
			self.RpcChatNew(nil, who, chatLog)
			if target.GetCallInfoPlayerId() == who.GetPlayerId() {
				return easygo.NewFailMsg("对方忙线中，无法呼叫", for_game.FAIL_MSG_CODE_1009)
			}
			return easygo.NewFailMsg("对方忙线中，无法呼叫", for_game.FAIL_MSG_CODE_1010)
		}
		who.SetCallInfo(targetId, reqMsg)
		target.SetCallInfo(who.GetPlayerId(), reqMsg)
		serverId := PlayerOnlineMgr.GetPlayerServerId(targetId)
		if PlayerOnlineMgr.CheckPlayerIsOnLine(targetId) && !PlayerOnlineMgr.CheckPlayerIsCutBackstage(targetId) {
			if serverId != PServerInfo.GetSid() { //不同大厅
				SendMsgToHallClientNew(targetId, "RpcRequestSpecialChatResponse", reqMsg)
			} else {
				ep1 := ClientEpMp.LoadEndpoint(targetId)
				if ep1 != nil {
					ep1.RpcRequestSpecialChatResponse(reqMsg)
				}
			}
		} else {
			if PlayerOnlineMgr.CheckPlayerIsCutBackstage(targetId) { //如果切后台要下发协议
				if serverId != PServerInfo.GetSid() { //不同大厅
					SendMsgToHallClientNew(targetId, "RpcRequestSpecialChatResponse", reqMsg)
				} else {
					ep1 := ClientEpMp.LoadEndpoint(targetId)
					if ep1 != nil {
						ep1.RpcRequestSpecialChatResponse(reqMsg)
					}
				}
			}
			isNotice := target.GetIsNewMessage()
			if isNotice {
				f := func() {
					isShow := target.GetIsMessageShow()
					ids := for_game.GetJGIds([]int64{targetId})
					tid := easygo.AnytoA(who.GetPlayerId())
					var content string
					if isShow {
						if chatType == for_game.TALK_CONTENT_AUDIO {
							content = "给您发送语音聊天请求"
						} else if chatType == for_game.TALK_CONTENT_VIDEO {
							content = "给您发送视频聊天请求"
						}
					} else {
						content = "给您发送了一条消息"
					}

					by, _ := json.Marshal(reqMsg)
					m := for_game.PushMessage{
						Title:       who.GetNickName(),
						Content:     content,
						ContentType: for_game.JG_TYPE_PERSONALCHAT,
						TargetId:    tid,
						ChatType:    easygo.AnytoA(chatType),
						Msg:         string(by),
					}
					for_game.JGSendMessage(ids, m)
				}
				easygo.Spawn(f)
			}
		}
	case for_game.CHAT_TYPE_TEAM: //todo 群聊
	}
	return nil
}

//操作语音视屏聊天
func (self *cls1) RpcOperateSpecialChat(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SpecialChatInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("======RpcOperateSpecialChat=========", reqMsg)
	chatType := reqMsg.GetType()
	targetId := reqMsg.GetTargetId()
	sendId := reqMsg.GetSendId()
	operate := reqMsg.GetOperate()
	sender := GetPlayerObj(sendId)
	target := GetPlayerObj(targetId)
	if sender.GetCallInfoPlayerId() == 0 && target.GetCallInfoPlayerId() == 0 { //如果已经挂断
		//ep.RpcOperateSpecialChatResponse(reqMsg)
		SendMsgToHallClientNew(sendId, "RpcOperateSpecialChatResponse", reqMsg)
		SendMsgToHallClientNew(targetId, "RpcOperateSpecialChatResponse", reqMsg)
		logs.Info("强制进行双方挂断-----------------")
		return nil
	}
	switch chatType {
	case for_game.CHAT_TYPE_PRIVATE: //私聊
		var pid PLAYER_ID
		var b bool
		reqMsg.OperateId = easygo.NewInt64(who.GetPlayerId())
		if operate == 1 { //接受
			if who.GetPlayerId() == sendId {
				panic("发起的人怎么能接受")
			}
			pid = sendId
			reqMsg.Time = easygo.NewInt64(time.Now().Unix())
			target.SetCallInfo(sendId, reqMsg)
			sender.SetCallInfo(targetId, reqMsg)
		} else {
			chatLog := &share_message.Chat{
				SessionId:   easygo.NewString(for_game.MakeSessionKey(sendId, targetId)),
				SourceId:    easygo.NewInt64(sendId),
				TargetId:    easygo.NewInt64(targetId),
				ChatType:    easygo.NewInt32(chatType),
				ContentType: easygo.NewInt32(reqMsg.GetChatType()),
			}
			var content string
			if operate == 2 { //拒绝
				pid = sendId
				content = base64.StdEncoding.EncodeToString([]byte("已拒绝"))
			} else if operate == 3 { //挂断
				if targetId == who.GetPlayerId() { //如果是接听方挂断
					pid = sendId
				} else { //如果是拨打方挂断
					pid = targetId
				}
				startTime := reqMsg.GetTime()
				overTime := time.Now().Unix()
				ti := overTime - startTime
				sti := GetCallTime(ti)

				content = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("通话时长:%s", sti)))
				easygo.Spawn(for_game.SaveVideoVoiceDurationLog, reqMsg.GetChatType(), ti, sendId, targetId) //写语音视频时长日志
			} else if operate == 4 || operate == 5 { //超时或者主动取消拨打
				pid = targetId
				if operate == 5 {
					b = true
				}
				content = base64.StdEncoding.EncodeToString([]byte("已取消"))
			}
			target.SetCallInfo(0, nil)
			sender.SetCallInfo(0, nil)
			chatLog.Content = easygo.NewString(content)
			//self.RpcChat(ep, who, chatLog)

			player := GetPlayerObj(sendId)
			if player == nil {
				panic("对象为空")
			}
			self.RpcChatNew(nil, player, chatLog)
			//ep1 := ClientEpMp.LoadEndpoint(sendId)
			//if ep1 != nil {
			//	self.RpcChatNew(nil, player, chatLog)
			//} else { // 通知其他大厅发送 	self.RpcChat(ep1, player, chatLog)
			//	msg := &server_server.ChatToOtherReq{
			//		Chat:     chatLog,
			//		PlayerId: easygo.NewInt64(sendId),
			//	}
			//	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcChatToOther", msg)
			//}
		}
		if PlayerOnlineMgr.CheckPlayerIsOnLine(pid) && !PlayerOnlineMgr.CheckPlayerIsCutBackstage(pid) { //如果在线并且没有切后台
			SendMsgToHallClientNew(pid, "RpcOperateSpecialChatResponse", reqMsg)
		} else { //如果不在线或者切后台
			if PlayerOnlineMgr.CheckPlayerIsCutBackstage(pid) {
				SendMsgToHallClientNew(pid, "RpcOperateSpecialChatResponse", reqMsg)
			}
			if b { //如果被拨打的人或者后台 并且拨打的人主动挂断
				f := func() {
					ids := for_game.GetJGIds([]int64{pid})
					by, _ := json.Marshal(reqMsg)
					m := for_game.PushMessage{
						Title:       sender.GetNickName(),
						Content:     "对方已挂断",
						ContentType: for_game.JG_TYPE_PERSONALCHAT,
						TargetId:    easygo.AnytoA(pid),
						Msg:         string(by),
					}
					for_game.JGSendMessage(ids, m)
				}
				easygo.Spawn(f)
			}
		}
		ep.RpcOperateSpecialChatResponse(reqMsg)
	case for_game.CHAT_TYPE_TEAM: //todo 群聊
	}
	return nil

}

//阅读信息协议
func (self *cls1) RpcReadMessage(ep IGameClientEndpoint, who *Player, reqMsg *client_server.ReadInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("===RpcReadMessage====reqMsg=%v", reqMsg)
	chatType := reqMsg.GetType()
	if chatType != for_game.CHAT_TYPE_PRIVATE && chatType != for_game.CHAT_TYPE_TEAM {
		logs.Error("RpcReadMessage chatType有误: ", chatType)
		return easygo.NewFailMsg("参数有误")
	}
	logId := reqMsg.GetLogId()
	session := for_game.GetRedisChatSessionObj(reqMsg.GetSessionId())
	if session == nil {
		return easygo.NewFailMsg("无效的会话id")
	}
	switch chatType {
	case for_game.CHAT_TYPE_PRIVATE: //私聊 )
		session.SetReadInfo(who.GetPlayerId(), logId)
	case for_game.CHAT_TYPE_TEAM: //群聊
		//teamId := reqMsg.GetTeamId()
		//teamObj := for_game.GetRedisTeamObj(teamId)
		//if teamObj == nil {
		//	panic(fmt.Sprintf("群对象为空：%d", teamId))
		//}
		//if !util.Int64InSlice(who.GetPlayerId(), teamObj.GetTeamMemberList()) {
		//	return nil
		//}
		if logId > session.GetMaxLogId() {
			logId = session.GetMaxLogId()
		}
		memberObj := for_game.GetRedisTeamPersonalObj(reqMsg.GetTeamId())
		memberObj.ReadTeamChatLog(who.GetPlayerId(), logId)
	}
	return nil
}

func AddTeamMemberOperate(teamId, playerId, adminId int64, playerList []int64, name string) string {
	newPlist := make([]int64, 0)
	//去重复
	for _, id := range playerList {
		if easygo.Contain(newPlist, id) {
			continue
		}
		newPlist = append(newPlist, id)
	}
	teamObj := for_game.GetRedisTeamObj(teamId)
	//去重处理
	if !teamObj.CheckTeamMaxMember(len(newPlist)) {
		//return "添加成员数量超过该群聊数量上限"
		return "群成员超出上限"
	}
	var newList []int64
	memberList := teamObj.GetTeamMemberList()
	for _, pid := range newPlist {
		if util.Int64InSlice(pid, memberList) {
			continue
		}
		newList = append(newList, pid)
	}
	if len(newList) == 0 {
		return "被邀请人已在群里面"
	}
	serverId := PlayerOnlineMgr.GetPlayerServerId(playerId)
	//后台添加群成员没有邀请人playerId，不验证邀请人
	if adminId != 0 {
		m := &share_message.TeamChannel{
			Name: easygo.NewString("客服"),
			Type: easygo.NewInt32(5),
		}
		AddTeamMember(teamId, teamObj.GetSessionLogMaxId(), newList, false, 0, "群主邀请入群", serverId, m)
		NoticeTeamMessage(teamId, playerId, for_game.INVITE_PLAYER, newList, adminId)
		for _, pid := range newList {
			player := for_game.GetRedisPlayerBase(pid)
			player.AddTeamId(teamId)
		}
		return ""
	}

	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		return "你不在这个群里面"
	}

	pos := for_game.GetTeamPlayerPos(teamId, playerId)
	if pos == 0 {
		return "你在这个群里没有职位，群id:%d"
	}
	setting := teamObj.GetTeamMessageSetting()

	reason := fmt.Sprintf("%s邀请入群", name)
	m := &share_message.TeamChannel{
		Name: easygo.NewString(name),
		Type: easygo.NewInt32(4),
	}
	if setting.GetIsInvite() && pos > for_game.TEAM_MANAGER { //群聊邀请确认开启
		teamObj.AddTeamInvite(playerId, name, newList, reason, m)
		NoticeTeamMessage(teamId, playerId, for_game.REQUEST_ADDTEAM, newList)
	} else {
		AddTeamMember(teamId, teamObj.GetSessionLogMaxId(), newList, false, 0, reason, serverId, m)
		NoticeTeamMessage(teamId, playerId, for_game.INVITE_PLAYER, newList)
		for _, pid := range newList {
			player := for_game.GetRedisPlayerBase(pid)
			player.AddTeamId(teamId)
		}
	}
	return ""
}

//群里通过+号添加群成员
func (self *cls1) RpcAddTeamMember(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamReq, common ...*base.Common) easygo.IMessage {
	logs.Info("======RpcAddTeamMember==========", reqMsg)
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic(fmt.Sprintf("群对象为空：%d", teamId))
	}

	playerlist := reqMsg.GetPlayerIdList()
	if len(playerlist) == 0 {
		panic("playerlist怎么会是个空")
	}
	playerId := who.GetPlayerId()
	adminId := reqMsg.GetAdminId()
	var name string
	if who.GetPlayerId() == teamObj.GetTeamOwner() {
		name = "群主"
	} else {
		name = who.GetNickName()
	}

	reason := AddTeamMemberOperate(teamId, playerId, adminId, playerlist, name)
	if reason != "" {
		return easygo.NewFailMsg(reason)
	}
	return nil
}

func DeleteTeamMemberOperate(teamId, playerId int64, playerlist []int64) string {
	//后台删除群成员没有操作人playerId，不验证操作人权限
	if playerId == 0 {
		DeleteTeamMember(teamId, playerlist, 0)
		NoticeTeamMessage(teamId, playerId, for_game.DEL_PLAYER, playerlist)
		for _, pid := range playerlist {
			player := GetPlayerObj(pid)
			player.DelTeamId(teamId)
		}
		return ""
	}
	teamObj := for_game.GetRedisTeamObj(teamId)
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		return fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
	}
	pos := for_game.GetTeamPlayerPos(teamId, playerId)
	if pos == for_game.TEAM_UNUSE {
		return fmt.Sprintf("你在这个群里没有职位，群id:%d", teamId)
	}

	if util.Int64InSlice(teamObj.GetTeamOwner(), playerlist) {
		return "你不能踢出群主"

	}
	if pos >= for_game.TEAM_MASSES {
		return "你没有权限做这个操作"

	}
	serverId := PlayerOnlineMgr.GetPlayerServerId(playerId)
	success := DeleteTeamMember(teamId, playerlist, serverId)
	if !success {
		return "删除群成员失败"
	}
	NoticeTeamMessage(teamId, playerId, for_game.DEL_PLAYER, playerlist)
	for _, pid := range playerlist {
		player := GetPlayerObj(pid)
		player.DelTeamId(teamId)
	}
	return ""
}

//刪除群成员
func (self *cls1) RpcRemoveTeamMember(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamReq, common ...*base.Common) easygo.IMessage {
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}

	plst := reqMsg.GetPlayerIdList()
	if len(plst) == 0 {
		panic("playerlist怎么会是个空")
	}
	reason := DeleteTeamMemberOperate(teamId, who.GetPlayerId(), plst)
	if reason != "" {
		return easygo.NewFailMsg(reason)
	}
	return nil
}

//退出群
func (self *cls1) RpcExitTeam(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamInfo, common ...*base.Common) easygo.IMessage {
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		logs.Error("不存在的群")
		return easygo.NewFailMsg(fmt.Sprintf("无效的群id:%d", teamId))
	}
	playerId := who.GetPlayerId()
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		return easygo.NewFailMsg(fmt.Sprintf("你不在这个群里面，群id:%d", teamId))
	}
	serverId := PlayerOnlineMgr.GetPlayerServerId(playerId)
	plst := []int64{playerId}
	success := DeleteTeamMember(teamId, plst, serverId)
	if !success {
		return easygo.NewFailMsg("退出失败")
	}
	NoticeTeamMessage(teamId, playerId, for_game.EXIT_PLAYER, plst)
	who.DelTeamId(teamId)
	sessionId := easygo.AnytoA(teamId)
	who.DeletePlayerSessions([]string{sessionId}) //删除会话
	return nil
}

//解散群，前端用
func (self *cls1) RpcDefunctTeam(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.DefunctTeam, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcDefunctTeam", reqMsg, who.GetPlayerId())
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		logs.Error("不存在的群")
		return easygo.NewFailMsg(fmt.Sprintf("无效的群id:%d", teamId))
	}
	playerId := who.GetPlayerId()
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		return easygo.NewFailMsg(fmt.Sprintf("你不在这个群里面，群id:%d", teamId))
	}
	if playerId != teamObj.GetTeamOwner() {
		return easygo.NewFailMsg("你不是群主")
	}
	//解散群
	info := &share_message.OperatorInfo{
		Operator: easygo.NewString(who.GetAccount()),
		Time:     easygo.NewInt64(util.GetMilliTime()),
		Flag:     easygo.NewInt32(1),
	}

	teamObj.SetRedisOperatorInfo(info)
	teamObj.SetTeamStatus(for_game.DISSOLVE)
	teamObj.SetTeamDissolveTime(time.Now().Unix())
	memberList := teamObj.GetTeamMemberList()
	DeleteTeamMemberOperate(teamId, 0, memberList)
	reqMsg.Result = easygo.NewInt32(1)
	return reqMsg
}

//聊天协议
//func (self *cls1) RpcChat(ep IGameClientEndpoint, who *Player, reqMsg *share_message.Chat, common ...*base.Common) easygo.IMessage {
//	logs.Info("RpcChat：", who.GetPlayerId(), reqMsg)
//	easygo.Spawn(func() { for_game.MakePlayerBehaviorReport(1, who.GetPlayerId(), nil, nil, nil, nil) }) //生成用户行为报表聊天相关 已优化到Redis
//	reqMsg.Time = easygo.NewInt64(for_game.GetMillSecond())
//	chatType := reqMsg.GetChatType()
//	targetId := reqMsg.GetTargetId()
//	contentType := reqMsg.GetContentType()
//	//文本聊天文字检测
//	if contentType == for_game.TALK_CONTENT_WORD || contentType == for_game.TALK_CONTENT_CITE {
//		//evilType := TextModeration(reqMsg.GetContent(), who.GetPlayerId(), targetId, true)
//		cite := reqMsg.GetCite()
//		content, _ := base64.StdEncoding.DecodeString(reqMsg.GetContent())
//		isDirty, dWord := for_game.PDirtyWordsMgr.CheckWord(string(content))
//		if cite != "" && !isDirty {
//			contentData := easygo.KWAT{}
//			err := json.Unmarshal([]byte(cite), &contentData)
//			easygo.PanicError(err)
//			s := contentData.GetString("text")
//			isDirty, dWord = for_game.PDirtyWordsMgr.CheckWord(s)
//		}
//		if isDirty {
//			//屏蔽不让发送
//			reqMsg.IsSuccessSend = easygo.NewInt32(4)
//			reqMsg.EvilType = easygo.NewInt32(20001)
//			reqMsg.DirtyWord = easygo.NewString(dWord)
//			ep.RpcChatToClient(reqMsg)
//			return nil
//		}
//	}
//	//图片敏感检测
//	if contentType == for_game.TALK_CONTENT_IMAGE || contentType == for_game.TALK_CONTENT_EMOTION {
//		js, err := base64.StdEncoding.DecodeString(reqMsg.GetContent())
//		easygo.PanicError(err)
//		contentData := easygo.KWAT{}
//		err = json.Unmarshal(js, &contentData)
//		url := contentData.GetString("url")
//		url1 := contentData.GetString("tn_url")
//		ImageModeration(url, who.GetPlayerId(), targetId, true)
//		ImageModeration(url1, who.GetPlayerId(), targetId, true)
//		//if evilType != 100 { //敏感图
//		//	reqMsg.EvilType = easygo.NewInt32(evilType)
//		//	contentData["url"] = for_game.BAN_PICTURE_BIG
//		//	contentData["tn_url"] = for_game.BAN_PICTURE_SMALL
//		//	js1, _ := json.Marshal(contentData)
//		//	newCon := base64.StdEncoding.EncodeToString(js1)
//		//	reqMsg.Content = easygo.NewString(newCon)
//		//}
//	}
//	if contentType == for_game.TALK_CONTENT_PERSONAL_CARD { // 个人名片信息
//		// id 从聊天内容中获取
//		bytes, _ := base64.StdEncoding.DecodeString(reqMsg.GetContent())
//		type con struct {
//			UserId int64 `json:"userId"`
//		}
//		var c *con
//		_ = json.Unmarshal(bytes, &c)
//		//id := reqMsg.GetCardPlayerId()
//		id := c.UserId
//		var state int32
//		if util.Int64InSlice(id, who.GetFriends()) {
//			state = 1
//		} else {
//			state = 2
//		}
//		nickName := for_game.GetFriendsReName(who.GetPlayerId(), id)
//		cardInfo := GetCardInfo(id, state, nickName)
//		reqMsg.CardInfo = cardInfo
//	}
//	switch chatType {
//	case for_game.CHAT_TYPE_PRIVATE: //私聊
//		//先检测聊天内容
//		base := for_game.GetRedisPlayerBase(targetId)
//		//检测账号注销状态
//		if base.GetStatus() == for_game.ACCOUNT_CANCELED {
//			reqMsg.IsSuccessSend = easygo.NewInt32(5) // 注销账号.
//			ep.RpcChatToClient(reqMsg)
//			return nil
//		}
//		if util.Int64InSlice(who.GetPlayerId(), base.GetBlackList()) {
//			//res := "消息已发出，但被对方拒收了"
//			//return easygo.NewFailMsg(res, for_game.FAIL_MSG_CODE_1002)
//			reqMsg.IsSuccessSend = easygo.NewInt32(2)
//			ep.RpcChatToClient(reqMsg)
//			return nil
//		}
//		fq := for_game.GetFriendBase(targetId)
//		if !util.Int64InSlice(who.GetPlayerId(), fq.GetFriendIds()) && base.GetIsAddFriend() { //如果被对方删除并且对方开启好友验证
//			reqMsg.IsSuccessSend = easygo.NewInt32(1)
//			ep.RpcChatToClient(reqMsg)
//			return nil
//		}
//		name := for_game.GetFriendsReName(targetId, who.GetPlayerId())
//		if name == "" {
//			name = who.GetNickName()
//		}
//		logId := for_game.AddPersonalChatLog(reqMsg) //添加聊天日志记录
//		reqMsg.LogId = easygo.NewInt64(logId)
//		reqMsg.SourceHeadIcon = easygo.NewString(who.GetHeadIcon())
//		reqMsg.SourceName = easygo.NewString(base64.StdEncoding.EncodeToString([]byte(name)))
//		reqMsg.Types = easygo.NewInt32(who.GetTypes())
//		logs.Info("玩家在线状态:", PlayerOnlineMgr.CheckPlayerIsOnLine(targetId), PlayerOnlineMgr.CheckPlayerIsCutBackstage(targetId))
//		if PlayerOnlineMgr.CheckPlayerIsOnLine(targetId) && !PlayerOnlineMgr.CheckPlayerIsCutBackstage(targetId) { //如果玩家在线并且现在不处于后台状态
//			//player := PlayerMgr.LoadPlayer(targetId)
//			//if player != nil { //在同一个大厅
//			ep1 := ClientEpMp.LoadEndpoint(targetId)
//			if ep1 != nil {
//				ep1.RpcChatToClient(reqMsg)
//			} else {
//				SendMsgToHallClientNew(targetId, "RpcChatToClient", reqMsg)
//			}
//		} else { //不在线 或者切后台
//			if PlayerOnlineMgr.CheckPlayerIsCutBackstage(targetId) { //如果切后台要下发协议
//				ep1 := ClientEpMp.LoadEndpoint(targetId)
//				if ep1 != nil {
//					ep1.RpcChatToClient(reqMsg)
//				}
//			}
//			fr := for_game.GetFriendBase(targetId)
//			setting := fr.GetFriend(who.GetPlayerId()).GetSetting()
//			isNotice := base.GetIsNewMessage()
//			if !setting.GetIsNoDisturb() && isNotice {
//				f := func() {
//					isShow := base.GetIsMessageShow()
//					ids := for_game.GetJGIds([]int64{targetId})
//					tid := easygo.AnytoA(who.GetPlayerId())
//					var content string
//					if isShow {
//						if contentType == for_game.TALK_CONTENT_REDPACKET {
//							content = "给您发送了一个[红包]"
//						} else if contentType == for_game.TALK_CONTENT_IMAGE {
//							content = "给您发送了一张图片"
//						} else if contentType == for_game.TALK_CONTENT_TRANSFER_MONEY {
//							content = "[转账]"
//						} else if contentType == for_game.TALK_CONTENT_WORD || contentType == for_game.TALK_CONTENT_CITE {
//							s, _ := base64.StdEncoding.DecodeString(reqMsg.GetContent())
//							// todo 对表情进行处理.
//							content = for_game.ReplaceEmotionStr(string(s))
//							//content = string(s)
//						} else if contentType == for_game.TALK_CONTENT_SOUND {
//							content = "给您发送了一段语音"
//						} else if contentType == for_game.TALK_CONTENT_TEAM_CARD {
//							content = "群名片"
//						} else if contentType == for_game.TALK_CONTENT_EMOTION {
//							content = "给您发送了一个动画表情"
//						} else if contentType == for_game.TALK_CONTENT_PERSONAL_CARD { // 个人名片
//							content = "[名片]"
//
//						}
//					} else {
//						content = "给您发送了一条消息"
//					}
//					m := for_game.PushMessage{
//						Title:       name,
//						Content:     content,
//						ContentType: for_game.JG_TYPE_PERSONALCHAT,
//						TargetId:    tid,
//						ChatType:    easygo.AnytoA(contentType),
//						JumpObject:  4, // 私聊
//					}
//					for_game.JGSendMessage(ids, m)
//				}
//				easygo.Spawn(f)
//			}
//		}
//	case for_game.CHAT_TYPE_TEAM: //群聊
//		teamId := reqMsg.GetTargetId()
//		teamObj := for_game.GetRedisTeamObj(teamId)
//		if teamObj == nil {
//			panic("群对象怎么会是个空")
//		}
//
//		talkId := who.GetPlayerId()
//		if !util.Int64InSlice(talkId, teamObj.GetTeamMemberList()) {
//			return easygo.NewFailMsg("你不是本群群员")
//		}
//		memberObj := for_game.GetRedisTeamPersonalObj(teamId)
//		pos := for_game.GetTeamPlayerPos(teamId, talkId)
//
//		member := memberObj.GetTeamMember(who.GetPlayerId())
//		var closeTime int64
//		for _, value := range member.GetOperatorInfoPer() {
//			if closeTime < value.GetCloseTime() {
//				closeTime = value.GetCloseTime()
//			}
//		}
//		if member.GetStatus() == 2 && closeTime >= util.GetMilliTime() {
//			return nil
//		}
//
//		if teamObj.GetTeamMessageSetting().GetIsBan() {
//			return nil
//		}
//		if teamObj.GetTeamMessageSetting().GetIsStopTalk() && pos > for_game.TEAM_MANAGER &&
//			(contentType != for_game.TALK_CONTENT_SYSTEM && contentType != for_game.TALK_CONTENT_REDPACKET_LOG) {
//			reqMsg.IsSuccessSend = easygo.NewInt32(3)
//			ep.RpcChatToClient(reqMsg)
//			return nil
//		}
//
//		talker := GetPlayerObj(talkId)
//		notice := &share_message.NoticeMsg{}
//		noticeInfo := reqMsg.GetNoticeInfo()
//		if noticeInfo != nil {
//			notice.PlayerId = noticeInfo.GetPlayerId()
//			notice.IsAll = easygo.NewBool(notice.GetIsAll())
//		} else {
//			notice = nil
//		}
//
//		rename := memberObj.GetTeamMemberReName(talkId)
//		if rename == "" {
//			rename = talker.GetNickName()
//		}
//		reqMsg.SourceHeadIcon = easygo.NewString(talker.GetHeadIcon())
//		reqMsg.SourceName = easygo.NewString(rename)
//		reqMsg.Types = easygo.NewInt32(talker.GetTypes())
//		isRead := false
//		if PlayerOnlineMgr.CheckPlayerIsOnLine(talkId) {
//			isRead = true
//		}
//		session := for_game.GetRedisChatSessionObj(easygo.AnytoA(teamId))
//		logId := session.GetNextMaxLogId()
//		for_game.AddTeamChatLog(teamId, talkId, logId, reqMsg, notice, isRead)
//		reqMsg.SourceName = easygo.NewString(base64.StdEncoding.EncodeToString([]byte(rename)))
//		reqMsg.LogId = easygo.NewInt64(logId)
//
//		plst := teamObj.GetTeamMemberList()
//		var offLine, onLine []PLAYER_ID
//		for _, pid := range plst { //群发 每个人都发
//			if pid == talkId {
//				continue
//			}
//			if PlayerOnlineMgr.CheckPlayerIsOnLine(pid) { //如果在线
//				// 统计成字典 然后统一发
//				onLine = append(onLine, pid)
//				if PlayerOnlineMgr.CheckPlayerIsCutBackstage(pid) { //在软件主界面
//					offLine = append(offLine, pid)
//				}
//			} else { //不在线
//				offLine = append(offLine, pid)
//			}
//		}
//		BroadCastMsgToHallClientNew(onLine, "RpcChatToClient", reqMsg)
//		if contentType == for_game.TALK_CONTENT_REDPACKET_LOG {
//			//不走推送
//			ep.RpcChatToClient(reqMsg)
//			return nil
//		}
//		if len(offLine) != 0 {
//			fun := func() {
//				var showIds, noshowIds []int64 //所有离线的人的极光id
//				for _, id := range offLine {
//					player := for_game.GetRedisPlayerBase(id)
//					if player == nil {
//						logs.Error("不存在的玩家id:", id)
//						continue
//					}
//					isNotice := player.GetIsNewMessage()
//					if !isNotice {
//						continue
//					}
//					member := memberObj.GetTeamMember(id)
//					if member == nil {
//						continue
//					}
//					isDisturb := member.GetSetting().GetIsNoDisturb()
//					if isDisturb {
//						continue
//					}
//					isShow := player.GetIsMessageShow()
//					if isShow {
//						showIds = append(showIds, id)
//					} else {
//						noshowIds = append(noshowIds, id)
//					}
//				}
//				tid := easygo.AnytoA(reqMsg.GetTargetId())
//				var content string
//				if contentType == for_game.TALK_CONTENT_REDPACKET {
//					content = rename + ":" + "发送了一个[红包]"
//				} else if contentType == for_game.TALK_CONTENT_IMAGE {
//					content = rename + ":" + "发送了一张图片"
//				} else if contentType == for_game.TALK_CONTENT_WORD || contentType == for_game.TALK_CONTENT_CITE {
//					s, _ := base64.StdEncoding.DecodeString(reqMsg.GetContent())
//					//content = rename + ":" + string(s)
//					// 表情数字替换成表情
//					content = rename + ":" + for_game.ReplaceEmotionStr(string(s))
//				} else if contentType == for_game.TALK_CONTENT_SOUND {
//					content = rename + ":" + "发送了一段语音"
//				} else if contentType == for_game.TALK_CONTENT_TEAM_CARD {
//					content = rename + ":" + "群名片"
//				} else if contentType == for_game.TALK_CONTENT_EMOTION {
//					content = rename + ":" + "发送了一个动画表情"
//				} else if contentType == for_game.TALK_CONTENT_PERSONAL_CARD { // 个人名片
//					content = rename + ":" + "[名片]"
//				}
//				m := for_game.PushMessage{
//					Title:       teamObj.GetTeamName(),
//					ContentType: for_game.JG_TYPE_TEAMCHAT,
//					TargetId:    tid,
//					Content:     content,
//					ChatType:    easygo.AnytoA(contentType),
//					JumpObject:  5, // 群聊
//				}
//				if len(showIds) != 0 {
//					ids := for_game.GetJGIds(showIds)
//					for_game.JGSendMessage(ids, m)
//				}
//				if len(noshowIds) != 0 {
//					ids := for_game.GetJGIds(noshowIds)
//					m.Content = "收到一条群消息"
//					for_game.JGSendMessage(ids, m)
//				}
//			}
//			easygo.Spawn(fun)
//		}
//	default:
//		res := fmt.Sprintf("错误的聊天类型:%d", chatType)
//		return easygo.NewFailMsg(res)
//	}
//	if ep != nil {
//		ep.RpcChatToClient(reqMsg)
//	}
//	return nil
//}

//新版聊天协议
func (self *cls1) RpcChatNew(ep IGameClientEndpoint, who *Player, reqMsg *share_message.Chat, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcChatNew：", who.GetPlayerId(), reqMsg)
	defer logs.Info("RpcChatNew响应:", who.GetPlayerId(), reqMsg)
	easygo.Spawn(func() { for_game.MakePlayerBehaviorReport(1, who.GetPlayerId(), nil, nil, nil, nil) }) //生成用户行为报表聊天相关 已优化到Redis
	reqMsg.Time = easygo.NewInt64(for_game.GetMillSecond())
	chatType := reqMsg.GetChatType()
	if chatType != for_game.CHAT_TYPE_PRIVATE && chatType != for_game.CHAT_TYPE_TEAM && chatType != for_game.CHAT_TYPE_GROUP && chatType != for_game.CHAT_TYPE_NEARBY {
		logs.Error("RpcChatNew chatType 有误, chatType为: ", chatType)
		return easygo.NewFailMsg("参数有误")
	}
	targetId := reqMsg.GetTargetId()
	if targetId == 0 {
		logs.Error("RpcChatNew targetId 为空")
		return easygo.NewFailMsg("参数有误")
	}
	contentType := reqMsg.GetContentType()
	content, _ := base64.StdEncoding.DecodeString(reqMsg.GetContent())
	//文本聊天文字检测
	if contentType == for_game.TALK_CONTENT_WORD || contentType == for_game.TALK_CONTENT_CITE {
		//evilType := TextModeration(reqMsg.GetContent(), who.GetPlayerId(), targetId, true)
		cite := reqMsg.GetCite()
		isDirty, dWord := for_game.PDirtyWordsMgr.CheckWord(string(content))
		if cite != "" && !isDirty {
			contentData := easygo.KWAT{}
			err := json.Unmarshal([]byte(cite), &contentData)
			easygo.PanicError(err)
			s := contentData.GetString("text")
			isDirty, dWord = for_game.PDirtyWordsMgr.CheckWord(s)
		}
		if isDirty {
			//屏蔽不让发送
			reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_4)
			reqMsg.EvilType = easygo.NewInt32(20001)
			reqMsg.DirtyWord = easygo.NewString(dWord)
			return reqMsg
		}
	}
	//图片敏感检测
	if contentType == for_game.TALK_CONTENT_IMAGE || contentType == for_game.TALK_CONTENT_EMOTION {
		js, err := base64.StdEncoding.DecodeString(reqMsg.GetContent())
		easygo.PanicError(err)
		contentData := easygo.KWAT{}
		err = json.Unmarshal(js, &contentData)
		url := contentData.GetString("url")
		url1 := contentData.GetString("tn_url")
		ImageModeration(url, who.GetPlayerId(), targetId, true)
		ImageModeration(url1, who.GetPlayerId(), targetId, true)
	}
	if contentType == for_game.TALK_CONTENT_PERSONAL_CARD { // 个人名片信息
		// id 从聊天内容中获取
		bytes, _ := base64.StdEncoding.DecodeString(reqMsg.GetContent())
		type con struct {
			UserId int64 `json:"userId"`
		}
		var c *con
		_ = json.Unmarshal(bytes, &c)
		//id := reqMsg.GetCardPlayerId()
		id := c.UserId
		var state int32
		if util.Int64InSlice(id, who.GetFriends()) {
			state = 1
		} else {
			state = 2
		}
		nickName := for_game.GetFriendsReName(who.GetPlayerId(), id)
		cardInfo := GetCardInfo(id, state, nickName)
		reqMsg.CardInfo = cardInfo
	}
	//获取会话信息
	session := for_game.GetRedisChatSessionObj(reqMsg.GetSessionId())
	if session == nil {
		session = for_game.CreateRedisChatSessionObj(who.GetPlayerId(), reqMsg)
		session.SaveToMongo() //新建的会话，及时写入库
		if reqMsg.GetChatType() == for_game.CHAT_TYPE_PRIVATE {
			who.AddPlayerSession(reqMsg.GetSessionId())
			target := for_game.GetRedisPlayerBase(targetId)
			target.AddPlayerSession(reqMsg.GetSessionId())
		}
	}
	reqMsg.IsSuccessSend = easygo.NewInt32(0)
	switch chatType {
	case for_game.CHAT_TYPE_PRIVATE: //私聊
		//先检测聊天内容
		base := for_game.GetRedisPlayerBase(targetId)
		if base == nil {
			logs.Error("redisPlayerBase 为空,targetId: ", targetId)
			return easygo.NewFailMsg("参数有误")
		}
		//检测账号注销状态
		if base.GetStatus() == for_game.ACCOUNT_CANCELED {
			reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_5) // 注销账号.
			return reqMsg
		}
		if util.Int64InSlice(who.GetPlayerId(), base.GetBlackList()) {
			//res := "消息已发出，但被对方拒收了"
			//return easygo.NewFailMsg(res, for_game.FAIL_MSG_CODE_1002)
			reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_2)
			return reqMsg
		}
		//好友逻辑处理
		if !easygo.Contain(who.GetFriends(), targetId) {
			logs.Info("对方不是我好友")
			//如果不是好友，则把对方添加为好友
			//检测自己的好友是否已达上线
			myFriend := for_game.GetFriendBase(who.GetPlayerId())
			if !myFriend.CheckMaxNum() {
				//自己好友达上限
				reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_11)
				return reqMsg
			}
			//对方是否开启陌生人打招呼
			fq := for_game.GetFriendBase(targetId)
			if !fq.CheckMaxNum() {
				//对方好友达上限
				reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_12)
				return reqMsg
			}
			//对方开启不允许陌生人打招呼，不是对方好友才生效
			if !easygo.Contain(base.GetFriends(), who.GetPlayerId()) {
				if base.GetIsBanSayHi() {
					reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_7)
					return reqMsg
				}
				if reqMsg.GetSayType() == int32(share_message.AddFriend_Type_TEAM) && !base.GetIsTeamChat() {
					//对方开启不允许通过群打招呼
					reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_8)
					return reqMsg
				}
				if reqMsg.GetSayType() == int32(share_message.AddFriend_Type_CODE) && !base.GetIsCode() {
					//对方开启不允许通过二维码打招呼
					reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_9)
					return reqMsg
				}
				if reqMsg.GetSayType() == int32(share_message.AddFriend_Type_CARD) && !base.GetIsCard() {
					//对方开启不允许通过名片打招呼
					reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_10)
					return reqMsg
				}
			}
			myFriend.AddFriend(targetId, reqMsg.GetSayType())
			who.AddAttention(targetId)
			who.AddFans(targetId)
			for_game.OperateRedisDynamicAttention(1, who.GetPlayerId(), targetId)
			//如果
			if !easygo.Contain(base.GetFriends(), who.GetPlayerId()) {
				nextId := for_game.NextId("Add_Friend")
				t := share_message.AddFriend_Type(reqMsg.GetSayType())
				msgFriend := &share_message.AddPlayerRequest{
					PlayerId:  easygo.NewInt64(who.GetPlayerId()),
					Time:      easygo.NewInt64(time.Now().Unix()),
					Text:      easygo.NewString(string(content)),
					Type:      &t,
					Result:    easygo.NewInt32(for_game.ADDFRIEND_NODEAL), //未处理
					Id:        easygo.NewInt64(nextId),
					IsRead:    easygo.NewBool(false),
					NickName:  easygo.NewString(who.GetNickName()),
					HeadIcon:  easygo.NewString(who.GetHeadIcon()),
					Photo:     easygo.NewString(who.GetHeadIcon()),
					Signature: easygo.NewString(base.GetSignature()),
					Sex:       easygo.NewInt32(who.GetSex()),
				}
				fq.AddFriendRequest(msgFriend)
				SendMsgToHallClientNew(targetId, "RpcNoticeAddFriendNew", msgFriend)
				NoticeAssistant2(who, targetId, int32(t), 0, 0)
			}
			//msgAdd := myFriend.GetNewVersionAllFriendRequestForOne() //发送所有好友申请信息
			//SendMsgToHallClientNew(targetId, "RpcNoticeAddFriend", msgAdd)
			addMsg := GetFriendInfo(targetId)
			SendMsgToHallClientNew(who.GetPlayerId(), "RpcNoticeAgreeFriend", addMsg)
		}
		//互为好友时亲密度才增加1
		if easygo.Contain(base.GetFriends(), who.GetPlayerId()) {
			logs.Info("我是对方好友")
			for_game.AddPlayerIntimacy(reqMsg.GetSessionId(), int64(1))
			//通知双发亲密度增加
			obj := for_game.GetRedisPlayerIntimacyObj(reqMsg.GetSessionId())
			respMsg := obj.GetPlayerIntimacyToClient()
			BroadCastMsgToHallClientNew([]int64{who.GetPlayerId(), targetId}, "RpcChangeIntimacy", respMsg)
		} else { //我不是对方好友
			if reqMsg.GetContentType() != for_game.TALK_CONTENT_SAY_HI {
				param := PSysParameterMgr.GetSysParameter(for_game.COMMON_PARAMETER)
				var cn int32
				if param == nil {
					cn = for_game.STRANGER_MAX_TALK_NUM //默认值
				} else {
					cn = int32(param.GetStrangerChatCount())
				}
				logs.Info("当前限制次数:", cn)
				num := session.CheckCanStrangerTalk(who.GetPlayerId(), cn)
				if num >= cn-2 && num < cn {
					reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_13)
					reqMsg.ExtentValue = easygo.NewString(easygo.AnytoA(cn))
				} else if num >= cn {
					//陌生人单日聊天超3次
					reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_6)
					reqMsg.ExtentValue = easygo.NewString(easygo.AnytoA(cn))
					return reqMsg
				}
			}
		}
		name := for_game.GetFriendsReName(targetId, who.GetPlayerId())
		if name == "" {
			name = who.GetNickName()
		}
		logId := for_game.AddPersonalChatLog(reqMsg) //添加聊天日志记录
		reqMsg.LogId = easygo.NewInt64(logId)
		reqMsg.SourceHeadIcon = easygo.NewString(who.GetHeadIcon())
		reqMsg.SourceName = easygo.NewString(base64.StdEncoding.EncodeToString([]byte(name)))
		recvSessionData := session.GetChatSessionDataForClient(targetId)
		SendMsgToHallClientNew(targetId, "RpcChatNewSession", recvSessionData)
		if !PlayerOnlineMgr.CheckPlayerIsOnLine(targetId) || PlayerOnlineMgr.CheckPlayerIsCutBackstage(targetId) {
			//如果玩家不在线或者切后台，需要推送通知
			fr := for_game.GetFriendBase(targetId)
			setting := fr.GetFriend(who.GetPlayerId()).GetSetting()
			isNotice := base.GetIsNewMessage()
			if !setting.GetIsNoDisturb() && isNotice {
				f := func() {
					isShow := base.GetIsMessageShow()
					ids := for_game.GetJGIds([]int64{targetId})
					tid := easygo.AnytoA(who.GetPlayerId())
					content := "给您发送了一条消息"
					if isShow {
						if contentType == for_game.TALK_CONTENT_REDPACKET {
							content = "给您发送了一个[红包]"
						} else if contentType == for_game.TALK_CONTENT_IMAGE {
							content = "给您发送了一张图片"
						} else if contentType == for_game.TALK_CONTENT_TRANSFER_MONEY {
							content = "[转账]"
						} else if contentType == for_game.TALK_CONTENT_WORD || contentType == for_game.TALK_CONTENT_CITE {
							s, _ := base64.StdEncoding.DecodeString(reqMsg.GetContent())
							// todo 对表情进行处理.
							content = for_game.ReplaceEmotionStr(string(s))
							//content = string(s)
						} else if contentType == for_game.TALK_CONTENT_SOUND {
							content = "给您发送了一段语音"
						} else if contentType == for_game.TALK_CONTENT_TEAM_CARD {
							content = "群名片"
						} else if contentType == for_game.TALK_CONTENT_EMOTION {
							content = "给您发送了一个动画表情"
						} else if contentType == for_game.TALK_CONTENT_PERSONAL_CARD { // 个人名片
							content = "[名片]"
						}
					}
					m := for_game.PushMessage{
						Title:       name,
						Content:     content,
						ContentType: for_game.JG_TYPE_PERSONALCHAT,
						TargetId:    tid,
						ChatType:    easygo.AnytoA(contentType),
					}
					for_game.JGSendMessage(ids, m)
				}
				easygo.Spawn(f)
			}
		}
	case for_game.CHAT_TYPE_TEAM: //群聊
		teamId := reqMsg.GetTargetId()
		teamObj := for_game.GetRedisTeamObj(teamId)
		if teamObj == nil {
			panic("群对象怎么会是个空")
		}

		talkId := who.GetPlayerId()
		if !util.Int64InSlice(talkId, teamObj.GetTeamMemberList()) {
			return easygo.NewFailMsg("你不是本群群员")
		}
		memberObj := for_game.GetRedisTeamPersonalObj(teamId)
		pos := for_game.GetTeamPlayerPos(teamId, talkId)

		member := memberObj.GetTeamMember(who.GetPlayerId())
		var closeTime int64
		for _, value := range member.GetOperatorInfoPer() {
			if closeTime < value.GetCloseTime() {
				closeTime = value.GetCloseTime()
			}
		}
		if member.GetStatus() == 2 && closeTime >= util.GetMilliTime() {
			return easygo.NewFailMsg("群禁言状态")
		}

		if teamObj.GetTeamMessageSetting().GetIsBan() {
			return easygo.NewFailMsg("群禁言状态")
		}
		if teamObj.GetTeamMessageSetting().GetIsStopTalk() && pos > for_game.TEAM_MANAGER &&
			(contentType != for_game.TALK_CONTENT_SYSTEM && contentType != for_game.TALK_CONTENT_REDPACKET_LOG) {
			reqMsg.IsSuccessSend = easygo.NewInt32(for_game.CHAT_REFUSE_TYPE_3)
			return reqMsg
		}

		talker := GetPlayerObj(talkId)
		notice := &share_message.NoticeMsg{}
		noticeInfo := reqMsg.GetNoticeInfo()
		if noticeInfo != nil {
			notice.PlayerId = noticeInfo.GetPlayerId()
			notice.IsAll = easygo.NewBool(noticeInfo.GetIsAll())
		} else {
			notice = nil
		}

		rename := memberObj.GetTeamMemberReName(talkId)
		if rename == "" {
			rename = talker.GetNickName()
		}
		reqMsg.SourceHeadIcon = easygo.NewString(talker.GetHeadIcon())
		reqMsg.SourceName = easygo.NewString(rename)
		isRead := false
		//if PlayerOnlineMgr.CheckPlayerIsOnLine(talkId) {
		//	isRead = true
		//}
		//logId := teamObj.GetNextLogMaxId()
		logId := session.GetNextMaxLogId()
		for_game.AddTeamChatLog(teamId, talkId, logId, reqMsg, notice, isRead)
		reqMsg.LogId = easygo.NewInt64(logId)
		plst := teamObj.GetTeamMemberList()
		var offLine, onLine []PLAYER_ID
		for _, pid := range plst { //群发 每个人都发
			if pid == talkId {
				continue
			}
			if PlayerOnlineMgr.CheckPlayerIsOnLine(pid) { //如果在线
				// 统计成字典 然后统一发
				onLine = append(onLine, pid)
				if PlayerOnlineMgr.CheckPlayerIsCutBackstage(pid) { //在软件主界面
					offLine = append(offLine, pid)
				}
			} else { //不在线
				offLine = append(offLine, pid)
			}
		}
		sendData := session.GetChatSessionDataForClient(0)
		BroadCastMsgToHallClientNew(onLine, "RpcChatNewSession", sendData)
		if contentType == for_game.TALK_CONTENT_REDPACKET_LOG {
			//领红包不走推送
			//ep.RpcChatToClient(reqMsg)
			reqMsg.IsSuccessSend = easygo.NewInt32(0)
			if ep == nil {
				sendSessionData := session.GetChatSessionDataForClient(who.GetPlayerId())
				SendMsgToHallClientNew(who.GetPlayerId(), "RpcChatNewSession", sendSessionData)
			}
			return reqMsg
		}
		if len(offLine) != 0 {
			fun := func() {
				var showIds, noshowIds []int64 //所有离线的人的极光id
				for _, id := range offLine {
					player := for_game.GetRedisPlayerBase(id)
					if player == nil {
						logs.Error("不存在的玩家id:", id)
						continue
					}
					isNotice := player.GetIsNewMessage()
					if !isNotice {
						continue
					}
					member := memberObj.GetTeamMember(id)
					if member == nil {
						continue
					}
					isDisturb := member.GetSetting().GetIsNoDisturb()
					if isDisturb {
						continue
					}
					isShow := player.GetIsMessageShow()
					if isShow {
						showIds = append(showIds, id)
					} else {
						noshowIds = append(noshowIds, id)
					}
				}
				tid := easygo.AnytoA(reqMsg.GetTargetId())
				var content string
				if contentType == for_game.TALK_CONTENT_REDPACKET {
					content = rename + ":" + "发送了一个[红包]"
				} else if contentType == for_game.TALK_CONTENT_IMAGE {
					content = rename + ":" + "发送了一张图片"
				} else if contentType == for_game.TALK_CONTENT_WORD || contentType == for_game.TALK_CONTENT_CITE {
					s, _ := base64.StdEncoding.DecodeString(reqMsg.GetContent())
					//content = rename + ":" + string(s)
					// 表情数字替换成表情
					content = rename + ":" + for_game.ReplaceEmotionStr(string(s))
				} else if contentType == for_game.TALK_CONTENT_SOUND {
					content = rename + ":" + "发送了一段语音"
				} else if contentType == for_game.TALK_CONTENT_TEAM_CARD {
					content = rename + ":" + "群名片"
				} else if contentType == for_game.TALK_CONTENT_EMOTION {
					content = rename + ":" + "发送了一个动画表情"
				} else if contentType == for_game.TALK_CONTENT_PERSONAL_CARD { // 个人名片
					content = rename + ":" + "[名片]"
				}
				m := for_game.PushMessage{
					Title:       teamObj.GetTeamName(),
					ContentType: for_game.JG_TYPE_TEAMCHAT,
					TargetId:    tid,
					Content:     content,
					ChatType:    easygo.AnytoA(contentType),
				}
				if len(showIds) != 0 {
					ids := for_game.GetJGIds(showIds)
					for_game.JGSendMessage(ids, m)
				}
				if len(noshowIds) != 0 {
					ids := for_game.GetJGIds(noshowIds)
					m.Content = "收到一条群消息"
					for_game.JGSendMessage(ids, m)
				}
			}
			easygo.Spawn(fun)
		}
	default:
		res := fmt.Sprintf("错误的聊天类型:%d", chatType)
		return easygo.NewFailMsg(res)
	}
	if ep == nil {
		sendSessionData := session.GetChatSessionDataForClient(who.GetPlayerId())
		SendMsgToHallClientNew(who.GetPlayerId(), "RpcChatNewSession", sendSessionData)
	}
	return reqMsg
}

//获取群设置主界面信息
func (self *cls1) RpcGetTeamSetting(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamInfo, common ...*base.Common) easygo.IMessage {
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}
	playerId := reqMsg.GetPlayerId()
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}
	msg := for_game.GetTeamSettingInfo(teamId, playerId)
	ep.RpcTeamSettingResponse(msg)
	return nil
}

//获取群设置主界面信息
func (self *cls1) RpcGetTeamSettingNew(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamInfo, common ...*base.Common) easygo.IMessage {
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}
	playerId := reqMsg.GetPlayerId()
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}
	msg := for_game.GetTeamSettingInfoEx(teamId, playerId)
	return msg
}

//获取群设置管理信息
func (self *cls1) RpcGetTeamManageSetting(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("========RpcGetTeamManageSetting========reqMsg=%v", reqMsg)
	teamId := reqMsg.GetTeamId()
	playerId := reqMsg.GetPlayerId()

	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}
	if for_game.GetTeamPlayerPos(teamId, playerId) > for_game.TEAM_MANAGER {
		res := "你没有权限操作"
		return easygo.NewFailMsg(res)
	}
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}
	msg := for_game.GetTeamManageSetting(teamId)
	ep.RpcTeamManageSettingResponse(msg)
	return nil
}
func (self *cls1) RpcGetTeamManageSettingNew(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("========RpcGetTeamManageSettingNew========reqMsg=%v", reqMsg)
	teamId := reqMsg.GetTeamId()
	playerId := reqMsg.GetPlayerId()

	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}
	if for_game.GetTeamPlayerPos(teamId, playerId) > for_game.TEAM_MANAGER {
		res := "你没有权限操作"
		return easygo.NewFailMsg(res)
	}
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}
	msg := for_game.GetTeamManageSettingEx(teamId)
	return msg
}

//发送红包
func (self *cls1) RpcSendRedPacket(ep IGameClientEndpoint, who *Player, reqMsg *share_message.RedPacket, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcSendRedPacket：", reqMsg)
	targetId := reqMsg.GetTargetId()
	redPacketType := reqMsg.GetType()
	//检测目标人id是否为0
	if targetId == 0 {
		return easygo.NewFailMsg("无效的TargetId")
	}
	// 私包,判断对方是否注销
	if redPacketType == for_game.RED_PACKET_PERSONAL {
		other := for_game.GetRedisPlayerBase(targetId)
		if other == nil || other.GetStatus() == for_game.ACCOUNT_CANCELED {
			return easygo.NewFailMsg("该账号异常", for_game.FAIL_MSG_CODE_1016)
		}
	}

	// 检测目标人是否存在.
	sysParams := PSysParameterMgr.GetSysParameter(for_game.LIMIT_PARAMETER)
	if !sysParams.GetIsRedPacket() {
		s := "红包功能已关闭，暂时无法使用该功能"
		logs.Error(s)
		return easygo.NewFailMsg(s)
	}

	//检测备注是否包含敏感词
	mark := reqMsg.GetContent()
	if len(mark) > 0 {
		evilType, _ := for_game.PDirtyWordsMgr.CheckWord(mark)
		if evilType {
			return easygo.NewFailMsg("备注包含敏感字，请重新编辑")
		}
	}
	if redPacketType == for_game.RED_PACKET_PERSONAL {
		other := for_game.GetRedisPlayerBase(targetId)

		//如果对方不是自己的好友
		if !util.Int64InSlice(targetId, who.GetFriends()) {
			res := "对方不是你的好友，操作失败"
			return easygo.NewFailMsg(res)
		}

		//如果对方在自己的黑名单
		if util.Int64InSlice(targetId, who.GetBlackList()) {
			res := "对方已被加入黑名单，操作失败"
			return easygo.NewFailMsg(res)
		}

		//如果自己不是对方好友
		if !util.Int64InSlice(who.GetPlayerId(), other.GetFriends()) {
			res := "你不是对方好友，操作失败"
			return easygo.NewFailMsg(res)
		}

		//如果自己在对方的黑名单中
		if util.Int64InSlice(who.GetPlayerId(), other.GetBlackList()) {
			res := "你已被对方拉黑，操作失败"
			return easygo.NewFailMsg(res)
		}
	}

	if reqMsg.GetPayWay() == for_game.PAY_TYPE_GOLD {
		if !who.GetIsPayPassword() {
			return easygo.NewFailMsg("请先设置支付密码")
		}
		if !who.CheckPayPassWord(for_game.Md5(reqMsg.GetPayPassWord())) {
			return easygo.NewFailMsg("支付密码错误")
		}
	}

	if reqMsg.GetTotalCount() < 1 {
		res := fmt.Sprintf("红包数量必须大于等于1")
		return easygo.NewFailMsg(res)
	}

	if redPacketType == for_game.RED_PACKET_PERSONAL && reqMsg.GetTotalMoney() > 20000 {
		return easygo.NewFailMsg("单个红包金额不可超过200.00元")
	}
	if reqMsg.GetTotalCount() > 100 {
		return easygo.NewFailMsg("红包数量不可超过100个")
	}
	perNum := float64(reqMsg.GetTotalMoney()) / float64(reqMsg.GetTotalCount())
	if redPacketType == for_game.RED_PACKET_TEAM_LUCKEY && perNum > 20000 {
		return easygo.NewFailMsg("单个红包金额不可超过200.00元")
	}
	maxNum := int64(for_game.RED_PACKET_MIN_VALUE) * int64(reqMsg.GetTotalCount())
	if maxNum > reqMsg.GetTotalMoney() {
		res := fmt.Sprintf("单个红包不能小于%.2f", for_game.RED_PACKET_MIN_VALUE/100.0)
		return easygo.NewFailMsg(res)
	}

	if redPacketType == for_game.RED_PACKET_TEAM_LUCKEY { //群红包判断是否开启禁言
		teamId := reqMsg.GetTargetId()
		teamObj := for_game.GetRedisTeamObj(teamId)
		if teamObj == nil {
			panic("群对象怎么会是个空")
		}
		pos := for_game.GetTeamPlayerPos(teamId, who.GetPlayerId())
		if teamObj.GetTeamMessageSetting().GetIsStopTalk() && pos > for_game.TEAM_MANAGER {
			return easygo.NewFailMsg("本群禁止群聊天")
		}
	}
	if reqMsg.GetPayWay() == for_game.PAY_TYPE_GOLD {
		if who.GetGold() < reqMsg.GetTotalMoney() {
			res := fmt.Sprintf("余额不足%d", reqMsg.GetTotalMoney()/100)
			return easygo.NewFailMsg(res)
		}
	} else {
		ep.RpcTunedUpPayInfo(reqMsg.GetPayOrderInfo()) //通知客户端调起支付
		return nil
	}
	chatLog := CreateRedPacket(reqMsg)
	self.RpcChatNew(nil, who, chatLog)
	return nil
}

func CreateRedPacket(reqMsg *share_message.RedPacket, orId ...string) *share_message.Chat {
	bMsg := ChooseOneBackStage("RpcDealReaPacket", reqMsg) //负载均衡选择一台后台服
	// 务器做拆包处理
	//logs.Info("后台生成红包:", bMsg)
	if bMsg == nil {
		panic("后台服务器出错")
	}
	orderId := append(orId, "")[0]
	if orderId == "" {
		//如果没有订单，则生成一个
		orderId = for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, for_game.GOLD_TYPE_SEND_REDPACKET)
	}
	newMsg, ok := bMsg.(*share_message.RedPacket) //数据强转
	if !ok {
		logs.Info("bMsg:", bMsg)
		panic("后台返回红包格式异常")
	}
	newMsg.OrderId = easygo.NewString(orderId)
	redPacket := for_game.CreateRedisRedPacket(newMsg)
	// for_game.RedPacketMgr.Store(redPacket.GetId(), redPacket)
	//logs.Info("产生的红包:", redPacket)
	content := redPacket.GetRedisSendRedPacketBase64Str()
	//扣除发送金额
	t := reqMsg.GetType()
	info := map[string]interface{}{
		"RedType": t,
	}
	reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_SEND_REDPACKET, reqMsg.GetTargetId(), info)
	msg := &share_message.GoldExtendLog{
		RedPacketId: easygo.NewInt64(redPacket.GetId()),
		OrderId:     easygo.NewString(orderId),
		PayType:     easygo.NewInt32(reqMsg.GetPayWay()),
		Title:       easygo.NewString(reason),
	}
	if reqMsg.GetBankCard() != "" {
		msg.BankName = easygo.NewString(reqMsg.GetBankCard())
	}
	base := GetPlayerObj(reqMsg.GetSender())
	//if reqMsg.GetPayWay() == for_game.PAY_TYPE_GOLD {
	msg.Gold = easygo.NewInt64(base.GetGold() - reqMsg.GetTotalMoney())
	//}
	if redPacket.GetType() == for_game.RED_PACKET_PERSONAL {
		tarBase := GetPlayerObj(reqMsg.GetTargetId())
		if tarBase != nil {
			msg.Account = easygo.NewString(tarBase.GetAccount())
		}
	}
	NotifyAddGold(base.GetPlayerId(), -reqMsg.GetTotalMoney(), reason, for_game.GOLD_TYPE_SEND_REDPACKET, msg)
	//FinishOrder(orderId, msg)
	for_game.UpdateRedisRedPacketTotal(reqMsg.GetSender(), for_game.GetMillSecond(), reqMsg.GetTotalMoney(), for_game.REDPACKET_STATISTICS_SEND)

	// 判断是否在集卡时间范围内
	tt := time.Now().Unix()
	act := for_game.GetActivityFromDB(1)
	if act.GetStartTime() <= tt && act.GetEndTime() > tt {
		// 发红包添加抽奖机会
		easygo.Spawn(func() {
			period := GetPlayerPeriod(reqMsg.GetSender())
			if !period.DayPeriod.FetchBool(for_game.LUCKY_PLAYER_IS_SENDREDPACK) {
				// 添加抽奖次数进数据库和redis
				for_game.RedisLuckyPlayer.IncrLuckyCountToRedis(reqMsg.GetSender(), 2)
				// 设置今天已发红包
				period.DayPeriod.Set(for_game.LUCKY_PLAYER_IS_SENDREDPACK, true)
				for_game.UpdateActivityReport("TaskRedPackCount", 1) // 完成发红包任务人数埋点
			}
		})
	}

	var sessionId string
	if t == for_game.CHAT_TYPE_PRIVATE {
		sessionId = for_game.MakeSessionKey(reqMsg.GetSender(), reqMsg.GetTargetId())
	} else if t == for_game.CHAT_TYPE_TEAM {
		sessionId = easygo.AnytoA(reqMsg.GetTargetId())
	}
	chatLog := &share_message.Chat{
		SessionId:   easygo.NewString(sessionId),
		SourceId:    easygo.NewInt64(reqMsg.GetSender()),
		TargetId:    easygo.NewInt64(reqMsg.GetTargetId()),
		Content:     easygo.NewString(content),
		ChatType:    easygo.NewInt32(reqMsg.GetType()),
		ContentType: easygo.NewInt32(for_game.TALK_CONTENT_REDPACKET),
	}

	easygo.Spawn(func() { for_game.MakePlayerBehaviorReport(7, 0, reqMsg, nil, nil, nil) }) //生成用户行为报表发红包相关 已优化到Redis

	return chatLog
}

//查看红包
func (self *cls1) RpcCheckRedPacket(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.OpenRedPacket, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcCheckRedPacket：", reqMsg)
	obj := for_game.GetRedisRedPacket(reqMsg.GetId())
	if obj == nil {
		return easygo.NewFailMsg("无效的红包id")
	}
	redPacket := obj.GetRedisRedPacket()

	if redPacket == nil {
		return easygo.NewFailMsg("无效的红包id")
	}
	//通知领玩家的人结果

	ep.RpcCheckRedPacketResult(redPacket)

	return easygo.EmptyMsg
}

//领取红包
func (self *cls1) RpcOpenRedPacket(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.OpenRedPacket, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcOpenRedPacket：", reqMsg)
	redPacket := for_game.GetRedisRedPacket(reqMsg.GetId())
	if redPacket == nil {
		logs.Info("没有该红包: reqMsg.GetId():", reqMsg.GetId())
		return nil
	}
	if redPacket.GetType() == for_game.RED_PACKET_TEAM_LUCKEY {
		//群红包看是否是群成员
		teamObj := for_game.GetRedisTeamObj(redPacket.GetTargetId())
		if teamObj != nil {
			if !easygo.Contain(teamObj.GetTeamMemberList(), who.GetPlayerId()) {
				logs.Info("不是群成员，无法领取红包", reqMsg.GetId())
				return nil
			}
		}
	}
	curTime := util.GetMilliTime()
	val, err := redPacket.RedisOpen(who.GetPlayerId(), curTime)
	//redPacket.Logs = for_game.GetRedisRedPacketLogList(reqMsg.GetId()) //领完红包以后带上领取日志
	sendRedPacket := redPacket.GetRedisRedPacket()
	if val == 0 || err != "" {
		//临时处理，手慢的人等会再返回
		time.Sleep(200 * time.Millisecond)
		sendRedPacket = redPacket.GetRedisRedPacket()
		ep.RpcCheckRedPacketResult(sendRedPacket)
		return nil
	}
	reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_GET_REDPACKET, redPacket.GetSender(), nil)
	//orderId, _ := for_game.PlaceOrder(who.GetPlayerId(), val, for_game.GOLD_TYPE_GET_REDPACKET)
	orderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, for_game.GOLD_TYPE_GET_REDPACKET)
	msg1 := &share_message.GoldExtendLog{
		RedPacketId: easygo.NewInt64(redPacket.GetId()),
		OrderId:     easygo.NewString(orderId),
		PayType:     easygo.NewInt32(for_game.PAY_TYPE_GOLD),
		Title:       easygo.NewString(reason),
		Gold:        easygo.NewInt64(who.GetGold() + val),
	}
	tarBase := GetPlayerObj(redPacket.GetSender())
	if tarBase != nil {
		msg1.Account = easygo.NewString(tarBase.GetAccount())
	}
	//FinishOrder(orderId, msg1)
	easygo.PCall(NotifyAddGold, who.GetPlayerId(), val, reason, for_game.GOLD_TYPE_GET_REDPACKET, msg1)
	for_game.UpdateRedisRedPacketTotal(who.GetPlayerId(), for_game.GetMillSecond(), val, for_game.REDPACKET_STATISTICS_RECV)
	name := for_game.GetTeamReName(sendRedPacket.GetType(), sendRedPacket.GetTargetId(), who.GetPlayerId(), who.GetNickName())
	msg := &client_hall.OpenRedPacket{
		Id:       easygo.NewInt64(redPacket.GetId()),
		PlayerId: easygo.NewInt64(who.GetPlayerId()),
		NickName: easygo.NewString(name),
		Gold:     easygo.NewInt64(val),
		Type:     easygo.NewInt32(reqMsg.GetType()),
		SenderId: easygo.NewInt64(redPacket.GetSender()),
		SendName: easygo.NewString(redPacket.GetSenderName()),
		Time:     easygo.NewInt64(util.GetMilliTime()),
		State:    easygo.NewInt32(redPacket.GetState()),
	}
	content := for_game.GetRedisOpenRedPacketBase64Str(msg)
	var sessionId string
	t := redPacket.GetType()
	if t == for_game.CHAT_TYPE_PRIVATE {
		sessionId = for_game.MakeSessionKey(redPacket.GetSender(), redPacket.GetTargetId())
	} else if t == for_game.CHAT_TYPE_TEAM {
		sessionId = easygo.AnytoA(redPacket.GetTargetId())
	}
	//领取提示通知
	chatLog := &share_message.Chat{
		SessionId: easygo.NewString(sessionId),
		SourceId:  easygo.NewInt64(who.GetPlayerId()),
		//TargetId:    easygo.NewInt64(redPacket.GetSender()),
		Content:     easygo.NewString(content),
		ChatType:    easygo.NewInt32(redPacket.GetType()),
		ContentType: easygo.NewInt32(for_game.TALK_CONTENT_REDPACKET_LOG),
		PlayIds:     []int64{who.GetPlayerId(), redPacket.GetSender()},
	}
	//群红包带上群号
	if redPacket.GetType() == for_game.RED_PACKET_TEAM_LUCKEY {
		msg.TeamId = easygo.NewInt64(redPacket.GetTargetId())
		chatLog.TargetId = easygo.NewInt64(redPacket.GetTargetId())
	} else {
		chatLog.TargetId = easygo.NewInt64(redPacket.GetSender())
	}
	//通知领玩家的人结果
	// newMsg := redPacket.GetSendRedPacketData()
	ep.RpcCheckRedPacketResult(sendRedPacket)
	self.RpcChatNew(nil, who, chatLog)
	//ep.RpcOpenRedPacketResult(msg)
	//logs.Info("发送通知前端完毕----------")
	////通知发包人有人领了红包
	//if redPacket.GetSender() != who.GetPlayerId() {
	//	if PlayerOnlineMgr.CheckPlayerIsOnLine(redPacket.GetSender()) {
	//		sendEp := ClientEpMp.LoadEndpoint(redPacket.GetSender())
	//		if sendEp != nil {
	//			sendEp.RpcOpenRedPacketResult(msg)
	//		} else {
	//			serverId := PlayerOnlineMgr.GetPlayerServerId(redPacket.GetSender())
	//			playerIds := []int64{redPacket.GetSender()}
	//			BroadCastMsgToOtherHallClient(serverId, playerIds, "RpcOpenRedPacketResult", msg)
	//		}
	//	}
	//}
	return nil
}

//个人转账
func (self *cls1) RpcTransferMoney(ep IGameClientEndpoint, who *Player, reqMsg *share_message.TransferMoney, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcTransferMoney：", reqMsg)
	pid := reqMsg.GetTargetId()
	other := GetPlayerObj(pid)
	if other == nil || other.GetStatus() == for_game.ACCOUNT_CANCELED {
		return easygo.NewFailMsg("该账号异常", for_game.FAIL_MSG_CODE_1016)
	}

	sysParams := PSysParameterMgr.GetSysParameter(for_game.LIMIT_PARAMETER)
	if !sysParams.GetIsTransfer() {
		s := "转账功能已关闭，暂时无法使用该功能"
		logs.Error(s)
		return easygo.NewFailMsg(s)
	}
	//gold := reqMsg.GetGold()
	t := reqMsg.GetWay()
	//检测备注是否包含敏感词
	mark := reqMsg.GetContent()
	if len(mark) > 0 {
		evilType, _ := for_game.PDirtyWordsMgr.CheckWord(mark)
		if evilType {
			return easygo.NewFailMsg("备注包含敏感字，请重新编辑")
		}
	}
	//如果对方不是自己的好友
	if !util.Int64InSlice(pid, who.GetFriends()) {
		res := "对方不是你的好友，操作失败"
		return easygo.NewFailMsg(res)
	}

	//如果对方在自己的黑名单
	if util.Int64InSlice(pid, who.GetBlackList()) {
		res := "对方已被加入黑名单，操作失败"
		return easygo.NewFailMsg(res)
	}

	//如果自己不是对方好友
	if !util.Int64InSlice(who.GetPlayerId(), other.GetFriends()) {
		res := "你不是对方好友，操作失败"
		return easygo.NewFailMsg(res)
	}

	//如果自己在对方的黑名单中
	if util.Int64InSlice(who.GetPlayerId(), other.GetBlackList()) {
		res := "你已被对方拉黑，操作失败"
		return easygo.NewFailMsg(res)
	}

	if t == for_game.PAY_TYPE_GOLD {
		if !who.GetIsPayPassword() {
			return easygo.NewFailMsg("请先设置支付密码")
		}
		if !who.CheckPayPassWord(for_game.Md5(reqMsg.GetPayPassWord())) {
			return easygo.NewFailMsg("支付密码错误")
		}
	}

	//零钱支付才需验证密码
	if t == for_game.PAY_TYPE_GOLD {
		//检测玩家身上的钱够不够
		if who.GetGold() < reqMsg.GetGold() {
			logs.Info("余额不足")
			res := fmt.Sprintf("余额不足%d", reqMsg.GetGold()/100)
			return easygo.NewFailMsg(res)
		}
	} else {
		ep.RpcTunedUpPayInfo(reqMsg.GetPayOrderInfo()) //通知客户端调起支付
		return nil
	}
	chatLog := CreateTransferMoney(reqMsg)
	self.RpcChatNew(nil, who, chatLog)
	return nil
}

func CreateTransferMoney(reqMsg *share_message.TransferMoney, orId ...string) *share_message.Chat {
	receiver := GetPlayerObj(reqMsg.GetTargetId())
	if receiver == nil {
		panic("收款玩家为空")
	}

	//扣除发送金额
	//orderId, _ := for_game.PlaceOrder(reqMsg.GetSender(), -reqMsg.GetGold(), for_game.GOLD_TYPE_SEND_TRANSFER_MONEY)
	orderId := append(orId, "")[0]
	if orderId == "" {
		//如果没有订单，则生成一个
		orderId = for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, for_game.GOLD_TYPE_SEND_TRANSFER_MONEY)
	}
	reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_SEND_TRANSFER_MONEY, reqMsg.GetTargetId(), nil)
	msg := &share_message.GoldExtendLog{
		OrderId:      easygo.NewString(orderId),
		PayType:      easygo.NewInt32(reqMsg.GetWay()),
		TransferText: easygo.NewString(reqMsg.GetContent()),
		HeadIcon:     easygo.NewString(receiver.GetHeadIcon()),
		Title:        easygo.NewString(reason),
		Account:      easygo.NewString(receiver.GetAccount()),
	}
	if reqMsg.GetBankInfo() != "" {
		msg.BankName = easygo.NewString(reqMsg.GetBankInfo())
	}
	base := GetPlayerObj(reqMsg.GetSender())
	//if reqMsg.GetWay() == for_game.PAY_TYPE_GOLD {
	msg.Gold = easygo.NewInt64(base.GetGold() - reqMsg.GetGold())
	//}
	//FinishOrder(orderId, msg)
	NotifyAddGold(base.GetPlayerId(), -reqMsg.GetGold(), reason, for_game.GOLD_TYPE_SEND_TRANSFER_MONEY, msg)
	//转账记录
	transferMoney := for_game.CreateTransferMoney(reqMsg, orderId)
	content := transferMoney.GetTransferMoneyBase64Str()
	chatLog := &share_message.Chat{
		SessionId:   easygo.NewString(for_game.MakeSessionKey(reqMsg.GetSender(), reqMsg.GetTargetId())),
		SourceId:    easygo.NewInt64(reqMsg.GetSender()),
		TargetId:    easygo.NewInt64(reqMsg.GetTargetId()),
		Content:     easygo.NewString(content),
		ChatType:    easygo.NewInt32(for_game.CHAT_TYPE_PRIVATE),
		ContentType: easygo.NewInt32(for_game.TALK_CONTENT_TRANSFER_MONEY),
	}

	easygo.Spawn(func() { for_game.MakePlayerBehaviorReport(6, 0, nil, nil, nil, reqMsg) }) //生成用户行为报表转账相关 已优化到Redis

	return chatLog
}

//领取个人转账
func (self *cls1) RpcOpenTransferMoney(ep IGameClientEndpoint, who *Player, reqMsg *share_message.TransferMoney, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcOpenTransferMoney：", reqMsg)
	transferMoney := for_game.GetRedisTransferMoneyObj(reqMsg.GetId())
	if transferMoney == nil {
		res := fmt.Sprintf("无效的转账记录id=%d", reqMsg.GetId())
		return easygo.NewFailMsg(res)
	}
	if reqMsg.GetOpenWay() == 3 {
		//查看转账信息
		// msg := transferMoney.GetTransferMoneyData()
		ep.RpcOpenTransferMoneyResult(transferMoney.GetTransferMoney())
		return nil
	}
	if who.GetPlayerId() != transferMoney.GetTargetId() {
		res := fmt.Sprintf("领取失败，不是转账对象的玩家id=%d", who.GetPlayerId())
		return easygo.NewFailMsg(res)
	}
	//打开红包、考虑跨服情况
	target := GetPlayerObj(transferMoney.GetTargetId())
	if target == nil {
		res := fmt.Sprintf("领取失败，不是转账对象的玩家id=%d", who.GetPlayerId())
		return easygo.NewFailMsg(res)
	}
	senderId := transferMoney.GetSender()
	sender := for_game.GetRedisPlayerBase(senderId)
	if reqMsg.GetOpenWay() == 1 {
		orderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, for_game.GOLD_TYPE_GET_TRANSFER_MONEY)
		reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_GET_TRANSFER_MONEY, transferMoney.GetSender(), nil)
		msg := &share_message.GoldExtendLog{
			OrderId:      easygo.NewString(orderId),
			TransferText: easygo.NewString(transferMoney.GetContent()),
			HeadIcon:     easygo.NewString(sender.GetHeadIcon()),
			PayType:      easygo.NewInt32(transferMoney.GetWay()),
			Title:        easygo.NewString(reason),
			Gold:         easygo.NewInt64(who.GetGold() + transferMoney.GetGold()),
			Account:      easygo.NewString(sender.GetAccount()),
		}
		transferMoney.SetState(for_game.TRANSFER_MONEY_FNISH) //已被领取
		NotifyAddGold(target.GetPlayerId(), transferMoney.GetGold(), reason, for_game.GOLD_TYPE_GET_TRANSFER_MONEY, msg)

	} else {
		orderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, for_game.GOLD_TYPE_TRANSFER_MONEY_OVER)
		reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_TRANSFER_MONEY_OVER, transferMoney.GetTargetId(), nil)
		msg := &share_message.GoldExtendLog{
			OrderId:      easygo.NewString(orderId),
			TransferText: easygo.NewString(transferMoney.GetContent()),
			HeadIcon:     easygo.NewString(target.GetHeadIcon()),
			PayType:      easygo.NewInt32(transferMoney.GetWay()),
			Title:        easygo.NewString(reason),
			Gold:         easygo.NewInt64(sender.GetGold() + transferMoney.GetGold()),
			Account:      easygo.NewString(sender.GetAccount()),
		}
		title := "转账退款到账通知"
		text := fmt.Sprintf("收到一笔转账退款，退款金额￥%.2f。通过“我的”→“零钱”→“账单”可查看详情", float64(transferMoney.GetGold())/100)
		NoticeAssistant(transferMoney.GetSender(), 1, title, text)
		easygo.PCall(NotifyAddGold, senderId, transferMoney.GetGold(), reason, for_game.GOLD_TYPE_TRANSFER_MONEY_OVER, msg)
		transferMoney.SetState(for_game.TRANSFER_MONEY_BACK) //主动退还
	}
	// transferMoney.SetOpenTime()
	transferMoney.SetOpenTime(time.Now().Unix())
	transferMoney.SaveToMongo()
	sendData := transferMoney.GetTransferMoney()
	//通知自己状态改变
	ep.RpcOpenTransferMoneyResult(sendData)
	//通知发送者状态改变
	SendMsgToHallClientNew(transferMoney.GetSender(), "RpcOpenTransferMoneyResult", sendData)
	msg1 := &client_hall.OpenTransfer{
		Id:       easygo.NewInt64(transferMoney.GetId()),
		SendName: easygo.NewString(sender.GetNickName()),
		NickName: easygo.NewString(target.GetNickName()),
		State:    easygo.NewInt32(transferMoney.GetState()),
		SendId:   easygo.NewInt64(transferMoney.GetSender()),
	}

	content := for_game.GetOpenTransferBase64Str(msg1)
	chatLog := &share_message.Chat{
		SessionId:   easygo.NewString(for_game.MakeSessionKey(who.GetPlayerId(), transferMoney.GetSender())),
		SourceId:    easygo.NewInt64(who.GetPlayerId()),
		TargetId:    easygo.NewInt64(transferMoney.GetSender()),
		Content:     easygo.NewString(content),
		ChatType:    easygo.NewInt32(for_game.CHAT_TYPE_PRIVATE),
		ContentType: easygo.NewInt32(for_game.TALK_CONTENT_TRAMSFER_LOG),
	}
	self.RpcChatNew(nil, who, chatLog)
	return nil
}

//修改群个人信息
func (self *cls1) RpcChangeTeamPersonalSetting(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CgTeamPerSetting, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcChangeTeamPersonalSetting", reqMsg)
	teamId := reqMsg.GetTeamId()

	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}
	playerId := reqMsg.GetOptPlayerId()
	//群主或者管理员操作,修改别人信息
	if playerId > 0 {
		if for_game.GetTeamPlayerPos(teamId, who.GetPlayerId()) >= for_game.TEAM_MASSES {
			return easygo.NewFailMsg("你没有权限执行该操作")
		}
	} else {
		playerId = who.GetPlayerId()
	}
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}
	t := reqMsg.GetType()
	memberObj := for_game.GetRedisTeamPersonalObj(teamId)
	if t == 4 { //修改群个人昵称
		name := reqMsg.GetName() //通知群中成员 某个人修改群昵称
		//检测备注是否包含敏感词
		if len(name) > 0 {
			evilType, _ := for_game.PDirtyWordsMgr.CheckWord(name)
			if evilType {
				return easygo.NewFailMsg("昵称包含敏感字，请重新编辑")
			}
		}

		memberObj.UpdateTeamPersonalSetting(playerId, t, reqMsg.GetName())
		/*		if name == "" {
				base := GetPlayerObj(playerId)
				name = base.GetNickName()
				reqMsg.Name = easygo.NewString(name)
			}*/

		msg := &client_hall.ChangeNameInfo{
			TeamId:   easygo.NewInt64(teamId),
			PlayerId: easygo.NewInt64(playerId),
			Name:     easygo.NewString(name),
		}
		serverId := PlayerOnlineMgr.GetPlayerServerId(who.GetPlayerId())
		TeamSendMessage(teamObj.GetTeamMemberList(), 0, serverId, "RpcRefreshTeamPersonalName", msg)
	} else {
		memberObj.UpdateTeamPersonalSetting(playerId, t, reqMsg.GetValue())
	}
	return reqMsg
}

//修改群管理信息
func (self *cls1) RpcChangeTeamManageSetting(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CgTeamManageSetting, common ...*base.Common) easygo.IMessage {
	logs.Info("======修改群管理信息 RpcChangeTeamManageSetting=======reqMsg=%v", reqMsg)
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}
	playerId := who.GetPlayerId()
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}
	if for_game.GetTeamPlayerPos(teamId, playerId) > for_game.TEAM_MANAGER {
		return easygo.NewFailMsg("你没有这个权限")
	}

	owner := teamObj.GetTeamOwner()
	info := &share_message.OperatorInfo{
		Operator: easygo.NewString(who.GetNickName()),
		Time:     easygo.NewInt64(util.GetMilliTime()),
		Flag:     easygo.NewInt32(2),
	}
	if owner != who.GetPlayerId() {
		info.Flag = easygo.NewInt32(3)
	}

	var v interface{}
	t := reqMsg.GetType()
	switch t {
	case for_game.TIME_CLEAR:
		v = reqMsg.GetValue()
		teamObj.SetIsTimeClean(reqMsg.GetValue())
	case for_game.READ_CLEAR:
		v = reqMsg.GetValue()
		teamObj.SetTeamIsReadClean(reqMsg.GetValue())
	case for_game.SCREENSHOT_NOTICE:
		v = reqMsg.GetValue()
		teamObj.SetTeamIsScreenShotNotify(reqMsg.GetValue())
	case for_game.STOP_TALK:
		v = reqMsg.GetValue()
		teamObj.SetTeamIsStopTalk(reqMsg.GetValue())
		if v == true { //封禁
			teamObj.SetTeamStatus(for_game.BANNED)
			teamObj.SetRedisOperatorInfo(info)
			teamObj.SaveToMongo()
		} else { //解封
			if owner != who.GetPlayerId() {
				teamObj.DelRedisOperatorInfo(for_game.MANAGER)
			} else {
				teamObj.DelRedisOperatorInfo(for_game.OWNER)
			}

			message := for_game.GetTeamManageSetting(teamId).GetMessageSetting()
			if message.GetIsBan() == false { //设置状态前判断后台是否封禁，如果没有封禁
				teamObj.SetTeamStatus(for_game.NORMAL)
			}

			teamObj.SaveToMongo()
		}

	case for_game.STOP_ADDFRIEND:
		pos := for_game.GetTeamPlayerPos(teamId, playerId)
		if pos != for_game.TEAM_OWNER {
			return easygo.NewFailMsg("只有群主能设置此项功能")
		}
		v = reqMsg.GetValue()
		teamObj.SetTeamIsAddFriend(reqMsg.GetValue())
	case for_game.TEAM_INVITE:
		v = reqMsg.GetValue()
		teamObj.SetTeamIsInvite(reqMsg.GetValue())
	case for_game.TEAM_NAME:
		v = reqMsg.GetValue2()
		s := easygo.AnytoA(v)
		evilType, _ := for_game.PDirtyWordsMgr.CheckWord(s)
		if evilType {
			return easygo.NewFailMsg("群聊名称包含敏感字，请重新编辑")
		}
		teamObj.SetTeamNikeName(reqMsg.GetValue2())
	case for_game.TEAM_GONGGAO:
		v = reqMsg.GetValue2()
		s := easygo.AnytoA(v)
		evilType, _ := for_game.PDirtyWordsMgr.CheckWord(s)
		if evilType {
			return easygo.NewFailMsg("群公告包含敏感字，请重新编辑")
		}
		teamObj.SetTeamGongGao(reqMsg.GetValue2())
	case for_game.TEAM_RECOMMEND:
		v = reqMsg.GetValue()
		teamObj.SetTeamIsRecommend(reqMsg.GetValue())
	case for_game.CHANGE_TEAMMONEYCODE:
		//检测备注是否包含敏感词
		js, _ := base64.StdEncoding.DecodeString(reqMsg.GetValue2())
		data := easygo.KWAT{}
		_ = json.Unmarshal(js, &data)
		remark := data.GetString("remark")
		if len(remark) > 0 {
			evilType, _ := for_game.PDirtyWordsMgr.CheckWord(remark)
			if evilType {
				return easygo.NewFailMsg("备注包含敏感字，请重新编辑")
			}
		}
		pos := for_game.GetTeamPlayerPos(teamId, playerId)
		if pos != for_game.TEAM_OWNER {
			return easygo.NewFailMsg("只有群主能设置此项功能")
		}
		v = reqMsg.GetValue2()
		teamObj.SetTeamQRCode(reqMsg.GetValue2())
	case for_game.STOP_ADDTEAM:
		pos := for_game.GetTeamPlayerPos(teamId, playerId)
		if pos > for_game.TEAM_MANAGER {
			return easygo.NewFailMsg("只有群主和管理员能设置此项功能")
		}
		v = reqMsg.GetValue()
		teamObj.SetTeamIsStopAddTeam(reqMsg.GetValue())
	case for_game.WELCOME_WORD: //欢迎语开关
		v = reqMsg.GetValue()
		pos := for_game.GetTeamPlayerPos(teamId, playerId)
		setting := teamObj.GetTeamMessageSetting()
		if setting.GetIsManagerEdit() {
			//管理员有权限
			if pos > for_game.TEAM_MANAGER {
				return easygo.NewFailMsg("只有群主和管理员能设置此项功能")
			}
		} else {
			//管理员无权限
			if pos >= for_game.TEAM_MANAGER {
				return easygo.NewFailMsg("只有群主能设置此项功能")
			}
		}
		setting.IsOpenWelcomeWord = easygo.NewBool(reqMsg.GetValue())
		teamObj.UpdateTeamMessageSetting(setting)
		return reqMsg

	case for_game.WELCOME_WORD_MANAGER: //管理员是否有权限编辑欢迎语
		v = reqMsg.GetValue()
		pos := for_game.GetTeamPlayerPos(teamId, playerId)
		setting := teamObj.GetTeamMessageSetting()
		if pos >= for_game.TEAM_MANAGER {
			return easygo.NewFailMsg("只有群主能设置此项功能")
		}
		setting.IsManagerEdit = easygo.NewBool(reqMsg.GetValue())
		teamObj.UpdateTeamMessageSetting(setting)

		// 异步通知群管理员修改了群设置,重新拉去新数据
		easygo.Spawn(notifyManage, playerId, teamObj, setting)

		return reqMsg
	case for_game.EDIT_WELCOME_WORD: //编辑新的欢迎语
		v = reqMsg.GetValue2()
		pos := for_game.GetTeamPlayerPos(teamId, playerId)
		setting := teamObj.GetTeamMessageSetting()
		if setting.GetIsManagerEdit() {
			//管理员有权限
			if pos > for_game.TEAM_MANAGER {
				return easygo.NewFailMsg("只有群主和管理员能设置此项功能")
			}
		} else {
			//管理员无权限
			if pos >= for_game.TEAM_MANAGER {
				return easygo.NewFailMsg("只有群主能设置此项功能")
			}
		}
		// 检测群欢迎语是否有违规字符
		if reqMsg.GetValue2() != "" {
			evilType, _ := for_game.PDirtyWordsMgr.CheckWord(reqMsg.GetValue2())
			logs.Info(evilType)
			if evilType {
				return easygo.NewFailMsg("签名包含敏感字，请重新编辑")
			}
		}

		teamObj.SetWelcomeWord(reqMsg.GetValue2())
		// 异步通知群管理员修改了群设置,重新拉去新数据
		easygo.Spawn(notifyManage, playerId, teamObj, setting)
		return reqMsg
	case for_game.TEAM_HEAD:
		teamHead := reqMsg.GetValue2()
		if who.GetPlayerId() != teamObj.GetTeamOwner() { // 只有群主才可以设置
			return easygo.NewFailMsg("只有群主能设置此项功能")
		}
		teamObj.SetHeadUrl(teamHead)
		////对应会话修改
		//sessionObj := for_game.GetRedisChatSessionObj(easygo.AnytoA(teamId))
		//if sessionObj != nil {
		//	sessionObj.SetSessionName(teamHead)
		//}
		return reqMsg
	case for_game.TOPIC_TEAM_DESC:
		//话题简介修改
		desc := reqMsg.GetValue2()
		pos := for_game.GetTeamPlayerPos(teamId, playerId)
		if pos > for_game.TEAM_MANAGER {
			return easygo.NewFailMsg("只有群主和管理员能设置此项功能")
		}
		teamObj.SetTopicDesc(desc)
	}

	if t == for_game.CHANGE_TEAMMONEYCODE {
		if teamObj.GetTeamMessageSetting().GetIsOpenTeamMoneyCode() == false { //第一次开启的时候才通知
			teamObj.SetTeamIsOpenTeamMoneyCode(true)
		}
	}
	NoticeTeamMessage(teamId, playerId, t, v)

	return reqMsg
}

// 通知群管理员修改了群设置,重新拉去新数据
func notifyManage(playerId int64, teamObj *for_game.RedisTeamObj, setting *share_message.MessageSetting) {
	logs.Info("notifyManage------->", "群设置通知管理员开始")
	notifyIds := make([]int64, 0)
	if playerId == teamObj.GetTeamOwner() {
		notifyIds = append(notifyIds, teamObj.GetTeamManageList()...)
	} else {
		notifyIds = append(notifyIds, teamObj.GetTeamOwner()) // 添加群主
		// 添加其他群管理
		for _, manageId := range teamObj.GetTeamManageList() {
			if playerId == manageId {
				continue
			}
			notifyIds = append(notifyIds, manageId) // 添加群管理.不包含自己.
		}
	}
	for _, manageId := range notifyIds {
		if !PlayerOnlineMgr.CheckPlayerIsOnLine(manageId) {
			continue
		}
		pep := ClientEpMp.LoadEndpoint(manageId)
		if pep == nil {
			continue
		}
		notify := &client_hall.TeamSettingNotify{
			TeamId:         easygo.NewInt64(teamObj.GetId()),
			WelcomeWord:    easygo.NewString(teamObj.GetWelcomeWord()),
			MessageSetting: setting,
		}
		pep.RpcTeamSettingNotify(notify)
	}
}

//修改好友设置信息
func (self *cls1) RpcChangeFriendSetting(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.FriendSetting, common ...*base.Common) easygo.IMessage {
	playerId := reqMsg.GetPlayerId()
	t := reqMsg.GetType()
	value := reqMsg.GetValue()
	wq := for_game.GetFriendBase(who.GetPlayerId())
	if !util.Int64InSlice(playerId, wq.GetFriendIds()) {
		res := "该玩家还不是你的好友"
		return easygo.NewFailMsg(res)
	}

	wq.SetFriendSetting(playerId, t, value)
	return reqMsg
}

//添加管理员
func (self *cls1) RpcAddTeamManager(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamReq, common ...*base.Common) easygo.IMessage {
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}
	playerId := who.GetPlayerId()
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}

	pos := for_game.GetTeamPlayerPos(teamId, playerId)
	if pos != for_game.TEAM_OWNER {
		res := "你不是群主不能进行这个操作"
		return easygo.NewFailMsg(res)
	}

	playerlst := reqMsg.GetPlayerIdList()
	if !teamObj.CheckTeamMaxManager(len(playerlst)) {
		res := "添加管理员人数已达上限"
		return easygo.NewFailMsg(res)

	}

	SetTeamMemberPosition(teamId, playerlst, for_game.TEAM_MANAGER)
	NoticeTeamMessage(teamId, playerId, for_game.ADD_MANAGER, playerlst)
	return nil
}

//删除管理员
func (self *cls1) RpcDelTeamManager(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamReq, common ...*base.Common) easygo.IMessage {
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}
	playerId := who.GetPlayerId()
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}

	pos := for_game.GetTeamPlayerPos(teamId, playerId)
	if pos != for_game.TEAM_OWNER {
		res := "你不是群主不能进行这个操作"
		return easygo.NewFailMsg(res)
	}

	playerlst := reqMsg.GetPlayerIdList()
	SetTeamMemberPosition(teamId, playerlst, for_game.TEAM_MASSES)
	NoticeTeamMessage(teamId, playerId, for_game.DEL_MANAGER, playerlst)
	return nil
}

//转让群主
func (self *cls1) RpcChangeTeamOwner(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamReq, common ...*base.Common) easygo.IMessage {
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}
	playerId := who.GetPlayerId()
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}

	pos := for_game.GetTeamPlayerPos(teamId, playerId)
	if pos != for_game.TEAM_OWNER {
		res := "你不是群主不能进行这个操作"
		return easygo.NewFailMsg(res)
	}
	playerlst := reqMsg.GetPlayerIdList()
	SetTeamMemberPosition(teamId, playerlst, for_game.TEAM_OWNER)
	NoticeTeamMessage(teamId, playerId, for_game.CHANGE_OWNER, playerlst)
	return nil
}

//通过进群请求
func (self *cls1) RpcAcceptAddTeam(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamInfo, common ...*base.Common) easygo.IMessage {
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}

	agreeId := reqMsg.GetPlayerId()
	playerId := who.GetPlayerId()
	logId := reqMsg.GetLogId()

	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}

	pos := for_game.GetTeamPlayerPos(teamId, playerId)
	if pos > for_game.TEAM_MANAGER {
		return easygo.NewFailMsg("你没有权限做这个操作")
	}

	if !teamObj.CheckTeamMaxMember(1) {
		//return easygo.NewFailMsg("添加成员数量超过该群聊数量上限")
		return easygo.NewFailMsg("群成员超出上限")
	}

	if teamObj.GetTeamInviteStateForLogId(logId) != for_game.INVITE_UNTREATED {
		res := "该请求已经被处理"
		msg := &client_hall.AllInviteInfo{
			AllInfo: teamObj.GetTeamInviteInfo(),
		}
		ep.RpcRefreshTeamInviteInfo(msg)
		return easygo.NewFailMsg(res)
	}
	if util.Int64InSlice(agreeId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("已经在群里面了，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}
	reason, inviteId, channel := teamObj.AcceptTeamRequest(logId)
	plst := []int64{agreeId}
	serverId := PlayerOnlineMgr.GetPlayerServerId(playerId)
	success := AddTeamMember(teamId, teamObj.GetSessionLogMaxId(), plst, false, 0, reason, serverId, channel)
	if !success {
		return easygo.NewFailMsg("添加群员失败")
	}
	NoticeTeamMessage(teamId, inviteId, for_game.ACTIVE_ADDTEAM, plst)

	target := GetPlayerObj(agreeId)
	if target != nil {
		target.AddTeamId(teamId)
	}

	msg := &client_hall.DealInviteInfo{
		LogId: easygo.NewInt64(logId),
		State: easygo.NewInt32(for_game.INVITE_ACCEPT),
	}
	ep.RpcDealAddTeamRequest(msg)
	return nil
}

//拒绝进群请求
func (self *cls1) RpcRefuseAddTeam(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamInfo, common ...*base.Common) easygo.IMessage {
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		panic("群对象怎么会是个空")
	}
	playerId := who.GetPlayerId()
	logId := reqMsg.GetLogId()
	if !util.Int64InSlice(playerId, teamObj.GetTeamMemberList()) {
		res := fmt.Sprintf("你不在这个群里面，群id:%d", teamId)
		return easygo.NewFailMsg(res)
	}

	pos := for_game.GetTeamPlayerPos(teamId, playerId)
	if pos > for_game.TEAM_MANAGER {
		return easygo.NewFailMsg("你没有权限做这个操作")
	}

	if teamObj.GetTeamInviteStateForLogId(logId) != for_game.INVITE_UNTREATED {
		res := "该请求已经被处理"
		msg := &client_hall.AllInviteInfo{
			AllInfo: teamObj.GetTeamInviteInfo(),
		}
		ep.RpcRefreshTeamInviteInfo(msg)
		return easygo.NewFailMsg(res)
	}

	teamObj.UpdateTeamInviteState(logId, for_game.INVITE_REFUSE)
	msg := &client_hall.DealInviteInfo{
		LogId: easygo.NewInt64(logId),
		State: easygo.NewInt32(for_game.INVITE_REFUSE),
	}
	ep.RpcDealAddTeamRequest(msg)
	return nil
}

//转发客户端连接报道
func (self *cls1) RpcTFToServer(ep IGameClientEndpoint, ctx interface{}, reqMsg *client_server.ClientInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcTFToServer:", reqMsg)
	return nil
}

//获取玩家在群中的最新信息
func (self *cls1) RpcGetTeamPlayerInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("==========获取玩家在群中的最新信息 RpcGetTeamPlayerInfo===============", reqMsg)
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		return easygo.NewFailMsg("该群不存在")
	}
	pid := reqMsg.GetPlayerId()   // 玩家id
	playerId := who.GetPlayerId() // 操作者id
	//if for_game.GetTeamPlayerPos(teamId, pid) == 0 { //这个人不在群中
	//	logs.Info(fmt.Sprintf("用户不在该群中", pid))
	//	return &share_message.TeamPlayerInfo{}
	//}
	currentPage := reqMsg.GetCurrentPage()
	pageSize := reqMsg.GetPageSize()
	if currentPage == 0 {
		currentPage = 1
	}
	if pageSize == 0 {
		pageSize = for_game.DefaultPageSize
	}
	msg := GetTeamPlayerInfo(teamId, pid, playerId, currentPage, pageSize)
	if msg == nil {
		logs.Error("玩家数据不存在:", pid)
		return easygo.NewFailMsg("玩家不存在")
	}
	msg.Type = easygo.NewInt32(reqMsg.GetType())
	// 获取我的关注列表
	if util.Int64InSlice(pid, who.GetAttention()) {
		msg.IsOnMyAttentionList = easygo.NewBool(true)
	}
	if util.Int64InSlice(pid, who.GetBlackList()) { // 是否在黑名单中.
		msg.IsOnMyBlackList = easygo.NewBool(true)
	}
	msg.LabelInfo = for_game.GetLabelInfo(for_game.GetRedisPlayerBase(pid).GetLabelList())
	return msg
}

//获取群的基本信息
func (self *cls1) RpcGetBaseTeamInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TeamReq, common ...*base.Common) easygo.IMessage {
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	team := teamObj.GetRedisTeam()
	inviteId := reqMsg.GetInviteId()
	base := GetPlayerObj(inviteId)
	msg := &client_hall.TeamDataInfo{
		Team:     team,
		Members:  for_game.GetAllTeamMember(teamId),
		PlayerId: easygo.NewInt64(inviteId),
		Type:     easygo.NewInt32(reqMsg.GetType()),
		Name:     easygo.NewString(base.GetNickName()),
		HeadUrl:  easygo.NewString(base.GetHeadIcon()),
		Sex:      easygo.NewInt32(base.GetSex()),
	}
	return msg
}

// 获取个人名片信息
func (self *cls1) RpcGetPlayerCardInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_server.PlayerReq, common ...*base.Common) easygo.IMessage {
	logs.Info("=======================大厅获取个人名片信息 RpcGetPlayerCardInfo=======================", reqMsg)
	msg := &share_message.TeamPlayerInfo{}
	pid := reqMsg.GetPlayerId()
	opId := who.GetPlayerId()                                                   // 操作者id
	if for_game.Min_Robot_PlayerId < pid && pid < for_game.Max_Robot_PlayerId { //机器人
		msg.PlayerId = easygo.NewInt64(0)
		return msg
	}

	var state int32
	if util.Int64InSlice(pid, who.GetFriends()) {
		state = 1
	} else {
		state = 2
	}

	nickName := for_game.GetFriendsReName(who.GetPlayerId(), pid)
	msg = GetCardInfo1(opId, pid, state, nickName)
	// 获取我的关注列表
	if util.Int64InSlice(pid, who.GetAttention()) {
		msg.IsOnMyAttentionList = easygo.NewBool(true)
	}
	if util.Int64InSlice(pid, who.GetBlackList()) { // 是否在黑名单中.
		msg.IsOnMyBlackList = easygo.NewBool(true)
	}
	msg.Type = easygo.NewInt32(reqMsg.GetType())
	player := for_game.GetRedisPlayerBase(pid)
	if player != nil {
		msg.LabelInfo = for_game.GetLabelInfo(player.GetLabelList())
	}
	if len(common) > 0 && common[0].GetVersion() == "2.7.4" {
		return msg
	}

	currentPage := reqMsg.GetCurrentPage()
	pageSize := reqMsg.GetPageSize()
	if currentPage == 0 { // 前端不传,默认第一页
		currentPage = 1
	}
	if pageSize == 0 { // 前端不传,默认50条
		pageSize = for_game.DefaultPageSize
	}

	base := for_game.GetRedisPlayerBase(pid)
	if base == nil {
		logs.Error("无效的玩家id", pid)
		return msg
	}
	ds, count := base.GetRedisPlayerDynamicListByPage(opId, currentPage, pageSize) // 分页
	dynamicList := for_game.GetRedisDynamicForSomeLogId1(opId, ds, pid)
	dynamicData := &share_message.DynamicDataListPage{
		DynamicData: dynamicList,
		TotalCount:  easygo.NewInt32(count),
	}
	msg.DynamicData = dynamicData

	//  异步处理动态的话题浏览量
	easygo.Spawn(func() {
		for_game.OperateViewNum(dynamicList)
	})
	return msg
}

func (self *cls1) RpcCheckout(ep IGameClientEndpoint, who *Player, reqMsg *share_message.OrderID, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcCheckout ", reqMsg)

	//定义锁数组
	var itemLockKeys []string = []string{}
	var orderLockKeys []string = []string{}
	var orderLockPayIngKeys []string = []string{}

	//取得订单的订单id和订单对应的商品id
	var bill *share_message.TableBill = &share_message.TableBill{}

	colBill, closeFunBill := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_BILLS)
	defer closeFunBill()

	colOrder, closeFunOrder := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFunOrder()

	eBill := colBill.Find(bson.M{"_id": reqMsg.GetOrderId()}).One(bill)
	if eBill != nil && eBill != mgo.ErrNotFound {
		logs.Error(eBill)
		return easygo.NewFailMsg("操作失败,刷新重试！")
	}

	//说明不是从购物车结算过来的支付操作
	if eBill == mgo.ErrNotFound {

		var shopOrder share_message.TableShopOrder = share_message.TableShopOrder{}
		//需要从订单表中取得数据

		eOrder := colOrder.Find(bson.M{"_id": reqMsg.GetOrderId()}).One(&shopOrder)
		if eOrder != nil && eOrder != mgo.ErrNotFound {
			logs.Error(eOrder)
			return easygo.NewFailMsg("操作失败,刷新重试！")
		}

		if eOrder == mgo.ErrNotFound {
			logs.Error(eOrder)
			return easygo.NewFailMsg("订单不存在")
		}

		if shopOrder.GetItems() != nil {

			if shopOrder.GetItems().GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY {
				if shopOrder.GetState() == for_game.SHOP_ORDER_EVALUTE {
					return easygo.NewFailMsg("重复支付,刷新重试！")
				}
			} else {
				if shopOrder.GetState() == for_game.SHOP_ORDER_WAIT_SEND {
					return easygo.NewFailMsg("重复支付,刷新重试！")
				}
			}
		} else {
			return easygo.NewFailMsg("订单无商品信息！")
		}

		if shopOrder.GetState() != for_game.SHOP_ORDER_WAIT_PAY {
			return easygo.NewFailMsg("操作失败,刷新重试！")
		}

		tempItemLockKey := for_game.MakeRedisKey(for_game.SHOP_ITEM_PAY_MUTEX, shopOrder.GetItems().GetItemId())
		itemLockKeys = []string{tempItemLockKey}

		tempOrderLockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_WAIT_PAY_MUTEX, shopOrder.GetOrderId())
		orderLockKeys = []string{tempOrderLockKey}

		tempOrderLockPayIngKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_PAYING_MUTEX, shopOrder.GetOrderId())
		orderLockPayIngKeys = []string{tempOrderLockPayIngKey}
		//说明是从购物车结算过来的支付操作
	} else {

		if bill.GetState() == for_game.SHOP_ORDER_WAIT_SEND {
			return easygo.NewFailMsg("重复支付,刷新重试！")
		}
		if bill.GetState() != for_game.SHOP_ORDER_WAIT_PAY {
			return easygo.NewFailMsg("操作失败,刷新重试！")
		}

		//取得各个子订单对应的商品的id
		var shopOrderList []*share_message.TableShopOrder
		//需要从订单表中取得数据
		eOrder := colOrder.Find(bson.M{"_id": bson.M{"$in": bill.GetOrderList()}}).All(&shopOrderList)
		if eOrder != nil && eOrder != mgo.ErrNotFound {
			logs.Error(eOrder)
			return easygo.NewFailMsg("操作失败,刷新重试！")
		}

		if eOrder == mgo.ErrNotFound {
			logs.Error(eOrder)
			return easygo.NewFailMsg("订单不存在")
		}

		for _, value := range shopOrderList {

			if value.GetItems() != nil {

				if value.GetItems().GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY {
					if value.GetState() == for_game.SHOP_ORDER_EVALUTE {
						return easygo.NewFailMsg("存在订单重复支付,刷新重试！")
					}
				} else {
					if value.GetState() == for_game.SHOP_ORDER_WAIT_SEND {
						return easygo.NewFailMsg("存在订单重复支付,刷新重试！")
					}
				}
			} else {
				return easygo.NewFailMsg("存在无商品信息的订单！")
			}

			if value.GetState() != for_game.SHOP_ORDER_WAIT_PAY {
				return easygo.NewFailMsg("操作失败,刷新重试！")
			}

			tempItemLockKey := for_game.MakeRedisKey(for_game.SHOP_ITEM_PAY_MUTEX, value.GetItems().GetItemId())
			itemLockKeys = append(itemLockKeys, tempItemLockKey)

			tempOrderLockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_WAIT_PAY_MUTEX, value.GetOrderId())
			orderLockKeys = append(orderLockKeys, tempOrderLockKey)

			tempOrderLockPayIngKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_PAYING_MUTEX, value.GetOrderId())
			orderLockPayIngKeys = append(orderLockPayIngKeys, tempOrderLockPayIngKey)
		}

	}

	//用户一些异常操作、进入即给与10s的分布式过期锁
	//这里不要主动清理、要让key自动失效
	errLockFirst := easygo.RedisMgr.GetC().DoBatchRedisLockNoRetry(orderLockPayIngKeys, 10)

	//未取得锁就直接不做了
	if errLockFirst != nil {
		s := fmt.Sprintf("RpcCheckout 取得商品多key分布式重试锁失败,redis keys is %v", orderLockPayIngKeys)
		logs.Error(s)
		logs.Error(errLockFirst)
		return easygo.NewFailMsg("下单失败,刷新重试！")
	}
	//取得分布式锁开始1、取得订单对应的商品的分布式锁(阻塞重试）2、取得订单对应的每个订单的分布式锁(不需要重试）
	//1、取得订单对应的商品的分布式锁，此锁需要重试，直到重试次数结束提示退出
	errLock := easygo.RedisMgr.GetC().DoBatchRedisLockWithRetry(itemLockKeys, 10)
	defer easygo.RedisMgr.GetC().DoBatchRedisUnlock(itemLockKeys)

	//如果重试后还未取得锁就直接不做了
	if errLock != nil {
		s := fmt.Sprintf("RpcCheckout 取得商品多key分布式重试锁失败,redis keys is %v", itemLockKeys)
		logs.Error(s)
		logs.Error(errLock)
		return easygo.NewFailMsg("下单失败,刷新重试！")
	}

	//因为订单失效、取消等操作恢复库存这里必须加订单
	//2、取得订单的分布式锁，此锁不重试,有一个取不到就返回
	errLock2 := easygo.RedisMgr.GetC().DoBatchRedisLockNoRetry(orderLockKeys, 10)
	defer easygo.RedisMgr.GetC().DoBatchRedisUnlock(orderLockKeys)

	//未取得锁就直接不做了
	if errLock2 != nil {
		s := fmt.Sprintf("RpcCheckout 取得订单多key分布式不重试锁失败,redis keys is %v", orderLockKeys)
		logs.Error(s)
		logs.Error(errLock)
		return easygo.NewFailMsg("下单失败,刷新重试！")
	}

	//判断库存和减少库存操作
	rst := for_game.ShopCheckAndSubStock(reqMsg.GetOrderId())
	if rst != "" {
		logs.Error(rst)
		return easygo.NewFailMsg(rst)
	}

	//开始扣钱操作
	result := ShopOrderPay(reqMsg.GetOrderId(), who, for_game.PAY_TYPE_GOLD, "")
	if result != "" {
		//恢复库存
		rst1 := for_game.ShopRecoverStock(reqMsg.GetOrderId())
		if rst1 != "" {
			logs.Error(rst1)
			return easygo.NewFailMsg(rst1)
		}
		return easygo.NewFailMsg(result)
	}

	//处理商城订单的后续操作
	return who.SendMsgToShop("RpcPayCallBack", reqMsg)
}

//客户端发起支付充值
func (self *cls1) RpcRechargeMoney(ep IGameClientEndpoint, who *Player, reqMsg *share_message.PayOrderInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcRechargeMoney:", reqMsg)
	resp := &share_message.PayOrderResult{
		Result: easygo.NewBool(false),
	}
	if reqMsg.GetPayId() == for_game.PAY_CHANNEL_YTS_WX {
		b, err := PWebYunTongShangPay.RechargeOrder(WebServreMgr, reqMsg)
		if err != nil {
			logs.Error("支付请求失败:", err)
			return resp
		}
		data := easygo.KWAT{}
		er := json.Unmarshal([]byte(b), &data)
		if er != nil {
			logs.Error("解析数据异常:err ", er.Error())
			return resp
		}
		payInfo := data.GetString("data")
		if payInfo != "" {
			resp.PayInfo = easygo.NewString(payInfo)
			resp.Result = easygo.NewBool(true)
		}
	}
	return resp
}

//请求群未读聊天记录
func (self *cls1) RpcRequestChatInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.ChatInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcRequestChatInfo:", reqMsg)
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		return easygo.NewFailMsg("不存在这个群")
	}
	playerId := who.GetPlayerId()
	if for_game.GetTeamPlayerPos(teamId, playerId) == for_game.TEAM_UNUSE {
		return easygo.NewFailMsg("你不在这个群中")
	}
	chatObj := for_game.GetRedisTeamChatLog(teamId)
	msg := chatObj.GetTeamChatLogListForLogId(reqMsg.GetStartId(), reqMsg.GetOverId(), reqMsg.GetCount())
	ep.RpcReturnChatInfo(msg)
	return nil
}

//银行卡支付入口
func (self *cls1) RpcReqPaySMS(ep IGameClientEndpoint, who *Player, reqMsg *share_message.PayOrderInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("reqMsg:", reqMsg)

	//if !for_game.IS_FORMAL_SERVER {
	//	return easygo.NewFailMsg("测试服不能充值，请联系客服")
	//}

	sysConfig := PSysParameterMgr.GetSysParameter(for_game.LIMIT_PARAMETER)
	if !sysConfig.GetIsRecharge() {
		return easygo.NewFailMsg("充值入口已关闭，请联系客服")
	}

	if reqMsg.GetAmount() == "" {
		return easygo.NewFailMsg("支付金额异常")
	}
	if reqMsg.GetProduceName() == "" {
		return easygo.NewFailMsg("商品描述不能为空")
	}
	bankInfo := who.GetBankMsg(reqMsg.GetPayBankNo())
	if bankInfo == nil {
		return easygo.NewFailMsg("无效的银行卡号")
	}
	//reqMsg.CardNo = easygo.NewString(bankInfo.GetSignNo())
	base := GetPlayerObj(who.GetPlayerId())
	if reqMsg.GetPayId() == for_game.PAY_CHANNEL_HUIJU {
		//汇聚支付
		if reqMsg.GetPayPassType() == 0 {
			//短信支付
			data, err := PWebHuiJuPay.ReqPaySMSApi(reqMsg, base)
			if err != nil {
				return err
			}
			logs.Info("ReqPaySMSApi data:", data)
			msg := &client_hall.BankPaySMS{
				OrderNo: easygo.NewString(data.GetString("mch_order_no")),
				PayId:   easygo.NewInt32(reqMsg.GetPayId()),
			}
			ep.RpcReqPaySMSResult(msg)
		} else {
			//免短信支付
			data, err := PWebHuiJuPay.ReqFastPayApi(reqMsg, base)
			if err != nil {
				return err
			}
			logs.Info("ReqFastPayApi data:", data)
			PWebHuiJuPay.DealBankPayResult(data)
		}
	} else if reqMsg.GetPayId() == for_game.PAY_CHANNEL_HUICHAO_YL {
		//汇潮支付
		data, err := PWebHuiChaoPay.ReqPaySMSApi(reqMsg, base)
		if err != nil {
			return err
		}
		logs.Info("ReqPaySMSApi data:", data)
		msg := &client_hall.BankPaySMS{
			OrderNo: easygo.NewString(data.GetString("mch_order_no")),
			PayId:   easygo.NewInt32(reqMsg.GetPayId()),
		}
		//存储下外部订单
		order := for_game.GetRedisOrderObj(msg.GetOrderNo())
		if order != nil {
			order.SetExternalNo(data.GetString("jp_order_no"))
		}
		ep.RpcReqPaySMSResult(msg)
	} else {
		//免短信支付
		data, err := PWebHuiJuPay.ReqFastPayApi(reqMsg, base)
		if err != nil {
			return err
		}
		logs.Info("ReqFastPayApi data:", data)
		//存储外部订单号
		order := for_game.GetRedisOrderObj(data.GetString("mch_order_no"))
		if order != nil {
			order.SetExternalNo(data.GetString("jp_order_no"))
		}
		PWebHuiJuPay.DealBankPayResult(data)
	}
	return nil
}

func (self *cls1) RpcReqSMSPay(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.BankPaySMS, common ...*base.Common) easygo.IMessage {

	//if !for_game.IS_FORMAL_SERVER {
	//	return easygo.NewFailMsg("测试服不能充值，请联系客服")
	//}
	logs.Info("RpcReqSMSPay--->>:", reqMsg)
	sysConfig := PSysParameterMgr.GetSysParameter(for_game.LIMIT_PARAMETER)
	if !sysConfig.GetIsRecharge() {
		return easygo.NewFailMsg("充值入口已关闭，请联系客服")
	}

	if reqMsg.GetSMS() == "" {
		return easygo.NewFailMsg("输入验证码错误")
	}
	if reqMsg.GetOrderNo() == "" {
		return easygo.NewFailMsg("订单号不能为空")
	}
	if reqMsg.GetPayId() == for_game.PAY_CHANNEL_HUIJU {
		data, err := PWebHuiJuPay.ReqSMSPayApi(reqMsg)
		if err != nil {
			return err
		}
		logs.Info("RpcReqSysPay data:", data)
		PWebHuiJuPay.DealBankPayResult(data, time.Second*5)
	} else if reqMsg.GetPayId() == for_game.PAY_CHANNEL_HUICHAO_YL {
		data, err := PWebHuiChaoPay.ReqSMSPayApi(reqMsg)
		if err != nil {
			return err
		}
		logs.Info("RpcReqSysPay data:", data)
		PWebHuiChaoPay.ReqCheckPayOrder(data.GetString("requestNo"), time.Second*2)
	}

	return nil
}

//检验账号有效
func (self *cls1) RpcCheckAccountVaild(ep IGameClientEndpoint, ctx interface{}, reqMsg *client_server.CheckInfo, common ...*base.Common) easygo.IMessage {
	account := reqMsg.GetAccount()
	accountInfo := for_game.GetRedisAccountByPhone(account)
	if accountInfo == nil {
		reqMsg.Vaild = easygo.NewBool(false)
	} else {
		pid := accountInfo.GetPlayerId()
		base := for_game.GetRedisPlayerBase(pid)
		if base == nil {
			reqMsg.Vaild = easygo.NewBool(false)
		} else {
			reqMsg.HeadIcon = easygo.NewString(base.GetHeadIcon())
			reqMsg.Sex = easygo.NewInt32(base.GetSex())
			reqMsg.Vaild = easygo.NewBool(true)
		}
	}
	return reqMsg
}

//=============================================================

//消息广播给其他大厅的玩家isSend:是否不解包直接发送给客户端，默认true
//func BroadCastMsgToOtherHallClient(serverId int32, playerIds []int64, memthName string, msg easygo.IMessage, isSend ...bool) {
//
//	otherHall := PServerInfoMgr.GetServerInfo(serverId)
//	if otherHall == nil {
//		logs.Error("大厅对象找不到，sid:", serverId)
//		return
//	}
//	var content []byte
//	if msg != nil {
//		b, err := msg.Marshal()
//		easygo.PanicError(err)
//		content = b
//	} else {
//		content = []byte{}
//	}
//	send := append(isSend, true)[0]
//	hallMsg := &hall_hall.MsgToClient{
//		PlayerIds: playerIds,
//		RpcName:   easygo.NewString(memthName),
//		MsgName:   easygo.NewString(proto.MessageName(msg)),
//		Msg:       content,
//		IsSend:    easygo.NewBool(send),
//	}
//	PWebApiForServer.SendToServer(otherHall, "RpcBroadCastMsgToClient", hallMsg)
//}

func GetCardInfo(pid PLAYER_ID, state int32, markName string) *share_message.TeamPlayerInfo {
	base := GetPlayerObj(pid)
	if base == nil {
		panic(fmt.Sprintf("base对象为空，id:%d", pid))
	}
	phone := base64.StdEncoding.EncodeToString([]byte(base.GetPhone()))
	dynamicData := &share_message.DynamicDataListPage{
		DynamicData: for_game.GetRedisDynamicForSomeLogId(false, pid, base.GetDynamicList()),
		TotalCount:  nil,
		PageCount:   nil,
	}
	trueZan := for_game.GetAllTrueZan(pid)
	info := &share_message.TeamPlayerInfo{
		PlayerId:     easygo.NewInt64(pid),
		Account:      easygo.NewString(base.GetAccount()),
		NickName:     easygo.NewString(base.GetNickName()),
		ReName:       easygo.NewString(markName),
		HeadIcon:     easygo.NewString(base.GetHeadIcon()),
		Sex:          easygo.NewInt32(base.GetSex()),
		Photo:        base.GetPhoto(),
		Phone:        easygo.NewString(phone),
		Signature:    easygo.NewString(base.GetSignature()),
		Provice:      easygo.NewString(base.GetProvice()),
		City:         easygo.NewString(base.GetCity()),
		State:        easygo.NewInt32(state),
		Zans:         easygo.NewInt32(base.GetZan() + trueZan + base.GetVCZanNum() + base.GetBsVCZanNum()),
		Fans:         easygo.NewInt32(len(base.GetFans())),
		Attentions:   easygo.NewInt32(len(base.GetAttention())),
		Icon:         easygo.NewInt32(0),
		DynamicData:  dynamicData, // 未分页数据
		AccountState: easygo.NewInt32(base.GetStatus()),
		//Constellation: easygo.NewString(for_game.GetConfigConstellationSortName(base.GetConstellation())),
		Constellation: easygo.NewInt32(base.GetConstellation()),

		//DynamicData: for_game.GetRedisDynamicForSomeLogId(pid, base.GetRedisPlayerDynamicListByPage(1, for_game.DefaultPageSize)),
	}
	return info
}

//func GetCardInfo1(pid PLAYER_ID, state, currentPage, pageSize int32, markName string, opId int64) *share_message.TeamPlayerInfo {
//	base := GetPlayerObj(pid)
//	if base == nil {
//		panic(fmt.Sprintf("base对象为空，id:%d", pid))
//	}
//	phone := base64.StdEncoding.EncodeToString([]byte(base.GetPhone()))
//	// 判断是否是第一页,如果是,叠加置顶动态.
//	pageMap := base.GetRedisPlayerDynamicListByPage(opId, currentPage, pageSize) // 分页
//	logIds := pageMap["arr"].([]int64)
//
//	dynamicList := for_game.GetRedisDynamicForSomeLogId1(opId, logIds, pid)
//	dynamicData := &share_message.DynamicDataListPage{
//		DynamicData: dynamicList,
//		TotalCount:  easygo.NewInt32(pageMap["totalCount"].(int32)),
//		PageCount:   easygo.NewInt32(pageMap["pageCount"].(int32)),
//	}
//	info := &share_message.TeamPlayerInfo{
//		PlayerId:           easygo.NewInt64(pid),
//		Account:            easygo.NewString(base.GetAccount()),
//		NickName:           easygo.NewString(base.GetNickName()),
//		ReName:             easygo.NewString(markName),
//		HeadIcon:           easygo.NewString(base.GetHeadIcon()),
//		Sex:                easygo.NewInt32(base.GetSex()),
//		Photo:              base.GetPhoto(),
//		Phone:              easygo.NewString(phone),
//		Signature:          easygo.NewString(base.GetSignature()),
//		Provice:            easygo.NewString(base.GetProvice()),
//		City:               easygo.NewString(base.GetCity()),
//		State:              easygo.NewInt32(state),
//		Zans:               easygo.NewInt32(base.GetZan()),
//		Fans:               easygo.NewInt32(len(base.GetFans())),
//		Attentions:         easygo.NewInt32(len(base.GetAttention())),
//		Icon:               easygo.NewInt32(0),
//		DynamicData:        dynamicData,
//		AccountState:       easygo.NewInt32(base.GetStatus()),
//		BackgroundImageURL: easygo.NewString(base.GetBackgroundImageURL()),
//		//DynamicData:        for_game.GetRedisDynamicForSomeLogId(pid, base.GetRedisPlayerDynamicListByPage(1, for_game.DefaultPageSize)),
//	}
//	return info
//}
func GetCardInfo1(opid, pid PLAYER_ID, state int32, markName string) *share_message.TeamPlayerInfo {
	base := for_game.GetRedisPlayerBase(pid)
	info := &share_message.TeamPlayerInfo{}
	if base == nil {
		logs.Error("base对象为空，id:%d", pid)
		return info
	}
	phone := base64.StdEncoding.EncodeToString([]byte(base.GetPhone()))
	trueZan := for_game.GetAllTrueZan(pid)
	// 判断是否是第一页,如果是,叠加置顶动态.
	var addFriendType int32
	if friendBase := for_game.GetFriendBase(opid); friendBase != nil {
		for _, v := range friendBase.GetFriends() {
			if v.GetPlayerId() == pid {
				addFriendType = v.GetType()
			}
		}
	}
	info = &share_message.TeamPlayerInfo{
		PlayerId:           easygo.NewInt64(pid),
		Account:            easygo.NewString(base.GetAccount()),
		NickName:           easygo.NewString(base.GetNickName()),
		ReName:             easygo.NewString(markName),
		HeadIcon:           easygo.NewString(base.GetHeadIcon()),
		Sex:                easygo.NewInt32(base.GetSex()),
		Photo:              base.GetPhoto(),
		Phone:              easygo.NewString(phone),
		Signature:          easygo.NewString(base.GetSignature()),
		Provice:            easygo.NewString(base.GetProvice()),
		City:               easygo.NewString(base.GetCity()),
		State:              easygo.NewInt32(state),
		Zans:               easygo.NewInt32(base.GetZan() + trueZan + base.GetVCZanNum() + base.GetBsVCZanNum()),
		Fans:               easygo.NewInt32(len(base.GetFans())),
		Attentions:         easygo.NewInt32(len(base.GetAttention())),
		Icon:               easygo.NewInt32(0),
		AccountState:       easygo.NewInt32(base.GetStatus()),
		BackgroundImageURL: easygo.NewString(base.GetBackgroundImageURL()),
		Types:              easygo.NewInt32(base.GetTypes()),
		AddFriendType:      easygo.NewInt32(addFriendType),
		Constellation:      easygo.NewInt32(base.GetConstellation()),
		//Constellation:      easygo.NewString(for_game.GetConfigConstellationSortName(base.GetConstellation())),
	}
	return info
}

func GetServerPlayerMap(playerList []PLAYER_ID, id PLAYER_ID) map[int32][]PLAYER_ID {
	info := make(map[int32][]PLAYER_ID)
	for _, pid := range playerList {
		if pid == id {
			continue
		}
		if PlayerOnlineMgr.CheckPlayerIsOnLine(pid) { //如果在线
			sid := PlayerOnlineMgr.GetPlayerServerId(pid)
			info[sid] = append(info[sid], pid)
		}
	}
	return info
}

func (self *cls1) RpcUserReturnApp(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	return easygo.EmptyMsg
}

//推送消息
func (self *cls1) RpcReturnTweets(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	tweetsList := who.GetSweets()
	if tweetsList == nil {
		return easygo.EmptyMsg
	}
	return tweetsList
}

//删除推文
func (self *cls1) RpcDelTweets(ep IGameClientEndpoint, who *Player, reqMsg *client_server.TweetsIdsRequest, common ...*base.Common) easygo.IMessage {

	sweetsIds := reqMsg.GetTweetsIdList()
	if len(sweetsIds) > 0 {
		who.DelSweetsByIds(sweetsIds)
	}
	return easygo.EmptyMsg
}

//推文客户端操作反馈
func (self *cls1) RpcSumClicksJumps(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.ArticleOptRequest, common ...*base.Common) easygo.IMessage {
	easygo.Spawn(func() {
		arc := for_game.QueryArticleById(reqMsg.GetArticleId())
		if arc != nil {
			switch reqMsg.GetType() {
			case 1:
				for_game.SetRedisArticleReportFildVal(reqMsg.GetArticleId(), 1, "Clicks") //添加文章报表消息点击数
			case 2:
				for_game.SetRedisArticleReportFildVal(reqMsg.GetArticleId(), 1, "Jumps") //添加文章报成功跳转数量
			}
		}
	})

	return easygo.EmptyMsg
}

//推送通知客户端操作反馈
func (self *cls1) RpcNoticeClicks(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.NoticeOptRequest, common ...*base.Common) easygo.IMessage {
	easygo.Spawn(func() {
		arc := for_game.QueryNoticeById(reqMsg.GetId())
		if arc != nil {
			switch reqMsg.GetType() {
			case 1:
				for_game.SetRedisNoticeReportFildVal(reqMsg.GetId(), 1, "Clicks") //添加推送通知报表消息点击数
			case 2:
				easygo.NewFailMsg("目前没有跳转")
			}
		}
	})

	return easygo.EmptyMsg
}

//修改、新增表情包
func (self *cls1) RpcModifyEmoticon(ep IGameClientEndpoint, who *Player, reqMsg *share_message.PlayerEmoticon, common ...*base.Common) easygo.IMessage {

	player := for_game.GetRedisPlayerBase(who.GetPlayerId())
	if player != nil {
		if player.ModifyRedisEmoticon(reqMsg) {
			return easygo.EmptyMsg
		}
	}
	return easygo.NewFailMsg("修改表情数据失败")
}

//删除表情包
func (self *cls1) RpcDelEmoticon(ep IGameClientEndpoint, who *Player, reqMsg *share_message.PlayerEmoticon, common ...*base.Common) easygo.IMessage {
	player := for_game.GetRedisPlayerBase(who.GetPlayerId())
	if player != nil {
		if player.DelRedisEmoticon(reqMsg.GetTypeId()) {
			return easygo.EmptyMsg
		}
	}
	return easygo.NewFailMsg("删除表情数据失败")
}

//==============================================================================================================================客服系统====>
//请求客服分类列表
func (self *cls1) RpcRequestWaiterTypes(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("请求客服分类列表,RpcRequestWaiterTypes:", reqMsg)
	lis := for_game.GetManagerTypesListForNormal()

	msg := &client_hall.WaiterTypesResponse{
		List: lis,
	}
	return msg
}

//请求客服人工服务
func (self *cls1) RpcRequestWaiterService(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.WaiterMsgRequest, common ...*base.Common) easygo.IMessage {
	logs.Info("请求人工服务,RpcRequestWaiterService:", reqMsg)
	// if reqMsg.Type == nil {
	// 	return easygo.NewFailMsg("咨询类型不能为空")
	// }
	if reqMsg.Type == nil || reqMsg.GetType() == 0 {
		reqMsg.Type = easygo.NewInt32(1)
	}
	userId := for_game.GetWaiterOnline(reqMsg.GetType()) //分配接收消息的客服Id

	msg := &client_hall.WaiterMsgResponse{
		WaiterId: easygo.NewInt64(userId),
	}

	return msg
}

// 查询指定玩家正在沟通的im消息
func (self *cls1) RpcGetWaiterMsg(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.WaiterMsgRequest, common ...*base.Common) easygo.IMessage {
	// logs.Info("==============RpcGetWaiterMsg", reqMsg)
	im := &share_message.IMmessage{}
	if reqMsg.Mid == nil || reqMsg.GetMid() == 0 {
		im = for_game.QueryIMmessageByPid(who.GetPlayerId(), for_game.WAITER_MESSAGE_ING)
	} else {
		im = for_game.QueryIMmessageByMid(reqMsg.GetMid())
	}

	if im != nil {
		if im.GetContent() != nil {
			im.SumContent = easygo.NewInt32(len(im.GetContent()))
		} else {
			im.SumContent = easygo.NewInt32(0)
		}
		new := im.GetCnew()
		if new != 0 && reqMsg.GetType() != 0 {
			content := im.GetContent()
			sum := len(content)
			start := sum - int(new)
			im.Content = content[start:]
			im.Cnew = easygo.NewInt32(0)
			for_game.UpdateMessageRead(im.GetId(), 1) //前端查看消息，修改消息已读状态
		} else {
			im.Content = []*share_message.IMcontent{}
		}
	} else {
		im = &share_message.IMmessage{SumContent: easygo.NewInt32(0)}
	}

	return im
}

//发送消息给后台服务器
func (self *cls1) RpcSendWaiterMsg(ep IGameClientEndpoint, who *Player, reqMsg *share_message.IMmessage, common ...*base.Common) easygo.IMessage {
	content := reqMsg.GetContent()[0]
	if content == nil {
		return easygo.NewFailMsg("不打算写点什么吗？")
	}

	if content.GetMtype() != for_game.WAITER_MSG_TYPE_C {
		return easygo.NewFailMsg("消息发送类型错误")
	}

	if content.Ctype == nil {
		return easygo.NewFailMsg("消息类型不能为空")
	}

	if reqMsg.WaiterId == nil || reqMsg.GetWaiterId() == 0 {
		return easygo.NewFailMsg("客服Id错误")
	}

	waiter := for_game.GetRedisWaiter(reqMsg.GetWaiterId())
	sid := int32(0)
	if waiter != nil {
		sid = waiter.ServerId
	}

	reqMsg.PlayerId = easygo.NewInt64(who.GetPlayerId())
	reqMsg.Account = easygo.NewString(who.GetAccount())
	reqMsg.NickName = easygo.NewString(who.GetNickName())
	reqMsg.HeadIcon = easygo.NewString(who.GetHeadIcon())
	if reqMsg.GetGrade() != 0 {
		reqMsg.Grade = easygo.NewInt32(0)
	}
	if reqMsg.GetStatus() == 0 {
		reqMsg.Status = easygo.NewInt32(1)
	}
	im := &share_message.IMmessage{}
	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		im = for_game.QueryIMmessageByPid(reqMsg.GetPlayerId(), for_game.WAITER_MESSAGE_ING)
		if im != nil {
			reqMsg.Id = easygo.NewInt64(im.GetId())
			reqMsg.WaiterId = easygo.NewInt64(im.GetWaiterId())
		} else {
			reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WAITER_MESSAGE))
			reqMsg.CreateTime = easygo.NewInt64(util.GetMilliTime())
		}
	} else {
		im = for_game.QueryIMmessageByMid(reqMsg.GetId())
		if im == nil {
			return easygo.NewFailMsg("消息ID错误")
		}
		if im.GetStatus() >= for_game.WAITER_MESSAGE_END {
			return easygo.NewFailMsg("本次服务已结束")
		}
	}

	reqMsg.Snew = easygo.NewInt32(im.GetSnew() + 1)

	easygo.Spawn(func() { SendMsgToServerNew(sid, "RpcHallSendIMmessage", reqMsg) })

	return &share_message.IMmessage{Id: easygo.NewInt64(reqMsg.GetId())}
}

//玩家给客服评分
func (self *cls1) RpcWaiterGrade(ep IGameClientEndpoint, who *Player, reqMsg *share_message.IMmessage, common ...*base.Common) easygo.IMessage {
	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("消息Id错误")
	} else {
		im := for_game.QueryIMmessageByMid(reqMsg.GetId())
		if im.GetStatus() < for_game.WAITER_MESSAGE_END {
			return easygo.NewFailMsg("本次服务未结束")
		}
		if im.GetStatus() == for_game.WAITER_MESSAGE_GRADE {
			return easygo.NewFailMsg("本次服务已评价")
		}
		if reqMsg.Grade == nil || reqMsg.GetGrade() < int32(-2) || reqMsg.GetGrade() > int32(2) {
			return easygo.NewFailMsg("评价分值错误")
		}

		easygo.Spawn(func() {
			for_game.GradeToWaiter(reqMsg.GetId(), reqMsg.GetGrade())             //修改消息评分
			for_game.UpdateWaiterPerformance(im.GetWaiterId(), reqMsg.GetGrade()) //修改客服绩效
		})

	}

	return easygo.EmptyMsg
}

//搜索常见问题
func (self *cls1) RpcSearchForKey(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SearchFaqRequest, common ...*base.Common) easygo.IMessage {
	key := reqMsg.GetKey()
	lis := for_game.QueryWaiterFAQList(key, for_game.FAQ_TITLE_SEARCH)
	if len(lis) == 0 {
		lis = for_game.QueryWaiterFAQList(key, for_game.FAQ_KEY_SEARCH)
	}

	msg := &client_hall.SearchFaqResponse{
		List: lis,
	}
	return msg
}

//打开常见问题详情
func (self *cls1) RpcOpenFaqById(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.OpenFaqRequest, common ...*base.Common) easygo.IMessage {
	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("Id错误")
	}
	one := for_game.OpenFaqById(reqMsg.GetId())
	return one
}
func (self *cls1) RpcMarkPlayer(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.MarkName, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcMarkPlayer:", reqMsg)
	name := reqMsg.GetName()
	if len(name) > 0 {
		b, _ := for_game.PDirtyWordsMgr.CheckWord(name)
		if b {
			return easygo.NewFailMsg("备注昵称存在敏感词")
		}
	}
	//给好友备注
	base := for_game.GetFriendBase(who.GetPlayerId())
	if base != nil {
		base.SetFriendReName(reqMsg.GetPlayerId(), name)
	}
	return easygo.EmptyMsg
}

// 打开我的主页
/*func (self *cls1) RpcOpenMyMainPage(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.OpenMyMainPageByPageReq, common... *base.Common) easygo.IMessage {
	logs.Info("===========打开我的主页 RpcOpenMyMainPage===========")
	currentPage := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()
	if currentPage == 0 {
		currentPage = 1
	}
	if pageSize == 0 {
		pageSize = for_game.DefaultPageSize
	}
	dynamicList := make([]*share_message.DynamicData, 0)
	// 如果是第一页,添加自己的置顶动态
	if currentPage == 1 {
		// todo 获取官方置顶动态
		// 获取后台置顶的keys
		keys := for_game.GetMyMainTopKeysFromRedis(who.GetPlayerId(), for_game.REDIS_SQUARE_BS_TOP_DYNAMIC)
		if len(keys) > 0 {
			slice := easygo.GetSliceByRandFromSlice(keys, for_game.BS_TOP_NUM)
			// 时间最新的最靠前
			sortSlice := for_game.SortDynamicSliceByTime(slice)
			logs.Info("返回后台置顶动态,排序后的----->", sortSlice)

			datas := for_game.GetRedisDynamicForSomeLogId(true, who.GetPlayerId(), sortSlice)
			dynamicList = append(dynamicList, datas...) // 置顶后的动态
		}
	}
	logIdMap := who.GetRedisPlayerDynamicListByPage(who.GetPlayerId(), currentPage, pageSize)
	logIds := logIdMap["arr"].([]int64)
	ds := for_game.GetRedisDynamicForSomeLogId(false, who.GetPlayerId(), logIds) // 没有置顶的动态
	dynamicList = append(dynamicList, ds...)
	dynamicData := &share_message.DynamicDataListPage{
		DynamicData: dynamicList,
		TotalCount:  easygo.NewInt32(logIdMap["totalCount"].(int32)),
		PageCount:   easygo.NewInt32(logIdMap["pageCount"].(int32)),
	}
	// 查询数据库获取所有动态的假赞数.
	trueZan := for_game.GetAllTrueZan(who.GetPlayerId())
	msg := &client_server.MyMainPageInfo{
		Fans:       easygo.NewInt32(len(who.GetFans())),
		Attentions: easygo.NewInt32(len(who.GetAttention())),
		Zans:       easygo.NewInt32(who.GetZan() + trueZan),
		//DynamicData: for_game.GetRedisDynamicForSomeLogId(who.GetPlayerId(), who.GetRedisPlayerDynamicList()), //  以前没有分页
		DynamicData: dynamicData,
		PlayerId:    easygo.NewInt64(who.GetPlayerId()),
	}
	return msg
}*/

// 新版本打开我的主页
func (self *cls1) RpcOpenMyMainPage(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.OpenMyMainPageByPageReq, common ...*base.Common) easygo.IMessage {
	logs.Info("===========新版本打开我的主页 RpcOpenMyMainPage===reqMst=%v,common=%v", reqMsg, common)

	// 查询数据库获取所有动态的假赞数.
	trueZan := for_game.GetAllTrueZan(who.GetPlayerId())
	msg := &client_server.MyMainPageInfo{
		Fans:       easygo.NewInt32(len(who.GetFans())),
		Attentions: easygo.NewInt32(len(who.GetAttention())),
		Zans:       easygo.NewInt32(who.GetZan() + trueZan),
		PlayerId:   easygo.NewInt64(who.GetPlayerId()),
	}
	com := append(common, nil)[0]
	if com != nil && com.GetVersion() == "2.7.4" { // 2.7.4 直接返回上半部分的内容.
		return msg
	}
	currentPage := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()
	if currentPage == 0 {
		currentPage = 1
	}
	if pageSize == 0 {
		pageSize = for_game.DefaultPageSize
	}
	dynamicList := make([]*share_message.DynamicData, 0)
	// 如果是第一页,添加自己的置顶动态
	if currentPage == 1 {
		// todo 获取官方置顶动态
		bsTopDynamicList := for_game.GetBSTopDynamicListByIDsFromDB(0, []int64{who.GetPlayerId()})
		if len(bsTopDynamicList) > 0 {
			slice := for_game.GetDynamicSliceByRandFromSlice(bsTopDynamicList, for_game.PLAYER_TOP_NUM)
			// 时间最新的最靠前
			sortSlice := for_game.SortDynamicSliceByTime1(slice)
			dynamicList = append(dynamicList, sortSlice...)
		}

	}
	//ds, count := for_game.GetNoTopDynamicByPIDsFromDB(int(currentPage), int(pageSize), []int64{who.GetPlayerId()})
	maxLogIdKey := for_game.MakeNewString(who.GetPlayerId(), reqMsg.GetPlayerId())
	ds, count := for_game.GetNoTopDynamicByPIDs(0, int(currentPage), int(pageSize), []int64{who.GetPlayerId()}, maxLogIdKey)

	dynamicList = append(dynamicList, ds...)
	ds = for_game.GetRedisDynamicForSomeLogId2(who.GetPlayerId(), dynamicList) // 没有置顶的动态
	dynamicData := &share_message.DynamicDataListPage{
		DynamicData: ds,
		TotalCount:  easygo.NewInt32(count),
	}
	// 异步处理动态的话题浏览量
	easygo.Spawn(func() {
		for_game.OperateViewNum(ds)
	})
	msg.DynamicData = dynamicData
	return msg
}

// 我的主页 分页获取动态信息
//func (self *cls1) RpcOpenMyMainPageByPage(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.OpenMyMainPageByPageReq, common... *base.Common) easygo.IMessage {
//	logs.Info("RpcOpenMyMainPageByPage 分页查询动态数据请求参数: %+v", reqMsg)
//	page := reqMsg.GetPage()
//	pageSize := reqMsg.GetPageSize()
//	playerID := reqMsg.GetPlayerId()
//	if page <= 0 || pageSize <= 0 {
//		return easygo.NewFailMsg("请求分页参数有误")
//	}
//	// 判断用户是否存在
//	player := for_game.GetRedisPlayerBase(playerID)
//	if player == nil {
//		logs.Error("获取玩家失败,playerID:---->", playerID)
//		return easygo.NewFailMsg("玩家不存在")
//	}
//	logIdMap := player.GetRedisPlayerDynamicListByPage(who.GetPlayerId(), page, pageSize) // 包含置顶的动态
//	logIds := logIdMap["arr"].([]int64)
//	var isTop bool
//	if page == 1 { // 需要置顶数据
//		isTop = true
//	}
//	dynamicData := &share_message.DynamicDataListPage{
//		DynamicData: for_game.GetRedisDynamicForSomeLogId(isTop, reqMsg.GetPlayerId(), logIds),
//		TotalCount:  easygo.NewInt32(logIdMap["totalCount"].(int32)),
//		PageCount:   easygo.NewInt32(logIdMap["pageCount"].(int32)),
//	}
//	msg := &client_server.MyMainPageInfo{
//		DynamicData: dynamicData,
//		PlayerId:    easygo.NewInt64(playerID),
//	}
//	return msg
//}

// 我的主页中, 分页查询动态数据请求参数
func (self *cls1) RpcOpenMyMainPageByPage(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.OpenMyMainPageByPageReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcOpenMyMainPageByPage 分页查询动态数据请求参数: %+v", reqMsg)
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()
	playerID := reqMsg.GetPlayerId()
	var opId int64
	if who.GetPlayerId() != reqMsg.GetPlayerId() { // 如果不相等,说明是看别人的主页
		opId = who.GetPlayerId()
	}

	if page <= 0 || pageSize <= 0 {
		return easygo.NewFailMsg("请求分页参数有误")
	}
	// 判断用户是否存在
	player := for_game.GetRedisPlayerBase(playerID)
	if player == nil {
		logs.Error("获取玩家失败,playerID:---->", playerID)
		return easygo.NewFailMsg("玩家不存在")
	}
	ds := make([]*share_message.DynamicData, 0)
	// 判断是否是第一页,是的话需要把置顶的动态加载出来
	if page == 1 {
		//获取后台置顶的动态列表(in)
		bsTopDynamicList := for_game.GetBSTopDynamicListByIDsFromDB(opId, []int64{playerID})
		if len(bsTopDynamicList) > 0 {
			slice := for_game.GetDynamicSliceByRandFromSlice(bsTopDynamicList, for_game.PLAYER_TOP_NUM)
			//  时间最新的最靠前
			sortSlice := for_game.SortDynamicSliceByTime1(slice)
			ds = append(ds, sortSlice...)
		}

	}
	//ds1, count := for_game.GetNoTopDynamicByPIDsFromDB(int(page), int(pageSize), []int64{playerID})
	maxLogIdKey := for_game.MakeNewString(who.GetPlayerId(), playerID) // 操作者的id_对方id
	ds1, count := for_game.GetNoTopDynamicByPIDs(opId, int(page), int(pageSize), []int64{playerID}, maxLogIdKey)
	ds = append(ds, ds1...)
	// 热门分数
	sysParam := PSysParameterMgr.GetSysParameter(for_game.SQUAREHOT_PARAMETER)
	var hotScore int32
	if sysParam != nil {
		hotScore = sysParam.GetHotScore()
	}
	// 计算热门.
	ds = for_game.ParseHotDynamic(ds, hotScore)
	dynamicData := &share_message.DynamicDataListPage{
		DynamicData: for_game.GetRedisDynamicForSomeLogId2(playerID, ds, who.GetPlayerId()), // 没有置顶的动态,
		TotalCount:  easygo.NewInt32(count),
	}
	msg := &client_server.MyMainPageInfo{
		DynamicData: dynamicData,
		PlayerId:    easygo.NewInt64(playerID),
	}

	//  异步处理动态的话题浏览量
	easygo.Spawn(func() {
		for_game.OperateViewNum(dynamicData.GetDynamicData())
	})

	return msg
}

// 获取粉丝关注信息列表
func (self *cls1) RpcGetMyFansAttentionInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.MainInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("=========获取粉丝关注信息列表RpcGetMyFansAttentionInfo========", reqMsg)
	if reqMsg.GetPlayerId() != who.GetPlayerId() {
		return easygo.NewFailMsg("只有本人可以查看这些消息")
	}

	msg := &client_hall.AllFansInfo{}
	t := reqMsg.GetType()
	page := reqMsg.GetPage()
	count := reqMsg.GetCount()
	start := int((page - 1) * count)
	end := int(page * count)
	var pids []int64
	if t == 1 { //粉丝
		pids = who.GetFans()
	} else { //关注
		pids = who.GetAttention()
	}
	//easygo.SortSliceInt64(pids, false)
	ns := pids
	leng := len(ns)
	if start >= leng {
		return msg
	}
	var ids []int64
	if leng > end { //如果超过请求数量
		ids = ns[start : end-1]
	} else {
		ids = ns[start:]
	}

	lst := []*client_hall.FansInfo{}
	playerInfo := for_game.GetAllPlayerBase(ids)
	for _, id := range ids {
		player, ok := playerInfo[id]
		if !ok {
			continue
		}
		m := &client_hall.FansInfo{
			PlayerId:  easygo.NewInt64(id),
			Sex:       easygo.NewInt32(player.GetSex()),
			Signature: easygo.NewString(player.GetSignature()),
			Name:      easygo.NewString(player.GetNickName()),
			HeadIcon:  easygo.NewString(player.GetHeadIcon()),
			Types:     easygo.NewInt32(player.GetTypes()),
		}
		lst = append(lst, m)
	}
	msg.FansInfo = lst
	return msg
}

//
func (self *cls1) RpcGetSomeDynamic(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.DynamicIdInfo, common ...*base.Common) easygo.IMessage {
	pid := reqMsg.GetPlayerId()
	t := reqMsg.GetType()
	logId := reqMsg.GetLogId()
	msg := for_game.GetRedisPlayerSomeDynamic(pid, logId, t)
	return msg
}
func (self *cls1) RpcGetTeamMembers(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.TeamMembers, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetTeamMembers:", reqMsg)
	members := for_game.GetAllTeamMember(reqMsg.GetTeamId(), who.GetPlayerId())
	reqMsg.Members = members
	return reqMsg
}

//检测敏感词
func (self *cls1) RpcCheckDirtyWord(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CheckDirtyWord, common ...*base.Common) easygo.IMessage {
	logs.Info("请求:", reqMsg)
	dirtyWords := for_game.PDirtyWordsMgr.CheckWords(reqMsg.GetWords())
	logs.Info("检测到:", dirtyWords)
	reqMsg.DirtyWords = dirtyWords
	return reqMsg
}

//注销账号检测
func (self *cls1) RpcCheckCancelAccount(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("请求CheckCancelAccount:", who.GetPlayerId())
	//绑定手机号和实名认证的号才能注销
	//phone := who.GetPhone() != ""
	//peopleId := who.GetPeopleId() != ""
	msg := &client_hall.CheckCancelAccount{
		//ThirtyDays:    easygo.NewBool(who.CheckThirtyDays()),
		//PhoneState:    easygo.NewBool(phone),
		//PeopleIdState: easygo.NewBool(peopleId),
		AccountState: easygo.NewBool(who.CheckAccountStatus()),
		TradeState:   easygo.NewBool(who.CheckShopState()),
		BalanceState: easygo.NewBool(who.CheckBalanceState()),
		//FriendState:   easygo.NewBool(who.CheckFriendTeamState()),
	}
	logs.Info("返回结果:", msg)
	return msg
}

//收到注销账号数据请求
func (self *cls1) RpcSubmitCancelAccount(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CancelAccountData, common ...*base.Common) easygo.IMessage {
	// 验证短信验证码
	//if for_game.IS_FORMAL_SERVER {
	//	if err := for_game.CheckMessageCode(who.GetPhone(), reqMsg.GetPhoneCode(), for_game.CLIENT_CODE_CANCELACCOUNT); err != nil {
	//		return err
	//	}
	//}

	//if reqMsg.GetRealName() != who.GetRealName() || reqMsg.GetPeopleId() != who.GetPeopleId() {
	//	return easygo.NewFailMsg("实名认证信息不符")
	//}

	log := &share_message.PlayerCancleAccount{
		Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_CANCEL_ACCOUNT)),
		Account:    easygo.NewString(who.GetAccount()),        // 账号
		Phone:      easygo.NewString(who.GetPhone()),          // 电话号码
		PlayerId:   easygo.NewInt64(who.GetPlayerId()),        //id
		CreateTime: easygo.NewInt64(for_game.GetMillSecond()), // 创建时间
		//RealName:          easygo.NewString(reqMsg.GetRealName()),           //真实姓名
		//PeopleId:          easygo.NewString(reqMsg.GetPeopleId()),           //身份证
		//PeopleIdBeforeUrl: easygo.NewString(reqMsg.GetPeopleIdBeforeUrl()),  //正面图
		//PeopleIdBackUrl:   easygo.NewString(reqMsg.GetPeopleIdBackUrl()),    //反面图
		//PeopleIdHandUrl:   easygo.NewString(reqMsg.PeopleIdHandUrl),         //手持图
		Status: easygo.NewInt32(for_game.ACCOUNT_CANCEL_WAITING), //0待处理，1完成，2已拒绝,3已取消
	}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CANCEL_ACCOUNT)
	defer closeFun()
	err := col.Insert(log)
	easygo.PanicError(err)
	//用户账号状态变为:注销中
	who.SetStatus(for_game.ACCOUNT_CANCELING)
	return nil
}

//关注某人
func (self *cls1) RpcAttentioPlayer(ep IGameClientEndpoint, who *Player, reqMsg *client_server.AttenInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("==================大厅关注某人RpcAttentioPlayer==================", reqMsg)
	pid := reqMsg.GetPlayerId() // 被关注人的id
	t := reqMsg.GetType()
	player := for_game.GetRedisPlayerBase(pid)
	if t == for_game.DYNAMIC_OPERATE { // 注销了的账号不给关注.
		//注销的账号不给关注
		if player == nil || player.GetStatus() == for_game.ACCOUNT_CANCELED {
			return easygo.NewFailMsg("该账号异常")
		}
	}
	// 1-关注,2-取消关注柱
	if t != for_game.DYNAMIC_OPERATE && t != for_game.DYNAMIC_DELOPERATE {
		logs.Error("大厅关注操作,操作类型期待的是1或者2,传过来的值为: ", t)
		return easygo.NewFailMsg("操作类型有误")
	}
	if pid == 0 {
		logs.Error("关注操作,前端传的playerId为空")
		return easygo.NewFailMsg("关注操作失败")
	}

	if pid == who.GetPlayerId() {
		return easygo.NewFailMsg("不能关注自己")
	}
	if t == for_game.DYNAMIC_OPERATE && util.Int64InSlice(pid, who.GetAttention()) {
		return easygo.NewFailMsg("已关注该玩家")
	}
	b := for_game.OperateRedisDynamicAttention(t, who.GetPlayerId(), pid)
	if !b {
		return easygo.NewFailMsg("关注操作失败")
	}

	if t == for_game.DYNAMIC_OPERATE { // 关注
		who.AddAttention(pid)
		player := for_game.GetRedisPlayerBase(pid)
		if player == nil {
			panic("玩家怎么会为空")
		}
		player.AddFans(who.GetPlayerId())
	} else { //取消关注
		who.DelAttention(pid)
		player.DelFans(who.GetPlayerId())
	}
	redMsgCallClient(ep, pid) // 大厅通知社交广场,显示红点.
	return reqMsg
}

//点赞操作
func (self *cls1) RpcZanOperateSquareDynamic(ep IGameClientEndpoint, who *Player, reqMsg *client_server.ZanInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("=================大厅点赞操作 RpcZanOperateSquareDynamic=================", reqMsg)
	logId := reqMsg.GetLogId()
	t := reqMsg.GetType()
	dynamic := for_game.GetRedisDynamic(logId)
	if dynamic == nil {
		return easygo.NewFailMsg("该动态已被删除")
	}
	if for_game.GetRedisDynamicIsZan(logId, who.GetPlayerId()) && t == for_game.DYNAMIC_OPERATE {
		return easygo.NewFailMsg("你已经点赞过该动态")
	}
	dynamicPid := dynamic.GetPlayerId()
	sysp := PSysParameterMgr.GetSysParameter(for_game.PUSH_PARAMETER)
	b := for_game.OperateRedisDynamicZan(t, who.GetPlayerId(), logId, dynamicPid, sysp)
	if !b {
		return easygo.NewFailMsg("点赞操作失败")
	}

	if dynamicPid != who.GetPlayerId() { //自己操作  不提示红点

		unReadMsg := &client_hall.NewUnReadMessageResp{
			UnreadComment:   easygo.NewInt32(for_game.GetPlayerUnreadInfo(dynamicPid, for_game.UNREAD_COMMENT)),
			UnreadZan:       easygo.NewInt32(for_game.GetPlayerUnreadInfo(dynamicPid, for_game.UNREAD_ZAN)),
			UnreadAttention: easygo.NewInt32(for_game.GetPlayerUnreadInfo(dynamicPid, for_game.UNREAD_ATTENTION)),
		}
		ep.RpcNewMessage(unReadMsg)
		//BroadCastMsgToOtherSquareClient([]int64{playerId}, "RpcNewMessage", true, unReadMsg)

		//redMsgCallClient(ep,dynamicPid)
	}
	for_game.AddSaveDynamicIdList(logId)
	if t == for_game.DYNAMIC_OPERATE {
		// 添加埋点数据
		easygo.Spawn(func() {
			// 获取动态热门参数
			p := for_game.QuerySysParameterById(for_game.SQUAREHOT_PARAMETER)
			if p == nil {
				return
			}
			// 对该动态的分数加累加
			for_game.UpsetDynamicScore(logId, p.GetZanScore())
		})
	}

	// 异步操作话题点赞
	easygo.Spawn(
		func() {
			content := dynamic.GetContent() //  动态内容.
			// 解码
			contentBytes, _ := base64.StdEncoding.DecodeString(content)
			switch t {
			case for_game.DYNAMIC_OPERATE: // 对话题加参与数
				for_game.OperateTopicParticipationNum(string(contentBytes), 1)
			case for_game.DYNAMIC_DELOPERATE: // 对话题减参与数
				for_game.OperateTopicParticipationNum(string(contentBytes), -1)
			}
		})

	return nil
}

// redMsgCallClient 未读消息下发给客户端公共方法
func redMsgCallClient(ep IGameClientEndpoint, playerId for_game.PLAYER_ID) {
	unReadMsg := &client_hall.NewUnReadMessageResp{
		UnreadComment:   easygo.NewInt32(for_game.GetPlayerUnreadInfo(playerId, for_game.UNREAD_COMMENT)),
		UnreadZan:       easygo.NewInt32(for_game.GetPlayerUnreadInfo(playerId, for_game.UNREAD_ZAN)),
		UnreadAttention: easygo.NewInt32(for_game.GetPlayerUnreadInfo(playerId, for_game.UNREAD_ATTENTION)),
	}

	//BroadCastMsgToOtherSquareClient([]int64{playerId}, "RpcNewMessage", true, unReadMsg)
	ep.RpcNewMessage(unReadMsg)
}

//添加评论
func (self *cls1) RpcAddCommentSquareDynamic(ep IGameClientEndpoint, who *Player, reqMsg *share_message.CommentData, common ...*base.Common) easygo.IMessage {
	logs.Info("================大厅添加评论 RpcAddCommentSquareDynamic================,msg=%v", reqMsg)
	logId := reqMsg.GetLogId()
	playerId := who.GetPlayerId()
	dynamic := for_game.GetRedisDynamic(logId)
	if dynamic == nil {
		return easygo.NewFailMsg("该动态已被删除")
	}
	belongId := reqMsg.GetBelongId()
	if belongId == 0 { //代表是主评论
		if reqMsg.GetTargetId() != reqMsg.GetOwnerId() {
			return easygo.NewFailMsg("主评论的targetId 怎么不是发动态的任务id")
		}
	} else {
		comment := for_game.GetRedisDynamicComment(logId, belongId)
		if comment == nil {
			return easygo.NewFailMsg("该评论不存在")
		}
	}

	Id := for_game.NextId(for_game.TABLE_SQUARE_COMMENT)
	reqMsg.Id = easygo.NewInt64(Id)
	reqMsg.PlayerId = easygo.NewInt64(playerId)
	reqMsg.CreateTime = easygo.NewInt64(time.Now().Unix())
	reqMsg.Statue = easygo.NewInt32(for_game.DYNAMIC_COMMENT_STATUE_COMMON)
	reqMsg.Score = easygo.NewInt32(0)
	text := base64.StdEncoding.EncodeToString([]byte(reqMsg.GetContent()))
	reqMsg.Content = easygo.NewString(text)
	b := for_game.AddRedisDynamicComment(reqMsg)
	if !b {
		return easygo.NewFailMsg("生成评论失败")
	}
	ownerId := reqMsg.GetOwnerId() //发布动态的人的id

	for_game.AddSaveDynamicIdList(logId)
	if ownerId != playerId { //如果发布动态的人跟评论动态的人不是同一个人，那就添加到消息列表中
		for_game.AddRedisPlayerMessageId(ownerId, logId, Id)
		for_game.AddPlayerUnreadInfo(ownerId, for_game.UNREAD_COMMENT, 1)
	}
	if reqMsg.GetTargetId() != playerId && ownerId != reqMsg.GetTargetId() { //如果被评论的人，不是评论的人 并且发布动态的人不等于被评论的人
		for_game.AddRedisPlayerMessageId(reqMsg.GetTargetId(), logId, Id)
		for_game.AddPlayerUnreadInfo(reqMsg.GetTargetId(), for_game.UNREAD_COMMENT, 1)
	}

	if playerId != ownerId { //自己评论自己不提示红点
		redMsgCallClient(ep, ownerId)
	}
	// A发布评论,B去评论,C回复B,B需要显示红点
	if belongId != 0 {
		if ep1 := ClientEpMp.LoadEndpoint(reqMsg.GetTargetId()); ep1 != nil {
			redMsgCallClient(ep, reqMsg.GetTargetId())
		}
	}

	// 添加埋点数据
	easygo.Spawn(func() {
		// 获取动态热门参数
		p := for_game.QuerySysParameterById(for_game.SQUAREHOT_PARAMETER)
		if p == nil {
			return
		}
		// 对该动态的分数加累加
		for_game.UpsetDynamicScore(logId, p.GetCommentScore())
	})

	// 异步操作话题参与数
	easygo.Spawn(
		func() {
			content := dynamic.GetContent() //  动态内容.
			// 解码
			contentBytes, _ := base64.StdEncoding.DecodeString(content)
			for_game.OperateTopicParticipationNum(string(contentBytes), 1)
		})

	// 添加评论推送
	easygo.Spawn(func() {
		if playerId == reqMsg.GetTargetId() {
			return
		}
		if owner := for_game.GetRedisPlayerBase(ownerId); owner != nil {
			// 如果开关关了,不推送
			if owner.GetIsOpenSquare() {
				return
			}
		}
		m := for_game.PushMessage{
			//Title:       "温馨提示",
			ContentType: for_game.JG_TYPE_SQUARE,
			Location:    14, // 跳转到评论
			ObjectId:    logId,
			CommentId:   Id,
		}
		sysp := PSysParameterMgr.GetSysParameter(for_game.PUSH_PARAMETER)
		var ids []string
		// 主评论
		if belongId == 0 {
			m.ItemId = for_game.PUSH_ITEM_202 // 点赞我的动态
			ids = for_game.GetJGIds([]int64{ownerId})
		}

		// 判断是主评论还是回复的评论
		//if belongId != 0 && reqMsg.GetTargetId() != ownerId { // 表示回复
		if belongId != 0 { // 表示回复
			m.ItemId = for_game.PUSH_ITEM_203 // 点赞我的动态
			ids = for_game.GetJGIds([]int64{reqMsg.GetTargetId()})
		}

		if sysp == nil || len(sysp.GetPushSet()) == 0 {
			return
		}
		oper := for_game.GetRedisPlayerBase(playerId)

		for _, ps := range sysp.GetPushSet() {
			if ps.GetObjId() == m.ItemId && ps.GetIsPush() {
				m.Content = fmt.Sprintf(ps.GetObjContent(), oper.GetNickName())
			}
		}
		for_game.JGSendMessage(ids, m, sysp)
	})

	return reqMsg
}

//删除评论
func (self *cls1) RpcDelCommentSquareDynamic(ep IGameClientEndpoint, who *Player, reqMsg *client_server.IdInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("==========大厅删除评论 RpcDelCommentSquareDynamic==================", reqMsg)
	logId := reqMsg.GetId()
	Id := reqMsg.GetMainId()

	comment := for_game.GetRedisDynamicComment(logId, Id)
	if comment != nil {
		ownerId := comment.GetOwnerId()
		for_game.AddPlayerUnreadInfo(ownerId, for_game.UNREAD_COMMENT, -1)
		/*		if !for_game.GetIsNewMessage(ownerId) {
				ep1 := ClientEpMp.LoadEndpoint(ownerId)
				if ep1 != nil {
					ep1.RpcNoNewMessage(nil)
				}
			}*/

		redMsgCallClient(ep, ownerId) // 通知前端有咩有红点

		targetId := comment.GetTargetId()
		if targetId != ownerId {
			for_game.AddPlayerUnreadInfo(targetId, for_game.UNREAD_COMMENT, -1)
			/*		if !for_game.GetIsNewMessage(targetId) {
					ep1 := ClientEpMp.LoadEndpoint(targetId)
					if ep1 != nil {
						ep1.RpcNoNewMessage(nil)
					}
				}*/
			redMsgCallClient(ep, ownerId) // 通知前端有咩有红点

		}
	}

	b := for_game.DelRedisDynamicComment(logId, Id, for_game.DYNAMIC_COMMENT_STATUE_DELETE_CLIENT)
	if !b {
		return easygo.NewFailMsg("删除评论失败")
	}

	// 异步操作话题参与数
	easygo.Spawn(
		func() {
			content := for_game.GetRedisDynamic(logId).GetContent() //动态内容.
			// 解码
			contentBytes, _ := base64.StdEncoding.DecodeString(content)
			for_game.OperateTopicParticipationNum(string(contentBytes), -1)
		})

	return nil
}

// 删除动态
func (self *cls1) RpcDelSquareDynamic(ep IGameClientEndpoint, who *Player, reqMsg *client_server.RequestInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcDelSquareDynamic,msg=%v", reqMsg)
	logId := reqMsg.GetId()
	b := for_game.DelRedisSquareDynamic(who.GetPlayerId(), logId, "")
	if !b {
		return easygo.NewFailMsg("删除动态失败")
	}
	who.DelRedisPlayerDynamicList(logId)

	// 异步操作话题参与数
	easygo.Spawn(
		func() {
			if d := for_game.GetRedisDynamic(logId); d != nil {
				content := d.GetContent() //动态内容.
				// 解码
				if content != "" {
					contentBytes, _ := base64.StdEncoding.DecodeString(content)
					for_game.OperateTopicParticipationNum(string(contentBytes), -1)
				}
			}
		})
	return nil
}

//获取动态详情
func (self *cls1) RpcGetDynamicInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_server.IdInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("=============大厅获取动态详情 RpcGetDynamicInfo=====================", reqMsg)
	logId := reqMsg.GetId()
	msg := for_game.GetRedisDynamicAllInfo(0, logId, who.GetPlayerId(), who.GetAttention())

	//  异步处理动态的话题浏览量
	easygo.Spawn(func() {
		data := make([]*share_message.DynamicData, 0)
		data = append(data, msg)
		for_game.OperateViewNum(data)
	})
	return msg
}

//获取下一页动态主评论
func (self *cls1) RpcGetDynamicMainComment(ep IGameClientEndpoint, who *Player, reqMsg *client_server.IdInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetDynamicMainComment", reqMsg)
	logId := reqMsg.GetId()
	mainId := reqMsg.GetMainId()
	msg := for_game.GetRedisDynamicCommentInfo(logId, mainId)
	return msg
}

//获取动态子评论
func (self *cls1) RpcGetDynamicSecondaryComment(ep IGameClientEndpoint, who *Player, reqMsg *client_server.IdInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("================大厅获取动态子评论RpcGetDynamicSecondaryComment================", reqMsg)
	logId := reqMsg.GetId()
	mainId := reqMsg.GetMainId()
	secondId := reqMsg.GetSecondId()
	msg := for_game.GetRedisDynamicSecondComment(logId, mainId, secondId)
	return msg
}

//====================================新版动态评论获取=================================
//获取动态详情
func (self *cls1) RpcGetDynamicInfoNew(ep IGameClientEndpoint, who *Player, reqMsg *client_server.IdInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("=============大厅获取动态详情 RpcGetDynamicInfoNew=====================", reqMsg)
	logId := reqMsg.GetId()
	msg := for_game.GetRedisDynamicAllInfo(reqMsg.GetJumpMainCommentId(), logId, who.GetPlayerId(), who.GetAttention())
	sysParam := PSysParameterMgr.GetSysParameter(for_game.SQUAREHOT_PARAMETER)
	var hotScore int32
	if sysParam != nil {
		hotScore = sysParam.GetHotScore()
	}
	if msg.GetHostScore() >= hotScore && easygo.NowTimestamp()-msg.GetSendTime() < 7*86400 { // 大于7天的不显示热门. *86400
		msg.HotType = easygo.NewInt32(for_game.HOT_TYPE_1)
	}
	base := who.GetRedisPlayerBase()
	easygo.Spawn(for_game.SetPvUvCount, logId, base.GetDeviceCode())
	//  异步处理动态的话题浏览量
	easygo.Spawn(func() {
		data := make([]*share_message.DynamicData, 0)
		data = append(data, msg)
		for_game.OperateViewNum(data)
	})
	return msg
}

//获取下一页动态主评论
func (self *cls1) RpcGetDynamicMainCommentNew(ep IGameClientEndpoint, who *Player, reqMsg *client_server.IdInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("================大厅获取动态评论 RpcGetDynamicMainCommentNew================", reqMsg)
	logId := reqMsg.GetId()
	params := PSysParameterMgr.GetSysParameter("squarehot_parameter")
	msg := for_game.GetRedisDynamicCommentInfoByPage(who.GetPlayerId(), logId, reqMsg, params)
	return msg
}

//获取动态子评论
func (self *cls1) RpcGetDynamicSecondaryCommentNew(ep IGameClientEndpoint, who *Player, reqMsg *client_server.IdInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("================大厅获取动态子评论RpcGetDynamicSecondaryCommentNew================", reqMsg)
	logId := reqMsg.GetId()
	mainId := reqMsg.GetMainId()
	secondId := reqMsg.GetSecondId()
	msg := for_game.GetRedisDynamicSecondComment(logId, mainId, secondId)
	return msg
}

//====================================新版动态评论获取=================================
//订单投诉详情取得
func (self *cls1) RpcShopOrderComplaintDetail(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.ComplaintID, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcShopOrderComplaintDetail ", reqMsg)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_COMPLAINT)
	defer closeFun()

	var complaint *share_message.PlayerComplaint = &share_message.PlayerComplaint{}
	e := col.Find(bson.M{"_id": reqMsg.GetComplaintId()}).Limit(1).One(complaint)

	if e != nil || nil == complaint {
		logs.Error(e)
		return nil
	}

	colOrder, closeFunOrder := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFunOrder()

	var shopOrder *share_message.TableShopOrder = &share_message.TableShopOrder{}
	errOrder := colOrder.Find(bson.M{"_id": complaint.GetOrderId()}).Limit(1).One(shopOrder)
	if errOrder != nil || nil == shopOrder {
		logs.Error(errOrder)
		return nil
	}

	playInfo := for_game.GetRedisPlayerBase(complaint.GetPlayerId())
	var complaintAvatar string = playInfo.GetHeadIcon()
	var complaintNickName string = playInfo.GetNickName()
	return &client_hall.ShopOrderComplaintDetailRsp{
		SponsorNickname:    easygo.NewString(shopOrder.GetSponsorNickname()),
		SponsorAvatar:      easygo.NewString(shopOrder.GetSponsorAvatar()),
		ItemFile:           shopOrder.GetItems().GetItemFile(),
		ItemName:           easygo.NewString(shopOrder.GetItems().GetName()),
		ItemTitle:          easygo.NewString(shopOrder.GetItems().GetTitle()),
		OrderCreateTime:    easygo.NewInt64(shopOrder.GetCreateTime()),
		OrderId:            easygo.NewInt64(shopOrder.GetOrderId()),
		ComplaintAvatar:    easygo.NewString(complaintAvatar),
		ComplaintNickname:  easygo.NewString(complaintNickName),
		ComplaintContent:   easygo.NewString(complaint.GetContent()),
		ComplaintReContent: easygo.NewString(complaint.GetReContent()),
	}
}

//获取广告数据
func (self *cls1) RpcGetAdvData(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.AdvSettingRequest, common ...*base.Common) easygo.IMessage {
	logs.Info("================RpcGetAdvData 请求广告数据================", reqMsg)
	loc := reqMsg.GetLocation()
	list1 := for_game.QueryAdvListToDB(loc, ep.GetConnection().RemoteAddr().String())
	list := make([]*share_message.AdvSetting, 0)

	for _, v := range list1 {
		if strings.HasSuffix(v.GetJumpUrl(), "wishingActivity") {
			b1, b2, b3 := for_game.CheckWishActIsOpen()
			if !b1 || !b2 || !b3 {
				continue
			}
		}

		if strings.HasSuffix(v.GetJumpUrl(), "payActivity") {
			at := for_game.GetActivityByType(6)
			t := time.Now().Unix()
			if t < at.GetStartTime() || t > at.GetEndTime() {
				continue
			}
			if at.GetStatus() == 1 {
				continue
			}
		}

		list = append(list, v)
	}
	msg := &client_hall.AdvSettingResponse{List: list}
	// banner 广告展示埋点+1
	// 展示次数+1,展示人数去重
	if len(list) > 0 {
		easygo.Spawn(func() {
			pId := who.GetPlayerId()
			opTime := time.Now().Unix()
			for _, adv := range list {
				id := for_game.NextId(for_game.TABLE_ADV_LOG)
				id2 := for_game.NextId(for_game.TABLE_ADV_LOG)
				advLog := new(share_message.AdvLogReq)
				advLog.Id = easygo.NewInt64(id)
				advLog.AdvId = easygo.NewInt64(adv.GetId())
				advLog.PlayerId = easygo.NewInt64(pId)
				advLog.OpType = easygo.NewInt32(for_game.ADV_LOG_OP_TYPE_1)
				advLog.OpTime = easygo.NewInt64(opTime)
				for_game.AddAdvLogToDB(advLog)
				// 展示人数去重
				if log := for_game.GetAdvLogByPidAndOpFromDB(pId, adv.GetId(), for_game.ADV_LOG_OP_TYPE_2); log != nil {
					continue
				}
				advLog.Id = easygo.NewInt64(id2)
				advLog.OpType = easygo.NewInt32(for_game.ADV_LOG_OP_TYPE_2)
				for_game.AddAdvLogToDB(advLog)
				continue
			}
		})
	}

	return msg
}

// 添加广告埋点数据
func (self *cls1) RpcAddAdvLog(ep IGameClientEndpoint, who *Player, reqMsg *share_message.AdvLogReq, common ...*base.Common) easygo.IMessage {
	logs.Info("================添加广告埋点数据================", reqMsg)
	opType := reqMsg.GetOpType()
	if opType != for_game.ADV_LOG_OP_TYPE_1 && opType != for_game.ADV_LOG_OP_TYPE_3 {
		logs.Error("大厅广告埋点数据,操作类型有误,opType: ", opType)
		return easygo.NewFailMsg("操作类型有误")
	}
	reqMsg.PlayerId = easygo.NewInt64(who.GetPlayerId())
	reqMsg.OpTime = easygo.NewInt64(time.Now().Unix())

	// 前端只传1-展示,3-点击
	switch opType {
	case for_game.ADV_LOG_OP_TYPE_1: // 展示次数+1,判断是否需要去重.
		id := for_game.NextId(for_game.TABLE_ADV_LOG)
		id2 := for_game.NextId(for_game.TABLE_ADV_LOG)
		reqMsg.Id = easygo.NewInt64(id)
		for_game.AddAdvLogToDB(reqMsg)
		// 展示人数去重
		if log := for_game.GetAdvLogByPidAndOpFromDB(who.GetPlayerId(), reqMsg.GetAdvId(), for_game.ADV_LOG_OP_TYPE_2); log != nil {
			return nil
		}
		reqMsg.Id = easygo.NewInt64(id2)
		reqMsg.OpType = easygo.NewInt32(for_game.ADV_LOG_OP_TYPE_2)
		for_game.AddAdvLogToDB(reqMsg)
		return nil
	case for_game.ADV_LOG_OP_TYPE_3:
		id := for_game.NextId(for_game.TABLE_ADV_LOG)
		id2 := for_game.NextId(for_game.TABLE_ADV_LOG)
		reqMsg.Id = easygo.NewInt64(id)
		for_game.AddAdvLogToDB(reqMsg)
		// 点击人数去重
		if log := for_game.GetAdvLogByPidAndOpFromDB(who.GetPlayerId(), reqMsg.GetAdvId(), for_game.ADV_LOG_OP_TYPE_4); log != nil {
			return nil
		}
		reqMsg.Id = easygo.NewInt64(id2)
		reqMsg.OpType = easygo.NewInt32(for_game.ADV_LOG_OP_TYPE_4)
		for_game.AddAdvLogToDB(reqMsg)
		return nil
	}
	return nil
}

// 获取最新的群设置信息
func (self *cls1) RpcNewTeamSetting(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.NewTeamSettingReq, common ...*base.Common) easygo.IMessage {
	logs.Info("===获取最新的群设置信息 RpcNewTeamSetting=====,reqMsg=%v", reqMsg)
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		logs.Error("获取最新的群设置信息,该群不存在,群id: ", teamId)
		return easygo.NewFailMsg("群不存在")
	}
	setting := teamObj.GetTeamMessageSetting()
	notify := &client_hall.TeamSettingNotify{
		TeamId:         easygo.NewInt64(teamObj.GetId()),
		WelcomeWord:    easygo.NewString(teamObj.GetWelcomeWord()),
		MessageSetting: setting,
	}
	return notify
}

// 2.7.6 附近的人.
func (self *cls1) RpcLocationInfoNew(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.LocationInfoNewReq, common ...*base.Common) easygo.IMessage {
	logs.Info("===2.7.6 附近的人 RpcLocationInfoNew=====,reqMsg=%v", reqMsg)
	begin := for_game.GetMillSecond()
	pl := for_game.GetRedisPlayerBase(who.GetPlayerId())
	if reqMsg.GetX() == 0 || reqMsg.GetY() == 0 {
		logs.Error("经纬度为空,默认设置成,经度 x:%v,维度: %v", for_game.NEAR_DEFAULT_LNG, for_game.NEAR_DEFAULT_LAT) // 深圳
		reqMsg.X = easygo.NewFloat64(for_game.NEAR_DEFAULT_LNG)
		reqMsg.Y = easygo.NewFloat64(for_game.NEAR_DEFAULT_LAT)
	}
	if reqMsg.GetSex() != for_game.PLAYER_SEX_ALL && reqMsg.GetSex() != for_game.PLAYER_SEX_BOY && reqMsg.GetSex() != for_game.PLAYER_SEX_GIRL {
		logs.Error("性别有误,sex :", reqMsg.GetSex())
		return easygo.NewFailMsg("性别有误")
	}
	// 计算距离是否大于1km
	var b bool
	geoJson := pl.GeRedisPoints()
	if geoJson != nil && len(geoJson.Coordinates) == 2 { // x,y
		distance := for_game.GetDistance(geoJson.Coordinates[1], reqMsg.GetY(), geoJson.Coordinates[0], reqMsg.GetX())
		logs.Info("和上次的距离对比,相差-------->", distance)
		if int(distance) > 1000 { // 大于1公里
			logs.Warn("大于一公里")
			b = true
		}
	}
	// 重置距离.
	pl.SetStation(reqMsg.GetX(), reqMsg.GetY())
	resp := for_game.GetLocationInfoNew(who.GetPlayerId(), b, reqMsg)
	// 下行通知是否还有未读
	logs.Info("--------->", pl.GetIsNearBy())
	unReq := &client_hall.HasUnReadNearResp{
		HasUnReadNear: easygo.NewBool(pl.GetIsNearBy()),
	}
	ep.RpcHasUnReadNear(unReq)
	end := for_game.GetMillSecond()

	// 广告展示次数埋点
	easygo.Spawn(func() {
		for _, v := range resp.GetLocationInfo() {
			if v.GetDataType() != for_game.NEAR_INFO_DATE_TYPE_LEAD {
				continue
			}
			nearSet := v.GetNearSet()
			pId := who.GetPlayerId()
			opTime := time.Now().Unix()

			id := for_game.NextId(for_game.TABLE_NEARBY_ADV_LOG)
			id2 := for_game.NextId(for_game.TABLE_NEARBY_ADV_LOG)
			advLog := new(share_message.AdvLogReq)
			advLog.Id = easygo.NewInt64(id)
			advLog.AdvId = easygo.NewInt64(nearSet.GetId())
			advLog.PlayerId = easygo.NewInt64(pId)
			advLog.OpType = easygo.NewInt32(for_game.ADV_LOG_OP_TYPE_1)
			advLog.OpTime = easygo.NewInt64(opTime)
			for_game.AddNearAdvLogToDB(advLog)
			// 展示人数去重
			if log := for_game.GetNearAdvLogByPidAndOpFromDB(pId, nearSet.GetId(), for_game.ADV_LOG_OP_TYPE_2); log != nil {
				continue
			}
			advLog.Id = easygo.NewInt64(id2)
			advLog.OpType = easygo.NewInt32(for_game.ADV_LOG_OP_TYPE_2)
			for_game.AddNearAdvLogToDB(advLog)
		}
	})

	logs.Warn("附近的人执行时间为--------->%d 毫秒", end-begin)
	return resp
}

// 附近的人好友推荐
func (self *cls1) RpcNearRecommend(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.NearRecommendReq, common ...*base.Common) easygo.IMessage {
	logs.Info("===附近的人好友推荐 RpcNearRecommend=====,reqMsg=%v", reqMsg)
	x := reqMsg.GetX()
	y := reqMsg.GetY()
	if x == 0 || y == 0 {
		logs.Error("附近的人好友推荐,默认设置成,经度 x:%v,维度: %v", for_game.NEAR_DEFAULT_LNG, for_game.NEAR_DEFAULT_LAT) // 深圳
		reqMsg.X = easygo.NewFloat64(for_game.NEAR_DEFAULT_LNG)
		reqMsg.Y = easygo.NewFloat64(for_game.NEAR_DEFAULT_LAT)
	}
	bases, count := for_game.GetNearRecommend(who.GetPlayerId(), reqMsg.GetPage(), reqMsg.GetPageSize(), x, y)
	resp := &client_hall.NearRecommendResp{
		RecommendList: bases,
		Count:         easygo.NewInt32(count),
	}
	return resp
}

// 2.7.6 给附近的人打招呼
func (self *cls1) RpcNearSayMessage(ep IGameClientEndpoint, who *Player, reqMsg *share_message.NearSessionList, common ...*base.Common) easygo.IMessage {
	logs.Info("===2.7.6 给附近的人打招呼 RpcNearSayMessage=====,reqMsg=%v", reqMsg)
	pid := reqMsg.GetReceivePlayerId()
	if pid >= for_game.Min_Robot_PlayerId && pid <= for_game.Max_Robot_PlayerId {
		return nil
	}
	if reqMsg.GetContent() == "" {
		logs.Error("打招呼的内容为空")
		return easygo.NewFailMsg("打招呼的内容为空")
	}
	if reqMsg.GetClientUnique() == "" {
		logs.Error("clientUnique 为空")
		return easygo.NewFailMsg("参数为空")
	}
	contentType := reqMsg.GetContentType()
	if contentType != for_game.NEAR_CONTENT_TYPE_COMMENT && contentType != for_game.NEAR_CONTENT_TYPE_EMOTI { // 内容类型 1- 普通内容; 2-图片或表情
		logs.Error("打招呼的消息内有误, contentType: %d", contentType)
		return easygo.NewFailMsg("消息类型有误")
	}

	base := GetPlayerObj(pid)
	if base == nil {
		logs.Error("接收人不存在pid :", pid)
		return easygo.NewFailMsg("参数有误")
	}
	if util.Int64InSlice(who.GetPlayerId(), base.GetBlackList()) {
		res := "玩家拒绝接受你的消息"
		return easygo.NewFailMsg(res)
	}

	// 存入数据库
	contentId := for_game.AddSessionList(who.GetPlayerId(), pid, reqMsg)
	// 下行通知对方
	content := reqMsg.GetContent()
	notifyReq := &client_hall.NotifyNewSayMessageReq{
		Id:           easygo.NewInt64(contentId),
		SendPlayerId: easygo.NewInt64(reqMsg.GetSendPlayerId()),
		Content:      easygo.NewString(content),
		ClientUnique: easygo.NewString(reqMsg.GetClientUnique()),
		ContentType:  easygo.NewInt32(reqMsg.GetContentType()),
		PropsId:      easygo.NewInt64(reqMsg.GetPropsId()),
	}
	if PlayerOnlineMgr.CheckPlayerIsOnLine(pid) {
		ep1 := ClientEpMp.LoadEndpoint(pid)
		if ep1 == nil {
			return nil
		}
		ep1.RpcNotifyNewSayMessage(notifyReq)
	}
	if base := GetPlayerObj(pid); base != nil {
		base.SetIsNearBy(true)
	}
	ep.RpcNotifyNewSayMessage(notifyReq)
	return nil
}

// 2.7.6 附近的人打招呼会话列表
func (self *cls1) RpcNearSessionList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.NearSessionListReq, common ...*base.Common) easygo.IMessage {
	logs.Info("===2.7.6 附近的人打招呼会话列表  RpcNearSessionList=====,reqMsg=%v", reqMsg)
	resp := for_game.NearSessionList(who.GetPlayerId(), reqMsg)
	if resp.GetCount() == 0 {
		who.SetIsNearBy(false)
	}
	return resp
}

// 具体到某个人打招呼的列表
func (self *cls1) RpcSendPlayerMessageList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SendPlayerMessageListReq, common ...*base.Common) easygo.IMessage {
	logs.Info("=== 具体到某个人打招呼的列表  RpcSendPlayerMessageList=====,reqMsg=%v", reqMsg)
	if reqMsg.GetSendPlayerId() == 0 {
		return easygo.NewFailMsg("参数有误")
	}

	if reqMsg.GetPage() != for_game.DEFAULT_PAGE && reqMsg.GetMaxId() == 0 {
		return easygo.NewFailMsg("参数有误")
	}

	list, count := for_game.GetNearSayMessageList(who.GetPlayerId(), reqMsg)
	var maxId int64
	if reqMsg.GetPage() != for_game.DEFAULT_PAGE {
		maxId = reqMsg.GetMaxId()
	} else if len(list) > 0 {
		maxId = list[0].GetId()
	}
	resp := &client_hall.NearSessionListResp{
		SessionList: list,
		Count:       easygo.NewInt32(count),
		MaxId:       easygo.NewInt64(maxId),
	}

	return resp
}

//广播群特效
func (self *cls1) RpcBroadCastQTX(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.BroadCastQTX, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcBroadCastQTX:", reqMsg)
	if reqMsg.GetTeamId() == 0 {
		logs.Error("teamId 为空")
		return easygo.NewFailMsg("teamId 为空")
	}
	//群特效设置为播放了
	//memberObj := for_game.GetRedisTeamPersonalObj(reqMsg.GetTeamId())
	//if memberObj.GetTeamQTX(who.GetPlayerId()) {
	//	logs.Info("已经播放过特效了")
	//	return nil
	//}
	//memberObj.SetTeamQTX(who.GetPlayerId(), true)
	//广播给其他成员
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetTeamId())
	memberList := teamObj.GetTeamMemberList()
	BroadCastMsgToHallClientNew(memberList, "RpcBroadCastQTXResp", reqMsg)
	return nil
}

// 附近的人聊天内容列表
func (self *cls1) RpcGetNearChatList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.GetNearChatListReq, common ...*base.Common) easygo.IMessage {
	logs.Info("=== 附近的人聊天内容列表  RpcGetNearChatList=====,reqMsg=%v", reqMsg)
	if reqMsg.GetReceivePlayerId() == 0 {
		return easygo.NewFailMsg("参数有误")
	}
	//list, count := for_game.GetNearChatList(who.GetPlayerId(), reqMsg.GetReceivePlayerId(), reqMsg.GetPage(), reqMsg.GetPageSize())
	list, count := for_game.GetNearChatList(who.GetPlayerId(), reqMsg)
	var maxId int64
	if reqMsg.GetPage() != for_game.DEFAULT_PAGE {
		maxId = reqMsg.GetMaxId()
	} else if len(list) > 0 {
		maxId = list[0].GetId()
	}
	resp := &client_hall.NearSessionListResp{
		SessionList: list,
		Count:       easygo.NewInt32(count),
		MaxId:       easygo.NewInt64(maxId),
	}

	return resp
}

// 删除附近的人打招呼消息
func (self *cls1) RpcDelNearMessage(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.DelNearMessageReq, common ...*base.Common) easygo.IMessage {
	logs.Info("=== 删除附近的人打招呼消息  RpcDelNearMessage=====,reqMsg=%v", reqMsg)
	if len(reqMsg.GetSendPlayerId()) == 0 {
		return easygo.NewFailMsg("参数有误")
	}
	for_game.DelNearMessage(reqMsg.GetSendPlayerId(), who.GetPlayerId())
	return nil
}

// 修改聊天内容已读
func (self *cls1) RpcUpdateIsReadReq(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.UpdateIsReadReq, common ...*base.Common) easygo.IMessage {
	logs.Info("=== 修改聊天内容已读  RpcUpdateIsReadReq=====,reqMsg=%v", reqMsg)
	if reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("参数有误")
	}
	for_game.UpdateIsRead(reqMsg.GetId())
	return nil
}

//设置青少年模式
func (self *cls1) RpcSetYoungPassWord(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.YoungPassWord, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcSetYoungPassWord,reqMsg=%v", reqMsg, who.GetPlayerId())
	if reqMsg.GetOpt() == for_game.YOUNG_STATUS_OPEN && reqMsg.GetPassWord() == "" {
		return easygo.NewFailMsg("密码不能为空")
	}
	if reqMsg.GetOpt() == for_game.YOUNG_STATUS_CLOSE && reqMsg.GetPassWord() != who.GetYoungPassWord() {
		return easygo.NewFailMsg("输入密码错误")
	}
	newPassWord := reqMsg.GetPassWord()
	if reqMsg.GetOpt() == for_game.YOUNG_STATUS_CLOSE {
		newPassWord = ""
	}
	who.SetYoungPassWord(newPassWord)
	return nil
}

// 附近的人添加广告埋点数据
func (self *cls1) RpcNearAddAdvLog(ep IGameClientEndpoint, who *Player, reqMsg *share_message.AdvLogReq, common ...*base.Common) easygo.IMessage {
	logs.Info("================添加附近的人广告埋点数据================", reqMsg)
	opType := reqMsg.GetOpType()
	if opType != for_game.ADV_LOG_OP_TYPE_1 && opType != for_game.ADV_LOG_OP_TYPE_3 {
		logs.Error("附近的人大厅广告埋点数据,操作类型有误,opType: ", opType)
		return easygo.NewFailMsg("操作类型有误")
	}
	reqMsg.PlayerId = easygo.NewInt64(who.GetPlayerId())
	reqMsg.OpTime = easygo.NewInt64(time.Now().Unix())

	// 前端只传1-展示,3-点击
	switch opType {
	case for_game.ADV_LOG_OP_TYPE_1: // 展示次数+1,判断是否需要去重.
		id := for_game.NextId(for_game.TABLE_NEARBY_ADV_LOG)
		id2 := for_game.NextId(for_game.TABLE_NEARBY_ADV_LOG)
		reqMsg.Id = easygo.NewInt64(id)
		for_game.AddNearAdvLogToDB(reqMsg)
		// 展示人数去重
		if log := for_game.GetNearAdvLogByPidAndOpFromDB(who.GetPlayerId(), reqMsg.GetAdvId(), for_game.ADV_LOG_OP_TYPE_2); log != nil {
			return nil
		}
		reqMsg.Id = easygo.NewInt64(id2)
		reqMsg.OpType = easygo.NewInt32(for_game.ADV_LOG_OP_TYPE_2)
		for_game.AddNearAdvLogToDB(reqMsg)
		return nil
	case for_game.ADV_LOG_OP_TYPE_3:
		id := for_game.NextId(for_game.TABLE_NEARBY_ADV_LOG)
		id2 := for_game.NextId(for_game.TABLE_NEARBY_ADV_LOG)
		reqMsg.Id = easygo.NewInt64(id)
		for_game.AddNearAdvLogToDB(reqMsg)
		// 点击人数去重
		if log := for_game.GetNearAdvLogByPidAndOpFromDB(who.GetPlayerId(), reqMsg.GetAdvId(), for_game.ADV_LOG_OP_TYPE_4); log != nil {
			return nil
		}
		reqMsg.Id = easygo.NewInt64(id2)
		reqMsg.OpType = easygo.NewInt32(for_game.ADV_LOG_OP_TYPE_4)
		for_game.AddNearAdvLogToDB(reqMsg)
		return nil
	}
	return nil
}

//玩家登录时，获取个人会话数据
func (self *cls1) RpcGetSessionData(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.AllSessionData, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetSessionData:", reqMsg, who.GetPlayerId())
	sessionDatas := make([]*client_hall.SessionData, 0)
	ids := who.GetPlayerSessions()
	if len(ids) == 0 {
		return reqMsg
	}
	data := for_game.GetAllSessionData(ids)
	if len(data) == 0 {
		return reqMsg
	}
	notExistTeam := make([]string, 0)
	pid := who.GetPlayerId()
	for _, s := range data {
		session := for_game.GetRedisChatSessionObj(s.GetId(), s)
		//撤回消息处理
		withDrawLogIds := make([]int64, 0)
		if session.GetType() == for_game.CHAT_TYPE_PRIVATE {
			withDrawLogIds = for_game.GetAllWithDrawLogIds(s.GetId(), pid, session.GetPersonalReadId(pid))
		} else {
			teamId := easygo.AtoInt64(s.GetId())
			chatObj := for_game.GetRedisTeamChatLog(teamId)
			memberObj := for_game.GetRedisTeamPersonalObj(teamId)
			readId := int64(0)
			if memberObj != nil {
				readId = memberObj.GetTeamReadChatLogId(pid)
			}
			withDrawLogIds = chatObj.GetAllTeamWithDrawLogIds(pid, readId)
		}
		if who.GetIsLoadedAllSessions() {
			if reqMsg.GetType() == 1 && !session.CheckSessionForClient(pid) {
				if len(withDrawLogIds) == 0 {
					continue
				}
			}
		}
		sessionData := session.GetChatSessionDataForClient(who.GetPlayerId(), true)
		if sessionData != nil {
			sessionData.WithdrawList = withDrawLogIds
			sessionDatas = append(sessionDatas, sessionData)
		} else {
			notExistTeam = append(notExistTeam, s.GetId())
		}
	}
	if len(notExistTeam) > 0 {
		who.DeletePlayerSessions(notExistTeam)
	}
	if !who.GetIsLoadedAllSessions() {
		who.SetIsLoadedAllSessions(true)
	}
	reqMsg.Sessions = sessionDatas
	//logs.Info("RpcGetSessionData 返回:", len(reqMsg.GetSessions()), reqMsg)
	return reqMsg
}

//获取指定会话数据
func (self *cls1) RpcGetOneSessionData(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SessionData, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetOneSessionData", reqMsg)
	session := for_game.GetRedisChatSessionObj(reqMsg.GetId())
	if session == nil {
		return reqMsg
	}
	sessionData := session.GetChatSessionDataForClient(who.GetPlayerId(), true)
	if sessionData == nil {
		return reqMsg
	}
	return sessionData
}

//获取指定会话的详细数据
func (self *cls1) RpcGetSessionDetail(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SessionData, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetSessionDetail", reqMsg, who.GetPlayerId())
	session := for_game.GetRedisChatSessionObj(reqMsg.GetId())
	sessionData := session.GetChatSessionDataForClient(who.GetPlayerId(), true)
	return sessionData
}

//获取指定会话聊天内容
func (self *cls1) RpcGetSessionChat(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SessionChatData, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetSessionChat:", reqMsg, who.GetPlayerId())
	session := for_game.GetRedisChatSessionObj(reqMsg.GetSessionId())
	//playerId := who.GetPlayerId()
	if session == nil {
		logs.Error("找不到该会话:id=", reqMsg.GetSessionId())
		return easygo.NewFailMsg("错误的会话id")
	}
	if session.GetType() == for_game.CHAT_TYPE_PRIVATE {
		//私聊
		//全部已读了，无需查询数据
		if reqMsg.GetStartId() == session.GetMaxLogId() {
			return reqMsg
		}
		log := for_game.GetPersonalSessionChatLog(reqMsg.GetStartId(), reqMsg.GetEndId(), reqMsg.GetSessionId())
		reqMsg.PersonalChatLogs = log
		logs.Info("返回数据:", len(log), log)

	} else if session.GetType() == for_game.CHAT_TYPE_TEAM {
		//群聊
		teamId := easygo.AtoInt64(reqMsg.GetSessionId())
		//全部已读了，无需查询数据
		if reqMsg.GetStartId() == session.GetMaxLogId() {
			return reqMsg
		}
		chatObj := for_game.GetRedisTeamChatLog(teamId)
		log := chatObj.GetTeamSessionChatLogs(reqMsg.GetStartId(), reqMsg.GetEndId(), reqMsg.GetSessionId())
		reqMsg.TeamChatLogs = log
	}
	//logs.Info("RpcGetSessionChat 返回:", reqMsg)
	return reqMsg
}

//群@数据,
func (self *cls1) RpcGetTeamAtData(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.TeamAtData, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetTeamAtData:", reqMsg)
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetTeamId())
	if teamObj == nil {
		panic("无效的群id")
	}
	members := for_game.GetTeamAtMember(who.GetPlayerId(), reqMsg.GetTeamId(), reqMsg.GetPage())
	reqMsg.Data = members
	return reqMsg
}

//群成员数据请求，每次30条
func (self *cls1) RpcGetTeamMemberData(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.TeamMemberData, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetTeamMemberData:", reqMsg)
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetTeamId())
	if teamObj == nil {
		panic("无效的群id")
	}
	members := for_game.GetTeamMemberDatas(who.GetPlayerId(), reqMsg.GetTeamId(), reqMsg.GetPage())
	ids := make([]int64, 0)
	for _, m := range members {
		ids = append(ids, m.GetPlayerId())
	}
	pMap := for_game.GetAllPlayerBase(ids, false)
	for _, m := range members {
		p := pMap[m.GetPlayerId()]
		if p != nil {
			m.HeadIcon = easygo.NewString(p.GetHeadIcon())
			if m.GetNickName() == "" {
				m.NickName = easygo.NewString(p.GetNickName())
			}
		}
	}
	reqMsg.Data = members
	return reqMsg
}

//获取群详细信息
func (self *cls1) RpcGetTeamDetailData(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.TeamDetailData, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetTeamDetailData:", reqMsg)
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetTeamId())
	if teamObj == nil {
		panic("无效的群id")
	}
	team := teamObj.GetRedisTeam()
	owner := for_game.GetRedisPlayerBase(teamObj.GetTeamOwner())
	if owner != nil {
		team.OwnerAccount = easygo.NewString(owner.GetAccount())
		team.OwnerNickName = easygo.NewString(owner.GetNickName())
	}
	reqMsg.Team = team

	return reqMsg
}

//检测玩家是否是指定群成员
func (self *cls1) RpcCheckIsTeamMember(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CheckTeamMember, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcCheckIsTeamMember:", reqMsg)
	teamIds := who.GetTeamIds()
	if easygo.Contain(teamIds, reqMsg.GetTeamId()) {
		reqMsg.IsMember = easygo.NewBool(true)
	} else {
		reqMsg.IsMember = easygo.NewBool(false)
	}
	session := for_game.GetRedisChatSessionObj(easygo.AnytoA(reqMsg.GetTeamId()))
	if session != nil {
		reqMsg.Session = session.GetChatSessionDataForTeam(who.GetPlayerId(), false)
	}
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetTeamId())
	if teamObj != nil {
		owner := for_game.GetRedisPlayerBase(teamObj.GetTeamOwner())
		if owner != nil {
			at := &client_hall.AtData{
				PlayerId: easygo.NewInt64(owner.GetPlayerId()),
				Name:     easygo.NewString(owner.GetNickName()),
				HeadUrl:  easygo.NewString(owner.GetHeadIcon()),
				Sex:      easygo.NewInt32(owner.GetSex()),
			}
			reqMsg.Owner = at
		}
	}
	return reqMsg
}

//获取玩家保存的群会话列表
func (self *cls1) RpcGetSaveTeamSessions(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetSaveTeamSessions:", reqMsg)
	sessionDatas := make([]*client_hall.SessionData, 0)
	ids := for_game.GetMySaveTeamIds(who.GetPlayerId())
	if len(ids) == 0 {
		return reqMsg
	}
	data := for_game.GetAllSessionData(ids)
	if len(data) == 0 {
		return reqMsg
	}
	for _, s := range data {
		session := for_game.GetRedisChatSessionObj(s.GetId(), s)
		sessionData := session.GetChatSessionDataForClient(who.GetPlayerId(), true)
		if sessionData != nil {
			sessionDatas = append(sessionDatas, sessionData)
		}
	}
	resp := &client_hall.AllSessionData{
		Sessions: sessionDatas,
	}
	return resp
}
func (self *cls1) RpcDeleteMessage(ep IGameClientEndpoint, who *Player, reqMsg *client_server.ReadInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcDeleteMessage:", reqMsg)
	session := for_game.GetRedisChatSessionObj(reqMsg.GetSessionId())
	if session == nil {
		logs.Error("无效的会话id:", reqMsg.GetSessionId())
		return nil
	}
	session.DeleteMessage(reqMsg.GetLogId())
	return nil
}
func (self *cls1) RpcCheckIsMyTeamSession(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.CheckIsMySession, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcCheckIsMyTeamSession:", reqMsg)
	session := for_game.GetRedisChatSessionObj(reqMsg.GetSessionId())
	reqMsg.Result = easygo.NewBool(false)
	if session == nil {
		return reqMsg
	}
	players := session.GetRedisChatSessionPlayers()
	if !easygo.Contain(players, reqMsg.GetPlayerId()) {
		return reqMsg
	}
	reqMsg.Result = easygo.NewBool(true)
	return reqMsg
}

//获取自己能显示的最近一条消息
func (self *cls1) RpcGetOneShowLog(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.GetOneShowLog, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetOneShowLog:", reqMsg)
	session := for_game.GetRedisChatSessionObj(reqMsg.GetSessionId())
	if session == nil {
		return reqMsg
	}
	log := session.GetOneShowLog(who.GetPlayerId(), reqMsg.GetLogId(), reqMsg.GetReadId())
	reqMsg.Log = log
	return reqMsg
}

//上报按钮点击数据
func (self *cls1) RpcBtnClick(ep IGameClientEndpoint, who *Player, reqMsg *client_server.BtnClickInfo, common ...*base.Common) easygo.IMessage {
	// logs.Info("RpcBtnClick,pid----->", reqMsg.GetBtnType(), who.GetPlayerId())
	player := for_game.GetRedisPlayerBase(who.GetPlayerId())
	if player == nil {
		return easygo.NewFailMsg("无效的玩家数据")
	}
	if reqMsg.BtnType == nil || reqMsg.GetBtnType() == 0 {
		return easygo.NewFailMsg("上报类型错误")
	}

	timeNow := easygo.GetToday0ClockTimestamp()
	playerid := who.GetPlayerId()
	btnType := reqMsg.GetBtnType()
	log := &share_message.ButtonClickLog{
		Id:         easygo.NewString(easygo.AnytoA(timeNow) + easygo.AnytoA(playerid) + easygo.AnytoA(btnType)),
		CreateTime: easygo.NewInt64(timeNow),
		PlayerId:   easygo.NewInt64(playerid),
		Type:       easygo.NewInt32(btnType),
	}
	err := for_game.InsertMgo(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_BUTTON_CLICK_LOG, log) //如果插入失败,说明数据已经存在.

	switch btnType {
	case for_game.InterestSurePV:
		for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "InterestOk")
		if err == nil {
			for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "InterestOkCount")
		}
	case for_game.InterestReturnPV:
		for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "InterestBack")
		if err == nil {
			for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "InterestBackCount")
		}
	case for_game.RecommendSkipPV:
		for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "RecommendSkip")
		if err == nil {
			for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "RecommendSkipCount")
		}
	case for_game.RecommendNextPV:
		for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "RecommendNext")
		if err == nil {
			for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "RecommendNextCount")
		}
	case for_game.InNmBtn:
		for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "InNmBtnCLick")
		if err == nil {
			for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "InNmBtnCount")
		}
	case for_game.InfoBack:
		for_game.SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, "InfoBack")
	default:
		return easygo.NewFailMsg("上报类型错误")
	}

	return nil
}

//获取新的群成员变化
func (self *cls1) RpcGetTeamMemberChange(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.TeamMembersChange, common ...*base.Common) easygo.IMessage {
	// logs.Info("RpcGetTeamMemberChange", reqMsg)
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetTeamId())
	if teamObj == nil {
		s := fmt.Sprintf("无效的群id=%d", reqMsg.GetTeamId())
		return easygo.NewFailMsg(s)
	}
	chatObj := for_game.GetRedisTeamChatLog(reqMsg.GetTeamId())
	if chatObj != nil {
		chatObj.SaveToMongo()
	}
	memberObj := for_game.GetRedisTeamPersonalObj(reqMsg.GetTeamId())
	if memberObj == nil {
		s := fmt.Sprintf("无效的群id=%d", reqMsg.GetTeamId())
		return easygo.NewFailMsg(s)
	}
	member := memberObj.GetTeamMember(who.GetPlayerId())
	if member == nil {
		s := fmt.Sprintf("无效的玩家id=%d", who.GetPlayerId())
		return easygo.NewFailMsg(s)
	}
	t := reqMsg.GetTime()
	delList := for_game.GetTeamLogDelPlayers(reqMsg.GetTeamId(), t)
	reqMsg.DelList = delList
	addList := for_game.GetTeamAddPlayers(reqMsg.GetTeamId(), t)
	ids := make([]int64, 0)
	for _, m := range addList {
		ids = append(ids, m.GetPlayerId())
	}
	mPlayers := for_game.GetAllPlayerBase(ids)
	for _, m := range addList {
		for_game.GetOneTeamMember(m, mPlayers[m.GetPlayerId()])
	}
	reqMsg.AddList = addList
	reqMsg.PosChange = for_game.GetTeamPlayerManagerChange(reqMsg.GetTeamId(), t)
	logs.Info("RpcGetTeamMemberChange返回信息:", reqMsg)
	return reqMsg
}

//上报注册登录页面加载数据
func (self *cls1) RpcPageRegLogLoad(ep IGameClientEndpoint, who *Player, reqMsg *client_server.PageRegLogLoad, common ...*base.Common) easygo.IMessage {
	// logs.Info("RpcPageRegLogLoad", reqMsg)

	if reqMsg.Type == nil || reqMsg.GetType() == 0 {
		// return easygo.NewFailMsg("上报类型错误")
		logs.Error("上报类型错误" + easygo.AnytoA(reqMsg.GetType()))
		return nil
	}
	if reqMsg.Code == nil {
		// return easygo.NewFailMsg("设备码不能为空")
		logs.Error("上报注册登录页面加载数据,设备码不能为空")
		return nil
	}

	easygo.Spawn(func() {
		data := for_game.FindOne(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_POS_DEVICECODE, bson.M{"DeviceCode": reqMsg.GetCode()})
		if data == nil {
			logs.Error("上报注册登录页面加载数据,设备码错误")
		} else {
			dicType := 2 //1-新设备,2-旧设备
			one := &share_message.PosDeviceCode{}
			for_game.StructToOtherStruct(data, one)
			if one.GetCreateTime() > easygo.Get0ClockMillTimestamp(easygo.NowTimestamp()) {
				dicType = 1
			}

			timeNow := util.GetMilliTime()
			log := &share_message.PageRegLog{
				Id:         easygo.NewString(easygo.AnytoA(timeNow) + reqMsg.GetCode()),
				CreateTime: easygo.NewInt64(timeNow),
				Code:       easygo.NewString(reqMsg.GetCode()),
				Type:       easygo.NewInt32(reqMsg.GetType()),
				DicType:    easygo.NewInt32(dicType),
				Channel:    easygo.NewString(reqMsg.GetChannel()),
			}
			err := for_game.InsertMgo(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PAGE_REGLOG, log) //如果插入失败,说明数据已经存在.
			if err != nil {
				logs.Error("设备码已经存在:", err.Error())
			}
		}
	})

	return nil
}

//获取支持的银行列表
func (self *cls1) RpcGetSupportBankList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SupportBankList, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetSupportBankList", reqMsg)
	switch reqMsg.GetPayId() {
	case for_game.PAY_CHANNEL_HUICHAO_WX, for_game.PAY_CHANNEL_HUICHAO_ZFB, for_game.PAY_CHANNEL_HUICHAO_YL, for_game.PAY_CHANNEL_HUICHAO_DF:
		banks := PWebHuiChaoPay.GetSupportBankList()
		reqMsg.Banks = banks
	case for_game.PAY_CHANNEL_PENGJU:
		banks := PWebPengJuPay.GetSupportBankList()
		reqMsg.Banks = banks
	}
	//logs.Info("RpcGetSupportBankList 返回:", reqMsg)
	return reqMsg
}

//获取在线人数
func (self *cls1) RpcGetOnLineNum(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetOnLineNum:", who.GetPlayerId())
	num := GetPlayerOnLineNum()
	resp := &client_hall.OnLineNum{
		Num: easygo.NewInt64(num),
	}
	return resp
}

//获取消息也广告信息
//func (self *cls1) RpcGetMsgPageAdvList(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
//	logs.Info("RpcGetMsgPageAdvList:", who.GetPlayerId())
//	advs := for_game.QueryAdvListToDB(for_game.ADV_LOCATION_BANNER_MSG, who.GetLastLoginIP())
//	resp := &client_hall.MsgAdv{
//		Advs: advs,
//	}
//	return resp
//}
//获取消息也广告信息
func (self *cls1) RpcGetMsgPageAdvList(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetMsgPageAdvList:", who.GetPlayerId())
	left, right := for_game.QueryMsgPageAdvToDB(who.GetLastLoginIP())
	resp := &client_hall.MsgAdv{
		Advs:      left,
		RightList: right,
	}
	return resp
}

//获取玩家钻石数
func (self *cls1) RpcGetDiamond(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetDiamond:", who.GetPlayerId())
	diamond := GetDiamondFromWishServer(who.GetPlayerId())
	resp := &client_hall.PlayerDiamond{
		PlayerId: easygo.NewInt64(who.GetPlayerId()),
		Diamond:  easygo.NewInt64(diamond),
	}
	return resp
}

//获取支持得渠道信息
func (self *cls1) RpcGetSupportPayChannel(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PayChannels, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetSupportPayChannel:", reqMsg, who.GetPlayerId())
	channel := PPayChannelMgr.GetPlatformChannelList(reqMsg.GetType())
	reqMsg.Channels = channel
	logs.Info("响应 RpcGetSupportPayChannel:", reqMsg)
	return reqMsg
}

//设置用户兴趣标签
func (self *cls1) RpcSetPlayerLabel(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.FirstInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcSetPlayerLabel: %v ,reqMsg: %v", reqMsg, who.GetPlayerId())
	label := reqMsg.GetLabel()
	if len(label) < 1 {
		easygo.NewFailMsg("参数有误")
	}

	easygo.Spawn(func() {
		if len(who.GetLabelList()) == 0 {
			for_game.SetRedisRegisterLoginReportFildVal(easygo.Get0ClockTimestamp(util.GetMilliTime()), 1, "LabelCount") //更新报表设置标签的人数
		}
	})

	who.SetRedisLabelList(label)
	var photo string
	var playTime int32

	m := for_game.QueryInterestGroupByGroup(label)
	if m == nil {
		rand.Seed(time.Now().Unix())
		index := rand.Intn(len(label))
		id := label[index]
		msg := for_game.GetInterestTag(id)
		photo = msg.GetPopIcon()
		playTime = msg.GetPlayTime()
	} else {
		photo = m.GetPopIcon()
		playTime = m.GetPlayTime()
	}
	newMsg := for_game.GetRecommendInfo(who.GetPlayerId(), 0)
	newMsg.PlayTime = easygo.NewInt32(playTime)
	newMsg.Photo = easygo.NewString(photo)
	ep.RpcReturnRecommendInfo(newMsg) //发送推荐的好友群信息
	ep.RpcPlayerAttrChange(who.GetPlayerInfo())

	return &base.Empty{}
}

//获取所有标签信息
func (self *cls1) RpcGetAllLabelMsg(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetAllLabelMsg,", reqMsg)
	msgList := for_game.GetInterestTagAllList()
	sysInfo := for_game.QuerySysParameterById(for_game.INTEREST_PARAMETER)
	msg := &client_server.LabelMsg{
		LabelInfo:    msgList,
		Max:          easygo.NewInt32(sysInfo.GetInterestMax()),
		Min:          easygo.NewInt32(sysInfo.GetInterestMin()),
		InterestType: for_game.GetRedisLabelInfo(),
	}
	return msg
}

// 话题群里面群动态api接口
func (self *cls1) RpcTopicTeamDynamicList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.TopicTeamDynamicReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcTopicTeamDynamicList:%v,reqMsg: %v", who.GetPlayerId(), reqMsg)
	// 判断话题群是否存在
	teamId := reqMsg.GetTopicTeamId()
	if teamId == 0 {
		logs.Error("RpcTopicTeamDynamicList 群id为0")
		return easygo.NewFailMsg("参数有误")
	}

	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		s := fmt.Sprintf("无效的群id=%d", teamId)
		return easygo.NewFailMsg(s)
	}
	resp := &client_hall.TopicTeamDynamicResp{}
	memberList := teamObj.GetTeamMemberList()
	if len(memberList) == 0 {
		return resp
	}

	page, pageSize := for_game.MakePageAndPageSize(reqMsg.GetPage(), reqMsg.GetPageSize())
	DynamicData, count := for_game.GetDynamicListByPids(page, pageSize, memberList)

	pid := who.GetPlayerId()
	player := for_game.GetRedisPlayerBase(pid)
	attentionList := player.GetAttention()
	ds := for_game.GetRedisSomeDynamic1(pid, DynamicData, attentionList) // 遍历动态,判断是否已关注.

	resp.TopicTeamId = easygo.NewInt64(teamId)
	resp.DynamicList = ds
	resp.Count = easygo.NewInt32(count)
	return resp
}

//获取话题群组分组信息
func (self *cls1) RpcGetTopicTeams(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.TopicTeams, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetTopicTeams,", reqMsg)
	teams := for_game.GetTeamsByTopic(reqMsg.GetTopic())
	for i := 0; i < len(teams); i++ {
		obj := for_game.GetRedisTeamObj(teams[i].GetId(), teams[i])
		if obj != nil {
			teams[i] = obj.GetRedisTeam()
			headUrls := make([]string, 0)
			membersList := teams[i].GetMemberList()
			if len(membersList) > 6 {
				membersList = membersList[:6]
			}
			for _, pid := range membersList {
				player := for_game.GetRedisPlayerBase(pid)
				if player != nil {
					headUrls = append(headUrls, player.GetHeadIcon())
				}
			}
			teams[i].TopicHeadUrls = headUrls
		}
	}
	reqMsg.Teams = teams
	return reqMsg
}

//获取主页菜单项信息
func (self *cls1) RpcGetAllMainMenu(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetAllMainMenu,", reqMsg)
	msg := for_game.GetIndexTipsToClient()
	logs.Info("返回 RpcGetAllMainMenu:", msg)
	return msg
}

//获取所有弹窗广告
func (self *cls1) RpcGetAllTipAdvs(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetAllTipAdvs,", reqMsg)
	msg := for_game.GetAllTipAdvsFromDB()
	logs.Info("相应 RpcGetAllTipAdvs:", msg)
	return msg
}

//获取启动页广告信息
func (self *cls1) RpcGetStartPageAdvList(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetStartPageAdvList:", who.GetPlayerId())
	adv := for_game.QueryRandAdvListToDB(for_game.ADV_LOCATION_START, 1, who.GetLastLoginIP())
	resp := &client_hall.MsgAdv{
		Advs: adv,
	}
	return resp
}
