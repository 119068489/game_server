package sport_api

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
)

//回调处理LOL游戏实时数据
func CallBackLOLRealTimeData(
	appLabelId int32,
	apiOrigin int32,
	gameIdStr string,
	eventId string,
	timestampStr string,
	LOLRealTimeData *share_message.TableESPortsLOLRealTimeData) {

	//s := fmt.Sprintf("=======CallBackLOLRealTimeData回调处理LOL游戏实时数据=========开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v,timestamp=%v",
	//	appLabelId, apiOrigin, gameIdStr, eventId, timestampStr)
	//for_game.WriteFile("ye_zi_api_real_time.log", s)

	//游戏实时数据只打错误日志
	if appLabelId == 0 || apiOrigin == 0 || gameIdStr == "" || eventId == "" || LOLRealTimeData == nil {
		s := fmt.Sprintf("=======CallBackLOLRealTimeData回调处理LOL游戏实时数据参数错误==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v,timestamp=%v,LOLRealTimeData=%v",
			appLabelId, apiOrigin, gameIdStr, eventId, timestampStr, LOLRealTimeData)
		for_game.WriteFile("ye_zi_api_real_time.log", s)
		return
	}

	//从redis取得比赛
	gameObject := for_game.GetRedisGameDetailHeadGroup(apiOrigin, appLabelId, gameIdStr)

	if gameObject == nil {
		s := fmt.Sprintf("=======CallBackLOLRealTimeData回调处理时未取得比赛的数据==========")
		for_game.WriteFile("ye_zi_api_real_time.log", s)
		return
	}

	gameRound := LOLRealTimeData.GetGameRound()
	gameId, err := strconv.ParseInt(gameIdStr, 10, 32)

	if err != nil {
		logs.Error(err)
		s := fmt.Sprintf("=======CallBackLOLRealTimeData回调处理LOL游戏实时数据参数gameIdStr错误==========:参数gameIdStr=%v", gameIdStr)
		for_game.WriteFile("ye_zi_api_real_time.log", s)
		return
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_LOL_REAL_TIME_DATA)
	defer closeFun()

	query := share_message.TableESPortsLOLRealTimeData{}
	//通过条件查询存在更新,不存在就插入
	errQuery := col.Find(bson.M{"app_label_id": appLabelId,
		"game_id":    int32(gameId),
		"api_origin": apiOrigin,
		"game_round": gameRound}).One(&query)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("======CallBackLOLRealTimeData LOL游戏实时数据表查询失败======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v,====game_round:%v",
			appLabelId, gameId, apiOrigin, gameRound)
		for_game.WriteFile("ye_zi_api_real_time.log", s)
		return
	}

	LOLRealTimeData.AppLabelId = easygo.NewInt32(appLabelId)
	LOLRealTimeData.AppLabelName = easygo.NewString(for_game.LabelToESportNameMap[appLabelId])
	LOLRealTimeData.ApiOrigin = easygo.NewInt32(apiOrigin)
	LOLRealTimeData.ApiOriginName = easygo.NewString(for_game.ApiOriginIdToNameMap[apiOrigin])

	//插入数据
	if errQuery == mgo.ErrNotFound {
		//新增数据取得自增id
		id := easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_LOL_REAL_TIME_DATA))
		LOLRealTimeData.Id = id

		//创建时间和更新时间
		nowTime := time.Now().Unix()
		LOLRealTimeData.CreateTime = easygo.NewInt64(nowTime)
		LOLRealTimeData.UpdateTime = easygo.NewInt64(nowTime)

		//取得比赛详情
		gameDetail := GetGamePlayerDetailInfos(apiOrigin, appLabelId, gameIdStr)
		//插入时自动匹配一次头像
		DoMapLOLPlayerPhotos(LOLRealTimeData, gameDetail)

		//记录游戏曲线
		if LOLRealTimeData.TeamA != nil {

			goldTimeDataTeamA := make([]*share_message.GoldTimeData, 0)
			goldTimeDataTeamA = append(goldTimeDataTeamA, &share_message.GoldTimeData{
				GameTimeDistance: easygo.NewInt32(LOLRealTimeData.GetDuration()),
				Gold:             easygo.NewInt32(LOLRealTimeData.TeamA.GetGlod()),
			})

			LOLRealTimeData.TeamA.GoldTimeData = goldTimeDataTeamA
		}

		if LOLRealTimeData.TeamB != nil {
			goldTimeDataTeamB := make([]*share_message.GoldTimeData, 0)
			goldTimeDataTeamB = append(goldTimeDataTeamB, &share_message.GoldTimeData{
				GameTimeDistance: easygo.NewInt32(LOLRealTimeData.GetDuration()),
				Gold:             easygo.NewInt32(LOLRealTimeData.TeamB.GetGlod()),
			})

			LOLRealTimeData.TeamB.GoldTimeData = goldTimeDataTeamB
		}

		if LOLRealTimeData.TeamA != nil && LOLRealTimeData.TeamB != nil {
			//一塔
			if LOLRealTimeData.TeamA.GetTowerState() != 0 && LOLRealTimeData.TeamB.GetTowerState() == 0 {
				LOLRealTimeData.FirstTower = easygo.NewInt32(1)
			} else if LOLRealTimeData.TeamA.GetTowerState() == 0 && LOLRealTimeData.TeamB.GetTowerState() != 0 {
				LOLRealTimeData.FirstTower = easygo.NewInt32(2)
			} else if LOLRealTimeData.TeamA.GetTowerState() != 0 && LOLRealTimeData.TeamB.GetTowerState() != 0 {
				if LOLRealTimeData.TeamA.GetTowerState() > LOLRealTimeData.TeamB.GetTowerState() {
					LOLRealTimeData.FirstTower = easygo.NewInt32(1)
				} else if LOLRealTimeData.TeamA.GetTowerState() < LOLRealTimeData.TeamB.GetTowerState() {
					LOLRealTimeData.FirstTower = easygo.NewInt32(2)
				}
			}

			//一小龙
			if LOLRealTimeData.TeamA.GetDrakes() != 0 && LOLRealTimeData.TeamB.GetDrakes() == 0 {
				LOLRealTimeData.FirstSmallDragon = easygo.NewInt32(1)
			} else if LOLRealTimeData.TeamA.GetDrakes() == 0 && LOLRealTimeData.TeamB.GetDrakes() != 0 {
				LOLRealTimeData.FirstSmallDragon = easygo.NewInt32(2)
			} else if LOLRealTimeData.TeamA.GetDrakes() != 0 && LOLRealTimeData.TeamB.GetDrakes() != 0 {
				if LOLRealTimeData.TeamA.GetDrakes() > LOLRealTimeData.TeamB.GetDrakes() {
					LOLRealTimeData.FirstSmallDragon = easygo.NewInt32(1)
				} else if LOLRealTimeData.TeamA.GetDrakes() < LOLRealTimeData.TeamB.GetDrakes() {
					LOLRealTimeData.FirstSmallDragon = easygo.NewInt32(2)
				}
			}

			//一大龙
			if LOLRealTimeData.TeamA.GetNahsorBarons() != 0 && LOLRealTimeData.TeamB.GetNahsorBarons() == 0 {
				LOLRealTimeData.FirstBigDragon = easygo.NewInt32(1)
			} else if LOLRealTimeData.TeamA.GetNahsorBarons() == 0 && LOLRealTimeData.TeamB.GetNahsorBarons() != 0 {
				LOLRealTimeData.FirstBigDragon = easygo.NewInt32(2)
			} else if LOLRealTimeData.TeamA.GetNahsorBarons() != 0 && LOLRealTimeData.TeamB.GetNahsorBarons() != 0 {
				if LOLRealTimeData.TeamA.GetNahsorBarons() > LOLRealTimeData.TeamB.GetNahsorBarons() {
					LOLRealTimeData.FirstBigDragon = easygo.NewInt32(1)
				} else if LOLRealTimeData.TeamA.GetNahsorBarons() < LOLRealTimeData.TeamB.GetNahsorBarons() {
					LOLRealTimeData.FirstBigDragon = easygo.NewInt32(2)
				}
			}
			//先五杀
			if LOLRealTimeData.TeamA.GetScore() >= 5 && LOLRealTimeData.TeamB.GetScore() < 5 {
				LOLRealTimeData.FirstFiveKill = easygo.NewInt32(1)
			} else if LOLRealTimeData.TeamA.GetScore() < 5 && LOLRealTimeData.TeamB.GetScore() >= 5 {
				LOLRealTimeData.FirstFiveKill = easygo.NewInt32(2)
			} else if LOLRealTimeData.TeamA.GetScore() >= 5 && LOLRealTimeData.TeamB.GetScore() >= 5 {
				if LOLRealTimeData.TeamA.GetScore() > LOLRealTimeData.TeamB.GetScore() {
					LOLRealTimeData.FirstFiveKill = easygo.NewInt32(1)
				} else if LOLRealTimeData.TeamA.GetScore() < LOLRealTimeData.TeamB.GetScore() {
					LOLRealTimeData.FirstFiveKill = easygo.NewInt32(2)
				}
			}
			//先十杀
			if LOLRealTimeData.TeamA.GetScore() >= 10 && LOLRealTimeData.TeamB.GetScore() < 10 {
				LOLRealTimeData.FirstTenKill = easygo.NewInt32(1)
			} else if LOLRealTimeData.TeamA.GetScore() < 10 && LOLRealTimeData.TeamB.GetScore() >= 10 {
				LOLRealTimeData.FirstTenKill = easygo.NewInt32(2)
			} else if LOLRealTimeData.TeamA.GetScore() >= 10 && LOLRealTimeData.TeamB.GetScore() >= 10 {
				if LOLRealTimeData.TeamA.GetScore() > LOLRealTimeData.TeamB.GetScore() {
					LOLRealTimeData.FirstTenKill = easygo.NewInt32(1)
				} else if LOLRealTimeData.TeamA.GetScore() < LOLRealTimeData.TeamB.GetScore() {
					LOLRealTimeData.FirstTenKill = easygo.NewInt32(2)
				}
			}
		}
		errIns := col.Insert(LOLRealTimeData)

		if errIns != nil {
			logs.Error(errIns)
			s := fmt.Sprintf("=======CallBackLOLRealTimeData LOL游戏实时数据表插入数据失败========")
			for_game.WriteFile("ye_zi_api_real_time.log", s)
			return
		}
	}

	//更新数据
	if nil == errQuery {

		//新的更新时间
		LOLRealTimeData.UpdateTime = easygo.NewInt64(time.Now().Unix())

		//将队伍中的数据将数据库的值重新设置到body
		if LOLRealTimeData.TeamA != nil && query.GetTeamA() != nil {

			if nil != query.GetTeamA().GetGoldTimeData() && len(query.GetTeamA().GetGoldTimeData()) > 0 {

				goldTimeDataTeamA := query.GetTeamA().GetGoldTimeData()

				//生产环境上满足条件就记录
				if for_game.IS_FORMAL_SERVER {
					//距离上次记录时间至少30秒以上
					if (LOLRealTimeData.GetDuration())-(goldTimeDataTeamA[len(goldTimeDataTeamA)-1].GetGameTimeDistance()) >= 30 {
						goldTimeDataTeamA = append(goldTimeDataTeamA, &share_message.GoldTimeData{
							GameTimeDistance: easygo.NewInt32(LOLRealTimeData.GetDuration()),
							Gold:             easygo.NewInt32(LOLRealTimeData.TeamA.GetGlod()),
						})
					}
				} else {
					//测试环境上如果点数超过120个就不记录了
					if len(goldTimeDataTeamA) <= 120 {
						//距离上次记录时间至少30秒以上
						if (LOLRealTimeData.GetDuration())-(goldTimeDataTeamA[len(goldTimeDataTeamA)-1].GetGameTimeDistance()) >= 30 {
							goldTimeDataTeamA = append(goldTimeDataTeamA, &share_message.GoldTimeData{
								GameTimeDistance: easygo.NewInt32(LOLRealTimeData.GetDuration()),
								Gold:             easygo.NewInt32(LOLRealTimeData.TeamA.GetGlod()),
							})
						}
					}
				}

				LOLRealTimeData.TeamA.GoldTimeData = goldTimeDataTeamA

			} else {
				goldTimeDataTeamA := make([]*share_message.GoldTimeData, 0)
				goldTimeDataTeamA = append(goldTimeDataTeamA, &share_message.GoldTimeData{
					GameTimeDistance: easygo.NewInt32(LOLRealTimeData.GetDuration()),
					Gold:             easygo.NewInt32(LOLRealTimeData.TeamA.GetGlod()),
				})

				LOLRealTimeData.TeamA.GoldTimeData = goldTimeDataTeamA
			}
		}

		if LOLRealTimeData.TeamB != nil && query.GetTeamB() != nil {

			if nil != query.GetTeamB().GetGoldTimeData() && len(query.GetTeamB().GetGoldTimeData()) > 0 {

				goldTimeDataTeamB := query.GetTeamB().GetGoldTimeData()

				//生产环境上满足条件就记录
				if for_game.IS_FORMAL_SERVER {
					//距离上次记录时间至少30秒以上
					if (LOLRealTimeData.GetDuration())-(goldTimeDataTeamB[len(goldTimeDataTeamB)-1].GetGameTimeDistance()) >= 30 {
						goldTimeDataTeamB = append(goldTimeDataTeamB, &share_message.GoldTimeData{
							GameTimeDistance: easygo.NewInt32(LOLRealTimeData.GetDuration()),
							Gold:             easygo.NewInt32(LOLRealTimeData.TeamB.GetGlod()),
						})
					}
				} else {
					//测试环境上如果点数超过120个就不记录了
					if len(goldTimeDataTeamB) <= 120 {
						//距离上次记录时间至少30秒以上
						if (LOLRealTimeData.GetDuration())-(goldTimeDataTeamB[len(goldTimeDataTeamB)-1].GetGameTimeDistance()) >= 30 {
							goldTimeDataTeamB = append(goldTimeDataTeamB, &share_message.GoldTimeData{
								GameTimeDistance: easygo.NewInt32(LOLRealTimeData.GetDuration()),
								Gold:             easygo.NewInt32(LOLRealTimeData.TeamB.GetGlod()),
							})
						}
					}
				}

				LOLRealTimeData.TeamB.GoldTimeData = goldTimeDataTeamB

			} else {
				goldTimeDataTeamB := make([]*share_message.GoldTimeData, 0)
				goldTimeDataTeamB = append(goldTimeDataTeamB, &share_message.GoldTimeData{
					GameTimeDistance: easygo.NewInt32(LOLRealTimeData.GetDuration()),
					Gold:             easygo.NewInt32(LOLRealTimeData.TeamB.GetGlod()),
				})

				LOLRealTimeData.TeamB.GoldTimeData = goldTimeDataTeamB
			}
		}

		//将队员信息结构中的数据将数据库值设置到body
		if LOLRealTimeData.PlayerAInfo != nil && query.PlayerAInfo != nil {
			for _, value := range LOLRealTimeData.PlayerAInfo {
				for _, dbValue := range query.PlayerAInfo {
					if value.GetName() == dbValue.GetName() {
						value.Photo = easygo.NewString(dbValue.GetPhoto())
						break
					}
				}
			}
		}

		if LOLRealTimeData.PlayerBInfo != nil && query.PlayerBInfo != nil {
			for _, value := range LOLRealTimeData.PlayerBInfo {
				for _, dbValue := range query.PlayerBInfo {
					if value.GetName() == dbValue.GetName() {
						value.Photo = easygo.NewString(dbValue.GetPhoto())
						break
					}
				}
			}
		}

		if LOLRealTimeData.TeamA != nil && LOLRealTimeData.TeamB != nil {
			//先判断数据是否设置过一塔
			if query.GetFirstTower() == 0 {
				if LOLRealTimeData.TeamA.GetTowerState() != 0 && LOLRealTimeData.TeamB.GetTowerState() == 0 {
					LOLRealTimeData.FirstTower = easygo.NewInt32(1)
				} else if LOLRealTimeData.TeamA.GetTowerState() == 0 && LOLRealTimeData.TeamB.GetTowerState() != 0 {
					LOLRealTimeData.FirstTower = easygo.NewInt32(2)
				} else if LOLRealTimeData.TeamA.GetTowerState() != 0 && LOLRealTimeData.TeamB.GetTowerState() != 0 {
					if LOLRealTimeData.TeamA.GetTowerState() > LOLRealTimeData.TeamB.GetTowerState() {
						LOLRealTimeData.FirstTower = easygo.NewInt32(1)
					} else if LOLRealTimeData.TeamA.GetTowerState() < LOLRealTimeData.TeamB.GetTowerState() {
						LOLRealTimeData.FirstTower = easygo.NewInt32(2)
					}
				}
			}

			//先判断数据是否设置过一小龙
			if query.GetFirstSmallDragon() == 0 {
				if LOLRealTimeData.TeamA.GetDrakes() != 0 && LOLRealTimeData.TeamB.GetDrakes() == 0 {
					LOLRealTimeData.FirstSmallDragon = easygo.NewInt32(1)
				} else if LOLRealTimeData.TeamA.GetDrakes() == 0 && LOLRealTimeData.TeamB.GetDrakes() != 0 {
					LOLRealTimeData.FirstSmallDragon = easygo.NewInt32(2)
				} else if LOLRealTimeData.TeamA.GetDrakes() != 0 && LOLRealTimeData.TeamB.GetDrakes() != 0 {
					if LOLRealTimeData.TeamA.GetDrakes() > LOLRealTimeData.TeamB.GetDrakes() {
						LOLRealTimeData.FirstSmallDragon = easygo.NewInt32(1)
					} else if LOLRealTimeData.TeamA.GetDrakes() < LOLRealTimeData.TeamB.GetDrakes() {
						LOLRealTimeData.FirstSmallDragon = easygo.NewInt32(2)
					}
				}
			}

			//先判断数据是否设置一大龙
			if query.GetFirstBigDragon() == 0 {
				if LOLRealTimeData.TeamA.GetNahsorBarons() != 0 && LOLRealTimeData.TeamB.GetNahsorBarons() == 0 {
					LOLRealTimeData.FirstBigDragon = easygo.NewInt32(1)
				} else if LOLRealTimeData.TeamA.GetNahsorBarons() == 0 && LOLRealTimeData.TeamB.GetNahsorBarons() != 0 {
					LOLRealTimeData.FirstBigDragon = easygo.NewInt32(2)
				} else if LOLRealTimeData.TeamA.GetNahsorBarons() != 0 && LOLRealTimeData.TeamB.GetNahsorBarons() != 0 {
					if LOLRealTimeData.TeamA.GetNahsorBarons() > LOLRealTimeData.TeamB.GetNahsorBarons() {
						LOLRealTimeData.FirstBigDragon = easygo.NewInt32(1)
					} else if LOLRealTimeData.TeamA.GetNahsorBarons() < LOLRealTimeData.TeamB.GetNahsorBarons() {
						LOLRealTimeData.FirstBigDragon = easygo.NewInt32(2)
					}
				}
			}

			//先判断数据是否设置先五杀
			if query.GetFirstFiveKill() == 0 {
				if LOLRealTimeData.TeamA.GetScore() >= 5 && LOLRealTimeData.TeamB.GetScore() < 5 {
					LOLRealTimeData.FirstFiveKill = easygo.NewInt32(1)
				} else if LOLRealTimeData.TeamA.GetScore() < 5 && LOLRealTimeData.TeamB.GetScore() >= 5 {
					LOLRealTimeData.FirstFiveKill = easygo.NewInt32(2)
				} else if LOLRealTimeData.TeamA.GetScore() >= 5 && LOLRealTimeData.TeamB.GetScore() >= 5 {
					if LOLRealTimeData.TeamA.GetScore() > LOLRealTimeData.TeamB.GetScore() {
						LOLRealTimeData.FirstFiveKill = easygo.NewInt32(1)
					} else if LOLRealTimeData.TeamA.GetScore() < LOLRealTimeData.TeamB.GetScore() {
						LOLRealTimeData.FirstFiveKill = easygo.NewInt32(2)
					}
				}
			}

			//先判断数据是否设置先十杀
			if query.GetFirstTenKill() == 0 {
				if LOLRealTimeData.TeamA.GetScore() >= 10 && LOLRealTimeData.TeamB.GetScore() < 10 {
					LOLRealTimeData.FirstTenKill = easygo.NewInt32(1)
				} else if LOLRealTimeData.TeamA.GetScore() < 10 && LOLRealTimeData.TeamB.GetScore() >= 10 {
					LOLRealTimeData.FirstTenKill = easygo.NewInt32(2)
				} else if LOLRealTimeData.TeamA.GetScore() >= 10 && LOLRealTimeData.TeamB.GetScore() >= 10 {
					if LOLRealTimeData.TeamA.GetScore() > LOLRealTimeData.TeamB.GetScore() {
						LOLRealTimeData.FirstTenKill = easygo.NewInt32(1)
					} else if LOLRealTimeData.TeamA.GetScore() < LOLRealTimeData.TeamB.GetScore() {
						LOLRealTimeData.FirstTenKill = easygo.NewInt32(2)
					}
				}
			}
		}

		errUpd := col.Update(bson.M{"_id": query.GetId()},
			bson.M{"$set": LOLRealTimeData})

		if errUpd != nil {
			logs.Error(errUpd)
			s := fmt.Sprintf("========CallBackLOLRealTimeData LOL游戏实时数据表更新数据失败======更新条件LOL唯一_id:%v",
				query.GetId())
			for_game.WriteFile("ye_zi_api_real_time.log", s)
			return
		}
	}

	//redis设置
	for_game.SetRedisGameRealTimeRounds(apiOrigin, appLabelId, int32(gameId), gameObject)

	for_game.SetRedisGameRealTime(gameObject.GetUniqueGameId(), gameRound)

	//s = fmt.Sprintf("=======CallBackLOLRealTimeData回调处理LOL游戏实时数据=========结束  ==============")
	//for_game.WriteFile("ye_zi_api_real_time.log", s)
}

