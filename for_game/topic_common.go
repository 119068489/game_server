package for_game

import (
	"encoding/base64"
	"errors"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
)

const (
	OPERATE_TOPIC_ATTENTION        = 1 // 关注话题
	OPERATE_TOPIC_CANCEL_ATTENTION = 2 // 取消关注话题
)

const (
	DEFAULT_PAGE      = 1
	DEFAULT_PAGE_SIZE = 10
)

const (
	TOPIC_DEFAULT_RELATION_NUM = 20
)

const (
	REQ_TYPE_NEW = 1 // 最新
	REQ_TYPE_HOT = 2 // 热门.
)

const (
	ATTENTION_RECOMMEND_PLAYER_NUM        = 6 // 关注页推荐用户的个数.
	ATTENTION_RECOMMEND_COMMON_PLAYER_NUM = 4 // 普通用户个数，运营的话就剩下2个人
	ATTENTION_RECOMMEND_DYNAMIC_NUM       = 3 // 动态大于2条的用户
)

// 推荐页中需要的话题数量
const (
	TOPIC_HOT_NUM       = 2  // 2个热门
	TOPIC_RECOMMEND_NUM = 2  // 2个推荐
	TOPIC_COMMON_NUM    = 6  // 6个普通话题
	TOPIC_TOTAL_NUM     = 10 // 总共10个话题
)

const (
	DYNAMIC_TOPIC_HOT_NUM       = 1  // 1个热门
	DYNAMIC_TOPIC_RECOMMEND_NUM = 1  // 1个推荐
	DYNAMIC_TOPIC_TOTAL_NUM     = 10 // 总共10个话题
)

// 关注页话题条数
const (
	ATTENTION_PAGE_TOPIC_NUM = 3 // 热门类型条数为3条,总条数也是3条.
	ATTENTION_PAGE_TOP_NUM   = 1 //  话题为1条,热门,推荐,普通各一条.
	ATTENTION__RECOMMEND_NUM = 6 // 有关注人,没关注话题,随机话题数量为6
)

// 关注或取消关注话题
func OperateTopic(operate int32, playerId int64, topicIds []int64) bool {
	switch operate {
	case OPERATE_TOPIC_ATTENTION: // 关注
		for _, topicId := range topicIds {
			//判断是否有关注过
			if att := GetPlayerAttentionFromDB(topicId, playerId); att != nil {
				logs.Error("用户id: %d 已经关注了该话题了,不用重复关注,话题id为:%d", playerId, topicId)
				return false
			}
			pat := &share_message.PlayerAttentionTopic{
				Id:         easygo.NewInt64(NextId(TABLE_TOPIC_PLAYER_ATTENTION)),
				TopicId:    easygo.NewInt64(topicId),
				PlayerId:   easygo.NewInt64(playerId),
				CreateTime: easygo.NewInt64(GetMillSecond()),
			}
			// 插入用户关注的话题表
			InsertPlayerAttentionToDB(pat)
			// 更新话题粉丝数
			UpdateTopicFansToDB(topicId, 1)
		}

	case OPERATE_TOPIC_CANCEL_ATTENTION: // 取消关注
		for _, topicId := range topicIds {
			if att := GetPlayerAttentionFromDB(topicId, playerId); att == nil {
				logs.Error("用户id: %d 没有关注话题 :%d,取消操作终止", playerId, topicId)
				return false
			}
			DelPlayerAttentionFromDB(playerId, topicId)
			UpdateTopicFansToDB(topicId, -1)
		}
	default:
		return false
	}
	return true
}

// 获取关注的列表
func GetPlayerAttentionList(pid, page, pageSize int64) ([]*share_message.Topic, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}

	attentionList, count := GetTopicPlayerAttentionListByPidFromDB(pid, int(page), int(pageSize))
	topicList := make([]*share_message.Topic, 0)
	for _, v := range attentionList {
		if topic := GetTopicByIdFromDB(v.GetTopicId()); topic != nil {
			topic.IsAttention = easygo.NewBool(true)
			topicList = append(topicList, topic)
		}
	}
	return topicList, count
}

