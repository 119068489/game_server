// 大厅服务器为[游戏客户端]提供的服务

package square

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"unicode/utf8"
)

// 全部话题(官方)
func (self *ServiceForHall) RpcGetAllTopic(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("====全部话题 RpcGetAllTopic=====,common=%v", common)
	// 查找官方话题类型
	bsTopicTypes := for_game.GetBSTopicTypeListByClassFormDB(for_game.TOPIC_CLASS_BS)
	topicTypeList := make([]*client_hall.OneTopicType, 0)
	for _, topicType := range bsTopicTypes {
		topicList := for_game.GetBSTopicTypeListByTypeIdFormDB(topicType.GetId())
		topicTypeList = append(topicTypeList, &client_hall.OneTopicType{
			TopicType: topicType,
			TopicList: topicList,
		})
	}
	resp := &client_hall.AllTopic{
		TopicTypeList: topicTypeList,
	}
	return resp
}

//  话题类别列表(不展示用户自定义的)
func (self *ServiceForHall) RpcGetTopicTypeList(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("====话题类别列表(不展示用户自定义的) RpcGetTopicTypeList=====,common=%v,reqMst=%v", common, reqMsg)
	resp := &client_hall.TopicTypeListResp{
		TopicTypeList: for_game.GetBSTopicTypeListByClassFormDB(for_game.TOPIC_CLASS_BS),
	}
	return resp
}

// 根据话题类别查找话题
func (self *ServiceForHall) RpcGetTopicList(common *base.Common, reqMsg *client_hall.TopicListReq) easygo.IMessage {
	logs.Info("====根据话题类别查找话题 RpcGetTopicList=====,common=%v,reqMst=%v", common, reqMsg)
	topics, count := for_game.GetTopicListByTypeIdPage(common.GetUserId(), reqMsg.GetTopicTypeId(), reqMsg.GetPage(), reqMsg.GetPageSize())
	resp := &client_hall.TopicListResp{
		TopicList: topics,
		Count:     easygo.NewInt64(count),
	}
	return resp
}

// 话题主页头部详细信息
func (self *ServiceForHall) RpcGetTopicDetailReq(common *base.Common, reqMsg *client_hall.TopicDetailReq) easygo.IMessage {
	logs.Info("====话题主页头部详细信息 RpcGetTopicDetailReq=====,common=%v,reqMst=%v", common, reqMsg)
	var topic *share_message.Topic
	if reqMsg.GetId() == 0 { // 通过名字取话题
		topic = for_game.GetTopicByNameNoStatusFromDB(reqMsg.GetName())
		if topic == nil {
			logs.Error("话题主页头部详细信息 RpcGetTopicDetailReq,该话题不存在, topicName: ", reqMsg.GetName())
			return easygo.NewFailMsg("该话题不存在")
		}
		reqMsg.Id = easygo.NewInt64(topic.GetId())
	}
	topic = for_game.GetTopicByIdNoStatusFromDB(reqMsg.GetId())
	if topic == nil {
		logs.Error("话题主页头部详细信息 RpcGetTopicDetailReq,该话题不存在, topicId: ", reqMsg.GetId())
		return easygo.NewFailMsg("该话题不存在")
	}
	if topic.GetStatus() == for_game.TOPIC_STATUS_CLOSE {
		logs.Error("违规话题整改中,状态为已关闭,id为: ", reqMsg.GetId())
		return easygo.NewFailMsg("违规话题整改中")
	}

	// 判断玩家是否关注了这个话题
	for_game.CheckPlayerIsAttentionTopic(common.GetUserId(), topic)

	//获取话题主头像
	topicMasterInfo := new(share_message.TeamPlayerInfo)
	if topic.GetTopicMaster() != 0 {
		playerInfo, err := for_game.GetPlayerInfo(topic.GetTopicMaster())
		if err != nil {
			logs.Error(err.Error(), topic.GetTopicMaster())
			return easygo.NewFailMsg(err.Error())
		}
		topicMasterInfo.PlayerId = easygo.NewInt64(playerInfo.GetPlayerId())
		topicMasterInfo.HeadIcon = easygo.NewString(playerInfo.GetHeadIcon())
		topicMasterInfo.Sex = easygo.NewInt32(playerInfo.GetSex())
	}

	// 关联的话题.
	relationTopicList := for_game.GetRelateTopicByTypeId(topic.GetTopicTypeId(), reqMsg.GetId())
	resp := &client_hall.TopicDetailResp{
		Topic:           topic,
		RelatedTopics:   relationTopicList,
		TopicMasterInfo: topicMasterInfo,
	}

	// 添加话题浏览数
	easygo.Spawn(func() { for_game.IncTopicViewingNumToDB(reqMsg.GetId(), 1) })

	return resp
}

