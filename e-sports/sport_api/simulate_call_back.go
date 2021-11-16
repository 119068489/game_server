package sport_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//=================该函数的值以及函数不要调用、只能在测试的时候才能用到==========
var LabelToESportEventIdMap = map[int32]string{
	for_game.ESPORTS_LABEL_WZRY: for_game.YEZI_ESPORTS_EVENT_WZRY, //王者荣耀(野子科技eventID)
	//for_game.ESPORTS_LABEL_DOTA2: for_game.YEZI_ESPORTS_EVENT_DOTA2, //dota2(野子科技eventID)
	for_game.ESPORTS_LABEL_LOL: for_game.YEZI_ESPORTS_EVENT_LOL, //英雄联盟lol(野子科技eventID)
	//for_game.ESPORTS_LABEL_CSGO:  for_game.YEZI_ESPORTS_EVENT_CSGO,  //CSGO(野子科技eventID)
}

var mapGames map[int64]*gameRealTimeRecord = make(map[int64]*gameRealTimeRecord)

type testBody struct {
	GameId     int32  `json:"game_id"`
	EventId    int32  `json:"event_id"`
	Type       string `json:"type"`
	Func       string `json:"func"`
	UpdateTime int64  `json:"update_time"`
}

type gameRealTimeRecord struct {
	round           int32
	recordFirstTime int64

	duration int32
}

func GameListCallBack() {
	s := fmt.Sprintf("=============GameListCallBack模拟回调比赛列表等初始化数据======开始")
	for_game.WriteFile("ye_zi_api.log", s)
	//初始化未开始比赛列表:王者荣耀、DOTA2、英雄联盟、CSGO
	DoGetYeZiAllGamesFunc()

	//初始化未开始、进行中的比赛的详情、比赛动态信息表(早盘、滚盘)信息
	//两队历史交锋、两队胜败统计、两队天敌克制统计(统计一次即可)
	//DoGetYeZiAllGameDetailsFunc()

	s1 := fmt.Sprintf("===========GameListCallBack模拟回调比赛列表等初始化数据==========结束===========")
	for_game.WriteFile("ye_zi_api.log", s1)

	easygo.AfterFunc(time.Duration(3600)*time.Second, func() {
		GameListCallBack()
	})
}

func GameDetailCallBack() {
	port := strconv.FormatInt(int64(PServerInfo.GetWebApiPort()), 10)
	yeziURL := "http://127.0.0.1:" + port + "/notice"

	s := fmt.Sprintf("===============GameDetailCallBack模拟回调比赛详情==========开始")
	for_game.WriteFile("ye_zi_api.log", s)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
	defer closeFun()
	//取得未开始、进行中的数据
	sportGameLists := make([]*share_message.TableESPortsGame, 0)

	errQuery := col.Find(bson.M{"game_status": for_game.GAME_STATUS_0,
		"api_origin": for_game.ESPORTS_API_ORIGIN_ID_YEZI}).All(&sportGameLists)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("====模拟回调比赛详情时、比赛列表接口查询失败===查询条件:game_status:%v,===api_origin:%v",
			for_game.GAME_STATUS_0, for_game.ESPORTS_API_ORIGIN_ID_YEZI)
		for_game.WriteFile("ye_zi_api.log", s)

		easygo.PanicError(errQuery)
	}

	if mgo.ErrNotFound == errQuery || nil == sportGameLists || len(sportGameLists) <= 0 {
		s := fmt.Sprintf("=====模拟回调比赛详情时、比赛列表接口没有数据要处理=======查询条件:game_status:%v,===api_origin:%v",
			for_game.GAME_STATUS_0, for_game.ESPORTS_API_ORIGIN_ID_YEZI)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	if nil == errQuery && nil != sportGameLists && len(sportGameLists) > 0 {

		//循环回调详情数据
		for _, gameValue := range sportGameLists {
			//初始化参数
			param := url.Values{}
			var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

			//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
			param.Set("func", "game_info")
			param.Set("game_id", gameValue.GetGameId())
			param.Set("type", "update")
			param.Set("timestamp", timeStr)
			param.Set("event_id", gameValue.GetEventId())

			gameId, errGameId := strconv.ParseInt(gameValue.GetGameId(), 10, 32)
			easygo.PanicError(errGameId)
			eventId, errEventId := strconv.ParseInt(gameValue.GetEventId(), 10, 32)
			easygo.PanicError(errEventId)

			body := testBody{
				GameId:     int32(gameId),
				EventId:    int32(eventId),
				Type:       "update",
				Func:       "game_info",
				UpdateTime: time.Now().Unix(),
			}

			data, jsErr := json.Marshal(body)
			easygo.PanicError(jsErr)
			rst, postErr := doBytesPost(yeziURL, data, param)
			easygo.PanicError(postErr)
			if rst == "success" {
				s := fmt.Sprintf("=====模拟回调比赛详情比赛id：gameId:%v===回调处理success结束",
					gameValue.GetGameId())
				for_game.WriteFile("ye_zi_api.log", s)
			}

			time.Sleep(1 * time.Second)
		}

		s := fmt.Sprintf("===============GameDetailCallBack模拟回调比赛详情===========结束")
		for_game.WriteFile("ye_zi_api.log", s)
	}
	easygo.AfterFunc(time.Duration(300)*time.Second, func() {
		GameDetailCallBack()
	})
}

