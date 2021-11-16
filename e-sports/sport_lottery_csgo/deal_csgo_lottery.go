package sport_lottery_csgo

import (
	"fmt"
	dal_common "game_server/e-sports/sport_common_dal"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
)

func DealESportCSGOLottery(uniqueGameId int64) {

	startStr := fmt.Sprintf("======DealESportCSGOLottery中比赛uniqueGameId:%v的csgo开奖  开始==========", uniqueGameId)
	for_game.WriteFile("csgo_sport_lottery.log", startStr)

	//函数退出就要做、不管是不是发布的比赛、所有比赛要设置成已开奖
	defer DealGameLottery(uniqueGameId)

	//通过比赛唯一id取得该比赛的信息
	dbESPortsGame := share_message.TableESPortsGame{}
	colGame, closeFunGame := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
	defer closeFunGame()

	errGame := colGame.Find(bson.M{"_id": uniqueGameId, "release_flag": for_game.GAME_RELEASE_FLAG_2}).One(&dbESPortsGame)

	if errGame != nil && errGame != mgo.ErrNotFound {
		s := fmt.Sprintf("=====csgo开奖取得比赛表数据库错误:errGame错误:%v=====,查询条件为uniqueGameId:%v,release_flag:%v",
			errGame, uniqueGameId, for_game.GAME_RELEASE_FLAG_2)
		logs.Error(s)
		for_game.WriteFile("csgo_sport_lottery.log", s)
		return
	}

	if errGame == mgo.ErrNotFound {
		s := fmt.Sprintf("===========csgo开奖未取得该比赛的发布的数据、清key和设置开奖状态后  直接返回==========,查询条件为uniqueGameId:%v,release_flag:%v", uniqueGameId, for_game.GAME_RELEASE_FLAG_2)
		for_game.WriteFile("csgo_sport_lottery.log", s)
		return
	}

	//先判断有没有投注、没有投注
	dbBetList := make([]*share_message.TableESPortsGuessBetRecord, 0)
	colBet, closeFunBet := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD)
	defer closeFunBet()
	errColBet := colBet.Find(bson.M{"UniqueGameId": uniqueGameId, "BetResult": for_game.GAME_GUESS_BET_RESULT_1}).All(&dbBetList)

	if errColBet != nil && errColBet != mgo.ErrNotFound {

		s := fmt.Sprintf("=======csgo开奖取得投注记录数据库错误:errColBet错误:%v========,查询条件为uniqueGameId:%v,BetResult:%v",
			errColBet, uniqueGameId, for_game.GAME_GUESS_BET_RESULT_1)
		logs.Error(s)
		for_game.WriteFile("csgo_sport_lottery.log", s)
		return
	}

	if errColBet == mgo.ErrNotFound || nil == dbBetList || len(dbBetList) <= 0 {
		s := fmt.Sprintf("=======csgo开奖未取得未结算的投注记录的数据、清key和设置开奖状态后  直接返回=======,查询条件为uniqueGameId:%v,BetResult:%v",
			uniqueGameId, for_game.GAME_GUESS_BET_RESULT_1)
		for_game.WriteFile("csgo_sport_lottery.log", s)

		return
	}

	//取得数据库该比赛的赔率信息
	dbGuessList := make([]*share_message.TableESPortsGameGuess, 0)

	colGuess, closeFunGuess := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_GUESS)
	defer closeFunGuess()

	errColGuess := colGuess.Find(bson.M{"app_label_id": dbESPortsGame.GetAppLabelId(),
		"game_id":    dbESPortsGame.GetGameId(),
		"api_origin": dbESPortsGame.GetApiOrigin()}).All(&dbGuessList)

	if errColGuess != nil && errColGuess != mgo.ErrNotFound {

		s := fmt.Sprintf("=======csgo开奖取得赔率数据库错误:errColGuess错误:%v=======,查询条件为app_label_id:%v,game_id:%v,api_origin:%v",
			errColGuess, dbESPortsGame.GetAppLabelId(), dbESPortsGame.GetGameId(), dbESPortsGame.GetApiOrigin())
		logs.Error(s)
		for_game.WriteFile("csgo_sport_lottery.log", s)
		return
	}

	if errColGuess == mgo.ErrNotFound || nil == dbGuessList || len(dbGuessList) <= 0 {
		s := fmt.Sprintf("======csgo开奖未取得赔率数据、清key和设置开奖状态后  直接返回======,查询条件为app_label_id:%v,game_id:%v,api_origin:%v",
			dbESPortsGame.GetAppLabelId(), dbESPortsGame.GetGameId(), dbESPortsGame.GetApiOrigin())
		for_game.WriteFile("csgo_sport_lottery.log", s)
		return
	}

	//将早盘的所有信息封装到一个itemRstMornMap结构体中
	itemRstMornMap := make(map[string]*client_hall.ItemResult)
	//将滚盘的所有信息封装到一个itemRstRollMap结构体中
	itemRstRollMap := make(map[string]*client_hall.ItemResult)

	if errColGuess == nil && len(dbGuessList) > 0 {
		for _, value := range dbGuessList {
			guessList := value.GetGuess()
			if guessList != nil && len(guessList) > 0 {
				for _, guessValue := range guessList {
					itemList := guessValue.GetItems()
					if itemList != nil && len(itemList) > 0 {
						for _, itemValue := range itemList {
							//早盘
							if value.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_1 {
								itemRstMornMap[itemValue.GetBetNum()] = &client_hall.ItemResult{
									BetNum:     easygo.NewString(itemValue.GetBetNum()),
									Status:     easygo.NewString(itemValue.GetStatus()),
									Win:        easygo.NewString(itemValue.GetWin()),
									StatusTime: easygo.NewInt64(itemValue.GetStatusTime()),
								}
							} else {
								//滚盘
								itemRstRollMap[itemValue.GetBetNum()] = &client_hall.ItemResult{
									BetNum:     easygo.NewString(itemValue.GetBetNum()),
									Status:     easygo.NewString(itemValue.GetStatus()),
									Win:        easygo.NewString(itemValue.GetWin()),
									StatusTime: easygo.NewInt64(itemValue.GetStatusTime()),
								}
							}
						}
					}
				}
			}
		}
	}

	//批量消息用、投注该比赛的所有用户
	playerMapMsgs := make(map[int64][]*share_message.TableESPortsGameOrderSysMsg)
	//正式开奖
	if dbBetList != nil && len(dbBetList) > 0 {

		for _, betValue := range dbBetList {
			//先判断比赛中的项
			//早盘比赛开始后投注为时段无效
			if betValue.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_1 {
				if dbESPortsGame.GetBeginTimeInt() < betValue.GetCreateTime() {
					//时段无效
					DealGameInvalid(betValue, uniqueGameId, for_game.GAME_GUESS_BET_DISABLE_CODE_2, "早盘比赛开始后投注", playerMapMsgs)
					continue
				}
			} else {
				//滚盘比赛开始前投注为时段无效
				if dbESPortsGame.GetBeginTimeInt() > betValue.GetCreateTime() {
					//时段无效
					DealGameInvalid(betValue, uniqueGameId, for_game.GAME_GUESS_BET_DISABLE_CODE_2, "滚盘比赛开始前投注", playerMapMsgs)
					continue
				}
			}

			//比赛已经结束后的投注 时段无效
			if dbESPortsGame.GetGameStatus() == for_game.GAME_STATUS_2 &&
				dbESPortsGame.GetGameStatusTime() > 0 &&
				betValue.GetCreateTime() >= dbESPortsGame.GetGameStatusTime() {
				//时段无效
				DealGameInvalid(betValue, uniqueGameId, for_game.GAME_GUESS_BET_DISABLE_CODE_2, "比赛结束后投注", playerMapMsgs)
				continue
			}

			//比赛取消结束的时候 比赛异常无效
			if dbESPortsGame.GetGameStatus() == for_game.GAME_STATUS_2 &&
				dbESPortsGame.GetGameStatusType() == for_game.GAME_STATUS_TYPE_3 {
				//比赛异常无效
				DealGameInvalid(betValue, uniqueGameId, for_game.GAME_GUESS_BET_DISABLE_CODE_1, "比赛取消结束", playerMapMsgs)
				continue
			}

			//投注项判断
			var betNumObj *client_hall.ItemResult
			if betValue.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_1 {
				betNumObj = itemRstMornMap[betValue.GetBetNum()]
			} else {
				betNumObj = itemRstRollMap[betValue.GetBetNum()]
			}

			if nil != betNumObj {
				//比赛无结果 比赛异常无效
				if betNumObj.GetStatus() == for_game.GAME_GUESS_ITEM_STATUS_0 {
					//比赛异常无效
					DealGameInvalid(betValue, uniqueGameId, for_game.GAME_GUESS_BET_DISABLE_CODE_1, "投注项无结果", playerMapMsgs)
					continue
				} else {
					//滚盘
					if betValue.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_2 {

						//盘口未封盘
						if betNumObj.GetStatusTime() == 0 {
							// 时段无效
							s := fmt.Sprintf("投注项未封盘")
							DealGameInvalid(betValue, uniqueGameId, for_game.GAME_GUESS_BET_DISABLE_CODE_2, s, playerMapMsgs)
							continue
						} else {
							//封盘截至提前时间段下注 时段无效
							if betValue.GetCreateTime() > (betNumObj.GetStatusTime() - blockedBeforeTime) {
								// 时段无效
								s := fmt.Sprintf("投注项封盘前%v秒的无效时间内投注", blockedBeforeTime)
								DealGameInvalid(betValue, uniqueGameId, for_game.GAME_GUESS_BET_DISABLE_CODE_2, s, playerMapMsgs)
								continue
							}
						}
					}

					//判断达成、未达成结果
					//达成
					if betNumObj.GetWin() == for_game.GAME_GUESS_ITEM_WIN_1 {
						//计算金额、整数返还给用户
						amountFloat, amountErr := strconv.ParseFloat(strconv.FormatInt(betValue.GetBetAmount(), 10), 64)
						if nil != amountErr {
							s := fmt.Sprintf("=======CSGO开奖投计算达成时、转换用户投注额出错=========,参数为==UniqueGameId:%v,MornRollGuessFlag:%v,BetNum:%v,BetAmount:%v",
								betValue.GetUniqueGameId(),
								betValue.GetMornRollGuessFlag(),
								betValue.GetBetNum(),
								betValue.GetBetAmount())
							logs.Error(s)
							for_game.WriteFile("csgo_sport_lottery.log", s)
							continue
						}
						oddsFloat, oddsErr := strconv.ParseFloat(betValue.GetOdds(), 64)
						if nil != oddsErr {
							s := fmt.Sprintf("=======CSGO开奖投计算达成时、转换用户投赔率出错=========,参数为==UniqueGameId:%v,MornRollGuessFlag:%v,BetNum:%v,Odds:%v",
								betValue.GetUniqueGameId(),
								betValue.GetMornRollGuessFlag(),
								betValue.GetBetNum(),
								betValue.GetOdds())
							logs.Error(s)
							for_game.WriteFile("csgo_sport_lottery.log", s)
							continue
						}
						settlementAmount := int64(amountFloat * oddsFloat)
						//成功
						DealGuessBetSuccess(betValue, uniqueGameId, settlementAmount, playerMapMsgs)
						continue
					} else if betNumObj.GetWin() == for_game.GAME_GUESS_ITEM_WIN_0 {
						//未达成
						//无需返回金额只要更新数据库和通知
						DealGuessBetFail(betValue, playerMapMsgs)
						continue
					} else if betNumObj.GetWin() == for_game.GAME_GUESS_ITEM_WIN_NORST {
						//比赛异常无效
						DealGameInvalid(betValue, uniqueGameId, for_game.GAME_GUESS_BET_DISABLE_CODE_1, "投注项无结果", playerMapMsgs)
						continue
					}
				}

			} else {
				s := fmt.Sprintf("=======csgo开奖投注项不存在=========,参数为==投注OrderId:%v,投注项bet_num:%v,比赛唯一uniqueGameId:%v,动态盘口唯一UniqueGameGuessId:%v",
					betValue.GetOrderId(), betValue.GetBetNum(), uniqueGameId, betValue.GetUniqueGameGuessId())
				logs.Error(s)
				for_game.WriteFile("csgo_sport_lottery.log", s)
				continue
			}
		}

		//批量发送消息给用户
		if nil != playerMapMsgs && len(playerMapMsgs) > 0 {
			for key, value := range playerMapMsgs {
				notifyRst := dal_common.PushGameMultipleOrderSysMsg(PServerInfoMgr, key, value)
				if notifyRst.GetCode() != for_game.C_OPT_SUCCESS {
					s := fmt.Sprintf("=======csgo开奖批量发送给用户消息失败,直接返回=======")
					for_game.WriteFile("csgo_sport_lottery.log", s)
					return
				}
			}
		}
	} else {
		s := fmt.Sprintf("=======csgo开奖未取得未结算的投注记录的数据、清key和设置开奖状态后,直接返回=======,查询条件为uniqueGameId:%v,BetResult:%v",
			uniqueGameId, for_game.GAME_GUESS_BET_RESULT_1)
		for_game.WriteFile("csgo_sport_lottery.log", s)
		return
	}

	endStr := fmt.Sprintf("======DealESportCSGOLottery中比赛uniqueGameId:%v的csgo开奖正常  结束==========", uniqueGameId)
	for_game.WriteFile("csgo_sport_lottery.log", endStr)

}

