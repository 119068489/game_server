package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"strings"
)

const (
	ESPORT_REDIS_GAME_DETAIL_HEAD_KEY = "redis_esport:game_detail_head" //比赛详情头部信息的key(比赛id)(3小时)
	ESPORT_REDIS_GAME_GUESS_MORN_KEY  = "redis_esport:game_guess_morn"  //比赛早盘的key(比赛id)(3小时)
	ESPORT_REDIS_GAME_GUESS_ROLL_KEY  = "redis_esport:game_guess_roll"  //比赛滚盘的key(比赛id)(3小时)

	ESPORT_REDIS_GAME_DETAIL_HEAD_GROUP_KEY = "redis_esport:game_detail_head_group" //比赛详情头部信息的组合用key(apiOrigin_appLabelId_gameId)(3小时)

)

//设置redis详情头部的信息
func SetRedisGameDetailHead(uniqueGameId int64) {

	gameDetailHead := GetDBGameDetailHead(uniqueGameId)
	if nil != gameDetailHead {
		//json序列化
		data, errJs := json.Marshal(gameDetailHead)
		easygo.PanicError(errJs)
		//redis 3小时
		err := easygo.RedisMgr.GetC().SetWithTime(
			MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_KEY, uniqueGameId),
			string(data),
			ESPORT_GAME_REDIS_EXPIRE_TIME)
		easygo.PanicError(err)

		//设置组合key
		groupKey := GetGameHeadGroupKey(gameDetailHead.GetApiOrigin(), gameDetailHead.GetAppLabelId(), gameDetailHead.GetGameId())
		err1 := easygo.RedisMgr.GetC().SetWithTime(
			MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_GROUP_KEY, groupKey),
			string(data),
			ESPORT_GAME_REDIS_EXPIRE_TIME)
		easygo.PanicError(err1)
	} else {
		keys := []interface{}{MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_KEY, uniqueGameId)}

		easygo.RedisMgr.GetC().Delete(keys...)
	}
}

//取得redis详情头部的信息
func GetRedisGameDetailHead(uniqueGameId int64) *client_hall.ESportGameObject {

	b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_KEY, uniqueGameId))
	easygo.PanicError(err)

	var gameDetailHead *client_hall.ESportGameObject
	if !b {
		gameDetailHead = GetDBGameDetailHead(uniqueGameId)
		if nil != gameDetailHead {
			//json序列化
			data, errJs := json.Marshal(gameDetailHead)
			easygo.PanicError(errJs)
			//redis 3小时
			err := easygo.RedisMgr.GetC().SetWithTime(
				MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_KEY, uniqueGameId),
				string(data),
				ESPORT_GAME_REDIS_EXPIRE_TIME)
			easygo.PanicError(err)

			//设置组合key
			groupKey := GetGameHeadGroupKey(gameDetailHead.GetApiOrigin(), gameDetailHead.GetAppLabelId(), gameDetailHead.GetGameId())
			err1 := easygo.RedisMgr.GetC().SetWithTime(
				MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_GROUP_KEY, groupKey),
				string(data),
				ESPORT_GAME_REDIS_EXPIRE_TIME)
			easygo.PanicError(err1)
		}
	} else {
		obj := &client_hall.ESportGameObject{}
		value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_KEY, uniqueGameId))
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), obj)
		easygo.PanicError(errJs)

		gameDetailHead = obj
	}
	return gameDetailHead
}

