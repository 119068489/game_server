// 大厅服务器为[游戏客户端]提供的服务

package hall

import (
	"encoding/base64"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"time"

	"github.com/astaxie/beego/logs"
)

//获取匹配的语音名牌
func (self *cls1) RpcGetVoiceCards(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetVoiceCards:", reqMsg)
	//TODO 检测当前是否满足100次
	param := PSysParameterMgr.GetSysParameter(for_game.COMMON_PARAMETER)
	maxNum := for_game.VOICE_CARD_FRESH_MAX_NUM
	if param != nil {
		maxNum = int(param.GetDayMaxMatchTimes())
	}
	logs.Info("当前限制次数:", maxNum)
	period := for_game.GetPlayerPeriod(who.GetPlayerId())
	freshNum := period.DayPeriod.FetchInt(for_game.CHECK_FRESH_NUM)
	if freshNum >= maxNum {
		logs.Error("玩家今日刷新已达上限！id=", who.GetPlayerId())
		return easygo.NewFailMsg("今日刷新次数已达上限!", for_game.FAIL_MSG_CODE_1017)
	}
	list := for_game.GetVoiceCardToDB(who.GetPlayerId(), false)
	resp := &client_hall.LoveMatchResp{
		Cards: list,
	}
	period.DayPeriod.AddInteger(for_game.CHECK_FRESH_NUM, 1)
	//logs.Info("RpcGetVoiceCards 返回数据:", resp)
	return resp
}

//新版语音匹配名片
func (self *cls1) RpcGetVoiceCardsNew(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.IsFirstLogin, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetVoiceCards:", reqMsg)
	//TODO 检测当前是否满足100次
	param := PSysParameterMgr.GetSysParameter(for_game.COMMON_PARAMETER)
	maxNum := for_game.VOICE_CARD_FRESH_MAX_NUM
	if param != nil {
		maxNum = int(param.GetDayMaxMatchTimes())
	}
	logs.Info("当前限制次数:", maxNum)
	period := for_game.GetPlayerPeriod(who.GetPlayerId())
	freshNum := period.DayPeriod.FetchInt(for_game.CHECK_FRESH_NUM)
	if freshNum >= maxNum {
		logs.Error("玩家今日刷新已达上限！id=", who.GetPlayerId())
		return easygo.NewFailMsg("今日刷新次数已达上限!", for_game.FAIL_MSG_CODE_1017)
	}
	list := for_game.GetVoiceCardToDB(who.GetPlayerId(), reqMsg.GetIsFirstReq())
	resp := &client_hall.LoveMatchResp{
		Cards: list,
	}
	period.DayPeriod.AddInteger(for_game.CHECK_FRESH_NUM, 1)
	//logs.Info("RpcGetVoiceCards 返回数据:", resp)
	return resp
}

//点赞别人的语音卡片
func (self *cls1) RpcZanVoiceCard(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PlayerInfoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcZanVoiceCard", reqMsg, who.GetPlayerId())
	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	resp := &client_hall.VCZanResult{
		Result: easygo.NewBool(false),
	}
	if player == nil {
		logs.Error("无效的玩家id:", reqMsg.GetPlayerId())
		return resp
	}
	if who.GetPlayerId() == player.GetPlayerId() {
		logs.Error("不能给自己点赞", reqMsg.GetPlayerId())
		return resp
	}
	//写行为日志
	zanLog := for_game.CheckVoiceCardZan(who.GetPlayerId(), reqMsg.GetPlayerId())
	if zanLog != nil {
		//已经赞过，最多只能点赞6次
		if zanLog.GetZanNum() >= for_game.VC_MAX_ZAN_NUM {
			return resp
		}
		for_game.AddVoiceCardZanNum(zanLog.GetId())
		player.AddVCZanNum(1)
		resp.Result = easygo.NewBool(true)
		return resp
	}
	//增加赞记录
	log := &share_message.PlayerVCZanLog{
		Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_PLAYER_VC_ZAN_LOG)),
		PlayerId:   easygo.NewInt64(who.GetPlayerId()),
		TargetId:   easygo.NewInt64(reqMsg.GetPlayerId()),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
		ZanNum:     easygo.NewInt32(1),
	}
	for_game.AddVoiceCardZanLog(log)

	//增关注记录
	for_game.AddAttentionLog(who.GetPlayerId(), reqMsg.GetPlayerId(), for_game.VC_ATTENTION_LIKE)

	player.AddVCZanNum(1)
	resp.Result = easygo.NewBool(true)
	return resp
}

