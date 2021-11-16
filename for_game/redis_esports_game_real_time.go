package for_game

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"strconv"
)

const (
	ESPORT_REDIS_GAME_HERO_KEY      = "redis_esport:game_hero"      //游戏英雄(加appLabelId)(常驻key、启动的时候读取json)
	ESPORT_REDIS_GAME_EQUIPMENT_KEY = "redis_esport:game_equipment" //游戏装备(加appLabelId)(常驻key、启动的时候读取json)

	ESPORT_REDIS_GAME_REAL_TIME_KEY        = "redis_esport:game_real_time"        //游戏每局实时数据(比赛id)(3小时)
	ESPORT_REDIS_GAME_REAL_TIME_ROUNDS_KEY = "redis_esport:game_real_time_rounds" //游戏实时数据的总局数(比赛id)(3小时)

)

//启动时候设置游戏英雄、装备数据(LOL、王者荣耀)
func SetRedisGameRealTimeBase() {

	var jsonFileHeroLOL string
	var jsonFileHeroWZRY string
	var jsonFileEquipLOL string
	var jsonFileEquipWZRY string

	if IS_FORMAL_SERVER {
		jsonFileHeroLOL = easygo.YamlCfg.GetValueAsString("YEZI_HERO_LOL_DATA_FORMAL")
		jsonFileHeroWZRY = easygo.YamlCfg.GetValueAsString("YEZI_HERO_WZRY_DATA_FORMAL")

		jsonFileEquipLOL = easygo.YamlCfg.GetValueAsString("YEZI_EQUIP_LOL_DATA_FORMAL")
		jsonFileEquipWZRY = easygo.YamlCfg.GetValueAsString("YEZI_EQUIP_WZRY_DATA_FORMAL")
	} else {
		jsonFileHeroLOL = easygo.YamlCfg.GetValueAsString("YEZI_HERO_LOL_DATA_TEST")
		jsonFileHeroWZRY = easygo.YamlCfg.GetValueAsString("YEZI_HERO_WZRY_DATA_TEST")

		jsonFileEquipLOL = easygo.YamlCfg.GetValueAsString("YEZI_EQUIP_LOL_DATA_TEST")
		jsonFileEquipWZRY = easygo.YamlCfg.GetValueAsString("YEZI_EQUIP_WZRY_DATA_TEST")
	}

	dataHeroLOL := FindHeroInfos(jsonFileHeroLOL)
	dataHeroWZRY := FindHeroInfos(jsonFileHeroWZRY)

	dataEquipLOL := FindEquipInfos(jsonFileEquipLOL)

	dataEquipWZRY := FindEquipInfos(jsonFileEquipWZRY)

	if nil != dataHeroLOL && len(dataHeroLOL) > 0 {
		//json序列化
		data, errJs := json.Marshal(&dataHeroLOL)
		easygo.PanicError(errJs)
		//redis 常驻key
		err := easygo.RedisMgr.GetC().Set(
			MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_LOL),
			string(data))
		easygo.PanicError(err)
	} else {
		keys := []interface{}{MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_LOL)}

		easygo.RedisMgr.GetC().Delete(keys...)

		s := fmt.Sprintf("LOL游戏的英雄基础json未配置正确")
		logs.Error(s)
	}

	if nil != dataEquipLOL && len(dataEquipLOL) > 0 {
		//json序列化
		data, errJs := json.Marshal(&dataEquipLOL)
		easygo.PanicError(errJs)
		//redis 常驻key
		err := easygo.RedisMgr.GetC().Set(
			MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_LOL),
			string(data))
		easygo.PanicError(err)
	} else {

		keys := []interface{}{MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_LOL)}

		easygo.RedisMgr.GetC().Delete(keys...)
		s := fmt.Sprintf("LOL游戏的英雄的装备基础json未配置正确")
		logs.Error(s)
	}

	if nil != dataHeroWZRY && len(dataHeroWZRY) > 0 {
		//json序列化
		data, errJs := json.Marshal(&dataHeroWZRY)
		easygo.PanicError(errJs)
		//redis 常驻key
		err := easygo.RedisMgr.GetC().Set(
			MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_WZRY),
			string(data))
		easygo.PanicError(err)
	} else {
		keys := []interface{}{MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_WZRY)}

		easygo.RedisMgr.GetC().Delete(keys...)
		s := fmt.Sprintf("王者荣耀游戏的英雄基础json未配置正确")
		logs.Error(s)
	}

	if nil != dataEquipWZRY && len(dataEquipWZRY) > 0 {
		//json序列化
		data, errJs := json.Marshal(&dataEquipWZRY)
		easygo.PanicError(errJs)
		//redis 常驻key
		err := easygo.RedisMgr.GetC().Set(
			MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_WZRY),
			string(data))
		easygo.PanicError(err)
	} else {
		keys := []interface{}{MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_WZRY)}

		easygo.RedisMgr.GetC().Delete(keys...)
		s := fmt.Sprintf("王者荣耀游戏的英雄的装备基础json未配置正确")
		logs.Error(s)
	}
}

