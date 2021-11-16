// 大厅发给子游戏服的消息处理

package hall

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"reflect"
	"strconv"
	"time"

	"github.com/akqp2019/mgo/bson"

	"github.com/astaxie/beego/logs"
)

//===================================================================

type ServiceForBackStage struct {
	Service reflect.Value
}

//分发消息
func (self *ServiceForBackStage) RpcServerReport(common *base.Common, reqMsg *share_message.ServerInfo) easygo.IMessage {
	logs.Info("RpcServerReport:", reqMsg)
	//ep.SetServerId(reqMsg.GetSid())
	//BackStageEpMgr.Store(reqMsg.GetSid(), ep)
	PServerInfoMgr.AddServerInfo(reqMsg)
	return nil
}

//后台顶号处理
//func (self *ServiceForBackStage) RpcReplaceLoginToHall(common *base.Common, reqMsg *server_server.AdminInfo) easygo.IMessage {
//	logs.Info("RpcReplaceLoginToHall:", reqMsg)
//	msg := &server_server.PlayerSI{
//		PlayerId: easygo.NewInt64(reqMsg.GetUserId()),
//	}
//	SendMsgToServerNew(reqMsg.GetServerId(), "RpcReplaceLogin", msg)
//	return nil
//}

//玩家数据发生变化
/*
func (self *ServiceForBackStage) RpcPlayerChangeHall(common *base.Common, reqMsg *backstage_hall.PlayerSI) easygo.IMessage {
	logs.Info("RpcPlayerChange ", reqMsg)
	player := PlayerMgr.LoadPlayer(reqMsg.GetPlayerId())
	if player != nil {
		player.OnLoadFromDB()
	}

	return nil
}*/

//后台请求修改群资料
func (self *ServiceForBackStage) RpcTeamChangeHall(common *base.Common, reqMsg *server_server.EditTeam) easygo.IMessage {
	logs.Info("RpcTeamChange ", reqMsg)
	teamId := reqMsg.GetId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		logs.Info("群不存在，id:", teamId)
		return nil
	}
	if reqMsg.IsRecommend != nil {
		teamObj.SetTeamIsRecommend(reqMsg.GetIsRecommend())
	}
	if reqMsg.MaxMember != nil && reqMsg.GetMaxMember() != 0 {
		teamObj.SetTeamMaxMember(reqMsg.GetMaxMember())
	}
	if reqMsg.GongGao != nil && reqMsg.GetGongGao() != "" {
		teamObj.SetTeamGongGao(reqMsg.GetGongGao())
	}
	if reqMsg.Name != nil && reqMsg.GetName() != "" {
		teamObj.SetTeamNikeName(reqMsg.GetName())
	}
	if reqMsg.Level != nil && reqMsg.GetLevel() != 0 {
		teamObj.SetTeamLevel(reqMsg.GetLevel())
	}

	return nil
}

//后台请求群解散
func (self *ServiceForBackStage) RpcDefunctTeamHall(common *base.Common, reqMsg *server_server.PlayerSI) easygo.IMessage {
	logs.Info("RpcDefunctTeam ", reqMsg)
	teamId := reqMsg.GetPlayerId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		logs.Info("群不存在，id:", teamId)
		return nil
	}

	info := &share_message.OperatorInfo{
		Operator: easygo.NewString(reqMsg.GetAccount()),
		Time:     easygo.NewInt64(util.GetMilliTime()),
		Flag:     easygo.NewInt32(1),
	}

	teamObj.SetRedisOperatorInfo(info)
	teamObj.SetTeamStatus(for_game.DISSOLVE)
	teamObj.SetTeamDissolveTime(time.Now().Unix())
	memberList := teamObj.GetTeamMemberList()
	DeleteTeamMemberOperate(teamId, 0, memberList)
	return nil
}

//后台请求增减群成员
func (self *ServiceForBackStage) RpcTeamMemberOptHall(common *base.Common, reqMsg *server_server.MemberOptRequest) easygo.IMessage {
	logs.Info("RpcTeamMemberOptToHall ", reqMsg)
	plst := reqMsg.GetPlayerIds()
	teamId := reqMsg.GetTeamId()
	teamObj := for_game.GetRedisTeamObj(teamId)
	if teamObj == nil {
		logs.Info("找不到这个群信息:", teamId)
		return nil
	}
	playerId := reqMsg.GetPlayerID()
	player := for_game.GetRedisPlayerBase(playerId)
	if player == nil {
		panic("玩家对象为空")
	}

	switch reqMsg.GetTypes() {
	case 1: //增加成员
		reason := AddTeamMemberOperate(teamId, playerId, reqMsg.GetAdminID(), plst, player.GetNickName())
		if reason != "" {
			return easygo.NewFailMsg(reason)
		}
	case 2: //删除成员
		reason := DeleteTeamMemberOperate(teamId, playerId, plst)
		if reason != "" {
			return easygo.NewFailMsg(reason)
		}
	default:
		return easygo.NewFailMsg("操作类型错误")
	}
	//后台请求直接保存成员数据
	memberObj := for_game.GetRedisTeamPersonalObj(teamId)
	memberObj.SaveToMongo()
	teamObj.SaveToMongo()
	return nil
}