// 话题主页动态列表
func (self *ServiceForHall) RpcGetTopicMainPageList(common *base.Common, reqMsg *client_hall.TopicMainPageListReq) easygo.IMessage {
	logs.Info("====话题主页动态列表 RpcGetTopicMainPageList=====,common=%v,reqMst=%v", common, reqMsg)
	topicId := reqMsg.GetId()
	if topicId == 0 { // 用名字获取
		t := for_game.GetTopicByNameNoStatusFromDB(reqMsg.GetName())
		if t == nil {
			logs.Error("话题主页动态列表 RpcGetTopicMainPageList,该话题不存在, topicName: ", reqMsg.GetName())
			return easygo.NewFailMsg("该话题不存在")
		}
		topicId = t.GetId()
	}
	// 判断该话题是否存在
	topic := for_game.GetTopicByIdNoStatusFromDB(topicId)
	if topic == nil {
		logs.Error("话题主页动态列表 RpcGetTopicMainPageList,该话题不存在, topicId: ", topicId)
		return easygo.NewFailMsg("该话题不存在")
	}
	if topic.GetStatus() == for_game.TOPIC_STATUS_CLOSE {
		logs.Error("违规话题整改中,状态为已关闭,id为: ", reqMsg.GetId())
		return easygo.NewFailMsg("违规话题整改中")
	}
	reqType := reqMsg.GetReqType() // 请求类型,1-最新,2-热门
	if reqType != for_game.REQ_TYPE_NEW && reqType != for_game.REQ_TYPE_HOT {
		logs.Error("话题主页动态列表 RpcGetTopicMainPageList 请求的类型有误,不是热门也不是最新,reqType :", reqType)
		return easygo.NewFailMsg("请求的类型有误")
	}
	// 获取热门得分
	sysParam := PSysParameterMgr.GetSysParameter(for_game.SQUAREHOT_PARAMETER)
	var hotScore int32
	if sysParam != nil {
		hotScore = sysParam.GetHotScore()
	}

	dynamicList, count := for_game.GetDynamicByTopicId(reqType, hotScore, common.GetUserId(), topicId, reqMsg.GetPage(), reqMsg.GetPageSize())
	logs.Info("dynamicList:", dynamicList)
	logs.Info("count:", count)
	resp := &client_hall.TopicMainPageListResp{
		DynamicList:  dynamicList,
		DynamicCount: easygo.NewInt64(count),
	}
	return resp
}

// 某话题参与详情列表请求(推荐用户栏)
func (self *ServiceForHall) RpcGetTopicParticipateList(common *base.Common, reqMsg *client_hall.TopicParticipateListReq) easygo.IMessage {
	logs.Info("====某话题参与详情列表请求 RpcGetTopicParticipateList=====,common=%v,reqMst=%v", common, reqMsg)
	// 判断该话题是否存在
	topicId := reqMsg.GetId()
	if topicId == 0 { // 用名字获取
		t := for_game.GetTopicByNameFromDB(reqMsg.GetName())
		if t == nil {
			logs.Error("话题主页动态列表 RpcGetTopicMainPageList,该话题不存在, topicName: ", reqMsg.GetName())
			return easygo.NewFailMsg("该话题不存在")
		}
		topicId = t.GetId()
	}
	topic := for_game.GetTopicByIdFromDB(topicId)
	if topic == nil {
		logs.Error("某话题参与详情列表请求 RpcGetTopicParticipateList,该话题不存在, topicId: ", topicId)
		return easygo.NewFailMsg("该话题不存在")
	}
	// 根据动态热门分排序
	players, count := for_game.GetPlayerListByTopicId(topicId, reqMsg.GetPage(), reqMsg.GetPageSize())
	resp := &client_hall.TopicParticipateListResp{
		PlayerList:  players,
		PlayerCount: easygo.NewInt64(count),
	}

	bytes, _ := json.Marshal(resp)
	logs.Info("----->", string(bytes))

	return resp
}