//数据库通过uniqueGameId取得详情头部的信息
//uniqueGameId:数据库唯一id
func GetDBGameDetailHead(uniqueGameId int64) *client_hall.ESportGameObject {

	//从数据库中取得数据
	dbESPortsGame := share_message.TableESPortsGame{}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_GAME)
	defer closeFun()

	var gameObject *client_hall.ESportGameObject
	err := col.Find(bson.M{"_id": uniqueGameId}).One(&dbESPortsGame)

	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return gameObject
	}

	if err == mgo.ErrNotFound {
		return gameObject
	}

	if err == nil {
		gameObject = &client_hall.ESportGameObject{
			UniqueGameId: easygo.NewInt64(dbESPortsGame.GetId()),
			GameName:     easygo.NewString(dbESPortsGame.GetMatchName() + dbESPortsGame.GetMatchStage() + "-BO" + dbESPortsGame.GetBo()),
			TeamAInfo: &client_hall.TeamObject{
				TeamId: easygo.NewString(dbESPortsGame.GetTeamA().GetTeamId()),
				Name:   easygo.NewString(dbESPortsGame.GetTeamA().GetName()),
				Icon:   easygo.NewString(dbESPortsGame.GetTeamA().GetIcon()),
			},
			ScoreA: easygo.NewString(dbESPortsGame.GetScoreA()),
			ScoreB: easygo.NewString(dbESPortsGame.GetScoreB()),
			TeamBInfo: &client_hall.TeamObject{
				TeamId: easygo.NewString(dbESPortsGame.GetTeamB().GetTeamId()),
				Name:   easygo.NewString(dbESPortsGame.GetTeamB().GetName()),
				Icon:   easygo.NewString(dbESPortsGame.GetTeamB().GetIcon()),
			},
			BeginTime:    easygo.NewInt64(dbESPortsGame.GetBeginTimeInt()),
			BeginTimeStr: easygo.NewString(dbESPortsGame.GetBeginTime()),
			GameStatus:   easygo.NewString(dbESPortsGame.GetGameStatus()),
			HaveRoll:     easygo.NewInt32(dbESPortsGame.GetHaveRoll()),
			AppLabelId:   easygo.NewInt32(dbESPortsGame.GetAppLabelId()),
			ApiOrigin:    easygo.NewInt32(dbESPortsGame.GetApiOrigin()),
			GameId:       easygo.NewString(dbESPortsGame.GetGameId()),
			IsLottery:    easygo.NewInt32(dbESPortsGame.GetIsLottery()),
			HistoryId:    easygo.NewInt64(dbESPortsGame.GetHistoryId()),
		}

		//比赛图标需要重配置表中取得后然后重新匹配
		gameLabelRedisObj := GetRedisGameLabel()
		if gameLabelRedisObj != nil {
			if gameObject.GetAppLabelId() == ESPORTS_LABEL_WZRY {
				gameObject.GameIcon = easygo.NewString(gameLabelRedisObj.GetWZRYIcon())
			} else if gameObject.GetAppLabelId() == ESPORTS_LABEL_DOTA2 {
				gameObject.GameIcon = easygo.NewString(gameLabelRedisObj.GetDOTAIcon())
			} else if gameObject.GetAppLabelId() == ESPORTS_LABEL_LOL {
				gameObject.GameIcon = easygo.NewString(gameLabelRedisObj.GetLOLIcon())
			} else if gameObject.GetAppLabelId() == ESPORTS_LABEL_CSGO {
				gameObject.GameIcon = easygo.NewString(gameLabelRedisObj.GetCSGOIcon())
			} else if gameObject.GetAppLabelId() == ESPORTS_LABEL_OTHER {
				gameObject.GameIcon = easygo.NewString(gameLabelRedisObj.GetOTHERIcon())
			}
		}
	}

	return gameObject
}

//取得redis详情头部的信息(key(apiOrigin_appLabelId_gameId))
func GetRedisGameDetailHeadGroup(apiOrigin int32, appLabelId int32, gameId string) *client_hall.ESportGameObject {

	groupKey := GetGameHeadGroupKey(apiOrigin, appLabelId, gameId)

	b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_GROUP_KEY, groupKey))
	easygo.PanicError(err)

	var gameDetailHead *client_hall.ESportGameObject
	if !b {
		gameDetailHead = GetDBGameDetailHeadGroup(apiOrigin, appLabelId, gameId)
		if nil != gameDetailHead {
			//json序列化
			data, errJs := json.Marshal(gameDetailHead)
			easygo.PanicError(errJs)
			//redis 3小时
			err := easygo.RedisMgr.GetC().SetWithTime(
				MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_GROUP_KEY, groupKey),
				string(data),
				ESPORT_GAME_REDIS_EXPIRE_TIME)
			easygo.PanicError(err)

			//设置唯一id key
			err1 := easygo.RedisMgr.GetC().SetWithTime(
				MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_KEY, gameDetailHead.GetUniqueGameId()),
				string(data),
				ESPORT_GAME_REDIS_EXPIRE_TIME)
			easygo.PanicError(err1)
		}
	} else {
		obj := &client_hall.ESportGameObject{}
		value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_GROUP_KEY, groupKey))
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), obj)
		easygo.PanicError(errJs)

		gameDetailHead = obj
	}
	return gameDetailHead
}

