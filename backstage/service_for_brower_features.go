// 管理后台为[浏览器]提供的服务
//功能管理

package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//查询推送消息
func (self *cls4) RpcQueryAppPushMessage(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryFeaturesRequest) easygo.IMessage {
	list, count := GetAppPushMessage(reqMsg)
	return &brower_backstage.QueryFeaturesResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//新增修改推送消息
func (self *cls4) RpcEditAppPushMessage(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.AppPushMessage) easygo.IMessage {
	if reqMsg.Content == nil && reqMsg.GetContent() == "" {
		return easygo.NewFailMsg("通知内容不能为空")
	}
	msg := "修改推送消息:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_FEATURES_APPPUSHMSG))
		msg = "添加推送消息:"
	} else {
		TimerMgr.DelTimerList(int64(reqMsg.GetId()))
	}

	reqMsg.Operator = user.Account
	reqMsg.Status = easygo.NewInt32(1)
	if reqMsg.SendTime == nil || reqMsg.GetSendTime() == 0 {
		reqMsg.SendTime = easygo.NewInt64(easygo.NowTimestamp())
	}
	EditAppPushMessage(reqMsg)
	easygo.Spawn(AddAppPush, reqMsg)
	msg = msg + easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询小助手消息
func (self *cls4) RpcQuerySystemNoticeMessage(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryFeaturesRequest) easygo.IMessage {
	list, count := GetSystemNoticeMessage(reqMsg)
	return &brower_backstage.QuerySystemNoticeResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//新增修改小助手消息
func (self *cls4) RpcEditSystemNoticeMessage(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.SystemNotice) easygo.IMessage {
	msg := "修改柠檬小助手消息:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_SYSTEM_NOTICE))
		reqMsg.CreateTime = easygo.NewInt64(util.GetMilliTime())
		reqMsg.SendState = easygo.NewInt32(0)
		msg = "添加柠檬小助手消息:"
	} else {
		data := QuerySystemNoticeMessage(reqMsg.GetId())
		if data == nil {
			s := fmt.Sprintf("找不到ID为[%v]小助手消息", reqMsg.GetId())
			return easygo.NewFailMsg(s)
		}

		reqMsg.SendState = easygo.NewInt32(data.GetSendState())
	}

	if reqMsg.GetState() == 2 {
		if reqMsg.GetSendState() == 0 {
			BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcEditSystemNoticeMessage", reqMsg)
		}
		reqMsg.SendState = easygo.NewInt32(1)
	}

	reqMsg.Operator = user.Account
	reqMsg.EditTime = easygo.NewInt64(util.GetMilliTime())

	EditSystemNoticeMessage(reqMsg)
	msg = msg + easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除小助手消息
func (self *cls4) RpcDelSystemNoticeMessage(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	err := DelDataById(for_game.TABLE_SYSTEM_NOTICE, idList)
	easygo.PanicError(err)

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("删除柠檬小助手消息: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//id查询系统参数设置
func (self *cls4) RpcQuerySysParameterById(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	id := reqMsg.GetIdStr()
	if id == "" {
		return easygo.NewFailMsg("id不能为空")
	}

	if id != for_game.LIMIT_PARAMETER && id != for_game.AVATAR_PARAMETER && id != for_game.INTEREST_PARAMETER &&
		id != for_game.OBJ_MODERATIONS && id != for_game.SQUAREHOT_PARAMETER && id != for_game.WARNING_PARAMETER &&
		id != for_game.TOPICHOT_PARAMETER && id != for_game.PUSH_PARAMETER && id != for_game.ESPORT_PARAMETER && id != for_game.COMMON_PARAMETER {
		return easygo.NewFailMsg("Id错误")
	}

	result := for_game.QuerySysParameterById(id)
	if result == nil {
		result = &share_message.SysParameter{
			Id: easygo.NewString(id),
		}
	}
	return result
}

//修改系统参数设置
func (self *cls4) RpcEditSysParameter(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.SysParameter) easygo.IMessage {
	// logs.Info("========RpcEditSysParameter===== reqMsg=%v", reqMsg)
	msg := ""
	switch reqMsg.GetId() {
	case for_game.LIMIT_PARAMETER:
		msg = "修改系统功能转账参数设置"
	case for_game.AVATAR_PARAMETER:
		msg = "修改系统功能头像参数设置"
	case for_game.INTEREST_PARAMETER:
		msg = "修改系统功能兴趣标签参数设置"
	case for_game.OBJ_MODERATIONS:
		msg = "修改系统功能屏蔽控制参数设置"
	case for_game.SQUAREHOT_PARAMETER:
		msg = "修改系统功能动态热门参数设置"
		if reqMsg.DampRatio == nil || reqMsg.GetDampRatio() <= 0 {
			return easygo.NewFailMsg("递减百分比不能小于等于0")
		}
	case for_game.WARNING_PARAMETER:
		msg = "修改预警参数设置"
	case for_game.TOPICHOT_PARAMETER:
		msg = "修改系统功能话题热门参数设置"
	case for_game.PUSH_PARAMETER:
		msg = "修改系统功能极光推送设置"
	case for_game.ESPORT_PARAMETER:
		msg = "修改电竞系统参数设置"
	case for_game.COMMON_PARAMETER:
		msg = "修改系统功能通用设置"
	default:
		return easygo.NewFailMsg("Id错误")
	}

	result := EditSysParameter(reqMsg)
	if result != nil {
		return result
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询兴趣标签列表查询
func (self *cls4) RpcInterestTypeList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetInterestTypeList(reqMsg)
	msg := &brower_backstage.InterestTypeResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//修改兴趣标签
func (self *cls4) RpcEditInterestType(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.InterestType) easygo.IMessage {
	msg := "修改兴趣标签:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_INTERESTTYPE))
		msg = "添加兴趣标签:"
	}
	reqMsg.UpdateTime = easygo.NewInt64(time.Now().Unix())
	EditInterestType(reqMsg)
	// BroadCastToAllHall("RpcPaySetChangeToHall", nil) //通知大厅重载支付配置
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)

	return easygo.EmptyMsg
}

//兴趣分类列表
func (self *cls4) RpcGetInterestTypeList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	list := GetInterestTypeListNopage()
	return &brower_backstage.KeyValueResponseTag{
		List: list,
	}
}

//兴趣查询
func (self *cls4) RpcInterestTagList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetInterestTagList(reqMsg)
	msg := &brower_backstage.InterestTagResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//兴趣列表
func (self *cls4) RpcGetInterestTagList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list := GetInterestTagListNopage(reqMsg)
	return &brower_backstage.KeyValueResponseTag{
		List: list,
	}
}

//修改兴趣分类
func (self *cls4) RpcEditInterestTag(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.InterestTag) easygo.IMessage {
	msg := "修改兴趣词:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_INTERESTTAG))
		msg = "添加兴趣词:"
	}
	reqMsg.UpdateTime = easygo.NewInt64(time.Now().Unix())
	EditInterestTag(reqMsg)
	// BroadCastToAllHall("RpcPaySetChangeToHall", nil) //通知大厅重载支付配置
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)

	return easygo.EmptyMsg
}

