package backstage

import (
	"encoding/base64"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"regexp"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

//查询用户列表
func GetPlayerList(user *share_message.Manager, reqMsg *brower_backstage.GetPlayerListRequest) ([]*share_message.PlayerBase, int) {
	var list []*share_message.PlayerBase
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	// 关键词查询账号
	// queryBson := bson.M{"$or": []bson.M{bson.M{"Phone": bson.M{"$not": bson.RegEx{Pattern: "^100.*$"}}}, bson.M{"Phone": bson.M{"$not": bson.RegEx{Pattern: "^200.*$"}}}}}
	var queryBson bson.M
	if reqMsg.GetPlayerType() == 0 || reqMsg.GetPlayerType() == 1 {
		if reqMsg.GetListType() > 1 {
			return nil, 0
		}

		queryBson = bson.M{"Types": 1}
		if user.GetRole() != 0 { //运营测试
			queryBson = bson.M{"Types": bson.M{"$in": []int32{1, 6}}}
		}
		// queryBson["Label"] = bson.M{"$ne": nil}
	} else {
		if reqMsg.GetListType() == 1 {
			return nil, 0
		}
		queryBson = bson.M{"Types": bson.M{"$ne": 1}}
		if user.GetRole() != 0 { //运营测试
			queryBson = bson.M{"Types": bson.M{"$in": []int32{2, 3, 4, 5}}}
		}
	}

	switch reqMsg.GetType() {
	case 1:
		queryBson["Account"] = easygo.If(reqMsg.GetKeyword() != "", reqMsg.GetKeyword(), bson.M{"$ne": nil})
	case 2:
		queryBson["NickName"] = easygo.If(reqMsg.GetKeyword() != "", bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}, bson.M{"$ne": nil})
	case 3:
		queryBson["Phone"] = easygo.If(reqMsg.GetKeyword() != "", reqMsg.GetKeyword(), bson.M{"$ne": nil})
	case 4:
		queryBson["ApiUrl"] = easygo.If(reqMsg.GetKeyword() != "", reqMsg.GetKeyword(), bson.M{"$ne": nil})
	case 5:
		querybs := "^" + reqMsg.GetKeyword() + ".*$"
		queryBson["Brand"] = easygo.If(reqMsg.GetKeyword() != "", bson.M{"$regex": bson.RegEx{Pattern: querybs}}, bson.M{"$ne": nil})
	case 6:
		queryBson["Version"] = easygo.If(reqMsg.GetKeyword() != "", reqMsg.GetKeyword(), bson.M{"$ne": nil})
	default:
		queryBson["Account"] = easygo.If(reqMsg.GetKeyword() != "", reqMsg.GetKeyword(), bson.M{"$ne": nil})
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		queryBson["Types"] = reqMsg.GetListType()
	}

	if reqMsg.Sex != nil && reqMsg.GetSex() != 0 {
		queryBson["Sex"] = reqMsg.GetSex()
	}

	switch reqMsg.GetIsOnline() {
	case 1:
		queryBson["IsOnline"] = true
	case 2:
		queryBson["IsOnline"] = false
	}

	// 判断有日期才按日期查询
	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		switch reqMsg.GetTimeType() {
		case 1:
			queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
		case 2:
			queryBson["VerifiedTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
		}

	}

	if reqMsg.Channel != nil && len(reqMsg.GetChannel()) > 0 {
		queryBson["Channel"] = bson.M{"$in": reqMsg.GetChannel()}
	}

	if reqMsg.Label != nil && len(reqMsg.GetLabel()) > 0 {
		label := reqMsg.GetLabel()
		easygo.SortSliceInt32(label, true)
		queryBson["Label"] = bson.M{"$elemMatch": bson.M{"$in": label}}
	}

	if reqMsg.CustomTag != nil && len(reqMsg.GetCustomTag()) > 0 {
		customtag := reqMsg.GetCustomTag()
		easygo.SortSliceInt32(customtag, true)
		queryBson["CustomTag"] = bson.M{"$elemMatch": bson.M{"$in": customtag}}
	}

	if reqMsg.GrabTag != nil && len(reqMsg.GetGrabTag()) > 0 {
		queryBson["GrabTag"] = bson.M{"$in": reqMsg.GetGrabTag()}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 1000 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.RegType != nil && reqMsg.GetRegType() != 0 {
		pls := []int64{}
		lis := for_game.GetRegisterLoginLogByType(reqMsg.GetRegType())
		for _, li := range lis {
			pls = append(pls, li.GetPlayerId())
		}
		if len(pls) > 0 {
			queryBson["_id"] = bson.M{"$in": pls}
		}

	}

	if reqMsg.DeviceType != nil && reqMsg.GetDeviceType() != 0 {
		queryBson["DeviceType"] = reqMsg.GetDeviceType()
	}
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改用户
func EditPlayer(reqMsg *share_message.PlayerBase, isCreate bool) {
	pMgr := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
	if reqMsg.GetPhone() != "" {
		playerAccount := for_game.GetRedisAccountObj(reqMsg.GetPlayerId())
		if playerAccount == nil {
			panic("不存在改玩家账号信息" + easygo.AnytoA(reqMsg.GetPlayerId()))
		}
		if playerAccount.GetAccount() != reqMsg.GetPhone() {
			playerAccount.SetAccount(reqMsg.GetPhone()) //改账号表手机号
			pMgr.SetPhone(reqMsg.GetPhone())            //改玩家表手机号
			//修改许愿池account
			msg := &server_server.PlayerSI{
				PlayerId: easygo.NewInt64(reqMsg.GetPlayerId()),
				Account:  easygo.NewString(reqMsg.GetPhone()),
			}
			CallRpcSetWishAccount(msg)

			if reqMsg.Password != nil && reqMsg.GetPassword() != "" {
				playerAccount.SetPassword(reqMsg.GetPassword())
			}
		} else {
			if reqMsg.Password != nil && reqMsg.GetPassword() != "" {
				playerAccount.SetPassword(reqMsg.GetPassword())
			}
		}
		playerAccount.SaveToMongo()
	}

	if reqMsg.YoungPassWord != nil && reqMsg.GetYoungPassWord() == "" {
		pMgr.SetYoungPassWord("")
	}

	pMgr.SetHeadIcon(reqMsg.GetHeadIcon())
	pMgr.SetNickName(reqMsg.GetNickName())
	pMgr.SetEmail(reqMsg.GetEmail())
	pMgr.SetBackgroundImageURL(reqMsg.GetBackgroundImageURL())
	if len(reqMsg.GetPhoto()) > 0 {
		pMgr.SetPhoto(reqMsg.GetPhoto())
	}
	pMgr.SetSex(reqMsg.GetSex())
	pMgr.SetApiUrl(reqMsg.GetApiUrl())
	pMgr.SetSecretKey(reqMsg.GetSecretKey())
	pMgr.SetIsCheckChatLog(reqMsg.GetIsCheckChatLog())

	pMgr.SaveToMongo()
}

//修改用户
func EditPlayerLable(reqMsg *brower_backstage.QueryDataByIds) {
	for _, item := range reqMsg.GetIds64() {
		pMgr := for_game.GetRedisPlayerBase(item)
		pMgr.SetRedisLabelList(reqMsg.GetIds32())
		pMgr.SaveToMongo()
	}
}

//帐号或手机号查询用户
func QueryPlayerByAccountOrPhone(key string) *share_message.PlayerBase {
	player := &share_message.PlayerBase{}
	one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE, bson.M{"$or": []bson.M{{"Account": key}, {"Phone": key}}})
	if one == nil {
		return nil
	}
	for_game.StructToOtherStruct(one, player)
	return player
}

//帐号查询用户
func QueryPlayerbyAccount(account string) *share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	siteOne := &share_message.PlayerBase{}
	err := col.Find(bson.M{"Account": account}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//帐号查询用户
func QueryPlayerbyPhone(phone string) *share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	siteOne := &share_message.PlayerBase{}
	err := col.Find(bson.M{"Phone": phone}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//id查询账号
func QueryPlayerAccountbyId(id PLAYER_ID) *share_message.PlayerAccount {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_ACCOUNT)
	defer closeFun()

	siteOne := &share_message.PlayerAccount{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//批量查询用户数据
func QueryplayerlistByIds(ids []int64) []*share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	queryBson := bson.M{"_id": bson.M{"$in": ids}}

	var list []*share_message.PlayerBase
	query := col.Find(queryBson)
	errc := query.Sort("-_id").All(&list)
	easygo.PanicError(errc)

	return list
}

//批量查询用户数据
func QueryplayerlistByAccounts(accounts []string) []*share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	queryBson := bson.M{"Account": bson.M{"$in": accounts}}

	var list []*share_message.PlayerBase
	query := col.Find(queryBson)
	errc := query.Sort("-Account").All(&list)
	easygo.PanicError(errc)

	return list
}

//修改用户
func EditPlayerAccount(reqMsg *share_message.PlayerAccount) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_ACCOUNT)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetPlayerId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//Id查询用户
func QueryPlayerbyId(id PLAYER_ID) *share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	siteOne := &share_message.PlayerBase{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//冻结解冻用户
func UpPlayerStatus(adminid []int64, status int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": adminid}}, bson.M{"$set": bson.M{"Status": status, "Note": nil}})
	easygo.PanicError(err)
}

//查询用户投诉
func GetPlayerComplaint(reqMsg *brower_backstage.ListRequest) ([]*share_message.PlayerComplaint, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_COMPLAINT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.GetListType() != 7 {
		queryTypes := bson.M{"Types": bson.M{"$nin": COMPLAINT_OTHER_TYPE}}
		queryRespondentId := bson.M{"RespondentId": bson.M{"$nin": GetOperateIds()}}
		queryReason := bson.M{"Reason": bson.M{"$nin": COMPLAINT_OTHER_REASON}}
		queryBson["$and"] = []bson.M{queryTypes, queryReason, queryRespondentId}
	}

	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		queryBson["Type"] = reqMsg.GetListType()
	}

	if reqMsg.Type != nil && reqMsg.GetType() != 0 {
		queryBson["Types"] = reqMsg.GetType()
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		queryBson["Status"] = reqMsg.GetStatus()
	}
	//不查关键词
	if reqMsg.GetKeyword() != "" && reqMsg.Keyword != nil {
		switch reqMsg.GetDownType() {
		case 1:
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			queryBson["PlayerId"] = player.GetPlayerId()
		case 2:
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			queryBson["RespondentId"] = player.GetPlayerId()
		case 3: //投诉人昵称
			player := for_game.GetPlayerByNickName(reqMsg.GetKeyword())
			queryBson["PlayerId"] = player.GetPlayerId()
		}
	}

	if reqMsg.Id != nil {
		queryBson["DynamicId"] = reqMsg.GetId()
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.PlayerComplaint
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	for i, item := range list {
		if item.ReTime != nil && item.GetReTime() < 9999999999 {
			item.ReTime = easygo.NewInt64(item.GetReTime() * 1000)
		}

		if item.GetCreateTime() < 9999999999 {
			item.CreateTime = easygo.NewInt64(item.GetCreateTime() * 1000)
		}
		player := QueryPlayerbyId(item.GetPlayerId())
		list[i].PlayerAcount = easygo.NewString(player.GetAccount())
		switch item.GetType() {
		case 5:
			team := QueryTeambyId(item.GetRespondentId())
			list[i].RespondentAcount = easygo.NewString(team.GetTeamChat())
		default:
			player := QueryPlayerbyId(item.GetRespondentId())
			list[i].RespondentAcount = easygo.NewString(player.GetAccount())
		}
	}

	return list, count
}

//查询其他投诉
func GetPlayerComplaintOther(reqMsg *brower_backstage.ListRequest) ([]*share_message.PlayerComplaint, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_COMPLAINT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryTypes := bson.M{"Types": bson.M{"$in": COMPLAINT_OTHER_TYPE}}
	queryRespondentId := bson.M{"RespondentId": bson.M{"$in": GetOperateIds()}}
	queryReason := bson.M{"Reason": bson.M{"$in": COMPLAINT_OTHER_REASON}}
	queryBson := bson.M{"$or": []bson.M{queryTypes, queryReason, queryRespondentId}}

	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		queryBson["Type"] = reqMsg.GetListType()
	}

	if reqMsg.Type != nil && reqMsg.GetType() != 0 {
		queryBson["Types"] = reqMsg.GetType()
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		queryBson["Status"] = reqMsg.GetStatus()
	}
	//不查关键词
	if reqMsg.GetKeyword() != "" && reqMsg.Keyword != nil {
		switch reqMsg.GetDownType() {
		case 1:
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			queryBson["PlayerId"] = player.GetPlayerId()
		case 2:
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			queryBson["RespondentId"] = player.GetPlayerId()
		case 3:
			queryBson["Operator"] = reqMsg.GetKeyword()

		}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.PlayerComplaint
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	for i, item := range list {
		if item.ReTime != nil && item.GetReTime() < 9999999999 {
			item.ReTime = easygo.NewInt64(item.GetReTime() * 1000)
		}

		if item.GetCreateTime() < 9999999999 {
			item.CreateTime = easygo.NewInt64(item.GetCreateTime() * 1000)
		}
		player := QueryPlayerbyId(item.GetPlayerId())
		list[i].PlayerAcount = easygo.NewString(player.GetAccount())
		switch item.GetType() {
		case 5:
			team := QueryTeambyId(item.GetRespondentId())
			list[i].RespondentAcount = easygo.NewString(team.GetTeamChat())
		default:
			player := QueryPlayerbyId(item.GetRespondentId())
			list[i].RespondentAcount = easygo.NewString(player.GetAccount())
		}
	}

	return list, count
}

//回复用户投诉
func ReplyPlayerComplaint(reqMsg *share_message.PlayerComplaint) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_COMPLAINT)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//查询批量添加用户列表
func GetFriendPlayerList(reqMsg *brower_backstage.GetPlayerFriendListRequest) ([]*share_message.PlayerBase, int) {
	var list []*share_message.PlayerBase
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	// 关键词查询账号
	queryBson := bson.M{}
	// 判断有日期才按日期查询
	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}
	queryBsonAnd := []bson.M{{"_id": bson.M{"$ne": reqMsg.GetPlayerId()}}, {"NickName": bson.M{"$ne": ""}}}
	friend_base := for_game.GetFriendBase(reqMsg.GetPlayerId())

	if friend_base != nil {
		friends := friend_base.GetFriendIds()
		for _, friend_id := range friends {
			queryBsonAnd = append(queryBsonAnd, bson.M{"_id": bson.M{"$ne": friend_id}})
		}
	}

	if reqMsg.Gender != nil && reqMsg.GetGender() != 0 {
		queryBsonAnd = append(queryBsonAnd, bson.M{"Sex": reqMsg.GetGender()})
	}

	if reqMsg.Province != nil && reqMsg.GetProvince() != "" {
		queryBsonAnd = append(queryBsonAnd, bson.M{"Provice": reqMsg.GetProvince()})
	}

	if reqMsg.City != nil && reqMsg.GetCity() != "" {
		queryBsonAnd = append(queryBsonAnd, bson.M{"City": reqMsg.GetCity()})
	}

	if reqMsg.Label != nil && reqMsg.GetLabel() != 0 {
		ids := []int32{}
		ids = append(ids, reqMsg.GetLabel())
		queryBsonAnd = append(queryBsonAnd, bson.M{"Label": bson.M{"$in": ids}})
	}

	if reqMsg.CustomTag != nil && reqMsg.GetCustomTag() != 0 {
		ids := []int32{}
		ids = append(ids, reqMsg.GetCustomTag())
		queryBsonAnd = append(queryBsonAnd, bson.M{"CustomTag": bson.M{"$in": ids}})
	}

	if reqMsg.GrabTag != nil && reqMsg.GetGrabTag() != 0 {
		queryBsonAnd = append(queryBsonAnd, bson.M{"GrabTag": reqMsg.GetGrabTag()})
	}

	if reqMsg.Region != nil && reqMsg.GetRegion() != "" {
		queryBsonAnd = append(queryBsonAnd, bson.M{"Area": reqMsg.GetRegion()})
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		if reqMsg.GetType() == 1 {
			queryBsonAnd = append(queryBsonAnd, bson.M{"Account": reqMsg.GetKeyword()})
		} else {
			queryBsonAnd = append(queryBsonAnd, bson.M{"NickName": reqMsg.GetKeyword()})
		}
	}

	if reqMsg.GrabTag != nil && reqMsg.GetGrabTag() != 0 {
		queryBsonAnd = append(queryBsonAnd, bson.M{"GrabTag": reqMsg.GetGrabTag()})
	}

	if reqMsg.CustomTag != nil && reqMsg.GetCustomTag() != 0 {
		ids := []int32{}
		ids = append(ids, reqMsg.GetCustomTag())
		queryBsonAnd = append(queryBsonAnd, bson.M{"CustomTag": bson.M{"$in": ids}})
	}

	if reqMsg.PlayerType != nil && reqMsg.GetPlayerType() > 0 {
		queryBsonAnd = append(queryBsonAnd, bson.M{"Types": reqMsg.GetPlayerType()})
	}

	queryBson["$and"] = queryBsonAnd
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	if reqMsg.PlayerId != nil {
		friend_base := for_game.GetFriendBase(reqMsg.GetPlayerId())
		if friend_base != nil {
			friend_ids := friend_base.GetFriendIds()

			for _, friend_id := range friend_ids {
				for _, player_base := range list {
					if player_base.GetPlayerId() == friend_id {
						player_base.IsFriend = easygo.NewInt32(1)
						break
					}
				}
			}
		}
	}

	for _, player_base := range list {
		if player_base.IsFriend == nil {
			player_base.IsFriend = easygo.NewInt32(0)
		}
	}

	return list, count
}