//数据库通过key(apiOrigin_appLabelId_gameId)取得详情头部的信息
func GetDBGameDetailHeadGroup(apiOrigin int32, appLabelId int32, gameId string) *client_hall.ESportGameObject {

	//从数据库中取得数据
	dbESPortsGame := share_message.TableESPortsGame{}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_GAME)
	defer closeFun()

	var gameObject *client_hall.ESportGameObject
	err := col.Find(bson.M{"app_label_id": appLabelId,
		"game_id":    gameId,
		"api_origin": apiOrigin}).One(&dbESPortsGame)

	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return gameObject
	}

	if err == mgo.ErrNotFound {
		return gameObject
	}

	if err == nil {
		gameObject = &client_hall.ESportGameObject{
			UniqueGameId: easygo.NewInt64(dbESPortsGame.GetId()),
			GameName:     easygo.NewString(dbESPortsGame.GetMatchName() + dbESPortsGame.GetMatchStage() + "-BO" + dbESPortsGame.GetBo()),
			TeamAInfo: &client_hall.TeamObject{
				TeamId: easygo.NewString(dbESPortsGame.GetTeamA().GetTeamId()),
				Name:   easygo.NewString(dbESPortsGame.GetTeamA().GetName()),
				Icon:   easygo.NewString(dbESPortsGame.GetTeamA().GetIcon()),
			},
			ScoreA: easygo.NewString(dbESPortsGame.GetScoreA()),
			ScoreB: easygo.NewString(dbESPortsGame.GetScoreB()),
			TeamBInfo: &client_hall.TeamObject{
				TeamId: easygo.NewString(dbESPortsGame.GetTeamB().GetTeamId()),
				Name:   easygo.NewString(dbESPortsGame.GetTeamB().GetName()),
				Icon:   easygo.NewString(dbESPortsGame.GetTeamB().GetIcon()),
			},
			BeginTime:    easygo.NewInt64(dbESPortsGame.GetBeginTimeInt()),
			BeginTimeStr: easygo.NewString(dbESPortsGame.GetBeginTime()),
			GameStatus:   easygo.NewString(dbESPortsGame.GetGameStatus()),
			HaveRoll:     easygo.NewInt32(dbESPortsGame.GetHaveRoll()),
			AppLabelId:   easygo.NewInt32(dbESPortsGame.GetAppLabelId()),
			ApiOrigin:    easygo.NewInt32(dbESPortsGame.GetApiOrigin()),
			GameId:       easygo.NewString(dbESPortsGame.GetGameId()),
			IsLottery:    easygo.NewInt32(dbESPortsGame.GetIsLottery()),
			HistoryId:    easygo.NewInt64(dbESPortsGame.GetHistoryId()),
		}

		//比赛图标需要重配置表中取得后然后重新匹配
		gameLabelRedisObj := GetRedisGameLabel()
		if gameLabelRedisObj != nil {
			if gameObject.GetAppLabelId() == ESPORTS_LABEL_WZRY {
				gameObject.GameIcon = easygo.NewString(gameLabelRedisObj.GetWZRYIcon())
			} else if gameObject.GetAppLabelId() == ESPORTS_LABEL_DOTA2 {
				gameObject.GameIcon = easygo.NewString(gameLabelRedisObj.GetDOTAIcon())
			} else if gameObject.GetAppLabelId() == ESPORTS_LABEL_LOL {
				gameObject.GameIcon = easygo.NewString(gameLabelRedisObj.GetLOLIcon())
			} else if gameObject.GetAppLabelId() == ESPORTS_LABEL_CSGO {
				gameObject.GameIcon = easygo.NewString(gameLabelRedisObj.GetCSGOIcon())
			} else if gameObject.GetAppLabelId() == ESPORTS_LABEL_OTHER {
				gameObject.GameIcon = easygo.NewString(gameLabelRedisObj.GetOTHERIcon())
			}
		}
	}

	return gameObject
}