func DealGameLottery(uniqueGameId int64) {
	startStr := fmt.Sprintf("======DealGameLottery中设置比赛uniqueGameId:%v为开奖状态  开始==========", uniqueGameId)
	for_game.WriteFile("csgo_sport_lottery.log", startStr)

	//第一步:不管有没有发布设置比赛已开奖
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
	defer closeFun()
	errUpd := col.Update(bson.M{"_id": uniqueGameId},
		bson.M{"$set": bson.M{"is_lottery": for_game.GAME_IS_LOTTERY_1, "update_time": time.Now().Unix()}})

	if errUpd != nil {
		s := fmt.Sprintf("========开奖设置比赛开奖状态时候比赛数据出错,errUpd错误:%v========更新条件为_id:%v", errUpd, uniqueGameId)
		logs.Error(s)
		for_game.WriteFile("csgo_sport_lottery.log", s)
		return
	}

	endStr := fmt.Sprintf("======DealGameLottery中设置比赛uniqueGameId:%v为开奖状态  结束", uniqueGameId)
	for_game.WriteFile("csgo_sport_lottery.log", endStr)
}

//开奖结算
func LotterySettlement(playId int64, coins int64, betOrderId int64, lotteryFlag int32) *client_hall.ESportCommonResult {
	rd := &client_hall.ESportCommonResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}
	//====返回电竞币开始====================
	st := for_game.ESPORTCOIN_TYPE_GUESS_BACK_IN
	var msg string
	if lotteryFlag == for_game.LOTTERY_FLAG_1 {
		msg = fmt.Sprintf("开奖无效返还电竞币[%d]个", coins)
	} else if lotteryFlag == for_game.LOTTERY_FLAG_2 {
		msg = fmt.Sprintf("开奖成功返还电竞币[%d]个", coins)
	}

	//取得流水订单号
	streamOrderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, st)
	req := &share_message.ESportCoinRecharge{
		PlayerId:     easygo.NewInt64(playId),
		RechargeCoin: easygo.NewInt64(coins),
		SourceType:   easygo.NewInt32(st),
		Note:         easygo.NewString(msg),
		ExtendLog: &share_message.GoldExtendLog{
			OrderId:    easygo.NewString(streamOrderId), //流水的订单号
			MerchantId: easygo.NewString(betOrderId),    //这里设置电竞的订单号
		},
	}

	result, err := dal_common.SendMsgToServerNewEx(PServerInfoMgr, playId, "RpcESportSendChangeESportCoins", req) //通知大厅
	if err != nil {
		logs.Error(err.GetReason())
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString(err.GetReason())
		return rd
	}
	if nil != result {
		rst, ok := result.(*client_hall.ESportCommonResult)
		if ok && nil != rst {
			if rst.GetCode() == for_game.C_SETTLEMENT_MONEY_FAIL {
				rd.Code = easygo.NewInt32(for_game.C_SETTLEMENT_MONEY_FAIL)
				rd.Msg = easygo.NewString(rst.GetMsg())
				return rd
			}
		}
	}

	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("结算成功")

	return rd
	//====返回电竞币结束================================
}