// 关注或取消关注话题
func (self *ServiceForHall) RpcAttentionTopic(common *base.Common, reqMsg *client_hall.AttentionTopicReq) easygo.IMessage {
	logs.Info("====关注或取消关注话题 RpcAttentionTopic=====,common=%v,reqMst=%v", common, reqMsg)
	// 判断用户是否存在
	pid := common.GetUserId()
	topicIds := reqMsg.GetId()
	if playerBase := for_game.GetRedisPlayerBase(pid); playerBase == nil {
		logs.Error("关注或取消关注话题 RpcAttentionTopic,用户不存在,pid: ", pid)
		return easygo.NewFailMsg("用户不存在")
	}
	// 判断该话题是否存在
	for _, topicId := range topicIds {
		topic := for_game.GetTopicByIdFromDB(topicId)
		if topic == nil {
			logs.Error("关注或取消关注话题 RpcAttentionTopic,该话题不存在, topicId: ", topicId)
			return easygo.NewFailMsg("该话题不存在")
		}

	}
	// 判断操作类型
	operate := reqMsg.GetOperate()
	if operate != for_game.OPERATE_TOPIC_ATTENTION && operate != for_game.OPERATE_TOPIC_CANCEL_ATTENTION {
		logs.Error("关注或取消关注话题 RpcAttentionTopic,操作类型有误, operate: ", reqMsg.GetOperate())
		return easygo.NewFailMsg("不存在的操作类型")
	}

	if !for_game.OperateTopic(operate, pid, topicIds) {
		return easygo.NewFailMsg("操作失败")
	}
	return nil
}

// 话题我的关注列表
func (self *ServiceForHall) RpcMyAttentionTopicList(common *base.Common, reqMsg *client_hall.MyAttentionTopicListReq) easygo.IMessage {
	logs.Info("====话题我的关注列表 RpcMyAttentionTopicList=====,common=%v,reqMst=%v", common, reqMsg)
	// 判断用户是否存在
	pid := common.GetUserId()
	if player := for_game.GetRedisPlayerBase(pid); player == nil {
		logs.Error("话题我的关注列表,用户不存在,pid: ", pid)
		return easygo.NewFailMsg("用户不存在")
	}
	topicList, count := for_game.GetPlayerAttentionList(pid, reqMsg.GetPage(), reqMsg.GetPageSize())
	resp := &client_hall.MyAttentionTopicListResp{
		TopicList: topicList,
		Count:     easygo.NewInt64(count),
	}
	return resp
}

// 模糊查询话题
func (self *ServiceForHall) RpcSearchTopic(common *base.Common, reqMsg *client_hall.SearchTopicReq) easygo.IMessage {
	logs.Info("====模糊查询话题 RpcSearchTopic=====,common=%v,reqMst=%v", common, reqMsg)
	name := reqMsg.GetName()
	if name == "" {
		return nil
	}

	resp := &client_hall.SearchTopicResp{
		TopicList: for_game.SearchTopic(name),
	}
	return resp
}

// 添加话题页的热门推荐,不包含自定义
func (self *ServiceForHall) RpcSearchHotTopic(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("====添加话题页的热门推荐,不包含自定义 RpcSearchHotTopic=====,common=%v,reqMst=%v", common, reqMsg)
	resp := &client_hall.SearchHotTopicResp{
		TopicList: for_game.SearchHotTopic(),
	}
	return resp
}