//取得redis取得竞猜详情
func GetRedisGuessDetail(uniqueGameId int64, appLabelId int32, gameId string, apiOrigin int32, mornRollGuessFlag int32) *share_message.GameGuessDetailObject {

	var guessDetail *share_message.GameGuessDetailObject
	if mornRollGuessFlag == GAME_IS_MORN_ROLL_1 {
		guessDetail = GetRedisGuessMornDetail(uniqueGameId, appLabelId, gameId, apiOrigin)
	} else if mornRollGuessFlag == GAME_IS_MORN_ROLL_2 {
		guessDetail = GetRedisGuessRollDetail(uniqueGameId, appLabelId, gameId, apiOrigin)
	}

	return guessDetail
}

//设置redis竞猜早盘详情
func SetRedisGuessMornDetail(uniqueGameId int64, appLabelId int32, gameId string, apiOrigin int32) {

	guessDetail := GetDBGuessDetail(appLabelId, gameId, apiOrigin, GAME_IS_MORN_ROLL_1)

	if nil != guessDetail {
		//过滤掉不显示的
		GetFilterGameDetailBet(guessDetail)

		//json序列化
		data, errJs := json.Marshal(guessDetail)
		easygo.PanicError(errJs)
		err := easygo.RedisMgr.GetC().SetWithTime(
			MakeRedisKey(ESPORT_REDIS_GAME_GUESS_MORN_KEY, uniqueGameId),
			string(data),
			ESPORT_GAME_REDIS_EXPIRE_TIME)
		easygo.PanicError(err)
	} else {
		keys := []interface{}{MakeRedisKey(ESPORT_REDIS_GAME_GUESS_MORN_KEY, uniqueGameId)}

		easygo.RedisMgr.GetC().Delete(keys...)
	}
}

//redis取得竞猜早盘详情
func GetRedisGuessMornDetail(uniqueGameId int64, appLabelId int32, gameId string, apiOrigin int32) *share_message.GameGuessDetailObject {

	var guessDetail *share_message.GameGuessDetailObject
	b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_GUESS_MORN_KEY, uniqueGameId))
	easygo.PanicError(err)

	if !b {
		guessDetail = GetDBGuessDetail(appLabelId, gameId, apiOrigin, GAME_IS_MORN_ROLL_1)
		if nil != guessDetail {

			//过滤掉不显示的
			GetFilterGameDetailBet(guessDetail)
			//json序列化
			data, errJs := json.Marshal(guessDetail)
			easygo.PanicError(errJs)
			//redis 3小时
			err := easygo.RedisMgr.GetC().SetWithTime(
				MakeRedisKey(ESPORT_REDIS_GAME_GUESS_MORN_KEY, uniqueGameId),
				string(data),
				ESPORT_GAME_REDIS_EXPIRE_TIME)
			easygo.PanicError(err)
		}
	} else {
		obj := &share_message.GameGuessDetailObject{}
		value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_GUESS_MORN_KEY, uniqueGameId))
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), obj)
		easygo.PanicError(errJs)
		guessDetail = obj
	}
	return guessDetail
}

//设置redis竞猜滚盘详情
func SetRedisGuessRollDetail(uniqueGameId int64, appLabelId int32, gameId string, apiOrigin int32) {

	guessDetail := GetDBGuessDetail(appLabelId, gameId, apiOrigin, GAME_IS_MORN_ROLL_2)

	if nil != guessDetail {
		//过滤掉不显示的
		GetFilterGameDetailBet(guessDetail)

		//json序列化
		data, errJs := json.Marshal(guessDetail)
		easygo.PanicError(errJs)
		err := easygo.RedisMgr.GetC().SetWithTime(
			MakeRedisKey(ESPORT_REDIS_GAME_GUESS_ROLL_KEY, uniqueGameId),
			string(data),
			ESPORT_GAME_REDIS_EXPIRE_TIME)
		easygo.PanicError(err)
	} else {
		keys := []interface{}{MakeRedisKey(ESPORT_REDIS_GAME_GUESS_ROLL_KEY, uniqueGameId)}

		easygo.RedisMgr.GetC().Delete(keys...)
	}

}

