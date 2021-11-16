package for_game

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"log"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

var _ = fmt.Sprintf
var _ = log.Println
var _ = easygo.Underline

//================================================================================用户留存报表  已优化到Redis
//生成玩家留存报表 types 1 登录，2补全资料
func MakePlayerKeepReport(pid PLAYER_ID, types int) {
	player := GetRedisPlayerBase(pid)             //上线用户资料
	loginTime := easygo.GetToday0ClockTimestamp() //今天0点时间戳
	createTime := player.GetCreateTime()          //注册时间
	lastLogOutTime := player.GetLastLogOutTime()  //最后登出时间

	if len(player.GetLabelList()) == 0 {
		return
	}

	//更新今日注册用户数
	if types == 2 {
		SetRedisPlayerKeepReportFildVal(loginTime, 1, "TodayRegister")
	} else {
		if easygo.Get0ClockTimestamp(lastLogOutTime) != loginTime && lastLogOutTime > 0 {
			days := easygo.GetDifferenceDay(createTime, loginTime)
			createtime0timestamp := easygo.Get0ClockTimestamp(createTime) //注册日0点
			//更新留存人数
			switch days {
			case 1:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "NextKeep")
			case 2:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "ThreeKeep")
			case 3:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "FourKeep")
			case 4:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "FiveKeep")
			case 5:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "SixKeep")
			case 6:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "SevenKeep")
			case 7:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "EightKeep")
			case 8:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "NineKeep")
			case 9:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TenKeep")
			case 10:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "ElevenKeep")
			case 11:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwelveKeep")
			case 12:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "ThirteenKeep")
			case 13:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "FourteenKeep")
			case 14:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "FifteenKeep")
			case 15:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "SixteenKeep")
			case 16:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "SeventeenKeep")
			case 17:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "EighteenKeep")
			case 18:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "NineteenKeep")
			case 19:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwentyKeep")
			case 20:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwentyOneKeep")
			case 21:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwentyTwoKeep")
			case 22:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwentyThreeKeep")
			case 23:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwentyFourKeep")
			case 24:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwentyFiveKeep")
			case 25:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwentySixKeep")
			case 26:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwentySevenKeep")
			case 27:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwentyEightKeep")
			case 28:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "TwentyNineKeep")
			case 29:
				SetRedisPlayerKeepReportFildVal(createtime0timestamp, 1, "Thirtykeep")
			}
		}
	}
}

//===============================================================================用户活跃日志
//更新玩家在线时长日志
func AddOnlineTimeLog(pid PLAYER_ID) {
	onlineTime := int64(0)            //在线时长
	player := GetRedisPlayerBase(pid) //从Redis中获取用户资料
	if player == nil {
		return
	}
	lastlogintime := player.GetLastOnLineTime() / 1000    //最后上线时间
	lastloginouttime := player.GetLastLogOutTime() / 1000 //最后下线时间
	sumonlinetime := lastloginouttime - lastlogintime     //总在线时长
	channelNo := player.GetChannel()
	days := easygo.GetDifferenceDay(lastlogintime, lastloginouttime)

	for i := int32(0); i <= days; i++ {
		if i == 0 {
			if days == 0 {
				onlineTime = lastloginouttime - lastlogintime
			} else {
				onlineTime = lastloginouttime - easygo.Get0ClockTimestamp(lastloginouttime) + 1
			}
			EditOnlineTimeLog(pid, lastloginouttime, onlineTime)
			if channelNo != "" {
				SetRedisOperationChannelReportFildVal(lastloginouttime, onlineTime, channelNo, "OnlineSum")
			}
		} else {
			sumonlinetime = sumonlinetime - onlineTime
			if sumonlinetime-86400 > 0 {
				lastloginouttime = lastloginouttime - onlineTime - 1
				onlineTime = 86400
				EditOnlineTimeLog(pid, lastloginouttime, onlineTime)
				if channelNo != "" {
					SetRedisOperationChannelReportFildVal(lastloginouttime, onlineTime, channelNo, "OnlineSum")
				}
			} else {
				onlineTime := easygo.Get24ClockTimestamp(lastlogintime) - lastlogintime
				EditOnlineTimeLog(pid, lastlogintime, onlineTime)
				if channelNo != "" {
					SetRedisOperationChannelReportFildVal(lastlogintime, onlineTime, channelNo, "OnlineSum")
				}
			}
		}
	}
}

