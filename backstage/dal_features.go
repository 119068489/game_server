//功能管理

package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/client_server"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"strconv"
	"time"

	"github.com/akqp2019/mgo"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//启服加载定时推送任务
func TimedAppPushMessage() {
	logs.Info("加载定时推送任务")
	list := GetAppPushbyStatus(1)
	for _, item := range list {
		AddAppPush(item)
	}
}

//添加推送任务定时器
func AddAppPush(item *share_message.AppPushMessage) {
	if item.GetSendState() == 1 {
		SendAppPushforHall(item)
		return
	}

	triggerTime := time.Duration(item.GetSendTime()-time.Now().Unix()) * time.Second
	if triggerTime > 0 {
		timer := easygo.AfterFunc(triggerTime, func() { SendAppPushforHall(item) })
		TimerMgr.AddTimerList(int64(item.GetId()), timer)
	} else {
		item.Status = easygo.NewInt32(3)
		EditAppPushMessage(item)
	}
}

// 起服将为已过定时时间但未推送的消息进行推送
func TimedAppPushTweets() {
	logs.Info("加载小助手定时推送推文")
	state := easygo.NewInt32(0)
	list := GetTweetsListByState(state)
	//将超过当前时间且处于未发送状态的推文重新发送给玩家
	for _, item := range list {
		TimePushTweets(item.GetSendTime(), item.GetID())
	}
}

//小助手推文定时器
func TimePushTweets(sendTime int64, tweetsId int64) {
	triggerTime := time.Duration((sendTime*int64(time.Millisecond))/1e9-time.Now().Unix()) * time.Second
	logs.Info("发送时间：", triggerTime)
	if triggerTime > 0 {
		timer := easygo.AfterFunc(triggerTime, func() { SendTweets(tweetsId) })
		ArticleTimeMgr.AddTimerList(tweetsId, timer)
	} else {
		SendTweets(tweetsId) //重新发送停服导致发送失败的那部分
	}
}

//根据有效期删除用户推文
func DelUserTimeTweets(validity float64, tweetsId []int64) {
	triggerTime := time.Duration(validity*60*60*1000) * time.Millisecond
	logs.Info("%v小时=%v毫秒 ", validity, triggerTime)
	timer := easygo.AfterFunc(triggerTime, func() {
		RemoveUserTweets(tweetsId)
	})
	UserSweetsTimeMgr.AddTimerList(tweetsId[0], timer)

}

//移除用户未收到的推文
func RemoveUserTweets(tweetsId []int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_TWEETS)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{}, bson.M{"$pullAll": bson.M{"TweetsIdList": tweetsId}})
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
}

//修改小助手推送状态
func EditPushTweetsState(tweets *share_message.Tweets) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TWEETS)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": tweets.GetID()}, bson.M{"$set": bson.M{"State": tweets.GetState()}})
	easygo.PanicError(err)

}

