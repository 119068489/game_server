package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"

	"github.com/astaxie/beego/logs"
)

/*===========================================================================================notify下行通知后台处理===*/
//接收大厅服务器推送新数据通知
func AllPush() {
	list := for_game.GetRedisAdminList()
	var ids []int64

	for _, v := range list {
		if v.ServerId == PServerInfo.GetSid() {
			ids = append(ids, v.UserId)
		}
	}

	for _, id := range ids {
		ep := BrowerEpMp.LoadEndpoint(id)
		if ep != nil {
			ep.RpcNewPush(nil)
		}
	}
}

/*===========================================================================================后台上行通知notify处理===*/
//从大厅服务器获取玩家基础信息
func GetPlayerBaseForNotify(id int64) *share_message.PlayerBase {
	ids := &server_server.PlayerSI{
		PlayerId: easygo.NewInt64(id),
	}
	playerBase := ChooseOneHall(0, "RpcGetPlayerBase", ids)
	return playerBase.(*share_message.PlayerBase)
}

//通知大厅修改玩家数据
//func EditPlayerForHall(id PLAYER_ID) {
//	msg := &backstage_hall.PlayerSI{
//		PlayerId: easygo.NewInt64(id),
//	}
//	BroadCastToAllHall("RpcPlayerChangeHall", msg)
//}

//通知大厅修改群基本数据
func EditTeamForHall(reqMsg *share_message.TeamData) {
	msg := &server_server.EditTeam{
		Id:          reqMsg.Id,
		Name:        reqMsg.Name,
		MaxMember:   reqMsg.MaxMember,
		GongGao:     reqMsg.GongGao,
		IsRecommend: reqMsg.IsRecommend,
		Level:       reqMsg.Level,
	}
	// BroadCastToAllHall("RpcTeamChangeHall", msg)
	ChooseOneHall(0, "RpcTeamChangeHall", msg)
}

//通知大厅解散群
func DefunctTeamForHall(reqMsg *brower_backstage.QueryDataById) {
	msg := &server_server.PlayerSI{
		PlayerId: reqMsg.Id64,
		Account:  easygo.NewString(reqMsg.GetIdStr()),
	}
	// BroadCastToAllHall("RpcDefunctTeamHall", msg)
	ChooseOneHall(0, "RpcDefunctTeamHall", msg)
}

//通知大厅增减群成员
func TeamMemberOptForHall(user *share_message.Manager, reqMsg *brower_backstage.MemberOptRequest) {
	PlayerIds := make([]int64, 0)
	for _, item := range reqMsg.GetAccount() {
		player := QueryPlayerbyAccount(item)
		PlayerIds = append(PlayerIds, player.GetPlayerId())
	}
	team := QueryTeambyId(reqMsg.GetTeamId())

	if team == nil {
		logs.Error("群[%d]不存在", reqMsg.GetTeamId())
		return
	}

	msg := &server_server.MemberOptRequest{
		TeamId:    reqMsg.TeamId,
		PlayerIds: PlayerIds,
		Types:     reqMsg.Types,
		AdminID:   user.Id,
		PlayerID:  team.Owner,
	}
	// BroadCastToAllHall("RpcTeamMemberOptHall", msg)
	ChooseOneHall(0, "RpcTeamMemberOptHall", msg)

}

//回复用户投诉到大厅
func ReplyPlayerComplaintForHall(reqMsg *share_message.PlayerComplaint) {
	SendToPlayer(reqMsg.GetPlayerId(), "RpcResponeComplainInfo", reqMsg)
}

//警告群主
func WarnLordForHall(reqMsg *brower_backstage.QueryDataByIds) {
	teams := QueryTeambyIds(reqMsg.GetIds64())
	var msg string
	for _, t := range teams {
		msg = fmt.Sprintf(reqMsg.GetNote(), t.GetTeamChat())
		reqMsg := &share_message.PlayerComplaint{
			PlayerId: easygo.NewInt64(t.GetOwner()),
			Content:  easygo.NewString(msg),
		}
		SendToPlayer(reqMsg.GetPlayerId(), "RpcWarnLordToHall", reqMsg)
	}
}

//推送app通知到大厅
func SendAppPushforHall(reqMsg *share_message.AppPushMessage) {
	logs.Info("进入通知推送")
	if reqMsg.GetStatus() == 1 {
		ChooseOneHall(0, "RpcSendAppPushHall", reqMsg)
		reqMsg.Status = easygo.NewInt32(2)
		EditAppPushMessage(reqMsg)
	}
	TimerMgr.DelTimerList(int64(reqMsg.GetId()))
}

//推送推文
func SendTweets(id int64) {
	reqMsg := QueryTweetsById(id)
	articleResponseList := []*share_message.Article{}
	now := easygo.NewInt64(util.GetMilliTime())
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
			ObjectId:        easygo.NewInt64(article.GetObjectId()),
		}
		articleResponseList = append(articleResponseList, art)
	}

	tweets := &share_message.Tweets{
		ID:          easygo.NewInt64(reqMsg.GetID()),
		List:        reqMsg.GetList(),
		UserType:    easygo.NewInt32(reqMsg.GetUserType()),
		SendState:   easygo.NewInt32(reqMsg.GetSendState()),
		CreateTime:  easygo.NewInt64(reqMsg.GetCreateTime()),
		Article:     articleResponseList,
		State:       easygo.NewInt32(reqMsg.GetState()),
		UpdateTime:  easygo.NewInt64(now),
		AllLabel:    easygo.NewInt32(reqMsg.GetAllLabel()),
		SendTime:    easygo.NewInt64(now), //发送时间
		CatchLabel:  reqMsg.GetCatchLabel(),
		CustomLabel: reqMsg.GetCustomLabel(),
		JgPush:      easygo.NewInt32(reqMsg.GetJgPush()),
		Validity:    easygo.NewFloat64(reqMsg.GetValidity()),
	}
	if tweets.GetState() == 0 {
		ChooseOneHall(0, "RpcEditArticle", tweets)
		tweets.State = easygo.NewInt32(1) //已发送
		EditPushTweetsState(tweets)
	}
	DelUserTimeTweets(tweets.GetValidity(), []int64{tweets.GetID()})
	ArticleTimeMgr.DelTimerList(tweets.GetID())
}