//指定时间查询玩家在线时长日志
func GetOnlineTimeLog(pid PLAYER_ID, logouttime int64) *share_message.OnlineTimeLog {
	logtime := easygo.Get0ClockTimestamp(logouttime)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ONLINETIMELOG)
	defer closeFun()
	var obj *share_message.OnlineTimeLog
	err := col.Find(bson.M{"PlayerId": pid, "CreateTime": logtime}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//修改玩家在线时长日志
func EditOnlineTimeLog(pid, lastloginouttime, onlinetime int64) {
	log := GetOnlineTimeLog(pid, lastloginouttime)
	if log == nil {
		log = &share_message.OnlineTimeLog{
			Id:         easygo.NewInt64(NextId(TABLE_ONLINETIMELOG)),
			PlayerId:   easygo.NewInt64(pid),
			CreateTime: easygo.NewInt64(easygo.Get0ClockTimestamp(lastloginouttime)),
			OnlineTime: easygo.NewInt64(onlinetime),
		}
	} else {
		log.OnlineTime = easygo.NewInt64(log.GetOnlineTime() + onlinetime)
	}

	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ONLINETIMELOG)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": log.GetId()}, bson.M{"$set": log})
	easygo.PanicError(err)
}

//================================================================================用户行为报表 已优化到Redis
//生成用户行为报表  如果是聊天pid就传0
//types： 1聊天，2玩家登录登出，3实名绑卡，4商城订单完成,5领取红包，6转账,7发红包
//pid：1聊天类型传0,4商城订单类型传0，5领红包传0,6转账类型传0,7发红包类型传0
func MakePlayerBehaviorReport(types int, pid PLAYER_ID, redPacket *share_message.RedPacket, shopOrder *share_message.TableShopOrder, RobRedpacket *share_message.RedPacketLog, transferMoney *share_message.TransferMoney) {
	querytime := easygo.GetToday0ClockTimestamp()
	switch types {
	case 1:
		SetRedisPlayerBehaviorReportFildVal(querytime, 1, "SendMsgCount")
	case 2:
		player := GetPlayerById(pid)
		if player.GetLastLogOutTime()-player.GetCreateTime() < 120000 && player.GetLoginTimes() == 1 {
			SetRedisPlayerBehaviorReportFildVal(querytime, 1, "OneDialogue")
		}
	case 3:
		player := GetPlayerById(pid)
		if len(player.GetBankInfo()) == 1 {
			SetRedisPlayerBehaviorReportFildVal(querytime, 1, "BindCard")
		}
	case 4:
		SetRedisPlayerBehaviorReportFildVal(querytime, 1, "ShopOrderCount")
		items := shopOrder.Items
		orderMoney := items.GetPrice() * items.GetCount()
		SetRedisPlayerBehaviorReportFildVal(querytime, int64(orderMoney), "ShopOrderMoney")

	case 5:
		if !IsSendRedpacket(RobRedpacket.GetPlayerId(), GOLD_TYPE_GET_REDPACKET) {
			SetRedisPlayerBehaviorReportFildVal(querytime, 1, "RobRedpacketPlayerCount")
		}
		SetRedisPlayerBehaviorReportFildVal(querytime, 1, "RobRedpacketCount")
		SetRedisPlayerBehaviorReportFildVal(querytime, RobRedpacket.GetMoney(), "RobRedpacketMoney")
	case 6:
		if !IsSendRedpacket(transferMoney.GetSender(), GOLD_TYPE_SEND_TRANSFER_MONEY) {
			SetRedisPlayerBehaviorReportFildVal(querytime, 1, "TransferPlayerCount")
		}
		SetRedisPlayerBehaviorReportFildVal(querytime, 1, "TransferCount")
		SetRedisPlayerBehaviorReportFildVal(querytime, transferMoney.GetGold(), "TransferMoney")
	case 7:
		if !IsSendRedpacket(redPacket.GetSender(), GOLD_TYPE_SEND_REDPACKET) {
			SetRedisPlayerBehaviorReportFildVal(querytime, 1, "SendRedpacketPlayerCount")
		}
		SetRedisPlayerBehaviorReportFildVal(querytime, 1, "SendRedpacketCount")
		SetRedisPlayerBehaviorReportFildVal(querytime, redPacket.GetTotalMoney(), "SendRedpacketMoney")
	}
}

