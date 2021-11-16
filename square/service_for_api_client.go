// 大厅服务器为[游戏客户端]提供的服务

package square

import (
	"encoding/base64"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/astaxie/beego/logs"
)

/*// redMsgCallClient 未读消息下发给客户端公共方法
func redMsgCallClient(ep IGameClientEndpoint, playerId for_game.PLAYER_ID) {
	unReadMsg := &client_hall.NewUnReadMessageResp{
		UnreadComment:   easygo.NewInt32(for_game.GetPlayerUnreadInfo(playerId, for_game.UNREAD_COMMENT)),
		UnreadZan:       easygo.NewInt32(for_game.GetPlayerUnreadInfo(playerId, for_game.UNREAD_ZAN)),
		UnreadAttention: easygo.NewInt32(for_game.GetPlayerUnreadInfo(playerId, for_game.UNREAD_ATTENTION)),
	}
	ep.RpcNewMessage(unReadMsg)
}*/

// 新版本的刷新社交广场内容
func (self *ServiceForHall) RpcNewVersionFlushSquareDynamic(common *base.Common, reqMsg *share_message.NewVersionFlushInfo) easygo.IMessage {
	logs.Info("=====RpcNewVersionFlushSquareDynamic======,common=%v,reqMsg=%v", common, reqMsg)
	t := reqMsg.GetType()
	playerId := common.GetUserId()
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()
	if page == 0 {
		page = int32(1)
	}
	if pageSize == 0 {
		pageSize = for_game.DYNAMIC_REQUEST_NUM
	}
	sysParam := PSysParameterMgr.GetSysParameter(for_game.SQUAREHOT_PARAMETER)
	var hotScore int32
	if sysParam != nil {
		hotScore = sysParam.GetHotScore()
	}
	// 调用通用法法.
	data := getNewVersionFlushSquareDynamic(t, page, pageSize, playerId, reqMsg.GetAdvId(), hotScore)
	msg := &client_hall.NewVersionAllInfo{
		SquareInfo:            data.SquareInfo,
		FirstAddSquareDynamic: data.FirstAddSquareDynamic,
		Type:                  easygo.NewInt32(reqMsg.GetType()),
	}
	//  异步处理动态的话题浏览量
	easygo.Spawn(func() {
		for_game.OperateViewNum(data.GetSquareInfo().GetDynamicData())
	})

	// todo ep 通过 api 发送给大厅.
	// 获取玩家,获取玩家对应的服务器            RpcNewVersionSquareAllDynamic
	SendMsgToHallClientNew([]int64{playerId}, "RpcNewVersionSquareAllDynamic", msg)
	// 通知玩家未读消息
	redMsgApiCallClient(playerId)
	return easygo.EmptyMsg
}

// 新版本的刷新社交广场内容包含话题列表.
func (self *ServiceForHall) RpcFlushSquareDynamicTopic(common *base.Common, reqMsg *share_message.FlushSquareDynamicTopicReq) easygo.IMessage {
	logs.Info("=====RpcFlushSquareDynamicTopic======,common=%v,reqMsg=%v", common, reqMsg)
	begin := for_game.GetMillSecond()
	playerId := common.GetUserId()
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()
	if page == 0 {
		page = int32(1)
	}
	if pageSize == 0 {
		pageSize = for_game.DYNAMIC_REQUEST_NUM
	}
	sysParam := PSysParameterMgr.GetSysParameter(for_game.SQUAREHOT_PARAMETER)
	var hotScore int32
	if sysParam != nil {
		hotScore = sysParam.GetHotScore()
	}

	// 调用通用法法.
	data := getNewVersionFlushSquareDynamicTopic(page, pageSize, playerId, reqMsg.GetAdvId(), hotScore)
	msg := &client_hall.NewVersionAllInfo{
		SquareInfo:            data.SquareInfo,
		FirstAddSquareDynamic: data.FirstAddSquareDynamic,
		Type:                  easygo.NewInt32(for_game.SQUARE_DYNAMIC),
	}

	//  异步处理动态的话题浏览量
	easygo.Spawn(func() {
		for_game.OperateViewNum(data.GetSquareInfo().GetDynamicData())
	})

	// todo ep 通过 api 发送给大厅.
	// 获取玩家,获取玩家对应的服务器            RpcNewVersionSquareAllDynamic
	SendMsgToHallClientNew([]int64{playerId}, "RpcNewVersionSquareAllDynamic", msg)
	// 通知玩家未读消息
	redMsgApiCallClient(playerId)
	end := for_game.GetMillSecond()
	logs.Warn("社交广场请求时间------------> %d 毫秒", end-begin)
	return easygo.EmptyMsg
}