// 话题,广场旁边的按钮.
func (self *ServiceForHall) RpcFlushTopic(common *base.Common, reqMsg *client_hall.FlushTopicReq) easygo.IMessage {
	logs.Info("====话题,广场旁边的按钮 RpcFlushTopic=====,common=%v,reqMst=%v", common, reqMsg)
	lists, count := for_game.FlushTopic(reqMsg.GetPage(), reqMsg.GetPageSize())
	resp := &client_hall.FlushTopicResp{
		TopicDynamicList: lists,
		Count:            easygo.NewInt64(count),
	}
	return resp
}

// 广场按钮旁边的话题头部话题列表
func (self *ServiceForHall) RpcTopicHead(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("====广场按钮旁边的话题头部话题列表 RpcTopicHead=====,common=%v,reqMst=%v", common, reqMsg)
	topics := for_game.GetTopicHeadTopic()
	resp := &client_hall.TopicHeadResp{
		TopicList: topics,
	}
	return resp
}

// 热门话题列表(官方)
func (self *ServiceForHall) RpcHotTopicList(common *base.Common, reqMsg *client_hall.HotTopicListReq) easygo.IMessage {
	logs.Info("====热门话题列表 RpcHotTopicList=====,common=%v,reqMst=%v", common, reqMsg)
	topics, count := for_game.GetHotTopicList(common.GetUserId(), reqMsg.GetPage(), reqMsg.GetPageSize())
	resp := &client_hall.HotTopicListResp{
		TopicList: topics,
		Count:     easygo.NewInt64(count),
	}
	return resp
}

// 关注页推荐用户
func (self *ServiceForHall) RpcAttentionRecommendPlayer(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("====热门话题列表 RpcAttentionRecommendPlayer=====,common=%v,reqMst=%v", common, reqMsg)
	playerList := for_game.GetAttentionRecommendPlayer()
	resp := &client_hall.AttentionRecommendPlayerResp{
		PlayerList: playerList,
	}
	return resp
}

// 话题主页,热门动态用户推荐栏
func (self *ServiceForHall) RpcTopicHotDynamicParticipatePlayer(common *base.Common, reqMsg *client_hall.TopicParticipateListReq) easygo.IMessage {
	logs.Info("====话题主页,热门动态用户推荐栏 RpcTopicHotDynamicParticipatePlayer=====,common=%v,reqMst=%v", common, reqMsg)
	// 判断该话题是否存在
	topicId := reqMsg.GetId()
	if topicId == 0 { // 用名字获取
		t := for_game.GetTopicByNameNoStatusFromDB(reqMsg.GetName())
		if t == nil {
			logs.Error("话题主页,热门动态用户推荐栏 RpcTopicHotDynamicParticipatePlayer,该话题不存在, topicName: ", reqMsg.GetName())
			return easygo.NewFailMsg("该话题不存在")
		}
		topicId = t.GetId()
	}
	topic := for_game.GetTopicByIdNoStatusFromDB(topicId)
	if topic == nil {
		logs.Error("话题主页,热门动态用户推荐栏 RpcTopicHotDynamicParticipatePlayer,该话题不存在, topicId: ", topicId)
		return easygo.NewFailMsg("该话题不存在")
	}
	if topic.GetStatus() == for_game.TOPIC_STATUS_CLOSE {
		logs.Error("违规话题整改中,状态为已关闭,id为: ", reqMsg.GetId())
		return easygo.NewFailMsg("违规话题整改中")
	}
	// 动态热门分
	sysParam := PSysParameterMgr.GetSysParameter(for_game.SQUAREHOT_PARAMETER)
	var hotScore int32
	if sysParam != nil {
		hotScore = sysParam.GetHotScore()
	}
	// 根据动态热门分排序
	players, count := for_game.GetDevicePlayerHotDynamicByTopicId(topicId, reqMsg.GetPage(), reqMsg.GetPageSize(), hotScore)
	resp := &client_hall.TopicParticipateListResp{
		PlayerList:  players,
		PlayerCount: easygo.NewInt64(count),
	}

	return resp
}