func MornBetCallBack() {
	port := strconv.FormatInt(int64(PServerInfo.GetWebApiPort()), 10)
	yeziURL := "http://127.0.0.1:" + port + "/notice"

	s := fmt.Sprintf("===============MornBetCallBack模拟回调早盘===========开始")
	for_game.WriteFile("ye_zi_api.log", s)

	colGuess, closeFunGuess := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_GUESS)
	defer closeFunGuess()

	GuessLists := make([]*share_message.TableESPortsGameGuess, 0)

	errQuery := colGuess.Find(bson.M{"game_status": for_game.GAME_STATUS_0,
		"mornRoll_guess_flag": for_game.GAME_IS_MORN_ROLL_1}).All(&GuessLists)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("==========模拟回调早盘时、早盘接口查询失败=====查询条件为game_status:%v,===mornRoll_guess_flag:%v",
			for_game.GAME_STATUS_0, for_game.GAME_IS_MORN_ROLL_1)
		for_game.WriteFile("ye_zi_api.log", s)
		easygo.PanicError(errQuery)
	}

	if mgo.ErrNotFound == errQuery || nil == GuessLists || len(GuessLists) <= 0 {
		s := fmt.Sprintf("======模拟回调早盘时、早盘接口没有数据要处理=========查询条件为game_status:%v,===mornRoll_guess_flag:%v",
			for_game.GAME_STATUS_0, for_game.GAME_IS_MORN_ROLL_1)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	if nil == errQuery && nil != GuessLists && len(GuessLists) > 0 {

		//循环
		for _, gameValue := range GuessLists {
			//初始化参数
			param := url.Values{}
			var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

			//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
			param.Set("func", "bet_info")
			param.Set("game_id", gameValue.GetGameId())
			param.Set("type", "update")
			param.Set("timestamp", timeStr)
			param.Set("event_id", LabelToESportEventIdMap[gameValue.GetAppLabelId()])

			gameId, errGameId := strconv.ParseInt(gameValue.GetGameId(), 10, 32)
			easygo.PanicError(errGameId)
			eventId, errEventId := strconv.ParseInt(LabelToESportEventIdMap[gameValue.GetAppLabelId()], 10, 32)
			easygo.PanicError(errEventId)

			body := testBody{
				GameId:     int32(gameId),
				EventId:    int32(eventId),
				Type:       "update",
				Func:       "bet_info",
				UpdateTime: time.Now().Unix(),
			}

			data, jsErr := json.Marshal(body)
			easygo.PanicError(jsErr)
			rst, postErr := doBytesPost(yeziURL, data, param)
			easygo.PanicError(postErr)
			if rst == "success" {
				s := fmt.Sprintf("========回调早盘时、比赛id：gameId:%v===回调处理success结束",
					gameValue.GetGameId())
				for_game.WriteFile("ye_zi_api.log", s)
			}

			time.Sleep(1 * time.Second)
		}

		s := fmt.Sprintf("============MornBetCallBack模拟回调早盘=============结束")
		for_game.WriteFile("ye_zi_api.log", s)
	}

	easygo.AfterFunc(time.Duration(600)*time.Second, func() {
		MornBetCallBack()
	})
}