//获取喜欢我的列表
func (self *cls1) RpcGetLoveMeList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.LoveMeReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetLoveMeList: ", reqMsg, who.GetPlayerId())
	cnt := 20 //默认每次拿20个
	myId := who.GetPlayerId()
	me := for_game.GetRedisPlayerBase(who.GetPlayerId())
	playerIds := make([]int64, 0)
	playerMixIds := make([]int64, 0)
	sessions := make([]string, 0)
	cards := make([]*client_hall.VoiceCard, 0)
	tagIds := make([]int32, 0)
	voiceMap := make(map[int64]string, 0)                              //玩家id，语音链接
	logsMap := make(map[int64]string, 0)                               //目标玩家ID，打招呼内容
	playerTagMap := make(map[int64][]int32, 0)                         //玩家id，用户最新三条标签
	commonTagMap := make(map[int64][]int32, 0)                         //玩家id，共同的标签
	tagDataMap := make(map[int32]*share_message.InterestTag, 0)        //标签id，标签内容
	constellationMap := make(map[int32]string, 0)                      //星座id，星座名字
	playersMap := make(map[int64]*share_message.PlayerBase, 0)         //星座id，星座名字
	attentionLogs := make(map[int64]*share_message.PlayerAttentionLog) //发起玩家id， 关注信息
	attentionList := for_game.GetAttentionPlayers(myId, int(reqMsg.GetPage()), cnt, for_game.VC_ATTENTION_TO_ME)
	for _, v := range attentionList {
		playerIds = append(playerIds, v.GetPlayerId())
		attentionLogs[v.GetPlayerId()] = v
	}

	players := for_game.GetAllPlayerBase(playerIds, false)
	for _, v := range players {
		playersMap[v.GetPlayerId()] = v
		playerMixIds = append(playerMixIds, v.GetMixId())
	}
	//获取用户语音名片
	voiceUrls := for_game.GetVoiceCardInfoByIds(playerMixIds)
	for _, v := range voiceUrls {
		voiceMap[v.GetPlayerId()] = v.GetMixVoiceUrl()
	}

	//获取所有星座
	for _, v := range for_game.GetConfigConstellationFormDB() {
		constellationMap[v.GetId()] = v.GetName()
	}
	for _, v := range playerIds {
		player := playersMap[v]
		matchDegree, commonTags := GetPlayerMatchDegree(me, player)
		tags := player.GetPersonalityTags()
		if len(tags) > 3 {
			tags = tags[:3]
		}
		commonTagMap[v] = commonTags
		playerTagMap[v] = tags
		tagIds = append(tagIds, commonTags...)
		tagIds = append(tagIds, tags...)

		if attentionLogs[v].GetOpt() == for_game.VC_ATTENTION_HI {
			sessions = append(sessions, for_game.MakeSessionKey(v, myId))
		}
		attentionTime := attentionLogs[v].GetSortTime()
		if attentionTime > 1000000000000 {
			attentionTime = attentionTime / 1000
		}
		cards = append(cards, &client_hall.VoiceCard{
			PlayerId:       easygo.NewInt64(v),
			NickName:       easygo.NewString(player.GetNickName()),
			HeadUrl:        easygo.NewString(player.GetHeadIcon()),
			Sex:            easygo.NewInt32(player.GetSex()),
			Constellation:  easygo.NewInt32(player.GetConstellation()),
			MatchingDegree: easygo.NewInt32(matchDegree),
			ZanNum:         easygo.NewInt32(player.GetVCZanNum()),
			VoiceUrl:       easygo.NewString(voiceMap[v]),
			IsOnLine:       easygo.NewBool(player.GetIsOnline()),
			BgUrl:          easygo.NewString(player.GetBgImageUrl()),
			AttentionType:  easygo.NewInt32(attentionLogs[v].GetOpt()),
			AttentionTime:  easygo.NewInt64(attentionTime),
		})
	}
	for _, v := range for_game.GetPlayerPersonalityTagAllData(tagIds) {
		tagDataMap[v.GetId()] = v
	}
	//获取用户sayHi内容
	sayHiLogs := for_game.GetAllSayHiLog(sessions, for_game.TALK_CONTENT_SAY_HI_WORD)
	for _, v := range sayHiLogs {
		logsMap[v.GetTalker()] = v.GetContent()
	}

	for _, v := range cards {
		commonTagStr := make([]string, 0)
		playerTagStr := make([]*client_hall.PersonTag, 0)
		for _, tag := range commonTagMap[v.GetPlayerId()] {
			commonTagStr = append(commonTagStr, tagDataMap[tag].GetName())
		}
		for _, tag := range playerTagMap[v.GetPlayerId()] {
			playerTagStr = append(playerTagStr, &client_hall.PersonTag{
				Id:   easygo.NewInt32(tagDataMap[tag].GetId()),
				Name: easygo.NewString(tagDataMap[tag].GetName()),
			})
		}
		v.PersonalityTags = playerTagStr
		v.CommonTags = commonTagStr
		v.Content = easygo.NewString(logsMap[v.GetPlayerId()])
	}
	resp := &client_hall.LoveMeResp{
		Page:  easygo.NewInt32(reqMsg.GetPage()),
		Cards: cards,
	}
	return resp
}

