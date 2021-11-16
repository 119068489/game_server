// 大厅服务器为[游戏客户端]提供的服务

package sport_apply

import (
	"fmt"
	dal_common "game_server/e-sports/sport_common_dal"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"sort"
	"strconv"
	"time"
)

//获取比赛列表
func (self *sfc) RpcESportGetGameList(common *base.Common, reqMsg *client_hall.ESportGameListRequest) *client_hall.ESportGameListResult {

	logs.Info("===api RpcESportGetGameList===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	rd := &client_hall.ESportGameListResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}
	page := int(reqMsg.GetPage())
	pageSize := int(reqMsg.GetPageSize())

	labelType := reqMsg.GetLabelType()
	labelId := reqMsg.GetLabelId()

	gameClass := reqMsg.GetGameClass()

	curPage := easygo.If(page > 1, page-1, 0).(int)

	var dbList []*share_message.TableESPortsGame
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
	defer closeFun()

	queryBson := bson.M{}
	//当日零点时间
	timeStart := for_game.GetGameTodayStartTime()
	//当日结束时间
	timeEnd := for_game.GetGameTodayEndTime()

	//现在时间
	timeNow := time.Now().Unix()

	var sort []string

	//全部游戏标签
	if labelType == for_game.ESPORTS_LABEL_TYPE_1 {

		if gameClass == client_hall.Game_Class_GAME_TODAY {

			queryBson = bson.M{"begin_time_int": bson.M{"$gte": timeStart, "$lte": timeEnd}}

			sort = make([]string, 0)
			sort = append(sort, "game_status")
			sort = append(sort, "begin_time_int")

		} else if gameClass == client_hall.Game_Class_GAME_BEFORE {

			queryBson = bson.M{"game_status": for_game.GAME_STATUS_0,
				"begin_time_int": bson.M{"$gt": timeNow}}

			sort = make([]string, 0)
			sort = append(sort, "begin_time_int")

		} else if gameClass == client_hall.Game_Class_GAME_ROLL {

			queryBson = bson.M{"game_status": for_game.GAME_STATUS_0,
				"begin_time_int": bson.M{"$lte": timeNow}}

			sort = make([]string, 0)
			sort = append(sort, "begin_time_int")

		} else if gameClass == client_hall.Game_Class_GAME_OVER {
			queryBson = bson.M{"game_status": for_game.GAME_STATUS_2}

			sort = make([]string, 0)
			sort = append(sort, "-begin_time_int")
		}

		queryBson["$and"] = []bson.M{{"app_label_id": bson.M{"$ne": for_game.ESPORTS_LABEL_DOTA2}},
			{"app_label_id": bson.M{"$ne": for_game.ESPORTS_LABEL_CSGO}},
			{"app_label_id": bson.M{"$ne": for_game.ESPORTS_LABEL_OTHER}}}

		//单个游戏标签 labelId
	} else if labelType == for_game.ESPORTS_LABEL_TYPE_3 {

		if gameClass == client_hall.Game_Class_GAME_TODAY {

			queryBson = bson.M{"app_label_id": labelId,
				"begin_time_int": bson.M{"$gte": timeStart, "$lte": timeEnd}}

			sort = make([]string, 0)
			sort = append(sort, "game_status")
			sort = append(sort, "begin_time_int")

		} else if gameClass == client_hall.Game_Class_GAME_BEFORE {
			queryBson = bson.M{"app_label_id": labelId,
				"game_status":    for_game.GAME_STATUS_0,
				"begin_time_int": bson.M{"$gt": timeNow}}

			sort = make([]string, 0)
			sort = append(sort, "begin_time_int")
		} else if gameClass == client_hall.Game_Class_GAME_ROLL {
			queryBson = bson.M{"app_label_id": labelId,
				"game_status":    for_game.GAME_STATUS_0,
				"begin_time_int": bson.M{"$lte": timeNow}}

			sort = make([]string, 0)
			sort = append(sort, "begin_time_int")
		} else if gameClass == client_hall.Game_Class_GAME_OVER {
			queryBson = bson.M{"app_label_id": labelId,
				"game_status": for_game.GAME_STATUS_2}

			sort = make([]string, 0)
			sort = append(sort, "-begin_time_int")
		}
	}

	queryBson["release_flag"] = for_game.GAME_RELEASE_FLAG_2
	queryBson["game_status_type"] = bson.M{"$ne": for_game.GAME_STATUS_TYPE_3}

	query := col.Find(queryBson)
	count, errCnt := query.Count()
	if errCnt != nil {
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("系统异常")
		return rd
	}

	errQuery := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&dbList)
	if errQuery != nil && errQuery != mgo.ErrNotFound {
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("系统异常")
		return rd
	}

	if errQuery == mgo.ErrNotFound || nil == dbList || len(dbList) <= 0 {
		rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
		rd.Msg = easygo.NewString("")
		rd.Total = easygo.NewInt32(0)
		rd.GameList = make([]*client_hall.ESportGameObject, 0)
		return rd
	}

	if errQuery == nil && dbList != nil && len(dbList) > 0 {
		rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
		rd.Msg = easygo.NewString("")
		rd.Total = easygo.NewInt32(count)
		//组织返回数据
		retGameList := make([]*client_hall.ESportGameObject, 0)
		for _, value := range dbList {
			retGameList = append(retGameList, &client_hall.ESportGameObject{
				UniqueGameId: easygo.NewInt64(value.GetId()),
				GameName:     easygo.NewString(value.GetMatchName() + value.GetMatchStage() + "-BO" + value.GetBo()),
				TeamAInfo: &client_hall.TeamObject{
					TeamId: easygo.NewString(value.GetTeamA().GetTeamId()),
					Name:   easygo.NewString(value.GetTeamA().GetName()),
					Icon:   easygo.NewString(value.GetTeamA().GetIcon()),
				},
				ScoreA: easygo.NewString(value.GetScoreA()),
				ScoreB: easygo.NewString(value.GetScoreB()),
				TeamBInfo: &client_hall.TeamObject{
					TeamId: easygo.NewString(value.GetTeamB().GetTeamId()),
					Name:   easygo.NewString(value.GetTeamB().GetName()),
					Icon:   easygo.NewString(value.GetTeamB().GetIcon()),
				},
				BeginTime:    easygo.NewInt64(value.GetBeginTimeInt()),
				BeginTimeStr: easygo.NewString(value.GetBeginTime()),
				GameStatus:   easygo.NewString(for_game.GetGameStatus(value.GetBeginTime(), value.GetGameStatus())),
				HaveRoll:     easygo.NewInt32(value.GetHaveRoll()),
				AppLabelId:   easygo.NewInt32(value.GetAppLabelId()),
				ApiOrigin:    easygo.NewInt32(value.GetApiOrigin()),
				GameId:       easygo.NewString(value.GetGameId()),
				HistoryId:    easygo.NewInt64(value.GetHistoryId()),
			})
		}

		rd.GameList = retGameList

		//比赛图标需要重配置表中取得后然后重新匹配
		gameLabelRedisObj := for_game.GetRedisGameLabel()
		if gameLabelRedisObj != nil {
			if nil != rd.GetGameList() && len(rd.GetGameList()) > 0 {
				for _, gameValue := range rd.GetGameList() {
					if gameValue.GetAppLabelId() == for_game.ESPORTS_LABEL_WZRY {
						gameValue.GameIcon = easygo.NewString(gameLabelRedisObj.GetWZRYIcon())
					} else if gameValue.GetAppLabelId() == for_game.ESPORTS_LABEL_DOTA2 {
						gameValue.GameIcon = easygo.NewString(gameLabelRedisObj.GetDOTAIcon())
					} else if gameValue.GetAppLabelId() == for_game.ESPORTS_LABEL_LOL {
						gameValue.GameIcon = easygo.NewString(gameLabelRedisObj.GetLOLIcon())
					} else if gameValue.GetAppLabelId() == for_game.ESPORTS_LABEL_CSGO {
						gameValue.GameIcon = easygo.NewString(gameLabelRedisObj.GetCSGOIcon())
					} else if gameValue.GetAppLabelId() == for_game.ESPORTS_LABEL_OTHER {
						gameValue.GameIcon = easygo.NewString(gameLabelRedisObj.GetOTHERIcon())
					}

				}
			}
		}
	}

	logs.Info("=========RpcESportGetGameList 返回===========")
	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("")
	return rd
}