// 话题贡献榜
func (self *ServiceForHall) RpcTopicDevoteList(common *base.Common, reqMsg *client_hall.TopicDevoteListReq) easygo.IMessage {
	logs.Info("====话题贡献榜 RpcTopicDevoteList=====,common=%v,reqMst=%v", common, reqMsg)
	data := make([]*share_message.TopicDevote, 0)
	var count int
	//日榜
	if reqMsg.GetDataType() == 1 {
		data, count = for_game.GetTopicPlayerDevoteDayList(reqMsg.GetTopicName(), reqMsg.GetTopicId(), int(reqMsg.GetPage()), int(reqMsg.GetPageSize()))
	}

	//月榜
	if reqMsg.GetDataType() == 2 {
		data, count = for_game.GetTopicPlayerDevoteMonthList(reqMsg.GetTopicName(), reqMsg.GetTopicId(), int(reqMsg.GetPage()), int(reqMsg.GetPageSize()))
	}

	//总榜
	if reqMsg.GetDataType() == 3 {
		data, count = for_game.GetTopicPlayerDevoteTotalList(reqMsg.GetTopicName(), reqMsg.GetTopicId(), int(reqMsg.GetPage()), int(reqMsg.GetPageSize()))
	}
	resp := &client_hall.TopicDevoteListResp{
		DevoteList: data,
		Count:      easygo.NewInt64(count),
	}
	return resp
}

// 获取申请话题主条件
func (self *ServiceForHall) RpcTopicMasterCondition(common *base.Common, reqMsg *client_hall.TopicMasterConditionReq) easygo.IMessage {
	logs.Info("====获取申请话题主条件 RpcTopicMasterCondition=====,common=%v,reqMst=%v", common, reqMsg)
	playerId := common.GetUserId()
	topicId := reqMsg.GetTopicId()
	playerRegisteredCondition := for_game.GetPlayerRegisteredCondition(playerId)

	follow := false
	playerAttentionTopic := for_game.GetPlayerAttentionFromDB(topicId, playerId)
	if playerAttentionTopic.GetId() > 0 {
		follow = true
	} else {
		follow = false
	}

	dynamicBool := false
	playerTopicDynamicCount := for_game.GetPlayerTopicDynamicCount(playerId, topicId)
	if playerTopicDynamicCount >= 10 {
		dynamicBool = true
	}

	generalStatus := false
	if playerRegisteredCondition && follow && dynamicBool {
		generalStatus = true
	}

	resp := &client_hall.TopicMasterConditionResp{
		Registered:    easygo.NewBool(playerRegisteredCondition),
		Follow:        easygo.NewBool(follow),
		Dynamic:       easygo.NewBool(dynamicBool),
		GeneralStatus: easygo.NewBool(generalStatus),
	}
	return resp
}

// 申请话题主
func (self *ServiceForHall) RpcApplyTopicMaster(common *base.Common, reqMsg *client_hall.ApplyTopicMasterReq) easygo.IMessage {
	logs.Info("====申请话题主 RpcTopicMasterCondition=====,common=%v,reqMst=%v", common, reqMsg)
	reasonCount := utf8.RuneCountInString(reqMsg.GetReason())
	contactDetailsCount := utf8.RuneCountInString(reqMsg.GetContactDetails())
	if reasonCount > 200 {
		logs.Error("申请理由过长", common.GetUserId())
		return easygo.NewFailMsg("申请理由过长")
	}
	if contactDetailsCount > 20 {
		logs.Error("联系方式过长", common.GetUserId())
		return easygo.NewFailMsg("联系方式过长")
	}
	playerInfo, err := for_game.GetPlayerInfo(common.GetUserId())
	if err != nil {
		logs.Error(err.Error(), common.GetUserId())
		return easygo.NewFailMsg(err.Error())
	}
	data := &share_message.ApplyTopicMaster{
		TopicId:        easygo.NewInt64(reqMsg.GetTopicId()),
		PlayerId:       easygo.NewInt64(common.GetUserId()),
		IsManageExp:    easygo.NewBool(reqMsg.GetIsManageExp()),
		Reason:         easygo.NewString(reqMsg.GetReason()),
		ContactDetails: easygo.NewString(reqMsg.GetContactDetails()),
		TopicName:      easygo.NewString(reqMsg.GetTopicName()),
		PlayerAccount:  easygo.NewString(playerInfo.GetAccount()),
	}
	err = for_game.AddApplyTopicMaster(data)
	if err != nil {
		logs.Error("写入失败", common.GetUserId())
		return easygo.NewFailMsg("写入失败")
	}

	resp := &client_hall.ApplyTopicMasterResp{
		Result: easygo.NewInt32(1),
	}
	return resp
}