//查询玩家今天是否发过红包或转账  types { 发红包:GOLD_TYPE_SEND_REDPACKET , 转出:GOLD_TYPE_SEND_TRANSFER_MONEY }
func IsSendRedpacket(id PLAYER_ID, types int32) bool {
	time0 := easygo.GetToday0ClockTimestamp() * 1000
	time24 := easygo.GetToday24ClockTimestamp() * 1000
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_GOLDCHANGELOG)
	defer closeFun()
	queryBson := bson.M{"PlayerId": id}
	queryBson["SourceType"] = types
	queryBson["CreateTime"] = bson.M{"$gte": time0, "$lte": time24}

	var obj *share_message.Order
	err := col.Find(queryBson).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return false
	}
	return true
}

//==================================================================================出入款汇总报表 已优化到Redis
//生成出入款汇总报表
func MakeInOutCashSumReport(sourceType int32, gold int64) {
	querytime := easygo.GetToday0ClockTimestamp()

	recharge := int64(0)
	withdraw := int64(0)

	switch sourceType {
	case GOLD_TYPE_CASH_AFIN:
		recharge += gold
	case GOLD_TYPE_CASH_IN:
		recharge += gold
	case GOLD_TYPE_CASH_AFOUT:
		withdraw += gold
	case GOLD_TYPE_CASH_OUT:
		withdraw += gold
	}
	redundant := recharge + withdraw
	if recharge != 0 {
		SetRedisInOutCashSumReportFildVal(querytime, recharge, "Recharge")
		SetRedisInOutCashSumReportFildVal(querytime, 1, "RechargeTimes")

	}
	if withdraw != 0 {
		SetRedisInOutCashSumReportFildVal(querytime, withdraw, "Withdraw")
		SetRedisInOutCashSumReportFildVal(querytime, 1, "WithdrawTimes")

	}
	if redundant != 0 {
		SetRedisInOutCashSumReportFildVal(querytime, redundant, "Redundant")
	}
}

//更新出入款汇总充值提现人数
func UpdateInOutCashCount() {
	querytime := easygo.GetYesterday0ClockTimestamp()
	endtime := easygo.GetYesterday24ClockTimestamp()
	M := []bson.M{
		{"$match": bson.M{"SourceType": bson.M{"$in": []int32{GOLD_TYPE_CASH_AFIN, GOLD_TYPE_CASH_IN}}, "CreateTime": bson.M{"$gte": querytime * 1000, "$lt": endtime * 1000}, "Status": ORDER_ST_FINISH}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	count := FindPipeAllCount(MONGODB_NINGMENG, TABLE_ORDER, M)
	if count > 0 {
		UpdateRedisInOutCashSumReportFildVal(querytime, count, "RechargeCount")
	}

	oM := []bson.M{
		{"$match": bson.M{"SourceType": bson.M{"$in": []int32{GOLD_TYPE_CASH_AFOUT, GOLD_TYPE_CASH_OUT}}, "CreateTime": bson.M{"$gte": querytime * 1000, "$lt": endtime * 1000}, "Status": ORDER_ST_FINISH}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": 1}}},
		{"$group": bson.M{"_id": nil, "Count": bson.M{"$sum": 1}}},
	}
	count1 := FindPipeAllCount(MONGODB_NINGMENG, TABLE_ORDER, oM)
	if count1 > 0 {
		UpdateRedisInOutCashSumReportFildVal(querytime, count1, "WithdrawCount")
	}
}

//=================================================================================埋点注册登录报表 已优化到Redis
//生成埋点注册登录报表 types{ 1:密码登录 2:验证码登录 4:一键登录 5:微信登录 6:自动登录 7:手机号码注册 8:一键登录注册 9:微信登录注册 }
func MakeRegisterLoginReport(reqMsg *share_message.LoginRegisterInfo) {
	// debug.PrintStack()
	today0timestamp := easygo.Get0ClockTimestamp(reqMsg.GetTime())
	player := GetRedisPlayerBase(reqMsg.GetPlayerId())

	switch reqMsg.GetType() {
	case 1, 2, 4, 5, 6:
		SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "LoginTimesCount")
		if easygo.Get0ClockTimestamp(player.GetLastLogOutTime()) < today0timestamp {
			SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "LoginSumCount")
		}
	case 7, 8, 9:
		SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "RegSumCount")
	}

	switch reqMsg.GetType() {
	case 1, 2:
		SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "PhoneLoginCount")
	case 4:
		SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "OneClickLoginCount")
	case 5:
		SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "WxLoginCount")
	case 6:
		SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "AoutLoginCount")
	case 7, 8:
		SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "PhoneRegCount")
	case 9:
		SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "WxRegCount")
	}
}

