package sport_api

import (
	"encoding/json"
	"fmt"
	dal_common "game_server/e-sports/sport_common_dal"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"net/url"
	"strconv"
	"time"
)

func DoGetYeZiAllGamesFunc() {
	s := fmt.Sprintf("========DoGetYeZiAllGamesFunc初始化未开始比赛列表 开始==========")
	for_game.WriteFile("ye_zi_api.log", s)
	//王者荣耀
	GetYeZiAllGames(for_game.YEZI_ESPORTS_EVENT_WZRY)
	////dota2 TODO
	//GetYeZiAllGames(for_game.YEZI_ESPORTS_EVENT_DOTA2)
	//英雄联盟lol
	GetYeZiAllGames(for_game.YEZI_ESPORTS_EVENT_LOL)
	////CSGO TODO
	//GetYeZiAllGames(for_game.YEZI_ESPORTS_EVENT_CSGO)

	s = fmt.Sprintf("========DoGetYeZiAllGamesFunc初始化未开始比赛列表 结束==========")
	for_game.WriteFile("ye_zi_api.log", s)
}

func DoGetYeZiAllGameDetailsFunc(incrementFlag int32) {

	s := fmt.Sprintf("========DoGetYeZiAllGameDetailsFunc初始化比赛详情、队伍比赛历史、比赛动态赔率 开始==========")
	for_game.WriteFile("ye_zi_api.log", s)

	//比赛表
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
	defer closeFun()

	//比赛详情表
	colDetail, closeFunDetail := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_DETAIL)
	defer closeFunDetail()

	//比赛使用滚盘表
	colUseRoll, closeFunUseRoll := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_USE_ROLL_GUESS)
	defer closeFunUseRoll()

	//比赛赔率表
	colGuess, closeFunGuess := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_GUESS)
	defer closeFunGuess()

	//api比赛历史表(跟爬虫爬回来的比赛历史区别)
	//目前app应用端未使用到比赛列表的数据
	colHistory, closeFunHistory := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_TEAM_HIS_INFO)
	defer closeFunHistory()

	//取得未开始、进行中的数据
	sportGameLists := make([]*share_message.TableESPortsGame, 0)
	queryBson := bson.M{}
	if incrementFlag == for_game.ESPORTS_INCREMENT_FLAG_1 {
		queryBson = bson.M{"game_status": for_game.GAME_STATUS_0,
			"api_origin": for_game.ESPORTS_API_ORIGIN_ID_YEZI}
	} else {
		timeNow := time.Now().Unix()
		queryBson = bson.M{"game_status": for_game.GAME_STATUS_0,
			"api_origin": for_game.ESPORTS_API_ORIGIN_ID_YEZI, "begin_time_int": bson.M{"$lte": timeNow}}
	}
	errQuery := col.Find(queryBson).All(&sportGameLists)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("======野子科技获取未开始、进行中比赛列表接口查询失败======,查询条件为:game_status:%v,api_origin:%v",
			for_game.GAME_STATUS_0, for_game.ESPORTS_API_ORIGIN_ID_YEZI)
		for_game.WriteFile("ye_zi_api.log", s)

		//easygo.PanicError(errQuery)
		return
	}

	if errQuery == mgo.ErrNotFound || sportGameLists == nil || len(sportGameLists) <= 0 {

		s := fmt.Sprintf("=====野子科技获取未开始、进行中比赛列表接口查询未取得数据=====,查询条件为:game_status:%v,api_origin:%v",
			for_game.GAME_STATUS_0, for_game.ESPORTS_API_ORIGIN_ID_YEZI)
		for_game.WriteFile("ye_zi_api.log", s)

		return
	}

	if errQuery == nil && nil != sportGameLists && len(sportGameLists) > 0 {

		for _, value := range sportGameLists {

			appLabelId := value.GetAppLabelId()
			apiOrigin := value.GetApiOrigin()
			gameId := value.GetGameId()

			eventId := value.GetEventId()

			//未结束的比赛redis信息设置、以下逻辑遇到有比赛更新要更新redis
			for_game.SetRedisGameDetailHead(value.GetId())

			//处理比赛详情、推流地址、两队历史交锋、两队胜败统计、两队天敌克制统计(统计一次即可)在比赛详情中更新
			//早盘、滚盘以及赔率表中的比赛时间、比赛状态在处理比赛时处理
			//初始化时候没有updateTime回调时间
			dealGameDetail(appLabelId, apiOrigin, gameId, eventId,
				colDetail, colUseRoll, colGuess, colHistory, col,
				0,
				for_game.INIT_CALLBACK_FLAG_1)

		}
	}

	s = fmt.Sprintf("========DoGetYeZiAllGameDetailsFunc初始化比赛详情、队伍比赛历史、比赛动态赔率  结束==========")
	for_game.WriteFile("ye_zi_api.log", s)
}

//处理比赛动态信息表(早盘)
func dealGameGuessMorn(appLabelId int32,
	apiOrigin int32,
	gameId string,
	eventId string,
	col *mgo.Collection,
	updateTime int64,
	initCallBackFlag int32) {

	if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {
		s := fmt.Sprintf("=======dealGameGuessMorn初始化处理比赛动态信息表(早盘)  开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
			appLabelId, apiOrigin, gameId, eventId)
		for_game.WriteFile("ye_zi_api.log", s)
	} else {
		s := fmt.Sprintf("=======dealGameGuessMorn回调处理比赛动态信息表(早盘)   开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v,updateTime=%v",
			appLabelId, apiOrigin, gameId, eventId, updateTime)
		for_game.WriteFile("ye_zi_api.log", s)
	}

	timeStr := strconv.FormatInt(time.Now().Unix(), 10)

	param := url.Values{}
	param.Set("app_secret", yeZiAppSecret)
	param.Set("app_key", yeZiAppKey)
	param.Set("timestamp", timeStr)
	param.Set("resource", "dynamic")
	param.Set("func", "game")
	param.Set("event_id", eventId)
	param.Set("game_id", gameId)

	sign := GetYeZiApiSign(param)

	param.Del("app_secret")
	param.Add("sign", sign)

	//发送请求
	data, errGet := Get(yeZiUrl, param)
	if errGet != nil {

		logs.Error(errGet)
		s := fmt.Sprintf("=======野子科技获取早盘数据接口请求api失败======,错误信息:%v,请求参数为::::param=%v", errGet, param)
		for_game.WriteFile("ye_zi_api.log", s)
		//easygo.PanicError(errGet)
		return
	} else {
		resultDB := YeZiESPortsGameGuess{}
		if "" == data {
			s := fmt.Sprintf("==========野子科技获取早盘api数据为空==========,请求参数为:::::param=%v", param)
			for_game.WriteFile("ye_zi_api.log", s)
			return
		} else {
			//s := fmt.Sprintf("=======野子科技获取早盘api数据为::::data=%v,========,请求参数为::::param=%v", data, param)
			//for_game.WriteFile("ye_zi_api.log", s)

			err := json.Unmarshal([]byte(data), &resultDB)

			if err != nil {
				logs.Error(err)
				s := fmt.Sprintf("=========野子科技获取早盘api数据后json转结构体出错::::data=%v========,请求参数为::::param=%v", data, param)
				for_game.WriteFile("ye_zi_api.log", s)
				//easygo.PanicError(err)
				return
			}

			if resultDB.ErrorCode != 0 && resultDB.ErrorMsg != "" {
				s := fmt.Sprintf("======野子科技获取早盘api数据请求返回错误信息:错误码ErrorCode:%v,错误信息ErrorMsg:%v========,请求参数为::::param=%v",
					resultDB.ErrorCode, resultDB.ErrorMsg, param)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}
		}

		//判断转化后结构处理业务
		if resultDB.Code == YEZI_SUCCESS_CODE && resultDB.Data != nil {

			gameGuessMornData := resultDB.Data

			gameGuessMornQuery := share_message.TableESPortsGameGuess{}
			//通过条件查询存在更新,不存在就插入
			errQuery := col.Find(bson.M{"app_label_id": appLabelId,
				"game_id":             gameId,
				"api_origin":          apiOrigin,
				"mornRoll_guess_flag": for_game.GAME_IS_MORN_ROLL_1}).One(&gameGuessMornQuery)

			if errQuery != nil && errQuery != mgo.ErrNotFound {
				logs.Error(errQuery)
				s := fmt.Sprintf("=======处理比赛动态信息表(早盘)查询失败=====,查询条件为:app_label_id:%v,===game_id:%v,===api_origin:%v,===mornRoll_guess_flag:%v",
					appLabelId, gameId, apiOrigin, for_game.GAME_IS_MORN_ROLL_1)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}

			//既然比赛状态修改的时候没有回调、比赛结束才改状态
			if gameGuessMornData.GetGameStatus() == for_game.GAME_STATUS_1 {
				gameGuessMornData.GameStatus = easygo.NewString(for_game.GAME_STATUS_0)
			}

			//不存在早盘数据就插入
			if errQuery == mgo.ErrNotFound {
				//新增数据取得自增id
				id := easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_GAME_GUESS))

				gameGuessMornData.Id = id
				//本项目的id
				gameGuessMornData.AppLabelId = easygo.NewInt32(appLabelId)
				//本项目的名称
				gameGuessMornData.AppLabelName = easygo.NewString(for_game.LabelToESportNameMap[appLabelId])

				//接口来源id
				gameGuessMornData.ApiOrigin = easygo.NewInt32(apiOrigin)
				//接口来源名
				gameGuessMornData.ApiOriginName = easygo.NewString(for_game.ApiOriginIdToNameMap[apiOrigin])

				gameGuessMornData.GameId = easygo.NewString(gameId)

				gameGuessMornData.MornRollGuessFlag = easygo.NewInt32(for_game.GAME_IS_MORN_ROLL_1)
				//创建时间和更新时间
				nowTime := time.Now().Unix()
				gameGuessMornData.CreateTime = easygo.NewInt64(nowTime)
				gameGuessMornData.UpdateTime = easygo.NewInt64(nowTime)

				//后台设置的投注内容的开关默认设置为关闭
				if nil != gameGuessMornData.GetGuess() {
					for _, apiValue := range gameGuessMornData.GetGuess() {

						//设置投注项默认盘口为封盘
						apiValue.AppGuessFlag = easygo.NewInt32(for_game.GAME_APP_GUESS_FLAG_1)

						//设置投注项为不显示
						apiValue.AppGuessViewFlag = easygo.NewInt32(for_game.GAME_APP_GUESS_VIEW_FLAG_1)

						apiItemInfo := apiValue.GetItems()
						//匹配bet_num
						if nil != apiItemInfo {
							for _, apiItemValue := range apiItemInfo {

								//记录封盘时间和结果产生时间
								//1、新增设置封盘时间
								//新增的时候只记录0关闭或者3暂停的时间
								if apiItemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_0 ||
									apiItemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_3 {
									if updateTime != 0 {
										apiItemValue.StatusTime = easygo.NewInt64(updateTime)
									} else {
										apiItemValue.StatusTime = easygo.NewInt64(time.Now().Unix())
									}
								}

								//2、结果产生时间
								if apiItemValue.GetStatus() == for_game.GAME_GUESS_ITEM_STATUS_1 &&
									apiItemValue.GetWin() != for_game.GAME_GUESS_ITEM_WIN_NORST {
									if updateTime != 0 {
										apiItemValue.ResultTime = easygo.NewInt64(updateTime)
									} else {
										apiItemValue.ResultTime = easygo.NewInt64(time.Now().Unix())
									}
								}
							}
						}
					}
				}

				errIns := col.Insert(gameGuessMornData)
				if errIns != nil {
					logs.Error(errIns)
					s := fmt.Sprintf("=======处理比赛动态信息表(早盘)插入数据失败========,插入数据相关信息:app_label_id:%v,===game_id:%v,===api_origin:%v,===mornRoll_guess_flag:%v",
						appLabelId, gameId, apiOrigin, for_game.GAME_IS_MORN_ROLL_1)
					for_game.WriteFile("ye_zi_api.log", s)
					return
				}
			}

			//存在早盘的数据就更新
			if errQuery == nil {

				//新的更新时间
				gameGuessMornData.UpdateTime = easygo.NewInt64(time.Now().Unix())

				//api结构数据
				apiGuessInfo := gameGuessMornData.GetGuess()
				//db查询的结构数据
				queryGuessInfo := gameGuessMornQuery.GetGuess()

				if nil != apiGuessInfo && nil != queryGuessInfo {

					//for _, apiValue := range apiGuessInfo {
					//	var flag bool = false
					//	for _, queryValue := range queryGuessInfo {
					//		//某场比赛中bet_id是唯一的
					//		if queryValue.GetBetId() == apiValue.GetBetId() {
					//			apiValue.AppGuessFlag = easygo.NewInt32(queryValue.GetAppGuessFlag())
					//			apiValue.AppGuessViewFlag = easygo.NewInt32(queryValue.GetAppGuessViewFlag())
					//			flag = true
					//			break
					//		}
					//	}
					//
					//	//代表新增的投注内容
					//	if !flag {
					//		apiValue.AppGuessFlag = easygo.NewInt32(for_game.GAME_APP_GUESS_FLAG_1)
					//		apiValue.AppGuessViewFlag = easygo.NewInt32(for_game.GAME_APP_GUESS_VIEW_FLAG_1)
					//	}
					//}
					//设置投注内容的总控盘口以及显示
					dealBetCntEnableUpdate(apiGuessInfo, queryGuessInfo)

					//设置投注项目的封盘时间和结果产生时间
					dealOddsTimeUpdate(apiGuessInfo, queryGuessInfo, updateTime)
				}

				//比赛的相关的设置
				//这三个值是通过详情来维护的、如果api过来的不是已经结束的状态就用数据库原来的
				if gameGuessMornData.GetGameStatus() != for_game.GAME_STATUS_2 {
					gameGuessMornData.GameStatus = easygo.NewString(gameGuessMornQuery.GetGameStatus())
					gameGuessMornData.GameStatusType = easygo.NewString(gameGuessMornQuery.GetGameStatusType())
					gameGuessMornData.BeginTime = easygo.NewString(gameGuessMornQuery.GetBeginTime())
				}

				errUpd := col.Update(bson.M{"_id": gameGuessMornQuery.GetId()},
					bson.M{"$set": gameGuessMornData})

				if errUpd != nil {
					logs.Error(errUpd)
					s := fmt.Sprintf("========处理比赛动态信息表(早盘)更新数据失败======更新条件盘口唯一_id:%v",
						gameGuessMornData.GetId())
					for_game.WriteFile("ye_zi_api.log", s)
					return
				}
			}
		}
	}

	if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {
		s := fmt.Sprintf("=======dealGameGuessMorn初始化处理比赛动态信息表(早盘)  结束==========")
		for_game.WriteFile("ye_zi_api.log", s)
	} else {
		s := fmt.Sprintf("=======dealGameGuessMorn回调处理比赛动态信息表(早盘)   结束==========")
		for_game.WriteFile("ye_zi_api.log", s)
	}
}

