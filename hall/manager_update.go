package hall

import (
	"fmt"
	"game_server/for_game"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo/bson"

	"game_server/easygo"
)

type UpdateManager struct {
	Mutex easygo.RLock
}

func NewUpdateManager() *UpdateManager {
	p := &UpdateManager{}
	p.Init()
	return p
}

func (self *UpdateManager) Init() {
	easygo.AfterFunc(time.Second*10, self.Update) //每10分钟检测一次红包过期 600
}
func (self *UpdateManager) Update() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	//指定大厅update
	if for_game.GetCurrentSaveServerSid(PServerInfo.GetSid(), for_game.REDIS_SAVE_HALL_SID) == PServerInfo.GetSid() {
		self.UpdateRedPackets()
		self.UpdateTransferMoney()
		//self.UpdateBCoinExpiration()
	}
	easygo.AfterFunc(time.Second*600, self.Update) //每10分钟检测一次红包过期 600
}

//红包定时器
func (self *UpdateManager) UpdateRedPackets() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var packets []*share_message.RedPacket
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_RED_PACKET)
	defer closeFun()
	t := time.Now().Unix()
	err := col.Find(bson.M{"State": 1, "CreateTime": bson.M{"$lt": t - for_game.DAY_SECOND}}).All(&packets)
	easygo.PanicError(err)
	if len(packets) > 0 {
		fun := func(p *share_message.RedPacket) {
			redPacket := for_game.GetRedisRedPacket(p.GetId()) //从Redis中取出红包
			//redPacket.SetOverTime() //红包标志已超时
			//退钱给发红包的人
			base := GetPlayerObj(redPacket.GetSender())
			if base != nil {
				info := map[string]interface{}{
					"RedType": redPacket.GetType(),
				}
				orderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, for_game.GOLD_TYPE_REDPACKET_OVERTIME)
				reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_REDPACKET_OVERTIME, redPacket.GetTargetId(), info)
				msg := &share_message.GoldExtendLog{
					OrderId:     easygo.NewString(orderId),
					RedPacketId: easygo.NewInt64(redPacket.GetId()),
					PayType:     easygo.NewInt32(redPacket.GetPayWay()),
					Title:       easygo.NewString(reason),
					Gold:        easygo.NewInt64(base.GetGold() + redPacket.GetCurMoney()),
				}
				NotifyAddGold(base.GetPlayerId(), redPacket.GetCurMoney(), reason, for_game.GOLD_TYPE_REDPACKET_OVERTIME, msg)
				//红包标识为超时状态
				title := "红包退款到账通知"
				text := fmt.Sprintf("收到一笔红包退款，退款金额￥%.2f。通过“我的”→“零钱”→“账单”可查看详情", float64(redPacket.GetCurMoney())/100)
				NoticeAssistant(redPacket.GetSender(), 1, title, text)
				redPacket.SetState(for_game.PACKET_MONEY_TIMEOUT)
				redPacket.SaveToMongo()
			}
		}
		for _, p := range packets {
			easygo.PCall(fun, p)
		}
	}
}

//转账定时器
func (self *UpdateManager) UpdateTransferMoney() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var packets []*share_message.TransferMoney
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TRANSFER_MONEY)
	defer closeFun()
	t := time.Now().Unix()
	err := col.Find(bson.M{"State": 1, "CreateTime": bson.M{"$lt": t - for_game.DAY_SECOND}}).All(&packets)
	easygo.PanicError(err)
	if len(packets) > 0 {
		fc := func(p *share_message.TransferMoney) {
			transferMoney := for_game.GetRedisTransferMoneyObj(p.GetId())
			//退钱给转账的人
			base := GetPlayerObj(transferMoney.GetSender())
			if base != nil {
				target := GetPlayerObj(transferMoney.GetTargetId())
				//orderId, _ := for_game.PlaceOrder(transferMoney.GetSender(), transferMoney.GetGold(), for_game.GOLD_TYPE_TRANSFER_MONEY_OVER)
				orderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, for_game.GOLD_TYPE_TRANSFER_MONEY_OVER)
				reason := for_game.GetGoldChangeNote(for_game.GOLD_TYPE_TRANSFER_MONEY_OVER, transferMoney.GetTargetId(), nil)
				msg := &share_message.GoldExtendLog{
					OrderId:      easygo.NewString(orderId),
					TransferText: easygo.NewString(transferMoney.GetContent()),
					HeadIcon:     easygo.NewString(target.GetHeadIcon()),
					PayType:      easygo.NewInt32(transferMoney.GetWay()),
					Title:        easygo.NewString(reason),
					Gold:         easygo.NewInt64(base.GetGold() + transferMoney.GetGold()),
				}
				NotifyAddGold(base.GetPlayerId(), transferMoney.GetGold(), reason, for_game.GOLD_TYPE_TRANSFER_MONEY_OVER, msg)
				//红包标识为超时状态
				title := "转账退款到账通知"
				text := fmt.Sprintf("收到一笔转账退款，退款金额￥%.2f。通过“我的”→“零钱”→“账单”可查看详情", float64(transferMoney.GetGold())/100)
				NoticeAssistant(transferMoney.GetSender(), 1, title, text)
				// transferMoney.SetState(for_game.TRANSFER_MONEY_BACK)
				transferMoney.SetState(for_game.TRANSFER_MONEY_BACK)
				transferMoney.SaveToMongo()
			}
		}
		for _, p := range packets {
			easygo.PCall(fc, p)
		}
	}
}