//取得redis游戏的英雄的信息
func GetRedisGameRealTimeHero(heroId int32, appLabelId int32) *share_message.RealTimeHeroObject {

	heroIdStr := fmt.Sprintf("%v", heroId)

	var heroObject *share_message.RealTimeHeroObject

	if appLabelId == ESPORTS_LABEL_LOL {
		b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_LOL))
		easygo.PanicError(err)

		if !b {

			//再设置一次redis
			SetRedisGameRealTimeBase()
			b1, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_LOL))
			easygo.PanicError(err)
			if !b1 {
				heroObject = nil
			} else {
				obj := make(map[string]*share_message.RealTimeHeroObject)
				value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_LOL))
				easygo.PanicError(err)
				errJs := json.Unmarshal([]byte(value), &obj)
				easygo.PanicError(errJs)

				if nil != obj && len(obj) > 0 {
					if IS_FORMAL_SERVER {
						heroObject = obj[heroIdStr]
					} else {
						heroObject = GetRandomHeroTest(obj)
					}
				}
			}
		} else {
			obj := make(map[string]*share_message.RealTimeHeroObject)
			value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_LOL))
			easygo.PanicError(err)
			errJs := json.Unmarshal([]byte(value), &obj)
			easygo.PanicError(errJs)
			if nil != obj && len(obj) > 0 {
				if IS_FORMAL_SERVER {
					heroObject = obj[heroIdStr]
				} else {
					heroObject = GetRandomHeroTest(obj)
				}
			}
		}
	} else if appLabelId == ESPORTS_LABEL_WZRY {
		b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_WZRY))
		easygo.PanicError(err)

		if !b {

			//再设置一次redis
			SetRedisGameRealTimeBase()
			b1, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_WZRY))
			easygo.PanicError(err)
			if !b1 {
				heroObject = nil
			} else {
				obj := make(map[string]*share_message.RealTimeHeroObject)
				value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_WZRY))
				easygo.PanicError(err)
				errJs := json.Unmarshal([]byte(value), &obj)
				easygo.PanicError(errJs)

				if nil != obj && len(obj) > 0 {
					if IS_FORMAL_SERVER {
						heroObject = obj[heroIdStr]
					} else {
						heroObject = GetRandomHeroTest(obj)
					}
				}
			}

		} else {
			obj := make(map[string]*share_message.RealTimeHeroObject)
			value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_HERO_KEY, ESPORTS_LABEL_WZRY))
			easygo.PanicError(err)
			errJs := json.Unmarshal([]byte(value), &obj)
			easygo.PanicError(errJs)

			if nil != obj && len(obj) > 0 {
				if IS_FORMAL_SERVER {
					heroObject = obj[heroIdStr]
				} else {
					heroObject = GetRandomHeroTest(obj)
				}
			}
		}
	}

	return heroObject
}

//取得redis游戏的装备的信息
func GetRedisGameRealTimeEquip(itemId int32, appLabelId int32) *share_message.RealTimeItemObject {

	itemIdStr := fmt.Sprintf("%v", itemId)

	var ItemObject *share_message.RealTimeItemObject

	if appLabelId == ESPORTS_LABEL_LOL {
		b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_LOL))
		easygo.PanicError(err)

		if !b {

			//再设置一次redis
			SetRedisGameRealTimeBase()
			b1, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_LOL))
			easygo.PanicError(err)
			if !b1 {
				ItemObject = nil
			} else {
				obj := make(map[string]*share_message.RealTimeItemObject)
				value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_LOL))
				easygo.PanicError(err)
				errJs := json.Unmarshal([]byte(value), &obj)
				easygo.PanicError(errJs)

				if nil != obj && len(obj) > 0 {
					if IS_FORMAL_SERVER {
						ItemObject = obj[itemIdStr]
					} else {
						ItemObject = GetRandomEquipTest(obj)
					}
				}
			}

		} else {
			obj := make(map[string]*share_message.RealTimeItemObject)
			value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_LOL))
			easygo.PanicError(err)
			errJs := json.Unmarshal([]byte(value), &obj)
			easygo.PanicError(errJs)

			if nil != obj && len(obj) > 0 {
				if IS_FORMAL_SERVER {
					ItemObject = obj[itemIdStr]
				} else {
					ItemObject = GetRandomEquipTest(obj)
				}
			}
		}
	} else if appLabelId == ESPORTS_LABEL_WZRY {
		b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_WZRY))
		easygo.PanicError(err)

		if !b {

			//再设置一次redis
			SetRedisGameRealTimeBase()
			b1, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_WZRY))
			easygo.PanicError(err)
			if !b1 {
				ItemObject = nil
			} else {
				obj := make(map[string]*share_message.RealTimeItemObject)
				value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_WZRY))
				easygo.PanicError(err)
				errJs := json.Unmarshal([]byte(value), &obj)
				easygo.PanicError(errJs)

				if nil != obj && len(obj) > 0 {
					if IS_FORMAL_SERVER {
						ItemObject = obj[itemIdStr]
					} else {
						ItemObject = GetRandomEquipTest(obj)
					}
				}
			}

		} else {
			obj := make(map[string]*share_message.RealTimeItemObject)
			value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_EQUIPMENT_KEY, ESPORTS_LABEL_WZRY))
			easygo.PanicError(err)
			errJs := json.Unmarshal([]byte(value), &obj)
			easygo.PanicError(errJs)

			if nil != obj && len(obj) > 0 {
				if IS_FORMAL_SERVER {
					ItemObject = obj[itemIdStr]
				} else {
					ItemObject = GetRandomEquipTest(obj)
				}
			}
		}
	}

	return ItemObject
}