// 获取该话题中关联的话题列表
func GetRelateTopicByTypeId(typeId, topicId int64) []*share_message.Topic {
	return GetRangeTopicByTypeIdFromDB(typeId, []int64{topicId}, TOPIC_DEFAULT_RELATION_NUM)
}

//  某话题动态参与详情列表请求(推荐用户栏)
func GetPlayerListByTopicId(topicId, page, pageSize int64) ([]*share_message.TopicParticipatePlayer, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	// 根据topicId 查找动态列表,根据热门分数倒序
	dynamicList, count := GetDeviceHotDynamicByTopicIdFromDB(topicId, int(page), int(pageSize))
	playerList := make([]*share_message.TopicParticipatePlayer, 0)
	for _, d := range dynamicList {
		if base := GetRedisPlayerBase(d.GetLogId()); base != nil {
			playerList = append(playerList, &share_message.TopicParticipatePlayer{
				PlayerId:  easygo.NewInt64(base.GetPlayerId()),
				TopicId:   easygo.NewInt64(topicId),
				Sex:       easygo.NewInt32(base.GetSex()),
				HeadIcon:  easygo.NewString(base.GetHeadIcon()),
				NickName:  easygo.NewString(base.GetNickName()),
				Signature: easygo.NewString(base.GetSignature()),
				Types:     easygo.NewInt32(base.GetTypes()),
				//DynamicId: easygo.NewInt64(d.GetLogId()),
			})
		}
	}
	return playerList, count
}

func GetDynamicByTopicId(reqType, hotScore int32, pid, topicId, page, pageSize int64) ([]*share_message.DynamicData, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	var count int
	dynamicList := make([]*share_message.DynamicData, 0)
	switch reqType {
	case REQ_TYPE_NEW: // 最新
		// 根据时间排序
		dynamicList, count = GetSortTimeDynamicByTopicIdFromDB(int(page), int(pageSize), topicId)
	case REQ_TYPE_HOT: // 热门
		dynamicList, count = GetIsHotDynamicByTopicIdFromDB(hotScore, int(page), int(pageSize), topicId)
	}
	// 找到自己关注的人.
	player := GetRedisPlayerBase(pid)
	attentionList := player.GetAttention()
	dynamicList = PerfectDynamicContainHot(pid, dynamicList, attentionList, hotScore) // 遍历动态,判断是否已关注.
	return dynamicList, count
}

func CheckPlayerIsAttentionTopic(pid int64, topic *share_message.Topic) {
	if pat := GetPlayerAttentionFromDB(topic.GetId(), pid); pat != nil { // 说明已经关注了此话题
		topic.IsAttention = easygo.NewBool(true)
	}
}

// 模糊搜索话题
func SearchTopic(name string) []*share_message.Topic {
	return LikeSearchTopicByNameFromDB(name)
}

// 热门推荐,不包含自定义
//热门话题由后台话题库抽取，取十个话题，组成为：2个热门+2个推荐+6个普通话题（非自定义类），热门和推荐放在前四个的位置，重新进入页面，话题刷新
func SearchHotTopic() []*share_message.Topic {
	topicList := make([]*share_message.Topic, 0)
	// 获取热门
	hotTopicList := GetHotTopicFromDB(TOPIC_HOT_NUM)
	hotTopIds := make([]int64, 0)
	// 找出id,推荐那里排除掉热门的话题
	for _, v := range hotTopicList {
		hotTopIds = append(hotTopIds, v.GetId())
	}
	topicList = append(topicList, hotTopicList...)
	//  获取推荐
	recommendTopicList := GetRangeIsRecommendTopicFromDB(TOPIC_RECOMMEND_NUM, hotTopIds)
	for _, v := range recommendTopicList {
		hotTopIds = append(hotTopIds, v.GetId())
	}
	topicList = append(topicList, recommendTopicList...)
	// 随机抽取话题
	num := TOPIC_TOTAL_NUM - len(hotTopicList) - len(recommendTopicList) // 需要抽取的普通话题条数
	commonTopicList := GetRangeCommonTopicFromDB(num, hotTopIds)
	topicList = append(topicList, commonTopicList...)
	return topicList
}