//获取我喜欢的列表
func (self *cls1) RpcGetMyLoveList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.MyLoveReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetMyLoveList: ", reqMsg, who.GetPlayerId())
	cnt := 20 //默认每次拿20个
	playerIds := make([]int64, 0)
	playerMixIds := make([]int64, 0)
	sessions := make([]string, 0)
	cards := make([]*client_hall.VoiceCard, 0)
	tagIds := make([]int32, 0)
	voiceMap := make(map[int64]string, 0) //玩家id，语音链接
	//logsMap := make(map[int64]string, 0)                        //目标玩家ID，打招呼内容
	playerTagMap := make(map[int64][]int32, 0)                  //玩家id，用户最新三条标签
	tagDataMap := make(map[int32]*share_message.InterestTag, 0) //标签id，标签内容
	//constellationMap := make(map[int32]string, 0)                         //星座id，星座名字
	playersMap := make(map[int64]*share_message.PlayerBase, 0)            //星座id，星座名字
	attentionLogs := make(map[int64]*share_message.PlayerAttentionLog, 0) //被关注用户id，关注内容
	myId := who.GetPlayerId()
	me := for_game.GetRedisPlayerBase(myId)

	attentionList := for_game.GetAttentionPlayers(myId, int(reqMsg.GetPage()), cnt, for_game.VC_ATTENTION_TO_OTHER)
	for _, v := range attentionList {
		playerIds = append(playerIds, v.GetTargetId())
		attentionLogs[v.GetTargetId()] = v
	}

	players := for_game.GetAllPlayerBase(playerIds, false)
	for _, v := range players {
		playersMap[v.GetPlayerId()] = v
		playerMixIds = append(playerMixIds, v.GetMixId())
	}
	//获取用户语音名片
	voiceUrls := for_game.GetVoiceCardInfoByIds(playerMixIds)
	for _, v := range voiceUrls {
		voiceMap[v.GetPlayerId()] = v.GetMixVoiceUrl()
	}

	//获取所有星座
	//for _, v := range for_game.GetConfigConstellationFormDB() {
	//	constellationMap[v.GetId()] = v.GetName()
	//}

	for _, v := range playerIds {
		player := playersMap[v]
		matchDegree, _ := GetPlayerMatchDegree(me, player)
		tags := player.GetPersonalityTags()
		if len(tags) > 3 {
			tags = tags[:3]
		}
		playerTagMap[v] = tags
		tagIds = append(tagIds, tags...)
		if attentionLogs[v].GetOpt() == for_game.VC_ATTENTION_HI {
			sessions = append(sessions, for_game.MakeSessionKey(v, myId))
		}
		attentionTime := attentionLogs[v].GetSortTime()
		if attentionTime > 1000000000000 {
			attentionTime = attentionTime / 1000
		}
		cards = append(cards, &client_hall.VoiceCard{
			PlayerId:       easygo.NewInt64(v),
			NickName:       easygo.NewString(player.GetNickName()),
			HeadUrl:        easygo.NewString(player.GetHeadIcon()),
			Sex:            easygo.NewInt32(player.GetSex()),
			MatchingDegree: easygo.NewInt32(matchDegree),
			VoiceUrl:       easygo.NewString(voiceMap[v]),
			IsOnLine:       easygo.NewBool(player.GetIsOnline()),
			AttentionTime:  easygo.NewInt64(attentionTime),
			AttentionType:  easygo.NewInt32(attentionLogs[v].GetOpt()),
			Constellation:  easygo.NewInt32(player.GetConstellation()),
		})
	}
	for _, v := range for_game.GetPlayerPersonalityTagAllData(tagIds) {
		tagDataMap[v.GetId()] = v
	}

	//获取用户sayHi内容
	//sayHiLogs := for_game.GetAllSayHiLog(sessions)
	//for _, v := range sayHiLogs {
	//	logsMap[v.GetTargetId()] = v.GetContent()
	//}

	for _, v := range cards {
		playerTag := make([]*client_hall.PersonTag, 0)

		for _, tag := range playerTagMap[v.GetPlayerId()] {
			playerTag = append(playerTag, &client_hall.PersonTag{
				Id:   easygo.NewInt32(tagDataMap[tag].GetId()),
				Name: easygo.NewString(tagDataMap[tag].GetName()),
			})
		}
		v.PersonalityTags = playerTag
		//v.Content = easygo.NewString(logsMap[v.GetPlayerId()])
	}
	resp := &client_hall.LoveMeResp{
		Page:  easygo.NewInt32(reqMsg.GetPage()),
		Cards: cards,
	}
	return resp
}