//设置redis游戏每局实时数据(LOL、王者荣耀)
func SetRedisGameRealTime(uniqueGameId int64, gameRound int32) {

	gameRealTime := GetDBGameRealTime(uniqueGameId, gameRound)
	if nil != gameRealTime {
		//json序列化
		data, errJs := json.Marshal(gameRealTime)
		easygo.PanicError(errJs)
		//redis 3小时
		err := easygo.RedisMgr.GetC().SetWithTime(
			MakeRedisKey(ESPORT_REDIS_GAME_REAL_TIME_KEY, uniqueGameId, gameRound),
			string(data),
			ESPORT_GAME_REDIS_EXPIRE_TIME)
		easygo.PanicError(err)
	} else {
		keys := []interface{}{MakeRedisKey(ESPORT_REDIS_GAME_REAL_TIME_KEY, uniqueGameId, gameRound)}

		easygo.RedisMgr.GetC().Delete(keys...)
	}
}

//取得redis游戏每局实时数据(LOL、王者荣耀)
func GetRedisGameRealTime(uniqueGameId int64, gameRound int32) *client_hall.GameRealTimeData {

	b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_REAL_TIME_KEY, uniqueGameId, gameRound))
	easygo.PanicError(err)

	var gameRealTimeData *client_hall.GameRealTimeData
	if !b {
		gameRealTimeData = GetDBGameRealTime(uniqueGameId, gameRound)
		if nil != gameRealTimeData {
			//json序列化
			data, errJs := json.Marshal(gameRealTimeData)
			easygo.PanicError(errJs)
			//redis 3小时
			err := easygo.RedisMgr.GetC().SetWithTime(
				MakeRedisKey(ESPORT_REDIS_GAME_REAL_TIME_KEY,
					uniqueGameId, gameRound), string(data), ESPORT_GAME_REDIS_EXPIRE_TIME)
			easygo.PanicError(err)
		}
	} else {
		obj := &client_hall.GameRealTimeData{}
		value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_REAL_TIME_KEY, uniqueGameId, gameRound))
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), obj)
		easygo.PanicError(errJs)

		gameRealTimeData = obj
	}
	return gameRealTimeData
}

//设置redis游戏实时数据总局数(LOL、王者荣耀)
func SetRedisGameRealTimeRounds(apiOrigin int32, appLabelId int32, gameId int32, gameObject *client_hall.ESportGameObject) {

	rounds := GetDBGameRealTimeRounds(apiOrigin, appLabelId, gameId)

	if rounds > 0 {
		tmpGameRounds := &client_hall.GameRounds{GameRounds: easygo.NewInt32(rounds)}
		//json序列化
		data, errJs := json.Marshal(tmpGameRounds)
		easygo.PanicError(errJs)
		//redis 3小时
		err := easygo.RedisMgr.GetC().SetWithTime(
			MakeRedisKey(ESPORT_REDIS_GAME_REAL_TIME_ROUNDS_KEY, gameObject.GetUniqueGameId()),
			string(data),
			ESPORT_GAME_REDIS_EXPIRE_TIME)
		easygo.PanicError(err)
	} else {
		keys := []interface{}{MakeRedisKey(ESPORT_REDIS_GAME_REAL_TIME_ROUNDS_KEY, gameObject.GetUniqueGameId())}

		easygo.RedisMgr.GetC().Delete(keys...)
	}
}