func GetTopicListByTypeIdPage(pid int64, topicTypeId, page, pageSize int64) ([]*share_message.Topic, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	//topics, count := GetBSTopicTypeListByTypeIdPageFormDB(topicTypeId, int(page), int(pageSize))
	topics, count := GetHighHotScoreTopicByTypeIdFromDB(topicTypeId, int(page), int(pageSize))
	//封装判断是否自己已关注
	attentionTopics := GetPlayerAttentionTopicsByPidFromDB(pid)
	topIds := make([]int64, 0)
	for _, v := range attentionTopics {
		topIds = append(topIds, v.GetTopicId())
	}

	if len(topIds) == 0 {
		return topics, count
	}
	for _, t := range topics {
		if util.Int64InSlice(t.GetId(), topIds) {
			t.IsAttention = easygo.NewBool(true)
		}
	}
	return topics, count
}

func FlushTopic(page, pageSize int64) ([]*client_hall.TopicDynamicList, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	// 获取推荐的话题
	recommendTopicList, count := GetRecommendTopicByPageFromDB(int(page), int(pageSize))
	// 获取对应话题的动态.
	topicDynamicList := make([]*client_hall.TopicDynamicList, 0)

	for _, topic := range recommendTopicList {
		dynamicList, _ := GetHotDynamicByTopicIdFromDB(topic.GetId(), DEFAULT_PAGE, TOPIC_MAIN_PAGE_SIZE) // 热门分倒序排序
		//dynamicList, _ := GetHotDynamicByTopicIdFromDB(topic.GetId(), DEFAULT_PAGE, DEFAULT_PAGE_SIZE) // 热门分倒序排序
		playerList := make([]*share_message.TopicParticipatePlayer, 0)
		pids := make([]int64, 0)
		pIds1 := make([]int64, 0)
		for _, v := range dynamicList {
			pIds1 = append(pIds1, v.GetPlayerId())
		}
		playerInforList := GetAllPlayerBase(pIds1, false)
		for _, d := range dynamicList { // 人物去重
			//base := GetRedisPlayerBase(d.GetPlayerId())
			base, ok := playerInforList[d.GetPlayerId()]
			if !ok {
				continue
			}
			// 封装性别和头像,给前端调用.
			d.Sex = easygo.NewInt32(base.GetSex())
			d.HeadIcon = easygo.NewString(base.GetHeadIcon())
			d.Types = easygo.NewInt32(base.GetTypes())
			if util.Int64InSlice(d.GetPlayerId(), pids) {
				continue
			}
			pids = append(pids, d.GetPlayerId())
			playerList = append(playerList, &share_message.TopicParticipatePlayer{
				PlayerId:  easygo.NewInt64(d.GetPlayerId()),
				TopicId:   easygo.NewInt64(topic.GetId()),
				Sex:       easygo.NewInt32(base.GetSex()),
				HeadIcon:  easygo.NewString(base.GetHeadIcon()),
				NickName:  easygo.NewString(base.GetNickName()),
				Signature: easygo.NewString(base.GetSignature()),
				DynamicId: easygo.NewInt64(d.GetLogId()),
				Types:     easygo.NewInt32(base.GetTypes()),
			})
		}

		// 封装结果.
		topicDynamicList = append(topicDynamicList, &client_hall.TopicDynamicList{
			Topic:                  topic,
			DynamicList:            dynamicList,
			TopicParticipatePlayer: playerList,
		})
	}
	return topicDynamicList, count
}