//验证是否是有效激活设备
func VerDeviceCode(code string, createtime int64) bool {
	if code == "" {
		return true
	}

	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_POS_DEVICECODE)
	defer closeFun()
	queryBson := bson.M{"DeviceCode": code}
	dc := &share_message.PosDeviceCode{}
	err := col.Find(queryBson).One(&dc)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	players := QueryPlayersByDeviceCode(code)
	if len(players) > 0 {
		return false
	}

	return true
}

//更新埋点报表绑卡人数
func SetRegisterLoginReportBankCardCount(pid PLAYER_ID) {
	today0timestamp := easygo.Get0ClockTimestamp(util.GetMilliTime())
	player := GetRedisPlayerBase(pid)
	lis := player.GetBankInfos()
	changeCount := 1
	for _, i := range lis {
		if easygo.Get0ClockTimestamp(i.GetTime()) == today0timestamp {
			changeCount = 0
		}
	}
	if changeCount == 0 {
		SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "BankCardCount") //更新埋点报表绑卡人数
	}
}

//设置pv uv数 id：社交广场动态id， code：设备码
func SetPvUvCount(id int64, code string) {
	today0timestamp := easygo.Get0ClockTimestamp(util.GetMilliTime())
	SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "PvCount") //更新pv数
	key := code + "_" + easygo.AnytoA(id)
	if !IsReadDynamicDevice(key) {
		SetRedisRegisterLoginReportFildVal(today0timestamp, 1, "UvCount") //更新uv数
		msg := &share_message.ReadDynamicDevice{
			Id:    easygo.NewString(key),
			Code:  easygo.NewString(code),
			LogId: easygo.NewInt64(id),
		}
		SaveReadDynamicDevice(msg)
	}
}

//检查设备是否读过动态
func IsReadDynamicDevice(id string) bool {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_READ_DYNAMIC_DEVICE)
	defer closeFun()

	count, err := col.Find(bson.M{"_id": id}).Count()
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	if count > 0 {
		return true
	}
	return false
}

//写设备读取动态记录
func SaveReadDynamicDevice(req *share_message.ReadDynamicDevice) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_READ_DYNAMIC_DEVICE)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": req.GetId()}, bson.M{"$set": req})
	easygo.PanicError(err)
}