//使用滚盘
func UseGameGuessRoll(appLabelId int32, apiOrigin int32, gameId string, eventId string, col *mgo.Collection) {

	s := fmt.Sprintf("========UseGameGuessRoll野子科技使用滚盘数据处理  开始=========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
		appLabelId, apiOrigin, gameId, eventId)
	for_game.WriteFile("ye_zi_api.log", s)

	//记录使用滚盘记录
	useRoLlQuery := share_message.TableESPortsUseRollGuess{}

	//通过条件查询不存在就插入并且发送api调用、存在不做处理
	errQuery := col.Find(bson.M{"app_label_id": appLabelId,
		"game_id":    gameId,
		"api_origin": apiOrigin}).One(&useRoLlQuery)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("=======处理使用滚盘数据时查询失败======,查询条件为:app_label_id:%v,===game_id:%v,===api_origin:%v",
			appLabelId, gameId, apiOrigin)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	if errQuery == mgo.ErrNotFound {
		var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

		param := url.Values{}
		param.Set("app_secret", yeZiAppSecret)
		param.Set("app_key", yeZiAppKey)
		param.Set("timestamp", timeStr)
		param.Set("resource", "roll")
		param.Set("func", "use_roll")
		param.Set("event_id", eventId)
		param.Set("game_id", gameId)

		sign := GetYeZiApiSign(param)

		param.Del("app_secret")
		param.Add("sign", sign)

		//发送请求
		data, errGet := Get(yeZiUrl, param)
		if errGet != nil {

			logs.Error(errGet)
			s := fmt.Sprintf("=========野子科技使用滚盘请求api失败========,错误信息:%v,请求参数为::::param=%v", errGet, param)
			for_game.WriteFile("ye_zi_api.log", s)
			//easygo.PanicError(errGet)

			return
		} else {
			useRstRoll := YeZiESPortsUSEGuessROLL{}
			if "" == data {
				s := fmt.Sprintf("=========野子科技使用滚盘请求api数据为空========,请求参数为::::param=%v", param)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			} else {
				//s := fmt.Sprintf("=========野子科技使用滚盘请求api数据为::::data=%v,=======,请求参数为::::param=%v", data, param)
				//for_game.WriteFile("ye_zi_api.log", s)

				err := json.Unmarshal([]byte(data), &useRstRoll)

				if err != nil {

					logs.Error(err)
					s := fmt.Sprintf("==========野子科技获取使用滚盘api数据后json转结构体出错::::data=%v=========,请求参数为::::param=%v", data, param)
					for_game.WriteFile("ye_zi_api.log", s)
					//easygo.PanicError(err)
					return
				}

				if useRstRoll.ErrorCode != 0 && useRstRoll.ErrorMsg != "" {
					s := fmt.Sprintf("===========野子科技使用滚盘数据请求返回错误信息:错误码ErrorCode:%v,错误信息ErrorMsg:%v,========,请求参数为::::param=%v",
						useRstRoll.ErrorCode, useRstRoll.ErrorMsg, param)
					for_game.WriteFile("ye_zi_api.log", s)
					return
				}

				//没有就插入
				//判断转化后结构处理业务
				if useRstRoll.Code == YEZI_SUCCESS_CODE {

					insUseRoll := share_message.TableESPortsUseRollGuess{}
					insUseRoll.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_USE_ROLL_GUESS))
					//本项目的id
					insUseRoll.AppLabelId = easygo.NewInt32(appLabelId)
					//本项目的名称
					insUseRoll.AppLabelName = easygo.NewString(for_game.LabelToESportNameMap[appLabelId])

					//接口来源id
					insUseRoll.ApiOrigin = easygo.NewInt32(apiOrigin)
					//接口来源名
					insUseRoll.ApiOriginName = easygo.NewString(for_game.ApiOriginIdToNameMap[apiOrigin])

					insUseRoll.GameId = easygo.NewString(gameId)

					//创建时间和更新时间
					nowTime := time.Now().Unix()
					insUseRoll.CreateTime = easygo.NewInt64(nowTime)
					insUseRoll.UpdateTime = easygo.NewInt64(nowTime)

					errIns := col.Insert(&insUseRoll)
					if errIns != nil {
						logs.Error(errIns)
						s := fmt.Sprintf("=======处理使用滚盘时插入数据失败========,插入相关信息为app_label_id:%v,===game_id:%v,===api_origin:%v",
							appLabelId, gameId, apiOrigin)
						for_game.WriteFile("ye_zi_api.log", s)
						return
					}
				}
			}
		}
	}

	s = fmt.Sprintf("========UseGameGuessRoll野子科技使用滚盘数据处理   结束==========")
	for_game.WriteFile("ye_zi_api.log", s)
}

//处理比赛动态信息表(滚盘)
//initCallBackFlag为2的时候是回调
//updateTime回调推送时间
func dealGameGuessRoll(appLabelId int32,
	apiOrigin int32,
	gameId string,
	eventId string,
	col *mgo.Collection,
	updateTime int64,
	initCallBackFlag int32) {

	if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {
		s := fmt.Sprintf("=======dealGameGuessRoll初始化处理比赛动态信息表(滚盘)  开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
			appLabelId, apiOrigin, gameId, eventId)
		for_game.WriteFile("ye_zi_api.log", s)
	} else {
		s := fmt.Sprintf("=======dealGameGuessRoll回调处理比赛动态信息表(滚盘)   开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v,updateTime=%v",
			appLabelId, apiOrigin, gameId, eventId, updateTime)
		for_game.WriteFile("ye_zi_api.log", s)
	}

	var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

	param := url.Values{}
	param.Set("app_secret", yeZiAppSecret)
	param.Set("app_key", yeZiAppKey)
	param.Set("timestamp", timeStr)
	param.Set("resource", "roll")
	param.Set("func", "get_info")
	param.Set("event_id", eventId)
	param.Set("game_id", gameId)

	sign := GetYeZiApiSign(param)

	param.Del("app_secret")
	param.Add("sign", sign)

	//发送请求
	data, errGet := Get(yeZiUrl, param)
	if errGet != nil {

		logs.Error(errGet)
		s := fmt.Sprintf("=======野子科技获取滚盘api数据接口请求失败======,错误信息:%v,请求参数为::::param=%v", errGet, param)
		for_game.WriteFile("ye_zi_api.log", s)
		//easygo.PanicError(errGet)
		return
	} else {
		resultDB := YeZiESPortsGameGuess{}
		if "" == data {
			s := fmt.Sprintf("========野子科技获取滚盘api数据为空=======,请求参数为::::param=%v", param)
			for_game.WriteFile("ye_zi_api.log", s)
			return
		} else {

			//s := fmt.Sprintf("=========野子科技获取滚盘api数据为::::data=%v,=======,请求参数为::::param=%v", data, param)
			//for_game.WriteFile("ye_zi_api.log", s)

			err := json.Unmarshal([]byte(data), &resultDB)

			if err != nil {
				logs.Error(err)
				s := fmt.Sprintf("==========野子科技获取滚盘api数据后json转结构体出错::::data=%v=========,请求参数为::::param=%v", data, param)
				for_game.WriteFile("ye_zi_api.log", s)
				//easygo.PanicError(err)
				return
			}

			if resultDB.ErrorCode != 0 && resultDB.ErrorMsg != "" {
				s := fmt.Sprintf("========野子科技获取滚盘数据接口请求返回错误信息:错误码:%v,错误信息:%v=======,请求参数为::::param=%v",
					resultDB.ErrorCode, resultDB.ErrorMsg, param)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}
		}

		//判断转化后结构处理业务
		if resultDB.Code == YEZI_SUCCESS_CODE && resultDB.Data != nil {

			gameGuessRollData := resultDB.Data
			gameGuessRollQuery := share_message.TableESPortsGameGuess{}
			//通过条件查询存在更新,不存在就插入
			errQuery := col.Find(bson.M{"app_label_id": appLabelId,
				"game_id":             gameId,
				"api_origin":          apiOrigin,
				"mornRoll_guess_flag": for_game.GAME_IS_MORN_ROLL_2}).One(&gameGuessRollQuery)

			if errQuery != nil && errQuery != mgo.ErrNotFound {
				logs.Error(errQuery)
				s := fmt.Sprintf("=========处理比赛动态信息表(滚盘)查询失败======,查询条件:app_label_id:%v,===game_id:%v,===api_origin:%v,===mornRoll_guess_flag:%v",
					appLabelId, gameId, apiOrigin, for_game.GAME_IS_MORN_ROLL_2)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}

			//既然比赛状态修改的时候没有回调、比赛结束才改状态
			if gameGuessRollData.GetGameStatus() == for_game.GAME_STATUS_1 {
				gameGuessRollData.GameStatus = easygo.NewString(for_game.GAME_STATUS_0)
			}

			//没有就插入
			if errQuery == mgo.ErrNotFound {
				//新增数据取得自增id
				id := easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_GAME_GUESS))

				gameGuessRollData.Id = id
				//本项目的id
				gameGuessRollData.AppLabelId = easygo.NewInt32(appLabelId)
				//本项目的名称
				gameGuessRollData.AppLabelName = easygo.NewString(for_game.LabelToESportNameMap[appLabelId])

				//接口来源id
				gameGuessRollData.ApiOrigin = easygo.NewInt32(apiOrigin)
				//接口来源名
				gameGuessRollData.ApiOriginName = easygo.NewString(for_game.ApiOriginIdToNameMap[apiOrigin])

				gameGuessRollData.GameId = easygo.NewString(gameId)

				gameGuessRollData.MornRollGuessFlag = easygo.NewInt32(for_game.GAME_IS_MORN_ROLL_2)
				//创建时间和更新时间
				nowTime := time.Now().Unix()
				gameGuessRollData.CreateTime = easygo.NewInt64(nowTime)
				gameGuessRollData.UpdateTime = easygo.NewInt64(nowTime)

				if nil != gameGuessRollData.GetGuess() {
					//后台设置的投注内容的开关默认设置为关闭
					for _, apiValue := range gameGuessRollData.GetGuess() {
						apiValue.AppGuessFlag = easygo.NewInt32(for_game.GAME_APP_GUESS_FLAG_1)
						apiValue.AppGuessViewFlag = easygo.NewInt32(for_game.GAME_APP_GUESS_VIEW_FLAG_1)

						apiItemInfo := apiValue.GetItems()
						//匹配bet_num
						if nil != apiItemInfo {
							for _, apiItemValue := range apiItemInfo {

								//1、新增设置封盘时间
								//新增的时候只记录0关闭或者3暂停的时间
								if apiItemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_0 ||
									apiItemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_3 {
									if updateTime != 0 {
										apiItemValue.StatusTime = easygo.NewInt64(updateTime)
									} else {
										apiItemValue.StatusTime = easygo.NewInt64(time.Now().Unix())
									}
								}

								//2、结果产生时间
								if apiItemValue.GetStatus() == for_game.GAME_GUESS_ITEM_STATUS_1 &&
									apiItemValue.GetWin() != for_game.GAME_GUESS_ITEM_WIN_NORST {
									if updateTime != 0 {
										apiItemValue.ResultTime = easygo.NewInt64(updateTime)
									} else {
										apiItemValue.ResultTime = easygo.NewInt64(time.Now().Unix())
									}
								}
							}
						}
					}
				}
				errIns := col.Insert(gameGuessRollData)
				if errIns != nil {
					logs.Error(errIns)
					s := fmt.Sprintf("=========处理比赛动态信息表(滚盘)插入数据失败=====,插入相关数据信息:app_label_id:%v,===game_id:%v,===api_origin:%v,====mornRoll_guess_flag:%v",
						appLabelId,
						gameId,
						apiOrigin,
						for_game.GAME_IS_MORN_ROLL_2)
					for_game.WriteFile("ye_zi_api.log", s)

				}
			}

			if errQuery == nil {

				//新的更新时间
				gameGuessRollData.UpdateTime = easygo.NewInt64(time.Now().Unix())

				apiGuessInfo := gameGuessRollData.GetGuess()
				queryGuessInfo := gameGuessRollQuery.GetGuess()

				if nil != apiGuessInfo && nil != queryGuessInfo {
					////后台设置的投注内容的开启与关闭需要重新设置到更新结构体中
					//for _, apiValue := range apiGuessInfo {
					//	var flag bool = false
					//	for _, queryValue := range queryGuessInfo {
					//
					//		//同一场比赛中bet_id是唯一的
					//		if queryValue.GetBetId() == apiValue.GetBetId() {
					//			apiValue.AppGuessFlag = easygo.NewInt32(queryValue.GetAppGuessFlag())
					//			apiValue.AppGuessViewFlag = easygo.NewInt32(queryValue.GetAppGuessViewFlag())
					//			flag = true
					//			break
					//		}
					//	}
					//	//代表新增的投注内容
					//	if !flag {
					//		apiValue.AppGuessFlag = easygo.NewInt32(for_game.GAME_APP_GUESS_FLAG_1)
					//		apiValue.AppGuessViewFlag = easygo.NewInt32(for_game.GAME_APP_GUESS_VIEW_FLAG_1)
					//	}
					//}

					//设置投注内容的总控盘口以及显示
					dealBetCntEnableUpdate(apiGuessInfo, queryGuessInfo)

					//设置投注项目的封盘时间和结果产生时间
					dealOddsTimeUpdate(apiGuessInfo, queryGuessInfo, updateTime)
				}

				//比赛相关的设置
				//这三个值是通过详情来维护的、如果api过来的不是已经结束的状态就用数据库原来的
				if gameGuessRollData.GetGameStatus() != for_game.GAME_STATUS_2 {
					gameGuessRollData.GameStatus = easygo.NewString(gameGuessRollQuery.GetGameStatus())
					gameGuessRollData.GameStatusType = easygo.NewString(gameGuessRollQuery.GetGameStatusType())
					gameGuessRollData.BeginTime = easygo.NewString(gameGuessRollQuery.GetBeginTime())
				}

				errUpd := col.Update(bson.M{"_id": gameGuessRollQuery.GetId()},
					bson.M{"$set": gameGuessRollData})

				if errUpd != nil {
					logs.Error(errUpd)
					s := fmt.Sprintf("=======处理比赛动态信息表(滚盘)更新失败=====盘口更新唯一_id:%v",
						gameGuessRollData.GetId())
					for_game.WriteFile("ye_zi_api.log", s)
					return
				}
			}
		}
	}

	if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {
		s := fmt.Sprintf("=======dealGameGuessRoll初始化处理比赛动态信息表(滚盘)  结束==========")
		for_game.WriteFile("ye_zi_api.log", s)
	} else {
		s := fmt.Sprintf("=======dealGameGuessRoll回调处理比赛动态信息表(滚盘)   结束==========")
		for_game.WriteFile("ye_zi_api.log", s)
	}
}

