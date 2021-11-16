package sport_common_dal

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"time"
	//   "time"
)

//添加评论的回复
func AddSportCommentReply(info *share_message.TableESportComment) (int32, string, int64) {

	table := ""
	table_c := ""
	table_p := ""
	if for_game.ESPORTMENU_REALTIME == info.GetMenuId() {
		table_c = for_game.TABLE_ESPORTS_NEWS
		table_p = for_game.TABLE_ESPORTS_COMMENT_NEWS
		table = for_game.TABLE_ESPORTS_COMMENT_NEWS_REPLY
	} else if for_game.ESPORTMENU_RECREATION == info.GetMenuId() {
		table_c = for_game.TABLE_ESPORTS_VIDEO
		table_p = for_game.TABLE_ESPORTS_COMMENT_VIDEO
		table = for_game.TABLE_ESPORTS_COMMENT_VIDEO_REPLY
	}
	code, msg, dataid := CreateTableESportCommentReply(table, info)
	if code == for_game.C_OPT_SUCCESS {
		//code, msg = UpdateFedAddition(table_p, "ReplyCount", info.GetCommentId(), 1)  //评资讯评论回复论回复+1
		//code, msg = UpdateFedAddition(table_c, "CommentCount", info.GetParentId(), 1) //+1
		UpdateFedAddition_xv(table_p, "ReplyCount", info.GetCommentId(), 1)  //评资讯评论回复论回复+1
		UpdateFedAddition_xv(table_c, "CommentCount", info.GetParentId(), 1) //+1
	}

	return code, msg, dataid
}

//添加评论的回复
func DeleteSportCommentReply(menuId int32, pid, pcid, cid int64) (int32, string) {

	table := ""
	table_c := ""
	table_p := ""
	if for_game.ESPORTMENU_REALTIME == menuId {
		table_c = for_game.TABLE_ESPORTS_NEWS
		table_p = for_game.TABLE_ESPORTS_COMMENT_NEWS
		table = for_game.TABLE_ESPORTS_COMMENT_NEWS_REPLY
	} else if for_game.ESPORTMENU_RECREATION == menuId {
		table_c = for_game.TABLE_ESPORTS_VIDEO
		table_p = for_game.TABLE_ESPORTS_COMMENT_VIDEO
		table = for_game.TABLE_ESPORTS_COMMENT_VIDEO_REPLY
	}
	b := DeleteTableESportCommentReply(table, cid)
	code := for_game.C_SYS_ERROR
	msg := ""
	if b {
		//code, msg = UpdateFedAddition(table_p, "ReplyCount", pcid, -1) //评资讯评论回复论回复+1
		//code, msg = UpdateFedAddition(table_c, "CommentCount", pid, -1)

		UpdateFedAddition_xv(table_p, "ReplyCount", pcid, -1) //评资讯评论回复论回复+1
		UpdateFedAddition_xv(table_c, "CommentCount", pid, -1)
	}

	return code, msg
}
func CreateTableESportCommentReply(table string, info *share_message.TableESportComment) (int32, string, int64) {
	id := NextId(table)
	dataid := id
	col, closeFun := GetC(table)
	defer closeFun()
	info.CreateTime = easygo.NewInt64(time.Now().Unix())
	info.Id = easygo.NewInt64(id)
	_, err := col.Upsert(bson.M{"_id": id}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误", 0
	}
	return for_game.C_OPT_SUCCESS, "創建成功", dataid
}
func DeleteTableESportCommentReply(table string, id int64) bool {
	col, closeFun := GetC(table)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": id})
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}
func UpdateTableESportCommentReply(table string, info *share_message.TableESportComment) (int32, string) {
	col, closeFun := GetC(table)
	defer closeFun()
	updatedata := bson.M{}
	if info.Content != nil {
		updatedata["Content"] = info.GetContent()
	}
	if info.PlayerId != nil {
		updatedata["PlayerId"] = info.GetPlayerId()
	}
	if info.PlayerNickName != nil {
		updatedata["PlayerNickName"] = info.GetPlayerNickName()
	}
	if info.ParentId != nil {
		updatedata["ParentId"] = info.GetParentId()
	}
	if info.CommentId != nil {
		updatedata["CommentId"] = info.GetCommentId()
	}
	if info.MenuId != nil {
		updatedata["MenuId"] = info.GetMenuId()
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
func GetTableESportCommentReplyById(table string, id int64) *share_message.TableESportComment {
	col, closeFun := GetC(table)
	defer closeFun()
	data := &share_message.TableESportComment{}
	err := col.Find(bson.M{"_id": id}).One(&data)
	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return nil
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return data
}

func GetTableESportCommentReplyList1(table string, offset, limit int, sort string, keyword string) ([]*share_message.TableESportComment, int) {
	var list []*share_message.TableESportComment
	col, closeFun := GetC(table)
	defer closeFun()
	queryBson := bson.M{}
	if keyword != "" {
		queryBson["Content"] = bson.M{"$regex": "^" + keyword + "+"}
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
		list = []*share_message.TableESportComment{}
	}
	return list, count
}
func GetTableESportCommentReplyList2(table string, cPage, pSize int, sort string, keyword string) ([]*share_message.TableESportComment, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESportComment
	col, closeFun := GetC(table)
	defer closeFun()
	queryBson := bson.M{}
	if keyword != "" {
		queryBson["Content"] = bson.M{"$regex": "^" + keyword + "+"}
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
func GetTableESportCommentList(table string, cPage, pSize int32, parentId, commentId int64, orderby string) ([]*share_message.TableESportComment, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESportComment
	col, closeFun := GetC(table)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["ParentId"] = parentId
	if commentId > 0 {
		queryBson["CommentId"] = commentId
	}
	queryBson["Status"] = 1
	query := col.Find(queryBson)
	count, err := query.Count()
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	err = query.Sort(orderby).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	return list, count
}
