package hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_server"
	"github.com/astaxie/beego/logs"
)

//关注某人-电竞入口
func (self *cls1) RpcESportAttentioPlayer(ep IGameClientEndpoint, who *Player, reqMsg *client_server.AttenInfo, common ...*base.Common) easygo.IMessage {
	logs.Info("==================大厅关注某人RpcAttentioPlayer==================", reqMsg)
	pid := reqMsg.GetPlayerId() // 被关注人的id
	t := reqMsg.GetType()
	player := for_game.GetRedisPlayerBase(pid)
	if t == for_game.DYNAMIC_OPERATE { // 注销了的账号不给关注.
		//注销的账号不给关注
		if player == nil || player.GetStatus() == for_game.ACCOUNT_CANCELED {
			return easygo.NewFailMsg("该账号异常")
		}
	}
	// 1-关注,2-取消关注柱
	if t != for_game.DYNAMIC_OPERATE && t != for_game.DYNAMIC_DELOPERATE {
		logs.Error("大厅关注操作,操作类型期待的是1或者2,传过来的值为: ", t)
		return easygo.NewFailMsg("操作类型有误")
	}
	if pid == 0 {
		logs.Error("关注操作,前端传的playerId为空")
		return easygo.NewFailMsg("关注操作失败")
	}

	if pid == who.GetPlayerId() {
		return easygo.NewFailMsg("不能关注自己")
	}
	if t == for_game.DYNAMIC_OPERATE && util.Int64InSlice(pid, who.GetAttention()) {
		return easygo.NewFailMsg("已关注该玩家")
	}
	b := for_game.OperateRedisDynamicAttentionEx(t, who.GetPlayerId(), pid, 1)
	if !b {
		return easygo.NewFailMsg("关注操作失败")
	}

	if t == for_game.DYNAMIC_OPERATE { // 关注
		who.AddAttention(pid)
		player := for_game.GetRedisPlayerBase(pid)
		if player == nil {
			panic("玩家怎么会为空")
		}
		player.AddFans(who.GetPlayerId())
	} else { //取消关注
		who.DelAttention(pid)
		player.DelFans(who.GetPlayerId())
	}
	redMsgCallClient(ep, pid) // 大厅通知社交广场,显示红点.
	return reqMsg
}
