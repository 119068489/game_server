package backstage

import (
	"game_server/for_game"
	"game_server/pb/share_message"
)

//红包处理拆包
func DealRedPacket(msg *share_message.RedPacket) {
	total := msg.GetTotalMoney()
	num := int64(msg.GetTotalCount())
	var safe_total, money int64
	if num > 1 {
		for i := int64(1); i < num; i++ {
			tem := (total - (num-i)*int64(for_game.RED_PACKET_MIN_VALUE))
			if tem < (num - i) {
				money = 1
			} else {
				safe_total = tem / (num - i) //随机安全上限
				money = int64(for_game.RandInt(int(for_game.RED_PACKET_MIN_VALUE), int(safe_total)))
			}
			msg.Packets = append(msg.GetPackets(), money)
			total = total - money
		}
	}
	msg.Packets = append(msg.GetPackets(), total)
	msg.PlayerList = []int64{}
}

/*//测试
func testchatlog() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PERSONAL_CHAT_LOG)
	defer closeFun()

	queryBson := bson.M{"_id": 1887436004, "ChatInfo.TargetId": 1887436015} //bson.M{"$gte": 10}
	m := []bson.M{
		{"$unwind": "$ChatInfo"},
		{"$match": queryBson},
		{"$project": bson.M{"ChatInfo": 1, "_id": 0}},
	}

	query := col.Pipe(m)
	var list []interface{}
	err := query.All(&list)
	easygo.PanicError(err)

	logs.Info("******", list)

}
*/