/* 优化前的代码
func FlushTopic(page, pageSize int64) ([]*client_hall.TopicDynamicList, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	// 获取推荐的话题
	recommendTopicList, count := GetRecommendTopicByPageFromDB(int(page), int(pageSize))
	// 获取对应话题的动态.
	topicDynamicList := make([]*client_hall.TopicDynamicList, 0)

	for _, topic := range recommendTopicList {
		dynamicList, _ := GetHotDynamicByTopicIdFromDB(topic.GetId(), DEFAULT_PAGE, TOPIC_MAIN_PAGE_SIZE) // 热门分倒序排序
		//dynamicList, _ := GetHotDynamicByTopicIdFromDB(topic.GetId(), DEFAULT_PAGE, DEFAULT_PAGE_SIZE) // 热门分倒序排序
		playerList := make([]*share_message.TopicParticipatePlayer, 0)
		pids := make([]int64, 0)
		for _, d := range dynamicList { // 人物去重
			base := GetRedisPlayerBase(d.GetPlayerId())
			if base == nil {
				continue
			}
			// 封装性别和头像,给前端调用.
			d.Sex = easygo.NewInt32(base.GetSex())
			d.HeadIcon = easygo.NewString(base.GetHeadIcon())
			d.Types = easygo.NewInt32(base.GetTypes())
			if util.Int64InSlice(d.GetPlayerId(), pids) {
				continue
			}
			pids = append(pids, d.GetPlayerId())
			playerList = append(playerList, &share_message.TopicParticipatePlayer{
				PlayerId:  easygo.NewInt64(d.GetPlayerId()),
				TopicId:   easygo.NewInt64(topic.GetId()),
				Sex:       easygo.NewInt32(base.GetSex()),
				HeadIcon:  easygo.NewString(base.GetHeadIcon()),
				NickName:  easygo.NewString(base.GetNickName()),
				Signature: easygo.NewString(base.GetSignature()),
				DynamicId: easygo.NewInt64(d.GetLogId()),
				Types:     easygo.NewInt32(base.GetTypes()),
			})
		}

		// 封装结果.
		topicDynamicList = append(topicDynamicList, &client_hall.TopicDynamicList{
			Topic:                  topic,
			DynamicList:            dynamicList,
			TopicParticipatePlayer: playerList,
		})
	}
	return topicDynamicList, count
}
*/
func GetHotTopicList(pid, page, pageSize int64) ([]*share_message.Topic, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	topics, count := GetHotTopicByPageFromDB(int(page), int(pageSize))
	if count == 0 {
		return topics, count
	}
	//封装判断是否自己已关注
	attentionTopics := GetPlayerAttentionTopicsByPidFromDB(pid)
	topIds := make([]int64, 0)
	for _, v := range attentionTopics {
		topIds = append(topIds, v.GetTopicId())
	}

	if len(topIds) == 0 {
		return topics, count
	}
	for _, t := range topics {
		if util.Int64InSlice(t.GetId(), topIds) {
			t.IsAttention = easygo.NewBool(true)
		}
	}
	return topics, count

}

func GetAttentionRecommendPlayer() []*share_message.TopicParticipatePlayer {
	// 抽取动态大于3条的动态
	dynamicList := GetRandPlayerByDynamicCountFromDB(ATTENTION_RECOMMEND_DYNAMIC_NUM, ATTENTION_RECOMMEND_COMMON_PLAYER_NUM)

	// 抽取两个运营号
	player := GetRandPlayer(ATTENTION_RECOMMEND_PLAYER_NUM - len(dynamicList)) // 运营号的个数为 总条数-动态用户人数.

	pb := make([]*share_message.PlayerBase, 0)
	pb = append(pb, player...)
	for _, v := range dynamicList {
		p1 := GetPlayerById(v.GetLogId())
		pb = append(pb, p1)
	}
	// 封装简单的用户给前端.
	pp := make([]*share_message.TopicParticipatePlayer, 0)
	for _, p := range pb {
		pp = append(pp, &share_message.TopicParticipatePlayer{
			PlayerId:  easygo.NewInt64(p.GetPlayerId()),
			Sex:       easygo.NewInt32(p.GetSex()),
			HeadIcon:  easygo.NewString(p.GetHeadIcon()),
			NickName:  easygo.NewString(p.GetNickName()),
			Signature: easygo.NewString(p.GetSignature()),
			FansNum:   easygo.NewInt64(len(p.GetFansList())),
			Types:     easygo.NewInt32(p.GetTypes()),
		})
	}
	return pp
}

// 浏览数
func OperateViewNum(dnList []*share_message.DynamicData) {
	for _, dn := range dnList {
		content := dn.GetContent()
		if content == "" {
			continue
		}
		// 解码
		contentBytes, _ := base64.StdEncoding.DecodeString(content)
		OperateTopicViewNum(string(contentBytes), 1)
	}
}

