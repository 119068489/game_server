package backstage

import (
	"encoding/base64"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/client_server"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"strconv"
	"strings"
	"time"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//	RpcQueryDynamic	查询社交动态
func (self *cls4) RpcQueryDynamic(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.DynamicListRequest) easygo.IMessage {
	// logs.Info("请求社交广场RpcQueryDynamic:", reqMsg)
	msg := &brower_backstage.DynamicListResponse{
		List:      nil,
		PageCount: easygo.NewInt32(0),
	}

	permission := true
	//先注释，后面再开放
	//permission := false
	//if user.GetRole() != 0 {
	//	permission = QueryPermissionById(user.GetSite(), user.GetRoleType(), "groupManage-oneself")
	//}

	if permission { //有权限
		list, count := QueryDynamic(reqMsg)
		msg.List = list
		msg.PageCount = easygo.NewInt32(count)
	}

	return msg
}

//	RpcQueryDynamicDetails	查询社交动态详情
func (self *cls4) RpcQueryDynamicDetails(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logId := reqMsg.GetId64()
	dbdy := QueryDynamicById(logId)
	config := for_game.QuerySysParameterById(for_game.SQUAREHOT_PARAMETER)
	if dbdy.GetStatue() == for_game.DYNAMIC_COMMENT_STATUE_COMMON {
		who := for_game.GetRedisPlayerBase(dbdy.GetPlayerId())
		dynamicInfo := for_game.GetRedisDynamicAllInfo(0, logId, who.GetPlayerId(), who.GetAttention())

		dynamicInfo.Account = easygo.NewString(who.GetAccount())
		dynamicInfo.NickName = easygo.NewString(who.GetNickName())

		if config != nil {
			if dbdy.GetHostScore() >= config.GetHotScore() {
				dynamicInfo.IsHot = easygo.NewBool(true)
			}
		}

		return dynamicInfo
	}

	if config != nil {
		if dbdy.GetHostScore() >= config.GetHotScore() {
			dbdy.IsHot = easygo.NewBool(true)
		}
	}

	return dbdy
}

//	RpcQueryCommentDetails	查询评论详情
func (self *cls4) RpcQueryCommentDetails(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.DynamicListRequest) easygo.IMessage {
	list, count := QueryDynamicComment(reqMsg)
	msg := &brower_backstage.CommentList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//批量删除评论
func (self *cls4) RpcDeleteCommentDatas(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	if len(reqMsg.GetIds64()) == 0 {
		return easygo.NewFailMsg("请先选择要删除的评论id64")
	}
	if len(reqMsg.GetIdsStr()) == 0 {
		return easygo.NewFailMsg("动态ID不能为空IdsStr")
	}
	i, _ := strconv.ParseInt(reqMsg.GetIdsStr()[0], 10, 64)
	if i < 0 {
		return easygo.NewFailMsg("动态ID错误")
	}

	for _, value := range reqMsg.GetIds64() {
		for_game.DelRedisDynamicComment(i, value, for_game.DYNAMIC_COMMENT_STATUE_DELETE, reqMsg.GetNote())
	}

	msg := "删除动态评论：" + easygo.Int64ArrayToString(reqMsg.GetIds64())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SQUARE_MANAGE, msg)
	return easygo.EmptyMsg
}

//	RpcUpdateDynamic	发布或修改社区动态
func (self *cls4) RpcUpdateDynamic(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.DynamicData) easygo.IMessage {
	// logs.Info("RpcUpdateDynamic,msg=%v", reqMsg)
	if reqMsg.GetIsBsTop() && reqMsg.GetTopOverTime() <= easygo.NowTimestamp() {
		return easygo.NewFailMsg("取消置顶时间不能小于当前时间")
	}
	if reqMsg.Statue == nil {
		return easygo.NewFailMsg("发布状态不能为空")
	}
	if reqMsg.GetContent() == "" && len(reqMsg.GetPhoto()) == 0 && reqMsg.GetVoice() == "" && reqMsg.GetVideo() == "" {
		return easygo.NewFailMsg("无效的动态，什么内容都没有！！！！")
	}
	if reqMsg.IsHot != nil && !reqMsg.GetIsHot() {
		return easygo.NewFailMsg("热门不能手动取消")
	}

	if reqMsg.GetIsBsTop() && reqMsg.GetIsHot() {
		return easygo.NewFailMsg("置顶状态的动态不能设置热门")
	}

	photolst := reqMsg.GetPhoto()
	if len(photolst) > 0 { //违规图检验
		for _, photo := range photolst {
			num := ImageModeration(photo)
			if num != 100 {
				return easygo.NewFailMsg("违规图片，不允许发送！")
			}
		}
	}
	content := reqMsg.GetContent()
	topics := for_game.CheckTopic(content)
	topicIds := make([]int64, 0)
	//取前5个话题
	if len(topics) > for_game.DYNAMIC_TOPIC_NUM {
		topics = topics[:for_game.DYNAMIC_TOPIC_NUM]
	}
	for _, topicName := range topics {
		topic := for_game.GetTopicByNameFromDB(topicName)
		if topic == nil {
			continue
		}
		topicIds = append(topicIds, topic.GetId())
		// 给该话题加参与数
		for_game.IncTopicParticipationNumToDB(topic.GetId(), 1)
	}
	if len(topicIds) > 0 {
		reqMsg.TopicId = topicIds
	}
	text := base64.StdEncoding.EncodeToString([]byte(content))
	reqMsg.Content = easygo.NewString(text)

	player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	msg := "发布社交动态"
	if reqMsg.LogId == nil && reqMsg.GetLogId() == 0 { //如果新增
		reqMsg.PlayerId = easygo.NewInt64(reqMsg.GetPlayerId())
		reqMsg.Sex = easygo.NewInt32(reqMsg.GetSex())
		reqMsg.CreateTime = easygo.NewInt64(time.Now().Unix())
		reqMsg.Zan = easygo.NewInt32(0)
		reqMsg.TrueZan = easygo.NewInt32(0)
		reqMsg.Check = easygo.NewInt32(for_game.DYNAMIC_CHECK_AUTOOK)
		if reqMsg.SendTime == nil || reqMsg.GetSendTime() == 0 {
			reqMsg.SendTime = reqMsg.CreateTime
		}
		eq := for_game.GetRedisPlayerEquipmentObj(reqMsg.GetPlayerId())
		equipment := eq.GetEquipmentForClient()
		qp := equipment.GetQP()
		if qp != nil {
			reqMsg.PropsId = easygo.NewInt64(qp.GetPropsId())
		}
		if reqMsg.GetIsBsTop() {
			switch reqMsg.GetStatue() {
			case for_game.DYNAMIC_STATUE_UNPUBLISHED:
				if (reqMsg.GetTopOverTime() - reqMsg.GetSendTime()) < 60*5 {
					return easygo.NewFailMsg("置顶时间必须大于发布时间至少5分钟！")
				}
			default:
				if (reqMsg.GetTopOverTime() - easygo.NowTimestamp()) < 60*5 {
					return easygo.NewFailMsg("置顶时间必须大于当前时间至少5分钟！")
				}
			}
		}
		if reqMsg.GetStatue() == for_game.DYNAMIC_STATUE_EXPIRED {
			reqMsg.Statue = easygo.NewInt32(for_game.DYNAMIC_STATUE_COMMON)
		}
		if reqMsg.GetStatue() != for_game.DYNAMIC_STATUE_UNPUBLISHED { // 不是延时的,就直接入库
			reqMsg.LogId = easygo.NewInt64(for_game.NextId(for_game.TABLE_SQUARE_DYNAMIC))
			AddSquareDynamic(player, reqMsg)
		}

		if reqMsg.GetStatue() == for_game.DYNAMIC_STATUE_UNPUBLISHED {
			reqMsg.LogId = easygo.NewInt64(for_game.NextId(for_game.TABLE_SQUARE_DYNAMIC_SNAP))
			AddSendDynamic(reqMsg) //定时发送
		}
	} else {
		msg = "修改社交动态"
		if SendDynamicTimeMgr.GetTimerById(reqMsg.GetLogId()) != nil {
			SendDynamicTimeMgr.DelTimerList(reqMsg.GetLogId())
		}

		topics := reqMsg.GetTopicTopSet()
		for _, ts := range topics {
			ts.TopicTopTime = easygo.NewInt64(easygo.NowTimestamp())
		}

		// 判断修改的动态是延时动态还是正常的动态
		if reqMsg.GetOldStatue() == for_game.DYNAMIC_STATUE_UNPUBLISHED {
			if reqMsg.GetStatue() == for_game.DYNAMIC_STATUE_EXPIRED || reqMsg.GetStatue() == for_game.DYNAMIC_STATUE_COMMON { // 从临时表中删除,从正式表中插入
				if reqMsg.GetStatue() == for_game.DYNAMIC_STATUE_COMMON {
					reqMsg.SendTime = easygo.NewInt64(easygo.NowTimestamp())
					SendDynamicTimeMgr.DelTimerList(reqMsg.GetLogId()) // 删除定时任务,临时表的id
				}
				reqMsg.Statue = easygo.NewInt32(for_game.DYNAMIC_STATUE_COMMON)
				//AddSquareDynamic(player, reqMsg)
				DelDynamicFromSnapDB(reqMsg.GetLogId())
				reqMsg.LogId = easygo.NewInt64(for_game.NextId(for_game.TABLE_SQUARE_DYNAMIC))
				// 把动态存在正式动态表中
				insertErr := for_game.InsertDynamicToDB(reqMsg)
				easygo.PanicError(insertErr)
			}
			if reqMsg.GetStatue() == for_game.DYNAMIC_STATUE_UNPUBLISHED {
				AddSendDynamic(reqMsg) //定时发送
			}

		} else {
			// 判断是否是从未来动态改成现有的动态.
			dynamic := for_game.GetRedisDynamic(reqMsg.GetLogId())
			if dynamic == nil {
				logs.Error("后台修改正常的动态,该动态不存在,正常动态的id为: %d", reqMsg.GetLogId())
				return easygo.NewFailMsg("该动态不存在")
			}
			dynamic.TrueZan = easygo.NewInt32(reqMsg.GetZan())
			if reqMsg.IsBsTop != nil {
				dynamic.IsBsTop = easygo.NewBool(reqMsg.GetIsBsTop())
				if reqMsg.TopOverTime != nil {
					dynamic.TopOverTime = easygo.NewInt64(reqMsg.GetTopOverTime())
					if reqMsg.GetTopOverTime() > 0 && (reqMsg.GetTopOverTime()-easygo.NowTimestamp()) < 60*5 {
						return easygo.NewFailMsg("置顶时间必须大于当前时间至少5分钟！")
					}
				}
			}
			dynamic.TopicTopSet = reqMsg.TopicTopSet
			for_game.UpdateRedisSquareDynamic(dynamic)
			for_game.SaveSquareDynamic(dynamic.GetLogId())
			easygo.Spawn(SetHotDynamic, reqMsg.GetLogId(), reqMsg.GetIsHot())
		}
	}
	// 不管是新增还是修改,未来的动态,都不在这里处理通知置顶,在定时任务那边处理
	if reqMsg.GetIsBsTop() && reqMsg.GetStatue() == for_game.DYNAMIC_STATUE_COMMON {
		msg := &server_server.TopRequest{
			LogId: reqMsg.LogId,
		}
		// BroadCastToAllHall("RpcBackstageTop", msg)
		ChooseOneHall(0, "RpcBackstageTop", msg) //通知大厅置顶社交动态
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SQUARE_MANAGE, msg)
	return easygo.EmptyMsg
}

//	RpcDeleteDynamic	删除社区动态
func (self *cls4) RpcDeleteDynamic(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.DelDynamicRequest) easygo.IMessage {
	var ids []int64
	for _, value := range reqMsg.GetList() {
		player := QueryPlayerbyAccount(value.GetIdStr())
		objLog := for_game.GetRedisDynamic(value.GetId64())
		b := for_game.DelRedisSquareDynamic(player.GetPlayerId(), value.GetId64(), reqMsg.GetNote(), for_game.DYNAMIC_STATUE_DELETE)
		if b {
			base := for_game.GetRedisPlayerBase(player.GetPlayerId())
			base.DelRedisPlayerDynamicList(value.GetId64())

			SendSystemNotice(player.GetPlayerId(), "社交广场通知", "您在"+easygo.Stamp2Str(objLog.GetSendTime())+"发布的动态涉嫌违规，已被删除，如有疑问请联系客服")
		}
		ids = append(ids, value.GetId64())
	}
	msg := "删除动态：" + easygo.Int64ArrayToString(ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SQUARE_MANAGE, msg)
	return easygo.EmptyMsg
}

//删除未发布的社区动态
func (self *cls4) RpcDeleteUnDynamic(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.DelDynamicRequest) easygo.IMessage {
	var ids []int64
	for _, value := range reqMsg.GetList() {
		if SendDynamicTimeMgr.GetTimerById(value.GetId64()) != nil {
			SendDynamicTimeMgr.DelTimerList(value.GetId64())
		}
		ids = append(ids, value.GetId64())
	}

	findBson := bson.M{"_id": bson.M{"$in": ids}}
	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC_SNAP, findBson)

	msg := "删除未发布的动态：" + easygo.Int64ArrayToString(ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SQUARE_MANAGE, msg)
	return easygo.EmptyMsg
}

