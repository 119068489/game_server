package sport_common_dal

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"time"
)

func GetMatchInfo(uniqueGameId int64) *share_message.UniqueGameInfo {
	if uniqueGameId < 1 {
		return nil
	}
	v := GetESPortsGameItem(uniqueGameId)
	if v == nil || v.GetId() < 1 {
		return nil
	}

	tan := "A队未知名称"
	if v.GetTeamA() != nil {
		tan = v.GetTeamA().GetName()
	}
	tbn := "B队未知名称"
	if v.GetTeamB() != nil {
		tbn = v.GetTeamB().GetName()
	}
	//matchVsName := for_game.GetMatchVSName(v.GetMatchName(), v.GetMatchStage(), v.GetBo(), tan, tbn)

	gameInfo := &share_message.UniqueGameInfo{
		MatchName:  easygo.NewString(v.GetMatchName()),
		MatchStage: easygo.NewString(v.GetMatchStage()),
		Bo:         easygo.NewString(v.GetBo()),
		TeamAName:  easygo.NewString(tan),
		TeamBName:  easygo.NewString(tbn),
	}
	return gameInfo
}

//参数参数顺序： 用户Id，标题，正文，赛事ID，游戏标签Id，视频url，封面图
func ApplyMylive(playerId int64, title, content string, uniqueGameId, applabelId int64, videoUrl, imageUrl string, status int32, note, uniqueGameName string) (int32, string) {
	applabelname := for_game.LabelToESportNameMap[int32(applabelId)]

	info := GetMyliveInfoByPlayerId(playerId)
	var gameInfo *share_message.UniqueGameInfo = nil
	if uniqueGameId > 0 {
		gameInfo = GetMatchInfo(uniqueGameId) //获取比赛和队伍信息
	}
	pobj := for_game.GetRedisPlayerBase(playerId)
	if pobj == nil {
		logs.Info("用戶", playerId, "redis信息不存在")
		return for_game.C_SYS_ERROR, "用戶不存在"
	}
	ptypes := int32(0)
	minfo := pobj.GetMyInfo()
	if minfo != nil {
		ptypes = minfo.GetTypes()
	} else {
		logs.Info("用戶", playerId, "用戶額外信息不存在")
		return for_game.C_SYS_ERROR, "用戶額外信息不存在"
	}
	if info == nil {
		return CreateTableESPortsVideoInfo(&share_message.TableESPortsVideoInfo{
			Status:             easygo.NewInt32(status), //0未发布(未审核) 1已发布(审核通过) 2已删除(审核拒绝) 3已禁用 4已过期
			CoverImageUrl:      easygo.NewString(imageUrl),
			Title:              easygo.NewString(title),
			VideoUrl:           easygo.NewString(videoUrl),
			AuthorPlayerId:     easygo.NewInt64(playerId),
			AuthorAccount:      easygo.NewString(""),
			Author:             easygo.NewString(""),
			DataSource:         easygo.NewString(""),
			LookCount:          easygo.NewInt32(0),
			LookCountSys:       easygo.NewInt32(0),
			ThumbsUpCount:      easygo.NewInt32(0),
			ThumbsUpCountSys:   easygo.NewInt32(0),
			AppLabelID:         easygo.NewInt64(applabelId),
			AppLabelName:       easygo.NewString(applabelname),
			BeginEffectiveTime: easygo.NewInt64(0),
			EffectiveType:      easygo.NewInt64(1), //// 1立刻有效 2 定时有效
			VideoType:          easygo.NewInt64(2), // 1视频 2 直播（放映厅）
			IsRecommend:        easygo.NewInt32(0),
			IsHot:              easygo.NewInt32(0),
			UniqueGameId:       easygo.NewInt64(uniqueGameId),
			UniqueGameName:     easygo.NewString(uniqueGameName),
			MenuId:             easygo.NewInt32(for_game.ESPORTMENU_LIVE),
			CommentCount:       easygo.NewInt32(0),
			LabelIds:           []int64{},
			Note:               easygo.NewString(note),
			Content:            easygo.NewString(content),
			UniqueGameInfo:     gameInfo,
			AuthorPlayerType:   easygo.NewInt32(ptypes),
			Operator:           easygo.NewString(""),
		})
	} else if info.GetId() > 0 {
		info.BeginEffectiveTime = easygo.NewInt64(time.Now().Unix())
		//info.Status = easygo.NewInt32(status) //0未发布(未审核) 1已发布(审核通过) 2已删除(审核拒绝) 3已禁用 4已过期
		info.CoverImageUrl = easygo.NewString(imageUrl)
		info.Title = easygo.NewString(title)
		info.VideoUrl = easygo.NewString(videoUrl)
		info.AppLabelID = easygo.NewInt64(applabelId)
		info.AppLabelName = easygo.NewString(applabelname)
		info.UniqueGameId = easygo.NewInt64(uniqueGameId)
		info.UniqueGameName = easygo.NewString(uniqueGameName)
		info.MenuId = easygo.NewInt32(for_game.ESPORTMENU_LIVE)
		info.Content = easygo.NewString(content)
		info.UniqueGameInfo = gameInfo
		return UpdateMyLiveInfo(info)
	}
	return for_game.C_SYS_ERROR, "系统错误"
}