//回调处理WZRY游戏实时数据
func CallBackWZRYRealTimeData(appLabelId int32,
	apiOrigin int32,
	gameIdStr string,
	eventId string,
	timestampStr string,
	WZRYRealTimeData *share_message.TableESPortsWZRYRealTimeData) {
	//游戏实时数据只打错误日志
	if appLabelId == 0 || apiOrigin == 0 || gameIdStr == "" || eventId == "" || WZRYRealTimeData == nil {
		s := fmt.Sprintf("=======CallBackWZRYRealTimeData回调处理WZRY游戏实时数据参数错误==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v,timestamp=%v,WZRYRealTimeData=%v",
			appLabelId, apiOrigin, gameIdStr, eventId, timestampStr, WZRYRealTimeData)
		for_game.WriteFile("ye_zi_api_real_time.log", s)
		return
	}

	//从redis取得比赛
	gameObject := for_game.GetRedisGameDetailHeadGroup(apiOrigin, appLabelId, gameIdStr)

	if gameObject == nil {
		s := fmt.Sprintf("=======CallBackWZRYRealTimeData回调处理时未取得比赛的数据==========")
		for_game.WriteFile("ye_zi_api_real_time.log", s)
		return
	}

	gameRound := WZRYRealTimeData.GetGameRound()
	gameId, err := strconv.ParseInt(gameIdStr, 10, 32)

	if err != nil {
		logs.Error(err)
		s := fmt.Sprintf("=======CallBackWZRYRealTimeData回调处理WZRY游戏转换实时数据参数gameIdStr错误==========:gameId=%v", gameIdStr)
		for_game.WriteFile("ye_zi_api_real_time.log", s)
		return
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_WZRY_REAL_TIME_DATA)
	defer closeFun()

	query := share_message.TableESPortsWZRYRealTimeData{}
	//通过条件查询存在更新,不存在就插入
	errQuery := col.Find(bson.M{"app_label_id": appLabelId,
		"game_id":    int32(gameId),
		"api_origin": apiOrigin,
		"game_round": gameRound}).One(&query)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("======CallBackWZRYRealTimeData WZRY游戏实时数据表查询失败======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v,====game_round:%v",
			appLabelId, gameId, apiOrigin, gameRound)
		for_game.WriteFile("ye_zi_api_real_time.log", s)
		return
	}

	WZRYRealTimeData.AppLabelId = easygo.NewInt32(appLabelId)
	WZRYRealTimeData.AppLabelName = easygo.NewString(for_game.LabelToESportNameMap[appLabelId])
	WZRYRealTimeData.ApiOrigin = easygo.NewInt32(apiOrigin)
	WZRYRealTimeData.ApiOriginName = easygo.NewString(for_game.ApiOriginIdToNameMap[apiOrigin])

	//插入数据
	if errQuery == mgo.ErrNotFound {
		//新增数据取得自增id
		id := easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_WZRY_REAL_TIME_DATA))
		WZRYRealTimeData.Id = id
		//创建时间和更新时间
		nowTime := time.Now().Unix()
		WZRYRealTimeData.CreateTime = easygo.NewInt64(nowTime)
		WZRYRealTimeData.UpdateTime = easygo.NewInt64(nowTime)

		//取得比赛详情
		gameDetail := GetGamePlayerDetailInfos(apiOrigin, appLabelId, gameIdStr)
		//插入时自动匹配一次头像
		DoMapWZRYPlayerPhotos(WZRYRealTimeData, gameDetail)

		//将队员信息结构中的数据将数据库值设置到body,同时计算每队的总经济
		var goldTeamA int32
		var goldTeamB int32

		if WZRYRealTimeData.PlayerAInfo != nil && nil != WZRYRealTimeData.TeamA {
			for _, value := range WZRYRealTimeData.PlayerAInfo {
				//加队员的钱作为队伍总经济
				goldTeamA = goldTeamA + value.GetGold()
			}

			WZRYRealTimeData.TeamA.Glod = easygo.NewInt32(goldTeamA)
		}

		if WZRYRealTimeData.PlayerBInfo != nil && nil != WZRYRealTimeData.TeamB {
			for _, value := range WZRYRealTimeData.PlayerBInfo {
				//加队员的钱作为队伍总经济
				goldTeamB = goldTeamB + value.GetGold()
			}

			WZRYRealTimeData.TeamB.Glod = easygo.NewInt32(goldTeamB)
		}

		//记录游戏曲线
		if WZRYRealTimeData.TeamA != nil {

			goldTimeDataTeamA := make([]*share_message.GoldTimeData, 0)
			goldTimeDataTeamA = append(goldTimeDataTeamA, &share_message.GoldTimeData{
				GameTimeDistance: easygo.NewInt32(WZRYRealTimeData.GetDuration()),
				Gold:             easygo.NewInt32(WZRYRealTimeData.TeamA.GetGlod()),
			})

			WZRYRealTimeData.TeamA.GoldTimeData = goldTimeDataTeamA
		}

		if WZRYRealTimeData.TeamB != nil {
			goldTimeDataTeamB := make([]*share_message.GoldTimeData, 0)
			goldTimeDataTeamB = append(goldTimeDataTeamB, &share_message.GoldTimeData{
				GameTimeDistance: easygo.NewInt32(WZRYRealTimeData.GetDuration()),
				Gold:             easygo.NewInt32(WZRYRealTimeData.TeamB.GetGlod()),
			})

			WZRYRealTimeData.TeamB.GoldTimeData = goldTimeDataTeamB
		}

		if WZRYRealTimeData.TeamA != nil && WZRYRealTimeData.TeamB != nil {
			//一塔
			if WZRYRealTimeData.TeamA.GetTowerState() != 0 && WZRYRealTimeData.TeamB.GetTowerState() == 0 {
				WZRYRealTimeData.FirstTower = easygo.NewInt32(1)
			} else if WZRYRealTimeData.TeamA.GetTowerState() == 0 && WZRYRealTimeData.TeamB.GetTowerState() != 0 {
				WZRYRealTimeData.FirstTower = easygo.NewInt32(2)
			} else if WZRYRealTimeData.TeamA.GetTowerState() != 0 && WZRYRealTimeData.TeamB.GetTowerState() != 0 {
				if WZRYRealTimeData.TeamA.GetTowerState() > WZRYRealTimeData.TeamB.GetTowerState() {
					WZRYRealTimeData.FirstTower = easygo.NewInt32(1)
				} else if WZRYRealTimeData.TeamA.GetTowerState() < WZRYRealTimeData.TeamB.GetTowerState() {
					WZRYRealTimeData.FirstTower = easygo.NewInt32(2)
				}
			}

			//一小龙
			if WZRYRealTimeData.TeamA.GetDrakes() != 0 && WZRYRealTimeData.TeamB.GetDrakes() == 0 {
				WZRYRealTimeData.FirstSmallDragon = easygo.NewInt32(1)
			} else if WZRYRealTimeData.TeamA.GetDrakes() == 0 && WZRYRealTimeData.TeamB.GetDrakes() != 0 {
				WZRYRealTimeData.FirstSmallDragon = easygo.NewInt32(2)
			} else if WZRYRealTimeData.TeamA.GetDrakes() != 0 && WZRYRealTimeData.TeamB.GetDrakes() != 0 {
				if WZRYRealTimeData.TeamA.GetDrakes() > WZRYRealTimeData.TeamB.GetDrakes() {
					WZRYRealTimeData.FirstSmallDragon = easygo.NewInt32(1)
				} else if WZRYRealTimeData.TeamA.GetDrakes() < WZRYRealTimeData.TeamB.GetDrakes() {
					WZRYRealTimeData.FirstSmallDragon = easygo.NewInt32(2)
				}
			}

			//一大龙
			if WZRYRealTimeData.TeamA.GetNahsorBarons() != 0 && WZRYRealTimeData.TeamB.GetNahsorBarons() == 0 {
				WZRYRealTimeData.FirstBigDragon = easygo.NewInt32(1)
			} else if WZRYRealTimeData.TeamA.GetNahsorBarons() == 0 && WZRYRealTimeData.TeamB.GetNahsorBarons() != 0 {
				WZRYRealTimeData.FirstBigDragon = easygo.NewInt32(2)
			} else if WZRYRealTimeData.TeamA.GetNahsorBarons() != 0 && WZRYRealTimeData.TeamB.GetNahsorBarons() != 0 {
				if WZRYRealTimeData.TeamA.GetNahsorBarons() > WZRYRealTimeData.TeamB.GetNahsorBarons() {
					WZRYRealTimeData.FirstBigDragon = easygo.NewInt32(1)
				} else if WZRYRealTimeData.TeamA.GetNahsorBarons() < WZRYRealTimeData.TeamB.GetNahsorBarons() {
					WZRYRealTimeData.FirstBigDragon = easygo.NewInt32(2)
				}
			}
			//先五杀
			if WZRYRealTimeData.TeamA.GetScore() >= 5 && WZRYRealTimeData.TeamB.GetScore() < 5 {
				WZRYRealTimeData.FirstFiveKill = easygo.NewInt32(1)
			} else if WZRYRealTimeData.TeamA.GetScore() < 5 && WZRYRealTimeData.TeamB.GetScore() >= 5 {
				WZRYRealTimeData.FirstFiveKill = easygo.NewInt32(2)
			} else if WZRYRealTimeData.TeamA.GetScore() >= 5 && WZRYRealTimeData.TeamB.GetScore() >= 5 {
				if WZRYRealTimeData.TeamA.GetScore() > WZRYRealTimeData.TeamB.GetScore() {
					WZRYRealTimeData.FirstFiveKill = easygo.NewInt32(1)
				} else if WZRYRealTimeData.TeamA.GetScore() < WZRYRealTimeData.TeamB.GetScore() {
					WZRYRealTimeData.FirstFiveKill = easygo.NewInt32(2)
				}
			}
			//先十杀
			if WZRYRealTimeData.TeamA.GetScore() >= 10 && WZRYRealTimeData.TeamB.GetScore() < 10 {
				WZRYRealTimeData.FirstTenKill = easygo.NewInt32(1)
			} else if WZRYRealTimeData.TeamA.GetScore() < 10 && WZRYRealTimeData.TeamB.GetScore() >= 10 {
				WZRYRealTimeData.FirstTenKill = easygo.NewInt32(2)
			} else if WZRYRealTimeData.TeamA.GetScore() >= 10 && WZRYRealTimeData.TeamB.GetScore() >= 10 {
				if WZRYRealTimeData.TeamA.GetScore() > WZRYRealTimeData.TeamB.GetScore() {
					WZRYRealTimeData.FirstTenKill = easygo.NewInt32(1)
				} else if WZRYRealTimeData.TeamA.GetScore() < WZRYRealTimeData.TeamB.GetScore() {
					WZRYRealTimeData.FirstTenKill = easygo.NewInt32(2)
				}
			}

		}

		errIns := col.Insert(WZRYRealTimeData)
		if errIns != nil {
			logs.Error(errIns)
			s := fmt.Sprintf("=======CallBackWZRYRealTimeData WZRY游戏实时数据表插入数据失败========")
			for_game.WriteFile("ye_zi_api_real_time.log", s)
			return
		}
	}

	//更新数据
	if nil == errQuery {
		//新的更新时间
		WZRYRealTimeData.UpdateTime = easygo.NewInt64(time.Now().Unix())

		//将队员信息结构中的数据将数据库值设置到body,同时计算每队的总经济
		var goldTeamA int32
		var goldTeamB int32
		if WZRYRealTimeData.PlayerAInfo != nil && query.PlayerAInfo != nil {
			for _, value := range WZRYRealTimeData.PlayerAInfo {

				//加队员的钱作为队伍总经济
				goldTeamA = goldTeamA + value.GetGold()

				for _, dbValue := range query.PlayerAInfo {
					if value.GetName() == dbValue.GetName() {
						value.Photo = easygo.NewString(dbValue.GetPhoto())
						break
					}
				}
			}
		}

		if WZRYRealTimeData.PlayerBInfo != nil && query.PlayerBInfo != nil {
			for _, value := range WZRYRealTimeData.PlayerBInfo {

				//加队员的钱作为队伍总经济
				goldTeamB = goldTeamB + value.GetGold()

				for _, dbValue := range query.PlayerBInfo {
					if value.GetName() == dbValue.GetName() {
						value.Photo = easygo.NewString(dbValue.GetPhoto())
						break
					}
				}
			}
		}

		//将队伍中的数据将数据库的值重新设置到body
		if WZRYRealTimeData.TeamA != nil && query.GetTeamA() != nil {
			WZRYRealTimeData.TeamA.Drakes = easygo.NewInt32(query.GetTeamA().GetDrakes())
			WZRYRealTimeData.TeamA.NahsorBarons = easygo.NewInt32(query.GetTeamA().GetNahsorBarons())
			WZRYRealTimeData.TeamA.Glod = easygo.NewInt32(goldTeamA)
		}

		if WZRYRealTimeData.TeamB != nil && query.GetTeamB() != nil {
			WZRYRealTimeData.TeamB.Drakes = easygo.NewInt32(query.GetTeamB().GetDrakes())
			WZRYRealTimeData.TeamB.NahsorBarons = easygo.NewInt32(query.GetTeamB().GetNahsorBarons())
			WZRYRealTimeData.TeamB.Glod = easygo.NewInt32(goldTeamB)
		}

		//将队伍中的数据将经济曲线数据库的值重新设置到body
		if WZRYRealTimeData.TeamA != nil && query.GetTeamA() != nil {

			if nil != query.GetTeamA().GetGoldTimeData() && len(query.GetTeamA().GetGoldTimeData()) > 0 {

				goldTimeDataTeamA := query.GetTeamA().GetGoldTimeData()

				//生产环境上满足条件就记录
				if for_game.IS_FORMAL_SERVER {
					//距离上次记录时间至少30秒以上
					if (WZRYRealTimeData.GetDuration())-(goldTimeDataTeamA[len(goldTimeDataTeamA)-1].GetGameTimeDistance()) >= 30 {
						goldTimeDataTeamA = append(goldTimeDataTeamA, &share_message.GoldTimeData{
							GameTimeDistance: easygo.NewInt32(WZRYRealTimeData.GetDuration()),
							Gold:             easygo.NewInt32(WZRYRealTimeData.TeamA.GetGlod()),
						})
					}
				} else {
					//测试环境上如果点数超过120个就不记录了
					if len(goldTimeDataTeamA) <= 120 {
						//距离上次记录时间至少30秒以上
						if (WZRYRealTimeData.GetDuration())-(goldTimeDataTeamA[len(goldTimeDataTeamA)-1].GetGameTimeDistance()) >= 30 {
							goldTimeDataTeamA = append(goldTimeDataTeamA, &share_message.GoldTimeData{
								GameTimeDistance: easygo.NewInt32(WZRYRealTimeData.GetDuration()),
								Gold:             easygo.NewInt32(WZRYRealTimeData.TeamA.GetGlod()),
							})
						}
					}
				}

				WZRYRealTimeData.TeamA.GoldTimeData = goldTimeDataTeamA

			} else {
				goldTimeDataTeamA := make([]*share_message.GoldTimeData, 0)
				goldTimeDataTeamA = append(goldTimeDataTeamA, &share_message.GoldTimeData{
					GameTimeDistance: easygo.NewInt32(WZRYRealTimeData.GetDuration()),
					Gold:             easygo.NewInt32(WZRYRealTimeData.TeamA.GetGlod()),
				})

				WZRYRealTimeData.TeamA.GoldTimeData = goldTimeDataTeamA
			}
		}

		if WZRYRealTimeData.TeamB != nil && query.GetTeamB() != nil {

			if nil != query.GetTeamB().GetGoldTimeData() && len(query.GetTeamB().GetGoldTimeData()) > 0 {

				goldTimeDataTeamB := query.GetTeamB().GetGoldTimeData()

				//生产环境上满足条件就记录
				if for_game.IS_FORMAL_SERVER {
					//距离上次记录时间至少30秒以上
					if (WZRYRealTimeData.GetDuration())-(goldTimeDataTeamB[len(goldTimeDataTeamB)-1].GetGameTimeDistance()) >= 30 {
						goldTimeDataTeamB = append(goldTimeDataTeamB, &share_message.GoldTimeData{
							GameTimeDistance: easygo.NewInt32(WZRYRealTimeData.GetDuration()),
							Gold:             easygo.NewInt32(WZRYRealTimeData.TeamB.GetGlod()),
						})
					}
				} else {
					//测试环境上如果点数超过120个就不记录了
					if len(goldTimeDataTeamB) <= 120 {
						//距离上次记录时间至少30秒以上
						if (WZRYRealTimeData.GetDuration())-(goldTimeDataTeamB[len(goldTimeDataTeamB)-1].GetGameTimeDistance()) >= 30 {
							goldTimeDataTeamB = append(goldTimeDataTeamB, &share_message.GoldTimeData{
								GameTimeDistance: easygo.NewInt32(WZRYRealTimeData.GetDuration()),
								Gold:             easygo.NewInt32(WZRYRealTimeData.TeamB.GetGlod()),
							})
						}
					}
				}

				WZRYRealTimeData.TeamB.GoldTimeData = goldTimeDataTeamB

			} else {
				goldTimeDataTeamB := make([]*share_message.GoldTimeData, 0)
				goldTimeDataTeamB = append(goldTimeDataTeamB, &share_message.GoldTimeData{
					GameTimeDistance: easygo.NewInt32(WZRYRealTimeData.GetDuration()),
					Gold:             easygo.NewInt32(WZRYRealTimeData.TeamB.GetGlod()),
				})

				WZRYRealTimeData.TeamB.GoldTimeData = goldTimeDataTeamB
			}
		}

		if WZRYRealTimeData.TeamA != nil && WZRYRealTimeData.TeamB != nil {
			//先判断数据是否设置过一塔
			if query.GetFirstTower() == 0 {
				if WZRYRealTimeData.TeamA.GetTowerState() != 0 && WZRYRealTimeData.TeamB.GetTowerState() == 0 {
					WZRYRealTimeData.FirstTower = easygo.NewInt32(1)
				} else if WZRYRealTimeData.TeamA.GetTowerState() == 0 && WZRYRealTimeData.TeamB.GetTowerState() != 0 {
					WZRYRealTimeData.FirstTower = easygo.NewInt32(2)
				} else if WZRYRealTimeData.TeamA.GetTowerState() != 0 && WZRYRealTimeData.TeamB.GetTowerState() != 0 {
					if WZRYRealTimeData.TeamA.GetTowerState() > WZRYRealTimeData.TeamB.GetTowerState() {
						WZRYRealTimeData.FirstTower = easygo.NewInt32(1)
					} else if WZRYRealTimeData.TeamA.GetTowerState() < WZRYRealTimeData.TeamB.GetTowerState() {
						WZRYRealTimeData.FirstTower = easygo.NewInt32(2)
					}
				}
			}

			//先判断数据是否设置过一小龙
			if query.GetFirstSmallDragon() == 0 {
				if WZRYRealTimeData.TeamA.GetDrakes() != 0 && WZRYRealTimeData.TeamB.GetDrakes() == 0 {
					WZRYRealTimeData.FirstSmallDragon = easygo.NewInt32(1)
				} else if WZRYRealTimeData.TeamA.GetDrakes() == 0 && WZRYRealTimeData.TeamB.GetDrakes() != 0 {
					WZRYRealTimeData.FirstSmallDragon = easygo.NewInt32(2)
				} else if WZRYRealTimeData.TeamA.GetDrakes() != 0 && WZRYRealTimeData.TeamB.GetDrakes() != 0 {
					if WZRYRealTimeData.TeamA.GetDrakes() > WZRYRealTimeData.TeamB.GetDrakes() {
						WZRYRealTimeData.FirstSmallDragon = easygo.NewInt32(1)
					} else if WZRYRealTimeData.TeamA.GetDrakes() < WZRYRealTimeData.TeamB.GetDrakes() {
						WZRYRealTimeData.FirstSmallDragon = easygo.NewInt32(2)
					}
				}
			}

			//先判断数据是否设置一大龙
			if query.GetFirstBigDragon() == 0 {
				if WZRYRealTimeData.TeamA.GetNahsorBarons() != 0 && WZRYRealTimeData.TeamB.GetNahsorBarons() == 0 {
					WZRYRealTimeData.FirstBigDragon = easygo.NewInt32(1)
				} else if WZRYRealTimeData.TeamA.GetNahsorBarons() == 0 && WZRYRealTimeData.TeamB.GetNahsorBarons() != 0 {
					WZRYRealTimeData.FirstBigDragon = easygo.NewInt32(2)
				} else if WZRYRealTimeData.TeamA.GetNahsorBarons() != 0 && WZRYRealTimeData.TeamB.GetNahsorBarons() != 0 {
					if WZRYRealTimeData.TeamA.GetNahsorBarons() > WZRYRealTimeData.TeamB.GetNahsorBarons() {
						WZRYRealTimeData.FirstBigDragon = easygo.NewInt32(1)
					} else if WZRYRealTimeData.TeamA.GetNahsorBarons() < WZRYRealTimeData.TeamB.GetNahsorBarons() {
						WZRYRealTimeData.FirstBigDragon = easygo.NewInt32(2)
					}
				}
			}

			//先判断数据是否设置先五杀
			if query.GetFirstFiveKill() == 0 {
				if WZRYRealTimeData.TeamA.GetScore() >= 5 && WZRYRealTimeData.TeamB.GetScore() < 5 {
					WZRYRealTimeData.FirstFiveKill = easygo.NewInt32(1)
				} else if WZRYRealTimeData.TeamA.GetScore() < 5 && WZRYRealTimeData.TeamB.GetScore() >= 5 {
					WZRYRealTimeData.FirstFiveKill = easygo.NewInt32(2)
				} else if WZRYRealTimeData.TeamA.GetScore() >= 5 && WZRYRealTimeData.TeamB.GetScore() >= 5 {
					if WZRYRealTimeData.TeamA.GetScore() > WZRYRealTimeData.TeamB.GetScore() {
						WZRYRealTimeData.FirstFiveKill = easygo.NewInt32(1)
					} else if WZRYRealTimeData.TeamA.GetScore() < WZRYRealTimeData.TeamB.GetScore() {
						WZRYRealTimeData.FirstFiveKill = easygo.NewInt32(2)
					}
				}
			}

			//先判断数据是否设置先十杀
			if query.GetFirstTenKill() == 0 {
				if WZRYRealTimeData.TeamA.GetScore() >= 10 && WZRYRealTimeData.TeamB.GetScore() < 10 {
					WZRYRealTimeData.FirstTenKill = easygo.NewInt32(1)
				} else if WZRYRealTimeData.TeamA.GetScore() < 10 && WZRYRealTimeData.TeamB.GetScore() >= 10 {
					WZRYRealTimeData.FirstTenKill = easygo.NewInt32(2)
				} else if WZRYRealTimeData.TeamA.GetScore() >= 10 && WZRYRealTimeData.TeamB.GetScore() >= 10 {
					if WZRYRealTimeData.TeamA.GetScore() > WZRYRealTimeData.TeamB.GetScore() {
						WZRYRealTimeData.FirstTenKill = easygo.NewInt32(1)
					} else if WZRYRealTimeData.TeamA.GetScore() < WZRYRealTimeData.TeamB.GetScore() {
						WZRYRealTimeData.FirstTenKill = easygo.NewInt32(2)
					}
				}
			}
		}

		errUpd := col.Update(bson.M{"_id": query.GetId()},
			bson.M{"$set": WZRYRealTimeData})

		if errUpd != nil {
			logs.Error(errUpd)
			s := fmt.Sprintf("========CallBackWZRYRealTimeData WZRY游戏实时数据表更新数据失败======更新条件WZRY唯一_id:%v",
				query.GetId())
			for_game.WriteFile("ye_zi_api_real_time.log", s)
			return
		}
	}

	//redis设置
	for_game.SetRedisGameRealTimeRounds(apiOrigin, appLabelId, int32(gameId), gameObject)
	for_game.SetRedisGameRealTime(gameObject.GetUniqueGameId(), gameRound)
}