//后台出入款通知
func (self *ServiceForBackStage) RpcRechargeToHall(common *base.Common, reqMsg *server_server.Recharge) easygo.IMessage {
	logs.Info("RpcRechargeToHall ", reqMsg)
	RechargeGoldToPlayer(reqMsg.GetPlayerId(), reqMsg.GetOrderId())
	return nil
}

//后台通知大厅硬币变化
func (self *ServiceForBackStage) RpcChangeCoinsToHall(common *base.Common, reqMsg *server_server.Recharge) easygo.IMessage {
	logs.Info("RpcChangeCoinsToHall ", reqMsg)

	if err := NotifyAddCoin(reqMsg.GetPlayerId(), reqMsg.GetRechargeGold(), reqMsg.GetNote(), reqMsg.GetSourceType(), nil); err != "" {
		return easygo.NewFailMsg(err)
	}
	return nil
}

func (self *ServiceForBackStage) RpcGetPlayerBase(common *base.Common, reqMsg *server_server.PlayerSI) easygo.IMessage {
	return for_game.GetPlayerById(reqMsg.GetPlayerId())
}

//后台请求给APP推送消息
func (self *ServiceForBackStage) RpcSendAppPushHall(common *base.Common, reqMsg *share_message.AppPushMessage) easygo.IMessage {
	logs.Info("RpcSendAppPushHall ", reqMsg)
	recipient := reqMsg.GetRecipient()
	lis := for_game.QueryPlayersByOfPush(recipient, reqMsg.GetLabel(), reqMsg.GetCustomTag(), reqMsg.GetGrabTag())
	ids := make([]string, 0)
	for _, i := range lis {
		if i.GetToken() != "" {
			ids = append(ids, i.GetToken())
		}
	}

	m := for_game.PushMessage{
		Title:       reqMsg.GetTitle(),
		Content:     reqMsg.GetContent(),
		ContentType: for_game.JG_TYPE_BACKSTAGE,
		JumpObject:  reqMsg.GetJumpObject(),
		ObjectId:    reqMsg.GetObjectId(),
		TargetId:    easygo.AnytoA(reqMsg.GetId()),
	}

	if reqMsg.GetJumpObject() == 9 {
		item := for_game.QueryShopItemById(reqMsg.GetObjectId())
		if item != nil {
			m.PlayerId = item.GetPlayerId()
		}
	}
	bytes, _ := json.Marshal(m)
	logs.Info("开始调用for_game的推送方法, pushMessage:--------> %s", string(bytes))
	for_game.JGSendMessage(ids, m)

	for_game.SetRedisNoticeReportFildVal(reqMsg.GetId(), int64(len(ids)), "PushPlayer") //添加文章报表推送用户数
	return nil
}

//后台回复投诉意见反馈
func (self *ServiceForBackStage) RpcResponeComplainInfo(common *base.Common, reqMsg *share_message.PlayerComplaint) easygo.IMessage {
	//1：意见反馈  2:投诉  3：订单投诉  4.商城物品投诉  5群投诉
	msg := "意见"
	switch reqMsg.GetType() {
	case 1:
		msg = "意见反馈"
	case 2:
		msg = "投诉"
	case 3:
		msg = "订单投诉"
	case 4:
		msg = "商城物品投诉"
	case 5:
		msg = "群投诉"
	}
	pid := reqMsg.GetPlayerId()
	title := fmt.Sprintf("%s回复", msg)
	NoticeAssistant(pid, 2, title, reqMsg.GetReContent())

	if !PlayerOnlineMgr.CheckPlayerIsOnLine(pid) || PlayerOnlineMgr.CheckPlayerIsCutBackstage(pid) {
		player := for_game.GetRedisPlayerBase(pid)
		isNotice := player.GetIsNewMessage()
		if isNotice {
			f := func() {
				isShow := player.GetIsMessageShow()
				ids := for_game.GetJGIds([]int64{pid})
				var content string
				if isShow {
					content = reqMsg.GetReContent()
				} else {
					content = "给您发送了一条消息"
				}
				m := for_game.PushMessage{
					Title:       "投诉反馈",
					Content:     content,
					ContentType: for_game.JG_TYPE_BACKSTAGE,
					JumpObject:  2,
				}
				for_game.JGSendMessage(ids, m)
			}
			easygo.Spawn(f)
		}
	}
	return nil
}

//后台修改支付配置通知到大厅重载
func (self *ServiceForBackStage) RpcPaySetChangeToHall(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	//查询当前可用的支付通道配置
	PPayChannelMgr.Init()
	//通知前端修改
	//msg := &client_server.AllPlayerMsg{
	//	Pay:       PPayChannelMgr.GetPlatformChannelList(),
	//	PayConfig: PPayChannelMgr.GetPaymentSettingList(),
	//}
	//BroadCastMsgToOnlineClient("RpcPayChannelChange", msg)
	return nil
}