//redis取得竞猜滚盘详情
func GetRedisGuessRollDetail(uniqueGameId int64, appLabelId int32, gameId string, apiOrigin int32) *share_message.GameGuessDetailObject {

	var guessDetail *share_message.GameGuessDetailObject
	b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_GUESS_ROLL_KEY, uniqueGameId))
	easygo.PanicError(err)

	if !b {
		guessDetail = GetDBGuessDetail(appLabelId, gameId, apiOrigin, GAME_IS_MORN_ROLL_2)

		if nil != guessDetail {

			//过滤掉不显示的
			GetFilterGameDetailBet(guessDetail)

			//json序列化
			data, errJs := json.Marshal(guessDetail)
			easygo.PanicError(errJs)
			//redis 3小时
			err := easygo.RedisMgr.GetC().SetWithTime(
				MakeRedisKey(ESPORT_REDIS_GAME_GUESS_ROLL_KEY, uniqueGameId),
				string(data),
				ESPORT_GAME_REDIS_EXPIRE_TIME)
			easygo.PanicError(err)
		}
	} else {
		obj := &share_message.GameGuessDetailObject{}
		value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_GUESS_ROLL_KEY, uniqueGameId))
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), obj)
		easygo.PanicError(errJs)
		guessDetail = obj
	}

	return guessDetail
}