// 动态里面的话题. 1个热门+1个推荐+8个普通话题（非自定义类）
func GetDynamicTopicList() []*share_message.Topic {
	topicList := make([]*share_message.Topic, 0)
	// 获取热门
	hotTopicList := GetHotTopicFromDB(DYNAMIC_TOPIC_HOT_NUM)
	hotTopIds := make([]int64, 0)
	// 找出id,推荐那里排除掉热门的话题
	for _, v := range hotTopicList {
		hotTopIds = append(hotTopIds, v.GetId())
	}
	topicList = append(topicList, hotTopicList...)
	//  获取推荐
	recommendTopicList := GetRangeIsRecommendTopicFromDB(DYNAMIC_TOPIC_RECOMMEND_NUM, hotTopIds)
	for _, v := range recommendTopicList {
		hotTopIds = append(hotTopIds, v.GetId())
	}
	topicList = append(topicList, recommendTopicList...)
	// 随机抽取话题
	num := DYNAMIC_TOPIC_TOTAL_NUM - len(hotTopicList) - len(recommendTopicList) // 需要抽取的普通话题条数
	commonTopicList := GetRangeCommonTopicFromDB(num, hotTopIds)
	topicList = append(topicList, commonTopicList...)
	return topicList
}

/**
    带话题的社交广场关注页的请求
	1.第一页
	校验当前用户是否有关注话题,或者是关注过人.
	(1)用户有关注人,有关注话题
	(2)用户无关注人,无关注话题
	(3)用户无关注人,有关注话题
	(4)用户有关注人,无关注话题
	2.不是第一页,判断是哪种情况.
`	对照着步骤1请求数据
*/
//func GetSquareAttentionData(pid int64, hotScore int32, req *client_hall.SquareAttentionReq) ([]*share_message.DynamicData, []*share_message.Topic,
//	[]*client_hall.OneTopicType, []*share_message.TopicParticipatePlayer, int) {
func GetSquareAttentionData(pid int64, hotScore int32, req *client_hall.SquareAttentionReq) *client_hall.SquareAttentionResp {
	page, pageSize := req.GetPage(), req.GetPageSize()
	if req.GetPage() == 0 {
		page = DEFAULT_PAGE
	}
	if req.GetPageSize() == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}

	hasAttentionTopic, hasAttentionPlayer := req.GetHasAttentionTopic(), req.GetHasAttentionPlayer() // 是否有关注用户
	if page == 1 {
		hasAttentionTopic, hasAttentionPlayer = false, false // 重置条件.
		player := GetRedisPlayerBase(pid)
		if player == nil {
			logs.Error("带话题的社交广场关注页的请求,用户不存在,pid:", pid)
			return nil
		}
		// 是否有关注话题
		if attentionTopicList := GetPlayerAttentionTopicsByPidFromDB(pid); len(attentionTopicList) > 0 {
			hasAttentionTopic = true
		}
		// 是否有关注过人.
		if attentionList := player.GetAttention(); len(attentionList) > 0 {
			hasAttentionPlayer = true
		}
	}
	data := make([]*share_message.DynamicData, 0)
	topics := make([]*share_message.Topic, 0)
	types := make([]*client_hall.OneTopicType, 0)
	players := make([]*share_message.TopicParticipatePlayer, 0)
	var count int
	if hasAttentionPlayer && hasAttentionTopic { // 用户有关注人,有关注话题
		data, topics, count = getHasPlayerHasTopic(pid, page, pageSize, hotScore)
	} else if !hasAttentionPlayer && !hasAttentionTopic { // 用户无关注人,无关注话题
		types, players = getNoPlayerNoTopic()
	} else if !hasAttentionPlayer && hasAttentionTopic { // 用户无关注人,有关注话题
		data, topics, count = getNoPlayerHasTopic(pid, page, pageSize)
	} else if hasAttentionPlayer && !hasAttentionTopic { // 用户有关注人,无关注话题
		data, topics, count = getHasPlayerNoTopic(pid, page, pageSize, hotScore)
	}
	//判断是否已关注
	if len(topics) > 0 {
		//封装判断是否自己已关注
		attentionTopics := GetPlayerAttentionTopicsByPidFromDB(pid)
		topIds := make([]int64, 0)
		for _, v := range attentionTopics {
			topIds = append(topIds, v.GetTopicId())
		}
		for _, t := range topics {
			if util.Int64InSlice(t.GetId(), topIds) {
				t.IsAttention = easygo.NewBool(true)
			}
		}
	}
	return &client_hall.SquareAttentionResp{
		TopicList:          topics,
		DynamicList:        data,
		TopicTypeList:      types,
		PlayerList:         players,
		HasAttentionTopic:  easygo.NewBool(hasAttentionTopic),
		HasAttentionPlayer: easygo.NewBool(hasAttentionPlayer),
		Count:              easygo.NewInt64(count),
	}
}

