package hall

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
	"strconv"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//获取已启用的注册推文
func GetRegisterPush(who *Player, state int32) *client_server.TweetsListResponse {
	articleUrl := easygo.YamlCfg.GetValueAsString("CLIENT_ARTICLE_URL") //测试服
	now := util.GetMilliTime()

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REGISTER_PUSH)
	defer closeFun()

	var registerPushList []*share_message.RegisterPush

	tweetsListResponse := &client_server.TweetsListResponse{}
	articleListRes := make([]*client_server.ArticleListResponse, 0)

	labelBson := []bson.M{bson.M{"AllLabel": 0}, bson.M{"CatchLabel": bson.M{"$in": []int32{who.GetGrabTag()}}}}
	if len(who.GetLabelList()) > 0 {
		labelBson = append(labelBson, bson.M{"List": bson.M{"$in": who.GetLabelList()}})
	}
	logs.Info("用户兴趣标签：", who.GetLabelList())
	logs.Info("用户设备类型:", who.GetDeviceType())
	if len(who.GetCustomTag()) > 0 {
		labelBson = append(labelBson, bson.M{"CustomLabel": bson.M{"$in": who.GetLabelList()}})
	}

	queryBson := bson.M{
		"$and": []bson.M{
			bson.M{"State": state},
			bson.M{"$or": []bson.M{bson.M{"User_type": who.GetDeviceType()}, bson.M{"User_type": 0}}},
			bson.M{"$or": labelBson},
		},
	}

	err := col.Find(queryBson).All(&registerPushList)

	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}

	for _, tweets := range registerPushList { //遍历推文

		tweetsInfo := &client_server.ArticleListResponse{}
		tweetsInfo.TweetsId = easygo.NewInt64(tweets.GetID())
		tweetsInfo.ArticleListId = easygo.NewInt64(tweets.GetUpdateTime())

		articleIds := []int64{}
		articleList1 := tweets.GetArticle()
		for _, article := range articleList1 { //遍历文章
			articleIds = append(articleIds, article.GetID())
		}
		articleList2 := QueryArticleByIds(articleIds) //得到文章列表

		articleList := make([]*client_server.ArticleResponse, 0)
		for _, article := range articleList2 {

			articleAdd := articleUrl + "?id=" + strconv.FormatInt(article.GetID(), 10) + "&t=1&pid=" //测试服
			if article.GetTransArticleUrl() != "" && article.TransArticleUrl != nil {
				articleAdd = article.GetTransArticleUrl()
			}

			articleRes := &client_server.ArticleResponse{
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
			articleList = append(articleList, articleRes)
		}

		tweetsInfo.ArticleList = articleList
		articleListRes = append(articleListRes, tweetsInfo)
	}

	tweetsListResponse.TweetsId = easygo.NewInt64(now)
	tweetsListResponse.TweetsList = articleListRes
	return tweetsListResponse
}

//根据Id查询群成员详情
func QueryTeamMemberById(teamId int64) []string {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAMMEMBER)
	defer closeFun()
	list := []*share_message.PersonalTeamData{}
	err := col.Find(bson.M{"TeamId": teamId}).Select(bson.M{"NickName": 1}).All(&list)
	easygo.PanicError(err)
	var nickNames []string
	for _, item := range list {
		nickNames = append(nickNames, item.GetNickName())
	}
	return nickNames
}

func QueryTeamMemInfoById(teamId int64, playerId int64) *share_message.PersonalTeamData {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAMMEMBER)
	defer closeFun()
	var info *share_message.PersonalTeamData
	err := col.Find(bson.M{"TeamId": teamId, "PlayerId": playerId}).One(&info)
	easygo.PanicError(err)
	return info
}

//根据群id列表批量更新群设置
func UpdateTeanInfoByIds(ids []int64, reqMsg *share_message.OperatorInfo) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$push": bson.M{"OperatorInfo": reqMsg}})
	if err != nil && err == mgo.ErrNotFound {
		easygo.PanicError(err)
	}
}

func DelTeanInfoByIds(ids []int64, flag int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$pull": bson.M{"OperatorInfo": bson.M{"Flag": flag}}, "$set": bson.M{"MessageSetting.IsBan": false, "Status": 0}})
	if err != nil && err == mgo.ErrNotFound {
		easygo.PanicError(err)
	}
}

func GetPlayerComplaint(playerId int64, obj int64, Type int32) *share_message.PlayerComplaint {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_COMPLAINT)
	playerComp := &share_message.PlayerComplaint{}
	defer closeFun()
	queryBson := bson.M{"PlayerId": playerId, "Type": Type}

	switch Type {
	case 2:
		queryBson["RespondentId"] = obj
	case 3:
		queryBson["order_id"] = obj
	case 4:
		queryBson["GoodsId"] = obj
	case 5:
		queryBson["RespondentId"] = obj
	case 7:
		queryBson["DynamicId"] = obj
	}

	err := col.Find(queryBson).Sort("-CreateTime").Limit(1).One(playerComp)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	return playerComp
}