//比赛异常、时段无效处理
//flag 1:比赛异常无效 2:时段无效
func DealGameInvalid(betValue *share_message.TableESPortsGuessBetRecord,
	uniqueGameId int64, flag int32, reasonDetail string,
	playerMapMsgs map[int64][]*share_message.TableESPortsGameOrderSysMsg) *client_hall.ESportCommonResult {

	rd := &client_hall.ESportCommonResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}
	//先更新数据库
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD)
	defer closeFun()

	betValue.DisableAmount = easygo.NewInt64(betValue.GetBetAmount())
	if for_game.GAME_GUESS_BET_DISABLE_CODE_1 == flag {
		betValue.Reason = easygo.NewString(for_game.GAME_GUESS_BET_DISABLE_REASON_1)
	} else {
		betValue.Reason = easygo.NewString(for_game.GAME_GUESS_BET_DISABLE_REASON_2)
	}

	betValue.ReasonDetail = easygo.NewString(reasonDetail)
	betValue.BetStatus = easygo.NewString(for_game.GAME_GUESS_BET_STATUS_3)
	betValue.BetResult = easygo.NewString(for_game.GAME_GUESS_BET_RESULT_4)
	betValue.UpdateTime = easygo.NewInt64(time.Now().Unix())

	errUpd := col.Update(bson.M{"_id": betValue.GetOrderId(), "BetResult": for_game.GAME_GUESS_BET_RESULT_1},
		bson.M{"$set": betValue})

	if errUpd != nil && errUpd != mgo.ErrNotFound {
		s := fmt.Sprintf("=====CSGO开奖设置===投注无效数据出错===,errUpd错误:%v==更新条件为OrderId:%v,BetResult:%v",
			errUpd, betValue.GetOrderId(), for_game.GAME_GUESS_BET_RESULT_1)
		logs.Error(s)
		for_game.WriteFile("csgo_sport_lottery.log", s)

		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("系统异常")
		return rd
	}

	if errUpd == mgo.ErrNotFound {
		s := fmt.Sprintf("======CSGO开奖设置===投注无效数据没有数据========更新条件为OrderId:%v,BetResult:%v",
			betValue.GetOrderId(), for_game.GAME_GUESS_BET_RESULT_1)

		for_game.WriteFile("csgo_sport_lottery.log", s)

		rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
		rd.Msg = easygo.NewString("信息不存在")
		return rd
	}

	//再返回无效金额
	settlement := LotterySettlement(betValue.GetPlayInfo().GetPlayId(), betValue.GetDisableAmount(), betValue.GetOrderId(), for_game.LOTTERY_FLAG_1)
	if nil != settlement && settlement.GetCode() != for_game.C_OPT_SUCCESS {
		var s string
		if for_game.GAME_GUESS_BET_DISABLE_CODE_1 == flag {
			s = fmt.Sprintf("csgo开奖返回比赛异常无效金额失败,==参数==投注OrderId:%v,投注项bet_num:%v,返回用户playId:%v,返回比赛异常无效金额amount:%v,比赛唯一uniqueGameId:%v,动态盘口唯一UniqueGameGuessId:%v",
				betValue.GetOrderId(), betValue.GetBetNum(), betValue.GetPlayInfo().GetPlayId(), betValue.GetDisableAmount(), uniqueGameId, betValue.GetUniqueGameGuessId())
		} else {
			s = fmt.Sprintf("csgo开奖返回时段无效金额失败,==参数==投注OrderId:%v,投注项bet_num:%v,返回用户playId:%v,返回时段无效金额amount:%v,比赛唯一uniqueGameId:%v,动态盘口唯一UniqueGameGuessId:%v",
				betValue.GetOrderId(), betValue.GetBetNum(), betValue.GetPlayInfo().GetPlayId(), betValue.GetDisableAmount(), uniqueGameId, betValue.GetUniqueGameGuessId())
		}
		logs.Error(s)
		for_game.WriteFile("csgo_sport_lottery.log", s)
	}

	//通知用户结构
	eSPortsGameOrderSysMsg := &share_message.TableESPortsGameOrderSysMsg{
		OrderId:      easygo.NewInt64(betValue.GetOrderId()),
		UniqueGameId: easygo.NewInt64(betValue.GetUniqueGameId()),
		BetTime:      easygo.NewInt64(betValue.GetCreateTime()),
		Odds:         easygo.NewString(betValue.GetOdds()),
		BetResult:    easygo.NewString(betValue.GetBetResult()),
		BetTitle:     easygo.NewString(betValue.GetBetTitle()),
		BetNum:       easygo.NewString(betValue.GetBetNum()),
		BetName:      easygo.NewString(betValue.GetBetName()),
		ResultAmount: easygo.NewInt64(betValue.GetDisableAmount()),
		PlayerId:     easygo.NewInt64(betValue.GetPlayInfo().GetPlayId()),
		BetAmount:    easygo.NewInt64(betValue.GetBetAmount()),
	}
	//重新设置比赛名称
	if nil != betValue.GetGameInfo() {
		eSPortsGameOrderSysMsg.GameName = easygo.NewString(betValue.GetGameInfo().GetGameName() + " " +
			betValue.GetGameInfo().GetTeamAName() +
			" VS " +
			betValue.GetGameInfo().GetTeamBName())
	}

	//该用户的消息
	playerMsgs := playerMapMsgs[eSPortsGameOrderSysMsg.GetPlayerId()]
	if playerMsgs != nil && len(playerMsgs) > 0 {
		playerMsgs = append(playerMsgs, eSPortsGameOrderSysMsg)
		playerMapMsgs[eSPortsGameOrderSysMsg.GetPlayerId()] = playerMsgs
	} else {
		tmpMsg := make([]*share_message.TableESPortsGameOrderSysMsg, 0)
		tmpMsg = append(tmpMsg, eSPortsGameOrderSysMsg)
		playerMapMsgs[eSPortsGameOrderSysMsg.GetPlayerId()] = tmpMsg
	}

	//notifyRst := dal_common.PushGameOrderSysMsg(PServerInfoMgr, eSPortsGameOrderSysMsg)
	//if notifyRst.GetCode() != for_game.C_OPT_SUCCESS {
	//	rd.Code = easygo.NewInt32(notifyRst.GetCode())
	//	rd.Msg = easygo.NewString(notifyRst.GetMsg())
	//	return rd
	//}

	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("")
	return rd
}