//查询推送消息
func GetAppPushMessage(reqMsg *brower_backstage.QueryFeaturesRequest) ([]*share_message.AppPushMessage, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_FEATURES_APPPUSHMSG)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		queryBson["SendTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	//不查关键词
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetTypes() {
		case 1: //标题
			queryBson["Title"] = reqMsg.GetKeyword()
		case 2: //操作者
			queryBson["Operator"] = reqMsg.GetKeyword()
		default:
			easygo.NewFailMsg("查询条件有误")
		}
	}

	if reqMsg.Recipient != nil {
		queryBson["Recipient"] = reqMsg.GetRecipient()

	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.AppPushMessage
	errc := query.Sort("-SendTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//新增修改推送消息
func EditAppPushMessage(reqMsg *share_message.AppPushMessage) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_FEATURES_APPPUSHMSG)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//根据状态查询推送消息
func GetAppPushbyStatus(status int) []*share_message.AppPushMessage {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_FEATURES_APPPUSHMSG)
	defer closeFun()
	queryBson := bson.M{"Status": status}
	query := col.Find(queryBson)
	var list []*share_message.AppPushMessage
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

//查询小助手消息
func GetSystemNoticeMessage(reqMsg *brower_backstage.QueryFeaturesRequest) ([]*share_message.SystemNotice, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_NOTICE)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		queryBson["create_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	//不查关键词
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetTypes() {
		case 1: //标题
			queryBson["title"] = reqMsg.GetKeyword()
		case 2: //操作者
			queryBson["operator"] = reqMsg.GetKeyword()
		default:
			logs.Info("查询条件有误")
			return []*share_message.SystemNotice{}, 0
		}
	}

	if reqMsg.Recipient != nil {
		queryBson["user_type"] = reqMsg.GetRecipient()

	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		queryBson["state"] = reqMsg.GetStatus()
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.SystemNotice
	errc := query.Sort("-create_time").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//新增修改小助手消息
func EditSystemNoticeMessage(reqMsg *share_message.SystemNotice) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_NOTICE)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

func QuerySystemNoticeMessage(id int64) *share_message.SystemNotice {
	data := &share_message.SystemNotice{}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_NOTICE)
	defer closeFun()
	err := col.Find(bson.M{"_id": id}).One(data)

	if err != nil {
		if err == mgo.ErrNotFound {
			easygo.PanicError(err)
		}
		return nil
	}

	return data
}

//修改系统参数设置
func EditSysParameter(reqMsg *share_message.SysParameter) *base.Fail {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SYS_PARAMETER)
	defer closeFun()
	if reqMsg.GetId() == "" || reqMsg.Id == nil {
		return easygo.NewFailMsg("参数错误")
	}
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)

	if reqMsg.GetId() == for_game.ESPORT_PARAMETER {
		for_game.SetRedisGuessBetRiskControl() //如果是电竞修改到redis
	} else {
		//通知大厅重载配置
		BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcSysParameterChangeToHall", &server_server.SysteamModId{Id: easygo.NewString(reqMsg.GetId())})
	}
	return nil
}

//查询兴趣分类列表查询
func GetInterestTypeList(reqMsg *brower_backstage.ListRequest) ([]*share_message.InterestType, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTTYPE)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Name"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() < 1000 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	var list []*share_message.InterestType
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("Sort").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改兴趣分类
func EditInterestType(reqMsg *share_message.InterestType) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTTYPE)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//兴趣分类列表
func GetInterestTypeListNopage() []*brower_backstage.KeyValueTag {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTTYPE)
	defer closeFun()

	queryBson := bson.M{}
	query := col.Find(queryBson)
	var list []*share_message.InterestType
	err := query.All(&list)
	easygo.PanicError(err)

	var lis []*brower_backstage.KeyValueTag

	for _, i := range list {
		li := &brower_backstage.KeyValueTag{
			//Key:   easygo.NewString(easygo.IntToString(int(i.GetId()))),
			Key:   i.Id,
			Value: i.Name,
		}
		lis = append(lis, li)
	}

	return lis
}

//查询兴趣列表
func GetInterestTagListNopage(reqMsg *brower_backstage.ListRequest) []*brower_backstage.KeyValueTag {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTTAG)
	defer closeFun()

	queryBson := bson.M{}
	if reqMsg.GetStatus() < 1000 {
		queryBson["Status"] = reqMsg.GetStatus()
	}
	if reqMsg.GetListType() > 0 {
		queryBson["InterestType"] = reqMsg.GetListType()
	}
	query := col.Find(queryBson)
	var list []*share_message.InterestTag
	err := query.All(&list)
	easygo.PanicError(err)

	var lis []*brower_backstage.KeyValueTag

	for _, i := range list {
		li := &brower_backstage.KeyValueTag{
			//Key:   easygo.NewString(easygo.IntToString(int(i.GetId()))),
			Key:   i.Id,
			Value: i.Name,
		}
		lis = append(lis, li)
	}

	return lis
}