//心跳
func (self *ServiceForBackStage) RpcHeartBeat(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	//logs.Info("RpcHeartBeat ")
	return nil
}
func (self *ServiceForBackStage) RpcEditSystemNoticeMessage(common *base.Common, reqMsg *share_message.SystemNotice) easygo.IMessage {
	end_points := ClientEpMgr.GetEndpoints()
	for _, ep := range end_points {
		player := ep.GetPlayer()
		if player != nil && player.GetIsRobot() {
			continue
		}
		if reqMsg.GetUserType() == 0 || player.GetDeviceType() == reqMsg.GetUserType() {
			NoticeAssistant3(player.GetPlayerId(), reqMsg.GetId(), reqMsg.GetTitle(), reqMsg.GetContent())
		}
	}

	return nil
}
func (self *ServiceForBackStage) RpcSendSystemNoticeToPlayer(common *base.Common, reqMsg *share_message.SystemNotice) easygo.IMessage {
	logs.Debug("RpcSendSystemNoticeToPlayer:", reqMsg)
	NoticeAssistant(reqMsg.GetId(), 1, reqMsg.GetTitle(), reqMsg.GetContent())
	return nil
}

//後台操作商城訂單    int64 OrderId	//订单ID   int32 Types	// 1取消订单，2完成发货，3商品上架，4商品下架
func (self *ServiceForBackStage) RpcBsOpShopOrder(common *base.Common, reqMsg *server_server.ShopOrderRequest) easygo.IMessage {
	logs.Info("RpcBsOpShopOrder ", reqMsg)
	shopSrv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_SHOP)
	if shopSrv == nil {
		//ep.Shutdown()
		logs.Info("登录大厅时，找不到有效的商场服务器")
		return nil
	}
	SendMsgToServerNew(shopSrv.GetSid(), "RpcBsOpShopOrder", &server_server.ShopOrderRequest{
		OrderId: reqMsg.OrderId,
		Types:   reqMsg.Types,
		UserId:  easygo.NewInt64(reqMsg.GetUserId())})
	return nil
}

//后台审核出款订单
func (self *ServiceForBackStage) RpcBsAuditOrder(common *base.Common, reqMsg *server_server.AuditOrder) easygo.IMessage {
	logs.Info("RpcBsAuditOrder ", reqMsg)
	order := for_game.GetRedisOrderObj(reqMsg.GetOrderId())
	//第三方发起提现请求
	errMsg := &base.Fail{}
	switch order.GetPayChannel() {
	case for_game.PAY_CHANNEL_PENGJU:
		errMsg = PWebPengJuPay.RechargeOrder(order.GetRedisOrder())
	case for_game.PAY_CHANNEL_HUIJU_DF:
		errMsg = PWebHuiJuPay.ReqSinglePay(order.GetRedisOrder())
	case for_game.PAY_CHANNEL_HUICHAO_DF:
		errMsg = PWebHuiChaoPay.TransferFixed(order.GetRedisOrder())
	}

	if errMsg.GetCode() == for_game.FAIL_MSG_CODE_SUCCESS {
		//第三方下单成功
		order.SetStatus(for_game.ORDER_ST_AUDIT)
		order.SetPayStatus(for_game.PAY_ST_DOING)
	} else {
		//渠道下单失败，如果已经生成订单，则转为人工处理
		order.SetChanneltype(for_game.CHANNEL_MAN_MAKE)
		order.SetExtendValue(errMsg.GetReason())
	}
	order.SaveToMongo()

	return nil
}

//后台请求冻结用户
func (self *ServiceForBackStage) RpcFreezePlayer(common *base.Common, reqMsg *server_server.PlayerIds) easygo.IMessage {
	logs.Info("RpcFreezePlayer ", reqMsg)
	eps := make([]IGameClientEndpoint, 0)
	for _, id := range reqMsg.GetPlayerIds() {
		player := for_game.GetRedisPlayerBase(id)
		if player == nil {
			continue
		}
		player.SetStatus(2)
		player.SetNote(reqMsg.GetNode())
		player.SetIsOnLine(false)
		player.SetOperator(reqMsg.GetOperator())
		player.SetBanOverTime(reqMsg.GetBanOverTime())
		for_game.AddFreezeAccount(player.GetAccount())
		player.SaveToMongo()

		ep := ClientEpMp.LoadEndpoint(id)
		if ep != nil {
			eps = append(eps, ep)
			ep.RpcFreezePlayerLogOut(nil)
		} else {
			SendMsgToHallClientNew(id, "RpcFreezePlayerLogOut", easygo.EmptyMsg)
		}
	}
	//3秒后服务器断开连接
	easygo.AfterFunc(time.Second*5, func() {
		for _, ep := range eps {
			ep.Shutdown()
		}
	})
	return nil
}

//后台修改系统参数设置
func (self *ServiceForBackStage) RpcSysParameterChangeToHall(common *base.Common, reqMsg *server_server.SysteamModId) easygo.IMessage {
	//查询当前可用的系统参数设置
	PSysParameterMgr.UpLoad(reqMsg.GetId())
	//通知前端系统参数改变了
	if reqMsg.GetId() == for_game.LIMIT_PARAMETER {
		msg := PSysParameterMgr.GetSysParameter(reqMsg.GetId())
		BroadCastMsgToOnlineClient("RpcSysParameterChange", msg)
	}
	//由主导大厅进行通知，避免重复通知
	if for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_HALL_SID) == PServerInfo.GetSid() {
		//通知社交广场配置改变
		BroadCastMsgToServerNew(for_game.SERVER_TYPE_SQUARE, "RpcHallNotifyParamChange", reqMsg)
		//通知登录服配置改变
		BroadCastMsgToServerNew(for_game.SERVER_TYPE_LOGIN, "RpcHallNotifyParamChange", reqMsg)
	}
	//修改电竞redis配置
	if reqMsg.GetId() == for_game.ESPORT_PARAMETER {
		//TODO
	}
	return nil
}