func RollBetCallBack() {
	port := strconv.FormatInt(int64(PServerInfo.GetWebApiPort()), 10)
	yeziURL := "http://127.0.0.1:" + port + "/notice"

	s := fmt.Sprintf("==========RollBetCallBack模拟回调滚盘==========开始")
	for_game.WriteFile("ye_zi_api.log", s)

	colGuess, closeFunGuess := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_GUESS)
	defer closeFunGuess()

	GuessLists := make([]*share_message.TableESPortsGameGuess, 0)

	errQuery := colGuess.Find(bson.M{"game_status": for_game.GAME_STATUS_0,
		"mornRoll_guess_flag": for_game.GAME_IS_MORN_ROLL_2}).All(&GuessLists)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("=======模拟回调滚盘时、滚盘接口查询失败======查询条件为game_status:%v,===mornRoll_guess_flag:%v",
			for_game.GAME_STATUS_0, for_game.GAME_IS_MORN_ROLL_2)
		for_game.WriteFile("ye_zi_api.log", s)

		easygo.PanicError(errQuery)
	}

	if mgo.ErrNotFound == errQuery || nil == GuessLists || len(GuessLists) <= 0 {
		s := fmt.Sprintf("========模拟回调滚盘时、滚盘接口查询没有数据要处理==========查询条件为game_status:%v,===mornRoll_guess_flag:%v",
			for_game.GAME_STATUS_0, for_game.GAME_IS_MORN_ROLL_2)
		for_game.WriteFile("ye_zi_api.log", s)
		return
	}

	if nil == errQuery && nil != GuessLists && len(GuessLists) > 0 {

		//循环
		for _, gameValue := range GuessLists {
			//初始化参数
			param := url.Values{}
			var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

			//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
			param.Set("func", "roll_bet_info")
			param.Set("game_id", gameValue.GetGameId())
			param.Set("type", "update")
			param.Set("timestamp", timeStr)
			param.Set("event_id", LabelToESportEventIdMap[gameValue.GetAppLabelId()])

			gameId, errGameId := strconv.ParseInt(gameValue.GetGameId(), 10, 32)
			easygo.PanicError(errGameId)
			eventId, errEventId := strconv.ParseInt(LabelToESportEventIdMap[gameValue.GetAppLabelId()], 10, 32)
			easygo.PanicError(errEventId)

			body := testBody{
				GameId:     int32(gameId),
				EventId:    int32(eventId),
				Type:       "update",
				Func:       "roll_bet_info",
				UpdateTime: time.Now().Unix(),
			}

			data, jsErr := json.Marshal(body)
			easygo.PanicError(jsErr)
			rst, postErr := doBytesPost(yeziURL, data, param)
			easygo.PanicError(postErr)
			if rst == "success" {
				s := fmt.Sprintf("=======模拟回调滚盘时、比赛id：gameId:%v===回调处理success结束",
					gameValue.GetGameId())
				for_game.WriteFile("ye_zi_api.log", s)
			}

			time.Sleep(1 * time.Second)
		}

		s := fmt.Sprintf("==========RollBetCallBack模拟回调滚盘==========结束")
		for_game.WriteFile("ye_zi_api.log", s)
	}

	easygo.AfterFunc(time.Duration(60)*time.Second, func() {
		RollBetCallBack()
	})
}

