/**
服务器间,客户端与服务器间的通用的获取数据的方法.
*/
package square

import (
	"encoding/base64"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_server"
	"game_server/pb/share_message"

	"github.com/astaxie/beego/logs"
)

// 新版获取刷新动态内容   hotScore:动态热门分数
func getNewVersionFlushSquareDynamic(t, page, pageSize int32, playerId int64, advId int64, hotScore int32) *share_message.NewVersionAllInfo {
	dataPage := for_game.GetRedisNewDynamic1(t, page, pageSize, playerId, hotScore, advId)
	// 判断自己是否发布过动态
	var firstAddSquareDynamic bool
	if !for_game.GetRedisPlayerBase(playerId).GetFirstAddSquareDynamic() {
		firstAddSquareDynamic = true
	}
	msg := &share_message.NewVersionAllInfo{
		FirstAddSquareDynamic: easygo.NewBool(firstAddSquareDynamic),
		SquareInfo:            dataPage,
	}
	return msg
}
func getNewVersionFlushSquareDynamicTopic(page, pageSize int32, playerId int64, advId int64, hotScore int32) *share_message.NewVersionAllInfo {
	dataPage := for_game.GetRedisNewDynamicTopic(page, pageSize, playerId, hotScore, advId)
	// 判断自己是否发布过动态
	var firstAddSquareDynamic bool
	if !for_game.GetRedisPlayerBase(playerId).GetFirstAddSquareDynamic() {
		firstAddSquareDynamic = true
	}
	msg := &share_message.NewVersionAllInfo{
		FirstAddSquareDynamic: easygo.NewBool(firstAddSquareDynamic),
		SquareInfo:            dataPage,
	}
	return msg
}

func addSquareDynamic(reqMsg *share_message.DynamicData) easygo.IMessage {
	sendTime := easygo.NowTimestamp()
	logId := for_game.NextId(for_game.TABLE_SQUARE_DYNAMIC)
	reqMsg.LogId = easygo.NewInt64(logId)
	reqMsg.CreateTime = easygo.NewInt64(sendTime)
	reqMsg.SendTime = easygo.NewInt64(sendTime)
	reqMsg.Statue = easygo.NewInt32(for_game.DYNAMIC_STATUE_COMMON)
	content := reqMsg.GetContent()

	//解析动态内容,抽取话题
	if topicIdList := for_game.ParseContentTopic(content, 1); topicIdList != nil {
		reqMsg.TopicId = topicIdList
	}

	//保存话题信息
	topicArr := make(map[int64]*share_message.Topic)
	for _, tid := range reqMsg.GetTopicId() {
		if topic := for_game.GetTopicByIdNoStatusFromDB(tid); topic != nil {
			if topic.GetStatus() == for_game.TOPIC_STATUS_CLOSE {
				logs.Error("发布动态,话题 %s 已关闭状态", topic.GetName())
				msg := fmt.Sprintf("%s %s", topic.GetName(), "涉嫌违规,暂停使用")
				return easygo.NewFailMsg(msg)
			}
			topicArr[tid] = topic
		}
	}

	text := base64.StdEncoding.EncodeToString([]byte(content))
	reqMsg.Content = easygo.NewString(text)
	for_game.AddRedisSquareDynamic(reqMsg)

	// 获取用户
	who := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	who.AddRedisPlayerDynamicList(logId)

	msg := &client_server.RequestInfo{
		Id: easygo.NewInt64(logId),
	}

	// 唯一码 和动态对应的记录在数据库.
	for_game.SaveDuplicateDynamicDataToRedis(reqMsg.GetClientUniqueCode(), reqMsg)
	// 设置初次发布动态
	for_game.GetRedisPlayerBase(who.GetPlayerId()).SetFirstAddSquareDynamic(true)

	// 异步通知粉丝我发布了动态
	fun := func(req *share_message.DynamicData) {
		p := for_game.GetRedisPlayerBase(req.GetPlayerId())
		if p == nil {
			return
		}
		fans := p.GetFans()
		if len(fans) == 0 {
			return
		}
		fid := make([]int64, 0)
		//判断粉丝是否开启了设置
		for _, v := range fans {
			f := for_game.GetRedisPlayerBase(v)
			if f == nil {
				continue
			}
			if f.GetIsOpenSquare() {
				continue
			}
			fid = append(fid, v)
		}
		if len(fid) == 0 {
			return
		}

		ids := for_game.GetJGIds(fid)
		m := for_game.PushMessage{
			//Title:       "温馨提示",
			ContentType: for_game.JG_TYPE_SQUARE,
			Location:    7,                      // 社交广场
			ItemId:      for_game.PUSH_ITEM_204, // 我关注的人发布新动态
			ObjectId:    logId,
		}
		sysp := PSysParameterMgr.GetSysParameter(for_game.PUSH_PARAMETER)

		for _, ps := range sysp.GetPushSet() {
			if ps.GetObjId() == m.ItemId && ps.GetIsPush() {
				m.Content = fmt.Sprintf(ps.GetObjContent(), p.GetNickName())
			}
		}
		for_game.JGSendMessage(ids, m, sysp)
	}
	easygo.Spawn(fun, reqMsg)

	//发布话题动态成功，增加贡献度
	for _, tid := range reqMsg.GetTopicId() {
		topicInfo := topicArr[tid]
		for_game.SetDevoteDynamic(reqMsg.GetPlayerId(), tid, topicInfo.GetName())
	}

	return msg
}

func delSquareDynamic(logId, playerId int64) *base.Fail {
	b := for_game.DelRedisSquareDynamic(playerId, logId, "")
	if !b {
		return easygo.NewFailMsg("删除动态失败")
	}
	dynamicData, _ := for_game.GetDynamicByLogIdFromDB(logId)
	//删除话题动态，扣除对应的贡献度
	for _, tid := range dynamicData.GetTopicId() {
		for_game.LessDevoteDynamic(dynamicData.GetPlayerId(), tid, logId)
	}
	who := for_game.GetRedisPlayerBase(playerId)
	who.DelRedisPlayerDynamicList(logId)
	return nil
}