//兴趣标签表查询
func GetInterestTagList(reqMsg *brower_backstage.ListRequest) ([]*share_message.InterestTag, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTTAG)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{"InterestType": reqMsg.GetType()}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetListType() {
		case 1:
			queryBson["Name"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() < 1000 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	var list []*share_message.InterestTag
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("Sort").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改兴趣标签
func EditInterestTag(reqMsg *share_message.InterestTag) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTTAG)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//查询自定义标签列表
func QueryArticle(articleId int64, contentId int64) *share_message.Article {

	data := &share_message.Article{}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ARTICLE)
	defer closeFun()

	err := col.Find(bson.M{"onumber": articleId}).Select(bson.M{"items": bson.M{"$elemMatch": bson.M{"ino": contentId}}}).One(data)

	if err != nil {
		if err == mgo.ErrNotFound {
			easygo.PanicError(err)
		}
		return nil
	}

	return data
}

//查询小助手文章
func GetArticleList(reqMsg *brower_backstage.QueryArticleOrTweetsRequest) (list []*share_message.Article, PageCount int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ARTICLE)
	defer closeFun()
	var Articlelist []*share_message.Article
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)
	// 关键词查询账号
	queryBson := bson.M{}

	if reqMsg.Querytype != nil && reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		flag := reqMsg.GetQuerytype()
		switch flag {
		case 1:
			//queryBson["Title"] = reqMsg.GetKeyword()
			queryBson["Title"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "i"}}
		case 2:
			queryBson["Operator"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.State != nil && reqMsg.GetState() != 0 {
		queryBson["State"] = reqMsg.GetState()
	}

	if reqMsg.IsMain != nil && reqMsg.GetIsMain() != 0 {
		queryBson["IsMain"] = reqMsg.GetIsMain()
	}

	if reqMsg.ArticleType != nil && reqMsg.GetArticleType() != 0 {
		queryBson["ArticleType"] = reqMsg.GetArticleType()
	}

	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		queryBson["Edit_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.State == nil && reqMsg.GetState() != 0 {
		queryBson["State"] = reqMsg.GetState()
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	err = query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&Articlelist)

	if err != nil && err == mgo.ErrNotFound {
		return nil, 0
	}
	easygo.PanicError(err)
	return Articlelist, int32(count)
}

//根据发送状态和时间查询文章
func GetTweetsListByState(state *int32) (list []*share_message.Tweets) {
	articleUrl := easygo.YamlCfg.GetValueAsString("CLIENT_ARTICLE_URL") //测试服
	var tweets []*share_message.Tweets
	tweetsMessage := []*share_message.Tweets{}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TWEETS)
	defer closeFun()
	queryBson := bson.M{"State": state}
	err := col.Find(queryBson).All(&tweets)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}

	for _, tweet := range tweets {

		articleInfos := tweet.GetArticle()
		articleIds := []int64{}

		for _, articleid := range articleInfos {
			articleIds = append(articleIds, articleid.GetID())
		}

		articleList := QueryArticleByIds(articleIds)

		articles := []*client_server.ArticleResponse{}

		for _, article := range articleList {
			//articleAdd := "http://192.168.150.194:8080/article.html?id=" + strconv.FormatInt(article.GetID(), 10) //本地
			articleAdd := articleUrl + "?id=" + strconv.FormatInt(article.GetID(), 10) //测试服
			if article.GetTransArticleUrl() != "" && article.TransArticleUrl != nil {
				articleAdd = article.GetTransArticleUrl()
			}

			art := &client_server.ArticleResponse{
				Id:          easygo.NewInt64(article.GetID()),
				Title:       easygo.NewString(article.GetTitle()),
				Icon:        easygo.NewString(article.GetIcon()),
				ArticleAdd:  easygo.NewString(articleAdd),
				ArticleType: easygo.NewInt32(article.GetArticleType()),
				Location:    easygo.NewInt32(article.GetLocation()),
				IsMain:      easygo.NewInt32(article.GetIsMain()),
				Profile:     easygo.NewString(article.GetProfile()),
			}
			articles = append(articles, art)
		}

		articleListResponse := &share_message.Tweets{
			ID:         easygo.NewInt64(tweet.GetID()),
			List:       tweet.GetList(),
			UserType:   easygo.NewInt32(tweet.GetUserType()),
			SendState:  easygo.NewInt32(tweet.GetState()),
			SendTime:   easygo.NewInt64(tweet.GetSendTime()),
			CreateTime: easygo.NewInt64(tweet.GetCreateTime()),
			Operator:   easygo.NewString(tweet.GetOperator()),
			State:      easygo.NewInt32(tweet.GetState()),
			UpdateTime: easygo.NewInt64(tweet.GetUpdateTime()),
		}

		tweetsMessage = append(tweetsMessage, articleListResponse)

	}
	return tweetsMessage
}