// redMsgApiCallClient 通过api下发通知前端
func redMsgApiCallClient(playerId for_game.PLAYER_ID) {
	unReadMsg := &client_hall.NewUnReadMessageRespForApi{
		UnreadComment:   easygo.NewInt32(for_game.GetPlayerUnreadInfo(playerId, for_game.UNREAD_COMMENT)),
		UnreadZan:       easygo.NewInt32(for_game.GetPlayerUnreadInfo(playerId, for_game.UNREAD_ZAN)),
		UnreadAttention: easygo.NewInt32(for_game.GetPlayerUnreadInfo(playerId, for_game.UNREAD_ATTENTION)),
	}
	SendMsgToHallClientNew([]int64{playerId}, "RpcNewMessageForApi", unReadMsg)
}

// 发布动态
func (self *ServiceForHall) RpcAddSquareDynamic(common *base.Common, reqMsg *share_message.DynamicData) easygo.IMessage {
	logs.Info("============api 社交广场发布动态 RpcAddSquareDynamic=================,msg=%v", reqMsg)
	// 根据唯一码判断是否有动态,有,直接返回,没有,执行下面的 操作
	clientUniqueCode := reqMsg.GetClientUniqueCode()
	dynamicData, err := for_game.CheckDuplicateDynamicData(clientUniqueCode)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == nil && dynamicData.GetLogId() > 0 { // redis 存在,说明前端断网了,请求不止一次了,
		//return dynamicData
		logs.Error("重复发布动态,logId为: %d", dynamicData.GetLogId())
		return easygo.NewFailMsg("该动态重复发布")
	}

	if reqMsg.GetContent() == "" && len(reqMsg.GetPhoto()) == 0 && reqMsg.GetVoice() == "" && reqMsg.GetVideo() == "" {
		return easygo.NewFailMsg("无效的动态，什么内容都没有！！！！")
	}
	photoLst := reqMsg.GetPhoto()
	if len(photoLst) > 0 { //违规图检验
		for _, photo := range photoLst {
			num := ImageModeration(photo)
			if num != 100 {
				//return easygo.NewFailMsg("违规图片，不允许发送！")
				return easygo.NewFailMsg("动态内容包含违规信息,请重新编辑")
			}
		}
	}
	videoThumbnailURL := reqMsg.GetVideoThumbnailURL()
	if videoThumbnailURL != "" { // 缩略图
		num := ImageModeration(videoThumbnailURL)
		if num != 100 {
			logs.Error("视频缩略图违规,videoThumbnailURL: %s", videoThumbnailURL)
			//return easygo.NewFailMsg("违规图片，不允许发送！")
			return easygo.NewFailMsg("动态内容包含违规信息,请重新编辑")
		}
	}
	reqMsg.PlayerId = easygo.NewInt64(common.GetUserId())
	return addSquareDynamic(reqMsg)
}

// 删除动态
func (self *ServiceForHall) RpcDelSquareDynamic(common *base.Common, reqMsg *client_server.RequestInfo) easygo.IMessage {
	logs.Info("RpcDelSquareDynamic,msg=%v", reqMsg)
	// 异步操作话题参与数
	easygo.Spawn(
		func() {
			content := for_game.GetRedisDynamic(reqMsg.GetId()).GetContent() //动态内容.
			// 解码
			contentBytes, _ := base64.StdEncoding.DecodeString(content)
			for_game.OperateTopicParticipationNum(string(contentBytes), -1)
		})

	err := delSquareDynamic(reqMsg.GetId(), common.GetUserId())
	if err == nil {
		return easygo.EmptyMsg
	}
	return err
}