func RealTimeCallBack() {
	port := strconv.FormatInt(int64(PServerInfo.GetWebApiPort()), 10)
	yeziURL := "http://127.0.0.1:" + port + "/notice"

	//s := fmt.Sprintf("==========RealTimeCallBack模拟回调==========开始")
	//for_game.WriteFile("ye_zi_api_real_time.log", s)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
	defer closeFun()

	gameList := make([]*share_message.TableESPortsGame, 0)

	errQuery := col.Find(bson.M{"game_status": for_game.GAME_STATUS_0,
		"begin_time_int": bson.M{"$lt": time.Now().Unix()},
		"$or": []bson.M{bson.M{"app_label_id": for_game.ESPORTS_LABEL_LOL},
			bson.M{"app_label_id": for_game.ESPORTS_LABEL_WZRY}}}).All(&gameList)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		s := fmt.Sprintf("=======RealTimeCallBack模拟回调接口查询比赛表失败======")
		for_game.WriteFile("ye_zi_api_real_time.log", s)

		easygo.PanicError(errQuery)
	}

	if mgo.ErrNotFound == errQuery || nil == gameList || len(gameList) <= 0 {
		//s := fmt.Sprintf("========RealTimeCallBack模拟回调接口查询没有数据要处理==========")
		//for_game.WriteFile("ye_zi_api_real_time.log", s)
		return
	}

	if nil == errQuery && nil != gameList && len(gameList) > 0 {

		//循环回调
		for _, gameValue := range gameList {
			if gameValue.GetAppLabelId() == for_game.ESPORTS_LABEL_LOL ||
				gameValue.GetAppLabelId() == for_game.ESPORTS_LABEL_WZRY {

				//初始化参数
				param := url.Values{}
				var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

				//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
				param.Set("func", "live_info")
				param.Set("game_id", gameValue.GetGameId())
				param.Set("timestamp", timeStr)
				param.Set("event_id", LabelToESportEventIdMap[gameValue.GetAppLabelId()])

				var data []byte
				var jsErr error

				timeRecord := mapGames[gameValue.GetId()]
				var duration int32
				if nil == timeRecord {

					//首先从redis中取得最新的局数
					rounds := for_game.GetRedisGameRealTimeRounds(gameValue.GetId())

					gameTimeRecord := gameRealTimeRecord{}
					if nil != rounds && rounds.GetGameRounds() > 0 {
						gameTimeRecord.round = rounds.GetGameRounds()

						//取得当前局数
						gameRoundInfo := for_game.GetRedisGameRealTime(gameValue.GetId(), gameTimeRecord.round)

						duration = duration + gameRoundInfo.GetDuration() + 30
					} else {
						gameTimeRecord.round = 1

						duration = 30
					}

					gameTimeRecord.recordFirstTime = time.Now().Unix()
					gameTimeRecord.duration = duration

					mapGames[gameValue.GetId()] = &gameTimeRecord
				} else {

					if gameValue.GetBo() != "" &&
						timeRecord.round > *easygo.NewInt32(gameValue.GetBo()) {
						continue
					}

					timeRecord2 := mapGames[gameValue.GetId()]

					timeRecord2.duration = timeRecord2.duration + 30

					nowTime := time.Now().Unix()

					//40分钟换一局
					if nowTime-timeRecord.recordFirstTime > 2400 {
						gameTimeRecord := gameRealTimeRecord{}
						gameTimeRecord.round = timeRecord.round + 1
						gameTimeRecord.recordFirstTime = time.Now().Unix()
						gameTimeRecord.duration = 30
						mapGames[gameValue.GetId()] = &gameTimeRecord
					}
				}

				timeRecord1 := mapGames[gameValue.GetId()]

				//s := fmt.Sprintf("========RealTimeCallBack处理uniqGameId:%v,round:%v==========", gameValue.GetId(), timeRecord1.round)
				//for_game.WriteFile("ye_zi_api_real_time.log", s)

				if gameValue.GetAppLabelId() == for_game.ESPORTS_LABEL_LOL {

					body := share_message.TableESPortsLOLRealTimeData{
						AppLabelId:    easygo.NewInt32(for_game.ESPORTS_LABEL_LOL),
						AppLabelName:  easygo.NewString(for_game.ESPORTS_LABEL_LOL_NAME),
						ApiOrigin:     easygo.NewInt32(for_game.ESPORTS_API_ORIGIN_ID_YEZI),
						ApiOriginName: easygo.NewString(for_game.ApiOriginIdToNameMap[for_game.ESPORTS_API_ORIGIN_ID_YEZI]),
						GameId:        easygo.NewInt32(gameValue.GetGameId()),
						GameRound:     easygo.NewInt32(timeRecord1.round),
						GameStatus:    easygo.NewInt32(1),
						EventId:       easygo.NewString(gameValue.GetEventId()),
						Duration:      easygo.NewInt32(timeRecord1.duration),
						//TeamA:            nil,
						//TeamB:            nil,
						//PlayerAInfo:      nil,
						//PlayerBInfo:      nil,
						CreateTime: easygo.NewInt64(time.Now().Unix()),
						UpdateTime: easygo.NewInt64(time.Now().Unix()),
					}

					TeamA := share_message.ApiLOLTeam{
						Faction:      easygo.NewString("blue"),
						Picks:        []int32{1, 2, 3, 4, 5},
						Bans:         []int32{1, 2, 3, 4, 5},
						Name:         easygo.NewString(gameValue.GetTeamA().GetName()),
						Id:           easygo.NewInt32(gameValue.GetTeamA().GetTeamId()),
						Score:        easygo.NewInt32(GetRandomTest()),
						Glod:         easygo.NewInt32(GetRandomTeamGoldTest()),
						Subsidy:      easygo.NewInt32(GetRandomTest()),
						TowerState:   easygo.NewInt32(GetRandomTest()),
						Drakes:       easygo.NewInt32(GetRandomTest()),
						NahsorBarons: easygo.NewInt32(GetRandomTest()),
					}

					TeamB := share_message.ApiLOLTeam{
						Faction:      easygo.NewString("red"),
						Picks:        []int32{1, 2, 3, 4, 5},
						Bans:         []int32{1, 2, 3, 4, 5},
						Name:         easygo.NewString(gameValue.GetTeamA().GetName()),
						Id:           easygo.NewInt32(gameValue.GetTeamA().GetTeamId()),
						Score:        easygo.NewInt32(GetRandomTest()),
						Glod:         easygo.NewInt32(GetRandomTeamGoldTest()),
						Subsidy:      easygo.NewInt32(GetRandomTest()),
						TowerState:   easygo.NewInt32(GetRandomTest()),
						Drakes:       easygo.NewInt32(GetRandomTest()),
						NahsorBarons: easygo.NewInt32(GetRandomTest()),
					}

					body.TeamA = &TeamA
					body.TeamB = &TeamB

					playerAInfos := make([]*share_message.ApiLOLPlayer, 0)
					for i := 0; i < 5; i++ {
						tmpPlay := share_message.ApiLOLPlayer{
							Name:    easygo.NewString("A队员"),
							HeroId:  easygo.NewInt32(1),
							Kills:   easygo.NewInt32(GetRandomTest()),
							Death:   easygo.NewInt32(GetRandomTest()),
							Assists: easygo.NewInt32(GetRandomTest()),
							Subsidy: easygo.NewInt32(GetRandomTest()),
							Gold:    easygo.NewInt32(GetRandomPlayGoldTest()),
							Item:    []int32{1, 2, 3, 4, 5, 6, 7},
						}
						playerAInfos = append(playerAInfos, &tmpPlay)
					}
					body.PlayerAInfo = playerAInfos

					playerBInfos := make([]*share_message.ApiLOLPlayer, 0)
					for i := 0; i < 5; i++ {
						tmpPlay := share_message.ApiLOLPlayer{
							Name:    easygo.NewString("B队员"),
							HeroId:  easygo.NewInt32(1),
							Kills:   easygo.NewInt32(GetRandomTest()),
							Death:   easygo.NewInt32(GetRandomTest()),
							Assists: easygo.NewInt32(GetRandomTest()),
							Subsidy: easygo.NewInt32(GetRandomTest()),
							Gold:    easygo.NewInt32(GetRandomPlayGoldTest()),
							Item:    []int32{1, 2, 3, 4, 5, 6, 7},
						}
						playerBInfos = append(playerBInfos, &tmpPlay)
					}
					body.PlayerBInfo = playerBInfos

					data, jsErr = json.Marshal(body)
					easygo.PanicError(jsErr)
					rst, postErr := doBytesPost(yeziURL, data, param)
					easygo.PanicError(postErr)
					if rst == "success" {
						//s := fmt.Sprintf("=======RealTimeCallBack模拟回调LOL比赛id：gameId:%v===回调处理success结束",
						//	gameValue.GetGameId())
						//for_game.WriteFile("ye_zi_api_real_time.log", s)
					}

				} else if gameValue.GetAppLabelId() == for_game.ESPORTS_LABEL_WZRY {
					body := share_message.TableESPortsWZRYRealTimeData{
						AppLabelId:    easygo.NewInt32(for_game.ESPORTS_LABEL_WZRY),
						AppLabelName:  easygo.NewString(for_game.ESPORTS_LABEL_WZRY_NAME),
						ApiOrigin:     easygo.NewInt32(for_game.ESPORTS_API_ORIGIN_ID_YEZI),
						ApiOriginName: easygo.NewString(for_game.ApiOriginIdToNameMap[for_game.ESPORTS_API_ORIGIN_ID_YEZI]),
						GameId:        easygo.NewInt32(gameValue.GetGameId()),
						GameRound:     easygo.NewInt32(timeRecord1.round),
						GameStatus:    easygo.NewInt32(1),
						EventId:       easygo.NewString(gameValue.GetEventId()),
						Duration:      easygo.NewInt32(timeRecord1.duration),
						//TeamA:            nil,
						//TeamB:            nil,
						//PlayerAInfo:      nil,
						//PlayerBInfo:      nil,
						CreateTime: easygo.NewInt64(time.Now().Unix()),
						UpdateTime: easygo.NewInt64(time.Now().Unix()),
					}

					TeamA := share_message.ApiWZRYTeam{
						Faction:    easygo.NewString("blue"),
						Picks:      []int32{1, 2, 3, 4, 5},
						Bans:       []int32{1, 2, 3, 4, 5},
						Name:       easygo.NewString(gameValue.GetTeamA().GetName()),
						Id:         easygo.NewInt32(gameValue.GetTeamA().GetTeamId()),
						Score:      easygo.NewInt32(GetRandomTest()),
						TowerState: easygo.NewInt32(GetRandomTest()),
					}

					TeamB := share_message.ApiWZRYTeam{
						Faction:    easygo.NewString("red"),
						Picks:      []int32{1, 2, 3, 4, 5},
						Bans:       []int32{1, 2, 3, 4, 5},
						Name:       easygo.NewString(gameValue.GetTeamA().GetName()),
						Id:         easygo.NewInt32(gameValue.GetTeamA().GetTeamId()),
						Score:      easygo.NewInt32(GetRandomTest()),
						TowerState: easygo.NewInt32(GetRandomTest()),
					}

					body.TeamA = &TeamA
					body.TeamB = &TeamB

					playerAInfos := make([]*share_message.ApiWZRYPlayer, 0)
					for i := 0; i < 5; i++ {
						tmpPlay := share_message.ApiWZRYPlayer{
							Name:    easygo.NewString("A队员"),
							HeroId:  easygo.NewInt32(1),
							Kills:   easygo.NewInt32(GetRandomTest()),
							Death:   easygo.NewInt32(GetRandomTest()),
							Assists: easygo.NewInt32(GetRandomTest()),
							Gold:    easygo.NewInt32(GetRandomPlayGoldTest()),
							Item:    []int32{1, 2, 3, 4, 5, 6, 7},
						}
						playerAInfos = append(playerAInfos, &tmpPlay)
					}
					body.PlayerAInfo = playerAInfos

					playerBInfos := make([]*share_message.ApiWZRYPlayer, 0)
					for i := 0; i < 5; i++ {
						tmpPlay := share_message.ApiWZRYPlayer{
							Name:    easygo.NewString("B队员"),
							HeroId:  easygo.NewInt32(1),
							Kills:   easygo.NewInt32(GetRandomTest()),
							Death:   easygo.NewInt32(GetRandomTest()),
							Assists: easygo.NewInt32(GetRandomTest()),
							Gold:    easygo.NewInt32(GetRandomPlayGoldTest()),
							Item:    []int32{1, 2, 3, 4, 5, 6, 7},
						}
						playerBInfos = append(playerBInfos, &tmpPlay)
					}
					body.PlayerBInfo = playerBInfos

					data, jsErr = json.Marshal(body)
					easygo.PanicError(jsErr)
					rst, postErr := doBytesPost(yeziURL, data, param)
					easygo.PanicError(postErr)
					if rst == "success" {
						//s := fmt.Sprintf("=======RealTimeCallBack模拟回调WZRY比赛id：gameId:%v===回调处理success结束",
						//	gameValue.GetGameId())
						//for_game.WriteFile("ye_zi_api_real_time.log", s)
					}
				}
			}

			//time.Sleep(500 * time.Millisecond)
		}

		//s := fmt.Sprintf("==========RealTimeCallBack模拟回调==========结束")
		//for_game.WriteFile("ye_zi_api_real_time.log", s)
	}

	easygo.AfterFunc(time.Duration(20)*time.Second, func() {
		RealTimeCallBack()
	})
}