//两队历史交锋、两队胜败统计、两队天敌克制统计(统计一次即可、回调时候如果是更新要比较是否变了队伍)
func dealTeamHistory(appLabelId int32, apiOrigin int32, gameId string, eventId string, col *mgo.Collection) {

	s := fmt.Sprintf("=======dealTeamHistory处理两队历史交锋、两队胜败统计、两队天敌克制统计   开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
		appLabelId, apiOrigin, gameId, eventId)
	for_game.WriteFile("ye_zi_api.log", s)

	//两组对战
	rstTeamBout := YeZiESPortsTeamBout{}
	//两队胜败
	rstTeamWinFail := YeZiESPortsTeamWinFailInfo{}
	//两队天敌克制
	rstTeamNatRes := YeZiESPortsTeamNatResInfo{}

	//两组对战
	var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

	param := url.Values{}
	param.Set("app_secret", yeZiAppSecret)
	param.Set("app_key", yeZiAppKey)
	param.Set("timestamp", timeStr)
	param.Set("resource", "game")
	param.Set("func", "info_extend")
	param.Set("event_id", eventId)
	param.Set("game_id", gameId)

	sign := GetYeZiApiSign(param)

	param.Del("app_secret")
	param.Add("sign", sign)

	//发送请求
	data, errGet := Get(yeZiUrl, param)
	if errGet != nil {

		logs.Error(errGet)
		s := fmt.Sprintf("========野子科技获取比赛历史对阵相关（两组对决）api接口请求失败======,错误信息:%v,请求参数为::::param=%v", errGet, param)
		for_game.WriteFile("ye_zi_api.log", s)
		//easygo.PanicError(errGet)
		return
	} else {

		if "" == data {
			s := fmt.Sprintf("======野子科技获取比赛历史对阵相关（两组对决）api接口数据为空======,请求参数为::::param=%v", param)
			for_game.WriteFile("ye_zi_api.log", s)
		} else {

			//s := fmt.Sprintf("=======野子科技获取比赛历史对阵相关（两组对决）api接口数据为::::data=%v,=======,请求参数为::::param=%v", data, param)
			//for_game.WriteFile("ye_zi_api.log", s)

			err := json.Unmarshal([]byte(data), &rstTeamBout)

			if err != nil {
				logs.Error(err)
				s := fmt.Sprintf("======野子科技获取比赛历史对阵相关（两组对决）api接口数据后json转结构体出错::::data=%v======,请求参数为::::param=%v", data, param)
				for_game.WriteFile("ye_zi_api.log", s)
				//easygo.PanicError(err)
				return
			}

			if rstTeamBout.ErrorCode != 0 && rstTeamBout.ErrorMsg != "" {
				s := fmt.Sprintf("=======野子科技获取比赛历史对阵相关（两组对决）接口请求返回错误信息:错误码:%v,错误信息:%v======,请求参数为::::param=%v",
					rstTeamBout.ErrorCode, rstTeamBout.ErrorMsg, param)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}
		}
	}

	//两队胜败
	var timeStr1 string = strconv.FormatInt(time.Now().Unix(), 10)

	param1 := url.Values{}
	param1.Set("app_secret", yeZiAppSecret)
	param1.Set("app_key", yeZiAppKey)
	param1.Set("timestamp", timeStr1)
	param1.Set("resource", "game")
	param1.Set("func", "continueWin")
	param1.Set("event_id", eventId)
	param1.Set("game_id", gameId)

	sign1 := GetYeZiApiSign(param1)

	param1.Del("app_secret")
	param1.Add("sign", sign1)

	//发送请求
	data1, errGet1 := Get(yeZiUrl, param1)
	if errGet1 != nil {

		logs.Error(errGet1)
		s := fmt.Sprintf("======野子科技获取两队胜败统计api接口请求失败=====,错误信息:%v,请求参数为::::param=%v", errGet1, param1)
		for_game.WriteFile("ye_zi_api.log", s)
		//easygo.PanicError(errGet1)
		return
	} else {

		if "" == data1 {
			s := fmt.Sprintf("=======野子科技获取两队胜败统计api接口数据为空======,请求参数为::::param=%v", param1)
			for_game.WriteFile("ye_zi_api.log", s)
		} else {

			//s := fmt.Sprintf("=======野子科技获取两队胜败统计api接口数据为::::data=%v,=======,请求参数为::::param=%v", data1, param1)
			//for_game.WriteFile("ye_zi_api.log", s)

			err1 := json.Unmarshal([]byte(data1), &rstTeamWinFail)

			if err1 != nil {
				logs.Error(err1)
				s := fmt.Sprintf("======野子科技获取两队胜败统计api接口数据后json转结构体出错::::data=%v======,请求参数为::::param=%v", data1, param1)
				for_game.WriteFile("ye_zi_api.log", s)
				//easygo.PanicError(err1)
				return
			}

			if rstTeamWinFail.ErrorCode != 0 && rstTeamWinFail.ErrorMsg != "" {
				s := fmt.Sprintf("野子科技获取两队胜败统计接口请求返回错误信息:错误码:%v,错误信息:%v,请求参数为::::param=%v",
					rstTeamWinFail.ErrorCode, rstTeamWinFail.ErrorMsg, param1)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}

			if rstTeamWinFail.Code == YEZI_SUCCESS_CODE && rstTeamWinFail.Data != nil {
				if nil != rstTeamWinFail.Data.TeamA {
					rstTeamBout.Data.TeamAWinFail = &share_message.APITeamWinFaiObject{
						IsContinueWin: easygo.NewInt32(rstTeamWinFail.Data.TeamA.IsContinueWin),
						Num:           easygo.NewInt32(rstTeamWinFail.Data.TeamA.Num),
						TeamId:        easygo.NewString(rstTeamWinFail.Data.TeamA.TeamID),
					}
				}

				if nil != rstTeamWinFail.Data.TeamB {
					rstTeamBout.Data.TeamBWinFail = &share_message.APITeamWinFaiObject{
						IsContinueWin: easygo.NewInt32(rstTeamWinFail.Data.TeamB.IsContinueWin),
						Num:           easygo.NewInt32(rstTeamWinFail.Data.TeamB.Num),
						TeamId:        easygo.NewString(rstTeamWinFail.Data.TeamB.TeamID),
					}
				}
			}
		}
	}

	//两队天敌克制
	var timeStr2 string = strconv.FormatInt(time.Now().Unix(), 10)

	param2 := url.Values{}
	param2.Set("app_secret", yeZiAppSecret)
	param2.Set("app_key", yeZiAppKey)
	param2.Set("timestamp", timeStr2)
	param2.Set("resource", "game")
	param2.Set("func", "teamOpponent")
	param2.Set("event_id", eventId)
	param2.Set("game_id", gameId)

	sign2 := GetYeZiApiSign(param2)

	param2.Del("app_secret")
	param2.Add("sign", sign2)

	//发送请求
	data2, errGet2 := Get(yeZiUrl, param2)
	if errGet2 != nil {

		logs.Error(errGet2)
		s := fmt.Sprintf("=======野子科技两队天敌克制统计api接口请求失败=====,错误信息:%v,请求参数为::::param=%v", errGet2, param2)
		for_game.WriteFile("ye_zi_api.log", s)
		//easygo.PanicError(errGet2)
		return
	} else {

		if "" == data2 {
			s := fmt.Sprintf("=========野子科技两队天敌克制统计api接口数据为空=======,请求参数为::::param=%v", param2)
			for_game.WriteFile("ye_zi_api.log", s)
		} else {

			//s := fmt.Sprintf("=======野子科技两队天敌克制统计api接口数据为::::data=%v,=======,请求参数为::::param=%v", data2, param2)
			//for_game.WriteFile("ye_zi_api.log", s)

			err2 := json.Unmarshal([]byte(data2), &rstTeamNatRes)

			if err2 != nil {
				logs.Error(err2)
				s := fmt.Sprintf("======野子科技两队天敌克制统计api接口数据后json转结构体出错::::data=%v======,请求参数为::::param=%v", data2, param2)
				for_game.WriteFile("ye_zi_api.log", s)
				//easygo.PanicError(err2)
				return
			}

			if rstTeamNatRes.ErrorCode != 0 && rstTeamNatRes.ErrorMsg != "" {
				s := fmt.Sprintf("=====野子科技两队天敌克制统计接口请求返回错误信息:错误码:%v,错误信息:%v=======,,请求参数为::::param=%v",
					rstTeamNatRes.ErrorCode, rstTeamNatRes.ErrorMsg, param2)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}

			if rstTeamNatRes.Code == YEZI_SUCCESS_CODE && rstTeamNatRes.Data != nil {
				if nil != rstTeamNatRes.Data.TeamA {
					rstTeamBout.Data.TeamANatRes = &share_message.APITeamNatResObject{
						NaturalTeam:  easygo.NewString(rstTeamNatRes.Data.TeamA.NaturalTeam),
						RestrainTeam: easygo.NewString(rstTeamNatRes.Data.TeamA.RestrainTeam),
						TeamId:       easygo.NewString(rstTeamNatRes.Data.TeamA.TeamID),
					}
				}

				if nil != rstTeamNatRes.Data.TeamB {
					rstTeamBout.Data.TeamBNatRes = &share_message.APITeamNatResObject{
						NaturalTeam:  easygo.NewString(rstTeamNatRes.Data.TeamB.NaturalTeam),
						RestrainTeam: easygo.NewString(rstTeamNatRes.Data.TeamB.RestrainTeam),
						TeamId:       easygo.NewString(rstTeamNatRes.Data.TeamB.TeamID),
					}
				}
			}
		}
	}

	//处理入库
	gameTeamHisData := rstTeamBout.Data
	gameTeamHisQuery := share_message.TableESPortsTeamBout{}
	//通过条件查询存在更新,不存在就插入
	errQuery := col.Find(bson.M{"app_label_id": appLabelId,
		"game_id":    gameId,
		"api_origin": apiOrigin}).One(&gameTeamHisQuery)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("=====处理两队历史交锋、两队胜败统计、两队天敌克制统计时查询失败=====查询条件为app_label_id:%v,===game_id:%v,===api_origin:%v",
			appLabelId, gameId, apiOrigin)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	//没有就插入
	if errQuery == mgo.ErrNotFound {

		//新增数据取得自增id
		gameTeamHisData.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_TEAM_HIS_INFO))
		//本项目的id
		gameTeamHisData.AppLabelId = easygo.NewInt32(appLabelId)
		//本项目的名称
		gameTeamHisData.AppLabelName = easygo.NewString(for_game.LabelToESportNameMap[appLabelId])

		//接口来源id
		gameTeamHisData.ApiOrigin = easygo.NewInt32(apiOrigin)
		//接口来源名
		gameTeamHisData.ApiOriginName = easygo.NewString(for_game.ApiOriginIdToNameMap[apiOrigin])

		//游戏id
		gameTeamHisData.GameId = easygo.NewString(gameId)

		//创建时间和更新时间
		nowTime := time.Now().Unix()
		gameTeamHisData.CreateTime = easygo.NewInt64(nowTime)
		gameTeamHisData.UpdateTime = easygo.NewInt64(nowTime)

		errIns := col.Insert(gameTeamHisData)
		if errIns != nil {
			s := fmt.Sprintf("=====处理两队历史交锋、两队胜败统计、两队天敌克制统计时插入数据失败===,插入相关信息:app_label_id:%v,===game_id:%v,===api_origin:%v",
				appLabelId, gameId, apiOrigin)
			for_game.WriteFile("ye_zi_api.log", s)

		}
	}

	if errQuery == nil {

		//新的更新时间
		gameTeamHisData.UpdateTime = easygo.NewInt64(time.Now().Unix())

		errUpd := col.Update(bson.M{"_id": gameTeamHisQuery.GetId()},
			bson.M{"$set": gameTeamHisData})

		if errUpd != nil {
			logs.Error(errUpd)
			s := fmt.Sprintf("处理两队历史交锋、两队胜败统计、两队天敌克制统计时更新数据失败=====更新唯一_id:%v",
				gameTeamHisQuery.GetId())
			for_game.WriteFile("ye_zi_api.log", s)
		}
	}

	s = fmt.Sprintf("=======dealTeamHistory处理两队历史交锋、两队胜败统计、两队天敌克制统计   结束==========")
	for_game.WriteFile("ye_zi_api.log", s)
}