//查询拉人进群用户列表
func GetTeamPlayerList(team *share_message.TeamData, reqMsg *brower_backstage.GetTeamPlayerListRequest) ([]*share_message.PlayerBase, int) {
	var list []*share_message.PlayerBase
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	members := []PLAYER_ID{}
	members = append(members, team.MemberList...)

	// 关键词查询账号
	queryBson := bson.M{}
	queryBson["NickName"] = bson.M{"$ne": ""}
	queryBson["_id"] = bson.M{"$nin": members}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		if reqMsg.GetType() == 1 {
			queryBson["Account"] = reqMsg.GetKeyword()
		} else {
			queryBson["NickName"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.PlayerType != nil && reqMsg.GetPlayerType() != 0 {
		queryBson["Types"] = reqMsg.GetPlayerType()
	}

	if reqMsg.Label != nil && reqMsg.GetLabel() != 0 {
		ids := []int32{}
		ids = append(ids, reqMsg.GetLabel())
		queryBson["Label"] = bson.M{"$in": ids}
	}

	if reqMsg.CustomTag != nil && reqMsg.GetCustomTag() != 0 {
		ids := []int32{}
		ids = append(ids, reqMsg.GetCustomTag())
		queryBson["CustomTag"] = bson.M{"$in": ids}
	}

	if reqMsg.GrabTag != nil && reqMsg.GetGrabTag() != 0 {
		queryBson["GrabTag"] = reqMsg.GetGrabTag()
	}

	if reqMsg.Channel != nil && reqMsg.GetChannel() != "" {
		queryBson["Channel"] = reqMsg.GetChannel()
	}

	// 判断有日期才按日期查询
	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改玩家自定义标签
func EditPlayerCustomTag(reqMsg *brower_backstage.QueryDataByIds) {
	player := reqMsg.GetIds64()[0]
	pMgr := for_game.GetRedisPlayerBase(player)
	v := reqMsg.GetIds32()
	if len(v) > 0 {
		easygo.SortSliceInt32(v, true)
		pMgr.SetRedisCustomTag(v)
		pMgr.SaveOneRedisDataToMongo("CustomTag", v)
	}
}

//写玩家冻结日志
func AddPlayerFreezeLogs(players []int64, note string) {
	pls := GetPlayerBaseByIds(players)
	var list []interface{}
	for _, p := range pls {
		lis := &share_message.PlayerFreezeLog{
			Account:    p.Account,
			CreateTime: easygo.NewInt64(util.GetMilliTime()),
			Status:     easygo.NewInt32(2),
			Note:       easygo.NewString(note),
		}
		list = append(list, lis)
	}

	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PLAYER_FREEZE_LOG)
	defer closeFun()
	bulk := col.Bulk()
	bulk.Insert(list...)
	_, err := bulk.Run()
	if err != nil {
		easygo.PanicError(err)
	}
}

//查询用户列表
func GetPlayerFreezeLogsList(cur, size int32) ([]*share_message.PlayerFreezeLog, int) {
	var list []*share_message.PlayerFreezeLog
	pageSize := int(size)
	curPage := easygo.If(int(cur) > 1, int(cur)-1, 0).(int)

	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PLAYER_FREEZE_LOG)
	defer closeFun()

	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//获取运营账号的id
func GetOperateIds() []int64 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	var pli []int64
	player := &share_message.PlayerBase{}
	iter := col.Find(bson.M{"Types": bson.M{"$ne": 1}}).Iter()
	for iter.Next(&player) {
		pli = append(pli, player.GetPlayerId())
	}

	if err := iter.Close(); err != nil {
		easygo.PanicError(err)
	}

	return pli
}

//注销账号记录列表
func GetPlayerCancleAccountList(reqMsg *brower_backstage.ListRequest) ([]*share_message.PlayerCancleAccount, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CANCEL_ACCOUNT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}

	if reqMsg.Type != nil && reqMsg.GetType() != 0 && reqMsg.GetBeginTimestamp() != 0 {
		switch reqMsg.GetType() {
		case 1: //创建时间
			queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
		case 2: //完成时间
			queryBson["FinishTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
		}
	}
	//1000不查询状态
	if reqMsg.Status != nil && reqMsg.GetStatus() != 1000 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	//不查关键词
	if reqMsg.GetKeyword() != "" && reqMsg.Keyword != nil {
		switch reqMsg.GetDownType() {
		case 1: //账号查询
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			queryBson["PlayerId"] = player.GetPlayerId()
		case 2: //昵称查询
			player := QuryPlayerByNickname(reqMsg.GetKeyword())
			queryBson["PlayerId"] = player.GetPlayerId()
		}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.PlayerCancleAccount
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改注销账号
func EditPlayerCancleAccount(reqMsg *share_message.PlayerCancleAccount) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CANCEL_ACCOUNT)
	defer closeFun()
	qurybson := bson.M{}
	switch reqMsg.GetStatus() {
	case 0: //审核
		qurybson = bson.M{"$set": bson.M{"Note": reqMsg.GetNote()}}
	case 2: //拒绝
		qurybson = bson.M{"$set": bson.M{"Note": reqMsg.GetNote(), "Status": reqMsg.GetStatus(), "FinishTime": util.GetMilliTime()}}
		playerMgr := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
		playerMgr.SetStatus(for_game.ACCOUNT_NORMAL)
		playerMgr.SaveOneRedisDataToMongo("Status", for_game.ACCOUNT_NORMAL)
		// 发送注销失败短信
		var areaCode string
		if playerMgr.GetAreaCode() == "" {
			areaCode = "+86"
			// 用腾讯运营商
			easygo.Spawn(for_game.NewSMSInst(for_game.SMS_BUSINESS_TC).SendMessageCodeEx, fmt.Sprintf("%s%s", areaCode, playerMgr.GetPhone()), playerMgr.GetPhone(), false, false)
		}
	}
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, qurybson)
	easygo.PanicError(err)
}