//通过联合主键取得比赛详情的数据
func GetGamePlayerDetailInfos(apiOrigin int32, appLabelId int32, gameIdStr string) *share_message.TableESPortsGameDetail {
	//取得详情表中的信心
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_DETAIL)
	defer closeFun()

	//通过name匹配
	query := share_message.TableESPortsGameDetail{}

	errQuery := col.Find(bson.M{"app_label_id": appLabelId,
		"game_id":    gameIdStr,
		"api_origin": apiOrigin}).One(&query)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("=====GetGamePlayerDetailInfos LOL游戏实时数据中通过联合主键取得比赛详情的数据查询失败======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v",
			appLabelId, gameIdStr, apiOrigin)
		for_game.WriteFile("ye_zi_api_real_time.log", s)
		return nil
	}

	if errQuery == mgo.ErrNotFound {
		s := fmt.Sprintf("=====GetGamePlayerDetailInfos LOL游戏实时数据中通过联合主键取得比赛详情的数据未查询到数据======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v",
			appLabelId, gameIdStr, apiOrigin)
		for_game.WriteFile("ye_zi_api_real_time.log", s)
		return nil
	}

	return &query
}

//LOL自动匹配比赛详情中用户头像
func DoMapLOLPlayerPhotos(LOLRealTimeData *share_message.TableESPortsLOLRealTimeData, gameDetail *share_message.TableESPortsGameDetail) {

	if nil != LOLRealTimeData && nil != gameDetail {
		playerAInfo := LOLRealTimeData.GetPlayerAInfo()
		playerBInfo := LOLRealTimeData.GetPlayerBInfo()

		apiTeamAPlayers := gameDetail.GetApiTeamAPlayers()
		apiTeamBPlayers := gameDetail.GetApiTeamBPlayers()

		if playerAInfo != nil && len(playerAInfo) > 0 && nil != apiTeamAPlayers && len(apiTeamAPlayers) > 0 {

			for _, value := range playerAInfo {
				for _, apiValue := range apiTeamAPlayers {
					if value.GetName() == apiValue.GetName() || value.GetName() == apiValue.GetSn() {
						value.Photo = easygo.NewString(apiValue.GetPhoto())

						break
					}
				}
			}
		}

		if playerBInfo != nil && len(playerBInfo) > 0 && nil != apiTeamBPlayers && len(apiTeamBPlayers) > 0 {

			for _, value := range playerBInfo {
				for _, apiValue := range apiTeamBPlayers {
					if value.GetName() == apiValue.GetName() || value.GetName() == apiValue.GetSn() {
						value.Photo = easygo.NewString(apiValue.GetPhoto())

						break
					}
				}
			}
		}
	}
}