//处理比赛详情
//initCallBackFlag 1为初始化、2的时候是回调 updateTime api推送时间
func dealGameDetail(appLabelId int32,
	apiOrigin int32,
	gameId string,
	eventId string,
	colDetail *mgo.Collection,
	colUseRoll *mgo.Collection,
	colGuess *mgo.Collection,
	colHistory *mgo.Collection,
	col *mgo.Collection,
	updateTime int64,
	initCallBackFlag int32) {

	if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {
		s := fmt.Sprintf("=======dealGameDetail初始化处理比赛详情  开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
			appLabelId, apiOrigin, gameId, eventId)
		for_game.WriteFile("ye_zi_api.log", s)
	} else {
		s := fmt.Sprintf("=======dealGameDetail回调处理比赛详情   开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v,updateTime=%v",
			appLabelId, apiOrigin, gameId, eventId, updateTime)
		for_game.WriteFile("ye_zi_api.log", s)
	}

	timeStr := strconv.FormatInt(time.Now().Unix(), 10)

	param := url.Values{}
	param.Set("app_secret", yeZiAppSecret)
	param.Set("app_key", yeZiAppKey)
	param.Set("timestamp", timeStr)
	param.Set("resource", "game")
	param.Set("func", "info")
	param.Set("event_id", eventId)
	param.Set("id", gameId)

	sign := GetYeZiApiSign(param)

	param.Del("app_secret")
	param.Add("sign", sign)

	//发送请求
	data, errGet := Get(yeZiUrl, param)
	if errGet != nil {

		logs.Error(errGet)
		s := fmt.Sprintf("=====野子科技获取比赛详情api接口请求失败======,错误信息:%v,请求参数为::::param=%v", errGet, param)
		for_game.WriteFile("ye_zi_api.log", s)
		//easygo.PanicError(errGet)
		return

	} else {
		resultDB := YeZiESPortsGameDetailInfo{}
		rstMiddleMap := YeZiESPortsMiddleDetailMapInfo{}
		if "" == data {
			s := fmt.Sprintf("========野子科技获取比赛详情api接口返回数据为空=====,请求参数为::::param=%v", param)
			for_game.WriteFile("ye_zi_api.log", s)
			return
		} else {

			//s := fmt.Sprintf("=======野子科技获取比赛详情api接口返回数据为::::data=%v,=======,请求参数为::::param=%v", data, param)
			//for_game.WriteFile("ye_zi_api.log", s)

			err := json.Unmarshal([]byte(data), &resultDB)

			if err != nil {
				logs.Error(err)
				s := fmt.Sprintf("======野子科技获取比赛详情api接口返回数据后json转结构体出错::::data=%v======,请求参数为::::param=%v", data, param)
				for_game.WriteFile("ye_zi_api.log", s)
				//easygo.PanicError(err)
				return
			}

			if resultDB.ErrorCode != 0 && resultDB.ErrorMsg != "" {
				s := fmt.Sprintf("=====野子科技获取比赛详情api接口返回错误信息:错误码:%v,错误信息:%v,请求参数为::::param=%v",
					resultDB.ErrorCode, resultDB.ErrorMsg, param)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}

			//取得中间map的值
			if resultDB.Code == YEZI_SUCCESS_CODE && resultDB.Data != nil {
				err1 := json.Unmarshal([]byte(data), &rstMiddleMap)

				if err1 != nil {
					logs.Error(err1)
					s := fmt.Sprintf("=====野子科技获取比赛详情api接口返回数据后json转中间map结构体出错::::data=%v======,请求参数为::::param=%v",
						data, param)
					for_game.WriteFile("ye_zi_api.log", s)
					//easygo.PanicError(err1)
					return
				}

				if rstMiddleMap.Code == YEZI_SUCCESS_CODE && nil != rstMiddleMap.Data {

					//转换map到数据库结构
					//设置队伍a的玩家详细信息
					if rstMiddleMap.Data.TeamAPlayers != nil {
						tempTeamAPlayers := make([]*share_message.APIPlayerDetail, 0)
						for _, teamAPlayer := range rstMiddleMap.Data.TeamAPlayers {

							tempTeamAPlayers = append(tempTeamAPlayers, teamAPlayer)
						}
						resultDB.Data.ApiTeamAPlayers = tempTeamAPlayers
					}

					//设置队伍b的玩家详细信息
					if rstMiddleMap.Data.TeamBPlayers != nil {
						tempTeamBPlayers := make([]*share_message.APIPlayerDetail, 0)
						for _, teamBPlayer := range rstMiddleMap.Data.TeamBPlayers {

							tempTeamBPlayers = append(tempTeamBPlayers, teamBPlayer)
						}
						resultDB.Data.ApiTeamBPlayers = tempTeamBPlayers
					}

					//设置直播信号源
					if rstMiddleMap.Data.LiveURL != nil {
						tempLiveURLs := make([]*share_message.APILiveURL, 0)
						for _, liveUrl := range rstMiddleMap.Data.LiveURL {
							var tempUrl = share_message.APILiveURL{}
							tempUrl.Name = easygo.NewString(liveUrl["name"])
							tempUrl.Url = easygo.NewString(liveUrl["url"])
							tempUrl.UrlH5 = easygo.NewString(liveUrl["url_h5"])
							tempUrl.NameH5 = easygo.NewString(liveUrl["name_h5"])
							tempLiveURLs = append(tempLiveURLs, &tempUrl)
						}
						resultDB.Data.ApiLiveUrls = tempLiveURLs
					}

					//设置两队历史交锋
					if rstMiddleMap.Data.WinProbability != nil {

						twoTeams := rstMiddleMap.Data.WinProbability["this_two_team"]
						allS := rstMiddleMap.Data.WinProbability["all"]
						dbTwoTeams := make([]*share_message.APIWinPBObject, 0)
						dbAllS := make([]*share_message.APIWinPBObject, 0)
						if twoTeams != nil {
							for twoKey, twoValue := range twoTeams {
								dbTwoTeams = append(dbTwoTeams, &share_message.APIWinPBObject{
									TeamId:  easygo.NewString(twoKey),
									WinRate: easygo.NewString(twoValue),
								})
							}
						}

						if allS != nil {
							for allKey, allValue := range allS {
								dbAllS = append(dbAllS, &share_message.APIWinPBObject{
									TeamId:  easygo.NewString(allKey),
									WinRate: easygo.NewString(allValue),
								})
							}
						}
						resultDB.Data.ApiWinProbability = &share_message.APIWinPB{
							ThisTwoTeam: dbTwoTeams,
							All:         dbAllS,
						}
					}
				}
			}
		}

		//判断转化后结构处理业务
		if resultDB.Code == YEZI_SUCCESS_CODE && resultDB.Data != nil {

			gameDetailData := resultDB.Data
			gameDetailQuery := &share_message.TableESPortsGameDetail{}
			//通过条件查询存在更新,不存在就插入
			errQuery := colDetail.Find(bson.M{"app_label_id": appLabelId,
				"game_id":    gameId,
				"api_origin": apiOrigin}).One(gameDetailQuery)

			if errQuery != nil && errQuery != mgo.ErrNotFound {
				logs.Error(errQuery)
				s := fmt.Sprintf("====处理比赛详情时查询失败===,查询条件:app_label_id:%v,===game_id:%v,===api_origin:%v",
					appLabelId, gameId, apiOrigin)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}

			//既然比赛状态修改的时候没有回调、比赛结束才改状态、接的时候就屏蔽掉进行中
			if gameDetailData.GetGameStatus() == for_game.GAME_STATUS_1 {
				gameDetailData.GameStatus = easygo.NewString(for_game.GAME_STATUS_0)
			}

			//本项目的id
			gameDetailData.AppLabelId = easygo.NewInt32(appLabelId)
			//本项目的名称
			gameDetailData.AppLabelName = easygo.NewString(for_game.LabelToESportNameMap[appLabelId])

			//接口来源id
			gameDetailData.ApiOrigin = easygo.NewInt32(apiOrigin)
			//接口来源名
			gameDetailData.ApiOriginName = easygo.NewString(for_game.ApiOriginIdToNameMap[apiOrigin])

			//没有就插入
			if errQuery == mgo.ErrNotFound {
				//新增数据取得自增id
				gameDetailData.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_GAME_DETAIL))
				//创建时间和更新时间
				nowTime := time.Now().Unix()
				gameDetailData.CreateTime = easygo.NewInt64(nowTime)
				gameDetailData.UpdateTime = easygo.NewInt64(nowTime)

				//新增的时候有直播推流地址设置到详情中
				if gameDetailData.GetHaveLive() == for_game.GAME_HAVE_LIVE_1 {
					gameDetailData.LivePaths = GetGameLiveUrl(gameId, eventId)
				}

				errIns := colDetail.Insert(gameDetailData)
				if errIns != nil {
					logs.Error(errIns)
					s := fmt.Sprintf("=====处理比赛详情时插入数据失败====插入相关信息:app_label_id:%v,===game_id:%v,===api_origin:%v",
						appLabelId, gameId, apiOrigin)
					for_game.WriteFile("ye_zi_api.log", s)
					return
				}

				//处理列表中对应的该场比赛、放在回调外防止服务器启动比赛详情数据有变更
				//赔率表中的比赛时间、比赛状态在处理比赛时处理
				dealSportGame(gameDetailData, col, colGuess, colUseRoll, updateTime, initCallBackFlag)

				//两队历史交锋、两队胜败统计、两队天敌克制统计(统计一次即可)
				dealTeamHistory(appLabelId, apiOrigin, gameId, eventId, colHistory)
			}

			if errQuery == nil {

				gameDetailData.GameStatusTime = easygo.NewInt64(gameDetailQuery.GetGameStatusTime())

				//回调时候设置、初始化的时候没有就不设置
				if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_2 {

					if gameDetailQuery.GetGameStatus() == for_game.GAME_STATUS_0 &&
						gameDetailData.GetGameStatus() == for_game.GAME_STATUS_2 {

						if updateTime != 0 {
							gameDetailData.GameStatusTime = easygo.NewInt64(updateTime)
						} else {
							gameDetailData.GameStatusTime = easygo.NewInt64(time.Now().Unix())
						}
					}
				}

				//新的更新时间
				gameDetailData.UpdateTime = easygo.NewInt64(time.Now().Unix())

				if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {
					//更新时有直播推流地址设置到详情
					if gameDetailData.GetHaveLive() == for_game.GAME_HAVE_LIVE_1 {

						gameDetailData.LivePaths = GetGameLiveUrl(gameId, eventId)
					}
				} else {
					//更新时有直播推流地址设置到详情
					if gameDetailQuery.GetHaveLive() == for_game.GAME_HAVE_LIVE_0 &&
						gameDetailData.GetHaveLive() == for_game.GAME_HAVE_LIVE_1 {

						gameDetailData.LivePaths = GetGameLiveUrl(gameId, eventId)
					}
				}

				errUpd := colDetail.Update(bson.M{"_id": gameDetailQuery.GetId()},
					bson.M{"$set": gameDetailData})

				if errUpd != nil {
					logs.Error(errUpd)
					s := fmt.Sprintf("处理比赛详情时更新失败=====更新唯一_id:%v", gameDetailQuery.GetId())
					for_game.WriteFile("ye_zi_api.log", s)
					return
				}

				//处理列表中对应的该场比赛、放在回调外防止服务器启动比赛详情数据有变更
				uniqueGameId := dealSportGame(gameDetailData, col, colGuess, colUseRoll, updateTime, initCallBackFlag)

				if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {

					//两队历史交锋、两队胜败统计、两队天敌克制统计(统计一次即可)
					dealTeamHistory(appLabelId, apiOrigin, gameId, eventId, colHistory)

				} else {
					//回调的时候更新的时候 判断一下查询的数据和api过来的队伍是否变化、若有队伍变化就调用历史相关
					if gameDetailData.GetTeamA() != gameDetailQuery.GetTeamA() ||
						gameDetailData.GetTeamB() != gameDetailQuery.GetTeamB() ||
						(gameDetailData.GetTeamAInfo() != nil && gameDetailQuery.GetTeamAInfo() != nil &&
							(gameDetailData.GetTeamAInfo().GetName() != gameDetailQuery.GetTeamAInfo().GetName() ||
								gameDetailData.GetTeamAInfo().GetNameEn() != gameDetailQuery.GetTeamAInfo().GetNameEn() ||
								gameDetailData.GetTeamAInfo().GetIcon() != gameDetailQuery.GetTeamAInfo().GetIcon())) ||
						(gameDetailData.GetTeamBInfo() != nil && gameDetailQuery.GetTeamBInfo() != nil &&
							(gameDetailData.GetTeamBInfo().GetName() != gameDetailQuery.GetTeamBInfo().GetName() ||
								gameDetailData.GetTeamBInfo().GetNameEn() != gameDetailQuery.GetTeamBInfo().GetNameEn() ||
								gameDetailData.GetTeamBInfo().GetIcon() != gameDetailQuery.GetTeamBInfo().GetIcon())) {

						//两队历史交锋、两队胜败统计、两队天敌克制统计(统计一次即可)
						dealTeamHistory(appLabelId, apiOrigin, gameId, eventId, colHistory)
					}
				}

				//回调的时候处理(初始化不能保证服务同时都启动、初始化在初始化的时候延时开奖)
				if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_2 {

					//通知开奖
					//状态从0变成2已经结束
					if gameDetailQuery.GetGameStatus() == for_game.GAME_STATUS_0 &&
						gameDetailData.GetGameStatus() == for_game.GAME_STATUS_2 {

						labelId := gameDetailQuery.GetAppLabelId()
						//通知各个模块开奖
						req := &client_hall.LotteryRequest{
							UniqueGameId: easygo.NewInt64(uniqueGameId),
						}

						s := fmt.Sprintf("=======%v请求开奖服发送  开始==========", for_game.LabelToESportNameMap[labelId])
						for_game.WriteFile("ye_zi_api.log", s)

						if labelId == for_game.ESPORTS_LABEL_WZRY {

							dal_common.SendMsgToIdOtherServer(PServerInfoMgr,
								for_game.SERVER_TYPE_SPORT_LOTTERY_WZRY,
								"RpcESportLottery",
								req)

						} else if labelId == for_game.ESPORTS_LABEL_DOTA2 {

							dal_common.SendMsgToIdOtherServer(
								PServerInfoMgr,
								for_game.SERVER_TYPE_SPORT_LOTTERY_DOTA,
								"RpcESportLottery",
								req)

						} else if labelId == for_game.ESPORTS_LABEL_LOL {

							dal_common.SendMsgToIdOtherServer(
								PServerInfoMgr,
								for_game.SERVER_TYPE_SPORT_LOTTERY_LOL,
								"RpcESportLottery",
								req)

						} else if labelId == for_game.ESPORTS_LABEL_CSGO {

							dal_common.SendMsgToIdOtherServer(
								PServerInfoMgr,
								for_game.SERVER_TYPE_SPORT_LOTTERY_CSGO,
								"RpcESportLottery",
								req)

						}

						s = fmt.Sprintf("=======%v请求开奖服发送   结束==========", for_game.LabelToESportNameMap[labelId])
						for_game.WriteFile("ye_zi_api.log", s)
					}
				}
			}
		}
	}

	if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {
		s := fmt.Sprintf("=======dealGameDetail初始化处理比赛详情   结束==========")
		for_game.WriteFile("ye_zi_api.log", s)
	} else {
		s := fmt.Sprintf("=======dealGameDetail回调处理比赛详情    结束==========")
		for_game.WriteFile("ye_zi_api.log", s)
	}
}