//投注项达成
func DealGuessBetSuccess(betValue *share_message.TableESPortsGuessBetRecord,
	uniqueGameId int64, successAmount int64, playerMapMsgs map[int64][]*share_message.TableESPortsGameOrderSysMsg) *client_hall.ESportCommonResult {
	rd := &client_hall.ESportCommonResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}
	//先更新数据库
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD)
	defer closeFun()

	betValue.SuccessAmount = easygo.NewInt64(successAmount)

	betValue.BetStatus = easygo.NewString(for_game.GAME_GUESS_BET_STATUS_2)
	betValue.BetResult = easygo.NewString(for_game.GAME_GUESS_BET_RESULT_2)
	betValue.UpdateTime = easygo.NewInt64(time.Now().Unix())

	errUpd := col.Update(bson.M{"_id": betValue.GetOrderId(), "BetResult": for_game.GAME_GUESS_BET_RESULT_1},
		bson.M{"$set": betValue})

	if errUpd != nil && errUpd != mgo.ErrNotFound {
		s := fmt.Sprintf("CSGO开奖设置===投注成功数据出错===,errUpd错误:%v=====更新条件为OrderId:%v,BetResult:%v",
			errUpd, betValue.GetOrderId(), for_game.GAME_GUESS_BET_RESULT_1)
		logs.Error(s)
		for_game.WriteFile("csgo_sport_lottery.log", s)

		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("系统异常")
		return rd
	}

	if errUpd == mgo.ErrNotFound {
		s := fmt.Sprintf("CSGO开奖设置===投注成功没有数据========更新条件为OrderId:%v,BetResult:%v",
			betValue.GetOrderId(), for_game.GAME_GUESS_BET_RESULT_1)
		for_game.WriteFile("csgo_sport_lottery.log", s)

		rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
		rd.Msg = easygo.NewString("信息不存在")
		return rd
	}

	//再返回成功金额
	settlement := LotterySettlement(betValue.GetPlayInfo().GetPlayId(), betValue.GetSuccessAmount(), betValue.GetOrderId(), for_game.LOTTERY_FLAG_2)
	if nil != settlement && settlement.GetCode() != for_game.C_OPT_SUCCESS {
		s := fmt.Sprintf("csgo开奖返回成功金额失败,==参数==投注OrderId:%v,投注项bet_num:%v,返回用户playId:%v,返回成功金额settlementAmount:%v,比赛唯一uniqueGameId:%v,动态盘口唯一UniqueGameGuessId:%v",
			betValue.GetOrderId(), betValue.GetBetNum(), betValue.GetPlayInfo().GetPlayId(), betValue.GetSuccessAmount(), uniqueGameId, betValue.GetUniqueGameGuessId())
		logs.Error(s)
		for_game.WriteFile("csgo_sport_lottery.log", s)
	}

	//通知用户结构
	eSPortsGameOrderSysMsg := &share_message.TableESPortsGameOrderSysMsg{
		OrderId:      easygo.NewInt64(betValue.GetOrderId()),
		UniqueGameId: easygo.NewInt64(betValue.GetUniqueGameId()),
		BetTime:      easygo.NewInt64(betValue.GetCreateTime()),
		Odds:         easygo.NewString(betValue.GetOdds()),
		BetResult:    easygo.NewString(betValue.GetBetResult()),
		BetTitle:     easygo.NewString(betValue.GetBetTitle()),
		BetNum:       easygo.NewString(betValue.GetBetNum()),
		BetName:      easygo.NewString(betValue.GetBetName()),
		ResultAmount: easygo.NewInt64(betValue.GetSuccessAmount()),
		PlayerId:     easygo.NewInt64(betValue.GetPlayInfo().GetPlayId()),
		BetAmount:    easygo.NewInt64(betValue.GetBetAmount()),
	}
	//重新设置比赛名称
	if nil != betValue.GetGameInfo() {
		eSPortsGameOrderSysMsg.GameName = easygo.NewString(betValue.GetGameInfo().GetGameName() + " " +
			betValue.GetGameInfo().GetTeamAName() +
			" VS " +
			betValue.GetGameInfo().GetTeamBName())
	}

	//该用户的消息
	playerMsgs := playerMapMsgs[eSPortsGameOrderSysMsg.GetPlayerId()]
	if playerMsgs != nil && len(playerMsgs) > 0 {
		playerMsgs = append(playerMsgs, eSPortsGameOrderSysMsg)
		playerMapMsgs[eSPortsGameOrderSysMsg.GetPlayerId()] = playerMsgs
	} else {
		tmpMsg := make([]*share_message.TableESPortsGameOrderSysMsg, 0)
		tmpMsg = append(tmpMsg, eSPortsGameOrderSysMsg)
		playerMapMsgs[eSPortsGameOrderSysMsg.GetPlayerId()] = tmpMsg
	}

	//notifyRst := dal_common.PushGameOrderSysMsg(PServerInfoMgr, eSPortsGameOrderSysMsg)
	//if notifyRst.GetCode() != for_game.C_OPT_SUCCESS {
	//	rd.Code = easygo.NewInt32(notifyRst.GetCode())
	//	rd.Msg = easygo.NewString(notifyRst.GetMsg())
	//	return rd
	//}

	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("")
	return rd
}