//话题主修改话题信息
func (self *ServiceForHall) RpcTopicMasterEdit(common *base.Common, reqMsg *client_hall.TopicMasterEditReq) easygo.IMessage {
	logs.Info("====话题主修改话题信息 RpcTopicMasterEdit=====,common=%v,reqMst=%v", common, reqMsg)

	isTopicMaster := for_game.IsTopicMaster(reqMsg.GetTopicId(), common.GetUserId())
	if !isTopicMaster {
		logs.Error("只有话题主可以修改", common.GetUserId())
		return easygo.NewFailMsg("只有话题主可以修改")
	}
	playerInfo, err := for_game.GetPlayerInfo(common.GetUserId())
	if err != nil {
		logs.Error(err.Error(), common.GetUserId())
		return easygo.NewFailMsg(err.Error())
	}
	data := &share_message.ApplyEditTopicInfo{
		TopicId:       easygo.NewInt64(reqMsg.GetTopicId()),
		HeadURL:       easygo.NewString(reqMsg.GetHeadURL()),
		Description:   easygo.NewString(reqMsg.GetDescription()),
		BgUrl:         easygo.NewString(reqMsg.GetBgUrl()),
		TopicRule:     easygo.NewString(reqMsg.GetTopicRule()),
		TopicName:     easygo.NewString(reqMsg.GetTopicName()),
		PlayerAccount: easygo.NewString(playerInfo.GetAccount()),
	}
	err = for_game.EditTopicInfo(data)
	if err != nil {
		logs.Error("写入失败", common.GetUserId())
		return easygo.NewFailMsg("写入失败")
	}

	resp := &client_hall.TopicMasterEditResp{
		Result: easygo.NewInt32(1),
	}
	return resp
}

//话题置顶
func (self *ServiceForHall) RpcTopicTop(common *base.Common, reqMsg *client_hall.TopicTopReq) easygo.IMessage {
	logs.Info("====话题置顶 RpcTopicTop=====,common=%v,reqMst=%v", common, reqMsg)

	var topicId int64
	if reqMsg.GetTopicId() != 0 {
		topicId = reqMsg.GetTopicId()
	} else {
		topicInfo := for_game.GetTopicByNameFromDB(reqMsg.GetTopicName())
		topicId = topicInfo.GetId()
	}

	isTopicMaster := for_game.IsTopicMaster(topicId, common.GetUserId())
	if !isTopicMaster {
		logs.Error("只有话题主可以置顶动态", common.GetUserId())
		return easygo.NewFailMsg("只有话题主可以置顶动态")
	}

	dynamic := for_game.GetRedisDynamic(reqMsg.GetLogId())
	if dynamic == nil {
		logs.Error("该动态不存在", reqMsg.GetLogId())
		return easygo.NewFailMsg("该动态不存在")
	}

	topicTopSet := dynamic.GetTopicTopSet()
	topicTop := &share_message.TopicTop{
		TopicId:          easygo.NewInt64(topicId),
		IsTopicTop:       easygo.NewBool(true),
		TopicTopOverTime: easygo.NewInt64(-1),
		TopicTopTime:     easygo.NewInt64(util.GetTime()),
	}
	topicTopSet = append(topicTopSet, topicTop)
	dynamic.TopicTopSet = topicTopSet
	for_game.UpdateRedisSquareDynamic(dynamic)
	for_game.SaveSquareDynamic(dynamic.GetLogId())

	resp := &client_hall.TopicTopResp{
		Result: easygo.NewInt32(1),
	}

	//推送小助手
	content := fmt.Sprintf("%s\n您发布的动态被话题主推荐啦，快去看看吧~", util.GetTimestamp())
	pushReq := &server_server.TopicDynamicTopLittleHelper{
		PlayerId: easygo.NewInt64(dynamic.GetPlayerId()),
		Title:    easygo.NewString("社交广场通知"),
		Content:  easygo.NewString(content),
	}
	pushTopicDynamicTopStatus(pushReq)

	return resp
}