//未开始比赛列表初始化
func GetYeZiAllGames(eventId string) {

	//配置的不是野子科技的项目参数传入的eventId
	if !ISAppLabel(eventId) {
		s := fmt.Sprintf("================传入的eventId:%v不在游戏配置列表中、直接返回=================", eventId)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	//取得相关label
	appLabelId := for_game.EventIdToESportLabelMap[eventId]
	s := fmt.Sprintf("================GetYeZiAllGames初始化%v未开始列表==================开始", for_game.LabelToESportNameMap[appLabelId])
	for_game.WriteFile("ye_zi_api.log", s)

	var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

	param := url.Values{}
	param.Set("app_secret", yeZiAppSecret)
	param.Set("app_key", yeZiAppKey)
	param.Set("timestamp", timeStr)
	param.Set("resource", "game")
	param.Set("func", "lists")
	param.Set("event_id", eventId)

	sign := GetYeZiApiSign(param)

	param.Del("app_secret")
	param.Add("sign", sign)

	//发送请求
	data, errGet := Get(yeZiUrl, param)
	if errGet != nil {

		logs.Error(errGet)
		s := fmt.Sprintf("======野子科技获取%v未开始比赛列表api接口请求失败=====,错误信息:%v,请求参数为::::param=%v",
			for_game.LabelToESportNameMap[appLabelId], errGet, param)
		for_game.WriteFile("ye_zi_api.log", s)
		//easygo.PanicError(errGet)
		return
	} else {
		result := YeZiESPortsGamesInfo{}
		if "" == data {
			s := fmt.Sprintf("========野子科技获取%v未开始比赛列表api接口返回数据为空=======,请求参数为::::param=%v",
				for_game.LabelToESportNameMap[appLabelId], param)
			for_game.WriteFile("ye_zi_api.log", s)
			return
		} else {
			//列表数据就初始化取得一次、正式环境就不需要记录了
			if !for_game.IS_FORMAL_SERVER {
				s := fmt.Sprintf("=============野子科技获取%v未开始比赛列表api接口返回的数据:data=%v=============,请求参数为::::param=%v",
					for_game.LabelToESportNameMap[appLabelId], data, param)
				for_game.WriteFile("ye_zi_api.log", s)
			}

			err := json.Unmarshal([]byte(data), &result)

			if err != nil {
				logs.Error(err)
				s = fmt.Sprintf("======野子科技获取%v未开始比赛列表api接口返回的数据后json转结构体出错::::data=%v======,请求参数为::::param=%v",
					for_game.LabelToESportNameMap[appLabelId], data, param)
				for_game.WriteFile("ye_zi_api.log", s)
				//easygo.PanicError(err)
				return
			}

			if result.ErrorCode != 0 && result.ErrorMsg != "" {
				s = fmt.Sprintf("=====野子科技获取%v未开始比赛列表api接口请求返回错误信息:错误码:%v,错误信息:%v======,请求参数为::::param=%v",
					for_game.LabelToESportNameMap[appLabelId], result.ErrorCode, result.ErrorMsg, param)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}
		}

		//判断转化后结构处理业务
		if result.Code == YEZI_SUCCESS_CODE {
			col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
			defer closeFun()
			if nil != result.Data && len(result.Data) > 0 {
				for i := 0; i < len(result.Data); i++ {
					sportGameData := result.Data[i]
					//目前只入库两组对战的比赛
					if nil != sportGameData && sportGameData.GetDimension() == for_game.GAME_DIMENSION_1 {
						sportGameQuery := share_message.TableESPortsGame{}
						//通过条件查询存在更新,不存在就插入
						errQuery := col.Find(bson.M{"app_label_id": appLabelId,
							"game_id":    sportGameData.GetGameId(),
							"api_origin": for_game.ESPORTS_API_ORIGIN_ID_YEZI}).One(&sportGameQuery)

						if errQuery != nil && errQuery != mgo.ErrNotFound {
							logs.Error(errQuery)
							s := fmt.Sprintf("========处理%v未开始列表查询比赛表数据失败=======,查询条件:app_label_id:%v,===game_id:%v,====api_origin:%v",
								for_game.LabelToESportNameMap[appLabelId], appLabelId, sportGameData.GetGameId(), for_game.ESPORTS_API_ORIGIN_ID_YEZI)
							for_game.WriteFile("ye_zi_api.log", s)

							continue
						}

						//没有就插入
						if errQuery == mgo.ErrNotFound {
							//新增数据取得自增id
							sportGameData.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_GAME))

							//本项目的id
							sportGameData.AppLabelId = easygo.NewInt32(appLabelId)
							//本项目的名称
							sportGameData.AppLabelName = easygo.NewString(for_game.LabelToESportNameMap[appLabelId])

							//接口来源id
							sportGameData.ApiOrigin = easygo.NewInt32(for_game.ESPORTS_API_ORIGIN_ID_YEZI)
							//接口来源名
							sportGameData.ApiOriginName = easygo.NewString(for_game.ApiOriginIdToNameMap[for_game.ESPORTS_API_ORIGIN_ID_YEZI])

							//发布未发布状态
							sportGameData.ReleaseFlag = easygo.NewInt32(for_game.GAME_RELEASE_FLAG_1)

							//开奖状态
							sportGameData.IsLottery = easygo.NewInt32(for_game.GAME_IS_LOTTERY_0)
							//创建时间和更新时间
							nowTime := time.Now().Unix()
							sportGameData.CreateTime = easygo.NewInt64(nowTime)
							sportGameData.UpdateTime = easygo.NewInt64(nowTime)

							//把比赛开始时间转成int64存到数据库
							sportGameData.BeginTimeInt = easygo.NewInt64(for_game.GetGameTimeStrToInt64(sportGameData.GetBeginTime()))

							errIns := col.Insert(sportGameData)
							if errIns != nil {
								logs.Error(errIns)
								s := fmt.Sprintf("========处理%v未开始列表插入比赛表数据失败=======,插入相关信息:app_label_id:%v,===game_id:%v,====api_origin:%v",
									for_game.LabelToESportNameMap[appLabelId], appLabelId, sportGameData.GetGameId(), for_game.ESPORTS_API_ORIGIN_ID_YEZI)
								for_game.WriteFile("ye_zi_api.log", s)

								continue
							}
						}

						//存在就更新 只设置api过来的数据
						if errQuery == nil {

							//新的更新时间
							sportGameData.UpdateTime = easygo.NewInt64(time.Now().Unix())
							//把比赛开始时间转成int64存到数据库
							sportGameData.BeginTimeInt = easygo.NewInt64(for_game.GetGameTimeStrToInt64(sportGameData.GetBeginTime()))

							errUpd := col.Update(bson.M{"_id": sportGameQuery.GetId()},
								bson.M{"$set": sportGameData})

							if errUpd != nil {
								logs.Error(errUpd)
								s := fmt.Sprintf("处理%v未开始列表数据更新比赛表数据失败=====更新唯一_id:%v",
									for_game.LabelToESportNameMap[appLabelId], sportGameQuery.GetId())
								for_game.WriteFile("ye_zi_api.log", s)

								continue
							}
						}
					}
				}
			}
		}
	}

	s1 := fmt.Sprintf("================GetYeZiAllGames初始化%v未开始列表==================结束", for_game.LabelToESportNameMap[appLabelId])
	for_game.WriteFile("ye_zi_api.log", s1)
}

//回调处理比赛详情
func CallBackDealGameDetail(appLabelId int32, apiOrigin int32, gameId string, eventId string, updateTime int64) {
	s := fmt.Sprintf("=======CallBackDealGameDetail回调处理比赛详情  开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v,updateTime=%v",
		appLabelId, apiOrigin, gameId, eventId, updateTime)
	for_game.WriteFile("ye_zi_api.log", s)

	//判断参数就可以、该回调的时候可能没有比赛
	if appLabelId == 0 || apiOrigin == 0 || gameId == "" || eventId == "" {
		s := fmt.Sprintf("=======CallBackDealGameDetail回调处理处理比赛详情参数错误==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
			appLabelId, apiOrigin, gameId, eventId)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	colDetail, closeFunDetail := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_DETAIL)
	defer closeFunDetail()

	colUseRoll, closeFunUseRoll := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_USE_ROLL_GUESS)
	defer closeFunUseRoll()

	colGuess, closeFunGuess := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_GUESS)
	defer closeFunGuess()

	colHistory, closeFunHistory := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_TEAM_HIS_INFO)
	defer closeFunHistory()

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
	defer closeFun()

	//处理比赛详情
	dealGameDetail(appLabelId, apiOrigin, gameId, eventId,
		colDetail, colUseRoll, colGuess, colHistory, col,
		updateTime, for_game.INIT_CALLBACK_FLAG_2)

	s = fmt.Sprintf("=======CallBackDealGameDetail回调处理比赛详情    结束==========")
	for_game.WriteFile("ye_zi_api.log", s)
}

