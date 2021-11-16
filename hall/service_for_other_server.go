package hall

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/server_server"

	"github.com/astaxie/beego/logs"
)

//TODO 接收分发方法，函数方法名自己定，对端发送时传什么，这里就接收什么,需要返回值，则直接返回，但在接收的地方要强转
func (self *WebHttpForServer) RpcNotifyAddCoin(common *base.Common, reqMsg *server_server.NotifyAddCoinReq) easygo.IMessage {
	logs.Info("=======公共服务器发过来的请求 RpcNotifyAddCoin=====reqMsg=%v", reqMsg)
	if reqMsg.GetNotifyAddCoin() == "" {
		return nil
	}
	playerBCoin := make(map[int64]int64)
	err := json.Unmarshal([]byte(reqMsg.GetNotifyAddCoin()), &playerBCoin)
	easygo.PanicError(err)
	for pid, coin := range playerBCoin {
		if coin == 0 {
			continue
		}
		NotifyAddCoin(pid, -coin, "过期回收", for_game.COIN_TYPE_EXPIRED_OUT, nil, false)
		// 过期小助手通知
		content := fmt.Sprintf("您有%d个平台赠送硬币已过期", coin)
		NoticeAssistant(pid, 1, "温馨提示", content)

	}
	return nil
}
func (self *WebHttpForServer) RpcNoticeAssistant(common *base.Common, reqMsg *server_server.NotifyAddCoinReq) easygo.IMessage {
	logs.Info("=======公共服务器发过来的过期提示请求 RpcNoticeAssistant=====reqMsg=%v", reqMsg)
	if reqMsg.GetNotifyAddCoin() == "" {
		return nil
	}
	playerBCoin := make(map[int64]int64)
	err := json.Unmarshal([]byte(reqMsg.GetNotifyAddCoin()), &playerBCoin)
	easygo.PanicError(err)

	players := make([]int64, 0)
	for pid, coin := range playerBCoin {
		if coin == 0 {
			continue
		}
		//content := fmt.Sprintf("你有%d硬币即将过期，请尽快使用!", coin)
		content := fmt.Sprintf("您有%d个平台赠送硬币明日0点即将过期", coin)
		NoticeAssistant(pid, 1, "温馨提示", content)
		// 判断用户是否开启了推送
		p := for_game.GetRedisPlayerBase(pid)
		if p == nil {
			continue
		}
		if p.GetIsOpenCoinShop() {
			continue
		}
		players = append(players, pid)
	}
	if len(players) == 0 {
		return nil
	}
	ids := for_game.GetJGIds(players)
	m := for_game.PushMessage{
		Title: "温馨提示",
		// Content:     "您有平台赠送硬币明日0点即将过期",
		ContentType: for_game.JG_TYPE_HALL,
		JumpObject:  10,                     // 硬币页面
		ItemId:      for_game.PUSH_ITEM_101, //硬币过期提醒
	}

	sysp := PSysParameterMgr.GetSysParameter(for_game.PUSH_PARAMETER)
	if sysp != nil {
		pushSet := sysp.PushSet
		for _, ps := range pushSet {
			if ps.GetObjId() == m.ItemId && ps.GetIsPush() {
				m.Content = ps.GetObjContent()
			}
		}
	}
	for_game.JGSendMessage(ids, m, sysp)
	return nil
}

// 通知物品过期
func (self *WebHttpForServer) RpcNoticeProductExp(common *base.Common, reqMsg *server_server.NotifyAddCoinReq) easygo.IMessage {
	logs.Info("=======通知物品过期 RpcNoticeProductExp=====reqMsg=%v", reqMsg)
	if reqMsg.GetNotifyAddCoin() == "" {
		return nil
	}
	reqMap := make(map[int64]string)
	err := json.Unmarshal([]byte(reqMsg.GetNotifyAddCoin()), &reqMap)
	easygo.PanicError(err)

	players := make([]int64, 0)
	for pid, name := range reqMap {
		//content := fmt.Sprintf("你有%d硬币即将过期，请尽快使用!", coin)
		content := fmt.Sprintf("你以下道具马上就要到期了：%s。请尽快使用或续费。", name)
		NoticeAssistant(pid, 1, "温馨提示", content)
		// 判断用户是否开启了推送
		p := for_game.GetRedisPlayerBase(pid)
		if p == nil {
			continue
		}
		if p.GetIsOpenCoinShop() {
			continue
		}
		players = append(players, pid)
	}
	if len(players) == 0 {
		return nil
	}
	ids := for_game.GetJGIds(players)
	m := for_game.PushMessage{
		Title:       "温馨提示",
		Content:     "你有道具将要过期，请及时查看",
		ContentType: for_game.JG_TYPE_HALL,
		JumpObject:  11,                     // 物品页面
		ItemId:      for_game.PUSH_ITEM_102, // 道具过期提醒-您的%v即将过期,请及时查看
	}

	sysp := PSysParameterMgr.GetSysParameter(for_game.PUSH_PARAMETER)
	if sysp != nil {
		pushSet := sysp.PushSet
		for _, ps := range pushSet {
			if ps.GetObjId() == m.ItemId && ps.GetIsPush() {
				m.Content = ps.GetObjContent()
			}
		}
	}

	for_game.JGSendMessage(ids, m)
	return nil
}