// 用户有关注人,有关注话题
func getHasPlayerHasTopic(pid, page, pageSize int64, hotScore int32) ([]*share_message.DynamicData, []*share_message.Topic, int) {
	// 关注话题
	attentionTopics, _ := GetTopicPlayerAttentionListByPidFromDB(pid, int(page), int(pageSize))
	topics := make([]*share_message.Topic, 0)
	for _, v := range attentionTopics {
		if topic := GetTopicByIdFromDB(v.GetTopicId()); topic != nil {
			topics = append(topics, topic)
		}
	}
	dsList, count := playerAttention(pid, page, pageSize, hotScore)
	return dsList, topics, count
}

// 玩家关注的动态
func playerAttention(pid int64, page int64, pageSize int64, hotScore int32) ([]*share_message.DynamicData, int) {
	dsList := make([]*share_message.DynamicData, 0)
	player := GetRedisPlayerBase(pid)
	attentionList := player.GetAttention()
	// 第一页的话,需要置顶消息
	if page == 1 {
		//获取后台置顶的动态列表(in)
		bsTopDynamicList := GetBSTopDynamicListByIDsFromDB(pid, attentionList)
		if len(bsTopDynamicList) > 0 {
			slice := GetDynamicSliceByRandFromSlice(bsTopDynamicList, BS_TOP_NUM)
			//  时间最新的最靠前
			sortSlice := SortDynamicSliceByTime1(slice)
			dsList = append(dsList, sortSlice...)
		}

		//获取app置顶的动态列表(in)
		appTopDynamicList := GetAppTopDynamicListByIDsFromDB(pid, attentionList)
		if len(appTopDynamicList) > 0 {
			slice := GetDynamicSliceByRandFromSlice(appTopDynamicList, APP_TOP_NUM)
			//// 时间最新的最靠前
			sortSlice := SortDynamicSliceByTime1(slice)
			dsList = append(dsList, sortSlice...)
		}
	}
	maxLogIdKey := MakeNewString(pid, "attention")
	// 关注.
	ds, count := GetNoTopDynamicByPIDs(pid, int(page), int(pageSize), attentionList, maxLogIdKey)

	dsList = append(dsList, ds...)
	dsList = ParseHotDynamic(dsList, hotScore)
	dsList = GetRedisSomeDynamic1(pid, dsList, attentionList) // 遍历动态,判断是否已关注.
	return dsList, count
}