//回调处理早盘
func CallBackDealGameGuessMorn(appLabelId int32, apiOrigin int32, gameId string, eventId string, updateTime int64) {

	s := fmt.Sprintf("========CallBackDealGameGuessMorn回调处理处理比赛动态信息表(早盘)   开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
		appLabelId, apiOrigin, gameId, eventId)
	for_game.WriteFile("ye_zi_api.log", s)

	//判断参数
	if appLabelId == 0 || apiOrigin == 0 || gameId == "" || eventId == "" {
		s := fmt.Sprintf("=======CallBackDealGameGuessMorn回调处理处理早盘时参数错误==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
			appLabelId, apiOrigin, gameId, eventId)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	//判断有无比赛数据
	gameObject := for_game.GetRedisGameDetailHeadGroup(apiOrigin, appLabelId, gameId)
	if gameObject == nil {
		s := fmt.Sprintf("=======CallBackDealGameGuessMorn回调处理处理早盘时缺少比赛信息==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v",
			appLabelId, apiOrigin, gameId)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	colGuess, closeFunGuess := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_GUESS)
	defer closeFunGuess()
	//处理早盘
	dealGameGuessMorn(appLabelId, apiOrigin, gameId, eventId, colGuess, updateTime, for_game.INIT_CALLBACK_FLAG_2)

	if gameObject != nil {
		//设置redis早盘信息
		for_game.SetRedisGuessMornDetail(gameObject.GetUniqueGameId(), appLabelId, gameId, apiOrigin)
	}

	s = fmt.Sprintf("========CallBackDealGameGuessMorn回调处理处理比赛动态信息表(早盘)   结束==========")
	for_game.WriteFile("ye_zi_api.log", s)
}

//回调处理使用滚盘
func CallBackUseGameGuessRoll(appLabelId int32, apiOrigin int32, gameId string, eventId string) {

	s := fmt.Sprintf("========CallBackUseGameGuessRoll回调处理处理使用滚盘   开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
		appLabelId, apiOrigin, gameId, eventId)
	for_game.WriteFile("ye_zi_api.log", s)

	//判断参数
	if appLabelId == 0 || apiOrigin == 0 || gameId == "" || eventId == "" {
		s := fmt.Sprintf("=======CallBackUseGameGuessRoll回调处理处理使用滚盘时参数错误==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
			appLabelId, apiOrigin, gameId, eventId)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	//判断有无比赛数据
	gameObject := for_game.GetRedisGameDetailHeadGroup(apiOrigin, appLabelId, gameId)
	if gameObject == nil {
		s := fmt.Sprintf("=======CallBackUseGameGuessRoll回调处理处理使用滚盘时缺少比赛信息==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v",
			appLabelId, apiOrigin, gameId)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	colUseRoll, closeFunUseRoll := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_USE_ROLL_GUESS)
	defer closeFunUseRoll()
	//处理使用滚盘
	UseGameGuessRoll(appLabelId, apiOrigin, gameId, eventId, colUseRoll)

	s = fmt.Sprintf("========CallBackUseGameGuessRoll回调处理处理使用滚盘   结束==========")
	for_game.WriteFile("ye_zi_api.log", s)
}

//回调处理调取滚盘数据
func CallBackDealGameGuessRoll(appLabelId int32, apiOrigin int32, gameId string, eventId string, updateTime int64) {

	s := fmt.Sprintf("=======CallBackDealGameGuessRoll回调处理比赛动态信息表(滚盘)   开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v,updateTime=%v",
		appLabelId, apiOrigin, gameId, eventId, updateTime)
	for_game.WriteFile("ye_zi_api.log", s)

	//判断参数
	if appLabelId == 0 || apiOrigin == 0 || gameId == "" || eventId == "" {
		s := fmt.Sprintf("=======CallBackDealGameGuessRoll回调处理比赛动态信息表(滚盘)时参数错误==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
			appLabelId, apiOrigin, gameId, eventId)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	//判断有无比赛数据
	gameObject := for_game.GetRedisGameDetailHeadGroup(apiOrigin, appLabelId, gameId)
	if gameObject == nil {
		s := fmt.Sprintf("=======CallBackDealGameGuessRoll回调处理比赛动态信息表(滚盘)时缺少比赛信息==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v",
			appLabelId, apiOrigin, gameId)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	colGuess, closeFunGuess := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_GUESS)
	defer closeFunGuess()
	//处理滚盘
	dealGameGuessRoll(appLabelId, apiOrigin, gameId, eventId, colGuess, updateTime, for_game.INIT_CALLBACK_FLAG_2)

	if gameObject != nil {
		//设置redis滚盘信息
		for_game.SetRedisGuessRollDetail(gameObject.GetUniqueGameId(), appLabelId, gameId, apiOrigin)
	}

	s = fmt.Sprintf("=======CallBackDealGameGuessRoll回调处理比赛动态信息表(滚盘)   结束==========")
	for_game.WriteFile("ye_zi_api.log", s)
}

//冲正回调
func CallBackDealReversal(appLabelId int32,
	apiOrigin int32,
	gameId string,
	eventId string,
	updateTime int64,
	betId int32) {

	s := fmt.Sprintf("=======CallBackDealReversal回调处理冲正信息   开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v,updateTime=%v,betId=%v",
		appLabelId, apiOrigin, gameId, eventId, updateTime, betId)
	for_game.WriteFile("ye_zi_api.log", s)

	//判断参数
	if appLabelId == 0 || apiOrigin == 0 || gameId == "" || eventId == "" || betId == 0 {
		s := fmt.Sprintf("=======CallBackDealReversal回调处理冲正信息时参数错误==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v,betId=%v",
			appLabelId, apiOrigin, gameId, eventId, betId)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	//判断有无比赛数据
	gameObject := for_game.GetRedisGameDetailHeadGroup(apiOrigin, appLabelId, gameId)
	if gameObject == nil {
		s := fmt.Sprintf("=======CallBackDealReversal回调处理冲正信息时缺少比赛信息==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v",
			appLabelId, apiOrigin, gameId)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	timeStr := strconv.FormatInt(time.Now().Unix(), 10)

	param := url.Values{}
	param.Set("app_secret", yeZiAppSecret)
	param.Set("app_key", yeZiAppKey)
	param.Set("timestamp", timeStr)
	param.Set("resource", "dynamic")
	param.Set("func", "gameBetInfo")
	param.Set("event_id", eventId)
	param.Set("game_id", gameId)
	param.Set("bet_id", string(betId))

	sign := GetYeZiApiSign(param)

	param.Del("app_secret")
	param.Add("sign", sign)

	//发送请求
	data, errGet := Get(yeZiUrl, param)
	if errGet != nil {

		logs.Error(errGet)
		s := fmt.Sprintf("=====野子科技冲正回调api接口请求失败,错误信息:%v,请求参数为::::param=%v", errGet, param)
		for_game.WriteFile("ye_zi_api.log", s)
		//easygo.PanicError(errGet)
		return
	} else {
		resultDB := YeZiESPortsGameBetInfo{}
		if "" == data {
			s := fmt.Sprintf("==========野子科技冲正回调api接口请求返回数据为空=======,请求参数为::::param=%v", param)
			for_game.WriteFile("ye_zi_api.log", s)
			return
		} else {
			//s := fmt.Sprintf("=============野子科技冲正回调api接口请求返回数据:data=%v=============,请求参数为::::param=%v", data, param)
			//for_game.WriteFile("ye_zi_api.log", s)
			err := json.Unmarshal([]byte(data), &resultDB)

			if err != nil {
				logs.Error(err)
				s = fmt.Sprintf("======野子科技冲正回调api接口请求返回数据后json转结构体出错::::data=%v======,请求参数为::::param=%v",
					data, param)
				for_game.WriteFile("ye_zi_api.log", s)
				//easygo.PanicError(err)
				return
			}

			if resultDB.ErrorCode != 0 && resultDB.ErrorMsg != "" {
				s := fmt.Sprintf("======野子科技冲正回调接口请求返回错误信息:错误码:%v,错误信息:%v=======,,请求参数为::::param=%v",
					resultDB.ErrorCode, resultDB.ErrorMsg, param)
				logs.Error(s)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}
		}

		//判断转化后结构处理业务
		if resultDB.Code == YEZI_SUCCESS_CODE && resultDB.Data != nil {

			apiData := resultDB.Data

			//查询早盘、滚盘的数据
			colGuess, closeFunGuess := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_GUESS)
			defer closeFunGuess()

			gameGuessQuery := make([]*share_message.TableESPortsGameGuess, 0)
			//查询所有该bet_id下的记录(可能有早盘和滚盘数据)
			errQuery := colGuess.Find(bson.M{"app_label_id": appLabelId,
				"game_id":    gameId,
				"api_origin": apiOrigin}).All(&gameGuessQuery)

			if errQuery != nil && errQuery != mgo.ErrNotFound {
				logs.Error(errQuery)
				s := fmt.Sprintf("=====处理冲正回调时查询竞猜数据失败=====查询条件为app_label_id:%v,===game_id:%v,===api_origin:%v",
					appLabelId, gameId, apiOrigin)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}

			if errQuery == mgo.ErrNotFound || gameGuessQuery == nil || len(gameGuessQuery) <= 0 {
				s := fmt.Sprintf("=====处理冲正回调时查询竞猜数据为空=====查询条件为app_label_id:%v,===game_id:%v,===api_origin:%v",
					appLabelId, gameId, apiOrigin)
				for_game.WriteFile("ye_zi_api.log", s)
				return
			}

			if errQuery == nil && nil != gameGuessQuery && len(gameGuessQuery) > 0 {

				//同时处理早盘滚盘
				for _, value := range gameGuessQuery {
					guessInfos := value.GetGuess()

					if nil != guessInfos && len(guessInfos) > 0 {

						for _, guessValue := range guessInfos {
							if guessValue.GetBetId() == apiData.GetBetId() {

								//num不要设置
								guessValue.BetId = easygo.NewString(apiData.GetBetId())
								guessValue.BetTitle = easygo.NewString(apiData.GetBetTitle())
								guessValue.BetTitleEn = easygo.NewString(apiData.GetBetTitleEn())
								guessValue.RiskLevel = easygo.NewString(apiData.GetRiskLevel())

								//设置具体投注项目
								itemInfos := guessValue.GetItems()
								apiItemInfos := apiData.GetItems()
								if nil != itemInfos && len(itemInfos) > 0 && nil != apiItemInfos && len(apiItemInfos) > 0 {
									for _, itemValue := range itemInfos {
										for _, apiItemValue := range apiItemInfos {
											if itemValue.GetBetNum() == apiItemValue.GetBetNum() {
												itemValue.BetNum = easygo.NewString(apiItemValue.GetBetNum())
												itemValue.TeamId = easygo.NewString(apiItemValue.GetTeamId())
												itemValue.Status = easygo.NewString(apiItemValue.GetStatus())
												itemValue.PlayerId = easygo.NewString(apiItemValue.GetPlayerId())
												itemValue.Win = easygo.NewString(apiItemValue.GetWin())
												itemValue.Odds = easygo.NewString(apiItemValue.GetOdds())
												itemValue.BetType = easygo.NewString(apiItemValue.GetBetType())
												itemValue.OddsName = easygo.NewString(apiItemValue.GetOddsName())
												itemValue.LimitBet = easygo.NewInt32(apiItemValue.GetLimitBet())
												itemValue.CanCustom = easygo.NewInt32(apiItemValue.GetCanCustom())
												itemValue.OddsStatus = easygo.NewString(apiItemValue.GetOddsStatus())
												itemValue.GroupId = easygo.NewString(apiItemValue.GetGroupId())
												itemValue.GroupValue = easygo.NewString(apiItemValue.GetGroupValue())
												itemValue.GroupFlag = easygo.NewString(apiItemValue.GetGroupFlag())
												itemValue.BetStart = easygo.NewString(apiItemValue.GetBetStart())
												itemValue.BetOver = easygo.NewString(apiItemValue.GetBetOver())
												itemValue.OddsTime = easygo.NewString(apiItemValue.GetOddsTime())

												//记录封盘、结果产生时间
												//1、封盘时间
												//回调才设置、投注状态从1到0  或  从1到3记录; 若3到1的时候不变 3到0的时候不变
												// (滚盘开奖的时候这个时间的前300秒投注为无效单)
												if itemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_1 &&
													apiItemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_0 {
													if updateTime != 0 {
														itemValue.StatusTime = easygo.NewInt64(updateTime)
													} else {
														itemValue.StatusTime = easygo.NewInt64(time.Now().Unix())
													}
												} else if itemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_1 &&
													apiItemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_3 {
													if updateTime != 0 {
														itemValue.StatusTime = easygo.NewInt64(updateTime)
													} else {
														itemValue.StatusTime = easygo.NewInt64(time.Now().Unix())
													}
												}

												//2、结果产生时间
												if (itemValue.GetStatus() == for_game.GAME_GUESS_ITEM_STATUS_0 ||
													itemValue.GetWin() == for_game.GAME_GUESS_ITEM_WIN_NORST) &&
													(apiItemValue.GetStatus() == for_game.GAME_GUESS_ITEM_STATUS_1 &&
														apiItemValue.GetWin() != for_game.GAME_GUESS_ITEM_WIN_NORST) {
													if updateTime != 0 {
														itemValue.ResultTime = easygo.NewInt64(updateTime)
													} else {
														itemValue.ResultTime = easygo.NewInt64(time.Now().Unix())
													}
												}

												//跳出做下一个投注项目
												break
											}
										}
									}
								}

								//直接跳出、说明投注内容已经匹配完毕(对于同一场比赛来说、同一个盘口中BetId是唯一的)
								break
							}
						}
					}

					//更新数据库
					value.UpdateTime = easygo.NewInt64(time.Now().Unix())

					errUpd := colGuess.Update(bson.M{"_id": value.GetId()},
						bson.M{"$set": value})

					if errUpd != nil {
						logs.Error(errUpd)
						s := fmt.Sprintf("=======野子科技处理冲正回调时更新数据失败=====更新盘口的唯一的_id:%v",
							value.GetId())
						for_game.WriteFile("ye_zi_api.log", s)
						return
					}

					//设置redis相关早盘、滚盘数据
					if gameObject != nil {
						//早盘
						if value.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_1 {
							//设置redis早盘信息
							for_game.SetRedisGuessMornDetail(gameObject.GetUniqueGameId(), appLabelId, gameId, apiOrigin)
							//滚盘
						} else {
							//设置redis滚盘信息
							for_game.SetRedisGuessRollDetail(gameObject.GetUniqueGameId(), appLabelId, gameId, apiOrigin)
						}
					}
				}
			}
		}
	}

	s = fmt.Sprintf("=======CallBackDealReversal回调处理冲正信息   结束==========")
	for_game.WriteFile("ye_zi_api.log", s)
}

//通过eventId确定是否本app里面配置的项目
func ISAppLabel(eventId string) bool {
	appLabelId := for_game.EventIdToESportLabelMap[eventId]
	//配置的不是野子科技的项目参数传入的eventId
	if appLabelId == 0 {
		return false
	} else {
		return true
	}
}

//通过比赛详情对象得到某场比赛列表中的某条记录的对象
func dealSportGame(detail *share_message.TableESPortsGameDetail,
	col *mgo.Collection,
	colGuess *mgo.Collection,
	colUseRoll *mgo.Collection,
	updateTime int64,
	initCallBackFlag int32) int64 {
	if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {
		s := fmt.Sprintf("=======dealSportGam初始化比赛详情后、处理比赛表   开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
			detail.GetAppLabelId(), detail.GetApiOrigin(), detail.GetGameId(), detail.GetEventId())
		for_game.WriteFile("ye_zi_api.log", s)
	} else {
		s := fmt.Sprintf("=======dealSportGame回调比赛详情后、处理比赛表   开始==========:参数appLabelId=%v,apiOrigin=%v,gameId=%v,eventId=%v",
			detail.GetAppLabelId(), detail.GetApiOrigin(), detail.GetGameId(), detail.GetEventId())
		for_game.WriteFile("ye_zi_api.log", s)
	}

	var uniqueGameId int64

	appLabelId := for_game.EventIdToESportLabelMap[detail.GetEventId()]
	gameId := detail.GetGameId()
	apiOrigin := for_game.ESPORTS_API_ORIGIN_ID_YEZI

	sportGameQuery := share_message.TableESPortsGame{}
	//通过条件查询存在更新,不存在就插入
	errQuery := col.Find(bson.M{"app_label_id": appLabelId,
		"game_id":    gameId,
		"api_origin": apiOrigin}).One(&sportGameQuery)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("回调比赛详情后处理比赛中某场的比赛时、查询比赛表失败===查询条件:app_label_id:%v,===game_id:%v,====api_origin:%v",
			appLabelId, gameId, apiOrigin)
		for_game.WriteFile("ye_zi_api.log", s)

		return uniqueGameId
	}

	sportsGame := share_message.TableESPortsGame{
		AppLabelId:     easygo.NewInt32(appLabelId),
		AppLabelName:   easygo.NewString(for_game.LabelToESportNameMap[appLabelId]),
		ApiOrigin:      easygo.NewInt32(apiOrigin),
		ApiOriginName:  easygo.NewString(for_game.ApiOriginIdToNameMap[apiOrigin]),
		Dimension:      easygo.NewString(for_game.GAME_DIMENSION_1),
		FightName:      easygo.NewString(detail.GetFightName()),
		EventName:      easygo.NewString(detail.GetEventName()),
		EventNameEn:    easygo.NewString(detail.GetEventNameEn()),
		EventId:        easygo.NewString(detail.GetEventId()),
		MatchStage:     easygo.NewString(detail.GetMatchStage()),
		MatchStageId:   easygo.NewString(detail.GetMatchStageId()),
		MatchName:      easygo.NewString(detail.GetMatchName()),
		MatchNameEn:    easygo.NewString(detail.GetMatchNameEn()),
		MatchId:        easygo.NewString(detail.GetMatchId()),
		GameType:       easygo.NewString(detail.GetGameType()),
		ScoreA:         easygo.NewString(detail.GetScoreA()),
		ScoreB:         easygo.NewString(detail.GetScoreB()),
		IsLive:         easygo.NewInt32(detail.GetIsLive()),
		Bo:             easygo.NewString(detail.GetBo()),
		HotGame:        easygo.NewString(detail.GetHotGame()),
		IsBet:          easygo.NewInt32(detail.GetIsBet()),
		BeginTime:      easygo.NewString(detail.GetBeginTime()),
		GameId:         easygo.NewString(detail.GetGameId()),
		GameStatus:     easygo.NewString(detail.GetGameStatus()),
		GameStatusType: easygo.NewString(detail.GetGameStatusType()),
		HaveLive:       easygo.NewInt32(detail.GetHaveLive()),
		HaveRoll:       easygo.NewInt32(detail.GetHaveRoll()),
		BeginTimeInt:   easygo.NewInt64(for_game.GetGameTimeStrToInt64(detail.GetBeginTime())),
		OverTime:       easygo.NewString(detail.GetOverTime()),
	}

	var teamA share_message.ApiTeam
	var teamB share_message.ApiTeam
	if nil != detail.GetTeamAInfo() {
		teamA = share_message.ApiTeam{
			TeamId: easygo.NewString(detail.GetTeamAInfo().GetTeamId()),
			Name:   easygo.NewString(detail.GetTeamAInfo().GetName()),
			NameEn: easygo.NewString(detail.GetTeamAInfo().GetNameEn()),
			Icon:   easygo.NewString(detail.GetTeamAInfo().GetIcon()),
		}
	}

	if nil != detail.GetTeamBInfo() {
		teamB = share_message.ApiTeam{
			TeamId: easygo.NewString(detail.GetTeamBInfo().GetTeamId()),
			Name:   easygo.NewString(detail.GetTeamBInfo().GetName()),
			NameEn: easygo.NewString(detail.GetTeamBInfo().GetNameEn()),
			Icon:   easygo.NewString(detail.GetTeamBInfo().GetIcon()),
		}
	}

	var playersA []*share_message.ApiPlayer
	var playersB []*share_message.ApiPlayer
	if detail.GetApiTeamAPlayers() != nil && len(detail.GetApiTeamAPlayers()) > 0 {
		playersA = make([]*share_message.ApiPlayer, 0)
		for _, value := range detail.GetApiTeamAPlayers() {
			playerA := share_message.ApiPlayer{
				Sn:       easygo.NewString(value.GetSn()),
				PlayerId: easygo.NewString(value.GetPlayerId()),
				Name:     easygo.NewString(value.GetName()),
			}
			playersA = append(playersA, &playerA)
		}
	}

	if detail.GetApiTeamBPlayers() != nil && len(detail.GetApiTeamBPlayers()) > 0 {
		playersB = make([]*share_message.ApiPlayer, 0)
		for _, value := range detail.GetApiTeamBPlayers() {
			playerB := share_message.ApiPlayer{
				Sn:       easygo.NewString(value.GetSn()),
				PlayerId: easygo.NewString(value.GetPlayerId()),
				Name:     easygo.NewString(value.GetName()),
			}
			playersB = append(playersB, &playerB)
		}
	}

	sportsGame.TeamA = &teamA
	sportsGame.TeamB = &teamB
	sportsGame.PlayerA = playersA
	sportsGame.PlayerB = playersB

	//没有就插入
	if errQuery == mgo.ErrNotFound {
		//新增数据取得自增id
		sportsGame.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_GAME))

		//发布未发布状态
		sportsGame.ReleaseFlag = easygo.NewInt32(for_game.GAME_RELEASE_FLAG_1)

		//开奖状态
		sportsGame.IsLottery = easygo.NewInt32(for_game.GAME_IS_LOTTERY_0)

		//创建时间和更新时间
		nowTime := time.Now().Unix()
		sportsGame.CreateTime = easygo.NewInt64(nowTime)
		sportsGame.UpdateTime = easygo.NewInt64(nowTime)

		uniqueGameId = sportsGame.GetId()

		errIns := col.Insert(sportsGame)

		if errIns != nil {
			logs.Error(errIns)
			s := fmt.Sprintf("回调比赛详情后处理比赛中某场的比赛时、插入比赛表失败===插入相关信息:app_label_id:%v,===game_id:%v,====api_origin:%v",
				appLabelId, gameId, for_game.ESPORTS_API_ORIGIN_ID_YEZI)
			for_game.WriteFile("ye_zi_api.log", s)

			return uniqueGameId
		}

		//设置redis
		for_game.SetRedisGameDetailHead(uniqueGameId)

		//新增的时候调取下盘口信息
		//早盘
		if sportsGame.GetIsBet() == for_game.GAME_IS_BET_1 {
			dealGameGuessMorn(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colGuess, updateTime, initCallBackFlag)
			//设置redis早盘信息
			for_game.SetRedisGuessMornDetail(uniqueGameId, appLabelId, gameId, apiOrigin)
		}

		//滚盘
		if sportsGame.GetHaveRoll() == for_game.GAME_HAVE_ROLL_1 {
			//授权
			UseGameGuessRoll(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colUseRoll)

			//取得滚盘
			dealGameGuessRoll(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colGuess, updateTime, initCallBackFlag)

			//设置redis滚盘信息
			for_game.SetRedisGuessRollDetail(uniqueGameId, appLabelId, gameId, apiOrigin)
		}
	}

	//存在就更新
	if errQuery == nil {

		//新的更新时间
		sportsGame.UpdateTime = easygo.NewInt64(time.Now().Unix())
		//把比赛结束时间记录下来开奖用
		sportsGame.GameStatusTime = easygo.NewInt64(detail.GetGameStatusTime())

		uniqueGameId = sportGameQuery.GetId()
		errUpd := col.Update(bson.M{"_id": sportGameQuery.GetId()},
			bson.M{"$set": sportsGame})

		if errUpd != nil {
			logs.Error(errUpd)
			s := fmt.Sprintf("回调比赛详情后处理比赛中某场的比赛时、更新比赛表失败=====更新比赛表的唯一_id:%v",
				sportGameQuery.GetId())
			for_game.WriteFile("ye_zi_api.log", s)

			return uniqueGameId
		}

		//设置redis
		for_game.SetRedisGameDetailHead(sportGameQuery.GetId())

		//初始化
		if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {
			//早盘
			//if sportsGame.GetIsBet() == for_game.GAME_IS_BET_1 {
			//	dealGameGuessMorn(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colGuess, updateTime, initCallBackFlag)
			//	//设置redis早盘信息
			//	for_game.SetRedisGuessMornDetail(uniqueGameId, appLabelId, gameId, apiOrigin)
			//}
			if sportsGame.GetIsBet() == for_game.GAME_IS_BET_1 {
				//早盘
				//如果是初始化存在比赛、说明可能存在回调没有做、不判断条件都拉取
				dealGameGuessMorn(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colGuess, updateTime, initCallBackFlag)
				//设置redis早盘信息
				for_game.SetRedisGuessMornDetail(uniqueGameId, appLabelId, gameId, apiOrigin)
			}
			////滚盘
			//if sportsGame.GetHaveRoll() == for_game.GAME_HAVE_ROLL_1 {
			//	//授权
			//	UseGameGuessRoll(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colUseRoll)
			//
			//	//取得滚盘
			//	dealGameGuessRoll(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colGuess, updateTime, initCallBackFlag)
			//
			//	//设置redis滚盘信息
			//	for_game.SetRedisGuessRollDetail(uniqueGameId, appLabelId, gameId, apiOrigin)
			//}

			if sportsGame.GetHaveRoll() == for_game.GAME_HAVE_ROLL_1 {
				//滚盘
				//如果是初始化存在比赛、说明可能存在回调没有做、不判断条件都拉取
				//授权
				UseGameGuessRoll(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colUseRoll)

				//取得滚盘
				dealGameGuessRoll(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colGuess, updateTime, initCallBackFlag)

				//设置redis滚盘信息
				for_game.SetRedisGuessRollDetail(uniqueGameId, appLabelId, gameId, apiOrigin)
			}

			//回调
		} else {
			//更新时候判断是否要获取
			//早盘
			if sportGameQuery.GetIsBet() == for_game.GAME_IS_BET_0 &&
				sportsGame.GetIsBet() == for_game.GAME_IS_BET_1 {

				dealGameGuessMorn(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colGuess, updateTime, initCallBackFlag)

				//设置redis早盘信息
				for_game.SetRedisGuessMornDetail(uniqueGameId, appLabelId, gameId, apiOrigin)
			}

			//滚盘
			if sportGameQuery.GetHaveRoll() == for_game.GAME_HAVE_ROLL_0 &&
				sportsGame.GetHaveRoll() == for_game.GAME_HAVE_ROLL_1 {
				//授权
				UseGameGuessRoll(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colUseRoll)

				//取得滚盘
				dealGameGuessRoll(appLabelId, apiOrigin, gameId, sportsGame.GetEventId(), colGuess,
					updateTime, initCallBackFlag)

				//设置redis滚盘信息
				for_game.SetRedisGuessRollDetail(uniqueGameId, appLabelId, gameId, apiOrigin)
			}
		}

		//赔率表中的比赛时间、比赛状态在处理比赛时处理
		//重新设置下盘口信息中的比赛的信息
		if sportGameQuery.GetBeginTime() != sportsGame.GetBeginTime() ||
			sportGameQuery.GetGameStatus() != sportsGame.GetGameStatus() ||
			sportGameQuery.GetGameStatusType() != sportsGame.GetGameStatusType() {

			gameGuessQuery := make([]*share_message.TableESPortsGameGuess, 0)
			//查询所有(可能有早盘和滚盘数据)
			errQuery := colGuess.Find(bson.M{"app_label_id": appLabelId,
				"game_id":    gameId,
				"api_origin": sportsGame.GetApiOrigin()}).All(&gameGuessQuery)

			if errQuery != nil && errQuery != mgo.ErrNotFound {
				logs.Error(errQuery)
				s := fmt.Sprintf("=====比赛详情回调查询竞猜数据失败===查询条件app_label_id:%v,===game_id:%v,===api_origin:%v",
					appLabelId, gameId, sportsGame.GetApiOrigin())
				for_game.WriteFile("ye_zi_api.log", s)
				return uniqueGameId
			}

			if errQuery == mgo.ErrNotFound || nil == gameGuessQuery || len(gameGuessQuery) <= 0 {
				s := fmt.Sprintf("=====比赛详情回调查询竞猜数据为空===查询条件app_label_id:%v,===game_id:%v,===api_origin:%v",
					appLabelId, gameId, sportsGame.GetApiOrigin())
				for_game.WriteFile("ye_zi_api.log", s)
				return uniqueGameId
			}

			if errQuery == nil && nil != gameGuessQuery && len(gameGuessQuery) > 0 {
				for _, guessValue := range gameGuessQuery {

					//重新设置更新时间
					guessValue.UpdateTime = easygo.NewInt64(time.Now().Unix())
					//重新设置比赛信息
					guessValue.BeginTime = easygo.NewString(sportsGame.GetBeginTime())
					guessValue.GameStatus = easygo.NewString(sportsGame.GetGameStatus())
					guessValue.GameStatusType = easygo.NewString(sportsGame.GetGameStatusType())

					//更新操作
					errUpd := colGuess.Update(bson.M{"_id": guessValue.GetId()},
						bson.M{"$set": guessValue})

					if errUpd != nil {
						logs.Error(errUpd)
						s := fmt.Sprintf("=====详情数据回调更新盘口中相关比赛信息数据失败=====盘口中唯一_id:%v",
							guessValue.GetId())
						for_game.WriteFile("ye_zi_api.log", s)
						continue
					}

					if guessValue.GetMornRollGuessFlag() == for_game.GAME_IS_MORN_ROLL_1 {
						//设置redis早盘信息
						for_game.SetRedisGuessMornDetail(uniqueGameId, appLabelId, gameId, sportsGame.GetApiOrigin())
					} else {
						//设置redis滚盘信息
						for_game.SetRedisGuessRollDetail(uniqueGameId, appLabelId, gameId, sportsGame.GetApiOrigin())
					}
				}
			}
		}
	}

	if initCallBackFlag == for_game.INIT_CALLBACK_FLAG_1 {
		s := fmt.Sprintf("=======dealSportGame初始化比赛详情后、处理比赛表   结束==========")
		for_game.WriteFile("ye_zi_api.log", s)
	} else {
		s := fmt.Sprintf("=======dealSportGame回调比赛详情后、处理比赛表   结束==========")
		for_game.WriteFile("ye_zi_api.log", s)
	}

	return uniqueGameId
}

//取得比赛表中唯一自增id
func GetUniqueGameId(appLabelId int32, apiOrigin int32, gameId string) *int64 {

	var uniqueGameId *int64
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
	defer closeFun()

	sportGameQuery := share_message.TableESPortsGame{}
	//通过条件查询
	errQuery := col.Find(bson.M{"app_label_id": appLabelId,
		"game_id":    gameId,
		"api_origin": apiOrigin}).One(&sportGameQuery)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("======比赛表查询失败======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v",
			appLabelId, gameId, apiOrigin)
		for_game.WriteFile("ye_zi_api.log", s)
		return nil
	}

	if errQuery == mgo.ErrNotFound {
		return nil
	}
	if errQuery == nil {
		uniqueGameId = easygo.NewInt64(sportGameQuery.GetId())
	}

	return uniqueGameId
}