//=========================================================================================================运营渠道数据汇总报表 已优化到Redis
//生成运营渠道数据汇总报表 types 1注册，6登录，querytime时间， 2下载，3uv，4商城完成订单，5充值提现
func MakeOperationChannelReport(types int32, playerId PLAYER_ID, channelNo string, shopOrder *share_message.TableShopOrder, order *share_message.Order) {
	today0timestamp := easygo.Get0ClockTimestamp(util.GetMilliTime())
	player := GetRedisPlayerBase(playerId)
	createtime0timestamp := easygo.Get0ClockTimestamp(player.GetCreateTime()) //注册时间0点时间戳

	if player.GetChannel() == "" {
		return
	} else {
		if channelNo == "" {
			channelNo = player.GetChannel()
		}
	}

	// report := GetRedisOperationChannelReport(channelNo, today0timestamp)
	// if report == nil {
	// 	channel := QueryOperationByNo(channelNo) //查询渠道信息
	// 	if channel == nil {
	// 		return
	// 	}
	// 	report = &share_message.OperationChannelReport{
	// 		Id:          easygo.NewInt64(NextId(TABLE_OPERATION_CHANNEL_REPORT)),
	// 		ChannelNo:   easygo.NewString(channelNo),
	// 		ChannelName: easygo.NewString(channel.GetName()),
	// 		Cooperation: easygo.NewInt32(channel.GetCooperation()),
	// 		CreateTime:  easygo.NewInt64(today0timestamp),
	// 	}
	// }

	switch types {
	case 1:
		// report.RegCount = easygo.NewInt64(report.GetRegCount() + 1)
		SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "RegCount")
		//设备是否是有效激活
		if VerDeviceCode(player.GetDeviceCode(), today0timestamp) {
			// report.ValidActDevCount = easygo.NewInt64(report.GetValidActDevCount() + 1)
			SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "ValidActDevCount")
		}
	case 6:
		if player.GetLastLogOutTime() == 0 {
			// report.RegCount = easygo.NewInt64(report.GetRegCount() + 1)
			SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "RegCount")
			//设备是否是有效激活
			if VerDeviceCode(player.GetDeviceCode(), today0timestamp) {
				// report.ValidActDevCount = easygo.NewInt64(report.GetValidActDevCount() + 1)
				SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "ValidActDevCount")
			}
		}
		// report.LoginCount = easygo.NewInt64(report.GetLoginCount() + 1)

		if easygo.Get0ClockTimestamp(player.GetLastLogOutTime()) < today0timestamp {
			SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "LoginCount")
		}

		if easygo.Get0ClockTimestamp(player.GetLastLogOutTime()) < today0timestamp && easygo.Get0ClockTimestamp(player.GetCreateTime()) == easygo.Get0ClockTimestamp(easygo.NowTimestamp()-86400) {
			SetRedisOperationChannelReportFildVal(easygo.Get0ClockTimestamp(easygo.NowTimestamp()-86400), 1, channelNo, "NextKeep")
		}
		if player.GetTypes() == 6 {
			rate := util.RandIntn(10)
			if rate < 4 {
				SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "NextKeep")
			}
		}
	case 2:
		// report.DownLoadCount = easygo.NewInt64(report.GetDownLoadCount() + 1)
		SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "DownLoadCount")
	case 3:
		// report.UvCount = easygo.NewInt64(report.GetUvCount() + 1)
		SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "UvCount")
	case 4:
		// report.ShopOrderSumCount = easygo.NewInt64(report.GetShopOrderSumCount() + 1)
		SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "ShopOrderSumCount")
		items := shopOrder.Items
		orderMoney := items.GetPrice() * items.GetCount()
		// report.ShopDealSumAmount = easygo.NewInt64(report.GetShopDealSumAmount() + int64(orderMoney))
		SetRedisOperationChannelReportFildVal(today0timestamp, int64(orderMoney), channelNo, "ShopDealSumAmount")

		if today0timestamp == createtime0timestamp {
			// report.ShopOrderNewCount = easygo.NewInt64(report.GetShopOrderNewCount() + 1)
			SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "ShopOrderNewCount")
			// report.ShopDealNewAmount = easygo.NewInt64(report.GetShopDealNewAmount() + int64(orderMoney))
			SetRedisOperationChannelReportFildVal(today0timestamp, int64(orderMoney), channelNo, "ShopDealNewAmount")
		} else {
			// report.ShopOrderOldCount = easygo.NewInt64(report.GetShopOrderOldCount() + 1)
			SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "ShopOrderOldCount")
			// report.ShopDealOldAmount = easygo.NewInt64(report.GetShopDealOldAmount() + int64(orderMoney))
			SetRedisOperationChannelReportFildVal(today0timestamp, int64(orderMoney), channelNo, "ShopDealOldAmount")
		}
		if IsFinishShopOrder(shopOrder.GetReceiverId()) {
			// report.ShopOrderPlayerCount = easygo.NewInt64(report.GetShopOrderPlayerCount() + 1)
			SetRedisOperationChannelReportFildVal(today0timestamp, 1, channelNo, "ShopOrderPlayerCount")
		}
	case 5:
		switch order.GetSourceType() {
		case GOLD_TYPE_CASH_IN:
			// report.RechargeSumAmount = easygo.NewInt64(report.GetRechargeSumAmount() + order.GetChangeGold())
			SetRedisOperationChannelReportFildVal(today0timestamp, order.GetChangeGold(), channelNo, "RechargeSumAmount")
			if today0timestamp == createtime0timestamp {
				// report.RechargeNewAmount = easygo.NewInt64(report.GetRechargeSumAmount() + order.GetChangeGold())
				SetRedisOperationChannelReportFildVal(today0timestamp, order.GetChangeGold(), channelNo, "RechargeNewAmount")
			} else {
				// report.RechargeOldAmount = easygo.NewInt64(report.GetRechargeOldAmount() + order.GetChangeGold())
				SetRedisOperationChannelReportFildVal(today0timestamp, order.GetChangeGold(), channelNo, "RechargeOldAmount")
			}
		case GOLD_TYPE_CASH_OUT:
			// report.WithdrawSumAmount = easygo.NewInt64(report.GetWithdrawSumAmount() - order.GetChangeGold())
			SetRedisOperationChannelReportFildVal(today0timestamp, order.GetChangeGold(), channelNo, "WithdrawSumAmount")
			if today0timestamp == createtime0timestamp {
				// report.WithdrawNewAmount = easygo.NewInt64(report.GetWithdrawNewAmount() - order.GetChangeGold())
				SetRedisOperationChannelReportFildVal(today0timestamp, order.GetChangeGold(), channelNo, "WithdrawNewAmount")
			} else {
				// report.WithdrawOldAmount = easygo.NewInt64(report.GetWithdrawOldAmount() - order.GetChangeGold())
				SetRedisOperationChannelReportFildVal(today0timestamp, order.GetChangeGold(), channelNo, "WithdrawOldAmount")
			}
		}
	}
}