//点赞操作
func (self *ServiceForHall) RpcZanOperateSquareDynamic(common *base.Common, reqMsg *client_server.ZanInfo) easygo.IMessage {
	logs.Info("RpcZanOperateSquareDynamic", reqMsg)
	logId := reqMsg.GetLogId()
	t := reqMsg.GetType()
	if t != for_game.DYNAMIC_OPERATE && t != for_game.DYNAMIC_DELOPERATE {
		logs.Error("社交广场 点赞操作,操作类型期待的是1或者2,传过来的类型为: ", t)
		return easygo.NewFailMsg("操作类型有误")
	}
	dynamic := for_game.GetRedisDynamic(logId)
	if dynamic == nil {
		return easygo.NewFailMsg("该动态已被删除")
	}
	if for_game.GetRedisDynamicIsZan(logId, common.GetUserId()) && t == for_game.DYNAMIC_OPERATE {
		return easygo.NewFailMsg("你已经点赞过该动态")
	}
	sysp := PSysParameterMgr.GetSysParameter(for_game.PUSH_PARAMETER)
	dynamicPid := dynamic.GetPlayerId()
	b := for_game.OperateRedisDynamicZan(t, common.GetUserId(), logId, dynamicPid, sysp)
	if !b {
		return easygo.NewFailMsg("点赞操作失败")
	}

	if dynamicPid != common.GetUserId() { //自己操作  不提示红点
		redMsgApiCallClient(dynamicPid)

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
			// 修改了动态,把redis删掉
			for_game.DelSquareDynamicById(logId)
		})

		//增加每日点赞贡献度
		for _, tid := range dynamic.GetTopicId() {
			if topic := for_game.GetTopicByIdNoStatusFromDB(tid); topic != nil {
				for_game.SetDevoteLike(dynamic.GetPlayerId(), tid, topic.GetName())
			}
		}
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
	return easygo.EmptyMsg
}