//添加、修改小助手文章消息
func (self *cls4) RpcEditArticle(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.Article) easygo.IMessage {
	msg := "修改柠檬小助手推送文章:"
	logs.Info("修改文章：", reqMsg)
	now := easygo.NewInt64(util.GetMilliTime())

	if reqMsg.ID == nil && reqMsg.GetID() == 0 { //如果新增
		reqMsg.ID = easygo.NewInt64(for_game.NextId(for_game.TABLE_ARTICLE)) //自增ID
		reqMsg.CreateTime = now                                              //初始化创建时间
		reqMsg.ReadedNum = easygo.NewInt64(0)                                //创建就加上阅读数
		reqMsg.ReadingNum = easygo.NewInt64(0)                               //正在阅读数
		reqMsg.ZanNum = easygo.NewInt64(0)                                   //赞数
		msg = "添加柠檬小助手推送文章:"
	} else { //如果修改，需要同步更新推文表和注册推送表

		article := for_game.QueryArticleById(reqMsg.GetID())

		if article.GetIsMain() != reqMsg.GetIsMain() { //判断主次是否修改
			_, count1 := GetTweetsListByArticleId(reqMsg.GetID()) //判断修改文章是否被应用到推文
			_, count2 := GetRegisterPushByArticleId(reqMsg.GetID())
			if count1 != 0 || count2 != 0 {
				return easygo.NewFailMsg("该文章已被使用")
			}
		}
		UpdateTweetsByArticleInfo(article, reqMsg)
		UpdateRegisterPushByArticleInfo(article, reqMsg)
	}

	switch reqMsg.GetLocation() {
	case 7, 9:
		if reqMsg.GetObjectId() == 0 {
			return easygo.NewFailMsg("跳转位置ID必须大于0")
		}
	}

	reqMsg.EditTime = now
	reqMsg.Operator = user.Account //初始化操作者
	if reqMsg.GetLocation() == 9 {
		//如果是商品跳转，查询商品详情
		shop := QueryShopItemById(reqMsg.GetObjectId())
		if shop == nil {
			return easygo.NewFailMsg("指定的商品不存在")
		}
		reqMsg.ObjPlayerId = easygo.NewInt64(shop.GetPlayerId())
	}
	for_game.EditArticle(reqMsg) //消息入库
	msg = msg + easygo.IntToString(int(reqMsg.GetID()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)
	return easygo.EmptyMsg
}

//搜索小助手文章消息
func (self *cls4) RpcQueryArticle(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryArticleOrTweetsRequest) easygo.IMessage {
	list, count := GetArticleList(reqMsg)
	for i, item := range list {
		list[i].CommentNum = easygo.NewInt32(len(for_game.GetArticleComment(item.GetID(), 0, 0)))
	}

	msg := &brower_backstage.QueryArticleResponse{
		List:      list,
		PageCount: &count,
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, "搜索文章")
	return msg
}

//删除小助手文章消息
func (self *cls4) RpcDelArticle(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	err := DelArticleById(reqMsg.GetIds64())
	easygo.PanicError(err)
	var ids string
	ids = easygo.Int64ArrayToString(reqMsg.GetIds64())
	msg := fmt.Sprintf("删除柠檬小助手文章: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)
	return easygo.EmptyMsg
}

//添加、修改小助手推文
func (self *cls4) RpcAddTweets(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.Tweets) easygo.IMessage {
	msg := "修改柠檬小助手推文:"

	articleResponseList := []*share_message.Article{}
	now := easygo.NewInt64(util.GetMilliTime())

	if reqMsg.ID == nil && reqMsg.GetID() == 0 { //如果新增
		reqMsg.ID = easygo.NewInt64(for_game.NextId(for_game.TABLE_TWEETS)) //自增ID
		reqMsg.CreateTime = now                                             //初始化创建时间
		reqMsg.State = easygo.NewInt32(0)                                   //初始化发送状态 未发送
		msg = "添加柠檬小助手推文:"
	} else {
		//如果是修改，一律删除定时任务
		ArticleTimeMgr.DelTimerList(reqMsg.GetID())    //删除小助手定时推送任务
		UserSweetsTimeMgr.DelTimerList(reqMsg.GetID()) //删除有效期的用户推文
	}

	for _, article := range reqMsg.GetArticle() {
		article.Content = easygo.NewString("")
		articleResponseList = append(articleResponseList, article)
	}

	tweets := &share_message.Tweets{
		ID:          easygo.NewInt64(reqMsg.GetID()),
		List:        reqMsg.GetList(),
		UserType:    easygo.NewInt32(reqMsg.GetUserType()),
		SendState:   easygo.NewInt32(reqMsg.GetSendState()),
		CreateTime:  easygo.NewInt64(reqMsg.GetCreateTime()),
		Operator:    easygo.NewString(user.Account),
		Article:     articleResponseList,
		State:       easygo.NewInt32(reqMsg.GetState()),
		UpdateTime:  easygo.NewInt64(now),
		SendTime:    easygo.NewInt64(now), //发送时间
		CatchLabel:  reqMsg.GetCatchLabel(),
		CustomLabel: reqMsg.GetCustomLabel(),
		JgPush:      easygo.NewInt32(reqMsg.GetJgPush()),
		Validity:    easygo.NewFloat64(reqMsg.GetValidity()),
		AllLabel:    easygo.NewInt32(reqMsg.GetAllLabel()),
	}

	// if reqMsg.GetAllLabel() == 0 {
	// 	tweets.AllLabel = easygo.NewInt32(reqMsg.GetAllLabel())
	// }

	if reqMsg.GetSendState() == 1 { //如果立即发送
		ChooseOneHall(0, "RpcEditArticle", tweets)
		DelUserTimeTweets(tweets.GetValidity(), []int64{tweets.GetID()})
		//DelUserTimeTweets(0.1, []int64{tweets.GetID()})
		tweets.State = easygo.NewInt32(1) //已发送
	} else { //定时发送
		tweets.SendTime = easygo.NewInt64(reqMsg.GetSendTime())
		logs.Info("定时发送时间：", tweets.GetSendTime())
		TimePushTweets(tweets.GetSendTime(), tweets.GetID())
	}
	EditTweets(tweets) //消息入库
	msg = msg + easygo.IntToString(int(reqMsg.GetID()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)

	return easygo.EmptyMsg
}

//添加、修改注册推文
func (self *cls4) RpcAddRegisterPush(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.RegisterPush) easygo.IMessage {
	msg := "修改用户注册推文:"
	logs.Info(reqMsg)

	articleResponseList := []*share_message.Article{}
	now := easygo.NewInt64(util.GetMilliTime())

	if reqMsg.ID == nil && reqMsg.GetID() == 0 { //如果新增
		reqMsg.ID = easygo.NewInt64(for_game.NextId(for_game.TABLE_REGISTER_PUSH)) //自增ID
		reqMsg.CreateTime = now                                                    //初始化创建时间
		msg = "添加用户注册推文:"
	}

	for _, article := range reqMsg.GetArticle() {
		art := &share_message.Article{
			ID:              easygo.NewInt64(article.GetID()),
			Title:           easygo.NewString(article.GetTitle()),
			Icon:            easygo.NewString(article.GetIcon()),
			ArticleType:     easygo.NewInt32(article.GetArticleType()),
			Location:        easygo.NewInt32(article.GetLocation()),
			IsMain:          easygo.NewInt32(article.GetIsMain()),
			Sort:            easygo.NewInt32(article.GetSort()),
			Profile:         easygo.NewString(article.GetProfile()),
			TransArticleUrl: easygo.NewString(article.GetTransArticleUrl()),
		}
		articleResponseList = append(articleResponseList, art)
	}

	registerPush := &share_message.RegisterPush{
		ID:          easygo.NewInt64(reqMsg.GetID()),
		List:        reqMsg.GetList(),
		UserType:    easygo.NewInt32(reqMsg.GetUserType()),
		CreateTime:  easygo.NewInt64(reqMsg.GetCreateTime()),
		Operator:    easygo.NewString(user.Account),
		Article:     articleResponseList,
		State:       easygo.NewInt32(reqMsg.GetState()),
		UpdateTime:  easygo.NewInt64(now),
		AllLabel:    easygo.NewInt32(reqMsg.GetAllLabel()),
		CatchLabel:  reqMsg.GetCatchLabel(),
		CustomLabel: reqMsg.GetCustomLabel(),
	}

	EditRegisterPush(registerPush) //消息入库
	msg = msg + easygo.IntToString(int(reqMsg.GetID()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询小助手推文消息
func (self *cls4) RpcQueryTweets(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryArticleOrTweetsRequest) easygo.IMessage {
	list, count := GetTweetsList(reqMsg)
	msg := &brower_backstage.QueryTweetsResponse{
		List:      list,
		PageCount: &count,
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, "搜索文章")
	return msg
}

//查询注册推文
func (self *cls4) RpcQueryRegisterPush(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryArticleOrTweetsRequest) easygo.IMessage {
	list, count := GetRegisterPushList(reqMsg)
	msg := &brower_backstage.QueryRegisterPushResponse{
		List:      list,
		PageCount: &count,
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, "搜索文章")
	return msg
}

//查询自定义标签列表查询
//查询自定义标签列表
func (self *cls4) RpcCustomTagList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetCustomTagList(reqMsg)
	msg := &brower_backstage.CustomTagResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//修改自定义标签
func (self *cls4) RpcEditCustomTag(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.CustomTag) easygo.IMessage {
	msg := "修改自定义标签:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_CUSTOMTAG))
		msg = "添加自定义标签:"
	}
	reqMsg.UpdateTime = easygo.NewInt64(time.Now().Unix())
	EditCustomTag(reqMsg)
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)

	return easygo.EmptyMsg
}

//自定义标签下拉列表
func (self *cls4) RpcGetCustomTagList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	list := GetCustomTagListNopage(reqMsg)
	return &brower_backstage.KeyValueResponseTag{
		List: list,
	}
}

//给玩家加自定义标签
func (self *cls4) RpcToPlayerCustomTag(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	msg := "给玩家打自定义标签:"
	customTag := reqMsg.GetIds32()
	playerids := reqMsg.GetIds64()
	if len(playerids) == 0 {
		return easygo.NewFailMsg("玩家ID不能为空")
	}
	playerid := playerids[0]
	pmg := for_game.GetRedisPlayerBase(playerid)
	pmg.SetRedisCustomTag(customTag)
	pmg.SaveToMongo()

	var item string
	for _, i := range customTag {
		item += easygo.IntToString(int(i))
	}
	msg += item
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)
	return easygo.EmptyMsg
}

//查询抓取标签列表
func (self *cls4) RpcGrabTagList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetGrabTagList(reqMsg)
	msg := &brower_backstage.GrabTagResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//修改抓取标签
func (self *cls4) RpcEditGrabTag(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.GrabTag) easygo.IMessage {
	msg := "修改抓取标签:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_GRABTAG))
		msg = "添加抓取标签:"
	}
	EditGrabTag(reqMsg)
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询抓取词列表
func (self *cls4) RpcCrawlWordsList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetCrawlWordsList(reqMsg)
	msg := &brower_backstage.CrawlWordsResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//修改抓取词
func (self *cls4) RpcEditCrawlWords(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.CrawlWords) easygo.IMessage {
	msg := "修改抓取词:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_CRAWL_WORDS))
		msg = "添加抓取词:"
	}
	if reqMsg.GrabTag == nil && reqMsg.GetGrabTag() == 0 {
		return easygo.NewFailMsg("抓取词标签必选")
	}

	EditCrawlWords(reqMsg)
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.BSPLAYER_MANAGE, msg)

	return easygo.EmptyMsg
}

