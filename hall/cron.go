package hall

import (
	"game_server/easygo"
	"time"
)

func StartCron() {
	defer easygo.RecoverAndLog()
	today0ClockTimestamp := easygo.GetToday0ClockTimestamp()
	now := time.Now().Unix()

	easygo.AfterFunc(time.Duration(today0ClockTimestamp+86400-now)*time.Second, HandleEvery000000)
	easygo.AfterFunc(time.Duration(today0ClockTimestamp+86400+300-now)*time.Second, HandleEvery000500)
	easygo.AfterFunc(time.Duration(today0ClockTimestamp+86400+600-now)*time.Second, HandleEvery001000)
}

// 每天00:00:00执行
func HandleEvery000000() {
	easygo.AfterFunc(time.Duration(easygo.GetToday0ClockTimestamp()+86400-time.Now().Unix())*time.Second, HandleEvery000000)
}

// 每天00:05:00执行
func HandleEvery000500() {

	easygo.AfterFunc(time.Duration(easygo.GetToday0ClockTimestamp()+86400+300-time.Now().Unix())*time.Second, HandleEvery000500)

	/*ClientEpMp.Range(func(key, value interface{}) bool {
		playerId := key.(PLAYER_ID)
		player := PlayerMgr.LoadPlayer(playerId)
		CheckNGiveVipCashBackWelfare(player.GetSite(), int64(playerId))
		return true
	})*/
}

// 每天00:10:00执行
func HandleEvery001000() {
	easygo.AfterFunc(time.Duration(easygo.GetToday0ClockTimestamp()+86400+600-time.Now().Unix())*time.Second, HandleEvery001000)
}