// 用户无关注人,无关注话题
func getNoPlayerNoTopic() ([]*client_hall.OneTopicType, []*share_message.TopicParticipatePlayer) {
	// 推荐用户  从数据库抽取6个用户，组成为：2个运营号+4个普通用户（带有3条动态以上）
	participatePlayer := GetAttentionRecommendPlayer()
	// 推荐话题 每个话题类别展示3个话题，组成为：热门（热度第一）+推荐（热度第一）+普通话题（热度第一) 原文档的需求,已废弃2020.11.16(谷尼)

	// 推荐话题 每个话题类别展示3个话题，组成为：热门+推荐+普通话题,没有就用普通的顶上.
	hotTopicList := GetHotTopicFromDB(ATTENTION_PAGE_TOPIC_NUM)
	if len(hotTopicList) < ATTENTION_PAGE_TOPIC_NUM { // 普通话题顶上
		topicIds := make([]int64, 0)
		for _, v := range hotTopicList {
			topicIds = append(topicIds, v.GetId())
		}
		tps := GetRangeCommonTopicFromDB(ATTENTION_PAGE_TOPIC_NUM-len(hotTopicList), topicIds)
		hotTopicList = append(hotTopicList, tps...)
	}

	// 获取所有类别(官方)
	topicTypeList := GetBSTopicTypeListByClassFormDB(TOPIC_CLASS_BS)
	ttList := make([]*client_hall.OneTopicType, 0)
	ttList = append(ttList, &client_hall.OneTopicType{TopicList: hotTopicList})
	for _, tt := range topicTypeList {
		tl := make([]*share_message.Topic, 0)
		// 根据类别找热门话题,没有就普通顶上
		hotTopicList := GetRangeHotTopicByTypeIdFromDB(tt.GetId(), []int64{}, ATTENTION_PAGE_TOP_NUM)
		// 排除的id
		ids := make([]int64, 0)
		for _, v := range hotTopicList {
			ids = append(ids, v.GetId())
		}
		tl = append(tl, hotTopicList...)
		// 根据类别找推荐的话题,没有就普通顶上
		recommendTopicList := GetRangeRecommendTopicByTypeIdFromDB(tt.GetId(), ids, ATTENTION_PAGE_TOP_NUM)
		for _, v := range recommendTopicList {
			ids = append(ids, v.GetId())
		}
		tl = append(tl, recommendTopicList...)
		// 根据类别找普通的话题,没有就普通顶上
		topicList := GetRangeCommentTopicByTypeIdFromDB(tt.GetId(), ids, ATTENTION_PAGE_TOPIC_NUM-len(hotTopicList)-len(recommendTopicList))
		tl = append(tl, topicList...)
		ttList = append(ttList, &client_hall.OneTopicType{
			TopicType: tt,
			TopicList: tl,
		})
	}
	return ttList, participatePlayer
}

// 用户无关注人,有关注话题
func getNoPlayerHasTopic(pid, page, pageSize int64) ([]*share_message.DynamicData, []*share_message.Topic, int) {
	// 关注话题
	attentionTopics, _ := GetTopicPlayerAttentionListByPidFromDB(pid, int(page), int(pageSize))
	if len(attentionTopics) == 0 {
		logs.Error("用户无关注人,有关注话题,话题页没有,pid:", pid)
		return nil, nil, 0
	}
	topics := make([]*share_message.Topic, 0)
	topicIds := make([]int64, 0)
	for _, v := range attentionTopics {
		if topic := GetTopicByIdFromDB(v.GetTopicId()); topic != nil {
			topics = append(topics, topic)
			topicIds = append(topicIds, v.GetTopicId())
		}
	}
	// 查找包含话题id的动态.
	dynamicData, count := GetSortTimeDynamicByTopicIdListFromDB(int(page), int(pageSize), topicIds)
	player := GetRedisPlayerBase(pid)
	attentionList := player.GetAttention()
	dynamicData = GetRedisSomeDynamic1(pid, dynamicData, attentionList) // 遍历动态,判断是否已关注.
	return dynamicData, topics, count
}

// 用户有关注人,无关注话题
func getHasPlayerNoTopic(pid, page, pageSize int64, hotScore int32) ([]*share_message.DynamicData, []*share_message.Topic, int) {
	// 随机从话题库抽6个推荐话题
	topics := GetRangeIsRecommendTopicFromDB(ATTENTION__RECOMMEND_NUM, []int64{})
	if len(topics) < 6 { // 普通话题不上
		tid := make([]int64, 0)
		for _, t := range topics {
			tid = append(tid, t.GetId())
		}
		commentTopic := GetRangeCommonTopicFromDB(ATTENTION__RECOMMEND_NUM-len(topics), tid)
		topics = append(topics, commentTopic...)
	}
	dsList, count := playerAttention(pid, page, pageSize, hotScore)
	return dsList, topics, count
}

