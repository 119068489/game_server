package sport_common_dal

import (
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	//   "time"
)

func CreateTableESPortsBpsClick(info *share_message.TableESPortsBpsClick) (int32, string) {

	col, closeFun := GetC(for_game.TABLE_ESPORTS_BPS_CLICK_LOG)
	defer closeFun()
	err := col.Insert(bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}

func CreateTableESPortsBpsDuration(info *share_message.TableESPortsBpsDuration) (int32, string) {

	col, closeFun := GetC(for_game.TABLE_ESPORTS_BPS_DURATION_LOG)
	defer closeFun()
	err := col.Insert(bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}