//向指定玩家发起sayHi操作
func (self *cls1) RpcSayHiToPlayer(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PlayerInfoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcSayHiToPlayer", reqMsg, who.GetPlayerId())
	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	if player == nil {
		return easygo.NewFailMsg("无效的玩家id")
	}
	if player.GetPlayerId() == who.GetPlayerId() {
		return easygo.NewFailMsg("自己不能给自己打招呼")
	}
	//写行为日志
	hiLog := for_game.CheckSayHiToPlayer(who.GetPlayerId(), reqMsg.GetPlayerId())
	if hiLog != nil {
		//已经打过招呼，不在处理
		logs.Info("已经打过招呼了")
		return nil
	}
	//增加打招呼sayHi记录
	log := &share_message.PlayerVCSayHiLog{
		Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_PLAYER_VC_SAY_HI_LOG)),
		PlayerId:   easygo.NewInt64(who.GetPlayerId()),
		TargetId:   easygo.NewInt64(reqMsg.GetPlayerId()),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
	}
	for_game.AddVoiceCardSayHiLog(log)
	//增关注记录
	for_game.AddAttentionLog(who.GetPlayerId(), reqMsg.GetPlayerId(), for_game.VC_ATTENTION_HI)
	id := for_game.MakeSessionKey(who.GetPlayerId(), reqMsg.GetPlayerId())
	obj := for_game.GetRedisPlayerIntimacyObj(id)
	t := int32(share_message.AddFriend_Type_VOICE_CARD)
	if obj != nil && !obj.GetIsSayHi() {
		//组装匹配卡片信息发送
		content := who.GetSayHiData(reqMsg.GetPlayerId())

		chatLog := &share_message.Chat{
			SessionId:   easygo.NewString(for_game.MakeSessionKey(reqMsg.GetPlayerId(), who.GetPlayerId())),
			SourceId:    easygo.NewInt64(who.GetPlayerId()),
			TargetId:    easygo.NewInt64(reqMsg.GetPlayerId()),
			Content:     easygo.NewString(content),
			ChatType:    easygo.NewInt32(for_game.CHAT_TYPE_PRIVATE),
			ContentType: easygo.NewInt32(for_game.TALK_CONTENT_SAY_HI),
			SayType:     easygo.NewInt32(t),
		}
		self.RpcChatNew(nil, who, chatLog)
		obj.SetIIsSayHi(true)
	}
	//打招呼内容发送
	sayHi := for_game.GetRandSayHi()
	content1 := base64.StdEncoding.EncodeToString([]byte(sayHi))
	chatLog1 := &share_message.Chat{
		SessionId:   easygo.NewString(for_game.MakeSessionKey(reqMsg.GetPlayerId(), who.GetPlayerId())),
		SourceId:    easygo.NewInt64(who.GetPlayerId()),
		TargetId:    easygo.NewInt64(reqMsg.GetPlayerId()),
		Content:     easygo.NewString(content1),
		ChatType:    easygo.NewInt32(for_game.CHAT_TYPE_PRIVATE),
		ContentType: easygo.NewInt32(for_game.TALK_CONTENT_SAY_HI_WORD),
		SayType:     easygo.NewInt32(t),
	}
	self.RpcChatNew(nil, who, chatLog1)
	return nil
}

//获取指定玩家的语音作品
func (self *cls1) RpcGetVoiceCardList(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PlayerInfoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetVoiceCardList:", reqMsg, who.GetPlayerId())
	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	if player == nil {
		return easygo.NewFailMsg("无效的玩家id")
	}
	list := for_game.GetMyMixVideo(who.GetPlayerId(), reqMsg.GetPlayerId())
	bgIds := make([]int64, 0)
	bgTagIds := make([]int32, 0)
	for _, mix := range list {
		bgIds = append(bgIds, mix.GetBgId())
	}
	data := for_game.GetBgVideoData(bgIds)
	mData := make(map[int64]*share_message.BgVoiceVideo)

	for _, d := range data {
		bgTagIds = append(bgTagIds, d.GetTags()...)
		mData[d.GetId()] = d
	}
	mTags := for_game.GetBgVideoTags(bgTagIds)
	mixList := make([]*client_hall.MixVideo, 0)

	for _, mix := range list {
		bgVideo := mData[mix.GetBgId()]
		tags := make([]string, 0)
		for _, t := range bgVideo.GetTags() {
			tags = append(tags, mTags[int32(t)].GetName())
		}
		newBgVideo := &client_hall.VoiceVideo{
			Id:        easygo.NewInt64(bgVideo.GetId()),
			Maker:     easygo.NewString(bgVideo.GetMaker()),
			Name:      easygo.NewString(bgVideo.GetName()),
			Tags:      tags,
			Content:   easygo.NewString(bgVideo.GetContent()),
			MusicUrl:  easygo.NewString(bgVideo.GetMusicUrl()),
			ImageUrl:  easygo.NewString(bgVideo.GetImageUrl()),
			Type:      easygo.NewInt32(bgVideo.GetType()),
			MusicTime: easygo.NewInt64(bgVideo.GetMusicTime()),
		}
		d := &client_hall.MixVideo{
			Id:          easygo.NewInt64(mix.GetId()),
			PlayerId:    easygo.NewInt64(mix.GetPlayerId()),
			BgVideo:     newBgVideo,
			MixVoiceUrl: easygo.NewString(mix.MixVoiceUrl),
			IsCard:      easygo.NewBool(mix.GetId() == player.GetMixId()),
			MixTime:     easygo.NewInt64(mix.GetMixTime()),
		}
		mixList = append(mixList, d)
	}
	backMsg := &client_hall.VoiceCardListResp{
		PlayerId: easygo.NewInt64(reqMsg.GetPlayerId()),
		Cards:    mixList,
	}
	return backMsg
}

