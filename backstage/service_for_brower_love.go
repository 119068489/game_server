// 管理后台为[浏览器]提供的服务

package backstage

//虚拟市场管理
import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"path"

	"github.com/astaxie/beego/logs"

	"github.com/akqp2019/mgo/bson"
)

const (
	//0-未审核,1-已发布,2-使用中,3-已删除
	PlayerVoiceWork_Status_Untreated int32 = 0
	PlayerVoiceWork_Status_Release   int32 = 1
	PlayerVoiceWork_Status_Use       int32 = 2
	PlayerVoiceWork_Status_Del       int32 = 3
)

//个性标签列表
func (l *cls4) RpcCharacterTagList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-_id"}
	if reqMsg.Status != nil && reqMsg.GetStatus() != 100 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		findBson["Name"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_CHARACTER_TAG, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.InterestTag
	for _, li := range lis {
		one := &share_message.InterestTag{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.InterestTagResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//个性标签保存
func (l *cls4) RpcSaveCharacterTag(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.InterestTag) easygo.IMessage {
	if reqMsg.Name == nil && reqMsg.GetName() == "" {
		return easygo.NewFailMsg("标签内容不能为空")
	}

	if reqMsg.Status == nil && reqMsg.GetStatus() == 0 {
		return easygo.NewFailMsg("标签状态不能为空")
	}

	msg := fmt.Sprintf("修改个性标签:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_CHARACTER_TAG))
		msg = fmt.Sprintf("添加个性标签:%d", reqMsg.GetId())
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_CHARACTER_TAG, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)

	return easygo.EmptyMsg
}

//用户作品管理列表
func (l *cls4) RpcPlayerVoiceWorkList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	logs.Info("RpcPlayerVoiceWorkList:", reqMsg)
	findBson := bson.M{}
	sort := []string{"-_id"}

	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.DownType != nil {
		lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_VIDEO, bson.M{"Type": reqMsg.GetDownType()}, 0, 0)
		if count > 0 {
			var types []int64
			for _, li := range lis {
				types = append(types, li.(bson.M)["_id"].(int64))
			}
			findBson["BgId"] = bson.M{"$in": types}
		}
	}

	//状态查询
	if reqMsg.Status != nil && reqMsg.GetStatus() != 1000 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	//用户类型查询
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["PlayerType"] = reqMsg.GetListType()
	}

	//用户查询
	if reqMsg.Type != nil && reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE, bson.M{"$or": []bson.M{{"Account": reqMsg.GetKeyword()}, {"Phone": reqMsg.GetKeyword()}}})
			if one == nil {
				return easygo.NewFailMsg("用户不存在")
			}

			findBson["PlayerId"] = one.(bson.M)["_id"].(int64)
		case 2:
			findBson["BgId"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		case 3:
			findBson["_id"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_MIX_VOICE_VIDEO, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.PlayerMixVoiceVideo
	for _, li := range lis {
		one := &share_message.PlayerMixVoiceVideo{}
		for_game.StructToOtherStruct(li, one)
		player := QueryPlayerbyId(one.GetPlayerId())
		one.PlayerAccount = easygo.NewString(player.GetAccount())
		bgIface := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_VIDEO, bson.M{"_id": one.GetBgId()})
		if bgIface != nil {
			bg := &share_message.BgVoiceVideo{}
			for_game.StructToOtherStruct(bgIface, bg)
			one.Type = easygo.NewInt32(bg.GetType())
			one.Content = easygo.NewString(bg.GetContent())
		}
		list = append(list, one)
	}
	msg := &brower_backstage.VoiceWorkListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	// logs.Info("msg:", msg)
	return msg
}

//用户作品审核
func (l *cls4) RpcReviewePlayerVoiceWork(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.GetId32() > 3 || reqMsg.GetId32() < 0 {
		return easygo.NewFailMsg("状态参数错误")
	}

	msg := fmt.Sprintf("审核名片语音作品:%d为%d", reqMsg.GetId64(), reqMsg.GetId32())
	queryBson := bson.M{"_id": reqMsg.GetId64()}
	updateBson := bson.M{"$set": bson.M{"Status": reqMsg.GetId32()}}
	one := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_MIX_VOICE_VIDEO, queryBson, updateBson, true)
	if one != nil {
		AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	}

	return easygo.EmptyMsg
}