//DB取得取得竞猜早盘、滚盘详情
//int64 盘口表中自增唯一id
//[]*share_message.GameGuessOddsNumObject 拼接好的页面显示的数据结构
func GetDBGuessDetail(appLabelId int32, gameId string, apiOrigin int32, mornRollGuessFlag int32) *share_message.GameGuessDetailObject {

	gameGuessDetailObject := &share_message.GameGuessDetailObject{}
	var obj []*share_message.GameGuessOddsNumObject

	//从数据库中取得数据
	dbGuessQuery := share_message.TableESPortsGameGuess{}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_GAME_GUESS)
	defer closeFun()

	err := col.Find(bson.M{"app_label_id": appLabelId,
		"game_id":             gameId,
		"api_origin":          apiOrigin,
		"mornRoll_guess_flag": mornRollGuessFlag}).One(&dbGuessQuery)

	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return nil
	}

	if err == mgo.ErrNotFound {
		return nil
	}

	if err == nil {
		//设置比赛状态、比赛时间
		gameGuessDetailObject.GameStatus = easygo.NewString(dbGuessQuery.GetGameStatus())
		gameGuessDetailObject.BeginTime = easygo.NewString(dbGuessQuery.GetBeginTime())
		gameGuessDetailObject.UniqueGameGuessId = easygo.NewInt64(dbGuessQuery.GetId())
		//组装页面要的数据
		guess := dbGuessQuery.GetGuess()
		if guess != nil && len(guess) > 0 {
			obj = make([]*share_message.GameGuessOddsNumObject, 0)

			gameGuessOddsNumObject0 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject1 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject2 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject3 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject4 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject5 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject6 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject7 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject8 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject9 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject10 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject11 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject12 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject13 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject14 := &share_message.GameGuessOddsNumObject{}
			gameGuessOddsNumObject15 := &share_message.GameGuessOddsNumObject{}

			for _, guessValue := range guess {
				switch guessValue.GetNum() {
				case "0":
					//设置局数的属性
					if gameGuessOddsNumObject0.GetNum() == "" {
						gameGuessOddsNumObject0.Num = easygo.NewString("0")
					}
					if gameGuessOddsNumObject0.GetNumName() == "" {
						gameGuessOddsNumObject0.NumName = easygo.NewString(GuessNumMap["0"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject0)
				case "1":
					if gameGuessOddsNumObject1.GetNum() == "" {
						gameGuessOddsNumObject1.Num = easygo.NewString("1")
					}
					if gameGuessOddsNumObject1.GetNumName() == "" {
						gameGuessOddsNumObject1.NumName = easygo.NewString(GuessNumMap["1"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject1)
				case "2":
					if gameGuessOddsNumObject2.GetNum() == "" {
						gameGuessOddsNumObject2.Num = easygo.NewString("2")
					}
					if gameGuessOddsNumObject2.GetNumName() == "" {
						gameGuessOddsNumObject2.NumName = easygo.NewString(GuessNumMap["2"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject2)
				case "3":
					if gameGuessOddsNumObject3.GetNum() == "" {
						gameGuessOddsNumObject3.Num = easygo.NewString("3")
					}
					if gameGuessOddsNumObject3.GetNumName() == "" {
						gameGuessOddsNumObject3.NumName = easygo.NewString(GuessNumMap["3"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject3)
				case "4":
					if gameGuessOddsNumObject4.GetNum() == "" {
						gameGuessOddsNumObject4.Num = easygo.NewString("4")
					}
					if gameGuessOddsNumObject4.GetNumName() == "" {
						gameGuessOddsNumObject4.NumName = easygo.NewString(GuessNumMap["4"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject4)
				case "5":
					if gameGuessOddsNumObject5.GetNum() == "" {
						gameGuessOddsNumObject5.Num = easygo.NewString("5")
					}
					if gameGuessOddsNumObject5.GetNumName() == "" {
						gameGuessOddsNumObject5.NumName = easygo.NewString(GuessNumMap["5"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject5)
				case "6":
					if gameGuessOddsNumObject6.GetNum() == "" {
						gameGuessOddsNumObject6.Num = easygo.NewString("6")
					}
					if gameGuessOddsNumObject6.GetNumName() == "" {
						gameGuessOddsNumObject6.NumName = easygo.NewString(GuessNumMap["6"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject6)
				case "7":
					if gameGuessOddsNumObject7.GetNum() == "" {
						gameGuessOddsNumObject7.Num = easygo.NewString("7")
					}
					if gameGuessOddsNumObject7.GetNumName() == "" {
						gameGuessOddsNumObject7.NumName = easygo.NewString(GuessNumMap["7"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject7)
				case "8":
					if gameGuessOddsNumObject8.GetNum() == "" {
						gameGuessOddsNumObject8.Num = easygo.NewString("8")
					}
					if gameGuessOddsNumObject8.GetNumName() == "" {
						gameGuessOddsNumObject8.NumName = easygo.NewString(GuessNumMap["8"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject8)
				case "9":
					if gameGuessOddsNumObject9.GetNum() == "" {
						gameGuessOddsNumObject9.Num = easygo.NewString("9")
					}
					if gameGuessOddsNumObject9.GetNumName() == "" {
						gameGuessOddsNumObject9.NumName = easygo.NewString(GuessNumMap["9"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject9)
				case "10":
					if gameGuessOddsNumObject10.GetNum() == "" {
						gameGuessOddsNumObject10.Num = easygo.NewString("10")
					}
					if gameGuessOddsNumObject10.GetNumName() == "" {
						gameGuessOddsNumObject10.NumName = easygo.NewString(GuessNumMap["10"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject10)
				case "11":
					if gameGuessOddsNumObject11.GetNum() == "" {
						gameGuessOddsNumObject11.Num = easygo.NewString("11")
					}
					if gameGuessOddsNumObject11.GetNumName() == "" {
						gameGuessOddsNumObject11.NumName = easygo.NewString(GuessNumMap["11"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject11)
				case "12":
					if gameGuessOddsNumObject12.GetNum() == "" {
						gameGuessOddsNumObject12.Num = easygo.NewString("12")
					}
					if gameGuessOddsNumObject12.GetNumName() == "" {
						gameGuessOddsNumObject12.NumName = easygo.NewString(GuessNumMap["12"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject12)
				case "13":
					if gameGuessOddsNumObject13.GetNum() == "" {
						gameGuessOddsNumObject13.Num = easygo.NewString("13")
					}
					if gameGuessOddsNumObject13.GetNumName() == "" {
						gameGuessOddsNumObject13.NumName = easygo.NewString(GuessNumMap["13"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject13)
				case "14":
					if gameGuessOddsNumObject14.GetNum() == "" {
						gameGuessOddsNumObject14.Num = easygo.NewString("14")
					}
					if gameGuessOddsNumObject14.GetNumName() == "" {
						gameGuessOddsNumObject14.NumName = easygo.NewString(GuessNumMap["14"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject14)
				case "15":
					if gameGuessOddsNumObject15.GetNum() == "" {
						gameGuessOddsNumObject15.Num = easygo.NewString("15")
					}
					if gameGuessOddsNumObject15.GetNumName() == "" {
						gameGuessOddsNumObject15.NumName = easygo.NewString(GuessNumMap["15"])
					}
					GetGameGuessOddsNumObject(guessValue, gameGuessOddsNumObject15)
				default:
				}
			}
			if gameGuessOddsNumObject0.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject0)
			}
			if gameGuessOddsNumObject1.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject1)
			}
			if gameGuessOddsNumObject2.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject2)
			}
			if gameGuessOddsNumObject3.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject3)
			}
			if gameGuessOddsNumObject4.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject4)
			}
			if gameGuessOddsNumObject5.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject5)
			}
			if gameGuessOddsNumObject6.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject6)
			}
			if gameGuessOddsNumObject7.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject7)
			}
			if gameGuessOddsNumObject8.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject8)
			}
			if gameGuessOddsNumObject9.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject9)
			}
			if gameGuessOddsNumObject10.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject10)
			}
			if gameGuessOddsNumObject11.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject11)
			}
			if gameGuessOddsNumObject12.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject12)
			}
			if gameGuessOddsNumObject13.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject13)
			}
			if gameGuessOddsNumObject14.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject14)
			}
			if gameGuessOddsNumObject15.GetNum() != "" {
				obj = append(obj, gameGuessOddsNumObject15)
			}
		}
	}
	gameGuessDetailObject.GuessOddsNums = obj
	return gameGuessDetailObject
}

