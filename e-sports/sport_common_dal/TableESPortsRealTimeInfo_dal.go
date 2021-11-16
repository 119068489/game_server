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

func CreateTableESPortsRealTimeInfo(info *share_message.TableESPortsRealTimeInfo) (int32, string) {
	id := for_game.NextId(for_game.TABLE_ESPORTS_NEWS)
	col, closeFun := GetC(for_game.TABLE_ESPORTS_NEWS)
	defer closeFun()

	info.Id = easygo.NewInt64(id)
	_, err := col.Upsert(bson.M{"_id": id}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}
func DeleteTableESPortsRealTimeInfo(id int64) bool {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_NEWS)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": id})
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}
func UpdateTableESPortsRealTimeInfo(info *share_message.TableESPortsRealTimeInfo) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_NEWS)
	defer closeFun()
	updatedata := bson.M{}
	updatedata["UpdateTime"] = time.Now().Unix()
	if info.Status != nil {
		updatedata["Status"] = info.GetStatus()
	}
	if info.IssueTime != nil {
		updatedata["IssueTime"] = info.GetIssueTime()
	}
	if info.CoverBigImageUrl != nil {
		updatedata["CoverBigImageUrl"] = info.GetCoverBigImageUrl()
	}
	if info.CoverSmallImageUrl != nil {
		updatedata["CoverSmallImageUrl"] = info.GetCoverSmallImageUrl()
	}
	if info.Title != nil {
		updatedata["Title"] = info.GetTitle()
	}
	if info.Content != nil {
		updatedata["Content"] = info.GetContent()
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
func GetTableESPortsRealTimeInfoById(id int64) *share_message.TableESPortsRealTimeInfo {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_NEWS)
	defer closeFun()
	data := &share_message.TableESPortsRealTimeInfo{}
	err := col.Find(bson.M{"_id": id}).One(&data)
	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return nil
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	UpdateFedAdditionEx(col, "LookCount", id, 1) //访问数+1
	//logs.Info("GetTableESPortsRealTimeInfoById", *rerurnVal.Value)
	return data
}

func GetTableESPortsRealTimeInfoList1(offset, limit int, sort string, keyword string) ([]*share_message.TableESPortsRealTimeInfo, int) {
	var list []*share_message.TableESPortsRealTimeInfo
	col, closeFun := GetC(for_game.TABLE_ESPORTS_NEWS)
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
		list = []*share_message.TableESPortsRealTimeInfo{}
	}
	return list, count
}
func GetTableESPortsRealTimeInfoList2(cPage, pSize int32, sort string, gameTypeId int32) ([]*share_message.TableESPortsRealTimeInfo, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESPortsRealTimeInfo
	col, closeFun := GetC(for_game.TABLE_ESPORTS_NEWS)
	defer closeFun()
	queryBson := bson.M{}

	queryBson["LabelId"] = gameTypeId

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
func GetESPortsRealTimeItemList(cPage, pSize int32, sort string, typeId, lableId int64) ([]*share_message.TableESPortsRealTimeInfo, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESPortsRealTimeInfo
	col, closeFun := GetC(for_game.TABLE_ESPORTS_NEWS)
	defer closeFun()
	queryBson := bson.M{}
	if typeId == 3 {
		queryBson["AppLabelID"] = lableId
	} else if typeId == 2 {
		queryBson["LabelIds"] = bson.M{"$elemMatch": bson.M{"$eq": lableId}}
	}
	queryBson["Status"] = 1
	now := time.Now().Unix()

	prdata := []bson.M{bson.M{"EffectiveType": 1}, bson.M{"EffectiveType": 2, "BeginEffectiveTime": bson.M{"$lte": now}}}

	queryBson["$or"] = prdata //bson.M{"$or": prdata}
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
	if list == nil {
		list = []*share_message.TableESPortsRealTimeInfo{}
	}
	return list, count
}