//屏蔽社交动态
func (self *cls4) RpcShieldDynamic(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	if len(reqMsg.GetIds64()) == 0 {
		return easygo.NewFailMsg("请先选择要屏蔽的项id64")
	}
	if len(reqMsg.GetIds32()) == 0 {
		return easygo.NewFailMsg("屏蔽操作不能为空id32")
	}
	for _, value := range reqMsg.GetIds64() {
		b := for_game.GetRedisDynamic(value)
		if b != nil {
			switch reqMsg.GetIds32()[0] {
			case 1:
				b.IsShield = easygo.NewBool(true)
			case 2:
				b.IsShield = easygo.NewBool(false)
			}
			b.Note = easygo.NewString(reqMsg.GetNote())
			for_game.UpdateRedisSquareDynamic(b)
			for_game.SaveSquareDynamic(b.GetLogId())
		}
	}

	msg := "屏蔽动态：" + easygo.Int64ArrayToString(reqMsg.GetIds64())
	if reqMsg.GetIds32()[0] == 2 {
		msg = "取消屏蔽动态：" + easygo.Int64ArrayToString(reqMsg.GetIds64())
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SQUARE_MANAGE, msg)
	return easygo.EmptyMsg
}

func AddSquareDynamic(who *for_game.RedisPlayerBaseObj, reqMsg *share_message.DynamicData) easygo.IMessage {
	logId := reqMsg.GetLogId()
	if logId == 0 {
		logId = for_game.NextId(for_game.TABLE_SQUARE_DYNAMIC)
	}

	reqMsg.LogId = easygo.NewInt64(logId)
	reqMsg.PlayerId = easygo.NewInt64(who.GetPlayerId())
	if reqMsg.GetStatue() == for_game.DYNAMIC_STATUE_COMMON {
		if reqMsg.SendTime == nil || reqMsg.GetSendTime() == 0 {
			reqMsg.SendTime = easygo.NewInt64(easygo.NowTimestamp())
		}
		for_game.AddRedisSquareDynamic(reqMsg)
		for_game.SaveSquareDynamic(reqMsg.GetLogId())
		who.AddRedisPlayerDynamicList(logId)
		easygo.Spawn(SetHotDynamic, logId, reqMsg.GetIsHot())
	} else {
		SaveDynamic(reqMsg)
	}

	msg := &client_server.RequestInfo{
		Id: easygo.NewInt64(logId),
	}
	return msg
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

//审核社交动态
func (self *cls4) RpcReviewDynamic(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logId := reqMsg.GetId64()
	if logId == 0 {
		return easygo.NewFailMsg("id错误")
	}
	if reqMsg.GetId32() == 0 {
		return easygo.NewFailMsg("审核状态错误")
	}

	b := QueryDynamicById(reqMsg.GetId64())
	if b == nil {
		return easygo.NewFailMsg("动态不存在")
	}

	UpdateDynamicCheck(logId, reqMsg.GetId32()) //修改审核状态
	playerMgr := for_game.GetRedisPlayerBase(b.GetPlayerId())
	if playerMgr == nil {
		return easygo.NewFailMsg("动态数据错误，玩家不存在")
	}
	msg := "审核拒绝动态：" + easygo.AnytoA(logId)
	switch reqMsg.GetId32() {
	case 1:
		playerMgr.IncrCheckNum(1)
		msg = "审核通过动态" + easygo.AnytoA(logId)
	case 2:
		if reqMsg.GetNote() == "" || reqMsg.Note == nil {
			return easygo.NewFailMsg("拒绝原因不能为空")
		}
		val := playerMgr.GetCheckNum()
		playerMgr.IncrCheckNum(-val)
		log := for_game.GetRedisDynamic(logId) //
		if log != nil {
			b := for_game.DelRedisSquareDynamic(playerMgr.GetPlayerId(), logId, reqMsg.GetNote(), for_game.DYNAMIC_STATUE_DELETE)
			if b {
				playerMgr.DelRedisPlayerDynamicList(logId)
				s := fmt.Sprintf("经系统检测,您涉嫌“%s”,相关动态已被删除,请遵守柠檬畅聊广场规则,共同维护良好的交友氛围", reqMsg.GetNote())
				SendSystemNotice(playerMgr.GetPlayerId(), "社交广场通知", s)
			}
		}
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SQUARE_MANAGE, msg)
	return easygo.EmptyMsg
}

//==话题相关===========================================================================>
//话题类型列表
func (self *cls4) RpcQueryTopicTypeList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	// logs.Info("RpcQueryTopicTypeList:", reqMsg)
	list, count := QueryTopicTypeList(reqMsg)
	msg := &brower_backstage.TopicTypeResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//更新话题类型
func (self *cls4) RpcUpdateTopicType(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.TopicType) easygo.IMessage {
	if reqMsg.Name == nil {
		return easygo.NewFailMsg("类型名不能为空")
	}
	if reqMsg.GetHotCount() < 0 {
		return easygo.NewFailMsg("热门数量不能小于0")
	}

	if reqMsg.GetStatus() <= 0 {
		return easygo.NewFailMsg("状态设置错误")
	}

	nameOk := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_TYPE, bson.M{"Name": reqMsg.Name})
	if nameOk != nil && nameOk.(bson.M)["_id"].(int64) != reqMsg.GetId() {
		return easygo.NewFailMsg("添加失败，已存在该话题类别")
	}

	msg := fmt.Sprintf("修改话题类型:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_TOPIC_TYPE))
		reqMsg.CreateTime = easygo.NewInt64(util.GetMilliTime())
		reqMsg.TopicClass = easygo.NewInt32(1)
		msg = fmt.Sprintf("添加话题类型:%d", reqMsg.GetId())
	}

	reqMsg.UpdateTime = easygo.NewInt64(util.GetMilliTime())

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_TYPE, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.TOPIC_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询话题列表
func (self *cls4) RpcQueryTopicList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	// logs.Info("RpcQueryTopicList:", reqMsg)
	list, count := QueryTopicList(reqMsg)
	msg := &brower_backstage.TopicResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return msg
}