func (self *ServiceForBackStage) RpcAddFriend(common *base.Common, reqMsg *server_server.AddPlayerFriendInfo) easygo.IMessage {
	pid := reqMsg.GetPlayerID()
	player := for_game.GetRedisPlayerBase(pid)
	for _, friend_id := range reqMsg.GetList() {
		AgreeAddFriend(pid, friend_id, 1)
		AgreeAddFriend(friend_id, pid, 1)
		player.AddAttention(friend_id)
		player.AddFans(friend_id)
		friend := for_game.GetRedisPlayerBase(friend_id)
		friend.AddAttention(pid)
		friend.AddFans(pid)
		for_game.OperateRedisDynamicAttention(1, pid, friend_id)
		for_game.OperateRedisDynamicAttention(1, friend_id, pid)
	}

	return nil
}

func (self *ServiceForBackStage) RpcCreateTeam(common *base.Common, reqMsg *server_server.CreateTeamInfo) easygo.IMessage {
	logs.Info("后台增加群:", reqMsg)
	playerId := reqMsg.GetPlayerID()
	player := for_game.GetRedisPlayerBase(playerId)
	if player == nil {
		panic("玩家对象为空")
	}
	info := &client_hall.CreateTeam{
		AdminId:  easygo.NewInt64(reqMsg.GetAdminID()),
		TeamName: easygo.NewString(reqMsg.GetTeamName()),
	}
	msg := CreateTeamOperate(playerId, []int64{}, info)
	//通知群主创建成功
	epOwner := ClientEpMp.LoadEndpoint(playerId)
	if epOwner != nil {
		epOwner.RpcCreateTeamResult(msg)
	}

	return &server_server.CreateTeamResult{TeamID: easygo.NewInt64(msg.GetTeam().GetId())}

}

func (self *ServiceForBackStage) RpcEditArticle(common *base.Common, reqMsg *share_message.Tweets) easygo.IMessage {
	players := GetAllPlayers(reqMsg)
	ids := []string{} //离线或者切换到后台的玩家
	//发送文章入库
	now := util.GetMilliTime()
	logs.Info("推送消息入库:", reqMsg.GetID(), len(players))
	var savePlayerTweets []interface{}

	articleResponseList := []*client_server.ArticleResponse{}
	articleUrl := easygo.YamlCfg.GetValueAsString("CLIENT_ARTICLE_URL") //测试服
	for _, article := range reqMsg.GetArticle() {
		articleAdd := articleUrl + "?id=" + strconv.FormatInt(article.GetID(), 10) + "&t=1&pid=" //测试服

		if article.GetTransArticleUrl() != "" && article.TransArticleUrl != nil {
			articleAdd = article.GetTransArticleUrl()
		}
		articleResponse := &client_server.ArticleResponse{
			Id:          easygo.NewInt64(article.GetID()),
			Title:       easygo.NewString(article.GetTitle()),
			Icon:        easygo.NewString(article.GetIcon()),
			ArticleAdd:  easygo.NewString(articleAdd),
			ArticleType: easygo.NewInt32(article.GetArticleType()),
			Location:    easygo.NewInt32(article.GetLocation()),
			IsMain:      easygo.NewInt32(article.GetIsMain()),
			Profile:     easygo.NewString(article.GetProfile()),
			ObjectId:    easygo.NewInt64(article.GetObjectId()),
			ObjPlayerId: easygo.NewInt64(article.GetObjPlayerId()),
		}

		articleResponseList = append(articleResponseList, articleResponse)
		for_game.SetRedisArticleReportFildVal(article.GetID(), int64(len(players)), "PushPlayer") //添加文章报表推送用户数
	}

	notice := &client_server.ArticleListResponse{
		ArticleListId: easygo.NewInt64(reqMsg.GetSendTime()),
		ArticleList:   articleResponseList,
		TweetsId:      easygo.NewInt64(reqMsg.GetID()),
	}
	playerIds := make([]int64, 0) //其他大厅在线玩家
	for _, player := range players {
		//机器人
		if player == nil || player.GetIsRobot() {
			continue
		}
		pid := player.GetPlayerId()
		if player.GetIsOnline() == true {
			playerIds = append(playerIds, pid)
		}
		//处理存储
		b1 := bson.M{"_id": pid}
		b2 := bson.M{}
		//处理在线推送
		if PlayerOnlineMgr.CheckPlayerIsOnLine(pid) && !PlayerOnlineMgr.CheckPlayerIsCutBackstage(pid) { //如果在线并且不处于后台状态
			b2 = bson.M{"$set": bson.M{"PlayerId": pid, "UpdateTime": now, "CreateTime": now}}
		} else {
			//需要推送的人才增加可库
			b2 = bson.M{"$push": bson.M{"TweetsIdList": reqMsg.GetID()}, "$set": bson.M{"PlayerId": pid, "UpdateTime": now, "CreateTime": now}}
			if player.GetToken() != "" {
				ids = append(ids, player.GetToken())
			}
		}
		savePlayerTweets = append(savePlayerTweets, b1, b2)
	}
	//通知其他大厅在线玩家
	easygo.Spawn(BroadCastMsgToHallClientNew, playerIds, "RpcAssistantNotifyArticle", notice)
	//保存
	easygo.Spawn(for_game.UpsertAll, MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_TWEETS, savePlayerTweets)
	//for_game.UpsertAll(MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_TWEETS, savePlayerTweets)
	//离线和切换后台 mob推送
	if reqMsg.GetJgPush() == 1 && reqMsg.JgPush != nil {
		easygo.Spawn(SendTweetsWithJG, ids, reqMsg)
	}
	return easygo.EmptyMsg
}