//取得redis实时数据总局数(LOL、王者荣耀)
func GetRedisGameRealTimeRounds(uniqueGameId int64) *client_hall.GameRounds {

	b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(ESPORT_REDIS_GAME_REAL_TIME_ROUNDS_KEY, uniqueGameId))
	easygo.PanicError(err)

	var rounds *client_hall.GameRounds
	if !b {
		gameDetailHead := GetRedisGameDetailHead(uniqueGameId)
		if nil != gameDetailHead {

			gameId, err := strconv.ParseInt(gameDetailHead.GetGameId(), 10, 32)

			if err != nil {
				logs.Error(err)
				s := fmt.Sprintf("=======GetRedisGameRealTimeRounds处理游戏实时数据参数gameDetailHead.GetGameId()错误==========:参数gameDetailHead.GetGameId()=%v", gameDetailHead.GetGameId())
				logs.Error(s)
				return nil
			}
			tempRounds := GetDBGameRealTimeRounds(gameDetailHead.GetApiOrigin(), gameDetailHead.GetAppLabelId(), int32(gameId))

			if tempRounds > 0 {
				rounds = &client_hall.GameRounds{GameRounds: easygo.NewInt32(tempRounds)}
				//json序列化
				data, errJs := json.Marshal(rounds)
				easygo.PanicError(errJs)
				//redis 3小时
				err1 := easygo.RedisMgr.GetC().SetWithTime(
					MakeRedisKey(ESPORT_REDIS_GAME_REAL_TIME_ROUNDS_KEY, uniqueGameId),
					string(data),
					ESPORT_GAME_REDIS_EXPIRE_TIME)
				easygo.PanicError(err1)
			}
		}
	} else {
		obj := &client_hall.GameRounds{}
		value, err := easygo.RedisMgr.GetC().Get(MakeRedisKey(ESPORT_REDIS_GAME_REAL_TIME_ROUNDS_KEY, uniqueGameId))
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), obj)
		easygo.PanicError(errJs)

		rounds = obj
	}
	return rounds
}