//更新话题
func (self *cls4) RpcUpdateTopic(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.Topic) easygo.IMessage {
	if reqMsg.Name == nil {
		return easygo.NewFailMsg("名称不能为空")
	}

	if reqMsg.GetStatus() == 0 {
		return easygo.NewFailMsg("状态设置错误")
	}

	msg := fmt.Sprintf("修改话题:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		// 判断前端是否有加#
		name := reqMsg.GetName()
		if !strings.HasPrefix(name, "#") {
			name = fmt.Sprintf("%s%s", "#", name)
		}
		if !strings.HasSuffix(name, "#") {
			name = fmt.Sprintf("%s%s", name, "#")
		}

		nameOk := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, bson.M{"Name": name})
		if nameOk != nil {
			return easygo.NewFailMsg("话题已存在")
		}

		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_TOPIC))
		reqMsg.CreateTime = easygo.NewInt64(util.GetMilliTime())
		reqMsg.Admin = easygo.NewString(user.GetAccount())
		reqMsg.FansNum = easygo.NewInt64(reqMsg.GetFansNum())
		reqMsg.ParticipationNum = easygo.NewInt64(reqMsg.GetParticipationNum())
		reqMsg.ViewingNum = easygo.NewInt64(reqMsg.GetViewingNum())
		reqMsg.AddViewingNum = easygo.NewInt64(reqMsg.GetAddViewingNum())
		reqMsg.AddParticipationNum = easygo.NewInt64(reqMsg.GetAddParticipationNum())
		reqMsg.AddFansNum = easygo.NewInt64(reqMsg.GetAddFansNum())
		// 前端添加话题时,没有带#号,服务器需要帮忙添加,因为数据库存的是带#号的
		reqMsg.Name = easygo.NewString(name)
		msg = fmt.Sprintf("添加话题:%d", reqMsg.GetId())
	}
	if reqMsg.GetAdmin() != "" && reqMsg.GetAdmin() != user.GetAccount() {
		admin := QueryManage(reqMsg.GetAdmin())
		if admin == nil {
			return easygo.NewFailMsg("官方管理员账号不存在")
		}
	}
	if reqMsg.Owner != nil && reqMsg.GetOwner() != "" {
		player := QueryPlayerbyAccount(reqMsg.GetOwner())
		if player == nil {
			return easygo.NewFailMsg("用户管理员账号不存在")
		}
		reqMsg.TopicMaster = easygo.NewInt64(player.GetPlayerId())
		reqMsg.IsOpen = easygo.NewBool(false)
	}
	// else {
	// 	reqMsg.TopicMaster = easygo.NewInt64(0)
	// 	reqMsg.IsOpen = easygo.NewBool(true)
	// }

	// if reqMsg.IsRecommend != nil && reqMsg.IsHot != nil {
	// 	easygo.NewFailMsg("推荐和热门不能同时设置")
	// }

	if reqMsg.GetStatus() == 2 {
		reqMsg.IsRecommend = easygo.NewBool(false)
		reqMsg.IsHot = easygo.NewBool(false)
	}

	reqMsg.UpdateTime = easygo.NewInt64(util.GetMilliTime())
	topicType := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_TYPE, bson.M{"_id": reqMsg.GetTopicTypeId()})
	reqMsg.TopicClass = easygo.NewInt32(topicType.(bson.M)["TopicClass"])

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.TOPIC_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询话题列表
func (self *cls4) RpcGetTopicByIds(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	list, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, bson.M{"_id": bson.M{"$in": reqMsg.GetIds64()}}, 0, 0)
	var lis []*share_message.Topic
	for _, li := range list {
		one := &share_message.Topic{}
		for_game.StructToOtherStruct(li, one)
		onenew := &share_message.Topic{
			Id:   easygo.NewInt64(one.GetId()),
			Name: easygo.NewString(one.GetName()),
		}
		lis = append(lis, onenew)
	}
	msg := &brower_backstage.TopicResponse{
		List:      lis,
		PageCount: easygo.NewInt32(count),
	}

	return msg
}