//获取玩家个性化标签
func (self *cls1) RpcGetPersonalityTags(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PlayerInfoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetPersonalityTags:", reqMsg, who.GetPlayerId())
	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	if player == nil {
		return easygo.NewFailMsg("无效的玩家id")
	}
	dbList := for_game.GetPlayerPersonalityTags(player.GetPersonalityTags())
	tags := make([]*client_hall.PersonTag, 0)
	for _, d := range player.GetPersonalityTags() {
		for _, db := range dbList {
			if db.GetId() == d {
				tag := &client_hall.PersonTag{
					Id:   easygo.NewInt32(db.GetId()),
					Name: easygo.NewString(db.GetName()),
				}
				tags = append(tags, tag)
				break
			}
		}
	}
	backMsg := &client_hall.PersonalityTagsResp{
		PlayerId:        easygo.NewInt64(reqMsg.GetPlayerId()),
		PersonalityTags: tags,
	}
	return backMsg
}

//获取新的喜欢我条目数量
func (self *cls1) RpcGetLoveMeNewNum(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetLoveMeNewNum:", reqMsg, who.GetPlayerId())
	t := who.GetReadLoveMeLogTime()
	num, totalNum := for_game.GetLoveMeNewNum(who.GetPlayerId(), t)
	msg := &client_hall.LoveMeData{
		Num:      easygo.NewInt32(num),
		TotalNum: easygo.NewInt32(totalNum),
	}
	return msg
}

//读取我喜欢列表
func (self *cls1) RpcReadLoveMeLog(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcReadLoveMeLog:", reqMsg, who.GetPlayerId())
	who.SetReadLoveMeLogTime(for_game.GetMillSecond())
	return nil
}

//获取指定玩家语音名片信息
func (self *cls1) RpcGetVoiceCard(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PlayerInfoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetVoiceCard:", reqMsg, who.GetPlayerId())
	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	if player == nil {
		return easygo.NewFailMsg("无效的玩家id")
	}
	card := player.GetVoiceCardData()
	return card
}

//翻页获取个性化系统标签
func (self *cls1) RpcSysPersonalityTags(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PersonalityTagReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcSysPersonalityTags:", reqMsg, who.GetPlayerId())
	dbList := for_game.GetPagePersonalityTags(int(reqMsg.GetNum()), int(reqMsg.GetPage()))
	tags := make([]*client_hall.PersonTag, 0)
	for _, d := range dbList {
		tag := &client_hall.PersonTag{
			Id:   easygo.NewInt32(d.GetId()),
			Name: easygo.NewString(d.GetName()),
		}
		tags = append(tags, tag)
	}
	reqMsg.PersonalityTags = tags
	return reqMsg
}

//获取系统背景图列表
func (self *cls1) RpcChangeSystemBgImage(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SystemBgReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcChangeSystemBgImage:", reqMsg, who.GetPlayerId())
	data := for_game.GetRandomBgImageUrl(reqMsg.GetUrl())
	resp := &client_hall.SystemBgResp{
		Url: easygo.NewString(data.GetUrl()),
	}
	return resp
}