/*//绑定硬币过期检测
func (self *UpdateManager) UpdateBCoinExpiration() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	logs.Info("UpdateBCoinExpiration ----》》")
	var logs, exLogs []*share_message.PlayerBCoinLog
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BCOIN_LOG)
	defer closeFun()
	t := time.Now().Unix()
	err := col.Find(bson.M{"Status": for_game.BCOIN_STATUS_UNUSE, "OverTime": bson.M{"$lt": t}}).Sort("OverTime").All(&logs)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	//处理到期绑定硬币
	self.DealBCoinExpiration(logs, col)
	//将要到期的:10分钟提示一次，提前一天
	err = col.Find(bson.M{"Status": for_game.BCOIN_STATUS_UNUSE, "OverTime": bson.M{"$gt": t, "$lt": t + 86400}}).Sort("OverTime").All(&exLogs)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	//处理一天内将要过期的硬币
	self.DealPreBCoinExpiration(exLogs)

}

//处理过期硬币:扣除掉指定绑定硬币
func (self *UpdateManager) DealBCoinExpiration(data []*share_message.PlayerBCoinLog, col *mgo.Collection) {
	if len(data) <= 0 {
		return
	}
	logs.Info("处理过期硬币", data)
	total := int64(0)
	var saveLogs []interface{}
	playerBCoin := make(map[int64]int64)
	for _, log := range data {
		total += log.GetCurBCoin()
		playerBCoin[log.GetPlayerId()] += log.GetCurBCoin()
		log.Status = easygo.NewInt32(for_game.BCOIN_STATUS_EXPIRATION)
		log.CurBCoin = easygo.NewInt64(0)
		saveLogs = append(saveLogs, bson.M{"_id": log.GetId()}, log)
	}
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BCOIN_LOG, saveLogs)
	//减掉玩家的绑定硬币
	for pid, coin := range playerBCoin {
		NotifyAddCoin(pid, -coin, "系统回收", for_game.COIN_TYPE_SYSTEM_OUT, nil, false)
	}
}

//处理1天内将要过期硬币:给玩家发送推送通知
func (self *UpdateManager) DealPreBCoinExpiration(data []*share_message.PlayerBCoinLog) {
	if len(data) <= 0 {
		return
	}
	logs.Info("处理即将过期硬币", data)
	playerBCoin := make(map[int64]int64)
	for _, log := range data {
		playerBCoin[log.GetPlayerId()] += log.GetCurBCoin()
	}
	players := make([]int64, 0)
	for pid, coin := range playerBCoin {
		players = append(players, pid)
		content := fmt.Sprintf("你有%d硬币即将过期，请尽快使用!", coin)
		NoticeAssistant(pid, 1, "过期提示", content)
	}
	ids := for_game.GetJGIds(players)
	m := for_game.PushMessage{
		Title:       "过期提示",
		Content:     "平台赠送硬币即将过期",
		ContentType: for_game.JG_TYPE_BACKSTAGE_ASS,
		JumpObject:  3,
	}
	for_game.JGSendMessage(ids, m)
}
*/