//删除小助手文章
func DelArticleById(idList []int64) error {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ARTICLE)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": idList}})
	return err
}

//根据id批量查询文章
func QueryArticleByIds(ids []int64) (articleList []*share_message.Article) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ARTICLE)
	defer closeFun()
	var articles []*share_message.Article
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&articles)

	if err != nil && err == mgo.ErrNotFound {
		return nil
	} else {
		easygo.PanicError(err)
	}
	return articles
}

//根据id查询推文
func QueryTweetsById(ids int64) *share_message.Tweets {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TWEETS)
	defer closeFun()
	var tweets *share_message.Tweets
	err := col.Find(bson.M{"_id": ids}).One(&tweets)

	if err != nil && err == mgo.ErrNotFound {
		return nil
	} else {
		easygo.PanicError(err)
	}
	return tweets
}

//添加、修改小助手推文
func EditTweets(reqMsg *share_message.Tweets) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TWEETS)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetID()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//获取推文列表
func GetTweetsList(reqMsg *brower_backstage.QueryArticleOrTweetsRequest) (list []*share_message.Tweets, PageCount int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TWEETS)
	defer closeFun()
	var tweetsList []*share_message.Tweets
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}

	// 关键词查询
	if reqMsg.Querytype != nil && reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		flag := reqMsg.GetQuerytype()
		switch flag {
		case 1:
			queryBson["Article.Title"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "i"}}
		case 2:
			queryBson["Operator"] = reqMsg.GetKeyword()
		}
	}

	// 状态

	if reqMsg.State != nil {
		queryBson["State"] = reqMsg.GetState()
	}

	//推送端
	//全部就是nil 、 0全体，1 IOS，2 Android
	if reqMsg.ArticleType != nil {
		queryBson["User_type"] = reqMsg.GetArticleType()
	}

	// 推送时间
	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		queryBson["Send_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	err = query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&tweetsList)

	if err != nil && err == mgo.ErrNotFound {
		return nil, 0
	}

	easygo.PanicError(err)
	return tweetsList, int32(count)
}

//查询注册推文
func GetRegisterPushList(reqMsg *brower_backstage.QueryArticleOrTweetsRequest) (list []*share_message.RegisterPush, PageCount int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REGISTER_PUSH)
	defer closeFun()
	var pushList []*share_message.RegisterPush
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}

	// 关键词查询
	if reqMsg.Querytype != nil && reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		flag := reqMsg.GetQuerytype()
		switch flag {
		case 1:
			queryBson["Article.Title"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "i"}}
		case 2:
			queryBson["Operator"] = reqMsg.GetKeyword()
		}
	}

	// 状态

	if reqMsg.State != nil {
		queryBson["State"] = reqMsg.GetState()
	}

	//推送端
	//全部就是nil 、 0全体，1 IOS，2 Android
	if reqMsg.ArticleType != nil {
		queryBson["User_type"] = reqMsg.GetArticleType()
	}

	// 推送时间
	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		queryBson["Create_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	err = query.Sort("-Create_time").Skip(curPage * pageSize).Limit(pageSize).All(&pushList)

	if err != nil && err == mgo.ErrNotFound {
		return nil, 0
	}

	easygo.PanicError(err)
	return pushList, int32(count)
}