//取消话题动态置顶
func (self *ServiceForHall) RpcTopicTopCancel(common *base.Common, reqMsg *client_hall.TopicTopCancelReq) easygo.IMessage {
	logs.Info("====取消话题置顶 RpcTopicTopCancel=====,common=%v,reqMst=%v", common, reqMsg)

	var topicId int64
	if reqMsg.GetTopicId() != 0 {
		topicId = reqMsg.GetTopicId()
	} else {
		topicInfo := for_game.GetTopicByNameFromDB(reqMsg.GetTopicName())
		topicId = topicInfo.GetId()
	}

	isTopicMaster := for_game.IsTopicMaster(topicId, common.GetUserId())
	if !isTopicMaster {
		logs.Error("只有话题主可以取消置顶动态", common.GetUserId())
		return easygo.NewFailMsg("只有话题主可以取消置顶动态")
	}

	dynamic := for_game.GetRedisDynamic(reqMsg.GetLogId())
	if dynamic == nil {
		logs.Error("该动态不存在", reqMsg.GetLogId())
		return easygo.NewFailMsg("该动态不存在")
	}

	topicTopSet := dynamic.GetTopicTopSet()
	for _, v := range topicTopSet {
		if v.GetTopicId() == topicId {
			v.IsTopicTop = easygo.NewBool(false)
			v.TopicTopTime = easygo.NewInt64(0)
		}
	}
	dynamic.TopicTopSet = topicTopSet
	for_game.UpdateRedisSquareDynamic(dynamic)
	for_game.SaveSquareDynamic(dynamic.GetLogId())

	resp := &client_hall.TopicTopCancelResp{
		Result: easygo.NewInt32(1),
	}
	return resp
}

//话题排行榜规则说明
func (self *ServiceForHall) RpcTopicLeaderBoardDescription(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("====话题排行榜规则说明 RpcTopicLeaderBoardDescription=====,common=%v,reqMst=%v", common, reqMsg)
	rule := []string{
		"1、关注话题且话题贡献度（至少大于10）排名前100的用户即可入榜",
		"2、贡献度日榜每日0点清零，月榜每月1日清零，总榜累计总分不清零。",
	}
	introduction := []string{
		"1、贡献度=发动态+动态获评数+动态获赞数",
		"2、带#话题发布动态，对应话题贡献度+5，每天最多20点",
		"3、带#话题动态每获得一条评论，对应话题贡献度+1，每天最多15点",
		"4、带#话题动态每获得一个点赞，对应话题贡献度+1，每天最多15点",
	}
	resp := &client_hall.TopicLeaderBoardDescriptionResp{
		Rule:         rule,
		Introduction: introduction,
	}
	return resp
}

//退出话题主
func (self *ServiceForHall) RpcQuitTopicMaster(common *base.Common, reqMsg *client_hall.QuitTopicMasterReq) easygo.IMessage {
	logs.Info("====退出话题主 RpcQuitTopicMaster=====,common=%v,reqMst=%v", common, reqMsg)
	err := for_game.QuitTopicMaster(reqMsg.GetTopicId(), common.GetUserId())
	if err != nil {
		logs.Error("退出话题主失败", common.GetUserId(), err)
		return easygo.NewFailMsg("退出话题主失败")
	}
	resp := &client_hall.QuitTopicMasterResp{
		Result: easygo.NewInt32(1),
	}
	return resp
}