//渠道列表
func GetChannelListNopage() []*brower_backstage.KeyValue {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_OPERATION_CHANNEL)
	defer closeFun()

	queryBson := bson.M{}
	query := col.Find(queryBson)
	var list []*share_message.OperationChannel
	err := query.Sort("-CreateTime").All(&list)
	easygo.PanicError(err)

	var lis []*brower_backstage.KeyValue
	for _, i := range list {
		li := &brower_backstage.KeyValue{
			Key:   i.ChannelNo,
			Value: i.Name,
		}
		lis = append(lis, li)
	}

	return lis
}

//查询玩家今天是否完成过商城订单
func IsFinishShopOrder(pid PLAYER_ID) bool {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_ORDERS)
	defer closeFun()

	startTime := easygo.GetToday0ClockTimestamp()
	endTime := easygo.GetToday24ClockTimestamp()

	queryBson := bson.M{"receiver_id": pid, "finish_time": bson.M{"$gte": startTime, "$lt": endTime}}
	query := col.Find(queryBson)
	var list []*share_message.TableShopOrder
	err := query.Sort("-_id").All(&list)
	easygo.PanicError(err)

	if len(list) == 0 {
		return true
	}

	return false
}

//从mongo中读取转账
func GetOperationChannelReport(channelNo string, createTime int64) *share_message.OperationChannelReport {
	querytime := easygo.Get0ClockTimestamp(createTime)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_OPERATION_CHANNEL_REPORT)
	defer closeFun()
	var obj *share_message.OperationChannelReport
	err := col.Find(bson.M{"CreateTime": querytime, "ChannelNo": channelNo}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//新增修改小助手文章
func EditArticle(reqMsg *share_message.Article) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ARTICLE)

	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetID()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//根据id查询文章
func QueryArticleById(id int64) *share_message.Article {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ARTICLE)
	defer closeFun()
	var article *share_message.Article
	err := col.Find(bson.M{"_id": id}).One(&article)

	if err != nil && err == mgo.ErrNotFound {
		return nil
	} else {
		easygo.PanicError(err)
	}
	return article
}

//更新文章阅读数，并返回
func ReadArticle(playerId, articleId int64) *share_message.Article {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ARTICLE)
	defer closeFun()
	//查出创建时间
	var article *share_message.Article
	err := col.Find(bson.M{"_id": articleId}).Select(bson.M{"Create_time": 1}).One(&article)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	virtual := 0
	if GetMillSecond()-article.GetCreateTime() < 6*3600000 {
		//6小时内随机增加3-5个虚拟增量
		virtual = RandInt(3, 6)
	}
	//修改新的阅读数
	_, err = col.Find(bson.M{"_id": articleId}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"ReadedNum": 1, "ReadedNumVirtual": virtual}},
		Upsert:    true,
		ReturnNew: true,
	}, &article)
	easygo.PanicError(err)
	//更新增在阅读数
	SetArticleReadingNum(article)
	//检测是否点赞了
	if playerId != 0 {
		article.IsZan = easygo.NewBool(GetArticleIsZan(playerId, articleId))
	}
	return article
}

