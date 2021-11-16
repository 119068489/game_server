package sport_lottery_lol

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"reflect"
	"time"
)

//===================================================================

type ServiceForESports struct {
	Service reflect.Value
}

func (self *ServiceForESports) RpcESportLottery(common *base.Common, reqMsg *client_hall.LotteryRequest) easygo.IMessage {
	uniqueGameId := reqMsg.GetUniqueGameId()

	startStr := fmt.Sprintf("====== RpcESportLottery LOL定时器开奖发送  开始==========reqMsg:%v", reqMsg)
	for_game.WriteFile("lol_sport_lottery.log", startStr)

	//即时回调过来的按照设置的推迟时间开奖
	if LOLSysMsgTimeMgr.GetTimerById(uniqueGameId) != nil {
		LOLSysMsgTimeMgr.DelTimerList(uniqueGameId)
	}

	triggerTime := time.Duration(delayLotteryTime) * time.Second
	timer := easygo.AfterFunc(triggerTime, func() {
		DealESportLOLLottery(uniqueGameId)
	})
	LOLSysMsgTimeMgr.AddTimerList(uniqueGameId, timer)

	endStr := fmt.Sprintf("====== RpcESportLottery  LOL定时器开奖发送  结束==========")
	for_game.WriteFile("lol_sport_lottery.log", endStr)
	return nil
}
