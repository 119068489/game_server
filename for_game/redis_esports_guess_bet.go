package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

const (
	//投注订单id分布式生成用的redis key(常驻redis的key)
	ESPORT_REDIS_CREATE_BET_ORDER_ID = "redis_esport:create_bet_order_id"

	//投注风控相关redis key(常驻redis的key)
	ESPORT_REDIS_BET_RISK_CONTROL = "redis_esport:bet_risk_control"

	//游戏标签redis key(常驻redis的key)
	ESPORT_REDIS_GAME_LABEL = "redis_esport:game_label"
)

//数据库游戏标签
func GetDBGameLabel() *share_message.GameLabelRedisObj {

	//从数据库中取得数据
	dbGameLabelList := make([]*share_message.TableESPortsLabel, 0)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_LABEL)
	defer closeFun()

	var gameLabelRedisObj *share_message.GameLabelRedisObj
	err := col.Find(bson.M{"LabelType": ESPORTS_LABEL_TYPE_3}).All(&dbGameLabelList)

	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return gameLabelRedisObj
	}

	if err == mgo.ErrNotFound {
		return gameLabelRedisObj
	}

	if err == nil && nil != dbGameLabelList && len(dbGameLabelList) > 0 {
		gameLabelRedisObj = &share_message.GameLabelRedisObj{}
		for _, value := range dbGameLabelList {
			if value.GetLabelId() == int64(ESPORTS_LABEL_WZRY) {
				gameLabelRedisObj.WZRYIcon = easygo.NewString(value.GetIconUrl())
			} else if value.GetLabelId() == int64(ESPORTS_LABEL_DOTA2) {
				gameLabelRedisObj.DOTAIcon = easygo.NewString(value.GetIconUrl())
			} else if value.GetLabelId() == int64(ESPORTS_LABEL_LOL) {
				gameLabelRedisObj.LOLIcon = easygo.NewString(value.GetIconUrl())
			} else if value.GetLabelId() == int64(ESPORTS_LABEL_CSGO) {
				gameLabelRedisObj.CSGOIcon = easygo.NewString(value.GetIconUrl())
			} else if value.GetLabelId() == int64(ESPORTS_LABEL_OTHER) {
				gameLabelRedisObj.OTHERIcon = easygo.NewString(value.GetIconUrl())
			}
		}
	}

	return gameLabelRedisObj
}

//设置redis游戏标签
func SetRedisGameLabel() {

	gameLabelDbObj := GetDBGameLabel()

	if nil != gameLabelDbObj {
		//json序列化
		data, errJs := json.Marshal(gameLabelDbObj)
		easygo.PanicError(errJs)
		err := easygo.RedisMgr.GetC().Set(ESPORT_REDIS_GAME_LABEL, string(data))
		easygo.PanicError(err)
	} else {
		keys := []interface{}{ESPORT_REDIS_GAME_LABEL}

		easygo.RedisMgr.GetC().Delete(keys...)
	}
}

//取得redis游戏标签
func GetRedisGameLabel() *share_message.GameLabelRedisObj {

	b, err := easygo.RedisMgr.GetC().Exist(ESPORT_REDIS_GAME_LABEL)
	easygo.PanicError(err)

	var gameLabelRedisObj *share_message.GameLabelRedisObj
	if !b {
		gameLabelDbObj := GetDBGameLabel()

		if nil != gameLabelDbObj {
			//json序列化
			data, errJs := json.Marshal(gameLabelDbObj)
			easygo.PanicError(errJs)
			err := easygo.RedisMgr.GetC().Set(ESPORT_REDIS_GAME_LABEL, string(data))
			easygo.PanicError(err)

			gameLabelRedisObj = gameLabelDbObj
		}
	} else {
		obj := &share_message.GameLabelRedisObj{}
		value, err := easygo.RedisMgr.GetC().Get(ESPORT_REDIS_GAME_LABEL)
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), obj)
		easygo.PanicError(errJs)

		gameLabelRedisObj = obj
	}
	return gameLabelRedisObj
}