//抓取标签下拉列表
func (self *cls4) RpcGetGrabTagList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	list := GetGrabTagListNopage()
	return &brower_backstage.KeyValueResponseTag{
		List: list,
	}
}

//个性标签下拉列表
func (self *cls4) RpcGetPersonalityTags(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	list := GetPersonalityTags()
	return &brower_backstage.KeyValueResponseTag{
		List: list,
	}
}

//删除抓取词
func (self *cls4) RpcDelCrawlWords(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds32()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的词")
	}
	DelCrawlWords(idList)

	var ids string
	idsarr := reqMsg.GetIds32()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除抓取词: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.PAY_MANAGE, msg)

	return easygo.EmptyMsg
}

//玩家抓取词列表
func (self *cls4) RpcQueryPlayerWordsList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	list := GetPlayerCrawlWordsList(reqMsg)
	return &brower_backstage.PlayerCrawlWordsResponse{
		List: list,
	}
}

//查询兴趣组合列表
func (self *cls4) RpcInterestGroupList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetInterestGroupList(reqMsg)
	msg := &brower_backstage.InterestGroupResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//修改兴趣组合
func (self *cls4) RpcEditInterestGroup(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.InterestGroup) easygo.IMessage {
	group := reqMsg.GetGroup()
	if len(group) == 0 {
		return easygo.NewFailMsg("请先选择组合")
	}

	easygo.SortSliceInt32(group, true) //排序

	one := for_game.QueryInterestGroupByGroup(group)
	if one != nil && one.GetId() != reqMsg.GetId() {
		return easygo.NewFailMsg("组合已经存在")
	}

	reqMsg.Group = group

	msg := "修改兴趣组合:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_INTERESTGROUP))
		msg = "添加兴趣组合:"
	}
	EditInterestGroup(reqMsg)
	msg += easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除兴趣组合