//获取文章正在阅读数
func SetArticleReadingNum(article *share_message.Article) {
	t := GetMillSecond() - article.GetCreateTime()
	readNum := article.GetReadedBase() + article.GetReadedNum() + article.GetReadedNumVirtual()
	if t < ONE_HOUR_MILLSECOND {
		//小于1小时，80%的正在阅读
		readNum = readNum * 80 / 100
	} else if t >= ONE_HOUR_MILLSECOND && t < SIX_HOUR_MILLSECOND {
		//1-6小时显示
		n := t / ONE_HOUR_MILLSECOND
		readNum = readNum * (50 - (n-1)*10) / 100
	} else if t >= SIX_HOUR_MILLSECOND && t < TWELVE_HOUR_MILLSECOND {
		//6-12小时显示
		n := t / ONE_HOUR_MILLSECOND
		readNum = readNum * (6 - (n - 6)) / 100
	} else {
		readNum = 1
	}
	article.ReadingNum = easygo.NewInt64(readNum)
}

//检测玩家是否点赞了
func GetArticleIsZan(playerId, articleId int64) bool {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ARTICLE_ZAN)
	defer closeFun()
	var data *share_message.ArticleZan
	err := col.Find(bson.M{"PlayerId": playerId, "ArticleId": articleId}).One(&data)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return false
	}
	return true
}

//点赞文章
func ZanArticle(playerId, articleId int64) {
	player := GetRedisPlayerBase(playerId)
	if player == nil {
		logs.Error("找不到点赞玩家的数据")
		return
	}
	msg := share_message.ArticleZan{
		Id:         easygo.NewInt64(NextId(TABLE_ARTICLE_ZAN)),
		ArticleId:  easygo.NewInt64(articleId),
		PlayerId:   easygo.NewInt64(playerId),
		Name:       easygo.NewString(player.GetNickName()),
		HeadUrl:    easygo.NewString(player.GetHeadIcon()),
		CreateTime: easygo.NewInt64(GetMillSecond()),
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ARTICLE_ZAN)
	defer closeFun()
	err := col.Insert(msg)
	easygo.PanicError(err)
	//点赞数+1
	UpdateArticleZanNum(articleId, 1)
}

//更新点赞数:num正数位增加，负数为减少
func UpdateArticleZanNum(articleId int64, num int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ARTICLE)
	defer closeFun()
	err := col.Update(bson.M{"_id": articleId}, bson.M{"$inc": bson.M{"ZanNum": num}})
	easygo.PanicError(err)
}

//获取评论:分页处理
func GetArticleComment(articleId int64, pagesize, curpage int64, status ...int32) []*share_message.ArticleComment {
	pageSize := int(pagesize)
	curPage := easygo.If(int(curpage) > 1, int(curpage)-1, 0).(int)

	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ARTICLE_COMMENT)
	defer closeFun()
	var comments []*share_message.ArticleComment
	qurybosn := bson.M{"ArticleId": articleId}
	if len(status) > 0 {
		qurybosn["Status"] = status[0]
	}
	query := col.Find(qurybosn)
	var err error
	if pagesize > 0 {
		err = query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&comments)
	} else {
		err = query.Sort("-CreateTime").All(&comments)
	}
	easygo.PanicError(err)
	return comments
}

//玩家评论文章
func CommentArticle(playerId, articleId int64, content string) {
	player := GetRedisPlayerBase(playerId)
	if player == nil {
		logs.Error("找不到评论玩家的数据")
		return
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ARTICLE_COMMENT)
	defer closeFun()
	msg := share_message.ArticleComment{
		Id:         easygo.NewInt64(NextId(TABLE_ARTICLE_COMMENT)),
		ArticleId:  easygo.NewInt64(articleId),
		Content:    easygo.NewString(content),
		PlayerId:   easygo.NewInt64(playerId),
		Name:       easygo.NewString(player.GetNickName()),
		HeadUrl:    easygo.NewString(player.GetHeadIcon()),
		CreateTime: easygo.NewInt64(GetMillSecond()),
		Status:     easygo.NewInt32(ARTICLE_COMMENT_SHOW), //默认隐藏
	}
	err := col.Insert(msg)
	easygo.PanicError(err)
}