//WZRY自动匹配比赛详情中用户头像
func DoMapWZRYPlayerPhotos(WZRYRealTimeData *share_message.TableESPortsWZRYRealTimeData, gameDetail *share_message.TableESPortsGameDetail) {

	if nil != WZRYRealTimeData && nil != gameDetail {
		playerAInfo := WZRYRealTimeData.GetPlayerAInfo()
		playerBInfo := WZRYRealTimeData.GetPlayerBInfo()

		apiTeamAPlayers := gameDetail.GetApiTeamAPlayers()
		apiTeamBPlayers := gameDetail.GetApiTeamBPlayers()

		if playerAInfo != nil && len(playerAInfo) > 0 && nil != apiTeamAPlayers && len(apiTeamAPlayers) > 0 {

			for _, value := range playerAInfo {
				for _, apiValue := range apiTeamAPlayers {
					if value.GetName() == apiValue.GetName() || value.GetName() == apiValue.GetSn() {
						value.Photo = easygo.NewString(apiValue.GetPhoto())

						break
					}
				}
			}
		}

		if playerBInfo != nil && len(playerBInfo) > 0 && nil != apiTeamBPlayers && len(apiTeamBPlayers) > 0 {

			for _, value := range playerBInfo {
				for _, apiValue := range apiTeamBPlayers {
					if value.GetName() == apiValue.GetName() || value.GetName() == apiValue.GetSn() {
						value.Photo = easygo.NewString(apiValue.GetPhoto())

						break
					}
				}
			}
		}
	}
}