//投注项失败
func DealGuessBetFail(betValue *share_message.TableESPortsGuessBetRecord, playerMapMsgs map[int64][]*share_message.TableESPortsGameOrderSysMsg) *client_hall.ESportCommonResult {
	rd := &client_hall.ESportCommonResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}
	//先更新数据库
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD)
	defer closeFun()

	betValue.FailAmount = easygo.NewInt64(betValue.GetBetAmount())

	betValue.BetStatus = easygo.NewString(for_game.GAME_GUESS_BET_STATUS_2)
	betValue.BetResult = easygo.NewString(for_game.GAME_GUESS_BET_RESULT_3)
	betValue.UpdateTime = easygo.NewInt64(time.Now().Unix())

	errUpd := col.Update(bson.M{"_id": betValue.GetOrderId(), "BetResult": for_game.GAME_GUESS_BET_RESULT_1},
		bson.M{"$set": betValue})

	if errUpd != nil && errUpd != mgo.ErrNotFound {
		s := fmt.Sprintf("CSGO开奖设置===投注失败数据出错===,errUpd错误:%v====更新条件为OrderId:%v,BetResult:%v",
			errUpd, betValue.GetOrderId(), for_game.GAME_GUESS_BET_RESULT_1)
		logs.Error(s)
		for_game.WriteFile("csgo_sport_lottery.log", s)

		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("系统异常")
		return rd
	}

	if errUpd == mgo.ErrNotFound {
		s := fmt.Sprintf("CSGO开奖设置===投注失败没有数据======更新条件为OrderId:%v,BetResult:%v",
			betValue.GetOrderId(), for_game.GAME_GUESS_BET_RESULT_1)
		for_game.WriteFile("csgo_sport_lottery.log", s)

		rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
		rd.Msg = easygo.NewString("信息不存在")
		return rd
	}

	//通知用户
	eSPortsGameOrderSysMsg := &share_message.TableESPortsGameOrderSysMsg{
		OrderId:      easygo.NewInt64(betValue.GetOrderId()),
		UniqueGameId: easygo.NewInt64(betValue.GetUniqueGameId()),
		BetTime:      easygo.NewInt64(betValue.GetCreateTime()),
		Odds:         easygo.NewString(betValue.GetOdds()),
		BetResult:    easygo.NewString(betValue.GetBetResult()),
		BetTitle:     easygo.NewString(betValue.GetBetTitle()),
		BetNum:       easygo.NewString(betValue.GetBetNum()),
		BetName:      easygo.NewString(betValue.GetBetName()),
		ResultAmount: easygo.NewInt64(0),
		PlayerId:     easygo.NewInt64(betValue.GetPlayInfo().GetPlayId()),
		BetAmount:    easygo.NewInt64(betValue.GetBetAmount()),
	}
	//重新设置比赛名称
	if nil != betValue.GetGameInfo() {
		eSPortsGameOrderSysMsg.GameName = easygo.NewString(betValue.GetGameInfo().GetGameName() + " " +
			betValue.GetGameInfo().GetTeamAName() +
			" VS " +
			betValue.GetGameInfo().GetTeamBName())
	}

	//notifyRst := dal_common.PushGameOrderSysMsg(PServerInfoMgr, eSPortsGameOrderSysMsg)
	//if notifyRst.GetCode() != for_game.C_OPT_SUCCESS {
	//	rd.Code = easygo.NewInt32(notifyRst.GetCode())
	//	rd.Msg = easygo.NewString(notifyRst.GetMsg())
	//	return rd
	//}

	//该用户的消息
	playerMsgs := playerMapMsgs[eSPortsGameOrderSysMsg.GetPlayerId()]
	if playerMsgs != nil && len(playerMsgs) > 0 {
		playerMsgs = append(playerMsgs, eSPortsGameOrderSysMsg)
		playerMapMsgs[eSPortsGameOrderSysMsg.GetPlayerId()] = playerMsgs
	} else {
		tmpMsg := make([]*share_message.TableESPortsGameOrderSysMsg, 0)
		tmpMsg = append(tmpMsg, eSPortsGameOrderSysMsg)
		playerMapMsgs[eSPortsGameOrderSysMsg.GetPlayerId()] = tmpMsg
	}

	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("")
	return rd
}