//取得推流直播地址
func GetGameLiveUrl(gameId string, eventId string) *share_message.ESPortsGameLivePathObj {

	s := fmt.Sprintf("========GetGameLiveUrl野子科技取得推流直播地址处理   开始=========:参数gameId=%v", gameId)
	for_game.WriteFile("ye_zi_api.log", s)

	var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

	param := url.Values{}
	param.Set("app_secret", yeZiAppSecret)
	param.Set("app_key", yeZiAppKey)
	param.Set("timestamp", timeStr)
	param.Set("resource", "live")
	param.Set("func", "get_video")
	param.Set("game_id", gameId)
	param.Set("event_id", eventId)

	sign := GetYeZiApiSign(param)

	param.Del("app_secret")
	param.Add("sign", sign)

	//发送请求
	data, errGet := Get(yeZiUrl, param)
	if errGet != nil {

		logs.Error(errGet)
		s := fmt.Sprintf("=========野子科技取得推流直播地址请求api失败========,错误信息:%v,请求参数为::::param=%v", errGet, param)
		for_game.WriteFile("ye_zi_api.log", s)
		//easygo.PanicError(errGet)
		return nil

	} else {
		liveUrl := YeZiESPortsGetGameLiveUrl{}
		if "" == data {
			s := fmt.Sprintf("========野子科技取得推流直播地址api数据为空========,请求参数为::::param=%v", param)
			for_game.WriteFile("ye_zi_api.log", s)
			return nil
		} else {
			s := fmt.Sprintf("=========野子科技取得推流直播地址api数据为::::data=%v,=======,请求参数为::::param=%v", data, param)
			for_game.WriteFile("ye_zi_api.log", s)

			err := json.Unmarshal([]byte(data), &liveUrl)

			if err != nil {

				logs.Error(err)
				s := fmt.Sprintf("==========野子科技取得推流直播地址api数据后json转结构体出错::::data=%v=========,请求参数为::::param=%v", data, param)
				for_game.WriteFile("ye_zi_api.log", s)
				easygo.PanicError(err)

				return nil
			}

			if liveUrl.ErrorCode != 0 && liveUrl.ErrorMsg != "" {
				s := fmt.Sprintf("===========野子科技取得推流直播地址请求返回错误信息:错误码ErrorCode:%v,错误信息ErrorMsg:%v,========,请求参数为::::param=%v",
					liveUrl.ErrorCode, liveUrl.ErrorMsg, param)
				for_game.WriteFile("ye_zi_api.log", s)
				return nil
			}

			//判断转化后结构处理业务
			if liveUrl.Code == YEZI_SUCCESS_CODE && liveUrl.Data != nil && liveUrl.Data.LivePaths != nil {
				return liveUrl.Data.LivePaths
			}
		}
	}

	s = fmt.Sprintf("========GetGameLiveUrl野子科技取得推流直播地址数据处理结束   开始==========")
	for_game.WriteFile("ye_zi_api.log", s)

	return &share_message.ESPortsGameLivePathObj{}
}