//话题主删除话题中的动态
func (self *ServiceForHall) RpcTopicMasterDelDynamic(common *base.Common, reqMsg *client_hall.TopicMasterDelDynamicReq) easygo.IMessage {
	logs.Info("====话题主删除话题中的动态 RpcTopicMasterDelDynamic=====,common=%v,reqMst=%v", common, reqMsg)
	var topicId int64
	if reqMsg.GetTopicId() != 0 {
		topicId = reqMsg.GetTopicId()
	} else {
		topicInfo := for_game.GetTopicByNameFromDB(reqMsg.GetTopicName())
		topicId = topicInfo.GetId()
	}
	IsTopicMaster := for_game.IsTopicMaster(topicId, common.GetUserId())
	if !IsTopicMaster {
		logs.Error("只有话题主可以删除动态", common.GetUserId())
		return easygo.NewFailMsg("只有话题主可以删除动态")
	}
	if reqMsg.GetDelReasonMsg() == "" {
		logs.Error("删除理由不能为空", common.GetUserId())
		return easygo.NewFailMsg("删除理由不能为空")
	}

	for_game.DelTopicDynamic(topicId, reqMsg.GetLogId())
	data := &share_message.TopicMasterDelDynamicLog{
		TopicId:      easygo.NewInt64(topicId),
		LogId:        easygo.NewInt64(reqMsg.GetLogId()),
		DelReasonId:  easygo.NewInt32(reqMsg.GetDelReasonId()),
		DelReasonMsg: easygo.NewString(reqMsg.GetDelReasonMsg()),
		TopicName:    easygo.NewString(reqMsg.GetTopicName()),
		PlayerId:     easygo.NewInt64(common.GetUserId()),
	}
	for_game.AddTopicMasterDelDynamicLog(data)
	resp := &client_hall.TopicMasterDelDynamicResp{
		Result: easygo.NewInt32(1),
	}
	dynamic := for_game.GetRedisDynamic(reqMsg.GetLogId())

	for_game.LessDevoteDynamic(dynamic.GetPlayerId(), topicId, reqMsg.GetLogId())
	//推送小助手
	content := fmt.Sprintf("%s\n经系统检测，您涉嫌“%s”，相关动态已被删除，请遵守柠檬畅聊广场规则，共同维护良好的交友氛围。", util.GetTimestamp(), reqMsg.GetDelReasonMsg())
	pushReq := &server_server.TopicDynamicTopLittleHelper{
		PlayerId: easygo.NewInt64(dynamic.GetPlayerId()),
		Title:    easygo.NewString("社交广场通知"),
		Content:  easygo.NewString(content),
	}
	pushTopicDynamicTopStatus(pushReq)
	return resp
}

// 推送话题动态置顶/取消置顶信息
func pushTopicDynamicTopStatus(pushReq *server_server.TopicDynamicTopLittleHelper) {
	logs.Info("===========推送话题助手=============pushReq: %v", pushReq)
	_, err1 := SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcTopicDynamicTopStatus", pushReq, pushReq.GetPlayerId())
	if err1 != nil {
		logs.Error("pushTopicDynamicTopStatus-err:", err1)
	}
}

//随机发送指定服务器类型
//此方法只对从im进入调用，其他不走这里pid为第三方玩家id
func SendMsgToIdelServer(t int32, methodName string, msg easygo.IMessage, pid ...int64) (easygo.IMessage, *base.Fail) {
	playerId := append(pid, 0)[0]
	var srv *share_message.ServerInfo
	if playerId != 0 {
		player := for_game.GetWishPlayerInfo(playerId)
		if player != nil {
			srv = PServerInfoMgr.GetServerInfo(player.GetHallSid())
		}
	}
	if srv == nil {
		srv = PServerInfoMgr.GetIdelServer(t)
	}
	if srv == nil {
		return nil, easygo.NewFailMsg("无法找到指定类型服务器")
	}
	return SendMsgToServerNew(srv.GetSid(), methodName, msg, pid...)
}