//删除用户作品
func (l *cls4) RpcDelPlayerVoiceWork(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}
	for_game.UpdateAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_MIX_VOICE_VIDEO, bson.M{"_id": bson.M{"$in": idList}}, bson.M{"$set": bson.M{"Status": PlayerVoiceWork_Status_Del}})

	count := len(idList)
	var ids string
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idList[i])) + ","
		} else {
			ids += easygo.IntToString(int(idList[i]))
		}
	}
	msg := fmt.Sprintf("批量删除用户作品: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	return easygo.EmptyMsg
}

//上传用户作品
func (l *cls4) RpcUploadPlayerVoiceWork(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.PlayerMixVoiceVideo) easygo.IMessage {
	mixUrl := reqMsg.GetMixVoiceUrl()
	if mixUrl == "" {
		return easygo.NewFailMsg("作品url不能为空")
	}

	bgId := reqMsg.GetBgId()
	if bgId == 0 {
		return easygo.NewFailMsg("对应资源id不能为空")
	}

	if for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_VIDEO, bson.M{"_id": bgId}) == nil {
		return easygo.NewFailMsg("对应资源id不存在")
	}

	var player *share_message.PlayerBase
	if reqMsg.GetPlayerId() > 0 {
		player = for_game.GetPlayerById(reqMsg.GetPlayerId())
	} else {
		if reqMsg.GetPlayerAccount() == "" {
			return easygo.NewFailMsg("用户ID或用户柠檬号至少要有一个")
		}
		player = QueryPlayerbyAccount(reqMsg.GetPlayerAccount())
	}

	if player.GetPlayerId() == 0 {
		return easygo.NewFailMsg("用户不存在")
	}

	msg := fmt.Sprintf("修改用户名片语音作品:%d", reqMsg.GetId())
	if reqMsg.Id == nil {
		reqMsg = &share_message.PlayerMixVoiceVideo{
			Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_PLAYER_MIX_VOICE_VIDEO)),
			PlayerId:   easygo.NewInt64(player.GetPlayerId()),
			PlayerType: easygo.NewInt64(player.GetTypes()),
			Status:     easygo.NewInt32(1), //状态：0-未审核,1-已发布,2-使用中,3-已删除
			CreateTime: easygo.NewInt64(easygo.NowTimestamp()),
		}
		msg = fmt.Sprintf("添加用户名片语音作品:%d", reqMsg.GetId())
	}

	reqMsg.BgId = easygo.NewInt64(bgId)
	reqMsg.MixVoiceUrl = easygo.NewString(mixUrl)

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	one := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_MIX_VOICE_VIDEO, queryBson, updateBson, true)
	if one != nil {
		pmg := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
		if pmg != nil {
			pmg.SetMixId(reqMsg.GetId())
			pmg.SaveOneRedisDataToMongo("MixId", reqMsg.GetId())
		}
		if reqMsg.GetIsUse() {
			if pmg != nil {
				path := path.Join("backstage", "match", "picture")
				lis := QQbucket.GetObjectList(path, "", 100)
				if len(lis) > 0 {
					i := for_game.RandInt(0, len(lis))
					bgimg := "https://im-resource-1253887233.file.myqcloud.com/" + *lis[i].Title
					pmg.SetBgImageUrl(bgimg)
					pmg.SaveOneRedisDataToMongo("BgImageUrl", bgimg)
				}
			}
		}
		AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	}

	return easygo.EmptyMsg
}

//Id查询用户作品
func (l *cls4) RpcGetPlayerVoiceWorkUse(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.GetId64() == 0 {
		return easygo.NewFailMsg("ID参数错误")
	}

	queryBson := bson.M{"_id": reqMsg.GetId64()}
	result := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_MIX_VOICE_VIDEO, queryBson)
	one := &share_message.PlayerMixVoiceVideo{}
	if result != nil {
		for_game.StructToOtherStruct(result, one)
	}

	return one
}

