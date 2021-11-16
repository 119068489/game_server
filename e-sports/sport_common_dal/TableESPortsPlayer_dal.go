package sport_common_dal

import (
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	//   "time"
)

func CreateTableESPortsPlayer(info *share_message.TableESPortsPlayer) (int32, string) {
	//id := for_game.NextId(for_game.TABLE_ESPORTS_PLAYER)
	col, closeFun := GetC(for_game.TABLE_ESPORTS_PLAYER)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": info.GetId()}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}
func DeleteTableESPortsPlayer(id int64) bool {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_PLAYER)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": id})
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}
func UpdateTableESPortsPlayer(info *share_message.TableESPortsPlayer) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_PLAYER)
	defer closeFun()
	updatedata := bson.M{}

	if info.Status != nil {
		updatedata["Status"] = info.GetStatus()
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

//根据PlayerId获取用户信息
func GetTableESPortsPlayerById(id int64) *share_message.TableESPortsPlayer {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_PLAYER)
	defer closeFun()
	data := &share_message.TableESPortsPlayer{}
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