//上传背景音乐，图片，
func (self *cls1) RpcMakeVoiceVideo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.VoiceVideo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcMakeVoiceVideo:", reqMsg, who.GetPlayerId())

	tagIds := make([]int32, 0)
	for _, t := range reqMsg.GetTagIds() {
		tagIds = append(tagIds, t)

	}
	player := for_game.GetRedisPlayerBase(who.GetPlayerId())
	//增加记录
	log := &share_message.BgVoiceVideo{
		Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_BG_VOICE_VIDEO)),
		PlayerId:   easygo.NewInt64(who.GetPlayerId()),
		Maker:      easygo.NewString(reqMsg.GetMaker()),
		Name:       easygo.NewString(reqMsg.GetName()),
		Tags:       tagIds,
		Content:    easygo.NewString(reqMsg.GetContent()),
		MusicUrl:   easygo.NewString(reqMsg.GetMusicUrl()),
		ImageUrl:   easygo.NewString(reqMsg.GetImageUrl()),
		Type:       easygo.NewInt32(reqMsg.GetType()),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
		PlayerType: easygo.NewInt64(player.GetTypes()),
		UseCount:   easygo.NewInt64(0),
		Status:     easygo.NewInt32(0),
		MusicTime:  easygo.NewInt64(reqMsg.GetMusicTime()),
	}
	for_game.AddBgVoiceVideo(log)
	logs.Info("RpcMakeVoiceVideo:制作完成")
	reqMsg.Id = easygo.NewInt64(log.GetId())
	return reqMsg
}

//搜索背景音乐
func (self *cls1) RpcSearchVoiceVideo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SearchVoiceVideoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcSearchVoiceVideo:", reqMsg, who.GetPlayerId())
	bgVoice := for_game.GetBgVoice(who.GetPlayerId(), reqMsg)
	data := make([]*client_hall.VoiceVideo, 0)
	for _, v := range bgVoice {
		data = append(data, &client_hall.VoiceVideo{
			Id:        easygo.NewInt64(v.GetId()),
			Maker:     easygo.NewString(v.GetMaker()),
			Name:      easygo.NewString(v.GetName()),
			Content:   easygo.NewString(v.GetContent()),
			Type:      easygo.NewInt32(v.GetType()),
			TagIds:    v.GetTags(),
			MusicUrl:  easygo.NewString(v.GetMusicUrl()),
			ImageUrl:  easygo.NewString(v.GetImageUrl()),
			MusicTime: easygo.NewInt64(v.GetMusicTime()),
		})
	}
	resp := &client_hall.SearchVoiceVideoResp{
		Data: data,
	}
	return resp
}

//获取所有背景音乐素材标签列表
func (self *cls1) RpcGetVoiceTags(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SearchVoiceVideoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetVoiceTags:", reqMsg, who.GetPlayerId())
	reqType := reqMsg.GetType()
	topTags, laterTags := for_game.GetVoiceTags(reqType, 10, -1)
	tags := make([]*client_hall.PersonTag, 0)
	for _, v := range topTags {
		tags = append(tags, &client_hall.PersonTag{
			Id:    easygo.NewInt32(v.GetId()),
			Name:  easygo.NewString(v.GetName()),
			IsHot: easygo.NewBool(true),
		})
	}
	for _, v := range laterTags {
		tags = append(tags, &client_hall.PersonTag{
			Id:   easygo.NewInt32(v.GetId()),
			Name: easygo.NewString(v.GetName()),
		})
	}
	resp := &client_hall.PersonalityTagsResp{
		PersonalityTags: tags,
	}
	return resp
}

//录制语音名片
func (self *cls1) RpcMixVoiceVideo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.MixVoiceVideo, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcMixVoiceVideo", reqMsg, who.GetPlayerId())
	bgm := for_game.GetOneBgVideoData(reqMsg.GetBgId())
	if bgm == nil {
		reqMsg.Result = easygo.NewInt32(1)
		logs.Error("无效的背景音乐id:", reqMsg.GetBgId())
		return reqMsg
	}
	if reqMsg.GetMyVoiceUrl() == "" {
		reqMsg.Result = easygo.NewInt32(2)
		logs.Error("录制的音频无效:")
		return reqMsg
	}
	//if for_game.GetCountMixVideo(who.GetPlayerId()) > 10 {
	//	reqMsg.Result = easygo.NewInt32(4)
	//	logs.Error("声音名片已达上限")
	//	return reqMsg
	//}
	newUrl := reqMsg.GetMyVoiceUrl()
	if reqMsg.GetBgVolume() != 0 && bgm.GetMusicUrl() != "" {
		//背景音乐不为0，去合成
		newUrl = for_game.MakeVoiceVideo(bgm.GetMusicUrl(), reqMsg.GetMyVoiceUrl(), reqMsg.GetBgVolume(), reqMsg.GetMixVolume())
		if newUrl == "" {
			reqMsg.Result = easygo.NewInt32(3)
			logs.Error("音频合成失败:")
			return reqMsg
		}
	}
	//数据库写入记录
	id := for_game.InsertNewMixVideo(who.GetPlayerId(), reqMsg.GetBgId(), reqMsg.GetMixTime(), newUrl, who.GetTypes())
	if id > 0 {
		easygo.Spawn(func() {
			vcObj := for_game.GetRedisVCBuryingPointReportObj(time.Now().Unix())
			vcObj.IncrFileVal("LZMPcgOK", 1) //录音上传成功
		})
	}
	reqMsg.MixVideoUrl = easygo.NewString(newUrl)
	reqMsg.MixId = easygo.NewInt64(id)

	//作品标签引用次数加一
	for_game.IncBgVoiceTag(bgm.GetTags())
	return reqMsg
}