//根据文章信息更新推文
func UpdateTweetsByArticleInfo(old *share_message.Article, new *share_message.Article) {
	if old.Title != new.Title || old.ArticleType != new.ArticleType {
		col2, closeFun2 := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TWEETS)
		defer closeFun2()
		selector := bson.M{"Article._id": new.GetID()}
		conditions := bson.M{"Article.$.ArticleType": new.GetArticleType(), "Article.$.Title": new.GetTitle()}
		update := bson.M{"$set": conditions}
		if new.GetIcon() != "" && new.Icon != nil {
			conditions["Article.$.Icon"] = new.GetIcon()
		}

		if new.Location != nil {
			conditions["Article.$.Location"] = new.GetLocation()
		}

		if new.GetTransArticleUrl() != "" && new.TransArticleUrl != nil {
			conditions["Article.$.TransArticleUrl"] = new.GetTransArticleUrl()
		}

		if new.Sort != nil {
			conditions["Article.$.Sort"] = new.GetSort()
		}

		if new.Profile != nil && new.GetProfile() != "" {
			conditions["Article.$.Profile"] = new.GetProfile()
		}
		_, err := col2.UpdateAll(selector, update)
		easygo.PanicError(err)
		//logs.Info("推文同步更新" + strconv.Itoa(changeInfo.Updated))

	}
}

//根据文章信息更新注册推送
func UpdateRegisterPushByArticleInfo(old *share_message.Article, new *share_message.Article) {
	if old.Title != new.Title || old.ArticleType != new.ArticleType {
		col2, closeFun2 := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REGISTER_PUSH)
		defer closeFun2()
		selector := bson.M{"Article._id": new.GetID()}
		conditions := bson.M{"Article.$.ArticleType": new.GetArticleType(), "Article.$.Title": new.GetTitle()}
		update := bson.M{"$set": conditions}
		if new.GetIcon() != "" && new.Icon != nil {
			conditions["Article.$.Icon"] = new.GetIcon()
		}

		if new.Location != nil {
			conditions["Article.$.Location"] = new.GetLocation()
		}

		if new.GetTransArticleUrl() != "" && new.TransArticleUrl != nil {
			conditions["Article.$.TransArticleUrl"] = new.GetTransArticleUrl()
		}

		if new.Sort != nil {
			conditions["Article.$.Sort"] = new.GetSort()
		}

		if new.Profile != nil && new.GetProfile() != "" {
			conditions["Article.$.Profile"] = new.GetProfile()
		}
		_, err := col2.UpdateAll(selector, update)
		easygo.PanicError(err)
		//logs.Info("推文同步更新" + strconv.Itoa(changeInfo.Updated))

	}
}

//添加、修改注册推文
func EditRegisterPush(reqMsg *share_message.RegisterPush) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REGISTER_PUSH)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetID()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//查询自定义标签列表查询
func GetCustomTagList(reqMsg *brower_backstage.ListRequest) ([]*share_message.CustomTag, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CUSTOMTAG)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Name"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() > 0 {
		switch reqMsg.GetListType() {
		case 1:
			queryBson["_id"] = bson.M{"$in": reqMsg.GetCustomTag()}
		case 2:
			queryBson["_id"] = bson.M{"$nin": reqMsg.GetCustomTag()}
		}
	}

	var list []*share_message.CustomTag
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改自定义标签
func EditCustomTag(reqMsg *share_message.CustomTag) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CUSTOMTAG)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//自定义标签列表
func GetCustomTagListNopage(reqMsg *brower_backstage.QueryDataById) []*brower_backstage.KeyValueTag {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CUSTOMTAG)
	defer closeFun()

	queryBson := bson.M{}
	name := reqMsg.GetIdStr()
	if name != "" {
		queryBson["Name"] = name
	}
	query := col.Find(queryBson)
	var list []*share_message.CustomTag
	err := query.All(&list)
	easygo.PanicError(err)

	var lis []*brower_backstage.KeyValueTag

	for _, i := range list {
		li := &brower_backstage.KeyValueTag{
			//Key:   easygo.NewString(easygo.IntToString(int(i.GetId()))),
			Key:   i.Id,
			Value: i.Name,
		}
		lis = append(lis, li)
	}

	return lis
}