//话题主申请列表 Status:0-待审核,1已通过,2已拒绝 Type:1-话题名称,2-申请人柠檬号 TimeType:1-申请时间,2-审核时间
func (s *cls4) RpcApplyTopicMasterList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{}
	if reqMsg.GetSort() != "" {
		sort = append(sort, reqMsg.GetSort())
	} else {
		sort = append(sort, "-CreateTime")
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			findBson["TopicName"] = reqMsg.GetKeyword()
		case 2:
			findBson["PlayerAccount"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() < 1000 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		timeBson := bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		switch reqMsg.GetTimeType() {
		case 1:
			findBson["CreateTime"] = timeBson
		case 2:
			findBson["UpdateTime"] = timeBson
		}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_APPLY_TOPIC_MASTER, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.ApplyTopicMaster
	for _, li := range lis {
		one := &share_message.ApplyTopicMaster{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.ApplyTopicMasterRes{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//话题主申请审核
func (s *cls4) RpcApplyTopicMaster(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.GetId64() == 0 || reqMsg.Id64 == nil {
		return easygo.NewFailMsg("Id64参数不能为空")
	}
	//查询待审核的申请是否存在
	one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_APPLY_TOPIC_MASTER, bson.M{"_id": reqMsg.GetId64(), "Status": 0})
	if one == nil {
		return easygo.NewFailMsg("未找到需要审核的申请")
	}

	status := 0
	switch reqMsg.GetId32() {
	case 1: //通过
		status = 1
	case 2: //拒绝
		status = 2
	default:
		return easygo.NewFailMsg("审核状态错误")
	}

	if status == 1 {
		apply := &share_message.ApplyTopicMaster{}
		for_game.StructToOtherStruct(one, apply)
		reqToHall := &server_server.PlayerSI{
			PlayerId: easygo.NewInt64(apply.GetPlayerId()),
			Account:  easygo.NewString(apply.GetTopicName()),
		}
		easygo.Spawn(ChooseOneHall, int32(0), "RpcChangeTopicOwner", reqToHall)
		findBson := bson.M{"_id": apply.GetTopicId()}
		upBson := bson.M{"$set": bson.M{"TopicMaster": apply.GetPlayerId(), "Owner": apply.GetPlayerAccount(), "IsOpen": false}}
		result := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, findBson, upBson, false)
		if result != nil {
			findBson := bson.M{"_id": reqMsg.GetId64()}
			upBson := bson.M{"$set": bson.M{"Status": status, "UpdateTime": easygo.NowTimestamp(), "Operator": user.GetAccount()}}
			resultA := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_APPLY_TOPIC_MASTER, findBson, upBson, false)
			if resultA != nil {
				for_game.UpdateAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_APPLY_TOPIC_MASTER, bson.M{"TopicId": apply.GetTopicId(), "Status": 0}, bson.M{"$set": bson.M{"Status": 2, "UpdateTime": easygo.NowTimestamp(), "Operator": "system"}})
			}
		}
	} else {
		findBson := bson.M{"_id": reqMsg.GetId64()}
		upBson := bson.M{"$set": bson.M{"Status": status, "UpdateTime": easygo.NowTimestamp(), "Operator": user.GetAccount()}}
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_APPLY_TOPIC_MASTER, findBson, upBson, false)
	}

	msg := fmt.Sprintf("话题主申请[%d]审核", reqMsg.GetId64())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.TOPIC_MANAGE, msg)

	return easygo.EmptyMsg
}