//修改设置玩家名片信息
func (self *cls1) RpcModifyVoiceCard(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SetVoiceCard, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcModifyVoiceCard", reqMsg, who.GetPlayerId())
	if reqMsg.GetMixId() != 0 {
		//更换录音
		who.SetMixId(reqMsg.GetMixId())
	}
	if reqMsg.GetPersonalityTags() != nil {
		//修改个性标签
		who.SetPersonalityTags(reqMsg.GetPersonalityTags())
	}
	if reqMsg.GetBgUrl() != "" {
		//更换背景
		who.SetBgImageUrl(reqMsg.GetBgUrl())
	}
	reqMsg.Result = easygo.NewInt32(1)
	logs.Info("设置语音:", reqMsg.GetMixId())
	return reqMsg
}

//获取玩家之间的亲密度:只有sayHi过的会话才能显示亲密度
func (self *cls1) RpcGetIntimacyInfo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.PlayerInfoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetIntimacyInfo", reqMsg, who.GetPlayerId())
	id := for_game.MakeSessionKey(who.GetPlayerId(), reqMsg.GetPlayerId())
	obj := for_game.GetRedisPlayerIntimacyObj(id)
	resp := &client_hall.IntimacyInfoResp{
		IsShow: easygo.NewBool(false),
	}
	if !obj.GetIsSayHi() {
		return resp
	}
	data := obj.GetRedisPlayerIntimacy()
	config := for_game.GetConfigIntimacy(data.GetIntimacyLv())
	resp = &client_hall.IntimacyInfoResp{
		Id:             easygo.NewString(id),
		IntimacyLv:     easygo.NewInt32(data.GetIntimacyLv()),
		IntimacyVal:    easygo.NewInt64(data.GetIntimacyVal()),
		IntimacyMaxVal: easygo.NewInt64(config.GetMaxVal()),
		IsShow:         easygo.NewBool(true),
	}
	if config != nil {
		resp.IntimacyMaxVal = easygo.NewInt64(config.GetMaxVal())
	}
	return resp
}

//获取所有星座信息
func (self *cls1) RpcGetAllConstellation(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetAllConstellation:", reqMsg, who.GetPlayerId())
	constellations := for_game.GetConfigConstellationFormDB()
	resp := &client_hall.ConstellationResp{
		Constellations: constellations,
	}
	//logs.Info("相应RpcGetAllConstellation:", resp)
	return resp
}

//修改设置自己的星座信息
func (self *cls1) RpcSetMyConstellation(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SetConstellation, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcSetMyConstellation", reqMsg, who.GetPlayerId())
	t := time.Now().Unix()
	//if t-who.GetConstellationTime() < 3*30*86400 {
	if t-who.GetConstellationTime() < 60 {
		return easygo.NewFailMsg("距上次修改星座未满3个月哦")
	}
	who.SetConstellation(reqMsg.GetId())
	who.SetConstellationTime(t)
	return reqMsg
}

//获取音频作品:自己的作品不需要审核
func (self *cls1) RpcGetVoiceVideo(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.GetVoiceVideoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetVoiceVideo:", reqMsg, who.GetPlayerId())
	voices := for_game.GetBgVoiceVideo(who.GetPlayerId(), reqMsg)
	var data []*client_hall.VoiceVideo
	for _, v := range voices {
		d := &client_hall.VoiceVideo{
			Id:        easygo.NewInt64(v.GetId()),
			Maker:     easygo.NewString(v.GetMaker()),
			Name:      easygo.NewString(v.GetName()),
			TagIds:    v.GetTags(),
			Content:   easygo.NewString(v.GetContent()),
			MusicUrl:  easygo.NewString(v.GetMusicUrl()),
			ImageUrl:  easygo.NewString(v.GetImageUrl()),
			Type:      easygo.NewInt32(v.GetType()),
			MusicTime: easygo.NewInt64(v.GetMusicTime()),
		}
		data = append(data, d)
	}

	resp := &client_hall.GetVoiceVideoResp{
		Type:     reqMsg.Type,
		Page:     reqMsg.Page,
		PageSize: reqMsg.PageSize,
		Data:     data,
	}
	return resp
}