//数据库通过uniqueGameId和gameRound取得游戏实时数据(LOL、王者荣耀)
func GetDBGameRealTime(uniqueGameId int64, gameRound int32) *client_hall.GameRealTimeData {

	var rst *client_hall.GameRealTimeData

	//从redis取得比赛的数据 、理论上比赛是不能未nil的
	gameObject := GetRedisGameDetailHead(uniqueGameId)
	//得到比赛相关的信息
	if nil == gameObject {
		s := fmt.Sprintf("取得比赛uniqueGameId:%v数据库实时数据时候、比赛信息丢失", uniqueGameId)
		logs.Error(s)
		return nil
	}

	appLabelId := gameObject.GetAppLabelId()
	apiOrigin := gameObject.GetApiOrigin()
	gameIdStr := gameObject.GetGameId()
	//LOL
	if appLabelId == ESPORTS_LABEL_LOL {

		col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_LOL_REAL_TIME_DATA)
		defer closeFun()

		gameId, err := strconv.ParseInt(gameIdStr, 10, 32)
		if err != nil {
			logs.Error(err)
			s := fmt.Sprintf("=======GetDBGameRealTime处理LOL游戏转换实时数据参数gameIdStr错误==========:gameIdStr=%v", gameIdStr)
			logs.Error(s)
			return nil
		}

		query := share_message.TableESPortsLOLRealTimeData{}
		//通过条件查询存在更新,不存在就插入
		errQuery := col.Find(bson.M{"app_label_id": appLabelId,
			"game_id":    int32(gameId),
			"api_origin": apiOrigin,
			"game_round": gameRound}).One(&query)

		if errQuery != nil && errQuery != mgo.ErrNotFound {
			logs.Error(errQuery)
			s := fmt.Sprintf("======GetDBGameRealTime LOL游戏实时数据表查询失败======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v,====game_round:%v",
				appLabelId, gameId, apiOrigin, gameRound)
			logs.Error(s)
			return nil
		}

		if errQuery == mgo.ErrNotFound {
			//因为轮询不打log
			//logs.Info(errQuery)
			//s := fmt.Sprintf("======GetDBGameRealTime LOL游戏实时数据表查询没有数据======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v,====game_round:%v",
			//	appLabelId, gameId, apiOrigin, gameRound)
			//logs.Info(s)
			return nil
		}

		if nil == errQuery {
			rst = &client_hall.GameRealTimeData{
				GameRound:        easygo.NewInt32(query.GetGameRound()),
				GameStatus:       easygo.NewInt32(query.GetGameStatus()),
				Duration:         easygo.NewInt32(query.GetDuration()),
				FirstTower:       easygo.NewInt32(query.GetFirstTower()),
				FirstSmallDragon: easygo.NewInt32(query.GetFirstSmallDragon()),
				FirstFiveKill:    easygo.NewInt32(query.GetFirstFiveKill()),
				FirstBigDragon:   easygo.NewInt32(query.GetFirstBigDragon()),
				FirstTenKill:     easygo.NewInt32(query.GetFirstTenKill()),
			}

			//战队a信息
			dbTeamA := query.GetTeamA()
			if nil != dbTeamA {
				teamA := client_hall.RealTimeTeamObject{
					Faction:      easygo.NewString(dbTeamA.GetFaction()),
					Score:        easygo.NewInt32(dbTeamA.GetScore()),
					Glod:         easygo.NewInt32(dbTeamA.GetGlod()),
					TowerState:   easygo.NewInt32(dbTeamA.GetTowerState()),
					Drakes:       easygo.NewInt32(dbTeamA.GetDrakes()),
					NahsorBarons: easygo.NewInt32(dbTeamA.GetNahsorBarons()),
					GoldTimeData: dbTeamA.GetGoldTimeData(),
				}

				//队伍A中选取的英雄设置
				if dbTeamA.GetPicks() != nil && len(dbTeamA.GetPicks()) > 0 {
					tempPickInfos := make([]*share_message.RealTimeHeroObject, 0)
					for _, pickValue := range dbTeamA.GetPicks() {
						tempPickHero := GetRedisGameRealTimeHero(pickValue, ESPORTS_LABEL_LOL)
						if nil != tempPickHero {
							tempPickInfos = append(tempPickInfos, tempPickHero)
						}
					}

					teamA.PickInfos = tempPickInfos
				}

				//队伍A中禁止的英雄设置
				if dbTeamA.GetBans() != nil && len(dbTeamA.GetBans()) > 0 {
					tempBanInfos := make([]*share_message.RealTimeHeroObject, 0)
					for _, banValue := range dbTeamA.GetBans() {

						tempBanHero := GetRedisGameRealTimeHero(banValue, ESPORTS_LABEL_LOL)
						if nil != tempBanHero {
							tempBanInfos = append(tempBanInfos, tempBanHero)
						}
					}

					teamA.BanInfos = tempBanInfos
				}

				rst.TeamA = &teamA
			}

			//战队b信息
			dbTeamB := query.GetTeamB()
			if nil != dbTeamB {
				teamB := client_hall.RealTimeTeamObject{
					Faction:      easygo.NewString(dbTeamB.GetFaction()),
					Score:        easygo.NewInt32(dbTeamB.GetScore()),
					Glod:         easygo.NewInt32(dbTeamB.GetGlod()),
					TowerState:   easygo.NewInt32(dbTeamB.GetTowerState()),
					Drakes:       easygo.NewInt32(dbTeamB.GetDrakes()),
					NahsorBarons: easygo.NewInt32(dbTeamB.GetNahsorBarons()),
					GoldTimeData: dbTeamB.GetGoldTimeData(),
				}

				//队伍B中选取的英雄设置
				if dbTeamB.GetPicks() != nil && len(dbTeamB.GetPicks()) > 0 {
					tempPickInfos := make([]*share_message.RealTimeHeroObject, 0)
					for _, pickValue := range dbTeamB.GetPicks() {

						tempPickHero := GetRedisGameRealTimeHero(pickValue, ESPORTS_LABEL_LOL)
						if nil != tempPickHero {
							tempPickInfos = append(tempPickInfos, tempPickHero)
						}
					}
					teamB.PickInfos = tempPickInfos
				}

				//队伍B中禁止的英雄设置
				if dbTeamB.GetBans() != nil && len(dbTeamB.GetBans()) > 0 {
					tempBanInfos := make([]*share_message.RealTimeHeroObject, 0)
					for _, banValue := range dbTeamB.GetBans() {

						tempBanHero := GetRedisGameRealTimeHero(banValue, ESPORTS_LABEL_LOL)
						if nil != tempBanHero {
							tempBanInfos = append(tempBanInfos, tempBanHero)
						}
					}
					teamB.BanInfos = tempBanInfos
				}
				rst.TeamB = &teamB
			}

			//战队A队员信息设置
			dbPlayerAInfos := query.GetPlayerAInfo()
			if dbPlayerAInfos != nil {
				tempPlayerAInfos := make([]*client_hall.RealTimePlayerObject, 0)

				for _, dbPlayA := range dbPlayerAInfos {
					tmpPlayA := client_hall.RealTimePlayerObject{
						Name:    easygo.NewString(dbPlayA.GetName()),
						Kills:   easygo.NewInt32(dbPlayA.GetKills()),
						Death:   easygo.NewInt32(dbPlayA.GetDeath()),
						Assists: easygo.NewInt32(dbPlayA.GetAssists()),
						Gold:    easygo.NewInt32(dbPlayA.GetGold()),
						Subsidy: easygo.NewInt32(dbPlayA.GetSubsidy()),
						Photo:   easygo.NewString(dbPlayA.GetPhoto()),
					}

					//设置选取的英雄
					tmpPlayA.HeroInfo = GetRedisGameRealTimeHero(dbPlayA.GetHeroId(), ESPORTS_LABEL_LOL)

					//设置装备栏道具
					if dbPlayA.GetItem() != nil && len(dbPlayA.GetItem()) > 0 {
						tempItemInfos := make([]*share_message.RealTimeItemObject, 0)
						for _, itemValue := range dbPlayA.GetItem() {

							tempItem := GetRedisGameRealTimeEquip(itemValue, ESPORTS_LABEL_LOL)
							if nil != tempItem {
								tempItemInfos = append(tempItemInfos, tempItem)
							}
						}
						tmpPlayA.ItemInfos = tempItemInfos
					}
					tempPlayerAInfos = append(tempPlayerAInfos, &tmpPlayA)
				}
				rst.PlayerAInfo = tempPlayerAInfos
			}

			//战队B队员信息设置
			dbPlayerBInfos := query.GetPlayerBInfo()
			if dbPlayerBInfos != nil {
				tempPlayerBInfos := make([]*client_hall.RealTimePlayerObject, 0)

				for _, dbPlayB := range dbPlayerBInfos {
					tmpPlayB := client_hall.RealTimePlayerObject{
						Name:    easygo.NewString(dbPlayB.GetName()),
						Kills:   easygo.NewInt32(dbPlayB.GetKills()),
						Death:   easygo.NewInt32(dbPlayB.GetDeath()),
						Assists: easygo.NewInt32(dbPlayB.GetAssists()),
						Gold:    easygo.NewInt32(dbPlayB.GetGold()),
						Subsidy: easygo.NewInt32(dbPlayB.GetSubsidy()),
						Photo:   easygo.NewString(dbPlayB.GetPhoto()),
					}

					//设置选取的英雄
					tmpPlayB.HeroInfo = GetRedisGameRealTimeHero(dbPlayB.GetHeroId(), ESPORTS_LABEL_LOL)

					//设置装备栏道具
					if dbPlayB.GetItem() != nil && len(dbPlayB.GetItem()) > 0 {
						tempItemInfos := make([]*share_message.RealTimeItemObject, 0)
						for _, itemValue := range dbPlayB.GetItem() {
							tempItem := GetRedisGameRealTimeEquip(itemValue, ESPORTS_LABEL_LOL)
							if nil != tempItem {
								tempItemInfos = append(tempItemInfos, tempItem)
							}
						}
						tmpPlayB.ItemInfos = tempItemInfos
					}
					tempPlayerBInfos = append(tempPlayerBInfos, &tmpPlayB)
				}
				rst.PlayerBInfo = tempPlayerBInfos
			}
		}
		return rst
		//王者荣耀
	} else if appLabelId == ESPORTS_LABEL_WZRY {
		col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_WZRY_REAL_TIME_DATA)
		defer closeFun()

		gameId, err := strconv.ParseInt(gameIdStr, 10, 32)
		if err != nil {
			logs.Error(err)
			s := fmt.Sprintf("=======GetDBGameRealTime处理WZRY游戏转换实时数据参数gameIdStr错误==========:gameIdStr=%v", gameIdStr)
			logs.Error(s)
			return nil
		}

		query := share_message.TableESPortsWZRYRealTimeData{}
		//通过条件查询存在更新,不存在就插入
		errQuery := col.Find(bson.M{"app_label_id": appLabelId,
			"game_id":    int32(gameId),
			"api_origin": apiOrigin,
			"game_round": gameRound}).One(&query)

		if errQuery != nil && errQuery != mgo.ErrNotFound {
			logs.Error(errQuery)
			s := fmt.Sprintf("======GetDBGameRealTime WZRY游戏实时数据表查询失败======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v,====game_round:%v",
				appLabelId, gameId, apiOrigin, gameRound)
			logs.Error(s)
			return nil
		}

		if errQuery == mgo.ErrNotFound {
			//因为轮询不打log
			//logs.Info(errQuery)
			//s := fmt.Sprintf("======GetDBGameRealTime WZRY游戏实时数据表查询没有数据======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v,====game_round:%v",
			//	appLabelId, gameId, apiOrigin, gameRound)
			//logs.Info(s)
			return nil
		}

		if nil == errQuery {
			rst = &client_hall.GameRealTimeData{
				GameRound:        easygo.NewInt32(query.GetGameRound()),
				GameStatus:       easygo.NewInt32(query.GetGameStatus()),
				Duration:         easygo.NewInt32(query.GetDuration()),
				FirstTower:       easygo.NewInt32(query.GetFirstTower()),
				FirstSmallDragon: easygo.NewInt32(query.GetFirstSmallDragon()),
				FirstFiveKill:    easygo.NewInt32(query.GetFirstFiveKill()),
				FirstBigDragon:   easygo.NewInt32(query.GetFirstBigDragon()),
				FirstTenKill:     easygo.NewInt32(query.GetFirstTenKill()),
			}

			//战队a信息
			dbTeamA := query.GetTeamA()
			if nil != dbTeamA {
				teamA := client_hall.RealTimeTeamObject{
					Faction:      easygo.NewString(dbTeamA.GetFaction()),
					Score:        easygo.NewInt32(dbTeamA.GetScore()),
					Glod:         easygo.NewInt32(dbTeamA.GetGlod()),
					TowerState:   easygo.NewInt32(dbTeamA.GetTowerState()),
					Drakes:       easygo.NewInt32(dbTeamA.GetDrakes()),
					NahsorBarons: easygo.NewInt32(dbTeamA.GetNahsorBarons()),
					GoldTimeData: dbTeamA.GetGoldTimeData(),
				}

				//队伍A中选取的英雄设置
				if dbTeamA.GetPicks() != nil && len(dbTeamA.GetPicks()) > 0 {
					tempPickInfos := make([]*share_message.RealTimeHeroObject, 0)
					for _, pickValue := range dbTeamA.GetPicks() {

						tempPickHero := GetRedisGameRealTimeHero(pickValue, ESPORTS_LABEL_WZRY)
						if nil != tempPickHero {
							tempPickInfos = append(tempPickInfos, tempPickHero)
						}
					}

					teamA.PickInfos = tempPickInfos
				}

				//队伍A中禁止的英雄设置
				if dbTeamA.GetBans() != nil && len(dbTeamA.GetBans()) > 0 {
					tempBanInfos := make([]*share_message.RealTimeHeroObject, 0)
					for _, banValue := range dbTeamA.GetBans() {

						tempBanHero := GetRedisGameRealTimeHero(banValue, ESPORTS_LABEL_WZRY)
						if nil != tempBanHero {
							tempBanInfos = append(tempBanInfos, tempBanHero)
						}
					}

					teamA.BanInfos = tempBanInfos
				}

				rst.TeamA = &teamA
			}

			//战队b信息
			dbTeamB := query.GetTeamB()
			if nil != dbTeamB {
				teamB := client_hall.RealTimeTeamObject{
					Faction:      easygo.NewString(dbTeamB.GetFaction()),
					Score:        easygo.NewInt32(dbTeamB.GetScore()),
					Glod:         easygo.NewInt32(dbTeamB.GetGlod()),
					TowerState:   easygo.NewInt32(dbTeamB.GetTowerState()),
					Drakes:       easygo.NewInt32(dbTeamB.GetDrakes()),
					NahsorBarons: easygo.NewInt32(dbTeamB.GetNahsorBarons()),
					GoldTimeData: dbTeamB.GetGoldTimeData(),
				}

				//队伍B中选取的英雄设置
				if dbTeamB.GetPicks() != nil && len(dbTeamB.GetPicks()) > 0 {
					tempPickInfos := make([]*share_message.RealTimeHeroObject, 0)
					for _, pickValue := range dbTeamB.GetPicks() {
						tempPickHero := GetRedisGameRealTimeHero(pickValue, ESPORTS_LABEL_WZRY)
						if nil != tempPickHero {
							tempPickInfos = append(tempPickInfos, tempPickHero)
						}
					}
					teamB.PickInfos = tempPickInfos
				}

				//队伍B中禁止的英雄设置
				if dbTeamB.GetBans() != nil && len(dbTeamB.GetBans()) > 0 {
					tempBanInfos := make([]*share_message.RealTimeHeroObject, 0)
					for _, banValue := range dbTeamB.GetBans() {

						tempBanHero := GetRedisGameRealTimeHero(banValue, ESPORTS_LABEL_WZRY)
						if nil != tempBanHero {
							tempBanInfos = append(tempBanInfos, tempBanHero)
						}
					}
					teamB.BanInfos = tempBanInfos
				}
				rst.TeamB = &teamB
			}

			//战队A队员信息设置
			dbPlayerAInfos := query.GetPlayerAInfo()
			if dbPlayerAInfos != nil {
				tempPlayerAInfos := make([]*client_hall.RealTimePlayerObject, 0)

				for _, dbPlayA := range dbPlayerAInfos {
					tmpPlayA := client_hall.RealTimePlayerObject{
						Name:    easygo.NewString(dbPlayA.GetName()),
						Kills:   easygo.NewInt32(dbPlayA.GetKills()),
						Death:   easygo.NewInt32(dbPlayA.GetDeath()),
						Assists: easygo.NewInt32(dbPlayA.GetAssists()),
						Gold:    easygo.NewInt32(dbPlayA.GetGold()),
						Photo:   easygo.NewString(dbPlayA.GetPhoto()),
					}

					//设置选取的英雄
					tmpPlayA.HeroInfo = GetRedisGameRealTimeHero(dbPlayA.GetHeroId(), ESPORTS_LABEL_WZRY)

					//设置装备栏道具
					if dbPlayA.GetItem() != nil && len(dbPlayA.GetItem()) > 0 {
						tempItemInfos := make([]*share_message.RealTimeItemObject, 0)
						for _, itemValue := range dbPlayA.GetItem() {
							tempItem := GetRedisGameRealTimeEquip(itemValue, ESPORTS_LABEL_WZRY)
							if nil != tempItem {
								tempItemInfos = append(tempItemInfos, tempItem)
							}
						}
						tmpPlayA.ItemInfos = tempItemInfos
					}
					tempPlayerAInfos = append(tempPlayerAInfos, &tmpPlayA)
				}
				rst.PlayerAInfo = tempPlayerAInfos
			}

			//战队B队员信息设置
			dbPlayerBInfos := query.GetPlayerBInfo()
			if dbPlayerBInfos != nil {
				tempPlayerBInfos := make([]*client_hall.RealTimePlayerObject, 0)

				for _, dbPlayB := range dbPlayerBInfos {
					tmpPlayB := client_hall.RealTimePlayerObject{
						Name:    easygo.NewString(dbPlayB.GetName()),
						Kills:   easygo.NewInt32(dbPlayB.GetKills()),
						Death:   easygo.NewInt32(dbPlayB.GetDeath()),
						Assists: easygo.NewInt32(dbPlayB.GetAssists()),
						Gold:    easygo.NewInt32(dbPlayB.GetGold()),
						Photo:   easygo.NewString(dbPlayB.GetPhoto()),
					}

					//设置选取的英雄
					tmpPlayB.HeroInfo = GetRedisGameRealTimeHero(dbPlayB.GetHeroId(), ESPORTS_LABEL_WZRY)

					//设置装备栏道具
					if dbPlayB.GetItem() != nil && len(dbPlayB.GetItem()) > 0 {
						tempItemInfos := make([]*share_message.RealTimeItemObject, 0)
						for _, itemValue := range dbPlayB.GetItem() {
							tempItem := GetRedisGameRealTimeEquip(itemValue, ESPORTS_LABEL_WZRY)
							if nil != tempItem {
								tempItemInfos = append(tempItemInfos, tempItem)
							}
						}
						tmpPlayB.ItemInfos = tempItemInfos
					}
					tempPlayerBInfos = append(tempPlayerBInfos, &tmpPlayB)
				}
				rst.PlayerBInfo = tempPlayerBInfos
			}
		}
		return rst
	}
	return rst
}