//设置redis投注风控相关
func SetRedisGuessBetRiskControl() {

	betRiskCtrlPara := QuerySysParameterById(ESPORT_PARAMETER)

	if nil != betRiskCtrlPara {
		redisBetRiskCtrlPara := &share_message.GameGuessBetRiskCtrlObj{
			EsOneBetGold:    easygo.NewInt64(betRiskCtrlPara.GetEsOneBetGold()),
			EsOneDayBetGold: easygo.NewInt64(betRiskCtrlPara.GetEsOneDayBetGold()),
			EsDaySumGold:    easygo.NewInt64(betRiskCtrlPara.GetEsDaySumGold()),
		}
		//json序列化
		data, errJs := json.Marshal(redisBetRiskCtrlPara)
		easygo.PanicError(errJs)
		err := easygo.RedisMgr.GetC().Set(ESPORT_REDIS_BET_RISK_CONTROL, string(data))
		easygo.PanicError(err)
	} else {
		keys := []interface{}{ESPORT_REDIS_BET_RISK_CONTROL}

		easygo.RedisMgr.GetC().Delete(keys...)
	}
}

//取得redis投注风控相关
func GetRedisGuessBetRiskControl() *share_message.GameGuessBetRiskCtrlObj {

	b, err := easygo.RedisMgr.GetC().Exist(ESPORT_REDIS_BET_RISK_CONTROL)
	easygo.PanicError(err)

	var betRiskCtrlObj *share_message.GameGuessBetRiskCtrlObj
	if !b {
		betRiskCtrlPara := QuerySysParameterById(ESPORT_PARAMETER)

		if nil != betRiskCtrlPara {
			redisBetRiskCtrlPara := &share_message.GameGuessBetRiskCtrlObj{
				EsOneBetGold:    easygo.NewInt64(betRiskCtrlPara.GetEsOneBetGold()),
				EsOneDayBetGold: easygo.NewInt64(betRiskCtrlPara.GetEsOneDayBetGold()),
				EsDaySumGold:    easygo.NewInt64(betRiskCtrlPara.GetEsDaySumGold()),
			}
			//json序列化
			data, errJs := json.Marshal(redisBetRiskCtrlPara)
			easygo.PanicError(errJs)
			err := easygo.RedisMgr.GetC().Set(ESPORT_REDIS_BET_RISK_CONTROL, string(data))
			easygo.PanicError(err)

			betRiskCtrlObj = redisBetRiskCtrlPara
		}

	} else {
		obj := &share_message.GameGuessBetRiskCtrlObj{}
		value, err := easygo.RedisMgr.GetC().Get(ESPORT_REDIS_BET_RISK_CONTROL)
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), obj)
		easygo.PanicError(errJs)

		betRiskCtrlObj = obj
	}
	return betRiskCtrlObj
}

func RedisCreateBetOrderID() int64 {
	b, err := easygo.RedisMgr.GetC().Exist(ESPORT_REDIS_CREATE_BET_ORDER_ID)
	easygo.PanicError(err)
	if !b {
		InitRedisCreateBetOrderId()
	}
	return easygo.RedisMgr.GetC().StringIncrForInt64(ESPORT_REDIS_CREATE_BET_ORDER_ID)
}

//投注订单id生成的初始化
func InitRedisCreateBetOrderId() {
	//取得投注表中的orderId
	betRecord := share_message.TableESPortsGuessBetRecord{}

	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_GUESS_BET_RECORD)
	defer closeFun()

	err := col.Find(bson.M{}).Sort("-_id").Limit(1).One(&betRecord)

	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	if err == mgo.ErrNotFound {
		err := easygo.RedisMgr.GetC().StringSet(ESPORT_REDIS_CREATE_BET_ORDER_ID, util.GetMilliTime()*1000)
		easygo.PanicError(err)
	} else if err == nil {
		err := easygo.RedisMgr.GetC().StringSet(ESPORT_REDIS_CREATE_BET_ORDER_ID, betRecord.GetOrderId())
		easygo.PanicError(err)
	}
}