//背景资源类型管理
func (l *cls4) RpcBgTagList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-_id"}
	if reqMsg.Status != nil && reqMsg.GetStatus() != 1000 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.ListType != nil {
		findBson["InterestType"] = reqMsg.GetListType()
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		findBson["Name"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_TAG, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.InterestTag
	for _, li := range lis {
		one := &share_message.InterestTag{}
		for_game.StructToOtherStruct(li, one)

		unReviewCount := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_VIDEO, bson.M{"Tags": one.GetId(), "Status": PlayerVoiceWork_Status_Untreated})
		one.UnReviewCount = easygo.NewInt64(unReviewCount)
		count := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_VIDEO, bson.M{"Tags": one.GetId()})
		one.Count = easygo.NewInt64(count)
		useM := []bson.M{
			{"$match": bson.M{"Tags": one.GetId()}},
			{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": "$UseCount"}}},
		}
		usecount := for_game.FindPipeAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_VIDEO, useM)
		one.UseCount = easygo.NewInt64(usecount)
		list = append(list, one)
	}
	msg := &brower_backstage.InterestTagResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//背景资源类型保存
func (l *cls4) RpcUpdateBgTag(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.InterestTag) easygo.IMessage {
	if reqMsg.InterestType == nil {
		return easygo.NewFailMsg("大类参数不能为空")
	}
	if reqMsg.GetInterestType() > 3 || reqMsg.GetInterestType() < 1 {
		return easygo.NewFailMsg("大类参数错误")
	}

	if reqMsg.Name == nil && reqMsg.GetName() == "" {
		return easygo.NewFailMsg("类别名字不能为空")
	}

	if reqMsg.Status == nil && reqMsg.GetStatus() == 0 {
		return easygo.NewFailMsg("状态不能为空")
	}

	msg := fmt.Sprintf("修改背景资源类型:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_BG_VOICE_TAG))
		msg = fmt.Sprintf("添加背景资源类型:%d", reqMsg.GetId())
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_TAG, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)

	return easygo.EmptyMsg
}

//背景资源管理列表
func (l *cls4) RpcBgVoiceVideoList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-_id"}

	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	//状态查询
	if reqMsg.Status != nil && reqMsg.GetStatus() != 1000 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	//用户类型查询
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["PlayerType"] = reqMsg.GetListType()
	}

	if reqMsg.Type != nil {
		findBson["Type"] = reqMsg.GetType()
	}

	if reqMsg.DownType != nil {
		findBson["Tags"] = reqMsg.GetDownType()
	}

	//用户查询
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE, bson.M{"$or": []bson.M{{"Account": reqMsg.GetKeyword()}, {"Phone": reqMsg.GetKeyword()}}})
		if one == nil {
			return easygo.NewFailMsg("用户不存在")
		}

		findBson["PlayerId"] = one.(bson.M)["_id"].(int64)
	}

	if reqMsg.Id != nil && reqMsg.GetId() > 0 {
		findBson["_id"] = reqMsg.GetId()
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_VIDEO, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.BgVoiceVideo
	for _, li := range lis {
		one := &share_message.BgVoiceVideo{}
		for_game.StructToOtherStruct(li, one)
		player := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE, bson.M{"_id": one.GetPlayerId()})
		if player != nil {
			one.PlayerAccount = easygo.NewString(player.(bson.M)["Account"])
		}
		list = append(list, one)
	}
	msg := &brower_backstage.BgVoiceVideoListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//背景资源审核
func (l *cls4) RpcRevieweBgVoiceVideo(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.GetId32() > 2 || reqMsg.GetId32() < 0 {
		return easygo.NewFailMsg("状态参数错误")
	}

	msg := fmt.Sprintf("审核名片背景资源:%d为%d", reqMsg.GetId64(), reqMsg.GetId32())
	queryBson := bson.M{"_id": reqMsg.GetId64()}
	updateBson := bson.M{"$set": bson.M{"Status": reqMsg.GetId32()}}
	one := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_VIDEO, queryBson, updateBson, true)
	if one != nil {
		AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	}

	return easygo.EmptyMsg
}