//数据库取得游戏实时数据的总局数(LOL、王者荣耀)
func GetDBGameRealTimeRounds(apiOrigin int32, appLabelId int32, gameId int32) int32 {

	var cnt int32

	if appLabelId == ESPORTS_LABEL_LOL {
		col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_LOL_REAL_TIME_DATA)
		defer closeFun()

		//通过条件查询
		query := col.Find(bson.M{"app_label_id": appLabelId,
			"game_id":    gameId,
			"api_origin": apiOrigin})

		count, err := query.Count()

		if err != nil {
			logs.Error(err)
			return cnt
		}

		cnt = int32(count)
		return cnt

	} else if appLabelId == ESPORTS_LABEL_WZRY {
		col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_WZRY_REAL_TIME_DATA)
		defer closeFun()

		//通过条件查询
		query := col.Find(bson.M{"app_label_id": appLabelId,
			"game_id":    gameId,
			"api_origin": apiOrigin})

		count, err := query.Count()

		if err != nil {
			logs.Error(err)
			return cnt
		}
		cnt = int32(count)
		return cnt
	}
	return cnt
}

//测试环境中随机取得英雄
func GetRandomHeroTest(heroMap map[string]*share_message.RealTimeHeroObject) *share_message.RealTimeHeroObject {
	mapKeys := make([]string, 0, len(heroMap))
	for key := range heroMap {
		mapKeys = append(mapKeys, key)
	}
	tempKey := mapKeys[rand.Intn(len(mapKeys))]

	return heroMap[tempKey]
}

//测试环境中随机取得装备
func GetRandomEquipTest(itemMap map[string]*share_message.RealTimeItemObject) *share_message.RealTimeItemObject {
	mapKeys := make([]string, 0, len(itemMap))
	for key := range itemMap {
		mapKeys = append(mapKeys, key)
	}
	tempKey := mapKeys[rand.Intn(len(mapKeys))]

	return itemMap[tempKey]
}