func (self *cls4) RpcDelInterestGroups(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds32()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的组合")
	}
	DelInterestGroup(idList)

	var ids string
	idsarr := reqMsg.GetIds32()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除兴趣组合: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询屏蔽词
func (self *cls4) RpcQueryDirtyWords(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := QueryDirtyWords(reqMsg)
	msg := &brower_backstage.DirtyWordsResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//删除屏蔽词
func (self *cls4) RpcDelDirtyWords(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	words := reqMsg.GetIdsStr()
	if len(words) == 0 {
		return easygo.NewFailMsg("写点什么吧")
	}

	DelDirtyWords(words)

	var ids string
	idsarr := reqMsg.GetIdsStr()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += idsarr[i] + ","

		} else {
			ids += idsarr[i]
		}
	}
	msg := fmt.Sprintf("批量删除屏蔽词: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//添加屏蔽词
func (self *cls4) RpcAddDirtyWords(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	words := reqMsg.GetIdsStr()
	if len(words) == 0 {
		return easygo.NewFailMsg("写点什么吧")
	}

	AddDirtyWords(words)

	var ids string
	idsarr := reqMsg.GetIdsStr()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += idsarr[i] + ","

		} else {
			ids += idsarr[i]
		}
	}
	msg := fmt.Sprintf("批量添加屏蔽词: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询个性签名
func (self *cls4) RpcQuerySignature(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := for_game.FindAll(for_game.MONGODB_NINGMENG, "signature", bson.M{}, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()))
	var lis []*share_message.Signature
	for _, li := range list {
		one := &share_message.Signature{}
		for_game.StructToOtherStruct(li, one)
		lis = append(lis, one)
	}
	msg := &brower_backstage.SignatureResponse{
		List:      lis,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//批量删除
func (self *cls4) RpcDelSignature(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	words := reqMsg.GetIdsStr()
	if len(words) == 0 {
		return easygo.NewFailMsg("写点什么吧")
	}

	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, "signature", bson.M{"_id": bson.M{"$in": words}})

	var ids string
	idsarr := reqMsg.GetIdsStr()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += idsarr[i] + ","

		} else {
			ids += idsarr[i]
		}
	}
	msg := fmt.Sprintf("批量删除个性签名: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//批量添加
func (self *cls4) RpcAddSignature(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	words := reqMsg.GetIdsStr()
	if len(words) == 0 {
		return easygo.NewFailMsg("写点什么吧")
	}

	for _, w := range words {
		one := &share_message.Signature{
			Title: easygo.NewString(w),
		}
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, "signature", bson.M{"_id": w}, bson.M{"$set": one}, true)
	}

	var ids string
	idsarr := reqMsg.GetIdsStr()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += idsarr[i] + ","

		} else {
			ids += idsarr[i]
		}
	}
	msg := fmt.Sprintf("批量添加个性签名: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询文章评论
func (self *cls4) RpcQueryArticleComment(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetArticleComment(reqMsg)
	for i, item := range list {
		player := QueryPlayerbyId(item.GetPlayerId())
		list[i].Account = player.Account
	}
	return &brower_backstage.ArticleCommentResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//删除文章评论
func (self *cls4) RpcDelArticleComment(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}
	DelArticleComment(idList, reqMsg.GetNote())

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除文章评论: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//附近的人引导列表
func (self *cls4) RpcQueryNearLead(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.GetKeyword() != "" && reqMsg.Keyword != nil {
		switch reqMsg.GetType() {
		case 1:
			id := easygo.StringToInt64noErr(reqMsg.GetKeyword())
			findBson["_id"] = id
		case 2:
			findBson["Name"] = reqMsg.GetKeyword()
		}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_NEAR_LEAD, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()))
	var list []*share_message.NearSet
	for _, li := range lis {
		one := &share_message.NearSet{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.QueryNearSetResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//附近的人引导保存
func (self *cls4) RpcSaveNearLead(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.NearSet) easygo.IMessage {
	if reqMsg.Name == nil && reqMsg.GetName() == "" {
		return easygo.NewFailMsg("名称不能为空")
	}

	if reqMsg.GetWeights() < 1 {
		return easygo.NewFailMsg("权重不能小于1")
	}

	status := int32(2)
	if reqMsg.Status != nil && reqMsg.GetStatus() > 0 {
		status = reqMsg.GetStatus()
	}
	reqMsg.Status = easygo.NewInt32(status)
	msg := fmt.Sprintf("修改附近的人引导项:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_NEAR_LEAD))
		msg = fmt.Sprintf("添加附近的人引导项:%d", reqMsg.GetId())
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_NEAR_LEAD, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除附近的人引导
func (self *cls4) RpcDelNearLead(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}

	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_NEAR_LEAD, bson.M{"_id": bson.M{"$in": idList}})

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除附近的人引导: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}

//附近的人快捷打招呼列表
func (self *cls4) RpcQueryNearFastTerm(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_NEAR_FAST_TERM, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()))
	var list []*share_message.NearSet
	for _, li := range lis {
		one := &share_message.NearSet{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.QueryNearSetResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//附近的人快捷打招呼保存
func (self *cls4) RpcSaveNearFastTerm(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.NearSet) easygo.IMessage {
	if reqMsg.Name == nil && reqMsg.GetName() == "" {
		return easygo.NewFailMsg("名称不能为空")
	}

	if reqMsg.GetWeights() < 1 {
		return easygo.NewFailMsg("价格不能小于1")
	}

	status := int32(2)
	if reqMsg.Status == nil && reqMsg.GetStatus() > 0 {
		status = reqMsg.GetStatus()
	}
	reqMsg.Status = easygo.NewInt32(status)
	msg := fmt.Sprintf("修改附近的人快捷用语:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_NEAR_FAST_TERM))
		msg = fmt.Sprintf("添加附近的人快捷用语:%d", reqMsg.GetId())
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_NEAR_FAST_TERM, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.FEATURES_MANAGE, msg)

	return easygo.EmptyMsg
}