//添加评论
func (self *ServiceForHall) RpcAddCommentSquareDynamic(common *base.Common, reqMsg *share_message.CommentData) easygo.IMessage {
	logs.Info("====RpcAddCommentSquareDynamic====,common=%v,msg=%v", common, reqMsg)
	logId := reqMsg.GetLogId()
	playerId := common.GetUserId()
	dynamic := for_game.GetRedisDynamic(logId)
	if dynamic == nil {
		return easygo.NewFailMsg("该动态已被删除")
	}
	belongId := reqMsg.GetBelongId()
	if belongId == 0 { //代表是主评论
		if reqMsg.GetTargetId() != reqMsg.GetOwnerId() {
			return easygo.NewFailMsg("主评论的targetId 怎么不是发动态的任务id")
		}
	} else { //非主评论
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
		redMsgApiCallClient(ownerId)
	}
	// A发布评论,B去评论,C回复B,B需要显示红点
	if belongId != 0 {
		redMsgApiCallClient(reqMsg.GetTargetId())
	}
	if belongId != 0 {
		//增加评论得分
		for_game.ModifyRedisDynamicCommentScore(belongId, 3, true)
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
		// 修改了动态,把redis删掉
		for_game.DelSquareDynamicById(logId)
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
		// 如果是自己评论,不推送
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
			m.ItemId = for_game.PUSH_ITEM_202 // 评论我的动态
			ids = for_game.GetJGIds([]int64{ownerId})
		}

		// 判断是主评论还是回复的评论
		//if belongId != 0 && reqMsg.GetTargetId() != ownerId { // 表示回复
		if belongId != 0 { // 表示回复
			m.ItemId = for_game.PUSH_ITEM_203 // 回复我的动态
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

	//设置用户每日评论获取的贡献度
	for _, tid := range dynamic.GetTopicId() {
		if topic := for_game.GetTopicByIdNoStatusFromDB(tid); topic != nil {
			for_game.SetDevoteComment(dynamic.GetPlayerId(), tid, topic.GetName())
		}
	}

	return reqMsg
}

//评论点赞/取消赞
func (self *ServiceForHall) RpcAddSquareCommentZan(common *base.Common, reqMsg *share_message.CommentDataZan) easygo.IMessage {
	logs.Info("=====RpcAddSquareCommentZan=====,common=%v,msg=%v", common, reqMsg)
	dynamicId := reqMsg.GetDynamicId()
	playerId := common.GetUserId()
	dynamic := for_game.GetRedisDynamic(dynamicId)
	if dynamic == nil {
		return easygo.NewFailMsg("该动态已被删除")
	}
	commentId := reqMsg.GetCommentId()
	comment := for_game.GetRedisDynamicComment(dynamicId, commentId)
	if comment == nil {
		return easygo.NewFailMsg("该评论不存在")
	}
	//检测是否已经点赞过
	logId := for_game.CheckIsDynamicCommentZan(playerId, dynamicId, commentId)
	if logId > 0 {
		for_game.DelDynamicCommentZan(logId)
		for_game.ModifyRedisDynamicCommentScore(commentId, -1)
		// 删除redis中我点赞过的评论id  redis_square_comment_myzanids_pid
		for_game.DelRedisPlayerMainCommentZanIds(playerId, commentId)
	} else {
		reqMsg.PlayerId = easygo.NewInt64(playerId)
		for_game.AddDynamicCommentZan(reqMsg)
		for_game.ModifyRedisDynamicCommentScore(commentId, 1)
		// 添加redis中我点赞过的评论id redis_square_comment_myzanids_pid
		for_game.AddRedisPlayerMainCommentZanIds(playerId, commentId)
	}
	return nil
}

//删除评论
func (self *ServiceForHall) RpcDelCommentSquareDynamic(common *base.Common, reqMsg *client_server.IdInfo) easygo.IMessage {
	logs.Info("RpcDelCommentSquareDynamic", reqMsg)
	logId := reqMsg.GetId()
	Id := reqMsg.GetMainId()
	dynamic := for_game.GetRedisDynamic(logId)

	comment := for_game.GetRedisDynamicComment(logId, Id)
	if comment != nil {
		ownerId := comment.GetOwnerId()
		for_game.AddPlayerUnreadInfo(ownerId, for_game.UNREAD_COMMENT, -1)

		message := for_game.GetNewUnReadMessageFromRedis(ownerId)

		result := &client_hall.NewUnReadMessageRespForApi{
			UnreadComment:   easygo.NewInt32(message.GetUnreadComment()),
			UnreadZan:       easygo.NewInt32(message.GetUnreadZan()),
			UnreadAttention: easygo.NewInt32(message.GetUnreadAttention()),
		}
		// 删除红点
		SendMsgToHallClientNew([]int64{ownerId}, "RpcNewMessageForApi", result)

		targetId := comment.GetTargetId()
		if targetId != ownerId {
			for_game.AddPlayerUnreadInfo(targetId, for_game.UNREAD_COMMENT, -1)

			message := for_game.GetNewUnReadMessageFromRedis(targetId)
			result1 := &client_hall.NewUnReadMessageRespForApi{
				UnreadComment:   easygo.NewInt32(message.GetUnreadComment()),
				UnreadZan:       easygo.NewInt32(message.GetUnreadZan()),
				UnreadAttention: easygo.NewInt32(message.GetUnreadAttention()),
			}
			SendMsgToHallClientNew([]int64{targetId}, "RpcNewMessageForApi", result1)

		}
	}

	b := for_game.DelRedisDynamicComment(logId, Id, for_game.DYNAMIC_COMMENT_STATUE_DELETE_CLIENT)
	if !b {
		return easygo.NewFailMsg("删除评论失败")
	}

	//减去用户每日评论获取的贡献度
	for _, tid := range dynamic.GetTopicId() {
		for_game.LessDevoteComment(dynamic.GetPlayerId(), tid)
	}

	//减少评论父级得分
	if comment != nil && comment.GetBelongId() != 0 {
		for_game.ModifyRedisDynamicCommentScore(comment.GetBelongId(), -3, true)
	}
	//删除评论相关点赞
	for_game.DelDynamicCommentZanByCommentId(Id)

	// 异步操作话题参与数
	easygo.Spawn(
		func() {
			content := for_game.GetRedisDynamic(logId).GetContent() //动态内容.
			// 解码
			contentBytes, _ := base64.StdEncoding.DecodeString(content)
			for_game.OperateTopicParticipationNum(string(contentBytes), -1)
		})
	return easygo.EmptyMsg
}

//关注某人
func (self *ServiceForHall) RpcAttentioPlayer(common *base.Common, reqMsg *client_server.AttenInfo) easygo.IMessage {
	logs.Info("RpcAttentioPlayer", reqMsg)
	pid := reqMsg.GetPlayerId()
	player := for_game.GetRedisPlayerBase(pid)
	who := for_game.GetRedisPlayerBase(common.GetUserId())
	t := reqMsg.GetType()
	if t == for_game.DYNAMIC_OPERATE { // 注销了的账号不给关注.
		//注销的账号不给关注
		player := for_game.GetRedisPlayerBase(pid)
		if player == nil || player.GetStatus() == for_game.ACCOUNT_CANCELED {
			return easygo.NewFailMsg("该账号异常")
		}
	}

	if t != for_game.DYNAMIC_OPERATE && t != for_game.DYNAMIC_DELOPERATE {
		logs.Error("社交广场关注操作,操作类型期待的是1或者2,传过来的值为: ", t)
		return easygo.NewFailMsg("操作类型有误")
	}
	if pid == 0 {
		return easygo.NewFailMsg("关注操作失败")
	}

	if pid == who.GetPlayerId() {
		return easygo.NewFailMsg("不能关注自己")
	}
	if util.Int64InSlice(pid, who.GetAttention()) {
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
	// 下发前端通知未读消息
	redMsgApiCallClient(pid)
	return easygo.EmptyMsg
}

func (self *ServiceForHall) RpcReadPlayerInfo(common *base.Common, reqMsg *client_hall.UnReadInfo) easygo.IMessage {
	logs.Info("RpcReadPlayerInfo", reqMsg)
	for_game.DelPlayerUnreadInfo(common.GetUserId(), reqMsg.GetType())
	return easygo.EmptyMsg
}

//获取动态详情
func (self *ServiceForHall) RpcGetDynamicInfo(common *base.Common, reqMsg *client_server.IdInfo) easygo.IMessage {
	logs.Info("========= api获取动态详情 RpcGetDynamicInfo============", reqMsg)
	logId := reqMsg.GetId()
	who := for_game.GetRedisPlayerBase(common.GetUserId())
	msg := for_game.GetRedisDynamicAllInfo(0, logId, who.GetPlayerId(), who.GetAttention())
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
func (self *ServiceForHall) RpcGetDynamicMainComment(common *base.Common, reqMsg *client_server.IdInfo) easygo.IMessage {
	logs.Info("api RpcGetDynamicMainComment", reqMsg)
	logId := reqMsg.GetId()
	mainId := reqMsg.GetMainId()
	msg := for_game.GetRedisDynamicCommentInfo(logId, mainId)
	return msg
}

//获取动态子评论
func (self *ServiceForHall) RpcGetDynamicSecondaryComment(common *base.Common, reqMsg *client_server.IdInfo) easygo.IMessage {
	logs.Info("RpcGetDynamicSecondaryComment", reqMsg)
	logId := reqMsg.GetId()
	mainId := reqMsg.GetMainId()
	secondId := reqMsg.GetSecondId()
	msg := for_game.GetRedisDynamicSecondComment(logId, mainId, secondId)
	return msg
}

//====================================新版动态评论获取=================================
//获取动态详情
func (self *ServiceForHall) RpcGetDynamicInfoNew(common *base.Common, reqMsg *client_server.IdInfo) easygo.IMessage {
	logs.Info("=============社交广场获取动态详情 RpcGetDynamicInfoNew=====================", reqMsg)
	logId := reqMsg.GetId()
	who := for_game.GetRedisPlayerBase(common.GetUserId())
	msg := for_game.GetRedisDynamicAllInfo(reqMsg.GetJumpMainCommentId(), logId, who.GetPlayerId(), who.GetAttention(), true)
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

//获取下一页动态主评论:通过页码来获取
func (self *ServiceForHall) RpcGetDynamicMainCommentNew(common *base.Common, reqMsg *client_server.IdInfo) easygo.IMessage {
	logs.Info("================社交广场获取动态评论 RpcGetDynamicMainCommentNew================", reqMsg)
	logId := reqMsg.GetId()
	params := PSysParameterMgr.GetSysParameter("squarehot_parameter")
	msg := for_game.GetRedisDynamicCommentInfoByPage(common.GetUserId(), logId, reqMsg, params)
	return msg
}

//获取动态子评论
func (self *ServiceForHall) RpcGetDynamicSecondaryCommentNew(common *base.Common, reqMsg *client_server.IdInfo) easygo.IMessage {
	logs.Info("================社交广场获取动态子评论RpcGetDynamicSecondaryCommentNew================", reqMsg)
	logId := reqMsg.GetId()
	mainId := reqMsg.GetMainId()
	secondId := reqMsg.GetSecondId()
	msg := for_game.GetRedisDynamicSecondComment(logId, mainId, secondId)
	return msg
}

//====================================新版动态评论获取=================================
//打开消息界面
func (self *ServiceForHall) RpcGetSquareMessage(common *base.Common, reqMsg *client_server.IdInfo) easygo.IMessage {
	logs.Info("====打开消息界面RpcGetSquareMessage====", reqMsg)
	id := reqMsg.GetSecondId()
	msg := for_game.GetRedisMessageAllInfo(common.GetUserId(), id)
	redMsgApiCallClient(common.GetUserId())
	return msg
}

//获取自己的被赞信息
func (self *ServiceForHall) RpcGetPlayerZanInfo(common *base.Common, reqMsg *client_server.RequestInfo) easygo.IMessage {
	logs.Info("RpcGetPlayerZanInfo", reqMsg)
	id := reqMsg.GetId()
	msg := for_game.GetRedisSquareAllZanInfo(common.GetUserId(), id)
	// 清除点赞未读列表
	for_game.DelPlayerUnreadInfo(common.GetUserId(), for_game.UNREAD_ZAN)
	redMsgApiCallClient(common.GetUserId())
	return msg
}

//获取自己的被关注信息
func (self *ServiceForHall) RpcGetPlayerAttentionInfo(common *base.Common, reqMsg *client_server.RequestInfo) easygo.IMessage {
	logs.Info("RpcGetPlayerAttentionInfo", reqMsg)
	id := reqMsg.GetId()
	msg := for_game.GetRedisSquareAttentionInfo(common.GetUserId(), id)
	redMsgApiCallClient(common.GetUserId())
	return msg
}

// 个人置顶操作.
/**
置顶:动态内容添加置顶状态(官方置顶和个人置顶两种状态). square.proto
	1.判断用户硬币是否足够
	2.判断是否超过有三条置顶了.有,返回错误不给置顶.
	3.判断该动态是否处于置顶状态.,如果是,返回:该动态已经是置顶状态了
	4.用户硬币-1
	5.修改动态为个人置顶状态.
	6.把动态置顶添加定时器.1小时后修改置顶状态.
	添加硬币流向的日志记录.
*/
func (self *ServiceForHall) RpcDynamicTop(common *base.Common, reqMsg *client_hall.DynamicTopReq) easygo.IMessage {
	logs.Info("==========个人置顶操作 RpcDynamicTop========", reqMsg)
	/*logId := reqMsg.GetLogId()
	coin := reqMsg.GetCoin()
	who := for_game.GetRedisPlayerBase(common.GetUserId())
	playerId := who.GetPlayerId()
	if logId <= 0 || coin <= 0 {
		logs.Error("个人置顶,参数有误,logId: %d,coin: %d", logId, coin)
		return easygo.NewFailMsg("参数有误")
	}
	// 判断用户硬币是否足够
	if who.GetCoin() < coin {
		return easygo.NewFailMsg("硬币不足")
	}
	//  判断置顶条数是否大于等于3条.
	if num := for_game.GetRedisDynamicTopNum(playerId); num >= for_game.DYNAMIC_TOP_NUM {
		return easygo.NewFailMsg("已存在3条置顶动态,无法置顶更多,请稍后再试")
	}
	// 判断该动态是否处于置顶状态
	dynamicData := for_game.GetRedisSquareDynamic(logId)
	if dynamicData == nil {
		return easygo.NewFailMsg("该动态不存在")
	}
	if dynamicData.GetIsBsTop() || dynamicData.GetIsTop() {
		return easygo.NewFailMsg("该动态已经是置顶状态了")
	}
	// 减用户硬币数量
	playerBase := for_game.GetRedisPlayerBase(playerId)
	if errStr := playerBase.AddCoin(0-coin, "社交广场置顶消费", for_game.COIN_TYPE_SQUARE_OUT, nil); errStr != "" {
		return easygo.NewFailMsg(errStr)
	}
	// 修改动态置顶状态.
	dynamicData.IsTop = easygo.NewBool(true)
	for_game.UpdateRedisSquareDynamic(dynamicData)
	for_game.SaveSquareDynamic(dynamicData.GetLogId())
	//===========处理置顶==============
	for_game.SetTopDynamicToRedis(for_game.REDIS_SQUARE_TOP_DYNAMIC, dynamicData) // 添加置顶动态进置顶列表
	topOverTime := time.Now().Add(time.Duration(coin) * time.Hour).Unix()         // 到期时间
	req := &share_message.BackstageNotifyTopReq{
		LogId:       easygo.NewInt64(logId),
		TopOverTime: easygo.NewInt64(topOverTime),
		IsTop:       easygo.NewBool(true),
	}
	// 判断是否是本机处理定时任务
	saveSid := for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_SQUARE_SID)
	if saveSid != PServerInfo.GetSid() { // 不在本服务器处理的,通知目标广场处理定时任务
		epSquare := SquareEpMgr.LoadEndpoint(saveSid)
		if ep1, ok := epSquare.(ISquareEndpoint); !ok {
			logs.Error("app置顶,找不到对应的广场ep,saveSid: ", saveSid)
			return nil
		} else {
			ep1.RpcSquareNotifyTop(req)
		}
		return nil
	}
	// 在同一个广场,直接处理定时任务
	if err := for_game.ProcessTopTimer(req); err != nil {
		logs.Error("动态置顶失败 RpcDynamicTop, err: ", err)
		return err
	}*/
	return nil
}

// 用户首次浏览社交广场，触发条件后弹出发布动态互动提示弹窗
func (self *ServiceForHall) RpcFirstLoginSquare(common *base.Common, reqMsg *client_hall.FirstLoginSquareReq) easygo.IMessage {
	logs.Info("======触发条件后弹出发布动态互动提示弹窗======,common=%v,reqMsg=%v", common, reqMsg) // 别删，永久留存
	player := for_game.GetRedisPlayerBase(common.GetUserId())
	if player == nil {
		return easygo.NewFailMsg("用户不存在")
	}
	if reqMsg.GetIsBrowse2Square() { // 存起来
		player.SetIsBrowse2Square(reqMsg.GetIsBrowse2Square())
	}
	// 获取用户赞的数量
	zanCount := for_game.GetRedisPlayerZanIds(common.GetUserId())
	// 获取用户评论的数量
	commentList, err := for_game.GetCommentListByPlayerIdFromDB(common.GetUserId())
	easygo.PanicError(err)
	resp := &client_hall.FirstLoginSquareReply{
		ZanCount:        easygo.NewInt64(len(zanCount)),
		CommentCount:    easygo.NewInt64(len(commentList)),
		IsBrowse2Square: easygo.NewBool(player.GetIsBrowse2Square()),
	}
	return resp
}

// 广告详情
func (self *ServiceForHall) RpcAdvDetail(common *base.Common, reqMsg *share_message.AdvSetting) easygo.IMessage {
	logs.Info("===========RpcAdvDetail 广告详情=========", reqMsg)
	if reqMsg.GetId() <= 0 {
		logs.Error("获取广告详情参数有误,id有误")
		return easygo.NewFailMsg("参数有误")
	}
	return &client_hall.AdvDetailReply{
		AdvSetting: for_game.QueryAdvToDB(reqMsg.GetId()),
		DataType:   easygo.NewInt32(1),
	}
}

// 添加广告埋点数据
func (self *ServiceForHall) RpcAddAdvLog(common *base.Common, reqMsg *share_message.AdvLogReq) easygo.IMessage {
	logs.Info("================添加广告埋点数据================", reqMsg)
	opType := reqMsg.GetOpType()
	if opType != for_game.ADV_LOG_OP_TYPE_1 && opType != for_game.ADV_LOG_OP_TYPE_3 {
		logs.Error("广场广告埋点数据,操作类型有误,opType: ", opType)
		return easygo.NewFailMsg("操作类型有误")
	}
	reqMsg.PlayerId = easygo.NewInt64(common.GetUserId())
	reqMsg.OpTime = easygo.NewInt64(time.Now().Unix())

	// 前端只传1-展示,3-点击
	switch opType {
	case for_game.ADV_LOG_OP_TYPE_1: // 展示次数+1,判断是否需要去重.
		id := for_game.NextId(for_game.TABLE_ADV_LOG)
		id2 := for_game.NextId(for_game.TABLE_ADV_LOG)
		reqMsg.Id = easygo.NewInt64(id)
		for_game.AddAdvLogToDB(reqMsg)
		// 展示人数去重
		if log := for_game.GetAdvLogByPidAndOpFromDB(common.GetUserId(), reqMsg.GetAdvId(), for_game.ADV_LOG_OP_TYPE_2); log != nil {
			return easygo.EmptyMsg
		}
		reqMsg.Id = easygo.NewInt64(id2)
		reqMsg.OpType = easygo.NewInt32(for_game.ADV_LOG_OP_TYPE_2)
		for_game.AddAdvLogToDB(reqMsg)
		return easygo.EmptyMsg
	case for_game.ADV_LOG_OP_TYPE_3:
		id := for_game.NextId(for_game.TABLE_ADV_LOG)
		id2 := for_game.NextId(for_game.TABLE_ADV_LOG)
		reqMsg.Id = easygo.NewInt64(id)
		for_game.AddAdvLogToDB(reqMsg)
		// 点击人数去重
		if log := for_game.GetAdvLogByPidAndOpFromDB(common.GetUserId(), reqMsg.GetAdvId(), for_game.ADV_LOG_OP_TYPE_4); log != nil {
			return easygo.EmptyMsg
		}
		reqMsg.Id = easygo.NewInt64(id2)
		reqMsg.OpType = easygo.NewInt32(for_game.ADV_LOG_OP_TYPE_4)
		for_game.AddAdvLogToDB(reqMsg)
		return easygo.EmptyMsg
	}
	return easygo.EmptyMsg
}

// 带话题的社交广场关注页的请求.
func (self *ServiceForHall) RpcSquareAttention(common *base.Common, reqMsg *client_hall.SquareAttentionReq) easygo.IMessage {
	logs.Info("======带话题的社交广场关注页的请求 RpcSquareAttention=====,common=%v,reqMsg=%v", common, reqMsg)
	begin := for_game.GetMillSecond()
	sysParam := PSysParameterMgr.GetSysParameter(for_game.SQUAREHOT_PARAMETER)
	var hotScore int32
	if sysParam != nil {
		hotScore = sysParam.GetHotScore()
	}
	resp := for_game.GetSquareAttentionData(common.GetUserId(), hotScore, reqMsg)
	//  异步处理动态的话题浏览量
	if len(resp.GetDynamicList()) > 0 {
		easygo.Spawn(func() {
			for_game.OperateViewNum(resp.GetDynamicList())
		})
	}
	end := for_game.GetMillSecond()
	logs.Info("关注页请求时间-------------> %d 毫秒", end-begin)
	return resp
}
func ImageModeration(url string) int32 {
	res := for_game.ImageModeration(url)
	if res == nil {
		return 0
	}
	logs.Info(" 图片检测结果:", res)
	if res.EvilFlag != 0 || res.EvilType != 100 {
		isGet := false
		score := PSysParameterMgr.GetImageModeration(res.EvilType)
		switch res.EvilType {
		case 20001:
			if res.PolityDetect.Score > score {
				isGet = true
			}
		case 20002:
			if res.PornDetect.Score > score {
				isGet = true
			}
		case 20006:
			if res.IllegalDetect.Score > score {
				isGet = true
			}
		case 20103:
			if res.HotDetect.Score > score {
				isGet = true
			}
		case 24001:
			if res.TerrorDetect.Score > score {
				isGet = true
			}
		}
		if isGet {
			return res.EvilType
		}
	}
	return 100
}