//========================================================================================================通知报表

//根据id查询文章
func QueryNoticeById(id int32) *share_message.AppPushMessage {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_FEATURES_APPPUSHMSG)
	defer closeFun()
	var articles *share_message.AppPushMessage
	err := col.Find(bson.M{"_id": id}).One(&articles)

	if err != nil && err == mgo.ErrNotFound {
		return nil
	} else {
		easygo.PanicError(err)
	}
	return articles
}

//=======================================================玩家在线报表
func UpdatePlayerOnlineReport(count int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYERONLINE_REPORT)
	defer closeFun()
	idTime := easygo.Get0ClockTimestamp(easygo.NowTimestamp())
	qurybson := bson.M{}
	fild := "Clock" + easygo.AnytoA(time.Now().Hour())
	qurybson[fild] = count

	_, err := col.Upsert(bson.M{"_id": idTime}, bson.M{"$set": qurybson})
	easygo.PanicError(err)
}

//=========================================================玩家登录地域报表
func MakePlayerLogLocationReport(deviceType int32, ip string) {
	if ip == "" || deviceType == 0 {
		return
	}
	data := IpSearch(ip)
	if data == nil {
		return
	}
	ctime := easygo.Get0ClockTimestamp(easygo.NowTimestamp())
	report := &share_message.PlayerLogLocationReport{
		DayTime:    easygo.NewInt64(ctime),
		DeviceType: easygo.NewInt32(deviceType),
		Piece:      easygo.NewString(data.CountryId),
		Position:   easygo.NewString(data.Region),
	}
	if data.CountryId != "CN" {
		report.Position = easygo.NewString(data.Country)
	}
	UpdatePlayerLogLocationReport(report)
}

func UpdatePlayerLogLocationReport(reqMsg *share_message.PlayerLogLocationReport) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYERLOG_LOCATION_REPORT)
	defer closeFun()

	_, err := col.Find(bson.M{"DayTime": reqMsg.GetDayTime(), "Position": reqMsg.GetPosition(), "DeviceType": reqMsg.GetDeviceType()}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"Count": 1}, "$set": reqMsg},
		Upsert:    true,
		ReturnNew: false,
	}, nil)
	easygo.PanicError(err)
}

//=============================================================================================广告报表
func QueryAdvLogByTime(startTime, endTime int64) []*share_message.AdvLogReq {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_ADV_LOG)
	defer closeFun()

	queryBson := bson.M{"OpTime": bson.M{"$gte": startTime, "$lt": endTime}}
	query := col.Find(queryBson)
	var list []*share_message.AdvLogReq
	err := query.Sort("OpTime").All(&list)
	easygo.PanicError(err)

	return list
}

//=============================================================================================附近的人引导项报表
func QueryNearbyAdvLogByTime(startTime, endTime int64) []*share_message.AdvLogReq {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_ADV_LOG)
	defer closeFun()

	queryBson := bson.M{"OpTime": bson.M{"$gte": startTime, "$lt": endTime}}
	query := col.Find(queryBson)
	var list []*share_message.AdvLogReq
	err := query.Sort("OpTime").All(&list)
	easygo.PanicError(err)

	return list
}

//==================================================================================短信召回统计报表
//写老用户回归日志
func AddRecallPlayerLog(pid PLAYER_ID) {
	if pid <= 0 {
		logs.Error("回归玩家id错误:" + easygo.AnytoA(pid))
		return
	}
	now := easygo.NowTimestamp()
	id := easygo.AnytoA(pid) + easygo.AnytoA(easygo.Get0ClockTimestamp(now))
	data := share_message.RecallPlayerLog{
		Id:         easygo.NewString(id),
		PlayerId:   easygo.NewInt64(pid),
		RecallTime: easygo.NewInt64(now),
	}
	queryBson := bson.M{"_id": id}
	updateBson := bson.M{"$set": data}
	one := FindAndModify(MONGODB_NINGMENG_LOG, TABLE_RECALLPLAYER_LOG, queryBson, updateBson, true)
	if one != nil {
		SetRedisRecallReportFildVal(easygo.Get0ClockTimestamp(now), 1, "RecallCount")
	}
}