func CreateTableESPortsVideoInfo(info *share_message.TableESPortsVideoInfo) (int32, string) {
	id := for_game.NextId(for_game.TABLE_ESPORTS_VIDEO)
	col, closeFun := GetC(for_game.TABLE_ESPORTS_VIDEO)
	defer closeFun()
	info.CreateTime = easygo.NewInt64(time.Now().Unix())
	info.UpdateTime = easygo.NewInt64(time.Now().Unix())
	info.BeginEffectiveTime = easygo.NewInt64(time.Now().Unix())
	info.Id = easygo.NewInt64(id)
	logs.Info("=====創建放映廳======", info)
	_, err := col.Upsert(bson.M{"_id": id}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "申请成功"
}
func DeleteTableESPortsVideoInfo(id int64) bool {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_VIDEO)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": id})
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}
func UpdateMyLiveInfo(info *share_message.TableESPortsVideoInfo) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_VIDEO)
	defer closeFun()
	info.UpdateTime = easygo.NewInt64(time.Now().Unix())
	_, err := col.Upsert(bson.M{"_id": info.GetId()}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	if err == nil {
		return for_game.C_OPT_SUCCESS, "申请成功"
	} else {
		return for_game.C_INFO_NOT_EXISTS, "数据不存在，修改失败"
	}
}

func UpdateTableESPortsVideoInfo(info *share_message.TableESPortsVideoInfo) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_VIDEO)
	defer closeFun()
	updatedata := bson.M{}
	updatedata["UpdateTime"] = time.Now().Unix()
	if info.IssueTime != nil {
		updatedata["IssueTime"] = info.GetIssueTime()
	}
	if info.Status != nil {
		updatedata["Status"] = info.GetStatus()
	}
	if info.CoverImageUrl != nil {
		updatedata["CoverImageUrl"] = info.GetCoverImageUrl()
	}
	if info.Title != nil {
		updatedata["Title"] = info.GetTitle()
	}
	if info.VideoUrl != nil {
		updatedata["VideoUrl"] = info.GetVideoUrl()
	}
	if info.AuthorPlayerId != nil {
		updatedata["AuthorPlayerId"] = info.GetAuthorPlayerId()
	}
	if info.AuthorAccount != nil {
		updatedata["AuthorAccount"] = info.GetAuthorAccount()
	}
	if info.Author != nil {
		updatedata["Author"] = info.GetAuthor()
	}
	if info.DataSource != nil {
		updatedata["DataSource"] = info.GetDataSource()
	}
	if info.LookCount != nil {
		updatedata["LookCount"] = info.GetLookCount()
	}
	if info.LookCountSys != nil {
		updatedata["LookCountSys"] = info.GetLookCountSys()
	}
	if info.ThumbsUpCount != nil {
		updatedata["ThumbsUpCount"] = info.GetThumbsUpCount()
	}
	if info.ThumbsUpCountSys != nil {
		updatedata["ThumbsUpCountSys"] = info.GetThumbsUpCountSys()
	}

	if info.BeginEffectiveTime != nil {
		updatedata["BeginEffectiveTime"] = info.GetBeginEffectiveTime()
	}
	if info.EffectiveType != nil {
		updatedata["EffectiveType"] = info.GetEffectiveType()
	}
	if info.VideoType != nil {
		updatedata["VideoType"] = info.GetVideoType()
	}
	if info.IsRecommend != nil {
		updatedata["IsRecommend"] = info.GetIsRecommend()
	}
	if info.IsHot != nil {
		updatedata["IsHot"] = info.GetIsHot()
	}
	if info.UniqueGameId != nil {
		updatedata["GameId"] = info.GetUniqueGameId()
	}
	cinfo, err := col.Upsert(bson.M{"_id": info.GetId()}, bson.M{"$set": updatedata})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	if cinfo.Updated > 0 {
		return for_game.C_OPT_SUCCESS, "修改成功"
	} else {
		return for_game.C_INFO_NOT_EXISTS, "数据不存在，修改失败"
	}
}
func GetTableESPortsVideoInfoById(id int64) *share_message.TableESPortsVideoInfo {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_VIDEO)
	defer closeFun()
	data := &share_message.TableESPortsVideoInfo{}
	err := col.Find(bson.M{"_id": id}).One(&data)
	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return nil
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	UpdateFedAdditionEx(col, "LookCount", id, 1) //访问数+1
	return data
}