//  某话题动态参与详情列表请求(推荐用户栏)
func GetDevicePlayerHotDynamicByTopicId(topicId, page, pageSize int64, hotScore int32) ([]*share_message.TopicParticipatePlayer, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	// 根据topicId 查找动态列表,根据热门分数倒序
	dynamicList, count := GetDevicePlayerHotDynamicByTopicIdFromDB(topicId, hotScore, int(page), int(pageSize))
	playerList := make([]*share_message.TopicParticipatePlayer, 0)
	for _, d := range dynamicList {
		if base := GetRedisPlayerBase(d.GetLogId()); base != nil {
			playerList = append(playerList, &share_message.TopicParticipatePlayer{
				PlayerId:  easygo.NewInt64(base.GetPlayerId()),
				TopicId:   easygo.NewInt64(topicId),
				Sex:       easygo.NewInt32(base.GetSex()),
				HeadIcon:  easygo.NewString(base.GetHeadIcon()),
				NickName:  easygo.NewString(base.GetNickName()),
				Signature: easygo.NewString(base.GetSignature()),
				Types:     easygo.NewInt32(base.GetTypes()),
			})
		}
	}
	return playerList, count
}

// 广场按钮旁边的话题头部话题列表,
func GetTopicHeadTopic() []*share_message.Topic {
	topics := make([]*share_message.Topic, 0)
	// 去重id列表
	tids := make([]int64, 0)
	// 1个推荐话题
	recommendTopics := GetRangeIsRecommendTopicFromDB(TOPIC_HEAD_TOPIC_RECOMD_NUM, []int64{})
	for _, v := range recommendTopics {
		tids = append(tids, v.GetId())
	}
	topics = append(topics, recommendTopics...)
	// 1 个热门话题
	hotTopics := GetRangeHotTopicFromDB(TOPIC_HEAD_TOPIC_HOT_NUM, tids)
	for _, v := range hotTopics {
		tids = append(tids, v.GetId())
	}
	topics = append(topics, hotTopics...)
	// 3个普通话题
	commonTopics := GetRangeCommonTopicFromDB(TOPIC_HEAD_TOPIC_ALL_NUM-len(recommendTopics)-len(hotTopics), tids)
	topics = append(topics, commonTopics...)
	return topics
}

//获取用户详情
func GetPlayerInfo(playerId int64) (*share_message.PlayerBase, error) {
	playerInfoOJB := GetRedisPlayerBase(playerId)
	if playerInfoOJB == nil {
		return nil, errors.New("用户不存在")
	}
	playerInfo := playerInfoOJB.QueryPlayerBase(playerId)
	if playerInfo == nil {
		return nil, errors.New("用户不存在")
	}
	return playerInfo, nil
}

//获取用户注册时间
func GetPlayerRegisteredCondition(playerId int64) bool {
	playerInfoOJB := GetRedisPlayerBase(playerId)
	if playerInfoOJB == nil {
		return false
	}
	playerInfo, err := GetPlayerInfo(playerId)
	if err != nil {
		logs.Error(err.Error(), playerId)
		return false
	}
	milliTime := util.GetMilliTime()
	milliTime30Day := milliTime - 30*24*3600*1000 //距离现在30天前的时间戳，毫秒
	//判断注册时间是否满足30天
	if playerInfo.GetCreateTime() < milliTime30Day {
		return true
	}
	return false
}

//判断用户是否为话题主
func IsTopicMaster(topicId, playerId int64) bool {
	topicInfo := GetTopicByIdFromDB(topicId)
	if topicInfo.GetTopicMaster() == playerId {
		return true
	}
	return false
}

//删除话题动态
func DelTopicDynamic(topicId, logId int64) bool {
	dynamicInfo, err := GetTopicDynamicByLogId(topicId, logId)
	if err != nil {
		return false
	}
	topicArr := make([]int64, 0)
	for _, v := range dynamicInfo.GetTopicId() {
		if v != topicId {
			topicArr = append(topicArr, v)
		}
	}
	err = UpTopicDynamicByLogId(topicId, logId, topicArr)
	if err != nil {
		return false
	}
	DelSquareDynamicById(logId)
	return true
}

func GetNewDynamicByTopicId(topicId, page, pageSize int64) ([]*share_message.DynamicData, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	var count int
	dynamicList := make([]*share_message.DynamicData, 0)
	dynamicList, count = GetSortTimeDynamicByTopicIdFromDB(int(page), int(pageSize), topicId)
	return dynamicList, count
}