func GetTweetsListByArticleId(articleId int64) (list []*share_message.Tweets, count int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TWEETS)
	defer closeFun()

	queryBson := bson.M{"Article._id": articleId}
	var tweetsList []*share_message.Tweets
	query := col.Find(queryBson)
	query.All(&tweetsList)
	count, err := query.Count()
	easygo.PanicError(err)
	return tweetsList, count
}

func GetRegisterPushByArticleId(articleId int64) (list []*share_message.RegisterPush, count int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REGISTER_PUSH)
	defer closeFun()

	queryBson := bson.M{"Article._id": articleId}
	var registerPushList []*share_message.RegisterPush
	query := col.Find(queryBson)
	query.All(&registerPushList)
	count, err := query.Count()
	easygo.PanicError(err)
	return registerPushList, count
}

//查询抓取标签列表
func GetGrabTagList(reqMsg *brower_backstage.ListRequest) ([]*share_message.GrabTag, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_GRABTAG)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	var list []*share_message.GrabTag
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改抓取标签
func EditGrabTag(reqMsg *share_message.GrabTag) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_GRABTAG)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//查询抓取词列表
func GetCrawlWordsList(reqMsg *brower_backstage.ListRequest) ([]*share_message.CrawlWords, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CRAWL_WORDS)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Name"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() > 0 {
		queryBson["GrabTag"] = reqMsg.GetListType()
	}

	var list []*share_message.CrawlWords
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改抓取词
func EditCrawlWords(reqMsg *share_message.CrawlWords) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CRAWL_WORDS)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//抓取标签列表
func GetGrabTagListNopage() []*brower_backstage.KeyValueTag {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_GRABTAG)
	defer closeFun()

	queryBson := bson.M{}
	query := col.Find(queryBson)
	var list []*share_message.GrabTag
	err := query.All(&list)
	easygo.PanicError(err)

	var lis []*brower_backstage.KeyValueTag

	for _, i := range list {
		li := &brower_backstage.KeyValueTag{
			Key:   i.Id,
			Value: i.Name,
		}
		lis = append(lis, li)
	}

	return lis
}

func GetPersonalityTags() []*brower_backstage.KeyValueTag {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CHARACTER_TAG)
	defer closeFun()

	queryBson := bson.M{}
	query := col.Find(queryBson)
	var list []*share_message.InterestTag
	err := query.All(&list)
	easygo.PanicError(err)

	var lis []*brower_backstage.KeyValueTag

	for _, i := range list {
		li := &brower_backstage.KeyValueTag{
			Key:   i.Id,
			Value: i.Name,
		}
		lis = append(lis, li)
	}

	return lis
}

//删除抓取词
func DelCrawlWords(ids []int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CRAWL_WORDS)
	defer closeFun()

	_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	easygo.PanicError(err)
}

//玩家抓取词列表
func GetPlayerCrawlWordsList(reqMsg *brower_backstage.QueryDataById) []*share_message.PlayerCrawlWords {
	id := reqMsg.GetId64()
	var keys string
	var types int32

	if reqMsg.IdStr != nil {
		keys = reqMsg.GetIdStr()
	}

	if reqMsg.Id32 != nil {
		types = reqMsg.GetId32()
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CRAWL_WORDS)
	defer closeFun()
	queryBson := bson.M{"_id": id}
	if keys != "" {
		queryBson["Words.Name"] = keys
	}
	if types != 0 {
		queryBson["Words.GrabTag"] = types
	}
	query := col.Find(queryBson)
	var list []*share_message.PlayerCrawlWords
	err := query.All(&list)
	easygo.PanicError(err)

	for i, l := range list {
		var cw []*share_message.CrawlWords
		for _, w := range l.GetWords() {
			if w.GetCount() > 0 {
				cw = append(cw, w)
			}
		}
		list[i].Words = cw
	}

	return list
}