//小助手极光推送
func SendTweetsWithJG(ids []string, reqMsg *share_message.Tweets) easygo.IMessage {
	var mainArticle *share_message.Article

	for _, article := range reqMsg.GetArticle() {
		if article.GetIsMain() == 2 {
			art := &share_message.Article{
				ID:              easygo.NewInt64(article.GetID()),
				Title:           easygo.NewString(article.GetTitle()),
				Icon:            easygo.NewString(article.GetIcon()),
				ArticleType:     easygo.NewInt32(article.GetArticleType()),
				Location:        easygo.NewInt32(article.GetLocation()),
				IsMain:          easygo.NewInt32(article.GetIsMain()),
				Sort:            easygo.NewInt32(article.GetSort()),
				TransArticleUrl: easygo.NewString(article.GetTransArticleUrl()),
				Profile:         easygo.NewString(article.GetProfile()),
				ObjectId:        easygo.NewInt64(article.GetObjectId()),
			}
			mainArticle = art
		}
	}

	articleUrl := easygo.YamlCfg.GetValueAsString("CLIENT_ARTICLE_URL")                          //测试服
	articleAdd := articleUrl + "?id=" + strconv.FormatInt(mainArticle.GetID(), 10) + "&t=1&pid=" //测试服

	//articleAdd := "http://192.168.150.194:8080/article.html?id=" + strconv.FormatInt(mainArticle.GetID(), 10) //本地

	if mainArticle.ArticleType != nil && mainArticle.GetArticleType() == 2 {
		articleAdd = mainArticle.GetTransArticleUrl()
	}

	m := for_game.PushMessage{
		Title:       mainArticle.GetTitle(),
		Content:     mainArticle.GetProfile(),
		ContentType: for_game.JG_TYPE_BACKSTAGE_ASS,
		Icon:        mainArticle.GetIcon(),
		JumpUrl:     articleAdd,
		ArticleType: mainArticle.GetArticleType(),
		Location:    mainArticle.GetLocation(),
		ArticleId:   mainArticle.GetID(),
		ObjectId:    mainArticle.GetObjectId(),
		PlayerId:    mainArticle.GetObjPlayerId(),
	}
	for_game.JGSendMessage(ids, m)

	return nil
}

//后台请求封禁群操作
func (self *ServiceForBackStage) RpcTeamBanHall(common *base.Common, reqMsg *server_server.TeamManager) easygo.IMessage {

	teamIds := reqMsg.GetTeamIds()
	for _, teamId := range teamIds {
		teamObj := for_game.GetRedisTeamObj(teamId)
		if teamObj == nil {
			logs.Info("找不到这个群信息:", teamId)
			continue
		}
		msg := &client_hall.OperatorMessage{
			Name:      easygo.NewString(reqMsg.GetName()),
			TeamId:    easygo.NewInt64(teamId),
			SendTime:  easygo.NewInt64(reqMsg.GetSendTime()),
			Flag:      easygo.NewInt64(reqMsg.GetFlag()),
			CloseTime: easygo.NewInt64(reqMsg.GetCloseTime()),
			LogId:     easygo.NewInt64(reqMsg.GetLogId()),
		}
		//player := PlayerMgr.LoadPlayer(reqMsg.GetTeamId()) //加载玩家
		//if player != nil {
		ep := ClientEpMp.LoadEndpoint(reqMsg.GetTeamId())
		if ep != nil {
			ep.RpcTeamChangeInfo(msg)
		}
		//}
	}
	return nil
}