////开奖结束后、删除redis中的key,存在就删除
//func DelRedisGameKeysForGameOver(uniqueGameId int64) {
//
//	keys := []interface{}{MakeRedisKey(ESPORT_REDIS_GAME_DETAIL_HEAD_KEY, uniqueGameId),
//		MakeRedisKey(ESPORT_REDIS_GAME_GUESS_MORN_KEY, uniqueGameId),
//		MakeRedisKey(ESPORT_REDIS_GAME_GUESS_ROLL_KEY, uniqueGameId)}
//
//	easygo.RedisMgr.GetC().Delete(keys...)
//}

//设置比赛局数共通
func GetGameGuessOddsNumObject(guessValue *share_message.ApiGuessObject, gameGuessOddsNumObject *share_message.GameGuessOddsNumObject) {

	contents := gameGuessOddsNumObject.GetContents()
	if nil == contents {
		contents = make([]*share_message.GameGuessOddsContentObject, 0)
	}

	tempOddsContent := share_message.GameGuessOddsContentObject{}
	tempOddsContent.BetTitle = easygo.NewString(guessValue.GetBetTitle())
	tempOddsContent.BetId = easygo.NewString(guessValue.GetBetId())
	tempOddsContent.AppGuessFlag = easygo.NewInt32(guessValue.GetAppGuessFlag())
	tempOddsContent.AppGuessViewFlag = easygo.NewInt32(guessValue.GetAppGuessViewFlag())

	if nil != guessValue.GetItems() && len(guessValue.GetItems()) > 0 {

		items := make([]*share_message.GameGuessOddsItemObject, 0)

		for _, itemValue := range guessValue.GetItems() {

			tempOddsItem := share_message.GameGuessOddsItemObject{}

			tempOddsItem.BetNum = easygo.NewString(itemValue.GetBetNum())
			tempOddsItem.Status = easygo.NewString(itemValue.GetStatus())
			tempOddsItem.Win = easygo.NewString(itemValue.GetWin())
			tempOddsItem.Odds = easygo.NewString(itemValue.GetOdds())

			tempOddsItem.StatusTime = easygo.NewString(itemValue.GetStatusTime())
			tempOddsItem.ResultTime = easygo.NewString(itemValue.GetResultTime())

			tempBetType := itemValue.GetBetType()
			//竞猜项名称(某些竞猜项目要通过组合计算得到)
			if tempBetType == GAME_GUESS_ITEM_BET_TYPE_1 ||
				tempBetType == GAME_GUESS_ITEM_BET_TYPE_2 ||
				tempBetType == GAME_GUESS_ITEM_BET_TYPE_3 {

				tempOddsItem.BetName = easygo.NewString(itemValue.GetOddsName())
			} else {

				if tempBetType == GAME_GUESS_ITEM_BET_TYPE_5 {
					tempStr := ""
					if itemValue.GetGroupFlag() == GAME_GUESS_ITEM_GROUP_FLAG_1 {
						tempStr = "+"
					} else {
						tempStr = "-"
					}
					tempOddsItem.BetName = easygo.NewString(itemValue.GetOddsName() + tempStr + strings.TrimLeft(itemValue.GetGroupValue(), "-"))

				} else if tempBetType == GAME_GUESS_ITEM_BET_TYPE_6 {

					tempOddsItem.BetName = easygo.NewString(itemValue.GetOddsName() + itemValue.GetGroupValue())

				} else if tempBetType == GAME_GUESS_ITEM_BET_TYPE_7 {
					tempStr := ""
					if itemValue.GetGroupFlag() == GAME_GUESS_ITEM_GROUP_FLAG_1 {
						tempStr = ">"
					} else if itemValue.GetGroupFlag() == GAME_GUESS_ITEM_GROUP_FLAG_2 {
						tempStr = "<"
					}

					tempOddsItem.BetName = easygo.NewString(itemValue.GetOddsName() + tempStr + itemValue.GetGroupValue())
				}
			}

			//后台总控关闭
			if guessValue.GetAppGuessFlag() == GAME_APP_GUESS_FLAG_1 || guessValue.GetAppGuessFlag() == GAME_APP_GUESS_FLAG_0 {
				tempOddsItem.BetStatus = easygo.NewString(GAME_GUESS_ITEM_ODDS_STATUS_2)
			} else {
				if itemValue.GetOddsStatus() == GAME_GUESS_ODDS_STATUS_0 || itemValue.GetOddsStatus() == GAME_GUESS_ODDS_STATUS_3 {
					tempOddsItem.BetStatus = easygo.NewString(GAME_GUESS_ITEM_ODDS_STATUS_2)
				} else if itemValue.GetOddsStatus() == GAME_GUESS_ODDS_STATUS_1 {
					tempOddsItem.BetStatus = easygo.NewString(GAME_GUESS_ITEM_ODDS_STATUS_1)
				}
			}

			//通过Status、Win组合计算出投注项的结果
			if tempOddsItem.GetStatus() != "" {
				if tempOddsItem.GetStatus() == GAME_GUESS_ITEM_STATUS_0 {
					tempOddsItem.Result = easygo.NewString(GAME_GUESS_ITEM_WIN_NORST)
				} else {
					if tempOddsItem.GetWin() != "" {
						if tempOddsItem.GetWin() == GAME_GUESS_ITEM_WIN_NORST {
							tempOddsItem.Result = easygo.NewString(GAME_GUESS_ITEM_WIN_NORST)
						} else if tempOddsItem.GetWin() == GAME_GUESS_ITEM_WIN_0 {
							tempOddsItem.Result = easygo.NewString(GAME_GUESS_ITEM_WIN_0)
						} else {
							tempOddsItem.Result = easygo.NewString(GAME_GUESS_ITEM_WIN_1)
						}
					} else {
						tempOddsItem.Result = easygo.NewString(GAME_GUESS_ITEM_WIN_NORST)
					}

				}
			} else {
				tempOddsItem.Result = easygo.NewString(GAME_GUESS_ITEM_WIN_NORST)
			}

			items = append(items, &tempOddsItem)
		}
		tempOddsContent.Items = items
	}
	contents = append(contents, &tempOddsContent)
	gameGuessOddsNumObject.Contents = contents
}

//过滤掉不显示的投注项目
func GetFilterGameDetailBet(detailObject *share_message.GameGuessDetailObject) {
	if detailObject.GetGuessOddsNums() != nil && len(detailObject.GetGuessOddsNums()) > 0 {

		oddsNumObjects := make([]*share_message.GameGuessOddsNumObject, 0)

		for _, value := range detailObject.GetGuessOddsNums() {
			oddsNumObject := share_message.GameGuessOddsNumObject{}
			oddsNumObject.Num = easygo.NewString(value.GetNum())
			oddsNumObject.NumName = easygo.NewString(value.GetNumName())

			rdContents := value.GetContents()

			conents := make([]*share_message.GameGuessOddsContentObject, 0)
			if rdContents != nil && len(rdContents) > 0 {
				for _, cntValue := range rdContents {
					if cntValue.GetAppGuessViewFlag() == GAME_APP_GUESS_VIEW_FLAG_2 {
						conents = append(conents, cntValue)
					}
				}
			}

			if nil != conents && len(conents) > 0 {
				oddsNumObject.Contents = conents

				oddsNumObjects = append(oddsNumObjects, &oddsNumObject)
			}
		}

		detailObject.GuessOddsNums = oddsNumObjects
	}
}