//获取热门片段的标签列表
func (self *cls1) RpcGetHotEpisode(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SearchVoiceVideoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetHotEpisode: ", reqMsg, who.GetPlayerId())
	episode := make([]*client_hall.PersonTag, 0)

	hotTags, laterTags := for_game.GetVoiceTags(reqMsg.GetType(), 3, 6)
	for _, v := range hotTags {
		episode = append(episode, &client_hall.PersonTag{
			Id:    easygo.NewInt32(v.GetId()),
			Name:  easygo.NewString(v.GetName()),
			IsHot: easygo.NewBool(true),
		})
	}
	for _, v := range laterTags {
		episode = append(episode, &client_hall.PersonTag{
			Id:   easygo.NewInt32(v.GetId()),
			Name: easygo.NewString(v.GetName()),
		})
	}

	respMsg := &client_hall.HotEpisodeResp{
		Episode: episode,
	}
	return respMsg
}

//随机20条猜你喜欢作品
func (self *cls1) RpcGetMayLikeEpisode(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SearchVoiceVideoReq, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetMayLikeEpisode", reqMsg, who.GetPlayerId())
	bgVoices := for_game.GetRandBgVoice(reqMsg.GetType(), 20)
	voiceVideo := make([]*client_hall.VoiceVideo, 0)
	tagMap := make(map[int64][]int32, 0) //[作品id]标签数组
	tags := make([]int32, 0)
	for _, v := range bgVoices {
		tagMap[v.GetId()] = v.GetTags()
		tags = append(tags, v.GetTags()...)
		voiceVideo = append(voiceVideo, &client_hall.VoiceVideo{
			Id:        easygo.NewInt64(v.GetId()),
			Maker:     easygo.NewString(v.GetMaker()),
			Name:      easygo.NewString(v.GetName()),
			TagIds:    v.GetTags(),
			Content:   easygo.NewString(v.GetContent()),
			MusicUrl:  easygo.NewString(v.GetMusicUrl()),
			ImageUrl:  easygo.NewString(v.GetImageUrl()),
			Type:      easygo.NewInt32(v.GetType()),
			MusicTime: easygo.NewInt64(v.GetMusicTime()),
		})
	}
	interestTag := for_game.GetBgVideoTags(tags)
	for _, v := range voiceVideo {
		tagStr := make([]string, 0)
		for _, tag := range v.GetTagIds() {
			tagStr = append(tagStr, interestTag[tag].GetName())
		}
		v.Tags = tagStr
	}

	return &client_hall.SearchVoiceVideoResp{
		Data: voiceVideo,
	}
}

//获取该标签下的作品
func (self *cls1) RpcGetVoiceProduct(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.VoiceProduct, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetVoiceProduct: ", reqMsg, who.GetPlayerId())
	pageSize := 20 //默认取
	voices := for_game.GetTagVoiceVideo(who.GetPlayerId(), reqMsg.GetTabId(), int(reqMsg.GetPage()), pageSize)
	voiceVideo := make([]*client_hall.VoiceVideo, 0)
	tagMap := make(map[int64][]int32, 0) //[作品id]标签数组
	tags := make([]int32, 0)
	for _, v := range voices {
		tagMap[v.GetId()] = v.GetTags()
		tags = append(tags, v.GetTags()...)
		voiceVideo = append(voiceVideo, &client_hall.VoiceVideo{
			Id:        easygo.NewInt64(v.GetId()),
			Maker:     easygo.NewString(v.GetMaker()),
			Name:      easygo.NewString(v.GetName()),
			TagIds:    v.GetTags(),
			Content:   easygo.NewString(v.GetContent()),
			MusicUrl:  easygo.NewString(v.GetMusicUrl()),
			ImageUrl:  easygo.NewString(v.GetImageUrl()),
			Type:      easygo.NewInt32(v.GetType()),
			MusicTime: easygo.NewInt64(v.GetMusicTime()),
		})
	}
	interestTag := for_game.GetBgVideoTags(tags)
	for _, v := range voiceVideo {
		tagStr := make([]string, 0)
		for _, tag := range v.GetTagIds() {
			tagStr = append(tagStr, interestTag[tag].GetName())
		}
		v.Tags = tagStr
	}
	resp := &client_hall.SearchVoiceVideoResp{
		Data: voiceVideo,
	}
	return resp
}

//删除语音名片
func (self *cls1) RpcDelVoiceCard(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.DelVoiceCard, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcDelVoiceCard:", reqMsg, who.GetPlayerId())
	newId := who.DelVoiceCard(reqMsg.GetMixId())
	reqMsg.NewMixId = easygo.NewInt64(newId)
	logs.Info("RpcDelVoiceCard结果:", reqMsg)
	return reqMsg
}

//获取sayHi信息
func (self *cls1) RpcGetSessionSayHiLog(ep IGameClientEndpoint, who *Player, reqMsg *client_hall.SayHiLog, common ...*base.Common) easygo.IMessage {
	logs.Info("RpcGetSessionSayHiLog:", reqMsg)
	sayHiLog := for_game.GetOneSayHiLog(reqMsg.GetSessionId(), for_game.TALK_CONTENT_SAY_HI)
	reqMsg.Log = sayHiLog
	return reqMsg
}
