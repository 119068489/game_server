package sport_common_dal

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	//   "time"
)

//初始化标签数据
func InitLabel() {
	/*
		ESPORTS_LABEL_WZRY  int32 = 10001 //王者荣耀
		ESPORTS_LABEL_DOTA2 int32 = 10002 //dota2
		ESPORTS_LABEL_LOL   int32 = 10003 //英雄联盟lol
		ESPORTS_LABEL_CSGO  int32 = 10004 //CSGO
		ESPORTS_LABEL_OTHER int32 = 20001 //其他
	*/
	AddTableESPortsLabel("推荐", 0, 0, 1)
	AddTableESPortsLabel(for_game.LabelToESportNameMap[for_game.ESPORTS_LABEL_WZRY], 0, int64(for_game.ESPORTS_LABEL_WZRY), 3)
	AddTableESPortsLabel(for_game.LabelToESportNameMap[for_game.ESPORTS_LABEL_DOTA2], 0, int64(for_game.ESPORTS_LABEL_DOTA2), 3)
	AddTableESPortsLabel(for_game.LabelToESportNameMap[for_game.ESPORTS_LABEL_LOL], 0, int64(for_game.ESPORTS_LABEL_LOL), 3)
	AddTableESPortsLabel(for_game.LabelToESportNameMap[for_game.ESPORTS_LABEL_CSGO], 0, int64(for_game.ESPORTS_LABEL_CSGO), 3)
	AddTableESPortsLabel("其他", 0, int64(for_game.ESPORTS_LABEL_OTHER), 3)

}

func AddTableESPortsLabel(title string, menuId int32, labelId int64, t int32) {

	if t == 1 {
		if GetTableESPortsLabelByType(t) {
			logs.Info("已存在【%s】标签", title)
			return
		}
	}
	if t == 3 {
		if GetTableESPortsLabelBylabelId(labelId) {
			logs.Info("已存在【%s】标签", title)
			return
		}
	}

	CreateTableESPortsLabel(&share_message.TableESPortsLabel{
		Title:     easygo.NewString(title),
		Status:    easygo.NewInt32(1),
		Weight:    easygo.NewInt32(10),
		MenuId:    easygo.NewInt32(menuId),
		LabelId:   easygo.NewInt64(labelId),
		LabelType: easygo.NewInt32(t),
	})
}
func CreateTableESPortsLabel(info *share_message.TableESPortsLabel) (int32, string) {

	id := for_game.NextId(for_game.TABLE_ESPORTS_LABEL)
	col, closeFun := GetC(for_game.TABLE_ESPORTS_LABEL)
	defer closeFun()

	info.Id = easygo.NewInt64(id)
	_, err := col.Upsert(bson.M{"_id": id}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	logs.Info("初始化【%s】标签", info.GetTitle())
	return for_game.C_OPT_SUCCESS, "創建成功"
}
func DeleteTableESPortsLabel(id int64) bool {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_LABEL)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": id})
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}

/*
func UpdateTableESPortsLabelWeight(id int64, weight int32, playerId int64, SysTypeId int32) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_LABEL)
	defer closeFun()

	updatedata := bson.M{}
	updatedata["Weight"] = weight
	cinfo, err := col.Upsert(bson.M{"_id": id, "PlayerId": playerId, "SysTypeId": SysTypeId}, bson.M{"$set": updatedata})
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
*/
func UpdateTableESPortsLabel(info *share_message.TableESPortsLabel) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_LABEL)
	defer closeFun()
	updatedata := bson.M{}
	if info.Title != nil {
		updatedata["Title"] = info.GetTitle()
	}
	if info.Status != nil {
		updatedata["Status"] = info.GetStatus()
	}
	if info.Weight != nil {
		updatedata["Weight"] = info.GetWeight()
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

func GetTableESPortsLabelByType(labelType int32) bool {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_LABEL)
	defer closeFun()
	n, err := col.Find(bson.M{"LabelType": labelType}).Count()
	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
	}
	return n > 0
}

func GetTableESPortsLabelBylabelId(labelId int64) bool {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_LABEL)
	defer closeFun()
	n, err := col.Find(bson.M{"LabelId": labelId}).Count()
	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
	}
	return n > 0
}

func GetTableESPortsLabelById(id int64) *share_message.TableESPortsLabel {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_LABEL)
	defer closeFun()
	data := &share_message.TableESPortsLabel{}
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

///一次性返回菜单的所有标签
func GetTableESPortsLabelList(menuId int32) []*share_message.TableESPortsLabel {
	var list []*share_message.TableESPortsLabel

	col, closeFun := GetC(for_game.TABLE_ESPORTS_LABEL)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["Status"] = 1

	prdata := []bson.M{
		bson.M{"LabelType": for_game.ESPORTS_LABEL_TYPE_1},
		bson.M{"LabelType": for_game.ESPORTS_LABEL_TYPE_3},
		bson.M{"LabelType": for_game.ESPORTS_LABEL_TYPE_2, "MenuId": menuId}}
	queryBson["$or"] = prdata

	if menuId == for_game.ESPORTMENU_GAME { //臨時處理，比賽菜單屏蔽 Csgo和dota2
		andqueryBson := []bson.M{
			bson.M{"LabelId": bson.M{"$ne": for_game.ESPORTS_LABEL_CSGO}},
			bson.M{"LabelId": bson.M{"$ne": for_game.ESPORTS_LABEL_DOTA2}},
			bson.M{"LabelId": bson.M{"$ne": for_game.ESPORTS_LABEL_OTHER}},
		}
		queryBson["$and"] = andqueryBson
	}

	query := col.Find(queryBson)

	err := query.Sort("+LabelType", "-Weight").All(&list)
	if err != nil {
		logs.Error(err)
		return nil
	}
	if list == nil {
		list = []*share_message.TableESPortsLabel{}
	}
	return list
}

//游戏标签
func GetTableESPortsGameLabelList() []*share_message.TableESPortsLabel {
	var list []*share_message.TableESPortsLabel

	col, closeFun := GetC(for_game.TABLE_ESPORTS_LABEL)
	defer closeFun()
	queryBson := bson.M{}

	queryBson["OpenFlag"] = 1

	query := col.Find(queryBson)

	err := query.Sort("+SortWeight").All(&list)
	if err != nil {
		logs.Error(err)
		return nil
	}
	if list == nil {
		list = []*share_message.TableESPortsLabel{}
	}
	return list
}