//后台请求解封群成员
func (self *ServiceForBackStage) RpcTeamMemCloseAndOpen(common *base.Common, reqMsg *server_server.TeamManager) easygo.IMessage {
	teamObj := for_game.GetRedisTeamObj(reqMsg.GetTeamId())
	if teamObj == nil {
		logs.Info("找不到这个群信息:", reqMsg.GetTeamId())
		return nil
	}
	logs.Info("通知前端取消封禁:", reqMsg)
	msg := &client_hall.OperatorMessage{
		Name:      easygo.NewString(reqMsg.GetName()),
		Members:   reqMsg.GetNickName(),
		TeamId:    easygo.NewInt64(reqMsg.GetTeamId()),
		Flag:      easygo.NewInt64(reqMsg.GetFlag()),
		PlayerId:  reqMsg.GetTeamIds(),
		CloseTime: easygo.NewInt64(reqMsg.GetCloseTime()),
		SendTime:  easygo.NewInt64(reqMsg.GetSendTime()),
		LogId:     easygo.NewInt64(reqMsg.GetLogId()),
	}

	//player := PlayerMgr.LoadPlayer(reqMsg.GetPlayerId()) //加载玩家
	//if player != nil {
	ep1 := ClientEpMp.LoadEndpoint(reqMsg.GetPlayerId())
	if ep1 != nil {
		ep1.RpcTeamMemChangeInfo(msg)
	}
	//}

	return nil

}

//===============================================================================================================客服消息处理===>
func (self *ServiceForBackStage) RpcSendMessageToPlayer(common *base.Common, reqMsg *share_message.IMmessage) easygo.IMessage {
	logs.Info("后台发送客服消息给用户:", reqMsg)
	epPlayer := ClientEpMp.LoadEndpoint(reqMsg.GetPlayerId())
	if epPlayer != nil {
		im := &share_message.IMmessage{}
		im.Id = easygo.NewInt64(reqMsg.GetId())
		im.Cnew = easygo.NewInt32(reqMsg.GetCnew())
		epPlayer.RpcNewWaiterMsg(reqMsg)
	}
	return nil

}

func (self *ServiceForBackStage) RpcEndMessageToPlayer(common *base.Common, reqMsg *share_message.IMmessage) easygo.IMessage {
	logs.Info("后台发送客服消息结束给用户:", reqMsg)
	epPlayer := ClientEpMp.LoadEndpoint(reqMsg.GetPlayerId())
	if epPlayer != nil {
		epPlayer.RpcEndWaiterMsg(reqMsg)
	}
	return nil

}

//后台修改屏蔽词通知到大厅重载
func (self *ServiceForBackStage) RpcEditDirtyWordsToHall(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	for_game.PDirtyWordsMgr.Init()
	return nil
}

// RpcBackstageTop 大厅通知社交广场处理动态定时器.
func (self *ServiceForBackStage) RpcBackstageTop(common *base.Common, reqMsg *server_server.TopRequest) easygo.IMessage {
	logs.Info("==========大厅通知社交广场处理动态定时器 RpcBackstageTop===========")
	dynamic := for_game.GetRedisDynamic(reqMsg.GetLogId())
	for_game.SetTopDynamicToRedis(for_game.REDIS_SQUARE_BS_TOP_DYNAMIC, dynamic) // 添加置顶动态进置顶列表
	serverInfo := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_SQUARE)
	if serverInfo == nil {
		logs.Error("社交广场对象找不到")
		return easygo.NewFailMsg("社交广场对象找不到")
	}

	notifyTopReq := &share_message.BackstageNotifyTopReq{
		LogId:       easygo.NewInt64(reqMsg.GetLogId()),
		TopOverTime: easygo.NewInt64(dynamic.GetTopOverTime()),
	}
	// 通知广场服务器
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_SQUARE, "RpcHallNotifyTop", notifyTopReq)
	return nil
}

//后台获取在线玩家数量
func (self *ServiceForBackStage) RpcGetPlayerOnline(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	lis := PlayerOnlineMgr.PlayerList
	msg := &server_server.PlayerSI{
		Count: easygo.NewInt64(len(lis)),
	}

	return msg
}

//警告群主
func (self *ServiceForBackStage) RpcWarnLordToHall(common *base.Common, reqMsg *share_message.PlayerComplaint) easygo.IMessage {
	pid := reqMsg.GetPlayerId()
	title := "警告"
	NoticeAssistant(pid, 2, title, reqMsg.GetContent())

	if !PlayerOnlineMgr.CheckPlayerIsOnLine(pid) || PlayerOnlineMgr.CheckPlayerIsCutBackstage(pid) {
		player := for_game.GetRedisPlayerBase(pid)
		isNotice := player.GetIsNewMessage()
		if isNotice {
			f := func() {
				isShow := player.GetIsMessageShow()
				ids := for_game.GetJGIds([]int64{pid})
				var content string
				if isShow {
					content = reqMsg.GetContent()
				} else {
					content = "给您发送了一条消息"
				}
				m := for_game.PushMessage{
					Title:       title,
					Content:     content,
					ContentType: for_game.JG_TYPE_BACKSTAGE,
					JumpObject:  2,
				}
				for_game.JGSendMessage(ids, m)
			}
			easygo.Spawn(f)
		}
	}
	return nil
}

//活动通知
func (self *ServiceForBackStage) RpcActivityNotic(common *base.Common, reqMsg *share_message.SystemNotice) easygo.IMessage {
	//id放的是玩家id
	NoticeAssistant(reqMsg.GetId(), 2, reqMsg.GetTitle(), reqMsg.GetContent())
	return nil
}