//初始化处理开奖
func InitDealCSGOLottery() {

	startStr := fmt.Sprintf("========InitDealCSGOLottery初始化处理csgo 定时发送开奖  开始==========")
	for_game.WriteFile("csgo_sport_lottery.log", startStr)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
	defer closeFun()

	//取得未开奖的数据(不管发布没有发布都开奖、为了清redis key)
	sportGameLists := make([]*share_message.TableESPortsGame, 0)

	errQuery := col.Find(bson.M{"game_status": for_game.GAME_STATUS_2,
		"is_lottery":   for_game.GAME_IS_LOTTERY_0,
		"app_label_id": for_game.ESPORTS_LABEL_CSGO}).All(&sportGameLists)

	if errQuery != nil && errQuery != mgo.ErrNotFound {

		s := fmt.Sprintf("======CSGO初始化开奖查询步骤2失败=====errQuery错误:%v,查询条件为game_status:%v,===is_lottery:%v,===app_label_id:%v",
			errQuery, for_game.GAME_STATUS_2, for_game.GAME_IS_LOTTERY_0, for_game.ESPORTS_LABEL_CSGO)

		logs.Error(s)
		for_game.WriteFile("csgo_sport_lottery.log", s)

		easygo.PanicError(errQuery)
	}

	if errQuery == mgo.ErrNotFound || nil == sportGameLists || len(sportGameLists) <= 0 {

		s := fmt.Sprintf("======CSGO初始化开奖查询步骤2未查询到数据、直接返回=====,查询条件为game_status:%v,===is_lottery:%v,===app_label_id:%v",
			for_game.GAME_STATUS_2, for_game.GAME_IS_LOTTERY_0, for_game.ESPORTS_LABEL_CSGO)

		for_game.WriteFile("csgo_sport_lottery.log", s)
		return
	}
	if errQuery == nil && nil != sportGameLists && len(sportGameLists) > 0 {

		if nil != sportGameLists && len(sportGameLists) > 0 {
			//必须要这样循环
			for i := 0; i < len(sportGameLists); i++ {
				gameId := sportGameLists[i].GetId()

				//初始化的时候按照设置的推迟时间开奖(给运营时间来设置)
				if CSGOSysMsgTimeMgr.GetTimerById(gameId) != nil {
					CSGOSysMsgTimeMgr.DelTimerList(gameId)
				}

				triggerTime := time.Duration(delayLotteryTime) * time.Second
				timer := easygo.AfterFunc(triggerTime, func() {
					DealESportCSGOLottery(gameId)
				})
				CSGOSysMsgTimeMgr.AddTimerList(gameId, timer)
			}
		}
	}

	endStr := fmt.Sprintf("========InitDealCSGOLottery初始化处理csgo 定时发送开奖  结束==========")
	for_game.WriteFile("csgo_sport_lottery.log", endStr)
}