//话题修改审核列表 Status:0-待审核,1已通过,2已拒绝 Type:1-话题名称,2-申请人柠檬号 TimeType:1-提交时间,2-审核时间
func (s *cls4) RpcQueryTopicApplyList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{}
	if reqMsg.GetSort() != "" {
		sort = append(sort, reqMsg.GetSort())
	} else {
		sort = append(sort, "-CreateTime")
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			findBson["TopicName"] = reqMsg.GetKeyword()
		case 2:
			findBson["PlayerAccount"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() < 1000 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		timeBson := bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		switch reqMsg.GetTimeType() {
		case 1:
			findBson["CreateTime"] = timeBson
		case 2:
			findBson["UpdateTime"] = timeBson
		}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_APPLY_LOG, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.ApplyEditTopicInfo
	for _, li := range lis {
		one := &share_message.ApplyEditTopicInfo{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.TopicApplyListRes{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//Id查询修改话题
func (s *cls4) RpcQueryTopicApply(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	id := bson.ObjectIdHex(reqMsg.GetIdStr())
	appliy := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_APPLY_LOG, bson.M{"_id": id})
	new := &share_message.ApplyEditTopicInfo{}
	for_game.StructToOtherStruct(appliy, new)
	topic := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, bson.M{"_id": new.GetTopicId()})
	old := &share_message.Topic{}
	for_game.StructToOtherStruct(topic, old)
	return &brower_backstage.QueryTopicApplyRes{
		New: new,
		Old: old,
	}
}

//审核话题申请
func (s *cls4) RpcAuditTopicApply(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.AuditTopicApplyReq) easygo.IMessage {
	if reqMsg.GetId() == "" || reqMsg.Id == nil {
		return easygo.NewFailMsg("Id参数不能为空")
	}
	id := bson.ObjectIdHex(reqMsg.GetId())
	one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_APPLY_LOG, bson.M{"_id": id})
	appliy := &share_message.ApplyEditTopicInfo{}
	for_game.StructToOtherStruct(one, appliy)
	msg := ""
	status := reqMsg.GetStatus()
	switch status {
	case 1: //通过
		findBson := bson.M{"_id": appliy.GetTopicId()}
		upBson := bson.M{"$set": bson.M{"Description": appliy.GetDescription(), "TopicRule": appliy.GetTopicRule(), "HeadURL": appliy.GetHeadURL(), "BgUrl": appliy.GetBgUrl(), "UpdateTime": util.GetMilliTime()}}
		result := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, findBson, upBson, false)
		if result != nil {
			appliyBson := bson.M{"_id": id}
			appliyUpBson := bson.M{"$set": bson.M{"Status": status, "UpdateTime": easygo.NowTimestamp(), "Operator": user.GetAccount()}}
			resultA := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_APPLY_LOG, appliyBson, appliyUpBson, false)
			if resultA != nil {
				msg = fmt.Sprintf("审核通过话题[%d]修改申请", appliy.GetTopicId())
			}
		}
	case 2: //拒绝
		findBson := bson.M{"_id": id}
		upBson := bson.M{"$set": bson.M{"Status": status, "Reason": reqMsg.GetNote(), "UpdateTime": easygo.NowTimestamp(), "Operator": user.GetAccount()}}
		result := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_APPLY_LOG, findBson, upBson, false)
		if result != nil {
			msg = fmt.Sprintf("审核拒绝话题[%d]修改申请", appliy.GetTopicId())
		}
	default:
		return easygo.NewFailMsg("审核状态错误")
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.TOPIC_MANAGE, msg)

	return easygo.EmptyMsg
}