//活动发奖
func (self *ServiceForBackStage) RpcActivityAddGold(common *base.Common, reqMsg *server_server.Recharge) easygo.IMessage {
	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	if player == nil {
		return easygo.NewFailMsg("用户不存在")
	}

	reason := "活动奖金"
	extendLog := &share_message.GoldExtendLog{
		Title:   easygo.NewString(reason),
		PayType: easygo.NewInt32(99), //99零钱
		Gold:    easygo.NewInt64(player.GetGold() + reqMsg.GetRechargeGold()),
	}

	NotifyAddGold(reqMsg.GetPlayerId(), reqMsg.GetRechargeGold(), reason, reqMsg.GetSourceType(), extendLog)
	return nil
}

//后台赠送商品道具给玩家
func (self *ServiceForBackStage) RpcSysGivePropsToHall(common *base.Common, reqMsg *server_server.SysGivePropsRequest) easygo.IMessage {
	logs.Debug("RpcSysGivePropsToHall", reqMsg)
	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	if player == nil {
		return easygo.NewFailMsg("用户不存在")
	}
	count := reqMsg.GetNum()
	if count <= 0 {
		return easygo.NewFailMsg("赠送道具商品不能为0")
	}

	product := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT, bson.M{"_id": reqMsg.GetProductId()})
	if product == nil {
		return easygo.NewFailMsg("赠送道具商品不存在")
	}

	var items []*share_message.CoinProduct
	one := &share_message.CoinProduct{}
	for_game.StructToOtherStruct(product, one)
	one.ProductNum = easygo.NewInt64(count)
	items = append(items, one)

	NotifyAddBagItems(reqMsg.GetPlayerId(), items, for_game.COIN_ITEM_GETTYPE_SEND, 0, "", reqMsg.GetOperator())

	days := fmt.Sprintf("%d天", one.GetEffectiveTime())
	if one.GetEffectiveTime() == for_game.COIN_PROPS_FOREVER {
		days = "永久"
	}

	content := fmt.Sprintf("亲爱的%s，恭喜你收到虚拟物品%s，有效期为%s，赶快去你的物品中查看吧。", player.GetNickName(), one.GetName(), days)
	NoticeAssistant(reqMsg.GetPlayerId(), 1, "温馨提示", content)

	return nil
}

//日志回收用户道具
func (self *ServiceForBackStage) RpcRecyclePropsToHall(common *base.Common, reqMsg *server_server.PlayerIds) easygo.IMessage {
	ids := reqMsg.GetPlayerIds() //要回收的道具日志ID
	if len(ids) == 0 {
		return easygo.NewFailMsg("回收的道具日志ID不能为空")
	}

	ls, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_GETPROPS_LOG, bson.M{"_id": bson.M{"$in": ids}}, 0, 0)
	if count == 0 {
		return easygo.NewFailMsg("回收的道具日志不存在")
	}
	var okIds []int64
	for _, li := range ls {
		one := &share_message.PlayerGetPropsLog{}
		for_game.StructToOtherStruct(li, one)

		if one.GetGetType() == 5 {
			continue //回收类型的日志直接跳过
		}

		t := int64(-1)
		if one.GetEffectiveTime() > 0 {
			t = one.GetEffectiveTime() * 86400
		}
		NoticeReduceBagItem(one.GetPlayerId(), one.GetBagId(), t)
		okIds = append(okIds, one.GetId())
	}

	msg := &server_server.PlayerIds{
		PlayerIds: okIds,
	}

	return msg
}

//背包回收用户道具
func (self *ServiceForBackStage) RpcBagRecyclePropsToHall(common *base.Common, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	id := reqMsg.GetId64() //要回收的道具背包ID
	if id == 0 {
		return easygo.NewFailMsg("回收的道具日志ID不能为空")
	}
	day := reqMsg.GetId32()
	if day == 0 {
		return easygo.NewFailMsg("回收天数错误")
	}

	one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BAGITEM, bson.M{"_id": id})
	if one == nil {
		return easygo.NewFailMsg("回收的道具不存在")
	}
	bagItem := &share_message.PlayerBagItem{}
	for_game.StructToOtherStruct(one, bagItem)

	t := int64(-1)
	if bagItem.GetOverTime() > 0 {
		if day > 0 {
			t = int64(day) * 86400
		} else {
			return easygo.NewFailMsg("回收的道具天数错误")
		}

	}

	playerid := bagItem.GetPlayerId()
	NoticeReduceBagItem(playerid, id, t)

	nowtime := easygo.NowTimestamp()
	data := &share_message.PlayerGetPropsLog{
		Id:            easygo.NewInt64(for_game.NextId(for_game.TABLE_PLAYER_GETPROPS_LOG)),
		PlayerId:      easygo.NewInt64(playerid),
		GivePlayerId:  easygo.NewInt64(0),
		PropsId:       easygo.NewInt64(bagItem.GetPropsId()),
		PropsNum:      easygo.NewInt64(1),
		GetType:       easygo.NewInt32(for_game.COIN_ITEM_GETTYPE_BACK),
		CreateTime:    easygo.NewInt64(nowtime),
		EffectiveTime: easygo.NewInt64(day),
		BagId:         easygo.NewInt64(id),
		RecycleTime:   easygo.NewInt64(nowtime),
		Operator:      easygo.NewString(reqMsg.GetIdStr()),
		Note:          easygo.NewString(reqMsg.GetNote()),
	}
	easygo.Spawn(for_game.AddGetPropsLog, data)

	return nil
}