//更新的时候处理封盘时间、结果产生时间
func dealOddsTimeUpdate(apiGuessInfo []*share_message.ApiGuessObject, queryGuessInfo []*share_message.ApiGuessObject, updateTime int64) {

	queryGuessItemMap := make(map[string]*share_message.ApiItemObject, 0)
	//将数据库的结构体放到map中
	for _, queryValue := range queryGuessInfo {
		queryItems := queryValue.GetItems()
		if queryItems != nil && len(queryItems) > 0 {
			for _, queryItemValue := range queryItems {
				queryGuessItemMap[queryItemValue.GetBetNum()] = queryItemValue
			}
		}
	}

	//将开奖的时候的api推送时间和本地时间设置到更新结构体
	for _, apiValue := range apiGuessInfo {
		apiItemInfo := apiValue.GetItems()
		if nil != apiItemInfo && len(apiItemInfo) > 0 {
			for _, apiItemValue := range apiItemInfo {
				tempMapValue := queryGuessItemMap[apiItemValue.GetBetNum()]
				if nil != tempMapValue {
					//1、封盘时间
					//回调才设置、投注状态从1到0  或  从1到3记录; 若3到1的时候不变 3到0的时候不变
					// (滚盘开奖的时候这个时间的前300秒投注为无效单)
					apiItemValue.StatusTime = easygo.NewInt64(tempMapValue.GetStatusTime())

					if tempMapValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_1 &&
						apiItemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_0 {
						if updateTime != 0 {
							apiItemValue.StatusTime = easygo.NewInt64(updateTime)
						} else {
							apiItemValue.StatusTime = easygo.NewInt64(time.Now().Unix())
						}
					} else if tempMapValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_1 &&
						apiItemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_3 {
						if updateTime != 0 {
							apiItemValue.StatusTime = easygo.NewInt64(updateTime)
						} else {
							apiItemValue.StatusTime = easygo.NewInt64(time.Now().Unix())
						}
					}

					//2、结果产生时间
					apiItemValue.ResultTime = easygo.NewInt64(tempMapValue.GetResultTime())
					if (tempMapValue.GetStatus() == for_game.GAME_GUESS_ITEM_STATUS_0 ||
						tempMapValue.GetWin() == for_game.GAME_GUESS_ITEM_WIN_NORST) &&
						(apiItemValue.GetStatus() == for_game.GAME_GUESS_ITEM_STATUS_1 &&
							apiItemValue.GetWin() != for_game.GAME_GUESS_ITEM_WIN_NORST) {
						if updateTime != 0 {
							apiItemValue.ResultTime = easygo.NewInt64(updateTime)
						} else {
							apiItemValue.ResultTime = easygo.NewInt64(time.Now().Unix())
						}
					}
				} else {
					//代表投注项目是第一次入库
					//1、新增设置封盘时间
					//新增的时候只记录0关闭或者3暂停的时间
					if apiItemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_0 ||
						apiItemValue.GetOddsStatus() == for_game.GAME_GUESS_ODDS_STATUS_3 {
						if updateTime != 0 {
							apiItemValue.StatusTime = easygo.NewInt64(updateTime)
						} else {
							apiItemValue.StatusTime = easygo.NewInt64(time.Now().Unix())
						}
					}

					//2、结果产生时间
					if apiItemValue.GetStatus() == for_game.GAME_GUESS_ITEM_STATUS_1 &&
						apiItemValue.GetWin() != for_game.GAME_GUESS_ITEM_WIN_NORST {
						if updateTime != 0 {
							apiItemValue.ResultTime = easygo.NewInt64(updateTime)
						} else {
							apiItemValue.ResultTime = easygo.NewInt64(time.Now().Unix())
						}
					}
				}
			}
		}
	}
}

//更新的时候处理后台设置的开启与封盘的设置
func dealBetCntEnableUpdate(apiGuessInfo []*share_message.ApiGuessObject, queryGuessInfo []*share_message.ApiGuessObject) {
	//key betId、value:数据库中对应的盘口状态和显示
	queryBetIdMap := make(map[string][]int32, 0)

	//将数据库中的值放到map中
	for _, queryValue := range queryGuessInfo {
		queryBetIdMap[queryValue.GetBetId()] = []int32{queryValue.GetAppGuessFlag(), queryValue.GetAppGuessViewFlag()}
	}

	for _, apiValue := range apiGuessInfo {
		tempApiBetValue := queryBetIdMap[apiValue.GetBetId()]
		if tempApiBetValue != nil && len(tempApiBetValue) == 2 {
			apiValue.AppGuessFlag = easygo.NewInt32(tempApiBetValue[0])
			apiValue.AppGuessViewFlag = easygo.NewInt32(tempApiBetValue[1])
		} else {
			apiValue.AppGuessFlag = easygo.NewInt32(for_game.GAME_APP_GUESS_FLAG_1)
			apiValue.AppGuessViewFlag = easygo.NewInt32(for_game.GAME_APP_GUESS_VIEW_FLAG_1)
		}
	}
}