//查询个人聊天记录
func GetPersonalChatLog(reqMsg *brower_backstage.ChatLogRequest) ([]*share_message.PersonalChatLog, int) {
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{"Type": bson.M{"$in": []int32{1, 2, 3}}}
	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		queryBson["Time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.Types != nil && reqMsg.GetTypes() != 1000 {
		queryBson["Type"] = reqMsg.GetTypes()
	}

	//不查关键词
	if reqMsg.GetKeyword1() != "" && reqMsg.Keyword1 != nil {
		player := QueryPlayerbyAccount(reqMsg.GetKeyword1())
		if player != nil {
			queryBson["Talker"] = player.GetPlayerId()
		}
	}

	if reqMsg.GetKeyword2() != "" && reqMsg.Keyword2 != nil {
		player := QueryPlayerbyAccount(reqMsg.GetKeyword2())
		if player != nil {
			queryBson["TargetId"] = player.GetPlayerId()
		}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.PersonalChatLog
	errc := query.Sort("-Time").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//查询个人聊天记录列表
func GetPersonalChatLogList(reqMsg *brower_backstage.ChatLogRequest) ([]*share_message.PersonalChatLog, int) {
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()

	talker := for_game.GetPlayerIdForAccount(reqMsg.GetKeyword1())
	targetId := for_game.GetPlayerIdForAccount(reqMsg.GetKeyword2())
	queryBson := bson.M{} //查询今天以前的所有数据

	var ids []int64
	if talker > 0 {
		ids = append(ids, talker)
		queryBson["Talker"] = talker
	}
	if targetId > 0 {
		ids = append(ids, targetId)
		queryBson["TargetId"] = targetId
	}

	//聚合查询
	m := []bson.M{
		{"$match": queryBson},
		{"$group": bson.M{"_id": "$Talker", "PlayerId": bson.M{"$addToSet": "$TargetId"}}},
		{"$unwind": "$PlayerId"},
	}

	query := col.Pipe(m)
	var list []*for_game.PlayerGroup
	err := query.All(&list)
	easygo.PanicError(err)
	if talker == 0 || targetId == 0 {
		for _, li := range list {
			if talker == 0 {
				ids = append(ids, li.Id)
			} else {
				ids = append(ids, li.PlayerId)
			}
		}
	}

	if targetId == 0 {
		for _, li := range list {
			if talker == 0 {
				ids = append(ids, li.Id)
			}
		}
	}

	players := QueryplayerlistByIds(ids)
	playerMap := make(map[int64]*share_message.PlayerBase)
	for _, p := range players {
		playerMap[p.GetPlayerId()] = p
	}

	var returnList = []*share_message.PersonalChatLog{}

	for _, item := range list {
		var one = &share_message.PersonalChatLog{}
		if talker != 0 && targetId != 0 {
			one.TalkerAccount = easygo.NewString(playerMap[talker].GetAccount())
			one.TalkerNickName = easygo.NewString(playerMap[talker].GetNickName())
			one.Talker = easygo.NewInt64(playerMap[talker].GetPlayerId())
			one.TargetAccount = easygo.NewString(playerMap[targetId].GetAccount())
			one.TargetNickName = easygo.NewString(playerMap[targetId].GetNickName())
			one.TargetId = easygo.NewInt64(playerMap[targetId].GetPlayerId())
		} else {
			one.TalkerAccount = easygo.NewString(playerMap[item.Id].GetAccount())
			one.TalkerNickName = easygo.NewString(playerMap[item.Id].GetNickName())
			one.Talker = easygo.NewInt64(playerMap[item.Id].GetPlayerId())
			one.TargetAccount = easygo.NewString(playerMap[item.PlayerId].GetAccount())
			one.TargetNickName = easygo.NewString(playerMap[item.PlayerId].GetNickName())
			one.TargetId = easygo.NewInt64(playerMap[item.PlayerId].GetPlayerId())
		}
		returnList = append(returnList, one)
	}

	return returnList, len(returnList)
}

//查询群聊天记录
func GetTeamChatLog(reqMsg *brower_backstage.ChatLogRequest) ([]*share_message.TeamChatLog, int) {
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_TEAM_CHAT_LOG)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{"Type": bson.M{"$in": []int32{1, 2, 3}}} //只查询文字，语音，图片
	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		queryBson["Time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.Types != nil && reqMsg.GetTypes() != 1000 {
		queryBson["Type"] = reqMsg.GetTypes()
	}

	//不查关键词
	if reqMsg.GetKeyword1() != "" && reqMsg.Keyword1 != nil {
		team := QueryTeambyNick(reqMsg.GetKeyword1())
		if team != nil {
			queryBson["TeamId"] = team.GetId()
		}
	}

	if reqMsg.GetKeyword2() != "" && reqMsg.Keyword2 != nil {
		player := QueryPlayerbyAccount(reqMsg.GetKeyword2())
		if player != nil {
			queryBson["Talker"] = player.GetPlayerId()
		}
	}

	if reqMsg.Keyword3 != nil && reqMsg.GetKeyword3() != "" {
		content := base64.StdEncoding.EncodeToString([]byte(reqMsg.GetKeyword3()))
		queryBson["Content"] = bson.M{"$regex": bson.RegEx{Pattern: content, Options: "im"}}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.TeamChatLog
	errc := query.Sort("-Time").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//查询当前在线的玩家数量
func QueryPlayerOnline() int64 {
	msg := ChooseOneHall(0, "RpcGetPlayerOnline", easygo.EmptyMsg)
	if p, ok := msg.(*server_server.PlayerSI); ok {
		return p.GetCount()
	}
	return 0
}

//分组查询个人聊天对象
func QueryPlayersByTalker(id PLAYER_ID) []*share_message.PlayerBase {
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()

	queryBson := bson.M{"Talker": id}
	//分组查询
	m := []bson.M{
		{"$match": queryBson},
		{"$group": bson.M{"_id": "$TargetId"}},
		{"$sort": bson.M{"_id": -1}},
	}

	query := col.Pipe(m)
	var list []*share_message.PlayerBase
	errc := query.All(&list)
	easygo.PanicError(errc)

	return list
}

//查询2个对象之间的聊天记录
func QueryPersonalChatLogByObj(reqMsg *brower_backstage.ChatLogRequest) ([]*share_message.PersonalChatLog, int, PLAYER_IDS) {
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{"Type": bson.M{"$in": []int32{1, 2, 3}}}
	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		queryBson["Time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.Types != nil && reqMsg.GetTypes() != 1000 {
		queryBson["Type"] = reqMsg.GetTypes()
	}

	if reqMsg.Keyword3 != nil && reqMsg.GetKeyword3() != "" {
		content := base64.StdEncoding.EncodeToString([]byte(reqMsg.GetKeyword3()))
		queryBson["Content"] = bson.M{"$regex": bson.RegEx{Pattern: content, Options: "im"}}
	}

	var ids PLAYER_IDS
	if reqMsg.Keyword1 != nil {
		player := QueryPlayerbyAccount(reqMsg.GetKeyword1())
		if player != nil {
			ids = append(ids, player.GetPlayerId())
		}
	}

	if reqMsg.Keyword2 != nil {
		player2 := QueryPlayerbyAccount(reqMsg.GetKeyword2())
		if player2 != nil {
			ids = append(ids, player2.GetPlayerId())
		}
	}

	if ids != nil {
		queryBson["Talker"] = bson.M{"$in": ids}
		queryBson["TargetId"] = bson.M{"$in": ids}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.PersonalChatLog
	errc := query.Sort("_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count, ids
}

//检查白名单
func CheckUserWhitelist(ids PLAYER_IDS) []*share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	queryBson := bson.M{"_id": bson.M{"$in": ids}, "IsCheckChatLog": true}

	var list []*share_message.PlayerBase
	query := col.Find(queryBson)
	errc := query.All(&list)
	easygo.PanicError(errc)

	return list
}

//生成测试用户
func TimingCreatePlayer() {
	brands := []string{"HUAWEI", "iPhone", "Xiaomi", "Meizu", "OPPO", "vivo", "HONOR", "samsung", "Redmi"}
	lis := for_game.QueryOperationChannleList()
	jsonFile := easygo.YamlCfg.GetValueAsString("CLIENT_VERSION_DATA")
	versiondata := for_game.FindVersion(jsonFile)
	for _, li := range lis {
		RandAddPlayer(int(li.GetAddCount()), li, versiondata, brands, "1111")
	}
}

//随机创建测试账号
func RandAddPlayer(addCount int, li *share_message.OperationChannel, versiondata *brower_backstage.VersionData, brands []string, pass string) {
	for i := 1; i <= addCount; i++ {
		sAccount := for_game.MakePhone()
		ip := for_game.GetRandPlayerIp()
		data := &share_message.CreateAccountData{
			Phone:    easygo.NewString(sAccount),
			PassWord: easygo.NewString(pass),
			Ip:       easygo.NewString(ip),
			Types:    easygo.NewInt32(6),
		}
		b, a := for_game.CreateAccount(data)
		if b {
			devType := util.RandIntn(2) + 1
			//检查是否苹果渠道
			channelType := []byte(li.GetName())
			pat := `苹果`
			reg1 := regexp.MustCompile(pat)
			channelReturn := reg1.Find(channelType)
			version := versiondata.Android.GetVersion()
			if channelReturn != nil {
				devType = 1
				version = versiondata.Ios.GetVersion()
			}
			lastLoginTime := util.GetMilliTime()
			onlineTime := int64(for_game.RandInt(60000, 600000))
			lastLoginOutTime := lastLoginTime + onlineTime

			//检查是否苹果渠道结束
			sexType := util.RandIntn(2) + 1
			player := for_game.GetRedisPlayerBase(a)
			player.SetNickName(for_game.GetRandNickName())
			player.SetDeviceType(int32(devType))
			player.SetCreateTime()
			player.SetRedisLabelList(for_game.GetRandLable())
			player.SetSex(int32(sexType))
			player.SetHeadIcon(for_game.GetRandRealHeadIcon(sexType))
			player.SetChannel(li.GetChannelNo())
			player.SetVersion(version)
			player.SetBrand(brands[util.RandIntn(9)])
			player.SetLastOnLineTime(lastLoginTime)
			player.SetLastLogOutTime(lastLoginOutTime)
			player.SetTodayOnlineTime(onlineTime)
			player.SetOnlineTime(onlineTime)
			player.AddLoginTimes()
			player.UpdateLastLoginIP(ip)
			player.SetSignature(for_game.GetRandSignature())
			player.SaveToMongo()

			for_game.AddStatisticsInfo(for_game.LOGINREGISTER_PHONEREGISTER, player.GetPlayerId(), 0, player.GetTypes())
			for_game.AddStatisticsInfo(for_game.LOGINREGISTER_MESSAGELOGIN, player.GetPlayerId(), 0, player.GetTypes())
			if li.ChannelNo != nil && li.GetChannelNo() != "" {
				for_game.SetRedisOperationChannelReportFildVal(easygo.Get0ClockTimestamp(lastLoginTime), int64(devType), li.GetChannelNo(), "ActDevCount")
				for_game.SetRedisOperationChannelReportFildVal(easygo.Get0ClockTimestamp(lastLoginTime), 1, li.GetChannelNo(), "ValidRegCount")
				for_game.SetRedisOperationChannelReportFildVal(easygo.Get0ClockTimestamp(lastLoginTime), 1, li.GetChannelNo(), "LoginCount")
			}

			for_game.AddOnlineTimeLog(player.GetPlayerId())
		}
	}
}