//背景资源保存
func (l *cls4) RpcUpdateBgVoiceVideo(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.BgVoiceVideo) easygo.IMessage {
	if reqMsg.Type == nil {
		return easygo.NewFailMsg("大类参数不能为空")
	}

	if reqMsg.GetType() > 3 || reqMsg.GetType() < 1 {
		return easygo.NewFailMsg("大类参数错误")
	}

	if reqMsg.Name == nil && reqMsg.GetName() == "" {
		return easygo.NewFailMsg("作品名不能为空")
	}

	if reqMsg.PlayerAccount == nil {
		return easygo.NewFailMsg("上传者不能为空")
	}

	if len(reqMsg.GetTags()) == 0 {
		return easygo.NewFailMsg("作品类型不能为空")
	}

	if reqMsg.Content == nil {
		return easygo.NewFailMsg("作品片段不能为空")
	}

	if reqMsg.ImageUrl == nil {
		return easygo.NewFailMsg("图片链接不能为空")
	}

	if reqMsg.MusicTime == nil {
		return easygo.NewFailMsg("音乐时长不能为空")
	}

	player := QueryPlayerbyAccount(reqMsg.GetPlayerAccount())
	if player == nil {
		return easygo.NewFailMsg("上传者不存在")
	}
	reqMsg.PlayerId = easygo.NewInt64(player.GetPlayerId())
	if reqMsg.GetStatus() == PlayerVoiceWork_Status_Untreated {
		reqMsg.Status = easygo.NewInt32(PlayerVoiceWork_Status_Release)
	}
	reqMsg.CreateTime = easygo.NewInt64(easygo.NowTimestamp())
	reqMsg.PlayerType = easygo.NewInt64(player.GetTypes())
	if reqMsg.Maker == nil {
		reqMsg.Maker = easygo.NewString(player.GetNickName())
	}

	msg := fmt.Sprintf("修改背景资源:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_BG_VOICE_VIDEO))
		reqMsg.UseCount = easygo.NewInt64(0)
		msg = fmt.Sprintf("添加背景资源:%d", reqMsg.GetId())
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_VIDEO, queryBson, updateBson, true)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)

	return easygo.EmptyMsg
}

//背景资源删除
func (l *cls4) RpcDelBgVoiceVideo(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}
	for_game.UpdateAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_BG_VOICE_VIDEO, bson.M{"_id": bson.M{"$in": idList}}, bson.M{"$set": bson.M{"Status": PlayerVoiceWork_Status_Del}})

	count := len(idList)
	var ids string
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idList[i])) + ","
		} else {
			ids += easygo.IntToString(int(idList[i]))
		}
	}
	msg := fmt.Sprintf("批量删除匹配文案: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	return easygo.EmptyMsg
}

//匹配文案列表
func (l *cls4) RpcMatchGuideList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-_id"}
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_MATCH_GUIDE, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.CommStrId
	for _, li := range lis {
		one := &share_message.CommStrId{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.MatchGuideListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//匹配文案更新
func (l *cls4) RpcUpdateMatchGuide(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	cont := reqMsg.GetIdsStr()
	if reqMsg.IdsStr == nil || len(cont) == 0 {
		return easygo.NewFailMsg("内容不能为空")
	}

	msg := fmt.Sprintf("更新匹配文案:%s为%s", reqMsg.GetNote(), cont[0])
	if reqMsg.Note != nil && reqMsg.GetNote() != "" {
		one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_MATCH_GUIDE, bson.M{"_id": reqMsg.GetNote()})
		if one == nil {
			return easygo.NewFailMsg("要修改的内容不存在")
		}
		for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_MATCH_GUIDE, bson.M{"_id": reqMsg.GetNote()})
		for_game.InsertAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_MATCH_GUIDE, &share_message.CommStrId{Id: easygo.NewString(cont[0])})
	} else {
		var ids []interface{}
		for _, i := range cont {
			one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_MATCH_GUIDE, bson.M{"_id": i})
			if one == nil {
				ids = append(ids, &share_message.CommStrId{Id: easygo.NewString(i)})
			}
		}

		for_game.InsertAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_MATCH_GUIDE, ids...)
		msg = fmt.Sprintf("新增匹配文案:%v", cont)
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	return easygo.EmptyMsg
}

//匹配文案删除
func (l *cls4) RpcDelMatchGuide(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIdsStr()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}
	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_MATCH_GUIDE, bson.M{"_id": bson.M{"$in": idList}})

	var ids string
	count := len(idList)
	for i := 0; i < count; i++ {
		if i < count {
			ids += idList[i] + ","

		} else {
			ids += idList[i]
		}
	}
	msg := fmt.Sprintf("批量删除匹配文案: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	return easygo.EmptyMsg
}