func GetInterestGroupList(reqMsg *brower_backstage.ListRequest) ([]*share_message.InterestGroup, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTGROUP)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			queryBson["Name"] = reqMsg.GetKeyword()
		}
	}

	var list []*share_message.InterestGroup
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改兴趣组合
func EditInterestGroup(reqMsg *share_message.InterestGroup) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTGROUP)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//删除兴趣组合
func DelInterestGroup(ids []int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_INTERESTGROUP)
	defer closeFun()

	_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	easygo.PanicError(err)
}

//查询屏蔽词
func QueryDirtyWords(reqMsg *brower_backstage.ListRequest) ([]*share_message.DirtyWords, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_DIRTY_WORDS)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	var list []*share_message.DirtyWords
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//删除屏蔽词
func DelDirtyWords(ids []string) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_DIRTY_WORDS)
	defer closeFun()

	_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	easygo.PanicError(err)

	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcEditDirtyWordsToHall", nil) //通知大厅重载屏蔽词
}

//添加屏蔽词
func AddDirtyWords(ids []string) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_DIRTY_WORDS)
	defer closeFun()

	var il []interface{}
	var idc []string
	for _, i := range ids {
		if for_game.IsContainsStr(i, idc) == -1 {
			idc = append(idc, i)
			l := &share_message.DirtyWords{
				Word: easygo.NewString(i),
			}
			il = append(il, l)
		}
	}

	bulk := col.Bulk()
	bulk.RemoveAll(il...)
	bulk.Insert(il...)
	_, err := bulk.Run()
	if err != nil {
		easygo.PanicError(err)
	}

	BroadCastMsgToServerNew(for_game.SERVER_TYPE_HALL, "RpcEditDirtyWordsToHall", nil) //通知大厅重载屏蔽词
}

//查询文章评论
func GetArticleComment(reqMsg *brower_backstage.ListRequest) ([]*share_message.ArticleComment, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ARTICLE_COMMENT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{"ArticleId": reqMsg.GetId()}
	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		// player := &share_message.PlayerBase{}
		switch reqMsg.GetType() {
		case 1:
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			queryBson["PlayerId"] = player.GetPlayerId()
		case 2:
			// player = QuryPlayerByNickname(reqMsg.GetKeyword())
			queryBson["Name"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() > 0 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	var list []*share_message.ArticleComment
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//删除文章评论
func DelArticleComment(ids []int64, note string) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ARTICLE_COMMENT)
	defer closeFun()

	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"Status": for_game.ARTICLE_COMMENT_HIDE, "Note": note}})
	easygo.PanicError(err)
}

//提现订单过期处理
func OrderExpiredDay() {
	sysp := for_game.QuerySysParameterById("limit_parameter")
	day := sysp.GetExpiredDay()
	if day < 1 {
		day = 3
	}
	expiredTime := util.GetMilliTime() - easygo.A_DAY_SECOND*1000*3
	list, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ORDER, bson.M{"ChangeType": 2, "Status": 0, "CreateTime": bson.M{"$lte": expiredTime}}, 0, 0)
	for _, l := range list {
		order := for_game.GetRedisOrderObj(l.(bson.M)["_id"].(string))
		if order == nil || order.GetStatus() > 0 {
			continue
		}
		err := OptOrder(order.GetOrderId(), 3, "system", "提现审核超时自动取消")
		if err != nil {
			logs.Error(err)
			continue
		}
		payNotice := fmt.Sprintf("您在%s发起的提现订单超时，如有疑问请咨询官方客服", easygo.Stamp2Str(order.GetCreateTime()))
		//提现超时通知
		SendSystemNotice(order.GetPlayerId(), "支付通知", payNotice)
	}
}