//获取比赛详情数据
func (self *sfc) RpcESportGetGameDetail(common *base.Common, reqMsg *client_hall.GameDetailRequest) *client_hall.ESportGameDetailResult {

	//有轮询不要打log
	//logs.Info("===api RpcESportGetGameDetail===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	rd := &client_hall.ESportGameDetailResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}

	//设置Head的数据
	rd.GameDetailHead = for_game.GetRedisGameDetailHead(reqMsg.GetUniqueGameId())

	//从redis取得风控数据
	betRiskControl := for_game.GetRedisGuessBetRiskControl()
	if betRiskControl != nil {
		rd.MaxAmount = easygo.NewInt64(betRiskControl.GetEsOneBetGold())
	}

	//设置竞猜项目数据
	if rd.GetGameDetailHead() != nil {

		//计算重置比赛状态
		rd.GameDetailHead.GameStatus = easygo.NewString(
			for_game.GetGameStatus(rd.GameDetailHead.GetBeginTimeStr(),
				rd.GameDetailHead.GetGameStatus()))

		//非轮询和轮询的时候都设置有、后期看需求、先预留判断
		if reqMsg.GetPollFlag() == 0 {
			// 设置放映厅按钮先设置有
			rd.GameDetailHead.HaveVideoHall = easygo.NewInt32(for_game.GAME_HAVE_VIDEO_HALL_STATUS_1)
		} else {
			rd.GameDetailHead.HaveVideoHall = easygo.NewInt32(for_game.GAME_HAVE_VIDEO_HALL_STATUS_1)
		}

		//比赛未开始
		if rd.GetGameDetailHead().GetGameStatus() == for_game.GAME_STATUS_0 {

			//显示早盘
			rd.MornRollGuessFlag = easygo.NewInt32(for_game.GAME_IS_MORN_ROLL_1)
			//查询动态盘口
			guessDetail := for_game.GetRedisGuessDetail(reqMsg.GetUniqueGameId(),
				rd.GetGameDetailHead().GetAppLabelId(),
				rd.GetGameDetailHead().GetGameId(),
				rd.GetGameDetailHead().GetApiOrigin(),
				for_game.GAME_IS_MORN_ROLL_1)
			if nil != guessDetail {
				rd.UniqueGameGuessId = easygo.NewInt64(guessDetail.GetUniqueGameGuessId())
				rd.GuessOddsNums = guessDetail.GetGuessOddsNums()
			}

			return rd
			//进行中
		} else if rd.GetGameDetailHead().GetGameStatus() == for_game.GAME_STATUS_1 {

			//查询滚盘动态盘口
			guessDetail := for_game.GetRedisGuessDetail(reqMsg.GetUniqueGameId(),
				rd.GetGameDetailHead().GetAppLabelId(),
				rd.GetGameDetailHead().GetGameId(),
				rd.GetGameDetailHead().GetApiOrigin(),
				for_game.GAME_IS_MORN_ROLL_2)
			//有滚盘显示滚盘、没有显示早盘
			if nil != guessDetail {
				rd.MornRollGuessFlag = easygo.NewInt32(for_game.GAME_IS_MORN_ROLL_2)
				rd.UniqueGameGuessId = easygo.NewInt64(guessDetail.GetUniqueGameGuessId())
				rd.GuessOddsNums = guessDetail.GetGuessOddsNums()

				return rd
			} else {
				//查询早盘盘口
				guessDetail = for_game.GetRedisGuessDetail(reqMsg.GetUniqueGameId(), rd.GetGameDetailHead().GetAppLabelId(),
					rd.GetGameDetailHead().GetGameId(),
					rd.GetGameDetailHead().GetApiOrigin(),
					for_game.GAME_IS_MORN_ROLL_1)
				if nil != guessDetail {
					rd.MornRollGuessFlag = easygo.NewInt32(for_game.GAME_IS_MORN_ROLL_1)
					rd.UniqueGameGuessId = easygo.NewInt64(guessDetail.GetUniqueGameGuessId())
					rd.GuessOddsNums = guessDetail.GetGuessOddsNums()

					//比赛进行中如果是早盘的时候、全封、再次给投注状态设置
					if nil != rd.GetGuessOddsNums() && len(rd.GetGuessOddsNums()) > 0 {
						for _, value := range rd.GetGuessOddsNums() {
							if value.GetContents() != nil && len(value.GetContents()) > 0 {
								for _, cntValue := range value.GetContents() {
									if nil != cntValue.GetItems() && len(cntValue.GetItems()) > 0 {
										for _, itemValue := range cntValue.GetItems() {
											itemValue.BetStatus = easygo.NewString(for_game.GAME_GUESS_ITEM_ODDS_STATUS_2)
										}
									}
								}
							}
						}
					}

					return rd
				}
			}

			//已结束
		} else if rd.GetGameDetailHead().GetGameStatus() == for_game.GAME_STATUS_2 {

			//查询滚盘动态盘口
			guessDetailRoll := for_game.GetRedisGuessDetail(reqMsg.GetUniqueGameId(), rd.GetGameDetailHead().GetAppLabelId(),
				rd.GetGameDetailHead().GetGameId(),
				rd.GetGameDetailHead().GetApiOrigin(),
				for_game.GAME_IS_MORN_ROLL_2)
			//查询早盘动态盘口
			guessDetailMorn := for_game.GetRedisGuessDetail(reqMsg.GetUniqueGameId(), rd.GetGameDetailHead().GetAppLabelId(),
				rd.GetGameDetailHead().GetGameId(),
				rd.GetGameDetailHead().GetApiOrigin(),
				for_game.GAME_IS_MORN_ROLL_1)

			if guessDetailMorn != nil && guessDetailRoll == nil {
				rd.MornRollGuessFlag = easygo.NewInt32(for_game.GAME_IS_MORN_ROLL_1)
				rd.UniqueGameGuessId = easygo.NewInt64(guessDetailMorn.GetUniqueGameGuessId())
				rd.GuessOddsNums = guessDetailMorn.GetGuessOddsNums()
			} else if guessDetailMorn == nil && guessDetailRoll != nil {
				rd.MornRollGuessFlag = easygo.NewInt32(for_game.GAME_IS_MORN_ROLL_2)
				rd.UniqueGameGuessId = easygo.NewInt64(guessDetailRoll.GetUniqueGameGuessId())
				rd.GuessOddsNums = guessDetailRoll.GetGuessOddsNums()
			} else if guessDetailMorn != nil && guessDetailRoll != nil {

				rollGuessOddsNums := guessDetailRoll.GetGuessOddsNums()
				mornGuessOddsNums := guessDetailMorn.GetGuessOddsNums()

				if mornGuessOddsNums != nil &&
					len(mornGuessOddsNums) > 0 &&
					(rollGuessOddsNums == nil || len(rollGuessOddsNums) <= 0) {
					rd.MornRollGuessFlag = easygo.NewInt32(for_game.GAME_IS_MORN_ROLL_1)
					rd.UniqueGameGuessId = easygo.NewInt64(guessDetailMorn.GetUniqueGameGuessId())
					rd.GuessOddsNums = mornGuessOddsNums

				} else if rollGuessOddsNums != nil &&
					len(rollGuessOddsNums) > 0 &&
					(mornGuessOddsNums == nil || len(mornGuessOddsNums) <= 0) {
					rd.MornRollGuessFlag = easygo.NewInt32(for_game.GAME_IS_MORN_ROLL_2)
					rd.UniqueGameGuessId = easygo.NewInt64(guessDetailRoll.GetUniqueGameGuessId())
					rd.GuessOddsNums = rollGuessOddsNums
				} else if rollGuessOddsNums != nil &&
					len(rollGuessOddsNums) > 0 &&
					(mornGuessOddsNums != nil && len(mornGuessOddsNums) > 0) {
					//将早盘和滚盘结合组合后返回给前端
					//以滚盘为基准、过滤重复项
					//key 第几局数、value:给予临时的一个值easygo.NewString("1")
					guessNumMaps := make(map[string]*string)
					//key betId、value:给予临时的一个值easygo.NewString("1")
					betContentMaps := make(map[string]*string)
					//key betNum、value:给予临时的一个值easygo.NewString("1")
					betNumMaps := make(map[string]*string)

					//将滚盘的所有投注项数据加入到Map中
					for _, rollNumValue := range rollGuessOddsNums {
						guessNumMaps[rollNumValue.GetNum()] = easygo.NewString("1")
						contents := rollNumValue.GetContents()
						if nil != contents && len(contents) > 0 {
							for _, cntValue := range contents {
								betContentMaps[cntValue.GetBetId()] = easygo.NewString("1")
								items := cntValue.GetItems()
								if nil != items && len(items) > 0 {
									for _, itemValue := range items {
										betNumMaps[itemValue.GetBetNum()] = easygo.NewString("1")
									}
								}
							}
						}
					}

					//把每个早盘的投注项目对象去查询滚盘中有没有、没有就加在滚盘的对象中、已经存在就忽略不加
					for _, mornNumValue := range mornGuessOddsNums {
						//每个投注项比较滚盘的数据
						tempGuessNumMapValue := guessNumMaps[mornNumValue.GetNum()]
						if tempGuessNumMapValue == nil || *tempGuessNumMapValue != "1" {
							// 加入到滚盘结构中
							guessNumMaps[mornNumValue.GetNum()] = easygo.NewString("1")
							rollGuessOddsNums = append(rollGuessOddsNums, mornNumValue)
						} else {
							mornContents := mornNumValue.GetContents()
							if nil != mornContents && len(mornContents) > 0 {
								for _, mornCntValue := range mornContents {
									tempBetContentValue := betContentMaps[mornCntValue.GetBetId()]
									if tempBetContentValue == nil || *tempBetContentValue != "1" {
										// 加入到滚盘结构中
										betContentMaps[mornCntValue.GetBetId()] = easygo.NewString("1")
										for _, tempValue := range rollGuessOddsNums {

											if mornNumValue.GetNum() == tempValue.GetNum() {
												tempCnts := tempValue.GetContents()
												tempCnts = append(tempCnts, mornCntValue)
												break
											}

										}
									} else {

										mornItems := mornCntValue.GetItems()

										if nil != mornItems && len(mornItems) > 0 {
											for _, mornItemValue := range mornItems {
												tempBetNumValue := betNumMaps[mornItemValue.GetBetNum()]
												if tempBetNumValue == nil || *tempBetNumValue != "1" {
													betNumMaps[mornItemValue.GetBetNum()] = easygo.NewString("1")

													for _, tempValue := range rollGuessOddsNums {

														if mornNumValue.GetNum() == tempValue.GetNum() {
															tempCnts := tempValue.GetContents()
															if nil != tempCnts && len(tempCnts) > 0 {
																for _, tempCntValue := range tempCnts {
																	if mornCntValue.GetBetId() == tempCntValue.GetBetId() {
																		tempItems := tempCntValue.GetItems()
																		tempItems = append(tempItems, mornItemValue)
																		break
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
					rd.GuessOddsNums = rollGuessOddsNums

					//重新给组合过的投注内容排序
					guessOddsNumsSort := GuessOddsNumsSort{}
					if nil != rd.GetGuessOddsNums() && len(rd.GetGuessOddsNums()) > 0 {
						for _, value := range rd.GetGuessOddsNums() {
							guessOddsNumsSort = append(guessOddsNumsSort, value)
						}

						sort.Sort(guessOddsNumsSort)

						//重新设置返回值得
						list := make([]*share_message.GameGuessOddsNumObject, 0)
						for _, rdValue := range guessOddsNumsSort {
							list = append(list, rdValue)
						}
						rd.GuessOddsNums = list
					}
				}
			}

			//比赛结束了全封、再次给投注状态设置
			if nil != rd.GetGuessOddsNums() && len(rd.GetGuessOddsNums()) > 0 {
				for _, value := range rd.GetGuessOddsNums() {
					if value.GetContents() != nil && len(value.GetContents()) > 0 {
						for _, cntValue := range value.GetContents() {
							if nil != cntValue.GetItems() && len(cntValue.GetItems()) > 0 {
								for _, itemValue := range cntValue.GetItems() {
									itemValue.BetStatus = easygo.NewString(for_game.GAME_GUESS_ITEM_ODDS_STATUS_2)
								}
							}
						}
					}
				}
			}
		}
	}
	//有轮询不要打log
	//logs.Info("=========RpcESportGetGameDetail 返回===========", rd)
	return rd
}

//比赛投注
func (self *sfc) RpcESportGameGuessBet(common *base.Common, reqMsg *client_hall.GameGuessBetRequest) *client_hall.GameGuessBetResult {

	logs.Info("===api RpcESportGameGuessBet===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	rd := &client_hall.GameGuessBetResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}

	// 风控数据取得
	//从redis取得风控数据
	betRiskControl := for_game.GetRedisGuessBetRiskControl()

	//生成投注订单
	//组装数据库对象
	var coins int32
	var orderIds string
	nowTime := time.Now().Unix()
	bets := make([]*share_message.TableESPortsGuessBetRecord, 0)
	for _, value := range reqMsg.GetGuessBets() {
		//对传输过来的数据作判断、防止前端传输错误以及被别人盗刷
		//取得redis中的gameHead信息
		gameObject := for_game.GetRedisGameDetailHead(value.GetUniqueGameId())
		if nil != gameObject {
			//取得某对应投注项的早盘或者滚盘的redis信息
			var guessDetailObject *share_message.GameGuessDetailObject
			if value.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_1 {
				guessDetailObject = for_game.GetRedisGuessMornDetail(value.GetUniqueGameId(),
					gameObject.GetAppLabelId(),
					gameObject.GetGameId(), gameObject.GetApiOrigin())
			} else {
				guessDetailObject = for_game.GetRedisGuessRollDetail(value.GetUniqueGameId(),
					gameObject.GetAppLabelId(),
					gameObject.GetGameId(), gameObject.GetApiOrigin())
			}

			if nil != guessDetailObject {

				//比赛结束
				if guessDetailObject.GetGameStatus() == for_game.GAME_STATUS_2 {
					rd.Code = easygo.NewInt32(for_game.C_INFO_ERROR)
					s := fmt.Sprintf("存在投注项的比赛已经结束、刷新重试")
					rd.Msg = easygo.NewString(s)
					return rd
				}

				//早盘已经开始
				if for_game.GetGameStatus(guessDetailObject.GetBeginTime(), guessDetailObject.GetGameStatus()) == for_game.GAME_STATUS_1 &&
					value.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_1 {
					rd.Code = easygo.NewInt32(for_game.C_INFO_ERROR)
					s := fmt.Sprintf("存在早盘投注项的比赛已经开始、刷新重试")
					rd.Msg = easygo.NewString(s)
					return rd
					//滚盘未开始
				} else if for_game.GetGameStatus(guessDetailObject.GetBeginTime(), guessDetailObject.GetGameStatus()) == for_game.GAME_STATUS_0 &&
					value.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_2 {
					rd.Code = easygo.NewInt32(for_game.C_INFO_ERROR)
					s := fmt.Sprintf("存在滚盘投注项的比赛未开始、刷新重试")
					rd.Msg = easygo.NewString(s)
					return rd
				}

				//判断断投注项
				if nil != guessDetailObject.GetGuessOddsNums() && len(guessDetailObject.GetGuessOddsNums()) > 0 {
					for _, guessValue := range guessDetailObject.GetGuessOddsNums() {
						if guessValue.GetContents() != nil && len(guessValue.GetContents()) > 0 {
							for _, contValue := range guessValue.GetContents() {
								if contValue.GetItems() != nil && len(contValue.GetItems()) > 0 {
									for _, itemValue := range contValue.GetItems() {
										if value.GetBetNum() == itemValue.GetBetNum() && itemValue.GetBetStatus() == for_game.GAME_GUESS_ITEM_ODDS_STATUS_2 {
											rd.Code = easygo.NewInt32(for_game.C_INFO_ERROR)
											s := fmt.Sprintf("存在投注项已经封盘、刷新重试")
											rd.Msg = easygo.NewString(s)
											return rd
										}
									}
								}
							}
						}
					}
				}
			}
		}

		//用户单次投注的限额判断
		if nil != betRiskControl {
			if value.GetBetAmount() > betRiskControl.GetEsOneBetGold() {
				rd.Code = easygo.NewInt32(for_game.C_OVER_QUOTA_ONE)
				s := fmt.Sprintf("单笔投注额度不能高于%v，超出限额的投注将自动调整至限额", betRiskControl.GetEsOneBetGold())
				rd.Msg = easygo.NewString(s)
				rd.MaxAmount = easygo.NewInt64(betRiskControl.GetEsOneBetGold())
				return rd
			}
		}

		coins = coins + int32(value.GetBetAmount())
		bet := &share_message.TableESPortsGuessBetRecord{
			OrderId:           easygo.NewInt64(for_game.RedisCreateBetOrderID()),
			UniqueGameId:      easygo.NewInt64(value.GetUniqueGameId()),
			UniqueGameGuessId: easygo.NewInt64(value.GetUniqueGameGuessId()),
			MornRollGuessFlag: easygo.NewInt32(value.GetMornRollGuessFlag()),
			GameInfo: &share_message.GuessBetGameInfo{
				GameName:  easygo.NewString(value.GetGameName()),
				TeamAName: easygo.NewString(value.GetTeamAName()),
				TeamBName: easygo.NewString(value.GetTeamBName()),
			},
			BetId:         easygo.NewString(value.BetId),
			BetTitle:      easygo.NewString(value.GetBetTitle()),
			BetNum:        easygo.NewString(value.GetBetNum()),
			BetName:       easygo.NewString(value.GetBetName()),
			Odds:          easygo.NewString(value.GetOdds()),
			BetAmount:     easygo.NewInt64(value.GetBetAmount()),
			SuccessAmount: easygo.NewInt64(0),
			FailAmount:    easygo.NewInt64(0),
			DisableAmount: easygo.NewInt64(0),
			IllegalAmount: easygo.NewInt64(0),
			Reason:        easygo.NewString(""),
			ReasonDetail:  easygo.NewString(""),
			BetStatus:     easygo.NewString(for_game.GAME_GUESS_BET_STATUS_1),
			BetResult:     easygo.NewString(for_game.GAME_GUESS_BET_RESULT_1),
			CreateTime:    easygo.NewInt64(nowTime),
			UpdateTime:    easygo.NewInt64(nowTime),
		}

		if orderIds == "" {
			orderIds = strconv.FormatInt(bet.GetOrderId(), 10)
		} else {
			orderIds = orderIds + "|" + strconv.FormatInt(bet.GetOrderId(), 10)
		}

		//设置部分比赛的数据
		gameInfo := for_game.GetRedisGameDetailHead(value.GetUniqueGameId())
		if nil != gameInfo {
			bet.AppLabelId = easygo.NewInt32(gameInfo.GetAppLabelId())
			bet.AppLabelName = easygo.NewString(for_game.LabelToESportNameMap[gameInfo.GetAppLabelId()])
			bet.ApiOrigin = easygo.NewInt32(gameInfo.GetApiOrigin())
			bet.ApiOriginName = easygo.NewString(for_game.ApiOriginIdToNameMap[gameInfo.GetApiOrigin()])
			bet.GameId = easygo.NewString(gameInfo.GetGameId())
		}

		//设置用户数据
		PlayInfo := share_message.GuessBetPlayerInfo{
			PlayId: easygo.NewInt64(common.GetUserId()),
		}
		player := for_game.GetRedisPlayerBase(common.GetUserId())

		if nil != player {
			PlayInfo.Account = easygo.NewString(player.GetAccount())
			PlayInfo.Phone = easygo.NewString(player.GetPhone())
		}

		bet.PlayInfo = &PlayInfo
		bets = append(bets, bet)
	}

	//风控用户当日投注额度判断
	betRiskOneDay := share_message.TableESPortsBetRiskOneDay{}
	colOneDay, closeFunOneDay := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_BET_RISK_ONE_DAY)
	defer closeFunOneDay()
	errQuery := colOneDay.Find(bson.M{"PlayerId": common.GetUserId(),
		"DateStr": util.GetYMD()}).One(&betRiskOneDay)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("系统异常")
		return rd
	}

	if errQuery == mgo.ErrNotFound {
		if betRiskControl != nil && int64(coins) > betRiskControl.GetEsOneDayBetGold() {
			rd.Code = easygo.NewInt32(for_game.C_OVER_QUOTA_ONE_DAY)
			rd.Msg = easygo.NewString("您当日可投注总量已达最大额度，无法继续投注")
			return rd
		}
		betRiskOneDay.PlayerId = easygo.NewInt64(common.GetUserId())
		betRiskOneDay.DateStr = easygo.NewString(util.GetYMD())
		betRiskOneDay.AmountDay = easygo.NewInt64(coins)
		betRiskOneDay.CreateTime = easygo.NewInt64(nowTime)
		betRiskOneDay.UpdateTime = easygo.NewInt64(nowTime)
	}

	if nil == errQuery {
		if betRiskControl != nil && ((int64(coins) + betRiskOneDay.GetAmountDay()) > betRiskControl.GetEsOneDayBetGold()) {
			rd.Code = easygo.NewInt32(for_game.C_OVER_QUOTA_ONE_DAY)
			rd.Msg = easygo.NewString("您当日可投注总量已达最大额度，无法继续投注")
			return rd
		}
		betRiskOneDay.AmountDay = easygo.NewInt64(int64(coins) + betRiskOneDay.GetAmountDay())
		betRiskOneDay.UpdateTime = easygo.NewInt64(time.Now().Unix())
	}

	//====先扣电竞币开始====================
	st := for_game.ESPORTCOIN_TYPE_GUESS_BET_OUT
	msg := fmt.Sprintf("竞猜投注电竞币[%d]个", coins)

	var tempCoins int32
	if coins > 0 {
		tempCoins = -coins
	} else {
		tempCoins = coins
	}

	//取得流水订单号
	streamOrderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, st)
	req := &share_message.ESportCoinRecharge{
		PlayerId:     easygo.NewInt64(common.GetUserId()),
		RechargeCoin: easygo.NewInt64(tempCoins),
		SourceType:   easygo.NewInt32(st),
		Note:         easygo.NewString(msg),
		ExtendLog: &share_message.GoldExtendLog{
			OrderId:    easygo.NewString(streamOrderId), //流水的订单号
			MerchantId: easygo.NewString(orderIds),      //这里设置电竞的订单号组合、中间用"|"分割
		},
	}

	result, err := dal_common.SendMsgToServerNewEx(PServerInfoMgr, common.GetUserId(), "RpcESportSendChangeESportCoins", req) //通知大厅
	if err != nil {
		logs.Error(err.GetReason())
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString(err.GetReason())
		return rd
	}
	if nil != result {
		rst, ok := result.(*client_hall.ESportCommonResult)
		if ok && nil != rst {
			if rst.GetCode() == for_game.C_DEDUCT_MONEY_FAIL {
				rd.Code = easygo.NewInt32(for_game.C_DEDUCT_MONEY_FAIL)
				rd.Msg = easygo.NewString(rst.GetMsg())
				return rd
			}
		}
	}
	//====先扣电竞币结束================================

	//批量插入数据库中
	var insLst []interface{}
	if nil != bets && len(bets) > 0 {
		for _, insValue := range bets {
			insLst = append(insLst, insValue)
		}
	}

	//批量插入
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD)
	defer closeFun()

	if nil != insLst && len(insLst) > 0 {
		errIns := col.Insert(insLst...)
		if errIns != nil {
			logs.Error(errIns)
			rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
			rd.Msg = easygo.NewString("系统异常")

			return rd
		}
	}

	//记录风控用户当日投注额度记录表
	_, errOneDay := colOneDay.Upsert(bson.M{"PlayerId": betRiskOneDay.GetPlayerId(),
		"DateStr": betRiskOneDay.GetDateStr()},
		bson.M{"$set": &betRiskOneDay})

	if nil != errOneDay {
		logs.Error(errOneDay)
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("系统异常")

		return rd
	}

	//记录风控平台当日投注额度记录表
	colDaySum, closeFunDaySum := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_BET_RISK_PLATFORM_DAY_SUM)
	defer closeFunDaySum()
	betRiskDaySum := share_message.TableESPortsBetRiskPlatFormDaySum{}
	errDaySum := colDaySum.Find(bson.M{"DateStr": util.GetYMD()}).One(&betRiskDaySum)

	if errDaySum != nil && errDaySum != mgo.ErrNotFound {
		logs.Error(errDaySum)
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("系统异常")
		return rd
	}
	if errDaySum == mgo.ErrNotFound {
		betRiskDaySum.DateStr = easygo.NewString(util.GetYMD())
		betRiskDaySum.AmountDaySum = easygo.NewInt64(coins)
		betRiskDaySum.CreateTime = easygo.NewInt64(nowTime)
		betRiskDaySum.UpdateTime = easygo.NewInt64(nowTime)

		_, errDaySumUpd := colDaySum.Upsert(bson.M{"DateStr": betRiskDaySum.GetDateStr()},
			bson.M{"$set": &betRiskDaySum})

		if nil != errDaySumUpd {
			logs.Error(errDaySumUpd)
			rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
			rd.Msg = easygo.NewString("系统异常")

			return rd
		}
	}

	if nil == errDaySum {
		//存在要在数据库中直接记录，防止并发
		_, errDaySumUpd := colDaySum.Upsert(bson.M{"DateStr": betRiskDaySum.GetDateStr()},
			bson.M{"$inc": bson.M{"AmountDaySum": coins}})

		if nil != errDaySumUpd {
			logs.Error(errDaySumUpd)
			rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
			rd.Msg = easygo.NewString("系统异常")

			return rd
		}
	}

	//组织通知用数据
	eSPortsGameOrderSysMsgs := make([]*share_message.TableESPortsGameOrderSysMsg, 0)
	for _, betValue := range bets {
		eSPortsGameOrderSysMsg := &share_message.TableESPortsGameOrderSysMsg{
			OrderId:      easygo.NewInt64(betValue.GetOrderId()),
			UniqueGameId: easygo.NewInt64(betValue.GetUniqueGameId()),
			BetTime:      easygo.NewInt64(betValue.GetCreateTime()),
			Odds:         easygo.NewString(betValue.GetOdds()),
			BetResult:    easygo.NewString(betValue.GetBetStatus()),
			BetTitle:     easygo.NewString(betValue.GetBetTitle()),
			BetNum:       easygo.NewString(betValue.GetBetNum()),
			BetName:      easygo.NewString(betValue.GetBetName()),
			ResultAmount: easygo.NewInt64(0),
			PlayerId:     easygo.NewInt64(common.GetUserId()),
			BetAmount:    easygo.NewInt64(betValue.GetBetAmount()),
		}
		//重新设置比赛名称
		if nil != betValue.GetGameInfo() {
			eSPortsGameOrderSysMsg.GameName = easygo.NewString(betValue.GetGameInfo().GetGameName() + " " +
				betValue.GetGameInfo().GetTeamAName() +
				" VS " +
				betValue.GetGameInfo().GetTeamBName())
		}

		eSPortsGameOrderSysMsgs = append(eSPortsGameOrderSysMsgs, eSPortsGameOrderSysMsg)
	}

	if nil != eSPortsGameOrderSysMsgs && len(eSPortsGameOrderSysMsgs) > 0 {
		notifyRst := dal_common.PushGameMultipleOrderSysMsg(PServerInfoMgr, common.GetUserId(), eSPortsGameOrderSysMsgs)
		if notifyRst.GetCode() != for_game.C_OPT_SUCCESS {
			logs.Error("======用户投注后批量发送给用户消息失败=======")
		}
	}

	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("投注成功")
	logs.Info("=========RpcESportGameGuessBet 返回===========", rd)
	return rd
}

//竞猜购物车轮询
func (self *sfc) RpcESportGetGameGuessCartPoll(common *base.Common, reqMsg *client_hall.GameGuessCartRequest) *client_hall.GameGuessCartResult {
	//有轮询不要打log
	//logs.Info("===RpcESportGetGameGuessCartPoll===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	rd := &client_hall.GameGuessCartResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}

	guessCartsRst := make([]*client_hall.GameGuessCartResultObj, 0)
	//循环处理
	guessCartsReq := reqMsg.GetGuessCartsReq()
	if nil != guessCartsReq && len(guessCartsReq) > 0 {

		for _, value := range guessCartsReq {
			gameGuessCartResultObj := client_hall.GameGuessCartResultObj{}

			//取得盘口数据
			//查询动态盘口
			guessDetail := for_game.GetRedisGuessDetail(value.GetUniqueGameId(), value.GetAppLabelId(),
				value.GetGameId(),
				value.GetApiOrigin(),
				value.GetMornRollGuessFlag())
			//设置比赛状态及部分数据
			if nil != guessDetail {
				gameGuessCartResultObj.GameStatus = easygo.NewString(for_game.GetGameStatus(guessDetail.GetBeginTime(), guessDetail.GetGameStatus()))
			}

			gameGuessCartResultObj.UniqueGameId = easygo.NewInt64(value.GetUniqueGameId())
			gameGuessCartResultObj.BetNum = easygo.NewString(value.GetBetNum())
			guesses := make([]*share_message.GameGuessOddsNumObject, 0)
			//设置实时赔率
			if nil != guessDetail {
				guesses = guessDetail.GuessOddsNums
			}

			itemOddsMap := make(map[string]string)
			itemBetStatusMap := make(map[string]string)
			if nil != guesses && len(guesses) > 0 {
				for _, guessValue := range guesses {
					contents := guessValue.GetContents()
					for _, contentValue := range contents {
						items := contentValue.GetItems()
						for _, itemValue := range items {
							itemOddsMap[itemValue.GetBetNum()] = itemValue.GetOdds()
							itemBetStatusMap[itemValue.GetBetNum()] = itemValue.GetBetStatus()
						}
					}
				}
			}
			gameGuessCartResultObj.Odds = easygo.NewString(itemOddsMap[value.GetBetNum()])

			//设置投注状态
			gameGuessCartResultObj.BetStatus = easygo.NewString(itemBetStatusMap[value.GetBetNum()])
			//通过比赛状态和盘口再次判断投注状态
			//该对象的状态已经在前面重新设置过直接拿来用
			if (gameGuessCartResultObj.GetGameStatus() == for_game.GAME_STATUS_1 ||
				gameGuessCartResultObj.GetGameStatus() == for_game.GAME_STATUS_2) &&
				value.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_1 {
				gameGuessCartResultObj.BetStatus = easygo.NewString(for_game.GAME_GUESS_ITEM_ODDS_STATUS_2)
			} else if (gameGuessCartResultObj.GetGameStatus() == for_game.GAME_STATUS_0 ||
				gameGuessCartResultObj.GetGameStatus() == for_game.GAME_STATUS_2) &&
				value.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_2 {
				gameGuessCartResultObj.BetStatus = easygo.NewString(for_game.GAME_GUESS_ITEM_ODDS_STATUS_2)
			}

			guessCartsRst = append(guessCartsRst, &gameGuessCartResultObj)
		}
		rd.GuessCartsRst = guessCartsRst
	}
	//有轮询不要打log
	//logs.Info("=========RpcESportGetGameGuessCartPoll 返回===========", rd)

	return rd
}

//比赛历史数据
func (self *sfc) RpcESportGameHistoryData(common *base.Common, reqMsg *client_hall.GameHistoryRequest) *client_hall.GameHistoryResult {

	logs.Info("=======RpcESportGameHistoryData===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	rd := &client_hall.GameHistoryResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_HISTORY)
	defer closeFun()

	query := share_message.RecentData{}
	//通过条件查询
	errQuery := col.Find(bson.M{"_id": reqMsg.GetHistoryId()}).One(&query)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		rd := &client_hall.GameHistoryResult{
			Code: easygo.NewInt32(for_game.C_SYS_ERROR),
			Msg:  easygo.NewString("系统异常"),
		}

		return rd
	}

	if errQuery == mgo.ErrNotFound {
		rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
		rd.Msg = easygo.NewString("")
		return rd
	}

	if errQuery == nil {
		rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
		rd.Msg = easygo.NewString("")
		rd.HisData = &query
		return rd
	}
	logs.Info("=========RpcESportGameHistoryData 返回===========")

	return rd
}

//比赛实时数据
func (self *sfc) RpcESportGameRealTimeData(common *base.Common, reqMsg *client_hall.GameRealTimeRequest) *client_hall.GameRealTimeResult {

	//有轮询不要打log
	//logs.Info("=======RpcESportGameRealTimeData===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	rd := &client_hall.GameRealTimeResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}

	gameRounds := for_game.GetRedisGameRealTimeRounds(reqMsg.GetUniqueGameId())

	if nil != gameRounds {
		rd.GameRounds = easygo.NewInt32(gameRounds.GetGameRounds())
	}

	realTimeData := for_game.GetRedisGameRealTime(reqMsg.GetUniqueGameId(), reqMsg.GetGameRound())
	if nil != realTimeData {
		rd.RealTimeData = realTimeData
	}
	//有轮询不要打log
	//logs.Info("=========RpcESportGameRealTimeData 返回===========")

	return rd
}

func (list GuessOddsNumsSort) Len() int { return len(list) }
func (list GuessOddsNumsSort) Swap(i, j int) {
	s := list[j]
	list[j] = list[i]
	list[i] = s
}

func (list GuessOddsNumsSort) Less(i, j int) bool {
	return list[i].GetNum() < list[j].GetNum()
}