//SayHi文案列表
func (l *cls4) RpcSayHiList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-_id"}
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_SAY_HI, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.CommStrId
	for _, li := range lis {
		one := &share_message.CommStrId{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.MatchGuideListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//SayHi文案更新
func (l *cls4) RpcUpdateSayHi(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	cont := reqMsg.GetIdsStr()
	if reqMsg.IdsStr == nil || len(cont) == 0 {
		return easygo.NewFailMsg("内容不能为空")
	}

	msg := fmt.Sprintf("更新匹配文案:%s为%s", reqMsg.GetNote(), cont[0])
	if reqMsg.Note != nil && reqMsg.GetNote() != "" {
		one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_SAY_HI, bson.M{"_id": reqMsg.GetNote()})
		if one == nil {
			return easygo.NewFailMsg("要修改的内容不存在")
		}
		for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_SAY_HI, bson.M{"_id": reqMsg.GetNote()})
		for_game.InsertAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_SAY_HI, &share_message.CommStrId{Id: easygo.NewString(cont[0])})
	} else {
		var ids []interface{}
		for _, i := range cont {
			one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_SAY_HI, bson.M{"_id": i})
			if one == nil {
				ids = append(ids, &share_message.CommStrId{Id: easygo.NewString(i)})
			}
		}

		for_game.InsertAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_SAY_HI, ids...)
		msg = fmt.Sprintf("新增匹配文案:%v", cont)
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	return easygo.EmptyMsg
}

//SayHi文案删除
func (l *cls4) RpcDelSayHi(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIdsStr()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}
	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_SAY_HI, bson.M{"_id": bson.M{"$in": idList}})

	var ids string
	count := len(idList)
	for i := 0; i < count; i++ {
		if i < count {
			ids += idList[i] + ","

		} else {
			ids += idList[i]
		}
	}
	msg := fmt.Sprintf("批量删除SayHi文案: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	return easygo.EmptyMsg
}

//查询系统背景资源图
func (l *cls4) RpcSystemBgImageList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-_id"}
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_BG_IMAGE, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.SystemBgImage
	for _, li := range lis {
		one := &share_message.SystemBgImage{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.SystemBgImageListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//保存系统背景资源图
func (l *cls4) RpcSaveSystemBgImage(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	cont := reqMsg.GetIdsStr()
	if reqMsg.IdsStr == nil || len(cont) == 0 {
		return easygo.NewFailMsg("内容不能为空")
	}

	var ids []interface{}
	for _, i := range cont {
		one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_BG_IMAGE, bson.M{"_id": i})
		if one == nil {
			ids = append(ids, &share_message.SystemBgImage{Url: easygo.NewString(i)})
		}
	}

	for_game.InsertAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_BG_IMAGE, ids...)
	msg := fmt.Sprintf("新增背景资源:%v", cont)

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	return easygo.EmptyMsg
}

//删除系统背景资源图
func (l *cls4) RpcDelSystemBgImage(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIdsStr()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}
	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_BG_IMAGE, bson.M{"_id": bson.M{"$in": idList}})
	//TODO 删除存储桶资源
	var ids string
	count := len(idList)
	for i := 0; i < count; i++ {
		if i < count {
			ids += idList[i] + ","

		} else {
			ids += idList[i]
		}
	}
	msg := fmt.Sprintf("批量删除背景资源: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, msg)
	return easygo.EmptyMsg
}

//查询亲密度分值配置
func (l *cls4) RpcQueryIntimacyConfig(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	findBson := bson.M{}
	var sort []string
	sort = append(sort, "_id")
	lis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_INTIMACY_COINFIG, findBson, 0, 0, sort...)
	var list []*share_message.IntimacyConfig
	for _, li := range lis {
		one := &share_message.IntimacyConfig{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	return &brower_backstage.IntimacyConfigRes{
		List: list,
	}
}

//修改亲密度分值配置
func (l *cls4) RpcUpdateIntimacyConfig(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.IntimacyConfigRes) easygo.IMessage {
	lis := reqMsg.GetList()
	if len(lis) != 6 {
		return easygo.NewFailMsg("配置条数错误")
	}

	var data []interface{}
	for _, v := range lis {
		b1 := bson.M{"_id": v.GetLv()}
		b2 := v
		data = append(data, b1, b2)
	}
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_INTIMACY_COINFIG, data)
	for_game.InitConfigIntimacy()
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.LOVE_MANAGE, "修改亲密度分值配置")
	return easygo.EmptyMsg
}

//匹配埋点报表
func (l *cls4) RpcVCBuryingPointReport(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	var sort []string
	sort = append(sort, "-_id")
	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_VC_BURYING_POINT_REPORT, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.VCBuryingPointReport
	for _, li := range lis {
		one := &share_message.VCBuryingPointReport{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	return &brower_backstage.VCBuryingPointReportRes{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}