//后台请求发送电竞系统消息
func (self *ServiceForBackStage) RpcSendSportSysNoticeToHall(common *base.Common, reqMsg *share_message.TableESPortsSysMsg) easygo.IMessage {
	logs.Info("RpcSendSportSysNoticeToHall", reqMsg.GetId())
	if reqMsg.GetStatus() != for_game.ESPORTS_STATUS_1 {
		reqMsg.Status = easygo.NewInt32(for_game.ESPORTS_STATUS_1)
	}

	findBson := bson.M{"IsRobot": false, "Status": for_game.ACCOUNT_NORMAL}
	if reqMsg.GetRecipientType() > 0 {
		findBson["DeviceType"] = reqMsg.GetRecipientType()
	}
	//批量获取用户，每次获取5000个
	var maxId int64 = for_game.INIT_PLAYER_ID
	for true {
		plist := for_game.GetSomeNoticePlayers(maxId, findBson)
		if len(plist) == 0 {
			break
		}
		ids := make([]int64, 0)
		for _, p := range plist {
			ids = append(ids, p.GetPlayerId())
		}
		DealJGNoticeToClient(ids, reqMsg)
		maxId = plist[len(plist)-1].GetPlayerId()
	}
	//前端推送在线玩家
	if reqMsg.GetIsMessageCenter() {
		//pids 推送对象
		//TODO 推送到前端
		NoticeAssistantEsportSysMsg(reqMsg)
	}
	logs.Info("推送处理完成:", reqMsg.GetId())
	return nil
}

//极光推送给所有用户
func DealJGNoticeToClient(pIds []int64, reqMsg *share_message.TableESPortsSysMsg) *base.Fail {
	ids := for_game.GetJGIds(pIds)
	if len(ids) == 0 {
		return easygo.NewFailMsg("推送用户不存在")
	}
	//极光推送
	if reqMsg.GetIsPush() {
		jumpInfo := reqMsg.GetJumpInfo()
		m := for_game.PushMessage{
			Title:       reqMsg.GetTitle(),
			Content:     reqMsg.GetContent(),
			ContentType: for_game.JG_TYPE_BACKSTAGE_ESPORT,
			Icon:        reqMsg.GetIcon(),
			JumpType:    jumpInfo.GetJumpType(),
			Location:    jumpInfo.GetJumpObject(),
			ObjectId:    jumpInfo.GetJumpObjId(),
			JumpUrl:     jumpInfo.GetJumpUrl(),
		}

		for_game.JGSendMessage(ids, m)
	}
	return nil
}

//后台请求人工查询充值订单更新订单状态并发货
func (self *ServiceForBackStage) RpcCheckOrderToHall(common *base.Common, reqMsg *brower_backstage.OptOrderRequest) easygo.IMessage {
	order := for_game.GetRedisOrderObj(reqMsg.GetOid())
	if order == nil {
		return easygo.NewFailMsg("订单不存在")
	}

	//第三方发起提现请求
	switch order.GetPayChannel() {
	case for_game.PAY_CHANNEL_HUICHAO_WX, for_game.PAY_CHANNEL_HUICHAO_ZFB, for_game.PAY_CHANNEL_HUICHAO_YL:
		PWebHuiChaoPay.ReqCheckPayOrder(reqMsg.GetOid(), time.Second*3598)
	default:
		return easygo.NewFailMsg("该渠道暂不支持此功能")
	}
	return nil
}

//后台通知大厅有人成为话题主,把话题主变成相应话题群组的群主
func (self *ServiceForBackStage) RpcChangeTopicOwner(common *base.Common, reqMsg *server_server.PlayerSI) easygo.IMessage {
	logs.Info("RpcChangeTopicOwner:", reqMsg)
	teams := for_game.GetTeamsByTopic(reqMsg.GetAccount())
	for _, team := range teams {
		teamId := team.GetId()
		teamObj := for_game.GetRedisTeamObj(teamId, team)
		if teamObj == nil {
			continue
		}
		ownerId := teamObj.GetTeamOwner()
		if ownerId == reqMsg.GetPlayerId() { //已经时群主了，不需要更换
			continue
		}
		//群主把群转让给话题主
		player := GetPlayerObj(ownerId)
		if player == nil {
			continue
		}
		playerList := []int64{reqMsg.GetPlayerId()}
		//先进群
		reason := AddTeamMemberOperate(teamId, reqMsg.GetPlayerId(), 10001, playerList, player.GetNickName())
		if reason != "" {
			return easygo.NewFailMsg(reason)
		}
		memberObj := for_game.GetRedisTeamPersonalObj(teamId)
		memberObj.SaveToMongo()
		teamObj.SaveToMongo()
		//更换群主
		SetTeamMemberPosition(teamId, playerList, for_game.TEAM_OWNER)
		NoticeTeamMessage(teamId, ownerId, for_game.CHANGE_OWNER, playerList)
	}
	return nil
}