//body提交二进制数据
func doBytesPost(url string, data []byte, param url.Values) (string, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	easygo.PanicError(err)
	request.Header.Set("Connection", "Keep-Alive")

	reqParam := param.Encode()
	request.Header.Set("xxe-request", reqParam)

	param.Set("app_secret", yeZiAppSecret)
	reqParam1 := param.Encode()
	signParam := for_game.Md5(reqParam1)
	request.Header.Set("xxe-sign", signParam)
	var resp *http.Response

	resp, err = http.DefaultClient.Do(request)
	easygo.PanicError(err)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err)
	return string(b), err
}

//测试环境中随机取得数
func GetRandomTest() int32 {
	var temps []int32 = []int32{3, 5, 6, 8, 10, 12, 17, 14, 18}
	tempKey := temps[rand.Intn(len(temps))]

	return tempKey
}

func GetRandomTeamGoldTest() int32 {
	var temps []int32 = []int32{2000, 3000, 4000, 5000, 5500, 6500, 7000, 8000, 9000, 12000, 13000}
	tempKey := temps[rand.Intn(len(temps))]

	return tempKey
}

func GetRandomPlayGoldTest() int32 {
	var temps []int32 = []int32{200, 300, 400, 500, 550, 650, 700, 800, 900, 1200, 1300}
	tempKey := temps[rand.Intn(len(temps))]

	return tempKey
}