///获取我的放映厅信息
func GetMyliveInfoByPlayerId(playerId int64) *share_message.TableESPortsVideoInfo {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_VIDEO)
	defer closeFun()
	data := &share_message.TableESPortsVideoInfo{}
	err := col.Find(bson.M{"AuthorPlayerId": playerId, "VideoType": 2}).One(&data) //只拿直播放映厅的
	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return nil
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return data
}

func GetTableESPortsVideoInfoList1(offset, limit int, sort string, keyword string) ([]*share_message.TableESPortsVideoInfo, int) {
	var list []*share_message.TableESPortsVideoInfo
	col, closeFun := GetC(for_game.TABLE_ESPORTS_VIDEO)
	defer closeFun()
	queryBson := bson.M{}
	if keyword != "" {
		queryBson["Title"] = bson.M{"$regex": "^" + keyword + "+"}
	}
	query := col.Find(queryBson)
	count, err := query.Count()
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	err = query.Sort(sort).Skip(offset).Limit(limit).All(&list)
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	if list == nil {
		list = []*share_message.TableESPortsVideoInfo{}
	}
	return list, count
}
func GetTableESPortsVideoInfoList2(cPage, pSize int, sort string, keyword string) ([]*share_message.TableESPortsVideoInfo, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESPortsVideoInfo
	col, closeFun := GetC(for_game.TABLE_ESPORTS_VIDEO)
	defer closeFun()
	queryBson := bson.M{}
	if keyword != "" {
		queryBson["Title"] = bson.M{"$regex": "^" + keyword + "+"}
	}
	query := col.Find(queryBson)
	count, err := query.Count()
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	err = query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	return list, count
}
func GetTableESPortsVideoItemList(cPage, pSize int32, sort string, typeId, lableId int64, videoType int32, playid int64) ([]*share_message.TableESPortsVideoInfo, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESPortsVideoInfo
	col, closeFun := GetC(for_game.TABLE_ESPORTS_VIDEO)
	defer closeFun()
	queryBson := bson.M{}

	if typeId == 3 {
		queryBson["AppLabelID"] = lableId
	} else if typeId == 2 {
		queryBson["LabelIds"] = bson.M{"$elemMatch": lableId}
	}

	if videoType == for_game.ESPORTS_VIDEO_TYPE_2 {
		queryBson["AuthorPlayerId"] = bson.M{"$ne": playid}
	}
	queryBson["VideoUrl"] = bson.M{"$ne": ""}
	queryBson["VideoType"] = videoType
	queryBson["Status"] = 1
	query := col.Find(queryBson)
	count, err := query.Count()
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	err = query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	return list, count
}

//获取比赛相关的
func GetTableESPortsGameVideoItemList(cPage, pSize int32, sort string, uniqueGameId int64) ([]*share_message.TableESPortsVideoInfo, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESPortsVideoInfo
	col, closeFun := GetC(for_game.TABLE_ESPORTS_VIDEO)
	defer closeFun()
	queryBson := bson.M{}

	queryBson["UniqueGameId"] = uniqueGameId

	queryBson["VideoType"] = for_game.ESPORTS_VIDEO_TYPE_2

	queryBson["Status"] = 1
	query := col.Find(queryBson)
	count, err := query.Count()
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	err = query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	return list, count
}
